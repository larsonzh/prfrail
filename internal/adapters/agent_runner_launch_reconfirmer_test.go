package adapters

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/tickets"
)

func reconfirmerFixture(t *testing.T) (AgentRunnerLaunchReconfirmer, chain.StepRequest) {
	t.Helper()
	admission, request := compositeAdmissionFixture(t)
	return AgentRunnerLaunchReconfirmer{
		RequestRecord:       admission.RequestRecord,
		AuthorizationLedger: admission.AuthorizationLedger,
		CostLedger:          admission.CostLedger,
		Clock:               admission.Clock,
	}, request
}

func TestAgentRunnerLaunchReconfirmerAcceptsBoundFreshFacts(t *testing.T) {
	reconfirmer, _ := reconfirmerFixture(t)
	if err := reconfirmer.ReconfirmLaunchAuthorization(context.Background()); err != nil {
		t.Fatalf("expected reconfirmation to pass, got %v", err)
	}
	if err := reconfirmer.ReconfirmLaunchAuthorization(context.Background()); err != nil {
		t.Fatalf("reconfirmation must be stateless, got %v", err)
	}
}

func TestAgentRunnerLaunchReconfirmerBlocksRevokedGrant(t *testing.T) {
	reconfirmer, _ := reconfirmerFixture(t)
	grant := reconfirmer.AuthorizationLedger.Grants[0]
	revocation, err := chain.NewRevocationAuthorizationRecord(chain.AuthorizationRevocation{
		RecordID:          "revocation-one",
		Kind:              "revocation",
		IssuedAt:          "2026-09-11T08:00:00.000Z",
		IssuedBy:          evidence.Actor{Type: "operator", ID: "operator-one"},
		AuthorizationID:   grant.Grant.AuthorizationID,
		AuthorizationHash: grant.RecordHash,
		ReasonCode:        "operator-request",
		StopDisposition:   "not-required",
	}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	reconfirmer.AuthorizationLedger.Revocations = []chain.AuthorizationRecord{revocation}
	err = reconfirmer.ReconfirmLaunchAuthorization(context.Background())
	if !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) || !strings.Contains(err.Error(), "revoked") {
		t.Fatalf("expected revoked reconfirmation failure, got %v", err)
	}
}

func TestAgentRunnerLaunchReconfirmerBlocksExpiredGrant(t *testing.T) {
	reconfirmer, _ := reconfirmerFixture(t)
	reconfirmer.Clock = func() time.Time {
		return time.Date(2026, 9, 11, 9, 0, 0, 0, time.UTC)
	}
	err := reconfirmer.ReconfirmLaunchAuthorization(context.Background())
	if !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expected expired reconfirmation failure, got %v", err)
	}
}

func TestAgentRunnerLaunchReconfirmerBlocksSettledReservation(t *testing.T) {
	reconfirmer, request := reconfirmerFixture(t)
	amount := int64(1000)
	calls := 1
	tokens := 1000
	if _, err := reconfirmer.CostLedger.Settle(tickets.SettlementInput{
		EntryID:             "settlement-one",
		OccurredAt:          time.Date(2026, 9, 11, 8, 0, 0, 0, time.UTC),
		RunID:               request.RunID,
		RequestID:           request.AgentRunnerFacts.RequestID,
		IdempotencyKey:      "settle-request-one",
		ReservationHash:     request.AgentRunnerFacts.BudgetHash,
		Status:              tickets.CostSettlementObserved,
		ChargedAmountMicros: &amount,
		ObservedCalls:       &calls,
		ObservedTokens:      &tokens,
		ProviderEvidence:    []string{queueHashFour},
	}); err != nil {
		t.Fatal(err)
	}
	err := reconfirmer.ReconfirmLaunchAuthorization(context.Background())
	if !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) || !errors.Is(err, tickets.ErrCostReservationSettled) {
		t.Fatalf("expected settled reservation to block reconfirmation, got %v", err)
	}
}

func TestAgentRunnerLaunchReconfirmerBlocksBindingMismatch(t *testing.T) {
	reconfirmer, _ := reconfirmerFixture(t)
	foreignAuthority := reconfirmer.RequestRecord.Request
	foreignAuthority.AuthorizationHash = queueHashOne
	rebuilt, err := NewAgentRunnerRequestRecord(foreignAuthority)
	if err != nil {
		t.Fatal(err)
	}
	reconfirmer.RequestRecord = rebuilt
	err = reconfirmer.ReconfirmLaunchAuthorization(context.Background())
	if !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) || !strings.Contains(err.Error(), "matched 0 grants") {
		t.Fatalf("expected authority mismatch failure, got %v", err)
	}

	reconfirmer, _ = reconfirmerFixture(t)
	foreignBudget := reconfirmer.RequestRecord.Request
	foreignBudget.BudgetHash = queueHashOne
	rebuilt, err = NewAgentRunnerRequestRecord(foreignBudget)
	if err != nil {
		t.Fatal(err)
	}
	reconfirmer.RequestRecord = rebuilt
	if err := reconfirmer.ReconfirmLaunchAuthorization(context.Background()); !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) {
		t.Fatalf("expected budget mismatch failure, got %v", err)
	}
}

func TestAgentRunnerLaunchReconfirmerFailsClosedOnInvalidInputs(t *testing.T) {
	reconfirmer, _ := reconfirmerFixture(t)
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := reconfirmer.ReconfirmLaunchAuthorization(cancelled); !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) {
		t.Fatalf("expected cancelled context failure, got %v", err)
	}

	noClock := reconfirmer
	noClock.Clock = nil
	if err := noClock.ReconfirmLaunchAuthorization(context.Background()); !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) {
		t.Fatalf("expected missing clock failure, got %v", err)
	}
	zeroClock := reconfirmer
	zeroClock.Clock = func() time.Time { return time.Time{} }
	if err := zeroClock.ReconfirmLaunchAuthorization(context.Background()); !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) {
		t.Fatalf("expected zero clock failure, got %v", err)
	}
	noBudget := reconfirmer
	noBudget.CostLedger = nil
	if err := noBudget.ReconfirmLaunchAuthorization(context.Background()); !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) {
		t.Fatalf("expected missing cost ledger failure, got %v", err)
	}
	var zero AgentRunnerLaunchReconfirmer
	if err := zero.ReconfirmLaunchAuthorization(context.Background()); !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) {
		t.Fatalf("expected zero-value reconfirmer failure, got %v", err)
	}
}
