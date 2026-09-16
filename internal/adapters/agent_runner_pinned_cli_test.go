package adapters

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

// agentRunnerPinnedCLITestVersion is the version the re-exec helper reports.
const agentRunnerPinnedCLITestVersion = "agent-stub-0.1.0-test"

// TestAgentRunnerPinnedCLIHelper is the deterministic pinned CLI used by the
// launcher tests. It only acts when re-executed with "--" plus a mode, so a normal
// test run simply returns.
func TestAgentRunnerPinnedCLIHelper(t *testing.T) {
	separator := -1
	for index, arg := range os.Args {
		if arg == "--" {
			separator = index
			break
		}
	}
	if separator < 0 || len(os.Args) <= separator+1 {
		return
	}
	switch os.Args[separator+1] {
	case "version":
		fmt.Fprintln(os.Stdout, agentRunnerPinnedCLITestVersion)
		os.Exit(0)
	case "version-slow":
		// A precheck that takes long enough to observe that no launch slot is held
		// while it runs.
		time.Sleep(1500 * time.Millisecond)
		fmt.Fprintln(os.Stdout, agentRunnerPinnedCLITestVersion)
		os.Exit(0)
	case "run":
		// A real agent run keeps working until it is stopped.
		helperEmit(true)
		time.Sleep(time.Minute)
	case "exit":
		helperEmit(true)
		os.Exit(0)
	case "mutate-workspace":
		// A run that changes its workspace before it exits cleanly, so a manifest
		// captured after the run would diff as "no change".
		workspace, err := os.Getwd()
		if err != nil {
			os.Exit(21)
		}
		if err := os.WriteFile(filepath.Join(workspace, mutatedWorkspaceFile), []byte("the run wrote this\n"), 0o600); err != nil {
			os.Exit(22)
		}
		helperEmit(true)
		os.Exit(0)
	case "torn-events":
		// A writer interrupted half-way: the event log is unreadable in full while
		// every other artifact is complete.
		helperEmit(true)
		helperWriteEvents("{\"type\":\"output\"}\n{\"type\":\"usa")
		os.Exit(0)
	case "exit-no-usage":
		helperEmit(false)
		os.Exit(0)
	case "ask":
		// An operator request outranks success and must survive a timeout too.
		helperEmit(true)
		helperWriteEvents("{\"type\":\"session-created\"}\n{\"type\":\"operator-action-required\"}\n")
		time.Sleep(time.Minute)
	case "no-usage":
		helperEmit(false)
		time.Sleep(time.Minute)
	case "crash":
		// A crash mid-write leaves the launcher-owned log without its marker.
		fmt.Fprintln(os.Stdout, "interrupted")
		os.Exit(7)
	case "spawn-exit":
		// The parent exits cleanly while the child it spawned keeps running: only the
		// process-tree verification can catch that, because the parent's status is 0.
		child := exec.Command(os.Args[0], "-test.run=^TestAgentRunnerPinnedCLIHelper$", "--", "sleep")
		if err := child.Start(); err != nil {
			os.Exit(19)
		}
		helperEmit(true)
		if err := os.WriteFile(filepath.Join(helperEvidenceDir(), agentRunnerChildPIDReportFileName), fmt.Appendf(nil, "%d\n", child.Process.Pid), 0o600); err != nil {
			os.Exit(20)
		}
		os.Exit(0)
	case "spawn":
		child := exec.Command(os.Args[0], "-test.run=^TestAgentRunnerPinnedCLIHelper$", "--", "sleep")
		if err := child.Start(); err != nil {
			os.Exit(9)
		}
		if err := os.WriteFile(filepath.Join(helperEvidenceDir(), "stub-child.pid"), fmt.Appendf(nil, "%d\n", child.Process.Pid), 0o600); err != nil {
			os.Exit(10)
		}
		time.Sleep(time.Hour)
	default: // "sleep"
		time.Sleep(time.Hour)
	}
	os.Exit(2)
}

// helperEvidenceDir is where the launcher told the CLI to publish its artifacts.
func helperEvidenceDir() string {
	if dir := strings.TrimSpace(os.Getenv(agentRunnerEvidenceDirEnvVar)); dir != "" {
		return dir
	}
	dir, err := os.Getwd()
	if err != nil {
		os.Exit(11)
	}
	return dir
}

// helperEmit writes the artifacts of a clean run and prints the log streams, which
// the launcher redirects into its own log artifacts.
func helperEmit(withUsage bool) {
	dir := helperEvidenceDir()
	if err := os.WriteFile(filepath.Join(dir, "stub.pid"), fmt.Appendf(nil, "%d\n", os.Getpid()), 0o600); err != nil {
		os.Exit(12)
	}
	fmt.Fprintln(os.Stdout, "hello")
	fmt.Fprintln(os.Stdout, agentRunnerLogCompletenessMarker)
	fmt.Fprintln(os.Stderr, agentRunnerLogCompletenessMarker)
	if err := os.WriteFile(filepath.Join(dir, agentRunnerEventsFileName), []byte("{\"type\":\"session-created\"}\n{\"type\":\"output\"}\n"), 0o600); err != nil {
		os.Exit(13)
	}
	if withUsage {
		if err := os.WriteFile(filepath.Join(dir, agentRunnerUsageFileName), []byte(`{"calls":2,"tokens":123,"durationMs":5}`), 0o600); err != nil {
			os.Exit(14)
		}
	}
}

