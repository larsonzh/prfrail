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
	// ErrAgentRunnerReplayStoreConvergence reports bounded replay-store rereads that
	// could not truthfully classify the winner after a no-replace collision.
	ErrAgentRunnerReplayStoreConvergence = errors.New("AgentRunner replay store convergence failure")
	// ErrAgentRunnerReplayStoreCorruption reports fail-closed on-disk replay-store corruption.
	ErrAgentRunnerReplayStoreCorruption = errors.New("AgentRunner replay store corruption")
	// ErrAgentRunnerReplayStoreDurabilityUnproven reports a write that cannot
	// prove replay publication durability on the current store instance.
	ErrAgentRunnerReplayStoreDurabilityUnproven = errors.New("AgentRunner replay store publication durability unproven")
)

// PublishDurability describes whether replay publication can be proved durable on
// the current platform implementation after the required post-publish steps.
type PublishDurability string

const (
	// PublishDurabilityProven means the implementation can prove the published
	// directory entry reached the local durability boundary required by contract.
	PublishDurabilityProven PublishDurability = "proven"
	// PublishDurabilityUnproven means publication may be atomically visible but
	// the implementation cannot prove directory-entry durability.
	PublishDurabilityUnproven PublishDurability = "unproven"
)

var (
	replayStoreLoadRequestRecord    = loadRequestRecord
	replayStoreLoadCompletionRecord = loadCompletionRecord
	replayStoreLoadLaunchReceipt    = loadAgentRunnerLaunchReceipt
	replayStoreLoadLaunchIdentity   = loadAgentRunnerLaunchIdentity
	replayStoreLoadTerminalIntent   = loadAgentRunnerTerminalIntent
	replayStoreLoadTerminalClosure  = loadAgentRunnerTerminalClosure
	replayStoreLoadOwnershipRecord  = loadAgentRunnerReplayOwnership
	replayStoreSyncParentDirectory  = syncReplayStoreParentDirectory
)

// Publish stages of writeReplayRecordNoReplace, in execution order. They are the
// crash points an experiment can stop at: everything before the record is linked
// into place is invisible to a reader, and the two stages after it differ only in
// whether the parent-directory durability step ran.
const (
	replayStorePublishStageTempWritten  = "temp-written"
	replayStorePublishStageTempSynced   = "temp-synced"
	replayStorePublishStageTempClosed   = "temp-closed"
	replayStorePublishStageRecordLinked = "record-linked"
	replayStorePublishStageParentSynced = "parent-synced"
)

// replayStorePublishStageHook is a test seam for the publication-durability
// falsification prototype. It is nil in production: a timing-based kill cannot
// stop a publish at an exact stage, so crash-point experiments need a seam that
// can. It never changes behaviour while nil, and it may not carry decisions.
//
// Seam contract: it fires for every no-replace publication (requests,
// completions, launch receipts and identities, terminal intents and closures,
// and the ownership record - the last of which is published before the store is
// constructed, so not every call runs under the store lock), and a hook runs
// while its caller holds the store lock whenever that caller is a store method,
// so a hook must not call back into the store. A failed step fires no further
// stage: the stages before the failure have already fired.
var replayStorePublishStageHook func(stage string, path string)

func replayStorePublishStage(stage string, path string) {
	if replayStorePublishStageHook != nil {
		replayStorePublishStageHook(stage, path)
	}
}

// replayStorePublishPrimitiveHook is a test seam for the B3b durability
// candidate experiments: it replaces only the no-replace primitive that makes a
// published record visible (the step os.Link performs today). It is nil in
// production, and while it is nil the primitive is os.Link, byte for byte the
// behaviour the store has always had.
//
// Seam contract: it fires once per no-replace publication, after the record is
// written, synced and closed and before the parent-directory durability step. A
// replacement must keep the no-replace contract (an existing target has to
// surface as fs.ErrExist) and must leave the target holding exactly the bytes the
// store already wrote. It may not carry decisions and may not call back into the
// store.
var replayStorePublishPrimitiveHook func(temp string, path string) error

