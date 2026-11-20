//go:build windows

package evidence

// P42-T027B-EVIDENCE-BIND Q2, Windows half: a directory handle opened with
// FILE_FLAG_BACKUP_SEMANTICS plus share-delete, and relative opens resolved
// through the object's current real name (GetFinalPathNameByHandle). Windows has
// no openat, so the handle-derived name is what makes the binding follow the
// directory object instead of the path.

import (
	"fmt"
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
	clean, err := cleanBoundRelativePath(rel)
	if err != nil {
		return nil, err
	}
	base, err := dir.currentName()
	if err != nil {
		return nil, fmt.Errorf("evidence: resolve the bound directory %s: %v", dir.path, err)
	}
	target := filepath.Join(base, clean)
	name, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: target, Err: err}
	}
	// T4 (P83 fix round 8): the leaf is opened as itself. A reparse point at the
	// artifact name is therefore not followed, and confirmOpenedArtifact refuses the
	// object it opened instead of reading the target it redirects to.
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, fileInformationShareAll, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL|fileFlagOpenReparsePoint, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: target, Err: err}
	}
	// P42 F1: the derived name above can already be stale by the time CreateFile
	// resolves it - the bound directory, or one of its ancestors, may be replaced
	// in between, and Windows has no openat to make the resolution atomic. So the
	// artifact is re-checked through the handle that was actually opened: its final
	// path must still sit under the bound directory's current real name, otherwise
	// the handle addresses a file from the replacement tree and is closed instead of
	// returned (fail-closed).
	if err := dir.confirmOpenedArtifact(handle, clean); err != nil {
		_ = syscall.CloseHandle(handle)
		return nil, err
	}
	return NewBoundFile(os.NewFile(uintptr(handle), target), bound, allowTruncation)
}

// BoundDirNames lists the entry names of one directory inside the bound directory
// object, sorted (P83 fix round 7, item T4b).
//
// The entries are read from the opened directory HANDLE (os.File.ReadDir enumerates
// a Windows handle through FileIdBothDirectoryInfo), so the names come from the
// directory object itself and never from a second name resolution. T4 (P83 fix
// round 8) made the admission of that handle object-grade: it must not be a reparse
// point, it must live on the bound volume, the bound directory's name must still
// address the bound object through the bound handle, and its handle-derived name
// must sit beneath that base. Ownership that cannot be proven refuses instead of
// enumerating a replacement tree. rel == "" (or ".") names the bound object itself.
func BoundDirNames(dir *BoundDir, rel string) ([]string, error) {
	if dir == nil || dir.handle == 0 {
		return nil, fmt.Errorf("evidence: nil bound directory")
	}
	clean := "."
	if strings.TrimSpace(rel) != "" {
		cleaned, err := cleanBoundRelativePath(rel)
		if err != nil {
			return nil, err
		}
		clean = cleaned
	}
	base, err := dir.currentName()
	if err != nil {
		return nil, fmt.Errorf("evidence: resolve the bound directory %s: %v", dir.path, err)
	}
	target := base
	if clean != "." {
		target = filepath.Join(base, clean)
	}
	name, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: target, Err: err}
	}
	// T4b (P83 fix round 9): a descendant directory is opened as ITSELF, exactly as
	// openBoundDirFile opens the artifact leaf, so a junction at the enumerated name
	// is not followed to its target. Without the flag the handle describes the
	// target, which carries no reparse attribute, so the classification in
	// confirmOpenedArtifact could never fire and a junction pointing inside the bound
	// tree was admitted by the containment comparison alone. rel "." (the bound
	// object itself) keeps the plain open: that object was already admitted by
	// OpenBoundDir, and confirmEnumeratedDirectory proves by volume and file index
	// that the name still addresses it.
	attributes := uint32(fileFlagBackupSemantics)
	if clean != "." {
		attributes |= fileFlagOpenReparsePoint
	}
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, fileInformationShareAll, nil, syscall.OPEN_EXISTING, attributes, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: target, Err: err}
	}
	if err := dir.confirmEnumeratedDirectory(handle, clean); err != nil {
		_ = syscall.CloseHandle(handle)
		return nil, err
	}
	file := os.NewFile(uintptr(handle), target)
	defer func() { _ = file.Close() }()
	entries, readErr := file.ReadDir(-1)
	if readErr != nil {
		return nil, readErr
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	return names, nil
}

// confirmEnumeratedDirectory proves the directory handle that will be enumerated
// still addresses the bound object (rel ".") or a strict descendant of it.
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

// isUnderDirectory reports whether candidate is a strict descendant of base by
// spelling. Both names come from the same API in the same form, so a case-folding
// comparison is enough and no volume-normalizing rewrite is needed.
func isUnderDirectory(base, candidate string) bool {
	trimmed := strings.TrimRight(filepath.Clean(base), `\/`)
	if trimmed == "" {
		return false
	}
	prefix := trimmed + string(filepath.Separator)
	cleaned := filepath.Clean(candidate)
	if len(cleaned) <= len(prefix) || !strings.EqualFold(cleaned[:len(prefix)], prefix) {
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
// named. The name form is compared first; a short (8.3) or differently cased
// spelling falls back to comparing the two directory objects, so a directory
// swapped between the open and this check is refused and only a spelling
// difference is accepted. A directory that was already replaced when the open
// happened cannot be detected here - the handle then addresses the replacement,
// which is exactly what the name resolves to - so the check is bounded to that
// window by construction.
func boundDirVerify(dir *BoundDir, expected string) error {
	if dir == nil || dir.handle == 0 {
		return fmt.Errorf("evidence: nil bound directory")
	}
	real, err := dir.currentName()
	if err != nil {
		return fmt.Errorf("evidence: derive the bound directory name: %v", err)
	}
	if strings.EqualFold(filepath.Clean(real), filepath.Clean(expected)) {
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
