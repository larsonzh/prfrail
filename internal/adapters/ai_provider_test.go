package adapters

import (
	"errors"
	"testing"
)

func TestAIProviderProfileHashAcceptsHostAndSecretAuthentication(t *testing.T) {
	profiles := []AIProviderProfile{
		{
			ProfileID:    "github",
			ProviderType: "github-copilot",
			Model:        "auto",
			AuthMode:     "host",
		},
		{
			ProfileID:    "deepseek-anthropic",
			ProviderType: "anthropic",
			BaseURL:      "https://api.deepseek.com/anthropic",
			Model:        "deepseek-flash",
			AuthMode:     "secret",
			SecretRef:    "windows-credential:proofrail/deepseek",
		},
	}
	for _, profile := range profiles {
		t.Run(profile.ProfileID, func(t *testing.T) {
			hash, err := AIProviderProfileHash(profile)
			if err != nil {
				t.Fatal(err)
			}
			if !isHash(hash) {
				t.Fatalf("invalid profile hash: %s", hash)
			}
		})
	}
}

func TestAIProviderProfileRejectsUnsafeOrIncompleteConfiguration(t *testing.T) {
	base := AIProviderProfile{
		ProfileID:    "deepseek-anthropic",
		ProviderType: "anthropic",
		BaseURL:      "https://api.deepseek.com/anthropic",
		Model:        "deepseek-flash",
		AuthMode:     "secret",
		SecretRef:    "windows-credential:proofrail/deepseek",
	}
	for name, mutate := range map[string]func(*AIProviderProfile){
		"missing-secret-ref": func(profile *AIProviderProfile) { profile.SecretRef = "" },
		"plaintext-http":     func(profile *AIProviderProfile) { profile.BaseURL = "http://api.deepseek.com" },
		"url-credential":     func(profile *AIProviderProfile) { profile.BaseURL = "https://token@api.deepseek.com" },
		"url-query":          func(profile *AIProviderProfile) { profile.BaseURL = "https://api.deepseek.com?key=value" },
		"url-backslash":      func(profile *AIProviderProfile) { profile.BaseURL = `https://api.deepseek.com\anthropic` },
		"invalid-host":       func(profile *AIProviderProfile) { profile.BaseURL = "https://-api.deepseek.com" },
		"host-with-secret": func(profile *AIProviderProfile) {
			profile.AuthMode = "host"
		},
		"invalid-model": func(profile *AIProviderProfile) { profile.Model = " deepseek-flash" },
	} {
		t.Run(name, func(t *testing.T) {
			profile := base
			mutate(&profile)
			if err := ValidateAIProviderProfile(profile); !errors.Is(err, ErrInvalidAIProviderProfile) {
				t.Fatalf("expected invalid profile, got %v", err)
			}
		})
	}
}

func isHash(value string) bool {
	return len(value) == len("sha256:")+64 && value[:len("sha256:")] == "sha256:"
}
