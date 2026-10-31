//go:build unix

package adapters

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/guard"
)

// Unix half of the spawn recovery contract
// (design: docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md
// section 2 Q3, section 5 risk 4; revision round F-2).
//
// The live-host decision is flock-based here and LockFileEx-based on Windows, so the
// two runtime paths that only exist on one platform live behind their own build tag.
// The platform-neutral recovery behaviour (I5 idempotence, token mismatch) stays in
// spawn_recovery_test.go.
//
// What this file adds over the Windows half is the Linux-only consequence of a free
// lock: the kernel does not kill the child with its host, so an orphan of a crashed
// host is still running when the next scan looks - and the scan has to adopt and stop
// it through the persisted identity.

// spawnRecoverySettleBound bounds every wait in the unix half of these tests. It lives
// here rather than in the shared const block because only this build-tagged file uses
// it, so a Windows-only build would otherwise carry an unused constant.
const spawnRecoverySettleBound = 10 * time.Second

// holdUnixStartingMarkerLock takes the exclusive host lock the design gives the live
// host (flock on managed-process.starting.json) and returns an idempotent release
// func, also registered as cleanup.
//
// flock is taken on a descriptor of its own: a lock is bound to the open file
// description, so a second open in this process conflicts with it exactly as another
// host's lock would, which is the condition the scan has to refuse.
func holdUnixStartingMarkerLock(t *testing.T, path string) func() {
	t.Helper()
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatalf("open the starting marker %s: %v", path, err)
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = file.Close()
		t.Fatalf("flock %s: %v", path, err)
	}
	var once sync.Once
	release := func() {
		once.Do(func() {
			_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
			_ = file.Close()
		})
	}
	t.Cleanup(release)
	return release
}

// startDetachedRecoveryOrphan starts the shared long-lived helper as a session
// leader, so it forms its own process group and outlives the test process - exactly
// the orphan a crashed Linux host leaves behind, and the same anchor guard uses for a
// spawned child. The returned identity is the process's real platform identity, read
// from the running process, so the marker names a live process the way a crashed
// host's marker would.
func startDetachedRecoveryOrphan(t *testing.T) guard.ProcessIdentity {
	t.Helper()
	cmd := exec.Command(os.Args[0], "-test.run=^TestSpawnRecoveryHelper$")
	cmd.Env = append(os.Environ(), spawnRecoveryHelperEnv+"=1")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the detached orphan: %v", err)
	}
	pid := cmd.Process.Pid
	var once sync.Once
	reap := func() {
		once.Do(func() {
			// Whatever the assertions concluded, the helper must not outlive this test:
			// kill it and reap the zombie so the run leaves no process behind.
			_ = syscall.Kill(pid, syscall.SIGKILL)
			_, _ = cmd.Process.Wait()
		})
	}
	t.Cleanup(reap)
	identity, err := guard.InspectProcess(pid)
	if err != nil {
		t.Fatalf("read the identity of the orphan %d: %v", pid, err)
	}
	// Bounded settle: the marker may only name a process that is really running,
	// otherwise nothing below proves that recovery stopped anything.
	deadline := time.Now().Add(spawnRecoverySettleBound)
	for {
		alive, aliveErr := guard.ProcessAlive(identity)
		if aliveErr != nil {
			t.Fatalf("inspect the orphan %d: %v", pid, aliveErr)
		}
		if alive {
			return identity
		}
		if time.Now().After(deadline) {
			t.Fatalf("the orphan %d never became alive", pid)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestRecoverStartingRefusesLiveHostUnix pins design risk 4 on Linux: a marker whose
// flock is still held belongs to a live host, so the scan must refuse - keep the
// marker and touch neither the process it names nor anything else. Only once the lock
// is gone and the process is dead may the marker be reclaimed, which is what makes the
// lock, not the marker's mere existence, the decision.
func TestRecoverStartingRefusesLiveHostUnix(t *testing.T) {
	runRoot := t.TempDir()
	live, stopLive := startRecoveryHelperProcess(t)
	markerPath := writeStartingMarker(t, runRoot, testStartingMarker{
		LaunchID:   "launch-request-one",
		RequestID:  "request-one",
		Identity:   live,
		StartToken: live.StartToken,
		StartedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	})
	releaseLock := holdUnixStartingMarkerLock(t, markerPath)
	launcher := recoveryLauncherFor(t, runRoot)
	// An error is an acceptable refusal as well; reclaimed must never be true.
	reclaimed, _ := launcher.RecoverStarting(context.Background())
	if reclaimed {
		t.Fatal("a marker whose lock is held by a live host must never be reclaimed")
	}
	if _, statErr := os.Stat(markerPath); statErr != nil {
		t.Fatalf("a refused scan must keep the marker for the live host: %v", statErr)
	}
	if alive, aliveErr := guard.ProcessAlive(live); aliveErr != nil || !alive {
		t.Fatalf("a refused scan must not touch the live process: alive=%v err=%v", alive, aliveErr)
	}
	releaseLock()
	stopLive()
	reclaimed, err := launcher.RecoverStarting(context.Background())
	if err != nil {
		t.Fatalf("recover the dead start once the lock is gone: %v", err)
	}
	if !reclaimed {
		t.Fatal("with the host lock released and the process dead the marker must be reclaimed")
	}
	if _, statErr := os.Stat(markerPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the reclaimed marker must be removed (stat err=%v)", statErr)
	}
}

// TestRecoverStartingAdoptsOrphanUnix pins the Linux reclaim path: a free-locked
// marker that names a still running orphan must be adopted through the persisted
// identity and stopped, and the marker and the slot the crashed start left behind must
// both be released, so the next launch does not wait out the freshness window. A
// second scan is a no-op.
func TestRecoverStartingAdoptsOrphanUnix(t *testing.T) {
	runRoot := t.TempDir()
	const launchID = "launch-orphan-one"
	orphan := startDetachedRecoveryOrphan(t)
	markerPath := writeStartingMarker(t, runRoot, testStartingMarker{
		LaunchID:   launchID,
		RequestID:  "request-orphan-one",
		Identity:   orphan,
		StartToken: orphan.StartToken,
		StartedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	})
	slot := filepath.Join(runRoot, managedProcessLaunchSlotFileName)
	if err := os.WriteFile(slot, []byte(launchID+"\ntoken\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	launcher := recoveryLauncherFor(t, runRoot)
	// The marker is published without its host lock, which is the state a crashed host
	// leaves behind: free lock, live orphan.
	reclaimed, err := launcher.RecoverStarting(context.Background())
	if err != nil {
		t.Fatalf("adopt the orphan of a crashed host: %v", err)
	}
	if !reclaimed {
		t.Fatal("a free-locked marker naming a live orphan must be reclaimed on Linux")
	}
	alive, aliveErr := guard.ProcessAlive(orphan)
	if aliveErr != nil {
		t.Fatalf("inspect the adopted orphan: %v", aliveErr)
	}
	if alive {
		t.Fatalf("recovery left the orphan %d running", orphan.PID)
	}
	if _, statErr := os.Stat(markerPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the reclaimed marker must be removed (stat err=%v)", statErr)
	}
	if _, statErr := os.Stat(slot); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the reclaimed launch slot must be released (stat err=%v)", statErr)
	}
	again, err := launcher.RecoverStarting(context.Background())
	if err != nil {
		t.Fatalf("the second recovery must not fail: %v", err)
	}
	if again {
		t.Fatal("a second recovery of the same marker must be a no-op")
	}
}
