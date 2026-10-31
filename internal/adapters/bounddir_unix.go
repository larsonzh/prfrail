//go:build !windows

package adapters

// P42-T027B-EVIDENCE-BIND Q2, Unix half: a directory descriptor opened once, and
// relative opens through openat(2) with O_NOFOLLOW and O_CLOEXEC. openat resolves
// relative to the descriptor, so a replaced ancestor cannot redirect the read and
// a leaf swapped for a symlink is refused instead of followed.

import (
	"fmt"
	"os"
	"syscall"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// BoundDir is one open directory descriptor. It owns the descriptor until Close.
type BoundDir struct {
	fd   int
	path string
}

// OpenBoundDir opens path as a directory descriptor. The descriptor keeps
// addressing the directory object even after the name is renamed or replaced, and
// it does not block removal of the entries an owner may still delete.
func OpenBoundDir(path string) (*BoundDir, error) {
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	var stat syscall.Stat_t
	if err := syscall.Fstat(fd, &stat); err != nil {
		_ = syscall.Close(fd)
		return nil, &os.PathError{Op: "stat", Path: path, Err: err}
	}
	if uint32(stat.Mode)&uint32(syscall.S_IFMT) != uint32(syscall.S_IFDIR) {
		_ = syscall.Close(fd)
		return nil, fmt.Errorf("adapters: %s is not a directory", path)
	}
	return &BoundDir{fd: fd, path: path}, nil
}

// OpenBoundDirFile opens one artifact relative to the bound directory. The
// returned handle is the artifact's identity: the caller reads through it and
// never resolves the name again.
func OpenBoundDirFile(dir *BoundDir, rel string, bound int64) (*evidence.BoundFile, error) {
	return openBoundDirFile(dir, rel, bound, false)
}

// openBoundDirFile is the class-aware form: allowTruncation keeps the log and
// event class readable past its bound instead of refusing the open.
func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*evidence.BoundFile, error) {
	if dir == nil || dir.fd < 0 {
		return nil, fmt.Errorf("adapters: nil bound directory")
	}
	clean, err := cleanBoundRelativePath(rel)
	if err != nil {
		return nil, err
	}
	fd, err := syscall.Openat(dir.fd, clean, syscall.O_RDONLY|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, &os.PathError{Op: "openat", Path: clean, Err: err}
	}
	// P42 F1: openat resolves the artifact against the descriptor in the same
	// syscall that opens it, so there is no name-then-open window left to re-check
	// afterwards - unlike the Windows half, which has to confirm the opened file
	// handle against the bound directory's current real name.
	return evidence.NewBoundFile(os.NewFile(uintptr(fd), clean), bound, allowTruncation)
}

// Close releases the directory descriptor.
func (dir *BoundDir) Close() error {
	if dir == nil || dir.fd < 0 {
		return nil
	}
	err := syscall.Close(dir.fd)
	dir.fd = -1
	return err
}

// boundDirVerify proves the handle still addresses the directory the caller
// named, by comparing the descriptor's device and inode with the one the expected
// path currently resolves to. A directory swapped before the bind is refused; a
// path whose ancestors are spelled differently still resolves to the same object
// and is accepted.
func boundDirVerify(dir *BoundDir, expected string) error {
	if dir == nil || dir.fd < 0 {
		return fmt.Errorf("adapters: nil bound directory")
	}
	var boundStat syscall.Stat_t
	if err := syscall.Fstat(dir.fd, &boundStat); err != nil {
		return fmt.Errorf("adapters: stat the bound directory: %v", err)
	}
	info, err := os.Stat(expected)
	if err != nil {
		return err
	}
	expectedStat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("adapters: %s carries no syscall stat", expected)
	}
	if expectedStat.Dev != boundStat.Dev || expectedStat.Ino != boundStat.Ino {
		return fmt.Errorf("adapters: bound directory %s is not %s", dir.path, expected)
	}
	return nil
}
