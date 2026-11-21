//go:build windows

package evidence

// P42-T027B-EVIDENCE-BIND Q2, Windows half: a directory handle opened with
// FILE_FLAG_BACKUP_SEMANTICS plus share-delete, and relative opens resolved
// through the object's current real name (GetFinalPathNameByHandle). Windows has
// no openat, so the handle-derived name is what makes the binding follow the
// directory object instead of the path.

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"unsafe"
)

const (
	// fileFlagBackupSemantics opens a directory handle instead of refusing it.
	fileFlagBackupSemantics = 0x02000000
	// fileInformationShareAll is FILE_SHARE_READ|WRITE|DELETE, the share mode that
	// keeps the bound object renamable and removable while it stays bound.
	fileInformationShareAll = syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE | syscall.FILE_SHARE_DELETE
	// fileFlagOpenReparsePoint opens a reparse point as itself instead of following
	// it to its target (T4, P83 fix round 8).
	fileFlagOpenReparsePoint = 0x00200000
	// volumeNameDos and fileNameNormalized ask for the ordinary DOS path form.
	volumeNameDos      = 0x00000000
	fileNameNormalized = 0x00000000
	// boundDirNameBufferStart is the first buffer size tried for the real name.
	boundDirNameBufferStart = 512
	// boundDirNameBufferLimit bounds the grow-and-retry loop.
	boundDirNameBufferLimit = 1 << 16
)

var procGetFinalPathNameByHandleW = syscall.NewLazyDLL("kernel32.dll").NewProc("GetFinalPathNameByHandleW")

// BoundDir is one open directory handle. It owns the handle until Close.
type BoundDir struct {
	handle syscall.Handle
	path   string
}

// OpenBoundDir opens path as a directory handle. The handle is share-delete, so
// the directory stays renamable and removable while it is bound, and the handle
// keeps addressing the directory object even after the name moves.
func OpenBoundDir(path string) (*BoundDir, error) {
	name, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, fileInformationShareAll, nil, syscall.OPEN_EXISTING, fileFlagBackupSemantics, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(handle, &info); err != nil {
		_ = syscall.CloseHandle(handle)
		return nil, &os.PathError{Op: "stat", Path: path, Err: err}
	}
	if info.FileAttributes&syscall.FILE_ATTRIBUTE_DIRECTORY == 0 {
		_ = syscall.CloseHandle(handle)
		return nil, fmt.Errorf("evidence: %s is not a directory", path)
	}
	return &BoundDir{handle: handle, path: path}, nil
}

// OpenBoundDirFile opens one artifact relative to the bound directory. The
// returned handle is the artifact's identity: the caller reads through it and
// never resolves the name again.
func OpenBoundDirFile(dir *BoundDir, rel string, bound int64) (*BoundFile, error) {
	return openBoundDirFile(dir, rel, bound, false)
}

// OpenBoundDirFileClass is the class-aware form: allowTruncation keeps the log and
// event class readable past its bound instead of refusing the open.
func OpenBoundDirFileClass(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*BoundFile, error) {
	return openBoundDirFile(dir, rel, bound, allowTruncation)
}

// openBoundDirFile is the class-aware form: allowTruncation keeps the log and
// event class readable past its bound instead of refusing the open.
func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*BoundFile, error) {
	if dir == nil || dir.handle == 0 {
		return nil, fmt.Errorf("evidence: nil bound directory")
	}
	// P84 T4(d): the bound object itself is not an artifact, and both root spellings
	// ("", ".") are refused through one shared error so the two halves cannot drift.
	if boundRootName(rel) {
		return nil, errBoundRootArtifact(rel)
	}
	clean, err := cleanBoundRelativePath(rel)
	if err != nil {
		return nil, err
	}
	// P84 T4(a): the artifact is resolved one component at a time, each step opened
	// inside the directory handle the previous step returned. There is no derived
	// full path any more, so the bound directory (or an ancestor) being replaced
	// between a name derivation and an open is no longer a window: the kernel only
	// ever looks inside the objects this chain already proved.
	//
	// T4 (P83 fix round 8) is preserved by the walk itself: the leaf is opened with
	// FILE_OPEN_REPARSE_POINT and classified, so a reparse point at the artifact name
	// is not followed, and confirmOpenedArtifact still refuses the object it opened
	// instead of reading the target it redirects to.
	//
	// P84 fix round 3 (P52 semantics): the walk also receives the bound directory, so a
	// reparse point ABOVE the leaf is admitted only when the target IT STORES is proven
	// inside the bound object, and the walk then continues from the bound handle along
	// that target's relative suffix (P84 fix round 6, item F1: see
	// openDescendantRelative/reparseTargetRelativeSuffix).
	handle, err := openDescendantRelative(dir, clean, false)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: clean, Err: err}
	}
	// P42 F1 belt: the chain proves the object was reached inside the bound tree, and
	// the volume serial plus the handle-derived containment below are kept as the same
	// second proof they were before the walk replaced the name derivation.
	if err := dir.confirmOpenedArtifact(handle, clean); err != nil {
		_ = syscall.CloseHandle(handle)
		return nil, err
	}
	return NewBoundFile(os.NewFile(uintptr(handle), clean), bound, allowTruncation)
}

// BoundDirNames lists the entry names of one directory inside the bound directory
// object, sorted (P83 fix round 7, item T4b).
//
// The entries are read from the opened directory HANDLE by readDirNames, which
// calls GetFileInformationByHandleEx(FileIdBothDirectoryInfo) on it (P84 T4 fix
// round 2: os.File.ReadDir is never used), so the names come from the directory
// object itself and never from a second name resolution. T4 (P83 fix
// round 8) made the admission of that handle object-grade: it must not be a reparse
// point, it must live on the bound volume, the bound directory's name must still
// address the bound object through the bound handle, and its handle-derived name
// must sit beneath that base. Ownership that cannot be proven refuses instead of
// enumerating a replacement tree. rel == "" (or ".") names the bound object itself
// (P84 T4(d): both spellings go through boundRootName, so the classifier - not the
// artifact-name guard - decides that), and that root is enumerated straight through
// the bound handle, with no name resolved whatsoever (P84 fix round 6, item F5).
func BoundDirNames(dir *BoundDir, rel string) ([]string, error) {
	if dir == nil || dir.handle == 0 {
		return nil, fmt.Errorf("evidence: nil bound directory")
	}
	clean := "."
	if !boundRootName(rel) {
		cleaned, err := cleanBoundRelativePath(rel)
		if err != nil {
			return nil, err
		}
		clean = cleaned
	}
	var handle syscall.Handle
	owned := false
	if clean == "." {
		// Root branch (P84 fix round 6, item F5): the bound handle IS the enumeration
		// handle, so the entries are read from the object the binding already owns and no
		// name is resolved at all. The branch used to derive the object's handle-derived
		// name and open that name again, which reintroduced exactly the
		// derivation-to-open window the rest of this file removes - and bought nothing,
		// because OpenBoundDir already admitted the object. Without a name resolution
		// there is no window to lose, which is strictly stronger than checking the
		// reopened handle afterwards.
		handle = dir.handle
	} else {
		// P84 T4(a)/(c): a descendant directory is walked inside the bound handle, so
		// a junction at any component above the enumerated name is followed only when
		// its target is proven inside the bound object, and a junction at the
		// enumerated name itself is refused by the leaf rule (P83 T4b), whether it
		// points inside the bound tree or out of it.
		opened, err := openDescendantRelative(dir, clean, true)
		if err != nil {
			return nil, &os.PathError{Op: "open", Path: clean, Err: err}
		}
		handle = opened
		owned = true
	}
	if owned {
		if err := dir.confirmEnumeratedDirectory(handle, clean); err != nil {
			_ = syscall.CloseHandle(handle)
			return nil, err
		}
	}
	// P84 T4 fix round 2: the names are enumerated through this very handle. A handle
	// the walk returned is owned here, so it is closed on both paths and never wrapped
	// in an os.File (see readDirNames for why neither ReadDir flavour is usable); the
	// bound handle itself is owned by the BoundDir.
	names, readErr := readDirNames(handle)
	if owned {
		_ = syscall.CloseHandle(handle)
	}
	if readErr != nil {
		return nil, readErr
	}
	return names, nil
}

