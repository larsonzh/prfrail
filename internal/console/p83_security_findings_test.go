package console

// P83 SECURITY-FINDINGS-REMEDIATION: the port-level regressions of F2, F3 and F4
// that can be driven without a live pinned CLI.
//
// F2's end-to-end leg (a revocation appended to the on-disk ledger after assembly)
// needs the production dispatch path, which only publishes R where publication
// durability is proven; it is TestAgentRunnerDispatchRereadsRevocationAfterAssembly
// in this file and skips elsewhere.
//
// F3 and F4 are pinned at the port level on purpose: each test drives exactly the
// port the finding names and asserts the refusal plus the evidence the acceptance
// path binds, so a mutation of that one branch turns the test red without needing
// a full offline run.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

const (
	p83RunID  = "run-p83"
	p83TaskID = "task-p83"
	p83StepID = "step-p83"
)

// p83Clock is the fixed instant every port-level test evaluates at.
func p83Clock() time.Time { return time.Date(2026, 11, 20, 10, 0, 0, 0, time.UTC) }

// p83Grant builds an active operator grant covering p83RunID/p83TaskID/p83StepID.
func p83Grant(t *testing.T, policyHash string) chain.AuthorizationRecord {
	t.Helper()
	base := p83Clock()
	grant, err := chain.NewGrantAuthorizationRecord(chain.AuthorizationGrant{
		RecordID:        "p83-grant-1",
		Kind:            "grant",
		IssuedAt:        base.Add(-time.Minute).UTC().Format(timestampLayout),
		IssuedBy:        evidence.Actor{Type: "operator", ID: "p83-operator"},
		AuthorizationID: "p83-authorization-1",
		RunID:           p83RunID,
		RunManifestHash: evidence.Digest("", []byte("manifest")),
		PolicyHash:      policyHash,
		Scope: chain.AuthorizationScope{
			TaskIDs: []string{p83TaskID},
			StepIDs: []string{p83StepID},
			Budget: chain.AuthorizationBudgetLimit{
				ModelCalls: 1, Tokens: 1000, WallClockMs: 600000, Attempts: 1, Currency: "USD",
			},
		},
		ExpiresAt: base.Add(time.Hour).UTC().Format(timestampLayout),
		Evidence:  []string{evidence.Digest("", []byte("grant-evidence"))},
	}, p83Clock, func() string { return "p83-grant-record-1" })
	if err != nil {
		t.Fatalf("build grant: %v", err)
	}
	return grant
}

// p83Runtime builds the minimal runtime state the review and publish ports need.
// The request record is a plain identity: these ports read only its task, step,
// attempt and request id, so building a full wire record would add no coverage.
func p83Runtime(t *testing.T, ledger chain.AuthorizationLedger) *agentRunnerRuntime {
	t.Helper()
	record := adapters.AgentRunnerRequestRecord{Request: adapters.AgentRunnerRequest{
		RequestID: "request-p83",
		RunID:     p83RunID,
		TaskID:    p83TaskID,
		StepID:    p83StepID,
		Attempt:   1,
	}}
	key := agentRunnerTaskStateKey(p83TaskID, 1)
	return &agentRunnerRuntime{
		runID:               p83RunID,
		requestRecord:       record,
		requestRecords:      map[string]adapters.AgentRunnerRequestRecord{key: record},
		authorizationLedger: ledger,
		clock:               p83Clock,
	}
}

// TestAgentRunnerReviewRefusesWithoutOperatorGrant is the F3 refusal: a review
// port with no covering, active operator grant must refuse instead of approving
// on a digest it computed itself.
// Mutation check: restoring the self-derived PolicyHash approve turns this red.
func TestAgentRunnerReviewRefusesWithoutOperatorGrant(t *testing.T) {
	runtime := p83Runtime(t, chain.AuthorizationLedger{})
	request := chain.ReviewRequest{
		RunID:                 p83RunID,
		TaskID:                p83TaskID,
		Attempt:               1,
		CandidateSnapshotHash: evidence.Digest("", []byte("candidate")),
		EvidenceRootHash:      evidence.Digest("", []byte("evidence-root")),
	}
	decision, err := runtime.Review(context.Background(), request)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("a review without an operator grant must be refused with ErrInvalidConfig, got %v", err)
	}
	if decision.Outcome != "" || decision.PolicyHash != "" {
		t.Fatalf("a refused review must produce no approve decision, got %+v", decision)
	}
}