// helperWriteEvents replaces the event log of a helper run.
func helperWriteEvents(content string) {
	if err := os.WriteFile(filepath.Join(helperEvidenceDir(), agentRunnerEventsFileName), []byte(content), 0o600); err != nil {
		os.Exit(18)
	}
}

func replayStoreForRunRoot(t *testing.T, runRoot string) *AgentRunnerReplayStore {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(runRoot, agentRunnerReplayRunEventDirName), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runRoot, agentRunnerReplayRunEventDirName, agentRunnerReplayRunEventFileName), []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	store, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one")
	if err != nil {
		t.Fatal(err)
	}
	return store
}

func agentRunnerPinnedCLIConfig(executable string) AgentRunnerPinnedCLIConfig {
	if executable == "" {
		executable = os.Args[0]
	}
	helper := []string{"-test.run=^TestAgentRunnerPinnedCLIHelper$", "--"}
	return AgentRunnerPinnedCLIConfig{
		Executable:   executable,
		Version:      agentRunnerPinnedCLITestVersion,
		VersionArgs:  append(append([]string{}, helper...), "version"),
		RunArgs:      append(append([]string{}, helper...), "run"),
		CheckTimeout: 10 * time.Second,
		Grace:        500 * time.Millisecond,
	}
}

func agentRunnerPinnedCLIConfigForMode(mode string) AgentRunnerPinnedCLIConfig {
	config := agentRunnerPinnedCLIConfig("")
	helper := []string{"-test.run=^TestAgentRunnerPinnedCLIHelper$", "--"}
	config.RunArgs = append(append([]string{}, helper...), mode)
	return config
}

func newAgentRunnerPinnedCLILauncher(t *testing.T, runRoot string) (*AgentRunnerPinnedCLILauncher, *AgentRunnerReplayStore) {
	t.Helper()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfig(""), store)
	if err != nil {
		t.Fatal(err)
	}
	return launcher, store
}

func agentRunnerPinnedCLIRequest(workspaceRoot string) chain.AgentRunnerLaunchRequest {
	return chain.AgentRunnerLaunchRequest{
		RequestID:     "request-one",
		RunID:         "run-one",
		AdapterID:     "adapter-one",
		Dir:           workspaceRoot,
		WorkspaceRoot: workspaceRoot,
	}
}

func agentRunnerPinnedCLIEvidenceDir(t *testing.T, store *AgentRunnerReplayStore) string {
	t.Helper()
	dir, err := store.RunEvidenceDir("request-one")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestAgentRunnerPinnedCLILauncherStartsAndStops(t *testing.T) {
	runRoot := t.TempDir()
	workspace := t.TempDir()
	launcher, store := newAgentRunnerPinnedCLILauncher(t, runRoot)
	result, err := launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(workspace))
	if err != nil {
		t.Fatal(err)
	}
	if result.LaunchID != "launch-request-one" {
		t.Fatalf("unexpected launch id %q", result.LaunchID)
	}
	identity, err := decodeAgentRunnerProcessID(result.ProcessID)
	if err != nil {
		t.Fatal(err)
	}
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("the spawned process must be alive: alive=%v err=%v", alive, err)
	}
	mirror, err := loadAgentRunnerProcessIdentityMirror(store.RunRoot())
	if err != nil {
		t.Fatal(err)
	}
	if mirror != identity {
		t.Fatalf("the identity mirror must match the reported process id: %+v vs %+v", mirror, identity)
	}
	evidenceDir := agentRunnerPinnedCLIEvidenceDir(t, store)
	if _, err := os.Stat(filepath.Join(evidenceDir, agentRunnerLogsDirectoryName, agentRunnerStdoutFileName)); err != nil {
		t.Fatalf("the launcher must own the log artifacts: %v", err)
	}
	// The CLI must publish into the evidence directory it was handed, which is how
	// the run keeps its artifacts out of the workspace. The artifact write is
	// asynchronous, so it is awaited with a bound.
	waitForAgentRunnerArtifact(t, filepath.Join(evidenceDir, "stub.pid"))
	// The child publishes its artifacts beside the launcher's logs, never inside
	// the workspace it was handed.
	if _, err := os.Stat(filepath.Join(workspace, agentRunnerUsageFileName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the evidence directory keeps the workspace clean: %v", err)
	}

	proof, err := launcher.StopAgentRunnerProcess(context.Background(), result.LaunchID)
	if err != nil {
		t.Fatalf("stop failed: %v (%+v)", err, proof)
	}
	if proof.Outcome != "stopped" {
		t.Fatalf("expected a proven stop, got %+v", proof)
	}
	if alive, err := guard.ProcessAlive(identity); err == nil && alive {
		t.Fatal("the process must be gone after a proven stop")
	}
	if launcher.registry.len() != 0 {
		t.Fatalf("a stopped launch must release its handle: %d", launcher.registry.len())
	}
}

