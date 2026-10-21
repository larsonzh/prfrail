package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/tickets"
)

var (
	// ErrAgentRunnerDispatchUnknownBlock reports durable dispatch state that
	// cannot authorize a launch: a request record exists without a terminal
	// receipt, or the single launch slot is already owned by a prior attempt.
	ErrAgentRunnerDispatchUnknownBlock = errors.New("AgentRunner dispatch unknown block")
	// ErrAgentRunnerDispatchTerminalReceiptPresent reports a terminal completion
	// already present for the request; any completion status only proves
	// terminal-receipt presence and never implies task PASS.
	ErrAgentRunnerDispatchTerminalReceiptPresent = errors.New("AgentRunner dispatch terminal receipt present")
	// ErrAgentRunnerLaunchFailed reports a start-port failure whose contract
	// guarantees no process was spawned. Callers must not retry a launch that
	// returned a result, and the retained launch intent must never be deleted to
	// "retry" a failed launch.
	ErrAgentRunnerLaunchFailed = errors.New("AgentRunner launch failed")
	// ErrAgentRunnerLaunchIdentityUnproven reports a spawn whose process identity
	// could not be published durably. The dispatcher reclaims the started process
	// through the launcher's stop capability when one exists and never retries or
	// relaunches in this state; settlement of the retained durable records is
	// owned by a later slice.
	ErrAgentRunnerLaunchIdentityUnproven = errors.New("AgentRunner launch identity unproven")
)

// AgentRunnerProcessConfig pins the offline process parameters the dispatcher
// hands to the launcher. It is additive to A2 and optional there: a dispatcher
// without it still dispatches, and the launcher fails closed on a request it
// cannot honour, so no replay or receipt semantics change.
type AgentRunnerProcessConfig struct {
	Command string
	Args    []string
	Env     []string
	Grace   time.Duration
}

// AgentRunnerReplayDispatcher owns the replay-aware dispatch path: publish R,
// reconfirm authorization and budget, claim the single launch slot, spawn, and
// publish the process identity. It implements chain.AgentRunnerPort and must be
// invoked after admission; it never launches when admission was skipped, and it
// never treats an admission verdict as a launch permit.
type AgentRunnerReplayDispatcher struct {
	Store          *AgentRunnerReplayStore
	Launcher       chain.AgentRunnerLauncher
	Process        AgentRunnerProcessConfig
	LedgerSnapshot func() (chain.AuthorizationLedger, error)
	CostLedger     *tickets.CostLedger
	RequestRecord  AgentRunnerRequestRecord
	Clock          func() time.Time
}

var _ chain.AgentRunnerPort = (*AgentRunnerReplayDispatcher)(nil)

// agentRunnerLaunchReclaimer is the optional capability a launcher may expose
// for stopping a process it already started. It is not a chain port and not a
// wire contract: a launcher without it keeps the pre-OB-57 behaviour, and the
// dispatcher still never touches an OS process directly.
// *AgentRunnerPinnedCLILauncher implements it.
type agentRunnerLaunchReclaimer interface {
	StopAgentRunnerProcess(ctx context.Context, launchID string) (guard.TerminationEvidence, error)
}

// reclaimUnprovenLaunch best-effort stops the process a successful
// StartAgentRunnerProcess left behind when the replay identity cannot be
// published. It never replaces the cause, never retries the stop, and reports
// explicitly when no reclaim happened.
func (dispatcher *AgentRunnerReplayDispatcher) reclaimUnprovenLaunch(ctx context.Context, launchID string, cause error) error {
	reclaimer, ok := dispatcher.Launcher.(agentRunnerLaunchReclaimer)
	// Defensive: unreachable today, but guards against a future Launcher that
	// stops implementing agentRunnerLaunchReclaimer. Any change that makes this
	// reachable MUST add a test.
	if !ok || reclaimer == nil {
		return fmt.Errorf("%w: the launcher exposes no stop capability; the spawned process was not reclaimed", cause)
	}
	grace := dispatcher.Process.Grace
	if grace <= 0 {
		grace = time.Second
	}
	stopCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), grace+time.Second)
	defer cancel()
	proof, stopErr := reclaimer.StopAgentRunnerProcess(stopCtx, launchID)
	switch {
	case stopErr != nil:
		return fmt.Errorf("%w: spawned process reclaim failed: %v", cause, stopErr)
	case proof.Outcome != "stopped":
		return fmt.Errorf("%w: spawned process reclaim unproven: outcome %q", cause, proof.Outcome)
	default:
		return fmt.Errorf("%w: spawned process reclaimed (%s)", cause, proof.Outcome)
	}
}

