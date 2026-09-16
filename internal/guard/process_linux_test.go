//go:build linux

package guard

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

// proofrailGuardZombieHelperEnv marks the child half of the zombie test: the
// child only has to stay alive until the test kills it.
const proofrailGuardZombieHelperEnv = "PROOFRAIL_GUARD_ZOMBIE_HELPER"

// TestLinuxStopProvesAnUnreapedTerminatedChild pins the platform rule that a
// process which has exited but whose parent has not reaped it is not alive. Its
// /proc entry, its start token and its process group all survive until then, so
// reading a zombie as alive made every identity-only stop report an uncertain
// termination whenever the stopper did not own the child handle - which is
// exactly the pinned-CLI launcher path, since it stops by persisted identity and
// never reaps.
func TestLinuxStopProvesAnUnreapedTerminatedChild(t *testing.T) {
	if os.Getenv(proofrailGuardZombieHelperEnv) != "" {
		time.Sleep(30 * time.Second)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestLinuxStopProvesAnUnreapedTerminatedChild$")
	cmd.Env = append(os.Environ(), proofrailGuardZombieHelperEnv+"=1")
	// The child gets its own process group, so the group probe can be asserted on
	// without touching the test binary's own group.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	// The child is deliberately not waited on before the assertions: staying
	// unreaped is the condition under test.
	reaped := false
	defer func() {
		if !reaped {
			_, _ = cmd.Process.Wait()
		}
	}()

	identity, err := platformIdentity(cmd.Process.Pid)
	if err != nil {
		t.Fatal(err)
	}
	if alive, err := platformIdentityAlive(identity); err != nil || !alive {
		t.Fatalf("a running child must read as alive: alive=%v err=%v", alive, err)
	}
	if err := syscall.Kill(cmd.Process.Pid, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		alive, err := platformIdentityAlive(identity)
		if err != nil {
			t.Fatal(err)
		}
		if !alive {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if alive, err := platformIdentityAlive(identity); err != nil || alive {
		t.Fatalf("an exited child must not read as alive before it is reaped: alive=%v err=%v", alive, err)
	}
	if processGroupAlive(cmd.Process.Pid) {
		t.Fatal("a process group whose only member is an unreaped zombie must not read as alive")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	proof, err := StopProcessIdentity(ctx, identity, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("stopping an already terminated child must be proven: %v (%+v)", err, proof)
	}
	if proof.Outcome != "stopped" {
		t.Fatalf("expected a proven stop, got %+v", proof)
	}
	if len(proof.Actions) != 1 || proof.Actions[0] != "already-stopped" {
		t.Fatalf("an exited child must be reported as already stopped, got %+v", proof.Actions)
	}

	// Reap only after the assertions, so the rule under test is not satisfied by
	// the parent's own wait.
	if _, err := cmd.Process.Wait(); err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("wait for the killed child: %v", err)
		}
	}
	reaped = true
}
