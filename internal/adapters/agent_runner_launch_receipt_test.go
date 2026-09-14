package adapters

import (
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func launchWire(t *testing.T, value any) []byte {
	t.Helper()
	wire, err := evidence.EncodeCanonical(value)
	if err != nil {
		t.Fatal(err)
	}
	return wire
}

func launchReceiptRecordFor(t *testing.T, request AgentRunnerRequestRecord) AgentRunnerLaunchReceiptRecord {
	t.Helper()
	receipt, err := NewAgentRunnerLaunchReceiptRecord(AgentRunnerLaunchReceipt{
		LaunchID:          "launch-" + request.Request.RequestID,
		RequestID:         request.Request.RequestID,
		RequestHash:       request.RecordHash,
		RunID:             request.Request.RunID,
		TaskID:            request.Request.TaskID,
		StepID:            request.Request.StepID,
		Attempt:           request.Request.Attempt,
		AdapterID:         request.Request.AdapterID,
		ReconfirmedAt:     "2026-09-11T08:05:00Z",
		AuthorizationHash: request.Request.AuthorizationHash,
		BudgetHash:        request.Request.BudgetHash,
	})
	if err != nil {
		t.Fatal(err)
	}
	return receipt
}

func launchIdentityRecordFor(t *testing.T, receipt AgentRunnerLaunchReceiptRecord) AgentRunnerLaunchIdentityRecord {
	t.Helper()
	identity, err := NewAgentRunnerLaunchIdentityRecord(AgentRunnerLaunchIdentity{
		LaunchID:         receipt.Receipt.LaunchID,
		RequestID:        receipt.Receipt.RequestID,
		IntentRecordHash: receipt.RecordHash,
		ProcessID:        "process-stub-one",
		StartedAt:        "2026-09-11T08:05:01Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	return identity
}

func TestAgentRunnerLaunchReceiptRecordsValidateBindings(t *testing.T) {
	record := replayStoreRequestRecord(t, "request-one")
	receipt := launchReceiptRecordFor(t, record)
	if err := ValidateAgentRunnerLaunchReceiptRecord(receipt); err != nil {
		t.Fatalf("expected valid receipt, got %v", err)
	}
	if decoded, err := DecodeAgentRunnerLaunchReceiptRecord(launchWire(t, receipt)); err != nil || decoded.RecordHash != receipt.RecordHash {
		t.Fatalf("receipt round trip failed: %v", err)
	}

	badKind := receipt.Receipt
	badKind.Kind = "other-receipt"
	if _, err := NewAgentRunnerLaunchReceiptRecord(badKind); !errors.Is(err, ErrInvalidAgentRunnerLaunchReceipt) {
		t.Fatalf("expected invalid receipt for bad kind, got %v", err)
	}
	badAttempt := receipt.Receipt
	badAttempt.Attempt = 0
	if _, err := NewAgentRunnerLaunchReceiptRecord(badAttempt); !errors.Is(err, ErrInvalidAgentRunnerLaunchReceipt) {
		t.Fatalf("expected invalid receipt for zero attempt, got %v", err)
	}
	badTime := receipt.Receipt
	badTime.ReconfirmedAt = "yesterday"
	if _, err := NewAgentRunnerLaunchReceiptRecord(badTime); !errors.Is(err, ErrInvalidAgentRunnerLaunchReceipt) {
		t.Fatalf("expected invalid receipt for bad reconfirmedAt, got %v", err)
	}
	tampered := receipt
	tampered.RecordHash = queueHashOne
	if err := ValidateAgentRunnerLaunchReceiptRecord(tampered); !errors.Is(err, ErrInvalidAgentRunnerLaunchReceipt) {
		t.Fatalf("expected invalid receipt for tampered hash, got %v", err)
	}

	identity := launchIdentityRecordFor(t, receipt)
	if err := ValidateAgentRunnerLaunchIdentityRecord(identity); err != nil {
		t.Fatalf("expected valid identity, got %v", err)
	}
	if decoded, err := DecodeAgentRunnerLaunchIdentityRecord(launchWire(t, identity)); err != nil || decoded.RecordHash != identity.RecordHash {
		t.Fatalf("identity round trip failed: %v", err)
	}
	badProcess := identity.Identity
	badProcess.ProcessID = "   "
	if _, err := NewAgentRunnerLaunchIdentityRecord(badProcess); !errors.Is(err, ErrInvalidAgentRunnerLaunchIdentity) {
		t.Fatalf("expected invalid identity for empty process id, got %v", err)
	}
	badStart := identity.Identity
	badStart.StartedAt = "never"
	if _, err := NewAgentRunnerLaunchIdentityRecord(badStart); !errors.Is(err, ErrInvalidAgentRunnerLaunchIdentity) {
		t.Fatalf("expected invalid identity for bad startedAt, got %v", err)
	}
}

func TestAgentRunnerReplayStoreLaunchReceiptRequiresPersistedRequest(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	record := replayStoreRequestRecord(t, "request-one")
	receipt := launchReceiptRecordFor(t, record)
	if _, err := store.RecordLaunchReceipt(receipt); !errors.Is(err, ErrAgentRunnerRequestConflict) {
		t.Fatalf("expected request conflict without persisted request, got %v", err)
	}
	if _, found, err := store.LaunchReceipt("request-one"); err != nil || found {
		t.Fatalf("no receipt may be published without a request: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerReplayStoreLaunchReceiptSingleWinnerReplayAndConflicts(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	record := replayStoreRequestRecord(t, "request-one")
	if decision, err := store.RecordRequest(record); err != nil || decision != AgentRunnerReplayDecisionFirstDispatch {
		t.Fatalf("request publication failed: decision=%v err=%v", decision, err)
	}
	receipt := launchReceiptRecordFor(t, record)
	if decision, err := store.RecordLaunchReceipt(receipt); err != nil || decision != AgentRunnerLaunchDecisionFirstLaunch {
		t.Fatalf("expected first launch, got decision=%v err=%v", decision, err)
	}
	if decision, err := store.RecordLaunchReceipt(receipt); err != nil || decision != AgentRunnerLaunchDecisionAlreadyLaunched {
		t.Fatalf("expected already-launched replay, got decision=%v err=%v", decision, err)
	}

	changedPayload, err := NewAgentRunnerLaunchReceiptRecord(func() AgentRunnerLaunchReceipt {
		body := receipt.Receipt
		body.ReconfirmedAt = "2026-09-11T08:06:00Z"
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordLaunchReceipt(changedPayload); !errors.Is(err, ErrAgentRunnerLaunchReceiptConflict) {
		t.Fatalf("expected receipt conflict for changed payload, got %v", err)
	}
	otherLaunch, err := NewAgentRunnerLaunchReceiptRecord(func() AgentRunnerLaunchReceipt {
		body := receipt.Receipt
		body.LaunchID = "launch-other-one"
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordLaunchReceipt(otherLaunch); !errors.Is(err, ErrAgentRunnerLaunchReceiptConflict) {
		t.Fatalf("expected receipt conflict for foreign launchId, got %v", err)
	}
	foreignHash, err := NewAgentRunnerLaunchReceiptRecord(func() AgentRunnerLaunchReceipt {
		body := receipt.Receipt
		body.RequestHash = queueHashOne
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordLaunchReceipt(foreignHash); !errors.Is(err, ErrAgentRunnerRequestConflict) {
		t.Fatalf("expected request conflict for foreign requestHash, got %v", err)
	}

	stored, found, err := store.LaunchReceipt("request-one")
	if err != nil || !found {
		t.Fatalf("stored receipt read failed: found=%v err=%v", found, err)
	}
	if stored.RecordHash != receipt.RecordHash {
		t.Fatalf("stored receipt hash = %s, want %s", stored.RecordHash, receipt.RecordHash)
	}
}

func TestAgentRunnerReplayStoreLaunchReceiptUnprovenStoreBlocksPublication(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	record := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	receipt := launchReceiptRecordFor(t, record)
	store.publicationDurability = PublishDurabilityUnproven
	if _, err := store.RecordLaunchReceipt(receipt); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("expected durability gate, got %v", err)
	}
	store.publicationDurability = PublishDurabilityProven
	if decision, err := store.RecordLaunchReceipt(receipt); err != nil || decision != AgentRunnerLaunchDecisionFirstLaunch {
		t.Fatalf("expected publication after durability restore, got decision=%v err=%v", decision, err)
	}
}

func TestAgentRunnerReplayStoreLaunchReceiptRunIDBinding(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	store.runID = "run-one"
	record := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	foreign, err := NewAgentRunnerLaunchReceiptRecord(func() AgentRunnerLaunchReceipt {
		body := launchReceiptRecordFor(t, record).Receipt
		body.RunID = "run-two"
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordLaunchReceipt(foreign); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("expected run binding conflict on write, got %v", err)
	}
	path, err := store.launchReceiptPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(path, foreign); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.LaunchReceipt("request-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for foreign run binding on disk, got %v", err)
	}
}

func TestAgentRunnerReplayStoreLaunchReceiptCorruptionWithoutRequest(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	record := replayStoreRequestRecord(t, "request-one")
	receipt := launchReceiptRecordFor(t, record)
	path, err := store.launchReceiptPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(path, receipt); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordLaunchReceipt(receipt); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for receipt without request, got %v", err)
	}
}

func TestAgentRunnerReplayStoreLaunchIdentityRequiresMatchingIntent(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	record := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	receipt := launchReceiptRecordFor(t, record)
	identity := launchIdentityRecordFor(t, receipt)
	if _, err := store.RecordLaunchIdentity(identity); !errors.Is(err, ErrAgentRunnerLaunchReceiptConflict) {
		t.Fatalf("expected conflict without launch intent, got %v", err)
	}
	if _, err := store.RecordLaunchReceipt(receipt); err != nil {
		t.Fatal(err)
	}
	foreignIntent, err := NewAgentRunnerLaunchIdentityRecord(func() AgentRunnerLaunchIdentity {
		body := identity.Identity
		body.IntentRecordHash = queueHashOne
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordLaunchIdentity(foreignIntent); !errors.Is(err, ErrAgentRunnerLaunchReceiptConflict) {
		t.Fatalf("expected conflict for foreign intent hash, got %v", err)
	}
}

func TestAgentRunnerReplayStoreLaunchIdentityReplayAndConflict(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	record := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	receipt := launchReceiptRecordFor(t, record)
	if _, err := store.RecordLaunchReceipt(receipt); err != nil {
		t.Fatal(err)
	}
	identity := launchIdentityRecordFor(t, receipt)
	if replay, err := store.RecordLaunchIdentity(identity); err != nil || replay {
		t.Fatalf("expected first identity publication, got replay=%v err=%v", replay, err)
	}
	if replay, err := store.RecordLaunchIdentity(identity); err != nil || !replay {
		t.Fatalf("expected identity replay, got replay=%v err=%v", replay, err)
	}
	changed, err := NewAgentRunnerLaunchIdentityRecord(func() AgentRunnerLaunchIdentity {
		body := identity.Identity
		body.ProcessID = "process-stub-two"
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordLaunchIdentity(changed); !errors.Is(err, ErrAgentRunnerLaunchIdentityConflict) {
		t.Fatalf("expected identity conflict for changed payload, got %v", err)
	}
	stored, found, err := store.LaunchIdentity("launch-request-one")
	if err != nil || !found || stored.RecordHash != identity.RecordHash {
		t.Fatalf("stored identity read failed: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerReplayStoreLaunchIdentityUnprovenReplayOnly(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	record := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	receipt := launchReceiptRecordFor(t, record)
	if _, err := store.RecordLaunchReceipt(receipt); err != nil {
		t.Fatal(err)
	}
	store.publicationDurability = PublishDurabilityUnproven
	identity := launchIdentityRecordFor(t, receipt)
	if _, err := store.RecordLaunchIdentity(identity); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("expected durability gate for new identity, got %v", err)
	}
	store.publicationDurability = PublishDurabilityProven
	if replay, err := store.RecordLaunchIdentity(identity); err != nil || replay {
		t.Fatalf("expected publication after durability restore, got replay=%v err=%v", replay, err)
	}
	store.publicationDurability = PublishDurabilityUnproven
	if replay, err := store.RecordLaunchIdentity(identity); err != nil || !replay {
		t.Fatalf("expected unproven store to allow exact identity replay, got replay=%v err=%v", replay, err)
	}
}

func TestAgentRunnerReplayStoreLaunchMethodsFailClosedOnZeroValueStore(t *testing.T) {
	var store AgentRunnerReplayStore
	record := replayStoreRequestRecord(t, "request-one")
	receipt := launchReceiptRecordFor(t, record)
	if _, err := store.RecordLaunchReceipt(receipt); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("zero-value store must reject receipt writes, got %v", err)
	}
	if _, _, err := store.LaunchReceipt("request-one"); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("zero-value store must reject receipt reads, got %v", err)
	}
	identity := launchIdentityRecordFor(t, receipt)
	if _, err := store.RecordLaunchIdentity(identity); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("zero-value store must reject identity writes, got %v", err)
	}
	if _, _, err := store.LaunchIdentity("launch-request-one"); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("zero-value store must reject identity reads, got %v", err)
	}
}

func TestAgentRunnerReplayStoreLaunchReceiptMessagesNameTheSlot(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	record := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	receipt := launchReceiptRecordFor(t, record)
	if _, err := store.RecordLaunchReceipt(receipt); err != nil {
		t.Fatal(err)
	}
	other, err := NewAgentRunnerLaunchReceiptRecord(func() AgentRunnerLaunchReceipt {
		body := receipt.Receipt
		body.LaunchID = "launch-other-one"
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.RecordLaunchReceipt(other)
	if err == nil || !strings.Contains(err.Error(), "launch-request-one") {
		t.Fatalf("expected conflict message to name the owning launchId, got %v", err)
	}
}

func TestAgentRunnerReplayStoreLaunchReceiptCrossStoreSingleWinner(t *testing.T) {
	root := t.TempDir()
	storeA := replayStoreMustNew(t, root)
	storeB := replayStoreMustNew(t, root)
	record := replayStoreRequestRecord(t, "request-one")
	if _, err := storeA.RecordRequest(record); err != nil {
		t.Fatal(err)
	}
	receipt := launchReceiptRecordFor(t, record)

	const workers = 24
	start := make(chan struct{})
	decisions := make(chan AgentRunnerLaunchDecision, workers)
	errorsByCall := make(chan error, workers)
	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		selected := storeA
		if worker%2 == 1 {
			selected = storeB
		}
		wg.Add(1)
		go func(active *AgentRunnerReplayStore) {
			defer wg.Done()
			<-start
			decision, err := active.RecordLaunchReceipt(receipt)
			decisions <- decision
			errorsByCall <- err
		}(selected)
	}
	close(start)
	wg.Wait()
	close(decisions)
	close(errorsByCall)

	for err := range errorsByCall {
		if err != nil {
			t.Fatal(err)
		}
	}
	firstLaunches := 0
	alreadyLaunched := 0
	for decision := range decisions {
		switch decision {
		case AgentRunnerLaunchDecisionFirstLaunch:
			firstLaunches++
		case AgentRunnerLaunchDecisionAlreadyLaunched:
			alreadyLaunched++
		default:
			t.Fatalf("unexpected decision %s", decision)
		}
	}
	if firstLaunches != 1 || alreadyLaunched != workers-1 {
		t.Fatalf("launch decisions = first:%d already:%d, want first:1 already:%d", firstLaunches, alreadyLaunched, workers-1)
	}
}

func TestAgentRunnerReplayStoreLaunchReceiptCollisionClassification(t *testing.T) {
	t.Run("visible winner with same payload resolves as replay", func(t *testing.T) {
		store := replayStoreMustNew(t, t.TempDir())
		record := replayStoreRequestRecord(t, "request-one")
		if _, err := store.RecordRequest(record); err != nil {
			t.Fatal(err)
		}
		receipt := launchReceiptRecordFor(t, record)
		path, err := store.launchReceiptPath("request-one")
		if err != nil {
			t.Fatal(err)
		}
		if err := writeReplayRecordNoReplace(path, receipt); err != nil {
			t.Fatal(err)
		}
		original := replayStoreLoadLaunchReceipt
		calls := 0
		replayStoreLoadLaunchReceipt = func(activePath, requestID string) (AgentRunnerLaunchReceiptRecord, bool, error) {
			calls++
			if calls == 1 {
				// Simulate the visibility gap where the winner exists on disk but
				// the first read has not yet observed it, forcing the collision path.
				return AgentRunnerLaunchReceiptRecord{}, false, nil
			}
			return original(activePath, requestID)
		}
		t.Cleanup(func() {
			replayStoreLoadLaunchReceipt = original
		})
		decision, err := store.RecordLaunchReceipt(receipt)
		if err != nil || decision != AgentRunnerLaunchDecisionAlreadyLaunched {
			t.Fatalf("expected collision replay, got decision=%v err=%v", decision, err)
		}
	})
	t.Run("visible winner with foreign payload surfaces conflict", func(t *testing.T) {
		store := replayStoreMustNew(t, t.TempDir())
		record := replayStoreRequestRecord(t, "request-one")
		if _, err := store.RecordRequest(record); err != nil {
			t.Fatal(err)
		}
		receipt := launchReceiptRecordFor(t, record)
		foreign, err := NewAgentRunnerLaunchReceiptRecord(func() AgentRunnerLaunchReceipt {
			body := receipt.Receipt
			body.ReconfirmedAt = "2026-09-11T08:06:00Z"
			return body
		}())
		if err != nil {
			t.Fatal(err)
		}
		path, err := store.launchReceiptPath("request-one")
		if err != nil {
			t.Fatal(err)
		}
		if err := writeReplayRecordNoReplace(path, foreign); err != nil {
			t.Fatal(err)
		}
		original := replayStoreLoadLaunchReceipt
		calls := 0
		replayStoreLoadLaunchReceipt = func(activePath, requestID string) (AgentRunnerLaunchReceiptRecord, bool, error) {
			calls++
			if calls == 1 {
				return AgentRunnerLaunchReceiptRecord{}, false, nil
			}
			return original(activePath, requestID)
		}
		t.Cleanup(func() {
			replayStoreLoadLaunchReceipt = original
		})
		if _, err := store.RecordLaunchReceipt(receipt); !errors.Is(err, ErrAgentRunnerLaunchReceiptConflict) {
			t.Fatalf("expected collision conflict, got %v", err)
		}
	})
}

func TestAgentRunnerReplayStoreLaunchIdentityOrphanIsCorruption(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	record := replayStoreRequestRecord(t, "request-one")
	receipt := launchReceiptRecordFor(t, record)
	identity := launchIdentityRecordFor(t, receipt)
	path, err := store.launchIdentityPath(receipt.Receipt.LaunchID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(path, identity); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordLaunchIdentity(identity); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for identity without intent, got %v", err)
	}
}
