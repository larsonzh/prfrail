package tickets

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// Test-only contract for the durable cost ledger (design:
// docs/validation/evidence/P39-T027B-DURABLE-BUDGET/01-arch-design/review.md).
// These tests intentionally reference an API that does not exist yet: the
// compile failure is the expected RED state of this slice.

const (
	durableHelperEnv = "PRFRAIL_DURABLE_HELPER"
	durableRootEnv   = "PRFRAIL_DURABLE_ROOT"
	durableStartEnv  = "PRFRAIL_DURABLE_START"
	durableIndexEnv  = "PRFRAIL_DURABLE_INDEX"
	durableHelperRun = "^TestDurableHelperReserve$"

	durableExitOK          = 0
	durableExitCapExceeded = 3
)

func durableConfig(root string, amountMicros *int64, callCap, tokenCap int) DurableCostLedgerConfig {
	return DurableCostLedgerConfig{
		Root: root,
		Ledger: CostLedgerConfig{
			LedgerID:       "durable-cost-ledger-one",
			CreatedAt:      costNow(),
			ScopeKind:      CostScopeProject,
			ScopeHash:      costHashOne,
			Currency:       "USD",
			PricingMode:    CostPricingDocumented,
			PricingVersion: "2026-09-08",
		},
		Initial: AllocationInput{
			EntryID:           "allocation-one",
			OccurredAt:        costNow(),
			AuthorizationHash: costHashTwo,
			Limits: CostLimits{
				AmountMicros: amountMicros,
				ModelCalls:   callCap,
				Tokens:       tokenCap,
				WallClockMs:  60000,
				Attempts:     3,
			},
		},
	}
}

func openDurable(t *testing.T, root string, amountMicros *int64, callCap, tokenCap int) *DurableCostLedger {
	t.Helper()
	ledger, err := OpenDurableCostLedger(durableConfig(root, amountMicros, callCap, tokenCap))
	if err != nil {
		t.Fatalf("open durable ledger: %v", err)
	}
	t.Cleanup(func() { _ = ledger.Close() })
	return ledger
}

func durableReservationInput(entryID, requestID, key string, amount int64, calls, tokens int) ReservationInput {
	return ReservationInput{
		EntryID:              entryID,
		OccurredAt:           costNow(),
		RunID:                "run-one",
		RequestID:            requestID,
		IdempotencyKey:       key,
		AuthorizationHash:    costHashThree,
		ReservedAmountMicros: int64Pointer(amount),
		ReservedCalls:        calls,
		ReservedTokens:       tokens,
		PricingEvidence:      []string{costHashFour},
	}
}

// reloadDurableLedger replays the canonical ledger file from disk, so every
// assertion below is a disk fact rather than an in-memory one.
func reloadDurableLedger(t *testing.T, path string) *CostLedger {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read durable ledger %s: %v", path, err)
	}
	record, err := DecodeCostLedgerRecord(raw)
	if err != nil {
		t.Fatalf("decode durable ledger: %v", err)
	}
	ledger, err := LoadCostLedger(record)
	if err != nil {
		t.Fatalf("load durable ledger: %v", err)
	}
	return ledger
}

func durableEntryCounts(t *testing.T, ledger *CostLedger) map[string]int {
	t.Helper()
	record, err := ledger.Snapshot()
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	counts := map[string]int{"allocation": 0, "reservation": 0, "settlement": 0}
	for _, entry := range record.Ledger.Entries {
		switch entry.(type) {
		case CostAllocationEntry:
			counts["allocation"]++
		case CostReservationEntry:
			counts["reservation"]++
		case CostSettlementEntry:
			counts["settlement"]++
		}
	}
	return counts
}

func durableReservedAmount(t *testing.T, ledger *CostLedger) int64 {
	t.Helper()
	summary := ledger.Summary()
	if summary.ReservedAmountMicros == nil {
		t.Fatalf("reserved amount is nil")
	}
	return *summary.ReservedAmountMicros
}

func durableUnknownAmount(t *testing.T, ledger *CostLedger) int64 {
	t.Helper()
	summary := ledger.Summary()
	if summary.UnknownReservedAmountMicros == nil {
		t.Fatalf("unknown reserved amount is nil")
	}
	return *summary.UnknownReservedAmountMicros
}

func durableSettlementEntry(t *testing.T, ledger *CostLedger) CostSettlementEntry {
	t.Helper()
	record, err := ledger.Snapshot()
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	for _, entry := range record.Ledger.Entries {
		if settlement, ok := entry.(CostSettlementEntry); ok {
			return settlement
		}
	}
	t.Fatalf("no settlement entry on reloaded ledger")
	return CostSettlementEntry{}
}

// 1. Sequential Reserve against one shared scope: the cap must reject the
// oversized reservation and every accepted reservation must be durable.
func TestDurableReserveSingleProcessSequence(t *testing.T) {
	amount := int64(100)
	ledger := openDurable(t, t.TempDir(), &amount, 10, 1000)
	ctx := context.Background()

	steps := []struct {
		name            string
		amountMicros    int64
		wantErr         error
		wantReserved    int64
		wantOutstanding int
	}{
		{"first", 60, nil, 60, 1},
		{"exceeds-cap", 50, ErrCostSharedCapExceeded, 60, 1},
		{"fills-to-cap", 40, nil, 100, 2},
		{"over-cap", 1, ErrCostSharedCapExceeded, 100, 2},
	}
	for index, step := range steps {
		outcome, ack, err := ledger.Reserve(ctx, durableReservationInput(
			"reservation-"+strconv.Itoa(index),
			"request-"+strconv.Itoa(index),
			"reserve-key-"+strconv.Itoa(index),
			step.amountMicros, 1, 10,
		))
		if step.wantErr != nil {
			if !errors.Is(err, step.wantErr) {
				t.Fatalf("%s: want %v, got %v", step.name, step.wantErr, err)
			}
			if ack.Valid() {
				t.Fatalf("%s: rejected reservation must not carry a valid ack", step.name)
			}
			if outcome.Appended || outcome.ReservationHash != "" {
				t.Fatalf("%s: rejected reservation must not report an appended outcome", step.name)
			}
		} else {
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", step.name, err)
			}
			if !ack.Valid() {
				t.Fatalf("%s: accepted reservation must carry a valid ack", step.name)
			}
			if err := ack.RequireEntry(outcome.ReservationHash); err != nil {
				t.Fatalf("%s: ack does not cover reservation entry: %v", step.name, err)
			}
			if !outcome.Appended {
				t.Fatalf("%s: first reservation for this key must be appended", step.name)
			}
		}

		reloaded := reloadDurableLedger(t, ledger.Path())
		if got := durableReservedAmount(t, reloaded); got != step.wantReserved {
			t.Fatalf("%s: reserved amount = %d, want %d", step.name, got, step.wantReserved)
		}
		if got := reloaded.OutstandingReservations(); got != step.wantOutstanding {
			t.Fatalf("%s: outstanding = %d, want %d", step.name, got, step.wantOutstanding)
		}
	}
}

