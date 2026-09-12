package adapters

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const replayStoreReadCollisionRetries = 16

var (
	// ErrInvalidAgentRunnerReplayStore reports invalid replay-store construction or usage.
	ErrInvalidAgentRunnerReplayStore = errors.New("invalid AgentRunner replay store")
	// ErrAgentRunnerReplayStoreCorruption reports fail-closed on-disk replay-store corruption.
	ErrAgentRunnerReplayStoreCorruption = errors.New("AgentRunner replay store corruption")
)

// AgentRunnerReplayState describes durable replay-state visibility for one requestId.
type AgentRunnerReplayState string

const (
	// AgentRunnerReplayStateAbsent means neither request nor completion was found.
	AgentRunnerReplayStateAbsent AgentRunnerReplayState = "absent"
	// AgentRunnerReplayStateDispatchedUnknown means request exists but no terminal receipt exists yet.
	AgentRunnerReplayStateDispatchedUnknown AgentRunnerReplayState = "dispatched-unknown"
	// AgentRunnerReplayStateTerminalReceiptPresent means a terminal completion receipt exists.
	// Any completion status, including uncertain, only proves terminal-receipt presence and never implies task PASS.
	AgentRunnerReplayStateTerminalReceiptPresent AgentRunnerReplayState = "terminal-receipt-present"
)

// AgentRunnerReplayDecision describes dispatch gating for a request write attempt.
type AgentRunnerReplayDecision string

const (
	// AgentRunnerReplayDecisionFirstDispatch means this call durably published the first request record.
	AgentRunnerReplayDecisionFirstDispatch AgentRunnerReplayDecision = "first-dispatch"
	// AgentRunnerReplayDecisionUnknownBlock means a request record exists without a terminal receipt.
	AgentRunnerReplayDecisionUnknownBlock AgentRunnerReplayDecision = "unknown-block"
	// AgentRunnerReplayDecisionTerminalReceiptPresent means a terminal receipt already exists for the request.
	// Any completion status, including uncertain, only proves terminal-receipt presence and never implies task PASS.
	AgentRunnerReplayDecisionTerminalReceiptPresent AgentRunnerReplayDecision = "terminal-receipt-present"
)

// AgentRunnerReplayStore persists request/completion replay records under one absolute root.
type AgentRunnerReplayStore struct {
	root string
	mu   sync.Mutex
}

// NewAgentRunnerReplayStore constructs a replay store rooted at an absolute path.
// Records are stored under requests/ and completions/ subdirectories.
func NewAgentRunnerReplayStore(root string) (*AgentRunnerReplayStore, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("%w: empty replay root", ErrInvalidAgentRunnerReplayStore)
	}
	cleanRoot := filepath.Clean(root)
	if !filepath.IsAbs(cleanRoot) {
		return nil, fmt.Errorf("%w: replay root must be absolute", ErrInvalidAgentRunnerReplayStore)
	}
	return &AgentRunnerReplayStore{root: cleanRoot}, nil
}

// State returns the durable replay state for requestID.
// It fails closed if an orphan completion exists without a request record.
func (store *AgentRunnerReplayStore) State(requestID string) (AgentRunnerReplayState, error) {
	if store == nil {
		return "", fmt.Errorf("%w: nil replay store", ErrInvalidAgentRunnerReplayStore)
	}
	store.mu.Lock()
	defer store.mu.Unlock()

	request, foundRequest, err := store.loadRequestLocked(requestID)
	if err != nil {
		return "", err
	}
	if !foundRequest {
		_, foundCompletion, err := store.loadCompletionLocked(requestID)
		if err != nil {
			return "", err
		}
		if foundCompletion {
			path, pathErr := store.completionPath(requestID)
			if pathErr != nil {
				return "", pathErr
			}
			return "", replayStoreCorruption(path, fmt.Sprintf("orphan completion present for requestId %s", requestID), nil)
		}
		return AgentRunnerReplayStateAbsent, nil
	}
	completion, foundCompletion, err := store.loadCompletionLocked(requestID)
	if err != nil {
		return "", err
	}
	if !foundCompletion {
		return AgentRunnerReplayStateDispatchedUnknown, nil
	}
	if err := ValidateAgentRunnerCompletionBinding(request, completion); err != nil {
		return "", err
	}
	return AgentRunnerReplayStateTerminalReceiptPresent, nil
}

