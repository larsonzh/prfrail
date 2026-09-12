package adapters

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const agentRunnerCapabilityDomain = "proofrail:agent-runner-capability:1\n"

var (
	ErrInvalidAgentRunnerCapability  = errors.New("invalid AgentRunner capability report")
	ErrAgentRunnerCapabilityConflict = errors.New("AgentRunner capability report conflict")
)

type AgentRunnerCapabilityFinding struct {
	Status   string   `json:"status"`
	Evidence []string `json:"evidence"`
}

type AgentRunnerCapabilityMatrix struct {
	Noninteractive            AgentRunnerCapabilityFinding `json:"noninteractive"`
	CWD                       AgentRunnerCapabilityFinding `json:"cwd"`
	EventStreamOrCompleteLogs AgentRunnerCapabilityFinding `json:"eventStreamOrCompleteLogs"`
	SessionCreate             AgentRunnerCapabilityFinding `json:"sessionCreate"`
	SessionResume             AgentRunnerCapabilityFinding `json:"sessionResume"`
	Cancellation              AgentRunnerCapabilityFinding `json:"cancellation"`
	ProcessTreeStop           AgentRunnerCapabilityFinding `json:"processTreeStop"`
	ToolControl               AgentRunnerCapabilityFinding `json:"toolControl"`
	NetworkControl            AgentRunnerCapabilityFinding `json:"networkControl"`
	PermissionControl         AgentRunnerCapabilityFinding `json:"permissionControl"`
	Usage                     AgentRunnerCapabilityFinding `json:"usage"`
	UnattendedConfirmations   AgentRunnerCapabilityFinding `json:"unattendedConfirmations"`
}

type AgentRunnerCapability struct {
	Kind           string                      `json:"kind"`
	ProbeID        string                      `json:"probeId"`
	AdapterID      string                      `json:"adapterId"`
	ProbedAt       string                      `json:"probedAt"`
	OS             string                      `json:"os"`
	Arch           string                      `json:"arch"`
	PlatformHash   string                      `json:"platformHash"`
	ExecutableHash string                      `json:"executableHash"`
	Version        string                      `json:"version"`
	ConfigHash     string                      `json:"configHash"`
	Disposition    string                      `json:"disposition"`
	Matrix         AgentRunnerCapabilityMatrix `json:"matrix"`
	Evidence       []string                    `json:"evidence"`
}

type AgentRunnerCapabilityRecord struct {
	SchemaVersion string                `json:"schemaVersion"`
	Capability    AgentRunnerCapability `json:"capability"`
	RecordHash    string                `json:"recordHash"`
}

type AgentRunnerCapabilityIndex struct {
	mu     sync.Mutex
	hashes map[string]string
}

func NewAgentRunnerCapabilityRecord(capability AgentRunnerCapability) (AgentRunnerCapabilityRecord, error) {
	if capability.Kind == "" {
		capability.Kind = "agent-runner-capability"
	}
	normalizeAgentRunnerCapability(&capability)
	if err := validateAgentRunnerCapability(capability); err != nil {
		return AgentRunnerCapabilityRecord{}, err
	}
	hash, err := digestMessage(agentRunnerCapabilityDomain, capability)
	if err != nil {
		return AgentRunnerCapabilityRecord{}, err
	}
	return AgentRunnerCapabilityRecord{SchemaVersion: evidence.SchemaVersion, Capability: capability, RecordHash: hash}, nil
}

func DecodeAgentRunnerCapabilityRecord(input []byte) (AgentRunnerCapabilityRecord, error) {
	var record AgentRunnerCapabilityRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AgentRunnerCapabilityRecord{}, err
	}
	if err := ValidateAgentRunnerCapabilityRecord(record); err != nil {
		return AgentRunnerCapabilityRecord{}, err
	}
	return record, nil
}

func ValidateAgentRunnerCapabilityRecord(record AgentRunnerCapabilityRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAgentRunnerCapability)
	}
	if err := validateAgentRunnerCapability(record.Capability); err != nil {
		return err
	}
	expected, err := digestMessage(agentRunnerCapabilityDomain, record.Capability)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAgentRunnerCapability)
	}
	return nil
}