// TestDurableHelperReserve is the child process used by the concurrent test.
// It must stay behind the environment switch so a normal run skips it.
func TestDurableHelperReserve(t *testing.T) {
	if os.Getenv(durableHelperEnv) != "reserve" {
		t.Skip("helper process only")
	}
	root := os.Getenv(durableRootEnv)
	startFile := os.Getenv(durableStartEnv)
	index := os.Getenv(durableIndexEnv)

	amount := int64(100)
	ledger, err := OpenDurableCostLedger(durableConfig(root, &amount, 10, 1000))
	if err != nil {
		os.Exit(12)
	}

	deadline := time.Now().Add(60 * time.Second)
	for {
		if _, statErr := os.Stat(startFile); statErr == nil {
			break
		}
		if time.Now().After(deadline) {
			os.Exit(11)
		}
		time.Sleep(5 * time.Millisecond)
	}

	outcome, ack, err := ledger.Reserve(context.Background(), durableReservationInput(
		"reservation-"+index, "request-"+index, "reserve-key-"+index, 70, 1, 10,
	))
	if err != nil {
		if errors.Is(err, ErrCostSharedCapExceeded) {
			os.Exit(durableExitCapExceeded)
		}
		os.Exit(13)
	}
	if !ack.Valid() {
		os.Exit(14)
	}
	if err := ack.RequireEntry(outcome.ReservationHash); err != nil {
		os.Exit(15)
	}
	if err := ledger.Close(); err != nil {
		os.Exit(16)
	}
	os.Exit(durableExitOK)
}

// 2. Two real processes reserve against the same scope; exactly one may win
// and the durable ledger must contain no duplicate entry.
func TestDurableReserveConcurrentProcesses(t *testing.T) {
	amount := int64(100)
	root := t.TempDir()
	ledger := openDurable(t, root, &amount, 10, 1000)
	path := ledger.Path()
	if err := ledger.Close(); err != nil {
		t.Fatalf("close bootstrap ledger: %v", err)
	}

	startFile := filepath.Join(root, "start.signal")
	cmds := make([]*exec.Cmd, 0, 2)
	for index := 0; index < 2; index++ {
		cmd := exec.Command(os.Args[0], "-test.run="+durableHelperRun)
		cmd.Env = append(os.Environ(),
			durableHelperEnv+"=reserve",
			durableRootEnv+"="+root,
			durableStartEnv+"="+startFile,
			durableIndexEnv+"="+strconv.Itoa(index),
		)
		if err := cmd.Start(); err != nil {
			t.Fatalf("start helper %d: %v", index, err)
		}
		cmds = append(cmds, cmd)
	}

	if err := os.WriteFile(startFile, []byte("go"), 0o600); err != nil {
		t.Fatalf("release start barrier: %v", err)
	}

	codes := make([]int, len(cmds))
	for index, cmd := range cmds {
		err := cmd.Wait()
		if err == nil {
			codes[index] = durableExitOK
			continue
		}
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("helper %d did not exit normally: %v", index, err)
		}
		codes[index] = exitErr.ExitCode()
	}

	winners, rejected := 0, 0
	for _, code := range codes {
		switch code {
		case durableExitOK:
			winners++
		case durableExitCapExceeded:
			rejected++
		}
	}
	if winners != 1 || rejected != 1 {
		t.Fatalf("want exactly one winner and one cap rejection, got exit codes %v", codes)
	}

	reloaded := reloadDurableLedger(t, path)
	if got := durableReservedAmount(t, reloaded); got != 70 {
		t.Fatalf("reserved amount after race = %d, want 70", got)
	}
	if got := reloaded.OutstandingReservations(); got != 1 {
		t.Fatalf("outstanding after race = %d, want 1", got)
	}
	counts := durableEntryCounts(t, reloaded)
	if counts["allocation"] != 1 || counts["reservation"] != 1 || counts["settlement"] != 0 {
		t.Fatalf("duplicate or missing entries after race: %+v", counts)
	}
}

