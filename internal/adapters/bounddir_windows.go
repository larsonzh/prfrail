//go:build windows

package adapters

// P42-T027B-EVIDENCE-BIND Q2, Windows half: a directory handle opened with
// FILE_FLAG_BACKUP_SEMANTICS plus share-delete, and relative opens resolved
// through the object's current real name (GetFinalPathNameByHandle). Windows has
// no openat, so the handle-derived name is what makes the binding follow the
// directory object instead of the path.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	// fileFlagBackupSemantics opens a directory handle instead of refusing it.
	fileFlagBackupSemantics = 0x02000000
	// fileInformationShareAll is FILE_SHARE_READ|WRITE|DELETE, the share mode that
	// keeps the bound object renamable and removable while it stays bound.
	fileInformationShareAll = syscall.FILE_SHARE_READ | syscall.FILE_SHARE_WRITE | syscall.FILE_SHARE_DELETE
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
		return nil, fmt.Errorf("adapters: %s is not a directory", path)
	}
	return &BoundDir{handle: handle, path: path}, nil
}

// OpenBoundDirFile opens one artifact relative to the bound directory. The
// returned handle is the artifact's identity: the caller reads through it and
// never resolves the name again.
func OpenBoundDirFile(dir *BoundDir, rel string, bound int64) (*evidence.BoundFile, error) {
	return openBoundDirFile(dir, rel, bound, false)
}

// openBoundDirFile is the class-aware form: allowTruncation keeps the log and
// event class readable past its bound instead of refusing the open.
func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*evidence.BoundFile, error) {
	if dir == nil || dir.handle == 0 {
		return nil, fmt.Errorf("adapters: nil bound directory")
	}
	clean, err := cleanBoundRelativePath(rel)
	if err != nil {
		return nil, err
	}
	base, err := dir.currentName()
	if err != nil {
		return nil, fmt.Errorf("adapters: resolve the bound directory %s: %v", dir.path, err)
	}
	target := filepath.Join(base, clean)
	name, err := syscall.UTF16PtrFromString(target)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: target, Err: err}
	}
	handle, err := syscall.CreateFile(name, syscall.GENERIC_READ, fileInformationShareAll, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
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
	return evidence.NewBoundFile(os.NewFile(uintptr(handle), target), bound, allowTruncation)
}

// confirmOpenedArtifact proves the artifact just opened still belongs to the bound
// directory object. The comparison is handle-to-handle: the file handle's final
// path against the directory handle's current real name, both derived by the API
// rather than by re-walking the caller's path, so a rename of the object or of its
// ancestors is followed and only a replacement is refused.
func (dir *BoundDir) confirmOpenedArtifact(handle syscall.Handle, rel string) error {
	final, err := handleFinalPath(handle)
	if err != nil {
		return fmt.Errorf("adapters: resolve the opened artifact %s: %v", rel, err)
	}
	base, err := dir.currentName()
	if err != nil {
		return fmt.Errorf("adapters: resolve the bound directory %s: %v", dir.path, err)
	}
	if !isUnderDirectory(base, final) {
		return fmt.Errorf("adapters: artifact %s resolved to %s, which is outside the bound directory %s", rel, final, base)
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
		return fmt.Errorf("adapters: nil bound directory")
	}
	real, err := dir.currentName()
	if err != nil {
		return fmt.Errorf("adapters: derive the bound directory name: %v", err)
	}
	if strings.EqualFold(filepath.Clean(real), filepath.Clean(expected)) {
		return nil
	}
	same, err := dir.sameObjectAs(expected)
	if err != nil {
		return fmt.Errorf("adapters: compare the bound directory with %s: %v", expected, err)
	}
	if !same {
		return fmt.Errorf("adapters: bound directory %q is not %q", real, expected)
	}
	return nil
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
	return "", fmt.Errorf("adapters: the handle name exceeds the %d byte buffer", boundDirNameBufferLimit)
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
