//go:build windows

package adapters

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// This file is the Windows C6 prebuilt-slot implementation (CONTRACTS.md §7
// 「C6 判定路径」, ADR-013). A slot is a fixed container:
//
//	offset 0..8   magic "PRSLOT1\0"
//	offset 8..16  payloadLen, uint64 little endian
//	offset 16..   payload: one canonical JSONL record, byte for byte the record
//	              the existing no-replace publication writes
//
// The directory entry is created once, before the first dispatch, and is never
// rebuilt, replaced or appended to afterwards. Publication writes the payload,
// flushes it, then commits the slot header (payloadLen 0 -> len) and flushes the
// header, so visibility is bounded by the header commit and not by the directory
// entry: "payload written, header not committed" is frozen as absent for a
// reader.
//
// Nothing here decides when a slot is used: this file only provides the
// primitives. The shared R/C publication wrapper and the R/C readers reach them
// through the hooks below, and the c6slot activation file is the only place that
// flips replayStoreSlotEnabled, so every other build keeps the ordinary record
// path.
//
// The slot creation chain lives in
// agent_runner_replay_store_c6_windows_create.go (windows && c6slot) so the
// default Windows build carries no unreferenced creation code.

const (
	// replaySlotHeaderSize is the fixed header size: 8 byte magic + uint64le.
	replaySlotHeaderSize = 16

	// Publication stages of one slot, in execution order. They are the crash
	// points a durability experiment stops at, and they are reported through the
	// existing replayStorePublishStage seam so a rig sees one token stream per
	// mechanism step.
	replaySlotStageCreated       = "slot-created"
	replaySlotStageWritten       = "slot-written"
	replaySlotStageDataSynced    = "slot-data-synced"
	replaySlotStageHeaderWritten = "slot-header-written"
	replaySlotStageHeaderSynced  = "slot-header-synced"
	replaySlotStageUnlocked      = "slot-unlocked"

	replaySlotGenericRead         = 0x80000000
	replaySlotGenericWrite        = 0x40000000
	replaySlotShareRead           = 0x00000001
	replaySlotShareWrite          = 0x00000002
	replaySlotOpenExisting        = 3
	replaySlotFileAttributeNormal = 0x00000080
	replaySlotLockExclusive       = 0x00000002
	replaySlotLockFailImmediately = 0x00000001
)

const (
	// replaySlotLockWaitLimit bounds the wait for a slot lock. A lock held past
	// the limit is reported as a convergence failure instead of blocking the run
	// forever: the caller's bounded-retry contract, not an unbounded kernel wait.
	replaySlotLockWaitLimit = 5 * time.Second
	// replaySlotLockPollInterval is the retry cadence of the non-blocking
	// LockFileEx probe inside that limit.
	replaySlotLockPollInterval = time.Millisecond
)

var (
	replaySlotMagic = [8]byte{'P', 'R', 'S', 'L', 'O', 'T', '1', 0x00}

	replaySlotErrorLockViolation = syscall.Errno(33)

	replaySlotKernel32         = syscall.NewLazyDLL("kernel32.dll")
	replaySlotProcLockFileEx   = replaySlotKernel32.NewProc("LockFileEx")
	replaySlotProcUnlockFileEx = replaySlotKernel32.NewProc("UnlockFileEx")
)

// The slot primitives are named from here rather than from the activation file
// so that every Windows build references them and no symbol depends on the
// c6slot build for being used. The build capability bit is replayStoreSlotEnabled,
// declared in the shared layer (agent_runner_replay_store.go) because the shared
// R/C wrapper and readers consult it, and only the c6slot activation file sets it
// to true. Carrying the primitives is not enough on its own: the store also has to
// have completed its own prebuild (store.slotsActive) before the slot path is
// taken for it (CONTRACTS.md §7 「A2 激活语义」).
func init() {
	// The primitives are named here, in a file every Windows build compiles, so
	// the shared R/C wrapper and readers can reach them without a build-tag fork.
	// They stay inert until the c6slot activation file flips
	// replayStoreSlotEnabled and a store completes its prebuild.
	replayStoreSlotPublishHook = writeReplaySlotNoReplace
	replayStoreSlotReadRequestHook = loadRequestRecordSlot
	replayStoreSlotReadCompletionHook = loadCompletionRecordSlot
}

