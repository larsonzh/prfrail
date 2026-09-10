package chain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func validOperatorInteractionRequest() OperatorInteractionRequest {
	return OperatorInteractionRequest{
		InteractionID:    "interaction-one",
		RequestedAt:      "2026-09-10T01:00:00.000Z",
		RequestedBy:      evidence.Actor{Type: "adapter", ID: "agent-one"},
		Operator:         evidence.Actor{Type: "operator", ID: "operator-one"},
		RunID:            "run-one",
		TaskID:           "task-one",
		StepID:           "step-one",
		Attempt:          2,
		WorkspaceHash:    baselineHash,
		ConversationID:   "conversation-one",
		RequestID:        "request-one",
		ContextHash:      acceptedOne,
		Question:         "Choose the supported migration target.",
		Reason:           "The task definition permits two targets.",
		AllowedResponses: []string{"use-java-21", "use-java-25"},
		Risk:             "The selected target changes the compatibility baseline.",
		RequiredAction:   "Select one target.",
		Evidence:         []string{stopHash},
	}
}

func TestOperatorInteractionLedgerAcceptsOneBoundResponse(t *testing.T) {
	requestRecord, err := NewOperatorInteractionRequestRecord(validOperatorInteractionRequest(), fixedClock, func() string { return "request-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	response := OperatorInteractionResponse{
		InteractionID:  requestRecord.Request.InteractionID,
		RespondedAt:    "2026-09-10T01:01:00.000Z",
		RespondedBy:    evidence.Actor{Type: "operator", ID: "operator-one"},
		RunID:          "run-one",
		TaskID:         "task-one",
		StepID:         "step-one",
		Attempt:        2,
		WorkspaceHash:  baselineHash,
		ConversationID: "conversation-one",
		RequestID:      "request-two",
		ContextHash:    acceptedOne,
		RequestHash:    requestRecord.RecordHash,
		Selection:      "use-java-25",
		Evidence:       []string{acceptedTwo},
	}
	responseRecord, err := NewOperatorInteractionResponseRecord(requestRecord, response, fixedClock, func() string { return "response-record-one" })
	if err != nil {
		t.Fatal(err)
	}

	ledger := OperatorInteractionLedger{}
	if err := ledger.Append(requestRecord); err != nil {
		t.Fatal(err)
	}
	if err := ledger.Append(responseRecord); err != nil {
		t.Fatal(err)
	}
	if pending := ledger.Pending(); len(pending) != 0 {
		t.Fatalf("answered request remains pending: %+v", pending)
	}
	if err := ledger.Append(responseRecord); !errors.Is(err, ErrDuplicateInteractionResponse) {
		t.Fatalf("expected duplicate response rejection, got %v", err)
	}
}

func TestOperatorInteractionResponseRejectsBrokenBinding(t *testing.T) {
	requestRecord, err := NewOperatorInteractionRequestRecord(validOperatorInteractionRequest(), fixedClock, func() string { return "request-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	base := OperatorInteractionResponse{
		InteractionID:  "interaction-one",
		RespondedAt:    "2026-09-10T01:01:00.000Z",
		RespondedBy:    evidence.Actor{Type: "operator", ID: "operator-one"},
		RunID:          "run-one",
		TaskID:         "task-one",
		StepID:         "step-one",
		Attempt:        2,
		WorkspaceHash:  baselineHash,
		ConversationID: "conversation-one",
		RequestID:      "request-two",
		ContextHash:    acceptedOne,
		RequestHash:    requestRecord.RecordHash,
		Selection:      "use-java-21",
		Evidence:       []string{acceptedTwo},
	}
	tests := map[string]func(*OperatorInteractionResponse){
		"wrong operator":  func(response *OperatorInteractionResponse) { response.RespondedBy.ID = "operator-two" },
		"stale attempt":   func(response *OperatorInteractionResponse) { response.Attempt = 1 },
		"wrong context":   func(response *OperatorInteractionResponse) { response.ContextHash = acceptedThree },
		"same request id": func(response *OperatorInteractionResponse) { response.RequestID = "request-one" },
		"free text":       func(response *OperatorInteractionResponse) { response.Selection = "something else" },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			response := base
			mutate(&response)
			if _, err := NewOperatorInteractionResponseRecord(requestRecord, response, fixedClock, func() string { return "response-record-one" }); !errors.Is(err, ErrInvalidOperatorInteraction) {
				t.Fatalf("expected binding rejection, got %v", err)
			}
		})
	}
}

