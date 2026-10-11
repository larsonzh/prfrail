package console

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/gates"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// This file is the production supervised-black-box assembly of the console
// package. It adds no chain state and no chain transition: the mode is driven
// entirely through the existing operator-interaction entries
// (OpenOperatorInteraction / ResumeOperatorInteraction) plus the supervised
// window entries (OpenSupervisedWorkWindow / SubmitSupervisedReturn). It never
// lets a visible delivery, a chat message or an operator confirmation reach a
// terminal state on its own.
//
// # How the run is driven without the step dispatch
//
// A supervised step routed through StepRouter.Execute inside engine.Run would
// make the engine treat the port's ErrAwaitingSupervisedConfirmation sentinel
// as a step failure (step-failed -> failTask). The assembly therefore parks the
// run from inside the engine's existing failpoint seam: the failpoint is
// invoked after the step reached RUNNING and before any dispatch happens, so it
// can deliver the notice, persist the binding interaction and pause the run
// while leaving a clean STEPS_RUNNING/RUNNING/PAUSED frontier. engine.Run is
// only ever used again for the tail chain (postflight -> review -> publish).
//
// # Disk state machine
//
// Every later invocation decides its phase from the run root alone:
//
//	phase A  no interaction request persisted  -> deliver, open the interaction, park
//	phase B  request persisted, no response    -> waiting for the operator
//	phase C  response persisted, no confirmation -> confirm, open the work window
//	phase D  window open, no return receipt    -> waiting for the explicit return
//	phase E  chain terminal                   -> nothing left to do
//
// Phase D is state-aware: a durable confirmation with no persisted session means
// the process died inside the confirm/resume/open/persist sequence, so the
// existing confirmation is replayed (never rebuilt) until the window and the
// session agree. A replay that cannot converge returns a structured error
// instead of stalling on the non-empty ledger.
//
// The residual is the window between persisting the request and writing the
// pause transition (plan section 9.7): a crash inside it leaves a RUNNING step
// whose re-entry fails closed as recovery-uncertain. The ordering is
// persist-then-transition on purpose; PILOT does not claim power-cut atomicity.

// The supervised store-local layout under <runRoot>/supervised/. These files are
// store-local protocol, not wire records: no schema and no external consumer
// depends on them.
const (
	supervisedDirectoryName          = "supervised"
	supervisedInteractionsName       = "operator-interactions.jsonl"
	supervisedConfirmationsName      = "risk-confirmations.jsonl"
	supervisedDeliveryName           = "delivery.json"
	supervisedSessionName            = "session.json"
	supervisedFrozenPolicyName       = "handoff-policy.json"
	supervisedReceiptName            = "handoff-receipt.json"
	supervisedReportName             = "report.json"
	supervisedScansDirectoryName     = "scans"
	supervisedBeforeManifestName     = "before-manifest.json"
	supervisedScanReportName         = "scan-report.json"
	supervisedGateEvidenceName       = "gate-evidence.json"
	supervisedSnapshotStoreDirectory = "snapshot-store"
	supervisedWorkspacesDirectory    = "workspaces"
	supervisedConfirmationResponse   = "confirm"
	supervisedRejectionResponse      = "reject"
	supervisedDeliveryAcceptedStatus = "accepted"
	supervisedStepFailpoint          = "step-running"
	supervisedAdapterID              = "sessbridge-visible"
	supervisedOperatorID             = "local-operator"
	supervisedRequiredActionText     = "supervise the external agent, then run: prfrail supervised return --run-id <id> --outcome complete|abort"
	supervisedRiskDescription        = "reduced-assurance-external-agent"
	supervisedReviewPolicyLabel      = "supervised-assembly-policy"
)

// Console-local evidence domains. They are internal digests, never wire records.
const (
	supervisedHandoffPolicyDomain         = "proofrail:supervised-handoff-policy:1\n"
	supervisedSessionLeaseDomain          = "proofrail:supervised-session-lease:1\n"
	supervisedConfirmationEvidenceDomain  = "proofrail:supervised-confirmation-evidence:1\n"
	supervisedContextDomain               = "proofrail:supervised-window-context:1\n"
	supervisedScanDomain                  = "proofrail:supervised-scan:1\n"
	supervisedProcessStopDomain           = "proofrail:supervised-process-stop:1\n"
	supervisedGateEvidenceDomain          = "proofrail:supervised-gate-evidence:1\n"
	supervisedReportEvidenceDomain        = "proofrail:supervised-report-evidence:1\n"
	supervisedDiffDomain                  = "proofrail:supervised-diff:1\n"
	supervisedCandidateEvidenceRootDomain = "proofrail:supervised-candidate-evidence-root:1\n"
	supervisedCandidateEvidenceDomain     = "proofrail:supervised-candidate-evidence:1\n"
	supervisedReviewPolicyDomain          = "proofrail:supervised-review-policy:1\n"
	supervisedReviewEvidenceDomain        = "proofrail:supervised-review-evidence:1\n"
	supervisedStopperDomain               = "proofrail:supervised-stopper:1\n"
	supervisedWriterStopDomain            = "proofrail:supervised-writer-stop:1\n"
	supervisedLeaseDomain                 = "proofrail:supervised-lease:1\n"
	supervisedPromotionEvidenceDomain     = "proofrail:supervised-promotion-evidence:1\n"
	supervisedPromotionBlockedDomain      = "proofrail:supervised-promotion-blocked:1\n"
	supervisedMissingProcessReason        = "no managed process was started by this run"
)

// The durable boundaries of the confirmation sequence, in order. A process that
// dies at one of them leaves a non-empty confirmation ledger with no persisted
// session; resumeSupervisedRun detects that and converges.
const (
	supervisedCrashAfterConfirmation = "after-confirmation"
	supervisedCrashAfterResume       = "after-resume"
	supervisedCrashAfterWindowOpen   = "after-window-open"
)

// The frozen closed sets mirrored from internal/chain. HAND-SYNCHRONISED
// COUPLING (unknown set): internal/chain owns the authoritative set
// (supervisedUnknownDimensionSet); the two copies must stay equal and must be
// re-evaluated together. The copy exists because the console must build a
// report without exporting the chain's unexported set; ValidateSupervisedReport
// still refuses anything outside the chain's own set, so a drift can only ever
// fail closed.
//
// supervisedVerifiedDimensions is deliberately NOT the chain's full closed set:
// it is the subset proven at the moment the report is frozen. The report is
// written before the tail chain (postflight -> review -> publish) runs, so
// "independent-review" is not yet true at write time and is omitted rather than
// asserted (P4: never present the unproven as proven). The chain's closed set
// still accepts it, so a later slice may record it once a post-review write
// point exists.
var (
	supervisedUnknownDimensions  = []string{"tool-use", "resource-reads", "network-access", "cost", "background-processes", "session-continuity", "external-side-effects"}
	supervisedVerifiedDimensions = []string{"workspace-isolation", "manifest-before", "manifest-after", "diff", "scope-scan", "secret-scan", "file-type-scan", "side-effect-scan", "process-stop-proof", "independent-gates"}
)

var (
	// errSupervisedParked is the console marker the failpoint returns to park the
	// run in the waiting state. It is never a step failure.
	errSupervisedParked = errors.New("supervised run parked awaiting operator confirmation")
	// ErrSupervisedStepRequired reports a definition that is not exactly one
	// supervised-black-box code step.
	ErrSupervisedStepRequired = errors.New("the supervised assembly requires exactly one supervised-black-box code step")
	// ErrSupervisedDeliveryBlocked reports a visible delivery that was not
	// accepted. The work window must never open.
	ErrSupervisedDeliveryBlocked = errors.New("supervised visible delivery is not accepted")
	// ErrSupervisedRejected reports an operator response that did not confirm.
	ErrSupervisedRejected = errors.New("supervised risk confirmation rejected")
	// ErrSupervisedTimeout reports a work window that outlived its frozen handoff
	// timeout: a complete return is refused and only abort remains possible.
	ErrSupervisedTimeout = errors.New("supervised handoff window timed out; only abort is permitted")
	// ErrSupervisedPolicyDrift reports a complete return attempted under a handoff
	// policy that no longer matches the policy frozen when the work window
	// opened. A changed chain configuration (or a changed config path) would
	// otherwise let the return enforce a swapped timeout, authorized scope or
	// after-return gate set, so the guard fails closed before any receipt, report
	// or chain transition is written.
	ErrSupervisedPolicyDrift = errors.New("supervised handoff policy drifted since the window opened")
	// ErrSupervisedGateUnavailable reports a missing after-return gate runner or
	// hook definition. The gate is not optional, so the return fails closed.
	ErrSupervisedGateUnavailable = errors.New("supervised after-return gate is unavailable")
	// ErrSupervisedGateFailed reports a failed after-return gate hook. A
	// confirmation never exempts a gate.
	ErrSupervisedGateFailed = errors.New("supervised after-return gate hook failed")
	// ErrSupervisedSecretFinding reports credential material in the returning
	// workspace.
	ErrSupervisedSecretFinding = errors.New("supervised post-return secret scan found credential material")
	// ErrSupervisedNotOpen reports a return attempted without a persisted window.
	ErrSupervisedNotOpen = errors.New("supervised work window is not open")
	// ErrSupervisedAbortOnly reports an unknown return outcome.
	ErrSupervisedAbortOnly = errors.New("supervised return outcome must be complete or abort")
	// ErrSupervisedScopeUnderivable reports an authorized scope that cannot be
	// derived from the frozen handoff policy (no declared allowed target, or an
	// unresolvable workspace root). An underivable scope is never read as "no
	// restriction": it fails closed before any comparison happens.
	ErrSupervisedScopeUnderivable = errors.New("supervised authorized scope cannot be derived")
	// ErrSupervisedScopeOutOfBounds reports a post-return change whose path leaves
	// the authorized write scope. ProofRail has both the duty and the ability to
	// prove this, so an observed out-of-scope path blocks the return (plan
	// section 8: detectable => block).
	ErrSupervisedScopeOutOfBounds = errors.New("supervised post-return change set leaves the authorized scope")
	// ErrSupervisedSideEffect reports an observable side effect on the isolated
	// workspace: a deletion of a pre-existing entry, a type change, a
	// permission/executable change, or a link whose target leaves the workspace
	// root.
	ErrSupervisedSideEffect = errors.New("supervised post-return side effect detected")
	// ErrSupervisedFileType reports an entry type the return scan does not accept.
	ErrSupervisedFileType = errors.New("supervised post-return file type is not accepted")
	// ErrSupervisedWindowUnproven reports a re-entry that carries a durable
	// confirmation but neither an open window nor a replayable resume: the run
	// cannot converge and the operator must inspect the run manually.
	ErrSupervisedWindowUnproven = errors.New("supervised confirmed window could not be reopened")
)

