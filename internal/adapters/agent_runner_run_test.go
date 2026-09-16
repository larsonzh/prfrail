package adapters

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/tickets"
)

func newAgentRunnerRunFixture(t *testing.T, mode string) (*AgentRunnerPinnedCLIRun, *AgentRunnerReplayStore, string) {
	t.Helper()
	return newAgentRunnerRunFixtureArgs(t, []string{mode})
}

func newAgentRunnerRunFixtureArgs(t *testing.T, args []string) (*AgentRunnerPinnedCLIRun, *AgentRunnerReplayStore, string) {
	t.Helper()
	runRoot := t.TempDir()
	workspace := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	config := agentRunnerPinnedCLIConfig("")
	helper := []string{"-test.run=^TestAgentRunnerPinnedCLIHelper$", "--"}
	config.RunArgs = append(append([]string{}, helper...), args...)
	launcher, err := NewAgentRunnerPinnedCLILauncher(config, store)
	if err != nil {
		t.Fatal(err)
	}
	return &AgentRunnerPinnedCLIRun{Launcher: launcher, Store: store}, store, workspace
}

func agentRunnerRunRequest(workspace string) AgentRunnerRunRequest {
	return AgentRunnerRunRequest{
		Launch:      agentRunnerPinnedCLIRequest(workspace),
		RequestHash: evidence.Digest("", []byte("request-record")),
		TaskID:      "task-two",
		StepID:      "code-isolated",
		Attempt:     1,
		Timeout:     15 * time.Second,
		MaxLogBytes: 1 << 20,
		Capturer:    &fakeManifestCapturer{before: []byte("before"), after: []byte("after"), diff: []byte("diff")},
	}
}

// requireValidCompletion proves the mapped outcome is itself a legal receipt,
// which is what the caller hands to the publisher.
func requireValidCompletion(t *testing.T, outcome AgentRunnerTerminalOutcome) {
	t.Helper()
	if _, err := NewAgentRunnerCompletionRecord(outcome.Completion); err != nil {
		t.Fatalf("the mapped completion must validate: %v", err)
	}
}

func TestAgentRunnerPinnedCLIRunCompletedOutcomeCarriesFacts(t *testing.T) {
	run, _, workspace := newAgentRunnerRunFixture(t, "exit")
	outcome, err := run.Run(context.Background(), agentRunnerRunRequest(workspace))
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Completion.Status != "completed" || outcome.Completion.ProcessTreeStatus != "stopped" {
		t.Fatalf("a clean exit must be completed with a stopped tree: %+v", outcome.Completion)
	}
	if outcome.Facts == nil {
		t.Fatal("a completed outcome must carry the frozen facts")
	}
	if err := outcome.Facts.Validate(); err != nil {
		t.Fatalf("the collected facts must validate: %v", err)
	}
	if outcome.Completion.OutputManifestHash == nil || *outcome.Completion.OutputManifestHash != outcome.Facts.ManifestHash {
		t.Fatalf("the completion must publish the captured manifest: %+v", outcome.Completion)
	}
	if outcome.Settlement.Status != tickets.CostSettlementObserved || outcome.Settlement.ObservedCalls == nil || *outcome.Settlement.ObservedCalls != 2 {
		t.Fatalf("a completed run settles the observed usage: %+v", outcome.Settlement)
	}
	requireValidCompletion(t, outcome)
}

func TestAgentRunnerPinnedCLIRunTimeoutIsUncertainWithAProvenStop(t *testing.T) {
	run, _, workspace := newAgentRunnerRunFixture(t, "run")
	request := agentRunnerRunRequest(workspace)
	request.Timeout = 300 * time.Millisecond
	outcome, err := run.Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Completion.Status != "uncertain" {
		t.Fatalf("a timed-out run must not be completed: %+v", outcome.Completion)
	}
	if outcome.Completion.ProcessTreeStatus != "stopped" {
		t.Fatalf("a timeout must still stop the tree: %+v", outcome.Completion)
	}
	if outcome.Facts != nil {
		t.Fatal("an uncertain outcome must carry no frozen facts")
	}
	if outcome.Settlement.Status != tickets.CostSettlementUnknown {
		t.Fatalf("an uncertain outcome settles as unknown: %+v", outcome.Settlement)
	}
	if len(outcome.Completion.ErrorEvidence) == 0 {
		t.Fatal("a non-completed receipt must carry error evidence")
	}
	// A stopped run did not exit on its own, so it has no exit status: reporting the
	// platform's -1 (or a stale zero) would be a fabricated fact.
	if outcome.Completion.ExitCode != nil {
		t.Fatalf("a stopped run must not report an exit code: %v", *outcome.Completion.ExitCode)
	}
	// The receipt must cite the stop evidence it archived, not just describe it.
	requireValidCompletion(t, outcome)
}

