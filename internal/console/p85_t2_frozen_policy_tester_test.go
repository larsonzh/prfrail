package console

// P85 T2 (finding F3): the counterexamples and the mutation targets of the
// durable FROZEN policy artifact.
//
// The design (tmp/P85-SECURITY-PROOF-SUBSYSTEM/01-design.md §2.4-§2.6) names
// three counterexample classes and two mutations:
//
//	missing          - the artifact does not exist when Review/Publish reread it;
//	unreadable       - the artifact exists but cannot be decoded (strict JSON);
//	well-formed but mismatched
//	                 - the artifact is internally consistent yet names another
//	                   policy than the operator grant, and the same class covers an
//	                   artifact replaced between Review and Publish.
//
// Mutation contracts:
//   - dropping the artifact comparison from operatorApprovedPolicyHashForAttempt
//     (runtime_agent.go) turns TestP85T2ReviewRefusesPolicyArtifactThatIsNotTheGrantedOne
//     green;
//   - dropping the disk reread from verifyReviewPolicyBinding (runtime_agent.go)
//     turns TestP85T2PublishRefusesArtifactReplacedAfterReview green.
//
// Every refusal assertion below states three facts, never just "some error came
// back": no approve decision, no accepted snapshot / no lease evidence, and the
// evidence that was on disk is still byte-identical.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85PolicyBody builds the policy body of the shape preview.go produces, so a
// fixture can state exactly the policy a run was previewed under.
func p85PolicyBody(chainID string, taskIDs ...string) previewPolicy {
	tasks := make([]previewPolicyTask, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		tasks = append(tasks, previewPolicyTask{
			TaskID:              taskID,
			Review:              "required",
			DocumentationPolicy: "required",
		})
	}
	return previewPolicy{
		ChainID:            chainID,
		ChainProfile:       "standard",
		ChainDocumentation: "required",
		Tasks:              tasks,
	}
}

// p85AttachFrozenPolicyStore binds the run root and the replay store the T2
// ports read the artifact through. Each test owns its run root, so the fixture
// never shares state with another test.
func p85AttachFrozenPolicyStore(t *testing.T, runtime *agentRunnerRuntime) {
	t.Helper()
	if runtime.runRoot == "" {
		runtime.runRoot = t.TempDir()
	}
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(runtime.runRoot)); err != nil {
		t.Fatal(err)
	}
	if runtime.store == nil {
		store, err := adapters.NewAgentRunnerReplayStoreForRun(runtime.runRoot, runtime.runID)
		if err != nil {
			t.Fatal(err)
		}
		runtime.store = store
	}
}

// p85FreezeAgentRunnerPolicy persists one frozen policy artifact for the runtime
// and returns the policy body plus the hash an operator grant must carry to be
// comparable with it.
func p85FreezeAgentRunnerPolicy(t *testing.T, runtime *agentRunnerRuntime, taskID, stepID string, attempt int) (previewPolicy, string) {
	t.Helper()
	p85AttachFrozenPolicyStore(t, runtime)
	body := p85PolicyBody(runtime.runID+"-chain", taskID)
	// F01 (P85 fix round 3): the freeze also proves the managed identity cannot
	// write the policy directory. That proof cannot succeed on this platform (the
	// assembly and the managed agent share one OS identity), so the fixture injects
	// the same test seam the production assembly fixture uses; the production
	// fail-closed path is pinned by the "without proof" counterexample below.
	//
	// P85 review r1 (new Medium): every load/Review/Publish of this artifact
	// re-observes the write domain and requires the anchor's value to be THAT
	// observation. The fixture therefore installs the matching stand-in, so the
	// cases below keep measuring their own property instead of the shared-identity
	// host; the real primitive is exercised by the M1/M2/M10 counterexamples.
	p85StubManagedWriteProofObservation(t, assemblyManagedWriteProof)
	persisted, err := persistAgentRunnerFrozenPolicyProven(runtime.store, runtime.runRoot, runtime.runID, taskID, stepID, attempt, body, assemblyManagedWriteProof)
	if err != nil {
		t.Fatalf("freeze the policy of task %q step %q attempt %d: %v", taskID, stepID, attempt, err)
	}
	return persisted.PolicyBody, persisted.PolicyHash
}

