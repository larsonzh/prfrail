package console

// P83 FIX ROUND 6 counterexamples (items T1a, T1b, T3b, T5a, T5d and T9a).
//
// Every test below is written so that reverting the one guard it names turns it
// red; the mutation contract is stated in each test's doc comment.

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/tickets"
)

const (
	p83Fix6T5aRunID = "run-p83-fix6-t5a"
	p83Fix6T5dRunID = "run-p83-fix6-t5d"
	// p83Fix6GonePID is above every pid a Windows host can report, so no process
	// is inspectable at it and the identity is provably gone.
	p83Fix6GonePID = 0x7FFFFFFF
)

// p83Fix6WriteAgentRunnerRecords writes one decode-complete AgentRunner record
// set and returns the flags that declare it. The records are intentionally not
// cross-checked by the CLI (that happens in the assembly), so the fixture only has
// to pass the decoders the CLI runs before it derives the ledger path.
func p83Fix6WriteAgentRunnerRecords(t *testing.T, dir, runID string) []string {
	t.Helper()
	createdAt := "2026-11-20T09:00:00.000Z"
	write := func(name string, value any) string {
		t.Helper()
		encoded, err := json.Marshal(value)
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, encoded, 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}

	capability, err := adapters.NewAgentRunnerCapabilityRecord(adapters.AgentRunnerCapability{
		ProbeID:        "p83-fix6-capability",
		AdapterID:      "p83-fix6-adapter",
		ProbedAt:       createdAt,
		OS:             "windows",
		Arch:           "amd64",
		PlatformHash:   evidence.Digest("", []byte("p83-fix6-platform")),
		ExecutableHash: evidence.Digest("", []byte("p83-fix6-executable")),
		Version:        "p83-fix6-agent-0.1.0",
		ConfigHash:     evidence.Digest("", []byte("p83-fix6-config")),
		Disposition:    "compatible",
		Matrix:         assemblyCapabilityMatrix("compatible"),
		Evidence:       []string{evidence.Digest("", []byte("p83-fix6-capability-evidence"))},
	})
	if err != nil {
		t.Fatalf("build the capability record: %v", err)
	}
	enforcement, err := adapters.NewAgentRunnerEnforcementRecord(adapters.AgentRunnerEnforcement{
		EnforcementID:         "p83-fix6-enforcement",
		ProviderID:            "p83-fix6-provider",
		ProbedAt:              createdAt,
		CapabilityRecordHash:  capability.RecordHash,
		ExecutableHash:        capability.Capability.ExecutableHash,
		ConfigHash:            capability.Capability.ConfigHash,
		PlatformHash:          capability.Capability.PlatformHash,
		EnforcementConfigHash: evidence.Digest("", []byte("p83-fix6-enforcement-config")),
		Disposition:           "verified",
		Matrix: adapters.AgentRunnerEnforcementMatrix{
			ToolControl:    assemblyVerifiedFinding("p83-fix6-enforcement-tool"),
			NetworkControl: assemblyVerifiedFinding("p83-fix6-enforcement-network"),
		},
		Evidence: []string{evidence.Digest("", []byte("p83-fix6-enforcement-evidence"))},
	})
	if err != nil {
		t.Fatalf("build the enforcement record: %v", err)
	}
	availability, err := adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
		ProbeID:           "p83-fix6-availability",
		ProfileID:         "deepseek-anthropic",
		ProfileConfigHash: evidence.Digest("", []byte("p83-fix6-profile")),
		Channel:           "agent-runner-cli",
		ProbedAt:          createdAt,
		Status:            "available",
		RequestsUsed:      1,
		Evidence:          []string{evidence.Digest("", []byte("p83-fix6-availability-evidence"))},
	})
	if err != nil {
		t.Fatalf("build the availability record: %v", err)
	}
	request := p51WriteRequestRecord(t, dir, "p83-fix6-request.json", runID, p51TaskOne, p51StepOne, "p83-fix6-request-one")
	return []string{
		"--request-record", request,
		"--capability-record", write("p83-fix6-capability.json", capability),
		"--enforcement-record", write("p83-fix6-enforcement.json", enforcement),
		"--availability-record", write("p83-fix6-availability.json", availability),
		"--agent-cli", filepath.Join(dir, "p83-fix6-agent-cli.exe"),
		"--agent-version", "p83-fix6-agent-0.1.0",
		"--probe-permit", filepath.Join(dir, "p83-fix6-probe-permit.marker"),
	}
}

