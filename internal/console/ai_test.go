package console

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// writeAIProbeApprovedLedger writes the operator grant the F6 gate requires, next
// to the chain configuration (its default ledger location). The grant names the
// receiver the configured profile targets (https://api.deepseek.com:443), the
// probe tool id and the digest of the executable under test.
func writeAIProbeApprovedLedger(t *testing.T, root, executablePath string) string {
	t.Helper()
	executableHash, err := adapters.CopilotCLIExecutableHash(executablePath)
	if err != nil {
		t.Fatalf("hash the approved probe executable: %v", err)
	}
	// T3 (finding F6): the approval must cover the FULL receiver identity, so the
	// grant evidence carries the profile hash (provider/baseUrl/model/secretRef),
	// the receiver binding digest (profile + channel + executable) and the
	// executable digest.
	profile, found := configuredAIProfile(configuredAIChain(), "agent-runner-cli")
	if !found {
		t.Fatal("the shared AI chain must configure the agent-runner-cli channel")
	}
	profileHash, err := adapters.AIProviderProfileHash(profile)
	if err != nil {
		t.Fatalf("hash the approved profile: %v", err)
	}
	binding := evidence.Digest(aiProbeReceiverBindingDomain, []byte(profileHash+"\n"+"agent-runner-cli"+"\n"+executableHash))
	// newTestCLI pins the clock at 2026-09-08T10:30:00Z, so the grant is issued
	// before it and expires after it.
	issued := time.Date(2026, 9, 8, 10, 0, 0, 0, time.UTC)
	record, err := chain.NewGrantAuthorizationRecord(chain.AuthorizationGrant{
		RecordID:        "ai-probe-grant-1",
		Kind:            "grant",
		IssuedAt:        issued.Format(timestampLayout),
		IssuedBy:        evidence.Actor{Type: "operator", ID: "ai-probe-operator"},
		AuthorizationID: "ai-probe-authorization-1",
		RunID:           "ai-probe-run",
		RunManifestHash: evidence.Digest("", []byte("ai-probe-manifest")),
		PolicyHash:      evidence.Digest("", []byte("ai-probe-policy")),
		Scope: chain.AuthorizationScope{
			ToolIDs: []string{aiProbeToolID},
			Network: []chain.AuthorizationNetworkScope{
				{Protocol: "https", Host: "api.deepseek.com", Port: 443, Purpose: "probe-provider"},
			},
			Budget: chain.AuthorizationBudgetLimit{
				ModelCalls: 1, Tokens: 1000, WallClockMs: 600000, Attempts: 1, Currency: "USD",
			},
		},
		ExpiresAt: issued.Add(time.Hour).Format(timestampLayout),
		Evidence:  []string{executableHash, profileHash, binding},
	}, nil, nil)
	if err != nil {
		t.Fatalf("build the operator grant: %v", err)
	}
	encoded, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	ledgerPath := filepath.Join(root, defaultApprovalsLedgerName)
	if err := os.WriteFile(ledgerPath, append(encoded, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	return ledgerPath
}

func writeAIAvailabilityFixture(t *testing.T, root, name string, record adapters.AIAvailabilityRecord) string {
	t.Helper()
	path := filepath.Join(root, name)
	if err := adapters.WriteAIAvailabilityRecord(path, record); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAICheckUsesConfiguredChannelAndReturnsAvailabilityRecord(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, DefaultChainConfigName)
	// F6: the probe receiver binding hashes the executable, and T3 additionally
	// proves its installation source immutable, so the target is a real system
	// binary instead of a stub written into a writable temporary directory.
	executablePath := p85T3ImmutableProbeExecutable(t)
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	// F6: a probe is only admitted by an active operator grant in the ledger.
	writeAIProbeApprovedLedger(t, root, executablePath)

	called := false
	cli := newTestCLI(root)
	cli.checkAIAvailability = func(_ context.Context, profile adapters.AIProviderProfile, channel string, maximumRequests int, receiver aiProbeReceiver, workingDirectory string) (adapters.AIAvailabilityRecord, error) {
		called = true
		if profile.ProfileID != "deepseek-anthropic" || channel != "agent-runner-cli" || maximumRequests != 1 {
			t.Fatalf("unexpected binding: profile=%s channel=%s requests=%d", profile.ProfileID, channel, maximumRequests)
		}
		if receiver.Executable != executablePath || workingDirectory != root {
			t.Fatalf("unexpected execution target: executable=%q cwd=%q", receiver.Executable, workingDirectory)
		}
		reason := "credential_rejected"
		return adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
			ProbeID:           "probe-cli-test",
			ProfileID:         profile.ProfileID,
			ProfileConfigHash: mustAIProfileHash(t, profile),
			Channel:           channel,
			ProbedAt:          "2026-09-11T08:00:00.000Z",
			Status:            "unavailable",
			Reason:            &reason,
			RequestsUsed:      1,
			Evidence:          []string{"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		})
	}

	code, stdout, stderr := runCLI(t, cli,
		"ai", "check",
		"--chain", configPath,
		"--channel", "agent-runner-cli",
		"--copilot", executablePath,
		"--max-requests", "1",
		"--out", filepath.Join(root, "availability.json"),
		"--json",
	)
	if code != exitFailure || stderr != "" || !called {
		t.Fatalf("unexpected command result: code=%d called=%t stdout=%q stderr=%q", code, called, stdout, stderr)
	}
	response := decodeResponse(t, stdout)
	if response.OK || response.ExitCode != exitFailure {
		t.Fatalf("unexpected response: %+v", response)
	}
	var record adapters.AIAvailabilityRecord
	if err := json.Unmarshal(response.Data, &record); err != nil {
		t.Fatal(err)
	}
	if record.Availability.Status != "unavailable" || record.Availability.Reason == nil || *record.Availability.Reason != "credential_rejected" {
		t.Fatalf("availability record was lost: %+v", record)
	}
	persisted, err := os.ReadFile(filepath.Join(root, "availability.json"))
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := adapters.DecodeAIAvailabilityRecord(persisted)
	if err != nil || decoded.RecordHash != record.RecordHash {
		t.Fatalf("availability record was not persisted: decoded=%+v err=%v", decoded, err)
	}
}

func TestAICheckRejectsUnconfiguredOrUnsupportedChannelWithoutProbe(t *testing.T) {
	root := t.TempDir()
	config := configuredAIChain()
	config.AI.Channels = config.AI.Channels[:1]
	configPath := filepath.Join(root, DefaultChainConfigName)
	encoded, err := EncodeChainConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	cli := newTestCLI(root)
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, aiProbeReceiver, string) (adapters.AIAvailabilityRecord, error) {
		t.Fatal("probe must not run for another or unconfigured channel")
		return adapters.AIAvailabilityRecord{}, nil
	}

	for _, channel := range []string{"sessionbridge-silent", "sessionbridge-visible"} {
		code, _, _ := runCLI(t, cli,
			"ai", "check", "--chain", configPath, "--channel", channel,
			"--copilot", filepath.Join(root, "copilot.exe"), "--max-requests", "1",
			"--out", filepath.Join(root, channel+".json"), "--json",
		)
		if code != exitFailure {
			t.Fatalf("channel %s did not fail closed: code=%d", channel, code)
		}
	}
}

func TestAICheckReturnsSuccessOnlyForAvailableRecord(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, DefaultChainConfigName)
	// F6 + T3: the receiver binding requires a hashable executable, and the
	// install source gate requires one this user cannot replace.
	executablePath := p85T3ImmutableProbeExecutable(t)
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	// F6: the operator grant is what admits the receiver, so the success path needs
	// one on disk.
	writeAIProbeApprovedLedger(t, root, executablePath)
	cli := newTestCLI(root)
	cli.checkAIAvailability = func(_ context.Context, profile adapters.AIProviderProfile, channel string, _ int, _ aiProbeReceiver, _ string) (adapters.AIAvailabilityRecord, error) {
		return adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
			ProbeID:           "probe-cli-available",
			ProfileID:         profile.ProfileID,
			ProfileConfigHash: mustAIProfileHash(t, profile),
			Channel:           channel,
			ProbedAt:          "2026-09-11T08:00:00.000Z",
			Status:            "available",
			RequestsUsed:      1,
			Evidence:          []string{"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		})
	}
	code, stdout, stderr := runCLI(t, cli,
		"ai", "check", "--chain", configPath, "--channel", "agent-runner-cli",
		"--copilot", executablePath, "--max-requests", "1", "--out", filepath.Join(root, "available.json"), "--json",
	)
	if code != exitSuccess || stderr != "" {
		t.Fatalf("unexpected command result: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if response := decodeResponse(t, stdout); !response.OK || response.ExitCode != exitSuccess {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestAICheckRequiresSingleRequestBudgetWithoutProbe(t *testing.T) {
	cli := newTestCLI(t.TempDir())
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, aiProbeReceiver, string) (adapters.AIAvailabilityRecord, error) {
		t.Fatal("probe must not run for invalid budget")
		return adapters.AIAvailabilityRecord{}, nil
	}
	for _, budget := range []string{"0", "2"} {
		code, _, _ := runCLI(t, cli,
			"ai", "check", "--channel", "agent-runner-cli", "--copilot", "copilot.exe",
			"--max-requests", budget, "--out", filepath.Join(t.TempDir(), "availability.json"), "--json",
		)
		if code != exitUsage {
			t.Fatalf("budget %s did not return usage failure: %d", budget, code)
		}
	}
}

func TestAICheckRequiresOutputWithoutProbe(t *testing.T) {
	cli := newTestCLI(t.TempDir())
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, aiProbeReceiver, string) (adapters.AIAvailabilityRecord, error) {
		t.Fatal("probe must not run without a durable output path")
		return adapters.AIAvailabilityRecord{}, nil
	}
	code, _, _ := runCLI(t, cli,
		"ai", "check", "--channel", "agent-runner-cli", "--copilot", "copilot.exe",
		"--max-requests", "1", "--json",
	)
	if code != exitUsage {
		t.Fatalf("missing --out did not return usage failure: %d", code)
	}
}

func TestAICheckRejectsExistingOutputBeforeProbe(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, DefaultChainConfigName)
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	outputPath := filepath.Join(root, "availability.json")
	if err := os.WriteFile(outputPath, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}
	cli := newTestCLI(root)
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, aiProbeReceiver, string) (adapters.AIAvailabilityRecord, error) {
		t.Fatal("existing output must be rejected before a paid probe")
		return adapters.AIAvailabilityRecord{}, nil
	}
	code, _, _ := runCLI(t, cli,
		"ai", "check", "--chain", configPath, "--channel", "agent-runner-cli",
		"--copilot", filepath.Join(root, "copilot.exe"), "--max-requests", "1", "--out", outputPath, "--json",
	)
	if code != exitFailure {
		t.Fatalf("existing output did not fail closed: code=%d", code)
	}
	wire, err := os.ReadFile(outputPath)
	if err != nil || string(wire) != "existing" {
		t.Fatalf("existing output changed: wire=%q err=%v", wire, err)
	}
}