// p85FrozenReviewRuntime builds the minimal review/publish runtime with a frozen
// policy artifact matching an active operator grant. It returns the runtime and
// the policy hash both sides agree on.
func p85FrozenReviewRuntime(t *testing.T) (*agentRunnerRuntime, string) {
	t.Helper()
	runtime := p83Runtime(t, chain.AuthorizationLedger{})
	_, policyHash := p85FreezeAgentRunnerPolicy(t, runtime, p83TaskID, p83StepID, 1)
	runtime.authorizationLedger = chain.AuthorizationLedger{
		Grants: []chain.AuthorizationRecord{p83Fix5Grant(t, func(grant *chain.AuthorizationGrant) {
			grant.PolicyHash = policyHash
		})},
	}
	return runtime, policyHash
}

// p85RuntimeOverFrozenFixture rebuilds a runtime that reads the artifact of an
// already-frozen fixture, so two ledgers can be evaluated against one artifact.
func p85RuntimeOverFrozenFixture(t *testing.T, fixture *agentRunnerRuntime, ledger chain.AuthorizationLedger) *agentRunnerRuntime {
	t.Helper()
	runtime := p83Runtime(t, ledger)
	runtime.runRoot = fixture.runRoot
	runtime.store = fixture.store
	return runtime
}

func p85ReviewRequest() chain.ReviewRequest {
	return chain.ReviewRequest{
		RunID:                 p83RunID,
		TaskID:                p83TaskID,
		Attempt:               1,
		CandidateSnapshotHash: evidence.Digest("", []byte("candidate")),
		EvidenceRootHash:      evidence.Digest("", []byte("evidence-root")),
	}
}

// p85ReadFrozenPolicyBytes reads the frozen artifact as the bytes on disk, so a
// test can prove a refusal left the evidence untouched.
func p85ReadFrozenPolicyBytes(t *testing.T, runtime *agentRunnerRuntime) []byte {
	t.Helper()
	data, err := os.ReadFile(agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1))
	if err != nil {
		t.Fatalf("read the frozen policy artifact: %v", err)
	}
	return data
}

// TestP85T2ReviewRefusesMissingPolicyArtifact is the "missing" counterexample:
// with the artifact gone, a covering active grant with the right hash must not
// approve, because the granted hash no longer names a policy that exists on disk.
func TestP85T2ReviewRefusesMissingPolicyArtifact(t *testing.T) {
	runtime, policyHash := p85FrozenReviewRuntime(t)
	request := p85ReviewRequest()

	// Positive control: with the artifact present the same runtime approves.
	control, err := runtime.Review(context.Background(), request)
	if err != nil {
		t.Fatalf("a frozen policy must let the covering grant approve, got %v", err)
	}
	if control.Outcome != "approve" || control.PolicyHash != policyHash {
		t.Fatalf("the review must bind the frozen policy hash %s, got %+v", policyHash, control)
	}

	path := agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1)
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove the frozen policy artifact: %v", err)
	}
	decision, err := runtime.Review(context.Background(), request)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("a missing frozen policy must refuse the review, got %+v, %v", decision, err)
	}
	if decision.Outcome != "" || decision.PolicyHash != "" || len(decision.Evidence) != 0 {
		t.Fatalf("a refused review must produce no decision, got %+v", decision)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the refusal must not recreate the artifact (stat err=%v)", statErr)
	}
}

// TestP85T2ReviewRefusesUnreadablePolicyArtifact is the "unreadable"
// counterexample: an artifact that cannot be decoded strictly is a refusal, never
// a zero value the review could read as approval.
func TestP85T2ReviewRefusesUnreadablePolicyArtifact(t *testing.T) {
	cases := map[string]string{
		"malformed json":        "{\"runId\": \"run-p83\", \"policyBody\": ",
		"wrong typed body":      "{\"runId\":\"run-p83\",\"taskId\":\"task-p83\",\"stepId\":\"step-p83\",\"attempt\":1,\"policyBody\":{\"chainId\":7},\"policyHash\":\"sha256:" + "00" + "\"}",
		"unknown strict fields": "{\"runId\":\"run-p83\",\"unexpected\":true}",
	}
	for name, payload := range cases {
		t.Run(name, func(t *testing.T) {
			runtime, _ := p85FrozenReviewRuntime(t)
			path := agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1)
			if err := os.WriteFile(path, []byte(payload), 0o600); err != nil {
				t.Fatalf("corrupt the frozen policy artifact: %v", err)
			}
			// Evidence integrity: the corruption is what the refusal reads, and the
			// guard must not rewrite it.
			before, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			decision, err := runtime.Review(context.Background(), p85ReviewRequest())
			if !errors.Is(err, ErrInvalidConfig) {
				t.Fatalf("an unreadable frozen policy must refuse the review, got %+v, %v", decision, err)
			}
			if decision.Outcome != "" || decision.PolicyHash != "" {
				t.Fatalf("a refused review must produce no decision, got %+v", decision)
			}
			after, readErr := os.ReadFile(path)
			if readErr != nil {
				t.Fatalf("the refusal must leave the artifact in place: %v", readErr)
			}
			if string(after) != string(before) {
				t.Fatalf("the refusal must not rewrite the evidence: got %q want %q", after, before)
			}
		})
	}
}

