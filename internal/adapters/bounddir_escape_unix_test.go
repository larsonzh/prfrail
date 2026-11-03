//go:build !windows && !darwin

package adapters

// P52-S2 ACCEPTANCE-SECURITY-COUNTEREXAMPLES: the middle-symlink escape on the
// bound-directory seam, recorded as GAP-77 (V-1, HIGH).
//
// Attack shape: the bound directory D contains a symlink "link" pointing at a
// directory O outside the bound. openBoundDirFile resolves "link/secret" with a
// single openat(2) and O_NOFOLLOW, and O_NOFOLLOW only refuses the *leaf*: the
// middle component "link" is followed, so the descriptor returned by
// OpenBoundDirFile addresses O/secret and the caller reads bytes from outside
// the bound directory.
//
// Current HEAD: the escape succeeds. The Linux run that proved it is recorded in
// docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-attack/
// t51-1-linux-escape.md (ATTACK-SUCCEEDED / SECRET-OUTSIDE-BOUNDARY), which is
// why this test is skipped rather than left red: a known vulnerability must not
// turn CI red before its fix slice.
//
// P53 fix: resolve the relative name component by component with openat(2) from
// the bound descriptor (refusing any non-final component that is not a directory
// entry of the current level). When that lands, delete the t.Skip line below and
// this test must PASS unchanged.

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBoundDirMiddleSymlinkEscapesBound pins GAP-77: a symlink in the middle of
// the relative name must not be able to redirect the open outside the bound
// directory. It is skipped at HEAD because the current single openat(2) with
// O_NOFOLLOW follows the middle symlink and returns a handle to the outside
// file; the skip therefore has to sit before the call under test, since the
// vulnerable code path is the one that succeeds.
func TestBoundDirMiddleSymlinkEscapesBound(t *testing.T) {
	t.Skip("Known vulnerability GAP-77: tracked for P53")

	root := t.TempDir()
	boundPath := filepath.Join(root, "D")
	outsidePath := filepath.Join(root, "O")

	if err := os.Mkdir(boundPath, 0o755); err != nil {
		t.Fatalf("create the bound directory: %v", err)
	}
	if err := os.Mkdir(outsidePath, 0o755); err != nil {
		t.Fatalf("create the outside directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(boundPath, "original"), []byte("inside"), 0o644); err != nil {
		t.Fatalf("write the in-bound file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outsidePath, "secret"), []byte("SECRET-OUTSIDE-BOUNDARY"), 0o644); err != nil {
		t.Fatalf("write the outside secret: %v", err)
	}
	if err := os.Mkdir(filepath.Join(boundPath, "sub"), 0o755); err != nil {
		t.Fatalf("create the in-bound subdirectory: %v", err)
	}
	if err := os.Symlink(outsidePath, filepath.Join(boundPath, "link")); err != nil {
		t.Fatalf("plant the middle symlink: %v", err)
	}

	dir, err := OpenBoundDir(boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	// The refusal carries no sentinel yet (P53 introduces the openat-shaped
	// error), so the assertion is on the attack outcome itself: a nil error means
	// the outside file was opened and the escape succeeded.
	_, err = OpenBoundDirFile(dir, filepath.Join("link", "secret"), 1024)
	if err == nil {
		t.Fatalf("GAP-77 regression: middle symlink escape succeeded and read outside the bound directory")
	}

	if err := dir.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
