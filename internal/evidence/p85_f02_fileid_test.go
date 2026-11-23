//go:build !darwin

package evidence

// P85 F02/F06/F20 (design rev 3, section 8, row M5): the directory identity has
// to be the object the PRE-CHECK approved, not a name that happens to resolve the
// same way twice.
//
// The counterexample is the one review r1 named: "pre-check approves D1, D1 is
// renamed away, D2 takes the same spelling, the caller then binds". The old shape
// compared the handle with a fresh os.Stat of the name, so after the swap both
// sides resolved to D2 and the check stayed green. VerifyBoundDir now compares the
// handle's own identity with the identity the binding produced, so an identity
// captured BEFORE the replacement is refused.

import (
	"os"
	"path/filepath"
	"testing"
)

func TestP85F02VerifyBoundDirRefusesAnIdentityCapturedBeforeAReplacement(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "policy")
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("create the approved directory: %v", err)
	}

	// The pre-check: one resolution that yields both the handle and the identity
	// of the object it bound. The handle is released so the name can be replaced.
	approvedDir, approved, err := OpenBoundDirVerified(path)
	if err != nil {
		t.Fatalf("bind the approved directory: %v", err)
	}
	if err := approvedDir.Close(); err != nil {
		t.Fatalf("release the approved handle: %v", err)
	}

	// The swap between the pre-check and the later bind.
	if err := os.Rename(path, path+".approved"); err != nil {
		t.Fatalf("rename the approved directory away: %v", err)
	}
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatalf("create the replacement directory: %v", err)
	}

	current, currentID, err := OpenBoundDirVerified(path)
	if err != nil {
		t.Fatalf("bind the replacement directory: %v", err)
	}
	defer func() { _ = current.Close() }()
	if currentID == approved {
		t.Fatalf("the fixture never replaced the object: both identities are %s", currentID)
	}

	// The decisive assertion: the handle holds D2 and the NAME also resolves to D2,
	// so a name-based comparison would agree - but the approved identity is D1, and
	// that is what the check has to answer.
	if err := VerifyBoundDir(current, approved); err == nil {
		t.Fatal("a handle bound after the replacement must not verify against the identity approved before it")
	}

	// Control: the same handle verifies against the identity its OWN binding
	// produced, so the refusal above is caused by the stale identity.
	if err := VerifyBoundDir(current, currentID); err != nil {
		t.Fatalf("the handle must verify against the identity its binding produced, got %v", err)
	}
}
