package tickets

import (
	"errors"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	hashOne   = "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	hashTwo   = "sha256:2222222222222222222222222222222222222222222222222222222222222222"
	hashThree = "sha256:3333333333333333333333333333333333333333333333333333333333333333"
	hashFour  = "sha256:4444444444444444444444444444444444444444444444444444444444444444"
)

func newLedger(t *testing.T, review, hard int) *Ledger {
	t.Helper()
	ledger, err := NewLedger(
		"ledger-one",
		"run-one",
		FingerprintInput{Code: "execution.exit-nonzero", Subject: Subject{Kind: "task", ID: "task-one"}, FailurePoint: "validation"},
		Policy{PolicyHash: hashOne, ReviewThreshold: review, HardBlockThreshold: hard},
	)
	if err != nil {
		t.Fatal(err)
	}
	return ledger
}

func now() time.Time { return time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC) }

func TestLedgerNoBudgetResetAndHardBlock(t *testing.T) {
	ledger := newLedger(t, 1, 2)
	if err := ledger.AppendFailure("ticket-one", 1, hashTwo, evidence.Actor{Type: "system", ID: "proofrail"}, nil, now()); err != nil {
		t.Fatal(err)
	}
	if err := ledger.AppendResolved("ticket-two", 1, evidence.Actor{Type: "operator", ID: "operator-one"}, nil, []string{hashThree}, now().Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	record, err := ledger.Snapshot(now().Add(2 * time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if record.TicketLedger.CurrentBudgetState != string(BudgetPendingReview) {
		t.Fatalf("state after resolve: %s", record.TicketLedger.CurrentBudgetState)
	}
	if err := ledger.AppendFailure("ticket-three", 2, hashThree, evidence.Actor{Type: "system", ID: "proofrail"}, nil, now().Add(3*time.Minute)); err != nil {
		t.Fatal(err)
	}
	record, err = ledger.Snapshot(now().Add(4 * time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if record.TicketLedger.CurrentBudgetState != string(BudgetHardBlock) {
		t.Fatalf("state after second failure: %s", record.TicketLedger.CurrentBudgetState)
	}
}

func TestLedgerOverrideAttemptExhaustion(t *testing.T) {
	ledger := newLedger(t, 1, 3)
	if err := ledger.AppendFailure("ticket-one", 1, hashTwo, evidence.Actor{Type: "system", ID: "proofrail"}, nil, now()); err != nil {
		t.Fatal(err)
	}
	if err := ledger.AppendOverride("ticket-override", 1, evidence.Actor{Type: "operator", ID: "operator-one"}, nil, hashFour, now().Add(10*time.Minute), 2, now().Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := ledger.ReserveAttempt(2, now().Add(2*time.Minute), false, false); err != nil {
		t.Fatal(err)
	}
	if err := ledger.ReserveAttempt(2, now().Add(2*time.Minute), false, false); !errors.Is(err, ErrBudgetAttemptExhausted) {
		t.Fatalf("duplicate attempt error: %v", err)
	}
	if err := ledger.ReserveAttempt(3, now().Add(2*time.Minute), false, false); err != nil {
		t.Fatal(err)
	}
	if err := ledger.ReserveAttempt(4, now().Add(2*time.Minute), false, false); !errors.Is(err, ErrBudgetOverrideExhausted) {
		t.Fatalf("override exhaustion error: %v", err)
	}
}

func TestLedgerWallClockAndCostBudgetExhaustion(t *testing.T) {
	ledger := newLedger(t, 2, 4)
	if err := ledger.AppendFailure("ticket-one", 1, hashTwo, evidence.Actor{Type: "system", ID: "proofrail"}, nil, now()); err != nil {
		t.Fatal(err)
	}
	if err := ledger.ReserveAttempt(2, now().Add(time.Minute), true, false); !errors.Is(err, ErrBudgetWallClockExhausted) {
		t.Fatalf("wall-clock exhaustion error: %v", err)
	}
	if err := ledger.ReserveAttempt(2, now().Add(time.Minute), false, true); !errors.Is(err, ErrBudgetCostExhausted) {
		t.Fatalf("cost exhaustion error: %v", err)
	}
}

func TestLedgerRejectsDuplicateFailureHash(t *testing.T) {
	ledger := newLedger(t, 1, 3)
	if err := ledger.AppendFailure("ticket-one", 1, hashTwo, evidence.Actor{Type: "system", ID: "proofrail"}, nil, now()); err != nil {
		t.Fatal(err)
	}
	if err := ledger.AppendFailure("ticket-two", 2, hashTwo, evidence.Actor{Type: "system", ID: "proofrail"}, nil, now().Add(time.Minute)); !errors.Is(err, ErrDuplicateFailure) {
		t.Fatalf("duplicate failure error: %v", err)
	}
}

func TestFingerprintIsStable(t *testing.T) {
	input := FingerprintInput{Code: "execution.exit-nonzero", Subject: Subject{Kind: "task", ID: "task-one"}, FailurePoint: "validation"}
	left, err := Fingerprint(input)
	if err != nil {
		t.Fatal(err)
	}
	right, err := Fingerprint(input)
	if err != nil {
		t.Fatal(err)
	}
	if left != right {
		t.Fatalf("fingerprint mismatch: %s vs %s", left, right)
	}
}

func TestLedgerRejectsDuplicateTicketID(t *testing.T) {
	ledger := newLedger(t, 1, 3)
	if err := ledger.AppendFailure("ticket-one", 1, hashTwo, evidence.Actor{Type: "system", ID: "proofrail"}, nil, now()); err != nil {
		t.Fatal(err)
	}
	if err := ledger.AppendResolved("ticket-one", 1, evidence.Actor{Type: "operator", ID: "operator-one"}, nil, []string{hashThree}, now().Add(time.Minute)); err == nil {
		t.Fatal("duplicate ticket ID must be rejected")
	}
}

func TestLedgerOverrideExhaustionFallsBackToPendingReview(t *testing.T) {
	ledger := newLedger(t, 1, 3)
	if err := ledger.AppendFailure("ticket-one", 1, hashTwo, evidence.Actor{Type: "system", ID: "proofrail"}, nil, now()); err != nil {
		t.Fatal(err)
	}
	if err := ledger.AppendOverride("ticket-override", 1, evidence.Actor{Type: "operator", ID: "operator-one"}, nil, hashFour, now().Add(10*time.Minute), 1, now().Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	if err := ledger.ReserveAttempt(2, now().Add(2*time.Minute), false, false); err != nil {
		t.Fatal(err)
	}
	record, err := ledger.Snapshot(now().Add(3 * time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if record.TicketLedger.CurrentBudgetState != string(BudgetPendingReview) {
		t.Fatalf("exhausted override must fall back: %s", record.TicketLedger.CurrentBudgetState)
	}
}