func TestAgentRunnerPinnedCLIRunCallerDeadlineIsCancelled(t *testing.T) {
	run, _, workspace := newAgentRunnerRunFixture(t, "run")
	request := agentRunnerRunRequest(workspace)
	// The caller's own deadline expires long before the run's timeout: that is a
	// cancellation, and the run may never be reported as completed just because the
	// killed CLI had already written every artifact. The deadline is comfortably after
	// the launch is registered, so the precheck cannot be the thing that fails.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	request.Timeout = 30 * time.Second
	outcome, err := run.Run(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Completion.Status != "cancelled" {
		t.Fatalf("a caller deadline must cancel the run, got %+v", outcome.Completion)
	}
	if outcome.Facts != nil {
		t.Fatal("a cancelled run must carry no frozen facts")
	}
	if outcome.Completion.ExitCode != nil {
		t.Fatalf("a cancelled run must not report an exit code: %v", *outcome.Completion.ExitCode)
	}
	if len(outcome.Completion.ErrorEvidence) == 0 || outcome.Settlement.Status != tickets.CostSettlementUnknown {
		t.Fatalf("a cancelled run must carry error evidence and settle unknown: %+v", outcome)
	}
	requireValidCompletion(t, outcome)
}

func TestAgentRunnerPinnedCLIRunTornEventLogIsUncertain(t *testing.T) {
	run, store, workspace := newAgentRunnerRunFixture(t, "torn-events")
	request := agentRunnerRunRequest(workspace)
	outcome, err := run.Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	// Every other artifact is complete here, so only the event log can withhold the
	// completion: a gap that reaches a completed receipt is a hole in the evidence.
	if outcome.Completion.Status != "uncertain" {
		t.Fatalf("a torn event log must withhold the completion: %+v", outcome.Completion)
	}
	if outcome.Facts != nil || len(outcome.Completion.ErrorEvidence) == 0 {
		t.Fatalf("a torn event log must degrade without facts: %+v", outcome)
	}
	if !outcome.Completion.LogsComplete || !outcome.Completion.UsageComplete {
		t.Fatalf("the failure must come from the event log, not the other artifacts: %+v", outcome.Completion)
	}
	evidenceDir, err := store.RunEvidenceDir(request.Launch.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(evidenceDir, agentRunnerManifestPreFileName)); err != nil {
		t.Fatalf("the other artifacts must still be collected: %v", err)
	}
	requireValidCompletion(t, outcome)
}

func TestAgentRunnerPinnedCLIRunCancellationIsCancelled(t *testing.T) {
	run, store, workspace := newAgentRunnerRunFixture(t, "run")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// The cancellation is fired once the process is live (the identity mirror is
	// published during the spawn), so the test measures the cancellation of a running
	// launch rather than the precheck's own responsiveness.
	go func() {
		mirror := filepath.Join(store.RunRoot(), managedProcessIdentityFileName)
		deadline := time.Now().Add(10 * time.Second)
		for time.Now().Before(deadline) {
			if _, err := os.Stat(mirror); err == nil {
				break
			}
			time.Sleep(5 * time.Millisecond)
		}
		cancel()
	}()
	outcome, err := run.Run(ctx, agentRunnerRunRequest(workspace))
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Completion.Status != "cancelled" {
		t.Fatalf("a cancelled run must be reported as cancelled: %+v", outcome.Completion)
	}
	if outcome.Settlement.Status != tickets.CostSettlementUnknown || len(outcome.Completion.ErrorEvidence) == 0 {
		t.Fatalf("a cancelled run must carry error evidence and settle unknown: %+v", outcome)
	}
	requireValidCompletion(t, outcome)
}

func TestAgentRunnerPinnedCLIRunCrashIsFailed(t *testing.T) {
	run, _, workspace := newAgentRunnerRunFixture(t, "crash")
	outcome, err := run.Run(context.Background(), agentRunnerRunRequest(workspace))
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Completion.Status != "failed" {
		t.Fatalf("a non-zero exit must be failed: %+v", outcome.Completion)
	}
	if outcome.Facts != nil {
		t.Fatal("a failed run must carry no frozen facts")
	}
	if outcome.Completion.LogsComplete {
		t.Fatal("a crash mid-write must not report complete logs")
	}
	requireValidCompletion(t, outcome)
}

func TestAgentRunnerPinnedCLIRunMissingUsageIsUncertain(t *testing.T) {
	run, _, workspace := newAgentRunnerRunFixture(t, "exit-no-usage")
	outcome, err := run.Run(context.Background(), agentRunnerRunRequest(workspace))
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Completion.Status != "uncertain" || outcome.Completion.UsageComplete {
		t.Fatalf("an unproven usage observation must degrade the outcome: %+v", outcome.Completion)
	}
	if outcome.Settlement.Status != tickets.CostSettlementUnknown {
		t.Fatalf("an unproven usage observation may not settle as observed: %+v", outcome.Settlement)
	}
	if outcome.Facts != nil {
		t.Fatal("a usage gap must withhold the frozen facts")
	}
	requireValidCompletion(t, outcome)
}

func TestAgentRunnerPinnedCLIRunOperatorAskOutranksSuccessAndTimeout(t *testing.T) {
	run, _, workspace := newAgentRunnerRunFixture(t, "ask")
	request := agentRunnerRunRequest(workspace)
	// The CLI blocks after asking, so the watchdog fires: the operator request
	// must still win, because ignoring it would hide a pending human decision.
	request.Timeout = 300 * time.Millisecond
	outcome, err := run.Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Completion.Status != "operator-action-required" {
		t.Fatalf("an operator request must outrank a timeout: %+v", outcome.Completion)
	}
	if outcome.Facts != nil {
		t.Fatal("an operator request must not publish frozen facts")
	}
	requireValidCompletion(t, outcome)
}

func TestAgentRunnerPinnedCLIRunTimeoutStopsTheWholeTree(t *testing.T) {
	run, store, workspace := newAgentRunnerRunFixtureWithMode(t, "spawn")
	request := agentRunnerRunRequest(workspace)
	request.Timeout = 700 * time.Millisecond
	outcome, err := run.Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Completion.Status != "uncertain" || outcome.Completion.ProcessTreeStatus != "stopped" {
		t.Fatalf("a timed-out tree must be reported as stopped: %+v", outcome.Completion)
	}
	// The CLI reported the child it spawned, so the test can prove the whole tree
	// died instead of observing only the parent.
	evidenceDir := agentRunnerPinnedCLIEvidenceDir(t, store)
	raw, err := os.ReadFile(filepath.Join(evidenceDir, agentRunnerChildPIDReportFileName))
	if err != nil {
		t.Fatalf("the CLI must report the child it spawned: %v", err)
	}
	pid := parseAgentRunnerPID(t, string(raw))
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		identity, err := guard.InspectProcess(pid)
		if err != nil {
			return
		}
		if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("a descendant survived the timeout stop")
}

func newAgentRunnerRunFixtureWithMode(t *testing.T, mode string) (*AgentRunnerPinnedCLIRun, *AgentRunnerReplayStore, string) {
	t.Helper()
	return newAgentRunnerRunFixture(t, mode)
}

func TestAgentRunnerPinnedCLIRunNaturalExitLeavesNoDescendant(t *testing.T) {
	run, store, workspace := newAgentRunnerRunFixture(t, "spawn-exit")
	outcome, err := run.Run(context.Background(), agentRunnerRunRequest(workspace))
	if err != nil {
		t.Fatal(err)
	}
	evidenceDir := agentRunnerPinnedCLIEvidenceDir(t, store)
	raw, err := os.ReadFile(filepath.Join(evidenceDir, agentRunnerChildPIDReportFileName))
	if err != nil {
		t.Fatalf("the CLI must report the child it spawned: %v", err)
	}
	pid := parseAgentRunnerPID(t, string(raw))
	childGone := waitForAgentRunnerProcessGone(pid)
	// The invariant is not the status but the guarantee behind it: a completed
	// outcome must never leave a descendant behind. On Windows the guard's job
	// close kills the whole job with the parent, on Unix the process group
	// verification catches the residue and degrades the outcome instead.
	if outcome.Completion.Status == "completed" {
		if !childGone {
			t.Fatalf("a completed run must not leave a descendant behind (pid %d)", pid)
		}
		return
	}
	if outcome.Facts != nil || len(outcome.Completion.ErrorEvidence) == 0 {
		t.Fatalf("a residue-degraded outcome must carry error evidence and no facts: %+v", outcome)
	}
	requireValidCompletion(t, outcome)
}

// waitForAgentRunnerProcessGone reports whether one pid is gone within a bound.
func waitForAgentRunnerProcessGone(pid int) bool {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		identity, err := guard.InspectProcess(pid)
		if err != nil {
			return true
		}
		if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

func TestAgentRunnerVerifyTreeGoneDetectsSurvivors(t *testing.T) {
	dir := t.TempDir()
	// No report: the CLI spawned nothing, so the tree is gone.
	if gone, err := agentRunnerVerifyTreeGone(context.Background(), dir, 200*time.Millisecond); err != nil || !gone {
		t.Fatalf("an absent report means no children: gone=%v err=%v", gone, err)
	}
	// A dead pid: gone.
	dead := exec.Command(os.Args[0], "-test.run=^$")
	if err := dead.Run(); err != nil {
		t.Fatal(err)
	}
	// A pid no longer owned by the reported child: gone. This is also the decision
	// the verifier makes when the OS refuses to open the pid at all - a run whose
	// CLI spawns a child and waits for it must not degrade merely because the child
	// has already exited, so an uninspectable pid counts as gone and the platform
	// containment (job close, process group) stays the primary guarantee.
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerChildPIDReportFileName, fmt.Appendf(nil, "%d\n", dead.Process.Pid))
	if gone, err := agentRunnerVerifyTreeGone(context.Background(), dir, 500*time.Millisecond); err != nil || !gone {
		t.Fatalf("an exited child must count as gone: gone=%v err=%v", gone, err)
	}
	// A live pid (this test process): the verification must refuse to call the tree
	// gone, otherwise a completed receipt could hide a surviving process.
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerChildPIDReportFileName, fmt.Appendf(nil, "%d\n", os.Getpid()))
	if gone, err := agentRunnerVerifyTreeGone(context.Background(), dir, 300*time.Millisecond); err != nil || gone {
		t.Fatalf("a live child must keep the tree unproven: gone=%v err=%v", gone, err)
	}
	// A malformed report is refused instead of being read as "no children".
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerChildPIDReportFileName, []byte("not-a-pid\n"))
	if _, err := agentRunnerVerifyTreeGone(context.Background(), dir, 100*time.Millisecond); !errors.Is(err, ErrInvalidAgentRunnerRun) {
		t.Fatalf("a malformed child report must fail closed, got %v", err)
	}
}

