package console

// T3 (finding F6): the console side of the install-source proof.
//
// The carrier has exactly one allowed location: the trusted root derived from the
// chain-configuration directory, the same durable root the probe's authorization
// ledger lives in. The caller never chooses the path - deriving it from the
// invocation directory would let a probe be admitted by a proof the caller wrote
// somewhere convenient, which is the self-approval hole this item exists to close.
//
// The carrier is never rewritten once published: an existing proof is re-read,
// re-validated and compared with a LIVE observation, and only an observation that
// still describes it converges. Anything else - a changed digest or identity, a
// DACL that now grants write access, a directory that became writable - is a
// refusal, because the proof exists to make exactly that change visible.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// ErrAIInstallSourceUnproven reports an installation source that may not be
// admitted. It is the cli_unavailable classification of the probe preflight: no
// credential is resolved and no process starts after it.
var ErrAIInstallSourceUnproven = errors.New("the probe install source cannot be admitted")

// aiInstallSourceProofPath derives the carrier path from the trusted root the chain
// configuration lives in. The derivation is the ledger's: same root, same
// discipline, no caller-selected path.
func aiInstallSourceProofPath(chainConfigDir string) (string, error) {
	path, err := trustedLedgerPath(chainConfigDir, adapters.AIInstallSourceProofFileName, "")
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrAIInstallSourceUnproven, err)
	}
	return path, nil
}

