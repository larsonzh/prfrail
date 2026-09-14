package adapters

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
)

// AgentRunnerTerminalIntentDecision describes the outcome of publishing a
// terminal-intent record for one request.
type AgentRunnerTerminalIntentDecision string

const (
	// AgentRunnerTerminalIntentDecisionFirst means this call published the first
	// terminal-intent; the caller owns the terminal slot.
	AgentRunnerTerminalIntentDecisionFirst AgentRunnerTerminalIntentDecision = "first-intent"
	// AgentRunnerTerminalIntentDecisionPresent means a terminal-intent already
	// exists; it is authoritative for recovery and its payload must match.
	AgentRunnerTerminalIntentDecisionPresent AgentRunnerTerminalIntentDecision = "intent-present"
)

// RecordTerminalIntent publishes the terminal-intent for a request that already
// has a persisted request record bound by the exact requestHash. The terminal
// slot is claimed with a no-replace write so concurrent publishers observe a
// single winner.
func (store *AgentRunnerReplayStore) RecordTerminalIntent(intent AgentRunnerTerminalIntentRecord) (AgentRunnerTerminalIntentDecision, error) {
	if err := store.validateConstructed(); err != nil {
		return "", err
	}
	if err := ValidateAgentRunnerTerminalIntentRecord(intent); err != nil {
		return "", err
	}
	body := intent.Intent
	if err := store.validateRecordRunBinding(body.RunID, body.RequestID); err != nil {
		return "", err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return "", err
	}
	requestPath, err := store.requestPath(body.RequestID)
	if err != nil {
		return "", err
	}
	path, err := store.terminalIntentPath(body.RequestID)
	if err != nil {
		return "", err
	}
	persisted, foundRequest, err := replayStoreLoadRequestRecord(requestPath, body.RequestID)
	if err != nil {
		return "", err
	}
	existing, foundIntent, err := replayStoreLoadTerminalIntent(path, body.RequestID)
	if err != nil {
		return "", err
	}
	// An existing terminal-intent is resolved without any write, so a replay or a
	// conflicting claim can never consume the request identity.
	if foundIntent {
		if !foundRequest {
			return "", replayStoreCorruption(path, "terminal intent present without a persisted request record", nil)
		}
		if err := store.validatePersistedRunID(requestPath, "request", persisted.Request.RunID); err != nil {
			return "", err
		}
		if err := store.validatePersistedRunID(path, "terminal intent", existing.Intent.RunID); err != nil {
			return "", err
		}
		if persisted.RecordHash != body.RequestHash {
			return "", fmt.Errorf("%w: terminal intent requestHash does not match persisted request for requestId %s", ErrAgentRunnerRequestConflict, body.RequestID)
		}
		if err := store.validateTerminalIntentRequestBinding(persisted, body); err != nil {
			return "", err
		}
		if existing.Intent.Completion.RecordHash != body.Completion.RecordHash {
			return "", fmt.Errorf("%w: terminal slot for requestId %s already decided a different completion %s", ErrAgentRunnerTerminalIntentConflict, body.RequestID, existing.Intent.Completion.RecordHash)
		}
		if !sameTerminalSettlementPlan(existing.Intent, body) {
			return "", fmt.Errorf("%w: terminal slot for requestId %s already decided a different settlement plan", ErrAgentRunnerTerminalIntentConflict, body.RequestID)
		}
		return AgentRunnerTerminalIntentDecisionPresent, nil
	}
	if store.PublicationDurability() != PublishDurabilityProven {
		return "", fmt.Errorf("%w: terminal intent publication", ErrAgentRunnerReplayStoreDurabilityUnproven)
	}
	if !foundRequest {
		return "", fmt.Errorf("%w: terminal intent requires a persisted request record for requestId %s", ErrAgentRunnerRequestConflict, body.RequestID)
	}
	if err := store.validatePersistedRunID(requestPath, "request", persisted.Request.RunID); err != nil {
		return "", err
	}
	if persisted.RecordHash != body.RequestHash {
		return "", fmt.Errorf("%w: terminal intent requestHash does not match persisted request for requestId %s", ErrAgentRunnerRequestConflict, body.RequestID)
	}
	if err := store.validateTerminalIntentRequestBinding(persisted, body); err != nil {
		return "", err
	}
	collided := false
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		if !collided {
			writeErr := writeReplayRecordNoReplace(path, intent)
			if writeErr == nil {
				return AgentRunnerTerminalIntentDecisionFirst, nil
			}
			if !errors.Is(writeErr, fs.ErrExist) {
				return "", writeErr
			}
			collided = true
		}
		visible, visibleOK, loadErr := replayStoreLoadTerminalIntent(path, body.RequestID)
		if loadErr != nil {
			return "", loadErr
		}
		if !visibleOK {
			runtime.Gosched()
			continue
		}
		if err := store.validatePersistedRunID(path, "terminal intent", visible.Intent.RunID); err != nil {
			return "", err
		}
		if visible.Intent.Completion.RecordHash != body.Completion.RecordHash ||
			!sameTerminalSettlementPlan(visible.Intent, body) {
			return "", fmt.Errorf("%w: terminal slot for requestId %s already decided a different outcome", ErrAgentRunnerTerminalIntentConflict, body.RequestID)
		}
		return AgentRunnerTerminalIntentDecisionPresent, nil
	}
	return "", fmt.Errorf("%w: terminal intent winner not visible after no-replace collision for requestId %s", ErrAgentRunnerReplayStoreConvergence, body.RequestID)
}

