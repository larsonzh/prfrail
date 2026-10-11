package console

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
	"github.com/larsonzh/prfrail/internal/gates"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// This file is the ⑨ evidence carrier for the production supervised-black-box
// assembly. Every delivery is answered by the deterministic scripted stub
// (adapters.StubVisibleClient): no network, no IPC and no external process is
// ever involved.

const (
	supervisedTestRunID  = "run-t028"
	supervisedTestTaskID = "task-supervised"
	supervisedTestStepID = "code-supervised"
	supervisedTestHookID = "hook-one"
)

// --- gates doubles -----------------------------------------------------------

// supervisedFakeExecutor declares the full capability set the gate policy
// demands and answers with a deterministic execution.
type supervisedFakeExecutor struct{ fail bool }

func (supervisedFakeExecutor) Capabilities() gates.Capabilities {
	return gates.Capabilities{Process: true, Timeout: true, OutputLimit: true, MemoryLimit: true, ProcessLimit: true, NetworkDeny: true}
}

func (executor supervisedFakeExecutor) Execute(context.Context, gates.Prepared, gates.Hook) gates.Execution {
	if executor.fail {
		exit := 1
		return gates.Execution{Outcome: "exited", ExitCode: &exit, RunnerEvidence: []string{digest("proofrail:test-gate-runner:1\n", "failed")}}
	}
	exit := 0
	return gates.Execution{Outcome: "exited", ExitCode: &exit, Stdout: []byte("ok\n"), RunnerEvidence: []string{digest("proofrail:test-gate-runner:1\n", "passed")}}
}

type supervisedMemoryObjects struct{}

func (supervisedMemoryObjects) PutObject(_ context.Context, data []byte) (string, error) {
	return evidence.Digest("", data), nil
}

type supervisedMemoryResults struct{ records []gates.ResultRecord }

func (results *supervisedMemoryResults) Append(_ context.Context, record gates.ResultRecord) error {
	results.records = append(results.records, record)
	return nil
}

func supervisedTestHook(id string) gates.Hook {
	return gates.Hook{
		ID:          id,
		Kind:        "verify",
		Runner:      gates.ProcessRunner{Type: "process", Executable: "proofrail-test-hook"},
		EffectClass: "read-only",
		OnFail:      gates.FailurePolicy{Action: "fail-stop"},
		Timeout:     time.Second,
		Resources:   gates.ResourceLimits{MemoryBytes: 1 << 20, OutputBytes: 1 << 16, ProcessCount: 1},
		Network:     gates.NetworkPolicy{Mode: "deny"},
	}
}

// --- fixture -----------------------------------------------------------------

func supervisedTestDefinition(timeoutMs int) chain.Definition {
	return chain.Definition{ID: "chain-supervised", Tasks: []chain.Task{{
		ID: supervisedTestTaskID,
		Steps: []chain.Step{{
			ID:   supervisedTestStepID,
			Kind: "code",
			Mode: chain.SupervisedBlackBox,
			HandoffPolicy: &chain.HandoffPolicy{
				AllowedTargets:   []string{"target-one"},
				InputPolicy:      "structured",
				HandoffTimeoutMs: timeoutMs,
				ReturnActions:    []string{"complete", "abort"},
				HooksAfterReturn: []string{supervisedTestHookID},
			},
		}},
	}}}
}

type supervisedHarness struct {
	t       *testing.T
	runRoot string
	visible *adapters.StubVisibleClient
	now     time.Time
	ids     int
	options SupervisedAssemblyOptions
}

