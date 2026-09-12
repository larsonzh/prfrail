package adapters

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const agentRunnerRequestDomain = "proofrail:agent-runner-request:1\n"

var (
	ErrInvalidAgentRunnerRequest  = errors.New("invalid AgentRunner request")
	ErrAgentRunnerRequestConflict = errors.New("AgentRunner request conflict")
)

type AgentRunnerRequest struct {
	Kind                string   `json:"kind"`
	RequestID           string   `json:"requestId"`
	RunID               string   `json:"runId"`
	TaskID              string   `json:"taskId"`
	StepID              string   `json:"stepId"`
	Attempt             int      `json:"attempt"`
	CreatedAt           string   `json:"createdAt"`
	AdapterID           string   `json:"adapterId"`
	Mode                string   `json:"mode"`
	WorkspaceHash       string   `json:"workspaceHash"`
	ContextHash         string   `json:"contextHash"`
	ParentSnapshotHash  string   `json:"parentSnapshotHash"`
	AuthorizationHash   string   `json:"authorizationHash"`
	BudgetHash          string   `json:"budgetHash"`
	AllowedTargets      []string `json:"allowedTargets"`
	AllowedEffects      []string `json:"allowedEffects"`
	PriorSessionID      string   `json:"priorSessionId,omitempty"`
	PriorCompletionHash string   `json:"priorCompletionHash,omitempty"`
	Evidence            []string `json:"evidence"`
}

type AgentRunnerRequestRecord struct {
	SchemaVersion string             `json:"schemaVersion"`
	Request       AgentRunnerRequest `json:"request"`
	RecordHash    string             `json:"recordHash"`
}

type AgentRunnerRequestIndex struct {
	mu     sync.Mutex
	hashes map[string]string
}

func NewAgentRunnerRequestRecord(request AgentRunnerRequest) (AgentRunnerRequestRecord, error) {
	if request.Kind == "" {
		request.Kind = "agent-runner-request"
	}
	request.AllowedTargets = sortedClone(request.AllowedTargets)
	request.AllowedEffects = sortedClone(request.AllowedEffects)
	request.Evidence = sortedClone(request.Evidence)
	if err := validateAgentRunnerRequest(request); err != nil {
		return AgentRunnerRequestRecord{}, err
	}
	hash, err := digestMessage(agentRunnerRequestDomain, request)
	if err != nil {
		return AgentRunnerRequestRecord{}, err
	}
	return AgentRunnerRequestRecord{SchemaVersion: evidence.SchemaVersion, Request: request, RecordHash: hash}, nil
}

func DecodeAgentRunnerRequestRecord(input []byte) (AgentRunnerRequestRecord, error) {
	var record AgentRunnerRequestRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AgentRunnerRequestRecord{}, err
	}
	if err := ValidateAgentRunnerRequestRecord(record); err != nil {
		return AgentRunnerRequestRecord{}, err
	}
	return record, nil
}

func ValidateAgentRunnerRequestRecord(record AgentRunnerRequestRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAgentRunnerRequest)
	}
	if err := validateAgentRunnerRequest(record.Request); err != nil {
		return err
	}
	expected, err := digestMessage(agentRunnerRequestDomain, record.Request)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAgentRunnerRequest)
	}
	return nil
}

func NewAgentRunnerRequestIndex(records []AgentRunnerRequestRecord) (*AgentRunnerRequestIndex, error) {
	index := &AgentRunnerRequestIndex{hashes: make(map[string]string, len(records))}
	for _, record := range records {
		if _, err := index.Record(record); err != nil {
			return nil, err
		}
	}
	return index, nil
}

func (index *AgentRunnerRequestIndex) Record(record AgentRunnerRequestRecord) (bool, error) {
	if index == nil {
		return false, fmt.Errorf("%w: nil replay index", ErrInvalidAgentRunnerRequest)
	}
	if err := ValidateAgentRunnerRequestRecord(record); err != nil {
		return false, err
	}
	index.mu.Lock()
	defer index.mu.Unlock()
	if index.hashes == nil {
		index.hashes = make(map[string]string)
	}
	if existing, found := index.hashes[record.Request.RequestID]; found {
		if existing != record.RecordHash {
			return false, fmt.Errorf("%w: requestId %s", ErrAgentRunnerRequestConflict, record.Request.RequestID)
		}
		return true, nil
	}
	index.hashes[record.Request.RequestID] = record.RecordHash
	return false, nil
}

func validateAgentRunnerRequest(request AgentRunnerRequest) error {
	if request.Kind != "agent-runner-request" {
		return fmt.Errorf("%w: invalid kind", ErrInvalidAgentRunnerRequest)
	}
	if !evidence.ValidID(request.RequestID) || !evidence.ValidID(request.RunID) || !evidence.ValidID(request.TaskID) || !evidence.ValidID(request.StepID) || !evidence.ValidID(request.AdapterID) || request.Attempt < 1 {
		return fmt.Errorf("%w: invalid identity", ErrInvalidAgentRunnerRequest)
	}
	if _, err := time.Parse(TimestampLayout, request.CreatedAt); err != nil {
		return fmt.Errorf("%w: invalid timestamp", ErrInvalidAgentRunnerRequest)
	}
	if !validAgentRunnerHashes(request) {
		return fmt.Errorf("%w: invalid immutable binding", ErrInvalidAgentRunnerRequest)
	}
	if !validSortedUniqueIDs(request.AllowedTargets) || !validSortedUniqueIDs(request.AllowedEffects) || !validSortedUniqueHashes(request.Evidence) {
		return fmt.Errorf("%w: invalid target, effect, or evidence set", ErrInvalidAgentRunnerRequest)
	}
	switch request.Mode {
	case "create":
		if request.PriorSessionID != "" || request.PriorCompletionHash != "" {
			return fmt.Errorf("%w: create cannot bind prior session", ErrInvalidAgentRunnerRequest)
		}
	case "resume":
		if !evidence.ValidID(request.PriorSessionID) || !evidence.ValidHash(request.PriorCompletionHash) {
			return fmt.Errorf("%w: resume requires prior session and completion", ErrInvalidAgentRunnerRequest)
		}
	default:
		return fmt.Errorf("%w: invalid mode", ErrInvalidAgentRunnerRequest)
	}
	return nil
}

func validAgentRunnerHashes(request AgentRunnerRequest) bool {
	return evidence.ValidHash(request.WorkspaceHash) &&
		evidence.ValidHash(request.ContextHash) &&
		evidence.ValidHash(request.ParentSnapshotHash) &&
		evidence.ValidHash(request.AuthorizationHash) &&
		evidence.ValidHash(request.BudgetHash)
}

func sortedClone(values []string) []string {
	cloned := slices.Clone(values)
	slices.Sort(cloned)
	return cloned
}

func validSortedUniqueIDs(values []string) bool {
	if len(values) == 0 || !slices.IsSorted(values) {
		return false
	}
	for index, value := range values {
		if !evidence.ValidID(value) || index > 0 && values[index-1] == value {
			return false
		}
	}
	return true
}

func validSortedUniqueHashes(values []string) bool {
	if len(values) == 0 || !slices.IsSorted(values) {
		return false
	}
	for index, value := range values {
		if !evidence.ValidHash(value) || index > 0 && values[index-1] == value {
			return false
		}
	}
	return true
}
