//go:build linux

package console

// P85 F02/F06/F20: the frozen-policy publication happens INSIDE the verified
// directory descriptor.
//
// The previous shape wrote a same-directory temporary file and hard-linked the
// final name BY PATH, after a binding that compared the handle with a fresh
// resolution of the same name. Two path operations after that comparison meant the
// name could be pointed at another directory in between, and the comparison itself
// only proved that the name agreed with itself.
//
// Every step below is relative to the descriptor the binding produced: the
// temporary file is created in it (openat, O_CREAT|O_EXCL), the final name is a
// hard link to it (linkat), the temporary name is removed from it (unlinkat) and
// the descriptor is fsynced. No path is resolved after the bind, so the object that
// was proven is the object that is written and made durable.
//
// F09 is unchanged: the hard link is the publication step, there is no rename
// fallback anywhere, and a hard-link failure that is not a no-replace collision is
// a refusal (fail-closed).

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	// agentRunnerPolicyTemporaryPrefix keeps a leftover from either shape
	// recognizable (the path-based half uses the same prefix).
	agentRunnerPolicyTemporaryPrefix = ".agent-runner-policy-"
	// agentRunnerPolicyTemporaryAttempts bounds the unique-name retry loop: an
	// EEXIST means another writer used the same random name, which is a collision
	// to retry, not a reason to publish over anything.
	agentRunnerPolicyTemporaryAttempts = 8
)

// publishAgentRunnerPolicyFileIn publishes one payload inside bound, without
// replacing an existing name. fileID is the identity the binding produced and is
// asserted again on the handle before the directory is fsynced.
func publishAgentRunnerPolicyFileIn(bound *evidence.BoundDir, fileID evidence.FileID, directory, base string, payload []byte, label string) error {
	path := filepath.Join(directory, base)
	for attempt := 0; attempt < agentRunnerPolicyTemporaryAttempts; attempt++ {
		temporaryName, err := agentRunnerPolicyTemporaryName()
		if err != nil {
			return fmt.Errorf("%w: generate a %s temporary name: %v", ErrInvalidConfig, label, err)
		}
		temporary, err := bound.CreateFileExclusive(temporaryName, 0o600)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			}
			return fmt.Errorf("%w: create the %s temporary file: %v", ErrInvalidConfig, label, err)
		}
		writeErr := temporary.Chmod(0o600)
		if writeErr == nil {
			_, writeErr = temporary.Write(payload)
		}
		if writeErr == nil {
			writeErr = temporary.Sync()
		}
		if closeErr := temporary.Close(); writeErr == nil {
			writeErr = closeErr
		}
		if writeErr != nil {
			_ = bound.Remove(temporaryName)
			return fmt.Errorf("%w: write the %s temporary file: %v", ErrInvalidConfig, label, writeErr)
		}
		if err := bound.LinkNoReplace(temporaryName, base); err != nil {
			_ = bound.Remove(temporaryName)
			if errors.Is(err, os.ErrExist) {
				return fmt.Errorf("%w: the %s %s already exists; a frozen policy is never replaced", ErrInvalidConfig, label, path)
			}
			return fmt.Errorf("%w: publish the %s %s without replacing it: %v", ErrInvalidConfig, label, path, err)
		}
		// A temporary name that cannot be removed is a stray entry, not a failed
		// publication, so the error is not turned into a refusal.
		_ = bound.Remove(temporaryName)
		if err := evidence.VerifyBoundDir(bound, fileID); err != nil {
			return fmt.Errorf("%w: the %s directory %s is no longer the approved object: %v", ErrInvalidConfig, label, directory, err)
		}
		// F10 (P85 design rev 3, section 2): the fsync acts on the SAME descriptor
		// the link was published through, so "the directory that was confirmed" and
		// "the directory whose entry became durable" are one object.
		if err := bound.SyncDirectory(); err != nil {
			return fmt.Errorf("%w: sync the %s directory %s: %v", ErrInvalidConfig, label, directory, err)
		}
		return nil
	}
	return fmt.Errorf("%w: no unique %s temporary name could be created", ErrInvalidConfig, label)
}

// agentRunnerPolicyTemporaryName returns one unguessable entry name. A name
// another writer could predict would let it pre-create the entry and turn the
// exclusive create into a refusal.
func agentRunnerPolicyTemporaryName() (string, error) {
	var buffer [8]byte
	if _, err := rand.Read(buffer[:]); err != nil {
		return "", err
	}
	return agentRunnerPolicyTemporaryPrefix + hex.EncodeToString(buffer[:]), nil
}