func newSupervisedHarness(t *testing.T, script []string, hookFail bool, timeoutMs int) *supervisedHarness {
	t.Helper()
	runRoot := filepath.Join(t.TempDir(), "run")
	if err := os.MkdirAll(runRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	harness := &supervisedHarness{
		t:       t,
		runRoot: runRoot,
		visible: &adapters.StubVisibleClient{Script: append([]string{}, script...)},
		now:     time.Date(2026, 10, 11, 9, 0, 0, 0, time.UTC),
	}
	results := &supervisedMemoryResults{}
	harness.options = SupervisedAssemblyOptions{
		RunID:         supervisedTestRunID,
		RunRoot:       runRoot,
		Definition:    supervisedTestDefinition(timeoutMs),
		VisibleSender: harness.visible,
		GateRunner: &gates.Runner{
			Executor: supervisedFakeExecutor{fail: hookFail},
			Scanner:  gates.DefaultSecretScanner(),
			Objects:  supervisedMemoryObjects{},
			Store:    results,
			Clock:    harness.clock,
		},
		GateHooks: map[string]gates.Hook{supervisedTestHookID: supervisedTestHook(supervisedTestHookID)},
		Clock:     harness.clock,
		IDs:       harness.nextID,
	}
	return harness
}

func (harness *supervisedHarness) clock() time.Time        { return harness.now }
func (harness *supervisedHarness) advance(d time.Duration) { harness.now = harness.now.Add(d) }

func (harness *supervisedHarness) nextID() string {
	harness.ids++
	return fmt.Sprintf("sv-%d", harness.ids)
}

func (harness *supervisedHarness) workspace() string {
	return filepath.Join(harness.runRoot, supervisedWorkspacesDirectory, supervisedTestTaskID+"-attempt-1")
}

func (harness *supervisedHarness) writeArtifact(name, content string) {
	harness.t.Helper()
	path := filepath.Join(harness.workspace(), name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		harness.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		harness.t.Fatal(err)
	}
}

// respond appends an operator answer to the binding interaction, exactly as
// `prfrail interactions respond --ledger ...` would.
func (harness *supervisedHarness) respond(selection string) {
	harness.t.Helper()
	ledger, err := loadOperatorInteractionLedger(SupervisedInteractionsLedgerPath(harness.runRoot))
	if err != nil {
		harness.t.Fatal(err)
	}
	pending := ledger.Pending()
	if len(pending) != 1 {
		harness.t.Fatalf("expected exactly one pending interaction, got %d", len(pending))
	}
	request := pending[0]
	harness.advance(time.Minute)
	response := chain.OperatorInteractionResponse{
		InteractionID:  request.Request.InteractionID,
		RequestID:      "sv-response",
		RespondedAt:    harness.now.Format(timestampLayout),
		RespondedBy:    request.Request.Operator,
		RunID:          request.Request.RunID,
		TaskID:         request.Request.TaskID,
		StepID:         request.Request.StepID,
		Attempt:        request.Request.Attempt,
		WorkspaceHash:  request.Request.WorkspaceHash,
		ConversationID: request.Request.ConversationID,
		ContextHash:    request.Request.ContextHash,
		RequestHash:    request.RecordHash,
		Selection:      selection,
		Evidence:       []string{request.RecordHash},
	}
	record, err := chain.NewOperatorInteractionResponseRecord(request, response, harness.clock, harness.nextID)
	if err != nil {
		harness.t.Fatalf("build operator response: %v", err)
	}
	if err := appendOperatorInteractionRecord(SupervisedInteractionsLedgerPath(harness.runRoot), record); err != nil {
		harness.t.Fatal(err)
	}
}

func (harness *supervisedHarness) beforeManifestPath() string {
	return filepath.Join(SupervisedDir(harness.runRoot), supervisedScansDirectoryName, supervisedBeforeManifestName)
}

// loadBeforeManifest reads the frozen before manifest the return-side diff is
// computed against.
func (harness *supervisedHarness) loadBeforeManifest() snapshot.SnapshotManifest {
	harness.t.Helper()
	var manifest snapshot.SnapshotManifest
	found, err := readSupervisedJSON(harness.beforeManifestPath(), &manifest)
	if err != nil {
		harness.t.Fatal(err)
	}
	if !found {
		harness.t.Fatal("the before manifest was not persisted")
	}
	return manifest
}

// loadSession reads the persisted handoff session of the open work window.
func (harness *supervisedHarness) loadSession() chain.HandoffSession {
	harness.t.Helper()
	var session chain.HandoffSession
	found, err := readSupervisedJSON(filepath.Join(SupervisedDir(harness.runRoot), supervisedSessionName), &session)
	if err != nil {
		harness.t.Fatal(err)
	}
	if !found {
		harness.t.Fatal("the handoff session was not persisted")
	}
	return session
}

// frozenPolicyPath is the store-local freeze of the handoff policy in force when
// the work window opened.
func (harness *supervisedHarness) frozenPolicyPath() string {
	return filepath.Join(SupervisedDir(harness.runRoot), supervisedFrozenPolicyName)
}

// loadScanReport reads the persisted post-return scan report, the scan evidence
// that feeds the frozen ScanEvidenceHash fact.
func (harness *supervisedHarness) loadScanReport() supervisedScanReport {
	harness.t.Helper()
	var report supervisedScanReport
	path := filepath.Join(SupervisedDir(harness.runRoot), supervisedScansDirectoryName, supervisedScanReportName)
	found, err := readSupervisedJSON(path, &report)
	if err != nil {
		harness.t.Fatal(err)
	}
	if !found {
		harness.t.Fatal("the post-return scan report was not persisted")
	}
	return report
}

func (harness *supervisedHarness) loadReport() chain.SupervisedReportRecord {
	harness.t.Helper()
	var record chain.SupervisedReportRecord
	found, err := readSupervisedJSON(SupervisedReportPath(harness.runRoot), &record)
	if err != nil {
		harness.t.Fatal(err)
	}
	if !found {
		harness.t.Fatal("the supervised report was not persisted")
	}
	if err := chain.ValidateSupervisedReport(record); err != nil {
		harness.t.Fatalf("persisted report does not validate: %v", err)
	}
	return record
}

func (harness *supervisedHarness) eventLog() string {
	harness.t.Helper()
	data, err := os.ReadFile(runEventLogPath(harness.runRoot))
	if err != nil {
		harness.t.Fatal(err)
	}
	return string(data)
}

// driveToWindow parks the run, answers the confirmation and opens the work
// window, leaving the harness one `supervised return` away from completion.
func (harness *supervisedHarness) driveToWindow() {
	harness.t.Helper()
	if _, err := ExecuteSupervisedRun(context.Background(), harness.options); err != nil {
		harness.t.Fatalf("parking the supervised run failed: %v", err)
	}
	harness.respond(supervisedConfirmationResponse)
	if _, err := ExecuteSupervisedRun(context.Background(), harness.options); err != nil {
		harness.t.Fatalf("opening the supervised work window failed: %v", err)
	}
	// The handoff close time must be strictly after the open time.
	harness.advance(time.Minute)
}

func hasStateName(counts []StateCount, state string) bool {
	for _, count := range counts {
		if count.State == state {
			return true
		}
	}
	return false
}

func containsName(values []string, name string) bool {
	for _, value := range values {
		if value == name {
			return true
		}
	}
	return false
}

// --- positive legs -----------------------------------------------------------

func TestSupervisedAssemblyFullChain(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()

	summary, err := ExecuteSupervisedRun(ctx, harness.options)
	if err != nil {
		t.Fatalf("first invocation must park, got %v", err)
	}
	if summary.ChainState != "PAUSED" {
		t.Fatalf("a parked supervised run must be PAUSED, got %s", summary.ChainState)
	}
	if hasStateName(summary.StepStateCounts, "PASSED") {
		t.Fatal("P2: a parked run must never have passed a step")
	}

	harness.respond(supervisedConfirmationResponse)
	summary, err = ExecuteSupervisedRun(ctx, harness.options)
	if err != nil {
		t.Fatalf("second invocation must open the window, got %v", err)
	}
	if summary.ChainState != "PAUSED" {
		t.Fatalf("the open window keeps the chain PAUSED, got %s", summary.ChainState)
	}
	if _, err := os.Stat(filepath.Join(SupervisedDir(harness.runRoot), supervisedSessionName)); err != nil {
		t.Fatalf("the handoff session must be persisted: %v", err)
	}

	harness.advance(time.Minute)
	harness.writeArtifact("artifact.txt", "hello from the supervised agent\n")
	summary, err = ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if err != nil {
		t.Fatalf("the explicit return must complete the run, got %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("the tail chain must complete the run, got %s", summary.ChainState)
	}
	if !hasStateName(summary.TaskStateCounts, "PASSED") {
		t.Fatalf("the task must pass: %+v", summary.TaskStateCounts)
	}

	report := harness.loadReport()
	if report.Record.AssuranceLevel != "reduced" {
		t.Fatalf("the assurance level is permanently reduced, got %q", report.Record.AssuranceLevel)
	}
	if report.Record.Mode != "supervised-black-box" {
		t.Fatalf("unexpected mode %q", report.Record.Mode)
	}
	for _, dimension := range []string{"cost", "background-processes", "external-side-effects"} {
		if !containsName(report.Record.UnknownDimensions, dimension) {
			t.Fatalf("the report must list %q as unknown: %v", dimension, report.Record.UnknownDimensions)
		}
	}
}

func TestSupervisedRunRestartRecovery(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.driveToWindow()

	// A restart rebuilds every projection from the durable event log and the
	// store-local supervised files; the window must still be re-entered.
	restarted := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	restarted.runRoot = harness.runRoot
	restarted.now = harness.now
	restarted.options.RunRoot = harness.runRoot
	restarted.options.Clock = restarted.clock
	restarted.options.IDs = restarted.nextID

	summary, err := ExecuteSupervisedRun(context.Background(), restarted.options)
	if err != nil {
		t.Fatalf("a restart during the window must not fail: %v", err)
	}
	if summary.ChainState != "PAUSED" {
		t.Fatalf("the restarted run must still wait for the return, got %s", summary.ChainState)
	}

	harness.writeArtifact("artifact.txt", "returned after the restart\n")
	summary, err = ExecuteSupervisedReturn(context.Background(), harness.options, "complete")
	if err != nil {
		t.Fatalf("the return must close after a restart, got %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("the restarted run must complete, got %s", summary.ChainState)
	}
}

func TestSupervisedReportRenderGolden(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.driveToWindow()
	harness.writeArtifact("artifact.txt", "rendered report\n")
	if _, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete"); err != nil {
		t.Fatalf("return failed: %v", err)
	}

	cli := NewCLI("test")
	var stdout, stderr strings.Builder
	if code := cli.Execute(context.Background(), []string{"supervised", "report", "--run-dir", harness.runRoot}, &stdout, &stderr); code != exitSuccess {
		t.Fatalf("supervised report failed (%d): %s", code, stderr.String())
	}
	output := stdout.String()
	if !strings.Contains(output, "unknown: cost=unknown") {
		t.Fatalf("P4: an unknown dimension must render literally as unknown:\n%s", output)
	}
	if !strings.Contains(output, "assuranceLevel: reduced") {
		t.Fatalf("the render must state the reduced assurance level:\n%s", output)
	}
	if strings.Contains(output, "independent-review") {
		t.Fatalf("P4: the report is frozen before the review and must not claim it:\n%s", output)
	}
	for _, forbidden := range []string{"safe", "verified"} {
		if strings.Contains(output, forbidden) {
			t.Fatalf("P4: the render must not carry the word %q:\n%s", forbidden, output)
		}
	}
}

// TestSupervisedReportOmitsUnprovenReviewDimension pins the F-3 fix: the
// reduced/unknown report is frozen before the tail chain (postflight -> review ->
// publish) runs, so it must not claim the independent-review dimension it has
// not earned yet. The tail chain completes on the same run and the persisted
// report still omits it.
func TestSupervisedReportOmitsUnprovenReviewDimension(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.driveToWindow()
	harness.writeArtifact("artifact.txt", "report written before the review\n")

	summary, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete")
	if err != nil {
		t.Fatalf("the return must complete, got %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("the tail chain must complete, got %s", summary.ChainState)
	}
	report := harness.loadReport()
	if containsName(report.Record.VerifiedDimensions, "independent-review") {
		t.Fatalf("P4: the frozen report must not claim the unproven review dimension: %v", report.Record.VerifiedDimensions)
	}
	for _, dimension := range []string{
		"workspace-isolation", "manifest-before", "manifest-after", "diff", "scope-scan",
		"secret-scan", "file-type-scan", "side-effect-scan", "process-stop-proof", "independent-gates",
	} {
		if !containsName(report.Record.VerifiedDimensions, dimension) {
			t.Fatalf("the write-time dimension %q must be recorded: %v", dimension, report.Record.VerifiedDimensions)
		}
	}
}

// --- negative legs -----------------------------------------------------------

func TestSupervisedAssemblyBlocksWithoutConfirmation(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("parking failed: %v", err)
	}
	summary, err := ExecuteSupervisedRun(ctx, harness.options)
	if err != nil {
		t.Fatalf("a parked run must stay in the waiting state, got %v", err)
	}
	if summary.ChainState != "PAUSED" || hasStateName(summary.StepStateCounts, "PASSED") {
		t.Fatalf("P2: without a confirmation the run cannot advance: %+v", summary)
	}
	if strings.Contains(harness.eventLog(), `"toState":"PASSED"`) {
		t.Fatal("P2: no PASSED transition may exist before the explicit return")
	}
	if _, err := os.Stat(filepath.Join(SupervisedDir(harness.runRoot), supervisedConfirmationsName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("P2: no confirmation record may exist without an operator answer")
	}
}

func TestSupervisedAssemblyRejectsSelectionReject(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("parking failed: %v", err)
	}
	harness.respond(supervisedRejectionResponse)
	summary, err := ExecuteSupervisedRun(ctx, harness.options)
	if !errors.Is(err, ErrSupervisedRejected) {
		t.Fatalf("a rejection must block the run, got %v", err)
	}
	if summary.ChainState != "PAUSED" {
		t.Fatalf("a rejected confirmation leaves the run paused, got %s", summary.ChainState)
	}
	if _, err := os.Stat(filepath.Join(SupervisedDir(harness.runRoot), supervisedConfirmationsName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("a rejection must never freeze a confirmation record")
	}
}

func TestSupervisedVisibleUnavailableFailsClosed(t *testing.T) {
	for _, scripted := range []string{adapters.VisibleStubUnavailable, adapters.VisibleStubUnknown} {
		t.Run(scripted, func(t *testing.T) {
			harness := newSupervisedHarness(t, []string{scripted}, false, 300000)
			summary, err := ExecuteSupervisedRun(context.Background(), harness.options)
			if !errors.Is(err, ErrSupervisedDeliveryBlocked) {
				t.Fatalf("a non-accepted delivery must fail closed, got %v", err)
			}
			if summary.ChainState == "COMPLETED" || hasStateName(summary.StepStateCounts, "PASSED") {
				t.Fatalf("no window and no pass may exist: %+v", summary)
			}
			if _, statErr := os.Stat(SupervisedInteractionsLedgerPath(harness.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatal("no binding interaction may be opened for an undelivered notice")
			}
		})
	}
}

// TestSupervisedConfirmationRefusesNonAcceptedPersistedDelivery pins the F-1
// console composition: even if the durable delivery record says the notice was
// not accepted (a record the capture path would never write), the confirmation
// phase must fail closed before it freezes a confirmation or resumes the
// operator. The run stays waiting and the chain work window never opens.
func TestSupervisedConfirmationRefusesNonAcceptedPersistedDelivery(t *testing.T) {
	for _, status := range []string{"unavailable", "unknown"} {
		t.Run(status, func(t *testing.T) {
			harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
			ctx := context.Background()
			if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
				t.Fatalf("parking failed: %v", err)
			}
			deliveryPath := filepath.Join(SupervisedDir(harness.runRoot), supervisedDeliveryName)
			var delivery chain.SupervisedDelivery
			if found, err := readSupervisedJSON(deliveryPath, &delivery); err != nil || !found {
				t.Fatalf("the accepted delivery record must be persisted: found=%v err=%v", found, err)
			}
			delivery.DeliveryStatus = status
			if err := writeSupervisedJSON(deliveryPath, delivery); err != nil {
				t.Fatal(err)
			}
			harness.respond(supervisedConfirmationResponse)
			summary, err := ExecuteSupervisedRun(ctx, harness.options)
			if !errors.Is(err, ErrSupervisedDeliveryBlocked) {
				t.Fatalf("a non-accepted persisted delivery must fail closed, got %v", err)
			}
			if summary.ChainState != "PAUSED" {
				t.Fatalf("the run must stay waiting, got %s", summary.ChainState)
			}
			if _, statErr := os.Stat(filepath.Join(SupervisedDir(harness.runRoot), supervisedConfirmationsName)); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatal("no confirmation record may be frozen for a non-accepted delivery")
			}
			if _, statErr := os.Stat(filepath.Join(SupervisedDir(harness.runRoot), supervisedSessionName)); !errors.Is(statErr, os.ErrNotExist) {
				t.Fatal("no session may be persisted for a non-accepted delivery")
			}
		})
	}
}

func TestSupervisedReturnRefusesFailedGate(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, true, 300000)
	harness.driveToWindow()
	harness.writeArtifact("artifact.txt", "gate will fail\n")

	summary, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete")
	if !errors.Is(err, ErrSupervisedGateFailed) {
		t.Fatalf("P3: a failed after-return gate must refuse the return, got %v", err)
	}
	if summary.ChainState == "COMPLETED" || hasStateName(summary.StepStateCounts, "PASSED") {
		t.Fatalf("P3: a confirmation must not exempt the gate: %+v", summary)
	}
	// The binding confirmation exists, which is exactly why the gate refusal
	// proves that a confirmation never exempts a safety gate.
	if _, statErr := os.Stat(filepath.Join(SupervisedDir(harness.runRoot), supervisedConfirmationsName)); statErr != nil {
		t.Fatal("the test requires a frozen confirmation record")
	}
}

func TestSupervisedReturnRefusesTimeout(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 1000)
	harness.driveToWindow()
	harness.writeArtifact("artifact.txt", "too late\n")
	harness.advance(2 * time.Second)

	if _, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete"); !errors.Is(err, ErrSupervisedTimeout) {
		t.Fatalf("an out-of-window complete return must be refused, got %v", err)
	}
	summary, err := ExecuteSupervisedReturn(context.Background(), harness.options, "abort")
	if errors.Is(err, ErrSupervisedTimeout) {
		t.Fatalf("abort must remain possible after a timeout, got %v", err)
	}
	if err != nil {
		t.Fatalf("abort must not fail: %v", err)
	}
	if summary.ChainState != "CANCELLED" {
		t.Fatalf("an aborted run is cancelled, got %s", summary.ChainState)
	}
}

