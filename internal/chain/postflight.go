package chain

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	// ErrInvalidPostflightFacts reports frozen facts that cannot be trusted:
	// a missing or malformed digest, or digests that are not mutually distinct.
	ErrInvalidPostflightFacts = errors.New("invalid postflight frozen facts")
	// ErrInvalidPostflightDecision reports a postflight decision that cannot be
	// accepted: an unknown outcome, or a passed decision whose evidence does not
	// prove that all six frozen facts were consumed.
	ErrInvalidPostflightDecision = errors.New("invalid postflight decision")
	// ErrPostflightUnavailable reports a missing postflight port. Postflight is
	// the only qualification for REVIEW_PENDING, so a run without it must fail
	// closed instead of silently skipping the gate.
	ErrPostflightUnavailable = errors.New("postflight port unavailable")
	// ErrUnqualifiedReview reports a task resting in REVIEW_PENDING whose review
	// transition does not prove that the postflight gate produced it. Such a
	// state cannot be resumed in place and is refused without writing.
	ErrUnqualifiedReview = errors.New("review state not qualified by postflight")
)

// postflightFactsPrefix is the fixed position of the six frozen-fact digests
// inside the terminal-passed routing evidence: request, completion, then
// manifest, diff, log, usage, process-stop and pid-reports. The layout is
// deliberately positional and documented here because the gate rebuilds the
// facts from the recorded event evidence without reading the replay store; the
// mutual distinctness enforced by Validate is what keeps the positions
// unambiguous.
const (
	postflightFactsPrefix      = 2
	postflightFactsCount       = 6
	postflightFactsMinEvidence = postflightFactsPrefix + postflightFactsCount
)

// AgentRunnerFrozenFacts are the six chain-owned digests of the facts an
// adapter froze after external execution ended. They carry no decision: an exit
// code, a process status, a settlement status, a review decision or a policy
// hash never belong here, because the chain decides from its own postflight.
type AgentRunnerFrozenFacts struct {
	ManifestHash            string
	DiffHash                string
	LogHash                 string
	UsageHash               string
	ProcessStopEvidenceHash string
	ProcessPidReportsHash   string
}

// hashes returns the digests in the fixed routing order.
func (facts AgentRunnerFrozenFacts) hashes() []string {
	return []string{facts.ManifestHash, facts.DiffHash, facts.LogHash, facts.UsageHash, facts.ProcessStopEvidenceHash, facts.ProcessPidReportsHash}
}

// Validate fails closed on missing digests and on digest collisions: the
// routing layout is positional, so repeated digests would make the rebuilt
// facts ambiguous.
func (facts AgentRunnerFrozenFacts) Validate() error {
	hashes := facts.hashes()
	for index, hash := range hashes {
		if !evidence.ValidHash(hash) {
			return fmt.Errorf("%w: digest %d is not a valid hash", ErrInvalidPostflightFacts, index)
		}
	}
	for index, hash := range hashes {
		if slices.Contains(hashes[:index], hash) {
			return fmt.Errorf("%w: digest %d repeats an earlier digest", ErrInvalidPostflightFacts, index)
		}
	}
	return nil
}

// DisjointFrom reports whether the facts collide with digests they must differ
// from (the request and completion digests of the same terminal).
func (facts AgentRunnerFrozenFacts) DisjointFrom(hashes ...string) error {
	for _, fact := range facts.hashes() {
		if slices.Contains(hashes, fact) {
			return fmt.Errorf("%w: fact digest collides with the terminal binding", ErrInvalidPostflightFacts)
		}
	}
	return nil
}

// PostflightOutcome is the raw verdict of the chain-owned postflight.
type PostflightOutcome string

const (
	// PostflightPassed means the frozen facts were re-derived and reconciled.
	PostflightPassed PostflightOutcome = "passed"
	// PostflightFailed means postflight ran with complete evidence and rejected
	// the frozen facts; the task goes to REPAIR_PENDING and the chain pauses.
	PostflightFailed PostflightOutcome = "failed"
	// PostflightUncertain means postflight could not prove the facts; the task
	// fails and the chain pauses without a retry.
	PostflightUncertain PostflightOutcome = "uncertain"
)

