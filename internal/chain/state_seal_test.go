package chain

// P43-T027B-STATE-SEAL ③ test-first (RED): the evidence seal barrier (C4), the
// immutable context digest (C5) and the resolved attempt (C5/N3).
//
// Design: docs/validation/evidence/P43-T027B-STATE-SEAL/01-arch-design/review.md
//
// The expected RED state of this slice is a compile failure, not a test failure:
// SealSink / LoadTaskSeal / TaskAttemptResolver / ContextKey / taskEntityAt /
// stepEntityAt do not exist yet. Every invariant below must survive the
// implementation unchanged.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// sealFailingEvents wraps a store and fails every seal.
type sealFailingEvents struct {
	EventStore
	err error
}

func (fake *sealFailingEvents) SealTaskEvidence(context.Context, SealRecord) (string, error) {
	return "", fake.err
}

func (fake *sealFailingEvents) LoadTaskSeal(context.Context, string, string, int) (SealRecord, string, error) {
	return SealRecord{}, "", fmt.Errorf("%w: no seal", ErrSealMissing)
}

// fixedAttempts resolves one configured attempt for every task.
type fixedAttempts struct {
	attempt int
	err     error
	tasks   []string
}

func (fake *fixedAttempts) ResolveTaskAttempt(_ context.Context, _ string, taskID string) (int, error) {
	fake.tasks = append(fake.tasks, taskID)
	if fake.err != nil {
		return 0, fake.err
	}
	return fake.attempt, nil
}

func sealTestOptions(t *testing.T, store EventStore) Options {
	t.Helper()
	options, _, _, _, _, _, _, _, _ := testOptions(store)
	return options
}

func fileStore(t *testing.T) *FileEventStore {
	t.Helper()
	store, err := NewFileEventStore(filepath.Join(t.TempDir(), "events.jsonl"))
	if err != nil {
		t.Fatalf("NewFileEventStore: %v", err)
	}
	return store
}

// I1: a task may not reach PASSED when the seal cannot be persisted. The seal
// barrier is the only writer of the accepted terminal, so a failed seal leaves
// the task in REPAIR_PENDING and never completes the chain.
func TestTaskSealBarrierPreventsAccepted(t *testing.T) {
	store := &sealFailingEvents{EventStore: fileStore(t), err: errors.New("seal disk full")}
	options := sealTestOptions(t, store)
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	err = engine.Run(context.Background())
	if err == nil {
		t.Fatal("run completed although the evidence seal failed")
	}
	for task, state := range engine.Projection().TaskStates {
		if state == "PASSED" {
			t.Fatalf("task %s passed without a sealed closure", task)
		}
	}
	if engine.Projection().ChainState == "COMPLETED" {
		t.Fatal("chain completed without a sealed closure")
	}
}

// I2: an invalid hash inside the accepted closure is a missing piece of
// evidence, not a formatting detail.
func TestClosureIncompleteCannotSeal(t *testing.T) {
	store := fileStore(t)
	options := sealTestOptions(t, store)
	acceptance := &fakeAcceptance{
		hashes:   []string{acceptedOne, acceptedTwo, acceptedThree},
		producer: evidence.Actor{Type: "agent", ID: "agent-one"},
	}
	acceptance.evidence = [][]string{{stopHash, "not-a-hash"}, {stopHash}, {stopHash}}
	options.Acceptance = acceptance
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := engine.Run(context.Background()); err == nil {
		t.Fatal("run accepted a closure that carries an invalid hash")
	}
	if engine.Projection().TaskStates[taskKey("task-one", 1)] == "PASSED" {
		t.Fatal("task passed with an incomplete evidence closure")
	}
}

// I3: recovery must verify every accepted terminal against its seal.
func TestPassedTaskSealMustVerifyOnRecovery(t *testing.T) {
	store := fileStore(t)
	options := sealTestOptions(t, store)
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	seals, err := filepath.Glob(filepath.Join(filepath.Dir(store.path), "seals", "run-one", "*.seal.json"))
	if err != nil {
		t.Fatalf("glob seals: %v", err)
	}
	if len(seals) == 0 {
		t.Fatal("no seal was written for the accepted tasks")
	}
	if err := os.Remove(seals[0]); err != nil {
		t.Fatalf("remove seal: %v", err)
	}
	reopened, err := New(context.Background(), options)
	if err == nil {
		err = reopened.Recover(context.Background())
	}
	if !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("verification error = %v, want ErrRecoveryUncertain", err)
	}
}

// I4: the resolved attempt reaches every port and every entity key.
func TestReopenUsesResolvedAttempt(t *testing.T) {
	store := fileStore(t)
	options := sealTestOptions(t, store)
	attempts := &fixedAttempts{attempt: 2}
	options.Attempts = attempts
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(attempts.tasks) == 0 {
		t.Fatal("the resolver was never consulted")
	}
	workspaces := options.Workspaces.(*fakeWorkspaces)
	for _, call := range workspaces.attempts {
		if call != 2 {
			t.Fatalf("Materialize attempt = %d, want 2", call)
		}
	}
	steps := options.Steps.(*fakeSteps)
	for _, request := range steps.requests {
		if request.Attempt != 2 {
			t.Fatalf("step request attempt = %d, want 2", request.Attempt)
		}
	}
	reviewer := options.Reviewer.(*fakeReviewer)
	if len(reviewer.attempts) == 0 || reviewer.attempts[0] != 2 {
		t.Fatalf("review attempts = %v, want 2", reviewer.attempts)
	}
	publisher := options.Publisher.(*fakePublisher)
	if len(publisher.attempts) == 0 || publisher.attempts[0] != 2 {
		t.Fatalf("promotion attempts = %v, want 2", publisher.attempts)
	}
	entity := taskEntityAt("run-one", "task-one", 2)
	if engine.Projection().TaskStates[taskKey("task-one", 2)] != "PASSED" {
		t.Fatalf("attempt 2 task state = %q, want PASSED", engine.Projection().TaskStates[taskKey("task-one", 2)])
	}
	if engine.state.current(entity) != "PASSED" {
		t.Fatalf("entity-carrying attempt reads %q", engine.state.current(entity))
	}
}

