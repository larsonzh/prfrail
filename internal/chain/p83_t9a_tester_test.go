package chain

// P83 FIX ROUND 5 counterexample (item T9a): the evidence-root binding of an
// approve decision belongs to the paths that PUBLISH, which review in "policy"
// mode. The generic manual receipt keeps its pre-slice shape, so a manual approve
// without the binding must still form a receipt.
//
// Mutation check: dropping the `decision.ReviewMode == "policy"` condition (making
// the binding unconditional) refuses the manual approve and turns this red.

import (
	"errors"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestP83T9aManualReviewReceiptNeedsNoEvidenceRootBinding(t *testing.T) {
	request := ReviewRequest{
		RunID:                 "run-one",
		TaskID:                "task-one",
		Attempt:               1,
		ParentSnapshotHash:    baselineHash,
		CandidateSnapshotHash: acceptedOne,
		EvidenceRootHash:      stopHash,
		CandidateProducer:     evidence.Actor{Type: "agent", ID: "agent-one"},
	}
	manual := ReviewDecision{
		ReceiptID:  "review-manual-nine",
		OccurredAt: "2026-09-07T01:02:03.000Z",
		RecordedBy: evidence.Actor{Type: "operator", ID: "reviewer-nine"},
		ReviewMode: "manual",
		Outcome:    "approve",
		// Deliberately WITHOUT ReviewEvidenceBinding: the manual receipt never
		// carried it, and the compatibility guarantee is that it still does not.
		Evidence: []string{evidence.Digest("", []byte("manual-observation"))},
	}
	result, err := buildReviewResult(request, manual, fixedClock, func() string { return "event-nine" })
	if err != nil {
		t.Fatalf("a manual approve without the evidence-root binding must still form a receipt, got %v", err)
	}
	if result.ReceiptHash == "" || result.Receipt.Outcome != "approve" || result.Receipt.ReviewMode != "manual" {
		t.Fatalf("unexpected manual review result: %+v", result)
	}

	// The policy path keeps the requirement: the same decision in "policy" mode is
	// still refused, so the relaxation is scoped to the manual receipt only.
	policy := manual
	policy.ReceiptID = "review-policy-nine"
	policy.RecordedBy = evidence.Actor{Type: "policy", ID: "policy-nine"}
	policy.ReviewMode = "policy"
	policy.PolicyHash = stopHash
	refused, err := buildReviewResult(request, policy, fixedClock, func() string { return "event-nine" })
	if !errors.Is(err, ErrInvalidState) {
		t.Fatalf("a policy approve without the evidence-root binding must be refused, got %v", err)
	}
	if refused.ReceiptHash != "" {
		t.Fatalf("a refused policy review must produce no receipt, got %+v", refused)
	}
}
