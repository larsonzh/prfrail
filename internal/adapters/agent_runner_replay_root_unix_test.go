//go:build !windows

package adapters

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestAgentRunnerReplayRootUnixRejectsSymlinkedReplayRoot(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	link := filepath.Join(runRoot, "agent-runner-replay")
	if err := os.Symlink(t.TempDir(), link); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixRejectsSymlinkedRunRoot(t *testing.T) {
	realRoot := replayRootTestRunRoot(t)
	linkRoot := filepath.Join(t.TempDir(), "run-link")
	if err := os.Symlink(realRoot, linkRoot); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(linkRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixRejectsSymlinkedRequestsDirectory(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	replayRoot := filepath.Join(runRoot, "agent-runner-replay")
	if err := os.MkdirAll(replayRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(replayRoot, "requests")); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixRejectsSymlinkedEventsDirectory(t *testing.T) {
	runRoot := t.TempDir()
	if err := os.Symlink(t.TempDir(), filepath.Join(runRoot, "events")); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixRejectsSymlinkedOwnershipMarker(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	replayRootTestStore(t, runRoot, "run-one")
	markerPath := filepath.Join(runRoot, "agent-runner-replay", "ownership.json")
	if err := os.Remove(markerPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(runRoot, "missing-marker-target"), markerPath); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for symlinked ownership marker, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixFailsClosedWhenProductionSyncFails(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	sentinel := errors.New("production sync failure")
	originalSync := replayStoreSyncParentDirectory
	replayStoreSyncParentDirectory = func(string) error {
		return sentinel
	}
	t.Cleanup(func() {
		replayStoreSyncParentDirectory = originalSync
	})
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, sentinel) {
		t.Fatalf("expected sync failure, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixSyncsCreatedParents(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	replayRoot := filepath.Join(runRoot, "agent-runner-replay")
	originalSync := replayStoreSyncParentDirectory
	var synced []string
	replayStoreSyncParentDirectory = func(path string) error {
		synced = append(synced, path)
		return nil
	}
	t.Cleanup(func() {
		replayStoreSyncParentDirectory = originalSync
	})
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); err != nil {
		t.Fatal(err)
	}
	if len(synced) < 3 {
		t.Fatalf("expected at least three parent syncs, got %v", synced)
	}
	if synced[0] != runRoot {
		t.Fatalf("first sync = %s, want run root %s (order %v)", synced[0], runRoot, synced)
	}
	for index := 1; index < len(synced); index++ {
		if synced[index] != replayRoot {
			t.Fatalf("sync %d = %s, want replay root %s (order %v)", index, synced[index], replayRoot, synced)
		}
	}
}

func TestAgentRunnerReplayRootUnixLeavesNoOwnershipWhenSyncFails(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	replayRoot := filepath.Join(runRoot, "agent-runner-replay")
	markerPath := filepath.Join(replayRoot, "ownership.json")
	sentinel := errors.New("sync failure")
	originalSync := replayStoreSyncParentDirectory
	replayStoreSyncParentDirectory = func(path string) error {
		if path == replayRoot {
			return sentinel
		}
		return nil
	}
	t.Cleanup(func() {
		replayStoreSyncParentDirectory = originalSync
	})
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, sentinel) {
		t.Fatalf("expected sync failure, got %v", err)
	}
	if _, err := os.Stat(markerPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("failed construction must not leave ownership metadata, stat err=%v", err)
	}
	replayStoreSyncParentDirectory = originalSync
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); err != nil {
		t.Fatalf("reopen after sync failure must publish ownership: %v", err)
	}
	if _, err := os.Stat(markerPath); err != nil {
		t.Fatalf("ownership marker not published on reopen: %v", err)
	}
}

func TestAgentRunnerReplayRootUnixRejectsRuntimeSubdirectorySwap(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	requests := filepath.Join(runRoot, "agent-runner-replay", "requests")
	if err := os.RemoveAll(requests); err != nil {
		t.Fatal(err)
	}
	external := t.TempDir()
	if err := os.Symlink(external, requests); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	record := replayStoreRequestRecord(t, "request-runtime-swap")
	if _, err := store.RecordRequest(record); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error after runtime swap, got %v", err)
	}
	if entries, err := os.ReadDir(external); err != nil || len(entries) != 0 {
		t.Fatalf("external target must stay empty: entries=%v err=%v", entries, err)
	}
	completions := filepath.Join(runRoot, "agent-runner-replay", "completions")
	if err := os.RemoveAll(completions); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, completions); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := store.State("request-runtime-swap"); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error for swapped completions, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixRejectsSymlinkedCallerRoots(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	linkRoot := filepath.Join(t.TempDir(), "caller-link")
	if err := os.Symlink(t.TempDir(), linkRoot); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if err := store.RejectWriterRoots(linkRoot); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error for symlinked writer root, got %v", err)
	}
	if err := store.RejectProtectedRoots(linkRoot); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error for symlinked protected root, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixOwnershipConvergenceFailure(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	replayRoot := filepath.Join(runRoot, "agent-runner-replay")
	if err := os.MkdirAll(replayRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	markerPath := filepath.Join(replayRoot, "ownership.json")
	if err := os.Symlink(filepath.Join(replayRoot, "missing-marker-target"), markerPath); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	originalLoad := replayStoreLoadOwnershipRecord
	replayStoreLoadOwnershipRecord = func(string) (AgentRunnerReplayOwnershipRecord, bool, error) {
		return AgentRunnerReplayOwnershipRecord{}, false, nil
	}
	t.Cleanup(func() {
		replayStoreLoadOwnershipRecord = originalLoad
	})
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayStoreConvergence) {
		t.Fatalf("expected convergence failure, got %v", err)
	}
}
