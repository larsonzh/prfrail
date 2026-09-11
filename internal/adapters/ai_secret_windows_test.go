//go:build windows

package adapters

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"testing"
)

func TestWindowsCredentialManagerLifecycle(t *testing.T) {
	randomID := make([]byte, 16)
	if _, err := rand.Read(randomID); err != nil {
		t.Fatal(err)
	}
	reference := "windows-credential:ProofRail/test/" + hex.EncodeToString(randomID)
	manager := WindowsCredentialManager{}
	ctx := context.Background()
	t.Cleanup(func() {
		_ = manager.Delete(ctx, reference)
	})

	if exists, err := manager.Exists(ctx, reference); err != nil || exists {
		t.Fatalf("new test credential unexpectedly exists: exists=%v err=%v", exists, err)
	}
	secret := []byte("proofrail-non-secret-test-canary")
	if err := manager.Set(ctx, reference, secret); err != nil {
		t.Fatal(err)
	}
	if exists, err := manager.Exists(ctx, reference); err != nil || !exists {
		t.Fatalf("stored credential is unavailable: exists=%v err=%v", exists, err)
	}
	resolved, err := manager.Resolve(ctx, reference)
	if err != nil {
		t.Fatal(err)
	}
	defer clear(resolved)
	if !bytes.Equal(resolved, secret) {
		t.Fatal("resolved credential differs from stored bytes")
	}
	if err := manager.Delete(ctx, reference); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.Resolve(ctx, reference); !errors.Is(err, ErrAISecretNotFound) {
		t.Fatalf("deleted credential should be absent, got %v", err)
	}
}

func TestWindowsCredentialManagerRejectsInvalidInputWithoutSecretDisclosure(t *testing.T) {
	manager := WindowsCredentialManager{}
	secret := []byte("proofrail-secret-canary-must-not-escape")
	for _, test := range []struct {
		name      string
		reference string
		secret    []byte
	}{
		{name: "wrong-scheme", reference: "environment:proofrail/test", secret: secret},
		{name: "empty-target", reference: "windows-credential:", secret: secret},
		{name: "control-character", reference: "windows-credential:ProofRail/test\nforged", secret: secret},
		{name: "oversized-target", reference: "windows-credential:" + string(make([]byte, 32768)), secret: secret},
		{name: "empty-secret", reference: "windows-credential:ProofRail/test", secret: nil},
		{name: "oversized-secret", reference: "windows-credential:ProofRail/test", secret: make([]byte, 2561)},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := manager.Set(context.Background(), test.reference, test.secret)
			if !errors.Is(err, ErrAISecretStore) {
				t.Fatalf("expected secret store rejection, got %v", err)
			}
			if bytes.Contains([]byte(err.Error()), secret) {
				t.Fatal("secret escaped through error text")
			}
		})
	}
}
