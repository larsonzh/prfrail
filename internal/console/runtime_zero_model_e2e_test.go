package console

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
)

// P45 zero-model E2E on the console composition root.
//
// ExecuteNoopRun is the only console path that builds a whole chain without the
// AgentRunner assembly: localRuntime serves every port itself, so a green run
// proves the console can drive T1..T3 with no model, no credential and no external
// process.

func zeroModelClock() time.Time {
	return time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
}

func zeroModelNoopTask(id string) chain.Task {
	return chain.Task{ID: id, Steps: []chain.Step{{ID: id + "-noop", Kind: "noop", Reason: "zero-model-e2e"}}}
}

func zeroModelTaskCount(summary RunSummary, state string) int {
	total := 0
	for _, counter := range summary.TaskStateCounts {
		if counter.State == state {
			total += counter.Count
		}
	}
	return total
}

// zeroModelRunProjection reloads the persisted event log of a finished run and
// returns its projection plus every task ID the log mentions.
func zeroModelRunProjection(t *testing.T, runDir string) (chain.Projection, []string) {
	t.Helper()
	store, err := chain.NewFileEventStore(runEventLogPath(runDir))
	if err != nil {
		t.Fatalf("open the run event log: %v", err)
	}
	records, err := store.Load(context.Background())
	if err != nil {
		t.Fatalf("load the run event log: %v", err)
	}
	projection, err := chain.Rebuild(records)
	if err != nil {
		t.Fatalf("rebuild the run projection: %v", err)
	}
	taskIDs := make([]string, 0, len(records))
	for _, record := range records {
		if record.Event.Entity.TaskID != "" {
			taskIDs = append(taskIDs, record.Event.Entity.TaskID)
		}
	}
	return projection, taskIDs
}