// recordOrConfirmAIInstallSourceProof observes the executable right now and either
// publishes the carrier (first use) or confirms that the published carrier still
// describes the observation.
//
// The observation always happens FIRST. A carrier can never admit a source that
// fails the live check, so a stale carrier cannot outlive the fact it records.
//
// F03 (P85 fix round 1): approvedExecutableHash is the digest the caller's
// receiver binding (or checker input) already approved. When it is non-empty the
// LIVE observation must carry that same digest before anything is persisted, so
// an unauthorized file can never be recorded as the install source and only
// rejected later by the receiver binding. The comparison uses the same
// executable domain the console binding hashes with
// (adapters.CopilotCLIExecutableHash).
//
// P85 F05: the function now returns the CARRIER's facts, not just an error. The
// carrier is the normative description of the object the operator approved, so it
// is what the executor compares its live observation against at the spawn boundary.
// A freshly re-derived observation would describe "whatever is at the path now" and
// silently re-approve a swapped object; returning the recorded facts closes that gap.
//
// P85 F15: the refusal text carries no free-text "(cli_unavailable)" marker any
// more. The classification is emitted as the structured errorClass field by the
// console command, so the message stays a plain detail string.
func recordOrConfirmAIInstallSourceProof(path, executable, approvedExecutableHash string, now time.Time) (adapters.AIInstallSourceFacts, error) {
	if strings.TrimSpace(path) == "" {
		return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: the install source proof path is empty", ErrAIInstallSourceUnproven)
	}
	facts, err := adapters.VerifyAIInstallSource(executable)
	if err != nil {
		return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: %v", ErrAIInstallSourceUnproven, err)
	}
	// F03: the live digest has to be the approved one BEFORE the carrier exists.
	// A refusal here leaves no proof behind, so the unauthorized source is never
	// published and never has to be caught by a later binding check.
	if approvedExecutableHash != "" && facts.ExecutableHash != approvedExecutableHash {
		return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: the live install source is not the approved executable", ErrAIInstallSourceUnproven)
	}
	directory := filepath.Dir(path)
	// F10 (P85 design rev 3, section 4): the carrier directory is bound ONCE for
	// the whole entry - read, confirmation and fsync all act on this handle - so
	// the three independent opens the previous shape performed cannot describe
	// three different directories.
	bound, fileID, err := aiInstallSourceBindCarrierDirectory(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return adapters.AIInstallSourceFacts{}, err
		}
		// The carrier directory does not exist yet: create it, then bind it once.
		if mkErr := os.MkdirAll(directory, 0o755); mkErr != nil {
			return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: create the install source proof directory: %v", ErrAIInstallSourceUnproven, mkErr)
		}
		bound, fileID, err = aiInstallSourceBindCarrierDirectory(path)
		if err != nil {
			return adapters.AIInstallSourceFacts{}, err
		}
	}
	defer func() { _ = bound.Close() }()

	existing, found, err := readAIInstallSourceProofFrom(bound, filepath.Base(path), adapters.AIInstallSourceProofByteBound)
	if err != nil {
		return adapters.AIInstallSourceFacts{}, err
	}
	if found {
		// F10 (P85 fix round 3): an existing carrier is only admissible after its
		// DIRECTORY has been proven and made persistent on this entry. The previous
		// publish hard-links the carrier into place and then syncs the directory; if
		// the link became visible but that sync failed, the earlier call refused -
		// and this re-entry would otherwise adopt a carrier whose directory entry is
		// not durable, silently undoing the refusal. Binding + re-confirming the
		// directory and repeating the sync here means the re-entry either re-proves
		// persistence or refuses.
		//
		// P85 design rev 3 section 4 collapses the three independent opens of this
		// branch (read carrier, confirm directory, sync directory) onto ONE bound
		// descriptor: the carrier is read relative to the handle, the identity is
		// asserted on that handle, and the fsync acts on it. "The directory that was
		// confirmed" and "the directory that was synced" are now one object instead
		// of a claim, and a sync failure is still a refusal rather than an adoption.
		//
		// Windows durability boundary (independent S1 registration, user ruling of
		// 2026-11-23, not an F04 exemption): the hard link gives ATOMIC
		// VISIBILITY - the name appears complete or not at all - but Windows cannot
		// sync a directory handle through Go, so this layer grants no power-loss
		// durability there. The two facts are stated apart instead of being merged
		// into one "atomic therefore durable" claim.
		if err := aiInstallSourceConfirmBoundDirectory(bound, fileID, directory); err != nil {
			return adapters.AIInstallSourceFacts{}, err
		}
		if syncErr := aiInstallSourceCarrierDirectorySync(bound); syncErr != nil && aiInstallSourceReentryRequiresSync {
			return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: sync the install source proof directory %s: %v", ErrAIInstallSourceUnproven, directory, syncErr)
		}
		if !facts.Matches(existing) {
			return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: the recorded install source proof %s no longer describes %q",
				ErrAIInstallSourceUnproven, path, facts.ExecutablePath)
		}
		return existing.Facts(), nil
	}
	proof := adapters.NewAIInstallSourceProof(facts, now)
	payload, err := adapters.EncodeAIInstallSourceProof(proof)
	if err != nil {
		return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: %v", ErrAIInstallSourceUnproven, err)
	}
	// F10 (P85 design rev 3, section 4): the first publication runs inside the same
	// handle the read used - create, write, sync, link, remove, fsync - so the
	// directory that is proven is the directory whose entry becomes durable, and a
	// failed fsync propagates as a refusal instead of being adopted on re-entry.
	if err := writeAIInstallSourceProofFileIn(bound, fileID, path, payload); err != nil {
		return adapters.AIInstallSourceFacts{}, err
	}
	published, found, err := readAIInstallSourceProofFrom(bound, filepath.Base(path), adapters.AIInstallSourceProofByteBound)
	if err != nil {
		return adapters.AIInstallSourceFacts{}, err
	}
	if !found {
		return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: the install source proof %s is missing right after it was published", ErrAIInstallSourceUnproven, path)
	}
	publishedPayload, err := adapters.EncodeAIInstallSourceProof(published)
	if err != nil {
		return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: %v", ErrAIInstallSourceUnproven, err)
	}
	if adapters.AIInstallSourceProofDigest(publishedPayload) != adapters.AIInstallSourceProofDigest(payload) {
		return adapters.AIInstallSourceFacts{}, fmt.Errorf("%w: the install source proof %s does not read back as the proof that was published", ErrAIInstallSourceUnproven, path)
	}
	return published.Facts(), nil
}

