package evidence

// P42-T027B-EVIDENCE-BIND Q2: the artifact-name guard shared by the windows and
// unix implementations.
//
// P83 (X-1): the build tag that used to exclude darwin was dropped. The guard is
// exported for the adapters alias, and an alias must compile on every platform;
// on darwin it simply has no internal caller because that half never resolves a
// name.
//
// P85 R=7 (user ruling of 2026-11-23): boundRootName and errBoundRootArtifact were
// moved to bounddir_clean_root.go (build tag !darwin, i.e. the windows half and
// every non-darwin unix half). The darwin half never references them, so keeping
// them in this shared file made them look unused in the darwin IDE view. Only the
// two root-name helpers moved; the artifact-name guard below stays shared.

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
