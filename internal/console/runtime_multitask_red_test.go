package console

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/chain"
)

// P50 CONSOLE-MULTITASK-ASSEMBLY -- accepted multi-task assembly shape.
//
// ③ pinned the single-task refusal wording while both assemblies still refused a
// definition with more than one task (03-tests/red-baseline.md §B). ② removed that
// gate, so this file is the green form of the same falsification: the two-task
// fixtures must now pass the step selection, and a definition with no task must
// still be refused with the dedicated zero-task wording. Each retired refusal
// text is kept as a negative assertion, so a regression that reintroduces the
// single-task gate turns this file red again.
//
// Nothing here writes to disk, spawns a process or reads a clock: every refusal
// asserted below happens before the run-root marker and before any I/O.

const (
	multiTaskDefinitionID = "chain-multitask-assembly"
	multiTaskAgentTaskOne = "task-one"
	multiTaskAgentStepOne = "code-one"
	multiTaskAgentTaskTwo = "task-two"
	multiTaskAgentStepTwo = "code-two"
)

// The retired single-task refusal wordings. They are asserted absent: a reworded
// or reintroduced gate is a behaviour change and must turn this file red.
const (
	agentRunnerRetiredMultiTaskRefusal = "the AgentRunner assembly requires exactly one task"
	agentRunnerSingleStepRefusal       = "the AgentRunner assembly requires exactly one step"
	agentRunnerStepShapeRefusal        = "the AgentRunner assembly requires one isolated-workspace code step"
	supervisedRetiredMultiTaskRefusal  = "the supervised assembly requires exactly one task"
	supervisedSingleStepRefusal        = "the supervised assembly requires exactly one step"
	supervisedStepShapeRefusal         = "the supervised assembly requires one supervised-black-box code step"

	// The zero-task refusal ② kept (and gave its own wording so it can never be
	// confused with multi-task acceptance).
	agentRunnerZeroTaskRefusal = "the AgentRunner assembly requires at least one task"
	supervisedZeroTaskRefusal  = "the supervised assembly requires at least one task"

	// Post-selection refusals: reached only when the definition passes the step
	// selection, so they prove the multi-task shape is accepted there.
	agentRunnerPostSelectionRefusal = "request record is missing or invalid"
	supervisedPostSelectionRefusal  = "run root is required"
)

// multiTaskIsolatedDefinition is the two-task AgentRunner fixture: two isolated
// workspaces, each task with a single code step. The chain-level validator must
// accept it (control leg), and the assembly must accept it too.
func multiTaskIsolatedDefinition() chain.Definition {
	return chain.Definition{ID: multiTaskDefinitionID, Tasks: []chain.Task{
		{
			ID:    multiTaskAgentTaskOne,
			Steps: []chain.Step{{ID: multiTaskAgentStepOne, Kind: "code", Mode: chain.IsolatedWorkspace}},
		},
		{
			ID:    multiTaskAgentTaskTwo,
			Steps: []chain.Step{{ID: multiTaskAgentStepTwo, Kind: "code", Mode: chain.IsolatedWorkspace}},
		},
	}}
}

// singleTaskIsolatedDefinition narrows the two-task fixture to one task, keeping
// the same step shape. It is the control leg for the AgentRunner assembly.
func singleTaskIsolatedDefinition(task chain.Task) chain.Definition {
	return chain.Definition{ID: multiTaskDefinitionID, Tasks: []chain.Task{task}}
}

// multiTaskSupervisedStep builds one supervised-black-box code step with a handoff
// policy that passes chain.validateHandoffPolicy.
func multiTaskSupervisedStep(id string) chain.Step {
	return chain.Step{
		ID:   id,
		Kind: "code",
		Mode: chain.SupervisedBlackBox,
		HandoffPolicy: &chain.HandoffPolicy{
			AllowedTargets:   []string{"target-one"},
			InputPolicy:      "structured",
			HandoffTimeoutMs: 60000,
			ReturnActions:    []string{"complete", "abort"},
			HooksAfterReturn: []string{"hook-one"},
		},
	}
}

// multiTaskSupervisedDefinition is the two-task supervised fixture.
func multiTaskSupervisedDefinition() chain.Definition {
	return chain.Definition{ID: multiTaskDefinitionID, Tasks: []chain.Task{
		{ID: multiTaskAgentTaskOne, Steps: []chain.Step{multiTaskSupervisedStep(multiTaskAgentStepOne)}},
		{ID: multiTaskAgentTaskTwo, Steps: []chain.Step{multiTaskSupervisedStep(multiTaskAgentStepTwo)}},
	}}
}

// requireDefinitionShapeAccepted fails when the fixture itself is malformed: a
// green must never be an artefact of an invalid definition.
func requireDefinitionShapeAccepted(t *testing.T, definition chain.Definition) {
	t.Helper()
	if err := chain.ValidateDefinition(definition); err != nil {
		t.Fatalf("the fixture must be a valid chain definition, got %v", err)
	}
}

