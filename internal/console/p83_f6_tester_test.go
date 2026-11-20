package console

// P83 F6 counterexample that runs on Windows.
//
// The existing availability-check test proves the command fails when the
// executable cannot be hashed, but the default checker hashes the executable too,
// so it cannot distinguish "the binding refused" from "the checker refused". This
// test replaces the checker with a spy that would resolve a credential: with the
// binding in place the spy is never entered, so removing validateAIProbeBinding
// makes the credential resolve and turns the test red.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
)

// p83SecretSpy records every credential resolution.
type p83SecretSpy struct {
	resolveCalls int
}

func (spy *p83SecretSpy) Resolve(context.Context, string) ([]byte, error) {
	spy.resolveCalls++
	return nil, adapters.ErrAISecretNotFound
}

func (spy *p83SecretSpy) Set(context.Context, string, []byte) error { return nil }

func (spy *p83SecretSpy) Exists(context.Context, string) (bool, error) { return false, nil }

func (spy *p83SecretSpy) Delete(context.Context, string) error { return nil }

// TestAICheckRefusesUnbindableReceiverWithoutResolvingSecret asserts the F6
// ordering: an unbindable receiver is refused before the availability checker
// runs, so the secret store is not read and no record is written.
// Mutation check: deleting the validateAIProbeBinding call enters the spy, which
// resolves a credential, and turns this red.
func TestAICheckRefusesUnbindableReceiverWithoutResolvingSecret(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, DefaultChainConfigName)
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}

	secrets := &p83SecretSpy{}
	cli := newTestCLI(root)
	cli.secretManager = secrets
	checkerCalls := 0
	cli.checkAIAvailability = func(ctx context.Context, _ adapters.AIProviderProfile, _ string, _ int, _ aiProbeReceiver, _ string) (adapters.AIAvailabilityRecord, error) {
		checkerCalls++
		_, _ = cli.secretManager.Resolve(ctx, "unused-secret-ref")
		return adapters.AIAvailabilityRecord{}, errors.New("the probe must not run")
	}

	outputPath := filepath.Join(root, "tester-availability.json")
	code, _, _ := runCLI(t, cli,
		"ai", "check", "--chain", configPath, "--channel", "agent-runner-cli",
		"--copilot", filepath.Join(root, "not-the-approved-copilot"),
		"--max-requests", "1", "--out", outputPath, "--json",
	)
	if code != exitFailure {
		t.Fatalf("an unbindable receiver must fail closed, got code=%d", code)
	}
	if checkerCalls != 0 {
		t.Fatalf("the availability checker must not run before the receiver binding holds, calls=%d", checkerCalls)
	}
	if secrets.resolveCalls != 0 {
		t.Fatalf("the secret store must not be read before the receiver binding holds, resolutions=%d", secrets.resolveCalls)
	}
	if _, statErr := os.Lstat(outputPath); !os.IsNotExist(statErr) {
		t.Fatalf("no availability record may be written, stat err=%v", statErr)
	}
}
