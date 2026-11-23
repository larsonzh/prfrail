package adapters

// P83 SECURITY-FINDINGS-REMEDIATION (X-1): the BoundDir primitive moved to
// internal/evidence so the gates runner can bind an artifact workspace handle
// (F1) without inverting the adapters <- gates import direction. The names below
// are thin aliases, not a second implementation: every existing caller in this
// package and every existing BoundDir test keeps compiling against the same
// surface, and the platform mechanics live in exactly one place.
//
// The unexported wrappers exist only because the package's own historical call
// sites and the residual tests spell these helpers in lower case; they forward
// unchanged and never re-implement the guard.

import "github.com/larsonzh/prfrail/internal/evidence"

// BoundDir is one open directory object. It owns the descriptor until Close.
type BoundDir = evidence.BoundDir

var (
	// OpenBoundDir opens path as a directory descriptor or handle.
	OpenBoundDir = evidence.OpenBoundDir
	// OpenBoundDirFile opens one artifact relative to the bound directory.
	OpenBoundDirFile = evidence.OpenBoundDirFile
	// ErrBoundDirUnsupported reports a platform without directory binding.
	ErrBoundDirUnsupported = evidence.ErrBoundDirUnsupported
)

// openBoundDirFile forwards the class-aware open (allowTruncation selects the log
// and event class).
func openBoundDirFile(dir *BoundDir, rel string, bound int64, allowTruncation bool) (*evidence.BoundFile, error) {
	return evidence.OpenBoundDirFileClass(dir, rel, bound, allowTruncation)
}

// boundDirVerify forwards the directory identity re-check. P85 F02/F06/F20: the
// expected value is the identity the binding produced (evidence.FileID), not a
// path, so the check compares objects instead of two resolutions of one name.
func boundDirVerify(dir *BoundDir, expected evidence.FileID) error {
	return evidence.VerifyBoundDir(dir, expected)
}
