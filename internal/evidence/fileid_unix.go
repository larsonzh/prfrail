//go:build !windows

package evidence

// T3 (P83 fix round 6): a digest alone does not identify the object a spawn will
// run, because two distinct files can be swapped for one another while keeping
// size, modification time and mode. Unix reports a stable identity per open
// descriptor - the device number plus the inode - which is the object identity
// the approval has to travel with.

import (
	"fmt"
	"os"
	"syscall"
)

// boundFileID returns the platform identity of the object the descriptor
// addresses. It is taken from fstat(2) on the descriptor, never from the path, so
// it describes exactly the bytes a bound read would see.
func boundFileID(file *os.File) (string, error) {
	if file == nil {
		return "", fmt.Errorf("%w: no handle to identify", ErrBoundFileClosed)
	}
	var info syscall.Stat_t
	if err := syscall.Fstat(int(file.Fd()), &info); err != nil {
		return "", fmt.Errorf("identify the bound file: %w", err)
	}
	return fmt.Sprintf("unix:%d:%d", info.Dev, info.Ino), nil
}

// pathFileID reports the identity of whatever path resolves to right now. It is
// the name-based half used by FileIDAt: a caller that declared a root needs to
// know whether that NAME still addresses the object it bound (P85 F02/F20 moved
// the binding itself off the name; this only answers the caller's-namespace
// question).
func pathFileID(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "", fmt.Errorf("evidence: %s carries no syscall stat", path)
	}
	return fmt.Sprintf("unix:%d:%d", stat.Dev, stat.Ino), nil
}
