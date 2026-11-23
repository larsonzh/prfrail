//go:build !darwin

package evidence

// P85 R=7 (user ruling of 2026-11-23): boundRootName and errBoundRootArtifact
// moved here from bounddir_clean.go. They are referenced by the windows half
// (bounddir_windows.go) and by the non-darwin unix half (bounddir_unix.go), never
// by bounddir_darwin.go, so the shared file made them look unused in the darwin
// IDE view. This relocation fixes that diagnostic only; it is not a precedent for
// keeping code that is otherwise dead.
//
// The tag is !darwin rather than "windows || linux" because bounddir_unix.go is
// itself tagged !windows && !darwin: every GOOS that compiles that file needs both
// symbols, and a windows||linux tag here would break the remaining unix targets.

import (
	"fmt"
	"path/filepath"
	"strings"
)

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
