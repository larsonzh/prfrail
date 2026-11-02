package console

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// P50 CONSOLE-MULTITASK-ASSEMBLY -- offline two-task AgentRunner harness.
//
// The production dispatcher refuses to publish R where publication durability is
// unproven (Windows), so the serial multi-task contract is proven here with an
// offline port set instead: the real chain engine, the real replay store and the
// real adapters intent registry, but stub step/accept/review/publish/postflight
// ports. Nothing spawns a process, reads a clock or touches the network.
//
// What it proves, per CONTRACTS.md A6 多任务串行装配（P50）:
//   - the two tasks are executed in definition order, serially;
//   - every task projects its own request record through the intent registry, so
//     no task can borrow another task's requestId;
//   - the six frozen facts of each task's terminal are recorded against that task
//     only, in their fixed positions, and never leak into the other task;
//   - the request/completion artifact streams are per requestId (distinct hashes,
//     distinct evidence directories, and - where durability is proven - distinct
//     replay states);
//   - the single event log is totally ordered: task one's terminal precedes task
//     two's dispatch.

const (
	multiTaskRunID = "run-multitask-agent"
	multiTaskClock = "2026-11-03T00:00:00.000Z"
)

func multiTaskTestDigest(parts ...string) string {
	return digest("proofrail:p50-multitask-test:1\n", parts...)
}

// multiTaskBaselineHash is the run-level baseline every record declares as its
// parent: the chain captures one baseline per run and materializes every task
// from the accepted parent it records.
func multiTaskBaselineHash() string {
	return multiTaskTestDigest("baseline")
}

