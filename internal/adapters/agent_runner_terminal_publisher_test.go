package adapters

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/tickets"
)

func terminalPublisherFixture(t *testing.T) (*AgentRunnerTerminalPublisher, AgentRunnerCompositeAdmission) {
	t.Helper()
	admission, _ := compositeAdmissionFixture(t)
	store := replayStoreMustNew(t, t.TempDir())
	if _, err := store.RecordRequest(admission.RequestRecord); err != nil {
		t.Fatal(err)
	}
	publisher := &AgentRunnerTerminalPublisher{
		Store:         store,
		CostLedger:    admission.CostLedger,
		RequestRecord: admission.RequestRecord,
		Clock: func() time.Time {
			return aiAvailabilityPolicy().EvaluatedAt
		},
	}
	return publisher, admission
}

func terminalObservedSettlement() AgentRunnerTerminalSettlement {
	amount := int64(500)
	calls := 1
	tokens := 500
	return AgentRunnerTerminalSettlement{
		Status:              tickets.CostSettlementObserved,
		ChargedAmountMicros: &amount,
		ObservedCalls:       &calls,
		ObservedTokens:      &tokens,
		ProviderEvidence:    []string{queueHashFour},
	}
}

func terminalUnknownSettlement() AgentRunnerTerminalSettlement {
	return AgentRunnerTerminalSettlement{Status: tickets.CostSettlementUnknown, ProviderEvidence: []string{queueHashFour}}
}

// queueHashFive is the fifth distinct digest used by the frozen-facts fixture.
const queueHashFive = "sha256:5555555555555555555555555555555555555555555555555555555555555555"

// terminalFrozenFacts is the five-digest fact set a completed terminal must
// carry; the chain refuses to route a completion without it.
func terminalFrozenFacts() *AgentRunnerFrozenFacts {
	return &AgentRunnerFrozenFacts{
		ManifestHash:            queueHashOne,
		DiffHash:                queueHashTwo,
		LogHash:                 queueHashThree,
		UsageHash:               queueHashFour,
		ProcessStopEvidenceHash: queueHashFive,
	}
}

// terminalFactsForCompletion attaches the facts only where the store requires
// them: a completed terminal intent must carry them, every other status must not.
func terminalFactsForCompletion(completion AgentRunnerCompletionRecord) *AgentRunnerFrozenFacts {
	if completion.Completion.Status != "completed" {
		return nil
	}
	return terminalFrozenFacts()
}

func terminalOutcomeFor(t *testing.T, request AgentRunnerRequestRecord, completionID string, settlement AgentRunnerTerminalSettlement) AgentRunnerTerminalOutcome {
	t.Helper()
	record := replayStoreCompletionRecord(t, request, completionID)
	outcome := AgentRunnerTerminalOutcome{Completion: record.Completion, Settlement: settlement}
	if record.Completion.Status == "completed" {
		outcome.Facts = terminalFrozenFacts()
	}
	return outcome
}

func terminalLedgerEntryCount(t *testing.T, ledger *tickets.CostLedger) int {
	t.Helper()
	snapshot, err := ledger.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	return len(snapshot.Ledger.Entries)
}

