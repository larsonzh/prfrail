package snapshot

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestCaptureBaselineUncommittedAndSourceUnchanged(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()

	store, err := NewStore(storeDir, 0)
	if err != nil {
		t.Fatalf("new store: %v", err)
	}

	// Create source files including uncommitted/untracked files
	files := map[string][]byte{
		"src/main.go":        []byte("package main\n\nfunc main() {}\n"),
		"src/util/helper.go": []byte("package util\n"),
		"package.json":       []byte(`{"name": "test-pkg"}`),
		"uncommitted.txt":    []byte("this is untracked or uncommitted work\n"),
		".git/config":        []byte("[core]\n"),
		"tmp/scratch.log":    []byte("should be excluded\n"),
	}

	for p, content := range files {
		fullPath := filepath.Join(srcDir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(fullPath, content, 0644); err != nil {
			t.Fatalf("write file: %v", err)
		}
	}

	// Record stat before capture
	mainPath := filepath.Join(srcDir, "src", "main.go")
	statBefore, err := os.Stat(mainPath)
	if err != nil {
		t.Fatalf("stat before: %v", err)
	}

	opts := CaptureOptions{
		SourceDir:     srcDir,
		Kind:          "baseline",
		SnapshotID:    "snapshot-baseline-0",
		RunID:         "run-test-001",
		Store:         store,
		MaxBytesQuota: 10 * 1024 * 1024,
	}

	manifest, err := Capture(context.Background(), opts)
	if err != nil {
		t.Fatalf("capture failed: %v", err)
	}

	// Verify source tree is completely unchanged
	statAfter, err := os.Stat(mainPath)
	if err != nil {
		t.Fatalf("stat after: %v", err)
	}
	if !statBefore.ModTime().Equal(statAfter.ModTime()) || statBefore.Size() != statAfter.Size() {
		t.Fatalf("source tree was modified during capture!")
	}

	// Verify manifest entries
	// Excluded paths (.git, tmp) must NOT be present
	entryPaths := make(map[string]Entry)
	for _, e := range manifest.Manifest.Entries {
		entryPaths[e.Path] = e
	}

	if _, ok := entryPaths[".git"]; ok {
		t.Errorf(".git must be excluded")
	}
	if _, ok := entryPaths[".git/config"]; ok {
		t.Errorf(".git/config must be excluded")
	}
	if _, ok := entryPaths["tmp"]; ok {
		t.Errorf("tmp must be excluded")
	}
	if _, ok := entryPaths["tmp/scratch.log"]; ok {
		t.Errorf("tmp/scratch.log must be excluded")
	}

	// uncommitted.txt must be captured
	if uncommitted, ok := entryPaths["uncommitted.txt"]; !ok {
		t.Errorf("uncommitted.txt must be captured")
	} else {
		if uncommitted.ContentHash != evidence.Digest("", files["uncommitted.txt"]) {
			t.Errorf("uncommitted.txt hash mismatch")
		}
	}

	// package.json must have RolePackageManifest
	if pkg, ok := entryPaths["package.json"]; !ok {
		t.Errorf("package.json missing")
	} else if pkg.Role != RolePackageManifest {
		t.Errorf("package.json role = %v, want RolePackageManifest", pkg.Role)
	}

	// Verify directory closure
	if _, ok := entryPaths["src"]; !ok {
		t.Errorf("directory 'src' missing")
	}
	if _, ok := entryPaths["src/util"]; !ok {
		t.Errorf("directory 'src/util' missing")
	}

	// Verify entries are sorted strictly ascending
	for i := 1; i < len(manifest.Manifest.Entries); i++ {
		if manifest.Manifest.Entries[i-1].Path >= manifest.Manifest.Entries[i].Path {
			t.Errorf("entries unsorted at %d: %s >= %s", i, manifest.Manifest.Entries[i-1].Path, manifest.Manifest.Entries[i].Path)
		}
	}
}

func TestCaptureReservedNameRejection(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	store, _ := NewStore(storeDir, 0)

	// Create file with Windows reserved name
	conFile := filepath.Join(srcDir, "con.txt")
	// On Windows, creating "con.txt" might directly be blocked by OS or allowed via \\?\
	// If OS blocks creation, test ValidatePath directly; if allowed, test Capture
	err := os.WriteFile(conFile, []byte("reserved"), 0644)
	if err == nil {
		opts := CaptureOptions{
			SourceDir:  srcDir,
			Kind:       "baseline",
			SnapshotID: "snap-res",
			RunID:      "run-res",
			Store:      store,
		}
		_, err := Capture(context.Background(), opts)
		if !errors.Is(err, ErrReservedPath) {
			t.Fatalf("expected ErrReservedPath, got: %v", err)
		}
	} else {
		// OS itself blocked reserved name, verify ValidatePath
		if err := ValidatePath("con.txt"); !errors.Is(err, ErrReservedPath) {
			t.Fatalf("expected ErrReservedPath, got: %v", err)
		}
	}
}

func TestCaptureQuotaExceeded(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	store, _ := NewStore(storeDir, 0)

	bigFile := filepath.Join(srcDir, "big.bin")
	if err := os.WriteFile(bigFile, make([]byte, 1000), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	opts := CaptureOptions{
		SourceDir:     srcDir,
		Kind:          "baseline",
		SnapshotID:    "snap-quota",
		RunID:         "run-quota",
		Store:         store,
		MaxBytesQuota: 500, // 500 < 1000
	}

	_, err := Capture(context.Background(), opts)
	if !errors.Is(err, ErrQuotaExceeded) {
		t.Fatalf("expected ErrQuotaExceeded, got: %v", err)
	}
}

func TestCaptureHardlinks(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	store, _ := NewStore(storeDir, 0)

	f1 := filepath.Join(srcDir, "original.txt")
	content := []byte("shared content across hardlinks")
	if err := os.WriteFile(f1, content, 0644); err != nil {
		t.Fatalf("write original: %v", err)
	}
	f2 := filepath.Join(srcDir, "linked.txt")
	if err := os.Link(f1, f2); err != nil {
		t.Skipf("hardlink not supported on this filesystem: %v", err)
	}

	opts := CaptureOptions{
		SourceDir:  srcDir,
		Kind:       "baseline",
		SnapshotID: "snap-hl",
		RunID:      "run-hl",
		Store:      store,
	}

	manifest, err := Capture(context.Background(), opts)
	if err != nil {
		t.Fatalf("capture: %v", err)
	}

	var group1, group2 *string
	for _, e := range manifest.Manifest.Entries {
		if e.Path == "original.txt" {
			group1 = e.HardlinkGroup
		}
		if e.Path == "linked.txt" {
			group2 = e.HardlinkGroup
		}
	}

	if group1 == nil || group2 == nil {
		t.Fatalf("hardlinkGroup must be set for both linked files")
	}
	if *group1 != *group2 {
		t.Fatalf("hardlinkGroup mismatch: %s vs %s", *group1, *group2)
	}
}

func TestCaptureSymlinkValid(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	store, _ := NewStore(storeDir, 0)

	targetFile := filepath.Join(srcDir, "target.txt")
	if err := os.WriteFile(targetFile, []byte("target data"), 0644); err != nil {
		t.Fatalf("write target: %v", err)
	}

	linkFile := filepath.Join(srcDir, "link.txt")
	if err := os.Symlink("target.txt", linkFile); err != nil {
		t.Skipf("symlink creation not permitted on this host: %v", err)
	}

	manifest, err := Capture(context.Background(), CaptureOptions{
		SourceDir:  srcDir,
		Kind:       "baseline",
		SnapshotID: "snap-sym",
		RunID:      "run-sym",
		Store:      store,
	})
	if err != nil {
		t.Fatalf("capture symlink: %v", err)
	}

	foundLink := false
	for _, e := range manifest.Manifest.Entries {
		if e.Path == "link.txt" {
			foundLink = true
			if e.Type != EntrySymbolicLink {
				t.Fatalf("expected EntrySymbolicLink, got %v", e.Type)
			}
			if e.Target != "target.txt" {
				t.Fatalf("expected target 'target.txt', got %q", e.Target)
			}
		}
	}
	if !foundLink {
		t.Fatalf("symlink entry not found in manifest")
	}
}

func TestCaptureConcurrentModification(t *testing.T) {
	srcDir := t.TempDir()
	storeDir := t.TempDir()
	store, _ := NewStore(storeDir, 0)

	// Create many files to give time for concurrent modification
	for i := 0; i < 50; i++ {
		p := filepath.Join(srcDir, filepath.FromSlash("file"+string(rune('a'+i%26))+".dat"))
		_ = os.WriteFile(p, []byte("steady content"), 0644)
	}

	changeFile := filepath.Join(srcDir, "changing.dat")
	_ = os.WriteFile(changeFile, []byte("initial"), 0644)

	// In a goroutine, modify changing.dat repeatedly
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_ = os.WriteFile(changeFile, []byte("modified content at "+time.Now().String()), 0644)
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	// Attempt capture multiple times if needed until a concurrent change is detected
	detected := false
	for attempt := 0; attempt < 20; attempt++ {
		_, err := Capture(context.Background(), CaptureOptions{
			SourceDir:  srcDir,
			Kind:       "baseline",
			SnapshotID: "snap-conc",
			RunID:      "run-conc",
			Store:      store,
		})
		if errors.Is(err, ErrConcurrentChange) {
			detected = true
			break
		}
	}
	cancel()
	if !detected {
		t.Logf("concurrent change timing did not trigger in 20 attempts (filesystem race dependent)")
	}
}
