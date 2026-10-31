//go:build !windows

package evidence

// P42-T027B-EVIDENCE-BIND: on Unix an open descriptor keeps addressing the file
// object after the name is renamed, replaced or unlinked, so a plain open already
// provides the same-handle binding the primitive needs.

import "os"

func openBoundFileHandle(path string) (*os.File, error) {
	return os.Open(path)
}