// --- frozen handoff policy (blind review High) --------------------------------

// TestSupervisedReturnEnforcesFrozenPolicy is the positive leg of the
// return-entry drift guard: opening the window freezes the handoff policy as a
// store-local snapshot whose canonical digest is exactly the digest the session
// recorded, and an unchanged policy still converges end to end.
func TestSupervisedReturnEnforcesFrozenPolicy(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.driveToWindow()

	var frozen supervisedHandoffPolicySnapshot
	found, err := readSupervisedJSON(harness.frozenPolicyPath(), &frozen)
	if err != nil || !found {
		t.Fatalf("opening the window must freeze the handoff policy: found=%v err=%v", found, err)
	}
	if got := supervisedHandoffPolicyHash(frozen.Policy); got != frozen.PolicyHash {
		t.Fatalf("the snapshot must match its own digest: %s != %s", got, frozen.PolicyHash)
	}
	session := harness.loadSession()
	if frozen.PolicyHash != session.HandoffPolicyHash {
		t.Fatalf("the frozen snapshot must match the session hash: %s != %s", frozen.PolicyHash, session.HandoffPolicyHash)
	}
	inForce := harness.options.Definition.Tasks[0].Steps[0].HandoffPolicy
	if got := supervisedHandoffPolicyHash(*inForce); got != session.HandoffPolicyHash {
		t.Fatalf("the policy in force must match the frozen hash: %s != %s", got, session.HandoffPolicyHash)
	}

	harness.writeArtifact("artifact.txt", "frozen policy honoured\n")
	summary, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete")
	if err != nil {
		t.Fatalf("an unchanged policy must still converge, got %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("the tail chain must complete, got %s", summary.ChainState)
	}
}

// TestSupervisedReturnRefusesDriftedPolicy is the negative leg of the
// return-entry drift guard: a handoff policy that changed between window-open
// and return is refused fail-closed with no durable write. The change is applied
// to the definition the return reloads, which is exactly how a swapped chain
// configuration (or config path) reaches the assembly.
func TestSupervisedReturnRefusesDriftedPolicy(t *testing.T) {
	// The after-return-hooks leg registers the swapped hook in its own harness so
	// the sub-case proves the strong proposition: with the swapped hook runnable,
	// a run without the drift guard would execute the wrong gate set and converge,
	// instead of merely failing closed on a missing hook definition. The extra
	// registration is local to that sub-case's harness and touches no shared
	// fixture.
	cases := []struct {
		name    string
		mutate  func(policy *chain.HandoffPolicy)
		prepare func(harness *supervisedHarness)
	}{
		{name: "different timeout", mutate: func(policy *chain.HandoffPolicy) { policy.HandoffTimeoutMs = 600000 }},
		{
			name:   "different after-return hooks",
			mutate: func(policy *chain.HandoffPolicy) { policy.HooksAfterReturn = []string{"hook-two"} },
			prepare: func(harness *supervisedHarness) {
				harness.options.GateHooks["hook-two"] = supervisedTestHook("hook-two")
			},
		},
		{name: "different allowed targets", mutate: func(policy *chain.HandoffPolicy) { policy.AllowedTargets = []string{"target-two"} }},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
			harness.driveToWindow()
			harness.writeArtifact("artifact.txt", "policy drift\n")
			if testCase.prepare != nil {
				testCase.prepare(harness)
			}

			sessionBefore := harness.loadSession()
			eventsBefore := harness.eventLog()

			testCase.mutate(harness.options.Definition.Tasks[0].Steps[0].HandoffPolicy)

			summary, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete")
			if !errors.Is(err, ErrSupervisedPolicyDrift) {
				t.Fatalf("a drifted policy must refuse the complete return, got %v", err)
			}
			if summary.ChainState == "COMPLETED" || hasStateName(summary.StepStateCounts, "PASSED") {
				t.Fatalf("a drifted policy must not pass the step: %+v", summary)
			}
			// The refusal is a fail-closed stop, not a partial write: the window
			// artifacts that already existed stay untouched, while every artifact
			// the complete return would have frozen is absent.
			for _, path := range []string{
				harness.frozenPolicyPath(),
				filepath.Join(SupervisedDir(harness.runRoot), supervisedSessionName),
				filepath.Join(SupervisedDir(harness.runRoot), supervisedConfirmationsName),
			} {
				if _, statErr := os.Stat(path); statErr != nil {
					t.Fatalf("the refused return must leave the window artifact %s in place: %v", path, statErr)
				}
			}
			for _, path := range []string{
				filepath.Join(SupervisedDir(harness.runRoot), supervisedReceiptName),
				SupervisedReportPath(harness.runRoot),
				filepath.Join(SupervisedDir(harness.runRoot), supervisedScansDirectoryName, supervisedScanReportName),
				filepath.Join(SupervisedDir(harness.runRoot), supervisedScansDirectoryName, supervisedGateEvidenceName),
			} {
				if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
					t.Fatalf("a refused return must write no %s: %v", path, statErr)
				}
			}
			frozen := readFrozenPolicy(t, harness)
			if frozen.PolicyHash == supervisedHandoffPolicyHash(*harness.options.Definition.Tasks[0].Steps[0].HandoffPolicy) {
				t.Fatal("the frozen snapshot must keep the window's policy, not the swapped one")
			}
			if sessionAfter := harness.loadSession(); sessionAfter.HandoffPolicyHash != sessionBefore.HandoffPolicyHash {
				t.Fatal("the refused return must not rewrite the frozen window session")
			}
			if harness.eventLog() != eventsBefore {
				t.Fatal("a refused return must append no state transition")
			}
		})
	}
}

