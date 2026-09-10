package adapters

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func capabilityFinding(status string) AgentRunnerCapabilityFinding {
	return AgentRunnerCapabilityFinding{Status: status, Evidence: []string{queueHashThree}}
}

func agentRunnerCapability() AgentRunnerCapability {
	verified := capabilityFinding("verified")
	return AgentRunnerCapability{
		ProbeID:        "probe-one",
		AdapterID:      "fixture-agent",
		ProbedAt:       "2026-09-10T02:00:00.000Z",
		ExecutableHash: queueHashOne,
		Version:        "fixture-agent 1.0.0",
		ConfigHash:     queueHashTwo,
		Disposition:    "compatible",
		Matrix: AgentRunnerCapabilityMatrix{
			Noninteractive:            verified,
			CWD:                       verified,
			EventStreamOrCompleteLogs: verified,
			SessionCreate:             verified,
			SessionResume:             verified,
			Cancellation:              verified,
			ProcessTreeStop:           verified,
			ToolControl:               verified,
			NetworkControl:            verified,
			PermissionControl:         verified,
			Usage:                     verified,
			UnattendedConfirmations:   verified,
		},
		Evidence: []string{queueHashFour},
	}
}

func TestAgentRunnerCapabilityDispositionFailsClosed(t *testing.T) {
	report := agentRunnerCapability()
	report.Matrix.NetworkControl = capabilityFinding("unknown")
	if _, err := NewAgentRunnerCapabilityRecord(report); !errors.Is(err, ErrInvalidAgentRunnerCapability) {
		t.Fatalf("expected compatible report rejection, got %v", err)
	}
	report.Disposition = "blocked"
	if _, err := NewAgentRunnerCapabilityRecord(report); err != nil {
		t.Fatalf("expected blocked report, got %v", err)
	}
	report.Matrix.NetworkControl.Evidence = []string{queueHashThree, queueHashThree}
	if _, err := NewAgentRunnerCapabilityRecord(report); !errors.Is(err, ErrInvalidAgentRunnerCapability) {
		t.Fatalf("expected duplicate evidence rejection, got %v", err)
	}
}

func TestAgentRunnerCapabilityReplayAndConflict(t *testing.T) {
	first, err := NewAgentRunnerCapabilityRecord(agentRunnerCapability())
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAgentRunnerCapabilityIndex([]AgentRunnerCapabilityRecord{first})
	if err != nil {
		t.Fatal(err)
	}
	if replayed, err := index.Record(first); err != nil || !replayed {
		t.Fatalf("expected replay, replayed=%v err=%v", replayed, err)
	}
	changed := agentRunnerCapability()
	changed.Version = "fixture-agent 1.0.1"
	conflicting, err := NewAgentRunnerCapabilityRecord(changed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := index.Record(conflicting); !errors.Is(err, ErrAgentRunnerCapabilityConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestPinnedAgentRunnerCapabilityReport(t *testing.T) {
	directory := filepath.Join("..", "..", "testdata", "agent-runner", "capability-probes")
	path := filepath.Join(directory, "github-copilot-cli-windows.json")
	wire, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var stored AgentRunnerCapabilityRecord
	if err := json.Unmarshal(wire, &stored); err != nil {
		t.Fatal(err)
	}
	generated, err := NewAgentRunnerCapabilityRecord(stored.Capability)
	if err != nil {
		t.Fatal(err)
	}
	if stored.RecordHash != generated.RecordHash {
		t.Fatalf("pinned recordHash mismatch: got %s want %s", stored.RecordHash, generated.RecordHash)
	}
	if _, err := DecodeAgentRunnerCapabilityRecord(wire); err != nil {
		t.Fatal(err)
	}
	evidenceWire, err := os.ReadFile(filepath.Join(directory, "github-copilot-cli-windows-evidence.json"))
	if err != nil {
		t.Fatal(err)
	}
	evidenceHash := fmt.Sprintf("sha256:%x", sha256.Sum256(evidenceWire))
	if len(stored.Capability.Evidence) != 1 || stored.Capability.Evidence[0] != evidenceHash {
		t.Fatalf("pinned evidence mismatch: got %v want %s", stored.Capability.Evidence, evidenceHash)
	}
	for _, finding := range agentRunnerCapabilityFindings(&stored.Capability.Matrix) {
		if len(finding.Evidence) != 1 || finding.Evidence[0] != evidenceHash {
			t.Fatalf("matrix evidence mismatch: got %v want %s", finding.Evidence, evidenceHash)
		}
	}
}
