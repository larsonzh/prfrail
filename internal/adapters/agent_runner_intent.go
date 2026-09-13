package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/larsonzh/prfrail/internal/chain"
)

// AgentRunnerStepIntentPreparer projects one validated persisted request into
// the Engine-owned execution path. Admission remains responsible for policy,
// authorization, budget, capability, and replay checks.
type AgentRunnerStepIntentPreparer struct {
	requestRecord AgentRunnerRequestRecord
}

var _ chain.StepIntentPreparer = (*AgentRunnerStepIntentPreparer)(nil)

func NewAgentRunnerStepIntentPreparer(record AgentRunnerRequestRecord, definition chain.Definition) (*AgentRunnerStepIntentPreparer, error) {
	if err := ValidateAgentRunnerRequestRecord(record); err != nil {
		return nil, err
	}
	matches := 0
	for _, task := range definition.Tasks {
		for _, step := range task.Steps {
			if task.ID != record.Request.TaskID || step.ID != record.Request.StepID {
				continue
			}
			matches++
			if step.Kind != "code" || step.Mode != chain.IsolatedWorkspace {
				return nil, fmt.Errorf("%w: target step must be isolated workspace code", ErrInvalidAgentRunnerRequest)
			}
		}
	}
	if matches != 1 {
		return nil, fmt.Errorf("%w: target step matched %d definition entries", ErrInvalidAgentRunnerRequest, matches)
	}
	return &AgentRunnerStepIntentPreparer{requestRecord: record}, nil
}

func (preparer *AgentRunnerStepIntentPreparer) PrepareStepIntent(ctx context.Context, request chain.StepRequest) (chain.StepExecutionIntent, error) {
	if preparer == nil {
		return chain.StepExecutionIntent{}, fmt.Errorf("%w: nil intent preparer", ErrInvalidAgentRunnerRequest)
	}
	if err := ctx.Err(); err != nil {
		return chain.StepExecutionIntent{}, err
	}
	if err := ValidateAgentRunnerRequestRecord(preparer.requestRecord); err != nil {
		return chain.StepExecutionIntent{}, err
	}
	body := preparer.requestRecord.Request
	if request.RunID != body.RunID {
		return chain.StepExecutionIntent{}, fmt.Errorf("%w: run binding mismatch", ErrInvalidAgentRunnerRequest)
	}
	if request.TaskID != body.TaskID || request.Step.ID != body.StepID {
		return chain.StepExecutionIntent{}, nil
	}
	if request.Attempt != body.Attempt || request.Step.Kind != "code" || request.Step.Mode != chain.IsolatedWorkspace || request.ParentHash != body.ParentSnapshotHash || strings.TrimSpace(request.Workspace.Root) == "" {
		return chain.StepExecutionIntent{}, fmt.Errorf("%w: target step binding mismatch", ErrInvalidAgentRunnerRequest)
	}
	return chain.StepExecutionIntent{
		ExecutionTarget: chain.AgentRunnerExecution,
		AgentRunnerFacts: &chain.AgentRunnerImmutableFacts{
			RequestID:         body.RequestID,
			WorkspaceHash:     body.WorkspaceHash,
			ContextHash:       body.ContextHash,
			AuthorizationHash: body.AuthorizationHash,
			BudgetHash:        body.BudgetHash,
		},
	}, nil
}
