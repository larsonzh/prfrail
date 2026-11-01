package console

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

// zeroModelAgentCLIVersion is the version the re-exec helper reports and the
// launcher pins. It never reaches a real agent CLI: the helper is this test binary.
const zeroModelAgentCLIVersion = "agent-stub-zero-model-0.1.0"

// TestZeroModelAgentCLIHelper is the deterministic pinned CLI of the concurrency
// leg. It only acts when the test binary is re-executed with "--" plus a mode, so a
// normal test run returns immediately.
func TestZeroModelAgentCLIHelper(t *testing.T) {
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
		fmt.Fprintln(os.Stdout, zeroModelAgentCLIVersion)
		os.Exit(0)
	case "run":
		// A real agent run keeps working until the launcher stops it. This sleep is
		// the child's own lifetime, not test synchronisation: the test itself never
		// sleeps and never polls.
		time.Sleep(time.Minute)
		os.Exit(0)
	}
	os.Exit(2)
}

func zeroModelLaunchRequest(requestID, workspaceRoot string) chain.AgentRunnerLaunchRequest {
	return chain.AgentRunnerLaunchRequest{
		RequestID:     requestID,
		RunID:         "run-zero-model-concurrent",
		AdapterID:     "zero-model-adapter",
		Dir:           workspaceRoot,
		WorkspaceRoot: workspaceRoot,
	}
}

// zeroModelProcessIdentity decodes the opaque process identity a launch result
// carries into the guard identity an operator stop would use.
func zeroModelProcessIdentity(processID string) (guard.ProcessIdentity, error) {
	var identity guard.ProcessIdentity
	if err := evidence.DecodeStrictJSON([]byte(processID), &identity); err != nil {
		return guard.ProcessIdentity{}, fmt.Errorf("decode process identity %q: %w", processID, err)
	}
	return identity, nil
}

// TestZeroModelConcurrentStartOneWinner proves the launch slot is exclusive: eight
// contenders on one launcher and one run root produce exactly one winner, and the
// loser verdicts are only the two exclusivity sentinels.
func TestZeroModelConcurrentStartOneWinner(t *testing.T) {
	runRoot := t.TempDir()
	// The replay store binds to the durable run root, which the run event log owns.
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(runRoot)); err != nil {
		t.Fatalf("create the run-root marker: %v", err)
	}
	store, err := adapters.NewAgentRunnerReplayStoreForRun(runRoot, "run-zero-model-concurrent")
	if err != nil {
		t.Fatalf("bind the replay store to the run root: %v", err)
	}
	launcher, err := adapters.NewAgentRunnerPinnedCLILauncher(adapters.AgentRunnerPinnedCLIConfig{
		Executable:   os.Args[0],
		Version:      zeroModelAgentCLIVersion,
		VersionArgs:  []string{"-test.run=^TestZeroModelAgentCLIHelper$", "--", "version"},
		RunArgs:      []string{"-test.run=^TestZeroModelAgentCLIHelper$", "--", "run"},
		CheckTimeout: 20 * time.Second,
		Grace:        time.Second,
	}, store)
	if err != nil {
		t.Fatalf("new pinned CLI launcher: %v", err)
	}

	workspace := t.TempDir()
	const contenders = 8
	type outcome struct {
		result chain.AgentRunnerLaunchResult
		err    error
	}
	release := make(chan struct{})
	outcomes := make(chan outcome, contenders)
	for index := 0; index < contenders; index++ {
		request := zeroModelLaunchRequest(fmt.Sprintf("request-%02d", index), workspace)
		go func() {
			<-release
			result, startErr := launcher.StartAgentRunnerProcess(context.Background(), request)
			outcomes <- outcome{result: result, err: startErr}
		}()
	}
	close(release)

	winner := chain.AgentRunnerLaunchResult{}
	winners, losers := 0, 0
	for index := 0; index < contenders; index++ {
		got := <-outcomes
		if got.err != nil {
			losers++
			// A loser may only see the two exclusivity sentinels: nothing else is
			// reachable while the slot is contended.
			if !errors.Is(got.err, adapters.ErrAgentRunnerLaunchConflict) && !errors.Is(got.err, adapters.ErrAgentRunnerProcessDuplicate) {
				t.Fatalf("a losing contender may only see the exclusivity sentinels, got %v", got.err)
			}
			continue
		}
		winners++
		winner = got.result
	}
	if winners != 1 || losers != contenders-1 {
		t.Fatalf("exactly one contender may win the launch slot, got %d winners and %d losers", winners, losers)
	}
	t.Cleanup(func() { _, _ = launcher.StopAgentRunnerProcess(context.Background(), winner.LaunchID) })

	identity, err := zeroModelProcessIdentity(winner.ProcessID)
	if err != nil {
		t.Fatalf("the winner must report a decodable process identity: %v", err)
	}
	if alive, err := guard.ProcessAlive(identity); err != nil || !alive {
		t.Fatalf("the winning launch must own a live process: %+v alive=%v err=%v", identity, alive, err)
	}

	// While the winner is live the run root is owned: one more start must conflict.
	if _, err := launcher.StartAgentRunnerProcess(context.Background(), zeroModelLaunchRequest("request-late", workspace)); !errors.Is(err, adapters.ErrAgentRunnerLaunchConflict) {
		t.Fatalf("a start while the winner is live must be refused with the launch conflict, got %v", err)
	}

	if _, err := launcher.StopAgentRunnerProcess(context.Background(), winner.LaunchID); err != nil {
		t.Fatalf("stop the winning launch: %v", err)
	}
	if alive, err := guard.ProcessAlive(identity); err != nil {
		t.Fatalf("inspect the stopped identity: %v", err)
	} else if alive {
		t.Fatalf("the stopped winner must be gone: %+v", identity)
	}

	// A released slot must be reusable, otherwise the refusal above was a lock.
	relaunched, err := launcher.StartAgentRunnerProcess(context.Background(), zeroModelLaunchRequest("request-released", workspace))
	if err != nil {
		t.Fatalf("a released launch slot must be reusable: %v", err)
	}
	t.Cleanup(func() { _, _ = launcher.StopAgentRunnerProcess(context.Background(), relaunched.LaunchID) })
	if relaunched.LaunchID == winner.LaunchID {
		t.Fatalf("the relaunch must be a distinct launch, got %q twice", relaunched.LaunchID)
	}
}