// PostflightRequest carries everything the port may rely on. The port must
// re-derive the facts itself (fresh stop evidence, a re-captured manifest
// diffed against the parent, scope/secret/type and side-effect checks) and
// reconcile them field by field; it must not write state and must not reach the
// Engine.
type PostflightRequest struct {
	RunID              string
	TaskID             string
	Attempt            int
	ParentSnapshotHash string
	Workspace          Workspace
	Facts              AgentRunnerFrozenFacts
}

// PostflightDecision is the port's verdict plus its evidence. A passed verdict
// must carry all six fact digests, so the engine can prove that the gate
// consumed the frozen facts instead of trusting the port.
type PostflightDecision struct {
	Outcome       PostflightOutcome
	Evidence      []string
	ErrorEvidence []string
}

// Validate checks the outcome and the fact binding of a passed verdict.
func (decision PostflightDecision) Validate(facts AgentRunnerFrozenFacts) error {
	switch decision.Outcome {
	case PostflightPassed, PostflightFailed, PostflightUncertain:
	default:
		return fmt.Errorf("%w: unknown outcome %q", ErrInvalidPostflightDecision, decision.Outcome)
	}
	for _, hash := range decision.Evidence {
		if !evidence.ValidHash(hash) {
			return fmt.Errorf("%w: invalid evidence digest", ErrInvalidPostflightDecision)
		}
	}
	for _, hash := range decision.ErrorEvidence {
		if !evidence.ValidHash(hash) {
			return fmt.Errorf("%w: invalid error evidence digest", ErrInvalidPostflightDecision)
		}
	}
	if decision.Outcome != PostflightPassed {
		return nil
	}
	for _, fact := range facts.hashes() {
		if !slices.Contains(decision.Evidence, fact) {
			return fmt.Errorf("%w: a passed postflight must carry every frozen fact digest", ErrInvalidPostflightDecision)
		}
	}
	return nil
}

// PostflightPort runs the chain-owned postflight for one task. It is a required
// port: implementation quality is the only trust boundary, and the port never
// writes chain state.
type PostflightPort interface {
	RunPostflight(context.Context, PostflightRequest) (PostflightDecision, error)
}

// SupervisedPostflightPort runs the chain-owned postflight of a supervised
// black-box task. The port re-derives the five supervised facts itself (a fresh
// workspace scan, a managed-process stop proof, the frozen gates, the reduced
// report and the returned handoff receipt) and reconciles them field by field;
// it must not write state and must not reach the Engine.
type SupervisedPostflightPort interface {
	RunSupervisedPostflight(context.Context, SupervisedPostflightRequest) (SupervisedPostflightDecision, error)
}

// SupervisedPostflightRequest carries everything the supervised postflight port
// may rely on.
type SupervisedPostflightRequest struct {
	RunID              string
	TaskID             string
	Attempt            int
	ParentSnapshotHash string
	Workspace          Workspace
	Facts              SupervisedFrozenFacts
}

// SupervisedPostflightDecision is the supervised port's verdict plus its
// evidence. A passed verdict must carry all five fact digests, so the engine can
// prove that the gate consumed the frozen facts instead of trusting the port.
type SupervisedPostflightDecision struct {
	Outcome       PostflightOutcome
	Evidence      []string
	ErrorEvidence []string
}

// Validate checks the outcome and the fact binding of a passed verdict.
func (decision SupervisedPostflightDecision) Validate(facts SupervisedFrozenFacts) error {
	switch decision.Outcome {
	case PostflightPassed, PostflightFailed, PostflightUncertain:
	default:
		return fmt.Errorf("%w: unknown outcome %q", ErrInvalidPostflightDecision, decision.Outcome)
	}
	for _, hash := range decision.Evidence {
		if !evidence.ValidHash(hash) {
			return fmt.Errorf("%w: invalid evidence digest", ErrInvalidPostflightDecision)
		}
	}
	for _, hash := range decision.ErrorEvidence {
		if !evidence.ValidHash(hash) {
			return fmt.Errorf("%w: invalid error evidence digest", ErrInvalidPostflightDecision)
		}
	}
	if decision.Outcome != PostflightPassed {
		return nil
	}
	for _, fact := range facts.hashes() {
		if !slices.Contains(decision.Evidence, fact) {
			return fmt.Errorf("%w: a passed supervised postflight must carry every frozen fact digest", ErrInvalidPostflightDecision)
		}
	}
	return nil
}

