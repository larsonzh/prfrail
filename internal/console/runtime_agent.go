package console

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
	"github.com/larsonzh/prfrail/internal/tickets"
)

// This file is the production AgentRunner assembly of the console package. It
// adds no chain port, no Go interface and no wire record: every port it hands to
// the engine is an implementation of an interface the chain (or the adapters
// layer) already defines. The assembly is fail-closed before the engine exists,
// and the only terminal it concludes is the one the adapters layer proved.

// Console-local evidence domains. They are internal digests, never wire records:
// the console reconciles its own assembly proof without touching a protocol.
const (
	agentRunnerEvidenceRootDomain              = "proofrail:agent-runner-evidence-root:1\n"
	agentRunnerCandidateEvidenceDomain         = "proofrail:agent-runner-candidate-evidence:1\n"
	agentRunnerReviewPolicyDomain              = "proofrail:agent-runner-review-policy:1\n"
	agentRunnerReviewEvidenceDomain            = "proofrail:agent-runner-review-evidence:1\n"
	agentRunnerStopperDomain                   = "proofrail:agent-runner-stopper:1\n"
	agentRunnerWriterStopDomain                = "proofrail:agent-runner-writer-stop:1\n"
	agentRunnerLeaseDomain                     = "proofrail:agent-runner-lease:1\n"
	agentRunnerPromotionEvidenceDomain         = "proofrail:agent-runner-promotion-evidence:1\n"
	agentRunnerPromotionBlockedDomain          = "proofrail:agent-runner-promotion-blocked:1\n"
	agentRunnerPlatformDomain                  = "proofrail:agent-runner-platform:1\n"
	agentRunnerDefaultEvidenceAge              = time.Hour
	agentRunnerDefaultMaximumRequests          = 1
	agentRunnerDefaultRunTimeout               = 30 * time.Minute
	agentRunnerDefaultMaxLogBytes        int64 = 1 << 20
	agentRunnerDefaultVersionFlag              = "--version"
	agentRunnerSnapshotStoreDirectory          = "snapshot-store"
	agentRunnerWorkspacesDirectory             = "workspaces"
	agentRunnerAssemblyReviewPolicyLabel       = "agent-runner-assembly-policy"
)

// AgentRunnerCLIOptions pins the offline CLI the assembly runs and the bounds it
// runs it under. Executable and Version are the pin; the launcher refuses any
// request that names another executable.
type AgentRunnerCLIOptions struct {
	Executable   string
	Version      string
	VersionArgs  []string
	RunArgs      []string
	Env          []string
	EnvAllow     []string
	CheckTimeout time.Duration
	Grace        time.Duration
	Timeout      time.Duration
	VerifyWait   time.Duration
	MaxLogBytes  int64
}

// AgentRunnerAssemblyOptions pins one production AgentRunner run. Every record is
// an already-decoded declaration (declared, not probed): the assembly validates
// it, it never probes a candidate and never touches the network.
type AgentRunnerAssemblyOptions struct {
	RunID      string
	RunRoot    string
	Definition chain.Definition

	RequestRecord      adapters.AgentRunnerRequestRecord
	CapabilityRecord   adapters.AgentRunnerCapabilityRecord
	EnforcementRecord  *adapters.AgentRunnerEnforcementRecord
	AvailabilityRecord adapters.AIAvailabilityRecord
	AvailabilityPolicy adapters.AIAvailabilityPolicy

	AuthorizationLedger chain.AuthorizationLedger
	CostLedger          *tickets.CostLedger

	AgentCLI AgentRunnerCLIOptions

	// WorkspaceRoot overrides the derived isolated workspace. The CLI never sets
	// it; the assembly precondition tests use it to force the workspace/replay
	// overlap the replay store must refuse.
	WorkspaceRoot string

	Clock func() time.Time
	IDs   func() string
}

