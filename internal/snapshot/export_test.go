package snapshot

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestExportAcceptedPackageRoundTrip(t *testing.T) {
	sourceDir := t.TempDir()
	storeDir := t.TempDir()
	runDir := t.TempDir()
	deliveryRoot := t.TempDir()
	destinationDir := filepath.Join(deliveryRoot, "delivery")

	writeFile(t, filepath.Join(sourceDir, "bin", "app.exe"), []byte("binary-content"), 0755)
	writeFile(t, filepath.Join(sourceDir, "docs", "readme.txt"), []byte("release notes"), 0644)

	store, err := NewStore(storeDir, 0)
	if err != nil {
		t.Fatal(err)
	}
	manifest := captureCandidateManifest(t, sourceDir, store)
	before, err := os.ReadFile(filepath.Join(sourceDir, "docs", "readme.txt"))
	if err != nil {
		t.Fatal(err)
	}

	exported, err := ExportAcceptedPackage(context.Background(), ExportOptions{
		DestinationDir:       destinationDir,
		SourceDir:            sourceDir,
		RunDir:               runDir,
		StoreDir:             storeDir,
		AcceptedSnapshotHash: manifest.ManifestHash,
		Manifest:             manifest,
		Store:                store,
		Exclusions: []evidence.ExportExclusion{{
			Category:     "cache",
			Pattern:      "tmp/cache",
			MatchedCount: 0,
		}},
	})
	if err != nil {
		t.Fatalf("export accepted package: %v", err)
	}
	if len(exported.Entries) == 0 || !evidence.ValidHash(exported.PackageManifestHash) || !evidence.ValidHash(exported.DestinationIdentityHash) {
		t.Fatalf("invalid export summary: %+v", exported)
	}
	record, err := evidence.NewExportRecord(evidence.DeliveryExport{
		ExportID:                "export-one",
		CreatedAt:               "2026-09-08T12:00:00.000Z",
		RecordedBy:              evidence.Actor{Type: "system", ID: "proofrail"},
		RunID:                   "run-one",
		TaskID:                  "task-one",
		Attempt:                 1,
		SourceSnapshotHash:      manifest.ManifestHash,
		ReviewReceiptHash:       evidence.Digest("", []byte("review")),
		PromotionReceiptHash:    evidence.Digest("", []byte("promotion")),
		VerificationReportHash:  evidence.Digest("", []byte("verification")),
		DestinationRef:          "delivery",
		DestinationIdentityHash: exported.DestinationIdentityHash,
		DestinationPrecondition: "absent",
		Entries:                 exported.Entries,
		Exclusions:              exported.Exclusions,
		Outcome:                 "completed",
		PackageManifestHash:     &exported.PackageManifestHash,
		Evidence:                exported.Evidence,
		ErrorEvidence:           []string{},
	})
	if err != nil {
		t.Fatalf("wrap completed export in record: %v", err)
	}
	if err := evidence.ValidateExportRecord(record); err != nil {
		t.Fatalf("validate completed export record: %v", err)
	}

	after, err := os.ReadFile(filepath.Join(sourceDir, "docs", "readme.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("source file changed after export")
	}

	for _, entry := range exported.Entries {
		fullPath := filepath.Join(destinationDir, filepath.FromSlash(entry.Path))
		data, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("read exported file %s: %v", entry.Path, err)
		}
		if evidence.Digest("", data) != entry.ContentHash {
			t.Fatalf("exported file hash mismatch: %s", entry.Path)
		}
	}
}

func TestExportAcceptedPackageRejectsHashMismatchAndExistingDestination(t *testing.T) {
	sourceDir := t.TempDir()
	storeDir := t.TempDir()
	runDir := t.TempDir()
	deliveryRoot := t.TempDir()
	destinationDir := filepath.Join(deliveryRoot, "delivery")

	writeFile(t, filepath.Join(sourceDir, "bin", "app.exe"), []byte("binary-content"), 0755)
	store, err := NewStore(storeDir, 0)
	if err != nil {
		t.Fatal(err)
	}
	manifest := captureCandidateManifest(t, sourceDir, store)

	_, err = ExportAcceptedPackage(context.Background(), ExportOptions{
		DestinationDir:       destinationDir,
		SourceDir:            sourceDir,
		RunDir:               runDir,
		StoreDir:             storeDir,
		AcceptedSnapshotHash: evidence.Digest("", []byte("other")),
		Manifest:             manifest,
		Store:                store,
	})
	if !errors.Is(err, ErrExportDenied) {
		t.Fatalf("expected ErrExportDenied, got %v", err)
	}

	if err := os.MkdirAll(destinationDir, 0755); err != nil {
		t.Fatal(err)
	}
	_, err = ExportAcceptedPackage(context.Background(), ExportOptions{
		DestinationDir:       destinationDir,
		SourceDir:            sourceDir,
		RunDir:               runDir,
		StoreDir:             storeDir,
		AcceptedSnapshotHash: manifest.ManifestHash,
		Manifest:             manifest,
		Store:                store,
	})
	if !errors.Is(err, ErrExportDestinationExists) {
		t.Fatalf("expected ErrExportDestinationExists, got %v", err)
	}
}

