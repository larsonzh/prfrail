package console

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
)

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
	executablePath := filepath.Join(root, "copilot.exe")
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}

	called := false
	cli := newTestCLI(root)
	cli.checkAIAvailability = func(_ context.Context, profile adapters.AIProviderProfile, channel string, maximumRequests int, executable, workingDirectory string) (adapters.AIAvailabilityRecord, error) {
		called = true
		if profile.ProfileID != "deepseek-anthropic" || channel != "agent-runner-cli" || maximumRequests != 1 {
			t.Fatalf("unexpected binding: profile=%s channel=%s requests=%d", profile.ProfileID, channel, maximumRequests)
		}
		if executable != executablePath || workingDirectory != root {
			t.Fatalf("unexpected execution target: executable=%q cwd=%q", executable, workingDirectory)
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
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, string, string) (adapters.AIAvailabilityRecord, error) {
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
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	cli := newTestCLI(root)
	cli.checkAIAvailability = func(_ context.Context, profile adapters.AIProviderProfile, channel string, _ int, _, _ string) (adapters.AIAvailabilityRecord, error) {
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
		"--copilot", filepath.Join(root, "copilot.exe"), "--max-requests", "1", "--out", filepath.Join(root, "available.json"), "--json",
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
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, string, string) (adapters.AIAvailabilityRecord, error) {
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
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, string, string) (adapters.AIAvailabilityRecord, error) {
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
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, string, string) (adapters.AIAvailabilityRecord, error) {
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
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, string, string) (adapters.AIAvailabilityRecord, error) {
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
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, string, string) (adapters.AIAvailabilityRecord, error) {
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