// TestAgentRunnerReviewBindsOperatorPolicyAndEvidenceRoot is the F3 acceptance
// half: the operator's policy hash is adopted verbatim and the decision carries
// the candidate/evidence-root binding the chain independently checks.
func TestAgentRunnerReviewBindsOperatorPolicyAndEvidenceRoot(t *testing.T) {
	// T2 (finding F3): the operator's policy hash must name the run's frozen policy
	// artifact, so the fixture freezes that policy first and grants its digest.
	runtime := p83Runtime(t, chain.AuthorizationLedger{})
	_, policyHash := p85FreezeAgentRunnerPolicy(t, runtime, p83TaskID, p83StepID, 1)
	runtime.authorizationLedger = chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{p83Grant(t, policyHash)}}
	request := chain.ReviewRequest{
		RunID:                 p83RunID,
		TaskID:                p83TaskID,
		Attempt:               1,
		CandidateSnapshotHash: evidence.Digest("", []byte("candidate")),
		EvidenceRootHash:      evidence.Digest("", []byte("evidence-root")),
	}
	decision, err := runtime.Review(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "approve" || decision.PolicyHash != policyHash {
		t.Fatalf("the review must adopt the operator policy hash %s, got %+v", policyHash, decision)
	}
	binding := chain.ReviewEvidenceBinding(request.CandidateSnapshotHash, request.EvidenceRootHash)
	found := false
	for _, item := range decision.Evidence {
		if item == binding {
			found = true
		}
	}
	if !found {
		t.Fatalf("the approve decision must carry the evidence-root binding, got %v", decision.Evidence)
	}
}

// p83BindReviewedPolicy seeds the state the F3 publish-time verification reads:
// the policy hash the review decision bound plus the operator grant that still
// binds it. The F4 tests drive Publish directly, without running a review, so
// they state that precondition explicitly instead of implying it.
func p83BindReviewedPolicy(t *testing.T, runtime *agentRunnerRuntime, taskID string, attempt int) string {
	t.Helper()
	stepID, err := runtime.stepIDForTask(taskID, attempt)
	if err != nil {
		t.Fatalf("resolve the step of task %q attempt %d: %v", taskID, attempt, err)
	}
	// T2 (finding F3): the reviewed policy hash must name a frozen artifact that
	// really exists on disk, so the fixture freezes one and grants exactly its
	// digest instead of an arbitrary digest the publish check could not verify.
	_, policyHash := p85FreezeAgentRunnerPolicy(t, runtime, taskID, stepID, attempt)
	runtime.authorizationLedger = chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{p83Grant(t, policyHash)}}
	if runtime.reviewPolicyHashes == nil {
		runtime.reviewPolicyHashes = map[string]string{}
	}
	runtime.reviewPolicyHashes[agentRunnerTaskStateKey(taskID, attempt)] = policyHash
	return policyHash
}