// DispatchAgentRunner performs the first-dispatch path for one step request.
// Replays never relaunch: R-only states, existing launch intents, and terminal
// states all block with errors.Is-identifiable outcomes, and the dispatcher
// publishes no completion.
func (dispatcher *AgentRunnerReplayDispatcher) DispatchAgentRunner(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	if dispatcher == nil {
		return chain.StepResult{}, fmt.Errorf("%w: nil dispatcher", ErrInvalidAgentRunnerRequest)
	}
	if err := ctx.Err(); err != nil {
		return chain.StepResult{}, fmt.Errorf("%w: dispatch context cancelled: %w", ErrInvalidAgentRunnerRequest, err)
	}
	if dispatcher.Store == nil {
		return chain.StepResult{}, fmt.Errorf("%w: replay store missing", ErrInvalidAgentRunnerRequest)
	}
	if dispatcher.Launcher == nil {
		return chain.StepResult{}, fmt.Errorf("%w: launcher missing", ErrInvalidAgentRunnerRequest)
	}
	if dispatcher.LedgerSnapshot == nil {
		return chain.StepResult{}, fmt.Errorf("%w: authorization ledger snapshot missing", ErrInvalidAgentRunnerRequest)
	}
	if dispatcher.Clock == nil {
		return chain.StepResult{}, fmt.Errorf("%w: evaluation clock missing", ErrInvalidAgentRunnerRequest)
	}
	if err := ValidateAgentRunnerRequestRecord(dispatcher.RequestRecord); err != nil {
		return chain.StepResult{}, err
	}
	if err := dispatcher.validateRequestBinding(request); err != nil {
		return chain.StepResult{}, err
	}
	body := dispatcher.RequestRecord.Request
	decision, err := dispatcher.Store.RecordRequest(dispatcher.RequestRecord)
	if err != nil {
		return chain.StepResult{}, err
	}
	switch decision {
	case AgentRunnerReplayDecisionTerminalReceiptPresent:
		return chain.StepResult{}, fmt.Errorf("%w: requestId %s", ErrAgentRunnerDispatchTerminalReceiptPresent, body.RequestID)
	case AgentRunnerReplayDecisionUnknownBlock:
		return chain.StepResult{}, fmt.Errorf("%w: requestId %s", ErrAgentRunnerDispatchUnknownBlock, body.RequestID)
	case AgentRunnerReplayDecisionFirstDispatch:
	default:
		return chain.StepResult{}, fmt.Errorf("%w: unknown replay decision %s", ErrInvalidAgentRunnerRequest, decision)
	}
	now := dispatcher.Clock().UTC()
	if now.IsZero() {
		return chain.StepResult{}, fmt.Errorf("%w: evaluation clock returned zero time", ErrAgentRunnerLaunchReconfirmation)
	}
	ledger, err := dispatcher.LedgerSnapshot()
	if err != nil {
		return chain.StepResult{}, fmt.Errorf("%w: authorization ledger snapshot unavailable: %w", ErrAgentRunnerLaunchReconfirmation, err)
	}
	reconfirmer := AgentRunnerLaunchReconfirmer{
		RequestRecord:       dispatcher.RequestRecord,
		AuthorizationLedger: ledger,
		CostLedger:          dispatcher.CostLedger,
		Clock:               func() time.Time { return now },
	}
	if err := reconfirmer.ReconfirmLaunchAuthorization(ctx); err != nil {
		return chain.StepResult{}, err
	}
	launchID, err := agentRunnerLaunchID(body.RequestID)
	if err != nil {
		return chain.StepResult{}, err
	}
	receipt, err := NewAgentRunnerLaunchReceiptRecord(AgentRunnerLaunchReceipt{
		LaunchID:          launchID,
		RequestID:         body.RequestID,
		RequestHash:       dispatcher.RequestRecord.RecordHash,
		RunID:             body.RunID,
		TaskID:            body.TaskID,
		StepID:            body.StepID,
		Attempt:           body.Attempt,
		AdapterID:         body.AdapterID,
		ReconfirmedAt:     now.Format(time.RFC3339),
		AuthorizationHash: body.AuthorizationHash,
		BudgetHash:        body.BudgetHash,
	})
	if err != nil {
		return chain.StepResult{}, err
	}
	launchDecision, err := dispatcher.Store.RecordLaunchReceipt(receipt)
	if err != nil {
		return chain.StepResult{}, err
	}
	if launchDecision != AgentRunnerLaunchDecisionFirstLaunch {
		return chain.StepResult{}, fmt.Errorf("%w: launch slot already owned for requestId %s", ErrAgentRunnerDispatchUnknownBlock, body.RequestID)
	}
	result, err := dispatcher.Launcher.StartAgentRunnerProcess(ctx, chain.AgentRunnerLaunchRequest{
		RequestID:     body.RequestID,
		RunID:         body.RunID,
		AdapterID:     body.AdapterID,
		Command:       dispatcher.Process.Command,
		Args:          append([]string(nil), dispatcher.Process.Args...),
		Dir:           request.Workspace.Root,
		Env:           append([]string(nil), dispatcher.Process.Env...),
		Grace:         dispatcher.Process.Grace,
		WorkspaceRoot: request.Workspace.Root,
	})
	if err != nil {
		return chain.StepResult{}, fmt.Errorf("%w: %w", ErrAgentRunnerLaunchFailed, err)
	}
	if result.LaunchID != launchID {
		return chain.StepResult{}, dispatcher.reclaimUnprovenLaunch(ctx, launchID,
			fmt.Errorf("%w: launcher returned launchId %s, want %s", ErrAgentRunnerLaunchIdentityUnproven, result.LaunchID, launchID))
	}
	identity, err := NewAgentRunnerLaunchIdentityRecord(AgentRunnerLaunchIdentity{
		LaunchID:         launchID,
		RequestID:        body.RequestID,
		IntentRecordHash: receipt.RecordHash,
		ProcessID:        result.ProcessID,
		StartedAt:        result.StartedAt.UTC().Format(time.RFC3339),
	})
	if err != nil {
		return chain.StepResult{}, dispatcher.reclaimUnprovenLaunch(ctx, launchID,
			fmt.Errorf("%w: %w", ErrAgentRunnerLaunchIdentityUnproven, err))
	}
	if _, err := dispatcher.Store.RecordLaunchIdentity(identity); err != nil {
		return chain.StepResult{}, dispatcher.reclaimUnprovenLaunch(ctx, launchID,
			fmt.Errorf("%w: %w", ErrAgentRunnerLaunchIdentityUnproven, err))
	}
	return chain.StepResult{Evidence: []string{dispatcher.RequestRecord.RecordHash, receipt.RecordHash, identity.RecordHash}}, nil
}