// ExecuteAgentRunnerRun is the production AgentRunner assembly entry. It builds
// the offline run, drives the chain through one dispatch, publishes the terminal
// the adapters layer proved and walks the downstream postflight/review/promotion
// path. Every failure before the engine exists is fail-closed: the run-root
// marker is created first (step 2) and no process is ever spawned, so a refusal
// leaves no live process and no terminal.
func ExecuteAgentRunnerRun(ctx context.Context, options AgentRunnerAssemblyOptions) (RunSummary, error) {
	if err := ctx.Err(); err != nil {
		return RunSummary{}, err
	}
	clock := options.Clock
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	ids := options.IDs

	// Step 1: assembly preconditions. The definition must be exactly one
	// isolated-workspace code step, and every required record must decode and
	// validate. A nil enforcement record is not refused here: it flows to the
	// frozen admission, which refuses a non-compatible capability without it.
	taskID, stepID, err := agentRunnerAssemblyStep(options.Definition)
	if err != nil {
		return RunSummary{}, err
	}
	if err := adapters.ValidateAgentRunnerRequestRecord(options.RequestRecord); err != nil {
		return RunSummary{}, fmt.Errorf("%w: request record is missing or invalid: %w", ErrInvalidConfig, err)
	}
	record := options.RequestRecord.Request
	if record.RunID != options.RunID || record.TaskID != taskID || record.StepID != stepID {
		return RunSummary{}, fmt.Errorf("%w: request record does not bind the assembled step", ErrInvalidConfig)
	}
	if err := adapters.ValidateAgentRunnerCapabilityRecord(options.CapabilityRecord); err != nil {
		return RunSummary{}, fmt.Errorf("%w: capability record is missing or invalid: %w", ErrInvalidConfig, err)
	}
	if options.EnforcementRecord != nil {
		if err := adapters.ValidateAgentRunnerEnforcementRecord(*options.EnforcementRecord); err != nil {
			return RunSummary{}, fmt.Errorf("%w: enforcement record is invalid: %w", ErrInvalidConfig, err)
		}
	}
	if err := adapters.ValidateAIAvailabilityRecord(options.AvailabilityRecord); err != nil {
		return RunSummary{}, fmt.Errorf("%w: availability record is missing or invalid: %w", ErrInvalidConfig, err)
	}
	if options.CostLedger == nil {
		return RunSummary{}, fmt.Errorf("%w: cost ledger is required", ErrInvalidConfig)
	}
	if strings.TrimSpace(options.RunRoot) == "" {
		return RunSummary{}, fmt.Errorf("%w: run root is required", ErrInvalidConfig)
	}
	if err := agentRunnerValidateCLI(options.AgentCLI); err != nil {
		return RunSummary{}, err
	}

	// Step 2: the run-root marker comes first. The replay store refuses to be
	// constructed without it, so the order is a correctness rule, not a script.
	eventLog := runEventLogPath(options.RunRoot)
	if err := ensureAgentRunnerEventLogMarker(eventLog); err != nil {
		return RunSummary{}, err
	}

	// Step 3: the production replay store, derived from the durable run root.
	store, err := adapters.NewAgentRunnerReplayStoreForRun(options.RunRoot, options.RunID)
	if err != nil {
		return RunSummary{}, err
	}

	workspaceRoot := strings.TrimSpace(options.WorkspaceRoot)
	if workspaceRoot == "" {
		workspaceRoot = agentRunnerWorkspacePath(options.RunRoot, taskID, record.Attempt)
	}

	// Step 4: the isolated workspace and the replay root may never overlap.
	if err := store.RejectWriterRoots(workspaceRoot); err != nil {
		return RunSummary{}, err
	}

	// Step 5: the snapshot store and the shared capture protocol.
	snapshotRoot := filepath.Join(options.RunRoot, agentRunnerSnapshotStoreDirectory)
	snapshotStore, err := snapshot.NewStore(snapshotRoot, 0)
	if err != nil {
		return RunSummary{}, err
	}
	capturer, err := adapters.NewAgentRunnerSnapshotCapturer(snapshotStore, options.RunID, taskID, record.Attempt)
	if err != nil {
		return RunSummary{}, err
	}

	// Step 6: the production pinned-CLI launcher.
	launcher, err := adapters.NewAgentRunnerPinnedCLILauncher(adapters.AgentRunnerPinnedCLIConfig{
		Executable:   options.AgentCLI.Executable,
		Version:      options.AgentCLI.Version,
		VersionArgs:  append([]string(nil), options.AgentCLI.VersionArgs...),
		RunArgs:      append([]string(nil), options.AgentCLI.RunArgs...),
		EnvAllow:     append([]string(nil), options.AgentCLI.EnvAllow...),
		CheckTimeout: options.AgentCLI.CheckTimeout,
		Grace:        options.AgentCLI.Grace,
	}, store)
	if err != nil {
		return RunSummary{}, err
	}

	// Step 7: the replay-aware dispatcher. Its process config and the (A′)
	// run request are built from the same values, so the launch is addressable.
	process := adapters.AgentRunnerProcessConfig{
		Command: options.AgentCLI.Executable,
		Args:    append([]string(nil), options.AgentCLI.RunArgs...),
		Env:     append([]string(nil), options.AgentCLI.Env...),
		Grace:   options.AgentCLI.Grace,
	}
	dispatcher := &adapters.AgentRunnerReplayDispatcher{
		Store:          store,
		Launcher:       launcher,
		Process:        process,
		LedgerSnapshot: func() (chain.AuthorizationLedger, error) { return options.AuthorizationLedger, nil },
		CostLedger:     options.CostLedger,
		RequestRecord:  options.RequestRecord,
		Clock:          clock,
	}

	// Step 8: the composite admission. Enforcement may be nil; the frozen
	// validator decides whether it is required.
	admission := &adapters.AgentRunnerCompositeAdmission{
		RequestRecord:       options.RequestRecord,
		AuthorizationLedger: options.AuthorizationLedger,
		CostLedger:          options.CostLedger,
		AvailabilityRecord:  options.AvailabilityRecord,
		AvailabilityPolicy:  options.AvailabilityPolicy,
		CapabilityRecord:    options.CapabilityRecord,
		EnforcementRecord:   options.EnforcementRecord,
		AdmissionPolicy:     agentRunnerAdmissionPolicy(clock()),
		Clock:               clock,
	}

	// Step 9: the engine-owned step intent projection.
	intents, err := adapters.NewAgentRunnerStepIntentPreparer(options.RequestRecord, options.Definition)
	if err != nil {
		return RunSummary{}, err
	}

	// Step 10: the production postflight port, bound to the request record.
	postflightPort, err := adapters.NewAgentRunnerPostflight(store, options.RequestRecord, snapshotRoot)
	if err != nil {
		return RunSummary{}, err
	}

	// Step 11: the production terminal publisher.
	publisher := &adapters.AgentRunnerTerminalPublisher{
		Store:         store,
		CostLedger:    options.CostLedger,
		RequestRecord: options.RequestRecord,
		Clock:         clock,
	}

	// Step 12: the runtime ports plus the (A′) wait carrier. The carrier shares the
	// dispatcher's launcher and store instances, otherwise the live launch cannot be
	// addressed.
	assemblyRuntime := &agentRunnerRuntime{
		runID:         options.RunID,
		store:         store,
		workspaceRoot: workspaceRoot,
		snapshotStore: snapshotStore,
		capturer:      capturer,
		requestRecord: options.RequestRecord,
		clock:         clock,
		ids:           ids,
	}
	assemblyRuntime.waitRun = &adapters.AgentRunnerPinnedCLIRun{Launcher: launcher, Store: store, Clock: clock}
	assemblyRuntime.steps = chain.StepRouter{
		AgentRunner:          &preManifestAgentRunnerPort{owner: assemblyRuntime, inner: dispatcher, capturer: capturer},
		AgentRunnerAdmission: admission,
	}

	// Step 13: the engine.
	eventStore, err := chain.NewFileEventStore(eventLog)
	if err != nil {
		return RunSummary{}, err
	}
	engine, err := chain.New(ctx, chain.Options{
		RunID:       options.RunID,
		Definition:  options.Definition,
		Events:      eventStore,
		Baselines:   assemblyRuntime,
		Workspaces:  assemblyRuntime,
		Steps:       assemblyRuntime,
		StepIntents: intents,
		Acceptance:  assemblyRuntime,
		Reviewer:    assemblyRuntime,
		Publisher:   assemblyRuntime,
		Postflight:  postflightPort,
		Stopper:     assemblyRuntime,
		Reconciler:  assemblyRuntime,
		Clock:       clock,
		IDs:         assemblyRuntime.nextEventID,
	})
	if err != nil {
		return RunSummary{}, err
	}

	// fail is the single exit for every failure after the engine exists. It is a
	// best-effort mitigation, NOT the root fix: the root cause (the launch receipt
	// is recorded before the spawn) lives in
	// internal/adapters/agent_runner_replay_dispatcher.go and is registered as
	// candidate OB-57. Because the closure is the ONLY place a failure returns a
	// projection built from an error, no post-spawn failure can bypass the
	// abandoned-launch mitigation. The source-level invariant
	// TestAgentRunnerAssemblyFunnelsEveryPostSpawnFailureThroughTheFallback pins
	// that shape; a new failure path must call fail(...) and bump that test's
	// failure-return count in the same change.
	fail := func(err error) (RunSummary, error) {
		err = agentRunnerMitigateAbandonedLaunch(ctx, store, launcher, options.RequestRecord.Request.RequestID, err)
		return summarizeProjection(engine.Projection(), eventLog), err
	}

	// Drive: dispatch parks the step, the assembly waits/collects/maps, the
	// publisher terminalizes and the chain routes the terminal fact.
	runErr := engine.Run(ctx)
	if runErr == nil {
		return summarizeProjection(engine.Projection(), eventLog), nil
	}
	if !errors.Is(runErr, chain.ErrAwaitingAgentRunnerTerminal) {
		// Every refusal (durability, admission, platform) lands here, but a
		// non-awaiting error can also follow a spawn: the launch receipt is written
		// before the process starts, so the shared fail closure attempts the
		// best-effort try-stop for any launch this request recorded. The original
		// error is always preserved.
		return fail(runErr)
	}
	runRequest := agentRunnerRunRequest(options.RequestRecord, process, workspaceRoot, options.AgentCLI, capturer)
	outcome, err := assemblyRuntime.waitRun.WaitAgentRunnerDispatched(ctx, runRequest, assemblyRuntime.beforeManifest)
	if err != nil {
		return fail(err)
	}
	if _, err := publisher.PublishTerminal(ctx, outcome); err != nil {
		return fail(err)
	}
	terminal, err := adapters.LoadChainAgentRunnerTerminal(store, options.RequestRecord)
	if err != nil {
		return fail(err)
	}
	if _, err := engine.SubmitAgentRunnerTerminal(ctx, terminal); err != nil {
		return fail(err)
	}
	if outcome.Completion.Status != "completed" {
		return fail(fmt.Errorf("%w: agent runner terminal %s", chain.ErrStepFailed, outcome.Completion.Status))
	}
	if err := engine.Run(ctx); err != nil {
		return fail(err)
	}
	return summarizeProjection(engine.Projection(), eventLog), nil
}

