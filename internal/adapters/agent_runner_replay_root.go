package adapters

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	agentRunnerReplayRootDirectoryName = "agent-runner-replay"
	agentRunnerReplayOwnershipFileName = "ownership.json"
	agentRunnerReplayOwnershipKind     = "agent-runner-replay-ownership"
	agentRunnerReplayOwnershipDomain   = "proofrail:agent-runner-replay-ownership:1\n"
	agentRunnerReplayRunEventDirName   = "events"
	agentRunnerReplayRunEventFileName  = "state-events.jsonl"
)

var (
	// ErrAgentRunnerReplayRootConflict reports a replay root bound to a different
	// run, moved or copied after creation, or overlapping a writer or protected root.
	ErrAgentRunnerReplayRootConflict = errors.New("AgentRunner replay root conflict")
	// ErrAgentRunnerReplayRootUnsafePath reports a replay-root path component that
	// is a symlink or reparse point, or whose safety cannot be decided.
	ErrAgentRunnerReplayRootUnsafePath = errors.New("AgentRunner replay root unsafe path")
)

// AgentRunnerReplayOwnershipRecord is store-local metadata published once per
// replay root. It never enters the wire contract and never participates in any
// R/C recordHash.
type AgentRunnerReplayOwnershipRecord struct {
	Kind          string `json:"kind"`
	SchemaVersion string `json:"schemaVersion"`
	RunID         string `json:"runId"`
	RunRootPath   string `json:"runRootPath"`
	RunRootHash   string `json:"runRootHash"`
	CreatedAt     string `json:"createdAt"`
}

// NewAgentRunnerReplayStoreForRun is the production replay-store constructor.
// It derives the replay root from the durable run root and binds the store to
// runID instead of accepting an arbitrary caller-chosen root. The run root must
// already exist and carry the durable run-root marker written by the run event
// store before any AgentRunner step can exist.
func NewAgentRunnerReplayStoreForRun(runRoot, runID string) (*AgentRunnerReplayStore, error) {
	cleanRoot, err := validateAgentRunnerReplayRunRoot(runRoot)
	if err != nil {
		return nil, err
	}
	if !evidence.ValidID(runID) {
		return nil, fmt.Errorf("%w: invalid runId", ErrInvalidAgentRunnerReplayStore)
	}
	if err := validateAgentRunnerRunRootMarker(cleanRoot); err != nil {
		return nil, err
	}
	if err := rejectUnsafeReplayPathComponents(cleanRoot); err != nil {
		return nil, err
	}
	replayRoot := filepath.Join(cleanRoot, agentRunnerReplayRootDirectoryName)
	if err := validateAgentRunnerReplayRootInvariants(replayRoot, cleanRoot); err != nil {
		return nil, err
	}
	if err := rejectUnsafeReplayPathComponents(replayRoot); err != nil {
		return nil, err
	}
	if err := ensureAgentRunnerReplayDirectory(replayRoot, cleanRoot); err != nil {
		return nil, err
	}
	// Subdirectory safety checks, bootstrap and the low-level constructor all
	// complete before ownership is published, so a failing construction never
	// leaves ownership metadata behind.
	subdirectories := []string{
		filepath.Join(replayRoot, "requests"),
		filepath.Join(replayRoot, "completions"),
	}
	for _, directory := range subdirectories {
		if err := rejectUnsafeReplayPathComponents(directory); err != nil {
			return nil, err
		}
		if err := ensureAgentRunnerReplayDirectory(directory, replayRoot); err != nil {
			return nil, err
		}
	}
	store, err := newAgentRunnerReplayStoreAt(replayRoot)
	if err != nil {
		return nil, err
	}
	if err := ensureAgentRunnerReplayOwnership(replayRoot, cleanRoot, runID); err != nil {
		return nil, err
	}
	store.runID = runID
	store.runRoot = cleanRoot
	return store, nil
}

