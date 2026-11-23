//go:build linux

package adapters

// T3 (finding F6), Linux half: the install source is immutable when EVERY
// component of the executable's absolute path (the installation directory
// included) and of its interpreter is non-writable by the executing user AND not
// owned by that user, and the executable object itself satisfies the same rule.
//
// Two facts are needed and neither alone is enough: "not writable" rejects an
// in-place rewrite of the bytes (which keeps the inode, so the identity check in
// requireSameExecutableFile cannot see it), and "not owned by this user" rejects a
// chmod-in-progress that would restore write access. Ownership is therefore
// required to be root or a foreign user, exactly as the contract states.
//
// P85 additions to the same walk:
//   - F06: every component is inspected with os.Lstat, so a symbolic link in the
//     NAME chain is refused instead of being followed to its target; the leaf
//     itself is opened through a bound directory handle with O_NOFOLLOW, and the
//     identity, digest and shebang/ELF header all come from that one handle.
//   - F08: "writable" is the EFFECTIVE permission of the executing user (owner,
//     other-write, or group-write for a group the user belongs to), not "some
//     subject carries a write bit": a 0775 root:root directory on a host where
//     this user is not in the group is not writable by this user.
//   - F07: an "/usr/bin/env" shebang is refused (the interpreter it picks depends
//     on the executing PATH, which the rule does not cover), the shebang chain is
//     verified recursively up to a fixed depth, and a dynamic ELF's PT_INTERP
//     loader is verified exactly like a shebang interpreter.
//
// P85 final-review additions (this round):
//   - F06: the leaf's mode and ownership are read from the BOUND HANDLE (Identity
//     and Owner), not from a second Lstat of the name, so path, mode, owner,
//     digest, identity and header are one observation of one object.
//   - F07: the interpreter's ORIGINAL name chain is what gets proven. There is no
//     EvalSymlinks that resolves the name and then checks only the target: the
//     chain is bound with openat2 RESOLVE_NO_SYMLINKS, so a link in any component
//     is refused instead of being followed, and the kernel-shebang split is
//     reproduced byte for byte (space and tab only, never strings.Fields).
//   - F07: every recursively verified interpreter level is recorded in the carrier
//     (AIInstallSourceInterpreter.Interpreter), so the inner identity is no longer
//     checked and then discarded.
//   - F08: an object that carries a POSIX access ACL is refused (fail-closed): a
//     mode/owner rule cannot exclude a named ACL grant, so it must not claim "this
//     user cannot write it".
//   - F07 (fix round 3): the name chain is proven in its ORIGINAL spelling. Every
//     path this file walks (the executable and every interpreter level) must
//     already be canonical - ".", "..", doubled separators and a trailing
//     separator are refused before anything is opened - and the shebang line is
//     parsed out of the kernel's own 256-byte window, which includes the "#!"
//     prefix, instead of the length of the remainder after it.
//
// Linux is the only supported Unix target. Every other Unix - darwin included -
// is compiled through ai_install_proof_unsupported.go and returns "cannot prove"
// instead of inventing an equivalent rule, so the probe fails closed there.

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"syscall"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	// aiInstallSourceOwnershipDomain is the digest domain of the Unix ownership
	// evidence recorded in the carrier. Declared here, not in the shared file,
	// because only the Unix half produces and consumes it.
	aiInstallSourceOwnershipDomain = "proofrail:ai-install-source-ownership:1\n"
	// aiInstallSourceInterpreterDomain is the digest domain of the interpreter
	// evidence recorded in the carrier. Declared here for the same reason as the
	// ownership domain.
	aiInstallSourceInterpreterDomain = "proofrail:ai-install-source-interpreter:1\n"
)

// installSourceShebangHeaderLength bounds the bytes inspected for a shebang and for
// the PT_INTERP entry of an ELF object. The kernel limits an interpreter line to a
// few hundred bytes; reading more would only describe a line the kernel will not use.
const installSourceShebangHeaderLength int64 = 4096

