package console

// P85 R6 (test round 6) / review r3 finding 2, console half: the M15/M17 legs
// measured the SHAPE decision through a stub that only returns a proof value, so
// nothing asserted which vector the production seam was actually asked to judge. The
// two anchors below therefore observe the seam's arguments and rebuild the four-slot
// vector from the chain's own handles:
//
//	creation => [policy directory | parent directory | ∅ | ∅]
//	reuse    => [policy directory | parent directory | artifact | anchor]
//
// Every entry of the rebuilt vector is compared with the identity the chain recorded
// when it bound the object, so a file slot that carries the wrong handle, or a
// role-less ordering that puts the anchor third, fails here rather than in a digest
// nobody can read.
//
// Mutation check: dropping the role labels, or handing the two file slots to the seam
// in path order, makes the rebuilt vector place the anchor third and turns this test
// red. The slot BUILDER itself (role-to-position folding) is pinned by the white-box
// sibling in internal/adapters.

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
)

// p85R6EmptyMarker is the canonical rendering of an absent slot (CONTRACTS.md:233).
const p85R6EmptyMarker = "∅"

// p85R6Observation is one judged vector plus the identities the chain recorded, so a
// leg can compare slot content with the binding instead of with itself.
type p85R6Observation struct {
	Vector     []string
	DirID      string
	ParentID   string
	ArtifactID string
	AnchorID   string
}

// p85R6Observe rebuilds the four-slot vector the observation seam received and
// refuses any vector that is not the contract's order.
func p85R6Observe(chain *agentRunnerPolicyDomainChain, files []adapters.ManagedPolicyDomainBoundFile) (p85R6Observation, error) {
	observation := p85R6Observation{ArtifactID: chain.artifactID, AnchorID: chain.anchorID}
	dirID, err := chain.dir.FileID()
	if err != nil {
		return observation, fmt.Errorf("identify the policy directory handle: %v", err)
	}
	if dirID != chain.dirID {
		return observation, fmt.Errorf("the policy directory handle is not the recorded object (%s != %s)", dirID, chain.dirID)
	}
	parentID, err := chain.parent.FileID()
	if err != nil {
		return observation, fmt.Errorf("identify the parent directory handle: %v", err)
	}
	if parentID != chain.parentID {
		return observation, fmt.Errorf("the parent directory handle is not the recorded object (%s != %s)", parentID, chain.parentID)
	}
	observation.DirID = string(dirID)
	observation.ParentID = string(parentID)
	observation.Vector = []string{
		"policy directory:" + string(dirID),
		"parent directory:" + string(parentID),
		"artifact:" + p85R6EmptyMarker,
		"anchor:" + p85R6EmptyMarker,
	}
	if len(files) > 2 {
		return observation, fmt.Errorf("the frozen policy domain carries at most two file slots, got %d", len(files))
	}
	for _, file := range files {
		identity, err := file.Bound.FileID()
		if err != nil {
			return observation, fmt.Errorf("identify the %s handle: %v", file.Role, err)
		}
		position := -1
		switch file.Role {
		case "artifact":
			position = 2
			if identity != chain.artifactID {
				return observation, fmt.Errorf("the artifact slot must carry the bound artifact handle (%s != %s)", identity, chain.artifactID)
			}
			if filepath.Clean(file.Path) != filepath.Clean(chain.artifactPath) {
				return observation, fmt.Errorf("the artifact slot must carry %q, got %q", chain.artifactPath, file.Path)
			}
		case "anchor":
			position = 3
			if identity != chain.anchorID {
				return observation, fmt.Errorf("the anchor slot must carry the bound anchor handle (%s != %s)", identity, chain.anchorID)
			}
			if filepath.Clean(file.Path) != filepath.Clean(chain.anchorPath) {
				return observation, fmt.Errorf("the anchor slot must carry %q, got %q", chain.anchorPath, file.Path)
			}
		default:
			return observation, fmt.Errorf("the seam must receive role-labelled file slots, got role %q", file.Role)
		}
		if !strings.HasSuffix(observation.Vector[position], p85R6EmptyMarker) {
			return observation, fmt.Errorf("slot %d is claimed twice (%s)", position+1, observation.Vector[position])
		}
		observation.Vector[position] = file.Role + ":" + identity
	}
	return observation, nil
}

