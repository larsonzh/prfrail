package adapters

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
)

func compositeAdmissionFixture(t *testing.T) (AgentRunnerCompositeAdmission, chain.StepRequest) {
	t.Helper()
	requestRecord, err := NewAgentRunnerRequestRecord(agentRunnerRequest("request-one"))
	if err != nil {
		t.Fatal(err)
	}
	availabilityRecord, err := NewAIAvailabilityRecord(aiAvailability())
	if err != nil {
		t.Fatal(err)
	}
	capabilityRecord, err := NewAgentRunnerCapabilityRecord(agentRunnerCapability())
	if err != nil {
		t.Fatal(err)
	}
	request := chain.StepRequest{
		RunID:           requestRecord.Request.RunID,
		TaskID:          requestRecord.Request.TaskID,
		Step:            chain.Step{ID: requestRecord.Request.StepID, Kind: "code", Mode: chain.IsolatedWorkspace},
		Attempt:         requestRecord.Request.Attempt,
		ParentHash:      requestRecord.Request.ParentSnapshotHash,
		ExecutionTarget: chain.AgentRunnerExecution,
		AgentRunnerFacts: &chain.AgentRunnerImmutableFacts{
			RequestID:         requestRecord.Request.RequestID,
			WorkspaceHash:     requestRecord.Request.WorkspaceHash,
			ContextHash:       requestRecord.Request.ContextHash,
			AuthorizationHash: requestRecord.Request.AuthorizationHash,
			BudgetHash:        requestRecord.Request.BudgetHash,
		},
	}
	return AgentRunnerCompositeAdmission{
		RequestRecord:      requestRecord,
		RequestIndex:       mustAgentRunnerRequestIndex(t),
		AvailabilityRecord: availabilityRecord,
		AvailabilityPolicy: aiAvailabilityPolicy(),
		CapabilityRecord:   capabilityRecord,
		AdmissionPolicy:    admissionPolicy(),
		Clock: func() time.Time {
			return aiAvailabilityPolicy().EvaluatedAt
		},
	}, request
}

func mustAgentRunnerRequestIndex(t *testing.T) *AgentRunnerRequestIndex {
	t.Helper()
	index, err := NewAgentRunnerRequestIndex(nil)
	if err != nil {
		t.Fatal(err)
	}
	return index
}

func TestAgentRunnerCompositeAdmissionAdmitsOfflineBoundRequest(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	if err := admission.AdmitAgentRunner(context.Background(), request); err != nil {
		t.Fatalf("expected offline composite admission, got %v", err)
	}
}

func TestAgentRunnerCompositeAdmissionRejectsMissingInvalidOrTamperedFacts(t *testing.T) {
	for name, mutate := range map[string]func(*chain.StepRequest){
		"missing": func(request *chain.StepRequest) { request.AgentRunnerFacts = nil },
		"invalid-hash": func(request *chain.StepRequest) {
			request.AgentRunnerFacts.WorkspaceHash = "invalid"
		},
		"invalid-request-id": func(request *chain.StepRequest) {
			request.AgentRunnerFacts.RequestID = "invalid request id"
		},
	} {
		t.Run(name, func(t *testing.T) {
			admission, request := compositeAdmissionFixture(t)
			mutate(&request)
			if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
				t.Fatalf("expected composite admission block, got %v", err)
			}
		})
	}
}

func TestAgentRunnerCompositeAdmissionRejectsMissingExecutionTarget(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	request.ExecutionTarget = chain.DefaultExecution
	if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("expected missing execution target to block, got %v", err)
	}
}

func TestAgentRunnerCompositeAdmissionRejectsMissingReplayIndexOrClock(t *testing.T) {
	for name, mutate := range map[string]func(*AgentRunnerCompositeAdmission){
		"request-index": func(admission *AgentRunnerCompositeAdmission) { admission.RequestIndex = nil },
		"clock":         func(admission *AgentRunnerCompositeAdmission) { admission.Clock = nil },
	} {
		t.Run(name, func(t *testing.T) {
			admission, request := compositeAdmissionFixture(t)
			mutate(&admission)
			if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
				t.Fatalf("expected missing %s to block, got %v", name, err)
			}
		})
	}
}

func TestAgentRunnerCompositeAdmissionRejectsAdapterCandidateMismatch(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	changed := admission.RequestRecord.Request
	changed.AdapterID = "fixture-agent-two"
	record, err := NewAgentRunnerRequestRecord(changed)
	if err != nil {
		t.Fatal(err)
	}
	admission.RequestRecord = record
	if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("expected adapter mismatch to block, got %v", err)
	}
}

func TestAgentRunnerCompositeAdmissionRequiresAgentRunnerCLIAvailability(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	availability := aiAvailability()
	availability.Channel = "sessionbridge-silent"
	record, err := NewAIAvailabilityRecord(availability)
	if err != nil {
		t.Fatal(err)
	}
	admission.AvailabilityRecord = record
	admission.AvailabilityPolicy.Channel = "sessionbridge-silent"
	if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("expected cross-channel availability to block, got %v", err)
	}
}

