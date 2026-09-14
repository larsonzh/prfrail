package adapters

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func replayStoreRequestRecord(t *testing.T, requestID string) AgentRunnerRequestRecord {
	t.Helper()
	record, err := NewAgentRunnerRequestRecord(agentRunnerRequest(requestID))
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func replayStoreCompletionRecord(t *testing.T, request AgentRunnerRequestRecord, completionID string) AgentRunnerCompletionRecord {
	t.Helper()
	body := agentRunnerCompletion(completionID)
	body.RequestID = request.Request.RequestID
	body.RequestHash = request.RecordHash
	body.RunID = request.Request.RunID
	body.TaskID = request.Request.TaskID
	body.StepID = request.Request.StepID
	body.Attempt = request.Request.Attempt
	body.AdapterID = request.Request.AdapterID
	record, err := NewAgentRunnerCompletionRecord(body)
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func replayStoreUncertainCompletionRecord(t *testing.T, request AgentRunnerRequestRecord, completionID string) AgentRunnerCompletionRecord {
	t.Helper()
	body := agentRunnerCompletion(completionID)
	body.RequestID = request.Request.RequestID
	body.RequestHash = request.RecordHash
	body.RunID = request.Request.RunID
	body.TaskID = request.Request.TaskID
	body.StepID = request.Request.StepID
	body.Attempt = request.Request.Attempt
	body.AdapterID = request.Request.AdapterID
	body.Status = "uncertain"
	body.ErrorEvidence = []string{queueHashFour}
	record, err := NewAgentRunnerCompletionRecord(body)
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func replayStoreMustNew(t *testing.T, root string) *AgentRunnerReplayStore {
	t.Helper()
	store, err := newAgentRunnerReplayStoreAt(root)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		store.publicationDurability = PublishDurabilityProven
	}
	return store
}

func TestAgentRunnerReplayStoreRequiresAbsoluteRoot(t *testing.T) {
	if _, err := newAgentRunnerReplayStoreAt(""); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("expected invalid replay root for empty path, got %v", err)
	}
	if _, err := newAgentRunnerReplayStoreAt("relative\\replay-store"); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("expected invalid replay root for relative path, got %v", err)
	}
	if _, err := newAgentRunnerReplayStoreAt(t.TempDir()); err != nil {
		t.Fatalf("expected absolute replay root to pass, got %v", err)
	}
}

func TestAgentRunnerReplayStoreUsesSafePrefixedFileNames(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	requestPath, err := store.requestPath("con")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(requestPath) != "request.con.jsonl" {
		t.Fatalf("request filename = %s, want request.con.jsonl", filepath.Base(requestPath))
	}
	completionPath, err := store.completionPath("prn")
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(completionPath) != "completion.prn.jsonl" {
		t.Fatalf("completion filename = %s, want completion.prn.jsonl", filepath.Base(completionPath))
	}
}

func TestAgentRunnerReplayStorePublicationDurabilityMatchesPlatform(t *testing.T) {
	store, err := newAgentRunnerReplayStoreAt(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	want := PublishDurabilityProven
	if runtime.GOOS == "windows" {
		want = PublishDurabilityUnproven
	}
	if got := store.PublicationDurability(); got != want {
		t.Fatalf("publication durability = %s, want %s", got, want)
	}
}

func TestAgentRunnerReplayStoreRecordRequestRejectsUnprovenDurabilityWithoutWriting(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	store.publicationDurability = PublishDurabilityUnproven
	request := replayStoreRequestRecord(t, "request-one")

	decision, err := store.RecordRequest(request)
	if !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("expected durability-unproven error, got decision=%s err=%v", decision, err)
	}
	if decision == AgentRunnerReplayDecisionFirstDispatch {
		t.Fatal("unproven durability must never return first dispatch")
	}
	path, pathErr := store.requestPath(request.Request.RequestID)
	if pathErr != nil {
		t.Fatal(pathErr)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("unproven request must not create a request file, stat error = %v", statErr)
	}
}

