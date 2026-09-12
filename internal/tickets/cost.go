package tickets

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	costLedgerDomain          = "proofrail:cost-ledger:1\n"
	costLimitsHashDomain      = "proofrail:cost-limits:1\n"
	costReservationHashDomain = "proofrail:cost-reservation:1\n"
	costSettlementHashDomain  = "proofrail:cost-settlement:1\n"
)

const (
	CostScopeOwner   = "owner"
	CostScopeProject = "project"
	CostScopeRun     = "run"
)

const (
	CostPricingDocumented   = "documented"
	CostPricingSubscription = "subscription"
	CostPricingUnknown      = "unknown"
)

const (
	CostSettlementObserved  = "observed"
	CostSettlementEstimated = "estimated"
	CostSettlementUnknown   = "unknown"
)

var (
	ErrInvalidCostLedger         = errors.New("invalid cost ledger")
	ErrCostAllocationRequired    = errors.New("cost allocation required")
	ErrCostSharedCapExceeded     = errors.New("shared cost cap exceeded")
	ErrCostReservationConflict   = errors.New("cost reservation conflict")
	ErrCostSettlementConflict    = errors.New("cost settlement conflict")
	ErrCostReservationNotFound   = errors.New("cost reservation not found")
	ErrCostReservationSettled    = errors.New("cost reservation already settled")
	ErrCostReservationUnbounded  = errors.New("cost reservation must include bounded amount")
	ErrCostUnknownCommittedUsage = errors.New("committed cost usage is unknown")
)

type CostLimits struct {
	AmountMicros *int64 `json:"amountMicros"`
	ModelCalls   int    `json:"modelCalls"`
	Tokens       int    `json:"tokens"`
	WallClockMs  int    `json:"wallClockMs"`
	Attempts     int    `json:"attempts"`
}

type CostSummary struct {
	ReservedAmountMicros        *int64 `json:"reservedAmountMicros"`
	SettledAmountMicros         *int64 `json:"settledAmountMicros"`
	UnknownReservedAmountMicros *int64 `json:"unknownReservedAmountMicros"`
	ReservedCalls               int    `json:"reservedCalls"`
	SettledCalls                int    `json:"settledCalls"`
	ReservedTokens              int    `json:"reservedTokens"`
	SettledTokens               int    `json:"settledTokens"`
}

func (summary CostSummary) CommittedAmountMicros() (int64, bool) {
	if summary.SettledAmountMicros == nil || summary.UnknownReservedAmountMicros == nil {
		return 0, false
	}
	return *summary.SettledAmountMicros + *summary.UnknownReservedAmountMicros, true
}

type CostAllocationEntry struct {
	EntryID            string     `json:"entryId"`
	Sequence           int        `json:"sequence"`
	Kind               string     `json:"kind"`
	OccurredAt         string     `json:"occurredAt"`
	AuthorizationHash  string     `json:"authorizationHash"`
	PreviousLimitsHash *string    `json:"previousLimitsHash"`
	Limits             CostLimits `json:"limits"`
}

type CostReservationEntry struct {
	EntryID              string   `json:"entryId"`
	Sequence             int      `json:"sequence"`
	Kind                 string   `json:"kind"`
	OccurredAt           string   `json:"occurredAt"`
	RunID                string   `json:"runId"`
	RequestID            string   `json:"requestId"`
	IdempotencyKey       string   `json:"idempotencyKey"`
	AuthorizationHash    string   `json:"authorizationHash"`
	ReservedAmountMicros *int64   `json:"reservedAmountMicros"`
	ReservedCalls        int      `json:"reservedCalls"`
	ReservedTokens       int      `json:"reservedTokens"`
	PricingEvidence      []string `json:"pricingEvidence"`
}

type CostSettlementEntry struct {
	EntryID             string   `json:"entryId"`
	Sequence            int      `json:"sequence"`
	Kind                string   `json:"kind"`
	OccurredAt          string   `json:"occurredAt"`
	RunID               string   `json:"runId"`
	RequestID           string   `json:"requestId"`
	IdempotencyKey      string   `json:"idempotencyKey"`
	ReservationHash     string   `json:"reservationHash"`
	Status              string   `json:"status"`
	ChargedAmountMicros *int64   `json:"chargedAmountMicros"`
	ObservedCalls       *int     `json:"observedCalls"`
	ObservedTokens      *int     `json:"observedTokens"`
	ProviderEvidence    []string `json:"providerEvidence"`
}

type CostLedgerBody struct {
	LedgerID       string      `json:"ledgerId"`
	CreatedAt      string      `json:"createdAt"`
	ScopeKind      string      `json:"scopeKind"`
	ScopeHash      string      `json:"scopeHash"`
	Currency       string      `json:"currency"`
	PricingMode    string      `json:"pricingMode"`
	PricingVersion string      `json:"pricingVersion"`
	Entries        []any       `json:"entries"`
	Summary        CostSummary `json:"summary"`
}

type CostLedgerRecord struct {
	SchemaVersion string         `json:"schemaVersion"`
	Ledger        CostLedgerBody `json:"ledger"`
	LedgerHash    string         `json:"ledgerHash"`
}

type CostLedgerConfig struct {
	LedgerID       string
	CreatedAt      time.Time
	ScopeKind      string
	ScopeHash      string
	Currency       string
	PricingMode    string
	PricingVersion string
}

type AllocationInput struct {
	EntryID            string
	OccurredAt         time.Time
	AuthorizationHash  string
	PreviousLimitsHash *string
	Limits             CostLimits
}

