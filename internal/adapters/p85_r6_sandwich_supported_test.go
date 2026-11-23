//go:build windows || linux

package adapters

// P85 R6 (test round 6) / review r3 finding 1: the identity SANDWICH of the
// handle-form write-domain observation. The observation reads a path-based probe
// (the POSIX access-ACL read, the per-right Windows access checks) between two
// name-versus-handle comparisons, so a name that no longer addresses the object
// the handle addresses refuses instead of judging another object.
//
// The counterexample is the swap the review names: the artifact is renamed away and
// a byte-identical stand-in is put at the same name. The handle keeps addressing the
// recorded object, the name addresses a different one, and no byte in the domain
// changed - so only an identity comparison can tell the two apart.
//
// Boundary (registered): a swap squeezed strictly BETWEEN the two comparisons of one
// observation cannot be injected without a production seam, and no seam is added
// here. The legs below therefore (a) prove the swap is a different object, and
// (b) prove each platform's observation refuses a name that moved, which is the
// discriminating fact both halves rely on.
//
// Mutation check (design rev 4 section 6, mutation 2): deleting one name comparison
// from the sandwich - the Windows handle-versus-name comparison, or the Linux
// managedAgentPolicyDomainAssertNameIdentity mismatch branch - makes the
// corresponding leg below accept the moved name and go red.

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85R6SandwichLeg is one platform-native leg of the sandwich counterexample.
type p85R6SandwichLeg struct {
	name string
	run  func(t *testing.T)
}

// p85R6SandwichPlatformLegs collects the legs each platform half registers in its
// own init, so the shared test never has to name a platform symbol that another
// build excludes.
var p85R6SandwichPlatformLegs []p85R6SandwichLeg

// p85R6SandwichFixture binds one artifact inside one bound directory and returns
// the handle-form slot, the identity its handle reports, and the swap that renames
// the artifact away and leaves a byte-identical stand-in at the same name.
func p85R6SandwichFixture(t *testing.T) (managedAgentPolicyDomainSlot, evidence.FileID, func()) {
	t.Helper()
	root := t.TempDir()
	policy := filepath.Join(root, "agent-runner-policy")
	if err := os.MkdirAll(policy, 0o755); err != nil {
		t.Fatalf("create the policy directory: %v", err)
	}
	name := "task.step.attempt.1.policy.json"
	path := filepath.Join(policy, name)
	if err := os.WriteFile(path, []byte("policy"), 0o600); err != nil {
		t.Fatalf("write the artifact: %v", err)
	}

	directory, _, err := evidence.OpenBoundDirVerified(policy)
	if err != nil {
		t.Skipf("this host's temporary root is not bindable as a verified directory: %v", err)
	}
	t.Cleanup(func() { _ = directory.Close() })
	bound, err := evidence.OpenBoundDirFile(directory, name, 1<<20)
	if err != nil {
		t.Fatalf("open the artifact relative to the bound directory: %v", err)
	}
	t.Cleanup(func() { _ = bound.Close() })
	identity, err := bound.FileID()
	if err != nil {
		t.Fatalf("identify the artifact on its handle: %v", err)
	}

	slot := managedAgentPolicyDomainSlot{
		label:   managedPolicyDomainArtifactRole,
		path:    path,
		isDir:   false,
		present: true,
		file:    bound,
	}
	swap := func() {
		t.Helper()
		// Rename first so the retired object keeps its identity; a create-after-delete
		// can hand the same inode back, which would make the fixture agree by accident.
		if err := os.Rename(path, path+".retired"); err != nil {
			t.Fatalf("rename the artifact away: %v", err)
		}
		if err := os.WriteFile(path, []byte("policy"), 0o600); err != nil {
			t.Fatalf("write the stand-in artifact: %v", err)
		}
	}
	return slot, evidence.FileID(identity), swap
}

// TestP85R6SandwichRefusesANameSwapBetweenAssertions is the counterexample.
func TestP85R6SandwichRefusesANameSwapBetweenAssertions(t *testing.T) {
	t.Run("the-swap-leaves-a-different-object-behind-the-same-name", func(t *testing.T) {
		slot, handleIdentity, swap := p85R6SandwichFixture(t)
		untouched, err := evidence.FileIDAt(slot.path)
		if err != nil {
			t.Fatalf("identify the untouched artifact by name: %v", err)
		}
		if untouched != handleIdentity {
			t.Fatalf("an untouched name must address the handle's object (%s != %s)", untouched, handleIdentity)
		}
		swap()
		moved, err := evidence.FileIDAt(slot.path)
		if err != nil {
			t.Fatalf("the stand-in must still be observable by name: %v", err)
		}
		if moved == handleIdentity {
			t.Fatalf("the fixture must put a DIFFERENT object behind the same name, got %s twice", moved)
		}
		// The name still resolves and the bytes are identical: only identity separates
		// the two objects, which is what makes the sandwich necessary at all.
		movedBytes, err := os.ReadFile(slot.path)
		if err != nil {
			t.Fatalf("read the stand-in by name: %v", err)
		}
		if string(movedBytes) != "policy" {
			t.Fatalf("the stand-in must carry the same bytes, got %q", movedBytes)
		}
	})

	for _, leg := range p85R6SandwichPlatformLegs {
		leg := leg
		t.Run(leg.name, leg.run)
	}
}