func replayStorePublishRecordNoReplace(temp string, path string) error {
	if replayStorePublishPrimitiveHook != nil {
		return replayStorePublishPrimitiveHook(temp, path)
	}
	return os.Link(temp, path)
}

// replayStoreParentSyncHook is the second B3b test seam: it replaces only the
// parent-directory durability step, the step syncReplayStoreParentDirectory
// performs on Unix and skips on Windows. It is nil in production, and while it is
// nil the behaviour is exactly the platform default.
//
// Seam contract: it fires once per no-replace publication, after the no-replace
// primitive published the record and after the record-linked stage, and it
// receives the directory holding the record. A failure aborts the publication:
// the caller drops its temp file and reports the error, while the already
// published record stays as visible as the primitive made it - visibility and
// durability are separate facts. It may not carry decisions and may not call back
// into the store.
//
// It is deliberately separate from replayStoreSyncParentDirectory, which existing
// tests replace to model the platform primitives, and it is deliberately not
// wired into the bootstrap path: bootstrap is not a publication, and on a platform
// whose durability is unproven the bootstrap durability step does not run at all.
var replayStoreParentSyncHook func(directory string) error

func replayStoreSyncParent(directory string) error {
	if replayStoreParentSyncHook != nil {
		return replayStoreParentSyncHook(directory)
	}
	return replayStoreSyncParentDirectory(directory)
}

// AgentRunnerReplayState describes replay-state visibility for one requestId.
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
	// AgentRunnerReplayDecisionFirstDispatch means this call published the first request record.
	AgentRunnerReplayDecisionFirstDispatch AgentRunnerReplayDecision = "first-dispatch"
	// AgentRunnerReplayDecisionUnknownBlock means a request record exists without a terminal receipt.
	AgentRunnerReplayDecisionUnknownBlock AgentRunnerReplayDecision = "unknown-block"
	// AgentRunnerReplayDecisionTerminalReceiptPresent means a terminal receipt already exists for the request.
	// Any completion status, including uncertain, only proves terminal-receipt presence and never implies task PASS.
	AgentRunnerReplayDecisionTerminalReceiptPresent AgentRunnerReplayDecision = "terminal-receipt-present"
)

// AgentRunnerReplayStore persists request/completion replay records under one absolute root.
type AgentRunnerReplayStore struct {
	root                   string
	runID                  string
	runRoot                string
	publicationDurability  PublishDurability
	mu                     sync.Mutex
	afterRequestReadLocked func()
}

// newAgentRunnerReplayStoreAt constructs a replay store at an exact absolute
// root. It is not a production entry point: production callers must use
// NewAgentRunnerReplayStoreForRun so the replay root is derived from a durable
// run root instead of an arbitrary caller-chosen path.
// Records are stored under requests/ and completions/ subdirectories.
func newAgentRunnerReplayStoreAt(root string) (*AgentRunnerReplayStore, error) {
	if strings.TrimSpace(root) == "" {
		return nil, fmt.Errorf("%w: empty replay root", ErrInvalidAgentRunnerReplayStore)
	}
	cleanRoot := filepath.Clean(root)
	if !filepath.IsAbs(cleanRoot) {
		return nil, fmt.Errorf("%w: replay root must be absolute", ErrInvalidAgentRunnerReplayStore)
	}
	store := &AgentRunnerReplayStore{
		root:                  cleanRoot,
		publicationDurability: replayStorePublicationDurability(),
	}
	if store.publicationDurability == PublishDurabilityProven {
		if err := bootstrapReplayStoreRoot(cleanRoot); err != nil {
			return nil, fmt.Errorf("%w: replay root bootstrap: %w", ErrInvalidAgentRunnerReplayStore, err)
		}
	}
	return store, nil
}

// validateConstructed fails closed for nil stores and for zero-value stores
// whose root was never set, so read paths cannot silently operate on the
// process working directory. Write paths keep the contract-mandated ordering
// where the durability gate rejects zero-value stores first.
func (store *AgentRunnerReplayStore) validateConstructed() error {
	if store == nil {
		return fmt.Errorf("%w: nil replay store", ErrInvalidAgentRunnerReplayStore)
	}
	if strings.TrimSpace(store.root) == "" {
		return fmt.Errorf("%w: replay store root is unset", ErrInvalidAgentRunnerReplayStore)
	}
	return nil
}

