package console

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// P50 CONSOLE-MULTITASK-ASSEMBLY -- production-entry coverage of the serial
// driver.
//
// The offline harness in runtime_multitask_agent_test.go re-implements the drive
// loop, so it proves the chain shape but not the production driver. These tests
// call ExecuteAgentRunnerRun itself: the definition plan, the record assembly,
// the run-root marker, the replay store, the engine, the event log, the serial
// driver and the inter-task stop proof are production code. Only chain ports are
// stubbed (the pinned CLI cannot be launched where publication durability is
// unproven), which is exactly what AgentRunnerAssemblyOptions.testStub allows.

// assemblyDriverStub serves the chain ports the production entry hands to the
// engine. The dispatch port still travels through the production
// pre-manifest wrapper, so the driver's dispatch observation is real.
type assemblyDriverStub struct {
	t       *testing.T
	records map[string]adapters.AgentRunnerRequestRecord
	facts   map[string]chain.AgentRunnerFrozenFacts

	dispatched     []string
	failByTask     map[string]error
	publishedNames []string
}

func (stub *assemblyDriverStub) DispatchAgentRunner(_ context.Context, request chain.StepRequest) (chain.StepResult, error) {
	record, ok := stub.records[request.TaskID]
	if !ok {
		return chain.StepResult{}, fmt.Errorf("task %q has no request record", request.TaskID)
	}
	if request.AgentRunnerFacts == nil {
		return chain.StepResult{}, fmt.Errorf("task %q was not projected as an AgentRunner step", request.TaskID)
	}
	if request.AgentRunnerFacts.RequestID != record.Request.RequestID {
		stub.t.Fatalf("task %q received another record's facts: %+v", request.TaskID, request.AgentRunnerFacts)
	}
	stub.dispatched = append(stub.dispatched, request.TaskID)
	if err := stub.failByTask[request.TaskID]; err != nil {
		return chain.StepResult{}, err
	}
	return chain.StepResult{Evidence: []string{record.RecordHash}}, nil
}

func (stub *assemblyDriverStub) AdmitAgentRunner(_ context.Context, request chain.StepRequest) error {
	if _, ok := stub.records[request.TaskID]; !ok {
		return fmt.Errorf("task %q has no request record", request.TaskID)
	}
	return nil
}

func (stub *assemblyDriverStub) Accept(_ context.Context, _, taskID string, _ int, _ chain.SnapshotRef, _ chain.Workspace) (chain.CandidateResult, error) {
	return chain.CandidateResult{
		Snapshot:          chain.SnapshotRef{Hash: multiTaskTestDigest("candidate", taskID)},
		EvidenceRootHash:  multiTaskTestDigest("evidence-root", taskID),
		Evidence:          []string{multiTaskTestDigest("candidate-evidence", taskID)},
		CandidateProducer: evidence.Actor{Type: "agent", ID: "assembly-driver-stub"},
	}, nil
}

func (stub *assemblyDriverStub) Review(_ context.Context, request chain.ReviewRequest) (chain.ReviewDecision, error) {
	return chain.ReviewDecision{
		ReceiptID:             "review-" + request.TaskID,
		OccurredAt:            multiTaskClock,
		RecordedBy:            evidence.Actor{Type: "policy", ID: "assembly-driver-policy"},
		ReviewMode:            "policy",
		PolicyHash:            multiTaskTestDigest("review-policy", request.TaskID),
		Outcome:               "approve",
		CandidateSnapshotHash: request.CandidateSnapshotHash,
		Evidence:              []string{multiTaskTestDigest("review-evidence", request.TaskID)},
	}, nil
}

func (stub *assemblyDriverStub) Publish(_ context.Context, request chain.PromotionRequest) (chain.PromotionDecision, error) {
	return chain.PromotionDecision{
		ReceiptID:            "promotion-" + request.TaskID,
		OccurredAt:           multiTaskClock,
		Outcome:              "completed",
		AcceptedSnapshotHash: request.CandidateSnapshotHash,
		WriterStopEvidence:   []string{multiTaskTestDigest("writer-stop", request.TaskID)},
		LeaseEvidence:        []string{multiTaskTestDigest("lease", request.TaskID)},
		Evidence:             []string{multiTaskTestDigest("promotion-evidence", request.TaskID)},
	}, nil
}

