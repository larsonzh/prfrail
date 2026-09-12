package chain

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type routeRecorder struct{ managed, isolated, handoff, hooks int }

func (recorder *routeRecorder) ApplyManagedChangeSet(context.Context, StepRequest) (StepResult, error) {
	recorder.managed++
	return StepResult{}, nil
}
func (recorder *routeRecorder) CaptureIsolatedWorkspace(context.Context, StepRequest) (StepResult, error) {
	recorder.isolated++
	return StepResult{}, nil
}
func (recorder *routeRecorder) OpenManualHandoff(context.Context, StepRequest) (StepResult, error) {
	recorder.handoff++
	return StepResult{}, nil
}
func (recorder *routeRecorder) RunHook(context.Context, StepRequest) (StepResult, error) {
	recorder.hooks++
	return StepResult{}, nil
}

type agentRunnerAdmissionRecorder struct {
	calls    int
	requests []StepRequest
	err      error
	sequence *[]string
}

func (recorder *agentRunnerAdmissionRecorder) AdmitAgentRunner(_ context.Context, request StepRequest) error {
	recorder.calls++
	recorder.requests = append(recorder.requests, request)
	if recorder.sequence != nil {
		*recorder.sequence = append(*recorder.sequence, "admit")
	}
	return recorder.err
}

type agentRunnerRecorder struct {
	calls    int
	requests []StepRequest
	result   StepResult
	err      error
	sequence *[]string
}

func (recorder *agentRunnerRecorder) DispatchAgentRunner(_ context.Context, request StepRequest) (StepResult, error) {
	recorder.calls++
	recorder.requests = append(recorder.requests, request)
	if recorder.sequence != nil {
		*recorder.sequence = append(*recorder.sequence, "dispatch")
	}
	if recorder.result.Evidence == nil {
		recorder.result = StepResult{Evidence: []string{"agent-runner-dispatched"}}
	}
	return recorder.result, recorder.err
}

func isolatedAgentRunnerRequest() StepRequest {
	return StepRequest{
		RunID:      "run-one",
		TaskID:     "task-one",
		Step:       Step{ID: "step-one", Kind: "code", Mode: IsolatedWorkspace},
		Attempt:    1,
		ParentHash: "sha256:1111111111111111111111111111111111111111111111111111111111111111",
		Workspace:  Workspace{Root: "run-one/task-one"},
	}
}

func TestStepRouterSeparatesChangeModesAndHooks(t *testing.T) {
	recorder := &routeRecorder{}
	router := StepRouter{Managed: recorder, Isolated: recorder, Handoff: recorder, Hooks: recorder}
	requests := []StepRequest{
		{Step: Step{Kind: "code", Mode: ManagedChangeSet}},
		{Step: Step{Kind: "code", Mode: IsolatedWorkspace}},
		{Step: Step{Kind: "code", Mode: ManualHandoff}},
		{Step: Step{Kind: "build"}},
		{Step: Step{Kind: "verify"}},
	}
	for _, request := range requests {
		if _, err := router.Execute(context.Background(), request); err != nil {
			t.Fatal(err)
		}
	}
	if recorder.managed != 1 || recorder.isolated != 1 || recorder.handoff != 1 || recorder.hooks != 2 {
		t.Fatalf("unexpected routes: %+v", recorder)
	}
	if _, err := router.Execute(context.Background(), StepRequest{Step: Step{Kind: "noop"}}); err == nil {
		t.Fatal("noop reached executor")
	}
}

func TestStepRouterExecuteAgentRunnerRejectsNilPort(t *testing.T) {
	admission := &agentRunnerAdmissionRecorder{}
	router := StepRouter{AgentRunnerAdmission: admission}

	if _, err := router.ExecuteAgentRunner(context.Background(), isolatedAgentRunnerRequest()); !errors.Is(err, ErrInvalidDefinition) {
		t.Fatalf("expected invalid definition for nil AgentRunner port, got %v", err)
	}
	if admission.calls != 0 {
		t.Fatalf("nil AgentRunner port invoked admission %d times", admission.calls)
	}
}

