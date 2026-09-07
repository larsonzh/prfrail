package chain

import (
	"fmt"
	"slices"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const handoffReceiptDomain = "proofrail:handoff-receipt:1\n"

type HandoffPolicy struct {
	AllowedTargets   []string
	InputPolicy      string
	HandoffTimeoutMs int
	ReturnActions    []string
	HooksAfterReturn []string
}

type HandoffSession struct {
	RunID                string
	TaskID               string
	StepID               string
	Attempt              int
	HandoffPolicyHash    string
	AllowedReturnActions []string
	Operator             evidence.Actor
	OpenedAt             string
	BeforeManifestHash   string
	LeaseEvidence        []string
}

type HandoffDecision struct {
	ReceiptID          string
	ClosedAt           string
	AfterManifestHash  string
	DiffHash           string
	Outcome            string
	HookResultEvidence []string
	Evidence           []string
}

type handoffReceipt struct {
	ReceiptID          string         `json:"receiptId"`
	RecordedBy         evidence.Actor `json:"recordedBy"`
	Operator           evidence.Actor `json:"operator"`
	RunID              string         `json:"runId"`
	TaskID             string         `json:"taskId"`
	StepID             string         `json:"stepId"`
	Attempt            int            `json:"attempt"`
	HandoffPolicyHash  string         `json:"handoffPolicyHash"`
	OpenedAt           string         `json:"openedAt"`
	ClosedAt           string         `json:"closedAt"`
	BeforeManifestHash string         `json:"beforeManifestHash"`
	AfterManifestHash  string         `json:"afterManifestHash"`
	DiffHash           string         `json:"diffHash"`
	LeaseEvidence      []string       `json:"leaseEvidence"`
	Outcome            string         `json:"outcome"`
	HookResultEvidence []string       `json:"hookResultEvidence"`
	Evidence           []string       `json:"evidence"`
}

type HandoffReceiptRecord struct {
	SchemaVersion string         `json:"schemaVersion"`
	Receipt       handoffReceipt `json:"receipt"`
	ReceiptHash   string         `json:"receiptHash"`
}

func BuildHandoffReceiptRecord(session HandoffSession, decision HandoffDecision, clock Clock, ids IDSource) (HandoffReceiptRecord, error) {
	if clock == nil || ids == nil {
		return HandoffReceiptRecord{}, fmt.Errorf("%w: missing clock or id source", ErrInvalidState)
	}
	if decision.ReceiptID == "" {
		decision.ReceiptID = ids()
	}
	if decision.ClosedAt == "" {
		decision.ClosedAt = clock().UTC().Format("2006-01-02T15:04:05.000Z")
	}
	if err := validateHandoffSession(session); err != nil {
		return HandoffReceiptRecord{}, err
	}
	if err := validateHandoffDecision(session, decision); err != nil {
		return HandoffReceiptRecord{}, err
	}
	receipt := handoffReceipt{
		ReceiptID:          decision.ReceiptID,
		RecordedBy:         evidence.Actor{Type: "system", ID: "proofrail"},
		Operator:           session.Operator,
		RunID:              session.RunID,
		TaskID:             session.TaskID,
		StepID:             session.StepID,
		Attempt:            session.Attempt,
		HandoffPolicyHash:  session.HandoffPolicyHash,
		OpenedAt:           session.OpenedAt,
		ClosedAt:           decision.ClosedAt,
		BeforeManifestHash: session.BeforeManifestHash,
		AfterManifestHash:  decision.AfterManifestHash,
		DiffHash:           decision.DiffHash,
		LeaseEvidence:      cloneStrings(session.LeaseEvidence),
		Outcome:            decision.Outcome,
		HookResultEvidence: cloneStrings(decision.HookResultEvidence),
		Evidence:           cloneStrings(decision.Evidence),
	}
	hash, err := digestRecord(handoffReceiptDomain, receipt)
	if err != nil {
		return HandoffReceiptRecord{}, err
	}
	return HandoffReceiptRecord{
		SchemaVersion: evidence.SchemaVersion,
		Receipt:       receipt,
		ReceiptHash:   hash,
	}, nil
}

func OpenHandoffTransition(taskState, stepState string) (nextTaskState, nextStepState string, err error) {
	if taskState != "STEPS_RUNNING" || stepState != "RUNNING" {
		return "", "", fmt.Errorf("%w: handoff open requires STEPS_RUNNING/RUNNING", ErrInvalidState)
	}
	return "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR", nil
}

func ReturnHandoffTransition(taskState, stepState, outcome string) (nextTaskState, nextStepState string, err error) {
	if taskState != "WAITING_FOR_OPERATOR" || stepState != "WAITING_FOR_OPERATOR" {
		return "", "", fmt.Errorf("%w: handoff return requires WAITING_FOR_OPERATOR", ErrInvalidState)
	}
	switch outcome {
	case "complete":
		return "STEPS_RUNNING", "RUNNING", nil
	case "abort", "request-agent":
		// abort keeps failure evidence and discards the candidate; request-agent
		// also enters FAILED so the repair flow can open a new attempt and inject
		// confirmed conclusions. CANCELLED would be terminal and poison the run.
		return "FAILED", "FAILED", nil
	default:
		return "", "", fmt.Errorf("%w: unknown handoff outcome", ErrInvalidState)
	}
}

func validateHandoffPolicy(policy HandoffPolicy) error {
	if !validUniqueIDs(policy.AllowedTargets) {
		return fmt.Errorf("%w: invalid allowed targets", ErrInvalidDefinition)
	}
	if policy.InputPolicy != "structured" && policy.InputPolicy != "secret-direct" {
		return fmt.Errorf("%w: invalid handoff input policy", ErrInvalidDefinition)
	}
	if policy.HandoffTimeoutMs < 1 {
		return fmt.Errorf("%w: invalid handoff timeout", ErrInvalidDefinition)
	}
	if !validUniqueReturnActions(policy.ReturnActions) {
		return fmt.Errorf("%w: invalid handoff return actions", ErrInvalidDefinition)
	}
	if !validUniqueIDs(policy.HooksAfterReturn) {
		return fmt.Errorf("%w: invalid handoff hooks", ErrInvalidDefinition)
	}
	return nil
}

func validateHandoffSession(session HandoffSession) error {
	if !evidence.ValidID(session.RunID) || !evidence.ValidID(session.TaskID) || !evidence.ValidID(session.StepID) || session.Attempt < 1 {
		return fmt.Errorf("%w: invalid handoff identity", ErrInvalidState)
	}
	if session.Operator.Type != "operator" || !evidence.ValidID(session.Operator.ID) {
		return fmt.Errorf("%w: invalid handoff operator", ErrInvalidState)
	}
	if !evidence.ValidHash(session.HandoffPolicyHash) || !evidence.ValidHash(session.BeforeManifestHash) {
		return fmt.Errorf("%w: invalid handoff session hashes", ErrInvalidState)
	}
	if !validUniqueHashSet(session.LeaseEvidence, true) {
		return fmt.Errorf("%w: invalid handoff lease evidence", ErrInvalidState)
	}
	if _, err := time.Parse("2006-01-02T15:04:05.000Z", session.OpenedAt); err != nil {
		return fmt.Errorf("%w: invalid handoff open time", ErrInvalidState)
	}
	if !validUniqueReturnActions(session.AllowedReturnActions) {
		return fmt.Errorf("%w: invalid allowed return action set", ErrInvalidState)
	}
	return nil
}

func validateHandoffDecision(session HandoffSession, decision HandoffDecision) error {
	if !evidence.ValidID(decision.ReceiptID) {
		return fmt.Errorf("%w: invalid handoff receipt id", ErrInvalidState)
	}
	openedAt, err := time.Parse("2006-01-02T15:04:05.000Z", session.OpenedAt)
	if err != nil {
		return fmt.Errorf("%w: invalid handoff open time", ErrInvalidState)
	}
	closedAt, err := time.Parse("2006-01-02T15:04:05.000Z", decision.ClosedAt)
	if err != nil || !closedAt.After(openedAt) {
		return fmt.Errorf("%w: invalid handoff close time", ErrInvalidState)
	}
	if !evidence.ValidHash(decision.AfterManifestHash) || !evidence.ValidHash(decision.DiffHash) {
		return fmt.Errorf("%w: invalid handoff decision hashes", ErrInvalidState)
	}
	if !validUniqueHashSet(decision.HookResultEvidence, false) || !validUniqueHashSet(decision.Evidence, true) {
		return fmt.Errorf("%w: invalid handoff decision evidence", ErrInvalidState)
	}
	if len(session.AllowedReturnActions) > 0 && !slices.Contains(session.AllowedReturnActions, decision.Outcome) {
		return fmt.Errorf("%w: handoff outcome not allowed by policy", ErrInvalidState)
	}
	switch decision.Outcome {
	case "complete":
		if len(decision.HookResultEvidence) == 0 {
			return fmt.Errorf("%w: complete handoff requires hook result evidence", ErrInvalidState)
		}
	case "abort", "request-agent":
		if len(decision.HookResultEvidence) != 0 {
			return fmt.Errorf("%w: abort/request-agent must not include hook result evidence", ErrInvalidState)
		}
	default:
		return fmt.Errorf("%w: unknown handoff outcome", ErrInvalidState)
	}
	return nil
}

func validUniqueIDs(values []string) bool {
	if len(values) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !evidence.ValidID(value) {
			return false
		}
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}

func validUniqueReturnActions(values []string) bool {
	if len(values) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		switch value {
		case "complete", "abort", "request-agent":
		default:
			return false
		}
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}

func validUniqueHashSet(values []string, requireNonEmpty bool) bool {
	if requireNonEmpty && len(values) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !evidence.ValidHash(value) {
			return false
		}
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}

func cloneStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return append([]string{}, values...)
}
