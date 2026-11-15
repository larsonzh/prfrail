package console

import (
	"context"
	"fmt"
	"os"

	"github.com/larsonzh/prfrail/internal/guard"
)

// sandboxEnforcementEnvVar is the single S1 trigger of the sandbox check layer
// (CONTRACTS.md 7.2 (e)). Unset and "off" keep the pre-slice control flow;
// "enforce" turns the check on for both dispatch paths; any other value is a
// usage error with no side effect.
const sandboxEnforcementEnvVar = "PRFRAIL_SANDBOX_CHECK"

// The dispatch labels the two assembly entries pass to the check layer. They
// are the audit tag of a refusal, so each entry keeps its own constant.
const (
	agentRunnerSandboxDispatchReason = "agent-runner-dispatch"
	supervisedSandboxDispatchReason  = "supervised-runtime-dispatch"
)

// sandboxForDispatch is the check implementation of the dispatch paths. It is
// a package variable so white-box tests can replace it with a recording
// double; production always uses the platform implementation.
var sandboxForDispatch guard.Sandbox = guard.NewSandbox()

// preflightSandbox runs the sandbox check of one dispatch before the probe
// permit is claimed. A non-Allowed verdict or a check error refuses the
// dispatch fail-closed: the caller returns before it consumes the one-shot
// quota, before it creates the run-root marker and before anything can spawn.
func preflightSandbox(ctx context.Context, reason string) error {
	return guard.CheckSandbox(ctx, sandboxForDispatch, guard.CheckSpec{Reason: reason})
}

// sandboxEnforcementFromEnv resolves the S1 trigger. S1 defaults to off, so a
// host without the variable behaves exactly like the pre-slice build. The
// comparison is exact: only "enforce" and "off" (or an unset/empty value) are
// accepted, and every other value is refused so a typo can never silently
// disable an enforcement request.
func sandboxEnforcementFromEnv() (bool, error) {
	value, found := os.LookupEnv(sandboxEnforcementEnvVar)
	if !found {
		return false, nil
	}
	switch value {
	case "", "off":
		return false, nil
	case "enforce":
		return true, nil
	default:
		return false, fmt.Errorf("%s=%q is not a valid sandbox enforcement setting: use enforce or off", sandboxEnforcementEnvVar, value)
	}
}
