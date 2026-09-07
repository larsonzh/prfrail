package chain

import (
	"context"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type Projection struct {
	RunID             string
	ChainState        string
	TaskStates        map[string]string
	StepStates        map[string]string
	Sequence          int
	LastEventHash     string
	AcceptedParent    string
	RecoveryUncertain bool
}

func Rebuild(events []evidence.StateEvent) (Projection, error) {
	projection := Projection{TaskStates: map[string]string{}, StepStates: map[string]string{}}
	if len(events) == 0 {
		return projection, nil
	}
	if err := evidence.VerifyEventChain(events); err != nil {
		return Projection{}, err
	}
	projection.RunID = events[0].Event.RunID
	for _, record := range events {
		event := record.Event
		switch event.Entity.Kind {
		case "chain":
			projection.ChainState = event.ToState
			projection.RecoveryUncertain = event.ToState == "PAUSED" && event.Reason.Code == "recovery-uncertain"
			if event.ToState == "BASELINED" && len(event.InputEvidence) > 0 {
				projection.AcceptedParent = event.InputEvidence[0]
			}
		case "task":
			projection.TaskStates[taskKey(event.Entity.TaskID, event.Entity.Attempt)] = event.ToState
			if event.ToState == "PASSED" && len(event.InputEvidence) > 0 {
				projection.AcceptedParent = event.InputEvidence[0]
			}
		case "step":
			projection.StepStates[stepKey(event.Entity.TaskID, event.Entity.StepID, event.Entity.Attempt)] = event.ToState
		}
		projection.Sequence = event.Sequence
		projection.LastEventHash = record.EventHash
	}
	return projection, nil
}

type stateWriter struct {
	store      EventStore
	projection Projection
	clock      Clock
	ids        IDSource
}

func (writer *stateWriter) append(ctx context.Context, actor evidence.Actor, entity evidence.Entity, to string, input []string, reason string) error {
	from := writer.current(entity)
	sequence := writer.projection.Sequence + 1
	inputEvidence := append([]string{}, input...)
	var previous *string
	if writer.projection.LastEventHash != "" {
		value := writer.projection.LastEventHash
		previous = &value
	}
	record, err := evidence.NewStateEvent(evidence.Event{
		EventID: writer.ids(), RunID: entity.RunID, Sequence: sequence,
		OccurredAt: writer.clock().UTC().Format("2006-01-02T15:04:05.000Z"),
		Actor:      actor, Entity: entity, FromState: from, ToState: to,
		PreviousEventHash: previous, InputEvidence: inputEvidence,
		Reason: evidence.Reason{Code: reason},
	})
	if err != nil {
		return fmt.Errorf("create state event: %w", err)
	}
	if err := writer.store.Append(ctx, record); err != nil {
		return err
	}
	writer.apply(record)
	return nil
}

func (writer *stateWriter) current(entity evidence.Entity) string {
	switch entity.Kind {
	case "chain":
		if writer.projection.ChainState != "" {
			return writer.projection.ChainState
		}
	case "task":
		if state := writer.projection.TaskStates[taskKey(entity.TaskID, entity.Attempt)]; state != "" {
			return state
		}
	case "step":
		if state := writer.projection.StepStates[stepKey(entity.TaskID, entity.StepID, entity.Attempt)]; state != "" {
			return state
		}
	}
	return "NONE"
}

func (writer *stateWriter) apply(record evidence.StateEvent) {
	event := record.Event
	writer.projection.RunID = event.RunID
	writer.projection.Sequence = event.Sequence
	writer.projection.LastEventHash = record.EventHash
	switch event.Entity.Kind {
	case "chain":
		writer.projection.ChainState = event.ToState
		writer.projection.RecoveryUncertain = event.ToState == "PAUSED" && event.Reason.Code == "recovery-uncertain"
		if event.ToState == "BASELINED" && len(event.InputEvidence) > 0 {
			writer.projection.AcceptedParent = event.InputEvidence[0]
		}
	case "task":
		writer.projection.TaskStates[taskKey(event.Entity.TaskID, event.Entity.Attempt)] = event.ToState
		if event.ToState == "PASSED" && len(event.InputEvidence) > 0 {
			writer.projection.AcceptedParent = event.InputEvidence[0]
		}
	case "step":
		writer.projection.StepStates[stepKey(event.Entity.TaskID, event.Entity.StepID, event.Entity.Attempt)] = event.ToState
	}
}

func taskKey(taskID string, attempt int) string { return fmt.Sprintf("%s\x00%d", taskID, attempt) }
func stepKey(taskID, stepID string, attempt int) string {
	return fmt.Sprintf("%s\x00%s\x00%d", taskID, stepID, attempt)
}

func defaultClock() time.Time { return time.Now() }
