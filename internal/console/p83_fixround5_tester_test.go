package console

// P83 FIX ROUND 5 counterexamples (items T1a, T1b, T2a, T2b, T3a, T3b, T5a, T9b).
//
// Every test below is written so that deleting the one guard it names turns it
// red; the mutation contract is stated in each test's doc comment.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

// p83Fix5Grant builds a grant identical to p83Grant but lets a test mutate it
// before it is sealed, so one fixture serves the supersession and expiry cases.
func p83Fix5Grant(t *testing.T, mutate func(*chain.AuthorizationGrant)) chain.AuthorizationRecord {
	t.Helper()
	base := p83Clock()
	grant := chain.AuthorizationGrant{
		RecordID:        "p83-grant-1",
		Kind:            "grant",
		IssuedAt:        base.Add(-time.Minute).UTC().Format(timestampLayout),
		IssuedBy:        evidence.Actor{Type: "operator", ID: "p83-operator"},
		AuthorizationID: "p83-authorization-1",
		RunID:           p83RunID,
		RunManifestHash: evidence.Digest("", []byte("manifest")),
		PolicyHash:      evidence.Digest("", []byte("p83-policy")),
		Scope: chain.AuthorizationScope{
			TaskIDs: []string{p83TaskID},
			StepIDs: []string{p83StepID},
			Budget: chain.AuthorizationBudgetLimit{
				ModelCalls: 1, Tokens: 1000, WallClockMs: 600000, Attempts: 1, Currency: "USD",
			},
		},
		ExpiresAt: base.Add(time.Hour).UTC().Format(timestampLayout),
		Evidence:  []string{evidence.Digest("", []byte("grant-evidence"))},
	}
	if mutate != nil {
		mutate(&grant)
	}
	record, err := chain.NewGrantAuthorizationRecord(grant, p83Clock, func() string { return grant.RecordID + "-record" })
	if err != nil {
		t.Fatalf("build grant %s: %v", grant.RecordID, err)
	}
	return record
}

// TestP83T1aProductionAssemblyWithoutLedgerPathRefuses is the T1a counterexample:
// a production assembly (not marked as an in-memory test assembly) with no
// authorization-ledger path must refuse at assembly time and leave the run root
// untouched, while the marked test assembly still reaches the later stages.
//
// Mutation check: deleting the `if !options.testInMemoryLedgers` block in
// ExecuteAgentRunnerRun lets the unflagged production assembly continue past the
// gate (in-memory ledger) and turns the refusal leg red.
func TestP83T1aProductionNeverAnswersFromTheInMemoryLedger(t *testing.T) {
	runRoot := t.TempDir()
	grant := p83Fix5Grant(t, nil)
	runtime := p83Runtime(t, chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{grant}})
	runtime.runRoot = runRoot
	// A production runtime derives the ledger from the durable run root; the derived
	// file is absent here, so the in-memory grant must not be consulted.
	runtime.authorizationLedgerPath = filepath.Join(runRoot, defaultApprovalsLedgerName)
	request := chain.ReviewRequest{
		RunID:                 p83RunID,
		TaskID:                p83TaskID,
		Attempt:               1,
		CandidateSnapshotHash: evidence.Digest("", []byte("candidate")),
		EvidenceRootHash:      evidence.Digest("", []byte("evidence-root")),
	}
	decision, err := runtime.Review(context.Background(), request)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("an absent derived ledger must refuse the review, got %+v, %v", decision, err)
	}
	if !strings.Contains(err.Error(), "no operator-approved policy grant covers this review") && !strings.Contains(err.Error(), "authorization ledger is unreadable") {
		t.Fatalf("the refusal must come from the trusted ledger source, got %v", err)
	}
	if decision.Outcome != "" || decision.PolicyHash != "" {
		t.Fatalf("a refused review must produce no decision, got %+v", decision)
	}
	if entries := agentRunnerDirEntries(t, runRoot); len(entries) != 0 {
		t.Fatalf("a refused authorization must write nothing to the run root, found %v", entries)
	}

	// Compatibility control: the marked in-memory assembly (no path) still answers
	// from its in-memory ledger, so the refusal above is attributable to the
	// production source and not to an empty ledger.
	control := p83Runtime(t, chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{grant}})
	controlDecision, err := control.Review(context.Background(), request)
	if err != nil {
		t.Fatalf("the marked in-memory assembly must still approve, got %v", err)
	}
	if controlDecision.Outcome != "approve" || controlDecision.PolicyHash != grant.Grant.PolicyHash {
		t.Fatalf("the in-memory assembly must bind the grant policy, got %+v", controlDecision)
	}
}