// TestP85T2ReviewRefusesPolicyArtifactThatIsNotTheGrantedOne is the
// "well-formed but mismatched" counterexample at review time: the artifact is
// internally consistent (its own body re-derives to its own hash) but it is not
// the policy the operator granted.
//
// Mutation check: deleting the verifyFrozenPolicyArtifact call from
// operatorApprovedPolicyHashForAttempt turns the refusal leg green.
func TestP85T2ReviewRefusesPolicyArtifactThatIsNotTheGrantedOne(t *testing.T) {
	runtime, frozenHash := p85FrozenReviewRuntime(t)
	otherPolicy := evidence.Digest("", []byte("another-operator-policy"))
	if otherPolicy == frozenHash {
		t.Fatal("the fixture needs two different policy digests")
	}
	runtime.authorizationLedger = chain.AuthorizationLedger{
		Grants: []chain.AuthorizationRecord{p83Fix5Grant(t, func(grant *chain.AuthorizationGrant) {
			grant.PolicyHash = otherPolicy
		})},
	}
	before := p85ReadFrozenPolicyBytes(t, runtime)
	decision, err := runtime.Review(context.Background(), p85ReviewRequest())
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("a grant naming another policy must refuse the review, got %+v, %v", decision, err)
	}
	if decision.Outcome != "" || decision.PolicyHash != "" {
		t.Fatalf("a refused review must produce no decision, got %+v", decision)
	}
	if after := p85ReadFrozenPolicyBytes(t, runtime); string(after) != string(before) {
		t.Fatalf("the refusal must not rewrite the evidence: got %q want %q", after, before)
	}
	if _, bound := runtime.reviewPolicyHashes[agentRunnerTaskStateKey(p83TaskID, 1)]; bound {
		t.Fatal("a refused review must not bind a reviewed policy hash")
	}
}

