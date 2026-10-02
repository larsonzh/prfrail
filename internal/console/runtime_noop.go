package console

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// localRuntime is the noop-only runtime used by ExecuteNoopRun. It is the
// composition root of the noop path, not a general runtime: the AgentRunner
// assembly lives in runtime_agent.go and never uses this type.
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

// RunPostflight is the noop-path-only fail-closed console stub. This stub serves
// the noop runtime exclusively (noop 路径专用); the noop-only CLI path never
// reaches postflight (no AgentRunner step can run there), and any definition
// that does reach it must fail closed (契约必须 fail-closed) instead of silently
// skipping the acceptance gate. The production AgentRunner postflight port lives
// in internal/adapters/agent_runner_postflight.go and is assembled in
// runtime_agent.go.
func (runtime *localRuntime) RunPostflight(context.Context, chain.PostflightRequest) (chain.PostflightDecision, error) {
	return chain.PostflightDecision{}, fmt.Errorf("%w: the console runtime cannot run postflight for run %s", chain.ErrPostflightUnavailable, runtime.runDir)
}

func (runtime *localRuntime) nextEventID() string {
	runtime.events++
	return fmt.Sprintf("event-%d", runtime.events)
}
