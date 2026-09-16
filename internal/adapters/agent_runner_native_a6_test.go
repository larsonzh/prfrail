//go:build a6native

// Native A6 experiments: they spawn real processes and therefore never run in the
// default test suite. Run them with:
//
//	go test -tags a6native -count=1 -run TestA6Native -v ./internal/adapters/
//
// Every experiment asserts one soundness claim about the offline run contract
// instead of a platform-specific outcome, because the containment mechanism
// differs: Windows kills the job when the launcher closes it, Unix verifies the
// process group.
package adapters

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/tickets"
)

const (
	agentStubVersion  = "agent-stub-0.1.0"
	agentStubPIDFile  = "stub.pid"
	agentStubChildPID = "stub-child.pid"
)

var (
	agentStubOnce sync.Once
	agentStubPath string
	agentStubErr  error
)

// buildAgentStub builds the deterministic CLI once per test binary. It is the
// artifact under experiment, so it is built from source rather than mocked.
func buildAgentStub(t *testing.T) string {
	t.Helper()
	agentStubOnce.Do(func() {
		root, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			agentStubErr = err
			return
		}
		target := filepath.Join(os.TempDir(), "proofrail-agent-stub-a6")
		if err := os.MkdirAll(target, 0o755); err != nil {
			agentStubErr = err
			return
		}
		name := "agent-stub"
		if strings.EqualFold(filepath.Ext(os.Args[0]), ".exe") {
			name += ".exe"
		}
		agentStubPath = filepath.Join(target, name)
		build := exec.Command("go", "build", "-o", agentStubPath, "./tools/agent-stub")
		build.Dir = root
		if output, err := build.CombinedOutput(); err != nil {
			agentStubErr = fmt.Errorf("build agent-stub: %v: %s", err, output)
		}
	})
	if agentStubErr != nil {
		t.Skipf("agent-stub is unavailable: %v", agentStubErr)
	}
	return agentStubPath
}

// nativeRun builds one run around the stub with the given mode arguments.
func nativeRun(t *testing.T, modeArgs []string, workspace string) (*AgentRunnerPinnedCLIRun, *AgentRunnerReplayStore, string) {
	t.Helper()
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(AgentRunnerPinnedCLIConfig{
		Executable:   buildAgentStub(t),
		Version:      agentStubVersion,
		VersionArgs:  []string{"-version"},
		RunArgs:      modeArgs,
		CheckTimeout: 20 * time.Second,
		Grace:        500 * time.Millisecond,
	}, store)
	if err != nil {
		t.Fatal(err)
	}
	return &AgentRunnerPinnedCLIRun{Launcher: launcher, Store: store}, store, workspace
}

// nativeRequest pins one round of an experiment.
func nativeRequest(workspace, requestID string, timeout time.Duration) AgentRunnerRunRequest {
	request := agentRunnerRunRequest(workspace)
	request.Launch.RequestID = requestID
	request.Timeout = timeout
	request.VerifyWait = 1500 * time.Millisecond
	return request
}

// reportPID reads one pid artifact and waits for the process to disappear.
func reportPID(t *testing.T, path string) (int, bool) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(raw)))
	if err != nil || pid <= 0 {
		return 0, false
	}
	return pid, true
}

