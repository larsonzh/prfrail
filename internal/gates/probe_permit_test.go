package gates

// P44-T027B-GATE-READINESS (N4) test-first (RED): the probe permit is a one-shot
// quota. Require reserves atomically, Complete spends it, and neither a failure, a
// timeout nor a process restart may hand the quota back.
//
// This file must fail to compile until internal/gates/probe_permit.go provides
// every symbol pinned below, with exactly the names and signatures the ①
// architecture design fixed. Compile failure - not assertion failure - is the
// expected RED shape for this slice.
//
// Falsification rule applied to every case below: remove the corresponding guard
// from the implementation and the case must turn red.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// Compile-time pins: the published surface of the N4 design.
var (
	_ func(string, string, int64) *ProbePermit = NewProbePermit
	_ func(string) (*ProbePermit, error)       = LoadProbePermit

	_ func(*ProbePermit, context.Context) error = (*ProbePermit).Require
	_ func(*ProbePermit) error                  = (*ProbePermit).Complete
)

// TestProbePermitConcurrentRequireSingleWinner closes the double-spend race: with
// 40 concurrent callers on one permit, exactly one may win.
func TestProbePermitConcurrentRequireSingleWinner(t *testing.T) {
	permit := NewProbePermit("permit-concurrent", filepath.Join(t.TempDir(), "probe-permit.json"), 1)

	const contenders = 40
	var (
		wait     sync.WaitGroup
		start    = make(chan struct{})
		outcomes = make([]error, contenders)
	)

	for index := 0; index < contenders; index++ {
		wait.Add(1)
		go func(slot int) {
			defer wait.Done()
			<-start
			outcomes[slot] = permit.Require(context.Background())
		}(index)
	}
	close(start)
	wait.Wait()

	winners := 0
	for index, err := range outcomes {
		switch {
		case err == nil:
			winners++
		case errors.Is(err, ErrProbePermitInUse), errors.Is(err, ErrProbePermitConsumed):
		default:
			t.Fatalf("contender %d got an unauthorised error: %v", index, err)
		}
	}
	if winners != 1 {
		t.Fatalf("exactly one contender may claim the permit, got %d winners out of %d", winners, contenders)
	}
}

// TestProbePermitFailureDoesNotResetQuota: an abandoned attempt keeps the claim,
// and a completed attempt keeps the permit spent.
func TestProbePermitFailureDoesNotResetQuota(t *testing.T) {
	permit := NewProbePermit("permit-failure", filepath.Join(t.TempDir(), "probe-permit.json"), 1)

	if err := permit.Require(context.Background()); err != nil {
		t.Fatalf("the first require must succeed: %v", err)
	}

	// The probe failed: the caller abandons the attempt without Complete.
	if err := permit.Require(context.Background()); !errors.Is(err, ErrProbePermitInUse) {
		t.Fatalf("an in-flight claim must survive a probe failure: %v", err)
	}

	if err := permit.Complete(); err != nil {
		t.Fatalf("complete must succeed on a claimed permit: %v", err)
	}
	if err := permit.Require(context.Background()); !errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("a spent permit must stay spent: %v", err)
	}
}

// TestProbePermitTimeoutDoesNotResetQuota: cancelling the caller context (a probe
// timeout) must not return the quota to the pool.
func TestProbePermitTimeoutDoesNotResetQuota(t *testing.T) {
	root := t.TempDir()

	claimed := NewProbePermit("permit-timeout-claimed", filepath.Join(root, "claimed.json"), 1)
	ctx, cancel := context.WithCancel(context.Background())
	if err := claimed.Require(ctx); err != nil {
		t.Fatalf("the first require must succeed: %v", err)
	}
	cancel() // the probe timed out; the caller abandons the attempt

	err := claimed.Require(context.Background())
	if err == nil {
		t.Fatal("a timed-out attempt must not return the quota to the pool")
	}
	if !errors.Is(err, ErrProbePermitInUse) && !errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("a timed-out attempt must still hold the claim: %v", err)
	}

	cancelled := NewProbePermit("permit-timeout-cancelled", filepath.Join(root, "cancelled.json"), 1)
	cancelledCtx, cancelBefore := context.WithCancel(context.Background())
	cancelBefore()
	if err := cancelled.Require(cancelledCtx); err == nil {
		t.Fatal("a require on an already cancelled context must fail closed")
	}
}

