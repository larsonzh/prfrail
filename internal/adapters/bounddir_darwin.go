//go:build darwin

package adapters

// P42-T027B-EVIDENCE-BIND Q2, darwin half: the directory-binding support matrix is
// windows and linux (validation report D7). The standard library exposes no
// syscall.Openat on darwin, so the unix implementation cannot be compiled there.
//
// The alternative - falling back to a path-based open - is exactly the shape the
// binding exists to remove, so darwin refuses instead: every entry point reports
// ErrBoundDirUnsupported and no artifact is ever read through a name that could be
// swapped. The package still builds on darwin so the rest of the runner can be
// inspected; the run itself cannot proceed.

import (
	"errors"
	"fmt"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// ErrBoundDirUnsupported reports that directory binding is unavailable on this
// platform.
var ErrBoundDirUnsupported = errors.New("adapters: directory binding is supported on windows and linux only")

// BoundDir is the darwin placeholder. It never wraps a descriptor because no
// directory can be bound here.
type BoundDir struct {
	path string
}

// OpenBoundDir refuses: there is no openat(2) to bind the directory with.
func OpenBoundDir(path string) (*BoundDir, error) {
	return nil, fmt.Errorf("%w: %s", ErrBoundDirUnsupported, path)
}

// OpenBoundDirFile refuses for the same reason as OpenBoundDir.
func OpenBoundDirFile(dir *BoundDir, rel string, bound int64) (*evidence.BoundFile, error) {
	return openBoundDirFile(dir, rel, bound, false)
}

// openBoundDirFile refuses for the same reason as OpenBoundDir.
func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*evidence.BoundFile, error) {
	_ = bound
	_ = allowTruncation
	return nil, fmt.Errorf("%w: %s", ErrBoundDirUnsupported, rel)
}

// Close releases nothing: no descriptor was opened.
func (dir *BoundDir) Close() error {
	_ = dir
	return nil
}

// boundDirVerify refuses: an unbound directory proves nothing about the path the
// caller named.
func boundDirVerify(dir *BoundDir, expected string) error {
	_ = dir
	return fmt.Errorf("%w: %s", ErrBoundDirUnsupported, expected)
}