func TestAIVerifyChecksPersistedRecordWithoutProbe(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, DefaultChainConfigName)
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	profile := configuredAIChain().AI.Profiles[0]
	record, err := adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
		ProbeID:           "probe-cli-persisted",
		ProfileID:         profile.ProfileID,
		ProfileConfigHash: mustAIProfileHash(t, profile),
		Channel:           "agent-runner-cli",
		ProbedAt:          "2026-09-11T08:00:00.000Z",
		Status:            "available",
		RequestsUsed:      1,
		Evidence:          []string{"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	})
	if err != nil {
		t.Fatal(err)
	}
	recordPath := filepath.Join(root, "availability.json")
	if err := adapters.WriteAIAvailabilityRecord(recordPath, record); err != nil {
		t.Fatal(err)
	}
	cli := newTestCLI(root)
	cli.now = func() time.Time { return time.Date(2026, 9, 11, 8, 5, 0, 0, time.UTC) }
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, aiProbeReceiver, string) (adapters.AIAvailabilityRecord, error) {
		t.Fatal("ai verify must not invoke a live probe")
		return adapters.AIAvailabilityRecord{}, nil
	}
	code, stdout, stderr := runCLI(t, cli,
		"ai", "verify", "--chain", configPath, "--record", recordPath,
		"--channel", "agent-runner-cli", "--max-age", "10m", "--max-requests", "1", "--json",
	)
	if code != exitSuccess || stderr != "" {
		t.Fatalf("fresh record failed verification: code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if response := decodeResponse(t, stdout); !response.OK {
		t.Fatalf("unexpected verify response: %+v", response)
	}

	cli.now = func() time.Time { return time.Date(2026, 9, 11, 8, 11, 0, 0, time.UTC) }
	code, _, _ = runCLI(t, cli,
		"ai", "verify", "--chain", configPath, "--record", recordPath,
		"--channel", "agent-runner-cli", "--max-age", "10m", "--max-requests", "1", "--json",
	)
	if code != exitFailure {
		t.Fatalf("stale record did not fail closed: code=%d", code)
	}
}