// SupervisedAssemblyOptions pins one production supervised-black-box run. Every
// record it consumes is already declared; the assembly never probes a candidate
// and never touches the network.
//
// GateHooks is additive to the plan's field list: the after-return gate hook
// definitions are supplied by the caller because the chain config only carries
// hook identifiers (StepConfig.Hooks is []string), so the console cannot derive
// a runnable hook definition from the config alone.
type SupervisedAssemblyOptions struct {
	RunID         string
	RunRoot       string
	Definition    chain.Definition
	VisibleSender adapters.VisibleClient
	GateRunner    *gates.Runner
	GateHooks     map[string]gates.Hook
	Clock         func() time.Time
	IDs           func() string
	// CrashPoint, when set, is called at each durable boundary of the
	// confirmation sequence (supervisedCrashAfterConfirmation /
	// supervisedCrashAfterResume / supervisedCrashAfterWindowOpen). It exists so
	// the crash-recovery tests can stop the sequence at an exact boundary and
	// prove that a re-entry converges; a non-nil error aborts the invocation.
	// Production callers leave it nil.
	CrashPoint func(point string) error
}

// SupervisedDir returns the store-local supervised directory of a run root.
func SupervisedDir(runRoot string) string { return filepath.Join(runRoot, supervisedDirectoryName) }

// SupervisedInteractionsLedgerPath returns the operator-interaction ledger the
// supervised assembly reads and appends. The operator answers the binding risk
// confirmation through the existing `prfrail interactions respond` command with
// `--ledger` pointing at this path.
func SupervisedInteractionsLedgerPath(runRoot string) string {
	return filepath.Join(SupervisedDir(runRoot), supervisedInteractionsName)
}

// SupervisedReportPath returns the persisted reduced/unknown report path.
func SupervisedReportPath(runRoot string) string {
	return filepath.Join(SupervisedDir(runRoot), supervisedReportName)
}

// ExecuteSupervisedRun is the production supervised-black-box entry. It drives
// the disk state machine described at the top of this file and returns the
// current projection. Reaching the waiting state is not an error: the run is
// legitimately parked for the operator.
func ExecuteSupervisedRun(ctx context.Context, options SupervisedAssemblyOptions) (RunSummary, error) {
	if err := ctx.Err(); err != nil {
		return RunSummary{}, err
	}
	clock, _ := supervisedClockAndIDs(options)
	taskID, stepID, policy, err := supervisedAssemblyStep(options.Definition)
	if err != nil {
		return RunSummary{}, err
	}
	if strings.TrimSpace(options.RunRoot) == "" {
		return RunSummary{}, fmt.Errorf("%w: run root is required", ErrInvalidConfig)
	}
	if options.VisibleSender == nil {
		return RunSummary{}, fmt.Errorf("%w: visible sender is required", ErrInvalidConfig)
	}
	eventLog := runEventLogPath(options.RunRoot)
	if err := ensureAgentRunnerEventLogMarker(eventLog); err != nil {
		return RunSummary{}, err
	}
	engine, runtime, err := newSupervisedEngine(ctx, options, taskID, stepID, *policy, clock)
	if err != nil {
		return RunSummary{}, err
	}
	projection := engine.Projection()
	if projection.ChainState == "PAUSED" {
		return runtime.resumeSupervisedRun(ctx)
	}
	runErr := engine.Run(ctx)
	if runErr != nil && !errors.Is(runErr, errSupervisedParked) {
		return summarizeProjection(engine.Projection(), eventLog), runErr
	}
	return summarizeProjection(engine.Projection(), eventLog), nil
}

// ExecuteSupervisedReturn consumes the operator's explicit return. It performs
// the post-return scan set, freezes the receipt and the reduced/unknown report,
// submits the return and drives the tail chain (postflight -> review ->
// publish). Only the engine's SubmitSupervisedReturn can pass the step.
func ExecuteSupervisedReturn(ctx context.Context, options SupervisedAssemblyOptions, outcome string) (RunSummary, error) {
	if err := ctx.Err(); err != nil {
		return RunSummary{}, err
	}
	if outcome != "complete" && outcome != "abort" {
		return RunSummary{}, ErrSupervisedAbortOnly
	}
	clock, _ := supervisedClockAndIDs(options)
	taskID, stepID, policy, err := supervisedAssemblyStep(options.Definition)
	if err != nil {
		return RunSummary{}, err
	}
	if strings.TrimSpace(options.RunRoot) == "" {
		return RunSummary{}, fmt.Errorf("%w: run root is required", ErrInvalidConfig)
	}
	eventLog := runEventLogPath(options.RunRoot)
	if err := ensureAgentRunnerEventLogMarker(eventLog); err != nil {
		return RunSummary{}, err
	}
	engine, runtime, err := newSupervisedEngine(ctx, options, taskID, stepID, *policy, clock)
	if err != nil {
		return RunSummary{}, err
	}
	session, found, err := runtime.loadSession()
	if err != nil {
		return RunSummary{}, err
	}
	if !found {
		return summarizeProjection(engine.Projection(), eventLog), ErrSupervisedNotOpen
	}
	if outcome == "abort" {
		return runtime.abortSupervisedRun(ctx, session, engine, eventLog)
	}
	// Return-entry guard (blind review High): the complete return enforces the
	// handoff policy -- timeout, authorized scope, after-return gates -- so it
	// must enforce the policy that was frozen when the window opened rather than
	// the configuration loaded at return time. Abort is deliberately exempt: it
	// consumes no policy and stays available as the operator's escape hatch.
	if err := runtime.enforceFrozenHandoffPolicy(session); err != nil {
		return summarizeProjection(engine.Projection(), eventLog), err
	}
	return runtime.completeSupervisedRun(ctx, session, engine, eventLog)
}

// supervisedAssemblyStep requires exactly one supervised-black-box code step
// with a handoff policy. Any other definition is refused before the engine
// exists, which is the assembly-level half of P1.
func supervisedAssemblyStep(definition chain.Definition) (string, string, *chain.HandoffPolicy, error) {
	if len(definition.Tasks) != 1 {
		return "", "", nil, fmt.Errorf("%w: the supervised assembly requires exactly one task", ErrInvalidConfig)
	}
	task := definition.Tasks[0]
	if len(task.Steps) != 1 {
		return "", "", nil, fmt.Errorf("%w: the supervised assembly requires exactly one step", ErrInvalidConfig)
	}
	step := task.Steps[0]
	if step.Kind != "code" || step.Mode != chain.SupervisedBlackBox {
		return "", "", nil, fmt.Errorf("%w: the supervised assembly requires one supervised-black-box code step", ErrInvalidConfig)
	}
	if step.HandoffPolicy == nil {
		return "", "", nil, fmt.Errorf("%w: the supervised step requires a handoff policy", ErrInvalidConfig)
	}
	return task.ID, step.ID, step.HandoffPolicy, nil
}

func supervisedClockAndIDs(options SupervisedAssemblyOptions) (func() time.Time, func() string) {
	clock := options.Clock
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	ids := options.IDs
	if ids == nil {
		ids = func() string { return fmt.Sprintf("supervised-%d", time.Now().UTC().UnixNano()) }
	}
	return clock, ids
}

// supervisedRuntime implements every chain port the supervised assembly needs.
// It is a real implementation over the filesystem and the snapshot store. It
// starts no process of its own.
type supervisedRuntime struct {
	runID   string
	runRoot string
	taskID  string
	stepID  string
	policy  chain.HandoffPolicy

	workspaceRoot string
	snapshotStore *snapshot.Store
	visible       adapters.VisibleClient
	gateRunner    *gates.Runner
	gateHooks     map[string]gates.Hook

	steps  chain.StepRouter
	engine *chain.Engine

	acceptedEntries []snapshot.Entry
	acceptedHash    string
	workspaceReady  bool
	events          int

	crashPoint func(point string) error
	clock      func() time.Time
	ids        func() string
}

var (
	_ chain.BaselinePort             = (*supervisedRuntime)(nil)
	_ chain.WorkspacePort            = (*supervisedRuntime)(nil)
	_ chain.StepPort                 = (*supervisedRuntime)(nil)
	_ chain.AcceptancePort           = (*supervisedRuntime)(nil)
	_ chain.ReviewerPort             = (*supervisedRuntime)(nil)
	_ chain.PublisherPort            = (*supervisedRuntime)(nil)
	_ chain.PostflightPort           = (*supervisedRuntime)(nil)
	_ chain.SupervisedPostflightPort = (*supervisedRuntime)(nil)
	_ chain.Stopper                  = (*supervisedRuntime)(nil)
	_ chain.Reconciler               = (*supervisedRuntime)(nil)
)

// supervisedStepPort is the StepRouter seam for the mode. In the parked flow the
// step is never dispatched (the failpoint stops the engine first); the port
// exists so a definition that somehow reaches the router fails closed with a
// precise reason instead of an untyped nil-port error. It performs the same
// delivery routine as the failpoint and writes no transition at all.
type supervisedStepPort struct {
	owner *supervisedRuntime
}

var _ chain.SupervisedBlackBoxPort = (*supervisedStepPort)(nil)

func (port *supervisedStepPort) OpenSupervisedBlackBox(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	if _, err := port.owner.deliverSupervisedNotice(ctx); err != nil {
		return chain.StepResult{}, err
	}
	return chain.StepResult{}, chain.ErrAwaitingSupervisedConfirmation
}