func TestAgentRunnerCompositeAdmissionUsesCurrentEvaluationTime(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	admission.Clock = func() time.Time {
		return admission.AvailabilityPolicy.EvaluatedAt.Add(11 * time.Minute)
	}
	if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAIUnavailable) {
		t.Fatalf("expected stale availability to block at current time, got %v", err)
	}
}

func TestAgentRunnerCompositeAdmissionRejectsResumeWithoutContinuityEvidence(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	changed := admission.RequestRecord.Request
	changed.Mode = "resume"
	changed.PriorSessionID = "session-one"
	changed.PriorCompletionHash = queueHashOne
	record, err := NewAgentRunnerRequestRecord(changed)
	if err != nil {
		t.Fatal(err)
	}
	admission.RequestRecord = record
	if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("expected unverifiable resume continuity to block, got %v", err)
	}
}

func TestAgentRunnerCompositeAdmissionBindsStepIdentityAndFacts(t *testing.T) {
	for name, mutate := range map[string]func(*chain.StepRequest){
		"run":                func(request *chain.StepRequest) { request.RunID = "run-two" },
		"task":               func(request *chain.StepRequest) { request.TaskID = "task-two" },
		"step":               func(request *chain.StepRequest) { request.Step.ID = "step-two" },
		"attempt":            func(request *chain.StepRequest) { request.Attempt = 2 },
		"parent":             func(request *chain.StepRequest) { request.ParentHash = queueHashOne },
		"request-id":         func(request *chain.StepRequest) { request.AgentRunnerFacts.RequestID = "request-two" },
		"workspace-hash":     func(request *chain.StepRequest) { request.AgentRunnerFacts.WorkspaceHash = queueHashTwo },
		"context-hash":       func(request *chain.StepRequest) { request.AgentRunnerFacts.ContextHash = queueHashThree },
		"authorization-hash": func(request *chain.StepRequest) { request.AgentRunnerFacts.AuthorizationHash = queueHashOne },
		"budget-hash":        func(request *chain.StepRequest) { request.AgentRunnerFacts.BudgetHash = queueHashOne },
	} {
		t.Run(name, func(t *testing.T) {
			admission, request := compositeAdmissionFixture(t)
			facts := *request.AgentRunnerFacts
			request.AgentRunnerFacts = &facts
			mutate(&request)
			if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
				t.Fatalf("expected binding mismatch to block, got %v", err)
			}
		})
	}
}

func TestAgentRunnerCompositeAdmissionRejectsInvalidRequestRecord(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	admission.RequestRecord.RecordHash = "invalid"
	if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) || !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected invalid request record to block, got %v", err)
	}
}

func TestAgentRunnerCompositeAdmissionCallsAvailabilityBeforeCapability(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	availability := aiAvailability()
	reason := "secret_missing"
	availability.Status = "unavailable"
	availability.Reason = &reason
	availability.RequestsUsed = 0
	record, err := NewAIAvailabilityRecord(availability)
	if err != nil {
		t.Fatal(err)
	}
	admission.AvailabilityRecord = record
	if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAIUnavailable) || !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("expected availability block before capability validation, got %v", err)
	}
}

func TestAgentRunnerCompositeAdmissionRejectsCapabilityAndEnforcementBlocks(t *testing.T) {
	t.Run("capability", func(t *testing.T) {
		admission, request := compositeAdmissionFixture(t)
		admission.CapabilityRecord = blockedAgentRunnerCapability(t)
		if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
			t.Fatalf("expected capability block, got %v", err)
		}
	})
	t.Run("enforcement", func(t *testing.T) {
		admission, request := compositeAdmissionFixture(t)
		capability := blockedAgentRunnerCapability(t)
		enforcement := agentRunnerEnforcement(capability)
		enforcement.ConfigHash = queueHashOne
		enforcementRecord, err := NewAgentRunnerEnforcementRecord(enforcement)
		if err != nil {
			t.Fatal(err)
		}
		admission.CapabilityRecord = capability
		admission.EnforcementRecord = &enforcementRecord
		if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerAdmissionBlocked) {
			t.Fatalf("expected enforcement block, got %v", err)
		}
	})
}

func TestAgentRunnerCompositeAdmissionRequestReplayBlocksRedispatchAndConflicts(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	index, err := NewAgentRunnerRequestIndex(nil)
	if err != nil {
		t.Fatal(err)
	}
	admission.RequestIndex = index
	if err := admission.AdmitAgentRunner(context.Background(), request); err != nil {
		t.Fatalf("expected first admission, got %v", err)
	}
	if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerRequestReplayBlocked) || !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("expected same-digest replay to block redispatch, got %v", err)
	}

	changed := agentRunnerRequest("request-one")
	changed.ContextHash = queueHashOne
	conflictingRecord, err := NewAgentRunnerRequestRecord(changed)
	if err != nil {
		t.Fatal(err)
	}
	admission.RequestRecord = conflictingRecord
	facts := *request.AgentRunnerFacts
	facts.ContextHash = changed.ContextHash
	request.AgentRunnerFacts = &facts
	if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerRequestConflict) || !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("expected replay conflict to block, got %v", err)
	}
}