// RejectWriterRoots fails closed when any writer root (an isolated workspace)
// is equal to or nested with the replay root in either direction.
func (store *AgentRunnerReplayStore) RejectWriterRoots(roots ...string) error {
	if store == nil {
		return fmt.Errorf("%w: nil replay store", ErrInvalidAgentRunnerReplayStore)
	}
	replayRoot, err := store.comparableReplayRoot()
	if err != nil {
		return err
	}
	for _, root := range roots {
		if err := rejectUnsafeReplayPathComponents(root); err != nil {
			return err
		}
		comparable, err := resolveComparableReplayPath(root)
		if err != nil {
			return replayRootConflict(root, "writer root cannot be resolved", err)
		}
		if replayPathContainsOrAliases(comparable, replayRoot) || replayPathContainsOrAliases(replayRoot, comparable) {
			return replayRootConflict(root, "writer root overlaps the replay root", nil)
		}
	}
	return nil
}

// RejectProtectedRoots fails closed when any protected root (source, store, and
// similar) equals the replay root or is located inside it. A protected root
// containing the replay root is allowed only when the containment chain also
// covers the verified run root.
func (store *AgentRunnerReplayStore) RejectProtectedRoots(roots ...string) error {
	if store == nil {
		return fmt.Errorf("%w: nil replay store", ErrInvalidAgentRunnerReplayStore)
	}
	replayRoot, err := store.comparableReplayRoot()
	if err != nil {
		return err
	}
	for _, root := range roots {
		if err := rejectUnsafeReplayPathComponents(root); err != nil {
			return err
		}
		comparable, err := resolveComparableReplayPath(root)
		if err != nil {
			return replayRootConflict(root, "protected root cannot be resolved", err)
		}
		if comparable == replayRoot || sameReplayDirectory(comparable, replayRoot) {
			return replayRootConflict(root, "protected root equals the replay root", nil)
		}
		if replayPathContainsOrAliases(replayRoot, comparable) {
			return replayRootConflict(root, "protected root is located inside the replay root", nil)
		}
		if replayPathContainsOrAliases(comparable, replayRoot) {
			if store.runRoot == "" {
				return replayRootConflict(root, "protected root contains the replay root without a bound run root", nil)
			}
			runRoot, err := resolveComparableReplayPath(store.runRoot)
			if err != nil {
				return replayRootConflict(root, "bound run root cannot be resolved", err)
			}
			if !replayPathContainsOrAliases(comparable, runRoot) {
				return replayRootConflict(root, "protected root contains the replay root without covering the run root", nil)
			}
		}
	}
	return nil
}

// validateRecordRunBinding fails closed when a record bound to a foreign run
// would be written into a store constructed for another runID.
func (store *AgentRunnerReplayStore) validateRecordRunBinding(runID, requestID string) error {
	if store == nil || store.runID == "" || runID == store.runID {
		return nil
	}
	return fmt.Errorf("%w: requestId %s belongs to runId %s but the replay store is bound to runId %s", ErrAgentRunnerReplayRootConflict, requestID, runID, store.runID)
}

// validatePersistedRunID fails closed as corruption when a record already
// published on disk belongs to a different run than the bound runID.
func (store *AgentRunnerReplayStore) validatePersistedRunID(path, kind, runID string) error {
	if store == nil || store.runID == "" || runID == store.runID {
		return nil
	}
	return replayStoreCorruption(path, fmt.Sprintf("persisted %s is bound to runId %s but the replay store is bound to runId %s", kind, runID, store.runID), nil)
}

func (store *AgentRunnerReplayStore) comparableReplayRoot() (string, error) {
	resolved, err := resolveComparableReplayPath(store.root)
	if err != nil {
		return "", replayRootConflict(store.root, "replay root cannot be resolved", err)
	}
	return resolved, nil
}

// verifyPathSafetyLocked repeats the component-level symlink and reparse
// rejection at runtime so a directory replaced after construction cannot
// silently redirect reads or writes outside the replay root.
func (store *AgentRunnerReplayStore) verifyPathSafetyLocked() error {
	if err := rejectUnsafeReplayPathComponents(store.root); err != nil {
		return err
	}
	for _, name := range []string{"requests", "completions", agentRunnerReplayOwnershipFileName} {
		if err := rejectUnsafeReplayPathComponents(filepath.Join(store.root, name)); err != nil {
			return err
		}
	}
	return nil
}

