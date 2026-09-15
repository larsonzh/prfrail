package chain

import (
	"context"
	"errors"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	agentRunnerRequestHash    = "sha256:0000000000000000000000000000000000000000000000000000000000000011"
	agentRunnerLaunchHash     = "sha256:0000000000000000000000000000000000000000000000000000000000000012"
	agentRunnerIdentityHash   = "sha256:0000000000000000000000000000000000000000000000000000000000000013"
	agentRunnerCompletionHash = "sha256:0000000000000000000000000000000000000000000000000000000000000014"
	agentRunnerTerminalHash   = "sha256:0000000000000000000000000000000000000000000000000000000000000015"
	agentRunnerResumeHash     = "sha256:0000000000000000000000000000000000000000000000000000000000000016"
	agentRunnerForeignHash    = "sha256:0000000000000000000000000000000000000000000000000000000000000017"
)

type agentRunnerFixture struct {
	engine     *Engine
	store      EventStore
	steps      *fakeSteps
	stopper    *fakeStopper
	acceptance *fakeAcceptance
	reviewer   *fakeReviewer
	publisher  *fakePublisher
}

// newAgentRunnerFixture runs the definition with the isolated-workspace step
// prepared as an AgentRunner step, so the step parks instead of passing.
func newAgentRunnerFixture(t *testing.T, store EventStore, tune func(*Options)) *agentRunnerFixture {
	t.Helper()
	if store == nil {
		store = &memoryEvents{}
	}
	options, _, _, steps, acceptance, reviewer, publisher, stopper, _ := testOptions(store)
	facts := AgentRunnerImmutableFacts{
		RequestID:         "request-one",
		WorkspaceHash:     baselineHash,
		ContextHash:       acceptedOne,
		AuthorizationHash: acceptedTwo,
		BudgetHash:        acceptedThree,
	}
	preparer := &fakeStepIntents{prepare: func(request StepRequest) (StepExecutionIntent, error) {
		if request.Step.Kind == "code" && request.Step.Mode == IsolatedWorkspace {
			current := facts
			return StepExecutionIntent{ExecutionTarget: AgentRunnerExecution, AgentRunnerFacts: &current}, nil
		}
		return StepExecutionIntent{}, nil
	}}
	options.StepIntents = preparer
	steps.evidence = agentRunnerDispatchEvidence()
	if tune != nil {
		tune(&options)
	}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	return &agentRunnerFixture{engine: engine, store: store, steps: steps, stopper: stopper, acceptance: acceptance, reviewer: reviewer, publisher: publisher}
}

func agentRunnerDispatchEvidence() []string {
	return []string{agentRunnerRequestHash, agentRunnerLaunchHash, agentRunnerIdentityHash}
}

// agentRunnerTerminal builds a valid terminal fact for the parked step. A
// completed terminal always carries the five frozen facts, because the chain
// refuses to route a completion without them.
func agentRunnerTerminal(status AgentRunnerTerminalStatus) AgentRunnerTerminal {
	terminal := AgentRunnerTerminal{
		RequestID:      "request-one",
		RequestHash:    agentRunnerRequestHash,
		CompletionHash: agentRunnerCompletionHash,
		RunID:          "run-one",
		TaskID:         "task-two",
		StepID:         "code-isolated",
		Attempt:        1,
		SessionID:      "session-one",
		Status:         status,
		Evidence:       []string{agentRunnerRequestHash, agentRunnerCompletionHash, agentRunnerTerminalHash},
	}
	if status == AgentRunnerTerminalCompleted {
		facts := testFrozenFacts()
		terminal.Facts = &facts
	}
	return terminal
}

// park dispatches the AgentRunner step and asserts that it waits instead of
// passing.
func (fixture *agentRunnerFixture) park(t *testing.T) {
	t.Helper()
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrAwaitingAgentRunnerTerminal) {
		t.Fatalf("expected the AgentRunner step to await its terminal, got %v", err)
	}
	fixture.requireStates(t, "TERMINAL_PENDING", "STEPS_RUNNING", "RUNNING")
}

