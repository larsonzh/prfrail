package adapters

import (
	"errors"
	"sync"
	"testing"
)

func terminalIntentRecordFor(t *testing.T, request AgentRunnerRequestRecord, completionID string) AgentRunnerTerminalIntentRecord {
	t.Helper()
	completion := replayStoreCompletionRecord(t, request, completionID)
	amount := int64(500)
	key := stableTerminalSettlementKey(request.Request.RequestID, completion.RecordHash)
	intent, err := NewAgentRunnerTerminalIntentRecord(AgentRunnerTerminalIntent{
		RequestID:                request.Request.RequestID,
		RequestHash:              request.RecordHash,
		RunID:                    request.Request.RunID,
		TaskID:                   request.Request.TaskID,
		StepID:                   request.Request.StepID,
		Attempt:                  request.Request.Attempt,
		AdapterID:                request.Request.AdapterID,
		Completion:               completion,
		Facts:                    terminalFactsForCompletion(completion),
		SettlementEntryID:        key,
		SettlementIdempotencyKey: key,
		ReservationHash:          request.Request.BudgetHash,
		SettlementStatus:         "observed",
		ChargedAmountMicros:      &amount,
		ObservedCalls:            intPointer(1),
		ObservedTokens:           intPointer(500),
		ProviderEvidence:         []string{completion.RecordHash, request.RecordHash, queueHashFour},
		DecidedAt:                "2026-09-11T08:05:00Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	return intent
}

func terminalClosureRecordFor(t *testing.T, intent AgentRunnerTerminalIntentRecord) AgentRunnerTerminalClosureRecord {
	t.Helper()
	closure, err := NewAgentRunnerTerminalClosureRecord(AgentRunnerTerminalClosure{
		RequestID:                intent.Intent.RequestID,
		IntentRecordHash:         intent.RecordHash,
		CompletionHash:           intent.Intent.Completion.RecordHash,
		SettlementEntryID:        intent.Intent.SettlementEntryID,
		SettlementIdempotencyKey: intent.Intent.SettlementIdempotencyKey,
		ClosedAt:                 intent.Intent.DecidedAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	return closure
}

func TestAgentRunnerTerminalRecordValidation(t *testing.T) {
	request := replayStoreRequestRecord(t, "request-one")
	intent := terminalIntentRecordFor(t, request, "completion-one")
	if err := ValidateAgentRunnerTerminalIntentRecord(intent); err != nil {
		t.Fatalf("expected valid intent, got %v", err)
	}
	if decoded, err := DecodeAgentRunnerTerminalIntentRecord(launchWire(t, intent)); err != nil || decoded.RecordHash != intent.RecordHash {
		t.Fatalf("intent round trip failed: %v", err)
	}
	closure := terminalClosureRecordFor(t, intent)
	if err := ValidateAgentRunnerTerminalClosureRecord(closure); err != nil {
		t.Fatalf("expected valid closure, got %v", err)
	}
	if decoded, err := DecodeAgentRunnerTerminalClosureRecord(launchWire(t, closure)); err != nil || decoded.RecordHash != closure.RecordHash {
		t.Fatalf("closure round trip failed: %v", err)
	}

	missingDigest := intent.Intent
	missingDigest.ProviderEvidence = []string{queueHashFour}
	if _, err := NewAgentRunnerTerminalIntentRecord(missingDigest); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("expected evidence binding rejection, got %v", err)
	}
	completedWithoutFacts := intent.Intent
	completedWithoutFacts.Facts = nil
	if _, err := NewAgentRunnerTerminalIntentRecord(completedWithoutFacts); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("expected completed-without-facts rejection, got %v", err)
	}
	unknownWithAmount := intent.Intent
	unknownWithAmount.SettlementStatus = "unknown"
	if _, err := NewAgentRunnerTerminalIntentRecord(unknownWithAmount); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("expected unknown-with-amount rejection, got %v", err)
	}
	tampered := intent
	tampered.RecordHash = queueHashOne
	if err := ValidateAgentRunnerTerminalIntentRecord(tampered); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("expected tampered hash rejection, got %v", err)
	}
	badClosure := closure.Closure
	badClosure.ClosedAt = "never"
	if _, err := NewAgentRunnerTerminalClosureRecord(badClosure); !errors.Is(err, ErrInvalidAgentRunnerTerminalClosure) {
		t.Fatalf("expected closure time rejection, got %v", err)
	}
}