func TestAgentRunnerReplayStoreZeroValueCannotBypassDurabilityGate(t *testing.T) {
	var store AgentRunnerReplayStore
	if _, err := store.RecordRequest(replayStoreRequestRecord(t, "request-one")); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("zero-value store bypassed durability gate: %v", err)
	}
}

func TestAgentRunnerReplayStoreUnprovenCompletionAllowsOnlyExistingReplay(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	completion := replayStoreCompletionRecord(t, request, "completion-one")
	store.publicationDurability = PublishDurabilityUnproven
	if replayed, err := store.RecordCompletion(request, completion); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) || replayed {
		t.Fatalf("unproven completion publication = replayed:%v err:%v, want durability error", replayed, err)
	}

	completionPath, err := store.completionPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(completionPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("unproven completion must not create a completion file, stat error = %v", err)
	}
	store.publicationDurability = PublishDurabilityProven
	if _, err := store.RecordCompletion(request, completion); err != nil {
		t.Fatal(err)
	}
	store.publicationDurability = PublishDurabilityUnproven
	if replayed, err := store.RecordCompletion(request, completion); err != nil || !replayed {
		t.Fatalf("existing completion replay = replayed:%v err:%v, want replay", replayed, err)
	}
}

func TestAgentRunnerReplayStoreFirstRequestAndIdempotentReplay(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")

	decision, err := store.RecordRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if decision != AgentRunnerReplayDecisionFirstDispatch {
		t.Fatalf("decision = %s, want %s", decision, AgentRunnerReplayDecisionFirstDispatch)
	}

	decision, err = store.RecordRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if decision != AgentRunnerReplayDecisionUnknownBlock {
		t.Fatalf("decision = %s, want %s", decision, AgentRunnerReplayDecisionUnknownBlock)
	}

	state, err := store.State(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("state = %s, want %s", state, AgentRunnerReplayStateDispatchedUnknown)
	}
}

func TestAgentRunnerReplayStoreRejectsRequestConflict(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}

	changed := agentRunnerRequest("request-one")
	changed.ContextHash = queueHashOne
	conflicting, err := NewAgentRunnerRequestRecord(changed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordRequest(conflicting); !errors.Is(err, ErrAgentRunnerRequestConflict) {
		t.Fatalf("expected request conflict, got %v", err)
	}
}

func TestAgentRunnerReplayStoreRejectsInvalidCompletionBinding(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}

	completionBody := agentRunnerCompletion("completion-one")
	completionBody.RequestID = request.Request.RequestID
	completionBody.RequestHash = queueHashTwo
	completionBody.RunID = request.Request.RunID
	completionBody.TaskID = request.Request.TaskID
	completionBody.StepID = request.Request.StepID
	completionBody.Attempt = request.Request.Attempt
	completionBody.AdapterID = request.Request.AdapterID
	completion, err := NewAgentRunnerCompletionRecord(completionBody)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordCompletion(request, completion); !errors.Is(err, ErrInvalidAgentRunnerCompletion) {
		t.Fatalf("expected completion binding rejection, got %v", err)
	}
}

func TestAgentRunnerReplayStoreInvalidCompletionDoesNotConsumeRequestIdentity(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")

	completionBody := agentRunnerCompletion("completion-one")
	completionBody.RequestID = request.Request.RequestID
	completionBody.RequestHash = queueHashTwo
	completionBody.RunID = request.Request.RunID
	completionBody.TaskID = request.Request.TaskID
	completionBody.StepID = request.Request.StepID
	completionBody.Attempt = request.Request.Attempt
	completionBody.AdapterID = request.Request.AdapterID
	invalid, err := NewAgentRunnerCompletionRecord(completionBody)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordCompletion(request, invalid); !errors.Is(err, ErrInvalidAgentRunnerCompletion) {
		t.Fatalf("expected completion binding rejection, got %v", err)
	}

	state, err := store.State(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if state != AgentRunnerReplayStateAbsent {
		t.Fatalf("state after rejected completion = %s, want %s", state, AgentRunnerReplayStateAbsent)
	}

	decision, err := store.RecordRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if decision != AgentRunnerReplayDecisionFirstDispatch {
		t.Fatalf("decision = %s, want %s", decision, AgentRunnerReplayDecisionFirstDispatch)
	}
}

func TestAgentRunnerReplayStoreCompletionFirstReplayAndConflict(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}

	completion := replayStoreCompletionRecord(t, request, "completion-one")
	replayed, err := store.RecordCompletion(request, completion)
	if err != nil {
		t.Fatal(err)
	}
	if replayed {
		t.Fatal("first completion write must not be replay")
	}

	replayed, err = store.RecordCompletion(request, completion)
	if err != nil {
		t.Fatal(err)
	}
	if !replayed {
		t.Fatal("same completion hash must be replay")
	}

	conflicting := replayStoreCompletionRecord(t, request, "completion-two")
	if _, err := store.RecordCompletion(request, conflicting); !errors.Is(err, ErrAgentRunnerCompletionConflict) {
		t.Fatalf("expected completion conflict, got %v", err)
	}
}