// p85R6InstallObservationSeam records every vector the production reverification
// asks the platform observation to judge, and answers with the frozen proof.
func p85R6InstallObservationSeam(t *testing.T) *[]p85R6Observation {
	t.Helper()
	observations := &[]p85R6Observation{}
	previous := agentRunnerObserveManagedWriteProof
	agentRunnerObserveManagedWriteProof = func(chain *agentRunnerPolicyDomainChain, files []adapters.ManagedPolicyDomainBoundFile) (string, error) {
		observation, err := p85R6Observe(chain, files)
		if err != nil {
			t.Fatalf("the seam was handed a vector that is not the contract's four-slot order: %v", err)
		}
		*observations = append(*observations, observation)
		return assemblyManagedWriteProof, nil
	}
	t.Cleanup(func() { agentRunnerObserveManagedWriteProof = previous })
	return observations
}

// TestP85R6SlotOrderMatchesTheContract pins the real vector of both anchor shapes.
func TestP85R6SlotOrderMatchesTheContract(t *testing.T) {
	t.Run("a-creation-anchor-seals-dir-parent-and-two-empty-file-slots", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		observations := p85R6InstallObservationSeam(t)

		if _, err := p85R4LoadFrozen(t, runtime); err != nil {
			t.Fatalf("an unchanged creation anchor must be adopted, got %v", err)
		}
		// The current domain is observed first, the sealed domain second.
		if len(*observations) != 2 {
			t.Fatalf("the reverification must observe the current and the sealed domain, got %d observation(s)", len(*observations))
		}
		current, sealed := (*observations)[0], (*observations)[1]
		if len(sealed.Vector) != 4 {
			t.Fatalf("a creation anchor must seal four slots, got %d: %v", len(sealed.Vector), sealed.Vector)
		}
		if sealed.Vector[2] != "artifact:"+p85R6EmptyMarker || sealed.Vector[3] != "anchor:"+p85R6EmptyMarker {
			t.Fatalf("a creation anchor must seal [dir|parent|∅|∅], got %v", sealed.Vector)
		}
		if !strings.HasPrefix(current.Vector[2], "artifact:") || !strings.HasPrefix(current.Vector[3], "anchor:") {
			t.Fatalf("the current domain must be judged as [dir|parent|artifact|anchor], got %v", current.Vector)
		}
		if sealed.DirID != current.DirID || sealed.ParentID != current.ParentID {
			t.Fatalf("both observations must carry the same directory objects (%v vs %v)", sealed.Vector[:2], current.Vector[:2])
		}
	})

	t.Run("a-reuse-anchor-seals-dir-parent-artifact-anchor", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		p85R4RewriteAnchorShape(t, runtime, agentRunnerPolicyWriteProofShapeReuse)
		observations := p85R6InstallObservationSeam(t)

		if _, err := p85R4LoadFrozen(t, runtime); err != nil {
			t.Fatalf("a reuse anchor whose two files are still bound must be adopted, got %v", err)
		}
		if len(*observations) < 2 {
			t.Fatalf("a reuse anchor must be observed twice (current and sealed), got %d", len(*observations))
		}
		for index, observation := range *observations {
			if observation.ArtifactID == "" || observation.AnchorID == "" {
				t.Fatalf("observation %d must bind both file slots", index+1)
			}
			if observation.Vector[2] != "artifact:"+observation.ArtifactID {
				t.Fatalf("observation %d must place the artifact in slot 3, got %v", index+1, observation.Vector)
			}
			if observation.Vector[3] != "anchor:"+observation.AnchorID {
				t.Fatalf("observation %d must place the anchor in slot 4, got %v", index+1, observation.Vector)
			}
		}
	})
}
