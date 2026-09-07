package snapshot

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestStorePutGetDedup(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewStore(tmpDir, 0)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	data := []byte("hello proofrail snapshot store")
	hash, err := store.PutObject(context.Background(), data)
	if err != nil {
		t.Fatalf("put object: %v", err)
	}
	expectedHash := evidence.Digest("", data)
	if hash != expectedHash {
		t.Fatalf("hash mismatch: got %s, want %s", hash, expectedHash)
	}

	// Verify deduplication
	hash2, err := store.PutObject(context.Background(), data)
	if err != nil {
		t.Fatalf("second put: %v", err)
	}
	if hash2 != hash {
		t.Fatalf("dedup hash mismatch: %s vs %s", hash2, hash)
	}

	// Get object
	readBack, err := store.GetObject(context.Background(), hash)
	if err != nil {
		t.Fatalf("get object: %v", err)
	}
	if string(readBack) != string(data) {
		t.Fatalf("data mismatch: got %q, want %q", readBack, data)
	}

	if !store.HasObject(hash) {
		t.Fatalf("expected HasObject(%s) = true", hash)
	}

	totalBytes, err := store.TotalBytes()
	if err != nil || totalBytes != int64(len(data)) {
		t.Fatalf("totalBytes = %d, want %d", totalBytes, len(data))
	}
}

func TestStoreNotFoundAndCorrupt(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewStore(tmpDir, 0)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	// Non-existent
	missingHash := "sha256:0000000000000000000000000000000000000000000000000000000000000000"
	if _, err := store.GetObject(context.Background(), missingHash); !errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("expected ErrObjectNotFound, got: %v", err)
	}
	if store.HasObject(missingHash) {
		t.Fatalf("HasObject for missing should be false")
	}

	// Corrupted object
	data := []byte("original data")
	hash, err := store.PutObject(context.Background(), data)
	if err != nil {
		t.Fatalf("put: %v", err)
	}

	objFile, _ := store.objectPath(hash)
	if err := os.WriteFile(objFile, []byte("tampered data"), 0644); err != nil {
		t.Fatalf("tamper file: %v", err)
	}

	if _, err := store.GetObject(context.Background(), hash); !errors.Is(err, ErrObjectCorrupt) {
		t.Fatalf("expected ErrObjectCorrupt, got: %v", err)
	}
	if store.HasObject(hash) {
		t.Fatalf("HasObject for corrupt object should be false")
	}
	if _, err := store.PutObject(context.Background(), data); !errors.Is(err, ErrObjectCorrupt) {
		t.Fatalf("corrupt immutable object must not be replaced: %v", err)
	}
	if _, err := store.GetObject(context.Background(), "../../outside"); !errors.Is(err, ErrInvalidManifest) {
		t.Fatalf("invalid hash must be rejected before path construction: %v", err)
	}
}

func TestStoreQuotaExceeded(t *testing.T) {
	tmpDir := t.TempDir()
	// Set quota to 50 bytes
	store, err := NewStore(tmpDir, 50)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	// Put 30 bytes
	data1 := make([]byte, 30)
	if _, err := store.PutObject(context.Background(), data1); err != nil {
		t.Fatalf("put 30 bytes: %v", err)
	}

	// Put 25 bytes -> total 55 > 50 -> must fail with ErrQuotaExceeded
	data2 := make([]byte, 25)
	data2[0] = 1
	if _, err := store.PutObject(context.Background(), data2); !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got: %v", err)
	}
}

func TestStoreGCDryRun(t *testing.T) {
	tmpDir := t.TempDir()
	store, err := NewStore(tmpDir, 0)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	hash1, _ := store.PutObject(context.Background(), []byte("reachable"))
	hash2, _ := store.PutObject(context.Background(), []byte("audit-locked"))
	hash3, _ := store.PutObject(context.Background(), []byte("orphaned-old"))
	hash4, _ := store.PutObject(context.Background(), []byte("orphaned-recent"))

	// Backdate hash3 object file to 2 hours ago
	oldTime := time.Now().Add(-2 * time.Hour)
	hash3Path, _ := store.objectPath(hash3)
	if err := os.Chtimes(hash3Path, oldTime, oldTime); err != nil {
		t.Fatalf("chtimes: %v", err)
	}

	reachable := map[string]bool{hash1: true}
	locked := map[string]bool{hash2: true}
	ttl := 1 * time.Hour

	plan, err := store.GCDryRun(context.Background(), reachable, locked, ttl)
	if err != nil {
		t.Fatalf("gc dry run: %v", err)
	}

	// hash3 is the only candidate
	if len(plan.CandidateHashes) != 1 || plan.CandidateHashes[0] != hash3 {
		t.Fatalf("expected candidates [hash3], got: %v", plan.CandidateHashes)
	}

	// Critical: verify no files were deleted
	if !store.HasObject(hash1) || !store.HasObject(hash2) || !store.HasObject(hash3) || !store.HasObject(hash4) {
		t.Fatalf("GCDryRun must never delete any files!")
	}
}
