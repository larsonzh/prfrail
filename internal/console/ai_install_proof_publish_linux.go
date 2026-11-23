//go:build linux

package console

// P85 F20 final review: the carrier is published INSIDE the directory descriptor.
//
// The previous shape bound the directory, re-confirmed it, and then created the
// temporary file and published the final name BY PATH. Two path operations after
// the confirmation means the name could be pointed at a different directory in
// between, so "the directory that was proven" and "the directory the carrier landed
// in" were two claims again - the same window the binding was introduced to close.
//
// Every step below is relative to the descriptor that was confirmed: the temporary
// file is created in it (openat2, O_CREAT|O_EXCL), the final name is a hard link to
// it (linkat), and the temporary name is removed from it (unlinkat). No path is
// resolved after the bind, so the object that was proven is the object that is
// written. The no-replace discipline is unchanged: the hard link is the publication
// step and there is no rename fallback anywhere.

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
	// aiInstallSourceProofTemporaryPrefix is the same prefix the path-based half
	// uses, so a leftover from either shape is recognizable.
	aiInstallSourceProofTemporaryPrefix = ".ai-install-source-proof-"
	// aiInstallSourceProofTemporaryAttempts bounds the unique-name retry loop: an
	// EEXIST means another writer used the same random name, which is a collision
	// to retry, not a reason to publish over anything.
	aiInstallSourceProofTemporaryAttempts = 8
)

// writeAIInstallSourceProofFileIn publishes one payload without replacing an
// existing name, entirely inside one bound directory descriptor (P85 F10: the
// descriptor is the one the caller already read and confirmed).
func writeAIInstallSourceProofFileIn(bound *evidence.BoundDir, fileID evidence.FileID, path string, payload []byte) error {
	directory := filepath.Dir(path)
	published := filepath.Base(path)
	for attempt := 0; attempt < aiInstallSourceProofTemporaryAttempts; attempt++ {
		temporaryName, err := aiInstallSourceProofTemporaryName()
		if err != nil {
			return fmt.Errorf("%w: generate an install source proof temporary name: %v", ErrAIInstallSourceUnproven, err)
		}
		temporary, err := bound.CreateFileExclusive(temporaryName, 0o600)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			}
			return fmt.Errorf("%w: create the install source proof temporary file: %v", ErrAIInstallSourceUnproven, err)
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
			return fmt.Errorf("%w: write the install source proof temporary file: %v", ErrAIInstallSourceUnproven, writeErr)
		}
		if err := bound.LinkNoReplace(temporaryName, published); err != nil {
			_ = bound.Remove(temporaryName)
			if errors.Is(err, os.ErrExist) {
				return fmt.Errorf("%w: the install source proof %s already exists; a published proof is never replaced", ErrAIInstallSourceUnproven, path)
			}
			return fmt.Errorf("%w: publish the install source proof %s without replacing it: %v", ErrAIInstallSourceUnproven, path, err)
		}
		// The carrier is published at this point. A temporary name that cannot be
		// removed is a stray entry, not a failed publication, so the error is not
		// turned into a refusal and the directory sync below still runs.
		_ = bound.Remove(temporaryName)
		if err := evidence.VerifyBoundDir(bound, fileID); err != nil {
			return fmt.Errorf("%w: the install source proof directory %s is no longer the approved object: %v", ErrAIInstallSourceUnproven, directory, err)
		}
		// F10 (P85 design rev 3, section 4.4): the fsync runs on the SAME descriptor
		// the link was published through, and a failure is a refusal. The seam is the
		// shared one, so the counterexample can inject the failure on the first
		// publication too.
		if err := aiInstallSourceCarrierDirectorySync(bound); err != nil {
			return fmt.Errorf("%w: sync the install source proof directory %s: %v", ErrAIInstallSourceUnproven, directory, err)
		}
		return nil
	}
	return fmt.Errorf("%w: no unique install source proof temporary name could be created", ErrAIInstallSourceUnproven)
}

// aiInstallSourceProofTemporaryName returns one unguessable entry name. A name
// another writer could predict would let it pre-create the entry and turn our
// exclusive create into a refusal.
func aiInstallSourceProofTemporaryName() (string, error) {
	var buffer [8]byte
	if _, err := rand.Read(buffer[:]); err != nil {
		return "", err
	}
	return aiInstallSourceProofTemporaryPrefix + hex.EncodeToString(buffer[:]), nil
}