func (stub *assemblyDriverStub) RunPostflight(_ context.Context, request chain.PostflightRequest) (chain.PostflightDecision, error) {
	return chain.PostflightDecision{Outcome: chain.PostflightPassed, Evidence: multiTaskFactsEvidence(request.Facts)}, nil
}

// WaitAgentRunnerDispatched answers the (A') wait with the frozen facts of the
// dispatched task, exactly the shape the adapters carrier returns for a
// completed launch.
func (stub *assemblyDriverStub) WaitAgentRunnerDispatched(_ context.Context, request adapters.AgentRunnerRunRequest, _ []byte) (adapters.AgentRunnerTerminalOutcome, error) {
	facts := adapterFrozenFacts(stub.facts[request.TaskID])
	return adapters.AgentRunnerTerminalOutcome{
		Completion: adapters.AgentRunnerCompletion{
			CompletionID: "completion-" + request.Launch.RequestID,
			RequestID:    request.Launch.RequestID,
			RequestHash:  request.RequestHash,
			RunID:        request.Launch.RunID,
			TaskID:       request.TaskID,
			StepID:       request.StepID,
			Attempt:      request.Attempt,
			Status:       "completed",
		},
		Facts: &facts,
	}, nil
}

func (stub *assemblyDriverStub) publishTerminal(_ context.Context, outcome adapters.AgentRunnerTerminalOutcome) error {
	stub.publishedNames = append(stub.publishedNames, outcome.Completion.TaskID+"/"+outcome.Completion.RequestID)
	return nil
}

func (stub *assemblyDriverStub) loadTerminal(record adapters.AgentRunnerRequestRecord) (chain.AgentRunnerTerminal, error) {
	return multiTaskTerminal(record, stub.facts[record.Request.TaskID]), nil
}

// adapterFrozenFacts rewrites the chain fact set as the adapters fact set. The
// two types are structurally identical; the console-local seam keeps them apart,
// so the test states the mapping explicitly.
func adapterFrozenFacts(facts chain.AgentRunnerFrozenFacts) adapters.AgentRunnerFrozenFacts {
	return adapters.AgentRunnerFrozenFacts{
		ManifestHash:            facts.ManifestHash,
		DiffHash:                facts.DiffHash,
		LogHash:                 facts.LogHash,
		UsageHash:               facts.UsageHash,
		ProcessStopEvidenceHash: facts.ProcessStopEvidenceHash,
		ProcessPidReportsHash:   facts.ProcessPidReportsHash,
	}
}

// driverRequestRecord builds one persisted request record for an arbitrary run id.
// It mirrors multiTaskRequestRecord, which is pinned to the offline harness run.
func driverRequestRecord(t *testing.T, runID, taskID, stepID, requestID, parentHash string) adapters.AgentRunnerRequestRecord {
	t.Helper()
	record, err := adapters.NewAgentRunnerRequestRecord(adapters.AgentRunnerRequest{
		RequestID:            requestID,
		RunID:                runID,
		TaskID:               taskID,
		StepID:               stepID,
		Attempt:              1,
		CreatedAt:            multiTaskClock,
		AdapterID:            "assembly-driver-stub",
		Mode:                 "create",
		WorkspaceHash:        multiTaskTestDigest("workspace", taskID),
		ContextHash:          multiTaskTestDigest("context", taskID),
		ParentSnapshotHash:   parentHash,
		AuthorizationHash:    multiTaskTestDigest("authorization", taskID),
		BudgetHash:           multiTaskTestDigest("budget", taskID),
		AllowedTargets:       []string{"workspace-files"},
		AllowedEffects:       []string{"workspace-write", "local-process"},
		EffectMappingVersion: assemblyEffectMappingVersion,
		EffectMappingHash:    assemblyEffectMappingHash,
		Evidence:             []string{multiTaskTestDigest("request-evidence", taskID)},
	})
	if err != nil {
		t.Fatal(err)
	}
	return record
}

