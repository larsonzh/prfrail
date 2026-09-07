package chain

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const reviewReceiptDomain = "proofrail:review-receipt:1\n"

type reviewReceipt struct {
	ReceiptID           string         `json:"receiptId"`
	OccurredAt          string         `json:"occurredAt"`
	RecordedBy          evidence.Actor `json:"recordedBy"`
	RunID               string         `json:"runId"`
	TaskID              string         `json:"taskId"`
	Attempt             int            `json:"attempt"`
	EvidenceRootHash    string         `json:"evidenceRootHash"`
	ReviewMode          string         `json:"reviewMode"`
	PolicyHash          any            `json:"policyHash"`
	Outcome             string         `json:"outcome"`
	Evidence            []string       `json:"evidence"`
	ErrorEvidence       []string       `json:"errorEvidence"`
	Reason              any            `json:"reason"`
	WaiverAuthorization any            `json:"waiverAuthorization"`
}

type reviewResult struct {
	Receipt     reviewReceipt
	ReceiptHash string
}

func buildReviewResult(request ReviewRequest, decision ReviewDecision, clock Clock, ids IDSource) (reviewResult, error) {
	if !evidence.ValidID(request.RunID) || !evidence.ValidID(request.TaskID) || request.Attempt < 1 {
		return reviewResult{}, fmt.Errorf("%w: invalid review request identity", ErrInvalidState)
	}
	if !evidence.ValidHash(request.CandidateSnapshotHash) || !evidence.ValidHash(request.EvidenceRootHash) {
		return reviewResult{}, fmt.Errorf("%w: invalid review request hashes", ErrInvalidState)
	}
	if decision.CandidateSnapshotHash == "" {
		decision.CandidateSnapshotHash = request.CandidateSnapshotHash
	}
	if decision.CandidateSnapshotHash != request.CandidateSnapshotHash {
		return reviewResult{}, fmt.Errorf("%w: review candidate mismatch", ErrInvalidState)
	}
	if decision.ReceiptID == "" {
		decision.ReceiptID = ids()
	}
	if decision.OccurredAt == "" {
		decision.OccurredAt = clock().UTC().Format("2006-01-02T15:04:05.000Z")
	}
	if err := validateReviewDecision(request, decision, clock().UTC()); err != nil {
		return reviewResult{}, err
	}

	receipt := reviewReceipt{
		ReceiptID:        decision.ReceiptID,
		OccurredAt:       decision.OccurredAt,
		RecordedBy:       decision.RecordedBy,
		RunID:            request.RunID,
		TaskID:           request.TaskID,
		Attempt:          request.Attempt,
		EvidenceRootHash: request.EvidenceRootHash,
		ReviewMode:       decision.ReviewMode,
		Outcome:          decision.Outcome,
		Evidence:         uniqueHashes(decision.Evidence),
		ErrorEvidence:    uniqueHashes(decision.ErrorEvidence),
	}
	if decision.PolicyHash == "" {
		receipt.PolicyHash = nil
	} else {
		receipt.PolicyHash = decision.PolicyHash
	}
	if decision.Reason == nil {
		receipt.Reason = nil
	} else {
		receipt.Reason = *decision.Reason
	}
	if decision.WaiverAuthorization == nil {
		receipt.WaiverAuthorization = nil
	} else {
		receipt.WaiverAuthorization = *decision.WaiverAuthorization
	}
	hash, err := digestRecord(reviewReceiptDomain, receipt)
	if err != nil {
		return reviewResult{}, err
	}
	return reviewResult{Receipt: receipt, ReceiptHash: hash}, nil
}

func validateReviewDecision(request ReviewRequest, decision ReviewDecision, now time.Time) error {
	if !evidence.ValidID(decision.ReceiptID) || !evidence.ValidID(decision.RecordedBy.ID) {
		return fmt.Errorf("%w: invalid review receipt identity", ErrInvalidState)
	}
	if decision.RecordedBy.Type != "operator" && decision.RecordedBy.Type != "policy" {
		return fmt.Errorf("%w: invalid review actor type", ErrInvalidState)
	}
	if request.CandidateProducer.ID != "" && decision.RecordedBy.ID == request.CandidateProducer.ID {
		return fmt.Errorf("%w: review separation violated", ErrInvalidState)
	}
	if _, err := time.Parse("2006-01-02T15:04:05.000Z", decision.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid review timestamp", ErrInvalidState)
	}
	if !allHashes(decision.Evidence) || !allHashes(decision.ErrorEvidence) {
		return fmt.Errorf("%w: invalid review evidence", ErrInvalidState)
	}
	switch decision.ReviewMode {
	case "manual":
		if decision.PolicyHash != "" {
			return fmt.Errorf("%w: manual review must not set policy hash", ErrInvalidState)
		}
	case "policy":
		if !evidence.ValidHash(decision.PolicyHash) {
			return fmt.Errorf("%w: policy review requires policy hash", ErrInvalidState)
		}
	default:
		return fmt.Errorf("%w: unknown review mode", ErrInvalidState)
	}
	switch decision.Outcome {
	case "approve":
		if decision.Reason != nil || decision.WaiverAuthorization != nil || len(decision.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: invalid approve review payload", ErrInvalidState)
		}
	case "reject":
		if decision.Reason == nil || decision.WaiverAuthorization != nil || len(decision.ErrorEvidence) == 0 {
			return fmt.Errorf("%w: invalid reject review payload", ErrInvalidState)
		}
		if !evidence.ValidID(decision.Reason.Code) || len(decision.Reason.Message) > 1024 {
			return fmt.Errorf("%w: invalid reject reason", ErrInvalidState)
		}
	case "waive":
		if decision.Reason == nil || decision.WaiverAuthorization == nil || len(decision.ErrorEvidence) == 0 {
			return fmt.Errorf("%w: invalid waiver payload", ErrInvalidState)
		}
		if !evidence.ValidID(decision.Reason.Code) || len(decision.Reason.Message) > 1024 {
			return fmt.Errorf("%w: invalid waiver reason", ErrInvalidState)
		}
		if decision.WaiverAuthorization.AuthorizedBy.Type != "operator" || !evidence.ValidID(decision.WaiverAuthorization.AuthorizedBy.ID) {
			return fmt.Errorf("%w: invalid waiver authorizer", ErrInvalidState)
		}
		if !evidence.ValidHash(decision.WaiverAuthorization.PolicyBasisHash) {
			return fmt.Errorf("%w: invalid waiver policy basis", ErrInvalidState)
		}
		if len(decision.WaiverAuthorization.Scope) == 0 || !slices.Contains(decision.WaiverAuthorization.Scope, request.TaskID) {
			return fmt.Errorf("%w: waiver scope does not include task", ErrInvalidState)
		}
		expiresAt, err := time.Parse("2006-01-02T15:04:05.000Z", decision.WaiverAuthorization.ExpiresAt)
		if err != nil || !expiresAt.After(now) {
			return fmt.Errorf("%w: waiver expired", ErrInvalidState)
		}
	default:
		return fmt.Errorf("%w: unknown review outcome", ErrInvalidState)
	}
	return nil
}

func digestRecord(domain string, body any) (string, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	canonical, err := evidence.Canonicalize(payload)
	if err != nil {
		return "", err
	}
	return evidence.Digest(domain, canonical), nil
}

func allHashes(items []string) bool {
	for _, item := range items {
		if !evidence.ValidHash(item) {
			return false
		}
	}
	return true
}
