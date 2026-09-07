package evidence

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"time"
)

const (
	SchemaVersion    = "1.0.0"
	stateEventDomain = "proofrail:state-event:1\n"
)

var (
	ErrInvalidRecord      = errors.New("invalid evidence record")
	ErrUnsupportedVersion = errors.New("unsupported evidence schema version")
	ErrBrokenEventChain   = errors.New("broken state event chain")
	idPattern             = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)
	hashPattern           = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

type Actor struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Entity struct {
	Kind    string `json:"kind"`
	RunID   string `json:"runId"`
	TaskID  string `json:"taskId,omitempty"`
	StepID  string `json:"stepId,omitempty"`
	Attempt int    `json:"attempt,omitempty"`
}

type Reason struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

type Event struct {
	EventID           string   `json:"eventId"`
	RunID             string   `json:"runId"`
	Sequence          int      `json:"sequence"`
	OccurredAt        string   `json:"occurredAt"`
	Actor             Actor    `json:"actor"`
	Entity            Entity   `json:"entity"`
	FromState         string   `json:"fromState"`
	ToState           string   `json:"toState"`
	PreviousEventHash *string  `json:"previousEventHash"`
	InputEvidence     []string `json:"inputEvidence"`
	Reason            Reason   `json:"reason"`
}

type StateEvent struct {
	SchemaVersion string `json:"schemaVersion"`
	Event         Event  `json:"event"`
	EventHash     string `json:"eventHash"`
}

func NewStateEvent(event Event) (StateEvent, error) {
	if err := validateEvent(event); err != nil {
		return StateEvent{}, err
	}
	canonical, err := canonicalValue(event)
	if err != nil {
		return StateEvent{}, err
	}
	return StateEvent{SchemaVersion: SchemaVersion, Event: event, EventHash: Digest(stateEventDomain, canonical)}, nil
}

func DecodeStateEvent(input []byte) (StateEvent, error) {
	var record StateEvent
	if err := decodeStrictJSON(input, &record); err != nil {
		return StateEvent{}, err
	}
	if record.SchemaVersion != SchemaVersion {
		return StateEvent{}, fmt.Errorf("%w: %q", ErrUnsupportedVersion, record.SchemaVersion)
	}
	canonical, err := canonicalValue(record.Event)
	if err != nil {
		return StateEvent{}, err
	}
	if err := validateEvent(record.Event); err != nil {
		return StateEvent{}, err
	}
	if record.EventHash != Digest(stateEventDomain, canonical) {
		return StateEvent{}, fmt.Errorf("%w: event hash mismatch", ErrInvalidRecord)
	}
	return record, nil
}

func VerifyEventChain(records []StateEvent) error {
	if len(records) == 0 {
		return fmt.Errorf("%w: empty chain", ErrBrokenEventChain)
	}
	seenIDs := make(map[string]struct{}, len(records))
	states := make(map[string]string)
	for index, record := range records {
		if record.SchemaVersion != SchemaVersion {
			return fmt.Errorf("%w: %q", ErrUnsupportedVersion, record.SchemaVersion)
		}
		if err := validateEvent(record.Event); err != nil {
			return fmt.Errorf("%w at index %d: %v", ErrBrokenEventChain, index, err)
		}
		if record.Event.Sequence != index+1 {
			return fmt.Errorf("%w: sequence %d at index %d", ErrBrokenEventChain, record.Event.Sequence, index)
		}
		if _, exists := seenIDs[record.Event.EventID]; exists {
			return fmt.Errorf("%w: duplicate event ID %q", ErrBrokenEventChain, record.Event.EventID)
		}
		seenIDs[record.Event.EventID] = struct{}{}
		canonical, err := canonicalValue(record.Event)
		if err != nil {
			return err
		}
		expectedHash := Digest(stateEventDomain, canonical)
		if record.EventHash != expectedHash {
			return fmt.Errorf("%w: event %q hash mismatch", ErrBrokenEventChain, record.Event.EventID)
		}
		if index == 0 {
			if record.Event.PreviousEventHash != nil {
				return fmt.Errorf("%w: genesis has previous hash", ErrBrokenEventChain)
			}
		} else {
			previous := records[index-1]
			if record.Event.RunID != previous.Event.RunID || record.Event.PreviousEventHash == nil || *record.Event.PreviousEventHash != previous.EventHash {
				return fmt.Errorf("%w: event %q does not reference its predecessor", ErrBrokenEventChain, record.Event.EventID)
			}
		}
		entityKey := eventEntityKey(record.Event.Entity)
		expectedState := states[entityKey]
		if expectedState == "" {
			expectedState = "NONE"
		}
		if record.Event.FromState != expectedState {
			return fmt.Errorf("%w: event %q state continuity mismatch", ErrBrokenEventChain, record.Event.EventID)
		}
		states[entityKey] = record.Event.ToState
	}
	return nil
}

