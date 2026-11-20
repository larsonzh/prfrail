package console

// P83 F4 fix round: the publish lease must survive the attempt boundary and the
// TTL, and must stop proving anything the moment the store no longer holds the
// exact record this run claimed.
//
// Mutation contracts:
//   - TestAgentRunnerAttemptTwoAdoptsRunScopedPublishLease: restoring the
//     per-attempt claim plus the permanent refusal on ErrLeaseHeld turns this red.
//   - TestClaimPublishLeaseRefusesRecordHeldByAnotherRun: relaxing the adoption
//     predicate (accept any held record) turns this red.
//   - TestAgentRunnerPublishRefusesStaleLeaseProof: dropping the token/hash
//     re-proof in renewPublishLease and verifyPublishLease turns this red.
//   - TestSupervisedReturnRenewsLeaseAtThreshold: removing the renewal call leaves
//     the fencing token at generation 1 and turns this red.
//   - TestSupervisedReturnRenewsExpiredLeaseInsteadOfRefusing: removing the
//     renewal call makes the promotion refuse an expired lease and turns this red.
//   - TestSupervisedReturnRefusesExpiredLeaseHeldByAnotherRun: renewing whatever
//     record is current (instead of re-proving the claimed token) turns this red.
//   - TestSupervisedReturnRefusesProofOfAnotherRun: removing the
//     publishLeaseProofBelongsToRun guard from renewPublishLease turns this red
//     (the foreign record is current and inside its TTL, so it renews), and
//     removing it from verifyPublishLease turns this red (the verification would
//     bind the foreign lease hash).

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// p83LeaseRuntime builds the per-attempt runtime state Publish needs inside a run
// root the caller owns, so several "processes" of one run can share one durable
// lease store.
func p83LeaseRuntime(t *testing.T, runRoot string, attempt int) *agentRunnerRuntime {
	t.Helper()
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(runRoot)); err != nil {
		t.Fatal(err)
	}
	store, err := adapters.NewAgentRunnerReplayStoreForRun(runRoot, p83RunID)
	if err != nil {
		t.Fatal(err)
	}
	record := adapters.AgentRunnerRequestRecord{Request: adapters.AgentRunnerRequest{
		RequestID: fmt.Sprintf("request-p83-attempt-%d", attempt),
		RunID:     p83RunID,
		TaskID:    p83TaskID,
		StepID:    p83StepID,
		Attempt:   attempt,
	}}
	evidenceDir, err := store.RunEvidenceDir(record.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{agentRunnerParentPIDReportFileName, agentRunnerChildPIDReportFileName} {
		if err := os.WriteFile(filepath.Join(evidenceDir, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	workspace := filepath.Join(runRoot, "workspaces", fmt.Sprintf("p83-attempt-%d", attempt))
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	snapshotStore, err := snapshot.NewStore(filepath.Join(runRoot, "snapshot-store"), 0)
	if err != nil {
		t.Fatal(err)
	}
	capturer, err := adapters.NewAgentRunnerSnapshotCapturer(snapshotStore, p83RunID, p83TaskID, attempt)
	if err != nil {
		t.Fatal(err)
	}
	key := agentRunnerTaskStateKey(p83TaskID, attempt)
	return &agentRunnerRuntime{
		runID:               p83RunID,
		runRoot:             runRoot,
		store:               store,
		requestRecord:       record,
		requestRecords:      map[string]adapters.AgentRunnerRequestRecord{key: record},
		authorizationLedger: chain.AuthorizationLedger{},
		clock:               p83Clock,
		workspaceRoots:      map[string]string{key: workspace},
		capturers:           map[string]*adapters.AgentRunnerSnapshotCapturer{key: capturer},
		acceptedEntries:     map[string][]snapshot.Entry{key: nil},
	}
}

func p83Promotion(attempt int) chain.PromotionRequest {
	return chain.PromotionRequest{
		RunID:                 p83RunID,
		TaskID:                p83TaskID,
		Attempt:               attempt,
		CandidateSnapshotHash: evidence.Digest("", []byte("candidate")),
		EvidenceRootHash:      evidence.Digest("", []byte("evidence-root")),
		ReviewReceiptHash:     evidence.Digest("", []byte("review-receipt")),
	}
}

// TestAgentRunnerAttemptTwoAdoptsRunScopedPublishLease is the F4 mixed-attempt
// liveness case: attempt 1 of a run claims the run-scoped lease and attempt 2 (a
// new process sharing the run root) re-enters, adopts that same lease instead of
// being refused for the rest of the run, and publishes on it.
func TestAgentRunnerAttemptTwoAdoptsRunScopedPublishLease(t *testing.T) {
	runRoot := t.TempDir()
	ctx := context.Background()
	first, err := claimPublishLease(ctx, runRoot, p83RunID)
	if err != nil {
		t.Fatalf("attempt 1 must claim the run lease: %v", err)
	}
	second, err := claimPublishLease(ctx, runRoot, p83RunID)
	if err != nil {
		t.Fatalf("attempt 2 must re-adopt the run lease instead of being refused: %v", err)
	}
	// T5 (P83 fix round 6): adoption is a transfer, but only when the record names a
	// DIFFERENT holder. Two call sites of one process are one holder - no OS liveness
	// proof can tell them apart and no transfer would separate anything - so this
	// re-entry keeps the record and every proof written for the run stays valid. A
	// successor PROCESS does transfer: it reclaims the record under its own identity
	// at generation+1 (guard.LeaseStore.Reclaim), which is what makes the second
	// successor of one stopped holder fail with ErrFencingTokenMismatch.
	if second != first {
		t.Fatalf("a re-entry of the same holder must reuse the run lease, got %+v want %+v", second, first)
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Validate(ctx, second.ResourceID, second.Token); err != nil {
		t.Fatalf("the reused generation must stay valid for its holder, got %v", err)
	}
	runtime := p83LeaseRuntime(t, runRoot, 2)
	runtime.publishLease = second
	// F3: Publish re-proves the reviewed policy and the frozen evidence root, so
	// the fixture seeds both instead of asserting completion without them.
	p83BindReviewedPolicy(t, runtime, p83TaskID, 2)
	request := p83Promotion(2)
	if request.EvidenceRootHash, err = runtime.evidenceRootHash(p83TaskID, 2); err != nil {
		t.Fatal(err)
	}
	decision, err := runtime.Publish(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "completed" {
		t.Fatalf("attempt 2 must publish on the adopted lease, got %+v", decision)
	}
	if len(decision.LeaseEvidence) != 1 || decision.LeaseEvidence[0] != second.LeaseHash {
		t.Fatalf("the promotion must bind the lease hash the attempt holds %s, got %v", second.LeaseHash, decision.LeaseEvidence)
	}
}

// TestClaimPublishLeaseRefusesRecordHeldByAnotherRun is the negative half of the
// adoption relaxation: a record at this run's resource that was written by
// another process identity is not adopted, so the run never obtains a proof.
func TestClaimPublishLeaseRefusesRecordHeldByAnotherRun(t *testing.T) {
	runRoot := t.TempDir()
	ctx := context.Background()
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	foreign := guard.ProcessIdentity{PID: os.Getpid(), StartToken: "run-foreign"}
	if _, err := store.Claim(ctx, publishLeaseResourceID(p83RunID), publishLeaseConsumerID, "run-foreign-claim", foreign, publishLeaseTTL); err != nil {
		t.Fatalf("seed a foreign holder: %v", err)
	}
	if _, err := claimPublishLease(ctx, runRoot, p83RunID); !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a lease held by another run must not be adopted, got %v", err)
	}
}

// TestAgentRunnerPublishRefusesStaleLeaseProof pins the publish-time re-proof: a
// proof whose generation and hash no longer match the stored record proves
// nothing, so the promotion is refused instead of completing on a stale token.
func TestAgentRunnerPublishRefusesStaleLeaseProof(t *testing.T) {
	runRoot := t.TempDir()
	ctx := context.Background()
	proof, err := claimPublishLease(ctx, runRoot, p83RunID)
	if err != nil {
		t.Fatal(err)
	}
	// A renewal that happened after this attempt wrote its proof (the window
	// outlived the renew threshold) advances the record past the stored proof.
	if _, err := renewPublishLease(ctx, runRoot, p83RunID, proof, p83Clock().Add(publishLeaseTTL)); err != nil {
		t.Fatalf("renew the record: %v", err)
	}
	runtime := p83LeaseRuntime(t, runRoot, 1)
	runtime.publishLease = proof
	decision, err := runtime.Publish(ctx, p83Promotion(1))
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "failed" || len(decision.LeaseEvidence) != 0 || len(decision.ErrorEvidence) == 0 {
		t.Fatalf("a stale lease proof must refuse the promotion, got %+v", decision)
	}
}

// p83SupervisedLeaseProofPath is the store-local proof of the single-task
// supervised harness.
func p83SupervisedLeaseProofPath(runRoot string) string {
	return filepath.Join(SupervisedDir(runRoot), publishLeaseProofName)
}

func p83LoadSupervisedLeaseProof(t *testing.T, runRoot string) publishLeaseProof {
	t.Helper()
	var proof publishLeaseProof
	found, err := readSupervisedJSON(p83SupervisedLeaseProofPath(runRoot), &proof)
	if err != nil || !found {
		t.Fatalf("the window must persist its publish lease proof, found=%v err=%v", found, err)
	}
	return proof
}

// p83FixedNow anchors the harness clock to the wall clock the lease store uses,
// so the TTL arithmetic in the test and in the store agree.
func p83FixedNow() time.Time { return time.Now().UTC().Truncate(time.Second) }

// TestSupervisedReturnRenewsLeaseAtThreshold is the F4 TTL liveness case that is
// still inside the TTL: a run past the renew threshold must renew (the fencing
// token advances) and complete, not sit on a proof that expires mid-window.
func TestSupervisedReturnRenewsLeaseAtThreshold(t *testing.T) {
	// A 100h window so the lease TTL, not the handoff timeout, is what moves.
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 360000000)
	harness.now = p83FixedNow()
	harness.driveToWindow()
	before := p83LoadSupervisedLeaseProof(t, harness.runRoot)
	if before.Token.Generation != 1 {
		t.Fatalf("the window must claim generation 1, got %d", before.Token.Generation)
	}
	harness.advance(13 * time.Hour)
	harness.writeArtifact("artifact.txt", "hello from the supervised agent\n")
	if _, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete"); err != nil {
		t.Fatalf("a run past the renew threshold must still complete, got %v", err)
	}
	after := p83LoadSupervisedLeaseProof(t, harness.runRoot)
	if after.Token.Generation != before.Token.Generation+1 {
		t.Fatalf("the renewal must advance the fencing token, generation %d -> %d", before.Token.Generation, after.Token.Generation)
	}
	if after.LeaseHash == before.LeaseHash {
		t.Fatal("the renewed lease must carry a new lease hash")
	}
}

// TestSupervisedReturnRenewsExpiredLeaseInsteadOfRefusing is the harder TTL
// liveness case: the window legitimately outlived the lease TTL, and the run still
// holds the exact record, so publication renews and completes instead of reading
// as an expired lease. The TTL is shortened through the store's own Renew, so the
// expiry moves without a clock conflict with the store.
func TestSupervisedReturnRenewsExpiredLeaseInsteadOfRefusing(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 360000000)
	harness.now = p83FixedNow()
	harness.driveToWindow()
	ctx := context.Background()
	before := p83LoadSupervisedLeaseProof(t, harness.runRoot)
	store, err := newPublishLeaseStore(harness.runRoot)
	if err != nil {
		t.Fatal(err)
	}
	// T5: the store renews only for the recorded holder, so the shortening renewal
	// presents the identity the record already names.
	held, err := store.Current(ctx, before.ResourceID)
	if err != nil {
		t.Fatalf("read the held lease: %v", err)
	}
	shortened, err := store.Renew(ctx, before.ResourceID, before.Token, held.Process, 30*time.Second)
	if err != nil {
		t.Fatalf("shorten the lease TTL: %v", err)
	}
	renewed := publishLeaseProof{ResourceID: shortened.ResourceID, Token: shortened.Token, LeaseHash: shortened.LeaseHash}
	if err := writeSupervisedJSON(p83SupervisedLeaseProofPath(harness.runRoot), renewed); err != nil {
		t.Fatal(err)
	}
	// The harness clock now reads past the shortened expiry.
	harness.advance(time.Minute)
	harness.writeArtifact("artifact.txt", "hello from the supervised agent\n")
	if _, err := ExecuteSupervisedReturn(ctx, harness.options, "complete"); err != nil {
		t.Fatalf("a run that outlived the TTL but still holds the lease must complete, got %v", err)
	}
	after := p83LoadSupervisedLeaseProof(t, harness.runRoot)
	if after.Token.Generation != renewed.Token.Generation+1 || after.LeaseHash == renewed.LeaseHash {
		t.Fatalf("the expired lease must be renewed and re-tokened, got %+v want generation %d", after, renewed.Token.Generation+1)
	}
}

// TestSupervisedReturnRefusesExpiredLeaseHeldByAnotherRun is the negative TTL
// case: the lease expired and the resource now names another holder's record, so
// no renewal of this run's claim is possible and publication must refuse.
func TestSupervisedReturnRefusesExpiredLeaseHeldByAnotherRun(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 360000000)
	harness.now = p83FixedNow()
	harness.driveToWindow()
	before := p83LoadSupervisedLeaseProof(t, harness.runRoot)
	ctx := context.Background()
	store, err := newPublishLeaseStore(harness.runRoot)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Release(ctx, before.ResourceID, before.Token); err != nil {
		t.Fatalf("release the expired claim: %v", err)
	}
	foreign := guard.ProcessIdentity{PID: os.Getpid(), StartToken: "run-foreign"}
	if _, err := store.Claim(ctx, before.ResourceID, publishLeaseConsumerID, "run-foreign-claim", foreign, 30*time.Second); err != nil {
		t.Fatalf("seed the foreign holder: %v", err)
	}
	// The foreign record is past its own expiry: nothing about it can be renewed
	// on behalf of this run.
	harness.advance(time.Minute)
	harness.writeArtifact("artifact.txt", "hello from the supervised agent\n")
	summary, err := ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if err == nil {
		t.Fatalf("a lease held by another run must not be renewed into a promotion, got %+v", summary)
	}
	if hasStateName(summary.TaskStateCounts, "PASSED") {
		t.Fatalf("a refused promotion must leave no passed task: %+v", summary.TaskStateCounts)
	}
	after := p83LoadSupervisedLeaseProof(t, harness.runRoot)
	if after != before {
		t.Fatalf("a refused promotion must not rewrite this run's proof, got %+v want %+v", after, before)
	}
}

