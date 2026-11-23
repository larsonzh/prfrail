package adapters

// P85 F18 (design rev 3, section 7 / test row M11): the platform write verdict must
// be the only gate between a COMPLIANT credential-bearing envelope and the start
// observation.
//
// Review r1 rejected the earlier shape because the Allowed stand-in never entered the
// real executor: the counterexample used a shell argv that the argv allowlist refused
// on its own, so no verdict change was measurable. The fix puts an injectable check
// layer on the executor and an observation point at the LAST gate
// (copilotCLIProbeRunManaged), so the two legs below differ in ONE fact - the verdict
// - and the child is never really started (the stub only counts).
//
// Mutation check (the single-variant body of M11): make the executor ignore its
// Sandbox field (or delete the final guard.CheckSandbox call) and the Unknown leg
// starts the child => "a-non-allowed-verdict-keeps-the-start-count-at-zero" goes red.
//
// The element-level control (M12: a non-compliant argv is still refused under an
// Allowed verdict, by its own element and not by the sandbox) lives in
// TestP85T8APlatformVerdictAloneControlsTheSpawn.

import (
	"context"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/guard"
)

// p85F18StubStart replaces the managed-process launcher with a counter, so "the child
// would have been created" becomes an integer and nothing is spawned.
func p85F18StubStart(t *testing.T) *int {
	t.Helper()
	previous := copilotCLIProbeRunManaged
	starts := 0
	copilotCLIProbeRunManaged = func(_ context.Context, _ guard.ProcessSpec, _ time.Duration) (guard.ProcessResult, error) {
		starts++
		return guard.ProcessResult{Started: true}, nil
	}
	t.Cleanup(func() { copilotCLIProbeRunManaged = previous })
	return &starts
}

func TestP85F18ThePlatformVerdictAloneDecidesTheStart(t *testing.T) {
	// The executable must pass the LIVE install-source re-verification Execute runs
	// before the BYOK gates, otherwise the legs would both be refused by that check
	// and neither would isolate the verdict.
	shell := p85T8ASystemShell(t)
	root := t.TempDir()
	profile := p85T8ABYOKProfile()
	invocation := CopilotCLIProbeInvocation{
		Executable:       shell,
		WorkingDirectory: root,
		Args:             copilotCLIProbeArgs(profile.Model, root, 1, "secret"),
		Environment:      copilotCLIProbeEnvironment(profile, []byte("p85-t8a-secret")),
		AuthMode:         "secret",
	}
	starts := p85F18StubStart(t)

	t.Run("an-allowed-verdict-reaches-the-start-observation", func(t *testing.T) {
		*starts = 0
		executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: root, Sandbox: p85T8AAllowedSandbox{}}
		execution := executor.Execute(context.Background(), invocation)
		if *starts != 1 {
			t.Fatalf("an Allowed verdict on a compliant envelope must reach the start observation exactly once, got %d (outcome=%q evidence=%v)",
				*starts, execution.Outcome, execution.Evidence)
		}
		if execution.Outcome == "start-failed" && execution.FailureCode == copilotCLIProbeBYOKFailure {
			t.Fatalf("the BYOK gate must not refuse a compliant envelope under an Allowed verdict, got %+v", execution)
		}
	})

	t.Run("a-non-allowed-verdict-keeps-the-start-count-at-zero", func(t *testing.T) {
		*starts = 0
		unknown := &p85T8ASandboxDouble{verdict: guard.Verdict{Status: guard.VerdictUnknown, Reason: "S1 sandbox enforcement checks not implemented"}}
		executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: root, Sandbox: unknown}
		execution := executor.Execute(context.Background(), invocation)
		if *starts != 0 {
			t.Fatalf("an Unknown platform verdict must keep the child from starting, got %d starts", *starts)
		}
		if execution.Outcome != "start-failed" || execution.FailureCode != copilotCLIProbeBYOKFailure {
			t.Fatalf("an Unknown verdict must refuse under the BYOK class, got outcome=%q failure=%q", execution.Outcome, execution.FailureCode)
		}
		want := copilotCLIExecutorEvidence("byok-confinement-unproven")
		if len(execution.Evidence) != 1 || execution.Evidence[0] != want {
			t.Fatalf("the refusal must be the closed classification %q, got %v", want, execution.Evidence)
		}
	})

	t.Run("the-observation-point-is-the-last-gate", func(t *testing.T) {
		*starts = 0
		executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: root, Sandbox: &p85T8ASandboxDouble{verdict: guard.Verdict{Status: guard.VerdictDenied, Reason: "the environment provably lacks write restrictions"}}}
		if execution := executor.Execute(context.Background(), invocation); *starts != 0 || execution.Outcome != "start-failed" {
			t.Fatalf("a Denied verdict must refuse before the launcher, got starts=%d outcome=%q", *starts, execution.Outcome)
		}
	})
}