// PublicationDurability reports whether this platform-specific replay-store
// implementation can prove request/completion publication durability.
func (store *AgentRunnerReplayStore) PublicationDurability() PublishDurability {
	if store == nil || store.publicationDurability != PublishDurabilityProven {
		return PublishDurabilityUnproven
	}
	return PublishDurabilityProven
}

// State returns the replay state visible for requestID.
// It fails closed if an orphan completion exists without a request record.
func (store *AgentRunnerReplayStore) State(requestID string) (AgentRunnerReplayState, error) {
	if err := store.validateConstructed(); err != nil {
		return "", err
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return "", err
	}

	request, foundRequest, completion, foundCompletion, err := store.loadConvergedRecordsLocked(requestID)
	if err != nil {
		return "", err
	}
	if !foundRequest {
		if foundCompletion {
			return "", store.orphanCompletionError(requestID)
		}
		return AgentRunnerReplayStateAbsent, nil
	}
	if !foundCompletion {
		return AgentRunnerReplayStateDispatchedUnknown, nil
	}
	if err := store.validatePersistedCompletionBinding(requestID, request, completion); err != nil {
		return "", err
	}
	return AgentRunnerReplayStateTerminalReceiptPresent, nil
}

// RecordRequest records one request and returns the replay decision for dispatch gating.
func (store *AgentRunnerReplayStore) RecordRequest(record AgentRunnerRequestRecord) (AgentRunnerReplayDecision, error) {
	if store == nil {
		return "", fmt.Errorf("%w: nil replay store", ErrInvalidAgentRunnerReplayStore)
	}
	if store.PublicationDurability() != PublishDurabilityProven {
		return "", fmt.Errorf("%w: request publication", ErrAgentRunnerReplayStoreDurabilityUnproven)
	}
	if err := ValidateAgentRunnerRequestRecord(record); err != nil {
		return "", err
	}
	if err := store.validateRecordRunBinding(record.Request.RunID, record.Request.RequestID); err != nil {
		return "", err
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return "", err
	}

	_, foundRequest, _, foundCompletion, err := store.loadConvergedRecordsLocked(record.Request.RequestID)
	if err != nil {
		return "", err
	}
	if !foundRequest && foundCompletion {
		return "", store.orphanCompletionError(record.Request.RequestID)
	}

	request, inserted, err := store.ensureRequestLocked(record)
	if err != nil {
		return "", err
	}
	_, _, completion, foundCompletion, err := store.loadConvergedRecordsLocked(record.Request.RequestID)
	if err != nil {
		return "", err
	}
	if !foundCompletion {
		if inserted {
			return AgentRunnerReplayDecisionFirstDispatch, nil
		}
		return AgentRunnerReplayDecisionUnknownBlock, nil
	}
	if err := store.validatePersistedCompletionBinding(record.Request.RequestID, request, completion); err != nil {
		return "", err
	}
	return AgentRunnerReplayDecisionTerminalReceiptPresent, nil
}

