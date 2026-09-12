package chain

import (
	"context"
	"fmt"
)

type ManagedChangeSetPort interface {
	ApplyManagedChangeSet(context.Context, StepRequest) (StepResult, error)
}

type IsolatedWorkspacePort interface {
	CaptureIsolatedWorkspace(context.Context, StepRequest) (StepResult, error)
}

type AgentRunnerPort interface {
	DispatchAgentRunner(context.Context, StepRequest) (StepResult, error)
}

// AgentRunnerAdmission owns the runtime preflight hook. Implementations must
// validate all persisted availability, capability, enforcement, platform,
// authorization, and budget bindings required by the AgentRunner contract.
type AgentRunnerAdmission interface {
	AdmitAgentRunner(context.Context, StepRequest) error
}

type ManualHandoffPort interface {
	OpenManualHandoff(context.Context, StepRequest) (StepResult, error)
}

type HookPort interface {
	RunHook(context.Context, StepRequest) (StepResult, error)
}

type StepRouter struct {
	Managed              ManagedChangeSetPort
	Isolated             IsolatedWorkspacePort
	AgentRunner          AgentRunnerPort
	AgentRunnerAdmission AgentRunnerAdmission
	Handoff              ManualHandoffPort
	Hooks                HookPort
}

func (router StepRouter) ExecuteAgentRunner(ctx context.Context, request StepRequest) (StepResult, error) {
	// This is an offline integration hook until Engine routing is explicitly
	// wired. Production AgentRunner dispatch must use this admission-first path.
	if request.Step.Kind != "code" || request.Step.Mode != IsolatedWorkspace {
		return StepResult{}, fmt.Errorf("%w: AgentRunner requires isolated workspace code step", ErrInvalidDefinition)
	}
	if router.AgentRunner == nil {
		return StepResult{}, fmt.Errorf("%w: AgentRunner port unavailable", ErrInvalidDefinition)
	}
	if router.AgentRunnerAdmission == nil {
		return StepResult{}, fmt.Errorf("%w: AgentRunner admission unavailable", ErrInvalidDefinition)
	}
	if err := router.AgentRunnerAdmission.AdmitAgentRunner(ctx, request); err != nil {
		return StepResult{}, err
	}
	return router.AgentRunner.DispatchAgentRunner(ctx, request)
}

func (router StepRouter) Execute(ctx context.Context, request StepRequest) (StepResult, error) {
	switch request.Step.Kind {
	case "code":
		switch request.Step.Mode {
		case ManagedChangeSet:
			if router.Managed == nil {
				return StepResult{}, fmt.Errorf("%w: managed change-set port unavailable", ErrInvalidDefinition)
			}
			return router.Managed.ApplyManagedChangeSet(ctx, request)
		case IsolatedWorkspace:
			if router.Isolated == nil {
				return StepResult{}, fmt.Errorf("%w: isolated workspace port unavailable", ErrInvalidDefinition)
			}
			return router.Isolated.CaptureIsolatedWorkspace(ctx, request)
		case ManualHandoff:
			if router.Handoff == nil {
				return StepResult{}, fmt.Errorf("%w: manual handoff port unavailable", ErrInvalidDefinition)
			}
			return router.Handoff.OpenManualHandoff(ctx, request)
		default:
			return StepResult{}, fmt.Errorf("%w: unknown change mode", ErrInvalidDefinition)
		}
	case "build", "verify":
		if router.Hooks == nil {
			return StepResult{}, fmt.Errorf("%w: hook port unavailable", ErrInvalidDefinition)
		}
		return router.Hooks.RunHook(ctx, request)
	case "noop":
		return StepResult{}, fmt.Errorf("%w: noop must not be executed", ErrInvalidState)
	default:
		return StepResult{}, fmt.Errorf("%w: unknown step kind", ErrInvalidDefinition)
	}
}
