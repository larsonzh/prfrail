package chain

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// SubmitAgentRunnerTerminal routes one external terminal fact into the core
// state. It is the only writer that may move an AgentRunner step out of
// TERMINAL_PENDING, it never dispatches or relaunches anything, and it never
// concludes a task: a completed terminal only passes the step, and the task
// still walks the existing freeze, gates, review and promotion path on the next
// Run.
//
// The route table is fixed: completed passes the step, failed fails the task
// and the chain, uncertain fails the task and keeps the chain paused without a
// retry, cancelled archives stop evidence before the terminal states are
// written, and operator-action-required parks the task and the step in
// WAITING_FOR_OPERATOR for the operator-interaction machine.
//
// A repeated terminal for the same request and completion converges on the
// recorded evidence and repairs any missing downstream transition; a repeated
// request with another completion, a terminal that does not match the parked
// dispatch, and an unproven resume continuity are refused before any write.
func (engine *Engine) SubmitAgentRunnerTerminal(ctx context.Context, terminal AgentRunnerTerminal) (AgentRunnerTerminalRoute, error) {
	if err := terminal.Validate(); err != nil {
		return AgentRunnerTerminalRoute{}, err
	}
	if terminal.RunID != engine.options.RunID || terminal.Attempt != 1 {
		return AgentRunnerTerminalRoute{}, fmt.Errorf("%w: terminal does not belong to this run and attempt", ErrInvalidAgentRunnerTerminal)
	}
	task, step, err := engine.lookupAgentRunnerStep(terminal)
	if err != nil {
		return AgentRunnerTerminalRoute{}, err
	}
	entity := stepEntity(engine.options.RunID, terminal.TaskID, terminal.StepID)
	events, err := engine.options.Events.Load(ctx)
	if err != nil {
		return AgentRunnerTerminalRoute{}, err
	}
	if terminal.resume() && !stepEvidenceHistoryContains(events, entity, terminal.PriorCompletionHash) {
		return AgentRunnerTerminalRoute{}, fmt.Errorf("%w: prior completion %s is not in step %s evidence history", ErrAgentRunnerResumeContinuityUnproven, terminal.PriorCompletionHash, step.ID)
	}
	routed, conflict, recorded := routedAgentTerminal(events, entity, terminal)
	if routed {
		if conflict {
			return AgentRunnerTerminalRoute{}, fmt.Errorf("%w: request %s already routed completion %s", ErrAgentRunnerTerminalConflict, terminal.RequestID, terminal.CompletionHash)
		}
		if terminal.Status == AgentRunnerTerminalCompleted {
			// The gate rebuilds the facts from the recorded evidence, so a replay
			// that claims the same request and completion but a different fact set
			// is a conflict: converging would hide the divergence from its caller.
			recordedFacts, factsErr := frozenFactsFromEvidence(recorded.InputEvidence)
			if factsErr != nil || terminal.Facts == nil || recordedFacts != *terminal.Facts {
				return AgentRunnerTerminalRoute{}, fmt.Errorf("%w: request %s already routed a different frozen fact set", ErrAgentRunnerTerminalConflict, terminal.RequestID)
			}
		}
		if err := engine.convergeAgentTerminalRoute(ctx, task, terminal, recorded.InputEvidence); err != nil {
			return AgentRunnerTerminalRoute{}, err
		}
		return engine.agentTerminalRouteReceipt(terminal, true), nil
	}
	if !latestEntityEventMatches(events, entity, "TERMINAL_PENDING", agentRunnerDispatchedReason, terminal.RequestHash) {
		return AgentRunnerTerminalRoute{}, fmt.Errorf("%w: step %s has no parked dispatch for request %s", ErrAgentRunnerTerminalConflict, step.ID, terminal.RequestID)
	}
	inputs := terminal.routeEvidence()
	if terminal.Status == AgentRunnerTerminalCancelled {
		// A cancelled terminal is only trustworthy once managed writers are
		// known to have stopped, so the stop evidence is archived in the very
		// transitions that record the terminal state.
		stopEvidence, err := engine.options.Stopper.Stop(ctx, engine.options.RunID)
		if err != nil {
			if engine.state.projection.ChainState == "RUNNING" {
				unproven := uniqueHashes(append([]string{terminal.CompletionHash}, stopEvidence...))
				_ = engine.transition(ctx, chainEntity(engine.options.RunID), "PAUSED", unproven, agentRunnerStopUncertainReason)
			}
			return AgentRunnerTerminalRoute{}, errors.Join(ErrRecoveryUncertain, fmt.Errorf("archive stop evidence before cancelling: %w", err))
		}
		inputs = uniqueHashes(append(stopEvidence, inputs...))
	}
	stepState, reason := agentTerminalStepRoute(terminal.Status)
	if err := engine.transition(ctx, entity, stepState, inputs, reason); err != nil {
		return AgentRunnerTerminalRoute{}, err
	}
	if err := engine.convergeAgentTerminalRoute(ctx, task, terminal, inputs); err != nil {
		return AgentRunnerTerminalRoute{}, err
	}
	return engine.agentTerminalRouteReceipt(terminal, false), nil
}

// agentTerminalStepRoute maps a terminal status to the single step transition
// that status is allowed to write.
func agentTerminalStepRoute(status AgentRunnerTerminalStatus) (string, string) {
	switch status {
	case AgentRunnerTerminalCompleted:
		return "PASSED", agentRunnerTerminalPassedReason
	case AgentRunnerTerminalOperatorActionRequired:
		return "WAITING_FOR_OPERATOR", agentRunnerTerminalWaitingReason
	case AgentRunnerTerminalCancelled:
		return "CANCELLED", agentRunnerTerminalCancelledReason
	case AgentRunnerTerminalFailed:
		return "FAILED", agentRunnerTerminalFailedReason
	default:
		return "FAILED", agentRunnerTerminalUncertainReason
	}
}

