//go:build !windows && !darwin

package adapters

import (
	"os"
	"path/filepath"
	"testing"
)

// P53 FIX-GAP-77-SYMLINK-ESCAPE: the shared payload for the middle-symlink red
// baseline and for the permanent regression that replaced it.
//
// It sits in its own file, tagged by platform only, so that both consumers compile
// in the default build: the red variant in bounddir_escape_redbaseline_test.go
// (armed by PRFRAIL_P53_RED_BASELINE=1, skipped otherwise) has to build the exact
// same tree and compare the exact same canaries as the green regression, otherwise
// a mutation run could be judged against a different payload than the one it
// mutated.
//
// Design authority: docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/01-arch/
// design-r2.md section 5 (red/green criteria) and section 6 (regression list). The
// original breach is recorded in
// docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-attack/
// t51-1-linux-escape.md.

const (
	// gap77OutsideCanary is the outside-the-bound content an escape would expose.
	// It is exactly 23 bytes and is compared byte for byte.
	gap77OutsideCanary = "SECRET-OUTSIDE-BOUNDARY"
	// gap77InsideCanary is the in-bound content the positive controls read. It is a
	// different length, so a swapped expectation cannot compare equal by accident.
	gap77InsideCanary = "IN-SIDE-BOUNDARY"

	// gap77RegressionMessage is the single named failure the red baseline raises when
	// the escape really happens: only this exact text means "the vulnerability
	// reproduced", never a setup failure.
	gap77RegressionMessage = "GAP-77 regression: middle symlink escape succeeded and read outside the bound directory"

	// gap77LeafName is the outside-the-bound secret the middle symlink reaches.
	gap77LeafName = "secret"
	// gap77LinkName is the symlink planted inside the bound directory; it is the
	// middle component of every attack spelling.
	gap77LinkName = "link"
)

// gap77Tree is the shared payload: the bound directory D holds the symlink
// "link" -> O, the outside directory O holds the 23 byte secret, and D also holds
// in-bound artifacts for the positive controls (one single-segment name and one
// three-level name).
type gap77Tree struct {
	boundPath   string
	outsidePath string
}

// buildGAP77Tree lays out the payload under a fresh temporary root and reports the
// two directories the test needs. Every failure is fatal here: a half-built tree
// would let the red baseline fail with a setup error instead of a verdict.
func buildGAP77Tree(t *testing.T) gap77Tree {
	t.Helper()

	root := t.TempDir()
	tree := gap77Tree{
		boundPath:   filepath.Join(root, "D"),
		outsidePath: filepath.Join(root, "O"),
	}
	if err := os.Mkdir(tree.boundPath, 0o755); err != nil {
		t.Fatalf("create the bound directory: %v", err)
	}
	if err := os.Mkdir(tree.outsidePath, 0o755); err != nil {
		t.Fatalf("create the outside directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tree.boundPath, "original"), []byte(gap77InsideCanary), 0o644); err != nil {
		t.Fatalf("write the in-bound file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tree.outsidePath, gap77LeafName), []byte(gap77OutsideCanary), 0o644); err != nil {
		t.Fatalf("write the outside secret: %v", err)
	}
	nested := filepath.Join(tree.boundPath, "sub", "nested")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("create the in-bound multi-level directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nested, "inside"), []byte(gap77InsideCanary), 0o644); err != nil {
		t.Fatalf("write the in-bound multi-level file: %v", err)
	}
	if err := os.Symlink(tree.outsidePath, filepath.Join(tree.boundPath, gap77LinkName)); err != nil {
		t.Fatalf("plant the middle symlink: %v", err)
	}
	return tree
}