func TestExportAcceptedPackageRejectsOverlapSecretMissingAndTamper(t *testing.T) {
	t.Run("overlap", func(t *testing.T) {
		sourceDir := t.TempDir()
		storeDir := t.TempDir()
		runDir := t.TempDir()
		writeFile(t, filepath.Join(sourceDir, "bin", "app.exe"), []byte("binary-content"), 0755)
		store, _ := NewStore(storeDir, 0)
		manifest := captureCandidateManifest(t, sourceDir, store)
		_, err := ExportAcceptedPackage(context.Background(), ExportOptions{
			DestinationDir:       filepath.Join(sourceDir, "delivery"),
			SourceDir:            sourceDir,
			RunDir:               runDir,
			StoreDir:             storeDir,
			AcceptedSnapshotHash: manifest.ManifestHash,
			Manifest:             manifest,
			Store:                store,
		})
		if !errors.Is(err, ErrExportOverlap) {
			t.Fatalf("expected ErrExportOverlap, got %v", err)
		}
	})

	t.Run("secret", func(t *testing.T) {
		sourceDir := t.TempDir()
		storeDir := t.TempDir()
		runDir := t.TempDir()
		deliveryRoot := t.TempDir()
		writeFile(t, filepath.Join(sourceDir, "secrets", "api-token.txt"), []byte("token"), 0644)
		writeFile(t, filepath.Join(sourceDir, "bin", "app.exe"), []byte("binary-content"), 0755)
		store, _ := NewStore(storeDir, 0)
		manifest := captureCandidateManifest(t, sourceDir, store)
		_, err := ExportAcceptedPackage(context.Background(), ExportOptions{
			DestinationDir:       filepath.Join(deliveryRoot, "delivery"),
			SourceDir:            sourceDir,
			RunDir:               runDir,
			StoreDir:             storeDir,
			AcceptedSnapshotHash: manifest.ManifestHash,
			Manifest:             manifest,
			Store:                store,
		})
		if !errors.Is(err, ErrExportSecret) {
			t.Fatalf("expected ErrExportSecret, got %v", err)
		}
	})

	t.Run("missing-object", func(t *testing.T) {
		sourceDir := t.TempDir()
		storeDir := t.TempDir()
		runDir := t.TempDir()
		deliveryRoot := t.TempDir()
		writeFile(t, filepath.Join(sourceDir, "bin", "app.exe"), []byte("binary-content"), 0755)
		store, _ := NewStore(storeDir, 0)
		manifest := captureCandidateManifest(t, sourceDir, store)

		for _, entry := range manifest.Manifest.Entries {
			if entry.Type != EntryRegularFile {
				continue
			}
			objectPath, pathErr := store.objectPath(entry.ContentHash)
			if pathErr != nil {
				t.Fatal(pathErr)
			}
			if err := os.Remove(objectPath); err != nil {
				t.Fatal(err)
			}
			break
		}

		_, err := ExportAcceptedPackage(context.Background(), ExportOptions{
			DestinationDir:       filepath.Join(deliveryRoot, "delivery"),
			SourceDir:            sourceDir,
			RunDir:               runDir,
			StoreDir:             storeDir,
			AcceptedSnapshotHash: manifest.ManifestHash,
			Manifest:             manifest,
			Store:                store,
		})
		if !errors.Is(err, ErrObjectNotFound) {
			t.Fatalf("expected ErrObjectNotFound, got %v", err)
		}
	})

	t.Run("tampered-manifest", func(t *testing.T) {
		sourceDir := t.TempDir()
		storeDir := t.TempDir()
		runDir := t.TempDir()
		deliveryRoot := t.TempDir()
		writeFile(t, filepath.Join(sourceDir, "bin", "app.exe"), []byte("binary-content"), 0755)
		store, _ := NewStore(storeDir, 0)
		manifest := captureCandidateManifest(t, sourceDir, store)
		manifest.ManifestHash = evidence.Digest("", []byte("tampered"))
		_, err := ExportAcceptedPackage(context.Background(), ExportOptions{
			DestinationDir:       filepath.Join(deliveryRoot, "delivery"),
			SourceDir:            sourceDir,
			RunDir:               runDir,
			StoreDir:             storeDir,
			AcceptedSnapshotHash: manifest.ManifestHash,
			Manifest:             manifest,
			Store:                store,
		})
		if !errors.Is(err, evidence.ErrInvalidRecord) {
			t.Fatalf("expected ErrInvalidRecord, got %v", err)
		}
	})
}

