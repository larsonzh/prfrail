//go:build windows

package adapters

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/larsonzh/prfrail/internal/guard"
)

// Windows-only half of the DR13-FIX-D3 ordering contract. The delete-sharing
// semantics that decide "can the published name be removed while a handle is
// open" only exist here, so the fail-closed refusal test and the share fact probe
// live behind this build tag.

const (
	// dr13ShareReadWriteOnly is FILE_SHARE_READ|FILE_SHARE_WRITE with
	// FILE_SHARE_DELETE deliberately left out: such a handle makes every delete of
	// the name fail with a sharing violation.
	dr13ShareReadWriteOnly = 0x00000001 | 0x00000002
)

// dr13OccupyStartingMarkerName opens path so that the name refuses deletion,
// which is how a third party (scanner, AV, filter) can hold the published marker.
func dr13OccupyStartingMarkerName(t *testing.T, path string) *os.File {
	t.Helper()
	pathp, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		t.Fatalf("encode the marker path %s: %v", path, err)
	}
	handle, createErr := syscall.CreateFile(
		pathp,
		startingMarkerGenericRead|startingMarkerGenericWrite,
		dr13ShareReadWriteOnly,
		nil,
		syscall.OPEN_EXISTING,
		startingMarkerFileAttributeNormal,
		0,
	)
	if createErr != nil {
		t.Fatalf("open %s without FILE_SHARE_DELETE: %v", path, createErr)
	}
	return os.NewFile(uintptr(handle), path)
}

// TestStartingMarkerReleaseFailClosedWhenTheNameIsOccupied pins the failure
// direction of the new order: when the delete is refused by another handle,
// release reports false and leaves the name in place - but it must still drop the
// host lock, because a leaked lock would block the run root until the process
// dies.
func TestStartingMarkerReleaseFailClosedWhenTheNameIsOccupied(t *testing.T) {
	runRoot := t.TempDir()
	claim, err := writeAgentRunnerStartingMarker(runRoot, newAgentRunnerStartingMarker(
		"launch-release-occupied",
		"request-release-occupied",
		guard.ProcessIdentity{PID: os.Getpid(), StartToken: "release-occupied-token"},
	))
	if err != nil {
		t.Fatalf("write the starting marker: %v", err)
	}
	occupier := dr13OccupyStartingMarkerName(t, claim.path)
	occupierClosed := false
	closeOccupier := func() {
		if occupierClosed {
			return
		}
		occupierClosed = true
		_ = occupier.Close()
	}
	t.Cleanup(closeOccupier)
	t.Cleanup(func() { _ = os.Remove(claim.path) })

	if removed := claim.release(); removed {
		t.Fatal("release must report false while another handle refuses the delete of the published name")
	}
	if _, statErr := os.Stat(claim.path); statErr != nil {
		t.Fatalf("a refused delete must leave the published name in place: %v", statErr)
	}

	// The lock is the fail-closed centre of gravity: remove failing may never keep
	// it held. Re-open the way the recovery clear does and prove the lock is free.
	reopened, openErr := openStartingMarkerForClear(claim.path)
	if openErr != nil {
		t.Fatalf("re-open the marker after a refused delete: %v", openErr)
	}
	locked, lockErr := tryLockStartingMarker(reopened)
	if lockErr != nil {
		unlockStartingMarker(reopened)
		_ = reopened.Close()
		t.Fatalf("lock the marker after a refused delete: %v", lockErr)
	}
	if !locked {
		unlockStartingMarker(reopened)
		_ = reopened.Close()
		t.Fatal("a refused delete must not leave the host lock held: the lock is still taken")
	}
	unlockStartingMarker(reopened)
	if closeErr := reopened.Close(); closeErr != nil {
		t.Fatalf("close the re-opened marker: %v", closeErr)
	}

	closeOccupier()
	if removeErr := os.Remove(claim.path); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
		t.Fatalf("clean up the residual marker: %v", removeErr)
	}
}

// TestStartingMarkerReleaseShareFactProbe records, without asserting, whether a
// Go os.OpenFile handle (the way RecoverStarting opens the marker) blocks a
// delete on this toolchain. The package comment at
// agent_runner_starting_marker_windows.go:97-101 claims it does; the design asks
// for the real behaviour to be pinned (GLM-2 R-N1). This test only logs, so a
// toolchain that changed the share mode cannot make the suite red.
func TestStartingMarkerReleaseShareFactProbe(t *testing.T) {
	runRoot := t.TempDir()
	path := filepath.Join(runRoot, "share-probe.txt")
	if err := os.WriteFile(path, []byte("probe"), 0o600); err != nil {
		t.Fatalf("write the probe file: %v", err)
	}
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatalf("open the probe file the way RecoverStarting does: %v", err)
	}
	removeErr := os.Remove(path)
	if removeErr == nil {
		t.Logf("os.OpenFile handle did NOT block os.Remove: a delete succeeded while the handle was open, so the windows.go comment (os.OpenFile asks for read/write sharing only) is stale for this toolchain")
	} else {
		var errno syscall.Errno
		if errors.As(removeErr, &errno) {
			t.Logf("os.OpenFile handle blocks os.Remove: err=%v errno=%d (0x%X), matching the windows.go comment", removeErr, uintptr(errno), uintptr(errno))
		} else {
			t.Logf("os.OpenFile handle blocks os.Remove: err=%v", removeErr)
		}
	}
	if closeErr := file.Close(); closeErr != nil {
		t.Fatalf("close the probe handle: %v", closeErr)
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("clean up the probe file: %v", err)
	}
}