// installSourceInterpreterDepthLimit bounds the recursive shebang chain. The kernel
// has no such limit, so the depth cap is a fail-closed decision: a chain deeper than
// this cannot be proven from here and is refused (registered as a residual).
const installSourceInterpreterDepthLimit = 4

// installSourceShebangWindow is BINPRM_BUF_SIZE. The kernel copies at most this
// many bytes of the script - including the two bytes of "#!" - into its own
// buffer and parses the interpreter out of THAT window, so this rule inspects the
// same window: anything beyond it is a line the kernel will never look at.
const installSourceShebangWindow = 256

// installSourceELFMagic is the first four bytes of an ELF object.
const installSourceELFMagic = "\x7fELF"

// verifyInstallSourceImmutability observes the executable's path components and,
// when it is a script or a dynamic ELF, its interpreter. leaf is the ONE handle the
// executable was opened through; every header-based fact is taken from it.
func verifyInstallSourceImmutability(executable string, leaf *evidence.BoundFile) (AIInstallSourceImmutability, error) {
	// Second line of defence: the file is only compiled for linux, but a future
	// refactor that moves it must still refuse instead of inventing a rule.
	if runtime.GOOS != "linux" {
		return AIInstallSourceImmutability{}, fmt.Errorf("%w: %s has no install source rule, so immutability cannot be proven",
			ErrAIInstallSourceProof, runtime.GOOS)
	}
	path := executable
	if err := installSourceRequireCanonicalSpelling(path); err != nil {
		return AIInstallSourceImmutability{}, err
	}
	// The superuser can rewrite any path regardless of mode and ownership, so the
	// same three-component rule cannot prove anything for it.
	if os.Geteuid() == 0 {
		return AIInstallSourceImmutability{}, fmt.Errorf("%w: the root user can write every path, so %q cannot be proven immutable",
			ErrAIInstallSourceProof, path)
	}
	ownershipEvidence, err := verifyUnwritableInstallPathAgainstLeaf(path, leaf)
	if err != nil {
		return AIInstallSourceImmutability{}, err
	}
	interpreter, err := installSourceInterpreterFromBound(leaf, path, 0)
	if err != nil {
		return AIInstallSourceImmutability{}, err
	}
	return AIInstallSourceImmutability{Linux: &AIInstallSourceLinuxFacts{
		UnwritablePath:        true,
		OwnershipEvidenceHash: ownershipEvidence,
		Interpreter:           interpreter,
	}}, nil
}

// verifyUnwritableInstallPathAgainstLeaf is the bound form: every ancestor is
// inspected by name (a link component is refused), and the LEAF is described by
// the handle the caller already read the digest and header from (P85 F06), so the
// permissions that are recorded cannot belong to a different object than the bytes
// that are recorded.
func verifyUnwritableInstallPathAgainstLeaf(path string, leaf *evidence.BoundFile) (string, error) {
	if leaf == nil {
		return "", fmt.Errorf("%w: %q has no bound handle, so its permissions cannot be observed", ErrAIInstallSourceProof, path)
	}
	components, err := installSourcePathComponents(path)
	if err != nil {
		return "", err
	}
	if len(components) == 0 {
		return "", fmt.Errorf("%w: %q has no path components", ErrAIInstallSourceProof, path)
	}
	observed := make([]string, 0, len(components))
	for index, component := range components {
		if index == len(components)-1 {
			// The leaf: mode and owner come from the bound handle's own stat.
			_, _, mode := leaf.Identity()
			uid, gid, ok := leaf.Owner()
			if !ok {
				return "", fmt.Errorf("%w: %q carries no ownership information on its bound handle", ErrAIInstallSourceProof, component)
			}
			if err := installSourceDenyAccessACLHandle(leaf, component); err != nil {
				return "", err
			}
			if err := installSourceEffectiveWritability(component, mode, uid, gid); err != nil {
				return "", err
			}
			observed = append(observed, fmt.Sprintf("%s:%04o:%d:%d", component, mode.Perm(), uid, gid))
			continue
		}
		info, err := os.Lstat(component)
		if err != nil {
			return "", fmt.Errorf("%w: lstat %q: %v", ErrAIInstallSourceProof, component, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("%w: %q is a symbolic link, so the name chain cannot be proven", ErrAIInstallSourceProof, component)
		}
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return "", fmt.Errorf("%w: %q carries no ownership information", ErrAIInstallSourceProof, component)
		}
		if err := installSourceDenyAccessACLPath(component); err != nil {
			return "", err
		}
		if err := installSourceEffectiveWritability(component, info.Mode(), int(stat.Uid), int(stat.Gid)); err != nil {
			return "", err
		}
		observed = append(observed, fmt.Sprintf("%s:%04o:%d:%d", component, info.Mode().Perm(), stat.Uid, stat.Gid))
	}
	return evidence.Digest(aiInstallSourceOwnershipDomain, []byte(strings.Join(observed, "\n"))), nil
}

