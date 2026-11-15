//go:build windows

package guard

import (
	"context"
	"testing"
)

// TestWindowsSandboxReportsUnknownHonestly pins the stage-S1 Windows contract
// (form alpha): the platform implementation runs no check, so it says exactly
// that -- Unknown, a non-empty reason, no error -- and the gate still refuses
// the dispatch fail-closed.
func TestWindowsSandboxReportsUnknownHonestly(t *testing.T) {
	verdict, err := NewSandbox().Check(context.Background(), CheckSpec{Reason: "unit"})
	if err != nil {
		t.Fatalf("the alpha implementation must not error, got %v", err)
	}
	if verdict.Status != VerdictUnknown {
		t.Fatalf("status = %s, want unknown: S1 must not claim that the environment provably complies or provably does not", verdict.Status)
	}
	if verdict.IsAllowed() {
		t.Fatal("the alpha implementation must never allow a dispatch")
	}
	if verdict.Reason == "" {
		t.Fatal("the unknown verdict must explain itself")
	}
	if err := CheckSandbox(context.Background(), NewSandbox(), CheckSpec{Reason: "unit"}); err == nil {
		t.Fatal("the gate must refuse a dispatch on the alpha implementation")
	}
}
