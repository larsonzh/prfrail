//go:build linux

package console

// P85 F02/F06/F20 (design rev 3, section 2 / test row M4): T2's publication must
// happen inside the VERIFIED directory descriptor, exactly like T3's carrier.
//
// The counterexample binds a directory, then renames it away and puts a decoy at the
// same spelling. A path-based publication (the pre-fix shape: os.CreateTemp +
// os.Link) writes into the decoy; a handle-based publication writes into the object
// the binding produced, which is the whole point of the identity chain.
//
// Mutation check: revert publishAgentRunnerPolicyFileIn to the path-based
// CreateTemp + Link form => the payload lands in the decoy directory and both
// assertions below go red.
//
// Windows note: the Windows half is deliberately path-based (Go cannot create or
// link relative to a directory handle there); those two steps are the registered S1
// user-mode path boundary, so this leg is Linux-only by design.

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestP85F02PublicationStaysInsideTheBoundDirectory(t *testing.T) {
	root := t.TempDir()
	live := filepath.Join(root, "agent-runner-policy")
	if err := os.MkdirAll(live, 0o755); err != nil {
		t.Fatalf("create the policy directory: %v", err)
	}

	bound, fileID, err := evidence.OpenBoundDirVerified(live)
	if err != nil {
		t.Fatalf("bind the policy directory: %v", err)
	}
	defer func() { _ = bound.Close() }()

	// Swap the name AFTER the binding. The handle still addresses the original
	// object; the spelling now addresses a decoy.
	retired := live + "-retired"
	if err := os.Rename(live, retired); err != nil {
		t.Fatalf("rename the bound directory away: %v", err)
	}
	if err := os.MkdirAll(live, 0o755); err != nil {
		t.Fatalf("create the decoy directory: %v", err)
	}

	payload := []byte("frozen policy payload")
	if err := publishAgentRunnerPolicyFileIn(bound, fileID, live, "artifact.json", payload, "frozen policy artifact"); err != nil {
		t.Fatalf("the publication must succeed inside the bound handle: %v", err)
	}

	published := filepath.Join(retired, "artifact.json")
	data, err := os.ReadFile(published)
	if err != nil {
		t.Fatalf("the payload must land in the object the binding produced: %v", err)
	}
	if string(data) != string(payload) {
		t.Fatalf("the published payload must be byte-identical, got %q want %q", data, payload)
	}
	if _, err := os.Stat(filepath.Join(live, "artifact.json")); !os.IsNotExist(err) {
		t.Fatalf("nothing may be written through the replaced spelling, stat err=%v", err)
	}
	// The decoy must also carry no leftover temporary entry: the whole publication
	// happens inside the bound descriptor.
	entries, err := os.ReadDir(live)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		t.Fatalf("the replaced spelling must receive no entries, found %v", names)
	}
}
