//go:build !windows

package adapters

import "os"

func replayStorePublicationDurability() PublishDurability {
	return PublishDurabilityProven
}

func syncReplayStoreParentDirectory(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return err
	}
	defer directory.Close()
	return directory.Sync()
}
