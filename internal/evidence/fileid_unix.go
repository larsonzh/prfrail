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