func TestAgentRunnerPinnedCLILauncherRefusesASecondLiveLaunch(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	first := agentRunnerPinnedCLIRequest(workspace)
	ctx := context.Background()
	launched, err := launcher.StartAgentRunnerProcess(ctx, first)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = launcher.StopAgentRunnerProcess(context.Background(), launched.LaunchID) })
	// The identity mirror is unique per run root, so a second Start must be refused
	// before it truncates a log artifact or overwrites the durable identity of the
	// process that is still running.
	if _, err := launcher.StartAgentRunnerProcess(ctx, first); !errors.Is(err, ErrAgentRunnerProcessDuplicate) {
		t.Fatalf("a duplicate launch must be refused, got %v", err)
	}
	other := agentRunnerPinnedCLIRequest(workspace)
	other.RequestID = "request-two"
	if _, err := launcher.StartAgentRunnerProcess(ctx, other); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
		t.Fatalf("a concurrent launch must be refused, got %v", err)
	}
	// The mirror still names the first process, which is what an operator stop uses.
	identity, err := loadAgentRunnerProcessIdentityMirror(store.RunRoot())
	if err != nil {
		t.Fatal(err)
	}
	var reported struct {
		PID        int    `json:"pid"`
		StartToken string `json:"startToken"`
	}
	if err := evidence.DecodeStrictJSON([]byte(launched.ProcessID), &reported); err != nil {
		t.Fatal(err)
	}
	if identity.PID != reported.PID || identity.StartToken != reported.StartToken {
		t.Fatalf("the mirror must still name the live process: %+v vs %+v", identity, reported)
	}
	// After a real stop the slot is free again, so the refusal is about liveness and
	// not a permanent lock.
	if _, err := launcher.StopAgentRunnerProcess(ctx, launched.LaunchID); err != nil {
		t.Fatal(err)
	}
	relaunched, err := launcher.StartAgentRunnerProcess(ctx, other)
	if err != nil {
		t.Fatalf("a released launch slot must be reusable: %v", err)
	}
	t.Cleanup(func() { _, _ = launcher.StopAgentRunnerProcess(context.Background(), relaunched.LaunchID) })
}

func TestAgentRunnerPinnedCLILauncherRefusesASecondInstanceWhileTheMirrorIsLive(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	first, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	ctx := context.Background()
	launched, err := first.StartAgentRunnerProcess(ctx, agentRunnerPinnedCLIRequest(workspace))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = first.StopAgentRunnerProcess(context.Background(), launched.LaunchID) })
	// A restart or a second launcher only sees the durable identity mirror, not the
	// first launcher's memory. The mirror is what keeps one run root to one live
	// launch, so a launch naming another process must be refused.
	second, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	other := agentRunnerPinnedCLIRequest(workspace)
	other.RequestID = "request-two"
	if _, err := second.StartAgentRunnerProcess(ctx, other); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
		t.Fatalf("a live mirrored process must refuse a second instance, got %v", err)
	}
	// A lost slot is not proof that nothing runs: the identity mirror still names the
	// live process, so a launch must be refused even with a free slot.
	if err := os.Remove(filepath.Join(runRoot, managedProcessLaunchSlotFileName)); err != nil {
		t.Fatal(err)
	}
	if _, err := second.StartAgentRunnerProcess(ctx, other); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
		t.Fatalf("a live mirrored process must refuse a launch even without a slot, got %v", err)
	}
	// Once the process is stopped the mirror is stale, so the slot is reusable.
	if _, err := first.StopAgentRunnerProcess(ctx, launched.LaunchID); err != nil {
		t.Fatal(err)
	}
	if relaunched, err := second.StartAgentRunnerProcess(ctx, other); err != nil {
		t.Fatalf("a stale mirror must not block a launch: %v", err)
	} else {
		// The second launcher owns its own handle, so it has to be retired too, which
		// also closes the log artifacts the test wants cleaned up.
		if _, err := second.StopAgentRunnerProcess(ctx, relaunched.LaunchID); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAgentRunnerPinnedCLILauncherStopGraceFollowsTheLaunch(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	// A launch that asked for a tighter grace must be stopped within it: an operator
	// stop that is more patient than the run it stops would report a stop that the
	// run side had already given up on.
	request := agentRunnerPinnedCLIRequest(t.TempDir())
	request.Grace = 200 * time.Millisecond
	launched, err := launcher.StartAgentRunnerProcess(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = launcher.StopAgentRunnerProcess(context.Background(), launched.LaunchID) })
	if grace := launcher.stopGrace(launched.LaunchID); grace != 200*time.Millisecond {
		t.Fatalf("the launch grace must win over the configuration: %v", grace)
	}
	if grace := launcher.stopGrace("launch-unknown"); grace != launcher.config.Grace {
		t.Fatalf("an unknown launch falls back to the configuration: %v", grace)
	}
	if _, err := launcher.StopAgentRunnerProcess(context.Background(), launched.LaunchID); err != nil {
		t.Fatal(err)
	}
}

