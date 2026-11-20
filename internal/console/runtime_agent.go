package console

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/gates"
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
	agentRunnerReviewEvidenceDomain            = "proofrail:agent-runner-review-evidence:1\n"
	agentRunnerStopperDomain                   = "proofrail:agent-runner-stopper:1\n"
	agentRunnerWriterStopDomain                = "proofrail:agent-runner-writer-stop:1\n"
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

	// RequestRecord binds the single-task compatibility path: a definition with
	// exactly one task may declare its record here. It is kept because callers
	// that assemble one offline run (the CLI among them) name exactly one record.
	RequestRecord adapters.AgentRunnerRequestRecord
	// RequestRecords registers one persisted request record per task for a
	// multi-task assembly. The key is the stable encoding
	// agentRunnerRequestRecordKey(taskID, stepID) = "<taskID>/<stepID>": task and
	// step ids are validated identifiers that never contain "/", so the encoding is
	// unambiguous. When the registry is non-empty the assembly binds every task from
	// it and drives the tasks serially (CONTRACTS.md A6 "多任务串行装配（P50）");
	// RequestRecord is then ignored. Every record must carry its own requestId.
	RequestRecords map[string]adapters.AgentRunnerRequestRecord

	CapabilityRecord   adapters.AgentRunnerCapabilityRecord
	EnforcementRecord  *adapters.AgentRunnerEnforcementRecord
	AvailabilityRecord adapters.AIAvailabilityRecord
	AvailabilityPolicy adapters.AIAvailabilityPolicy

	AuthorizationLedger chain.AuthorizationLedger
	CostLedger          *tickets.CostLedger

	// AuthorizationLedgerPath and CostLedgerPath are the on-disk ledgers the
	// operator CLI appends to. When set, the dispatch rereads them just before a
	// launch is reconfirmed instead of trusting the assembly-time copies, so a
	// revocation or a settlement persisted after assembly is still observed
	// (CONTRACTS.md:212 "dispatch ... must reconfirm ... after publishing R"). A
	// zero value keeps the in-memory ledger, which is the compatibility path for
	// callers that never persisted one.
	AuthorizationLedgerPath string
	CostLedgerPath          string

	AgentCLI AgentRunnerCLIOptions

	// WorkspaceRoot overrides the derived isolated workspace. The CLI never sets
	// it; the assembly precondition tests use it to force the workspace/replay
	// overlap the replay store must refuse.
	WorkspaceRoot string

	// ProbePermit is the one-shot probe permit bound by the CLI. The assembly
	// evaluates its presence before any precondition and consumes it only after
	// every precondition passed; nil refuses fail-closed with no footprint.
	ProbePermit *gates.ProbePermit

	// SandboxEnforced requests the sandbox check-layer preflight of this
	// dispatch. The zero value (every caller that does not opt in, and every
	// production caller unless PRFRAIL_SANDBOX_CHECK=enforce is set) keeps the
	// pre-slice control flow: the check is skipped and nothing is consumed. When
	// set, the check runs before the probe permit is claimed, so a check refusal
	// never burns the one-shot quota and never creates the run-root marker
	// (CONTRACTS.md 7.2 (d)).
	SandboxEnforced bool

	Clock func() time.Time
	IDs   func() string

	// Progress is the optional per-task/per-step progress sink resolved from the
	// CLI --progress flag. A nil value means progress is off (the zero value),
	// which keeps every existing run byte-identical; a non-nil value decorates the
	// event store with the per-event progress line.
	Progress io.Writer

	// testStub is the in-package test seam of the production entry. It is
	// unexported, so no caller outside this package can reach it and every
	// production caller leaves it nil: the plan, the record validation, the
	// run-root marker, the replay store, the engine, the event log, the serial
	// driver and the inter-task stop proof are always production code. A stub
	// only supplies chain ports (and the two terminal calls the driver performs
	// directly), which lets the driver be exercised end to end on a platform
	// where the pinned CLI cannot be launched without proven publication
	// durability.
	testStub *agentRunnerAssemblyStub

	// testInMemoryLedgers marks an assembly that may answer the launch
	// reconfirmation from the in-memory AuthorizationLedger/CostLedger copies
	// instead of the durable files. T1 (finding F2): a production assembly must
	// never do that, because an in-memory copy cannot observe a revocation or a
	// settlement persisted after it was made, so the fallback is gated on this
	// flag and the flag is unexported - only the in-package test constructors set
	// it, and every production caller leaves it false.
	testInMemoryLedgers bool
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
	// stub is nil in production (the field is unexported): it only replaces
	// chain ports for in-package tests.
	stub := options.testStub

	// Step 0 (P69 T027B-E2E-WIRING): the probe permit readiness gate. Presence is
	// evaluated before any precondition so a missing permit refuses with zero
	// footprint. Quota consumption is the Require call after Step 1.
	if _, err := gates.EvaluateT027BReadiness(gates.GoalProbe, gates.ReadinessInput{Permit: options.ProbePermit}); err != nil {
		return RunSummary{}, err
	}

	// Step 1: assembly preconditions. The definition may name several tasks (each
	// one exactly one isolated-workspace code step, CONTRACTS.md A6), and every task
	// must bind its own persisted request record, so no task can borrow another
	// task's launch. A nil enforcement record is not refused here: it flows to the
	// frozen admission, which refuses a non-compatible capability without it.
	plan, err := agentRunnerAssemblyPlan(options.Definition)
	if err != nil {
		return RunSummary{}, err
	}
	records, err := agentRunnerAssemblyRecords(options, plan)
	if err != nil {
		return RunSummary{}, err
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

	// T1 (finding F2): a production assembly must be authorised by a LIVE ledger
	// source. The two on-disk ledgers are the operator CLI's append-only files, so
	// the dispatch can reread them after publishing R; an in-memory copy cannot
	// observe a revocation or a settlement that lands after assembly. The paths
	// are therefore required, and they must agree with the location the durable
	// run root implies (the shared trusted-source rule of T1/T2/T3), so a caller
	// cannot point the launch reconfirmation at an arbitrary file. Only an
	// explicitly marked test assembly keeps the in-memory compatibility path.
	if !options.testInMemoryLedgers {
		// T1 (P83 fix round 6): the refusal happens at assembly, before the run-root
		// marker, the replay store, the probe permit and the lease exist, so the
		// message states that no run state was created. A caller that has to tell
		// "refused and clean" from "refused after damage" can rely on that.
		authorizationLedgerPath, err := trustedLedgerPath(options.RunRoot, defaultApprovalsLedgerName, options.AuthorizationLedgerPath)
		if err != nil {
			return RunSummary{}, fmt.Errorf("%w: the AgentRunner assembly requires the authorization ledger: %v (refused at assembly; no run state was created)", ErrInvalidConfig, err)
		}
		options.AuthorizationLedgerPath = authorizationLedgerPath
		costLedgerPath, err := trustedLedgerPath(options.RunRoot, defaultCostLedgerName, options.CostLedgerPath)
		if err != nil {
			return RunSummary{}, fmt.Errorf("%w: the AgentRunner assembly requires the cost ledger: %v (refused at assembly; no run state was created)", ErrInvalidConfig, err)
		}
		options.CostLedgerPath = costLedgerPath
	}

	// Sandbox check layer (CONTRACTS.md 7.2): the check runs before the probe
	// permit is claimed, so an environment or configuration refusal never burns
	// the one-shot quota, never creates the run-root marker and never spawns.
	// The zero value of SandboxEnforced keeps the pre-slice control flow intact.
	if options.SandboxEnforced {
		if err := preflightSandbox(ctx, agentRunnerSandboxDispatchReason); err != nil {
			return RunSummary{}, fmt.Errorf("%w: sandbox preflight failed (fail-closed): %w", ErrInvalidConfig, err)
		}
	}

	// The probe permit quota is claimed only after every precondition passed: a
	// configuration refusal must not burn the one-shot quota. The release defers
	// Complete so the quota is spent on every path past this point, success or
	// failure (N4: a claimed permit is never refunded).
	releaseProbePermit, err := requireProbePermit(ctx, options.ProbePermit)
	if err != nil {
		return RunSummary{}, err
	}
	defer releaseProbePermit()

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

	// Step 3b (CONTRACTS.md §7 「C6 判定路径」, ADR-013): prebuild the fixed R/C
	// slots with the full declared requestId set before anything can dispatch.
	// The call site is shared by every platform: without an activated C6
	// implementation the method is a no-op that returns nil, so its error branch
	// is unreachable and Unix behaviour is unchanged. With one activated, any
	// failure is fail-closed for the whole run: no dispatch happens and a
	// partially prebuilt set stays unusable.
	if err := store.PrebuildReplaySlots(agentRunnerAssemblyRequestIDs(records)); err != nil {
		return RunSummary{}, err
	}

	// Step 4: the snapshot store and the shared capture protocol.
	snapshotRoot := filepath.Join(options.RunRoot, agentRunnerSnapshotStoreDirectory)
	snapshotStore, err := snapshot.NewStore(snapshotRoot, 0)
	if err != nil {
		return RunSummary{}, err
	}

	// Step 5: the production pinned-CLI launcher. The launch slot and the run-root
	// identity mirror are run-level (CONTRACTS.md A6): each task's launch rotates
	// them one launch at a time, which is exactly why the driver below is serial.
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

	// Step 6: the per-task ports. Every task owns its isolated workspace, snapshot
	// capturer, replay-aware dispatcher, admission, intent projection, postflight
	// and terminal publisher, all bound to that task's own request record. The
	// isolated workspace and the replay root may never overlap, so the refusal is
	// per task and the process config stays shared (it is the pinned CLI pin, not a
	// record).
	tasks, err := agentRunnerBindTasks(options, plan, records, store, launcher, snapshotStore, snapshotRoot, clock)
	if err != nil {
		return RunSummary{}, err
	}

	// Step 7: the runtime ports plus the (A′) wait carrier. The carrier shares the
	// launcher and store instances, otherwise a live launch cannot be addressed.
	// Every per-task runtime state is keyed by (taskID, attempt), so two tasks never
	// share a manifest, a candidate capture or a capturer.
	recordsByTask := make(map[string]adapters.AgentRunnerRequestRecord, len(tasks))
	workspaceRoots := make(map[string]string, len(tasks))
	capturers := make(map[string]*adapters.AgentRunnerSnapshotCapturer, len(tasks))
	admissions := make(map[string]chain.AgentRunnerAdmission, len(tasks))
	intents := make(map[string]chain.StepIntentPreparer, len(tasks))
	postflights := make(map[string]chain.PostflightPort, len(tasks))
	stepPorts := make(map[string]chain.AgentRunnerPort, len(tasks))
	for _, task := range tasks {
		taskKey := agentRunnerTaskStateKey(task.binding.TaskID, task.record.Request.Attempt)
		stepKey := agentRunnerRequestRecordKey(task.binding.TaskID, task.binding.StepID)
		recordsByTask[taskKey] = task.record
		workspaceRoots[taskKey] = task.workspaceRoot
		capturers[taskKey] = task.capturer
		admissions[stepKey] = task.admission
		intents[stepKey] = task.intent
		postflights[taskKey] = task.postflight
	}

	// The chain baseline is run-level: it is captured on the first task's workspace
	// and every later task materializes from the accepted parent the chain records.
	assemblyRuntime := &agentRunnerRuntime{
		runID:           options.RunID,
		store:           store,
		snapshotStore:   snapshotStore,
		stub:            stub,
		workspaceRoot:   tasks[0].workspaceRoot,
		requestRecord:   tasks[0].record,
		requestRecords:  recordsByTask,
		workspaceRoots:  workspaceRoots,
		capturers:       capturers,
		beforeManifests: make(map[string][]byte, len(tasks)),
		acceptedEntries: make(map[string][]snapshot.Entry, len(tasks)),
		clock:           clock,
		ids:             ids,

		authorizationLedgerPath: options.AuthorizationLedgerPath,
		authorizationLedger:     options.AuthorizationLedger,

		runRoot: options.RunRoot,
	}
	for _, task := range tasks {
		// The pre-manifest wrapper is always production code: it takes the dispatch
		// observation before any receipt can be written, which is what makes the
		// observation authoritative for the failure mitigation (H1). Only the inner
		// dispatch port may be replaced by the in-package stub.
		dispatcher := task.dispatcher
		if stub != nil && stub.Dispatcher != nil {
			dispatcher = stub.Dispatcher
		}
		stepPorts[agentRunnerRequestRecordKey(task.binding.TaskID, task.binding.StepID)] = &preManifestAgentRunnerPort{
			owner:    assemblyRuntime,
			inner:    dispatcher,
			capturer: task.capturer,
		}
	}
	assemblyRuntime.waitRun = agentRunnerDispatchWaiter(&adapters.AgentRunnerPinnedCLIRun{Launcher: launcher, Store: store, Clock: clock})
	if stub != nil && stub.Waiter != nil {
		assemblyRuntime.waitRun = stub.Waiter
	}
	admission := chain.AgentRunnerAdmission(&agentRunnerAdmissionMux{ports: admissions})
	if stub != nil && stub.Admission != nil {
		admission = stub.Admission
	}
	assemblyRuntime.steps = chain.StepRouter{
		AgentRunner:          &agentRunnerDispatchMux{ports: stepPorts},
		AgentRunnerAdmission: admission,
	}

	// Step 7b (F4): claim the exclusive publish lease of this run. The lease is
	// taken before the engine exists, so a run that cannot prove exclusive
	// publisher ownership never launches, and Publish re-proves the same record
	// against the durable store (CONTRACTS.md:150-151). A claim failure is
	// fail-closed: no engine, no dispatch, no promotion. The claim is run-scoped
	// (P83 F4 liveness), so a resume whose first record has advanced adopts the
	// lease this run already holds instead of being refused for good.
	leaseProof, err := claimPublishLease(ctx, options.RunRoot, options.RunID)
	if err != nil {
		return RunSummary{}, err
	}
	assemblyRuntime.publishLease = leaseProof

	// Step 8: the engine. Every port is the assembly runtime unless the
	// in-package stub replaces it.
	acceptance := chain.AcceptancePort(assemblyRuntime)
	reviewer := chain.ReviewerPort(assemblyRuntime)
	publisher := chain.PublisherPort(assemblyRuntime)
	postflight := chain.PostflightPort(&agentRunnerPostflightMux{ports: postflights})
	if stub != nil {
		if stub.Acceptance != nil {
			acceptance = stub.Acceptance
		}
		if stub.Reviewer != nil {
			reviewer = stub.Reviewer
		}
		if stub.Publisher != nil {
			publisher = stub.Publisher
		}
		if stub.Postflight != nil {
			postflight = stub.Postflight
		}
	}
	eventStore, err := chain.NewFileEventStore(eventLog)
	if err != nil {
		return RunSummary{}, err
	}
	// --progress only decorates the store; a nil writer keeps the plain store so
	// the run stays byte-identical to the pre-flag behaviour.
	var events chain.EventStore = eventStore
	if options.Progress != nil {
		events = newProgressStore(eventStore, options.Progress)
	}
	engine, err := chain.New(ctx, chain.Options{
		RunID:       options.RunID,
		Definition:  options.Definition,
		Events:      events,
		Baselines:   assemblyRuntime,
		Workspaces:  assemblyRuntime,
		Steps:       assemblyRuntime,
		StepIntents: &agentRunnerIntentMux{ports: intents},
		Acceptance:  acceptance,
		Reviewer:    reviewer,
		Publisher:   publisher,
		Postflight:  postflight,
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
	// abandoned-launch mitigation. The mitigation is bound to the requestId of the
	// task the driver is currently driving: with several tasks that is the only
	// launch that can have been spawned. The source-level invariant
	// TestAgentRunnerAssemblyFunnelsEveryPostSpawnFailureThroughTheFallback pins
	// that shape; a new failure path must call fail(...) and bump that test's
	// failure-return count in the same change.
	//
	// The reclaimer is the production store and launcher; the in-package stub may
	// replace them so a test can prove which requestId the closure was handed.
	receipts := agentRunnerLaunchReceiptReader(store)
	stopper := agentRunnerLaunchStopper(launcher)
	if stub != nil {
		if stub.Receipts != nil {
			receipts = stub.Receipts
		}
		if stub.Stopper != nil {
			stopper = stub.Stopper
		}
	}
	fail := func(requestID string, err error) (RunSummary, error) {
		err = agentRunnerMitigateAbandonedLaunch(ctx, receipts, stopper, requestID, err)
		return summarizeProjection(engine.Projection(), eventLog), err
	}

	// Drive: the engine parks one AgentRunner step at a time and the assembly
	// concludes it. Every task is driven in turn - locate the task the dispatch
	// port just parked, wait for that task's launch, collect its outcome, publish
	// its terminal, submit it - and the next launch may only start once the previous
	// one is proven stopped. The proof is mandatory between tasks: a run that cannot
	// prove the previous launch stopped fails closed instead of stacking launches on
	// the run-root slot and identity mirror.
	// currentRequestID names the requestId of the task the driver is driving right
	// now. The dispatch observation always wins: the pre-manifest port records it
	// before the dispatcher can write a launch receipt, so a failure that follows a
	// spawn of a task that has not become "active" yet (the engine failed inside
	// Run, before the awaiting branch updated active) still reclaims THAT task's
	// launch instead of the previous task's. The active record is only the
	// fallback for a failure that never reached a dispatch at all.
	active := tasks[0]
	currentRequestID := func() string {
		return agentRunnerDrivingRequestID(tasks, assemblyRuntime, active.record.Request.RequestID)
	}

	for {
		runErr := engine.Run(ctx)
		if runErr == nil {
			return summarizeProjection(engine.Projection(), eventLog), nil
		}
		if !errors.Is(runErr, chain.ErrAwaitingAgentRunnerTerminal) {
			// Every refusal (durability, admission, platform) lands here, but a
			// non-awaiting error can also follow a spawn: the launch receipt is written
			// before the process starts, so the shared fail closure attempts the
			// best-effort try-stop for the launch the latest dispatch observation names.
			// The original error is always preserved.
			return fail(currentRequestID(), runErr)
		}
		task, err := agentRunnerAwaitingTask(tasks, assemblyRuntime)
		if err != nil {
			return fail(currentRequestID(), err)
		}
		active = task
		runRequest := agentRunnerRunRequest(task.record, task.process, task.workspaceRoot, options.AgentCLI, task.capturer)
		outcome, err := assemblyRuntime.waitRun.WaitAgentRunnerDispatched(ctx, runRequest, assemblyRuntime.beforeManifestFor(task.binding.TaskID, task.record.Request.Attempt))
		if err != nil {
			return fail(currentRequestID(), err)
		}
		if err := assemblyRuntime.publishTerminal(ctx, task, outcome); err != nil {
			return fail(currentRequestID(), err)
		}
		terminal, err := assemblyRuntime.loadTerminal(store, task.record)
		if err != nil {
			return fail(currentRequestID(), err)
		}
		if _, err := engine.SubmitAgentRunnerTerminal(ctx, terminal); err != nil {
			return fail(currentRequestID(), err)
		}
		if outcome.Completion.Status != "completed" {
			return fail(currentRequestID(), fmt.Errorf("%w: agent runner terminal %s", chain.ErrStepFailed, outcome.Completion.Status))
		}
		// Between tasks the previous launch must be proven stopped before the next
		// dispatch may publish an R and a launch-intent on the same run root. The
		// final task has no successor, so it keeps the single-task shape exactly.
		if task.binding != tasks[len(tasks)-1].binding {
			if err := assemblyRuntime.proveTaskStopped(task.binding.TaskID, task.record.Request.Attempt); err != nil {
				return fail(currentRequestID(), err)
			}
		}
	}
}

// preManifestAgentRunnerPort captures the workspace manifest before the
// dispatcher may touch it. A failed capture stays a gap (stored as nil) and the
// dispatch still happens, exactly as the offline run degrades instead of
// aborting. The port also records which task/step the engine just parked: the
// driver resolves the awaiting task from that observation instead of assuming a
// definition order, so a skipped or resumed task cannot be waited on twice.
type preManifestAgentRunnerPort struct {
	owner    *agentRunnerRuntime
	inner    chain.AgentRunnerPort
	capturer *adapters.AgentRunnerSnapshotCapturer
}

var _ chain.AgentRunnerPort = (*preManifestAgentRunnerPort)(nil)

func (port *preManifestAgentRunnerPort) DispatchAgentRunner(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	port.owner.observeDispatch(request.TaskID, request.Step.ID, request.Attempt)
	if port.capturer != nil {
		if manifest, err := port.capturer.CaptureWorkspaceManifest(ctx, request.Workspace.Root); err == nil {
			port.owner.recordBeforeManifest(request.TaskID, request.Attempt, manifest)
		} else {
			port.owner.recordBeforeManifest(request.TaskID, request.Attempt, nil)
		}
	}
	return port.inner.DispatchAgentRunner(ctx, request)
}

// agentRunnerRuntime implements every chain port the assembly needs. It is a
// real implementation over the filesystem, the snapshot store and the production
// launcher: no port answers from memory except the pre-dispatch manifest the
// wrapper captured.
//
// One runtime serves every task of the run, so each piece of per-task state is
// keyed by (taskID, attempt): a manifest captured for one task can never be
// consumed by another, and two tasks never share a candidate capture. The
// run-level fields (run id, store, snapshot store, event ids) stay single.
type agentRunnerRuntime struct {
	runID         string
	store         *adapters.AgentRunnerReplayStore
	snapshotStore *snapshot.Store

	// requestRecord is the first (single-task compatibility: the only) request
	// record; every per-task lookup prefers requestRecords and falls back to it.
	requestRecord  adapters.AgentRunnerRequestRecord
	requestRecords map[string]adapters.AgentRunnerRequestRecord

	// workspaceRoot is the workspace the chain baseline is captured on; the
	// per-task workspaces live in workspaceRoots.
	workspaceRoot  string
	workspaceRoots map[string]string

	capturers       map[string]*adapters.AgentRunnerSnapshotCapturer
	beforeManifests map[string][]byte
	acceptedEntries map[string][]snapshot.Entry

	// reviewPolicyHashes is the policy hash every review decision bound, per task
	// attempt. F3: Publish re-reads the operator grant and refuses a promotion
	// whose freshly granted policy hash differs from the one the review bound, so
	// a grant revoked or re-issued between review and publish can never be
	// promoted on the strength of the earlier verdict.
	reviewPolicyHashes map[string]string

	// authorizationLedgerPath is the operator-written authorization ledger the
	// review port rereads. When empty the assembly-time ledger below is used, which
	// keeps every caller that never persisted a ledger on the compatibility path.
	// F3: the review port must read a source the executing agent cannot write.
	authorizationLedgerPath string
	authorizationLedger     chain.AuthorizationLedger

	// runRoot and publishLease are the F4 publish-lease state: the lease was claimed
	// before the engine existed, and Publish re-proves it against the durable store
	// instead of deriving evidence from the run id.
	runRoot      string
	publishLease publishLeaseProof

	// dispatchedTaskID/dispatchedStepID record the last step the dispatch port
	// parked. The driver is serial, so the observation is the awaiting task.
	//
	// PREMISE (pinned by I5): the observation is deliberately never reset. That is
	// sound only because every run this assembly drives is serial and offline: at
	// most one dispatch is outstanding at a time and an offline run cannot be
	// resumed into a second driver, so the last observation always names the step
	// the driver is currently handling. A future resume path or a concurrent
	// driver must reset the flag before the next engine.Run instead of reusing
	// this shape - the stale observation would then be read as the current one.
	dispatchedTaskID  string
	dispatchedStepID  string
	dispatchedAttempt int
	dispatched        bool

	// stub is the in-package test seam of the production entry; nil in production.
	// It only ever replaces chain ports and the two terminal calls the driver
	// performs directly.
	stub *agentRunnerAssemblyStub

	steps chain.StepRouter

	waitRun agentRunnerDispatchWaiter

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
// before it materializes the task workspace. The chain baseline is run-level and
// carries the first task's binding; every task after it materializes from the
// accepted parent the chain records.
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

// Materialize returns the isolated workspace bound to this task attempt. Each
// task gets its own directory, so one task's writes can never be read as
// another's candidate.
func (runtime *agentRunnerRuntime) Materialize(_ context.Context, _ string, taskID string, attempt int, _ chain.SnapshotRef) (chain.Workspace, error) {
	if err := runtime.ensureWorkspace(); err != nil {
		return chain.Workspace{}, err
	}
	root, err := runtime.workspaceRootFor(taskID, attempt)
	if err != nil {
		return chain.Workspace{}, err
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return chain.Workspace{}, fmt.Errorf("%w: create workspace: %v", ErrInvalidConfig, err)
	}
	return chain.Workspace{Root: root}, nil
}

// Execute routes one step. The AgentRunner route is admission-first inside the
// router; the step port mux selects the port bound to that task/step and the
// pre-manifest wrapper only captures before dispatch.
func (runtime *agentRunnerRuntime) Execute(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	return runtime.steps.Execute(ctx, request)
}

// Accept freezes the candidate: a real candidate capture of the workspace plus
// the digest of the recorded evidence inventory. The capture is recorded against
// (taskID, attempt) so a later task can never inherit it.
func (runtime *agentRunnerRuntime) Accept(ctx context.Context, runID, taskID string, attempt int, parent chain.SnapshotRef, workspace chain.Workspace) (chain.CandidateResult, error) {
	record, err := runtime.recordFor(taskID, attempt)
	if err != nil {
		return chain.CandidateResult{}, err
	}
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
	runtime.recordAcceptedEntries(taskID, attempt, captured.Manifest.Entries)
	evidenceRoot, err := runtime.evidenceRootHash(taskID, attempt)
	if err != nil {
		return chain.CandidateResult{}, err
	}
	return chain.CandidateResult{
		Snapshot:          chain.SnapshotRef{Hash: captured.ManifestHash},
		EvidenceRootHash:  evidenceRoot,
		Evidence:          []string{evidence.Digest(agentRunnerCandidateEvidenceDomain, []byte(captured.ManifestHash))},
		CandidateProducer: evidence.Actor{Type: "agent", ID: record.Request.AdapterID},
	}, nil
}

// Review is a real policy review over the frozen candidate. The chain only ever
// reaches it from REVIEW_PENDING, so the postflight qualification is already
// proven; the port records the policy it applied.
//
// F3: the policy hash is never derived from the candidate itself. A digest that
// the executing assembly computes proves only that the candidate is
// self-consistent, so the port instead adopts the policy hash of an operator
// grant that is still active and that covers this run/task/step, read from the
// operator-written ledger (CONTRACTS.md:147-148). No covering grant, an
// unreadable ledger, an expired/revoked grant or a grant without a valid policy
// hash refuses the review, which parks the task for repair instead of approving
// it. The decision also carries the candidate/evidence-root binding the chain
// independently checks.
func (runtime *agentRunnerRuntime) Review(_ context.Context, request chain.ReviewRequest) (chain.ReviewDecision, error) {
	if !evidence.ValidHash(request.CandidateSnapshotHash) || !evidence.ValidHash(request.EvidenceRootHash) {
		return chain.ReviewDecision{}, fmt.Errorf("%w: review request carries an invalid binding", ErrInvalidConfig)
	}
	stepID, err := runtime.stepIDForTask(request.TaskID, request.Attempt)
	if err != nil {
		return chain.ReviewDecision{}, err
	}
	policyHash, err := runtime.operatorApprovedPolicyHash(request.RunID, request.TaskID, stepID)
	if err != nil {
		return chain.ReviewDecision{}, err
	}
	// F3: the hash this decision binds is remembered, so Publish can prove the
	// policy it promotes under is still the operator-granted one.
	if runtime.reviewPolicyHashes == nil {
		runtime.reviewPolicyHashes = map[string]string{}
	}
	runtime.reviewPolicyHashes[agentRunnerTaskStateKey(request.TaskID, request.Attempt)] = policyHash
	return chain.ReviewDecision{
		ReceiptID:  runtime.nextEventID(),
		OccurredAt: runtime.clock().UTC().Format(timestampLayout),
		RecordedBy: evidence.Actor{Type: "policy", ID: agentRunnerAssemblyReviewPolicyLabel},
		ReviewMode: "policy",
		PolicyHash: policyHash,
		Outcome:    "approve",
		Evidence: []string{
			chain.ReviewEvidenceBinding(request.CandidateSnapshotHash, request.EvidenceRootHash),
			digest(agentRunnerReviewEvidenceDomain, request.CandidateSnapshotHash, request.EvidenceRootHash),
		},
	}, nil
}

// stepIDForTask resolves the step of the task the review belongs to. The request
// carries no step id, and a grant is scoped per step, so the binding is resolved
// from the assembly's own record registry rather than guessed.
func (runtime *agentRunnerRuntime) stepIDForTask(taskID string, attempt int) (string, error) {
	if record, ok := runtime.requestRecords[agentRunnerTaskStateKey(taskID, attempt)]; ok {
		if record.Request.TaskID == taskID && record.Request.StepID != "" {
			return record.Request.StepID, nil
		}
	}
	if runtime.requestRecord.Request.TaskID == taskID && runtime.requestRecord.Request.StepID != "" {
		return runtime.requestRecord.Request.StepID, nil
	}
	return "", fmt.Errorf("%w: no request record is bound to task %q", ErrInvalidConfig, taskID)
}

// operatorApprovedPolicyHash returns the policy hash of the newest active grant
// that covers the given run/task/step.
//
// F3: the ledger is reread from its file when a path is configured, because the
// grant source must be the operator CLI's append-only file and never an
// in-memory copy the executing assembly could have assembled itself. Every
// refusal is fail-closed: an unreadable ledger, no covering grant, a non-active
// grant status and a grant whose policy hash is not a valid hash all block the
// review.
//
// T2 (finding F3) is implemented partially, on purpose:
//   - DONE: only the latest grant record of an authorizationId may authorise, so
//     an older record's scope can never be backed by the newer record's status.
//   - NOT DONE: requiring grant.PolicyHash to equal the run's FROZEN policy digest.
//     This revision has no durable frozen-policy artifact for an AgentRunner run
//     (schemas/run-manifest.schema.json is declarative only: nothing in this tree
//     writes it, and the replay store freezes requests, not policies), so
//     implementing that half would invent a new record type - which the fix
//     round's hard rule forbids. It is reported as a STOP instead.
func (runtime *agentRunnerRuntime) operatorApprovedPolicyHash(runID, taskID, stepID string) (string, error) {
	ledger := runtime.authorizationLedger
	if runtime.authorizationLedgerPath != "" {
		loaded, err := loadAuthorizationLedger(runtime.authorizationLedgerPath)
		if err != nil {
			return "", fmt.Errorf("%w: authorization ledger is unreadable: %v", ErrInvalidConfig, err)
		}
		ledger = loaded
	}
	now := runtime.clock().UTC()
	policyHash := ""
	for _, record := range ledger.Grants {
		grant := record.Grant
		if grant == nil || grant.RunID != runID {
			continue
		}
		if !slices.Contains(grant.Scope.TaskIDs, taskID) || !slices.Contains(grant.Scope.StepIDs, stepID) {
			continue
		}
		// T2 (finding F3): a superseded grant record is history, not a current
		// approval. A re-issued authorization may narrow its scope, so only the
		// latest record for an authorizationId may be read as that operator's
		// intent; an older record's scope must never be mixed with the newer
		// record's status.
		latest, _, latestExists := ledger.LatestGrant(grant.AuthorizationID)
		if !latestExists || latest.RecordHash != record.RecordHash {
			continue
		}
		status, err := chain.EvaluateGrant(grant.AuthorizationID, ledger, now)
		if err != nil {
			return "", fmt.Errorf("%w: authorization grant %q cannot be evaluated: %v", ErrInvalidConfig, grant.AuthorizationID, err)
		}
		if status != chain.GrantActive {
			continue
		}
		if !evidence.ValidHash(grant.PolicyHash) {
			return "", fmt.Errorf("%w: authorization grant %q carries no valid policy hash", ErrInvalidConfig, grant.AuthorizationID)
		}
		policyHash = grant.PolicyHash
	}
	if policyHash == "" {
		return "", fmt.Errorf("%w: no operator-approved policy grant covers this review", ErrInvalidConfig)
	}
	return policyHash, nil
}

// Publish promotes the candidate only after proving the workspace still matches
// the frozen candidate and no managed process of this task is alive. The process
// proof is fail-closed, exactly like postflight leg 5: a run that left no pid
// report blocks the promotion instead of reading as "nothing alive".
func (runtime *agentRunnerRuntime) Publish(ctx context.Context, request chain.PromotionRequest) (chain.PromotionDecision, error) {
	// F3: the promotion re-proves the policy it was reviewed under and the
	// candidate evidence root it was frozen from. Both were true at review time;
	// this is the verification "before publishing" the contract requires, so a
	// revoked grant or a drifted evidence root refuses the promotion instead of
	// completing it on the earlier verdict.
	if err := runtime.verifyReviewPolicyBinding(request); err != nil {
		return runtime.blockedPromotion(err.Error()), nil
	}
	freshEvidenceRoot, err := runtime.evidenceRootHash(request.TaskID, request.Attempt)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	if freshEvidenceRoot != request.EvidenceRootHash {
		return runtime.blockedPromotion("the candidate evidence root changed after the review"), nil
	}
	workspaceRoot, err := runtime.workspaceRootFor(request.TaskID, request.Attempt)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	capturer, err := runtime.capturerFor(request.TaskID, request.Attempt)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	fresh, err := capturer.CaptureWorkspaceManifest(ctx, workspaceRoot)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	freshManifest, err := snapshot.DecodeSnapshotManifest(fresh)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	if changed := agentRunnerEntryMismatches(runtime.acceptedEntriesFor(request.TaskID, request.Attempt), freshManifest.Manifest.Entries); len(changed) > 0 {
		return chain.PromotionDecision{
			ReceiptID:     runtime.nextEventID(),
			OccurredAt:    runtime.clock().UTC().Format(timestampLayout),
			Outcome:       "failed",
			ErrorEvidence: []string{digest(agentRunnerPromotionBlockedDomain, strings.Join(changed, "\n"))},
		}, nil
	}
	stopped, err := runtime.managedProcessesGone(request.TaskID, request.Attempt)
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
	// F4: the promotion binds the lease record this run actually holds. A lease
	// that is gone, taken over, replaced or expired proves no exclusive writer, so
	// the promotion is refused with evidence instead of completing on a run-id
	// derived digest (CONTRACTS.md:150-151).
	now := runtime.clock().UTC()
	// F4 liveness finding (c): a run whose window outlived the TTL is not wedged.
	// The renewal extends only the record this run's exact token and hash match;
	// anything else is refused right here.
	// The proof here is never read from disk: it is the in-memory record obtained
	// from claimPublishLease at assembly time (Step 7b), so it cannot arrive from
	// another run's store-local file. The run-scoped identity is still re-checked
	// by both helpers below, so a future on-disk proof path stays fail-closed.
	refreshed, err := renewPublishLease(ctx, runtime.runRoot, runtime.runID, runtime.publishLease, now)
	if err != nil {
		return chain.PromotionDecision{
			ReceiptID:     runtime.nextEventID(),
			OccurredAt:    runtime.clock().UTC().Format(timestampLayout),
			Outcome:       "failed",
			ErrorEvidence: []string{digest(agentRunnerPromotionBlockedDomain, err.Error())},
		}, nil
	}
	runtime.publishLease = refreshed
	leaseHash, err := verifyPublishLease(ctx, runtime.runRoot, runtime.runID, refreshed, now)
	if err != nil {
		return chain.PromotionDecision{
			ReceiptID:     runtime.nextEventID(),
			OccurredAt:    runtime.clock().UTC().Format(timestampLayout),
			Outcome:       "failed",
			ErrorEvidence: []string{digest(agentRunnerPromotionBlockedDomain, err.Error())},
		}, nil
	}
	return chain.PromotionDecision{
		ReceiptID:            runtime.nextEventID(),
		OccurredAt:           runtime.clock().UTC().Format(timestampLayout),
		Outcome:              "completed",
		AcceptedSnapshotHash: request.CandidateSnapshotHash,
		WriterStopEvidence:   []string{digest(agentRunnerWriterStopDomain, request.RunID, request.TaskID)},
		LeaseEvidence:        []string{leaseHash},
		Evidence:             []string{digest(agentRunnerPromotionEvidenceDomain, request.ReviewReceiptHash)},
	}, nil
}

// blockedPromotion is the fail-closed promotion answer of one F3 refusal: an
// explicit failed outcome with no accepted snapshot and no lease evidence, so no
// caller can read it as a completed promotion.
func (runtime *agentRunnerRuntime) blockedPromotion(reason string) chain.PromotionDecision {
	return chain.PromotionDecision{
		ReceiptID:     runtime.nextEventID(),
		OccurredAt:    runtime.clock().UTC().Format(timestampLayout),
		Outcome:       "failed",
		ErrorEvidence: []string{digest(agentRunnerPromotionBlockedDomain, reason)},
	}
}

// verifyReviewPolicyBinding is the F3 publish-time verification of the agent
// assembly. The review rereads the operator ledger and adopts the granted policy
// hash; publishing proves the same fact once more, as close to the promotion as
// the port allows: the trusted ledger is read again and the active grant that
// covers this run/task/step must still bind the exact policy hash the review
// decision recorded. An unreadable ledger, a lost grant or a different hash all
// refuse the promotion.
func (runtime *agentRunnerRuntime) verifyReviewPolicyBinding(request chain.PromotionRequest) error {
	recorded, ok := runtime.reviewPolicyHashes[agentRunnerTaskStateKey(request.TaskID, request.Attempt)]
	if !ok || !evidence.ValidHash(recorded) {
		return fmt.Errorf("%w: no reviewed policy hash is bound to task %q", ErrInvalidConfig, request.TaskID)
	}
	stepID, err := runtime.stepIDForTask(request.TaskID, request.Attempt)
	if err != nil {
		return err
	}
	granted, err := runtime.operatorApprovedPolicyHash(request.RunID, request.TaskID, stepID)
	if err != nil {
		return err
	}
	if granted != recorded {
		return fmt.Errorf("%w: the operator grant no longer binds the reviewed policy", ErrInvalidConfig)
	}
	return nil
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

// evidenceRootHash digests the sorted artifact inventory of the task's request
// evidence directory: it is the root the review binds, not a policy decision.
func (runtime *agentRunnerRuntime) evidenceRootHash(taskID string, attempt int) (string, error) {
	root, err := runtime.replayEvidenceDir(taskID, attempt)
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

// replayEvidenceDir resolves the requestId of the task and returns its evidence
// directory. The requestId - never the task order - is what separates two
// tasks' artifacts (CONTRACTS.md A6).
func (runtime *agentRunnerRuntime) replayEvidenceDir(taskID string, attempt int) (string, error) {
	if runtime.store == nil {
		return "", fmt.Errorf("%w: replay store is not assembled", ErrInvalidConfig)
	}
	record, err := runtime.recordFor(taskID, attempt)
	if err != nil {
		return "", err
	}
	return runtime.store.RunEvidenceDir(record.Request.RequestID)
}

// agentRunnerRequestRecordKey is the stable registry key of one request record:
// task and step ids are validated identifiers that never contain "/", so
// "<taskID>/<stepID>" is unambiguous. Callers building
// AgentRunnerAssemblyOptions.RequestRecords use the same encoding.
func agentRunnerRequestRecordKey(taskID, stepID string) string {
	return taskID + "/" + stepID
}

// agentRunnerTaskStateKey is the internal key of the per-task runtime state
// (request record, workspace, capturer, manifest, accepted entries). It never
// leaves the package and is not a wire identity.
func agentRunnerTaskStateKey(taskID string, attempt int) string {
	return fmt.Sprintf("%s#%d", taskID, attempt)
}

// recordFor resolves the request record bound to one task attempt. The per-task
// registry wins; the single-record field is the compatibility path.
func (runtime *agentRunnerRuntime) recordFor(taskID string, attempt int) (adapters.AgentRunnerRequestRecord, error) {
	if record, ok := runtime.requestRecords[agentRunnerTaskStateKey(taskID, attempt)]; ok {
		return record, nil
	}
	if runtime.requestRecord.Request.TaskID == taskID && runtime.requestRecord.Request.Attempt == attempt {
		return runtime.requestRecord, nil
	}
	return adapters.AgentRunnerRequestRecord{}, fmt.Errorf("%w: no request record is bound to task %q attempt %d", ErrInvalidConfig, taskID, attempt)
}

// workspaceRootFor resolves the isolated workspace of one task attempt. The
// per-task table is authoritative; the single baseline workspace only answers
// when no table was assembled at all (the single-task compatibility path), so a
// task missing from a populated table is refused instead of silently sharing
// another task's workspace.
func (runtime *agentRunnerRuntime) workspaceRootFor(taskID string, attempt int) (string, error) {
	if root, ok := runtime.workspaceRoots[agentRunnerTaskStateKey(taskID, attempt)]; ok && strings.TrimSpace(root) != "" {
		return root, nil
	}
	if len(runtime.workspaceRoots) == 0 && strings.TrimSpace(runtime.workspaceRoot) != "" {
		return runtime.workspaceRoot, nil
	}
	return "", fmt.Errorf("%w: no workspace is bound to task %q attempt %d", ErrInvalidConfig, taskID, attempt)
}

func (runtime *agentRunnerRuntime) capturerFor(taskID string, attempt int) (*adapters.AgentRunnerSnapshotCapturer, error) {
	if capturer, ok := runtime.capturers[agentRunnerTaskStateKey(taskID, attempt)]; ok && capturer != nil {
		return capturer, nil
	}
	return nil, fmt.Errorf("%w: no snapshot capturer is bound to task %q attempt %d", ErrInvalidConfig, taskID, attempt)
}

func (runtime *agentRunnerRuntime) beforeManifestFor(taskID string, attempt int) []byte {
	return runtime.beforeManifests[agentRunnerTaskStateKey(taskID, attempt)]
}

// recordBeforeManifest stores the pre-dispatch manifest for one task attempt, so
// the wait carrier of another task can never consume it.
func (runtime *agentRunnerRuntime) recordBeforeManifest(taskID string, attempt int, manifest []byte) {
	if runtime.beforeManifests == nil {
		runtime.beforeManifests = make(map[string][]byte)
	}
	runtime.beforeManifests[agentRunnerTaskStateKey(taskID, attempt)] = manifest
}

func (runtime *agentRunnerRuntime) acceptedEntriesFor(taskID string, attempt int) []snapshot.Entry {
	return runtime.acceptedEntries[agentRunnerTaskStateKey(taskID, attempt)]
}

func (runtime *agentRunnerRuntime) recordAcceptedEntries(taskID string, attempt int, entries []snapshot.Entry) {
	if runtime.acceptedEntries == nil {
		runtime.acceptedEntries = make(map[string][]snapshot.Entry)
	}
	runtime.acceptedEntries[agentRunnerTaskStateKey(taskID, attempt)] = entries
}

// observeDispatch records the step the dispatch port just parked. The driver is
// serial, so this observation names the only task that can be awaiting.
func (runtime *agentRunnerRuntime) observeDispatch(taskID, stepID string, attempt int) {
	runtime.dispatchedTaskID = taskID
	runtime.dispatchedStepID = stepID
	runtime.dispatchedAttempt = attempt
	runtime.dispatched = true
}

// awaitingDispatch returns the last observed dispatch. It is never reset (see the
// PREMISE note on the dispatched fields): the driver is serial and offline, so
// the last observation is always the step being handled.
func (runtime *agentRunnerRuntime) awaitingDispatch() (string, string, int, bool) {
	return runtime.dispatchedTaskID, runtime.dispatchedStepID, runtime.dispatchedAttempt, runtime.dispatched
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
func (runtime *agentRunnerRuntime) managedProcessesGone(taskID string, attempt int) (bool, error) {
	root, err := runtime.replayEvidenceDir(taskID, attempt)
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

// proveTaskStopped is the serial driver's inter-task barrier: the next task's
// launch may only publish an R and a launch-intent on the same run root once the
// previous launch is proven stopped (CONTRACTS.md A6 多任务串行装配). A live
// process and an unprovable stop are both fail-closed.
func (runtime *agentRunnerRuntime) proveTaskStopped(taskID string, attempt int) error {
	stopped, err := runtime.managedProcessesGone(taskID, attempt)
	if err != nil {
		return err
	}
	if !stopped {
		return fmt.Errorf("%w: a managed process of task %q is still alive, so the next launch cannot start", ErrInvalidConfig, taskID)
	}
	return nil
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

// agentRunnerTaskBinding names one task of the assembly plan: the task and the
// single step its request record and its launch bind to.
type agentRunnerTaskBinding struct {
	TaskID string
	StepID string
}

// agentRunnerPlannedTask is one task of the assembly bound to its persisted
// request record and to the ports that record owns.
type agentRunnerPlannedTask struct {
	binding       agentRunnerTaskBinding
	record        adapters.AgentRunnerRequestRecord
	workspaceRoot string
	process       adapters.AgentRunnerProcessConfig
	capturer      *adapters.AgentRunnerSnapshotCapturer
	dispatcher    chain.AgentRunnerPort
	admission     chain.AgentRunnerAdmission
	intent        chain.StepIntentPreparer
	postflight    chain.PostflightPort
	publisher     *adapters.AgentRunnerTerminalPublisher
}

// agentRunnerAssemblyPlan returns one binding per task of the definition.
//
// A definition with no task is refused: an empty assembly would produce a chain
// that never launches anything. Multi-task definitions are accepted (CONTRACTS.md
// A6 多任务串行装配), and every task must still carry exactly one
// isolated-workspace code step, which is the only step shape the AgentRunner
// route can prove. The zero-task refusal has its own wording so it can never be
// confused with the multi-task acceptance.
func agentRunnerAssemblyPlan(definition chain.Definition) ([]agentRunnerTaskBinding, error) {
	if len(definition.Tasks) == 0 {
		return nil, fmt.Errorf("%w: the AgentRunner assembly requires at least one task", ErrInvalidConfig)
	}
	plan := make([]agentRunnerTaskBinding, 0, len(definition.Tasks))
	for _, task := range definition.Tasks {
		if len(task.Steps) != 1 {
			return nil, fmt.Errorf("%w: the AgentRunner assembly requires exactly one step", ErrInvalidConfig)
		}
		step := task.Steps[0]
		if step.Kind != "code" || step.Mode != chain.IsolatedWorkspace {
			return nil, fmt.Errorf("%w: the AgentRunner assembly requires one isolated-workspace code step", ErrInvalidConfig)
		}
		plan = append(plan, agentRunnerTaskBinding{TaskID: task.ID, StepID: step.ID})
	}
	return plan, nil
}

// agentRunnerAssemblyRecords resolves the request record of every planned task.
// The multi-record registry is authoritative when it is non-empty; otherwise the
// single compatibility record answers, and only for a single-task plan. Each
// record is validated and must bind the task/step it is registered for, and the
// requestIds must be distinct: an R/C/terminal stream is per requestId, so two
// tasks sharing one would alias their artifacts.
func agentRunnerAssemblyRecords(options AgentRunnerAssemblyOptions, plan []agentRunnerTaskBinding) ([]adapters.AgentRunnerRequestRecord, error) {
	records := make([]adapters.AgentRunnerRequestRecord, len(plan))
	if len(options.RequestRecords) > 0 {
		for index, binding := range plan {
			record, ok := options.RequestRecords[agentRunnerRequestRecordKey(binding.TaskID, binding.StepID)]
			if !ok {
				return nil, fmt.Errorf("%w: request record is missing or invalid: no record is registered for task %q step %q", ErrInvalidConfig, binding.TaskID, binding.StepID)
			}
			records[index] = record
		}
	} else if len(plan) == 1 {
		records[0] = options.RequestRecord
	} else {
		return nil, fmt.Errorf("%w: request record is missing or invalid: %d tasks need a request record registry", ErrInvalidConfig, len(plan))
	}
	bound := make(map[string]struct{}, len(records))
	for index, binding := range plan {
		record := records[index]
		if err := adapters.ValidateAgentRunnerRequestRecord(record); err != nil {
			return nil, fmt.Errorf("%w: request record is missing or invalid: %w", ErrInvalidConfig, err)
		}
		if record.Request.RunID != options.RunID || record.Request.TaskID != binding.TaskID || record.Request.StepID != binding.StepID {
			return nil, fmt.Errorf("%w: request record does not bind the assembled step", ErrInvalidConfig)
		}
		// The assembly drives attempt 1 only: the chain resolves the attempt of a
		// task that never started as 1, the isolated workspace and the inter-task
		// stop proof are keyed by (taskID, attempt), and the per-task runtime state
		// is bound to the record's attempt. A record declared for another attempt
		// would therefore be rejected later, after the run-root marker exists and
		// after the engine was built. It is refused here instead, so the refusal
		// stays a red baseline: nothing is written and no summary is projected.
		if record.Request.Attempt != 1 {
			return nil, fmt.Errorf("%w: request record of task %q declares attempt %d, but the AgentRunner assembly binds attempt 1 only", ErrInvalidConfig, binding.TaskID, record.Request.Attempt)
		}
		if _, duplicate := bound[record.Request.RequestID]; duplicate {
			return nil, fmt.Errorf("%w: request records must bind distinct request ids", ErrInvalidConfig)
		}
		bound[record.Request.RequestID] = struct{}{}
	}
	return records, nil
}

// agentRunnerAssemblyRequestIDs returns the requestId set of the validated
// assembly records in dictionary order. It is the full declared set the C6
// prebuild contract requires: never a batch and never a subset.
func agentRunnerAssemblyRequestIDs(records []adapters.AgentRunnerRequestRecord) []string {
	requestIDs := make([]string, 0, len(records))
	for _, record := range records {
		requestIDs = append(requestIDs, record.Request.RequestID)
	}
	sort.Strings(requestIDs)
	return requestIDs
}

// agentRunnerBindTasks builds the per-task ports. Everything a task can launch,
// admit, project, postflight or publish is bound to that task's own record; the
// process config is the shared pinned-CLI pin and stays one value.
func agentRunnerBindTasks(options AgentRunnerAssemblyOptions, plan []agentRunnerTaskBinding, records []adapters.AgentRunnerRequestRecord, store *adapters.AgentRunnerReplayStore, launcher *adapters.AgentRunnerPinnedCLILauncher, snapshotStore *snapshot.Store, snapshotRoot string, clock func() time.Time) ([]agentRunnerPlannedTask, error) {
	process := adapters.AgentRunnerProcessConfig{
		Command: options.AgentCLI.Executable,
		Args:    append([]string(nil), options.AgentCLI.RunArgs...),
		Env:     append([]string(nil), options.AgentCLI.Env...),
		Grace:   options.AgentCLI.Grace,
	}
	workspaceOverride := strings.TrimSpace(options.WorkspaceRoot)
	// T1 (finding F2 / P83 fix round 6): ONE derivation of the trusted ledger
	// sources, shared by the pre-admission and the dispatch reconfirmation. Both
	// checks therefore read the same durable files, so they cannot disagree about
	// what the operator currently authorises or what budget is still available.
	ledgerSource := func() (chain.AuthorizationLedger, error) {
		if options.AuthorizationLedgerPath != "" {
			return loadAuthorizationLedger(options.AuthorizationLedgerPath)
		}
		if options.testInMemoryLedgers {
			return options.AuthorizationLedger, nil
		}
		return chain.AuthorizationLedger{}, fmt.Errorf("%w: no authorization ledger source is configured", ErrInvalidConfig)
	}
	costSource := func() (*tickets.CostLedger, error) {
		if options.CostLedgerPath != "" {
			return loadCostLedgerSnapshot(options.CostLedgerPath)
		}
		if options.testInMemoryLedgers {
			return options.CostLedger, nil
		}
		return nil, fmt.Errorf("%w: no cost ledger source is configured", ErrInvalidConfig)
	}
	tasks := make([]agentRunnerPlannedTask, 0, len(plan))
	for index, binding := range plan {
		record := records[index]
		workspaceRoot := agentRunnerWorkspacePath(options.RunRoot, binding.TaskID, record.Request.Attempt)
		// The override exists for the single-run precondition tests, which force the
		// workspace/replay overlap the replay store must refuse. Two tasks sharing one
		// workspace would break the per-task candidate capture, so it is ignored for
		// a multi-task plan.
		if workspaceOverride != "" && len(plan) == 1 {
			workspaceRoot = workspaceOverride
		}
		// Step 4 of the old single-task path: the isolated workspace and the replay
		// root may never overlap.
		if err := store.RejectWriterRoots(workspaceRoot); err != nil {
			return nil, err
		}
		capturer, err := adapters.NewAgentRunnerSnapshotCapturer(snapshotStore, options.RunID, binding.TaskID, record.Request.Attempt)
		if err != nil {
			return nil, err
		}
		intent, err := adapters.NewAgentRunnerStepIntentPreparer(record, options.Definition)
		if err != nil {
			return nil, err
		}
		postflight, err := adapters.NewAgentRunnerPostflight(store, record, snapshotRoot)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, agentRunnerPlannedTask{
			binding:       binding,
			record:        record,
			workspaceRoot: workspaceRoot,
			process:       process,
			capturer:      capturer,
			dispatcher: &adapters.AgentRunnerReplayDispatcher{
				Store:    store,
				Launcher: launcher,
				Process:  process,
				// T1 (finding F2): the authorization ledger is reread from its trusted
				// file at every dispatch, so a revocation appended after assembly is
				// visible before the launch is reconfirmed (CONTRACTS.md:212). The
				// in-memory copy is no longer a fallback: it is reachable only from an
				// explicitly marked test assembly, and a production closure without a
				// path refuses instead of approving on a snapshot it cannot refresh.
				LedgerSnapshot: ledgerSource,
				// T1 (finding F2): the cost ledger is reread from its durable file so a
				// settlement persisted in another process removes the outstanding
				// reservation before the launch is permitted (CONTRACTS.md:212). The
				// in-memory copy is limited to the same marked test assemblies.
				CostLedgerSnapshot: costSource,
				CostLedger:         options.CostLedger,
				RequestRecord:      record,
				Clock:              clock,
			},
			admission: &adapters.AgentRunnerCompositeAdmission{
				RequestRecord: record,
				// T1 (P83 fix round 6): the pre-admission reads the same derived trusted
				// sources as the dispatch reconfirmation, so a revocation or settlement
				// that landed after assembly refuses the step before it is admitted
				// instead of being judged against the assembly-time snapshot.
				AuthorizationLedger:       options.AuthorizationLedger,
				AuthorizationLedgerSource: ledgerSource,
				CostLedger:                options.CostLedger,
				CostLedgerSource:          costSource,
				AvailabilityRecord:        options.AvailabilityRecord,
				AvailabilityPolicy:        options.AvailabilityPolicy,
				CapabilityRecord:          options.CapabilityRecord,
				EnforcementRecord:         options.EnforcementRecord,
				AdmissionPolicy:           agentRunnerAdmissionPolicy(clock()),
				Clock:                     clock,
			},
			intent:     intent,
			postflight: postflight,
			publisher: &adapters.AgentRunnerTerminalPublisher{
				Store:         store,
				CostLedger:    options.CostLedger,
				RequestRecord: record,
				Clock:         clock,
			},
		})
	}
	return tasks, nil
}

// agentRunnerAwaitingTask resolves the task the engine just parked from the
// dispatch the pre-manifest port observed. The observation is taken the moment
// the port is entered, so it names the only task that can be awaiting; a plan
// without a matching binding is fail-closed instead of guessed.
func agentRunnerAwaitingTask(tasks []agentRunnerPlannedTask, runtime *agentRunnerRuntime) (agentRunnerPlannedTask, error) {
	taskID, stepID, _, observed := runtime.awaitingDispatch()
	if !observed {
		return agentRunnerPlannedTask{}, fmt.Errorf("%w: the AgentRunner assembly cannot bind the awaiting step to a dispatch", ErrInvalidConfig)
	}
	for _, task := range tasks {
		if task.binding.TaskID == taskID && task.binding.StepID == stepID {
			return task, nil
		}
	}
	return agentRunnerPlannedTask{}, fmt.Errorf("%w: dispatched step %q of task %q has no request record", ErrInvalidConfig, stepID, taskID)
}

// agentRunnerDispatchMux routes a dispatch to the port bound to that task/step.
// The assembly is serial by contract, so the mux is a lookup, never a scheduler:
// it is the only reason two tasks can share one chain router.
type agentRunnerDispatchMux struct {
	ports map[string]chain.AgentRunnerPort
}

var _ chain.AgentRunnerPort = (*agentRunnerDispatchMux)(nil)

func (mux *agentRunnerDispatchMux) DispatchAgentRunner(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	port, ok := mux.ports[agentRunnerRequestRecordKey(request.TaskID, request.Step.ID)]
	if !ok || port == nil {
		return chain.StepResult{}, fmt.Errorf("%w: no AgentRunner launch port is bound to task %q step %q", ErrInvalidConfig, request.TaskID, request.Step.ID)
	}
	return port.DispatchAgentRunner(ctx, request)
}

// agentRunnerAdmissionMux routes the runtime preflight to the admission bound to
// that task/step, so each task is admitted against its own record.
type agentRunnerAdmissionMux struct {
	ports map[string]chain.AgentRunnerAdmission
}

var _ chain.AgentRunnerAdmission = (*agentRunnerAdmissionMux)(nil)

func (mux *agentRunnerAdmissionMux) AdmitAgentRunner(ctx context.Context, request chain.StepRequest) error {
	port, ok := mux.ports[agentRunnerRequestRecordKey(request.TaskID, request.Step.ID)]
	if !ok || port == nil {
		return fmt.Errorf("%w: no AgentRunner admission is bound to task %q step %q", ErrInvalidConfig, request.TaskID, request.Step.ID)
	}
	return port.AdmitAgentRunner(ctx, request)
}

// agentRunnerIntentMux routes the engine-owned intent projection to the preparer
// bound to that task/step. The adapters preparer answers a non-matching step with
// a zero intent, so the selection must happen here rather than inside it.
type agentRunnerIntentMux struct {
	ports map[string]chain.StepIntentPreparer
}

var _ chain.StepIntentPreparer = (*agentRunnerIntentMux)(nil)

func (mux *agentRunnerIntentMux) PrepareStepIntent(ctx context.Context, request chain.StepRequest) (chain.StepExecutionIntent, error) {
	port, ok := mux.ports[agentRunnerRequestRecordKey(request.TaskID, request.Step.ID)]
	if !ok || port == nil {
		return chain.StepExecutionIntent{}, fmt.Errorf("%w: no step intent is bound to task %q step %q", ErrInvalidConfig, request.TaskID, request.Step.ID)
	}
	return port.PrepareStepIntent(ctx, request)
}

// agentRunnerPostflightMux routes the chain-owned postflight to the port bound to
// that task attempt: the postflight reads the requestId-scoped evidence of its
// own task.
type agentRunnerPostflightMux struct {
	ports map[string]chain.PostflightPort
}

var _ chain.PostflightPort = (*agentRunnerPostflightMux)(nil)

func (mux *agentRunnerPostflightMux) RunPostflight(ctx context.Context, request chain.PostflightRequest) (chain.PostflightDecision, error) {
	port, ok := mux.ports[agentRunnerTaskStateKey(request.TaskID, request.Attempt)]
	if !ok || port == nil {
		return chain.PostflightDecision{}, fmt.Errorf("%w: no postflight is bound to task %q attempt %d", ErrInvalidConfig, request.TaskID, request.Attempt)
	}
	return port.RunPostflight(ctx, request)
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

// agentRunnerDispatchWaiter is the (A') wait carrier contract. It is not a chain
// port and not a wire contract: *adapters.AgentRunnerPinnedCLIRun satisfies it in
// production, and the in-package stub satisfies it in tests that cannot spawn the
// pinned CLI.
type agentRunnerDispatchWaiter interface {
	WaitAgentRunnerDispatched(ctx context.Context, request adapters.AgentRunnerRunRequest, beforeManifest []byte) (adapters.AgentRunnerTerminalOutcome, error)
}

// agentRunnerAssemblyStub is the in-package port double of the production
// driver. It exists only on AgentRunnerAssemblyOptions.testStub, so production
// code paths can never reach it. Every field is optional: a nil field keeps the
// production wiring for that port, so a test replaces exactly what it must and
// the untouched ports stay real.
type agentRunnerAssemblyStub struct {
	// Dispatcher replaces the inner AgentRunner port of every task. The
	// pre-manifest wrapper (and with it the dispatch observation) stays production
	// code.
	Dispatcher chain.AgentRunnerPort
	// Admission replaces the runtime admission router.
	Admission chain.AgentRunnerAdmission
	// Acceptance, Reviewer, Publisher and Postflight replace the run-level chain
	// ports.
	Acceptance chain.AcceptancePort
	Reviewer   chain.ReviewerPort
	Publisher  chain.PublisherPort
	Postflight chain.PostflightPort
	// Waiter replaces the (A') wait carrier.
	Waiter agentRunnerDispatchWaiter
	// PublishTerminal and LoadTerminal replace the two terminal calls the driver
	// performs directly on the planned task.
	PublishTerminal func(ctx context.Context, outcome adapters.AgentRunnerTerminalOutcome) error
	LoadTerminal    func(record adapters.AgentRunnerRequestRecord) (chain.AgentRunnerTerminal, error)
	// Receipts and Stopper replace the best-effort launch reclaimer of the fail
	// closure.
	Receipts agentRunnerLaunchReceiptReader
	Stopper  agentRunnerLaunchStopper
}

// agentRunnerDrivingRequestID resolves the requestId of the task the driver is
// currently handling. The dispatch observation wins over the active record
// because the observation is taken before the dispatcher can write a launch
// receipt: a failure raised inside engine.Run after a spawn of a task that never
// became active therefore still reclaims that task's launch. fallback is used
// only when no dispatch was observed at all.
func agentRunnerDrivingRequestID(tasks []agentRunnerPlannedTask, runtime *agentRunnerRuntime, fallback string) string {
	taskID, stepID, _, observed := runtime.awaitingDispatch()
	if observed {
		for _, task := range tasks {
			if task.binding.TaskID == taskID && task.binding.StepID == stepID {
				return task.record.Request.RequestID
			}
		}
	}
	return fallback
}

// publishTerminal publishes the terminal of the task the driver is concluding.
// Production binds the adapters publisher of that task; the in-package stub may
// replace the call, which is what lets the serial driver be exercised without a
// live launch.
func (runtime *agentRunnerRuntime) publishTerminal(ctx context.Context, task agentRunnerPlannedTask, outcome adapters.AgentRunnerTerminalOutcome) error {
	if runtime.stub != nil && runtime.stub.PublishTerminal != nil {
		return runtime.stub.PublishTerminal(ctx, outcome)
	}
	_, err := task.publisher.PublishTerminal(ctx, outcome)
	return err
}

// loadTerminal loads the chain terminal of one task from the replay store.
// Production reads the record the publisher just wrote; the in-package stub may
// answer from the same outcome it scripted.
func (runtime *agentRunnerRuntime) loadTerminal(store *adapters.AgentRunnerReplayStore, record adapters.AgentRunnerRequestRecord) (chain.AgentRunnerTerminal, error) {
	if runtime.stub != nil && runtime.stub.LoadTerminal != nil {
		return runtime.stub.LoadTerminal(record)
	}
	return adapters.LoadChainAgentRunnerTerminal(store, record)
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

// loadCostLedgerSnapshot rereads the durable cost ledger file the operator CLI
// writes. It exists for F2: the launch reconfirmation must consult the ledger as
// it is now, not the copy the assembly loaded before any dispatch happened
// (CONTRACTS.md:212). Every read or decode failure is returned as an error so the
// caller can refuse the launch; there is no "ledger unreadable therefore empty"
// fallback, because an empty ledger would silently mean "no outstanding
// reservation" and block by accident rather than by evidence.
func loadCostLedgerSnapshot(path string) (*tickets.CostLedger, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	record, err := tickets.DecodeCostLedgerRecord(data)
	if err != nil {
		return nil, err
	}
	return tickets.LoadCostLedger(record)
}

// Publish lease (F4).
//
// CONTRACTS.md:150-151 binds the shutdown and the lease evidence of a promotion
// into one transaction. Deriving the evidence from the run id proves nothing
// about who currently holds the writer, so both assemblies instead claim a real
// guard.LeaseStore lease before the chain can reach Publish and bind the current
// record's hash at publish time. A missing, replaced, taken-over or expired lease
// refuses the promotion.
const (
	publishLeaseTTL            = 24 * time.Hour
	publishLeaseConsumerID     = "proofrail-publisher"
	publishLeaseResourceDomain = "proofrail:publish-lease-resource:1\n"
	// publishLeaseClaimDomain derives the run-scoped claim identity. It is a
	// digest of the run id and never of an attempt, so every re-entry of the same
	// run addresses the same claim (P83 F4 liveness finding (a)/(b)).
	publishLeaseClaimDomain = "proofrail:publish-lease-claim:1\n"
	// publishLeaseRenewBefore is the remaining lifetime below which a lease that
	// this run still provably holds is renewed before it is verified. A run whose
	// publish window legitimately stays open for longer than publishLeaseTTL must
	// not be wedged by expiry (P83 F4 liveness finding (c)).
	publishLeaseRenewBefore = publishLeaseTTL / 2
	// publishLeaseProofName is the store-local file a re-entrant supervised run
	// reads its claim back from. It is not a wire record.
	publishLeaseProofName = "publish-lease.json"
	publishLeaseTimestamp = "2006-01-02T15:04:05.000Z"
)

// ErrPublishLeaseUnproven reports a promotion whose exclusive lease cannot be
// proven. It is never a warning: the promotion is refused.
var ErrPublishLeaseUnproven = errors.New("publish lease proof missing or inconsistent")

// publishLeaseProof is the claim a run persists so a later process can prove it
// still holds the same lease.
type publishLeaseProof struct {
	ResourceID string             `json:"resourceId"`
	Token      guard.FencingToken `json:"token"`
	LeaseHash  string             `json:"leaseHash"`
}

// publishLeaseRoot is the lease directory of a run root.
func publishLeaseRoot(runRoot string) string { return filepath.Join(runRoot, "leases") }

// publishLeaseResourceID derives the lease resource of a run. The resource id is
// a run-scoped identifier, never a task or step, because the publication it
// fences is the run's terminal promotion.
func publishLeaseResourceID(runID string) string {
	digest := evidence.Digest(publishLeaseResourceDomain, []byte(runID))
	return "ws-" + digest[len("sha256:"):len("sha256:")+24]
}

// newPublishLeaseStore opens the lease store under the run root.
func newPublishLeaseStore(runRoot string) (*guard.LeaseStore, error) {
	store, err := guard.NewLeaseStore(publishLeaseRoot(runRoot))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPublishLeaseUnproven, err)
	}
	return store, nil
}

// claimPublishLease claims the run's publish lease and reports the proof to bind
// at publish time. The claim identity is run-scoped, never per attempt (P83 F4
// liveness finding (a)/(b)): a supervised run whose task reaches attempt 2 and an
// agent-runner resume whose first record has advanced both re-enter the same run,
// so they must adopt the lease the run already holds instead of being refused for
// the rest of the run.
//
// Adoption accepts only a lease this run provably wrote: the same run-scoped
// resource and the same consumer. It deliberately does not compare the claim id,
// because older records carry a per-attempt id and every renewal keeps the id but
// advances the generation. The relaxation cannot admit a foreign owner: the
// resource id is derived from the run id, so another run addresses a different
// resource. A deleted, replaced or hash-swapped record is still refused here (the
// store read fails) and verifyPublishLease re-binds the exact (claim id,
// generation) token plus lease hash before a promotion may complete.
//
// T5 (P83 fix round): adoption is additionally gated on process-level exclusivity.
// A second PID of the same run may only adopt the record after the recorded holder
// is proven stopped, and a holder whose liveness cannot be disproved refuses
// adoption. Without that proof two processes of one run would both believe they
// hold the writer.
func claimPublishLease(ctx context.Context, runRoot, runID string) (publishLeaseProof, error) {
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		return publishLeaseProof{}, err
	}
	resourceID := publishLeaseResourceID(runID)
	holder := publishLeaseProcess(runID)
	lease, err := store.Claim(ctx, resourceID, publishLeaseConsumerID, publishLeaseClaimID(runID), holder, publishLeaseTTL)
	if err != nil {
		if !errors.Is(err, guard.ErrLeaseHeld) {
			return publishLeaseProof{}, fmt.Errorf("%w: claim %s: %v", ErrPublishLeaseUnproven, resourceID, err)
		}
		current, currentErr := store.Current(ctx, resourceID)
		if currentErr != nil {
			return publishLeaseProof{}, fmt.Errorf("%w: claim %s: %v", ErrPublishLeaseUnproven, resourceID, currentErr)
		}
		if !publishLeaseBelongsToRun(current, runID) {
			return publishLeaseProof{}, fmt.Errorf("%w: the publish lease of %s is held by another holder", ErrPublishLeaseUnproven, resourceID)
		}
		// T5 (P83 fix round 6): adoption is a TRANSFER, not a hand-out of the old
		// record. Returning the holder's token let two successors of one stopped
		// holder both adopt it and both believe they held the writer. The record is
		// reclaimed instead: it names this process and a new generation, so the
		// second successor presents the previous token and is refused by the store.
		//
		// A record that already names THIS process is not a successor at all - it is
		// a re-entry of the one holder (another task of the same supervised run, or a
		// repeated assembly inside one process), and no OS-level liveness proof could
		// tell two call sites of one process apart. That case keeps the record, so
		// the proofs already written for this run's other tasks stay valid instead of
		// being invalidated by a transfer that separates no holders.
		if current.Process == holder {
			return publishLeaseProof{ResourceID: resourceID, Token: current.Token, LeaseHash: current.LeaseHash}, nil
		}
		reclaimed, reclaimErr := store.Reclaim(ctx, resourceID, current.Token, holder, publishLeaseTTL, func(locked guard.Lease) error {
			// The proof runs under the store lock, so the "previous holder is stopped"
			// decision cannot be separated from the transfer by a competing successor.
			return publishLeasePreviousHolderStopped(locked.Process)
		})
		if reclaimErr != nil {
			return publishLeaseProof{}, fmt.Errorf("%w: adopt %s: %v", ErrPublishLeaseUnproven, resourceID, reclaimErr)
		}
		lease = reclaimed
	}
	return publishLeaseProof{ResourceID: resourceID, Token: lease.Token, LeaseHash: lease.LeaseHash}, nil
}

// publishLeaseProcess is the process identity a claim records. T5: it is the real
// OS identity this repo already inspects (guard.InspectProcess), so a later process
// can ask the OS whether the recorded holder is still running instead of trusting a
// token that only names the run. A platform that cannot report an identity keeps
// the pre-T5 run-token shape, which no other process of the run can distinguish
// from itself; on those platforms the adoption proof below therefore refuses.
func publishLeaseProcess(runID string) guard.ProcessIdentity {
	if identity, err := guard.InspectProcess(os.Getpid()); err == nil {
		return identity
	}
	return guard.ProcessIdentity{PID: os.Getpid(), StartToken: runID}
}

// publishLeasePreviousHolderStopped is the T5 exclusivity proof of adoption. It
// admits the record only when the recorded holder is this very process (a
// re-entry inside one process) or when the OS proves that identity is gone: the
// identity is not alive AND no process is inspectable at that PID any more. A live
// holder, a PID that now carries a different process, and an unobservable identity
// all refuse adoption - the previous holder's stop cannot be disproved for them.
func publishLeasePreviousHolderStopped(previous guard.ProcessIdentity) error {
	if self, err := guard.InspectProcess(os.Getpid()); err == nil && previous == self {
		return nil
	}
	alive, aliveErr := guard.ProcessAlive(previous)
	if aliveErr != nil {
		return fmt.Errorf("the previous publish lease holder %d cannot be inspected: %v", previous.PID, aliveErr)
	}
	if alive {
		return fmt.Errorf("the previous publish lease holder %d is still running", previous.PID)
	}
	if occupant, inspectErr := guard.InspectProcess(previous.PID); inspectErr == nil {
		return fmt.Errorf("the previous publish lease holder cannot be proven stopped: pid %d is now held by process %d", previous.PID, occupant.PID)
	}
	return nil
}

// publishLeaseClaimID is the run-scoped claim identity of a lease record. The
// same run always derives the same id, so claim and adoption agree across
// attempts, processes and resumes.
func publishLeaseClaimID(runID string) string {
	digest := evidence.Digest(publishLeaseClaimDomain, []byte(runID))
	return "run-" + digest[len("sha256:"):len("sha256:")+24]
}

// publishLeaseBelongsToRun reports whether a lease record is this run's own: the
// run-scoped resource and the publish consumer. T5: the recorded process identity
// is no longer compared with the run id, because it now names the real OS process
// that wrote the record; whether that holder stopped is a separate proof
// (publishLeasePreviousHolderStopped), and the resource id - a digest of the run
// id - still keeps another run's record out.
func publishLeaseBelongsToRun(lease guard.Lease, runID string) bool {
	return lease.ResourceID == publishLeaseResourceID(runID) &&
		lease.ConsumerID == publishLeaseConsumerID
}

// publishLeaseProofBelongsToRun is the F4 identity half of the publish proof: a
// proof binds a record only when it names this run's own lease resource. The
// resource id is a digest of the run id, so a proof file replaced with another
// run's legitimately issued lease (same store root) fails here even though its
// token and lease hash would match that other run's record.
func publishLeaseProofBelongsToRun(proof publishLeaseProof, runID string) bool {
	return proof.ResourceID != "" && proof.ResourceID == publishLeaseResourceID(runID)
}

// renewPublishLease re-proves that this run still owns the lease record and
// extends it when less than publishLeaseRenewBefore remains (P83 F4 liveness
// finding (c)): a run whose window stays open beyond the TTL would otherwise be
// refused at promotion forever. The renewal is not a weakening of the proof
// because it runs only after the record was proven to be this run's own record
// (run-scoped resource, F4) and matched against this run's exact token and lease
// hash, and the store renews that exact token again under the resource lock. A
// foreign writer, a swapped record or a deleted record fails before any renewal
// is attempted.
//
// It returns the proof to use from now on; it is unchanged when no renewal was
// due, so callers persist it only when it actually moved.
func renewPublishLease(ctx context.Context, runRoot, runID string, proof publishLeaseProof, now time.Time) (publishLeaseProof, error) {
	if proof.Token.ClaimID == "" || proof.Token.Generation == 0 || !evidence.ValidHash(proof.LeaseHash) {
		return publishLeaseProof{}, ErrPublishLeaseUnproven
	}
	// F4: the proof must bind THIS run's own record. Without this check a proof
	// file swapped for another run's legitimate lease of the same store root would
	// be renewed and published although this run never held it.
	if !publishLeaseProofBelongsToRun(proof, runID) {
		return publishLeaseProof{}, fmt.Errorf("%w: the proof names lease %q, not this run's %q", ErrPublishLeaseUnproven, proof.ResourceID, publishLeaseResourceID(runID))
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		return publishLeaseProof{}, err
	}
	lease, err := store.Current(ctx, proof.ResourceID)
	if err != nil {
		return publishLeaseProof{}, fmt.Errorf("%w: %v", ErrPublishLeaseUnproven, err)
	}
	holder := publishLeaseProcess(runID)
	if lease.Token != proof.Token || lease.LeaseHash != proof.LeaseHash {
		// T5#4 (P83 fix round 7): a generation commit is followed by a best-effort
		// rewrite of every task's store-local proof, so an interruption between the two
		// leaves a task on the superseded token although the record still names THIS
		// process as its holder. Refusing there was a dead end: re-entry presented the
		// same stale proof forever. A record that names this process IS this process's
		// authority, so the proof is re-synced to it - which is exactly the ownership
		// condition verifyPublishLease demands before a promotion may complete.
		//
		// T5#4 (P83 fix round 8, item C): "the record names this process" alone is not
		// enough to accept any token a proof file carries - a proof swapped for another
		// run's token of the same store root would be re-synced just as happily. The
		// stale proof is therefore reconciled only when it is a verifiable EARLIER
		// GENERATION OF THE SAME CLAIM: same lease resource (checked above), same claim
		// id, exactly one generation behind, the current record naming exactly that hash
		// as its predecessor, and the retained body for that hash reading back with the
		// proof's own token. Anything else keeps the refusal.
		if lease.Process != holder {
			return publishLeaseProof{}, fmt.Errorf("%w: the current lease of %s is not the one this run claimed", ErrPublishLeaseUnproven, proof.ResourceID)
		}
		if err := verifyPublishLeasePredecessor(store, proof, lease); err != nil {
			return publishLeaseProof{}, fmt.Errorf("%w: %v", ErrPublishLeaseUnproven, err)
		}
		proof = publishLeaseProof{ResourceID: proof.ResourceID, Token: lease.Token, LeaseHash: lease.LeaseHash}
	}
	// T5 (P83 fix round 6): matching the token proves only that this run once saw
	// that token; it does not prove who holds the record now. Any process able to
	// read the store-local proof could otherwise extend a lease it never claimed -
	// including a second successor of a stopped holder, which would then verify the
	// same record and publish. The ownership proof below is therefore the same one
	// adoption uses, and it runs before any expiry decision: a holder that is not
	// this process has to be provably stopped and is then reclaimed, which is what
	// makes the second successor fail.
	if lease.Process != holder {
		reclaimed, reclaimErr := store.Reclaim(ctx, proof.ResourceID, proof.Token, holder, publishLeaseTTL, func(locked guard.Lease) error {
			return publishLeasePreviousHolderStopped(locked.Process)
		})
		if reclaimErr != nil {
			return publishLeaseProof{}, fmt.Errorf("%w: take the publish lease of %s: %v", ErrPublishLeaseUnproven, proof.ResourceID, reclaimErr)
		}
		return publishLeaseProof{ResourceID: proof.ResourceID, Token: reclaimed.Token, LeaseHash: reclaimed.LeaseHash}, nil
	}
	expiresAt, err := time.Parse(publishLeaseTimestamp, lease.LeaseExpiresAt)
	if err != nil {
		return publishLeaseProof{}, fmt.Errorf("%w: lease expiry is unreadable: %v", ErrPublishLeaseUnproven, err)
	}
	if expiresAt.Sub(now) > publishLeaseRenewBefore {
		return proof, nil
	}
	renewed, err := store.Renew(ctx, proof.ResourceID, proof.Token, holder, publishLeaseTTL)
	if err != nil {
		return publishLeaseProof{}, fmt.Errorf("%w: renew %s: %v", ErrPublishLeaseUnproven, proof.ResourceID, err)
	}
	return publishLeaseProof{ResourceID: proof.ResourceID, Token: renewed.Token, LeaseHash: renewed.LeaseHash}, nil
}

// verifyPublishLeasePredecessor proves that a stale store-local proof is a genuine
// earlier generation of the SAME claim, so a generation commit that was not followed
// by the proof rewrite stays recoverable without accepting anything else.
//
// T5#4 (P83 fix round 8, item C): the generation commit (Renew/Reclaim) and the
// rewrite of the task proofs are two separate writes, and a process that dies
// between them leaves a task on the superseded token. Refusing there would strand
// the task forever, so the next entry reconciles - but only through the retained
// records: the proof must name the same claim, exactly one generation behind, the
// hash the current record itself records as its predecessor, and a retained body
// that verifies and carries the proof's token. A proof carrying a token this store
// never issued for this resource (another run's token, a forged generation, a
// fabricated hash) fails one of those checks and is refused.
func verifyPublishLeasePredecessor(store *guard.LeaseStore, proof publishLeaseProof, current guard.Lease) error {
	if proof.Token.ClaimID != current.Token.ClaimID {
		return fmt.Errorf("the stale proof names claim %q, not the current claim %q of %s", proof.Token.ClaimID, current.Token.ClaimID, current.ResourceID)
	}
	if proof.Token.Generation+1 != current.Token.Generation {
		return fmt.Errorf("the stale proof names generation %d, which is not the generation before %d of %s", proof.Token.Generation, current.Token.Generation, current.ResourceID)
	}
	if current.PreviousHash == nil || *current.PreviousHash != proof.LeaseHash {
		return fmt.Errorf("the current lease of %s does not name %s as its predecessor", current.ResourceID, proof.LeaseHash)
	}
	retained, err := store.RetainedLease(proof.ResourceID, proof.LeaseHash)
	if err != nil {
		return fmt.Errorf("the predecessor of %s cannot be verified from the retained records: %v", current.ResourceID, err)
	}
	if retained.Token != proof.Token {
		return fmt.Errorf("the retained predecessor %s carries token %s/%d, not the proof's %s/%d", proof.LeaseHash, retained.Token.ClaimID, retained.Token.Generation, proof.Token.ClaimID, proof.Token.Generation)
	}
	return nil
}

// verifyPublishLease re-reads the lease and proves it is still the one this run
// claimed. It returns the current lease hash, which is what the promotion binds.
func verifyPublishLease(ctx context.Context, runRoot, runID string, proof publishLeaseProof, now time.Time) (string, error) {
	if proof.Token.ClaimID == "" || proof.Token.Generation == 0 || !evidence.ValidHash(proof.LeaseHash) {
		return "", ErrPublishLeaseUnproven
	}
	// F4: a proof only binds this run's own resource. A repaired proof naming any
	// other run's lease is refused, so publication cannot bind a foreign hash.
	if !publishLeaseProofBelongsToRun(proof, runID) {
		return "", fmt.Errorf("%w: the proof names lease %q, not this run's %q", ErrPublishLeaseUnproven, proof.ResourceID, publishLeaseResourceID(runID))
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		return "", err
	}
	lease, err := store.Current(ctx, proof.ResourceID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrPublishLeaseUnproven, err)
	}
	if lease.Token != proof.Token || lease.LeaseHash != proof.LeaseHash {
		return "", fmt.Errorf("%w: the current lease of %s is not the one this run claimed", ErrPublishLeaseUnproven, proof.ResourceID)
	}
	// T5 (P83 fix round 6): the proof file is a claim, not an ownership proof. The
	// record has to name THIS process as the holder, so a proof that was read from
	// disk by a process which never held the lease (a second successor of a stopped
	// holder, or any reader of the run root) cannot verify it and publish.
	if holder := publishLeaseProcess(runID); lease.Process != holder {
		return "", fmt.Errorf("%w: the publish lease of %s names process %d as its holder, not this process", ErrPublishLeaseUnproven, proof.ResourceID, lease.Process.PID)
	}
	expiresAt, err := time.Parse(publishLeaseTimestamp, lease.LeaseExpiresAt)
	if err != nil {
		return "", fmt.Errorf("%w: lease expiry is unreadable: %v", ErrPublishLeaseUnproven, err)
	}
	if !expiresAt.After(now) {
		return "", fmt.Errorf("%w: the publish lease of %s expired at %s", ErrPublishLeaseUnproven, proof.ResourceID, lease.LeaseExpiresAt)
	}
	return lease.LeaseHash, nil
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
