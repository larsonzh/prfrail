//go:build unix

package guard

// P84 fix round 8 (Rev 3b): the POSIX half of the lease lock substrate.
//
// The previous shape named the lock file after the holder's pid and used O_EXCL as the
// "take the lock" step. That tied identity to a name each process derives from its own
// pid, so two processes created two different names and never blocked each other
// (tester round 4: a child holding the lock did not stop the parent's claim). The lock
// is therefore an advisory OS lock on ONE contended name, "<record>.lock": every
// process opens the same file and asks the kernel for LOCK_EX on it (flock(2)). The
// kernel owns the exclusion, so there is no pid marker, no liveness query and no
// healing - and a process that dies, including one that crashes, has its lock released
// by the kernel.

import (
	"errors"
	"os"
	"syscall"
)

// leaseLock is one held resource lock: the shared lock file and the descriptor whose
// flock the kernel grants to exactly one holder at a time.
type leaseLock struct {
	file *os.File
}

// createLeaseLock takes the resource lock, or reports that another holder owns it.
//
// O_CREATE without O_EXCL is deliberate: the lock is the flock, not the name, so every
// competitor opens the same existing-or-new file and the kernel decides. LOCK_NB turns
// contention into EWOULDBLOCK/EAGAIN instead of a blocking flock, so the caller keeps
// its own 5ms/ctx retry semantics.
func createLeaseLock(recordPath string) (*leaseLock, error) {
	file, err := os.OpenFile(leaseLockPath(recordPath), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		// P84 fix round 9 (N8): O_CREATE|O_RDWR on an existing directory fails with
		// EISDIR (ENOTDIR is the neighbouring shape), which is exactly how a legacy
		// mkdir-style lock directory left by an older run surfaces here. It is reported
		// as its own sentinel - never deleted and never waited on - rather than being
		// mistaken for contention or for a broken filesystem.
		if errors.Is(err, syscall.EISDIR) || errors.Is(err, syscall.ENOTDIR) {
			if leaseLockPathIsDirectory(recordPath) {
				return nil, legacyLeaseLockDirectoryError(recordPath)
			}
		}
		return nil, err
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = file.Close()
		// EWOULDBLOCK and EAGAIN may be the same errno; both mean "held elsewhere".
		if errors.Is(err, syscall.EWOULDBLOCK) || errors.Is(err, syscall.EAGAIN) {
			return nil, errLeaseLockHeld
		}
		return nil, err
	}
	return &leaseLock{file: file}, nil
}

// release drops the advisory lock and closes the descriptor.
//
// Closing alone already releases a flock; the explicit LOCK_UN keeps the intent
// readable. The file itself is deliberately NOT removed: unlinking a lock file a waiter
// already has open would let a third process create a fresh inode and hold a second,
// independent lock on the same resource.
func (lock *leaseLock) release() {
	_ = syscall.Flock(int(lock.file.Fd()), syscall.LOCK_UN)
	_ = lock.file.Close()
}