// newSupervisedEngine builds the engine, the runtime ports and the failpoint
// closure that parks the run.
func newSupervisedEngine(ctx context.Context, options SupervisedAssemblyOptions, taskID, stepID string, policy chain.HandoffPolicy, clock func() time.Time) (*chain.Engine, *supervisedRuntime, error) {
	snapshotRoot := filepath.Join(options.RunRoot, supervisedSnapshotStoreDirectory)
	snapshotStore, err := snapshot.NewStore(snapshotRoot, 0)
	if err != nil {
		return nil, nil, err
	}
	runtime := &supervisedRuntime{
		runID:         options.RunID,
		runRoot:       options.RunRoot,
		taskID:        taskID,
		stepID:        stepID,
		policy:        policy,
		workspaceRoot: filepath.Join(options.RunRoot, supervisedWorkspacesDirectory, fmt.Sprintf("%s-attempt-1", taskID)),
		snapshotStore: snapshotStore,
		visible:       options.VisibleSender,
		gateRunner:    options.GateRunner,
		gateHooks:     options.GateHooks,
		crashPoint:    options.CrashPoint,
		clock:         clock,
		ids:           options.IDs,
	}
	runtime.steps = chain.StepRouter{Supervised: &supervisedStepPort{owner: runtime}}

	eventStore, err := chain.NewFileEventStore(runEventLogPath(options.RunRoot))
	if err != nil {
		return nil, nil, err
	}
	engine, err := chain.New(ctx, chain.Options{
		RunID:      options.RunID,
		Definition: options.Definition,
		Events:     eventStore,
		Baselines:  runtime,
		Workspaces: runtime,
		Steps:      runtime,
		Acceptance: runtime,
		Reviewer:   runtime,
		Publisher:  runtime,
		// Postflight is required for every run by chain.New. A supervised-only run
		// never records AgentRunner facts, so the chain never consults this port;
		// it is wired fail-closed so a definition that did reach it would refuse
		// instead of skipping the acceptance gate.
		Postflight:           runtime,
		SupervisedPostflight: runtime,
		Stopper:              runtime,
		Reconciler:           runtime,
		Clock:                clock,
		IDs:                  runtime.nextEventID,
		Failpoint: func(point string) error {
			if point != supervisedStepFailpoint {
				return nil
			}
			return runtime.parkSupervisedStep(ctx)
		},
	})
	if err != nil {
		return nil, nil, err
	}
	runtime.engine = engine
	return engine, runtime, nil
}

// parkSupervisedStep is the front half of the run: it captures the before
// manifest, delivers the versioned notice and opens the binding confirmation
// interaction, then parks with errSupervisedParked. It writes no work window.
func (runtime *supervisedRuntime) parkSupervisedStep(ctx context.Context) error {
	manifestHash, err := runtime.deliverSupervisedNotice(ctx)
	if err != nil {
		return err
	}
	if err := runtime.openConfirmationInteraction(ctx, manifestHash); err != nil {
		return err
	}
	// The interaction is durable and the run is paused: park the engine so the
	// step is never dispatched and can never pass without an explicit return.
	return errSupervisedParked
}

// deliverSupervisedNotice captures the before manifest, persists it for the
// post-return diff, delivers the versioned notice through the visible channel
// and persists the delivery fact. A non-accepted delivery is refused here, so
// the window can never open on an undelivered notice.
//
// The delivery record is returned nowhere on purpose: the authoritative source
// of the delivery fact is the persisted delivery.json, which
// confirmAndOpenWindow re-reads when the window opens. Returning an in-memory
// copy would invite a second, divergent reader of the same fact.
func (runtime *supervisedRuntime) deliverSupervisedNotice(ctx context.Context) (string, error) {
	manifest, err := runtime.captureManifest(ctx, "baseline", "supervised-before", "")
	if err != nil {
		return "", err
	}
	if err := writeSupervisedJSON(runtime.beforeManifestPath(), manifest); err != nil {
		return "", err
	}
	now := runtime.clock().UTC()
	envelope, err := adapters.NewRequestEnvelope(adapters.RequestMessage{
		Type:        "request",
		RequestID:   runtime.nextID(),
		RunID:       runtime.runID,
		TaskID:      runtime.taskID,
		Attempt:     1,
		CreatedAt:   now.Format(timestampLayout),
		Adapter:     supervisedAdapterID,
		ContextHash: digest(supervisedContextDomain, runtime.runID, runtime.taskID, runtime.stepID, manifest.ManifestHash),
	})
	if err != nil {
		return "", err
	}
	delivery, err := adapters.DeliverVisible(ctx, runtime.visible, envelope)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrSupervisedDeliveryBlocked, err)
	}
	record := chain.SupervisedDelivery{
		Channel:            supervisedRequiredChannel,
		DeliveredAt:        now.Format(timestampLayout),
		DeliveryStatus:     supervisedDeliveryAcceptedStatus,
		TransportRequestID: delivery.TransportRequestID,
		ErrorEvidence:      delivery.ErrorEvidence,
	}
	if err := writeSupervisedJSON(runtime.deliveryPath(), record); err != nil {
		return "", err
	}
	return manifest.ManifestHash, nil
}

// openConfirmationInteraction persists (or reuses) the binding operator
// interaction request and pauses the run. The request binds the before-manifest
// digest, the notice version and hash, and the structured confirm/reject answer
// set.
func (runtime *supervisedRuntime) openConfirmationInteraction(ctx context.Context, manifestHash string) error {
	ledger, err := runtime.loadLedger()
	if err != nil {
		return err
	}
	request, found := supervisedPendingRequest(ledger, runtime.runID, runtime.taskID, runtime.stepID)
	if !found {
		now := runtime.clock().UTC()
		requestRecord, err := chain.NewOperatorInteractionRequestRecord(chain.OperatorInteractionRequest{
			Kind:             "request",
			InteractionID:    runtime.nextID(),
			RequestedAt:      now.Format(timestampLayout),
			RequestedBy:      evidence.Actor{Type: "adapter", ID: supervisedAdapterID},
			Operator:         evidence.Actor{Type: "operator", ID: supervisedOperatorID},
			RunID:            runtime.runID,
			TaskID:           runtime.taskID,
			StepID:           runtime.stepID,
			Attempt:          1,
			WorkspaceHash:    manifestHash,
			ConversationID:   runtime.nextID(),
			RequestID:        runtime.nextID(),
			ContextHash:      digest(supervisedContextDomain, runtime.runID, runtime.taskID, runtime.stepID, manifestHash),
			Question:         chain.SupervisedRiskNoticeText,
			Reason:           "supervised-black-box risk notice version " + fmt.Sprintf("%d", chain.SupervisedRiskNoticeVersion),
			Risk:             supervisedRiskDescription,
			RequiredAction:   supervisedRequiredActionText,
			AllowedResponses: []string{supervisedConfirmationResponse, supervisedRejectionResponse},
			Evidence:         []string{chain.SupervisedRiskNoticeHash()},
		}, runtime.clock, runtime.idSource())
		if err != nil {
			return err
		}
		if err := appendOperatorInteractionRecord(runtime.ledgerPath(), requestRecord); err != nil {
			return err
		}
		request = requestRecord
	}
	return runtime.engine.OpenOperatorInteraction(ctx, request)
}

// resumeSupervisedRun handles phases B, C and D from the persisted state.
func (runtime *supervisedRuntime) resumeSupervisedRun(ctx context.Context) (RunSummary, error) {
	eventLog := runEventLogPath(runtime.runRoot)
	ledger, err := runtime.loadLedger()
	if err != nil {
		return RunSummary{}, err
	}
	request, found := supervisedPendingRequest(ledger, runtime.runID, runtime.taskID, runtime.stepID)
	if !found {
		return summarizeProjection(runtime.engine.Projection(), eventLog), fmt.Errorf("%w: a paused supervised run has no binding interaction request", ErrInvalidConfig)
	}
	response, found := supervisedInteractionResponse(ledger, request.Request.InteractionID)
	if !found {
		// Phase B: the operator has not answered yet.
		return summarizeProjection(runtime.engine.Projection(), eventLog), nil
	}
	confirmations, err := runtime.loadConfirmationLedger()
	if err != nil {
		return RunSummary{}, err
	}
	if len(confirmations.Records) == 0 {
		// Phase C: validate the answer, freeze the confirmation and open the window.
		if err := runtime.confirmAndOpenWindow(ctx, request, response); err != nil {
			return summarizeProjection(runtime.engine.Projection(), eventLog), err
		}
		return summarizeProjection(runtime.engine.Projection(), eventLog), nil
	}
	if _, found, err := runtime.loadSession(); err != nil {
		return RunSummary{}, err
	} else if found {
		// Phase D: the window is open and waits for the explicit return.
		return summarizeProjection(runtime.engine.Projection(), eventLog), nil
	}
	// The confirmation is durable but the session is not: the process died
	// inside the confirm/resume/open/persist sequence (plan section 9.7).
	// Replaying the durable confirmation converges the run instead of stalling
	// on a non-empty ledger.
	if err := runtime.replayConfirmedWindow(ctx, request, response, confirmations); err != nil {
		return summarizeProjection(runtime.engine.Projection(), eventLog), err
	}
	return summarizeProjection(runtime.engine.Projection(), eventLog), nil
}

// replayConfirmedWindow re-drives the confirmation sequence after a crash that
// left the confirmation durable but the handoff session unwritten. It replays
// the existing confirmation: rebuilding it is impossible and would be wrong,
// because the record ids are non-deterministic and a second, divergent record
// for the same step attempt is refused by the ledger's conflict rule.
func (runtime *supervisedRuntime) replayConfirmedWindow(ctx context.Context, request, response chain.OperatorInteractionRecord, ledger chain.SupervisedRiskConfirmationLedger) error {
	confirmation, found := supervisedConfirmationForBinding(ledger, runtime.runID, runtime.taskID, runtime.stepID)
	if !found {
		return fmt.Errorf("%w: the durable confirmation ledger holds no record for this step attempt", ErrInvalidConfig)
	}
	// The replay opens the window again, so it freezes the policy the same way a
	// first attempt does: when the crash was after the window opened the snapshot
	// already exists and this call re-validates the policy in force against it
	// (failing closed on drift), otherwise it writes the freeze the original
	// attempt did not reach.
	if err := runtime.persistFrozenHandoffPolicy(); err != nil {
		return err
	}
	session := supervisedSessionFromRequest(request, runtime.policy)
	// The original open time was never persisted, so the replay stamps the
	// recovery time; the already-open branch below recovers the recorded time.
	session.OpenedAt = runtime.clock().UTC().Format(timestampLayout)
	if _, err := runtime.engine.ResumeOperatorInteraction(ctx, request, response, supervisedInteractionBinding(request)); err != nil {
		// The one legitimate refusal is the window already being open: the
		// waiting transitions then bind the confirmation instead of the
		// interaction, so the resume has nothing to lift. Detect that case by
		// its durable marker and converge; anything else is a real dead-end.
		openedAt, open, openErr := runtime.supervisedOpenWindowAt(ctx, confirmation.RecordHash)
		if openErr != nil {
			return openErr
		}
		if !open {
			return fmt.Errorf("%w: replaying the operator resume failed: %v", ErrSupervisedWindowUnproven, err)
		}
		session.OpenedAt = openedAt
		return writeSupervisedJSON(runtime.sessionPath(), session)
	}
	if err := runtime.engine.OpenSupervisedWorkWindow(ctx, session, confirmation); err != nil {
		return err
	}
	return writeSupervisedJSON(runtime.sessionPath(), session)
}

