//go:build !windows

package tickets

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
)

// acquireDurableLock takes the scope-wide exclusive lock on the stable side-car
// lock file. The lock file is never replaced and never removed, so every
// process in the scope contends on the same file identity. Read-only callers
// pass create=false, which fails closed when the lock file is missing instead of
// creating it. The wait uses non-blocking flock plus backoff so ctx can cancel it.
func acquireDurableLock(ctx context.Context, path string, create bool) (func() error, error) {
	if ctx == nil {
		return nil, fmt.Errorf("%w: durable lock requires a context", ErrInvalidCostLedger)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
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
	fd := int(file.Fd())
	for {
		err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
			break
		}
		if err == syscall.EINTR {
			continue
		}
		if err != syscall.EWOULDBLOCK && err != syscall.EAGAIN {
			_ = file.Close()
			return nil, err
		}
		if backoffErr := durableLockBackoff(ctx); backoffErr != nil {
			_ = file.Close()
			return nil, backoffErr
		}
	}
	return func() error {
		unlockErr := syscall.Flock(fd, syscall.LOCK_UN)
		closeErr := file.Close()
		if unlockErr != nil {
			return unlockErr
		}
		return closeErr
	}, nil
}

// syncDurableDir is the parent directory barrier after an atomic replace.
func syncDurableDir(dir string) error {
	handle, err := os.Open(dir)
	if err != nil {
		return err
	}
	syncErr := handle.Sync()
	closeErr := handle.Close()
	if syncErr != nil {
		return syncErr
	}
	return closeErr
}
