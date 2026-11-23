//go:build windows || linux

package adapters

// P85 R7 (review r6 finding): a PRESENT slot whose observation cannot render a line
// must keep its position in the fixed four-slot vector, recording the canonical
// empty marker. Before this, the slot was dropped and the refusal digest described a
// shorter vector than the contract fixes - a digest that looked like a two-slot
// domain.

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// TestP85R7UnobservableSlotKeepsItsPositionInTheDigest compares the refusal digest
// with the digest of the same domain whose only slot is the canonical empty marker.
func TestP85R7UnobservableSlotKeepsItsPositionInTheDigest(t *testing.T) {
	if runtime.GOOS != "windows" && os.Geteuid() == 0 {
		t.Skip("the Linux half refuses as root before it observes any slot")
	}
	missing := filepath.Join(t.TempDir(), "gone.policy.json")
	digest, err := managedAgentPolicyDomainWriteProofSlots([]managedAgentPolicyDomainSlot{
		{label: managedPolicyDomainArtifactRole, path: missing, present: true},
	})
	if err == nil {
		t.Fatal("a present slot that cannot be observed must refuse")
	}
	if digest == "" {
		t.Fatal("the refusal must still carry the observation digest")
	}
	want := evidence.Digest(managedAgentPolicyWriteProofDomain, []byte(managedAgentPolicyDomainEmptySlot))
	if digest != want {
		t.Fatalf("an unobservable slot must record the canonical empty marker in place (%s), got %s", want, digest)
	}
	if empty := evidence.Digest(managedAgentPolicyWriteProofDomain, nil); digest == empty {
		t.Fatal("an unobservable present slot must keep its position instead of shortening the vector")
	}
}
