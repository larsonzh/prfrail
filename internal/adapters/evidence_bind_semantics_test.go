package adapters

// P42-T027B-EVIDENCE-BIND ③ test-first (RED): adapter-level binding semantics
// that hold on every platform (design section 2 "Q2" and "failure and
// partial-read semantics", section 4 I2/I3).
//
// These tests are platform-neutral on purpose: the store-level seam
// (readBoundEvidence) and the log truncation fallback must behave identically on
// Windows and Unix, while the platform primitives they rest on are pinned
// separately by evidence_bind_dir_windows_test.go / evidence_bind_dir_unix_test.go.

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

const evidenceBindStoreTestBound int64 = 1 << 20

// TestStoreReadBoundEvidenceSurvivesAncestorReplacement pins I3 at the seam the
// adapters will actually call: after RunEvidenceDir bound the evidence directory,
// an ancestor renamed away and re-created by an attacker cannot redirect the read.
func TestStoreReadBoundEvidenceSurvivesAncestorReplacement(t *testing.T) {
	root := t.TempDir()
	store := replayStoreMustNew(t, root)
	evidenceDir, err := store.RunEvidenceDir("request-one")
	if err != nil {
		t.Fatalf("RunEvidenceDir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(evidenceDir, "f.txt"), []byte("bound"), 0o644); err != nil {
		t.Fatalf("write bound evidence file: %v", err)
	}

	// Replace the ancestor directory of the bound evidence directory with an
	// attacker-controlled tree at the same path.
	ancestor := filepath.Dir(evidenceDir)
	if err := os.Rename(ancestor, ancestor+".bak"); err != nil {
		t.Fatalf("rename the bound ancestor away (the dir handle must be share-delete): %v", err)
	}
	injected := filepath.Join(ancestor, filepath.Base(evidenceDir))
	if err := os.MkdirAll(injected, 0o755); err != nil {
		t.Fatalf("create the injected evidence directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(injected, "f.txt"), []byte("attacker"), 0o644); err != nil {
		t.Fatalf("write the injected file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(injected, "injected.txt"), []byte("attacker"), 0o644); err != nil {
		t.Fatalf("write the injected-only file: %v", err)
	}

	data, truncated, err := store.readBoundEvidence("request-one", "f.txt", evidenceBindStoreTestBound)
	if err != nil {
		t.Fatalf("readBoundEvidence through the bound directory: %v", err)
	}
	if string(data) != "bound" {
		t.Fatalf("bound evidence read must follow the bound directory object, got %q", data)
	}
	if truncated {
		t.Fatal("a small bound evidence read must not be truncated")
	}

	// A name that only the injected tree has must stay missing for the bound
	// directory: an ErrNotExist is the contract, injected bytes are a breach.
	if _, _, err := store.readBoundEvidence("request-one", "injected.txt", evidenceBindStoreTestBound); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("a file absent from the bound evidence directory must be fs.ErrNotExist, got %v", err)
	}
}

// TestBoundedArtifactKeepsLogTruncationSemantics pins I2 for the log/event class:
// the record and PID class refuses an over-bound artifact at open time, but the
// log class must keep degrading to a truncated read (design section 2 and risk 3).
// Pinning it here is what stops the new open-time refusal from drifting the
// existing log truncation behaviour.
func TestBoundedArtifactKeepsLogTruncationSemantics(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stdout.log")
	payload := bytes.Repeat([]byte("a"), 64)
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write log artifact: %v", err)
	}
	const bound int64 = 16

	data, truncated, err := readBoundedAgentRunnerArtifact(path, bound)
	if err != nil {
		t.Fatalf("an over-bound log artifact must degrade to a truncated read, got %v", err)
	}
	if !truncated {
		t.Fatal("an over-bound log artifact must report truncated, never a complete read")
	}
	if int64(len(data)) > bound+1 {
		t.Fatalf("a truncated read must stop within bound+1 bytes (%d), got %d", bound+1, len(data))
	}
	if !bytes.Equal(data, payload[:len(data)]) {
		t.Fatal("a truncated read must return a prefix of the artifact, never substituted bytes")
	}
}
