//go:build windows && c6slot

package adapters

// Tests for the Windows C6 prebuilt-slot implementation (CONTRACTS.md §7
// 「C6 判定路径」, ADR-013). The build constraint is the isolation mechanism
// itself: this file is compiled only for windows && c6slot, so `GOOS=linux`
// builds and every unactivated Windows build never see it, and the c6slot
// activation file is the only place that flips the shared flag.
//
// Every test here is written to be falsifiable against the pre-C6 baseline: the
// primitives, the frozen validation order, the prebuild lifecycle, the allowed
// set table and the two-arm BA3 comparison all fail to compile or fail to pass
// before the C6 implementation exists.

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// c6Store builds a store whose publication durability is forced to proven, so
// the R/C write path is reachable. The gate tests below build a store through
// the platform constructor instead and keep the real (unproven) claim.
func c6Store(t *testing.T) *AgentRunnerReplayStore {
	t.Helper()
	return replayStoreMustNew(t, t.TempDir())
}

// c6UnprovenStore builds a store with the durability claim the platform
// constructor gives it, which on Windows is unproven.
func c6UnprovenStore(t *testing.T) *AgentRunnerReplayStore {
	t.Helper()
	store, err := newAgentRunnerReplayStoreAt(t.TempDir())
	if err != nil {
		t.Fatalf("construct replay store: %v", err)
	}
	if got := store.PublicationDurability(); got != PublishDurabilityUnproven {
		t.Fatalf("platform durability = %s, want unproven on this build", got)
	}
	return store
}

// c6PrebuiltStore builds a store and runs the prebuild for the declared set, so
// a test can start from the state the assembly produces before the first
// dispatch.
func c6PrebuiltStore(t *testing.T, requestIDs ...string) *AgentRunnerReplayStore {
	t.Helper()
	store := c6Store(t)
	if err := store.PrebuildReplaySlots(requestIDs); err != nil {
		t.Fatalf("prebuild %v: %v", requestIDs, err)
	}
	return store
}

// c6SlotPaths returns the two slots one requestId owns.
func c6SlotPaths(t *testing.T, store *AgentRunnerReplayStore, requestID string) (string, string) {
	t.Helper()
	request, err := store.requestPath(requestID)
	if err != nil {
		t.Fatalf("request path for %s: %v", requestID, err)
	}
	completion, err := store.completionPath(requestID)
	if err != nil {
		t.Fatalf("completion path for %s: %v", requestID, err)
	}
	return request, completion
}

// c6Payload encodes one record exactly the way the store publishes it: canonical
// JSON plus the JSONL newline.
func c6Payload(t *testing.T, record any) []byte {
	t.Helper()
	canonical, err := evidence.EncodeCanonical(record)
	if err != nil {
		t.Fatalf("encode canonical record: %v", err)
	}
	return append(canonical, '\n')
}

// c6RawSlotBytes builds raw slot bytes with an arbitrary magic, so a test can
// write the exact on-disk shape it wants to classify.
func c6RawSlotBytes(magic []byte, payloadLen uint64, payload []byte) []byte {
	raw := make([]byte, replaySlotHeaderSize)
	copy(raw, magic)
	binary.LittleEndian.PutUint64(raw[len(replaySlotMagic):], payloadLen)
	return append(raw, payload...)
}

// c6WriteFile writes bytes at path, creating the parent directory.
func c6WriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create parent of %s: %v", path, err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

// c6ReadFile reads one file and fails the test instead of returning an error.
func c6ReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

// c6PlainRecord returns the bytes of an ordinary record file, or nil when the
// file does not exist. That nil is the same "absent" the slot reader reports for
// a prebuilt slot that never committed a payload.
func c6PlainRecord(t *testing.T, path string) []byte {
	t.Helper()
	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatalf("read plain record %s: %v", path, err)
	}
	return raw
}

// c6SlotPayload returns the payload behind one slot header, or nil when the slot
// file does not exist or carries no committed payload. The 16 byte header is
// deliberately not returned: BA3 compares record content, and the header belongs
// to the activation mechanism, not to the record.
func c6SlotPayload(t *testing.T, path string) []byte {
	t.Helper()
	raw, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		t.Fatalf("read slot %s: %v", path, err)
	}
	if len(raw) < replaySlotHeaderSize {
		t.Fatalf("slot %s is %d bytes, shorter than its %d byte header: the arm did not take the slot path", path, len(raw), replaySlotHeaderSize)
	}
	if !bytes.Equal(raw[:len(replaySlotMagic)], replaySlotMagic[:]) {
		t.Fatalf("slot %s magic = %q, want %q: the arm did not take the slot path", path, raw[:len(replaySlotMagic)], replaySlotMagic[:])
	}
	payloadLen := binary.LittleEndian.Uint64(raw[len(replaySlotMagic):replaySlotHeaderSize])
	if payloadLen == 0 {
		return nil
	}
	if want := replaySlotHeaderSize + int(payloadLen); len(raw) != want {
		t.Fatalf("slot %s total length = %d, want %d", path, len(raw), want)
	}
	return raw[replaySlotHeaderSize:]
}

// c6RequireArmPayloadsAgree is the BA3 record comparison: the control arm's
// ordinary JSONL bytes and the call arm's slot payload must be identical, and
// absence must agree with absence.
func c6RequireArmPayloadsAgree(t *testing.T, label string, control, call []byte) {
	t.Helper()
	if control == nil && call == nil {
		return
	}
	if control == nil || call == nil {
		t.Fatalf("%s: one arm recorded a payload and the other did not (control=%d bytes, call=%d bytes)", label, len(control), len(call))
	}
	if !bytes.Equal(control, call) {
		t.Fatalf("%s: record bytes differ between the arms:\ncontrol=%q\ncall=%q", label, control, call)
	}
}

// c6StageStop models a process that stops existing at one exact slot stage.
type c6StageStop struct{ stage string }

// c6RecordRequestAndStopAtStage drives one request publication and stops the
// process at the given slot stage through the publish stage seam, because a
// timing-based kill cannot hit a chosen stage deterministically. The panic
// unwinds the deferred unlocks, exactly like a crash releasing kernel locks.
func c6RecordRequestAndStopAtStage(t *testing.T, store *AgentRunnerReplayStore, record AgentRunnerRequestRecord, stage string) {
	t.Helper()
	original := replayStorePublishStageHook
	replayStorePublishStageHook = func(current string, _ string) {
		if current == stage {
			panic(c6StageStop{stage: current})
		}
	}
	defer func() { replayStorePublishStageHook = original }()
	defer func() {
		recovered := recover()
		if _, ok := recovered.(c6StageStop); !ok {
			t.Fatalf("request publication did not stop at stage %s: recovered %v", stage, recovered)
		}
	}()
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatalf("request publication before stage %s failed: %v", stage, err)
	}
	t.Fatalf("request publication returned even though it must have stopped at stage %s", stage)
}

// ---------------------------------------------------------------------------
// 1. slot creation, publication and reading
// ---------------------------------------------------------------------------