func TestOperatorInteractionTransitions(t *testing.T) {
	taskState, stepState, err := OpenOperatorInteractionTransition("STEPS_RUNNING", "RUNNING")
	if err != nil || taskState != "WAITING_FOR_OPERATOR" || stepState != "WAITING_FOR_OPERATOR" {
		t.Fatalf("open transition mismatch: task=%s step=%s err=%v", taskState, stepState, err)
	}
	taskState, stepState, err = ResumeOperatorInteractionTransition(taskState, stepState)
	if err != nil || taskState != "STEPS_RUNNING" || stepState != "RUNNING" {
		t.Fatalf("resume transition mismatch: task=%s step=%s err=%v", taskState, stepState, err)
	}
	if _, _, err := ResumeOperatorInteractionTransition("STEPS_RUNNING", "RUNNING"); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected invalid resume rejection, got %v", err)
	}
}

func TestOperatorInteractionRecordStrictDecodeAndTamperDetection(t *testing.T) {
	record, err := NewOperatorInteractionRequestRecord(validOperatorInteractionRequest(), fixedClock, func() string { return "request-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeOperatorInteractionRecord(payload)
	if err != nil || decoded.RecordHash != record.RecordHash {
		t.Fatalf("round trip failed: record=%+v err=%v", decoded, err)
	}
	var wire map[string]any
	if err := json.Unmarshal(payload, &wire); err != nil {
		t.Fatal(err)
	}
	interaction := wire["interaction"].(map[string]any)
	interaction["question"] = "tampered"
	tampered, _ := json.Marshal(wire)
	if _, err := DecodeOperatorInteractionRecord(tampered); !errors.Is(err, ErrInvalidOperatorInteraction) {
		t.Fatalf("expected tamper rejection, got %v", err)
	}
	unknown := append(payload[:len(payload)-1], []byte(`,"unknown":true}`)...)
	if _, err := DecodeOperatorInteractionRecord(unknown); err == nil {
		t.Fatal("expected unknown field rejection")
	}
}

func TestPrepareOperatorInteractionResumeRevalidatesCurrentBinding(t *testing.T) {
	requestRecord, err := NewOperatorInteractionRequestRecord(validOperatorInteractionRequest(), fixedClock, func() string { return "request-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	request := requestRecord.Request
	responseRecord, err := NewOperatorInteractionResponseRecord(requestRecord, OperatorInteractionResponse{
		InteractionID: request.InteractionID, RespondedAt: "2026-09-10T01:01:00.000Z",
		RespondedBy: request.Operator, RunID: request.RunID, TaskID: request.TaskID, StepID: request.StepID,
		Attempt: request.Attempt, WorkspaceHash: request.WorkspaceHash, ConversationID: request.ConversationID,
		RequestID: "request-two", ContextHash: request.ContextHash, RequestHash: requestRecord.RecordHash,
		Selection: "use-java-25", Evidence: []string{acceptedTwo},
	}, fixedClock, func() string { return "response-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	binding := OperatorInteractionBinding{
		RunID: request.RunID, TaskID: request.TaskID, StepID: request.StepID, Attempt: request.Attempt,
		WorkspaceHash: request.WorkspaceHash, ConversationID: request.ConversationID, ContextHash: request.ContextHash,
	}
	command, err := PrepareOperatorInteractionResume(requestRecord, responseRecord, binding, "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR")
	if err != nil {
		t.Fatal(err)
	}
	if command.ConversationID != request.ConversationID || command.RequestID != "request-two" || command.ResponseHash != responseRecord.RecordHash {
		t.Fatalf("unexpected resume command: %+v", command)
	}

	stale := binding
	stale.Attempt++
	if _, err := PrepareOperatorInteractionResume(requestRecord, responseRecord, stale, "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR"); !errors.Is(err, ErrInvalidOperatorInteraction) {
		t.Fatalf("expected stale binding rejection, got %v", err)
	}
	if _, err := PrepareOperatorInteractionResume(requestRecord, responseRecord, binding, "STEPS_RUNNING", "RUNNING"); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected state rejection, got %v", err)
	}
}

func TestEngineOperatorInteractionControlPersistsAndReplays(t *testing.T) {
	store := &memoryEvents{}
	engine := newOperatorInteractionEngine(t, store)
	requestRecord, responseRecord, binding := engineOperatorInteractionRecords(t)

	if err := engine.OpenOperatorInteraction(context.Background(), requestRecord); err != nil {
		t.Fatal(err)
	}
	projection := engine.Projection()
	if projection.ChainState != "PAUSED" || projection.TaskStates[taskKey("task-one", 1)] != "WAITING_FOR_OPERATOR" || projection.StepStates[stepKey("task-one", "code-managed", 1)] != "WAITING_FOR_OPERATOR" {
		t.Fatalf("interaction did not take control: %+v", projection)
	}
	command, err := engine.ResumeOperatorInteraction(context.Background(), requestRecord, responseRecord, binding)
	if err != nil {
		t.Fatal(err)
	}
	projection = engine.Projection()
	if command.RequestID != "request-two" || projection.ChainState != "PAUSED" || projection.TaskStates[taskKey("task-one", 1)] != "STEPS_RUNNING" || projection.StepStates[stepKey("task-one", "code-managed", 1)] != "RUNNING" {
		t.Fatalf("interaction did not return control safely: command=%+v projection=%+v", command, projection)
	}

	options, _, _, _, _, _, _, _, _ := testOptions(store)
	sequence := len(store.events)
	options.IDs = func() string { sequence++; return fmt.Sprintf("event-%d", sequence) }
	restarted, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if got := restarted.Projection(); got.ChainState != "PAUSED" || got.TaskStates[taskKey("task-one", 1)] != "STEPS_RUNNING" || got.StepStates[stepKey("task-one", "code-managed", 1)] != "RUNNING" {
		t.Fatalf("replay lost resumed state: %+v", got)
	}
	replayedCommand, err := restarted.ResumeOperatorInteraction(context.Background(), requestRecord, responseRecord, binding)
	if err != nil || replayedCommand != command {
		t.Fatalf("restart did not rebuild the same resume command: command=%+v err=%v", replayedCommand, err)
	}
}

func TestEngineOperatorInteractionWriteFailureDoesNotReleaseControl(t *testing.T) {
	store := &memoryEvents{}
	engine := newOperatorInteractionEngine(t, store)
	requestRecord, responseRecord, binding := engineOperatorInteractionRecords(t)

	store.err = errors.New("disk full")
	if err := engine.OpenOperatorInteraction(context.Background(), requestRecord); err == nil {
		t.Fatal("expected open persistence failure")
	}
	if engine.Projection().ChainState != "RUNNING" {
		t.Fatal("failed pause append changed projection")
	}
	store.err = nil
	if err := engine.OpenOperatorInteraction(context.Background(), requestRecord); err != nil {
		t.Fatal(err)
	}

	store.err = errors.New("disk full")
	command, err := engine.ResumeOperatorInteraction(context.Background(), requestRecord, responseRecord, binding)
	if err == nil || command != (OperatorResumeCommand{}) {
		t.Fatalf("failed resume returned executable command: command=%+v err=%v", command, err)
	}
	projection := engine.Projection()
	if projection.TaskStates[taskKey("task-one", 1)] != "WAITING_FOR_OPERATOR" || projection.StepStates[stepKey("task-one", "code-managed", 1)] != "WAITING_FOR_OPERATOR" {
		t.Fatalf("failed resume released control: %+v", projection)
	}
	store.err = nil
	if _, err := engine.ResumeOperatorInteraction(context.Background(), requestRecord, responseRecord, binding); err != nil {
		t.Fatal(err)
	}
}

type failOnceEvents struct {
	memoryEvents
	appendCalls int
	failAt      int
}

func (store *failOnceEvents) Append(ctx context.Context, event evidence.StateEvent) error {
	store.appendCalls++
	if store.appendCalls == store.failAt {
		return errors.New("injected append failure")
	}
	return store.memoryEvents.Append(ctx, event)
}

func TestEngineOperatorInteractionResumeConvergesAfterPartialWrite(t *testing.T) {
	store := &failOnceEvents{}
	engine := newOperatorInteractionEngine(t, store)
	requestRecord, responseRecord, binding := engineOperatorInteractionRecords(t)
	if err := engine.OpenOperatorInteraction(context.Background(), requestRecord); err != nil {
		t.Fatal(err)
	}
	store.failAt = store.appendCalls + 2
	command, err := engine.ResumeOperatorInteraction(context.Background(), requestRecord, responseRecord, binding)
	if err == nil || command != (OperatorResumeCommand{}) {
		t.Fatalf("partial resume returned executable command: command=%+v err=%v", command, err)
	}
	projection := engine.Projection()
	if projection.TaskStates[taskKey("task-one", 1)] != "WAITING_FOR_OPERATOR" || projection.StepStates[stepKey("task-one", "code-managed", 1)] != "RUNNING" {
		t.Fatalf("unexpected partial state: %+v", projection)
	}
	request := requestRecord.Request
	alternateResponse, err := NewOperatorInteractionResponseRecord(requestRecord, OperatorInteractionResponse{
		InteractionID: request.InteractionID, RespondedAt: "2026-09-10T01:02:00.000Z", RespondedBy: request.Operator,
		RunID: request.RunID, TaskID: request.TaskID, StepID: request.StepID, Attempt: request.Attempt,
		WorkspaceHash: request.WorkspaceHash, ConversationID: request.ConversationID, RequestID: "request-three",
		ContextHash: request.ContextHash, RequestHash: requestRecord.RecordHash, Selection: "use-java-21", Evidence: []string{acceptedThree},
	}, fixedClock, func() string { return "response-record-two" })
	if err != nil {
		t.Fatal(err)
	}
	if _, err := engine.ResumeOperatorInteraction(context.Background(), requestRecord, alternateResponse, binding); !errors.Is(err, ErrInvalidOperatorInteraction) {
		t.Fatalf("expected alternate response rejection after partial write, got %v", err)
	}
	command, err = engine.ResumeOperatorInteraction(context.Background(), requestRecord, responseRecord, binding)
	if err != nil || command.RequestID != "request-two" {
		t.Fatalf("resume did not converge: command=%+v err=%v", command, err)
	}
}

func TestEngineOperatorInteractionRejectsForeignRunAndStaleState(t *testing.T) {
	store := &memoryEvents{}
	engine := newOperatorInteractionEngine(t, store)
	requestRecord, responseRecord, binding := engineOperatorInteractionRecords(t)

	foreign := requestRecord
	foreignRequest := *requestRecord.Request
	foreignRequest.RunID = "run-two"
	foreign, _ = NewOperatorInteractionRequestRecord(foreignRequest, fixedClock, func() string { return foreignRequest.RecordID })
	if err := engine.OpenOperatorInteraction(context.Background(), foreign); !errors.Is(err, ErrInvalidOperatorInteraction) {
		t.Fatalf("expected foreign run rejection, got %v", err)
	}
	if _, err := engine.ResumeOperatorInteraction(context.Background(), requestRecord, responseRecord, binding); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected stale state rejection, got %v", err)
	}
}

func newOperatorInteractionEngine(t *testing.T, store EventStore) *Engine {
	t.Helper()
	options, _, _, _, _, _, _, _, _ := testOptions(store)
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	transitions := []struct {
		entity evidence.Entity
		state  string
	}{
		{chainEntity("run-one"), "BASELINED"}, {chainEntity("run-one"), "RUNNING"},
		{taskEntity("run-one", "task-one"), "PENDING"}, {taskEntity("run-one", "task-one"), "PRECHECK"},
		{taskEntity("run-one", "task-one"), "STEPS_RUNNING"}, {stepEntity("run-one", "task-one", "code-managed"), "PENDING"},
		{stepEntity("run-one", "task-one", "code-managed"), "RUNNING"},
	}
	for _, transition := range transitions {
		inputs := []string(nil)
		if transition.state == "BASELINED" {
			inputs = []string{baselineHash}
		}
		if err := engine.transition(context.Background(), transition.entity, transition.state, inputs, "test-setup"); err != nil {
			t.Fatal(err)
		}
	}
	return engine
}

func engineOperatorInteractionRecords(t *testing.T) (OperatorInteractionRecord, OperatorInteractionRecord, OperatorInteractionBinding) {
	t.Helper()
	request := validOperatorInteractionRequest()
	request.Attempt = 1
	request.StepID = "code-managed"
	requestRecord, err := NewOperatorInteractionRequestRecord(request, fixedClock, func() string { return "request-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	responseRecord, err := NewOperatorInteractionResponseRecord(requestRecord, OperatorInteractionResponse{
		InteractionID: request.InteractionID, RespondedAt: "2026-09-10T01:01:00.000Z", RespondedBy: request.Operator,
		RunID: request.RunID, TaskID: request.TaskID, StepID: request.StepID, Attempt: request.Attempt,
		WorkspaceHash: request.WorkspaceHash, ConversationID: request.ConversationID, RequestID: "request-two",
		ContextHash: request.ContextHash, RequestHash: requestRecord.RecordHash, Selection: "use-java-25", Evidence: []string{acceptedTwo},
	}, fixedClock, func() string { return "response-record-one" })
	if err != nil {
		t.Fatal(err)
	}
	binding := OperatorInteractionBinding{
		RunID: request.RunID, TaskID: request.TaskID, StepID: request.StepID, Attempt: request.Attempt,
		WorkspaceHash: request.WorkspaceHash, ConversationID: request.ConversationID, ContextHash: request.ContextHash,
	}
	return requestRecord, responseRecord, binding
}
