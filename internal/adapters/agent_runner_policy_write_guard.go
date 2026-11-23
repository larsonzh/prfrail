package adapters

// T2 / F01 (P85 fix round 3): the write guard of the frozen-policy directory.
//
// The assembly anchor proves that the artifact on disk is still the object the
// assembly froze, but it lives in the same directory as the artifact. An anchor
// can therefore only be an integrity record, never a boundary: a writer that may
// write the artifact may write its anchor too. CONTRACTS.md:230 says the managed
// CLI agent "must not write source trees, the store, POLICIES, gate
// implementations or acceptance records", so the assembly has to prove that the
// identity the managed agent runs under cannot write the policy directory at all.
//
// The identity is the one the launcher actually uses: guard.StartManaged spawns
// the pinned CLI as the invoking process's own token/user (S1 has no privilege
// separation - see docs/SECURITY.md, the trust boundary section). Asking the
// KERNEL what THIS token may do to the directory therefore asks exactly what the
// managed agent may do to it, and the answer is the proof - not a path
// disjointness argument, and not the "same OS user is not a malicious process
// boundary" reading that SECURITY.md:11 reserves for a different question.
//
// On S1 the managed identity IS the operator's identity, and the assembly creates
// the policy directory itself, so the owner holds the mutation rights and the
// proof cannot succeed. That is the intended outcome, not a bug: the assembly
// refuses (fail-closed) and no dispatch happens. "Cannot prove it" is "refuse it".
//
// The platform halves own the "how":
//
//   - Windows (agent_runner_policy_write_guard_windows.go): one access check per
//     mutation right against the directory with this process token, reusing the
//     T3 primitive family (installSourceMutationAccess and its right tables);
//   - Linux (agent_runner_policy_write_guard_unix.go): mode bits, ownership and
//     the POSIX access-ACL rule (a named grant cannot be excluded), reusing the
//     F08 verdict of the install-source proof;
//   - every other platform fails closed: an unproven layout is refused.

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// ErrManagedAgentPolicyWriteProof reports a frozen-policy directory the managed
// agent's own identity could obtain a mutation right on, or one whose protection
// cannot be determined at all. It is raised at assembly, before anything is
// dispatched.
var ErrManagedAgentPolicyWriteProof = errors.New("the managed agent can write the frozen policy directory")

// managedAgentPolicyDomainEmptySlot is the canonical rendering of an absent slot.
// It keeps the vector shape fixed, so "no artifact yet" and "no anchor yet" are
// positions in a stable sequence instead of a shorter sequence that could collide
// with a different domain.
const managedAgentPolicyDomainEmptySlot = "∅"

// ManagedPolicyDomainBoundFile is one file slot of the handle-form observation:
// the handle the bytes are read through, plus the path the path-based probes (the
// POSIX ACL read, the Windows per-right access checks, the root-chain classifier)
// address. The identity of the slot is taken from Bound, never from Path.
type ManagedPolicyDomainBoundFile struct {
	// Path is the filesystem name the probes address.
	Path string
	// Role names the slot in refusals ("artifact", "anchor").
	Role string
	// Bound is the open handle whose identity is the slot's identity.
	Bound *evidence.BoundFile
}