func TestAgentRunnerReplayStoreCollisionNeverReplacesExistingCompletion(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")

	completion := replayStoreCompletionRecord(t, request, "completion-one")
	if _, err := store.RecordCompletion(request, completion); err != nil {
		t.Fatal(err)
	}
	path, err := store.completionPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	conflicting := replayStoreCompletionRecord(t, request, "completion-two")
	if _, err := store.RecordCompletion(request, conflicting); !errors.Is(err, ErrAgentRunnerCompletionConflict) {
		t.Fatalf("expected completion conflict, got %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("an occupied terminal slot must never be replaced by a conflicting completion")
	}
}

func TestAgentRunnerReplayStoreRequestCollisionVisibilityExhaustionIsConvergenceFailure(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	path, err := store.requestPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeCanonicalLineNoReplace(path, request); err != nil {
		t.Fatal(err)
	}
	originalLoad := replayStoreLoadRequestRecord
	replayStoreLoadRequestRecord = func(string, string) (AgentRunnerRequestRecord, bool, error) {
		return AgentRunnerRequestRecord{}, false, nil
	}
	t.Cleanup(func() {
		replayStoreLoadRequestRecord = originalLoad
	})

	if _, err := store.RecordRequest(request); !errors.Is(err, ErrAgentRunnerReplayStoreConvergence) {
		t.Fatalf("expected convergence failure, got %v", err)
	} else if errors.Is(err, ErrAgentRunnerRequestConflict) {
		t.Fatalf("request collision visibility exhaustion must not report request conflict: %v", err)
	}
}

func TestAgentRunnerReplayStoreCompletionCollisionVisibilityExhaustionIsConvergenceFailure(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	completion := replayStoreCompletionRecord(t, request, "completion-one")
	path, err := store.completionPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeCanonicalLineNoReplace(path, completion); err != nil {
		t.Fatal(err)
	}
	originalLoad := replayStoreLoadCompletionRecord
	replayStoreLoadCompletionRecord = func(string, string) (AgentRunnerCompletionRecord, bool, error) {
		return AgentRunnerCompletionRecord{}, false, nil
	}
	t.Cleanup(func() {
		replayStoreLoadCompletionRecord = originalLoad
	})

	if _, err := store.RecordCompletion(request, completion); !errors.Is(err, ErrAgentRunnerReplayStoreConvergence) {
		t.Fatalf("expected convergence failure, got %v", err)
	} else if errors.Is(err, ErrAgentRunnerCompletionConflict) {
		t.Fatalf("completion collision visibility exhaustion must not report completion conflict: %v", err)
	}
}

func TestAgentRunnerReplayStorePublishSyncFailureFailsClosedWithoutRollback(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	if store.PublicationDurability() != PublishDurabilityProven {
		t.Skip("native platform cannot prove replay publication durability")
	}
	request := replayStoreRequestRecord(t, "request-one")
	sentinel := errors.New("sync failure")
	originalSync := replayStoreSyncParentDirectory
	replayStoreSyncParentDirectory = func(string) error {
		return sentinel
	}
	t.Cleanup(func() {
		replayStoreSyncParentDirectory = originalSync
	})

	if _, err := store.RecordRequest(request); !errors.Is(err, sentinel) {
		t.Fatalf("expected sync failure, got %v", err)
	}
	replayStoreSyncParentDirectory = originalSync
	path, err := store.requestPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("published request must remain visible after sync failure, stat error = %v", err)
	}
	restarted := replayStoreMustNew(t, store.root)
	state, err := restarted.State(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("state = %s, want %s", state, AgentRunnerReplayStateDispatchedUnknown)
	}
}