func validateAgentRunnerReplayRunRoot(runRoot string) (string, error) {
	if strings.TrimSpace(runRoot) == "" {
		return "", fmt.Errorf("%w: empty run root", ErrInvalidAgentRunnerReplayStore)
	}
	cleanRoot := filepath.Clean(runRoot)
	if !filepath.IsAbs(cleanRoot) {
		return "", fmt.Errorf("%w: run root must be absolute", ErrInvalidAgentRunnerReplayStore)
	}
	info, err := os.Stat(cleanRoot)
	if err != nil {
		return "", fmt.Errorf("%w: run root unavailable: %v", ErrInvalidAgentRunnerReplayStore, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%w: run root is not a directory", ErrInvalidAgentRunnerReplayStore)
	}
	return cleanRoot, nil
}

// validateAgentRunnerRunRootMarker requires the durable run-root marker written
// by the run event store before any AgentRunner step can exist. Without it any
// directory could masquerade as a durable run root.
func validateAgentRunnerRunRootMarker(runRoot string) error {
	eventsDir := filepath.Join(runRoot, agentRunnerReplayRunEventDirName)
	if err := rejectUnsafeReplayPathComponents(eventsDir); err != nil {
		return err
	}
	info, err := os.Stat(eventsDir)
	if err != nil {
		return fmt.Errorf("%w: run root events directory unavailable: %v", ErrInvalidAgentRunnerReplayStore, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%w: run root events path is not a directory", ErrInvalidAgentRunnerReplayStore)
	}
	eventLog := filepath.Join(eventsDir, agentRunnerReplayRunEventFileName)
	if err := rejectUnsafeReplayPathComponents(eventLog); err != nil {
		return err
	}
	logInfo, err := os.Lstat(eventLog)
	if err != nil {
		return fmt.Errorf("%w: run root event log unavailable: %v", ErrInvalidAgentRunnerReplayStore, err)
	}
	if !logInfo.Mode().IsRegular() {
		return fmt.Errorf("%w: run root event log is not a regular file", ErrInvalidAgentRunnerReplayStore)
	}
	return nil
}

func validateAgentRunnerReplayRootInvariants(replayRoot, runRoot string) error {
	if replayRoot == runRoot {
		return fmt.Errorf("%w: replay root must not equal the run root", ErrInvalidAgentRunnerReplayStore)
	}
	if !replayPathContains(runRoot, replayRoot) {
		return fmt.Errorf("%w: replay root must live inside the run root", ErrInvalidAgentRunnerReplayStore)
	}
	events := filepath.Join(runRoot, agentRunnerReplayRunEventDirName)
	if replayPathContains(replayRoot, events) || replayPathContains(events, replayRoot) {
		return fmt.Errorf("%w: replay root must not overlap the run event store", ErrInvalidAgentRunnerReplayStore)
	}
	return nil
}

// ensureAgentRunnerReplayDirectory creates a replay-store directory that must
// live under an owned parent, then syncs the parent after creation or confirmation.
func ensureAgentRunnerReplayDirectory(directory, parent string) error {
	if err := os.Mkdir(directory, 0o755); err != nil && !errors.Is(err, fs.ErrExist) {
		return fmt.Errorf("%w: create %s: %v", ErrInvalidAgentRunnerReplayStore, directory, err)
	}
	info, err := os.Stat(directory)
	if err != nil {
		return fmt.Errorf("%w: confirm %s: %v", ErrInvalidAgentRunnerReplayStore, directory, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%w: %s is not a directory", ErrInvalidAgentRunnerReplayStore, directory)
	}
	if err := replayStoreSyncParentDirectory(parent); err != nil {
		return fmt.Errorf("%w: sync %s: %w", ErrInvalidAgentRunnerReplayStore, parent, err)
	}
	return nil
}

func ensureAgentRunnerReplayOwnership(replayRoot, runRoot, runID string) error {
	markerPath := filepath.Join(replayRoot, agentRunnerReplayOwnershipFileName)
	expectedHash := agentRunnerReplayRunRootHash(runRoot)
	record, found, err := replayStoreLoadOwnershipRecord(markerPath)
	if err != nil {
		return err
	}
	if found {
		return validateAgentRunnerReplayOwnership(record, markerPath, expectedHash, runRoot, runID)
	}
	if err := rejectPublishedReplayRecords(replayRoot); err != nil {
		return err
	}
	marker := AgentRunnerReplayOwnershipRecord{
		Kind:          agentRunnerReplayOwnershipKind,
		SchemaVersion: evidence.SchemaVersion,
		RunID:         runID,
		RunRootPath:   runRoot,
		RunRootHash:   expectedHash,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
	}
	collided := false
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		if !collided {
			writeErr := writeReplayRecordNoReplace(markerPath, marker)
			if writeErr == nil {
				return nil
			}
			if !errors.Is(writeErr, fs.ErrExist) {
				return fmt.Errorf("%w: publish ownership marker at %s: %w", ErrInvalidAgentRunnerReplayStore, markerPath, writeErr)
			}
			// A no-replace collision proves another writer owns the marker slot.
			// Only bounded rereads follow so a visibility gap cannot replace it.
			collided = true
		}
		existing, found, loadErr := replayStoreLoadOwnershipRecord(markerPath)
		if loadErr != nil {
			return loadErr
		}
		if found {
			return validateAgentRunnerReplayOwnership(existing, markerPath, expectedHash, runRoot, runID)
		}
		runtime.Gosched()
	}
	return fmt.Errorf("%w: replay ownership marker not visible after no-replace collision at %s", ErrAgentRunnerReplayStoreConvergence, markerPath)
}

func loadAgentRunnerReplayOwnership(path string) (AgentRunnerReplayOwnershipRecord, bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return AgentRunnerReplayOwnershipRecord{}, false, nil
	}
	if err != nil {
		return AgentRunnerReplayOwnershipRecord{}, false, replayStoreCorruption(path, "ownership record stat failed", err)
	}
	if !info.Mode().IsRegular() {
		return AgentRunnerReplayOwnershipRecord{}, false, replayStoreCorruption(path, "ownership record is not a regular file", nil)
	}
	wire, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return AgentRunnerReplayOwnershipRecord{}, false, nil
	}
	if err != nil {
		return AgentRunnerReplayOwnershipRecord{}, false, replayStoreCorruption(path, "ownership record read failed", err)
	}
	var record AgentRunnerReplayOwnershipRecord
	if err := evidence.DecodeStrictJSON(wire, &record); err != nil {
		return AgentRunnerReplayOwnershipRecord{}, false, replayStoreCorruption(path, "ownership record decode failed", err)
	}
	return record, true, nil
}