func readFrozenPolicy(t *testing.T, harness *supervisedHarness) supervisedHandoffPolicySnapshot {
	t.Helper()
	var frozen supervisedHandoffPolicySnapshot
	found, err := readSupervisedJSON(harness.frozenPolicyPath(), &frozen)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal("the frozen handoff policy snapshot was not persisted")
	}
	return frozen
}

// TestSupervisedReturnAllowsAbortUnderDriftedPolicy pins the deliberate
// exemption of the return-entry drift guard. The guard applies to complete
// returns only, and the exemption is load-bearing rather than an oversight:
// abort consumes no handoff policy at all. It builds its receipt from the
// already-persisted window session, writes no reduced/unknown report and runs
// neither the post-return scans nor the after-return gates, so there is nothing
// for a drifted policy to bend. That is exactly why abort must stay available
// when a run is stuck in drift or has outlived its frozen window: it is the
// operator's only remaining exit. A future refactor that routed abort through
// the drift guard would silently remove that exit, and this test would fail.
func TestSupervisedReturnAllowsAbortUnderDriftedPolicy(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.driveToWindow()

	// Drift the policy exactly as a swapped chain configuration would: the
	// definition the return reloads no longer matches the frozen window policy.
	// Both a field the guard checks (timeout) and the after-return hook set
	// change, so the drift is unambiguous.
	policy := harness.options.Definition.Tasks[0].Steps[0].HandoffPolicy
	policy.HandoffTimeoutMs = 600000
	policy.HooksAfterReturn = []string{"hook-two"}

	frozenHash := readFrozenPolicy(t, harness).PolicyHash
	if frozenHash == supervisedHandoffPolicyHash(*policy) {
		t.Fatal("the test requires the policy in force to differ from the frozen snapshot")
	}

	summary, err := ExecuteSupervisedReturn(context.Background(), harness.options, "abort")
	if err != nil {
		t.Fatalf("abort must be accepted under a drifted policy, got %v", err)
	}
	if summary.ChainState != "CANCELLED" {
		t.Fatalf("an aborted run must reach the terminal CANCELLED state, got %s", summary.ChainState)
	}
	if hasStateName(summary.StepStateCounts, "PASSED") {
		t.Fatalf("abort must not pass the step: %+v", summary.StepStateCounts)
	}

	// The abort wrote its own receipt (outcome abort) and none of the artifacts
	// the complete path owns: no reduced/unknown report, no scan evidence and no
	// gate evidence, hence no PASS and no review/promotion input.
	var receipt chain.HandoffReceiptRecord
	found, err := readSupervisedJSON(filepath.Join(SupervisedDir(harness.runRoot), supervisedReceiptName), &receipt)
	if err != nil || !found {
		t.Fatalf("abort must persist its handoff receipt: found=%v err=%v", found, err)
	}
	if receipt.Receipt.Outcome != "abort" {
		t.Fatalf("the abort receipt must record the abort outcome, got %q", receipt.Receipt.Outcome)
	}
	for _, artifact := range []string{
		SupervisedReportPath(harness.runRoot),
		filepath.Join(SupervisedDir(harness.runRoot), supervisedScansDirectoryName, supervisedScanReportName),
		filepath.Join(SupervisedDir(harness.runRoot), supervisedScansDirectoryName, supervisedGateEvidenceName),
	} {
		if _, statErr := os.Stat(artifact); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("abort must write no %s: %v", artifact, statErr)
		}
	}
}