func TestAgentRunnerReplayStoreOrphanCompletionRejectsWithoutWritingRequest(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")

	orphan := replayStoreCompletionRecord(t, request, "completion-one")
	completionPath, err := store.completionPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeCanonicalLineNoReplace(completionPath, orphan); err != nil {
		t.Fatal(err)
	}

	// An orphan terminal receipt is a corruption signal: the call must fail
	// closed and must not heal the store by persisting the missing request.
	if _, err := store.RecordCompletion(request, orphan); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected replay store corruption, got %v", err)
	}

	requestPath, err := store.requestPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(requestPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a rejected completion must not persist a request record, stat error = %v", err)
	}
}

func TestAgentRunnerReplayStoreOrphanCompletionRejectsRecordRequestWithoutWritingRequest(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")

	orphan := replayStoreCompletionRecord(t, request, "completion-one")
	completionPath, err := store.completionPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeCanonicalLineNoReplace(completionPath, orphan); err != nil {
		t.Fatal(err)
	}

	if _, err := store.RecordRequest(request); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected replay store corruption, got %v", err)
	}

	requestPath, err := store.requestPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(requestPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a rejected request must not persist a request record, stat error = %v", err)
	}
}

func TestAgentRunnerReplayStoreRestartRecoveryForUnknownAndCompleted(t *testing.T) {
	root := t.TempDir()
	store := replayStoreMustNew(t, root)

	requestUnknown := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(requestUnknown); err != nil {
		t.Fatal(err)
	}
	requestCompleted := replayStoreRequestRecord(t, "request-two")
	if _, err := store.RecordRequest(requestCompleted); err != nil {
		t.Fatal(err)
	}
	completion := replayStoreCompletionRecord(t, requestCompleted, "completion-two")
	if _, err := store.RecordCompletion(requestCompleted, completion); err != nil {
		t.Fatal(err)
	}

	restarted := replayStoreMustNew(t, root)
	unknownState, err := restarted.State("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if unknownState != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("unknown state = %s, want %s", unknownState, AgentRunnerReplayStateDispatchedUnknown)
	}
	completedState, err := restarted.State("request-two")
	if err != nil {
		t.Fatal(err)
	}
	if completedState != AgentRunnerReplayStateTerminalReceiptPresent {
		t.Fatalf("completed state = %s, want %s", completedState, AgentRunnerReplayStateTerminalReceiptPresent)
	}
}

func TestAgentRunnerReplayStoreUncertainCompletionYieldsTerminalReceiptPresent(t *testing.T) {
	root := t.TempDir()
	store := replayStoreMustNew(t, root)
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	completion := replayStoreUncertainCompletionRecord(t, request, "completion-one")
	if replayed, err := store.RecordCompletion(request, completion); err != nil {
		t.Fatal(err)
	} else if replayed {
		t.Fatal("first completion write must not be replay")
	}

	state, err := store.State(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if state != AgentRunnerReplayStateTerminalReceiptPresent {
		t.Fatalf("state = %s, want %s", state, AgentRunnerReplayStateTerminalReceiptPresent)
	}

	decision, err := store.RecordRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if decision != AgentRunnerReplayDecisionTerminalReceiptPresent {
		t.Fatalf("decision = %s, want %s", decision, AgentRunnerReplayDecisionTerminalReceiptPresent)
	}
}

func TestAgentRunnerReplayStoreRequestOnlyCrashWindowBlocksRedispatch(t *testing.T) {
	root := t.TempDir()
	store := replayStoreMustNew(t, root)
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}

	restarted := replayStoreMustNew(t, root)
	state, err := restarted.State("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("state = %s, want %s", state, AgentRunnerReplayStateDispatchedUnknown)
	}
	decision, err := restarted.RecordRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if decision != AgentRunnerReplayDecisionUnknownBlock {
		t.Fatalf("decision = %s, want %s", decision, AgentRunnerReplayDecisionUnknownBlock)
	}
}

