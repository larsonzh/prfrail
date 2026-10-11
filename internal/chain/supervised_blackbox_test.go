package chain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	supervisedWorkspaceHash   = "sha256:0000000000000000000000000000000000000000000000000000000000000041"
	supervisedScanHash        = "sha256:0000000000000000000000000000000000000000000000000000000000000042"
	supervisedStopHash        = "sha256:0000000000000000000000000000000000000000000000000000000000000043"
	supervisedGateHash        = "sha256:0000000000000000000000000000000000000000000000000000000000000044"
	supervisedReportHash      = "sha256:0000000000000000000000000000000000000000000000000000000000000045"
	supervisedReceiptEvidence = "sha256:0000000000000000000000000000000000000000000000000000000000000046"
	supervisedHookHash        = "sha256:0000000000000000000000000000000000000000000000000000000000000047"
	supervisedAfterManifest   = "sha256:0000000000000000000000000000000000000000000000000000000000000048"
	supervisedDiffHash        = "sha256:0000000000000000000000000000000000000000000000000000000000000049"
	supervisedInteractionOne  = "sha256:0000000000000000000000000000000000000000000000000000000000000050"
	supervisedInteractionTwo  = "sha256:0000000000000000000000000000000000000000000000000000000000000051"
	supervisedConfirmEvidence = "sha256:0000000000000000000000000000000000000000000000000000000000000052"
	supervisedReportEvidence  = "sha256:0000000000000000000000000000000000000000000000000000000000000053"
	supervisedSpareHash       = "sha256:0000000000000000000000000000000000000000000000000000000000000054"
)

type supervisedFixture struct {
	engine     *Engine
	store      EventStore
	acceptance *fakeAcceptance
	reviewer   *fakeReviewer
	publisher  *fakePublisher
	postflight *fakeSupervisedPostflight
	session    HandoffSession
	confirm    SupervisedRiskConfirmationRecord
}

type fakeSupervisedPostflight struct {
	decision *SupervisedPostflightDecision
	err      error
	calls    int
	requests []SupervisedPostflightRequest
}

func (fake *fakeSupervisedPostflight) RunSupervisedPostflight(_ context.Context, request SupervisedPostflightRequest) (SupervisedPostflightDecision, error) {
	fake.requests = append(fake.requests, request)
	fake.calls++
	if fake.err != nil {
		return SupervisedPostflightDecision{}, fake.err
	}
	if fake.decision != nil {
		return *fake.decision, nil
	}
	return SupervisedPostflightDecision{Outcome: PostflightPassed, Evidence: supervisedFactsEvidence(request.Facts)}, nil
}

type supervisedRouteRecorder struct {
	calls    int
	requests []StepRequest
}

func (recorder *supervisedRouteRecorder) OpenSupervisedBlackBox(_ context.Context, request StepRequest) (StepResult, error) {
	recorder.calls++
	recorder.requests = append(recorder.requests, request)
	return StepResult{}, nil
}

func supervisedHandoffPolicy() *HandoffPolicy {
	return &HandoffPolicy{
		AllowedTargets:   []string{"target-one"},
		InputPolicy:      "structured",
		HandoffTimeoutMs: 300000,
		ReturnActions:    []string{"complete", "abort"},
		HooksAfterReturn: []string{"hook-one"},
	}
}

func supervisedDefinition() Definition {
	return Definition{ID: "chain-one", Tasks: []Task{{
		ID:    "task-one",
		Steps: []Step{{ID: "supervised-one", Kind: "code", Mode: SupervisedBlackBox, HandoffPolicy: supervisedHandoffPolicy()}},
	}}}
}

// newSupervisedFixture builds an engine over a single supervised task with the
// chain created but no state driven yet.
func newSupervisedFixture(t *testing.T) *supervisedFixture {
	t.Helper()
	store := &memoryEvents{}
	options, _, _, _, acceptance, reviewer, publisher, _, _ := testOptions(store)
	options.Definition = supervisedDefinition()
	postflight := &fakeSupervisedPostflight{}
	options.SupervisedPostflight = postflight
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	return &supervisedFixture{
		engine:     engine,
		store:      store,
		acceptance: acceptance,
		reviewer:   reviewer,
		publisher:  publisher,
		postflight: postflight,
		session:    supervisedSession(),
		confirm:    supervisedConfirmation(t, SupervisedRiskNoticeVersion, SupervisedRiskNoticeHash()),
	}
}

func supervisedSession() HandoffSession {
	return HandoffSession{
		RunID:                "run-one",
		TaskID:               "task-one",
		StepID:               "supervised-one",
		Attempt:              1,
		HandoffPolicyHash:    baselineHash,
		AllowedReturnActions: []string{"complete", "abort"},
		Operator:             evidence.Actor{Type: "operator", ID: "operator-one"},
		OpenedAt:             "2026-09-07T01:00:00.000Z",
		BeforeManifestHash:   supervisedWorkspaceHash,
		LeaseEvidence:        []string{stopHash},
	}
}

