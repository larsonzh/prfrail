//go:build windows

package console

// P85 design rev 4, section 3 (publication re-check) / test row M14: the check that
// runs before and after a path-based publication must be a RE-BIND through the same
// root-chain walk, not a name-based identity comparison.
//
// The discriminating file-system shape is a junction that resolves to the very same
// object the binding produced (rename the directory away, put a junction at its old
// name pointing at the new location). A name-based comparison reports the bound
// identity - it follows the junction and lands on the same object - so it approves;
// the root chain sees a reparse point in the last component and refuses.
//
// Mutation check: replacing the re-bind with evidence.FileIDAt in
// agentRunnerAssertDirectoryIdentity makes this leg approve and turns it red.

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// TestP85R4M14PublishRecheckRefusesASameObjectJunction is the console half of M14.
func TestP85R4M14PublishRecheckRefusesASameObjectJunction(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "agent-runner-policy")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatalf("create the policy directory: %v", err)
	}

	bound, fileID, err := evidence.OpenBoundDirVerified(directory)
	if err != nil {
		t.Skipf("this host's temporary root is not link-free, so the binding half of the fixture cannot be established: %v", err)
	}
	defer func() { _ = bound.Close() }()

	// Positive control: the untouched directory passes the re-check.
	if err := agentRunnerAssertDirectoryIdentity(bound, fileID, directory, "frozen policy artifact"); err != nil {
		t.Fatalf("an untouched directory must pass the publication re-check, got %v", err)
	}

	renamed := filepath.Join(root, "agent-runner-policy-renamed")
	if err := os.Rename(directory, renamed); err != nil {
		t.Fatalf("rename the policy directory away: %v", err)
	}
	if output, err := exec.Command("cmd", "/c", "mklink", "/J", directory, renamed).CombinedOutput(); err != nil {
		t.Skipf("this host cannot create the junction %s -> %s without elevation: %v (%s)", directory, renamed, err, strings.TrimSpace(string(output)))
	}
	t.Cleanup(func() { _ = os.Remove(directory) })

	// The shape the pre-fix re-check accepted: the NAME still reports the bound
	// identity, because the junction is followed to the same object.
	named, err := evidence.FileIDAt(directory)
	if err != nil {
		t.Fatalf("identify the junction by name: %v", err)
	}
	if named != fileID {
		t.Fatalf("the fixture must make the name-based comparison agree with the binding (%s != %s), otherwise this leg does not discriminate", named, fileID)
	}

	err = agentRunnerAssertDirectoryIdentity(bound, fileID, directory, "frozen policy artifact")
	if err == nil {
		t.Fatal("a directory whose name is now a junction must refuse the publication re-check even though it resolves to the bound object")
	}
	if !strings.Contains(err.Error(), "root chain") {
		t.Fatalf("the refusal must come from the root-chain re-bind, got %v", err)
	}
}