// installSourceDenyAccessACLPath refuses an ancestor that carries a POSIX access
// ACL (P85 F08). A named ACL entry can grant this user write access that the
// mode/owner rule never sees, so an object whose ACL cannot be excluded from the
// verdict is not provable by that verdict.
func installSourceDenyAccessACLPath(component string) error {
	present, err := evidence.AccessACLPath(component)
	if err != nil {
		return fmt.Errorf("%w: the access ACL of %q cannot be read: %v", ErrAIInstallSourceProof, component, err)
	}
	if present {
		return fmt.Errorf("%w: %q carries a POSIX access ACL, so a named grant to this user cannot be excluded", ErrAIInstallSourceProof, component)
	}
	return nil
}

// installSourceDenyAccessACLHandle is the leaf form: the ACL is read from the same
// bound object the digest and identity come from (P85 F06/F08).
func installSourceDenyAccessACLHandle(leaf *evidence.BoundFile, component string) error {
	present, err := leaf.AccessACL()
	if err != nil {
		return fmt.Errorf("%w: the access ACL of %q cannot be read on its bound handle: %v", ErrAIInstallSourceProof, component, err)
	}
	if present {
		return fmt.Errorf("%w: %q carries a POSIX access ACL, so a named grant to this user cannot be excluded", ErrAIInstallSourceProof, component)
	}
	return nil
}

// installSourceEffectiveWritability refuses one component that the EXECUTING user
// can make writable (P85 F08):
//
//   - the user owns it: an owner restores any permission with chmod, so the mode
//     bits do not matter;
//   - it is other-writable: every user writes it directly;
//   - it is group-writable AND the user is in that group (effective gid or the
//     supplementary groups): a group member writes it directly.
//
// The owner write bit is deliberately NOT a separate criterion: when the owner is
// this user the first rule already refuses, and when the owner is someone else the
// owner bit grants this user nothing. Judging "any subject carries a write bit"
// refused a 0775 root:root path that this user cannot write at all.
//
// This is a mode/ownership rule, not a revocation of an already-open write handle by
// the same OS user: that case is outside the S1 threat model (SECURITY.md) and is
// registered, never claimed. Named POSIX ACL grants are excluded separately, by
// refusing any component that carries an access ACL (installSourceDenyAccessACL*).
func installSourceEffectiveWritability(component string, mode os.FileMode, uid, gid int) error {
	euid := os.Geteuid()
	if uid == euid {
		return fmt.Errorf("%w: %q is owned by the executing user (mode %04o, owner %d), so the mode bits can be restored",
			ErrAIInstallSourceProof, component, mode.Perm(), uid)
	}
	if mode.Perm()&0o002 != 0 {
		return fmt.Errorf("%w: %q is writable by every user (mode %04o)", ErrAIInstallSourceProof, component, mode.Perm())
	}
	if mode.Perm()&0o020 == 0 {
		return nil
	}
	inGroup, err := installSourceEuidInGroup(gid)
	if err != nil {
		return fmt.Errorf("%w: %q carries a group write bit and the executing user's groups cannot be read: %v",
			ErrAIInstallSourceProof, component, err)
	}
	if inGroup {
		return fmt.Errorf("%w: %q grants write access to group %d, which the executing user belongs to (mode %04o)",
			ErrAIInstallSourceProof, component, gid, mode.Perm())
	}
	return nil
}

