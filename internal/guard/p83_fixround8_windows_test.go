//go:build windows

package guard

// P83 FIX ROUND 8 counterexamples (item T5#3) on the Windows half: the lease
// record is opened as its own object and classified on the SAME handle, so a
// reparse point at the record path is refused as a reparse point and a DIRECTORY
// there is refused as "not a regular file" instead of surfacing the platform's
// access-denied error.
//
// Mutation contracts are stated on each test.

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

// p83Fix8Junction creates one directory junction. An unprivileged account can
// create a junction, so the reparse branch is reachable without the
// symbolic-link privilege.
func p83Fix8Junction(t *testing.T, link, target string) {
	t.Helper()
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); err != nil {
		t.Fatalf("create the junction %s -> %s: %v (%s)", link, target, err, output)
	}
	// Runs before the TempDir cleanup (cleanup order is LIFO).
	t.Cleanup(func() { _ = os.Remove(link) })
}

// TestP83Fix8T5c3ReparseRecordIsRefused is the reparse half of T5#3: a record name
// that resolves to a reparse point refuses with the classification, not with the
// target the point redirects to.
//
// Mutation check: dropping leaseRecordOpenReparsePoint from openLeaseRecord makes
// the junction resolve to its target, which is then classified as a directory -
// the assertion on the reparse classification turns red.
func TestP83Fix8T5c3ReparseRecordIsRefused(t *testing.T) {
	root := t.TempDir()
	recordPath := filepath.Join(root, "workspace-p83-fix8-reparse.lease.json")
	p83Fix8Junction(t, recordPath, filepath.Join(root, "elsewhere"))

	_, err := readLeaseRecord(recordPath)
	if err == nil {
		t.Fatal("a reparse point at the record path must be refused")
	}
	if !strings.Contains(err.Error(), "reparse point") {
		t.Fatalf("the refusal must name the reparse point, got %v", err)
	}
}

// TestP83Fix8T5c3DirectoryRecordIsRefusedAsNonRegular is the directory half of
// T5#3: the record path holding a DIRECTORY is refused as "not a regular file",
// which is the classification the backup-semantics open exists for - not the
// platform's access-denied error.
//
// Mutation check: dropping leaseRecordOpenBackupSemantics makes CreateFile fail on
// the directory with ERROR_ACCESS_DENIED before any classification happens, so the
// error carries the platform message and the assertions turn red.
func TestP83Fix8T5c3DirectoryRecordIsRefusedAsNonRegular(t *testing.T) {
	root := t.TempDir()
	recordPath := filepath.Join(root, "workspace-p83-fix8-dir.lease.json")
	if err := os.Mkdir(recordPath, 0o700); err != nil {
		t.Fatal(err)
	}

	// What the platform answers when the directory is opened WITHOUT the backup
	// semantics - the error a reverted guard would surface.
	pointer, err := syscall.UTF16PtrFromString(recordPath)
	if err != nil {
		t.Fatal(err)
	}
	plainShare := uint32(syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE | syscall.FILE_SHARE_DELETE)
	plain, plainErr := syscall.CreateFile(pointer, syscall.GENERIC_READ, plainShare, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if plainErr == nil {
		_ = syscall.CloseHandle(plain)
		t.Fatal("precondition failed: a directory must not open without backup semantics")
	}
	if !errors.Is(plainErr, syscall.ERROR_ACCESS_DENIED) {
		t.Fatalf("precondition failed: expected an access-denied open, got %v", plainErr)
	}

	_, err = readLeaseRecord(recordPath)
	if err == nil {
		t.Fatal("a directory at the record path must be refused")
	}
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("the directory must be refused as non-regular, got %v", err)
	}
	if strings.Contains(err.Error(), plainErr.Error()) {
		t.Fatalf("the refusal must not be the platform access-denied error, got %v", err)
	}
}
