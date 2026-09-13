package chain

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	baselineHash  = "sha256:0000000000000000000000000000000000000000000000000000000000000001"
	acceptedOne   = "sha256:0000000000000000000000000000000000000000000000000000000000000002"
	acceptedTwo   = "sha256:0000000000000000000000000000000000000000000000000000000000000003"
	acceptedThree = "sha256:0000000000000000000000000000000000000000000000000000000000000004"
	stopHash      = "sha256:0000000000000000000000000000000000000000000000000000000000000005"
)

type fakeBaseline struct{ calls int }

func (fake *fakeBaseline) CaptureBaseline(context.Context, string) (SnapshotRef, error) {
	fake.calls++
	return SnapshotRef{Hash: baselineHash}, nil
}

type materialization struct{ task, parent string }
type fakeWorkspaces struct{ calls []materialization }

func (fake *fakeWorkspaces) Materialize(_ context.Context, _ string, task string, _ int, parent SnapshotRef) (Workspace, error) {
	fake.calls = append(fake.calls, materialization{task: task, parent: parent.Hash})
	return Workspace{Root: task}, nil
}

type fakeSteps struct {
	requests []StepRequest
	failTask string
	pause    func()
}

type fakeStepIntents struct {
	requests []StepRequest
	intent   StepExecutionIntent
	err      error
	prepare  func(StepRequest) (StepExecutionIntent, error)
}

func (fake *fakeStepIntents) PrepareStepIntent(_ context.Context, request StepRequest) (StepExecutionIntent, error) {
	fake.requests = append(fake.requests, request)
	if fake.prepare != nil {
		return fake.prepare(request)
	}
	return fake.intent, fake.err
}

func (fake *fakeSteps) Execute(_ context.Context, request StepRequest) (StepResult, error) {
	fake.requests = append(fake.requests, request)
	if fake.pause != nil {
		fake.pause()
		fake.pause = nil
	}
	if request.TaskID == fake.failTask {
		return StepResult{}, errors.New("runner failed")
	}
	return StepResult{Evidence: []string{stopHash}}, nil
}

type fakeAcceptance struct {
	hashes   []string
	producer evidence.Actor
	calls    int
}

func (fake *fakeAcceptance) Accept(context.Context, string, string, int, SnapshotRef, Workspace) (CandidateResult, error) {
	hash := fake.hashes[fake.calls]
	fake.calls++
	return CandidateResult{
		Snapshot:          SnapshotRef{Hash: hash},
		EvidenceRootHash:  stopHash,
		Evidence:          []string{stopHash},
		CandidateProducer: fake.producer,
	}, nil
}

type fakeReviewer struct {
	decisions []ReviewDecision
	err       error
	calls     int
}

func (fake *fakeReviewer) Review(context.Context, ReviewRequest) (ReviewDecision, error) {
	if fake.err != nil {
		return ReviewDecision{}, fake.err
	}
	decision := fake.decisions[fake.calls]
	fake.calls++
	return decision, nil
}

type fakePublisher struct {
	decisions []PromotionDecision
	err       error
	calls     int
}

func (fake *fakePublisher) Publish(context.Context, PromotionRequest) (PromotionDecision, error) {
	if fake.err != nil {
		return PromotionDecision{}, fake.err
	}
	decision := fake.decisions[fake.calls]
	fake.calls++
	return decision, nil
}

type fakeStopper struct {
	calls int
	err   error
}

func (fake *fakeStopper) Stop(context.Context, string) ([]string, error) {
	fake.calls++
	return []string{stopHash}, fake.err
}

type fakeReconciler struct {
	calls int
	err   error
}

func (fake *fakeReconciler) Reconcile(context.Context, string) error {
	fake.calls++
	return fake.err
}

func threeTaskDefinition() Definition {
	return Definition{ID: "chain-one", Tasks: []Task{
		{ID: "task-one", Steps: []Step{{ID: "code-managed", Kind: "code", Mode: ManagedChangeSet}, {ID: "build-one", Kind: "build"}}},
		{ID: "task-two", Steps: []Step{{ID: "noop-two", Kind: "noop", Reason: "not-needed"}, {ID: "code-isolated", Kind: "code", Mode: IsolatedWorkspace}}},
		{ID: "task-three", Steps: []Step{{ID: "verify-three", Kind: "verify"}}},
	}}
}