// installSourceEuidInGroup reports whether the effective user is a member of gid,
// through the effective group id or the supplementary groups. An unreadable group
// list is an error, never "not a member": the verdict has to be fail-closed.
func installSourceEuidInGroup(gid int) (bool, error) {
	if gid == os.Getegid() || gid == syscall.Getegid() {
		return true, nil
	}
	groups, err := syscall.Getgroups()
	if err != nil {
		return false, err
	}
	return slices.Contains(groups, gid), nil
}

// installSourcePathComponents walks one absolute path from the filesystem root down
// to the object itself, so a writable ancestor cannot hide behind a clean leaf.
//
// P85 F07 (fix round 3): the components come from the ORIGINAL spelling, and the
// spelling is required to be canonical first. The previous form cleaned the path
// and walked the cleaned form, so a chain carrying ".." or a doubled separator was
// proven as its folded spelling while the kernel walked the chain that was written.
// On a canonical spelling filepath.Join is an identity for one component, so the
// walk describes exactly the bytes that were named.
func installSourcePathComponents(path string) ([]string, error) {
	if err := installSourceRequireCanonicalSpelling(path); err != nil {
		return nil, err
	}
	separator := string(filepath.Separator)
	if path == separator {
		// The root itself: one component, the only spelling with a trailing
		// separator that is canonical.
		return []string{separator}, nil
	}
	components := []string{separator}
	current := separator
	for _, part := range strings.Split(strings.TrimPrefix(path, separator), separator) {
		if part == "" {
			// Unreachable on a canonical spelling; refuse rather than skip, because
			// skipping would drop a component from the proof.
			return nil, fmt.Errorf("%w: %q carries an empty component, so its components cannot be proven", ErrAIInstallSourceProof, path)
		}
		current = filepath.Join(current, part)
		components = append(components, current)
	}
	return components, nil
}

// installSourceInterpreter returns the interpreter identity when the executable is
// a script or a dynamic ELF, and nil when it is a static ELF object. It is the thin
// wrapper kept for tests and for callers that only have a path; the production path
// goes through installSourceInterpreterFromBound so the header and the identity come
// from the SAME handle the digest came from (P85 F06).
func installSourceInterpreter(path string) (*AIInstallSourceInterpreter, error) {
	bound, err := openInstallSourceLeaf(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = bound.Close() }()
	return installSourceInterpreterFromBound(bound, path, 0)
}

// installSourceInterpreterFromBound classifies one already-bound executable from its
// own bytes and returns the interpreter the kernel would re-resolve at exec time:
//
//   - "#!" script: the shebang interpreter (an env form is refused, see below);
//   - ELF object: the PT_INTERP loader of a dynamic object, nil for a static one;
//   - anything else: refused, because no rule this adapter can observe would exec it.
//
// depth counts the shebang chain already verified; an interpreter that is itself a
// script is verified one level deeper (the kernel re-resolves that chain too), up to
// installSourceInterpreterDepthLimit.
func installSourceInterpreterFromBound(bound *evidence.BoundFile, path string, depth int) (*AIInstallSourceInterpreter, error) {
	if bound == nil {
		return nil, fmt.Errorf("%w: %q has no bound handle, so its interpreter cannot be observed", ErrAIInstallSourceProof, path)
	}
	if depth > installSourceInterpreterDepthLimit {
		return nil, fmt.Errorf("%w: the interpreter chain of %q is deeper than %d levels", ErrAIInstallSourceProof, path, installSourceInterpreterDepthLimit)
	}
	data, truncated, err := bound.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w: read %q: %v", ErrAIInstallSourceProof, path, err)
	}
	if truncated {
		// A partial object cannot describe the header the kernel would read: refuse
		// instead of parsing a prefix.
		return nil, fmt.Errorf("%w: %q exceeds the bound used to inspect it", ErrAIInstallSourceProof, path)
	}
	if len(data) < 4 {
		// Too short to carry either header: no shebang and no ELF magic, so how the
		// kernel would exec this file cannot be observed. Fail closed.
		return nil, fmt.Errorf("%w: %q carries neither an ELF header nor a shebang", ErrAIInstallSourceProof, path)
	}
	if string(data[:4]) == installSourceELFMagic {
		loader, err := installSourceELFInterpreter(data, path)
		if err != nil {
			return nil, err
		}
		if loader == "" {
			// A static ELF object is exec'd directly; the path and ownership checks
			// already cover the object itself.
			return nil, nil
		}
		return installSourceVerifyInterpreter(loader, path, depth)
	}
	if data[0] != '#' || data[1] != '!' {
		// The contract proves a probe executable either as an ELF object or through a
		// shebang interpreter that satisfies the same immutability rule. A file with
		// neither property is exec'd by no rule this adapter can observe, so it must
		// not be treated as runnable: refuse it instead of assuming.
		return nil, fmt.Errorf("%w: %q is neither an ELF object nor a shebang script", ErrAIInstallSourceProof, path)
	}
	return installSourceShebangInterpreter(data, path, depth)
}

