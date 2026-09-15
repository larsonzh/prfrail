package adapters

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
)

// terminalIntentRecord builds a terminal intent for one request from a
// completion record, mirroring what the publisher would have written.
func terminalIntentRecord(t *testing.T, request AgentRunnerRequestRecord, completion AgentRunnerCompletionRecord) AgentRunnerTerminalIntentRecord {
	t.Helper()
	amount := int64(500)
	calls, tokens := 1, 500
	intent, err := NewAgentRunnerTerminalIntentRecord(AgentRunnerTerminalIntent{
		RequestID:                request.Request.RequestID,
		RequestHash:              request.RecordHash,
		RunID:                    request.Request.RunID,
		TaskID:                   request.Request.TaskID,
		StepID:                   request.Request.StepID,
		Attempt:                  request.Request.Attempt,
		AdapterID:                request.Request.AdapterID,
		Completion:               completion,
		SettlementEntryID:        "settlement-one",
		SettlementIdempotencyKey: "settlement-key-one",
		ReservationHash:          request.Request.BudgetHash,
		SettlementStatus:         "observed",
		ChargedAmountMicros:      &amount,
		ObservedCalls:            &calls,
		ObservedTokens:           &tokens,
		ProviderEvidence:         []string{completion.RecordHash, request.RecordHash},
		DecidedAt:                "2026-08-01T07:08:09.000Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	return intent
}

func terminalClosureRecord(t *testing.T, intent AgentRunnerTerminalIntentRecord) AgentRunnerTerminalClosureRecord {
	t.Helper()
	closure, err := NewAgentRunnerTerminalClosureRecord(AgentRunnerTerminalClosure{
		RequestID:                intent.Intent.RequestID,
		IntentRecordHash:         intent.RecordHash,
		CompletionHash:           intent.Intent.Completion.RecordHash,
		SettlementEntryID:        intent.Intent.SettlementEntryID,
		SettlementIdempotencyKey: intent.Intent.SettlementIdempotencyKey,
		ClosedAt:                 "2026-08-01T07:09:00.000Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	return closure
}

func TestToChainAgentRunnerTerminalMapsThePublishedChain(t *testing.T) {
	publisher, _ := terminalPublisherFixture(t)
	request := publisher.RequestRecord
	if _, err := publisher.PublishTerminal(context.Background(), terminalOutcomeFor(t, request, "completion-one", terminalObservedSettlement())); err != nil {
		t.Fatal(err)
	}
	intent, found, err := publisher.Store.TerminalIntent(request.Request.RequestID)
	if err != nil || !found {
		t.Fatalf("intent missing: found=%v err=%v", found, err)
	}
	closure, found, err := publisher.Store.TerminalClosure(request.Request.RequestID)
	if err != nil || !found {
		t.Fatalf("closure missing: found=%v err=%v", found, err)
	}
	terminal, err := ToChainAgentRunnerTerminal(request, intent, closure)
	if err != nil {
		t.Fatal(err)
	}
	if terminal.RequestID != request.Request.RequestID || terminal.RequestHash != request.RecordHash ||
		terminal.CompletionHash != intent.Intent.Completion.RecordHash || terminal.RunID != request.Request.RunID ||
		terminal.TaskID != request.Request.TaskID || terminal.StepID != request.Request.StepID ||
		terminal.Attempt != request.Request.Attempt || terminal.Status != chain.AgentRunnerTerminalCompleted {
		t.Fatalf("terminal identity mismatch: %+v", terminal)
	}
	if terminal.SessionID != intent.Intent.Completion.Completion.SessionID {
		t.Fatalf("session was not carried over: %+v", terminal)
	}
	if terminal.PriorSessionID != "" || terminal.PriorCompletionHash != "" {
		t.Fatalf("create mode must not carry prior binding: %+v", terminal)
	}
	if len(terminal.Evidence) != len(dedupeTerminalEvidence(terminal.Evidence, nil)) {
		t.Fatalf("terminal evidence must be deduplicated: %+v", terminal.Evidence)
	}
	for _, want := range []string{request.RecordHash, intent.Intent.Completion.RecordHash, intent.RecordHash, closure.RecordHash, queueHashFour} {
		found := false
		for _, hash := range terminal.Evidence {
			if hash == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("terminal evidence missing %s: %+v", want, terminal.Evidence)
		}
	}
}

func TestLoadChainAgentRunnerTerminalRequiresASettledClosure(t *testing.T) {
	admission, _ := compositeAdmissionFixture(t)
	store := replayStoreMustNew(t, t.TempDir())
	if _, err := store.RecordRequest(admission.RequestRecord); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadChainAgentRunnerTerminal(store, admission.RequestRecord); !errors.Is(err, ErrInvalidAgentRunnerTerminalTranslation) {
		t.Fatalf("a request without a terminal must be refused, got %v", err)
	}
	publisher := &AgentRunnerTerminalPublisher{
		Store:         store,
		CostLedger:    admission.CostLedger,
		RequestRecord: admission.RequestRecord,
		Clock:         func() time.Time { return aiAvailabilityPolicy().EvaluatedAt },
	}
	if _, err := publisher.PublishTerminal(context.Background(), terminalOutcomeFor(t, admission.RequestRecord, "completion-one", terminalObservedSettlement())); err != nil {
		t.Fatal(err)
	}
	requestID := admission.RequestRecord.Request.RequestID
	intentBefore, found, err := store.TerminalIntent(requestID)
	if err != nil || !found {
		t.Fatalf("intent missing: found=%v err=%v", found, err)
	}
	closureBefore, found, err := store.TerminalClosure(requestID)
	if err != nil || !found {
		t.Fatalf("closure missing: found=%v err=%v", found, err)
	}
	if closureBefore.Closure.CompletionHash != intentBefore.Intent.Completion.RecordHash {
		t.Fatal("stored closure does not bind the stored completion")
	}
	snapshotBefore := storeFileSizes(t, store)
	terminal, err := LoadChainAgentRunnerTerminal(store, admission.RequestRecord)
	if err != nil {
		t.Fatal(err)
	}
	if terminal.Status != chain.AgentRunnerTerminalCompleted || terminal.CompletionHash != intentBefore.Intent.Completion.RecordHash {
		t.Fatalf("loaded terminal mismatch: %+v", terminal)
	}
	snapshotAfter := storeFileSizes(t, store)
	if len(snapshotAfter) != len(snapshotBefore) {
		t.Fatalf("translation changed the store file set: before=%d after=%d", len(snapshotBefore), len(snapshotAfter))
	}
	for path, size := range snapshotBefore {
		if snapshotAfter[path] != size {
			t.Fatalf("translation rewrote %s", path)
		}
	}
	if _, err := LoadChainAgentRunnerTerminal(nil, admission.RequestRecord); !errors.Is(err, ErrInvalidAgentRunnerTerminalTranslation) {
		t.Fatalf("a missing store must be refused, got %v", err)
	}
}

// storeFileSizes snapshots every durable replay-store file with its size, so a
// read-only path can be proven write-free across all record kinds.
func storeFileSizes(t *testing.T, store *AgentRunnerReplayStore) map[string]int64 {
	t.Helper()
	sizes := map[string]int64{}
	err := filepath.WalkDir(store.root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		sizes[path] = info.Size()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return sizes
}

func TestToChainAgentRunnerTerminalFailsClosed(t *testing.T) {
	admission, _ := compositeAdmissionFixture(t)
	request := admission.RequestRecord
	completion := mustCompletionRecord(t, request, "completion-one")
	intent := terminalIntentRecord(t, request, completion)
	closure := terminalClosureRecord(t, intent)

	foreign := closure
	foreignBody, err := NewAgentRunnerTerminalClosureRecord(AgentRunnerTerminalClosure{
		RequestID:                closure.Closure.RequestID,
		IntentRecordHash:         queueHashOne,
		CompletionHash:           closure.Closure.CompletionHash,
		SettlementEntryID:        closure.Closure.SettlementEntryID,
		SettlementIdempotencyKey: closure.Closure.SettlementIdempotencyKey,
		ClosedAt:                 closure.Closure.ClosedAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	foreign = foreignBody

	foreignSettlement, err := NewAgentRunnerTerminalClosureRecord(AgentRunnerTerminalClosure{
		RequestID:                closure.Closure.RequestID,
		IntentRecordHash:         closure.Closure.IntentRecordHash,
		CompletionHash:           closure.Closure.CompletionHash,
		SettlementEntryID:        "settlement-two",
		SettlementIdempotencyKey: closure.Closure.SettlementIdempotencyKey,
		ClosedAt:                 closure.Closure.ClosedAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	foreignKey, err := NewAgentRunnerTerminalClosureRecord(AgentRunnerTerminalClosure{
		RequestID:                closure.Closure.RequestID,
		IntentRecordHash:         closure.Closure.IntentRecordHash,
		CompletionHash:           closure.Closure.CompletionHash,
		SettlementEntryID:        closure.Closure.SettlementEntryID,
		SettlementIdempotencyKey: "settlement-key-two",
		ClosedAt:                 closure.Closure.ClosedAt,
	})
	if err != nil {
		t.Fatal(err)
	}

	resumeBody := request.Request
	resumeBody.Mode = "resume"
	resumeBody.PriorSessionID = "session-two"
	resumeBody.PriorCompletionHash = queueHashOne
	resumeRecord, err := NewAgentRunnerRequestRecord(resumeBody)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		request AgentRunnerRequestRecord
		intent  AgentRunnerTerminalIntentRecord
		closure AgentRunnerTerminalClosureRecord
	}{
		"zero request":           {request: AgentRunnerRequestRecord{}, intent: intent, closure: closure},
		"zero intent":            {request: request, intent: AgentRunnerTerminalIntentRecord{}, closure: closure},
		"zero closure":           {request: request, intent: intent, closure: AgentRunnerTerminalClosureRecord{}},
		"foreign intent binding": {request: request, intent: intent, closure: foreign},
		"foreign settlement entry": {
			request: request,
			intent:  intent,
			closure: foreignSettlement,
		},
		"foreign settlement key": {
			request: request,
			intent:  intent,
			closure: foreignKey,
		},
		"resume session mismatch": {
			request: resumeRecord,
			intent:  terminalIntentRecord(t, resumeRecord, mustCompletionRecord(t, resumeRecord, "completion-one")),
			closure: closure,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := ToChainAgentRunnerTerminal(test.request, test.intent, test.closure); !errors.Is(err, ErrInvalidAgentRunnerTerminalTranslation) {
				t.Fatalf("expected fail-closed translation, got %v", err)
			}
		})
	}
}

func TestToChainAgentRunnerTerminalCarriesResumeBinding(t *testing.T) {
	admission, _ := compositeAdmissionFixture(t)
	body := admission.RequestRecord.Request
	body.Mode = "resume"
	body.PriorSessionID = "session-one"
	body.PriorCompletionHash = queueHashOne
	request, err := NewAgentRunnerRequestRecord(body)
	if err != nil {
		t.Fatal(err)
	}
	completion := mustCompletionRecord(t, request, "completion-one")
	intent := terminalIntentRecord(t, request, completion)
	terminal, err := ToChainAgentRunnerTerminal(request, intent, terminalClosureRecord(t, intent))
	if err != nil {
		t.Fatal(err)
	}
	if terminal.PriorSessionID != "session-one" || terminal.PriorCompletionHash != queueHashOne {
		t.Fatalf("resume binding was not translated: %+v", terminal)
	}
	if err := terminal.Validate(); err != nil {
		t.Fatalf("resume terminal must validate: %v", err)
	}
}

func mustCompletionRecord(t *testing.T, request AgentRunnerRequestRecord, completionID string) AgentRunnerCompletionRecord {
	t.Helper()
	return replayStoreCompletionRecord(t, request, completionID)
}