func testOptions(store EventStore) (Options, *fakeBaseline, *fakeWorkspaces, *fakeSteps, *fakeAcceptance, *fakeReviewer, *fakePublisher, *fakeStopper, *fakeReconciler) {
	baseline := &fakeBaseline{}
	workspaces := &fakeWorkspaces{}
	steps := &fakeSteps{}
	acceptance := &fakeAcceptance{hashes: []string{acceptedOne, acceptedTwo, acceptedThree}, producer: evidence.Actor{Type: "agent", ID: "agent-one"}}
	reviewer := &fakeReviewer{decisions: []ReviewDecision{
		{ReceiptID: "review-one", OccurredAt: "2026-09-07T01:02:03.000Z", RecordedBy: evidence.Actor{Type: "operator", ID: "reviewer-one"}, ReviewMode: "manual", Outcome: "approve", Evidence: []string{stopHash}},
		{ReceiptID: "review-two", OccurredAt: "2026-09-07T01:02:03.000Z", RecordedBy: evidence.Actor{Type: "operator", ID: "reviewer-one"}, ReviewMode: "manual", Outcome: "approve", Evidence: []string{stopHash}},
		{ReceiptID: "review-three", OccurredAt: "2026-09-07T01:02:03.000Z", RecordedBy: evidence.Actor{Type: "operator", ID: "reviewer-one"}, ReviewMode: "manual", Outcome: "approve", Evidence: []string{stopHash}},
	}}
	publisher := &fakePublisher{decisions: []PromotionDecision{
		{ReceiptID: "promotion-one", OccurredAt: "2026-09-07T01:02:03.000Z", Outcome: "completed", AcceptedSnapshotHash: acceptedOne, WriterStopEvidence: []string{stopHash}, LeaseEvidence: []string{stopHash}, Evidence: []string{stopHash}},
		{ReceiptID: "promotion-two", OccurredAt: "2026-09-07T01:02:03.000Z", Outcome: "completed", AcceptedSnapshotHash: acceptedTwo, WriterStopEvidence: []string{stopHash}, LeaseEvidence: []string{stopHash}, Evidence: []string{stopHash}},
		{ReceiptID: "promotion-three", OccurredAt: "2026-09-07T01:02:03.000Z", Outcome: "completed", AcceptedSnapshotHash: acceptedThree, WriterStopEvidence: []string{stopHash}, LeaseEvidence: []string{stopHash}, Evidence: []string{stopHash}},
	}}
	stopper := &fakeStopper{}
	reconciler := &fakeReconciler{}
	sequence := 0
	options := Options{RunID: "run-one", Definition: threeTaskDefinition(), Events: store, Baselines: baseline, Workspaces: workspaces, Steps: steps, Acceptance: acceptance, Reviewer: reviewer, Publisher: publisher, Stopper: stopper, Reconciler: reconciler, Clock: fixedClock, IDs: func() string { sequence++; return fmt.Sprintf("event-%d", sequence) }}
	return options, baseline, workspaces, steps, acceptance, reviewer, publisher, stopper, reconciler
}

func TestEngineRunsThreeTasksFourKindsAndBothModes(t *testing.T) {
	store := &memoryEvents{}
	options, baseline, workspaces, steps, acceptance, reviewer, publisher, _, _ := testOptions(store)
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if engine.Projection().ChainState != "COMPLETED" || baseline.calls != 1 || acceptance.calls != 3 || reviewer.calls != 3 || publisher.calls != 3 {
		t.Fatalf("unexpected completion: %+v", engine.Projection())
	}
	if len(steps.requests) != 4 {
		t.Fatalf("noop started a runner: %d requests", len(steps.requests))
	}
	if steps.requests[0].Step.Mode != ManagedChangeSet || steps.requests[2].Step.Mode != IsolatedWorkspace {
		t.Fatal("change modes were not preserved")
	}
	wantParents := []string{baselineHash, acceptedOne, acceptedTwo}
	for index, call := range workspaces.calls {
		if call.parent != wantParents[index] {
			t.Fatalf("task %s parent %s, want %s", call.task, call.parent, wantParents[index])
		}
	}
	if err := evidence.VerifyEventChain(store.events); err != nil {
		t.Fatal(err)
	}
}

