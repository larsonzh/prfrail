//go:build windows

package adapters

// T3 (finding F6), Windows half: the install source is immutable when this user's
// token cannot obtain write, append, delete, delete-child or permission-changing
// access to the executable, to its installation directory, AND to every ancestor
// directory up to the volume root, and when no component of the name chain is a
// reparse point.
//
// The verdict is taken from the OS itself: one CreateFile per RIGHT PER OBJECT asks
// for exactly that single mutation right, and a refusal (ERROR_ACCESS_DENIED) is the
// proof we want, while a granted handle - closed immediately, nothing is written -
// means the source is mutable. The DACL is additionally fingerprinted through
// GetNamedSecurityInfoW so the carrier records WHICH security descriptor was
// observed, and a later change of that descriptor is visible as a changed fact.
//
// P85 F03: the earlier form asked for all mutation rights in ONE access check, and a
// single ERROR_ACCESS_DENIED only proved "not all rights at once". An ACL that grants
// just WRITE_OWNER (or just DELETE) was therefore reported as immutable. The rights
// are now enumerated one at a time (see installSourceFileMutationRights /
// installSourceDirectoryMutationRights) so any single granted right refuses.
//
// P85 F04: the observed chain is the FULL ancestor chain, not just the direct parent.
// The volume root is part of the chain but is anchored by IDENTITY (the P84 proof
// that the bound handle is the root of a fixed local NTFS volume) instead of by
// mutation rights - excluding it would miss "a token that may delete or re-permission
// the drive root replaces the first-level ancestor".
//
// The verdict deliberately does NOT come from GetEffectiveRightsFromAclW: that
// function does not resolve CREATOR OWNER correctly and reports GENERIC_ALL for an
// inherited CREATOR OWNER ace, which nearly every Windows directory carries. It
// would therefore refuse every legitimate install under %ProgramFiles% as
// "writable". Asking the same access check the OS will run for the actual token is
// both stronger and cheaper than reimplementing ACE semantics here. Every outcome
// that is neither "granted" nor "ACCESS_DENIED" is a refusal (fail-closed), never
// "assume immutable".
//
// Every errno is read directly from the Go syscall result (errors.Is on
// syscall.ERROR_ACCESS_DENIED). No site re-derives a Windows error from a
// translated/last-error value across a language or library boundary: a translated
// errno is exactly how a "denied" can silently arrive as "granted".
//
// No third-party dependency: GetNamedSecurityInfoW, GetSecurityDescriptorLength
// and LocalFree are declared here; the token and the access check come from the
// standard syscall package.

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	// aiInstallSourceACLDomain is the digest domain of the Windows ACL evidence
	// recorded in the carrier. Declared here, not in the shared file, because only
	// the Windows half produces and consumes it.
	aiInstallSourceACLDomain = "proofrail:ai-install-source-acl:1\n"

	// installSourceSEFileObject is SE_FILE_OBJECT: the security descriptor of a
	// file or directory.
	installSourceSEFileObject = 1
	// installSourceDACLSecurityInformation asks for the DACL only; the descriptor
	// fingerprint therefore covers the ACEs and not the owner or group.
	installSourceDACLSecurityInformation = 0x00000004
	// installSourceFileFlagOpenReparsePoint opens a reparse point AS ITSELF instead
	// of following it, so the name chain can be classified (P85 F04).
	installSourceFileFlagOpenReparsePoint = 0x00200000
	// installSourceFileReadAttributes is FILE_READ_ATTRIBUTES, the least access that
	// still answers "what kind of object is this name".
	installSourceFileReadAttributes = 0x00000080
)

// P85 F03: the mutation rights are enumerated ONE BY ONE. A merged mask asked the
// kernel "may I obtain all of these at once", and a denial of that single question
// proved nothing about any individual right.
//
// File table (unchanged by the F04 ruling): writing content, appending, writing
// extended attributes, writing attributes, deleting the object, changing the DACL,
// changing the owner.
var installSourceFileMutationRights = []uint32{
	installSourceFileWriteData,       // 0x00000002 FILE_WRITE_DATA
	installSourceFileAppendData,      // 0x00000004 FILE_APPEND_DATA
	installSourceFileWriteEA,         // 0x00000010 FILE_WRITE_EA
	installSourceFileWriteAttributes, // 0x00000100 FILE_WRITE_ATTRIBUTES
	installSourceStandardDelete,      // 0x00010000 DELETE
	installSourceStandardWriteDAC,    // 0x00040000 WRITE_DAC
	installSourceStandardWriteOwner,  // 0x00080000 WRITE_OWNER
}