// findTask resolves a task of the run definition.
func (engine *Engine) findTask(taskID string) (Task, bool) {
	for _, task := range engine.options.Definition.Tasks {
		if task.ID == taskID {
			return task, true
		}
	}
	return Task{}, false
}

// taskHasSupervisedStep reports whether a task contains a supervised-black-box
// step. Such a task may only reach review from the supervised postflight branch.
func taskHasSupervisedStep(task Task) bool {
	for _, step := range task.Steps {
		if step.Mode == SupervisedBlackBox {
			return true
		}
	}
	return false
}

// agentRunnerFactsEvidence returns the fact digests in routing order, for the
// terminal-passed routing evidence layout.
func agentRunnerFactsEvidence(facts AgentRunnerFrozenFacts) []string {
	return facts.hashes()
}

// frozenFactsFromEvidence rebuilds the frozen facts from the recorded
// terminal-passed routing evidence. Missing or ambiguous evidence fails closed:
// an unprovable fact set must never open the postflight gate.
func frozenFactsFromEvidence(input []string) (AgentRunnerFrozenFacts, error) {
	if len(input) < postflightFactsMinEvidence {
		return AgentRunnerFrozenFacts{}, fmt.Errorf("%w: routing evidence carries only %d digests", ErrInvalidPostflightFacts, len(input))
	}
	facts := AgentRunnerFrozenFacts{
		ManifestHash:            input[postflightFactsPrefix],
		DiffHash:                input[postflightFactsPrefix+1],
		LogHash:                 input[postflightFactsPrefix+2],
		UsageHash:               input[postflightFactsPrefix+3],
		ProcessStopEvidenceHash: input[postflightFactsPrefix+4],
		ProcessPidReportsHash:   input[postflightFactsPrefix+5],
	}
	if err := facts.Validate(); err != nil {
		return AgentRunnerFrozenFacts{}, err
	}
	return facts, nil
}

// agentRunnerTaskFacts returns the frozen facts a task's routed AgentRunner
// terminal left in the event log. A task whose steps never routed a completed
// terminal has no facts and keeps the pre-existing downstream path; a task whose
// terminals came from several steps is refused instead of guessing which facts
// the postflight should audit.
func agentRunnerTaskFacts(events []evidence.StateEvent, runID, taskID string, attempt int) (AgentRunnerFrozenFacts, bool, error) {
	var (
		facts  AgentRunnerFrozenFacts
		found  bool
		stepID string
	)
	for _, record := range events {
		event := record.Event
		if event.Entity.Kind != "step" || event.Entity.RunID != runID || event.Entity.TaskID != taskID {
			continue
		}
		// P43 N3: facts belong to the attempt that recorded them.
		if event.Entity.Attempt != attempt {
			continue
		}
		if event.Reason.Code != agentRunnerTerminalPassedReason {
			continue
		}
		rebuilt, err := frozenFactsFromEvidence(event.InputEvidence)
		if err != nil {
			return AgentRunnerFrozenFacts{}, false, err
		}
		if found && stepID != event.Entity.StepID {
			return AgentRunnerFrozenFacts{}, false, fmt.Errorf("%w: task %s routed completed terminals for several steps", ErrInvalidPostflightFacts, taskID)
		}
		facts, found, stepID = rebuilt, true, event.Entity.StepID
	}
	return facts, found, nil
}

