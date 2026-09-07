package chain

import "testing"

func TestBuildPromotionResultCompleted(t *testing.T) {
	request := PromotionRequest{
		RunID:                 "run-one",
		TaskID:                "task-one",
		Attempt:               1,
		ParentSnapshotHash:    baselineHash,
		CandidateSnapshotHash: acceptedOne,
		EvidenceRootHash:      stopHash,
		ReviewReceiptHash:     acceptedTwo,
	}
	decision := PromotionDecision{
		ReceiptID:            "promotion-one",
		OccurredAt:           "2026-09-07T01:02:03.000Z",
		Outcome:              "completed",
		AcceptedSnapshotHash: acceptedOne,
		WriterStopEvidence:   []string{stopHash},
		LeaseEvidence:        []string{stopHash},
	}
	result, err := buildPromotionResult(request, decision, fixedClock, func() string { return "event-one" })
	if err != nil {
		t.Fatal(err)
	}
	if result.ReceiptHash == "" || result.Receipt.Outcome != "completed" {
		t.Fatalf("unexpected promotion result: %+v", result)
	}
}

func TestBuildPromotionResultRejectsMismatchedAcceptance(t *testing.T) {
	request := PromotionRequest{
		RunID:                 "run-one",
		TaskID:                "task-one",
		Attempt:               1,
		ParentSnapshotHash:    baselineHash,
		CandidateSnapshotHash: acceptedOne,
		EvidenceRootHash:      stopHash,
		ReviewReceiptHash:     acceptedTwo,
	}
	decision := PromotionDecision{
		ReceiptID:            "promotion-one",
		OccurredAt:           "2026-09-07T01:02:03.000Z",
		Outcome:              "completed",
		AcceptedSnapshotHash: acceptedTwo,
		WriterStopEvidence:   []string{stopHash},
		LeaseEvidence:        []string{stopHash},
	}
	if _, err := buildPromotionResult(request, decision, fixedClock, func() string { return "event-one" }); err == nil {
		t.Fatal("expected mismatched accepted hash rejection")
	}
}
