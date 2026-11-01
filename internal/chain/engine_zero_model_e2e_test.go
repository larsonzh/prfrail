package chain

import (
	"context"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// P45 zero-model E2E at the engine level.
//
// Every port in these tests is a package-local fake (baseline, workspaces, steps,
// acceptance, review, publish, postflight, stop, reconcile), so a green run proves
// the chain assembly itself drives T1..T3 without reaching a model, a network or an
// external process. The tests reuse testOptions/threeTaskDefinition from
// engine_test.go on purpose: the same fixture the positive engine tests use is what
// the zero-model chain must accept.

// TestZeroModelE2ET1PassT2FailT3NeverStarted pins the stop semantics: a failing T2
// must fail the chain immediately and T3 must never be materialized or dispatched.
func TestZeroModelE2ET1PassT2FailT3NeverStarted(t *testing.T) {
	store := &memoryEvents{}
	options, _, workspaces, steps, _, _, _, _, _ := testOptions(store)
	steps.failTask = "task-two"

	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	runErr := engine.Run(context.Background())
	if runErr == nil {
		t.Fatal("a failed task must fail the run, got a nil error")
	}

	projection := engine.Projection()
	if got := projection.TaskStates[taskKey("task-one", 1)]; got != "PASSED" {
		t.Fatalf("task-one must pass before the failure, got %q", got)
	}
	if got := projection.TaskStates[taskKey("task-two", 1)]; got != "FAILED" {
		t.Fatalf("task-two must carry the failure, got %q", got)
	}
	if projection.ChainState != "FAILED" {
		t.Fatalf("chain must stop FAILED, got %q", projection.ChainState)
	}

	// T3 never starts: the workspace port is the first thing runTask calls.
	if len(workspaces.calls) != 2 {
		t.Fatalf("only T1 and T2 may materialize a workspace, got %+v", workspaces.calls)
	}
	if workspaces.calls[0].task != "task-one" || workspaces.calls[1].task != "task-two" {
		t.Fatalf("workspace order must be T1 then T2, got %+v", workspaces.calls)
	}
	for _, call := range workspaces.calls {
		if call.task == "task-three" {
			t.Fatalf("task-three must never be materialized: %+v", call)
		}
	}

	// No task-level or step-level event may name task-three.
	for _, record := range store.events {
		if record.Event.Entity.TaskID == "task-three" {
			t.Fatalf("task-three wrote event %q (%s); it must never start", record.Event.EventID, record.Event.Reason.Code)
		}
	}
	if _, exists := projection.TaskStates[taskKey("task-three", 1)]; exists {
		t.Fatalf("task-three must have no projected task state: %+v", projection.TaskStates)
	}
	for _, request := range steps.requests {
		if request.TaskID == "task-three" {
			t.Fatalf("task-three step %q must never be dispatched", request.Step.ID)
		}
	}
	if len(steps.requests) != 3 {
		// T1 dispatches two steps (code-managed, build-one), T2 dispatches one
		// (code-isolated) and then fails; the T2 noop never reaches the step port.
		t.Fatalf("expected exactly three dispatched steps, got %d: %+v", len(steps.requests), steps.requests)
	}
}

// TestZeroModelT1ParentBinding pins the accepted-parent chain: T1 starts from the
// baseline and every later task starts from the snapshot its predecessor accepted.
func TestZeroModelT1ParentBinding(t *testing.T) {
	store := &memoryEvents{}
	options, _, workspaces, _, _, _, _, _, _ := testOptions(store)

	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatalf("the zero-model chain must complete: %v", err)
	}

	if len(workspaces.calls) != 3 {
		t.Fatalf("three tasks must materialize three workspaces, got %+v", workspaces.calls)
	}
	if workspaces.calls[0].task != "task-one" || workspaces.calls[0].parent != baselineHash {
		t.Fatalf("T1 must start from the captured baseline %s, got %+v", baselineHash, workspaces.calls[0])
	}
	if workspaces.calls[1].task != "task-two" || workspaces.calls[1].parent != acceptedOne {
		t.Fatalf("T2 must start from the T1 accepted hash %s, got %+v", acceptedOne, workspaces.calls[1])
	}
	if workspaces.calls[2].task != "task-three" || workspaces.calls[2].parent != acceptedTwo {
		t.Fatalf("T3 must start from the T2 accepted hash %s, got %+v", acceptedTwo, workspaces.calls[2])
	}
	if got := engine.Projection().AcceptedParent; got != acceptedThree {
		t.Fatalf("the run must end on the T3 accepted hash %s, got %s", acceptedThree, got)
	}
}

// TestZeroModelEventChainSealedMonotonic proves the audit trail the zero-model run
// leaves behind is verifiable: a hash-linked, gapless event log plus one seal per
// accepted task.
func TestZeroModelEventChainSealedMonotonic(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, _, _, _, _, _, _ := testOptions(store)

	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatalf("the zero-model chain must complete: %v", err)
	}
	if len(store.events) == 0 {
		t.Fatal("a completed run must write events")
	}
	if err := evidence.VerifyEventChain(store.events); err != nil {
		t.Fatalf("the recorded event chain must verify: %v", err)
	}
	previous := 0
	for index, record := range store.events {
		if record.Event.Sequence != index+1 {
			t.Fatalf("event %d carries sequence %d, want %d", index, record.Event.Sequence, index+1)
		}
		if record.Event.Sequence <= previous {
			t.Fatalf("sequence must increase at index %d: %d after %d", index, record.Event.Sequence, previous)
		}
		previous = record.Event.Sequence
	}
	projection := engine.Projection()
	if projection.Sequence != len(store.events) {
		t.Fatalf("projection sequence %d must equal the recorded event count %d", projection.Sequence, len(store.events))
	}
	if projection.LastEventHash != store.events[len(store.events)-1].EventHash {
		t.Fatalf("projection last hash %s must be the final event hash %s", projection.LastEventHash, store.events[len(store.events)-1].EventHash)
	}
	seals := 0
	for _, record := range store.seals {
		if record.ReviewDecision == "approve" && record.PromotionDecision == "completed" {
			seals++
		}
	}
	if seals != 3 {
		t.Fatalf("each accepted task must leave a sealed closure, got %d of 3", seals)
	}
}

// TestZeroModelNoAgentRunnerDispatch proves the zero-model graph has no model
// channel: without a StepIntentPreparer the engine never derives an AgentRunner
// target, so every dispatched step stays on the default execution.
func TestZeroModelNoAgentRunnerDispatch(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, steps, _, _, _, _, _ := testOptions(store)
	if options.StepIntents != nil {
		t.Fatal("the zero-model graph must not assemble a step-intent preparer: that is the only way an AgentRunner target can appear")
	}

	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatalf("the zero-model chain must complete: %v", err)
	}
	if len(steps.requests) == 0 {
		t.Fatal("the run dispatched no step, so this test proves nothing")
	}
	for _, request := range steps.requests {
		if request.ExecutionTarget != DefaultExecution {
			t.Fatalf("step %s/%s must target the default execution, got %q", request.TaskID, request.Step.ID, request.ExecutionTarget)
		}
		if request.AgentRunnerFacts != nil {
			t.Fatalf("step %s/%s must carry no AgentRunner facts off the zero-model path, got %+v", request.TaskID, request.Step.ID, *request.AgentRunnerFacts)
		}
	}
	for _, record := range store.events {
		if record.Event.Reason.Code == agentRunnerDispatchedReason {
			t.Fatalf("event %q records an AgentRunner dispatch on the zero-model path", record.Event.EventID)
		}
	}
}
