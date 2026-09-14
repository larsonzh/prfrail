package adapters

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"
)

// AgentRunnerLaunchDecision describes the outcome of publishing a launch-intent
// receipt for one request.
type AgentRunnerLaunchDecision string

const (
	// AgentRunnerLaunchDecisionFirstLaunch means this call published the first
	// launch-intent receipt; the caller owns the launch slot.
	AgentRunnerLaunchDecisionFirstLaunch AgentRunnerLaunchDecision = "first-launch"
	// AgentRunnerLaunchDecisionAlreadyLaunched means a launch-intent receipt for
	// the same launch identity already exists. It proves intent only: whether the
	// process was actually spawned requires the identity record.
	AgentRunnerLaunchDecisionAlreadyLaunched AgentRunnerLaunchDecision = "already-launched"
)

// RecordLaunchReceipt publishes the pre-spawn launch-intent receipt for a
// request that already has a persisted request record. The request record must
// be visible with the exact bound requestHash; the launch slot is claimed with
// a no-replace write so concurrent dispatchers observe a single winner.
func (store *AgentRunnerReplayStore) RecordLaunchReceipt(receipt AgentRunnerLaunchReceiptRecord) (AgentRunnerLaunchDecision, error) {
	if err := store.validateConstructed(); err != nil {
		return "", err
	}
	if err := ValidateAgentRunnerLaunchReceiptRecord(receipt); err != nil {
		return "", err
	}
	body := receipt.Receipt
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
	path, err := store.launchReceiptPath(body.RequestID)
	if err != nil {
		return "", err
	}
	persisted, foundRequest, err := replayStoreLoadRequestRecord(requestPath, body.RequestID)
	if err != nil {
		return "", err
	}
	existing, foundReceipt, err := replayStoreLoadLaunchReceipt(path, body.RequestID)
	if err != nil {
		return "", err
	}
	// An existing launch intent is resolved without any write, so a replay or a
	// conflicting claim can never consume the request identity.
	if foundReceipt {
		if !foundRequest {
			return "", replayStoreCorruption(path, "launch receipt present without a persisted request record", nil)
		}
		if err := store.validatePersistedRunID(requestPath, "request", persisted.Request.RunID); err != nil {
			return "", err
		}
		if err := store.validatePersistedRunID(path, "launch receipt", existing.Receipt.RunID); err != nil {
			return "", err
		}
		if persisted.RecordHash != body.RequestHash {
			return "", fmt.Errorf("%w: launch receipt requestHash does not match persisted request for requestId %s", ErrAgentRunnerRequestConflict, body.RequestID)
		}
		if err := store.validateLaunchReceiptRequestBinding(persisted, body); err != nil {
			return "", err
		}
		if existing.Receipt.LaunchID != body.LaunchID || existing.Receipt.RequestHash != body.RequestHash {
			return "", fmt.Errorf("%w: launch slot for requestId %s is owned by launchId %s", ErrAgentRunnerLaunchReceiptConflict, body.RequestID, existing.Receipt.LaunchID)
		}
		if existing.RecordHash != receipt.RecordHash {
			return "", fmt.Errorf("%w: launch receipt payload differs for launchId %s", ErrAgentRunnerLaunchReceiptConflict, body.LaunchID)
		}
		return AgentRunnerLaunchDecisionAlreadyLaunched, nil
	}
	if store.PublicationDurability() != PublishDurabilityProven {
		return "", fmt.Errorf("%w: launch receipt publication", ErrAgentRunnerReplayStoreDurabilityUnproven)
	}
	// No launch intent exists yet: the request binding must hold before the
	// first durable write of the launch slot.
	if !foundRequest {
		return "", fmt.Errorf("%w: launch receipt requires a persisted request record for requestId %s", ErrAgentRunnerRequestConflict, body.RequestID)
	}
	if err := store.validatePersistedRunID(requestPath, "request", persisted.Request.RunID); err != nil {
		return "", err
	}
	if persisted.RecordHash != body.RequestHash {
		return "", fmt.Errorf("%w: launch receipt requestHash does not match persisted request for requestId %s", ErrAgentRunnerRequestConflict, body.RequestID)
	}
	if err := store.validateLaunchReceiptRequestBinding(persisted, body); err != nil {
		return "", err
	}
	collided := false
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		if !collided {
			writeErr := writeReplayRecordNoReplace(path, receipt)
			if writeErr == nil {
				return AgentRunnerLaunchDecisionFirstLaunch, nil
			}
			if !errors.Is(writeErr, fs.ErrExist) {
				return "", writeErr
			}
			collided = true
		}
		visible, visibleOK, loadErr := replayStoreLoadLaunchReceipt(path, body.RequestID)
		if loadErr != nil {
			return "", loadErr
		}
		if !visibleOK {
			runtime.Gosched()
			continue
		}
		if err := store.validatePersistedRunID(path, "launch receipt", visible.Receipt.RunID); err != nil {
			return "", err
		}
		if visible.Receipt.LaunchID != body.LaunchID || visible.RecordHash != receipt.RecordHash {
			return "", fmt.Errorf("%w: launch slot for requestId %s is owned by launchId %s", ErrAgentRunnerLaunchReceiptConflict, body.RequestID, visible.Receipt.LaunchID)
		}
		return AgentRunnerLaunchDecisionAlreadyLaunched, nil
	}
	return "", fmt.Errorf("%w: launch receipt winner not visible after no-replace collision for requestId %s", ErrAgentRunnerReplayStoreConvergence, body.RequestID)
}