func supervisedStrPtr(value string) *string { return &value }

func supervisedConfirmation(t *testing.T, noticeVersion int, noticeHash string) SupervisedRiskConfirmationRecord {
	t.Helper()
	record, err := NewSupervisedRiskConfirmation(SupervisedRiskConfirmationInput{
		RecordID:          "risk-confirmation-one",
		OperatorID:        "operator-one",
		RunID:             "run-one",
		TaskID:            "task-one",
		StepID:            "supervised-one",
		Attempt:           1,
		WorkspaceHash:     supervisedWorkspaceHash,
		RiskNoticeVersion: noticeVersion,
		RiskNoticeHash:    noticeHash,
		Delivery: SupervisedDelivery{
			Channel:            "sessionbridge-visible",
			DeliveredAt:        "2026-09-07T00:59:00.000Z",
			DeliveryStatus:     "accepted",
			TransportRequestID: supervisedStrPtr("sessbridge-0001"),
			ErrorEvidence:      []string{},
		},
		InteractionEvidence: []string{supervisedInteractionOne, supervisedInteractionTwo},
		Evidence:            []string{supervisedConfirmEvidence},
	}, fixedClock, func() string { return "risk-confirmation-one" })
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func supervisedFacts(receiptHash string) SupervisedFrozenFacts {
	return SupervisedFrozenFacts{
		HandoffReceiptHash:      receiptHash,
		ScanEvidenceHash:        supervisedScanHash,
		ProcessStopEvidenceHash: supervisedStopHash,
		GateEvidenceHash:        supervisedGateHash,
		ReportHash:              supervisedReportHash,
	}
}

func supervisedReturnReceipt(t *testing.T, outcome string) HandoffReceiptRecord {
	t.Helper()
	decision := HandoffDecision{
		ReceiptID:         "handoff-return-one",
		ClosedAt:          "2026-09-07T01:05:00.000Z",
		AfterManifestHash: supervisedAfterManifest,
		DiffHash:          supervisedDiffHash,
		Outcome:           outcome,
		Evidence:          []string{supervisedReceiptEvidence},
	}
	if outcome == "complete" {
		decision.HookResultEvidence = []string{supervisedHookHash}
	}
	record, err := BuildHandoffReceiptRecord(supervisedSession(), decision, fixedClock, func() string { return "handoff-return-one" })
	if err != nil {
		t.Fatal(err)
	}
	return record
}

// driveToRunning parks the supervised step in RUNNING with a paused chain,
// exactly the state a prepared supervised step waits in.
func (fixture *supervisedFixture) driveToRunning(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	chain := chainEntity("run-one")
	task := taskEntity("run-one", "task-one")
	step := stepEntity("run-one", "task-one", "supervised-one")
	for _, transition := range []struct {
		entity evidence.Entity
		to     string
		input  []string
		reason string
	}{
		{chain, "BASELINED", []string{baselineHash}, "baseline-captured"},
		{chain, "RUNNING", nil, "chain-started"},
		{task, "PENDING", nil, "task-created"},
		{task, "PRECHECK", nil, "task-precheck"},
		{task, "STEPS_RUNNING", nil, "steps-started"},
		{step, "PENDING", nil, "step-created"},
		{step, "RUNNING", nil, "step-started"},
		{chain, "PAUSED", nil, "operator-interaction-required"},
	} {
		if err := fixture.engine.transition(ctx, transition.entity, transition.to, transition.input, transition.reason); err != nil {
			t.Fatal(err)
		}
	}
}

func (fixture *supervisedFixture) openWindow(t *testing.T) {
	t.Helper()
	if err := fixture.engine.OpenSupervisedWorkWindow(context.Background(), fixture.session, fixture.confirm); err != nil {
		t.Fatalf("opening the supervised window failed: %v", err)
	}
}

func (fixture *supervisedFixture) enterOpenWindow(t *testing.T) {
	t.Helper()
	fixture.driveToRunning(t)
	fixture.openWindow(t)
}

func (fixture *supervisedFixture) eventCount(t *testing.T) int {
	t.Helper()
	events, err := fixture.store.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return len(events)
}

func (fixture *supervisedFixture) requireStates(t *testing.T, wantStep, wantTask, wantChain string) {
	t.Helper()
	projection := fixture.engine.Projection()
	if got := projection.StepStates[stepKey("task-one", "supervised-one", 1)]; got != wantStep {
		t.Fatalf("step state %s, want %s (%+v)", got, wantStep, projection)
	}
	if got := projection.TaskStates[taskKey("task-one", 1)]; got != wantTask {
		t.Fatalf("task state %s, want %s (%+v)", got, wantTask, projection)
	}
	if got := projection.ChainState; got != wantChain {
		t.Fatalf("chain state %s, want %s (%+v)", got, wantChain, projection)
	}
}

// appendSupervisedTaskReview records a task-level REVIEW_PENDING transition.
func (fixture *supervisedFixture) appendSupervisedTaskReview(t *testing.T, input []string) {
	t.Helper()
	records, err := fixture.store.Load(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	previous := records[len(records)-1].EventHash
	predecessor := previous
	if input == nil {
		input = []string{}
	}
	record, err := evidence.NewStateEvent(evidence.Event{
		EventID:           "event-supervised-task-review",
		RunID:             "run-one",
		Sequence:          len(records) + 1,
		OccurredAt:        fixedClock().Format("2006-01-02T15:04:05.000Z"),
		Actor:             evidence.Actor{Type: "system", ID: "engine"},
		Entity:            evidence.Entity{Kind: "task", RunID: "run-one", TaskID: "task-one", Attempt: 1},
		FromState:         "STEPS_RUNNING",
		ToState:           "REVIEW_PENDING",
		PreviousEventHash: &predecessor,
		InputEvidence:     input,
		Reason:            evidence.Reason{Code: "review-pending"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := fixture.store.Append(context.Background(), record); err != nil {
		t.Fatal(err)
	}
}

func TestSupervisedRiskConfirmationRoundTrip(t *testing.T) {
	record := supervisedConfirmation(t, SupervisedRiskNoticeVersion, SupervisedRiskNoticeHash())
	if record.Record.Selection != "confirm" || !evidence.ValidHash(record.RecordHash) {
		t.Fatalf("unexpected confirmation: %+v", record)
	}
	payload, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeSupervisedRiskConfirmation(payload)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.RecordHash != record.RecordHash || !reflect.DeepEqual(decoded.Record, record.Record) {
		t.Fatalf("roundtrip mismatch: %+v", decoded)
	}
	// A tampered body whose stored hash no longer matches must be refused.
	tampered := record
	tampered.Record.TaskID = "task-two"
	if err := ValidateSupervisedRiskConfirmation(tampered); !errors.Is(err, ErrInvalidSupervisedConfirmation) {
		t.Fatalf("a tampered body must be refused, got %v", err)
	}
}

func TestSupervisedRiskConfirmationLedgerIdempotentAppend(t *testing.T) {
	ledger := &SupervisedRiskConfirmationLedger{}
	first := supervisedConfirmation(t, SupervisedRiskNoticeVersion, SupervisedRiskNoticeHash())
	if err := ledger.Append(first); err != nil {
		t.Fatal(err)
	}
	if err := ledger.Append(first); err != nil {
		t.Fatalf("an identical append must converge, got %v", err)
	}
	if len(ledger.Records) != 1 {
		t.Fatalf("an identical append must not duplicate: %d records", len(ledger.Records))
	}
	if !ledger.Confirmed(SupervisedRiskConfirmationBinding{RunID: "run-one", TaskID: "task-one", StepID: "supervised-one", Attempt: 1}) {
		t.Fatal("the recorded binding must be confirmed")
	}
	if ledger.Confirmed(SupervisedRiskConfirmationBinding{RunID: "run-one", TaskID: "task-one", StepID: "supervised-one", Attempt: 2}) {
		t.Fatal("another attempt must not be confirmed")
	}
	divergent := first
	divergent.Record.RecordID = "risk-confirmation-two"
	hash, err := digestRecord(supervisedRiskConfirmationDomain, divergent.Record)
	if err != nil {
		t.Fatal(err)
	}
	divergent.RecordHash = hash
	if err := ledger.Append(divergent); !errors.Is(err, ErrInvalidSupervisedConfirmation) {
		t.Fatalf("a divergent confirmation for the same step attempt must conflict, got %v", err)
	}
	if len(ledger.Records) != 1 {
		t.Fatal("a divergent append replaced the durable confirmation")
	}
}

func TestSupervisedReportNewAlwaysReduced(t *testing.T) {
	record, err := NewSupervisedReport(SupervisedReportInput{
		ReportID:           "supervised-report-one",
		RunID:              "run-one",
		TaskID:             "task-one",
		StepID:             "supervised-one",
		Attempt:            1,
		VerifiedDimensions: []string{"independent-gates", "process-stop-proof"},
		UnknownDimensions:  []string{"tool-use"},
		Evidence:           []string{supervisedReportEvidence},
	}, fixedClock, func() string { return "supervised-report-one" })
	if err != nil {
		t.Fatal(err)
	}
	if record.Record.AssuranceLevel != "reduced" {
		t.Fatalf("the assurance level is permanently reduced, got %q", record.Record.AssuranceLevel)
	}
	if record.Record.Mode != "supervised-black-box" {
		t.Fatalf("unexpected mode %q", record.Record.Mode)
	}
	if !containsAll(record.Record.UnknownDimensions, "cost") {
		t.Fatalf("the unknown dimension set must always carry cost: %v", record.Record.UnknownDimensions)
	}
	// An empty unknown set falls back to the full closed set, still carrying cost.
	empty, err := NewSupervisedReport(SupervisedReportInput{ReportID: "supervised-report-two", RunID: "run-one", TaskID: "task-one", StepID: "supervised-one", Attempt: 1, Evidence: []string{supervisedReportEvidence}}, fixedClock, func() string { return "supervised-report-two" })
	if err != nil {
		t.Fatal(err)
	}
	if len(empty.Record.UnknownDimensions) != len(supervisedUnknownDimensionSet) {
		t.Fatalf("an empty unknown set must fall back to the closed set: %v", empty.Record.UnknownDimensions)
	}
	// A full assurance level, an empty unknown set and a cost-free unknown set
	// are all refused even when the record hash is recomputed over them.
	for name, mutate := range map[string]func(*SupervisedReport){
		"full assurance": func(body *SupervisedReport) { body.AssuranceLevel = "full" },
		"empty unknown":  func(body *SupervisedReport) { body.UnknownDimensions = []string{} },
		"cost-free":      func(body *SupervisedReport) { body.UnknownDimensions = []string{"tool-use"} },
		"unknown id":     func(body *SupervisedReport) { body.UnknownDimensions = []string{"cost", "not-a-dimension"} },
	} {
		t.Run(name, func(t *testing.T) {
			body := record.Record
			mutate(&body)
			hash, err := digestRecord(supervisedReportDomain, body)
			if err != nil {
				t.Fatal(err)
			}
			broken := SupervisedReportRecord{SchemaVersion: supervisedReportSchema, Record: body, RecordHash: hash}
			if err := ValidateSupervisedReport(broken); !errors.Is(err, ErrInvalidSupervisedReport) {
				t.Fatalf("expected a report rejection, got %v", err)
			}
		})
	}
}

func TestOpenSupervisedWorkWindowWithValidConfirmation(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.driveToRunning(t)
	if err := fixture.engine.OpenSupervisedWorkWindow(context.Background(), fixture.session, fixture.confirm); err != nil {
		t.Fatal(err)
	}
	fixture.requireStates(t, "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR", "PAUSED")
	found := false
	for _, record := range mustEvents(t, fixture.store) {
		if record.Event.Entity.Kind != "step" || record.Event.ToState != "WAITING_FOR_OPERATOR" {
			continue
		}
		if record.Event.Reason.Code != supervisedWorkWindowOpenReason || !containsAll(record.Event.InputEvidence, fixture.confirm.RecordHash) {
			t.Fatalf("the open transition must bind the confirmation hash: %+v", record.Event)
		}
		found = true
	}
	if !found {
		t.Fatal("no window-opening transition was recorded")
	}
}

func TestOpenSupervisedWorkWindowWithoutConfirmation(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.driveToRunning(t)
	events := fixture.eventCount(t)
	if err := fixture.engine.OpenSupervisedWorkWindow(context.Background(), fixture.session, SupervisedRiskConfirmationRecord{}); !errors.Is(err, ErrInvalidSupervisedConfirmation) {
		t.Fatalf("a missing confirmation must be refused, got %v", err)
	}
	if fixture.eventCount(t) != events {
		t.Fatal("the refused open wrote state")
	}
	fixture.requireStates(t, "RUNNING", "STEPS_RUNNING", "PAUSED")
}

func TestOpenSupervisedWorkWindowWrongNoticeVersion(t *testing.T) {
	tests := map[string]SupervisedRiskConfirmationRecord{
		"version drift": supervisedConfirmation(t, SupervisedRiskNoticeVersion+1, SupervisedRiskNoticeHash()),
		"hash drift":    supervisedConfirmation(t, SupervisedRiskNoticeVersion, supervisedSpareHash),
	}
	for name, confirmation := range tests {
		t.Run(name, func(t *testing.T) {
			fixture := newSupervisedFixture(t)
			fixture.driveToRunning(t)
			events := fixture.eventCount(t)
			if err := fixture.engine.OpenSupervisedWorkWindow(context.Background(), fixture.session, confirmation); !errors.Is(err, ErrInvalidSupervisedConfirmation) {
				t.Fatalf("a notice drift must be refused, got %v", err)
			}
			if fixture.eventCount(t) != events {
				t.Fatal("the refused open wrote state")
			}
			fixture.requireStates(t, "RUNNING", "STEPS_RUNNING", "PAUSED")
		})
	}
}

func TestSupervisedRiskConfirmationRejectsBindingDrift(t *testing.T) {
	mutations := map[string]func(*SupervisedRiskConfirmationInput){
		"run":       func(input *SupervisedRiskConfirmationInput) { input.RunID = "run-two" },
		"task":      func(input *SupervisedRiskConfirmationInput) { input.TaskID = "task-two" },
		"step":      func(input *SupervisedRiskConfirmationInput) { input.StepID = "step-two" },
		"attempt":   func(input *SupervisedRiskConfirmationInput) { input.Attempt = 2 },
		"workspace": func(input *SupervisedRiskConfirmationInput) { input.WorkspaceHash = acceptedTwo },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			fixture := newSupervisedFixture(t)
			base := SupervisedRiskConfirmationInput{
				RecordID:            "risk-confirmation-one",
				OperatorID:          "operator-one",
				RunID:               "run-one",
				TaskID:              "task-one",
				StepID:              "supervised-one",
				Attempt:             1,
				WorkspaceHash:       supervisedWorkspaceHash,
				RiskNoticeVersion:   SupervisedRiskNoticeVersion,
				RiskNoticeHash:      SupervisedRiskNoticeHash(),
				Delivery:            fixture.confirm.Record.Delivery,
				InteractionEvidence: []string{supervisedInteractionOne, supervisedInteractionTwo},
				Evidence:            []string{supervisedConfirmEvidence},
			}
			mutate(&base)
			drifted, err := NewSupervisedRiskConfirmation(base, fixedClock, func() string { return "risk-confirmation-one" })
			if err != nil {
				t.Fatalf("the drifted confirmation must still be a well-formed record: %v", err)
			}
			fixture.driveToRunning(t)
			events := fixture.eventCount(t)
			if err := fixture.engine.OpenSupervisedWorkWindow(context.Background(), fixture.session, drifted); !errors.Is(err, ErrInvalidSupervisedConfirmation) {
				t.Fatalf("a binding drift must be refused, got %v", err)
			}
			if fixture.eventCount(t) != events {
				t.Fatal("the refused open wrote state")
			}
			fixture.requireStates(t, "RUNNING", "STEPS_RUNNING", "PAUSED")
		})
	}
}

func TestSupervisedRiskConfirmationRejectsRejectSelection(t *testing.T) {
	record := supervisedConfirmation(t, SupervisedRiskNoticeVersion, SupervisedRiskNoticeHash())
	record.Record.Selection = "reject"
	hash, err := digestRecord(supervisedRiskConfirmationDomain, record.Record)
	if err != nil {
		t.Fatal(err)
	}
	record.RecordHash = hash
	if err := ValidateSupervisedRiskConfirmation(record); !errors.Is(err, ErrInvalidSupervisedConfirmation) {
		t.Fatalf("a reject selection must never validate as a confirmation, got %v", err)
	}
}

func TestOpenSupervisedWorkWindowWithoutActiveStep(t *testing.T) {
	fixture := newSupervisedFixture(t)
	// The chain is still CREATED: a confirmation cannot open a window.
	if err := fixture.engine.OpenSupervisedWorkWindow(context.Background(), fixture.session, fixture.confirm); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("an inactive step must be refused, got %v", err)
	}
	if fixture.eventCount(t) != 1 {
		t.Fatalf("the refusal wrote state: %d events", fixture.eventCount(t))
	}
}

func TestSubmitSupervisedReturnCompleteFlow(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.enterOpenWindow(t)
	receipt := supervisedReturnReceipt(t, "complete")
	facts := supervisedFacts(receipt.ReceiptHash)
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, facts); err != nil {
		t.Fatal(err)
	}
	fixture.requireStates(t, "PASSED", "STEPS_RUNNING", "PAUSED")
	for _, record := range mustEvents(t, fixture.store) {
		if record.Event.Reason.Code != supervisedReturnCompletedReason {
			continue
		}
		rebuilt, err := supervisedFactsFromEvidence(record.Event.InputEvidence)
		if err != nil {
			t.Fatalf("the step PASSED evidence must rebuild the frozen facts: %v", err)
		}
		if rebuilt != facts {
			t.Fatalf("rebuilt facts mismatch: %+v", rebuilt)
		}
		return
	}
	t.Fatal("no supervised return completion was recorded")
}

func TestSubmitSupervisedReturnIdempotentReplay(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.enterOpenWindow(t)
	receipt := supervisedReturnReceipt(t, "complete")
	facts := supervisedFacts(receipt.ReceiptHash)
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, facts); err != nil {
		t.Fatal(err)
	}
	events := fixture.eventCount(t)
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, facts); err != nil {
		t.Fatalf("an identical replay must converge, got %v", err)
	}
	if fixture.eventCount(t) != events {
		t.Fatal("the identical replay wrote state")
	}
	divergent := facts
	divergent.ScanEvidenceHash = supervisedSpareHash
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, divergent); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("a replay with a different fact set must conflict, got %v", err)
	}
	if fixture.eventCount(t) != events {
		t.Fatal("the conflicting replay wrote state")
	}
}

