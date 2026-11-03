//go:build windows

package adapters

// P52-S2 ACCEPTANCE-SECURITY-COUNTEREXAMPLES: junction traversal across the
// bound-directory seam (GAP-52 control shape, Windows side).
//
// The Windows half has no openat(2): openBoundDirFile resolves the relative name
// against the bound directory's current real name and then re-checks the opened
// artifact handle against that same base (confirmOpenedArtifact). A junction is
// a reparse point that the name-based open follows, so the question is whether
// the post-open handle check still holds the boundary:
//
//   - a junction that leaves the bound directory must be refused, because the
//     opened handle resolves outside the bound base;
//   - a junction that stays inside the bound directory must still be admitted,
//     because it resolves under the same base.
//
// Junction creation needs no elevation (mklink /J), unlike symbolic links; the
// tests skip when the environment cannot create one.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	escapeTestBound int64 = 1 << 10
	escapeTestLink        = "j"
)

// TestBoundDirRefusesJunctionEscape pins the refusal: an artifact reached
// through a junction pointing outside the bound directory must not be opened.
// The refusal is the boundary message from confirmOpenedArtifact; there is no
// sentinel for this seam, so the error text is asserted instead.
func TestBoundDirRefusesJunctionEscape(t *testing.T) {
	root := t.TempDir()
	boundPath := filepath.Join(root, "D")
	outsidePath := filepath.Join(root, "O")
	if err := os.Mkdir(boundPath, 0o755); err != nil {
		t.Fatalf("create the bound directory: %v", err)
	}
	if err := os.Mkdir(outsidePath, 0o755); err != nil {
		t.Fatalf("create the outside directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(outsidePath, "secret"), []byte("SECRET-OUTSIDE-BOUNDARY"), 0o644); err != nil {
		t.Fatalf("write the outside secret: %v", err)
	}

	link := filepath.Join(boundPath, escapeTestLink)
	command := exec.Command("cmd", "/c", "mklink", "/J", link, outsidePath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Skipf("junction creation unavailable: %v (%s)", err, strings.TrimSpace(string(output)))
	}
	// A junction is a directory entry; removing it removes the reparse point and
	// does not touch the outside target.
	defer func() { _ = os.Remove(link) }()

	dir, err := OpenBoundDir(boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	file, err := OpenBoundDirFile(dir, filepath.Join(escapeTestLink, "secret"), escapeTestBound)
	if err == nil {
		_ = file.Close()
		t.Fatalf("a junction leaving the bound directory must be refused, got a handle")
	}
	if !strings.Contains(err.Error(), "outside the bound directory") {
		t.Fatalf("the refusal must name the boundary violation, got %v", err)
	}

	if err := dir.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// TestBoundDirJunctionInsideBoundAdmitted is the control: a junction that
// resolves inside the bound directory is not an escape, so the artifact opens
// and its bytes are read back unchanged.
func TestBoundDirJunctionInsideBoundAdmitted(t *testing.T) {
	root := t.TempDir()
	boundPath := filepath.Join(root, "D")
	insidePath := filepath.Join(boundPath, "sub")
	if err := os.MkdirAll(insidePath, 0o755); err != nil {
		t.Fatalf("create the in-bound subdirectory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(insidePath, "f"), []byte("inside"), 0o644); err != nil {
		t.Fatalf("write the in-bound file: %v", err)
	}

	link := filepath.Join(boundPath, escapeTestLink)
	command := exec.Command("cmd", "/c", "mklink", "/J", link, insidePath)
	if output, err := command.CombinedOutput(); err != nil {
		t.Skipf("junction creation unavailable: %v (%s)", err, strings.TrimSpace(string(output)))
	}
	defer func() { _ = os.Remove(link) }()

	dir, err := OpenBoundDir(boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	file, err := OpenBoundDirFile(dir, filepath.Join(escapeTestLink, "f"), escapeTestBound)
	if err != nil {
		t.Fatalf("a junction staying inside the bound directory must be admitted: %v", err)
	}
	defer func() { _ = file.Close() }()

	data, truncated, err := file.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll through the in-bound junction: %v", err)
	}
	if !bytes.Equal(data, []byte("inside")) {
		t.Fatalf("the in-bound junction must read the bound bytes, got %q", data)
	}
	if truncated {
		t.Fatal("a small read must not be truncated")
	}

	if err := file.Close(); err != nil {
		t.Fatalf("Close the artifact: %v", err)
	}
	if err := dir.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}
