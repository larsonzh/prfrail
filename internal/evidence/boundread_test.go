package evidence

// P42-T027B-EVIDENCE-BIND ③ test-first (RED): these tests pin the bound-read
// primitive of the C3/N2 design (docs/validation/evidence/P42-T027B-EVIDENCE-BIND/
// 01-arch-design/review.md, section 2 "Q1 unified read primitive"). The API they
// reference does not exist yet on purpose: the package must fail to compile until
// the implementation lands.
//
// Pinned invariants (design section 4):
//   I1 same-handle read of the opened file object, never the replaced path.
//   I2 over-bound or non-regular artifacts are refused at open time, with no data.
//   I4 identity, length and digest all describe the same opened handle.
//
// Note (registered as not-pinned-by-design): the design lists OpenBoundFile,
// ReadAll, Hash and Identity only. It does not define a Close method on BoundFile,
// so these tests never call one; the closed-handle lifetime is pinned for the
// adapter-owned directory handles by CloseBoundDirs instead.

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

const boundFileTestDomain = "proofrail:test:bound-file:1\n"

// TestBoundFileReadsOriginalAfterPathReplacement pins I1: once opened, the read
// follows the file object, not the path. A path replaced after open must not be
// able to swap in attacker bytes.
func TestBoundFileReadsOriginalAfterPathReplacement(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "artifact.jsonl")
	const bound int64 = 1 << 20
	if err := os.WriteFile(path, []byte("original"), 0o644); err != nil {
		t.Fatalf("write original artifact: %v", err)
	}
	// Computed before the handle exists, so it is independent of the primitive.
	wantHash := Digest(boundFileTestDomain, []byte("original"))

	boundFile, err := OpenBoundFile(path, bound)
	if err != nil {
		t.Fatalf("OpenBoundFile: %v", err)
	}

	// Replace the path with a hostile file: rename the bound file away, then
	// publish a different file under the same name.
	if err := os.Rename(path, path+".bak"); err != nil {
		t.Fatalf("rename bound artifact away (the handle must be share-delete): %v", err)
	}
	if err := os.WriteFile(path, []byte("attacker"), 0o644); err != nil {
		t.Fatalf("write replacement artifact: %v", err)
	}

	data, truncated, err := boundFile.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on the bound handle: %v", err)
	}
	if string(data) != "original" {
		t.Fatalf("bound read must return the opened file object's bytes, got %q", data)
	}
	if truncated {
		t.Fatalf("a %d byte read under a %d byte bound must not be truncated", len(data), bound)
	}

	hash, err := boundFile.Hash(boundFileTestDomain)
	if err != nil {
		t.Fatalf("Hash on the bound handle: %v", err)
	}
	if hash != wantHash {
		t.Fatalf("digest must cover the opened file object: got %s want %s", hash, wantHash)
	}

	// I4 folded in: identity, length and digest describe one object.
	size, _, mode := boundFile.Identity()
	if size != int64(len(data)) {
		t.Fatalf("identity size %d must match the %d bytes read", size, len(data))
	}
	if !mode.IsRegular() {
		t.Fatalf("identity mode must be regular, got %v", mode)
	}
}

