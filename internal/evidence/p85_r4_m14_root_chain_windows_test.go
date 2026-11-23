//go:build windows

package evidence

// P85 design rev 4, section 3 / test row M14 (Windows native): OpenBoundDirVerified
// must prove the NAME CHAIN it binds, not just the object one CreateFile happens to
// resolve.
//
// Before rev 4 the binding handed the whole path to a single CreateFile, so a
// junction in ANY component of the directory path was followed silently by the
// kernel; the handle identity and the name-based FileIDAt then agreed with each
// other because both had landed on the same followed object. The walk refuses the
// component instead.
//
// Two discriminating shapes are used:
//
//   - a junction that leads to a DIFFERENT tree (the classic redirect);
//   - a junction that leads to the VERY SAME object (rename the directory away, put
//     a junction at its old name pointing at the new name). Here FileIDAt reports the
//     bound identity, so only a name-chain classifier can tell the difference - which
//     is exactly the state the pre-fix binding accepted.
//
// Mutation check: deleting the reparse refusal from openBoundDirChainHandle makes the
// walk descend through the junction; the refusal assertions below then fail (either
// the bind succeeds or it fails for a reason other than the reparse classification).
//
// Junctions are creatable without elevation; if this host cannot create one the legs
// skip rather than pretend.

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// p85R4Junction creates one directory junction, or skips the caller when the host
// cannot create one without elevation.
func p85R4Junction(t *testing.T, link, target string) {
	t.Helper()
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("create the junction target %s: %v", target, err)
	}
	if output, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); err != nil {
		t.Skipf("this host cannot create the junction %s -> %s without elevation: %v (%s)", link, target, err, strings.TrimSpace(string(output)))
	}
	// Registered AFTER t.TempDir's own cleanup, so it runs first (LIFO) and the
	// recursive removal never walks through the reparse point.
	t.Cleanup(func() { _ = os.Remove(link) })
}

// p85R4RequireReparseRefusal asserts the refusal is the name-chain classification
// and not an unrelated open failure.
func p85R4RequireReparseRefusal(t *testing.T, label string, err error) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s must be refused by the root-chain walk, got a successful bind", label)
	}
	if !strings.Contains(err.Error(), "is a reparse point") {
		t.Fatalf("%s must be refused BECAUSE it is a reparse point, got %v", label, err)
	}
}

// TestP85R4M14RootChainRefusesAReparseComponent is M14.
func TestP85R4M14RootChainRefusesAReparseComponent(t *testing.T) {
	root := t.TempDir()
	real := filepath.Join(root, "real")
	policy := filepath.Join(real, "policy")
	if err := os.MkdirAll(policy, 0o755); err != nil {
		t.Fatalf("create the policy directory: %v", err)
	}

	// Positive control: the same directory, reached without any reparse component,
	// binds and reports an identity.
	direct, directID, err := OpenBoundDirVerified(policy)
	if err != nil {
		t.Skipf("this host's temporary root is not link-free, so the positive control cannot be established: %v", err)
	}
	if directID == "" {
		t.Fatal("a successful binding must report the bound identity")
	}
	if err := direct.Close(); err != nil {
		t.Fatalf("close the bound directory: %v", err)
	}

	t.Run("an-intermediate-junction-is-refused", func(t *testing.T) {
		link := filepath.Join(root, "linked")
		p85R4Junction(t, link, real)

		_, _, err := OpenBoundDirVerified(filepath.Join(link, "policy"))
		p85R4RequireReparseRefusal(t, "a directory bound through an intermediate junction", err)

		// The standalone classifier - the one the write-domain observation uses - has
		// to reach the same verdict, so "the object that was judged" and "the object
		// that was bound" cannot drift apart.
		err = ClassifyRootChainNoReparse(filepath.Join(link, "policy"), true)
		p85R4RequireReparseRefusal(t, "the root-chain classifier over an intermediate junction", err)

		// Removing the junction restores the acceptance of the very same target, so
		// the refusal is attributable to the junction.
		if err := os.Remove(link); err != nil {
			t.Fatalf("remove the junction: %v", err)
		}
		restored, restoredID, err := OpenBoundDirVerified(policy)
		if err != nil {
			t.Fatalf("the same target must bind once the junction is gone, got %v", err)
		}
		defer func() { _ = restored.Close() }()
		if restoredID != directID {
			t.Fatalf("the target must still be the same object, got %s want %s", restoredID, directID)
		}
	})

	t.Run("a-last-component-junction-is-refused", func(t *testing.T) {
		link := filepath.Join(root, "policy-link")
		p85R4Junction(t, link, policy)

		_, _, err := OpenBoundDirVerified(link)
		p85R4RequireReparseRefusal(t, "a directory whose own name is a junction", err)
	})

	t.Run("a-junction-to-the-very-same-object-is-refused", func(t *testing.T) {
		// The pre-fix blind spot: the name still resolves to the bound object, so a
		// name-based identity comparison agrees with itself. Only the name CHAIN can
		// refuse this shape.
		renamed := filepath.Join(root, "policy-renamed")
		if err := os.Rename(policy, renamed); err != nil {
			t.Fatalf("rename the policy directory: %v", err)
		}
		link := policy
		p85R4Junction(t, link, renamed)

		named, err := FileIDAt(link)
		if err != nil {
			t.Fatalf("identify the junction by name: %v", err)
		}
		if string(named) != string(directID) {
			t.Fatalf("the fixture must make the name resolve to the SAME object, got %s want %s", named, directID)
		}

		_, _, err = OpenBoundDirVerified(link)
		p85R4RequireReparseRefusal(t, "a junction that resolves to the bound object", err)
	})
}
