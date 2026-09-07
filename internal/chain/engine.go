package chain

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type Failpoint func(string) error

type Options struct {
	RunID      string
	Definition Definition
	Events     EventStore
	Baselines  BaselinePort
	Workspaces WorkspacePort
	Steps      StepPort
	Acceptance AcceptancePort
	Reviewer   ReviewerPort
	Publisher  PublisherPort
	Stopper    Stopper
	Reconciler Reconciler
	Clock      Clock
	IDs        IDSource
	Failpoint  Failpoint
}

type Engine struct {
	options Options
	state   stateWriter
	pause   atomic.Bool
	cancel  atomic.Bool
}

func New(ctx context.Context, options Options) (*Engine, error) {
	if !evidence.ValidID(options.RunID) || options.Events == nil || options.Baselines == nil || options.Workspaces == nil || options.Steps == nil || options.Acceptance == nil || options.Reviewer == nil || options.Publisher == nil || options.Stopper == nil || options.Reconciler == nil {
		return nil, ErrInvalidDefinition
	}
	if err := options.Definition.validate(); err != nil {
		return nil, err
	}
	if options.Clock == nil {
		options.Clock = defaultClock
	}
	if options.IDs == nil {
		options.IDs = randomID
	}
	events, err := options.Events.Load(ctx)
	if err != nil {
		return nil, err
	}
	projection, err := Rebuild(events)
	if err != nil {
		return nil, err
	}
	if projection.RunID != "" && projection.RunID != options.RunID {
		return nil, fmt.Errorf("%w: event log belongs to %q", ErrInvalidState, projection.RunID)
	}
	engine := &Engine{options: options, state: stateWriter{store: options.Events, projection: projection, clock: options.Clock, ids: options.IDs}}
	if projection.ChainState == "" {
		if err := engine.transition(ctx, chainEntity(options.RunID), "CREATED", nil, "run-created"); err != nil {
			return nil, err
		}
	}
	return engine, nil
}

func (engine *Engine) Projection() Projection { return cloneProjection(engine.state.projection) }
func (engine *Engine) RequestPause()          { engine.pause.Store(true) }
func (engine *Engine) RequestCancel()         { engine.cancel.Store(true) }

func (engine *Engine) Initialize(ctx context.Context) error {
	if engine.state.projection.ChainState != "CREATED" {
		return fmt.Errorf("%w: initialize from %s", ErrInvalidState, engine.state.projection.ChainState)
	}
	baseline, err := engine.options.Baselines.CaptureBaseline(ctx, engine.options.RunID)
	if err != nil {
		return engine.failChain(ctx, "baseline-failed", err)
	}
	if !evidence.ValidHash(baseline.Hash) {
		return engine.failChain(ctx, "baseline-invalid", ErrInvalidState)
	}
	if err := engine.transition(ctx, chainEntity(engine.options.RunID), "BASELINED", []string{baseline.Hash}, "baseline-captured"); err != nil {
		return err
	}
	engine.state.projection.AcceptedParent = baseline.Hash
	return nil
}

func (engine *Engine) Run(ctx context.Context) error {
	switch engine.state.projection.ChainState {
	case "CREATED":
		if err := engine.Initialize(ctx); err != nil {
			return err
		}
	case "PAUSED":
		if engine.state.projection.RecoveryUncertain {
			return ErrRecoveryUncertain
		}
		if err := engine.transition(ctx, chainEntity(engine.options.RunID), "RUNNING", nil, "operator-resumed"); err != nil {
			return err
		}
	case "BASELINED":
	case "RUNNING":
	case "COMPLETED":
		return nil
	case "FAILED", "CANCELLED":
		return fmt.Errorf("%w: terminal chain state %s", ErrInvalidState, engine.state.projection.ChainState)
	default:
		return fmt.Errorf("%w: unknown chain state", ErrInvalidState)
	}
	if engine.state.projection.ChainState == "BASELINED" {
		if err := engine.transition(ctx, chainEntity(engine.options.RunID), "RUNNING", nil, "chain-started"); err != nil {
			return err
		}
	}
	for _, task := range engine.options.Definition.Tasks {
		if engine.state.projection.TaskStates[taskKey(task.ID, 1)] == "PASSED" {
			continue
		}
		if err := engine.boundary(ctx); err != nil {
			return err
		}
		if err := engine.runTask(ctx, task); err != nil {
			return err
		}
	}
	return engine.transition(ctx, chainEntity(engine.options.RunID), "COMPLETED", nil, "chain-completed")
}

