//go:build !darwin

package adapters

// P42-T027B-EVIDENCE-BIND Q2: the artifact-name guard shared by the windows and
// unix implementations. It is kept out of the darwin build because the darwin half
// is a fail-closed stub that never resolves a name, so the guard would have no
// caller there and only register as dead code.

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
