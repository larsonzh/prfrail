package console

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

const timestampLayout = "2006-01-02T15:04:05.000Z"

type RunSummary struct {
	RunID           string       `json:"runId"`
	ChainState      string       `json:"chainState"`
	Sequence        int          `json:"sequence"`
	LastEventHash   string       `json:"lastEventHash"`
	TaskStateCounts []StateCount `json:"taskStateCounts"`
	StepStateCounts []StateCount `json:"stepStateCounts"`
	EventLogPath    string       `json:"eventLogPath"`
}

type StateCount struct {
	State string `json:"state"`
	Count int    `json:"count"`
}

func ExecuteNoopRun(ctx context.Context, runID, runDir string, definition chain.Definition, workspaceRoot string, now func() time.Time) (RunSummary, error) {
	store, err := chain.NewFileEventStore(runEventLogPath(runDir))
	if err != nil {
		return RunSummary{}, err
	}

	runtime := &localRuntime{runDir: runDir, workspaceRoot: workspaceRoot, now: now}
	options := chain.Options{
		RunID:      runID,
		Definition: definition,
		Events:     store,
		Baselines:  runtime,
		Workspaces: runtime,
		Steps:      runtime,
		Acceptance: runtime,
		Reviewer:   runtime,
		Publisher:  runtime,
		Stopper:    runtime,
		Reconciler: runtime,
		Clock:      now,
		IDs:        runtime.nextEventID,
	}
	engine, err := chain.New(ctx, options)
	if err != nil {
		return RunSummary{}, err
	}
	if err := engine.Run(ctx); err != nil {
		return RunSummary{}, err
	}
	return summarizeProjection(engine.Projection(), runEventLogPath(runDir)), nil
}

func LoadRunSummary(ctx context.Context, runDir string) (RunSummary, error) {
	path := runEventLogPath(runDir)
	if _, err := os.Stat(path); err != nil {
		return RunSummary{}, err
	}
	store, err := chain.NewFileEventStore(path)
	if err != nil {
		return RunSummary{}, err
	}
	events, err := store.Load(ctx)
	if err != nil {
		return RunSummary{}, err
	}
	if len(events) == 0 {
		return RunSummary{}, fmt.Errorf("event log %q is empty", path)
	}
	projection, err := chain.Rebuild(events)
	if err != nil {
		return RunSummary{}, err
	}
	return summarizeProjection(projection, path), nil
}

func runEventLogPath(runDir string) string {
	return filepath.Join(runDir, "events", "state-events.jsonl")
}

func summarizeProjection(projection chain.Projection, eventLogPath string) RunSummary {
	return RunSummary{
		RunID:           projection.RunID,
		ChainState:      projection.ChainState,
		Sequence:        projection.Sequence,
		LastEventHash:   projection.LastEventHash,
		TaskStateCounts: countStates(projection.TaskStates),
		StepStateCounts: countStates(projection.StepStates),
		EventLogPath:    eventLogPath,
	}
}

func countStates(states map[string]string) []StateCount {
	if len(states) == 0 {
		return []StateCount{}
	}
	counters := map[string]int{}
	for _, state := range states {
		counters[state]++
	}
	names := make([]string, 0, len(counters))
	for name := range counters {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make([]StateCount, 0, len(names))
	for _, name := range names {
		result = append(result, StateCount{State: name, Count: counters[name]})
	}
	return result
}

type localRuntime struct {
	runDir        string
	workspaceRoot string
	now           func() time.Time
	events        int
	reviewCount   int
	publishCount  int
	acceptCount   int
}

func (runtime *localRuntime) CaptureBaseline(context.Context, string) (chain.SnapshotRef, error) {
	hash := digest("proofrail:local-baseline:1\n", runtime.runDir, runtime.workspaceRoot)
	return chain.SnapshotRef{Hash: hash}, nil
}

func (runtime *localRuntime) Materialize(_ context.Context, _ string, taskID string, attempt int, parent chain.SnapshotRef) (chain.Workspace, error) {
	dir := filepath.Join(runtime.runDir, "workspaces", fmt.Sprintf("%s-attempt-%d", taskID, attempt))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return chain.Workspace{}, err
	}
	_ = parent
	return chain.Workspace{Root: dir}, nil
}

func (runtime *localRuntime) Execute(_ context.Context, request chain.StepRequest) (chain.StepResult, error) {
	return chain.StepResult{}, fmt.Errorf("local runtime cannot execute step kind %q for step %q", request.Step.Kind, request.Step.ID)
}

func (runtime *localRuntime) Accept(_ context.Context, runID string, taskID string, attempt int, parent chain.SnapshotRef, workspace chain.Workspace) (chain.CandidateResult, error) {
	runtime.acceptCount++
	snapshotHash := digest(
		"proofrail:local-candidate:1\n",
		runID,
		taskID,
		fmt.Sprintf("%d", attempt),
		parent.Hash,
		workspace.Root,
		fmt.Sprintf("%d", runtime.acceptCount),
	)
	evidenceRootHash := digest("proofrail:local-evidence-root:1\n", snapshotHash)
	evidenceHash := digest("proofrail:local-evidence:1\n", runID, taskID)
	return chain.CandidateResult{
		Snapshot:          chain.SnapshotRef{Hash: snapshotHash},
		EvidenceRootHash:  evidenceRootHash,
		Evidence:          []string{evidenceHash},
		CandidateProducer: evidence.Actor{Type: "agent", ID: "local-runner"},
	}, nil
}

func (runtime *localRuntime) Review(_ context.Context, request chain.ReviewRequest) (chain.ReviewDecision, error) {
	runtime.reviewCount++
	return chain.ReviewDecision{
		ReceiptID:  fmt.Sprintf("review-%d", runtime.reviewCount),
		OccurredAt: runtime.now().UTC().Format(timestampLayout),
		RecordedBy: evidence.Actor{Type: "operator", ID: "local-reviewer"},
		ReviewMode: "manual",
		Outcome:    "approve",
		Evidence:   []string{digest("proofrail:local-review:1\n", request.RunID, request.TaskID)},
	}, nil
}

func (runtime *localRuntime) Publish(_ context.Context, request chain.PromotionRequest) (chain.PromotionDecision, error) {
	runtime.publishCount++
	return chain.PromotionDecision{
		ReceiptID:            fmt.Sprintf("promotion-%d", runtime.publishCount),
		OccurredAt:           runtime.now().UTC().Format(timestampLayout),
		Outcome:              "completed",
		AcceptedSnapshotHash: request.CandidateSnapshotHash,
		WriterStopEvidence:   []string{digest("proofrail:local-stop:1\n", request.RunID, request.TaskID)},
		LeaseEvidence:        []string{digest("proofrail:local-lease:1\n", request.RunID)},
		Evidence:             []string{digest("proofrail:local-publish:1\n", request.ReviewReceiptHash)},
	}, nil
}

func (runtime *localRuntime) Stop(context.Context, string) ([]string, error) {
	return []string{digest("proofrail:local-stopper:1\n", runtime.runDir)}, nil
}

func (runtime *localRuntime) Reconcile(context.Context, string) error {
	return nil
}

func (runtime *localRuntime) nextEventID() string {
	runtime.events++
	return fmt.Sprintf("event-%d", runtime.events)
}

func digest(domain string, values ...string) string {
	return evidence.Digest(domain, []byte(strings.Join(values, "\n")))
}
