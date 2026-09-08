package gates

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func effectTestHash(seed byte) string {
	return "sha256:" + strings.Repeat(string(seed), 64)
}

func validEffectClassification() EffectClassification {
	return EffectClassification{
		RecordID:              "class-hook-go-test",
		Kind:                  "classification",
		CreatedAt:             "2026-09-08T10:00:00.000Z",
		PolicyHash:            effectTestHash('a'),
		SubjectKind:           "hook",
		SubjectID:             "hook-go-test",
		SubjectDefinitionHash: effectTestHash('b'),
		EffectClass:           "read-only",
		ExternalSystems:       []string{},
		RecoveryGuarantee:     "none-required",
		Evidence:              []string{effectTestHash('c')},
	}
}

func validEffectObservation() EffectObservation {
	return EffectObservation{
		RecordID:           "obs-one",
		Kind:               "observation",
		OccurredAt:         "2026-09-08T11:00:00.000Z",
		RunID:              "run-one",
		TaskID:             "task-one",
		Attempt:            1,
		ClassificationHash: effectTestHash('d'),
		AuthorizationHash:  effectTestHash('e'),
		EffectClass:        "read-only",
		OperationKey:       "op-run-go-test",
		Outcome:            "completed",
		Evidence:           []string{effectTestHash('f')},
		ErrorEvidence:      []string{},
	}
}

func TestEffectClassificationRoundTrip(t *testing.T) {
	record, err := NewEffectClassificationRecord(validEffectClassification())
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeEffectClassificationRecord(wire)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.RecordHash != record.RecordHash || decoded.Record.SubjectID != "hook-go-test" {
		t.Fatalf("round trip drift: %+v", decoded)
	}
	second, err := NewEffectClassificationRecord(validEffectClassification())
	if err != nil {
		t.Fatal(err)
	}
	if second.RecordHash != record.RecordHash {
		t.Fatalf("hash not stable: %q vs %q", second.RecordHash, record.RecordHash)
	}
}

func TestEffectClassificationRejectsInvalidFields(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*EffectClassification)
	}{
		{"unknown effect class", func(c *EffectClassification) { c.EffectClass = "read-write-everywhere" }},
		{"unknown subject kind", func(c *EffectClassification) { c.SubjectKind = "model" }},
		{"unknown recovery guarantee", func(c *EffectClassification) { c.RecoveryGuarantee = "magic" }},
		{"empty evidence", func(c *EffectClassification) { c.Evidence = []string{} }},
		{"duplicate external system", func(c *EffectClassification) {
			c.ExternalSystems = []string{"ci-server", "ci-server"}
		}},
		{"bad timestamp", func(c *EffectClassification) { c.CreatedAt = "2026-09-08 10:00:00Z" }},
		{"external-write claims none-required", func(c *EffectClassification) {
			c.EffectClass = "external-write"
			c.RecoveryGuarantee = "none-required"
		}},
		{"external-write claims discard-workspace", func(c *EffectClassification) {
			c.EffectClass = "external-write"
			c.RecoveryGuarantee = "discard-workspace"
		}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			classification := validEffectClassification()
			test.mutate(&classification)
			if _, err := NewEffectClassificationRecord(classification); !errors.Is(err, ErrInvalidEffectRecord) {
				t.Fatalf("expected rejection, got %v", err)
			}
		})
	}
	// external-write with reconcile-only is the only honest combination.
	external := validEffectClassification()
	external.EffectClass = "external-write"
	external.RecoveryGuarantee = "reconcile-only"
	if _, err := NewEffectClassificationRecord(external); err != nil {
		t.Fatalf("external-write + reconcile-only must be accepted: %v", err)
	}
}

func TestAdmitEffectEnforcesS1Rules(t *testing.T) {
	readOnly := validEffectClassification()
	if decision := AdmitEffect(readOnly); decision != EffectAdmit {
		t.Fatalf("read-only must admit, got %q", decision)
	}
	local := validEffectClassification()
	local.EffectClass = "local-discardable"
	if decision := AdmitEffect(local); decision != EffectAdmit {
		t.Fatalf("local-discardable must admit, got %q", decision)
	}
	external := validEffectClassification()
	external.EffectClass = "external-write"
	if decision := AdmitEffect(external); decision != EffectDeny {
		t.Fatalf("external-write must deny in S1, got %q", decision)
	}
	unknown := validEffectClassification()
	unknown.EffectClass = "sneaky-dry-run"
	if decision := AdmitEffect(unknown); decision != EffectDeny {
		t.Fatalf("unknown class must fail closed, got %q", decision)
	}
}

func TestEffectObservationValidationAndRecovery(t *testing.T) {
	completed := validEffectObservation()
	if _, err := NewEffectObservationRecord(completed); err != nil {
		t.Fatal(err)
	}
	if action := ObservationRecoveryAction(completed); action != "none-required" {
		t.Fatalf("completed action %q", action)
	}

	completedNoEvidence := validEffectObservation()
	completedNoEvidence.Evidence = []string{}
	if _, err := NewEffectObservationRecord(completedNoEvidence); !errors.Is(err, ErrInvalidEffectRecord) {
		t.Fatalf("completed without evidence must be rejected, got %v", err)
	}

	uncertain := validEffectObservation()
	uncertain.Outcome = "uncertain"
	uncertain.ErrorEvidence = []string{effectTestHash('f')}
	uncertain.Evidence = []string{}
	if _, err := NewEffectObservationRecord(uncertain); err != nil {
		t.Fatalf("uncertain observation must be accepted: %v", err)
	}
	if action := ObservationRecoveryAction(uncertain); action != "reconcile-only" {
		t.Fatalf("uncertain must reconcile, got %q", action)
	}

	uncertainNoError := validEffectObservation()
	uncertainNoError.Outcome = "failed"
	uncertainNoError.ErrorEvidence = []string{}
	if _, err := NewEffectObservationRecord(uncertainNoError); !errors.Is(err, ErrInvalidEffectRecord) {
		t.Fatalf("failed without error evidence must be rejected, got %v", err)
	}

	duplicate := validEffectObservation()
	duplicate.Evidence = []string{effectTestHash('f'), effectTestHash('f')}
	if _, err := NewEffectObservationRecord(duplicate); !errors.Is(err, ErrInvalidEffectRecord) {
		t.Fatalf("duplicate evidence must be rejected, got %v", err)
	}
}