// 3. Settlement idempotency is keyed by the derived terminal key.
func TestDurableSettleIdempotentByKey(t *testing.T) {
	amount := int64(1000)
	ledger := openDurable(t, t.TempDir(), &amount, 10, 1000)
	ctx := context.Background()

	outcome, ack, err := ledger.Reserve(ctx, durableReservationInput("reservation-one", "request-one", "reserve-key-one", 60, 1, 10))
	if err != nil {
		t.Fatalf("reserve: %v", err)
	}
	if !ack.Valid() {
		t.Fatalf("reserve ack invalid")
	}

	charged := int64(40)
	input := SettlementInput{
		EntryID:             "settlement-one",
		OccurredAt:          costNow(),
		RunID:               "run-one",
		RequestID:           "request-one",
		IdempotencyKey:      stableTerminalSettlementKey("request-one", costHashFour),
		ReservationHash:     outcome.ReservationHash,
		Status:              CostSettlementObserved,
		ChargedAmountMicros: int64Pointer(charged),
		ProviderEvidence:    []string{costHashThree},
	}

	first, firstAck, err := ledger.Settle(ctx, input)
	if err != nil {
		t.Fatalf("first settle: %v", err)
	}
	if !firstAck.Valid() || !first.Appended {
		t.Fatalf("first settle must be appended and confirmed: appended=%v ack=%v", first.Appended, firstAck.Valid())
	}
	if err := firstAck.RequireEntry(first.EntryHash); err != nil {
		t.Fatalf("first settle ack does not cover entry: %v", err)
	}

	second, secondAck, err := ledger.Settle(ctx, input)
	if err != nil {
		t.Fatalf("idempotent settle must not fail: %v", err)
	}
	if second.Appended {
		t.Fatalf("idempotent settle must not append a second entry")
	}
	if second.EntryHash != first.EntryHash {
		t.Fatalf("idempotent settle returned %q, want %q", second.EntryHash, first.EntryHash)
	}
	if !secondAck.Valid() {
		t.Fatalf("idempotent hit must still be durably confirmed")
	}
	if err := secondAck.RequireEntry(first.EntryHash); err != nil {
		t.Fatalf("idempotent ack does not cover entry: %v", err)
	}

	reloaded := reloadDurableLedger(t, ledger.Path())
	if counts := durableEntryCounts(t, reloaded); counts["settlement"] != 1 {
		t.Fatalf("settlement entries = %d, want 1: %+v", counts["settlement"], counts)
	}

	other := input
	other.EntryID = "settlement-two"
	other.IdempotencyKey = stableTerminalSettlementKey("request-one", costHashTwo)
	if _, otherAck, err := ledger.Settle(ctx, other); !errors.Is(err, ErrCostReservationSettled) {
		t.Fatalf("second key for same reservation: want %v, got %v", ErrCostReservationSettled, err)
	} else if otherAck.Valid() {
		t.Fatalf("rejected settlement must not carry a valid ack")
	}

	reloaded = reloadDurableLedger(t, ledger.Path())
	if counts := durableEntryCounts(t, reloaded); counts["settlement"] != 1 {
		t.Fatalf("rejected settlement mutated the ledger: %+v", counts)
	}
}

// 4. Unknown settlement ends outstanding eligibility but keeps the reservation
// charged against the shared cap across a real reload.
func TestDurableUnknownKeepsOccupancyAfterReload(t *testing.T) {
	amount := int64(1000)
	root := t.TempDir()
	ledger := openDurable(t, root, &amount, 3, 300)
	ctx := context.Background()

	outcome, _, err := ledger.Reserve(ctx, durableReservationInput("reservation-one", "request-one", "reserve-key-one", 400, 1, 80))
	if err != nil {
		t.Fatalf("reserve: %v", err)
	}
	unknown := SettlementInput{
		EntryID:          "settlement-unknown",
		OccurredAt:       costNow(),
		RunID:            "run-one",
		RequestID:        "request-one",
		IdempotencyKey:   stableTerminalSettlementKey("request-one", costHashFour),
		ReservationHash:  outcome.ReservationHash,
		Status:           CostSettlementUnknown,
		ProviderEvidence: []string{costHashThree},
	}
	if _, ack, err := ledger.Settle(ctx, unknown); err != nil {
		t.Fatalf("unknown settle: %v", err)
	} else if !ack.Valid() {
		t.Fatalf("unknown settle ack invalid")
	}

	path := ledger.Path()
	if err := ledger.Close(); err != nil {
		t.Fatalf("close before reopen: %v", err)
	}

	reopened, err := OpenDurableCostLedger(durableConfig(root, &amount, 3, 300))
	if err != nil {
		t.Fatalf("reopen durable ledger: %v", err)
	}
	t.Cleanup(func() { _ = reopened.Close() })

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read ledger bytes: %v", err)
	}
	if !strings.Contains(string(raw), "\"chargedAmountMicros\":null") {
		t.Fatalf("unknown settlement must serialize chargedAmountMicros as null")
	}

	reloaded := reloadDurableLedger(t, path)
	settlement := durableSettlementEntry(t, reloaded)
	if settlement.Status != CostSettlementUnknown {
		t.Fatalf("settlement status = %q, want %q", settlement.Status, CostSettlementUnknown)
	}
	if settlement.ChargedAmountMicros != nil || settlement.ObservedCalls != nil || settlement.ObservedTokens != nil {
		t.Fatalf("unknown settlement must keep nil amount/calls/tokens")
	}
	if got := reloaded.OutstandingReservations(); got != 0 {
		t.Fatalf("outstanding after unknown = %d, want 0", got)
	}
	if got := reloaded.UnknownHoldReservations(); got != 1 {
		t.Fatalf("unknown hold after reload = %d, want 1", got)
	}
	if got := durableUnknownAmount(t, reloaded); got != 400 {
		t.Fatalf("unknown reserved amount = %d, want 400", got)
	}
	if _, err := reopened.RequireOutstandingReservation(ctx, outcome.ReservationHash); !errors.Is(err, ErrCostReservationSettled) {
		t.Fatalf("unknown settlement must end outstanding eligibility, got %v", err)
	}

	if _, ack, err := reopened.Reserve(ctx, durableReservationInput("reservation-two", "request-two", "reserve-key-two", 700, 1, 80)); !errors.Is(err, ErrCostSharedCapExceeded) {
		t.Fatalf("unknown hold must still consume the cap, got %v (ack valid=%v)", err, ack.Valid())
	}
	if _, _, err := reopened.Reserve(ctx, durableReservationInput("reservation-three", "request-three", "reserve-key-three", 600, 1, 80)); err != nil {
		t.Fatalf("reservation up to the remaining cap must succeed: %v", err)
	}
	if got := durableReservedAmount(t, reloadDurableLedger(t, path)); got != 1000 {
		t.Fatalf("reserved amount after refill = %d, want 1000", got)
	}
}

