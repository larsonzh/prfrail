package adapters

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

// newAgentRunnerTestProcess starts a short-lived real process so the registry is
// exercised with a genuine guard handle. The caller always terminates it.
func newAgentRunnerTestProcess(t *testing.T) *guard.ManagedProcess {
	t.Helper()
	managed, err := guard.StartManaged(context.Background(), guard.ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=^$"},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if _, err := managed.Terminate(context.Background(), 200*time.Millisecond); err != nil {
			t.Logf("terminate test process: %v", err)
		}
	})
	return managed
}

func TestAgentRunnerProcessIDRoundTripFailsClosed(t *testing.T) {
	managed := newAgentRunnerTestProcess(t)
	encoded, err := agentRunnerProcessID(managed.Identity())
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := decodeAgentRunnerProcessID(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != managed.Identity() {
		t.Fatalf("identity round trip mismatch: %+v vs %+v", decoded, managed.Identity())
	}
	for name, value := range map[string]string{
		"empty":           "",
		"blank":           "   ",
		"not json":        "pid=1",
		"unknown field":   `{"pid":1,"startToken":"t","extra":true}`,
		"missing token":   `{"pid":1,"startToken":""}`,
		"nonpositive pid": `{"pid":0,"startToken":"t"}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := decodeAgentRunnerProcessID(value); !errors.Is(err, ErrAgentRunnerProcessIdentity) {
				t.Fatalf("expected a fail-closed identity rejection, got %v", err)
			}
		})
	}
	if _, err := agentRunnerProcessID(guard.ProcessIdentity{}); !errors.Is(err, ErrAgentRunnerProcessIdentity) {
		t.Fatalf("an incomplete identity must never be encoded, got %v", err)
	}
}

func TestAgentRunnerProcessRegistryRefusesDuplicatesAndRestart(t *testing.T) {
	managed := newAgentRunnerTestProcess(t)
	registry := newAgentRunnerProcessRegistry()
	handle := &agentRunnerProcessHandle{
		launchID:  "launch-request-one",
		requestID: "request-one",
		process:   managed,
		phase:     agentRunnerProcessStarting,
	}
	if err := registry.register(handle); err != nil {
		t.Fatal(err)
	}
	if err := registry.register(handle); !errors.Is(err, ErrAgentRunnerProcessDuplicate) {
		t.Fatalf("a duplicate launch registration must be refused, got %v", err)
	}
	if err := registry.register(&agentRunnerProcessHandle{launchID: "launch-request-two", requestID: "request-two", process: nil}); !errors.Is(err, ErrAgentRunnerProcessIdentity) {
		t.Fatalf("a handle without a process must be refused, got %v", err)
	}
	if err := registry.advance("launch-request-one", agentRunnerProcessRunning); err != nil {
		t.Fatal(err)
	}
	if err := registry.advance("launch-request-one", agentRunnerProcessStopping); err != nil {
		t.Fatal(err)
	}
	if err := registry.advance("launch-request-one", agentRunnerProcessStopped); err != nil {
		t.Fatal(err)
	}
	// A stopped process is terminal: nothing may move it forward again.
	if err := registry.advance("launch-request-one", agentRunnerProcessRunning); !errors.Is(err, ErrAgentRunnerProcessIdentity) {
		t.Fatalf("a stopped launch must never advance again, got %v", err)
	}
	if err := registry.advance("launch-request-missing", agentRunnerProcessRunning); !errors.Is(err, ErrAgentRunnerProcessMissing) {
		t.Fatalf("an unknown launch must be refused, got %v", err)
	}
	if err := registry.advance("launch-request-one", "resumed"); !errors.Is(err, ErrAgentRunnerProcessIdentity) {
		t.Fatalf("an unknown phase must be refused, got %v", err)
	}
	if handle, found := registry.lookup("launch-request-one"); !found || handle.phase != agentRunnerProcessStopped {
		t.Fatalf("lookup lost the handle: found=%v handle=%+v", found, handle)
	}
	registry.release("launch-request-one")
	if registry.len() != 0 {
		t.Fatalf("release left %d handles", registry.len())
	}
	if _, found := registry.lookup("launch-request-one"); found {
		t.Fatal("a released handle must not be found")
	}
}

func TestAgentRunnerProcessIdentityMirrorIsReadableByTheOperatorStopPath(t *testing.T) {
	managed := newAgentRunnerTestProcess(t)
	runRoot := t.TempDir()
	if err := writeAgentRunnerProcessIdentityMirror(runRoot, managed.Identity()); err != nil {
		t.Fatal(err)
	}
	mirrorPath := filepath.Join(runRoot, managedProcessIdentityFileName)
	data, err := os.ReadFile(mirrorPath)
	if err != nil {
		t.Fatal(err)
	}
	// The console stop path decodes the mirror strictly and rejects an incomplete
	// identity; this asserts the mirror satisfies exactly that reader.
	var identity guard.ProcessIdentity
	if err := evidence.DecodeStrictJSON(data, &identity); err != nil {
		t.Fatal(err)
	}
	if identity.PID <= 0 || identity.StartToken == "" {
		t.Fatalf("mirror identity is incomplete: %+v", identity)
	}
	if identity != managed.Identity() {
		t.Fatalf("mirror identity mismatch: %+v vs %+v", identity, managed.Identity())
	}
	encoded, err := agentRunnerProcessID(managed.Identity())
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != encoded {
		t.Fatalf("mirror bytes must equal the encoded process id: %q vs %q", string(data), encoded)
	}
	entries, err := os.ReadDir(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("the mirror write left temporary files behind: %v", entries)
	}
	if err := writeAgentRunnerProcessIdentityMirror("", managed.Identity()); !errors.Is(err, ErrAgentRunnerProcessIdentity) {
		t.Fatalf("an empty run root must be refused, got %v", err)
	}
}
