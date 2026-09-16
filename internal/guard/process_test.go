package guard

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestManagedProcessTerminatesDescendantTree(t *testing.T) {
	if os.Getenv("PROOFRAIL_PROCESS_HELPER") != "" {
		runProcessHelper()
		return
	}
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	managed, err := StartManaged(context.Background(), ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=TestManagedProcessTerminatesDescendantTree", "--", "parent", pidFile},
		Env:     append(os.Environ(), "PROOFRAIL_PROCESS_HELPER=parent"),
	})
	if err != nil {
		t.Fatal(err)
	}
	var childIdentity ProcessIdentity
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		data, readErr := os.ReadFile(pidFile)
		if readErr == nil {
			pid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
			if parseErr == nil {
				childIdentity, err = platformIdentity(pid)
				if err == nil {
					break
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	if childIdentity.PID == 0 {
		t.Fatal("child process identity was not observed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	proof, err := managed.Terminate(ctx, 100*time.Millisecond)
	if err != nil {
		t.Fatalf("terminate tree: %v (%+v)", err, proof)
	}
	if proof.Outcome != "stopped" || !evidence.ValidHash(proof.Hash) {
		t.Fatalf("invalid termination evidence: %+v", proof)
	}
	alive, err := platformIdentityAlive(childIdentity)
	if err != nil || alive {
		t.Fatalf("descendant remained alive: alive=%v err=%v", alive, err)
	}
}

func TestStopProcessIdentityStopsRunningProcess(t *testing.T) {
	if os.Getenv("PROOFRAIL_PROCESS_HELPER") != "" {
		runProcessHelper()
		return
	}
	pidFile := filepath.Join(t.TempDir(), "unused.pid")
	managed, err := StartManaged(context.Background(), ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=^TestStopProcessIdentityStopsRunningProcess$", "--", "child", pidFile},
		Env:     append(os.Environ(), "PROOFRAIL_PROCESS_HELPER=child"),
	})
	if err != nil {
		t.Fatal(err)
	}
	waitDone := make(chan error, 1)
	go func() {
		_, waitErr := managed.cmd.Process.Wait()
		waitDone <- waitErr
	}()
	defer func() {
		_ = managed.cmd.Process.Kill()
		managed.platform.close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	proof, err := StopProcessIdentity(ctx, managed.Identity(), 100*time.Millisecond)
	if err != nil {
		t.Fatalf("stop identity failed: %v (%+v)", err, proof)
	}
	waitErr := <-waitDone
	var exitErr *exec.ExitError
	if waitErr != nil && !errors.As(waitErr, &exitErr) {
		t.Fatalf("wait stopped process: %v", waitErr)
	}
	if proof.Outcome != "stopped" || !evidence.ValidHash(proof.Hash) {
		t.Fatalf("invalid stop proof: %+v", proof)
	}
	alive, err := ProcessAlive(managed.Identity())
	if err != nil || alive {
		t.Fatalf("process should be stopped: alive=%v err=%v", alive, err)
	}
}

// errUnboundedVerifyContext marks a stop-verification call that was handed a
// context without a deadline. Bounded verification is a hard requirement, so the
// injected verifier fails fast instead of blocking forever.
var errUnboundedVerifyContext = errors.New("stop verification context is unbounded")

// TestRunManagedCancellationVerifiesStopWithinBoundedContext pins the
// cancellation path. The caller's context is already cancelled there, so the
// verification must receive its own bounded context: replacing it with an
// unbounded background context makes every such call unbounded, which is exactly
// what an unkillable descendant would turn into a permanent hang.
func TestRunManagedCancellationVerifiesStopWithinBoundedContext(t *testing.T) {
	if os.Getenv("PROOFRAIL_PROCESS_HELPER") != "" {
		runProcessHelper()
		return
	}
	pidFile := filepath.Join(t.TempDir(), "unused.pid")
	process, err := StartManaged(context.Background(), ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=^TestRunManagedCancellationVerifiesStopWithinBoundedContext$", "--", "child", pidFile},
		Env:     append(os.Environ(), "PROOFRAIL_PROCESS_HELPER=child"),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = process.cmd.Process.Kill()
		process.platform.close()
	}()

	const grace = 200 * time.Millisecond
	var (
		mu            sync.Mutex
		calls         int
		unboundedSeen bool
		longestBound  time.Duration
	)
	process.verify = func(ctx context.Context) error {
		deadline, ok := ctx.Deadline()
		mu.Lock()
		calls++
		if !ok {
			unboundedSeen = true
		} else if bound := time.Until(deadline); bound > longestBound {
			longestBound = bound
		}
		mu.Unlock()
		if !ok {
			return errUnboundedVerifyContext
		}
		<-ctx.Done()
		return ctx.Err()
	}

	ctx, cancel := context.WithCancel(context.Background())
	results := make(chan ProcessResult, 1)
	errs := make(chan error, 1)
	go func() {
		result, runErr := process.runManaged(ctx, grace, time.Now().UTC())
		results <- result
		errs <- runErr
	}()
	time.Sleep(300 * time.Millisecond)
	cancel()

	select {
	case result := <-results:
		// A verified stop that cannot complete is reported as an uncertain
		// termination; the cancellation itself is the caller-facing reason.
		if runErr := <-errs; !errors.Is(runErr, ErrTerminationUncertain) && !errors.Is(runErr, context.Canceled) {
			t.Fatalf("a cancelled run must report the cancellation or an uncertain termination, got %v", runErr)
		}
		if result.Termination == nil || result.Termination.Outcome != "uncertain" {
			t.Fatalf("a stop that could not be verified must be uncertain: %+v", result.Termination)
		}
	case <-time.After(grace + 10*time.Second):
		t.Fatal("runManaged did not return after cancellation")
	}

	mu.Lock()
	defer mu.Unlock()
	if unboundedSeen {
		t.Fatal("stop verification was handed a context without a deadline")
	}
	if calls == 0 {
		t.Fatal("stop verification never ran")
	}
	if longestBound <= 0 || longestBound > grace+time.Second {
		t.Fatalf("stop verification bound %v is not derived from the grace period %v", longestBound, grace)
	}
}

func runProcessHelper() {
	separator := 0
	for index, arg := range os.Args {
		if arg == "--" {
			separator = index
			break
		}
	}
	if separator == 0 || len(os.Args) <= separator+2 {
		os.Exit(2)
	}
	mode, pidFile := os.Args[separator+1], os.Args[separator+2]
	if mode == "parent" {
		child := exec.Command(os.Args[0], "-test.run=TestManagedProcessTerminatesDescendantTree", "--", "child", pidFile)
		child.Env = append(os.Environ(), "PROOFRAIL_PROCESS_HELPER=child")
		if err := child.Start(); err != nil {
			os.Exit(3)
		}
		if err := os.WriteFile(pidFile, fmt.Appendf(nil, "%d\n", child.Process.Pid), 0600); err != nil {
			os.Exit(4)
		}
	}
	for {
		time.Sleep(time.Hour)
	}
}
