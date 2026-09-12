package adapters

import (
	"errors"
	"sync"
	"testing"
)

func intPointer(value int) *int { return &value }

func stringPointer(value string) *string { return &value }

func agentRunnerCompletion(completionID string) AgentRunnerCompletion {
	return AgentRunnerCompletion{
		CompletionID:       completionID,
		RequestID:          "request-one",
		RequestHash:        queueHashOne,
		RunID:              "run-one",
		TaskID:             "task-one",
		StepID:             "step-one",
		Attempt:            1,
		AdapterID:          "fixture-agent",
		SessionID:          "session-one",
		CompletedAt:        "2026-09-10T02:05:00.000Z",
		Status:             "completed",
		ExitCode:           intPointer(0),
		ProcessTreeStatus:  "stopped",
		LogsComplete:       true,
		UsageComplete:      true,
		OutputManifestHash: stringPointer(queueHashTwo),
		Evidence:           []string{queueHashThree},
		ErrorEvidence:      []string{},
	}
}

func TestAgentRunnerCompletionRequiresCompleteExternalEvidence(t *testing.T) {
	completion := agentRunnerCompletion("completion-one")
	completion.ProcessTreeStatus = "unknown"
	if _, err := NewAgentRunnerCompletionRecord(completion); !errors.Is(err, ErrInvalidAgentRunnerCompletion) {
		t.Fatalf("expected incomplete completion rejection, got %v", err)
	}
	completion.Status = "uncertain"
	completion.ErrorEvidence = []string{queueHashFour}
	if _, err := NewAgentRunnerCompletionRecord(completion); err != nil {
		t.Fatalf("expected uncertain receipt, got %v", err)
	}
}

func TestAgentRunnerCompletionBinding(t *testing.T) {
	requestBody := agentRunnerRequest("request-one")
	requestBody.Mode = "resume"
	requestBody.PriorSessionID = "session-one"
	requestBody.PriorCompletionHash = queueHashFour
	request, err := NewAgentRunnerRequestRecord(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	completionBody := agentRunnerCompletion("completion-one")
	completionBody.RequestHash = request.RecordHash
	completion, err := NewAgentRunnerCompletionRecord(completionBody)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAgentRunnerCompletionBinding(request, completion); err != nil {
		t.Fatal(err)
	}
}

func TestAgentRunnerCompletionReplayAndRequestConflict(t *testing.T) {
	first, err := NewAgentRunnerCompletionRecord(agentRunnerCompletion("completion-one"))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerCompletionIndex([]AgentRunnerCompletionRecord{first})
	if err != nil {
		t.Fatal(err)
	}
	if replayed, err := index.Record(first); err != nil || !replayed {
		t.Fatalf("expected replay, replayed=%v err=%v", replayed, err)
	}
	conflicting, err := NewAgentRunnerCompletionRecord(agentRunnerCompletion("completion-two"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := index.Record(conflicting); !errors.Is(err, ErrAgentRunnerCompletionConflict) {
		t.Fatalf("expected request terminal conflict, got %v", err)
	}
}

func TestAgentRunnerCompletionConcurrentReplay(t *testing.T) {
	record, err := NewAgentRunnerCompletionRecord(agentRunnerCompletion("completion-one"))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerCompletionIndex(nil)
	if err != nil {
		t.Fatal(err)
	}

	const workers = 24
	start := make(chan struct{})
	results := make(chan struct {
		replayed bool
		err      error
	}, workers)
	var waitGroup sync.WaitGroup
	waitGroup.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer waitGroup.Done()
			<-start
			replayed, err := index.Record(record)
			results <- struct {
				replayed bool
				err      error
			}{replayed: replayed, err: err}
		}()
	}
	close(start)
	waitGroup.Wait()
	close(results)

	insertedCount := 0
	replayedCount := 0
	for result := range results {
		if result.err != nil {
			t.Errorf("unexpected concurrent replay error: %v", result.err)
			continue
		}
		if result.replayed {
			replayedCount++
		} else {
			insertedCount++
		}
	}
	if insertedCount != 1 || replayedCount != workers-1 {
		t.Fatalf("expected one insert and %d replays, got inserts=%d replays=%d", 1, insertedCount, replayedCount)
	}
}

func TestAgentRunnerCompletionConcurrentRequestConflict(t *testing.T) {
	first, err := NewAgentRunnerCompletionRecord(agentRunnerCompletion("completion-one"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewAgentRunnerCompletionRecord(agentRunnerCompletion("completion-two"))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerCompletionIndex(nil)
	if err != nil {
		t.Fatal(err)
	}

	start := make(chan struct{})
	results := make(chan struct {
		record   AgentRunnerCompletionRecord
		replayed bool
		err      error
	}, 2)
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)
	for _, record := range []AgentRunnerCompletionRecord{first, second} {
		go func(record AgentRunnerCompletionRecord) {
			defer waitGroup.Done()
			<-start
			replayed, err := index.Record(record)
			results <- struct {
				record   AgentRunnerCompletionRecord
				replayed bool
				err      error
			}{record: record, replayed: replayed, err: err}
		}(record)
	}
	close(start)
	waitGroup.Wait()
	close(results)

	successCount := 0
	conflictCount := 0
	var failedRecord AgentRunnerCompletionRecord
	failedRecordSet := false
	for result := range results {
		switch {
		case result.err == nil:
			successCount++
			if result.replayed {
				t.Errorf("expected concurrent winner to be inserted, got replay for %s", result.record.Completion.CompletionID)
			}
		case errors.Is(result.err, ErrAgentRunnerCompletionConflict):
			conflictCount++
			failedRecord = result.record
			failedRecordSet = true
		default:
			t.Errorf("unexpected concurrent request conflict error: %v", result.err)
		}
	}
	if successCount != 1 || conflictCount != 1 || !failedRecordSet {
		t.Fatalf("expected one success and one conflict, got successes=%d conflicts=%d", successCount, conflictCount)
	}

	replayed, err := index.Record(failedRecord)
	if replayed {
		t.Fatalf("expected failed completion replay to remain a conflict, got replay=true")
	}
	if !errors.Is(err, ErrAgentRunnerCompletionConflict) {
		t.Fatalf("expected failed completion replay conflict, got %v", err)
	}
}