// RecordRequest records one request and returns the replay decision for dispatch gating.
func (store *AgentRunnerReplayStore) RecordRequest(record AgentRunnerRequestRecord) (AgentRunnerReplayDecision, error) {
	if store == nil {
		return "", fmt.Errorf("%w: nil replay store", ErrInvalidAgentRunnerReplayStore)
	}
	if err := ValidateAgentRunnerRequestRecord(record); err != nil {
		return "", err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	request, inserted, err := store.ensureRequestLocked(record)
	if err != nil {
		return "", err
	}
	completion, foundCompletion, err := store.loadCompletionLocked(record.Request.RequestID)
	if err != nil {
		return "", err
	}
	if !foundCompletion {
		if inserted {
			return AgentRunnerReplayDecisionFirstDispatch, nil
		}
		return AgentRunnerReplayDecisionUnknownBlock, nil
	}
	if err := ValidateAgentRunnerCompletionBinding(request, completion); err != nil {
		return "", err
	}
	return AgentRunnerReplayDecisionTerminalReceiptPresent, nil
}

// RecordCompletion records one terminal completion receipt.
// The returned bool is true only when the same completion record was already durably present (replay).
// The returned bool is false when this call durably published a new completion record.
func (store *AgentRunnerReplayStore) RecordCompletion(request AgentRunnerRequestRecord, completion AgentRunnerCompletionRecord) (bool, error) {
	if store == nil {
		return false, fmt.Errorf("%w: nil replay store", ErrInvalidAgentRunnerReplayStore)
	}
	if err := ValidateAgentRunnerRequestRecord(request); err != nil {
		return false, err
	}
	if err := ValidateAgentRunnerCompletionRecord(completion); err != nil {
		return false, err
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	requestID := request.Request.RequestID
	persisted, foundRequest, err := store.loadRequestLocked(requestID)
	if err != nil {
		return false, err
	}
	existing, foundCompletion, err := store.loadCompletionLocked(requestID)
	if err != nil {
		return false, err
	}

	// A terminal receipt that is already on disk is resolved without any write,
	// so a rejected or replayed completion can never consume the request identity.
	if foundCompletion {
		if !foundRequest {
			path, pathErr := store.completionPath(requestID)
			if pathErr != nil {
				return false, pathErr
			}
			return false, replayStoreCorruption(path, fmt.Sprintf("orphan completion present for requestId %s", requestID), nil)
		}
		if persisted.RecordHash != request.RecordHash {
			return false, fmt.Errorf("%w: requestId %s", ErrAgentRunnerRequestConflict, requestID)
		}
		if existing.RecordHash != completion.RecordHash {
			return false, fmt.Errorf("%w: requestId %s already completed", ErrAgentRunnerCompletionConflict, requestID)
		}
		if err := ValidateAgentRunnerCompletionBinding(persisted, existing); err != nil {
			return false, err
		}
		return true, nil
	}

	// No terminal receipt exists yet. The binding must hold before the first
	// durable write, so no rejected completion leaves a request behind.
	if foundRequest {
		if persisted.RecordHash != request.RecordHash {
			return false, fmt.Errorf("%w: requestId %s", ErrAgentRunnerRequestConflict, requestID)
		}
		if err := ValidateAgentRunnerCompletionBinding(persisted, completion); err != nil {
			return false, err
		}
	} else if err := ValidateAgentRunnerCompletionBinding(request, completion); err != nil {
		return false, err
	}

	persistedRequest, _, err := store.ensureRequestLocked(request)
	if err != nil {
		return false, err
	}
	path, err := store.completionPath(requestID)
	if err != nil {
		return false, err
	}
	collided := false
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		if !collided {
			err = writeCanonicalLineNoReplace(path, completion)
			if err == nil {
				return false, nil
			}
			if !errors.Is(err, fs.ErrExist) {
				return false, err
			}
			// A no-replace collision proves the terminal slot is already owned.
			// Only re-reads follow; retrying the write could take over the slot
			// during a visibility gap and replace a different terminal receipt.
			collided = true
		}
		visible, found, loadErr := loadCompletionRecord(path, requestID)
		if loadErr != nil {
			return false, loadErr
		}
		if !found {
			runtime.Gosched()
			continue
		}
		if visible.RecordHash != completion.RecordHash {
			return false, fmt.Errorf("%w: requestId %s already completed", ErrAgentRunnerCompletionConflict, requestID)
		}
		if err := ValidateAgentRunnerCompletionBinding(persistedRequest, visible); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, fmt.Errorf("%w: completion file not visible after no-replace collision for requestId %s", ErrAgentRunnerCompletionConflict, requestID)
}

func (store *AgentRunnerReplayStore) ensureRequestLocked(record AgentRunnerRequestRecord) (AgentRunnerRequestRecord, bool, error) {
	path, err := store.requestPath(record.Request.RequestID)
	if err != nil {
		return AgentRunnerRequestRecord{}, false, err
	}
	collided := false
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		if !collided {
			err = writeCanonicalLineNoReplace(path, record)
			if err == nil {
				return record, true, nil
			}
			if !errors.Is(err, fs.ErrExist) {
				return AgentRunnerRequestRecord{}, false, err
			}
			// The request slot is already owned. Only re-reads follow so a
			// visibility gap cannot replace the durable request record.
			collided = true
		}
		existing, found, loadErr := loadRequestRecord(path, record.Request.RequestID)
		if loadErr != nil {
			return AgentRunnerRequestRecord{}, false, loadErr
		}
		if !found {
			runtime.Gosched()
			continue
		}
		if existing.RecordHash != record.RecordHash {
			return AgentRunnerRequestRecord{}, false, fmt.Errorf("%w: requestId %s", ErrAgentRunnerRequestConflict, record.Request.RequestID)
		}
		return existing, false, nil
	}
	return AgentRunnerRequestRecord{}, false, fmt.Errorf("%w: request file not visible after no-replace collision for requestId %s", ErrAgentRunnerRequestConflict, record.Request.RequestID)
}

func (store *AgentRunnerReplayStore) loadRequestLocked(requestID string) (AgentRunnerRequestRecord, bool, error) {
	path, err := store.requestPath(requestID)
	if err != nil {
		return AgentRunnerRequestRecord{}, false, err
	}
	return loadRequestRecord(path, requestID)
}

func (store *AgentRunnerReplayStore) loadCompletionLocked(requestID string) (AgentRunnerCompletionRecord, bool, error) {
	path, err := store.completionPath(requestID)
	if err != nil {
		return AgentRunnerCompletionRecord{}, false, err
	}
	return loadCompletionRecord(path, requestID)
}

func (store *AgentRunnerReplayStore) requestPath(requestID string) (string, error) {
	return store.pathForRequestID(requestID, "requests", "request.")
}

func (store *AgentRunnerReplayStore) completionPath(requestID string) (string, error) {
	return store.pathForRequestID(requestID, "completions", "completion.")
}

func (store *AgentRunnerReplayStore) pathForRequestID(requestID, folder, prefix string) (string, error) {
	if !evidence.ValidID(requestID) {
		return "", fmt.Errorf("%w: invalid requestId", ErrInvalidAgentRunnerRequest)
	}
	filename := prefix + requestID + ".jsonl"
	root := filepath.Clean(store.root)
	path := filepath.Clean(filepath.Join(root, folder, filename))
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return "", err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("%w: requestId escaped replay root", ErrInvalidAgentRunnerReplayStore)
	}
	return path, nil
}