func (fixture *agentRunnerFixture) requireStates(t *testing.T, wantStep, wantTask, wantChain string) {
	t.Helper()
	projection := fixture.engine.Projection()
	if got := projection.StepStates[stepKey("task-two", "code-isolated", 1)]; got != wantStep {
		t.Fatalf("step state %s, want %s (%+v)", got, wantStep, projection)
	}
	if got := projection.TaskStates[taskKey("task-two", 1)]; got != wantTask {
		t.Fatalf("task state %s, want %s (%+v)", got, wantTask, projection)
	}
	if got := projection.ChainState; got != wantChain {
		t.Fatalf("chain state %s, want %s (%+v)", got, wantChain, projection)
	}
}

func (fixture *agentRunnerFixture) eventCount() int {
	events, _ := fixture.store.Load(context.Background())
	return len(events)
}

func TestAgentRunnerTerminalValidateFailsClosed(t *testing.T) {
	tests := map[string]func(*AgentRunnerTerminal){
		"zero value":      func(terminal *AgentRunnerTerminal) { *terminal = AgentRunnerTerminal{} },
		"unknown status":  func(terminal *AgentRunnerTerminal) { terminal.Status = "succeeded" },
		"missing request": func(terminal *AgentRunnerTerminal) { terminal.RequestHash = "" },
		"missing completion": func(terminal *AgentRunnerTerminal) {
			terminal.CompletionHash = ""
		},
		"mixed mode":              func(terminal *AgentRunnerTerminal) { terminal.PriorSessionID = "session-one" },
		"resume session mismatch": func(terminal *AgentRunnerTerminal) { terminal.PriorSessionID = "session-two" },
		"invalid attempt":         func(terminal *AgentRunnerTerminal) { terminal.Attempt = 0 },
		"invalid session":         func(terminal *AgentRunnerTerminal) { terminal.SessionID = "session one" },
		"evidence duplicate": func(terminal *AgentRunnerTerminal) {
			terminal.Evidence = []string{agentRunnerRequestHash, agentRunnerRequestHash, agentRunnerCompletionHash}
		},
		"evidence without completion": func(terminal *AgentRunnerTerminal) {
			terminal.Evidence = []string{agentRunnerRequestHash, agentRunnerTerminalHash}
		},
		"evidence without request": func(terminal *AgentRunnerTerminal) {
			terminal.Evidence = []string{agentRunnerCompletionHash, agentRunnerTerminalHash}
		},
		"evidence with invalid digest": func(terminal *AgentRunnerTerminal) {
			terminal.Evidence = []string{agentRunnerRequestHash, agentRunnerCompletionHash, "not-a-hash"}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			terminal := agentRunnerTerminal(AgentRunnerTerminalCompleted)
			mutate(&terminal)
			if err := terminal.Validate(); !errors.Is(err, ErrInvalidAgentRunnerTerminal) {
				t.Fatalf("expected fail-closed validation, got %v", err)
			}
		})
	}
	resume := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	resume.PriorSessionID = "session-one"
	resume.PriorCompletionHash = agentRunnerResumeHash
	if err := resume.Validate(); err != nil {
		t.Fatalf("resume continuation must validate: %v", err)
	}
}

func TestAgentRunnerDispatchParksStepInsteadOfPassingIt(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	// Only task-one reached the downstream flow; the parked AgentRunner step did
	// not freeze, review, or promote anything.
	if fixture.acceptance.calls != 1 || fixture.reviewer.calls != 1 || fixture.publisher.calls != 1 {
		t.Fatalf("parked step reached downstream ports: accept=%d review=%d publish=%d", fixture.acceptance.calls, fixture.reviewer.calls, fixture.publisher.calls)
	}
	executions := len(fixture.steps.requests)
	events := fixture.eventCount()
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrAwaitingAgentRunnerTerminal) {
		t.Fatalf("second run must refuse to redispatch a parked step, got %v", err)
	}
	if len(fixture.steps.requests) != executions || fixture.eventCount() != events {
		t.Fatalf("re-entry redispatch or wrote state: executions=%d events=%d", len(fixture.steps.requests), fixture.eventCount())
	}
	fixture.requireStates(t, "TERMINAL_PENDING", "STEPS_RUNNING", "RUNNING")
}