func TestSubmitSupervisedReturnRejectsAbortReceipt(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.enterOpenWindow(t)
	receipt := supervisedReturnReceipt(t, "abort")
	events := fixture.eventCount(t)
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, supervisedFacts(receipt.ReceiptHash)); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("an abort receipt must never complete a supervised step, got %v", err)
	}
	if fixture.eventCount(t) != events {
		t.Fatal("the refused return wrote state")
	}
	fixture.requireStates(t, "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR", "PAUSED")
}

func TestSubmitSupervisedReturnRejectsFactCollisions(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.enterOpenWindow(t)
	receipt := supervisedReturnReceipt(t, "complete")
	facts := supervisedFacts(receipt.ReceiptHash)
	facts.ScanEvidenceHash = facts.HandoffReceiptHash
	events := fixture.eventCount(t)
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, facts); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("colliding facts must be refused, got %v", err)
	}
	if fixture.eventCount(t) != events {
		t.Fatal("the refused return wrote state")
	}
}

func TestSubmitSupervisedReturnRejectsDeliveryOnlyEvidence(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.enterOpenWindow(t)
	receipt := supervisedReturnReceipt(t, "complete")
	// The confirmation that unlocked the window is a delivery-era fact: it can
	// prove supervision was acknowledged, never that work completed.
	facts := supervisedFacts(receipt.ReceiptHash)
	facts.ScanEvidenceHash = fixture.confirm.RecordHash
	events := fixture.eventCount(t)
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, facts); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("delivery/confirmation digests must never become frozen facts, got %v", err)
	}
	if fixture.eventCount(t) != events {
		t.Fatal("the refused return wrote state")
	}
	fixture.requireStates(t, "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR", "PAUSED")
}

