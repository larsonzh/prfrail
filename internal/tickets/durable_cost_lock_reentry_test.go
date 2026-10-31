package tickets

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
)

// TestSameLedgerNestedAcquisitionRefused pins the reentrancy rule of the global
// lock order (docs/validation/evidence/P40-T027B-LOCK-REAPER/01-arch-design/review.md
// §2 and §4 risk 2): a goroutine that already holds a scope lock may not take the
// same scope again, because the second, independent file lock would deadlock
// against the first.
func TestSameLedgerNestedAcquisitionRefused(t *testing.T) {
	ledger := p40Open(t, t.TempDir(), "durable-cost-ledger-reentry", costHashOne, "allocation-reentry")

	err := ledger.runWithLock(context.Background(), true, func() error {
		return ledger.runWithLock(context.Background(), true, func() error { return nil })
	})
	if !errors.Is(err, ErrDurableLockReentrant) {
		t.Fatalf("nested acquisition of the same scope = %v, want ErrDurableLockReentrant", err)
	}

	// The refused inner call must not have released the outer acquisition.
	if err := ledger.runWithLock(context.Background(), true, func() error { return nil }); err != nil {
		t.Fatalf("re-acquire after the refusal: %v", err)
	}
}

// TestPanicInLockedFnPopsRegistry pins that the registry entry is dropped even
// when the locked function panics: the release defer has to run during the
// unwind, otherwise the goroutine would keep a stale acquisition recorded and
// every later acquisition of that path would be refused as reentrant.
func TestPanicInLockedFnPopsRegistry(t *testing.T) {
	ledger := p40Open(t, t.TempDir(), "durable-cost-ledger-panic", costHashTwo, "allocation-panic")

	func() {
		defer func() {
			if recovered := recover(); recovered == nil {
				t.Fatal("the locked function panic did not propagate")
			}
		}()
		_ = ledger.runWithLock(context.Background(), true, func() error {
			panic("locked function panic")
		})
	}()

	if err := ledger.runWithLock(context.Background(), true, func() error { return nil }); err != nil {
		t.Fatalf("re-acquire after a panic inside the locked function: %v", err)
	}
}

// TestConcurrentSameScopeNonNestedSerialized pins that two goroutines using the
// same ledger without nesting are serialised rather than refused: both
// reservations must be durably committed and the canonical ledger must hold
// exactly one allocation and both reservation entries.
func TestConcurrentSameScopeNonNestedSerialized(t *testing.T) {
	ledger := p40Open(t, t.TempDir(), "durable-cost-ledger-concurrent", costHashThree, "allocation-concurrent")

	const reservations = 2
	results := make(chan error, reservations)
	var wait sync.WaitGroup
	for index := 0; index < reservations; index++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			_, _, err := ledger.Reserve(context.Background(), p40Reservation(
				fmt.Sprintf("reserve-%d", index),
				fmt.Sprintf("request-%d", index),
				fmt.Sprintf("reserve-key-%d", index),
			))
			results <- err
		}(index)
	}
	wait.Wait()
	close(results)
	for err := range results {
		if err != nil {
			t.Fatalf("concurrent same-scope reservation failed: %v", err)
		}
	}

	reloaded := reloadDurableLedger(t, ledger.Path())
	counts := durableEntryCounts(t, reloaded)
	if counts["allocation"] != 1 || counts["reservation"] != reservations {
		t.Fatalf("canonical ledger entries = %+v, want 1 allocation and %d reservations", counts, reservations)
	}
}