func TestAgentRunnerTerminalRoutingMatrix(t *testing.T) {
	tests := map[AgentRunnerTerminalStatus]struct {
		step   string
		task   string
		chain  string
		reason string
	}{
		AgentRunnerTerminalCompleted:              {step: "PASSED", task: "STEPS_RUNNING", chain: "RUNNING", reason: agentRunnerTerminalPassedReason},
		AgentRunnerTerminalFailed:                 {step: "FAILED", task: "FAILED", chain: "FAILED", reason: agentRunnerTerminalFailedReason},
		AgentRunnerTerminalUncertain:              {step: "FAILED", task: "FAILED", chain: "PAUSED", reason: agentRunnerTerminalUncertainReason},
		AgentRunnerTerminalCancelled:              {step: "CANCELLED", task: "CANCELLED", chain: "CANCELLED", reason: agentRunnerTerminalCancelledReason},
		AgentRunnerTerminalOperatorActionRequired: {step: "WAITING_FOR_OPERATOR", task: "WAITING_FOR_OPERATOR", chain: "PAUSED", reason: agentRunnerTerminalWaitingReason},
	}
	for status, want := range tests {
		t.Run(string(status), func(t *testing.T) {
			fixture := newAgentRunnerFixture(t, nil, nil)
			fixture.park(t)
			acceptance, reviewer, publisher := fixture.acceptance.calls, fixture.reviewer.calls, fixture.publisher.calls
			route, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(status))
			if err != nil {
				t.Fatal(err)
			}
			if fixture.acceptance.calls != acceptance || fixture.reviewer.calls != reviewer || fixture.publisher.calls != publisher {
				t.Fatalf("terminal routing drove downstream work: accept=%d review=%d publish=%d", fixture.acceptance.calls, fixture.reviewer.calls, fixture.publisher.calls)
			}
			fixture.requireStates(t, want.step, want.task, want.chain)
			if route.Status != status || route.Converged || route.StepState != want.step || route.TaskState != want.task || route.ChainState != want.chain {
				t.Fatalf("unexpected route receipt: %+v", route)
			}
			events, err := fixture.store.Load(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if err := evidence.VerifyEventChain(events); err != nil {
				t.Fatal(err)
			}
			for _, record := range events {
				if record.Event.Reason.Code == want.reason {
					if !containsAll(record.Event.InputEvidence, agentRunnerRequestHash, agentRunnerCompletionHash) {
						t.Fatalf("routing evidence missing the request or completion digest: %+v", record.Event.InputEvidence)
					}
					return
				}
			}
			t.Fatalf("routing transition with reason %s not recorded", want.reason)
		})
	}
}

func TestAgentRunnerCompletedTerminalAloneDoesNotCompleteTask(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	acceptance, reviewer, publisher := fixture.acceptance.calls, fixture.reviewer.calls, fixture.publisher.calls
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	fixture.requireStates(t, "PASSED", "STEPS_RUNNING", "RUNNING")
	if fixture.acceptance.calls != acceptance || fixture.reviewer.calls != reviewer || fixture.publisher.calls != publisher {
		t.Fatalf("terminal receipt drove downstream work: accept=%d review=%d publish=%d", fixture.acceptance.calls, fixture.reviewer.calls, fixture.publisher.calls)
	}
	if err := fixture.engine.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fixture.acceptance.calls <= acceptance || fixture.reviewer.calls <= reviewer || fixture.publisher.calls <= publisher {
		t.Fatalf("downstream flow did not run after routing: accept=%d review=%d publish=%d", fixture.acceptance.calls, fixture.reviewer.calls, fixture.publisher.calls)
	}
	projection := fixture.engine.Projection()
	if projection.ChainState != "COMPLETED" || projection.TaskStates[taskKey("task-two", 1)] != "PASSED" {
		t.Fatalf("run did not complete through the existing flow: %+v", projection)
	}
	if err := evidence.VerifyEventChain(mustEvents(t, fixture.store)); err != nil {
		t.Fatal(err)
	}
}