// TestC6SlotCreateWritesSixteenByteEmptyHeader pins the empty-slot shape: magic
// followed by payloadLen 0, no payload, no temp residue.
func TestC6SlotCreateWritesSixteenByteEmptyHeader(t *testing.T) {
	store := c6Store(t)
	path, _ := c6SlotPaths(t, store, "request-one")
	if err := createReplaySlotWindows(path); err != nil {
		t.Fatalf("create slot: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat slot: %v", err)
	}
	if info.Size() != replaySlotHeaderSize {
		t.Fatalf("empty slot size = %d, want %d", info.Size(), replaySlotHeaderSize)
	}
	raw := c6ReadFile(t, path)
	if !bytes.Equal(raw[:len(replaySlotMagic)], replaySlotMagic[:]) {
		t.Fatalf("empty slot magic = %q, want %q", raw[:len(replaySlotMagic)], replaySlotMagic[:])
	}
	if got := binary.LittleEndian.Uint64(raw[len(replaySlotMagic):replaySlotHeaderSize]); got != 0 {
		t.Fatalf("empty slot payloadLen = %d, want 0", got)
	}

	payload, found, err := readReplaySlot(path)
	if err != nil {
		t.Fatalf("read empty slot: %v", err)
	}
	if found || payload != nil {
		t.Fatalf("an empty slot must read as absent, got found=%t payload=%q", found, payload)
	}

	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatalf("read slot directory: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != filepath.Base(path) {
		t.Fatalf("slot creation left residue in %s: %v", filepath.Dir(path), entries)
	}

	// The directory entry is created once: a second create must collide, not rebuild.
	before := c6ReadFile(t, path)
	if err := createReplaySlotWindows(path); !errors.Is(err, fs.ErrExist) {
		t.Fatalf("second create error = %v, want fs.ErrExist", err)
	}
	if after := c6ReadFile(t, path); !bytes.Equal(before, after) {
		t.Fatal("a refused second create must leave the slot bytes unchanged")
	}
}

// TestC6SlotPublishTotalLengthIsHeaderPlusPayload pins the published shape and
// the reader round trip through the record reader.
func TestC6SlotPublishTotalLengthIsHeaderPlusPayload(t *testing.T) {
	store := c6Store(t)
	path, _ := c6SlotPaths(t, store, "request-one")
	if err := createReplaySlotWindows(path); err != nil {
		t.Fatalf("create slot: %v", err)
	}
	record := replayStoreRequestRecord(t, "request-one")
	payload := c6Payload(t, record)
	if err := writeReplaySlotNoReplace(path, payload); err != nil {
		t.Fatalf("publish slot: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat slot: %v", err)
	}
	if want := int64(replaySlotHeaderSize + len(payload)); info.Size() != want {
		t.Fatalf("published slot size = %d, want %d", info.Size(), want)
	}
	raw := c6ReadFile(t, path)
	if got := binary.LittleEndian.Uint64(raw[len(replaySlotMagic):replaySlotHeaderSize]); got != uint64(len(payload)) {
		t.Fatalf("published payloadLen = %d, want %d", got, len(payload))
	}
	if !bytes.Equal(raw[replaySlotHeaderSize:], payload) {
		t.Fatal("published payload bytes differ from the bytes the caller handed over")
	}

	stored, found, err := readReplaySlot(path)
	if err != nil || !found {
		t.Fatalf("read published slot: found=%t err=%v", found, err)
	}
	if !bytes.Equal(stored, payload) {
		t.Fatal("slot reader returned different payload bytes")
	}
	decoded, found, err := loadRequestRecordSlot(path, "request-one")
	if err != nil || !found {
		t.Fatalf("request reader: found=%t err=%v", found, err)
	}
	if decoded.RecordHash != record.RecordHash {
		t.Fatalf("request reader hash = %s, want %s", decoded.RecordHash, record.RecordHash)
	}
}

// TestC6SlotCommittedSlotRejectsSecondWriteWithoutChangingBytes pins the
// no-replace contract inside the slot: a committed slot is never rewritten.
func TestC6SlotCommittedSlotRejectsSecondWriteWithoutChangingBytes(t *testing.T) {
	store := c6Store(t)
	path, _ := c6SlotPaths(t, store, "request-one")
	if err := createReplaySlotWindows(path); err != nil {
		t.Fatalf("create slot: %v", err)
	}
	first := c6Payload(t, replayStoreRequestRecord(t, "request-one"))
	if err := writeReplaySlotNoReplace(path, first); err != nil {
		t.Fatalf("first publication: %v", err)
	}
	before := c6ReadFile(t, path)

	second := c6Payload(t, replayStoreRequestRecord(t, "request-two"))
	if err := writeReplaySlotNoReplace(path, second); !errors.Is(err, fs.ErrExist) {
		t.Fatalf("second publication error = %v, want fs.ErrExist", err)
	}
	if after := c6ReadFile(t, path); !bytes.Equal(before, after) {
		t.Fatal("a refused second publication must leave the committed slot unchanged")
	}
	stored, found, err := readReplaySlot(path)
	if err != nil || !found || !bytes.Equal(stored, first) {
		t.Fatalf("committed slot must keep its first payload: found=%t err=%v", found, err)
	}
}

// TestC6SlotRejectsEmptyPayload pins that payloadLen 0 can never be published as
// a committed payload.
func TestC6SlotRejectsEmptyPayload(t *testing.T) {
	store := c6Store(t)
	path, _ := c6SlotPaths(t, store, "request-one")
	if err := createReplaySlotWindows(path); err != nil {
		t.Fatalf("create slot: %v", err)
	}
	before := c6ReadFile(t, path)
	for _, payload := range [][]byte{nil, {}} {
		if err := writeReplaySlotNoReplace(path, payload); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
			t.Fatalf("empty payload error = %v, want ErrInvalidAgentRunnerReplayStore", err)
		}
	}
	if after := c6ReadFile(t, path); !bytes.Equal(before, after) {
		t.Fatal("a refused empty payload must leave the slot unchanged")
	}
}

// ---------------------------------------------------------------------------
// 2. frozen validation order (A2 ① -> ② -> ③)
// ---------------------------------------------------------------------------

// TestC6SlotValidationOrderBadMagicIsCorruption pins ①: magic is checked before
// any length arithmetic, so a bad magic reports magic even when the length is
// also wrong.
func TestC6SlotValidationOrderBadMagicIsCorruption(t *testing.T) {
	path := filepath.Join(t.TempDir(), "slot.request-one.jsonl")
	c6WriteFile(t, path, c6RawSlotBytes([]byte("BROKEN!!"), 5, []byte("abc")))
	_, _, err := readReplaySlot(path)
	if !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("bad magic error = %v, want corruption", err)
	}
	if !strings.Contains(err.Error(), "slot magic mismatch") {
		t.Fatalf("bad magic must be reported as a magic mismatch before the length class: %v", err)
	}
}

// TestC6SlotValidationOrderShortFileIsCorruption pins the length floor: fewer
// than 16 bytes cannot even carry a header.
func TestC6SlotValidationOrderShortFileIsCorruption(t *testing.T) {
	path := filepath.Join(t.TempDir(), "slot.request-one.jsonl")
	c6WriteFile(t, path, []byte("PRSLOT1\x00"))
	_, _, err := readReplaySlot(path)
	if !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("short slot error = %v, want corruption", err)
	}
	if !strings.Contains(err.Error(), "shorter than its") {
		t.Fatalf("short slot must be reported as a header length violation: %v", err)
	}

	// Absence stays distinguishable from every abnormal shape.
	if _, _, err := inspectReplaySlot(filepath.Join(t.TempDir(), "missing")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("missing slot error = %v, want fs.ErrNotExist", err)
	}
}

