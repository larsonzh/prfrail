package adapters

import (
	"context"
	"errors"
	"testing"

	"github.com/larsonzh/prfrail/internal/chain"
)

func agentRunnerIntentRequest(record AgentRunnerRequestRecord) chain.StepRequest {
	body := record.Request
	return chain.StepRequest{
		RunID:      body.RunID,
		TaskID:     body.TaskID,
		Step:       chain.Step{ID: body.StepID, Kind: "code", Mode: chain.IsolatedWorkspace},
		Attempt:    body.Attempt,
		ParentHash: body.ParentSnapshotHash,
		Workspace:  chain.Workspace{Root: "workspace/task-one"},
	}
}

func agentRunnerIntentDefinition(record AgentRunnerRequestRecord) chain.Definition {
	return chain.Definition{ID: "chain-one", Tasks: []chain.Task{
		{ID: record.Request.TaskID, Steps: []chain.Step{
			{ID: "managed-step", Kind: "code", Mode: chain.ManagedChangeSet},
			{ID: record.Request.StepID, Kind: "code", Mode: chain.IsolatedWorkspace},
		}},
	}}
}

func mustAgentRunnerIntentPreparer(t *testing.T, record AgentRunnerRequestRecord) *AgentRunnerStepIntentPreparer {
	t.Helper()
	preparer, err := NewAgentRunnerStepIntentPreparer(record, agentRunnerIntentDefinition(record))
	if err != nil {
		t.Fatal(err)
	}
	return preparer
}

func TestAgentRunnerStepIntentPreparerProjectsValidatedRequest(t *testing.T) {
	record, err := NewAgentRunnerRequestRecord(agentRunnerRequest("request-one"))
	if err != nil {
		t.Fatal(err)
	}
	preparer := mustAgentRunnerIntentPreparer(t, record)
	intent, err := preparer.PrepareStepIntent(context.Background(), agentRunnerIntentRequest(record))
	if err != nil {
		t.Fatal(err)
	}
	if intent.ExecutionTarget != chain.AgentRunnerExecution || intent.AgentRunnerFacts == nil {
		t.Fatalf("AgentRunner intent missing: %+v", intent)
	}
	facts := intent.AgentRunnerFacts
	body := record.Request
	if facts.RequestID != body.RequestID || facts.WorkspaceHash != body.WorkspaceHash || facts.ContextHash != body.ContextHash || facts.AuthorizationHash != body.AuthorizationHash || facts.BudgetHash != body.BudgetHash {
		t.Fatalf("immutable facts do not match persisted request: %+v", facts)
	}
}

func TestAgentRunnerStepIntentPreparerLeavesOtherStepsOnDefaultRoute(t *testing.T) {
	record, err := NewAgentRunnerRequestRecord(agentRunnerRequest("request-one"))
	if err != nil {
		t.Fatal(err)
	}
	preparer := mustAgentRunnerIntentPreparer(t, record)
	request := agentRunnerIntentRequest(record)
	request.Step.ID = "other-step"
	intent, err := preparer.PrepareStepIntent(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if intent.ExecutionTarget != chain.DefaultExecution || intent.AgentRunnerFacts != nil {
		t.Fatalf("other step received AgentRunner intent: %+v", intent)
	}
}

func TestAgentRunnerStepIntentPreparerRejectsTargetBindingDrift(t *testing.T) {
	record, err := NewAgentRunnerRequestRecord(agentRunnerRequest("request-one"))
	if err != nil {
		t.Fatal(err)
	}
	preparer := mustAgentRunnerIntentPreparer(t, record)
	tests := map[string]func(*chain.StepRequest){
		"attempt":   func(request *chain.StepRequest) { request.Attempt++ },
		"mode":      func(request *chain.StepRequest) { request.Step.Mode = chain.ManagedChangeSet },
		"parent":    func(request *chain.StepRequest) { request.ParentHash = queueHashFour },
		"workspace": func(request *chain.StepRequest) { request.Workspace.Root = "" },
	}
	for name, change := range tests {
		t.Run(name, func(t *testing.T) {
			request := agentRunnerIntentRequest(record)
			change(&request)
			if _, err := preparer.PrepareStepIntent(context.Background(), request); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
				t.Fatalf("expected binding rejection, got %v", err)
			}
		})
	}
}

func TestAgentRunnerStepIntentPreparerRejectsRunDriftAndInvalidRecord(t *testing.T) {
	record, err := NewAgentRunnerRequestRecord(agentRunnerRequest("request-one"))
	if err != nil {
		t.Fatal(err)
	}
	preparer := mustAgentRunnerIntentPreparer(t, record)
	request := agentRunnerIntentRequest(record)
	request.RunID = "other-run"
	if _, err := preparer.PrepareStepIntent(context.Background(), request); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected run binding rejection, got %v", err)
	}
	preparer.requestRecord.RecordHash = queueHashFour
	if _, err := preparer.PrepareStepIntent(context.Background(), agentRunnerIntentRequest(record)); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected invalid record rejection, got %v", err)
	}
}

func TestNewAgentRunnerStepIntentPreparerRejectsMissingOrInvalidTarget(t *testing.T) {
	record, err := NewAgentRunnerRequestRecord(agentRunnerRequest("request-one"))
	if err != nil {
		t.Fatal(err)
	}
	missing := agentRunnerIntentDefinition(record)
	missing.Tasks[0].Steps[1].ID = "other-step"
	if _, err := NewAgentRunnerStepIntentPreparer(record, missing); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected missing target rejection, got %v", err)
	}
	invalid := agentRunnerIntentDefinition(record)
	invalid.Tasks[0].Steps[1].Mode = chain.ManagedChangeSet
	if _, err := NewAgentRunnerStepIntentPreparer(record, invalid); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected invalid target rejection, got %v", err)
	}
}
