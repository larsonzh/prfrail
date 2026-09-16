//go:build a7native && windows

package adapters

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"unsafe"
)

// These probes test the Windows primitives a publication protocol could use for
// the step the current implementation leaves as a no-op. They record what the
// platform actually does - the errors included - because "cannot be done without
// privileges" and "is documented to do it" are different claims, and only the
// first one is observable here.

var (
	a7Kernel32    = syscall.NewLazyDLL("kernel32.dll")
	a7MoveFileExW = a7Kernel32.NewProc("MoveFileExW")
)

// a7MoveFileWriteThrough is MOVEFILE_WRITE_THROUGH from WinBase.h. It is spelled
// out because syscall does not export MoveFileEx at all: the API is reached
// through kernel32 so the probe stays on the standard library.
const a7MoveFileWriteThrough = 0x00000008

func a7MoveFileExWriteThrough(from, to string) error {
	fromPointer := a7UTF16Ptr(from)
	toPointer := a7UTF16Ptr(to)
	result, _, callErr := a7MoveFileExW.Call(
		uintptr(unsafe.Pointer(fromPointer)),
		uintptr(unsafe.Pointer(toPointer)),
		uintptr(a7MoveFileWriteThrough),
	)
	if result == 0 {
		if callErr == syscall.Errno(0) {
			return errors.New("MoveFileExW failed without reporting an error")
		}
		return callErr
	}
	return nil
}

func a7UTF16Ptr(path string) *uint16 {
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}
	return pointer
}

// TestA7NativeDirectoryHandleFlushProbe records whether a directory can be opened
// with the access a flush needs. The result overturned the assumption this slice
// started from: only the read-only attempt is denied, and a directory opened with
// GENERIC_READ|GENERIC_WRITE plus FILE_FLAG_BACKUP_SEMANTICS accepts
// FlushFileBuffers without privileges. That success is an observation about the
// API, not about stable storage - whether NTFS then guarantees the directory entry
// survives a power loss is exactly what B3 has to inject.
func TestA7NativeDirectoryHandleFlushProbe(t *testing.T) {
	if os.Getenv(a7StageEnv) != "" {
		a7RunPublishHelper(t)
		return
	}
	directory := t.TempDir()
	shareMode := uint32(syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE | syscall.FILE_SHARE_DELETE)
	attempts := []struct {
		name   string
		access uint32
	}{
		{"read-only-backup-semantics", syscall.GENERIC_READ},
		{"read-write-backup-semantics", syscall.GENERIC_READ | syscall.GENERIC_WRITE},
		{"write-only-backup-semantics", syscall.GENERIC_WRITE},
	}
	readOnlyDenied := false
	flushAccepted := 0
	// The expected writable-mode count is fixed up front: counting only the modes
	// that happened to open would let a writable mode fail to open and still leave
	// the assertion green, which is exactly the over-claim this probe must prevent.
	const a7WritableModes = 2
	writableFlushed := 0
	for _, attempt := range attempts {
		writable := attempt.access&syscall.GENERIC_WRITE != 0
		handle, err := syscall.CreateFile(
			a7UTF16Ptr(directory),
			attempt.access,
			shareMode,
			nil,
			syscall.OPEN_EXISTING,
			syscall.FILE_FLAG_BACKUP_SEMANTICS,
			0,
		)
		if err != nil {
			t.Logf("A7 directory handle probe %s: CreateFile denied: %v", attempt.name, err)
			if writable {
				t.Fatalf("the writable access mode %s must be openable on a directory; if Windows changed that, the ADR-013 finding and the report must be revisited: %v", attempt.name, err)
			}
			readOnlyDenied = true
			continue
		}
		flushErr := syscall.FlushFileBuffers(handle)
		_ = syscall.CloseHandle(handle)
		t.Logf("A7 directory handle probe %s: FlushFileBuffers -> %v", attempt.name, flushErr)
		if flushErr == nil {
			flushAccepted++
			if writable {
				writableFlushed++
			}
		} else if !writable {
			readOnlyDenied = true
		}
	}
	if writableFlushed != a7WritableModes {
		// Asserting "at least one mode succeeded" would let the report claim that
		// both writable access modes work while only one of them ever did.
		t.Fatalf("every writable access mode must accept a directory flush, got %d of %d: the ADR-013 finding and its dependency on write access must be revisited", writableFlushed, a7WritableModes)
	}
	if flushAccepted == 0 {
		// This experiment runs on the Windows host only, never in CI, so a silent
		// log here would let the ADR-013 finding rot without any test going red.
		t.Fatal("no access mode accepted a directory flush: the ADR-013 finding and the observed dependency on write access must be revisited")
	}
	if !readOnlyDenied {
		t.Fatal("the read-only access mode must be denied; otherwise the observed dependency on write access is not established")
	}
	t.Logf("A7 directory handle probe FINDING: %d access mode(s) accepted FlushFileBuffers on a directory without privileges; whether NTFS makes the entry durable is unresolved and belongs to B3", flushAccepted)
}