func validateEvent(event Event) error {
	if !validID(event.EventID) || !validID(event.RunID) || !validID(event.Actor.ID) || !validID(event.Reason.Code) {
		return fmt.Errorf("%w: invalid ID", ErrInvalidRecord)
	}
	if event.Actor.Type != "system" && event.Actor.Type != "operator" && event.Actor.Type != "agent" && event.Actor.Type != "policy" {
		return fmt.Errorf("%w: invalid actor type %q", ErrInvalidRecord, event.Actor.Type)
	}
	if len(event.Reason.Message) > 1024 {
		return fmt.Errorf("%w: reason message too long", ErrInvalidRecord)
	}
	if _, err := time.Parse("2006-01-02T15:04:05.000Z", event.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid timestamp", ErrInvalidRecord)
	}
	if event.Sequence < 1 || event.Entity.RunID != event.RunID || !validID(event.Entity.RunID) {
		return fmt.Errorf("%w: invalid sequence or entity run ID", ErrInvalidRecord)
	}
	if err := validateEntity(event.Entity); err != nil {
		return err
	}
	if event.InputEvidence == nil {
		return fmt.Errorf("%w: input evidence must be an array", ErrInvalidRecord)
	}
	seenEvidence := make(map[string]struct{}, len(event.InputEvidence))
	for _, hash := range event.InputEvidence {
		if !validHash(hash) {
			return fmt.Errorf("%w: invalid input evidence hash", ErrInvalidRecord)
		}
		if _, exists := seenEvidence[hash]; exists {
			return fmt.Errorf("%w: duplicate input evidence hash", ErrInvalidRecord)
		}
		seenEvidence[hash] = struct{}{}
	}
	if event.Sequence == 1 {
		if event.PreviousEventHash != nil || event.Entity.Kind != "chain" || event.FromState != "NONE" || event.ToState != "CREATED" {
			return fmt.Errorf("%w: invalid genesis", ErrInvalidRecord)
		}
	} else if event.PreviousEventHash == nil || !validHash(*event.PreviousEventHash) {
		return fmt.Errorf("%w: missing previous event hash", ErrInvalidRecord)
	}
	if !validTransition(event.Entity.Kind, event.FromState, event.ToState) {
		return fmt.Errorf("%w: invalid %s transition %s -> %s", ErrInvalidRecord, event.Entity.Kind, event.FromState, event.ToState)
	}
	return nil
}

func validateEntity(entity Entity) error {
	switch entity.Kind {
	case "chain":
		if entity.TaskID != "" || entity.StepID != "" || entity.Attempt != 0 {
			return fmt.Errorf("%w: invalid chain entity", ErrInvalidRecord)
		}
	case "task":
		if !validID(entity.TaskID) || entity.StepID != "" || entity.Attempt < 1 {
			return fmt.Errorf("%w: invalid task entity", ErrInvalidRecord)
		}
	case "step":
		if !validID(entity.TaskID) || !validID(entity.StepID) || entity.Attempt < 1 {
			return fmt.Errorf("%w: invalid step entity", ErrInvalidRecord)
		}
	default:
		return fmt.Errorf("%w: invalid entity kind %q", ErrInvalidRecord, entity.Kind)
	}
	return nil
}

func validTransition(kind, from, to string) bool {
	transitions := map[string]map[string]map[string]bool{
		"chain": {
			"NONE": {"CREATED": true}, "CREATED": {"BASELINED": true, "FAILED": true, "CANCELLED": true},
			"BASELINED": {"RUNNING": true, "FAILED": true, "CANCELLED": true}, "RUNNING": {"PAUSED": true, "COMPLETED": true, "FAILED": true, "CANCELLED": true},
			"PAUSED": {"RUNNING": true, "FAILED": true, "CANCELLED": true},
		},
		"task": {
			"NONE": {"PENDING": true}, "PENDING": {"PRECHECK": true, "FAILED": true, "CANCELLED": true},
			"PRECHECK": {"STEPS_RUNNING": true, "FAILED": true, "CANCELLED": true}, "STEPS_RUNNING": {"WAITING_FOR_OPERATOR": true, "REVIEW_PENDING": true, "FAILED": true, "CANCELLED": true},
			"WAITING_FOR_OPERATOR": {"STEPS_RUNNING": true, "FAILED": true, "CANCELLED": true}, "REVIEW_PENDING": {"PASSED": true, "REPAIR_PENDING": true, "FAILED": true, "CANCELLED": true},
			"FAILED": {"REPAIR_PENDING": true}, "REPAIR_PENDING": {"STEPS_RUNNING": true, "FAILED": true, "CANCELLED": true},
		},
		"step": {
			"NONE": {"PENDING": true}, "PENDING": {"RUNNING": true, "NOOP_RECORDED": true, "FAILED": true, "CANCELLED": true},
			"RUNNING": {"WAITING_FOR_OPERATOR": true, "PASSED": true, "FAILED": true, "CANCELLED": true}, "WAITING_FOR_OPERATOR": {"RUNNING": true, "FAILED": true, "CANCELLED": true},
		},
	}
	return transitions[kind][from][to]
}

func eventEntityKey(entity Entity) string {
	return fmt.Sprintf("%s\x00%s\x00%s\x00%s\x00%d", entity.Kind, entity.RunID, entity.TaskID, entity.StepID, entity.Attempt)
}

func canonicalValue(value any) ([]byte, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return Canonicalize(encoded)
}

func decodeStrictJSON(input []byte, destination any) error {
	if _, err := Canonicalize(input); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: trailing value", ErrInvalidJSON)
	}
	return nil
}

func validID(value string) bool   { return idPattern.MatchString(value) }
func validHash(value string) bool { return hashPattern.MatchString(value) }
