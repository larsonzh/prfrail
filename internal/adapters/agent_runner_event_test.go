package adapters

import (
	"errors"
	"testing"
)

func agentRunnerEvent(eventID string, sequence int) AgentRunnerEvent {
	return AgentRunnerEvent{
		EventID:     eventID,
		RequestID:   "request-one",
		RequestHash: queueHashOne,
		RunID:       "run-one",
		TaskID:      "task-one",
		StepID:      "step-one",
		Attempt:     1,
		AdapterID:   "fixture-agent",
		SessionID:   "session-one",
		Sequence:    sequence,
		OccurredAt:  "2026-09-10T02:00:01.000Z",
		EventType:   "output",
		PayloadHash: queueHashTwo,
		Evidence:    []string{queueHashThree},
	}
}

func TestAgentRunnerEventBinding(t *testing.T) {
	requestBody := agentRunnerRequest("request-one")
	requestBody.Mode = "resume"
	requestBody.PriorSessionID = "session-one"
	requestBody.PriorCompletionHash = queueHashFour
	request, err := NewAgentRunnerRequestRecord(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	eventBody := agentRunnerEvent("event-one", 1)
	eventBody.RequestHash = request.RecordHash
	event, err := NewAgentRunnerEventRecord(eventBody)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAgentRunnerEventBinding(request, event); err != nil {
		t.Fatal(err)
	}
	event.Event.SessionID = "session-two"
	event.RecordHash, err = digestMessage(agentRunnerEventDomain, event.Event)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAgentRunnerEventBinding(request, event); !errors.Is(err, ErrInvalidAgentRunnerEvent) {
		t.Fatalf("expected session mismatch, got %v", err)
	}
}

func TestAgentRunnerEventReplayAndSequenceConflict(t *testing.T) {
	first, err := NewAgentRunnerEventRecord(agentRunnerEvent("event-one", 1))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerEventIndex([]AgentRunnerEventRecord{first})
	if err != nil {
		t.Fatal(err)
	}
	if replayed, err := index.Record(first); err != nil || !replayed {
		t.Fatalf("expected replay, replayed=%v err=%v", replayed, err)
	}
	conflictingBody := agentRunnerEvent("event-two", 1)
	conflicting, err := NewAgentRunnerEventRecord(conflictingBody)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := index.Record(conflicting); !errors.Is(err, ErrAgentRunnerEventConflict) {
		t.Fatalf("expected sequence conflict, got %v", err)
	}
	driftedBody := agentRunnerEvent("event-three", 2)
	driftedBody.SessionID = "session-two"
	drifted, err := NewAgentRunnerEventRecord(driftedBody)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := index.Record(drifted); !errors.Is(err, ErrAgentRunnerEventConflict) {
		t.Fatalf("expected request session conflict, got %v", err)
	}
}
