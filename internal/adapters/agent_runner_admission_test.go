package adapters

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/tickets"
)

func compositeAdmissionFixture(t *testing.T) (AgentRunnerCompositeAdmission, chain.StepRequest) {
	t.Helper()
	grant, err := chain.NewGrantAuthorizationRecord(chain.AuthorizationGrant{
		RecordID:        "grant-one",
		Kind:            "grant",
		IssuedAt:        "2026-09-11T07:00:00.000Z",
		IssuedBy:        evidence.Actor{Type: "operator", ID: "operator-one"},
		AuthorizationID: "authorization-one",
		RunID:           "run-one",
		RunManifestHash: queueHashOne,
		PolicyHash:      queueHashTwo,
		Scope: chain.AuthorizationScope{
			TaskIDs:       []string{"task-one"},
			StepIDs:       []string{"step-one"},
			TargetIDs:     []string{"workspace-files"},
			EffectClasses: []string{"local-discardable"},
			Budget: chain.AuthorizationBudgetLimit{
				ModelCalls:  1,
				Tokens:      1000,
				WallClockMs: 60000,
				Attempts:    1,
				Currency:    "USD",
			},
		},
		ExpiresAt: "2026-09-11T09:00:00.000Z",
		Evidence:  []string{queueHashThree},
	}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	amount := int64(1000)
	costLedger, err := tickets.NewCostLedger(tickets.CostLedgerConfig{
		LedgerID:       "cost-ledger-one",
		CreatedAt:      time.Date(2026, 9, 11, 7, 0, 0, 0, time.UTC),
		ScopeKind:      tickets.CostScopeRun,
		ScopeHash:      queueHashOne,
		Currency:       "USD",
		PricingMode:    tickets.CostPricingDocumented,
		PricingVersion: "2026-09-11",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := costLedger.AppendAllocation(tickets.AllocationInput{
		EntryID:           "allocation-one",
		OccurredAt:        time.Date(2026, 9, 11, 7, 1, 0, 0, time.UTC),
		AuthorizationHash: grant.RecordHash,
		Limits: tickets.CostLimits{
			AmountMicros: &amount,
			ModelCalls:   1,
			Tokens:       1000,
			WallClockMs:  60000,
			Attempts:     1,
		},
	}); err != nil {
		t.Fatal(err)
	}
	requestBody := agentRunnerRequest("request-one")
	requestBody.AuthorizationHash = grant.RecordHash
	reservation, err := costLedger.Reserve(tickets.ReservationInput{
		EntryID:              "reservation-one",
		OccurredAt:           time.Date(2026, 9, 11, 7, 2, 0, 0, time.UTC),
		RunID:                requestBody.RunID,
		RequestID:            requestBody.RequestID,
		IdempotencyKey:       "reserve-request-one",
		AuthorizationHash:    grant.RecordHash,
		ReservedAmountMicros: &amount,
		ReservedCalls:        1,
		ReservedTokens:       1000,
		PricingEvidence:      []string{queueHashFour},
	})
	if err != nil {
		t.Fatal(err)
	}
	requestBody.BudgetHash = reservation.ReservationHash
	requestRecord, err := NewAgentRunnerRequestRecord(requestBody)
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
		RequestRecord:       requestRecord,
		RequestIndex:        mustAgentRunnerRequestIndex(t),
		AuthorizationLedger: chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{grant}},
		CostLedger:          costLedger,
		AvailabilityRecord:  availabilityRecord,
		AvailabilityPolicy:  aiAvailabilityPolicy(),
		CapabilityRecord:    capabilityRecord,
		AdmissionPolicy:     admissionPolicy(),
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

func admissionCostLedger(t *testing.T, runID, requestID, authorizationHash string) (*tickets.CostLedger, string) {
	t.Helper()
	amount := int64(1000)
	ledger, err := tickets.NewCostLedger(tickets.CostLedgerConfig{
		LedgerID:       "cost-ledger-test",
		CreatedAt:      time.Date(2026, 9, 11, 7, 0, 0, 0, time.UTC),
		ScopeKind:      tickets.CostScopeRun,
		ScopeHash:      queueHashOne,
		Currency:       "USD",
		PricingMode:    tickets.CostPricingDocumented,
		PricingVersion: "2026-09-11",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ledger.AppendAllocation(tickets.AllocationInput{
		EntryID:           "allocation-test",
		OccurredAt:        time.Date(2026, 9, 11, 7, 1, 0, 0, time.UTC),
		AuthorizationHash: authorizationHash,
		Limits: tickets.CostLimits{
			AmountMicros: &amount,
			ModelCalls:   1,
			Tokens:       1000,
			WallClockMs:  60000,
			Attempts:     1,
		},
	}); err != nil {
		t.Fatal(err)
	}
	reservation, err := ledger.Reserve(tickets.ReservationInput{
		EntryID:              "reservation-test",
		OccurredAt:           time.Date(2026, 9, 11, 7, 2, 0, 0, time.UTC),
		RunID:                runID,
		RequestID:            requestID,
		IdempotencyKey:       "reserve-test",
		AuthorizationHash:    authorizationHash,
		ReservedAmountMicros: &amount,
		ReservedCalls:        1,
		ReservedTokens:       1000,
		PricingEvidence:      []string{queueHashFour},
	})
	if err != nil {
		t.Fatal(err)
	}
	return ledger, reservation.ReservationHash
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

func TestAgentRunnerCompositeAdmissionRequiresActiveScopedAuthorizationBeforeReplay(t *testing.T) {
	t.Run("missing grant", func(t *testing.T) {
		admission, request := compositeAdmissionFixture(t)
		validLedger := admission.AuthorizationLedger
		admission.AuthorizationLedger = chain.AuthorizationLedger{}
		if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
			t.Fatalf("expected missing grant to block, got %v", err)
		}
		admission.AuthorizationLedger = validLedger
		if err := admission.AdmitAgentRunner(context.Background(), request); err != nil {
			t.Fatalf("rejected authorization consumed request identity: %v", err)
		}
	})

	t.Run("revoked grant", func(t *testing.T) {
		admission, request := compositeAdmissionFixture(t)
		grant := admission.AuthorizationLedger.Grants[0]
		revocation, err := chain.NewRevocationAuthorizationRecord(chain.AuthorizationRevocation{
			RecordID:          "revocation-one",
			Kind:              "revocation",
			IssuedAt:          "2026-09-11T08:00:00.000Z",
			IssuedBy:          evidence.Actor{Type: "operator", ID: "operator-one"},
			AuthorizationID:   grant.Grant.AuthorizationID,
			AuthorizationHash: grant.RecordHash,
			ReasonCode:        "operator-request",
			StopDisposition:   "not-required",
		}, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		validLedger := admission.AuthorizationLedger
		admission.AuthorizationLedger.Revocations = []chain.AuthorizationRecord{revocation}
		if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
			t.Fatalf("expected revoked grant to block, got %v", err)
		}
		admission.AuthorizationLedger = validLedger
		if err := admission.AdmitAgentRunner(context.Background(), request); err != nil {
			t.Fatalf("rejected revocation consumed request identity: %v", err)
		}
	})

	t.Run("target outside scope", func(t *testing.T) {
		admission, request := compositeAdmissionFixture(t)
		original := admission.RequestRecord
		changed := original.Request
		changed.AllowedTargets = []string{"source-tree"}
		var err error
		admission.RequestRecord, err = NewAgentRunnerRequestRecord(changed)
		if err != nil {
			t.Fatal(err)
		}
		if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
			t.Fatalf("expected target scope mismatch to block, got %v", err)
		}
		admission.RequestRecord = original
		if err := admission.AdmitAgentRunner(context.Background(), request); err != nil {
			t.Fatalf("rejected scope mismatch consumed request identity: %v", err)
		}
	})

	t.Run("bound expired grant with newer active grant", func(t *testing.T) {
		admission, request := compositeAdmissionFixture(t)
		base := *admission.AuthorizationLedger.Grants[0].Grant
		expired := base
		expired.RecordID = "grant-expired"
		expired.IssuedAt = "2026-09-11T06:00:00.000Z"
		expired.ExpiresAt = "2026-09-11T08:00:00.000Z"
		expiredRecord, err := chain.NewGrantAuthorizationRecord(expired, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		active := base
		active.RecordID = "grant-active"
		active.IssuedAt = "2026-09-11T07:30:00.000Z"
		active.ExpiresAt = "2026-09-11T09:00:00.000Z"
		activeRecord, err := chain.NewGrantAuthorizationRecord(active, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		admission.AuthorizationLedger.Grants = []chain.AuthorizationRecord{expiredRecord, activeRecord}
		admission.CostLedger, request.AgentRunnerFacts.BudgetHash = admissionCostLedger(t, request.RunID, request.AgentRunnerFacts.RequestID, expiredRecord.RecordHash)
		request.AgentRunnerFacts.AuthorizationHash = expiredRecord.RecordHash
		changed := admission.RequestRecord.Request
		changed.AuthorizationHash = expiredRecord.RecordHash
		changed.BudgetHash = request.AgentRunnerFacts.BudgetHash
		admission.RequestRecord, err = NewAgentRunnerRequestRecord(changed)
		if err != nil {
			t.Fatal(err)
		}
		if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) || !strings.Contains(err.Error(), "expired") {
			t.Fatalf("expected the hash-bound expired grant to block, got %v", err)
		}
	})
}

func TestAgentRunnerCompositeAdmissionRequiresOutstandingBoundBudgetBeforeReplay(t *testing.T) {
	t.Run("missing ledger", func(t *testing.T) {
		admission, request := compositeAdmissionFixture(t)
		validLedger := admission.CostLedger
		admission.CostLedger = nil
		if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
			t.Fatalf("expected missing cost ledger to block, got %v", err)
		}
		admission.CostLedger = validLedger
		if err := admission.AdmitAgentRunner(context.Background(), request); err != nil {
			t.Fatalf("rejected budget check consumed request identity: %v", err)
		}
	})

	t.Run("settled reservation", func(t *testing.T) {
		admission, request := compositeAdmissionFixture(t)
		amount := int64(1000)
		calls := 1
		tokens := 1000
		if _, err := admission.CostLedger.Settle(tickets.SettlementInput{
			EntryID:             "settlement-one",
			OccurredAt:          time.Date(2026, 9, 11, 8, 0, 0, 0, time.UTC),
			RunID:               request.RunID,
			RequestID:           request.AgentRunnerFacts.RequestID,
			IdempotencyKey:      "settle-request-one",
			ReservationHash:     request.AgentRunnerFacts.BudgetHash,
			Status:              tickets.CostSettlementObserved,
			ChargedAmountMicros: &amount,
			ObservedCalls:       &calls,
			ObservedTokens:      &tokens,
			ProviderEvidence:    []string{queueHashFour},
		}); err != nil {
			t.Fatal(err)
		}
		if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, tickets.ErrCostReservationSettled) || !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) {
			t.Fatalf("expected settled reservation to block, got %v", err)
		}

		healthy, _ := compositeAdmissionFixture(t)
		admission.CostLedger = healthy.CostLedger
		if err := admission.AdmitAgentRunner(context.Background(), request); err != nil {
			t.Fatalf("settled reservation rejection consumed request identity: %v", err)
		}
	})

	for name, binding := range map[string]struct {
		runID             string
		requestID         string
		authorizationHash string
	}{
		"run":           {runID: "run-two", requestID: "request-one"},
		"request":       {runID: "run-one", requestID: "request-two"},
		"authorization": {runID: "run-one", requestID: "request-one", authorizationHash: queueHashOne},
	} {
		t.Run(name+" binding mismatch", func(t *testing.T) {
			admission, request := compositeAdmissionFixture(t)
			authorizationHash := binding.authorizationHash
			if authorizationHash == "" {
				authorizationHash = request.AgentRunnerFacts.AuthorizationHash
			}
			ledger, reservationHash := admissionCostLedger(t, binding.runID, binding.requestID, authorizationHash)
			admission.CostLedger = ledger
			request.AgentRunnerFacts.BudgetHash = reservationHash
			changed := admission.RequestRecord.Request
			changed.BudgetHash = reservationHash
			var err error
			admission.RequestRecord, err = NewAgentRunnerRequestRecord(changed)
			if err != nil {
				t.Fatal(err)
			}
			if err := admission.AdmitAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerCompositeAdmissionBlocked) || !strings.Contains(err.Error(), "binding mismatch") {
				t.Fatalf("expected %s mismatch to block, got %v", name, err)
			}
		})
	}
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

	changed := admission.RequestRecord.Request
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