// requirePostflightQualifiedReview proves that a task already resting in
// REVIEW_PENDING earned that state from the gate. A review state a previous
// build wrote without the gate cannot be honoured, because the gate is the only
// qualification for review; an unprovable one is refused with zero writes so the
// run is never resumed in place. A task without routed AgentRunner facts keeps
// the pre-existing path.
func (engine *Engine) requirePostflightQualifiedReview(ctx context.Context, taskID string, attempt int) error {
	events, err := engine.options.Events.Load(ctx)
	if err != nil {
		return err
	}
	agentFacts, agentFound, err := agentRunnerTaskFacts(events, engine.options.RunID, taskID, attempt)
	if err != nil {
		return err
	}
	supervisedFacts, supervisedFound, err := supervisedTaskFacts(events, engine.options.RunID, taskID, attempt)
	if err != nil {
		return err
	}
	if agentFound && supervisedFound {
		return fmt.Errorf("%w: task %s carries mixed postflight facts", ErrInvalidPostflightFacts, taskID)
	}
	if !agentFound && !supervisedFound {
		// A task whose steps never routed frozen facts keeps the pre-existing
		// path; a supervised task with no return facts can never have earned
		// REVIEW_PENDING at all.
		if task, ok := engine.findTask(taskID); ok && taskHasSupervisedStep(task) {
			return fmt.Errorf("%w: supervised task %s is in REVIEW_PENDING without a return fact set", ErrUnqualifiedReview, taskID)
		}
		return nil
	}
	var recorded *evidence.StateEvent
	for index := range events {
		event := events[index].Event
		if event.Entity.Kind != "task" || event.Entity.TaskID != taskID || event.ToState != "REVIEW_PENDING" {
			continue
		}
		recorded = &events[index]
	}
	if recorded == nil {
		return fmt.Errorf("%w: task %s is in REVIEW_PENDING without a recorded review transition", ErrUnqualifiedReview, taskID)
	}
	required := agentRunnerFactsEvidence(agentFacts)
	if supervisedFound {
		required = supervisedFactsEvidence(supervisedFacts)
	}
	for _, digest := range required {
		if !slices.Contains(recorded.Event.InputEvidence, digest) {
			return fmt.Errorf("%w: task %s reached REVIEW_PENDING without the postflight gate", ErrUnqualifiedReview, taskID)
		}
	}
	return nil
}

// runTaskPostflight runs the chain-owned postflight gate before a task may be
// reviewed. Only a passed postflight qualifies the task for REVIEW_PENDING, so
// an exit code, a completed receipt or an adapter record can never reach task
// acceptance on its own.
func (engine *Engine) runTaskPostflight(ctx context.Context, task Task, parent SnapshotRef, workspace Workspace, attempt int) ([]string, string, error) {
	entity := taskEntityAt(engine.options.RunID, task.ID, attempt)
	if engine.state.current(entity) != "STEPS_RUNNING" {
		return nil, "steps-completed", nil
	}
	events, err := engine.options.Events.Load(ctx)
	if err != nil {
		return nil, "", err
	}
	facts, found, err := agentRunnerTaskFacts(events, engine.options.RunID, task.ID, attempt)
	if err != nil {
		return nil, "", engine.failTask(ctx, task.ID, attempt, "postflight-facts-missing", err)
	}
	supervisedFacts, supervisedFound, err := supervisedTaskFacts(events, engine.options.RunID, task.ID, attempt)
	if err != nil {
		return nil, "", engine.failTask(ctx, task.ID, attempt, "supervised-postflight-facts-missing", err)
	}
	if found && supervisedFound {
		return nil, "", engine.failTask(ctx, task.ID, attempt, "postflight-mixed-facts", fmt.Errorf("%w: task %s mixes AgentRunner and supervised facts", ErrInvalidPostflightFacts, task.ID))
	}
	if supervisedFound {
		return engine.runSupervisedTaskPostflight(ctx, task, parent, workspace, supervisedFacts, attempt)
	}
	if taskHasSupervisedStep(task) {
		// A supervised task can never infer completion from a plain step pass:
		// without return facts the review state is not reachable (P3).
		return nil, "", engine.failTask(ctx, task.ID, attempt, "supervised-return-facts-missing", fmt.Errorf("%w: supervised task %s has no return fact set", ErrInvalidPostflightFacts, task.ID))
	}
	if !found {
		return nil, "steps-completed", nil
	}
	decision, err := engine.options.Postflight.RunPostflight(ctx, PostflightRequest{
		RunID:              engine.options.RunID,
		TaskID:             task.ID,
		Attempt:            attempt,
		ParentSnapshotHash: parent.Hash,
		Workspace:          workspace,
		Facts:              facts,
	})
	if err != nil {
		return nil, "", engine.failTask(ctx, task.ID, attempt, "postflight-failed", err)
	}
	if err := decision.Validate(facts); err != nil {
		return nil, "", engine.failTask(ctx, task.ID, attempt, "postflight-invalid", err)
	}
	inputs := uniqueHashes(append(append(agentRunnerFactsEvidence(facts), decision.Evidence...), decision.ErrorEvidence...))
	switch decision.Outcome {
	case PostflightFailed:
		return nil, "", engine.markTaskForRepair(ctx, task.ID, attempt, "postflight-rejected", inputs)
	case PostflightUncertain:
		if err := engine.transition(ctx, entity, "FAILED", inputs, "postflight-uncertain"); err != nil {
			return nil, "", err
		}
		if engine.state.projection.ChainState == "RUNNING" {
			if err := engine.transition(ctx, chainEntity(engine.options.RunID), "PAUSED", inputs, agentRunnerRecoveryUncertainReason); err != nil {
				return nil, "", err
			}
		}
		return nil, "", fmt.Errorf("%w: postflight uncertain for task %q", ErrRecoveryUncertain, task.ID)
	default:
		return inputs, "postflight-passed", nil
	}
}

