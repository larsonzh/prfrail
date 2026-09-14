//go:build !windows

package adapters

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/tickets"
)

func TestAgentRunnerReplayRootUnixRejectsSymlinkedTerminalsDirectory(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	replayRoot := filepath.Join(runRoot, "agent-runner-replay")
	if err := os.MkdirAll(replayRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(replayRoot, "terminals")); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error for symlinked terminals, got %v", err)
	}
}

func TestAgentRunnerReplayRootUnixRejectsRuntimeTerminalsSwap(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	request := replayStoreRequestRecord(t, "request-runtime-terminal-swap")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	terminals := filepath.Join(runRoot, "agent-runner-replay", "terminals")
	if err := os.RemoveAll(terminals); err != nil {
		t.Fatal(err)
	}
	external := t.TempDir()
	if err := os.Symlink(external, terminals); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	intent := terminalIntentRecordFor(t, request, "completion-one")
	if _, err := store.RecordTerminalIntent(intent); !errors.Is(err, ErrAgentRunnerReplayRootUnsafePath) {
		t.Fatalf("expected unsafe path error after runtime swap, got %v", err)
	}
	if entries, err := os.ReadDir(external); err != nil || len(entries) != 0 {
		t.Fatalf("external target must stay empty: entries=%v err=%v", entries, err)
	}
}

func TestAgentRunnerReplayStoreUnixTerminalIntentConvergence(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	request := replayStoreRequestRecord(t, "request-terminal-convergence")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	intent := terminalIntentRecordFor(t, request, "completion-one")
	path, err := store.terminalIntentPath("request-terminal-convergence")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(filepath.Dir(path), "missing-target"), path); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	if _, err := store.RecordTerminalIntent(intent); !errors.Is(err, ErrAgentRunnerReplayStoreConvergence) {
		t.Fatalf("expected convergence failure for dangling terminal slot, got %v", err)
	}
}

func TestAgentRunnerReplayStoreUnixTerminalSyncFailureFailsClosedWithoutRollback(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	if store.PublicationDurability() != PublishDurabilityProven {
		t.Skip("native platform cannot prove replay publication durability")
	}
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	intent := terminalIntentRecordFor(t, request, "completion-one")
	sentinel := errors.New("sync failure")
	originalSync := replayStoreSyncParentDirectory
	replayStoreSyncParentDirectory = func(string) error {
		return sentinel
	}
	t.Cleanup(func() {
		replayStoreSyncParentDirectory = originalSync
	})

	if _, err := store.RecordTerminalIntent(intent); !errors.Is(err, sentinel) {
		t.Fatalf("expected sync failure, got %v", err)
	}
	replayStoreSyncParentDirectory = originalSync
	path, err := store.terminalIntentPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("published terminal intent must remain visible after sync failure, stat error = %v", err)
	}
	restarted := replayStoreMustNew(t, store.root)
	decision, err := restarted.RecordTerminalIntent(intent)
	if err != nil || decision != AgentRunnerTerminalIntentDecisionPresent {
		t.Fatalf("reopened store must resolve the retained intent as present, got decision=%v err=%v", decision, err)
	}
}

func TestAgentRunnerTerminalPublisherUnixProvenChainCompletes(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	result, err := publisher.PublishTerminal(context.Background(), outcome)
	if err != nil || result.Replayed {
		t.Fatalf("expected a fresh proven publication, got replayed=%v err=%v", result.Replayed, err)
	}
	intentPath, err := publisher.Store.terminalIntentPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(intentPath); err != nil {
		t.Fatalf("terminal intent must be on disk: %v", err)
	}
	closurePath, err := publisher.Store.terminalClosurePath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(closurePath); err != nil {
		t.Fatalf("terminal closure must be on disk: %v", err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); !errors.Is(err, tickets.ErrCostReservationSettled) {
		t.Fatalf("expected reservation settled, got %v", err)
	}
	if holds := admission.CostLedger.UnknownHoldReservations(); holds != 0 {
		t.Fatalf("observed settlement must not create unknown holds, got %d", holds)
	}
}