// requireNoRetiredRefusal fails when an error still names any of the retired
// single-task gate wordings.
func requireNoRetiredRefusal(t *testing.T, err error, refusals ...string) {
	t.Helper()
	if err == nil {
		t.Fatal("a refused assembly must return an error")
	}
	for _, refusal := range refusals {
		if strings.Contains(err.Error(), refusal) {
			t.Fatalf("the single-task gate must be gone, got %q", err.Error())
		}
	}
}

// TestAgentRunnerAssemblyAcceptsMultiTaskDefinitionShape pins the accepted shape
// of the production AgentRunner assembly: a two-task definition passes the step
// selection and stops at the per-task request-record registry, never at a
// single-task gate. A definition with no task keeps its dedicated refusal.
func TestAgentRunnerAssemblyAcceptsMultiTaskDefinitionShape(t *testing.T) {
	definition := multiTaskIsolatedDefinition()
	requireDefinitionShapeAccepted(t, definition)
	if len(definition.Tasks) != 2 {
		t.Fatalf("the fixture must carry two tasks, got %d", len(definition.Tasks))
	}

	summary, err := ExecuteAgentRunnerRun(context.Background(), AgentRunnerAssemblyOptions{Definition: definition})
	if err == nil {
		t.Fatalf("a definition without request records must not reach a live run, got summary %+v", summary)
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("the refusal must wrap ErrInvalidConfig, got %v", err)
	}
	requireNoRetiredRefusal(t, err, agentRunnerRetiredMultiTaskRefusal, agentRunnerSingleStepRefusal, agentRunnerStepShapeRefusal)
	if !strings.Contains(err.Error(), agentRunnerPostSelectionRefusal) {
		t.Fatalf("the two-task definition must stop at the request-record registry (%q), got %q", agentRunnerPostSelectionRefusal, err.Error())
	}
	// The refusal happens before the run-root marker and before the engine: no
	// partial summary and no chain state may be reported.
	if summary.RunID != "" || summary.ChainState != "" || summary.Sequence != 0 || summary.EventLogPath != "" {
		t.Fatalf("a refused assembly must return an empty summary, got %+v", summary)
	}

	// The zero-task refusal is retained with its own wording, and it must run
	// before the step check (an empty definition has no task to inspect).
	_, zeroErr := ExecuteAgentRunnerRun(context.Background(), AgentRunnerAssemblyOptions{Definition: chain.Definition{ID: multiTaskDefinitionID}})
	if !errors.Is(zeroErr, ErrInvalidConfig) || !strings.Contains(zeroErr.Error(), agentRunnerZeroTaskRefusal) {
		t.Fatalf("a zero-task definition must be refused with %q, got %v", agentRunnerZeroTaskRefusal, zeroErr)
	}
	requireNoRetiredRefusal(t, zeroErr, agentRunnerSingleStepRefusal, agentRunnerStepShapeRefusal)

	// A task with two steps is still refused per task: dropping the single-task
	// gate must not drop the per-task step rule, and the rule must still run for
	// every task of the definition.
	overloaded := multiTaskIsolatedDefinition()
	overloaded.Tasks[1].Steps = append(overloaded.Tasks[1].Steps, chain.Step{ID: "code-extra", Kind: "code", Mode: chain.IsolatedWorkspace})
	_, stepErr := ExecuteAgentRunnerRun(context.Background(), AgentRunnerAssemblyOptions{Definition: overloaded})
	if !errors.Is(stepErr, ErrInvalidConfig) || !strings.Contains(stepErr.Error(), agentRunnerSingleStepRefusal) {
		t.Fatalf("a task with two steps must be refused per task, got %v", stepErr)
	}
	requireNoRetiredRefusal(t, stepErr, agentRunnerRetiredMultiTaskRefusal)
}