// supervisedOpenWindowAt reports whether the work window is already open for a
// confirmation digest, returning the recorded open time. The window opener is
// the only writer that binds the confirmation digest into a WAITING_FOR_OPERATOR
// transition, so that binding is the durable marker of an open window.
func (runtime *supervisedRuntime) supervisedOpenWindowAt(ctx context.Context, confirmationHash string) (string, bool, error) {
	store, err := chain.NewFileEventStore(runEventLogPath(runtime.runRoot))
	if err != nil {
		return "", false, err
	}
	events, err := store.Load(ctx)
	if err != nil {
		return "", false, err
	}
	openedAt := ""
	open := false
	for _, record := range events {
		event := record.Event
		if event.ToState != "WAITING_FOR_OPERATOR" || !slices.Contains(event.InputEvidence, confirmationHash) {
			continue
		}
		if event.Entity.Kind != "step" || event.Entity.RunID != runtime.runID ||
			event.Entity.TaskID != runtime.taskID || event.Entity.StepID != runtime.stepID {
			continue
		}
		openedAt, open = event.OccurredAt, true
	}
	return openedAt, open, nil
}

// supervisedConfirmationForBinding returns the durable confirmation of a step
// attempt. The ledger holds at most one per attempt.
func supervisedConfirmationForBinding(ledger chain.SupervisedRiskConfirmationLedger, runID, taskID, stepID string) (chain.SupervisedRiskConfirmationRecord, bool) {
	for _, record := range ledger.Records {
		if record.Record.RunID == runID && record.Record.TaskID == taskID && record.Record.StepID == stepID && record.Record.Attempt == 1 {
			return record, true
		}
	}
	return chain.SupervisedRiskConfirmationRecord{}, false
}

// confirmAndOpenWindow builds the binding risk confirmation, resumes the
// interaction and opens the supervised work window.
func (runtime *supervisedRuntime) confirmAndOpenWindow(ctx context.Context, request, response chain.OperatorInteractionRecord) error {
	if response.Response.Selection != supervisedConfirmationResponse {
		// P2/P3: a rejection is never a confirmation and never opens a window.
		return fmt.Errorf("%w: operator selection %q", ErrSupervisedRejected, response.Response.Selection)
	}
	delivery, found, err := runtime.loadDelivery()
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("%w: the delivery fact is missing", ErrSupervisedDeliveryBlocked)
	}
	// Contracts section 7.1a: only an accepted delivery may unlock the window.
	// The capture path refuses a non-accepted transport answer, but a durable
	// record that contradicts it (crafted or written by an older build) must
	// fail closed here too, before the confirmation is frozen or the operator
	// resume is written. The chain re-refuses at the window opener, so this is
	// the console half of the same invariant.
	if delivery.DeliveryStatus != supervisedDeliveryAcceptedStatus {
		return fmt.Errorf("%w: delivery status %q is not accepted", ErrSupervisedDeliveryBlocked, delivery.DeliveryStatus)
	}
	now := runtime.clock().UTC()
	confirmation, err := chain.NewSupervisedRiskConfirmation(chain.SupervisedRiskConfirmationInput{
		RecordID:            runtime.nextID(),
		IssuedAt:            now.Format(timestampLayout),
		OperatorID:          response.Response.RespondedBy.ID,
		RunID:               runtime.runID,
		TaskID:              runtime.taskID,
		StepID:              runtime.stepID,
		Attempt:             1,
		WorkspaceHash:       request.Request.WorkspaceHash,
		RiskNoticeVersion:   chain.SupervisedRiskNoticeVersion,
		RiskNoticeHash:      chain.SupervisedRiskNoticeHash(),
		Delivery:            delivery,
		InteractionEvidence: []string{request.RecordHash, response.RecordHash},
		Evidence:            []string{digest(supervisedConfirmationEvidenceDomain, runtime.runID, runtime.taskID, runtime.stepID, request.RecordHash, response.RecordHash)},
	}, runtime.clock, runtime.idSource())
	if err != nil {
		return err
	}
	if err := appendSupervisedConfirmation(runtime.confirmationsPath(), confirmation); err != nil {
		return err
	}
	if err := runtime.crash(supervisedCrashAfterConfirmation); err != nil {
		return err
	}
	if _, err := runtime.engine.ResumeOperatorInteraction(ctx, request, response, supervisedInteractionBinding(request)); err != nil {
		return err
	}
	if err := runtime.crash(supervisedCrashAfterResume); err != nil {
		return err
	}
	session := supervisedSessionFromRequest(request, runtime.policy)
	session.OpenedAt = now.Format(timestampLayout)
	// Freeze the handoff policy of this window before the window becomes durable
	// (blind review High). The snapshot is written first, so an open window always
	// has its freeze artifact and the return entry can enforce the frozen policy
	// instead of the configuration loaded at return time.
	if err := runtime.persistFrozenHandoffPolicy(); err != nil {
		return err
	}
	if err := runtime.engine.OpenSupervisedWorkWindow(ctx, session, confirmation); err != nil {
		return err
	}
	if err := runtime.crash(supervisedCrashAfterWindowOpen); err != nil {
		return err
	}
	return writeSupervisedJSON(runtime.sessionPath(), session)
}

// supervisedInteractionBinding is the operator-interaction binding carried by a
// request record.
func supervisedInteractionBinding(request chain.OperatorInteractionRecord) chain.OperatorInteractionBinding {
	return chain.OperatorInteractionBinding{
		RunID:          request.Request.RunID,
		TaskID:         request.Request.TaskID,
		StepID:         request.Request.StepID,
		Attempt:        request.Request.Attempt,
		WorkspaceHash:  request.Request.WorkspaceHash,
		ConversationID: request.Request.ConversationID,
		ContextHash:    request.Request.ContextHash,
	}
}

// supervisedSessionFromRequest rebuilds the handoff session of a binding
// confirmation request. It is used both when the window opens for the first time
// and when a crash replay recovers an unwritten session; OpenedAt is set by the
// caller.
func supervisedSessionFromRequest(request chain.OperatorInteractionRecord, policy chain.HandoffPolicy) chain.HandoffSession {
	return chain.HandoffSession{
		RunID:                request.Request.RunID,
		TaskID:               request.Request.TaskID,
		StepID:               request.Request.StepID,
		Attempt:              request.Request.Attempt,
		HandoffPolicyHash:    supervisedHandoffPolicyHash(policy),
		AllowedReturnActions: append([]string{}, policy.ReturnActions...),
		Operator:             request.Request.Operator,
		BeforeManifestHash:   request.Request.WorkspaceHash,
		LeaseEvidence:        []string{digest(supervisedSessionLeaseDomain, request.Request.RunID, request.Request.TaskID)},
	}
}

// crash invokes the optional crash-injection seam of the confirmation sequence.
func (runtime *supervisedRuntime) crash(point string) error {
	if runtime.crashPoint == nil {
		return nil
	}
	return runtime.crashPoint(point)
}

// completeSupervisedRun performs the post-return scan set, freezes the five
// facts, submits the return and drives the tail chain.
func (runtime *supervisedRuntime) completeSupervisedRun(ctx context.Context, session chain.HandoffSession, engine *chain.Engine, eventLog string) (RunSummary, error) {
	now := runtime.clock().UTC()
	openedAt, err := time.Parse(timestampLayout, session.OpenedAt)
	if err != nil {
		return RunSummary{}, err
	}
	if now.Sub(openedAt) > time.Duration(runtime.policy.HandoffTimeoutMs)*time.Millisecond {
		// The window outlived its frozen timeout: complete is refused, abort stays
		// possible.
		return summarizeProjection(engine.Projection(), eventLog), fmt.Errorf("%w: window exceeded %d ms", ErrSupervisedTimeout, runtime.policy.HandoffTimeoutMs)
	}
	afterManifest, err := runtime.captureManifest(ctx, "candidate", "supervised-after", session.BeforeManifestHash)
	if err != nil {
		return RunSummary{}, err
	}
	beforeManifest, found, err := runtime.loadBeforeManifest()
	if err != nil {
		return RunSummary{}, err
	}
	if !found {
		return RunSummary{}, fmt.Errorf("%w: the before manifest is missing", ErrInvalidConfig)
	}
	diffHash, changed := supervisedDiff(beforeManifest, afterManifest)
	// GAP-1: the return-side scan set covers the whole AT-24 wording. The scope
	// scan compares every changed path against the authorized scope, and the
	// side-effect scan compares the two manifests for observable side effects.
	// Both fail closed and refuse the return before any receipt is built: an
	// underivable scope, an out-of-scope change and a detected side effect are
	// all blocking (plan section 8).
	if err := runtime.runScopeScan(supervisedChangedPaths(changed)); err != nil {
		return summarizeProjection(engine.Projection(), eventLog), err
	}
	if err := runtime.runSideEffectScan(beforeManifest.Manifest.Entries, afterManifest.Manifest.Entries); err != nil {
		return summarizeProjection(engine.Projection(), eventLog), err
	}
	if err := runtime.runSecretAndTypeScan(ctx, afterManifest); err != nil {
		return summarizeProjection(engine.Projection(), eventLog), err
	}
	processStopHash, err := runtime.processStopEvidence()
	if err != nil {
		return summarizeProjection(engine.Projection(), eventLog), err
	}
	gateHashes, gateEvidenceHash, err := runtime.runAfterReturnGates(ctx)
	if err != nil {
		return summarizeProjection(engine.Projection(), eventLog), err
	}
	// Every flag below is reachable only because the corresponding scan
	// returned no blocking error above: a scope violation, a side effect, a
	// secret finding or an unaccepted file type refuses the return before this
	// report is built. The report is the scan evidence that feeds the frozen
	// ScanEvidenceHash fact, so the new scope/side-effect scans reach the frozen
	// five facts through the same wiring as the pre-existing scans.
	scanReport := supervisedScanReport{
		RunID:              runtime.runID,
		TaskID:             runtime.taskID,
		StepID:             runtime.stepID,
		Attempt:            1,
		BeforeManifestHash: session.BeforeManifestHash,
		AfterManifestHash:  afterManifest.ManifestHash,
		DiffHash:           diffHash,
		ChangedPaths:       changed,
		EntryCount:         len(afterManifest.Manifest.Entries),
		ChangeSetInScope:   true,
		SecretScanClean:    true,
		FileTypeScanClean:  true,
		SideEffectClean:    true,
		ProcessStopHash:    processStopHash,
		GateHashes:         gateHashes,
	}
	if err := writeSupervisedJSON(runtime.scanReportPath(), scanReport); err != nil {
		return RunSummary{}, err
	}
	scanEvidenceHash, err := supervisedScanEvidenceHash(runtime.scanReportPath())
	if err != nil {
		return RunSummary{}, err
	}
	receipt, err := chain.BuildHandoffReceiptRecord(session, chain.HandoffDecision{
		ReceiptID:          runtime.nextID(),
		ClosedAt:           now.Format(timestampLayout),
		AfterManifestHash:  afterManifest.ManifestHash,
		DiffHash:           diffHash,
		Outcome:            "complete",
		HookResultEvidence: gateHashes,
		Evidence:           []string{scanEvidenceHash, processStopHash},
	}, runtime.clock, runtime.idSource())
	if err != nil {
		return RunSummary{}, err
	}
	report, err := chain.NewSupervisedReport(chain.SupervisedReportInput{
		ReportID:           runtime.nextID(),
		RecordedAt:         now.Format(timestampLayout),
		RunID:              runtime.runID,
		TaskID:             runtime.taskID,
		StepID:             runtime.stepID,
		Attempt:            1,
		VerifiedDimensions: append([]string{}, supervisedVerifiedDimensions...),
		UnknownDimensions:  append([]string{}, supervisedUnknownDimensions...),
		Evidence:           []string{receipt.ReceiptHash, digest(supervisedReportEvidenceDomain, runtime.runID, runtime.taskID, runtime.stepID)},
	}, runtime.clock, runtime.idSource())
	if err != nil {
		return RunSummary{}, err
	}
	if err := writeSupervisedJSON(runtime.receiptPath(), receipt); err != nil {
		return RunSummary{}, err
	}
	if err := writeSupervisedJSON(runtime.reportPath(), report); err != nil {
		return RunSummary{}, err
	}
	facts := chain.SupervisedFrozenFacts{
		HandoffReceiptHash:      receipt.ReceiptHash,
		ScanEvidenceHash:        scanEvidenceHash,
		ProcessStopEvidenceHash: processStopHash,
		GateEvidenceHash:        gateEvidenceHash,
		ReportHash:              report.RecordHash,
	}
	if err := engine.SubmitSupervisedReturn(ctx, receipt, facts); err != nil {
		return summarizeProjection(engine.Projection(), eventLog), err
	}
	if err := engine.Run(ctx); err != nil {
		return summarizeProjection(engine.Projection(), eventLog), err
	}
	return summarizeProjection(engine.Projection(), eventLog), nil
}

