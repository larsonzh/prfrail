//go:build !windows && !darwin

package evidence

// P42-T027B-EVIDENCE-BIND Q2, Unix half: a directory descriptor opened once, and
// artifact reads resolved against that descriptor instead of a re-resolved path.
//
// P53 FIX-GAP-77-SYMLINK-ESCAPE (design r2 section 2 and section 4, user rulings
// (a)~(e) of 2026-11-04) replaced the old single openat(2) with O_NOFOLLOW, which
// only refuses a swapped *leaf* and follows a middle symlink straight out of the
// bound tree (GAP-77). Resolution now happens in one openat2(2) call relative to
// the bound descriptor with
//
//	RESOLVE_NO_XDEV | RESOLVE_NO_MAGICLINKS | RESOLVE_NO_SYMLINKS | RESOLVE_BENEATH
//
// so every component is resolved and the artifact is opened inside the same
// system call: there is no name-then-open window for a concurrent rename (F1) and
// no mount crossing, bind mounts included (F2). Ruling (a) also fixed the
// behaviour on hosts without openat2: the fallback is fail-closed with ENOSYS, not
// a step-by-step openat walk, because a walk over several system calls cannot keep
// the containment guarantee once external renames are inside the threat model.
// allowUnsafeFallback is not read here: this slice does not offer it, so a host
// without openat2 simply cannot bind reads.
//
// darwin is excluded because the standard library has no syscall.Openat there; it
// gets the fail-closed half in bounddir_darwin.go instead (validation report D7).

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"
)

const (
	// sysOpenat2Nr is SYS_openat2. go1.22.12 does not export it (syscall's
	// zsysnum_linux_amd64.go stops at SYS_PRLIMIT64), and this repository has no
	// golang.org/x/sys dependency, so the number lives here as a private constant.
	// It is only used behind the runtime linux/amd64 gate below: on other unix
	// targets 437 means something else entirely.
	sysOpenat2Nr = 437

	// openHowSize is the ABI size of struct open_how: three uint64 fields. openat2
	// rejects any other size with EINVAL.
	openHowSize = 24

	// resolveNoXdev refuses any mount crossing, bind mounts included. Ruling (b)
	// put mount imports on the attack surface, so the cost of refusing a legal
	// sub-mount inside the bound tree is accepted.
	resolveNoXdev = 0x01
	// resolveNoMagiclinks refuses /proc style magic links explicitly.
	resolveNoMagiclinks = 0x02
	// resolveNoSymlinks refuses a symlink in any component, not just the leaf.
	resolveNoSymlinks = 0x04
	// resolveBeneath confines the whole resolution beneath the bound descriptor.
	resolveBeneath = 0x08

	// boundResolveMask is the single mask every bound read uses: NO_XDEV 0x01 |
	// NO_MAGICLINKS 0x02 | NO_SYMLINKS 0x04 | BENEATH 0x08.
	boundResolveMask = resolveNoXdev | resolveNoMagiclinks | resolveNoSymlinks | resolveBeneath

	// resolveNoSymlinksMask is the mask used when a directory OUTSIDE the bound
	// tree has to be bound link-free from the filesystem root (P85 F02/F20).
	// NO_XDEV is deliberately absent: a policy or proof root may legally live on
	// its own mount (a tmpfs run root, a separate /home), and refusing a mount
	// crossing there would refuse every legal layout instead of only a link. What
	// matters is that no component is a symbolic link, which RESOLVE_NO_SYMLINKS
	// enforces for the whole chain inside the one resolution.
	resolveNoSymlinksMask = resolveNoMagiclinks | resolveNoSymlinks | resolveBeneath
)

// openHow mirrors struct open_how from linux/openat2.h: three little-endian
// uint64 fields, no padding on linux/amd64.
type openHow struct {
	flags   uint64
	mode    uint64
	resolve uint64
}

