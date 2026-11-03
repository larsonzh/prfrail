//go:build !windows && !darwin

package adapters

// P53 FIX-GAP-77-SYMLINK-ESCAPE: the fail-closed half of the seam contract and the
// Close/read lifetime contract of a shared BoundDir.
//
// Ruling (a) of 2026-11-04 removed every fallback. A host that cannot resolve with
// openat2(2) - a kernel below 5.6, a sandbox that filters the syscall, or any build
// that is not linux/amd64 - refuses every bound read with ENOSYS instead of walking
// the path component by component, because a multi-step walk would reopen the very
// name-then-open window this slice closes once external renames are inside the
// threat model. That refusal is a contract of the seam, not an accident, so it is
// pinned here: *os.PathError, op "openat", errno ENOSYS, no panic, and a latched
// openat2Disabled flag that keeps the answer stable for the life of the process.
//
// The second test pins design r2 section 5 item F4 (descriptor lifetime). BoundDir
// Close and a concurrent read share dir.mu, so a read either finishes against the
// still-open descriptor or is refused with the nil-bound-directory error - never
// with EBADF, and never with bytes from outside the bound.
//
// Both tests are in-package tests because they pin private symbols
// (openat2Disabled, openBoundArtifactLocked) that an external test package cannot
// reach.

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// gap77NilBoundDirMessage is the exact refusal a nil or already closed bound
// directory produces: openBoundDirFile rejects fd == -1 under dir.mu instead of
// handing the released number to the kernel. The seam exports no sentinel for this
// state, so the tests compare the message; a change to either side shows up here.
const gap77NilBoundDirMessage = "adapters: nil bound directory"

// TestBoundDirFailClosedWithoutOpenat2 is the ENOSYS contract of ruling (a).
//
// openat2Disabled is forced rather than waiting for a kernel that lacks openat2:
// the flag is exactly what the main path consults after the first refusal, so the
// latched state and the refusal it implies are the contract under test. The seam
// must answer with one *os.PathError, op "openat", carrying ENOSYS - not a walk, not
// a plain *syscall.Errno, not a panic, and not a handle.
func TestBoundDirFailClosedWithoutOpenat2(t *testing.T) {
	// The latch is package-global. Tests in this package run serially (nothing here
	// calls t.Parallel), so no other test can observe the transient true; the
	// cleanup restores the pre-test state on every exit path anyway.
	openat2Disabled.Store(true)
	t.Cleanup(func() { openat2Disabled.Store(false) })

	tree := buildGAP77Tree(t)
	dir, err := OpenBoundDir(tree.boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	// A legal in-bound name is used on purpose: the refusal must come from the
	// fail-closed gate, not from the name being hostile.
	const rel = "original"

	var (
		openErr   error
		gotHandle bool
	)
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("the fail-closed refusal must return an error, not panic: %v", r)
			}
		}()
		file, err := OpenBoundDirFile(dir, rel, 1024)
		if file != nil {
			gotHandle = true
			_ = file.Close()
		}
		openErr = err
	}()

	if gotHandle {
		t.Fatal("a fail-closed host answered ENOSYS but still handed out an artifact handle; the refusal must be total")
	}

	assertFailClosedRefusal(t, openErr, rel)

	openat2PlatformGate := runtime.GOOS == "linux" && runtime.GOARCH == "amd64"
	if !openat2PlatformGate {
		// Ruling (a) on a build that is not linux/amd64: the gate is a constant
		// comparison, so the refusal is produced without any descriptor being touched.
		// Clearing the latch isolates that branch from the ENOSYS latch above; dirfd -1
		// can never reach the kernel because the platform check precedes the syscall.
		openat2Disabled.Store(false)
		fd, err := openBoundArtifactLocked(-1, rel)
		if fd != -1 || !errors.Is(err, syscall.ENOSYS) {
			t.Fatalf("a %s/%s build must answer ENOSYS from the platform gate before touching the descriptor, got fd %d err %v",
				runtime.GOOS, runtime.GOARCH, fd, err)
		}
		openat2Disabled.Store(true)
	}

	// The latch keeps the refusal stable: a second read answers the same way instead
	// of retrying a syscall that cannot succeed.
	if _, err := OpenBoundDirFile(dir, rel, 1024); err == nil {
		t.Fatal("the latched fail-closed state must refuse the second read too")
	} else {
		assertFailClosedRefusal(t, err, rel)
	}
}