type ReservationInput struct {
	EntryID              string
	OccurredAt           time.Time
	RunID                string
	RequestID            string
	IdempotencyKey       string
	AuthorizationHash    string
	ReservedAmountMicros *int64
	ReservedCalls        int
	ReservedTokens       int
	PricingEvidence      []string
}

type SettlementInput struct {
	EntryID             string
	OccurredAt          time.Time
	RunID               string
	RequestID           string
	IdempotencyKey      string
	ReservationHash     string
	Status              string
	ChargedAmountMicros *int64
	ObservedCalls       *int
	ObservedTokens      *int
	ProviderEvidence    []string
}

type ReservationOutcome struct {
	Entry           CostReservationEntry `json:"entry"`
	ReservationHash string               `json:"reservationHash"`
	Appended        bool                 `json:"appended"`
}

type SettlementOutcome struct {
	Entry     CostSettlementEntry `json:"entry"`
	EntryHash string              `json:"entryHash"`
	Appended  bool                `json:"appended"`
}

type costReservationState struct {
	entry CostReservationEntry
	hash  string
}

type costSettlementState struct {
	entry CostSettlementEntry
	hash  string
}

type CostLedger struct {
	mu                      sync.RWMutex
	body                    CostLedgerBody
	entryIDs                map[string]struct{}
	reservationByKey        map[string]costReservationState
	reservationByHash       map[string]costReservationState
	reservationByRequestID  map[string]costReservationState
	settlementByKey         map[string]costSettlementState
	settlementByReservation map[string]costSettlementState
	latestLimits            *CostLimits
	latestLimitsHash        *string
}

func NewCostLedger(config CostLedgerConfig) (*CostLedger, error) {
	if err := validateCostLedgerConfig(config); err != nil {
		return nil, err
	}
	createdAt := config.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	ledger := &CostLedger{
		body: CostLedgerBody{
			LedgerID:       config.LedgerID,
			CreatedAt:      createdAt.UTC().Format(timestampLayout),
			ScopeKind:      config.ScopeKind,
			ScopeHash:      config.ScopeHash,
			Currency:       config.Currency,
			PricingMode:    config.PricingMode,
			PricingVersion: config.PricingVersion,
			Entries:        []any{},
			Summary:        knownZeroCostSummary(),
		},
		entryIDs:                map[string]struct{}{},
		reservationByKey:        map[string]costReservationState{},
		reservationByHash:       map[string]costReservationState{},
		reservationByRequestID:  map[string]costReservationState{},
		settlementByKey:         map[string]costSettlementState{},
		settlementByReservation: map[string]costSettlementState{},
	}
	return ledger, nil
}

func DecodeCostLedgerRecord(input []byte) (CostLedgerRecord, error) {
	type costLedgerBodyWire struct {
		LedgerID       string            `json:"ledgerId"`
		CreatedAt      string            `json:"createdAt"`
		ScopeKind      string            `json:"scopeKind"`
		ScopeHash      string            `json:"scopeHash"`
		Currency       string            `json:"currency"`
		PricingMode    string            `json:"pricingMode"`
		PricingVersion string            `json:"pricingVersion"`
		Entries        []json.RawMessage `json:"entries"`
		Summary        CostSummary       `json:"summary"`
	}
	type costLedgerWire struct {
		SchemaVersion string             `json:"schemaVersion"`
		Ledger        costLedgerBodyWire `json:"ledger"`
		LedgerHash    string             `json:"ledgerHash"`
	}

	var wire costLedgerWire
	if err := evidence.DecodeStrictJSON(input, &wire); err != nil {
		return CostLedgerRecord{}, err
	}
	entries := make([]any, 0, len(wire.Ledger.Entries))
	for _, raw := range wire.Ledger.Entries {
		var kind struct {
			Kind string `json:"kind"`
		}
		if err := json.Unmarshal(raw, &kind); err != nil {
			return CostLedgerRecord{}, err
		}
		switch kind.Kind {
		case "allocation":
			var entry CostAllocationEntry
			if err := evidence.DecodeStrictJSON(raw, &entry); err != nil {
				return CostLedgerRecord{}, err
			}
			entries = append(entries, entry)
		case "reservation":
			var entry CostReservationEntry
			if err := evidence.DecodeStrictJSON(raw, &entry); err != nil {
				return CostLedgerRecord{}, err
			}
			entries = append(entries, entry)
		case "settlement":
			var entry CostSettlementEntry
			if err := evidence.DecodeStrictJSON(raw, &entry); err != nil {
				return CostLedgerRecord{}, err
			}
			entries = append(entries, entry)
		default:
			return CostLedgerRecord{}, fmt.Errorf("%w: unknown entry kind %q", ErrInvalidCostLedger, kind.Kind)
		}
	}

	record := CostLedgerRecord{
		SchemaVersion: wire.SchemaVersion,
		Ledger: CostLedgerBody{
			LedgerID:       wire.Ledger.LedgerID,
			CreatedAt:      wire.Ledger.CreatedAt,
			ScopeKind:      wire.Ledger.ScopeKind,
			ScopeHash:      wire.Ledger.ScopeHash,
			Currency:       wire.Ledger.Currency,
			PricingMode:    wire.Ledger.PricingMode,
			PricingVersion: wire.Ledger.PricingVersion,
			Entries:        entries,
			Summary:        wire.Ledger.Summary,
		},
		LedgerHash: wire.LedgerHash,
	}
	if err := validateCostLedgerRecord(record); err != nil {
		return CostLedgerRecord{}, err
	}
	return record, nil
}

