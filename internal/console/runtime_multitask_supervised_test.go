package console

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
)

// P50 CONSOLE-MULTITASK-ASSEMBLY -- supervised mirror of the multi-task harness.
//
// The same serial contract is proven for the supervised-black-box assembly: two
// tasks park, confirm, open, return and pass in definition order, each task
// freezes its own five supervised facts, and every durable artifact is written
// under that task's own directory. The visible channel is the deterministic
// scripted stub, so no network or external process is involved.
func TestSupervisedMultiTaskSerialHarness(t *testing.T) {
	harness := newSupervisedHarness(t, []string{
		adapters.VisibleStubAccepted,
		adapters.VisibleStubAccepted,
	}, false, 60000)
	harness.options.Definition = multiTaskSupervisedDefinition()
	ctx := context.Background()

	// Task one: park, confirm, open.
	summary, err := ExecuteSupervisedRun(ctx, harness.options)
	if err != nil {
		t.Fatalf("task one must park, got %v", err)
	}
	if summary.ChainState != "PAUSED" {
		t.Fatalf("task one must park the run, got %s", summary.ChainState)
	}
	harness.respond(supervisedConfirmationResponse)
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("task one must open its work window, got %v", err)
	}
	harness.writeTaskArtifact(multiTaskAgentTaskOne, "artifact.txt", "task one work\n")

	// Task one return drives the tail chain and parks task two.
	harness.advance(time.Minute)
	summary, err = ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if err != nil {
		t.Fatalf("task one must complete its return, got %v", err)
	}
	if summary.ChainState != "PAUSED" {
		t.Fatalf("the run must hand over to task two after task one, got %s", summary.ChainState)
	}

	// Task two: park, confirm, open, return.
	harness.respond(supervisedConfirmationResponse)
	if _, err := ExecuteSupervisedRun(ctx, harness.options); err != nil {
		t.Fatalf("task two must open its work window, got %v", err)
	}
	harness.writeTaskArtifact(multiTaskAgentTaskTwo, "artifact.txt", "task two work\n")
	harness.advance(time.Minute)
	summary, err = ExecuteSupervisedReturn(ctx, harness.options, "complete")
	if err != nil {
		t.Fatalf("task two must complete its return, got %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("the two-task supervised chain must complete, got %s", summary.ChainState)
	}

	// Both tasks passed, each on its own step.
	if got := zeroModelTaskCount(summary, "PASSED"); got != 2 {
		t.Fatalf("both tasks must be PASSED, got %d: %+v", got, summary.TaskStateCounts)
	}

	// Per-task durable artifacts live under that task's own directory.
	for _, taskID := range []string{multiTaskAgentTaskOne, multiTaskAgentTaskTwo} {
		taskDir := filepath.Join(SupervisedDir(harness.runRoot), supervisedTasksDirectoryName, taskID+"-attempt-1")
		if _, err := os.Stat(taskDir); err != nil {
			t.Fatalf("task %q must own a supervised directory: %v", taskID, err)
		}
		for _, name := range []string{supervisedReceiptName, supervisedSessionName} {
			if _, err := os.Stat(filepath.Join(taskDir, name)); err != nil {
				t.Fatalf("task %q must persist %s in its own directory: %v", taskID, name, err)
			}
		}
	}

	// The single event log is totally ordered and every task freezes its own five
	// supervised facts.
	eventStore, err := chain.NewFileEventStore(runEventLogPath(harness.runRoot))
	if err != nil {
		t.Fatal(err)
	}
	events, err := eventStore.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	previous := 0
	for _, record := range events {
		if record.Event.Sequence <= previous {
			t.Fatalf("event sequence must increase, got %d after %d", record.Event.Sequence, previous)
		}
		previous = record.Event.Sequence
	}
	factsByTask := map[string][]string{}
	completedAt := map[string]int{}
	windowOpenedAt := map[string]int{}
	for _, record := range events {
		event := record.Event
		if event.Entity.TaskID == "" || event.Entity.Attempt != 1 {
			continue
		}
		switch event.Reason.Code {
		case "supervised-return-completed":
			completedAt[event.Entity.TaskID] = event.Sequence
			if len(event.InputEvidence) < 5 {
				t.Fatalf("task %q PASSED evidence carries only %d digests", event.Entity.TaskID, len(event.InputEvidence))
			}
			factsByTask[event.Entity.TaskID] = append([]string(nil), event.InputEvidence[:5]...)
			for index, digest := range factsByTask[event.Entity.TaskID] {
				if len(digest) == 0 {
					t.Fatalf("task %q fact %d is empty", event.Entity.TaskID, index)
				}
			}
		case "supervised-work-window-open":
			if _, seen := windowOpenedAt[event.Entity.TaskID]; !seen {
				windowOpenedAt[event.Entity.TaskID] = event.Sequence
			}
		}
	}
	if len(factsByTask[multiTaskAgentTaskOne]) != 5 || len(factsByTask[multiTaskAgentTaskTwo]) != 5 {
		t.Fatalf("both tasks must freeze five supervised facts, got %v", factsByTask)
	}
	if !(completedAt[multiTaskAgentTaskOne] != 0 && completedAt[multiTaskAgentTaskOne] < windowOpenedAt[multiTaskAgentTaskTwo] &&
		windowOpenedAt[multiTaskAgentTaskTwo] < completedAt[multiTaskAgentTaskTwo]) {
		t.Fatalf("the supervised tasks must interleave serially in one log: windows=%v completions=%v", windowOpenedAt, completedAt)
	}

	// The frozen five-tuple is positional, so "non-empty and distinct" alone would
	// stay green through a reordering regression. Each task's tuple must equal the
	// ordered result of re-deriving that task's facts from its OWN persisted
	// artifacts, through the same production re-derivation path the postflight
	// uses (newSupervisedRuntime + deriveSupervisedFacts + supervisedFactsList).
	clock, _ := supervisedClockAndIDs(harness.options)
	plan, err := supervisedAssemblyPlan(harness.options)
	if err != nil {
		t.Fatalf("the two-task plan must bind: %v", err)
	}
	replay, err := newSupervisedRuntime(harness.options, plan, clock)
	if err != nil {
		t.Fatalf("the re-derivation runtime must assemble: %v", err)
	}
	for _, taskID := range []string{multiTaskAgentTaskOne, multiTaskAgentTaskTwo} {
		derived, err := replay.deriveSupervisedFacts(ctx, taskID, 1)
		if err != nil {
			t.Fatalf("task %q facts must re-derive from its own artifacts: %v", taskID, err)
		}
		want := supervisedFactsList(derived)
		got := factsByTask[taskID]
		if len(got) != len(want) {
			t.Fatalf("task %q freezes %d facts, want %d (%v)", taskID, len(got), len(want), got)
		}
		for index := range want {
			if got[index] != want[index] {
				t.Fatalf("task %q fact %d = %q, want %q (rederived %v)", taskID, index, got[index], want[index], want)
			}
		}
	}
	for _, first := range factsByTask[multiTaskAgentTaskOne] {
		for _, second := range factsByTask[multiTaskAgentTaskTwo] {
			if first == second {
				t.Fatalf("the two tasks must not share a supervised fact: %s", first)
			}
		}
	}
}

// writeTaskArtifact writes into the isolated workspace of one supervised task.
func (harness *supervisedHarness) writeTaskArtifact(taskID, name, content string) {
	harness.t.Helper()
	path := filepath.Join(harness.runRoot, supervisedWorkspacesDirectory, taskID+"-attempt-1", name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		harness.t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		harness.t.Fatal(err)
	}
}