func TestAgentRunnerTerminalRefusesUnboundAndConflictingReceipts(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	events := fixture.eventCount()
	foreign := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), foreign); !errors.Is(err, ErrAgentRunnerTerminalConflict) {
		t.Fatalf("expected a conflict for a step without a parked dispatch, got %v", err)
	}
	if fixture.eventCount() != events {
		t.Fatal("refused terminal wrote state")
	}
	fixture.park(t)
	events = fixture.eventCount()
	unbound := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	unbound.RequestHash = agentRunnerForeignHash
	unbound.Evidence = []string{agentRunnerForeignHash, agentRunnerCompletionHash}
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), unbound); !errors.Is(err, ErrAgentRunnerTerminalConflict) {
		t.Fatalf("expected a conflict for a terminal that does not match the dispatch, got %v", err)
	}
	foreignStep := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	foreignStep.StepID = "code-managed"
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), foreignStep); !errors.Is(err, ErrInvalidAgentRunnerTerminal) {
		t.Fatalf("expected a rejection for a non-AgentRunner step, got %v", err)
	}
	if fixture.eventCount() != events {
		t.Fatal("refused terminal wrote state")
	}
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	routed := fixture.eventCount()
	conflicting := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	conflicting.CompletionHash = agentRunnerResumeHash
	conflicting.Evidence = []string{agentRunnerRequestHash, agentRunnerResumeHash}
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), conflicting); !errors.Is(err, ErrAgentRunnerTerminalConflict) {
		t.Fatalf("expected a conflict for a second completion of the same request, got %v", err)
	}
	if fixture.eventCount() != routed {
		t.Fatal("conflicting terminal wrote state")
	}
	converged, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted))
	if err != nil || !converged.Converged {
		t.Fatalf("repeated terminal must converge: route=%+v err=%v", converged, err)
	}
	if fixture.eventCount() != routed {
		t.Fatal("converged terminal wrote state")
	}
	fixture.requireStates(t, "PASSED", "STEPS_RUNNING", "RUNNING")
}

func TestAgentRunnerTerminalConvergesAfterPartialWrite(t *testing.T) {
	store := &failOnceEvents{}
	fixture := newAgentRunnerFixture(t, store, nil)
	fixture.park(t)
	store.failAt = store.appendCalls + 2
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalFailed)); err == nil {
		t.Fatal("expected the partial write to surface")
	}
	fixture.requireStates(t, "FAILED", "STEPS_RUNNING", "RUNNING")
	route, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalFailed))
	if err != nil {
		t.Fatal(err)
	}
	if !route.Converged {
		t.Fatalf("partial route must converge on replay: %+v", route)
	}
	fixture.requireStates(t, "FAILED", "FAILED", "FAILED")
}

func TestAgentRunnerOperatorRouteSurvivesRecoveryAndInteractionOpen(t *testing.T) {
	store := &memoryEvents{}
	fixture := newAgentRunnerFixture(t, store, nil)
	fixture.park(t)
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalOperatorActionRequired)); err != nil {
		t.Fatal(err)
	}
	fixture.requireStates(t, "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR", "PAUSED")

	// A crash-and-restart before the interaction record exists must still be
	// able to hand control to the operator and converge.
	recovered := newAgentRunnerFixture(t, store, nil)
	if err := recovered.engine.Recover(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("a waiting step must be reported as uncertain: %v", err)
	}
	recovered.requireStates(t, "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR", "PAUSED")
	if recovered.engine.Projection().RecoveryUncertain {
		t.Fatal("the route pause must not be recorded as recovery uncertainty")
	}
	requestRecord, responseRecord, binding := agentRunnerOperatorRecords(t)
	if err := recovered.engine.OpenOperatorInteraction(context.Background(), requestRecord); err != nil {
		t.Fatal(err)
	}
	if _, err := recovered.engine.ResumeOperatorInteraction(context.Background(), requestRecord, responseRecord, binding); err != nil {
		t.Fatal(err)
	}
	recovered.requireStates(t, "RUNNING", "STEPS_RUNNING", "PAUSED")
	// The resumed run continues the same loop: it dispatches again and parks.
	if err := recovered.engine.Run(context.Background()); !errors.Is(err, ErrAwaitingAgentRunnerTerminal) {
		t.Fatalf("resumed run must dispatch again and park, got %v", err)
	}
	recovered.requireStates(t, "TERMINAL_PENDING", "STEPS_RUNNING", "RUNNING")
}

