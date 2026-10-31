package guard

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

// reclaimHelperEnv marks the child half of the reclaim and identity-stop tests.
const reclaimHelperEnv = "PROOFRAIL_RECLAIM_HELPER"

// TestReclaimHelperStayAlive is the child half of the reclaim tests. It blocks on
// a timer rather than on an empty select, so the runtime cannot declare a
// deadlock: the helper really keeps running until the parent stops it.
func TestReclaimHelperStayAlive(t *testing.T) {
	if os.Getenv(reclaimHelperEnv) == "" {
		t.Skip("reclaim helper only runs as a child process")
	}
	for {
		time.Sleep(time.Hour)
	}
}

func reclaimHelperSpec() ProcessSpec {
	return ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=^TestReclaimHelperStayAlive$"},
		Env:     append(os.Environ(), reclaimHelperEnv+"=1"),
	}
}

// TestReclaimReleasesHandleOnce pins the ownership hand-over rule of the bounded
// termination design (docs/validation/evidence/P40-T027B-LOCK-REAPER/01-arch-design/review.md
// §4 risk 1): an uncertain Terminate keeps the platform handle, and Reclaim is the
// only path that may release it afterwards. Reclaim must release exactly once,
// stay idempotent, and refuse a process whose single reaper has not finished.
func TestReclaimReleasesHandleOnce(t *testing.T) {
	process, err := StartManaged(context.Background(), reclaimHelperSpec())
	if err != nil {
		t.Fatalf("start managed process: %v", err)
	}
	t.Cleanup(func() {
		_ = process.cmd.Process.Kill()
		process.platform.close()
	})

	// The tree cannot be proven gone during the stop, which is what makes the
	// termination uncertain and hands the handle over to a later reclaim.
	treeGone := false
	process.verify = func(ctx context.Context) error {
		if treeGone {
			return nil
		}
		<-ctx.Done()
		return ctx.Err()
	}

	proof, err := process.Terminate(context.Background(), 20*time.Millisecond)
	if !errors.Is(err, ErrTerminationUncertain) {
		t.Fatalf("terminate = %v, want ErrTerminationUncertain", err)
	}
	if proof.Outcome != "uncertain" {
		t.Fatalf("termination evidence outcome = %q, want uncertain", proof.Outcome)
	}
	if got := process.platform.closeCount(); got != 0 {
		t.Fatalf("an uncertain termination must keep the handle, closeCount = %d", got)
	}

	// Force the helper down so the single reaper finishes; a reclaim may only act
	// on a reaped tree.
	_ = process.cmd.Process.Kill()
	select {
	case <-process.waitResult():
	case <-time.After(10 * time.Second):
		t.Fatal("the killed helper was never reaped")
	}

	treeGone = true
	if err := process.Reclaim(context.Background()); err != nil {
		t.Fatalf("reclaim a reaped tree: %v", err)
	}
	if got := process.platform.closeCount(); got != 1 {
		t.Fatalf("reclaim must release the platform handle exactly once, closeCount = %d", got)
	}
	if err := process.Reclaim(context.Background()); err != nil {
		t.Fatalf("a second reclaim must be a no-op, got %v", err)
	}
	if got := process.platform.closeCount(); got != 1 {
		t.Fatalf("a second reclaim released the handle again, closeCount = %d", got)
	}

	// A process whose reaper has not finished is not reclaimable: ownership cannot
	// be released before the tree is proven reaped.
	running, err := StartManaged(context.Background(), reclaimHelperSpec())
	if err != nil {
		t.Fatalf("start second managed process: %v", err)
	}
	t.Cleanup(func() {
		_ = running.cmd.Process.Kill()
		running.platform.close()
	})
	if err := running.Reclaim(context.Background()); !errors.Is(err, ErrTerminationUncertain) {
		t.Fatalf("reclaim of a running process = %v, want ErrTerminationUncertain", err)
	}
	if got := running.platform.closeCount(); got != 0 {
		t.Fatalf("a running process must not be released, closeCount = %d", got)
	}
}