// Private seams of the C6 slot pipeline. They are nil in production, and while
// nil each step runs the mechanism the platform proves with. A seam replaces one
// mechanism token at a time so a falsification rig can build a degraded arm
// (skip a flush, skip the header commit) without editing the product path.
//
// Seam contract: they fire while the caller holds the slot lock and, for store
// methods, the store lock. They may not carry decisions and may not call back
// into the store.
//
// slotCreateFn lives in agent_runner_replay_store_c6_windows_create.go (windows && c6slot).
var (
	slotFlushFn       func(file *os.File) error
	slotHeaderWriteFn func(file *os.File, payloadLen uint64) error
)

// replaySlotState is the classification a reader derives from the slot header,
// in the frozen order: magic first, then the length class.
type replaySlotState int

const (
	// replaySlotStateEmpty covers both the 16-byte prebuilt slot and the
	// uncommitted intermediate state (payload written, header not committed).
	// Both read as absent; neither is corruption.
	replaySlotStateEmpty replaySlotState = iota
	// replaySlotStatePublished means the committed payloadLen is non-zero and
	// the total length matches it exactly.
	replaySlotStatePublished
)

// writeReplaySlotNoReplace publishes one payload into a prebuilt slot. The slot
// must already exist: this path never creates the directory entry.
//
// The six stages are reported in order - slot-created (the empty slot is
// confirmed under the lock), slot-written, slot-data-synced, slot-header-written,
// slot-header-synced and the unlock. A failure before the header commit leaves
// the slot absent for a reader; a failure after it leaves the slot published
// with the payload it carried, which a reopen converges on as published.
func writeReplaySlotNoReplace(path string, payload []byte) error {
	if len(payload) == 0 {
		// An empty payload would commit payloadLen = 0, which is exactly the
		// "not committed" state: refuse it instead of publishing a slot that can
		// never be read.
		return fmt.Errorf("%w: replay slot payload is empty", ErrInvalidAgentRunnerReplayStore)
	}
	if int64(len(payload)) > replayStoreRecordBound {
		return fmt.Errorf("%w: replay slot payload exceeds the %d byte bound", ErrInvalidAgentRunnerReplayStore, replayStoreRecordBound)
	}
	file, err := openReplaySlot(path, replaySlotGenericRead|replaySlotGenericWrite)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	if err := lockReplaySlot(file); err != nil {
		return err
	}
	unlocked := false
	defer func() {
		if !unlocked {
			unlockReplaySlot(file)
		}
	}()
	state, err := classifyReplaySlotFile(file, path)
	if err != nil {
		return err
	}
	if state == replaySlotStatePublished {
		// A committed slot is never rewritten or truncated: the no-replace
		// contract hands the collision to the caller's existing convergence path.
		return fs.ErrExist
	}
	replayStorePublishStage(replaySlotStageCreated, path)
	// Reclaim an uncommitted intermediate state inside the lock: truncate to the
	// header so a longer failed payload cannot leave a tail behind a shorter
	// retry, then write the payload at its fixed offset.
	if err := file.Truncate(replaySlotHeaderSize); err != nil {
		return fmt.Errorf("replay slot %s: reclaim: %w", path, err)
	}
	if _, err := file.WriteAt(payload, replaySlotHeaderSize); err != nil {
		return fmt.Errorf("replay slot %s: payload write: %w", path, err)
	}
	replayStorePublishStage(replaySlotStageWritten, path)
	if err := flushSlotFile(file); err != nil {
		// No header on a failed data flush: the slot stays absent.
		return fmt.Errorf("replay slot %s: payload flush: %w", path, err)
	}
	replayStorePublishStage(replaySlotStageDataSynced, path)
	if err := writeSlotHeader(file, uint64(len(payload))); err != nil {
		return fmt.Errorf("replay slot %s: header commit: %w", path, err)
	}
	replayStorePublishStage(replaySlotStageHeaderWritten, path)
	if err := flushSlotFile(file); err != nil {
		// The header is already written, so the slot may be visible as published
		// from here on; a reopen converges on it instead of repairing it.
		return fmt.Errorf("replay slot %s: header flush: %w", path, err)
	}
	replayStorePublishStage(replaySlotStageHeaderSynced, path)
	unlockReplaySlot(file)
	unlocked = true
	replayStorePublishStage(replaySlotStageUnlocked, path)
	return nil
}