// 5. A failed commit must return an error with no valid ack and must not roll
// back or delete the snapshot that was already visible on disk.
func TestDurableCommitFailureKeepsVisibleRecord(t *testing.T) {
	amount := int64(1000)
	ledger := openDurable(t, t.TempDir(), &amount, 10, 1000)
	ctx := context.Background()

	path := ledger.Path()
	if _, ack, err := ledger.Reserve(ctx, durableReservationInput("reservation-one", "request-one", "reserve-key-one", 60, 1, 10)); err != nil {
		t.Fatalf("bootstrap reserve: %v", err)
	} else if !ack.Valid() {
		t.Fatalf("bootstrap reserve ack invalid")
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read ledger before fault: %v", err)
	}

	// The pre-commit fault is injected through the design's fault surface
	// (Config.CommitFault) instead of a permission heuristic, so this leg runs
	// on every platform without skipping.
	injected := errors.New("injected pre-commit write fault")
	faultConfig := durableConfig(filepath.Dir(path), &amount, 10, 1000)
	faultConfig.CommitFault = func(stage string) error {
		if stage == StageWrite {
			return injected
		}
		return nil
	}
	faulted, err := OpenDurableCostLedger(faultConfig)
	if err != nil {
		t.Fatalf("open faulted ledger: %v", err)
	}
	defer func() { _ = faulted.Close() }()

	_, ack, err := faulted.Reserve(ctx, durableReservationInput("reservation-two", "request-two", "reserve-key-two", 10, 1, 10))
	if err == nil {
		t.Fatalf("commit failure must be reported")
	}
	if !errors.Is(err, ErrDurableNotCommitted) {
		t.Fatalf("want %v, got %v", ErrDurableNotCommitted, err)
	}
	if !errors.Is(err, injected) {
		t.Fatalf("injected cause must be preserved: %v", err)
	}
	if ack.Valid() {
		t.Fatalf("failed commit must not carry a valid ack")
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("visible record was removed by a failed commit: %v", err)
	}
	if string(before) != string(after) {
		t.Fatalf("failed commit mutated the visible ledger snapshot")
	}

	reloaded := reloadDurableLedger(t, path)
	if got := durableReservedAmount(t, reloaded); got != 60 {
		t.Fatalf("reserved amount after failed commit = %d, want 60", got)
	}
	if counts := durableEntryCounts(t, reloaded); counts["reservation"] != 1 {
		t.Fatalf("failed commit added entries: %+v", counts)
	}
	if got := reloaded.OutstandingReservations(); got != 1 {
		t.Fatalf("outstanding after failed commit = %d, want 1", got)
	}
}

// ---------------------------------------------------------------------------
// P39 durable budget matrix tests (see
// docs/validation/evidence/P39-T027B-DURABLE-BUDGET/01-arch-design/review-batch-b.md).
// Every test below is fail-to-pass: removing the scope lock, the in-lock
// re-read, the prefix verification, the read-only fail-closed path or the
// commit-stage ordering turns the corresponding assertion red.
// ---------------------------------------------------------------------------

const p39DeadlineHelperRun = "^TestP39DurableHelperDeadlineReserve$"

func p39FileCount(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir %s: %v", dir, err)
	}
	return len(entries)
}

func p39LockFileCount(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir %s: %v", dir, err)
	}
	count := 0
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".lock") {
			count++
		}
	}
	return count
}

// p39LoadRecord decodes the canonical ledger bytes without going through the
// durable wrapper, so tampering helpers can rebuild a self-consistent record.
func p39LoadRecord(t *testing.T, path string) CostLedgerRecord {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	record, err := DecodeCostLedgerRecord(raw)
	if err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return record
}

// p39WriteRecord re-seals a record with a fresh ledger hash and writes it back,
// which is exactly what an attacker with disk access would do.
func p39WriteRecord(t *testing.T, path string, record CostLedgerRecord) []byte {
	t.Helper()
	hash, err := digest(costLedgerDomain, record.Ledger)
	if err != nil {
		t.Fatalf("recompute ledger hash: %v", err)
	}
	record.LedgerHash = hash
	canonical, err := evidence.EncodeCanonical(record)
	if err != nil {
		t.Fatalf("encode tampered record: %v", err)
	}
	if _, err := DecodeCostLedgerRecord(canonical); err != nil {
		t.Fatalf("tampered fixture must still decode self-consistently: %v", err)
	}
	if err := os.WriteFile(path, canonical, 0o600); err != nil {
		t.Fatalf("write tampered record: %v", err)
	}
	return canonical
}

// TestP39DurableHelperDeadlineReserve is the child of
// TestP39DurableReserveTwoProcessesOneWinner. It opens the ledger (old
// snapshot) first and only then waits for the shared deadline barrier, so both
// processes contend on the same scope lock at the same wall-clock instant.
func TestP39DurableHelperDeadlineReserve(t *testing.T) {
	if os.Getenv(durableHelperEnv) != "deadline" {
		t.Skip("helper process only")
	}
	root := os.Getenv(durableRootEnv)
	barrier := os.Getenv(durableStartEnv)
	index := os.Getenv(durableIndexEnv)

	raw, err := os.ReadFile(barrier)
	if err != nil {
		os.Exit(11)
	}
	deadline, err := strconv.ParseInt(strings.TrimSpace(string(raw)), 10, 64)
	if err != nil {
		os.Exit(11)
	}

	amount := int64(100)
	ledger, err := OpenDurableCostLedger(durableConfig(root, &amount, 10, 1000))
	if err != nil {
		os.Exit(12)
	}

	target := time.Unix(0, deadline)
	for time.Now().Before(target) {
		time.Sleep(2 * time.Millisecond)
	}

	outcome, ack, err := ledger.Reserve(context.Background(), durableReservationInput(
		"reservation-barrier-"+index, "request-barrier-"+index, "reserve-barrier-key-"+index, 70, 1, 10,
	))
	if err != nil {
		if errors.Is(err, ErrCostSharedCapExceeded) {
			os.Exit(durableExitCapExceeded)
		}
		os.Exit(13)
	}
	if !ack.Valid() {
		os.Exit(14)
	}
	if err := ack.RequireEntry(outcome.ReservationHash); err != nil {
		os.Exit(15)
	}
	if err := ledger.Close(); err != nil {
		os.Exit(16)
	}
	os.Exit(durableExitOK)
}