func (engine *Engine) runTask(ctx context.Context, task Task) error {
	entity := taskEntity(engine.options.RunID, task.ID)
	state := engine.state.current(entity)
	if state == "NONE" {
		if err := engine.transition(ctx, entity, "PENDING", nil, "task-created"); err != nil {
			return err
		}
		state = "PENDING"
	}
	if state == "PENDING" {
		if err := engine.transition(ctx, entity, "PRECHECK", nil, "task-precheck"); err != nil {
			return err
		}
		state = "PRECHECK"
	}
	if state == "REPAIR_PENDING" || state == "WAITING_FOR_OPERATOR" {
		if engine.state.projection.ChainState == "RUNNING" {
			if err := engine.transition(ctx, chainEntity(engine.options.RunID), "PAUSED", nil, "repair-required"); err != nil {
				return err
			}
		}
		return ErrPaused
	}
	parent := SnapshotRef{Hash: engine.state.projection.AcceptedParent}
	if !evidence.ValidHash(parent.Hash) {
		return engine.failTask(ctx, task.ID, "accepted-parent-missing", ErrInvalidState)
	}
	workspace, err := engine.options.Workspaces.Materialize(ctx, engine.options.RunID, task.ID, 1, parent)
	if err != nil {
		return engine.failTask(ctx, task.ID, "workspace-materialize-failed", err)
	}
	if state == "PRECHECK" {
		if err := engine.transition(ctx, entity, "STEPS_RUNNING", nil, "steps-started"); err != nil {
			return err
		}
	}
	for _, step := range task.Steps {
		if err := engine.boundary(ctx); err != nil {
			return err
		}
		if err := engine.runStep(ctx, task.ID, step, parent, workspace); err != nil {
			return err
		}
	}
	if err := engine.transition(ctx, entity, "REVIEW_PENDING", nil, "steps-completed"); err != nil {
		return err
	}
	candidate, err := engine.options.Acceptance.Accept(ctx, engine.options.RunID, task.ID, 1, parent, workspace)
	if err != nil {
		return engine.failTask(ctx, task.ID, "candidate-freeze-failed", err)
	}
	if !evidence.ValidHash(candidate.Snapshot.Hash) || !evidence.ValidHash(candidate.EvidenceRootHash) {
		return engine.failTask(ctx, task.ID, "candidate-invalid", ErrInvalidState)
	}
	review, err := engine.options.Reviewer.Review(ctx, ReviewRequest{
		RunID:                 engine.options.RunID,
		TaskID:                task.ID,
		Attempt:               1,
		ParentSnapshotHash:    parent.Hash,
		CandidateSnapshotHash: candidate.Snapshot.Hash,
		EvidenceRootHash:      candidate.EvidenceRootHash,
		CandidateProducer:     candidate.CandidateProducer,
	})
	if err != nil {
		return engine.failTask(ctx, task.ID, "review-failed", err)
	}
	reviewResult, err := buildReviewResult(ReviewRequest{
		RunID:                 engine.options.RunID,
		TaskID:                task.ID,
		Attempt:               1,
		ParentSnapshotHash:    parent.Hash,
		CandidateSnapshotHash: candidate.Snapshot.Hash,
		EvidenceRootHash:      candidate.EvidenceRootHash,
		CandidateProducer:     candidate.CandidateProducer,
	}, review, engine.options.Clock, engine.options.IDs)
	if err != nil {
		return engine.markTaskForRepair(ctx, task.ID, "review-invalid", append(candidate.Evidence, candidate.EvidenceRootHash))
	}
	if review.Outcome == "reject" {
		inputs := append([]string{reviewResult.ReceiptHash, candidate.Snapshot.Hash, candidate.EvidenceRootHash}, review.ErrorEvidence...)
		inputs = uniqueHashes(inputs)
		return engine.markTaskForRepair(ctx, task.ID, "review-rejected", inputs)
	}
	promotion, err := engine.options.Publisher.Publish(ctx, PromotionRequest{
		RunID:                 engine.options.RunID,
		TaskID:                task.ID,
		Attempt:               1,
		ParentSnapshotHash:    parent.Hash,
		CandidateSnapshotHash: candidate.Snapshot.Hash,
		EvidenceRootHash:      candidate.EvidenceRootHash,
		ReviewReceiptHash:     reviewResult.ReceiptHash,
	})
	if err != nil {
		return engine.failTask(ctx, task.ID, "promotion-failed", err)
	}
	promotionResult, err := buildPromotionResult(PromotionRequest{
		RunID:                 engine.options.RunID,
		TaskID:                task.ID,
		Attempt:               1,
		ParentSnapshotHash:    parent.Hash,
		CandidateSnapshotHash: candidate.Snapshot.Hash,
		EvidenceRootHash:      candidate.EvidenceRootHash,
		ReviewReceiptHash:     reviewResult.ReceiptHash,
	}, promotion, engine.options.Clock, engine.options.IDs)
	if err != nil {
		return engine.markTaskForRepair(ctx, task.ID, "promotion-invalid", append(candidate.Evidence, candidate.EvidenceRootHash))
	}
	if promotion.Outcome != "completed" {
		inputs := append([]string{reviewResult.ReceiptHash, promotionResult.ReceiptHash, candidate.Snapshot.Hash}, promotion.ErrorEvidence...)
		inputs = uniqueHashes(inputs)
		return engine.markTaskForRepair(ctx, task.ID, "promotion-incomplete", inputs)
	}
	inputs := []string{candidate.Snapshot.Hash, reviewResult.ReceiptHash, promotionResult.ReceiptHash, candidate.EvidenceRootHash}
	inputs = append(inputs, candidate.Evidence...)
	inputs = append(inputs, review.Evidence...)
	inputs = append(inputs, promotion.Evidence...)
	inputs = uniqueHashes(inputs)
	if err := engine.transition(ctx, entity, "PASSED", inputs, "task-accepted"); err != nil {
		return err
	}
	engine.state.projection.AcceptedParent = candidate.Snapshot.Hash
	return nil
}

