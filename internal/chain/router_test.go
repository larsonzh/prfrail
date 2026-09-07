package chain

import (
	"context"
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
