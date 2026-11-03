//go:build !windows && !darwin

package adapters

// P52-S2 ACCEPTANCE-SECURITY-COUNTEREXAMPLES recorded the middle-symlink escape on
// the bound-directory seam as GAP-77 (V-1, HIGH): the old openBoundDirFile resolved
// "link/secret" with a single openat(2) plus O_NOFOLLOW, and O_NOFOLLOW only refuses
// the *leaf*, so the middle component "link" was followed and the returned
// descriptor addressed O/secret outside the bound directory. The Linux run that
// proved it is recorded in
// docs/validation/evidence/ACCEPTANCE-SECURITY-COUNTEREXAMPLES/05-attack/
// t51-1-linux-escape.md (ATTACK-SUCCEEDED / SECRET-OUTSIDE-BOUNDARY).
//
// P53 FIX-GAP-77-SYMLINK-ESCAPE replaced that resolution with one confined
// openat2(2) call (design r2 sections 2 and 4, user rulings (a)~(e) of 2026-11-04).
// This file now carries the permanent, always-green regression; the payload it
// shares with the red baseline lives in bounddir_escape_payload_unix_test.go.
//
//   - TestBoundDirMiddleSymlinkEscapesBound keeps its red-baseline name so the
//     T51-1 evidence chain stays traceable, but it now asserts the post-fix
//     criterion of design r2 section 5: the escape is refused with an
//     *os.PathError whose errno is inside the security closure {ELOOP, ENOTDIR,
//     EXDEV}, and an in-bound three-level name still reads byte-exactly, so a
//     blanket refusal cannot pass as a fix.
//   - TestBoundDirMiddleSymlinkRefusedWithClosedErrno complements it on a different
//     axis and deliberately does not repeat the closure assertion: it pins the
//     *os.PathError contract (op and path), the attack-spelling matrix under the
//     path cleaner, the determinism of the refusal, and a single-segment positive
//     control.
//
// The red baseline itself is not lost: it moved to
// bounddir_escape_redbaseline_test.go, where the ⑧ role replays it as the mutation
// check. The separation is by name plus a runtime gate, not by build tag: both
// files belong to the ordinary unix build (a bespoke tag made an editor tooling
// such as gopls report "no packages found"), the baseline owns the distinct name
// TestBoundDirRedBaselineMiddleSymlinkEscape, and it skips itself unless
// PRFRAIL_P53_RED_BASELINE=1. The default suite therefore never sees it fail.

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
	"testing"
)

// gap77ErrnoInClosure reports whether a refusal errno is inside the security
// closure {ELOOP, ENOTDIR, EXDEV}. ENOENT, EACCES, EPERM, EBADF, EMFILE, EINTR,
// EIO and ENOSYS are all outside it, so an arbitrary error cannot pass as a
// refusal.
func gap77ErrnoInClosure(errno error) bool {
	return errors.Is(errno, syscall.ELOOP) ||
		errors.Is(errno, syscall.ENOTDIR) ||
		errors.Is(errno, syscall.EXDEV)
}

// gap77SkipUnlessOpenat2 skips on builds where the openat2 main path of design r2
// section 2 is not the resolution path: those hosts fail closed with ENOSYS, which
// is pinned by their own criteria, so asserting the refusal closure there would pin
// the wrong contract rather than a weaker one.
//
// The platform check alone is not enough. openat2 landed in Linux 5.6, and a
// sandbox that filters the syscall (seccomp) answers ENOSYS or EPERM on a kernel
// that has it, so the helper also probes the kernel with the production resolution
// call once per process. Without that probe every success-shaped assertion in this
// package ("a successful read is byte-exact", "the storm produced a read") would
// turn red on such a host for an environment reason.
func gap77SkipUnlessOpenat2(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skipf("the GAP-77 refusal closure is claimed on linux/amd64 only (openat2 main path, design r2 section 2); other unix targets fail closed with ENOSYS. This build is %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	switch err := gap77ProbeOpenat2(t); {
	case errors.Is(err, fs.ErrNotExist):
		// The kernel resolved the absent name: the main path is live.
	case errors.Is(err, syscall.ENOSYS):
		t.Skipf("fail-closed host (no openat2), refusal closure not applicable here: the production path answers ENOSYS for every bound read")
	default:
		t.Skipf("the openat2 main path is not usable on this host (probe answered %v); the closure criteria need the kernel resolution path, and a filtered syscall must not be reported as a failed guard", err)
	}
}