func TestSubmitSupervisedReturnRejectsReceiptNotBoundToFacts(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.enterOpenWindow(t)
	receipt := supervisedReturnReceipt(t, "complete")
	facts := supervisedFacts(supervisedSpareHash)
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, facts); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("the handoff fact must be the return receipt digest, got %v", err)
	}
}

func TestSupervisedPostflightDecisionMustConsumeAllFacts(t *testing.T) {
	facts := supervisedFacts(supervisedSpareHash)
	passed := SupervisedPostflightDecision{Outcome: PostflightPassed, Evidence: supervisedFactsEvidence(facts)}
	if err := passed.Validate(facts); err != nil {
		t.Fatalf("a passed decision carrying every fact must validate: %v", err)
	}
	incomplete := SupervisedPostflightDecision{Outcome: PostflightPassed, Evidence: supervisedFactsEvidence(facts)[:4]}
	if err := incomplete.Validate(facts); !errors.Is(err, ErrInvalidPostflightDecision) {
		t.Fatalf("a passed decision without every fact must be refused, got %v", err)
	}
	unknown := SupervisedPostflightDecision{Outcome: "approved", Evidence: supervisedFactsEvidence(facts)}
	if err := unknown.Validate(facts); !errors.Is(err, ErrInvalidPostflightDecision) {
		t.Fatalf("an unknown outcome must be refused, got %v", err)
	}
	for _, outcome := range []PostflightOutcome{PostflightFailed, PostflightUncertain} {
		decision := SupervisedPostflightDecision{Outcome: outcome, ErrorEvidence: []string{supervisedSpareHash}}
		if err := decision.Validate(facts); err != nil {
			t.Fatalf("outcome %s must validate without the fact digests: %v", outcome, err)
		}
	}
}

