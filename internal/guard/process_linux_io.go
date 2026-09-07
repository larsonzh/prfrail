//go:build linux

package guard

import "os"

func readFileString(path string) (string, error) {
	data, err := os.ReadFile(path)
	return string(data), err
}
