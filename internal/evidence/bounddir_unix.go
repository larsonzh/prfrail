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
	"runtime"
	"sort"
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

// openat2ResolveFlags performs the single confined resolution. glibc and musl expose
// no openat2 wrapper (the man page says so) and this is a plain kernel ABI call, so
// the syscall goes out directly through Syscall6: no cgo, no new dependency, no
// go.mod change.
func openat2ResolveFlags(dirfd int, clean string, flags uint64) (int, error) {
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
		mode:    0,
		resolve: boundResolveMask,
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

// boundDirVerify proves the handle still addresses the directory the caller
// named, by comparing the descriptor's device and inode with the one the expected
// path currently resolves to. A directory swapped before the bind is refused; a
// path whose ancestors are spelled differently still resolves to the same object
// and is accepted.
func boundDirVerify(dir *BoundDir, expected string) error {
	if dir == nil {
		return fmt.Errorf("evidence: nil bound directory")
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
	info, err := os.Stat(expected)
	if err != nil {
		return err
	}
	expectedStat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return fmt.Errorf("evidence: %s carries no syscall stat", expected)
	}
	if expectedStat.Dev != boundStat.Dev || expectedStat.Ino != boundStat.Ino {
		return fmt.Errorf("evidence: bound directory %s is not %s", dir.path, expected)
	}
	return nil
}

// VerifyBoundDir is the exported form of the identity re-check, kept for the
// adapters alias of the original seam.
func VerifyBoundDir(dir *BoundDir, expected string) error {
	return boundDirVerify(dir, expected)
}