func ValidateCostLedgerRecord(record CostLedgerRecord) error {
	canonical, err := evidence.EncodeCanonical(record)
	if err != nil {
		return err
	}
	_, err = DecodeCostLedgerRecord(canonical)
	return err
}

func LoadCostLedger(record CostLedgerRecord) (*CostLedger, error) {
	canonical, err := evidence.EncodeCanonical(record)
	if err != nil {
		return nil, err
	}
	decoded, err := DecodeCostLedgerRecord(canonical)
	if err != nil {
		return nil, err
	}
	return replayCostLedger(decoded.Ledger)
}

func (ledger *CostLedger) LedgerID() string {
	if ledger == nil {
		return ""
	}
	ledger.mu.RLock()
	defer ledger.mu.RUnlock()
	return ledger.body.LedgerID
}

func (ledger *CostLedger) Summary() CostSummary {
	if ledger == nil {
		return knownZeroCostSummary()
	}
	ledger.mu.RLock()
	defer ledger.mu.RUnlock()
	return cloneCostSummary(ledger.body.Summary)
}

func (ledger *CostLedger) CurrentLimits() (CostLimits, bool) {
	if ledger == nil {
		return CostLimits{}, false
	}
	ledger.mu.RLock()
	defer ledger.mu.RUnlock()
	if ledger.latestLimits == nil {
		return CostLimits{}, false
	}
	return cloneLimits(*ledger.latestLimits), true
}

func (ledger *CostLedger) OutstandingReservations() int {
	if ledger == nil {
		return 0
	}
	ledger.mu.RLock()
	defer ledger.mu.RUnlock()
	count := 0
	for _, reservation := range ledger.reservationByHash {
		if _, settled := ledger.settlementByReservation[reservation.hash]; !settled {
			count++
		}
	}
	return count
}

// RequireOutstandingReservation returns the reservation bound to reservationHash
// only while no settlement, including an unknown settlement, exists for it.
func (ledger *CostLedger) RequireOutstandingReservation(reservationHash string) (CostReservationEntry, error) {
	if ledger == nil {
		return CostReservationEntry{}, fmt.Errorf("%w: ledger is nil", ErrInvalidCostLedger)
	}
	ledger.mu.RLock()
	defer ledger.mu.RUnlock()
	if !evidence.ValidHash(reservationHash) {
		return CostReservationEntry{}, fmt.Errorf("%w: invalid reservationHash", ErrInvalidCostLedger)
	}
	reservation, found := ledger.reservationByHash[reservationHash]
	if !found {
		return CostReservationEntry{}, ErrCostReservationNotFound
	}
	if _, settled := ledger.settlementByReservation[reservationHash]; settled {
		return CostReservationEntry{}, ErrCostReservationSettled
	}
	entry := reservation.entry
	entry.ReservedAmountMicros = cloneInt64Pointer(entry.ReservedAmountMicros)
	entry.PricingEvidence = cloneStrings(entry.PricingEvidence)
	return entry, nil
}

func (ledger *CostLedger) UnknownHoldReservations() int {
	if ledger == nil {
		return 0
	}
	ledger.mu.RLock()
	defer ledger.mu.RUnlock()
	count := 0
	for _, reservation := range ledger.reservationByHash {
		settlement, settled := ledger.settlementByReservation[reservation.hash]
		if !settled || settlement.entry.Status == CostSettlementUnknown {
			count++
		}
	}
	return count
}

func (ledger *CostLedger) Snapshot() (CostLedgerRecord, error) {
	if ledger == nil {
		return CostLedgerRecord{}, fmt.Errorf("%w: ledger is nil", ErrInvalidCostLedger)
	}
	ledger.mu.RLock()
	defer ledger.mu.RUnlock()
	if len(ledger.body.Entries) == 0 {
		return CostLedgerRecord{}, fmt.Errorf("%w: empty entries", ErrInvalidCostLedger)
	}
	body := ledger.body
	body.Entries = append([]any{}, ledger.body.Entries...)
	body.Summary = cloneCostSummary(ledger.body.Summary)
	hash, err := digest(costLedgerDomain, body)
	if err != nil {
		return CostLedgerRecord{}, err
	}
	return CostLedgerRecord{
		SchemaVersion: evidence.SchemaVersion,
		Ledger:        body,
		LedgerHash:    hash,
	}, nil
}

func (ledger *CostLedger) AppendAllocation(input AllocationInput) (CostAllocationEntry, error) {
	if ledger == nil {
		return CostAllocationEntry{}, fmt.Errorf("%w: ledger is nil", ErrInvalidCostLedger)
	}
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	occurredAt := input.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}
	previous := cloneStringPointer(input.PreviousLimitsHash)
	if previous == nil {
		previous = cloneStringPointer(ledger.latestLimitsHash)
	}
	entry := CostAllocationEntry{
		EntryID:            input.EntryID,
		Sequence:           len(ledger.body.Entries) + 1,
		Kind:               "allocation",
		OccurredAt:         occurredAt.UTC().Format(timestampLayout),
		AuthorizationHash:  input.AuthorizationHash,
		PreviousLimitsHash: previous,
		Limits:             cloneLimits(input.Limits),
	}
	if err := ledger.applyAllocation(entry); err != nil {
		return CostAllocationEntry{}, err
	}
	return entry, nil
}