func TestAgentRunnerReplayStoreCompletionPersistedWindowConvergesReplay(t *testing.T) {
	root := t.TempDir()
	store := replayStoreMustNew(t, root)
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	completion := replayStoreCompletionRecord(t, request, "completion-one")
	if _, err := store.RecordCompletion(request, completion); err != nil {
		t.Fatal(err)
	}

	restarted := replayStoreMustNew(t, root)
	decision, err := restarted.RecordRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if decision != AgentRunnerReplayDecisionTerminalReceiptPresent {
		t.Fatalf("decision = %s, want %s", decision, AgentRunnerReplayDecisionTerminalReceiptPresent)
	}
}

func TestAgentRunnerReplayStoreConcurrentSameRequestHasSingleFirstDispatch(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")

	const workers = 24
	start := make(chan struct{})
	decisions := make(chan AgentRunnerReplayDecision, workers)
	errorsByCall := make(chan error, workers)
	var wg sync.WaitGroup
	for worker := 0; worker < workers; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			decision, err := store.RecordRequest(request)
			decisions <- decision
			errorsByCall <- err
		}()
	}
	close(start)
	wg.Wait()
	close(decisions)
	close(errorsByCall)

	firstDispatches := 0
	unknownBlocks := 0
	for err := range errorsByCall {
		if err != nil {
			t.Fatal(err)
		}
	}
	for decision := range decisions {
		switch decision {
		case AgentRunnerReplayDecisionFirstDispatch:
			firstDispatches++
		case AgentRunnerReplayDecisionUnknownBlock:
			unknownBlocks++
		default:
			t.Fatalf("unexpected decision %s", decision)
		}
	}
	if firstDispatches != 1 {
		t.Fatalf("first dispatches = %d, want 1", firstDispatches)
	}
	if unknownBlocks != workers-1 {
		t.Fatalf("unknown blocks = %d, want %d", unknownBlocks, workers-1)
	}
}

func TestAgentRunnerReplayStoreConcurrentSameRequestAcrossStoresHasSingleFirstDispatch(t *testing.T) {
	root := t.TempDir()
	storeA := replayStoreMustNew(t, root)
	storeB := replayStoreMustNew(t, root)
	request := replayStoreRequestRecord(t, "request-one")

	const workers = 24
	start := make(chan struct{})
	decisions := make(chan AgentRunnerReplayDecision, workers)
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
			decision, err := active.RecordRequest(request)
			decisions <- decision
			errorsByCall <- err
		}(selected)
	}
	close(start)
	wg.Wait()
	close(decisions)
	close(errorsByCall)

	firstDispatches := 0
	unknownBlocks := 0
	for err := range errorsByCall {
		if err != nil {
			t.Fatal(err)
		}
	}
	for decision := range decisions {
		switch decision {
		case AgentRunnerReplayDecisionFirstDispatch:
			firstDispatches++
		case AgentRunnerReplayDecisionUnknownBlock:
			unknownBlocks++
		default:
			t.Fatalf("unexpected decision %s", decision)
		}
	}
	if firstDispatches != 1 {
		t.Fatalf("first dispatches = %d, want 1", firstDispatches)
	}
	if unknownBlocks != workers-1 {
		t.Fatalf("unknown blocks = %d, want %d", unknownBlocks, workers-1)
	}
}