// abortSupervisedRun persists the abort receipt and cancels the run through the
// existing engine cancellation path. Deviation from the plan: the plan mapped
// abort to a FAILED task/step, but no exported engine entry writes that
// transition without the console writing chain state directly, so the PILOT
// abort uses the engine's existing RequestCancel path (terminal CANCELLED).
func (runtime *supervisedRuntime) abortSupervisedRun(ctx context.Context, session chain.HandoffSession, engine *chain.Engine, eventLog string) (RunSummary, error) {
	now := runtime.clock().UTC()
	receipt, err := chain.BuildHandoffReceiptRecord(session, chain.HandoffDecision{
		ReceiptID:         runtime.nextID(),
		ClosedAt:          now.Format(timestampLayout),
		AfterManifestHash: session.BeforeManifestHash,
		DiffHash:          digest(supervisedDiffDomain, "aborted", runtime.runID, runtime.taskID, runtime.stepID),
		Outcome:           "abort",
		Evidence:          []string{digest(supervisedReportEvidenceDomain, "abort", runtime.runID)},
	}, runtime.clock, runtime.idSource())
	if err != nil {
		return RunSummary{}, err
	}
	if err := writeSupervisedJSON(runtime.receiptPath(), receipt); err != nil {
		return RunSummary{}, err
	}
	engine.RequestCancel()
	if err := engine.Run(ctx); err != nil && !errors.Is(err, chain.ErrCancelled) {
		return summarizeProjection(engine.Projection(), eventLog), err
	}
	return summarizeProjection(engine.Projection(), eventLog), nil
}

// captureManifest captures the workspace manifest and verifies it. A non-empty
// parent hash makes it a candidate capture (the snapshot kinds differ: a
// candidate requires a parent digest and a task binding).
func (runtime *supervisedRuntime) captureManifest(ctx context.Context, kind, snapshotIDPrefix, parentHash string) (snapshot.SnapshotManifest, error) {
	if err := runtime.ensureWorkspace(); err != nil {
		return snapshot.SnapshotManifest{}, err
	}
	options := snapshot.CaptureOptions{
		SourceDir:  runtime.workspaceRoot,
		Kind:       kind,
		SnapshotID: fmt.Sprintf("%s-%s", snapshotIDPrefix, runtime.runID),
		RunID:      runtime.runID,
		Store:      runtime.snapshotStore,
	}
	if parentHash != "" {
		parent := parentHash
		options.ParentSnapshotHash = &parent
		options.Task = &snapshot.TaskBinding{TaskID: runtime.taskID, Attempt: 1}
	}
	manifest, err := snapshot.Capture(ctx, options)
	if err != nil {
		return snapshot.SnapshotManifest{}, err
	}
	return manifest, nil
}

// supervisedDiff lists the entries that differ between the persisted before
// manifest and the post-return capture and freezes a diff digest. The caller
// loads the before manifest once; the returned mismatches are the change set the
// scope scan and the side-effect scan consume.
func supervisedDiff(before, after snapshot.SnapshotManifest) (string, []string) {
	changed := supervisedEntryMismatches(before.Manifest.Entries, after.Manifest.Entries)
	return digest(supervisedDiffDomain, append([]string{before.ManifestHash, after.ManifestHash}, changed...)...), changed
}

// supervisedChangedPaths extracts the raw path of every mismatch produced by
// supervisedEntryMismatches ("added x", "changed x", "removed x").
func supervisedChangedPaths(mismatches []string) []string {
	paths := make([]string, 0, len(mismatches))
	for _, mismatch := range mismatches {
		paths = append(paths, supervisedChangedPath(mismatch))
	}
	return paths
}

func supervisedChangedPath(mismatch string) string {
	for _, prefix := range []string{"added ", "changed ", "removed "} {
		if strings.HasPrefix(mismatch, prefix) {
			return strings.TrimPrefix(mismatch, prefix)
		}
	}
	return mismatch
}

// supervisedScope is the authorized write scope of a supervised return. It is
// derived from the frozen handoff policy, not invented here: HandoffPolicy
// .AllowedTargets is the repository's declared allow-list for a handoff (the
// read-write target IDs a handoff authorizes, CONTRACTS section 2). The console
// assembly holds no target registry, so the authorized targets can only be
// observed through the isolated workspace they materialize into; the scope is
// therefore the workspace root, and a changed path is in scope only when it is a
// contained relative path that resolves inside it (the same containment rule the
// gate runner applies to a hook cwd, gates.ErrWorkspaceEscape, t010-gates item 4).
type supervisedScope struct {
	root string
}

// newSupervisedScope derives the authorized scope. An empty allow-list is an
// underivable scope: it fails closed rather than being read as "no restriction".
func newSupervisedScope(policy chain.HandoffPolicy, workspaceRoot string) (supervisedScope, error) {
	if len(policy.AllowedTargets) == 0 || strings.TrimSpace(workspaceRoot) == "" {
		return supervisedScope{}, fmt.Errorf("%w: the frozen handoff policy declares no allowed target", ErrSupervisedScopeUnderivable)
	}
	absolute, err := filepath.Abs(workspaceRoot)
	if err != nil {
		return supervisedScope{}, fmt.Errorf("%w: %v", ErrSupervisedScopeUnderivable, err)
	}
	resolved, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return supervisedScope{}, fmt.Errorf("%w: the workspace root cannot be resolved: %v", ErrSupervisedScopeUnderivable, err)
	}
	return supervisedScope{root: resolved}, nil
}

// allows reports whether a changed workspace-relative path is inside the
// authorized scope: it must be relative, must not traverse upwards, and must not
// resolve outside the workspace root through a link.
func (scope supervisedScope) allows(path string) error {
	cleaned := filepath.Clean(filepath.FromSlash(path))
	if path == "" || filepath.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(filepath.Separator)) {
		return fmt.Errorf("%w: %q", ErrSupervisedScopeOutOfBounds, path)
	}
	resolved, err := filepath.EvalSymlinks(filepath.Join(scope.root, cleaned))
	if err != nil {
		// The path is gone (a deletion, for example) or is not materialized. Its
		// lexical position is already proven contained; deletions are judged by
		// the side-effect scan, not here.
		return nil
	}
	relative, err := filepath.Rel(scope.root, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("%w: %q resolves outside the workspace root", ErrSupervisedScopeOutOfBounds, path)
	}
	return nil
}

// runScopeScan compares the post-return changed paths against the authorized
// scope. An out-of-scope change blocks the return; an underivable scope blocks
// before any comparison happens instead of being treated as "no restriction".
// The scan judges the manifest change set only: an out-of-scope write anywhere
// outside the observed workspace is unobservable and is carried by the report's
// external-side-effects unknown dimension rather than being claimed absent.
func (runtime *supervisedRuntime) runScopeScan(changed []string) error {
	scope, err := newSupervisedScope(runtime.policy, runtime.workspaceRoot)
	if err != nil {
		return err
	}
	for _, path := range changed {
		if err := scope.allows(path); err != nil {
			return err
		}
	}
	return nil
}

