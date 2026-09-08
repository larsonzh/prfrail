package gates

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	effectRecordDomain    = "proofrail:effect-record:1\n"
	effectTimestampLayout = "2006-01-02T15:04:05.000Z"
	effectSchemaVersion   = "1.0.0"
)

var ErrInvalidEffectRecord = errors.New("invalid effect record")

// EffectClass separates read-only, locally discardable and external-write
// side effects (S0 capability declaration, trusted policy origin only).
type EffectClass string

const (
	EffectReadOnly         EffectClass = "read-only"
	EffectLocalDiscardable EffectClass = "local-discardable"
	EffectExternalWrite    EffectClass = "external-write"
)

// EffectAdmission is the S1 enforcement decision for a classified effect.
type EffectAdmission string

const (
	EffectAdmit     EffectAdmission = "admit"
	EffectDeny      EffectAdmission = "deny"
	EffectUncertain EffectAdmission = "uncertain"
)

var (
	validEffectClasses = map[EffectClass]bool{
		EffectReadOnly: true, EffectLocalDiscardable: true, EffectExternalWrite: true,
	}
	validEffectSubjectKinds = map[string]bool{"hook": true, "adapter": true, "manual-action": true}
	validRecoveryGuarantees = map[string]bool{
		"none-required": true, "discard-workspace": true, "reconcile-only": true, "compensation-only": true,
	}
	validObservationOutcomes = map[string]bool{"completed": true, "failed": true, "uncertain": true}
)

// EffectClassification mirrors the classification body of effect-record.schema.json.
type EffectClassification struct {
	RecordID              string   `json:"recordId"`
	Kind                  string   `json:"kind"`
	CreatedAt             string   `json:"createdAt"`
	PolicyHash            string   `json:"policyHash"`
	SubjectKind           string   `json:"subjectKind"`
	SubjectID             string   `json:"subjectId"`
	SubjectDefinitionHash string   `json:"subjectDefinitionHash"`
	EffectClass           string   `json:"effectClass"`
	ExternalSystems       []string `json:"externalSystems"`
	RecoveryGuarantee     string   `json:"recoveryGuarantee"`
	Evidence              []string `json:"evidence"`
}

// EffectClassificationRecord is the wire record for a classification.
type EffectClassificationRecord struct {
	SchemaVersion string               `json:"schemaVersion"`
	Record        EffectClassification `json:"record"`
	RecordHash    string               `json:"recordHash"`
}

// EffectObservation mirrors the observation body of effect-record.schema.json.
type EffectObservation struct {
	RecordID           string   `json:"recordId"`
	Kind               string   `json:"kind"`
	OccurredAt         string   `json:"occurredAt"`
	RunID              string   `json:"runId"`
	TaskID             string   `json:"taskId"`
	Attempt            int      `json:"attempt"`
	ClassificationHash string   `json:"classificationHash"`
	AuthorizationHash  string   `json:"authorizationHash"`
	EffectClass        string   `json:"effectClass"`
	OperationKey       string   `json:"operationKey"`
	Outcome            string   `json:"outcome"`
	Evidence           []string `json:"evidence"`
	ErrorEvidence      []string `json:"errorEvidence"`
}

// EffectObservationRecord is the wire record for an observation.
type EffectObservationRecord struct {
	SchemaVersion string            `json:"schemaVersion"`
	Record        EffectObservation `json:"record"`
	RecordHash    string            `json:"recordHash"`
}

// NewEffectClassificationRecord validates and hashes a classification.
func NewEffectClassificationRecord(classification EffectClassification) (EffectClassificationRecord, error) {
	if err := validateEffectClassification(classification); err != nil {
		return EffectClassificationRecord{}, err
	}
	hash, err := effectRecordHash(classification)
	if err != nil {
		return EffectClassificationRecord{}, err
	}
	return EffectClassificationRecord{SchemaVersion: effectSchemaVersion, Record: classification, RecordHash: hash}, nil
}

// DecodeEffectClassificationRecord strictly decodes and validates a wire record.
func DecodeEffectClassificationRecord(input []byte) (EffectClassificationRecord, error) {
	var record EffectClassificationRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return EffectClassificationRecord{}, err
	}
	if err := ValidateEffectClassificationRecord(record); err != nil {
		return EffectClassificationRecord{}, err
	}
	return record, nil
}