func TestAgentRunnerTerminalPublisherHappyPathPublishesChainInOrder(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	intentHookCalled := false
	settleHookCalled := false

	publisher.afterIntentWrite = func() error {
		intentHookCalled = true
		if _, found, err := publisher.Store.TerminalIntent(request.Request.RequestID); err != nil || !found {
			t.Fatalf("intent must be visible before settlement: found=%v err=%v", found, err)
		}
		if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
			t.Fatalf("settlement must not run before the intent is durable: %v", err)
		}
		state, err := publisher.Store.State(request.Request.RequestID)
		if err != nil || state != AgentRunnerReplayStateDispatchedUnknown {
			t.Fatalf("C must not be published before settlement: state=%v err=%v", state, err)
		}
		return nil
	}
	publisher.afterSettle = func() error {
		settleHookCalled = true
		if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); !errors.Is(err, tickets.ErrCostReservationSettled) {
			t.Fatalf("expected the reservation to be settled at this point, got %v", err)
		}
		state, err := publisher.Store.State(request.Request.RequestID)
		if err != nil || state != AgentRunnerReplayStateDispatchedUnknown {
			t.Fatalf("C must not be published before its settlement: state=%v err=%v", state, err)
		}
		return nil
	}
	result, err := publisher.PublishTerminal(context.Background(), outcome)
	if err != nil {
		t.Fatalf("expected terminal publication to succeed, got %v", err)
	}
	if !intentHookCalled || !settleHookCalled {
		t.Fatalf("order assertions require both hooks to fire: intent=%v settle=%v", intentHookCalled, settleHookCalled)
	}
	if result.Replayed {
		t.Fatalf("fresh publication must not report a replay")
	}
	if len(result.Evidence) != 3 {
		t.Fatalf("expected a three-part evidence chain, got %v", result.Evidence)
	}
	closure, found, err := publisher.Store.TerminalClosure(request.Request.RequestID)
	if err != nil || !found {
		t.Fatalf("expected a terminal closure: found=%v err=%v", found, err)
	}
	if result.Evidence[2] != closure.RecordHash {
		t.Fatalf("evidence closure hash = %s, want %s", result.Evidence[2], closure.RecordHash)
	}
	state, err := publisher.Store.State(request.Request.RequestID)
	if err != nil || state != AgentRunnerReplayStateTerminalReceiptPresent {
		t.Fatalf("expected terminal state, got state=%v err=%v", state, err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); !errors.Is(err, tickets.ErrCostReservationSettled) {
		t.Fatalf("expected reservation settled, got %v", err)
	}
}

func TestAgentRunnerTerminalPublisherReplayIsIdempotent(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	outcome := terminalOutcomeFor(t, publisher.RequestRecord, "completion-one", terminalObservedSettlement())
	first, err := publisher.PublishTerminal(context.Background(), outcome)
	if err != nil {
		t.Fatal(err)
	}
	entriesAfterFirst := terminalLedgerEntryCount(t, admission.CostLedger)
	second, err := publisher.PublishTerminal(context.Background(), outcome)
	if err != nil || !second.Replayed {
		t.Fatalf("expected idempotent replay, got replayed=%v err=%v", second.Replayed, err)
	}
	if terminalLedgerEntryCount(t, admission.CostLedger) != entriesAfterFirst {
		t.Fatalf("replay must not append settlement entries")
	}
	if second.Evidence[2] != first.Evidence[2] {
		t.Fatalf("replay closure hash = %s, want %s", second.Evidence[2], first.Evidence[2])
	}
}

func TestAgentRunnerTerminalPublisherConcurrentSingleWinner(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	outcome := terminalOutcomeFor(t, publisher.RequestRecord, "completion-one", terminalObservedSettlement())
	const workers = 8
	results := make([]error, workers)
	var wg sync.WaitGroup
	for index := 0; index < workers; index++ {
		wg.Add(1)
		go func(slot int) {
			defer wg.Done()
			_, results[slot] = publisher.PublishTerminal(context.Background(), outcome)
		}(index)
	}
	wg.Wait()
	for _, err := range results {
		if err != nil {
			t.Fatalf("concurrent terminal publication failed: %v", err)
		}
	}
	if count := terminalLedgerEntryCount(t, admission.CostLedger); count != 3 {
		t.Fatalf("ledger entries = %d, want 3 (allocation + reservation + one settlement)", count)
	}
	closure, found, err := publisher.Store.TerminalClosure(publisher.RequestRecord.Request.RequestID)
	if err != nil || !found {
		t.Fatalf("expected a single closure: found=%v err=%v", found, err)
	}
	if _, found, err := publisher.Store.TerminalIntent(publisher.RequestRecord.Request.RequestID); err != nil || !found {
		t.Fatalf("expected a single intent: found=%v err=%v", found, err)
	}
	_ = closure
}

