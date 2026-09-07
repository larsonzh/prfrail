//go:build !windows

package snapshot

import "os"

func isReparsePointOrJunction(_ os.FileInfo, _ string) bool {
	return false
}