// TestC6SlotValidationOrderLengthMismatchIsCorruption pins "committed length
// must equal total minus header".
func TestC6SlotValidationOrderLengthMismatchIsCorruption(t *testing.T) {
	path := filepath.Join(t.TempDir(), "slot.request-one.jsonl")
	c6WriteFile(t, path, c6RawSlotBytes(replaySlotMagic[:], 8, []byte("abc")))
	_, _, err := readReplaySlot(path)
	if !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("length mismatch error = %v, want corruption", err)
	}
	if !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("length mismatch must be reported as a length violation: %v", err)
	}
}

// TestC6SlotValidationOrderUncommittedPayloadIsAbsent is the freeze clause:
// payloadLen 0 with a longer file is the uncommitted intermediate state, and it
// reads as absent - never corruption.
func TestC6SlotValidationOrderUncommittedPayloadIsAbsent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "slot.request-one.jsonl")
	c6WriteFile(t, path, c6RawSlotBytes(replaySlotMagic[:], 0, []byte("0123456789")))

	payload, found, err := readReplaySlot(path)
	if err != nil {
		t.Fatalf("uncommitted intermediate state must not be corruption: %v", err)
	}
	if found || payload != nil {
		t.Fatalf("uncommitted intermediate state must read as absent, got found=%t", found)
	}
	_, state, err := inspectReplaySlot(path)
	if err != nil {
		t.Fatalf("inspect uncommitted state: %v", err)
	}
	if state != replaySlotStateEmpty {
		t.Fatalf("uncommitted state = %v, want empty", state)
	}
	if _, found, err := loadRequestRecordSlot(path, "request-one"); err != nil || found {
		t.Fatalf("record reader on uncommitted state: found=%t err=%v", found, err)
	}
}

// TestC6SlotValidationOrderUndecodablePayloadIsCorruption pins ③ for both
// shapes: bytes that are not strict JSON, and strict JSON that is not the
// record the slot is supposed to hold.
func TestC6SlotValidationOrderUndecodablePayloadIsCorruption(t *testing.T) {
	notJSON := []byte("this is not json")
	path := filepath.Join(t.TempDir(), "slot.request-one.jsonl")
	c6WriteFile(t, path, c6RawSlotBytes(replaySlotMagic[:], uint64(len(notJSON)), notJSON))
	if _, _, err := readReplaySlot(path); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("non-JSON payload error = %v, want corruption", err)
	}

	wrongShape := []byte("{\"nope\":1}\n")
	other := filepath.Join(t.TempDir(), "slot.request-two.jsonl")
	c6WriteFile(t, other, c6RawSlotBytes(replaySlotMagic[:], uint64(len(wrongShape)), wrongShape))
	_, found, err := loadRequestRecordSlot(other, "request-two")
	if found || !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("wrong-shape payload: found=%t err=%v, want corruption", found, err)
	}
	if !strings.Contains(err.Error(), "decode failed") {
		t.Fatalf("wrong-shape payload must be reported as a decode failure: %v", err)
	}
}

// TestC6SlotValidationOrderRequestIDMismatchIsCorruption pins the last check: a
// committed, decodable record that belongs to another requestId is corruption.
func TestC6SlotValidationOrderRequestIDMismatchIsCorruption(t *testing.T) {
	store := c6Store(t)
	path, _ := c6SlotPaths(t, store, "request-two")
	if err := createReplaySlotWindows(path); err != nil {
		t.Fatalf("create slot: %v", err)
	}
	if err := writeReplaySlotNoReplace(path, c6Payload(t, replayStoreRequestRecord(t, "request-one"))); err != nil {
		t.Fatalf("publish slot: %v", err)
	}

	// The same slot read under its own requestId is a valid record: the corruption
	// below is the mismatch, not an undecodable payload.
	if _, found, err := loadRequestRecordSlot(path, "request-one"); err != nil || !found {
		t.Fatalf("control read: found=%t err=%v", found, err)
	}
	_, found, err := loadRequestRecordSlot(path, "request-two")
	if found || !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("requestId mismatch: found=%t err=%v, want corruption", found, err)
	}
	if !strings.Contains(err.Error(), "requestId mismatch") {
		t.Fatalf("requestId mismatch must be reported as such: %v", err)
	}
}

// ---------------------------------------------------------------------------
// 3. uncommitted intermediate state and failed retries
// ---------------------------------------------------------------------------

// TestC6SlotFreshPayloadFlushFailureLeavesHeaderUncommitted pins the freeze rule
// at the pipeline level: a failed data flush writes no header, so a reader sees
// absent even though the payload bytes exist.
func TestC6SlotFreshPayloadFlushFailureLeavesHeaderUncommitted(t *testing.T) {
	store := c6Store(t)
	path, _ := c6SlotPaths(t, store, "request-one")
	if err := createReplaySlotWindows(path); err != nil {
		t.Fatalf("create slot: %v", err)
	}
	payload := c6Payload(t, replayStoreRequestRecord(t, "request-one"))

	original := slotFlushFn
	failures := 1
	slotFlushFn = func(*os.File) error {
		if failures > 0 {
			failures--
			return errors.New("injected data flush failure")
		}
		return nil
	}
	t.Cleanup(func() { slotFlushFn = original })

	err := writeReplaySlotNoReplace(path, payload)
	if err == nil || !strings.Contains(err.Error(), "payload flush") {
		t.Fatalf("publication with a failed data flush = %v, want a payload flush error", err)
	}
	raw := c6ReadFile(t, path)
	if want := replaySlotHeaderSize + len(payload); len(raw) != want {
		t.Fatalf("payload bytes present = %d, want %d", len(raw), want)
	}
	if got := binary.LittleEndian.Uint64(raw[len(replaySlotMagic):replaySlotHeaderSize]); got != 0 {
		t.Fatalf("header payloadLen = %d after a failed data flush, want 0", got)
	}
	if _, found, err := readReplaySlot(path); err != nil || found {
		t.Fatalf("a failed data flush must leave the slot absent: found=%t err=%v", found, err)
	}
}

