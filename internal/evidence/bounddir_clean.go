package evidence

// P42-T027B-EVIDENCE-BIND Q2: the artifact-name guard shared by the windows and
// unix implementations.
//
// P83 (X-1): the build tag that used to exclude darwin was dropped. The guard is
// exported for the adapters alias, and an alias must compile on every platform;
// on darwin it simply has no internal caller because that half never resolves a
// name.

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
		return "", fmt.Errorf("evidence: empty bound artifact name")
	}
	if filepath.IsAbs(rel) || strings.HasPrefix(rel, "/") || strings.HasPrefix(rel, `\`) {
		return "", fmt.Errorf("evidence: bound artifact name %q must be relative", rel)
	}
	clean := filepath.Clean(filepath.FromSlash(rel))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("evidence: bound artifact name %q escapes the bound directory", rel)
	}
	return clean, nil
}

// CleanBoundRelativePath is the exported form of the artifact-name guard, kept for
// callers outside the package (the adapters alias of the original seam).
func CleanBoundRelativePath(rel string) (string, error) {
	return cleanBoundRelativePath(rel)
}

// boundRootName reports whether rel names the bound directory itself instead of an
// entry inside it.
//
// P84 T4(d): the empty string and "." are the same name, and both halves have to
// read it the same way. BoundDirNames enumerates the bound object for either
// spelling, while an artifact open refuses either spelling - an artifact name that
// resolves to the bound directory is not an artifact. cleanBoundRelativePath keeps
// its own refusal of "." as the artifact-name guard; this classifier exists so the
// enumeration half does not inherit that refusal for the root spelling.
func boundRootName(rel string) bool {
	return strings.TrimSpace(rel) == "" || filepath.Clean(filepath.FromSlash(rel)) == "."
}

// errBoundRootArtifact is the one refusal both halves return when an artifact open
// names the bound directory itself. One constructor keeps the windows and unix
// messages from drifting apart for the two root spellings.
func errBoundRootArtifact(rel string) error {
	return fmt.Errorf("evidence: bound artifact name %q names the bound directory itself", rel)
}
