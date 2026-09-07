//go:build !windows

package taskdef

import "os"

const defaultFileMode os.FileMode = 0644
