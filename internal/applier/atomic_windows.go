//go:build windows

package applier

import (
	"syscall"
	"unsafe"
)

const (
	moveFileReplaceExisting = 0x1
	moveFileWriteThrough    = 0x8
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
	ok, _, callErr := moveFileEx.Call(uintptr(unsafe.Pointer(sourcePointer)), uintptr(unsafe.Pointer(targetPointer)), moveFileReplaceExisting|moveFileWriteThrough)
	if ok == 0 {
		return callErr
	}
	return nil
}

func syncDirectory(string) error { return nil }