func (ledger *CostLedger) Reserve(input ReservationInput) (ReservationOutcome, error) {
	if ledger == nil {
		return ReservationOutcome{}, fmt.Errorf("%w: ledger is nil", ErrInvalidCostLedger)
	}
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	occurredAt := input.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}
	entry := CostReservationEntry{
		EntryID:              input.EntryID,
		Sequence:             len(ledger.body.Entries) + 1,
		Kind:                 "reservation",
		OccurredAt:           occurredAt.UTC().Format(timestampLayout),
		RunID:                input.RunID,
		RequestID:            input.RequestID,
		IdempotencyKey:       input.IdempotencyKey,
		AuthorizationHash:    input.AuthorizationHash,
		ReservedAmountMicros: cloneInt64Pointer(input.ReservedAmountMicros),
		ReservedCalls:        input.ReservedCalls,
		ReservedTokens:       input.ReservedTokens,
		PricingEvidence:      cloneStrings(input.PricingEvidence),
	}
	reservationHash, appended, err := ledger.applyReservation(entry, true)
	if err != nil {
		return ReservationOutcome{}, err
	}
	if !appended {
		state := ledger.reservationByKey[entry.IdempotencyKey]
		entry = state.entry
		reservationHash = state.hash
	}
	return ReservationOutcome{Entry: entry, ReservationHash: reservationHash, Appended: appended}, nil
}

func (ledger *CostLedger) Settle(input SettlementInput) (SettlementOutcome, error) {
	if ledger == nil {
		return SettlementOutcome{}, fmt.Errorf("%w: ledger is nil", ErrInvalidCostLedger)
	}
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	occurredAt := input.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}
	entry := CostSettlementEntry{
		EntryID:             input.EntryID,
		Sequence:            len(ledger.body.Entries) + 1,
		Kind:                "settlement",
		OccurredAt:          occurredAt.UTC().Format(timestampLayout),
		RunID:               input.RunID,
		RequestID:           input.RequestID,
		IdempotencyKey:      input.IdempotencyKey,
		ReservationHash:     input.ReservationHash,
		Status:              input.Status,
		ChargedAmountMicros: cloneInt64Pointer(input.ChargedAmountMicros),
		ObservedCalls:       cloneIntPointer(input.ObservedCalls),
		ObservedTokens:      cloneIntPointer(input.ObservedTokens),
		ProviderEvidence:    cloneStrings(input.ProviderEvidence),
	}
	entryHash, appended, err := ledger.applySettlement(entry, true)
	if err != nil {
		return SettlementOutcome{}, err
	}
	if !appended {
		state := ledger.settlementByKey[entry.IdempotencyKey]
		entry = state.entry
		entryHash = state.hash
	}
	return SettlementOutcome{Entry: entry, EntryHash: entryHash, Appended: appended}, nil
}

func (ledger *CostLedger) applyAllocation(entry CostAllocationEntry) error {
	if err := validateCostAllocationEntry(entry); err != nil {
		return err
	}
	if err := ledger.validateAppendMetadata(entry.EntryID, entry.Sequence); err != nil {
		return err
	}
	if ledger.latestLimitsHash == nil {
		if entry.PreviousLimitsHash != nil {
			return fmt.Errorf("%w: first allocation must use null previousLimitsHash", ErrInvalidCostLedger)
		}
	} else {
		if entry.PreviousLimitsHash == nil || *entry.PreviousLimitsHash != *ledger.latestLimitsHash {
			return fmt.Errorf("%w: allocation previousLimitsHash mismatch", ErrInvalidCostLedger)
		}
	}
	limitsHash, err := digest(costLimitsHashDomain, entry.Limits)
	if err != nil {
		return err
	}
	ledger.body.Entries = append(ledger.body.Entries, entry)
	ledger.entryIDs[entry.EntryID] = struct{}{}
	limits := cloneLimits(entry.Limits)
	ledger.latestLimits = &limits
	ledger.latestLimitsHash = &limitsHash
	ledger.recomputeSummary()
	return nil
}

func (ledger *CostLedger) applyReservation(entry CostReservationEntry, allowIdempotent bool) (string, bool, error) {
	if state, exists := ledger.reservationByKey[entry.IdempotencyKey]; exists {
		if !allowIdempotent {
			return "", false, fmt.Errorf("%w: duplicate reservation idempotencyKey %q", ErrInvalidCostLedger, entry.IdempotencyKey)
		}
		if !sameReservationPayload(state.entry, entry) {
			return "", false, ErrCostReservationConflict
		}
		return state.hash, false, nil
	}
	if err := validateCostReservationEntry(entry); err != nil {
		return "", false, err
	}
	if err := ledger.validateAppendMetadata(entry.EntryID, entry.Sequence); err != nil {
		return "", false, err
	}
	if ledger.latestLimits == nil {
		return "", false, ErrCostAllocationRequired
	}
	if existing, exists := ledger.reservationByRequestID[entry.RequestID]; exists {
		if existing.entry.IdempotencyKey != entry.IdempotencyKey {
			return "", false, ErrCostReservationConflict
		}
	}
	if err := ledger.validateSharedCaps(entry); err != nil {
		return "", false, err
	}
	hash, err := digest(costReservationHashDomain, entry)
	if err != nil {
		return "", false, err
	}
	if _, exists := ledger.reservationByHash[hash]; exists {
		return "", false, fmt.Errorf("%w: duplicate reservation hash", ErrInvalidCostLedger)
	}
	state := costReservationState{entry: entry, hash: hash}
	ledger.body.Entries = append(ledger.body.Entries, entry)
	ledger.entryIDs[entry.EntryID] = struct{}{}
	ledger.reservationByKey[entry.IdempotencyKey] = state
	ledger.reservationByHash[hash] = state
	ledger.reservationByRequestID[entry.RequestID] = state
	ledger.recomputeSummary()
	return hash, true, nil
}

