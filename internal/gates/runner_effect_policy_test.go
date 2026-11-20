package gates

// P83 F5 regressions: the effect class of a run is derived from the frozen hook
// policy, never from an optional declared scope.
//
// Three properties are pinned here, and each one is the mutation contract of the
// corresponding branch in deriveEffectClassification:
//
//   - TestRunnerDerivesClassificationFromFrozenHookPolicy: request.Effect == nil
//     still executes AND produces an effect observation. Restoring the
//     "classification != nil" short-circuit in Run drops the observation and turns
//     this red.
//   - TestRunnerRefusesContradictoryEffectScope: a scope class that contradicts
//     the hook never reaches Execute. Deleting the contradiction check calls the
//     executor and turns this red.
//   - TestRunnerRefusesInvalidScopePolicyHashWithoutPlaceholder: an invalid scope
//     policy hash is refused instead of being replaced by a fabricated placeholder
//     digest. Restoring the placeholder turns this red.

import (
	"context"
	"errors"
	"testing"
	"time"
)

// requestClock is the fixed instant the classification tests classify at.
func requestClock() time.Time { return time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC) }

func TestRunnerDerivesClassificationFromFrozenHookPolicy(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Hook.EffectClass = "read-only"
	request.Effect = nil
	refreshHookDefinition(t, &request)
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	record, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if executor.calls != 1 {
		t.Fatalf("a read-only hook without a declared scope must still execute, calls=%d", executor.calls)
	}
	if record.Result.Assessment != "passed" {
		t.Fatalf("expected a passed result, got %+v", record.Result)
	}
	if record.Result.EffectObservationHash == nil || record.Result.EffectRecoveryAction == nil {
		t.Fatalf("an unscoped read-only run must still be classified and observed: %+v", record.Result)
	}
	if *record.Result.EffectRecoveryAction != "none-required" {
		t.Fatalf("a completed read-only effect needs no recovery, got %q", *record.Result.EffectRecoveryAction)
	}
}

func TestRunnerRefusesContradictoryEffectScope(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Hook.EffectClass = "read-only"
	request.Effect = effectScopeForTest("external-write")
	refreshHookDefinition(t, &request)
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	_, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if !errors.Is(err, ErrEffectClassificationConflict) {
		t.Fatalf("a scope that contradicts the frozen hook must be refused, got %v", err)
	}
	if executor.calls != 0 {
		t.Fatalf("a contradictory classification must never execute, calls=%d", executor.calls)
	}
	if len(store.records) != 0 {
		t.Fatalf("a contradictory classification must persist nothing, got %+v", store.records)
	}
}

func TestRunnerRefusesInvalidScopePolicyHashWithoutPlaceholder(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Hook.EffectClass = "read-only"
	request.Effect = effectScopeForTest("read-only")
	request.Effect.PolicyHash = "not-a-hash"
	refreshHookDefinition(t, &request)
	// The classification layer is the guard that used to invent a placeholder, so it
	// is exercised directly: Prepare already refuses the scope shape earlier, which
	// would hide the placeholder branch.
	if classification, err := deriveEffectClassification(request, requestClock()); !errors.Is(err, ErrEffectClassificationUnavailable) || classification != nil {
		t.Fatalf("an invalid scope policy hash must be refused with no classification, got %+v, %v", classification, err)
	}
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	_, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if err == nil {
		t.Fatal("an invalid scope policy hash must refuse the run")
	}
	if executor.calls != 0 {
		t.Fatalf("an unclassifiable request must never execute, calls=%d", executor.calls)
	}
	if len(store.records) != 0 {
		t.Fatalf("an unclassifiable request must persist nothing, got %+v", store.records)
	}
}
