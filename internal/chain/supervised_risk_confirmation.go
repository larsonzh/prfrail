package chain

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// The binding risk confirmation is the start-up precondition of the
// supervised-black-box mode (contracts section 7.1a). It fixes run/task/step/
// attempt, the isolated-workspace digest, the risk-notice version and hash, the
// visible delivery fact and the request/response interaction evidence. Its
// selection is permanently "confirm": a rejection is never a confirmation, and
// a notice version or hash that drifts from the current constants is refused
// before any work window opens.
const (
	supervisedRiskConfirmationDomain = "proofrail:supervised-risk-confirmation:1\n"
	supervisedRiskConfirmationKind   = "supervised-risk-confirmation"
	supervisedRiskConfirmationSchema = "1.0.0"
	supervisedVisibleChannel         = "sessionbridge-visible"
	supervisedTimeLayout             = "2006-01-02T15:04:05.000Z"
	// supervisedDeliveryAccepted is the only delivery status that may open a
	// work window (contracts section 7.1a). The record shape also permits
	// "unavailable" and "unknown" so a delivery fact can be recorded honestly,
	// but the window opener refuses both.
	supervisedDeliveryAccepted = "accepted"
)

var (
	// ErrInvalidSupervisedConfirmation reports a risk confirmation that cannot
	// be trusted: a wrong schema version or hash, a malformed binding, a
	// selection other than "confirm", an illegal delivery fact or malformed
	// evidence.
	ErrInvalidSupervisedConfirmation = errors.New("invalid supervised risk confirmation")
	// ErrInvalidSupervisedReport reports a reduced/unknown report that cannot be
	// trusted: a wrong assurance level, an unknown dimension outside the frozen
	// closed set, an empty or cost-free unknown set, or a malformed hash.
	ErrInvalidSupervisedReport = errors.New("invalid supervised report")
	// supervisedTransportIDPattern is the accepted transport request id shape.
	// It is deliberately broader than evidence.ValidID because the visible
	// transport id is owned by SessionBridge, not by the ProofRail id grammar.
	supervisedTransportIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{0,127}$`)
)

// SupervisedDelivery is the factual record of one visible delivery attempt. It
// carries no completion claim: the delivery only proves that a notice was
// handed to the operator UI, never that any work happened.
type SupervisedDelivery struct {
	Channel            string   `json:"channel"`
	DeliveredAt        string   `json:"deliveredAt"`
	DeliveryStatus     string   `json:"deliveryStatus"`
	TransportRequestID *string  `json:"transportRequestId"`
	ErrorEvidence      []string `json:"errorEvidence"`
}

// SupervisedRiskConfirmation is the record body. It is exported so callers can
// read the binding, but it must always be built and validated as a whole.
type SupervisedRiskConfirmation struct {
	Kind                string             `json:"kind"`
	RecordID            string             `json:"recordId"`
	IssuedAt            string             `json:"issuedAt"`
	IssuedBy            evidence.Actor     `json:"issuedBy"`
	RunID               string             `json:"runId"`
	TaskID              string             `json:"taskId"`
	StepID              string             `json:"stepId"`
	Attempt             int                `json:"attempt"`
	WorkspaceHash       string             `json:"workspaceHash"`
	RiskNoticeVersion   int                `json:"riskNoticeVersion"`
	RiskNoticeHash      string             `json:"riskNoticeHash"`
	Selection           string             `json:"selection"`
	Delivery            SupervisedDelivery `json:"delivery"`
	InteractionEvidence []string           `json:"interactionEvidence"`
	Evidence            []string           `json:"evidence"`
}

// SupervisedRiskConfirmationRecord is the persisted envelope. The canonical
// record hash domain separator is applied by the hash layer, never by the
// stored JSON.
type SupervisedRiskConfirmationRecord struct {
	SchemaVersion string                     `json:"schemaVersion"`
	Record        SupervisedRiskConfirmation `json:"record"`
	RecordHash    string                     `json:"recordHash"`
}

// SupervisedRiskConfirmationInput is the caller-supplied body of a new
// confirmation. Selection is not an input: it is always "confirm".
type SupervisedRiskConfirmationInput struct {
	RecordID            string
	IssuedAt            string
	OperatorID          string
	RunID               string
	TaskID              string
	StepID              string
	Attempt             int
	WorkspaceHash       string
	RiskNoticeVersion   int
	RiskNoticeHash      string
	Delivery            SupervisedDelivery
	InteractionEvidence []string
	Evidence            []string
}

// NewSupervisedRiskConfirmation builds and validates a confirmation record. The
// record hash is recomputed over the canonical body so a caller can never hand
// the engine a body that does not match its own digest.
func NewSupervisedRiskConfirmation(input SupervisedRiskConfirmationInput, clock Clock, ids IDSource) (SupervisedRiskConfirmationRecord, error) {
	if clock == nil {
		clock = defaultClock
	}
	if ids == nil {
		ids = randomID
	}
	body := SupervisedRiskConfirmation{
		Kind:                supervisedRiskConfirmationKind,
		RecordID:            input.RecordID,
		IssuedAt:            input.IssuedAt,
		IssuedBy:            evidence.Actor{Type: "operator", ID: input.OperatorID},
		RunID:               input.RunID,
		TaskID:              input.TaskID,
		StepID:              input.StepID,
		Attempt:             input.Attempt,
		WorkspaceHash:       input.WorkspaceHash,
		RiskNoticeVersion:   input.RiskNoticeVersion,
		RiskNoticeHash:      input.RiskNoticeHash,
		Selection:           "confirm",
		Delivery:            input.Delivery,
		InteractionEvidence: cloneStrings(input.InteractionEvidence),
		Evidence:            cloneStrings(input.Evidence),
	}
	if body.RecordID == "" {
		body.RecordID = ids()
	}
	if body.IssuedAt == "" {
		body.IssuedAt = clock().UTC().Format(supervisedTimeLayout)
	}
	if body.Delivery.ErrorEvidence == nil {
		body.Delivery.ErrorEvidence = []string{}
	}
	hash, err := digestRecord(supervisedRiskConfirmationDomain, body)
	if err != nil {
		return SupervisedRiskConfirmationRecord{}, err
	}
	record := SupervisedRiskConfirmationRecord{SchemaVersion: supervisedRiskConfirmationSchema, Record: body, RecordHash: hash}
	if err := ValidateSupervisedRiskConfirmation(record); err != nil {
		return SupervisedRiskConfirmationRecord{}, err
	}
	return record, nil
}

// DecodeSupervisedRiskConfirmation strictly decodes and validates one record.
func DecodeSupervisedRiskConfirmation(input []byte) (SupervisedRiskConfirmationRecord, error) {
	var record SupervisedRiskConfirmationRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return SupervisedRiskConfirmationRecord{}, fmt.Errorf("%w: %v", ErrInvalidSupervisedConfirmation, err)
	}
	if err := ValidateSupervisedRiskConfirmation(record); err != nil {
		return SupervisedRiskConfirmationRecord{}, err
	}
	return record, nil
}

// ValidateSupervisedRiskConfirmation checks the envelope and the body, and
// recomputes the canonical record hash so a tampered body can never validate.
func ValidateSupervisedRiskConfirmation(record SupervisedRiskConfirmationRecord) error {
	if record.SchemaVersion != supervisedRiskConfirmationSchema || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidSupervisedConfirmation)
	}
	if err := validateSupervisedRiskConfirmationBody(record.Record); err != nil {
		return err
	}
	expected, err := digestRecord(supervisedRiskConfirmationDomain, record.Record)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidSupervisedConfirmation)
	}
	return nil
}

func validateSupervisedRiskConfirmationBody(body SupervisedRiskConfirmation) error {
	if body.Kind != supervisedRiskConfirmationKind {
		return fmt.Errorf("%w: invalid kind", ErrInvalidSupervisedConfirmation)
	}
	if !evidence.ValidID(body.RecordID) || !evidence.ValidID(body.RunID) || !evidence.ValidID(body.TaskID) || !evidence.ValidID(body.StepID) || body.Attempt < 1 {
		return fmt.Errorf("%w: invalid binding identity", ErrInvalidSupervisedConfirmation)
	}
	if body.IssuedBy.Type != "operator" || !evidence.ValidID(body.IssuedBy.ID) {
		return fmt.Errorf("%w: invalid issuing operator", ErrInvalidSupervisedConfirmation)
	}
	if !validSupervisedTimestamp(body.IssuedAt) {
		return fmt.Errorf("%w: invalid issue time", ErrInvalidSupervisedConfirmation)
	}
	if !evidence.ValidHash(body.WorkspaceHash) {
		return fmt.Errorf("%w: invalid workspace digest", ErrInvalidSupervisedConfirmation)
	}
	if body.RiskNoticeVersion < 1 || !evidence.ValidHash(body.RiskNoticeHash) {
		return fmt.Errorf("%w: invalid risk notice binding", ErrInvalidSupervisedConfirmation)
	}
	if body.Selection != "confirm" {
		return fmt.Errorf("%w: selection must be confirm", ErrInvalidSupervisedConfirmation)
	}
	if err := validateSupervisedDelivery(body.Delivery); err != nil {
		return err
	}
	if len(body.InteractionEvidence) != 2 || !validUniqueHashSet(body.InteractionEvidence, true) {
		return fmt.Errorf("%w: interaction evidence must bind the request and the response", ErrInvalidSupervisedConfirmation)
	}
	if !validUniqueHashSet(body.Evidence, true) {
		return fmt.Errorf("%w: confirmation evidence must be a non-empty unique hash set", ErrInvalidSupervisedConfirmation)
	}
	return nil
}

func validateSupervisedDelivery(delivery SupervisedDelivery) error {
	if delivery.Channel != supervisedVisibleChannel {
		return fmt.Errorf("%w: unknown delivery channel", ErrInvalidSupervisedConfirmation)
	}
	if !validSupervisedTimestamp(delivery.DeliveredAt) {
		return fmt.Errorf("%w: invalid delivery time", ErrInvalidSupervisedConfirmation)
	}
	switch delivery.DeliveryStatus {
	case "accepted", "unavailable", "unknown":
	default:
		return fmt.Errorf("%w: unknown delivery status", ErrInvalidSupervisedConfirmation)
	}
	if delivery.TransportRequestID != nil && !supervisedTransportIDPattern.MatchString(*delivery.TransportRequestID) {
		return fmt.Errorf("%w: invalid transport request id", ErrInvalidSupervisedConfirmation)
	}
	if !validUniqueHashSet(delivery.ErrorEvidence, false) {
		return fmt.Errorf("%w: invalid delivery error evidence", ErrInvalidSupervisedConfirmation)
	}
	return nil
}

func validSupervisedTimestamp(value string) bool {
	_, err := time.Parse(supervisedTimeLayout, value)
	return err == nil
}

// SupervisedRiskConfirmationBinding identifies the step attempt a confirmation
// belongs to.
type SupervisedRiskConfirmationBinding struct {
	RunID   string
	TaskID  string
	StepID  string
	Attempt int
}

// SupervisedRiskConfirmationLedger is the append-only, no-replace collection of
// confirmations. Append is idempotent for the same step attempt and the same
// digest, and refuses a second, divergent confirmation for the same step
// attempt instead of silently replacing it.
type SupervisedRiskConfirmationLedger struct {
	Records []SupervisedRiskConfirmationRecord
}

// Append validates and records one confirmation.
func (ledger *SupervisedRiskConfirmationLedger) Append(record SupervisedRiskConfirmationRecord) error {
	if err := ValidateSupervisedRiskConfirmation(record); err != nil {
		return err
	}
	for _, existing := range ledger.Records {
		if existing.Record.RunID != record.Record.RunID || existing.Record.TaskID != record.Record.TaskID ||
			existing.Record.StepID != record.Record.StepID || existing.Record.Attempt != record.Record.Attempt {
			continue
		}
		if existing.RecordHash == record.RecordHash {
			// Same step attempt, same digest: a replay converges on the record
			// that is already durable.
			return nil
		}
		return fmt.Errorf("%w: a different confirmation already binds this step attempt", ErrInvalidSupervisedConfirmation)
	}
	ledger.Records = append(ledger.Records, record)
	return nil
}

// Confirmed reports whether a valid confirmation binds the step attempt.
func (ledger SupervisedRiskConfirmationLedger) Confirmed(binding SupervisedRiskConfirmationBinding) bool {
	for _, record := range ledger.Records {
		if record.Record.RunID == binding.RunID && record.Record.TaskID == binding.TaskID &&
			record.Record.StepID == binding.StepID && record.Record.Attempt == binding.Attempt {
			return true
		}
	}
	return false
}
