package chain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestBuildReviewResultApprove(t *testing.T) {
	request := ReviewRequest{
		RunID:                 "run-one",
		TaskID:                "task-one",
		Attempt:               1,
		ParentSnapshotHash:    baselineHash,
		CandidateSnapshotHash: acceptedOne,
		EvidenceRootHash:      stopHash,
		CandidateProducer:     evidence.Actor{Type: "agent", ID: "agent-one"},
	}
	decision := ReviewDecision{
		ReceiptID:  "review-one",
		OccurredAt: "2026-09-07T01:02:03.000Z",
		RecordedBy: evidence.Actor{Type: "operator", ID: "reviewer-one"},
		ReviewMode: "manual",
		Outcome:    "approve",
	}
	result, err := buildReviewResult(request, decision, fixedClock, func() string { return "event-one" })
	if err != nil {
		t.Fatal(err)
	}
	if result.ReceiptHash == "" || result.Receipt.Outcome != "approve" {
		t.Fatalf("unexpected review result: %+v", result)
	}
}

func TestBuildReviewResultWaiverRequiresValidAuthorization(t *testing.T) {
	request := ReviewRequest{
		RunID:                 "run-one",
		TaskID:                "task-one",
		Attempt:               1,
		ParentSnapshotHash:    baselineHash,
		CandidateSnapshotHash: acceptedOne,
		EvidenceRootHash:      stopHash,
	}
	decision := ReviewDecision{
		ReceiptID:     "review-waive",
		OccurredAt:    "2026-09-07T01:02:03.000Z",
		RecordedBy:    evidence.Actor{Type: "operator", ID: "reviewer-one"},
		ReviewMode:    "manual",
		Outcome:       "waive",
		Reason:        &ReviewReason{Code: "approved-exception"},
		ErrorEvidence: []string{stopHash},
		WaiverAuthorization: &WaiverAuthorization{
			AuthorizedBy:    evidence.Actor{Type: "operator", ID: "release-manager"},
			PolicyBasisHash: stopHash,
			Scope:           []string{"task-one"},
			ExpiresAt:       fixedClock().UTC().Add(time.Hour).Format("2006-01-02T15:04:05.000Z"),
		},
	}
	if _, err := buildReviewResult(request, decision, fixedClock, func() string { return "event-one" }); err != nil {
		t.Fatal(err)
	}
	decision.WaiverAuthorization.ExpiresAt = fixedClock().UTC().Add(-time.Minute).Format("2006-01-02T15:04:05.000Z")
	if _, err := buildReviewResult(request, decision, fixedClock, func() string { return "event-one" }); err == nil {
		t.Fatal("expected expired waiver rejection")
	}
}

func TestBuiltReceiptsSerializeToSchemaShape(t *testing.T) {
	request := ReviewRequest{
		RunID:                 "run-one",
		TaskID:                "task-one",
		Attempt:               1,
		ParentSnapshotHash:    baselineHash,
		CandidateSnapshotHash: acceptedOne,
		EvidenceRootHash:      stopHash,
	}
	decision := ReviewDecision{
		ReceiptID:     "review-reject",
		OccurredAt:    "2026-09-07T01:02:03.000Z",
		RecordedBy:    evidence.Actor{Type: "operator", ID: "reviewer-one"},
		ReviewMode:    "manual",
		Outcome:       "reject",
		Reason:        &ReviewReason{Code: "docs-missing"},
		ErrorEvidence: []string{stopHash},
	}
	result, err := buildReviewResult(request, decision, fixedClock, func() string { return "event-one" })
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(result.Receipt)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil {
		t.Fatal(err)
	}
	if _, ok := fields["receiptId"]; !ok {
		t.Fatal("missing receiptId key")
	}
	reason, ok := fields["reason"].(map[string]any)
	if !ok || reason["code"] != "docs-missing" {
		t.Fatalf("unexpected reason shape: %v", fields["reason"])
	}
	if _, exists := reason["message"]; exists {
		t.Fatal("empty message must be omitted for schema minLength")
	}
	if _, exists := reason["Code"]; exists {
		t.Fatal("reason keys must be lower-case")
	}
	if _, exists := fields["policyHash"]; !exists {
		t.Fatal("missing policyHash key")
	}
	if _, exists := fields["waiverAuthorization"]; !exists {
		t.Fatal("missing waiverAuthorization key")
	}
}