func TestAgentRunnerTerminalIntentFactsBindingFailsClosed(t *testing.T) {
	tests := map[string]struct {
		status string
		facts  *AgentRunnerFrozenFacts
		refuse bool
	}{
		"completed with facts":    {status: "completed", facts: terminalFrozenFacts()},
		"completed without facts": {status: "completed", facts: nil, refuse: true},
		"failed with facts":       {status: "failed", facts: terminalFrozenFacts(), refuse: true},
		"failed without facts":    {status: "failed", facts: nil},
		"uncertain with facts":    {status: "uncertain", facts: terminalFrozenFacts(), refuse: true},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			intent := AgentRunnerTerminalIntent{
				Completion: AgentRunnerCompletionRecord{Completion: AgentRunnerCompletion{Status: test.status}},
				Facts:      test.facts,
			}
			err := validateAgentRunnerTerminalIntentFacts(intent)
			if test.refuse && !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
				t.Fatalf("expected a fail-closed rejection, got %v", err)
			}
			if !test.refuse && err != nil {
				t.Fatalf("expected acceptance, got %v", err)
			}
		})
	}
	t.Run("facts must be mutually distinct", func(t *testing.T) {
		colliding := terminalFrozenFacts()
		colliding.LogHash = colliding.DiffHash
		intent := AgentRunnerTerminalIntent{
			Completion: AgentRunnerCompletionRecord{Completion: AgentRunnerCompletion{Status: "completed"}},
			Facts:      colliding,
		}
		if err := validateAgentRunnerTerminalIntentFacts(intent); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
			t.Fatalf("expected a collision rejection, got %v", err)
		}
	})
}