func TestAgentRunnerNeedsDurableStopCoversUnsettledRuns(t *testing.T) {
	settled := guard.ProcessResult{ExitCode: -1, Termination: &guard.TerminationEvidence{Outcome: "stopped"}}
	// A natural exit and a stopped run both owe the durable stop proof, an unsettled
	// run owes a real stop, and a settled guard stop is complete by itself.
	cases := []struct {
		name   string
		watch  AgentRunnerWatchResult
		result guard.ProcessResult
		want   bool
	}{
		{name: "natural exit", watch: AgentRunnerWatchResult{}, result: guard.ProcessResult{ExitCode: 0}, want: true},
		{name: "settled guard stop", watch: AgentRunnerWatchResult{Err: guard.ErrTerminationUncertain}, result: settled, want: false},
		{name: "unsettled", watch: AgentRunnerWatchResult{Err: fmt.Errorf("%w: timeout", ErrAgentRunnerLaunchUnsettled), TimedOut: true}, result: settled, want: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := agentRunnerNeedsDurableStop(testCase.watch, testCase.result); got != testCase.want {
				t.Fatalf("agentRunnerNeedsDurableStop=%v, want %v", got, testCase.want)
			}
		})
	}
}

func TestAgentRunnerPinnedCLILauncherStopNeverReachesAnotherLaunch(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	first, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	ctx := context.Background()
	launchedFirst, err := first.StartAgentRunnerProcess(ctx, agentRunnerPinnedCLIRequest(workspace))
	if err != nil {
		t.Fatal(err)
	}
	firstIdentity, err := loadAgentRunnerProcessIdentityMirror(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	// The first launch dies without its launcher noticing, and its slot ages past the
	// freshness window, so a second launcher may reclaim the run root.
	if _, err := guard.StopProcessIdentity(ctx, firstIdentity, 2*time.Second); err != nil {
		t.Fatal(err)
	}
	slot := filepath.Join(runRoot, managedProcessLaunchSlotFileName)
	old := time.Now().Add(-2 * agentRunnerLaunchSlotFreshness)
	if err := os.Chtimes(slot, old, old); err != nil {
		t.Fatal(err)
	}
	second, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	other := agentRunnerPinnedCLIRequest(workspace)
	other.RequestID = "request-two"
	launchedSecond, err := second.StartAgentRunnerProcess(ctx, other)
	if err != nil {
		t.Fatalf("the expired slot must be reclaimable: %v", err)
	}
	t.Cleanup(func() { _, _ = second.StopAgentRunnerProcess(context.Background(), launchedSecond.LaunchID) })
	secondIdentity, err := loadAgentRunnerProcessIdentityMirror(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	// The first launcher now stops its own (already dead) launch. It may only reach the
	// process it started: the mirror names the second launch now, and the slot belongs
	// to it as well.
	if _, err := first.StopAgentRunnerProcess(ctx, launchedFirst.LaunchID); err != nil {
		t.Fatalf("stopping its own dead launch must be idempotent: %v", err)
	}
	if alive, err := guard.ProcessAlive(secondIdentity); err != nil || !alive {
		t.Fatalf("a stop may never reach another launch's process (alive=%v err=%v)", alive, err)
	}
	owner, err := os.ReadFile(slot)
	if err != nil {
		t.Fatalf("a stop may never release another launch's slot: %v", err)
	}
	if fields := strings.Fields(string(owner)); len(fields) == 0 || fields[0] != launchedSecond.LaunchID {
		t.Fatalf("the slot must still belong to the second launch, got %q", owner)
	}
}

func TestAgentRunnerPinnedCLILauncherHoldsNoSlotDuringThePrecheck(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	config := agentRunnerPinnedCLIConfig("")
	config.CheckTimeout = 10 * time.Second
	config.VersionArgs = []string{"-test.run=^TestAgentRunnerPinnedCLIHelper$", "--", "version-slow"}
	launcher, err := NewAgentRunnerPinnedCLILauncher(config, store)
	if err != nil {
		t.Fatal(err)
	}
	started := make(chan chain.AgentRunnerLaunchResult, 1)
	failed := make(chan error, 1)
	go func() {
		launched, err := launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(t.TempDir()))
		if err != nil {
			failed <- err
			return
		}
		started <- launched
	}()
	slot := filepath.Join(runRoot, managedProcessLaunchSlotFileName)
	// The version precheck runs long, and the claimed window must not include it:
	// otherwise a slow precheck would hold the run root without any process to show
	// for it.
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(slot); err == nil {
			t.Fatal("the launch slot must not be held during the version precheck")
		}
		time.Sleep(20 * time.Millisecond)
	}
	select {
	case launched := <-started:
		if _, err := os.Stat(slot); err != nil {
			t.Fatalf("the slot must be held once the launch is registered: %v", err)
		}
		if _, err := launcher.StopAgentRunnerProcess(context.Background(), launched.LaunchID); err != nil {
			t.Fatal(err)
		}
	case err := <-failed:
		t.Fatalf("the slow precheck must still launch: %v", err)
	}
}