// RecordCompletion records one terminal completion receipt.
// The returned bool is true only when the same completion record was already present (replay).
// The returned bool is false when this call published a new completion record.
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
	if err := store.validateRecordRunBinding(request.Request.RunID, request.Request.RequestID); err != nil {
		return false, err
	}
	if err := store.validateRecordRunBinding(completion.Completion.RunID, completion.Completion.RequestID); err != nil {
		return false, err
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return false, err
	}

	requestID := request.Request.RequestID
	persisted, foundRequest, existing, foundCompletion, err := store.loadConvergedRecordsLocked(requestID)
	if err != nil {
		return false, err
	}

	// A terminal receipt that is already on disk is resolved without any write,
	// so a rejected or replayed completion can never consume the request identity.
	if foundCompletion {
		if !foundRequest {
			return false, store.orphanCompletionError(requestID)
		}
		if persisted.RecordHash != request.RecordHash {
			return false, fmt.Errorf("%w: requestId %s", ErrAgentRunnerRequestConflict, requestID)
		}
		if existing.RecordHash != completion.RecordHash {
			return false, fmt.Errorf("%w: requestId %s already completed", ErrAgentRunnerCompletionConflict, requestID)
		}
		if err := store.validatePersistedCompletionBinding(requestID, persisted, existing); err != nil {
			return false, err
		}
		return true, nil
	}
	if store.PublicationDurability() != PublishDurabilityProven {
		return false, fmt.Errorf("%w: completion publication", ErrAgentRunnerReplayStoreDurabilityUnproven)
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
			err = writeReplayRecordNoReplace(path, completion)
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
		visible, found, loadErr := replayStoreLoadCompletionRecord(path, requestID)
		if loadErr != nil {
			return false, loadErr
		}
		if !found {
			runtime.Gosched()
			continue
		}
		if err := store.validatePersistedRunID(path, "completion", visible.Completion.RunID); err != nil {
			return false, err
		}
		if visible.RecordHash != completion.RecordHash {
			return false, fmt.Errorf("%w: requestId %s already completed", ErrAgentRunnerCompletionConflict, requestID)
		}
		if err := store.validatePersistedCompletionBinding(requestID, persistedRequest, visible); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, fmt.Errorf("%w: completion winner not visible after no-replace collision for requestId %s", ErrAgentRunnerReplayStoreConvergence, requestID)
}

func bootstrapReplayStoreRoot(root string) error {
	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("replay root is not a directory: %s", root)
	}
	for _, name := range []string{"requests", "completions", "launches", "terminals"} {
		directory := filepath.Join(root, name)
		if err := os.Mkdir(directory, 0o755); err != nil && !errors.Is(err, fs.ErrExist) {
			return err
		}
		info, err := os.Stat(directory)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("replay store path is not a directory: %s", directory)
		}
		if err := replayStoreSyncParentDirectory(root); err != nil {
			return err
		}
	}
	return nil
}

func (store *AgentRunnerReplayStore) ensureRequestLocked(record AgentRunnerRequestRecord) (AgentRunnerRequestRecord, bool, error) {
	path, err := store.requestPath(record.Request.RequestID)
	if err != nil {
		return AgentRunnerRequestRecord{}, false, err
	}
	collided := false
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		if !collided {
			err = writeReplayRecordNoReplace(path, record)
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
		existing, found, loadErr := replayStoreLoadRequestRecord(path, record.Request.RequestID)
		if loadErr != nil {
			return AgentRunnerRequestRecord{}, false, loadErr
		}
		if !found {
			runtime.Gosched()
			continue
		}
		if err := store.validatePersistedRunID(path, "request", existing.Request.RunID); err != nil {
			return AgentRunnerRequestRecord{}, false, err
		}
		if existing.RecordHash != record.RecordHash {
			return AgentRunnerRequestRecord{}, false, fmt.Errorf("%w: requestId %s", ErrAgentRunnerRequestConflict, record.Request.RequestID)
		}
		return existing, false, nil
	}
	return AgentRunnerRequestRecord{}, false, fmt.Errorf("%w: request winner not visible after no-replace collision for requestId %s", ErrAgentRunnerReplayStoreConvergence, record.Request.RequestID)
}

func (store *AgentRunnerReplayStore) loadRequestLocked(requestID string) (AgentRunnerRequestRecord, bool, error) {
	path, err := store.requestPath(requestID)
	if err != nil {
		return AgentRunnerRequestRecord{}, false, err
	}
	record, found, err := replayStoreLoadRequestRecord(path, requestID)
	if err != nil || !found {
		return record, found, err
	}
	if err := store.validatePersistedRunID(path, "request", record.Request.RunID); err != nil {
		return AgentRunnerRequestRecord{}, false, err
	}
	return record, true, nil
}