// TestP83T1bGrantExpiringDuringLedgerReadRefused is the T1b counterexample. The
// ledger is read at the moment of use and the clock is read AFTER that read, so a
// grant that expires while the read is in flight cannot authorise the review.
//
// Two facts are asserted, each with its own mutation contract:
//  1. ordering - the clock callback truncates the ledger file and returns an
//     ACTIVE instant; the grant still approves, which is only possible when the
//     load happened before the clock read. Moving `now := runtime.clock()` above
//     the load makes this leg refuse and turns it red.
//  2. expiry-at-use - with the ledger restored, the same runtime refuses once the
//     clock advances past the expiry. Evaluating expiry at assembly/grant-issue
//     time keeps the approve and turns this leg red.
func TestP83T1bGrantExpiringDuringLedgerReadRefused(t *testing.T) {
	root := t.TempDir()
	ledgerPath := filepath.Join(root, defaultApprovalsLedgerName)
	grant := p83Fix5Grant(t, nil)
	writeP83Fix4Ledger(t, ledgerPath, grant)

	runtime := p83Runtime(t, chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{grant}})
	runtime.authorizationLedgerPath = ledgerPath

	var clockErr error
	runtime.clock = func() time.Time {
		// The read has already happened when this runs, so rewriting the file here
		// can only be observed by a later load - which is exactly the ordering fact.
		if clockErr = os.WriteFile(ledgerPath, nil, 0o600); clockErr != nil {
			return p83Clock()
		}
		return p83Clock()
	}
	policyHash, err := runtime.operatorApprovedPolicyHash(p83RunID, p83TaskID, p83StepID)
	if clockErr != nil {
		t.Fatalf("rewrite the ledger from the clock: %v", clockErr)
	}
	if err != nil {
		t.Fatalf("the grant loaded before the clock read must still approve, got %v", err)
	}
	if policyHash != grant.Grant.PolicyHash {
		t.Fatalf("the approved policy must be the grant policy %s, got %s", grant.Grant.PolicyHash, policyHash)
	}
	// The truncation is now visible to a fresh read: the ledger file is the source.
	runtime.clock = func() time.Time { return p83Clock() }
	if _, err := runtime.operatorApprovedPolicyHash(p83RunID, p83TaskID, p83StepID); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("an emptied ledger must refuse on the next read, got %v", err)
	}

	// Leg 2: the ledger holds the grant again and the clock is past its expiry.
	writeP83Fix4Ledger(t, ledgerPath, grant)
	runtime.clock = func() time.Time { return p83Clock().Add(2 * time.Hour) }
	policyHash, err = runtime.operatorApprovedPolicyHash(p83RunID, p83TaskID, p83StepID)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("a grant expired during the read must be refused, got %v", err)
	}
	if policyHash != "" {
		t.Fatalf("a refused authorization must release no policy hash, got %q", policyHash)
	}
}

// TestP83T2aSupersededGrantCannotLendItsScope is the T2a counterexample: an older
// grant record of an authorizationId whose newer record no longer covers this
// task/step must not authorise the review.
//
// Mutation check: deleting the `ledger.LatestGrant(...)` supersession check makes
// the older record evaluate as active (EvaluateGrant reads the newest record) and
// the review approves with the older policy hash, turning this red.
func TestP83T2aSupersededGrantCannotLendItsScope(t *testing.T) {
	older := p83Fix5Grant(t, nil)
	newer := p83Fix5Grant(t, func(grant *chain.AuthorizationGrant) {
		grant.RecordID = "p83-grant-2"
		grant.PolicyHash = evidence.Digest("", []byte("p83-narrowed-policy"))
		grant.Scope.TaskIDs = []string{"task-other"}
		grant.Scope.StepIDs = []string{"step-other"}
	})

	request := chain.ReviewRequest{
		RunID:                 p83RunID,
		TaskID:                p83TaskID,
		Attempt:               1,
		CandidateSnapshotHash: evidence.Digest("", []byte("candidate")),
		EvidenceRootHash:      evidence.Digest("", []byte("evidence-root")),
	}

	superseded := p83Runtime(t, chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{older, newer}})
	decision, err := superseded.Review(context.Background(), request)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("a superseded grant must not authorise the review, got %+v, %v", decision, err)
	}
	if decision.Outcome != "" || decision.PolicyHash != "" {
		t.Fatalf("a superseded grant must produce no approve decision, got %+v", decision)
	}

	// Positive control: the SAME two records approve when the newer record covers
	// the task, so the refusal above is the supersession and nothing else.
	widened := p83Fix5Grant(t, func(grant *chain.AuthorizationGrant) {
		grant.RecordID = "p83-grant-2"
		grant.PolicyHash = evidence.Digest("", []byte("p83-narrowed-policy"))
	})
	control := p83Runtime(t, chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{older, widened}})
	controlDecision, err := control.Review(context.Background(), request)
	if err != nil {
		t.Fatalf("the newest covering grant must approve, got %v", err)
	}
	if controlDecision.Outcome != "approve" || controlDecision.PolicyHash != widened.Grant.PolicyHash {
		t.Fatalf("the review must bind the newest grant policy %s, got %+v", widened.Grant.PolicyHash, controlDecision)
	}
}

