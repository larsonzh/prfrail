//go:build linux

package adapters

// F01, Linux half: the managed agent runs under this process's uid/gid, so the
// proof is the F08 verdict of the install-source proof applied to the policy
// directory itself:
//
//   - a directory owned by the executing user is refused, because the owner can
//     restore any permission with chmod, so its mode bits prove nothing;
//   - other-writable is refused directly; group-writable is refused when the
//     executing user is in that group;
//   - a directory that carries a POSIX access ACL is refused, because a named
//     entry can grant this user write access that the mode/owner rule never sees.
//
// The symbolic-link classification comes from os.Lstat: the NAME chain is what the
// kernel walks, so a link component would make the judged object differ from the
// object a write would reach.

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// managedAgentPolicyDomainWriteProofSlots is the Linux half of F01 (design rev 3
// section 1) with the fixed four-slot vector of design rev 4 section 2. Every
// PRESENT slot - the policy directory, its parent and each existing
// artifact/anchor file - is judged with the same three rules the directory alone
// used to get:
//
//   - a symbolic link is refused, because the judged name would not be the object
//     a write reaches (the name form only; the handle form was opened without
//     following one);
//   - a POSIX access ACL is refused, because a named entry can grant this user
//     write access the mode/owner rule never sees;
//   - effective writability (owner, other-write, group-write in group) is refused,
//     because an owner restores any permission with chmod.
//
// A slot that is not present contributes the canonical empty marker, so the vector
// shape is stable. The observation line of each object carries dev:ino, so
// replacing the directory, its parent or a file changes the digest.
func managedAgentPolicyDomainWriteProofSlots(slots []managedAgentPolicyDomainSlot) (string, error) {
	// The superuser writes every path regardless of mode and ownership, so the
	// mode/owner rule can prove nothing for it.
	if os.Geteuid() == 0 {
		return "", fmt.Errorf("%w: the root user can write every path, so the policy domain cannot be proven closed", ErrManagedAgentPolicyWriteProof)
	}
	observed := make([]string, 0, len(slots))
	var failures []error
	for _, slot := range slots {
		if !slot.present {
			observed = append(observed, managedAgentPolicyDomainEmptySlot)
			continue
		}
		line, objectErr := managedAgentPolicyDomainSlotObservation(slot)
		if line == "" {
			// P85 R7: an observation that could not render a line keeps its POSITION in
			// the vector and records the canonical empty marker, so a refusal digest still
			// carries the fixed four-slot shape instead of silently becoming shorter.
			line = managedAgentPolicyDomainEmptySlot
		}
		observed = append(observed, line)
		if objectErr != nil {
			failures = append(failures, objectErr)
		}
	}
	digest := evidence.Digest(managedAgentPolicyWriteProofDomain, []byte(strings.Join(observed, "\n")))
	if len(failures) > 0 {
		return digest, fmt.Errorf("%w: %v", ErrManagedAgentPolicyWriteProof, errors.Join(failures...))
	}
	return digest, nil
}

