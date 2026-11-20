//go:build windows

package guard

// P83 FIX ROUND 7 counterexamples (items T3a, T3b and T3c) on the Windows half.
//
// The three tests attack the binding from the outside: T3a proves the verified
// handle keeps the image unwritable and unrenamable for the whole verify -> spawn
// window, T3b proves a bare command name is refused before any PATH lookup, T3c
// proves a mismatching approval refuses before a child exists and leaves no
// handle behind.
//
// Mutation contracts are stated on each test.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// p83Fix7ApprovedImage writes one approved image and returns its path and the
// digest of the exact bytes it wrote.
func p83Fix7ApprovedImage(t *testing.T) (string, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "agent-runner.exe")
	content := []byte("p83-fix-round-7 approved image")
	if err := os.WriteFile(path, content, 0o700); err != nil {
		t.Fatal(err)
	}
	return path, evidence.Digest(p83Fix7ExecutableDomain, content)
}

// TestP83Fix7T3aHeldApprovalDeniesWriteAndReplace is the T3a counterexample: while
// the guard holds the read-only-share handle of the verified executable, no other
// opener can write the image in place or replace it by a rename, so the object that
// was hashed is still the object CreateProcessW resolves.
//
// Mutation check: widening openApprovedExecutable to FILE_SHARE_READ|WRITE|DELETE
// lets both the write and the rename below succeed and turns this red.
func TestP83Fix7T3aHeldApprovalDeniesWriteAndReplace(t *testing.T) {
	path, digest := p83Fix7ApprovedImage(t)
	approved, err := verifyApprovedExecutable(path, &ProcessSpec{
		ExecutableDigest:       digest,
		ExecutableDigestDomain: p83Fix7ExecutableDomain,
	})
	if err != nil {
		t.Fatalf("the approved image must verify, got %v", err)
	}
	defer func() { _ = approved.Close() }()

	for _, flag := range []int{os.O_WRONLY, os.O_RDWR} {
		file, openErr := os.OpenFile(path, flag, 0)
		if openErr == nil {
			_ = file.Close()
			t.Fatalf("opening the approved image for write (flag %#o) must fail while the verified handle is held", flag)
		}
	}
	if err := os.WriteFile(path, []byte("rewritten in place"), 0o700); err == nil {
		t.Fatal("an in-place rewrite of the approved image must fail while the verified handle is held")
	}
	if err := os.Rename(path, path+".swapped"); err == nil {
		t.Fatal("replacing the approved image by a rename must fail while the verified handle is held")
	}
	// The verified object is unchanged: the refusals above were not a lucky ACL.
	if content, readErr := os.ReadFile(path); readErr != nil || string(content) != "p83-fix-round-7 approved image" {
		t.Fatalf("the approved image must still carry the approved bytes, got %q (err %v)", content, readErr)
	}

	// Control: once the verified handle is released the same write is allowed, so
	// the denial above was the share mode of the held handle and nothing else.
	if err := approved.Close(); err != nil {
		t.Fatalf("closing the verified handle: %v", err)
	}
	if err := os.WriteFile(path, []byte("after release"), 0o700); err != nil {
		t.Fatalf("the release control must be able to write the image, got %v", err)
	}
}

// TestP83Fix7T3bBareApprovedCommandIsRefused is the T3b counterexample: an
// approval names one object, and a bare command name would be resolved through
// PATH at spawn time, so it is refused before any verification or spawn.
//
// Mutation check: deleting the `strings.ContainsAny(spec.Command, ...)` refusal in
// platformProcess.start lets the bare name reach openApprovedExecutable; the error
// then reports the failed open instead of the path requirement and turns this red.
func TestP83Fix7T3bBareApprovedCommandIsRefused(t *testing.T) {
	_, digest := p83Fix7ApprovedImage(t)
	process, err := preparePlatformProcess()
	if err != nil {
		t.Fatal(err)
	}
	defer process.close()

	_, err = process.start(&ProcessSpec{
		Command:                "agent-runner",
		ExecutableDigest:       digest,
		ExecutableDigestDomain: p83Fix7ExecutableDomain,
	})
	if !errors.Is(err, ErrExecutableBinding) {
		t.Fatalf("a bare approved command name must refuse with the binding error, got %v", err)
	}
	if !strings.Contains(err.Error(), "named by path") {
		t.Fatalf("the refusal must report the path requirement, got %v", err)
	}

	// Control: the same digest with a path goes PAST the name check and fails only
	// because the named file does not exist, so the assertion above is specific to
	// the bare name and not to any refusal of the approved path.
	_, err = process.start(&ProcessSpec{
		Command:                filepath.Join(t.TempDir(), "missing.exe"),
		ExecutableDigest:       digest,
		ExecutableDigestDomain: p83Fix7ExecutableDomain,
	})
	if !errors.Is(err, ErrExecutableBinding) {
		t.Fatalf("a missing approved path must refuse with the binding error, got %v", err)
	}
	if strings.Contains(err.Error(), "named by path") {
		t.Fatalf("a path must not be refused for being a bare name, got %v", err)
	}
}

// TestP83Fix7T3cWrongApprovedDigestIsRefusedBeforeSpawn is the T3c counterexample:
// a digest that does not describe the opened object refuses, no child is created
// and no handle is retained.
//
// Mutation check: skipping the `actual != spec.ExecutableDigest` comparison makes
// verifyApprovedExecutable return a bound object and turns this red.
func TestP83Fix7T3cWrongApprovedDigestIsRefusedBeforeSpawn(t *testing.T) {
	path, _ := p83Fix7ApprovedImage(t)
	wrong := evidence.Digest(p83Fix7ExecutableDomain, []byte("an image the approval never saw"))

	_, err := verifyApprovedExecutable(path, &ProcessSpec{
		ExecutableDigest:       wrong,
		ExecutableDigestDomain: p83Fix7ExecutableDomain,
	})
	if !errors.Is(err, ErrExecutableBinding) {
		t.Fatalf("a wrong approved digest must refuse with the binding error, got %v", err)
	}
	// The refusal path must release the handle it opened: with the deny-delete share
	// mode still held the rename below would fail.
	if renameErr := os.Rename(path, path+".moved"); renameErr != nil {
		t.Fatalf("a refused verification must not keep the executable handle open, got %v", renameErr)
	}
	// No child: the refusal happened before any spawn, so the process identity the
	// guard would have created does not exist. The assertion is structural - the
	// function has no spawn call - and the counterexample above is what proves the
	// refusal precedes it.
}
