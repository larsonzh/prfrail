package chain

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	postflightManifestHash = "sha256:0000000000000000000000000000000000000000000000000000000000000021"
	postflightDiffHash     = "sha256:0000000000000000000000000000000000000000000000000000000000000022"
	postflightLogHash      = "sha256:0000000000000000000000000000000000000000000000000000000000000023"
	postflightUsageHash    = "sha256:0000000000000000000000000000000000000000000000000000000000000024"
	postflightStopHash     = "sha256:0000000000000000000000000000000000000000000000000000000000000025"
)

func testFrozenFacts() AgentRunnerFrozenFacts {
	return AgentRunnerFrozenFacts{
		ManifestHash:            postflightManifestHash,
		DiffHash:                postflightDiffHash,
		LogHash:                 postflightLogHash,
		UsageHash:               postflightUsageHash,
		ProcessStopEvidenceHash: postflightStopHash,
	}
}

func TestAgentRunnerFrozenFactsValidateFailsClosed(t *testing.T) {
	tests := map[string]func(*AgentRunnerFrozenFacts){
		"missing manifest":  func(facts *AgentRunnerFrozenFacts) { facts.ManifestHash = "" },
		"missing diff":      func(facts *AgentRunnerFrozenFacts) { facts.DiffHash = "" },
		"missing log":       func(facts *AgentRunnerFrozenFacts) { facts.LogHash = "" },
		"missing usage":     func(facts *AgentRunnerFrozenFacts) { facts.UsageHash = "" },
		"missing stop":      func(facts *AgentRunnerFrozenFacts) { facts.ProcessStopEvidenceHash = "" },
		"invalid digest":    func(facts *AgentRunnerFrozenFacts) { facts.LogHash = "not-a-hash" },
		"duplicate digest":  func(facts *AgentRunnerFrozenFacts) { facts.DiffHash = facts.ManifestHash },
		"collision on stop": func(facts *AgentRunnerFrozenFacts) { facts.ProcessStopEvidenceHash = facts.UsageHash },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			facts := testFrozenFacts()
			mutate(&facts)
			if err := facts.Validate(); !errors.Is(err, ErrInvalidPostflightFacts) {
				t.Fatalf("expected fail-closed fact validation, got %v", err)
			}
		})
	}
	if err := testFrozenFacts().Validate(); err != nil {
		t.Fatalf("valid facts must pass: %v", err)
	}
	if err := testFrozenFacts().DisjointFrom(agentRunnerRequestHash); err != nil {
		t.Fatalf("disjoint terminal binding must pass: %v", err)
	}
	if err := testFrozenFacts().DisjointFrom(postflightUsageHash); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("expected a collision rejection, got %v", err)
	}
}

func TestFrozenFactsRebuildFromRoutingEvidence(t *testing.T) {
	facts := testFrozenFacts()
	evidence := append([]string{agentRunnerRequestHash, agentRunnerCompletionHash}, agentRunnerFactsEvidence(facts)...)
	evidence = append(evidence, agentRunnerTerminalHash)
	rebuilt, err := frozenFactsFromEvidence(evidence)
	if err != nil {
		t.Fatal(err)
	}
	if rebuilt != facts {
		t.Fatalf("rebuilt facts mismatch: %+v", rebuilt)
	}
	if _, err := frozenFactsFromEvidence(evidence[:6]); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("short evidence must fail closed, got %v", err)
	}
	ambiguous := append([]string{agentRunnerRequestHash, agentRunnerCompletionHash}, agentRunnerFactsEvidence(facts)...)
	ambiguous[postflightFactsPrefix+1] = ambiguous[postflightFactsPrefix]
	if _, err := frozenFactsFromEvidence(ambiguous); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("ambiguous evidence must fail closed, got %v", err)
	}
}

// installPostflight swaps the fixture's port for a scripted one.
func (fixture *agentRunnerFixture) installPostflight(port *fakePostflight) *fakePostflight {
	fixture.engine.options.Postflight = port
	return port
}

func TestEngineRequiresAPostflightPort(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, _, _, _, _, _, _ := testOptions(store)
	options.Postflight = nil
	if _, err := New(context.Background(), options); !errors.Is(err, ErrPostflightUnavailable) {
		t.Fatalf("a missing postflight port must fail closed, got %v", err)
	}
}

