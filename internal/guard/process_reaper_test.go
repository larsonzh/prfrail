package guard

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"
)

// Test-only contract for the bounded-termination reaper slice
// (design: docs/validation/evidence/P40-T027B-LOCK-REAPER/01-arch-design/review.md §1/§3).
//
// These tests intentionally reference an API that does not exist yet:
// ManagedProcess.waitDone, ManagedProcess.waitResult and the bounded Terminate
// contract. The expected RED state of this slice is a compile failure.
const reaperHelperEnv = "PROOFRAIL_REAPER_HELPER"

// TestReaperHelperNeverExit is the child half of the reaper tests. It blocks
// forever so the parent can exercise the bounded stop/verification paths.
func TestReaperHelperNeverExit(t *testing.T) {
	if os.Getenv(reaperHelperEnv) == "" {
		t.Skip("reaper helper only runs as a child process")
	}
	select {}
}

func reaperHelperSpec() ProcessSpec {
	return ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=^TestReaperHelperNeverExit$"},
		Env:     append(os.Environ(), reaperHelperEnv+"=1"),
	}
}

// TestTerminationDeadlineReturnsUncertain pins the "no unbounded wait" rule:
// a tree the guard cannot prove gone must make Terminate return within the
// deadline with ErrTerminationUncertain, and it must keep ownership of the
// handle (the single reaper and the platform handle are not released).
func TestTerminationDeadlineReturnsUncertain(t *testing.T) {
	process, err := StartManaged(context.Background(), reaperHelperSpec())
	if err != nil {
		t.Fatalf("start managed process: %v", err)
	}
	t.Cleanup(func() {
		_ = process.cmd.Process.Kill()
		process.platform.close()
	})

	// Simulate an unkillable tree: verification can never prove it gone.
	process.verify = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}

	waitDone := process.waitDone // new field: the single reaper broadcast channel
	if waitDone == nil {
		t.Fatal("managed process must expose exactly one waitDone channel")
	}
	if got := process.waitResult(); got != waitDone {
		t.Fatal("waitResult must return the single reaper channel, not a new one")
	}
	identity := process.Identity()

	const deadline = 150 * time.Millisecond
	type outcome struct {
		proof TerminationEvidence
		err   error
	}
	result := make(chan outcome, 1)
	started := time.Now()
	go func() {
		proof, terminateErr := process.Terminate(context.Background(), deadline)
		result <- outcome{proof: proof, err: terminateErr}
	}()

	var got outcome
	select {
	case got = <-result:
	case <-time.After(5 * time.Second):
		t.Fatal("Terminate blocked past its deadline: unbounded wait on an unkillable tree")
	}
	elapsed := time.Since(started)
	if elapsed > deadline+2*time.Second {
		t.Fatalf("Terminate returned after %s, far past its %s deadline", elapsed, deadline)
	}
	if !errors.Is(got.err, ErrTerminationUncertain) {
		t.Fatalf("Terminate must report ErrTerminationUncertain at the deadline, got %v", got.err)
	}
	if got.proof.Outcome != "uncertain" {
		t.Fatalf("termination evidence outcome = %q, want uncertain", got.proof.Outcome)
	}
	// Ownership is retained: the single reaper is the same channel and the
	// platform handle is still attached, so a later Reclaim can finish the job.
	if again := process.waitResult(); again != waitDone {
		t.Fatal("an uncertain termination must not replace the single reaper channel")
	}
	if process.Identity() != identity {
		t.Fatal("an uncertain termination must keep the managed process identity owned")
	}
}

// TestConcurrentRunAndTerminateUseOneWaiter pins the single-reaper invariant:
// concurrent Run and Terminate calls share one cmd.Wait, so no second Wait can
// panic and no caller can race the reaper for the terminal state.
func TestConcurrentRunAndTerminateUseOneWaiter(t *testing.T) {
	process, err := StartManaged(context.Background(), reaperHelperSpec())
	if err != nil {
		t.Fatalf("start managed process: %v", err)
	}
	t.Cleanup(func() {
		_ = process.cmd.Process.Kill()
	})
	process.verify = func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-process.waitDone:
			return nil
		}
	}
	waitDone := process.waitDone
	if waitDone == nil {
		t.Fatal("managed process must expose exactly one waitDone channel")
	}

	startedAt := time.Now().UTC()
	panics := make(chan any, 4)
	var wg sync.WaitGroup
	launch := func(fn func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if recovered := recover(); recovered != nil {
					panics <- recovered
				}
			}()
			fn()
		}()
	}

	launch(func() { _, _ = process.Run(context.Background(), 100*time.Millisecond, startedAt) })
	launch(func() { _, _ = process.Terminate(context.Background(), 100*time.Millisecond) })
	launch(func() { _, _ = process.Terminate(context.Background(), 100*time.Millisecond) })

	converged := make(chan struct{})
	go func() { wg.Wait(); close(converged) }()
	select {
	case <-converged:
	case <-time.After(10 * time.Second):
		t.Fatal("concurrent Run/Terminate did not converge: the single reaper is not shared")
	}
	select {
	case recovered := <-panics:
		t.Fatalf("concurrent wait panicked (a second cmd.Wait was issued): %v", recovered)
	default:
	}
	if again := process.waitResult(); again != waitDone {
		t.Fatal("concurrent callers must observe the same single reaper channel")
	}
}
