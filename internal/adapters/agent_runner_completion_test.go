package adapters

import (
	"errors"
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