// preManifestAgentRunnerPort captures the workspace manifest before the
// dispatcher may touch it. A failed capture stays a gap (stored as nil) and the
// dispatch still happens, exactly as the offline run degrades instead of
// aborting.
type preManifestAgentRunnerPort struct {
	owner    *agentRunnerRuntime
	inner    chain.AgentRunnerPort
	capturer *adapters.AgentRunnerSnapshotCapturer
}

var _ chain.AgentRunnerPort = (*preManifestAgentRunnerPort)(nil)

func (port *preManifestAgentRunnerPort) DispatchAgentRunner(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	if port.capturer != nil {
		if manifest, err := port.capturer.CaptureWorkspaceManifest(ctx, request.Workspace.Root); err == nil {
			port.owner.beforeManifest = manifest
		} else {
			port.owner.beforeManifest = nil
		}
	}
	return port.inner.DispatchAgentRunner(ctx, request)
}

// agentRunnerRuntime implements every chain port the assembly needs. It is a
// real implementation over the filesystem, the snapshot store and the production
// launcher: no port answers from memory except the pre-dispatch manifest the
// wrapper captured.
type agentRunnerRuntime struct {
	runID         string
	store         *adapters.AgentRunnerReplayStore
	workspaceRoot string
	snapshotStore *snapshot.Store
	capturer      *adapters.AgentRunnerSnapshotCapturer
	requestRecord adapters.AgentRunnerRequestRecord

	steps chain.StepRouter

	waitRun        *adapters.AgentRunnerPinnedCLIRun
	beforeManifest []byte

	acceptedEntries []snapshot.Entry

	workspaceReady bool
	events         int

	clock func() time.Time
	ids   func() string
}