func TestAgentRunnerReplayStoreConcurrentDifferentCompletionsAcrossStoresConflicts(t *testing.T) {
	root := t.TempDir()
	storeA := replayStoreMustNew(t, root)
	storeB := replayStoreMustNew(t, root)
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := storeA.RecordRequest(request); err != nil {
		t.Fatal(err)
	}

	completionA := replayStoreCompletionRecord(t, request, "completion-one")
	completionB := replayStoreCompletionRecord(t, request, "completion-two")

	type completionResult struct {
		replayed bool
		err      error
	}

	start := make(chan struct{})
	results := make(chan completionResult, 2)
	go func() {
		<-start
		replayed, err := storeA.RecordCompletion(request, completionA)
		results <- completionResult{replayed: replayed, err: err}
	}()
	go func() {
		<-start
		replayed, err := storeB.RecordCompletion(request, completionB)
		results <- completionResult{replayed: replayed, err: err}
	}()
	close(start)

	first := <-results
	second := <-results

	successes := 0
	conflicts := 0
	for _, result := range []completionResult{first, second} {
		if result.err == nil {
			if result.replayed {
				t.Fatal("newly published completion must not be replay")
			}
			successes++
			continue
		}
		if errors.Is(result.err, ErrAgentRunnerCompletionConflict) {
			conflicts++
			continue
		}
		t.Fatalf("unexpected completion result error: %v", result.err)
	}
	if successes != 1 {
		t.Fatalf("successes = %d, want 1", successes)
	}
	if conflicts != 1 {
		t.Fatalf("conflicts = %d, want 1", conflicts)
	}
}

func TestAgentRunnerReplayStoreStateRejectsOrphanCompletion(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	orphan := replayStoreCompletionRecord(t, request, "completion-one")
	path, err := store.completionPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeCanonicalLineNoReplace(path, orphan); err != nil {
		t.Fatal(err)
	}

	if _, err := store.State(request.Request.RequestID); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected orphan completion fail-closed error, got %v", err)
	}
}

func TestAgentRunnerReplayStoreStateRejectsCorruptRequestJSON(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	path, err := store.requestPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("{\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := store.State("request-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption sentinel for decode failure, got %v", err)
	} else if !strings.Contains(err.Error(), path) {
		t.Fatalf("expected error to include path %s, got %v", path, err)
	}
}

func TestAgentRunnerReplayStoreStateRejectsDiskIDMismatch(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	requestPath, err := store.requestPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	requestTwo := replayStoreRequestRecord(t, "request-two")
	if err := writeCanonicalLineNoReplace(requestPath, requestTwo); err != nil {
		t.Fatal(err)
	}

	if _, err := store.State("request-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption sentinel for requestId mismatch, got %v", err)
	} else if !strings.Contains(err.Error(), requestPath) {
		t.Fatalf("expected mismatch error to include path %s, got %v", requestPath, err)
	}
}

func TestAgentRunnerReplayStoreStateRejectsCompletionIDMismatch(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}

	body := agentRunnerCompletion("completion-one")
	body.RequestID = "request-two"
	body.RequestHash = queueHashFour
	body.RunID = request.Request.RunID
	body.TaskID = request.Request.TaskID
	body.StepID = request.Request.StepID
	body.Attempt = request.Request.Attempt
	body.AdapterID = request.Request.AdapterID
	mismatch, err := NewAgentRunnerCompletionRecord(body)
	if err != nil {
		t.Fatal(err)
	}
	completionPath, err := store.completionPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeCanonicalLineNoReplace(completionPath, mismatch); err != nil {
		t.Fatal(err)
	}

	if _, err := store.State("request-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption sentinel for completion requestId mismatch, got %v", err)
	} else if !strings.Contains(err.Error(), completionPath) {
		t.Fatalf("expected mismatch error to include path %s, got %v", completionPath, err)
	}
}