// TerminalIntent returns the stored terminal-intent for requestID.
func (store *AgentRunnerReplayStore) TerminalIntent(requestID string) (AgentRunnerTerminalIntentRecord, bool, error) {
	if err := store.validateConstructed(); err != nil {
		return AgentRunnerTerminalIntentRecord{}, false, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return AgentRunnerTerminalIntentRecord{}, false, err
	}
	path, err := store.terminalIntentPath(requestID)
	if err != nil {
		return AgentRunnerTerminalIntentRecord{}, false, err
	}
	record, found, err := replayStoreLoadTerminalIntent(path, requestID)
	if err != nil || !found {
		return record, found, err
	}
	if err := store.validatePersistedRunID(path, "terminal intent", record.Intent.RunID); err != nil {
		return AgentRunnerTerminalIntentRecord{}, false, err
	}
	return record, true, nil
}

// RecordTerminalClosure marks a terminal chain complete. Publishing requires a
// matching terminal-intent; a returned true means the exact closure was already
// present (replay), never a new terminal decision.
func (store *AgentRunnerReplayStore) RecordTerminalClosure(closure AgentRunnerTerminalClosureRecord) (bool, error) {
	if err := store.validateConstructed(); err != nil {
		return false, err
	}
	if err := ValidateAgentRunnerTerminalClosureRecord(closure); err != nil {
		return false, err
	}
	body := closure.Closure
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return false, err
	}
	intentPath, err := store.terminalIntentPath(body.RequestID)
	if err != nil {
		return false, err
	}
	path, err := store.terminalClosurePath(body.RequestID)
	if err != nil {
		return false, err
	}
	intent, foundIntent, err := replayStoreLoadTerminalIntent(intentPath, body.RequestID)
	if err != nil {
		return false, err
	}
	if foundIntent {
		if err := store.validatePersistedRunID(intentPath, "terminal intent", intent.Intent.RunID); err != nil {
			return false, err
		}
	}
	existing, foundClosure, err := replayStoreLoadTerminalClosure(path, body.RequestID)
	if err != nil {
		return false, err
	}
	if foundClosure {
		if !foundIntent {
			return false, replayStoreCorruption(path, "terminal closure present without a terminal intent", nil)
		}
		if err := store.validateTerminalClosureBinding(intent, body); err != nil {
			return false, err
		}
		if existing.RecordHash != closure.RecordHash {
			return false, fmt.Errorf("%w: terminal closure payload differs for requestId %s", ErrAgentRunnerTerminalClosureConflict, body.RequestID)
		}
		return true, nil
	}
	if store.PublicationDurability() != PublishDurabilityProven {
		return false, fmt.Errorf("%w: terminal closure publication", ErrAgentRunnerReplayStoreDurabilityUnproven)
	}
	if !foundIntent {
		return false, fmt.Errorf("%w: terminal closure requires a terminal intent for requestId %s", ErrAgentRunnerTerminalIntentConflict, body.RequestID)
	}
	if err := store.validateTerminalClosureBinding(intent, body); err != nil {
		return false, err
	}
	collided := false
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		if !collided {
			writeErr := writeReplayRecordNoReplace(path, closure)
			if writeErr == nil {
				return false, nil
			}
			if !errors.Is(writeErr, fs.ErrExist) {
				return false, writeErr
			}
			collided = true
		}
		visible, visibleOK, loadErr := replayStoreLoadTerminalClosure(path, body.RequestID)
		if loadErr != nil {
			return false, loadErr
		}
		if !visibleOK {
			runtime.Gosched()
			continue
		}
		if err := store.validateTerminalClosureBinding(intent, visible.Closure); err != nil {
			return false, err
		}
		if visible.RecordHash != closure.RecordHash {
			return false, fmt.Errorf("%w: terminal closure payload differs for requestId %s", ErrAgentRunnerTerminalClosureConflict, body.RequestID)
		}
		return true, nil
	}
	return false, fmt.Errorf("%w: terminal closure winner not visible after no-replace collision for requestId %s", ErrAgentRunnerReplayStoreConvergence, body.RequestID)
}

// TerminalClosure returns the stored terminal-closure for requestID.
func (store *AgentRunnerReplayStore) TerminalClosure(requestID string) (AgentRunnerTerminalClosureRecord, bool, error) {
	if err := store.validateConstructed(); err != nil {
		return AgentRunnerTerminalClosureRecord{}, false, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return AgentRunnerTerminalClosureRecord{}, false, err
	}
	path, err := store.terminalClosurePath(requestID)
	if err != nil {
		return AgentRunnerTerminalClosureRecord{}, false, err
	}
	return replayStoreLoadTerminalClosure(path, requestID)
}

