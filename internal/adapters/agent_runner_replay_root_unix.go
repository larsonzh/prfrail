//go:build !windows

package adapters

import (
	"errors"
	"io/fs"
	"os"
)

// replayPathComponentUnsafe reports whether an existing path component is a
// symlink. Missing components are reported as safe because only already-existing
// components can redirect reads or writes.
func replayPathComponentUnsafe(path string) (bool, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return info.Mode()&os.ModeSymlink != 0, nil
}
