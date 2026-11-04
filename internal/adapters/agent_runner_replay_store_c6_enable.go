//go:build windows && c6slot

package adapters

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// This file activates the C6 prebuilt-slot implementation for the c6slot build
// (CONTRACTS.md §7 「C6 判定路径」, ADR-013). It is the only place that assigns
// replayStoreSlotEnabled and replayStorePrebuildSlotsFn, so every other build
// keeps the shared no-op and the Unix path untouched.
//
// The prebuild runs before the first dispatch with the full declared requestId
// set and is the only step that creates a slot's directory entry. It runs in two
// orders, and the order is the point:
//
//   - first prebuild (no prebuild-set.jsonl): verify existing slots, create the
//     missing ones, and publish the set record only after every slot is in
//     place. A failure therefore leaves no record, so a reopen reads the run as
//     "never prebuilt" and converges by idempotent verification/creation, while
//     the failed run itself is already fail-closed.
//   - reopen (prebuild-set.jsonl present): verify the declared set matches the
//     record item for item, then verify every declared slot. A declared slot
//     that is missing is corruption, never a slot to recreate - that is the only
//     way "missing slot" can be told apart from "never prebuilt".
//
// A lost prebuild-set therefore loses plan-drift detection and nothing else: R/C
// binding, orphan detection and runId validation still run, which is the
// degradation the contract accepts and registers.

const (
	agentRunnerReplayPrebuildSetFileName      = "prebuild-set.jsonl"
	agentRunnerReplayPrebuildSetKind          = "agent-runner-replay-prebuild-set"
	agentRunnerReplayPrebuildSetSchemaVersion = "1"
	agentRunnerReplayPrebuildSetDomain        = "proofrail:agent-runner-replay-prebuild-set:1\n"
)

// agentRunnerReplayPrebuildSetBody is the hashed part of the store-local record:
// everything except the record hash itself.
type agentRunnerReplayPrebuildSetBody struct {
	Kind          string   `json:"kind"`
	SchemaVersion string   `json:"schemaVersion"`
	RunID         string   `json:"runId"`
	CreatedAt     string   `json:"createdAt"`
	Count         int      `json:"count"`
	RequestIDs    []string `json:"requestIds"`
}

// agentRunnerReplayPrebuildSetRecord is one canonical JSONL store-local record.
// It never enters the wire and never participates in an R/C recordHash; its own
// hash is checked by validateReplayPrebuildSet.
type agentRunnerReplayPrebuildSetRecord struct {
	Kind          string   `json:"kind"`
	SchemaVersion string   `json:"schemaVersion"`
	RunID         string   `json:"runId"`
	CreatedAt     string   `json:"createdAt"`
	Count         int      `json:"count"`
	RequestIDs    []string `json:"requestIds"`
	RecordHash    string   `json:"recordHash"`
}

func init() {
	replayStoreSlotEnabled = true
	replayStorePrebuildSlotsFn = prebuildReplaySlots
}

// prebuildReplaySlots is the C6 prebuild implementation behind the shared seam.
// A failure is fail-closed for the whole run: partially created empty slots may
// stay (a reopen verifies them), but nothing dispatches.
//
// Success is also the only thing that activates the slot path for this store
// (CONTRACTS.md §7 「A2 激活语义」): slotsActive is set after the first prebuild
// published the set record, and after a reopen verified the declared set and
// every declared slot. Both branches therefore leave the store on the slot path
// only when every slot is in place, and a failure leaves it on the ordinary
// record path.
func prebuildReplaySlots(store *AgentRunnerReplayStore, requestIDs []string) error {
	if err := store.validateConstructed(); err != nil {
		return err
	}
	declared, err := replayPrebuildSetIDs(requestIDs)
	if err != nil {
		return err
	}
	// Lock order is store.mu -> slot lock, so the store lock is taken here and
	// the slot locks are taken inside. No slot lock is ever taken first.
	store.mu.Lock()
	defer store.mu.Unlock()
	if err := store.verifyPathSafetyLocked(); err != nil {
		return err
	}
	recordPath := filepath.Join(store.root, agentRunnerReplayPrebuildSetFileName)
	record, found, err := loadReplayPrebuildSet(recordPath)
	if err != nil {
		return err
	}
	if found {
		if err := verifyReplayPrebuildSetSlots(store, record, recordPath, declared); err != nil {
			return err
		}
		store.slotsActive = true
		return nil
	}
	if err := firstReplayPrebuild(store, recordPath, declared); err != nil {
		return err
	}
	store.slotsActive = true
	return nil
}