func TestAgentRunnerTerminalStatusMapsEveryDegradation(t *testing.T) {
	stopped := &guard.TerminationEvidence{Outcome: "stopped"}
	cleanEvidence := AgentRunnerEvidence{
		LogsComplete:            true,
		UsageComplete:           true,
		EventsComplete:          true,
		ManifestsComplete:       true,
		ProcessStopEvidenceHash: evidence.Digest("", []byte("stop")),
	}
	stoppedResult := guard.ProcessResult{ExitCode: -1, Termination: stopped}
	cleanResult := guard.ProcessResult{ExitCode: 0}
	cases := []struct {
		name        string
		watch       AgentRunnerWatchResult
		result      guard.ProcessResult
		stopProof   *guard.TerminationEvidence
		collected   AgentRunnerEvidence
		treeGone    bool
		wantStatus  string
		wantTree    string
		wantNoProof bool
	}{
		{name: "clean", watch: AgentRunnerWatchResult{}, result: cleanResult, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "completed", wantTree: "stopped", wantNoProof: true},
		{name: "operator ask", watch: AgentRunnerWatchResult{TimedOut: true}, result: cleanResult, stopProof: stopped, collected: AgentRunnerEvidence{LogsComplete: true, UsageComplete: true, EventsComplete: true, ManifestsComplete: true, OperatorAsk: true}, treeGone: true, wantStatus: "operator-action-required", wantTree: "stopped"},
		{name: "timeout", watch: AgentRunnerWatchResult{TimedOut: true}, result: stoppedResult, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "stopped"},
		{name: "cancel", watch: AgentRunnerWatchResult{Err: context.Canceled}, result: stoppedResult, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "cancelled", wantTree: "stopped"},
		{name: "caller deadline", watch: AgentRunnerWatchResult{Err: context.DeadlineExceeded}, result: stoppedResult, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "cancelled", wantTree: "stopped"},
		{name: "tree remained", watch: AgentRunnerWatchResult{Err: guard.ErrProcessTreeRemained}, result: stoppedResult, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "stopped"},
		{name: "termination uncertain", watch: AgentRunnerWatchResult{Err: guard.ErrTerminationUncertain}, result: stoppedResult, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "stopped"},
		{name: "unclassified stop", watch: AgentRunnerWatchResult{Err: errors.New("stop reason unknown")}, result: stoppedResult, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "stopped"},
		{
			// A run the guard stopped reports a platform code (often -1). That code is
			// not an exit status of the run, so it must not make the run "failed".
			name: "stopped, no watchdog verdict", watch: AgentRunnerWatchResult{}, result: guard.ProcessResult{ExitCode: 5, Termination: stopped}, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "stopped",
		},
		{name: "tree alive", watch: AgentRunnerWatchResult{}, result: cleanResult, stopProof: stopped, collected: cleanEvidence, treeGone: false, wantStatus: "uncertain", wantTree: "unknown"},
		{name: "stop unproven", watch: AgentRunnerWatchResult{}, result: cleanResult, stopProof: nil, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "unknown"},
		{name: "stop forced", watch: AgentRunnerWatchResult{}, result: cleanResult, stopProof: &guard.TerminationEvidence{Outcome: "terminated"}, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "unknown"},
		{name: "exit code", watch: AgentRunnerWatchResult{}, result: guard.ProcessResult{ExitCode: 3}, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "failed", wantTree: "stopped"},
		{name: "log gap", watch: AgentRunnerWatchResult{}, result: cleanResult, stopProof: stopped, collected: AgentRunnerEvidence{UsageComplete: true, EventsComplete: true, ManifestsComplete: true}, treeGone: true, wantStatus: "uncertain", wantTree: "stopped"},
		{name: "events gap", watch: AgentRunnerWatchResult{}, result: cleanResult, stopProof: stopped, collected: AgentRunnerEvidence{LogsComplete: true, UsageComplete: true, ManifestsComplete: true}, treeGone: true, wantStatus: "uncertain", wantTree: "stopped"},
		{name: "manifest gap", watch: AgentRunnerWatchResult{}, result: cleanResult, stopProof: stopped, collected: AgentRunnerEvidence{LogsComplete: true, UsageComplete: true, EventsComplete: true}, treeGone: true, wantStatus: "uncertain", wantTree: "stopped"},
		{
			// A watchdog verdict with a placeholder process result (no exit status at
			// all) must still degrade: this is the shape a deadline-killed run had when
			// the launcher read its process result before the run had reconciled it.
			name: "timeout without a stopped process", watch: AgentRunnerWatchResult{TimedOut: true, Err: errors.New("AgentRunner run timed out")}, result: guard.ProcessResult{}, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "stopped",
		},
		{
			// A watchdog verdict races the process exit: the timeout must win, because
			// the run was killed rather than failed by its own exit status.
			name: "timeout on a natural failure", watch: AgentRunnerWatchResult{TimedOut: true, Err: errors.New("AgentRunner run timed out")}, result: guard.ProcessResult{ExitCode: 7}, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "stopped",
		},
		{
			// The same shape with an error the run cannot classify: fail closed rather
			// than letting an unknown stop reason fall through to a completion.
			name: "unclassified error without a stopped process", watch: AgentRunnerWatchResult{Err: errors.New("stop reason unknown")}, result: guard.ProcessResult{}, stopProof: stopped, collected: cleanEvidence, treeGone: true, wantStatus: "uncertain", wantTree: "stopped",
		},
		{
			// A natural failure outranks a gap: the exit status is a first-class fact, and
			// a gap can only block a completion. Each artifact gap is checked here.
			name: "usage gap on a natural failure", watch: AgentRunnerWatchResult{}, result: guard.ProcessResult{ExitCode: 7}, stopProof: stopped, collected: AgentRunnerEvidence{LogsComplete: true, EventsComplete: true, ManifestsComplete: true}, treeGone: true, wantStatus: "failed", wantTree: "stopped",
		},
		{
			name: "events gap on a natural failure", watch: AgentRunnerWatchResult{}, result: guard.ProcessResult{ExitCode: 7}, stopProof: stopped, collected: AgentRunnerEvidence{LogsComplete: true, UsageComplete: true, ManifestsComplete: true}, treeGone: true, wantStatus: "failed", wantTree: "stopped",
		},
		{
			name: "manifest gap on a natural failure", watch: AgentRunnerWatchResult{}, result: guard.ProcessResult{ExitCode: 7}, stopProof: stopped, collected: AgentRunnerEvidence{LogsComplete: true, UsageComplete: true, EventsComplete: true}, treeGone: true, wantStatus: "failed", wantTree: "stopped",
		},
		{
			// A natural failure with an unproven tree: the exit status is a fact and the
			// tree status is a separate fact, so the receipt must carry both rather than
			// hiding the failure behind the tree. The tree gap only blocks a completion.
			name: "natural failure with an unproven tree", watch: AgentRunnerWatchResult{}, result: guard.ProcessResult{ExitCode: 7}, stopProof: stopped, collected: cleanEvidence, treeGone: false, wantStatus: "failed", wantTree: "unknown",
		},
		{
			// The receipt schema requires sorted, unique digests, so the mapping must
			// normalize whatever the collection and the watchdog reported.
			name:       "evidence normalization",
			watch:      AgentRunnerWatchResult{TimedOut: true, Evidence: []string{"sha256:dd"}},
			result:     stoppedResult,
			stopProof:  stopped,
			collected:  AgentRunnerEvidence{ErrorEvidence: []string{"sha256:bb", "sha256:aa", "sha256:bb", "sha256:dd"}},
			treeGone:   true,
			wantStatus: "uncertain",
			wantTree:   "stopped",
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			status, treeStatus, errorEvidence := agentRunnerTerminalStatus("request-one", testCase.watch, testCase.result, testCase.stopProof, testCase.collected, testCase.treeGone)
			if status != testCase.wantStatus || treeStatus != testCase.wantTree {
				t.Fatalf("status=%s tree=%s, want %s/%s", status, treeStatus, testCase.wantStatus, testCase.wantTree)
			}
			// A completion never carries error evidence, and every degraded outcome
			// must carry some: a receipt that says "not completed" without naming a
			// reason would be unactionable.
			if testCase.wantNoProof {
				if len(errorEvidence) != 0 {
					t.Fatalf("a completed outcome must not carry error evidence: %v", errorEvidence)
				}
				return
			}
			if len(errorEvidence) == 0 {
				t.Fatal("a degraded outcome must carry error evidence")
			}
			if !slices.IsSorted(errorEvidence) || len(slices.Compact(slices.Clone(errorEvidence))) != len(errorEvidence) {
				t.Fatalf("error evidence must be sorted and unique: %v", errorEvidence)
			}
		})
	}
}