// TestSupervisedFrozenPolicySnapshotIsCanonical pins the representation
// independence of the frozen-policy digest. The snapshot is written through the
// repository's canonical JSON encoder and read back through the strict decoder,
// so the digest depends on the policy content alone, never on JSON key order,
// whitespace or the file's on-disk formatting. The comparison anchor is
// supervisedHandoffPolicyHash -- the single convention that also populates
// HandoffSession.HandoffPolicyHash -- so the mode keeps exactly one policy
// digest and invents no second hashing scheme.
//
// HONEST LIMITATION: the digest is deliberately order-sensitive over the three
// list fields (AllowedTargets, ReturnActions, HooksAfterReturn), because that is
// the convention the session hash already uses. A reordered list is therefore
// read as drift and refused fail-closed rather than silently normalised; the
// last leg below pins that behaviour instead of claiming a canonicalisation the
// convention does not perform.
func TestSupervisedFrozenPolicySnapshotIsCanonical(t *testing.T) {
	policy := chain.HandoffPolicy{
		AllowedTargets:   []string{"target-one", "target-two"},
		InputPolicy:      "structured",
		HandoffTimeoutMs: 300000,
		ReturnActions:    []string{"complete", "abort"},
		HooksAfterReturn: []string{"hook-one"},
	}
	direct := supervisedHandoffPolicyHash(policy)

	path := filepath.Join(t.TempDir(), supervisedFrozenPolicyName)
	if err := writeSupervisedJSON(path, supervisedHandoffPolicySnapshot{Policy: policy, PolicyHash: direct}); err != nil {
		t.Fatal(err)
	}
	var roundTripped supervisedHandoffPolicySnapshot
	if found, err := readSupervisedJSON(path, &roundTripped); err != nil || !found {
		t.Fatalf("the snapshot must round-trip: found=%v err=%v", found, err)
	}
	if got := supervisedHandoffPolicyHash(roundTripped.Policy); got != direct {
		t.Fatalf("a canonical JSON round-trip must preserve the digest: %s != %s", got, direct)
	}

	// The same content in a different representation (leading whitespace and a
	// trailing blank line) must not move the digest: the strict reader
	// canonicalises before decoding.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	padded := filepath.Join(t.TempDir(), supervisedFrozenPolicyName)
	if err := os.WriteFile(padded, append([]byte("\n\t"), raw...), 0o600); err != nil {
		t.Fatal(err)
	}
	var paddedSnapshot supervisedHandoffPolicySnapshot
	if found, err := readSupervisedJSON(padded, &paddedSnapshot); err != nil || !found {
		t.Fatalf("a differently formatted snapshot must still read: found=%v err=%v", found, err)
	}
	if got := supervisedHandoffPolicyHash(paddedSnapshot.Policy); got != direct {
		t.Fatalf("representation-only changes must not move the digest: %s != %s", got, direct)
	}

	reordered := policy
	reordered.AllowedTargets = []string{"target-two", "target-one"}
	if supervisedHandoffPolicyHash(reordered) == direct {
		t.Fatal("the frozen convention is order-sensitive; a reordered list must be treated as drift")
	}
}

func TestSupervisedReturnRefusesSecretFinding(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.driveToWindow()
	harness.writeArtifact("leak.txt", "token=github_pat_11ABCDEFG0123456789abcdefghijklmnopqrstuvwxyz\n")

	if _, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete"); !errors.Is(err, ErrSupervisedSecretFinding) {
		t.Fatalf("a secret finding must refuse the return, got %v", err)
	}
}

func TestRunCommandRefusesMixedSupervisedAgentRunnerInputs(t *testing.T) {
	chainPath := writeSupervisedChainConfig(t)
	cli := NewCLI("test")
	var stdout, stderr strings.Builder
	code := cli.Execute(context.Background(), []string{"run", "--chain", chainPath, "--request-record", "declared.json"}, &stdout, &stderr)
	if code == exitSuccess {
		t.Fatal("P1: a supervised definition with AgentRunner inputs must be refused")
	}
	if !strings.Contains(stderr.String(), "mutually exclusive") {
		t.Fatalf("P1: the refusal must name the mutual exclusion: %s", stderr.String())
	}
}

func TestAgentRunnerAssemblyRefusesSupervisedDefinition(t *testing.T) {
	runRoot := t.TempDir()
	_, err := ExecuteAgentRunnerRun(context.Background(), AgentRunnerAssemblyOptions{
		RunID:      supervisedTestRunID,
		RunRoot:    runRoot,
		Definition: supervisedTestDefinition(300000),
	})
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("P1: the AgentRunner assembly must refuse a supervised definition, got %v", err)
	}
	entries, readErr := os.ReadDir(runRoot)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("P1: a refused AgentRunner assembly must write nothing: %v", entries)
	}
}

func TestSupervisedAssemblyRefusesAgentRunnerDefinition(t *testing.T) {
	definition := chain.Definition{ID: "chain-isolated", Tasks: []chain.Task{{
		ID:    supervisedTestTaskID,
		Steps: []chain.Step{{ID: supervisedTestStepID, Kind: "code", Mode: chain.IsolatedWorkspace}},
	}}}
	_, err := ExecuteSupervisedRun(context.Background(), SupervisedAssemblyOptions{
		RunID:      supervisedTestRunID,
		RunRoot:    t.TempDir(),
		Definition: definition,
	})
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("the supervised assembly must refuse an AgentRunner definition, got %v", err)
	}
}

// --- GAP-1: scope scan -------------------------------------------------------

// TestSupervisedReturnRefusesOutOfScopeChange proves that a post-return change
// whose path leaves the authorized scope blocks the return. The scope is derived
// from the frozen handoff policy (its declared allowed target), so a change set
// that reaches outside the isolated workspace is refused before any receipt.
func TestSupervisedReturnRefusesOutOfScopeChange(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.driveToWindow()

	// The frozen before manifest is evidence, and evidence is untrusted input to
	// the return: inject a changed path that resolves outside the workspace root.
	before := harness.loadBeforeManifest()
	before.Manifest.Entries = append(before.Manifest.Entries, snapshot.Entry{
		Path:        "../outside.txt",
		Type:        snapshot.EntryRegularFile,
		ContentHash: evidence.Digest("", []byte("outside the scope")),
		SizeBytes:   int64(len("outside the scope")),
	})
	if err := writeSupervisedJSON(harness.beforeManifestPath(), before); err != nil {
		t.Fatal(err)
	}
	harness.writeArtifact("artifact.txt", "in scope\n")

	summary, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete")
	if !errors.Is(err, ErrSupervisedScopeOutOfBounds) {
		t.Fatalf("an out-of-scope change must refuse the return, got %v", err)
	}
	if summary.ChainState == "COMPLETED" || hasStateName(summary.StepStateCounts, "PASSED") {
		t.Fatalf("an out-of-scope change must not pass the step: %+v", summary)
	}
	if _, statErr := os.Stat(SupervisedReportPath(harness.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatal("no complete receipt or report may be frozen for an out-of-scope return")
	}
}

// TestSupervisedScopeUnderivableFailsClosed proves the fail-closed half of the
// scope guard: an underivable scope is refused rather than read as "no
// restriction". HONEST LIMITATION: through the production entry a definition
// with no allowed target is refused even earlier by chain.New
// (ErrInvalidDefinition), so this test asserts both the guard itself and the
// assembly-level refusal, which is stricter than the guard, not weaker.
func TestSupervisedScopeUnderivableFailsClosed(t *testing.T) {
	if _, err := newSupervisedScope(chain.HandoffPolicy{}, t.TempDir()); !errors.Is(err, ErrSupervisedScopeUnderivable) {
		t.Fatalf("an empty allowed-target allow-list must be underivable, got %v", err)
	}

	definition := supervisedTestDefinition(300000)
	definition.Tasks[0].Steps[0].HandoffPolicy.AllowedTargets = nil
	_, err := ExecuteSupervisedRun(context.Background(), SupervisedAssemblyOptions{
		RunID:         supervisedTestRunID,
		RunRoot:       t.TempDir(),
		Definition:    definition,
		VisibleSender: &adapters.StubVisibleClient{Script: []string{adapters.VisibleStubAccepted}},
		Clock:         func() time.Time { return time.Date(2026, 10, 11, 9, 0, 0, 0, time.UTC) },
	})
	if !errors.Is(err, chain.ErrInvalidDefinition) {
		t.Fatalf("an underivable scope must fail closed, got %v", err)
	}
}

// TestSupervisedScopeRejectsUnresolvableRoot pins the EvalSymlinks failure
// branch of newSupervisedScope (runtime_supervised.go ~864): a workspace root
// that cannot be resolved -- a non-existent path fails EvalSymlinks
// deterministically on every supported platform -- is an underivable scope and
// fails closed. The empty allow-list branch is pinned separately by
// TestSupervisedScopeUnderivableFailsClosed.
func TestSupervisedScopeRejectsUnresolvableRoot(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "not-materialized")
	if _, err := newSupervisedScope(chain.HandoffPolicy{AllowedTargets: []string{"target-one"}}, missing); !errors.Is(err, ErrSupervisedScopeUnderivable) {
		t.Fatalf("an unresolvable workspace root must fail closed, got %v", err)
	}
}