func TestAgentRunnerPinnedCLILauncherVersionMismatchNeverSpawns(t *testing.T) {
	runRoot := t.TempDir()
	workspace := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	config := agentRunnerPinnedCLIConfig("")
	config.Version = "agent-stub-9.9.9"
	launcher, err := NewAgentRunnerPinnedCLILauncher(config, store)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(workspace)); !errors.Is(err, ErrAgentRunnerLaunchVersion) {
		t.Fatalf("a version mismatch must refuse the launch, got %v", err)
	}
	// The precheck runs before the evidence directory is even created, so its
	// absence proves the refusal happened before any spawn.
	if _, err := os.Stat(filepath.Join(store.root, agentRunnerReplayEvidenceDirectoryName, "request-one")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a refused launch must not prepare an evidence directory: %v", err)
	}
	if launcher.registry.len() != 0 {
		t.Fatal("a refused launch must not register a handle")
	}
}

func TestAgentRunnerPinnedCLILauncherRefusesAnotherExecutable(t *testing.T) {
	runRoot := t.TempDir()
	workspace := t.TempDir()
	launcher, _ := newAgentRunnerPinnedCLILauncher(t, runRoot)
	request := agentRunnerPinnedCLIRequest(workspace)
	request.Command = "some-other-cli"
	if _, err := launcher.StartAgentRunnerProcess(context.Background(), request); !errors.Is(err, ErrInvalidAgentRunnerLauncher) {
		t.Fatalf("a request naming another executable must be refused, got %v", err)
	}
	request = agentRunnerPinnedCLIRequest(workspace)
	request.RequestID = "Request One"
	if _, err := launcher.StartAgentRunnerProcess(context.Background(), request); !errors.Is(err, ErrInvalidAgentRunnerLauncher) {
		t.Fatalf("an invalid request id must be refused, got %v", err)
	}
	request = agentRunnerPinnedCLIRequest("")
	request.Dir = ""
	if _, err := launcher.StartAgentRunnerProcess(context.Background(), request); !errors.Is(err, ErrInvalidAgentRunnerLauncher) {
		t.Fatalf("a launch without a workspace root must be refused, got %v", err)
	}
}

func TestAgentRunnerPinnedCLILauncherStopsThroughPersistedIdentity(t *testing.T) {
	runRoot := t.TempDir()
	workspace := t.TempDir()
	launcher, store := newAgentRunnerPinnedCLILauncher(t, runRoot)
	result, err := launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(workspace))
	if err != nil {
		t.Fatal(err)
	}
	// A restart erases the live handle, so the test owns retiring this one: the
	// log artifacts stay open until the handle is released.
	t.Cleanup(func() { launcher.releaseLaunch(result.LaunchID) })
	identity, err := decodeAgentRunnerProcessID(result.ProcessID)
	if err != nil {
		t.Fatal(err)
	}
	// A second launcher on the same run root has no live handle: exactly the
	// restart case, where the mirror is the only durable identity.
	restarted, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfig(""), store)
	if err != nil {
		t.Fatal(err)
	}
	if restarted.registry.len() != 0 {
		t.Fatal("a restarted launcher must not inherit live handles")
	}
	proof, err := restarted.StopAgentRunnerProcess(context.Background(), result.LaunchID)
	if err != nil {
		t.Fatalf("stop through the mirror failed: %v (%+v)", err, proof)
	}
	if proof.Outcome != "stopped" {
		t.Fatalf("expected a proven stop, got %+v", proof)
	}
	if alive, err := guard.ProcessAlive(identity); err == nil && alive {
		t.Fatal("the restarted stop must reach the original process")
	}
}

func TestAgentRunnerPinnedCLILauncherUnprovenIdentityStopsTheSpawn(t *testing.T) {
	runRoot := t.TempDir()
	workspace := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfig(""), store)
	if err != nil {
		t.Fatal(err)
	}
	// Make the identity publication fail through the writer seam, which works on
	// every platform: the spawn must then be stopped again instead of being reported
	// as a launch whose identity an operator stop could never find.
	original := agentRunnerIdentityMirrorWriter
	agentRunnerIdentityMirrorWriter = func(string, guard.ProcessIdentity) error {
		return errors.New("identity publication unavailable")
	}
	t.Cleanup(func() { agentRunnerIdentityMirrorWriter = original })
	if _, err := launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(workspace)); !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) {
		t.Fatalf("an unpublished identity must fail closed, got %v", err)
	} else if !strings.Contains(err.Error(), "spawn stop outcome:") {
		t.Fatalf("the refusal must report what happened to the abandoned spawn, got %v", err)
	}
	if launcher.registry.len() != 0 {
		t.Fatal("an unproven spawn must not leave a handle behind")
	}
	// The failed launch must release its slot, otherwise the run root would stay
	// claimed with nothing in it.
	agentRunnerIdentityMirrorWriter = original
	launched, err := launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(workspace))
	if err != nil {
		t.Fatalf("a released slot must allow a new launch: %v", err)
	}
	if _, err := launcher.StopAgentRunnerProcess(context.Background(), launched.LaunchID); err != nil {
		t.Fatal(err)
	}
}

