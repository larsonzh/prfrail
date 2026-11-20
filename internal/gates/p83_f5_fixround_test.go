package gates

// P83 F5 fix round: the declared scope policy identity is a cross-check of the
// frozen hook policy, never an alternative truth, and the adapter never invents
// an identity the hook did not declare.
//
// Mutation contracts:
//   - TestRunnerRefusesConflictingScopePolicyHash: deleting the mismatch branch
//     in deriveEffectClassification lets the hook execute and turns this red.
//   - TestRunnerAcceptsMatchingScopePolicyHash: treating every declared scope
//     policy hash as a conflict turns this red.
//   - TestAdapterEffectScopeInventsNoPlaceholderIdentity: restoring the fabricated
//     placeholder digests in buildEffectScope turns this red.

import (
	"context"
	"errors"
	"testing"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// TestRunnerRefusesConflictingScopePolicyHash pins the conflict branch: a scope
// that names a different (but well-formed) policy than the frozen hook is a
// contradiction, so nothing executes and nothing is persisted.
func TestRunnerRefusesConflictingScopePolicyHash(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Hook.EffectClass = "read-only"
	request.Hook.PolicyHash = effectTestHash('a')
	refreshHookDefinition(t, &request)
	request.Effect = effectScopeForTest("read-only")
	request.Effect.PolicyHash = effectTestHash('b')
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	_, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if !errors.Is(err, ErrEffectClassificationConflict) {
		t.Fatalf("a scope policy hash that contradicts the frozen hook must be refused as a conflict, got %v", err)
	}
	if executor.calls != 0 {
		t.Fatalf("a conflicting policy identity must never execute, calls=%d", executor.calls)
	}
	if len(store.records) != 0 {
		t.Fatalf("a conflicting policy identity must persist nothing, got %+v", store.records)
	}
}

// TestRunnerAcceptsMatchingScopePolicyHash is the acceptance half: an identical
// policy hash is not a conflict, the classification binds the hook's policy and
// the run executes.
func TestRunnerAcceptsMatchingScopePolicyHash(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Hook.EffectClass = "read-only"
	request.Hook.PolicyHash = effectTestHash('a')
	refreshHookDefinition(t, &request)
	request.Effect = effectScopeForTest("read-only")
	request.Effect.PolicyHash = request.Hook.PolicyHash
	classification, err := deriveEffectClassification(request, requestClock())
	if err != nil {
		t.Fatalf("an equal policy hash must classify, got %v", err)
	}
	if classification.Record.PolicyHash != request.Hook.PolicyHash {
		t.Fatalf("the classification must bind the frozen hook policy %s, got %s", request.Hook.PolicyHash, classification.Record.PolicyHash)
	}
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	record, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if err != nil {
		t.Fatalf("an equal policy hash must not refuse the run, got %v", err)
	}
	if executor.calls != 1 || record.Result.Assessment != "passed" {
		t.Fatalf("an equal policy hash must execute once and pass, calls=%d record=%+v", executor.calls, record.Result)
	}
}

// TestAdapterEffectScopeInventsNoPlaceholderIdentity pins buildEffectScope: the
// scope restates the hook's declared identities verbatim, and a hook that
// declares no authorization identity yields no scope at all instead of a
// fabricated one.
func TestAdapterEffectScopeInventsNoPlaceholderIdentity(t *testing.T) {
	request := chain.StepRequest{RunID: "run-one", TaskID: "task-one", Step: chain.Step{ID: "step-one"}, Attempt: 1}
	undeclared := Hook{ID: "go-test", Kind: "test", EffectClass: "read-only"}
	if scope := buildEffectScope(request, undeclared); scope != nil {
		t.Fatalf("a hook without a declared authorization hash must yield no scope, got %+v", scope)
	}
	declared := Hook{
		ID:                "go-test",
		Kind:              "test",
		EffectClass:       "read-only",
		PolicyHash:        effectTestHash('a'),
		AuthorizationHash: effectTestHash('e'),
	}
	scope := buildEffectScope(request, declared)
	if scope == nil {
		t.Fatal("a hook with a declared identity must yield a scope")
	}
	if scope.PolicyHash != declared.PolicyHash || scope.AuthorizationHash != declared.AuthorizationHash {
		t.Fatalf("the scope must restate the declared identities, got policy=%s authorization=%s", scope.PolicyHash, scope.AuthorizationHash)
	}
	if len(scope.Evidence) != 1 || !evidence.ValidHash(scope.Evidence[0]) {
		t.Fatalf("the scope evidence must be a real hash set, got %v", scope.Evidence)
	}
}
