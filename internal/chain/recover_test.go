package chain

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestRecoverProducesBlockedDiagnosisPlan(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, _, _, _, _, _, _ := testOptions(store)
	options.Failpoint = func(point string) error {
		if point == "step-running" {
			return errors.New("simulated crash")
		}
		return nil
	}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); err == nil {
		t.Fatal("expected simulated crash")
	}

	recoveryOptions, _, _, _, _, _, _, _, _ := testOptions(store)
	sequence := len(store.events)
	recoveryOptions.IDs = func() string { sequence++; return fmt.Sprintf("event-%d", sequence) }
	recovered, err := New(context.Background(), recoveryOptions)
	if err != nil {
		t.Fatal(err)
	}
	if err := recovered.Recover(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("recovery result: %v", err)
	}
	projection := recovered.Projection()
	if projection.ChainState != "PAUSED" || !projection.RecoveryUncertain {
		t.Fatalf("recovery must pause uncertain: %+v", projection)
	}

	accepted := projection.AcceptedParent
	diagnosis := evidence.RecoveryDiagnosis{
		RecordID:                 "diag-run-one",
		Kind:                     "recovery-diagnosis",
		CreatedAt:                "2026-09-08T12:00:00.000Z",
		RunID:                    projection.RunID,
		Mode:                     "read-only-redacted",
		LastAcceptedSnapshotHash: &accepted,
		JournalEvidence:          []string{stopHash},
		Processes: []evidence.DiagnosisProcess{
			{ProcessRef: "proc-runner", State: "unknown", Evidence: []string{stopHash}},
		},
		UnknownEffectEvidence: []string{stopHash},
		AllowedActions:        []string{"inspect", "controlled-stop"},
		MutationsPerformed:    false,
		Outcome:               "blocked",
		ErrorEvidence:         []string{stopHash},
	}
	record, err := evidence.NewEffectDiagnosisRecord(diagnosis)
	if err != nil {
		t.Fatal(err)
	}
	if record.Record.Outcome != "blocked" || record.Record.MutationsPerformed {
		t.Fatalf("diagnosis plan must be blocked and read-only: %+v", record.Record)
	}
	if record.Record.LastAcceptedSnapshotHash == nil || *record.Record.LastAcceptedSnapshotHash != accepted {
		t.Fatal("diagnosis must anchor the last accepted snapshot")
	}
}

func TestRecoverDoesNotMutateEventLog(t *testing.T) {
	store := &memoryEvents{}
	options, _, _, _, _, _, _, _, _ := testOptions(store)
	options.Failpoint = func(point string) error {
		if point == "step-running" {
			return errors.New("simulated crash")
		}
		return nil
	}
	engine, err := New(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	if err := engine.Run(context.Background()); err == nil {
		t.Fatal("expected simulated crash")
	}
	before := append([]evidence.StateEvent(nil), store.events...)

	recoveryOptions, _, _, _, _, _, _, _, _ := testOptions(store)
	sequence := len(store.events)
	recoveryOptions.IDs = func() string { sequence++; return fmt.Sprintf("event-%d", sequence) }
	recovered, err := New(context.Background(), recoveryOptions)
	if err != nil {
		t.Fatal(err)
	}
	if err := recovered.Recover(context.Background()); !errors.Is(err, ErrRecoveryUncertain) {
		t.Fatalf("recovery result: %v", err)
	}
	if len(store.events) < len(before) {
		t.Fatal("recover must never delete events")
	}
	for index := range before {
		if !reflect.DeepEqual(store.events[index], before[index]) {
			t.Fatalf("recover mutated event %d: %+v -> %+v", index, before[index], store.events[index])
		}
	}
}