// assertFailClosedRefusal pins the shape of the ENOSYS refusal: one *os.PathError
// from the bound open, naming the cleaned artifact, with ENOSYS inside.
func assertFailClosedRefusal(t *testing.T, err error, wantPath string) {
	t.Helper()
	if err == nil {
		t.Fatalf("the fail-closed gate must refuse %q with ENOSYS, got no error", wantPath)
	}
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("the fail-closed refusal must be an *os.PathError (the unchanged seam contract), got %T: %v", err, err)
	}
	if pathErr.Op != "openat" {
		t.Fatalf("the fail-closed refusal must come from the bound open (op %q), got op %q: %v", "openat", pathErr.Op, err)
	}
	if pathErr.Path != wantPath {
		t.Fatalf("the fail-closed refusal must name the cleaned artifact %q, got %q", wantPath, pathErr.Path)
	}
	if !errors.Is(pathErr.Err, syscall.ENOSYS) {
		t.Fatalf("the fail-closed refusal must be ENOSYS (ruling (a): no walk fallback, no allowUnsafeFallback here), got %v (%q)", pathErr.Err, pathErr.Err)
	}
}

// TestBoundDirConcurrentCloseAndRead is design r2 section 5 item F4.
//
// Readers keep opening and reading an in-bound artifact through one shared
// BoundDir while the test closes that directory repeatedly. Every outcome has to be
// one of exactly two things: the byte-exact in-bound canary (the descriptor was
// still valid, and the artifact handle owns its own file), or the
// nil-bound-directory refusal (the descriptor was already invalidated). A panic, an
// EBADF leak - the released number reaching the kernel - or any other byte string is
// a failure. Run it under -race to cover the dir.fd / dir.mu pairing.
func TestBoundDirConcurrentCloseAndRead(t *testing.T) {
	gap77SkipUnlessOpenat2(t)

	tree := buildGAP77Tree(t)
	dir, err := OpenBoundDir(tree.boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	const (
		readers = 4
		closes  = 16
		// in-bound artifact of the shared payload, one segment, exact bytes below.
		rel = "original"
	)

	var (
		mu       sync.Mutex
		reads    int
		refusals int
		problems []string
	)
	record := func(format string, args ...any) {
		mu.Lock()
		defer mu.Unlock()
		if len(problems) < 8 {
			problems = append(problems, fmt.Sprintf(format, args...))
		}
	}

	ready := make(chan struct{})
	var readyOnce sync.Once
	stop := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}

				file, openErr := OpenBoundDirFile(dir, rel, 1024)
				if openErr != nil {
					mu.Lock()
					refusals++
					mu.Unlock()
					if openErr.Error() != gap77NilBoundDirMessage {
						record("a read racing Close must be refused with %q, got %v: EBADF and every other errno must not leak out of a closed bound directory",
							gap77NilBoundDirMessage, openErr)
					}
					continue
				}

				data, truncated, readErr := file.ReadAll()
				closeErr := file.Close()
				mu.Lock()
				reads++
				mu.Unlock()
				readyOnce.Do(func() { close(ready) })

				if readErr != nil {
					record("a handle opened before Close must stay readable afterwards (the artifact descriptor is its own file), got %v", readErr)
					continue
				}
				if truncated {
					record("the %d byte in-bound artifact must not be truncated at the 1024 byte bound", len(gap77InsideCanary))
				}
				if string(data) != gap77InsideCanary {
					record("a read under a concurrent Close returned %d bytes %q, want the in-bound canary %q: outside bytes must never surface through a closed bound directory",
						len(data), data, gap77InsideCanary)
				}
				if closeErr != nil {
					record("closing the artifact handle failed: %v", closeErr)
				}
			}
		}()
	}

	select {
	case <-ready:
	case <-time.After(10 * time.Second):
		close(stop)
		wg.Wait()
		t.Fatal("no reader completed a read before Close could win the race: the concurrency case proved nothing")
	}

	// Close is idempotent and serialized by dir.mu, so the extra calls are no-ops
	// that keep the race window open while the readers are still looping.
	for i := 0; i < closes; i++ {
		_ = dir.Close()
	}
	close(stop)
	wg.Wait()

	mu.Lock()
	gotProblems := problems
	gotReads, gotRefusals := reads, refusals
	mu.Unlock()
	if len(gotProblems) > 0 {
		t.Fatalf("%s", strings.Join(gotProblems, "; "))
	}

	// Deterministic terminal state: whichever Close ran last, the directory is closed
	// and the seam says so with its own refusal rather than EBADF or a handle.
	switch _, err := OpenBoundDirFile(dir, rel, 1024); {
	case err == nil:
		t.Fatal("the bound directory is still usable after the Close storm: Close must invalidate the descriptor")
	case err.Error() != gap77NilBoundDirMessage:
		t.Fatalf("a closed bound directory must be refused with %q, got %v", gap77NilBoundDirMessage, err)
	}

	if gotRefusals == 0 {
		// A legitimate schedule: every reader completed its read before the first
		// Close. Not a failure - the byte-exact reads and the terminal state above
		// still hold - but it is worth recording.
		t.Logf("F4: no reader observed the closed state (all %d reads completed before the close), which is a valid schedule", gotReads)
	}
	t.Logf("F4: %d byte-exact reads, %d refusals with %q, %d Close calls", gotReads, gotRefusals, gap77NilBoundDirMessage, closes)
}
