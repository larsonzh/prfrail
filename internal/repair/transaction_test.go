package repair

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	hashOne   = "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	hashTwo   = "sha256:2222222222222222222222222222222222222222222222222222222222222222"
	hashThree = "sha256:3333333333333333333333333333333333333333333333333333333333333333"
	hashFour  = "sha256:4444444444444444444444444444444444444444444444444444444444444444"
	hashFive  = "sha256:5555555555555555555555555555555555555555555555555555555555555555"
	hashSix   = "sha256:6666666666666666666666666666666666666666666666666666666666666666"
	hashSeven = "sha256:7777777777777777777777777777777777777777777777777777777777777777"
	hashEight = "sha256:8888888888888888888888888888888888888888888888888888888888888888"
	hashNine  = "sha256:9999999999999999999999999999999999999999999999999999999999999999"
)

func fixedNow() time.Time { return time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC) }

func newTransaction(t *testing.T) *Transaction {
	t.Helper()
	tx, err := NewTransaction("repair-one", "run-one", "task-one", 1, 2, "ledger-one", hashOne, hashTwo, hashThree, fixedNow())
	if err != nil {
		t.Fatal(err)
	}
	return tx
}

func appendHappyPathToValidate(t *testing.T, tx *Transaction) {
	t.Helper()
	actor := evidence.Actor{Type: "system", ID: "proofrail"}
	if err := tx.AppendPrepare(PrepareInput{
		StageID:               "prepare-one",
		StartedAt:             fixedNow(),
		CompletedAt:           fixedNow().Add(time.Minute),
		Actor:                 actor,
		Outcome:               "completed",
		Evidence:              []string{hashFour},
		ParentSnapshotHash:    hashFive,
		BeforeManifestHash:    hashSix,
		TargetIDs:             []string{"target-one"},
		WriterStopEvidence:    []string{hashSeven},
		LeaseEvidence:         []string{hashEight},
		CandidateManifestHash: hashNine,
	}); err != nil {
		t.Fatal(err)
	}
	if err := tx.AppendInspect(InspectInput{
		StageID:               "inspect-one",
		StartedAt:             fixedNow().Add(2 * time.Minute),
		CompletedAt:           fixedNow().Add(3 * time.Minute),
		Actor:                 actor,
		Outcome:               "passed",
		Evidence:              []string{hashFour},
		CandidateManifestHash: hashNine,
		DiffHash:              hashFive,
		ScopeEvidence:         []string{hashSix},
		OwnershipEvidence:     []string{hashSeven},
	}); err != nil {
		t.Fatal(err)
	}
	if err := tx.AppendValidate(ValidateInput{
		StageID:               "validate-one",
		StartedAt:             fixedNow().Add(4 * time.Minute),
		CompletedAt:           fixedNow().Add(5 * time.Minute),
		Actor:                 actor,
		Outcome:               "passed",
		Evidence:              []string{hashFour},
		CandidateManifestHash: hashNine,
		ValidationPlanHash:    hashEight,
		HookResultEvidence:    []string{hashSix},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestRepairRejectsStaleCandidateAndFakePromoteBinding(t *testing.T) {
	tx := newTransaction(t)
	actor := evidence.Actor{Type: "system", ID: "proofrail"}
	if err := tx.AppendPrepare(PrepareInput{
		StageID:               "prepare-one",
		StartedAt:             fixedNow(),
		CompletedAt:           fixedNow().Add(time.Minute),
		Actor:                 actor,
		Outcome:               "completed",
		Evidence:              []string{hashFour},
		ParentSnapshotHash:    hashFive,
		BeforeManifestHash:    hashSix,
		TargetIDs:             []string{"target-one"},
		WriterStopEvidence:    []string{hashSeven},
		LeaseEvidence:         []string{hashEight},
		CandidateManifestHash: hashNine,
	}); err != nil {
		t.Fatal(err)
	}
	if err := tx.AppendInspect(InspectInput{
		StageID:               "inspect-one",
		StartedAt:             fixedNow().Add(2 * time.Minute),
		CompletedAt:           fixedNow().Add(3 * time.Minute),
		Actor:                 actor,
		Outcome:               "passed",
		Evidence:              []string{hashFour},
		CandidateManifestHash: hashEight,
		DiffHash:              hashFive,
		ScopeEvidence:         []string{hashSix},
		OwnershipEvidence:     []string{hashSeven},
	}); !errors.Is(err, ErrStaleCandidate) {
		t.Fatalf("stale candidate error: %v", err)
	}

	tx = newTransaction(t)
	appendHappyPathToValidate(t, tx)
	if err := tx.AppendPromote(PromoteInput{
		StageID:               "promote-one",
		StartedAt:             fixedNow().Add(6 * time.Minute),
		CompletedAt:           fixedNow().Add(7 * time.Minute),
		Actor:                 actor,
		Outcome:               "completed",
		Evidence:              []string{hashFour},
		CandidateManifestHash: hashNine,
		ValidateStageHash:     hashFive,
		TicketLedgerHash:      hashOne,
		WriterStopEvidence:    []string{hashSeven},
		LeaseEvidence:         []string{hashEight},
		ResultingManifestHash: hashTwo,
	}); !errors.Is(err, ErrInvalidStage) {
		t.Fatalf("fake binding error: %v", err)
	}
}

func TestRepairPromoteInterruptionStopsFurtherStages(t *testing.T) {
	tx := newTransaction(t)
	appendHappyPathToValidate(t, tx)
	actor := evidence.Actor{Type: "system", ID: "proofrail"}
	if err := tx.AppendPromote(PromoteInput{
		StageID:               "promote-one",
		StartedAt:             fixedNow().Add(6 * time.Minute),
		CompletedAt:           fixedNow().Add(7 * time.Minute),
		Actor:                 actor,
		Outcome:               "uncertain",
		Evidence:              []string{hashFour},
		ErrorEvidence:         []string{hashFive},
		CandidateManifestHash: hashNine,
		ValidateStageHash:     tx.body.Stages[len(tx.body.Stages)-1].StageHash,
		TicketLedgerHash:      hashOne,
		WriterStopEvidence:    []string{hashSeven},
		LeaseEvidence:         []string{hashEight},
	}); err != nil {
		t.Fatal(err)
	}
	record, err := tx.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if record.Transaction.Status != "uncertain" {
		t.Fatalf("status: %s", record.Transaction.Status)
	}
	if err := tx.AppendPromote(PromoteInput{}); !errors.Is(err, ErrTerminalTransaction) {
		t.Fatalf("terminal promote error: %v", err)
	}
}

func TestRepairStagesSerializeToSchemaShape(t *testing.T) {
	tx := newTransaction(t)
	actor := evidence.Actor{Type: "system", ID: "proofrail"}
	if err := tx.AppendPrepare(PrepareInput{
		StageID:            "prepare-one",
		StartedAt:          fixedNow(),
		CompletedAt:        fixedNow().Add(time.Minute),
		Actor:              actor,
		Outcome:            "failed",
		Evidence:           []string{hashFour},
		ErrorEvidence:      []string{hashFive},
		ParentSnapshotHash: hashFive,
		BeforeManifestHash: hashSix,
		TargetIDs:          []string{"target-one"},
		WriterStopEvidence: []string{hashSeven},
		LeaseEvidence:      []string{hashEight},
	}); err != nil {
		t.Fatal(err)
	}
	if err := tx.AppendInspect(InspectInput{
		StageID:               "inspect-one",
		StartedAt:             fixedNow(),
		CompletedAt:           fixedNow().Add(time.Minute),
		Actor:                 actor,
		Outcome:               "passed",
		Evidence:              []string{hashFour},
		CandidateManifestHash: hashNine,
		DiffHash:              hashFive,
		ScopeEvidence:         []string{hashSix},
		OwnershipEvidence:     []string{hashSeven},
	}); !errors.Is(err, ErrTerminalTransaction) {
		t.Fatalf("prepare failed must terminate: %v", err)
	}
	record, err := tx.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	if record.Transaction.Status != "failed" || len(record.Transaction.Stages) != 1 {
		t.Fatalf("failed transaction: %+v", record.Transaction)
	}
	payload, err := json.Marshal(record.Transaction.Stages[0])
	if err != nil {
		t.Fatal(err)
	}
	var fields map[string]any
	if err := json.Unmarshal(payload, &fields); err != nil {
		t.Fatal(err)
	}
	if value, exists := fields["candidateManifestHash"]; !exists || value != nil {
		t.Fatalf("failed prepare candidate must be null: %v", fields["candidateManifestHash"])
	}
	if _, exists := fields["resultingManifestHash"]; exists {
		t.Fatal("prepare must not include resultingManifestHash")
	}
	if _, exists := fields["previousStageHash"]; !exists {
		t.Fatal("prepare must include previousStageHash as null")
	}
}

func TestRepairRejectsEmptyEvidenceAndDuplicateTargets(t *testing.T) {
	tx := newTransaction(t)
	actor := evidence.Actor{Type: "system", ID: "proofrail"}
	if err := tx.AppendPrepare(PrepareInput{
		StageID:               "prepare-one",
		StartedAt:             fixedNow(),
		CompletedAt:           fixedNow().Add(time.Minute),
		Actor:                 actor,
		Outcome:               "completed",
		ParentSnapshotHash:    hashFive,
		BeforeManifestHash:    hashSix,
		TargetIDs:             []string{"target-one"},
		WriterStopEvidence:    []string{hashSeven},
		LeaseEvidence:         []string{hashEight},
		CandidateManifestHash: hashNine,
	}); !errors.Is(err, ErrInvalidStage) {
		t.Fatalf("empty evidence must be rejected: %v", err)
	}
	tx = newTransaction(t)
	if err := tx.AppendPrepare(PrepareInput{
		StageID:               "prepare-one",
		StartedAt:             fixedNow(),
		CompletedAt:           fixedNow().Add(time.Minute),
		Actor:                 actor,
		Outcome:               "completed",
		Evidence:              []string{hashFour},
		ParentSnapshotHash:    hashFive,
		BeforeManifestHash:    hashSix,
		TargetIDs:             []string{"target-one", "target-one"},
		WriterStopEvidence:    []string{hashSeven},
		LeaseEvidence:         []string{hashEight},
		CandidateManifestHash: hashNine,
	}); !errors.Is(err, ErrInvalidStage) {
		t.Fatalf("duplicate targets must be rejected: %v", err)
	}
}
