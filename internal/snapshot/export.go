package snapshot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	exportDestinationIdentityDomain = "proofrail:export-destination-identity:1\n"
	exportPackageManifestDomain     = "proofrail:delivery-package-manifest:1\n"
)

var (
	ErrExportDenied            = errors.New("export denied")
	ErrExportDestinationExists = errors.New("export destination must be absent")
	ErrExportOverlap           = errors.New("export destination overlaps protected paths")
	ErrExportSecret            = errors.New("export contains secret artifact")
	ErrExportInterrupted       = errors.New("export interrupted")
)

type ExportOptions struct {
	DestinationDir       string
	SourceDir            string
	RunDir               string
	StoreDir             string
	AcceptedSnapshotHash string
	Manifest             SnapshotManifest
	Store                *Store
	Exclusions           []evidence.ExportExclusion
	FailAfterEntries     int
}

type ExportPackage struct {
	Entries                 []evidence.ExportEntry     `json:"entries"`
	Exclusions              []evidence.ExportExclusion `json:"exclusions"`
	DestinationIdentityHash string                     `json:"destinationIdentityHash"`
	PackageManifestHash     string                     `json:"packageManifestHash"`
	Evidence                []string                   `json:"evidence"`
}

func ExportAcceptedPackage(ctx context.Context, options ExportOptions) (ExportPackage, error) {
	if options.Store == nil {
		return ExportPackage{}, errors.New("snapshot store is required")
	}
	if err := VerifySnapshotManifest(options.Manifest); err != nil {
		return ExportPackage{}, err
	}
	if !evidence.ValidHash(options.AcceptedSnapshotHash) || options.AcceptedSnapshotHash != options.Manifest.ManifestHash {
		return ExportPackage{}, fmt.Errorf("%w: accepted snapshot hash mismatch", ErrExportDenied)
	}

	destinationPath, err := resolvePathForCompare(options.DestinationDir)
	if err != nil {
		return ExportPackage{}, err
	}
	sourcePath, err := resolvePathForCompare(options.SourceDir)
	if err != nil {
		return ExportPackage{}, err
	}
	runPath, err := resolvePathForCompare(options.RunDir)
	if err != nil {
		return ExportPackage{}, err
	}
	storePath, err := resolvePathForCompare(options.StoreDir)
	if err != nil {
		return ExportPackage{}, err
	}
	if err := ensureDestinationAbsent(destinationPath); err != nil {
		return ExportPackage{}, err
	}
	if err := ensureNoOverlap(destinationPath, sourcePath, runPath, storePath); err != nil {
		return ExportPackage{}, err
	}

	exportEntries, err := collectExportEntries(options.Manifest, options.Store)
	if err != nil {
		return ExportPackage{}, err
	}
	if len(exportEntries) == 0 {
		return ExportPackage{}, fmt.Errorf("%w: no exportable regular files", ErrExportDenied)
	}

	exclusions := normalizeExclusions(options.Exclusions)
	stageParent := filepath.Dir(destinationPath)
	if err := os.MkdirAll(stageParent, 0755); err != nil {
		return ExportPackage{}, fmt.Errorf("create export destination parent: %w", err)
	}
	stageDir, err := os.MkdirTemp(stageParent, ".prfrail-export-*")
	if err != nil {
		return ExportPackage{}, fmt.Errorf("create export staging directory: %w", err)
	}
	defer os.RemoveAll(stageDir)

	if err := materializeExport(ctx, stageDir, options.Manifest, options.Store, options.FailAfterEntries); err != nil {
		return ExportPackage{}, err
	}
	if err := verifyExportStage(stageDir, exportEntries); err != nil {
		return ExportPackage{}, err
	}

	destinationIdentityHash, packageManifestHash, err := computeExportHashes(exportEntries, exclusions)
	if err != nil {
		return ExportPackage{}, err
	}
	if err := os.Rename(stageDir, destinationPath); err != nil {
		return ExportPackage{}, fmt.Errorf("publish exported package: %w", err)
	}

	return ExportPackage{
		Entries:                 exportEntries,
		Exclusions:              exclusions,
		DestinationIdentityHash: destinationIdentityHash,
		PackageManifestHash:     packageManifestHash,
		Evidence: []string{
			options.Manifest.ManifestHash,
			destinationIdentityHash,
			packageManifestHash,
		},
	}, nil
}

func resolvePathForCompare(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("%w: empty export path", ErrInvalidPath)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	clean := filepath.Clean(abs)
	if resolved, err := filepath.EvalSymlinks(clean); err == nil {
		clean = resolved
	}
	return clean, nil
}

func ensureDestinationAbsent(destinationPath string) error {
	if _, err := os.Lstat(destinationPath); err == nil {
		return fmt.Errorf("%w: %q", ErrExportDestinationExists, destinationPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("check export destination: %w", err)
	}
	return nil
}

func ensureNoOverlap(destinationPath string, protectedPaths ...string) error {
	for _, protectedPath := range protectedPaths {
		if protectedPath == "" {
			continue
		}
		if pathsOverlap(destinationPath, protectedPath) {
			return fmt.Errorf("%w: destination %q overlaps %q", ErrExportOverlap, destinationPath, protectedPath)
		}
	}
	return nil
}

func pathsOverlap(left, right string) bool {
	return pathContains(left, right) || pathContains(right, left)
}

func pathContains(root, target string) bool {
	root = normalizePathForCompare(root)
	target = normalizePathForCompare(target)
	if root == target {
		return true
	}
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return false
	}
	return true
}

func normalizePathForCompare(path string) string {
	clean := filepath.Clean(path)
	if runtime.GOOS == "windows" {
		clean = strings.ToLower(clean)
	}
	return clean
}