// readAIInstallSourceProofFrom reads and validates one carrier RELATIVE to a
// directory handle the caller already holds (P85 design rev 3, section 4).
//
// Splitting the read off the binding is what lets one entry use one descriptor for
// everything: the carrier bytes, the identity assertion and the fsync all describe
// the same directory object, so no later step can rediscover a different one.
//
// found=false means "no carrier at that name"; every other failure - unreadable,
// oversized, malformed, failing validation - is an error, so a carrier that cannot
// be read never degrades into "recompute and continue".
func readAIInstallSourceProofFrom(bound *evidence.BoundDir, base string, byteBound int64) (adapters.AIInstallSourceProof, bool, error) {
	file, err := evidence.OpenBoundDirFile(bound, base, byteBound)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return adapters.AIInstallSourceProof{}, false, nil
		}
		return adapters.AIInstallSourceProof{}, false, fmt.Errorf("%w: the install source proof %s cannot be opened inside its bound directory: %v", ErrAIInstallSourceUnproven, base, err)
	}
	defer func() { _ = file.Close() }()
	data, truncated, err := file.ReadAll()
	if err != nil {
		return adapters.AIInstallSourceProof{}, false, fmt.Errorf("%w: the install source proof %s cannot be read: %v", ErrAIInstallSourceUnproven, base, err)
	}
	if truncated {
		return adapters.AIInstallSourceProof{}, false, fmt.Errorf("%w: the install source proof %s exceeds its byte bound", ErrAIInstallSourceUnproven, base)
	}
	proof, err := adapters.DecodeAIInstallSourceProof(data)
	if err != nil {
		return adapters.AIInstallSourceProof{}, false, fmt.Errorf("%w: the install source proof %s is invalid: %v", ErrAIInstallSourceUnproven, base, err)
	}
	return proof, true, nil
}

// aiInstallSourceBindCarrierDirectory binds the derived carrier directory ONCE and
// returns the handle together with the identity the binding produced (P85 design
// rev 3, section 4.1). Every later step of the entry - reading the carrier,
// confirming the directory, syncing it, publishing it - uses this handle, so
// "which directory" is answered once instead of three times.
func aiInstallSourceBindCarrierDirectory(path string) (*evidence.BoundDir, evidence.FileID, error) {
	directory := filepath.Dir(path)
	bound, fileID, err := evidence.OpenBoundDirVerified(directory)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, "", err
		}
		return nil, "", fmt.Errorf("%w: the install source proof directory %s cannot be bound without following a symbolic link: %v", ErrAIInstallSourceUnproven, directory, err)
	}
	return bound, fileID, nil
}

// aiInstallSourceConfirmBoundDirectory asserts, on the handle the caller already
// holds, that it still addresses the object the binding produced (P85 design rev 3,
// sections 2 and 4). It replaces the previous "bind the path again and compare" -
// which compared two resolutions of one name instead of the approved object - and
// it opens nothing.
func aiInstallSourceConfirmBoundDirectory(bound *evidence.BoundDir, fileID evidence.FileID, directory string) error {
	if err := evidence.VerifyBoundDir(bound, fileID); err != nil {
		return fmt.Errorf("%w: the install source proof directory %s is no longer the approved object: %v", ErrAIInstallSourceUnproven, directory, err)
	}
	return nil
}

// aiInstallSourceCarrierDirectorySync is the ONE fsync every carrier path uses. It
// is a package variable so the F10 counterexample can inject a failing sync (the
// sync guard is otherwise unreachable from outside), and it is NEVER read from
// configuration: production always uses (*evidence.BoundDir).SyncDirectory.
var aiInstallSourceCarrierDirectorySync = (*evidence.BoundDir).SyncDirectory

// aiInstallSourceReentryRequiresSync is the F10 observability seam: it is always
// true in production, and only an in-package test sets it false to prove that a
// re-entry refuses BECAUSE of the sync guard and not because the carrier read
// fails first (P85 design rev 3, section 4.5).
var aiInstallSourceReentryRequiresSync = true

// writeAIInstallSourceProofFileIn publishes one payload without replacing an
// existing name. The platform halves own the "how":
//
//   - Linux (ai_install_proof_publish_linux.go) creates the temporary file, writes,
//     syncs, hard-links it to the final name, removes the temporary name and fsyncs
//     the descriptor, all INSIDE one bound directory descriptor (P85 F20: the
//     directory that was proven is the directory the bytes land in, and no name is
//     resolved again);
//   - every other platform keeps the path-based publication with an identity
//     assertion around it (ai_install_proof_publish_other.go). Windows cannot
//     create or link relative to a directory handle with the standard library, and
//     those two steps are the registered S1 user-mode path boundary (user ruling of
//     2026-11-23), not an F04 exemption; darwin and the remaining unix targets
//     already refuse the install-source proof outright.
//
// Persisting is the linearization point: only after the carrier is in place and its
// directory is synced may the check run.
