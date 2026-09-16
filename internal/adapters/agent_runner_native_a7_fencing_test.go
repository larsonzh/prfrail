//go:build a7native

package adapters

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Each experiment spawns a child that runs only its own helper half. A child that
// is handed a filter covering several tests runs the other experiments as well,
// which is exactly how the first version of the completion experiment ended up
// measuring the request path.
const (
	a7RequestHelperFilter    = "^TestA7NativePublishCrashMatrix$"
	a7CompletionHelperFilter = "^TestA7NativeCompletionPublicationCrashMatrix$"
	a7FencingHelperFilter    = "^TestA7NativeConcurrentProcessFencing$"
)

// The pre-review required two experiments the first implementation lacked: a
// crash matrix over the completion publication (not just the request) and a real
// multi-process fencing experiment (the in-process race only proved the lock and
// the no-replace collision handling). Both reuse the existing publish-stage seam.

const (
	a7CompletionModeEnv = "PROOFRAIL_A7_COMPLETION"
	a7ConcurrentModeEnv = "PROOFRAIL_A7_CONCURRENT"
	a7BarrierEnv        = "PROOFRAIL_A7_BARRIER"
)

func a7SpawnHelper(t *testing.T, runFilter string, environment []string) (*exec.Cmd, *bytes.Buffer) {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run="+runFilter)
	cmd.Env = append(os.Environ(), environment...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	if err := cmd.Start(); err != nil {
		t.Fatalf("start helper: %v", err)
	}
	return cmd, &buffer
}

func a7WaitForFile(t *testing.T, path string, buffer *bytes.Buffer) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("helper never reached the agreed point %s: %s", filepath.Base(path), buffer.String())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestA7NativeCompletionPublicationCrashMatrix kills a child while it publishes a
// completion record and classifies the replay root afterwards. The request must
// never be lost and a completion must never appear without its request: the
// protocol writes R first and publishes C through the same no-replace primitive,
// so the legal outcomes are "request only" and "request plus completion".
func TestA7NativeCompletionPublicationCrashMatrix(t *testing.T) {
	if os.Getenv(a7CompletionModeEnv) != "" {
		a7RunCompletionHelper(t)
		return
	}
	rounds := a7ExperimentRounds(t)
	observed := map[string]map[AgentRunnerReplayState]int{}
	for _, stage := range a7PublishStages() {
		observed[stage] = map[AgentRunnerReplayState]int{}
		for round := 0; round < rounds; round++ {
			root := t.TempDir()
			ready := filepath.Join(t.TempDir(), "ready")
			cmd, buffer := a7SpawnHelper(t, a7CompletionHelperFilter, []string{
				a7CompletionModeEnv + "=1",
				a7StageEnv + "=" + stage,
				a7RootEnv + "=" + root,
				a7ReadyEnv + "=" + ready,
			})
			a7WaitForFile(t, ready, buffer)
			if err := cmd.Process.Kill(); err != nil {
				t.Fatalf("stage %s round %d: kill completion helper: %v", stage, round, err)
			}
			_ = cmd.Wait()

			restarted := a7StoreAt(t, root)
			state, err := restarted.State(a7RequestID)
			if err != nil {
				t.Fatalf("stage %s round %d: an orphan completion must not be produced, got %v", stage, round, err)
			}
			switch state {
			case AgentRunnerReplayStateDispatchedUnknown, AgentRunnerReplayStateTerminalReceiptPresent:
			default:
				t.Fatalf("stage %s round %d: illegal state %s after a completion crash", stage, round, state)
			}
			visible := a7StageLeavesRecordVisible(stage)
			if want := AgentRunnerReplayStateTerminalReceiptPresent; visible && state != want {
				t.Fatalf("stage %s round %d: a visible completion must be %s, got %s", stage, round, want, state)
			}
			if !visible && state != AgentRunnerReplayStateDispatchedUnknown {
				t.Fatalf("stage %s round %d: the request must survive without its completion, got %s", stage, round, state)
			}
			observed[stage][state]++
		}
	}
	for _, stage := range a7PublishStages() {
		t.Logf("A7 completion kill matrix stage=%-14s outcomes=%v rounds=%d", stage, observed[stage], rounds)
	}
}

// a7IsCompletionPath keeps the completion experiment inside the completion
// publication. Publishing a completion first re-publishes the request through the
// same primitive (the collision decides the winner), so a hook that matched only
// the stage name would stop the child before the completion was ever attempted.
func a7IsCompletionPath(path string) bool {
	separator := string(filepath.Separator)
	return strings.Contains(path, separator+"completions"+separator)
}

// a7RunCompletionHelper is the child half of the completion crash matrix: it
// publishes the request for real, then stops existing inside the completion
// publication at the named stage.
func a7RunCompletionHelper(t *testing.T) {
	t.Helper()
	stage := os.Getenv(a7StageEnv)
	root := os.Getenv(a7RootEnv)
	ready := os.Getenv(a7ReadyEnv)
	store := a7StoreAt(t, root)
	request := replayStoreRequestRecord(t, a7RequestID)
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatalf("completion helper could not publish the request: %v", err)
	}
	replayStorePublishStageHook = func(current string, path string) {
		if current != stage || !a7IsCompletionPath(path) {
			return
		}
		if err := os.WriteFile(ready, []byte(current), 0o600); err != nil {
			panic(err)
		}
		time.Sleep(a7HelperBlockDuration)
	}
	completion := replayStoreCompletionRecord(t, request, "completion-one")
	if _, err := store.RecordCompletion(request, completion); err != nil {
		t.Fatalf("completion helper publish failed: %v", err)
	}
	panic("completion helper returned from a publish it was told to stop inside")
}

