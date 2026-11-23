//go:build !darwin

package evidence

// P52-S1 ACCEPTANCE-SECURITY-COUNTEREXAMPLES: directory-identity counterexamples.
//
// Two refusals on the bound-directory seam: a directory swapped after the bind is
// not the bound object (boundDirVerify), and an artifact name that escapes the
// bound directory is never opened (cleanBoundRelativePath). Both are package
// internal seams, so the tests live in the package rather than behind the API.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestBoundDirVerifyRefusesSwappedDirectory pins the identity check: the open
// handle keeps addressing the original object, so a rename away plus a fresh
// directory at the same path is reported as "is not" that path.
func TestBoundDirVerifyRefusesSwappedDirectory(t *testing.T) {
	root := t.TempDir()
	boundPath := filepath.Join(root, "ev")
	if err := os.Mkdir(boundPath, 0o755); err != nil {
		t.Fatalf("create the bound directory: %v", err)
	}

	dir, err := OpenBoundDir(boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()
	boundID, err := dir.FileID()
	if err != nil {
		t.Fatalf("identify the bound directory: %v", err)
	}

	if err := os.Rename(boundPath, boundPath+".bak"); err != nil {
		t.Fatalf("rename the bound directory away (the handle must stay renamable): %v", err)
	}
	if err := os.Mkdir(boundPath, 0o755); err != nil {
		t.Fatalf("create the injected directory at the same path: %v", err)
	}

	// P85 F02/F06/F20: the expected value is the identity the binding produced - the
	// object the caller approved - and the check compares the handle with it. The
	// injected directory at the same spelling is a DIFFERENT object, so it must be
	// refused.
	injectedID, err := FileIDAt(boundPath)
	if err != nil {
		t.Fatalf("identify the injected directory: %v", err)
	}
	swapped := boundDirVerify(dir, injectedID)
	if swapped == nil {
		t.Fatal("a directory swapped in after the bind must be refused")
	}
	if !strings.Contains(swapped.Error(), "is not") {
		t.Fatalf("the refusal must name the identity mismatch, got %v", swapped)
	}

	// Control: the very same check accepts the object the handle really holds, so
	// the refusal above is caused by the swap and not by any verification error.
	if kept := boundDirVerify(dir, boundID); kept != nil {
		t.Fatalf("the bound object must still verify against its own identity, got %v", kept)
	}

	if err := dir.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

// TestCleanBoundRelativePathRefusesEscapeShapes pins the name guard: empty,
// absolute and parent-escaping names are refused, while a benign name is cleaned
// (single dots folded, separators normalized) and accepted.
func TestCleanBoundRelativePathRefusesEscapeShapes(t *testing.T) {
	for _, rel := range []string{"", "/etc/passwd", "..", "a/../../b"} {
		clean, err := cleanBoundRelativePath(rel)
		if err == nil {
			t.Fatalf("cleanBoundRelativePath(%q) must be refused, got %q", rel, clean)
		}
	}

	clean, err := cleanBoundRelativePath("sub/./x")
	if err != nil {
		t.Fatalf("cleanBoundRelativePath(\"sub/./x\") must be accepted: %v", err)
	}
	if want := filepath.FromSlash("sub/x"); clean != want {
		t.Fatalf("cleanBoundRelativePath(\"sub/./x\") = %q, want %q", clean, want)
	}
}
