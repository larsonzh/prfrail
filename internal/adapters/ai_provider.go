package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const aiProviderProfileDomain = "proofrail:ai-provider-profile:1\n"

var (
	ErrInvalidAIProviderProfile = errors.New("invalid AI provider profile")
	ErrAISecretNotFound         = errors.New("AI secret not found")
)

type AIProviderProfile struct {
	ProfileID    string `json:"profileId"`
	ProviderType string `json:"providerType"`
	BaseURL      string `json:"baseUrl,omitempty"`
	Model        string `json:"model"`
	WireAPI      string `json:"wireApi,omitempty"`
	AuthMode     string `json:"authMode"`
	SecretRef    string `json:"secretRef,omitempty"`
}

type AISecretResolver interface {
	Resolve(context.Context, string) ([]byte, error)
}

type AISecretManager interface {
	AISecretResolver
	Set(context.Context, string, []byte) error
	Exists(context.Context, string) (bool, error)
	Delete(context.Context, string) error
}

func AIProviderProfileHash(profile AIProviderProfile) (string, error) {
	if err := ValidateAIProviderProfile(profile); err != nil {
		return "", err
	}
	return digestMessage(aiProviderProfileDomain, profile)
}

func ValidateAIProviderProfile(profile AIProviderProfile) error {
	if !evidence.ValidID(profile.ProfileID) || !evidence.ValidID(profile.ProviderType) {
		return fmt.Errorf("%w: invalid identity", ErrInvalidAIProviderProfile)
	}
	if profile.Model == "" || len(profile.Model) > 128 || strings.TrimSpace(profile.Model) != profile.Model {
		return fmt.Errorf("%w: invalid model", ErrInvalidAIProviderProfile)
	}
	if profile.WireAPI != "" && !evidence.ValidID(profile.WireAPI) {
		return fmt.Errorf("%w: invalid wire API", ErrInvalidAIProviderProfile)
	}
	switch profile.AuthMode {
	case "host":
		if profile.SecretRef != "" {
			return fmt.Errorf("%w: host authentication cannot contain secretRef", ErrInvalidAIProviderProfile)
		}
	case "secret":
		if !validAISecretReference(profile.SecretRef) {
			return fmt.Errorf("%w: invalid secretRef", ErrInvalidAIProviderProfile)
		}
	default:
		return fmt.Errorf("%w: invalid auth mode", ErrInvalidAIProviderProfile)
	}
	if profile.BaseURL == "" {
		if profile.AuthMode == "secret" {
			return fmt.Errorf("%w: secret authentication requires baseUrl", ErrInvalidAIProviderProfile)
		}
		return nil
	}
	if !validHTTPSBaseURL(profile.BaseURL) {
		return fmt.Errorf("%w: baseUrl must be an HTTPS origin or path without credentials, query, or fragment", ErrInvalidAIProviderProfile)
	}
	return nil
}

func validHTTPSBaseURL(value string) bool {
	const prefix = "https://"
	if !strings.HasPrefix(value, prefix) || strings.ContainsAny(value, "@?#\\\r\n\t ") {
		return false
	}
	remainder := strings.TrimPrefix(value, prefix)
	authority, _, _ := strings.Cut(remainder, "/")
	if authority == "" || strings.HasPrefix(authority, ".") || strings.HasSuffix(authority, ".") || strings.Contains(authority, "..") {
		return false
	}
	host, port, hasPort := strings.Cut(authority, ":")
	if host == "" || strings.Contains(host, ":") || !validDNSName(host) {
		return false
	}
	if hasPort && (port == "" || strings.Trim(port, "0123456789") != "") {
		return false
	}
	return true
}

func validDNSName(host string) bool {
	for _, label := range strings.Split(host, ".") {
		if label == "" || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, character := range label {
			if character != '-' && (character < 'a' || character > 'z') && (character < 'A' || character > 'Z') && (character < '0' || character > '9') {
				return false
			}
		}
	}
	return true
}

func validAISecretReference(reference string) bool {
	if len(reference) < 3 || len(reference) > 256 || strings.TrimSpace(reference) != reference || strings.ContainsAny(reference, "\r\n\t") {
		return false
	}
	separator := strings.IndexByte(reference, ':')
	return separator > 0 && separator < len(reference)-1
}
