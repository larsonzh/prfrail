package console

// P85 T3 (finding F6) - the install source counterexamples at the command boundary.
//
// `ai check` admits an agent-runner-cli probe only after it has PROVEN that the
// probe executable's installation source is immutable: a stable absolute path whose
// install directory denies this user write, delete and permission changes
// (Windows DACL) or whose every path component is unwritable and foreign-owned
// (Linux). The proof is memoised into the store-local carrier
// `ai-install-source-proof.json` under the chain-config directory, but the verdict
// always comes from a fresh observation - so a missing, unreadable or stale carrier
// can only ever REFUSE, never admit.
//
// The counterexamples below assert the three refusals named by the design (§3.4)
// plus the positive control that the carrier is published and then converges. Each
// one asserts the full consequence, not just a non-zero exit: no credential is
// resolved, the availability checker is never entered, no availability record is
// published, and the carrier on disk is either absent (nothing was proven) or
// byte-identical to what the test wrote (a published proof is never repaired).

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85T3ImmutableProbeExecutable returns an executable whose installation source the
// live check can actually prove on this host, or skips the caller.
//
// The console gate is hardwired to the platform observation (there is no injectable
// verifier), so a positive fixtures has only two honest choices: a binary shipped by
// the OS in an unwritable directory, or a skip. A stub written into t.TempDir() is
// exactly the shape T3 refuses - that is the finding, not a test artefact.
func p85T3ImmutableProbeExecutable(t *testing.T) string {
	t.Helper()
	var candidate string
	switch runtime.GOOS {
	case "windows":
		root := os.Getenv("SystemRoot")
		if root == "" {
			root = `C:\Windows`
		}
		candidate = filepath.Join(root, "System32", "cmd.exe")
	case "linux":
		candidate = "/bin/sh"
	default:
		t.Skipf("the install source proof has a platform rule only for windows and linux, not %s", runtime.GOOS)
	}
	if _, err := adapters.VerifyAIInstallSource(candidate); err != nil {
		// euid 0 cannot prove immutability by design (see ai_install_proof_unix.go),
		// so a root CI leg must skip the positive path instead of failing it.
		t.Skipf("the fixture executable %s is not provably immutable for the executing user on this host (%v); this positive path needs a native non-root Ubuntu run", candidate, err)
	}
	return candidate
}

// p85T3CarrierPath is the carrier location the command derives on its own: the
// trusted root of the chain configuration. A test never passes this as a flag,
// because the command does not accept one - that is the self-approval hole T3
// closes.
func p85T3CarrierPath(root string) string {
	return filepath.Join(root, adapters.AIInstallSourceProofFileName)
}

// p85T3CheckOutcome is what one `ai check` invocation did, from the outside.
type p85T3CheckOutcome struct {
	code         int
	output       string
	checkerCalls int
	resolveCalls int
	outputPath   string
}