// gap77Openat2Probe caches the single kernel probe of the process.
var gap77Openat2Probe struct {
	once sync.Once
	err  error
}

// gap77ProbeOpenat2 asks the kernel for one non-existent name beneath a fresh
// temporary directory, through the very call the production path uses. ENOENT
// means openat2 works, ENOSYS means the kernel does not have it, and any other
// errno means something (seccomp, seccomp plus a policy) obstructs it.
//
// The probe has no side effect: the probed name never exists, nothing is created or
// written, the temporary directory belongs to the test framework, and
// openat2Disabled is deliberately left untouched - the probe must not change the
// behaviour of the code under test.
func gap77ProbeOpenat2(t *testing.T) error {
	t.Helper()
	gap77Openat2Probe.once.Do(func() {
		root := t.TempDir()
		dirfd, err := syscall.Open(root, syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
		if err != nil {
			gap77Openat2Probe.err = fmt.Errorf("open the probe directory: %w", err)
			return
		}
		defer func() { _ = syscall.Close(dirfd) }()
		_, gap77Openat2Probe.err = openat2Resolve(dirfd, "gap77-openat2-probe-absent")
	})
	return gap77Openat2Probe.err
}

// TestBoundDirMiddleSymlinkEscapesBound is the permanent GAP-77 regression: the
// T51-1 red baseline, kept under its original name and turned into the post-fix
// assertion (design r2 section 5, green criterion; user rulings (a)~(e)).
//
// The escape-reproducing form of this criterion moved to
// TestBoundDirRedBaselineMiddleSymlinkEscape in bounddir_escape_redbaseline_test.go;
// it is armed by PRFRAIL_P53_RED_BASELINE=1 and skipped otherwise, so the two
// never report a verdict in the same run.
func TestBoundDirMiddleSymlinkEscapesBound(t *testing.T) {
	gap77SkipUnlessOpenat2(t)

	tree := buildGAP77Tree(t)

	dir, err := OpenBoundDir(tree.boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	attackName := filepath.Join(gap77LinkName, gap77LeafName)
	file, err := OpenBoundDirFile(dir, attackName, 1024)
	if err == nil {
		data, truncated, readErr := file.ReadAll()
		_ = file.Close()
		t.Fatalf("GAP-77 regression: the middle symlink %q must be refused, but the open succeeded and returned %d bytes %q (truncated %v, read error %v)",
			attackName, len(data), data, truncated, readErr)
	}

	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("the refusal must be an *os.PathError, got %T: %v", err, err)
	}
	if !gap77ErrnoInClosure(pathErr.Err) {
		t.Fatalf("the refusal errno %v (%q) is outside the security closure {ELOOP, ENOTDIR, EXDEV}: an arbitrary error must not pass as a refusal", pathErr.Err, pathErr.Err)
	}

	// Positive control: the same bound directory, the same call shape, a legal
	// three-level in-bound name. Without this a "refuse every relative name" change
	// would look like a fix.
	insideName := filepath.Join("sub", "nested", "inside")
	inside, err := OpenBoundDirFile(dir, insideName, 1024)
	if err != nil {
		t.Fatalf("a legal in-bound multi-level name %q must still open: a blanket refusal is not a fix (%v)", insideName, err)
	}
	defer func() { _ = inside.Close() }()

	data, truncated, err := inside.ReadAll()
	if err != nil {
		t.Fatalf("reading the in-bound multi-level artifact %q failed: %v", insideName, err)
	}
	if truncated {
		t.Fatalf("the %d byte in-bound artifact must not be truncated at the 1024 byte bound", len(gap77InsideCanary))
	}
	if string(data) != gap77InsideCanary {
		t.Fatalf("the in-bound multi-level read returned %d bytes %q, want %q", len(data), data, gap77InsideCanary)
	}
}

// TestBoundDirMiddleSymlinkRefusedWithClosedErrno complements
// TestBoundDirMiddleSymlinkEscapesBound. The two were written as a red/green pair in
// the P53 baseline round; now that both run on the fixed tree the closure assertion
// is owned by the regression above, so this test does not repeat it.
//
// What it pins instead:
//   - the refusal is an *os.PathError from the bound open with op "openat", naming
//     the cleaned artifact: the error contract of the seam is unchanged;
//   - every spelling the path cleaner folds to the same artifact ("link/secret",
//     "./link/secret", "link/../link/secret", "sub/../link/secret") is refused with
//     the same errno, so a per-spelling special case shows up as drift here;
//   - repeated attempts are refused as well, so a one-shot refusal cannot pass as a
//     closed door;
//   - a single-segment in-bound name still opens and reads byte-exactly.
func TestBoundDirMiddleSymlinkRefusedWithClosedErrno(t *testing.T) {
	gap77SkipUnlessOpenat2(t)

	tree := buildGAP77Tree(t)

	dir, err := OpenBoundDir(tree.boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	spellings := []string{
		"link/secret",
		"./link/secret",
		"link/../link/secret",
		"sub/../link/secret",
	}
	var canonical error
	for _, spelling := range spellings {
		wantPath := filepath.Clean(filepath.FromSlash(spelling))
		for attempt := 1; attempt <= 2; attempt++ {
			file, err := OpenBoundDirFile(dir, spelling, 1024)
			if err == nil {
				data, _, readErr := file.ReadAll()
				_ = file.Close()
				t.Fatalf("the middle symlink spelling %q must be refused (attempt %d), but the open succeeded and returned %q (read error %v)",
					spelling, attempt, data, readErr)
			}

			var pathErr *os.PathError
			if !errors.As(err, &pathErr) {
				t.Fatalf("the refusal of %q (attempt %d) must be an *os.PathError, got %T: %v", spelling, attempt, err, err)
			}
			if pathErr.Op != "openat" {
				t.Fatalf("the refusal of %q must come from the bound open (op %q), got op %q: %v", spelling, "openat", pathErr.Op, err)
			}
			if pathErr.Path != wantPath {
				t.Fatalf("the refusal of %q must name the cleaned artifact %q, got %q", spelling, wantPath, pathErr.Path)
			}
			if canonical == nil {
				canonical = pathErr.Err
				continue
			}
			if !errors.Is(pathErr.Err, canonical) {
				t.Fatalf("the refusal errno drifted with the spelling or the attempt: %q answered %v, the canonical spelling answered %v", spelling, pathErr.Err, canonical)
			}
		}
	}

	// Single-segment positive control, complementary to the three-level control in
	// the regression: a fix that only admits sub-directory paths would show up here.
	originalName := "original"
	original, err := OpenBoundDirFile(dir, originalName, 1024)
	if err != nil {
		t.Fatalf("a legal single-segment in-bound name %q must still open: %v", originalName, err)
	}
	defer func() { _ = original.Close() }()

	data, truncated, err := original.ReadAll()
	if err != nil {
		t.Fatalf("reading the in-bound artifact %q failed: %v", originalName, err)
	}
	if truncated {
		t.Fatalf("the %d byte in-bound artifact %q must not be truncated at the 1024 byte bound", len(gap77InsideCanary), originalName)
	}
	if string(data) != gap77InsideCanary {
		t.Fatalf("the in-bound read of %q returned %d bytes %q, want %q", originalName, len(data), data, gap77InsideCanary)
	}
}