func TestAIVerifyRejectsConflictingPriorRecord(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, DefaultChainConfigName)
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	profile := configuredAIChain().AI.Profiles[0]
	profileHash := mustAIProfileHash(t, profile)
	current, err := adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
		ProbeID:           "probe-cli-conflict",
		ProfileID:         profile.ProfileID,
		ProfileConfigHash: profileHash,
		Channel:           "agent-runner-cli",
		ProbedAt:          "2026-09-11T08:00:00.000Z",
		Status:            "available",
		RequestsUsed:      1,
		Evidence:          []string{"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	})
	if err != nil {
		t.Fatal(err)
	}
	reason := "timeout"
	prior, err := adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
		ProbeID:           current.Availability.ProbeID,
		ProfileID:         profile.ProfileID,
		ProfileConfigHash: profileHash,
		Channel:           "agent-runner-cli",
		ProbedAt:          "2026-09-11T07:59:00.000Z",
		Status:            "unavailable",
		Reason:            &reason,
		RequestsUsed:      0,
		Evidence:          []string{"sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
	})
	if err != nil {
		t.Fatal(err)
	}
	currentPath := writeAIAvailabilityFixture(t, root, "current.json", current)
	priorPath := writeAIAvailabilityFixture(t, root, "prior.json", prior)
	cli := newTestCLI(root)
	cli.now = func() time.Time { return time.Date(2026, 9, 11, 8, 5, 0, 0, time.UTC) }
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, aiProbeReceiver, string) (adapters.AIAvailabilityRecord, error) {
		t.Fatal("ai verify must remain offline")
		return adapters.AIAvailabilityRecord{}, nil
	}
	code, _, _ := runCLI(t, cli,
		"ai", "verify", "--chain", configPath, "--record", currentPath,
		"--prior-record", priorPath, "--channel", "agent-runner-cli",
		"--max-age", "10m", "--max-requests", "1", "--json",
	)
	if code != exitFailure {
		t.Fatalf("conflicting prior record did not fail closed: code=%d", code)
	}
	code, _, _ = runCLI(t, cli,
		"ai", "verify", "--chain", configPath, "--record", currentPath,
		"--prior-record", currentPath, "--channel", "agent-runner-cli",
		"--max-age", "10m", "--max-requests", "1", "--json",
	)
	if code != exitSuccess {
		t.Fatalf("identical prior record was not idempotent: code=%d", code)
	}
}

