//go:build windows

package adapters

func replayStorePublicationDurability() PublishDurability {
	return PublishDurabilityUnproven
}

func syncReplayStoreParentDirectory(string) error {
	return nil
}