func (ledger *CostLedger) applySettlement(entry CostSettlementEntry, allowIdempotent bool) (string, bool, error) {
	if state, exists := ledger.settlementByKey[entry.IdempotencyKey]; exists {
		reservation, ok := ledger.reservationByHash[entry.ReservationHash]
		if !ok {
			return "", false, ErrCostReservationNotFound
		}
		normalized, err := normalizeSettlementEntry(entry, reservation.entry)
		if err != nil {
			return "", false, err
		}
		if !allowIdempotent {
			return "", false, fmt.Errorf("%w: duplicate settlement idempotencyKey %q", ErrInvalidCostLedger, entry.IdempotencyKey)
		}
		if !sameSettlementPayload(state.entry, normalized) {
			return "", false, ErrCostSettlementConflict
		}
		return state.hash, false, nil
	}
	if err := validateCostSettlementEntry(entry); err != nil {
		return "", false, err
	}
	if err := ledger.validateAppendMetadata(entry.EntryID, entry.Sequence); err != nil {
		return "", false, err
	}
	reservation, exists := ledger.reservationByHash[entry.ReservationHash]
	if !exists {
		return "", false, ErrCostReservationNotFound
	}
	if reservation.entry.RunID != entry.RunID || reservation.entry.RequestID != entry.RequestID {
		return "", false, fmt.Errorf("%w: settlement run/request mismatch", ErrInvalidCostLedger)
	}
	if previous, settled := ledger.settlementByReservation[entry.ReservationHash]; settled {
		if allowIdempotent && sameSettlementPayload(previous.entry, entry) {
			return previous.hash, false, nil
		}
		return "", false, ErrCostReservationSettled
	}

	normalized, err := normalizeSettlementEntry(entry, reservation.entry)
	if err != nil {
		return "", false, err
	}
	hash, err := digest(costSettlementHashDomain, normalized)
	if err != nil {
		return "", false, err
	}
	state := costSettlementState{entry: normalized, hash: hash}
	ledger.body.Entries = append(ledger.body.Entries, normalized)
	ledger.entryIDs[normalized.EntryID] = struct{}{}
	ledger.settlementByKey[normalized.IdempotencyKey] = state
	ledger.settlementByReservation[normalized.ReservationHash] = state
	ledger.recomputeSummary()
	return hash, true, nil
}

func normalizeSettlementEntry(entry CostSettlementEntry, reservation CostReservationEntry) (CostSettlementEntry, error) {
	switch entry.Status {
	case CostSettlementObserved, CostSettlementEstimated:
		if entry.ChargedAmountMicros == nil {
			return CostSettlementEntry{}, fmt.Errorf("%w: %s settlement requires chargedAmountMicros", ErrInvalidCostLedger, entry.Status)
		}
		if *entry.ChargedAmountMicros < 0 {
			return CostSettlementEntry{}, fmt.Errorf("%w: chargedAmountMicros must be >= 0", ErrInvalidCostLedger)
		}
		if len(entry.ProviderEvidence) == 0 {
			return CostSettlementEntry{}, fmt.Errorf("%w: %s settlement requires providerEvidence", ErrInvalidCostLedger, entry.Status)
		}
		if reservation.ReservedAmountMicros != nil && *entry.ChargedAmountMicros > *reservation.ReservedAmountMicros {
			return CostSettlementEntry{}, fmt.Errorf("%w: chargedAmountMicros exceeds reservation bound", ErrCostSettlementConflict)
		}
		if entry.ObservedCalls == nil {
			entry.ObservedCalls = intPointer(reservation.ReservedCalls)
		}
		if entry.ObservedTokens == nil {
			entry.ObservedTokens = intPointer(reservation.ReservedTokens)
		}
		if entry.ObservedCalls != nil && *entry.ObservedCalls < 0 {
			return CostSettlementEntry{}, fmt.Errorf("%w: observedCalls must be >= 0", ErrInvalidCostLedger)
		}
		if entry.ObservedTokens != nil && *entry.ObservedTokens < 0 {
			return CostSettlementEntry{}, fmt.Errorf("%w: observedTokens must be >= 0", ErrInvalidCostLedger)
		}
	case CostSettlementUnknown:
		if entry.ChargedAmountMicros != nil {
			return CostSettlementEntry{}, fmt.Errorf("%w: unknown settlement cannot carry chargedAmountMicros", ErrInvalidCostLedger)
		}
		entry.ObservedCalls = nil
		entry.ObservedTokens = nil
	default:
		return CostSettlementEntry{}, fmt.Errorf("%w: unknown settlement status %q", ErrInvalidCostLedger, entry.Status)
	}
	entry.ProviderEvidence = cloneStrings(entry.ProviderEvidence)
	return entry, nil
}

func (ledger *CostLedger) validateAppendMetadata(entryID string, sequence int) error {
	if _, exists := ledger.entryIDs[entryID]; exists {
		return fmt.Errorf("%w: duplicate entry id %q", ErrInvalidCostLedger, entryID)
	}
	if sequence != len(ledger.body.Entries)+1 {
		return fmt.Errorf("%w: sequence %d is not append-only", ErrInvalidCostLedger, sequence)
	}
	return nil
}

