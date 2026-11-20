package gates

// P83 FIX ROUND 6 counterexample (item T4a).
//
// T4a: the workspace root is re-proven against the bound directory handle on
// EVERY attempt, so a retry whose root was swapped refuses before the attempt's
// hook runs and leaves no side effect in the replacement directory.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// p83Fix6SwappingExecutor performs the root swap inside the first attempt's hook
// and records how many attempts reached the hook at all.
type p83Fix6SwappingExecutor struct {
	calls   int
	onFirst func()
	results []Execution
}

func (executor *p83Fix6SwappingExecutor) Capabilities() Capabilities { return allCapabilities() }

func (executor *p83Fix6SwappingExecutor) Execute(context.Context, Prepared, Hook) Execution {
	executor.calls++
	if executor.calls == 1 && executor.onFirst != nil {
		executor.onFirst()
	}
	index := executor.calls - 1
	if index >= len(executor.results) {
		index = len(executor.results) - 1
	}
	return executor.results[index]
}

// TestP83Fix6T4aRetryAfterRootSwapRefusesBeforeTheHook is the T4a counterexample:
// the first attempt fails and the policy retries, and between the two attempts
// the declared root is renamed away and a replacement directory takes its name.
// The retry must refuse with the root-moved error before the hook runs, and the
// replacement directory must stay empty.
//
// Mutation check: deleting the per-attempt VerifyBoundDir/RootPath block in
// Runner.Run lets the second attempt execute in the replacement directory
// (calls==2) and turns this red.
func TestP83Fix6T4aRetryAfterRootSwapRefusesBeforeTheHook(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "report.txt"), []byte("report"), 0o600); err != nil {
		t.Fatal(err)
	}
	moved := root + "-moved-away"
	request := validRequest(root)
	request.Hook.ArtifactGlobs = []string{"report.txt"}
	request.Hook.OnFail = FailurePolicy{Action: "retry", MaxAttempts: 2}
	refreshHookDefinition(t, &request)

	failed := successfulExecution()
	exitCode := 2
	failed.ExitCode = &exitCode
	executor := &p83Fix6SwappingExecutor{
		results: []Execution{failed, successfulExecution()},
		onFirst: func() {
			if err := os.Rename(root, moved); err != nil {
				t.Errorf("swap the workspace root: %v", err)
				return
			}
			if err := os.Mkdir(root, 0o755); err != nil {
				t.Errorf("materialise the replacement root: %v", err)
			}
		},
	}
	store := &memoryResults{}

	record, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if !errors.Is(err, ErrWorkspaceRootMoved) {
		t.Fatalf("a retry whose root was swapped must refuse with the root-moved error, got record=%+v err=%v", record, err)
	}
	if executor.calls != 1 {
		t.Fatalf("the retry must not reach the hook, got %d attempts", executor.calls)
	}
	if len(store.records) != 1 {
		t.Fatalf("only the first attempt may be recorded, got %d records", len(store.records))
	}
	if got := store.records[0].Result.ExecutionAttempt; got != 1 {
		t.Fatalf("the recorded attempt must be the one that ran, got %d", got)
	}
	// Zero side effects for the refused attempt: the replacement directory is the
	// directory the second attempt would have executed in, so it must be untouched.
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		names := make([]string, 0, len(entries))
		for _, entry := range entries {
			names = append(names, entry.Name())
		}
		t.Fatalf("the refused attempt must leave the replacement root empty, found %v", names)
	}
}