// Directory table: the full seven rights of the rev2 design (FILE_ADD_FILE,
// FILE_ADD_SUBDIRECTORY, FILE_DELETE_CHILD, FILE_WRITE_ATTRIBUTES, DELETE,
// WRITE_DAC, WRITE_OWNER). A non-root ancestor directory that lets the user
// create a child CAN be swapped: create a sibling, delete the entry
// (FILE_DELETE_CHILD) and rename the sibling into place. So creation rights
// stay in the table for every directory EXCEPT the true volume root, which
// the operating system reserves and cannot be renamed, swapped or deleted.
var installSourceDirectoryMutationRights = []uint32{
	installSourceFileWriteData,       // 0x00000002 FILE_ADD_FILE (directory spelling)
	installSourceFileAppendData,      // 0x00000004 FILE_ADD_SUBDIRECTORY
	installSourceFileDeleteChild,     // 0x00000040 FILE_DELETE_CHILD
	installSourceFileWriteAttributes, // 0x00000100 FILE_WRITE_ATTRIBUTES
	installSourceStandardDelete,      // 0x00010000 DELETE
	installSourceStandardWriteDAC,    // 0x00040000 WRITE_DAC
	installSourceStandardWriteOwner,  // 0x00080000 WRITE_OWNER
}

// Volume-root table: the directory table minus the two CREATION rights. The
// operator's F04 ruling (2026-11-23, advisor-triggered, ① single-round) removes
// FILE_ADD_FILE and FILE_ADD_SUBDIRECTORY ONLY for a volume root whose identity
// is anchored by the P84 primitive (FileIndexLow==5). Measured on this host with
// a NON-ELEVATED medium token: C:\ grants FILE_ADD_SUBDIRECTORY (the default
// "Authenticated Users:(AD)" ace; "AD" is append-data/add-subdirectory, no
// inheritance) while FILE_ADD_FILE and GENERIC_WRITE are ACCESS_DENIED. A
// creation right on the reserved root says nothing about whether an EXISTING
// object on the chain can be replaced, rewritten or deleted. When the identity
// anchor FAILS, the root is judged by the full directory table instead - the
// exclusion is a privilege, never an assumption.
var installSourceVolumeRootMutationRights = []uint32{
	installSourceFileDeleteChild,     // 0x00000040 FILE_DELETE_CHILD
	installSourceFileWriteAttributes, // 0x00000100 FILE_WRITE_ATTRIBUTES
	installSourceStandardDelete,      // 0x00010000 DELETE
	installSourceStandardWriteDAC,    // 0x00040000 WRITE_DAC
	installSourceStandardWriteOwner,  // 0x00080000 WRITE_OWNER
}

const (
	installSourceFileWriteData       = 0x00000002
	installSourceFileAppendData      = 0x00000004
	installSourceFileWriteEA         = 0x00000010
	installSourceFileDeleteChild     = 0x00000040
	installSourceFileWriteAttributes = 0x00000100
	installSourceStandardDelete      = 0x00010000
	installSourceStandardWriteDAC    = 0x00040000
	installSourceStandardWriteOwner  = 0x00080000
)

var (
	installSourceKernel32 = syscall.NewLazyDLL("kernel32.dll")
	installSourceAdvapi32 = syscall.NewLazyDLL("advapi32.dll")

	procInstallSourceLocalFree                = installSourceKernel32.NewProc("LocalFree")
	procInstallSourceGetNamedSecurityInfoW    = installSourceAdvapi32.NewProc("GetNamedSecurityInfoW")
	procInstallSourceGetSecurityDescriptorLen = installSourceAdvapi32.NewProc("GetSecurityDescriptorLength")
)

