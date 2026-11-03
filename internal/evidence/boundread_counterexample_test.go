package evidence

// P52-S1 ACCEPTANCE-SECURITY-COUNTEREXAMPLES: primitive-layer counterexamples.
//
// The three sentinels in boundread.go (ErrBoundArtifactExceedsBound,
// ErrBoundFileClosed, ErrBoundFileNotRegular) are acceptance evidence: a caller
// distinguishes "absent", "over bound" and "unusable handle" by errors.Is. Each
// test below drives one refusal shape and pins the sentinel, never just err != nil.
// Every fixture is built in t.TempDir() so no repository fixture is written.

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// TestBoundFileReadAfterCloseRefused pins the closed-handle contract: once Close
// released the handle, both reading surfaces must report ErrBoundFileClosed
// instead of touching a released descriptor.
func TestBoundFileReadAfterCloseRefused(t *testing.T) {
	path := filepath.Join(t.TempDir(), "artifact.jsonl")
	payload := []byte("0123456789abcdef")
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write the artifact: %v", err)
	}

	bound, err := OpenBoundFile(path, 64)
	if err != nil {
		t.Fatalf("OpenBoundFile: %v", err)
	}
	if data, truncated, err := bound.ReadAll(); err != nil || truncated || string(data) != string(payload) {
		t.Fatalf("first read must succeed whole: len=%d truncated=%v err=%v", len(data), truncated, err)
	}
	if err := bound.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	if data, truncated, err := bound.ReadAll(); !errors.Is(err, ErrBoundFileClosed) {
		t.Fatalf("ReadAll after Close must be ErrBoundFileClosed: len=%d truncated=%v err=%v", len(data), truncated, err)
	}
	if digest, err := bound.Hash("probe"); !errors.Is(err, ErrBoundFileClosed) {
		t.Fatalf("Hash after Close must be ErrBoundFileClosed: digest=%q err=%v", digest, err)
	}
}

// TestBoundFileHashRefusesGrowthPastBound pins the growth shape: the open-time
// bound decision and the read share the handle, so a byte appended by another
// writer degrades the read to a truncation flag and makes the digest refuse.
func TestBoundFileHashRefusesGrowthPastBound(t *testing.T) {
	path := filepath.Join(t.TempDir(), "record.jsonl")
	if err := os.WriteFile(path, []byte("abcd"), 0o644); err != nil {
		t.Fatalf("write the artifact: %v", err)
	}

	bound, err := OpenBoundFile(path, 4)
	if err != nil {
		t.Fatalf("OpenBoundFile at the exact bound: %v", err)
	}
	defer func() { _ = bound.Close() }()

	writer, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		t.Skipf("the platform refused a second writer on the bound artifact, growth is not testable here: %v", err)
	}
	if _, err := writer.Write([]byte("e")); err != nil {
		_ = writer.Close()
		t.Fatalf("append past the bound: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close the second writer: %v", err)
	}

	data, truncated, err := bound.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll over a grown artifact must degrade, not fail: %v", err)
	}
	if !truncated || len(data) != 4 || string(data) != "abcd" {
		t.Fatalf("a grown artifact must read as the bound prefix with truncated=true: len=%d truncated=%v data=%q", len(data), truncated, data)
	}
	if digest, err := bound.Hash("probe"); !errors.Is(err, ErrBoundArtifactExceedsBound) {
		t.Fatalf("Hash over a grown artifact must be ErrBoundArtifactExceedsBound: digest=%q err=%v", digest, err)
	}
}

// TestBoundFileNewRefusesNilAndNegativeBound pins the two constructor guards: a
// nil handle is not a regular artifact, and a negative bound can never be met.
func TestBoundFileNewRefusesNilAndNegativeBound(t *testing.T) {
	if bound, err := NewBoundFile(nil, 0, false); !errors.Is(err, ErrBoundFileNotRegular) {
		t.Fatalf("NewBoundFile(nil) must be ErrBoundFileNotRegular: bound=%v err=%v", bound, err)
	}

	path := filepath.Join(t.TempDir(), "artifact.jsonl")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatalf("write the artifact: %v", err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open the artifact: %v", err)
	}
	if bound, err := NewBoundFile(file, -1, false); !errors.Is(err, ErrBoundArtifactExceedsBound) {
		t.Fatalf("NewBoundFile with a negative bound must be ErrBoundArtifactExceedsBound: bound=%v err=%v", bound, err)
	}
}

// TestBoundFileOpenRefusesDirectoryWithSentinel pins the type class at the
// open-bound entry point: a directory is never a readable artifact.
//
// Platform note: the sentinel has two legitimate shapes. Where the handle lets a
// directory be opened (unix), NewBoundFile refuses it with ErrBoundFileNotRegular;
// on Windows the opener itself fails with ACCESS_DENIED before the type check can
// run. Both are refusals and both are pinned with errors.Is - a bare err != nil
// would accept an unrelated failure.
func TestBoundFileOpenRefusesDirectoryWithSentinel(t *testing.T) {
	dir := t.TempDir()

	bound, err := OpenBoundFile(dir, 64)
	if err == nil {
		_ = bound.Close()
		t.Fatal("opening a directory as a bounded artifact must be refused")
	}
	if bound != nil {
		t.Fatalf("a refused open must yield no handle, got %p", bound)
	}
	switch {
	case errors.Is(err, ErrBoundFileNotRegular):
		// Handle opened, type check refused it (the unix shape).
	case errors.Is(err, fs.ErrPermission):
		// The platform opener refused the directory itself (the Windows shape).
	default:
		t.Fatalf("a directory must be refused as ErrBoundFileNotRegular (or the platform open refusal): %v", err)
	}
}

// TestBoundFileOpenRefusesOverBoundWithSentinel pins the over-bound refusal at
// open time: an artifact that already exceeds its bound yields no handle and no
// data, because a digest over a prefix would name bytes the artifact lacks.
func TestBoundFileOpenRefusesOverBoundWithSentinel(t *testing.T) {
	path := filepath.Join(t.TempDir(), "record.jsonl")
	if err := os.WriteFile(path, []byte("123456789"), 0o644); err != nil {
		t.Fatalf("write the artifact: %v", err)
	}

	bound, err := OpenBoundFile(path, 8)
	if !errors.Is(err, ErrBoundArtifactExceedsBound) {
		t.Fatalf("a 9 byte artifact at an 8 byte bound must be ErrBoundArtifactExceedsBound: bound=%v err=%v", bound, err)
	}
	if bound != nil {
		t.Fatalf("a refused open must yield no handle, got %p", bound)
	}
}