// fileIDBothDirInfo mirrors FILE_ID_BOTH_DIR_INFORMATION as the x64 kernel lays it
// out. FileName is the first UTF-16 unit of the trailing name buffer, whose length
// is FileNameLength bytes. FileId sits between ShortName and FileName, and Go's
// natural alignment (ShortName at 70, FileId at 96, FileName at 104) reproduces the
// C offsets, so unsafe.Offsetof(FileName) is the header length the parser skips.
//
// P84 fix round 4 (review L-1): the layout is an ABI contract with the kernel, not a
// local choice. The kernel lays FileName out at offset 104 on x64 (the header ends
// after ShortName at 94, then FileId, then the name buffer), so any deviation is an
// ABI mismatch and must be refused rather than parsed on a guess - see
// assertFileIDBothDirInfoLayout, which readDirNames calls before it reads one entry.
type fileIDBothDirInfo struct {
	NextEntryOffset uint32     // 0: relative offset of the next entry, 0 for the last
	FileIndex       uint32     // 4
	CreationTime    int64      // 8
	LastAccessTime  int64      // 16
	LastWriteTime   int64      // 24
	ChangeTime      int64      // 32
	EndOfFile       int64      // 40
	AllocationSize  int64      // 48
	FileAttributes  uint32     // 56
	FileNameLength  uint32     // 60: FileName length in BYTES
	EaSize          uint32     // 64
	ShortNameLength uint8      // 68: ShortName length in bytes, 0 when absent
	ShortName       [12]uint16 // 70
	FileId          int64      // 96
	FileName        [1]uint16  // 104: first unit of the name buffer
}

const (
	// fileIDBothDirInfoFileNameOffset is the kernel's x64 offset of FileName: the
	// header ends at 94 (after ShortName), FileId occupies 96..103, and the name
	// buffer starts at 104.
	fileIDBothDirInfoFileNameOffset = 104
	// fileIDBothDirInfoMinSize is the smallest safe header size: everything up to and
	// including the first name unit.
	fileIDBothDirInfoMinSize = 104
)

// assertFileIDBothDirInfoLayout refuses to enumerate when the compiled layout of
// fileIDBothDirInfo diverges from the kernel's FILE_ID_BOTH_DIR_INFORMATION.
//
// P84 fix round 4 (review L-1): the parser below derives both the header size it
// skips and the byte it stops at from this struct alone, so a compiler or
// architecture that packs the fields differently (GOARCH=386 lowers int64
// alignment to 4) would read a name out of the middle of FileId without any error
// surfacing. The check is fail-closed: a mismatch returns an error and no entry is
// ever parsed.
func assertFileIDBothDirInfoLayout() error {
	if offset := unsafe.Offsetof(fileIDBothDirInfo{}.FileName); offset != fileIDBothDirInfoFileNameOffset {
		return fmt.Errorf("evidence: enumerate the directory: FILE_ID_BOTH_DIR_INFORMATION.FileName sits at offset %d, not the kernel layout's %d: the compiled layout does not match the ABI, refusing to parse directory entries", offset, fileIDBothDirInfoFileNameOffset)
	}
	if size := unsafe.Sizeof(fileIDBothDirInfo{}); size < fileIDBothDirInfoMinSize {
		return fmt.Errorf("evidence: enumerate the directory: FILE_ID_BOTH_DIR_INFORMATION is %d bytes, smaller than the kernel layout's %d: the compiled layout does not match the ABI, refusing to parse directory entries", size, fileIDBothDirInfoMinSize)
	}
	return nil
}

// readDirNames lists the entry names of the directory the handle addresses. It is
// the only enumerator the Windows half uses (P84 T4 fix round 2).
//
// os.File.ReadDir is unusable here twice over: on Go 1.23+ its handle path drives
// FileFullDirectoryRestartInfo, which blocks forever on a handle NtCreateFile opened,
// and on the Go 1.22 toolchain it falls back to a by-name FindFirstFileW, which
// enumerates a relative name against the wrong directory. The handle-driven call
// therefore stays explicit: GetFileInformationByHandleEx(FileIdBothDirectoryInfo)
// returns as many entries as fit and keeps the position in the handle, so the loop
// below resumes until the kernel reports the directory exhausted.
//
// "." and ".." are not entries of the directory object and are dropped. The result
// is sorted so callers see one stable order. No failure is ever reported as an empty
// list: a refused or unintelligible answer returns the error that produced it.
func readDirNames(handle syscall.Handle) ([]string, error) {
	if err := assertFileIDBothDirInfoLayout(); err != nil {
		return nil, err
	}
	names := make([]string, 0, 32)
	// The first call of every scan restarts it; the calls that continue inside the
	// same scan use the non-restart class so they resume where the buffer ended.
	restart := true
	for size := uint32(readDirBufferStart); ; {
		buffer := make([]byte, size)
		class := uintptr(fileIdBothDirectoryInfo)
		if restart {
			class = uintptr(fileIdBothDirectoryInfoRestart)
		}
		// GetFileInformationByHandleEx takes no out byte count: a successful answer IS
		// the entry chain in the buffer, and the end of the directory arrives as an
		// error from a later call. The first entry is therefore checked for actual
		// data instead of trusting a count that does not exist.
		result, _, callErr := procGetFileInformationByHandleEx.Call(
			uintptr(handle),
			class,
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(size),
		)
		if result == 0 {
			errno, _ := callErr.(syscall.Errno)
			switch errno {
			case syscall.ERROR_NO_MORE_FILES, syscall.ERROR_FILE_NOT_FOUND:
				// The directory is exhausted: the first call on an empty directory
				// reports this, and some filesystems spell the same end state
				// ERROR_FILE_NOT_FOUND. Both are the normal end of the walk.
				sort.Strings(names)
				return names, nil
			case syscall.ERROR_MORE_DATA, syscall.ERROR_INSUFFICIENT_BUFFER, syscall.ERROR_BUFFER_OVERFLOW:
				// No entry fit (or the answer overflowed), so the position did not
				// move: a larger buffer resumes exactly where this call started.
				if size >= readDirBufferLimit {
					return nil, fmt.Errorf("evidence: enumerate the directory: %v at the %d byte buffer limit", callErr, readDirBufferLimit)
				}
				size *= 2
				continue
			}
			return nil, fmt.Errorf("evidence: enumerate the directory: %v", callErr)
		}
		base := unsafe.Pointer(&buffer[0])
		// The scan has started, so the follow-up calls must not restart it.
		restart = false
		nameOffset := unsafe.Offsetof(fileIDBothDirInfo{}.FileName)
		// Entries are chained by NextEntryOffset; a zero offset marks the last one
		// in this buffer, after which the next call continues the enumeration.
		for offset := uintptr(0); ; {
			// P84 fix round 4 (review L-2): the chain is written by the kernel, but a
			// malformed or hostile answer must not turn into an out-of-bounds read. The
			// header has to fit before it is dereferenced, and the name it declares has
			// to fit before it is sliced.
			if offset+nameOffset > uintptr(len(buffer)) {
				return nil, fmt.Errorf("evidence: enumerate the directory: malformed directory entry at offset %d: the %d byte header does not fit the %d byte buffer", offset, nameOffset, len(buffer))
			}
			entry := (*fileIDBothDirInfo)(unsafe.Add(base, offset))
			if offset+nameOffset+uintptr(entry.FileNameLength) > uintptr(len(buffer)) {
				return nil, fmt.Errorf("evidence: enumerate the directory: malformed directory entry at offset %d: a %d byte name does not fit the %d byte buffer", offset, entry.FileNameLength, len(buffer))
			}
			units := uintptr(entry.FileNameLength) / 2
			if units == 0 {
				// A nameless entry cannot exist in a directory, so this is an answer
				// the kernel did not fill. Refusing it keeps the walk from looping on
				// an unfilled buffer and from returning a list emptier than the
				// directory: a failure is never signalled as an empty result.
				return nil, fmt.Errorf("evidence: enumerate the directory: the kernel answered without a directory entry")
			}
			if name := syscall.UTF16ToString(unsafe.Slice((*uint16)(unsafe.Add(base, offset+nameOffset)), units)); name != "." && name != ".." {
				names = append(names, name)
			}
			if entry.NextEntryOffset == 0 {
				break
			}
			offset += uintptr(entry.NextEntryOffset)
		}
	}
}

