package snapshot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	backupSourceIdentityDomain      = "proofrail:backup-source-identity:1\n"
	backupDestinationIdentityDomain = "proofrail:backup-destination-identity:1\n"
	backupManifestDomain            = "proofrail:backup-manifest:1\n"
	backupRestoreEvidenceDomain     = "proofrail:backup-restore-evidence:1\n"
)

var (
	ErrBackupDenied            = errors.New("backup denied")
	ErrBackupActiveWriters     = errors.New("backup requires stopped writers")
	ErrBackupUnknownSchema     = errors.New("backup source schema is unsupported")
	ErrBackupIncompleteClosure = errors.New("backup reference closure is incomplete")
	ErrBackupDestinationExists = errors.New("backup destination must be absent")
	ErrBackupInterrupted       = errors.New("backup interrupted")
)

type BackupBlob struct {
	Hash string
	Data []byte
}

type BackupOptions struct {
	DestinationDir         string
	SourceStore            *Store
	SourceSchemaVersion    string
	WritersStopped         bool
	WriterStopEvidence     []string
	ObjectHashes           []string
	EventSegments          []BackupBlob
	ReferenceRoots         []BackupBlob
	SecretScanEvidence     []string
	VerificationReportHash string
	FailAfterObjects       int
}

type BackupManifest struct {
	SchemaVersion           string   `json:"schemaVersion"`
	SourceStoreHash         string   `json:"sourceStoreHash"`
	DestinationIdentityHash string   `json:"destinationIdentityHash"`
	ObjectHashes            []string `json:"objectHashes"`
	EventSegmentHashes      []string `json:"eventSegmentHashes"`
	ReferenceRootHashes     []string `json:"referenceRootHashes"`
	SecretDisposition       string   `json:"secretDisposition"`
	SecretScanEvidence      []string `json:"secretScanEvidence"`
	VerificationReportHash  string   `json:"verificationReportHash"`
	RestoreEvidence         []string `json:"restoreEvidence"`
}

type BackupResult struct {
	Manifest           BackupManifest `json:"manifest"`
	BackupManifestHash string         `json:"backupManifestHash"`
}

func CreateBackup(ctx context.Context, options BackupOptions) (BackupResult, error) {
	manifest, objects, events, roots, err := prepareBackup(options)
	if err != nil {
		return BackupResult{}, err
	}
	destination, err := resolvePathForCompare(options.DestinationDir)
	if err != nil {
		return BackupResult{}, err
	}
	if err := ensureBackupDestinationAbsent(destination); err != nil {
		return BackupResult{}, err
	}
	sourceRoot, err := resolvePathForCompare(options.SourceStore.root)
	if err != nil {
		return BackupResult{}, err
	}
	if pathsOverlap(destination, sourceRoot) {
		return BackupResult{}, fmt.Errorf("%w: destination overlaps source store", ErrBackupDenied)
	}

	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return BackupResult{}, fmt.Errorf("create backup parent: %w", err)
	}
	stageDir, err := os.MkdirTemp(filepath.Dir(destination), ".prfrail-backup-*")
	if err != nil {
		return BackupResult{}, fmt.Errorf("create backup staging directory: %w", err)
	}
	defer os.RemoveAll(stageDir)

	destinationStore, err := NewStore(stageDir, 0)
	if err != nil {
		return BackupResult{}, err
	}
	for index, hash := range objects {
		data, err := options.SourceStore.GetObject(ctx, hash)
		if err != nil {
			return BackupResult{}, fmt.Errorf("copy backup object %s: %w", hash, err)
		}
		writtenHash, err := destinationStore.PutObject(ctx, data)
		if err != nil {
			return BackupResult{}, fmt.Errorf("write backup object %s: %w", hash, err)
		}
		if writtenHash != hash {
			return BackupResult{}, fmt.Errorf("%w: copied object hash mismatch", ErrBackupIncompleteClosure)
		}
		if options.FailAfterObjects > 0 && index+1 >= options.FailAfterObjects {
			return BackupResult{}, fmt.Errorf("%w: injected after %d object(s)", ErrBackupInterrupted, index+1)
		}
	}
	if err := writeBackupBlobs(stageDir, "event-segments", events); err != nil {
		return BackupResult{}, err
	}
	if err := writeBackupBlobs(stageDir, "reference-roots", roots); err != nil {
		return BackupResult{}, err
	}

	destinationIdentity, restoreEvidence, err := verifyBackupStage(ctx, stageDir, manifest, events, roots)
	if err != nil {
		return BackupResult{}, err
	}
	manifest.DestinationIdentityHash = destinationIdentity
	manifest.RestoreEvidence = []string{restoreEvidence}
	canonical, err := evidence.EncodeCanonical(manifest)
	if err != nil {
		return BackupResult{}, err
	}
	manifestHash := evidence.Digest(backupManifestDomain, canonical)
	if err := os.WriteFile(filepath.Join(stageDir, "backup-manifest.json"), append(canonical, '\n'), 0644); err != nil {
		return BackupResult{}, fmt.Errorf("write backup manifest: %w", err)
	}
	if err := os.Rename(stageDir, destination); err != nil {
		return BackupResult{}, fmt.Errorf("publish backup: %w", err)
	}
	return BackupResult{Manifest: manifest, BackupManifestHash: manifestHash}, nil
}