// openat2Disabled records that this host answered ENOSYS once. Resolution is
// already fail-closed there; the flag only stops repeating a syscall that cannot
// succeed and keeps the refusal stable for the life of the process.
var openat2Disabled atomic.Bool

// BoundDir is one open directory descriptor. It owns the descriptor until Close.
//
// mu closes the descriptor lifetime (F4): openBoundDirFile, boundDirVerify and
// Close all read or replace fd under it, so the descriptor cannot be closed while
// a resolution is using it and the number cannot be recycled underneath the
// caller.
type BoundDir struct {
	fd   int
	path string
	mu   sync.Mutex
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
		return nil, fmt.Errorf("evidence: %s is not a directory", path)
	}
	return &BoundDir{fd: fd, path: path}, nil
}

// OpenBoundDirNoFollow opens path as a directory descriptor WITHOUT following a
// symbolic link at any component (P85 F02/F06/F20).
//
// OpenBoundDir opens the name as given, so a link anywhere in the chain is
// followed; the descriptor then addresses the link TARGET, and a later
// VerifyBoundDir compares that target with a fresh stat of the path - which
// follows the same link, so the two agree no matter which directory the name was
// pointed at. Agreement between two resolutions of the same name therefore proves
// nothing about WHICH directory the caller believed it was checking.
//
// This entry point removes the resolution from the caller's hands: the whole
// absolute path is resolved from the filesystem root in ONE openat2(2) call with
// RESOLVE_NO_SYMLINKS | RESOLVE_NO_MAGICLINKS | RESOLVE_BENEATH, so the descriptor
// is the object the name addresses at that instant and no component of it can be
// a link. A path that carries a link is refused (ELOOP) instead of being followed,
// which is the same rule the per-component pre-check applies to the name - the
// difference is that this verdict is taken on the object the caller then uses.
//
// Hosts without openat2 are refused with ENOSYS exactly like a bound artifact
// open: a step-by-step walk cannot keep the guarantee.
func OpenBoundDirNoFollow(path string) (*BoundDir, error) {
	return openBoundDirNoFollow(path)
}

// OpenBoundDirVerified binds path AND reports the identity of the object it bound,
// from ONE resolution (P85 F02/F06/F20).
//
// This is the shape a caller needs when the binding itself is the pre-check: the
// single openat2(2) call below performs the component-level link refusal and
// returns the descriptor, and the identity is read from that descriptor, so
// "the object that was checked" and "the handle that is used" cannot be two
// different objects - there is no name resolution between the two.
func OpenBoundDirVerified(path string) (*BoundDir, FileID, error) {
	dir, err := openBoundDirNoFollow(path)
	if err != nil {
		return nil, "", err
	}
	identity, err := dir.FileID()
	if err != nil {
		_ = dir.Close()
		return nil, "", err
	}
	return dir, identity, nil
}

// OpenBoundDirFile opens one artifact relative to the bound directory. The
// returned handle is the artifact's identity: the caller reads through it and
// never resolves the name again.
func OpenBoundDirFile(dir *BoundDir, rel string, bound int64) (*BoundFile, error) {
	return openBoundDirFile(dir, rel, bound, false)
}

// OpenBoundDirFileClass is the class-aware form: allowTruncation keeps the log and
// event class readable past its bound instead of refusing the open.
func OpenBoundDirFileClass(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*BoundFile, error) {
	return openBoundDirFile(dir, rel, bound, allowTruncation)
}