// p83RunRootBindsHash reports whether any durable artifact of the run root
// mentions a lease hash, skipping the lease store and the store-local proofs:
// those are exactly the planted foreign records this test has to leave alone.
func p83RunRootBindsHash(t *testing.T, runRoot, hash string) (string, bool) {
	t.Helper()
	binding := ""
	leaseRoot := publishLeaseRoot(runRoot)
	walkErr := filepath.WalkDir(runRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if path == leaseRoot {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.Name() == publishLeaseProofName {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(data), hash) {
			binding = path
			return filepath.SkipAll
		}
		return nil
	})
	if walkErr != nil && !errors.Is(walkErr, filepath.SkipAll) {
		t.Fatalf("scan the run root for %s: %v", hash, walkErr)
	}
	return binding, binding != ""
}

// TestSupervisedReturnRefusesProofOfAnotherRun is the F4 identity case: another
// run's legitimate lease is issued in the same store root and its proof is
// dropped into this task's publish-lease.json. The token and the lease hash of
// that proof are genuine, so only the run-scoped resource can tell the runs
// apart. Both proof halves must refuse it, and the supervised return must refuse
// to publish instead of binding a foreign lease hash.
//
// The positive control lives next to it: the run's own proof still renews and
// verifies here, and TestSupervisedReturnRenewsLeaseAtThreshold covers the same
// control end to end.
func TestSupervisedReturnRefusesProofOfAnotherRun(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 360000000)
	harness.now = p83FixedNow()
	harness.driveToWindow()
	ctx := context.Background()
	own := p83LoadSupervisedLeaseProof(t, harness.runRoot)

	// The other run is a legitimate holder: it claims its own run-scoped resource
	// in the same store root through the production claim path.
	const foreignRunID = "run-t028-foreign"
	foreign, err := claimPublishLease(ctx, harness.runRoot, foreignRunID)
	if err != nil {
		t.Fatalf("issue the other run's lease in the same store root: %v", err)
	}
	if foreign.ResourceID == own.ResourceID || foreign.LeaseHash == own.LeaseHash {
		t.Fatalf("the planted proof must belong to another run, got %+v (own %+v)", foreign, own)
	}
	if err := writeSupervisedJSON(p83SupervisedLeaseProofPath(harness.runRoot), foreign); err != nil {
		t.Fatal(err)
	}

	// Each proof half refuses the planted record on its own, so a mutation of
	// either guard is visible without the full return path.
	if _, err := renewPublishLease(ctx, harness.runRoot, supervisedTestRunID, foreign, harness.now); !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("renewing another run's proof must be refused, got %v", err)
	}
	if hash, err := verifyPublishLease(ctx, harness.runRoot, supervisedTestRunID, foreign, harness.now); !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("verifying another run's proof must be refused, got %q %v", hash, err)
	}
	// Positive control: this run's own record is untouched by the planting.
	if _, err := renewPublishLease(ctx, harness.runRoot, supervisedTestRunID, own, harness.now); err != nil {
		t.Fatalf("this run's own proof must still renew, got %v", err)
	}
	if hash, err := verifyPublishLease(ctx, harness.runRoot, supervisedTestRunID, own, harness.now); err != nil || hash != own.LeaseHash {
		t.Fatalf("this run's own proof must still verify to its hash, got %q %v", hash, err)
	}

	harness.writeArtifact("artifact.txt", "hello from the supervised agent\n")
	summary, err := ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if err == nil {
		t.Fatalf("a proof of another run must not publish, got %+v", summary)
	}
	if hasStateName(summary.TaskStateCounts, "PASSED") {
		t.Fatalf("a refused promotion must leave no passed task: %+v", summary.TaskStateCounts)
	}
	// No completed promotion receipt: the seal is written only for an accepted
	// attempt and records the promotion decision that accepted it.
	store, err := chain.NewFileEventStore(runEventLogPath(harness.runRoot))
	if err != nil {
		t.Fatal(err)
	}
	if seal, _, sealErr := store.LoadTaskSeal(ctx, supervisedTestRunID, supervisedTestTaskID, 1); sealErr == nil && seal.PromotionDecision == "completed" {
		t.Fatalf("a refused promotion must not leave a completed seal, got %+v", seal)
	}
	// Lease evidence: no durable artifact outside the planted proof and the lease
	// store may bind the foreign lease hash.
	if path, bound := p83RunRootBindsHash(t, harness.runRoot, foreign.LeaseHash); bound {
		t.Fatalf("the foreign lease hash must not be bound by any artifact, found in %s", path)
	}
	after := p83LoadSupervisedLeaseProof(t, harness.runRoot)
	if after != foreign {
		t.Fatalf("a refused promotion must leave the planted proof untouched, got %+v want %+v", after, foreign)
	}
}
