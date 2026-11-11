package console

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/chain"
)

// TestSupervisedScopeAllowsRejectsJunctionResolvingOutsideRoot is the Windows
// equivalent of the symlink-resolution escape leg: on hosts where symlink
// creation is unavailable, a directory junction (created without elevation)
// resolving outside the workspace root must be rejected by scope.allows with
// the same out-of-bounds error the symlink leg pins, and an in-root junction
// must pass as a positive control. Added by P64 ruling 2 branch A.
func TestSupervisedScopeAllowsRejectsJunctionResolvingOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.WriteFile(filepath.Join(outside, "target.txt"), []byte("outside the root\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "escape-dir")
	out, err := exec.Command("cmd", "/c", "mklink", "/J", link, outside).CombinedOutput()
	if err != nil {
		t.Skipf("junction creation unavailable: %v: %s", err, out)
	}
	scope, err := newSupervisedScope(chain.HandoffPolicy{AllowedTargets: []string{"target-one"}}, root)
	if err != nil {
		t.Fatalf("the scope must derive for an existing root, got %v", err)
	}
	if err := scope.allows("escape-dir"); !errors.Is(err, ErrSupervisedScopeOutOfBounds) {
		t.Fatalf("a junction resolving outside the root must be out of bounds, got %v", err)
	}
	inner := filepath.Join(root, "inner")
	if err := os.Mkdir(inner, 0o755); err != nil {
		t.Fatal(err)
	}
	if out, err := exec.Command("cmd", "/c", "mklink", "/J", filepath.Join(root, "in-link"), inner).CombinedOutput(); err != nil {
		t.Fatalf("inner junction creation failed: %v: %s", err, out)
	}
	if err := scope.allows("in-link"); err != nil {
		t.Fatalf("an in-root junction must be allowed, got %v", err)
	}
}
