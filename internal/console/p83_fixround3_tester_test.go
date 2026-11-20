package console

// P83 fix round 3: the publish-time counterexamples for F3 on both assemblies.
//
// The review decides under a policy and freezes a candidate; Publish must re-prove
// both facts immediately before promoting. Every test below changes the fact after
// the review and asserts the refusal the contract demands (failed outcome, no
// accepted snapshot, no lease evidence), with the unchanged fact as the positive
// control in the same runtime.
//
// Mutation contracts:
//   - TestP83F3AgentPublishRefusesGrantRevokedAfterReview: deleting the
//     verifyReviewPolicyBinding call from agentRunnerRuntime.Publish turns the
//     revoked leg green and this red.
//   - TestP83F3SupervisedPublishRefusesConfirmationWithdrawnAfterReview: deleting
//     the same call from supervisedRuntime.Publish does the same.
//   - TestP83F3SupervisedPublishRefusesCandidateEvidenceRootDrift: deleting the
//     re-derived evidence-root comparison from supervisedRuntime.Publish turns the
//     drift leg green and this red.

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// TestP83F3AgentPublishRefusesGrantRevokedAfterReview is the F3 counterexample on
// the AgentRunner assembly: the operator revokes the grant that backed the review
// after the review approved, and the promotion must refuse instead of completing.
func TestP83F3AgentPublishRefusesGrantRevokedAfterReview(t *testing.T) {
	ctx := context.Background()
	runtime, runRoot := p83PublishRuntime(t)
	proof, err := claimPublishLease(ctx, runRoot, p83RunID)
	if err != nil {
		t.Fatal(err)
	}
	runtime.publishLease = proof
	grant := runtime.authorizationLedger.Grants[0]
	evidenceRoot, err := runtime.evidenceRootHash(p83TaskID, 1)
	if err != nil {
		t.Fatal(err)
	}
	request := chain.PromotionRequest{
		RunID:                 p83RunID,
		TaskID:                p83TaskID,
		Attempt:               1,
		CandidateSnapshotHash: evidence.Digest("", []byte("candidate")),
		EvidenceRootHash:      evidenceRoot,
		ReviewReceiptHash:     evidence.Digest("", []byte("review-receipt")),
	}

	// Positive control: the unchanged grant still publishes.
	decision, err := runtime.Publish(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "completed" || len(decision.LeaseEvidence) != 1 {
		t.Fatalf("an unchanged grant must publish, got %+v", decision)
	}

	// Revocation lands after the review: the grant is still present, but no longer
	// active, so the publish-time verification must refuse.
	revocation, err := chain.NewRevocationAuthorizationRecord(chain.AuthorizationRevocation{
		RecordID:          "p83-f3-agent-revocation",
		Kind:              "revocation",
		IssuedAt:          p83Clock().Add(-time.Minute).UTC().Format(timestampLayout),
		IssuedBy:          evidence.Actor{Type: "operator", ID: "p83-operator"},
		AuthorizationID:   grant.Grant.AuthorizationID,
		AuthorizationHash: grant.RecordHash,
		ReasonCode:        "operator-revoked",
		StopDisposition:   "not-required",
	}, p83Clock, func() string { return "p83-f3-agent-revocation-record" })
	if err != nil {
		t.Fatal(err)
	}
	runtime.authorizationLedger = chain.AuthorizationLedger{
		Grants:      []chain.AuthorizationRecord{grant},
		Revocations: []chain.AuthorizationRecord{revocation},
	}

	decision, err = runtime.Publish(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "failed" {
		t.Fatalf("a grant revoked after the review must refuse the promotion, got %+v", decision)
	}
	if decision.AcceptedSnapshotHash != "" || len(decision.LeaseEvidence) != 0 {
		t.Fatalf("a refused promotion must carry no accepted snapshot and no lease evidence, got %+v", decision)
	}
	if len(decision.ErrorEvidence) == 0 {
		t.Fatal("a refused promotion must carry error evidence")
	}
}

// p83F3SupervisedRuntime rebuilds the supervised runtime over a run root a harness
// already drove to the open work window, so the ports read the artifacts the
// production window wrote.
func p83F3SupervisedRuntime(t *testing.T, harness *supervisedHarness) (*supervisedRuntime, *supervisedTaskState) {
	t.Helper()
	clock, _ := supervisedClockAndIDs(harness.options)
	plan, err := supervisedAssemblyPlan(harness.options)
	if err != nil {
		t.Fatalf("the harness plan must bind: %v", err)
	}
	runtime, err := newSupervisedRuntime(harness.options, plan, clock)
	if err != nil {
		t.Fatalf("the supervised runtime must assemble: %v", err)
	}
	state, err := runtime.stateFor(supervisedTestTaskID, 1)
	if err != nil {
		t.Fatalf("the harness task state must resolve: %v", err)
	}
	return runtime, state
}

// p83F3SupervisedReview freezes a candidate and binds the policy hash the review
// would have bound, returning the promotion request the chain would hand Publish.
func p83F3SupervisedReview(t *testing.T, harness *supervisedHarness, runtime *supervisedRuntime, state *supervisedTaskState) chain.PromotionRequest {
	t.Helper()
	ctx := context.Background()
	candidate, err := runtime.Accept(ctx, supervisedTestRunID, supervisedTestTaskID, 1,
		chain.SnapshotRef{Hash: harness.loadBeforeManifest().ManifestHash},
		chain.Workspace{Root: harness.workspace()})
	if err != nil {
		t.Fatalf("freezing the candidate must succeed: %v", err)
	}
	session, found, err := runtime.loadSession(state)
	if err != nil || !found {
		t.Fatalf("the open window must persist its handoff session: found=%v err=%v", found, err)
	}
	ledger, err := runtime.loadConfirmationLedger()
	if err != nil {
		t.Fatal(err)
	}
	// T2 (finding F3): the risk confirmation is still a separate requirement of the
	// window (both ports refuse without it), but it no longer shapes the policy
	// hash itself - that is the frozen handoff policy digest.
	if _, ok := supervisedConfirmationForBinding(ledger, supervisedTestRunID, supervisedTestTaskID, supervisedTestStepID); !ok {
		t.Fatal("the window must persist the operator risk confirmation")
	}
	state.reviewPolicyHash = supervisedReviewPolicyHash(session)
	return chain.PromotionRequest{
		RunID:                 supervisedTestRunID,
		TaskID:                supervisedTestTaskID,
		Attempt:               1,
		CandidateSnapshotHash: candidate.Snapshot.Hash,
		EvidenceRootHash:      candidate.EvidenceRootHash,
		ReviewReceiptHash:     evidence.Digest("", []byte("supervised-review-receipt")),
	}
}

// TestP83F3SupervisedPublishRefusesConfirmationWithdrawnAfterReview is the F3
// counterexample on the supervised assembly: the operator risk confirmation that
// the review bound is withdrawn before the promotion, so Publish must refuse.
func TestP83F3SupervisedPublishRefusesConfirmationWithdrawnAfterReview(t *testing.T) {
	ctx := context.Background()
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.now = p83FixedNow()
	harness.driveToWindow()
	runtime, state := p83F3SupervisedRuntime(t, harness)
	request := p83F3SupervisedReview(t, harness, runtime, state)

	// Positive control: the unchanged confirmation still publishes.
	decision, err := runtime.Publish(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "completed" || len(decision.LeaseEvidence) != 1 {
		t.Fatalf("an unchanged risk confirmation must publish, got %+v", decision)
	}

	// Withdraw the operator confirmation after the review.
	confirmations := runtime.confirmationsPath()
	saved, err := os.ReadFile(confirmations)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(confirmations); err != nil {
		t.Fatal(err)
	}
	decision, err = runtime.Publish(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "failed" {
		t.Fatalf("a withdrawn confirmation must refuse the promotion, got %+v", decision)
	}
	if decision.AcceptedSnapshotHash != "" || len(decision.LeaseEvidence) != 0 || len(decision.ErrorEvidence) == 0 {
		t.Fatalf("a refused promotion must carry no accepted snapshot, no lease evidence and error evidence, got %+v", decision)
	}

	// Restoring the same confirmation publishes again: the refusal came from the
	// withdrawal and nothing else.
	if err := os.WriteFile(confirmations, saved, 0o600); err != nil {
		t.Fatal(err)
	}
	decision, err = runtime.Publish(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "completed" || len(decision.LeaseEvidence) != 1 {
		t.Fatalf("the restored confirmation must publish again, got %+v", decision)
	}
}

// TestP83F3SupervisedPublishRefusesCandidateEvidenceRootDrift is the second half
// of F3 on the supervised assembly: the frozen candidate drifts after the review,
// so the re-derived evidence root no longer matches the reviewed one.
func TestP83F3SupervisedPublishRefusesCandidateEvidenceRootDrift(t *testing.T) {
	ctx := context.Background()
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.now = p83FixedNow()
	harness.driveToWindow()
	runtime, state := p83F3SupervisedRuntime(t, harness)
	request := p83F3SupervisedReview(t, harness, runtime, state)

	// The candidate is re-frozen to a different manifest after the review.
	state.acceptedHash = evidence.Digest("", []byte("p83-drifted-candidate"))
	drifted := supervisedCandidateEvidenceRoot(runtime.runID, supervisedTestTaskID, 1, state.acceptedParentHash, state.acceptedHash)
	if drifted == request.EvidenceRootHash {
		t.Fatal("test precondition failed: the drifted candidate must derive a different evidence root")
	}

	decision, err := runtime.Publish(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "failed" {
		t.Fatalf("a drifted candidate evidence root must refuse the promotion, got %+v", decision)
	}
	wantEvidence := digest(supervisedPromotionBlockedDomain, "the candidate evidence root changed after the review")
	if len(decision.ErrorEvidence) != 1 || decision.ErrorEvidence[0] != wantEvidence {
		t.Fatalf("the refusal must be the evidence-root drift branch, got %v want %v", decision.ErrorEvidence, wantEvidence)
	}
	if decision.AcceptedSnapshotHash != "" || len(decision.LeaseEvidence) != 0 {
		t.Fatalf("a refused promotion must carry no accepted snapshot and no lease evidence, got %+v", decision)
	}
}
