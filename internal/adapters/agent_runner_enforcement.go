package adapters

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const agentRunnerEnforcementDomain = "proofrail:agent-runner-enforcement:1\n"

var (
	ErrInvalidAgentRunnerEnforcement  = errors.New("invalid AgentRunner enforcement report")
	ErrAgentRunnerEnforcementConflict = errors.New("AgentRunner enforcement report conflict")
	ErrAgentRunnerAdmissionBlocked    = errors.New("AgentRunner admission blocked")
)

type AgentRunnerEnforcementMatrix struct {
	ToolControl    AgentRunnerCapabilityFinding `json:"toolControl"`
	NetworkControl AgentRunnerCapabilityFinding `json:"networkControl"`
}

type AgentRunnerEnforcement struct {
	Kind                  string                       `json:"kind"`
	EnforcementID         string                       `json:"enforcementId"`
	ProviderID            string                       `json:"providerId"`
	ProbedAt              string                       `json:"probedAt"`
	CapabilityRecordHash  string                       `json:"capabilityRecordHash"`
	ExecutableHash        string                       `json:"executableHash"`
	ConfigHash            string                       `json:"configHash"`
	PlatformHash          string                       `json:"platformHash"`
	EnforcementConfigHash string                       `json:"enforcementConfigHash"`
	Disposition           string                       `json:"disposition"`
	Matrix                AgentRunnerEnforcementMatrix `json:"matrix"`
	Evidence              []string                     `json:"evidence"`
}

type AgentRunnerEnforcementRecord struct {
	SchemaVersion string                 `json:"schemaVersion"`
	Enforcement   AgentRunnerEnforcement `json:"enforcement"`
	RecordHash    string                 `json:"recordHash"`
}

type AgentRunnerEnforcementIndex struct {
	mu     sync.Mutex
	hashes map[string]string
}

type AgentRunnerAdmissionPolicy struct {
	OS           string
	Arch         string
	PlatformHash string
	EvaluatedAt  time.Time
	MaximumAge   time.Duration
}

func NewAgentRunnerEnforcementRecord(enforcement AgentRunnerEnforcement) (AgentRunnerEnforcementRecord, error) {
	if enforcement.Kind == "" {
		enforcement.Kind = "agent-runner-enforcement"
	}
	normalizeAgentRunnerEnforcement(&enforcement)
	if err := validateAgentRunnerEnforcement(enforcement); err != nil {
		return AgentRunnerEnforcementRecord{}, err
	}
	hash, err := digestMessage(agentRunnerEnforcementDomain, enforcement)
	if err != nil {
		return AgentRunnerEnforcementRecord{}, err
	}
	return AgentRunnerEnforcementRecord{SchemaVersion: evidence.SchemaVersion, Enforcement: enforcement, RecordHash: hash}, nil
}

func DecodeAgentRunnerEnforcementRecord(input []byte) (AgentRunnerEnforcementRecord, error) {
	var record AgentRunnerEnforcementRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AgentRunnerEnforcementRecord{}, err
	}
	if err := ValidateAgentRunnerEnforcementRecord(record); err != nil {
		return AgentRunnerEnforcementRecord{}, err
	}
	return record, nil
}

func ValidateAgentRunnerEnforcementRecord(record AgentRunnerEnforcementRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAgentRunnerEnforcement)
	}
	if err := validateAgentRunnerEnforcement(record.Enforcement); err != nil {
		return err
	}
	expected, err := digestMessage(agentRunnerEnforcementDomain, record.Enforcement)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAgentRunnerEnforcement)
	}
	return nil
}

func NewAgentRunnerEnforcementIndex(records []AgentRunnerEnforcementRecord) (*AgentRunnerEnforcementIndex, error) {
	index := &AgentRunnerEnforcementIndex{hashes: make(map[string]string, len(records))}
	for _, record := range records {
		if _, err := index.Record(record); err != nil {
			return nil, err
		}
	}
	return index, nil
}

func (index *AgentRunnerEnforcementIndex) Record(record AgentRunnerEnforcementRecord) (bool, error) {
	if index == nil {
		return false, fmt.Errorf("%w: nil replay index", ErrInvalidAgentRunnerEnforcement)
	}
	if err := ValidateAgentRunnerEnforcementRecord(record); err != nil {
		return false, err
	}
	index.mu.Lock()
	defer index.mu.Unlock()
	if index.hashes == nil {
		index.hashes = make(map[string]string)
	}
	key := record.Enforcement.EnforcementID
	if existing, found := index.hashes[key]; found {
		if existing != record.RecordHash {
			return false, fmt.Errorf("%w: enforcementId %s", ErrAgentRunnerEnforcementConflict, key)
		}
		return true, nil
	}
	index.hashes[key] = record.RecordHash
	return false, nil
}