func (ledger *CostLedger) validateSharedCaps(reservation CostReservationEntry) error {
	limits := ledger.latestLimits
	if limits == nil {
		return ErrCostAllocationRequired
	}
	outstandingAmount, outstandingCalls, outstandingTokens, knownAmount := ledger.outstandingUsage()
	settledAmount, knownSettled := optionalMicrosValue(ledger.body.Summary.SettledAmountMicros)
	if limits.AmountMicros != nil {
		if reservation.ReservedAmountMicros == nil {
			return ErrCostReservationUnbounded
		}
		if !knownAmount || !knownSettled {
			return ErrCostUnknownCommittedUsage
		}
		projected := settledAmount + outstandingAmount + *reservation.ReservedAmountMicros
		if projected > *limits.AmountMicros {
			return ErrCostSharedCapExceeded
		}
	}
	projectedCalls := ledger.body.Summary.SettledCalls + outstandingCalls + reservation.ReservedCalls
	if projectedCalls > limits.ModelCalls {
		return ErrCostSharedCapExceeded
	}
	projectedTokens := ledger.body.Summary.SettledTokens + outstandingTokens + reservation.ReservedTokens
	if projectedTokens > limits.Tokens {
		return ErrCostSharedCapExceeded
	}
	return nil
}

func (ledger *CostLedger) outstandingUsage() (amount int64, calls int, tokens int, knownAmount bool) {
	knownAmount = true
	for _, reservation := range ledger.reservationByHash {
		settlement, settled := ledger.settlementByReservation[reservation.hash]
		if settled && settlement.entry.Status != CostSettlementUnknown {
			continue
		}
		calls += reservation.entry.ReservedCalls
		tokens += reservation.entry.ReservedTokens
		if reservation.entry.ReservedAmountMicros == nil {
			knownAmount = false
			continue
		}
		amount += *reservation.entry.ReservedAmountMicros
	}
	return amount, calls, tokens, knownAmount
}

func (ledger *CostLedger) recomputeSummary() {
	reservedAmount := int64(0)
	reservedAmountKnown := true
	reservedCalls := 0
	reservedTokens := 0
	unknownReservedAmount := int64(0)
	unknownReservedKnown := true

	for _, reservation := range ledger.reservationByHash {
		reservedCalls += reservation.entry.ReservedCalls
		reservedTokens += reservation.entry.ReservedTokens
		if reservation.entry.ReservedAmountMicros == nil {
			reservedAmountKnown = false
		} else {
			reservedAmount += *reservation.entry.ReservedAmountMicros
		}

		settlement, settled := ledger.settlementByReservation[reservation.hash]
		if settled && settlement.entry.Status != CostSettlementUnknown {
			continue
		}
		if reservation.entry.ReservedAmountMicros == nil {
			unknownReservedKnown = false
		} else {
			unknownReservedAmount += *reservation.entry.ReservedAmountMicros
		}
	}

	settledAmount := int64(0)
	settledAmountKnown := true
	settledCalls := 0
	settledTokens := 0
	for _, settlement := range ledger.settlementByReservation {
		if settlement.entry.Status == CostSettlementUnknown {
			continue
		}
		if settlement.entry.ChargedAmountMicros == nil {
			settledAmountKnown = false
		} else {
			settledAmount += *settlement.entry.ChargedAmountMicros
		}
		if settlement.entry.ObservedCalls != nil {
			settledCalls += *settlement.entry.ObservedCalls
		}
		if settlement.entry.ObservedTokens != nil {
			settledTokens += *settlement.entry.ObservedTokens
		}
	}

	ledger.body.Summary = CostSummary{
		ReservedAmountMicros:        pointerIfKnown(reservedAmount, reservedAmountKnown),
		SettledAmountMicros:         pointerIfKnown(settledAmount, settledAmountKnown),
		UnknownReservedAmountMicros: pointerIfKnown(unknownReservedAmount, unknownReservedKnown),
		ReservedCalls:               reservedCalls,
		SettledCalls:                settledCalls,
		ReservedTokens:              reservedTokens,
		SettledTokens:               settledTokens,
	}
}

func validateCostLedgerConfig(config CostLedgerConfig) error {
	if !validCostID(config.LedgerID) {
		return fmt.Errorf("%w: invalid ledgerId", ErrInvalidCostLedger)
	}
	if !validCostScope(config.ScopeKind) {
		return fmt.Errorf("%w: invalid scopeKind %q", ErrInvalidCostLedger, config.ScopeKind)
	}
	if !evidence.ValidHash(config.ScopeHash) {
		return fmt.Errorf("%w: invalid scopeHash", ErrInvalidCostLedger)
	}
	if !validCurrency(config.Currency) {
		return fmt.Errorf("%w: invalid currency %q", ErrInvalidCostLedger, config.Currency)
	}
	if !validPricingMode(config.PricingMode) {
		return fmt.Errorf("%w: invalid pricingMode %q", ErrInvalidCostLedger, config.PricingMode)
	}
	if strings.TrimSpace(config.PricingVersion) == "" || len(config.PricingVersion) > 128 {
		return fmt.Errorf("%w: invalid pricingVersion", ErrInvalidCostLedger)
	}
	return nil
}

func validateCostAllocationEntry(entry CostAllocationEntry) error {
	if entry.Kind != "allocation" {
		return fmt.Errorf("%w: allocation kind mismatch", ErrInvalidCostLedger)
	}
	if !validCostID(entry.EntryID) || entry.Sequence < 1 {
		return fmt.Errorf("%w: invalid allocation identity", ErrInvalidCostLedger)
	}
	if !evidence.ValidHash(entry.AuthorizationHash) {
		return fmt.Errorf("%w: invalid allocation authorizationHash", ErrInvalidCostLedger)
	}
	if entry.PreviousLimitsHash != nil && !evidence.ValidHash(*entry.PreviousLimitsHash) {
		return fmt.Errorf("%w: invalid previousLimitsHash", ErrInvalidCostLedger)
	}
	if _, err := time.Parse(timestampLayout, entry.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid allocation timestamp", ErrInvalidCostLedger)
	}
	if err := validateLimits(entry.Limits); err != nil {
		return err
	}
	return nil
}