// Matrix #1: two real processes hold the pre-race snapshot, then cross the same
// deadline barrier and both Reserve 70 against a 100 cap. Exactly one may win.
func TestP39DurableReserveTwoProcessesOneWinner(t *testing.T) {
	amount := int64(100)
	root := t.TempDir()
	seed := openDurable(t, root, &amount, 10, 1000)
	path := seed.Path()
	if err := seed.Close(); err != nil {
		t.Fatalf("close seed ledger: %v", err)
	}

	barrier := filepath.Join(root, "deadline.barrier")
	deadline := time.Now().Add(1500 * time.Millisecond).UnixNano()
	if err := os.WriteFile(barrier, []byte(strconv.FormatInt(deadline, 10)), 0o600); err != nil {
		t.Fatalf("write deadline barrier: %v", err)
	}

	cmds := make([]*exec.Cmd, 0, 2)
	for index := 0; index < 2; index++ {
		cmd := exec.Command(os.Args[0], "-test.run="+p39DeadlineHelperRun)
		cmd.Env = append(os.Environ(),
			durableHelperEnv+"=deadline",
			durableRootEnv+"="+root,
			durableStartEnv+"="+barrier,
			durableIndexEnv+"="+strconv.Itoa(index),
		)
		if err := cmd.Start(); err != nil {
			t.Fatalf("start helper %d: %v", index, err)
		}
		cmds = append(cmds, cmd)
	}

	codes := make([]int, len(cmds))
	for index, cmd := range cmds {
		err := cmd.Wait()
		if err == nil {
			codes[index] = durableExitOK
			continue
		}
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("helper %d did not exit normally: %v", index, err)
		}
		codes[index] = exitErr.ExitCode()
	}

	winners, rejected := 0, 0
	for _, code := range codes {
		switch code {
		case durableExitOK:
			winners++
		case durableExitCapExceeded:
			rejected++
		}
	}
	if winners != 1 || rejected != 1 {
		t.Fatalf("want exactly one winner and one cap rejection, got exit codes %v", codes)
	}

	reloaded := reloadDurableLedger(t, path)
	if got := durableReservedAmount(t, reloaded); got != 70 {
		t.Fatalf("reserved amount after barrier race = %d, want 70", got)
	}
	if got := reloaded.OutstandingReservations(); got != 1 {
		t.Fatalf("outstanding after barrier race = %d, want 1", got)
	}
	counts := durableEntryCounts(t, reloaded)
	if counts["allocation"] != 1 || counts["reservation"] != 1 || counts["settlement"] != 0 {
		t.Fatalf("duplicate or missing entries after barrier race: %+v", counts)
	}
	record, err := reloaded.Snapshot()
	if err != nil {
		t.Fatalf("snapshot after barrier race: %v", err)
	}
	seen := map[string]int{}
	for _, entry := range record.Ledger.Entries {
		if reservation, ok := entry.(CostReservationEntry); ok {
			seen[reservation.IdempotencyKey]++
		}
	}
	if len(seen) != 1 {
		t.Fatalf("want exactly one reservation identity, got %v", seen)
	}
	for key, count := range seen {
		if count != 1 {
			t.Fatalf("duplicate reservation entry %q appears %d times", key, count)
		}
	}
}

// Matrix #2: the side-car lock file identity is stable across openings,
// replacements and close, and an alias root can never mint a second lock.
func TestP39DurableScopeLockStaysStable(t *testing.T) {
	amount := int64(1000)
	root := t.TempDir()
	first := openDurable(t, root, &amount, 10, 1000)
	second := openDurable(t, root, &amount, 10, 1000)

	if first.lockPath != second.lockPath {
		t.Fatalf("two openings produced different lock paths: %q vs %q", first.lockPath, second.lockPath)
	}
	if first.Path() != second.Path() {
		t.Fatalf("two openings produced different ledger paths: %q vs %q", first.Path(), second.Path())
	}

	// A successful commit replaces the ledger snapshot; the lock file must not
	// be replaced or removed by that sequence.
	beforeInfo, err := os.Stat(first.lockPath)
	if err != nil {
		t.Fatalf("stat lock file before commit: %v", err)
	}
	outcome, ack, err := first.Reserve(context.Background(), durableReservationInput("reservation-lock-one", "request-lock-one", "reserve-lock-one", 60, 1, 10))
	if err != nil {
		t.Fatalf("reserve: %v", err)
	}
	if !ack.Valid() {
		t.Fatalf("reserve ack invalid")
	}
	second.Close()
	midInfo, err := os.Stat(first.lockPath)
	if err != nil {
		t.Fatalf("lock file disappeared after commit/close: %v", err)
	}
	if !os.SameFile(beforeInfo, midInfo) {
		t.Fatalf("lock file identity changed across commit and second close")
	}

	// A lexically different but Clean-equivalent root must resolve to the very
	// same lock file, never a second independent one.
	alias := root + string(filepath.Separator) + "unused-dir" + string(filepath.Separator) + ".."
	if alias == root {
		t.Fatalf("test bug: alias root is not lexically distinct from root")
	}
	aliased, err := OpenDurableCostLedger(durableConfig(alias, &amount, 10, 1000))
	if err != nil {
		// Fail-closed is an accepted outcome, but only for a classified reason.
		if !errors.Is(err, ErrInvalidCostLedger) {
			t.Fatalf("alias root must either resolve to the same lock or fail closed, got %v", err)
		}
	} else {
		defer func() { _ = aliased.Close() }()
		if aliased.lockPath != first.lockPath {
			t.Fatalf("alias root produced a different lock path: %q vs %q", aliased.lockPath, first.lockPath)
		}
		if got := p39LockFileCount(t, root); got != 1 {
			t.Fatalf("alias root minted %d lock files, want 1", got)
		}
		aliasOutcome, aliasAck, err := aliased.Reserve(context.Background(), durableReservationInput("reservation-lock-two", "request-lock-two", "reserve-lock-two", 10, 1, 10))
		if err != nil {
			t.Fatalf("reserve through alias root: %v", err)
		}
		if !aliasAck.Valid() {
			t.Fatalf("alias reserve ack invalid")
		}
		if err := aliasAck.RequireEntry(aliasOutcome.ReservationHash); err != nil {
			t.Fatalf("alias ack does not cover its reservation: %v", err)
		}
		// Same lock and same ledger: both writes land in one canonical history,
		// and the original handle still answers for its own confirmed prefix.
		reloaded := reloadDurableLedger(t, first.Path())
		if got := durableEntryCounts(t, reloaded)["reservation"]; got != 2 {
			t.Fatalf("alias write did not join the canonical ledger: %d reservations", got)
		}
		if got := reloaded.OutstandingReservations(); got != 2 {
			t.Fatalf("outstanding after alias write = %d, want 2", got)
		}
		if got := durableReservedAmount(t, reloaded); got != 70 {
			t.Fatalf("reserved amount after alias write = %d, want 70", got)
		}
		if _, err := first.RequireOutstandingReservation(context.Background(), outcome.ReservationHash); err != nil {
			t.Fatalf("original reservation lost after alias write: %v", err)
		}
	}

	t.Run("symlinkAlias", func(t *testing.T) {
		link := filepath.Join(t.TempDir(), "scope-link")
		if err := os.Symlink(root, link); err != nil {
			// Registered platform limitation: Windows needs Developer Mode or an
			// elevated principal to create symlinks.
			t.Skipf("symlink alias unavailable on this platform: %v", err)
		}
		linked, err := OpenDurableCostLedger(durableConfig(link, &amount, 10, 1000))
		if err != nil {
			if !errors.Is(err, ErrInvalidCostLedger) {
				t.Fatalf("symlinked root must resolve to the same lock or fail closed, got %v", err)
			}
			return
		}
		defer func() { _ = linked.Close() }()
		linkedInfo, err := os.Stat(linked.lockPath)
		if err != nil {
			t.Fatalf("stat symlinked lock file: %v", err)
		}
		canonicalInfo, err := os.Stat(first.lockPath)
		if err != nil {
			t.Fatalf("stat canonical lock file: %v", err)
		}
		if !os.SameFile(linkedInfo, canonicalInfo) {
			t.Fatalf("symlinked root produced a second independent lock file")
		}
		if got := p39LockFileCount(t, root); got != 1 {
			t.Fatalf("symlink alias minted %d lock files, want 1", got)
		}
	})
}

