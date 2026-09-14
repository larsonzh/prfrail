package adapters

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	agentRunnerLaunchReceiptDomain  = "proofrail:agent-runner-launch-receipt:1\n"
	agentRunnerLaunchIdentityDomain = "proofrail:agent-runner-process-identity:1\n"
)

var (
	// ErrInvalidAgentRunnerLaunchReceipt reports an invalid launch-intent receipt.
	ErrInvalidAgentRunnerLaunchReceipt = errors.New("invalid AgentRunner launch receipt")
	// ErrInvalidAgentRunnerLaunchIdentity reports an invalid process-identity record.
	ErrInvalidAgentRunnerLaunchIdentity = errors.New("invalid AgentRunner launch identity")
	// ErrAgentRunnerLaunchReceiptConflict reports a launch slot already owned by
	// a different launch-intent record.
	ErrAgentRunnerLaunchReceiptConflict = errors.New("AgentRunner launch receipt conflict")
	// ErrAgentRunnerLaunchIdentityConflict reports an identity slot already owned
	// by a different process-identity record.
	ErrAgentRunnerLaunchIdentityConflict = errors.New("AgentRunner launch identity conflict")
)

// AgentRunnerLaunchReceipt is the launch-intent body published before spawning.
// It is store-local metadata: never on the wire and never part of any R/C
// recordHash.
type AgentRunnerLaunchReceipt struct {
	Kind              string `json:"kind"`
	LaunchID          string `json:"launchId"`
	RequestID         string `json:"requestId"`
	RequestHash       string `json:"requestHash"`
	RunID             string `json:"runId"`
	TaskID            string `json:"taskId"`
	StepID            string `json:"stepId"`
	Attempt           int    `json:"attempt"`
	AdapterID         string `json:"adapterId"`
	ReconfirmedAt     string `json:"reconfirmedAt"`
	AuthorizationHash string `json:"authorizationHash"`
	BudgetHash        string `json:"budgetHash"`
}

// AgentRunnerLaunchReceiptRecord wraps the launch-intent body with its local digest.
type AgentRunnerLaunchReceiptRecord struct {
	SchemaVersion string                   `json:"schemaVersion"`
	Receipt       AgentRunnerLaunchReceipt `json:"receipt"`
	RecordHash    string                   `json:"recordHash"`
}

// AgentRunnerLaunchIdentity is the process-identity body published after a spawn.
// It is store-local metadata: never on the wire and never part of any R/C recordHash.
type AgentRunnerLaunchIdentity struct {
	Kind             string `json:"kind"`
	LaunchID         string `json:"launchId"`
	RequestID        string `json:"requestId"`
	IntentRecordHash string `json:"intentRecordHash"`
	ProcessID        string `json:"processId"`
	StartedAt        string `json:"startedAt"`
}

// AgentRunnerLaunchIdentityRecord wraps the identity body with its local digest.
type AgentRunnerLaunchIdentityRecord struct {
	SchemaVersion string                    `json:"schemaVersion"`
	Identity      AgentRunnerLaunchIdentity `json:"identity"`
	RecordHash    string                    `json:"recordHash"`
}

// NewAgentRunnerLaunchReceiptRecord builds a validated launch-intent record.
func NewAgentRunnerLaunchReceiptRecord(receipt AgentRunnerLaunchReceipt) (AgentRunnerLaunchReceiptRecord, error) {
	if receipt.Kind == "" {
		receipt.Kind = "agent-runner-launch-receipt"
	}
	if err := validateAgentRunnerLaunchReceipt(receipt); err != nil {
		return AgentRunnerLaunchReceiptRecord{}, err
	}
	hash, err := digestMessage(agentRunnerLaunchReceiptDomain, receipt)
	if err != nil {
		return AgentRunnerLaunchReceiptRecord{}, err
	}
	return AgentRunnerLaunchReceiptRecord{SchemaVersion: evidence.SchemaVersion, Receipt: receipt, RecordHash: hash}, nil
}

// DecodeAgentRunnerLaunchReceiptRecord strictly decodes and validates a record.
func DecodeAgentRunnerLaunchReceiptRecord(input []byte) (AgentRunnerLaunchReceiptRecord, error) {
	var record AgentRunnerLaunchReceiptRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AgentRunnerLaunchReceiptRecord{}, err
	}
	if err := ValidateAgentRunnerLaunchReceiptRecord(record); err != nil {
		return AgentRunnerLaunchReceiptRecord{}, err
	}
	return record, nil
}