// TestSupervisedScopeAllowsRejectsLinkResolvingOutsideRoot pins the
// symlink-resolution escape branch of allows (runtime_supervised.go ~887): a
// lexically contained relative path that resolves through a link to a location
// outside the authorized root is rejected even though its lexical position is
// inside the root. HONEST LIMITATION (reachability): the guard is
// defence-in-depth over an untrusted manifest (snapshot.Capture only yields
// workspace-relative paths), and this leg additionally needs a real link on
// disk, so it skips on hosts that forbid symlink creation (the Ubuntu CI leg
// runs it); the lexical half of the same guard is pinned on every platform by
// TestSupervisedSideEffectScanBlocksTamperedManifests.
func TestSupervisedScopeAllowsRejectsLinkResolvingOutsideRoot(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	target := filepath.Join(outside, "target.txt")
	if err := os.WriteFile(target, []byte("outside the root\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(root, "escape.txt")); err != nil {
		t.Skipf("symlink creation unavailable: %v", err)
	}
	scope, err := newSupervisedScope(chain.HandoffPolicy{AllowedTargets: []string{"target-one"}}, root)
	if err != nil {
		t.Fatalf("the scope must derive for an existing root, got %v", err)
	}
	if err := scope.allows("escape.txt"); !errors.Is(err, ErrSupervisedScopeOutOfBounds) {
		t.Fatalf("a link resolving outside the root must be out of bounds, got %v", err)
	}
}

// --- GAP-1: side-effect scan -------------------------------------------------

// TestSupervisedReturnRefusesSideEffect proves that an observable side effect on
// the isolated workspace (here: the deletion of a pre-existing entry) blocks the
// return.
func TestSupervisedReturnRefusesSideEffect(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()
	// Seed a tracked entry before the run so the frozen before manifest carries
	// it, then delete it during the supervised window.
	harness.writeArtifact("seed.txt", "tracked before the window\n")
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("parking the supervised run failed: %v", err)
	}
	harness.respond(supervisedConfirmationResponse)
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("opening the supervised work window failed: %v", err)
	}
	harness.advance(time.Minute)
	if err := os.Remove(filepath.Join(harness.workspace(), "seed.txt")); err != nil {
		t.Fatal(err)
	}

	summary, err := ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if !errors.Is(err, ErrSupervisedSideEffect) {
		t.Fatalf("a detected side effect must refuse the return, got %v", err)
	}
	if summary.ChainState == "COMPLETED" {
		t.Fatalf("a detected side effect must not complete the run: %+v", summary)
	}
}

// TestSupervisedSideEffectScanBlocksTamperedManifests pins the four blocking
// branches of runSideEffectScan (runtime_supervised.go ~921) that the deletion
// leg does not cover: a type change, a permission/executable change, an
// absolute symbolic-link target, and a relative symbolic-link target that
// escapes the workspace root. HONEST LIMITATION (reachability): the scan is
// defence-in-depth over an untrusted manifest -- snapshot.Capture only emits
// workspace-relative, lexically contained paths, so a "../" style target can
// only arrive through a tampered before/after manifest, which is exactly how
// these rows inject it. The positive control (an unmodified pair plus an
// in-bounds relative link) proves the rows are load-bearing rather than "any
// error at all" being accepted.
func TestSupervisedSideEffectScanBlocksTamperedManifests(t *testing.T) {
	regular := func(path string, executable bool) snapshot.Entry {
		return snapshot.Entry{Path: path, Type: snapshot.EntryRegularFile, Executable: executable}
	}
	cases := []struct {
		name   string
		before []snapshot.Entry
		after  []snapshot.Entry
		want   error
	}{
		{
			name:   "positive control: unmodified pair passes",
			before: []snapshot.Entry{regular("kept.txt", false)},
			after: []snapshot.Entry{
				regular("kept.txt", false),
				regular("added.txt", false),
				{Path: "dir/link", Type: snapshot.EntrySymbolicLink, Target: "kept.txt"},
			},
			want: nil,
		},
		{
			name:   "type change",
			before: []snapshot.Entry{regular("entry.txt", false)},
			after:  []snapshot.Entry{{Path: "entry.txt", Type: snapshot.EntryDirectory}},
			want:   ErrSupervisedSideEffect,
		},
		{
			name:   "permission change",
			before: []snapshot.Entry{regular("script.sh", false)},
			after:  []snapshot.Entry{regular("script.sh", true)},
			want:   ErrSupervisedSideEffect,
		},
		{
			name:  "absolute link target",
			after: []snapshot.Entry{{Path: "link", Type: snapshot.EntrySymbolicLink, Target: "/etc/passwd"}},
			want:  ErrSupervisedSideEffect,
		},
		{
			name:  "relative link escaping the workspace root",
			after: []snapshot.Entry{{Path: "link", Type: snapshot.EntrySymbolicLink, Target: "../../outside"}},
			want:  ErrSupervisedSideEffect,
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			runtime := &supervisedRuntime{
				workspaceRoot: t.TempDir(),
				policy:        chain.HandoffPolicy{AllowedTargets: []string{"target-one"}},
			}
			err := runtime.runSideEffectScan(testCase.before, testCase.after)
			if testCase.want == nil {
				if err != nil {
					t.Fatalf("an unmodified manifest pair must pass, got %v", err)
				}
				return
			}
			if !errors.Is(err, testCase.want) {
				t.Fatalf("%s must block the return, got %v", testCase.name, err)
			}
		})
	}
}

// --- GAP-3: file-type scan ---------------------------------------------------

// TestSupervisedFileTypeScanRejectsUnacceptedType proves that an unaccepted
// entry type blocks the post-return scan membership. The violation is injected at
// the manifest boundary because snapshot.Capture, by construction, can only emit
// the three accepted types; the guard is therefore defense in depth over an
// untrusted manifest, and runSecretAndTypeScan is the exact function the return
// path calls.
func TestSupervisedFileTypeScanRejectsUnacceptedType(t *testing.T) {
	runtime := &supervisedRuntime{workspaceRoot: t.TempDir()}
	manifest := snapshot.SnapshotManifest{Manifest: snapshot.ManifestBody{Entries: []snapshot.Entry{
		{Path: "pipe", Type: snapshot.EntryType("named-pipe")},
	}}}
	err := runtime.runSecretAndTypeScan(context.Background(), manifest)
	if !errors.Is(err, ErrSupervisedFileType) {
		t.Fatalf("an unaccepted entry type must block the scan, got %v", err)
	}
}

