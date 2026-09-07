package snapshot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type RestoreOptions struct {
	TargetDir string
	Manifest  SnapshotManifest
	Store     *Store
}

func Restore(ctx context.Context, opts RestoreOptions) error {
	if opts.Store == nil {
		return errors.New("snapshot store is required")
	}
	if err := VerifySnapshotManifest(opts.Manifest); err != nil {
		return fmt.Errorf("verify manifest before restore: %w", err)
	}

	// Pre-flight closure verification: all objects must exist in store before writing
	for _, entry := range opts.Manifest.Manifest.Entries {
		if entry.Type == EntryRegularFile {
			if !opts.Store.HasObject(entry.ContentHash) {
				return fmt.Errorf("%w: object %s missing for %q", ErrObjectNotFound, entry.ContentHash, entry.Path)
			}
		}
	}

	// Target directory check. Materialization happens in a sibling staging
	// directory so any later failure leaves the requested target untouched.
	targetExists := false
	if info, err := os.Stat(opts.TargetDir); err == nil {
		targetExists = true
		if !info.IsDir() {
			return fmt.Errorf("%w: target path is not a directory", ErrInvalidPath)
		}
		items, err := os.ReadDir(opts.TargetDir)
		if err != nil {
			return fmt.Errorf("read target dir: %w", err)
		}
		if len(items) > 0 {
			return fmt.Errorf("%w: target directory %q is not empty", ErrTargetNotEmpty, opts.TargetDir)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat target dir: %w", err)
	}
	parentDir := filepath.Dir(filepath.Clean(opts.TargetDir))
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("create target parent: %w", err)
	}
	stageDir, err := os.MkdirTemp(parentDir, ".prfrail-restore-*")
	if err != nil {
		return fmt.Errorf("create restore staging directory: %w", err)
	}
	defer os.RemoveAll(stageDir)

	hardlinkPrimaries := make(map[string]string) // groupID -> full path of first restored file
	var readOnlyDirs []struct {
		path string
		perm os.FileMode
	}

	// Restore entries in strict manifest order
	for _, entry := range opts.Manifest.Manifest.Entries {
		destPath := filepath.Join(stageDir, filepath.FromSlash(entry.Path))

		switch entry.Type {
		case EntryDirectory:
			perm := os.FileMode(0755)
			if entry.ReadOnly {
				perm = 0555
			}
			if err := os.MkdirAll(destPath, perm); err != nil {
				return fmt.Errorf("mkdir %q: %w", entry.Path, err)
			}
			if entry.ReadOnly {
				readOnlyDirs = append(readOnlyDirs, struct {
					path string
					perm os.FileMode
				}{destPath, perm})
			}

		case EntryRegularFile:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("mkdir parent of %q: %w", entry.Path, err)
			}

			// Check if we can hardlink to already restored primary file
			linked := false
			if entry.HardlinkGroup != nil {
				if primaryPath, exists := hardlinkPrimaries[*entry.HardlinkGroup]; exists {
					if err := os.Link(primaryPath, destPath); err != nil {
						return fmt.Errorf("%w: create hardlink %q: %v", ErrUnsupportedLink, entry.Path, err)
					}
					linked = true
				}
			}

			if !linked {
				data, err := opts.Store.GetObject(ctx, entry.ContentHash)
				if err != nil {
					return fmt.Errorf("get object for %q: %w", entry.Path, err)
				}

				perm := os.FileMode(0644)
				if entry.Executable && entry.ReadOnly {
					perm = 0555
				} else if entry.Executable {
					perm = 0755
				} else if entry.ReadOnly {
					perm = 0444
				}

				if err := os.WriteFile(destPath, data, perm); err != nil {
					return fmt.Errorf("write regular file %q: %w", entry.Path, err)
				}
				if err := os.Chmod(destPath, perm); err != nil {
					return fmt.Errorf("restore permissions for %q: %w", entry.Path, err)
				}

				if entry.HardlinkGroup != nil {
					hardlinkPrimaries[*entry.HardlinkGroup] = destPath
				}
			}

		case EntrySymbolicLink:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("mkdir parent of symlink %q: %w", entry.Path, err)
			}
			// In Go, os.Symlink(oldname, newname) creates newname as a symbolic link to oldname
			target := filepath.FromSlash(entry.Target)
			if err := os.Symlink(target, destPath); err != nil {
				return fmt.Errorf("%w: create symlink %q -> %q: %v", ErrUnsupportedLink, entry.Path, entry.Target, err)
			}

		default:
			return fmt.Errorf("%w: unknown entry type %q for %q", ErrInvalidManifest, entry.Type, entry.Path)
		}
	}
	for i := len(readOnlyDirs) - 1; i >= 0; i-- {
		if err := os.Chmod(readOnlyDirs[i].path, readOnlyDirs[i].perm); err != nil {
			return fmt.Errorf("restore directory permissions: %w", err)
		}
	}
	if targetExists {
		if err := os.Remove(opts.TargetDir); err != nil {
			return fmt.Errorf("remove empty target before publication: %w", err)
		}
	}
	if err := os.Rename(stageDir, opts.TargetDir); err != nil {
		if targetExists {
			_ = os.Mkdir(opts.TargetDir, 0755)
		}
		return fmt.Errorf("publish restored target: %w", err)
	}

	return nil
}