func validateAgentRunnerReplayOwnership(record AgentRunnerReplayOwnershipRecord, path, expectedHash, runRoot, runID string) error {
	if record.Kind != agentRunnerReplayOwnershipKind ||
		record.SchemaVersion != evidence.SchemaVersion ||
		!evidence.ValidID(record.RunID) ||
		!evidence.ValidHash(record.RunRootHash) ||
		strings.TrimSpace(record.RunRootPath) == "" {
		return replayStoreCorruption(path, "ownership record shape invalid", nil)
	}
	if _, err := time.Parse(time.RFC3339, record.CreatedAt); err != nil {
		return replayStoreCorruption(path, "ownership record createdAt is invalid", err)
	}
	if agentRunnerReplayRunRootHash(record.RunRootPath) != record.RunRootHash {
		return replayRootConflict(path, "ownership record run root hash is inconsistent", nil)
	}
	if record.RunID != runID {
		return replayRootConflict(path, fmt.Sprintf("ownership record is bound to runId %s but the store requested %s", record.RunID, runID), nil)
	}
	if record.RunRootHash != expectedHash && !sameReplayDirectory(record.RunRootPath, runRoot) {
		return replayRootConflict(path, "ownership record does not match the current run root", nil)
	}
	return nil
}

func agentRunnerReplayRunRootHash(runRoot string) string {
	return evidence.Digest(agentRunnerReplayOwnershipDomain, []byte(filepath.Clean(runRoot)))
}