// installSourceShebangInterpreter parses the shebang line the way the kernel does
// and resolves the interpreter it names.
//
// P85 F07: the split is the kernel's, not strings.Fields'. binfmt_script skips
// leading spaces and tabs, takes every byte up to the next space or tab as the
// interpreter, and keeps the rest of the line as ONE optional argument. It never
// tokenizes on \v, \f or a non-ASCII space - and strings.Fields does, so a path that
// carried one of those bytes used to be truncated to a prefix the kernel would
// never exec, and that prefix could be a different, acceptable object.
//
// P85 F07: an "/usr/bin/env" interpreter is REFUSED. The kernel execs env, and env
// then resolves the real interpreter through the executing PATH, which this rule
// does not cover; re-implementing env's option and PATH semantics would drift from
// the kernel and reopen the same hole. The shipped probe is an ELF object, so the
// only thing refused here is an unusual external layout.
// P85 F07 (fix round 3): the limit is the kernel's own WINDOW, counted over the
// bytes the kernel buffers - "#!" included - not the length of the remainder after
// the two magic bytes. The interpreter is parsed out of the first
// installSourceShebangWindow bytes of the file only: a line whose terminating
// newline does not fall inside that window is a line the kernel never terminates,
// so the object it would exec cannot be described from here and is refused.
func installSourceShebangInterpreter(data []byte, path string, depth int) (*AIInstallSourceInterpreter, error) {
	if len(data) < 2 || data[0] != '#' || data[1] != '!' {
		return nil, fmt.Errorf("%w: %q does not open with a shebang", ErrAIInstallSourceProof, path)
	}
	window := data
	if len(window) > installSourceShebangWindow {
		window = window[:installSourceShebangWindow]
	}
	// The newline has to be INSIDE the window (file offset <= 255). A window with no
	// newline is a shebang the kernel would parse as a truncated line - the object
	// that runs would be the prefix, which this rule cannot describe.
	newline := bytes.IndexByte(window, '\n')
	if newline < 0 {
		return nil, fmt.Errorf("%w: the shebang of %q carries no newline within the kernel's %d-byte window, so what the kernel would execute cannot be proven",
			ErrAIInstallSourceProof, path, installSourceShebangWindow)
	}
	// window[2:newline] is the interpreter line WITHOUT the "#!" prefix: newline <= 255
	// bounds it at 253 bytes, exactly what the kernel keeps of it.
	line := window[2:newline]
	// Kernel rule: leading spaces and tabs only.
	specification := strings.TrimLeft(string(line), " \t")
	if specification == "" {
		return nil, fmt.Errorf("%w: %q carries an empty shebang", ErrAIInstallSourceProof, path)
	}
	interpreterPath := specification
	if index := strings.IndexAny(specification, " \t"); index >= 0 {
		// The remainder is the kernel's single optional argument (trailing blanks
		// trimmed). It is passed verbatim to the interpreter, never a second path,
		// so it does not change which object is exec'd.
		interpreterPath = specification[:index]
	}
	if interpreterPath == "" {
		return nil, fmt.Errorf("%w: %q carries an empty shebang interpreter", ErrAIInstallSourceProof, path)
	}
	if filepath.Base(interpreterPath) == "env" {
		return nil, fmt.Errorf("%w: %q delegates to env, whose interpreter depends on the executing PATH and cannot be proven",
			ErrAIInstallSourceProof, path)
	}
	if !filepath.IsAbs(interpreterPath) {
		return nil, fmt.Errorf("%w: the shebang interpreter %q of %q is not an absolute path",
			ErrAIInstallSourceProof, interpreterPath, path)
	}
	return installSourceVerifyInterpreter(interpreterPath, path, depth)
}

