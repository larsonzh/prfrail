package adapters

import (
	"context"
	"errors"
	"testing"

	"github.com/larsonzh/prfrail/internal/chain"
)

// P50 CONSOLE-MULTITASK-ASSEMBLY -- multi-record intent registry.

func agentRunnerRegistryDefinition() chain.Definition {
	return chain.Definition{ID: "chain-registry", Tasks: []chain.Task{
		{ID: "task-one", Steps: []chain.Step{{ID: "step-one", Kind: "code", Mode: chain.IsolatedWorkspace}}},
		{ID: "task-two", Steps: []chain.Step{{ID: "step-two", Kind: "code", Mode: chain.IsolatedWorkspace}}},
	}}
}

func mustAgentRunnerRegistryRecord(t *testing.T, requestID, taskID, stepID string) AgentRunnerRequestRecord {
	t.Helper()
	body := agentRunnerRequest(requestID)
	body.TaskID = taskID
	body.StepID = stepID
	record, err := NewAgentRunnerRequestRecord(body)
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func agentRunnerRegistryRequest(record AgentRunnerRequestRecord) chain.StepRequest {
	body := record.Request
	return chain.StepRequest{
		RunID:      body.RunID,
		TaskID:     body.TaskID,
		Step:       chain.Step{ID: body.StepID, Kind: "code", Mode: chain.IsolatedWorkspace},
		Attempt:    body.Attempt,
		ParentHash: body.ParentSnapshotHash,
		Workspace:  chain.Workspace{Root: "workspace/" + body.TaskID},
	}
}

// TestAgentRunnerStepIntentRegistryBindsEachRecordToItsOwnStep proves that a
// multi-task registry projects every request through the record it is bound to:
// each task receives its own requestId and immutable facts, never another task's.
func TestAgentRunnerStepIntentRegistryBindsEachRecordToItsOwnStep(t *testing.T) {
	records := []AgentRunnerRequestRecord{
		mustAgentRunnerRegistryRecord(t, "request-one", "task-one", "step-one"),
		mustAgentRunnerRegistryRecord(t, "request-two", "task-two", "step-two"),
	}
	definition := agentRunnerRegistryDefinition()
	registry, err := NewAgentRunnerStepIntentRegistry(records, definition)
	if err != nil {
		t.Fatal(err)
	}
	if registry.Len() != 2 {
		t.Fatalf("the registry must bind two task/step entries, got %d", registry.Len())
	}
	for _, record := range records {
		intent, err := registry.PrepareStepIntent(context.Background(), agentRunnerRegistryRequest(record))
		if err != nil {
			t.Fatalf("task %q must project its own record: %v", record.Request.TaskID, err)
		}
		if intent.ExecutionTarget != chain.AgentRunnerExecution || intent.AgentRunnerFacts == nil {
			t.Fatalf("task %q received no AgentRunner intent: %+v", record.Request.TaskID, intent)
		}
		if intent.AgentRunnerFacts.RequestID != record.Request.RequestID {
			t.Fatalf("task %q received request %q, want %q", record.Request.TaskID, intent.AgentRunnerFacts.RequestID, record.Request.RequestID)
		}
		if intent.AgentRunnerFacts.ContextHash != record.Request.ContextHash {
			t.Fatalf("task %q facts do not carry its own context hash", record.Request.TaskID)
		}
		preparer, err := registry.Preparer(record.Request.TaskID, record.Request.StepID)
		if err != nil || preparer == nil {
			t.Fatalf("task %q step %q must be selectable: %v", record.Request.TaskID, record.Request.StepID, err)
		}
	}
	if _, err := registry.Preparer("task-one", "step-two"); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("an unbound task/step pair must be refused, got %v", err)
	}
	var nilRegistry *AgentRunnerStepIntentRegistry
	if _, err := nilRegistry.Preparer("task-one", "step-one"); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("a nil registry must be refused, got %v", err)
	}
	if nilRegistry.Len() != 0 {
		t.Fatal("a nil registry must report zero bindings")
	}
}

// TestAgentRunnerStepIntentRegistryRejectsCrossRecordDrift pins the drift rules
// the multi-task assembly depends on: two records may never bind the same
// task/step, share a requestId, or claim different runs.
func TestAgentRunnerStepIntentRegistryRejectsCrossRecordDrift(t *testing.T) {
	definition := agentRunnerRegistryDefinition()
	base := mustAgentRunnerRegistryRecord(t, "request-one", "task-one", "step-one")

	duplicateBinding := mustAgentRunnerRegistryRecord(t, "request-three", "task-one", "step-one")
	if _, err := NewAgentRunnerStepIntentRegistry([]AgentRunnerRequestRecord{base, duplicateBinding}, definition); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("two records binding one task/step must be refused, got %v", err)
	}

	duplicateRequestID := mustAgentRunnerRegistryRecord(t, "request-one", "task-two", "step-two")
	if _, err := NewAgentRunnerStepIntentRegistry([]AgentRunnerRequestRecord{base, duplicateRequestID}, definition); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("two records sharing one requestId must be refused, got %v", err)
	}

	driftedRun := agentRunnerRequest("request-three")
	driftedRun.TaskID = "task-two"
	driftedRun.StepID = "step-two"
	driftedRun.RunID = "run-two"
	drifted, err := NewAgentRunnerRequestRecord(driftedRun)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewAgentRunnerStepIntentRegistry([]AgentRunnerRequestRecord{base, drifted}, definition); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("records from different runs must be refused, got %v", err)
	}
	if _, err := NewAgentRunnerStepIntentRegistry(nil, definition); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("an empty registry must be refused, got %v", err)
	}
}

// TestNewAgentRunnerStepIntentRegistryReusesPerRecordAdmission proves the
// registry does not re-derive the single-record rules: a record whose
// task/step is absent from the definition, or whose step shape is wrong, is
// refused by the same constructor a single-task assembly uses.
func TestNewAgentRunnerStepIntentRegistryReusesPerRecordAdmission(t *testing.T) {
	definition := agentRunnerRegistryDefinition()

	missing := mustAgentRunnerRegistryRecord(t, "request-one", "task-one", "other-step")
	if _, err := NewAgentRunnerStepIntentRegistry([]AgentRunnerRequestRecord{missing}, definition); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("a record with no definition entry must be refused, got %v", err)
	}

	managed := chain.Definition{ID: definition.ID, Tasks: []chain.Task{
		{ID: "task-one", Steps: []chain.Step{{ID: "step-one", Kind: "code", Mode: chain.ManagedChangeSet}}},
	}}
	record := mustAgentRunnerRegistryRecord(t, "request-one", "task-one", "step-one")
	if _, err := NewAgentRunnerStepIntentRegistry([]AgentRunnerRequestRecord{record}, managed); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("a non isolated-workspace step must be refused, got %v", err)
	}
}
