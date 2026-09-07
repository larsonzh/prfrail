package chain

import (
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const promotionReceiptDomain = "proofrail:promotion-receipt:1\n"

type promotionReceipt struct {
	ReceiptID             string         `json:"receiptId"`
	OccurredAt            string         `json:"occurredAt"`
	RecordedBy            evidence.Actor `json:"recordedBy"`
	RunID                 string         `json:"runId"`
	TaskID                string         `json:"taskId"`
	Attempt               int            `json:"attempt"`
	ParentSnapshotHash    string         `json:"parentSnapshotHash"`
	CandidateSnapshotHash string         `json:"candidateSnapshotHash"`
	EvidenceRootHash      string         `json:"evidenceRootHash"`
	ReviewReceiptHash     string         `json:"reviewReceiptHash"`
	WriterStopEvidence    []string       `json:"writerStopEvidence"`
	LeaseEvidence         []string       `json:"leaseEvidence"`
	Outcome               string         `json:"outcome"`
	AcceptedSnapshotHash  any            `json:"acceptedSnapshotHash"`
	Evidence              []string       `json:"evidence"`
	ErrorEvidence         []string       `json:"errorEvidence"`
}

type promotionResult struct {
	Receipt     promotionReceipt
	ReceiptHash string
}

func buildPromotionResult(request PromotionRequest, decision PromotionDecision, clock Clock, ids IDSource) (promotionResult, error) {
	if !evidence.ValidID(request.RunID) || !evidence.ValidID(request.TaskID) || request.Attempt < 1 {
		return promotionResult{}, fmt.Errorf("%w: invalid promotion identity", ErrInvalidState)
	}
	if !evidence.ValidHash(request.ParentSnapshotHash) || !evidence.ValidHash(request.CandidateSnapshotHash) || !evidence.ValidHash(request.EvidenceRootHash) || !evidence.ValidHash(request.ReviewReceiptHash) {
		return promotionResult{}, fmt.Errorf("%w: invalid promotion request hashes", ErrInvalidState)
	}
	if decision.ReceiptID == "" {
		decision.ReceiptID = ids()
	}
	if decision.OccurredAt == "" {
		decision.OccurredAt = clock().UTC().Format("2006-01-02T15:04:05.000Z")
	}
	if err := validatePromotionDecision(request, decision); err != nil {
		return promotionResult{}, err
	}

	receipt := promotionReceipt{
		ReceiptID:             decision.ReceiptID,
		OccurredAt:            decision.OccurredAt,
		RecordedBy:            evidence.Actor{Type: "system", ID: "proofrail"},
		RunID:                 request.RunID,
		TaskID:                request.TaskID,
		Attempt:               request.Attempt,
		ParentSnapshotHash:    request.ParentSnapshotHash,
		CandidateSnapshotHash: request.CandidateSnapshotHash,
		EvidenceRootHash:      request.EvidenceRootHash,
		ReviewReceiptHash:     request.ReviewReceiptHash,
		WriterStopEvidence:    uniqueHashes(decision.WriterStopEvidence),
		LeaseEvidence:         uniqueHashes(decision.LeaseEvidence),
		Outcome:               decision.Outcome,
		Evidence:              uniqueHashes(decision.Evidence),
		ErrorEvidence:         uniqueHashes(decision.ErrorEvidence),
	}
	if decision.AcceptedSnapshotHash == "" {
		receipt.AcceptedSnapshotHash = nil
	} else {
		receipt.AcceptedSnapshotHash = decision.AcceptedSnapshotHash
	}
	hash, err := digestRecord(promotionReceiptDomain, receipt)
	if err != nil {
		return promotionResult{}, err
	}
	return promotionResult{Receipt: receipt, ReceiptHash: hash}, nil
}

func validatePromotionDecision(request PromotionRequest, decision PromotionDecision) error {
	if !evidence.ValidID(decision.ReceiptID) {
		return fmt.Errorf("%w: invalid promotion receipt id", ErrInvalidState)
	}
	if _, err := time.Parse("2006-01-02T15:04:05.000Z", decision.OccurredAt); err != nil {
		return fmt.Errorf("%w: invalid promotion timestamp", ErrInvalidState)
	}
	if len(decision.WriterStopEvidence) == 0 || len(decision.LeaseEvidence) == 0 {
		return fmt.Errorf("%w: promotion requires stop and lease evidence", ErrInvalidState)
	}
	if !allHashes(decision.WriterStopEvidence) || !allHashes(decision.LeaseEvidence) || !allHashes(decision.Evidence) || !allHashes(decision.ErrorEvidence) {
		return fmt.Errorf("%w: invalid promotion evidence", ErrInvalidState)
	}
	switch decision.Outcome {
	case "completed":
		if decision.AcceptedSnapshotHash != request.CandidateSnapshotHash || len(decision.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: invalid completed promotion", ErrInvalidState)
		}
	case "failed", "uncertain":
		if decision.AcceptedSnapshotHash != "" || len(decision.ErrorEvidence) == 0 {
			return fmt.Errorf("%w: invalid failed/uncertain promotion", ErrInvalidState)
		}
	default:
		return fmt.Errorf("%w: unknown promotion outcome", ErrInvalidState)
	}
	return nil
}