// convergeAgentTerminalRoute completes the task and chain transitions a routed
// terminal implies. Each transition is skipped when the entity already holds
// its routed state, so a replay after a partial write repairs exactly the
// missing transitions and nothing else.
func (engine *Engine) convergeAgentTerminalRoute(ctx context.Context, task Task, terminal AgentRunnerTerminal, inputs []string) error {
	taskEntityRef := taskEntity(engine.options.RunID, task.ID)
	chainEntityRef := chainEntity(engine.options.RunID)
	taskState, reason := "", ""
	switch terminal.Status {
	case AgentRunnerTerminalCompleted:
		// A completed terminal never touches the task: freeze, gates, review
		// and promotion stay the only path to a passed task.
		return nil
	case AgentRunnerTerminalOperatorActionRequired:
		// The waiting route is converged only while the step is still waiting:
		// once the operator machine answered it (or the step moved on for any
		// other proven reason), a repeated terminal must not push the task and
		// the run back into a waiting pause.
		if engine.state.current(stepEntity(engine.options.RunID, terminal.TaskID, terminal.StepID)) != "WAITING_FOR_OPERATOR" {
			return nil
		}
		taskState, reason = "WAITING_FOR_OPERATOR", agentRunnerTerminalWaitingReason
	case AgentRunnerTerminalCancelled:
		taskState, reason = "CANCELLED", agentRunnerTerminalCancelledReason
	case AgentRunnerTerminalFailed:
		taskState, reason = "FAILED", agentRunnerTerminalFailedReason
	default:
		taskState, reason = "FAILED", agentRunnerTerminalUncertainReason
	}
	if current := engine.state.current(taskEntityRef); current != taskState {
		if err := engine.transition(ctx, taskEntityRef, taskState, inputs, reason); err != nil {
			return err
		}
	}
	switch terminal.Status {
	case AgentRunnerTerminalFailed:
		if engine.state.projection.ChainState != "FAILED" {
			return engine.transition(ctx, chainEntityRef, "FAILED", inputs, reason)
		}
	case AgentRunnerTerminalUncertain:
		// Uncertain never fails or resumes the chain: the run stays paused and
		// a new attempt must be built instead of a retry.
		if engine.state.projection.ChainState == "RUNNING" {
			return engine.transition(ctx, chainEntityRef, "PAUSED", inputs, agentRunnerRecoveryUncertainReason)
		}
	case AgentRunnerTerminalOperatorActionRequired:
		// Hand control to the operator exactly like the classic interaction
		// path does: the run pauses while a human decision is pending, so a
		// crash between this route and the interaction record still converges
		// instead of leaving a parked run nobody can resume.
		if engine.state.projection.ChainState == "RUNNING" {
			return engine.transition(ctx, chainEntityRef, "PAUSED", inputs, agentRunnerTerminalWaitingReason)
		}
	case AgentRunnerTerminalCancelled:
		if engine.state.projection.ChainState != "CANCELLED" {
			return engine.transition(ctx, chainEntityRef, "CANCELLED", inputs, reason)
		}
	}
	return nil
}

func (engine *Engine) agentTerminalRouteReceipt(terminal AgentRunnerTerminal, converged bool) AgentRunnerTerminalRoute {
	return AgentRunnerTerminalRoute{
		Status:     terminal.Status,
		StepState:  engine.state.current(stepEntity(engine.options.RunID, terminal.TaskID, terminal.StepID)),
		TaskState:  engine.state.current(taskEntity(engine.options.RunID, terminal.TaskID)),
		ChainState: engine.state.projection.ChainState,
		Converged:  converged,
	}
}

// lookupAgentRunnerStep resolves the definition entry a terminal claims and
// verifies it is an AgentRunner-eligible step: only isolated-workspace code
// steps carry an external session.
func (engine *Engine) lookupAgentRunnerStep(terminal AgentRunnerTerminal) (Task, Step, error) {
	for _, task := range engine.options.Definition.Tasks {
		if task.ID != terminal.TaskID {
			continue
		}
		for _, step := range task.Steps {
			if step.ID != terminal.StepID {
				continue
			}
			if step.Kind != "code" || step.Mode != IsolatedWorkspace {
				return Task{}, Step{}, fmt.Errorf("%w: step %s is not an AgentRunner step", ErrInvalidAgentRunnerTerminal, step.ID)
			}
			return task, step, nil
		}
		return Task{}, Step{}, fmt.Errorf("%w: step %s is not in the run definition", ErrInvalidAgentRunnerTerminal, terminal.StepID)
	}
	return Task{}, Step{}, fmt.Errorf("%w: task %s is not in the run definition", ErrInvalidAgentRunnerTerminal, terminal.TaskID)
}

// routedAgentTerminal reports whether the step already routed a terminal for
// the submitted request, and returns the recorded routing event so a replay can
// converge on the evidence that was actually written. A conflicting completion
// for the same request is reported separately so callers fail closed instead of
// overwriting a decision.
func routedAgentTerminal(events []evidence.StateEvent, entity evidence.Entity, terminal AgentRunnerTerminal) (routed bool, conflict bool, recorded evidence.Event) {
	for index := len(events) - 1; index >= 0; index-- {
		event := events[index].Event
		if event.Entity != entity || !agentTerminalRoutingReason(event.Reason.Code) {
			continue
		}
		if !slices.Contains(event.InputEvidence, terminal.RequestHash) {
			continue
		}
		if !slices.Contains(event.InputEvidence, terminal.CompletionHash) {
			return true, true, event
		}
		return true, false, event
	}
	return false, false, evidence.Event{}
}