// TestP85T2PublishRefusesArtifactReplacedAfterReview is the publish-time half of
// "well-formed but mismatched": the review binds the frozen policy and the
// artifact is then replaced with a different, internally consistent artifact.
// Publish must reread the artifact from disk and refuse with no accepted
// snapshot and no lease evidence.
//
// Mutation check: replacing the disk reread in verifyReviewPolicyBinding with the
// in-memory reviewed hash turns both legs green.
func TestP85T2PublishRefusesArtifactReplacedAfterReview(t *testing.T) {
	promotion := func(t *testing.T, runtime *agentRunnerRuntime) chain.PromotionRequest {
		t.Helper()
		evidenceRoot, err := runtime.evidenceRootHash(p83TaskID, 1)
		if err != nil {
			t.Fatal(err)
		}
		return chain.PromotionRequest{
			RunID:                 p83RunID,
			TaskID:                p83TaskID,
			Attempt:               1,
			CandidateSnapshotHash: evidence.Digest("", []byte("candidate")),
			EvidenceRootHash:      evidenceRoot,
			ReviewReceiptHash:     evidence.Digest("", []byte("review-receipt")),
		}
	}
	// setup drives the production publish fixture to a state where every other
	// publish obligation holds (lease, evidence root, workspace, pid reports), so
	// the ONLY fact left to vary below is the frozen policy artifact itself.
	setup := func(t *testing.T) (*agentRunnerRuntime, chain.PromotionRequest, []byte) {
		t.Helper()
		runtime, runRoot := p83PublishRuntime(t)
		proof, err := claimPublishLease(context.Background(), runRoot, p83RunID)
		if err != nil {
			t.Fatal(err)
		}
		runtime.publishLease = proof
		if _, err := runtime.Review(context.Background(), p85ReviewRequest()); err != nil {
			t.Fatalf("the review must approve before the publish legs: %v", err)
		}
		original := p85ReadFrozenPolicyBytes(t, runtime)
		return runtime, promotion(t, runtime), original
	}
	// assertBlocked pins the fail-closed promotion shape, then proves the artifact
	// is the sole cause by restoring it and requiring the same runtime to complete.
	assertBlocked := func(t *testing.T, runtime *agentRunnerRuntime, request chain.PromotionRequest, path string, original []byte) {
		t.Helper()
		decision, err := runtime.Publish(context.Background(), request)
		if err != nil {
			t.Fatalf("a blocked promotion is a decision, not an error: %v", err)
		}
		if decision.Outcome != "failed" {
			t.Fatalf("a tampered artifact must block the promotion, got %+v", decision)
		}
		if decision.AcceptedSnapshotHash != "" || len(decision.LeaseEvidence) != 0 || len(decision.WriterStopEvidence) != 0 {
			t.Fatalf("a blocked promotion must accept no snapshot and bind no lease/writer evidence, got %+v", decision)
		}
		if len(decision.ErrorEvidence) == 0 {
			t.Fatal("a blocked promotion must carry error evidence")
		}
		if after, readErr := os.ReadFile(path); readErr == nil && string(after) == string(original) {
			t.Fatal("the fixture never changed the artifact")
		}
		if err := os.WriteFile(path, original, 0o600); err != nil {
			t.Fatalf("restore the frozen policy artifact: %v", err)
		}
		restored, err := runtime.Publish(context.Background(), request)
		if err != nil {
			t.Fatalf("the restored artifact must be publishable, got %v", err)
		}
		if restored.Outcome != "completed" {
			t.Fatalf("restoring the artifact must complete the promotion, got %+v", restored)
		}
	}

	t.Run("the artifact is replaced with another policy", func(t *testing.T) {
		runtime, request, original := setup(t)
		divergentBody := p85PolicyBody("chain-p85-divergent", p83TaskID)
		divergentHash, err := agentRunnerFrozenPolicyHash(divergentBody)
		if err != nil {
			t.Fatal(err)
		}
		payload, err := agentRunnerFrozenPolicyPayload(agentRunnerFrozenPolicy{
			RunID:      p83RunID,
			TaskID:     p83TaskID,
			StepID:     p83StepID,
			Attempt:    1,
			PolicyBody: divergentBody,
			PolicyHash: divergentHash,
		})
		if err != nil {
			t.Fatal(err)
		}
		// The replacement is internally consistent: its own body re-derives to its
		// own hash, so only a disk reread compared against the reviewed value can
		// detect it.
		path := agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1)
		if err := os.WriteFile(path, payload, 0o600); err != nil {
			t.Fatalf("replace the frozen policy artifact: %v", err)
		}
		assertBlocked(t, runtime, request, path, original)
	})

	t.Run("the artifact is deleted after the review", func(t *testing.T) {
		runtime, request, original := setup(t)
		path := agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1)
		if err := os.Remove(path); err != nil {
			t.Fatalf("remove the frozen policy artifact: %v", err)
		}
		assertBlocked(t, runtime, request, path, original)
	})
}

// TestP85T2AssemblyRefusesWithoutFrozenPolicy is the assembly-level "missing"
// counterexample: an AgentRunner run without the previewed policy body must fail
// closed before anything is dispatched or spawned.
func TestP85T2AssemblyRefusesWithoutFrozenPolicy(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{})
	fixture.options.PolicySnapshot = nil
	summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("an assembly without a frozen policy must be refused, got %+v, %v", summary, err)
	}
	if summary.ChainState != "" || summary.RunID != "" {
		t.Fatalf("a refused assembly must return no summary, got %+v", summary)
	}
	if entries := agentRunnerDirEntries(t, agentRunnerPolicyDirectory(fixture.runRoot)); len(entries) != 0 {
		t.Fatalf("a refused assembly must freeze nothing, found %v", entries)
	}
	if _, statErr := os.Stat(filepath.Join(fixture.runRoot, managedProcessIdentityFileName)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a refused assembly must leave no managed-process identity mirror (stat err=%v)", statErr)
	}
	data, readErr := os.ReadFile(runEventLogPath(fixture.runRoot))
	if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
		t.Fatal(readErr)
	}
	if strings.Contains(string(data), "agent-dispatched") || strings.Contains(string(data), "TERMINAL_PENDING") {
		t.Fatalf("a policy-freeze refusal must record no dispatch transition: %s", data)
	}
}