func validateCostReservationEntry(entry CostReservationEntry) error {
	if entry.Kind != "reservation" {
		return fmt.Errorf("%w: reservation kind mismatch", ErrInvalidCostLedger)
	}
	if !validCostID(entry.EntryID) || !validCostID(entry.RunID) || !validCostID(entry.RequestID) || !validCostID(entry.IdempotencyKey) || entry.Sequence < 1 {
		return fmt.Errorf("%w: invalid reservation identity", ErrInvalidCostLedger)
	}
	if !evidence.ValidHash(entry.AuthorizationHash) {
		return fmt.Errorf("%w: invalid reservation authorizationHash", ErrInvalidCostLedger)
	}
	if _, err := time.Parse(timestampLayout, entry.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid reservation timestamp", ErrInvalidCostLedger)
	}
	if entry.ReservedAmountMicros != nil && *entry.ReservedAmountMicros < 0 {
		return fmt.Errorf("%w: reservedAmountMicros must be >= 0", ErrInvalidCostLedger)
	}
	if entry.ReservedCalls < 1 || entry.ReservedTokens < 0 {
		return fmt.Errorf("%w: invalid reserved calls/tokens", ErrInvalidCostLedger)
	}
	if err := validateHashSet(entry.PricingEvidence, "pricingEvidence"); err != nil {
		return err
	}
	return nil
}

func validateCostSettlementEntry(entry CostSettlementEntry) error {
	if entry.Kind != "settlement" {
		return fmt.Errorf("%w: settlement kind mismatch", ErrInvalidCostLedger)
	}
	if !validCostID(entry.EntryID) || !validCostID(entry.RunID) || !validCostID(entry.RequestID) || !validCostID(entry.IdempotencyKey) || entry.Sequence < 1 {
		return fmt.Errorf("%w: invalid settlement identity", ErrInvalidCostLedger)
	}
	if !evidence.ValidHash(entry.ReservationHash) {
		return fmt.Errorf("%w: invalid settlement reservationHash", ErrInvalidCostLedger)
	}
	if _, err := time.Parse(timestampLayout, entry.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid settlement timestamp", ErrInvalidCostLedger)
	}
	if entry.ObservedCalls != nil && *entry.ObservedCalls < 0 {
		return fmt.Errorf("%w: observedCalls must be >= 0", ErrInvalidCostLedger)
	}
	if entry.ObservedTokens != nil && *entry.ObservedTokens < 0 {
		return fmt.Errorf("%w: observedTokens must be >= 0", ErrInvalidCostLedger)
	}
	if err := validateHashSet(entry.ProviderEvidence, "providerEvidence"); err != nil {
		return err
	}
	return nil
}

func validateLimits(limits CostLimits) error {
	if limits.AmountMicros != nil && *limits.AmountMicros < 0 {
		return fmt.Errorf("%w: limits.amountMicros must be >= 0", ErrInvalidCostLedger)
	}
	if limits.ModelCalls < 0 || limits.Tokens < 0 || limits.WallClockMs < 1 || limits.Attempts < 1 {
		return fmt.Errorf("%w: invalid limits", ErrInvalidCostLedger)
	}
	return nil
}

func validateCostLedgerRecord(record CostLedgerRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion {
		return fmt.Errorf("%w: unsupported schema version", ErrInvalidCostLedger)
	}
	if !evidence.ValidHash(record.LedgerHash) {
		return fmt.Errorf("%w: invalid ledgerHash", ErrInvalidCostLedger)
	}
	if err := validateCostLedgerConfig(CostLedgerConfig{
		LedgerID:       record.Ledger.LedgerID,
		ScopeKind:      record.Ledger.ScopeKind,
		ScopeHash:      record.Ledger.ScopeHash,
		Currency:       record.Ledger.Currency,
		PricingMode:    record.Ledger.PricingMode,
		PricingVersion: record.Ledger.PricingVersion,
	}); err != nil {
		return err
	}
	if _, err := time.Parse(timestampLayout, record.Ledger.CreatedAt); err != nil {
		return fmt.Errorf("%w: invalid createdAt", ErrInvalidCostLedger)
	}
	if len(record.Ledger.Entries) == 0 {
		return fmt.Errorf("%w: empty entries", ErrInvalidCostLedger)
	}

	replayed, err := replayCostLedger(record.Ledger)
	if err != nil {
		return err
	}
	if !equalSummary(record.Ledger.Summary, replayed.body.Summary) {
		return fmt.Errorf("%w: summary mismatch", ErrInvalidCostLedger)
	}
	expectedHash, err := digest(costLedgerDomain, record.Ledger)
	if err != nil {
		return err
	}
	if expectedHash != record.LedgerHash {
		return fmt.Errorf("%w: ledger hash mismatch", ErrInvalidCostLedger)
	}
	return nil
}

