package evidence

import (
	"errors"
	"fmt"
	"time"
)

const (
	diagnosisRecordDomain    = "proofrail:effect-record:1\n"
	diagnosisTimestampLayout = "2006-01-02T15:04:05.000Z"
)

var ErrInvalidRecoveryDiagnosis = errors.New("invalid recovery diagnosis")

var (
	validDiagnosisProcessStates  = map[string]bool{"running": true, "stopped": true, "unknown": true}
	validDiagnosisAllowedActions = map[string]bool{
		"inspect": true, "controlled-stop": true, "reconcile": true, "abort": true, "resume": true,
	}
)

// DiagnosisProcess mirrors the process body of the recovery-diagnosis record.
type DiagnosisProcess struct {
	ProcessRef string   `json:"processRef"`
	State      string   `json:"state"`
	Evidence   []string `json:"evidence"`
}

// RecoveryDiagnosis mirrors the diagnosis body of effect-record.schema.json.
// It is read-only and redacted by contract: mutationsPerformed must be false.
type RecoveryDiagnosis struct {
	RecordID                 string             `json:"recordId"`
	Kind                     string             `json:"kind"`
	CreatedAt                string             `json:"createdAt"`
	RunID                    string             `json:"runId"`
	Mode                     string             `json:"mode"`
	LastAcceptedSnapshotHash *string            `json:"lastAcceptedSnapshotHash"`
	JournalEvidence          []string           `json:"journalEvidence"`
	Processes                []DiagnosisProcess `json:"processes"`
	UnknownEffectEvidence    []string           `json:"unknownEffectEvidence"`
	AllowedActions           []string           `json:"allowedActions"`
	MutationsPerformed       bool               `json:"mutationsPerformed"`
	Outcome                  string             `json:"outcome"`
	ErrorEvidence            []string           `json:"errorEvidence"`
}

// EffectDiagnosisRecord is the wire record wrapping a recovery diagnosis.
type EffectDiagnosisRecord struct {
	SchemaVersion string            `json:"schemaVersion"`
	Record        RecoveryDiagnosis `json:"record"`
	RecordHash    string            `json:"recordHash"`
}

// NewEffectDiagnosisRecord validates and hashes a recovery diagnosis.
func NewEffectDiagnosisRecord(diagnosis RecoveryDiagnosis) (EffectDiagnosisRecord, error) {
	if diagnosis.Kind == "" {
		diagnosis.Kind = "recovery-diagnosis"
	}
	if diagnosis.Mode == "" {
		diagnosis.Mode = "read-only-redacted"
	}
	if err := validateRecoveryDiagnosis(diagnosis); err != nil {
		return EffectDiagnosisRecord{}, err
	}
	canonical, err := EncodeCanonical(diagnosis)
	if err != nil {
		return EffectDiagnosisRecord{}, fmt.Errorf("%w: canonicalize diagnosis: %v", ErrInvalidRecoveryDiagnosis, err)
	}
	return EffectDiagnosisRecord{
		SchemaVersion: SchemaVersion,
		Record:        diagnosis,
		RecordHash:    Digest(diagnosisRecordDomain, canonical),
	}, nil
}

// DecodeEffectDiagnosisRecord strictly decodes and validates a wire record.
func DecodeEffectDiagnosisRecord(input []byte) (EffectDiagnosisRecord, error) {
	var record EffectDiagnosisRecord
	if err := decodeStrictJSON(input, &record); err != nil {
		return EffectDiagnosisRecord{}, err
	}
	if err := ValidateEffectDiagnosisRecord(record); err != nil {
		return EffectDiagnosisRecord{}, err
	}
	return record, nil
}

// ValidateEffectDiagnosisRecord re-validates body and record hash.
func ValidateEffectDiagnosisRecord(record EffectDiagnosisRecord) error {
	if record.SchemaVersion != SchemaVersion {
		return fmt.Errorf("%w: unsupported schemaVersion %q", ErrInvalidRecoveryDiagnosis, record.SchemaVersion)
	}
	if err := validateRecoveryDiagnosis(record.Record); err != nil {
		return err
	}
	canonical, err := EncodeCanonical(record.Record)
	if err != nil {
		return fmt.Errorf("%w: canonicalize diagnosis: %v", ErrInvalidRecoveryDiagnosis, err)
	}
	if record.RecordHash != Digest(diagnosisRecordDomain, canonical) {
		return fmt.Errorf("%w: diagnosis hash mismatch", ErrInvalidRecoveryDiagnosis)
	}
	return nil
}