// TestA7NativeConcurrentProcessFencing lets two real processes publish the same
// request at the same instant. The in-process race only proved the mutex and the
// no-replace collision path; the slice's fencing item is about separate processes
// sharing one replay root, which is what this experiment isolates.
func TestA7NativeConcurrentProcessFencing(t *testing.T) {
	if os.Getenv(a7ConcurrentModeEnv) != "" {
		a7RunConcurrentHelper(t)
		return
	}
	rounds := a7ExperimentRounds(t)
	for round := 0; round < rounds; round++ {
		root := t.TempDir()
		barrier := filepath.Join(t.TempDir(), "go")
		type helper struct {
			cmd    *exec.Cmd
			buffer *bytes.Buffer
			ready  string
		}
		helperList := make([]helper, 0, 2)
		for index := 0; index < 2; index++ {
			ready := filepath.Join(t.TempDir(), "ready")
			cmd, buffer := a7SpawnHelper(t, a7FencingHelperFilter, []string{
				a7ConcurrentModeEnv + "=1",
				a7RootEnv + "=" + root,
				a7BarrierEnv + "=" + barrier,
				a7ReadyEnv + "=" + ready,
			})
			helperList = append(helperList, helper{cmd: cmd, buffer: buffer, ready: ready})
		}
		// Two-way rendezvous: both children park immediately before their publish and
		// are released by one write, so the race is started together rather than in
		// turn. A one-way start gate can pass with the two processes running serially.
		for _, item := range helperList {
			a7WaitForFile(t, item.ready, item.buffer)
		}
		if err := os.WriteFile(barrier, []byte("go"), 0o600); err != nil {
			t.Fatalf("round %d: release the fencing helpers: %v", round, err)
		}
		first := 0
		unknown := 0
		for _, item := range helperList {
			if err := item.cmd.Wait(); err != nil {
				t.Fatalf("round %d: helper failed: %s", round, item.buffer.String())
			}
			decision := ""
			for _, line := range strings.Split(item.buffer.String(), "\n") {
				if strings.HasPrefix(line, "A7-DECISION ") {
					decision = strings.TrimSpace(strings.TrimPrefix(line, "A7-DECISION "))
				}
			}
			switch decision {
			case string(AgentRunnerReplayDecisionFirstDispatch):
				first++
			case string(AgentRunnerReplayDecisionUnknownBlock):
				unknown++
			default:
				t.Fatalf("round %d: unexpected decision %q in %s", round, decision, item.buffer.String())
			}
		}
		if first != 1 || unknown != 1 {
			t.Fatalf("round %d: two processes produced first=%d unknown-block=%d, want exactly one of each", round, first, unknown)
		}
		restarted := a7StoreAt(t, root)
		state, err := restarted.State(a7RequestID)
		if err != nil {
			t.Fatalf("round %d: fencing left the store unreadable: %v", round, err)
		}
		if state != AgentRunnerReplayStateDispatchedUnknown {
			t.Fatalf("round %d: state = %s, want %s", round, state, AgentRunnerReplayStateDispatchedUnknown)
		}
	}
	t.Logf("A7 multi-process fencing: %d rounds, two-way rendezvous, every round had exactly one first dispatch", rounds)
}

// a7RunConcurrentHelper is the child half of the fencing experiment: it parks as
// soon as it is ready to publish and waits for the parent's single release, so
// both siblings start their publish from the same rendezvous.
func a7RunConcurrentHelper(t *testing.T) {
	t.Helper()
	root := os.Getenv(a7RootEnv)
	barrier := os.Getenv(a7BarrierEnv)
	ready := os.Getenv(a7ReadyEnv)
	store := a7StoreAt(t, root)
	request := replayStoreRequestRecord(t, a7RequestID)
	if err := os.WriteFile(ready, []byte("ready"), 0o600); err != nil {
		t.Fatalf("fencing helper could not announce readiness: %v", err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for {
		if _, err := os.Stat(barrier); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("fencing helper never saw the barrier")
		}
		time.Sleep(2 * time.Millisecond)
	}
	decision, err := store.RecordRequest(request)
	if err != nil {
		t.Fatalf("fencing helper publish failed: %v", err)
	}
	// Printed on stdout, not through t.Logf: a passing child suppresses t.Logf
	// output and the parent would see no decision at all.
	fmt.Printf("A7-DECISION %s\n", decision)
}