func TestAgentRunnerTerminalStatusArchivesTheStopEvidence(t *testing.T) {
	stop := evidence.Digest("", []byte("stop-evidence"))
	complete := AgentRunnerEvidence{LogsComplete: true, UsageComplete: true, EventsComplete: true, ManifestsComplete: true, ProcessStopEvidenceHash: stop}
	proof := &guard.TerminationEvidence{Outcome: "stopped"}
	status, _, errorEvidence := agentRunnerTerminalStatus("request-one", AgentRunnerWatchResult{TimedOut: true}, guard.ProcessResult{ExitCode: -1, Termination: proof}, proof, complete, true)
	if status != "uncertain" {
		t.Fatalf("the fixture must degrade: %s", status)
	}
	// A stop that is only held inside the collection is not archived evidence: the
	// operator has to be able to find the termination record from the receipt.
	if !slices.Contains(errorEvidence, stop) {
		t.Fatalf("a degraded receipt must cite the stop evidence: %v", errorEvidence)
	}
	// A completed receipt carries no error evidence at all, so it never cites it.
	completedStatus, _, completedEvidence := agentRunnerTerminalStatus("request-one", AgentRunnerWatchResult{}, guard.ProcessResult{ExitCode: 0}, proof, complete, true)
	if completedStatus != "completed" || len(completedEvidence) != 0 {
		t.Fatalf("a completed run must stay clean: %s %v", completedStatus, completedEvidence)
	}
}

