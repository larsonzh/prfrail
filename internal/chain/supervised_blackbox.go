package chain

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	// ErrAwaitingSupervisedConfirmation reports that no supervised work window
	// exists yet: the visible risk notice was delivered and the run waits for
	// the operator's binding confirmation. Nothing was written.
	ErrAwaitingSupervisedConfirmation = errors.New("awaiting supervised risk confirmation")
	// ErrAwaitingSupervisedReturn reports that the supervised work window is
	// open: an external agent may be writing inside the isolated workspace and
	// the chain waits for the operator to explicitly return it. Only a complete
	// handoff receipt can end that wait; a delivery or a chat message never can.
	ErrAwaitingSupervisedReturn = errors.New("awaiting supervised return")
)

// Step-transition reasons owned by the supervised-black-box mode.
const (
	supervisedWorkWindowOpenReason  = "supervised-work-window-open"
	supervisedReturnAcceptedReason  = "supervised-return-accepted"
	supervisedReturnCompletedReason = "supervised-return-completed"
)

// The frozen-fact layout inside a supervised step's PASSED evidence is
// positional: handoff receipt, post-return scan evidence, managed-process stop
// proof, independent gate evidence, reduced/unknown report. The chain rebuilds
// the facts from the recorded event evidence without reading any store; the
// mutual distinctness enforced by Validate keeps the positions unambiguous.
const (
	supervisedFactsPrefix      = 0
	supervisedFactsCount       = 5
	supervisedFactsMinEvidence = supervisedFactsPrefix + supervisedFactsCount
)

// Supervised risk notice. The version and the exact text are part of the
// binding: a confirmation that quotes another version or another hash is a
// notice drift and is refused before any work window opens. The text names the
// seven dimensions ProofRail cannot prove about an external agent, matching
// contracts section 7.1.
const (
	// SupervisedRiskNoticeVersion is the current notice version constant.
	SupervisedRiskNoticeVersion = 2
	supervisedRiskNoticeDomain  = "proofrail:supervised-risk-notice:1\n"
	// SupervisedRiskNoticeText enumerates what the mode cannot prove and what
	// the operator is responsible for. It never implies that any unknown
	// dimension is safe, and it never claims an isolation boundary ProofRail
	// does not build: the workspace is an ordinary directory ProofRail
	// materializes under the run root, not an OS-level sandbox.
	SupervisedRiskNoticeText = "supervised-black-box risk notice (version 2): ProofRail cannot prove which tools or commands the external agent used; which resources it read; whether it issued network requests; its actual token, cost and retry usage; every background process it started; the continuity of its session; or external side effects such as commit, push or publish. The workspace ProofRail hands to the agent is a plain directory it materialized under the run root, not an OS-level sandbox: nothing there confines filesystem, process or network reach, so the agent can act beyond it and ProofRail would not observe it. ProofRail guarantees only its own obligations: workspace materialization, post-return artifact scanning, independent gates and the acceptance transaction. The operator selects the trusted host, supervises the process, limits credentials and network access, and confirms the unknown external side effects. Unknown dimensions are never shown as safe."
)

// SupervisedRiskNoticeHash returns the digest the risk notice binds. Changing
// the notice text or the domain changes this digest, which turns a stale
// confirmation into a refused version drift.
func SupervisedRiskNoticeHash() string {
	return evidence.Digest(supervisedRiskNoticeDomain, []byte(SupervisedRiskNoticeText))
}

// SupervisedFrozenFacts are the five chain-owned digests a supervised return
// freezes. They carry no decision: a delivery status, a chat message, an exit
// code or an operator's verbal confirmation never belong here.
type SupervisedFrozenFacts struct {
	HandoffReceiptHash      string
	ScanEvidenceHash        string
	ProcessStopEvidenceHash string
	GateEvidenceHash        string
	ReportHash              string
}

// hashes returns the digests in the fixed routing order.
func (facts SupervisedFrozenFacts) hashes() []string {
	return []string{facts.HandoffReceiptHash, facts.ScanEvidenceHash, facts.ProcessStopEvidenceHash, facts.GateEvidenceHash, facts.ReportHash}
}