// readReplaySlot reads one slot. Absence and the uncommitted intermediate state
// both answer (nil, false, nil); an abnormal length, a bad magic or a payload
// that is not strict JSON is corruption.
func readReplaySlot(path string) ([]byte, bool, error) {
	payload, state, err := inspectReplaySlot(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if state != replaySlotStatePublished {
		return nil, false, nil
	}
	return payload, true, nil
}

// loadRequestRecordSlot is the slot form of the request reader: the header is
// classified first (bad magic, an abnormal length or a payload that is not strict
// JSON is corruption), the uncommitted intermediate state reads as absent, and
// only a committed payload is decoded and checked against requestID - the same
// order and the same corruption classification the plain record reader uses.
func loadRequestRecordSlot(path, requestID string) (AgentRunnerRequestRecord, bool, error) {
	payload, found, err := readReplaySlot(path)
	if err != nil || !found {
		return AgentRunnerRequestRecord{}, false, err
	}
	record, err := DecodeAgentRunnerRequestRecord(payload)
	if err != nil {
		return AgentRunnerRequestRecord{}, false, replayStoreCorruption(path, "request record decode failed", err)
	}
	if record.Request.RequestID != requestID {
		return AgentRunnerRequestRecord{}, false, replayStoreCorruption(path, fmt.Sprintf("requestId mismatch on disk (want %s, got %s)", requestID, record.Request.RequestID), nil)
	}
	return record, true, nil
}

// loadCompletionRecordSlot is the completion form of loadRequestRecordSlot:
// identical classification order, identical corruption wording.
func loadCompletionRecordSlot(path, requestID string) (AgentRunnerCompletionRecord, bool, error) {
	payload, found, err := readReplaySlot(path)
	if err != nil || !found {
		return AgentRunnerCompletionRecord{}, false, err
	}
	record, err := DecodeAgentRunnerCompletionRecord(payload)
	if err != nil {
		return AgentRunnerCompletionRecord{}, false, replayStoreCorruption(path, "completion record decode failed", err)
	}
	if record.Completion.RequestID != requestID {
		return AgentRunnerCompletionRecord{}, false, replayStoreCorruption(path, fmt.Sprintf("completion requestId mismatch on disk (want %s, got %s)", requestID, record.Completion.RequestID), nil)
	}
	return record, true, nil
}

// inspectReplaySlot is the readReplaySlot form that keeps the missing-file case
// distinguishable from the empty-slot case: a caller that must tell "never
// prebuilt" from "prebuilt but not committed" needs that difference, and the
// classification below is shared by both forms so there is one validation order.
func inspectReplaySlot(path string) ([]byte, replaySlotState, error) {
	// The frozen arbitration protocol gives reader and writer the same handle
	// configuration (GENERIC_READ|GENERIC_WRITE, share read+write); only the
	// writer takes the slot lock, but both ask for the same access so neither can
	// be starved by a holder that shares less than the protocol assumes.
	file, err := openReplaySlot(path, replaySlotGenericRead|replaySlotGenericWrite)
	if err != nil {
		return nil, replaySlotStateEmpty, err
	}
	defer func() { _ = file.Close() }()
	state, err := classifyReplaySlotFile(file, path)
	if err != nil {
		return nil, replaySlotStateEmpty, err
	}
	if state != replaySlotStatePublished {
		return nil, state, nil
	}
	info, err := file.Stat()
	if err != nil {
		return nil, replaySlotStateEmpty, replayStoreCorruption(path, "slot stat failed", err)
	}
	payloadLen := info.Size() - replaySlotHeaderSize
	payload := make([]byte, payloadLen)
	if _, err := file.ReadAt(payload, replaySlotHeaderSize); err != nil {
		return nil, replaySlotStateEmpty, replayStoreCorruption(path, "slot payload read failed", err)
	}
	// Strict decode at the container boundary: the payload must be exactly one
	// strict JSON value (canonicalization refuses duplicate keys, invalid UTF-8
	// and trailing values). The record-level decode and the requestId check stay
	// with the caller that knows which record type the slot holds.
	if _, err := evidence.Canonicalize(payload); err != nil {
		return nil, replaySlotStateEmpty, replayStoreCorruption(path, "slot payload is not strict JSON", err)
	}
	return payload, state, nil
}

// classifyReplaySlotFile applies the frozen validation order to one open slot:
// magic first, then the length class, and the payload length bound before any
// length arithmetic.
func classifyReplaySlotFile(file *os.File, path string) (replaySlotState, error) {
	if file == nil {
		return replaySlotStateEmpty, replayStoreCorruption(path, "slot file is nil", nil)
	}
	info, err := file.Stat()
	if err != nil {
		return replaySlotStateEmpty, replayStoreCorruption(path, "slot stat failed", err)
	}
	total := info.Size()
	if total < replaySlotHeaderSize {
		return replaySlotStateEmpty, replayStoreCorruption(path, fmt.Sprintf("slot is shorter than its %d-byte header (%d bytes)", replaySlotHeaderSize, total), nil)
	}
	header := make([]byte, replaySlotHeaderSize)
	if _, err := file.ReadAt(header, 0); err != nil {
		return replaySlotStateEmpty, replayStoreCorruption(path, "slot header read failed", err)
	}
	if !bytes.Equal(header[:len(replaySlotMagic)], replaySlotMagic[:]) {
		return replaySlotStateEmpty, replayStoreCorruption(path, "slot magic mismatch", nil)
	}
	payloadLen := binary.LittleEndian.Uint64(header[len(replaySlotMagic):])
	switch {
	case payloadLen == 0:
		// 16 bytes is the prebuilt empty slot; more than 16 is a payload written
		// without a header commit. Both are absent, and neither is corruption.
		return replaySlotStateEmpty, nil
	case payloadLen > uint64(replayStoreRecordBound):
		return replaySlotStateEmpty, replayStoreCorruption(path, fmt.Sprintf("slot payload length %d exceeds the %d byte bound", payloadLen, replayStoreRecordBound), nil)
	case total != int64(replaySlotHeaderSize)+int64(payloadLen):
		return replaySlotStateEmpty, replayStoreCorruption(path, fmt.Sprintf("slot length %d does not match its committed payload length %d", total, payloadLen), nil)
	default:
		return replaySlotStatePublished, nil
	}
}

// openReplaySlot opens a slot with the access and sharing the arbitration
// protocol fixes: readers take read access, writers read and write, and both
// share read and write so the lock, not the share mode, arbitrates.
func openReplaySlot(path string, access uint32) (*os.File, error) {
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	handle, err := syscall.CreateFile(
		pointer,
		access,
		replaySlotShareRead|replaySlotShareWrite,
		nil,
		replaySlotOpenExisting,
		replaySlotFileAttributeNormal,
		0,
	)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(handle), path), nil
}