// verifyInstallSourceImmutability observes the whole ancestor chain of the
// executable for this user's token: every directory above the leaf is judged by its
// per-right mutation grants, the volume root is anchored by identity, and the leaf
// is judged by the file table. leaf is the handle the executable was opened through
// (P85 F06); the Windows rule needs no fact from its bytes.
//
// P85 R7 (review r6 finding): the permission facts are tied to THAT object instead
// of to whatever the name addresses. The identity of the leaf handle is captured
// and the leaf NAME has to reproduce it before the probes and again after them, so
// the ACL/DACL observation cannot describe a different object than the bytes that
// were hashed.
func verifyInstallSourceImmutability(executable string, leaf *evidence.BoundFile) (AIInstallSourceImmutability, error) {
	// P85 F07: the ORIGINAL spelling is judged. A path that is not already
	// canonical names a chain the kernel walks literally and this rule would
	// otherwise fold away ("\..", ".", doubled or trailing separators), so it is
	// refused before any component is inspected.
	path := executable
	if err := installSourceRequireCanonicalSpelling(path); err != nil {
		return AIInstallSourceImmutability{}, err
	}
	if leaf == nil {
		return AIInstallSourceImmutability{}, fmt.Errorf("%w: %q has no bound handle, so its permission facts cannot be tied to the object whose bytes were hashed", ErrAIInstallSourceProof, path)
	}
	leafIdentity, err := leaf.FileID()
	if err != nil {
		return AIInstallSourceImmutability{}, fmt.Errorf("%w: identify the bound executable %q: %v", ErrAIInstallSourceProof, path, err)
	}
	leafNameBefore, err := evidence.FileIDAt(path)
	if err != nil {
		return AIInstallSourceImmutability{}, fmt.Errorf("%w: identify the install source %q before its permission probes: %v", ErrAIInstallSourceProof, path, err)
	}
	if string(leafNameBefore) != leafIdentity {
		return AIInstallSourceImmutability{}, fmt.Errorf("%w: the install source %q does not name the executable whose bytes were hashed (%s != %s)",
			ErrAIInstallSourceProof, path, leafNameBefore, leafIdentity)
	}
	sidText, err := installSourceCurrentUserSIDText()
	if err != nil {
		return AIInstallSourceImmutability{}, err
	}
	components := installSourceWindowsComponents(path)
	objectStates := make([]string, 0, len(components))
	denies := true
	for index, component := range components {
		// P85 F04: a reparse point anywhere on the chain - junction above the install
		// directory included - means the name the kernel walks is not the name that was
		// judged. Refuse instead of following.
		reparse, err := installSourceComponentIsReparse(component)
		if err != nil {
			return AIInstallSourceImmutability{}, err
		}
		if reparse {
			return AIInstallSourceImmutability{}, fmt.Errorf("%w: %q on the install source chain is a reparse point",
				ErrAIInstallSourceProof, component)
		}
		if index == 0 {
			// Volume root: identity-anchored by the P84 primitive. The operator's
			// F04 ruling removes the two creation rights ONLY for an anchored true
			// volume root; an anchor failure falls back to the full directory
			// table (fail-closed, the exclusion is never assumed).
			anchored, err := installSourceAnchorsFixedVolumeRoot(component)
			if err != nil {
				return AIInstallSourceImmutability{}, err
			}
			rights := installSourceDirectoryMutationRights
			anchorTag := "anchor:unproven"
			if anchored {
				rights = installSourceVolumeRootMutationRights
				anchorTag = "anchor:fixed-volume"
			}
			granted, denied, err := installSourceMutationAccess(component, true, rights)
			if err != nil {
				return AIInstallSourceImmutability{}, err
			}
			descriptor, err := installSourceDescriptorDigest(component)
			if err != nil {
				return AIInstallSourceImmutability{}, err
			}
			objectStates = append(objectStates, fmt.Sprintf("%s|%s|denied:%s|dacl:%s",
				component, anchorTag, installSourceRightsText(denied), descriptor))
			if len(granted) > 0 {
				denies = false
				break
			}
			continue
		}
		isLeaf := index == len(components)-1
		rights := installSourceDirectoryMutationRights
		if isLeaf {
			rights = installSourceFileMutationRights
		}
		granted, denied, err := installSourceMutationAccess(component, !isLeaf, rights)
		if err != nil {
			return AIInstallSourceImmutability{}, err
		}
		descriptor, err := installSourceDescriptorDigest(component)
		if err != nil {
			return AIInstallSourceImmutability{}, err
		}
		objectStates = append(objectStates, fmt.Sprintf("%s|denied:%s|dacl:%s",
			component, installSourceRightsText(denied), descriptor))
		if len(granted) > 0 {
			denies = false
			break
		}
	}
	// The closing half of the leaf verification: after every by-name probe the leaf
	// name must still address the object the handle addresses. A name that moved
	// during the walk refuses instead of contributing another object's ACL/DACL to
	// the evidence (P85 R7).
	leafNameAfter, err := evidence.FileIDAt(path)
	if err != nil {
		return AIInstallSourceImmutability{}, fmt.Errorf("%w: re-identify the install source %q after its permission probes: %v", ErrAIInstallSourceProof, path, err)
	}
	if string(leafNameAfter) != leafIdentity {
		return AIInstallSourceImmutability{}, fmt.Errorf("%w: the install source %q no longer addresses the object whose bytes were hashed (%s != %s)",
			ErrAIInstallSourceProof, path, leafNameAfter, leafIdentity)
	}
	facts := AIInstallSourceImmutability{Windows: &AIInstallSourceWindowsFacts{
		DirACLDeniesWriteDelete: denies,
		ACLEvidenceHash: evidence.Digest(aiInstallSourceACLDomain, []byte(strings.ToLower(path)+"\n"+
			"objects="+strings.Join(objectStates, ";")+"\n"+
			"trustee="+sidText)),
	}}
	if !denies {
		return facts, fmt.Errorf("%w: the install source %q grants this user a mutation right somewhere on its chain (see the recorded object states)",
			ErrAIInstallSourceProof, path)
	}
	return facts, nil
}

