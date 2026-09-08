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

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	exitSuccess = 0
	exitFailure = 1
	exitUsage   = 2
)

type CLI struct {
	version string
	now     func() time.Time
	getwd   func() (string, error)
	stopRun func(context.Context, string, string) ([]string, error)
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
	cli := CLI{version: version, now: time.Now, getwd: os.Getwd}
	cli.stopRun = newDefaultStopRun(cli.getwd)
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
	case "cost":
		return cli.executeCost(args[1:], stdout, stderr)
	case "run":
		return cli.executeRun(ctx, args[1:], stdout, stderr)
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
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "run", *jsonOutput, parseErr.String(), "usage: prfrail run [--chain <path>] [--run-id <id>] [--run-dir <path>] [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "run", *jsonOutput, "run does not accept positional arguments", "usage: prfrail run [--chain <path>] [--run-id <id>] [--run-dir <path>] [--json]")
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
	if len(executable) > 0 {
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
	summary, err := ExecuteNoopRun(ctx, *runID, resolvedRunDir, definition, filepath.Dir(resolvedPath), cli.now)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", *jsonOutput, err)
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "run", OK: true, ExitCode: exitSuccess, Message: "run completed", Data: summary})
	}
	fmt.Fprintf(stdout, "run: completed\n")
	fmt.Fprintf(stdout, "runId: %s\n", summary.RunID)
	fmt.Fprintf(stdout, "runDir: %s\n", resolvedRunDir)
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
	fmt.Fprintln(output, "commands: version, init, validate, preview, approvals, cost, run, report, config explain")
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