var (
	_ chain.BaselinePort   = (*agentRunnerRuntime)(nil)
	_ chain.WorkspacePort  = (*agentRunnerRuntime)(nil)
	_ chain.StepPort       = (*agentRunnerRuntime)(nil)
	_ chain.AcceptancePort = (*agentRunnerRuntime)(nil)
	_ chain.ReviewerPort   = (*agentRunnerRuntime)(nil)
	_ chain.PublisherPort  = (*agentRunnerRuntime)(nil)
	_ chain.Stopper        = (*agentRunnerRuntime)(nil)
	_ chain.Reconciler     = (*agentRunnerRuntime)(nil)
)

// CaptureBaseline digests the workspace entries (entries only, never the capture
// timestamp) so the request record's declared parent snapshot is the measured
// one. The workspace is created first, because the engine captures the baseline
// before it materializes the task workspace.
func (runtime *agentRunnerRuntime) CaptureBaseline(ctx context.Context, runID string) (chain.SnapshotRef, error) {
	if err := runtime.ensureWorkspace(); err != nil {
		return chain.SnapshotRef{}, err
	}
	record := runtime.requestRecord.Request
	hash, err := adapters.AgentRunnerWorkspaceBaselineDigest(ctx, runtime.snapshotStore, runtime.workspaceRoot, runID, record.TaskID, record.Attempt)
	if err != nil {
		return chain.SnapshotRef{}, err
	}
	return chain.SnapshotRef{Hash: hash}, nil
}

// Materialize returns the one isolated workspace of this run root.
func (runtime *agentRunnerRuntime) Materialize(context.Context, string, string, int, chain.SnapshotRef) (chain.Workspace, error) {
	if err := runtime.ensureWorkspace(); err != nil {
		return chain.Workspace{}, err
	}
	return chain.Workspace{Root: runtime.workspaceRoot}, nil
}

// Execute routes the one step. The AgentRunner route is admission-first inside
// the router; the pre-manifest wrapper only captures before dispatch.
func (runtime *agentRunnerRuntime) Execute(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	return runtime.steps.Execute(ctx, request)
}