func TestAgentRunnerPinnedCLILauncherStopsItsOwnProcessWhenTheMirrorIsUnreadable(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	launched, err := launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(t.TempDir()))
	if err != nil {
		t.Fatal(err)
	}
	identity, err := loadAgentRunnerProcessIdentityMirror(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, stopErr := guard.StopProcessIdentity(context.Background(), identity, 2*time.Second); stopErr != nil {
			t.Logf("leftover process %d could not be stopped: %v", identity.PID, stopErr)
		}
	})
	// Corrupt the mirror after the launch started. The launch this process holds is
	// stopped through its own identity, so a damaged durable record cannot make a live
	// process unstoppable - and the launch, its logs and its slot are all retired.
	mirrorPath := filepath.Join(runRoot, managedProcessIdentityFileName)
	if err := os.WriteFile(mirrorPath, []byte("not-an-identity"), 0o600); err != nil {
		t.Fatal(err)
	}
	proof, err := launcher.StopAgentRunnerProcess(context.Background(), launched.LaunchID)
	if err != nil {
		t.Fatalf("a held launch must be stoppable through its own identity: %v", err)
	}
	if proof.Outcome != "stopped" {
		t.Fatalf("the stop must be proven, got %+v", proof)
	}
	if alive, err := guard.ProcessAlive(identity); err == nil && alive {
		t.Fatal("the held process must be gone")
	}
	if launcher.registry.len() != 0 {
		t.Fatal("the launch must be retired")
	}
	if _, err := os.Stat(filepath.Join(runRoot, managedProcessLaunchSlotFileName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the launch slot must be released with the launch: %v", err)
	}
	// A launcher that did not start the process has only the mirror, so an unreadable
	// mirror must fail closed there instead of guessing.
	restarted, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := restarted.StopAgentRunnerProcess(context.Background(), launched.LaunchID); err == nil {
		t.Fatal("a restarted launcher may not stop without a readable identity")
	}
}

func TestAgentRunnerPinnedCLILauncherRefusesToStopAForeignLaunch(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	launched, err := launcher.StartAgentRunnerProcess(ctx, agentRunnerPinnedCLIRequest(t.TempDir()))
	if err != nil {
		t.Fatal(err)
	}
	identity, err := loadAgentRunnerProcessIdentityMirror(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = launcher.StopAgentRunnerProcess(context.Background(), launched.LaunchID) })
	// A launcher may only stop what it started: an unknown or foreign launch id must
	// not resolve to the process its own live launch owns.
	if _, err := launcher.StopAgentRunnerProcess(ctx, "launch-foreign"); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
		t.Fatalf("a foreign launch id must be refused, got %v", err)
	}
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("the foreign request must not touch the live process: %v", err)
	}
	// The same holds for a launcher that holds nothing: the run root names its owner,
	// so a stale id cannot make it stop the process that owner is running either.
	restarted, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := restarted.StopAgentRunnerProcess(ctx, "launch-foreign"); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
		t.Fatalf("a foreign id must not stop the run root owner, got %v", err)
	}
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("the owner must still be running: %v", err)
	}
	// Without the slot the ownership record is gone, but the launcher still holds its
	// own live launch: that alone must keep a foreign request away from the process.
	if err := os.Remove(filepath.Join(runRoot, managedProcessLaunchSlotFileName)); err != nil {
		t.Fatal(err)
	}
	if _, err := launcher.StopAgentRunnerProcess(ctx, "launch-foreign"); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
		t.Fatalf("a held launch must shield its process from a foreign request, got %v", err)
	}
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("the held process must still be running: %v", err)
	}
}

