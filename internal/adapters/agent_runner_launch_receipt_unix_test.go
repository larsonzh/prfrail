//go:build !windows

package adapters

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestAgentRunnerReplayRootUnixRejectsSymlinkedLaunchesDirectory(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	replayRoot := filepath.Join(runRoot, "agent-runner-replay")
	if err := os.MkdirAll(replayRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(replayRoot, "launches")); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error for symlinked launches, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixRejectsRuntimeLaunchesSwap(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	record := replayStoreRequestRecord(t, "request-runtime-launch-swap")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	launches := filepath.Join(runRoot, "agent-runner-replay", "launches")
	if err := os.RemoveAll(launches); err != nil {
		t.Fatal(err)
	}
	external := t.TempDir()
	if err := os.Symlink(external, launches); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	receipt := launchReceiptRecordFor(t, record)
	if _, err := store.RecordLaunchReceipt(receipt); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error after runtime swap, got %v", err)
	}
	if entries, err := os.ReadDir(external); err != nil || len(entries) != 0 {
		t.Fatalf("external target must stay empty: entries=%v err=%v", entries, err)
	}
}

func TestAgentRunnerReplayStoreUnixLaunchReceiptConvergence(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	record := replayStoreRequestRecord(t, "request-launch-convergence")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	receipt := launchReceiptRecordFor(t, record)
	path, err := store.launchReceiptPath("request-launch-convergence")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(filepath.Dir(path), "missing-target"), path); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := store.RecordLaunchReceipt(receipt); !errors.Is(err, ErrAgentRunnerReplayStoreConvergence) {
		t.Fatalf("expected convergence failure for dangling launch slot, got %v", err)
	}
}

func TestAgentRunnerReplayStoreUnixLaunchReceiptSyncFailureFailsClosedWithoutRollback(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	if store.PublicationDurability() != PublishDurabilityProven {
		t.Skip("native platform cannot prove replay publication durability")
	}
	record := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	receipt := launchReceiptRecordFor(t, record)
	sentinel := errors.New("sync failure")
	originalSync := replayStoreSyncParentDirectory
	replayStoreSyncParentDirectory = func(string) error {
		return sentinel
	}
	t.Cleanup(func() {
		replayStoreSyncParentDirectory = originalSync
	})

	if _, err := store.RecordLaunchReceipt(receipt); !errors.Is(err, sentinel) {
		t.Fatalf("expected sync failure, got %v", err)
	}
	replayStoreSyncParentDirectory = originalSync
	path, err := store.launchReceiptPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("published launch receipt must remain visible after sync failure, stat error = %v", err)
	}
	restarted := replayStoreMustNew(t, store.root)
	decision, err := restarted.RecordLaunchReceipt(receipt)
	if err != nil || decision != AgentRunnerLaunchDecisionAlreadyLaunched {
		t.Fatalf("reopened store must resolve the retained intent as replay, got decision=%v err=%v", decision, err)
	}
}