// TestSupervisedAssemblyAcceptsMultiTaskDefinitionShape is the supervised mirror:
// the two-task definition passes the step selection and stops at the run-root
// precondition, and the zero-task refusal keeps its dedicated wording.
func TestSupervisedAssemblyAcceptsMultiTaskDefinitionShape(t *testing.T) {
	definition := multiTaskSupervisedDefinition()
	requireDefinitionShapeAccepted(t, definition)
	if len(definition.Tasks) != 2 {
		t.Fatalf("the fixture must carry two tasks, got %d", len(definition.Tasks))
	}

	summary, err := ExecuteSupervisedRun(context.Background(), SupervisedAssemblyOptions{Definition: definition})
	if err == nil {
		t.Fatalf("a definition without a run root must not reach a live run, got summary %+v", summary)
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("the refusal must wrap ErrInvalidConfig, got %v", err)
	}
	requireNoRetiredRefusal(t, err, supervisedRetiredMultiTaskRefusal, supervisedSingleStepRefusal, supervisedStepShapeRefusal)
	if !strings.Contains(err.Error(), supervisedPostSelectionRefusal) {
		t.Fatalf("the two-task definition must stop at the run-root precondition (%q), got %q", supervisedPostSelectionRefusal, err.Error())
	}
	if summary.RunID != "" || summary.ChainState != "" || summary.Sequence != 0 || summary.EventLogPath != "" {
		t.Fatalf("a refused assembly must return an empty summary, got %+v", summary)
	}

	_, zeroErr := ExecuteSupervisedRun(context.Background(), SupervisedAssemblyOptions{Definition: chain.Definition{ID: multiTaskDefinitionID}})
	if !errors.Is(zeroErr, ErrInvalidConfig) || !strings.Contains(zeroErr.Error(), supervisedZeroTaskRefusal) {
		t.Fatalf("a zero-task definition must be refused with %q, got %v", supervisedZeroTaskRefusal, zeroErr)
	}
	requireNoRetiredRefusal(t, zeroErr, supervisedSingleStepRefusal, supervisedStepShapeRefusal)

	overloaded := multiTaskSupervisedDefinition()
	overloaded.Tasks[1].Steps = append(overloaded.Tasks[1].Steps, multiTaskSupervisedStep("code-extra"))
	_, stepErr := ExecuteSupervisedRun(context.Background(), SupervisedAssemblyOptions{Definition: overloaded})
	if !errors.Is(stepErr, ErrInvalidConfig) || !strings.Contains(stepErr.Error(), supervisedSingleStepRefusal) {
		t.Fatalf("a task with two steps must be refused per task, got %v", stepErr)
	}
	requireNoRetiredRefusal(t, stepErr, supervisedRetiredMultiTaskRefusal)
}

// TestMultiTaskDefinitionShapeIsOtherwiseValid is the control leg. It proves the
// accepted multi-task shape is not accepted by accident:
//
//  1. the two-task definitions pass chain.ValidateDefinition, the same semantic
//     validation the engine applies before a run starts;
//  2. each task of the two-task fixture, taken alone, passes the step selection
//     of both assemblies and is refused by the check that follows it;
//  3. the two-task errors carry none of the retired single-task/step/step-shape
//     wordings, i.e. neither assembly refuses the shape for another reason.
func TestMultiTaskDefinitionShapeIsOtherwiseValid(t *testing.T) {
	agentDefinition := multiTaskIsolatedDefinition()
	supervisedDefinition := multiTaskSupervisedDefinition()
	requireDefinitionShapeAccepted(t, agentDefinition)
	requireDefinitionShapeAccepted(t, supervisedDefinition)

	// 2. Every single-task slice of the two-task fixtures passes the step
	// selection, so each task's step shape is accepted on its own.
	for _, task := range agentDefinition.Tasks {
		single := singleTaskIsolatedDefinition(task)
		_, err := ExecuteAgentRunnerRun(context.Background(), AgentRunnerAssemblyOptions{Definition: single})
		if err == nil {
			t.Fatalf("task %q alone must not reach a live run", task.ID)
		}
		requireNoRetiredRefusal(t, err, agentRunnerRetiredMultiTaskRefusal, agentRunnerSingleStepRefusal, agentRunnerStepShapeRefusal)
		if !errors.Is(err, ErrInvalidConfig) || !strings.Contains(err.Error(), agentRunnerPostSelectionRefusal) {
			t.Fatalf("task %q alone must stop at the request-record check (%q), got %v", task.ID, agentRunnerPostSelectionRefusal, err)
		}
	}
	for _, task := range supervisedDefinition.Tasks {
		single := chain.Definition{ID: multiTaskDefinitionID, Tasks: []chain.Task{task}}
		_, err := ExecuteSupervisedRun(context.Background(), SupervisedAssemblyOptions{Definition: single})
		if err == nil {
			t.Fatalf("task %q alone must not reach a live run", task.ID)
		}
		requireNoRetiredRefusal(t, err, supervisedRetiredMultiTaskRefusal, supervisedSingleStepRefusal, supervisedStepShapeRefusal)
		if !errors.Is(err, ErrInvalidConfig) || !strings.Contains(err.Error(), supervisedPostSelectionRefusal) {
			t.Fatalf("task %q alone must stop at the run-root check (%q), got %v", task.ID, supervisedPostSelectionRefusal, err)
		}
	}

	// 3. The two-task errors name no retired gate.
	_, agentErr := ExecuteAgentRunnerRun(context.Background(), AgentRunnerAssemblyOptions{Definition: agentDefinition})
	requireNoRetiredRefusal(t, agentErr, agentRunnerRetiredMultiTaskRefusal, agentRunnerSingleStepRefusal, agentRunnerStepShapeRefusal)
	_, supervisedErr := ExecuteSupervisedRun(context.Background(), SupervisedAssemblyOptions{Definition: supervisedDefinition})
	requireNoRetiredRefusal(t, supervisedErr, supervisedRetiredMultiTaskRefusal, supervisedSingleStepRefusal, supervisedStepShapeRefusal)
}
