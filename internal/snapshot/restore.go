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

	// Target directory check
	if info, err := os.Stat(opts.TargetDir); err == nil {
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
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(opts.TargetDir, 0755); err != nil {
			return fmt.Errorf("create target dir: %w", err)
		}
	} else {
		return fmt.Errorf("stat target dir: %w", err)
	}

	hardlinkPrimaries := make(map[string]string) // groupID -> full path of first restored file

	// Restore entries in strict manifest order
	for _, entry := range opts.Manifest.Manifest.Entries {
		destPath := filepath.Join(opts.TargetDir, filepath.FromSlash(entry.Path))

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
				_ = os.Chmod(destPath, perm)
			}

		case EntryRegularFile:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("mkdir parent of %q: %w", entry.Path, err)
			}

			// Check if we can hardlink to already restored primary file
			linked := false
			if entry.HardlinkGroup != nil {
				if primaryPath, exists := hardlinkPrimaries[*entry.HardlinkGroup]; exists {
					if err := os.Link(primaryPath, destPath); err == nil {
						linked = true
					}
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
				if entry.ReadOnly {
					_ = os.Chmod(destPath, perm)
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

	return nil
}