// TestC6SlotFailedLongWriteThenShortRetryLeavesNoTail pins the reclaim rule: the
// retry truncates to the header inside the slot lock, so a longer failed payload
// cannot leave a tail behind a shorter committed one.
func TestC6SlotFailedLongWriteThenShortRetryLeavesNoTail(t *testing.T) {
	store := c6Store(t)
	path, _ := c6SlotPaths(t, store, "request-one")
	if err := createReplaySlotWindows(path); err != nil {
		t.Fatalf("create slot: %v", err)
	}
	long := []byte("{\"pad\":\"" + strings.Repeat("p", 512) + "\"}\n")
	short := c6Payload(t, replayStoreRequestRecord(t, "request-one"))

	original := slotFlushFn
	failures := 1
	slotFlushFn = func(*os.File) error {
		if failures > 0 {
			failures--
			return errors.New("injected data flush failure")
		}
		return nil
	}
	if err := writeReplaySlotNoReplace(path, long); err == nil {
		t.Fatal("the long publication must fail on the injected data flush")
	}
	slotFlushFn = original
	if raw := c6ReadFile(t, path); len(raw) != replaySlotHeaderSize+len(long) {
		t.Fatalf("intermediate state size = %d, want %d", len(raw), replaySlotHeaderSize+len(long))
	}

	if err := writeReplaySlotNoReplace(path, short); err != nil {
		t.Fatalf("short retry: %v", err)
	}
	raw := c6ReadFile(t, path)
	if want := replaySlotHeaderSize + len(short); len(raw) != want {
		t.Fatalf("retried slot size = %d, want %d (the long tail must be truncated away)", len(raw), want)
	}
	if !bytes.Equal(raw[replaySlotHeaderSize:], short) {
		t.Fatal("retried slot payload is not exactly the short payload")
	}
	stored, found, err := readReplaySlot(path)
	if err != nil || !found || !bytes.Equal(stored, short) {
		t.Fatalf("retried slot read: found=%t err=%v", found, err)
	}
}

// ---------------------------------------------------------------------------
// 4. prebuild (PrebuildReplaySlots)
// ---------------------------------------------------------------------------

// TestC6PrebuildCreatesFullSetAndIsIdempotent pins the full-set prebuild and the
// reopen no-op: a second call verifies, never rewrites.
func TestC6PrebuildCreatesFullSetAndIsIdempotent(t *testing.T) {
	store := c6Store(t)
	ids := []string{"request-one", "request-two", "request-three"}
	if err := store.PrebuildReplaySlots(ids); err != nil {
		t.Fatalf("first prebuild: %v", err)
	}
	if !store.slotsActive {
		t.Fatal("a completed prebuild must activate the slot path for this store")
	}
	recordPath := filepath.Join(store.root, agentRunnerReplayPrebuildSetFileName)
	slots := make(map[string][]byte, len(ids)*2)
	for _, requestID := range ids {
		requestPath, completionPath := c6SlotPaths(t, store, requestID)
		for _, path := range []string{requestPath, completionPath} {
			raw := c6ReadFile(t, path)
			if len(raw) != replaySlotHeaderSize {
				t.Fatalf("prebuilt slot %s size = %d, want %d", path, len(raw), replaySlotHeaderSize)
			}
			if got := binary.LittleEndian.Uint64(raw[len(replaySlotMagic):replaySlotHeaderSize]); got != 0 {
				t.Fatalf("prebuilt slot %s payloadLen = %d, want 0", path, got)
			}
			slots[path] = raw
		}
	}
	firstRecord := c6ReadFile(t, recordPath)

	if err := store.PrebuildReplaySlots(ids); err != nil {
		t.Fatalf("second prebuild: %v", err)
	}
	if !store.slotsActive {
		t.Fatal("the reopen branch must keep the store on the slot path")
	}
	if after := c6ReadFile(t, recordPath); !bytes.Equal(firstRecord, after) {
		t.Fatal("the reopen branch must verify the prebuild set record, not rewrite it")
	}
	for path, before := range slots {
		if after := c6ReadFile(t, path); !bytes.Equal(before, after) {
			t.Fatalf("the reopen branch must not touch the prebuilt slot %s", path)
		}
	}

	// The activation bit is not only a field: what this store publishes lands in
	// the slot container, header included, instead of a plain JSONL file.
	record := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(record); err != nil {
		t.Fatalf("record request on an activated store: %v", err)
	}
	activatedPath, _ := c6SlotPaths(t, store, "request-one")
	raw := c6ReadFile(t, activatedPath)
	if len(raw) < replaySlotHeaderSize {
		t.Fatalf("activated store wrote %d bytes, want a slot carrying its %d byte header", len(raw), replaySlotHeaderSize)
	}
	if !bytes.Equal(raw[:len(replaySlotMagic)], replaySlotMagic[:]) {
		t.Fatalf("activated store wrote a plain record with magic %q, want the slot magic %q", raw[:len(replaySlotMagic)], replaySlotMagic[:])
	}
	if got, want := binary.LittleEndian.Uint64(raw[len(replaySlotMagic):replaySlotHeaderSize]), uint64(len(c6Payload(t, record))); got != want {
		t.Fatalf("activated store slot payloadLen = %d, want %d", got, want)
	}
}

// TestC6PrebuildFailsClosedWhenAnySlotCannotBeCreated pins the fail-closed
// semantics: one failed slot fails the whole prebuild and publishes no record, so
// a reopen reads the run as never prebuilt.
func TestC6PrebuildFailsClosedWhenAnySlotCannotBeCreated(t *testing.T) {
	store := c6Store(t)
	ids := []string{"request-one", "request-two", "request-three"}

	original := slotCreateFn
	slotCreateFn = func(directory, pattern string) (*os.File, error) {
		if strings.Contains(pattern, "completion.request-two") {
			return nil, errors.New("injected slot creation failure")
		}
		return os.CreateTemp(directory, pattern)
	}
	t.Cleanup(func() { slotCreateFn = original })

	if err := store.PrebuildReplaySlots(ids); err == nil {
		t.Fatal("a prebuild with one failing slot must fail the whole run")
	}
	if store.slotsActive {
		t.Fatal("a failed prebuild must leave the store on the ordinary record path")
	}
	recordPath := filepath.Join(store.root, agentRunnerReplayPrebuildSetFileName)
	if _, err := os.Stat(recordPath); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a failed prebuild must publish no set record, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(store.root, "requests")); err != nil {
		t.Fatalf("partial slots may stay on disk for a reopen to verify, stat err = %v", err)
	}

	// An unactivated store keeps the ordinary record reader, so a leftover empty
	// slot on that path is corruption and never a silent misread: that is exactly
	// why a failed prebuild may leave its partial slots behind.
	requestTwo := replayStoreRequestRecord(t, "request-two")
	requestTwoPath, _ := c6SlotPaths(t, store, "request-two")
	if decision, err := store.RecordRequest(requestTwo); decision != "" || !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("ordinary read against a leftover slot = %q err=%v, want corruption", decision, err)
	}
	if raw := c6ReadFile(t, requestTwoPath); len(raw) != replaySlotHeaderSize {
		t.Fatalf("the leftover slot must stay untouched, size = %d, want %d", len(raw), replaySlotHeaderSize)
	}

	// Recovery converges because no record exists: the run is treated as never
	// prebuilt and the prebuild runs again.
	slotCreateFn = original
	if err := store.PrebuildReplaySlots(ids); err != nil {
		t.Fatalf("prebuild after the failure: %v", err)
	}
	if !store.slotsActive {
		t.Fatal("the prebuild that completed after the failure must activate the store")
	}
	if _, err := os.Stat(recordPath); err != nil {
		t.Fatalf("the successful prebuild must publish the set record: %v", err)
	}
}

