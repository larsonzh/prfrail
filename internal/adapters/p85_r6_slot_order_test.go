package adapters

// P85 R6 (test round 6) / review r3 finding 2: the observed vector of the
// managed-identity write proof must be the FOUR fixed slots of the contract -
//
//	[policy directory | parent directory | artifact | anchor]
//
// with the absent slot rendered as the canonical empty marker - and the two file
// slots must be placed by ROLE, never by path order. The review found that
// `.anchor.json` sorts before `.policy.json` and that an anchor-only set slid the
// anchor into the artifact slot, so the vector that was actually hashed differed
// from the vector CONTRACTS.md:233 describes.
//
// The white-box assertions below read the slot list the production helpers return,
// so the property is pinned where it is produced and not through a digest whose
// preimage a later refactor could change.
//
// Mutation check (design rev 4 section 6, mutation 1): making
// managedAgentPolicyDomainPlaceFileSlots ignore the resolved role and append the
// candidates in path order puts the anchor third and the artifact fourth; the
// full-object-vector, anchor-only and handle-form legs then go red.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85R6DescribeSlots renders one four-slot vector as "label=path" with the
// canonical empty marker for an absent slot, so an assertion can state the whole
// vector literally.
func p85R6DescribeSlots(slots []managedAgentPolicyDomainSlot) []string {
	described := make([]string, 0, len(slots))
	for _, slot := range slots {
		if !slot.present {
			described = append(described, slot.label+"="+managedAgentPolicyDomainEmptySlot)
			continue
		}
		described = append(described, slot.label+"="+slot.path)
	}
	return described
}

// p85R6RequireSlots compares one vector against the expected slots position by
// position, so an order swap cannot pass by comparing sets.
func p85R6RequireSlots(t *testing.T, label string, got []string, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s must keep the fixed %d-slot vector, got %d slots: %v", label, len(want), len(got), got)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("%s: slot %d must be %q, got %q (whole vector %v)", label, index+1, want[index], got[index], got)
		}
	}
}

