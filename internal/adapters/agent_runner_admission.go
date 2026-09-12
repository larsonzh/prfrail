package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrAgentRunnerCompositeAdmissionBlocked = errors.New("AgentRunner composite admission blocked")
	ErrAgentRunnerRequestReplayBlocked      = errors.New("AgentRunner request replay blocked")
)

// AgentRunnerCompositeAdmission composes offline, persisted preflight facts.
// EnforcementRecord is optional when the capability record is compatible;
// AdmissionPolicy is still required by ValidateAgentRunnerAdmission.
// AuthorizationHash and BudgetHash only bind this request to frozen records.
// This slice does not decide whether authorization or cost ledger entries are
// active or unexhausted; that remains a later chain/tickets responsibility.
// Clock is the only evaluation-time authority; policy timestamps are ignored.
type AgentRunnerCompositeAdmission struct {
	RequestRecord      AgentRunnerRequestRecord
	RequestIndex       *AgentRunnerRequestIndex
	AvailabilityRecord AIAvailabilityRecord
	AvailabilityPolicy AIAvailabilityPolicy
	CapabilityRecord   AgentRunnerCapabilityRecord
	EnforcementRecord  *AgentRunnerEnforcementRecord
	AdmissionPolicy    AgentRunnerAdmissionPolicy
	Clock              func() time.Time
}

var _ chain.AgentRunnerAdmission = (*AgentRunnerCompositeAdmission)(nil)

func (admission *AgentRunnerCompositeAdmission) AdmitAgentRunner(ctx context.Context, request chain.StepRequest) error {
	if admission == nil {
		return compositeAdmissionBlocked(nil, "nil admission")
	}
	if err := ctx.Err(); err != nil {
		return compositeAdmissionBlocked(err, "admission context cancelled")
	}
	if err := ValidateAgentRunnerRequestRecord(admission.RequestRecord); err != nil {
		return compositeAdmissionBlocked(err, "invalid request record")
	}
	if admission.RequestIndex == nil {
		return compositeAdmissionBlocked(nil, "request replay index missing")
	}
	if admission.Clock == nil {
		return compositeAdmissionBlocked(nil, "evaluation clock missing")
	}

	facts := request.AgentRunnerFacts
	if facts == nil {
		return compositeAdmissionBlocked(nil, "immutable facts missing")
	}
	if !evidence.ValidID(facts.RequestID) ||
		!evidence.ValidHash(facts.WorkspaceHash) ||
		!evidence.ValidHash(facts.ContextHash) ||
		!evidence.ValidHash(facts.AuthorizationHash) ||
		!evidence.ValidHash(facts.BudgetHash) {
		return compositeAdmissionBlocked(nil, "immutable facts contain an invalid binding")
	}

	body := admission.RequestRecord.Request
	if request.Step.Kind != "code" || request.Step.Mode != chain.IsolatedWorkspace || request.ExecutionTarget != chain.AgentRunnerExecution {
		return compositeAdmissionBlocked(nil, "AgentRunner requires isolated workspace code step")
	}
	if body.Mode != "create" {
		return compositeAdmissionBlocked(nil, "resume continuity evidence is not implemented")
	}
	if request.RunID != body.RunID ||
		request.TaskID != body.TaskID ||
		request.Step.ID != body.StepID ||
		request.Attempt != body.Attempt ||
		request.ParentHash != body.ParentSnapshotHash ||
		body.AdapterID != admission.CapabilityRecord.Capability.AdapterID {
		return compositeAdmissionBlocked(nil, "step request or adapter candidate binding mismatch")
	}
	if facts.RequestID != body.RequestID ||
		facts.WorkspaceHash != body.WorkspaceHash ||
		facts.ContextHash != body.ContextHash ||
		facts.AuthorizationHash != body.AuthorizationHash ||
		facts.BudgetHash != body.BudgetHash {
		return compositeAdmissionBlocked(nil, "immutable facts binding mismatch")
	}

	if admission.AvailabilityPolicy.Channel != "agent-runner-cli" {
		return compositeAdmissionBlocked(nil, "AgentRunner requires agent-runner-cli availability")
	}
	now := admission.Clock().UTC()
	if now.IsZero() {
		return compositeAdmissionBlocked(nil, "evaluation clock returned zero time")
	}
	availabilityPolicy := admission.AvailabilityPolicy
	availabilityPolicy.EvaluatedAt = now
	if err := ValidateAIAvailabilityAdmission(admission.AvailabilityRecord, availabilityPolicy); err != nil {
		return compositeAdmissionBlocked(err, "AI availability admission blocked")
	}
	admissionPolicy := admission.AdmissionPolicy
	admissionPolicy.EvaluatedAt = now
	if err := ValidateAgentRunnerAdmission(admission.CapabilityRecord, admission.EnforcementRecord, admissionPolicy); err != nil {
		return compositeAdmissionBlocked(err, "AgentRunner capability or enforcement admission blocked")
	}
	replayed, err := admission.RequestIndex.Record(admission.RequestRecord)
	if err != nil {
		return compositeAdmissionBlocked(err, "request replay or conflict")
	}
	if replayed {
		return compositeAdmissionBlocked(ErrAgentRunnerRequestReplayBlocked, "durable dispatch or completion state unavailable")
	}
	return nil
}

func compositeAdmissionBlocked(cause error, reason string) error {
	if cause == nil {
		return fmt.Errorf("%w: %s", ErrAgentRunnerCompositeAdmissionBlocked, reason)
	}
	return fmt.Errorf("%w: %s: %w", ErrAgentRunnerCompositeAdmissionBlocked, reason, cause)
}