// TestA7NativeMoveFileExWriteThroughProbe records the behaviour of the one
// candidate that needs no privileges. It must publish atomically, must refuse to
// replace an existing record with an error the convergence logic already
// understands, and must report whether an open reader changes anything.
func TestA7NativeMoveFileExWriteThroughProbe(t *testing.T) {
	if os.Getenv(a7StageEnv) != "" {
		a7RunPublishHelper(t)
		return
	}
	root := t.TempDir()
	requests := filepath.Join(root, "requests")
	if err := os.MkdirAll(requests, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(requests, "request."+a7RequestID+".jsonl")
	temp := filepath.Join(requests, ".request."+a7RequestID+".jsonl.1.tmp")
	if err := os.WriteFile(temp, []byte("{\"a7\":\"move-file-ex\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := a7MoveFileExWriteThrough(temp, target); err != nil {
		t.Fatalf("MoveFileEx(MOVEFILE_WRITE_THROUGH) failed: %v", err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Fatalf("the move must publish the record: %v", err)
	}
	t.Log("A7 MoveFileEx probe: MOVEFILE_WRITE_THROUGH published the record without privileges")

	// A second move onto the same path must fail as "exists" so the caller keeps
	// the no-replace convergence semantics it already has.
	second := filepath.Join(requests, ".request."+a7RequestID+".jsonl.2.tmp")
	if err := os.WriteFile(second, []byte("{\"a7\":\"second\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	replaceErr := a7MoveFileExWriteThrough(second, target)
	if replaceErr == nil {
		t.Fatal("a no-replace move must not overwrite an existing record")
	}
	if !errors.Is(replaceErr, fs.ErrExist) {
		t.Fatalf("a no-replace move must report an existing target as fs.ErrExist so the caller's convergence logic survives, got %v", replaceErr)
	}
	t.Log("A7 MoveFileEx probe: existing-target error maps to fs.ErrExist (asserted, not just observed)")

	// Isolate the variable the first version of this probe conflated: bind a
	// no-replace move onto a target that already exists but is NOT open. If that
	// also fails, the refusal is about the target existing, not about an open
	// handle, and the observation must be worded accordingly.
	closedTarget := filepath.Join(requests, "request.closed.jsonl")
	if err := os.WriteFile(closedTarget, []byte("{\"a7\":\"closed\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	closedMove := filepath.Join(requests, ".request.closed.jsonl.4.tmp")
	if err := os.WriteFile(closedMove, []byte("{\"a7\":\"fourth\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	closedErr := a7MoveFileExWriteThrough(closedMove, closedTarget)
	t.Logf("A7 MoveFileEx probe: publishing onto an existing target with no open handle -> %v", closedErr)
	if closedErr == nil {
		t.Fatal("a no-replace move must refuse an existing target even when nobody holds it open")
	}
	// The control case must land on the same error class as the open-reader case,
	// otherwise "the open handle made no difference" is only a guess.
	if !errors.Is(closedErr, fs.ErrExist) {
		t.Fatalf("an existing target must map to fs.ErrExist even with no open handle, got %v", closedErr)
	}

	// An open reader must not turn a publication into data loss: record whether the
	// primitive refuses or succeeds while a handle is open.
	openTarget := filepath.Join(requests, "request.reader.jsonl")
	if err := os.WriteFile(openTarget, []byte("{\"a7\":\"reader\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	reader, err := os.Open(openTarget)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	third := filepath.Join(requests, ".request.reader.jsonl.3.tmp")
	if err := os.WriteFile(third, []byte("{\"a7\":\"third\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	readerErr := a7MoveFileExWriteThrough(third, openTarget)
	t.Logf("A7 MoveFileEx probe: publishing onto a path with an open reader -> %v", readerErr)
	if !errors.Is(readerErr, fs.ErrExist) {
		t.Fatalf("the refusal to overwrite must surface as fs.ErrExist in this case too, got %v", readerErr)
	}
	_ = os.Remove(third)
	t.Log("A7 MoveFileEx probe: the refusal is not isolated to the open-reader case - an existing target is already enough, so the report must not claim a cause it did not measure")
}
