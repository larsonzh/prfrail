package adapters

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const agentRunnerCompletionDomain = "proofrail:agent-runner-completion:1\n"

var (
	ErrInvalidAgentRunnerCompletion  = errors.New("invalid AgentRunner completion")
	ErrAgentRunnerCompletionConflict = errors.New("AgentRunner completion conflict")
)

type AgentRunnerCompletion struct {
	Kind               string   `json:"kind"`
	CompletionID       string   `json:"completionId"`
	RequestID          string   `json:"requestId"`
	RequestHash        string   `json:"requestHash"`
	RunID              string   `json:"runId"`
	TaskID             string   `json:"taskId"`
	StepID             string   `json:"stepId"`
	Attempt            int      `json:"attempt"`
	AdapterID          string   `json:"adapterId"`
	SessionID          string   `json:"sessionId"`
	CompletedAt        string   `json:"completedAt"`
	Status             string   `json:"status"`
	ExitCode           *int     `json:"exitCode"`
	ProcessTreeStatus  string   `json:"processTreeStatus"`
	LogsComplete       bool     `json:"logsComplete"`
	UsageComplete      bool     `json:"usageComplete"`
	OutputManifestHash *string  `json:"outputManifestHash"`
	Evidence           []string `json:"evidence"`
	ErrorEvidence      []string `json:"errorEvidence"`
}

type AgentRunnerCompletionRecord struct {
	SchemaVersion string                `json:"schemaVersion"`
	Completion    AgentRunnerCompletion `json:"completion"`
	RecordHash    string                `json:"recordHash"`
}

type AgentRunnerCompletionIndex struct {
	mu               sync.Mutex
	completionHashes map[string]string
	requestHashes    map[string]string
}

func NewAgentRunnerCompletionRecord(completion AgentRunnerCompletion) (AgentRunnerCompletionRecord, error) {
	if completion.Kind == "" {
		completion.Kind = "agent-runner-completion"
	}
	completion.Evidence = sortedClone(completion.Evidence)
	completion.ErrorEvidence = sortedClone(completion.ErrorEvidence)
	if err := validateAgentRunnerCompletion(completion); err != nil {
		return AgentRunnerCompletionRecord{}, err
	}
	hash, err := digestMessage(agentRunnerCompletionDomain, completion)
	if err != nil {
		return AgentRunnerCompletionRecord{}, err
	}
	return AgentRunnerCompletionRecord{SchemaVersion: evidence.SchemaVersion, Completion: completion, RecordHash: hash}, nil
}

func DecodeAgentRunnerCompletionRecord(input []byte) (AgentRunnerCompletionRecord, error) {
	var record AgentRunnerCompletionRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AgentRunnerCompletionRecord{}, err
	}
	if err := ValidateAgentRunnerCompletionRecord(record); err != nil {
		return AgentRunnerCompletionRecord{}, err
	}
	return record, nil
}

func ValidateAgentRunnerCompletionRecord(record AgentRunnerCompletionRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAgentRunnerCompletion)
	}
	if err := validateAgentRunnerCompletion(record.Completion); err != nil {
		return err
	}
	expected, err := digestMessage(agentRunnerCompletionDomain, record.Completion)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAgentRunnerCompletion)
	}
	return nil
}

func ValidateAgentRunnerCompletionBinding(request AgentRunnerRequestRecord, completion AgentRunnerCompletionRecord) error {
	if err := ValidateAgentRunnerRequestRecord(request); err != nil {
		return err
	}
	if err := ValidateAgentRunnerCompletionRecord(completion); err != nil {
		return err
	}
	body := completion.Completion
	if body.RequestID != request.Request.RequestID || body.RequestHash != request.RecordHash || body.RunID != request.Request.RunID || body.TaskID != request.Request.TaskID || body.StepID != request.Request.StepID || body.Attempt != request.Request.Attempt || body.AdapterID != request.Request.AdapterID {
		return fmt.Errorf("%w: request binding mismatch", ErrInvalidAgentRunnerCompletion)
	}
	if request.Request.Mode == "resume" && body.SessionID != request.Request.PriorSessionID {
		return fmt.Errorf("%w: resumed session mismatch", ErrInvalidAgentRunnerCompletion)
	}
	return nil
}