// TestZeroModelConsoleNoopThreeTaskE2E is the console positive leg: three noop
// tasks reach COMPLETED with every task PASSED.
func TestZeroModelConsoleNoopThreeTaskE2E(t *testing.T) {
	runDir := t.TempDir()
	workspace := t.TempDir()
	definition := chain.Definition{ID: "chain-zero-model-noop", Tasks: []chain.Task{
		zeroModelNoopTask("task-one"),
		zeroModelNoopTask("task-two"),
		zeroModelNoopTask("task-three"),
	}}
	now := zeroModelClock

	summary, err := ExecuteNoopRun(context.Background(), "run-zero-model-noop", runDir, definition, workspace, now)
	if err != nil {
		t.Fatalf("the noop console path must run without any model: %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("chain state must be COMPLETED, got %q", summary.ChainState)
	}
	if got := zeroModelTaskCount(summary, "PASSED"); got != 3 {
		t.Fatalf("all three tasks must be PASSED, got %d: %+v", got, summary.TaskStateCounts)
	}
	if summary.RunID != "run-zero-model-noop" {
		t.Fatalf("summary must report the run it executed, got %q", summary.RunID)
	}
	// The persisted log must agree with the returned summary, otherwise the console
	// would report a completion the durable record does not carry.
	loaded, err := LoadRunSummary(context.Background(), runDir)
	if err != nil {
		t.Fatalf("reload the run summary: %v", err)
	}
	if loaded.ChainState != "COMPLETED" || loaded.Sequence != summary.Sequence || loaded.LastEventHash != summary.LastEventHash {
		t.Fatalf("the persisted summary must match the executed one: %+v vs %+v", loaded, summary)
	}
}

// TestZeroModelConsoleNaturalFailureStopsLaterTasks proves the console runtime
// fails naturally instead of skipping: a non-noop step is refused by localRuntime,
// the chain ends FAILED and task-three never appears.
func TestZeroModelConsoleNaturalFailureStopsLaterTasks(t *testing.T) {
	runDir := t.TempDir()
	workspace := t.TempDir()
	definition := chain.Definition{ID: "chain-zero-model-fail", Tasks: []chain.Task{
		zeroModelNoopTask("task-one"),
		{ID: "task-two", Steps: []chain.Step{
			{ID: "task-two-noop", Kind: "noop", Reason: "zero-model-e2e"},
			{ID: "task-two-code", Kind: "code", Mode: chain.ManagedChangeSet},
		}},
		zeroModelNoopTask("task-three"),
	}}

	_, err := ExecuteNoopRun(context.Background(), "run-zero-model-fail", runDir, definition, workspace, zeroModelClock)
	if err == nil {
		t.Fatal("a non-noop step must fail the console run instead of being skipped")
	}
	if !strings.Contains(err.Error(), "cannot execute step kind") {
		t.Fatalf("the failure must be the natural localRuntime refusal, got %v", err)
	}

	projection, taskIDs := zeroModelRunProjection(t, runDir)
	if projection.ChainState != "FAILED" {
		t.Fatalf("chain state must be FAILED, got %q", projection.ChainState)
	}
	// Only the two started tasks may have projected states; task-three must be absent
	// from both the projection and the event stream.
	if len(projection.TaskStates) != 2 {
		t.Fatalf("only T1 and T2 may hold a task state, got %+v", projection.TaskStates)
	}
	states := map[string]int{}
	for _, state := range projection.TaskStates {
		states[state]++
	}
	if states["PASSED"] != 1 || states["FAILED"] != 1 {
		t.Fatalf("T1 must pass and T2 must fail, got %+v", projection.TaskStates)
	}
	for _, taskID := range taskIDs {
		if taskID == "task-three" {
			t.Fatalf("task-three wrote events (%v); it must never start", taskIDs)
		}
	}
}

// TestZeroModelPreviewStaysRunnableWithoutAI pins the falsifiable half of the
// no-AI preview promise on the documented example: the zero-model chain must
// preview as runnable, ready and free of blocking evidence, with every previewed
// step resolving to noop.
//
// Why this test is shaped the way it is (P45 review finding): the call counters
// asserted below are static constants that BuildPreviewReport writes as literal
// zeros, so an assertion that only pins them is a tautology and can never turn
// red on an implementation change. The counter assertions are therefore kept as a
// documentation lock on the README operator promise ("preview never reaches a
// model"), NOT as a mechanical model-channel sentinel. The mechanical, falsifiable
// proofs are RunnableInCLI, Outcome, BlockingEvidence and the previewed step
// kinds: turning any fixture step from noop into an executable kind drives all of
// them red.
func TestZeroModelPreviewStaysRunnableWithoutAI(t *testing.T) {
	configPath := filepath.Join("..", "..", "examples", "no-ai-preview", "proofrail.chain.json")
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("the documented no-AI preview example config is missing (%s): %v (README.md:10 tells the operator to run `--chain .\\examples\\no-ai-preview\\proofrail.chain.json`)", configPath, err)
	}
	cfg, rawConfig, err := LoadChainConfig(configPath)
	if err != nil {
		t.Fatalf("load the example chain config: %v", err)
	}
	report, err := BuildPreviewReport(configPath, cfg, rawConfig, zeroModelClock())
	if err != nil {
		t.Fatalf("build the preview report: %v", err)
	}
	counters := report.CallCounters
	// Static fields, not a sentinel: BuildPreviewReport hardcodes these five
	// counters to 0, so this block pins the README promise to the operator
	// ("go run ./cmd/prfrail preview --chain .\examples\no-ai-preview\proofrail.chain.json"
	// reports no AI capability use) and would stay green even if the preview layer silently
	// gained a model channel. The mechanical zero-model guarantees are asserted
	// further down (RunnableInCLI / Outcome / BlockingEvidence / step kinds).
	if counters.CommandCalls != 0 {
		t.Fatalf("commandCalls must be 0, got %d", counters.CommandCalls)
	}
	if counters.NetworkCalls != 0 {
		t.Fatalf("networkCalls must be 0, got %d", counters.NetworkCalls)
	}
	if counters.ModelCalls != 0 {
		t.Fatalf("modelCalls must be 0, got %d", counters.ModelCalls)
	}
	if counters.CredentialReads != 0 {
		t.Fatalf("credentialReads must be 0, got %d", counters.CredentialReads)
	}
	if counters.VersionProbes != 0 {
		t.Fatalf("versionProbes must be 0, got %d", counters.VersionProbes)
	}
	preview := report.PreviewRecord.Preview
	if got := preview.Mode; got != "offline-read-only" {
		t.Fatalf("the preview must declare mode offline-read-only, got %q", got)
	}
	// RunnableInCLI is the derived quantity len(executableSteps) == 0: it is true
	// only when the chain holds nothing that would need a model or a runtime port.
	if !report.RunnableInCLI {
		t.Fatalf("the zero-model example must preview as runnable in the CLI (no executable step is allowed), outcome=%q warnings=%v", preview.Outcome, report.Warnings)
	}
	if got := preview.Outcome; got != "ready" {
		t.Fatalf("the zero-model preview must stay ready (no AI needed), got outcome %q with blockingEvidence %v", got, preview.BlockingEvidence)
	}
	if len(preview.BlockingEvidence) != 0 {
		t.Fatalf("a ready zero-model preview must carry no blocking evidence, got %v (each entry names a step that would require authorization)", preview.BlockingEvidence)
	}
	if len(report.Warnings) != 0 {
		t.Fatalf("a runnable zero-model preview must emit no executable-step warning, got %v", report.Warnings)
	}
	// buildPreviewUnknowns always reports the static capability gaps
	// (version-probe/network-access/credential-input), so only the
	// command-execution unknown is falsifiable here: it appears exactly when the
	// chain gained an executable step.
	for _, unknown := range report.Unknowns {
		if unknown.UnknownID == "command-execution" {
			t.Fatalf("the zero-model preview must not report a command-execution unknown (it would mean the chain needs a runner), got %+v", unknown)
		}
	}
	// Steps are taken from the preview report itself (already mapped by the
	// preview layer), never recomputed in the test: a single non-noop kind here
	// means the chain would need a runtime capability instead of staying no-AI.
	if len(preview.Steps) == 0 {
		t.Fatalf("the example chain must preview at least one step, got none")
	}
	for _, step := range preview.Steps {
		if step.Kind != "noop" {
			t.Fatalf("every previewed step must stay noop so the chain remains zero-model, but %s/%s previewed as %q", step.TaskID, step.StepID, step.Kind)
		}
	}
}
