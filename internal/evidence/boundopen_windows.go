//go:build windows

package evidence

// P42-T027B-EVIDENCE-BIND: the Windows bound read must keep the artifact
// renamable and removable while it stays open. Go's syscall.Open shares read and
// write but not delete, so os.Open alone would let an open handle pin the name and
// hide a replacement from the same-handle read this primitive exists to provide.
// The handle is opened with FILE_SHARE_READ|WRITE|DELETE instead.

import (
	"os"
	"syscall"
)

const boundFileShareAll = syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE | syscall.FILE_SHARE_DELETE

// openBoundFileHandle opens one artifact handle with delete sharing. The error is
// a *os.PathError, so errors.Is(err, fs.ErrNotExist) keeps working.
func openBoundFileHandle(path string) (*os.File, error) {
	name, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, boundFileShareAll, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	return os.NewFile(uintptr(handle), path), nil
}