func TestAgentRunnerSettlementOnlyObservesCompletedUsage(t *testing.T) {
	completed := AgentRunnerEvidence{
		LogsComplete:  true,
		UsageComplete: true,
		UsageHash:     evidence.Digest("", []byte("usage")),
		Usage:         AgentRunnerUsage{Calls: 2, Tokens: 30, DurationMs: 40},
	}
	settlement := agentRunnerSettlementFor("completed", completed)
	if settlement.Status != tickets.CostSettlementObserved {
		t.Fatalf("a completed run with complete usage must settle as observed: %+v", settlement)
	}
	if settlement.ObservedCalls == nil || *settlement.ObservedCalls != 2 || settlement.ObservedTokens == nil || *settlement.ObservedTokens != 30 {
		t.Fatalf("the observation must carry the measured usage: %+v", settlement)
	}
	if len(settlement.ProviderEvidence) != 1 || settlement.ProviderEvidence[0] != completed.UsageHash {
		t.Fatalf("the observation must cite the usage artifact: %+v", settlement)
	}
	// An incomplete usage artifact may never be settled as observed, even when the
	// run itself completed and the numbers happen to be present.
	gapped := completed
	gapped.UsageComplete = false
	for _, testCase := range []struct {
		name      string
		status    string
		collected AgentRunnerEvidence
	}{
		{name: "usage gap", status: "completed", collected: gapped},
		{name: "uncertain run", status: "uncertain", collected: completed},
		{name: "cancelled run", status: "cancelled", collected: completed},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			unknown := agentRunnerSettlementFor(testCase.status, testCase.collected)
			if unknown.Status != tickets.CostSettlementUnknown {
				t.Fatalf("status=%s must settle as unknown: %+v", testCase.status, unknown)
			}
			if unknown.ObservedCalls != nil || unknown.ObservedTokens != nil || len(unknown.ProviderEvidence) != 0 {
				t.Fatalf("an unknown settlement must not carry observations: %+v", unknown)
			}
		})
	}
}