func TestPostflightGateConvergesAfterReviewPendingWasPersisted(t *testing.T) {
	store := &failOnceEvents{}
	fixture := newAgentRunnerFixture(t, store, nil)
	fixture.park(t)
	port := fixture.installPostflight(&fakePostflight{})
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	// Fail the write that follows REVIEW_PENDING (the task-passed event), so the
	// task is left resting in REVIEW_PENDING exactly like a crash between the
	// review write and acceptance.
	store.failAt = store.appendCalls + 2
	if err := fixture.engine.Run(context.Background()); err == nil {
		t.Fatal("expected the injected failure to surface")
	}
	fixture.requireStates(t, "PASSED", "REVIEW_PENDING", "RUNNING")
	if err := fixture.engine.Run(context.Background()); err != nil {
		t.Fatalf("re-entry into REVIEW_PENDING must converge, got %v", err)
	}
	fixture.requireStates(t, "PASSED", "PASSED", "COMPLETED")
	if port.calls != 1 {
		t.Fatalf("the gate must not re-run once the review state is durable: %d calls", port.calls)
	}
	if err := evidence.VerifyEventChain(mustEvents(t, fixture.store)); err != nil {
		t.Fatal(err)
	}
}

// appendTaskReview records a task-level REVIEW_PENDING transition, the way a
// build that wrote the review state without running the gate would.
func (fixture *agentRunnerFixture) appendTaskReview(t *testing.T, inputs []string) {
	t.Helper()
	records := mustEvents(t, fixture.store)
	previous := records[len(records)-1].EventHash
	predecessor := previous
	if inputs == nil {
		inputs = []string{}
	}
	record, err := evidence.NewStateEvent(evidence.Event{
		EventID:           fmt.Sprintf("event-task-review-%d", len(records)+1),
		RunID:             "run-one",
		Sequence:          len(records) + 1,
		OccurredAt:        fixedClock().Format("2006-01-02T15:04:05.000Z"),
		Actor:             evidence.Actor{Type: "system", ID: "engine"},
		Entity:            evidence.Entity{Kind: "task", RunID: "run-one", TaskID: "task-two", Attempt: 1},
		FromState:         "STEPS_RUNNING",
		ToState:           "REVIEW_PENDING",
		PreviousEventHash: &predecessor,
		InputEvidence:     inputs,
		Reason:            evidence.Reason{Code: "review-pending"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := fixture.store.Append(context.Background(), record); err != nil {
		t.Fatal(err)
	}
}

// TestRequirePostflightQualifiedReviewFailsClosed covers the three shapes of the
// re-entry proof: no recorded review transition, a transition without the frozen
// facts, and a transition that carries them.
func TestRequirePostflightQualifiedReviewFailsClosed(t *testing.T) {
	tests := map[string]struct {
		inputs []string
		append bool
		want   error
	}{
		"no review transition": {append: false, want: ErrUnqualifiedReview},
		"without the facts":    {inputs: []string{agentRunnerTerminalHash}, append: true, want: ErrUnqualifiedReview},
		"with the facts":       {inputs: agentRunnerFactsEvidence(testFrozenFacts()), append: true, want: nil},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fixture := newAgentRunnerFixture(t, nil, nil)
			fixture.park(t)
			if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
				t.Fatal(err)
			}
			if test.append {
				fixture.appendTaskReview(t, test.inputs)
				if err := evidence.VerifyEventChain(mustEvents(t, fixture.store)); err != nil {
					t.Fatalf("the appended review transition must be a valid record: %v", err)
				}
			}
			if err := fixture.engine.requirePostflightQualifiedReview(context.Background(), "task-two"); !errors.Is(err, test.want) {
				t.Fatalf("guard returned %v, want %v", err, test.want)
			}
		})
	}
}