// Matrix #3: a local edit of an already-confirmed prefix entry must be rejected
// by the in-lock re-read, never silently accepted.
func TestP39DurableReloadRejectsHistoryDrift(t *testing.T) {
	amount := int64(1000)
	root := t.TempDir()
	ledger := openDurable(t, root, &amount, 10, 1000)
	ctx := context.Background()

	if _, ack, err := ledger.Reserve(ctx, durableReservationInput("reservation-drift-one", "request-drift-one", "reserve-drift-one", 60, 1, 10)); err != nil {
		t.Fatalf("seed reserve: %v", err)
	} else if !ack.Valid() {
		t.Fatalf("seed reserve ack invalid")
	}
	path := ledger.Path()

	record := p39LoadRecord(t, path)
	allocation, ok := record.Ledger.Entries[0].(CostAllocationEntry)
	if !ok {
		t.Fatalf("first entry is %T, want CostAllocationEntry", record.Ledger.Entries[0])
	}
	if allocation.AuthorizationHash == costHashOne {
		t.Fatalf("test bug: drift replacement equals the original authorization hash")
	}
	allocation.AuthorizationHash = costHashOne
	record.Ledger.Entries[0] = allocation
	tampered := p39WriteRecord(t, path, record)

	if _, ack, err := ledger.Reserve(ctx, durableReservationInput("reservation-drift-two", "request-drift-two", "reserve-drift-two", 10, 1, 10)); err == nil {
		t.Fatalf("reload must reject a drifted confirmed prefix")
	} else {
		if !errors.Is(err, ErrInvalidCostLedger) {
			t.Fatalf("want %v, got %v", ErrInvalidCostLedger, err)
		}
		if !strings.Contains(err.Error(), "drift") {
			t.Fatalf("want an identity/prefix drift classification, got %v", err)
		}
		if ack.Valid() {
			t.Fatalf("rejected drift must not carry a valid ack")
		}
	}
	if _, err := ledger.RequireOutstandingReservation(ctx, costHashTwo); err == nil {
		t.Fatalf("read-only query must reject the drifted prefix as well")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read ledger after drift rejection: %v", err)
	}
	if !bytes.Equal(tampered, after) {
		t.Fatalf("rejected drift must not overwrite the failed input bytes")
	}
	if counts := durableEntryCounts(t, reloadDurableLedger(t, path)); counts["reservation"] != 1 {
		t.Fatalf("rejected drift committed entries: %+v", counts)
	}

	t.Run("identityDrift", func(t *testing.T) {
		idRoot := t.TempDir()
		idLedger := openDurable(t, idRoot, &amount, 10, 1000)
		idPath := idLedger.Path()
		if _, _, err := idLedger.Reserve(ctx, durableReservationInput("reservation-ident-one", "request-ident-one", "reserve-ident-one", 60, 1, 10)); err != nil {
			t.Fatalf("seed reserve: %v", err)
		}
		_ = idLedger.Close()

		idRecord := p39LoadRecord(t, idPath)
		idRecord.Ledger.LedgerID = "durable-cost-ledger-other"
		p39WriteRecord(t, idPath, idRecord)

		reopened, err := OpenDurableCostLedger(durableConfig(idRoot, &amount, 10, 1000))
		if err == nil {
			_ = reopened.Close()
			t.Fatalf("identity drift must not be adopted on reopen")
		}
		if !errors.Is(err, ErrInvalidCostLedger) {
			t.Fatalf("want %v, got %v", ErrInvalidCostLedger, err)
		}
		if !strings.Contains(err.Error(), "identity drift") {
			t.Fatalf("want an identity drift classification, got %v", err)
		}
	})
}