// LaunchReceipt returns the stored launch-intent receipt for requestID.
func (store *AgentRunnerReplayStore) LaunchReceipt(requestID string) (AgentRunnerLaunchReceiptRecord, bool, error) {
	if err := store.validateConstructed(); err != nil {
		return AgentRunnerLaunchReceiptRecord{}, false, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return AgentRunnerLaunchReceiptRecord{}, false, err
	}
	path, err := store.launchReceiptPath(requestID)
	if err != nil {
		return AgentRunnerLaunchReceiptRecord{}, false, err
	}
	record, found, err := replayStoreLoadLaunchReceipt(path, requestID)
	if err != nil || !found {
		return record, found, err
	}
	if err := store.validatePersistedRunID(path, "launch receipt", record.Receipt.RunID); err != nil {
		return AgentRunnerLaunchReceiptRecord{}, false, err
	}
	return record, true, nil
}

// RecordLaunchIdentity publishes the post-spawn process identity. Publication
// requires a matching launch-intent receipt; a returned true means the exact
// identity was already present (replay), never a new spawn.
func (store *AgentRunnerReplayStore) RecordLaunchIdentity(identity AgentRunnerLaunchIdentityRecord) (bool, error) {
	if err := store.validateConstructed(); err != nil {
		return false, err
	}
	if err := ValidateAgentRunnerLaunchIdentityRecord(identity); err != nil {
		return false, err
	}
	body := identity.Identity
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return false, err
	}
	intentPath, err := store.launchReceiptPath(body.RequestID)
	if err != nil {
		return false, err
	}
	path, err := store.launchIdentityPath(body.LaunchID)
	if err != nil {
		return false, err
	}
	receipt, foundReceipt, err := replayStoreLoadLaunchReceipt(intentPath, body.RequestID)
	if err != nil {
		return false, err
	}
	existing, foundIdentity, err := replayStoreLoadLaunchIdentity(path, body.LaunchID)
	if err != nil {
		return false, err
	}
	// An existing process identity is resolved without any write, so a replay or
	// a conflicting claim can never consume the launch slot.
	if foundIdentity {
		if !foundReceipt {
			return false, replayStoreCorruption(path, "process identity present without a launch-intent receipt", nil)
		}
		if receipt.Receipt.LaunchID != body.LaunchID || receipt.Receipt.RequestID != body.RequestID || receipt.RecordHash != body.IntentRecordHash {
			return false, fmt.Errorf("%w: process identity intent binding mismatch for launchId %s", ErrAgentRunnerLaunchReceiptConflict, body.LaunchID)
		}
		if existing.RecordHash != identity.RecordHash {
			return false, fmt.Errorf("%w: process identity payload differs for launchId %s", ErrAgentRunnerLaunchIdentityConflict, body.LaunchID)
		}
		return true, nil
	}
	if store.PublicationDurability() != PublishDurabilityProven {
		return false, fmt.Errorf("%w: process identity publication", ErrAgentRunnerReplayStoreDurabilityUnproven)
	}
	// No identity exists yet: the launch intent must exist and bind before the
	// first durable write of the identity record.
	if !foundReceipt {
		return false, fmt.Errorf("%w: process identity requires a launch-intent receipt for launchId %s", ErrAgentRunnerLaunchReceiptConflict, body.LaunchID)
	}
	if receipt.Receipt.LaunchID != body.LaunchID || receipt.Receipt.RequestID != body.RequestID || receipt.RecordHash != body.IntentRecordHash {
		return false, fmt.Errorf("%w: process identity intent binding mismatch for launchId %s", ErrAgentRunnerLaunchReceiptConflict, body.LaunchID)
	}
	collided := false
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		if !collided {
			writeErr := writeReplayRecordNoReplace(path, identity)
			if writeErr == nil {
				return false, nil
			}
			if !errors.Is(writeErr, fs.ErrExist) {
				return false, writeErr
			}
			collided = true
		}
		visible, visibleOK, loadErr := replayStoreLoadLaunchIdentity(path, body.LaunchID)
		if loadErr != nil {
			return false, loadErr
		}
		if !visibleOK {
			runtime.Gosched()
			continue
		}
		if visible.RecordHash != identity.RecordHash {
			return false, fmt.Errorf("%w: process identity payload differs for launchId %s", ErrAgentRunnerLaunchIdentityConflict, body.LaunchID)
		}
		return true, nil
	}
	return false, fmt.Errorf("%w: process identity winner not visible after no-replace collision for launchId %s", ErrAgentRunnerReplayStoreConvergence, body.LaunchID)
}