// TestBoundFileRefusesOverBoundOnOpen pins I2 for the record/PID class: an
// artifact that already exceeds its bound is refused at open time and yields no
// data (design section 2 "failure and partial-read semantics").
func TestBoundFileRefusesOverBoundOnOpen(t *testing.T) {
	dir := t.TempDir()
	const bound int64 = 8

	over := filepath.Join(dir, "record.jsonl")
	if err := os.WriteFile(over, bytes.Repeat([]byte("x"), int(bound)+1), 0o644); err != nil {
		t.Fatalf("write over-bound artifact: %v", err)
	}
	boundFile, err := OpenBoundFile(over, bound)
	if err == nil {
		t.Fatalf("a %d byte artifact must be refused at the %d byte bound instead of read", bound+1, bound)
	}
	if boundFile != nil {
		t.Fatalf("a refused open must yield no handle and no data, got %p", boundFile)
	}

	// Boundary pin: exactly bound bytes is legal and untruncated, so the refusal
	// threshold is bound+1 and not "any artifact near the bound".
	exact := filepath.Join(dir, "exact.jsonl")
	payload := bytes.Repeat([]byte("y"), int(bound))
	if err := os.WriteFile(exact, payload, 0o644); err != nil {
		t.Fatalf("write bound-sized artifact: %v", err)
	}
	okFile, err := OpenBoundFile(exact, bound)
	if err != nil {
		t.Fatalf("an artifact of exactly the bound must open: %v", err)
	}
	data, truncated, err := okFile.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on the bound-sized artifact: %v", err)
	}
	if !bytes.Equal(data, payload) || truncated {
		t.Fatalf("bound-sized artifact must read whole and untruncated, got %d bytes truncated=%v", len(data), truncated)
	}
}

// TestBoundFileRefusesDirectory pins I2 for the type class: a directory is not a
// readable artifact, even on platforms where opening one succeeds.
func TestBoundFileRefusesDirectory(t *testing.T) {
	dir := t.TempDir()
	boundFile, err := OpenBoundFile(dir, 1<<20)
	if err == nil {
		t.Fatal("opening a directory as a bounded artifact must fail")
	}
	if boundFile != nil {
		t.Fatalf("a refused open must yield no handle, got %p", boundFile)
	}
}

// TestBoundFileMissingPreservesNotExist pins the mapping every replacement site
// relies on: the callers turn a missing artifact into "absent, not an error", so
// the primitive must keep fs.ErrNotExist visible through errors.Is.
func TestBoundFileMissingPreservesNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "absent.jsonl")
	boundFile, err := OpenBoundFile(path, 1<<20)
	if err == nil {
		t.Fatal("opening a missing artifact must fail")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a missing artifact must stay fs.ErrNotExist for the callers, got %v", err)
	}
	if boundFile != nil {
		t.Fatalf("a refused open must yield no handle, got %p", boundFile)
	}
}

// TestBoundFileIdentityMatchesRead pins I4 without any replacement: identity,
// byte length, truncation flag and digest all come from the same handle, and the
// bound+1 decision is the one the open time stat applied.
func TestBoundFileIdentityMatchesRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "usage.json")
	payload := []byte("identity-payload")
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write artifact: %v", err)
	}

	boundFile, err := OpenBoundFile(path, int64(len(payload)))
	if err != nil {
		t.Fatalf("OpenBoundFile: %v", err)
	}
	// Identity before the read: it must not consume the handle the read needs.
	size, modTime, mode := boundFile.Identity()
	if size != int64(len(payload)) {
		t.Fatalf("identity size %d must equal the on-disk %d bytes", size, len(payload))
	}
	if modTime.IsZero() {
		t.Fatal("identity must carry the handle's modification time")
	}
	if !mode.IsRegular() {
		t.Fatalf("identity mode must be regular, got %v", mode)
	}

	data, truncated, err := boundFile.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if truncated {
		t.Fatalf("a bound-sized read must not be truncated")
	}
	if size != int64(len(data)) {
		t.Fatalf("identity size %d must match the %d bytes read", size, len(data))
	}

	hash, err := boundFile.Hash(boundFileTestDomain)
	if err != nil {
		t.Fatalf("Hash after ReadAll: %v", err)
	}
	if want := Digest(boundFileTestDomain, payload); hash != want {
		t.Fatalf("digest must be re-readable from the same handle: got %s want %s", hash, want)
	}

	// The same bytes one byte over the bound must be refused: the truncation
	// decision is bound+1, not "the artifact fits whatever it needs".
	if _, err := OpenBoundFile(path, int64(len(payload))-1); err == nil {
		t.Fatalf("a %d byte artifact must be refused at a %d byte bound", len(payload), len(payload)-1)
	}
}
