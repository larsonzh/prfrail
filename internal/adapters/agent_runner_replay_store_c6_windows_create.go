//go:build windows && c6slot

package adapters

// Slot creation chain of the C6 prebuilt-slot implementation
// (agent_runner_replay_store_c6_windows.go). It is compiled only for
// windows && c6slot, where the activation prebuild path consumes it.

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

const replaySlotMoveFileWriteThrough = 0x00000008

var (
	replaySlotErrorFileExists    = syscall.Errno(80)
	replaySlotErrorAlreadyExists = syscall.Errno(183)

	replaySlotProcMoveFileExW = replaySlotKernel32.NewProc("MoveFileExW")
)

// slotCreateFn is the creation seam of the C6 slot pipeline. It is nil in
// production; the seam contract is documented with the other private seams in
// agent_runner_replay_store_c6_windows.go.
var slotCreateFn func(directory, pattern string) (*os.File, error)

// createReplaySlotWindows creates one empty slot: a temporary file in the target
// directory holding the 16-byte zero header, flushed, then published under the
// final name with MoveFileExW(MOVEFILE_WRITE_THROUGH) and without
// MOVEFILE_REPLACE_EXISTING. An existing target therefore surfaces as
// fs.ErrExist, the no-replace contract the store already relies on.
func createReplaySlotWindows(path string) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return err
	}
	file, err := createSlotFile(directory, filepath.Base(path))
	if err != nil {
		return err
	}
	temp := file.Name()
	header := replaySlotHeader(0)
	_, writeErr := file.Write(header[:])
	if writeErr == nil {
		writeErr = flushSlotFile(file)
	}
	if closeErr := file.Close(); writeErr == nil {
		writeErr = closeErr
	}
	if writeErr != nil {
		_ = os.Remove(temp)
		return writeErr
	}
	if err := moveReplaySlotWriteThrough(temp, path); err != nil {
		_ = os.Remove(temp)
		return err
	}
	return nil
}

// moveReplaySlotWriteThrough publishes from as to without replacing an existing
// target. syscall does not export MoveFileEx at all, so the call goes through
// kernel32 to stay on the standard library.
func moveReplaySlotWriteThrough(from, to string) error {
	fromPointer, err := syscall.UTF16PtrFromString(from)
	if err != nil {
		return err
	}
	toPointer, err := syscall.UTF16PtrFromString(to)
	if err != nil {
		return err
	}
	result, _, callErr := replaySlotProcMoveFileExW.Call(
		uintptr(unsafe.Pointer(fromPointer)),
		uintptr(unsafe.Pointer(toPointer)),
		uintptr(replaySlotMoveFileWriteThrough),
	)
	if result == 0 {
		if callErr == syscall.Errno(0) {
			return errors.New("MoveFileExW failed without reporting an error")
		}
		if errors.Is(callErr, replaySlotErrorAlreadyExists) || errors.Is(callErr, replaySlotErrorFileExists) {
			return fs.ErrExist
		}
		return callErr
	}
	return nil
}

func createSlotFile(directory, base string) (*os.File, error) {
	if slotCreateFn != nil {
		return slotCreateFn(directory, base)
	}
	return os.CreateTemp(directory, "."+base+".*.tmp")
}