func (dispatcher *AgentRunnerReplayDispatcher) validateRequestBinding(request chain.StepRequest) error {
	facts := request.AgentRunnerFacts
	if facts == nil {
		return fmt.Errorf("%w: immutable facts missing", ErrInvalidAgentRunnerRequest)
	}
	if !evidence.ValidID(facts.RequestID) ||
		!evidence.ValidHash(facts.WorkspaceHash) ||
		!evidence.ValidHash(facts.ContextHash) ||
		!evidence.ValidHash(facts.AuthorizationHash) ||
		!evidence.ValidHash(facts.BudgetHash) {
		return fmt.Errorf("%w: immutable facts contain an invalid binding", ErrInvalidAgentRunnerRequest)
	}
	body := dispatcher.RequestRecord.Request
	if request.Step.Kind != "code" || request.Step.Mode != chain.IsolatedWorkspace || request.ExecutionTarget != chain.AgentRunnerExecution {
		return fmt.Errorf("%w: AgentRunner requires isolated workspace code step", ErrInvalidAgentRunnerRequest)
	}
	if body.Mode != "create" {
		return fmt.Errorf("%w: resume continuity evidence is not implemented", ErrInvalidAgentRunnerRequest)
	}
	if request.RunID != body.RunID ||
		request.TaskID != body.TaskID ||
		request.Step.ID != body.StepID ||
		request.Attempt != body.Attempt ||
		request.ParentHash != body.ParentSnapshotHash {
		return fmt.Errorf("%w: step request binding mismatch", ErrInvalidAgentRunnerRequest)
	}
	if facts.RequestID != body.RequestID ||
		facts.WorkspaceHash != body.WorkspaceHash ||
		facts.ContextHash != body.ContextHash ||
		facts.AuthorizationHash != body.AuthorizationHash ||
		facts.BudgetHash != body.BudgetHash {
		return fmt.Errorf("%w: immutable facts binding mismatch", ErrInvalidAgentRunnerRequest)
	}
	return nil
}

func agentRunnerLaunchID(requestID string) (string, error) {
	launchID := "launch-" + requestID
	if !evidence.ValidID(launchID) {
		return "", fmt.Errorf("%w: derived launchId %q is invalid", ErrInvalidAgentRunnerLaunchReceipt, launchID)
	}
	return launchID, nil
}