// TestPostflightGateRefusesAnUnqualifiedReviewState pins the re-entry boundary: a
// task already resting in REVIEW_PENDING may only continue when its review
// transition proves the gate produced it. A review state written without the gate
// is refused with zero writes instead of silently reaching acceptance.
func TestPostflightGateRefusesAnUnqualifiedReviewState(t *testing.T) {
	store := &failOnceEvents{}
	fixture := newAgentRunnerFixture(t, store, nil)
	fixture.park(t)
	port := fixture.installPostflight(&fakePostflight{})
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	store.failAt = store.appendCalls + 2
	if err := fixture.engine.Run(context.Background()); err == nil {
		t.Fatal("expected the injected failure to surface")
	}
	fixture.requireStates(t, "PASSED", "REVIEW_PENDING", "RUNNING")
	// Strip the frozen-fact proof from the recorded review transition. The event
	// chain is deliberately left inconsistent here: the point is that the gate
	// refuses a review state whose qualification it cannot prove.
	rewritten := false
	for index := range store.events {
		event := &store.events[index].Event
		if event.Entity.Kind == "task" && event.ToState == "REVIEW_PENDING" {
			event.InputEvidence = []string{}
			rewritten = true
		}
	}
	if !rewritten {
		t.Fatal("the fixture recorded no review transition to strip")
	}
	acceptance, reviewer, publisher := fixture.acceptance.calls, fixture.reviewer.calls, fixture.publisher.calls
	events := fixture.eventCount()
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrUnqualifiedReview) {
		t.Fatalf("an unproven review state must be refused, got %v", err)
	}
	fixture.requireStates(t, "PASSED", "REVIEW_PENDING", "RUNNING")
	if port.calls != 1 {
		t.Fatalf("the gate must not re-run for a task already in review: %d calls", port.calls)
	}
	if fixture.eventCount() != events || fixture.acceptance.calls != acceptance || fixture.reviewer.calls != reviewer || fixture.publisher.calls != publisher {
		t.Fatal("the refusal wrote state or reached the acceptance flow")
	}
}

// TestPostflightGateHandsAParentDigestFactSetToThePort pins the boundary that
// only the port can discriminate: a no-change execution legitimately re-captures
// the parent manifest, so a fact set whose manifest digest equals the accepted
// parent snapshot digest must still reach the port. The chain cannot tell such a
// no-op from an adapter echoing the parent, so it neither refuses it beforehand
// nor lets it pass on its own: the port's own verdict decides.
func TestPostflightGateHandsAParentDigestFactSetToThePort(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	port := fixture.installPostflight(&fakePostflight{decisions: []PostflightDecision{{Outcome: PostflightFailed, Evidence: agentRunnerFactsEvidence(testFrozenFacts())}}})
	parentDigest := testFrozenFacts()
	parentDigest.ManifestHash = acceptedOne
	terminal := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	terminal.Facts = &parentDigest
	if err := terminal.Validate(); err != nil {
		t.Fatalf("the chain must not refuse a parent-equal manifest before the port: %v", err)
	}
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), terminal); err != nil {
		t.Fatal(err)
	}
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrPaused) {
		t.Fatalf("the port verdict must decide the outcome, got %v", err)
	}
	fixture.requireStates(t, "PASSED", "REPAIR_PENDING", "PAUSED")
	if port.calls != 1 || len(port.requests) != 1 || port.requests[0].Facts != parentDigest {
		t.Fatalf("the port must reconcile the parent-equal fact set exactly once: %+v", port.requests)
	}
}

func TestAgentRunnerTerminalFactsMustNotCollideWithItsBindings(t *testing.T) {
	terminal := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	facts := testFrozenFacts()
	facts.ManifestHash = terminal.RequestHash
	terminal.Facts = &facts
	if err := terminal.Validate(); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("a fact digest equal to the request digest must be refused, got %v", err)
	}
	facts = testFrozenFacts()
	facts.DiffHash = terminal.CompletionHash
	terminal.Facts = &facts
	if err := terminal.Validate(); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("a fact digest equal to the completion digest must be refused, got %v", err)
	}
}

func TestAgentRunnerReplayWithDifferentFactsConflicts(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	first := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	events := fixture.eventCount()
	divergent := testFrozenFacts()
	divergent.LogHash = agentRunnerResumeHash
	replay := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	replay.Facts = &divergent
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), replay); !errors.Is(err, ErrAgentRunnerTerminalConflict) {
		t.Fatalf("a replay with a different fact set must conflict, got %v", err)
	}
	if fixture.eventCount() != events {
		t.Fatal("the conflicting replay wrote state")
	}
	converged, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), first)
	if err != nil || !converged.Converged {
		t.Fatalf("the identical replay must still converge: route=%+v err=%v", converged, err)
	}
}

