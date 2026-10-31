//go:build linux

package guard

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
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

// proofrailGuardOrphanHelperEnv marks the two child halves of the remaining-tree
// test: "parent" starts a child in the same process group and exits at once,
// "child" stays alive so the tree outlives its parent.
const proofrailGuardOrphanHelperEnv = "PROOFRAIL_GUARD_ORPHAN_HELPER"

func proofrailGuardOrphanSpec(value string) ProcessSpec {
	env := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, proofrailGuardOrphanHelperEnv+"=") {
			continue
		}
		env = append(env, entry)
	}
	env = append(env, proofrailGuardOrphanHelperEnv+"="+value)
	return ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=^TestLinuxNaturalExitReportsRemainingTree$"},
		Env:     env,
	}
}

// TestLinuxNaturalExitReportsRemainingTree pins the natural-exit branch that has
// no Windows equivalent: the managed parent exits while a descendant keeps the
// process group alive, so the exit is reported as ErrProcessTreeRemained with a
// termination proof, and the handle is released exactly once after the group is
// proven gone.
func TestLinuxNaturalExitReportsRemainingTree(t *testing.T) {
	switch os.Getenv(proofrailGuardOrphanHelperEnv) {
	case "parent":
		child := exec.Command(os.Args[0], "-test.run=^TestLinuxNaturalExitReportsRemainingTree$")
		child.Env = proofrailGuardOrphanSpec("child").Env
		if err := child.Start(); err != nil {
			os.Exit(3)
		}
		// The parent exits without waiting: the child is left behind in the same
		// process group, which is exactly the shape runManaged must report.
		return
	case "child":
		time.Sleep(30 * time.Second)
		return
	}

	process, err := StartManaged(context.Background(), proofrailGuardOrphanSpec("parent"))
	if err != nil {
		t.Fatalf("start managed process: %v", err)
	}
	t.Cleanup(func() { process.platform.close() })

	startedAt := time.Now().UTC()
	result, err := process.Run(context.Background(), time.Second, startedAt)
	if !errors.Is(err, ErrProcessTreeRemained) {
		t.Fatalf("a natural exit that leaves a descendant = %v, want ErrProcessTreeRemained", err)
	}
	if result.Termination == nil {
		t.Fatal("a remaining tree must still carry termination evidence")
	}
	if result.Termination.Outcome != "stopped" {
		t.Fatalf("termination outcome = %q, want stopped after the tree was stopped", result.Termination.Outcome)
	}
	if !evidence.ValidHash(result.Termination.Hash) {
		t.Fatalf("termination evidence hash is invalid: %+v", result.Termination)
	}
	stopped := false
	for _, action := range result.Termination.Actions {
		if action == "signal-term-process-group" {
			stopped = true
		}
	}
	if !stopped {
		t.Fatalf("the remaining tree must be stopped through the process group: %v", result.Termination.Actions)
	}
	if got := process.platform.closeCount(); got != 1 {
		t.Fatalf("the handle must be released exactly once, closeCount = %d", got)
	}
}
