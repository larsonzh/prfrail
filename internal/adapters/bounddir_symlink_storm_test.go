//go:build linux && amd64

package adapters

// P53 FIX-GAP-77-SYMLINK-ESCAPE: the two red-team counterexamples that the fix's
// main path has to close, written before any production change (design r2 §1, §5;
// user rulings (a)~(e) of 2026-11-04).
//
//   - F1: a concurrent rename storm around the bound tree. Ruling (a) put external
//     rename on the threat boundary, so the guarantee "a successful read yields
//     in-bound bytes only" has to survive an attacker renaming entries of the bound
//     tree while a read is in flight. The storm writes a unique outside canary on
//     every iteration and asserts that no successful read ever returns it. This is
//     the regression guardrail: it passes at HEAD (a single openat is atomic) and it
//     must keep passing once the multi-step resolution exists.
//   - F2: an imported bind mount (mount --bind O D/mnt). A mount point is not a
//     symlink, so O_NOFOLLOW and O_DIRECTORY both let it through; ruling (b) put
//     mount imports on the attack surface, which is why the fix keeps
//     RESOLVE_NO_XDEV and why the refusal must be EXDEV.
//
// These tests are linux-only: they use Linux syscalls (unshare, mount) and the
// errno/namespace semantics of design r2 §1. The mount case runs in a re-exec'd
// child process so the user/mount namespace work never leaks into the parent test
// process.
//
// The build tag narrows linux to linux/amd64, the only target whose resolution path
// is the openat2 main path (design r2 §2); the tests also call
// gap77SkipUnlessOpenat2, which probes the kernel. Both matter: on linux/arm64 the
// production gate answers ENOSYS for every bound read, so "at least one read
// succeeds" and "the mount crossing is refused with EXDEV" would be guaranteed
// false reds instead of findings.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
)

const (
	// bindMountHelperEnv marks the re-exec'd child that owns the namespace work.
	bindMountHelperEnv = "PRFRAIL_P53_BIND_MOUNT_HELPER"
	// bindMountRootEnv carries the prepared tree from the parent to the child.
	bindMountRootEnv = "PRFRAIL_P53_BIND_MOUNT_ROOT"
	// bindMountUnsupportedMarker is printed by the child when the host refuses the
	// namespace or mount setup; the parent then skips with that reason on record.
	bindMountUnsupportedMarker = "BIND-SETUP-UNSUPPORTED: "
)