// newAssemblyDriverFixture builds a two-task production entry wired to the stub
// port set and returns the evidence directory of every task's requestId.
func newAssemblyDriverFixture(t *testing.T) (AgentRunnerAssemblyOptions, *assemblyDriverStub, map[string]string) {
	t.Helper()
	fixture := newAssemblyFixture(t, assemblyKnobs{})
	runID := fixture.options.RunID

	// The chain captures the run baseline on the first task's isolated workspace,
	// and the intent preparer refuses a record whose declared parent differs from
	// the hash it is handed. The declaration therefore has to be the production
	// baseline digest of that same workspace.
	firstWorkspace := agentRunnerWorkspacePath(fixture.runRoot, multiTaskAgentTaskOne, 1)
	if err := os.MkdirAll(firstWorkspace, 0o755); err != nil {
		t.Fatal(err)
	}
	baselineStore, err := snapshot.NewStore(filepath.Join(fixture.runRoot, "baseline-probe"), 0)
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := adapters.AgentRunnerWorkspaceBaselineDigest(context.Background(), baselineStore, firstWorkspace, runID, multiTaskAgentTaskOne, 1)
	if err != nil {
		t.Fatal(err)
	}

	records := []adapters.AgentRunnerRequestRecord{
		driverRequestRecord(t, runID, multiTaskAgentTaskOne, multiTaskAgentStepOne, "request-one", baseline),
		driverRequestRecord(t, runID, multiTaskAgentTaskTwo, multiTaskAgentStepTwo, "request-two", multiTaskTestDigest("candidate", multiTaskAgentTaskOne)),
	}
	stub := &assemblyDriverStub{
		t:          t,
		records:    map[string]adapters.AgentRunnerRequestRecord{},
		facts:      map[string]chain.AgentRunnerFrozenFacts{},
		failByTask: map[string]error{},
	}
	registry := map[string]adapters.AgentRunnerRequestRecord{}
	for _, record := range records {
		stub.records[record.Request.TaskID] = record
		stub.facts[record.Request.TaskID] = multiTaskFacts(record.Request.TaskID)
		registry[agentRunnerRequestRecordKey(record.Request.TaskID, record.Request.StepID)] = record
	}

	options := fixture.options
	options.Definition = multiTaskIsolatedDefinition()
	options.RequestRecord = adapters.AgentRunnerRequestRecord{}
	options.RequestRecords = registry
	options.testStub = &agentRunnerAssemblyStub{
		Dispatcher:      stub,
		Admission:       stub,
		Acceptance:      stub,
		Reviewer:        stub,
		Publisher:       stub,
		Postflight:      stub,
		Waiter:          stub,
		PublishTerminal: stub.publishTerminal,
		LoadTerminal:    stub.loadTerminal,
	}

	// The evidence directories of both requestIds are resolved through the same
	// production store API the assembly uses, so the pid reports seeded into them
	// are read by the real inter-task stop proof.
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(fixture.runRoot)); err != nil {
		t.Fatal(err)
	}
	store, err := adapters.NewAgentRunnerReplayStoreForRun(fixture.runRoot, runID)
	if err != nil {
		t.Fatal(err)
	}
	evidenceDirs := map[string]string{}
	for _, record := range records {
		dir, err := store.RunEvidenceDir(record.Request.RequestID)
		if err != nil {
			t.Fatal(err)
		}
		evidenceDirs[record.Request.TaskID] = dir
	}
	return options, stub, evidenceDirs
}

// agentRunnerDeadPID returns the pid of a process that has provably exited, so a
// seeded pid report proves "reported, and gone".
func agentRunnerDeadPID(t *testing.T) int {
	t.Helper()
	dead := exec.Command(os.Args[0], "-test.run=^$")
	if err := dead.Run(); err != nil {
		t.Fatal(err)
	}
	return dead.Process.Pid
}

