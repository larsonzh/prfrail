//go:build darwin

package evidence

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
	"fmt"
)

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
func OpenBoundDirFile(dir *BoundDir, rel string, bound int64) (*BoundFile, error) {
	return openBoundDirFile(dir, rel, bound, false)
}

// OpenBoundDirNoFollow refuses for the same reason as OpenBoundDir: darwin has
// no openat2 and a name-based walk cannot keep the containment guarantee.
func OpenBoundDirNoFollow(path string) (*BoundDir, error) {
	return nil, fmt.Errorf("%w: %s", ErrBoundDirUnsupported, path)
}

// OpenBoundDirFileClass is the class-aware form: allowTruncation keeps the log and
// event class readable past its bound instead of refusing the open.
func OpenBoundDirFileClass(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*BoundFile, error) {
	return openBoundDirFile(dir, rel, bound, allowTruncation)
}

// openBoundDirFile refuses for the same reason as OpenBoundDir.
func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*BoundFile, error) {
	_ = dir
	_ = bound
	_ = allowTruncation
	return nil, fmt.Errorf("%w: %s", ErrBoundDirUnsupported, rel)
}

// RootPath reports the name the placeholder was created with. T4 (P83 fix round
// 6): darwin never binds, so no enumeration can run here and the name is only
// reported for symmetry with the other platforms.
func (dir *BoundDir) RootPath() (string, error) {
	if dir == nil {
		return "", fmt.Errorf("evidence: nil bound directory")
	}
	return dir.path, nil
}

// BoundDirNames refuses for the same reason as OpenBoundDir (P83 fix round 7, item
// T4b): there is no descriptor to enumerate through, and a path-based listing is
// exactly the shape the binding exists to remove.
func BoundDirNames(dir *BoundDir, rel string) ([]string, error) {
	_ = dir
	return nil, fmt.Errorf("%w: %s", ErrBoundDirUnsupported, rel)
}

// Close releases nothing: no descriptor was opened.
func (dir *BoundDir) Close() error {
	_ = dir
	return nil
}

// boundDirVerify refuses: an unbound directory proves nothing about the path the
// caller named.
func boundDirVerify(dir *BoundDir, expected FileID) error {
	_ = dir
	return fmt.Errorf("%w: %s", ErrBoundDirUnsupported, expected)
}

// VerifyBoundDir is the exported form of the identity re-check. darwin never
// binds, so every identity assertion refuses.
func VerifyBoundDir(dir *BoundDir, expected FileID) error {
	return boundDirVerify(dir, expected)
}

// OpenBoundDirVerified refuses for the same reason as OpenBoundDirNoFollow: there
// is no descriptor to take an identity from (P85 F02/F06/F20).
func OpenBoundDirVerified(path string) (*BoundDir, FileID, error) {
	return nil, "", fmt.Errorf("%w: %s", ErrBoundDirUnsupported, path)
}

// FileID refuses: the placeholder wraps no descriptor.
func (dir *BoundDir) FileID() (FileID, error) {
	if dir == nil {
		return "", fmt.Errorf("evidence: nil bound directory")
	}
	return "", fmt.Errorf("%w: %s", ErrBoundDirUnsupported, dir.path)
}

// SyncDirectory refuses for the same reason as the other entry points: an
// unbound directory cannot be flushed (P85 F10).
func (dir *BoundDir) SyncDirectory() error {
	if dir == nil {
		return fmt.Errorf("evidence: nil bound directory")
	}
	return fmt.Errorf("%w: %s", ErrBoundDirUnsupported, dir.path)
}