// ValidateEffectClassificationRecord re-validates body and record hash.
func ValidateEffectClassificationRecord(record EffectClassificationRecord) error {
	if record.SchemaVersion != effectSchemaVersion {
		return fmt.Errorf("%w: unsupported schemaVersion %q", ErrInvalidEffectRecord, record.SchemaVersion)
	}
	if err := validateEffectClassification(record.Record); err != nil {
		return err
	}
	expected, err := effectRecordHash(record.Record)
	if err != nil {
		return err
	}
	if record.RecordHash != expected {
		return fmt.Errorf("%w: classification hash mismatch", ErrInvalidEffectRecord)
	}
	return nil
}

// NewEffectObservationRecord validates and hashes an observation.
func NewEffectObservationRecord(observation EffectObservation) (EffectObservationRecord, error) {
	if err := validateEffectObservation(observation); err != nil {
		return EffectObservationRecord{}, err
	}
	hash, err := effectRecordHash(observation)
	if err != nil {
		return EffectObservationRecord{}, err
	}
	return EffectObservationRecord{SchemaVersion: effectSchemaVersion, Record: observation, RecordHash: hash}, nil
}

// DecodeEffectObservationRecord strictly decodes and validates a wire record.
func DecodeEffectObservationRecord(input []byte) (EffectObservationRecord, error) {
	var record EffectObservationRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return EffectObservationRecord{}, err
	}
	if err := ValidateEffectObservationRecord(record); err != nil {
		return EffectObservationRecord{}, err
	}
	return record, nil
}

// ValidateEffectObservationRecord re-validates body and record hash.
func ValidateEffectObservationRecord(record EffectObservationRecord) error {
	if record.SchemaVersion != effectSchemaVersion {
		return fmt.Errorf("%w: unsupported schemaVersion %q", ErrInvalidEffectRecord, record.SchemaVersion)
	}
	if err := validateEffectObservation(record.Record); err != nil {
		return err
	}
	expected, err := effectRecordHash(record.Record)
	if err != nil {
		return err
	}
	if record.RecordHash != expected {
		return fmt.Errorf("%w: observation hash mismatch", ErrInvalidEffectRecord)
	}
	return nil
}

// AdmitEffect enforces the S1 rule: external writes are denied by default,
// unknown classes fail closed, read-only and local-discardable are admitted.
// Command names or dry-run wording never influence the decision.
func AdmitEffect(classification EffectClassification) EffectAdmission {
	switch EffectClass(classification.EffectClass) {
	case EffectReadOnly, EffectLocalDiscardable:
		return EffectAdmit
	case EffectExternalWrite:
		return EffectDeny
	default:
		return EffectDeny
	}
}

// ObservationRecoveryAction maps an observation outcome to the recovery
// action; uncertain and failed outcomes reconcile instead of blind retry.
func ObservationRecoveryAction(observation EffectObservation) string {
	switch observation.Outcome {
	case "completed":
		return "none-required"
	default:
		return "reconcile-only"
	}
}

func effectRecordHash(body any) (string, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	canonical, err := evidence.Canonicalize(payload)
	if err != nil {
		return "", fmt.Errorf("%w: canonicalize effect record: %v", ErrInvalidEffectRecord, err)
	}
	return evidence.Digest(effectRecordDomain, canonical), nil
}

func validEffectID(value string) bool {
	return evidence.ValidID(value) && len(value) <= 64
}

func parseEffectTimestamp(value string) (time.Time, error) {
	return time.Parse(effectTimestampLayout, value)
}

func validateEffectHashSet(values []string, field string) error {
	seen := map[string]bool{}
	for _, value := range values {
		if !evidence.ValidHash(value) {
			return fmt.Errorf("%w: invalid %s hash", ErrInvalidEffectRecord, field)
		}
		if seen[value] {
			return fmt.Errorf("%w: duplicate %s hash", ErrInvalidEffectRecord, field)
		}
		seen[value] = true
	}
	return nil
}

