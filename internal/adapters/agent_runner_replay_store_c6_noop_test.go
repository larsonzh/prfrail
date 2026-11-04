//go:build !(windows && c6slot)

package adapters

// This file pins the unactivated behaviour of PrebuildReplaySlots
// (CONTRACTS.md §7 「C6 判定路径」, ADR-013). Every build that does not compile
// agent_runner_replay_store_c6_enable.go leaves the prebuild seam nil, so the
// shared method is a pure no-op: it must not validate the declared set, must not
// touch the receiver (a nil receiver included) and must not touch the
// filesystem. The build constraint mirrors the activation file exactly, so the
// two facts are pinned together: the no-op claim is asserted only where the C6
// implementation is statically absent.

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// noopReplayTreeSnapshot renders the replay root as one sorted line per entry:
// directories by relative path, files by relative path plus the sha256 of their
// bytes. An absent root renders as the empty snapshot.
func noopReplayTreeSnapshot(t *testing.T, root string) string {
	t.Helper()
	var lines []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		if entry.IsDir() {
			lines = append(lines, "dir "+filepath.ToSlash(relative))
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		lines = append(lines, fmt.Sprintf("file %s %x", filepath.ToSlash(relative), sha256.Sum256(data)))
		return nil
	})
	if errors.Is(err, fs.ErrNotExist) {
		return ""
	}
	if err != nil {
		t.Fatalf("snapshot replay root %s: %v", root, err)
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// TestC6PrebuildNoopWithoutActivation pins that an unactivated build carries no
// C6 behaviour at all: the nil receiver is accepted (nil, no panic) and an
// ordinary store's prebuild leaves the replay root byte for byte unchanged -
// same directory entries, same per-file sha256.
func TestC6PrebuildNoopWithoutActivation(t *testing.T) {
	if replayStoreSlotEnabled {
		t.Fatal("this build must not carry the C6 slot implementation")
	}
	if replayStorePrebuildSlotsFn != nil {
		t.Fatal("this build must leave the prebuild seam nil")
	}
	for _, requestIDs := range [][]string{nil, {"request-one", "request-two"}} {
		if err := (*AgentRunnerReplayStore)(nil).PrebuildReplaySlots(requestIDs); err != nil {
			t.Fatalf("nil receiver PrebuildReplaySlots(%v) = %v, want nil", requestIDs, err)
		}
	}

	root := filepath.Join(t.TempDir(), "replay")
	if err := os.MkdirAll(filepath.Join(root, "requests"), 0o755); err != nil {
		t.Fatalf("seed replay root: %v", err)
	}
	seed := filepath.Join(root, "requests", "request.seed.jsonl")
	if err := os.WriteFile(seed, []byte("{\"seed\":true}\n"), 0o644); err != nil {
		t.Fatalf("seed replay file: %v", err)
	}
	store, err := newAgentRunnerReplayStoreAt(root)
	if err != nil {
		t.Fatalf("construct replay store: %v", err)
	}

	before := noopReplayTreeSnapshot(t, root)
	if err := store.PrebuildReplaySlots([]string{"request-one", "request-two"}); err != nil {
		t.Fatalf("prebuild on an unactivated build = %v, want nil", err)
	}
	if after := noopReplayTreeSnapshot(t, root); after != before {
		t.Fatalf("an unactivated prebuild must leave the replay root unchanged:\nbefore:\n%s\nafter:\n%s", before, after)
	}
	if store.slotsActive {
		t.Fatal("an unactivated prebuild must not activate the slot path")
	}
	if _, statErr := os.Stat(filepath.Join(root, "prebuild-set.jsonl")); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("an unactivated prebuild must not publish a prebuild set record: %v", statErr)
	}
}
