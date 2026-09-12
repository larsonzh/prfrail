package adapters

import (
	"errors"
	"os"
	"path/filepath"
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
	store, err := NewAgentRunnerReplayStore(root)
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func TestAgentRunnerReplayStoreRequiresAbsoluteRoot(t *testing.T) {
	if _, err := NewAgentRunnerReplayStore(""); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("expected invalid replay root for empty path, got %v", err)
	}
	if _, err := NewAgentRunnerReplayStore("relative\\replay-store"); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("expected invalid replay root for relative path, got %v", err)
	}
	if _, err := NewAgentRunnerReplayStore(t.TempDir()); err != nil {
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
