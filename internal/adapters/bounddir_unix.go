//go:build !windows && !darwin

package adapters

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
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"github.com/larsonzh/prfrail/internal/evidence"
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
	if dir == nil {
		return nil, fmt.Errorf("adapters: nil bound directory")
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
		return nil, fmt.Errorf("adapters: nil bound directory")
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
	// P53: openat2 confined the whole name beneath the descriptor in the syscall
	// that opened it, so containment is decided at the instant resolution completes
	// - unlike the Windows half, which has to confirm the opened file handle against
	// the bound directory's current real name.
	return evidence.NewBoundFile(os.NewFile(uintptr(fd), clean), bound, allowTruncation)
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

// openat2Resolve performs the single confined resolution. glibc and musl expose no
// openat2 wrapper (the man page says so) and this is a plain kernel ABI call, so
// the syscall goes out directly through Syscall6: no cgo, no new dependency, no
// go.mod change.
func openat2Resolve(dirfd int, clean string) (int, error) {
	name, err := syscall.BytePtrFromString(clean)
	if err != nil {
		return -1, err
	}
	how := openHow{
		// O_NOFOLLOW stays as the leaf-level belt next to RESOLVE_NO_SYMLINKS.
		flags:   syscall.O_RDONLY | syscall.O_CLOEXEC | syscall.O_NOFOLLOW,
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
		return fmt.Errorf("adapters: nil bound directory")
	}
	dir.mu.Lock()
	defer dir.mu.Unlock()
	if dir.fd < 0 {
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
