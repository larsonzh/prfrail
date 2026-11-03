package chain

// P52-S4 K-binding counterexamples (ACCEPTANCE-SECURITY-COUNTEREXAMPLES).
//
// Tests only: this file adds no production behaviour and touches no existing
// test file. The assembly below is deliberately self-contained (it does not
// reuse the helpers of state_seal_test.go or engine_test.go) so the
// counterexamples stay readable on their own.
//
// Three cases are covered:
//  1. a seal record whose ContextKey is not a hash is refused as incomplete;
//  2. a seal record with an empty ContextKey is still admitted (current
//     behaviour, GAP-78 candidate: the seal contract allows an empty K and each
//     consumer fails closed on its own);
//  3. an acceptance replay carrying a foreign K is refused, while the same
//     record with the engine-derived K replays successfully.

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	cxDomain = "proofrail:counterexample:1\n"
	cxRunID  = "run-one"
	cxTaskID = "task-one"
)

var (
	cxParentHash           = evidence.Digest(cxDomain, []byte("parent"))
	cxCandidateHash        = evidence.Digest(cxDomain, []byte("candidate"))
	cxEvidenceRootHash     = evidence.Digest(cxDomain, []byte("evidence-root"))
	cxReviewReceiptHash    = evidence.Digest(cxDomain, []byte("review-receipt"))
	cxPromotionReceiptHash = evidence.Digest(cxDomain, []byte("promotion-receipt"))
	cxStopHash             = evidence.Digest(cxDomain, []byte("stop"))
	cxForeignContextKey    = evidence.Digest(cxDomain, []byte("foreign-context-key"))
)

type cxBaseline struct{}

func (cxBaseline) CaptureBaseline(context.Context, string) (SnapshotRef, error) {
	return SnapshotRef{Hash: cxParentHash}, nil
}

type cxWorkspaces struct{}

func (cxWorkspaces) Materialize(_ context.Context, _ string, task string, _ int, _ SnapshotRef) (Workspace, error) {
	return Workspace{Root: task}, nil
}

type cxSteps struct{}

func (cxSteps) Execute(context.Context, StepRequest) (StepResult, error) {
	return StepResult{Evidence: []string{cxStopHash}}, nil
}

type cxAcceptance struct{}

func (cxAcceptance) Accept(context.Context, string, string, int, SnapshotRef, Workspace) (CandidateResult, error) {
	return CandidateResult{
		Snapshot:          SnapshotRef{Hash: cxCandidateHash},
		EvidenceRootHash:  cxEvidenceRootHash,
		Evidence:          []string{cxStopHash},
		CandidateProducer: evidence.Actor{Type: "agent", ID: "agent-one"},
	}, nil
}

type cxReviewer struct{}

func (cxReviewer) Review(context.Context, ReviewRequest) (ReviewDecision, error) {
	return ReviewDecision{
		ReceiptID:  "review-one",
		OccurredAt: "2026-11-03T00:00:00.000Z",
		RecordedBy: evidence.Actor{Type: "operator", ID: "reviewer-one"},
		ReviewMode: "manual",
		Outcome:    "approve",
		Evidence:   []string{cxStopHash},
	}, nil
}

type cxPublisher struct{}

func (cxPublisher) Publish(context.Context, PromotionRequest) (PromotionDecision, error) {
	return PromotionDecision{
		ReceiptID:            "promotion-one",
		OccurredAt:           "2026-11-03T00:00:00.000Z",
		Outcome:              "completed",
		AcceptedSnapshotHash: cxCandidateHash,
		WriterStopEvidence:   []string{cxStopHash},
		LeaseEvidence:        []string{cxStopHash},
		Evidence:             []string{cxStopHash},
	}, nil
}

type cxPostflight struct{}

func (cxPostflight) RunPostflight(context.Context, PostflightRequest) (PostflightDecision, error) {
	return PostflightDecision{Outcome: PostflightPassed}, nil
}

type cxStopper struct{}

func (cxStopper) Stop(context.Context, string) ([]string, error) {
	return []string{cxStopHash}, nil
}

type cxReconciler struct{}

func (cxReconciler) Reconcile(context.Context, string) error { return nil }

// cxOptions builds the smallest valid assembly for a single managed code task.
func cxOptions(store EventStore) Options {
	return Options{
		RunID: cxRunID,
		Definition: Definition{ID: "chain-one", Tasks: []Task{{
			ID:    cxTaskID,
			Steps: []Step{{ID: "code-managed", Kind: "code", Mode: ManagedChangeSet}},
		}}},
		Events:     store,
		Baselines:  cxBaseline{},
		Workspaces: cxWorkspaces{},
		Steps:      cxSteps{},
		Acceptance: cxAcceptance{},
		Reviewer:   cxReviewer{},
		Publisher:  cxPublisher{},
		Postflight: cxPostflight{},
		Stopper:    cxStopper{},
		Reconciler: cxReconciler{},
	}
}