// TestP83Fix6T1aRunDerivesTheLedgerPathsWithoutFlags is the T1a counterexample:
// `run` no longer requires --authorization-ledger and --cost-ledger, the derived
// location under the durable run root is used, and a flag that names anything else
// is still refused.
//
// Mutation check: adding the two ledger fields back to agentRunnerRunInputs.complete
// makes the no-flag leg stop at the noop-only refusal and turns this red.
func TestP83Fix6T1aRunDerivesTheLedgerPathsWithoutFlags(t *testing.T) {
	runRoot := t.TempDir()
	recordsDir := t.TempDir()
	cli := newTestCLI(t.TempDir())
	chainPath := p51WriteChainConfig(t, p51SingleTaskChain()...)
	records := p83Fix6WriteAgentRunnerRecords(t, recordsDir, p51RunID)

	// Leg 1: no ledger flags at all. The declaration is complete, so the command
	// must reach the derived ledgers; the absent cost ledger is what stops it, and
	// the refusal names the file under the durable run root.
	args := append([]string{"run", "--chain", chainPath, "--run-id", p51RunID, "--run-dir", runRoot}, records...)
	code, stdout, stderr := runCLI(t, cli, args...)
	if code == exitSuccess {
		t.Fatalf("an absent derived cost ledger must fail closed, stdout=%s stderr=%s", stdout, stderr)
	}
	if strings.Contains(stderr, "local runtime only supports noop steps") {
		t.Fatalf("the two ledger flags must not be required any more, stderr=%q", stderr)
	}
	wantDerived := filepath.Join(runRoot, defaultCostLedgerName)
	if !strings.Contains(stderr, wantDerived) {
		t.Fatalf("the refusal must name the derived cost ledger %q, stderr=%q", wantDerived, stderr)
	}

	// Leg 2: a disagreeing --authorization-ledger is refused, not honoured.
	away := filepath.Join(t.TempDir(), defaultApprovalsLedgerName)
	disagreeing := append(append([]string{}, args...), "--authorization-ledger", away)
	code, stdout, stderr = runCLI(t, cli, disagreeing...)
	if code == exitSuccess {
		t.Fatalf("a disagreeing --authorization-ledger must fail closed, stdout=%s stderr=%s", stdout, stderr)
	}
	if !strings.Contains(stderr, "disagrees with the trusted") {
		t.Fatalf("the disagreement must be the refusing guard, stderr=%q", stderr)
	}

	// Control: the same invocation without --agent-cli keeps the pre-existing
	// noop-only refusal, so leg 1 reached the ledgers because the declaration was
	// complete and not because the completeness gate disappeared.
	control := make([]string, 0, len(args))
	for index := 0; index < len(args); index++ {
		if args[index] == "--agent-cli" {
			index++
			continue
		}
		control = append(control, args[index])
	}
	code, _, stderr = runCLI(t, cli, control...)
	if code == exitSuccess || !strings.Contains(stderr, "local runtime only supports noop steps") {
		t.Fatalf("an incomplete declaration must keep the noop-only refusal, code=%d stderr=%q", code, stderr)
	}
}