// openBoundDirFile is the class-aware form: allowTruncation keeps the log and
// event class readable past its bound instead of refusing the open.
func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*BoundFile, error) {
	if dir == nil {
		return nil, fmt.Errorf("evidence: nil bound directory")
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return nil, fmt.Errorf("evidence: nil bound directory")
	}
	// P84 T4(d): the bound object itself is not an artifact, and both root spellings
	// ("", ".") are refused through one shared error so the two halves cannot drift.
	if boundRootName(rel) {
		return nil, errBoundRootArtifact(rel)
	}
	clean, err := cleanBoundRelativePath(rel)
	if err != nil {
		return nil, err
	}
	fd, err := openBoundArtifactLocked(dir.fd, clean)
	if err != nil {
		// Error contract unchanged: one *os.PathError from the bound open, whose
		// Err carries the raw errno, so errors.Is(err, fs.ErrNotExist) keeps
		// working for a missing artifact.
		return nil, &os.PathError{Op: "openat", Path: clean, Err: err}
	}
	// P83 FIX ROUND 4 (Item B): the open above is non-blocking, so a FIFO with no
	// writer returns an fd instead of parking the collector forever. The type is
	// then checked on the returned handle, before any read, and only a regular file
	// is admitted: anything else is refused as an unsafe artifact with the same
	// sentinel the other platforms return.
	if err := boundArtifactRegular(fd, clean); err != nil {
		_ = syscall.Close(fd)
		return nil, err
	}
	// P53: openat2 confined the whole name beneath the descriptor in the syscall
	// that opened it, so containment is decided at the instant resolution completes
	// - unlike the Windows half, which has to confirm the opened file handle against
	// the bound directory's current real name.
	return NewBoundFile(os.NewFile(uintptr(fd), clean), bound, allowTruncation)
}

// boundArtifactRegular fstats an already open artifact descriptor and refuses
// anything that is not a regular file.
//
// P83 FIX ROUND 4 (Item B): the check lives on the returned handle rather than on
// the name, so it describes the very object the caller will read and cannot be
// raced by a rename or a type swap between an open and a later stat. A FIFO that
// has no writer is the case that used to hang the collector: the non-blocking open
// returns, and this refusal makes it a missing/unsafe artifact instead of a read
// that never ends.
func boundArtifactRegular(fd int, clean string) error {
	var stat syscall.Stat_t
	if err := syscall.Fstat(fd, &stat); err != nil {
		return fmt.Errorf("evidence: stat %s: %v", clean, err)
	}
	if uint32(stat.Mode)&uint32(syscall.S_IFMT) != uint32(syscall.S_IFREG) {
		return fmt.Errorf("%w: %s has mode %#o", ErrBoundFileNotRegular, clean, uint32(stat.Mode))
	}
	return nil
}

// openBoundArtifactLocked resolves clean against dirfd. The caller holds dir.mu.
//
// Main path: one openat2(2) call with boundResolveMask (design r2 section 2). Every
// errno is returned as-is: ELOOP, EXDEV, ENOTDIR, ENOENT, EACCES, EAGAIN and EINTR
// all fail the open rather than being retried or second-guessed.
//
// Ruling (a): openat2 is required. A host that cannot provide it - a kernel below
// 5.6 answering ENOSYS, or any unix build that is not linux/amd64 - is refused with
// ENOSYS instead of being served by a step-by-step openat walk. The walk is
// deliberately not implemented: with external renames inside the threat model it
// would reopen the very window this slice closes, and this slice offers no
// allowUnsafeFallback switch.
func openBoundArtifactLocked(dirfd int, clean string) (int, error) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return -1, syscall.ENOSYS
	}
	if openat2Disabled.Load() {
		return -1, syscall.ENOSYS
	}
	fd, err := openat2Resolve(dirfd, clean)
	if err != nil {
		if errors.Is(err, syscall.ENOSYS) {
			openat2Disabled.Store(true)
		}
		return -1, err
	}
	return fd, nil
}

// openat2Resolve performs the artifact resolution: the confined open with the
// artifact open flags (see artifactOpenFlags).
func openat2Resolve(dirfd int, clean string) (int, error) {
	return openat2ResolveFlags(dirfd, clean, artifactOpenFlags)
}

