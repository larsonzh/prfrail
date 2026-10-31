//go:build windows

package adapters

// P42-T027B-EVIDENCE-BIND ③ test-first (RED): Windows half of the directory
// binding tests (design section 2 "Q2 ancestor directory replacement").
//
// Pinned invariants (design section 4):
//   I3 a bound directory keeps reading the directory object it opened, so an
//      ancestor renamed away and re-created by an attacker cannot redirect it.
//   Risk 1 CloseBoundDirs releases every bound handle so the replay root can be
//      removed afterwards.
//
// The attack shape is rename-only: Windows junction creation needs privileges CI
// does not guarantee, and a rename plus a fresh directory at the same path models
// the same replacement without them (the design explicitly allows either shape).

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

// TestBoundDirReadSurvivesAncestorReplacement pins I3 at the platform primitive.
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

	// Replace the bound directory by moving it away and re-creating a directory
	// that an attacker controls at the same path.
	if err := os.Rename(boundPath, boundPath+".bak"); err != nil {
		t.Fatalf("rename the bound directory away (the handle must be share-delete): %v", err)
	}
	if err := os.MkdirAll(boundPath, 0o755); err != nil {
		t.Fatalf("create the injected directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(boundPath, "f.txt"), []byte("attacker"), 0o644); err != nil {
		t.Fatalf("write the injected file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(boundPath, "injected.txt"), []byte("attacker"), 0o644); err != nil {
		t.Fatalf("write the injected-only file: %v", err)
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

	// A name that exists only in the injected directory must stay missing: the
	// bound directory has no such file, so no fallback to the path may happen.
	if _, err := OpenBoundDirFile(dir, "injected.txt", evidenceBindTestBound); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a file absent from the bound directory must be fs.ErrNotExist, got %v", err)
	}
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
