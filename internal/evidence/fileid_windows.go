//go:build windows

package evidence

// T3 (P83 fix round 6): a digest alone does not identify the object a spawn will
// run, because two distinct files can be swapped for one another while keeping
// size, modification time and mode. Windows reports a stable identity per open
// handle - the volume serial number plus the file index - which is the object
// identity the approval has to travel with.

import (
	"fmt"
	"os"
	"syscall"
)

// boundFileID returns the platform identity of the object the handle addresses.
// It is derived from the handle, never from the path, so it describes exactly the
// bytes a bound read would see.
func boundFileID(file *os.File) (string, error) {
	if file == nil {
		return "", fmt.Errorf("%w: no handle to identify", ErrBoundFileClosed)
	}
	handle := syscall.Handle(file.Fd())
	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(handle, &info); err != nil {
		return "", fmt.Errorf("identify the bound file: %w", err)
	}
	return fmt.Sprintf("win32:%d:%d:%d", info.VolumeSerialNumber, info.FileIndexHigh, info.FileIndexLow), nil
}
