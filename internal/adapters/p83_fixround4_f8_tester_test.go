package adapters

// P83 FIX ROUND 4 (F8) - observation evidence must not be derived from the raw
// child streams, and the logging-destination screen must not mistake data for a
// flag.
//
// CONTRACTS.md:214/216 forbids API keys and tokens in evidence, and the probe is
// the one place where unredacted child output exists at all. Round 4 replaced the
// "digest of stdout/stderr/usage" evidence with a digest over the validated
// structure (copilotCLIProbeObservationEvidence). The counterexample below plants
// a canary in the child's own output and asserts that neither the canary nor
// anything derived from it reaches evidence - and, because two runs that differ
// only in the canary must produce identical evidence, that no digest of the raw
// stream can be reintroduced without turning this red.
//
// Mutation checks:
//   - re-adding any digest of stdout/stderr/usage to execution.Evidence makes the
//     two runs disagree and turns TestP83F8ObservationEvidenceCarriesNoRawChildOutput red;
//   - deleting the copilotCLIProbeObservationEvidence call leaves no evidence entry
//     at all, which the same test also rejects;
//   - dropping the "-"/"/" prefix requirement in copilotCLIProbeRefusedLoggingArg
//     makes a bare "log" value refuse and turns
//     TestP83F8LoggingScreenDistinguishesFlagsFromData red.

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
)

// p83Fix4CanaryCommand returns the platform shell and the argv that echoes the
// canary. The counterexample must run on every CI leg: a Windows-only child would
// make the Linux leg fail with start-failed instead of judging the evidence.
//
// P85 (T3, finding F6): the executor re-verifies the child's INSTALLATION SOURCE
// immediately before the spawn, so the fixture may no longer name a bare "cmd.exe"
// (no stable absolute path) and may not live in a directory this user can write.
// The shell is therefore addressed by its absolute system path, and the test skips
// - never passes - when the host cannot prove that source immutable (a Linux run
// as uid 0 cannot prove it by design; see ai_install_proof_unix.go).
func p83Fix4CanaryCommand(t *testing.T, canary string, redirectStderr bool) (string, []string) {
	t.Helper()
	var shell string
	switch runtime.GOOS {
	case "windows":
		root := os.Getenv("SystemRoot")
		if root == "" {
			root = `C:\Windows`
		}
		shell = filepath.Join(root, "System32", "cmd.exe")
	case "linux":
		shell = "/bin/sh"
	default:
		t.Skipf("the install source proof has a platform rule only for windows and linux, not %s", runtime.GOOS)
	}
	if _, err := VerifyAIInstallSource(shell); err != nil {
		t.Skipf("the canary shell %s is not provably immutable for the executing user on this host (%v); this leg needs a native non-root Ubuntu run or a privileged fixture", shell, err)
	}
	if redirectStderr {
		return shell, []string{"-c", "echo " + canary + " 1>&2"}
	}
	return shell, []string{"-c", "echo " + canary}
}

// p83Fix4RunProbeWithCanary starts a real child (cmd.exe) whose output carries the
// canary, so the evidence is produced by the production execution path rather than
// by a hand-built struct.
func p83Fix4RunProbeWithCanary(t *testing.T, canary string, redirectStderr bool) CopilotCLIProbeExecution {
	t.Helper()
	root := t.TempDir()
	executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: root}
	executable, args := p83Fix4CanaryCommand(t, canary, redirectStderr)
	execution := executor.Execute(context.Background(), CopilotCLIProbeInvocation{
		Executable:       executable,
		WorkingDirectory: root,
		Args:             args,
	})
	if execution.Outcome != "exited" || execution.ExitCode != 0 {
		t.Fatalf("the canary child must exit cleanly so the observation covers a completed run, got outcome=%q exit=%d evidence=%v",
			execution.Outcome, execution.ExitCode, execution.Evidence)
	}
	return execution
}

