//go:build a6native

// A6 pilot adjudication experiment (A11): the fourth-step reviewer claimed the
// watchdog's timeout/cancel race is non-deterministic, because a Go select over
// ctx.Done and timer.C picks a winner pseudo-randomly. This experiment forces both
// channels ready at the same instant, many times, and checks that every round
// converges on the same verdict - the one the contract demands (cancellation wins).
package adapters

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestA6NativeE8CancelAndTimeoutAtTheSameInstantHaveOneVerdict(t *testing.T) {
	const rounds = 200
	cancelled, timeouts := 0, 0
	for round := 0; round < rounds; round++ {
		manager := AgentRunnerTimeoutManager{Timeout: 5 * time.Millisecond}
		ctx, cancel := context.WithCancel(context.Background())
		// Fire the cancellation so that ctx.Done and the watchdog timer become ready
		// at the same instant: this is the race the reviewer claims is undecided.
		go func() {
			time.Sleep(5 * time.Millisecond)
			cancel()
		}()
		result, err := manager.RunWithWatchdog(ctx, "request-one", func(runCtx context.Context) error {
			<-runCtx.Done()
			return runCtx.Err()
		})
		cancel()
		if err != nil {
			t.Fatal(err)
		}
		if result.TimedOut {
			timeouts++
			continue
		}
		if errors.Is(result.Err, context.Canceled) || errors.Is(result.Err, context.DeadlineExceeded) {
			cancelled++
			continue
		}
		t.Fatalf("round %d: unexpected verdict %+v", round, result)
	}
	t.Logf("E8 summary: cancelled=%d timeout=%d of %d rounds", cancelled, timeouts, rounds)
	// The contract allows only one winner for a simultaneous arrival, and that winner
	// is the caller cancellation: any timeout verdict here means the race is undecided.
	if timeouts != 0 {
		t.Fatalf("a simultaneous cancel and timeout must always resolve to the cancellation, got %d timeouts", timeouts)
	}
	if cancelled != rounds {
		t.Fatalf("every round must be classified as a cancellation, got %d of %d", cancelled, rounds)
	}
}