// Accept freezes the candidate: a real candidate capture of the workspace plus
// the digest of the recorded evidence inventory.
func (runtime *agentRunnerRuntime) Accept(ctx context.Context, runID, taskID string, attempt int, parent chain.SnapshotRef, workspace chain.Workspace) (chain.CandidateResult, error) {
	parentHash := parent.Hash
	captured, err := snapshot.Capture(ctx, snapshot.CaptureOptions{
		SourceDir:          workspace.Root,
		Kind:               "candidate",
		SnapshotID:         fmt.Sprintf("candidate-%s-%s-%d", runID, taskID, attempt),
		RunID:              runID,
		ParentSnapshotHash: &parentHash,
		Task:               &snapshot.TaskBinding{TaskID: taskID, Attempt: attempt},
		Store:              runtime.snapshotStore,
	})
	if err != nil {
		return chain.CandidateResult{}, err
	}
	runtime.acceptedEntries = captured.Manifest.Entries
	evidenceRoot, err := runtime.evidenceRootHash()
	if err != nil {
		return chain.CandidateResult{}, err
	}
	return chain.CandidateResult{
		Snapshot:          chain.SnapshotRef{Hash: captured.ManifestHash},
		EvidenceRootHash:  evidenceRoot,
		Evidence:          []string{evidence.Digest(agentRunnerCandidateEvidenceDomain, []byte(captured.ManifestHash))},
		CandidateProducer: evidence.Actor{Type: "agent", ID: runtime.requestRecord.Request.AdapterID},
	}, nil
}

// Review is a real policy review over the frozen candidate. The chain only ever
// reaches it from REVIEW_PENDING, so the postflight qualification is already
// proven; the port records the policy it applied.
func (runtime *agentRunnerRuntime) Review(_ context.Context, request chain.ReviewRequest) (chain.ReviewDecision, error) {
	if !evidence.ValidHash(request.CandidateSnapshotHash) || !evidence.ValidHash(request.EvidenceRootHash) {
		return chain.ReviewDecision{}, fmt.Errorf("%w: review request carries an invalid binding", ErrInvalidConfig)
	}
	return chain.ReviewDecision{
		ReceiptID:  runtime.nextEventID(),
		OccurredAt: runtime.clock().UTC().Format(timestampLayout),
		RecordedBy: evidence.Actor{Type: "policy", ID: agentRunnerAssemblyReviewPolicyLabel},
		ReviewMode: "policy",
		PolicyHash: digest(agentRunnerReviewPolicyDomain, agentRunnerAssemblyReviewPolicyLabel, "postflight-qualified-candidate"),
		Outcome:    "approve",
		Evidence:   []string{digest(agentRunnerReviewEvidenceDomain, request.CandidateSnapshotHash, request.EvidenceRootHash)},
	}, nil
}

// Publish promotes the candidate only after proving the workspace still matches
// the frozen candidate and no managed process of this run is alive. The process
// proof is fail-closed, exactly like postflight leg 5: a run that left no pid
// report blocks the promotion instead of reading as "nothing alive".
func (runtime *agentRunnerRuntime) Publish(ctx context.Context, request chain.PromotionRequest) (chain.PromotionDecision, error) {
	fresh, err := runtime.capturer.CaptureWorkspaceManifest(ctx, runtime.workspaceRoot)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	freshManifest, err := snapshot.DecodeSnapshotManifest(fresh)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	if changed := agentRunnerEntryMismatches(runtime.acceptedEntries, freshManifest.Manifest.Entries); len(changed) > 0 {
		return chain.PromotionDecision{
			ReceiptID:     runtime.nextEventID(),
			OccurredAt:    runtime.clock().UTC().Format(timestampLayout),
			Outcome:       "failed",
			ErrorEvidence: []string{digest(agentRunnerPromotionBlockedDomain, strings.Join(changed, "\n"))},
		}, nil
	}
	stopped, err := runtime.managedProcessesGone()
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	if !stopped {
		return chain.PromotionDecision{
			ReceiptID:     runtime.nextEventID(),
			OccurredAt:    runtime.clock().UTC().Format(timestampLayout),
			Outcome:       "failed",
			ErrorEvidence: []string{digest(agentRunnerPromotionBlockedDomain, "a managed process of this run is still alive")},
		}, nil
	}
	return chain.PromotionDecision{
		ReceiptID:            runtime.nextEventID(),
		OccurredAt:           runtime.clock().UTC().Format(timestampLayout),
		Outcome:              "completed",
		AcceptedSnapshotHash: request.CandidateSnapshotHash,
		WriterStopEvidence:   []string{digest(agentRunnerWriterStopDomain, request.RunID, request.TaskID)},
		LeaseEvidence:        []string{digest(agentRunnerLeaseDomain, runtime.runID)},
		Evidence:             []string{digest(agentRunnerPromotionEvidenceDomain, request.ReviewReceiptHash)},
	}, nil
}

// Stop reports the recorded stop evidence location. The only managed writers are
// the pinned CLI launches; their stop path is owned by the launcher.
func (runtime *agentRunnerRuntime) Stop(context.Context, string) ([]string, error) {
	return []string{digest(agentRunnerStopperDomain, runtime.runID)}, nil
}

// Reconcile is a read-only assembly inspection: the assembly holds no partial
// state to repair.
func (runtime *agentRunnerRuntime) Reconcile(context.Context, string) error { return nil }