// lockReplaySlot takes the exclusive slot lock whose range covers the whole
// file, waiting at most replaySlotLockWaitLimit. The wait is bounded because the
// protocol requires an uncertain outcome instead of an unbounded block.
func lockReplaySlot(file *os.File) error {
	deadline := time.Now().Add(replaySlotLockWaitLimit)
	for {
		locked, err := tryLockReplaySlot(file)
		if err != nil {
			return err
		}
		if locked {
			return nil
		}
		if !time.Now().Before(deadline) {
			return fmt.Errorf("%w: replay slot lock was not acquired within %s", ErrAgentRunnerReplayStoreConvergence, replaySlotLockWaitLimit)
		}
		time.Sleep(replaySlotLockPollInterval)
	}
}

// tryLockReplaySlot is the non-blocking LockFileEx probe. A held lock reports
// (false, nil): that is contention, not a failure.
func tryLockReplaySlot(file *os.File) (bool, error) {
	if file == nil {
		return false, os.ErrInvalid
	}
	overlapped := new(syscall.Overlapped)
	locked, _, callErr := replaySlotProcLockFileEx.Call(
		file.Fd(),
		replaySlotLockExclusive|replaySlotLockFailImmediately,
		0,
		0xffffffff,
		0xffffffff,
		uintptr(unsafe.Pointer(overlapped)),
	)
	if locked == 0 {
		if errors.Is(callErr, replaySlotErrorLockViolation) {
			return false, nil
		}
		return false, callErr
	}
	return true, nil
}

// unlockReplaySlot releases the same exact range that was locked. Closing the
// handle releases it as well, so this is the explicit form the protocol names.
func unlockReplaySlot(file *os.File) {
	if file == nil {
		return
	}
	overlapped := new(syscall.Overlapped)
	_, _, _ = replaySlotProcUnlockFileEx.Call(file.Fd(), 0, 0xffffffff, 0xffffffff, uintptr(unsafe.Pointer(overlapped)))
}

// replaySlotHeader builds the fixed header: magic followed by the little-endian
// payload length.
func replaySlotHeader(payloadLen uint64) [replaySlotHeaderSize]byte {
	var header [replaySlotHeaderSize]byte
	copy(header[:], replaySlotMagic[:])
	binary.LittleEndian.PutUint64(header[len(replaySlotMagic):], payloadLen)
	return header
}

func flushSlotFile(file *os.File) error {
	if slotFlushFn != nil {
		return slotFlushFn(file)
	}
	return syscall.FlushFileBuffers(syscall.Handle(file.Fd()))
}

func writeSlotHeader(file *os.File, payloadLen uint64) error {
	if slotHeaderWriteFn != nil {
		return slotHeaderWriteFn(file, payloadLen)
	}
	header := replaySlotHeader(payloadLen)
	_, err := file.WriteAt(header[len(replaySlotMagic):], int64(len(replaySlotMagic)))
	return err
}