// p85DivergentFrozenPolicyPayload builds a replacement artifact that is internally
// consistent - its own policy body re-derives to its own policy hash - so only a
// disk reread compared against the reviewed value can tell it apart from the
// artifact the review bound.
func p85DivergentFrozenPolicyPayload(t *testing.T) []byte {
	t.Helper()
	body := p85PolicyBody("chain-p85-divergent", p83TaskID)
	policyHash, err := agentRunnerFrozenPolicyHash(body)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := agentRunnerFrozenPolicyPayload(agentRunnerFrozenPolicy{
		RunID:      p83RunID,
		TaskID:     p83TaskID,
		StepID:     p83StepID,
		Attempt:    1,
		PolicyBody: body,
		PolicyHash: policyHash,
	})
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

// TestP85T2PublishRefusesArtifactReplacedAtPromotionPoint is the F01 (P85 fix
// round 1) counterexample: the artifact is replaced WHILE Publish runs - after the
// entry-time reread has already passed and immediately before the promotion is
// linearized. The entry-time verification cannot see that replacement, so only the
// second reread at the promotion point blocks the promotion.
//
// The replacement is injected through the runtime clock: Publish reads it once
// outside the policy-binding verification, between the entry-time policy binding
// and the promotion-point binding. The window is therefore stated deterministically
// instead of raced.
//
// Mutation check: deleting the promotion-point verifyReviewPolicyBinding call in
// runtime_agent.go makes Publish complete and turns this test red.
func TestP85T2PublishRefusesArtifactReplacedAtPromotionPoint(t *testing.T) {
	runtime, runRoot := p83PublishRuntime(t)
	proof, err := claimPublishLease(context.Background(), runRoot, p83RunID)
	if err != nil {
		t.Fatal(err)
	}
	runtime.publishLease = proof
	if _, err := runtime.Review(context.Background(), p85ReviewRequest()); err != nil {
		t.Fatalf("the review must approve before the publish legs: %v", err)
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
	path := agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1)
	original := p85ReadFrozenPolicyBytes(t, runtime)
	replacement := p85DivergentFrozenPolicyPayload(t)

	// The injection point is the first clock read that does NOT happen inside the
	// policy-binding verification itself. The entry-time binding reads the clock too
	// (grant expiry, before its own disk reread), so counting calls would land in the
	// wrong place; the stack is consulted instead, which keeps the replacement in the
	// F01 window: after the entry-time binding returned and before the
	// promotion-point binding.
	outsidePolicyBinding := func() bool {
		buffer := make([]byte, 8192)
		stack := string(buffer[:goruntime.Stack(buffer, false)])
		return !strings.Contains(stack, "verifyReviewPolicyBinding") &&
			!strings.Contains(stack, "operatorApprovedPolicyHashForAttempt")
	}
	baseClock := runtime.clock
	armed := true
	injections := 0
	runtime.clock = func() time.Time {
		if armed && outsidePolicyBinding() {
			armed = false
			injections++
			if writeErr := os.WriteFile(path, replacement, 0o600); writeErr != nil {
				t.Errorf("replace the frozen policy artifact at the promotion point: %v", writeErr)
			}
		}
		return baseClock()
	}

	decision, err := runtime.Publish(context.Background(), request)
	if err != nil {
		t.Fatalf("a blocked promotion is a decision, not an error: %v", err)
	}
	if injections != 1 {
		t.Fatalf("the fixture never reached the window between the entry check and the promotion point (%d injections)", injections)
	}
	if decision.Outcome != "failed" {
		t.Fatalf("an artifact replaced at the promotion point must block the promotion, got %+v", decision)
	}
	if decision.AcceptedSnapshotHash != "" || len(decision.LeaseEvidence) != 0 || len(decision.WriterStopEvidence) != 0 {
		t.Fatalf("a blocked promotion must accept no snapshot and bind no lease/writer evidence, got %+v", decision)
	}
	if len(decision.ErrorEvidence) == 0 {
		t.Fatal("a blocked promotion must carry error evidence")
	}
	if after, readErr := os.ReadFile(path); readErr != nil || string(after) != string(replacement) {
		t.Fatalf("the refusal must read the replaced artifact and leave it untouched (read err=%v)", readErr)
	}
	// The replacement is the sole cause: the same runtime completes once the
	// reviewed artifact is back on disk.
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatalf("restore the frozen policy artifact: %v", err)
	}
	restored, err := runtime.Publish(context.Background(), request)
	if err != nil {
		t.Fatalf("the restored artifact must be publishable, got %v", err)
	}
	if restored.Outcome != "completed" {
		t.Fatalf("restoring the artifact must complete the promotion, got %+v", restored)
	}
}

