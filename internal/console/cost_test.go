package console

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/tickets"
)

const (
	costReportHashOne   = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	costReportHashTwo   = "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	costReportHashThree = "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	costReportHashFour  = "sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
)

func writeCostLedgerFixture(t *testing.T, root string, pricingMode string, amount *int64) string {
	t.Helper()
	ledger, err := tickets.NewCostLedger(tickets.CostLedgerConfig{
		LedgerID:       "cost-ledger-one",
		CreatedAt:      time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
		ScopeKind:      tickets.CostScopeRun,
		ScopeHash:      costReportHashOne,
		Currency:       "USD",
		PricingMode:    pricingMode,
		PricingVersion: "2026-09-08",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.AppendAllocation(tickets.AllocationInput{
		EntryID:           "allocation-one",
		OccurredAt:        time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
		AuthorizationHash: costReportHashTwo,
		Limits: tickets.CostLimits{
			AmountMicros: amount,
			ModelCalls:   3,
			Tokens:       500,
			WallClockMs:  60000,
			Attempts:     3,
		},
	}); err != nil {
		t.Fatal(err)
	}
	reserved := int64(400)
	reservation, err := ledger.Reserve(tickets.ReservationInput{
		EntryID:              "reservation-one",
		OccurredAt:           time.Date(2026, 9, 8, 12, 1, 0, 0, time.UTC),
		RunID:                "run-one",
		RequestID:            "request-one",
		IdempotencyKey:       "idem-one",
		AuthorizationHash:    costReportHashTwo,
		ReservedAmountMicros: &reserved,
		ReservedCalls:        1,
		ReservedTokens:       120,
		PricingEvidence:      []string{costReportHashThree},
	})
	if err != nil {
		t.Fatal(err)
	}
	charged := int64(260)
	if _, err := ledger.Settle(tickets.SettlementInput{
		EntryID:             "settlement-one",
		OccurredAt:          time.Date(2026, 9, 8, 12, 2, 0, 0, time.UTC),
		RunID:               "run-one",
		RequestID:           "request-one",
		IdempotencyKey:      "settle-one",
		ReservationHash:     reservation.ReservationHash,
		Status:              tickets.CostSettlementObserved,
		ChargedAmountMicros: &charged,
		ProviderEvidence:    []string{costReportHashFour},
	}); err != nil {
		t.Fatal(err)
	}
	record, err := ledger.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	canonical, err := evidence.EncodeCanonical(record)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "cost-ledger.json")
	if err := os.WriteFile(path, append(canonical, '\n'), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCostReportJSONAndText(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	amount := int64(1000)
	ledgerPath := writeCostLedgerFixture(t, root, tickets.CostPricingDocumented, &amount)

	code, stdout, stderr := runCLI(t, cli, "cost", "report", "--ledger", ledgerPath, "--json")
	if code != 0 {
		t.Fatalf("cost report code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result := decodeResponse(t, stdout)
	if !result.OK {
		t.Fatalf("cost report failed: %s", stdout)
	}
	var report CostReport
	if err := json.Unmarshal(result.Data, &report); err != nil {
		t.Fatal(err)
	}
	if report.TelemetryExported {
		t.Fatal("cost report must never export telemetry by default")
	}
	if report.OutstandingReservations != 0 || report.UnknownHoldReservations != 0 {
		t.Fatalf("unexpected hold counts: %+v", report)
	}
	if report.SettledAmountMicros == nil || *report.SettledAmountMicros != 260 {
		t.Fatalf("unexpected settled amount: %+v", report)
	}

	code, textOut, textErr := runCLI(t, cli, "cost", "report", "--ledger", ledgerPath)
	if code != 0 {
		t.Fatalf("cost report text code=%d stdout=%s stderr=%s", code, textOut, textErr)
	}
	if strings.Contains(textOut, "\x1b") || strings.Contains(textErr, "\x1b") {
		t.Fatalf("unexpected ANSI in text output stdout=%q stderr=%q", textOut, textErr)
	}
	if !strings.Contains(textOut, "telemetryExported: false") || !strings.Contains(textOut, "billing:") {
		t.Fatalf("text output missing key facts: %s", textOut)
	}
}

func TestCostReportMissingLedgerIsNonFatal(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "missing-cost-ledger.json")

	code, stdout, stderr := runCLI(t, cli, "cost", "report", "--ledger", ledgerPath, "--json")
	if code != 0 {
		t.Fatalf("cost report code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	var report CostReport
	if err := json.Unmarshal(decodeResponse(t, stdout).Data, &report); err != nil {
		t.Fatal(err)
	}
	if len(report.Warnings) == 0 {
		t.Fatal("expected warning for missing cost ledger")
	}
	if _, err := os.Stat(ledgerPath); !os.IsNotExist(err) {
		t.Fatal("report must not create missing ledger")
	}
}

func TestCostReportSubscriptionDisclaimer(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledger, err := tickets.NewCostLedger(tickets.CostLedgerConfig{
		LedgerID:       "cost-ledger-two",
		CreatedAt:      time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
		ScopeKind:      tickets.CostScopeRun,
		ScopeHash:      costReportHashOne,
		Currency:       "USD",
		PricingMode:    tickets.CostPricingSubscription,
		PricingVersion: "subscription-v1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.AppendAllocation(tickets.AllocationInput{
		EntryID:           "allocation-one",
		OccurredAt:        time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
		AuthorizationHash: costReportHashTwo,
		Limits:            tickets.CostLimits{AmountMicros: nil, ModelCalls: 5, Tokens: 1000, WallClockMs: 60000, Attempts: 3},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.Reserve(tickets.ReservationInput{
		EntryID:              "reservation-one",
		OccurredAt:           time.Date(2026, 9, 8, 12, 1, 0, 0, time.UTC),
		RunID:                "run-one",
		RequestID:            "request-one",
		IdempotencyKey:       "idem-one",
		AuthorizationHash:    costReportHashTwo,
		ReservedAmountMicros: nil,
		ReservedCalls:        1,
		ReservedTokens:       100,
		PricingEvidence:      []string{costReportHashThree},
	}); err != nil {
		t.Fatal(err)
	}
	record, err := ledger.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := evidence.EncodeCanonical(record)
	if err != nil {
		t.Fatal(err)
	}
	ledgerPath := filepath.Join(root, "cost-ledger-subscription.json")
	if err := os.WriteFile(ledgerPath, append(encoded, '\n'), 0644); err != nil {
		t.Fatal(err)
	}

	code, stdout, stderr := runCLI(t, cli, "cost", "report", "--ledger", ledgerPath, "--json")
	if code != 0 {
		t.Fatalf("cost report code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	var report CostReport
	if err := json.Unmarshal(decodeResponse(t, stdout).Data, &report); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(report.BillingDisclaimer, "does not guarantee exact monetary bills") {
		t.Fatalf("expected subscription disclaimer, got %q", report.BillingDisclaimer)
	}
}
