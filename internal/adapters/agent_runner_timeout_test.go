package adapters

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestAgentRunnerWatchdogTimesOutWithBoundedEvidence(t *testing.T) {
	manager := AgentRunnerTimeoutManager{Timeout: 120 * time.Millisecond}
	var childCancelled atomic.Bool
	var observed error
	var sawDeadline atomic.Bool
	started := time.Now()
	result, err := manager.RunWithWatchdog(context.Background(), "request-one", func(ctx context.Context) error {
		if _, hasDeadline := ctx.Deadline(); hasDeadline {
			sawDeadline.Store(true)
		}
		<-ctx.Done()
		observed = ctx.Err()
		childCancelled.Store(true)
		return observed
	})
	if err != nil {
		t.Fatal(err)
	}
	// The watchdog owns the deadline, so the run side can only ever be cancelled by
	// it - never expired by a derived deadline of its own. If the run could expire
	// itself, the watchdog would be racing a second timeout source and could report
	// a caller deadline instead of its own timeout. The assertions come after the
	// bounded wait below, because the run goroutine may not have been scheduled yet.
	defer func() {
		if sawDeadline.Load() {
			t.Fatal("the run side must not carry a deadline of its own")
		}
		if !errors.Is(observed, context.Canceled) {
			t.Fatalf("the run side must be cancelled, not expired: %v", observed)
		}
	}()
	elapsed := time.Since(started)
	if !result.TimedOut || !errors.Is(result.Err, ErrAgentRunnerTimeout) {
		t.Fatalf("an expired watchdog must report a timeout: %+v", result)
	}
	if len(result.Evidence) != 1 || !evidence.ValidHash(result.Evidence[0]) {
		t.Fatalf("a timeout must carry one evidence digest: %v", result.Evidence)
	}
	if !childCancelled.Load() {
		// The watched goroutine may not have been scheduled yet when the watchdog
		// returns, so the cancellation is observed with a bounded wait.
		waitUntil := time.Now().Add(5 * time.Second)
		for !childCancelled.Load() && time.Now().Before(waitUntil) {
			time.Sleep(5 * time.Millisecond)
		}
	}
	if !childCancelled.Load() {
		t.Fatal("the watched run must observe the derived cancellation")
	}
	// The lower bound proves the watchdog waited for its own timeout instead of
	// returning early. The upper bound only guards against a hang: under a fully
	// loaded host a timer fires late, so a tight bound would flake rather than
	// catch a defect (Windows monotonic granularity is about 15.6ms).
	if elapsed < 100*time.Millisecond || elapsed > 10*time.Second {
		t.Fatalf("watchdog elapsed %v is outside the tolerated range", elapsed)
	}
}

func TestAgentRunnerWatchdogPrefersCallerCancellation(t *testing.T) {
	manager := AgentRunnerTimeoutManager{Timeout: 5 * time.Second}
	ctx, cancel := context.WithCancel(context.Background())
	released := make(chan struct{})
	result := make(chan AgentRunnerWatchResult, 1)
	go func() {
		watch, err := manager.RunWithWatchdog(ctx, "request-one", func(ctx context.Context) error {
			<-released
			return ctx.Err()
		})
		if err != nil {
			t.Errorf("unexpected watchdog error: %v", err)
		}
		result <- watch
	}()
	time.Sleep(50 * time.Millisecond)
	cancel()
	watch := <-result
	close(released)
	if watch.TimedOut {
		t.Fatalf("a caller cancellation must not be reported as a timeout: %+v", watch)
	}
	if !errors.Is(watch.Err, context.Canceled) {
		t.Fatalf("a caller cancellation must be reported as a cancellation, got %v", watch.Err)
	}
	if len(watch.Evidence) != 0 {
		t.Fatalf("a cancellation carries no timeout evidence: %v", watch.Evidence)
	}
}

func TestAgentRunnerWatchdogReturnsTheRunError(t *testing.T) {
	manager := AgentRunnerTimeoutManager{Timeout: 2 * time.Second}
	sentinel := errors.New("run failed")
	result, err := manager.RunWithWatchdog(context.Background(), "request-one", func(context.Context) error {
		return sentinel
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.TimedOut || !errors.Is(result.Err, sentinel) {
		t.Fatalf("the run error must pass through unchanged: %+v", result)
	}
	if len(result.Evidence) != 0 {
		t.Fatalf("no timeout evidence may be claimed for a plain run error: %v", result.Evidence)
	}
}

func TestAgentRunnerWatchdogRefusesUnusableTimeouts(t *testing.T) {
	calls := 0
	for name, manager := range map[string]AgentRunnerTimeoutManager{
		"zero timeout": {Timeout: 0},
		"negative":     {Timeout: -time.Second},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := manager.RunWithWatchdog(context.Background(), "request-one", func(context.Context) error {
				calls++
				return nil
			}); !errors.Is(err, ErrInvalidAgentRunnerTimeout) {
				t.Fatalf("a watchdog without a timeout must fail closed, got %v", err)
			}
		})
	}
	if _, err := (AgentRunnerTimeoutManager{Timeout: time.Second}).RunWithWatchdog(context.Background(), "request-one", nil); !errors.Is(err, ErrInvalidAgentRunnerTimeout) {
		t.Fatalf("a nil run function must be refused, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("a refused watchdog must never run the function: %d calls", calls)
	}
	// An already-cancelled caller must not start the run at all: starting it would
	// spawn work whose result nobody can use. The assertion waits for the call
	// instead of trusting an ordering, so disabling the short-circuit reddens it.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	started := make(chan struct{}, 1)
	result, err := (AgentRunnerTimeoutManager{Timeout: time.Second}).RunWithWatchdog(ctx, "request-one", func(context.Context) error {
		started <- struct{}{}
		return nil
	})
	if err != nil || !errors.Is(result.Err, context.Canceled) {
		t.Fatalf("an already-cancelled caller must short-circuit: result=%+v err=%v", result, err)
	}
	select {
	case <-started:
		t.Fatal("an already-cancelled caller must never start the run")
	case <-time.After(200 * time.Millisecond):
	}
}

func TestAgentRunnerWatchdogRaceClassificationIsSingleValued(t *testing.T) {
	for round := 0; round < 20; round++ {
		manager := AgentRunnerTimeoutManager{Timeout: time.Duration(15+round) * time.Millisecond}
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15+round)*time.Millisecond)
		result, err := manager.RunWithWatchdog(ctx, "request-one", func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		})
		cancel()
		if err != nil {
			t.Fatal(err)
		}
		if result.Err == nil {
			t.Fatalf("round %d produced no outcome", round)
		}
		if result.TimedOut {
			if !errors.Is(result.Err, ErrAgentRunnerTimeout) || len(result.Evidence) != 1 {
				t.Fatalf("round %d: a timeout must carry its evidence: %+v", round, result)
			}
			continue
		}
		if errors.Is(result.Err, ErrAgentRunnerTimeout) || len(result.Evidence) != 0 {
			t.Fatalf("round %d: a cancellation must not look like a timeout: %+v", round, result)
		}
	}
}