// TestP85T2AssemblyRefusesMismatchedPolicySnapshot is the F05 (P85 fix round 1)
// counterexample of the pre-dispatch freeze: the snapshot is well formed but it is
// not the policy of THIS definition - another chain, another task count, another
// task at a position, or the same tasks in a different order - and every shape is
// refused before anything is frozen or dispatched.
//
// Mutation check: restoring the loose "a chain id is present and the tasks are
// covered" validation in agent_runner_policy.go turns every subtest green.
func TestP85T2AssemblyRefusesMismatchedPolicySnapshot(t *testing.T) {
	assertRefused := func(t *testing.T, runRoot string, summary RunSummary, err error, want string, dispatchCount func() int) {
		t.Helper()
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("a mismatched policy snapshot must be refused at assembly, got %+v, %v", summary, err)
		}
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("the refusal must name the mismatch %q, got %v", want, err)
		}
		if summary.ChainState != "" || summary.RunID != "" {
			t.Fatalf("a refused assembly must return no summary, got %+v", summary)
		}
		if entries := agentRunnerDirEntries(t, agentRunnerPolicyDirectory(runRoot)); len(entries) != 0 {
			t.Fatalf("a refused assembly must freeze nothing, found %v", entries)
		}
		if dispatchCount != nil {
			if count := dispatchCount(); count != 0 {
				t.Fatalf("a refused assembly must dispatch nothing, got %d", count)
			}
		}
		if _, statErr := os.Stat(filepath.Join(runRoot, managedProcessIdentityFileName)); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("a refused assembly must leave no managed-process identity mirror (stat err=%v)", statErr)
		}
		data, readErr := os.ReadFile(runEventLogPath(runRoot))
		if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
			t.Fatal(readErr)
		}
		if strings.Contains(string(data), "agent-dispatched") || strings.Contains(string(data), "TERMINAL_PENDING") {
			t.Fatalf("a policy-freeze refusal must record no dispatch transition: %s", data)
		}
	}

	t.Run("the snapshot names another chain", func(t *testing.T) {
		fixture := newAssemblyFixture(t, assemblyKnobs{})
		body := p85PolicyBody("chain-p85-other", assemblyTaskID)
		fixture.options.PolicySnapshot = &body
		summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
		assertRefused(t, fixture.runRoot, summary, err, "names chain", nil)
	})

	t.Run("the snapshot describes another task count", func(t *testing.T) {
		fixture := newAssemblyFixture(t, assemblyKnobs{})
		body := p85PolicyBody(assemblyDefinition().ID, assemblyTaskID, "task-p85-extra")
		fixture.options.PolicySnapshot = &body
		summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
		assertRefused(t, fixture.runRoot, summary, err, "but the assembly definition declares", nil)
	})

	t.Run("the snapshot names another task at a position", func(t *testing.T) {
		fixture := newAssemblyFixture(t, assemblyKnobs{})
		body := p85PolicyBody(assemblyDefinition().ID, "task-p85-other")
		fixture.options.PolicySnapshot = &body
		summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
		assertRefused(t, fixture.runRoot, summary, err, "at position 0", nil)
	})

	t.Run("the snapshot carries the same tasks in another order", func(t *testing.T) {
		options, stub, _ := newAssemblyDriverFixture(t)
		body := p85PolicyBody(multiTaskDefinitionID, multiTaskAgentTaskTwo, multiTaskAgentTaskOne)
		options.PolicySnapshot = &body
		summary, err := ExecuteAgentRunnerRun(context.Background(), options)
		assertRefused(t, options.RunRoot, summary, err, "at position 0", func() int { return len(stub.dispatched) })
	})
}

