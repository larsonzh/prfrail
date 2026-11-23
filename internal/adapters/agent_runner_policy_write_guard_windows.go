//go:build windows

package adapters

// F01, Windows half: the managed agent runs under this process's token, so the
// proof is "can THIS token obtain any mutation right on the policy directory".
//
// The probe is the T3 primitive family reused verbatim: one CreateFile per right
// per object (installSourceMutationAccess with the full directory table), so a
// directory that grants just one right - FILE_DELETE_CHILD alone, WRITE_DAC
// alone - is refused instead of being averaged away by a merged mask. A granted
// handle is closed immediately and nothing is written, so the probe leaves no
// footprint; only ERROR_ACCESS_DENIED counts as a denial and every other failure
// refuses.
//
// The directory is opened AS ITSELF first: a reparse point would make the name
// the kernel walks differ from the object that was judged.

import (
	"errors"
	"fmt"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// managedAgentPolicyDomainWriteProofSlots is the Windows half of F01 (design rev
// 3 section 1) with the fixed four-slot vector of design rev 4 section 2: one
// access check per mutation right, per PRESENT slot of the write domain (the
// policy directory, its parent, and every existing artifact/anchor file). A single
// granted right on any object refuses the whole domain, so a layout that grants
// FILE_DELETE_CHILD on the parent or WRITE_DATA on the artifact is refused instead
// of being averaged away. A granted handle is closed immediately and nothing is
// written, so the probe leaves no footprint.
//
// Each present object also contributes its identity (volume serial plus file
// index), so replacing a directory, a parent or a file changes the digest; an
// absent slot contributes the canonical empty marker and keeps the vector shape.
func managedAgentPolicyDomainWriteProofSlots(slots []managedAgentPolicyDomainSlot) (string, error) {
	sidText, err := installSourceCurrentUserSIDText()
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrManagedAgentPolicyWriteProof, err)
	}
	observed := make([]string, 0, len(slots))
	var failures []error
	for _, slot := range slots {
		if !slot.present {
			observed = append(observed, managedAgentPolicyDomainEmptySlot)
			continue
		}
		line, granted, objectErr := managedAgentPolicyDomainSlotObservation(slot, sidText)
		if line == "" {
			// P85 R7: an observation that could not render a line keeps its POSITION in
			// the vector and records the canonical empty marker, so a refusal digest still
			// carries the fixed four-slot shape instead of silently becoming shorter.
			line = managedAgentPolicyDomainEmptySlot
		}
		observed = append(observed, line)
		if objectErr != nil {
			failures = append(failures, objectErr)
			continue
		}
		if granted > 0 {
			failures = append(failures, fmt.Errorf("%w: the managed identity can obtain %d mutation right(s) on the %s %q",
				ErrManagedAgentPolicyWriteProof, granted, slot.label, slot.path))
		}
	}
	digest := evidence.Digest(managedAgentPolicyWriteProofDomain, []byte(strings.Join(observed, "\n")))
	if len(failures) > 0 {
		return digest, fmt.Errorf("%w: %v", ErrManagedAgentPolicyWriteProof, errors.Join(failures...))
	}
	return digest, nil
}

// managedAgentPolicyDomainSlotObservation judges one present slot and renders its
// observation line (which is returned even when the verdict refuses).
//
// P85 design rev 4, sections 1 and 3:
//   - a handle-form slot takes its identity from the handle, and the path must
//     still address that same object before AND after the probes (the identity
//     sandwich), so a name that moved off the handle refuses;
//   - the name CHAIN is classified with the same root-chain walk the verified
//     binding uses, so a reparse point anywhere in the chain refuses instead of
//     being followed by the kernel.
//
// P85 R7 (review r6 finding): every BY-NAME read of this observation - the
// per-right access checks and the DACL digest (GetNamedSecurityInfoW) - sits
// between the two chain classifications and the two identity comparisons, so the
// DACL that enters the digest describes the object the handle addresses.
func managedAgentPolicyDomainSlotObservation(slot managedAgentPolicyDomainSlot, sidText string) (string, int, error) {
	handleIdentity := ""
	if slot.dir != nil {
		identity, err := slot.dir.FileID()
		if err != nil {
			return "", 0, fmt.Errorf("%w: identify the %s handle %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		handleIdentity = string(identity)
	} else if slot.file != nil {
		identity, err := slot.file.FileID()
		if err != nil {
			return "", 0, fmt.Errorf("%w: identify the %s handle %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
		}
		handleIdentity = identity
	}
	if err := evidence.ClassifyRootChainNoReparse(slot.path, slot.isDir); err != nil {
		return "", 0, fmt.Errorf("%w: classify the root chain of the %s %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
	}
	before, err := evidence.FileIDAt(slot.path)
	if err != nil {
		return "", 0, fmt.Errorf("%w: identify the %s %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
	}
	if handleIdentity != "" && string(before) != handleIdentity {
		return "", 0, fmt.Errorf("%w: the %s %q no longer names the object its handle addresses (%s != %s)",
			ErrManagedAgentPolicyWriteProof, slot.label, slot.path, before, handleIdentity)
	}
	rights := installSourceFileMutationRights
	if slot.isDir {
		rights = installSourceDirectoryMutationRights
	}
	granted, denied, err := installSourceMutationAccess(slot.path, slot.isDir, rights)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", ErrManagedAgentPolicyWriteProof, err)
	}
	// P85 R7 (review r6 finding): the DACL digest is a BY-NAME read
	// (GetNamedSecurityInfoW), so it belongs INSIDE the sandwich. Taking it after the
	// closing comparison would describe a name that nothing re-checked afterwards.
	descriptor, err := installSourceDescriptorDigest(slot.path)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %v", ErrManagedAgentPolicyWriteProof, err)
	}
	// The closing half of the sandwich runs after every by-name read: the root chain
	// must still be reparse-free and the name must still address the object it started
	// on, so the recorded DACL describes the object that was judged.
	if err := evidence.ClassifyRootChainNoReparse(slot.path, slot.isDir); err != nil {
		return "", 0, fmt.Errorf("%w: re-classify the root chain of the %s %q after its access probes: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
	}
	after, err := evidence.FileIDAt(slot.path)
	if err != nil {
		return "", 0, fmt.Errorf("%w: re-identify the %s %q: %v", ErrManagedAgentPolicyWriteProof, slot.label, slot.path, err)
	}
	if after != before {
		return "", 0, fmt.Errorf("%w: the %s %q changed identity during its access probes (%s != %s)",
			ErrManagedAgentPolicyWriteProof, slot.label, slot.path, after, before)
	}
	identity := handleIdentity
	if identity == "" {
		identity = string(before)
	}
	line := strings.ToLower(slot.path) + "\n" +
		"id:" + identity + "\n" +
		"granted:" + installSourceRightsText(granted) + "\n" +
		"denied:" + installSourceRightsText(denied) + "\n" +
		"dacl:" + descriptor + "\n" +
		"trustee:" + sidText
	return line, len(granted), nil
}