func TestAgentRunnerTerminalPublisherCrashAfterIntentRecovers(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	sentinel := errors.New("crash after intent")
	publisher.afterIntentWrite = func() error { return sentinel }
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, sentinel) {
		t.Fatalf("expected injected crash, got %v", err)
	}
	if _, found, err := publisher.Store.TerminalIntent(request.Request.RequestID); err != nil || !found {
		t.Fatalf("intent must survive the crash: found=%v err=%v", found, err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
		t.Fatalf("no settlement may happen before the crash point: %v", err)
	}
	publisher.afterIntentWrite = nil
	result, err := publisher.PublishTerminal(context.Background(), outcome)
	if err != nil || !result.Replayed {
		t.Fatalf("expected recovery replay, got replayed=%v err=%v", result.Replayed, err)
	}
	if count := terminalLedgerEntryCount(t, admission.CostLedger); count != 3 {
		t.Fatalf("recovery must settle exactly once, entries=%d", count)
	}
	state, err := publisher.Store.State(request.Request.RequestID)
	if err != nil || state != AgentRunnerReplayStateTerminalReceiptPresent {
		t.Fatalf("expected recovered terminal state, got state=%v err=%v", state, err)
	}
}

func TestAgentRunnerTerminalPublisherCrashAfterSettleRecovers(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	sentinel := errors.New("crash after settle")
	publisher.afterSettle = func() error { return sentinel }
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, sentinel) {
		t.Fatalf("expected injected crash, got %v", err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); !errors.Is(err, tickets.ErrCostReservationSettled) {
		t.Fatalf("settlement must be booked before the crash point: %v", err)
	}
	state, err := publisher.Store.State(request.Request.RequestID)
	if err != nil || state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("C must be missing after the crash: state=%v err=%v", state, err)
	}
	intent, found, err := publisher.Store.TerminalIntent(request.Request.RequestID)
	if err != nil || !found {
		t.Fatalf("intent must survive the crash: found=%v err=%v", found, err)
	}
	// Simulate the later crash window where C is published but the closure is not.
	if _, err := publisher.Store.RecordCompletion(request, intent.Intent.Completion); err != nil {
		t.Fatal(err)
	}
	publisher.afterSettle = nil
	result, err := publisher.PublishTerminal(context.Background(), outcome)
	if err != nil || !result.Replayed {
		t.Fatalf("expected recovery replay, got replayed=%v err=%v", result.Replayed, err)
	}
	if count := terminalLedgerEntryCount(t, admission.CostLedger); count != 3 {
		t.Fatalf("recovery must not double-settle, entries=%d", count)
	}
	if _, found, err := publisher.Store.TerminalClosure(request.Request.RequestID); err != nil || !found {
		t.Fatalf("recovery must complete the closure: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerTerminalPublisherOrphanCompletionFailsClosed(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	completion, err := NewAgentRunnerCompletionRecord(outcome.Completion)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := publisher.Store.RecordCompletion(request, completion); err != nil {
		t.Fatal(err)
	}
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrAgentRunnerTerminalOrphan) {
		t.Fatalf("expected orphan completion rejection, got %v", err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
		t.Fatalf("orphan rejection must not settle: %v", err)
	}
	if _, found, err := publisher.Store.TerminalIntent(request.Request.RequestID); err != nil || found {
		t.Fatalf("orphan rejection must not adopt the completion: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerTerminalPublisherRejectsForeignSettlement(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	amount := int64(500)
	calls := 1
	tokens := 500
	if _, err := admission.CostLedger.Settle(tickets.SettlementInput{
		EntryID:             "settle-foreign",
		OccurredAt:          time.Date(2026, 9, 11, 8, 0, 0, 0, time.UTC),
		RunID:               request.Request.RunID,
		RequestID:           request.Request.RequestID,
		IdempotencyKey:      "settle-foreign-key",
		ReservationHash:     request.Request.BudgetHash,
		Status:              tickets.CostSettlementObserved,
		ChargedAmountMicros: &amount,
		ObservedCalls:       &calls,
		ObservedTokens:      &tokens,
		ProviderEvidence:    []string{queueHashFour},
	}); err != nil {
		t.Fatal(err)
	}
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	_, err := publisher.PublishTerminal(context.Background(), outcome)
	if !errors.Is(err, ErrAgentRunnerTerminalSettlementBlocked) || !errors.Is(err, tickets.ErrCostReservationSettled) {
		t.Fatalf("expected foreign-settlement rejection, got %v", err)
	}
	if _, found, err := publisher.Store.TerminalIntent(request.Request.RequestID); err != nil || found {
		t.Fatalf("foreign settlement must be rejected before any intent: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerTerminalPublisherUnknownSettlementKeepsReservation(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalUnknownSettlement())
	if _, err := publisher.PublishTerminal(context.Background(), outcome); err != nil {
		t.Fatalf("expected unknown settlement publication, got %v", err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); !errors.Is(err, tickets.ErrCostReservationSettled) {
		t.Fatalf("unknown settlement must consume the reservation slot: %v", err)
	}
	if holds := admission.CostLedger.UnknownHoldReservations(); holds != 1 {
		t.Fatalf("unknown holds = %d, want 1", holds)
	}
	if outstanding := admission.CostLedger.OutstandingReservations(); outstanding != 0 {
		t.Fatalf("outstanding reservations = %d, want 0", outstanding)
	}
	summary := admission.CostLedger.Summary()
	if summary.UnknownReservedAmountMicros == nil || *summary.UnknownReservedAmountMicros != 1000 {
		t.Fatalf("unknown reserved amount = %v, want 1000", summary.UnknownReservedAmountMicros)
	}
	if replay, err := publisher.PublishTerminal(context.Background(), outcome); err != nil || !replay.Replayed {
		t.Fatalf("expected unknown settlement replay, got replayed=%v err=%v", replay.Replayed, err)
	}
}

func TestAgentRunnerTerminalPublisherSettlementMatrix(t *testing.T) {
	uncertainCompletion := func(t *testing.T, request AgentRunnerRequestRecord) AgentRunnerCompletion {
		return replayStoreUncertainCompletionRecord(t, request, "completion-one").Completion
	}
	incompleteFailure := func(t *testing.T, request AgentRunnerRequestRecord) AgentRunnerCompletion {
		body := replayStoreCompletionRecord(t, request, "completion-one").Completion
		body.Status = "failed"
		body.ExitCode = nil
		body.OutputManifestHash = nil
		body.UsageComplete = false
		body.ErrorEvidence = []string{queueHashFour}
		return body
	}
	cases := []struct {
		name       string
		completion func(*testing.T, AgentRunnerRequestRecord) AgentRunnerCompletion
		settlement AgentRunnerTerminalSettlement
		wantErr    bool
	}{
		{"uncertain requires unknown", uncertainCompletion, terminalObservedSettlement(), true},
		{"uncertain accepts unknown", uncertainCompletion, terminalUnknownSettlement(), false},
		{"incomplete usage requires unknown", incompleteFailure, terminalObservedSettlement(), true},
		{"incomplete usage accepts unknown", incompleteFailure, terminalUnknownSettlement(), false},
		{"observed requires a charged amount", func(t *testing.T, request AgentRunnerRequestRecord) AgentRunnerCompletion {
			return replayStoreCompletionRecord(t, request, "completion-one").Completion
		}, AgentRunnerTerminalSettlement{Status: tickets.CostSettlementObserved, ProviderEvidence: []string{queueHashFour}}, true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			publisher, _ := terminalPublisherFixture(t)
			outcome := AgentRunnerTerminalOutcome{
				Completion: testCase.completion(t, publisher.RequestRecord),
				Settlement: testCase.settlement,
			}
			_, err := publisher.PublishTerminal(context.Background(), outcome)
			if testCase.wantErr {
				if !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
					t.Fatalf("expected settlement matrix rejection, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected settlement matrix acceptance, got %v", err)
			}
		})
	}
}

func TestAgentRunnerTerminalPublisherRejectsChargedOverReserved(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	amount := int64(5000)
	calls := 1
	tokens := 500
	outcome := terminalOutcomeFor(t, request, "completion-one", AgentRunnerTerminalSettlement{
		Status:              tickets.CostSettlementObserved,
		ChargedAmountMicros: &amount,
		ObservedCalls:       &calls,
		ObservedTokens:      &tokens,
		ProviderEvidence:    []string{queueHashFour},
	})
	_, err := publisher.PublishTerminal(context.Background(), outcome)
	if !errors.Is(err, ErrAgentRunnerTerminalSettlementBlocked) || !errors.Is(err, tickets.ErrCostSettlementConflict) {
		t.Fatalf("expected charged-over-reserved rejection, got %v", err)
	}
	if _, found, err := publisher.Store.TerminalIntent(request.Request.RequestID); err != nil || !found {
		t.Fatalf("the intent is the settlement decision record and must remain: found=%v err=%v", found, err)
	}
	if _, found, err := publisher.Store.TerminalClosure(request.Request.RequestID); err != nil || found {
		t.Fatalf("no closure may exist after a rejected settlement: found=%v err=%v", found, err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
		t.Fatalf("rejected settlement must keep the reservation outstanding: %v", err)
	}
}

func TestAgentRunnerTerminalPublisherRefusesDifferentOutcomeSameRequest(t *testing.T) {
	publisher, _ := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	if _, err := publisher.PublishTerminal(context.Background(), outcome); err != nil {
		t.Fatal(err)
	}
	other := terminalOutcomeFor(t, request, "completion-two", terminalObservedSettlement())
	if _, err := publisher.PublishTerminal(context.Background(), other); !errors.Is(err, ErrAgentRunnerTerminalIntentConflict) {
		t.Fatalf("expected conflict for a different terminal outcome, got %v", err)
	}
}

func TestAgentRunnerTerminalPublisherUnprovenRejectsBeforeSettlement(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	publisher.Store.publicationDurability = PublishDurabilityUnproven
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("expected durability gate, got %v", err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
		t.Fatalf("unproven rejection must happen before any settlement: %v", err)
	}
	if entries := terminalLedgerEntryCount(t, admission.CostLedger); entries != 2 {
		t.Fatalf("unproven rejection must not append ledger entries, got %d", entries)
	}
	if _, found, err := publisher.Store.TerminalIntent(request.Request.RequestID); err != nil || found {
		t.Fatalf("unproven rejection must not publish an intent: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerTerminalPublisherRejectsResumeSessionMismatch(t *testing.T) {
	admission, _ := compositeAdmissionFixture(t)
	store := replayStoreMustNew(t, t.TempDir())
	resumeBody := admission.RequestRecord.Request
	resumeBody.Mode = "resume"
	resumeBody.PriorSessionID = "session-two"
	resumeBody.PriorCompletionHash = queueHashOne
	resumeRecord, err := NewAgentRunnerRequestRecord(resumeBody)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordRequest(resumeRecord); err != nil {
		t.Fatal(err)
	}
	publisher := &AgentRunnerTerminalPublisher{
		Store:         store,
		CostLedger:    admission.CostLedger,
		RequestRecord: resumeRecord,
		Clock:         func() time.Time { return aiAvailabilityPolicy().EvaluatedAt },
	}
	outcome := terminalOutcomeFor(t, resumeRecord, "completion-one", terminalObservedSettlement())
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("expected resume session mismatch rejection, got %v", err)
	}
	if _, found, err := publisher.Store.TerminalIntent(resumeRecord.Request.RequestID); err != nil || found {
		t.Fatalf("resume mismatch must be rejected before any intent: found=%v err=%v", found, err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(resumeRecord.Request.BudgetHash); err != nil {
		t.Fatalf("resume mismatch must be rejected before any settlement: %v", err)
	}

	matchedBody := resumeBody
	matchedBody.PriorSessionID = "session-one"
	matchedRecord, err := NewAgentRunnerRequestRecord(matchedBody)
	if err != nil {
		t.Fatal(err)
	}
	matchedStore := replayStoreMustNew(t, t.TempDir())
	if _, err := matchedStore.RecordRequest(matchedRecord); err != nil {
		t.Fatal(err)
	}
	matchedPublisher := &AgentRunnerTerminalPublisher{
		Store:         matchedStore,
		CostLedger:    admission.CostLedger,
		RequestRecord: matchedRecord,
		Clock:         func() time.Time { return aiAvailabilityPolicy().EvaluatedAt },
	}
	if _, err := matchedPublisher.PublishTerminal(context.Background(), terminalOutcomeFor(t, matchedRecord, "completion-one", terminalObservedSettlement())); err != nil {
		t.Fatalf("resume with the matching session must publish, got %v", err)
	}
}

func TestAgentRunnerTerminalPublisherDifferentClocksConverge(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	second := *publisher
	second.Clock = func() time.Time { return aiAvailabilityPolicy().EvaluatedAt.Add(time.Second) }
	outcome := terminalOutcomeFor(t, publisher.RequestRecord, "completion-one", terminalObservedSettlement())
	start := make(chan struct{})
	results := make([]error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		_, results[0] = publisher.PublishTerminal(context.Background(), outcome)
	}()
	go func() {
		defer wg.Done()
		<-start
		_, results[1] = second.PublishTerminal(context.Background(), outcome)
	}()
	close(start)
	wg.Wait()
	for index, err := range results {
		if err != nil {
			t.Fatalf("publisher %d failed: %v", index, err)
		}
	}
	if count := terminalLedgerEntryCount(t, admission.CostLedger); count != 3 {
		t.Fatalf("ledger entries = %d, want exactly one settlement", count)
	}
	if _, found, err := publisher.Store.TerminalClosure(publisher.RequestRecord.Request.RequestID); err != nil || !found {
		t.Fatalf("expected a single closure: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerTerminalPublisherDifferentClocksSequentialConverge(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	outcome := terminalOutcomeFor(t, publisher.RequestRecord, "completion-one", terminalObservedSettlement())
	if _, err := publisher.PublishTerminal(context.Background(), outcome); err != nil {
		t.Fatal(err)
	}
	second := *publisher
	second.Clock = func() time.Time { return aiAvailabilityPolicy().EvaluatedAt.Add(time.Second) }
	result, err := second.PublishTerminal(context.Background(), outcome)
	if err != nil || !result.Replayed {
		t.Fatalf("a different clock must converge on the stored slot, got replayed=%v err=%v", result.Replayed, err)
	}
	if count := terminalLedgerEntryCount(t, admission.CostLedger); count != 3 {
		t.Fatalf("sequential convergence must not settle again, entries=%d", count)
	}
}

func TestAgentRunnerTerminalPublisherRejectsInvalidEvidenceBeforeIntent(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	outcome.Settlement.ProviderEvidence = []string{"not-a-digest"}
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("expected invalid evidence rejection, got %v", err)
	}
	if _, found, err := publisher.Store.TerminalIntent(request.Request.RequestID); err != nil || found {
		t.Fatalf("invalid evidence must be rejected before any intent: found=%v err=%v", found, err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
		t.Fatalf("invalid evidence must be rejected before any settlement: %v", err)
	}

	negativeCallsOutcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	negativeCalls := -1
	negativeCallsOutcome.Settlement.ObservedCalls = &negativeCalls
	if _, err := publisher.PublishTerminal(context.Background(), negativeCallsOutcome); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("expected negative observed calls rejection, got %v", err)
	}
}

func TestAgentRunnerTerminalPublisherRefusesACompletedOutcomeWithoutFacts(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	outcome.Facts = nil
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("expected a completed-without-facts rejection, got %v", err)
	}
	if _, found, err := publisher.Store.TerminalIntent(request.Request.RequestID); err != nil || found {
		t.Fatalf("missing facts must be rejected before any intent: found=%v err=%v", found, err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
		t.Fatalf("missing facts must be rejected before any settlement: %v", err)
	}

	notCompleted := AgentRunnerTerminalOutcome{
		Completion: replayStoreUncertainCompletionRecord(t, request, "completion-two").Completion,
		Settlement: terminalUnknownSettlement(),
		Facts:      terminalFrozenFacts(),
	}
	if _, err := publisher.PublishTerminal(context.Background(), notCompleted); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("a non-completed outcome must never carry frozen facts: %v", err)
	}
	if _, found, err := publisher.Store.TerminalIntent(request.Request.RequestID); err != nil || found {
		t.Fatalf("facts on a non-completed outcome must be rejected before any intent: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerTerminalPublisherDeduplicatesEvidence(t *testing.T) {
	publisher, _ := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	outcome.Settlement.ProviderEvidence = []string{queueHashFour, queueHashFour}
	if _, err := publisher.PublishTerminal(context.Background(), outcome); err != nil {
		t.Fatalf("expected duplicate evidence to be deduplicated, got %v", err)
	}
	intent, found, err := publisher.Store.TerminalIntent(request.Request.RequestID)
	if err != nil || !found {
		t.Fatalf("expected an intent: found=%v err=%v", found, err)
	}
	occurrences := 0
	for _, digest := range intent.Intent.ProviderEvidence {
		if digest == queueHashFour {
			occurrences++
		}
	}
	if occurrences != 1 {
		t.Fatalf("duplicate evidence must be stored once, got %d occurrences of %s", occurrences, queueHashFour)
	}
	if len(intent.Intent.ProviderEvidence) != 3 {
		t.Fatalf("expected completion + request + one caller digest, got %v", intent.Intent.ProviderEvidence)
	}
}

func TestAgentRunnerTerminalPublisherConvergeUnprovenRefusesWithoutSettlement(t *testing.T) {
	t.Run("intent without settlement", func(t *testing.T) {
		publisher, admission := terminalPublisherFixture(t)
		request := publisher.RequestRecord
		outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
		sentinel := errors.New("crash after intent")
		publisher.afterIntentWrite = func() error { return sentinel }
		if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, sentinel) {
			t.Fatalf("expected injected crash, got %v", err)
		}
		publisher.afterIntentWrite = nil
		publisher.Store.publicationDurability = PublishDurabilityUnproven
		if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
			t.Fatalf("expected durability gate on the recovery path, got %v", err)
		}
		if entries := terminalLedgerEntryCount(t, admission.CostLedger); entries != 2 {
			t.Fatalf("recovery on an unproven store must not settle, entries=%d", entries)
		}
		if _, found, err := publisher.Store.TerminalClosure(request.Request.RequestID); err != nil || found {
			t.Fatalf("no closure may exist: found=%v err=%v", found, err)
		}
	})
	t.Run("settled without closure", func(t *testing.T) {
		publisher, admission := terminalPublisherFixture(t)
		request := publisher.RequestRecord
		outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
		sentinel := errors.New("crash after settle")
		publisher.afterSettle = func() error { return sentinel }
		if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, sentinel) {
			t.Fatalf("expected injected crash, got %v", err)
		}
		publisher.afterSettle = nil
		intent, found, err := publisher.Store.TerminalIntent(request.Request.RequestID)
		if err != nil || !found {
			t.Fatalf("expected the intent to survive: found=%v err=%v", found, err)
		}
		if _, err := publisher.Store.RecordCompletion(request, intent.Intent.Completion); err != nil {
			t.Fatal(err)
		}
		entriesBefore := terminalLedgerEntryCount(t, admission.CostLedger)
		publisher.Store.publicationDurability = PublishDurabilityUnproven
		if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
			t.Fatalf("expected durability gate on the recovery path, got %v", err)
		}
		if entries := terminalLedgerEntryCount(t, admission.CostLedger); entries != entriesBefore {
			t.Fatalf("recovery on an unproven store must not settle again, entries=%d want %d", entries, entriesBefore)
		}
		if _, found, err := publisher.Store.TerminalClosure(request.Request.RequestID); err != nil || found {
			t.Fatalf("no closure may exist: found=%v err=%v", found, err)
		}
	})
}

func TestAgentRunnerTerminalPublisherClosureWithoutCompletionIsCorruption(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	intent := terminalIntentRecordFor(t, request, "completion-one")
	if _, err := publisher.Store.RecordTerminalIntent(intent); err != nil {
		t.Fatal(err)
	}
	if _, err := publisher.Store.RecordTerminalClosure(terminalClosureRecordFor(t, intent)); err != nil {
		t.Fatal(err)
	}
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for a closure without a completion, got %v", err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
		t.Fatalf("corruption rejection must not settle: %v", err)
	}
}

func TestAgentRunnerTerminalPublisherUnprovenCompleteChainIsWriteFree(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	outcome := terminalOutcomeFor(t, publisher.RequestRecord, "completion-one", terminalObservedSettlement())
	first, err := publisher.PublishTerminal(context.Background(), outcome)
	if err != nil {
		t.Fatal(err)
	}
	entriesBefore := terminalLedgerEntryCount(t, admission.CostLedger)
	publisher.Store.publicationDurability = PublishDurabilityUnproven
	// Both probes are wired to the single settlement choke point and the
	// post-settlement point, so any branch that touches the ledger is caught.
	publisher.beforeSettle = func() error {
		t.Fatalf("a write-free replay on an unproven store must not call the ledger")
		return nil
	}
	publisher.afterSettle = func() error {
		t.Fatalf("a write-free replay on an unproven store must not settle")
		return nil
	}
	result, err := publisher.PublishTerminal(context.Background(), outcome)
	if err != nil || !result.Replayed {
		t.Fatalf("expected write-free replay, got replayed=%v err=%v", result.Replayed, err)
	}
	if result.Evidence[2] != first.Evidence[2] {
		t.Fatalf("write-free replay must return the stored closure hash")
	}
	if entries := terminalLedgerEntryCount(t, admission.CostLedger); entries != entriesBefore {
		t.Fatalf("write-free replay must not append ledger entries, entries=%d want %d", entries, entriesBefore)
	}
}

func TestAgentRunnerTerminalPublisherClosureBindingMismatchIsCorruption(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	intent := terminalIntentRecordFor(t, request, "completion-one")
	if _, err := publisher.Store.RecordTerminalIntent(intent); err != nil {
		t.Fatal(err)
	}
	// Publish C so the state check would pass; the only remaining failing cause
	// is the settlement-key binding of the mismatched closure.
	if _, err := publisher.Store.RecordCompletion(request, intent.Intent.Completion); err != nil {
		t.Fatal(err)
	}
	mismatched, err := NewAgentRunnerTerminalClosureRecord(AgentRunnerTerminalClosure{
		RequestID:                "request-one",
		IntentRecordHash:         intent.RecordHash,
		CompletionHash:           intent.Intent.Completion.RecordHash,
		SettlementEntryID:        "settle-other",
		SettlementIdempotencyKey: intent.Intent.SettlementIdempotencyKey,
		ClosedAt:                 intent.Intent.DecidedAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	path, err := publisher.Store.terminalClosurePath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(path, mismatched); err != nil {
		t.Fatal(err)
	}
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for a closure whose settlement keys do not bind, got %v", err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
		t.Fatalf("binding mismatch must not settle: %v", err)
	}
}

func TestAgentRunnerTerminalPublisherRecoveryRequiresPersistedRequest(t *testing.T) {
	publisher, admission := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	intent := terminalIntentRecordFor(t, request, "completion-one")
	if _, err := publisher.Store.RecordTerminalIntent(intent); err != nil {
		t.Fatal(err)
	}
	requestPath, err := publisher.Store.requestPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(requestPath); err != nil {
		t.Fatal(err)
	}
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())
	if _, err := publisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrAgentRunnerRequestConflict) {
		t.Fatalf("expected recovery to require the persisted request record, got %v", err)
	}
	if _, err := admission.CostLedger.RequireOutstandingReservation(request.Request.BudgetHash); err != nil {
		t.Fatalf("missing-request rejection must not settle: %v", err)
	}
}

func TestAgentRunnerTerminalPublisherFailsClosedOnInvalidInputs(t *testing.T) {
	publisher, _ := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	outcome := terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())

	var nilPublisher *AgentRunnerTerminalPublisher
	if _, err := nilPublisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("nil publisher must fail closed, got %v", err)
	}
	cases := map[string]func(*AgentRunnerTerminalPublisher){
		"store":  func(target *AgentRunnerTerminalPublisher) { target.Store = nil },
		"ledger": func(target *AgentRunnerTerminalPublisher) { target.CostLedger = nil },
		"clock":  func(target *AgentRunnerTerminalPublisher) { target.Clock = nil },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			copyPublisher := *publisher
			mutate(&copyPublisher)
			if _, err := copyPublisher.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
				t.Fatalf("missing %s must fail closed, got %v", name, err)
			}
		})
	}
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := publisher.PublishTerminal(cancelled, outcome); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("cancelled context must fail closed, got %v", err)
	}
	zeroClock := *publisher
	zeroClock.Clock = func() time.Time { return time.Time{} }
	if _, err := zeroClock.PublishTerminal(context.Background(), outcome); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("zero clock must fail closed, got %v", err)
	}
	bindingMismatch := outcome
	bindingMismatch.Completion.RequestID = "request-two"
	if _, err := publisher.PublishTerminal(context.Background(), bindingMismatch); !errors.Is(err, ErrInvalidAgentRunnerTerminalIntent) {
		t.Fatalf("completion binding mismatch must fail closed, got %v", err)
	}
}