func TestSupervisedPostflightQualifiesReviewPending(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.enterOpenWindow(t)
	receipt := supervisedReturnReceipt(t, "complete")
	facts := supervisedFacts(receipt.ReceiptHash)
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, facts); err != nil {
		t.Fatal(err)
	}
	if err := fixture.engine.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if fixture.postflight.calls != 1 || len(fixture.postflight.requests) != 1 {
		t.Fatalf("the supervised postflight must run exactly once: %d", fixture.postflight.calls)
	}
	if fixture.postflight.requests[0].Facts != facts || fixture.postflight.requests[0].TaskID != "task-one" || fixture.postflight.requests[0].ParentSnapshotHash != baselineHash {
		t.Fatalf("unexpected supervised postflight request: %+v", fixture.postflight.requests[0])
	}
	fixture.requireStates(t, "PASSED", "PASSED", "COMPLETED")
	for _, record := range mustEvents(t, fixture.store) {
		if record.Event.Entity.Kind != "task" || record.Event.ToState != "REVIEW_PENDING" {
			continue
		}
		if record.Event.Reason.Code != "supervised-postflight-passed" || !containsAll(record.Event.InputEvidence, supervisedFactsEvidence(facts)...) {
			t.Fatalf("the review transition must carry the supervised facts: %+v", record.Event)
		}
		return
	}
	t.Fatal("the supervised postflight recorded no review transition")
}

