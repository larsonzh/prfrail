//go:build windows

package guard

import (
	"context"
	"errors"
	"testing"
	"time"
)

// hasAction reports whether the termination evidence lists one action marker.
func hasAction(actions []string, want string) bool {
	for _, action := range actions {
		if action == want {
			return true
		}
	}
	return false
}

// TestPlatformStopIdentityExpiredBoundReportsDetached pins the detached-stop
// contract of the Windows identity stop path: taskkill is an external command and
// is never killed when it outlives its synchronous bound, so reaching the bound
// must be reported as "unknown" (errStopCommandExpired plus the detached action)
// instead of as a failure. The caller proves the tree itself afterwards.
func TestPlatformStopIdentityExpiredBoundReportsDetached(t *testing.T) {
	process, err := StartManaged(context.Background(), reclaimHelperSpec())
	if err != nil {
		t.Fatalf("start managed process: %v", err)
	}
	t.Cleanup(func() {
		_ = process.cmd.Process.Kill()
		process.platform.close()
	})

	// A nanosecond bound expires before taskkill can even start, so "the command
	// outlived its bound" is the only possible outcome of this call.
	actions, err := platformStopIdentity(context.Background(), process.Identity(), time.Second, time.Nanosecond)
	if !errors.Is(err, errStopCommandExpired) {
		t.Fatalf("stop command past its bound = %v, want errStopCommandExpired", err)
	}
	if !hasAction(actions, "taskkill-detached") {
		t.Fatalf("detached action missing from %v", actions)
	}
	if !hasAction(actions, "taskkill-process-tree") {
		t.Fatalf("the command that was left running must still be named: %v", actions)
	}
}

// TestStopProcessIdentitySettlesVerifiedGone pins that the identity stop settles
// on a proven result: a live tree is stopped through taskkill and then proven
// gone by the bounded identity wait, and an already finished tree is reported as
// already stopped. Neither case may report an uncertain outcome, and no assertion
// here depends on wall-clock latency.
func TestStopProcessIdentitySettlesVerifiedGone(t *testing.T) {
	t.Run("live-tree-is-stopped-and-proven-gone", func(t *testing.T) {
		process, err := StartManaged(context.Background(), reclaimHelperSpec())
		if err != nil {
			t.Fatalf("start managed process: %v", err)
		}
		t.Cleanup(func() {
			_ = process.cmd.Process.Kill()
			process.platform.close()
		})
		identity := process.Identity()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		proof, err := StopProcessIdentity(ctx, identity, 500*time.Millisecond)
		if err != nil {
			t.Fatalf("stop identity = %v (%+v)", err, proof)
		}
		if proof.Outcome != "stopped" {
			t.Fatalf("termination outcome = %q, want stopped", proof.Outcome)
		}
		if !hasAction(proof.Actions, "taskkill-process-tree") {
			t.Fatalf("the taskkill path must be named in %v", proof.Actions)
		}
		if alive, aliveErr := ProcessAlive(identity); aliveErr != nil || alive {
			t.Fatalf("a settled stop must leave the identity gone: alive=%v err=%v", alive, aliveErr)
		}
	})

	t.Run("already-exited-tree-reports-already-stopped", func(t *testing.T) {
		process, err := StartManaged(context.Background(), reclaimHelperSpec())
		if err != nil {
			t.Fatalf("start managed process: %v", err)
		}
		t.Cleanup(func() { process.platform.close() })
		identity := process.Identity()

		_ = process.cmd.Process.Kill()
		deadline := time.Now().Add(10 * time.Second)
		for time.Now().Before(deadline) {
			alive, aliveErr := ProcessAlive(identity)
			if aliveErr != nil {
				t.Fatalf("probe the killed helper: %v", aliveErr)
			}
			if !alive {
				break
			}
			time.Sleep(20 * time.Millisecond)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		proof, err := StopProcessIdentity(ctx, identity, 500*time.Millisecond)
		if err != nil {
			t.Fatalf("stop an already exited identity = %v (%+v)", err, proof)
		}
		if proof.Outcome != "stopped" {
			t.Fatalf("termination outcome = %q, want stopped", proof.Outcome)
		}
		if !hasAction(proof.Actions, "already-stopped") {
			t.Fatalf("an exited tree must be reported as already stopped: %v", proof.Actions)
		}
	})
}