func TestAgentRunnerVerifyTreeGoneUsesTheLivenessObservation(t *testing.T) {
	dir := t.TempDir()
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerChildPIDReportFileName, []byte("4242\n"))
	identity := guard.ProcessIdentity{PID: 4242, StartToken: "start-4242"}
	inspect := func(pid int) (guard.ProcessIdentity, error) { return identity, nil }

	// An observed live process keeps the tree unproven.
	liveCalls := 0
	live := func(guard.ProcessIdentity) (bool, error) { liveCalls++; return true, nil }
	if gone, err := agentRunnerVerifyTreeGoneWith(context.Background(), dir, 150*time.Millisecond, inspect, live); err != nil || gone {
		t.Fatalf("a live reported child must keep the tree unproven: gone=%v err=%v", gone, err)
	}
	if liveCalls < 2 {
		t.Fatalf("a live child must be re-checked until the bound expires, got %d observation(s)", liveCalls)
	}
	// An observed dead process proves the tree is gone.
	if gone, err := agentRunnerVerifyTreeGoneWith(context.Background(), dir, time.Second, inspect, func(guard.ProcessIdentity) (bool, error) { return false, nil }); err != nil || !gone {
		t.Fatalf("a dead reported child must prove the tree gone: gone=%v err=%v", gone, err)
	}
	// An observation error is conservative: it may not be read as proof of a stop.
	if gone, err := agentRunnerVerifyTreeGoneWith(context.Background(), dir, 150*time.Millisecond, inspect, func(guard.ProcessIdentity) (bool, error) { return false, errors.New("cannot observe") }); err != nil || gone {
		t.Fatalf("an unobservable liveness must fail closed: gone=%v err=%v", gone, err)
	}
	// Every reported pid is checked, not just the first.
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerChildPIDReportFileName, []byte("4242\n4243\n"))
	seen := make([]int, 0, 2)
	perPID := func(pid int) (guard.ProcessIdentity, error) {
		seen = append(seen, pid)
		return guard.ProcessIdentity{PID: pid, StartToken: "start"}, nil
	}
	aliveForSecond := func(identity guard.ProcessIdentity) (bool, error) { return identity.PID == 4243, nil }
	if gone, err := agentRunnerVerifyTreeGoneWith(context.Background(), dir, 100*time.Millisecond, perPID, aliveForSecond); err != nil || gone {
		t.Fatalf("a live second child must keep the tree unproven: gone=%v err=%v", gone, err)
	}
	if len(seen) < 2 || seen[0] != 4242 || seen[1] != 4243 {
		t.Fatalf("every reported pid must be observed, saw %v", seen)
	}
}