func TestAgentRunnerTerminalConvergesAfterTaskWriteBeforeChainWrite(t *testing.T) {
	store := &failOnceEvents{}
	fixture := newAgentRunnerFixture(t, store, nil)
	fixture.park(t)
	// Fail the third write of the route: step and task are recorded, the chain is not.
	store.failAt = store.appendCalls + 3
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalFailed)); err == nil {
		t.Fatal("expected the partial write to surface")
	}
	fixture.requireStates(t, "FAILED", "FAILED", "RUNNING")
	events := fixture.eventCount()
	route, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalFailed))
	if err != nil {
		t.Fatal(err)
	}
	if !route.Converged {
		t.Fatalf("partial route must converge on replay: %+v", route)
	}
	fixture.requireStates(t, "FAILED", "FAILED", "FAILED")
	if fixture.eventCount() != events+1 {
		t.Fatalf("convergence must write exactly the missing transition: %d events", fixture.eventCount()-events)
	}
}

func TestAgentRunnerTerminalResumeContinuityProof(t *testing.T) {
	fj := newAgentRunnerFixture(t, nil, nil)
	fj.park(t)
	if _, err := fj.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalOperatorActionRequired)); err != nil {
		t.Fatal(err)
	}
	requestRecord, responseRecord, binding := agentRunnerOperatorRecords(t)
	if err := fj.engine.OpenOperatorInteraction(context.Background(), requestRecord); err != nil {
		t.Fatal(err)
	}
	fj.requireStates(t, "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR", "PAUSED")
	if _, err := fj.engine.ResumeOperatorInteraction(context.Background(), requestRecord, responseRecord, binding); err != nil {
		t.Fatal(err)
	}
	fj.requireStates(t, "RUNNING", "STEPS_RUNNING", "PAUSED")
	// A repeated operator terminal must not undo the answer that already
	// consumed the waiting route.
	operatorTerminal := agentRunnerTerminal(AgentRunnerTerminalOperatorActionRequired)
	events := fj.eventCount()
	replay, err := fj.engine.SubmitAgentRunnerTerminal(context.Background(), operatorTerminal)
	if err != nil || !replay.Converged {
		t.Fatalf("superseded operator route must converge: route=%+v err=%v", replay, err)
	}
	fj.requireStates(t, "RUNNING", "STEPS_RUNNING", "PAUSED")
	if fj.eventCount() != events {
		t.Fatal("superseded operator replay wrote state")
	}
	fj.steps.evidence = []string{agentRunnerResumeHash, agentRunnerLaunchHash, agentRunnerIdentityHash}
	if err := fj.engine.Run(context.Background()); !errors.Is(err, ErrAwaitingAgentRunnerTerminal) {
		t.Fatalf("resumed step must dispatch again and park, got %v", err)
	}
	fj.requireStates(t, "TERMINAL_PENDING", "STEPS_RUNNING", "RUNNING")
	events = fj.eventCount()
	unproven := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	unproven.RequestHash = agentRunnerResumeHash
	unproven.PriorSessionID = "session-one"
	unproven.PriorCompletionHash = agentRunnerForeignHash
	unproven.Evidence = []string{agentRunnerResumeHash, agentRunnerCompletionHash}
	if _, err := fj.engine.SubmitAgentRunnerTerminal(context.Background(), unproven); !errors.Is(err, ErrAgentRunnerResumeContinuityUnproven) {
		t.Fatalf("expected the unproven resume to be refused, got %v", err)
	}
	if fj.eventCount() != events {
		t.Fatal("unproven resume wrote state")
	}
	proven := unproven
	proven.PriorCompletionHash = agentRunnerCompletionHash
	route, err := fj.engine.SubmitAgentRunnerTerminal(context.Background(), proven)
	if err != nil {
		t.Fatal(err)
	}
	if route.Status != AgentRunnerTerminalCompleted || route.Converged {
		t.Fatalf("proven resume must route: %+v", route)
	}
	fj.requireStates(t, "PASSED", "STEPS_RUNNING", "RUNNING")
}

func TestAgentRunnerTerminalCancelledArchiveFailuresStayParked(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	fixture.stopper.err = errors.New("stop unavailable")
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCancelled)); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("expected a stop-evidence failure, got %v", err)
	}
	fixture.requireStates(t, "TERMINAL_PENDING", "STEPS_RUNNING", "PAUSED")
	if fixture.stopper.calls != 1 {
		t.Fatalf("stop was not attempted once: %d", fixture.stopper.calls)
	}
}