// TestP83T2bSupervisedConfirmationIsSeparateFromPolicyApproval is the T2b
// counterexample: the risk confirmation is a separate requirement, not the policy
// approval. The reviewed policy hash is the frozen handoff policy digest (the
// confirmation is not folded into it), and removing the confirmation still refuses
// the review while that policy hash is unchanged.
//
// Mutation check: folding the confirmation hash into supervisedReviewPolicyHash
// breaks the independent digest comparison; deleting the confirmation requirement
// from supervisedRuntime.Review makes the withdrawal leg approve. Either turns this
// red.
func TestP83T2bSupervisedConfirmationIsSeparateFromPolicyApproval(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.now = p83FixedNow()
	harness.driveToWindow()
	runtime, state := p83F3SupervisedRuntime(t, harness)

	session, found, err := runtime.loadSession(state)
	if err != nil || !found {
		t.Fatalf("the open window must persist its handoff session: found=%v err=%v", found, err)
	}
	// The reviewed policy hash is exactly the frozen handoff policy digest.
	expected := digest(supervisedReviewPolicyDomain, session.HandoffPolicyHash)
	if got := supervisedReviewPolicyHash(session); got != expected {
		t.Fatalf("the review policy hash must be the frozen policy digest %s, got %s", expected, got)
	}

	request := chain.ReviewRequest{
		RunID:                 supervisedTestRunID,
		TaskID:                supervisedTestTaskID,
		Attempt:               1,
		CandidateSnapshotHash: harness.loadBeforeManifest().ManifestHash,
		EvidenceRootHash:      evidence.Digest("", []byte("p83-t2b-evidence-root")),
	}
	decision, err := runtime.Review(context.Background(), request)
	if err != nil {
		t.Fatalf("the unchanged window must approve, got %v", err)
	}
	if decision.Outcome != "approve" || decision.PolicyHash != expected {
		t.Fatalf("the review must approve under the frozen policy digest %s, got %+v", expected, decision)
	}

	// Withdraw only the risk confirmation: the policy hash still matches, and the
	// review must still refuse.
	confirmations := runtime.confirmationsPath()
	if err := os.Remove(confirmations); err != nil {
		t.Fatal(err)
	}
	decision, err = runtime.Review(context.Background(), request)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("a withdrawn risk confirmation must refuse the review, got %+v, %v", decision, err)
	}
	if decision.Outcome != "" || decision.PolicyHash != "" {
		t.Fatalf("a refused review must produce no decision, got %+v", decision)
	}
}