// mutatedWorkspaceFile is the file the stub's mutate-workspace mode writes into
// the workspace it was given.
const mutatedWorkspaceFile = "agent-stub-mutation.txt"

// workspaceManifestCapturer is a real manifest capturer: it lists the workspace
// contents, so a run that changes the workspace changes the manifest and the diff.
// It also records whether the spawn had already happened when the first capture
// was taken, which is how the caller's ordering can be proven rather than assumed.
type workspaceManifestCapturer struct {
	root       string
	mirrorPath string
	calls      int
	spawned    bool
}

func (capturer *workspaceManifestCapturer) CaptureWorkspaceManifest(_ context.Context, root string) ([]byte, error) {
	capturer.calls++
	if capturer.calls == 1 && capturer.mirrorPath != "" {
		// The launcher publishes the process identity mirror as part of starting the
		// process, so its presence proves the spawn already happened.
		if _, err := os.Stat(capturer.mirrorPath); err == nil {
			capturer.spawned = true
		}
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		lines = append(lines, fmt.Sprintf("%s %d", entry.Name(), info.Size()))
	}
	slices.Sort(lines)
	return []byte(strings.Join(lines, "\n") + "\n"), nil
}

func (capturer *workspaceManifestCapturer) DiffWorkspaceManifests(_ context.Context, before, after []byte) ([]byte, error) {
	if bytes.Equal(before, after) {
		return []byte("no change\n"), nil
	}
	return append(append([]byte("changed\n"), before...), after...), nil
}