// TestBoundDirRenameStormNeverYieldsOutsideBytes is red-team finding F1.
//
// The bound tree is D/sub/nested/f = gap77InsideCanary. Every storm iteration
// rewrites the outside secret O/secret with a unique canary (the bytes an escape
// would surface) and renames D/sub out of the bound and back in, while two readers
// keep opening "sub/nested/f" through the bound handle.
//
// Assertions: every successful read equals the in-bound canary exactly; at least
// one read succeeds (otherwise the test would prove nothing); and the tree is
// intact afterwards. Refusals are tolerated - a name that is momentarily absent is
// a correct outcome, silence is not.
func TestBoundDirRenameStormNeverYieldsOutsideBytes(t *testing.T) {
	gap77SkipUnlessOpenat2(t)

	const (
		iterations = 200
		readers    = 2
		boundLimit = 1024
	)

	root := t.TempDir()
	boundPath := filepath.Join(root, "D")
	outsidePath := filepath.Join(root, "O")
	if err := os.Mkdir(boundPath, 0o755); err != nil {
		t.Fatalf("create the bound directory: %v", err)
	}
	if err := os.Mkdir(outsidePath, 0o755); err != nil {
		t.Fatalf("create the outside directory: %v", err)
	}
	insideDir := filepath.Join(boundPath, "sub", "nested")
	if err := os.MkdirAll(insideDir, 0o755); err != nil {
		t.Fatalf("create the in-bound multi-level directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(insideDir, "f"), []byte(gap77InsideCanary), 0o644); err != nil {
		t.Fatalf("write the in-bound artifact: %v", err)
	}
	outsideSecret := filepath.Join(outsidePath, "secret")
	if err := os.WriteFile(outsideSecret, []byte("OUTSIDE-CANARY-INIT"), 0o644); err != nil {
		t.Fatalf("write the outside secret: %v", err)
	}

	dir, err := OpenBoundDir(boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	rel := filepath.Join("sub", "nested", "f")

	var (
		mu         sync.Mutex
		successes  int
		refusals   int
		violations []string
	)

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

				file, openErr := OpenBoundDirFile(dir, rel, boundLimit)
				if openErr != nil {
					mu.Lock()
					refusals++
					mu.Unlock()
					continue
				}
				data, _, readErr := file.ReadAll()
				_ = file.Close()

				mu.Lock()
				if readErr != nil {
					refusals++
				} else {
					successes++
					if string(data) != gap77InsideCanary && len(violations) < 8 {
						violations = append(violations, string(data))
					}
				}
				mu.Unlock()
			}
		}()
	}

	// Deterministic storm: each iteration parks the in-bound subtree outside the
	// bound and moves it straight back, so a reader that opens while the name is
	// absent is refused and one that opens while it is present must see the
	// in-bound bytes. The outside secret is rewritten first so that any leakage
	// carries the canary of that very iteration.
	parked := filepath.Join(outsidePath, "parked")
	subPath := filepath.Join(boundPath, "sub")
	var stormErr error
	for i := 0; i < iterations; i++ {
		canary := fmt.Sprintf("OUTSIDE-CANARY-%08d", i)
		if err := os.WriteFile(outsideSecret, []byte(canary), 0o644); err != nil {
			stormErr = fmt.Errorf("rewrite the outside canary on iteration %d: %w", i, err)
			break
		}
		if err := os.Rename(subPath, parked); err != nil {
			stormErr = fmt.Errorf("park the in-bound subtree on iteration %d: %w", i, err)
			break
		}
		if err := os.Rename(parked, subPath); err != nil {
			stormErr = fmt.Errorf("restore the in-bound subtree on iteration %d: %w", i, err)
			break
		}
	}

	close(stop)
	wg.Wait()

	if stormErr != nil {
		t.Fatalf("rename storm failed: %v", stormErr)
	}
	t.Logf("F1 storm: %d iterations, %d successful reads, %d refusals", iterations, successes, refusals)
	if violations != nil {
		t.Fatalf("F1: a read under the rename storm returned bytes that are not the in-bound canary: %q (want %q); outside bytes must never surface through the bound handle",
			violations, gap77InsideCanary)
	}
	if successes == 0 {
		t.Fatalf("F1: the storm never produced a successful read (%d refusals), so the guardrail proved nothing", refusals)
	}

	// The storm must leave the payload intact: a damaged tree would let a later
	// assertion pass for the wrong reason.
	file, err := OpenBoundDirFile(dir, rel, boundLimit)
	if err != nil {
		t.Fatalf("the in-bound artifact must still open after the storm: %v", err)
	}
	defer func() { _ = file.Close() }()
	data, _, err := file.ReadAll()
	if err != nil {
		t.Fatalf("the in-bound artifact must still read after the storm: %v", err)
	}
	if string(data) != gap77InsideCanary {
		t.Fatalf("the in-bound artifact changed during the storm: got %q want %q", data, gap77InsideCanary)
	}
}