// confirmEnumeratedDirectory proves the directory handle that will be enumerated
// still addresses the bound object (rel ".") or a strict descendant of it.
//
// P84 fix round 6 (F5): BoundDirNames no longer routes its root branch through this
// function - the root is enumerated through the bound handle itself, so no name is
// derived and reopened. The rel "." form is kept because it is the direct-call seam
// the T4a counterexample exercises (a handle opened at the bound NAME after the
// object was renamed away must still be refused), and because it states the
// object-level rule without depending on where the handle came from.
func (dir *BoundDir) confirmEnumeratedDirectory(handle syscall.Handle, rel string) error {
	if rel != "." {
		return dir.confirmOpenedArtifact(handle, rel)
	}
	boundInfo, err := dir.byHandleFileInformation()
	if err != nil {
		return fmt.Errorf("evidence: identify the bound directory %s: %v", dir.path, err)
	}
	openedInfo, err := fileInformationOf(handle)
	if err != nil {
		return fmt.Errorf("evidence: identify the enumerated directory: %v", err)
	}
	if boundInfo.VolumeSerialNumber != openedInfo.VolumeSerialNumber ||
		boundInfo.FileIndexHigh != openedInfo.FileIndexHigh ||
		boundInfo.FileIndexLow != openedInfo.FileIndexLow {
		return fmt.Errorf("evidence: the enumerated directory is not the bound directory %s", dir.path)
	}
	// T4 (P83 fix round 8): the handle is the bound object, but the name it will
	// enumerate was derived separately. Ownership is only proven when that name
	// still addresses this very object, so a bound directory that was replaced
	// after it was opened refuses here instead of listing a replacement tree.
	base, err := handleFinalPath(handle)
	if err != nil {
		return fmt.Errorf("evidence: resolve the enumerated directory: %v", err)
	}
	return dir.verifyBoundObjectStillNamed(base)
}

// confirmOpenedArtifact proves the artifact just opened still belongs to the bound
// directory object.
//
// T4 (P83 fix round 8): the comparison is no longer purely a comparison of two path
// strings. The opened object is classified through its own handle (a reparse point
// is refused), its volume is compared with the bound handle's volume, the bound
// directory's name is proven to still address the bound object through the bound
// handle, and only then are the two handle-derived names compared for containment.
// A replacement of the bound directory, a reparse point at the artifact name, and a
// junction to another volume therefore all refuse before any byte is read, and a
// case where ownership cannot be proven refuses instead of being admitted.
func (dir *BoundDir) confirmOpenedArtifact(handle syscall.Handle, rel string) error {
	openedInfo, err := fileInformationOf(handle)
	if err != nil {
		return fmt.Errorf("evidence: identify the opened artifact %s: %v", rel, err)
	}
	if openedInfo.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0 {
		return fmt.Errorf("evidence: artifact %s is a reparse point, which cannot be proven to belong to the bound directory %s", rel, dir.path)
	}
	base, err := dir.currentName()
	if err != nil {
		return fmt.Errorf("evidence: resolve the bound directory %s: %v", dir.path, err)
	}
	if err := dir.verifyBoundObjectStillNamed(base); err != nil {
		return err
	}
	boundInfo, err := dir.byHandleFileInformation()
	if err != nil {
		return fmt.Errorf("evidence: identify the bound directory %s: %v", dir.path, err)
	}
	if boundInfo.VolumeSerialNumber != openedInfo.VolumeSerialNumber {
		return fmt.Errorf("evidence: artifact %s is on volume %d, which is not the volume %d of the bound directory %s", rel, openedInfo.VolumeSerialNumber, boundInfo.VolumeSerialNumber, dir.path)
	}
	final, err := handleFinalPath(handle)
	if err != nil {
		return fmt.Errorf("evidence: resolve the opened artifact %s: %v", rel, err)
	}
	if !isUnderDirectory(base, final) {
		return fmt.Errorf("evidence: artifact %s resolved to %s, which is outside the bound directory %s", rel, final, base)
	}
	return nil
}

// verifyBoundObjectStillNamed proves, through the bound handle itself, that the
// name derived from that handle still addresses the bound object.
//
// T4 (P83 fix round 8): Windows exposes no openat(2), so a descendant is opened
// under a name. That name is derived from the bound handle, but the object it
// resolves to can be replaced between the derivation and the open. Comparing the
// volume serial and file index of the bound handle with the object the name
// resolves to refuses exactly that case: ownership of the name can no longer be
// proven, so the caller refuses instead of describing a replacement tree.
func (dir *BoundDir) verifyBoundObjectStillNamed(base string) error {
	same, err := dir.sameObjectAs(base)
	if err != nil {
		return fmt.Errorf("evidence: the directory named %s can no longer be identified: %v", base, err)
	}
	if !same {
		return fmt.Errorf("evidence: the directory named %s is no longer the bound directory %s", base, dir.path)
	}
	return nil
}