// TestC6PrebuildSetRecordFieldsAreExact pins the store-local record format field
// by field, including its canonical JSONL shape and its own recordHash.
func TestC6PrebuildSetRecordFieldsAreExact(t *testing.T) {
	store := c6Store(t)
	store.runID = "run-one"
	input := []string{"request-two", "request-one", "request-three"}
	if err := store.PrebuildReplaySlots(input); err != nil {
		t.Fatalf("prebuild: %v", err)
	}
	raw := c6ReadFile(t, filepath.Join(store.root, agentRunnerReplayPrebuildSetFileName))

	if !bytes.HasSuffix(raw, []byte("\n")) {
		t.Fatal("the prebuild set record must be JSONL terminated by one newline")
	}
	trimmed := raw[:len(raw)-1]
	if bytes.Contains(trimmed, []byte("\n")) {
		t.Fatal("the prebuild set record must be exactly one JSONL line")
	}
	canonical, err := evidence.Canonicalize(trimmed)
	if err != nil {
		t.Fatalf("prebuild set record is not strict JSON: %v", err)
	}
	if !bytes.Equal(canonical, trimmed) {
		t.Fatal("the prebuild set record must be canonical JSON")
	}

	var generic map[string]any
	if err := json.Unmarshal(trimmed, &generic); err != nil {
		t.Fatalf("decode prebuild set record: %v", err)
	}
	keys := make([]string, 0, len(generic))
	for key := range generic {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	wantKeys := []string{"count", "createdAt", "kind", "recordHash", "requestIds", "runId", "schemaVersion"}
	if strings.Join(keys, ",") != strings.Join(wantKeys, ",") {
		t.Fatalf("prebuild set record fields = %v, want %v", keys, wantKeys)
	}

	var record agentRunnerReplayPrebuildSetRecord
	if err := evidence.DecodeStrictJSON(trimmed, &record); err != nil {
		t.Fatalf("strict decode prebuild set record: %v", err)
	}
	if record.Kind != agentRunnerReplayPrebuildSetKind {
		t.Fatalf("kind = %q, want %q", record.Kind, agentRunnerReplayPrebuildSetKind)
	}
	if record.SchemaVersion != agentRunnerReplayPrebuildSetSchemaVersion {
		t.Fatalf("schemaVersion = %q, want %q", record.SchemaVersion, agentRunnerReplayPrebuildSetSchemaVersion)
	}
	if record.RunID != "run-one" {
		t.Fatalf("runId = %q, want run-one", record.RunID)
	}
	if record.Count != len(input) {
		t.Fatalf("count = %d, want %d", record.Count, len(input))
	}
	// The declared set is canonicalized on the way in: dictionary order, distinct.
	wantIDs := []string{"request-one", "request-three", "request-two"}
	if strings.Join(record.RequestIDs, ",") != strings.Join(wantIDs, ",") {
		t.Fatalf("requestIds = %v, want dictionary order %v", record.RequestIDs, wantIDs)
	}
	if _, err := time.Parse(time.RFC3339, record.CreatedAt); err != nil {
		t.Fatalf("createdAt = %q is not RFC3339: %v", record.CreatedAt, err)
	}
	if !evidence.ValidHash(record.RecordHash) {
		t.Fatalf("recordHash = %q is not a valid hash", record.RecordHash)
	}
	expectedHash, err := digestMessage(agentRunnerReplayPrebuildSetDomain, record.body())
	if err != nil {
		t.Fatalf("recompute recordHash: %v", err)
	}
	if expectedHash != record.RecordHash {
		t.Fatalf("recordHash = %s, want the hash of its own body %s", record.RecordHash, expectedHash)
	}
	drifted := record.body()
	drifted.RequestIDs = []string{"request-one"}
	otherHash, err := digestMessage(agentRunnerReplayPrebuildSetDomain, drifted)
	if err != nil {
		t.Fatalf("hash drifted body: %v", err)
	}
	if otherHash == record.RecordHash {
		t.Fatal("recordHash must bind the requestId set")
	}
}

// TestC6PrebuildDeclaredSlotMissingIsCorruption pins "record present + slot
// missing = corruption": the record is the only proof a slot was allocated, so a
// missing one can never be silently recreated.
func TestC6PrebuildDeclaredSlotMissingIsCorruption(t *testing.T) {
	for _, missingRequestSlot := range []bool{true, false} {
		name := "completion"
		if missingRequestSlot {
			name = "request"
		}
		t.Run(name, func(t *testing.T) {
			store := c6PrebuiltStore(t, "request-one", "request-two")
			requestPath, completionPath := c6SlotPaths(t, store, "request-two")
			target := requestPath
			if !missingRequestSlot {
				target = completionPath
			}
			if err := os.Remove(target); err != nil {
				t.Fatalf("remove %s: %v", target, err)
			}
			err := store.PrebuildReplaySlots([]string{"request-one", "request-two"})
			if !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
				t.Fatalf("missing declared slot %s error = %v, want corruption", filepath.Base(target), err)
			}
			if !strings.Contains(err.Error(), "declared prebuilt slot is missing") {
				t.Fatalf("missing declared slot must be reported as such: %v", err)
			}
		})
	}
}

// TestC6PrebuildMissingRecordAndMissingSlotsCreatesThem pins "record missing +
// slot missing = normal": the first prebuild branch verifies and creates.
func TestC6PrebuildMissingRecordAndMissingSlotsCreatesThem(t *testing.T) {
	store := c6PrebuiltStore(t, "request-one")
	if err := os.Remove(filepath.Join(store.root, agentRunnerReplayPrebuildSetFileName)); err != nil {
		t.Fatalf("remove prebuild set record: %v", err)
	}
	for _, directory := range []string{"requests", "completions"} {
		if err := os.RemoveAll(filepath.Join(store.root, directory)); err != nil {
			t.Fatalf("remove %s: %v", directory, err)
		}
	}

	if err := store.PrebuildReplaySlots([]string{"request-one"}); err != nil {
		t.Fatalf("prebuild from an empty root: %v", err)
	}
	requestPath, completionPath := c6SlotPaths(t, store, "request-one")
	for _, path := range []string{requestPath, completionPath} {
		raw := c6ReadFile(t, path)
		if len(raw) != replaySlotHeaderSize {
			t.Fatalf("recreated slot %s size = %d, want %d", path, len(raw), replaySlotHeaderSize)
		}
	}
	if _, err := os.Stat(filepath.Join(store.root, agentRunnerReplayPrebuildSetFileName)); err != nil {
		t.Fatalf("the first prebuild must publish the set record: %v", err)
	}
}

// TestC6PrebuildLostSetRecordDegradesDriftDetectionOnly pins the accepted
// degradation: losing prebuild-set.jsonl loses plan-drift detection, while the
// R/C binding checks still fail closed.
func TestC6PrebuildLostSetRecordDegradesDriftDetectionOnly(t *testing.T) {
	store := c6PrebuiltStore(t, "request-one")
	if err := os.Remove(filepath.Join(store.root, agentRunnerReplayPrebuildSetFileName)); err != nil {
		t.Fatalf("remove prebuild set record: %v", err)
	}

	// Drift is no longer detectable: a wider declared set is accepted.
	if err := store.PrebuildReplaySlots([]string{"request-one", "request-two"}); err != nil {
		t.Fatalf("a lost set record must degrade to a first prebuild, got %v", err)
	}

	// R/C binding still fails closed: a completion without its request is an orphan.
	_, completionPath := c6SlotPaths(t, store, "request-two")
	if err := writeReplaySlotNoReplace(completionPath, c6Payload(t, replayStoreCompletionRecord(t, replayStoreRequestRecord(t, "request-two"), "completion-two"))); err != nil {
		t.Fatalf("seed orphan completion: %v", err)
	}
	if _, err := store.State("request-two"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("orphan completion error = %v, want corruption", err)
	}

	// The requestId check on the slot reader still runs as well: a committed record
	// that belongs to another requestId is corruption, not a silent match.
	requestPath, _ := c6SlotPaths(t, store, "request-one")
	if err := writeReplaySlotNoReplace(requestPath, c6Payload(t, replayStoreRequestRecord(t, "request-two"))); err != nil {
		t.Fatalf("seed mismatched request slot: %v", err)
	}
	if _, found, err := loadRequestRecordSlot(requestPath, "request-one"); found || !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("slot requestId check after degradation: found=%t err=%v", found, err)
	}
}

