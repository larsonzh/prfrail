//go:build linux

package evidence

// P85 F20: handle-relative creation and publication inside a bound directory.
//
// The carrier write used to bind and re-confirm the directory, then create the
// temporary file and publish the final name BY PATH. Between the confirmation and
// the two path operations the name could be pointed somewhere else, so "the
// directory that was proven" and "the directory the bytes landed in" were two
// claims, not one object. These three methods remove the name from the operation:
// the creation, the hard link and the cleanup are all relative to the descriptor
// the caller bound, so the object that was proven is the object that is written.
//
// Everything is deliberately limited to ONE path component. A multi-component
// name would reintroduce a resolution walk inside the operation, and the whole
// point is that no resolution happens at all: the kernel looks up the single
// entry in the directory this descriptor addresses.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

// CreateFileExclusive creates one new file named rel inside the bound directory
// and returns its handle. An existing entry - including a symbolic link, because
// O_EXCL is evaluated before any link is followed - is refused with fs.ErrExist.
//
// perm is the creation mode. The caller writes, syncs and closes the handle, and
// publishes it with LinkNoReplace; the temporary name is never the published one.
func (dir *BoundDir) CreateFileExclusive(rel string, perm uint32) (*os.File, error) {
	if dir == nil {
		return nil, fmt.Errorf("evidence: nil bound directory")
	}
	clean, err := cleanBoundEntryName(rel)
	if err != nil {
		return nil, err
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return nil, fmt.Errorf("evidence: nil bound directory")
	}
	fd, err := syscall.Openat(dir.fd, clean, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_EXCL|syscall.O_CLOEXEC|syscall.O_NOFOLLOW, perm)
	if err != nil {
		return nil, &os.PathError{Op: "openat", Path: clean, Err: err}
	}
	return os.NewFile(uintptr(fd), clean), nil
}

// LinkNoReplace publishes newRel as a hard link to the already created entry
// existingRel, both inside the bound directory. An existing newRel is refused with
// fs.ErrExist: a published proof is never replaced, and there is no rename
// fallback that would replace it.
func (dir *BoundDir) LinkNoReplace(existingRel, newRel string) error {
	if dir == nil {
		return fmt.Errorf("evidence: nil bound directory")
	}
	existing, err := cleanBoundEntryName(existingRel)
	if err != nil {
		return err
	}
	published, err := cleanBoundEntryName(newRel)
	if err != nil {
		return err
	}
	existingPtr, err := syscall.BytePtrFromString(existing)
	if err != nil {
		return err
	}
	publishedPtr, err := syscall.BytePtrFromString(published)
	if err != nil {
		return err
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return fmt.Errorf("evidence: nil bound directory")
	}
	// go1.22's syscall package does not export Linkat, so the ABI call goes out
	// directly, exactly like openat2 above: no cgo, no new dependency.
	_, _, errno := syscall.Syscall6(
		syscall.SYS_LINKAT,
		uintptr(dir.fd),
		uintptr(unsafe.Pointer(existingPtr)),
		uintptr(dir.fd),
		uintptr(unsafe.Pointer(publishedPtr)),
		0,
		0,
	)
	if errno != 0 {
		return &os.LinkError{Op: "linkat", Old: existing, New: published, Err: errno}
	}
	return nil
}

// Remove deletes one regular entry, named rel, from the bound directory. It is the
// cleanup half of the temporary name; a failure is reported so the caller can
// decide, and it never removes a directory (syscall.Unlinkat without
// AT_REMOVEDIR).
func (dir *BoundDir) Remove(rel string) error {
	if dir == nil {
		return fmt.Errorf("evidence: nil bound directory")
	}
	clean, err := cleanBoundEntryName(rel)
	if err != nil {
		return err
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return fmt.Errorf("evidence: nil bound directory")
	}
	if err := syscall.Unlinkat(dir.fd, clean); err != nil {
		return &os.PathError{Op: "unlinkat", Path: clean, Err: err}
	}
	return nil
}

// cleanBoundEntryName applies the shared artifact-name guard and additionally
// requires a SINGLE directory entry: no separator may survive, because a nested
// name would be resolved component by component and could cross a symbolic link
// the bound handle does not cover.
func cleanBoundEntryName(rel string) (string, error) {
	clean, err := cleanBoundRelativePath(rel)
	if err != nil {
		return "", err
	}
	if strings.Contains(clean, string(filepath.Separator)) {
		return "", fmt.Errorf("evidence: %q is not a single directory entry", rel)
	}
	return clean, nil
}