func (engine *Engine) markTaskForRepair(ctx context.Context, taskID, reason string, input []string) error {
	entity := taskEntity(engine.options.RunID, taskID)
	input = uniqueHashes(input)
	if err := engine.transition(ctx, entity, "REPAIR_PENDING", input, reason); err != nil {
		return err
	}
	if engine.state.projection.ChainState == "RUNNING" {
		if err := engine.transition(ctx, chainEntity(engine.options.RunID), "PAUSED", input, reason); err != nil {
			return err
		}
	}
	return ErrPaused
}

func (engine *Engine) runStep(ctx context.Context, taskID string, step Step, parent SnapshotRef, workspace Workspace) error {
	entity := stepEntity(engine.options.RunID, taskID, step.ID)
	state := engine.state.current(entity)
	if state == "PASSED" || state == "NOOP_RECORDED" {
		return nil
	}
	if state == "NONE" {
		if err := engine.transition(ctx, entity, "PENDING", nil, "step-created"); err != nil {
			return err
		}
		state = "PENDING"
	}
	if step.Kind == "noop" {
		return engine.transition(ctx, entity, "NOOP_RECORDED", nil, "noop-recorded")
	}
	if state != "PENDING" {
		return fmt.Errorf("%w: step %q is %s", ErrRecoveryUncertain, step.ID, state)
	}
	if err := engine.transition(ctx, entity, "RUNNING", nil, "step-started"); err != nil {
		return err
	}
	if engine.options.Failpoint != nil {
		if err := engine.options.Failpoint("step-running"); err != nil {
			return err
		}
	}
	result, err := engine.options.Steps.Execute(ctx, StepRequest{RunID: engine.options.RunID, TaskID: taskID, Step: step, Attempt: 1, ParentHash: parent.Hash, Workspace: workspace})
	if err != nil {
		if transitionErr := engine.transition(ctx, entity, "FAILED", result.Evidence, "step-failed"); transitionErr != nil {
			return errors.Join(err, transitionErr)
		}
		return engine.failTask(ctx, taskID, "step-failed", err)
	}
	return engine.transition(ctx, entity, "PASSED", result.Evidence, "step-passed")
}