func replayCostLedger(body CostLedgerBody) (*CostLedger, error) {
	createdAt, err := time.Parse(timestampLayout, body.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid createdAt", ErrInvalidCostLedger)
	}
	ledger, err := NewCostLedger(CostLedgerConfig{
		LedgerID:       body.LedgerID,
		CreatedAt:      createdAt,
		ScopeKind:      body.ScopeKind,
		ScopeHash:      body.ScopeHash,
		Currency:       body.Currency,
		PricingMode:    body.PricingMode,
		PricingVersion: body.PricingVersion,
	})
	if err != nil {
		return nil, err
	}
	for _, raw := range body.Entries {
		switch entry := raw.(type) {
		case CostAllocationEntry:
			if err := ledger.applyAllocation(entry); err != nil {
				return nil, err
			}
		case CostReservationEntry:
			if _, _, err := ledger.applyReservation(entry, false); err != nil {
				return nil, err
			}
		case CostSettlementEntry:
			if _, _, err := ledger.applySettlement(entry, false); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("%w: unsupported entry type %T", ErrInvalidCostLedger, raw)
		}
	}
	return ledger, nil
}

func validCostScope(value string) bool {
	switch value {
	case CostScopeOwner, CostScopeProject, CostScopeRun:
		return true
	default:
		return false
	}
}

func validPricingMode(value string) bool {
	switch value {
	case CostPricingDocumented, CostPricingSubscription, CostPricingUnknown:
		return true
	default:
		return false
	}
}

func validCurrency(value string) bool {
	if len(value) != 3 || strings.ToUpper(value) != value {
		return false
	}
	for _, r := range value {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}

func validCostID(value string) bool {
	return len(value) > 0 && len(value) <= 64 && evidence.ValidID(value)
}

func validateHashSet(values []string, field string) error {
	seen := map[string]struct{}{}
	for _, value := range values {
		if !evidence.ValidHash(value) {
			return fmt.Errorf("%w: invalid %s hash", ErrInvalidCostLedger, field)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%w: duplicate %s hash", ErrInvalidCostLedger, field)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func sameReservationPayload(existing CostReservationEntry, incoming CostReservationEntry) bool {
	if existing.RunID != incoming.RunID || existing.RequestID != incoming.RequestID || existing.IdempotencyKey != incoming.IdempotencyKey || existing.AuthorizationHash != incoming.AuthorizationHash {
		return false
	}
	if existing.ReservedCalls != incoming.ReservedCalls || existing.ReservedTokens != incoming.ReservedTokens {
		return false
	}
	if !equalOptionalMicros(existing.ReservedAmountMicros, incoming.ReservedAmountMicros) {
		return false
	}
	return equalStringSlice(existing.PricingEvidence, incoming.PricingEvidence)
}

func sameSettlementPayload(existing CostSettlementEntry, incoming CostSettlementEntry) bool {
	if existing.RunID != incoming.RunID || existing.RequestID != incoming.RequestID || existing.IdempotencyKey != incoming.IdempotencyKey || existing.ReservationHash != incoming.ReservationHash || existing.Status != incoming.Status {
		return false
	}
	if !equalOptionalMicros(existing.ChargedAmountMicros, incoming.ChargedAmountMicros) {
		return false
	}
	if !equalOptionalInt(existing.ObservedCalls, incoming.ObservedCalls) || !equalOptionalInt(existing.ObservedTokens, incoming.ObservedTokens) {
		return false
	}
	return equalStringSlice(existing.ProviderEvidence, incoming.ProviderEvidence)
}

func equalSummary(left, right CostSummary) bool {
	return equalOptionalMicros(left.ReservedAmountMicros, right.ReservedAmountMicros) &&
		equalOptionalMicros(left.SettledAmountMicros, right.SettledAmountMicros) &&
		equalOptionalMicros(left.UnknownReservedAmountMicros, right.UnknownReservedAmountMicros) &&
		left.ReservedCalls == right.ReservedCalls &&
		left.SettledCalls == right.SettledCalls &&
		left.ReservedTokens == right.ReservedTokens &&
		left.SettledTokens == right.SettledTokens
}

func equalOptionalMicros(left, right *int64) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func equalOptionalInt(left, right *int) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func equalStringSlice(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func knownZeroCostSummary() CostSummary {
	return CostSummary{
		ReservedAmountMicros:        int64Pointer(0),
		SettledAmountMicros:         int64Pointer(0),
		UnknownReservedAmountMicros: int64Pointer(0),
		ReservedCalls:               0,
		SettledCalls:                0,
		ReservedTokens:              0,
		SettledTokens:               0,
	}
}

func pointerIfKnown(value int64, known bool) *int64 {
	if !known {
		return nil
	}
	return int64Pointer(value)
}

func int64Pointer(value int64) *int64 {
	v := value
	return &v
}

func intPointer(value int) *int {
	v := value
	return &v
}

func cloneCostSummary(summary CostSummary) CostSummary {
	return CostSummary{
		ReservedAmountMicros:        cloneInt64Pointer(summary.ReservedAmountMicros),
		SettledAmountMicros:         cloneInt64Pointer(summary.SettledAmountMicros),
		UnknownReservedAmountMicros: cloneInt64Pointer(summary.UnknownReservedAmountMicros),
		ReservedCalls:               summary.ReservedCalls,
		SettledCalls:                summary.SettledCalls,
		ReservedTokens:              summary.ReservedTokens,
		SettledTokens:               summary.SettledTokens,
	}
}

func cloneLimits(limits CostLimits) CostLimits {
	return CostLimits{
		AmountMicros: cloneInt64Pointer(limits.AmountMicros),
		ModelCalls:   limits.ModelCalls,
		Tokens:       limits.Tokens,
		WallClockMs:  limits.WallClockMs,
		Attempts:     limits.Attempts,
	}
}

func cloneInt64Pointer(value *int64) *int64 {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}

func cloneIntPointer(value *int) *int {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}

func cloneStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}

func cloneStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}

func optionalMicrosValue(value *int64) (int64, bool) {
	if value == nil {
		return 0, false
	}
	return *value, true
}
