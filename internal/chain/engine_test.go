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
	hashes []string
	calls  int
}

func (fake *fakeAcceptance) Accept(context.Context, string, string, int, SnapshotRef, Workspace) (SnapshotRef, []string, error) {
	hash := fake.hashes[fake.calls]
	fake.calls++
	return SnapshotRef{Hash: hash}, []string{stopHash}, nil
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

func testOptions(store EventStore) (Options, *fakeBaseline, *fakeWorkspaces, *fakeSteps, *fakeAcceptance, *fakeStopper, *fakeReconciler) {
	baseline := &fakeBaseline{}
	workspaces := &fakeWorkspaces{}
	steps := &fakeSteps{}
	acceptance := &fakeAcceptance{hashes: []string{acceptedOne, acceptedTwo, acceptedThree}}
	stopper := &fakeStopper{}
	reconciler := &fakeReconciler{}
	sequence := 0
	options := Options{RunID: "run-one", Definition: threeTaskDefinition(), Events: store, Baselines: baseline, Workspaces: workspaces, Steps: steps, Acceptance: acceptance, Stopper: stopper, Reconciler: reconciler, Clock: fixedClock, IDs: func() string { sequence++; return fmt.Sprintf("event-%d", sequence) }}
	return options, baseline, workspaces, steps, acceptance, stopper, reconciler
}

func TestEngineRunsThreeTasksFourKindsAndBothModes(t *testing.T) {
	store := &memoryEvents{}
	options, baseline, workspaces, steps, acceptance, _, _ := testOptions(store)
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if engine.Projection().ChainState != "COMPLETED" || baseline.calls != 1 || acceptance.calls != 3 {
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

func TestTaskFailureStopsLaterTasks(t *testing.T) {
	store := &memoryEvents{}
	options, _, workspaces, steps, _, _, _ := testOptions(store)
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
	options, _, _, steps, _, _, _ := testOptions(store)
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
	options, _, _, _, _, stopper, _ := testOptions(store)
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
	options, _, _, steps, _, _, reconciler := testOptions(store)
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
	recoveryOptions, _, _, recoverySteps, _, _, recoveryReconciler := testOptions(store)
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