func TestEnginePreparesAgentRunnerIntentBeforeSingleExecutionPort(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, steps, _, _, _, _, _ := testOptions(store)
	facts := &AgentRunnerImmutableFacts{
		RequestID:         "request-one",
		WorkspaceHash:     baselineHash,
		ContextHash:       acceptedOne,
		AuthorizationHash: acceptedTwo,
		BudgetHash:        acceptedThree,
	}
	preparer := &fakeStepIntents{prepare: func(request StepRequest) (StepExecutionIntent, error) {
		if request.Step.Kind == "code" && request.Step.Mode == IsolatedWorkspace {
			return StepExecutionIntent{ExecutionTarget: AgentRunnerExecution, AgentRunnerFacts: facts}, nil
		}
		return StepExecutionIntent{}, nil
	}}
	options.StepIntents = preparer
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(preparer.requests) != 4 || len(steps.requests) != 4 {
		t.Fatalf("unexpected prepare/execute calls: prepare=%d execute=%d", len(preparer.requests), len(steps.requests))
	}
	for index, request := range steps.requests {
		prepared := preparer.requests[index]
		if prepared.ExecutionTarget != DefaultExecution || prepared.AgentRunnerFacts != nil {
			t.Fatalf("preparer received mutable intent fields: %+v", prepared)
		}
		if request.Step.Kind == "code" && request.Step.Mode == IsolatedWorkspace {
			if request.ExecutionTarget != AgentRunnerExecution || request.AgentRunnerFacts == nil || *request.AgentRunnerFacts != *facts || request.AgentRunnerFacts == facts {
				t.Fatalf("request %d missing prepared intent: %+v", index, request)
			}
		} else if request.ExecutionTarget != DefaultExecution || request.AgentRunnerFacts != nil {
			t.Fatalf("request %d unexpectedly received AgentRunner intent: %+v", index, request)
		}
	}
}

func TestEngineRejectsIntentThatDoesNotMatchStep(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, steps, _, _, _, _, _ := testOptions(store)
	preparer := &fakeStepIntents{intent: StepExecutionIntent{
		ExecutionTarget:  AgentRunnerExecution,
		AgentRunnerFacts: &AgentRunnerImmutableFacts{RequestID: "request-one"},
	}}
	options.StepIntents = preparer
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); !errors.Is(err, ErrInvalidDefinition) {
		t.Fatalf("expected invalid intent failure, got %v", err)
	}
	if len(steps.requests) != 0 {
		t.Fatalf("invalid managed-step intent reached execution port: %d requests", len(steps.requests))
	}
	projection := engine.Projection()
	if projection.StepStates[stepKey("task-one", "code-managed", 1)] != "FAILED" || projection.ChainState != "FAILED" {
		t.Fatalf("invalid intent did not fail closed: %+v", projection)
	}
}

func TestEngineIntentPreparationFailureBlocksExecution(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, steps, acceptance, reviewer, publisher, _, _ := testOptions(store)
	prepareErr := errors.New("intent unavailable")
	preparer := &fakeStepIntents{err: prepareErr}
	options.StepIntents = preparer
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); !errors.Is(err, prepareErr) {
		t.Fatalf("expected preparation failure, got %v", err)
	}
	projection := engine.Projection()
	if projection.ChainState != "FAILED" || projection.StepStates[stepKey("task-one", "code-managed", 1)] != "FAILED" {
		t.Fatalf("preparation failure did not fail closed: %+v", projection)
	}
	if len(steps.requests) != 0 || acceptance.calls != 0 || reviewer.calls != 0 || publisher.calls != 0 {
		t.Fatalf("preparation failure reached later ports: steps=%d accept=%d review=%d publish=%d", len(steps.requests), acceptance.calls, reviewer.calls, publisher.calls)
	}
}

func TestTaskFailureStopsLaterTasks(t *testing.T) {
	store := &memoryEvents{}
	options, _, workspaces, steps, _, _, _, _, _ := testOptions(store)
	steps.failTask = "task-two"
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); err == nil {
		t.Fatal("expected task failure")
	}
	if engine.Projection().ChainState != "FAILED" || len(workspaces.calls) != 2 {
		t.Fatalf("later task started after failure: %+v", workspaces.calls)
	}
}

