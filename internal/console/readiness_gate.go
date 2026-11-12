package console

import (
	"context"

	"github.com/larsonzh/prfrail/internal/gates"
)

// requireProbePermit is the shared probe-permit gate of the production assembly
// entries. It evaluates level-1 readiness (a bound permit), claims the one-shot
// quota with Require and returns a release function that completes (spends) the
// permit exactly once. A failed evaluation or claim refuses the run fail-closed:
// no run-root marker and no spawned process exist.
func requireProbePermit(ctx context.Context, permit *gates.ProbePermit) (release func(), err error) {
	if _, err := gates.EvaluateT027BReadiness(gates.GoalProbe, gates.ReadinessInput{Permit: permit}); err != nil {
		return nil, err
	}
	if err := permit.Require(ctx); err != nil {
		return nil, err
	}
	return func() {
		// Complete can only fail when the permit was never required. Require
		// returned nil just above, so the permit is Claimed and the CAS to Spent
		// always succeeds; a permit that is already Spent is idempotent. The
		// error is therefore discarded on purpose: there is no observable
		// failure path to record.
		_ = permit.Complete()
	}, nil
}
