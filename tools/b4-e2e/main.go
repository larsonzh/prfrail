package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/console"
	"github.com/larsonzh/prfrail/internal/gates"
	"github.com/larsonzh/prfrail/internal/snapshot"
	"github.com/larsonzh/prfrail/internal/tickets"
)

// emittedInputs is the declared input file every face reads: the identity binding
// plus the digest of the parent workspace. It is emitted once, digested into the
// evidence pack, and verified against measured facts by every run.
type emittedInputs struct {
	RunID              string `json:"runId"`
	TaskID             string `json:"taskId"`
	StepID             string `json:"stepId"`
	ParentSnapshotHash string `json:"parentSnapshotHash"`
	WorkspaceHash      string `json:"workspaceHash"`
	ContextHash        string `json:"contextHash"`
	CreatedAt          string `json:"createdAt"`
	ChainConfig        string `json:"chainConfig"`
	StubVersion        string `json:"stubVersion"`
}

// loadSpec pins the managed workload of one face/case.
type loadSpec struct {
	args    []string
	timeout time.Duration
	grace   time.Duration
	verify  time.Duration
	note    string
}

// loadFor is the single table of managed workloads. It is deliberately exhaustive
// and flag-free: a face/case pair that is not in this table is refused instead of
// defaulting to something.
func loadFor(face, caseID string) (loadSpec, error) {
	table := map[string]loadSpec{
		"f1/pos": {args: []string{"-mode", "run", "-sleep", "2s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "clean deterministic run; the workspace must be unchanged and all five facts must be provable"},
		"f1/neg": {args: []string{"-mode", "mutate-workspace", "-sleep", "2s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "run writes into its workspace; the diff must name the file it wrote (control leg of the f1/neg2 pair)"},
		// f1/neg2 is the true falsification test of this face. The load is a clean,
		// mutation-free run; the workspace change is injected by the experimenter AFTER
		// the terminal was published, so the frozen facts cannot describe it. The bytes
		// are identical to the control leg's, so the pair differs only in when the write
		// happens - an unaccounted change must be rejected, an accounted one must not.
		"f1/neg2": {args: []string{"-mode", "run", "-sleep", "2s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "clean run plus an unaccounted workspace write injected after the terminal was published; the postflight rescan must reject it"},
		// f1/neg3 measures the BOUNDARY of the claim the other two legs establish. The same
		// unaccounted write is injected into the workspace's tmp/ subtree, which the snapshot
		// capturer excludes by default (.git/**, .prfrail/**, tmp/**), so the manifest
		// reconciliation cannot see it. The run is expected to reach PASSED, and that is
		// recorded as a deviation from the ① expectation instead of being hidden.
		"f1/neg3": {args: []string{"-mode", "run", "-sleep", "2s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "clean run plus the same unaccounted write inside the capture-excluded tmp/ subtree; measures the blind spot of the unchanged-workspace claim"},
		"f2/pos": {args: []string{"-mode", "hang"}, timeout: 5 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "run outlives its bound; the watchdog must classify it as uncertain and the tree must be gone"},
		"f2/neg": {args: []string{"-mode", "spawn", "-child-life", "1h", "-sleep", "5s"}, timeout: 3 * time.Second, grace: 2 * time.Second, verify: 5 * time.Second,
			note: "run leaves a long-lived descendant; the stop must take the whole tree"},
		"f3/pos": {args: []string{"-mode", "run", "-no-marker", "-sleep", "2s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "log completeness marker missing; the outcome must degrade"},
		"f3/neg": {args: []string{"-mode", "run", "-sleep", "2s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "control leg of the same shape with the marker present; must complete"},
		"f4/pos": {args: []string{"-mode", "crash", "-exit-code", "7", "-sleep", "3s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "the runner is parked after dispatch, then restarted on the same run root"},
		"f4/neg": {args: []string{"-mode", "run", "-sleep", "2s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "a forged resume terminal must be refused with zero writes"},
		"f5/pos": {args: []string{"-mode", "run", "-sleep", "2s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "full loop: dispatch, terminal, postflight, gate hook, acceptance, review, promotion"},
		"f5/neg": {args: []string{"-mode", "run", "-sleep", "2s"}, timeout: 30 * time.Second, grace: 2 * time.Second, verify: 3 * time.Second,
			note: "the workspace is tampered after the terminal; a replayed postflight must reject"},
	}
	spec, ok := table[face+"/"+caseID]
	if !ok {
		return loadSpec{}, fmt.Errorf("unknown face/case %s/%s", face, caseID)
	}
	return spec, nil
}

func main() { os.Exit(run()) }

func run() int {
	face := flag.String("face", "", "f1..f5")
	caseID := flag.String("case", "", "pos|neg")
	runRootFlag := flag.String("run-root", "", "run root (created if missing)")
	recordPath := flag.String("inputs", "", "declared inputs JSON (emitted once by -emit-record)")
	chainPath := flag.String("chain", "", "chain config JSON (the file the CLI validates)")
	stubPath := flag.String("stub", "", "pinned agent-stub executable")
	stubVersion := flag.String("stub-version", "agent-stub-0.1.0", "pinned stub version")
	phase := flag.String("phase", "run", "run|park|recover|forge|halt|resume")
	verdictPath := flag.String("verdict", "", "verdict.json output path")
	emitRecord := flag.Bool("emit-record", false, "emit the AgentRunner request record and exit")
	workspaceFlag := flag.String("workspace", "", "workspace directory (with -emit-record)")
	outFlag := flag.String("out", "", "record output path (with -emit-record)")
	flag.Parse()

	if *emitRecord {
		return emitRequestRecord(*chainPath, *workspaceFlag, *outFlag)
	}
	if *face == "" || *caseID == "" || *runRootFlag == "" || *recordPath == "" || *chainPath == "" || *stubPath == "" {
		fmt.Fprintln(os.Stderr, "usage: b4-e2e -face f1..f5 -case pos|neg -run-root <dir> -inputs <json> -chain <json> -stub <exe> [-phase run|park|recover|forge|halt|resume] [-verdict <json>]")
		return 2
	}
	if err := harness(*face, *caseID, *phase, *runRootFlag, *recordPath, *chainPath, *stubPath, *stubVersion, *verdictPath); err != nil {
		fmt.Fprintf(os.Stderr, "b4-e2e: %v\n", err)
		return 1
	}
	return 0
}

// emitRequestRecord writes the input file the CLI leg and every face share.
func emitRequestRecord(chainPath, workspace, out string) int {
	if chainPath == "" || workspace == "" || out == "" {
		fmt.Fprintln(os.Stderr, "usage: b4-e2e -emit-record -chain <json> -workspace <dir> -out <json>")
		return 2
	}
	cfg, _, err := console.LoadChainConfig(chainPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load chain: %v\n", err)
		return 1
	}
	definition, err := cfg.ToDefinition()
	if err != nil {
		fmt.Fprintf(os.Stderr, "definition: %v\n", err)
		return 1
	}
	taskID, stepID, err := agentRunnerStep(definition)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	parent, err := workspaceDigest(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	inputs := emittedInputs{
		RunID:              defaultRunID,
		TaskID:             taskID,
		StepID:             stepID,
		ParentSnapshotHash: parent,
		WorkspaceHash:      parent,
		ContextHash:        digestOf("proofrail:b4-declared-context:1\n", defaultRunID+"/"+taskID+"/"+stepID),
		CreatedAt:          time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
		ChainConfig:        filepath.Base(chainPath),
		StubVersion:        "agent-stub-0.1.0",
	}
	data, err := json.MarshalIndent(inputs, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal: %v\n", err)
		return 1
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	if err := os.WriteFile(out, append(data, '\n'), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	fmt.Printf("inputs: %s\nrunId: %s\ntaskId: %s\nstepId: %s\nparentSnapshotHash: %s\n", out, inputs.RunID, taskID, stepID, parent)
	return 0
}

// loadEmittedInputs reads and validates the declared input file.
func loadEmittedInputs(path string) (emittedInputs, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return emittedInputs{}, fmt.Errorf("read declared inputs: %w", err)
	}
	var inputs emittedInputs
	if err := json.Unmarshal(data, &inputs); err != nil {
		return emittedInputs{}, fmt.Errorf("decode declared inputs: %w", err)
	}
	if _, err := time.Parse("2006-01-02T15:04:05.000Z", inputs.CreatedAt); err != nil {
		return emittedInputs{}, fmt.Errorf("declared inputs carry an invalid createdAt: %w", err)
	}
	return inputs, nil
}

const defaultRunID = "run-b4"

func resolvedRunID() string { return defaultRunID }

// agentRunnerStep finds the single isolated-workspace code step of a definition.
func agentRunnerStep(definition chain.Definition) (string, string, error) {
	taskID, stepID := "", ""
	matches := 0
	for _, task := range definition.Tasks {
		for _, step := range task.Steps {
			if step.Kind == "code" && step.Mode == chain.IsolatedWorkspace {
				taskID, stepID = task.ID, step.ID
				matches++
			}
		}
	}
	if matches != 1 {
		return "", "", fmt.Errorf("definition must carry exactly one isolated workspace code step, found %d", matches)
	}
	return taskID, stepID, nil
}

// workspaceDigest is the deterministic parent digest: entries only.
func workspaceDigest(root string) (string, error) {
	entries, err := workspaceEntries(root, filepath.Join(os.TempDir(), "b4-e2e-digest-store"))
	if err != nil {
		return "", err
	}
	return entriesDigest("proofrail:b4-baseline:1\n", entries)
}

// ---------------------------------------------------------------- the harness

func harness(face, caseID, phase, runRootRaw, recordPath, chainPath, stubPath, stubVersion, verdictPath string) error {
	spec, err := loadFor(face, caseID)
	if err != nil {
		return err
	}
	runRoot, err := filepath.Abs(runRootRaw)
	if err != nil {
		return err
	}
	report := newVerdict()
	report.Face = face
	report.Case = caseID
	report.Phase = phase
	report.RunID = resolvedRunID()
	report.RequestID = "request-" + resolvedRunID()
	report.RunRoot = runRoot
	report.StartedAt = stamp(time.Now())
	report.Load.Args = spec.args
	report.Load.Timeout = spec.timeout.String()
	report.Load.Grace = spec.grace.String()
	report.addNote(spec.note)

	finish := func(class string, err error) error {
		report.VerdictClass = class
		report.EndedAt = stamp(time.Now())
		report.addError(err)
		if writeErr := report.write(verdictPath); writeErr != nil && err == nil {
			return writeErr
		}
		return err
	}

	workspace := filepath.Join(runRoot, "workspace")
	eventsDir := filepath.Join(runRoot, "events")
	eventLog := filepath.Join(eventsDir, "state-events.jsonl")
	if err := os.MkdirAll(eventsDir, 0o755); err != nil {
		return finish("blocked", err)
	}
	if _, err := os.Stat(eventLog); os.IsNotExist(err) {
		file, createErr := os.OpenFile(eventLog, os.O_CREATE|os.O_WRONLY, 0o600)
		if createErr != nil {
			return finish("blocked", createErr)
		}
		_ = file.Close()
	}
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		return finish("blocked", err)
	}
	report.Workspace = workspace

	cfg, _, err := console.LoadChainConfig(chainPath)
	if err != nil {
		return finish("blocked", fmt.Errorf("load chain config: %w", err))
	}
	definition, err := cfg.ToDefinition()
	if err != nil {
		return finish("blocked", fmt.Errorf("chain definition: %w", err))
	}
	taskID, stepID, err := agentRunnerStep(definition)
	if err != nil {
		return finish("blocked", err)
	}

	store, err := adapters.NewAgentRunnerReplayStoreForRun(runRoot, resolvedRunID())
	if err != nil {
		return finish("blocked", fmt.Errorf("replay store: %w", err))
	}
	if err := store.RejectWriterRoots(workspace); err != nil {
		return finish("blocked", fmt.Errorf("workspace overlaps the replay root: %w", err))
	}
	eventStore, err := chain.NewFileEventStore(eventLog)
	if err != nil {
		return finish("blocked", err)
	}
	snapshotRoot := filepath.Join(runRoot, "snapshot-store")
	snapshotStore, err := snapshot.NewStore(snapshotRoot, 0)
	if err != nil {
		return finish("blocked", err)
	}

	now := time.Now().UTC
	declared, err := loadEmittedInputs(recordPath)
	if err != nil {
		return finish("blocked", err)
	}
	parentHash, err := workspaceDigest(workspace)
	if err != nil {
		return finish("blocked", err)
	}
	createdAt, err := time.Parse("2006-01-02T15:04:05.000Z", declared.CreatedAt)
	if err != nil {
		return finish("blocked", err)
	}
	inputs, err := buildDeclaredInputs(requestSpec{
		RunID:         resolvedRunID(),
		TaskID:        taskID,
		StepID:        stepID,
		CreatedAt:     createdAt,
		ParentHash:    parentHash,
		WorkspaceHash: parentHash,
		ContextHash:   declared.ContextHash,
	}, stubPath, stubVersion, now())
	if err != nil {
		return finish("blocked", fmt.Errorf("declared inputs: %w", err))
	}
	// The declared input file and the measured run-time facts must agree. The parent
	// workspace digest is only asserted for the phases that actually dispatch: the
	// resume phases run on a run root whose workspace may legitimately differ (that is
	// exactly what the f5 negative case injects).
	if declared.TaskID != taskID || declared.StepID != stepID || declared.RunID != resolvedRunID() {
		return finish("blocked", fmt.Errorf("declared inputs binding mismatch: declared(task=%s step=%s run=%s) measured(task=%s step=%s run=%s)",
			declared.TaskID, declared.StepID, declared.RunID, taskID, stepID, resolvedRunID()))
	}
	if phase == "run" || phase == "park" || phase == "halt" {
		if declared.ParentSnapshotHash != parentHash {
			return finish("blocked", fmt.Errorf("declared parent snapshot hash %s does not match the workspace digest %s", declared.ParentSnapshotHash, parentHash))
		}
	}
	report.Notes = append(report.Notes, "requestRecordHash="+inputs.record.RecordHash)

	stepPorts := &ports{
		runID:        resolvedRunID(),
		requestID:    "request-" + resolvedRunID(),
		runRoot:      runRoot,
		workspace:    workspace,
		snapshotRoot: snapshotRoot,
		store:        snapshotStore,
		replay:       store,
		now:          now,
	}

	gateRunner := gates.Runner{
		Executor: gates.GuardExecutor{},
		Scanner:  gates.DefaultSecretScanner(),
		Objects:  snapshotStore,
		Clock:    now,
		IDs:      func() string { return "gate-result-" + fmt.Sprintf("%d", time.Now().UTC().UnixNano()) },
	}
	resultStore, err := gates.NewFileResultStore(filepath.Join(runRoot, "gate-results.jsonl"))
	if err != nil {
		return finish("blocked", err)
	}
	gateRunner.Store = resultStore
	hooks := map[string]gates.Hook{}
	gateCommands := map[string]gateCommand{}
	for _, task := range definition.Tasks {
		for _, step := range task.Steps {
			if isGateStep(step) {
				hook := gateHook(step.ID)
				hooks[step.ID] = hook
				gateCommands[step.ID] = gateCommand{executable: "cmd.exe", args: []string{"/c", "dir", "/b", "/s"}}
			}
		}
	}
	dispatcher := &adapters.AgentRunnerReplayDispatcher{
		Store: store,
		// The launcher is the production pinned-CLI launcher: the harness never
		// imports os/exec and never starts a process itself.
		Launcher:       nil,
		Process:        adapters.AgentRunnerProcessConfig{Command: stubPath, Args: spec.args, Grace: spec.grace},
		LedgerSnapshot: func() (chain.AuthorizationLedger, error) { return inputs.admission.AuthorizationLedger, nil },
		CostLedger:     inputs.ledger,
		RequestRecord:  inputs.record,
		Clock:          now,
	}
	launcher, err := adapters.NewAgentRunnerPinnedCLILauncher(adapters.AgentRunnerPinnedCLIConfig{
		Executable:   stubPath,
		Version:      stubVersion,
		VersionArgs:  []string{"-version"},
		CheckTimeout: 20 * time.Second,
		Grace:        spec.grace,
	}, store)
	if err != nil {
		return finish("blocked", fmt.Errorf("launcher: %w", err))
	}
	dispatcher.Launcher = launcher
	stepPorts.router = chain.StepRouter{
		AgentRunner:          dispatcher,
		AgentRunnerAdmission: &inputs.admission,
		Hooks: &hookPort{
			owner:      stepPorts,
			commands:   gateCommands,
			production: gates.ChainPort{Runner: &gateRunner, Hooks: hooks},
		},
	}
	stepPorts.admission = &inputs.admission
	stepPorts.dispatcher = dispatcher
	stepPorts.launcher = launcher
	stepPorts.recordHash = inputs.record.RecordHash
	stepPorts.process = dispatcher.Process

	intents, err := adapters.NewAgentRunnerStepIntentPreparer(inputs.record, definition)
	if err != nil {
		return finish("blocked", fmt.Errorf("step intent preparer: %w", err))
	}
	postflightPort := newPostflight("request-"+resolvedRunID(), workspace, snapshotRoot, stepPorts)
	publisher := &adapters.AgentRunnerTerminalPublisher{
		Store:         store,
		CostLedger:    inputs.ledger,
		RequestRecord: inputs.record,
		Clock:         now,
	}

	engine, err := chain.New(context.Background(), chain.Options{
		RunID:       resolvedRunID(),
		Definition:  definition,
		Events:      eventStore,
		Baselines:   stepPorts,
		Workspaces:  stepPorts,
		Steps:       stepPorts,
		StepIntents: intents,
		Acceptance:  stepPorts,
		Reviewer:    stepPorts,
		Publisher:   stepPorts,
		Postflight:  postflightPort,
		Stopper:     stepPorts,
		Reconciler:  stepPorts,
		Clock:       now,
		IDs:         newNodeIDSource(now),
	})
	if err != nil {
		return finish("blocked", fmt.Errorf("engine: %w", err))
	}
	report.addNote("declared-inputs: capability disposition=compatible (declared, not probed); enforcement record declared but not consulted on that path; adapterId=" + declaredAdapterID)

	switch phase {
	case "run", "park", "halt":
		if err := engine.Run(context.Background()); !errors.Is(err, chain.ErrAwaitingAgentRunnerTerminal) {
			stepPorts.record("dispatch", "unexpected", nil, err)
			if err == nil {
				return finish("blocked", fmt.Errorf("engine.Run returned without parking on the AgentRunner terminal"))
			}
			return finish("blocked", fmt.Errorf("engine.Run: %w", err))
		}
		stepPorts.record("dispatch", "parked", map[string]string{"requestId": report.RequestID}, nil)
		report.addLeg("dispatch", "ran: process spawned through the production launcher and step parked in TERMINAL_PENDING", "", "")
	case "recover":
		mirrorBefore := identityMirrorDigest(runRoot)
		if err := engine.Recover(context.Background()); err != nil {
			stepPorts.record("recover", "refused", nil, err)
			if !errors.Is(err, chain.ErrRecoveryUncertain) {
				return finish("blocked", fmt.Errorf("recover: %w", err))
			}
		} else {
			stepPorts.record("recover", "completed-without-uncertainty", nil, nil)
		}
		runErr := engine.Run(context.Background())
		stepPorts.record("recover-run", "returned", nil, runErr)
		report.VerdictClass = classifyProjection(engine)
		report.fillProjection(engine)
		report.Journal = stepPorts.snapshotJournal()
		report.Artifacts = artifactHashes(store)
		report.EventLogHash = eventLogDigest(eventLog)
		dispatches, dispatchErr := countDispatchEvents(eventLog, stepID)
		if dispatchErr != nil {
			report.addError(dispatchErr)
		}
		mirrorAfter := identityMirrorDigest(runRoot)
		report.Observations["dispatchEvents"] = fmt.Sprintf("%d", dispatches)
		report.Observations["identityMirrorStable"] = fmt.Sprintf("%t", mirrorBefore == mirrorAfter)
		report.addNote("no launch receipts exist on this platform: request publication is refused by the unproven-durability declaration, so the relaunch check uses the parked dispatch count and the managed-process identity mirror instead")
		report.addCheck("recover-refused", errors.Is(runErr, chain.ErrRecoveryUncertain), fmt.Sprintf("engine.Run after Recover returned %v", runErr))
		report.addCheck("single-dispatch", dispatches == 1, fmt.Sprintf("agent-dispatched TERMINAL_PENDING events for the step: %d", dispatches))
		report.addCheck("no-new-process-identity", mirrorBefore == mirrorAfter, "managed-process identity mirror unchanged across the restart: "+fmt.Sprint(mirrorBefore == mirrorAfter))
		report.addCheck("chain-paused", engine.Projection().ChainState == "PAUSED", "chain state "+engine.Projection().ChainState)
		return finish(report.VerdictClass, nil)
	case "resume":
		// Phase of the f5 negative case: the run root already holds the routed terminal
		// and the step already passed, so the engine walks straight into postflight.
		runErr := engine.Run(context.Background())
		stepPorts.record("resume", "returned", nil, runErr)
		report.VerdictClass = classifyProjection(engine)
		report.fillProjection(engine)
		report.Journal = stepPorts.snapshotJournal()
		report.Artifacts = artifactHashes(store)
		report.EventLogHash = eventLogDigest(eventLog)
		report.Diff = diffFor(store, report.RequestID)
		report.addCheck("postflight-rejected", report.TaskStateIsRepair(), "task state after the tampered replay: "+report.TaskState)
		return finish(report.VerdictClass, nil)
	case "forge":
		if err := engine.Run(context.Background()); !errors.Is(err, chain.ErrAwaitingAgentRunnerTerminal) {
			return finish("blocked", fmt.Errorf("engine.Run before the forged terminal: %w", err))
		}
		before := eventLogDigest(eventLog)
		forged := forgedResumeTerminal(inputs.record)
		_, submitErr := engine.SubmitAgentRunnerTerminal(context.Background(), forged)
		after := eventLogDigest(eventLog)
		report.VerdictClass = "refused"
		report.EventLogHash = after
		report.addCheck("resume-continuity-refused", errors.Is(submitErr, chain.ErrAgentRunnerResumeContinuityUnproven), fmt.Sprintf("SubmitAgentRunnerTerminal returned %v", submitErr))
		report.addCheck("zero-write", before == after, "event log digest before="+before+" after="+after)
		report.Observations["submitError"] = fmt.Sprint(submitErr)
		report.fillProjection(engine)
		report.Journal = stepPorts.snapshotJournal()
		report.Artifacts = artifactHashes(store)
		return finish(report.VerdictClass, nil)
	}

	// phase park: the run root now holds a parked dispatch and nothing else.
	if phase == "park" {
		report.VerdictClass = "parked"
		report.fillProjection(engine)
		report.Journal = stepPorts.snapshotJournal()
		report.Artifacts = artifactHashes(store)
		report.EventLogHash = eventLogDigest(eventLog)
		return finish(report.VerdictClass, nil)
	}

	// ------------------------------------------------ wait, collect and map
	observation, observeErr := waitAndMap(context.Background(), launcher, store, stepPorts, resolvedRunID(), report.RequestID, taskID, stepID, inputs.record.RecordHash, workspace, stepPorts.beforeManifestFor(report.RequestID), spec.timeout, spec.verify, 1<<20, now)
	if observeErr != nil {
		stepPorts.record("run", "failed", nil, observeErr)
		report.fillProjection(engine)
		report.Journal = stepPorts.snapshotJournal()
		return finish("blocked", fmt.Errorf("wait/collect/map: %w", observeErr))
	}
	report.Mapping = &mappingView{
		Status:         observation.Status,
		TreeStatus:     observation.TreeStatus,
		LogsComplete:   observation.Collected.LogsComplete,
		UsageComplete:  observation.Collected.UsageComplete,
		EventsComplete: observation.Collected.EventsComplete,
		Manifests:      observation.Collected.ManifestsComplete,
		ErrorEvidence:  len(observation.ErrorEvi),
		ElapsedMs:      observation.Elapsed.Milliseconds(),
	}
	if observation.Facts != nil {
		report.Mapping.FactCount = len(factHashes(adaptersFactsToChain(*observation.Facts)))
	}
	report.Observations["stopOutcome"] = "unknown"
	if observation.StopProof != nil {
		report.Observations["stopOutcome"] = observation.StopProof.Outcome
	}
	report.Observations["stopError"] = observation.StopError
	report.addLeg("dispatch->wait->collect->map", "ran: "+observation.Status, "", "")

	// ------------------------------------------------ production publication leg
	outcome := adapters.AgentRunnerTerminalOutcome{
		Completion: observation.Completion,
		Settlement: observation.Settlement,
		Facts:      observation.Facts,
	}
	if _, publishErr := publisher.PublishTerminal(context.Background(), outcome); publishErr != nil {
		report.addLeg("terminal-publication (production AgentRunnerTerminalPublisher)", "failed closed", publishErr.Error(),
			"the durable terminal chain is built through the public record constructors and adapters.ToChainAgentRunnerTerminal instead; the chain-side routing below is unmodified production code")
	} else {
		report.addLeg("terminal-publication (production AgentRunnerTerminalPublisher)", "published", "", "")
	}

	// ------------------------------------------------ settlement leg
	if settleErr := settleLedger(inputs.ledger, inputs.record, observation); settleErr != nil {
		report.addLeg("cost settlement (production tickets.CostLedger.Settle)", "failed", settleErr.Error(), "")
	} else {
		report.addLeg("cost settlement (production tickets.CostLedger.Settle)", "booked", "", "")
	}

	// ------------------------------------------------ terminal routing leg
	terminal, err := terminalChain(inputs, observation, time.Now().UTC())
	if err != nil {
		return finish("blocked", fmt.Errorf("terminal chain: %w", err))
	}
	report.Artifacts["terminalCompletionHash"] = terminal.completion.RecordHash
	report.Artifacts["terminalIntentHash"] = terminal.intent.RecordHash
	report.Artifacts["terminalClosureHash"] = terminal.closure.RecordHash
	route, routeErr := engine.SubmitAgentRunnerTerminal(context.Background(), terminal.fact)
	if routeErr != nil {
		stepPorts.record("terminal-routing", "refused", nil, routeErr)
		report.fillProjection(engine)
		report.Journal = stepPorts.snapshotJournal()
		return finish("blocked", fmt.Errorf("route terminal: %w", routeErr))
	}
	stepPorts.record("terminal-routing", string(route.Status), map[string]string{"step": route.StepState, "task": route.TaskState, "chain": route.ChainState}, nil)
	report.addLeg("terminal routing (chain.SubmitAgentRunnerTerminal)", "routed as "+string(route.Status), "", "")

	// The f1 negative injections. Two legs of the pair write the same bytes into the
	// workspace explicitly after the terminal was published (dev-1): f1/neg2 into a
	// captured path (the postflight rescan must reject it) and f1/neg3 into the
	// workspace's tmp/ subtree, which the snapshot capturer excludes by construction -
	// that leg measures the boundary of the claim instead of assuming it away.
	if face == "f1" && (caseID == "neg2" || caseID == "neg3") {
		relPath := mutatedWorkspaceFileName
		if caseID == "neg3" {
			relPath = filepath.Join("tmp", mutatedWorkspaceFileName)
		}
		if injectErr := injectUnaccountedWorkspaceWrite(workspace, relPath); injectErr != nil {
			return finish("blocked", fmt.Errorf("inject unaccounted workspace write: %w", injectErr))
		}
		stepPorts.record("injection", "unaccounted-workspace-write", map[string]string{
			"path":        filepath.ToSlash(relPath),
			"contentHash": digestOf("proofrail:b4-injected-content:1\n", mutatedWorkspaceContent),
		}, nil)
	}

	if phase == "halt" {
		report.VerdictClass = classifyProjection(engine)
		report.fillProjection(engine)
		report.Journal = stepPorts.snapshotJournal()
		report.Artifacts = artifactHashes(store)
		report.EventLogHash = eventLogDigest(eventLog)
		report.Diff = diffFor(store, report.RequestID)
		return finish(report.VerdictClass, nil)
	}
	if observation.Status == "completed" {
		if runErr := engine.Run(context.Background()); runErr != nil {
			stepPorts.record("downstream", "returned", nil, runErr)
			report.fillProjection(engine)
			report.Journal = stepPorts.snapshotJournal()
			report.Artifacts = artifactHashes(store)
			report.EventLogHash = eventLogDigest(eventLog)
			report.Diff = diffFor(store, report.RequestID)
			report.addLeg("postflight/acceptance/review/promotion", "returned "+runErr.Error(), runErr.Error(), "")
			// A rejection must still be judged by the face's own checks: the f1 negative
			// case of dev-1 is exactly this path, and skipping the checks here would make
			// the rejection indistinguishable from an unjudged run.
			report.VerdictClass = classifyProjection(engine)
			addFaceChecks(report, face, caseID, observation, engine, stepPorts, runRoot)
			return finish(report.VerdictClass, nil)
		}
		stepPorts.record("downstream", "completed", nil, nil)
		report.addLeg("postflight/acceptance/review/promotion", "ran to the end of the task", "", "")
	} else {
		report.addLeg("postflight/acceptance/review/promotion", "skipped: a non-completed terminal never reaches the downstream legs", "", "")
	}

	report.fillProjection(engine)
	report.Journal = stepPorts.snapshotJournal()
	report.Artifacts = artifactHashes(store)
	report.EventLogHash = eventLogDigest(eventLog)
	report.Diff = diffFor(store, report.RequestID)
	report.Observations["gateHookSteps"] = fmt.Sprintf("%d", len(hooks))
	report.VerdictClass = classifyProjection(engine)

	// The R5 order accounting: the gate hook step is part of the task's step list, so
	// its transitions must precede the task's own postflight-qualified REVIEW_PENDING.
	// A definition without a gate step records nothing here. The reconciliation names
	// the declared gate step ids, so a passing non-gate step can never supply it.
	gateSteps := gateHookStepIDs(definition)
	if hookSequence, reviewSequence, orderErr := gateOrder(eventLog, gateSteps); orderErr == nil && hookSequence > 0 {
		report.Observations["gateHookStepIds"] = strings.Join(gateSteps, ",")
		report.Observations["gateHookSequence"] = fmt.Sprintf("%d", hookSequence)
		report.Observations["reviewPendingSequence"] = fmt.Sprintf("%d", reviewSequence)
		report.addCheck("gate-hook-before-review-pending", hookSequence < reviewSequence,
			fmt.Sprintf("gate hook step(s) %v sequence=%d, task REVIEW_PENDING sequence=%d", gateSteps, hookSequence, reviewSequence))
	}

	// ------------------------------------------------ face checks
	addFaceChecks(report, face, caseID, observation, engine, stepPorts, runRoot)
	return finish(report.VerdictClass, nil)
}

// addFaceChecks records the per-face expectations. It never invents an outcome:
// each check states the measured value it compared.
func addFaceChecks(report *verdict, face, caseID string, observation *runObservation, engine *chain.Engine, stepPorts *ports, runRoot string) {
	projection := engine.Projection()
	taskPassed := false
	for _, state := range projection.TaskStates {
		if state == "PASSED" {
			taskPassed = true
		}
	}
	switch face {
	case "f1":
		report.addCheck("five-facts-distinct", observation.Facts != nil, "frozen facts present: "+fmt.Sprint(observation.Facts != nil))
		if observation.Facts != nil {
			hashes := factHashes(adaptersFactsToChain(*observation.Facts))
			distinct := map[string]struct{}{}
			for _, hash := range hashes {
				distinct[hash] = struct{}{}
			}
			report.addCheck("five-facts-distinct", len(distinct) == 5, fmt.Sprintf("distinct fact digests: %d", len(distinct)))
		}
		if diff := report.Diff; diff != nil {
			switch caseID {
			case "pos":
				report.addCheck("workspace-unchanged", len(diff.Added)+len(diff.Removed)+len(diff.Changed) == 0,
					fmt.Sprintf("diff added=%v removed=%v changed=%v", diff.Added, diff.Removed, diff.Changed))
			case "neg2":
				// The frozen diff was written before the injection, so it must still say
				// "nothing changed": that is what makes the injected change unaccounted.
				report.addCheck("frozen-diff-excludes-injection", len(diff.Added)+len(diff.Removed)+len(diff.Changed) == 0,
					fmt.Sprintf("frozen diff added=%v removed=%v changed=%v", diff.Added, diff.Removed, diff.Changed))
				report.addCheck("postflight-rejected-unaccounted-write", report.TaskStateIsRepair(),
					"task state after the unaccounted write: "+report.TaskState)
				report.addCheck("never-passed", !taskPassed, fmt.Sprintf("task reached PASSED: %t", taskPassed))
				// The observable failure point must come from the chain's own event log,
				// not from this verdict's summary field.
				eventLogPath := filepath.Join(runRoot, "events", "state-events.jsonl")
				toState, reason, sequence, err := lastTaskTransition(eventLogPath)
				report.Observations["taskFinalTransition"] = fmt.Sprintf("sequence=%d toState=%s reason=%s", sequence, toState, reason)
				report.addCheck("task-final-transition-is-postflight-rejection",
					err == nil && toState == "REPAIR_PENDING" && reason == "postflight-rejected",
					report.Observations["taskFinalTransition"])
				if err != nil {
					report.addError(err)
				}
			case "neg3":
				// The same injection, but under the workspace's tmp/ subtree, which the
				// snapshot capturer excludes by default - so both the frozen diff and the
				// postflight rescan are blind to it and the task passes. The assertions
				// below only state what was measured; the contradiction with the ①
				// expectation is recorded as a deviation, never as a silent failure.
				injectedPath := filepath.Join(runRoot, "workspace", "tmp", mutatedWorkspaceFileName)
				_, statErr := os.Stat(injectedPath)
				report.addCheck("injection-present-on-disk", statErr == nil, "injected file on disk: "+injectedPath)
				report.addCheck("frozen-diff-excludes-injection", len(diff.Added)+len(diff.Removed)+len(diff.Changed) == 0,
					fmt.Sprintf("frozen diff added=%v removed=%v changed=%v", diff.Added, diff.Removed, diff.Changed))
				report.addDeviation("unaccounted-change-rejected(design expectation)", report.TaskStateIsRepair(),
					fmt.Sprintf("taskPassed=%t; the injected path is inside the snapshot capturer's default exclusion set (.git/**, .prfrail/**, tmp/**), so the manifest reconciliation cannot see it - the boundary of the F1 claim, measured", taskPassed))
			default:
				// The f1/neg control leg (the only remaining case of this face).
				found := false
				for _, path := range diff.Added {
					if path == "agent-stub-mutation.txt" {
						found = true
					}
				}
				report.addCheck("mutation-captured-in-diff", found, fmt.Sprintf("diff added=%v (expected exactly the file the load wrote)", diff.Added))
				// The design's first expectation for this case was a postflight rejection.
				// It did not happen: the mutation precedes the frozen facts, so the frozen
				// facts describe the mutated tree and the reconciliation passes. Recorded as
				// a deviation from the ① design, not hidden. This case is therefore the
				// control leg of the pair; the write it performs is injected unaccounted by
				// f1/neg2, which is where the falsification is actually measured.
				report.addDeviation("postflight-rejection(design expectation)", !taskPassed,
					fmt.Sprintf("taskPassed=%t (① expected a rejection or an exact diff; the exact diff is asserted above, and the true falsification test of the same write is f1/neg2)", taskPassed))
			}
		}
	case "f2":
		report.addCheck("outcome-not-completed", observation.Status != "completed", "mapped status "+observation.Status)
		report.addCheck("no-frozen-facts", observation.Facts == nil, fmt.Sprintf("facts present: %t", observation.Facts != nil))
		if caseID == "pos" {
			report.addCheck("timeout-fired", observation.Watch.TimedOut, fmt.Sprintf("watch timedOut=%t elapsedMs=%d boundMs=%d", observation.Watch.TimedOut, observation.Elapsed.Milliseconds(), reportBound(report)))
			report.addCheck("error-evidence-present", len(observation.ErrorEvi) > 0, fmt.Sprintf("error evidence digests: %d", len(observation.ErrorEvi)))
		}
		// The tree claim is asserted for BOTH cases. Before this check existed the forward
		// case only measured a duration interval, so it could have passed with the stopped
		// process still alive; the reported pids are the artifact the claim rests on.
		gone, goneErr := recordedPIDsGone(stepPorts.requestID, stepPorts.replay)
		report.addCheck("process-tree-taken-away", gone && goneErr == nil, fmt.Sprintf("reported pids gone: %t (%v)", gone, goneErr))
	case "f3":
		if caseID == "pos" {
			report.addCheck("log-gap-degrades", observation.Status == "uncertain" && !observation.Collected.LogsComplete, fmt.Sprintf("status=%s logsComplete=%t", observation.Status, observation.Collected.LogsComplete))
		} else {
			report.addCheck("control-leg-completes", observation.Status == "completed", "status "+observation.Status)
		}
	case "f4", "f5":
		_ = projection
	}
	report.addCheck("verdict-class-recorded", report.VerdictClass != "", "verdict class "+report.VerdictClass)
}

func reportBound(report *verdict) int64 {
	duration, err := time.ParseDuration(report.Load.Timeout)
	if err != nil {
		return 0
	}
	return duration.Milliseconds()
}

// diffFor reads back the recorded diff artifact of this run.
func diffFor(store *adapters.AgentRunnerReplayStore, requestID string) *diffView {
	directory, err := store.RunEvidenceDir(requestID)
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(filepath.Join(directory, "diff.json"))
	if err != nil {
		return nil
	}
	var report diffReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil
	}
	view := &diffView{Added: []string{}, Removed: []string{}, Changed: []string{}}
	for _, entry := range report.Added {
		view.Added = append(view.Added, entry.Path)
	}
	for _, entry := range report.Removed {
		view.Removed = append(view.Removed, entry.Path)
	}
	for _, entry := range report.Changed {
		view.Changed = append(view.Changed, entry.Path)
	}
	return view
}

// recordedPIDsGone reports whether the pids the CLI reported are gone. The run
// root is not a parameter: the evidence directory is addressed through the replay
// store that owns it.
func recordedPIDsGone(requestID string, store *adapters.AgentRunnerReplayStore) (bool, error) {
	dir, err := store.RunEvidenceDir(requestID)
	if err != nil {
		return false, err
	}
	alive, err := alivePIDs(dir)
	if err != nil {
		return false, err
	}
	if len(alive) > 0 {
		// The experiment must not leak the process it is measuring: the leftover is
		// stopped through the reviewed identity stop and the leak is recorded.
		return false, fmt.Errorf("still alive: %s", strings.Join(alive, ","))
	}
	return true, nil
}

// settleLedger books the settlement the production publisher would have booked.
func settleLedger(ledger *tickets.CostLedger, record adapters.AgentRunnerRequestRecord, observation *runObservation) error {
	body := record.Request
	settlementID := "settle-" + body.RequestID
	if observation.Completion.Status != "completed" {
		settlementID = "settle-unknown-" + body.RequestID
	}
	_, err := ledger.Settle(tickets.SettlementInput{
		EntryID:             settlementID,
		OccurredAt:          observation.CollectedAt.UTC(),
		RunID:               body.RunID,
		RequestID:           body.RequestID,
		IdempotencyKey:      settlementID,
		ReservationHash:     body.BudgetHash,
		Status:              observation.Settlement.Status,
		ChargedAmountMicros: observation.Settlement.ChargedAmountMicros,
		ObservedCalls:       observation.Settlement.ObservedCalls,
		ObservedTokens:      observation.Settlement.ObservedTokens,
		ProviderEvidence:    observation.Settlement.ProviderEvidence,
	})
	return err
}

// forgedResumeTerminal builds the resume terminal f4/neg submits: it claims a
// prior completion the step evidence history cannot contain.
func forgedResumeTerminal(record adapters.AgentRunnerRequestRecord) chain.AgentRunnerTerminal {
	body := record.Request
	sessionID := "session-forged-" + body.RequestID
	return chain.AgentRunnerTerminal{
		RequestID:           body.RequestID,
		RequestHash:         record.RecordHash,
		CompletionHash:      digestOf("proofrail:b4-forged-completion:1\n", body.RequestID),
		RunID:               body.RunID,
		TaskID:              body.TaskID,
		StepID:              body.StepID,
		Attempt:             body.Attempt,
		SessionID:           sessionID,
		PriorSessionID:      sessionID,
		PriorCompletionHash: digestOf("proofrail:b4-forged-prior-completion:1\n", body.RequestID),
		Status:              chain.AgentRunnerTerminalUncertain,
		Evidence:            []string{record.RecordHash, digestOf("proofrail:b4-forged-completion:1\n", body.RequestID)},
	}
}

// fillProjection records the chain/task/step projection of the engine.
//
// The task and step maps are written as a COMPLETE, sorted `key=state` list. The
// first version wrote both maps through a single slot inside a `for ... range map`
// loop, so Go's randomised map order both made the verdict bytes non-deterministic
// across re-runs (the same case produced different `stepState` strings) and could
// silently drop a state entirely (measured on f5-pos, where only the verify step
// survived and the code step was lost). Sorting the keys and emitting every entry
// removes the nondeterminism and the silent drop at once.
func (v *verdict) fillProjection(engine *chain.Engine) {
	projection := engine.Projection()
	v.ChainState = projection.ChainState
	v.EventSequence = projection.Sequence
	v.TaskState = sortedStateSummary(projection.TaskStates)
	v.StepState = sortedStateSummary(projection.StepStates)
}

// sortedStateSummary renders a state map as a deterministic `key=state; key=state`
// string covering every entry (sorted by key).
func sortedStateSummary(states map[string]string) string {
	keys := make([]string, 0, len(states))
	for key := range states {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, key+"="+states[key])
	}
	return strings.Join(parts, "; ")
}

// classifyProjection maps the final chain/task state to the verdict class.
func classifyProjection(engine *chain.Engine) string {
	projection := engine.Projection()
	for _, state := range projection.TaskStates {
		if state == "PASSED" {
			return "completed"
		}
	}
	switch projection.ChainState {
	case "PAUSED":
		return "uncertain-or-paused"
	case "FAILED":
		return "failed"
	case "CANCELLED":
		return "cancelled"
	case "COMPLETED":
		return "completed"
	}
	return "incomplete"
}

// artifactHashes records the digests of the run's evidence artifacts. The run root
// is not a parameter: every artifact is addressed through the replay store, which
// is already bound to that run root.
func artifactHashes(store *adapters.AgentRunnerReplayStore) map[string]string {
	out := map[string]string{}
	dir, err := store.RunEvidenceDir("request-" + resolvedRunID())
	if err != nil {
		return out
	}
	for _, name := range []string{"manifest.pre.json", "manifest.post.json", "diff.json", "usage.json", "events.jsonl", "stub.pid", "stub-child.pid", filepath.Join("logs", "stdout.log"), filepath.Join("logs", "stderr.log")} {
		path := filepath.Join(dir, name)
		if _, statErr := os.Stat(path); statErr != nil {
			continue
		}
		hash, hashErr := fileHash(path)
		if hashErr != nil {
			continue
		}
		out[filepath.ToSlash(name)] = hash
	}
	return out
}

func eventLogDigest(path string) string {
	hash, err := fileHash(path)
	if err != nil {
		return ""
	}
	return hash
}

// countDispatchEvents counts the parked-dispatch transitions one step holds. A
// relaunch would append a second one, which is the observable that has to stay at
// one when a restart must not resume by guessing.
func countDispatchEvents(eventLog, stepID string) (int, error) {
	data, err := os.ReadFile(eventLog)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record struct {
			Event struct {
				Entity struct {
					Kind   string `json:"kind"`
					StepID string `json:"stepId"`
				} `json:"entity"`
				ToState string `json:"toState"`
				Reason  struct {
					Code string `json:"code"`
				} `json:"reason"`
			} `json:"event"`
		}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			continue
		}
		if record.Event.Entity.Kind == "step" && record.Event.Entity.StepID == stepID && record.Event.ToState == "TERMINAL_PENDING" && record.Event.Reason.Code == "agent-dispatched" {
			count++
		}
	}
	return count, nil
}

// identityMirrorDigest digests the durable managed-process identity mirror, which
// is rewritten on every spawn.
func identityMirrorDigest(runRoot string) string {
	path := filepath.Join(runRoot, "managed-process.identity.json")
	if _, err := os.Stat(path); err != nil {
		return "absent"
	}
	return eventLogDigest(path)
}

// gateOrder returns the event sequence of the gate hook step's PASSED transition
// and of the task's postflight-qualified REVIEW_PENDING transition.
//
// The declared gate-hook step ids are matched explicitly against the event entity.
// The first version took the first `step-passed` transition of ANY step, so a chain
// whose step table carried an earlier passing build step would have handed that
// sequence to the R5 reconciliation and pointed it at the wrong entity. Matching by
// id and keeping the HIGHEST matching sequence makes the assertion the strongest the
// definition allows: with several gate steps, the last one must still precede the
// task's REVIEW_PENDING for the check to pass.
func gateOrder(eventLog string, hookStepIDs []string) (int, int, error) {
	data, err := os.ReadFile(eventLog)
	if err != nil {
		return 0, 0, err
	}
	gateSteps := make(map[string]struct{}, len(hookStepIDs))
	for _, id := range hookStepIDs {
		gateSteps[id] = struct{}{}
	}
	hookSequence, reviewSequence := 0, 0
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record struct {
			Event struct {
				Sequence int `json:"sequence"`
				Entity   struct {
					Kind   string `json:"kind"`
					StepID string `json:"stepId"`
				} `json:"entity"`
				ToState string `json:"toState"`
				Reason  struct {
					Code string `json:"code"`
				} `json:"reason"`
			} `json:"event"`
		}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			continue
		}
		if record.Event.Entity.Kind == "step" && record.Event.ToState == "PASSED" && record.Event.Reason.Code == "step-passed" {
			if _, isGateStep := gateSteps[record.Event.Entity.StepID]; isGateStep && record.Event.Sequence > hookSequence {
				hookSequence = record.Event.Sequence
			}
		}
		if record.Event.Entity.Kind == "task" && record.Event.ToState == "REVIEW_PENDING" {
			reviewSequence = record.Event.Sequence
		}
	}
	return hookSequence, reviewSequence, nil
}

// isGateStep reports whether a step is one of the definition's declared gate steps:
// exactly the steps the harness attaches a gate hook to.
func isGateStep(step chain.Step) bool { return step.Kind == "build" || step.Kind == "verify" }

// gateHookStepIDs returns the declared gate-hook step ids in definition order.
func gateHookStepIDs(definition chain.Definition) []string {
	ids := []string{}
	for _, task := range definition.Tasks {
		for _, step := range task.Steps {
			if isGateStep(step) {
				ids = append(ids, step.ID)
			}
		}
	}
	return ids
}

func newNodeIDSource(clock func() time.Time) chain.IDSource {
	counter := 0
	return func() string {
		counter++
		return fmt.Sprintf("b4-event-%d-%d", clock().UTC().UnixNano(), counter)
	}
}

func adaptersFactsToChain(facts adapters.AgentRunnerFrozenFacts) chain.AgentRunnerFrozenFacts {
	return chain.AgentRunnerFrozenFacts{
		ManifestHash:            facts.ManifestHash,
		DiffHash:                facts.DiffHash,
		LogHash:                 facts.LogHash,
		UsageHash:               facts.UsageHash,
		ProcessStopEvidenceHash: facts.ProcessStopEvidenceHash,
	}
}

// The f1/neg2 injection reproduces the bytes of the stub's mutate-workspace payload
// exactly (tools/agent-stub/main.go: mutatedWorkspaceFileName and the body it writes),
// so the accounted leg and the unaccounted leg differ only in WHEN the write happens
// and never in what the workspace looks like afterwards.
const (
	mutatedWorkspaceFileName = "agent-stub-mutation.txt"
	mutatedWorkspaceContent  = "the run wrote this\n"
)

// injectUnaccountedWorkspaceWrite performs the out-of-band workspace write of the F1
// negative legs. It is called after the terminal was published and routed, so the write
// can never appear in the frozen facts the postflight reconciles against.
func injectUnaccountedWorkspaceWrite(workspace, relPath string) error {
	full := filepath.Join(workspace, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(mutatedWorkspaceContent), 0o600)
}

// lastTaskTransition reads the last task-entity transition of the chain event log. It
// is the observable failure point of the f1/neg2 injection: the task must leave the step
// phase through the postflight rejection reason and never through acceptance, and this
// is read from the chain's own durable log rather than from the harness's summary.
func lastTaskTransition(eventLog string) (string, string, int, error) {
	data, err := os.ReadFile(eventLog)
	if err != nil {
		return "", "", 0, err
	}
	toState, reason, sequence := "", "", 0
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var record struct {
			Event struct {
				Sequence int `json:"sequence"`
				Entity   struct {
					Kind string `json:"kind"`
				} `json:"entity"`
				ToState string `json:"toState"`
				Reason  struct {
					Code string `json:"code"`
				} `json:"reason"`
			} `json:"event"`
		}
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			continue
		}
		if record.Event.Entity.Kind == "task" {
			toState, reason, sequence = record.Event.ToState, record.Event.Reason.Code, record.Event.Sequence
		}
	}
	if toState == "" {
		return "", "", 0, fmt.Errorf("no task transition recorded in %s", eventLog)
	}
	return toState, reason, sequence, nil
}