func TestPostflightDecisionRejectsMalformedEvidence(t *testing.T) {
	facts := testFrozenFacts()
	// Every fact is present and only one extra digest is malformed, so nothing
	// but the evidence-digest loop can reject this decision.
	salted := PostflightDecision{Outcome: PostflightPassed, Evidence: append(agentRunnerFactsEvidence(facts), "not-a-hash")}
	if err := salted.Validate(facts); !errors.Is(err, ErrInvalidPostflightDecision) {
		t.Fatalf("a malformed evidence digest must be refused even when every fact is present, got %v", err)
	}
	incomplete := PostflightDecision{Outcome: PostflightPassed, Evidence: agentRunnerFactsEvidence(facts)[:4]}
	if err := incomplete.Validate(facts); !errors.Is(err, ErrInvalidPostflightDecision) {
		t.Fatalf("a passed decision without every fact must be refused, got %v", err)
	}
	malformedError := PostflightDecision{Outcome: PostflightFailed, ErrorEvidence: []string{"not-a-hash"}}
	if err := malformedError.Validate(facts); !errors.Is(err, ErrInvalidPostflightDecision) {
		t.Fatalf("a malformed error evidence digest must be refused, got %v", err)
	}
}

func TestPostflightGatePassesAndBindsTheFrozenFacts(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	port := fixture.installPostflight(&fakePostflight{})
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	if port.calls != 0 {
		t.Fatal("postflight must not run before the task is ready for review")
	}
	if err := fixture.engine.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if port.calls != 1 || len(port.requests) != 1 {
		t.Fatalf("postflight must run exactly once: %d", port.calls)
	}
	if port.requests[0].Facts != testFrozenFacts() || port.requests[0].TaskID != "task-two" || port.requests[0].ParentSnapshotHash != acceptedOne {
		t.Fatalf("unexpected postflight request: %+v", port.requests[0])
	}
	// task-two passed through the gate and the existing flow, and task-three
	// then completed the run through the same flow.
	if fixture.acceptance.calls != 3 || fixture.reviewer.calls != 3 || fixture.publisher.calls != 3 {
		t.Fatalf("acceptance flow did not run exactly once per task: accept=%d review=%d publish=%d", fixture.acceptance.calls, fixture.reviewer.calls, fixture.publisher.calls)
	}
	want := []string{"task-one", "task-two", "task-three"}
	for name, got := range map[string][]string{"acceptance": fixture.acceptance.tasks, "review": fixture.reviewer.tasks, "promotion": fixture.publisher.tasks} {
		if len(got) != len(want) {
			t.Fatalf("%s ran for %v, want %v", name, got, want)
		}
		for index, taskID := range want {
			if got[index] != taskID {
				t.Fatalf("%s order %v, want %v", name, got, want)
			}
		}
	}
	projection := fixture.engine.Projection()
	if projection.TaskStates[taskKey("task-two", 1)] != "PASSED" || projection.ChainState != "COMPLETED" {
		t.Fatalf("run did not complete after a passed postflight: %+v", projection)
	}
	for _, record := range mustEvents(t, fixture.store) {
		if record.Event.Reason.Code != "postflight-passed" {
			continue
		}
		if record.Event.ToState != "REVIEW_PENDING" || !containsAll(record.Event.InputEvidence, agentRunnerFactsEvidence(testFrozenFacts())...) {
			t.Fatalf("review transition does not prove the frozen facts: %+v", record.Event)
		}
		return
	}
	t.Fatal("the postflight gate recorded no review transition")
}

func TestPostflightGateRefusesAnUnprovenPass(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	incomplete := agentRunnerFactsEvidence(testFrozenFacts())[:4]
	fixture.installPostflight(&fakePostflight{decisions: []PostflightDecision{{Outcome: PostflightPassed, Evidence: incomplete}}})
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	acceptance := fixture.acceptance.calls
	if err := fixture.engine.Run(context.Background()); err == nil {
		t.Fatal("a passed verdict without every fact digest must be refused")
	}
	fixture.requireStates(t, "PASSED", "FAILED", "FAILED")
	if fixture.acceptance.calls != acceptance {
		t.Fatal("an unproven pass reached the acceptance flow")
	}
}

func TestPostflightGateRejectionParksTheTaskForRepair(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	fixture.installPostflight(&fakePostflight{decisions: []PostflightDecision{{Outcome: PostflightFailed, Evidence: agentRunnerFactsEvidence(testFrozenFacts())}}})
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	acceptance := fixture.acceptance.calls
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrPaused) {
		t.Fatalf("a postflight rejection must pause the run, got %v", err)
	}
	fixture.requireStates(t, "PASSED", "REPAIR_PENDING", "PAUSED")
	if fixture.acceptance.calls != acceptance {
		t.Fatal("a rejected postflight reached the acceptance flow")
	}
}

