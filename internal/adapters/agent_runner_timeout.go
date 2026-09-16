package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	// ErrAgentRunnerTimeout reports a run whose watchdog expired. The run side
	// still has to stop the process tree and prove it; a timeout never becomes a
	// completed outcome on its own.
	ErrAgentRunnerTimeout = errors.New("AgentRunner run timed out")
	// ErrInvalidAgentRunnerTimeout reports a watchdog without a usable timeout.
	ErrInvalidAgentRunnerTimeout = errors.New("AgentRunner timeout invalid")
)

// agentRunnerTimeoutEvidenceDomain prefixes the digest a watchdog records when
// it expires, so the outcome can carry provable error evidence.
const agentRunnerTimeoutEvidenceDomain = "proofrail:agent-runner-timeout:1\n"

// AgentRunnerWatchResult reports how one watched run ended. TimedOut is true only
// when the watchdog itself expired: a caller cancellation is reported as a
// cancellation so the two races converge on one winner.
type AgentRunnerWatchResult struct {
	Err      error
	TimedOut bool
	Evidence []string
}

// AgentRunnerTimeoutManager bounds one run. It owns no process and never kills
// anything: the watched function receives the derived context and is responsible
// for stopping what it started, which is how a timeout becomes a provable stop
// instead of a silent abandon.
type AgentRunnerTimeoutManager struct {
	Timeout time.Duration
}

func (manager AgentRunnerTimeoutManager) RunWithWatchdog(ctx context.Context, requestID string, run func(context.Context) error) (AgentRunnerWatchResult, error) {
	if manager.Timeout <= 0 {
		return AgentRunnerWatchResult{}, fmt.Errorf("%w: timeout must be positive", ErrInvalidAgentRunnerTimeout)
	}
	if run == nil {
		return AgentRunnerWatchResult{}, fmt.Errorf("%w: nil run function", ErrInvalidAgentRunnerTimeout)
	}
	if err := ctx.Err(); err != nil {
		return AgentRunnerWatchResult{Err: err}, nil
	}
	// The watchdog's own deadline is the timer, and the derived context carries only
	// the cancellation. Deriving a second deadline here would create a race in which
	// the run returns the derived context's expiry and the watchdog then reports a
	// caller deadline instead of its own timeout: the two would be indistinguishable,
	// and a timeout misread as a cancellation is a wrong verdict.
	runCtx, cancelRun := context.WithCancel(ctx)
	defer cancelRun()
	done := make(chan error, 1)
	go func() { done <- run(runCtx) }()
	timer := time.NewTimer(manager.Timeout)
	defer timer.Stop()
	select {
	case err := <-done:
		// The run finished on its own. If the caller's context is already done the
		// caller stopped it, and that is the verdict the caller asked for.
		if callerErr := ctx.Err(); callerErr != nil {
			return AgentRunnerWatchResult{Err: callerErr}, nil
		}
		return AgentRunnerWatchResult{Err: err}, nil
	case <-ctx.Done():
		// Caller cancellation is not a timeout: the run was asked to stop.
		return AgentRunnerWatchResult{Err: ctx.Err()}, nil
	case <-timer.C:
		// The deadline is owned by the watchdog, so tell the run side to stop and let
		// it reconcile. A caller cancellation that lands at the very same instant must
		// still win, so it is given a bounded settling window instead of being decided
		// by which channel the scheduler happened to deliver first: an undecided race
		// would produce two different verdicts from one edge event.
		cancelRun()
		if err := ctx.Err(); err != nil {
			return AgentRunnerWatchResult{Err: err}, nil
		}
		if err := waitForCallerStop(ctx, agentRunnerTimeoutSettleWindow); err != nil {
			return AgentRunnerWatchResult{Err: err}, nil
		}
		deadline := time.Now().UTC().Add(-manager.Timeout).Format(time.RFC3339Nano)
		return AgentRunnerWatchResult{
			Err:      fmt.Errorf("%w after %s from %s", ErrAgentRunnerTimeout, manager.Timeout, deadline),
			TimedOut: true,
			Evidence: []string{evidence.Digest(agentRunnerTimeoutEvidenceDomain, []byte(requestID+"\n"+manager.Timeout.String()))},
		}, nil
	}
}

// agentRunnerTimeoutSettleWindow is how long a timeout waits for a simultaneous
// caller cancellation before it declares itself the winner. It only bounds the tie
// break, never the run: the window is far shorter than any real run bound.
const agentRunnerTimeoutSettleWindow = 20 * time.Millisecond

// waitForCallerStop reports the caller's cancellation when it lands inside the
// window, and nil when the caller stayed live: that is what makes the timeout the
// single winner whenever the cancellation did not actually arrive.
func waitForCallerStop(ctx context.Context, window time.Duration) error {
	if window <= 0 {
		return nil
	}
	timer := time.NewTimer(window)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
