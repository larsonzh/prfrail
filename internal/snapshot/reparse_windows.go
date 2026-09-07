//go:build windows

package snapshot

import (
	"os"
	"syscall"
)

const fileAttributeReparsePoint = 0x00000400

func isReparsePointOrJunction(info os.FileInfo, _ string) bool {
	sys, ok := info.Sys().(*syscall.Win32FileAttributeData)
	if !ok {
		return false
	}
	if sys.FileAttributes&fileAttributeReparsePoint == 0 {
		return false
	}
	// On Windows, if it's a directory with reparse point, it is a junction or mount point
	if info.IsDir() {
		return true
	}
	// If it is a reparse point but not a standard symlink, it is an irregular reparse point (e.g. OneDrive placeholder, HSM)
	if info.Mode()&os.ModeSymlink == 0 {
		return true
	}
	return false
}
