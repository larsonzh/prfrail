package console

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
)

type panicSecretManager struct{}

func (panicSecretManager) Set(context.Context, string, []byte) error {
	panic("static config command attempted a secret write")
}

func (panicSecretManager) Resolve(context.Context, string) ([]byte, error) {
	panic("static config command attempted a secret read")
}

func (panicSecretManager) Exists(context.Context, string) (bool, error) {
	panic("static config command attempted a secret existence check")
}

func (panicSecretManager) Delete(context.Context, string) error {
	panic("static config command attempted a secret delete")
}

const aiConfigSecretReferenceCanary = "windows-credential:ProofRail/config-canary"

func configuredAIChain() ChainConfig {
	config := DefaultChainConfig("ai-config")
	config.AI = &AIConfig{
		Profiles: []adapters.AIProviderProfile{
			{
				ProfileID:    "deepseek-anthropic",
				ProviderType: "anthropic",
				BaseURL:      "https://api.deepseek.com/anthropic",
				Model:        "deepseek-flash",
				AuthMode:     "secret",
				SecretRef:    aiConfigSecretReferenceCanary,
			},
			{
				ProfileID:    "github",
				ProviderType: "github-copilot",
				Model:        "auto",
				AuthMode:     "host",
			},
		},
		Channels: []AIChannelBinding{
			{Channel: "agent-runner-cli", ProfileID: "deepseek-anthropic"},
			{Channel: "sessionbridge-silent", ProfileID: "github"},
		},
	}
	return config
}

func TestAIConfigRoundTripAndExplainRedactsSecretReference(t *testing.T) {
	config := configuredAIChain()
	encoded, err := EncodeChainConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	var decoded ChainConfig
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if err := ValidateChainConfig(decoded); err != nil {
		t.Fatal(err)
	}
	report := BuildConfigExplain("proofrail.chain.json", decoded, encoded)
	if len(report.AIProfiles) != 2 || len(report.AIChannels) != 2 {
		t.Fatalf("unexpected AI explanation: %+v", report)
	}
	explained, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(explained), aiConfigSecretReferenceCanary) || strings.Contains(string(explained), "secretRef") {
		t.Fatalf("config explain disclosed secret reference: %s", explained)
	}
	if report.AIChannels[0].ProfileConfigHash != report.AIProfiles[0].ProfileConfigHash {
		t.Fatalf("channel/profile hash mismatch: %+v", report)
	}
}

func TestAIConfigRejectsInvalidOrderingBindingAndProfile(t *testing.T) {
	for name, mutate := range map[string]func(*ChainConfig){
		"profiles-unsorted": func(config *ChainConfig) {
			config.AI.Profiles[0], config.AI.Profiles[1] = config.AI.Profiles[1], config.AI.Profiles[0]
		},
		"duplicate-profile": func(config *ChainConfig) {
			config.AI.Profiles[1].ProfileID = config.AI.Profiles[0].ProfileID
		},
		"channels-unsorted": func(config *ChainConfig) {
			config.AI.Channels[0], config.AI.Channels[1] = config.AI.Channels[1], config.AI.Channels[0]
		},
		"duplicate-channel": func(config *ChainConfig) {
			config.AI.Channels[1].Channel = config.AI.Channels[0].Channel
		},
		"dangling-profile": func(config *ChainConfig) {
			config.AI.Channels[0].ProfileID = "missing"
		},
		"unsafe-base-url": func(config *ChainConfig) {
			config.AI.Profiles[0].BaseURL = "http://api.deepseek.com"
		},
		"host-with-secret": func(config *ChainConfig) {
			config.AI.Profiles[1].SecretRef = aiConfigSecretReferenceCanary
		},
	} {
		t.Run(name, func(t *testing.T) {
			config := configuredAIChain()
			mutate(&config)
			if err := ValidateChainConfig(config); !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("expected invalid config, got %v", err)
			}
		})
	}
}

func TestAIConfigRequiresProfilesAndChannels(t *testing.T) {
	for name, ai := range map[string]*AIConfig{
		"empty":         {},
		"profiles-only": {Profiles: configuredAIChain().AI.Profiles},
		"channels-only": {Channels: configuredAIChain().AI.Channels},
	} {
		t.Run(name, func(t *testing.T) {
			config := DefaultChainConfig("ai-config")
			config.AI = ai
			if err := ValidateChainConfig(config); !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("expected invalid config, got %v", err)
			}
		})
	}
}

func TestAIConfigCLIValidateAndExplainDoNotReadSecretStore(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, DefaultChainConfigName)
	encoded, err := EncodeChainConfig(configuredAIChain())
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, encoded, 0600); err != nil {
		t.Fatal(err)
	}
	cli := newTestCLI(root)
	cli.secretManager = panicSecretManager{}

	for _, command := range [][]string{
		{"validate", "--chain", configPath, "--json"},
		{"config", "explain", "--chain", configPath, "--json"},
		{"config", "explain", "--chain", configPath},
		{"preview", "--chain", configPath, "--json"},
	} {
		code, stdout, stderr := runCLI(t, cli, command...)
		if code != exitSuccess || stderr != "" {
			t.Fatalf("static config command failed: command=%v code=%d stdout=%q stderr=%q", command, code, stdout, stderr)
		}
		if strings.Contains(stdout, aiConfigSecretReferenceCanary) || strings.Contains(stdout, "secretRef") {
			t.Fatalf("static config command disclosed secret reference: %s", stdout)
		}
	}
}