func TestAgentRunnerPinnedCLIRunUsesThePreRunManifest(t *testing.T) {
	run, store, workspace := newAgentRunnerRunFixtureArgs(t, []string{"mutate-workspace", "-sleep", "100ms"})
	request := agentRunnerRunRequest(workspace)
	capturer := &workspaceManifestCapturer{root: workspace, mirrorPath: filepath.Join(store.RunRoot(), managedProcessIdentityFileName)}
	request.Capturer = capturer
	outcome, err := run.Run(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if capturer.spawned {
		// The pre manifest may only describe the workspace as it was before the CLI
		// could touch it: a capture taken after the spawn is not a pre-state.
		t.Fatal("the pre-run manifest was captured after the process was started")
	}
	if outcome.Completion.Status != "completed" {
		t.Fatalf("a clean run that changed the workspace must complete: %+v", outcome.Completion)
	}
	evidenceDir, err := store.RunEvidenceDir(request.Launch.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	pre, err := os.ReadFile(filepath.Join(evidenceDir, agentRunnerManifestPreFileName))
	if err != nil {
		t.Fatal(err)
	}
	post, err := os.ReadFile(filepath.Join(evidenceDir, agentRunnerManifestPostFileName))
	if err != nil {
		t.Fatal(err)
	}
	// The pre manifest must describe the workspace as it was before the run: a
	// capture taken afterwards would make the diff claim "no change" for a run that
	// did change the workspace.
	if strings.Contains(string(pre), mutatedWorkspaceFile) {
		t.Fatalf("the pre manifest was captured after the run: %s", pre)
	}
	if !strings.Contains(string(post), mutatedWorkspaceFile) {
		t.Fatalf("the post manifest must show what the run wrote: %s", post)
	}
	diff, err := os.ReadFile(filepath.Join(evidenceDir, agentRunnerDiffFileName))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(diff), "no change") {
		t.Fatalf("a run that changed the workspace must not diff as no change: %s", diff)
	}
	if outcome.Facts == nil || outcome.Facts.ManifestHash != evidence.Digest("", post) {
		t.Fatal("the manifest fact must be the post-manifest digest")
	}
}

func TestAgentRunnerPinnedCLIRunManifestGapIsUncertainNotAnError(t *testing.T) {
	run, store, workspace := newAgentRunnerRunFixture(t, "exit")
	request := agentRunnerRunRequest(workspace)
	// The diff cannot be computed: the contract makes that uncertain, not fatal, so
	// the caller still receives an outcome naming the gap.
	request.Capturer = &fakeManifestCapturer{after: []byte("after"), diffErr: errors.New("diff failed")}
	outcome, err := run.Run(context.Background(), request)
	if err != nil {
		t.Fatalf("a manifest gap must be an outcome, not an error: %v", err)
	}
	if outcome.Completion.Status != "uncertain" || outcome.Facts != nil {
		t.Fatalf("a manifest gap must degrade without facts: %+v", outcome)
	}
	if len(outcome.Completion.ErrorEvidence) == 0 {
		t.Fatal("the gap must be named as error evidence")
	}
	evidenceDir, err := store.RunEvidenceDir(request.Launch.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(evidenceDir, agentRunnerDiffFileName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a gap must not leave a diff artifact behind: %v", err)
	}
}

func TestAgentRunnerPinnedCLIRunRefusesUnusableRequests(t *testing.T) {
	run, _, workspace := newAgentRunnerRunFixture(t, "exit")
	for name, mutate := range map[string]func(*AgentRunnerRunRequest){
		"missing request hash": func(request *AgentRunnerRunRequest) { request.RequestHash = "" },
		"missing task":         func(request *AgentRunnerRunRequest) { request.TaskID = "" },
		"missing step":         func(request *AgentRunnerRunRequest) { request.StepID = "" },
		"missing attempt":      func(request *AgentRunnerRunRequest) { request.Attempt = 0 },
		"missing timeout":      func(request *AgentRunnerRunRequest) { request.Timeout = 0 },
		"unbounded logs":       func(request *AgentRunnerRunRequest) { request.MaxLogBytes = 0 },
	} {
		t.Run(name, func(t *testing.T) {
			request := agentRunnerRunRequest(workspace)
			mutate(&request)
			if _, err := run.Run(context.Background(), request); !errors.Is(err, ErrInvalidAgentRunnerRun) {
				t.Fatalf("expected a fail-closed request rejection, got %v", err)
			}
		})
	}
	if _, err := (&AgentRunnerPinnedCLIRun{}).Run(context.Background(), agentRunnerRunRequest(workspace)); !errors.Is(err, ErrInvalidAgentRunnerRun) {
		t.Fatalf("a run without a launcher must be refused, got %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := run.Run(ctx, agentRunnerRunRequest(workspace)); !errors.Is(err, context.Canceled) {
		t.Fatalf("a cancelled run must refuse to start, got %v", err)
	}
}