// supervisedFactsEvidence returns the fact digests in routing order, for the
// step PASSED evidence layout.
func supervisedFactsEvidence(facts SupervisedFrozenFacts) []string {
	return facts.hashes()
}

// Validate fails closed on missing digests and on digest collisions: the
// layout is positional, so repeated digests would make the rebuilt facts
// ambiguous.
func (facts SupervisedFrozenFacts) Validate() error {
	hashes := facts.hashes()
	for index, hash := range hashes {
		if !evidence.ValidHash(hash) {
			return fmt.Errorf("%w: supervised fact %d is not a valid hash", ErrInvalidPostflightFacts, index)
		}
	}
	for index, hash := range hashes {
		if slices.Contains(hashes[:index], hash) {
			return fmt.Errorf("%w: supervised fact %d repeats an earlier digest", ErrInvalidPostflightFacts, index)
		}
	}
	return nil
}

// DisjointFrom reports whether any frozen fact collides with a digest that
// already carries another meaning in the step history (the confirmation that
// unlocked the window, or the return receipt's own auxiliary evidence). This is
// the chain-level proof that a visible delivery or a risk confirmation can
// never masquerade as a frozen fact.
func (facts SupervisedFrozenFacts) DisjointFrom(hashes ...string) error {
	for _, fact := range facts.hashes() {
		if slices.Contains(hashes, fact) {
			return fmt.Errorf("%w: a supervised fact collides with a non-fact binding digest", ErrInvalidPostflightFacts)
		}
	}
	return nil
}

// supervisedFactsFromEvidence rebuilds the frozen facts from the recorded
// step PASSED evidence. Missing or ambiguous evidence fails closed: an
// unprovable fact set must never open the postflight gate.
func supervisedFactsFromEvidence(input []string) (SupervisedFrozenFacts, error) {
	if len(input) < supervisedFactsMinEvidence {
		return SupervisedFrozenFacts{}, fmt.Errorf("%w: supervised evidence carries only %d digests", ErrInvalidPostflightFacts, len(input))
	}
	facts := SupervisedFrozenFacts{
		HandoffReceiptHash:      input[supervisedFactsPrefix],
		ScanEvidenceHash:        input[supervisedFactsPrefix+1],
		ProcessStopEvidenceHash: input[supervisedFactsPrefix+2],
		GateEvidenceHash:        input[supervisedFactsPrefix+3],
		ReportHash:              input[supervisedFactsPrefix+4],
	}
	if err := facts.Validate(); err != nil {
		return SupervisedFrozenFacts{}, err
	}
	return facts, nil
}

// supervisedTaskFacts returns the frozen facts a task's supervised return left
// in the event log, rebuilt from the recorded step PASSED evidence. A task
// whose supervised step never completed a return has no facts; a task whose
// returns came from several steps, or several receipts for one step, is refused
// instead of guessing which facts the postflight should audit.
func supervisedTaskFacts(events []evidence.StateEvent, runID, taskID string) (SupervisedFrozenFacts, bool, error) {
	var (
		facts SupervisedFrozenFacts
		found bool
		seen  = map[string]int{}
		order []string
	)
	for _, record := range events {
		event := record.Event
		if event.Entity.Kind != "step" || event.Entity.RunID != runID || event.Entity.TaskID != taskID {
			continue
		}
		if event.Reason.Code != supervisedReturnCompletedReason {
			continue
		}
		if seen[event.Entity.StepID] > 0 {
			return SupervisedFrozenFacts{}, false, fmt.Errorf("%w: task %s recorded several supervised receipts for step %s", ErrInvalidPostflightFacts, taskID, event.Entity.StepID)
		}
		seen[event.Entity.StepID]++
		order = append(order, event.Entity.StepID)
		if len(order) > 1 {
			return SupervisedFrozenFacts{}, false, fmt.Errorf("%w: task %s completed supervised returns for several steps", ErrInvalidPostflightFacts, taskID)
		}
		rebuilt, err := supervisedFactsFromEvidence(event.InputEvidence)
		if err != nil {
			return SupervisedFrozenFacts{}, false, err
		}
		facts, found = rebuilt, true
	}
	return facts, found, nil
}