// managedAgentPolicyDomainSlotObservation judges one present slot and renders its
// observation line. The line is returned even when the verdict refuses, so a
// refusal records what was observed.
//
// P85 design rev 4 section 1: when the slot carries a handle, the permission facts
// come from fstat(2) on THAT descriptor - the same object the reads and writes go
// through.
//
// P85 R6 (review r3 finding): two fstats on the same descriptor only prove the
// DESCRIPTOR did not change; they say nothing about which object the NAME-based ACL
// probe actually read. Every path-based probe therefore runs inside an identity
// sandwich whose two comparisons are name against handle, literally:
//
//	fstat(id0) -> FileIDAt(name) == handle FileID -> ACL probe
//	           -> fstat(id1) -> FileIDAt(name) == handle FileID
//
// so a name that no longer addresses the descriptor's object refuses instead of
// judging another object, and the caller's terminal ReassertIdentities closes the
// observation afterwards. P85 R7 moved the ACL probe itself onto the descriptor for
// both handle forms, so a handle-form slot has no path-based probe left; the two
// comparisons stay because they pin the recorded path to the object that was judged.
//
// P85 R7 (review r6 findings): the access ACL is read from the DESCRIPTOR for both
// handle forms - BoundFile.AccessACL and BoundDir.AccessACL both issue
// fgetxattr(2) on the descriptor the binding opened - so a handle-form slot never
// resolves the name for the ACL probe. The two name comparisons therefore no longer
// bracket a by-name read for a directory slot; they pin the recorded path to the
// descriptor so the observation line (and any refusal that quotes it) describes the
// object that was judged.
//
// The only remaining by-name ACL read is the NOMINAL entry
// (VerifyManagedAgentCannotWritePolicyDomain), which holds no descriptor at all.
// That window is registered as an INTERNAL API BOUNDARY - evidence has no
// directory handle to offer there - not as a kernel limitation: fgetxattr(2) works
// on any readable directory descriptor, and the handle-form callers use it. The
// nominal entry re-lstats the name after the read, so the unverified instant is
// the resolution inside getxattr(2) itself; the probe writes nothing.
func managedAgentPolicyDomainSlotObservation(slot managedAgentPolicyDomainSlot) (string, error) {
	fd, hasHandle, err := managedAgentPolicyDomainSlotDescriptor(slot)
	if err != nil {
		return "", err
	}
	var handleIdentity evidence.FileID
	if hasHandle {
		handleIdentity, err = managedAgentPolicyDomainSlotHandleIdentity(slot)
		if err != nil {
			return "", err
		}
	}
	var stat syscall.Stat_t
	var mode os.FileMode
	if hasHandle {
		if err := syscall.Fstat(fd, &stat); err != nil {
			return "", fmt.Errorf("%w: fstat the %s %q on its handle: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		mode = os.FileMode(stat.Mode)
	} else {
		info, err := os.Lstat(slot.path)
		if err != nil {
			return "", fmt.Errorf("%w: lstat the %s %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("%w: the %s %q is a symbolic link, so the judged name is not the object the writes would reach", ErrManagedAgentPolicyWriteProof, slot.label, slot.path)
		}
		owned, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return "", fmt.Errorf("%w: the %s %q carries no ownership information", ErrManagedAgentPolicyWriteProof, slot.label, slot.path)
		}
		stat = *owned
		mode = info.Mode()
	}
	fileType := uint32(stat.Mode) & uint32(syscall.S_IFMT)
	if slot.isDir && fileType != uint32(syscall.S_IFDIR) {
		return "", fmt.Errorf("%w: the %s %q is not a directory", ErrManagedAgentPolicyWriteProof, slot.label, slot.path)
	}
	if !slot.isDir && fileType != uint32(syscall.S_IFREG) {
		return "", fmt.Errorf("%w: the %s %q is not a regular file", ErrManagedAgentPolicyWriteProof, slot.label, slot.path)
	}
	line := fmt.Sprintf("%s:%d:%d:%04o:%d:%d", slot.path, stat.Dev, stat.Ino, mode.Perm(), stat.Uid, stat.Gid)
	// The opening name comparison: the observation line quotes this path, so the
	// path has to address the descriptor's object on BOTH sides of the probe.
	if err := managedAgentPolicyDomainAssertNameIdentity(slot, handleIdentity, "before"); err != nil {
		return line, err
	}
	aclErr := managedAgentPolicyDomainDenyAccessACL(slot)
	// P85 R7: the closing verification runs even when the probe refused. A refusal
	// has to describe the object the observation was pinned to, so the refused branch
	// re-checks the descriptor (and, on the handle form, the name) before returning
	// the probe error instead of returning straight out of the probe.
	probeErr := managedAgentPolicyDomainAfterProbe(slot, fd, hasHandle, handleIdentity, stat)
	if aclErr != nil {
		return line, aclErr
	}
	if probeErr != nil {
		return line, probeErr
	}
	if err := installSourceEffectiveWritability(slot.path, mode, int(stat.Uid), int(stat.Gid)); err != nil {
		return line, fmt.Errorf("%w: %v", ErrManagedAgentPolicyWriteProof, err)
	}
	return line, nil
}

// managedAgentPolicyDomainAfterProbe re-verifies, AFTER the ACL probe, that the
// observation still describes the object it started on (P85 R7).
//
//   - handle form: fstat(2) must report the same identity, and the NAME must still
//     address the descriptor's object. With the opening comparison this brackets the
//     probe, so the recorded line, the probe and the descriptor are one object.
//   - name form (the nominal entry, which has no descriptor): a second Lstat must
//     reproduce the first observation. The ACL read still resolves the name, so the
//     window shrinks to the resolution inside getxattr(2); that residue is
//     registered as an internal API boundary (see the observation comment).
func managedAgentPolicyDomainAfterProbe(slot managedAgentPolicyDomainSlot, fd int, hasHandle bool, handleIdentity evidence.FileID, before syscall.Stat_t) error {
	if !hasHandle {
		info, err := os.Lstat(slot.path)
		if err != nil {
			return fmt.Errorf("%w: lstat the %s %q after its ACL check: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		after, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("%w: the %s %q carries no ownership information after its ACL check", ErrManagedAgentPolicyWriteProof, slot.label, slot.path)
		}
		if after.Dev != before.Dev || after.Ino != before.Ino || after.Mode != before.Mode || after.Uid != before.Uid || after.Gid != before.Gid {
			return fmt.Errorf("%w: the %s %q changed identity during its ACL check, so the ACL was not read from the observed object", ErrManagedAgentPolicyWriteProof, slot.label, slot.path)
		}
		return nil
	}
	var after syscall.Stat_t
	if err := syscall.Fstat(fd, &after); err != nil {
		return fmt.Errorf("%w: fstat the %s %q after its ACL check: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
	}
	if after.Dev != before.Dev || after.Ino != before.Ino || after.Mode != before.Mode || after.Uid != before.Uid || after.Gid != before.Gid {
		return fmt.Errorf("%w: the %s %q changed identity between its descriptor observation and the ACL check", ErrManagedAgentPolicyWriteProof, slot.label, slot.path)
	}
	// The closing name comparison. With the opening one it proves the recorded path
	// addressed the descriptor on both sides of the probe; the caller's
	// ReassertIdentities is the terminal verdict over the whole domain.
	return managedAgentPolicyDomainAssertNameIdentity(slot, handleIdentity, "after")
}

// managedAgentPolicyDomainSlotHandleIdentity reports the identity of the handle
// that owns the slot. It is the value every name-based comparison has to reproduce
// literally.
func managedAgentPolicyDomainSlotHandleIdentity(slot managedAgentPolicyDomainSlot) (evidence.FileID, error) {
	if slot.dir != nil {
		identity, err := slot.dir.FileID()
		if err != nil {
			return "", fmt.Errorf("%w: identify the %s handle %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		return identity, nil
	}
	if slot.file != nil {
		identity, err := slot.file.FileID()
		if err != nil {
			return "", fmt.Errorf("%w: identify the %s handle %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		return evidence.FileID(identity), nil
	}
	return "", nil
}

// managedAgentPolicyDomainAssertNameIdentity requires the slot's NAME to still
// address the object its handle addresses (P85 R6). evidence.FileIDAt resolves the
// name the way the write would, the handle reports its own identity from fstat(2),
// and the two are compared as strings: an equal pair means the path-based probe and
// the descriptor describe one object at that instant.
func managedAgentPolicyDomainAssertNameIdentity(slot managedAgentPolicyDomainSlot, handleIdentity evidence.FileID, when string) error {
	if handleIdentity == "" {
		return nil
	}
	current, err := evidence.FileIDAt(slot.path)
	if err != nil {
		return fmt.Errorf("%w: the %s %q can no longer be observed %s its ACL check: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, when, err)
	}
	if current != handleIdentity {
		return fmt.Errorf("%w: the %s %q no longer addresses the object its handle addresses %s its ACL check (%s != %s)",
			ErrManagedAgentPolicyWriteProof, slot.label, slot.path, when, current, handleIdentity)
	}
	return nil
}

// managedAgentPolicyDomainDenyAccessACL is the ACL probe of one slot.
//
// Both handle forms read the ACL from their own descriptor through fgetxattr(2),
// so no name is resolved (P85 R7). Only the nominal slot - neither a bound file nor
// a bound directory - falls back to the name-based read, which is the registered
// internal API boundary.
func managedAgentPolicyDomainDenyAccessACL(slot managedAgentPolicyDomainSlot) error {
	if slot.file != nil {
		present, err := slot.file.AccessACL()
		if err != nil {
			return fmt.Errorf("%w: the access ACL of the %s %q cannot be read on its handle: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		if present {
			return fmt.Errorf("%w: the %s %q carries a POSIX access ACL, so a named grant to this user cannot be excluded", ErrManagedAgentPolicyWriteProof, slot.label, slot.path)
		}
		return nil
	}
	if slot.dir != nil {
		present, err := slot.dir.AccessACL()
		if err != nil {
			return fmt.Errorf("%w: the access ACL of the %s %q cannot be read on its handle: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		if present {
			return fmt.Errorf("%w: the %s %q carries a POSIX access ACL, so a named grant to this user cannot be excluded", ErrManagedAgentPolicyWriteProof, slot.label, slot.path)
		}
		return nil
	}
	if err := installSourceDenyAccessACLPath(slot.path); err != nil {
		return fmt.Errorf("%w: %v", ErrManagedAgentPolicyWriteProof, err)
	}
	return nil
}

// managedAgentPolicyDomainSlotDescriptor resolves the descriptor of a handle-form
// slot; hasHandle=false selects the name form.
func managedAgentPolicyDomainSlotDescriptor(slot managedAgentPolicyDomainSlot) (int, bool, error) {
	if slot.dir != nil {
		fd, err := slot.dir.Descriptor()
		if err != nil {
			return -1, false, fmt.Errorf("%w: read the descriptor of the %s %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		return int(fd), true, nil
	}
	if slot.file != nil {
		fd, err := slot.file.Descriptor()
		if err != nil {
			return -1, false, fmt.Errorf("%w: read the descriptor of the %s %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		return int(fd), true, nil
	}
	return -1, false, nil
}