func TestPauseAndResume(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, steps, _, _, _, _, _ := testOptions(store)
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	steps.pause = engine.RequestPause
	if err := engine.Run(context.Background()); !errors.Is(err, ErrPaused) {
		t.Fatalf("pause result: %v", err)
	}
	if engine.Projection().ChainState != "PAUSED" {
		t.Fatal("chain did not pause")
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if engine.Projection().ChainState != "COMPLETED" {
		t.Fatal("chain did not resume to completion")
	}
}

func TestCancelStopsBeforeTerminalState(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, _, _, _, _, stopper, _ := testOptions(store)
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	engine.RequestCancel()
	if err := engine.Run(context.Background()); !errors.Is(err, ErrCancelled) {
		t.Fatalf("cancel result: %v", err)
	}
	if engine.Projection().ChainState != "CANCELLED" || stopper.calls != 1 {
		t.Fatal("cancel did not stop and terminate")
	}
}

func TestCrashReplayPausesRunningStepWithoutRedispatch(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, steps, _, _, _, _, reconciler := testOptions(store)
	options.Failpoint = func(point string) error {
		if point == "step-running" {
			return errors.New("simulated crash")
		}
		return nil
	}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); err == nil {
		t.Fatal("expected simulated crash")
	}
	if len(steps.requests) != 0 {
		t.Fatal("runner executed past crash point")
	}
	recoveryOptions, _, _, recoverySteps, _, _, _, _, recoveryReconciler := testOptions(store)
	sequence := len(store.events)
	recoveryOptions.IDs = func() string { sequence++; return fmt.Sprintf("event-%d", sequence) }
	recovered, err := New(context.Background(), recoveryOptions)
	if err != nil {
		t.Fatal(err)
	}
	if err := recovered.Recover(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("recovery result: %v", err)
	}
	if recovered.Projection().ChainState != "PAUSED" || !recovered.Projection().RecoveryUncertain || len(recoverySteps.requests) != 0 || recoveryReconciler.calls != 1 || reconciler.calls != 0 {
		t.Fatalf("unsafe recovery: %+v", recovered.Projection())
	}
	if err := recovered.Run(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("uncertain run resumed: %v", err)
	}
}

func TestReviewRejectMovesTaskToRepairPending(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, _, _, reviewer, publisher, _, _ := testOptions(store)
	reviewer.decisions[0] = ReviewDecision{
		ReceiptID:     "review-reject",
		OccurredAt:    "2026-09-07T01:02:03.000Z",
		RecordedBy:    evidence.Actor{Type: "operator", ID: "reviewer-two"},
		ReviewMode:    "manual",
		Outcome:       "reject",
		Reason:        &ReviewReason{Code: "docs-missing"},
		ErrorEvidence: []string{stopHash},
		Evidence:      []string{stopHash},
	}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); !errors.Is(err, ErrPaused) {
		t.Fatalf("run result: %v", err)
	}
	projection := engine.Projection()
	if projection.ChainState != "PAUSED" || projection.TaskStates[taskKey("task-one", 1)] != "REPAIR_PENDING" {
		t.Fatalf("unexpected projection: %+v", projection)
	}
	if publisher.calls != 0 {
		t.Fatal("publisher must not run when review rejected")
	}
}

func TestRejectSelfReview(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, _, acceptance, reviewer, publisher, _, _ := testOptions(store)
	acceptance.producer = evidence.Actor{Type: "agent", ID: "reviewer-one"}
	reviewer.decisions[0] = ReviewDecision{
		ReceiptID:  "review-self",
		OccurredAt: "2026-09-07T01:02:03.000Z",
		RecordedBy: evidence.Actor{Type: "operator", ID: "reviewer-one"},
		ReviewMode: "manual",
		Outcome:    "approve",
	}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); !errors.Is(err, ErrPaused) {
		t.Fatalf("run result: %v", err)
	}
	if engine.Projection().TaskStates[taskKey("task-one", 1)] != "REPAIR_PENDING" || publisher.calls != 0 {
		t.Fatalf("unexpected projection: %+v", engine.Projection())
	}
}

