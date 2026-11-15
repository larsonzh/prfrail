package console

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/guard"
)

// sandboxCheckRecorder is the white-box double of the dispatch-path sandbox
// check. It answers one scripted verdict and records every CheckSpec it saw, so
// a test can prove both the refusal and the position of the call point.
type sandboxCheckRecorder struct {
	verdict guard.Verdict
	err     error
	reasons []string
}

func (recorder *sandboxCheckRecorder) Check(_ context.Context, spec guard.CheckSpec) (guard.Verdict, error) {
	recorder.reasons = append(recorder.reasons, spec.Reason)
	return recorder.verdict, recorder.err
}

// installSandboxDispatch replaces the package-level check implementation for
// the duration of one test. The replacement is restored through t.Cleanup, and
// no test that installs it may run in parallel (the seam is a package
// variable).
func installSandboxDispatch(t *testing.T, fake guard.Sandbox) {
	t.Helper()
	previous := sandboxForDispatch
	sandboxForDispatch = fake
	t.Cleanup(func() { sandboxForDispatch = previous })
}

// sandboxSummaryIsZero reports whether a summary carries no run at all. It is
// the comparability-safe form of `summary == RunSummary{}` (the summary holds
// slices).
func sandboxSummaryIsZero(summary RunSummary) bool {
	return summary.RunID == "" &&
		summary.ChainState == "" &&
		summary.Sequence == 0 &&
		summary.LastEventHash == "" &&
		len(summary.TaskStateCounts) == 0 &&
		len(summary.StepStateCounts) == 0 &&
		summary.EventLogPath == ""
}

// newAssemblySandboxStub wires a single-task fixture to the full in-package
// port double. Every chain port is stubbed, so a run that is NOT refused by the
// sandbox check can only reach recording doubles: it never spawns a process.
func newAssemblySandboxStub(t *testing.T, fixture *assemblyFixture) (*agentRunnerAssemblyStub, *assemblyDriverStub) {
	t.Helper()
	record := fixture.options.RequestRecord
	driver := &assemblyDriverStub{
		t:          t,
		records:    map[string]adapters.AgentRunnerRequestRecord{record.Request.TaskID: record},
		facts:      map[string]chain.AgentRunnerFrozenFacts{record.Request.TaskID: multiTaskFacts(record.Request.TaskID)},
		failByTask: map[string]error{},
	}
	return &agentRunnerAssemblyStub{
		Dispatcher:      driver,
		Admission:       driver,
		Acceptance:      driver,
		Reviewer:        driver,
		Publisher:       driver,
		Postflight:      driver,
		Waiter:          driver,
		PublishTerminal: driver.publishTerminal,
		LoadTerminal:    driver.loadTerminal,
	}, driver
}

// TestSandboxEnforcementFromEnv pins the S1 trigger semantics of CONTRACTS.md
// 7.2 (e): unset, empty and "off" mean off; "enforce" means on; anything else
// is refused so a typo can never silently disable an enforcement request.
func TestSandboxEnforcementFromEnv(t *testing.T) {
	cases := []struct {
		name    string
		set     bool
		value   string
		want    bool
		wantErr bool
	}{
		{name: "unset", set: false},
		{name: "empty", set: true, value: ""},
		{name: "off", set: true, value: "off"},
		{name: "enforce", set: true, value: "enforce", want: true},
		{name: "unknown value", set: true, value: "on", wantErr: true},
		{name: "wrong case", set: true, value: "Enforce", wantErr: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.set {
				t.Setenv(sandboxEnforcementEnvVar, testCase.value)
			} else {
				t.Setenv(sandboxEnforcementEnvVar, "")
				os.Unsetenv(sandboxEnforcementEnvVar)
			}
			enforced, err := sandboxEnforcementFromEnv()
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("value %q must be refused, got enforced=%t", testCase.value, enforced)
				}
				if enforced {
					t.Fatal("a refused value must never enable the check")
				}
				return
			}
			if err != nil {
				t.Fatalf("value %q must be accepted: %v", testCase.value, err)
			}
			if enforced != testCase.want {
				t.Fatalf("enforced = %t, want %t for value %q", enforced, testCase.want, testCase.value)
			}
		})
	}
}

// TestAgentRunnerRunSandboxRefusesBeforePermit proves that the AgentRunner
// assembly entry consults the sandbox check layer before it claims the one-shot
// probe permit: a Denied or Unknown verdict and an unavailable checker all
// refuse fail-closed, the run-root marker is never created, no AgentRunner port
// is reached, and the quota stays claimable.
//
// Deleting the check call point of ExecuteAgentRunnerRun makes every case fail:
// the run then claims the permit, creates the step-2 marker, reaches the
// recording dispatcher and returns something other than a sandbox refusal.
func TestAgentRunnerRunSandboxRefusesBeforePermit(t *testing.T) {
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
			fixture := newAssemblyFixture(t, assemblyKnobs{})
			doubles, driver := newAssemblySandboxStub(t, fixture)
			fake := &sandboxCheckRecorder{verdict: testCase.verdict, err: testCase.checkErr}
			installSandboxDispatch(t, fake)

			options := fixture.options
			options.SandboxEnforced = true
			options.testStub = doubles

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			summary, err := ExecuteAgentRunnerRun(ctx, options)

			if err == nil || !errors.Is(err, testCase.want) {
				t.Fatalf("the dispatch must be refused with %v, got %v", testCase.want, err)
			}
			if !sandboxSummaryIsZero(summary) {
				t.Fatalf("a refused dispatch must return the zero summary, got %+v", summary)
			}
			if len(fake.reasons) != 1 {
				t.Fatalf("the sandbox check must run exactly once, got %d call(s): %v", len(fake.reasons), fake.reasons)
			}
			if fake.reasons[0] != "agent-runner-dispatch" {
				t.Fatalf("the sandbox check must carry the dispatch label, got %q", fake.reasons[0])
			}
			if len(driver.dispatched) != 0 {
				t.Fatalf("a refused dispatch must reach no AgentRunner port, got %v", driver.dispatched)
			}
			if _, statErr := os.Stat(runEventLogPath(fixture.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatalf("a sandbox refusal must not create the run-root marker, stat error = %v", statErr)
			}

			// The check point precedes the probe permit: the one-shot quota is
			// still claimable after the refusal. This is the last assertion
			// because claiming it consumes the quota for good.
			release, permitErr := requireProbePermit(ctx, options.ProbePermit)
			if permitErr != nil {
				t.Fatalf("a sandbox refusal must not burn the one-shot probe quota: %v", permitErr)
			}
			release()
		})
	}
}