// --- GAP-4: unknown / live managed process -----------------------------------

// TestSupervisedReturnRefusesLiveManagedProcess proves that a live managed
// process identity at return time fails the return closed. "Unknown process"
// cannot be positively observed (see processStopEvidence): only the recorded
// managed identity is inspectable, and the unobservable external agent tree is
// carried by the report's background-processes unknown dimension.
func TestSupervisedReturnRefusesLiveManagedProcess(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.driveToWindow()
	harness.writeArtifact("artifact.txt", "written by a live process\n")

	identity, err := guard.InspectProcess(os.Getpid())
	if err != nil {
		t.Fatalf("the test process identity must be inspectable: %v", err)
	}
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("the test process must be observed alive: alive=%v err=%v", alive, err)
	}
	identityPath := filepath.Join(harness.runRoot, managedProcessIdentityFileName)
	if err := writeSupervisedJSON(identityPath, identity); err != nil {
		t.Fatal(err)
	}

	summary, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete")
	if !errors.Is(err, ErrSupervisedGateFailed) {
		t.Fatalf("a live managed process must refuse the return, got %v", err)
	}
	if summary.ChainState == "COMPLETED" || hasStateName(summary.StepStateCounts, "PASSED") {
		t.Fatalf("a live managed process must not pass the step: %+v", summary)
	}
}

// --- GAP-1a: isolation coverage ----------------------------------------------