// TestProbePermitRestartDoesNotResetQuota: the claim is durable on disk, so a
// rebuilt permit (process restart) still refuses a second claim.
func TestProbePermitRestartDoesNotResetQuota(t *testing.T) {
	path := filepath.Join(t.TempDir(), "probe-permit.json")
	permit := NewProbePermit("permit-restart", path, 1)

	if err := permit.Require(context.Background()); err != nil {
		t.Fatalf("the first require must succeed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("the claim must be durable on disk before a restart: %v", err)
	}

	restarted, err := LoadProbePermit(path)
	if err != nil {
		t.Fatalf("load the persisted permit: %v", err)
	}
	if err := restarted.Require(context.Background()); !errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("a restarted permit must not hand out a second claim: %v", err)
	}
}

// TestProbePermitDamagedMarkerStaysSpent pins the fail-closed reading of the
// marker: its existence is the claim bit, so an unparseable body must still load
// as a spent permit. A damaged marker that demoted the permit back to "unused"
// would hand the same quota out a second time.
func TestProbePermitDamagedMarkerStaysSpent(t *testing.T) {
	root := t.TempDir()

	// Control: with no marker at all the loader yields a claimable permit, so the
	// assertion below discriminates "marker present" from "load always spent".
	controlPath := filepath.Join(root, "control.json")
	control, err := LoadProbePermit(controlPath)
	if err != nil {
		t.Fatalf("a missing marker must load as a fresh quota: %v", err)
	}
	if control == nil {
		t.Fatal("a missing marker must still return a permit bound to the path")
	}
	if err := control.Require(context.Background()); err != nil {
		t.Fatalf("a missing marker must leave the quota claimable: %v", err)
	}

	damagedPath := filepath.Join(root, "damaged.json")
	if err := os.WriteFile(damagedPath, []byte("{"), 0o644); err != nil {
		t.Fatalf("write the damaged marker: %v", err)
	}

	damaged, err := LoadProbePermit(damagedPath)
	if err != nil {
		// A hard-loading implementation is acceptable, but it must not hand back
		// a permit alongside the error: that would be an unauthorised quota.
		if damaged != nil {
			t.Fatalf("a failed load must not hand back a permit: %v", err)
		}
		return
	}
	if damaged == nil {
		t.Fatal("a successful load must return a permit")
	}
	if err := damaged.Require(context.Background()); !errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("a damaged marker must stay spent, require got: %v", err)
	}
}

// TestProbePermitSecondInstanceOnSamePathLosesOnly models two processes that each
// build their own permit (their own in-memory state) over the same marker path:
// the durable O_CREATE|O_EXCL claim must admit exactly one of them.
func TestProbePermitSecondInstanceOnSamePathLosesOnly(t *testing.T) {
	path := filepath.Join(t.TempDir(), "probe-permit.json")

	first := NewProbePermit("permit-cross-first", path, 1)
	second := NewProbePermit("permit-cross-second", path, 1)

	outcomes := []error{
		first.Require(context.Background()),
		second.Require(context.Background()),
	}

	winners := 0
	for index, err := range outcomes {
		switch {
		case err == nil:
			winners++
		case errors.Is(err, ErrProbePermitConsumed), errors.Is(err, ErrProbePermitInUse):
		default:
			t.Fatalf("instance %d got an unauthorised error: %v", index, err)
		}
	}
	if winners != 1 {
		t.Fatalf("the durable claim must admit exactly one of two instances, got %d winners", winners)
	}
}