func validateRecoveryDiagnosis(diagnosis RecoveryDiagnosis) error {
	if !ValidID(diagnosis.RecordID) || len(diagnosis.RecordID) > 64 {
		return fmt.Errorf("%w: invalid recordId", ErrInvalidRecoveryDiagnosis)
	}
	if diagnosis.Kind != "recovery-diagnosis" {
		return fmt.Errorf("%w: kind must be recovery-diagnosis", ErrInvalidRecoveryDiagnosis)
	}
	if _, err := time.Parse(diagnosisTimestampLayout, diagnosis.CreatedAt); err != nil {
		return fmt.Errorf("%w: invalid createdAt", ErrInvalidRecoveryDiagnosis)
	}
	if !ValidID(diagnosis.RunID) {
		return fmt.Errorf("%w: invalid runId", ErrInvalidRecoveryDiagnosis)
	}
	if diagnosis.Mode != "read-only-redacted" {
		return fmt.Errorf("%w: mode must be read-only-redacted", ErrInvalidRecoveryDiagnosis)
	}
	if diagnosis.LastAcceptedSnapshotHash != nil && !ValidHash(*diagnosis.LastAcceptedSnapshotHash) {
		return fmt.Errorf("%w: invalid lastAcceptedSnapshotHash", ErrInvalidRecoveryDiagnosis)
	}
	if err := validateDiagnosisHashSet(diagnosis.JournalEvidence, "journalEvidence"); err != nil {
		return err
	}
	for _, process := range diagnosis.Processes {
		if !ValidID(process.ProcessRef) || len(process.ProcessRef) > 64 {
			return fmt.Errorf("%w: invalid processRef", ErrInvalidRecoveryDiagnosis)
		}
		if !validDiagnosisProcessStates[process.State] {
			return fmt.Errorf("%w: unknown process state %q", ErrInvalidRecoveryDiagnosis, process.State)
		}
		if len(process.Evidence) < 1 {
			return fmt.Errorf("%w: process evidence must not be empty", ErrInvalidRecoveryDiagnosis)
		}
		if err := validateDiagnosisHashSet(process.Evidence, "process evidence"); err != nil {
			return err
		}
	}
	if err := validateDiagnosisHashSet(diagnosis.UnknownEffectEvidence, "unknownEffectEvidence"); err != nil {
		return err
	}
	seenActions := map[string]bool{}
	for _, action := range diagnosis.AllowedActions {
		if !validDiagnosisAllowedActions[action] {
			return fmt.Errorf("%w: unknown allowed action %q", ErrInvalidRecoveryDiagnosis, action)
		}
		if seenActions[action] {
			return fmt.Errorf("%w: duplicate allowed action %q", ErrInvalidRecoveryDiagnosis, action)
		}
		seenActions[action] = true
	}
	if diagnosis.MutationsPerformed {
		return fmt.Errorf("%w: diagnosis must not perform mutations", ErrInvalidRecoveryDiagnosis)
	}
	if err := validateDiagnosisHashSet(diagnosis.ErrorEvidence, "errorEvidence"); err != nil {
		return err
	}
	switch diagnosis.Outcome {
	case "clear":
		if len(diagnosis.UnknownEffectEvidence) != 0 || len(diagnosis.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: clear diagnosis needs empty unknown/error evidence", ErrInvalidRecoveryDiagnosis)
		}
	case "blocked":
		if len(diagnosis.ErrorEvidence) < 1 {
			return fmt.Errorf("%w: blocked diagnosis needs error evidence", ErrInvalidRecoveryDiagnosis)
		}
	default:
		return fmt.Errorf("%w: unknown outcome %q", ErrInvalidRecoveryDiagnosis, diagnosis.Outcome)
	}
	return nil
}

func validateDiagnosisHashSet(values []string, field string) error {
	seen := map[string]bool{}
	for _, value := range values {
		if !ValidHash(value) {
			return fmt.Errorf("%w: invalid %s hash", ErrInvalidRecoveryDiagnosis, field)
		}
		if seen[value] {
			return fmt.Errorf("%w: duplicate %s hash", ErrInvalidRecoveryDiagnosis, field)
		}
		seen[value] = true
	}
	return nil
}