func NewAgentRunnerCompletionIndex(records []AgentRunnerCompletionRecord) (*AgentRunnerCompletionIndex, error) {
	index := &AgentRunnerCompletionIndex{
		completionHashes: make(map[string]string, len(records)),
		requestHashes:    make(map[string]string, len(records)),
	}
	for _, record := range records {
		if _, err := index.Record(record); err != nil {
			return nil, err
		}
	}
	return index, nil
}

func (index *AgentRunnerCompletionIndex) Record(record AgentRunnerCompletionRecord) (bool, error) {
	if index == nil {
		return false, fmt.Errorf("%w: nil replay index", ErrInvalidAgentRunnerCompletion)
	}
	if err := ValidateAgentRunnerCompletionRecord(record); err != nil {
		return false, err
	}
	index.mu.Lock()
	defer index.mu.Unlock()
	if index.completionHashes == nil {
		index.completionHashes = make(map[string]string)
	}
	if index.requestHashes == nil {
		index.requestHashes = make(map[string]string)
	}
	if existing, found := index.completionHashes[record.Completion.CompletionID]; found {
		if existing != record.RecordHash {
			return false, fmt.Errorf("%w: completionId %s", ErrAgentRunnerCompletionConflict, record.Completion.CompletionID)
		}
		return true, nil
	}
	if existing, found := index.requestHashes[record.Completion.RequestID]; found && existing != record.RecordHash {
		return false, fmt.Errorf("%w: requestId %s already completed", ErrAgentRunnerCompletionConflict, record.Completion.RequestID)
	}
	index.completionHashes[record.Completion.CompletionID] = record.RecordHash
	index.requestHashes[record.Completion.RequestID] = record.RecordHash
	return false, nil
}

func validateAgentRunnerCompletion(completion AgentRunnerCompletion) error {
	if completion.Kind != "agent-runner-completion" || !evidence.ValidID(completion.CompletionID) || !evidence.ValidID(completion.RequestID) || !evidence.ValidID(completion.RunID) || !evidence.ValidID(completion.TaskID) || !evidence.ValidID(completion.StepID) || !evidence.ValidID(completion.AdapterID) || !evidence.ValidID(completion.SessionID) || completion.Attempt < 1 {
		return fmt.Errorf("%w: invalid identity", ErrInvalidAgentRunnerCompletion)
	}
	if !evidence.ValidHash(completion.RequestHash) || !validSortedUniqueHashes(completion.Evidence) || !validSortedUniqueHashesAllowEmpty(completion.ErrorEvidence) {
		return fmt.Errorf("%w: invalid hash or evidence", ErrInvalidAgentRunnerCompletion)
	}
	if _, err := time.Parse(TimestampLayout, completion.CompletedAt); err != nil {
		return fmt.Errorf("%w: invalid timestamp", ErrInvalidAgentRunnerCompletion)
	}
	switch completion.Status {
	case "completed":
		if completion.ExitCode == nil || *completion.ExitCode != 0 || completion.ProcessTreeStatus != "stopped" || !completion.LogsComplete || !completion.UsageComplete || completion.OutputManifestHash == nil || !evidence.ValidHash(*completion.OutputManifestHash) || len(completion.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: incomplete completed receipt", ErrInvalidAgentRunnerCompletion)
		}
	case "failed", "cancelled", "operator-action-required", "uncertain":
		if len(completion.ErrorEvidence) == 0 || completion.ProcessTreeStatus != "stopped" && completion.ProcessTreeStatus != "unknown" {
			return fmt.Errorf("%w: incomplete non-success receipt", ErrInvalidAgentRunnerCompletion)
		}
		if completion.OutputManifestHash != nil && !evidence.ValidHash(*completion.OutputManifestHash) {
			return fmt.Errorf("%w: invalid output manifest hash", ErrInvalidAgentRunnerCompletion)
		}
	default:
		return fmt.Errorf("%w: invalid status", ErrInvalidAgentRunnerCompletion)
	}
	return nil
}

func validSortedUniqueHashesAllowEmpty(values []string) bool {
	return len(values) == 0 || validSortedUniqueHashes(values)
}