func sameReplayDirectory(first, second string) bool {
	firstInfo, firstErr := os.Stat(first)
	secondInfo, secondErr := os.Stat(second)
	return firstErr == nil && secondErr == nil && os.SameFile(firstInfo, secondInfo)
}

func rejectPublishedReplayRecords(replayRoot string) error {
	for _, name := range []string{"requests", "completions"} {
		directory := filepath.Join(replayRoot, name)
		entries, err := os.ReadDir(directory)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return replayStoreCorruption(directory, "replay root layout is not readable", err)
		}
		if len(entries) > 0 {
			return replayStoreCorruption(directory, "published replay records exist without ownership metadata", nil)
		}
	}
	return nil
}

func replayRootConflict(path, detail string, cause error) error {
	if cause == nil {
		return fmt.Errorf("%w: %s at %s", ErrAgentRunnerReplayRootConflict, detail, path)
	}
	return fmt.Errorf("%w: %s at %s: %w", ErrAgentRunnerReplayRootConflict, detail, path, cause)
}

// rejectUnsafeReplayPathComponents fails closed when any already-existing
// component of path (up to the volume root) is a symlink or reparse point.
func rejectUnsafeReplayPathComponents(path string) error {
	clean := filepath.Clean(path)
	for component := clean; ; {
		unsafe, err := replayPathComponentUnsafe(component)
		if err != nil {
			return fmt.Errorf("%w: component check failed for %s: %w", ErrAgentRunnerReplayRootUnsafePath, component, err)
		}
		if unsafe {
			return fmt.Errorf("%w: component is a symlink or reparse point: %s", ErrAgentRunnerReplayRootUnsafePath, component)
		}
		parent := filepath.Dir(component)
		if parent == component {
			return nil
		}
		component = parent
	}
}

// resolveComparableReplayPath normalizes a path for containment comparison by
// resolving the deepest existing prefix through symlinks and re-joining the
// remaining segments. Errors are returned instead of being treated as disjoint.
func resolveComparableReplayPath(path string) (string, error) {
	clean := filepath.Clean(path)
	if !filepath.IsAbs(clean) {
		return "", fmt.Errorf("%w: path must be absolute: %s", ErrInvalidAgentRunnerReplayStore, path)
	}
	current := clean
	suffix := make([]string, 0, 4)
	for {
		resolved, err := filepath.EvalSymlinks(current)
		if err == nil {
			parts := append([]string{resolved}, reverseReplaySegments(suffix)...)
			return filepath.Join(parts...), nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", err
		}
		suffix = append(suffix, filepath.Base(current))
		current = parent
	}
}

func reverseReplaySegments(segments []string) []string {
	reversed := make([]string, len(segments))
	for index, segment := range segments {
		reversed[len(segments)-1-index] = segment
	}
	return reversed
}

// replayPathContains reports whether parent is equal to or an ancestor of child.
func replayPathContains(parent, child string) bool {
	parent = filepath.Clean(parent)
	child = filepath.Clean(child)
	if runtime.GOOS == "windows" {
		parent = strings.ToLower(parent)
		child = strings.ToLower(child)
	}
	if parent == child {
		return true
	}
	separator := string(os.PathSeparator)
	if !strings.HasSuffix(parent, separator) {
		parent += separator
	}
	return strings.HasPrefix(child, parent)
}

// replayPathContainsOrAliases extends replayPathContains with file-identity
// containment so directory aliases that do not share an identical spelling
// (for example Windows 8.3 short names) are still recognized. The identity walk
// starts at an existing parent and looks for the same directory among the
// already-existing ancestors of child.
func replayPathContainsOrAliases(parent, child string) bool {
	if replayPathContains(parent, child) {
		return true
	}
	return replayPathContainedByIdentity(parent, child)
}

func replayPathContainedByIdentity(parent, child string) bool {
	parentInfo, err := os.Stat(parent)
	if err != nil {
		return false
	}
	for current := filepath.Clean(child); ; {
		info, statErr := os.Stat(current)
		if statErr == nil && os.SameFile(parentInfo, info) {
			return true
		}
		parentDir := filepath.Dir(current)
		if parentDir == current {
			return false
		}
		current = parentDir
	}
}
