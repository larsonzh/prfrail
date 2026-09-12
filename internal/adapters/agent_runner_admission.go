package adapters

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/tickets"
)

var (
	ErrAgentRunnerCompositeAdmissionBlocked = errors.New("AgentRunner composite admission blocked")
	ErrAgentRunnerRequestReplayBlocked      = errors.New("AgentRunner request replay blocked")
)

// AgentRunnerCompositeAdmission composes offline, persisted preflight facts.
// EnforcementRecord is optional when the capability record is compatible;
// AdmissionPolicy is still required by ValidateAgentRunnerAdmission.
// Clock is the only evaluation-time authority; policy timestamps are ignored.
type AgentRunnerCompositeAdmission struct {
	RequestRecord       AgentRunnerRequestRecord
	RequestIndex        *AgentRunnerRequestIndex
	AuthorizationLedger chain.AuthorizationLedger
	CostLedger          *tickets.CostLedger
	AvailabilityRecord  AIAvailabilityRecord
	AvailabilityPolicy  AIAvailabilityPolicy
	CapabilityRecord    AgentRunnerCapabilityRecord
	EnforcementRecord   *AgentRunnerEnforcementRecord
	AdmissionPolicy     AgentRunnerAdmissionPolicy
	Clock               func() time.Time
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
	if err := admission.validateAuthorization(body, now); err != nil {
		return compositeAdmissionBlocked(err, "AgentRunner authorization blocked")
	}
	if err := admission.validateBudget(body); err != nil {
		return compositeAdmissionBlocked(err, "AgentRunner budget blocked")
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

func (admission *AgentRunnerCompositeAdmission) validateAuthorization(request AgentRunnerRequest, now time.Time) error {
	var matched chain.AuthorizationRecord
	matches := 0
	for _, record := range admission.AuthorizationLedger.Grants {
		if record.RecordHash != request.AuthorizationHash {
			continue
		}
		if err := chain.ValidateAuthorizationRecord(record); err != nil {
			return err
		}
		matched = record
		matches++
	}
	if matches != 1 || matched.Grant == nil {
		return fmt.Errorf("authorizationHash matched %d grants", matches)
	}
	for _, record := range admission.AuthorizationLedger.GrantRecordsFor(matched.Grant.AuthorizationID) {
		if err := chain.ValidateAuthorizationRecord(record); err != nil {
			return err
		}
	}
	for _, record := range admission.AuthorizationLedger.Revocations {
		if record.Revocation != nil && record.Revocation.AuthorizationID == matched.Grant.AuthorizationID {
			if err := chain.ValidateAuthorizationRecord(record); err != nil {
				return err
			}
		}
	}
	issuedAt, err := time.Parse(time.RFC3339, matched.Grant.IssuedAt)
	if err != nil {
		return fmt.Errorf("invalid authorization issuedAt: %w", err)
	}
	expiresAt, err := time.Parse(time.RFC3339, matched.Grant.ExpiresAt)
	if err != nil {
		return fmt.Errorf("invalid authorization expiresAt: %w", err)
	}
	if now.Before(issuedAt) {
		return errors.New("authorization grant is pending")
	}
	if !now.Before(expiresAt) {
		return errors.New("authorization grant is expired")
	}
	if _, revocation, revoked := admission.AuthorizationLedger.LatestRevocation(matched.Grant.AuthorizationID); revoked {
		matchedRevocation := false
		for _, record := range admission.AuthorizationLedger.GrantRecordsFor(matched.Grant.AuthorizationID) {
			if record.RecordHash == revocation.AuthorizationHash {
				matchedRevocation = true
				break
			}
		}
		if !matchedRevocation {
			return errors.New("authorization revocation hash mismatch")
		}
		return errors.New("authorization grant is revoked")
	}
	if matched.Grant.RunID != request.RunID ||
		!slices.Contains(matched.Grant.Scope.TaskIDs, request.TaskID) ||
		!slices.Contains(matched.Grant.Scope.StepIDs, request.StepID) {
		return errors.New("authorization scope mismatch")
	}
	for _, target := range request.AllowedTargets {
		if !slices.Contains(matched.Grant.Scope.TargetIDs, target) {
			return errors.New("authorization target scope mismatch")
		}
	}
	return nil
}

func (admission *AgentRunnerCompositeAdmission) validateBudget(request AgentRunnerRequest) error {
	reservation, err := admission.CostLedger.RequireOutstandingReservation(request.BudgetHash)
	if err != nil {
		return err
	}
	if reservation.RunID != request.RunID ||
		reservation.RequestID != request.RequestID ||
		reservation.AuthorizationHash != request.AuthorizationHash {
		return errors.New("cost reservation binding mismatch")
	}
	return nil
}

func compositeAdmissionBlocked(cause error, reason string) error {
	if cause == nil {
		return fmt.Errorf("%w: %s", ErrAgentRunnerCompositeAdmissionBlocked, reason)
	}
	return fmt.Errorf("%w: %s: %w", ErrAgentRunnerCompositeAdmissionBlocked, reason, cause)
}