// VerifyManagedAgentCannotWritePolicyDomain proves that the identity the managed
// agent runs under cannot obtain a mutation right on the frozen-policy DOMAIN, and
// returns the stable digest of the observation for the assembly anchor.
//
// P85 design rev 3 section 1 introduced the domain (the policy directory, its
// PARENT, and every existing artifact/anchor file). P85 design rev 4 section 2
// fixes its SHAPE: the observed vector always has four slots -
//
//	[policy directory | parent directory | artifact | anchor]
//
// - and a slot whose object is not part of the observed set is rendered with the
// canonical empty marker instead of being omitted. The two file slots are placed by
// ROLE, never by path order (P85 R6): the artifact is always slot 3 and the anchor
// always slot 4, so a set that carries only the anchor leaves the artifact slot
// empty instead of sliding the anchor into it.
//
// This is the NAME-based entry: it takes its identities from the filesystem names.
// Callers that hold handles must use VerifyManagedAgentCannotWritePolicyDomainHandles,
// which takes the identities from the very handles the bytes are read and written
// through (P85 design rev 4 section 1).
//
// The returned digest is produced even when the proof FAILS, so a refusal can
// carry the observation it refused on; the caller must treat a non-nil error as a
// refusal regardless of the digest.
func VerifyManagedAgentCannotWritePolicyDomain(directory string, existing []string) (string, error) {
	if strings.TrimSpace(directory) == "" {
		return "", fmt.Errorf("%w: the policy directory is empty", ErrManagedAgentPolicyWriteProof)
	}
	slots, err := managedAgentPolicyDomainNameSlots(directory, existing)
	if err != nil {
		return "", err
	}
	return managedAgentPolicyDomainWriteProofSlots(slots)
}

// VerifyManagedAgentCannotWritePolicyDomainHandles is the HANDLE-form observation
// of P85 design rev 4 section 1: every identity in the vector comes from an open
// handle - the two bound directories and the bound artifact/anchor files - while
// the path-based probes (POSIX ACL read on Linux, per-right access checks and the
// root-chain classifier on Windows) run inside an identity sandwich, so a name
// that no longer addresses the handle's object refuses instead of judging another
// object.
//
// files carries at most the artifact and the anchor. The vector keeps the fixed
// four-slot shape, so a missing file is the canonical empty slot - and each file
// goes to the slot its ROLE names (P85 R6), never to the slot a path sort would
// give it.
func VerifyManagedAgentCannotWritePolicyDomainHandles(directory *evidence.BoundDir, directoryID evidence.FileID, parent *evidence.BoundDir, parentID evidence.FileID, files []ManagedPolicyDomainBoundFile) (string, error) {
	slots, err := managedAgentPolicyDomainHandleSlots(directory, directoryID, parent, parentID, files)
	if err != nil {
		return "", err
	}
	return managedAgentPolicyDomainWriteProofSlots(slots)
}

// managedAgentPolicyDomainSlot is one slot of the fixed four-slot vector. A slot
// that is not part of the observed set has present=false and contributes the
// canonical empty marker; dir/file is non-nil only on the handle form.
type managedAgentPolicyDomainSlot struct {
	// label names the role in refusals ("policy directory", "parent directory",
	// "artifact", "anchor") so an operator can tell which slot failed the proof.
	label   string
	path    string
	isDir   bool
	present bool
	dir     *evidence.BoundDir
	file    *evidence.BoundFile
}

// The role markers of the fixed four-slot vector (P85 R6): the file slots are
// addressed by the ROLE the frozen policy contract gives them, never by their
// position in a sorted path list. The name suffixes are the role markers the
// contract itself spells (the artifact is <task>.<step>.<attempt>.policy.json and
// the anchor <task>.<step>.<attempt>.anchor.json), so a name carries its role even
// when the caller passes no role field.
const (
	managedPolicyDomainArtifactRole   = "artifact"
	managedPolicyDomainAnchorRole     = "anchor"
	managedPolicyDomainArtifactSuffix = ".policy.json"
	managedPolicyDomainAnchorSuffix   = ".anchor.json"
)

// managedPolicyDomainFileSlotCount is the number of FILE slots of the fixed
// vector: [policy directory | parent directory | artifact | anchor].
const managedPolicyDomainFileSlotCount = 2

// managedPolicyDomainFileCandidate is one file slot before its role position is
// resolved. role is 0 for the artifact slot, 1 for the anchor slot and -1 when
// neither the caller's role field nor the path suffix names a role.
type managedPolicyDomainFileCandidate struct {
	role int
	slot managedAgentPolicyDomainSlot
}

