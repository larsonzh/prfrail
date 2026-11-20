package adapters

// P83 FIX ROUND 8 counterexample (item T8b): the final start-failed evidence is
// built from the CLOSED classification only. The executor is an interface, so the
// Evidence it returns and the text of its FailureCode are caller-supplied data; a
// refusal must neither carry them nor let them select the record's content, and
// the executor's Measurements must not be merged into a refusal.
//
// Mutation contracts are stated on the test.

import (
	"context"
	"slices"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// TestP83Fix8T8bStartFailedEvidenceIsClosedClassification is the T8b counterexample.
//
// Mutation check: merging execution.Evidence (or execution.Measurements) into the
// refusal evidence, or digesting the raw FailureCode instead of its class, turns
// the corresponding assertion red.
func TestP83Fix8T8bStartFailedEvidenceIsClosedClassification(t *testing.T) {
	executorHash := evidence.Digest("proofrail:test-executor-evidence:1\n", []byte("executor supplied value"))
	measurementHash := evidence.Digest("proofrail:test-executor-measurement:1\n", []byte("measured value"))
	execution := CopilotCLIProbeExecution{
		Outcome:      "start-failed",
		FailureCode:  `logging flag --log-file C:\secret\probe.log refused`,
		Evidence:     []string{executorHash},
		Measurements: []string{measurementHash},
	}
	result := p83Fix8Probe(t, execution)
	if result.Status != "unavailable" || result.Reason != "cli_unavailable" {
		t.Fatalf("a start-failed probe must be an unavailable refusal, got %+v", result)
	}
	if want := copilotCLIProbeRefusalEvidence(execution); !slices.Equal(result.Evidence, want) {
		t.Fatalf("the refusal evidence must be the closed classification: got %#v want %#v", result.Evidence, want)
	}
	if slices.Contains(result.Evidence, executorHash) {
		t.Fatalf("the executor's Evidence must not reach the refusal record: %#v", result.Evidence)
	}
	if slices.Contains(result.Evidence, measurementHash) {
		t.Fatalf("the executor's Measurements must not be merged into a refusal: %#v", result.Evidence)
	}

	// Value independence: another refusal of the same class, with another free-text
	// code and another executor-supplied digest, produces the same record.
	otherExecution := CopilotCLIProbeExecution{
		Outcome:      "start-failed",
		FailureCode:  "a completely different free text naming /tmp/another-value",
		Evidence:     []string{evidence.Digest("proofrail:test-executor-evidence:1\n", []byte("another value"))},
		Measurements: []string{evidence.Digest("proofrail:test-executor-measurement:1\n", []byte("another measurement"))},
	}
	other := p83Fix8Probe(t, otherExecution)
	if !slices.Equal(result.Evidence, other.Evidence) {
		t.Fatalf("the refusal record must not depend on the refused values: %#v vs %#v", result.Evidence, other.Evidence)
	}

	// The classification is read, not constant: a classified failure code names its
	// class and therefore produces a different record.
	classified := execution
	classified.FailureCode = "quota_exhausted"
	if got := p83Fix8Probe(t, classified); slices.Equal(got.Evidence, result.Evidence) {
		t.Fatalf("a classified failure code must be named in the record, got %#v", got.Evidence)
	}
}

// p83Fix8Probe runs one start-failed probe through the Copilot CLI seam.
func p83Fix8Probe(t *testing.T, execution CopilotCLIProbeExecution) AIProbeResult {
	t.Helper()
	probe := CopilotCLIProbe{
		Executable:       `C:\bin\copilot.exe`,
		WorkingDirectory: `C:\probe`,
		Executor:         &fakeCopilotCLIProbeExecutor{execution: execution},
	}
	result, err := probe.Probe(context.Background(), deepSeekProfile(), []byte("not-a-real-secret"), "agent-runner-cli", 1)
	if err != nil {
		t.Fatalf("the probe must return a result, got %v", err)
	}
	return result
}
