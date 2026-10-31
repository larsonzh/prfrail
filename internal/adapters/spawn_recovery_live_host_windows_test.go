//go:build windows

package adapters

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"github.com/larsonzh/prfrail/internal/guard"
)

// Windows-only half of the spawn recovery contract
// (design: docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md
// §二 Q3, §五 risk 4).
//
// The live-host decision is LockFileEx-based and only exists on Windows, so the
// two tests that depend on it live behind this build tag. The platform-neutral
// recovery behaviour (I5 idempotence, token mismatch) is in spawn_recovery_test.go.
var (
	spawnRecoveryKernel32         = syscall.NewLazyDLL("kernel32.dll")
	spawnRecoveryProcLockFileEx   = spawnRecoveryKernel32.NewProc("LockFileEx")
	spawnRecoveryProcUnlockFileEx = spawnRecoveryKernel32.NewProc("UnlockFileEx")
)

const (
	lockFileExclusiveLock   = 0x00000002
	lockFileFailImmediately = 0x00000001
)

// holdStartingMarkerLock takes the exclusive byte-range lock the design gives the
// live host (LockFileEx on managed-process.starting.json) and returns an
// idempotent release func, also registered as cleanup.
func holdStartingMarkerLock(t *testing.T, path string) func() {
	t.Helper()
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatalf("open the starting marker %s: %v", path, err)
	}
	overlapped := new(syscall.Overlapped)
	locked, _, callErr := spawnRecoveryProcLockFileEx.Call(
		file.Fd(),
		lockFileExclusiveLock|lockFileFailImmediately,
		0,
		1,
		0,
		uintptr(unsafe.Pointer(overlapped)),
	)
	if locked == 0 {
		_ = file.Close()
		t.Fatalf("LockFileEx on %s: %v", path, callErr)
	}
	var once sync.Once
	release := func() {
		once.Do(func() {
			_, _, _ = spawnRecoveryProcUnlockFileEx.Call(file.Fd(), 0, 1, 0, uintptr(unsafe.Pointer(overlapped)))
			_ = file.Close()
		})
	}
	t.Cleanup(release)
	return release
}

// TestRecoverStartingCleansCrashedHostWindows pins I3/I5 on Windows: job-at-
// creation already killed the suspended tree with the crashed host, so the marker
// is proof of a dead start and is cleaned immediately - marker, slot and mirror -
// instead of waiting out the launch-slot freshness window. The freed slot must be
// directly reusable, and a second scan must be a no-op.
func TestRecoverStartingCleansCrashedHostWindows(t *testing.T) {
	runRoot := t.TempDir()
	dead := deadProcessIdentity(t)
	markerPath := writeStartingMarker(t, runRoot, testStartingMarker{
		LaunchID:   "launch-request-one",
		RequestID:  "request-one",
		Identity:   dead,
		StartToken: dead.StartToken,
		StartedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	})
	slot := filepath.Join(runRoot, managedProcessLaunchSlotFileName)
	if err := os.WriteFile(slot, []byte("launch-request-one\ntoken\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	launcher := recoveryLauncherFor(t, runRoot)
	reclaimed, err := launcher.RecoverStarting(context.Background())
	if err != nil {
		t.Fatalf("recover a crashed start: %v", err)
	}
	if !reclaimed {
		t.Fatal("a marker left by a crashed Windows host must be reclaimed")
	}
	if _, statErr := os.Stat(markerPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the reclaimed marker must be removed (stat err=%v)", statErr)
	}
	if _, statErr := os.Stat(slot); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the reclaimed launch slot must be released (stat err=%v)", statErr)
	}
	// The slot must be usable right away: a crashed start may not cost the run root
	// the freshness window.
	file, err := os.OpenFile(slot, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatalf("the reclaimed slot must be immediately reusable: %v", err)
	}
	_ = file.Close()
	_ = os.Remove(slot)
	again, err := launcher.RecoverStarting(context.Background())
	if err != nil {
		t.Fatalf("the second recovery must not fail: %v", err)
	}
	if again {
		t.Fatal("a second recovery of the same marker must be a no-op")
	}
}

// TestRecoverStartingRefusesLiveHost pins design risk 4: a marker whose
// LockFileEx lock is still held belongs to a live host, so the scan must refuse.
// Only once the lock is gone and the process is dead may the marker be reclaimed -
// the lock, not the marker's mere existence, is what decides.
func TestRecoverStartingRefusesLiveHost(t *testing.T) {
	runRoot := t.TempDir()
	live, stopLive := startRecoveryHelperProcess(t)
	markerPath := writeStartingMarker(t, runRoot, testStartingMarker{
		LaunchID:   "launch-request-one",
		RequestID:  "request-one",
		Identity:   live,
		StartToken: live.StartToken,
		StartedAt:  time.Now().UTC().Format(time.RFC3339Nano),
	})
	releaseLock := holdStartingMarkerLock(t, markerPath)
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
