//go:build windows

package adapters

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"unsafe"
)

// These are the two publication primitives A7 left as candidates for the step the
// Windows durability implementation does not yet perform: publishing a record with
// a write-through move (C2), and dropping the write access mode would be a silent
// downgrade, so the probe below pins the access mode as well.
//
// They live in test scope on purpose. Nothing here is wired into the store: the
// default primitive stays os.Link and the default parent-directory step stays the
// Windows no-op until a candidate is proven, and shipping an unused publication
// primitive into the product binary would be exactly the speculative interface the
// architecture rules forbid.
var (
	b3bKernel32    = syscall.NewLazyDLL("kernel32.dll")
	b3bMoveFileExW = b3bKernel32.NewProc("MoveFileExW")
)

// b3bMoveFileWriteThrough is MOVEFILE_WRITE_THROUGH from WinBase.h. syscall does
// not export MoveFileEx at all, so the call goes through kernel32 to stay on the
// standard library.
const b3bMoveFileWriteThrough = 0x00000008

// b3bMoveFileExWriteThrough publishes from as to without replacing an existing
// target: MOVEFILE_REPLACE_EXISTING is deliberately absent, and an existing target
// therefore has to surface as fs.ErrExist, the same no-replace contract os.Link
// satisfies today.
func b3bMoveFileExWriteThrough(from string, to string) error {
	fromPointer, err := syscall.UTF16PtrFromString(from)
	if err != nil {
		return err
	}
	toPointer, err := syscall.UTF16PtrFromString(to)
	if err != nil {
		return err
	}
	result, _, callErr := b3bMoveFileExW.Call(
		uintptr(unsafe.Pointer(fromPointer)),
		uintptr(unsafe.Pointer(toPointer)),
		uintptr(b3bMoveFileWriteThrough),
	)
	if result == 0 {
		if callErr == syscall.Errno(0) {
			return errors.New("MoveFileExW failed without reporting an error")
		}
		if errors.Is(callErr, syscall.ERROR_ALREADY_EXISTS) || errors.Is(callErr, syscall.ERROR_FILE_EXISTS) {
			return fs.ErrExist
		}
		return callErr
	}
	return nil
}