func loadRequestRecord(path, requestID string) (AgentRunnerRequestRecord, bool, error) {
	wire, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return AgentRunnerRequestRecord{}, false, nil
	}
	if err != nil {
		return AgentRunnerRequestRecord{}, false, err
	}
	record, err := DecodeAgentRunnerRequestRecord(wire)
	if err != nil {
		return AgentRunnerRequestRecord{}, false, replayStoreCorruption(path, "request record decode failed", err)
	}
	if record.Request.RequestID != requestID {
		return AgentRunnerRequestRecord{}, false, replayStoreCorruption(path, fmt.Sprintf("requestId mismatch on disk (want %s, got %s)", requestID, record.Request.RequestID), nil)
	}
	return record, true, nil
}

func loadCompletionRecord(path, requestID string) (AgentRunnerCompletionRecord, bool, error) {
	wire, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return AgentRunnerCompletionRecord{}, false, nil
	}
	if err != nil {
		return AgentRunnerCompletionRecord{}, false, err
	}
	record, err := DecodeAgentRunnerCompletionRecord(wire)
	if err != nil {
		return AgentRunnerCompletionRecord{}, false, replayStoreCorruption(path, "completion record decode failed", err)
	}
	if record.Completion.RequestID != requestID {
		return AgentRunnerCompletionRecord{}, false, replayStoreCorruption(path, fmt.Sprintf("completion requestId mismatch on disk (want %s, got %s)", requestID, record.Completion.RequestID), nil)
	}
	return record, true, nil
}

func replayStoreCorruption(path, detail string, cause error) error {
	if cause == nil {
		return fmt.Errorf("%w: %s at %s", ErrAgentRunnerReplayStoreCorruption, detail, path)
	}
	return fmt.Errorf("%w: %s at %s: %w", ErrAgentRunnerReplayStoreCorruption, detail, path, cause)
}
