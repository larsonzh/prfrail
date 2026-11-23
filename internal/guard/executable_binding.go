//go:build windows || linux

package guard

// P83 fix round 7, item T3: the verified executable object must be the launched
// object.
//
// The earlier shape opened and hashed the executable, closed it, and then let the
// platform create the child from the PATH - two independent resolutions of one
// name, with a window in between that a replacer owns. The binding below removes
// the window per platform instead of narrowing it:
//
//   - Windows: the executable is opened with a share mode that denies write and
//     delete, and that handle stays open across the whole verify -> spawn
//     sequence; the image cannot be rewritten or replaced while it is held.
//   - Linux: the child is created from the very descriptor that was hashed
//     (descriptor exec), so the name is never resolved again; a host that cannot
//     serve the descriptor refuses rather than falling back to the path.
//   - Other platforms (darwin): no descriptor exec is reachable from the standard
//     library, so an approved digest refuses outright.
//
// The digest is computed through the same handle that the spawn uses, so a
// mismatch is a refusal with no child and no footprint.

import (
	"fmt"
	"os"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// managedExecutableMaxBytes bounds the approved executable the guard digests. A
// larger image is refused instead of digested: the digest has to describe the
// whole object the child is created from.
const managedExecutableMaxBytes = int64(64 << 20)

// ErrExecutableBinding lives in executable_binding_error.go so the unsupported
// platform can refuse with it without this helper being compiled.

// approvedExecutable is an approved executable held by the single handle that was
// verified. The handle stays open until Close, which is what keeps the verified
// object the launched object for the whole verify -> spawn sequence.
type approvedExecutable struct {
	// file is the same handle bound holds; it is kept so the platform start can
	// reach the descriptor it has to execute.
	file  *os.File
	bound *evidence.BoundFile
}

// Close releases the verified handle. It is idempotent.
func (approved *approvedExecutable) Close() error {
	if approved == nil || approved.bound == nil {
		return nil
	}
	return approved.bound.Close()
}

// verifyApprovedExecutable opens path once, digests THAT handle and compares it
// with the approval. The returned handle is the object the platform start must
// create the child from; a caller that closes it before the spawn reopens the
// window this binding exists to remove.
func verifyApprovedExecutable(path string, spec *ProcessSpec) (*approvedExecutable, error) {
	if spec.ExecutableDigest == "" || spec.ExecutableDigestDomain == "" {
		return nil, fmt.Errorf("%w: the approved executable carries no digest or no digest domain", ErrExecutableBinding)
	}
	file, err := openApprovedExecutable(path)
	if err != nil {
		return nil, fmt.Errorf("%w: open the approved executable: %v", ErrExecutableBinding, err)
	}
	// The bound file owns the handle from here on: a refused artifact (a
	// directory, a non-regular object, an oversized image) is closed by the
	// constructor, so no failure path below leaves a handle behind.
	bound, err := evidence.NewBoundFile(file, managedExecutableMaxBytes, false)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrExecutableBinding, err)
	}
	actual, err := bound.Hash(spec.ExecutableDigestDomain)
	if err != nil {
		_ = bound.Close()
		return nil, fmt.Errorf("%w: digest the approved executable: %v", ErrExecutableBinding, err)
	}
	if actual != spec.ExecutableDigest {
		_ = bound.Close()
		return nil, fmt.Errorf("%w: the opened executable is not the approved object", ErrExecutableBinding)
	}
	// P85 F05: when the caller pins the object identity too, the SAME handle that
	// was digested must carry the approved file id. A digest comparison cannot see
	// a same-digest-different-inode copy; the identity comparison can.
	if spec.ExecutableFileID != "" {
		actualID, idErr := bound.FileID()
		if idErr != nil {
			_ = bound.Close()
			return nil, fmt.Errorf("%w: read the identity of the approved executable: %v", ErrExecutableBinding, idErr)
		}
		if actualID != spec.ExecutableFileID {
			_ = bound.Close()
			return nil, fmt.Errorf("%w: the opened executable is not the approved object identity", ErrExecutableBinding)
		}
	}
	return &approvedExecutable{file: file, bound: bound}, nil
}
