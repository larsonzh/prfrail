package guard

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

// Test-only contract for the spawn-ownership slice
// (design: docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md
// §二 Q1/Q2, §三, §四 I1/I4).
//
// These tests intentionally reference API that does not exist yet:
// ProcessSpec.OnSpawn, ErrSpawnOwnership and platformProcess.resume. The expected
// RED state of this slice is a compile failure, not a test failure.
const (
	// spawnOwnershipHelperEnv marks the re-executed helper child.
	spawnOwnershipHelperEnv = "PROOFRAIL_SPAWN_OWNERSHIP_HELPER"
	// spawnOwnershipMarkerEnv names the file the helper writes as its first user
	// instruction, so the parent can tell "never executed" from "executed".
	spawnOwnershipMarkerEnv = "PROOFRAIL_SPAWN_OWNERSHIP_MARKER"
	// spawnOwnershipSettleBound bounds every wait in these tests: a hook that never
	// runs and a child that never resumes must fail instead of hanging.
	spawnOwnershipSettleBound = 10 * time.Second
	// spawnOwnershipBlockedWindow is how long the hook stays blocked before the
	// test concludes the child never executed a user instruction. It is far longer
	// than a resumed Go test binary needs to reach the marker write.
	spawnOwnershipBlockedWindow = 500 * time.Millisecond
)

// The platform split is part of the design contract: ownership binding (attach)
// and running the process must be separate steps, and resume must exist on every
// platform (Windows: NtResumeProcess, Linux: no-op). Referencing the method
// expression pins both the existence and the exact signature.
var _ func(*platformProcess) error = (*platformProcess).resume

// TestSpawnOwnershipHelper is the child half of the spawn-ownership tests. Its
// first user instruction writes the marker; it then keeps running with a bounded
// sleep. select{} is deliberately not used (DR-12): an unreaped helper that never
// returns makes Windows cleanup hang on an older toolchain.
func TestSpawnOwnershipHelper(t *testing.T) {
	if os.Getenv(spawnOwnershipHelperEnv) == "" {
		t.Skip("spawn ownership helper only runs as a child process")
	}
	marker := os.Getenv(spawnOwnershipMarkerEnv)
	if marker == "" {
		os.Exit(31)
	}
	if err := os.WriteFile(marker, []byte("executed\n"), 0o600); err != nil {
		os.Exit(32)
	}
	time.Sleep(30 * time.Second)
}

// spawnOwnershipHelperSpec builds a managed process whose first user instruction
// is observable through marker.
func spawnOwnershipHelperSpec(marker string) ProcessSpec {
	return ProcessSpec{
		Command: os.Args[0],
		Args:    []string{"-test.run=^TestSpawnOwnershipHelper$"},
		Env:     append(os.Environ(), spawnOwnershipHelperEnv+"=1", spawnOwnershipMarkerEnv+"="+marker),
	}
}

// spawnStartResult is the outcome of one StartManaged call parked in OnSpawn.
type spawnStartResult struct {
	process *ManagedProcess
	err     error
}

// blockedSpawn is one StartManaged call parked inside its OnSpawn hook. It lets a
// test observe the world while the hook is still blocked.
type blockedSpawn struct {
	identity ProcessIdentity
	release  chan struct{}
	done     chan struct{}
	once     sync.Once
	mu       sync.Mutex
	result   spawnStartResult
}

// beginBlockedSpawn starts a helper whose OnSpawn blocks until the test releases
// it, and returns once the hook has been entered. The cleanup releases the hook
// and stops whatever the start produced, so a failing assertion never leaves a
// suspended or running child behind.
func beginBlockedSpawn(t *testing.T, marker string) *blockedSpawn {
	t.Helper()
	blocked := &blockedSpawn{release: make(chan struct{}), done: make(chan struct{})}
	entered := make(chan ProcessIdentity, 1)
	spec := spawnOwnershipHelperSpec(marker)
	spec.OnSpawn = func(_ context.Context, identity ProcessIdentity) error {
		entered <- identity
		<-blocked.release
		return nil
	}
	go func() {
		process, err := StartManaged(context.Background(), spec)
		blocked.mu.Lock()
		blocked.result = spawnStartResult{process: process, err: err}
		blocked.mu.Unlock()
		close(blocked.done)
	}()
	t.Cleanup(func() {
		blocked.releaseHook()
		select {
		case <-blocked.done:
		case <-time.After(spawnOwnershipSettleBound):
			return
		}
		blocked.mu.Lock()
		process := blocked.result.process
		blocked.mu.Unlock()
		if process != nil {
			_, _ = process.Terminate(context.Background(), time.Second)
		}
	})
	select {
	case blocked.identity = <-entered:
	case <-time.After(spawnOwnershipSettleBound):
		t.Fatal("OnSpawn was never called: StartManaged does not expose the hook")
	}
	return blocked
}