// TestP85R6SlotOrderMatchesTheContract pins the real vector that enters the digest.
func TestP85R6SlotOrderMatchesTheContract(t *testing.T) {
	// The two names are the ones the frozen policy contract spells. The anchor name
	// sorts BEFORE the artifact name, which is what made a path-sorted vector differ
	// from the contract's order.
	const (
		artifactName = "task.step.attempt.1.policy.json"
		anchorName   = "task.step.attempt.1.anchor.json"
	)
	if !(anchorName < artifactName) {
		t.Fatalf("the fixture must use an anchor name that sorts first, otherwise the path-order counterexample proves nothing")
	}

	t.Run("the-full-object-vector-is-dir-parent-artifact-anchor", func(t *testing.T) {
		root := t.TempDir()
		policy := filepath.Join(root, "agent-runner-policy")
		if err := os.MkdirAll(policy, 0o755); err != nil {
			t.Fatalf("create the policy directory: %v", err)
		}
		// The caller passes the anchor first, so neither the argument order nor the
		// path order may decide the slot.
		slots, err := managedAgentPolicyDomainNameSlots(policy, []string{anchorName, artifactName})
		if err != nil {
			t.Fatalf("a two-file domain must produce the four-slot vector, got %v", err)
		}
		p85R6RequireSlots(t, "the full object set", p85R6DescribeSlots(slots),
			"policy directory="+policy,
			"parent directory="+filepath.Dir(policy),
			"artifact="+artifactName,
			"anchor="+anchorName,
		)
		if slots[2].path != artifactName || slots[3].path != anchorName {
			t.Fatalf("the artifact must be slot 3 and the anchor slot 4, got %q and %q", slots[2].path, slots[3].path)
		}
		if slots[2].isDir || slots[3].isDir {
			t.Fatal("the two file slots must be file slots")
		}
	})

	t.Run("an-anchor-only-set-leaves-the-artifact-slot-empty", func(t *testing.T) {
		root := t.TempDir()
		policy := filepath.Join(root, "agent-runner-policy")
		if err := os.MkdirAll(policy, 0o755); err != nil {
			t.Fatalf("create the policy directory: %v", err)
		}
		slots, err := managedAgentPolicyDomainNameSlots(policy, []string{anchorName})
		if err != nil {
			t.Fatalf("an anchor-only domain must keep the four-slot shape, got %v", err)
		}
		p85R6RequireSlots(t, "the anchor-only set", p85R6DescribeSlots(slots),
			"policy directory="+policy,
			"parent directory="+filepath.Dir(policy),
			"artifact="+managedAgentPolicyDomainEmptySlot,
			"anchor="+anchorName,
		)
		if slots[2].present {
			t.Fatalf("the artifact slot must stay empty instead of receiving %q", slots[2].path)
		}
	})

	t.Run("the-handle-form-places-files-by-role-not-by-path-order", func(t *testing.T) {
		root := t.TempDir()
		policy := filepath.Join(root, "agent-runner-policy")
		if err := os.MkdirAll(policy, 0o755); err != nil {
			t.Fatalf("create the policy directory: %v", err)
		}
		artifactPath := filepath.Join(policy, artifactName)
		anchorPath := filepath.Join(policy, anchorName)
		for _, path := range []string{artifactPath, anchorPath} {
			if err := os.WriteFile(path, []byte("policy"), 0o600); err != nil {
				t.Fatalf("write %s: %v", path, err)
			}
		}

		directory, directoryID, err := evidence.OpenBoundDirVerified(policy)
		if err != nil {
			t.Skipf("this host's temporary root is not bindable as a verified directory: %v", err)
		}
		defer func() { _ = directory.Close() }()
		parent, parentID, err := evidence.OpenBoundDirVerified(filepath.Dir(policy))
		if err != nil {
			t.Skipf("this host's temporary root is not bindable as a verified directory: %v", err)
		}
		defer func() { _ = parent.Close() }()

		artifact, err := evidence.OpenBoundDirFile(directory, artifactName, 1<<20)
		if err != nil {
			t.Fatalf("open the artifact relative to the bound directory: %v", err)
		}
		defer func() { _ = artifact.Close() }()
		anchor, err := evidence.OpenBoundDirFile(directory, anchorName, 1<<20)
		if err != nil {
			t.Fatalf("open the anchor relative to the bound directory: %v", err)
		}
		defer func() { _ = anchor.Close() }()

		// The anchor is passed FIRST, exactly as a caller that walked the directory
		// in name order would.
		files := []ManagedPolicyDomainBoundFile{
			{Path: anchorPath, Role: managedPolicyDomainAnchorRole, Bound: anchor},
			{Path: artifactPath, Role: managedPolicyDomainArtifactRole, Bound: artifact},
		}
		slots, err := managedAgentPolicyDomainHandleSlots(directory, directoryID, parent, parentID, files)
		if err != nil {
			t.Fatalf("the handle form must produce the four-slot vector, got %v", err)
		}
		if len(slots) != 4 {
			t.Fatalf("the handle form must keep the fixed four-slot vector, got %d slots", len(slots))
		}
		for index, want := range []string{"policy directory", "parent directory", managedPolicyDomainArtifactRole, managedPolicyDomainAnchorRole} {
			if slots[index].label != want {
				t.Fatalf("slot %d must be labelled %q, got %q", index+1, want, slots[index].label)
			}
		}
		if slots[2].path != artifactPath {
			t.Fatalf("slot 3 must carry the artifact %q, got %q", artifactPath, slots[2].path)
		}
		if slots[3].path != anchorPath {
			t.Fatalf("slot 4 must carry the anchor %q, got %q", anchorPath, slots[3].path)
		}
		for index, want := range []evidence.FileID{directoryID, parentID} {
			if slots[index].dir == nil {
				t.Fatalf("slot %d must carry its bound directory handle", index+1)
			}
			identity, err := slots[index].dir.FileID()
			if err != nil {
				t.Fatalf("identify the handle of slot %d: %v", index+1, err)
			}
			if identity != want {
				t.Fatalf("slot %d must carry the recorded directory identity %s, got %s", index+1, want, identity)
			}
		}
		if slots[2].file != artifact || slots[3].file != anchor {
			t.Fatal("the artifact slot must carry the artifact handle and the anchor slot the anchor handle")
		}
	})

	t.Run("a-cross-role-claim-on-one-slot-is-refused", func(t *testing.T) {
		root := t.TempDir()
		policy := filepath.Join(root, "agent-runner-policy")
		if err := os.MkdirAll(policy, 0o755); err != nil {
			t.Fatalf("create the policy directory: %v", err)
		}
		declared := filepath.Join(policy, "declared-name-policy.json")
		fromName := filepath.Join(policy, "task.step.attempt.1.anchor.json")
		for _, path := range []string{declared, fromName} {
			if err := os.WriteFile(path, []byte("policy"), 0o600); err != nil {
				t.Fatalf("write %s: %v", path, err)
			}
		}

		directory, directoryID, err := evidence.OpenBoundDirVerified(policy)
		if err != nil {
			t.Skipf("this host's temporary root is not bindable as a verified directory: %v", err)
		}
		defer func() { _ = directory.Close() }()
		parent, parentID, err := evidence.OpenBoundDirVerified(filepath.Dir(policy))
		if err != nil {
			t.Skipf("this host's temporary root is not bindable as a verified directory: %v", err)
		}
		defer func() { _ = parent.Close() }()

		declaredBound, err := evidence.OpenBoundDirFile(directory, filepath.Base(declared), 1<<20)
		if err != nil {
			t.Fatalf("open the declared slot file: %v", err)
		}
		defer func() { _ = declaredBound.Close() }()
		nameBound, err := evidence.OpenBoundDirFile(directory, filepath.Base(fromName), 1<<20)
		if err != nil {
			t.Fatalf("open the name-marked slot file: %v", err)
		}
		defer func() { _ = nameBound.Close() }()

		// The first file DECLARES itself the anchor while its name says nothing; the
		// second says anchor by name. Both claim slot 4, so the pair cannot be folded.
		files := []ManagedPolicyDomainBoundFile{
			{Path: declared, Role: managedPolicyDomainAnchorRole, Bound: declaredBound},
			{Path: fromName, Role: "", Bound: nameBound},
		}
		_, err = managedAgentPolicyDomainHandleSlots(directory, directoryID, parent, parentID, files)
		if !errors.Is(err, ErrManagedAgentPolicyWriteProof) {
			t.Fatalf("two candidates claiming one slot must be refused, got %v", err)
		}
		if !strings.Contains(err.Error(), "claimed by both") {
			t.Fatalf("the refusal must be the slot-claim conflict, got %v", err)
		}
		if !strings.Contains(err.Error(), managedPolicyDomainAnchorRole) {
			t.Fatalf("the refusal must name the contested slot, got %v", err)
		}
	})

	t.Run("the-digest-tracks-the-slot-content", func(t *testing.T) {
		root := t.TempDir()
		policy := filepath.Join(root, "agent-runner-policy")
		if err := os.MkdirAll(policy, 0o755); err != nil {
			t.Fatalf("create the policy directory: %v", err)
		}
		artifactPath := filepath.Join(policy, artifactName)
		anchorPath := filepath.Join(policy, anchorName)
		for _, path := range []string{artifactPath, anchorPath} {
			if err := os.WriteFile(path, []byte("policy"), 0o600); err != nil {
				t.Fatalf("write %s: %v", path, err)
			}
		}
		full, err := managedAgentPolicyDomainNameSlots(policy, []string{artifactPath, anchorPath})
		if err != nil {
			t.Fatalf("observe the full domain: %v", err)
		}
		anchorOnly, err := managedAgentPolicyDomainNameSlots(policy, []string{anchorPath})
		if err != nil {
			t.Fatalf("observe the anchor-only domain: %v", err)
		}
		fullDigest, _ := managedAgentPolicyDomainWriteProofSlots(full)
		anchorOnlyDigest, _ := managedAgentPolicyDomainWriteProofSlots(anchorOnly)
		if fullDigest == "" || anchorOnlyDigest == "" {
			t.Skipf("this host cannot render the observation digest (the proof refuses before observing it): full=%q anchor-only=%q", fullDigest, anchorOnlyDigest)
		}
		if fullDigest == anchorOnlyDigest {
			t.Fatalf("an empty artifact slot and a filled one must not hash to the same digest, got %s twice", fullDigest)
		}
	})
}
