package snapshot

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestCreateBackupVerifiesNewStoreAndPreservesSource(t *testing.T) {
	ctx := context.Background()
	sourceDir := t.TempDir()
	store, err := NewStore(sourceDir, 0)
	if err != nil {
		t.Fatal(err)
	}
	first, err := store.PutObject(ctx, []byte("object-one"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.PutObject(ctx, []byte("object-two"))
	if err != nil {
		t.Fatal(err)
	}
	before, err := store.GetObject(ctx, first)
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "backup")
	result, err := CreateBackup(ctx, validBackupOptions(store, destination, []string{second, first}))
	if err != nil {
		t.Fatal(err)
	}
	if !evidence.ValidHash(result.BackupManifestHash) || !evidence.ValidHash(result.Manifest.DestinationIdentityHash) || len(result.Manifest.RestoreEvidence) != 1 {
		t.Fatalf("invalid completed backup result: %+v", result)
	}
	restored, err := NewStore(destination, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, hash := range []string{first, second} {
		if !restored.HasObject(hash) {
			t.Fatalf("restored store is missing %s", hash)
		}
	}
	after, err := store.GetObject(ctx, first)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("source store changed during backup")
	}
	if _, err := os.Stat(filepath.Join(destination, "backup-manifest.json")); err != nil {
		t.Fatalf("backup manifest missing: %v", err)
	}
}

func TestCreateBackupRejectsActiveWriterUnknownSchemaAndMissingClosure(t *testing.T) {
	ctx := context.Background()
	store, err := NewStore(t.TempDir(), 0)
	if err != nil {
		t.Fatal(err)
	}
	objectHash, err := store.PutObject(ctx, []byte("object"))
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		mutate func(*BackupOptions)
		target error
	}{
		{name: "active-writers", mutate: func(options *BackupOptions) { options.WritersStopped = false }, target: ErrBackupActiveWriters},
		{name: "unknown-schema", mutate: func(options *BackupOptions) { options.SourceSchemaVersion = "2.0.0" }, target: ErrBackupUnknownSchema},
		{name: "missing-object", mutate: func(options *BackupOptions) { options.ObjectHashes = []string{evidence.Digest("", []byte("missing"))} }, target: ErrBackupIncompleteClosure},
		{name: "missing-events", mutate: func(options *BackupOptions) { options.EventSegments = nil }, target: ErrBackupIncompleteClosure},
		{name: "missing-roots", mutate: func(options *BackupOptions) { options.ReferenceRoots = nil }, target: ErrBackupIncompleteClosure},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			destination := filepath.Join(t.TempDir(), "backup")
			options := validBackupOptions(store, destination, []string{objectHash})
			test.mutate(&options)
			result, err := CreateBackup(ctx, options)
			if !errors.Is(err, test.target) {
				t.Fatalf("expected %v, got result=%+v err=%v", test.target, result, err)
			}
			if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
				t.Fatalf("denied backup must not publish destination: %v", statErr)
			}
		})
	}
}

func TestCreateBackupInterruptedHasNoCompletedManifest(t *testing.T) {
	ctx := context.Background()
	store, err := NewStore(t.TempDir(), 0)
	if err != nil {
		t.Fatal(err)
	}
	objectHash, err := store.PutObject(ctx, []byte("object"))
	if err != nil {
		t.Fatal(err)
	}
	destination := filepath.Join(t.TempDir(), "backup")
	options := validBackupOptions(store, destination, []string{objectHash})
	options.FailAfterObjects = 1
	result, err := CreateBackup(ctx, options)
	if !errors.Is(err, ErrBackupInterrupted) {
		t.Fatalf("expected interruption, got result=%+v err=%v", result, err)
	}
	if result.BackupManifestHash != "" {
		t.Fatal("interrupted backup must not return a completed manifest hash")
	}
	if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
		t.Fatalf("interrupted backup must not publish destination: %v", statErr)
	}
}

func validBackupOptions(store *Store, destination string, objects []string) BackupOptions {
	eventData := []byte("event-segment")
	rootData := []byte("reference-root")
	return BackupOptions{
		DestinationDir:         destination,
		SourceStore:            store,
		SourceSchemaVersion:    SchemaVersion,
		WritersStopped:         true,
		WriterStopEvidence:     []string{evidence.Digest("", []byte("stopped"))},
		ObjectHashes:           objects,
		EventSegments:          []BackupBlob{{Hash: evidence.Digest("", eventData), Data: eventData}},
		ReferenceRoots:         []BackupBlob{{Hash: evidence.Digest("", rootData), Data: rootData}},
		SecretScanEvidence:     []string{evidence.Digest("", []byte("secret-scan"))},
		VerificationReportHash: evidence.Digest("", []byte("verification")),
	}
}