// p85T3RunCheck runs `ai check` against the shared configured chain with an
// instrumented secret store and an availability checker that never spawns anything.
func p85T3RunCheck(t *testing.T, root, configPath, executablePath, outName string) p85T3CheckOutcome {
	t.Helper()
	secrets := &recordingProbeSecrets{credential: []byte("p85-t3-probe-credential")}
	cli := newTestCLI(root)
	cli.secretManager = secrets
	entries := 0
	cli.checkAIAvailability = func(_ context.Context, profile adapters.AIProviderProfile, channel string, _ int, _ aiProbeReceiver, _ string) (adapters.AIAvailabilityRecord, error) {
		entries++
		if _, resolveErr := cli.secretManager.Resolve(context.Background(), profile.SecretRef); resolveErr != nil {
			return adapters.AIAvailabilityRecord{}, resolveErr
		}
		return adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
			ProbeID:           "p85-t3-probe",
			ProfileID:         profile.ProfileID,
			ProfileConfigHash: mustAIProfileHash(t, profile),
			Channel:           channel,
			ProbedAt:          "2026-09-11T08:00:00.000Z",
			Status:            "available",
			RequestsUsed:      1,
			Evidence:          []string{"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		})
	}
	outputPath := filepath.Join(root, outName)
	code, stdout, stderr := runCLI(t, cli,
		"ai", "check", "--chain", configPath, "--channel", "agent-runner-cli",
		"--copilot", executablePath, "--max-requests", "1", "--out", outputPath, "--json",
	)
	return p85T3CheckOutcome{
		code:         code,
		output:       stdout + stderr,
		checkerCalls: entries,
		resolveCalls: secrets.resolveCalls,
		outputPath:   outputPath,
	}
}

// p85T3AssertRefusedBeforeAnyCredential pins the whole consequence of a refusal.
func p85T3AssertRefusedBeforeAnyCredential(t *testing.T, outcome p85T3CheckOutcome, carrierPath string) {
	t.Helper()
	if outcome.code != exitFailure {
		t.Fatalf("an unproven install source must fail closed, got code=%d output=%q", outcome.code, outcome.output)
	}
	if !strings.Contains(outcome.output, "install source") {
		t.Fatalf("the install source gate must be the refusing guard, got %q", outcome.output)
	}
	if outcome.checkerCalls != 0 {
		t.Fatalf("the availability checker must not be entered, entries=%d", outcome.checkerCalls)
	}
	if outcome.resolveCalls != 0 {
		t.Fatalf("no credential may be resolved, resolutions=%d", outcome.resolveCalls)
	}
	if _, statErr := os.Lstat(outcome.outputPath); !os.IsNotExist(statErr) {
		t.Fatalf("no availability record may be published, stat err=%v", statErr)
	}
	if _, statErr := os.Lstat(carrierPath); statErr != nil && !os.IsNotExist(statErr) {
		t.Fatalf("unexpected carrier state: %v", statErr)
	}
}

// TestP85T3RefusesUnprovableInstallSourceWithoutPublishingCarrier is the "carrier
// missing and the live observation cannot replace it" counterexample: the probe
// executable sits in a directory this user can write, so nothing can be proven and
// nothing may be memoised.
//
// Mutation check: dropping the ACL capability probe (Windows) or the mode/owner
// check (Linux) makes the observation succeed, the checker run and this red.
func TestP85T3RefusesUnprovableInstallSourceWithoutPublishingCarrier(t *testing.T) {
	root := t.TempDir()
	configPath := writeAIProbeTestChain(t, root)
	executablePath := filepath.Join(root, "copilot.exe")
	// A writable directory this user owns: writable on every platform T3 supports.
	if err := os.WriteFile(executablePath, []byte("p85-t3-unprovable-probe"), 0o700); err != nil {
		t.Fatal(err)
	}
	// The operator grant exists, so the refusal below is attributable to the install
	// source gate and not to an unapproved receiver.
	writeAIProbeApprovedLedger(t, root, executablePath)
	carrierPath := p85T3CarrierPath(root)

	outcome := p85T3RunCheck(t, root, configPath, executablePath, "p85-t3-unprovable.json")
	p85T3AssertRefusedBeforeAnyCredential(t, outcome, carrierPath)
	if _, statErr := os.Lstat(carrierPath); !os.IsNotExist(statErr) {
		t.Fatalf("an unproven source must not publish a carrier, stat err=%v", statErr)
	}
}

// TestP85T3RefusesMalformedCarrierInsteadOfRecomputing is the "carrier unreadable"
// counterexample: a carrier that cannot be decoded is a refusal, never a cue to
// recompute and continue. The carrier bytes must survive untouched, so an operator
// can still see what was rejected.
func TestP85T3RefusesMalformedCarrierInsteadOfRecomputing(t *testing.T) {
	root := t.TempDir()
	configPath := writeAIProbeTestChain(t, root)
	executablePath := p85T3ImmutableProbeExecutable(t)
	writeAIProbeApprovedLedger(t, root, executablePath)

	carrierPath := p85T3CarrierPath(root)
	const malformed = `{"schemaVersion":1,"executablePath":`
	if err := os.WriteFile(carrierPath, []byte(malformed), 0o600); err != nil {
		t.Fatal(err)
	}

	outcome := p85T3RunCheck(t, root, configPath, executablePath, "p85-t3-malformed.json")
	p85T3AssertRefusedBeforeAnyCredential(t, outcome, carrierPath)
	wire, err := os.ReadFile(carrierPath)
	if err != nil {
		t.Fatalf("the rejected carrier must stay readable: %v", err)
	}
	if string(wire) != malformed {
		t.Fatalf("a malformed carrier must never be rewritten, got %q", string(wire))
	}
}

// TestP85T3RefusesCarrierThatNoLongerDescribesTheSource is the "well-formed but
// mismatched" counterexample: the carrier parses and validates, but it describes a
// different object (a different digest and identity), so the live observation does
// not match it. The refusal must happen without touching the file.
func TestP85T3RefusesCarrierThatNoLongerDescribesTheSource(t *testing.T) {
	root := t.TempDir()
	configPath := writeAIProbeTestChain(t, root)
	executablePath := p85T3ImmutableProbeExecutable(t)
	writeAIProbeApprovedLedger(t, root, executablePath)

	facts, err := adapters.VerifyAIInstallSource(executablePath)
	if err != nil {
		t.Fatalf("the fixture executable must be provable: %v", err)
	}
	stale := facts
	stale.ExecutableHash = evidence.Digest("", []byte("p85-t3-a-different-approved-object"))
	proof := adapters.NewAIInstallSourceProof(stale, time.Date(2026, 9, 8, 10, 0, 0, 0, time.UTC))
	carrierPath := p85T3CarrierPath(root)
	payload, err := adapters.EncodeAIInstallSourceProof(proof)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(carrierPath, payload, 0o600); err != nil {
		t.Fatal(err)
	}

	outcome := p85T3RunCheck(t, root, configPath, executablePath, "p85-t3-stale.json")
	p85T3AssertRefusedBeforeAnyCredential(t, outcome, carrierPath)
	published, err := os.ReadFile(carrierPath)
	if err != nil {
		t.Fatalf("the stale carrier must stay readable: %v", err)
	}
	if string(published) != string(payload) {
		t.Fatalf("a carrier that no longer describes the source must never be rewritten")
	}
}

// TestP85T3PublishesTheCarrierAndConverges is the positive control: a provable
// source publishes the carrier once, and a second run with the same source converges
// on the published bytes instead of replacing them. Without this control the three
// refusals above could be satisfied by refusing everything.
func TestP85T3PublishesTheCarrierAndConverges(t *testing.T) {
	root := t.TempDir()
	configPath := writeAIProbeTestChain(t, root)
	executablePath := p85T3ImmutableProbeExecutable(t)
	writeAIProbeApprovedLedger(t, root, executablePath)
	carrierPath := p85T3CarrierPath(root)

	first := p85T3RunCheck(t, root, configPath, executablePath, "p85-t3-available-1.json")
	if first.code != exitSuccess || first.checkerCalls != 1 || first.resolveCalls != 1 {
		t.Fatalf("a provable source must be admitted, got code=%d checker=%d resolve=%d output=%q",
			first.code, first.checkerCalls, first.resolveCalls, first.output)
	}
	published, err := os.ReadFile(carrierPath)
	if err != nil {
		t.Fatalf("an admitted probe must publish the carrier: %v", err)
	}
	decoded, err := adapters.DecodeAIInstallSourceProof(published)
	if err != nil {
		t.Fatalf("the published carrier must decode: %v", err)
	}
	facts, err := adapters.VerifyAIInstallSource(executablePath)
	if err != nil {
		t.Fatal(err)
	}
	if !facts.Matches(decoded) {
		t.Fatalf("the published carrier must describe the live source, got %+v", decoded)
	}

	second := p85T3RunCheck(t, root, configPath, executablePath, "p85-t3-available-2.json")
	if second.code != exitSuccess || second.checkerCalls != 1 {
		t.Fatalf("a converged carrier must still admit the probe, got code=%d checker=%d output=%q",
			second.code, second.checkerCalls, second.output)
	}
	again, err := os.ReadFile(carrierPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(again) != string(published) {
		t.Fatalf("a converged carrier must be byte-identical, got %q and %q", string(published), string(again))
	}
}
