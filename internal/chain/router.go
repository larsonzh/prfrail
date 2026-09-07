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

type HookPort interface {
	RunHook(context.Context, StepRequest) (StepResult, error)
}

type StepRouter struct {
	Managed  ManagedChangeSetPort
	Isolated IsolatedWorkspacePort
	Hooks    HookPort
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