func ValidateAgentRunnerAdmission(capability AgentRunnerCapabilityRecord, enforcement *AgentRunnerEnforcementRecord, policy AgentRunnerAdmissionPolicy) error {
	if err := ValidateAgentRunnerCapabilityRecord(capability); err != nil {
		return err
	}
	if !validAgentRunnerPlatform(policy.OS, policy.Arch) || !evidence.ValidHash(policy.PlatformHash) || policy.EvaluatedAt.IsZero() || policy.MaximumAge <= 0 {
		return fmt.Errorf("%w: invalid admission policy", ErrAgentRunnerAdmissionBlocked)
	}
	if capability.Capability.OS != policy.OS || capability.Capability.Arch != policy.Arch || capability.Capability.PlatformHash != policy.PlatformHash {
		return fmt.Errorf("%w: capability platform mismatch", ErrAgentRunnerAdmissionBlocked)
	}
	if capability.Capability.Disposition == "compatible" {
		return nil
	}
	for _, finding := range nonControlCapabilityFindings(&capability.Capability.Matrix) {
		if finding.Status != "verified" {
			return fmt.Errorf("%w: uncompensated capability", ErrAgentRunnerAdmissionBlocked)
		}
	}
	if enforcement == nil {
		return fmt.Errorf("%w: external enforcement required", ErrAgentRunnerAdmissionBlocked)
	}
	if err := ValidateAgentRunnerEnforcementRecord(*enforcement); err != nil {
		return fmt.Errorf("%w: %v", ErrAgentRunnerAdmissionBlocked, err)
	}
	boundary := enforcement.Enforcement
	if boundary.CapabilityRecordHash != capability.RecordHash || boundary.ExecutableHash != capability.Capability.ExecutableHash || boundary.ConfigHash != capability.Capability.ConfigHash || boundary.PlatformHash != capability.Capability.PlatformHash || boundary.PlatformHash != policy.PlatformHash {
		return fmt.Errorf("%w: enforcement binding mismatch", ErrAgentRunnerAdmissionBlocked)
	}
	probedAt, _ := time.Parse(TimestampLayout, boundary.ProbedAt)
	capabilityProbedAt, _ := time.Parse(TimestampLayout, capability.Capability.ProbedAt)
	if probedAt.Before(capabilityProbedAt) || probedAt.After(policy.EvaluatedAt) || policy.EvaluatedAt.Sub(probedAt) > policy.MaximumAge {
		return fmt.Errorf("%w: stale or future enforcement evidence", ErrAgentRunnerAdmissionBlocked)
	}
	toolVerified := capability.Capability.Matrix.ToolControl.Status == "verified" || boundary.Matrix.ToolControl.Status == "verified"
	networkVerified := capability.Capability.Matrix.NetworkControl.Status == "verified" || boundary.Matrix.NetworkControl.Status == "verified"
	if !toolVerified || !networkVerified {
		return fmt.Errorf("%w: incomplete external enforcement", ErrAgentRunnerAdmissionBlocked)
	}
	return nil
}

func normalizeAgentRunnerEnforcement(enforcement *AgentRunnerEnforcement) {
	enforcement.Evidence = sortedClone(enforcement.Evidence)
	enforcement.Matrix.ToolControl.Evidence = sortedClone(enforcement.Matrix.ToolControl.Evidence)
	enforcement.Matrix.NetworkControl.Evidence = sortedClone(enforcement.Matrix.NetworkControl.Evidence)
}

func validateAgentRunnerEnforcement(enforcement AgentRunnerEnforcement) error {
	if enforcement.Kind != "agent-runner-enforcement" || !evidence.ValidID(enforcement.EnforcementID) || !evidence.ValidID(enforcement.ProviderID) {
		return fmt.Errorf("%w: invalid identity", ErrInvalidAgentRunnerEnforcement)
	}
	if _, err := time.Parse(TimestampLayout, enforcement.ProbedAt); err != nil {
		return fmt.Errorf("%w: invalid timestamp", ErrInvalidAgentRunnerEnforcement)
	}
	if !evidence.ValidHash(enforcement.CapabilityRecordHash) || !evidence.ValidHash(enforcement.ExecutableHash) || !evidence.ValidHash(enforcement.ConfigHash) || !evidence.ValidHash(enforcement.PlatformHash) || !evidence.ValidHash(enforcement.EnforcementConfigHash) || !validSortedUniqueHashes(enforcement.Evidence) {
		return fmt.Errorf("%w: invalid binding or evidence", ErrInvalidAgentRunnerEnforcement)
	}
	toolVerified := validAgentRunnerCapabilityFinding(enforcement.Matrix.ToolControl) && enforcement.Matrix.ToolControl.Status == "verified"
	networkVerified := validAgentRunnerCapabilityFinding(enforcement.Matrix.NetworkControl) && enforcement.Matrix.NetworkControl.Status == "verified"
	if !validAgentRunnerCapabilityFinding(enforcement.Matrix.ToolControl) || !validAgentRunnerCapabilityFinding(enforcement.Matrix.NetworkControl) {
		return fmt.Errorf("%w: invalid matrix finding", ErrInvalidAgentRunnerEnforcement)
	}
	allVerified := toolVerified && networkVerified
	if enforcement.Disposition == "verified" && !allVerified || enforcement.Disposition == "blocked" && allVerified {
		return fmt.Errorf("%w: disposition does not match matrix", ErrInvalidAgentRunnerEnforcement)
	}
	if enforcement.Disposition != "verified" && enforcement.Disposition != "blocked" {
		return fmt.Errorf("%w: invalid disposition", ErrInvalidAgentRunnerEnforcement)
	}
	return nil
}

func nonControlCapabilityFindings(matrix *AgentRunnerCapabilityMatrix) []*AgentRunnerCapabilityFinding {
	return []*AgentRunnerCapabilityFinding{
		&matrix.Noninteractive,
		&matrix.CWD,
		&matrix.EventStreamOrCompleteLogs,
		&matrix.SessionCreate,
		&matrix.SessionResume,
		&matrix.Cancellation,
		&matrix.ProcessTreeStop,
		&matrix.PermissionControl,
		&matrix.Usage,
		&matrix.UnattendedConfirmations,
	}
}