func TestPostflightGateUncertainFailsAndPausesWithoutRetry(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	port := fixture.installPostflight(&fakePostflight{decisions: []PostflightDecision{{Outcome: PostflightUncertain, ErrorEvidence: []string{agentRunnerTerminalHash}}}})
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("an uncertain postflight must not resume, got %v", err)
	}
	fixture.requireStates(t, "PASSED", "FAILED", "PAUSED")
	if !fixture.engine.Projection().RecoveryUncertain {
		t.Fatal("an uncertain postflight must keep the run unresumable")
	}
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("a second run must stay refused, got %v", err)
	}
	if port.calls != 1 {
		t.Fatalf("an uncertain postflight must not be retried: %d calls", port.calls)
	}
}

func TestPostflightGatePortFailureFailsTheTask(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	port := fixture.installPostflight(&fakePostflight{err: errors.New("postflight unavailable")})
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	if err := fixture.engine.Run(context.Background()); err == nil {
		t.Fatal("a postflight port failure must fail the task")
	}
	fixture.requireStates(t, "PASSED", "FAILED", "FAILED")
	if len(port.requests) != 1 {
		t.Fatalf("unexpected postflight attempts: %d", len(port.requests))
	}
}

func TestPostflightGateConvergesAfterPartialWrite(t *testing.T) {
	store := &failOnceEvents{}
	fixture := newAgentRunnerFixture(t, store, nil)
	fixture.park(t)
	port := fixture.installPostflight(&fakePostflight{})
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	store.failAt = store.appendCalls + 1
	if err := fixture.engine.Run(context.Background()); err == nil {
		t.Fatal("expected the partial review write to surface")
	}
	fixture.requireStates(t, "PASSED", "STEPS_RUNNING", "RUNNING")
	if port.calls != 1 || fixture.acceptance.calls != 1 {
		t.Fatalf("the failed review write must not reach acceptance: postflight=%d accept=%d", port.calls, fixture.acceptance.calls)
	}
	if err := fixture.engine.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	fixture.requireStates(t, "PASSED", "PASSED", "COMPLETED")
	// The gate re-ran once for task-two, whose acceptance flow then completed
	// together with task-three's own flow.
	if port.calls != 2 || fixture.acceptance.calls != 3 || fixture.reviewer.calls != 3 || fixture.publisher.calls != 3 {
		t.Fatalf("re-entry did not converge exactly once: postflight=%d accept=%d review=%d publish=%d", port.calls, fixture.acceptance.calls, fixture.reviewer.calls, fixture.publisher.calls)
	}
	if err := evidence.VerifyEventChain(mustEvents(t, fixture.store)); err != nil {
		t.Fatal(err)
	}
}

func TestCompletedTerminalWithoutFrozenFactsIsRefused(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	terminal := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	terminal.Facts = nil
	events := fixture.eventCount()
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), terminal); !errors.Is(err, ErrInvalidAgentRunnerTerminal) {
		t.Fatalf("a completion without frozen facts must be refused, got %v", err)
	}
	if fixture.eventCount() != events {
		t.Fatal("the refused completion wrote state")
	}
	fixture.requireStates(t, "TERMINAL_PENDING", "STEPS_RUNNING", "RUNNING")
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrAwaitingAgentRunnerTerminal) {
		t.Fatalf("a parked step must never pass on its own, got %v", err)
	}
	fixture.requireStates(t, "TERMINAL_PENDING", "STEPS_RUNNING", "RUNNING")
}

func TestPostflightDecisionValidate(t *testing.T) {
	facts := testFrozenFacts()
	passed := PostflightDecision{Outcome: PostflightPassed, Evidence: agentRunnerFactsEvidence(facts)}
	if err := passed.Validate(facts); err != nil {
		t.Fatalf("a passed decision carrying every fact must validate: %v", err)
	}
	incomplete := PostflightDecision{Outcome: PostflightPassed, Evidence: agentRunnerFactsEvidence(facts)[:4]}
	if err := incomplete.Validate(facts); !errors.Is(err, ErrInvalidPostflightDecision) {
		t.Fatalf("a passed decision without every fact must be refused, got %v", err)
	}
	unknown := PostflightDecision{Outcome: "approved", Evidence: agentRunnerFactsEvidence(facts)}
	if err := unknown.Validate(facts); !errors.Is(err, ErrInvalidPostflightDecision) {
		t.Fatalf("an unknown outcome must be refused, got %v", err)
	}
	for _, outcome := range []PostflightOutcome{PostflightFailed, PostflightUncertain} {
		decision := PostflightDecision{Outcome: outcome, ErrorEvidence: []string{agentRunnerTerminalHash}}
		if err := decision.Validate(facts); err != nil {
			t.Fatalf("outcome %s must validate without the fact digests: %v", outcome, err)
		}
	}
}

