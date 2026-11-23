//go:build linux

package adapters

// P85 R6 / the Linux native leg of the identity sandwich (review r3 finding 1).
//
// managedAgentPolicyDomainSlotObservation runs, per present slot,
//
//	fstat(handle) -> name == handle identity -> ACL probe
//	              -> fstat(handle) again      -> name == handle identity
//
// The leg below states the bracket at its smallest discriminating unit: the opening
// comparison passes on the untouched name, the swap happens, and the closing
// comparison refuses. It then runs the production observation over a name that
// already moved, so the refusal is attributable to the sandwich and not to the
// ACL or the writability rules (both of which this host would also refuse, which is
// why the assertion is on the message).
//
// Verification status: this file is skipped on Windows; its authoritative run is
// the Ubuntu VM leg (本机 skip 待 VM).

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func init() {
	p85R6SandwichPlatformLegs = append(p85R6SandwichPlatformLegs, p85R6SandwichLeg{
		name: "the-linux-native-sandwich-refuses-a-swap-between-its-two-comparisons",
		run: func(t *testing.T) {
			slot, handleIdentity, swap := p85R6SandwichFixture(t)
			if err := managedAgentPolicyDomainAssertNameIdentity(slot, handleIdentity, "before"); err != nil {
				t.Fatalf("the untouched name must pass the opening comparison, got %v", err)
			}
			swap()
			err := managedAgentPolicyDomainAssertNameIdentity(slot, handleIdentity, "after")
			if err == nil {
				t.Fatal("a name swapped between the two comparisons must refuse the closing comparison")
			}
			if !errors.Is(err, ErrManagedAgentPolicyWriteProof) {
				t.Fatalf("the refusal must be the write-proof sentinel, got %v", err)
			}
			if !strings.Contains(err.Error(), "no longer addresses the object its handle addresses") {
				t.Fatalf("the refusal must be the name-versus-handle comparison, got %v", err)
			}
			if !strings.Contains(err.Error(), "after") {
				t.Fatalf("the refusal must say which side of the sandwich failed, got %v", err)
			}

			// The production observation over the same moved name refuses as well, so
			// the bracket above is not a parallel re-implementation.
			movedSlot, _, moveAgain := p85R6SandwichFixture(t)
			moveAgain()
			if _, err := managedAgentPolicyDomainSlotObservation(movedSlot); err == nil {
				t.Fatal("the production observation must refuse a name that no longer addresses its handle")
			} else if !strings.Contains(err.Error(), "before its ACL check") {
				t.Fatalf("the production refusal must come from the opening comparison of the sandwich, got %v", err)
			}
		},
	})
}

// p85R7SetAccessACL writes an EXTENDED POSIX access ACL (version 2: user object
// rwx, a named user rwx, group object r-x, mask rwx, other r-x) onto path, or
// skips the caller when this host's filesystem cannot store one. The named user
// entry matters: a MINIMAL ACL (no named entries) is equivalent to the mode bits,
// and the kernel does not persist it, so getxattr/fgetxattr would answer ENODATA
// and the fd-versus-name discrimination below could never be established. Without
// a storable ACL the discrimination cannot run, so skipping is the honest outcome.
func p85R7SetAccessACL(t *testing.T, path string) {
	t.Helper()
	blob := []byte{2, 0, 0, 0} // ACL version 2, little-endian
	appendEntry := func(tag, perm uint16, id uint32) {
		blob = append(blob,
			byte(tag), byte(tag>>8),
			byte(perm), byte(perm>>8),
			byte(id), byte(id>>8), byte(id>>16), byte(id>>24))
	}
	appendEntry(0x01, 0o7, 0xFFFFFFFF) // ACL_USER_OBJ
	appendEntry(0x02, 0o7, 4242)       // ACL_USER (named; forces kernel storage)
	appendEntry(0x04, 0o5, 0xFFFFFFFF) // ACL_GROUP_OBJ
	appendEntry(0x10, 0o7, 0xFFFFFFFF) // ACL_MASK (required with named entries)
	appendEntry(0x20, 0o5, 0xFFFFFFFF) // ACL_OTHER
	if err := syscall.Setxattr(path, "system.posix_acl_access", blob, 0); err != nil {
		t.Skipf("this host cannot store a POSIX access ACL on %s: %v", path, err)
	}
}

// p85R7DirectorySlotFixture binds one policy directory and returns the DIRECTORY
// slot of the fixed four-slot vector, the identity its handle reports, and the swap
// that moves the name onto a fresh directory. The descriptor keeps addressing the
// object that was bound, which is what the slot's identity has to follow.
func p85R7DirectorySlotFixture(t *testing.T) (managedAgentPolicyDomainSlot, evidence.FileID, func()) {
	t.Helper()
	root := t.TempDir()
	policy := filepath.Join(root, "agent-runner-policy")
	if err := os.MkdirAll(policy, 0o755); err != nil {
		t.Fatalf("create the policy directory: %v", err)
	}
	directory, identity, err := evidence.OpenBoundDirVerified(policy)
	if err != nil {
		t.Skipf("this host's temporary root is not bindable as a verified directory: %v", err)
	}
	t.Cleanup(func() { _ = directory.Close() })
	slot := managedAgentPolicyDomainSlot{
		label:   "policy directory",
		path:    policy,
		isDir:   true,
		present: true,
		dir:     directory,
	}
	swap := func() {
		t.Helper()
		// Rename first so the retired directory keeps its identity; a create-after-delete
		// can hand the same inode back and make the fixture agree by accident.
		if err := os.Rename(policy, policy+".retired"); err != nil {
			t.Fatalf("rename the policy directory away: %v", err)
		}
		if err := os.Mkdir(policy, 0o755); err != nil {
			t.Fatalf("create the replacement policy directory: %v", err)
		}
	}
	return slot, identity, swap
}