// managedPolicyDomainRoleFromRoleMarker maps a declared role onto its slot index.
func managedPolicyDomainRoleFromRoleMarker(role string) int {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case managedPolicyDomainArtifactRole:
		return 0
	case managedPolicyDomainAnchorRole:
		return 1
	}
	return -1
}

// managedPolicyDomainRoleFromName maps a file NAME onto its slot index through the
// role markers the frozen policy contract spells. Matching is on the suffix, so the
// task/step/attempt prefix never has to be parsed here.
//
// The comparison is the contract's lowercase spelling (P85 R7): a name that matches
// only case-insensitively is refused by managedPolicyDomainResolveFileRole instead
// of being placed.
func managedPolicyDomainRoleFromName(path string) int {
	base := filepath.Base(path)
	switch {
	case strings.HasSuffix(base, managedPolicyDomainAnchorSuffix):
		return 1
	case strings.HasSuffix(base, managedPolicyDomainArtifactSuffix):
		return 0
	}
	return -1
}

// managedPolicyDomainRoleFromNameIgnoringCase reports the role a name would carry
// if the contract's suffix were matched case-insensitively. It never PLACES a slot;
// it exists only so a differently-cased suffix can be refused explicitly.
func managedPolicyDomainRoleFromNameIgnoringCase(path string) int {
	base := strings.ToLower(filepath.Base(path))
	switch {
	case strings.HasSuffix(base, managedPolicyDomainAnchorSuffix):
		return 1
	case strings.HasSuffix(base, managedPolicyDomainArtifactSuffix):
		return 0
	}
	return -1
}

// managedPolicyDomainResolveFileRole resolves one handle-form file onto its
// contract role (P85 R7).
//
// The declared role wins only when the name AGREES with it. An explicit
// Role:"artifact" on a ".anchor.json" name - or the mirror image - describes two
// different slots for one object; picking either silently would place the object in
// the slot the other side of the contract does not expect, so the pair is refused.
//
// A name whose suffix matches only case-insensitively is not the spelling the
// contract freezes, so it is refused instead of being either accepted or left to
// path order. A name with no role marker at all keeps role -1 (the caller places it
// in the first free slot, deterministically by path).
func managedPolicyDomainResolveFileRole(declared, path string) (int, error) {
	declaredRole := managedPolicyDomainRoleFromRoleMarker(declared)
	nameRole := managedPolicyDomainRoleFromName(path)
	// A suffix that matches only case-insensitively is refused with or without a
	// declared role: the contract freezes the lowercase spelling, and accepting the
	// other spelling would let a name mean "the anchor" to one reader and "no role"
	// to another.
	if nameRole < 0 {
		if casing := managedPolicyDomainRoleFromNameIgnoringCase(path); casing >= 0 {
			return -1, fmt.Errorf("%w: the %q slot carries the %s role suffix in a spelling other than the contract's lowercase form",
				ErrManagedAgentPolicyWriteProof, path, managedPolicyDomainFileSlotLabel(casing))
		}
	}
	switch {
	case declaredRole >= 0 && nameRole >= 0 && declaredRole != nameRole:
		return -1, fmt.Errorf("%w: the %q slot declares the %s role, but its name is the %s role",
			ErrManagedAgentPolicyWriteProof, path,
			managedPolicyDomainFileSlotLabel(declaredRole), managedPolicyDomainFileSlotLabel(nameRole))
	case declaredRole >= 0:
		return declaredRole, nil
	case nameRole >= 0:
		return nameRole, nil
	}
	return -1, nil
}

// managedPolicyDomainFileSlotLabel names a file slot by its ROLE, so a refusal
// says which position of the contract failed instead of "some existing file".
func managedPolicyDomainFileSlotLabel(index int) string {
	if index == 1 {
		return managedPolicyDomainAnchorRole
	}
	return managedPolicyDomainArtifactRole
}

