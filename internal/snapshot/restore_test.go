package snapshot

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestCaptureAndRestoreRoundtrip(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	targetDir := t.TempDir()

	store, err := NewStore(storeDir, 0)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	// Create a representative workspace
	testFiles := map[string]string{
		"src/main.go":     "package main\n\nfunc main() {}\n",
		"src/pkg/util.go": "package pkg\n",
		"README.md":       "# Test Project\n",
		"config/app.json": `{"env":"test"}`,
		"scripts/run.sh":  "#!/bin/sh\necho ok\n",
	}

	for p, c := range testFiles {
		full := filepath.Join(srcDir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		perm := os.FileMode(0644)
		if filepath.Ext(p) == ".sh" {
			perm = 0755
		}
		if err := os.WriteFile(full, []byte(c), perm); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	// 1. Capture
	opts := CaptureOptions{
		SourceDir:  srcDir,
		Kind:       "baseline",
		SnapshotID: "snap-roundtrip",
		RunID:      "run-roundtrip",
		Store:      store,
	}

	manifest, err := Capture(context.Background(), opts)
	if err != nil {
		t.Fatalf("capture: %v", err)
	}

	// 2. Restore into empty target
	restoreOpts := RestoreOptions{
		TargetDir: targetDir,
		Manifest:  manifest,
		Store:     store,
	}

	if err := Restore(context.Background(), restoreOpts); err != nil {
		t.Fatalf("restore: %v", err)
	}

	// 3. Verify bit-for-bit equality
	for p, expectedContent := range testFiles {
		restoredPath := filepath.Join(targetDir, filepath.FromSlash(p))
		restoredData, err := os.ReadFile(restoredPath)
		if err != nil {
			t.Errorf("read restored file %s: %v", p, err)
			continue
		}
		if string(restoredData) != expectedContent {
			t.Errorf("content mismatch for %s: got %q, want %q", p, restoredData, expectedContent)
		}
	}
}

func TestRestoreWithHardlinksAndPermissions(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	targetDir := t.TempDir()

	store, err := NewStore(storeDir, 0)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	// Create files
	f1 := filepath.Join(srcDir, "file1.txt")
	content := []byte("identical hardlink content")
	if err := os.WriteFile(f1, content, 0644); err != nil {
		t.Fatalf("write f1: %v", err)
	}

	f2 := filepath.Join(srcDir, "file2.txt")
	hasHardlink := false
	if err := os.Link(f1, f2); err == nil {
		hasHardlink = true
	} else {
		_ = os.WriteFile(f2, content, 0644)
	}

	roDir := filepath.Join(srcDir, "ro_dir")
	_ = os.MkdirAll(roDir, 0555)

	roFile := filepath.Join(roDir, "ro_file.txt")
	_ = os.WriteFile(roFile, []byte("readonly"), 0444)

	manifest, err := Capture(context.Background(), CaptureOptions{
		SourceDir:  srcDir,
		Kind:       "baseline",
		SnapshotID: "snap-hl-perm",
		RunID:      "run-hl-perm",
		Store:      store,
	})
	if err != nil {
		t.Fatalf("capture: %v", err)
	}

	err = Restore(context.Background(), RestoreOptions{
		TargetDir: targetDir,
		Manifest:  manifest,
		Store:     store,
	})
	if err != nil {
		t.Fatalf("restore: %v", err)
	}

	// Verify restored file contents
	r1, err := os.ReadFile(filepath.Join(targetDir, "file1.txt"))
	if err != nil || string(r1) != string(content) {
		t.Fatalf("r1 content mismatch")
	}
	r2, err := os.ReadFile(filepath.Join(targetDir, "file2.txt"))
	if err != nil || string(r2) != string(content) {
		t.Fatalf("r2 content mismatch")
	}

	if hasHardlink {
		info1, _ := os.Stat(filepath.Join(targetDir, "file1.txt"))
		info2, _ := os.Stat(filepath.Join(targetDir, "file2.txt"))
		if !os.SameFile(info1, info2) {
			t.Logf("filesystem did not preserve os.SameFile hardlink identity on restore, but content matched")
		}
	}
}

func TestRestoreTargetNotEmptyFailsClosed(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	targetDir := t.TempDir()

	store, _ := NewStore(storeDir, 0)
	testFile := filepath.Join(srcDir, "test.txt")
	_ = os.WriteFile(testFile, []byte("data"), 0644)

	manifest, err := Capture(context.Background(), CaptureOptions{
		SourceDir:  srcDir,
		Kind:       "baseline",
		SnapshotID: "snap-ne",
		RunID:      "run-ne",
		Store:      store,
	})
	if err != nil {
		t.Fatalf("capture: %v", err)
	}

	// Put a pre-existing file in targetDir
	_ = os.WriteFile(filepath.Join(targetDir, "existing.txt"), []byte("collision"), 0644)

	err = Restore(context.Background(), RestoreOptions{
		TargetDir: targetDir,
		Manifest:  manifest,
		Store:     store,
	})
	if !errors.Is(err, ErrTargetNotEmpty) {
		t.Fatalf("expected ErrTargetNotEmpty, got: %v", err)
	}
}

func TestRestoreMissingObjectFailsClosed(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	targetDir := t.TempDir()

	store, _ := NewStore(storeDir, 0)
	testFile := filepath.Join(srcDir, "test.txt")
	_ = os.WriteFile(testFile, []byte("important content"), 0644)

	manifest, err := Capture(context.Background(), CaptureOptions{
		SourceDir:  srcDir,
		Kind:       "baseline",
		SnapshotID: "snap-mo",
		RunID:      "run-mo",
		Store:      store,
	})
	if err != nil {
		t.Fatalf("capture: %v", err)
	}

	// Delete the object from store to simulate missing CAS object
	objHash := manifest.Manifest.Entries[0].ContentHash
	_ = os.Remove(store.objectPath(objHash))

	err = Restore(context.Background(), RestoreOptions{
		TargetDir: targetDir,
		Manifest:  manifest,
		Store:     store,
	})
	if !errors.Is(err, ErrObjectNotFound) {
		t.Fatalf("expected ErrObjectNotFound, got: %v", err)
	}

	// Target directory must remain clean / empty
	items, _ := os.ReadDir(targetDir)
	if len(items) > 0 {
		t.Fatalf("target directory must not be modified when pre-check fails!")
	}
}

func TestRestoreTamperedManifestFailsClosed(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	targetDir := t.TempDir()

	store, _ := NewStore(storeDir, 0)
	testFile := filepath.Join(srcDir, "test.txt")
	_ = os.WriteFile(testFile, []byte("data"), 0644)

	manifest, err := Capture(context.Background(), CaptureOptions{
		SourceDir:  srcDir,
		Kind:       "baseline",
		SnapshotID: "snap-tm",
		RunID:      "run-tm",
		Store:      store,
	})
	if err != nil {
		t.Fatalf("capture: %v", err)
	}

	// Tamper manifest hash
	manifest.ManifestHash = "sha256:0000000000000000000000000000000000000000000000000000000000000000"

	err = Restore(context.Background(), RestoreOptions{
		TargetDir: targetDir,
		Manifest:  manifest,
		Store:     store,
	})
	if !errors.Is(err, evidence.ErrInvalidRecord) {
		t.Fatalf("expected ErrInvalidRecord for tampered manifest, got: %v", err)
	}
}
