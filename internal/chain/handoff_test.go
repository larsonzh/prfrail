package chain

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func validHandoffSession() HandoffSession {
	return HandoffSession{
		RunID:             "run-one",
		TaskID:            "task-one",
		StepID:            "step-one",
		Attempt:           1,
		HandoffPolicyHash: baselineHash,
		AllowedReturnActions: []string{
			"complete",
			"abort",
			"request-agent",
		},
		Operator:           evidence.Actor{Type: "operator", ID: "operator-one"},
		OpenedAt:           "2026-09-07T01:00:00.000Z",
		BeforeManifestHash: acceptedOne,
		LeaseEvidence:      []string{stopHash},
	}
}

func TestBuildHandoffReceiptRecordComplete(t *testing.T) {
	session := validHandoffSession()
	decision := HandoffDecision{
		ReceiptID:          "handoff-one",
		ClosedAt:           "2026-09-07T01:01:00.000Z",
		AfterManifestHash:  acceptedTwo,
		DiffHash:           acceptedThree,
		Outcome:            "complete",
		HookResultEvidence: []string{stopHash},
		Evidence:           []string{baselineHash},
	}
	record, err := BuildHandoffReceiptRecord(session, decision, fixedClock, func() string { return "handoff-id" })
	if err != nil {
		t.Fatal(err)
	}
	if record.SchemaVersion != evidence.SchemaVersion || record.ReceiptHash == "" {
		t.Fatalf("unexpected record: %+v", record)
	}
	if record.Receipt.Outcome != "complete" || len(record.Receipt.HookResultEvidence) != 1 {
		t.Fatalf("unexpected receipt payload: %+v", record.Receipt)
	}
}

func TestBuildHandoffReceiptRecordRejectsOutcomeRules(t *testing.T) {
	session := validHandoffSession()
	invalidComplete := HandoffDecision{
		ReceiptID:         "handoff-one",
		ClosedAt:          "2026-09-07T01:01:00.000Z",
		AfterManifestHash: acceptedTwo,
		DiffHash:          acceptedThree,
		Outcome:           "complete",
		Evidence:          []string{baselineHash},
	}
	if _, err := BuildHandoffReceiptRecord(session, invalidComplete, fixedClock, func() string { return "handoff-id" }); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected complete rejection, got: %v", err)
	}
	invalidAbort := invalidComplete
	invalidAbort.Outcome = "abort"
	invalidAbort.HookResultEvidence = []string{stopHash}
	if _, err := BuildHandoffReceiptRecord(session, invalidAbort, fixedClock, func() string { return "handoff-id" }); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected abort rejection, got: %v", err)
	}
}

func TestBuildHandoffReceiptRecordRejectsDisallowedOutcome(t *testing.T) {
	session := validHandoffSession()
	session.AllowedReturnActions = []string{"complete", "abort"}
	decision := HandoffDecision{
		ReceiptID:         "handoff-one",
		ClosedAt:          "2026-09-07T01:01:00.000Z",
		AfterManifestHash: acceptedTwo,
		DiffHash:          acceptedThree,
		Outcome:           "request-agent",
		Evidence:          []string{baselineHash},
	}
	if _, err := BuildHandoffReceiptRecord(session, decision, fixedClock, func() string { return "handoff-id" }); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected disallowed outcome rejection, got: %v", err)
	}
}

func TestBuildHandoffReceiptRecordRejectsEmptyAllowedReturnActions(t *testing.T) {
	session := validHandoffSession()
	session.AllowedReturnActions = []string{}
	decision := HandoffDecision{
		ReceiptID:         "handoff-one",
		ClosedAt:          "2026-09-07T01:01:00.000Z",
		AfterManifestHash: acceptedTwo,
		DiffHash:          acceptedThree,
		Outcome:           "abort",
		Evidence:          []string{baselineHash},
	}
	if _, err := BuildHandoffReceiptRecord(session, decision, fixedClock, func() string { return "handoff-id" }); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected empty allowed actions rejection, got: %v", err)
	}
}

func TestBuiltHandoffReceiptSerializesToSchemaShape(t *testing.T) {
	session := validHandoffSession()
	decision := HandoffDecision{
		ReceiptID:         "handoff-one",
		ClosedAt:          "2026-09-07T01:01:00.000Z",
		AfterManifestHash: acceptedTwo,
		DiffHash:          acceptedThree,
		Outcome:           "abort",
		Evidence:          []string{baselineHash},
	}
	record, err := BuildHandoffReceiptRecord(session, decision, fixedClock, func() string { return "handoff-id" })
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"schemaVersion", "receipt", "receiptHash"} {
		if _, ok := fields[key]; !ok {
			t.Fatalf("missing top-level key %q", key)
		}
	}
	receipt, ok := fields["receipt"].(map[string]any)
	if !ok {
		t.Fatal("missing receipt object")
	}
	for _, key := range []string{"receiptId", "recordedBy", "operator", "handoffPolicyHash", "openedAt", "closedAt", "beforeManifestHash", "afterManifestHash", "diffHash", "leaseEvidence", "outcome", "hookResultEvidence", "evidence"} {
		if _, ok := receipt[key]; !ok {
			t.Fatalf("missing receipt key %q", key)
		}
	}
	if _, ok := receipt["ReceiptID"]; ok {
		t.Fatal("receipt keys must be lower-case camelCase")
	}
	hookEvidence, ok := receipt["hookResultEvidence"].([]any)
	if !ok || len(hookEvidence) != 0 {
		t.Fatalf("abort hookResultEvidence must serialize as empty array, got: %v", receipt["hookResultEvidence"])
	}
}

func TestHandoffTransitions(t *testing.T) {
	taskState, stepState, err := OpenHandoffTransition("STEPS_RUNNING", "RUNNING")
	if err != nil || taskState != "WAITING_FOR_OPERATOR" || stepState != "WAITING_FOR_OPERATOR" {
		t.Fatalf("open transition mismatch: task=%s step=%s err=%v", taskState, stepState, err)
	}
	cases := []struct {
		outcome  string
		taskWant string
		stepWant string
	}{
		{outcome: "complete", taskWant: "STEPS_RUNNING", stepWant: "RUNNING"},
		{outcome: "abort", taskWant: "FAILED", stepWant: "FAILED"},
		{outcome: "request-agent", taskWant: "FAILED", stepWant: "FAILED"},
	}
	for _, item := range cases {
		t.Run(item.outcome, func(t *testing.T) {
			nextTask, nextStep, err := ReturnHandoffTransition("WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR", item.outcome)
			if err != nil || nextTask != item.taskWant || nextStep != item.stepWant {
				t.Fatalf("return transition mismatch: task=%s step=%s err=%v", nextTask, nextStep, err)
			}
		})
	}
	if _, _, err := ReturnHandoffTransition("STEPS_RUNNING", "RUNNING", "complete"); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected state rejection, got: %v", err)
	}
}
