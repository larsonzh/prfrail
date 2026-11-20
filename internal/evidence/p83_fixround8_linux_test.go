//go:build linux

package evidence

// P83 FIX ROUND 8 counterexample (item T4c): BoundDirNames reads the entry names
// from the bound DESCRIPTOR, never from a second resolution of the caller's path.
//
// The Ubuntu leg of the CI executes this file; the Windows host cannot run it.
//
// Mutation check: replacing the descriptor enumeration with a path glob under
// dir.path makes the replacement directory's entry show up after the rename below
// and turns the assertion red.

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestP83Fix8T4cBoundDirNamesReadsFromDescriptor(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	if err := os.MkdirAll(filepath.Join(bound, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"b.txt", "a.txt"} {
		if err := os.WriteFile(filepath.Join(bound, "sub", name), []byte(name), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	// The object moves away and the old NAME is taken by another directory.
	if err := os.Rename(bound, filepath.Join(root, "bound-moved")); err != nil {
		t.Fatalf("rename the bound directory away: %v", err)
	}
	if err := os.Mkdir(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bound, "replacement.txt"), []byte("replacement"), 0o600); err != nil {
		t.Fatal(err)
	}

	names, err := BoundDirNames(dir, "sub")
	if err != nil {
		t.Fatalf("the bound descriptor must still enumerate, got %v", err)
	}
	if want := []string{"a.txt", "b.txt"}; !slices.Equal(names, want) {
		t.Fatalf("the enumeration must come from the descriptor: got %#v want %#v", names, want)
	}

	names, err = BoundDirNames(dir, "")
	if err != nil {
		t.Fatalf("the bound object must enumerate, got %v", err)
	}
	if slices.Contains(names, "replacement.txt") {
		t.Fatalf("the enumeration described the replacement at the old name: %#v", names)
	}
	if !slices.Contains(names, "sub") {
		t.Fatalf("the enumeration must describe the bound object, got %#v", names)
	}
}
