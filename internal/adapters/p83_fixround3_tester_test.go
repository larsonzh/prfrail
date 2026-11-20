package adapters

// P83 fix round 3: falsifiable counterexamples for the four blocking items the
// independent scan raised against the F2, F6 and F8 fixes.
//
// Every test here is written so that reverting the corresponding production
// guard turns it red:
//   - TestP83F2ReconfirmationIsTheLastObservationBeforeLaunch: moving the launch
//     intent write before ReconfirmLaunchAuthorization turns it red.
//   - TestP83F6BindingRefusalPrecedesSecretResolution: moving the receiver-binding
//     verification after the credential resolve turns it red.
//   - TestP83F6SwappedExecutablePrecedesSecretResolution: the same move, reached
//     through the pinned executable hash instead of a stub binder.
//   - TestP83F8ExecutorRefusesLoggingArgumentsBeforeTempDir: deleting the
//     logging-argument refusal turns it red.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/tickets"
)

// TestP83F2ReconfirmationIsTheLastObservationBeforeLaunch proves the F2 ordering
// on Windows: the authorization verdict is taken from the dispatch-time snapshot
// after every observable preparation step and before the launch intent exists, so
// a revocation that lands at the reconfirmation step still prevents the spawn.
func TestP83F2ReconfirmationIsTheLastObservationBeforeLaunch(t *testing.T) {
	dispatcher, admission, launcher, request := replayDispatcherFixture(t)
	grant := admission.AuthorizationLedger.Grants[0]
	revocation, err := chain.NewRevocationAuthorizationRecord(chain.AuthorizationRevocation{
		RecordID:          "p83-f2-revocation",
		Kind:              "revocation",
		IssuedAt:          "2026-09-11T08:00:00.000Z",
		IssuedBy:          evidence.Actor{Type: "operator", ID: "p83-operator"},
		AuthorizationID:   grant.Grant.AuthorizationID,
		AuthorizationHash: grant.RecordHash,
		ReasonCode:        "operator-request",
		StopDisposition:   "not-required",
	}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	revoked := chain.AuthorizationLedger{
		Grants:      []chain.AuthorizationRecord{grant},
		Revocations: []chain.AuthorizationRecord{revocation},
	}
	// The observations below are taken from inside the dispatcher, so they also
	// assert the ordering: at the moment the ledgers are reread the launch intent
	// must not exist yet, which is only true when the verdict is taken before the
	// single durable preparation step that precedes the spawn.
	requestID := dispatcher.RequestRecord.Request.RequestID
	var observed []string
	requireNoIntentYet := func(step string) {
		if _, found, err := dispatcher.Store.LaunchReceipt(requestID); err != nil || found {
			t.Fatalf("the launch intent must not exist before the %s reconfirmation observation: found=%v err=%v", step, found, err)
		}
	}
	dispatcher.LedgerSnapshot = func() (chain.AuthorizationLedger, error) {
		requireNoIntentYet("authorization")
		observed = append(observed, "authorization")
		return revoked, nil
	}
	dispatcher.CostLedgerSnapshot = func() (*tickets.CostLedger, error) {
		requireNoIntentYet("budget")
		observed = append(observed, "budget")
		return admission.CostLedger, nil
	}

	_, err = dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) || !strings.Contains(err.Error(), "revoked") {
		t.Fatalf("a revocation observed at the reconfirmation must refuse the launch, got %v", err)
	}
	if want := []string{"authorization", "budget"}; len(observed) != len(want) ||
		observed[0] != want[0] || observed[1] != want[1] {
		t.Fatalf("the reconfirmation must observe the authorization and budget ledgers in that order before the launch, got %v", observed)
	}
	if calls := launcher.callCount(); calls != 0 {
		t.Fatalf("a refused reconfirmation must never reach StartAgentRunnerProcess, got %d starts", calls)
	}
	if _, found, err := dispatcher.Store.LaunchReceipt(requestID); err != nil || found {
		t.Fatalf("a refused reconfirmation must publish no launch intent: found=%v err=%v", found, err)
	}
	if state, err := dispatcher.Store.State(requestID); err != nil || state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("R must stay retained after the refusal, got state=%v err=%v", state, err)
	}
}

// p83F6BindingRefusedProbe is a probe that fails its receiver binding and would
// resolve a credential-bearing probe if it were ever entered.
type p83F6BindingRefusedProbe struct {
	probeCalls int
}

func (probe *p83F6BindingRefusedProbe) VerifyReceiverBinding() error {
	return ErrCopilotCLIExecutableBinding
}

func (probe *p83F6BindingRefusedProbe) Probe(context.Context, AIProviderProfile, []byte, string, int) (AIProbeResult, error) {
	probe.probeCalls++
	return AIProbeResult{Status: "available", RequestsUsed: 1, Evidence: []string{queueHashTwo}}, nil
}

// p83F6SecretSpy counts credential resolutions so the ordering is observable.
type p83F6SecretSpy struct {
	resolveCalls int
}

func (spy *p83F6SecretSpy) Resolve(context.Context, string) ([]byte, error) {
	spy.resolveCalls++
	return []byte("not-a-real-secret"), nil
}