// Matrix #4: the read-only leg fails closed without creating files, and a
// committed-but-unconfirmed tail can never be promoted by a query.
func TestP39DurableReadCannotMintConfirmation(t *testing.T) {
	ctx := context.Background()
	reservationHash := "sha256:" + strings.Repeat("ab", 32)

	t.Run("missingLockFile", func(t *testing.T) {
		amount := int64(1000)
		root := t.TempDir()
		ledger := openDurable(t, root, &amount, 10, 1000)
		if err := os.Remove(ledger.lockPath); err != nil {
			t.Fatalf("remove lock file: %v", err)
		}
		before := p39FileCount(t, root)
		if _, err := ledger.RequireOutstandingReservation(ctx, reservationHash); err == nil {
			t.Fatalf("missing lock file must fail closed")
		} else if !errors.Is(err, ErrInvalidCostLedger) {
			t.Fatalf("want %v, got %v", ErrInvalidCostLedger, err)
		}
		if after := p39FileCount(t, root); after != before {
			t.Fatalf("read-only query created files: %d -> %d", before, after)
		}
		if _, statErr := os.Stat(ledger.lockPath); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("read-only query recreated the lock file")
		}
	})

	t.Run("missingLedger", func(t *testing.T) {
		amount := int64(1000)
		root := t.TempDir()
		ledger := openDurable(t, root, &amount, 10, 1000)
		path := ledger.Path()
		if err := os.Remove(path); err != nil {
			t.Fatalf("remove ledger: %v", err)
		}
		before := p39FileCount(t, root)
		if _, err := ledger.RequireOutstandingReservation(ctx, reservationHash); err == nil {
			t.Fatalf("missing ledger must fail closed")
		} else if !errors.Is(err, ErrInvalidCostLedger) {
			t.Fatalf("want %v, got %v", ErrInvalidCostLedger, err)
		}
		if after := p39FileCount(t, root); after != before {
			t.Fatalf("read-only query created files: %d -> %d", before, after)
		}
		if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("read-only query recreated the ledger: minted confirmation from nothing")
		}
	})

	t.Run("unconfirmedTail", func(t *testing.T) {
		amount := int64(1000)
		root := t.TempDir()
		writer := openDurable(t, root, &amount, 10, 1000)
		reader := openDurable(t, root, &amount, 10, 1000)

		firstOutcome, firstAck, err := writer.Reserve(ctx, durableReservationInput("reservation-read-one", "request-read-one", "reserve-read-one", 60, 1, 10))
		if err != nil {
			t.Fatalf("reserve one: %v", err)
		}
		if !firstAck.Valid() {
			t.Fatalf("reserve one ack invalid")
		}
		if _, err := reader.RequireOutstandingReservation(ctx, firstOutcome.ReservationHash); err == nil {
			t.Fatalf("reader opened before the write must not confirm it")
		}
		// Re-open so the reader's confirmed prefix covers reservation one only;
		// the tail written next must stay unconfirmed for this reader.
		reader = openDurable(t, root, &amount, 10, 1000)
		confirmedBefore := len(reader.confirmed)

		secondOutcome, secondAck, err := writer.Reserve(ctx, durableReservationInput("reservation-read-two", "request-read-two", "reserve-read-two", 10, 1, 10))
		if err != nil {
			t.Fatalf("reserve two: %v", err)
		}
		if !secondAck.Valid() {
			t.Fatalf("reserve two ack invalid")
		}
		beforeBytes, err := os.ReadFile(writer.Path())
		if err != nil {
			t.Fatalf("read ledger before query: %v", err)
		}
		if _, err := reader.RequireOutstandingReservation(ctx, secondOutcome.ReservationHash); err == nil {
			t.Fatalf("query promoted a committed-but-unconfirmed tail into a confirmation")
		} else if !errors.Is(err, ErrCostReservationNotFound) {
			t.Fatalf("want %v, got %v", ErrCostReservationNotFound, err)
		}
		if got := len(reader.confirmed); got != confirmedBefore {
			t.Fatalf("query changed the confirmed prefix: %d -> %d", confirmedBefore, got)
		}
		afterBytes, err := os.ReadFile(writer.Path())
		if err != nil {
			t.Fatalf("read ledger after query: %v", err)
		}
		if !bytes.Equal(beforeBytes, afterBytes) {
			t.Fatalf("read-only query wrote to the ledger")
		}
		if _, err := reader.RequireOutstandingReservation(ctx, firstOutcome.ReservationHash); err != nil {
			t.Fatalf("existing confirmed prefix must stay usable: %v", err)
		}
	})
}