// TestBoundDirBindMountImportRefused is red-team finding F2: an imported bind
// mount inside the bound directory must not become a read path to the outside.
//
// The namespace work is delegated to a child process (see runBindMountImportHelper)
// so a failure to obtain the namespace leaves the parent test process untouched.
// When the host refuses the setup the case is skipped and the reason is printed, so
// the gap is registered instead of silently passing.
func TestBoundDirBindMountImportRefused(t *testing.T) {
	if os.Getenv(bindMountHelperEnv) == "1" {
		runBindMountImportHelper(t)
		return
	}
	// Gate the parent, not the child: a child that skipped here would still exit 0
	// and the parent would read that as a pass, silently losing the case.
	gap77SkipUnlessOpenat2(t)

	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "D"), 0o755); err != nil {
		t.Fatalf("create the bound directory: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "O"), 0o755); err != nil {
		t.Fatalf("create the outside directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "O", "secret"), []byte(gap77OutsideCanary), 0o644); err != nil {
		t.Fatalf("write the outside secret: %v", err)
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestBoundDirBindMountImportRefused$", "-test.v", "-test.count=1")
	cmd.Env = append(os.Environ(), bindMountHelperEnv+"=1", bindMountRootEnv+"="+root)
	output, err := cmd.CombinedOutput()
	text := string(output)

	if idx := strings.Index(text, bindMountUnsupportedMarker); idx >= 0 {
		reason := text[idx+len(bindMountUnsupportedMarker):]
		if nl := strings.IndexByte(reason, '\n'); nl >= 0 {
			reason = reason[:nl]
		}
		t.Skipf("bind mount import case skipped: host refused the setup (%s)", strings.TrimSpace(reason))
	}
	if err != nil {
		t.Fatalf("bind mount import case failed in the isolated child process: %v\n%s", err, text)
	}
}

// runBindMountImportHelper is the child body: it moves itself into a new user and
// mount namespace, binds O onto D/mnt, and then asks the bound open for
// "mnt/secret". After the P53 fix the open must be refused with EXDEV; at HEAD the
// helper fails with the outside bytes it read, which is the recorded red.
func runBindMountImportHelper(t *testing.T) {
	root := os.Getenv(bindMountRootEnv)
	if root == "" {
		t.Fatalf("helper started without %s", bindMountRootEnv)
	}
	boundPath := filepath.Join(root, "D")
	outsidePath := filepath.Join(root, "O")

	// Capture the identity before the namespace switch: inside a fresh user
	// namespace the same uid has no meaning until the map is written.
	uid, gid := os.Getuid(), os.Getgid()

	// Preferred shape: a user namespace (which grants CAP_SYS_ADMIN inside it) plus a
	// mount namespace. Hosts that disable unprivileged user namespaces fail here with
	// EPERM/EINVAL, but a test binary run as root needs only the mount namespace, so
	// that fallback is tried before the case is reported as unsupported.
	if err := syscall.Unshare(syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS); err != nil {
		if !errors.Is(err, syscall.EPERM) && !errors.Is(err, syscall.EINVAL) {
			unsupported(t, "unshare(CLONE_NEWUSER|CLONE_NEWNS)", err)
			return
		}
		if mountErr := syscall.Unshare(syscall.CLONE_NEWNS); mountErr != nil {
			unsupported(t, "unshare(CLONE_NEWUSER|CLONE_NEWNS), then unshare(CLONE_NEWNS) as a privileged fallback",
				fmt.Errorf("user namespace: %v; mount namespace: %v", err, mountErr))
			return
		}
	} else if err := mapSelfIDs(uid, gid); err != nil {
		unsupported(t, "write /proc/self/{uid,gid}_map", err)
		return
	}

	mountTarget := filepath.Join(boundPath, "mnt")
	if err := os.Mkdir(mountTarget, 0o755); err != nil {
		t.Fatalf("create the mount point inside the bound directory: %v", err)
	}
	if err := syscall.Mount(outsidePath, mountTarget, "", syscall.MS_BIND, ""); err != nil {
		unsupported(t, "mount --bind", err)
		return
	}
	defer func() { _ = syscall.Unmount(mountTarget, syscall.MNT_DETACH) }()

	dir, err := OpenBoundDir(boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	relative := filepath.Join("mnt", "secret")
	file, err := OpenBoundDirFile(dir, relative, 1024)
	if err == nil {
		data, _, readErr := file.ReadAll()
		_ = file.Close()
		t.Fatalf("F2: the imported bind mount %q was followed and read %q (read error %v); RESOLVE_NO_XDEV must refuse the mount crossing with EXDEV",
			relative, data, readErr)
	}

	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Fatalf("the mount import refusal must be an *os.PathError, got %T: %v", err, err)
	}
	if !errors.Is(pathErr.Err, syscall.EXDEV) {
		t.Fatalf("the mount import refusal must be EXDEV (RESOLVE_NO_XDEV), got %v (%q)", pathErr.Err, pathErr.Err)
	}
}

// mapSelfIDs maps the pre-namespace identity to root inside the new user
// namespace, which is what lets the child hold CAP_SYS_ADMIN there and create the
// mount point in the bound directory.
func mapSelfIDs(uid, gid int) error {
	if err := os.WriteFile("/proc/self/setgroups", []byte("deny"), 0o644); err != nil {
		return fmt.Errorf("setgroups: %w", err)
	}
	if err := os.WriteFile("/proc/self/uid_map", fmt.Appendf(nil, "0 %d 1", uid), 0o644); err != nil {
		return fmt.Errorf("uid_map: %w", err)
	}
	if err := os.WriteFile("/proc/self/gid_map", fmt.Appendf(nil, "0 %d 1", gid), 0o644); err != nil {
		return fmt.Errorf("gid_map: %w", err)
	}
	return nil
}

// unsupported records a host that refuses the namespace or the mount, which is an
// environment limit (no CAP_SYS_ADMIN / userns restricted), not a test failure.
func unsupported(t *testing.T, step string, err error) {
	t.Helper()
	fmt.Printf("%s%s: %v\n", bindMountUnsupportedMarker, step, err)
	t.Skipf("F2 bind mount import case cannot run here: %s: %v", step, err)
}