// installSourceVerifyInterpreter proves one interpreter: the ORIGINAL name chain is
// bound link-free from the filesystem root, every component of that chain is
// unwritable by this user, the leaf's identity, digest, permissions and ACL all come
// from that one handle, and (when the interpreter is itself a script or a dynamic
// object) its own interpreter is verified one level deeper and RECORDED.
//
// P85 F07 removed the EvalSymlinks step this function used to start with. Resolving
// the name and then checking only the resolved target described an object the kernel
// would reach through a chain this rule never inspected, so rewriting a link
// component silently changed what ran while the recorded target stayed provable.
// The kernel walks the ORIGINAL chain, so that chain is what has to be proven:
// openat2 with RESOLVE_NO_SYMLINKS refuses a link at any component, and when the
// chain carries no link the original spelling IS the real chain (there is nothing
// left for a resolver to change), which is why Path records the original spelling
// rather than a post-resolution target.
//
// Refusing a linked layout is the fail-closed branch the review allows: a chain this
// rule cannot describe is not admitted. It is a real behaviour change on hosts whose
// PT_INTERP or /bin/sh lives behind a link (merged-/usr distributions), where a
// dynamic object now refuses instead of being proven through its link target.
func installSourceVerifyInterpreter(interpreterPath, path string, depth int) (*AIInstallSourceInterpreter, error) {
	// P85 F07 (fix round 3): the interpreter's ORIGINAL spelling is validated before
	// anything is bound. A shebang or PT_INTERP entry that carries "..", ".", a
	// doubled separator or a trailing separator names a chain the kernel walks
	// literally and this rule would otherwise fold away, so it is refused whole.
	if err := installSourceRequireCanonicalSpelling(interpreterPath); err != nil {
		return nil, fmt.Errorf("%w: the interpreter %q of %q: %v", ErrAIInstallSourceProof, interpreterPath, path, err)
	}
	if !filepath.IsAbs(interpreterPath) {
		return nil, fmt.Errorf("%w: the interpreter %q of %q is not an absolute path",
			ErrAIInstallSourceProof, interpreterPath, path)
	}
	bound, err := openInstallSourceLeaf(interpreterPath)
	if err != nil {
		return nil, fmt.Errorf("%w: the interpreter %q of %q cannot be opened without following a symbolic link: %v",
			ErrAIInstallSourceProof, interpreterPath, path, err)
	}
	defer func() { _ = bound.Close() }()
	if _, err := verifyUnwritableInstallPathAgainstLeaf(interpreterPath, bound); err != nil {
		return nil, err
	}
	fileID, err := bound.FileID()
	if err != nil {
		return nil, fmt.Errorf("%w: identify the interpreter %q: %v", ErrAIInstallSourceProof, interpreterPath, err)
	}
	digest, err := bound.Hash(aiInstallSourceInterpreterDomain)
	if err != nil {
		return nil, fmt.Errorf("%w: digest the interpreter %q: %v", ErrAIInstallSourceProof, interpreterPath, err)
	}
	// P85 F07: the inner level is not just checked, it is RECORDED. The carrier has
	// to describe every interpreter the kernel will exec in turn, otherwise the
	// inner identity is verified and then dropped.
	inner, err := installSourceInterpreterFromBound(bound, interpreterPath, depth+1)
	if err != nil {
		return nil, err
	}
	return &AIInstallSourceInterpreter{Path: interpreterPath, Hash: digest, FileID: fileID, Interpreter: inner}, nil
}

