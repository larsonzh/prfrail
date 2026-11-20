//go:build windows

package evidence

// P83 FIX ROUND 8 counterexamples (items T4a and T4b) on the Windows half.
//
// T4a: the enumeration a bound directory hands out must be the BOUND object. A
// directory whose name still resolves to that object is accepted, while a
// directory carrying the same name but a different file index (the original was
// renamed away and a new one was created at its name) is refused by the handle
// classification.
//
// T4b: an artifact that resolves to a reparse point is refused even when its name
// is unchanged.
//
// Mutation contracts are stated on each test.

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
)

// p83Fix8Junction creates one directory junction. A junction is a reparse point
// that an unprivileged account can create, so the reparse branch is reachable on a
// host without the symbolic-link privilege.
func p83Fix8Junction(t *testing.T, link, target string) {
	t.Helper()
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); err != nil {
		t.Fatalf("create the junction %s -> %s: %v (%s)", link, target, err, output)
	}
	// Runs before the TempDir cleanup (cleanup order is LIFO), so a recursive
	// removal never walks through the reparse point.
	t.Cleanup(func() { _ = os.Remove(link) })
}

// TestP83Fix8T4aEnumeratedDirectoryMustBeTheBoundObject is the T4a counterexample.
// The rename below moves the bound object away and re-creates the same NAME; the
// accepted enumeration still describes the bound object, and the replacement at
// that name is refused because its file index belongs to another object.
//
// Mutation check: dropping the volume/file-index comparison in
// confirmEnumeratedDirectory (rel ".") lets the replacement handle through and
// turns the final refusal assertion red.
func TestP83Fix8T4aEnumeratedDirectoryMustBeTheBoundObject(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	if err := os.Mkdir(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bound, "original.txt"), []byte("original"), 0o600); err != nil {
		t.Fatal(err)
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	// Accept side: the name still resolves to the original object.
	names, err := BoundDirNames(dir, "")
	if err != nil {
		t.Fatalf("an unchanged bound directory must enumerate, got %v", err)
	}
	if !slices.Contains(names, "original.txt") {
		t.Fatalf("the enumeration must describe the bound object, got %#v", names)
	}

	// Renamed away and re-created under the SAME name.
	if err := os.Rename(bound, filepath.Join(root, "bound-old")); err != nil {
		t.Fatalf("rename the bound directory away: %v", err)
	}
	if err := os.Mkdir(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(bound, "replacement.txt"), []byte("replacement"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Still accepted - and it still names the BOUND object, not the replacement.
	names, err = BoundDirNames(dir, "")
	if err != nil {
		t.Fatalf("the bound object must stay enumerable after its name moved, got %v", err)
	}
	if slices.Contains(names, "replacement.txt") {
		t.Fatalf("the enumeration described the replacement at the old name: %#v", names)
	}
	if !slices.Contains(names, "original.txt") {
		t.Fatalf("the enumeration must still describe the bound object, got %#v", names)
	}

	// Refuse side: the same NAME now addresses a different file index.
	pointer, err := syscall.UTF16PtrFromString(bound)
	if err != nil {
		t.Fatal(err)
	}
	handle, err := syscall.CreateFile(pointer, syscall.GENERIC_READ, fileInformationShareAll, nil, syscall.OPEN_EXISTING, fileFlagBackupSemantics, 0)
	if err != nil {
		t.Fatalf("open the replacement directory: %v", err)
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	refusal := dir.confirmEnumeratedDirectory(handle, ".")
	if refusal == nil {
		t.Fatal("a directory with the bound NAME but a different file index must be refused")
	}
	if !strings.Contains(refusal.Error(), "not the bound directory") {
		t.Fatalf("the refusal must name the object mismatch, got %v", refusal)
	}
}

// TestP83Fix8T4bReparseArtifactIsRefused is the T4b counterexample for the file
// half: a name that resolves to a reparse point is refused although the name
// itself never changed.
//
// A file symbolic link needs SeCreateSymbolicLinkPrivilege, so on a host without
// it the discriminating leg below is skipped and only the junction leg runs.
//
// Mutation check (file leg): dropping FILE_FLAG_OPEN_REPARSE_POINT from
// openBoundDirFile (or the reparse refusal in confirmOpenedArtifact) lets the link
// be followed and turns the assertion red.
// Mutation check (junction leg): dropping the reparse refusal in
// confirmOpenedArtifact makes the junction resolve inside the bound tree and turn
// green, which is what this leg asserts against.
func TestP83Fix8T4bReparseArtifactIsRefused(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	if err := os.Mkdir(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	t.Run("file-symlink", func(t *testing.T) {
		target := filepath.Join(bound, "target.txt")
		if err := os.WriteFile(target, []byte("target"), 0o600); err != nil {
			t.Fatal(err)
		}
		link := filepath.Join(bound, "artifact.txt")
		if symlinkErr := os.Symlink(target, link); symlinkErr != nil {
			t.Skipf("this host cannot create a symbolic link: %v", symlinkErr)
		}
		_, openErr := OpenBoundDirFile(dir, "artifact.txt", 1<<20)
		if openErr == nil {
			t.Fatal("a reparse-point artifact must be refused although its name is unchanged")
		}
		if !strings.Contains(openErr.Error(), "reparse point") {
			t.Fatalf("the refusal must name the reparse point, got %v", openErr)
		}
		// Control: the very same bytes reachable through a real name open fine.
		file, openErr := OpenBoundDirFile(dir, "target.txt", 1<<20)
		if openErr != nil {
			t.Fatalf("a regular artifact must still open, got %v", openErr)
		}
		_ = file.Close()
	})

}

// TestP83Fix8T4bEnumerationReparseIsRefused is the T4b counterexample for the
// enumeration half of the same rule. BoundDirNames documents that the handle it
// admits "must not be a reparse point", so a junction INSIDE the bound tree - the
// one shape the containment comparison cannot catch, because it resolves to a
// legal descendant - has to be refused by the reparse classification.
//
// FINDING (P83 fix round 8 baseline): this test is RED. BoundDirNames opens the
// enumerated directory without FILE_FLAG_OPEN_REPARSE_POINT, so the junction is
// followed and confirmOpenedArtifact never sees a reparse point. The claim in the
// BoundDirNames doc comment is therefore not implemented on this path; a fix that
// adds the flag (and keeps the reparse refusal in confirmOpenedArtifact) turns it
// green, and removing either half then turns it red again.
func TestP83Fix8T4bEnumerationReparseIsRefused(t *testing.T) {
	root := t.TempDir()
	bound := filepath.Join(root, "bound")
	if err := os.Mkdir(bound, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(bound, "real"), 0o755); err != nil {
		t.Fatal(err)
	}
	dir, err := OpenBoundDir(bound)
	if err != nil {
		t.Fatalf("bind the directory: %v", err)
	}
	defer func() { _ = dir.Close() }()

	p83Fix8Junction(t, filepath.Join(bound, "artifacts"), filepath.Join(bound, "real"))
	if _, namesErr := BoundDirNames(dir, "artifacts"); namesErr == nil {
		t.Fatal("an enumerated directory that is a reparse point must be refused, even when it resolves inside the bound tree")
	}
}