// TestC6PrebuildPlanDriftIsCorruption pins the reopen plan check: a store that
// reopens the same root with a different declared requestId set is refused as
// corruption (plan drift), and the drifted set's slots are never created.
func TestC6PrebuildPlanDriftIsCorruption(t *testing.T) {
	store := c6Store(t)
	store.runID = "run-one"
	if err := store.PrebuildReplaySlots([]string{"request-one", "request-two"}); err != nil {
		t.Fatalf("first prebuild: %v", err)
	}

	reopened, err := newAgentRunnerReplayStoreAt(store.root)
	if err != nil {
		t.Fatalf("reopen the same root: %v", err)
	}
	reopened.runID = "run-one"
	err = reopened.PrebuildReplaySlots([]string{"request-one", "request-three"})
	if !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("drifted declared set error = %v, want corruption", err)
	}
	if !strings.Contains(err.Error(), "plan drift") {
		t.Fatalf("a drifted declared set must be reported as plan drift: %v", err)
	}
	if reopened.slotsActive {
		t.Fatal("a refused reopen must leave the store on the ordinary record path")
	}
	requestPath, completionPath := c6SlotPaths(t, store, "request-three")
	for _, path := range []string{requestPath, completionPath} {
		if _, statErr := os.Stat(path); !errors.Is(statErr, fs.ErrNotExist) {
			t.Fatalf("plan drift must not create the drifted slot %s: %v", path, statErr)
		}
	}
}

// TestC6PrebuildRunIDMismatchIsRootConflict pins the run binding of the set
// record: a reopen with a different runId is a root conflict, never the
// corruption class, and the store stays unactivated.
func TestC6PrebuildRunIDMismatchIsRootConflict(t *testing.T) {
	store := c6Store(t)
	store.runID = "run-one"
	if err := store.PrebuildReplaySlots([]string{"request-one"}); err != nil {
		t.Fatalf("first prebuild: %v", err)
	}

	reopened, err := newAgentRunnerReplayStoreAt(store.root)
	if err != nil {
		t.Fatalf("reopen the same root: %v", err)
	}
	reopened.runID = "run-two"
	err = reopened.PrebuildReplaySlots([]string{"request-one"})
	if !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("runId mismatch error = %v, want root conflict", err)
	}
	if errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("a runId mismatch must be a root conflict, never corruption: %v", err)
	}
	if !strings.Contains(err.Error(), "prebuild set is bound to runId run-one but the store requested run-two") {
		t.Fatalf("a runId mismatch must name both run ids: %v", err)
	}
	if reopened.slotsActive {
		t.Fatal("a refused reopen must leave the store on the ordinary record path")
	}
}

// ---------------------------------------------------------------------------
// 5. allowed set table (附件 C)
// ---------------------------------------------------------------------------

// c6AllowedSetCase is one row of the allowed set table: the cut point and the
// states the table permits in each scenario.
type c6AllowedSetCase struct {
	stage     string
	kill      []AgentRunnerReplayState
	powerLoss []AgentRunnerReplayState
}

func c6AllowedSetCases() []c6AllowedSetCase {
	return []c6AllowedSetCase{
		{
			stage:     replaySlotStageCreated,
			kill:      []AgentRunnerReplayState{AgentRunnerReplayStateAbsent},
			powerLoss: []AgentRunnerReplayState{AgentRunnerReplayStateAbsent},
		},
		{
			stage:     replaySlotStageWritten,
			kill:      []AgentRunnerReplayState{AgentRunnerReplayStateAbsent},
			powerLoss: []AgentRunnerReplayState{AgentRunnerReplayStateAbsent},
		},
		{
			stage:     replaySlotStageDataSynced,
			kill:      []AgentRunnerReplayState{AgentRunnerReplayStateAbsent},
			powerLoss: []AgentRunnerReplayState{AgentRunnerReplayStateAbsent},
		},
		{
			stage: replaySlotStageHeaderWritten,
			kill:  []AgentRunnerReplayState{AgentRunnerReplayStateDispatchedUnknown},
			// The header may or may not have reached the disk: absent is still
			// legal, and so is the committed state.
			powerLoss: []AgentRunnerReplayState{AgentRunnerReplayStateAbsent, AgentRunnerReplayStateDispatchedUnknown},
		},
		{
			stage:     replaySlotStageHeaderSynced,
			kill:      []AgentRunnerReplayState{AgentRunnerReplayStateDispatchedUnknown},
			powerLoss: []AgentRunnerReplayState{AgentRunnerReplayStateDispatchedUnknown},
		},
	}
}

func c6StateIn(state AgentRunnerReplayState, allowed []AgentRunnerReplayState) bool {
	for _, candidate := range allowed {
		if candidate == state {
			return true
		}
	}
	return false
}

// TestC6AllowedSetTableProcessKill walks the kill column: page cache is intact,
// so a reader sees exactly what the process wrote before it stopped.
func TestC6AllowedSetTableProcessKill(t *testing.T) {
	for _, testCase := range c6AllowedSetCases() {
		t.Run(testCase.stage, func(t *testing.T) {
			if len(testCase.kill) != 1 {
				t.Fatalf("the kill column must be a single allowed state, got %v", testCase.kill)
			}
			store := c6PrebuiltStore(t, "request-one")
			c6RecordRequestAndStopAtStage(t, store, replayStoreRequestRecord(t, "request-one"), testCase.stage)

			state, err := store.State("request-one")
			if err != nil {
				t.Fatalf("kill at %s must never be corruption, got %v", testCase.stage, err)
			}
			if state != testCase.kill[0] {
				t.Fatalf("kill at %s: state = %s, want %s", testCase.stage, state, testCase.kill[0])
			}
		})
	}
}

// c6ModelPowerLoss rewrites the slot as the bytes a power loss could have left
// durable. dropLastStep models losing the newest unflushed step; the other arm
// models it surviving the power loss.
func c6ModelPowerLoss(t *testing.T, path, stage string, dropLastStep bool) {
	t.Helper()
	raw := c6ReadFile(t, path)
	switch stage {
	case replaySlotStageCreated:
		raw = raw[:replaySlotHeaderSize]
	case replaySlotStageWritten:
		if dropLastStep {
			// The payload was never flushed, so it need not have reached the disk.
			raw = raw[:replaySlotHeaderSize]
		}
	case replaySlotStageHeaderWritten:
		if dropLastStep {
			// The header commit was never flushed: the disk may still hold payloadLen 0.
			binary.LittleEndian.PutUint64(raw[len(replaySlotMagic):replaySlotHeaderSize], 0)
		}
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		t.Fatalf("model power loss at %s: %v", stage, err)
	}
}