// installSourceELFInterpreter reads the PT_INTERP program header of one ELF object
// from the bytes of the SAME handle the digest came from (P85 F07). It returns the
// absolute loader path, "" for a static object (no PT_INTERP), and an error for any
// malformed header: every offset and size is range-checked before it is used, so a
// crafted header is a refusal and never an out-of-range read.
func installSourceELFInterpreter(data []byte, path string) (string, error) {
	const (
		elfClass32             = 1
		elfClass64             = 2
		elfDataLSB             = 1
		elfDataMSB             = 2
		elfProgramHeaderInterp = 3
	)
	malformed := func(detail string) (string, error) {
		return "", fmt.Errorf("%w: the ELF header of %q is malformed (%s)", ErrAIInstallSourceProof, path, detail)
	}
	if len(data) < 16 {
		return malformed("truncated identification")
	}
	var order binary.ByteOrder
	switch data[5] {
	case elfDataLSB:
		order = binary.LittleEndian
	case elfDataMSB:
		order = binary.BigEndian
	default:
		return malformed("unknown data encoding")
	}
	class := data[4]
	var programOffset uint64
	var entrySize, entryCount uint16
	switch class {
	case elfClass64:
		if len(data) < 64 {
			return malformed("truncated 64-bit header")
		}
		programOffset = order.Uint64(data[32:40])
		entrySize = order.Uint16(data[54:56])
		entryCount = order.Uint16(data[56:58])
	case elfClass32:
		if len(data) < 52 {
			return malformed("truncated 32-bit header")
		}
		programOffset = uint64(order.Uint32(data[28:32]))
		entrySize = order.Uint16(data[42:44])
		entryCount = order.Uint16(data[44:46])
	default:
		return malformed("unknown class")
	}
	if entryCount == 0 {
		return "", nil
	}
	if entrySize == 0 || programOffset > uint64(len(data)) {
		return malformed("program header table outside the object")
	}
	for index := 0; index < int(entryCount); index++ {
		offset := programOffset + uint64(index)*uint64(entrySize)
		if offset > uint64(len(data)) || uint64(entrySize) > uint64(len(data))-offset {
			return malformed("program header entry outside the object")
		}
		entry := data[offset : offset+uint64(entrySize)]
		var kind uint32
		var fileOffset, fileSize uint64
		if class == elfClass64 {
			if entrySize < 56 {
				return malformed("short 64-bit program header entry")
			}
			kind = order.Uint32(entry[0:4])
			fileOffset = order.Uint64(entry[8:16])
			fileSize = order.Uint64(entry[32:40])
		} else {
			if entrySize < 32 {
				return malformed("short 32-bit program header entry")
			}
			kind = order.Uint32(entry[0:4])
			fileOffset = uint64(order.Uint32(entry[4:8]))
			fileSize = uint64(order.Uint32(entry[16:20]))
		}
		if kind != elfProgramHeaderInterp {
			continue
		}
		if fileSize == 0 || fileSize > uint64(installSourceShebangHeaderLength) {
			return malformed("implausible PT_INTERP length")
		}
		if fileOffset > uint64(len(data)) || fileSize > uint64(len(data))-fileOffset {
			return malformed("PT_INTERP path outside the object")
		}
		raw := data[fileOffset : fileOffset+fileSize]
		end := bytes.IndexByte(raw, 0)
		if end <= 0 {
			return malformed("unterminated PT_INTERP path")
		}
		interpreter := string(raw[:end])
		if !filepath.IsAbs(interpreter) {
			return "", fmt.Errorf("%w: the PT_INTERP loader %q of %q is not an absolute path",
				ErrAIInstallSourceProof, interpreter, path)
		}
		return interpreter, nil
	}
	return "", nil
}