// I5: K is derived from immutable inputs only, so a re-entry recomputes the same
// value and a different attempt can never reuse another attempt's context.
func TestContextKeyStableAcrossReentry(t *testing.T) {
	first := ContextKey("run-one", "task-one", "code-managed", 1, baselineHash, "/w", "isolated-workspace", "request-one", acceptedOne)
	second := ContextKey("run-one", "task-one", "code-managed", 1, baselineHash, "/w", "isolated-workspace", "request-one", acceptedOne)
	if first != second {
		t.Fatalf("context key drifted across re-entry: %s vs %s", first, second)
	}
	if !evidence.ValidHash(first) {
		t.Fatalf("context key %q is not a valid hash", first)
	}
	other := ContextKey("run-one", "task-one", "code-managed", 2, baselineHash, "/w", "isolated-workspace", "request-one", acceptedOne)
	if other == first {
		t.Fatal("context key ignored the attempt")
	}
	moved := ContextKey("run-one", "task-one", "code-managed", 1, acceptedOne, "/w", "isolated-workspace", "request-one", acceptedOne)
	if moved == first {
		t.Fatal("context key ignored the parent snapshot")
	}
}

// I6: the seal record binds the candidate and both receipt hashes that the
// accepted event carries, so a reader can prove the closure instead of trusting it.
func TestSealRecordBindsAcceptedClosure(t *testing.T) {
	store := fileStore(t)
	options := sealTestOptions(t, store)
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	record, hash, err := store.LoadTaskSeal(context.Background(), "run-one", "task-one", 1)
	if err != nil {
		t.Fatalf("LoadTaskSeal: %v", err)
	}
	if !evidence.ValidHash(hash) {
		t.Fatalf("seal hash %q is not a valid hash", hash)
	}
	if record.CandidateHash != acceptedOne || record.EvidenceRootHash != stopHash {
		t.Fatalf("seal record = %+v, want candidate %s and closure %s", record, acceptedOne, stopHash)
	}
	if !evidence.ValidHash(record.ReviewReceiptHash) || !evidence.ValidHash(record.PromotionReceiptHash) {
		t.Fatalf("seal record does not bind both receipts: %+v", record)
	}
	if record.Attempt != 1 || record.TaskID != "task-one" || record.RunID != "run-one" {
		t.Fatalf("seal identity = %+v", record)
	}
}

// I7: the barrier refuses to accept when the resolver cannot name the attempt.
func TestAttemptResolverFailureBlocksAcceptance(t *testing.T) {
	store := fileStore(t)
	options := sealTestOptions(t, store)
	options.Attempts = &fixedAttempts{err: fmt.Errorf("attempt ledger unavailable")}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := engine.Run(context.Background()); err == nil {
		t.Fatal("run completed although the attempt resolver failed")
	}
	if engine.Projection().TaskStates[taskKey("task-one", 1)] == "PASSED" {
		t.Fatal("task passed without a resolved attempt")
	}
	_ = store
}

// I8: a terminal bound to another attempt must be refused, so the routed identity
// can never address a sibling attempt.
func TestTerminalRouteRefusesForeignAttempt(t *testing.T) {
	store := fileStore(t)
	options := sealTestOptions(t, store)
	options.Attempts = &fixedAttempts{attempt: 2}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	terminal := AgentRunnerTerminal{
		RunID:          "run-one",
		TaskID:         "task-one",
		StepID:         "code-isolated",
		Attempt:        1,
		Status:         AgentRunnerTerminalFailed,
		RequestID:      "request-one",
		RequestHash:    baselineHash,
		CompletionHash: acceptedOne,
		SessionID:      "session-one",
		Evidence:       []string{baselineHash, acceptedOne},
	}
	_, err = engine.SubmitAgentRunnerTerminal(context.Background(), terminal)
	if err == nil {
		t.Fatal("terminal for attempt 1 was routed while the run resolves attempt 2")
	}
	// The refusal must be the attempt binding, not a shape error: a terminal that
	// fails Validate never reaches the comparison this test exists to pin.
	if !strings.Contains(err.Error(), "does not belong to this run and attempt") {
		t.Fatalf("refusal = %v, want the attempt binding rejection", err)
	}
}

// I9: K is sealed and re-derived, so a seal written under another context key is
// refused instead of replayed.
func TestSealBindsTaskContextKey(t *testing.T) {
	store := fileStore(t)
	options := sealTestOptions(t, store)
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	record, _, err := store.LoadTaskSeal(context.Background(), "run-one", "task-one", 1)
	if err != nil {
		t.Fatalf("LoadTaskSeal: %v", err)
	}
	want := TaskContextKey("run-one", "task-one", 1, baselineHash, "task-one")
	if record.ContextKey != want {
		t.Fatalf("sealed context key = %q, want %q", record.ContextKey, want)
	}
	if other := TaskContextKey("run-one", "task-one", 2, baselineHash, "task-one"); other == want {
		t.Fatal("context key ignored the attempt")
	}
}
