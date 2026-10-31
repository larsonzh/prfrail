//go:build !windows

package adapters

// P42-T027B-EVIDENCE-BIND ③ test-first (RED): Unix half of the directory binding
// tests (design section 2 "Q2 ancestor directory replacement").
//
// The design states the Unix implementation under "//go:build !windows", so this
// file uses the same constraint instead of the narrower "unix" tag: the invariant
// must hold on every platform the implementation claims.
//
// Pinned invariants (design section 4):
//   I3 a bound directory keeps reading the directory object it opened, so an
//      ancestor replaced by a symlink to an injected directory cannot redirect it,
//      and a leaf replaced by a symlink is refused (O_NOFOLLOW).
//   Risk 1 CloseBoundDirs releases every bound handle so the replay root can be
//      removed afterwards.

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	evidenceBindTestBound  int64 = 1 << 20
	evidenceBindTestDomain       = "proofrail:test:bound-dir:1\n"
)

// TestBoundDirReadSurvivesAncestorReplacement pins I3 at the platform primitive
// with the symlink replacement shape.
func TestBoundDirReadSurvivesAncestorReplacement(t *testing.T) {
	root := t.TempDir()
	boundPath := filepath.Join(root, "ev")
	if err := os.MkdirAll(boundPath, 0o755); err != nil {
		t.Fatalf("create bound directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(boundPath, "f.txt"), []byte("bound"), 0o644); err != nil {
		t.Fatalf("write bound file: %v", err)
	}

	dir, err := OpenBoundDir(boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}

	// Move the bound directory away and point the old path at an injected tree.
	injected := filepath.Join(root, "injected")
	if err := os.MkdirAll(injected, 0o755); err != nil {
		t.Fatalf("create the injected directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(injected, "f.txt"), []byte("attacker"), 0o644); err != nil {
		t.Fatalf("write the injected file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(injected, "injected.txt"), []byte("attacker"), 0o644); err != nil {
		t.Fatalf("write the injected-only file: %v", err)
	}
	if err := os.Rename(boundPath, boundPath+".bak"); err != nil {
		t.Fatalf("rename the bound directory away: %v", err)
	}
	if err := os.Symlink(injected, boundPath); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}

	file, err := OpenBoundDirFile(dir, "f.txt", evidenceBindTestBound)
	if err != nil {
		t.Fatalf("OpenBoundDirFile on the bound directory: %v", err)
	}
	data, truncated, err := file.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll through the bound directory: %v", err)
	}
	if string(data) != "bound" {
		t.Fatalf("a dir-bound read must follow the bound directory object, got %q", data)
	}
	if truncated {
		t.Fatal("a small dir-bound read must not be truncated")
	}
	hash, err := file.Hash(evidenceBindTestDomain)
	if err != nil {
		t.Fatalf("Hash through the bound directory: %v", err)
	}
	if want := evidence.Digest(evidenceBindTestDomain, []byte("bound")); hash != want {
		t.Fatalf("dir-bound digest must cover the bound bytes: got %s want %s", hash, want)
	}

	// A name that exists only in the injected tree must stay missing: the bound
	// directory has no such entry, so no fallback to the replaced path may happen.
	if _, err := OpenBoundDirFile(dir, "injected.txt", evidenceBindTestBound); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a file absent from the bound directory must be fs.ErrNotExist, got %v", err)
	}
}

// TestBoundDirRejectsSymlinkLeaf pins the O_NOFOLLOW half of I3: a leaf swapped
// for a symlink after the directory was bound must be refused, not followed.
func TestBoundDirRejectsSymlinkLeaf(t *testing.T) {
	root := t.TempDir()
	boundPath := filepath.Join(root, "ev")
	if err := os.MkdirAll(boundPath, 0o755); err != nil {
		t.Fatalf("create bound directory: %v", err)
	}
	leaf := filepath.Join(boundPath, "g.txt")
	if err := os.WriteFile(leaf, []byte("bound"), 0o644); err != nil {
		t.Fatalf("write bound leaf: %v", err)
	}

	dir, err := OpenBoundDir(boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}

	target := filepath.Join(root, "attacker.txt")
	if err := os.WriteFile(target, []byte("attacker"), 0o644); err != nil {
		t.Fatalf("write the attacker target: %v", err)
	}
	if err := os.Remove(leaf); err != nil {
		t.Fatalf("remove the bound leaf: %v", err)
	}
	if err := os.Symlink(target, leaf); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}

	file, err := OpenBoundDirFile(dir, "g.txt", evidenceBindTestBound)
	if err != nil {
		return // refused at open: the pinned outcome.
	}
	data, _, readErr := file.ReadAll()
	t.Fatalf("a symlinked leaf must be refused by O_NOFOLLOW, got %q (read error %v)", data, readErr)
}

// TestCloseBoundDirsAllowsRemoval pins risk 1: the store owns the bound directory
// handles for the whole run, and a run close must release them.
func TestCloseBoundDirsAllowsRemoval(t *testing.T) {
	root := t.TempDir()
	store := replayStoreMustNew(t, root)
	dir, err := store.RunEvidenceDir("request-one")
	if err != nil {
		t.Fatalf("RunEvidenceDir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "f.txt"), []byte("evidence"), 0o644); err != nil {
		t.Fatalf("write evidence file: %v", err)
	}

	store.CloseBoundDirs()

	if err := os.RemoveAll(root); err != nil {
		t.Fatalf("after CloseBoundDirs the replay root must be removable, got %v", err)
	}
}
