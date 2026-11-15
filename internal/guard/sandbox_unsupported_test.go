//go:build !windows

package guard

import (
	"context"
	"errors"
	"testing"
)

// TestUnsupportedPlatformReturnsUnknown pins the non-Windows platform contract:
// the platform has no check capability, so every check reports Unknown together
// with ErrSandboxUnavailable (fail-closed) and says why.
func TestUnsupportedPlatformReturnsUnknown(t *testing.T) {
	verdict, err := NewSandbox().Check(context.Background(), CheckSpec{Reason: "unit"})
	if verdict.Status != VerdictUnknown {
		t.Fatalf("status = %s, want unknown", verdict.Status)
	}
	if !errors.Is(err, ErrSandboxUnavailable) {
		t.Fatalf("err = %v, want ErrSandboxUnavailable", err)
	}
	if verdict.IsAllowed() {
		t.Fatal("an unsupported platform must never allow a dispatch")
	}
	if verdict.Reason == "" {
		t.Fatal("the unknown verdict must explain itself")
	}
}
