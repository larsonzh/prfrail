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

// pathFileID reports the identity of whatever path resolves to right now. It is
// the name-based half used by FileIDAt: a caller that declared a root needs to
// know whether that NAME still addresses the object it bound (P85 F02/F20 moved
// the binding itself off the name; this only answers the caller's-namespace
// question). The file is opened share-all so a declared root stays renamable
// while the observation runs.
func pathFileID(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(syscall.Handle(file.Fd()), &info); err != nil {
		return "", fmt.Errorf("identify %s: %w", path, err)
	}
	return fmt.Sprintf("win32:%d:%d:%d", info.VolumeSerialNumber, info.FileIndexHigh, info.FileIndexLow), nil
}