// installSourceWindowsComponents splits one absolute Windows path into its chain
// from the volume root down to the leaf: `C:\` first, the executable last.
func installSourceWindowsComponents(path string) []string {
	volume := filepath.VolumeName(path)
	root := volume + string(filepath.Separator)
	components := []string{root}
	current := root
	remainder := strings.TrimPrefix(path, root)
	for _, part := range strings.Split(remainder, string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		components = append(components, current)
	}
	return components
}

// installSourceAnchorsFixedVolumeRoot proves, on a handle this call opens, that the
// volume root name addresses the root directory of a fixed local NTFS volume. It
// reuses the P84 handle-based anchor (evidence.BoundDir.IsFixedVolumeRoot) instead of
// inventing a path-string rule.
func installSourceAnchorsFixedVolumeRoot(root string) (bool, error) {
	bound, err := evidence.OpenBoundDir(root)
	if err != nil {
		return false, fmt.Errorf("%w: bind the volume root %q: %v", ErrAIInstallSourceProof, root, err)
	}
	defer func() { _ = bound.Close() }()
	return bound.IsFixedVolumeRoot(), nil
}

// installSourceComponentIsReparse classifies one name on the chain WITHOUT following
// it: the object is opened as itself (FILE_FLAG_OPEN_REPARSE_POINT) and its
// attributes are read from the handle. A directory, a file and a reparse point are
// all openable this way; an object that cannot be opened at all is a refusal.
func installSourceComponentIsReparse(name string) (bool, error) {
	namePointer, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return false, fmt.Errorf("%w: encode the path %q: %v", ErrAIInstallSourceProof, name, err)
	}
	share := uint32(syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE | syscall.FILE_SHARE_DELETE)
	attributes := uint32(syscall.FILE_FLAG_BACKUP_SEMANTICS | installSourceFileFlagOpenReparsePoint)
	handle, err := syscall.CreateFile(namePointer, uint32(installSourceFileReadAttributes), share, nil, syscall.OPEN_EXISTING, attributes, 0)
	if err != nil {
		return false, fmt.Errorf("%w: classify %q: %v", ErrAIInstallSourceProof, name, err)
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(handle, &info); err != nil {
		return false, fmt.Errorf("%w: read the attributes of %q: %v", ErrAIInstallSourceProof, name, err)
	}
	return info.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0, nil
}

// installSourceMutationAccess probes each right in the table with its OWN access
// check and returns the granted and the denied rights, both in table order. A
// granted handle is closed immediately and nothing is written, so the probe leaves
// no footprint. Only ERROR_ACCESS_DENIED counts as a denial; every other failure is
// "cannot tell" and therefore a refusal, never "assume denied".
func installSourceMutationAccess(name string, directory bool, rights []uint32) ([]uint32, []uint32, error) {
	namePointer, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: encode the path %q: %v", ErrAIInstallSourceProof, name, err)
	}
	share := uint32(syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE | syscall.FILE_SHARE_DELETE)
	attributes := uint32(syscall.FILE_ATTRIBUTE_NORMAL)
	if directory {
		// A directory can only be opened with FILE_FLAG_BACKUP_SEMANTICS.
		attributes = syscall.FILE_FLAG_BACKUP_SEMANTICS
	}
	granted := make([]uint32, 0, len(rights))
	denied := make([]uint32, 0, len(rights))
	for _, right := range rights {
		handle, err := syscall.CreateFile(namePointer, right, share, nil, syscall.OPEN_EXISTING, attributes, 0)
		if err == nil {
			_ = syscall.CloseHandle(handle)
			granted = append(granted, right)
			continue
		}
		if errors.Is(err, syscall.ERROR_ACCESS_DENIED) {
			denied = append(denied, right)
			continue
		}
		return granted, denied, fmt.Errorf("%w: determine right 0x%08x of %q: %v", ErrAIInstallSourceProof, right, name, err)
	}
	return granted, denied, nil
}