// p83PublishRuntime builds a runtime whose Publish reaches the lease proof: an
// empty workspace matches an empty accepted entry set, and both pid reports exist
// and list no process.
func p83PublishRuntime(t *testing.T) (*agentRunnerRuntime, string) {
	t.Helper()
	runRoot := t.TempDir()
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(runRoot)); err != nil {
		t.Fatal(err)
	}
	store, err := adapters.NewAgentRunnerReplayStoreForRun(runRoot, p83RunID)
	if err != nil {
		t.Fatal(err)
	}
	evidenceDir, err := store.RunEvidenceDir("request-p83")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{agentRunnerParentPIDReportFileName, agentRunnerChildPIDReportFileName} {
		if err := os.WriteFile(filepath.Join(evidenceDir, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	workspace := filepath.Join(runRoot, "workspaces", "p83")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	snapshotStore, err := snapshot.NewStore(filepath.Join(runRoot, "snapshot-store"), 0)
	if err != nil {
		t.Fatal(err)
	}
	capturer, err := adapters.NewAgentRunnerSnapshotCapturer(snapshotStore, p83RunID, p83TaskID, 1)
	if err != nil {
		t.Fatal(err)
	}
	key := agentRunnerTaskStateKey(p83TaskID, 1)
	base := p83Runtime(t, chain.AuthorizationLedger{})
	base.runRoot = runRoot
	base.store = store
	base.workspaceRoots = map[string]string{key: workspace}
	base.capturers = map[string]*adapters.AgentRunnerSnapshotCapturer{key: capturer}
	base.acceptedEntries = map[string][]snapshot.Entry{key: nil}
	p83BindReviewedPolicy(t, base, p83TaskID, 1)
	return base, runRoot
}

// TestAgentRunnerPublishRefusesWithoutCurrentLease is the F4 regression: the
// promotion completes only while this run still holds its lease record.
// Mutation check: replacing the lease re-read with a digest derived only from
// the run ID makes the second publish complete and turns this red.
func TestAgentRunnerPublishRefusesWithoutCurrentLease(t *testing.T) {
	runtime, runRoot := p83PublishRuntime(t)
	proof, err := claimPublishLease(context.Background(), runRoot, p83RunID)
	if err != nil {
		t.Fatalf("claim the publish lease: %v", err)
	}
	runtime.publishLease = proof
	// F3: the promotion also binds the candidate evidence root the review froze,
	// so the fixture states that root instead of an unrelated digest.
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
	decision, err := runtime.Publish(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "completed" || len(decision.LeaseEvidence) != 1 || decision.LeaseEvidence[0] != proof.LeaseHash {
		t.Fatalf("the promotion must bind the claimed lease record, got %+v", decision)
	}

	// Remove the lease: the run no longer proves exclusive ownership.
	if err := os.RemoveAll(publishLeaseRoot(runRoot)); err != nil {
		t.Fatal(err)
	}
	decision, err = runtime.Publish(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if decision.Outcome != "failed" {
		t.Fatalf("a promotion without a current lease must fail, got %+v", decision)
	}
	if len(decision.LeaseEvidence) != 0 {
		t.Fatalf("a failed promotion must bind no lease evidence, got %v", decision.LeaseEvidence)
	}
	if len(decision.ErrorEvidence) == 0 {
		t.Fatal("a failed promotion must carry error evidence")
	}
}

// TestSupervisedReturnRefusesDriftedPublishLease is the F4 re-entry half: the
// return runs in a new process, so it re-proves the lease the window claimed. A
// lease that disappeared between the confirmation and the return refuses the
// promotion instead of completing.
func TestSupervisedReturnRefusesDriftedPublishLease(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("first invocation must park, got %v", err)
	}
	harness.respond(supervisedConfirmationResponse)
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("second invocation must open the window, got %v", err)
	}
	harness.advance(time.Minute)
	harness.writeArtifact("artifact.txt", "hello from the supervised agent\n")
	if err := os.RemoveAll(publishLeaseRoot(harness.runRoot)); err != nil {
		t.Fatal(err)
	}
	summary, err := ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if err == nil {
		t.Fatalf("a return without a current lease must not complete, got %+v", summary)
	}
	if hasStateName(summary.TaskStateCounts, "PASSED") {
		t.Fatalf("a refused promotion must leave no passed task: %+v", summary.TaskStateCounts)
	}
}

// TestAgentRunnerDispatchRereadsRevocationAfterAssembly is the F2 end-to-end leg:
// the assembly holds an active grant in memory, the operator revokes it on disk
// after assembly, and the dispatch must observe the revocation before launching.
func TestAgentRunnerDispatchRereadsRevocationAfterAssembly(t *testing.T) {
	if durability := agentRunnerPublicationDurability(t); durability != adapters.PublishDurabilityProven {
		t.Skipf("SKIP: publication durability is %s on %s, so the production dispatch refuses before R and the reread cannot be observed. The reread runs on the Ubuntu CI leg.", durability, runtime.GOOS)
	}
	stub := buildAssemblyStub(t)
	fixture := newAssemblyFixture(t, assemblyKnobs{
		declaredExecutable: stub,
		runArgs:            []string{"-mode", "mutate-workspace", "-sleep", "100ms"},
	})
	grant := fixture.options.AuthorizationLedger.Grants[0]
	ledgerPath := filepath.Join(fixture.runRoot, "authorization-ledger.jsonl")
	if err := appendAuthorizationRecord(ledgerPath, grant); err != nil {
		t.Fatal(err)
	}
	revocation, err := chain.NewRevocationAuthorizationRecord(chain.AuthorizationRevocation{
		RecordID:          "p83-revocation-1",
		Kind:              "revocation",
		IssuedAt:          time.Now().UTC().Format(timestampLayout),
		IssuedBy:          evidence.Actor{Type: "operator", ID: "t027-declared-operator"},
		AuthorizationID:   grant.Grant.AuthorizationID,
		AuthorizationHash: grant.RecordHash,
		ReasonCode:        "operator-revoked",
		StopDisposition:   "not-required",
	}, time.Now, func() string { return "p83-revocation-record-1" })
	if err != nil {
		t.Fatalf("build revocation: %v", err)
	}
	if err := appendAuthorizationRecord(ledgerPath, revocation); err != nil {
		t.Fatal(err)
	}
	// The assembly keeps the in-memory grant: only the reread can see the revocation.
	fixture.options.AuthorizationLedgerPath = ledgerPath
	summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	// F2: the revocation must be observed before any launch. Two fail-closed gates
	// can see it - the admission ledger read (before R) and the dispatch
	// reconfirmation (after R) - and on a host where the entry gate reads the same
	// ledger the earlier one fires first. The dispatch reread itself is pinned by
	// TestP83F2ReconfirmationIsTheLastObservationBeforeLaunch; here the assertion is
	// that whichever gate fires, nothing is launched.
	reconfirmation := errors.Is(err, adapters.ErrAgentRunnerLaunchReconfirmation)
	admission := err != nil && strings.Contains(err.Error(), "authorization grant is revoked")
	if !reconfirmation && !admission {
		t.Fatalf("the run must refuse the launch after the revocation, got %v (summary %+v)", err, summary)
	}
	// The request record R, when it was published, stays published: a refusal never
	// deletes or retries it.
	requests := agentRunnerDirEntries(t, filepath.Join(fixture.runRoot, "agent-runner-replay", "requests"))
	if reconfirmation && len(requests) == 0 {
		t.Fatal("the published request record R must survive the refused launch")
	}
	launches := agentRunnerDirEntries(t, filepath.Join(fixture.runRoot, "agent-runner-replay", "launches"))
	if len(launches) != 0 {
		t.Fatalf("a refused launch must publish no launch receipt, found %v", launches)
	}
	if _, statErr := os.Stat(filepath.Join(fixture.runRoot, managedProcessIdentityFileName)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a refused launch must leave no process identity, stat err=%v", statErr)
	}
}
