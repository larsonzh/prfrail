//go:build windows

package adapters

import (
	"syscall"
)

const fileAttributeReparsePoint = 0x400

// replayPathComponentUnsafe reports whether an existing path component is a
// reparse point (symlink, junction, mounted folder, or similar). Missing
// components are reported as safe because only already-existing components can
// redirect reads or writes.
func replayPathComponentUnsafe(path string) (bool, error) {
	name, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}
	attributes, err := syscall.GetFileAttributes(name)
	if err != nil {
		if err == syscall.ERROR_FILE_NOT_FOUND || err == syscall.ERROR_PATH_NOT_FOUND {
			return false, nil
		}
		return false, err
	}
	return attributes&fileAttributeReparsePoint != 0, nil
}