func TestSupervisedUnqualifiedReviewRefused(t *testing.T) {
	fixture := newSupervisedFixture(t)
	fixture.enterOpenWindow(t)
	receipt := supervisedReturnReceipt(t, "complete")
	if err := fixture.engine.SubmitSupervisedReturn(context.Background(), receipt, supervisedFacts(receipt.ReceiptHash)); err != nil {
		t.Fatal(err)
	}
	// A review state written without the gate cannot be honoured.
	fixture.appendSupervisedTaskReview(t, []string{stopHash})
	events := fixture.eventCount(t)
	if err := fixture.engine.requirePostflightQualifiedReview(context.Background(), "task-one"); !errors.Is(err, ErrUnqualifiedReview) {
		t.Fatalf("an unproven supervised review state must be refused, got %v", err)
	}
	if fixture.eventCount(t) != events {
		t.Fatal("the refusal wrote state")
	}
}

func TestSupervisedTaskWithoutReturnFactsFailsClosed(t *testing.T) {
	fixture := newSupervisedFixture(t)
	ctx := context.Background()
	chain := chainEntity("run-one")
	task := taskEntity("run-one", "task-one")
	step := stepEntity("run-one", "task-one", "supervised-one")
	for _, transition := range []struct {
		entity evidence.Entity
		to     string
		input  []string
		reason string
	}{
		{chain, "BASELINED", []string{baselineHash}, "baseline-captured"},
		{chain, "RUNNING", nil, "chain-started"},
		{task, "PENDING", nil, "task-created"},
		{task, "PRECHECK", nil, "task-precheck"},
		{task, "STEPS_RUNNING", nil, "steps-started"},
		{step, "PENDING", nil, "step-created"},
		{step, "RUNNING", nil, "step-started"},
		// A plain step pass without a supervised return: it must never qualify
		// the task for review.
		{step, "PASSED", []string{stopHash}, "step-passed"},
	} {
		if err := fixture.engine.transition(ctx, transition.entity, transition.to, transition.input, transition.reason); err != nil {
			t.Fatal(err)
		}
	}
	acceptance := fixture.acceptance.calls
	if err := fixture.engine.Run(ctx); !errors.Is(err, ErrInvalidPostflightFacts) {
		t.Fatalf("a supervised step without return facts must fail closed, got %v", err)
	}
	if fixture.acceptance.calls != acceptance {
		t.Fatal("the fail-closed path reached the acceptance flow")
	}
	fixture.requireStates(t, "PASSED", "FAILED", "FAILED")
	for _, record := range mustEvents(t, fixture.store) {
		if record.Event.Entity.Kind == "task" && record.Event.ToState == "REVIEW_PENDING" {
			t.Fatal("REVIEW_PENDING was written without return facts")
		}
	}
}