func TestExportAcceptedPackageInterruptedWriteHasNoCompletion(t *testing.T) {
	sourceDir := t.TempDir()
	storeDir := t.TempDir()
	runDir := t.TempDir()
	deliveryRoot := t.TempDir()
	destinationDir := filepath.Join(deliveryRoot, "delivery")

	writeFile(t, filepath.Join(sourceDir, "bin", "app.exe"), []byte("binary-content"), 0755)
	writeFile(t, filepath.Join(sourceDir, "docs", "readme.txt"), []byte("release notes"), 0644)
	store, err := NewStore(storeDir, 0)
	if err != nil {
		t.Fatal(err)
	}
	manifest := captureCandidateManifest(t, sourceDir, store)

	result, err := ExportAcceptedPackage(context.Background(), ExportOptions{
		DestinationDir:       destinationDir,
		SourceDir:            sourceDir,
		RunDir:               runDir,
		StoreDir:             storeDir,
		AcceptedSnapshotHash: manifest.ManifestHash,
		Manifest:             manifest,
		Store:                store,
		FailAfterEntries:     1,
	})
	if !errors.Is(err, ErrExportInterrupted) {
		t.Fatalf("expected ErrExportInterrupted, got %v", err)
	}
	if result.PackageManifestHash != "" || result.DestinationIdentityHash != "" || len(result.Entries) != 0 {
		t.Fatalf("interrupted export must not emit completion summary: %+v", result)
	}
	if _, recErr := evidence.NewExportRecord(evidence.DeliveryExport{
		ExportID:                "export-interrupted",
		CreatedAt:               "2026-09-08T12:00:00.000Z",
		RecordedBy:              evidence.Actor{Type: "system", ID: "proofrail"},
		RunID:                   "run-one",
		TaskID:                  "task-one",
		Attempt:                 1,
		SourceSnapshotHash:      manifest.ManifestHash,
		ReviewReceiptHash:       evidence.Digest("", []byte("review")),
		PromotionReceiptHash:    evidence.Digest("", []byte("promotion")),
		VerificationReportHash:  evidence.Digest("", []byte("verification")),
		DestinationRef:          "delivery",
		DestinationIdentityHash: evidence.Digest("", []byte("destination")),
		DestinationPrecondition: "absent",
		Entries:                 result.Entries,
		Exclusions:              result.Exclusions,
		Outcome:                 "completed",
		PackageManifestHash:     nil,
		Evidence:                result.Evidence,
		ErrorEvidence:           []string{},
	}); recErr == nil {
		t.Fatal("interrupted export must not be representable as completed export record")
	}
	if _, statErr := os.Stat(destinationDir); !os.IsNotExist(statErr) {
		t.Fatalf("destination should stay absent on interrupted export: %v", statErr)
	}
}

func captureCandidateManifest(t *testing.T, sourceDir string, store *Store) SnapshotManifest {
	t.Helper()
	parentHash := evidence.Digest("", []byte("parent"))
	manifest, err := Capture(context.Background(), CaptureOptions{
		SourceDir:          sourceDir,
		Kind:               "candidate",
		SnapshotID:         "candidate-one",
		RunID:              "run-one",
		ParentSnapshotHash: &parentHash,
		Task:               &TaskBinding{TaskID: "task-one", Attempt: 1},
		Store:              store,
	})
	if err != nil {
		t.Fatal(err)
	}
	return manifest
}

func writeFile(t *testing.T, path string, data []byte, perm os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		t.Fatal(err)
	}
}