// openBoundDirectoryLocked resolves a DIRECTORY name beneath dirfd with the same
// confinement an artifact open uses (P83 fix round 7, item T4b). rel "." names the
// bound directory itself.
func openBoundDirectoryLocked(dirfd int, clean string) (int, error) {
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return -1, syscall.ENOSYS
	}
	if openat2Disabled.Load() {
		return -1, syscall.ENOSYS
	}
	fd, err := openat2ResolveFlags(dirfd, clean, artifactOpenFlags|syscall.O_DIRECTORY)
	if err != nil {
		if errors.Is(err, syscall.ENOSYS) {
			openat2Disabled.Store(true)
		}
		return -1, err
	}
	return fd, nil
}

// artifactOpenFlags are the flags every bound artifact open uses: read-only, close
// on exec, no leaf symlink, and non-blocking so a FIFO or a device node without a
// writer cannot park the caller (P83 fix round 4, item B).
const artifactOpenFlags = syscall.O_RDONLY | syscall.O_CLOEXEC | syscall.O_NOFOLLOW | syscall.O_NONBLOCK

// BoundDirNames lists the entry names of one directory inside the bound directory
// object, sorted (P83 fix round 7, item T4b).
//
// The listing is read from a descriptor resolved beneath the bound descriptor - not
// from the caller's path - so the enumeration and the artifact reads describe the
// same object even when the directory was renamed or replaced after it was bound.
// rel == "" (or ".") names the bound object itself: both spellings go through
// boundRootName, so the classifier - not the artifact-name guard, which refuses
// "." - decides that (P84 T4(d)). A platform that cannot enumerate through the
// descriptor returns ErrBoundDirUnsupported instead of falling back to a path read.
func BoundDirNames(dir *BoundDir, rel string) ([]string, error) {
	if dir == nil {
		return nil, fmt.Errorf("evidence: nil bound directory")
	}
	clean := "."
	if !boundRootName(rel) {
		cleaned, err := cleanBoundRelativePath(rel)
		if err != nil {
			return nil, err
		}
		clean = cleaned
	}
	dir.mu.Lock()
	if dir.fd < 0 {
		dir.mu.Unlock()
		return nil, fmt.Errorf("evidence: nil bound directory")
	}
	fd, err := openBoundDirectoryLocked(dir.fd, clean)
	dir.mu.Unlock()
	if err != nil {
		return nil, &os.PathError{Op: "openat", Path: rel, Err: err}
	}
	file := os.NewFile(uintptr(fd), clean)
	defer func() { _ = file.Close() }()
	// T4 (P83 fix round 8): Readdirnames reads the entries with getdents on the
	// descriptor itself and returns names only. os.File.ReadDir must classify each
	// entry, so an entry whose directory record carries DT_UNKNOWN - which several
	// file systems report - is lstat'ed through its relative name (and the directory
	// itself through "."), i.e. resolved against the process working directory
	// instead of against the object this descriptor addresses. Names therefore come
	// from the directory the descriptor names, and every candidate is still resolved
	// with openat2 beneath the bound descriptor when it is opened, so the listing and
	// the reads describe one object.
	names, readErr := file.Readdirnames(-1)
	if readErr != nil {
		return nil, readErr
	}
	sort.Strings(names)
	return names, nil
}