func (blocked *blockedSpawn) releaseHook() {
	blocked.once.Do(func() { close(blocked.release) })
}

// finish releases the hook and returns the process StartManaged resumed.
func (blocked *blockedSpawn) finish(t *testing.T) *ManagedProcess {
	t.Helper()
	blocked.releaseHook()
	select {
	case <-blocked.done:
	case <-time.After(spawnOwnershipSettleBound):
		t.Fatal("StartManaged did not continue after OnSpawn returned")
	}
	blocked.mu.Lock()
	result := blocked.result
	blocked.mu.Unlock()
	if result.err != nil {
		t.Fatalf("StartManaged must resume once OnSpawn returns nil: %v", result.err)
	}
	return result.process
}

// waitForPath fails unless path appears within the settle bound.
func waitForPath(t *testing.T, path string) {
	t.Helper()
	deadline := time.Now().Add(spawnOwnershipSettleBound)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("%s was never created: the child never executed its first user instruction", path)
}

// TestOnSpawnFailureRollsBack pins I4: an OnSpawn error means "not started". The
// error must stay identifiable with errors.Is, the child must be stopped and no
// managed process may escape StartManaged.
func TestOnSpawnFailureRollsBack(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "spawned.marker")
	identityCh := make(chan ProcessIdentity, 1)
	spec := spawnOwnershipHelperSpec(marker)
	spec.OnSpawn = func(_ context.Context, identity ProcessIdentity) error {
		identityCh <- identity
		return ErrSpawnOwnership
	}
	process, err := StartManaged(context.Background(), spec)
	if err == nil {
		t.Fatal("StartManaged must fail when OnSpawn fails")
	}
	if process != nil {
		t.Fatal("a failed OnSpawn must not leave a managed process behind")
	}
	if !errors.Is(err, ErrSpawnOwnership) {
		t.Fatalf("the OnSpawn failure must stay identifiable with errors.Is(ErrSpawnOwnership), got %v", err)
	}
	var identity ProcessIdentity
	select {
	case identity = <-identityCh:
	case <-time.After(spawnOwnershipSettleBound):
		t.Fatal("OnSpawn was never called, so the rollback cannot be attributed to it")
	}
	if identity.PID <= 0 || identity.StartToken == "" {
		t.Fatalf("OnSpawn must receive the captured identity, got %+v", identity)
	}
	// The process was never allowed to escape a failed start: it must be gone.
	deadline := time.Now().Add(spawnOwnershipSettleBound)
	for {
		alive, aliveErr := ProcessAlive(identity)
		if aliveErr != nil {
			t.Fatalf("inspect the rolled-back process: %v", aliveErr)
		}
		if !alive {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("the process %d survived a failed OnSpawn: no rollback", identity.PID)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// TestStartManagedResumeIsSeparated pins the ordering half of I1 that every
// platform can express: the hook sees a complete identity, and the child only
// runs after the hook returned nil.
//
// On Windows the ordering is directly observable (the child is created suspended
// inside its job). On Linux there is no suspend step - the child is covered by
// pgid/Pdeathsig ownership instead - so "before resume" has no observable moment
// there; that part is registered as not directly testable (see the slice report
// and TestOnSpawnHookRunsBeforeResumeWindows).
func TestStartManagedResumeIsSeparated(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "spawned.marker")
	blocked := beginBlockedSpawn(t, marker)
	if blocked.identity.PID <= 0 || blocked.identity.StartToken == "" {
		t.Fatalf("OnSpawn must run after the identity is captured, got %+v", blocked.identity)
	}
	if runtime.GOOS == "windows" {
		// Give a wrongly resumed child ample time to write the marker.
		time.Sleep(spawnOwnershipBlockedWindow)
		if _, statErr := os.Stat(marker); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("the child executed a user instruction before OnSpawn returned (stat err=%v)", statErr)
		}
	}
	process := blocked.finish(t)
	if process.Identity() != blocked.identity {
		t.Fatalf("resume changed the identity: got %+v want %+v", process.Identity(), blocked.identity)
	}
	waitForPath(t, marker)
}