func collectExportEntries(manifest SnapshotManifest, store *Store) ([]evidence.ExportEntry, error) {
	entries := make([]evidence.ExportEntry, 0)
	for _, entry := range manifest.Manifest.Entries {
		if isSecretPath(entry.Path) {
			return nil, fmt.Errorf("%w: %q", ErrExportSecret, entry.Path)
		}
		if entry.Type != EntryRegularFile {
			continue
		}
		if !store.HasObject(entry.ContentHash) {
			return nil, fmt.Errorf("%w: object %s missing for %q", ErrObjectNotFound, entry.ContentHash, entry.Path)
		}
		entries = append(entries, evidence.ExportEntry{
			Path:        entry.Path,
			Type:        string(entry.Type),
			ContentHash: entry.ContentHash,
			SizeBytes:   entry.SizeBytes,
		})
	}
	sort.Slice(entries, func(left, right int) bool {
		return entries[left].Path < entries[right].Path
	})
	return entries, nil
}

func materializeExport(ctx context.Context, stageDir string, manifest SnapshotManifest, store *Store, failAfterEntries int) error {
	writtenEntries := 0
	for _, entry := range manifest.Manifest.Entries {
		targetPath := filepath.Join(stageDir, filepath.FromSlash(entry.Path))
		switch entry.Type {
		case EntryDirectory:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("create export directory %q: %w", entry.Path, err)
			}
		case EntryRegularFile:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("create export parent %q: %w", entry.Path, err)
			}
			data, err := store.GetObject(ctx, entry.ContentHash)
			if err != nil {
				return fmt.Errorf("load export object %q: %w", entry.Path, err)
			}
			perm := os.FileMode(0644)
			if entry.Executable && entry.ReadOnly {
				perm = 0555
			} else if entry.Executable {
				perm = 0755
			} else if entry.ReadOnly {
				perm = 0444
			}
			if err := os.WriteFile(targetPath, data, perm); err != nil {
				return fmt.Errorf("write export file %q: %w", entry.Path, err)
			}
			if err := os.Chmod(targetPath, perm); err != nil {
				return fmt.Errorf("set export permissions %q: %w", entry.Path, err)
			}
			writtenEntries++
			if failAfterEntries > 0 && writtenEntries >= failAfterEntries {
				return fmt.Errorf("%w: injected after %d file writes", ErrExportInterrupted, writtenEntries)
			}
		case EntrySymbolicLink:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("create export symlink parent %q: %w", entry.Path, err)
			}
			if err := os.Symlink(filepath.FromSlash(entry.Target), targetPath); err != nil {
				return fmt.Errorf("%w: create export symlink %q: %v", ErrUnsupportedLink, entry.Path, err)
			}
		default:
			return fmt.Errorf("%w: unsupported entry type %q", ErrInvalidManifest, entry.Type)
		}
	}
	return nil
}

func verifyExportStage(stageDir string, entries []evidence.ExportEntry) error {
	for _, entry := range entries {
		fullPath := filepath.Join(stageDir, filepath.FromSlash(entry.Path))
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("read exported file %q: %w", entry.Path, err)
		}
		if int64(len(data)) != entry.SizeBytes {
			return fmt.Errorf("%w: exported file size mismatch for %q", ErrObjectCorrupt, entry.Path)
		}
		if evidence.Digest("", data) != entry.ContentHash {
			return fmt.Errorf("%w: exported file hash mismatch for %q", ErrObjectCorrupt, entry.Path)
		}
	}
	return nil
}

func computeExportHashes(entries []evidence.ExportEntry, exclusions []evidence.ExportExclusion) (string, string, error) {
	identityPayload := struct {
		Entries []evidence.ExportEntry `json:"entries"`
	}{Entries: entries}
	identityCanonical, err := evidence.EncodeCanonical(identityPayload)
	if err != nil {
		return "", "", err
	}
	destinationIdentityHash := evidence.Digest(exportDestinationIdentityDomain, identityCanonical)

	manifestPayload := struct {
		DestinationIdentityHash string                     `json:"destinationIdentityHash"`
		Entries                 []evidence.ExportEntry     `json:"entries"`
		Exclusions              []evidence.ExportExclusion `json:"exclusions"`
	}{
		DestinationIdentityHash: destinationIdentityHash,
		Entries:                 entries,
		Exclusions:              exclusions,
	}
	manifestCanonical, err := evidence.EncodeCanonical(manifestPayload)
	if err != nil {
		return "", "", err
	}
	packageManifestHash := evidence.Digest(exportPackageManifestDomain, manifestCanonical)
	return destinationIdentityHash, packageManifestHash, nil
}

func normalizeExclusions(exclusions []evidence.ExportExclusion) []evidence.ExportExclusion {
	if len(exclusions) == 0 {
		return []evidence.ExportExclusion{}
	}
	out := make([]evidence.ExportExclusion, len(exclusions))
	copy(out, exclusions)
	sort.Slice(out, func(left, right int) bool {
		if out[left].Category != out[right].Category {
			return out[left].Category < out[right].Category
		}
		return out[left].Pattern < out[right].Pattern
	})
	return out
}

func isSecretPath(path string) bool {
	lower := strings.ToLower(path)
	segments := strings.Split(lower, "/")
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		if strings.HasPrefix(segment, ".env") {
			return true
		}
		tokens := []string{
			"secret", "token", "password", "passwd", "credential", "private-key",
			"id_rsa", "id_ed25519", "api-key", "apikey",
		}
		for _, token := range tokens {
			if segment == token || strings.Contains(segment, token) {
				return true
			}
		}
	}
	return false
}