// p85AssertAssemblyRefusedWithoutWriteProof pins the fail-closed shape of a
// policy freeze that could not prove the managed identity is unable to write the
// policy directory: zero summary, zero frozen artifacts, no managed-process
// identity mirror and no dispatch transition.
func p85AssertAssemblyRefusedWithoutWriteProof(t *testing.T, runRoot string, summary RunSummary, err error) {
	t.Helper()
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("an assembly that cannot prove the policy write domain closed must be refused, got %+v, %v", summary, err)
	}
	if summary.ChainState != "" || summary.RunID != "" {
		t.Fatalf("a refused assembly must return no summary, got %+v", summary)
	}
	if entries := agentRunnerDirEntries(t, agentRunnerPolicyDirectory(runRoot)); len(entries) != 0 {
		t.Fatalf("a refused assembly must freeze no policy artifact, found %v", entries)
	}
	if _, statErr := os.Stat(filepath.Join(runRoot, managedProcessIdentityFileName)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a refused assembly must leave no managed-process identity mirror (stat err=%v)", statErr)
	}
	data, readErr := os.ReadFile(runEventLogPath(runRoot))
	if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
		t.Fatal(readErr)
	}
	if strings.Contains(string(data), "agent-dispatched") || strings.Contains(string(data), "TERMINAL_PENDING") {
		t.Fatalf("a policy-freeze refusal must record no dispatch transition: %s", data)
	}
}

// p85RewriteAnchorProof rewrites ONLY the managed-write proof field of the
// assembly anchor of the one task attempt the T2 fixtures use. The anchor stays
// strict-decodable and the artifact is not touched, so the mutation isolates the
// anchor as the single varying fact.
func p85RewriteAnchorProof(t *testing.T, runtime *agentRunnerRuntime, proof string) (original, rewritten []byte) {
	t.Helper()
	anchorPath := agentRunnerFrozenPolicyAnchorPath(runtime.runRoot, p83TaskID, p83StepID, 1)
	original, err := os.ReadFile(anchorPath)
	if err != nil {
		t.Fatalf("read the frozen policy anchor: %v", err)
	}
	var anchor agentRunnerFrozenPolicyAnchor
	if err := evidence.DecodeStrictJSON(original, &anchor); err != nil {
		t.Fatalf("decode the frozen policy anchor: %v", err)
	}
	anchor.ManagedWriteProof = proof
	rewritten, err = agentRunnerFrozenPolicyAnchorPayload(anchor)
	if err != nil {
		t.Fatalf("encode the rewritten frozen policy anchor: %v", err)
	}
	if err := os.WriteFile(anchorPath, rewritten, 0o600); err != nil {
		t.Fatalf("rewrite the frozen policy anchor: %v", err)
	}
	return original, rewritten
}

// TestP85T2AssemblyRefusesWithoutManagedWriteProof is the F01 (P85 fix round 3)
// fail-closed counterexample: without the injected test seam the real platform
// guard runs, and on this platform it can never prove the managed identity is
// unable to write the directory the assembly itself created. The assembly must
// therefore refuse before any dispatch and freeze nothing.
//
// Mutation check: returning a constant from agentRunnerProvePolicyWriteProtection
// instead of running the guard turns this refusal green.
func TestP85T2AssemblyRefusesWithoutManagedWriteProof(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{})
	fixture.options.managedWriteProof = ""
	summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil || !strings.Contains(err.Error(), "cannot be proven unwritable") {
		t.Fatalf("the refusal must name the unproven write domain, got %v", err)
	}
	p85AssertAssemblyRefusedWithoutWriteProof(t, fixture.runRoot, summary, err)
}

// TestP85T2AssemblyRefusesInvalidManagedWriteProof is the F01 counterexample of
// the test seam itself: an injected value that is not a valid evidence digest
// proves nothing, so the assembly refuses instead of adopting it as a proof.
//
// Mutation check: dropping the evidence.ValidHash check of the injected value in
// agentRunnerProvePolicyWriteProtection lets this assembly freeze and turns the
// refusal green.
func TestP85T2AssemblyRefusesInvalidManagedWriteProof(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{})
	fixture.options.managedWriteProof = "not-a-hash"
	summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil || !strings.Contains(err.Error(), "not a valid proof digest") {
		t.Fatalf("the refusal must name the invalid proof digest, got %v", err)
	}
	p85AssertAssemblyRefusedWithoutWriteProof(t, fixture.runRoot, summary, err)
}