func (runtime *agentRunnerRuntime) nextEventID() string {
	runtime.events++
	if runtime.ids != nil {
		return runtime.ids()
	}
	return fmt.Sprintf("agent-runner-event-%d", runtime.events)
}

func (runtime *agentRunnerRuntime) ensureWorkspace() error {
	if runtime.workspaceReady {
		return nil
	}
	if err := os.MkdirAll(runtime.workspaceRoot, 0o755); err != nil {
		return fmt.Errorf("%w: create workspace: %v", ErrInvalidConfig, err)
	}
	runtime.workspaceReady = true
	return nil
}

// evidenceRootHash digests the sorted artifact inventory of the request's
// evidence directory: it is the root the review binds, not a policy decision.
func (runtime *agentRunnerRuntime) evidenceRootHash() (string, error) {
	root, err := runtime.replayEvidenceDir()
	if err != nil {
		return "", err
	}
	var lines []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines = append(lines, filepath.ToSlash(relative)+" "+evidence.Digest("", data))
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(lines)
	return digest(agentRunnerEvidenceRootDomain, strings.Join(lines, "\n")), nil
}

func (runtime *agentRunnerRuntime) replayEvidenceDir() (string, error) {
	if runtime.store == nil {
		return "", fmt.Errorf("%w: replay store is not assembled", ErrInvalidConfig)
	}
	return runtime.store.RunEvidenceDir(runtime.requestRecord.Request.RequestID)
}

// The agent-stub pid report file names. They are the stub's convention, not a
// wire contract: this is a hand-synchronised coupling point, mirroring the
// unexported agentRunnerParentPIDReportFileName / agentRunnerChildPIDReportFileName
// constants of internal/adapters (declared in agent_runner_evidence.go and
// agent_runner_postflight.go). The two packages must stay equal, and both sides
// must be re-evaluated together when a real CLI candidate is wired in.
const (
	agentRunnerParentPIDReportFileName = "stub.pid"
	agentRunnerChildPIDReportFileName  = "stub-child.pid"
)

// agentRunnerMissingPIDReportDetail is the exact verdict text the adapters
// postflight also uses for the same condition (internal/adapters/agent_runner_postflight.go,
// agentRunnerMissingPIDReportDetail, leg 5). The two copies must stay
// byte-identical: it is a hand-synchronised coupling point, like the pid report
// file names above.
const agentRunnerMissingPIDReportDetail = "the run is missing a required pid report so the managed processes cannot be proven gone"

// managedProcessesGone reports whether every pid the pinned CLI reported is gone.
//
// The criterion is word-for-word leg 5 of
// internal/adapters/agent_runner_postflight.go (agentRunnerAlivePIDs): BOTH pid
// reports (parent and child) must exist, and any missing report or unobservable
// liveness answer is fail-closed. Both reports are written by the candidate, so a
// run that left either missing cannot be proven to have stopped its managed
// processes and the assembly must not publish. The caller reuses the adapters
// evidence sentinel rather than minting a new one (zero new sentinels).
func (runtime *agentRunnerRuntime) managedProcessesGone() (bool, error) {
	root, err := runtime.replayEvidenceDir()
	if err != nil {
		return false, err
	}
	missing := make([]string, 0, 2)
	alive := false
	names := []string{agentRunnerParentPIDReportFileName, agentRunnerChildPIDReportFileName}
	for _, name := range names {
		pids, present, err := readAgentRunnerEvidencePIDs(filepath.Join(root, name))
		if err != nil {
			return false, err
		}
		if !present {
			missing = append(missing, name)
			continue
		}
		for _, pid := range pids {
			identity, err := guard.InspectProcess(pid)
			if err != nil {
				// A pid that cannot be opened is the normal shape of a process that has
				// already exited. See the adapters counterpart agentRunnerAlivePIDsWith
				// for the cross-platform rationale; the liveness answer below is what
				// must never be silently ignored.
				continue
			}
			processAlive, err := guard.ProcessAlive(identity)
			if err != nil {
				return false, fmt.Errorf("%w: cannot determine whether pid %d from %s is alive: %v", adapters.ErrInvalidAgentRunnerEvidence, pid, name, err)
			}
			if processAlive {
				alive = true
			}
		}
	}
	// The missing check runs first, exactly like leg 5, so the two sides answer
	// the same verdict for the same evidence.
	if len(missing) > 0 {
		// Reuse the existing adapters evidence sentinel rather than minting a new
		// one: the fault is AgentRunner evidence that cannot be read, which is what
		// ErrInvalidAgentRunnerEvidence names, and Publish's caller can identify it
		// with errors.Is.
		return false, fmt.Errorf("%w: %s", adapters.ErrInvalidAgentRunnerEvidence, agentRunnerMissingPIDReportDetail)
	}
	if alive {
		return false, nil
	}
	return true, nil
}