func (store *AgentRunnerReplayStore) terminalIntentPath(requestID string) (string, error) {
	return store.pathForRequestID(requestID, "terminals", "terminal-intent.")
}

func (store *AgentRunnerReplayStore) terminalClosurePath(requestID string) (string, error) {
	return store.pathForRequestID(requestID, "terminals", "terminal-closure.")
}

// validateTerminalIntentRequestBinding cross-checks the intent identity fields
// against the persisted request record (the requestHash equality already binds
// them transitively; the explicit comparison keeps failures obvious).
func (store *AgentRunnerReplayStore) validateTerminalIntentRequestBinding(persisted AgentRunnerRequestRecord, intent AgentRunnerTerminalIntent) error {
	if persisted.Request.RunID != intent.RunID ||
		persisted.Request.TaskID != intent.TaskID ||
		persisted.Request.StepID != intent.StepID ||
		persisted.Request.Attempt != intent.Attempt ||
		persisted.Request.AdapterID != intent.AdapterID ||
		persisted.Request.BudgetHash != intent.ReservationHash {
		return fmt.Errorf("%w: terminal intent binding mismatch for requestId %s", ErrAgentRunnerTerminalIntentConflict, intent.RequestID)
	}
	return nil
}

func (store *AgentRunnerReplayStore) validateTerminalClosureBinding(intent AgentRunnerTerminalIntentRecord, closure AgentRunnerTerminalClosure) error {
	if intent.Intent.RequestID != closure.RequestID ||
		intent.RecordHash != closure.IntentRecordHash ||
		intent.Intent.Completion.RecordHash != closure.CompletionHash ||
		intent.Intent.SettlementEntryID != closure.SettlementEntryID ||
		intent.Intent.SettlementIdempotencyKey != closure.SettlementIdempotencyKey {
		return fmt.Errorf("%w: terminal closure binding mismatch for requestId %s", ErrAgentRunnerTerminalClosureConflict, closure.RequestID)
	}
	return nil
}

// sameTerminalSettlementPlan compares the settlement plan fields of two intent
// bodies. Field equality is sufficient because the plan is a deterministic
// function of the completion digest and the usage observation.
func sameTerminalSettlementPlan(left, right AgentRunnerTerminalIntent) bool {
	if left.SettlementEntryID != right.SettlementEntryID ||
		left.SettlementIdempotencyKey != right.SettlementIdempotencyKey ||
		left.ReservationHash != right.ReservationHash ||
		left.SettlementStatus != right.SettlementStatus {
		return false
	}
	if !sameInt64Pointer(left.ChargedAmountMicros, right.ChargedAmountMicros) {
		return false
	}
	if !sameIntPointer(left.ObservedCalls, right.ObservedCalls) || !sameIntPointer(left.ObservedTokens, right.ObservedTokens) {
		return false
	}
	if len(left.ProviderEvidence) != len(right.ProviderEvidence) {
		return false
	}
	for index := range left.ProviderEvidence {
		if left.ProviderEvidence[index] != right.ProviderEvidence[index] {
			return false
		}
	}
	return true
}

func sameInt64Pointer(left, right *int64) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func sameIntPointer(left, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func loadAgentRunnerTerminalIntent(path, requestID string) (AgentRunnerTerminalIntentRecord, bool, error) {
	wire, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return AgentRunnerTerminalIntentRecord{}, false, nil
	}
	if err != nil {
		return AgentRunnerTerminalIntentRecord{}, false, err
	}
	record, err := DecodeAgentRunnerTerminalIntentRecord(wire)
	if err != nil {
		return AgentRunnerTerminalIntentRecord{}, false, replayStoreCorruption(path, "terminal intent decode failed", err)
	}
	if record.Intent.RequestID != requestID {
		return AgentRunnerTerminalIntentRecord{}, false, replayStoreCorruption(path, fmt.Sprintf("terminal intent requestId mismatch on disk (want %s, got %s)", requestID, record.Intent.RequestID), nil)
	}
	return record, true, nil
}

func loadAgentRunnerTerminalClosure(path, requestID string) (AgentRunnerTerminalClosureRecord, bool, error) {
	wire, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return AgentRunnerTerminalClosureRecord{}, false, nil
	}
	if err != nil {
		return AgentRunnerTerminalClosureRecord{}, false, err
	}
	record, err := DecodeAgentRunnerTerminalClosureRecord(wire)
	if err != nil {
		return AgentRunnerTerminalClosureRecord{}, false, replayStoreCorruption(path, "terminal closure decode failed", err)
	}
	if record.Closure.RequestID != requestID {
		return AgentRunnerTerminalClosureRecord{}, false, replayStoreCorruption(path, fmt.Sprintf("terminal closure requestId mismatch on disk (want %s, got %s)", requestID, record.Closure.RequestID), nil)
	}
	return record, true, nil
}
