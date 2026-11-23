//go:build windows

package adapters

// P85 R6 (test round 6) / review r3 finding 3: the T3 executable forensics.
//
// openInstallSourceLeaf used to bind the executable's directory with
// evidence.OpenBoundDirNoFollow, whose Windows half falls through to a plain
// OpenBoundDir: the kernel followed a reparse point in the chain and the digest was
// then taken from whatever object the name had been redirected to. The binding is now
// evidence.OpenBoundDirVerified, whose Windows half walks the root chain and refuses
// every reparse component, and the identity it returns belongs to the SAME resolution
// that produced the handle.
//
// The counterexample is an intermediate junction: `...\linked\bin\agent-stub.exe`
// where `linked` is a junction. Before the fix the leaf was opened through the
// junction; now the bind refuses and the refusal names the offending component, i.e.
// it carries the directory identity the caller has to fix.
//
// Mutation check: going back to evidence.OpenBoundDirNoFollow (or any binding that
// lets the kernel follow a component) opens the leaf through the junction and turns
// this test red.

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// p85R6Junction creates one directory junction, or skips the caller when this host
// cannot create one without elevation.
func p85R6Junction(t *testing.T, link, target string) {
	t.Helper()
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("create the junction target %s: %v", target, err)
	}
	if output, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); err != nil {
		t.Skipf("this host cannot create the junction %s -> %s without elevation: %v (%s)", link, target, err, strings.TrimSpace(string(output)))
	}
	// Registered AFTER t.TempDir's cleanup so it runs first (LIFO): the recursive
	// removal must never walk through the reparse point.
	t.Cleanup(func() { _ = os.Remove(link) })
}

// TestP85R6WindowsCarrierPublishRecheckRefusesAReparse carries the executable half of
// the T3 binding (the carrier-publication half lives in internal/console, which owns
// the publication helper).
func TestP85R6WindowsCarrierPublishRecheckRefusesAReparse(t *testing.T) {
	root := t.TempDir()
	realDir := filepath.Join(root, "real")
	binDir := filepath.Join(realDir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create the install-source directory: %v", err)
	}
	executable := filepath.Join(binDir, "agent-stub.exe")
	if err := os.WriteFile(executable, []byte("MZ p85-r6-forensics"), 0o700); err != nil {
		t.Fatalf("write the executable: %v", err)
	}

	// Positive control: the same leaf, reached without any reparse component, binds.
	bound, err := openInstallSourceLeaf(executable)
	if err != nil {
		t.Skipf("this host's temporary root is not link-free, so the positive control cannot be established: %v", err)
	}
	if err := bound.Close(); err != nil {
		t.Fatalf("close the bound executable: %v", err)
	}

	t.Run("an-intermediate-junction-is-refused-with-its-component", func(t *testing.T) {
		link := filepath.Join(root, "linked")
		p85R6Junction(t, link, realDir)

		throughLink := filepath.Join(link, "bin", "agent-stub.exe")
		boundLeaf, err := openInstallSourceLeaf(throughLink)
		if err == nil {
			_ = boundLeaf.Close()
			t.Fatal("an executable reached through an intermediate junction must not bind")
		}
		if !errors.Is(err, ErrAIInstallSourceProof) {
			t.Fatalf("the refusal must be the install-source sentinel, got %v", err)
		}
		if !strings.Contains(err.Error(), "is a reparse point") {
			t.Fatalf("the refusal must be the root-chain classification, got %v", err)
		}
		// The refusal has to carry the directory identity that failed, not just the
		// leaf: an operator needs to know which component to fix.
		if component := filepath.Base(link); !strings.Contains(err.Error(), component) {
			t.Fatalf("the refusal must name the reparse component %q, got %v", component, err)
		}

		// Removing the junction restores the acceptance of the same target, so the
		// refusal is attributable to the junction and not to the fixture.
		if err := os.Remove(link); err != nil {
			t.Fatalf("remove the junction: %v", err)
		}
		restored, err := openInstallSourceLeaf(executable)
		if err != nil {
			t.Fatalf("the same leaf must bind once the junction is gone, got %v", err)
		}
		if err := restored.Close(); err != nil {
			t.Fatalf("close the re-bound executable: %v", err)
		}
	})
}