// managedAgentPolicyDomainPlaceFileSlots folds the file candidates into the two
// ROLE-ORDERED file slots of the four-slot vector (P85 R6).
//
// The contract fixes both the slot ORDER and the slot COUNT: the artifact slot is
// slot 3 and the anchor slot slot 4 no matter how many objects exist, and a slot
// with no object stays in the vector as the canonical empty marker. Sorting the
// files by path - the previous shape - put .anchor.json ahead of .policy.json and
// slid the anchor into the artifact slot when only the anchor existed, which made
// the observed vector differ from the vector the contract describes.
//
// A path that carries no role marker is not part of the frozen policy naming
// contract; it is placed in the first free slot in path order so the shape stays
// deterministic, and it can never displace a role-marked slot.
func managedAgentPolicyDomainPlaceFileSlots(base []managedAgentPolicyDomainSlot, candidates []managedPolicyDomainFileCandidate) ([]managedAgentPolicyDomainSlot, error) {
	slots := make([]managedAgentPolicyDomainSlot, len(base), len(base)+managedPolicyDomainFileSlotCount)
	copy(slots, base)
	for index := 0; index < managedPolicyDomainFileSlotCount; index++ {
		slots = append(slots, managedAgentPolicyDomainSlot{label: managedPolicyDomainFileSlotLabel(index)})
	}
	unplaced := make([]managedAgentPolicyDomainSlot, 0, len(candidates))
	for _, candidate := range candidates {
		if candidate.role < 0 {
			unplaced = append(unplaced, candidate.slot)
			continue
		}
		position := len(base) + candidate.role
		if slots[position].present {
			return nil, fmt.Errorf("%w: the %s slot of the frozen policy domain is claimed by both %q and %q",
				ErrManagedAgentPolicyWriteProof, managedPolicyDomainFileSlotLabel(candidate.role), slots[position].path, candidate.slot.path)
		}
		candidate.slot.label = managedPolicyDomainFileSlotLabel(candidate.role)
		slots[position] = candidate.slot
	}
	sort.Slice(unplaced, func(i, j int) bool { return unplaced[i].path < unplaced[j].path })
	for _, slot := range unplaced {
		placed := false
		for index := 0; index < managedPolicyDomainFileSlotCount; index++ {
			if slots[len(base)+index].present {
				continue
			}
			slot.label = managedPolicyDomainFileSlotLabel(index)
			slots[len(base)+index] = slot
			placed = true
			break
		}
		if !placed {
			return nil, fmt.Errorf("%w: the frozen policy domain carries at most two files (artifact and anchor)", ErrManagedAgentPolicyWriteProof)
		}
	}
	return slots, nil
}

// managedAgentPolicyDomainNameSlots builds the four-slot vector from names: the
// policy directory, its parent, and the existing entries placed by role.
func managedAgentPolicyDomainNameSlots(directory string, existing []string) ([]managedAgentPolicyDomainSlot, error) {
	policy := filepath.Clean(directory)
	if !filepath.IsAbs(policy) {
		return nil, fmt.Errorf("%w: %q is not a stable absolute path", ErrManagedAgentPolicyWriteProof, policy)
	}
	if len(existing) > 2 {
		return nil, fmt.Errorf("%w: the frozen policy domain carries at most two files (artifact and anchor), got %d", ErrManagedAgentPolicyWriteProof, len(existing))
	}
	base := []managedAgentPolicyDomainSlot{
		{label: "policy directory", path: policy, isDir: true, present: true},
	}
	if parent := filepath.Dir(policy); parent != policy {
		// The parent gates rename/unlink of the policy directory entry itself:
		// DELETE_CHILD or a rename right there replaces the directory no matter how
		// well protected the directory's own DACL is.
		base = append(base, managedAgentPolicyDomainSlot{label: "parent directory", path: parent, isDir: true, present: true})
	} else {
		// A root directory has no distinct parent; the slot stays in the vector so the
		// shape never changes.
		base = append(base, managedAgentPolicyDomainSlot{label: "parent directory", isDir: true})
	}
	candidates := make([]managedPolicyDomainFileCandidate, 0, len(existing))
	for _, name := range existing {
		if strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("%w: an existing policy file name is empty", ErrManagedAgentPolicyWriteProof)
		}
		clean := filepath.Clean(name)
		candidates = append(candidates, managedPolicyDomainFileCandidate{
			role: managedPolicyDomainRoleFromName(clean),
			slot: managedAgentPolicyDomainSlot{path: clean, isDir: false, present: true},
		})
	}
	return managedAgentPolicyDomainPlaceFileSlots(base, candidates)
}