// readAgentRunnerEvidencePIDs reads a one-pid-per-line report and reports
// whether the file exists at all; a missing report is not an error, the caller
// decides. An unparsable line is fail-closed instead of being ignored, and it
// reuses the adapters evidence sentinel (zero new sentinels). The present flag
// mirrors agentRunnerReadPIDs in the adapters postflight: it is the fail-closed
// support both legs need.
func readAgentRunnerEvidencePIDs(path string) ([]int, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}
	// COUPLING: this parse loop must stay word-for-word equal to
	// agentRunnerReadPIDs in internal/adapters/agent_runner_postflight.go. Both
	// sides answer the same verdict for the same report, so the parse expression is
	// kept byte-identical by hand, like agentRunnerMissingPIDReportDetail and the
	// pid report file names above. strconv.Atoi(strings.TrimSpace(line)) rejects
	// the numeric-prefix shapes fmt.Sscanf silently accepted ("12abc", "7.0",
	// "1_0"), closing the fail-open gap.
	pids := make([]int, 0, 2)
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || pid <= 0 {
			return nil, false, fmt.Errorf("%w: unreadable pid %q in %s", adapters.ErrInvalidAgentRunnerEvidence, line, path)
		}
		pids = append(pids, pid)
	}
	return pids, true, nil
}