// b3bFlushDirectoryHandle opens directory for the access a directory flush needs
// and flushes it. The access mode is not decoration: A7 measured that a read-only
// directory handle is denied, so a primitive that silently relaxed the mode could
// not open the directory at all.
func b3bFlushDirectoryHandle(directory string) error {
	handle, err := syscall.CreateFile(
		mustB3bUTF16Ptr(directory),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		uint32(syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE),
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		return err
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	return syscall.FlushFileBuffers(handle)
}

// b3bVolumeDevicePath turns an absolute path with a drive letter into the volume
// device path FlushFileBuffers needs ("R:\agent-runner-replay\x.jsonl" -> `\\.\R:`).
// A path without a drive letter is refused rather than guessed: flushing the wrong
// device would turn the positive control into a silent no-op with a clean-looking
// result, which is the one failure mode the control exists to rule out.
func b3bVolumeDevicePath(path string) (string, error) {
	if len(path) < 3 || path[1] != ':' || (path[2] != '\\' && path[2] != '/') {
		return "", errors.New("b3bVolumeDevicePath: not an absolute drive path: " + path)
	}
	return `\\.\` + strings.ToUpper(path[:2]), nil
}

// b3bFlushDeviceHandle flushes the write cache of the volume that holds path. This is
// the positive control of the B3b matrix: it is deliberately not a publication
// primitive the product could adopt (it needs a volume handle and write access to it),
// and it exists to show that the *device* can be drained before a cut. Without that,
// "nothing was lost" could not be told apart from "this rig cannot make the device
// lose anything in the first place".
//
// The error is returned, never swallowed: an ACCESS_DENIED here means the control
// could not run at that privilege level, and the round has to say so instead of
// counting as a clean result.
func b3bFlushDeviceHandle(path string) error {
	device, err := b3bVolumeDevicePath(path)
	if err != nil {
		return err
	}
	handle, err := syscall.CreateFile(
		mustB3bUTF16Ptr(device),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		uint32(syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE),
		nil,
		syscall.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		return err
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	return syscall.FlushFileBuffers(handle)
}

// b3bFlushTargetFile opens the published record itself and flushes it. It is the
// control's fallback when the volume handle cannot be opened, and it reports under its
// own token so the trace never claims a device flush that did not happen.
func b3bFlushTargetFile(path string) error {
	handle, err := syscall.CreateFile(
		mustB3bUTF16Ptr(path),
		syscall.GENERIC_READ|syscall.GENERIC_WRITE,
		uint32(syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE),
		nil,
		syscall.OPEN_EXISTING,
		0,
		0,
	)
	if err != nil {
		return err
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	return syscall.FlushFileBuffers(handle)
}

// TestB3bCandidateDeviceFlush pins the volume-path derivation and requires the flush to
// report what happened. Whether the volume handle can be opened depends on the
// privileges of the caller, so the test accepts a refusal - but it has to be a
// *reported* refusal, which is what keeps the control honest on a non-elevated host.
func TestB3bCandidateDeviceFlush(t *testing.T) {
	accepted := []struct{ in, want string }{
		{`R:\agent-runner-replay\requests\request.b3b-r01.jsonl`, `\\.\R:`},
		{`c:/tmp/x.jsonl`, `\\.\C:`},
	}
	for _, testCase := range accepted {
		got, err := b3bVolumeDevicePath(testCase.in)
		if err != nil {
			t.Fatalf("b3bVolumeDevicePath(%q): %v", testCase.in, err)
		}
		if got != testCase.want {
			t.Fatalf("b3bVolumeDevicePath(%q) = %q, want %q", testCase.in, got, testCase.want)
		}
	}
	for _, refused := range []string{`relative\x.jsonl`, `R:`, `\\.\R:`} {
		if _, err := b3bVolumeDevicePath(refused); err == nil {
			t.Fatalf("b3bVolumeDevicePath(%q) must be refused", refused)
		}
	}
	target := filepath.Join(t.TempDir(), "record.jsonl")
	if err := os.WriteFile(target, []byte("payload"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := b3bFlushTargetFile(target); err != nil {
		t.Fatalf("flushing the published record must work for an ordinary user: %v", err)
	}
	if err := b3bFlushDeviceHandle(target); err != nil {
		t.Logf("volume handle not available at this privilege level (reported, not swallowed): %v", err)
	}
}

func mustB3bUTF16Ptr(path string) *uint16 {
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}
	return pointer
}

// TestB3bCandidateMoveFileExNoReplace pins the no-replace contract of the C2
// candidate. Removing MOVEFILE_REPLACE_EXISTING is what makes the second case
// fail: with the replace bit set the move would silently take over an owned slot.
func TestB3bCandidateMoveFileExNoReplace(t *testing.T) {
	directory := t.TempDir()
	source := filepath.Join(directory, "source.jsonl")
	target := filepath.Join(directory, "target.jsonl")
	payload := []byte("{\"record\":\"c2\"}\n")
	if err := os.WriteFile(source, payload, 0o600); err != nil {
		t.Fatal(err)
	}

	if err := b3bMoveFileExWriteThrough(source, target); err != nil {
		t.Fatalf("a write-through move onto an absent target must succeed: %v", err)
	}
	if _, err := os.Stat(source); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a successful move must consume the source, stat = %v", err)
	}
	moved, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(moved) != string(payload) {
		t.Fatalf("moved bytes = %q, want %q", moved, payload)
	}

	// An existing target must be refused, whether or not it is open.
	secondSource := filepath.Join(directory, "source-two.jsonl")
	if err := os.WriteFile(secondSource, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := b3bMoveFileExWriteThrough(secondSource, target); !errors.Is(err, fs.ErrExist) {
		t.Fatalf("an existing target must surface as fs.ErrExist, got %v", err)
	}
	if _, err := os.Stat(secondSource); err != nil {
		t.Fatalf("a refused move must leave the source in place: %v", err)
	}

	openTarget, err := os.Open(target)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = openTarget.Close() }()
	if err := b3bMoveFileExWriteThrough(secondSource, target); !errors.Is(err, fs.ErrExist) {
		t.Fatalf("an open target must still surface as fs.ErrExist, got %v", err)
	}

	thirdSource := filepath.Join(directory, "source-three.jsonl")
	missingTarget := filepath.Join(directory, "missing", "target.jsonl")
	if err := b3bMoveFileExWriteThrough(thirdSource, missingTarget); err == nil {
		t.Fatal("moving a source that does not exist must fail")
	} else if errors.Is(err, fs.ErrExist) {
		t.Fatalf("a missing source is not a no-replace collision: %v", err)
	}
}

// TestB3bCandidateDirectoryFlush pins that the C3 candidate really opens the
// directory with the access mode a flush requires, and that it reports failures
// instead of swallowing them. Relaxing the access mode to read-only, which reads
// as harmless, stops the directory from opening at all.
func TestB3bCandidateDirectoryFlush(t *testing.T) {
	directory := t.TempDir()
	if err := b3bFlushDirectoryHandle(directory); err != nil {
		t.Fatalf("flushing a directory through a write-access handle must succeed: %v", err)
	}
	missing := filepath.Join(directory, "missing")
	err := b3bFlushDirectoryHandle(missing)
	if err == nil {
		t.Fatal("flushing a directory that does not exist must fail")
	}
	if strings.Contains(err.Error(), "no-replace") {
		t.Fatalf("unexpected error shape: %v", err)
	}
}
