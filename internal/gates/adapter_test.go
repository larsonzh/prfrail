package gates

import (
	"context"
	"errors"
	"testing"

	"github.com/larsonzh/prfrail/internal/chain"
)

func TestChainPortReturnsDurableResultEvidence(t *testing.T) {
	hookRequest := validRequest(t.TempDir())
	hookRequest.Hook.Kind = "build"
	store := &memoryResults{}
	runner := testRunner(&fakeExecutor{executions: []Execution{successfulExecution()}}, fakeScanner{}, store)
	port := ChainPort{Runner: &runner, Hooks: map[string]Hook{"step-one": hookRequest.Hook}}
	result, err := port.RunHook(context.Background(), chain.StepRequest{RunID: "run-one", TaskID: "task-one", Step: chain.Step{ID: "step-one", Kind: "build"}, Attempt: 1, Workspace: chain.Workspace{Root: hookRequest.WorkspaceRoot}})
	if err != nil || len(result.Evidence) != 1 || result.Evidence[0] != store.records[0].ResultHash {
		t.Fatalf("chain evidence: %+v err=%v", result, err)
	}
}

func TestChainPortReturnsFailedResultEvidence(t *testing.T) {
	hookRequest := validRequest(t.TempDir())
	hookRequest.Hook.Kind = "verify"
	failed := successfulExecution()
	exitCode := 1
	failed.ExitCode = &exitCode
	store := &memoryResults{}
	runner := testRunner(&fakeExecutor{executions: []Execution{failed}}, fakeScanner{}, store)
	port := ChainPort{Runner: &runner, Hooks: map[string]Hook{"step-one": hookRequest.Hook}}
	result, err := port.RunHook(context.Background(), chain.StepRequest{RunID: "run-one", TaskID: "task-one", Step: chain.Step{ID: "step-one", Kind: "verify"}, Attempt: 1, Workspace: chain.Workspace{Root: hookRequest.WorkspaceRoot}})
	if !errors.Is(err, ErrGateFailed) || len(result.Evidence) != 1 || result.Evidence[0] != store.records[0].ResultHash {
		t.Fatalf("failed chain evidence: %+v err=%v", result, err)
	}
}

func TestChainPortRejectsMissingOrMismatchedHook(t *testing.T) {
	port := ChainPort{}
	if _, err := port.RunHook(context.Background(), chain.StepRequest{Step: chain.Step{ID: "missing", Kind: "build"}}); !errors.Is(err, ErrInvalidHook) {
		t.Fatalf("got %v", err)
	}
}

func TestChainPortDeniesExternalWriteHook(t *testing.T) {
	hookRequest := validRequest(t.TempDir())
	hookRequest.Hook.Kind = "build"
	hookRequest.Hook.EffectClass = "external-write"
	hookRequest.Hook.RecoveryGuarantee = "reconcile-only"
	hookRequest.Hook.AuthorizationHash = effectTestHash('e')
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	runner := testRunner(executor, fakeScanner{}, store)
	port := ChainPort{Runner: &runner, Hooks: map[string]Hook{"step-one": hookRequest.Hook}}
	result, err := port.RunHook(context.Background(), chain.StepRequest{RunID: "run-one", TaskID: "task-one", Step: chain.Step{ID: "step-one", Kind: "build"}, Attempt: 1, Workspace: chain.Workspace{Root: hookRequest.WorkspaceRoot}})
	if !errors.Is(err, ErrGateFailed) {
		t.Fatalf("external-write hook must fail the gate: %v", err)
	}
	if executor.calls != 0 {
		t.Fatalf("denied external-write hook must not execute: calls=%d", executor.calls)
	}
	if len(store.records) != 1 || len(result.Evidence) != 1 || result.Evidence[0] != store.records[0].ResultHash {
		t.Fatalf("denied hook must persist durable evidence: %+v", result)
	}
	persisted := store.records[0]
	if persisted.Result.EffectRecoveryAction == nil || *persisted.Result.EffectRecoveryAction != "reconcile-only" {
		t.Fatalf("denied hook result must demand reconciliation: %+v", persisted.Result)
	}
}