// managedAgentPolicyDomainHandleSlots builds the same four-slot vector from
// handles. The two directory identities are asserted against the handles here, so
// a caller cannot pair a handle with a mismatched recorded identity.
func managedAgentPolicyDomainHandleSlots(directory *evidence.BoundDir, directoryID evidence.FileID, parent *evidence.BoundDir, parentID evidence.FileID, files []ManagedPolicyDomainBoundFile) ([]managedAgentPolicyDomainSlot, error) {
	if directory == nil || parent == nil {
		return nil, fmt.Errorf("%w: the handle form of the policy domain requires both directory handles", ErrManagedAgentPolicyWriteProof)
	}
	if len(files) > 2 {
		return nil, fmt.Errorf("%w: the frozen policy domain carries at most two files (artifact and anchor), got %d", ErrManagedAgentPolicyWriteProof, len(files))
	}
	directoryPath, err := directory.RootPath()
	if err != nil {
		return nil, fmt.Errorf("%w: read the policy directory handle name: %v", ErrManagedAgentPolicyWriteProof, err)
	}
	parentPath, err := parent.RootPath()
	if err != nil {
		return nil, fmt.Errorf("%w: read the parent directory handle name: %v", ErrManagedAgentPolicyWriteProof, err)
	}
	if actual, err := directory.FileID(); err != nil {
		return nil, fmt.Errorf("%w: identify the policy directory handle: %v", ErrManagedAgentPolicyWriteProof, err)
	} else if actual != directoryID {
		return nil, fmt.Errorf("%w: the policy directory handle is not the recorded object (%s != %s)", ErrManagedAgentPolicyWriteProof, actual, directoryID)
	}
	if actual, err := parent.FileID(); err != nil {
		return nil, fmt.Errorf("%w: identify the parent directory handle: %v", ErrManagedAgentPolicyWriteProof, err)
	} else if actual != parentID {
		return nil, fmt.Errorf("%w: the parent directory handle is not the recorded object (%s != %s)", ErrManagedAgentPolicyWriteProof, actual, parentID)
	}
	base := []managedAgentPolicyDomainSlot{
		{label: "policy directory", path: filepath.Clean(directoryPath), isDir: true, present: true, dir: directory},
		{label: "parent directory", path: filepath.Clean(parentPath), isDir: true, present: true, dir: parent},
	}
	// P85 R6: the slot position comes from the ROLE (the caller's role field, then
	// the contract's name suffix), never from the path order.
	candidates := make([]managedPolicyDomainFileCandidate, 0, len(files))
	for _, file := range files {
		if file.Bound == nil {
			return nil, fmt.Errorf("%w: the %s slot carries no bound handle", ErrManagedAgentPolicyWriteProof, file.Role)
		}
		if strings.TrimSpace(file.Path) == "" {
			return nil, fmt.Errorf("%w: the %s slot carries no path", ErrManagedAgentPolicyWriteProof, file.Role)
		}
		clean := filepath.Clean(file.Path)
		role, err := managedPolicyDomainResolveFileRole(file.Role, clean)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, managedPolicyDomainFileCandidate{
			role: role,
			slot: managedAgentPolicyDomainSlot{path: clean, isDir: false, present: true, file: file.Bound},
		})
	}
	return managedAgentPolicyDomainPlaceFileSlots(base, candidates)
}
