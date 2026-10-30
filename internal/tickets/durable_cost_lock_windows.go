//go:build windows

package tickets

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	durableLockExclusiveFlag = 0x00000002
	// durableLockFailImmediately makes LockFileEx non-blocking so the wait can be
	// polled and cancelled through ctx.
	durableLockFailImmediately = 0x00000001
	durableLockViolationErrno  = syscall.Errno(33)
)

var (
	durableKernel32     = syscall.NewLazyDLL("kernel32.dll")
	durableLockFileEx   = durableKernel32.NewProc("LockFileEx")
	durableUnlockFileEx = durableKernel32.NewProc("UnlockFileEx")
)

// acquireDurableLock takes the scope-wide exclusive lock on the stable side-car
// lock file via LockFileEx over a fixed one-byte range. The lock file is never
// replaced, removed or inherited by child processes; the kernel releases the
// lock when the owning process exits. Read-only callers pass create=false so a
// missing lock file fails closed, and the non-blocking retry loop honours ctx.
func acquireDurableLock(ctx context.Context, path string, create bool) (func() error, error) {
	flags := os.O_RDWR
	if create {
		flags |= os.O_CREATE
	}
	file, err := os.OpenFile(path, flags, 0o600)
	if err != nil {
		if !create && errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: durable lock file %q is missing", ErrInvalidCostLedger, path)
		}
		return nil, err
	}
	lockFlags := uintptr(durableLockExclusiveFlag | durableLockFailImmediately)
	for {
		overlapped := new(syscall.Overlapped)
		acquired, _, callErr := durableLockFileEx.Call(
			file.Fd(),
			lockFlags,
			0,
			1,
			0,
			uintptr(unsafe.Pointer(overlapped)),
		)
		if acquired != 0 {
			return func() error {
				released, _, unlockErr := durableUnlockFileEx.Call(
					file.Fd(),
					0,
					1,
					0,
					uintptr(unsafe.Pointer(overlapped)),
				)
				closeErr := file.Close()
				if released == 0 && unlockErr != nil && unlockErr != syscall.Errno(0) {
					return unlockErr
				}
				return closeErr
			}, nil
		}
		if callErr != nil && callErr != syscall.Errno(0) && callErr != durableLockViolationErrno {
			_ = file.Close()
			return nil, callErr
		}
		if backoffErr := durableLockBackoff(ctx); backoffErr != nil {
			_ = file.Close()
			return nil, backoffErr
		}
	}
}

// syncDurableDir is the parent directory barrier after an atomic replace. The
// Windows backend keeps the platform durability grade "unproven" (see
// docs/CONTRACTS.md §7), so there is no portable directory handle flush here.
func syncDurableDir(dir string) error {
	_ = dir
	return nil
}
