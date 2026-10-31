//go:build windows

package adapters

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/larsonzh/prfrail/internal/guard"
)

// The Windows liveness proof of a starting marker is a byte-range lock on the
// marker file itself (LockFileEx, byte 0, one byte, exclusive). The kernel drops
// it when the holding process dies, so a free lock is exactly "the host is gone"
// - no clock, no freshness window
// (docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md
// section 2 Q3).
const (
	startingMarkerLockExclusive        = 0x00000002
	startingMarkerLockFailImmediately  = 0x00000001
	startingMarkerGenericRead          = 0x80000000
	startingMarkerGenericWrite         = 0x40000000
	startingMarkerShareRead            = 0x00000001
	startingMarkerShareWrite           = 0x00000002
	startingMarkerShareDelete          = 0x00000004
	startingMarkerCreateNew            = 1
	startingMarkerFileAttributeNormal  = 0x00000080
	startingMarkerErrorLockViolation   = syscall.Errno(33)
	startingMarkerErrorFileExists      = syscall.Errno(80)
	startingMarkerErrorAlreadyExists   = syscall.Errno(183)
	startingMarkerTempFileAttemptLimit = 32
)

var (
	startingMarkerKernel32         = syscall.NewLazyDLL("kernel32.dll")
	startingMarkerProcLockFileEx   = startingMarkerKernel32.NewProc("LockFileEx")
	startingMarkerProcUnlockFileEx = startingMarkerKernel32.NewProc("UnlockFileEx")
)

// tryLockStartingMarker takes the exclusive host lock without waiting. It reports
// (false, nil) when another host holds it, which is a refusal, not a failure.
func tryLockStartingMarker(file *os.File) (bool, error) {
	if file == nil {
		return false, os.ErrInvalid
	}
	overlapped := new(syscall.Overlapped)
	locked, _, callErr := startingMarkerProcLockFileEx.Call(
		file.Fd(),
		startingMarkerLockExclusive|startingMarkerLockFailImmediately,
		0,
		1,
		0,
		uintptr(unsafe.Pointer(overlapped)),
	)
	if locked == 0 {
		if errors.Is(callErr, startingMarkerErrorLockViolation) {
			return false, nil
		}
		return false, callErr
	}
	return true, nil
}

// unlockStartingMarker releases a lock taken by tryLockStartingMarker. Closing the
// file releases it as well, so this is only the explicit form.
func unlockStartingMarker(file *os.File) {
	if file == nil {
		return
	}
	overlapped := new(syscall.Overlapped)
	_, _, _ = startingMarkerProcUnlockFileEx.Call(file.Fd(), 0, 1, 0, uintptr(unsafe.Pointer(overlapped)))
}

// reclaimStartingProcess decides what a free-locked starting marker means for the
// process it names.
//
// On Windows the ownership anchor is the job the child was created in, and that job
// dies with its host (kill-on-job-close). A host that is gone therefore took the
// whole suspended tree with it, so nothing may be killed here - and a process that
// still answers with the marker's identity cannot be that child. The reclaim is
// proven by the anchor, not by a cleanup action.
func reclaimStartingProcess(_ context.Context, _ guard.ProcessIdentity, _ time.Duration, alive bool) (bool, error) {
	if alive {
		// Fail closed: never kill a process this anchor cannot claim as ours.
		return false, nil
	}
	return true, nil
}

// openStartingMarkerForClear re-opens a published marker so the cleanup path can
// verify the record under its host lock and remove it while that proof is still held.
//
// The handle must share delete: Go's os.OpenFile asks for read and write sharing only,
// and Windows then refuses to delete a name any handle keeps without
// FILE_SHARE_DELETE. The recovery clear would silently fail on every marker without
// this handle (review F-1).
func openStartingMarkerForClear(path string) (*os.File, error) {
	pathp, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	handle, createErr := syscall.CreateFile(
		pathp,
		startingMarkerGenericRead|startingMarkerGenericWrite,
		startingMarkerShareRead|startingMarkerShareWrite|startingMarkerShareDelete,
		nil,
		syscall.OPEN_EXISTING,
		startingMarkerFileAttributeNormal,
		0,
	)
	if createErr != nil {
		return nil, createErr
	}
	return os.NewFile(uintptr(handle), path), nil
}

// createLockedStartingMarkerFile creates a temp file that is already locked, with a
// share mode that stays renameable: the share mode allows the delete that a rename
// performs, so the marker can be renamed into place while the lock stays held.
func createLockedStartingMarkerFile(dir string) (*os.File, error) {
	for attempt := 0; attempt < startingMarkerTempFileAttemptLimit; attempt++ {
		probe, err := os.CreateTemp(dir, "starting-*.tmp")
		if err != nil {
			return nil, err
		}
		name := probe.Name()
		if err := probe.Close(); err != nil {
			_ = os.Remove(name)
			return nil, err
		}
		if err := os.Remove(name); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		path, err := syscall.UTF16PtrFromString(name)
		if err != nil {
			return nil, err
		}
		handle, createErr := syscall.CreateFile(
			path,
			startingMarkerGenericRead|startingMarkerGenericWrite,
			startingMarkerShareRead|startingMarkerShareWrite|startingMarkerShareDelete,
			nil,
			startingMarkerCreateNew,
			startingMarkerFileAttributeNormal,
			0,
		)
		if createErr != nil {
			if errors.Is(createErr, startingMarkerErrorFileExists) || errors.Is(createErr, startingMarkerErrorAlreadyExists) {
				// Another writer took the same name between the probe and the create.
				continue
			}
			return nil, createErr
		}
		file := os.NewFile(uintptr(handle), name)
		// A (false, nil) answer is a refusal, not a success: publishing a claim on a
		// file another host still holds would hand out an unlocked marker and defeat the
		// whole live-host rule (review F-3).
		locked, lockErr := tryLockStartingMarker(file)
		if lockErr != nil {
			_ = file.Close()
			_ = os.Remove(name)
			return nil, fmt.Errorf("lock the starting marker: %w", lockErr)
		}
		if !locked {
			_ = file.Close()
			_ = os.Remove(name)
			return nil, errors.New("the starting marker file is already locked by another host")
		}
		return file, nil
	}
	return nil, errors.New("cannot create a unique starting marker file")
}