func (engine *Engine) boundary(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if engine.cancel.Load() {
		return engine.cancelRun(ctx)
	}
	if engine.pause.Load() {
		engine.pause.Store(false)
		if engine.state.projection.ChainState == "RUNNING" {
			if err := engine.transition(ctx, chainEntity(engine.options.RunID), "PAUSED", nil, "operator-paused"); err != nil {
				return err
			}
		}
		return ErrPaused
	}
	return nil
}

func (engine *Engine) cancelRun(ctx context.Context) error {
	evidenceHashes, err := engine.options.Stopper.Stop(ctx, engine.options.RunID)
	if err != nil {
		if engine.state.projection.ChainState == "RUNNING" {
			_ = engine.transition(ctx, chainEntity(engine.options.RunID), "PAUSED", evidenceHashes, "stop-uncertain")
		}
		return errors.Join(ErrRecoveryUncertain, err)
	}
	for _, task := range engine.options.Definition.Tasks {
		for _, step := range task.Steps {
			entity := stepEntity(engine.options.RunID, task.ID, step.ID)
			state := engine.state.current(entity)
			if state == "PENDING" || state == "RUNNING" || state == "WAITING_FOR_OPERATOR" {
				if err := engine.transition(ctx, entity, "CANCELLED", evidenceHashes, "operator-cancelled"); err != nil {
					return err
				}
			}
		}
		entity := taskEntity(engine.options.RunID, task.ID)
		state := engine.state.current(entity)
		if state == "PENDING" || state == "PRECHECK" || state == "STEPS_RUNNING" || state == "WAITING_FOR_OPERATOR" || state == "REVIEW_PENDING" || state == "REPAIR_PENDING" {
			if err := engine.transition(ctx, entity, "CANCELLED", evidenceHashes, "operator-cancelled"); err != nil {
				return err
			}
		}
	}
	if err := engine.transition(ctx, chainEntity(engine.options.RunID), "CANCELLED", evidenceHashes, "operator-cancelled"); err != nil {
		return err
	}
	return ErrCancelled
}

func (engine *Engine) failTask(ctx context.Context, taskID, reason string, cause error) error {
	entity := taskEntity(engine.options.RunID, taskID)
	state := engine.state.current(entity)
	if state != "FAILED" {
		if err := engine.transition(ctx, entity, "FAILED", nil, reason); err != nil {
			return errors.Join(cause, err)
		}
	}
	return engine.failChain(ctx, reason, cause)
}

func (engine *Engine) failChain(ctx context.Context, reason string, cause error) error {
	if engine.state.projection.ChainState != "FAILED" {
		if err := engine.transition(ctx, chainEntity(engine.options.RunID), "FAILED", nil, reason); err != nil {
			return errors.Join(cause, err)
		}
	}
	return cause
}

func (engine *Engine) transition(ctx context.Context, entity evidence.Entity, to string, input []string, reason string) error {
	return engine.state.append(ctx, evidence.Actor{Type: "system", ID: "engine"}, entity, to, input, reason)
}

func chainEntity(runID string) evidence.Entity { return evidence.Entity{Kind: "chain", RunID: runID} }
func taskEntity(runID, taskID string) evidence.Entity {
	return evidence.Entity{Kind: "task", RunID: runID, TaskID: taskID, Attempt: 1}
}
func stepEntity(runID, taskID, stepID string) evidence.Entity {
	return evidence.Entity{Kind: "step", RunID: runID, TaskID: taskID, StepID: stepID, Attempt: 1}
}

func cloneProjection(source Projection) Projection {
	result := source
	result.TaskStates = make(map[string]string, len(source.TaskStates))
	result.StepStates = make(map[string]string, len(source.StepStates))
	for key, value := range source.TaskStates {
		result.TaskStates[key] = value
	}
	for key, value := range source.StepStates {
		result.StepStates[key] = value
	}
	return result
}

func uniqueHashes(values []string) []string {
	if len(values) < 2 {
		return values
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

var eventCounter atomic.Uint64

func randomID() string {
	return fmt.Sprintf("event-%d-%d", defaultClock().UTC().UnixNano(), eventCounter.Add(1))
}