// openBoundDirNoFollow is the Unix half of OpenBoundDirNoFollow.
func openBoundDirNoFollow(path string) (*BoundDir, error) {
	clean := filepath.Clean(path)
	if !filepath.IsAbs(clean) {
		return nil, &os.PathError{Op: "open", Path: path, Err: syscall.EINVAL}
	}
	// The root itself has no component to follow, so the plain open is exact.
	if clean == string(filepath.Separator) {
		return OpenBoundDir(clean)
	}
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		return nil, &os.PathError{Op: "open", Path: path, Err: syscall.ENOSYS}
	}
	if openat2Disabled.Load() {
		return nil, &os.PathError{Op: "open", Path: path, Err: syscall.ENOSYS}
	}
	root, err := syscall.Open(string(filepath.Separator), syscall.O_RDONLY|syscall.O_CLOEXEC|syscall.O_DIRECTORY, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: string(filepath.Separator), Err: err}
	}
	defer func() { _ = syscall.Close(root) }()
	relative := strings.TrimPrefix(clean, string(filepath.Separator))
	fd, err := openat2ResolveFlagsMask(root, relative, artifactOpenFlags|syscall.O_DIRECTORY, 0, resolveNoSymlinksMask)
	if err != nil {
		if errors.Is(err, syscall.ENOSYS) {
			openat2Disabled.Store(true)
		}
		return nil, &os.PathError{Op: "openat", Path: path, Err: err}
	}
	var stat syscall.Stat_t
	if err := syscall.Fstat(fd, &stat); err != nil {
		_ = syscall.Close(fd)
		return nil, &os.PathError{Op: "stat", Path: path, Err: err}
	}
	if uint32(stat.Mode)&uint32(syscall.S_IFMT) != uint32(syscall.S_IFDIR) {
		_ = syscall.Close(fd)
		return nil, fmt.Errorf("evidence: %s is not a directory", path)
	}
	return &BoundDir{fd: fd, path: path}, nil
}

// openat2ResolveFlags performs the single confined resolution. glibc and musl expose
// no openat2 wrapper (the man page says so) and this is a plain kernel ABI call, so
// the syscall goes out directly through Syscall6: no cgo, no new dependency, no
// go.mod change.
func openat2ResolveFlags(dirfd int, clean string, flags uint64) (int, error) {
	return openat2ResolveFlagsMask(dirfd, clean, flags, 0, boundResolveMask)
}

// openat2ResolveFlagsMask is the explicit form: the caller supplies the open
// flags, the creation mode (open_how.mode, only read when O_CREAT/O_TMPFILE is
// set) and the resolution mask. The bound artifact and directory opens use
// boundResolveMask; OpenBoundDirNoFollow uses resolveNoSymlinksMask.
func openat2ResolveFlagsMask(dirfd int, clean string, flags, mode, resolve uint64) (int, error) {
	name, err := syscall.BytePtrFromString(clean)
	if err != nil {
		return -1, err
	}
	how := openHow{
		// The caller supplies the open flags; artifactOpenFlags carries the
		// leaf-level O_NOFOLLOW belt next to RESOLVE_NO_SYMLINKS and the O_NONBLOCK
		// that keeps a FIFO without a writer from blocking the collector (P83 FIX
		// ROUND 4, item B). O_NONBLOCK is a no-op for a regular file, which is the
		// only type the caller then accepts (checked with fstat on the returned
		// descriptor). A directory listing adds O_DIRECTORY to the same set.
		flags:   flags,
		mode:    mode,
		resolve: resolve,
	}
	fd, _, errno := syscall.Syscall6(
		sysOpenat2Nr,
		uintptr(dirfd),
		uintptr(unsafe.Pointer(name)),
		uintptr(unsafe.Pointer(&how)),
		uintptr(openHowSize),
		0,
		0,
	)
	if errno != 0 {
		return -1, errno
	}
	return int(fd), nil
}

// RootPath reports the directory name the bound descriptor was opened at. T4 (P83
// fix round 6): it is the name the attempt's work directory is derived from, and
// VerifyBoundDir is what proves that name still addresses this descriptor. T4 (P83
// fix round 7) removed its other use: the artifact enumeration now reads the entries
// from the descriptor itself (BoundDirNames) instead of globbing this path, so no
// enumeration depends on a name. Unix has no portable handle-to-name mapping, so the
// name of record is the one captured at bind time.
func (dir *BoundDir) RootPath() (string, error) {
	if dir == nil {
		return "", fmt.Errorf("evidence: nil bound directory")
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return "", fmt.Errorf("evidence: nil bound directory")
	}
	return dir.path, nil
}

