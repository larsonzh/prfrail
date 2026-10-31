package tickets

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// Test-only contract for the global lock order slice
// (design: docs/validation/evidence/P40-T027B-LOCK-REAPER/01-arch-design/review.md §2/§3).
//
// These tests intentionally reference an API that does not exist yet:
// OpenDurableCostLedger(ctx, config) and ErrDurableLockReentrant. The expected
// RED state of this slice is a compile failure.
const (
	p40LockHoldEnv   = "PRFRAIL_P40_LOCK_HOLD_ROOT"
	p40LockReadyEnv  = "PRFRAIL_P40_LOCK_HOLD_READY"
	p40LockHelperRun = "^TestDurableLockHelperHold$"
)

func p40DurableConfig(root, ledgerID, scopeHash, entryID string) DurableCostLedgerConfig {
	amount := int64(10_000_000)
	return DurableCostLedgerConfig{
		Root: root,
		Ledger: CostLedgerConfig{
			LedgerID:       ledgerID,
			CreatedAt:      costNow(),
			ScopeKind:      CostScopeProject,
			ScopeHash:      scopeHash,
			Currency:       "USD",
			PricingMode:    CostPricingDocumented,
			PricingVersion: "2026-09-08",
		},
		Initial: AllocationInput{
			EntryID:           entryID,
			OccurredAt:        costNow(),
			AuthorizationHash: costHashTwo,
			Limits: CostLimits{
				AmountMicros: &amount,
				ModelCalls:   10,
				Tokens:       1000,
				WallClockMs:  60000,
				Attempts:     3,
			},
		},
	}
}

func p40Reservation(entryID, requestID, key string) ReservationInput {
	amount := int64(1000)
	return ReservationInput{
		EntryID:              entryID,
		OccurredAt:           costNow(),
		RunID:                "run-lock-order",
		RequestID:            requestID,
		IdempotencyKey:       key,
		AuthorizationHash:    costHashThree,
		ReservedAmountMicros: &amount,
		ReservedCalls:        1,
		ReservedTokens:       10,
	}
}

func p40Open(t *testing.T, root, ledgerID, scopeHash, entryID string) *DurableCostLedger {
	t.Helper()
	ledger, err := OpenDurableCostLedger(context.Background(), p40DurableConfig(root, ledgerID, scopeHash, entryID))
	if err != nil {
		t.Fatalf("open durable ledger %s: %v", ledgerID, err)
	}
	t.Cleanup(func() { _ = ledger.Close() })
	return ledger
}

// TestLockOrderNoCrossDeadlock pins the global order T0 < T1 < T2: two paths
// taking the same two scope locks in opposite order must not both block. At
// least one path must converge, either by completing or by the reentrancy
// guard returning ErrDurableLockReentrant.
func TestLockOrderNoCrossDeadlock(t *testing.T) {
	root := t.TempDir()
	ledgerA := p40Open(t, root, "durable-cost-ledger-a", costHashOne, "allocation-a")
	ledgerB := p40Open(t, root, "durable-cost-ledger-b", costHashThree, "allocation-b")

	ctx := context.Background()
	if _, _, err := ledgerA.Reserve(ctx, p40Reservation("reserve-a", "request-a", "reserve-a")); err != nil {
		t.Fatalf("reserve on scope a: %v", err)
	}
	if _, _, err := ledgerB.Reserve(ctx, p40Reservation("reserve-b", "request-b", "reserve-b")); err != nil {
		t.Fatalf("reserve on scope b: %v", err)
	}

	// Nested acquisition mirrors "hold budget lock, then need the next lock in
	// the chain". Opposite orders across goroutines is the deadlock shape.
	nested := func(first, second *DurableCostLedger) error {
		return first.runWithLock(ctx, true, func() error {
			return second.runWithLock(ctx, true, func() error { return nil })
		})
	}

	results := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); results <- nested(ledgerA, ledgerB) }()
	go func() { defer wg.Done(); results <- nested(ledgerB, ledgerA) }()

	converged := make(chan struct{})
	go func() { wg.Wait(); close(converged) }()
	select {
	case <-converged:
	case <-time.After(5 * time.Second):
		t.Fatal("opposite lock orders both hung: cross deadlock (lock order is not monotonic)")
	}
	close(results)

	completed := 0
	for err := range results {
		switch {
		case err == nil:
			completed++
		case errors.Is(err, ErrDurableLockReentrant),
			errors.Is(err, context.DeadlineExceeded),
			errors.Is(err, context.Canceled):
			// A refusal that names the order violation is the allowed outcome.
		default:
			t.Fatalf("unexpected lock-order error: %v", err)
		}
	}
	if completed == 0 {
		t.Fatal("no path converged: neither completed nor rejected with ErrDurableLockReentrant")
	}
}