// TestC6AllowedSetTablePowerLoss walks the power-loss column: the durable bytes
// decide what a reader sees, so the header-written row admits both outcomes and
// nothing outside the table.
func TestC6AllowedSetTablePowerLoss(t *testing.T) {
	for _, testCase := range c6AllowedSetCases() {
		t.Run(testCase.stage, func(t *testing.T) {
			observed := make(map[AgentRunnerReplayState]struct{})
			for _, dropLastStep := range []bool{true, false} {
				store := c6PrebuiltStore(t, "request-one")
				c6RecordRequestAndStopAtStage(t, store, replayStoreRequestRecord(t, "request-one"), testCase.stage)
				requestPath, _ := c6SlotPaths(t, store, "request-one")
				c6ModelPowerLoss(t, requestPath, testCase.stage, dropLastStep)

				state, err := store.State("request-one")
				if err != nil {
					t.Fatalf("power loss at %s (drop=%t) must never be corruption, got %v", testCase.stage, dropLastStep, err)
				}
				if !c6StateIn(state, testCase.powerLoss) {
					t.Fatalf("power loss at %s (drop=%t): state = %s, allowed %v", testCase.stage, dropLastStep, state, testCase.powerLoss)
				}
				observed[state] = struct{}{}
			}
			// The two models must together cover the table entry, otherwise the row
			// would be satisfied by a scenario that can never occur.
			for _, allowed := range testCase.powerLoss {
				if _, ok := observed[allowed]; !ok {
					t.Fatalf("power loss at %s never produced the allowed state %s: observed %v", testCase.stage, allowed, observed)
				}
			}
		})
	}
}

// TestC6AllowedSetCorruptionIsFailClosed pins the note under the table:
// corruption is never in the allowed set, and the reader fails closed instead of
// classifying.
func TestC6AllowedSetCorruptionIsFailClosed(t *testing.T) {
	store := c6PrebuiltStore(t, "request-one")
	requestPath, _ := c6SlotPaths(t, store, "request-one")
	c6WriteFile(t, requestPath, c6RawSlotBytes([]byte("BADMAGIC"), 0, nil))
	if _, err := store.State("request-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("corrupted slot state error = %v, want corruption", err)
	}
}

// ---------------------------------------------------------------------------
// 6. BA3: two-arm comparison
// ---------------------------------------------------------------------------

// c6RunSequence runs one legal operation sequence S and returns its observable
// return sequence.
func c6RunSequence(t *testing.T, store *AgentRunnerReplayStore) []string {
	t.Helper()
	record := replayStoreRequestRecord(t, "request-one")
	completion := replayStoreCompletionRecord(t, record, "completion-one")
	var seq []string
	appendState := func(requestID string) {
		state, err := store.State(requestID)
		seq = append(seq, fmt.Sprintf("state:%s=%s err=%v", requestID, state, err))
	}
	appendState("request-one")
	decision, err := store.RecordRequest(record)
	seq = append(seq, fmt.Sprintf("recordRequest:first=%s err=%v", decision, err))
	decision, err = store.RecordRequest(record)
	seq = append(seq, fmt.Sprintf("recordRequest:replay=%s err=%v", decision, err))
	appendState("request-one")
	replayed, err := store.RecordCompletion(record, completion)
	seq = append(seq, fmt.Sprintf("recordCompletion:first replayed=%t err=%v", replayed, err))
	replayed, err = store.RecordCompletion(record, completion)
	seq = append(seq, fmt.Sprintf("recordCompletion:replay replayed=%t err=%v", replayed, err))
	appendState("request-one")
	appendState("request-two")
	return seq
}

// TestC6BA3PrebuildArmsAgreeOnReturnSequenceAndTree is the BA3 comparison under
// store-level activation: both arms start from an equivalent unprebuilt replay
// root, the control arm runs S on the ordinary record path, and the call arm runs
// PrebuildReplaySlots first and runs the same S on the slot path.
//
// The comparison is the return sequence item for item plus the R/C record
// content per requestId and kind: the control arm's plain JSONL bytes against the
// call arm's slot payload with the 16 byte header stripped. prebuild-set.jsonl
// and the slot headers are activation evidence and stay out of it - a whole-tree
// byte comparison would compare time.Now() stamps in the set record instead of
// record content.
func TestC6BA3PrebuildArmsAgreeOnReturnSequenceAndTree(t *testing.T) {
	declared := []string{"request-one", "request-two"}

	// Equivalent initial state: two independent empty replay roots, neither
	// prebuilt and neither carrying a slot file.
	control := replayStoreMustNew(t, t.TempDir())
	call := replayStoreMustNew(t, t.TempDir())

	controlSequence := c6RunSequence(t, control)
	if err := call.PrebuildReplaySlots(declared); err != nil {
		t.Fatalf("call arm prebuild: %v", err)
	}
	callSequence := c6RunSequence(t, call)

	if control.slotsActive {
		t.Fatal("the control arm never prebuilds, so it must stay on the ordinary record path")
	}
	if !call.slotsActive {
		t.Fatal("the call arm must be activated by its own prebuild")
	}
	if strings.Join(controlSequence, "\n") != strings.Join(callSequence, "\n") {
		t.Fatalf("BA3 return sequences differ:\ncontrol=%v\ncall=%v", controlSequence, callSequence)
	}

	for _, armCase := range []struct {
		requestID  string
		wantRecord bool
	}{
		{requestID: "request-one", wantRecord: true},
		{requestID: "request-two", wantRecord: false},
	} {
		controlRequestPath, controlCompletionPath := c6SlotPaths(t, control, armCase.requestID)
		callRequestPath, callCompletionPath := c6SlotPaths(t, call, armCase.requestID)
		controlRequest := c6PlainRecord(t, controlRequestPath)
		callRequest := c6SlotPayload(t, callRequestPath)
		controlCompletion := c6PlainRecord(t, controlCompletionPath)
		callCompletion := c6SlotPayload(t, callCompletionPath)
		c6RequireArmPayloadsAgree(t, armCase.requestID+" request", controlRequest, callRequest)
		c6RequireArmPayloadsAgree(t, armCase.requestID+" completion", controlCompletion, callCompletion)
		if !armCase.wantRecord {
			// Nothing was published for this requestId: an empty plain record and a
			// prebuilt but uncommitted slot must both read as absent.
			if controlRequest != nil || callRequest != nil || controlCompletion != nil || callCompletion != nil {
				t.Fatalf("%s: no record was published, both arms must read absent", armCase.requestID)
			}
			continue
		}
		// Two absences must not be able to satisfy the comparison: the published
		// records have to be present in both arms.
		if controlRequest == nil || callRequest == nil {
			t.Fatalf("%s: both arms must hold the published request record", armCase.requestID)
		}
		if controlCompletion == nil || callCompletion == nil {
			t.Fatalf("%s: both arms must hold the published completion record", armCase.requestID)
		}
	}
}

