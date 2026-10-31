package evidence

// P42-T027B-EVIDENCE-BIND: the C3 bound-read primitive.
//
// docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md
// section 2 "Q1 unified read primitive" requires one open per evidence artifact:
// the bound decision, the bytes, the truncation flag, the digest and the identity
// must all describe the same opened file object, never a path that was resolved a
// second time. A caller that stats a path and then opens it lets an artifact be
// swapped between the two calls; a caller that opens and then reads by path lets
// it be swapped between the open and the read. BoundFile removes both shapes by
// holding the handle and never resolving the name again.
//
// The primitive deliberately has no path re-resolution of any kind: ReadAll and
// Hash only Seek and Read the handle captured by the constructor.

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	// ErrBoundArtifactExceedsBound reports an artifact that grew past the bound the
	// open-time stat accepted. It is an error, never a silent short read, because a
	// digest computed over a truncated artifact would describe the wrong bytes.
	ErrBoundArtifactExceedsBound = errors.New("bound artifact exceeds its byte bound")
	// ErrBoundFileClosed reports use of a BoundFile after Close.
	ErrBoundFileClosed = errors.New("bound file is closed")
	// ErrBoundFileNotRegular reports a bound artifact that is not a regular file.
	ErrBoundFileNotRegular = errors.New("bound artifact is not a regular file")
)

// BoundFile is one evidence artifact held open by a single handle. Every reading
// surface (bytes, truncation flag, digest, identity) is derived from that handle,
// so a path swapped after the open cannot redirect a read.
//
// BoundFile owns the handle it was built from: the constructor closes it when the
// artifact is refused, and Close releases it on the success path.
type BoundFile struct {
	file   *os.File
	info   os.FileInfo
	bound  int64
	closed bool
}

// OpenBoundFile opens path once and binds it to bound for the record and pid
// class: an artifact that already exceeds bound is refused at open time and
// yields no handle and no data. A missing artifact keeps fs.ErrNotExist visible
// through errors.Is so callers keep their "absent is not an error" mapping.
func OpenBoundFile(path string, bound int64) (*BoundFile, error) {
	return openBoundFile(path, bound, false)
}

// OpenBoundFileTruncating opens path once for the log and event class: an
// artifact that already exceeds bound is still opened and degrades to a
// truncated read (ReadAll reports truncated at bound+1). It exists so the record
// class refusal never drifts the log class truncation contract.
func OpenBoundFileTruncating(path string, bound int64) (*BoundFile, error) {
	return openBoundFile(path, bound, true)
}

func openBoundFile(path string, bound int64, allowTruncation bool) (*BoundFile, error) {
	// The platform opener supplies the handle sharing the class needs: on Windows
	// an open artifact must stay renamable and removable, otherwise the handle
	// itself would hide the replacement it is supposed to observe.
	file, err := openBoundFileHandle(path)
	if err != nil {
		// The opener reports a *PathError, and errors.Is(err, fs.ErrNotExist) keeps
		// working through it: no remapping may hide the absence contract.
		return nil, err
	}
	return NewBoundFile(file, bound, allowTruncation)
}

// NewBoundFile adopts an already open handle (a directory-relative open, or a
// caller-owned open) and applies the same type, bound and identity checks
// OpenBoundFile applies. It takes ownership of file: a refused artifact is
// closed before the error returns.
//
// allowTruncation selects the class contract: false refuses an artifact whose
// open-time size exceeds bound (record and pid class), true keeps it readable and
// lets ReadAll report the truncation (log and event class).
func NewBoundFile(file *os.File, bound int64, allowTruncation bool) (*BoundFile, error) {
	if file == nil {
		return nil, fmt.Errorf("%w: an open handle is required", ErrBoundFileNotRegular)
	}
	if bound < 0 {
		_ = file.Close()
		return nil, fmt.Errorf("%w: negative byte bound %d", ErrBoundArtifactExceedsBound, bound)
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if info.IsDir() {
		_ = file.Close()
		return nil, fmt.Errorf("%w: %s is a directory", ErrBoundFileNotRegular, info.Name())
	}
	if !info.Mode().IsRegular() {
		_ = file.Close()
		return nil, fmt.Errorf("%w: %s has mode %s", ErrBoundFileNotRegular, info.Name(), info.Mode())
	}
	if !allowTruncation && info.Size() > bound {
		_ = file.Close()
		return nil, fmt.Errorf("%w: %s is %d bytes at the %d byte bound", ErrBoundArtifactExceedsBound, info.Name(), info.Size(), bound)
	}
	return &BoundFile{file: file, info: info, bound: bound}, nil
}

// ReadAll returns the artifact bytes read through the bound handle. truncated
// reports that the artifact continued past bound; in that case the returned
// prefix stops at the bound, exactly as the log class has always degraded.
//
// A read error discards every byte read so far: a partial read never masquerades
// as data.
func (bound *BoundFile) ReadAll() ([]byte, bool, error) {
	if err := bound.checkReadable(); err != nil {
		return nil, false, err
	}
	if _, err := bound.file.Seek(0, io.SeekStart); err != nil {
		return nil, false, err
	}
	data, err := io.ReadAll(io.LimitReader(bound.file, readLimit(bound.bound)))
	if err != nil {
		return nil, false, err
	}
	if int64(len(data)) > bound.bound {
		return data[:bound.bound], true, nil
	}
	return data, false, nil
}

// Hash returns the domain-separated digest of the artifact, computed in one pass
// over the same handle ReadAll uses. A truncated or unreadable artifact yields an
// error and no digest: a digest that covered only a prefix would name bytes the
// artifact does not have.
func (bound *BoundFile) Hash(domain string) (string, error) {
	if err := bound.checkReadable(); err != nil {
		return "", err
	}
	if _, err := bound.file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	hasher := sha256.New()
	hasher.Write([]byte(domain))
	written, err := io.Copy(hasher, io.LimitReader(bound.file, readLimit(bound.bound)))
	if err != nil {
		return "", err
	}
	if written > bound.bound {
		return "", fmt.Errorf("%w: digest refused for %s", ErrBoundArtifactExceedsBound, bound.info.Name())
	}
	return "sha256:" + hex.EncodeToString(hasher.Sum(nil)), nil
}

// Identity returns the size, modification time and mode observed on the handle
// when it was opened. It never re-stats the path, so it describes the object the
// reads describe.
func (bound *BoundFile) Identity() (size int64, mod time.Time, mode os.FileMode) {
	if bound == nil || bound.info == nil {
		return 0, time.Time{}, 0
	}
	return bound.info.Size(), bound.info.ModTime(), bound.info.Mode()
}

// Close releases the handle. It is idempotent; ReadAll and Hash report
// ErrBoundFileClosed afterwards instead of touching a released handle.
func (bound *BoundFile) Close() error {
	if bound == nil || bound.file == nil || bound.closed {
		return nil
	}
	bound.closed = true
	return bound.file.Close()
}

func (bound *BoundFile) checkReadable() error {
	if bound == nil || bound.file == nil {
		return fmt.Errorf("%w: nil bound file", ErrBoundFileClosed)
	}
	if bound.closed {
		return fmt.Errorf("%w: %s", ErrBoundFileClosed, bound.info.Name())
	}
	return nil
}

// readLimit is the bound-plus-one read that makes the truncation decision
// observable without unbounded buffering. The sentinel byte cannot be
// represented at the int64 ceiling, so the ceiling keeps the bound itself and
// the caller reports no truncation instead of wrapping to a negative limit.
func readLimit(bound int64) int64 {
	limit := bound + 1
	if limit <= bound {
		return bound
	}
	return limit
}