// supervisedWindowBindingDigests returns the digests recorded when the
// supervised work window was opened (the binding risk confirmation hash). They
// must never be reused as frozen facts.
func supervisedWindowBindingDigests(events []evidence.StateEvent, entity evidence.Entity) []string {
	var digests []string
	for _, record := range events {
		if record.Event.Entity != entity || record.Event.ToState != "WAITING_FOR_OPERATOR" {
			continue
		}
		digests = append(digests, record.Event.InputEvidence...)
	}
	return uniqueHashes(digests)
}

// OpenSupervisedWorkWindow opens the supervised-black-box work window. It is
// the only entry that may move a supervised step from RUNNING to the waiting
// state, and it refuses to write anything unless a valid binding risk
// confirmation fixes this run, task, step, attempt and workspace, quotes the
// current risk notice, and the chain is already paused for the operator.
func (engine *Engine) OpenSupervisedWorkWindow(ctx context.Context, session HandoffSession, confirmation SupervisedRiskConfirmationRecord) error {
	if err := ValidateSupervisedRiskConfirmation(confirmation); err != nil {
		return err
	}
	if err := validateHandoffSession(session); err != nil {
		return err
	}
	body := confirmation.Record
	if body.RunID != engine.options.RunID || session.RunID != engine.options.RunID {
		return fmt.Errorf("%w: confirmation does not belong to this run", ErrInvalidSupervisedConfirmation)
	}
	if body.TaskID != session.TaskID || body.StepID != session.StepID || body.Attempt != session.Attempt {
		return fmt.Errorf("%w: confirmation does not bind the handoff session identity", ErrInvalidSupervisedConfirmation)
	}
	if body.WorkspaceHash != session.BeforeManifestHash {
		return fmt.Errorf("%w: confirmation does not bind the handoff workspace digest", ErrInvalidSupervisedConfirmation)
	}
	// Re-asserted defensively: Validate already constrains the selection, but
	// the window opener is the boundary that must never open on a rejection.
	if body.Selection != "confirm" {
		return fmt.Errorf("%w: risk confirmation selection must be confirm", ErrInvalidSupervisedConfirmation)
	}
	// Contracts section 7.1a: a non-accepted delivery must never open the write
	// window. The record shape deliberately permits the unavailable/unknown
	// statuses so that a delivery fact can be recorded truthfully, but the
	// window opener is the place where that truth becomes a refusal.
	if body.Delivery.DeliveryStatus != supervisedDeliveryAccepted {
		return fmt.Errorf("%w: risk confirmation delivery status must be accepted", ErrInvalidSupervisedConfirmation)
	}
	if body.RiskNoticeVersion != SupervisedRiskNoticeVersion || body.RiskNoticeHash != SupervisedRiskNoticeHash() {
		return fmt.Errorf("%w: risk notice version or hash drifted", ErrInvalidSupervisedConfirmation)
	}
	task := taskEntity(engine.options.RunID, body.TaskID)
	step := stepEntity(engine.options.RunID, body.TaskID, body.StepID)
	taskState, stepState := engine.state.current(task), engine.state.current(step)
	if taskState != "STEPS_RUNNING" || stepState != "RUNNING" {
		return fmt.Errorf("%w: supervised work window requires STEPS_RUNNING/RUNNING", ErrInvalidState)
	}
	if engine.state.projection.ChainState != "PAUSED" {
		return fmt.Errorf("%w: supervised work window requires a paused chain", ErrInvalidState)
	}
	nextTask, nextStep, err := OpenHandoffTransition(taskState, stepState)
	if err != nil {
		return err
	}
	inputs := []string{confirmation.RecordHash}
	if err := engine.transition(ctx, task, nextTask, inputs, supervisedWorkWindowOpenReason); err != nil {
		return err
	}
	return engine.transition(ctx, step, nextStep, inputs, supervisedWorkWindowOpenReason)
}