// runSideEffectScan compares the frozen before manifest with the post-return
// manifest and blocks every observable side effect on the isolated workspace: a
// deletion of a pre-existing entry, a type change, a permission/executable
// change, or a symbolic link whose target leaves the workspace root.
//
// Dividing line (plan section 8): these are facts ProofRail can observe, so a
// detection blocks. Physical effects outside the observable surface -- network
// access, external commits/pushes/publications, cost, and the external agent's
// own process tree -- cannot be observed from ProofRail's position and are NOT
// claimed absent here; they are carried by the report's unknownDimensions
// (external-side-effects, network-access, cost, background-processes) rather
// than being claimed clean. Creating a new file inside the workspace is the
// declared handoff work, not a side effect.
func (runtime *supervisedRuntime) runSideEffectScan(before, after []snapshot.Entry) error {
	afterByPath := make(map[string]snapshot.Entry, len(after))
	for _, entry := range after {
		afterByPath[entry.Path] = entry
	}
	for _, prior := range before {
		current, exists := afterByPath[prior.Path]
		if !exists {
			return fmt.Errorf("%w: pre-existing entry %q was deleted", ErrSupervisedSideEffect, prior.Path)
		}
		if current.Type != prior.Type {
			return fmt.Errorf("%w: entry %q changed type from %q to %q", ErrSupervisedSideEffect, prior.Path, prior.Type, current.Type)
		}
		if current.ReadOnly != prior.ReadOnly || current.Executable != prior.Executable {
			return fmt.Errorf("%w: entry %q changed permissions", ErrSupervisedSideEffect, prior.Path)
		}
	}
	scope, err := newSupervisedScope(runtime.policy, runtime.workspaceRoot)
	if err != nil {
		return err
	}
	for _, entry := range after {
		if entry.Type != snapshot.EntrySymbolicLink {
			continue
		}
		if strings.HasPrefix(entry.Target, "/") || filepath.IsAbs(entry.Target) {
			return fmt.Errorf("%w: symbolic link %q has an absolute target", ErrSupervisedSideEffect, entry.Path)
		}
		joined := path.Clean(path.Join(path.Dir(entry.Path), entry.Target))
		if err := scope.allows(joined); err != nil {
			return fmt.Errorf("%w: symbolic link %q leaves the workspace root", ErrSupervisedSideEffect, entry.Path)
		}
	}
	return nil
}

// runSecretAndTypeScan refuses credential material and unaccepted entry types in
// the returning workspace.
func (runtime *supervisedRuntime) runSecretAndTypeScan(ctx context.Context, manifest snapshot.SnapshotManifest) error {
	scanner := gates.DefaultSecretScanner()
	for _, entry := range manifest.Manifest.Entries {
		switch entry.Type {
		case snapshot.EntryRegularFile, snapshot.EntryDirectory, snapshot.EntrySymbolicLink:
		default:
			return fmt.Errorf("%w: unexpected entry type %q at %q", ErrSupervisedFileType, entry.Type, entry.Path)
		}
		if entry.Type != snapshot.EntryRegularFile {
			continue
		}
		data, err := os.ReadFile(filepath.Join(runtime.workspaceRoot, filepath.FromSlash(entry.Path)))
		if err != nil {
			return fmt.Errorf("%w: read %q: %v", ErrSupervisedGateFailed, entry.Path, err)
		}
		if _, err := scanner.Scan(ctx, entry.Path, data); err != nil {
			return fmt.Errorf("%w: %q", ErrSupervisedSecretFinding, entry.Path)
		}
	}
	return nil
}

// processStopEvidence proves that no managed process of this run is alive. A run
// with no managed-process identity has the trivial proof: no managed process was
// started. A recorded identity that is still alive fails closed.
//
// GAP-4: "unknown process" cannot be positively observed. ProofRail can inspect
// only the process identity it recorded itself; the external agent's process
// tree is unobservable from ProofRail's position. It is therefore not scanned
// and not claimed absent -- it is carried by the report's "background-processes"
// unknown dimension (plan section 8). The detectable subset is a *managed*
// identity that is still alive, which blocks the return here.
func (runtime *supervisedRuntime) processStopEvidence() (string, error) {
	identityPath := filepath.Join(runtime.runRoot, managedProcessIdentityFileName)
	identity, err := loadManagedProcessIdentity(identityPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return digest(supervisedProcessStopDomain, runtime.runID, runtime.taskID, supervisedMissingProcessReason), nil
		}
		return "", fmt.Errorf("%w: %v", ErrSupervisedGateFailed, err)
	}
	live, err := guard.ProcessAlive(identity)
	if err != nil {
		return "", fmt.Errorf("%w: managed process liveness is unobservable: %v", ErrSupervisedGateFailed, err)
	}
	if live {
		return "", fmt.Errorf("%w: a managed process of this run is still alive", ErrSupervisedGateFailed)
	}
	return digest(supervisedProcessStopDomain, runtime.runID, runtime.taskID, identity.StartToken), nil
}

// runAfterReturnGates runs the frozen after-return gate hooks. A missing runner
// or hook definition fails closed, and a failed hook refuses the return.
func (runtime *supervisedRuntime) runAfterReturnGates(ctx context.Context) ([]string, string, error) {
	if runtime.gateRunner == nil {
		return nil, "", fmt.Errorf("%w: no gate runner was assembled", ErrSupervisedGateUnavailable)
	}
	hashes := make([]string, 0, len(runtime.policy.HooksAfterReturn))
	for _, hookID := range runtime.policy.HooksAfterReturn {
		hook, ok := runtime.gateHooks[hookID]
		if !ok {
			return nil, "", fmt.Errorf("%w: hook %q", ErrSupervisedGateUnavailable, hookID)
		}
		hook.ID = hookID
		definitionHash, err := gates.DefinitionHash(hook)
		if err != nil {
			return nil, "", err
		}
		record, err := runtime.gateRunner.Run(ctx, gates.Request{
			RunID:              runtime.runID,
			TaskID:             runtime.taskID,
			StepID:             hookID,
			Attempt:            1,
			ExecutionAttempt:   1,
			WorkspaceRoot:      runtime.workspaceRoot,
			HookDefinitionHash: definitionHash,
			Hook:               hook,
		})
		if err != nil {
			return nil, "", fmt.Errorf("%w: hook %q: %v", ErrSupervisedGateFailed, hookID, err)
		}
		if !evidence.ValidHash(record.ResultHash) {
			return nil, "", fmt.Errorf("%w: hook %q returned no result hash", ErrSupervisedGateFailed, hookID)
		}
		hashes = append(hashes, record.ResultHash)
	}
	sort.Strings(hashes)
	if err := writeSupervisedJSON(runtime.gateEvidencePath(), supervisedGateEvidence{RunID: runtime.runID, Hashes: hashes}); err != nil {
		return nil, "", err
	}
	return hashes, digest(supervisedGateEvidenceDomain, append([]string{runtime.runID}, hashes...)...), nil
}

// --- chain port implementations -------------------------------------------------

// CaptureBaseline digests the workspace before the task starts.
func (runtime *supervisedRuntime) CaptureBaseline(ctx context.Context, _ string) (chain.SnapshotRef, error) {
	manifest, err := runtime.captureManifest(ctx, "baseline", "supervised-baseline", "")
	if err != nil {
		return chain.SnapshotRef{}, err
	}
	return chain.SnapshotRef{Hash: manifest.ManifestHash}, nil
}

// Materialize returns the one isolated workspace of this run root.
func (runtime *supervisedRuntime) Materialize(context.Context, string, string, int, chain.SnapshotRef) (chain.Workspace, error) {
	if err := runtime.ensureWorkspace(); err != nil {
		return chain.Workspace{}, err
	}
	return chain.Workspace{Root: runtime.workspaceRoot}, nil
}

// Execute routes the one step through the supervised port.
func (runtime *supervisedRuntime) Execute(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	return runtime.steps.Execute(ctx, request)
}

// Accept freezes the candidate workspace after the postflight gate qualified it.
func (runtime *supervisedRuntime) Accept(ctx context.Context, runID, taskID string, attempt int, parent chain.SnapshotRef, workspace chain.Workspace) (chain.CandidateResult, error) {
	parentHash := parent.Hash
	captured, err := snapshot.Capture(ctx, snapshot.CaptureOptions{
		SourceDir:          workspace.Root,
		Kind:               "candidate",
		SnapshotID:         fmt.Sprintf("supervised-candidate-%s-%s-%d", runID, taskID, attempt),
		RunID:              runID,
		ParentSnapshotHash: &parentHash,
		Task:               &snapshot.TaskBinding{TaskID: taskID, Attempt: attempt},
		Store:              runtime.snapshotStore,
	})
	if err != nil {
		return chain.CandidateResult{}, err
	}
	runtime.acceptedEntries = captured.Manifest.Entries
	runtime.acceptedHash = captured.ManifestHash
	evidenceRoot := digest(supervisedCandidateEvidenceRootDomain, runID, taskID, fmt.Sprintf("%d", attempt), parentHash, captured.ManifestHash)
	return chain.CandidateResult{
		Snapshot:          chain.SnapshotRef{Hash: captured.ManifestHash},
		EvidenceRootHash:  evidenceRoot,
		Evidence:          []string{digest(supervisedCandidateEvidenceDomain, captured.ManifestHash)},
		CandidateProducer: evidence.Actor{Type: "agent", ID: "supervised-operator"},
	}, nil
}

// Review is a real policy review over the frozen candidate. The chain only
// reaches it from REVIEW_PENDING, so the supervised postflight qualification is
// already proven; there is no fast path.
func (runtime *supervisedRuntime) Review(_ context.Context, request chain.ReviewRequest) (chain.ReviewDecision, error) {
	if !evidence.ValidHash(request.CandidateSnapshotHash) || !evidence.ValidHash(request.EvidenceRootHash) {
		return chain.ReviewDecision{}, fmt.Errorf("%w: review request carries an invalid binding", ErrInvalidConfig)
	}
	return chain.ReviewDecision{
		ReceiptID:  runtime.nextID(),
		OccurredAt: runtime.clock().UTC().Format(timestampLayout),
		RecordedBy: evidence.Actor{Type: "policy", ID: supervisedReviewPolicyLabel},
		ReviewMode: "policy",
		PolicyHash: digest(supervisedReviewPolicyDomain, supervisedReviewPolicyLabel, "supervised-postflight-qualified-candidate"),
		Outcome:    "approve",
		Evidence:   []string{digest(supervisedReviewEvidenceDomain, request.CandidateSnapshotHash, request.EvidenceRootHash)},
	}, nil
}

// Publish promotes the candidate only after proving the workspace still matches
// the frozen candidate and no managed process is alive.
func (runtime *supervisedRuntime) Publish(ctx context.Context, request chain.PromotionRequest) (chain.PromotionDecision, error) {
	fresh, err := runtime.captureManifest(ctx, "candidate", "supervised-promotion", runtime.acceptedHash)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	if changed := supervisedEntryMismatches(runtime.acceptedEntries, fresh.Manifest.Entries); len(changed) > 0 {
		return chain.PromotionDecision{
			ReceiptID:     runtime.nextID(),
			OccurredAt:    runtime.clock().UTC().Format(timestampLayout),
			Outcome:       "failed",
			ErrorEvidence: []string{digest(supervisedPromotionBlockedDomain, strings.Join(changed, "\n"))},
		}, nil
	}
	if _, err := runtime.processStopEvidence(); err != nil {
		return chain.PromotionDecision{
			ReceiptID:     runtime.nextID(),
			OccurredAt:    runtime.clock().UTC().Format(timestampLayout),
			Outcome:       "failed",
			ErrorEvidence: []string{digest(supervisedPromotionBlockedDomain, "a managed process of this run could not be proven stopped")},
		}, nil
	}
	return chain.PromotionDecision{
		ReceiptID:            runtime.nextID(),
		OccurredAt:           runtime.clock().UTC().Format(timestampLayout),
		Outcome:              "completed",
		AcceptedSnapshotHash: request.CandidateSnapshotHash,
		WriterStopEvidence:   []string{digest(supervisedWriterStopDomain, request.RunID, request.TaskID)},
		LeaseEvidence:        []string{digest(supervisedLeaseDomain, runtime.runID)},
		Evidence:             []string{digest(supervisedPromotionEvidenceDomain, request.ReviewReceiptHash)},
	}, nil
}

