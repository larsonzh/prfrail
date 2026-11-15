package console

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/guard"
)

// TestSupervisedRunSandboxRefusesBeforePermit proves that the supervised
// fresh-start branch consults the sandbox check layer before it claims the
// one-shot probe permit: a Denied or Unknown verdict and an unavailable checker
// all refuse fail-closed, the operator-interaction ledger and the run-root
// marker are never created, and the quota stays claimable. A second invocation
// after the refusal is still a fresh start and is refused the same way.
//
// Deleting the check call point of the fresh-start branch makes every case
// fail: the run then claims the permit, writes the marker and parks in the
// waiting state with a nil error.
func TestSupervisedRunSandboxRefusesBeforePermit(t *testing.T) {
	checkFailure := errors.New("sandbox checker unavailable")
	cases := []struct {
		name     string
		verdict  guard.Verdict
		checkErr error
		want     error
	}{
		{
			name:    "denied",
			verdict: guard.Verdict{Status: guard.VerdictDenied, Reason: "environment does not comply"},
			want:    guard.ErrSandboxRejected,
		},
		{
			name:    "unknown",
			verdict: guard.Verdict{Status: guard.VerdictUnknown, Reason: "compliance unknown"},
			want:    guard.ErrSandboxRejected,
		},
		{
			name:     "unavailable",
			checkErr: checkFailure,
			want:     guard.ErrSandboxUnavailable,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			harness := newSupervisedHarness(t, []string{"complete"}, false, 1000)
			fake := &sandboxCheckRecorder{verdict: testCase.verdict, err: testCase.checkErr}
			installSandboxDispatch(t, fake)

			options := harness.options
			options.SandboxEnforced = true

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// A second invocation proves that the refusal leaves no state
			// behind: the run is still a fresh start and is refused the same
			// way.
			for invocation := 1; invocation <= 2; invocation++ {
				summary, err := ExecuteSupervisedRun(ctx, options)
				if err == nil || !errors.Is(err, testCase.want) {
					t.Fatalf("invocation %d: the dispatch must be refused with %v, got %v", invocation, testCase.want, err)
				}
				if !sandboxSummaryIsZero(summary) {
					t.Fatalf("invocation %d: a refused dispatch must return the zero summary, got %+v", invocation, summary)
				}
				if len(fake.reasons) != invocation {
					t.Fatalf("invocation %d: the sandbox check must run exactly once per invocation, got %v", invocation, fake.reasons)
				}
				if fake.reasons[invocation-1] != "supervised-runtime-dispatch" {
					t.Fatalf("invocation %d: the sandbox check must carry the dispatch label, got %q", invocation, fake.reasons[invocation-1])
				}
			}

			if _, statErr := os.Stat(SupervisedInteractionsLedgerPath(harness.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatalf("a sandbox refusal must not create the operator-interaction ledger, stat error = %v", statErr)
			}
			if _, statErr := os.Stat(runEventLogPath(harness.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatalf("a sandbox refusal must not create the run-root marker, stat error = %v", statErr)
			}

			// The check point precedes the probe permit: the one-shot quota is
			// still claimable after both refusals. This is the last assertion
			// because claiming it consumes the quota for good.
			release, permitErr := requireProbePermit(ctx, options.ProbePermit)
			if permitErr != nil {
				t.Fatalf("a sandbox refusal must not burn the one-shot probe quota: %v", permitErr)
			}
			release()
		})
	}
}
