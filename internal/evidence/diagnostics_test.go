package evidence

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func diagnosisTestHash(seed byte) string {
	return "sha256:" + strings.Repeat(string(seed), 64)
}

func validRecoveryDiagnosis() RecoveryDiagnosis {
	snapshot := diagnosisTestHash('a')
	return RecoveryDiagnosis{
		RecordID:                 "diag-one",
		Kind:                     "recovery-diagnosis",
		CreatedAt:                "2026-09-08T12:00:00.000Z",
		RunID:                    "run-one",
		Mode:                     "read-only-redacted",
		LastAcceptedSnapshotHash: &snapshot,
		JournalEvidence:          []string{diagnosisTestHash('b')},
		Processes: []DiagnosisProcess{
			{ProcessRef: "proc-runner", State: "unknown", Evidence: []string{diagnosisTestHash('c')}},
		},
		UnknownEffectEvidence: []string{diagnosisTestHash('d')},
		AllowedActions:        []string{"inspect", "controlled-stop"},
		MutationsPerformed:    false,
		Outcome:               "blocked",
		ErrorEvidence:         []string{diagnosisTestHash('e')},
	}
}

func TestRecoveryDiagnosisRoundTrip(t *testing.T) {
	record, err := NewEffectDiagnosisRecord(validRecoveryDiagnosis())
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeEffectDiagnosisRecord(wire)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.RecordHash != record.RecordHash || decoded.Record.Kind != "recovery-diagnosis" {
		t.Fatalf("round trip drift: %+v", decoded)
	}
	if decoded.Record.LastAcceptedSnapshotHash == nil || *decoded.Record.LastAcceptedSnapshotHash != diagnosisTestHash('a') {
		t.Fatalf("snapshot anchor lost: %+v", decoded.Record)
	}
}

func TestRecoveryDiagnosisClearAndBlockedRules(t *testing.T) {
	clear := validRecoveryDiagnosis()
	clear.Outcome = "clear"
	clear.UnknownEffectEvidence = []string{}
	clear.ErrorEvidence = []string{}
	if _, err := NewEffectDiagnosisRecord(clear); err != nil {
		t.Fatalf("clear diagnosis must be accepted: %v", err)
	}

	clearWithUnknown := validRecoveryDiagnosis()
	clearWithUnknown.Outcome = "clear"
	clearWithUnknown.ErrorEvidence = []string{}
	if _, err := NewEffectDiagnosisRecord(clearWithUnknown); !errors.Is(err, ErrInvalidRecoveryDiagnosis) {
		t.Fatalf("clear with unknown effects must be rejected, got %v", err)
	}

	blockedNoError := validRecoveryDiagnosis()
	blockedNoError.ErrorEvidence = []string{}
	if _, err := NewEffectDiagnosisRecord(blockedNoError); !errors.Is(err, ErrInvalidRecoveryDiagnosis) {
		t.Fatalf("blocked without error evidence must be rejected, got %v", err)
	}
}

func TestRecoveryDiagnosisRejectsMutationsAndBadFields(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*RecoveryDiagnosis)
	}{
		{"mutations performed", func(d *RecoveryDiagnosis) { d.MutationsPerformed = true }},
		{"non-redacted mode", func(d *RecoveryDiagnosis) { d.Mode = "full-access" }},
		{"unknown allowed action", func(d *RecoveryDiagnosis) { d.AllowedActions = []string{"delete-logs"} }},
		{"duplicate allowed action", func(d *RecoveryDiagnosis) { d.AllowedActions = []string{"inspect", "inspect"} }},
		{"empty process evidence", func(d *RecoveryDiagnosis) {
			d.Processes = []DiagnosisProcess{{ProcessRef: "proc-x", State: "running", Evidence: []string{}}}
		}},
		{"unknown process state", func(d *RecoveryDiagnosis) {
			d.Processes = []DiagnosisProcess{{ProcessRef: "proc-x", State: "zombie", Evidence: []string{diagnosisTestHash('c')}}}
		}},
		{"bad timestamp", func(d *RecoveryDiagnosis) { d.CreatedAt = "2026-09-08 12:00:00Z" }},
		{"unknown outcome", func(d *RecoveryDiagnosis) { d.Outcome = "maybe" }},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			diagnosis := validRecoveryDiagnosis()
			test.mutate(&diagnosis)
			if _, err := NewEffectDiagnosisRecord(diagnosis); !errors.Is(err, ErrInvalidRecoveryDiagnosis) {
				t.Fatalf("expected rejection, got %v", err)
			}
		})
	}
}

func TestDecodeRecoveryDiagnosisRejectsUnknownField(t *testing.T) {
	record, err := NewEffectDiagnosisRecord(validRecoveryDiagnosis())
	if err != nil {
		t.Fatal(err)
	}
	wire, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	tampered := strings.Replace(string(wire), `"recordId":"diag-one"`, `"recordId":"diag-one","secretMutation":true`, 1)
	if _, err := DecodeEffectDiagnosisRecord([]byte(tampered)); err == nil {
		t.Fatal("expected unknown field rejection")
	}
}