// TestSupervisedIsolationReadsOnlyWorkspace asserts (a) that the black-box work
// and the acceptance/promotion observe only the isolated workspace root under
// <runRoot>/workspaces/..., and (b) that a candidate placed outside that
// workspace is neither read nor accepted nor promoted.
func TestSupervisedIsolationReadsOnlyWorkspace(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()

	// A candidate outside the isolated workspace must be invisible to the run.
	outsidePath := filepath.Join(harness.runRoot, "outside-candidate.txt")
	const outsideContent = "candidate outside the isolated workspace\n"
	if err := os.WriteFile(outsidePath, []byte(outsideContent), 0o600); err != nil {
		t.Fatal(err)
	}

	harness.driveToWindow()
	want := filepath.Join(harness.runRoot, supervisedWorkspacesDirectory, supervisedTestTaskID+"-attempt-1")
	if got := harness.workspace(); got != want {
		t.Fatalf("the black-box workspace must be the isolated root under workspaces/: %s", got)
	}
	harness.writeArtifact("artifact.txt", "inside the isolated workspace\n")

	summary, err := ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if err != nil {
		t.Fatalf("the return must complete against the isolated workspace, got %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("the tail chain must complete, got %s", summary.ChainState)
	}

	report := harness.loadScanReport()
	if len(report.ChangedPaths) != 1 || report.ChangedPaths[0] != "added artifact.txt" {
		t.Fatalf("the manifest/diff read must cover only the isolated workspace: %v", report.ChangedPaths)
	}
	if report.EntryCount != 1 {
		t.Fatalf("only the in-workspace artifact may be captured, got entry count %d", report.EntryCount)
	}
	data, err := os.ReadFile(outsidePath)
	if err != nil {
		t.Fatalf("the outside candidate must not be consumed: %v", err)
	}
	if string(data) != outsideContent {
		t.Fatal("the outside candidate was modified")
	}
}

// --- extended scan set positive ----------------------------------------------

// TestSupervisedAssemblyFullChainExtendedScans is the sibling of
// TestSupervisedAssemblyFullChain: it proves that the extended scan set keeps the
// full chain completing and that the reduced report records the new scans as
// proved while still carrying the unobservable dimensions as unknown.
func TestSupervisedAssemblyFullChainExtendedScans(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.driveToWindow()
	harness.writeArtifact("artifact.txt", "extended scan set\n")

	summary, err := ExecuteSupervisedReturn(context.Background(), harness.options, "complete")
	if err != nil {
		t.Fatalf("the extended scan set must complete the run, got %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("the tail chain must complete, got %s", summary.ChainState)
	}

	report := harness.loadReport()
	if report.Record.AssuranceLevel != "reduced" {
		t.Fatalf("the assurance level is permanently reduced, got %q", report.Record.AssuranceLevel)
	}
	for _, dimension := range []string{"scope-scan", "file-type-scan", "side-effect-scan"} {
		if !containsName(report.Record.VerifiedDimensions, dimension) {
			t.Fatalf("the report must record %q as proved: %v", dimension, report.Record.VerifiedDimensions)
		}
	}
	for _, dimension := range []string{"cost", "background-processes", "external-side-effects"} {
		if !containsName(report.Record.UnknownDimensions, dimension) {
			t.Fatalf("the report must list %q as unknown: %v", dimension, report.Record.UnknownDimensions)
		}
	}
	scan := harness.loadScanReport()
	if !scan.ChangeSetInScope || !scan.SideEffectClean || !scan.FileTypeScanClean || !scan.SecretScanClean {
		t.Fatalf("every scan must be recorded clean after a passing return: %+v", scan)
	}
}

// --- F-2: confirmation-phase crash recovery ---------------------------------

// errSupervisedCrashInjected marks a crash injected at a durable boundary of the
// confirmation sequence.
var errSupervisedCrashInjected = errors.New("injected supervised confirmation crash")

// crashAfter arms the harness so the next invocation fails at one durable
// boundary, exactly as a process death there would.
func (harness *supervisedHarness) crashAfter(point string) {
	harness.t.Helper()
	harness.options.CrashPoint = func(candidate string) error {
		if candidate == point {
			return errSupervisedCrashInjected
		}
		return nil
	}
}

// clearCrash disarms the crash seam.
func (harness *supervisedHarness) clearCrash() { harness.options.CrashPoint = nil }

// requireConvergedWindow asserts that a re-entry replays the durable
// confirmation into an open, persisted window and that the run still completes.
func requireConvergedWindow(t *testing.T, harness *supervisedHarness, label string) {
	t.Helper()
	ctx := context.Background()
	summary, err := ExecuteSupervisedRun(ctx, harness.options)
	if err != nil {
		t.Fatalf("%s: re-entry must converge, got %v", label, err)
	}
	if summary.ChainState != "PAUSED" {
		t.Fatalf("%s: the converged run waits for the return, got %s", label, summary.ChainState)
	}
	if _, err := os.Stat(filepath.Join(SupervisedDir(harness.runRoot), supervisedSessionName)); err != nil {
		t.Fatalf("%s: the replay must persist the handoff session: %v", label, err)
	}
	harness.advance(time.Minute)
	harness.writeArtifact("artifact.txt", label+"\n")
	summary, err = ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if err != nil {
		t.Fatalf("%s: the run must complete after the replay, got %v", label, err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("%s: the tail chain must complete, got %s", label, summary.ChainState)
	}
}

// TestSupervisedConfirmationCrashAfterAppendConverges injects a crash after the
// confirmation became durable but before the operator resume: the parked window
// is replayed through the confirmation, the resume and the window opener.
func TestSupervisedConfirmationCrashAfterAppendConverges(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("parking failed: %v", err)
	}
	harness.respond(supervisedConfirmationResponse)
	harness.crashAfter(supervisedCrashAfterConfirmation)
	if _, err := ExecuteSupervisedRun(ctx, harness.options); !errors.Is(err, errSupervisedCrashInjected) {
		t.Fatalf("the injected crash must abort the confirmation sequence, got %v", err)
	}
	harness.clearCrash()
	if _, err := os.Stat(filepath.Join(SupervisedDir(harness.runRoot), supervisedSessionName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the crash must leave no session behind: %v", err)
	}
	requireConvergedWindow(t, harness, "after-confirmation crash")
}

// TestSupervisedConfirmationCrashAfterResumeConverges injects a crash after the
// operator resume but before the window opens: the re-entry resumes idempotently
// and opens the window.
func TestSupervisedConfirmationCrashAfterResumeConverges(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("parking failed: %v", err)
	}
	harness.respond(supervisedConfirmationResponse)
	harness.crashAfter(supervisedCrashAfterResume)
	if _, err := ExecuteSupervisedRun(ctx, harness.options); !errors.Is(err, errSupervisedCrashInjected) {
		t.Fatalf("the injected crash must abort the confirmation sequence, got %v", err)
	}
	harness.clearCrash()
	requireConvergedWindow(t, harness, "after-resume crash")
}

// TestSupervisedConfirmationCrashAfterWindowOpenConverges injects a crash after
// the window opened but before the session was persisted: the window is already
// durable, so the re-entry recovers the open time and persists the session
// instead of stalling on the non-empty confirmation ledger.
func TestSupervisedConfirmationCrashAfterWindowOpenConverges(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("parking failed: %v", err)
	}
	harness.respond(supervisedConfirmationResponse)
	harness.crashAfter(supervisedCrashAfterWindowOpen)
	if _, err := ExecuteSupervisedRun(ctx, harness.options); !errors.Is(err, errSupervisedCrashInjected) {
		t.Fatalf("the injected crash must abort the confirmation sequence, got %v", err)
	}
	harness.clearCrash()
	requireConvergedWindow(t, harness, "after-window-open crash")
}

// TestNoSupervisedReferencesInAgentRunnerPaths is a source-level invariant. It
// proves that the AgentRunner edge of the adapters layer and the AgentRunner
// console assembly carry no supervised token, and that the visible adapter
// references no chain transition. HONEST LIMITATION: this is a lexical check
// over the current source text; an equivalent rewrite (for example building the
// token at run time) could evade it. It is a guard against accidental coupling,
// not a proof of impossibility.
func TestNoSupervisedReferencesInAgentRunnerPaths(t *testing.T) {
	tokens := []string{"supervised-black-box", "SupervisedBlackBox"}
	adaptersDir := filepath.Join("..", "adapters")
	entries, err := os.ReadDir(adaptersDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		assertNoToken(t, filepath.Join(adaptersDir, entry.Name()), tokens)
	}
	assertNoToken(t, "runtime_agent.go", tokens)

	visible, err := os.ReadFile(filepath.Join(adaptersDir, "sessbridge_visible.go"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(visible), "chain.") {
		t.Fatal("the visible adapter must reference no chain transition")
	}
}

// TestSupervisedPathNeverRoutesThroughAgentRunner is the semantic sibling of
// TestNoSupervisedReferencesInAgentRunnerPaths, which scans source text and could
// be evaded by an equivalent rewrite. This test never inspects source text: it
// asserts on the assembled object graph and on the durable state a real run
// produced.
//
//  1. The production assembly (newSupervisedEngine, the exact constructor both
//     production entries call) wires no AgentRunner step port and no AgentRunner
//     admission port, and does wire its supervised step port. Because the engine
//     routes a step through StepRouter.Execute, a nil AgentRunner port makes the
//     AgentRunner runner entry unreachable by construction, not merely
//     unreferenced.
//  2. A full supervised run -- park, confirm, open, return, tail chain -- never
//     records the AgentRunner-only TERMINAL_PENDING step state, never records an
//     agent-dispatched reason, and never materialises the AgentRunner replay
//     store directory where a dispatched step's request/completion records would
//     live.
func TestSupervisedPathNeverRoutesThroughAgentRunner(t *testing.T) {
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	ctx := context.Background()

	taskID, stepID, policy, err := supervisedAssemblyStep(harness.options.Definition)
	if err != nil {
		t.Fatalf("the supervised definition must resolve: %v", err)
	}
	assemblyOptions := harness.options
	assemblyOptions.RunRoot = t.TempDir()
	_, runtime, err := newSupervisedEngine(ctx, assemblyOptions, taskID, stepID, *policy, harness.clock)
	if err != nil {
		t.Fatalf("building the production supervised assembly failed: %v", err)
	}
	if runtime.steps.AgentRunner != nil {
		t.Fatal("the supervised assembly must wire no AgentRunner step port")
	}
	if runtime.steps.AgentRunnerAdmission != nil {
		t.Fatal("the supervised assembly must wire no AgentRunner admission port")
	}
	if runtime.steps.Supervised == nil {
		t.Fatal("the supervised assembly must wire its supervised step port")
	}

	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("parking the supervised run failed: %v", err)
	}
	harness.respond(supervisedConfirmationResponse)
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("opening the supervised work window failed: %v", err)
	}
	harness.advance(time.Minute)
	harness.writeArtifact("artifact.txt", "supervised only\n")
	summary, err := ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if err != nil {
		t.Fatalf("the supervised return must complete, got %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("the tail chain must complete, got %s", summary.ChainState)
	}

	// Positive control: the supervised leg really ran, so the absence asserted
	// below is not the absence of an execution that never happened.
	if harness.visible.Calls() == 0 {
		t.Fatal("the supervised visible delivery must have run")
	}

	store, err := chain.NewFileEventStore(runEventLogPath(harness.runRoot))
	if err != nil {
		t.Fatal(err)
	}
	events, err := store.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 {
		t.Fatal("the completed run must have persisted state events")
	}
	for _, record := range events {
		if record.Event.ToState == "TERMINAL_PENDING" {
			t.Fatalf("the supervised path must never park in the AgentRunner-only state: %+v", record.Event)
		}
		if record.Event.Reason.Code == "agent-dispatched" {
			t.Fatalf("the supervised path must never record an AgentRunner dispatch: %+v", record.Event)
		}
	}
	if hasStateName(summary.StepStateCounts, "TERMINAL_PENDING") {
		t.Fatalf("no supervised step may end in TERMINAL_PENDING: %+v", summary.StepStateCounts)
	}
	// adapters.agentRunnerReplayRootDirectoryName: the store that holds a
	// dispatched AgentRunner step's request/completion records.
	replayRoot := filepath.Join(harness.runRoot, "agent-runner-replay")
	if _, statErr := os.Stat(replayRoot); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the supervised path must create no AgentRunner replay store: %v", statErr)
	}
}

func assertNoToken(t *testing.T, path string, tokens []string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	for _, token := range tokens {
		if strings.Contains(string(data), token) {
			t.Fatalf("%s must not reference %q", path, token)
		}
	}
}

// writeSupervisedChainConfig writes a valid supervised chain configuration and
// returns its path.
func writeSupervisedChainConfig(t *testing.T) string {
	t.Helper()
	config := DefaultChainConfig("supervised-config")
	config.AI = &AIConfig{
		Profiles: []adapters.AIProviderProfile{{
			ProfileID:    "local-visible",
			ProviderType: "github-copilot",
			Model:        "auto",
			AuthMode:     "host",
		}},
		Channels: []AIChannelBinding{{Channel: supervisedRequiredChannel, ProfileID: "local-visible"}},
	}
	config.Tasks = []TaskConfig{{
		ID:     supervisedTestTaskID,
		Review: "policy",
		Steps: []StepConfig{{
			ID:        supervisedTestStepID,
			Kind:      "code",
			Execution: supervisedExecutionMode,
			HandoffPolicy: &HandoffPolicyInput{
				AllowedTargets:   []string{"target-one"},
				InputPolicy:      "structured",
				HandoffTimeoutMs: 300000,
				ReturnActions:    []string{"complete", "abort"},
				HooksAfterReturn: []string{supervisedTestHookID},
			},
		}},
	}}
	if err := ValidateChainConfig(config); err != nil {
		t.Fatalf("the supervised test configuration must validate: %v", err)
	}
	content, err := EncodeChainConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), DefaultChainConfigName)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