// seedStopReports writes the two pid reports of one task's evidence directory.
// An empty value writes the "the candidate reported nothing" report, which is
// what a completed run leaves behind; only a test that wants the report MISSING
// skips the call entirely.
func seedStopReports(t *testing.T, dir, parent, child string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, agentRunnerParentPIDReportFileName), []byte(parent), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, agentRunnerChildPIDReportFileName), []byte(child), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestAgentRunnerDriverTwoTaskSerialAssembly drives two tasks through the real
// production entry: the two terminals are published in definition order, both
// tasks end PASSED, and the inter-task stop proof gates the second dispatch.
func TestAgentRunnerDriverTwoTaskSerialAssembly(t *testing.T) {
	t.Run("a proven stop lets the next task start", func(t *testing.T) {
		options, stub, evidenceDirs := newAssemblyDriverFixture(t)
		dead := fmt.Sprintf("%d\n", agentRunnerDeadPID(t))
		for _, taskID := range []string{multiTaskAgentTaskOne, multiTaskAgentTaskTwo} {
			seedStopReports(t, evidenceDirs[taskID], dead, "")
		}

		summary, err := ExecuteAgentRunnerRun(context.Background(), options)
		if err != nil {
			t.Fatalf("the two-task assembly must complete: %v", err)
		}
		if summary.ChainState != "COMPLETED" {
			t.Fatalf("chain state = %s, want COMPLETED (%+v)", summary.ChainState, summary)
		}
		if passed := zeroModelTaskCount(summary, "PASSED"); passed != 2 {
			t.Fatalf("both tasks must end PASSED, got %d: %+v", passed, summary.TaskStateCounts)
		}
		wantOrder := []string{multiTaskAgentTaskOne, multiTaskAgentTaskTwo}
		if len(stub.dispatched) != len(wantOrder) {
			t.Fatalf("both tasks must dispatch exactly once, got %v", stub.dispatched)
		}
		for index, want := range wantOrder {
			if stub.dispatched[index] != want {
				t.Fatalf("dispatch %d = %q, want %q (%v)", index, stub.dispatched[index], want, stub.dispatched)
			}
		}
		wantTerminals := []string{
			multiTaskAgentTaskOne + "/request-one",
			multiTaskAgentTaskTwo + "/request-two",
		}
		if len(stub.publishedNames) != len(wantTerminals) {
			t.Fatalf("both tasks must publish their own terminal, got %v", stub.publishedNames)
		}
		for index, want := range wantTerminals {
			if stub.publishedNames[index] != want {
				t.Fatalf("terminal %d = %q, want %q (%v)", index, stub.publishedNames[index], want, stub.publishedNames)
			}
		}
	})

	t.Run("a live launch blocks the next task", func(t *testing.T) {
		options, stub, evidenceDirs := newAssemblyDriverFixture(t)
		dead := fmt.Sprintf("%d\n", agentRunnerDeadPID(t))
		// Task one reports the test process itself: a proven-alive managed process.
		seedStopReports(t, evidenceDirs[multiTaskAgentTaskOne], fmt.Sprintf("%d\n", os.Getpid()), "")
		seedStopReports(t, evidenceDirs[multiTaskAgentTaskTwo], dead, "")

		summary, err := ExecuteAgentRunnerRun(context.Background(), options)
		if err == nil {
			t.Fatalf("a live managed process must block the next dispatch: %+v", summary)
		}
		if !errors.Is(err, ErrInvalidConfig) || !strings.Contains(err.Error(), "still alive") {
			t.Fatalf("the barrier refusal must be the console config sentinel, got %v", err)
		}
		if len(stub.dispatched) != 1 || stub.dispatched[0] != multiTaskAgentTaskOne {
			t.Fatalf("the second task must never dispatch, got %v", stub.dispatched)
		}
	})

	t.Run("a missing stop report fails closed", func(t *testing.T) {
		options, stub, _ := newAssemblyDriverFixture(t)
		// Task one leaves no pid report at all, so the stop cannot be proven.
		summary, err := ExecuteAgentRunnerRun(context.Background(), options)
		if err == nil {
			t.Fatalf("an unprovable stop must fail closed: %+v", summary)
		}
		if !errors.Is(err, adapters.ErrInvalidAgentRunnerEvidence) || !strings.Contains(err.Error(), "cannot be proven gone") {
			t.Fatalf("the barrier refusal must carry the fail-closed evidence sentinel, got %v", err)
		}
		if len(stub.dispatched) != 1 || stub.dispatched[0] != multiTaskAgentTaskOne {
			t.Fatalf("the second task must never dispatch, got %v", stub.dispatched)
		}
	})
}

// assemblyDriverReceipts answers the launch-receipt lookup of the mitigation and
// pins the requestId it is asked about.
type assemblyDriverReceipts struct {
	t    *testing.T
	want string
	seen []string
}