// P84 T4(a): handle-relative descendant resolution.
//
// Windows has no openat(2), but ntdll!NtCreateFile accepts a directory handle in
// OBJECT_ATTRIBUTES.RootDirectory, which is the same capability: the component is
// resolved by the kernel INSIDE the object that handle names. The walk below opens
// a relative name one component at a time, feeding each verified directory handle
// in as the root of the next step, so:
//
//   - no full path is ever derived and re-resolved, which removes the derivation-to-
//     open window P42's comment describes;
//   - a case-variant twin directory sharing the bound object's spelling up to case
//     is unreachable, because it is a SIBLING of the object the parent handle names
//     and the kernel never leaves that object;
//   - a reparse point at the LAST component is opened as ITSELF and refused
//     (P83 T4b), so the artifact name is never redirected;
//   - a reparse point in the MIDDLE of the name is opened as itself, classified,
//     then reopened with FILE_OPEN_REPARSE_POINT dropped so the kernel follows it,
//     and the followed handle must still be proven inside the bound object before
//     the walk continues on it (P84 fix round 3, the P52 control shape). A redirect
//     that leaves the bound tree is refused with the boundary message, and one that
//     stays inside it is admitted, which is the only reading that keeps a
//     name-relative walk usable on Windows at all.
//
// Design alternatives rejected by the P84 architecture design: one whole-name open
// with FILE_OPEN_REPARSE_POINT (only the last component is affected) and
// "GetFinalPathNameByHandle, reopen, compare" (the exact name-then-open shape).
const (
	// objCaseInsensitive is OBJ_CASE_INSENSITIVE, the flag that keeps component
	// matching case-insensitive inside the parent object.
	objCaseInsensitive = 0x00000040
	// createDispositionOpen is FILE_OPEN: every component must already exist.
	createDispositionOpen = 0x00000001
	// createOptionDirectoryFile and createOptionNonDirectoryFile pin the expected
	// type of each step at the kernel, so a swap of a component's type fails at the
	// step that opens it.
	createOptionDirectoryFile    = 0x00000001
	createOptionNonDirectoryFile = 0x00000040
	// createOptionOpenForBackupIntent is FILE_OPEN_FOR_BACKUP_INTENT, the
	// CreateOptions form of the FILE_FLAG_BACKUP_SEMANTICS intent OpenBoundDir uses,
	// so a directory handle opens with the same reach the binding itself has.
	createOptionOpenForBackupIntent = 0x00004000
	// createOptionOpenReparsePoint is FILE_OPEN_REPARSE_POINT in the CreateOptions
	// namespace; the numeric value is the same one fileFlagOpenReparsePoint carries
	// for CreateFile.
	createOptionOpenReparsePoint = fileFlagOpenReparsePoint
	// statusSuccess is STATUS_SUCCESS.
	statusSuccess = 0
	// win32InvalidParameter is ERROR_INVALID_PARAMETER, the fallback errno for an
	// NTSTATUS the Win32 mapping does not translate.
	win32InvalidParameter = 87

	// P84 T4 fix round 2: the desired-access bits of every walk step are spelled
	// out instead of asking for GENERIC_READ. Measured on this host, a GENERIC_READ
	// handle answers GetFileInformationByHandleEx(FileIdBothDirectoryInfo) with
	// ERROR_INVALID_PARAMETER (87), and the Go >= 1.23 handle path of os.File.ReadDir
	// (FileFullDirectoryRestartInfo) blocks forever on a handle NtCreateFile opened.
	// FILE_LIST_DIRECTORY|FILE_READ_ATTRIBUTES|SYNCHRONIZE is the combination that
	// both enumerates and identifies, and the artifact leaf needs FILE_READ_DATA
	// instead of the list bit (same value, file sense).
	//
	// SYNCHRONIZE is required because every step is waited on through its
	// IO_STATUS_BLOCK. FILE_SYNCHRONOUS_IO_ALERT/NONALERT cannot replace it: adding
	// either makes NtCreateFile answer STATUS_INVALID_PARAMETER.
	accessListOrRead     = 0x00000001 // FILE_LIST_DIRECTORY (directory) / FILE_READ_DATA (file)
	accessReadAttributes = 0x00000080 // FILE_READ_ATTRIBUTES
	accessSynchronize    = 0x00100000 // SYNCHRONIZE

	// fileIdBothDirectoryInfo is FILE_INFO_BY_HANDLE_CLASS's
	// FileIdBothDirectoryInfo, the class of FILE_ID_BOTH_DIR_INFORMATION: the entry
	// name plus its file index, returned one buffer at a time with the enumeration
	// position kept in the handle between calls (it is not a restart class).
	//
	// The value is 10, not 3: FILE_INFO_BY_HANDLE_CLASS (GetFileInformationByHandleEx)
	// and FileInformationClass (NtQueryInformationFile) number their enums
	// differently, and 3 means FileRenameInfo to the former. This host answers the
	// wrong number with ERROR_INVALID_PARAMETER (87).
	fileIdBothDirectoryInfo = 10
	// fileIdBothDirectoryInfoRestart is FileIdBothDirectoryRestartInfo: the same
	// answer as FileIdBothDirectoryInfo, but the scan restarts at the first entry.
	// P84 fix round 6 (F5): the root of a bound directory is enumerated through the
	// BOUND handle, which outlives one call and carries the position of the previous
	// scan, so a second enumeration would otherwise answer "no more entries". Leading
	// every scan with the restart class keeps one call self-contained and makes the
	// bound handle reusable without ever re-resolving its name.
	fileIdBothDirectoryInfoRestart = 11
	// readDirBufferStart is the first enumeration buffer tried;
	// readDirBufferLimit bounds the grow-and-retry loop that ERROR_MORE_DATA drives.
	readDirBufferStart = 64 << 10
	readDirBufferLimit = 1 << 20

	// P84 fix round 6 (F1): FSCTL_GET_REPARSE_POINT asks a reparse point about ITSELF -
	// the target it stores - instead of letting the kernel apply that target to the
	// rest of a name. The walk below uses it as the mechanism of the containment
	// proof: the redirect is read as data, the target is compared with the bound
	// object, and only what the comparison admits is opened, component by component,
	// inside the bound handle.
	fsctlGetReparsePoint = 0x000900A8 // FSCTL_GET_REPARSE_POINT
	// ioReparseTagMountPoint is IO_REPARSE_TAG_MOUNT_POINT (a junction). Its stored
	// target is always an absolute path, so it is the tag the walk can prove.
	ioReparseTagMountPoint = 0xA0000003
	// ioReparseTagSymlink is IO_REPARSE_TAG_SYMLINK. An absolute target is handled like
	// a junction; a relative one refuses (see reparseTargetRelativeSuffix).
	ioReparseTagSymlink = 0xA000000C
	// symlinkFlagRelative is SYMLINK_FLAG_RELATIVE, the bit that marks a symbolic
	// link's target as relative to the link's own directory.
	symlinkFlagRelative = 0x00000001
	// reparseDataHeaderSize is the fixed part of REPARSE_DATA_BUFFER: ReparseTag
	// (ULONG), ReparseDataLength (USHORT) and Reserved (USHORT).
	reparseDataHeaderSize = 8
	// reparseMountPointPathOffset and reparseSymlinkPathOffset are the offsets of the
	// tag-specific PathBuffer. A mount point body is four USHORTs (16 bytes from the
	// start), a symbolic link body adds its ULONG Flags (20 bytes from the start).
	reparseMountPointPathOffset = reparseDataHeaderSize + 8
	reparseSymlinkPathOffset    = reparseDataHeaderSize + 12
	// reparseDataBufferLimit is MAXIMUM_REPARSE_DATA_BUFFER_SIZE, the size the kernel
	// documents for one reparse buffer.
	reparseDataBufferLimit = 16 << 10
	// reparseWalkLimit bounds how many reparse points one name may cross. A junction
	// whose target contains itself would otherwise restart the walk forever.
	reparseWalkLimit = 32
)

// unicodeString mirrors UNICODE_STRING. Length and MaximumLength are in bytes and
// exclude/includ the terminating NUL respectively.
type unicodeString struct {
	Length        uint16
	MaximumLength uint16
	Buffer        *uint16
}

// objectAttributes mirrors OBJECT_ATTRIBUTES. A zero Length makes NtCreateFile
// answer STATUS_INVALID_PARAMETER, so Length is always unsafe.Sizeof of this type.
type objectAttributes struct {
	Length                   uint32
	RootDirectory            syscall.Handle
	ObjectName               *unicodeString
	Attributes               uint32
	SecurityDescriptor       unsafe.Pointer
	SecurityQualityOfService unsafe.Pointer
}

