//go:build !windows

package adapters

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/larsonzh/prfrail/internal/guard"
)

// The Unix liveness proof of a starting marker is an advisory lock on the marker
// file itself (flock, exclusive, non-blocking). The kernel drops it when the
// holding process dies, so a free lock is exactly "the host is gone" - the same
// decision the Windows byte-range lock makes
// (docs/validation/evidence/P41-T027B-SPAWN-OWNERSHIP/01-arch-design/review.md
// section 2 Q3).
//
// flock is used instead of a lock file because it is inherited by no child, kept
// across a rename and released by the kernel on any kind of host death.

// tryLockStartingMarker takes the exclusive host lock without waiting. It reports
// (false, nil) when another host holds it, which is a refusal, not a failure.
func tryLockStartingMarker(file *os.File) (bool, error) {
	if file == nil {
		return false, os.ErrInvalid
	}
	err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, syscall.EWOULDBLOCK), errors.Is(err, syscall.EAGAIN):
		return false, nil
	default:
		return false, err
	}
}

// unlockStartingMarker releases a lock taken by tryLockStartingMarker. Closing the
// file releases it as well, so this is only the explicit form.
func unlockStartingMarker(file *os.File) {
	if file == nil {
		return
	}
	_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
}

// reclaimStartingProcess decides what a free-locked starting marker means for the
// process it names.
//
// Linux has no job object: the process group is the anchor and nothing kills the
// child when the host dies, so an orphan of a crashed host may still be running.
// It is adopted through the persisted identity - process-group kill plus
// start-token comparison, both already reviewed in the guard - and the reclaim is
// proven by the bounded stop verification.
func reclaimStartingProcess(ctx context.Context, identity guard.ProcessIdentity, grace time.Duration, alive bool) (bool, error) {
	if !alive {
		return true, nil
	}
	if _, err := guard.StopProcessIdentity(ctx, identity, grace); err != nil {
		return false, err
	}
	return true, nil
}

// openStartingMarkerForClear re-opens a published marker so the cleanup path can
// verify the record under its host lock and remove it while that proof is still held.
// Unlinking an open file is always allowed on Unix, so a plain open is enough.
func openStartingMarkerForClear(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, 0)
}

// createLockedStartingMarkerFile creates a temp file that is already locked. A
// rename keeps the inode, so the lock stays held on the published marker path too.
func createLockedStartingMarkerFile(dir string) (*os.File, error) {
	file, err := os.CreateTemp(dir, "starting-*.tmp")
	if err != nil {
		return nil, err
	}
	// A (false, nil) answer is a refusal, not a success: publishing a claim on a file
	// another host still holds would hand out an unlocked marker and defeat the whole
	// live-host rule (review F-3).
	locked, lockErr := tryLockStartingMarker(file)
	if lockErr != nil {
		name := file.Name()
		_ = file.Close()
		_ = os.Remove(name)
		return nil, fmt.Errorf("lock the starting marker: %w", lockErr)
	}
	if !locked {
		name := file.Name()
		_ = file.Close()
		_ = os.Remove(name)
		return nil, errors.New("the starting marker file is already locked by another host")
	}
	return file, nil
}