// TestP83F8ObservationEvidenceCarriesNoRawChildOutput is the F8 counterexample for
// the observation evidence: facts, never bytes.
func TestP83F8ObservationEvidenceCarriesNoRawChildOutput(t *testing.T) {
	const (
		alphaCanary = "P83CANARYALPHA"
		betaCanary  = "P83CANARYBETA"
	)
	alpha := p83Fix4RunProbeWithCanary(t, alphaCanary, false)
	beta := p83Fix4RunProbeWithCanary(t, betaCanary, false)

	// The observation evidence is still reported: a run that records nothing would
	// satisfy a "no canary" assertion vacuously.
	if len(alpha.Evidence) != 1 {
		t.Fatalf("a completed probe must carry exactly the observation evidence, got %v", alpha.Evidence)
	}
	if want := copilotCLIProbeObservationEvidence(alpha); alpha.Evidence[0] != want {
		t.Fatalf("the evidence must be the observation digest, got %q want %q", alpha.Evidence[0], want)
	}
	if !strings.HasPrefix(alpha.Evidence[0], "sha256:") {
		t.Fatalf("the observation evidence must be a digest, got %q", alpha.Evidence[0])
	}
	// Outcome/exit-class facts survive the redaction.
	if alpha.Outcome != "exited" || alpha.ExitCode != 0 {
		t.Fatalf("the observation must still report the exit class, got outcome=%q exit=%d", alpha.Outcome, alpha.ExitCode)
	}

	joined := strings.Join(alpha.Evidence, "\n") + "\n" + strings.Join(beta.Evidence, "\n")
	for _, canary := range []string{alphaCanary, betaCanary, strings.ToLower(alphaCanary), strings.ToLower(betaCanary)} {
		if strings.Contains(joined, canary) {
			t.Fatalf("evidence must not carry the child output %q, got %q", canary, joined)
		}
	}
	// The decisive assertion: two runs that differ only in what the child printed
	// must produce byte-identical evidence. Any digest over the raw stream would
	// differ here (the streams differ, and so does the temporary usage path).
	if !slices.Equal(alpha.Evidence, beta.Evidence) {
		t.Fatalf("evidence must not depend on the raw child output, got %v and %v", alpha.Evidence, beta.Evidence)
	}
}

// TestP83F8ObservationEvidenceIgnoresChildStderr is the stderr half of the same
// finding: the second stream is also unredacted, so a canary printed there must
// change nothing either.
func TestP83F8ObservationEvidenceIgnoresChildStderr(t *testing.T) {
	const canary = "P83STDERRCANARY"
	execution := p83Fix4RunProbeWithCanary(t, canary, true)
	for _, entry := range execution.Evidence {
		if strings.Contains(entry, canary) {
			t.Fatalf("evidence must not carry the child stderr %q, got %q", canary, entry)
		}
	}
	if len(execution.Evidence) != 1 || execution.Evidence[0] != copilotCLIProbeObservationEvidence(execution) {
		t.Fatalf("a stderr canary must leave the observation evidence unchanged, got %v", execution.Evidence)
	}
}

// TestP83F8LoggingScreenDistinguishesFlagsFromData pins the round-4 refinement of
// the F8 screen: only a token that starts like a flag can name a logging
// destination, so "-p log" (a prompt whose value is the word "log") and a bare
// positional "log" stay legal while every flag spelling is refused.
func TestP83F8LoggingScreenDistinguishesFlagsFromData(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		refused  string
		wantFail bool
	}{
		{"prompt-flag-kept", []string{"-p"}, "", false},
		{"prompt-value-log-kept", []string{"-p", "log"}, "", false},
		{"bare-positional-log-kept", []string{"log"}, "", false},
		{"bare-positional-logdir-kept", []string{"log-dir"}, "", false},
		{"long-flag-refused", []string{"--log-dir"}, "--log-dir", true},
		{"long-flag-with-value-refused", []string{"--log-dir=X"}, "--log-dir=X", true},
		{"single-dash-refused", []string{"-logdir"}, "-logdir", true},
		{"windows-slash-refused", []string{"/log-dir"}, "/log-dir", true},
		{"mixed-case-refused", []string{"--LogDir"}, "--LogDir", true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			refused, ok := copilotCLIProbeRefusedLoggingArg(testCase.args)
			if ok != testCase.wantFail {
				t.Fatalf("refusal decision for %v = %v, want %v", testCase.args, ok, testCase.wantFail)
			}
			if refused != testCase.refused {
				t.Fatalf("refused argument for %v = %q, want %q", testCase.args, refused, testCase.refused)
			}
		})
	}
}

// TestP83F8BenignProbeArgumentsStillReachTheSpawn is the acceptance half: a probe
// whose prompt value is literally "log" is not refused by the executor, so the
// screen above cannot be satisfied by refusing every argument.
func TestP83F8BenignProbeArgumentsStillReachTheSpawn(t *testing.T) {
	root := t.TempDir()
	executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: root}
	executable, args := p83Fix4CanaryCommand(t, "log", false)
	execution := executor.Execute(context.Background(), CopilotCLIProbeInvocation{
		Executable:       executable,
		WorkingDirectory: root,
		Args:             args,
	})
	if execution.Outcome != "exited" {
		t.Fatalf("a benign probe whose value is the word \"log\" must still run, got %+v", execution)
	}
	want := copilotCLIExecutorEvidence("logging-argument-refused:log")
	for _, entry := range execution.Evidence {
		if entry == want {
			t.Fatalf("a value token must not be read as a logging flag, got %v", execution.Evidence)
		}
	}
	if _, err := os.Stat(root); err != nil {
		t.Fatalf("the probe temporary root must remain usable: %v", err)
	}
}
