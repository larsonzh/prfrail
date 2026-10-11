package console

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/gates"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// This file wires the supervised-black-box assembly into the CLI: the `run`
// routing entry (called from executeRun) and the `prfrail supervised
// <return|report>` command group.

// unavailableVisibleClient is the fail-closed stand-in for the visible channel.
// The real SessionBridge visible IPC wiring is T027-B and is
// out of this slice, so the CLI refuses a supervised run at the delivery step
// with visible_delivery_unavailable instead of pretending the notice reached an
// operator UI.
type unavailableVisibleClient struct{}

var _ adapters.VisibleClient = unavailableVisibleClient{}

func (unavailableVisibleClient) SendVisible(context.Context, adapters.RequestEnvelope) (adapters.VisibleResponse, error) {
	return adapters.VisibleResponse{Status: "unavailable"}, adapters.ErrSessbridgeClientUnavailable
}

// executeSupervisedRun is the `run` routing entry of a supervised definition.
func (cli CLI) executeSupervisedRun(ctx context.Context, runID, runDir string, definition chain.Definition, jsonOutput bool, stdout, stderr io.Writer) int {
	options, err := cli.supervisedAssemblyOptions(runID, runDir, definition)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	summary, err := ExecuteSupervisedRun(ctx, options)
	if err != nil {
		return writeCommandError(stdout, stderr, "run", jsonOutput, err)
	}
	return writeRunSummary(stdout, jsonOutput, runDir, summary)
}

// supervisedAssemblyOptions builds the assembly options for a run root. The
// after-return gate runner is assembled here; hook definitions are left empty
// because the chain config only carries hook identifiers, so the return fails
// closed until a hook-definition source is wired in.
func (cli CLI) supervisedAssemblyOptions(runID, runDir string, definition chain.Definition) (SupervisedAssemblyOptions, error) {
	runner, err := cli.supervisedGateRunner(runDir)
	if err != nil {
		return SupervisedAssemblyOptions{}, err
	}
	return SupervisedAssemblyOptions{
		RunID:         runID,
		RunRoot:       runDir,
		Definition:    definition,
		VisibleSender: cli.visibleSender,
		GateRunner:    runner,
		Clock:         cli.now,
	}, nil
}

// supervisedGateRunner assembles the production gate runner over the run root.
func (cli CLI) supervisedGateRunner(runDir string) (*gates.Runner, error) {
	snapshotStore, err := snapshot.NewStore(filepath.Join(runDir, supervisedSnapshotStoreDirectory), 0)
	if err != nil {
		return nil, err
	}
	resultStore, err := gates.NewFileResultStore(filepath.Join(SupervisedDir(runDir), supervisedScansDirectoryName, "gate-results.jsonl"))
	if err != nil {
		return nil, err
	}
	return &gates.Runner{
		Executor: gates.GuardExecutor{},
		Scanner:  gates.DefaultSecretScanner(),
		Objects:  snapshotStore,
		Store:    resultStore,
		Clock:    cli.now,
	}, nil
}

// executeSupervised resolves the run root and definition of a supervised run.
func (cli CLI) executeSupervised(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, supervisedUsage)
		return exitUsage
	}
	switch args[0] {
	case "return":
		return cli.executeSupervisedReturn(ctx, args[1:], stdout, stderr)
	case "report":
		return cli.executeSupervisedReport(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown supervised subcommand: %s\n", args[0])
		fmt.Fprintln(stderr, supervisedUsage)
		return exitUsage
	}
}

const supervisedUsage = "usage: prfrail supervised <return|report> [options]"

func (cli CLI) executeSupervisedReturn(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("supervised return")
	chainPath := set.String("chain", "", "path to chain config file")
	runDir := set.String("run-dir", "", "run directory")
	runID := set.String("run-id", "", "run identifier")
	outcome := set.String("outcome", "", "return outcome: complete or abort")
	jsonOutput := set.Bool("json", false, "output as JSON")
	usage := "usage: prfrail supervised return --run-id <id> --outcome complete|abort [--run-dir <path>] [--chain <path>] [--json]"
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "supervised return", *jsonOutput, parseErr.String(), usage)
	}
	if len(set.Args()) != 0 || strings.TrimSpace(*runID) == "" || (*outcome != "complete" && *outcome != "abort") {
		return writeUsageError(stdout, stderr, "supervised return", *jsonOutput, "run-id and a complete|abort outcome are required", usage)
	}
	resolved, err := cli.resolveSupervisedRun(*chainPath, *runDir, *runID)
	if err != nil {
		return writeCommandError(stdout, stderr, "supervised return", *jsonOutput, err)
	}
	options, err := cli.supervisedAssemblyOptions(*runID, resolved.runDir, resolved.definition)
	if err != nil {
		return writeCommandError(stdout, stderr, "supervised return", *jsonOutput, err)
	}
	summary, err := ExecuteSupervisedReturn(ctx, options, *outcome)
	if err != nil {
		return writeCommandError(stdout, stderr, "supervised return", *jsonOutput, err)
	}
	return writeRunSummary(stdout, *jsonOutput, resolved.runDir, summary)
}

