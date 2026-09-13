package adapters

import (
	"encoding/json"
	"errors"
	"testing"
)

func agentRunnerRequest(requestID string) AgentRunnerRequest {
	return AgentRunnerRequest{
		RequestID:            requestID,
		RunID:                "run-one",
		TaskID:               "task-one",
		StepID:               "step-one",
		Attempt:              1,
		CreatedAt:            "2026-09-10T02:00:00.000Z",
		AdapterID:            "fixture-agent",
		Mode:                 "create",
		WorkspaceHash:        queueHashOne,
		ContextHash:          queueHashTwo,
		ParentSnapshotHash:   queueHashThree,
		AuthorizationHash:    queueHashFour,
		BudgetHash:           "sha256:5555555555555555555555555555555555555555555555555555555555555555",
		AllowedTargets:       []string{"workspace-files"},
		AllowedEffects:       []string{"workspace-write", "local-process"},
		EffectMappingVersion: agentRunnerEffectMappingVersion,
		EffectMappingHash:    agentRunnerEffectMappingHash,
		Evidence:             []string{"sha256:6666666666666666666666666666666666666666666666666666666666666666"},
	}
}

func TestAgentRunnerRequestRecordCanonicalRoundTrip(t *testing.T) {
	request := agentRunnerRequest("request-one")
	record, err := NewAgentRunnerRequestRecord(request)
	if err != nil {
		t.Fatal(err)
	}
	if record.Request.AllowedEffects[0] != "local-process" {
		t.Fatalf("allowed effects were not normalized: %v", record.Request.AllowedEffects)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeAgentRunnerRequestRecord(wire)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.RecordHash != record.RecordHash {
		t.Fatalf("record hash mismatch: %s != %s", decoded.RecordHash, record.RecordHash)
	}
}

func TestAgentRunnerEffectMappingContract(t *testing.T) {
	body := map[string]string{
		"local-process":   "local-discardable",
		"workspace-write": "local-discardable",
	}
	hash, err := digestMessage("proofrail:agent-runner-effect-mapping:1\n", body)
	if err != nil {
		t.Fatal(err)
	}
	if hash != agentRunnerEffectMappingHash {
		t.Fatalf("effect mapping hash mismatch: %s", hash)
	}
	for operation, want := range body {
		got, found := agentRunnerEffectClass(operation)
		if !found || got != want {
			t.Fatalf("unexpected v1 effect mapping for %q: %q, found=%v", operation, got, found)
		}
	}
}

func TestAgentRunnerRequestRejectsInvalidEffectMappingBindings(t *testing.T) {
	for name, mutate := range map[string]func(*AgentRunnerRequest){
		"missing version": func(request *AgentRunnerRequest) { request.EffectMappingVersion = "" },
		"unknown version": func(request *AgentRunnerRequest) { request.EffectMappingVersion = "2" },
		"missing hash":    func(request *AgentRunnerRequest) { request.EffectMappingHash = "" },
		"tampered hash":   func(request *AgentRunnerRequest) { request.EffectMappingHash = queueHashOne },
		"unknown effect":  func(request *AgentRunnerRequest) { request.AllowedEffects = []string{"external-publish"} },
	} {
		t.Run(name, func(t *testing.T) {
			request := agentRunnerRequest("request-one")
			mutate(&request)
			if _, err := NewAgentRunnerRequestRecord(request); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
				t.Fatalf("expected invalid effect mapping rejection, got %v", err)
			}
		})
	}
}

func TestAgentRunnerRequestRecordRejectsTamperingAndUnknownFields(t *testing.T) {
	record, err := NewAgentRunnerRequestRecord(agentRunnerRequest("request-one"))
	if err != nil {
		t.Fatal(err)
	}
	record.Request.ContextHash = queueHashOne
	if err := ValidateAgentRunnerRequestRecord(record); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected hash mismatch, got %v", err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	wire = append(wire[:len(wire)-1], []byte(`,"unknown":true}`)...)
	if _, err := DecodeAgentRunnerRequestRecord(wire); err == nil {
		t.Fatal("expected unknown field rejection")
	}
}

func TestAgentRunnerRequestModeBindings(t *testing.T) {
	create := agentRunnerRequest("request-one")
	create.PriorSessionID = "session-one"
	if _, err := NewAgentRunnerRequestRecord(create); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected create continuity rejection, got %v", err)
	}
	resume := agentRunnerRequest("request-two")
	resume.Mode = "resume"
	if _, err := NewAgentRunnerRequestRecord(resume); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected incomplete resume rejection, got %v", err)
	}
	resume.PriorSessionID = "session-one"
	resume.PriorCompletionHash = queueHashOne
	if _, err := NewAgentRunnerRequestRecord(resume); err != nil {
		t.Fatalf("expected valid resume, got %v", err)
	}
}

func TestAgentRunnerRequestReplayAndConflict(t *testing.T) {
	first, err := NewAgentRunnerRequestRecord(agentRunnerRequest("request-one"))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerRequestIndex([]AgentRunnerRequestRecord{first})
	if err != nil {
		t.Fatal(err)
	}
	replayed, err := index.Record(first)
	if err != nil || !replayed {
		t.Fatalf("expected idempotent replay, replayed=%v err=%v", replayed, err)
	}
	changed := agentRunnerRequest("request-one")
	changed.ContextHash = queueHashOne
	conflicting, err := NewAgentRunnerRequestRecord(changed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := index.Record(conflicting); !errors.Is(err, ErrAgentRunnerRequestConflict) {
		t.Fatalf("expected replay conflict, got %v", err)
	}
}

func TestAgentRunnerRequestIndexConcurrentReplay(t *testing.T) {
	record, err := NewAgentRunnerRequestRecord(agentRunnerRequest("request-one"))
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerRequestIndex(nil)
	if err != nil {
		t.Fatal(err)
	}
	errorsByCall := make(chan error, 16)
	replaysByCall := make(chan bool, 16)
	for call := 0; call < cap(errorsByCall); call++ {
		go func() {
			replayed, recordErr := index.Record(record)
			replaysByCall <- replayed
			errorsByCall <- recordErr
		}()
	}
	replays := 0
	for call := 0; call < cap(errorsByCall); call++ {
		if recordErr := <-errorsByCall; recordErr != nil {
			t.Fatal(recordErr)
		}
		if <-replaysByCall {
			replays++
		}
	}
	if replays != cap(errorsByCall)-1 {
		t.Fatalf("expected one insertion and %d replays, got %d replays", cap(errorsByCall)-1, replays)
	}
}