// TestP83Fix6T1bRunDirDerivesTheDefaultLedgerPath is the T1b counterexample:
// `approvals list` and `cost report` with --run-dir read the ledger of that run
// root, not a file in the working directory.
//
// Mutation check: deleting the runRoot branch of resolveApprovalsLedgerPath (or of
// resolveCostLedgerPath) falls back to the working directory, where no ledger
// exists, and turns the corresponding leg red.
func TestP83Fix6T1bRunDirDerivesTheDefaultLedgerPath(t *testing.T) {
	runRoot := t.TempDir()
	writeP83Fix4Ledger(t, filepath.Join(runRoot, defaultApprovalsLedgerName), p83Fix5Grant(t, nil))
	amount := int64(1000)
	writeCostLedgerFixture(t, runRoot, tickets.CostPricingDocumented, &amount)
	// The working directory is a DIFFERENT directory, so a cwd default cannot pass.
	cli := newTestCLI(t.TempDir())

	code, stdout, stderr := runCLI(t, cli, "approvals", "list", "--run-dir", runRoot, "--json")
	if code != exitSuccess {
		t.Fatalf("approvals list code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	var approvals ApprovalsReport
	if err := json.Unmarshal(decodeResponse(t, stdout).Data, &approvals); err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(runRoot, defaultApprovalsLedgerName); approvals.LedgerPath != want {
		t.Fatalf("approvals list must read %q, got %q", want, approvals.LedgerPath)
	}
	if len(approvals.Grants) != 1 || len(approvals.Warnings) != 0 {
		t.Fatalf("the run's own grant must be listed without a missing-ledger warning, got %+v", approvals)
	}

	code, stdout, stderr = runCLI(t, cli, "cost", "report", "--run-dir", runRoot, "--json")
	if code != exitSuccess {
		t.Fatalf("cost report code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	var report CostReport
	if err := json.Unmarshal(decodeResponse(t, stdout).Data, &report); err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(runRoot, defaultCostLedgerName); report.LedgerPath != want {
		t.Fatalf("cost report must read %q, got %q", want, report.LedgerPath)
	}
	// The counters prove the derived file was actually read, not merely named.
	if report.ReservedCalls != 1 || report.SettledCalls != 1 || len(report.Warnings) != 0 {
		t.Fatalf("the run's own cost ledger must be reported, got %+v", report)
	}
}

// TestP83Fix6T3bEmptyApprovedDigestDoesNotDegradeToAnUnboundPath is the T3b
// counterexample: a receiver without an approved digest is refused instead of
// falling back to the executor's unpinned compatibility path (which would start a
// path that no approval ever named).
//
// Mutation check: deleting the empty-digest refusal in defaultAICheck builds an
// executor with no pinned digest, resolves the credential and turns this red.
func TestP83Fix6T3bEmptyApprovedDigestDoesNotDegradeToAnUnboundPath(t *testing.T) {
	root := t.TempDir()
	profile, found := configuredAIProfile(configuredAIChain(), "agent-runner-cli")
	if !found {
		t.Fatal("the shared AI chain must configure the agent-runner-cli channel")
	}
	executable := filepath.Join(root, "copilot.exe")
	if err := os.WriteFile(executable, []byte("p83-fix6-t3b-approved"), 0o700); err != nil {
		t.Fatal(err)
	}
	secrets := &recordingProbeSecrets{credential: []byte("p83-fix6-t3b-credential")}
	cli := newTestCLI(root)
	cli.secretManager = secrets

	_, err := cli.defaultAICheck(context.Background(), profile, "agent-runner-cli", 1, aiProbeReceiver{Executable: executable}, root)
	if !errors.Is(err, ErrAIReceiverBinding) {
		t.Fatalf("a receiver without an approved digest must be refused as a binding failure, got %v", err)
	}
	if !strings.Contains(err.Error(), "approved executable digest") {
		t.Fatalf("the refusal must name the missing approved digest, got %v", err)
	}
	if secrets.resolveCalls != 0 {
		t.Fatalf("an unbound receiver must not resolve a credential, got %d resolutions", secrets.resolveCalls)
	}
}

// TestP83Fix6T9aHostAuthProfileSucceedsWithoutAnyLedgerFile is the T9a
// counterexample: a host-auth profile names no receiver, so it needs no grant and
// must reach the availability check even when no ledger file exists at all.
//
// Mutation check: dropping the `profile.BaseURL != ""` guard in executeAICheck
// loads the absent ledger and turns this red.
func TestP83Fix6T9aHostAuthProfileSucceedsWithoutAnyLedgerFile(t *testing.T) {
	root := t.TempDir()
	config := configuredAIChain()
	// Bind the CLI channel to the host-auth profile: it authenticates by host
	// identity, names no base URL and contacts no receiver.
	for index := range config.AI.Channels {
		if config.AI.Channels[index].Channel == "agent-runner-cli" {
			config.AI.Channels[index].ProfileID = "github"
		}
	}
	if err := ValidateChainConfig(config); err != nil {
		t.Fatalf("the host-auth fixture chain must validate: %v", err)
	}
	profile, found := configuredAIProfile(config, "agent-runner-cli")
	if !found || profile.BaseURL != "" || profile.AuthMode != "host" {
		t.Fatalf("test precondition failed: the bound profile must be host-auth without a base URL, got %+v", profile)
	}
	encoded, err := EncodeChainConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(root, DefaultChainConfigName)
	if err := os.WriteFile(configPath, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(root, "copilot.exe")
	if err := os.WriteFile(executable, []byte("p83-fix6-t9a-host-auth"), 0o700); err != nil {
		t.Fatal(err)
	}
	ledgerPath := filepath.Join(root, defaultApprovalsLedgerName)
	if _, statErr := os.Stat(ledgerPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("test precondition failed: no ledger file may exist yet, stat err=%v", statErr)
	}

	cli := newTestCLI(root)
	checked := false
	cli.checkAIAvailability = func(_ context.Context, bound adapters.AIProviderProfile, channel string, maximumRequests int, receiver aiProbeReceiver, workingDirectory string) (adapters.AIAvailabilityRecord, error) {
		checked = true
		if bound.ProfileID != "github" || channel != "agent-runner-cli" || maximumRequests != 1 {
			t.Fatalf("unexpected binding: profile=%s channel=%s requests=%d", bound.ProfileID, channel, maximumRequests)
		}
		if receiver.ExecutableHash == "" {
			t.Fatal("the host-auth receiver must still pin the executable digest")
		}
		return adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
			ProbeID:           "probe-p83-fix6-t9a",
			ProfileID:         bound.ProfileID,
			ProfileConfigHash: mustAIProfileHash(t, bound),
			Channel:           channel,
			ProbedAt:          "2026-09-08T10:29:00.000Z",
			Status:            "available",
			RequestsUsed:      1,
			Evidence:          []string{evidence.Digest("", []byte("p83-fix6-t9a-availability"))},
		})
	}
	outputPath := filepath.Join(root, "p83-fix6-t9a-availability.json")
	code, stdout, stderr := runCLI(t, cli,
		"ai", "check", "--chain", configPath, "--channel", "agent-runner-cli",
		"--copilot", executable, "--max-requests", "1", "--out", outputPath, "--json",
	)
	if code != exitSuccess {
		t.Fatalf("a host-auth profile must be admitted without a ledger, code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if !checked {
		t.Fatal("the availability check must run for a host-auth profile")
	}
	if _, statErr := os.Stat(ledgerPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the probe must neither create nor require a ledger file, stat err=%v", statErr)
	}
	if _, statErr := os.Stat(outputPath); statErr != nil {
		t.Fatalf("the availability record must be written, stat err=%v", statErr)
	}
}

// p83Fix6SupervisedRuntime builds the minimal supervised runtime that owns the
// publish lease files of one run, without an engine.
func p83Fix6SupervisedRuntime(t *testing.T, runRoot, runID string) (*supervisedRuntime, *supervisedTaskState) {
	t.Helper()
	state := &supervisedTaskState{taskID: "task-p83-fix6", stepID: "step-p83-fix6", attempt: 1}
	runtime := &supervisedRuntime{
		runID:   runID,
		runRoot: runRoot,
		plan:    []*supervisedTaskState{state},
		tasks:   map[string]*supervisedTaskState{state.key(): state},
		clock:   func() time.Time { return time.Now().UTC() },
		ids:     func() string { return "p83-fix6-id" },
	}
	if err := os.MkdirAll(runtime.taskDir(state), 0o755); err != nil {
		t.Fatal(err)
	}
	return runtime, state
}

// TestP83Fix6T5aThirdSuccessorIsRefusedByTheFencingToken is the T5a
// counterexample: A claims and stops, B adopts the record (a transfer at a new
// generation), and a THIRD holder presenting the predecessor token is refused with
// the fencing token mismatch.
//
// Mutation checks. (1) Returning the old record from claimPublishLease instead of
// reclaiming it lets B keep A's generation. (2) Deleting the token comparison in
// LeaseStore.Reclaim lets C take the record. Either turns this red.
func TestP83Fix6T5aThirdSuccessorIsRefusedByTheFencingToken(t *testing.T) {
	runRoot := t.TempDir()
	ctx := context.Background()
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	resource := publishLeaseResourceID(p83Fix6T5aRunID)
	// A is a holder that stopped and left its record behind. On this host an exited
	// process stays openable for a while, so "stopped" is expressed the way the
	// store observes it: an identity no process is inspectable at any more.
	gone := guard.ProcessIdentity{PID: p83Fix6GonePID, StartToken: "p83-fix6-gone-holder"}
	if _, err := guard.InspectProcess(gone.PID); err == nil {
		t.Fatalf("test precondition failed: pid %d is inspectable, so it is not a gone identity", gone.PID)
	}
	predecessor, err := store.Claim(ctx, resource, publishLeaseConsumerID, publishLeaseClaimID(p83Fix6T5aRunID), gone, publishLeaseTTL)
	if err != nil {
		t.Fatalf("the predecessor must be able to claim the lease: %v", err)
	}

	// B (this process) adopts; adoption is a transfer, not a hand-out.
	adopted, err := claimPublishLease(ctx, runRoot, p83Fix6T5aRunID)
	if err != nil {
		t.Fatalf("a stopped holder must be adoptable, got %v", err)
	}
	current, err := store.Current(ctx, resource)
	if err != nil {
		t.Fatal(err)
	}
	if current.Token.Generation != predecessor.Token.Generation+1 {
		t.Fatalf("adoption must advance the generation to %d, got %+v", predecessor.Token.Generation+1, current.Token)
	}
	if current.Process != publishLeaseProcess(p83Fix6T5aRunID) || current.Process == predecessor.Process {
		t.Fatalf("adoption must name this process as the holder, got %+v", current.Process)
	}
	if adopted.LeaseHash != current.LeaseHash || adopted.Token != current.Token {
		t.Fatalf("the successor's proof must name the transferred record, got %+v", adopted)
	}

	// C presents the predecessor token it could read from the run root.
	third := guard.ProcessIdentity{PID: os.Getpid(), StartToken: "p83-fix6-third-successor"}
	if _, err := store.Reclaim(ctx, resource, predecessor.Token, third, publishLeaseTTL, func(guard.Lease) error { return nil }); !errors.Is(err, guard.ErrFencingTokenMismatch) {
		t.Fatalf("a third holder presenting the predecessor token must be refused, got %v", err)
	}
	after, err := store.Current(ctx, resource)
	if err != nil {
		t.Fatal(err)
	}
	if after.LeaseHash != current.LeaseHash || after.Token != current.Token || after.Process != current.Process {
		t.Fatalf("the refused third holder must not move the record: got %+v", after)
	}
}

// TestP83Fix6T5dSupervisedPublishRefusesAForeignOnDiskProof is the T5d
// counterexample: the store-local publish proof is a claim, not an ownership
// proof, so a proof whose record names another live process is refused both when
// the window is opened and when the promotion verifies it.
//
// Mutation check: deleting the `lease.Process != holder` refusal in
// verifyPublishLease makes the publish leg return a lease hash and turns this red.
func TestP83Fix6T5dSupervisedPublishRefusesAForeignOnDiskProof(t *testing.T) {
	runRoot := t.TempDir()
	ctx := context.Background()
	holder := exec.Command(os.Args[0], "-test.run=TestP83Fix6T5dForeignHolderHelperProcess")
	holder.Env = append(os.Environ(), "P83_FIX6_T5D_HOLDER=1")
	if err := holder.Start(); err != nil {
		t.Skipf("the host cannot start the foreign holder process: %v", err)
	}
	defer func() {
		_ = holder.Process.Kill()
		_, _ = holder.Process.Wait()
	}()
	identity, err := guard.InspectProcess(holder.Process.Pid)
	if err != nil {
		t.Skipf("the host cannot inspect the foreign holder process: %v", err)
	}

	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	resource := publishLeaseResourceID(p83Fix6T5dRunID)
	lease, err := store.Claim(ctx, resource, publishLeaseConsumerID, publishLeaseClaimID(p83Fix6T5dRunID), identity, publishLeaseTTL)
	if err != nil {
		t.Fatalf("the foreign holder must be able to claim the lease: %v", err)
	}
	foreign := publishLeaseProof{ResourceID: lease.ResourceID, Token: lease.Token, LeaseHash: lease.LeaseHash}

	runtime, state := p83Fix6SupervisedRuntime(t, runRoot, p83Fix6T5dRunID)
	// The proof is read from the task's own directory, exactly like a re-entrant
	// supervised run reads it back.
	if err := writeSupervisedJSON(runtime.publishLeasePath(state), foreign); err != nil {
		t.Fatal(err)
	}
	if err := runtime.ensurePublishLease(ctx, state); !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a proof held by a live foreign process must not be adopted, got %v", err)
	}
	if _, err := verifyPublishLease(ctx, runRoot, p83Fix6T5dRunID, foreign, time.Now().UTC()); !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("publish must refuse a lease that names a foreign holder, got %v", err)
	} else if !strings.Contains(err.Error(), "not this process") {
		t.Fatalf("the refusal must state that the record names another holder, got %v", err)
	}
	untouched, err := store.Current(ctx, resource)
	if err != nil {
		t.Fatal(err)
	}
	if untouched.Process != identity || untouched.LeaseHash != lease.LeaseHash {
		t.Fatalf("the refused adoption must leave the foreign record in place, got %+v", untouched)
	}

	// Control: the same verification succeeds for a lease this process claimed, so
	// the refusal above is the holder proof and nothing else.
	controlRunID := "run-p83-fix6-t5d-control"
	controlResource := publishLeaseResourceID(controlRunID)
	controlLease, err := store.Claim(ctx, controlResource, publishLeaseConsumerID, publishLeaseClaimID(controlRunID), publishLeaseProcess(controlRunID), publishLeaseTTL)
	if err != nil {
		t.Fatalf("the control claim must succeed: %v", err)
	}
	control := publishLeaseProof{ResourceID: controlLease.ResourceID, Token: controlLease.Token, LeaseHash: controlLease.LeaseHash}
	controlHash, err := verifyPublishLease(ctx, runRoot, controlRunID, control, time.Now().UTC())
	if err != nil || controlHash != controlLease.LeaseHash {
		t.Fatalf("a lease this process holds must verify, got hash=%q err=%v", controlHash, err)
	}
}

// TestP83Fix6T5dForeignHolderHelperProcess is the live foreign holder the T5d
// test inspects. It only sleeps; every other run of this function skips.
func TestP83Fix6T5dForeignHolderHelperProcess(t *testing.T) {
	if os.Getenv("P83_FIX6_T5D_HOLDER") != "1" {
		t.Skip("helper process for the T5d publish ownership test")
	}
	time.Sleep(30 * time.Second)
}