// RunPostflight is the required agent postflight port. A supervised-only run
// never records AgentRunner facts, so the chain never consults it; it refuses
// closed rather than skipping the acceptance gate (mirrors localRuntime).
func (runtime *supervisedRuntime) RunPostflight(context.Context, chain.PostflightRequest) (chain.PostflightDecision, error) {
	return chain.PostflightDecision{}, fmt.Errorf("%w: the supervised assembly has no agent postflight for run %s", chain.ErrPostflightUnavailable, runtime.runID)
}

// RunSupervisedPostflight re-derives the five frozen facts from the persisted
// artifacts and reconciles them field by field with the facts the step carries.
func (runtime *supervisedRuntime) RunSupervisedPostflight(ctx context.Context, request chain.SupervisedPostflightRequest) (chain.SupervisedPostflightDecision, error) {
	derived, err := runtime.deriveSupervisedFacts(ctx)
	if err != nil {
		return chain.SupervisedPostflightDecision{
			Outcome:       chain.PostflightUncertain,
			ErrorEvidence: []string{digest(supervisedPromotionBlockedDomain, err.Error())},
		}, nil
	}
	if err := supervisedFactsEqual(derived, request.Facts); err != nil {
		return chain.SupervisedPostflightDecision{
			Outcome:       chain.PostflightFailed,
			ErrorEvidence: []string{digest(supervisedPromotionBlockedDomain, err.Error())},
		}, nil
	}
	return chain.SupervisedPostflightDecision{
		Outcome:  chain.PostflightPassed,
		Evidence: supervisedFactsList(request.Facts),
	}, nil
}

// deriveSupervisedFacts recomputes the five fact digests from the persisted
// artifacts. It reads only files this assembly wrote.
func (runtime *supervisedRuntime) deriveSupervisedFacts(context.Context) (chain.SupervisedFrozenFacts, error) {
	receipt, found, err := runtime.loadReceipt()
	if err != nil {
		return chain.SupervisedFrozenFacts{}, err
	}
	if !found {
		return chain.SupervisedFrozenFacts{}, fmt.Errorf("%w: the handoff receipt is missing", ErrSupervisedNotOpen)
	}
	report, found, err := runtime.loadReport()
	if err != nil {
		return chain.SupervisedFrozenFacts{}, err
	}
	if !found {
		return chain.SupervisedFrozenFacts{}, fmt.Errorf("%w: the reduced/unknown report is missing", ErrSupervisedNotOpen)
	}
	scanEvidenceHash, err := supervisedScanEvidenceHash(runtime.scanReportPath())
	if err != nil {
		return chain.SupervisedFrozenFacts{}, err
	}
	scanReport, _, err := runtime.loadScanReport()
	if err != nil {
		return chain.SupervisedFrozenFacts{}, err
	}
	gateEvidence, found, err := runtime.loadGateEvidence()
	if err != nil {
		return chain.SupervisedFrozenFacts{}, err
	}
	if !found {
		return chain.SupervisedFrozenFacts{}, fmt.Errorf("%w: the gate evidence is missing", ErrSupervisedGateUnavailable)
	}
	return chain.SupervisedFrozenFacts{
		HandoffReceiptHash:      receipt.ReceiptHash,
		ScanEvidenceHash:        scanEvidenceHash,
		ProcessStopEvidenceHash: scanReport.ProcessStopHash,
		GateEvidenceHash:        digest(supervisedGateEvidenceDomain, append([]string{gateEvidence.RunID}, gateEvidence.Hashes...)...),
		ReportHash:              report.RecordHash,
	}, nil
}

// Stop reports the recorded stop evidence location. This run starts no process
// of its own.
func (runtime *supervisedRuntime) Stop(context.Context, string) ([]string, error) {
	return []string{digest(supervisedStopperDomain, runtime.runID)}, nil
}

// Reconcile is a read-only assembly inspection: the assembly holds no partial
// state to repair.
func (runtime *supervisedRuntime) Reconcile(context.Context, string) error { return nil }

// --- helpers -------------------------------------------------------------------

func (runtime *supervisedRuntime) ensureWorkspace() error {
	if runtime.workspaceReady {
		return nil
	}
	if err := os.MkdirAll(runtime.workspaceRoot, 0o755); err != nil {
		return fmt.Errorf("%w: create workspace: %v", ErrInvalidConfig, err)
	}
	runtime.workspaceReady = true
	return nil
}

func (runtime *supervisedRuntime) nextEventID() string { return runtime.nextID() }

func (runtime *supervisedRuntime) nextID() string {
	if runtime.ids != nil {
		return runtime.ids()
	}
	runtime.events++
	return fmt.Sprintf("supervised-id-%d", runtime.events)
}

func (runtime *supervisedRuntime) idSource() chain.IDSource {
	if runtime.ids != nil {
		return runtime.ids
	}
	return nil
}

func (runtime *supervisedRuntime) ledgerPath() string {
	return filepath.Join(SupervisedDir(runtime.runRoot), supervisedInteractionsName)
}

func (runtime *supervisedRuntime) confirmationsPath() string {
	return filepath.Join(SupervisedDir(runtime.runRoot), supervisedConfirmationsName)
}

func (runtime *supervisedRuntime) deliveryPath() string {
	return filepath.Join(SupervisedDir(runtime.runRoot), supervisedDeliveryName)
}

func (runtime *supervisedRuntime) sessionPath() string {
	return filepath.Join(SupervisedDir(runtime.runRoot), supervisedSessionName)
}

func (runtime *supervisedRuntime) frozenPolicyPath() string {
	return filepath.Join(SupervisedDir(runtime.runRoot), supervisedFrozenPolicyName)
}

func (runtime *supervisedRuntime) receiptPath() string {
	return filepath.Join(SupervisedDir(runtime.runRoot), supervisedReceiptName)
}

func (runtime *supervisedRuntime) reportPath() string { return SupervisedReportPath(runtime.runRoot) }

func (runtime *supervisedRuntime) scansDir() string {
	return filepath.Join(SupervisedDir(runtime.runRoot), supervisedScansDirectoryName)
}

func (runtime *supervisedRuntime) beforeManifestPath() string {
	return filepath.Join(runtime.scansDir(), supervisedBeforeManifestName)
}

func (runtime *supervisedRuntime) scanReportPath() string {
	return filepath.Join(runtime.scansDir(), supervisedScanReportName)
}

func (runtime *supervisedRuntime) gateEvidencePath() string {
	return filepath.Join(runtime.scansDir(), supervisedGateEvidenceName)
}

func (runtime *supervisedRuntime) loadLedger() (chain.OperatorInteractionLedger, error) {
	return loadOperatorInteractionLedger(runtime.ledgerPath())
}

func (runtime *supervisedRuntime) loadDelivery() (chain.SupervisedDelivery, bool, error) {
	var delivery chain.SupervisedDelivery
	found, err := readSupervisedJSON(runtime.deliveryPath(), &delivery)
	return delivery, found, err
}

func (runtime *supervisedRuntime) loadSession() (chain.HandoffSession, bool, error) {
	var session chain.HandoffSession
	found, err := readSupervisedJSON(runtime.sessionPath(), &session)
	return session, found, err
}

// loadFrozenPolicy reads the persisted freeze of the handoff policy that was in
// force when the work window opened.
func (runtime *supervisedRuntime) loadFrozenPolicy() (supervisedHandoffPolicySnapshot, bool, error) {
	var snapshot supervisedHandoffPolicySnapshot
	found, err := readSupervisedJSON(runtime.frozenPolicyPath(), &snapshot)
	return snapshot, found, err
}

func (runtime *supervisedRuntime) loadReceipt() (chain.HandoffReceiptRecord, bool, error) {
	var receipt chain.HandoffReceiptRecord
	found, err := readSupervisedJSON(runtime.receiptPath(), &receipt)
	return receipt, found, err
}

func (runtime *supervisedRuntime) loadReport() (chain.SupervisedReportRecord, bool, error) {
	var report chain.SupervisedReportRecord
	found, err := readSupervisedJSON(runtime.reportPath(), &report)
	return report, found, err
}

func (runtime *supervisedRuntime) loadBeforeManifest() (snapshot.SnapshotManifest, bool, error) {
	var manifest snapshot.SnapshotManifest
	found, err := readSupervisedJSON(runtime.beforeManifestPath(), &manifest)
	return manifest, found, err
}

func (runtime *supervisedRuntime) loadScanReport() (supervisedScanReport, bool, error) {
	var report supervisedScanReport
	found, err := readSupervisedJSON(runtime.scanReportPath(), &report)
	return report, found, err
}

func (runtime *supervisedRuntime) loadGateEvidence() (supervisedGateEvidence, bool, error) {
	var gateEvidence supervisedGateEvidence
	found, err := readSupervisedJSON(runtime.gateEvidencePath(), &gateEvidence)
	return gateEvidence, found, err
}

func (runtime *supervisedRuntime) loadConfirmationLedger() (chain.SupervisedRiskConfirmationLedger, error) {
	ledger := chain.SupervisedRiskConfirmationLedger{}
	_, err := readSupervisedJSONLines(runtime.confirmationsPath(), func(line []byte) error {
		record, err := chain.DecodeSupervisedRiskConfirmation(line)
		if err != nil {
			return err
		}
		return ledger.Append(record)
	})
	if err != nil {
		return chain.SupervisedRiskConfirmationLedger{}, err
	}
	return ledger, nil
}

// supervisedPendingRequest finds the binding risk-confirmation request of a step.
func supervisedPendingRequest(ledger chain.OperatorInteractionLedger, runID, taskID, stepID string) (chain.OperatorInteractionRecord, bool) {
	for _, record := range ledger.Records {
		if record.Request == nil {
			continue
		}
		request := record.Request
		if request.RunID == runID && request.TaskID == taskID && request.StepID == stepID && request.Attempt == 1 {
			return record, true
		}
	}
	return chain.OperatorInteractionRecord{}, false
}

