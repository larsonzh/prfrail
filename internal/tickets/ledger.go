package tickets

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	ticketLedgerDomain = "proofrail:ticket-ledger:1\n"
	timestampLayout    = "2006-01-02T15:04:05.000Z"
)

var (
	ErrInvalidPolicy          = errors.New("invalid ticket policy")
	ErrInvalidLedger          = errors.New("invalid ticket ledger")
	ErrDuplicateFailure       = errors.New("duplicate failure error hash")
	ErrInvalidLedgerOperation = errors.New("invalid ticket ledger operation")
)

type Policy struct {
	PolicyHash         string `json:"policyHash"`
	ReviewThreshold    int    `json:"reviewThreshold"`
	HardBlockThreshold int    `json:"hardBlockThreshold"`
}

type Entry struct {
	TicketID              string         `json:"ticketId"`
	Sequence              int            `json:"sequence"`
	RecordedAt            string         `json:"recordedAt"`
	Attempt               int            `json:"attempt"`
	Fingerprint           string         `json:"fingerprint"`
	Kind                  string         `json:"kind"`
	Actor                 evidence.Actor `json:"actor"`
	Evidence              []string       `json:"evidence"`
	ErrorHash             string         `json:"errorHash,omitempty"`
	AuthorizationHash     string         `json:"authorizationHash,omitempty"`
	ExpiresAt             string         `json:"expiresAt,omitempty"`
	MaxAdditionalAttempts int            `json:"maxAdditionalAttempts,omitempty"`
	ResolutionEvidence    []string       `json:"resolutionEvidence,omitempty"`
}

type LedgerBody struct {
	LedgerID           string           `json:"ledgerId"`
	RunID              string           `json:"runId"`
	Fingerprint        string           `json:"fingerprint"`
	FingerprintInput   FingerprintInput `json:"fingerprintInput"`
	Policy             Policy           `json:"policy"`
	Entries            []Entry          `json:"entries"`
	CurrentBudgetState string           `json:"currentBudgetState"`
}

type LedgerRecord struct {
	SchemaVersion    string     `json:"schemaVersion"`
	TicketLedger     LedgerBody `json:"ticketLedger"`
	TicketLedgerHash string     `json:"ticketLedgerHash"`
}

type Ledger struct {
	body              LedgerBody
	overrideUsage     map[int]int
	consumedByAttempt map[int]struct{}
}

func NewLedger(ledgerID, runID string, input FingerprintInput, policy Policy) (*Ledger, error) {
	if !evidence.ValidID(ledgerID) || !evidence.ValidID(runID) {
		return nil, fmt.Errorf("%w: invalid ledger identity", ErrInvalidLedger)
	}
	fingerprint, err := Fingerprint(input)
	if err != nil {
		return nil, err
	}
	if err := validatePolicy(policy); err != nil {
		return nil, err
	}
	return &Ledger{
		body: LedgerBody{
			LedgerID:           ledgerID,
			RunID:              runID,
			Fingerprint:        fingerprint,
			FingerprintInput:   input,
			Policy:             policy,
			Entries:            []Entry{},
			CurrentBudgetState: string(BudgetPendingReview),
		},
		overrideUsage:     map[int]int{},
		consumedByAttempt: map[int]struct{}{},
	}, nil
}

func (ledger *Ledger) AppendFailure(ticketID string, attempt int, errorHash string, actor evidence.Actor, evidenceHashes []string, recordedAt time.Time) error {
	if !evidence.ValidID(ticketID) || attempt < 1 || !evidence.ValidHash(errorHash) || !validActor(actor) {
		return fmt.Errorf("%w: invalid failure entry", ErrInvalidLedgerOperation)
	}
	if ledger.hasTicketID(ticketID) {
		return fmt.Errorf("%w: duplicate ticket ID", ErrInvalidLedgerOperation)
	}
	for _, entry := range ledger.body.Entries {
		if entry.Kind == "failure" && entry.ErrorHash == errorHash {
			return ErrDuplicateFailure
		}
	}
	entry := Entry{
		TicketID:    ticketID,
		Sequence:    len(ledger.body.Entries) + 1,
		RecordedAt:  recordedAt.UTC().Format(timestampLayout),
		Attempt:     attempt,
		ErrorHash:   errorHash,
		Fingerprint: ledger.body.Fingerprint,
		Kind:        "failure",
		Actor:       actor,
		Evidence:    uniqueHashes(evidenceHashes),
	}
	ledger.body.Entries = append(ledger.body.Entries, entry)
	return nil
}

