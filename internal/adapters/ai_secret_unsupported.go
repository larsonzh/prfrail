//go:build !windows

package adapters

import (
	"context"
	"errors"
)

var ErrAISecretStore = errors.New("AI secret store operation failed")

type unsupportedAISecretManager struct{}

func NewPlatformAISecretManager() AISecretManager {
	return unsupportedAISecretManager{}
}

func (unsupportedAISecretManager) Set(context.Context, string, []byte) error {
	return ErrAISecretStore
}

func (unsupportedAISecretManager) Resolve(context.Context, string) ([]byte, error) {
	return nil, ErrAISecretStore
}

func (unsupportedAISecretManager) Exists(context.Context, string) (bool, error) {
	return false, ErrAISecretStore
}

func (unsupportedAISecretManager) Delete(context.Context, string) error {
	return ErrAISecretStore
}