func (cli CLI) executeSupervisedReport(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("supervised report")
	runDir := set.String("run-dir", "", "run directory")
	jsonOutput := set.Bool("json", false, "output as JSON")
	usage := "usage: prfrail supervised report --run-dir <path> [--json]"
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "supervised report", *jsonOutput, parseErr.String(), usage)
	}
	if len(set.Args()) != 0 || strings.TrimSpace(*runDir) == "" {
		return writeUsageError(stdout, stderr, "supervised report", *jsonOutput, "supervised report requires --run-dir", usage)
	}
	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "supervised report", *jsonOutput, err)
	}
	resolvedRunDir, err := resolvePathFromCWD(cwd, *runDir)
	if err != nil {
		return writeCommandError(stdout, stderr, "supervised report", *jsonOutput, err)
	}
	path := SupervisedReportPath(resolvedRunDir)
	var record chain.SupervisedReportRecord
	found, err := readSupervisedJSON(path, &record)
	if err != nil {
		return writeCommandError(stdout, stderr, "supervised report", *jsonOutput, err)
	}
	if !found {
		return writeCommandError(stdout, stderr, "supervised report", *jsonOutput, fmt.Errorf("no supervised report at %s", path))
	}
	if err := chain.ValidateSupervisedReport(record); err != nil {
		return writeCommandError(stdout, stderr, "supervised report", *jsonOutput, err)
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "supervised report", OK: true, ExitCode: exitSuccess, Data: record})
	}
	writeSupervisedReportText(stdout, path, record)
	return exitSuccess
}

// supervisedRunResolution is the resolved run root and definition of a supervised run.
type supervisedRunResolution struct {
	runDir     string
	definition chain.Definition
}

func (cli CLI) resolveSupervisedRun(chainPath, runDir, runID string) (supervisedRunResolution, error) {
	cwd, err := cli.getwd()
	if err != nil {
		return supervisedRunResolution{}, err
	}
	resolvedPath, err := ResolveChainConfigPath(chainPath, cwd)
	if err != nil {
		return supervisedRunResolution{}, err
	}
	cfg, _, err := LoadChainConfig(resolvedPath)
	if err != nil {
		return supervisedRunResolution{}, err
	}
	definition, err := cfg.ToDefinition()
	if err != nil {
		return supervisedRunResolution{}, err
	}
	resolvedRunDir := runDir
	if strings.TrimSpace(resolvedRunDir) == "" {
		resolvedRunDir = filepath.Join(filepath.Dir(resolvedPath), "tmp", "prfrail-runs", runID)
	}
	resolvedRunDir, err = resolvePathFromCWD(cwd, resolvedRunDir)
	if err != nil {
		return supervisedRunResolution{}, err
	}
	return supervisedRunResolution{runDir: resolvedRunDir, definition: definition}, nil
}

// writeSupervisedReportText renders the reduced/unknown report. Unknown
// dimensions are rendered literally as `unknown` and the render never claims
// that any dimension is safe or verified (P4).
func writeSupervisedReportText(stdout io.Writer, path string, record chain.SupervisedReportRecord) {
	body := record.Record
	fmt.Fprintf(stdout, "supervised report: %s\n", path)
	fmt.Fprintf(stdout, "runId: %s\n", body.RunID)
	fmt.Fprintf(stdout, "taskId: %s\n", body.TaskID)
	fmt.Fprintf(stdout, "stepId: %s\n", body.StepID)
	fmt.Fprintf(stdout, "attempt: %d\n", body.Attempt)
	fmt.Fprintf(stdout, "mode: %s\n", body.Mode)
	fmt.Fprintf(stdout, "assuranceLevel: %s\n", body.AssuranceLevel)
	for _, dimension := range body.VerifiedDimensions {
		fmt.Fprintf(stdout, "proved: %s\n", dimension)
	}
	for _, dimension := range body.UnknownDimensions {
		fmt.Fprintf(stdout, "unknown: %s=unknown\n", dimension)
	}
}