func (ledger *Ledger) AppendOverride(ticketID string, attempt int, actor evidence.Actor, evidenceHashes []string, authorizationHash string, expiresAt time.Time, maxAdditionalAttempts int, recordedAt time.Time) error {
	if !evidence.ValidID(ticketID) || attempt < 1 || !evidence.ValidHash(authorizationHash) || maxAdditionalAttempts < 1 {
		return fmt.Errorf("%w: invalid override entry", ErrInvalidLedgerOperation)
	}
	if ledger.hasTicketID(ticketID) {
		return fmt.Errorf("%w: duplicate ticket ID", ErrInvalidLedgerOperation)
	}
	if actor.Type != "operator" && actor.Type != "policy" {
		return fmt.Errorf("%w: invalid override actor", ErrInvalidLedgerOperation)
	}
	if !evidence.ValidID(actor.ID) || !expiresAt.After(recordedAt.UTC()) {
		return fmt.Errorf("%w: invalid override timing", ErrInvalidLedgerOperation)
	}
	entry := Entry{
		TicketID:              ticketID,
		Sequence:              len(ledger.body.Entries) + 1,
		RecordedAt:            recordedAt.UTC().Format(timestampLayout),
		Attempt:               attempt,
		Fingerprint:           ledger.body.Fingerprint,
		Kind:                  "override-granted",
		Actor:                 actor,
		Evidence:              uniqueHashes(evidenceHashes),
		AuthorizationHash:     authorizationHash,
		ExpiresAt:             expiresAt.UTC().Format(timestampLayout),
		MaxAdditionalAttempts: maxAdditionalAttempts,
	}
	ledger.body.Entries = append(ledger.body.Entries, entry)
	return nil
}

func (ledger *Ledger) AppendResolved(ticketID string, attempt int, actor evidence.Actor, evidenceHashes, resolutionEvidence []string, recordedAt time.Time) error {
	if !evidence.ValidID(ticketID) || attempt < 1 || !validActor(actor) {
		return fmt.Errorf("%w: invalid resolved entry", ErrInvalidLedgerOperation)
	}
	if ledger.hasTicketID(ticketID) {
		return fmt.Errorf("%w: duplicate ticket ID", ErrInvalidLedgerOperation)
	}
	if len(resolutionEvidence) == 0 {
		return fmt.Errorf("%w: missing resolution evidence", ErrInvalidLedgerOperation)
	}
	entry := Entry{
		TicketID:           ticketID,
		Sequence:           len(ledger.body.Entries) + 1,
		RecordedAt:         recordedAt.UTC().Format(timestampLayout),
		Attempt:            attempt,
		Fingerprint:        ledger.body.Fingerprint,
		Kind:               "resolved",
		Actor:              actor,
		Evidence:           uniqueHashes(evidenceHashes),
		ResolutionEvidence: uniqueHashes(resolutionEvidence),
	}
	ledger.body.Entries = append(ledger.body.Entries, entry)
	return nil
}

func (ledger *Ledger) Snapshot(now time.Time) (LedgerRecord, error) {
	if len(ledger.body.Entries) == 0 {
		return LedgerRecord{}, fmt.Errorf("%w: empty entries", ErrInvalidLedger)
	}
	body := ledger.body
	body.CurrentBudgetState = string(ledger.BudgetState(now))
	hash, err := digest(ticketLedgerDomain, body)
	if err != nil {
		return LedgerRecord{}, err
	}
	return LedgerRecord{SchemaVersion: evidence.SchemaVersion, TicketLedger: body, TicketLedgerHash: hash}, nil
}

func (ledger *Ledger) Fingerprint() string { return ledger.body.Fingerprint }
func (ledger *Ledger) LedgerID() string    { return ledger.body.LedgerID }

func (ledger *Ledger) hasTicketID(ticketID string) bool {
	for _, entry := range ledger.body.Entries {
		if entry.TicketID == ticketID {
			return true
		}
	}
	return false
}

func validatePolicy(policy Policy) error {
	if !evidence.ValidHash(policy.PolicyHash) || policy.ReviewThreshold < 1 || policy.HardBlockThreshold < 2 || policy.HardBlockThreshold <= policy.ReviewThreshold {
		return ErrInvalidPolicy
	}
	return nil
}

func validActor(actor evidence.Actor) bool {
	if !evidence.ValidID(actor.ID) {
		return false
	}
	switch actor.Type {
	case "system", "operator", "agent", "policy":
		return true
	default:
		return false
	}
}

func digest(domain string, body any) (string, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	canonical, err := evidence.Canonicalize(payload)
	if err != nil {
		return "", err
	}
	return evidence.Digest(domain, canonical), nil
}

func uniqueHashes(values []string) []string {
	result := make([]string, 0, len(values))
	if len(values) < 2 {
		for _, value := range values {
			if evidence.ValidHash(value) {
				result = append(result, value)
			}
		}
		return result
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !evidence.ValidHash(value) {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
