//go:build windows

package console

// P85 R6 (test round 6) / review r3 finding 3, console half: the carrier
// publication re-check.
//
// aiInstallSourceAssertDirectoryIdentity used to compare the directory by name with
// evidence.FileIDAt, whose stat FOLLOWS a reparse point. A chain that became a
// junction after the binding therefore still reported the approved identity, and the
// re-check approved a directory the kernel would walk differently. It now re-BINDS
// the name through the same root-chain walk the initial binding uses.
//
// The counterexample is a junction that resolves to the very object the binding
// produced (rename the directory away, put a junction at its old name). Only the
// re-bind can tell the difference, because both the name-based comparison and the
// handle comparison land on the approved object.
//
// Mutation check (design rev 4 section 6, mutation 3): replacing the re-bind with
// evidence.FileIDAt makes the re-check approve, the publication proceed, and this
// test red.

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// TestP85R6WindowsCarrierPublishRecheckRefusesAReparse carries the publication half
// of the T3 binding (the executable-forensics half lives in internal/adapters, which
// owns openInstallSourceLeaf).
func TestP85R6WindowsCarrierPublishRecheckRefusesAReparse(t *testing.T) {
	root := t.TempDir()
	payload := []byte("{\"carrier\":\"p85-r6\"}\n")

	t.Run("an-untouched-directory-publishes-the-carrier", func(t *testing.T) {
		directory := filepath.Join(root, "publish-ok")
		if err := os.Mkdir(directory, 0o755); err != nil {
			t.Fatalf("create the publication directory: %v", err)
		}
		bound, fileID, err := evidence.OpenBoundDirVerified(directory)
		if err != nil {
			t.Skipf("this host's temporary root is not link-free, so the binding cannot be established: %v", err)
		}
		defer func() { _ = bound.Close() }()

		target := filepath.Join(directory, "install-source-proof.json")
		if err := writeAIInstallSourceProofFileIn(bound, fileID, target, payload); err != nil {
			t.Fatalf("an untouched directory must accept the publication, got %v", err)
		}
		published, err := os.ReadFile(target)
		if err != nil {
			t.Fatalf("read the published carrier: %v", err)
		}
		if string(published) != string(payload) {
			t.Fatalf("the published carrier must be the payload, got %q", published)
		}
	})

	t.Run("a-junction-inserted-after-the-binding-refuses-the-recheck", func(t *testing.T) {
		directory := filepath.Join(root, "publish-link")
		if err := os.Mkdir(directory, 0o755); err != nil {
			t.Fatalf("create the publication directory: %v", err)
		}
		bound, fileID, err := evidence.OpenBoundDirVerified(directory)
		if err != nil {
			t.Skipf("this host's temporary root is not link-free, so the binding cannot be established: %v", err)
		}
		defer func() { _ = bound.Close() }()

		renamed := filepath.Join(root, "publish-link-renamed")
		if err := os.Rename(directory, renamed); err != nil {
			t.Fatalf("rename the publication directory away: %v", err)
		}
		if output, err := exec.Command("cmd", "/c", "mklink", "/J", directory, renamed).CombinedOutput(); err != nil {
			t.Skipf("this host cannot create the junction %s -> %s without elevation: %v (%s)", directory, renamed, err, strings.TrimSpace(string(output)))
		}
		t.Cleanup(func() { _ = os.Remove(directory) })

		// The shape the pre-fix re-check accepted: the name-based comparison reports
		// the bound identity because the junction is followed to the same object.
		named, err := evidence.FileIDAt(directory)
		if err != nil {
			t.Fatalf("identify the junction by name: %v", err)
		}
		if named != fileID {
			t.Fatalf("the fixture must make the name-based comparison agree with the binding (%s != %s), otherwise this leg does not discriminate", named, fileID)
		}

		target := filepath.Join(directory, "install-source-proof.json")
		err = writeAIInstallSourceProofFileIn(bound, fileID, target, payload)
		if err == nil {
			t.Fatal("a directory whose name became a junction must refuse the publication, even when it resolves to the bound object")
		}
		if !errors.Is(err, ErrAIInstallSourceUnproven) {
			t.Fatalf("the refusal must be the install-source sentinel, got %v", err)
		}
		if !strings.Contains(err.Error(), "re-bound without following a link") {
			t.Fatalf("the refusal must come from the root-chain re-bind, got %v", err)
		}
		// A refused publication leaves no carrier, neither through the junction nor
		// behind it.
		for _, path := range []string{target, filepath.Join(renamed, "install-source-proof.json")} {
			if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatalf("a refused publication must leave no carrier at %s (stat error %v)", path, statErr)
			}
		}
	})
}