func TestSupervisedStepRejectsAgentRunnerExecutionTarget(t *testing.T) {
	step := Step{ID: "supervised-one", Kind: "code", Mode: SupervisedBlackBox, HandoffPolicy: supervisedHandoffPolicy()}
	facts := &AgentRunnerImmutableFacts{RequestID: "request-one"}
	if err := (StepExecutionIntent{ExecutionTarget: AgentRunnerExecution, AgentRunnerFacts: facts}).validate(step); !errors.Is(err, ErrInvalidDefinition) {
		t.Fatalf("a supervised step must never accept an AgentRunner execution target, got %v", err)
	}
	if err := (StepExecutionIntent{ExecutionTarget: DefaultExecution, AgentRunnerFacts: facts}).validate(step); !errors.Is(err, ErrInvalidDefinition) {
		t.Fatalf("a supervised step must never carry AgentRunner facts, got %v", err)
	}
	// Positive control: the same step with the default target is routable.
	if err := (StepExecutionIntent{ExecutionTarget: DefaultExecution}).validate(step); err != nil {
		t.Fatalf("a supervised step with the default target must validate: %v", err)
	}
}

// requireKeys pins the exact key set of a marshalled record to the frozen
// schema, including the additionalProperties:false contract.
func requireKeys(t *testing.T, fields map[string]any, want []string) {
	t.Helper()
	if len(fields) != len(want) {
		t.Fatalf("key count %d, want %d: %v", len(fields), len(want), fields)
	}
	for _, key := range want {
		if _, ok := fields[key]; !ok {
			t.Fatalf("missing key %q in %v", key, fields)
		}
	}
}

func marshalFields(t *testing.T, value any) map[string]any {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil {
		t.Fatal(err)
	}
	return fields
}

func nestedFields(t *testing.T, fields map[string]any, key string) map[string]any {
	t.Helper()
	nested, ok := fields[key].(map[string]any)
	if !ok {
		t.Fatalf("%q is not an object: %v", key, fields)
	}
	return nested
}

func TestSupervisedRecordJSONShapesMatchFrozenSchemas(t *testing.T) {
	confirmation := marshalFields(t, supervisedConfirmation(t, SupervisedRiskNoticeVersion, SupervisedRiskNoticeHash()))
	requireKeys(t, confirmation, []string{"schemaVersion", "record", "recordHash"})
	requireKeys(t, nestedFields(t, confirmation, "record"), []string{
		"kind", "recordId", "issuedAt", "issuedBy", "runId", "taskId", "stepId", "attempt",
		"workspaceHash", "riskNoticeVersion", "riskNoticeHash", "selection", "delivery",
		"interactionEvidence", "evidence",
	})
	requireKeys(t, nestedFields(t, nestedFields(t, confirmation, "record"), "delivery"), []string{
		"channel", "deliveredAt", "deliveryStatus", "transportRequestId", "errorEvidence",
	})
	requireKeys(t, nestedFields(t, nestedFields(t, confirmation, "record"), "issuedBy"), []string{"type", "id"})

	report, err := NewSupervisedReport(SupervisedReportInput{
		ReportID:           "supervised-report-one",
		RunID:              "run-one",
		TaskID:             "task-one",
		StepID:             "supervised-one",
		Attempt:            1,
		VerifiedDimensions: []string{"independent-gates"},
		UnknownDimensions:  []string{"cost"},
		Evidence:           []string{supervisedReportEvidence},
	}, fixedClock, func() string { return "supervised-report-one" })
	if err != nil {
		t.Fatal(err)
	}
	record := marshalFields(t, report)
	requireKeys(t, record, []string{"schemaVersion", "record", "recordHash"})
	requireKeys(t, nestedFields(t, record, "record"), []string{
		"kind", "reportId", "recordedAt", "recordedBy", "runId", "taskId", "stepId", "attempt",
		"mode", "assuranceLevel", "verifiedDimensions", "unknownDimensions", "evidence",
	})
	requireKeys(t, nestedFields(t, nestedFields(t, record, "record"), "recordedBy"), []string{"type", "id"})
}