// replayPrebuildSetIDs canonicalizes the declared set: dictionary order,
// distinct, every item a valid id. The record field requires that order, so the
// caller's order is normalized rather than trusted.
func replayPrebuildSetIDs(requestIDs []string) ([]string, error) {
	declared := append([]string(nil), requestIDs...)
	slices.Sort(declared)
	for index, requestID := range declared {
		if !evidence.ValidID(requestID) {
			return nil, fmt.Errorf("%w: prebuild requestId is invalid", ErrInvalidAgentRunnerReplayStore)
		}
		if index > 0 && declared[index-1] == requestID {
			return nil, fmt.Errorf("%w: prebuild requestIds must be distinct", ErrInvalidAgentRunnerReplayStore)
		}
	}
	return declared, nil
}

// replayPrebuildSlotPaths returns the two slots a requestId owns. Only R/C
// records use the slot container; every other store-local record keeps its
// original publication and read path.
func replayPrebuildSlotPaths(store *AgentRunnerReplayStore, requestID string) ([]string, error) {
	request, err := store.requestPath(requestID)
	if err != nil {
		return nil, err
	}
	completion, err := store.completionPath(requestID)
	if err != nil {
		return nil, err
	}
	return []string{request, completion}, nil
}

// firstReplayPrebuild runs the prebuild when no set record exists yet: verify
// what exists, create what is missing, and publish the record last.
func firstReplayPrebuild(store *AgentRunnerReplayStore, recordPath string, declared []string) error {
	for _, requestID := range declared {
		paths, err := replayPrebuildSlotPaths(store, requestID)
		if err != nil {
			return err
		}
		for _, path := range paths {
			if err := ensureReplaySlot(path); err != nil {
				return err
			}
		}
	}
	return publishReplayPrebuildSet(recordPath, store.runID, declared)
}

// ensureReplaySlot verifies an existing slot and creates a missing one. It never
// rebuilds or replaces a slot that already exists; an uncommitted intermediate
// state is left to the writer, which reclaims it under the slot lock.
func ensureReplaySlot(path string) error {
	if _, _, err := inspectReplaySlot(path); err == nil {
		return nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := createReplaySlotWindows(path); err != nil {
		if errors.Is(err, fs.ErrExist) {
			// A concurrent prebuild won the name: verify the slot it created
			// instead of replacing it.
			return verifyReplaySlotExists(path)
		}
		return err
	}
	return nil
}

// verifyReplayPrebuildSetSlots is the reopen branch: the record decides which
// slots must exist. A declared slot that is missing is corruption, because the
// record is the only proof that it was ever allocated.
func verifyReplayPrebuildSetSlots(store *AgentRunnerReplayStore, record agentRunnerReplayPrebuildSetRecord, recordPath string, declared []string) error {
	if err := validateReplayPrebuildSet(record, recordPath, store.runID, declared); err != nil {
		return err
	}
	for _, requestID := range declared {
		paths, err := replayPrebuildSlotPaths(store, requestID)
		if err != nil {
			return err
		}
		for _, path := range paths {
			if err := verifyReplaySlotExists(path); err != nil {
				return err
			}
		}
	}
	return nil
}

// verifyReplaySlotExists verifies one existing slot. Absence is corruption here:
// the caller only reaches this after the set record declared the slot.
func verifyReplaySlotExists(path string) error {
	if _, _, err := inspectReplaySlot(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return replayStoreCorruption(path, "declared prebuilt slot is missing", nil)
		}
		return err
	}
	return nil
}

// loadReplayPrebuildSet reads the store-local set record. Absence is normal (a
// first prebuild) and is reported as (zero, false, nil).
func loadReplayPrebuildSet(path string) (agentRunnerReplayPrebuildSetRecord, bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return agentRunnerReplayPrebuildSetRecord{}, false, nil
	}
	if err != nil {
		return agentRunnerReplayPrebuildSetRecord{}, false, replayStoreCorruption(path, "prebuild set stat failed", err)
	}
	if !info.Mode().IsRegular() {
		return agentRunnerReplayPrebuildSetRecord{}, false, replayStoreCorruption(path, "prebuild set is not a regular file", nil)
	}
	wire, err := readAgentRunnerArtifactBytes(path, replayStoreRecordBound)
	if errors.Is(err, fs.ErrNotExist) {
		return agentRunnerReplayPrebuildSetRecord{}, false, nil
	}
	if err != nil {
		return agentRunnerReplayPrebuildSetRecord{}, false, replayStoreCorruption(path, "prebuild set read failed", err)
	}
	var record agentRunnerReplayPrebuildSetRecord
	if err := evidence.DecodeStrictJSON(wire, &record); err != nil {
		return agentRunnerReplayPrebuildSetRecord{}, false, replayStoreCorruption(path, "prebuild set decode failed", err)
	}
	return record, true, nil
}