// TestP83T3aAwayAuthorizationLedgerPathRefused is the T3a counterexample: an
// --authorization-ledger that points away from the chain-config-derived path is a
// refusal, never an override, and the secret store is never reached.
//
// Mutation check: restoring "honour the caller-selected path" (using selectedLedger
// directly) admits the probe from the foreign ledger and turns this red.
func TestP83T3aAwayAuthorizationLedgerPathRefused(t *testing.T) {
	root := t.TempDir()
	configPath := writeAIProbeTestChain(t, root)
	executablePath := filepath.Join(root, "copilot.exe")
	if err := os.WriteFile(executablePath, []byte("p83-fix5-approved-probe"), 0o700); err != nil {
		t.Fatal(err)
	}
	executableHash, err := adapters.CopilotCLIExecutableHash(executablePath)
	if err != nil {
		t.Fatal(err)
	}
	// The trusted ledger (next to the chain config) approves the receiver, so only
	// the disagreement can be the refusal.
	writeP83Fix4Ledger(t, filepath.Join(root, defaultApprovalsLedgerName), p83Fix4Grant(t, executableHash, nil))
	away := filepath.Join(root, "elsewhere", defaultApprovalsLedgerName)
	if err := os.MkdirAll(filepath.Dir(away), 0o755); err != nil {
		t.Fatal(err)
	}
	writeP83Fix4Ledger(t, away, p83Fix4Grant(t, executableHash, nil))

	secrets := &recordingProbeSecrets{credential: []byte("p83-fix5-credential")}
	cli := newTestCLI(root)
	cli.secretManager = secrets
	checkerEntries := 0
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, aiProbeReceiver, string) (adapters.AIAvailabilityRecord, error) {
		checkerEntries++
		return adapters.AIAvailabilityRecord{}, nil
	}
	outputPath := filepath.Join(root, "p83-fix5-availability.json")
	code, stdout, stderr := runCLI(t, cli,
		"ai", "check", "--chain", configPath,
		"--authorization-ledger", away,
		"--channel", "agent-runner-cli", "--copilot", executablePath,
		"--max-requests", "1", "--out", outputPath, "--json",
	)
	if code != exitFailure {
		t.Fatalf("a ledger path that disagrees with the derived one must fail closed, got code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout+stderr, "disagrees with the trusted") {
		t.Fatalf("the disagreement must be the refusing guard, got stdout=%q stderr=%q", stdout, stderr)
	}
	if checkerEntries != 0 || secrets.resolveCalls != 0 {
		t.Fatalf("a refused receiver must not resolve a credential, checker=%d resolutions=%d", checkerEntries, secrets.resolveCalls)
	}
	if _, statErr := os.Lstat(outputPath); !os.IsNotExist(statErr) {
		t.Fatalf("no availability record may be written, stat err=%v", statErr)
	}
}