func TestRouterRefusesSupervisedWithoutPort(t *testing.T) {
	request := StepRequest{RunID: "run-one", TaskID: "task-one", Step: Step{ID: "supervised-one", Kind: "code", Mode: SupervisedBlackBox}, Attempt: 1}
	if _, err := (StepRouter{}).Execute(context.Background(), request); !errors.Is(err, ErrInvalidDefinition) {
		t.Fatalf("a supervised step without its port must be refused, got %v", err)
	}
	recorder := &supervisedRouteRecorder{requests: nil}
	router := StepRouter{Managed: &routeRecorder{}, Isolated: &routeRecorder{}, Handoff: &routeRecorder{}, Supervised: recorder, Hooks: &routeRecorder{}}
	if _, err := router.Execute(context.Background(), request); err != nil {
		t.Fatalf("a supervised step must route to its own port: %v", err)
	}
	if recorder.calls != 1 || len(recorder.requests) != 1 || recorder.requests[0].Step.Mode != SupervisedBlackBox {
		t.Fatalf("the supervised port was not the sole route: %+v", recorder)
	}
}

// supervisedConfirmationWithDelivery re-freezes a valid confirmation whose
// delivery fact carries the given status. The record shape permits the
// unavailable/unknown statuses; the window opener is the place that must refuse
// them.
func supervisedConfirmationWithDelivery(t *testing.T, status string) SupervisedRiskConfirmationRecord {
	t.Helper()
	record := supervisedConfirmation(t, SupervisedRiskNoticeVersion, SupervisedRiskNoticeHash())
	record.Record.Delivery.DeliveryStatus = status
	hash, err := digestRecord(supervisedRiskConfirmationDomain, record.Record)
	if err != nil {
		t.Fatal(err)
	}
	record.RecordHash = hash
	if err := ValidateSupervisedRiskConfirmation(record); err != nil {
		t.Fatalf("a %s delivery fact must still be a well-formed record: %v", status, err)
	}
	return record
}

// TestOpenSupervisedWorkWindowRejectsUnacceptedDelivery pins the delivery-status
// invariant of contracts section 7.1a at the window-opening entry: a
// confirmation whose delivery was not accepted must be refused with zero writes,
// even though the record itself is well-formed.
func TestOpenSupervisedWorkWindowRejectsUnacceptedDelivery(t *testing.T) {
	for _, status := range []string{"unavailable", "unknown"} {
		t.Run(status, func(t *testing.T) {
			fixture := newSupervisedFixture(t)
			fixture.driveToRunning(t)
			events := fixture.eventCount(t)
			err := fixture.engine.OpenSupervisedWorkWindow(context.Background(), fixture.session, supervisedConfirmationWithDelivery(t, status))
			if !errors.Is(err, ErrInvalidSupervisedConfirmation) {
				t.Fatalf("a non-accepted delivery must be refused at the window opener, got %v", err)
			}
			if fixture.eventCount(t) != events {
				t.Fatal("the refused open wrote state")
			}
			fixture.requireStates(t, "RUNNING", "STEPS_RUNNING", "PAUSED")
		})
	}
}

// TestSupervisedRiskNoticeIsTruthfulAndVersioned pins the F-4 contract of the
// versioned risk notice: the text must name the current notice version, and it
// must not claim an OS-level isolation boundary ProofRail does not build.
func TestSupervisedRiskNoticeIsTruthfulAndVersioned(t *testing.T) {
	if !strings.Contains(SupervisedRiskNoticeText, fmt.Sprintf("version %d", SupervisedRiskNoticeVersion)) {
		t.Fatalf("the notice text must name the current version %d: %q", SupervisedRiskNoticeVersion, SupervisedRiskNoticeText)
	}
	if strings.Contains(SupervisedRiskNoticeText, "OS-isolated") {
		t.Fatalf("the notice must not claim an OS-level isolation boundary it does not build: %q", SupervisedRiskNoticeText)
	}
	if !strings.Contains(SupervisedRiskNoticeText, "not an OS-level sandbox") {
		t.Fatalf("the notice must state the absence of OS-level isolation: %q", SupervisedRiskNoticeText)
	}
}