func validateEffectClassification(classification EffectClassification) error {
	if !validEffectID(classification.RecordID) {
		return fmt.Errorf("%w: invalid recordId", ErrInvalidEffectRecord)
	}
	if classification.Kind != "classification" {
		return fmt.Errorf("%w: kind must be classification", ErrInvalidEffectRecord)
	}
	if _, err := parseEffectTimestamp(classification.CreatedAt); err != nil {
		return fmt.Errorf("%w: invalid createdAt", ErrInvalidEffectRecord)
	}
	if !evidence.ValidHash(classification.PolicyHash) {
		return fmt.Errorf("%w: invalid policyHash", ErrInvalidEffectRecord)
	}
	if !validEffectSubjectKinds[classification.SubjectKind] {
		return fmt.Errorf("%w: unknown subjectKind %q", ErrInvalidEffectRecord, classification.SubjectKind)
	}
	if !validEffectID(classification.SubjectID) {
		return fmt.Errorf("%w: invalid subjectId", ErrInvalidEffectRecord)
	}
	if !evidence.ValidHash(classification.SubjectDefinitionHash) {
		return fmt.Errorf("%w: invalid subjectDefinitionHash", ErrInvalidEffectRecord)
	}
	if !validEffectClasses[EffectClass(classification.EffectClass)] {
		return fmt.Errorf("%w: unknown effectClass %q", ErrInvalidEffectRecord, classification.EffectClass)
	}
	seenSystems := map[string]bool{}
	for _, system := range classification.ExternalSystems {
		if !validEffectID(system) {
			return fmt.Errorf("%w: invalid external system id", ErrInvalidEffectRecord)
		}
		if seenSystems[system] {
			return fmt.Errorf("%w: duplicate external system %q", ErrInvalidEffectRecord, system)
		}
		seenSystems[system] = true
	}
	if !validRecoveryGuarantees[classification.RecoveryGuarantee] {
		return fmt.Errorf("%w: unknown recoveryGuarantee %q", ErrInvalidEffectRecord, classification.RecoveryGuarantee)
	}
	if classification.EffectClass == string(EffectExternalWrite) &&
		(classification.RecoveryGuarantee == "none-required" || classification.RecoveryGuarantee == "discard-workspace") {
		return fmt.Errorf("%w: external-write cannot claim %q recovery", ErrInvalidEffectRecord, classification.RecoveryGuarantee)
	}
	if len(classification.Evidence) < 1 {
		return fmt.Errorf("%w: classification evidence must not be empty", ErrInvalidEffectRecord)
	}
	return validateEffectHashSet(classification.Evidence, "evidence")
}

func validateEffectObservation(observation EffectObservation) error {
	if !validEffectID(observation.RecordID) {
		return fmt.Errorf("%w: invalid recordId", ErrInvalidEffectRecord)
	}
	if observation.Kind != "observation" {
		return fmt.Errorf("%w: kind must be observation", ErrInvalidEffectRecord)
	}
	if _, err := parseEffectTimestamp(observation.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid occurredAt", ErrInvalidEffectRecord)
	}
	if !evidence.ValidID(observation.RunID) || !evidence.ValidID(observation.TaskID) || observation.Attempt < 1 {
		return fmt.Errorf("%w: invalid observation identity", ErrInvalidEffectRecord)
	}
	if !evidence.ValidHash(observation.ClassificationHash) || !evidence.ValidHash(observation.AuthorizationHash) {
		return fmt.Errorf("%w: invalid observation hashes", ErrInvalidEffectRecord)
	}
	if !validEffectClasses[EffectClass(observation.EffectClass)] {
		return fmt.Errorf("%w: unknown effectClass %q", ErrInvalidEffectRecord, observation.EffectClass)
	}
	if !validEffectID(observation.OperationKey) {
		return fmt.Errorf("%w: invalid operationKey", ErrInvalidEffectRecord)
	}
	if !validObservationOutcomes[observation.Outcome] {
		return fmt.Errorf("%w: unknown outcome %q", ErrInvalidEffectRecord, observation.Outcome)
	}
	if err := validateEffectHashSet(observation.Evidence, "evidence"); err != nil {
		return err
	}
	if err := validateEffectHashSet(observation.ErrorEvidence, "errorEvidence"); err != nil {
		return err
	}
	switch observation.Outcome {
	case "completed":
		if len(observation.Evidence) < 1 || len(observation.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: completed observation needs evidence and no error evidence", ErrInvalidEffectRecord)
		}
	default:
		if len(observation.ErrorEvidence) < 1 {
			return fmt.Errorf("%w: %s observation needs error evidence", ErrInvalidEffectRecord, observation.Outcome)
		}
	}
	return nil
}
