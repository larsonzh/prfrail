//go:build windows

package guard

// P84 fix round 8 (Rev 3b): the Windows half of the lease lock substrate.
//
// The previous shape named the lock file after the holder's pid and used CREATE_NEW as
// the "take the lock" step. That tied identity to a name each process derives from its
// own pid, so two processes created two different names and never blocked each other
// (tester round 4: a child holding the lock did not stop the parent's claim). The lock
// is therefore an OS-level exclusive lock on ONE contended name, "<record>.lock": every
// process opens the same file and asks the kernel for an exclusive byte range on it
// (LockFileEx). The kernel owns the exclusion, so there is no pid marker, no liveness
// query and no healing - and a process that dies, including one that crashes, has its
// lock released by the OS.

import (
	"errors"
	"syscall"
	"unsafe"
)

const (
	// LOCKFILE_EXCLUSIVE_LOCK makes the range exclusive, and LOCKFILE_FAIL_IMMEDIATELY
	// turns contention into an immediate ERROR_LOCK_VIOLATION rather than a wait inside
	// the kernel, so the caller keeps its own 5ms/ctx retry semantics.
	leaseLockFileExclusiveLock   = 0x00000002
	leaseLockFileFailImmediately = 0x00000001
	// leaseLockViolation is ERROR_LOCK_VIOLATION: another open handle already holds the
	// requested range.
	leaseLockViolation = syscall.Errno(33)
	// The locked range is one byte at offset 0 - enough to make the lock exclusive, and
	// the zero OVERLAPPED below means the range starts at the beginning of the file.
	leaseLockBytesLow  = 1
	leaseLockBytesHigh = 0
)

var (
	procLockFileEx   = syscall.NewLazyDLL("kernel32.dll").NewProc("LockFileEx")
	procUnlockFileEx = syscall.NewLazyDLL("kernel32.dll").NewProc("UnlockFileEx")
)

// leaseLock is one held resource lock: the shared lock file and the handle whose
// exclusive byte range the kernel grants to exactly one holder at a time.
type leaseLock struct {
	handle syscall.Handle
}

// createLeaseLock takes the resource lock, or reports that another holder owns it.
//
// OPEN_ALWAYS (not CREATE_NEW) is deliberate: the lock is the byte range, not the name,
// so every competitor opens the same existing-or-new file and the kernel decides. The
// share mode stays fully permissive because the exclusion comes from LockFileEx;
// FILE_SHARE_DELETE keeps the file unlinkable by external tooling without ever
// weakening the lock itself.
func createLeaseLock(recordPath string) (*leaseLock, error) {
	name, err := syscall.UTF16PtrFromString(leaseLockPath(recordPath))
	if err != nil {
		return nil, err
	}
	handle, err := syscall.CreateFile(
		name,
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_ALWAYS,
		0,
		0,
	)
	if err != nil {
		// P84 fix round 9 (N8): CreateFile cannot open a DIRECTORY with this access and
		// share shape, so a legacy mkdir-style lock directory left by an older run lands
		// here as an access failure. It is reported as its own sentinel - never deleted
		// and never waited on - instead of being mistaken for contention.
		if leaseLockPathIsDirectory(recordPath) {
			return nil, legacyLeaseLockDirectoryError(recordPath)
		}
		return nil, err
	}
	overlapped := new(syscall.Overlapped) // zero value: Offset = OffsetHigh = 0
	ok, _, callErr := procLockFileEx.Call(
		uintptr(handle),
		uintptr(leaseLockFileExclusiveLock|leaseLockFileFailImmediately),
		0,
		leaseLockBytesLow,
		leaseLockBytesHigh,
		uintptr(unsafe.Pointer(overlapped)),
	)
	if ok == 0 {
		_ = syscall.CloseHandle(handle)
		if errors.Is(callErr, leaseLockViolation) {
			return nil, errLeaseLockHeld
		}
		return nil, callErr
	}
	return &leaseLock{handle: handle}, nil
}

// release unlocks the byte range and closes the handle.
//
// Closing alone already releases a kernel lock; the explicit UnlockFileEx keeps the
// intent readable and its failure does not block the close, because releasing the
// handle is what actually frees the range for the next holder. The file itself is
// deliberately NOT removed: unlinking a lock file a waiter already has open would let a
// third process create a fresh file and hold a second, independent lock on the same
// resource.
func (lock *leaseLock) release() {
	overlapped := new(syscall.Overlapped)
	_, _, _ = procUnlockFileEx.Call(
		uintptr(lock.handle),
		0,
		leaseLockBytesLow,
		leaseLockBytesHigh,
		uintptr(unsafe.Pointer(overlapped)),
	)
	_ = syscall.CloseHandle(lock.handle)
}