func prepareBackup(options BackupOptions) (BackupManifest, []string, []BackupBlob, []BackupBlob, error) {
	if options.SourceStore == nil {
		return BackupManifest{}, nil, nil, nil, errors.New("snapshot store is required")
	}
	if options.SourceSchemaVersion != SchemaVersion {
		return BackupManifest{}, nil, nil, nil, fmt.Errorf("%w: %q", ErrBackupUnknownSchema, options.SourceSchemaVersion)
	}
	if !options.WritersStopped || len(options.WriterStopEvidence) == 0 || !allBackupHashes(options.WriterStopEvidence) {
		return BackupManifest{}, nil, nil, nil, ErrBackupActiveWriters
	}
	if len(options.SecretScanEvidence) == 0 || !allBackupHashes(options.SecretScanEvidence) || !evidence.ValidHash(options.VerificationReportHash) {
		return BackupManifest{}, nil, nil, nil, fmt.Errorf("%w: invalid scan or verification evidence", ErrBackupDenied)
	}
	objects, err := normalizeBackupHashes(options.ObjectHashes)
	if err != nil || len(objects) == 0 {
		return BackupManifest{}, nil, nil, nil, fmt.Errorf("%w: object set", ErrBackupIncompleteClosure)
	}
	events, err := normalizeBackupBlobs(options.EventSegments)
	if err != nil || len(events) == 0 {
		return BackupManifest{}, nil, nil, nil, fmt.Errorf("%w: event segments", ErrBackupIncompleteClosure)
	}
	roots, err := normalizeBackupBlobs(options.ReferenceRoots)
	if err != nil || len(roots) == 0 {
		return BackupManifest{}, nil, nil, nil, fmt.Errorf("%w: reference roots", ErrBackupIncompleteClosure)
	}
	for _, hash := range objects {
		if !options.SourceStore.HasObject(hash) {
			return BackupManifest{}, nil, nil, nil, fmt.Errorf("%w: missing object %s", ErrBackupIncompleteClosure, hash)
		}
	}
	eventHashes := backupBlobHashes(events)
	rootHashes := backupBlobHashes(roots)
	sourcePayload := struct {
		SchemaVersion       string   `json:"schemaVersion"`
		ObjectHashes        []string `json:"objectHashes"`
		EventSegmentHashes  []string `json:"eventSegmentHashes"`
		ReferenceRootHashes []string `json:"referenceRootHashes"`
	}{options.SourceSchemaVersion, objects, eventHashes, rootHashes}
	canonical, err := evidence.EncodeCanonical(sourcePayload)
	if err != nil {
		return BackupManifest{}, nil, nil, nil, err
	}
	return BackupManifest{
		SchemaVersion:           options.SourceSchemaVersion,
		SourceStoreHash:         evidence.Digest(backupSourceIdentityDomain, canonical),
		DestinationIdentityHash: evidence.Digest("", []byte("pending")),
		ObjectHashes:            objects,
		EventSegmentHashes:      eventHashes,
		ReferenceRootHashes:     rootHashes,
		SecretDisposition:       "excluded",
		SecretScanEvidence:      cloneBackupStrings(options.SecretScanEvidence),
		VerificationReportHash:  options.VerificationReportHash,
		RestoreEvidence:         []string{},
	}, objects, events, roots, nil
}

