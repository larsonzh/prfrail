package repair

import (
	"errors"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/tickets"
)

func TestRepairRequiresBudgetReservation(t *testing.T) {
	ledger, err := tickets.NewLedger(
		"ledger-one",
		"run-one",
		tickets.FingerprintInput{Code: "execution.exit-nonzero", Subject: tickets.Subject{Kind: "task", ID: "task-one"}, FailurePoint: "validation"},
		tickets.Policy{PolicyHash: hashOne, ReviewThreshold: 1, HardBlockThreshold: 2},
	)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	if err := ledger.AppendFailure("ticket-one", 1, hashTwo, evidence.Actor{Type: "system", ID: "proofrail"}, nil, now); err != nil {
		t.Fatal(err)
	}
	if err := ledger.ReserveAttempt(2, now.Add(time.Minute), false, false); !errors.Is(err, tickets.ErrBudgetNeedsOverride) {
		t.Fatalf("expected manual review gate: %v", err)
	}
	if err := ledger.AppendOverride("ticket-override", 1, evidence.Actor{Type: "operator", ID: "operator-one"}, nil, hashThree, now.Add(10*time.Minute), 1, now.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	before, err := ledger.Snapshot(now.Add(3 * time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if before.TicketLedger.CurrentBudgetState != string(tickets.BudgetOverrideWindow) {
		t.Fatalf("unexpected budget state: %s", before.TicketLedger.CurrentBudgetState)
	}
	if err := ledger.ReserveAttempt(2, now.Add(3*time.Minute), false, false); err != nil {
		t.Fatalf("override reservation failed: %v", err)
	}
	record, err := ledger.Snapshot(now.Add(3 * time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if record.TicketLedger.CurrentBudgetState != string(tickets.BudgetPendingReview) {
		t.Fatalf("exhausted override must fall back: %s", record.TicketLedger.CurrentBudgetState)
	}
}