// SubmitSupervisedReturn consumes the operator's explicit return. It is the
// only entry that may pass a supervised step, and it requires a complete
// handoff receipt plus five distinct frozen facts that are disjoint from the
// window binding. A delivery, a confirmation or a chat message can never pass a
// supervised step (P2). A replay with an equal fact set converges; a replay
// with a different fact set is refused.
func (engine *Engine) SubmitSupervisedReturn(ctx context.Context, receipt HandoffReceiptRecord, facts SupervisedFrozenFacts) error {
	if receipt.Receipt.Outcome != "complete" {
		return fmt.Errorf("%w: supervised return requires a complete handoff receipt", ErrInvalidState)
	}
	if err := facts.Validate(); err != nil {
		return err
	}
	// The handoff receipt itself is one of the five frozen facts: the step's
	// PASSED evidence must carry the very receipt that proves the return.
	if facts.HandoffReceiptHash != receipt.ReceiptHash {
		return fmt.Errorf("%w: the handoff receipt fact must be the return receipt digest", ErrInvalidPostflightFacts)
	}
	task := taskEntity(engine.options.RunID, receipt.Receipt.TaskID)
	step := stepEntity(engine.options.RunID, receipt.Receipt.TaskID, receipt.Receipt.StepID)
	events, err := engine.options.Events.Load(ctx)
	if err != nil {
		return err
	}
	if err := facts.DisjointFrom(supervisedWindowBindingDigests(events, step)...); err != nil {
		return err
	}
	stepState, taskState := engine.state.current(step), engine.state.current(task)
	if stepState == "PASSED" {
		// Idempotent convergence: the same return already passed this step.
		if !latestEntityEventMatches(events, step, "PASSED", supervisedReturnCompletedReason, supervisedFactsEvidence(facts)...) {
			return fmt.Errorf("%w: step %s is PASSED with a different supervised fact set", ErrInvalidPostflightFacts, receipt.Receipt.StepID)
		}
		return engine.convergeSupervisedReturn(ctx, task, receipt, facts, true)
	}
	if stepState == "RUNNING" {
		// An interrupted return: the step was already lifted back to RUNNING.
		if !latestEntityEventMatches(events, step, "RUNNING", supervisedReturnAcceptedReason, receipt.ReceiptHash) {
			return fmt.Errorf("%w: step %s is RUNNING without a proven supervised return", ErrInvalidState, receipt.Receipt.StepID)
		}
		return engine.convergeSupervisedReturn(ctx, task, receipt, facts, false)
	}
	if stepState != "WAITING_FOR_OPERATOR" || taskState != "WAITING_FOR_OPERATOR" {
		return fmt.Errorf("%w: supervised return requires a waiting task and step", ErrInvalidState)
	}
	_, nextStep, err := ReturnHandoffTransition(taskState, stepState, "complete")
	if err != nil {
		return err
	}
	if err := engine.transition(ctx, step, nextStep, []string{receipt.ReceiptHash}, supervisedReturnAcceptedReason); err != nil {
		return err
	}
	return engine.convergeSupervisedReturn(ctx, task, receipt, facts, false)
}

// convergeSupervisedReturn writes the step PASSED transition (the frozen fact
// layout) and lifts the task back to STEPS_RUNNING. When the step was already
// passed the fact evidence is not rewritten, so a replay after a partial write
// repairs exactly the missing task transition.
func (engine *Engine) convergeSupervisedReturn(ctx context.Context, task evidence.Entity, receipt HandoffReceiptRecord, facts SupervisedFrozenFacts, stepAlreadyPassed bool) error {
	if !stepAlreadyPassed {
		if err := engine.transition(ctx, stepEntity(engine.options.RunID, receipt.Receipt.TaskID, receipt.Receipt.StepID), "PASSED", supervisedFactsEvidence(facts), supervisedReturnCompletedReason); err != nil {
			return err
		}
	}
	if engine.state.current(task) == "WAITING_FOR_OPERATOR" {
		if err := engine.transition(ctx, task, "STEPS_RUNNING", []string{receipt.ReceiptHash}, supervisedReturnAcceptedReason); err != nil {
			return err
		}
	}
	return nil
}