// waitGone reports whether one pid is gone within the bound and never leaves a
// process behind for the rest of the suite.
func waitGone(t *testing.T, pid int, bound time.Duration) bool {
	t.Helper()
	deadline := time.Now().Add(bound)
	for time.Now().Before(deadline) {
		identity, err := guard.InspectProcess(pid)
		if err != nil {
			return true
		}
		alive, err := guard.ProcessAlive(identity)
		if err != nil || !alive {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	// The experiment must not leak the process it is measuring.
	identity, err := guard.InspectProcess(pid)
	if err == nil {
		if _, stopErr := guard.StopProcessIdentity(context.Background(), identity, 2*time.Second); stopErr != nil {
			t.Logf("leftover process %d could not be stopped: %v", pid, stopErr)
		}
	}
	return false
}

// TestA6NativeE1 checks the natural-exit tree rule: a parent that exits cleanly
// while a descendant keeps running may never be reported as a completed run with a
// live tree. Windows contains the child through the job object, Unix through the
// process group, so the experiment asserts the containment claim itself.
func TestA6NativeE1(t *testing.T) {
	workspace := t.TempDir()
	const rounds = 10
	contained, refused := 0, 0
	for round := 0; round < rounds; round++ {
		run, store, _ := nativeRun(t, []string{"-mode", "spawn-exit", "-child-life", "20s"}, workspace)
		requestID := fmt.Sprintf("request-%02d", round)
		outcome, err := run.Run(context.Background(), nativeRequest(workspace, requestID, 15*time.Second))
		if err != nil {
			t.Fatalf("round %d: %v", round, err)
		}
		evidenceDir, err := store.RunEvidenceDir(requestID)
		if err != nil {
			t.Fatal(err)
		}
		pid, reported := reportPID(t, filepath.Join(evidenceDir, agentStubChildPID))
		if !reported {
			t.Fatalf("round %d: the CLI must report the child it spawned", round)
		}
		gone := waitGone(t, pid, 3*time.Second)
		status := outcome.Completion.Status
		t.Logf("round %d: status=%s tree=%s child=%d gone=%v logs=%v usage=%v errors=%d",
			round, status, outcome.Completion.ProcessTreeStatus, pid, gone,
			outcome.Completion.LogsComplete, outcome.Completion.UsageComplete, len(outcome.Completion.ErrorEvidence))
		switch {
		case status == "completed" && gone:
			contained++
			if outcome.Facts == nil {
				t.Fatalf("round %d: a completed run must publish frozen facts", round)
			}
		case status == "completed" && !gone:
			t.Fatalf("round %d: a completed run left the descendant %d alive", round, pid)
		default:
			refused++
			if outcome.Facts != nil || len(outcome.Completion.ErrorEvidence) == 0 {
				t.Fatalf("round %d: a refused run must carry error evidence and no facts: %+v", round, outcome)
			}
		}
	}
	t.Logf("E1 summary: contained=%d refused=%d of %d rounds", contained, refused, rounds)
	if contained == 0 && refused == 0 {
		t.Fatal("E1 measured nothing")
	}
}

// TestA6NativeE2 walks the timeout boundary in both directions: a run that finishes
// inside the window completes with a proven stop, one that does not may never be
// reported as completed, and no round may leave the stub alive.
func TestA6NativeE2(t *testing.T) {
	workspace := t.TempDir()
	const rounds = 20
	completed, degraded := 0, 0
	for round := 0; round < rounds; round++ {
		inside := round%2 == 0
		args := []string{"-mode", "run", "-sleep", "300ms"}
		timeout := 3 * time.Second
		if !inside {
			args = []string{"-mode", "run", "-sleep", "3s"}
			timeout = 300 * time.Millisecond
		}
		run, store, _ := nativeRun(t, args, workspace)
		requestID := fmt.Sprintf("request-%02d", round)
		outcome, err := run.Run(context.Background(), nativeRequest(workspace, requestID, timeout))
		if err != nil {
			t.Fatalf("round %d: %v", round, err)
		}
		evidenceDir, err := store.RunEvidenceDir(requestID)
		if err != nil {
			t.Fatal(err)
		}
		pid, reported := reportPID(t, filepath.Join(evidenceDir, agentStubPIDFile))
		if !reported {
			t.Fatalf("round %d: the stub must report its own pid", round)
		}
		if !waitGone(t, pid, 3*time.Second) {
			t.Fatalf("round %d: the stub %d survived the run", round, pid)
		}
		status := outcome.Completion.Status
		t.Logf("round %d: inside=%v timeout=%s status=%s tree=%s errors=%d",
			round, inside, timeout, status, outcome.Completion.ProcessTreeStatus, len(outcome.Completion.ErrorEvidence))
		if status == "completed" {
			completed++
			if !inside {
				t.Fatalf("round %d: a run that outlived its timeout was reported as completed", round)
			}
			if outcome.Completion.ProcessTreeStatus != "stopped" || outcome.Facts == nil {
				t.Fatalf("round %d: a completion needs a proven stop and facts: %+v", round, outcome.Completion)
			}
			continue
		}
		degraded++
		if inside {
			t.Fatalf("round %d: a run that finished inside the window was degraded to %s", round, status)
		}
		if outcome.Facts != nil {
			t.Fatalf("round %d: a degraded run must not publish facts", round)
		}
	}
	t.Logf("E2 summary: completed=%d degraded=%d of %d rounds", completed, degraded, rounds)
	if completed == 0 || degraded == 0 {
		t.Fatalf("E2 must exercise both sides of the boundary: completed=%d degraded=%d", completed, degraded)
	}
}

// TestA6NativeE3 races cancellation against the timeout. Exactly one winner may
// classify the run, the receipt must say which, and the process must be gone.
func TestA6NativeE3(t *testing.T) {
	workspace := t.TempDir()
	const rounds = 20
	cancelled, timedOut := 0, 0
	for round := 0; round < rounds; round++ {
		run, store, _ := nativeRun(t, []string{"-mode", "hang"}, workspace)
		requestID := fmt.Sprintf("request-%02d", round)
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			// Fire the cancellation into the timeout window: the jitter is wide enough
			// that both winners appear, which is what makes the race observable. The
			// window is far larger than the version precheck, so the precheck cannot
			// swallow the round.
			time.Sleep(1450*time.Millisecond + time.Duration(round%11)*10*time.Millisecond)
			cancel()
		}()
		outcome, err := run.Run(ctx, nativeRequest(workspace, requestID, 1500*time.Millisecond))
		cancel()
		if err != nil {
			t.Fatalf("round %d: %v", round, err)
		}
		evidenceDir, err := store.RunEvidenceDir(requestID)
		if err != nil {
			t.Fatal(err)
		}
		pid, reported := reportPID(t, filepath.Join(evidenceDir, agentStubPIDFile))
		if !reported {
			t.Fatalf("round %d: the stub must report its own pid", round)
		}
		if !waitGone(t, pid, 3*time.Second) {
			t.Fatalf("round %d: the stub %d survived the cancel/timeout race", round, pid)
		}
		status := outcome.Completion.Status
		t.Logf("round %d: status=%s tree=%s errors=%d", round, status, outcome.Completion.ProcessTreeStatus, len(outcome.Completion.ErrorEvidence))
		switch status {
		case "cancelled":
			cancelled++
		case "uncertain":
			timedOut++
		default:
			t.Fatalf("round %d: a raced run may only be cancelled or uncertain, got %s", round, status)
		}
		if outcome.Facts != nil || len(outcome.Completion.ErrorEvidence) == 0 {
			t.Fatalf("round %d: a raced run must carry error evidence and no facts", round)
		}
	}
	t.Logf("E3 summary: cancelled=%d uncertain=%d of %d rounds", cancelled, timedOut, rounds)
}

// TestA6NativeE4 injects a failure at each phase - start, stop and collection - and
// checks that the mapping never reports a completion and never loses the reason.
func TestA6NativeE4(t *testing.T) {
	workspace := t.TempDir()
	cases := []struct {
		name       string
		args       []string
		timeout    time.Duration
		wantStatus string
	}{
		{name: "start fails", args: []string{"-mode", "fail-fast", "-exit-code", "9"}, timeout: 10 * time.Second, wantStatus: "failed"},
		{name: "stop leaks a child", args: []string{"-mode", "spawn-exit", "-child-life", "20s"}, timeout: 15 * time.Second, wantStatus: "uncertain-or-contained"},
		{name: "collection loses the log marker", args: []string{"-mode", "run", "-no-marker", "-sleep", "100ms"}, timeout: 10 * time.Second, wantStatus: "uncertain"},
		{name: "collection loses the usage artifact", args: []string{"-mode", "no-usage", "-sleep", "100ms"}, timeout: 10 * time.Second, wantStatus: "uncertain"},
		{name: "collection tears the usage artifact", args: []string{"-mode", "run", "-torn-usage", "-sleep", "100ms"}, timeout: 10 * time.Second, wantStatus: "uncertain"},
		{name: "operator asks for a human", args: []string{"-mode", "ask", "-sleep", "100ms"}, timeout: 10 * time.Second, wantStatus: "operator-action-required"},
	}
	for index, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			run, store, _ := nativeRun(t, testCase.args, workspace)
			requestID := fmt.Sprintf("request-%02d", index)
			outcome, err := run.Run(context.Background(), nativeRequest(workspace, requestID, testCase.timeout))
			if err != nil {
				t.Fatal(err)
			}
			evidenceDir, err := store.RunEvidenceDir(requestID)
			if err != nil {
				t.Fatal(err)
			}
			if pid, reported := reportPID(t, filepath.Join(evidenceDir, agentStubPIDFile)); reported {
				if !waitGone(t, pid, 3*time.Second) {
					t.Fatalf("the stub %d survived the run", pid)
				}
			}
			if child, reported := reportPID(t, filepath.Join(evidenceDir, agentStubChildPID)); reported {
				if !waitGone(t, child, 3*time.Second) {
					t.Fatalf("the spawned child %d survived the run", child)
				}
			}
			status := outcome.Completion.Status
			t.Logf("%s: status=%s tree=%s errors=%d logs=%v usage=%v",
				testCase.name, status, outcome.Completion.ProcessTreeStatus,
				len(outcome.Completion.ErrorEvidence), outcome.Completion.LogsComplete, outcome.Completion.UsageComplete)
			if testCase.wantStatus == "uncertain-or-contained" {
				// The residue is either refused (Unix) or contained by the platform
				// (Windows, where the job close already killed the child).
				if status == "completed" && outcome.Facts == nil {
					t.Fatal("a completed run must publish frozen facts")
				}
				return
			}
			if status != testCase.wantStatus {
				t.Fatalf("status=%s, want %s", status, testCase.wantStatus)
			}
			if outcome.Facts != nil {
				t.Fatal("no phase may publish facts for this case")
			}
			if len(outcome.Completion.ErrorEvidence) == 0 {
				t.Fatal("the reason must be recorded as error evidence")
			}
		})
	}
}