// TestP83F6BindingRefusalPrecedesSecretResolution proves the F6 ordering: a probe
// whose receiver binding cannot be proven is refused before the secret store is
// touched at all.
func TestP83F6BindingRefusalPrecedesSecretResolution(t *testing.T) {
	spy := &p83F6SecretSpy{}
	probe := &p83F6BindingRefusedProbe{}
	checker := availabilityChecker(probe, spy)
	record, err := checker.Check(context.Background(), deepSeekProfile(), "agent-runner-cli", 1)
	if err != nil {
		t.Fatal(err)
	}
	if spy.resolveCalls != 0 {
		t.Fatalf("the secret store must not be read before the receiver binding holds, resolutions=%d", spy.resolveCalls)
	}
	if probe.probeCalls != 0 {
		t.Fatalf("the probe must not run when its binding is refused, calls=%d", probe.probeCalls)
	}
	if record.Availability.Status == "available" {
		t.Fatalf("a refused binding must not read as available, got %+v", record.Availability)
	}
}

// TestP83F6SwappedExecutablePrecedesSecretResolution proves the same ordering for
// a real swap: the pinned executable no longer matches, so the availability check
// refuses before the credential for the approved profile is materialized.
func TestP83F6SwappedExecutablePrecedesSecretResolution(t *testing.T) {
	root := t.TempDir()
	executablePath := filepath.Join(root, "copilot-probe.exe")
	if err := os.WriteFile(executablePath, []byte("approved-content"), 0o700); err != nil {
		t.Fatal(err)
	}
	approved, err := CopilotCLIExecutableHash(executablePath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(executablePath, []byte("swapped-content!"), 0o700); err != nil {
		t.Fatal(err)
	}
	tempRoot := filepath.Join(root, "probe-temp")
	if err := os.MkdirAll(tempRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	probe := CopilotCLIProbe{
		Executable:       executablePath,
		WorkingDirectory: root,
		Executor:         ManagedCopilotCLIProbeExecutor{TemporaryDirectory: tempRoot, ExpectedExecutableHash: approved},
	}
	spy := &p83F6SecretSpy{}
	record, err := availabilityChecker(probe, spy).Check(context.Background(), deepSeekProfile(), "agent-runner-cli", 1)
	if err != nil {
		t.Fatal(err)
	}
	if spy.resolveCalls != 0 {
		t.Fatalf("a swapped executable must be detected before the secret is resolved, resolutions=%d", spy.resolveCalls)
	}
	if record.Availability.Status == "available" {
		t.Fatalf("a swapped executable must not read as available, got %+v", record.Availability)
	}
	if entries, err := os.ReadDir(tempRoot); err == nil && len(entries) != 0 {
		t.Fatalf("a binding refusal must leave no probe footprint, found %v", entries)
	}
}

// TestP83F8ExecutorRefusesLoggingArgumentsBeforeTempDir is the F8 counterexample
// over every spelling of a logging destination: the executor refuses before it
// creates any temporary directory, and the evidence it returns is a digest
// instead of the caller's raw argument.
func TestP83F8ExecutorRefusesLoggingArgumentsBeforeTempDir(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		refused  string
		class    string
		wantFail bool
	}{
		{"long flag", []string{"--log-dir"}, "--log-dir", "log-dir", true},
		{"long flag with value", []string{"--log-dir=C:/p83/logs"}, "--log-dir=C:/p83/logs", "log-dir", true},
		{"long flag separate value", []string{"--log-dir", "C:/p83/logs"}, "--log-dir", "log-dir", true},
		{"long log file", []string{"--log-file=C:/p83/out.log"}, "--log-file=C:/p83/out.log", "log-file", true},
		{"single dash logdir", []string{"-logdir"}, "-logdir", "logdir", true},
		{"windows slash log dir", []string{"/log-dir"}, "/log-dir", "log-dir", true},
		{"mixed case", []string{"--LogDir"}, "--LogDir", "logdir", true},
		{"benign prompt", []string{"-p", "prompt"}, "", "", false},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			refused, ok := copilotCLIProbeRefusedLoggingArg(testCase.args)
			if ok != testCase.wantFail {
				t.Fatalf("refusal decision for %v = %v, want %v", testCase.args, ok, testCase.wantFail)
			}
			if !testCase.wantFail {
				return
			}
			if refused != testCase.refused {
				t.Fatalf("refused argument = %q, want %q", refused, testCase.refused)
			}
			root := t.TempDir()
			executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: root}
			execution := executor.Execute(context.Background(), CopilotCLIProbeInvocation{
				Executable:       filepath.Join(root, "copilot.exe"),
				WorkingDirectory: root,
				Args:             testCase.args,
			})
			if execution.Outcome != "start-failed" {
				t.Fatalf("a logging destination must refuse before the spawn, got %+v", execution)
			}
			if len(execution.Evidence) != 1 {
				t.Fatalf("the refusal must carry exactly one evidence digest, got %v", execution.Evidence)
			}
			// T8b (P83 fix round): the evidence names the flag class, never the value,
			// so two refused values of one flag are indistinguishable in evidence.
			want := copilotCLIExecutorEvidence("logging-argument-refused:" + testCase.class)
			if execution.Evidence[0] != want {
				t.Fatalf("the evidence must be the digest of the refusal, got %q want %q", execution.Evidence[0], want)
			}
			if !strings.HasPrefix(execution.Evidence[0], "sha256:") || strings.Contains(strings.ToLower(execution.Evidence[0][len("sha256:"):]), "log") {
				t.Fatalf("the evidence must not leak the raw argument, got %q", execution.Evidence[0])
			}
			entries, err := os.ReadDir(root)
			if err != nil {
				t.Fatal(err)
			}
			for _, entry := range entries {
				if strings.Contains(entry.Name(), "proofrail-ai-probe-") {
					t.Fatalf("a refused logging argument must create no temporary directory, found %q", entry.Name())
				}
			}
		})
	}
}