// TestProbePermitCompleteOnUnusedRejected: completion cannot invent a claim, and
// the refusal must not spend the untouched quota.
func TestProbePermitCompleteOnUnusedRejected(t *testing.T) {
	permit := NewProbePermit("permit-complete-unused", filepath.Join(t.TempDir(), "probe-permit.json"), 1)

	if err := permit.Complete(); !errors.Is(err, ErrProbePermitInUse) {
		t.Fatalf("completing an unused permit must be refused as in use, got: %v", err)
	}
	if err := permit.Require(context.Background()); err != nil {
		t.Fatalf("a refused complete must not burn the quota, require got: %v", err)
	}
	if err := permit.Complete(); err != nil {
		t.Fatalf("complete after the claim must succeed: %v", err)
	}
}

// TestProbePermitRestartThenCompleteIsIdempotent: a reloaded (already spent)
// permit completes without error and still refuses a claim.
func TestProbePermitRestartThenCompleteIsIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "probe-permit.json")

	permit := NewProbePermit("permit-restart-complete", path, 1)
	if err := permit.Require(context.Background()); err != nil {
		t.Fatalf("the first require must succeed: %v", err)
	}
	if err := permit.Complete(); err != nil {
		t.Fatalf("complete must succeed on a claimed permit: %v", err)
	}

	reloaded, err := LoadProbePermit(path)
	if err != nil {
		t.Fatalf("load the persisted permit: %v", err)
	}
	if err := reloaded.Complete(); err != nil {
		t.Fatalf("complete on an already spent permit must be idempotent, got: %v", err)
	}
	if err := reloaded.Require(context.Background()); !errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("an idempotent complete must not restore the quota, require got: %v", err)
	}
}

// TestProbePermitSuccessImpliesMarkerOnDisk pins the persistence-first invariant
// of Require: a nil return means the durable claim marker was already written and
// synced, so a process that dies right after the admission can never re-issue the
// same quota.
//
// Falsification: move persistClaim after the in-memory state flip and return nil
// on the cache alone, and the os.Stat below turns red.
func TestProbePermitSuccessImpliesMarkerOnDisk(t *testing.T) {
	path := filepath.Join(t.TempDir(), "probe-permit.json")
	permit := NewProbePermit("permit-persist-first", path, 1)

	if err := permit.Require(context.Background()); err != nil {
		t.Fatalf("the first require must succeed: %v", err)
	}

	// Assertion order matters: nil first, then the disk, then the reloaded view.
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("a successful Require must have persisted the marker before returning: %v", err)
	}

	reloaded, err := LoadProbePermit(path)
	if err != nil {
		t.Fatalf("load the persisted permit: %v", err)
	}
	if reloaded == nil {
		t.Fatal("loading a persisted permit must return a permit bound to the path")
	}
	if err := reloaded.Require(context.Background()); !errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("a permit reloaded from a persisted claim must be spent, require got: %v", err)
	}
}

// TestProbePermitCancelledBeforeClaimDoesNotConsume pins the other half of the
// persistence-first invariant: a cancellation observed before the claim touches
// the disk must not burn the quota, because no probe ever ran.
//
// Falsification: drop the early ctx.Err() pre-check, or persist before checking
// the context, and the second require below turns red (the quota would already
// be claimed or spent).
func TestProbePermitCancelledBeforeClaimDoesNotConsume(t *testing.T) {
	path := filepath.Join(t.TempDir(), "probe-permit.json")
	permit := NewProbePermit("permit-cancel-before-claim", path, 1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // the caller gave up before the claim was ever attempted

	err := permit.Require(ctx)
	if err == nil {
		t.Fatal("a require on an already cancelled context must fail closed")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("a pre-claim cancellation must be reported as context.Canceled, got: %v", err)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a pre-claim cancellation must not persist a claim marker, stat err = %v", statErr)
	}

	// The quota was never taken, so a healthy context must still be admitted.
	if err := permit.Require(context.Background()); err != nil {
		t.Fatalf("a cancellation before the claim must not burn the quota, require got: %v", err)
	}
	if _, statErr := os.Stat(path); statErr != nil {
		t.Fatalf("the admitted require must have created the claim marker: %v", statErr)
	}
}