func TestStepRouterExecuteAgentRunnerRejectsOtherStepsBeforeAdmission(t *testing.T) {
	for _, step := range []Step{
		{Kind: "build"},
		{Kind: "noop"},
		{Kind: "code", Mode: ManagedChangeSet},
		{Kind: "code", Mode: ManualHandoff},
	} {
		t.Run(step.Kind+"/"+string(step.Mode), func(t *testing.T) {
			admission := &agentRunnerAdmissionRecorder{}
			port := &agentRunnerRecorder{}
			router := StepRouter{AgentRunner: port, AgentRunnerAdmission: admission}
			request := isolatedAgentRunnerRequest()
			request.Step = step

			if _, err := router.ExecuteAgentRunner(context.Background(), request); !errors.Is(err, ErrInvalidDefinition) {
				t.Fatalf("expected invalid definition, got %v", err)
			}
			if admission.calls != 0 || port.calls != 0 {
				t.Fatalf("invalid step reached AgentRunner: admission=%d dispatch=%d", admission.calls, port.calls)
			}
		})
	}
}

func TestStepRouterExecuteAgentRunnerRejectsNilAdmission(t *testing.T) {
	port := &agentRunnerRecorder{}
	router := StepRouter{AgentRunner: port}

	if _, err := router.ExecuteAgentRunner(context.Background(), isolatedAgentRunnerRequest()); !errors.Is(err, ErrInvalidDefinition) {
		t.Fatalf("expected invalid definition for nil AgentRunner admission, got %v", err)
	}
	if port.calls != 0 {
		t.Fatalf("nil AgentRunner admission dispatched %d times", port.calls)
	}
}

func TestStepRouterExecuteAgentRunnerBlocksAdmissionError(t *testing.T) {
	blocked := errors.New("admission blocked")
	admission := &agentRunnerAdmissionRecorder{err: blocked}
	port := &agentRunnerRecorder{}
	router := StepRouter{AgentRunner: port, AgentRunnerAdmission: admission}

	if _, err := router.ExecuteAgentRunner(context.Background(), isolatedAgentRunnerRequest()); !errors.Is(err, blocked) {
		t.Fatalf("expected admission error, got %v", err)
	}
	if port.calls != 0 {
		t.Fatalf("blocked AgentRunner admission dispatched %d times", port.calls)
	}
}

func TestStepRouterExecuteAgentRunnerDispatchesAfterAdmission(t *testing.T) {
	sequence := make([]string, 0, 2)
	admission := &agentRunnerAdmissionRecorder{sequence: &sequence}
	port := &agentRunnerRecorder{sequence: &sequence}
	router := StepRouter{AgentRunner: port, AgentRunnerAdmission: admission}

	result, err := router.ExecuteAgentRunner(context.Background(), isolatedAgentRunnerRequest())
	if err != nil {
		t.Fatal(err)
	}
	if port.calls != 1 {
		t.Fatalf("expected one AgentRunner dispatch, got %d", port.calls)
	}
	if !reflect.DeepEqual(sequence, []string{"admit", "dispatch"}) {
		t.Fatalf("unexpected AgentRunner call order: %v", sequence)
	}
	if len(result.Evidence) != 1 || result.Evidence[0] != "agent-runner-dispatched" {
		t.Fatalf("unexpected AgentRunner result: %+v", result)
	}
}

func TestStepRouterExecuteAgentRunnerReturnsDispatchError(t *testing.T) {
	dispatchErr := errors.New("dispatch failed")
	port := &agentRunnerRecorder{err: dispatchErr}
	router := StepRouter{AgentRunner: port, AgentRunnerAdmission: &agentRunnerAdmissionRecorder{}}

	if _, err := router.ExecuteAgentRunner(context.Background(), isolatedAgentRunnerRequest()); !errors.Is(err, dispatchErr) {
		t.Fatalf("expected dispatch error, got %v", err)
	}
	if port.calls != 1 {
		t.Fatalf("expected one dispatch attempt, got %d", port.calls)
	}
}

func TestStepRouterExecuteAgentRunnerPassesSameRequest(t *testing.T) {
	request := isolatedAgentRunnerRequest()
	want := request
	admission := &agentRunnerAdmissionRecorder{}
	port := &agentRunnerRecorder{}
	router := StepRouter{AgentRunner: port, AgentRunnerAdmission: admission}

	if _, err := router.ExecuteAgentRunner(context.Background(), request); err != nil {
		t.Fatal(err)
	}
	if len(admission.requests) != 1 || len(port.requests) != 1 {
		t.Fatalf("unexpected request counts: admission=%d port=%d", len(admission.requests), len(port.requests))
	}
	if !reflect.DeepEqual(admission.requests[0], want) || !reflect.DeepEqual(port.requests[0], want) || !reflect.DeepEqual(admission.requests[0], port.requests[0]) {
		t.Fatalf("admission and port did not receive the same request: admission=%+v port=%+v want=%+v", admission.requests[0], port.requests[0], want)
	}
}