// Matrix #5: every pre-commit stage failure rejects the ack and leaves the
// canonical bytes untouched; post-replace failures are classified as visible
// but unconfirmed; an inconclusive replace is classified as outcome unknown.
func TestP39DurablePrecommitFaultsRejectAck(t *testing.T) {
	ctx := context.Background()
	type faultCase struct {
		name        string
		stage       string
		tamper      bool
		want        error
		wantChanged bool
	}
	cases := []faultCase{
		{"write", StageWrite, false, ErrDurableNotCommitted, false},
		{"sync", StageSync, false, ErrDurableNotCommitted, false},
		{"close", StageClose, false, ErrDurableNotCommitted, false},
		{"rename", StageRename, false, ErrDurableNotCommitted, false},
		{"dirBarrier", StageDirBarrier, false, ErrDurableVisibleUnconfirmed, true},
		{"dirClose", StageDirClose, false, ErrDurableVisibleUnconfirmed, true},
		{"renameUncertain", StageRename, true, ErrDurableOutcomeUnknown, true},
	}
	for _, subject := range cases {
		subject := subject
		t.Run(subject.name, func(t *testing.T) {
			amount := int64(1000)
			root := t.TempDir()
			seed := openDurable(t, root, &amount, 10, 1000)
			path := seed.Path()
			if _, ack, err := seed.Reserve(ctx, durableReservationInput("reservation-seed", "request-seed", "reserve-seed", 60, 1, 10)); err != nil {
				t.Fatalf("seed reserve: %v", err)
			} else if !ack.Valid() {
				t.Fatalf("seed reserve ack invalid")
			}
			_ = seed.Close()

			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("read ledger before fault: %v", err)
			}
			injected := fmt.Errorf("injected %s fault", subject.stage)
			config := durableConfig(root, &amount, 10, 1000)
			config.CommitFault = func(stage string) error {
				if stage != subject.stage {
					return nil
				}
				if subject.tamper {
					if err := os.WriteFile(path, []byte("{not-a-ledger}"), 0o600); err != nil {
						return fmt.Errorf("tamper canonical ledger: %w", err)
					}
				}
				return injected
			}
			faulted, err := OpenDurableCostLedger(config)
			if err != nil {
				t.Fatalf("open faulted ledger: %v", err)
			}
			defer func() { _ = faulted.Close() }()

			_, ack, err := faulted.Reserve(ctx, durableReservationInput("reservation-fault", "request-fault", "reserve-fault", 10, 1, 10))
			if err == nil {
				t.Fatalf("injected %s fault must be reported", subject.stage)
			}
			if !errors.Is(err, subject.want) {
				t.Fatalf("want %v, got %v", subject.want, err)
			}
			if !errors.Is(err, injected) {
				t.Fatalf("injected cause must stay identifiable: %v", err)
			}
			if ack.Valid() {
				t.Fatalf("faulted commit must not carry a valid ack")
			}
			if subject.want == ErrDurableNotCommitted {
				if errors.Is(err, ErrDurableVisibleUnconfirmed) || errors.Is(err, ErrDurableOutcomeUnknown) {
					t.Fatalf("provably uncommitted failure must not be classified as visible/unknown: %v", err)
				}
			} else if !errors.Is(err, ErrDurableNotCommitted) {
				t.Fatalf("umbrella %v must hold: %v", ErrDurableNotCommitted, err)
			}
			if got, want := len(faulted.confirmed), 2; got != want {
				t.Fatalf("confirmed prefix advanced on a faulted commit: %d, want %d", got, want)
			}

			after, readErr := os.ReadFile(path)
			if !subject.wantChanged {
				if readErr != nil {
					t.Fatalf("uncommitted failure must leave the canonical ledger readable: %v", readErr)
				}
				if !bytes.Equal(before, after) {
					t.Fatalf("uncommitted failure mutated the canonical bytes")
				}
				recovered, err := OpenDurableCostLedger(durableConfig(root, &amount, 10, 1000))
				if err != nil {
					t.Fatalf("retry after a provably uncommitted failure must succeed: %v", err)
				}
				if _, ack, err := recovered.Reserve(ctx, durableReservationInput("reservation-retry", "request-retry", "reserve-retry", 10, 1, 10)); err != nil {
					t.Fatalf("retry reserve: %v", err)
				} else if !ack.Valid() {
					t.Fatalf("retry reserve ack invalid")
				}
				_ = recovered.Close()
			} else if readErr == nil && bytes.Equal(before, after) {
				t.Fatalf("visible-unconfirmed path must leave a new snapshot on disk")
			}
		})
	}

	t.Run("initialAllocation", func(t *testing.T) {
		amount := int64(1000)
		root := t.TempDir()
		injected := errors.New("injected initial allocation fault")
		config := durableConfig(root, &amount, 10, 1000)
		config.CommitFault = func(stage string) error {
			if stage == StageWrite {
				return injected
			}
			return nil
		}
		ledger, err := OpenDurableCostLedger(config)
		if err == nil {
			_ = ledger.Close()
			t.Fatalf("faulted initial allocation must fail")
		}
		if !errors.Is(err, ErrDurableNotCommitted) || !errors.Is(err, injected) {
			t.Fatalf("want injected pre-commit failure, got %v", err)
		}
		path := filepath.Join(root, durableLedgerFileName(config.Ledger))
		if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("faulted initial allocation left a canonical ledger behind")
		}
		if _, statErr := os.Stat(path + ".init"); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("provably uncommitted initial allocation must drop its init marker")
		}
		fresh := openDurable(t, root, &amount, 10, 1000)
		if counts := durableEntryCounts(t, reloadDurableLedger(t, fresh.Path())); counts["allocation"] != 1 {
			t.Fatalf("retried initial allocation is wrong: %+v", counts)
		}
		fresh.Close()
		again := openDurable(t, root, &amount, 10, 1000)
		if counts := durableEntryCounts(t, reloadDurableLedger(t, again.Path())); counts["allocation"] != 1 {
			t.Fatalf("re-open re-issued the initial allocation: %+v", counts)
		}
	})
}

// Matrix #6: an initialized scope whose canonical ledger vanished must fail
// closed instead of re-issuing the initial allocation.
func TestP39DurableLedgerMissingCannotReissueAllocation(t *testing.T) {
	amount := int64(1000)
	root := t.TempDir()
	ledger := openDurable(t, root, &amount, 10, 1000)
	path := ledger.Path()
	if _, ack, err := ledger.Reserve(context.Background(), durableReservationInput("reservation-missing-one", "request-missing-one", "reserve-missing-one", 60, 1, 10)); err != nil {
		t.Fatalf("seed reserve: %v", err)
	} else if !ack.Valid() {
		t.Fatalf("seed reserve ack invalid")
	}
	_ = ledger.Close()

	marker := path + ".init"
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("init marker must exist after initialization: %v", err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove canonical ledger: %v", err)
	}

	reopened, err := OpenDurableCostLedger(durableConfig(root, &amount, 10, 1000))
	if err == nil {
		_ = reopened.Close()
		t.Fatalf("missing ledger after initialization must fail closed")
	}
	if !errors.Is(err, ErrInvalidCostLedger) {
		t.Fatalf("want %v, got %v", ErrInvalidCostLedger, err)
	}
	if !strings.Contains(err.Error(), "refusing to re-issue") {
		t.Fatalf("want an explicit re-issue refusal, got %v", err)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("failed reopen recreated the canonical ledger (initial allocation re-issued)")
	}

	second, err := OpenDurableCostLedger(durableConfig(root, &amount, 10, 1000))
	if err == nil {
		_ = second.Close()
		t.Fatalf("repeated reopen must keep failing closed")
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("repeated reopen recreated the canonical ledger")
	}
}