// FileID reports the platform identity of the directory object this descriptor
// addresses, taken from fstat(2) on the descriptor and never from a name (P85
// F02/F06/F20). It is the value OpenBoundDirVerified hands back and the value
// VerifyBoundDir compares against.
func (dir *BoundDir) FileID() (FileID, error) {
	if dir == nil {
		return "", fmt.Errorf("evidence: nil bound directory")
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return "", fmt.Errorf("evidence: nil bound directory")
	}
	var stat syscall.Stat_t
	if err := syscall.Fstat(dir.fd, &stat); err != nil {
		return "", fmt.Errorf("evidence: identify the bound directory: %v", err)
	}
	return FileID(fmt.Sprintf("unix:%d:%d", stat.Dev, stat.Ino)), nil
}

// Descriptor returns the numeric descriptor behind the bound directory handle.
//
// P85 design rev 4, section 1: the write-domain observation has to take its
// permission facts (mode, ownership, effective write) from the SAME object the
// reads and writes go through. The descriptor is the only way to ask the kernel
// about that object, so it is exported next to BoundFile.Descriptor instead of
// being re-derived from a name. The caller must not close it: it stays owned by
// the BoundDir and the accessor exists for read-only observations.
func (dir *BoundDir) Descriptor() (uintptr, error) {
	if dir == nil {
		return 0, fmt.Errorf("evidence: nil bound directory")
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return 0, fmt.Errorf("evidence: nil bound directory")
	}
	return uintptr(dir.fd), nil
}

// SyncDirectory flushes this directory object with fsync(2) on the bound
// descriptor (P85 F10). The previous shape re-opened the directory BY NAME to
// sync it, which proved nothing about which object was synced; with the
// descriptor held, the object that was read is the object that is flushed.
func (dir *BoundDir) SyncDirectory() error {
	if dir == nil {
		return fmt.Errorf("evidence: nil bound directory")
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return fmt.Errorf("evidence: nil bound directory")
	}
	if err := syscall.Fsync(dir.fd); err != nil {
		return fmt.Errorf("evidence: sync the bound directory: %v", err)
	}
	return nil
}

// Close releases the directory descriptor. It is idempotent and safe against a
// concurrent read: the descriptor is closed and the field invalidated under mu, so
// a read either finishes before the close or sees fd == -1 and is refused.
func (dir *BoundDir) Close() error {
	if dir == nil {
		return nil
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return nil
	}
	err := syscall.Close(dir.fd)
	dir.fd = -1
	return err
}

// boundDirVerify proves the handle still addresses the object whose identity the
// caller captured when the binding was made (P85 F02/F06/F20).
//
// The comparison used to be handle-versus-path: the descriptor's dev:ino was
// compared with a fresh os.Stat of the name, which only proved that two
// resolutions of the SAME name agreed - after a directory was swapped, both sides
// resolved to the replacement and the check stayed green. The expected value is
// now the identity the pre-check produced (FileID), and it is compared with the
// identity of the handle, so the check answers "is this still the object that was
// approved" instead of "does the name still resolve to itself".
func boundDirVerify(dir *BoundDir, expected FileID) error {
	if dir == nil {
		return fmt.Errorf("evidence: nil bound directory")
	}
	if expected == "" {
		return fmt.Errorf("evidence: the expected directory identity is empty")
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return fmt.Errorf("evidence: nil bound directory")
	}
	var boundStat syscall.Stat_t
	if err := syscall.Fstat(dir.fd, &boundStat); err != nil {
		return fmt.Errorf("evidence: stat the bound directory: %v", err)
	}
	actual := FileID(fmt.Sprintf("unix:%d:%d", boundStat.Dev, boundStat.Ino))
	if actual != expected {
		return fmt.Errorf("evidence: bound directory %s is not the approved object (%s != %s)", dir.path, actual, expected)
	}
	return nil
}

// VerifyBoundDir is the exported form of the identity re-check: the handle's own
// identity must be the one the binding produced (FileID), never a fresh
// resolution of a name.
func VerifyBoundDir(dir *BoundDir, expected FileID) error {
	return boundDirVerify(dir, expected)
}