func verifyBackupStage(ctx context.Context, stageDir string, manifest BackupManifest, events, roots []BackupBlob) (string, string, error) {
	store, err := NewStore(stageDir, 0)
	if err != nil {
		return "", "", err
	}
	for _, hash := range manifest.ObjectHashes {
		if _, err := store.GetObject(ctx, hash); err != nil {
			return "", "", fmt.Errorf("%w: restored object %s: %v", ErrBackupIncompleteClosure, hash, err)
		}
	}
	if err := verifyBackupBlobs(stageDir, "event-segments", events); err != nil {
		return "", "", err
	}
	if err := verifyBackupBlobs(stageDir, "reference-roots", roots); err != nil {
		return "", "", err
	}
	identityPayload := struct {
		SourceStoreHash     string   `json:"sourceStoreHash"`
		ObjectHashes        []string `json:"objectHashes"`
		EventSegmentHashes  []string `json:"eventSegmentHashes"`
		ReferenceRootHashes []string `json:"referenceRootHashes"`
	}{manifest.SourceStoreHash, manifest.ObjectHashes, manifest.EventSegmentHashes, manifest.ReferenceRootHashes}
	canonical, err := evidence.EncodeCanonical(identityPayload)
	if err != nil {
		return "", "", err
	}
	destinationIdentity := evidence.Digest(backupDestinationIdentityDomain, canonical)
	restoreEvidence := evidence.Digest(backupRestoreEvidenceDomain, canonical)
	return destinationIdentity, restoreEvidence, nil
}

func writeBackupBlobs(root, kind string, blobs []BackupBlob) error {
	dir := filepath.Join(root, kind)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	for _, blob := range blobs {
		if err := os.WriteFile(filepath.Join(dir, blob.Hash[7:]), blob.Data, 0444); err != nil {
			return fmt.Errorf("write backup %s %s: %w", kind, blob.Hash, err)
		}
	}
	return nil
}

func verifyBackupBlobs(root, kind string, blobs []BackupBlob) error {
	for _, blob := range blobs {
		data, err := os.ReadFile(filepath.Join(root, kind, blob.Hash[7:]))
		if err != nil || evidence.Digest("", data) != blob.Hash {
			return fmt.Errorf("%w: invalid %s %s", ErrBackupIncompleteClosure, kind, blob.Hash)
		}
	}
	return nil
}

func normalizeBackupHashes(values []string) ([]string, error) {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if !evidence.ValidHash(value) {
			return nil, ErrBackupIncompleteClosure
		}
		if _, exists := seen[value]; exists {
			return nil, ErrBackupIncompleteClosure
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result, nil
}

func normalizeBackupBlobs(values []BackupBlob) ([]BackupBlob, error) {
	result := make([]BackupBlob, len(values))
	seen := make(map[string]struct{}, len(values))
	for index, value := range values {
		if !evidence.ValidHash(value.Hash) || evidence.Digest("", value.Data) != value.Hash {
			return nil, ErrBackupIncompleteClosure
		}
		if _, exists := seen[value.Hash]; exists {
			return nil, ErrBackupIncompleteClosure
		}
		seen[value.Hash] = struct{}{}
		result[index] = BackupBlob{Hash: value.Hash, Data: append([]byte(nil), value.Data...)}
	}
	sort.Slice(result, func(left, right int) bool { return result[left].Hash < result[right].Hash })
	return result, nil
}

func backupBlobHashes(values []BackupBlob) []string {
	result := make([]string, len(values))
	for index, value := range values {
		result[index] = value.Hash
	}
	return result
}

func ensureBackupDestinationAbsent(path string) error {
	if _, err := os.Lstat(path); err == nil {
		return fmt.Errorf("%w: %q", ErrBackupDestinationExists, path)
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func allBackupHashes(values []string) bool {
	for _, value := range values {
		if !evidence.ValidHash(value) {
			return false
		}
	}
	return true
}

func cloneBackupStrings(values []string) []string {
	result := make([]string, len(values))
	copy(result, values)
	return result
}