func mustAIProfileHash(t *testing.T, profile adapters.AIProviderProfile) string {
	t.Helper()
	hash, err := adapters.AIProviderProfileHash(profile)
	if err != nil {
		t.Fatal(err)
	}
	return hash
}

// recordingProbeSecrets counts credential resolutions so a test can prove the
// secret store was never touched. onResolve and credential let a test act inside
// the resolution window - between the receiver binding and the spawn - without
// replacing the production availability checker.
type recordingProbeSecrets struct {
	resolveCalls int
	onResolve    func()
	credential   []byte
}

func (manager *recordingProbeSecrets) Resolve(context.Context, string) ([]byte, error) {
	manager.resolveCalls++
	if manager.onResolve != nil {
		manager.onResolve()
	}
	if len(manager.credential) == 0 {
		return nil, adapters.ErrAISecretNotFound
	}
	return manager.credential, nil
}

func (manager *recordingProbeSecrets) Set(context.Context, string, []byte) error { return nil }

func (manager *recordingProbeSecrets) Exists(context.Context, string) (bool, error) {
	return false, nil
}

func (manager *recordingProbeSecrets) Delete(context.Context, string) error { return nil }

// writeAIProbeTestChain writes the shared configured chain next to the probe.
func writeAIProbeTestChain(t *testing.T, root string) string {
	t.Helper()
	configPath := filepath.Join(root, DefaultChainConfigName)
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	return configPath
}