// TestA6NativeE5 keeps the stub honest: the version precheck must refuse a stub
// that reports a different version, so the experiments cannot silently run a
// different binary than the one they pinned.
func TestA6NativeE5(t *testing.T) {
	workspace := t.TempDir()
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(AgentRunnerPinnedCLIConfig{
		Executable:   buildAgentStub(t),
		Version:      "agent-stub-9.9.9",
		VersionArgs:  []string{"-version"},
		RunArgs:      []string{"-mode", "run", "-sleep", "100ms"},
		CheckTimeout: 10 * time.Second,
		Grace:        500 * time.Millisecond,
	}, store)
	if err != nil {
		t.Fatal(err)
	}
	run := &AgentRunnerPinnedCLIRun{Launcher: launcher, Store: store}
	_, err = run.Run(context.Background(), nativeRequest(workspace, "request-00", 10*time.Second))
	if !errors.Is(err, ErrAgentRunnerLaunchVersion) {
		t.Fatalf("a mismatched pin must refuse the run, got %v", err)
	}
}

// TestA6NativeE6 proves the evidence digests describe the artifacts the stub
// actually wrote, which is what makes a replay comparable.
func TestA6NativeE6(t *testing.T) {
	workspace := t.TempDir()
	run, store, _ := nativeRun(t, []string{"-mode", "run", "-sleep", "100ms", "-calls", "7", "-tokens", "70"}, workspace)
	requestID := "request-00"
	outcome, err := run.Run(context.Background(), nativeRequest(workspace, requestID, 10*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Completion.Status != "completed" {
		t.Fatalf("the deterministic stub must complete: %+v", outcome.Completion)
	}
	evidenceDir, err := store.RunEvidenceDir(requestID)
	if err != nil {
		t.Fatal(err)
	}
	usage, err := os.ReadFile(filepath.Join(evidenceDir, agentRunnerUsageFileName))
	if err != nil {
		t.Fatal(err)
	}
	if outcome.Facts == nil || outcome.Facts.UsageHash != evidence.Digest("", usage) {
		t.Fatalf("the facts must describe the written usage artifact: %+v", outcome.Facts)
	}
	if outcome.Settlement.Status != tickets.CostSettlementObserved {
		t.Fatalf("a completed run with usage must settle as observed: %+v", outcome.Settlement)
	}
	if outcome.Settlement.ObservedCalls == nil || *outcome.Settlement.ObservedCalls != 7 || outcome.Settlement.ObservedTokens == nil || *outcome.Settlement.ObservedTokens != 70 {
		t.Fatalf("the settlement must carry the stub's numbers: %+v", outcome.Settlement)
	}
}
