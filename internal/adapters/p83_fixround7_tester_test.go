package adapters

// P83 FIX ROUND 7 counterexamples (items T1a and T8b).
//
// T1a drives the real AdmitAgentRunner with a ledger source that consumes the
// grant's remaining lifetime while it loads: the timestamp has to be taken AFTER
// every input, otherwise a grant that expired during the read still passes.
//
// T8b drives the real CopilotCLIProbe.Probe refusal path and compares the final
// evidence of two different refused configurations byte for byte.
//
// Mutation contracts are stated on each test.

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
)

// TestP83Fix7T1aExpiryDuringTheLedgerReadIsRefused is the T1a counterexample: the
// grant is alive when the call starts and expires while the authorization ledger is
// being read, so the admission decision has to judge the clock taken after the read
// and refuse.
//
// Mutation check: moving `now := admission.Clock().UTC()` above the ledger source
// reads judges the grant against the pre-read timestamp, admits the step and turns
// this red.
func TestP83Fix7T1aExpiryDuringTheLedgerReadIsRefused(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	ctx := context.Background()
	// Control: the untouched fixture admits, so the refusal below is the clock order.
	if err := admission.AdmitAgentRunner(ctx, request); err != nil {
		t.Fatalf("the fixture must admit before the clock is advanced, got %v", err)
	}
	// The grant expires at 2026-09-11T09:00:00.000Z and the fixture clock reads
	// 08:05:00.000Z; the read consumes the remaining 55 minutes. The availability
	// policy is widened so the availability guard cannot be the refusal under test:
	// the expiration of the AUTHORIZATION is what has to be judged with the
	// post-read clock.
	expiresAt := time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC)
	admission.AvailabilityPolicy.MaximumAge = 4 * time.Hour
	ledger := admission.AuthorizationLedger
	if len(ledger.Grants) != 1 || ledger.Grants[0].Grant == nil {
		t.Fatal("test precondition failed: the fixture must carry one grant")
	}
	admission.AuthorizationLedgerSource = func() (chain.AuthorizationLedger, error) {
		admission.Clock = func() time.Time { return expiresAt }
		return ledger, nil
	}

	err := admission.AdmitAgentRunner(ctx, request)
	if !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("an authorization that expires during the ledger read must be refused, got %v", err)
	}
	if !strings.Contains(err.Error(), "expired") {
		t.Fatalf("the refusal must name the expiry, got %v", err)
	}
}

// TestP83Fix7T8bProbeRefusalEvidenceIsValueIndependent is the T8b counterexample at
// the Probe seam: two start-failed refusals of DIFFERENT configurations produce the
// same final evidence, because a refusal names its class and never the values it
// refused.
//
// Mutation check: restoring the value-carrying builder in the start-failed branch of
// Probe (copilotCLIProbeEvidence(invocation, execution) instead of the refusal
// evidence) digests the executable, the working directory and the argv, so the two
// evidence slices differ and this turns red.
func TestP83Fix7T8bProbeRefusalEvidenceIsValueIndependent(t *testing.T) {
	results := make([]AIProbeResult, 0, 2)
	for _, configuration := range []struct {
		executable       string
		workingDirectory string
		model            string
	}{
		{executable: `C:\bin\copilot-a.exe`, workingDirectory: `C:\probe-a`, model: "deepseek-flash"},
		{executable: `/opt/copilot-b/copilot`, workingDirectory: `/srv/probe-b`, model: "deepseek-pro"},
	} {
		profile := deepSeekProfile()
		profile.Model = configuration.model
		executor := &fakeCopilotCLIProbeExecutor{execution: CopilotCLIProbeExecution{
			Outcome:     "start-failed",
			FailureCode: "",
			Evidence:    nil,
		}}
		probe := CopilotCLIProbe{
			Executable:       configuration.executable,
			WorkingDirectory: configuration.workingDirectory,
			Executor:         executor,
		}
		result, err := probe.Probe(context.Background(), profile, []byte("not-a-real-secret"), "agent-runner-cli", 1)
		if err != nil {
			t.Fatalf("a refused probe reports a result, not an error: %v", err)
		}
		if !executor.called {
			t.Fatal("the probe must reach the executor")
		}
		if result.Status != "unavailable" || result.Reason != "cli_unavailable" {
			t.Fatalf("a start-failed probe must be the cli_unavailable refusal, got %+v", result)
		}
		if len(result.Evidence) == 0 {
			t.Fatal("a refusal must still name its class")
		}
		results = append(results, result)
	}

	if !slices.Equal(results[0].Evidence, results[1].Evidence) {
		t.Fatalf("two refused configurations of one class must yield identical evidence:\n%v\n%v", results[0].Evidence, results[1].Evidence)
	}
	if len(results[0].Evidence) != 1 {
		t.Fatalf("a refusal must carry exactly its classification, got %v", results[0].Evidence)
	}
	joined := strings.Join(results[0].Evidence, " ")
	for _, echoed := range []string{"copilot-a", "copilot-b", "probe-a", "probe-b", "deepseek-flash", "deepseek-pro"} {
		if strings.Contains(joined, echoed) {
			t.Fatalf("the refusal evidence must not echo the refused value %q, got %v", echoed, results[0].Evidence)
		}
	}
}
