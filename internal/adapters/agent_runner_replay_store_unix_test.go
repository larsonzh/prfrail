//go:build !windows

package adapters

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestAgentRunnerReplayStoreConstructorBootstrapsAndSyncsEachParentEntry(t *testing.T) {
	root := t.TempDir()
	originalSync := replayStoreSyncParentDirectory
	var synced []string
	replayStoreSyncParentDirectory = func(path string) error {
		synced = append(synced, path)
		return nil
	}
	t.Cleanup(func() {
		replayStoreSyncParentDirectory = originalSync
	})

	if _, err := newAgentRunnerReplayStoreAt(root); err != nil {
		t.Fatal(err)
	}
	if len(synced) != 4 || synced[0] != root || synced[1] != root || synced[2] != root || synced[3] != root {
		t.Fatalf("bootstrap sync order = %v, want four syncs of %s", synced, root)
	}
	for _, directory := range []string{"requests", "completions", "launches", "terminals"} {
		info, err := os.Stat(filepath.Join(root, directory))
		if err != nil || !info.IsDir() {
			t.Fatalf("bootstrap directory %s unavailable: info=%v err=%v", directory, info, err)
		}
	}
}

func TestAgentRunnerReplayStoreConstructorFailsClosedWhenBootstrapSyncFails(t *testing.T) {
	root := t.TempDir()
	sentinel := errors.New("bootstrap sync failure")
	originalSync := replayStoreSyncParentDirectory
	replayStoreSyncParentDirectory = func(string) error {
		return sentinel
	}
	t.Cleanup(func() {
		replayStoreSyncParentDirectory = originalSync
	})

	if _, err := newAgentRunnerReplayStoreAt(root); !errors.Is(err, sentinel) {
		t.Fatalf("expected bootstrap sync failure, got %v", err)
	}
}