func TestAgentRunnerReplayStoreConvergesCrossStorePublicationAfterAbsentRequestRead(t *testing.T) {
	root := t.TempDir()
	request := replayStoreRequestRecord(t, "request-one")
	completion := replayStoreCompletionRecord(t, request, "completion-one")

	tests := []struct {
		name string
		run  func(*AgentRunnerReplayStore) error
	}{
		{
			name: "state",
			run: func(store *AgentRunnerReplayStore) error {
				state, err := store.State(request.Request.RequestID)
				if err == nil && state != AgentRunnerReplayStateTerminalReceiptPresent {
					t.Fatalf("state = %s, want %s", state, AgentRunnerReplayStateTerminalReceiptPresent)
				}
				return err
			},
		},
		{
			name: "record request",
			run: func(store *AgentRunnerReplayStore) error {
				decision, err := store.RecordRequest(request)
				if err == nil && decision != AgentRunnerReplayDecisionTerminalReceiptPresent {
					t.Fatalf("decision = %s, want %s", decision, AgentRunnerReplayDecisionTerminalReceiptPresent)
				}
				return err
			},
		},
		{
			name: "record completion",
			run: func(store *AgentRunnerReplayStore) error {
				replayed, err := store.RecordCompletion(request, completion)
				if err == nil && !replayed {
					t.Fatal("cross-store completion must converge as replay")
				}
				return err
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := filepath.Join(root, strings.ReplaceAll(test.name, " ", "-"))
			if err := os.Mkdir(root, 0o755); err != nil {
				t.Fatal(err)
			}
			store := replayStoreMustNew(t, root)
			publisher := replayStoreMustNew(t, root)
			store.afterRequestReadLocked = func() {
				if _, err := publisher.RecordCompletion(request, completion); err != nil {
					t.Fatal(err)
				}
			}
			if err := test.run(store); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAgentRunnerReplayStoreMonotonicStateInvariants(t *testing.T) {
	request := replayStoreRequestRecord(t, "request-one")
	completion := replayStoreCompletionRecord(t, request, "completion-one")
	tests := []struct {
		name      string
		write     func(*AgentRunnerReplayStore) error
		wantState AgentRunnerReplayState
		wantError error
	}{
		{name: "absent", wantState: AgentRunnerReplayStateAbsent},
		{
			name: "dispatched",
			write: func(store *AgentRunnerReplayStore) error {
				_, err := store.RecordRequest(request)
				return err
			},
			wantState: AgentRunnerReplayStateDispatchedUnknown,
		},
		{
			name: "terminal",
			write: func(store *AgentRunnerReplayStore) error {
				_, err := store.RecordCompletion(request, completion)
				return err
			},
			wantState: AgentRunnerReplayStateTerminalReceiptPresent,
		},
		{
			name: "orphan fails closed",
			write: func(store *AgentRunnerReplayStore) error {
				path, err := store.completionPath(request.Request.RequestID)
				if err != nil {
					return err
				}
				return writeCanonicalLineNoReplace(path, completion)
			},
			wantError: ErrAgentRunnerReplayStoreCorruption,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store := replayStoreMustNew(t, t.TempDir())
			if test.write != nil {
				if err := test.write(store); err != nil {
					t.Fatal(err)
				}
			}
			state, err := store.State(request.Request.RequestID)
			if !errors.Is(err, test.wantError) {
				t.Fatalf("state error = %v, want %v", err, test.wantError)
			}
			if err == nil && state != test.wantState {
				t.Fatalf("state = %s, want %s", state, test.wantState)
			}
		})
	}
}

func TestAgentRunnerReplayStoreStateClassifiesDiskBindingMismatchAsCorruption(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	body := agentRunnerCompletion("completion-one")
	body.RequestID = request.Request.RequestID
	body.RequestHash = queueHashTwo
	body.RunID = request.Request.RunID
	body.TaskID = request.Request.TaskID
	body.StepID = request.Request.StepID
	body.Attempt = request.Request.Attempt
	body.AdapterID = request.Request.AdapterID
	completion, err := NewAgentRunnerCompletionRecord(body)
	if err != nil {
		t.Fatal(err)
	}
	path, err := store.completionPath(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeCanonicalLineNoReplace(path, completion); err != nil {
		t.Fatal(err)
	}

	if _, err := store.State(request.Request.RequestID); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected disk binding corruption, got %v", err)
	}
}

func TestAgentRunnerReplayStoreAbsentState(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	state, err := store.State("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if state != AgentRunnerReplayStateAbsent {
		t.Fatalf("state = %s, want %s", state, AgentRunnerReplayStateAbsent)
	}
}