// ioStatusBlock mirrors IO_STATUS_BLOCK; NtCreateFile writes the completion
// status and information here.
type ioStatusBlock struct {
	Status      uintptr
	Information uintptr
}

var (
	procNtCreateFile          = syscall.NewLazyDLL("ntdll.dll").NewProc("NtCreateFile")
	procRtlNtStatusToDosError = syscall.NewLazyDLL("ntdll.dll").NewProc("RtlNtStatusToDosError")
	// procGetFileInformationByHandleEx drives readDirNames; it is the only way to
	// enumerate a handle NtCreateFile opened without blocking or re-resolving a name.
	procGetFileInformationByHandleEx = syscall.NewLazyDLL("kernel32.dll").NewProc("GetFileInformationByHandleEx")
	// procDeviceIoControl reads a reparse point's stored target (P84 F1). It is a
	// kernel32 call like the other two, so no dependency is added.
	procDeviceIoControl = syscall.NewLazyDLL("kernel32.dll").NewProc("DeviceIoControl")
)

// openDescendantRelative opens clean - a cleaned relative name whose components
// cannot escape - as a descendant of the bound directory object, one component at
// a time.
//
// wantDir selects the type of the LAST component: a directory when the caller
// enumerates it, a regular entry when the caller opens an artifact. Every
// intermediate component is a directory by construction. Each step's handle is
// closed as soon as the next step exists, so no chain of handles outlives the walk.
//
// P84 fix round 3 (the P52 contract shape) survived the P84 fix round 6 mechanism
// change: a reparse point at the LAST component is refused outright, whether the
// caller opens an artifact or enumerates a directory (P83 T4b), and a reparse point
// in the MIDDLE is followed only when its target is proven inside the bound object.
// What changed is HOW that target is obtained (F1): the reparse point is asked for
// the path it stores (FSCTL_GET_REPARSE_POINT) and the stored path - never the
// kernel's redirect - is compared with the bound object's own handle-derived name.
// An in-bound target contributes its relative suffix, and the walk restarts at the
// BOUND handle with that suffix in front of the components that are still left.
// Nothing is therefore opened by a name derived from the target, so there is no
// "follow, then compare" window in which the target could be swapped: the redirect
// is decided as data, and every component is thereafter opened inside the object the
// binding already owns. An escape, an unprovable target and an unusable tag are
// refused with a message naming the boundary.
func openDescendantRelative(dir *BoundDir, clean string, wantDir bool) (syscall.Handle, error) {
	components := strings.Split(clean, string(filepath.Separator))
	parent := dir.handle
	owned := false
	relative := ""
	restarts := 0
	for index := 0; index < len(components); index++ {
		component := components[index]
		last := index == len(components)-1
		if relative == "" {
			relative = component
		} else {
			relative = relative + string(filepath.Separator) + component
		}
		handle, reparse, err := openRelativeComponent(parent, component, !last || wantDir, false)
		if err == nil && reparse {
			// The object the name addresses is a reparse point. It was opened as
			// itself, so closing it never touches the target.
			_ = syscall.CloseHandle(handle)
			if last {
				err = fmt.Errorf("evidence: %s is a reparse point inside the bound directory, which cannot be proven to belong to it", component)
			} else {
				var suffix []string
				suffix, err = reparseTargetRelativeSuffix(dir, parent, component, relative)
				if err == nil {
					// The redirect is replaced by the stored target's suffix, walked
					// from the bound handle. The accumulated chain is dropped first, so
					// the continuation cannot run inside an object the exit left.
					if owned {
						_ = syscall.CloseHandle(parent)
						owned = false
					}
					restarts++
					if restarts > reparseWalkLimit {
						err = fmt.Errorf("evidence: %s crosses more than %d reparse points inside the bound directory", relative, reparseWalkLimit)
					} else {
						parent = dir.handle
						relative = ""
						components = append(suffix, components[index+1:]...)
						index = -1
						continue
					}
				}
			}
		}
		// The previous step's handle is closed before anything is reported: a failed
		// walk must not leak the handles it already opened, and a closed number must
		// not stay addressable while the caller handles the error.
		if owned {
			_ = syscall.CloseHandle(parent)
		}
		if err != nil {
			return 0, err
		}
		parent = handle
		owned = true
	}
	return parent, nil
}

// reparseTargetRelativeSuffix reads the target a middle-of-the-name reparse point
// stores and reports it as the components that remain to be walked under the bound
// object (P84 fix round 6, item F1).
//
// The mechanism is the fix: the old shape let the kernel follow the redirect and
// then compared the followed handle's name, which left one window - the object the
// name addressed could be swapped between the kernel's resolution and the
// comparison. Here the redirect is never applied to the remaining name at all. The
// reparse point is opened as itself, its stored target is read as data, the target
// is compared with the bound object's handle-derived name by Ordinal bytes, and only
// the admitted suffix is walked - inside the bound handle.
//
// A junction (IO_REPARSE_TAG_MOUNT_POINT) always stores an absolute target, which is
// the shape that can be compared. A symbolic link (IO_REPARSE_TAG_SYMLINK) is
// admitted the same way when its target is absolute and refused when it is relative:
// a target relative to the link's own directory cannot be compared against the bound
// object without resolving it against a name, which is the shape being removed here,
// so it refuses (fail-closed). Any other tag, an unreadable buffer, a target that is
// not absolute, a target that is not already in canonical form (P84 fix round 10 /
// N10) and a target outside the bound object all refuse with the boundary message.
func reparseTargetRelativeSuffix(dir *BoundDir, parent syscall.Handle, component, relative string) ([]string, error) {
	point, _, err := openRelativeComponent(parent, component, true, false)
	if err != nil {
		return nil, fmt.Errorf("evidence: intermediate reparse point %s resolves outside the bound directory or cannot be opened: %v", relative, err)
	}
	target, readErr := reparseStoredTarget(point)
	_ = syscall.CloseHandle(point)
	if readErr != nil {
		return nil, fmt.Errorf("evidence: intermediate reparse point %s resolves outside the bound directory or cannot be read: %v", relative, readErr)
	}
	base, err := dir.currentName()
	if err != nil {
		return nil, fmt.Errorf("evidence: intermediate reparse point %s cannot be proven inside the bound directory %s: %v", relative, dir.path, err)
	}
	// P84 fix round 10 (N10): a stored target that filepath.Clean would rewrite is
	// refused here with the reason spelled out, so the operator sees that the shape -
	// not the boundary - is what failed. The boundary comparison below then works on
	// the stored bytes alone.
	if reason := nonCanonicalReparseTarget(target); reason != "" {
		return nil, fmt.Errorf("evidence: intermediate reparse point %s stores the target %q, which cannot be interpreted losslessly: it is not the exact form filepath.Clean produces (%s); a target is only compared byte for byte when its own spelling is already canonical", relative, target, reason)
	}
	suffix, ok := boundRelativeSuffix(base, target)
	if !ok {
		return nil, fmt.Errorf("evidence: intermediate reparse point %s resolves to %s, which is outside the bound directory %s", relative, target, base)
	}
	return suffix, nil
}