// multiTaskRequestRecord builds one persisted request record bound to a task.
func multiTaskRequestRecord(t *testing.T, taskID, stepID, requestID, parentHash string) adapters.AgentRunnerRequestRecord {
	t.Helper()
	record, err := adapters.NewAgentRunnerRequestRecord(adapters.AgentRunnerRequest{
		RequestID:            requestID,
		RunID:                multiTaskRunID,
		TaskID:               taskID,
		StepID:               stepID,
		Attempt:              1,
		CreatedAt:            multiTaskClock,
		AdapterID:            "multitask-stub",
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

func multiTaskFacts(taskID string) chain.AgentRunnerFrozenFacts {
	return chain.AgentRunnerFrozenFacts{
		ManifestHash:            multiTaskTestDigest("manifest", taskID),
		DiffHash:                multiTaskTestDigest("diff", taskID),
		LogHash:                 multiTaskTestDigest("log", taskID),
		UsageHash:               multiTaskTestDigest("usage", taskID),
		ProcessStopEvidenceHash: multiTaskTestDigest("process-stop", taskID),
		ProcessPidReportsHash:   multiTaskTestDigest("pid-reports", taskID),
	}
}

func multiTaskFactsEvidence(facts chain.AgentRunnerFrozenFacts) []string {
	return []string{
		facts.ManifestHash, facts.DiffHash, facts.LogHash,
		facts.UsageHash, facts.ProcessStopEvidenceHash, facts.ProcessPidReportsHash,
	}
}

// multiTaskStubPorts serves every non-AgentRunner chain port offline. The step
// port is the only one that observes the per-task projection.
type multiTaskStubPorts struct {
	t        *testing.T
	runRoot  string
	baseline string
	records  map[string]adapters.AgentRunnerRequestRecord
	facts    map[string]chain.AgentRunnerFrozenFacts

	executed   []string
	parkedTask string
}

func (stub *multiTaskStubPorts) CaptureBaseline(context.Context, string) (chain.SnapshotRef, error) {
	return chain.SnapshotRef{Hash: stub.baseline}, nil
}

func (stub *multiTaskStubPorts) Materialize(_ context.Context, _ string, taskID string, _ int, _ chain.SnapshotRef) (chain.Workspace, error) {
	root := agentRunnerWorkspacePath(stub.runRoot, taskID, 1)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return chain.Workspace{}, err
	}
	return chain.Workspace{Root: root}, nil
}

func (stub *multiTaskStubPorts) Execute(_ context.Context, request chain.StepRequest) (chain.StepResult, error) {
	record, ok := stub.records[request.TaskID]
	if !ok {
		return chain.StepResult{}, fmt.Errorf("task %q has no request record", request.TaskID)
	}
	if request.AgentRunnerFacts == nil {
		return chain.StepResult{}, fmt.Errorf("task %q was not projected as an AgentRunner step", request.TaskID)
	}
	if request.AgentRunnerFacts.RequestID != record.Request.RequestID || request.AgentRunnerFacts.ContextHash != record.Request.ContextHash {
		stub.t.Fatalf("task %q received another record's facts: %+v", request.TaskID, request.AgentRunnerFacts)
	}
	stub.executed = append(stub.executed, request.TaskID+"/"+request.Step.ID)
	stub.parkedTask = request.TaskID
	return chain.StepResult{Evidence: []string{record.RecordHash, multiTaskTestDigest("launch", request.TaskID)}}, nil
}

func (stub *multiTaskStubPorts) Accept(_ context.Context, _, taskID string, _ int, _ chain.SnapshotRef, _ chain.Workspace) (chain.CandidateResult, error) {
	return chain.CandidateResult{
		Snapshot:          chain.SnapshotRef{Hash: multiTaskTestDigest("candidate", taskID)},
		EvidenceRootHash:  multiTaskTestDigest("evidence-root", taskID),
		Evidence:          []string{multiTaskTestDigest("candidate-evidence", taskID)},
		CandidateProducer: evidence.Actor{Type: "agent", ID: "multitask-stub"},
	}, nil
}

func (stub *multiTaskStubPorts) Review(_ context.Context, request chain.ReviewRequest) (chain.ReviewDecision, error) {
	return chain.ReviewDecision{
		ReceiptID:             "review-" + request.TaskID,
		OccurredAt:            multiTaskClock,
		RecordedBy:            evidence.Actor{Type: "policy", ID: "multitask-policy"},
		ReviewMode:            "policy",
		PolicyHash:            multiTaskTestDigest("review-policy", request.TaskID),
		Outcome:               "approve",
		CandidateSnapshotHash: request.CandidateSnapshotHash,
		Evidence:              []string{multiTaskTestDigest("review-evidence", request.TaskID)},
	}, nil
}

func (stub *multiTaskStubPorts) Publish(_ context.Context, request chain.PromotionRequest) (chain.PromotionDecision, error) {
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

func (stub *multiTaskStubPorts) RunPostflight(_ context.Context, request chain.PostflightRequest) (chain.PostflightDecision, error) {
	return chain.PostflightDecision{Outcome: chain.PostflightPassed, Evidence: multiTaskFactsEvidence(request.Facts)}, nil
}

func (stub *multiTaskStubPorts) Stop(context.Context, string) ([]string, error) {
	return []string{multiTaskTestDigest("stop-evidence")}, nil
}

func (stub *multiTaskStubPorts) Reconcile(context.Context, string) error { return nil }

// multiTaskTerminal builds the completed terminal of one task: its own request
// and completion digests plus its own six frozen facts.
func multiTaskTerminal(record adapters.AgentRunnerRequestRecord, facts chain.AgentRunnerFrozenFacts) chain.AgentRunnerTerminal {
	completionHash := multiTaskTestDigest("completion", record.Request.TaskID)
	return chain.AgentRunnerTerminal{
		RequestID:      record.Request.RequestID,
		RequestHash:    record.RecordHash,
		CompletionHash: completionHash,
		RunID:          record.Request.RunID,
		TaskID:         record.Request.TaskID,
		StepID:         record.Request.StepID,
		Attempt:        record.Request.Attempt,
		SessionID:      "session-" + record.Request.TaskID,
		Status:         chain.AgentRunnerTerminalCompleted,
		Facts:          &facts,
		Evidence:       []string{record.RecordHash, completionHash},
	}
}

// multiTaskCompletionRecord builds the C record of one task for the replay store.
func multiTaskCompletionRecord(t *testing.T, record adapters.AgentRunnerRequestRecord, facts chain.AgentRunnerFrozenFacts) adapters.AgentRunnerCompletionRecord {
	t.Helper()
	exit := 0
	outputManifest := facts.ManifestHash
	completion, err := adapters.NewAgentRunnerCompletionRecord(adapters.AgentRunnerCompletion{
		CompletionID:       "completion-" + record.Request.TaskID,
		RequestID:          record.Request.RequestID,
		RequestHash:        record.RecordHash,
		RunID:              record.Request.RunID,
		TaskID:             record.Request.TaskID,
		StepID:             record.Request.StepID,
		Attempt:            record.Request.Attempt,
		AdapterID:          record.Request.AdapterID,
		SessionID:          "session-" + record.Request.TaskID,
		CompletedAt:        multiTaskClock,
		Status:             "completed",
		ExitCode:           &exit,
		ProcessTreeStatus:  "stopped",
		LogsComplete:       true,
		UsageComplete:      true,
		OutputManifestHash: &outputManifest,
		Evidence:           []string{record.RecordHash},
	})
	if err != nil {
		t.Fatal(err)
	}
	return completion
}

// TestAgentRunnerMultiTaskSerialOfflineHarness drives two isolated-workspace
// tasks through one chain engine and one event log.
func TestAgentRunnerMultiTaskSerialOfflineHarness(t *testing.T) {
	runRoot := t.TempDir()
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(runRoot)); err != nil {
		t.Fatal(err)
	}
	store, err := adapters.NewAgentRunnerReplayStoreForRun(runRoot, multiTaskRunID)
	if err != nil {
		t.Fatal(err)
	}

	definition := multiTaskIsolatedDefinition()
	records := []adapters.AgentRunnerRequestRecord{
		multiTaskRequestRecord(t, multiTaskAgentTaskOne, multiTaskAgentStepOne, "request-one", multiTaskBaselineHash()),
		multiTaskRequestRecord(t, multiTaskAgentTaskTwo, multiTaskAgentStepTwo, "request-two", multiTaskTestDigest("candidate", multiTaskAgentTaskOne)),
	}
	registry, err := adapters.NewAgentRunnerStepIntentRegistry(records, definition)
	if err != nil {
		t.Fatalf("the intent registry must bind both records: %v", err)
	}
	if registry.Len() != 2 {
		t.Fatalf("the intent registry must hold two bindings, got %d", registry.Len())
	}

	stub := &multiTaskStubPorts{
		t:        t,
		runRoot:  runRoot,
		baseline: multiTaskBaselineHash(),
		records:  map[string]adapters.AgentRunnerRequestRecord{},
		facts:    map[string]chain.AgentRunnerFrozenFacts{},
	}
	for _, record := range records {
		stub.records[record.Request.TaskID] = record
		stub.facts[record.Request.TaskID] = multiTaskFacts(record.Request.TaskID)
	}

	eventStore, err := chain.NewFileEventStore(runEventLogPath(runRoot))
	if err != nil {
		t.Fatal(err)
	}
	ids := 0
	engine, err := chain.New(context.Background(), chain.Options{
		RunID:       multiTaskRunID,
		Definition:  definition,
		Events:      eventStore,
		Baselines:   stub,
		Workspaces:  stub,
		Steps:       stub,
		StepIntents: registry,
		Acceptance:  stub,
		Reviewer:    stub,
		Publisher:   stub,
		Postflight:  stub,
		Stopper:     stub,
		Reconciler:  stub,
		Clock:       func() time.Time { return time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC) },
		IDs: func() string {
			ids++
			return fmt.Sprintf("multitask-event-%d", ids)
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Drive: the engine parks one task at a time; the harness routes exactly that
	// task's terminal. No task may be terminalised twice.
	ctx := context.Background()
	submitted := map[string]bool{}
	for iteration := 0; iteration < 2*len(definition.Tasks)+2; iteration++ {
		runErr := engine.Run(ctx)
		if runErr == nil {
			break
		}
		if !errors.Is(runErr, chain.ErrAwaitingAgentRunnerTerminal) {
			t.Fatalf("the harness reached an unexpected chain failure: %v", runErr)
		}
		taskID := stub.parkedTask
		if submitted[taskID] {
			t.Fatalf("task %q was terminalised twice", taskID)
		}
		submitted[taskID] = true
		terminal := multiTaskTerminal(stub.records[taskID], stub.facts[taskID])
		if _, err := engine.SubmitAgentRunnerTerminal(ctx, terminal); err != nil {
			t.Fatalf("task %q terminal must route: %v", taskID, err)
		}
	}
	projection := engine.Projection()
	if projection.ChainState != "COMPLETED" {
		t.Fatalf("the two-task chain must complete, got %s", projection.ChainState)
	}

	// 1. Serial definition order.
	wantOrder := []string{"task-one/code-one", "task-two/code-two"}
	if len(stub.executed) != len(wantOrder) {
		t.Fatalf("both tasks must execute exactly once, got %v", stub.executed)
	}
	for index, want := range wantOrder {
		if stub.executed[index] != want {
			t.Fatalf("execution %d = %q, want %q (%v)", index, stub.executed[index], want, stub.executed)
		}
	}

	// 2. Per-task frozen facts and the total order of the single event log.
	events, err := eventStore.Load(ctx)
	if err != nil {
		t.Fatal(err)
	}
	previous := 0
	for _, record := range events {
		if record.Event.Sequence <= previous {
			t.Fatalf("event sequence must increase, got %d after %d", record.Event.Sequence, previous)
		}
		previous = record.Event.Sequence
	}
	stepStart := map[string]int{}
	terminalAt := map[string]int{}
	for _, record := range events {
		event := record.Event
		if event.Entity.TaskID == "" || event.Entity.Attempt != 1 {
			continue
		}
		switch event.Reason.Code {
		case "step-started":
			if _, seen := stepStart[event.Entity.TaskID]; !seen {
				stepStart[event.Entity.TaskID] = event.Sequence
			}
		case "agent-terminal-completed":
			terminalAt[event.Entity.TaskID] = event.Sequence
			facts := stub.facts[event.Entity.TaskID]
			want := append([]string{stub.records[event.Entity.TaskID].RecordHash, multiTaskTestDigest("completion", event.Entity.TaskID)}, multiTaskFactsEvidence(facts)...)
			if len(event.InputEvidence) != len(want) {
				t.Fatalf("task %q terminal evidence = %v, want %v", event.Entity.TaskID, event.InputEvidence, want)
			}
			for index := range want {
				if event.InputEvidence[index] != want[index] {
					t.Fatalf("task %q terminal evidence[%d] = %q, want %q", event.Entity.TaskID, index, event.InputEvidence[index], want[index])
				}
			}
			// No fact of the other task may appear.
			for _, other := range []string{multiTaskAgentTaskOne, multiTaskAgentTaskTwo} {
				if other == event.Entity.TaskID {
					continue
				}
				for _, foreign := range multiTaskFactsEvidence(stub.facts[other]) {
					for _, digest := range event.InputEvidence {
						if digest == foreign {
							t.Fatalf("task %q terminal carries a fact of task %q", event.Entity.TaskID, other)
						}
					}
				}
			}
		}
	}
	if !(stepStart[multiTaskAgentTaskOne] < terminalAt[multiTaskAgentTaskOne] &&
		terminalAt[multiTaskAgentTaskOne] < stepStart[multiTaskAgentTaskTwo] &&
		stepStart[multiTaskAgentTaskTwo] < terminalAt[multiTaskAgentTaskTwo]) {
		t.Fatalf("the tasks must interleave serially in one log: starts=%v terminals=%v", stepStart, terminalAt)
	}

	// 3. R/C artifact streams are per requestId.
	firstDir, err := store.RunEvidenceDir(records[0].Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	secondDir, err := store.RunEvidenceDir(records[1].Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if records[0].Request.RequestID == records[1].Request.RequestID || records[0].RecordHash == records[1].RecordHash {
		t.Fatal("the two tasks must own distinct requestIds and request hashes")
	}
	if firstDir == secondDir {
		t.Fatalf("the two tasks must own distinct evidence directories, got %q twice", firstDir)
	}
	if durability := store.PublicationDurability(); durability == adapters.PublishDurabilityProven {
		for _, record := range records {
			if _, err := store.RecordRequest(record); err != nil {
				t.Fatalf("task %q request publish must succeed: %v", record.Request.TaskID, err)
			}
			if _, err := store.RecordCompletion(record, multiTaskCompletionRecord(t, record, stub.facts[record.Request.TaskID])); err != nil {
				t.Fatalf("task %q completion publish must succeed: %v", record.Request.TaskID, err)
			}
			if state, err := store.State(record.Request.RequestID); err != nil || state != adapters.AgentRunnerReplayStateTerminalReceiptPresent {
				t.Fatalf("task %q replay state = %q err=%v", record.Request.TaskID, state, err)
			}
		}
	} else {
		// Windows cannot publish R (publication durability is unproven), so the
		// per-requestId artifact separation is asserted through the derived paths
		// only; the write path is covered by the Ubuntu CI leg.
		t.Logf("R/C publication skipped: durability is %s on this platform", durabilityOf(store))
	}
}

// durabilityOf names the store's publication durability for the log line above.
func durabilityOf(store *adapters.AgentRunnerReplayStore) adapters.PublishDurability {
	return store.PublicationDurability()
}