func TestAgentRunnerTerminalCancelledArchivesStopEvidence(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCancelled)); err != nil {
		t.Fatal(err)
	}
	if fixture.stopper.calls != 1 {
		t.Fatalf("cancel route must archive stop evidence: %d calls", fixture.stopper.calls)
	}
	fixture.requireStates(t, "CANCELLED", "CANCELLED", "CANCELLED")
	events := mustEvents(t, fixture.store)
	for _, record := range events {
		if record.Event.Reason.Code == agentRunnerTerminalCancelledReason {
			if !containsAll(record.Event.InputEvidence, stopHash, agentRunnerCompletionHash) {
				t.Fatalf("cancel evidence missing stop or completion digest: %+v", record.Event.InputEvidence)
			}
		}
	}
	if err := evidence.VerifyEventChain(events); err != nil {
		t.Fatal(err)
	}
}

func TestAgentRunnerUncertainTerminalPausesWithoutRetry(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	executions := len(fixture.steps.requests)
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalUncertain)); err != nil {
		t.Fatal(err)
	}
	fixture.requireStates(t, "FAILED", "FAILED", "PAUSED")
	projection := fixture.engine.Projection()
	if !projection.RecoveryUncertain {
		t.Fatalf("uncertain terminal must stay unresumable: %+v", projection)
	}
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("uncertain run must refuse to resume, got %v", err)
	}
	if len(fixture.steps.requests) != executions {
		t.Fatal("uncertain terminal re-dispatched the step")
	}
}

func TestAgentRunnerRecoverTreatsParkedStepAsUncertain(t *testing.T) {
	store := &memoryEvents{}
	fixture := newAgentRunnerFixture(t, store, nil)
	fixture.park(t)
	recovered := newAgentRunnerFixture(t, store, nil)
	if err := recovered.engine.Recover(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("recovery of a parked dispatch must be uncertain: %v", err)
	}
	recovered.requireStates(t, "TERMINAL_PENDING", "STEPS_RUNNING", "PAUSED")
	if err := recovered.engine.Run(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("recovery pause must not resume by itself, got %v", err)
	}
	if _, err := recovered.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	recovered.requireStates(t, "PASSED", "STEPS_RUNNING", "PAUSED")
}

// agentRunnerOperatorRecords binds the operator-interaction machine to the
// AgentRunner step the terminal route parked.
func agentRunnerOperatorRecords(t *testing.T) (OperatorInteractionRecord, OperatorInteractionRecord, OperatorInteractionBinding) {
	t.Helper()
	request := validOperatorInteractionRequest()
	request.TaskID = "task-two"
	request.StepID = "code-isolated"
	request.Attempt = 1
	request.RequestID = "request-two"
	requestRecord, err := NewOperatorInteractionRequestRecord(request, fixedClock, func() string { return "request-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	response := OperatorInteractionResponse{
		InteractionID:  request.InteractionID,
		RespondedAt:    "2026-09-10T01:01:00.000Z",
		RespondedBy:    request.Operator,
		RunID:          request.RunID,
		TaskID:         request.TaskID,
		StepID:         request.StepID,
		Attempt:        request.Attempt,
		WorkspaceHash:  request.WorkspaceHash,
		ConversationID: request.ConversationID,
		RequestID:      "request-three",
		ContextHash:    request.ContextHash,
		RequestHash:    requestRecord.RecordHash,
		Selection:      request.AllowedResponses[0],
		Evidence:       []string{acceptedTwo},
	}
	responseRecord, err := NewOperatorInteractionResponseRecord(requestRecord, response, fixedClock, func() string { return "response-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	binding := OperatorInteractionBinding{
		RunID: request.RunID, TaskID: request.TaskID, StepID: request.StepID, Attempt: request.Attempt,
		WorkspaceHash: request.WorkspaceHash, ConversationID: request.ConversationID, ContextHash: request.ContextHash,
	}
	return requestRecord, responseRecord, binding
}

func containsAll(values []string, want ...string) bool {
	for _, target := range want {
		found := false
		for _, value := range values {
			if value == target {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func mustEvents(t *testing.T, store EventStore) []evidence.StateEvent {
	t.Helper()
	events, err := store.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return events
}