func (store *AgentRunnerReplayStore) loadCompletionLocked(requestID string) (AgentRunnerCompletionRecord, bool, error) {
	path, err := store.completionPath(requestID)
	if err != nil {
		return AgentRunnerCompletionRecord{}, false, err
	}
	record, found, err := replayStoreLoadCompletionRecord(path, requestID)
	if err != nil || !found {
		return record, found, err
	}
	if err := store.validatePersistedRunID(path, "completion", record.Completion.RunID); err != nil {
		return AgentRunnerCompletionRecord{}, false, err
	}
	return record, true, nil
}

func (store *AgentRunnerReplayStore) loadConvergedRecordsLocked(requestID string) (AgentRunnerRequestRecord, bool, AgentRunnerCompletionRecord, bool, error) {
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		request, foundRequest, err := store.loadRequestLocked(requestID)
		if err != nil {
			return AgentRunnerRequestRecord{}, false, AgentRunnerCompletionRecord{}, false, err
		}
		if store.afterRequestReadLocked != nil {
			hook := store.afterRequestReadLocked
			store.afterRequestReadLocked = nil
			hook()
		}
		completion, foundCompletion, err := store.loadCompletionLocked(requestID)
		if err != nil {
			return AgentRunnerRequestRecord{}, false, AgentRunnerCompletionRecord{}, false, err
		}
		if foundRequest || !foundCompletion {
			return request, foundRequest, completion, foundCompletion, nil
		}
		runtime.Gosched()
	}
	return AgentRunnerRequestRecord{}, false, AgentRunnerCompletionRecord{}, true, nil
}

func (store *AgentRunnerReplayStore) orphanCompletionError(requestID string) error {
	path, err := store.completionPath(requestID)
	if err != nil {
		return err
	}
	return replayStoreCorruption(path, fmt.Sprintf("orphan completion present for requestId %s", requestID), nil)
}

func (store *AgentRunnerReplayStore) validatePersistedCompletionBinding(requestID string, request AgentRunnerRequestRecord, completion AgentRunnerCompletionRecord) error {
	if err := ValidateAgentRunnerCompletionBinding(request, completion); err != nil {
		path, pathErr := store.completionPath(requestID)
		if pathErr != nil {
			return pathErr
		}
		return replayStoreCorruption(path, "request/completion binding mismatch on disk", err)
	}
	return nil
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

func writeReplayRecordNoReplace(path string, record any) error {
	canonical, err := evidence.EncodeCanonical(record)
	if err != nil {
		return err
	}
	payload := append(canonical, '\n')
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	file, err := os.CreateTemp(directory, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	temp := file.Name()
	writeErr := error(nil)
	if _, err := file.Write(payload); err != nil {
		writeErr = err
	} else {
		replayStorePublishStage(replayStorePublishStageTempWritten, path)
		if err := file.Sync(); err != nil {
			writeErr = err
		} else {
			replayStorePublishStage(replayStorePublishStageTempSynced, path)
		}
	}
	closeErr := file.Close()
	if writeErr == nil {
		writeErr = closeErr
	}
	if writeErr != nil {
		_ = os.Remove(temp)
		return writeErr
	}
	replayStorePublishStage(replayStorePublishStageTempClosed, path)
	// Atomic no-replace visibility and durability are separate facts. The
	// no-replace primitive publishes visibility without replacement; only a
	// successful parent directory durability step may later justify a proven
	// durability claim.
	if err := replayStorePublishRecordNoReplace(temp, path); err != nil {
		_ = os.Remove(temp)
		if errors.Is(err, fs.ErrExist) {
			return fs.ErrExist
		}
		return err
	}
	replayStorePublishStage(replayStorePublishStageRecordLinked, path)
	if err := replayStoreSyncParent(directory); err != nil {
		_ = os.Remove(temp)
		return err
	}
	replayStorePublishStage(replayStorePublishStageParentSynced, path)
	_ = os.Remove(temp)
	return nil
}

func replayStoreCorruption(path, detail string, cause error) error {
	if cause == nil {
		return fmt.Errorf("%w: %s at %s", ErrAgentRunnerReplayStoreCorruption, detail, path)
	}
	return fmt.Errorf("%w: %s at %s: %w", ErrAgentRunnerReplayStoreCorruption, detail, path, cause)
}
