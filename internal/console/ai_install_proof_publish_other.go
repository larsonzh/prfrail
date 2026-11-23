//go:build !linux

package console

// The path-based carrier publication for every platform that is not Linux.
//
// P85 F20 is a Linux finding: there the directory can be bound and every
// publication step can be expressed relative to that descriptor, so the fixed
// shape lives in ai_install_proof_publish_linux.go. This half keeps the previous
// behaviour unchanged:
//
//   - Windows cannot create or hard-link relative to a directory handle with the
//     standard library, so the creation and link steps stay path-based. P85 R7
//     registers that residue as an INDEPENDENT S1 boundary of the user-mode path
//     publications (user ruling of 2026-11-23): it is not an F04 exemption and must
//     not be cited as one, and no durability is claimed from it. P85 R6 narrowed
//     what is left: the directory that receives the temporary file and the link is
//     re-bound through the root-chain walk before and after the publication (see
//     aiInstallSourceAssertDirectoryIdentity), so a chain that becomes a reparse
//     point refuses.
//   - darwin and the remaining unix targets already refuse the install-source proof
//     before this point (adapters.VerifyAIInstallSource fails closed there), so this
//     function is unreachable on them.
//
// The no-replace discipline is identical on both halves: the payload is written and
// synced in a same-directory temporary file, claimed with a hard link that fails
// when the name exists, and the directory is synced. There is no os.Rename fallback
// (F09: os.Rename replaces an existing name, so a fallback would silently turn
// "never replaced" into "last writer wins").

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// writeAIInstallSourceProofFileIn publishes one payload without replacing an
// existing name: the payload is written and synced in a same-directory temporary
// file, then claimed with a hard link, and the directory is synced.
//
// P85 F10 (design rev 3, section 4): the caller passes the descriptor it already
// read the carrier through, so this half never has to decide "which directory"
// again. The create/link steps stay path-based (Windows cannot create or link
// relative to a directory handle through Go); that residue is the registered S1
// user-mode path boundary, and the identity of the bound object is re-proved around
// them, so a swap is refused instead of receiving the carrier.
func writeAIInstallSourceProofFileIn(bound *evidence.BoundDir, fileID evidence.FileID, path string, payload []byte) error {
	directory := filepath.Dir(path)
	if err := aiInstallSourceAssertDirectoryIdentity(bound, fileID, directory); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(directory, ".ai-install-source-proof-*")
	if err != nil {
		return fmt.Errorf("%w: create the install source proof temporary file: %v", ErrAIInstallSourceUnproven, err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
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
		return fmt.Errorf("%w: write the install source proof temporary file: %v", ErrAIInstallSourceUnproven, writeErr)
	}
	if err := os.Link(temporaryPath, path); err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%w: the install source proof %s already exists; a published proof is never replaced", ErrAIInstallSourceUnproven, path)
		}
		return fmt.Errorf("%w: publish the install source proof %s without replacing it: %v", ErrAIInstallSourceUnproven, path, err)
	}
	if err := aiInstallSourceAssertDirectoryIdentity(bound, fileID, directory); err != nil {
		return err
	}
	// Atomic visibility is not power-loss durability: the hard link above is the
	// durability point this platform can offer (F10, registered; the path-based link
	// itself is the registered S1 user-mode boundary). The seam is
	// still used so both halves report a failing sync the same way.
	return aiInstallSourceCarrierDirectorySync(bound)
}

// aiInstallSourceAssertDirectoryIdentity proves that the name the caller derived
// still resolves to the object the verified binding produced.
//
// P85 R6 (review r3 finding): the name is re-BOUND through the same root-chain walk
// the initial binding uses (evidence.OpenBoundDirVerified -> openBoundDirChainHandle
// on Windows: every component below the volume root opened as itself, so a junction,
// symbolic link or any other reparse point refuses), and the rebound object's
// identity is compared with the approved one. The previous form used
// evidence.FileIDAt, whose stat FOLLOWS a reparse point, so a chain that became a
// link after the bind could still report the approved identity - the publish
// confirmation and the M14 test measured different paths.
func aiInstallSourceAssertDirectoryIdentity(bound *evidence.BoundDir, fileID evidence.FileID, directory string) error {
	if err := evidence.VerifyBoundDir(bound, fileID); err != nil {
		return fmt.Errorf("%w: the install source proof directory %s is no longer the approved object: %v", ErrAIInstallSourceUnproven, directory, err)
	}
	rebound, reboundID, err := evidence.OpenBoundDirVerified(directory)
	if err != nil {
		return fmt.Errorf("%w: the install source proof directory %s can no longer be re-bound without following a link: %v", ErrAIInstallSourceUnproven, directory, err)
	}
	defer func() { _ = rebound.Close() }()
	if reboundID != fileID {
		return fmt.Errorf("%w: the install source proof directory %s no longer names the object that was bound (%s != %s)", ErrAIInstallSourceUnproven, directory, reboundID, fileID)
	}
	return nil
}
