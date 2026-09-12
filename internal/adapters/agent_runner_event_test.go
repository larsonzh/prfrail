package adapters

import (
	"errors"
	"sync"
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

func TestAgentRunnerEventIndexConcurrentReplay(t *testing.T) {
	record, err := NewAgentRunnerEventRecord(agentRunnerEvent("event-one", 1))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerEventIndex(nil)
	if err != nil {
		t.Fatal(err)
	}
	const callers = 24
	results := make(chan struct {
		replayed bool
		err      error
	}, callers)
	start := make(chan struct{})
	var waitGroup sync.WaitGroup
	for call := 0; call < callers; call++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			<-start
			replayed, recordErr := index.Record(record)
			results <- struct {
				replayed bool
				err      error
			}{replayed: replayed, err: recordErr}
		}()
	}
	close(start)
	waitGroup.Wait()
	close(results)

	insertions := 0
	replays := 0
	for result := range results {
		if result.err != nil {
			t.Fatalf("unexpected concurrent replay error: %v", result.err)
		}
		if result.replayed {
			replays++
		} else {
			insertions++
		}
	}
	if insertions != 1 {
		t.Fatalf("expected one insertion, got %d", insertions)
	}
	if replays != callers-1 {
		t.Fatalf("expected %d replays, got %d", callers-1, replays)
	}
}

func TestAgentRunnerEventIndexConcurrentSequenceConflict(t *testing.T) {
	first, err := NewAgentRunnerEventRecord(agentRunnerEvent("event-one", 1))
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewAgentRunnerEventRecord(agentRunnerEvent("event-two", 1))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerEventIndex(nil)
	if err != nil {
		t.Fatal(err)
	}
	results := make(chan error, 2)
	start := make(chan struct{})
	var waitGroup sync.WaitGroup
	for _, record := range []AgentRunnerEventRecord{first, second} {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			<-start
			_, recordErr := index.Record(record)
			results <- recordErr
		}()
	}
	close(start)
	waitGroup.Wait()
	close(results)

	successes := 0
	conflicts := 0
	for recordErr := range results {
		switch {
		case recordErr == nil:
			successes++
		case errors.Is(recordErr, ErrAgentRunnerEventConflict):
			conflicts++
		default:
			t.Fatalf("unexpected concurrent conflict error: %v", recordErr)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("expected one successful insert and one conflict, successes=%d conflicts=%d", successes, conflicts)
	}
}