func TestAgentRunnerReplayStoreTerminalIntentRequiresPersistedRequest(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	intent := terminalIntentRecordFor(t, request, "completion-one")
	if _, err := store.RecordTerminalIntent(intent); !errors.Is(err, ErrAgentRunnerRequestConflict) {
		t.Fatalf("expected request conflict without persisted request, got %v", err)
	}
	if _, found, err := store.TerminalIntent("request-one"); err != nil || found {
		t.Fatalf("no intent may be published without a request: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerReplayStoreTerminalIntentSingleWinnerReplayAndConflicts(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	intent := terminalIntentRecordFor(t, request, "completion-one")
	if decision, err := store.RecordTerminalIntent(intent); err != nil || decision != AgentRunnerTerminalIntentDecisionFirst {
		t.Fatalf("expected first intent, got decision=%v err=%v", decision, err)
	}
	if decision, err := store.RecordTerminalIntent(intent); err != nil || decision != AgentRunnerTerminalIntentDecisionPresent {
		t.Fatalf("expected intent replay, got decision=%v err=%v", decision, err)
	}

	otherCompletion := terminalIntentRecordFor(t, request, "completion-two")
	if _, err := store.RecordTerminalIntent(otherCompletion); !errors.Is(err, ErrAgentRunnerTerminalIntentConflict) {
		t.Fatalf("expected conflict for a different completion, got %v", err)
	}
	changedPlan, err := NewAgentRunnerTerminalIntentRecord(func() AgentRunnerTerminalIntent {
		body := intent.Intent
		amount := int64(700)
		body.ChargedAmountMicros = &amount
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordTerminalIntent(changedPlan); !errors.Is(err, ErrAgentRunnerTerminalIntentConflict) {
		t.Fatalf("expected conflict for a different settlement plan, got %v", err)
	}
}

func TestAgentRunnerReplayStoreTerminalClosureLifecycle(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	intent := terminalIntentRecordFor(t, request, "completion-one")
	closure := terminalClosureRecordFor(t, intent)
	if _, err := store.RecordTerminalClosure(closure); !errors.Is(err, ErrAgentRunnerTerminalIntentConflict) {
		t.Fatalf("expected conflict without an intent, got %v", err)
	}
	if _, err := store.RecordTerminalIntent(intent); err != nil {
		t.Fatal(err)
	}
	if replayed, err := store.RecordTerminalClosure(closure); err != nil || replayed {
		t.Fatalf("expected first closure, got replayed=%v err=%v", replayed, err)
	}
	if replayed, err := store.RecordTerminalClosure(closure); err != nil || !replayed {
		t.Fatalf("expected closure replay, got replayed=%v err=%v", replayed, err)
	}
	changed, err := NewAgentRunnerTerminalClosureRecord(func() AgentRunnerTerminalClosure {
		body := closure.Closure
		body.ClosedAt = "2026-09-11T08:06:00Z"
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordTerminalClosure(changed); !errors.Is(err, ErrAgentRunnerTerminalClosureConflict) {
		t.Fatalf("expected closure conflict, got %v", err)
	}
}

func TestAgentRunnerReplayStoreTerminalClosureOrphanIsCorruption(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	intent := terminalIntentRecordFor(t, request, "completion-one")
	closure := terminalClosureRecordFor(t, intent)
	path, err := store.terminalClosurePath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(path, closure); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordTerminalClosure(closure); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for closure without intent, got %v", err)
	}
}

func TestAgentRunnerReplayStoreTerminalRunIDBinding(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	store.runID = "run-one"
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	foreignBody := terminalIntentRecordFor(t, request, "completion-one").Intent
	foreignBody.RunID = "run-two"
	foreignCompletion, err := NewAgentRunnerCompletionRecord(func() AgentRunnerCompletion {
		body := foreignBody.Completion.Completion
		body.RunID = "run-two"
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	foreignBody.Completion = foreignCompletion
	foreignBody.ProviderEvidence = []string{foreignCompletion.RecordHash, request.RecordHash, queueHashFour}
	foreign, err := NewAgentRunnerTerminalIntentRecord(foreignBody)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordTerminalIntent(foreign); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("expected run binding conflict on write, got %v", err)
	}
	path, err := store.terminalIntentPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(path, foreign); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.TerminalIntent("request-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for foreign run binding on disk, got %v", err)
	}
}

func TestAgentRunnerReplayStoreTerminalUnprovenReplayOnly(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	intent := terminalIntentRecordFor(t, request, "completion-one")
	store.publicationDurability = PublishDurabilityUnproven
	if _, err := store.RecordTerminalIntent(intent); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("expected durability gate for new intent, got %v", err)
	}
	store.publicationDurability = PublishDurabilityProven
	if _, err := store.RecordTerminalIntent(intent); err != nil {
		t.Fatal(err)
	}
	closure := terminalClosureRecordFor(t, intent)
	store.publicationDurability = PublishDurabilityUnproven
	if _, err := store.RecordTerminalClosure(closure); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("expected durability gate for new closure, got %v", err)
	}
	store.publicationDurability = PublishDurabilityProven
	if _, err := store.RecordTerminalClosure(closure); err != nil {
		t.Fatal(err)
	}
	store.publicationDurability = PublishDurabilityUnproven
	if replayed, err := store.RecordTerminalClosure(closure); err != nil || !replayed {
		t.Fatalf("expected unproven store to allow exact closure replay, got replayed=%v err=%v", replayed, err)
	}
	if decision, err := store.RecordTerminalIntent(intent); err != nil || decision != AgentRunnerTerminalIntentDecisionPresent {
		t.Fatalf("expected unproven store to allow exact intent replay, got decision=%v err=%v", decision, err)
	}
}

func TestAgentRunnerReplayStoreTerminalCrossStoreSingleWinner(t *testing.T) {
	root := t.TempDir()
	storeA := replayStoreMustNew(t, root)
	storeB := replayStoreMustNew(t, root)
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := storeA.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	intent := terminalIntentRecordFor(t, request, "completion-one")

	const workers = 24
	start := make(chan struct{})
	decisions := make(chan AgentRunnerTerminalIntentDecision, workers)
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
			decision, err := active.RecordTerminalIntent(intent)
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
	firsts := 0
	presents := 0
	for decision := range decisions {
		switch decision {
		case AgentRunnerTerminalIntentDecisionFirst:
			firsts++
		case AgentRunnerTerminalIntentDecisionPresent:
			presents++
		default:
			t.Fatalf("unexpected decision %s", decision)
		}
	}
	if firsts != 1 || presents != workers-1 {
		t.Fatalf("intent decisions = first:%d present:%d, want first:1 present:%d", firsts, presents, workers-1)
	}

	closure := terminalClosureRecordFor(t, intent)
	replays := make(chan bool, workers)
	closureErrors := make(chan error, workers)
	start = make(chan struct{})
	for worker := 0; worker < workers; worker++ {
		selected := storeA
		if worker%2 == 1 {
			selected = storeB
		}
		wg.Add(1)
		go func(active *AgentRunnerReplayStore) {
			defer wg.Done()
			<-start
			replayed, err := active.RecordTerminalClosure(closure)
			replays <- replayed
			closureErrors <- err
		}(selected)
	}
	close(start)
	wg.Wait()
	close(replays)
	close(closureErrors)
	for err := range closureErrors {
		if err != nil {
			t.Fatal(err)
		}
	}
	firstClosures := 0
	replayedClosures := 0
	for replayed := range replays {
		if replayed {
			replayedClosures++
		} else {
			firstClosures++
		}
	}
	if firstClosures != 1 || replayedClosures != workers-1 {
		t.Fatalf("closure outcomes = first:%d replay:%d, want first:1 replay:%d", firstClosures, replayedClosures, workers-1)
	}
}

func TestAgentRunnerReplayStoreTerminalMethodsFailClosedOnZeroValueStore(t *testing.T) {
	var store AgentRunnerReplayStore
	request := replayStoreRequestRecord(t, "request-one")
	intent := terminalIntentRecordFor(t, request, "completion-one")
	if _, err := store.RecordTerminalIntent(intent); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("zero-value store must reject intent writes, got %v", err)
	}
	if _, _, err := store.TerminalIntent("request-one"); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("zero-value store must reject intent reads, got %v", err)
	}
	closure := terminalClosureRecordFor(t, intent)
	if _, err := store.RecordTerminalClosure(closure); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("zero-value store must reject closure writes, got %v", err)
	}
	if _, _, err := store.TerminalClosure("request-one"); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("zero-value store must reject closure reads, got %v", err)
	}
}

func TestAgentRunnerReplayStoreTerminalClosureIntentRunIDMismatchIsCorruption(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	store.runID = "run-one"
	request := replayStoreRequestRecord(t, "request-one")
	foreignBody := terminalIntentRecordFor(t, request, "completion-one").Intent
	foreignBody.RunID = "run-two"
	foreignCompletion, err := NewAgentRunnerCompletionRecord(func() AgentRunnerCompletion {
		body := foreignBody.Completion.Completion
		body.RunID = "run-two"
		return body
	}())
	if err != nil {
		t.Fatal(err)
	}
	foreignBody.Completion = foreignCompletion
	foreignBody.ProviderEvidence = []string{foreignCompletion.RecordHash, request.RecordHash, queueHashFour}
	foreign, err := NewAgentRunnerTerminalIntentRecord(foreignBody)
	if err != nil {
		t.Fatal(err)
	}
	path, err := store.terminalIntentPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(path, foreign); err != nil {
		t.Fatal(err)
	}
	closure := terminalClosureRecordFor(t, foreign)
	if _, err := store.RecordTerminalClosure(closure); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for a closure over a foreign-run intent, got %v", err)
	}
}
