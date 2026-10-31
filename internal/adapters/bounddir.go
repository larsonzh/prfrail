package adapters

// P42-T027B-EVIDENCE-BIND: the Q2 ancestor-directory binding primitive.
//
// docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md
// section 2 "Q2 ancestor directory replacement" requires a directory handle that
// outlives the path: once the per-request evidence directory is bound, an
// ancestor renamed away and re-created by an attacker must not be able to
// redirect a read. Both platform implementations resolve a relative artifact
// name against the bound directory object, never against a re-resolved path:
//
//   - Windows: CreateFile with FILE_FLAG_BACKUP_SEMANTICS opens the directory
//     object; a relative open derives the object's current real name with
//     GetFinalPathNameByHandle and opens under it.
//   - Unix: openat(2) with O_NOFOLLOW and O_CLOEXEC resolves relative to the
//     directory descriptor, so a replaced ancestor cannot redirect it and a leaf
//     swapped for a symlink is refused.
//
// BoundDir, OpenBoundDir and OpenBoundDirFile are the shared surface; the
// platform files own the mechanics and the post-open verification.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// cleanBoundRelativePath validates an artifact name that is opened relative to a
// bound directory. Absolute names and names that escape the bound directory are
// refused: a bound read may never resolve outside the object it bound.
func cleanBoundRelativePath(rel string) (string, error) {
	if strings.TrimSpace(rel) == "" {
		return "", fmt.Errorf("adapters: empty bound artifact name")
	}
	if filepath.IsAbs(rel) || strings.HasPrefix(rel, "/") || strings.HasPrefix(rel, `\`) {
		return "", fmt.Errorf("adapters: bound artifact name %q must be relative", rel)
	}
	clean := filepath.Clean(filepath.FromSlash(rel))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("adapters: bound artifact name %q escapes the bound directory", rel)
	}
	return clean, nil
}
