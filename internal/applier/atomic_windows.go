//go:build windows

package applier

import (
	"syscall"
	"time"
	"unsafe"
)

const (
	moveFileReplaceExisting = 0x1
	moveFileWriteThrough    = 0x8

	replaceAttempts      = 10
	replaceRetryInterval = 5 * time.Millisecond
)

const (
	errorAccessDenied     = syscall.Errno(5)
	errorSharingViolation = syscall.Errno(32)
	errorLockViolation    = syscall.Errno(33)
)

var (
	kernel32   = syscall.NewLazyDLL("kernel32.dll")
	moveFileEx = kernel32.NewProc("MoveFileExW")
)

func replacePath(source, target string) error {
	sourcePointer, err := syscall.UTF16PtrFromString(source)
	if err != nil {
		return err
	}
	targetPointer, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return err
	}
	var lastErr error
	for attempt := 0; attempt < replaceAttempts; attempt++ {
		ok, _, callErr := moveFileEx.Call(uintptr(unsafe.Pointer(sourcePointer)), uintptr(unsafe.Pointer(targetPointer)), moveFileReplaceExisting|moveFileWriteThrough)
		if ok != 0 {
			return nil
		}
		lastErr = callErr
		if !transientReplaceError(callErr) {
			return callErr
		}
		time.Sleep(time.Duration(attempt+1) * replaceRetryInterval)
	}
	return lastErr
}

func transientReplaceError(err error) bool {
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false
	}
	switch errno {
	case errorAccessDenied, errorSharingViolation, errorLockViolation:
		return true
	default:
		return false
	}
}

func syncDirectory(string) error { return nil }