// agentRunnerEntryMismatches lists the entries whose signature differs between
// the frozen candidate and the fresh capture, including added and removed ones.
func agentRunnerEntryMismatches(frozen, fresh []snapshot.Entry) []string {
	frozenByPath := make(map[string]string, len(frozen))
	for _, entry := range frozen {
		frozenByPath[entry.Path] = agentRunnerEntrySignature(entry)
	}
	freshByPath := make(map[string]string, len(fresh))
	for _, entry := range fresh {
		freshByPath[entry.Path] = agentRunnerEntrySignature(entry)
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

func agentRunnerEntrySignature(entry snapshot.Entry) string {
	return fmt.Sprintf("%s|%s|%s|%d|%t|%t|%s|%s|%s",
		entry.Path, entry.Type, entry.ContentHash, entry.SizeBytes, entry.Executable, entry.ReadOnly, entry.Role, entry.Target, entry.TargetHash)
}

// agentRunnerRunRequest builds the (A′) request from the same process config and
// workspace root the dispatcher used, so the live launch is addressable.
func agentRunnerRunRequest(record adapters.AgentRunnerRequestRecord, process adapters.AgentRunnerProcessConfig, workspaceRoot string, cli AgentRunnerCLIOptions, capturer *adapters.AgentRunnerSnapshotCapturer) adapters.AgentRunnerRunRequest {
	body := record.Request
	return adapters.AgentRunnerRunRequest{
		Launch: chain.AgentRunnerLaunchRequest{
			RequestID:     body.RequestID,
			RunID:         body.RunID,
			AdapterID:     body.AdapterID,
			Command:       process.Command,
			Args:          append([]string(nil), process.Args...),
			Dir:           workspaceRoot,
			Env:           append([]string(nil), process.Env...),
			Grace:         process.Grace,
			WorkspaceRoot: workspaceRoot,
		},
		RequestHash: record.RecordHash,
		TaskID:      body.TaskID,
		StepID:      body.StepID,
		Attempt:     body.Attempt,
		Timeout:     cli.Timeout,
		Capturer:    capturer,
		MaxLogBytes: cli.MaxLogBytes,
		VerifyWait:  cli.VerifyWait,
	}
}

// agentRunnerAssemblyStep requires exactly one isolated-workspace code step. A
// multi-task, multi-step or hook definition is refused: this slice assembles one
// offline run, not a general chain.
func agentRunnerAssemblyStep(definition chain.Definition) (string, string, error) {
	if len(definition.Tasks) != 1 {
		return "", "", fmt.Errorf("%w: the AgentRunner assembly requires exactly one task", ErrInvalidConfig)
	}
	task := definition.Tasks[0]
	if len(task.Steps) != 1 {
		return "", "", fmt.Errorf("%w: the AgentRunner assembly requires exactly one step", ErrInvalidConfig)
	}
	step := task.Steps[0]
	if step.Kind != "code" || step.Mode != chain.IsolatedWorkspace {
		return "", "", fmt.Errorf("%w: the AgentRunner assembly requires one isolated-workspace code step", ErrInvalidConfig)
	}
	return task.ID, step.ID, nil
}

// agentRunnerValidateCLI fails closed on an unusable CLI pin.
func agentRunnerValidateCLI(cli AgentRunnerCLIOptions) error {
	if strings.TrimSpace(cli.Executable) == "" || strings.TrimSpace(cli.Version) == "" {
		return fmt.Errorf("%w: agent executable and version are required", ErrInvalidConfig)
	}
	if cli.Timeout <= 0 {
		return fmt.Errorf("%w: a positive run timeout is required", ErrInvalidConfig)
	}
	if cli.MaxLogBytes <= 0 {
		return fmt.Errorf("%w: a positive log bound is required", ErrInvalidConfig)
	}
	return nil
}

// agentRunnerAdmissionPolicy derives the current-evaluation-environment policy.
// The policy is measured here, never taken from the capability record: the
// platform mismatch check is only meaningful against the running environment.
func agentRunnerAdmissionPolicy(now time.Time) adapters.AgentRunnerAdmissionPolicy {
	return adapters.AgentRunnerAdmissionPolicy{
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		PlatformHash: agentRunnerPlatformHash(),
		EvaluatedAt:  now.UTC(),
		MaximumAge:   agentRunnerDefaultEvidenceAge,
	}
}

// AgentRunnerDefaultAvailabilityPolicy derives the availability policy from a
// declared availability record and the current evaluation time. It is exported
// because the CLI caller must bind the same policy the record was declared for.
func AgentRunnerDefaultAvailabilityPolicy(record adapters.AIAvailabilityRecord, now time.Time) adapters.AIAvailabilityPolicy {
	return adapters.AIAvailabilityPolicy{
		ProfileID:         record.Availability.ProfileID,
		ProfileConfigHash: record.Availability.ProfileConfigHash,
		Channel:           record.Availability.Channel,
		EvaluatedAt:       now.UTC(),
		MaximumAge:        agentRunnerDefaultEvidenceAge,
		MaximumRequests:   agentRunnerDefaultMaximumRequests,
	}
}

// agentRunnerPlatformHash digests the running platform identity. It is the
// console's own measurement of the evaluation environment.
func agentRunnerPlatformHash() string {
	return digest(agentRunnerPlatformDomain, runtime.GOOS, runtime.GOARCH)
}

func agentRunnerWorkspacePath(runRoot, taskID string, attempt int) string {
	return filepath.Join(runRoot, agentRunnerWorkspacesDirectory, fmt.Sprintf("%s-attempt-%d", taskID, attempt))
}

// agentRunnerLaunchReceiptReader and agentRunnerLaunchStopper are the private
// seam of the best-effort launch mitigation. They are not chain ports and not
// wire contracts: they exist so the mitigation path can be driven by a recording
// double in unit tests. *adapters.AgentRunnerReplayStore and
// *adapters.AgentRunnerPinnedCLILauncher satisfy them in production.
type agentRunnerLaunchReceiptReader interface {
	LaunchReceipt(requestID string) (adapters.AgentRunnerLaunchReceiptRecord, bool, error)
}

type agentRunnerLaunchStopper interface {
	StopAgentRunnerProcess(ctx context.Context, launchID string) (guard.TerminationEvidence, error)
}

// agentRunnerMitigateAbandonedLaunch is a best-effort mitigation, NOT the root
// fix. The root cause lives in internal/adapters/agent_runner_replay_dispatcher.go:
// the launch receipt is recorded before the spawn
// (agent_runner_replay_store_launches.go, RecordLaunchReceipt), so an engine
// error that is not ErrAwaitingAgentRunnerTerminal can still leave a live launch
// behind. The root-cause fix is registered as candidate OB-57 and owned by
// another slice; this helper only tries to stop the launch and never hides the
// failure.
//
// The failure semantics are unproven: a missing receipt, an unreadable receipt,
// a stop error, or a stop proof that does not read "stopped" all mean the stop
// could not be proven, and the returned error chain then says so explicitly. The
// original engine error always stays on the chain (errors.Is(runErr) remains
// true) and no new sentinel is minted.
func agentRunnerMitigateAbandonedLaunch(ctx context.Context, receipts agentRunnerLaunchReceiptReader, stopper agentRunnerLaunchStopper, requestID string, runErr error) error {
	if receipts == nil || stopper == nil {
		return fmt.Errorf("%w: mitigation stop outcome unproven: the launch reclaimer is unavailable", runErr)
	}
	receipt, found, err := receipts.LaunchReceipt(requestID)
	if err != nil {
		return fmt.Errorf("%w: mitigation stop outcome unproven: the launch receipt is unreadable: %v", runErr, err)
	}
	if !found {
		return fmt.Errorf("%w: mitigation stop outcome unproven: no launch receipt was recorded", runErr)
	}
	proof, stopErr := stopper.StopAgentRunnerProcess(ctx, receipt.Receipt.LaunchID)
	if stopErr != nil {
		return fmt.Errorf("%w: mitigation stop outcome unproven: %v", runErr, stopErr)
	}
	if proof.Outcome != "stopped" {
		return fmt.Errorf("%w: mitigation stop outcome unproven: the stop proof reads %q", runErr, proof.Outcome)
	}
	return runErr
}

// ensureAgentRunnerEventLogMarker creates the run-root marker the replay store
// requires, without truncating an existing log (a resume keeps its history).
func ensureAgentRunnerEventLogMarker(eventLog string) error {
	if _, err := os.Stat(eventLog); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(eventLog), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(eventLog, os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	return file.Close()
}
