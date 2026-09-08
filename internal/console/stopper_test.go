package console

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

func TestDefaultStopRunReturnsTerminalEvidence(t *testing.T) {
	root := t.TempDir()
	cfg := DefaultChainConfig("chain-one")
	definition, err := cfg.ToDefinition()
	if err != nil {
		t.Fatal(err)
	}
	runID := "run-terminal"
	runDir := filepath.Join(root, "tmp", "prfrail-runs", runID)
	if _, err := ExecuteNoopRun(context.Background(), runID, runDir, definition, root, time.Now); err != nil {
		t.Fatal(err)
	}

	stopRun := newDefaultStopRun(func() (string, error) { return root, nil })
	evidenceHashes, err := stopRun(context.Background(), runID, "")
	if err != nil {
		t.Fatalf("terminal run should not require managed identity: %v", err)
	}
	if len(evidenceHashes) != 1 || !evidence.ValidHash(evidenceHashes[0]) {
		t.Fatalf("invalid terminal evidence: %+v", evidenceHashes)
	}
}

func TestDefaultStopRunFailsClosedWhenManagedIdentityMissing(t *testing.T) {
	root := t.TempDir()
	stopRun := newDefaultStopRun(func() (string, error) { return root, nil })
	evidenceHashes, err := stopRun(context.Background(), "run-missing", "")
	if err == nil {
		t.Fatal("expected missing managed identity failure")
	}
	if len(evidenceHashes) == 0 || !evidence.ValidHash(evidenceHashes[0]) {
		t.Fatalf("expected uncertainty evidence hash, got %+v", evidenceHashes)
	}
}

func TestDefaultStopRunRespectsRunDirOverride(t *testing.T) {
	root := t.TempDir()
	cfg := DefaultChainConfig("chain-two")
	definition, err := cfg.ToDefinition()
	if err != nil {
		t.Fatal(err)
	}
	runID := "run-overridden"
	runDir := filepath.Join(root, "custom-runs", runID)
	if _, err := ExecuteNoopRun(context.Background(), runID, runDir, definition, root, time.Now); err != nil {
		t.Fatal(err)
	}

	stopRun := newDefaultStopRun(func() (string, error) { return root, nil })
	evidenceHashes, err := stopRun(context.Background(), runID, runDir)
	if err != nil {
		t.Fatalf("run-dir override should resolve terminal run: %v", err)
	}
	if len(evidenceHashes) != 1 || !evidence.ValidHash(evidenceHashes[0]) {
		t.Fatalf("invalid terminal evidence for override: %+v", evidenceHashes)
	}
}

func TestDefaultStopRunAppliesBoundedStopTimeout(t *testing.T) {
	root := t.TempDir()
	runID := "run-with-identity"
	runDir := filepath.Join(root, "tmp", "prfrail-runs", runID)
	if err := os.MkdirAll(runDir, 0755); err != nil {
		t.Fatal(err)
	}
	identity := guard.ProcessIdentity{PID: 4321, StartToken: "start-token"}
	identityBytes, err := json.Marshal(identity)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runDir, managedProcessIdentityFileName), identityBytes, 0600); err != nil {
		t.Fatal(err)
	}

	called := 0
	var gotDeadline bool
	var gotGrace time.Duration
	stopRun := newDefaultStopRunWithStopper(
		func() (string, error) { return root, nil },
		func(ctx context.Context, gotIdentity guard.ProcessIdentity, grace time.Duration) (guard.TerminationEvidence, error) {
			called++
			gotGrace = grace
			_, gotDeadline = ctx.Deadline()
			if gotIdentity != identity {
				t.Fatalf("identity mismatch: %+v", gotIdentity)
			}
			return guard.TerminationEvidence{
				Identity:  identity,
				Requested: "2026-09-08T12:00:00.000Z",
				Verified:  "2026-09-08T12:00:00.500Z",
				Outcome:   "stopped",
				Actions:   []string{"already-stopped"},
				Hash:      digest("proofrail:termination-evidence:1\n", runID),
			}, nil
		},
	)

	evidenceHashes, err := stopRun(context.Background(), runID, "")
	if err != nil {
		t.Fatalf("stop run failed: %v", err)
	}
	if called != 1 {
		t.Fatalf("stopper called %d times", called)
	}
	if !gotDeadline {
		t.Fatal("stopper context should include a timeout deadline")
	}
	if gotGrace != defaultStopGrace {
		t.Fatalf("grace mismatch: got %s want %s", gotGrace, defaultStopGrace)
	}
	if len(evidenceHashes) != 1 || !evidence.ValidHash(evidenceHashes[0]) {
		t.Fatalf("invalid stop evidence: %+v", evidenceHashes)
	}
}