// runSupervisedTaskPostflight runs the chain-owned postflight gate of a
// supervised-black-box task. Only a passed supervised postflight qualifies the
// task for REVIEW_PENDING, and the review transition carries the supervised
// facts - never the agent ones.
func (engine *Engine) runSupervisedTaskPostflight(ctx context.Context, task Task, parent SnapshotRef, workspace Workspace, facts SupervisedFrozenFacts, attempt int) ([]string, string, error) {
	if engine.options.SupervisedPostflight == nil {
		return nil, "", engine.failTask(ctx, task.ID, attempt, "supervised-postflight-unavailable", fmt.Errorf("%w: %w", ErrInvalidDefinition, ErrPostflightUnavailable))
	}
	decision, err := engine.options.SupervisedPostflight.RunSupervisedPostflight(ctx, SupervisedPostflightRequest{
		RunID:              engine.options.RunID,
		TaskID:             task.ID,
		Attempt:            attempt,
		ParentSnapshotHash: parent.Hash,
		Workspace:          workspace,
		Facts:              facts,
	})
	if err != nil {
		return nil, "", engine.failTask(ctx, task.ID, attempt, "supervised-postflight-failed", err)
	}
	if err := decision.Validate(facts); err != nil {
		return nil, "", engine.failTask(ctx, task.ID, attempt, "supervised-postflight-invalid", err)
	}
	inputs := uniqueHashes(append(append(supervisedFactsEvidence(facts), decision.Evidence...), decision.ErrorEvidence...))
	switch decision.Outcome {
	case PostflightFailed:
		return nil, "", engine.markTaskForRepair(ctx, task.ID, attempt, "supervised-postflight-rejected", inputs)
	case PostflightUncertain:
		if err := engine.transition(ctx, taskEntityAt(engine.options.RunID, task.ID, attempt), "FAILED", inputs, "supervised-postflight-uncertain"); err != nil {
			return nil, "", err
		}
		if engine.state.projection.ChainState == "RUNNING" {
			if err := engine.transition(ctx, chainEntity(engine.options.RunID), "PAUSED", inputs, agentRunnerRecoveryUncertainReason); err != nil {
				return nil, "", err
			}
		}
		return nil, "", fmt.Errorf("%w: supervised postflight uncertain for task %q", ErrRecoveryUncertain, task.ID)
	default:
		return inputs, "supervised-postflight-passed", nil
	}
}
