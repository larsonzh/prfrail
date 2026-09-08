package tickets

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	costHashOne   = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	costHashTwo   = "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	costHashThree = "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	costHashFour  = "sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
)

func costNow() time.Time {
	return time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
}

func mustCostLedger(t *testing.T, scopeKind string, pricingMode string, amountMicros *int64, callCap, tokenCap int) *CostLedger {
	t.Helper()
	ledger, err := NewCostLedger(CostLedgerConfig{
		LedgerID:       "cost-ledger-one",
		CreatedAt:      costNow(),
		ScopeKind:      scopeKind,
		ScopeHash:      costHashOne,
		Currency:       "USD",
		PricingMode:    pricingMode,
		PricingVersion: "2026-09-08",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.AppendAllocation(AllocationInput{
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
	}); err != nil {
		t.Fatal(err)
	}
	return ledger
}

func TestCostLedgerSharedCapAcrossRuns(t *testing.T) {
	amount := int64(1000)
	ledger := mustCostLedger(t, CostScopeOwner, CostPricingDocumented, &amount, 2, 300)

	firstAmount := int64(600)
	if _, err := ledger.Reserve(ReservationInput{
		EntryID:              "reservation-one",
		OccurredAt:           costNow().Add(time.Minute),
		RunID:                "run-one",
		RequestID:            "request-one",
		IdempotencyKey:       "idem-one",
		AuthorizationHash:    costHashTwo,
		ReservedAmountMicros: &firstAmount,
		ReservedCalls:        1,
		ReservedTokens:       100,
		PricingEvidence:      []string{costHashThree},
	}); err != nil {
		t.Fatal(err)
	}

	secondAmount := int64(500)
	if _, err := ledger.Reserve(ReservationInput{
		EntryID:              "reservation-two",
		OccurredAt:           costNow().Add(2 * time.Minute),
		RunID:                "run-two",
		RequestID:            "request-two",
		IdempotencyKey:       "idem-two",
		AuthorizationHash:    costHashTwo,
		ReservedAmountMicros: &secondAmount,
		ReservedCalls:        1,
		ReservedTokens:       120,
		PricingEvidence:      []string{costHashThree},
	}); !errors.Is(err, ErrCostSharedCapExceeded) {
		t.Fatalf("expected shared cap error, got %v", err)
	}

	fitAmount := int64(400)
	if _, err := ledger.Reserve(ReservationInput{
		EntryID:              "reservation-three",
		OccurredAt:           costNow().Add(3 * time.Minute),
		RunID:                "run-two",
		RequestID:            "request-three",
		IdempotencyKey:       "idem-three",
		AuthorizationHash:    costHashTwo,
		ReservedAmountMicros: &fitAmount,
		ReservedCalls:        1,
		ReservedTokens:       120,
		PricingEvidence:      []string{costHashThree},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestCostLedgerUnknownHoldSurvivesRestart(t *testing.T) {
	amount := int64(1000)
	ledger := mustCostLedger(t, CostScopeProject, CostPricingDocumented, &amount, 3, 300)

	reserved := int64(400)
	reservation, err := ledger.Reserve(ReservationInput{
		EntryID:              "reservation-one",
		OccurredAt:           costNow().Add(time.Minute),
		RunID:                "run-one",
		RequestID:            "request-one",
		IdempotencyKey:       "idem-one",
		AuthorizationHash:    costHashTwo,
		ReservedAmountMicros: &reserved,
		ReservedCalls:        1,
		ReservedTokens:       80,
		PricingEvidence:      []string{costHashThree},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.Settle(SettlementInput{
		EntryID:          "settlement-one",
		OccurredAt:       costNow().Add(2 * time.Minute),
		RunID:            "run-one",
		RequestID:        "request-one",
		IdempotencyKey:   "settle-one",
		ReservationHash:  reservation.ReservationHash,
		Status:           CostSettlementUnknown,
		ProviderEvidence: []string{costHashFour},
	}); err != nil {
		t.Fatal(err)
	}

	record, err := ledger.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadCostLedger(record)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.OutstandingReservations() != 0 || loaded.UnknownHoldReservations() != 1 {
		t.Fatalf("unexpected hold counts: outstanding=%d holds=%d", loaded.OutstandingReservations(), loaded.UnknownHoldReservations())
	}
	summary := loaded.Summary()
	if summary.UnknownReservedAmountMicros == nil || *summary.UnknownReservedAmountMicros != 400 {
		t.Fatalf("unknown hold drift: %+v", summary)
	}
	if summary.SettledAmountMicros == nil || *summary.SettledAmountMicros != 0 {
		t.Fatalf("settled amount drift: %+v", summary)
	}

	blocked := int64(700)
	if _, err := loaded.Reserve(ReservationInput{
		EntryID:              "reservation-two",
		OccurredAt:           costNow().Add(3 * time.Minute),
		RunID:                "run-two",
		RequestID:            "request-two",
		IdempotencyKey:       "idem-two",
		AuthorizationHash:    costHashTwo,
		ReservedAmountMicros: &blocked,
		ReservedCalls:        1,
		ReservedTokens:       80,
		PricingEvidence:      []string{costHashThree},
	}); !errors.Is(err, ErrCostSharedCapExceeded) {
		t.Fatalf("expected shared cap block under unknown hold, got %v", err)
	}
}

func TestCostLedgerSettlementDeduplicates(t *testing.T) {
	amount := int64(1200)
	ledger := mustCostLedger(t, CostScopeRun, CostPricingDocumented, &amount, 3, 400)

	reserved := int64(500)
	reservation, err := ledger.Reserve(ReservationInput{
		EntryID:              "reservation-one",
		OccurredAt:           costNow().Add(time.Minute),
		RunID:                "run-one",
		RequestID:            "request-one",
		IdempotencyKey:       "idem-one",
		AuthorizationHash:    costHashTwo,
		ReservedAmountMicros: &reserved,
		ReservedCalls:        1,
		ReservedTokens:       100,
		PricingEvidence:      []string{costHashThree},
	})
	if err != nil {
		t.Fatal(err)
	}
	charged := int64(320)
	first, err := ledger.Settle(SettlementInput{
		EntryID:             "settlement-one",
		OccurredAt:          costNow().Add(2 * time.Minute),
		RunID:               "run-one",
		RequestID:           "request-one",
		IdempotencyKey:      "settle-one",
		ReservationHash:     reservation.ReservationHash,
		Status:              CostSettlementObserved,
		ChargedAmountMicros: &charged,
		ProviderEvidence:    []string{costHashFour},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !first.Appended {
		t.Fatal("first settlement must append")
	}

	second, err := ledger.Settle(SettlementInput{
		EntryID:             "settlement-two",
		OccurredAt:          costNow().Add(3 * time.Minute),
		RunID:               "run-one",
		RequestID:           "request-one",
		IdempotencyKey:      "settle-one",
		ReservationHash:     reservation.ReservationHash,
		Status:              CostSettlementObserved,
		ChargedAmountMicros: &charged,
		ProviderEvidence:    []string{costHashFour},
	})
	if err != nil {
		t.Fatal(err)
	}
	if second.Appended || second.EntryHash != first.EntryHash {
		t.Fatalf("duplicate settlement must be idempotent: %+v vs %+v", first, second)
	}

	otherAmount := int64(300)
	if _, err := ledger.Settle(SettlementInput{
		EntryID:             "settlement-three",
		OccurredAt:          costNow().Add(4 * time.Minute),
		RunID:               "run-one",
		RequestID:           "request-one",
		IdempotencyKey:      "settle-two",
		ReservationHash:     reservation.ReservationHash,
		Status:              CostSettlementObserved,
		ChargedAmountMicros: &otherAmount,
		ProviderEvidence:    []string{costHashFour},
	}); !errors.Is(err, ErrCostReservationSettled) {
		t.Fatalf("expected reservation settled error, got %v", err)
	}

	changed := int64(321)
	if _, err := ledger.Settle(SettlementInput{
		EntryID:             "settlement-four",
		OccurredAt:          costNow().Add(5 * time.Minute),
		RunID:               "run-one",
		RequestID:           "request-one",
		IdempotencyKey:      "settle-one",
		ReservationHash:     reservation.ReservationHash,
		Status:              CostSettlementObserved,
		ChargedAmountMicros: &changed,
		ProviderEvidence:    []string{costHashFour},
	}); !errors.Is(err, ErrCostSettlementConflict) {
		t.Fatalf("expected settlement conflict, got %v", err)
	}
}

func TestCostLedgerCrashPointReservationLifecycle(t *testing.T) {
	amount := int64(900)
	ledger := mustCostLedger(t, CostScopeRun, CostPricingDocumented, &amount, 2, 200)

	beforeReserve, err := ledger.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if len(beforeReserve.Ledger.Entries) != 1 {
		t.Fatalf("expected allocation-only snapshot before reservation, got %d entries", len(beforeReserve.Ledger.Entries))
	}

	reserved := int64(300)
	reservation, err := ledger.Reserve(ReservationInput{
		EntryID:              "reservation-one",
		OccurredAt:           costNow().Add(time.Minute),
		RunID:                "run-one",
		RequestID:            "request-one",
		IdempotencyKey:       "idem-one",
		AuthorizationHash:    costHashTwo,
		ReservedAmountMicros: &reserved,
		ReservedCalls:        1,
		ReservedTokens:       60,
		PricingEvidence:      []string{costHashThree},
	})
	if err != nil {
		t.Fatal(err)
	}
	afterReserve, err := ledger.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if len(afterReserve.Ledger.Entries) != 2 {
		t.Fatalf("expected reservation snapshot after reserve, got %d entries", len(afterReserve.Ledger.Entries))
	}
	if ledger.OutstandingReservations() != 1 {
		t.Fatalf("reservation must remain outstanding until settlement, got %d", ledger.OutstandingReservations())
	}

	charged := int64(260)
	if _, err := ledger.Settle(SettlementInput{
		EntryID:             "settlement-one",
		OccurredAt:          costNow().Add(2 * time.Minute),
		RunID:               "run-one",
		RequestID:           "request-one",
		IdempotencyKey:      "settle-one",
		ReservationHash:     reservation.ReservationHash,
		Status:              CostSettlementObserved,
		ChargedAmountMicros: &charged,
		ProviderEvidence:    []string{costHashFour},
	}); err != nil {
		t.Fatal(err)
	}
	if ledger.OutstandingReservations() != 0 || ledger.UnknownHoldReservations() != 0 {
		t.Fatalf("holds must clear after observed settlement: outstanding=%d holds=%d", ledger.OutstandingReservations(), ledger.UnknownHoldReservations())
	}
}

func TestDecodeCostLedgerRecordRejectsSummaryDrift(t *testing.T) {
	amount := int64(1000)
	ledger := mustCostLedger(t, CostScopeRun, CostPricingDocumented, &amount, 2, 200)
	reserved := int64(300)
	if _, err := ledger.Reserve(ReservationInput{
		EntryID:              "reservation-one",
		OccurredAt:           costNow().Add(time.Minute),
		RunID:                "run-one",
		RequestID:            "request-one",
		IdempotencyKey:       "idem-one",
		AuthorizationHash:    costHashTwo,
		ReservedAmountMicros: &reserved,
		ReservedCalls:        1,
		ReservedTokens:       60,
		PricingEvidence:      []string{costHashThree},
	}); err != nil {
		t.Fatal(err)
	}
	record, err := ledger.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	record.Ledger.Summary.UnknownReservedAmountMicros = int64Pointer(0)
	record.LedgerHash, err = digest(costLedgerDomain, record.Ledger)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := evidence.EncodeCanonical(record)
	if err != nil {
		t.Fatal(err)
	}
	_, err = DecodeCostLedgerRecord(encoded)
	if err == nil || !strings.Contains(err.Error(), "summary mismatch") {
		t.Fatalf("expected summary mismatch, got %v", err)
	}
}

func TestCostBudgetExhaustionHelpers(t *testing.T) {
	capMicros := int64(200)
	summary := CostSummary{
		SettledAmountMicros:         int64Pointer(120),
		UnknownReservedAmountMicros: int64Pointer(50),
		ReservedTokens:              300,
		SettledTokens:               150,
	}
	if CostBudgetExhausted(&capMicros, summary) {
		t.Fatal("committed 170 must not exhaust 200")
	}
	capMicros = 170
	if !CostBudgetExhausted(&capMicros, summary) {
		t.Fatal("committed 170 must exhaust 170")
	}

	unknown := summary
	unknown.UnknownReservedAmountMicros = nil
	if !CostBudgetExhausted(&capMicros, unknown) {
		t.Fatal("unknown committed amount must fail closed as exhausted")
	}
	if CostTokenBudgetExhausted(400, summary) {
		t.Fatal("tokens below cap should not exhaust")
	}
	if !CostTokenBudgetExhausted(300, summary) {
		t.Fatal("tokens at cap should exhaust")
	}
}
