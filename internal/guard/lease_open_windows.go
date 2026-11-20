//go:build windows

package guard

import (
	"fmt"
	"os"
	"syscall"
)

// leaseRecordOpenReparsePoint opens the record object itself instead of the target a
// reparse point redirects to.
const leaseRecordOpenReparsePoint = 0x00200000

// leaseRecordOpenBackupSemantics lets the open reach a DIRECTORY so the refusal can
// classify it as "not a regular file" instead of surfacing the share/access error a
// directory open otherwise produces. The record is never read through this flag:
// the type check on the handle below refuses it first.
const leaseRecordOpenBackupSemantics = 0x02000000

// openLeaseRecord opens one lease record for the bounded record read (P83 fix
// round 7, item T5#3; tightened in fix round 8 for the same item).
//
// The record is opened with FILE_FLAG_OPEN_REPARSE_POINT, so a symbolic link or a
// junction is opened AS ITSELF and never followed, and the opened object is then
// classified on the SAME handle: a reparse point is refused, and so is anything that
// is not a regular file. The previous form classified the NAME with os.Lstat and
// opened the name afterwards, so the record could be swapped for a link between the
// two calls and the read would then describe the link target. Nothing about the
// opened object is re-resolved by name here or by the caller: readLeaseRecord stat's
// the returned handle and bounds the read.
//
// A fully cancellable read is not reachable with the standard library: os.File on a
// synchronous Windows handle offers no per-read deadline, so the byte bound in
// readLeaseRecord stays the only limit on a hostile record. Windows exposes a
// relative open (NtCreateFile with RootDirectory) only outside the standard library,
// and this repository takes no third-party dependency; opening the reparse point as
// itself and classifying it on the same handle is what the standard library allows,
// and it removes the Lstat/Open window entirely.
func openLeaseRecord(path string) (*os.File, error) {
	name, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL|leaseRecordOpenReparsePoint|leaseRecordOpenBackupSemantics, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(handle, &info); err != nil {
		_ = syscall.CloseHandle(handle)
		return nil, &os.PathError{Op: "stat", Path: path, Err: err}
	}
	if info.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		_ = syscall.CloseHandle(handle)
		return nil, fmt.Errorf("lease record %s is a reparse point", path)
	}
	if info.FileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY != 0 {
		_ = syscall.CloseHandle(handle)
		return nil, fmt.Errorf("lease record %s is not a regular file", path)
	}
	return os.NewFile(uintptr(handle), path), nil
}
