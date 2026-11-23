//go:build linux

package evidence

// P85 F08: does this object carry a POSIX ACCESS ACL?
//
// The install-source rule decides "this user cannot write it" from the mode bits
// and the owner/group facts. That verdict silently assumes the mode bits are the
// whole story, which is exactly what a POSIX ACL denies: an ACL can grant write to
// a NAMED user or group that the owning uid/gid/mode triple never mentions. Rather
// than parse ACL semantics here (and drift from the kernel's own evaluation), the
// rule takes the fail-closed branch the review allows: an object that carries an
// access ACL cannot be proven unwritable by a mode-based rule, so it is refused.
//
// "Carries" is asked of the kernel: system.posix_acl_access exists only when an
// access ACL has been set. ENODATA means there is none, and ENOTSUP means the
// filesystem cannot store one, so in both cases no named grant can exist. Any
// other error is a refusal: a question that could not be answered is not a "no".
//
// The handle form is what a bound object uses (P85 F06: the permission facts
// describe the object the digest and identity come from). P85 R7 added the
// DIRECTORY handle form (BoundDir.AccessACL): fgetxattr(2) is defined for any
// descriptor opened readable, an O_RDONLY|O_DIRECTORY descriptor included, so the
// write-domain observation no longer has to read a directory ACL by name. The path
// form remains for the two call sites that hold no descriptor: the nominal
// write-domain entry, which re-lstats the name after the read, and the
// install-source ancestor walk, which lstats each component once before the read
// (see AccessACLPath).

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

// accessACLAttribute is the xattr name Linux uses for the POSIX access ACL.
const accessACLAttribute = "system.posix_acl_access"

// AccessACL reports whether the object this handle addresses carries a POSIX
// access ACL.
func (bound *BoundFile) AccessACL() (bool, error) {
	if err := bound.checkReadable(); err != nil {
		return false, err
	}
	return accessACLOnDescriptor(bound.file.Fd())
}

// AccessACL reports whether the DIRECTORY this descriptor addresses carries a
// POSIX access ACL (P85 R7).
//
// It is the directory twin of BoundFile.AccessACL and exists for the same reason:
// the write-domain observation has to judge the ACL of the object the descriptor
// addresses, not of whatever the name resolves to a moment later. fgetxattr(2) is
// defined for any descriptor opened readable - an O_RDONLY|O_DIRECTORY descriptor
// included - so no path is needed here and none is resolved.
//
// An O_PATH descriptor would answer EBADF to fgetxattr(2); the bindings never
// create one (both OpenBoundDir and the verified binding open with O_RDONLY), so
// EBADF is reported as the refusal it is instead of being folded into a verdict.
func (bound *BoundDir) AccessACL() (bool, error) {
	if bound == nil {
		return false, fmt.Errorf("evidence: nil bound directory")
	}
	bound.mu.Lock()
	defer bound.mu.Unlock()
	if bound.fd < 0 {
		return false, fmt.Errorf("evidence: nil bound directory")
	}
	return accessACLOnDescriptor(uintptr(bound.fd))
}

// accessACLOnDescriptor is the shared fgetxattr(2) probe behind both handle
// forms. go1.22's syscall package exposes Getxattr but not Fgetxattr, so the ABI
// call goes out directly: fgetxattr(fd, name, NULL, 0) reports the size the value
// would need, or the errno that decides the verdict.
func accessACLOnDescriptor(fd uintptr) (bool, error) {
	name, err := syscall.BytePtrFromString(accessACLAttribute)
	if err != nil {
		return false, err
	}
	_, _, errno := syscall.Syscall6(
		syscall.SYS_FGETXATTR,
		fd,
		uintptr(unsafe.Pointer(name)),
		0,
		0,
		0,
		0,
	)
	if errno != 0 {
		return accessACLVerdict(errno)
	}
	return true, nil
}

// AccessACLPath reports whether the object the path currently names carries a
// POSIX access ACL. It is the name form, used only where no descriptor exists:
//
//   - the nominal write-domain entry re-lstats the name AFTER the read and
//     compares the identity against its before observation (P85 R7), so a name
//     that moved is refused; the read itself still resolves the name, and that
//     residue is the registered internal API boundary;
//   - the install-source ancestor walk lstats each component ONCE before the
//     read and refuses a link component, so the identity and mode facts precede
//     the ACL read; no post-read re-lstat is performed on that path.
func AccessACLPath(path string) (bool, error) {
	if _, err := syscall.Getxattr(path, accessACLAttribute, nil); err != nil {
		return accessACLVerdict(err)
	}
	return true, nil
}

// accessACLVerdict maps one errno onto "no ACL can exist" or a refusal.
func accessACLVerdict(err error) (bool, error) {
	switch {
	case errors.Is(err, syscall.ENODATA):
		// No access ACL was ever set on this object.
		return false, nil
	case errors.Is(err, syscall.ENOTSUP):
		// The filesystem cannot carry POSIX ACLs at all (ENOTSUP and EOPNOTSUPP
		// are the same value on Linux, so one case covers both).
		return false, nil
	default:
		return false, fmt.Errorf("evidence: the access ACL cannot be read: %w", err)
	}
}