// TestC6BA3RecordRequestFirstDispatchAndReplayUnknownBlock pins the interface BA3
// must not treat as a side-effect-free probe: RecordRequest writes on the first
// call and reports a replay on the second, in both arms and from equivalent
// unprebuilt initial states - the control arm on the ordinary record path, the
// call arm on the slot path after its own prebuild.
func TestC6BA3RecordRequestFirstDispatchAndReplayUnknownBlock(t *testing.T) {
	declared := []string{"request-one", "request-two"}
	record := replayStoreRequestRecord(t, "request-one")

	// Equivalent initial state: two independent empty replay roots, neither
	// prebuilt.
	control := replayStoreMustNew(t, t.TempDir())
	call := replayStoreMustNew(t, t.TempDir())
	if err := call.PrebuildReplaySlots(declared); err != nil {
		t.Fatalf("call arm prebuild: %v", err)
	}

	type c6ArmResult struct {
		sequence []string
		payload  []byte
	}
	runArm := func(store *AgentRunnerReplayStore, slotPath bool) c6ArmResult {
		t.Helper()
		var sequence []string
		decision, err := store.RecordRequest(record)
		sequence = append(sequence, fmt.Sprintf("first=%s err=%v", decision, err))
		if err != nil || decision != AgentRunnerReplayDecisionFirstDispatch {
			t.Fatalf("first RecordRequest = %s err=%v, want first-dispatch", decision, err)
		}
		decision, err = store.RecordRequest(record)
		sequence = append(sequence, fmt.Sprintf("replay=%s err=%v", decision, err))
		if err != nil || decision != AgentRunnerReplayDecisionUnknownBlock {
			t.Fatalf("replayed RecordRequest = %s err=%v, want unknown-block", decision, err)
		}
		state, err := store.State("request-one")
		sequence = append(sequence, fmt.Sprintf("state=%s err=%v", state, err))
		if err != nil || state != AgentRunnerReplayStateDispatchedUnknown {
			t.Fatalf("state = %s err=%v, want dispatched-unknown", state, err)
		}
		requestPath, _ := c6SlotPaths(t, store, "request-one")
		payload := c6PlainRecord(t, requestPath)
		if slotPath {
			payload = c6SlotPayload(t, requestPath)
		}
		return c6ArmResult{sequence: sequence, payload: payload}
	}

	controlArm := runArm(control, false)
	callArm := runArm(call, true)
	if strings.Join(controlArm.sequence, "\n") != strings.Join(callArm.sequence, "\n") {
		t.Fatalf("arms disagree on RecordRequest: control=%v call=%v", controlArm.sequence, callArm.sequence)
	}
	c6RequireArmPayloadsAgree(t, "request-one R", controlArm.payload, callArm.payload)

	// The published bytes must be the record itself in both containers, so an
	// agreement between two wrong payloads cannot pass unnoticed.
	want := c6Payload(t, record)
	if !bytes.Equal(controlArm.payload, want) {
		t.Fatalf("control arm request record = %q, want the canonical JSONL payload %q", controlArm.payload, want)
	}
	if !bytes.Equal(callArm.payload, want) {
		t.Fatalf("call arm request payload = %q, want the same canonical JSONL payload %q", callArm.payload, want)
	}
}

// ---------------------------------------------------------------------------
// 7. gate invariants under c6slot activation
// ---------------------------------------------------------------------------

// TestC6GateActivationKeepsDurabilityUnproven pins that activating the slot
// implementation does not move the platform durability claim.
func TestC6GateActivationKeepsDurabilityUnproven(t *testing.T) {
	if !replayStoreSlotEnabled {
		t.Fatal("the c6slot build must activate the slot implementation, otherwise this file proves nothing")
	}
	store := c6UnprovenStore(t)
	if got := store.PublicationDurability(); got != PublishDurabilityUnproven {
		t.Fatalf("PublicationDurability = %s, want unproven", got)
	}
}

// TestC6GateRecordRequestStillRejectsUnprovenDurability pins that the first
// dispatch is still refused: no first-dispatch decision and no R record.
func TestC6GateRecordRequestStillRejectsUnprovenDurability(t *testing.T) {
	store := c6UnprovenStore(t)
	if err := store.PrebuildReplaySlots([]string{"request-one"}); err != nil {
		t.Fatalf("prebuild: %v", err)
	}
	decision, err := store.RecordRequest(replayStoreRequestRecord(t, "request-one"))
	if !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("RecordRequest error = %v, want durability unproven", err)
	}
	if decision != "" {
		t.Fatalf("RecordRequest decision = %q, want empty", decision)
	}
	if decision == AgentRunnerReplayDecisionFirstDispatch {
		t.Fatal("RecordRequest must never return first-dispatch while durability is unproven")
	}
	requestPath, _ := c6SlotPaths(t, store, "request-one")
	if payload, found, readErr := readReplaySlot(requestPath); readErr != nil || found || payload != nil {
		t.Fatalf("a refused dispatch must not create R: found=%t err=%v", found, readErr)
	}
	if state, stateErr := store.State("request-one"); stateErr != nil || state != AgentRunnerReplayStateAbsent {
		t.Fatalf("state after a refused dispatch = %s err=%v, want absent", state, stateErr)
	}
}

// TestC6GateRecordCompletionReplaysWithoutWrite pins that a published completion
// replays without any durability claim or write, while a completion that would
// need a write is still refused.
func TestC6GateRecordCompletionReplaysWithoutWrite(t *testing.T) {
	store := c6UnprovenStore(t)
	if err := store.PrebuildReplaySlots([]string{"request-one", "request-two"}); err != nil {
		t.Fatalf("prebuild: %v", err)
	}

	request := replayStoreRequestRecord(t, "request-one")
	completion := replayStoreCompletionRecord(t, request, "completion-one")
	requestPath, completionPath := c6SlotPaths(t, store, "request-one")
	if err := writeReplaySlotNoReplace(requestPath, c6Payload(t, request)); err != nil {
		t.Fatalf("seed R slot: %v", err)
	}
	if err := writeReplaySlotNoReplace(completionPath, c6Payload(t, completion)); err != nil {
		t.Fatalf("seed C slot: %v", err)
	}
	before := c6ReadFile(t, completionPath)

	replayed, err := store.RecordCompletion(request, completion)
	if err != nil || !replayed {
		t.Fatalf("replayed completion = %t err=%v, want a replay without a durability claim", replayed, err)
	}
	if after := c6ReadFile(t, completionPath); !bytes.Equal(before, after) {
		t.Fatal("a replayed completion must not write the C slot")
	}

	// A completion that needs a write is still refused while durability is unproven.
	requestTwo := replayStoreRequestRecord(t, "request-two")
	completionTwo := replayStoreCompletionRecord(t, requestTwo, "completion-two")
	requestPathTwo, completionPathTwo := c6SlotPaths(t, store, "request-two")
	if err := writeReplaySlotNoReplace(requestPathTwo, c6Payload(t, requestTwo)); err != nil {
		t.Fatalf("seed R slot for request-two: %v", err)
	}
	if _, err := store.RecordCompletion(requestTwo, completionTwo); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("new completion error = %v, want durability unproven", err)
	}
	if payload, found, readErr := readReplaySlot(completionPathTwo); readErr != nil || found || payload != nil {
		t.Fatalf("a refused completion must not write C: found=%t err=%v", found, readErr)
	}
}