// validateReplayPrebuildSet checks the record shape, its own hash, the run
// binding, and that the declared set matches it item for item. A mismatch is
// plan drift and therefore corruption.
func validateReplayPrebuildSet(record agentRunnerReplayPrebuildSetRecord, path, runID string, declared []string) error {
	if record.Kind != agentRunnerReplayPrebuildSetKind {
		return replayStoreCorruption(path, "prebuild set kind is not "+agentRunnerReplayPrebuildSetKind, nil)
	}
	if record.SchemaVersion != agentRunnerReplayPrebuildSetSchemaVersion {
		return replayStoreCorruption(path, "prebuild set schemaVersion is unknown", nil)
	}
	if _, err := time.Parse(time.RFC3339, record.CreatedAt); err != nil {
		return replayStoreCorruption(path, "prebuild set createdAt is invalid", err)
	}
	if record.Count != len(record.RequestIDs) {
		return replayStoreCorruption(path, "prebuild set count does not match requestIds", nil)
	}
	if !slices.IsSorted(record.RequestIDs) {
		return replayStoreCorruption(path, "prebuild set requestIds are not in dictionary order", nil)
	}
	for index, requestID := range record.RequestIDs {
		if !evidence.ValidID(requestID) || index > 0 && record.RequestIDs[index-1] == requestID {
			return replayStoreCorruption(path, "prebuild set requestIds are invalid or duplicated", nil)
		}
	}
	expected, err := digestMessage(agentRunnerReplayPrebuildSetDomain, record.body())
	if err != nil {
		return err
	}
	if expected != record.RecordHash {
		return replayStoreCorruption(path, "prebuild set recordHash mismatch", nil)
	}
	if record.RunID != runID {
		return replayRootConflict(path, fmt.Sprintf("prebuild set is bound to runId %s but the store requested %s", record.RunID, runID), nil)
	}
	if !slices.Equal(record.RequestIDs, declared) {
		return replayStoreCorruption(path, "prebuild set does not match the declared requestId set (plan drift)", nil)
	}
	return nil
}

// publishReplayPrebuildSet publishes the set record with the ordinary no-replace
// store-local path: not through a slot, no slot header, no wire, no R/C hash. A
// collision converges on the winner by bounded rereads, exactly as the ownership
// record does, and does not require its own proven durability.
func publishReplayPrebuildSet(path, runID string, declared []string) error {
	body := agentRunnerReplayPrebuildSetBody{
		Kind:          agentRunnerReplayPrebuildSetKind,
		SchemaVersion: agentRunnerReplayPrebuildSetSchemaVersion,
		RunID:         runID,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
		Count:         len(declared),
		RequestIDs:    append([]string{}, declared...),
	}
	hash, err := digestMessage(agentRunnerReplayPrebuildSetDomain, body)
	if err != nil {
		return err
	}
	record := agentRunnerReplayPrebuildSetRecord{
		Kind:          body.Kind,
		SchemaVersion: body.SchemaVersion,
		RunID:         body.RunID,
		CreatedAt:     body.CreatedAt,
		Count:         body.Count,
		RequestIDs:    body.RequestIDs,
		RecordHash:    hash,
	}
	collided := false
	for attempt := 0; attempt < replayStoreReadCollisionRetries; attempt++ {
		if !collided {
			writeErr := writeReplayRecordNoReplace(path, record)
			if writeErr == nil {
				return nil
			}
			if !errors.Is(writeErr, fs.ErrExist) {
				return fmt.Errorf("%w: publish prebuild set at %s: %w", ErrInvalidAgentRunnerReplayStore, path, writeErr)
			}
			// A no-replace collision proves another writer owns the record slot.
			// Only bounded rereads follow, so a visibility gap cannot replace it.
			collided = true
		}
		existing, found, loadErr := loadReplayPrebuildSet(path)
		if loadErr != nil {
			return loadErr
		}
		if found {
			return validateReplayPrebuildSet(existing, path, runID, declared)
		}
		runtime.Gosched()
	}
	return fmt.Errorf("%w: prebuild set not visible after no-replace collision at %s", ErrAgentRunnerReplayStoreConvergence, path)
}

func (record agentRunnerReplayPrebuildSetRecord) body() agentRunnerReplayPrebuildSetBody {
	return agentRunnerReplayPrebuildSetBody{
		Kind:          record.Kind,
		SchemaVersion: record.SchemaVersion,
		RunID:         record.RunID,
		CreatedAt:     record.CreatedAt,
		Count:         record.Count,
		RequestIDs:    record.RequestIDs,
	}
}