// reparseStoredTarget reads the target one reparse point stores, as the ordinary
// absolute DOS path the handle-derived names use.
//
// REPARSE_DATA_BUFFER is parsed by offset rather than through a Go struct: its body
// is a tag-dependent union whose PathBuffer is the last member, and the tag decides
// both where that buffer starts and whether a Flags word precedes it, so an offset
// read states the layout explicitly and cannot be silently changed by field packing.
//
// P84 fix round 11 (N14): the UTF-16 names are decoded by reparseNames, which scans
// each declared range for an embedded NUL before UTF16ToString and refuses, so a name
// that would otherwise decode as a prefix cannot be judged on that prefix.
func reparseStoredTarget(handle syscall.Handle) (string, error) {
	buffer := make([]byte, reparseDataBufferLimit)
	var returned uint32
	result, _, callErr := procDeviceIoControl.Call(
		uintptr(handle),
		uintptr(fsctlGetReparsePoint),
		0, // input buffer: FSCTL_GET_REPARSE_POINT takes none
		0,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&returned)),
		0, // no OVERLAPPED: synchronous
	)
	if result == 0 {
		return "", callErr
	}
	if int(returned) < reparseDataHeaderSize {
		return "", fmt.Errorf("evidence: the reparse buffer is %d bytes, shorter than the %d byte header", returned, reparseDataHeaderSize)
	}
	size := int(returned)
	tag := binary.LittleEndian.Uint32(buffer[0:4])
	dataLength := int(binary.LittleEndian.Uint16(buffer[4:6]))
	if dataLength > size-reparseDataHeaderSize {
		return "", fmt.Errorf("evidence: the reparse buffer declares %d bytes of body in %d bytes", dataLength, size-reparseDataHeaderSize)
	}
	body := buffer[reparseDataHeaderSize:size]
	var pathOffset int
	var relative bool
	switch tag {
	case ioReparseTagMountPoint:
		pathOffset = reparseMountPointPathOffset
	case ioReparseTagSymlink:
		pathOffset = reparseSymlinkPathOffset
		if len(body) < pathOffset-reparseDataHeaderSize {
			return "", fmt.Errorf("evidence: the symbolic link reparse buffer is %d bytes, shorter than the %d byte header", size, pathOffset)
		}
		relative = binary.LittleEndian.Uint32(body[reparseSymlinkPathOffset-4-reparseDataHeaderSize:pathOffset-reparseDataHeaderSize])&symlinkFlagRelative != 0
	default:
		return "", fmt.Errorf("evidence: reparse tag 0x%08X is not a mount point or a symbolic link, so its target cannot be proven", tag)
	}
	// The name offsets are relative to PathBuffer but the bounds are checked against
	// the whole buffer, so the absolute form (buffer[:size]) is what is passed down.
	substitute, printName, err := reparseNames(buffer[:size], pathOffset)
	if err != nil {
		return "", err
	}
	if relative {
		// The stored target is relative to the link's own directory; see the caller.
		return "", fmt.Errorf("evidence: the symbolic link stores the relative target %q, which cannot be compared with the bound object", substitute)
	}
	// P84 fix round 7 (N3): the SUBSTITUTE name is the only authority for where the
	// reparse point leads. It is the name the object manager resolves; the print name is
	// a display string that the same link may carry a different value for, so deciding
	// on it - or falling back to it - would let a reparse point whose two names disagree
	// be judged by the friendlier one, and a boundary decision must not be made on a
	// field the kernel does not follow. The print name is kept for diagnostics only.
	if substitute == "" {
		return "", fmt.Errorf("evidence: the reparse point stores no substitute name (print name %q), so its target cannot be proven", printName)
	}
	normalized := normalizeReparseTarget(substitute)
	if !filepath.IsAbs(normalized) {
		return "", fmt.Errorf("evidence: the stored reparse target %q (print name %q) is not an absolute path", substitute, printName)
	}
	return normalized, nil
}

// reparseNames extracts the substitute and print names of a mount point or symbolic
// link body. buffer is the whole REPARSE_DATA_BUFFER and pathOffset is the absolute
// offset of its PathBuffer; the stored name offsets are relative to that buffer, so
// each declared range is translated to absolute offsets and checked before it is read.
func reparseNames(buffer []byte, pathOffset int) (substitute, printName string, err error) {
	if len(buffer) < pathOffset {
		return "", "", fmt.Errorf("evidence: the reparse buffer is %d bytes, shorter than the %d byte path header", len(buffer), pathOffset)
	}
	fields := buffer[reparseDataHeaderSize:pathOffset]
	substituteOffset := int(binary.LittleEndian.Uint16(fields[0:2]))
	substituteLength := int(binary.LittleEndian.Uint16(fields[2:4]))
	printOffset := int(binary.LittleEndian.Uint16(fields[4:6]))
	printLength := int(binary.LittleEndian.Uint16(fields[6:8]))
	pathBuffer := buffer[pathOffset:]
	read := func(kind string, offset, length int) (string, error) {
		if length == 0 {
			return "", nil
		}
		if length%2 != 0 || offset < 0 || offset+length > len(pathBuffer) {
			return "", fmt.Errorf("evidence: the reparse name at offset %d with %d bytes does not fit the %d byte path buffer", offset, length, len(pathBuffer))
		}
		units := make([]uint16, length/2)
		for unit := range units {
			units[unit] = binary.LittleEndian.Uint16(pathBuffer[offset+2*unit : offset+2*unit+2])
		}
		// P84 fix round 11 (N14): UTF16ToString stops at the first NUL unit it meets, so
		// a NUL INSIDE the declared range would decode a prefix of the stored name and
		// the boundary would then be judged on that prefix. The declared range is
		// scanned before decoding and an embedded NUL refuses (fail-closed, so the
		// refusal can only widen); the last unit of the range is read as the terminator
		// NTFS writes when it counts that terminator in the length.
		for unit, value := range units {
			if value == 0 && unit != len(units)-1 {
				return "", fmt.Errorf("evidence: the %s reparse name at offset %d with %d bytes contains an embedded NUL unit at offset %d, so the stored name cannot be decoded losslessly", kind, offset, length, offset+2*unit)
			}
		}
		return syscall.UTF16ToString(units), nil
	}
	if substitute, err = read("substitute", substituteOffset, substituteLength); err != nil {
		return "", "", err
	}
	if printName, err = read("print", printOffset, printLength); err != nil {
		return "", "", err
	}
	return substitute, printName, nil
}