// LaunchIdentity returns the stored process-identity record for launchID.
func (store *AgentRunnerReplayStore) LaunchIdentity(launchID string) (AgentRunnerLaunchIdentityRecord, bool, error) {
	if err := store.validateConstructed(); err != nil {
		return AgentRunnerLaunchIdentityRecord{}, false, err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return AgentRunnerLaunchIdentityRecord{}, false, err
	}
	path, err := store.launchIdentityPath(launchID)
	if err != nil {
		return AgentRunnerLaunchIdentityRecord{}, false, err
	}
	return replayStoreLoadLaunchIdentity(path, launchID)
}

func (store *AgentRunnerReplayStore) launchReceiptPath(requestID string) (string, error) {
	return store.pathForRequestID(requestID, "launches", "intent.")
}

// validateLaunchReceiptRequestBinding cross-checks the launch receipt identity
// fields against the persisted request record. The requestHash equality check
// already binds every field transitively; the explicit comparison keeps the
// failure mode obvious if hash domains ever diverge.
func (store *AgentRunnerReplayStore) validateLaunchReceiptRequestBinding(persisted AgentRunnerRequestRecord, receipt AgentRunnerLaunchReceipt) error {
	if persisted.Request.RunID != receipt.RunID ||
		persisted.Request.TaskID != receipt.TaskID ||
		persisted.Request.StepID != receipt.StepID ||
		persisted.Request.Attempt != receipt.Attempt ||
		persisted.Request.AdapterID != receipt.AdapterID ||
		persisted.Request.AuthorizationHash != receipt.AuthorizationHash ||
		persisted.Request.BudgetHash != receipt.BudgetHash {
		return fmt.Errorf("%w: launch receipt binding mismatch for requestId %s", ErrAgentRunnerLaunchReceiptConflict, receipt.RequestID)
	}
	return nil
}

func (store *AgentRunnerReplayStore) launchIdentityPath(launchID string) (string, error) {
	return store.pathForRequestID(launchID, "launches", "identity.")
}

func loadAgentRunnerLaunchReceipt(path, requestID string) (AgentRunnerLaunchReceiptRecord, bool, error) {
	wire, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return AgentRunnerLaunchReceiptRecord{}, false, nil
	}
	if err != nil {
		return AgentRunnerLaunchReceiptRecord{}, false, err
	}
	record, err := DecodeAgentRunnerLaunchReceiptRecord(wire)
	if err != nil {
		return AgentRunnerLaunchReceiptRecord{}, false, replayStoreCorruption(path, "launch receipt decode failed", err)
	}
	if record.Receipt.RequestID != requestID {
		return AgentRunnerLaunchReceiptRecord{}, false, replayStoreCorruption(path, fmt.Sprintf("launch receipt requestId mismatch on disk (want %s, got %s)", requestID, record.Receipt.RequestID), nil)
	}
	return record, true, nil
}

func loadAgentRunnerLaunchIdentity(path, launchID string) (AgentRunnerLaunchIdentityRecord, bool, error) {
	wire, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return AgentRunnerLaunchIdentityRecord{}, false, nil
	}
	if err != nil {
		return AgentRunnerLaunchIdentityRecord{}, false, err
	}
	record, err := DecodeAgentRunnerLaunchIdentityRecord(wire)
	if err != nil {
		return AgentRunnerLaunchIdentityRecord{}, false, replayStoreCorruption(path, "launch identity decode failed", err)
	}
	if record.Identity.LaunchID != launchID {
		return AgentRunnerLaunchIdentityRecord{}, false, replayStoreCorruption(path, fmt.Sprintf("launch identity launchId mismatch on disk (want %s, got %s)", launchID, record.Identity.LaunchID), nil)
	}
	return record, true, nil
}