// TestP83T3bSameReceiverDifferentSecretRefRefused is the T3b counterexample: two
// profiles can share host/port (and even model) and differ only in secretRef, so
// the grant must approve the profile hash too. Here the grant carries the receiver
// binding and the executable digest of the profile that is actually configured but
// names the OTHER profile's hash (the one with the different secretRef), so the
// profile-hash comparison is the only identity left that can refuse it.
//
// Mutation check: deleting the `authorizationCarriesHash(grant.Evidence,
// profileHash)` check admits the same-address profile and turns this red.
func TestP83T3bSameReceiverDifferentSecretRefRefused(t *testing.T) {
	root := t.TempDir()
	configPath := writeAIProbeTestChain(t, root)
	executablePath := filepath.Join(root, "copilot.exe")
	if err := os.WriteFile(executablePath, []byte("p83-fix5-approved-probe"), 0o700); err != nil {
		t.Fatal(err)
	}
	executableHash, err := adapters.CopilotCLIExecutableHash(executablePath)
	if err != nil {
		t.Fatal(err)
	}
	configured, found := configuredAIProfile(configuredAIChain(), "agent-runner-cli")
	if !found {
		t.Fatal("the shared AI chain must configure the agent-runner-cli channel")
	}
	other := configured
	other.SecretRef = "windows-credential:ProofRail/p83-fix5-other-secret-ref"
	if other.BaseURL != configured.BaseURL || other.Model != configured.Model {
		t.Fatal("test precondition failed: only secretRef may differ")
	}
	otherHash := mustAIProfileHash(t, other)
	configuredHash := mustAIProfileHash(t, configured)
	if otherHash == configuredHash {
		t.Fatal("test precondition failed: a different secretRef must change the profile hash")
	}
	// The binding names the profile that is actually configured, so the grant looks
	// right for this receiver except for the profile hash it approved.
	binding := evidence.Digest(aiProbeReceiverBindingDomain, []byte(configuredHash+"\n"+"agent-runner-cli"+"\n"+executableHash))
	writeP83Fix4Ledger(t, filepath.Join(root, defaultApprovalsLedgerName), p83Fix4Grant(t, executableHash, func(grant *chain.AuthorizationGrant) {
		grant.Evidence = []string{executableHash, otherHash, binding}
	}))

	secrets := &recordingProbeSecrets{credential: []byte("p83-fix5-credential")}
	cli := newTestCLI(root)
	cli.secretManager = secrets
	checkerEntries := 0
	cli.checkAIAvailability = func(context.Context, adapters.AIProviderProfile, string, int, aiProbeReceiver, string) (adapters.AIAvailabilityRecord, error) {
		checkerEntries++
		return adapters.AIAvailabilityRecord{}, nil
	}
	outputPath := filepath.Join(root, "p83-fix5-availability.json")
	code, stdout, stderr := runCLI(t, cli,
		"ai", "check", "--chain", configPath,
		"--channel", "agent-runner-cli", "--copilot", executablePath,
		"--max-requests", "1", "--out", outputPath, "--json",
	)
	if code != exitFailure {
		t.Fatalf("a grant approving another profile must refuse, got code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !strings.Contains(stdout+stderr, "receiver binding") {
		t.Fatalf("the receiver binding must be the refusing guard, got stdout=%q stderr=%q", stdout, stderr)
	}
	if checkerEntries != 0 || secrets.resolveCalls != 0 {
		t.Fatalf("a refused profile must not resolve a credential, checker=%d resolutions=%d", checkerEntries, secrets.resolveCalls)
	}
	if _, statErr := os.Lstat(outputPath); !os.IsNotExist(statErr) {
		t.Fatalf("no availability record may be written, stat err=%v", statErr)
	}
}

// TestP83T9bHostAuthProfileWithoutBaseURLAdmitted is the T9b counterexample: the
// receiver gate is not applied to a profile that names no base URL, because such a
// profile authenticates by host identity and contacts no receiver. The gate must
// admit it instead of refusing every host-auth channel.
//
// Mutation check: deleting the `if profile.BaseURL == "" { return nil }` early
// return refuses this profile with ErrAIReceiverBinding and turns this red.
func TestP83T9bHostAuthProfileWithoutBaseURLAdmitted(t *testing.T) {
	profile, found := configuredAIProfile(configuredAIChain(), "sessionbridge-silent")
	if !found {
		t.Fatal("the shared AI chain must configure the sessionbridge-silent channel")
	}
	if profile.BaseURL != "" || profile.AuthMode != "host" {
		t.Fatalf("test precondition failed: the host-auth profile must name no base URL, got %+v", profile)
	}
	executablePath := filepath.Join(t.TempDir(), "copilot.exe")
	if err := os.WriteFile(executablePath, []byte("p83-fix5-host-auth-probe"), 0o700); err != nil {
		t.Fatal(err)
	}
	binding, executableHash, err := validateAIProbeBinding(profile, "sessionbridge-silent", executablePath)
	if err != nil {
		t.Fatalf("the host-auth profile must bind, got %v", err)
	}
	if err := authorizeAIProbeReceiver(chain.AuthorizationLedger{}, profile, binding, executableHash, p83Clock()); err != nil {
		t.Fatalf("a host-auth profile without a base URL must not be refused, got %v", err)
	}
}

// TestP83T5aSecondLiveHolderCannotAdoptLease is the T5a counterexample: a lease
// record written by a holder whose process is still alive may not be adopted by a
// second process of the same run.
//
// The "other process" is a re-executed copy of this test binary that sleeps, so
// its identity is a real, live OS process that is not this one.
//
// Mutation check: deleting the publishLeasePreviousHolderStopped call lets the
// second claim adopt the record and turns this red.
func TestP83T5aSecondLiveHolderCannotAdoptLease(t *testing.T) {
	runRoot := t.TempDir()
	holder := exec.Command(os.Args[0], "-test.run=TestP83Fix5T5HolderHelperProcess")
	holder.Env = append(os.Environ(), "P83_FIX5_T5_HOLDER=1")
	if err := holder.Start(); err != nil {
		t.Skipf("the host cannot start the helper holder process: %v", err)
	}
	defer func() {
		_ = holder.Process.Kill()
		_, _ = holder.Process.Wait()
	}()
	identity, err := guard.InspectProcess(holder.Process.Pid)
	if err != nil {
		t.Skipf("the host cannot inspect the helper holder process: %v", err)
	}
	if self, selfErr := guard.InspectProcess(os.Getpid()); selfErr == nil && self == identity {
		t.Skip("the host reports the same identity for the helper and this process")
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Claim(context.Background(), publishLeaseResourceID(p83RunID), publishLeaseConsumerID, publishLeaseClaimID(p83RunID), identity, publishLeaseTTL); err != nil {
		t.Fatalf("the helper holder must be able to claim the lease: %v", err)
	}
	if _, err := claimPublishLease(context.Background(), runRoot, p83RunID); !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a lease held by a live process must not be adopted, got %v", err)
	}
}

// TestP83Fix5T5HolderHelperProcess is the sleeping holder the T5a test inspects.
// It only sleeps; every other run of this function skips.
func TestP83Fix5T5HolderHelperProcess(t *testing.T) {
	if os.Getenv("P83_FIX5_T5_HOLDER") != "1" {
		t.Skip("helper process for the T5a lease exclusivity test")
	}
	time.Sleep(30 * time.Second)
}