// normalizeReparseTarget turns the object-manager form a reparse target is stored in
// into the ordinary DOS form handleFinalPath reports, so the two can be compared by
// bytes.
//
// P84 fix round 9 (N7): the target is NOT trimmed. Trimming moved a name whose last
// character is a legal-but-blank code point (U+00A0 is legal in an NTFS file name) onto
// the name without it, so a sibling that merely looked like the bound object - e.g.
// "<bound>\u00A0" stored as the redirect target - normalised to the bound name itself
// and was admitted as the bound object. From here on only the exact object-manager
// prefixes are stripped, once each, and every remaining byte, blanks included, is what
// boundRelativeSuffix judges by Ordinal equality / strict-descendant bytes. A target
// that is nothing but a prefix strips to "" and is refused by the caller's
// absolute-path check.
func normalizeReparseTarget(target string) string {
	const ntPrefix = `\??\`
	const globalPrefix = `\GLOBAL??\`
	switch {
	case strings.HasPrefix(target, ntPrefix+`UNC\`):
		target = `\\` + target[len(ntPrefix+`UNC\`):]
	case strings.HasPrefix(target, ntPrefix):
		target = target[len(ntPrefix):]
	case strings.HasPrefix(target, globalPrefix+`UNC\`):
		target = `\\` + target[len(globalPrefix+`UNC\`):]
	case strings.HasPrefix(target, globalPrefix):
		target = target[len(globalPrefix):]
	}
	return stripExtendedPathPrefix(target)
}

// nonCanonicalReparseTarget explains why a stored reparse target is not already the
// exact form filepath.Clean produces, or returns "" when it is.
//
// P84 fix round 10 (N10): the target is judged by its own bytes. filepath.Clean folds
// "." and ".." components, repeated separators and - on Windows - forward slashes, so
// cleaning the target before comparing it let a stored target that names a DIFFERENT
// object be rewritten onto an admitted name (e.g. "<bound>\x\..\outside" collapsed to
// "<bound>\outside"). A target whose bytes change under Clean cannot be interpreted
// losslessly, so it is refused (fail-closed) rather than folded. Forward slashes are
// refused by the same rule even though the object manager accepts them: the
// handle-derived names are backslash forms, so admitting slash spellings would mean
// comparing a rewritten string instead of the stored one. A trailing separator is
// refused for the same reason.
func nonCanonicalReparseTarget(target string) string {
	if filepath.Clean(target) == target {
		return ""
	}
	// The leading UNC marker is part of the canonical form, so it is excluded before
	// repeated separators are looked for.
	body := strings.TrimPrefix(target, `\\`)
	shapes := make([]string, 0, 4)
	seen := make(map[string]bool, 4)
	add := func(shape string) {
		if !seen[shape] {
			seen[shape] = true
			shapes = append(shapes, shape)
		}
	}
	if strings.Contains(body, "/") {
		add("a forward slash, which is not the separator the handle-derived names use")
	}
	if strings.Contains(body, `\\`) {
		add("a repeated separator")
	}
	if strings.HasSuffix(target, `\`) || strings.HasSuffix(target, "/") {
		add("a trailing separator")
	}
	for _, part := range strings.Split(strings.ReplaceAll(body, "/", `\`), `\`) {
		switch part {
		case ".":
			add(`a "." component`)
		case "..":
			add(`a ".." component`)
		}
	}
	if len(shapes) == 0 {
		return "it is not the form filepath.Clean produces"
	}
	return strings.Join(shapes, ", ")
}

// boundRelativeSuffix reports the components that lead from the bound object's
// handle-derived name to a reparse target, and whether the target is inside it.
//
// The comparison is Ordinal (P84 T4(b)): both names are absolute paths in the same
// DOS form, so a target that is genuinely inside the bound object differs from it by
// whole components only. A target that is the bound object itself contributes no
// components and is admitted; a case-variant twin (the shape W1 pins) differs in
// bytes and is refused, as is an absolute target on another volume, a relative one
// and one whose remaining components could walk upwards.
//
// P84 fix round 10 (N10): the target is NOT cleaned before the comparison. The
// canonical-form precondition is checked first, and once it holds the stored bytes
// themselves are compared - so the object the target names is the object judged, and
// a non-canonical spelling can never be folded into a different admitted path.
func boundRelativeSuffix(base, target string) ([]string, bool) {
	if nonCanonicalReparseTarget(target) != "" {
		return nil, false
	}
	if !filepath.IsAbs(target) {
		return nil, false
	}
	// P84 fix round 11 (N13): equality is judged on the canonical form AS IS, so a
	// volume root keeps the trailing separator Clean preserves ("D:\" stays "D:\").
	// Trimming first compared the root's "D:" against the target "D:\", failed the
	// equality branch and fell into the strict-descendant branch, where the target IS
	// the prefix - a deterministic in-bounds refusal. The trim is now applied only
	// when the strict-descendant prefix is built, so the root's own descendants still
	// compare against the prefix "D:\"; no refusal condition is relaxed.
	cleanBase := filepath.Clean(base)
	if target == cleanBase {
		return nil, true
	}
	trimmedBase := strings.TrimRight(cleanBase, `\/`)
	if trimmedBase == "" {
		return nil, false
	}
	prefix := trimmedBase + string(filepath.Separator)
	if len(target) <= len(prefix) || target[:len(prefix)] != prefix {
		return nil, false
	}
	rest := target[len(prefix):]
	parts := strings.Split(rest, string(filepath.Separator))
	for _, part := range parts {
		if part == "" || part == "." || part == ".." {
			return nil, false
		}
	}
	return parts, true
}

// utf16ComponentLength reports the byte length one component occupies in a
// UNICODE_STRING, and whether that length fits the 16-bit Length field (P84 F6).
//
// The bound is MaxUint16-2 rather than MaxUint16 because MaximumLength carries the
// same field and has to hold Length plus the terminating NUL; a length that fits
// Length but not MaximumLength would wrap and describe a name the kernel never sees.
func utf16ComponentLength(name string) (uint16, bool) {
	units, err := syscall.UTF16FromString(name)
	if err != nil {
		return 0, false
	}
	size := uint64(len(units)-1) * 2
	if size > math.MaxUint16-2 {
		return 0, false
	}
	return uint16(size), true
}

// openRelativeComponent opens one component inside parent and reports whether the
// object the name addresses is itself a reparse point.
//
// followReparse drops FILE_OPEN_REPARSE_POINT, which makes the kernel follow the
// reparse point and hand back the target object instead. The walk uses the flag for
// every component it opens by name, and omits it only to follow a reparse point it
// has already classified and is about to prove.
func openRelativeComponent(parent syscall.Handle, component string, directory, followReparse bool) (syscall.Handle, bool, error) {
	units, err := syscall.UTF16FromString(component)
	if err != nil {
		return 0, false, err
	}
	// P84 fix round 6 (F6): UNICODE_STRING.Length is 16 bits, so a component whose
	// UTF-16 form does not fit was silently truncated before - the kernel would then
	// resolve a DIFFERENT, shorter name. The length is now proven to fit and the open
	// refuses otherwise; a component that cannot be named truthfully cannot be opened.
	length, ok := utf16ComponentLength(component)
	if !ok {
		return 0, false, fmt.Errorf("evidence: the component name of %d UTF-16 bytes cannot be opened: a UNICODE_STRING name is limited to %d bytes", (len(units)-1)*2, math.MaxUint16)
	}
	objectName := unicodeString{
		Length:        length,
		MaximumLength: length + 2,
		Buffer:        &units[0],
	}
	attributes := objectAttributes{
		Length:        uint32(unsafe.Sizeof(objectAttributes{})),
		RootDirectory: parent,
		ObjectName:    &objectName,
		Attributes:    objCaseInsensitive,
	}
	createOptions := uint32(createOptionOpenForBackupIntent)
	if !followReparse {
		createOptions |= createOptionOpenReparsePoint
	}
	if directory {
		createOptions |= createOptionDirectoryFile
	} else {
		createOptions |= createOptionNonDirectoryFile
	}
	// P84 T4 fix round 2: the access bits are explicit (see the access constants
	// above). A directory component is opened to be LISTED and identified, the leaf
	// of an artifact open to be READ and identified; the first bit is
	// FILE_LIST_DIRECTORY in the directory sense and FILE_READ_DATA in the file
	// sense, so one expression states both steps truthfully. GENERIC_READ is
	// deliberately never requested.
	desiredAccess := uintptr(accessListOrRead | accessReadAttributes | accessSynchronize)
	var handle syscall.Handle
	var status ioStatusBlock
	result, _, _ := procNtCreateFile.Call(
		uintptr(unsafe.Pointer(&handle)),
		desiredAccess,
		uintptr(unsafe.Pointer(&attributes)),
		uintptr(unsafe.Pointer(&status)),
		0, // AllocationSize
		0, // FileAttributes
		uintptr(fileInformationShareAll),
		uintptr(createDispositionOpen),
		uintptr(createOptions),
		0, // EaBuffer
		0, // EaLength
	)
	if result != statusSuccess {
		return 0, false, ntStatusError(result)
	}
	info, err := fileInformationOf(handle)
	if err != nil {
		_ = syscall.CloseHandle(handle)
		return 0, false, err
	}
	return handle, info.FileAttributes&syscall.FILE_ATTRIBUTE_REPARSE_POINT != 0, nil
}

// ntStatusError converts one NTSTATUS into the error a caller already understands.
// ntdll!RtlNtStatusToDosError yields the Win32 code, so ERROR_FILE_NOT_FOUND and
// ERROR_PATH_NOT_FOUND keep errors.Is(err, fs.ErrNotExist) working exactly as the
// path-based open did (same contract as the unix half). A status the mapping does
// not translate is still refused rather than read as success.
func ntStatusError(status uintptr) error {
	if status == statusSuccess {
		return nil
	}
	code, _, _ := procRtlNtStatusToDosError.Call(status)
	if code == 0 {
		// ERROR_INVALID_PARAMETER: an unmapped status is still refused, never read as
		// success through a zero errno.
		code = win32InvalidParameter
	}
	return syscall.Errno(code)
}

// isUnderDirectory reports whether candidate is a strict descendant of base by
// spelling. Both names come from the same API (handleFinalPath) in the same form,
// so a real descendant's prefix is byte-identical and the comparison is Ordinal
// (P84 T4(b)). Case-folding would admit a case-variant twin directory whose name
// only differs by case, which is exactly the shape that must be refused.
func isUnderDirectory(base, candidate string) bool {
	trimmed := strings.TrimRight(filepath.Clean(base), `\/`)
	if trimmed == "" {
		return false
	}
	prefix := trimmed + string(filepath.Separator)
	cleaned := filepath.Clean(candidate)
	if len(cleaned) <= len(prefix) || cleaned[:len(prefix)] != prefix {
		return false
	}
	rest := cleaned[len(prefix):]
	return rest != "" && rest != "." && rest != ".." && !strings.HasPrefix(rest, ".."+string(filepath.Separator))
}

// Close releases the directory handle.
func (dir *BoundDir) Close() error {
	if dir == nil || dir.handle == 0 {
		return nil
	}
	err := syscall.CloseHandle(dir.handle)
	dir.handle = 0
	return err
}

// boundDirVerify proves the handle still addresses the directory the caller
// named. The name form is compared first by exact bytes; any other spelling - a
// short (8.3) form, or a name that differs only by case - falls back to comparing
// the two directory objects, so a directory swapped between the open and this check
// is refused. P84 T4(b) removed the case-folding fast path: a differently cased
// name is precisely how a case-variant twin directory presents itself, so it is
// admitted only when the object comparison proves it is the same object. A
// directory that was already replaced when the open happened cannot be detected
// here - the handle then addresses the replacement, which is exactly what the name
// resolves to - so the check is bounded to that window by construction.
func boundDirVerify(dir *BoundDir, expected string) error {
	if dir == nil || dir.handle == 0 {
		return fmt.Errorf("evidence: nil bound directory")
	}
	real, err := dir.currentName()
	if err != nil {
		return fmt.Errorf("evidence: derive the bound directory name: %v", err)
	}
	if filepath.Clean(real) == filepath.Clean(expected) {
		return nil
	}
	same, err := dir.sameObjectAs(expected)
	if err != nil {
		return fmt.Errorf("evidence: compare the bound directory with %s: %v", expected, err)
	}
	if !same {
		return fmt.Errorf("evidence: bound directory %q is not %q", real, expected)
	}
	return nil
}

// VerifyBoundDir is the exported form of the identity re-check, kept for the
// adapters alias of the original seam.
func VerifyBoundDir(dir *BoundDir, expected string) error {
	return boundDirVerify(dir, expected)
}

// RootPath reports the filesystem name the bound directory object carries right
// now, derived from the handle. T4 (P83 fix round 6): the attempt's work directory is
// derived from this name, so a directory renamed after it was bound still executes
// inside the object the artifacts are read through. T4 (P83 fix round 7) removed its
// other use: the artifact enumeration now reads the entries through the bound handle
// (BoundDirNames) instead of globbing this path.
func (dir *BoundDir) RootPath() (string, error) {
	if dir == nil || dir.handle == 0 {
		return "", fmt.Errorf("evidence: nil bound directory")
	}
	return dir.currentName()
}

// currentName reports the current filesystem name of the bound directory object.
// It follows renames of the object and of its ancestors, which is exactly why the
// binding survives an ancestor replacement.
func (dir *BoundDir) currentName() (string, error) {
	return handleFinalPath(dir.handle)
}

// handleFinalPath resolves the current filesystem name of whatever the handle
// addresses, in the ordinary DOS form. It is shared by the directory handle and
// by the artifact handles opened under it so both names are derived the same way.
func handleFinalPath(handle syscall.Handle) (string, error) {
	for size := boundDirNameBufferStart; size <= boundDirNameBufferLimit; size *= 2 {
		buffer := make([]uint16, size)
		written, _, callErr := procGetFinalPathNameByHandleW.Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(size),
			uintptr(volumeNameDos|fileNameNormalized),
		)
		if written == 0 {
			return "", callErr
		}
		if int(written) < size {
			return stripExtendedPathPrefix(syscall.UTF16ToString(buffer[:written])), nil
		}
	}
	return "", fmt.Errorf("evidence: the handle name exceeds the %d byte buffer", boundDirNameBufferLimit)
}

// sameObjectAs reports whether the directory currently reachable at path is the
// bound object, by volume serial and file index rather than by spelling.
func (dir *BoundDir) sameObjectAs(path string) (bool, error) {
	boundInfo, err := dir.byHandleFileInformation()
	if err != nil {
		return false, err
	}
	name, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, fileInformationShareAll, nil, syscall.OPEN_EXISTING, fileFlagBackupSemantics, 0)
	if err != nil {
		return false, err
	}
	defer func() { _ = syscall.CloseHandle(handle) }()
	otherInfo, err := fileInformationOf(handle)
	if err != nil {
		return false, err
	}
	return boundInfo.VolumeSerialNumber == otherInfo.VolumeSerialNumber &&
		boundInfo.FileIndexHigh == otherInfo.FileIndexHigh &&
		boundInfo.FileIndexLow == otherInfo.FileIndexLow, nil
}

func (dir *BoundDir) byHandleFileInformation() (syscall.ByHandleFileInformation, error) {
	return fileInformationOf(dir.handle)
}

func fileInformationOf(handle syscall.Handle) (syscall.ByHandleFileInformation, error) {
	var info syscall.ByHandleFileInformation
	if err := syscall.GetFileInformationByHandle(handle, &info); err != nil {
		return syscall.ByHandleFileInformation{}, err
	}
	return info, nil
}

// stripExtendedPathPrefix removes the \\?\ / \\?\UNC\ form the API returns so the
// name can be joined with filepath and compared with the caller's path.
func stripExtendedPathPrefix(name string) string {
	const extendedUNC = `\\?\UNC\`
	if strings.HasPrefix(name, extendedUNC) {
		return `\\` + name[len(extendedUNC):]
	}
	const extended = `\\?\`
	if strings.HasPrefix(name, extended) {
		return name[len(extended):]
	}
	return name
}
