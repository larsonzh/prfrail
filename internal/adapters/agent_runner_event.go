package adapters

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const agentRunnerEventDomain = "proofrail:agent-runner-event:1\n"

var (
	ErrInvalidAgentRunnerEvent  = errors.New("invalid AgentRunner event")
	ErrAgentRunnerEventConflict = errors.New("AgentRunner event conflict")
)

type AgentRunnerEvent struct {
	Kind        string   `json:"kind"`
	EventID     string   `json:"eventId"`
	RequestID   string   `json:"requestId"`
	RequestHash string   `json:"requestHash"`
	RunID       string   `json:"runId"`
	TaskID      string   `json:"taskId"`
	StepID      string   `json:"stepId"`
	Attempt     int      `json:"attempt"`
	AdapterID   string   `json:"adapterId"`
	SessionID   string   `json:"sessionId"`
	Sequence    int      `json:"sequence"`
	OccurredAt  string   `json:"occurredAt"`
	EventType   string   `json:"eventType"`
	PayloadHash string   `json:"payloadHash"`
	Evidence    []string `json:"evidence"`
}

type AgentRunnerEventRecord struct {
	SchemaVersion string           `json:"schemaVersion"`
	Event         AgentRunnerEvent `json:"event"`
	RecordHash    string           `json:"recordHash"`
}

type AgentRunnerEventIndex struct {
	mu              sync.Mutex
	eventHashes     map[string]string
	sequenceHashes  map[string]string
	requestSessions map[string]string
}

func NewAgentRunnerEventRecord(event AgentRunnerEvent) (AgentRunnerEventRecord, error) {
	if event.Kind == "" {
		event.Kind = "agent-runner-event"
	}
	event.Evidence = sortedClone(event.Evidence)
	if err := validateAgentRunnerEvent(event); err != nil {
		return AgentRunnerEventRecord{}, err
	}
	hash, err := digestMessage(agentRunnerEventDomain, event)
	if err != nil {
		return AgentRunnerEventRecord{}, err
	}
	return AgentRunnerEventRecord{SchemaVersion: evidence.SchemaVersion, Event: event, RecordHash: hash}, nil
}

func DecodeAgentRunnerEventRecord(input []byte) (AgentRunnerEventRecord, error) {
	var record AgentRunnerEventRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AgentRunnerEventRecord{}, err
	}
	if err := ValidateAgentRunnerEventRecord(record); err != nil {
		return AgentRunnerEventRecord{}, err
	}
	return record, nil
}

func ValidateAgentRunnerEventRecord(record AgentRunnerEventRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAgentRunnerEvent)
	}
	if err := validateAgentRunnerEvent(record.Event); err != nil {
		return err
	}
	expected, err := digestMessage(agentRunnerEventDomain, record.Event)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAgentRunnerEvent)
	}
	return nil
}

func ValidateAgentRunnerEventBinding(request AgentRunnerRequestRecord, event AgentRunnerEventRecord) error {
	if err := ValidateAgentRunnerRequestRecord(request); err != nil {
		return err
	}
	if err := ValidateAgentRunnerEventRecord(event); err != nil {
		return err
	}
	body := event.Event
	if body.RequestID != request.Request.RequestID || body.RequestHash != request.RecordHash || body.RunID != request.Request.RunID || body.TaskID != request.Request.TaskID || body.StepID != request.Request.StepID || body.Attempt != request.Request.Attempt || body.AdapterID != request.Request.AdapterID {
		return fmt.Errorf("%w: request binding mismatch", ErrInvalidAgentRunnerEvent)
	}
	if request.Request.Mode == "resume" && body.SessionID != request.Request.PriorSessionID {
		return fmt.Errorf("%w: resumed session mismatch", ErrInvalidAgentRunnerEvent)
	}
	return nil
}

func NewAgentRunnerEventIndex(records []AgentRunnerEventRecord) (*AgentRunnerEventIndex, error) {
	index := &AgentRunnerEventIndex{
		eventHashes:     make(map[string]string, len(records)),
		sequenceHashes:  make(map[string]string, len(records)),
		requestSessions: make(map[string]string, len(records)),
	}
	for _, record := range records {
		if _, err := index.Record(record); err != nil {
			return nil, err
		}
	}
	return index, nil
}

func (index *AgentRunnerEventIndex) Record(record AgentRunnerEventRecord) (bool, error) {
	if index == nil {
		return false, fmt.Errorf("%w: nil replay index", ErrInvalidAgentRunnerEvent)
	}
	if err := ValidateAgentRunnerEventRecord(record); err != nil {
		return false, err
	}
	index.mu.Lock()
	defer index.mu.Unlock()
	if index.eventHashes == nil {
		index.eventHashes = make(map[string]string)
	}
	if index.sequenceHashes == nil {
		index.sequenceHashes = make(map[string]string)
	}
	if index.requestSessions == nil {
		index.requestSessions = make(map[string]string)
	}
	if existing, found := index.eventHashes[record.Event.EventID]; found {
		if existing != record.RecordHash {
			return false, fmt.Errorf("%w: eventId %s", ErrAgentRunnerEventConflict, record.Event.EventID)
		}
		return true, nil
	}
	sequenceKey := fmt.Sprintf("%s/%d", record.Event.SessionID, record.Event.Sequence)
	if existing, found := index.sequenceHashes[sequenceKey]; found && existing != record.RecordHash {
		return false, fmt.Errorf("%w: session sequence %s", ErrAgentRunnerEventConflict, sequenceKey)
	}
	if existing, found := index.requestSessions[record.Event.RequestID]; found && existing != record.Event.SessionID {
		return false, fmt.Errorf("%w: requestId %s changed session", ErrAgentRunnerEventConflict, record.Event.RequestID)
	}
	index.eventHashes[record.Event.EventID] = record.RecordHash
	index.sequenceHashes[sequenceKey] = record.RecordHash
	index.requestSessions[record.Event.RequestID] = record.Event.SessionID
	return false, nil
}

func validateAgentRunnerEvent(event AgentRunnerEvent) error {
	if event.Kind != "agent-runner-event" || !evidence.ValidID(event.EventID) || !evidence.ValidID(event.RequestID) || !evidence.ValidID(event.RunID) || !evidence.ValidID(event.TaskID) || !evidence.ValidID(event.StepID) || !evidence.ValidID(event.AdapterID) || !evidence.ValidID(event.SessionID) || event.Attempt < 1 || event.Sequence < 1 {
		return fmt.Errorf("%w: invalid identity or sequence", ErrInvalidAgentRunnerEvent)
	}
	if !evidence.ValidHash(event.RequestHash) || !evidence.ValidHash(event.PayloadHash) || !validSortedUniqueHashes(event.Evidence) {
		return fmt.Errorf("%w: invalid hash or evidence", ErrInvalidAgentRunnerEvent)
	}
	if _, err := time.Parse(TimestampLayout, event.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid timestamp", ErrInvalidAgentRunnerEvent)
	}
	switch event.EventType {
	case "session-created", "output", "tool", "usage", "diagnostic", "operator-action-required", "process":
		return nil
	default:
		return fmt.Errorf("%w: invalid event type", ErrInvalidAgentRunnerEvent)
	}
}