// installSourceRightsText renders one rights vector in table order, so the carrier
// evidence distinguishes WHICH rights were denied and not merely how many.
func installSourceRightsText(rights []uint32) string {
	parts := make([]string, 0, len(rights))
	for _, right := range rights {
		parts = append(parts, fmt.Sprintf("0x%08x", right))
	}
	return strings.Join(parts, ",")
}

// installSourceCurrentUserSIDText renders the invoking token's user SID. It is the
// trustee the access check runs for, so it belongs in the recorded evidence.
func installSourceCurrentUserSIDText() (string, error) {
	process, err := syscall.GetCurrentProcess()
	if err != nil {
		return "", fmt.Errorf("%w: obtain the process handle: %v", ErrAIInstallSourceProof, err)
	}
	var token syscall.Token
	if err := syscall.OpenProcessToken(process, syscall.TOKEN_QUERY, &token); err != nil {
		return "", fmt.Errorf("%w: open the process token: %v", ErrAIInstallSourceProof, err)
	}
	defer func() { _ = token.Close() }()
	user, err := token.GetTokenUser()
	if err != nil || user == nil || user.User.Sid == nil {
		return "", fmt.Errorf("%w: read the token user: %v", ErrAIInstallSourceProof, err)
	}
	sidText, err := user.User.Sid.String()
	if err != nil {
		return "", fmt.Errorf("%w: render the token user SID: %v", ErrAIInstallSourceProof, err)
	}
	return sidText, nil
}

// installSourceDescriptorDigest fingerprints the DACL of one object. A self-
// relative security descriptor returned by GetNamedSecurityInfoW is a stable byte
// rendering of that DACL, so a later loosening - or tightening - produces a
// different digest and the carrier comparison turns red.
func installSourceDescriptorDigest(name string) (string, error) {
	namePointer, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return "", fmt.Errorf("%w: encode the path %q: %v", ErrAIInstallSourceProof, name, err)
	}
	var dacl uintptr
	// descriptor is an out-parameter: GetNamedSecurityInfoW writes a pointer to a
	// self-relative security descriptor, and keeping it as *byte avoids converting a
	// uintptr back into a pointer (which go vet rightly refuses).
	var descriptor *byte
	result, _, _ := procInstallSourceGetNamedSecurityInfoW.Call(
		uintptr(unsafe.Pointer(namePointer)),
		uintptr(installSourceSEFileObject),
		uintptr(installSourceDACLSecurityInformation),
		0,
		0,
		uintptr(unsafe.Pointer(&dacl)),
		0,
		uintptr(unsafe.Pointer(&descriptor)),
	)
	if descriptor != nil {
		defer procInstallSourceLocalFree.Call(uintptr(unsafe.Pointer(descriptor)))
	}
	if result != 0 {
		return "", fmt.Errorf("%w: read the security descriptor of %q (win32 error %d)", ErrAIInstallSourceProof, name, result)
	}
	if descriptor == nil || dacl == 0 {
		return "", fmt.Errorf("%w: %q carries no DACL, so immutability cannot be recorded", ErrAIInstallSourceProof, name)
	}
	length, _, _ := procInstallSourceGetSecurityDescriptorLen.Call(uintptr(unsafe.Pointer(descriptor)))
	if length == 0 || length > 1<<20 {
		return "", fmt.Errorf("%w: the security descriptor of %q reports an unusable length %d", ErrAIInstallSourceProof, name, length)
	}
	descriptorBytes := unsafe.Slice(descriptor, int(length))
	return evidence.Digest(aiInstallSourceACLDomain, descriptorBytes), nil
}