// cxSealRecord builds a well-formed seal record; the caller only varies the
// field under test.
func cxSealRecord() SealRecord {
	return SealRecord{
		SchemaVersion:        sealSchemaVersion,
		RunID:                cxRunID,
		TaskID:               cxTaskID,
		Attempt:              1,
		ParentHash:           cxParentHash,
		CandidateHash:        cxCandidateHash,
		EvidenceRootHash:     cxEvidenceRootHash,
		ReviewReceiptHash:    cxReviewReceiptHash,
		PromotionReceiptHash: cxPromotionReceiptHash,
		ReviewDecision:       "approve",
		PromotionDecision:    "completed",
		Evidence:             []string{cxStopHash, cxCandidateHash},
	}
}

func cxNewFileStore(t *testing.T) *FileEventStore {
	t.Helper()
	store, err := NewFileEventStore(filepath.Join(t.TempDir(), "events.jsonl"))
	if err != nil {
		t.Fatalf("NewFileEventStore: %v", err)
	}
	return store
}

func cxNewEngine(t *testing.T, store EventStore) *Engine {
	t.Helper()
	engine, err := New(context.Background(), cxOptions(store))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return engine
}

// TestSealValidateRefusesNonHashContextKey pins the incomplete-closure branch:
// a context key that is not an evidence hash means the seal has nothing
// trustworthy to bind.
func TestSealValidateRefusesNonHashContextKey(t *testing.T) {
	record := cxSealRecord()
	record.ContextKey = "not-a-hash"
	err := record.validate()
	if !errors.Is(err, ErrSealIncomplete) {
		t.Fatalf("validate() = %v, want ErrSealIncomplete", err)
	}
}

// TestSealValidateAdmitsEmptyContextKey pins the CURRENT behaviour and is
// recorded as a GAP-78 candidate: validate() accepts an empty context key, so
// the seal contract itself does not require K. Each consumer must therefore
// fail closed on the missing binding on its own (acceptFromExistingSeal
// compares the sealed K against the engine-derived K; that comparison refuses
// an empty K for an attempt that derives one).
func TestSealValidateAdmitsEmptyContextKey(t *testing.T) {
	record := cxSealRecord()
	record.ContextKey = ""
	if err := record.validate(); err != nil {
		t.Fatalf("validate() with empty context key = %v, want nil (GAP-78 current behaviour)", err)
	}
}

// TestAcceptanceReplayRefusesForeignContextKey drives the real engine replay
// path (Engine.Run -> runTask -> acceptFromExistingSeal). acceptFromExistingSeal
// is unexported, but the test lives in package chain and reaches it through the
// production entry point: the foreign seal is persisted by the real
// FileEventStore before Run starts, so the engine's own replay path is the one
// under test (internal/chain/engine.go:230 calls it, engine.go:352-354 performs
// the K comparison).
func TestAcceptanceReplayRefusesForeignContextKey(t *testing.T) {
	t.Run("foreign K is refused", func(t *testing.T) {
		store := cxNewFileStore(t)
		record := cxSealRecord()
		record.ContextKey = cxForeignContextKey
		if record.ContextKey == TaskContextKey(cxRunID, cxTaskID, 1, cxParentHash, cxTaskID) {
			t.Fatal("the foreign context key must differ from the engine-derived key")
		}
		if _, err := store.SealTaskEvidence(context.Background(), record); err != nil {
			t.Fatalf("SealTaskEvidence: %v", err)
		}

		engine := cxNewEngine(t, store)
		err := engine.Run(context.Background())
		if err == nil {
			t.Fatal("a replay carrying a foreign context key must not complete the run")
		}
		if !errors.Is(err, ErrPaused) {
			t.Fatalf("Run() = %v, want ErrPaused", err)
		}
		if state := engine.Projection().TaskStates[taskKey(cxTaskID, 1)]; state != "REPAIR_PENDING" {
			t.Fatalf("task state = %q, want REPAIR_PENDING", state)
		}
		events, err := store.Load(context.Background())
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		reason := ""
		for _, record := range events {
			if record.Event.Entity.Kind == "task" && record.Event.Entity.TaskID == cxTaskID && record.Event.ToState == "REPAIR_PENDING" {
				reason = record.Event.Reason.Code
			}
		}
		if reason != "evidence-seal-mismatch" {
			t.Fatalf("repair reason = %q, want evidence-seal-mismatch", reason)
		}
	})

	// Control: the same record shape with the engine-derived K replays
	// successfully, so the refusal above is attributable to the K difference and
	// not to the assembly.
	t.Run("bound K replays", func(t *testing.T) {
		store := cxNewFileStore(t)
		record := cxSealRecord()
		record.ContextKey = TaskContextKey(cxRunID, cxTaskID, 1, cxParentHash, cxTaskID)
		if _, err := store.SealTaskEvidence(context.Background(), record); err != nil {
			t.Fatalf("SealTaskEvidence: %v", err)
		}

		engine := cxNewEngine(t, store)
		if err := engine.Run(context.Background()); err != nil {
			t.Fatalf("Run: %v", err)
		}
		if state := engine.Projection().TaskStates[taskKey(cxTaskID, 1)]; state != "PASSED" {
			t.Fatalf("task state = %q, want PASSED", state)
		}
		// The replay accepts the sealed closure instead of re-freezing a new one.
		if engine.Projection().ChainState != "COMPLETED" {
			t.Fatalf("chain state = %q, want COMPLETED", engine.Projection().ChainState)
		}
	})
}