// TestDurableLockHelperHold is the child half of TestLockWaitHonoursCancel. It
// opens the ledger, holds the scope lock and signals readiness on disk.
func TestDurableLockHelperHold(t *testing.T) {
	root := os.Getenv(p40LockHoldEnv)
	if root == "" {
		t.Skip("lock-hold helper only runs as a child process")
	}
	ready := os.Getenv(p40LockReadyEnv)
	ledger, err := OpenDurableCostLedger(context.Background(),
		p40DurableConfig(root, "durable-cost-ledger-hold", costHashOne, "allocation-hold"))
	if err != nil {
		t.Fatalf("helper open ledger: %v", err)
	}
	defer func() { _ = ledger.Close() }()
	if err := ledger.runWithLock(context.Background(), true, func() error {
		if writeErr := os.WriteFile(ready, []byte("holding"), 0o600); writeErr != nil {
			return writeErr
		}
		time.Sleep(30 * time.Second)
		return nil
	}); err != nil {
		t.Fatalf("helper hold scope lock: %v", err)
	}
}

// TestLockWaitHonoursCancel pins that a wait for a scope lock held by another
// process stays interruptible: with the lock stuck in a child process, a
// cancelled parent context must make the open return ctx.Err() well inside the
// bounded window instead of blocking forever.
func TestLockWaitHonoursCancel(t *testing.T) {
	root := t.TempDir()
	ready := filepath.Join(root, "lock-held.marker")

	cmd := exec.Command(os.Args[0], "-test.run="+p40LockHelperRun)
	cmd.Env = append(os.Environ(), p40LockHoldEnv+"="+root, p40LockReadyEnv+"="+ready)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start lock-hold helper: %v", err)
	}
	killed := false
	kill := func() {
		if killed {
			return
		}
		killed = true
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	}
	t.Cleanup(kill)

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if _, err := os.Stat(ready); err != nil {
		data, _ := os.ReadFile(ready)
		t.Fatalf("helper never acquired the scope lock (ready=%q stat err=%v)", strings.TrimSpace(string(data)), err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	type openResult struct {
		ledger *DurableCostLedger
		err    error
	}
	opened := make(chan openResult, 1)
	go func() {
		ledger, err := OpenDurableCostLedger(ctx,
			p40DurableConfig(root, "durable-cost-ledger-hold", costHashOne, "allocation-hold"))
		opened <- openResult{ledger: ledger, err: err}
	}()

	select {
	case result := <-opened:
		if result.ledger != nil {
			_ = result.ledger.Close()
		}
		if result.err == nil {
			t.Fatal("open must not succeed while another process holds the scope lock")
		}
		if !errors.Is(result.err, context.DeadlineExceeded) && !errors.Is(result.err, context.Canceled) {
			t.Fatalf("lock wait returned %v, want the context error", result.err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("lock wait ignored context cancellation: unbounded wait on a stuck scope lock")
	}
}

// TestCancelWinsOverTimeout pins the convergence rule: when cancellation and
// the deadline land together, cancellation wins. A pre-cancelled context with an
// already expired deadline must surface ctx.Err() as context.Canceled.
func TestCancelWinsOverTimeout(t *testing.T) {
	root := t.TempDir()
	cancelled, cancel := context.WithCancel(context.Background())
	cancel() // cancellation happened first and must win
	timedCtx, deadlineCancel := context.WithTimeout(cancelled, 20*time.Millisecond)
	defer deadlineCancel()

	ledger, err := OpenDurableCostLedger(timedCtx,
		p40DurableConfig(root, "durable-cost-ledger-cancel", costHashTwo, "allocation-cancel"))
	if ledger != nil {
		_ = ledger.Close()
		t.Fatal("open must not succeed with a cancelled context")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("open returned %v, want context.Canceled", err)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		t.Fatal("timeout must not win over an already-cancelled context")
	}
}