func NewAgentRunnerCapabilityIndex(records []AgentRunnerCapabilityRecord) (*AgentRunnerCapabilityIndex, error) {
	index := &AgentRunnerCapabilityIndex{hashes: make(map[string]string, len(records))}
	for _, record := range records {
		if _, err := index.Record(record); err != nil {
			return nil, err
		}
	}
	return index, nil
}

func (index *AgentRunnerCapabilityIndex) Record(record AgentRunnerCapabilityRecord) (bool, error) {
	if index == nil {
		return false, fmt.Errorf("%w: nil replay index", ErrInvalidAgentRunnerCapability)
	}
	if err := ValidateAgentRunnerCapabilityRecord(record); err != nil {
		return false, err
	}
	index.mu.Lock()
	defer index.mu.Unlock()
	if index.hashes == nil {
		index.hashes = make(map[string]string)
	}
	key := record.Capability.ProbeID
	if existing, found := index.hashes[key]; found {
		if existing != record.RecordHash {
			return false, fmt.Errorf("%w: probeId %s", ErrAgentRunnerCapabilityConflict, key)
		}
		return true, nil
	}
	index.hashes[key] = record.RecordHash
	return false, nil
}

func normalizeAgentRunnerCapability(capability *AgentRunnerCapability) {
	capability.Evidence = sortedClone(capability.Evidence)
	for _, finding := range agentRunnerCapabilityFindings(&capability.Matrix) {
		finding.Evidence = sortedClone(finding.Evidence)
	}
}

func validateAgentRunnerCapability(capability AgentRunnerCapability) error {
	if capability.Kind != "agent-runner-capability" || !evidence.ValidID(capability.ProbeID) || !evidence.ValidID(capability.AdapterID) {
		return fmt.Errorf("%w: invalid identity", ErrInvalidAgentRunnerCapability)
	}
	if _, err := time.Parse(TimestampLayout, capability.ProbedAt); err != nil {
		return fmt.Errorf("%w: invalid timestamp", ErrInvalidAgentRunnerCapability)
	}
	if !validAgentRunnerPlatform(capability.OS, capability.Arch) || !evidence.ValidHash(capability.PlatformHash) || !evidence.ValidHash(capability.ExecutableHash) || !evidence.ValidHash(capability.ConfigHash) || capability.Version == "" || len(capability.Version) > 128 || !validSortedUniqueHashes(capability.Evidence) {
		return fmt.Errorf("%w: invalid platform, executable, config, version, or evidence", ErrInvalidAgentRunnerCapability)
	}
	allVerified := true
	for _, finding := range agentRunnerCapabilityFindings(&capability.Matrix) {
		if !validAgentRunnerCapabilityFinding(*finding) {
			return fmt.Errorf("%w: invalid matrix finding", ErrInvalidAgentRunnerCapability)
		}
		allVerified = allVerified && finding.Status == "verified"
	}
	if capability.Disposition == "compatible" && !allVerified || capability.Disposition == "blocked" && allVerified {
		return fmt.Errorf("%w: disposition does not match matrix", ErrInvalidAgentRunnerCapability)
	}
	if capability.Disposition != "compatible" && capability.Disposition != "blocked" {
		return fmt.Errorf("%w: invalid disposition", ErrInvalidAgentRunnerCapability)
	}
	return nil
}

func validAgentRunnerPlatform(osName, arch string) bool {
	validOS := osName == "windows" || osName == "linux" || osName == "darwin"
	validArch := arch == "amd64" || arch == "arm64"
	return validOS && validArch
}

func validAgentRunnerCapabilityFinding(finding AgentRunnerCapabilityFinding) bool {
	if finding.Status != "verified" && finding.Status != "unsupported" && finding.Status != "unknown" {
		return false
	}
	return validSortedUniqueHashes(finding.Evidence)
}

func agentRunnerCapabilityFindings(matrix *AgentRunnerCapabilityMatrix) []*AgentRunnerCapabilityFinding {
	return []*AgentRunnerCapabilityFinding{
		&matrix.Noninteractive,
		&matrix.CWD,
		&matrix.EventStreamOrCompleteLogs,
		&matrix.SessionCreate,
		&matrix.SessionResume,
		&matrix.Cancellation,
		&matrix.ProcessTreeStop,
		&matrix.ToolControl,
		&matrix.NetworkControl,
		&matrix.PermissionControl,
		&matrix.Usage,
		&matrix.UnattendedConfirmations,
	}
}