func TestAgentRunnerPinnedCLILauncherRefusesStopWhenTheSlotIsUnreadable(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	owner, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	launched, err := owner.StartAgentRunnerProcess(ctx, agentRunnerPinnedCLIRequest(t.TempDir()))
	if err != nil {
		t.Fatal(err)
	}
	identity, err := loadAgentRunnerProcessIdentityMirror(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, stopErr := guard.StopProcessIdentity(context.Background(), identity, 2*time.Second); stopErr != nil {
			t.Logf("leftover process %d could not be stopped: %v", identity.PID, stopErr)
		}
	})
	// A launcher that did not start the process has only the slot and the mirror. An
	// unreadable slot is not proof that nothing runs - a live launch always holds its
	// slot - so the stop must fail closed instead of falling through to the mirror.
	slot := filepath.Join(runRoot, managedProcessLaunchSlotFileName)
	if err := os.Remove(slot); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(slot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(slot, "blocker"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	restarted, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	for _, launchID := range []string{launched.LaunchID, "launch-foreign"} {
		if _, err := restarted.StopAgentRunnerProcess(ctx, launchID); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
			t.Fatalf("an unreadable slot must refuse the stop of %s, got %v", launchID, err)
		}
	}
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("the live process must not be touched: %v", err)
	}
	// An empty slot names nobody, which is equally no proof of ownership.
	if err := os.RemoveAll(slot); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(slot, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := restarted.StopAgentRunnerProcess(ctx, "launch-foreign"); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
		t.Fatalf("an ownerless slot must refuse a foreign stop, got %v", err)
	}
	// A slot that names nobody is no proof for the owner either: without the launch
	// record the owner's own id buys nothing, and the refusal must not have touched
	// the process on the way out.
	if _, err := restarted.StopAgentRunnerProcess(ctx, launched.LaunchID); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
		t.Fatalf("an ownerless slot must refuse the owner's stop too, got %v", err)
	}
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("a refused stop must leave the owner's process running: %v", err)
	}
	// Without a slot the mirror itself decides: a live process may not be stopped by a
	// request that does not own it, while a dead one may be reached (already-stopped).
	if err := os.Remove(slot); err != nil {
		t.Fatal(err)
	}
	// A missing slot with a live process refuses every id, the owner's included: the
	// slot is what proves ownership, and without it nobody can prove any.
	for _, launchID := range []string{"launch-foreign", launched.LaunchID} {
		if _, err := restarted.StopAgentRunnerProcess(ctx, launchID); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
			t.Fatalf("a missing slot with a live process must refuse the stop of %s, got %v", launchID, err)
		}
	}
	// Refusing twice must not have been a stop in disguise: the process the mirror
	// still names has to be alive when the dead branch is exercised below.
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("a refused stop must leave the mirrored process running: %v", err)
	}
	if _, err := guard.StopProcessIdentity(ctx, identity, 2*time.Second); err != nil {
		t.Fatal(err)
	}
	proof, err := restarted.StopAgentRunnerProcess(ctx, "launch-foreign")
	if err != nil {
		t.Fatalf("a dead process must be reachable so the stop can report already-stopped: %v", err)
	}
	// A stop that nothing owns must still be provable, and the idempotent marker must
	// be the only action the proof records: the guard keeps the outcome a stop, and a
	// second action (a kill) would mean the dead branch did touch a process.
	if proof.Outcome != "stopped" || len(proof.Actions) != 1 || proof.Actions[0] != "already-stopped" {
		t.Fatalf("a stop that nothing owns must report exactly one already-stopped proof, got %+v", proof)
	}
	// The owner still holds its handle (the test killed the process behind its back),
	// so it retires the launch: that closes the log artifacts this test must not leak.
	if _, err := owner.StopAgentRunnerProcess(ctx, launched.LaunchID); err != nil {
		t.Fatal(err)
	}
}

func TestAgentRunnerPinnedCLILauncherConcurrentStartsHaveOneWinner(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	const contenders = 8
	start := make(chan struct{})
	results := make(chan chain.AgentRunnerLaunchResult, contenders)
	errs := make(chan error, contenders)
	for index := 0; index < contenders; index++ {
		request := agentRunnerPinnedCLIRequest(workspace)
		request.RequestID = fmt.Sprintf("request-%02d", index)
		go func() {
			<-start
			result, err := launcher.StartAgentRunnerProcess(context.Background(), request)
			if err != nil {
				errs <- err
				return
			}
			results <- result
		}()
	}
	close(start)
	winners := 0
	losers := 0
	for index := 0; index < contenders; index++ {
		select {
		case result := <-results:
			winners++
			t.Cleanup(func() { _, _ = launcher.StopAgentRunnerProcess(context.Background(), result.LaunchID) })
		case err := <-errs:
			losers++
			// A loser may only see the two exclusivity sentinels: nothing else may be
			// reachable while the slot is contended.
			if !errors.Is(err, ErrAgentRunnerLaunchConflict) && !errors.Is(err, ErrAgentRunnerProcessDuplicate) {
				t.Fatalf("an unexpected contender error: %v", err)
			}
		}
	}
	// The slot is claimed atomically, so exactly one launch may win even when the
	// contenders start together.
	if winners != 1 {
		t.Fatalf("exactly one contender may win the launch slot, got %d (losers %d)", winners, losers)
	}
	if launcher.registry.len() != 1 {
		t.Fatalf("one winner means one registered handle, got %d", launcher.registry.len())
	}
	// The winner owns the run root until it is retired: the slot must exist while it
	// runs, otherwise the window between the spawn and the identity publication would
	// be open to another process.
	slot := filepath.Join(runRoot, managedProcessLaunchSlotFileName)
	if _, err := os.Stat(slot); err != nil {
		t.Fatalf("a live launch must hold the launch slot: %v", err)
	}
	identity, err := loadAgentRunnerProcessIdentityMirror(store.RunRoot())
	if err != nil {
		t.Fatal(err)
	}
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("the mirror must name the winning process (%+v, alive=%v err=%v)", identity, alive, err)
	}
	if _, err := launcher.StopAgentRunnerProcess(context.Background(), identityLaunchID(t, launcher)); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(slot); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("retiring the launch must release the slot: %v", err)
	}
}

// identityLaunchID returns the only registered launch id, which is what the tests
// above have already proven to be unique.
func identityLaunchID(t *testing.T, launcher *AgentRunnerPinnedCLILauncher) string {
	t.Helper()
	launchID, found := launcher.registry.live()
	if !found {
		t.Fatal("no live launch to retire")
	}
	return launchID
}