// TestP85R7LinuxDirectorySlotReadsTheACLFromItsDescriptor is the discriminating
// case for the fd-ification of the directory slot: the descriptor addresses a
// directory that CARRIES an access ACL while the NAME addresses one that does not.
// Only a descriptor-based read refuses; a regression to evidence.AccessACLPath
// accepts, and the assertion below goes red.
func TestP85R7LinuxDirectorySlotReadsTheACLFromItsDescriptor(t *testing.T) {
	slot, _, swap := p85R7DirectorySlotFixture(t)
	p85R7SetAccessACL(t, slot.path)
	swap()

	byName, err := evidence.AccessACLPath(slot.path)
	if err != nil {
		t.Skipf("this host cannot read the access ACL of %s by name: %v", slot.path, err)
	}
	if byName {
		t.Skipf("the replacement directory inherited an access ACL, so the fd-versus-name discrimination cannot be established here")
	}
	onHandle, err := slot.dir.AccessACL()
	if err != nil {
		t.Fatalf("read the directory access ACL on its descriptor: %v", err)
	}
	if !onHandle {
		t.Fatal("the descriptor addresses the ACL-carrying directory; a name-based read would report the ACL-free stand-in")
	}
	if err := managedAgentPolicyDomainDenyAccessACL(slot); err == nil {
		t.Fatal("the production ACL probe must refuse the directory that carries an access ACL")
	} else if !strings.Contains(err.Error(), "carries a POSIX access ACL") {
		t.Fatalf("the refusal must be the access-ACL verdict, got %v", err)
	}
}

// TestP85R7LinuxDirectorySlotRefusesAReplacedName pins both identity comparisons of
// the DIRECTORY slot. The opening comparison passes on the untouched name, the swap
// happens, the closing comparison refuses, and the production observation refuses
// with the identity verdict rather than with the writability rule (the temporary
// root is owned by the executing user, which that rule would also refuse, so the
// message is what attributes the refusal to the sandwich).
//
// Mutation check: short-circuiting the directory slot's name comparison (or
// removing the closing call) leaves the observation to refuse with a different
// verdict, and this leg goes red.
func TestP85R7LinuxDirectorySlotRefusesAReplacedName(t *testing.T) {
	slot, identity, swap := p85R7DirectorySlotFixture(t)
	if err := managedAgentPolicyDomainAssertNameIdentity(slot, identity, "before"); err != nil {
		t.Fatalf("the untouched name must pass the opening comparison, got %v", err)
	}
	swap()
	closingErr := managedAgentPolicyDomainAssertNameIdentity(slot, identity, "after")
	if closingErr == nil {
		t.Fatal("a moved name must refuse the closing comparison")
	}
	if !errors.Is(closingErr, ErrManagedAgentPolicyWriteProof) {
		t.Fatalf("the closing refusal must be the write-proof sentinel, got %v", closingErr)
	}
	if !strings.Contains(closingErr.Error(), "no longer addresses the object its handle addresses") {
		t.Fatalf("the closing refusal must be the name-versus-handle comparison, got %v", closingErr)
	}

	_, err := managedAgentPolicyDomainSlotObservation(slot)
	if err == nil {
		t.Fatal("the production observation must refuse a directory whose name moved off its handle")
	}
	if !strings.Contains(err.Error(), "before its ACL check") {
		t.Fatalf("the refusal must come from the opening comparison of the sandwich, got %v", err)
	}
	if strings.Contains(err.Error(), "owned by the executing user") {
		t.Fatal("the refusal must be the identity comparison, not the writability rule")
	}
}

// TestP85R7LinuxDirectorySlotClosingComparisonRefusesAMovedName exercises the
// CLOSING half of the directory sandwich directly. The post-probe verification runs
// on the descriptor and then repeats the name comparison; with the name moved off
// the descriptor that comparison must refuse.
//
// Mutation check: deleting the closing comparison from
// managedAgentPolicyDomainAfterProbe leaves this call returning nil and the leg goes
// red (the descriptor itself did not change, so nothing else can refuse here).
func TestP85R7LinuxDirectorySlotClosingComparisonRefusesAMovedName(t *testing.T) {
	slot, identity, swap := p85R7DirectorySlotFixture(t)
	fd, hasHandle, err := managedAgentPolicyDomainSlotDescriptor(slot)
	if err != nil {
		t.Fatalf("read the directory descriptor: %v", err)
	}
	if !hasHandle {
		t.Fatal("the directory slot must carry a descriptor")
	}
	var before syscall.Stat_t
	if err := syscall.Fstat(fd, &before); err != nil {
		t.Fatalf("fstat the directory descriptor: %v", err)
	}
	swap()
	closingErr := managedAgentPolicyDomainAfterProbe(slot, fd, hasHandle, identity, before)
	if closingErr == nil {
		t.Fatal("the closing half must refuse a name that moved off the descriptor")
	}
	if !strings.Contains(closingErr.Error(), "no longer addresses the object its handle addresses") {
		t.Fatalf("the refusal must be the closing name-versus-handle comparison, got %v", closingErr)
	}
}