// appendStepRouting records a chained completed routing for another step of the
// same task, exactly like an adapter that routed a second AgentRunner step.
func (fixture *agentRunnerFixture) appendStepRouting(t *testing.T, stepID string, factSet AgentRunnerFrozenFacts) {
	t.Helper()
	records := mustEvents(t, fixture.store)
	previous := records[len(records)-1].EventHash
	sequence := len(records)
	for _, transition := range []struct{ from, to, reason string }{
		{"NONE", "PENDING", "step-pending"},
		{"PENDING", "RUNNING", "step-running"},
		{"RUNNING", "TERMINAL_PENDING", "step-terminal-pending"},
		{"TERMINAL_PENDING", "PASSED", agentRunnerTerminalPassedReason},
	} {
		sequence++
		predecessor := previous
		event := evidence.Event{
			EventID:           fmt.Sprintf("event-%s-%d", stepID, sequence),
			RunID:             "run-one",
			Sequence:          sequence,
			OccurredAt:        fixedClock().Format("2006-01-02T15:04:05.000Z"),
			Actor:             evidence.Actor{Type: "agent", ID: "agent-runner"},
			Entity:            evidence.Entity{Kind: "step", RunID: "run-one", TaskID: "task-two", StepID: stepID, Attempt: 1},
			FromState:         transition.from,
			ToState:           transition.to,
			PreviousEventHash: &predecessor,
			InputEvidence:     []string{},
			Reason:            evidence.Reason{Code: transition.reason},
		}
		if transition.to == "PASSED" {
			event.InputEvidence = append([]string{agentRunnerRequestHash, agentRunnerCompletionHash}, agentRunnerFactsEvidence(factSet)...)
		}
		record, err := evidence.NewStateEvent(event)
		if err != nil {
			t.Fatal(err)
		}
		if err := fixture.store.Append(context.Background(), record); err != nil {
			t.Fatal(err)
		}
		previous = record.EventHash
	}
}

// TestPostflightGateRefusesTasksWithSeveralRoutedSteps pins the rule that a task
// whose completed terminals came from more than one step is refused whole
// instead of guessing which facts the postflight should audit.
func TestPostflightGateRefusesTasksWithSeveralRoutedSteps(t *testing.T) {
	fixture := newAgentRunnerFixture(t, nil, nil)
	fixture.park(t)
	port := fixture.installPostflight(&fakePostflight{})
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), agentRunnerTerminal(AgentRunnerTerminalCompleted)); err != nil {
		t.Fatal(err)
	}
	other := testFrozenFacts()
	other.LogHash = agentRunnerResumeHash
	fixture.appendStepRouting(t, "code-isolated-b", other)
	if err := evidence.VerifyEventChain(mustEvents(t, fixture.store)); err != nil {
		t.Fatalf("the second routing must be a valid record: %v", err)
	}
	if err := fixture.engine.Run(context.Background()); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("facts routed by several steps must be refused, got %v", err)
	}
	fixture.requireStates(t, "PASSED", "FAILED", "FAILED")
	if port.calls != 0 {
		t.Fatalf("ambiguous facts must be refused before the port runs: %d calls", port.calls)
	}
}

// TestAgentRunnerReplayWithDamagedRoutingEvidenceConflicts covers the fail-closed
// path only a damaged event store can reach: a completed routing whose recorded
// evidence is too short to rebuild the facts may never converge.
func TestAgentRunnerReplayWithDamagedRoutingEvidenceConflicts(t *testing.T) {
	store := &memoryEvents{}
	fixture := newAgentRunnerFixture(t, store, nil)
	fixture.park(t)
	terminal := agentRunnerTerminal(AgentRunnerTerminalCompleted)
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), terminal); err != nil {
		t.Fatal(err)
	}
	damaged := false
	for index := range store.events {
		if store.events[index].Event.Reason.Code != agentRunnerTerminalPassedReason {
			continue
		}
		store.events[index].Event.InputEvidence = store.events[index].Event.InputEvidence[:postflightFactsPrefix+2]
		damaged = true
	}
	if !damaged {
		t.Fatal("the fixture recorded no routing event to damage")
	}
	if _, err := fixture.engine.SubmitAgentRunnerTerminal(context.Background(), terminal); !errors.Is(err, ErrAgentRunnerTerminalConflict) {
		t.Fatalf("a completed replay without rebuildable facts must conflict, got %v", err)
	}
}