// TestP85T2PublishRefusesAnchorRewrittenAlone is the F01 counterexample that
// isolates the ANCHOR: the frozen artifact is left byte-identical and the grant
// still binds the same policy hash, so the only mutated fact is the anchor.
//
// Two legs state the contract:
//  1. a reopened assembly (a resumed run adopting the existing artifact) refuses
//     on the anchor byte comparison - confirmAgentRunnerFrozenPolicyAnchor - even
//     when the rewritten proof is itself a valid digest, because the anchor is
//     never replaced;
//  2. Publish refuses when the rewritten anchor no longer carries a valid
//     managed-write proof, because the load path requires that proof to reach the
//     artifact.
//
// Mutation check: dropping the confirmAgentRunnerFrozenPolicyAnchor call from the
// reuse branch turns leg 1 green; dropping the ManagedWriteProof check from
// agentRunnerVerifyFrozenPolicyAnchor turns leg 2 green.
func TestP85T2PublishRefusesAnchorRewrittenAlone(t *testing.T) {
	t.Run("a rewritten anchor refuses the reopening assembly", func(t *testing.T) {
		runtime, _ := p85FrozenReviewRuntime(t)
		body := p85PolicyBody(runtime.runID+"-chain", p83TaskID)
		artifactBefore := p85ReadFrozenPolicyBytes(t, runtime)
		// A different but still valid proof digest: the anchor stays decodable and
		// semantically plausible, so only the byte comparison can catch the rewrite.
		otherProof := "sha256:" + strings.Repeat("2", 64)
		original, rewritten := p85RewriteAnchorProof(t, runtime, otherProof)
		if string(original) == string(rewritten) {
			t.Fatal("the fixture never changed the anchor")
		}
		_, err := persistAgentRunnerFrozenPolicyProven(runtime.store, runtime.runRoot, p83RunID, p83TaskID, p83StepID, 1, body, assemblyManagedWriteProof)
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("a reopened assembly must refuse a rewritten anchor, got %v", err)
		}
		if !strings.Contains(err.Error(), "anchor") || !strings.Contains(err.Error(), "different content") {
			t.Fatalf("the refusal must come from the anchor byte comparison, got %v", err)
		}
		if after := p85ReadFrozenPolicyBytes(t, runtime); string(after) != string(artifactBefore) {
			t.Fatalf("the artifact must stay byte-identical, got %q want %q", after, artifactBefore)
		}
		anchorPath := agentRunnerFrozenPolicyAnchorPath(runtime.runRoot, p83TaskID, p83StepID, 1)
		if after, readErr := os.ReadFile(anchorPath); readErr != nil || string(after) != string(rewritten) {
			t.Fatalf("the refusal must not repair the anchor (read err=%v)", readErr)
		}
	})

	t.Run("a rewritten anchor refuses publish", func(t *testing.T) {
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
		artifactBefore := p85ReadFrozenPolicyBytes(t, runtime)
		anchorPath := agentRunnerFrozenPolicyAnchorPath(runtime.runRoot, p83TaskID, p83StepID, 1)
		original, _ := p85RewriteAnchorProof(t, runtime, "not-a-hash")

		decision, err := runtime.Publish(context.Background(), request)
		if err != nil {
			t.Fatalf("a blocked promotion is a decision, not an error: %v", err)
		}
		if decision.Outcome != "failed" {
			t.Fatalf("an anchor without a valid write proof must block the promotion, got %+v", decision)
		}
		if decision.AcceptedSnapshotHash != "" || len(decision.LeaseEvidence) != 0 || len(decision.WriterStopEvidence) != 0 {
			t.Fatalf("a blocked promotion must accept no snapshot and bind no lease/writer evidence, got %+v", decision)
		}
		if len(decision.ErrorEvidence) == 0 {
			t.Fatal("a blocked promotion must carry error evidence")
		}
		if after := p85ReadFrozenPolicyBytes(t, runtime); string(after) != string(artifactBefore) {
			t.Fatalf("the artifact must stay byte-identical, got %q want %q", after, artifactBefore)
		}

		// The anchor is the sole cause: restoring it completes the same runtime.
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
	})
}