func TestAgentRunnerPinnedCLILauncherRefusesAnUnreadableIdentityMirror(t *testing.T) {
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfig(""), store)
	if err != nil {
		t.Fatal(err)
	}
	slot := filepath.Join(runRoot, managedProcessLaunchSlotFileName)
	mirrorPath := filepath.Join(runRoot, managedProcessIdentityFileName)
	backdate := func(path string) {
		t.Helper()
		old := time.Now().Add(-2 * agentRunnerLaunchSlotFreshness)
		if err := os.Chtimes(path, old, old); err != nil {
			t.Fatal(err)
		}
	}
	// An expired slot whose mirror cannot be examined is not proof that nothing is
	// running: the launch is refused before anything is spawned. The slot is
	// backdated so the freshness rule cannot be the reason for the refusal.
	if err := os.MkdirAll(mirrorPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(mirrorPath, "blocker"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(slot, []byte("launch-slot\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	backdate(slot)
	if _, err := launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(t.TempDir())); !errors.Is(err, ErrAgentRunnerLaunchConflict) {
		t.Fatalf("an unreadable claimed mirror must refuse the launch, got %v", err)
	}
	if launcher.registry.len() != 0 {
		t.Fatal("the refusal must happen before anything is spawned")
	}
	// The same expired slot with a mirror that names a dead process is stale state:
	// the slot is reclaimed and a fresh readable identity is published.
	if err := os.RemoveAll(mirrorPath); err != nil {
		t.Fatal(err)
	}
	finished := exec.Command(os.Args[0], "-test.run=^$")
	if err := finished.Start(); err != nil {
		t.Fatal(err)
	}
	identity, err := guard.InspectProcess(finished.Process.Pid)
	if err != nil {
		t.Fatal(err)
	}
	_ = finished.Wait()
	if err := writeAgentRunnerProcessIdentityMirror(runRoot, identity); err != nil {
		t.Fatal(err)
	}
	launched, err := launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(t.TempDir()))
	if err != nil {
		t.Fatalf("an expired slot with a dead mirror must be reclaimable: %v", err)
	}
	published, err := loadAgentRunnerProcessIdentityMirror(runRoot)
	if err != nil {
		t.Fatalf("the launch must publish a readable identity: %v", err)
	}
	if published.PID == identity.PID {
		t.Fatal("the reclaimed slot must be backed by a fresh identity")
	}
	// Without a claim the stale mirror is simply replaced: nothing proves a process
	// is running once the slot is gone.
	if _, err := launcher.StopAgentRunnerProcess(context.Background(), launched.LaunchID); err != nil {
		t.Fatal(err)
	}
	launched, err = launcher.StartAgentRunnerProcess(context.Background(), agentRunnerPinnedCLIRequest(t.TempDir()))
	if err != nil {
		t.Fatalf("an unclaimed dead mirror must not block a launch: %v", err)
	}
	if _, err := launcher.StopAgentRunnerProcess(context.Background(), launched.LaunchID); err != nil {
		t.Fatal(err)
	}
}

func TestAgentRunnerPinnedCLILauncherRefusesUnusableConfig(t *testing.T) {
	store := replayStoreForRunRoot(t, t.TempDir())
	for name, config := range map[string]AgentRunnerPinnedCLIConfig{
		"missing executable": {Version: "v1"},
		"missing version":    {Executable: os.Args[0]},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := NewAgentRunnerPinnedCLILauncher(config, store); !errors.Is(err, ErrInvalidAgentRunnerLauncher) {
				t.Fatalf("expected a fail-closed config rejection, got %v", err)
			}
		})
	}
	if _, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfig(""), nil); !errors.Is(err, ErrInvalidAgentRunnerLauncher) {
		t.Fatalf("a nil store must be refused, got %v", err)
	}
}

func TestAgentRunnerPinnedCLIChildEnvAllowlist(t *testing.T) {
	launcher, _ := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfig(""), replayStoreForRunRoot(t, t.TempDir()))
	env := launcher.childEnv([]string{"PATH=one", "SECRET=two"})
	if len(env) != 2 {
		t.Fatalf("an empty allowlist must pass the request environment through: %v", env)
	}
	launcher.config.EnvAllow = []string{"path"}
	filtered := launcher.childEnv([]string{"PATH=one", "SECRET=two"})
	if len(filtered) != 1 || filtered[0] != "PATH=one" {
		t.Fatalf("the allowlist must filter the child environment: %v", filtered)
	}
	if inherited := launcher.childEnv(nil); len(inherited) == 0 {
		t.Fatal("an empty request environment falls back to the inherited one")
	}
}

func TestAgentRunnerPinnedCLILauncherStopWithoutIdentityFailsClosed(t *testing.T) {
	launcher, _ := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfig(""), replayStoreForRunRoot(t, t.TempDir()))
	if _, err := launcher.StopAgentRunnerProcess(context.Background(), "launch-request-missing"); !errors.Is(err, ErrAgentRunnerProcessIdentity) {
		t.Fatalf("stopping an unknown launch without a mirror must fail closed, got %v", err)
	}
}

func parseAgentRunnerPID(t *testing.T, value string) int {
	t.Helper()
	pid := 0
	if _, err := fmt.Sscanf(strings.TrimSpace(value), "%d", &pid); err != nil {
		t.Fatalf("unparseable pid %q: %v", value, err)
	}
	return pid
}

// waitForAgentRunnerArtifact waits for an artifact the spawned CLI writes.
func waitForAgentRunnerArtifact(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("the CLI never published %s", path)
}