// TestAICheckRefusesReceiverMismatchWithoutReadingStore is the F6 regression: a
// probe whose receiver cannot be bound must be refused by the receiver binding
// or by the executor's pre-spawn re-verification.
//
// Both counterexamples hand the command an executable that EXISTS, so a refusal
// can never be explained by a path that was never there (the shape that made the
// first version of this test unable to bite).
//
// Mutation checks. (1) Deleting the validateAIProbeBinding call enters the
// availability check in the first subtest. (2) Deleting the executor's
// ExpectedExecutableHash re-verification spawns the swapped image in the second
// subtest and downgrades the record from "unavailable" to "unknown".
func TestAICheckRefusesReceiverMismatchWithoutReadingStore(t *testing.T) {
	t.Run("receiver-binding-refuses-before-the-availability-check", func(t *testing.T) {
		root := t.TempDir()
		configPath := writeAIProbeTestChain(t, root)
		// The executable exists and is readable, so its digest is computable; it is
		// one byte past the 64 MiB cap the receiver binding admits, so the binding
		// itself is what refuses it, not a missing path.
		executablePath := filepath.Join(root, "copilot.exe")
		oversized, err := os.Create(executablePath)
		if err != nil {
			t.Fatal(err)
		}
		if err := oversized.Truncate(64<<20 + 1); err != nil {
			_ = oversized.Close()
			t.Fatal(err)
		}
		if err := oversized.Close(); err != nil {
			t.Fatal(err)
		}

		secrets := &recordingProbeSecrets{}
		cli := newTestCLI(root)
		cli.secretManager = secrets
		// The seam wraps the production checker instead of replacing it: F6 promises
		// that the availability check does not run at all, so counting entries keeps
		// the real path in place while making that promise observable.
		checkerEntries := 0
		cli.checkAIAvailability = func(ctx context.Context, profile adapters.AIProviderProfile, channel string, maximumRequests int, receiver aiProbeReceiver, workingDirectory string) (adapters.AIAvailabilityRecord, error) {
			checkerEntries++
			return cli.defaultAICheck(ctx, profile, channel, maximumRequests, receiver, workingDirectory)
		}

		outputPath := filepath.Join(root, "binding-refused.json")
		code, stdout, _ := runCLI(t, cli,
			"ai", "check", "--chain", configPath, "--channel", "agent-runner-cli",
			"--copilot", executablePath, "--max-requests", "1", "--out", outputPath, "--json",
		)
		if code != exitFailure {
			t.Fatalf("an unbindable receiver must fail closed, got code=%d", code)
		}
		if !strings.Contains(stdout, "receiver binding") {
			t.Fatalf("the receiver binding must be the refusing guard, got %q", stdout)
		}
		if checkerEntries != 0 {
			t.Fatalf("the availability check must not run before the receiver binding holds, entries=%d", checkerEntries)
		}
		if secrets.resolveCalls != 0 {
			t.Fatalf("the secret store must not be read before the receiver binding holds, got %d resolutions", secrets.resolveCalls)
		}
		if _, statErr := os.Lstat(outputPath); !os.IsNotExist(statErr) {
			t.Fatalf("no availability record may be written, got err=%v", statErr)
		}
	})

	t.Run("executor-refuses-a-swapped-executable-before-the-spawn", func(t *testing.T) {
		// T3 (finding F6) defers this leg, and P85 r1 F4 turns the deferral into an
		// EXEMPTION REGISTRATION (four elements in the skip message, same text as the
		// entry for this slice's report §8 and the ⑦ input package).
		t.Skip("EXEMPT-REGISTERED P85/F4 | reason: the console install source gate is hardwired to the live platform observation, so every executable that reaches the spawn is one this user cannot replace - and a leg whose whole point is to replace the approved image cannot use such an executable; no injection seam exists at the console layer (accepted by the ⑤ r1 review) | carried-by: TestP83T3cExecutableReplacedAfterApprovalRefused (internal/adapters/p83_fixround5_tester_test.go:40), which pins the pre-spawn identity-and-digest re-verification on the executor path | destination: this slice's report §8 and the ⑦ review input package")
	})
}
