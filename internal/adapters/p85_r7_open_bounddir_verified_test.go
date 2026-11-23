//go:build windows || linux

package adapters

// P85 R7 (mutation coverage, review r6 finding): openInstallSourceLeaf must bind
// its directory through evidence.OpenBoundDirVerified. The mutant is a plain open
// (whose resolution follows a link) or a name-level identity check
// (evidence.FileIDAt, whose stat follows a link too) - both accept the linked chain
// below and turn this leg red.
//
// The counterexample uses a link that resolves to the SAME object the direct path
// reaches, so an identity comparison alone cannot refuse it: only a binding that
// refuses a link component can.

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85R7CreateDirectoryLink creates a directory link that the root-chain walk must
// refuse: a junction on Windows (no elevation needed) and a symbolic link on Linux.
// The caller skips when the host cannot create one.
func p85R7CreateDirectoryLink(t *testing.T, link, target string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		if output, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); err != nil {
			t.Skipf("this host cannot create the junction %s -> %s: %v (%s)", link, target, err, strings.TrimSpace(string(output)))
		}
		// Registered after t.TempDir's cleanup so it runs first (LIFO): the recursive
		// removal must never walk through the reparse point.
		t.Cleanup(func() { _ = os.Remove(link) })
		return
	}
	if err := os.Symlink(target, link); err != nil {
		t.Skipf("this host cannot create the symbolic link %s -> %s: %v", link, target, err)
	}
}

// TestP85R7OpenInstallSourceLeafRefusesALinkedChain binds the same executable twice:
// once as itself (the positive control) and once through a link directory.
func TestP85R7OpenInstallSourceLeafRefusesALinkedChain(t *testing.T) {
	root := t.TempDir()
	realDir := filepath.Join(root, "real")
	binDir := filepath.Join(realDir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create the install-source directory: %v", err)
	}
	executable := filepath.Join(binDir, "agent-stub.exe")
	if err := os.WriteFile(executable, []byte("MZ p85-r7"), 0o700); err != nil {
		t.Fatalf("write the executable: %v", err)
	}

	// Positive control: the same leaf reached without any link component binds, so
	// this host can run the counterexample at all.
	control, err := openInstallSourceLeaf(executable)
	if err != nil {
		t.Skipf("this host's temporary root is not link-free, so the positive control cannot be established: %v", err)
	}
	controlIdentity, err := control.FileID()
	if err != nil {
		t.Fatalf("identify the bound executable: %v", err)
	}
	if err := control.Close(); err != nil {
		t.Fatalf("close the bound executable: %v", err)
	}

	link := filepath.Join(root, "linked")
	p85R7CreateDirectoryLink(t, link, realDir)
	throughLink := filepath.Join(link, "bin", "agent-stub.exe")

	// A name-level identity check follows the link and reports the same object as the
	// direct path, so identity alone cannot refuse this chain.
	followed, err := evidence.FileIDAt(throughLink)
	if err != nil {
		t.Fatalf("identify the linked name: %v", err)
	}
	if string(followed) != controlIdentity {
		t.Fatalf("the fixture must make the linked name resolve to the linked-to object (%s != %s)", followed, controlIdentity)
	}

	leaf, err := openInstallSourceLeaf(throughLink)
	if err == nil {
		_ = leaf.Close()
		t.Fatal("an executable reached through a symbolic-link component must not bind")
	}
	if !errors.Is(err, ErrAIInstallSourceProof) {
		t.Fatalf("the refusal must be the install-source sentinel, got %v", err)
	}

	// Removing the link restores the acceptance of the same target, so the refusal is
	// attributable to the link and not to the fixture.
	if err := os.Remove(link); err != nil {
		t.Fatalf("remove the directory link: %v", err)
	}
	restored, err := openInstallSourceLeaf(executable)
	if err != nil {
		t.Fatalf("the same leaf must bind once the link is gone, got %v", err)
	}
	if err := restored.Close(); err != nil {
		t.Fatalf("close the re-bound executable: %v", err)
	}
}
