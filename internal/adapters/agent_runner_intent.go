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

// AgentRunnerStepIntentRegistry projects several validated persisted requests
// into the Engine-owned execution path, one per (taskID, stepID).
//
// A multi-task assembly (CONTRACTS.md A6 多任务串行装配（P50）) binds one request
// record per task. Every record is admitted through the single-record
// constructor above, so the "matched exactly one definition entry" rule and the
// isolated-workspace step-shape rule are reused verbatim instead of re-derived
// here. On top of that the registry refuses cross-record drift: two records
// naming the same task/step, two records sharing one requestId (their R/C
// artifact streams would alias), and records that do not all belong to one run.
type AgentRunnerStepIntentRegistry struct {
	preparers map[string]chain.StepIntentPreparer
}

var _ chain.StepIntentPreparer = (*AgentRunnerStepIntentRegistry)(nil)

// NewAgentRunnerStepIntentRegistry builds the per-(taskID, stepID) preparer
// table. The order of records is irrelevant: the registry answers by binding,
// never by order, and the chain drives the tasks serially.
func NewAgentRunnerStepIntentRegistry(records []AgentRunnerRequestRecord, definition chain.Definition) (*AgentRunnerStepIntentRegistry, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("%w: at least one request record is required", ErrInvalidAgentRunnerRequest)
	}
	registry := &AgentRunnerStepIntentRegistry{preparers: make(map[string]chain.StepIntentPreparer, len(records))}
	runID := ""
	requestIDs := make(map[string]string, len(records))
	for _, record := range records {
		// Reuse the single-record admission: validates the record, requires
		// exactly one matching definition entry and the isolated-workspace code
		// step shape.
		preparer, err := NewAgentRunnerStepIntentPreparer(record, definition)
		if err != nil {
			return nil, err
		}
		key := AgentRunnerStepIntentKey(record.Request.TaskID, record.Request.StepID)
		if _, duplicate := registry.preparers[key]; duplicate {
			return nil, fmt.Errorf("%w: two request records bind task %q step %q", ErrInvalidAgentRunnerRequest, record.Request.TaskID, record.Request.StepID)
		}
		if previous, duplicate := requestIDs[record.Request.RequestID]; duplicate {
			return nil, fmt.Errorf("%w: request id %q is registered by both task %q and task %q", ErrInvalidAgentRunnerRequest, record.Request.RequestID, previous, record.Request.TaskID)
		}
		if runID == "" {
			runID = record.Request.RunID
		} else if record.Request.RunID != runID {
			return nil, fmt.Errorf("%w: request records drift across runs (%q and %q)", ErrInvalidAgentRunnerRequest, runID, record.Request.RunID)
		}
		requestIDs[record.Request.RequestID] = record.Request.TaskID
		registry.preparers[key] = preparer
	}
	return registry, nil
}

// AgentRunnerStepIntentKey is the stable registry key of one step binding: task
// and step ids are validated identifiers that never contain "/", so the encoding
// is unambiguous. It mirrors the console assembly's request-record key.
func AgentRunnerStepIntentKey(taskID, stepID string) string {
	return taskID + "/" + stepID
}

// Preparer returns the preparer bound to one task/step. An unbound step is
// refused instead of silently falling back to another task's request.
func (registry *AgentRunnerStepIntentRegistry) Preparer(taskID, stepID string) (chain.StepIntentPreparer, error) {
	if registry == nil {
		return nil, fmt.Errorf("%w: nil intent registry", ErrInvalidAgentRunnerRequest)
	}
	preparer, ok := registry.preparers[AgentRunnerStepIntentKey(taskID, stepID)]
	if !ok || preparer == nil {
		return nil, fmt.Errorf("%w: no step intent is bound to task %q step %q", ErrInvalidAgentRunnerRequest, taskID, stepID)
	}
	return preparer, nil
}

// Len reports how many task/step bindings the registry holds.
func (registry *AgentRunnerStepIntentRegistry) Len() int {
	if registry == nil {
		return 0
	}
	return len(registry.preparers)
}

// PrepareStepIntent selects the preparer by binding and delegates, so the
// registry can be handed to the engine directly as a StepIntentPreparer.
func (registry *AgentRunnerStepIntentRegistry) PrepareStepIntent(ctx context.Context, request chain.StepRequest) (chain.StepExecutionIntent, error) {
	preparer, err := registry.Preparer(request.TaskID, request.Step.ID)
	if err != nil {
		return chain.StepExecutionIntent{}, err
	}
	return preparer.PrepareStepIntent(ctx, request)
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