func (reader *assemblyDriverReceipts) LaunchReceipt(requestID string) (adapters.AgentRunnerLaunchReceiptRecord, bool, error) {
	reader.seen = append(reader.seen, requestID)
	if requestID != reader.want {
		reader.t.Errorf("the mitigation asked for requestId %q, want %q", requestID, reader.want)
		return adapters.AgentRunnerLaunchReceiptRecord{}, false, nil
	}
	return adapters.AgentRunnerLaunchReceiptRecord{Receipt: adapters.AgentRunnerLaunchReceipt{LaunchID: "launch-" + requestID}}, true, nil
}

// assemblyDriverStopper records the launch ids the mitigation tried to stop.
type assemblyDriverStopper struct {
	calls []string
}

func (stopper *assemblyDriverStopper) StopAgentRunnerProcess(_ context.Context, launchID string) (guard.TerminationEvidence, error) {
	stopper.calls = append(stopper.calls, launchID)
	return guard.TerminationEvidence{Outcome: "stopped"}, nil
}

// TestAgentRunnerDriverStopsTheLaunchOfTheTaskThatJustFailed is the H1
// regression: a task that fails inside engine.Run after its dispatch has been
// observed must have ITS OWN launch reclaimed. Binding the mitigation to the
// last active task instead would read the previous task's receipt and leave the
// failing task's live process running.
func TestAgentRunnerDriverStopsTheLaunchOfTheTaskThatJustFailed(t *testing.T) {
	options, stub, evidenceDirs := newAssemblyDriverFixture(t)
	seedStopReports(t, evidenceDirs[multiTaskAgentTaskOne], fmt.Sprintf("%d\n", agentRunnerDeadPID(t)), "")

	dispatchFailure := errors.New("task two launch refused after its receipt was written")
	stub.failByTask[multiTaskAgentTaskTwo] = dispatchFailure

	receipts := &assemblyDriverReceipts{t: t, want: "request-two"}
	stopper := &assemblyDriverStopper{}
	options.testStub.Receipts = receipts
	options.testStub.Stopper = stopper

	summary, err := ExecuteAgentRunnerRun(context.Background(), options)
	if err == nil {
		t.Fatalf("the failing second dispatch must fail the run: %+v", summary)
	}
	if !errors.Is(err, dispatchFailure) {
		t.Fatalf("the original failure must stay on the chain, got %v", err)
	}
	if len(stub.dispatched) != 2 || stub.dispatched[1] != multiTaskAgentTaskTwo {
		t.Fatalf("task two must have dispatched before failing, got %v", stub.dispatched)
	}
	if len(receipts.seen) != 1 || receipts.seen[0] != "request-two" {
		t.Fatalf("the mitigation must reclaim task two, saw %v", receipts.seen)
	}
	if len(stopper.calls) != 1 || stopper.calls[0] != "launch-request-two" {
		t.Fatalf("task two's launch must be stopped, got %v", stopper.calls)
	}
}

// TestAgentRunnerAssemblyRefusesNonInitialAttempt pins the M1 refusal: a request
// record declared for another attempt is rejected before the run-root marker
// exists, so no engine is built and no summary is projected.
func TestAgentRunnerAssemblyRefusesNonInitialAttempt(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{})
	body := fixture.options.RequestRecord.Request
	body.Attempt = 2
	record, err := adapters.NewAgentRunnerRequestRecord(body)
	if err != nil {
		t.Fatalf("the attempt-2 record must be a well-formed record: %v", err)
	}
	fixture.options.RequestRecord = record

	summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil {
		t.Fatalf("an attempt-2 record must be refused: %+v", summary)
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("the refusal must be the console config sentinel, got %v", err)
	}
	if !strings.Contains(err.Error(), "attempt 2") {
		t.Fatalf("the refusal must name the refused attempt, got %v", err)
	}
	if summary.ChainState != "" || len(summary.TaskStateCounts) != 0 || len(summary.StepStateCounts) != 0 {
		t.Fatalf("a precondition refusal must project an empty summary, got %+v", summary)
	}
	if _, statErr := os.Stat(runEventLogPath(fixture.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a precondition refusal must not create the run-root marker (stat err=%v)", statErr)
	}
}
