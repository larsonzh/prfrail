package console

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/gates"
	"github.com/larsonzh/prfrail/internal/tickets"
)

const (
	exitSuccess = 0
	exitFailure = 1
	exitUsage   = 2
)

type CLI struct {
	version             string
	now                 func() time.Time
	getwd               func() (string, error)
	stdin               io.Reader
	stopRun             func(context.Context, string, string) ([]string, error)
	appendInteraction   func(string, chain.OperatorInteractionRecord) error
	secretManager       adapters.AISecretManager
	readSecret          func() ([]byte, error)
	checkAIAvailability func(context.Context, adapters.AIProviderProfile, string, int, string, string) (adapters.AIAvailabilityRecord, error)
	// visibleSender is the visible channel of the supervised-black-box assembly.
	// It defaults to the fail-closed placeholder until the real SessionBridge
	// visible wiring lands (T027-B).
	visibleSender adapters.VisibleClient
}

type commandResponse struct {
	Command  string `json:"command"`
	OK       bool   `json:"ok"`
	ExitCode int    `json:"exitCode"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
	Usage    string `json:"usage,omitempty"`
	Data     any    `json:"data,omitempty"`
}

func NewCLI(version string) CLI {
	cli := CLI{
		version:       version,
		now:           time.Now,
		getwd:         os.Getwd,
		stdin:         os.Stdin,
		secretManager: adapters.NewPlatformAISecretManager(),
		readSecret:    readTerminalSecret,
		visibleSender: unavailableVisibleClient{},
	}
	cli.stopRun = newDefaultStopRun(cli.getwd)
	cli.appendInteraction = appendOperatorInteractionRecord
	cli.checkAIAvailability = cli.defaultAICheck
	return cli
}

func (cli CLI) Execute(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		writeMainUsage(stderr)
		return exitUsage
	}
	switch args[0] {
	case "help", "--help", "-h":
		writeMainUsage(stdout)
		return exitSuccess
	case "version":
		return cli.executeVersion(args[1:], stdout, stderr)
	case "init":
		return cli.executeInit(args[1:], stdout, stderr)
	case "validate":
		return cli.executeValidate(args[1:], stdout, stderr)
	case "preview":
		return cli.executePreview(args[1:], stdout, stderr)
	case "approvals":
		return cli.executeApprovals(ctx, args[1:], stdout, stderr)
	case "interactions":
		return cli.executeInteractions(args[1:], stdout, stderr)
	case "cost":
		return cli.executeCost(args[1:], stdout, stderr)
	case "secret":
		return cli.executeSecret(ctx, args[1:], stdout, stderr)
	case "ai":
		return cli.executeAI(ctx, args[1:], stdout, stderr)
	case "run":
		return cli.executeRun(ctx, args[1:], stdout, stderr)
	case "supervised":
		return cli.executeSupervised(ctx, args[1:], stdout, stderr)
	case "report":
		return cli.executeReport(ctx, args[1:], stdout, stderr)
	case "config":
		if len(args) < 2 || args[1] != "explain" {
			usage := "usage: prfrail config explain [--chain <path>] [--json]"
			fmt.Fprintln(stderr, usage)
			return exitUsage
		}
		return cli.executeConfigExplain(args[2:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		writeMainUsage(stderr)
		return exitUsage
	}
}

func (cli CLI) executeVersion(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("version")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "version", *jsonOutput, parseErr.String(), "usage: prfrail version [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "version", *jsonOutput, "version does not accept positional arguments", "usage: prfrail version [--json]")
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "version", OK: true, ExitCode: exitSuccess, Data: map[string]string{"version": cli.version}})
	}
	fmt.Fprintf(stdout, "prfrail %s\n", cli.version)
	return exitSuccess
}

func (cli CLI) executeInit(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("init")
	workspacePath := set.String("workspace", ".", "workspace root")
	chainPath := set.String("chain", "", "path to chain config file")
	force := set.Bool("force", false, "overwrite output config file")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "init", *jsonOutput, parseErr.String(), "usage: prfrail init [--workspace <path>] [--chain <path>] [--force] [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "init", *jsonOutput, "init does not accept positional arguments", "usage: prfrail init [--workspace <path>] [--chain <path>] [--force] [--json]")
	}

	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "init", *jsonOutput, err)
	}
	workspaceAbs, err := resolvePathFromCWD(cwd, *workspacePath)
	if err != nil {
		return writeCommandError(stdout, stderr, "init", *jsonOutput, err)
	}
	info, err := os.Stat(workspaceAbs)
	if err != nil {
		return writeCommandError(stdout, stderr, "init", *jsonOutput, err)
	}
	if !info.IsDir() {
		return writeCommandError(stdout, stderr, "init", *jsonOutput, fmt.Errorf("workspace path %q is not a directory", workspaceAbs))
	}

	outputPath := *chainPath
	if outputPath == "" {
		outputPath = filepath.Join(workspaceAbs, DefaultChainConfigName)
	}
	outputPath, err = resolvePathFromCWD(cwd, outputPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "init", *jsonOutput, err)
	}

	chainID := NormalizeChainID(filepath.Base(workspaceAbs))
	cfg := DefaultChainConfig(chainID)
	content, err := EncodeChainConfig(cfg)
	if err != nil {
		return writeCommandError(stdout, stderr, "init", *jsonOutput, err)
	}
	if err := writeFile(outputPath, content, *force); err != nil {
		return writeCommandError(stdout, stderr, "init", *jsonOutput, err)
	}

	data := map[string]any{"chainPath": outputPath, "chainId": cfg.Chain.ID, "profile": cfg.Chain.Profile}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "init", OK: true, ExitCode: exitSuccess, Message: "template chain config created", Data: data})
	}
	fmt.Fprintf(stdout, "init: wrote %s\n", outputPath)
	fmt.Fprintf(stdout, "chain: %s profile=%s\n", cfg.Chain.ID, cfg.Chain.Profile)
	return exitSuccess
}

func (cli CLI) executeValidate(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("validate")
	chainPath := set.String("chain", "", "path to chain config file")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "validate", *jsonOutput, parseErr.String(), "usage: prfrail validate [--chain <path>] [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "validate", *jsonOutput, "validate does not accept positional arguments", "usage: prfrail validate [--chain <path>] [--json]")
	}

	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "validate", *jsonOutput, err)
	}
	resolvedPath, err := ResolveChainConfigPath(*chainPath, cwd)
	if err != nil {
		return writeCommandError(stdout, stderr, "validate", *jsonOutput, err)
	}
	cfg, _, err := LoadChainConfig(resolvedPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "validate", *jsonOutput, err)
	}
	summary := ValidateSummary(resolvedPath, cfg)
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "validate", OK: true, ExitCode: exitSuccess, Message: "configuration is valid", Data: summary})
	}
	fmt.Fprintln(stdout, "validate: ok")
	fmt.Fprintf(stdout, "config: %s\n", summary.ConfigPath)
	fmt.Fprintf(stdout, "chain: %s tasks=%d steps=%d\n", summary.ChainID, summary.TaskCount, summary.StepCount)
	fmt.Fprintf(stdout, "runnableInCli: %t\n", summary.RunnableInCLI)
	for _, warning := range summary.Warnings {
		fmt.Fprintf(stdout, "warning: %s\n", warning)
	}
	return exitSuccess
}

func (cli CLI) executeConfigExplain(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("config explain")
	chainPath := set.String("chain", "", "path to chain config file")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "config explain", *jsonOutput, parseErr.String(), "usage: prfrail config explain [--chain <path>] [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "config explain", *jsonOutput, "config explain does not accept positional arguments", "usage: prfrail config explain [--chain <path>] [--json]")
	}

	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "config explain", *jsonOutput, err)
	}
	resolvedPath, err := ResolveChainConfigPath(*chainPath, cwd)
	if err != nil {
		return writeCommandError(stdout, stderr, "config explain", *jsonOutput, err)
	}
	cfg, rawConfig, err := LoadChainConfig(resolvedPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "config explain", *jsonOutput, err)
	}
	report := BuildConfigExplain(resolvedPath, cfg, rawConfig)
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "config explain", OK: true, ExitCode: exitSuccess, Data: report})
	}
	fmt.Fprintf(stdout, "config: %s\n", report.ConfigPath)
	fmt.Fprintf(stdout, "chain: %s profile=%s\n", report.ChainID, report.Profile)
	fmt.Fprintf(stdout, "documentationPolicy: %s\n", report.DocumentationPolicy)
	fmt.Fprintf(stdout, "runnableInCli: %t\n", report.RunnableInCLI)
	for _, profile := range report.AIProfiles {
		fmt.Fprintf(stdout, "aiProfile: %s hash=%s\n", profile.ProfileID, profile.ProfileConfigHash)
	}
	for _, binding := range report.AIChannels {
		fmt.Fprintf(stdout, "aiChannel: %s profile=%s hash=%s\n", binding.Channel, binding.ProfileID, binding.ProfileConfigHash)
	}
	for _, policy := range report.TaskPolicies {
		source := policy.SourceKind
		if policy.SourcePointer != nil {
			source = source + " " + *policy.SourcePointer
		}
		fmt.Fprintf(stdout, "taskPolicy: %s=%s (%s)\n", policy.TaskID, policy.Policy, source)
	}
	return exitSuccess
}

func (cli CLI) executeRun(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("run")
	chainPath := set.String("chain", "", "path to chain config file")
	runDir := set.String("run-dir", "", "run directory")
	runID := set.String("run-id", "", "run identifier")
	jsonOutput := set.Bool("json", false, "output as JSON")
	progressOutput := set.Bool("progress", false, "stream per-task and per-step progress to stderr")
	requestRecord := set.String("request-record", "", "path to a declared AgentRunner request record (enables the AgentRunner assembly)")
	requestRecordPaths := make([]string, 0)
	set.Func("request-records", "path to a declared AgentRunner request record (repeatable; one per task/step)", func(value string) error {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("request record path is empty")
		}
		requestRecordPaths = append(requestRecordPaths, value)
		return nil
	})
	capabilityRecord := set.String("capability-record", "", "path to a declared AgentRunner capability record")
	enforcementRecord := set.String("enforcement-record", "", "path to a declared AgentRunner enforcement record")
	availabilityRecord := set.String("availability-record", "", "path to a declared AI availability record")
	authorizationLedger := set.String("authorization-ledger", "", "path to an authorization ledger")
	costLedger := set.String("cost-ledger", "", "path to a cost ledger")
	agentCLIPath := set.String("agent-cli", "", "path to the pinned agent CLI executable")
	agentVersion := set.String("agent-version", "", "pinned agent CLI version")
	probePermit := set.String("probe-permit", "", "path to the one-shot probe permit marker (required for the AgentRunner and supervised assembly routes)")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "run", *jsonOutput, parseErr.String(), runUsage)
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "run", *jsonOutput, "run does not accept positional arguments", runUsage)
	}

	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}
	resolvedPath, err := ResolveChainConfigPath(*chainPath, cwd)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}
	cfg, _, err := LoadChainConfig(resolvedPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}
	executable := FindExecutableSteps(cfg)
	supervised := FindSupervisedSteps(cfg)
	inputs := agentRunnerRunInputs{
		requestRecord:       *requestRecord,
		requestRecords:      requestRecordPaths,
		capabilityRecord:    *capabilityRecord,
		enforcementRecord:   *enforcementRecord,
		availabilityRecord:  *availabilityRecord,
		authorizationLedger: *authorizationLedger,
		costLedger:          *costLedger,
		agentCLI:            *agentCLIPath,
		agentVersion:        *agentVersion,
	}
	if strings.TrimSpace(inputs.requestRecord) != "" && len(inputs.requestRecords) > 0 {
		return writeUsageError(stdout, stderr, "run", *jsonOutput, "--request-record and --request-records are mutually exclusive", runUsage)
	}
	// A chain without executable steps never enters the AgentRunner assembly, so
	// a per-task registry declared against it has no consumer. Refuse instead of
	// running the noop path and silently dropping the declaration.
	if len(inputs.requestRecords) > 0 && len(executable) == 0 && len(supervised) == 0 {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, fmt.Errorf("request records require at least one executable step in the chain definition"))
	}
	// P1: a supervised definition can never be entered by supplying AgentRunner
	// input flags, and it never falls back to the AgentRunner assembly.
	if len(supervised) > 0 && inputs.anyProvided() {
		err := fmt.Errorf("supervised-black-box definitions are mutually exclusive with AgentRunner input flags")
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}
	if len(executable) > 0 && !inputs.complete() && len(supervised) == 0 {
		first := executable[0]
		err := fmt.Errorf("local runtime only supports noop steps; first executable step is %s/%s (%s)", first.TaskID, first.StepID, first.Kind)
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}

	if *runID == "" {
		*runID = defaultRunID(cli.now())
	}
	if !evidence.ValidID(*runID) {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, fmt.Errorf("run-id %q is invalid", *runID))
	}

	resolvedRunDir := *runDir
	if resolvedRunDir == "" {
		resolvedRunDir = filepath.Join(filepath.Dir(resolvedPath), "tmp", "prfrail-runs", *runID)
	}
	resolvedRunDir, err = resolvePathFromCWD(cwd, resolvedRunDir)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}
	eventLogPath := runEventLogPath(resolvedRunDir)
	if _, err := os.Stat(eventLogPath); err == nil {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, fmt.Errorf("event log already exists at %q", eventLogPath))
	} else if !errors.Is(err, os.ErrNotExist) {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}

	definition, err := cfg.ToDefinition()
	if err != nil {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}
	// The probe permit gate is a CLI-side presence check only. It resolves and
	// loads the marker so a bad path or an unreadable marker fails here, but it
	// never consumes the quota: consumption belongs to gates.ProbePermit.Require
	// at the assembly entry. A permit that passes this check can still be spent,
	// and the assembly entry is where that is decided.
	var gatePermit *gates.ProbePermit
	if len(supervised) > 0 || len(executable) > 0 {
		if strings.TrimSpace(*probePermit) == "" {
			err := fmt.Errorf("--probe-permit is required when the chain routes to the AgentRunner or supervised assembly")
			return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
		}
		permitPath, err := resolvePathFromCWD(cwd, *probePermit)
		if err != nil {
			return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
		}
		permit, err := gates.LoadProbePermit(permitPath)
		if err != nil {
			return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
		}
		if _, err := gates.EvaluateT027BReadiness(gates.GoalProbe, gates.ReadinessInput{Permit: permit}); err != nil {
			return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
		}
		gatePermit = permit
	}
	// The progress writer is resolved once, before routing, so all three routes
	// receive the same nil-or-stderr decision. A nil writer keeps every existing
	// output byte-identical: no route writes progress when the flag is absent.
	var progressWriter io.Writer
	if *progressOutput {
		progressWriter = stderr
	}
	if len(supervised) > 0 {
		return cli.executeSupervisedRun(ctx, *runID, resolvedRunDir, definition, gatePermit, *jsonOutput, stdout, stderr, progressWriter)
	}
	if len(executable) > 0 {
		return cli.executeAgentRunnerRun(ctx, *runID, resolvedRunDir, definition, inputs, gatePermit, *jsonOutput, stdout, stderr, progressWriter)
	}
	summary, err := ExecuteNoopRun(ctx, *runID, resolvedRunDir, definition, filepath.Dir(resolvedPath), cli.now, progressWriter)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}
	return writeRunSummary(stdout, *jsonOutput, resolvedRunDir, summary)
}

const runUsage = "usage: prfrail run [--chain <path>] [--run-id <id>] [--run-dir <path>] [--json] [--progress] (--request-record <path> | --request-records <path>)... --capability-record <path> --enforcement-record <path> --availability-record <path> --authorization-ledger <path> --cost-ledger <path> --agent-cli <path> --agent-version <version> --probe-permit <path>"

// agentRunnerRunInputs is the all-or-nothing AgentRunner CLI surface. Before the
// run command routes to the AgentRunner assembly the request-record dimension
// must be supplied in exactly one of its two forms, the single-value
// requestRecord or the repeatable requestRecords, and all seven remaining values
// must be present. A partially supplied set keeps the existing noop-only refusal
// instead of assembling from an incomplete declaration.
type agentRunnerRunInputs struct {
	requestRecord       string
	requestRecords      []string
	capabilityRecord    string
	enforcementRecord   string
	availabilityRecord  string
	authorizationLedger string
	costLedger          string
	agentCLI            string
	agentVersion        string
}

func (inputs agentRunnerRunInputs) complete() bool {
	if strings.TrimSpace(inputs.requestRecord) == "" && len(inputs.requestRecords) == 0 {
		return false
	}
	for _, value := range []string{
		inputs.capabilityRecord,
		inputs.enforcementRecord,
		inputs.availabilityRecord,
		inputs.authorizationLedger,
		inputs.costLedger,
		inputs.agentCLI,
		inputs.agentVersion,
	} {
		if strings.TrimSpace(value) == "" {
			return false
		}
	}
	return true
}

// anyProvided reports whether at least one AgentRunner input flag was supplied.
// It is the hazard half of the supervised/AgentRunner mutual exclusion.
func (inputs agentRunnerRunInputs) anyProvided() bool {
	if len(inputs.requestRecords) > 0 {
		return true
	}
	for _, value := range []string{
		inputs.requestRecord,
		inputs.capabilityRecord,
		inputs.enforcementRecord,
		inputs.availabilityRecord,
		inputs.authorizationLedger,
		inputs.costLedger,
		inputs.agentCLI,
		inputs.agentVersion,
	} {
		if strings.TrimSpace(value) != "" {
			return true
		}
	}
	return false
}

// executeAgentRunnerRun decodes the declared record files and drives the
// production assembly. Decoding is CLI-side file IO; the assembly re-validates
// every record before the engine exists.
func (cli CLI) executeAgentRunnerRun(ctx context.Context, runID, runDir string, definition chain.Definition, inputs agentRunnerRunInputs, probePermit *gates.ProbePermit, jsonOutput bool, stdout, stderr io.Writer, progress io.Writer) int {
	// probePermit is only forwarded: the assembly entry evaluates and consumes it.
	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	readRecord := func(label, value string) ([]byte, error) {
		path, err := resolvePathFromCWD(cwd, value)
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", label, err)
		}
		return data, nil
	}
	var (
		requestRecord  adapters.AgentRunnerRequestRecord
		requestRecords map[string]adapters.AgentRunnerRequestRecord
	)
	if len(inputs.requestRecords) > 0 {
		// The repeatable form owns the whole registry: every planned task must be
		// declared here, and the single compatibility record stays zero so no task
		// can silently borrow it.
		requestRecords, err = readAgentRunnerRequestRecords(cwd, runID, definition, inputs.requestRecords)
		if err != nil {
			return writeCommandError(stdout, stderr, "run", jsonOutput, err)
		}
	} else {
		requestBytes, err := readRecord("request record", inputs.requestRecord)
		if err != nil {
			return writeCommandError(stdout, stderr, "run", jsonOutput, err)
		}
		requestRecord, err = adapters.DecodeAgentRunnerRequestRecord(requestBytes)
		if err != nil {
			return writeCommandError(stdout, stderr, "run", jsonOutput, err)
		}
	}
	capabilityBytes, err := readRecord("capability record", inputs.capabilityRecord)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	capabilityRecord, err := adapters.DecodeAgentRunnerCapabilityRecord(capabilityBytes)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	enforcementBytes, err := readRecord("enforcement record", inputs.enforcementRecord)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	enforcementRecord, err := adapters.DecodeAgentRunnerEnforcementRecord(enforcementBytes)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	availabilityBytes, err := readRecord("availability record", inputs.availabilityRecord)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	availabilityRecord, err := adapters.DecodeAIAvailabilityRecord(availabilityBytes)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	ledgerPath, err := resolvePathFromCWD(cwd, inputs.authorizationLedger)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	authorizationLedger, err := loadAuthorizationLedger(ledgerPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	costLedgerPath, err := resolvePathFromCWD(cwd, inputs.costLedger)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	costLedgerBytes, err := os.ReadFile(costLedgerPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	costLedgerRecord, err := tickets.DecodeCostLedgerRecord(costLedgerBytes)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	costLedger, err := tickets.LoadCostLedger(costLedgerRecord)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	agentCLIPath, err := resolvePathFromCWD(cwd, inputs.agentCLI)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	options := AgentRunnerAssemblyOptions{
		RunID:               runID,
		RunRoot:             runDir,
		Definition:          definition,
		ProbePermit:         probePermit,
		RequestRecord:       requestRecord,
		RequestRecords:      requestRecords,
		CapabilityRecord:    capabilityRecord,
		EnforcementRecord:   &enforcementRecord,
		AvailabilityRecord:  availabilityRecord,
		AvailabilityPolicy:  AgentRunnerDefaultAvailabilityPolicy(availabilityRecord, cli.now()),
		AuthorizationLedger: authorizationLedger,
		CostLedger:          costLedger,
		AgentCLI: AgentRunnerCLIOptions{
			Executable:  agentCLIPath,
			Version:     inputs.agentVersion,
			VersionArgs: []string{agentRunnerDefaultVersionFlag},
			Timeout:     agentRunnerDefaultRunTimeout,
			MaxLogBytes: agentRunnerDefaultMaxLogBytes,
		},
		Clock:    cli.now,
		Progress: progress,
	}
	summary, err := ExecuteAgentRunnerRun(ctx, options)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	return writeRunSummary(stdout, jsonOutput, runDir, summary)
}

// readAgentRunnerRequestRecords decodes the repeatable --request-records files
// into the per-task registry of the AgentRunner assembly. Every record is read,
// decoded, bound to the run and keyed by agentRunnerRequestRecordKey; the
// registry must then cover exactly the task/step pairs the assembly plan binds.
// The order of the refusals is fixed so the CLI failure is deterministic: the
// run binding and the duplicate key are judged per record while decoding, the
// key/definition match after every record is decoded, and the coverage of the
// definition last.
func readAgentRunnerRequestRecords(cwd, runID string, definition chain.Definition, paths []string) (map[string]adapters.AgentRunnerRequestRecord, error) {
	records := make(map[string]adapters.AgentRunnerRequestRecord, len(paths))
	// keys keeps the decoded keys in input order, so the stray-key refusal below
	// reports a deterministic key instead of a random map iteration order.
	keys := make([]string, 0, len(paths))
	for _, value := range paths {
		path, err := resolvePathFromCWD(cwd, value)
		if err != nil {
			return nil, fmt.Errorf("request record: %w", err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("request record: %w", err)
		}
		record, err := adapters.DecodeAgentRunnerRequestRecord(data)
		if err != nil {
			return nil, err
		}
		if record.Request.RunID != runID {
			return nil, fmt.Errorf("request record runId %q does not match run id %q", record.Request.RunID, runID)
		}
		key := agentRunnerRequestRecordKey(record.Request.TaskID, record.Request.StepID)
		if _, duplicate := records[key]; duplicate {
			return nil, fmt.Errorf("request records must bind distinct task/step keys: duplicate key %q", key)
		}
		records[key] = record
		keys = append(keys, key)
	}
	bindings := agentRunnerExecutableBindings(definition)
	if len(bindings) == 0 {
		return nil, fmt.Errorf("request records require at least one executable step in the chain definition")
	}
	planned := make(map[string]struct{}, len(bindings))
	for _, binding := range bindings {
		planned[agentRunnerRequestRecordKey(binding.TaskID, binding.StepID)] = struct{}{}
	}
	for _, key := range keys {
		if _, ok := planned[key]; !ok {
			return nil, fmt.Errorf("request record does not bind a task/step of the chain definition: %q", key)
		}
	}
	for _, binding := range bindings {
		if _, ok := records[agentRunnerRequestRecordKey(binding.TaskID, binding.StepID)]; !ok {
			return nil, fmt.Errorf("request record is missing for task %q step %q", binding.TaskID, binding.StepID)
		}
	}
	return records, nil
}

// agentRunnerExecutableBindings enumerates the task/step pairs of the definition
// that the AgentRunner assembly binds. The enumeration mirrors
// agentRunnerAssemblyPlan: a task contributes a binding only when it carries
// exactly one step and that step is the isolated-workspace code step the
// AgentRunner route can prove. A definition the plan refuses is left to the
// assembly's own refusal; the registry only has to cover what the plan binds.
func agentRunnerExecutableBindings(definition chain.Definition) []agentRunnerTaskBinding {
	bindings := make([]agentRunnerTaskBinding, 0, len(definition.Tasks))
	for _, task := range definition.Tasks {
		if len(task.Steps) != 1 {
			continue
		}
		step := task.Steps[0]
		if step.Kind != "code" || step.Mode != chain.IsolatedWorkspace {
			continue
		}
		bindings = append(bindings, agentRunnerTaskBinding{TaskID: task.ID, StepID: step.ID})
	}
	return bindings
}

func writeRunSummary(stdout io.Writer, jsonOutput bool, runDir string, summary RunSummary) int {
	if jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "run", OK: true, ExitCode: exitSuccess, Message: "run completed", Data: summary})
	}
	fmt.Fprintf(stdout, "run: completed\n")
	fmt.Fprintf(stdout, "runId: %s\n", summary.RunID)
	fmt.Fprintf(stdout, "runDir: %s\n", runDir)
	fmt.Fprintf(stdout, "chainState: %s\n", summary.ChainState)
	fmt.Fprintf(stdout, "events: %d\n", summary.Sequence)
	fmt.Fprintf(stdout, "eventLog: %s\n", summary.EventLogPath)
	return exitSuccess
}

func (cli CLI) executeReport(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("report")
	runDir := set.String("run-dir", "", "run directory")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "report", *jsonOutput, parseErr.String(), "usage: prfrail report --run-dir <path> [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "report", *jsonOutput, "report does not accept positional arguments", "usage: prfrail report --run-dir <path> [--json]")
	}
	if strings.TrimSpace(*runDir) == "" {
		return writeUsageError(stdout, stderr, "report", *jsonOutput, "report requires --run-dir", "usage: prfrail report --run-dir <path> [--json]")
	}

	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "report", *jsonOutput, err)
	}
	resolvedRunDir, err := resolvePathFromCWD(cwd, *runDir)
	if err != nil {
		return writeCommandError(stdout, stderr, "report", *jsonOutput, err)
	}
	summary, err := LoadRunSummary(ctx, resolvedRunDir)
	if err != nil {
		return writeCommandError(stdout, stderr, "report", *jsonOutput, err)
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "report", OK: true, ExitCode: exitSuccess, Data: summary})
	}
	fmt.Fprintf(stdout, "report: %s\n", summary.RunID)
	fmt.Fprintf(stdout, "chainState: %s\n", summary.ChainState)
	fmt.Fprintf(stdout, "events: %d\n", summary.Sequence)
	fmt.Fprintf(stdout, "eventLog: %s\n", summary.EventLogPath)
	for _, item := range summary.TaskStateCounts {
		fmt.Fprintf(stdout, "taskState: %s=%d\n", item.State, item.Count)
	}
	for _, item := range summary.StepStateCounts {
		fmt.Fprintf(stdout, "stepState: %s=%d\n", item.State, item.Count)
	}
	return exitSuccess
}

func defaultRunID(now time.Time) string {
	return "run-" + now.UTC().Format("20060102-150405")
}

func resolvePathFromCWD(cwd, value string) (string, error) {
	if filepath.IsAbs(value) {
		return filepath.Clean(value), nil
	}
	return filepath.Abs(filepath.Join(cwd, value))
}

func writeMainUsage(output io.Writer) {
	fmt.Fprintln(output, "usage: prfrail <command> [options]")
	fmt.Fprintln(output, "commands: version, init, validate, preview, approvals, interactions, cost, secret, ai, run, report, config explain")
}

func writeFile(path string, content []byte, force bool) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	flags := os.O_CREATE | os.O_WRONLY
	if force {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}
	file, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(content); err != nil {
		return err
	}
	return file.Sync()
}

func writeUsageError(stdout, stderr io.Writer, command string, jsonOutput bool, errText string, usage string) int {
	errText = strings.TrimSpace(errText)
	if errText == "" {
		errText = "invalid arguments"
	}
	if jsonOutput {
		_ = writeJSONResponse(stdout, commandResponse{Command: command, OK: false, ExitCode: exitUsage, Error: errText, Usage: usage})
		return exitUsage
	}
	fmt.Fprintf(stderr, "%s: %s\n", command, errText)
	fmt.Fprintln(stderr, usage)
	return exitUsage
}

func writeCommandError(stdout, stderr io.Writer, command string, jsonOutput bool, err error) int {
	if jsonOutput {
		_ = writeJSONResponse(stdout, commandResponse{Command: command, OK: false, ExitCode: exitFailure, Error: err.Error()})
		return exitFailure
	}
	fmt.Fprintf(stderr, "%s: %v\n", command, err)
	return exitFailure
}

func writeJSONResponse(output io.Writer, response commandResponse) int {
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(response)
	return response.ExitCode
}

func newFlagSet(name string) (*flag.FlagSet, *bytes.Buffer) {
	buffer := &bytes.Buffer{}
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	set.SetOutput(buffer)
	return set, buffer
}
