package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/guard"
)

// Test-only contract for the spawn-ownership recovery slice
// (design: docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md
// §二 Q3, §三, §四 I3/I5).
//
// (*AgentRunnerPinnedCLILauncher).RecoverStarting does not exist yet, so the
// expected RED state of this slice is a compile failure, not a test failure.
const (
	// testStartingMarkerFileName is the durable starting marker the design fixes
	// for the run root (§二 Q1, §六). It is a restatement of the protocol name on
	// purpose: a rename in production without a design update must break this test.
	testStartingMarkerFileName = "managed-process.starting.json"
	// spawnRecoveryHelperEnv marks the re-executed long-lived helper child.
	spawnRecoveryHelperEnv = "PROOFRAIL_SPAWN_RECOVERY_HELPER"
	// spawnRecoverySettleBound bounds every wait in these tests.
	spawnRecoverySettleBound = 10 * time.Second
)

// testStartingMarker is the durable starting-phase record from the design: the
// host writes it atomically before the child is resumed, so a crash leaves proof
// that a start was in flight.
type testStartingMarker struct {
	LaunchID   string                `json:"launchId"`
	RequestID  string                `json:"requestId"`
	Identity   guard.ProcessIdentity `json:"identity"`
	StartToken string                `json:"startToken"`
	StartedAt  string                `json:"startedAt"`
}

// TestSpawnRecoveryHelper is a long-lived child used to obtain a real process
// identity for the recovery tests. It sleeps with a bound (never select{}), so an
// unreaped helper cannot hang an older Windows toolchain (DR-12).
func TestSpawnRecoveryHelper(t *testing.T) {
	if os.Getenv(spawnRecoveryHelperEnv) == "" {
		t.Skip("spawn recovery helper only runs as a child process")
	}
	time.Sleep(time.Minute)
}

// writeStartingMarker publishes a marker atomically (temp file plus rename), the
// protocol the design fixes for the starting phase.
func writeStartingMarker(t *testing.T, runRoot string, marker testStartingMarker) string {
	t.Helper()
	if err := os.MkdirAll(runRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(marker)
	if err != nil {
		t.Fatal(err)
	}
	temporary, err := os.CreateTemp(runRoot, "starting-*.tmp")
	if err != nil {
		t.Fatal(err)
	}
	name := temporary.Name()
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		t.Fatal(err)
	}
	if err := temporary.Sync(); err != nil {
		_ = temporary.Close()
		t.Fatal(err)
	}
	if err := temporary.Close(); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(runRoot, testStartingMarkerFileName)
	if err := os.Rename(name, target); err != nil {
		_ = os.Remove(name)
		t.Fatal(err)
	}
	return target
}

// startRecoveryHelperProcess starts one live helper and returns its identity plus
// an idempotent stop func (also registered as cleanup).
func startRecoveryHelperProcess(t *testing.T) (guard.ProcessIdentity, func()) {
	t.Helper()
	process, err := guard.StartManaged(context.Background(), guard.ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=^TestSpawnRecoveryHelper$"},
		Env:     append(os.Environ(), spawnRecoveryHelperEnv+"=1"),
	})
	if err != nil {
		t.Fatalf("start the recovery helper: %v", err)
	}
	identity := process.Identity()
	stop := func() {
		_, _ = process.Terminate(context.Background(), 2*time.Second)
	}
	t.Cleanup(stop)
	return identity, stop
}

// deadProcessIdentity returns an identity that no longer names a live process,
// which is what a crashed host leaves in the marker on Windows (job-at-creation
// kills the tree with the host).
func deadProcessIdentity(t *testing.T) guard.ProcessIdentity {
	t.Helper()
	identity, stop := startRecoveryHelperProcess(t)
	stop()
	alive, err := guard.ProcessAlive(identity)
	if err != nil {
		t.Fatalf("inspect the stopped helper: %v", err)
	}
	if alive {
		t.Fatalf("helper %d is still alive after Terminate", identity.PID)
	}
	return identity
}

// recoveryLauncherFor builds a launcher bound to runRoot; recovery never needs a
// spawned launch, only the run root the markers live in.
func recoveryLauncherFor(t *testing.T, runRoot string) *AgentRunnerPinnedCLILauncher {
	t.Helper()
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), replayStoreForRunRoot(t, runRoot))
	if err != nil {
		t.Fatal(err)
	}
	return launcher
}

// TestRecoverStartingIdempotent pins I5: recovering the same crashed start twice
// reclaims it once, and the second scan is a no-op rather than a second kill or a
// second cleanup.
func TestRecoverStartingIdempotent(t *testing.T) {
	runRoot := t.TempDir()
	dead := deadProcessIdentity(t)
	markerPath := writeStartingMarker(t, runRoot, testStartingMarker{
		LaunchID:   "launch-request-one",
		RequestID:  "request-one",
		Identity:   dead,
		StartToken: dead.StartToken,
		StartedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	})
	launcher := recoveryLauncherFor(t, runRoot)
	reclaimed, err := launcher.RecoverStarting(context.Background())
	if err != nil {
		t.Fatalf("recover a crashed start: %v", err)
	}
	if !reclaimed {
		t.Fatal("a marker left by a dead host must be reclaimed")
	}
	if _, statErr := os.Stat(markerPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the reclaimed marker must be removed (stat err=%v)", statErr)
	}
	again, err := launcher.RecoverStarting(context.Background())
	if err != nil {
		t.Fatalf("the second recovery must not fail: %v", err)
	}
	if again {
		t.Fatal("a second recovery of the same marker must be a no-op (I5)")
	}
}

// TestRecoverStartingTokenMismatchFailsClosed pins the anti-kill half of I5: a
// marker whose start token does not match the process now holding that pid is a
// reused pid, so recovery must not kill it and must not drop the identity mirror.
func TestRecoverStartingTokenMismatchFailsClosed(t *testing.T) {
	runRoot := t.TempDir()
	live, _ := startRecoveryHelperProcess(t)
	if err := writeAgentRunnerProcessIdentityMirror(runRoot, live); err != nil {
		t.Fatalf("publish the identity mirror: %v", err)
	}
	// Same pid, different start token: exactly what pid reuse looks like.
	reused := guard.ProcessIdentity{PID: live.PID, StartToken: "ffffffffffffffff"}
	markerPath := writeStartingMarker(t, runRoot, testStartingMarker{
		LaunchID:   "launch-request-one",
		RequestID:  "request-one",
		Identity:   reused,
		StartToken: reused.StartToken,
		StartedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	})
	launcher := recoveryLauncherFor(t, runRoot)
	reclaimed, err := launcher.RecoverStarting(context.Background())
	if err != nil {
		t.Fatalf("a token mismatch is a fail-closed outcome, not an error: %v", err)
	}
	if reclaimed {
		t.Fatal("a marker whose token does not match the live pid must not be reclaimed")
	}
	alive, aliveErr := guard.ProcessAlive(live)
	if aliveErr != nil {
		t.Fatalf("inspect the live process: %v", aliveErr)
	}
	if !alive {
		t.Fatalf("recovery killed pid %d, which the marker only matched by a reused pid", live.PID)
	}
	if _, statErr := os.Stat(markerPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the stale marker must be cleared even when the kill is refused (stat err=%v)", statErr)
	}
	if _, statErr := os.Stat(filepath.Join(runRoot, managedProcessIdentityFileName)); statErr != nil {
		t.Fatalf("fail-closed keeps the identity mirror for the operator stop path: %v", statErr)
	}
}