// ValidateAgentRunnerLaunchReceiptRecord validates schema, hashes and body.
func ValidateAgentRunnerLaunchReceiptRecord(record AgentRunnerLaunchReceiptRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAgentRunnerLaunchReceipt)
	}
	if err := validateAgentRunnerLaunchReceipt(record.Receipt); err != nil {
		return err
	}
	expected, err := digestMessage(agentRunnerLaunchReceiptDomain, record.Receipt)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAgentRunnerLaunchReceipt)
	}
	return nil
}

func validateAgentRunnerLaunchReceipt(receipt AgentRunnerLaunchReceipt) error {
	if receipt.Kind != "agent-runner-launch-receipt" {
		return fmt.Errorf("%w: invalid kind", ErrInvalidAgentRunnerLaunchReceipt)
	}
	if !evidence.ValidID(receipt.LaunchID) || !evidence.ValidID(receipt.RequestID) ||
		!evidence.ValidID(receipt.RunID) || !evidence.ValidID(receipt.TaskID) ||
		!evidence.ValidID(receipt.StepID) || !evidence.ValidID(receipt.AdapterID) {
		return fmt.Errorf("%w: invalid identity binding", ErrInvalidAgentRunnerLaunchReceipt)
	}
	if !evidence.ValidHash(receipt.RequestHash) || !evidence.ValidHash(receipt.AuthorizationHash) || !evidence.ValidHash(receipt.BudgetHash) {
		return fmt.Errorf("%w: invalid digest binding", ErrInvalidAgentRunnerLaunchReceipt)
	}
	if receipt.Attempt < 1 {
		return fmt.Errorf("%w: invalid attempt", ErrInvalidAgentRunnerLaunchReceipt)
	}
	if _, err := time.Parse(time.RFC3339, receipt.ReconfirmedAt); err != nil {
		return fmt.Errorf("%w: invalid reconfirmedAt", ErrInvalidAgentRunnerLaunchReceipt)
	}
	return nil
}

// NewAgentRunnerLaunchIdentityRecord builds a validated process-identity record.
func NewAgentRunnerLaunchIdentityRecord(identity AgentRunnerLaunchIdentity) (AgentRunnerLaunchIdentityRecord, error) {
	if identity.Kind == "" {
		identity.Kind = "agent-runner-process-identity"
	}
	if err := validateAgentRunnerLaunchIdentity(identity); err != nil {
		return AgentRunnerLaunchIdentityRecord{}, err
	}
	hash, err := digestMessage(agentRunnerLaunchIdentityDomain, identity)
	if err != nil {
		return AgentRunnerLaunchIdentityRecord{}, err
	}
	return AgentRunnerLaunchIdentityRecord{SchemaVersion: evidence.SchemaVersion, Identity: identity, RecordHash: hash}, nil
}

// DecodeAgentRunnerLaunchIdentityRecord strictly decodes and validates a record.
func DecodeAgentRunnerLaunchIdentityRecord(input []byte) (AgentRunnerLaunchIdentityRecord, error) {
	var record AgentRunnerLaunchIdentityRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AgentRunnerLaunchIdentityRecord{}, err
	}
	if err := ValidateAgentRunnerLaunchIdentityRecord(record); err != nil {
		return AgentRunnerLaunchIdentityRecord{}, err
	}
	return record, nil
}

// ValidateAgentRunnerLaunchIdentityRecord validates schema, hashes and body.
func ValidateAgentRunnerLaunchIdentityRecord(record AgentRunnerLaunchIdentityRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAgentRunnerLaunchIdentity)
	}
	if err := validateAgentRunnerLaunchIdentity(record.Identity); err != nil {
		return err
	}
	expected, err := digestMessage(agentRunnerLaunchIdentityDomain, record.Identity)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAgentRunnerLaunchIdentity)
	}
	return nil
}

func validateAgentRunnerLaunchIdentity(identity AgentRunnerLaunchIdentity) error {
	if identity.Kind != "agent-runner-process-identity" {
		return fmt.Errorf("%w: invalid kind", ErrInvalidAgentRunnerLaunchIdentity)
	}
	if !evidence.ValidID(identity.LaunchID) || !evidence.ValidID(identity.RequestID) {
		return fmt.Errorf("%w: invalid identity binding", ErrInvalidAgentRunnerLaunchIdentity)
	}
	if !evidence.ValidHash(identity.IntentRecordHash) {
		return fmt.Errorf("%w: invalid intent digest", ErrInvalidAgentRunnerLaunchIdentity)
	}
	if strings.TrimSpace(identity.ProcessID) == "" {
		return fmt.Errorf("%w: empty process identity", ErrInvalidAgentRunnerLaunchIdentity)
	}
	if _, err := time.Parse(time.RFC3339, identity.StartedAt); err != nil {
		return fmt.Errorf("%w: invalid startedAt", ErrInvalidAgentRunnerLaunchIdentity)
	}
	return nil
}