func TestRejectExpiredOrWrongCandidateWaiver(t *testing.T) {
	for _, scenario := range []struct {
		name     string
		decision ReviewDecision
	}{
		{
			name: "expired waiver",
			decision: ReviewDecision{
				ReceiptID:     "review-waive-expired",
				OccurredAt:    "2026-09-07T01:02:03.000Z",
				RecordedBy:    evidence.Actor{Type: "operator", ID: "reviewer-one"},
				ReviewMode:    "manual",
				Outcome:       "waive",
				Reason:        &ReviewReason{Code: "approved-exception"},
				ErrorEvidence: []string{stopHash},
				WaiverAuthorization: &WaiverAuthorization{
					AuthorizedBy:    evidence.Actor{Type: "operator", ID: "release-manager"},
					PolicyBasisHash: stopHash,
					Scope:           []string{"task-one"},
					ExpiresAt:       "2026-09-07T01:00:00.000Z",
				},
			},
		},
		{
			name: "wrong candidate",
			decision: ReviewDecision{
				ReceiptID:             "review-waive-candidate",
				OccurredAt:            "2026-09-07T01:02:03.000Z",
				RecordedBy:            evidence.Actor{Type: "operator", ID: "reviewer-one"},
				ReviewMode:            "manual",
				Outcome:               "waive",
				CandidateSnapshotHash: baselineHash,
				Reason:                &ReviewReason{Code: "approved-exception"},
				ErrorEvidence:         []string{stopHash},
				WaiverAuthorization: &WaiverAuthorization{
					AuthorizedBy:    evidence.Actor{Type: "operator", ID: "release-manager"},
					PolicyBasisHash: stopHash,
					Scope:           []string{"task-one"},
					ExpiresAt:       "2026-09-08T01:02:03.000Z",
				},
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			store := &memoryEvents{}
			options, _, _, _, _, reviewer, publisher, _, _ := testOptions(store)
			reviewer.decisions[0] = scenario.decision
			engine, err := New(context.Background(), options)
			if err != nil {
				t.Fatal(err)
			}
			if err := engine.Run(context.Background()); !errors.Is(err, ErrPaused) {
				t.Fatalf("run result: %v", err)
			}
			if engine.Projection().TaskStates[taskKey("task-one", 1)] != "REPAIR_PENDING" || publisher.calls != 0 {
				t.Fatalf("unexpected projection: %+v", engine.Projection())
			}
		})
	}
}

func TestValidIndependentApprovalPublishesOnce(t *testing.T) {
	store := &memoryEvents{}
	definition := Definition{ID: "chain-one", Tasks: []Task{{ID: "task-one", Steps: []Step{{ID: "code-one", Kind: "code", Mode: ManagedChangeSet}}}}}
	options, _, _, _, _, reviewer, publisher, _, _ := testOptions(store)
	options.Definition = definition
	reviewer.decisions = []ReviewDecision{{ReceiptID: "review-approve", OccurredAt: "2026-09-07T01:02:03.000Z", RecordedBy: evidence.Actor{Type: "operator", ID: "reviewer-two"}, ReviewMode: "manual", Outcome: "approve", Evidence: []string{stopHash}}}
	publisher.decisions = []PromotionDecision{{ReceiptID: "promotion-complete", OccurredAt: "2026-09-07T01:02:03.000Z", Outcome: "completed", AcceptedSnapshotHash: acceptedOne, WriterStopEvidence: []string{stopHash}, LeaseEvidence: []string{stopHash}, Evidence: []string{stopHash}}}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	projection := engine.Projection()
	if projection.ChainState != "COMPLETED" || projection.TaskStates[taskKey("task-one", 1)] != "PASSED" {
		t.Fatalf("unexpected projection: %+v", projection)
	}
	if publisher.calls != 1 {
		t.Fatalf("publisher calls: %d", publisher.calls)
	}
}

func TestRepairPendingCannotResumeWithoutNewAttempt(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, _, _, reviewer, publisher, _, _ := testOptions(store)
	reviewer.decisions[0] = ReviewDecision{
		ReceiptID:     "review-reject",
		OccurredAt:    "2026-09-07T01:02:03.000Z",
		RecordedBy:    evidence.Actor{Type: "operator", ID: "reviewer-two"},
		ReviewMode:    "manual",
		Outcome:       "reject",
		Reason:        &ReviewReason{Code: "docs-missing"},
		ErrorEvidence: []string{stopHash},
	}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); !errors.Is(err, ErrPaused) {
		t.Fatalf("run result: %v", err)
	}
	before := engine.Projection()
	if err := engine.Run(context.Background()); !errors.Is(err, ErrPaused) {
		t.Fatalf("resume result: %v", err)
	}
	after := engine.Projection()
	if after.ChainState != "PAUSED" || after.TaskStates[taskKey("task-one", 1)] != before.TaskStates[taskKey("task-one", 1)] || publisher.calls != 0 {
		t.Fatalf("resume mutated repair-pending task: %+v", after)
	}
	if err := evidence.VerifyEventChain(store.events); err != nil {
		t.Fatal(err)
	}
}