// supervisedInteractionResponse finds the operator answer of an interaction.
func supervisedInteractionResponse(ledger chain.OperatorInteractionLedger, interactionID string) (chain.OperatorInteractionRecord, bool) {
	for _, record := range ledger.Records {
		if record.Response != nil && record.Response.InteractionID == interactionID {
			return record, true
		}
	}
	return chain.OperatorInteractionRecord{}, false
}

// persistFrozenHandoffPolicy freezes the handoff policy in force before the work
// window becomes durable. It is idempotent: an existing snapshot with the same
// canonical digest converges, while a divergent snapshot is refused rather than
// silently overwritten, so a later configuration change can never rewrite an
// already-frozen window.
func (runtime *supervisedRuntime) persistFrozenHandoffPolicy() error {
	hash := supervisedHandoffPolicyHash(runtime.policy)
	existing, found, err := runtime.loadFrozenPolicy()
	if err != nil {
		return err
	}
	if found {
		if existing.PolicyHash != hash {
			return fmt.Errorf("%w: the persisted frozen handoff policy differs from the policy in force", ErrSupervisedPolicyDrift)
		}
		return nil
	}
	return writeSupervisedJSON(runtime.frozenPolicyPath(), supervisedHandoffPolicySnapshot{Policy: runtime.policy, PolicyHash: hash})
}

// enforceFrozenHandoffPolicy is the complete-return guard against handoff-policy
// drift (blind review High). The work window was opened under a policy that was
// frozen then; the timeout, the authorized scope and the after-return gates are
// all enforced from that policy, so a chain configuration change between open
// and return must not be able to bend them. The guard fails closed before any
// durable write: a mismatch freezes no receipt, no report and no chain
// transition.
//
// Three conditions must hold, all anchored on the same canonical digest. First,
// the store-local snapshot must exist and match its own recorded digest. Second,
// the snapshot digest must equal the handoffPolicyHash the session recorded at
// window-open, so the session and the snapshot cannot disagree. Third, the
// policy in force now must digest to that same value.
//
// On success the runtime adopts the snapshot, so the scans and gates below are
// enforced from the window's policy and not from any swapped configuration.
func (runtime *supervisedRuntime) enforceFrozenHandoffPolicy(session chain.HandoffSession) error {
	snapshot, found, err := runtime.loadFrozenPolicy()
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("%w: the frozen handoff policy snapshot is missing", ErrSupervisedPolicyDrift)
	}
	if supervisedHandoffPolicyHash(snapshot.Policy) != snapshot.PolicyHash {
		return fmt.Errorf("%w: the frozen handoff policy snapshot does not match its own digest", ErrSupervisedPolicyDrift)
	}
	if snapshot.PolicyHash != session.HandoffPolicyHash {
		return fmt.Errorf("%w: the frozen handoff policy does not match the window session", ErrSupervisedPolicyDrift)
	}
	if supervisedHandoffPolicyHash(runtime.policy) != session.HandoffPolicyHash {
		return fmt.Errorf("%w: the policy in force changed since the window opened", ErrSupervisedPolicyDrift)
	}
	runtime.policy = snapshot.Policy
	return nil
}

// supervisedHandoffPolicyHash digests the frozen handoff policy of the step.
func supervisedHandoffPolicyHash(policy chain.HandoffPolicy) string {
	parts := []string{fmt.Sprintf("%d", policy.HandoffTimeoutMs), policy.InputPolicy}
	parts = append(parts, policy.AllowedTargets...)
	parts = append(parts, policy.ReturnActions...)
	parts = append(parts, policy.HooksAfterReturn...)
	return digest(supervisedHandoffPolicyDomain, parts...)
}

// supervisedFactsList returns the five frozen facts in the frozen positional
// routing order. It mirrors chain's unexported supervisedFactsEvidence; the
// order is part of the frozen contract (contracts section 7.1c).
func supervisedFactsList(facts chain.SupervisedFrozenFacts) []string {
	return []string{facts.HandoffReceiptHash, facts.ScanEvidenceHash, facts.ProcessStopEvidenceHash, facts.GateEvidenceHash, facts.ReportHash}
}

// supervisedFactsEqual reconciles two frozen fact sets field by field.
func supervisedFactsEqual(left, right chain.SupervisedFrozenFacts) error {
	for index, pair := range [][2]string{
		{left.HandoffReceiptHash, right.HandoffReceiptHash},
		{left.ScanEvidenceHash, right.ScanEvidenceHash},
		{left.ProcessStopEvidenceHash, right.ProcessStopEvidenceHash},
		{left.GateEvidenceHash, right.GateEvidenceHash},
		{left.ReportHash, right.ReportHash},
	} {
		if pair[0] != pair[1] {
			return fmt.Errorf("%w: supervised fact %d does not match the persisted artifact", chain.ErrInvalidPostflightFacts, index)
		}
	}
	return nil
}

// supervisedScanEvidenceHash recomputes the scan-evidence digest from the
// persisted scan report, so the postflight re-derives it rather than trusting
// the step evidence.
func supervisedScanEvidenceHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%w: read scan report: %v", ErrInvalidConfig, err)
	}
	var report supervisedScanReport
	if err := evidence.DecodeStrictJSON(data, &report); err != nil {
		return "", fmt.Errorf("%w: decode scan report: %v", ErrInvalidConfig, err)
	}
	canonical, err := evidence.EncodeCanonical(report)
	if err != nil {
		return "", err
	}
	return evidence.Digest(supervisedScanDomain, canonical), nil
}

// supervisedEntryMismatches lists the entries whose signature differs between
// two manifests, mirroring the AgentRunner assembly helper.
func supervisedEntryMismatches(frozen, fresh []snapshot.Entry) []string {
	frozenByPath := make(map[string]string, len(frozen))
	for _, entry := range frozen {
		frozenByPath[entry.Path] = supervisedEntrySignature(entry)
	}
	freshByPath := make(map[string]string, len(fresh))
	for _, entry := range fresh {
		freshByPath[entry.Path] = supervisedEntrySignature(entry)
	}
	var changed []string
	for path, signature := range freshByPath {
		if prior, existed := frozenByPath[path]; !existed {
			changed = append(changed, "added "+path)
		} else if prior != signature {
			changed = append(changed, "changed "+path)
		}
	}
	for path := range frozenByPath {
		if _, exists := freshByPath[path]; !exists {
			changed = append(changed, "removed "+path)
		}
	}
	sort.Strings(changed)
	return changed
}

func supervisedEntrySignature(entry snapshot.Entry) string {
	return fmt.Sprintf("%s|%s|%s|%d|%t|%t|%s|%s|%s",
		entry.Path, entry.Type, entry.ContentHash, entry.SizeBytes, entry.Executable, entry.ReadOnly, entry.Role, entry.Target, entry.TargetHash)
}

// supervisedScanReport is the post-return scan artifact.
type supervisedScanReport struct {
	RunID              string   `json:"runId"`
	TaskID             string   `json:"taskId"`
	StepID             string   `json:"stepId"`
	Attempt            int      `json:"attempt"`
	BeforeManifestHash string   `json:"beforeManifestHash"`
	AfterManifestHash  string   `json:"afterManifestHash"`
	DiffHash           string   `json:"diffHash"`
	ChangedPaths       []string `json:"changedPaths"`
	EntryCount         int      `json:"entryCount"`
	// ChangeSetInScope is true only when every path in the post-return change
	// set lies inside the authorized scope (the resolved isolated workspace
	// root). It is NOT a claim that no out-of-scope write happened anywhere:
	// writes outside the observable workspace are unobservable from ProofRail's
	// position and are carried by the report's unknown dimensions
	// (external-side-effects), never asserted clean here.
	ChangeSetInScope  bool     `json:"changeSetInScope"`
	SecretScanClean   bool     `json:"secretScanClean"`
	FileTypeScanClean bool     `json:"fileTypeScanClean"`
	SideEffectClean   bool     `json:"sideEffectClean"`
	ProcessStopHash   string   `json:"processStopHash"`
	GateHashes        []string `json:"gateHashes"`
}

// supervisedGateEvidence is the persisted after-return gate result set.
type supervisedGateEvidence struct {
	RunID  string   `json:"runId"`
	Hashes []string `json:"hashes"`
}

// supervisedHandoffPolicySnapshot is the store-local freeze of the handoff
// policy that was in force when the work window opened. Like every other
// <runRoot>/supervised/ file it is store-local protocol, not a wire record: no
// schema and no external consumer depends on it. PolicyHash is produced by
// supervisedHandoffPolicyHash, the same single convention that populates
// HandoffSession.HandoffPolicyHash, so the mode keeps exactly one policy digest.
type supervisedHandoffPolicySnapshot struct {
	Policy     chain.HandoffPolicy `json:"policy"`
	PolicyHash string              `json:"policyHash"`
}

// writeSupervisedJSON writes one canonical JSON document atomically. It never
// writes with a BOM and always terminates with a single newline.
func writeSupervisedJSON(path string, value any) error {
	if value == nil {
		return fmt.Errorf("%w: nil document", ErrInvalidConfig)
	}
	canonical, err := evidence.EncodeCanonical(value)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".supervised-*")
	if err != nil {
		return err
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if err = temporary.Chmod(0o600); err == nil {
		_, err = temporary.Write(append(canonical, '\n'))
	}
	if err == nil {
		err = temporary.Sync()
	}
	if closeErr := temporary.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(temporaryPath, path)
}

// readSupervisedJSON reads one canonical JSON document, reporting whether the
// file exists at all.
func readSupervisedJSON(path string, target any) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	if err := evidence.DecodeStrictJSON(data, target); err != nil {
		return false, fmt.Errorf("%w: %s: %v", ErrInvalidConfig, path, err)
	}
	return true, nil
}

// readSupervisedJSONLines reads an append-only JSONL store and applies one
// record per non-empty line.
func readSupervisedJSONLines(path string, apply func([]byte) error) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if err := apply([]byte(trimmed)); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// appendSupervisedConfirmation appends one confirmation to the no-replace
// confirmation ledger. An identical record converges; a divergent record for the
// same step attempt is refused by the in-memory ledger check.
func appendSupervisedConfirmation(path string, record chain.SupervisedRiskConfirmationRecord) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	canonical, err := evidence.EncodeCanonical(record)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(append(canonical, '\n')); err != nil {
		return err
	}
	return file.Sync()
}
