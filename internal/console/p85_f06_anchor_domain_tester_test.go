package console

// P85 review r1 (new Medium), design rev 3 section 6 / test row M10: Review and
// Publish must re-verify the anchor's managed-identity write proof against a FRESH
// observation of THIS policy domain, not against the digest's FORMAT.
//
// The counterexample has to be a well-formed digest, so that the format check
// (evidence.ValidHash) cannot be what refuses it: the anchor is rewritten with the
// observation of ANOTHER directory - a layout this identity can write, so its
// observation is a real, valid digest of a different write domain.
//
// Mutation check: dropping the agentRunnerReverifyManagedWriteProof call from
// agentRunnerVerifyFrozenPolicyAnchor (back to the ValidHash-only check) makes the
// review - and the publish - approve, turning both legs red.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85F06ForeignDomainProof returns a WELL-FORMED digest produced from a different
// write domain. The observation is taken from the real primitive on a throwaway
// directory, so the value is a genuine proof of another layout rather than a random
// hash; hosts where the primitive cannot even produce a digest fall back to a
// well-formed constant, which still passes the format check.
func p85F06ForeignDomainProof(t *testing.T) string {
	t.Helper()
	directory := filepath.Join(t.TempDir(), "another-policy-domain")
	if err := os.MkdirAll(directory, 0o755); err != nil {
		t.Fatalf("create the foreign policy domain: %v", err)
	}
	observed, _ := adapters.VerifyManagedAgentCannotWritePolicyDomain(directory, nil)
	if observed == "" || observed == assemblyManagedWriteProof || !evidence.ValidHash(observed) {
		observed = "sha256:" + strings.Repeat("3", 64)
	}
	return observed
}

func p85F06RequireAnchorProofMismatch(t *testing.T, err error) {
	t.Helper()
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("an anchor whose proof describes another domain must be refused with ErrInvalidConfig, got %v", err)
	}
	if !strings.Contains(err.Error(), "is not the proof of the domain") {
		t.Fatalf("the refusal must come from the fresh observation, not from the digest format, got %v", err)
	}
}

func TestP85F06AnchorProofFromAnotherDomainRefusesReview(t *testing.T) {
	runtime, _ := p85FrozenReviewRuntime(t)
	anchorPath := agentRunnerFrozenPolicyAnchorPath(runtime.runRoot, p83TaskID, p83StepID, 1)
	foreign := p85F06ForeignDomainProof(t)
	original, _ := p85RewriteAnchorProof(t, runtime, foreign)

	decision, err := runtime.Review(context.Background(), p85ReviewRequest())
	p85F06RequireAnchorProofMismatch(t, err)
	if decision.Outcome != "" || decision.PolicyHash != "" || len(decision.Evidence) != 0 {
		t.Fatalf("a refused review must produce no decision, got %+v", decision)
	}

	// The anchor is the sole cause: restoring the bytes approves the same runtime.
	if err := os.WriteFile(anchorPath, original, 0o600); err != nil {
		t.Fatalf("restore the frozen policy anchor: %v", err)
	}
	restored, err := runtime.Review(context.Background(), p85ReviewRequest())
	if err != nil {
		t.Fatalf("the restored anchor must approve, got %v", err)
	}
	if restored.Outcome != "approve" {
		t.Fatalf("restoring the anchor must approve the review, got %+v", restored)
	}
}

func TestP85F06AnchorProofFromAnotherDomainRefusesPublish(t *testing.T) {
	runtime, runRoot := p83PublishRuntime(t)
	proof, err := claimPublishLease(context.Background(), runRoot, p83RunID)
	if err != nil {
		t.Fatal(err)
	}
	runtime.publishLease = proof
	if _, err := runtime.Review(context.Background(), p85ReviewRequest()); err != nil {
		t.Fatalf("the review must approve before the publish leg: %v", err)
	}
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
	anchorPath := agentRunnerFrozenPolicyAnchorPath(runtime.runRoot, p83TaskID, p83StepID, 1)
	artifactBefore := p85ReadFrozenPolicyBytes(t, runtime)
	original, _ := p85RewriteAnchorProof(t, runtime, p85F06ForeignDomainProof(t))

	decision, err := runtime.Publish(context.Background(), request)
	if err != nil {
		t.Fatalf("a blocked promotion is a decision, not an error: %v", err)
	}
	if decision.Outcome != "failed" || decision.AcceptedSnapshotHash != "" || len(decision.LeaseEvidence) != 0 || len(decision.WriterStopEvidence) != 0 {
		t.Fatalf("an anchor proven for another domain must block the promotion, got %+v", decision)
	}
	if len(decision.ErrorEvidence) == 0 {
		t.Fatal("a blocked promotion must carry error evidence")
	}
	if after := p85ReadFrozenPolicyBytes(t, runtime); string(after) != string(artifactBefore) {
		t.Fatalf("the artifact must stay byte-identical, got %q want %q", after, artifactBefore)
	}

	if err := os.WriteFile(anchorPath, original, 0o600); err != nil {
		t.Fatalf("restore the frozen policy anchor: %v", err)
	}
	restored, err := runtime.Publish(context.Background(), request)
	if err != nil {
		t.Fatalf("the restored anchor must be publishable, got %v", err)
	}
	if restored.Outcome != "completed" {
		t.Fatalf("restoring the anchor must complete the promotion, got %+v", restored)
	}
}
