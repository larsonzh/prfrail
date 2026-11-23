package console

// P85 design rev 4, sections 2 and 5 / test rows M15 and M17: the explicit
// write-proof SHAPE of the assembly anchor, and the removal of the "a sealed short
// domain matched, so return" shortcut.
//
// Two independent questions were conflated before rev 4:
//
//	(a) is the domain still closed for every object that exists NOW? - the current
//	    write rights of the artifact and the anchor are always judged;
//	(b) is this anchor's proof the proof of the domain its SHAPE declares? - the
//	    sealed domain is rebuilt from the shape (creation: both file slots empty;
//	    reuse: artifact + anchor) and compared byte for byte.
//
// On this host the assembly and the managed agent share one OS identity, so the real
// domain observation can never SUCCEED and the positive legs would be unreachable.
// The stub below therefore replaces only the platform observation and lets a case
// state, per domain shape, what the observation says - which is what makes the two
// questions separable at all:
//
//	creationProof - the seal of a creation-shaped domain (no files bound);
//	fullProof     - the seal of a reuse-shaped domain / of the current object set;
//	failFullSet   - the current object set is writable NOW (the artifact/anchor are
//	                owned by this identity while the directory is not).
//
// Mutation contracts (design rev 4 section 6, rows M15/M17):
//   - deleting the current-domain observation from agentRunnerReverifyManagedWriteProof
//     turns M15's "writable-now" leg - and M17's creation leg - green;
//   - restoring the short-domain shortcut (return nil as soon as the sealed domain
//     matches) turns M17's creation leg green.

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85R4DomainObservationStub is the shape-aware observation stand-in. It never runs
// the platform guard; it answers per slot count so a case can separate "the sealed
// domain still matches" from "the domain is writable now".
type p85R4DomainObservationStub struct {
	creationProof     string
	fullProof         string
	failFullSet       bool
	failFirstFullCall bool
	fullCalls         int
	calls             int
}

// install replaces the observation seam for the duration of the test.
func (stub *p85R4DomainObservationStub) install(t *testing.T) {
	t.Helper()
	previous := agentRunnerObserveManagedWriteProof
	agentRunnerObserveManagedWriteProof = func(_ *agentRunnerPolicyDomainChain, files []adapters.ManagedPolicyDomainBoundFile) (string, error) {
		stub.calls++
		if len(files) == 0 {
			if stub.creationProof == "" {
				return "", fmt.Errorf("%w: stub refuses the creation domain", adapters.ErrManagedAgentPolicyWriteProof)
			}
			return stub.creationProof, nil
		}
		stub.fullCalls++
		if stub.failFullSet || (stub.failFirstFullCall && stub.fullCalls == 1) {
			return "", fmt.Errorf("%w: stub: an object of the current domain is writable", adapters.ErrManagedAgentPolicyWriteProof)
		}
		if stub.fullProof == "" {
			return "", fmt.Errorf("%w: stub refuses the full domain", adapters.ErrManagedAgentPolicyWriteProof)
		}
		return stub.fullProof, nil
	}
	t.Cleanup(func() { agentRunnerObserveManagedWriteProof = previous })
}

// p85R4AnchorPath is the assembly anchor of the fixture's one task attempt.
func p85R4AnchorPath(runtime *agentRunnerRuntime) string {
	return agentRunnerFrozenPolicyAnchorPath(runtime.runRoot, p83TaskID, p83StepID, 1)
}

// p85R4RewriteAnchorShape rewrites ONLY the declared write-proof shape of the
// anchor. The proof, the digests and the artifact stay byte-identical, so a
// refusal can only be about the shape.
func p85R4RewriteAnchorShape(t *testing.T, runtime *agentRunnerRuntime, shape string) {
	t.Helper()
	path := p85R4AnchorPath(runtime)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the frozen policy anchor: %v", err)
	}
	var anchor agentRunnerFrozenPolicyAnchor
	if err := evidence.DecodeStrictJSON(raw, &anchor); err != nil {
		t.Fatalf("decode the frozen policy anchor: %v", err)
	}
	anchor.ManagedWriteProofShape = shape
	payload, err := agentRunnerFrozenPolicyAnchorPayload(anchor)
	if err != nil {
		t.Fatalf("encode the rewritten anchor: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatalf("rewrite the frozen policy anchor: %v", err)
	}
}

// p85R4RemoveAnchorShapeField physically removes the shape key, which is the exact
// on-disk shape of an anchor written BEFORE the field existed. An empty string is a
// different fact (an unknown shape) and is covered by the rewrite leg above.
func p85R4RemoveAnchorShapeField(t *testing.T, runtime *agentRunnerRuntime) {
	t.Helper()
	path := p85R4AnchorPath(runtime)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the frozen policy anchor: %v", err)
	}
	var document map[string]any
	if err := json.Unmarshal(raw, &document); err != nil {
		t.Fatalf("decode the frozen policy anchor as a document: %v", err)
	}
	if _, present := document["managedWriteProofShape"]; !present {
		t.Fatal("the fixture anchor must carry the shape field before it can be removed")
	}
	delete(document, "managedWriteProofShape")
	rewritten, err := json.Marshal(document)
	if err != nil {
		t.Fatalf("re-encode the old-shaped anchor: %v", err)
	}
	if err := os.WriteFile(path, append(rewritten, '\n'), 0o600); err != nil {
		t.Fatalf("write the old-shaped anchor: %v", err)
	}
}

// TestP85R4M15ShapeDecidesWhichDomainIsSealed covers M15.
func TestP85R4M15ShapeDecidesWhichDomainIsSealed(t *testing.T) {
	t.Run("a-creation-anchor-reenters-unchanged", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		artifact, err := p85R4LoadFrozen(t, runtime)
		if err != nil {
			t.Fatalf("an unchanged creation anchor must be adopted on re-entry, got %v", err)
		}
		if artifact.PolicyHash == "" {
			t.Fatal("the adopted artifact must carry its policy hash")
		}
	})

	t.Run("a-reuse-anchor-reenters-unchanged", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		p85R4RewriteAnchorShape(t, runtime, agentRunnerPolicyWriteProofShapeReuse)
		// The constant stand-in reports the recorded proof for every domain, so this
		// leg measures the SHAPE branch (both files must still be bound) and not the
		// platform's verdict.
		p85StubManagedWriteProofObservation(t, assemblyManagedWriteProof)
		if _, err := p85R4LoadFrozen(t, runtime); err != nil {
			t.Fatalf("a reuse anchor whose two files are still present must be adopted, got %v", err)
		}
	})

	t.Run("a-creation-anchor-rewritten-as-reuse-is-refused", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		p85R4RewriteAnchorShape(t, runtime, agentRunnerPolicyWriteProofShapeReuse)
		stub := &p85R4DomainObservationStub{
			creationProof: assemblyManagedWriteProof,
			fullProof:     "sha256:" + strings.Repeat("b", 64),
		}
		stub.install(t)

		_, err := p85R4LoadFrozen(t, runtime)
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("an anchor that declares a domain it did not seal must be refused, got %v", err)
		}
		if !strings.Contains(err.Error(), "is not the proof of the domain") {
			t.Fatalf("the refusal must be the sealed-domain comparison, got %v", err)
		}
		// The creation seal was the recorded one, so the mismatch is provably the
		// shape's doing: the same stub accepts the unmodified anchor.
		runtime2 := p85R4FrozenCreationRuntime(t)
		stub2 := &p85R4DomainObservationStub{creationProof: assemblyManagedWriteProof, fullProof: "sha256:" + strings.Repeat("b", 64)}
		stub2.install(t)
		if _, err := p85R4LoadFrozen(t, runtime2); err != nil {
			t.Fatalf("the same stub must accept an unmodified creation anchor, got %v", err)
		}
	})

	t.Run("a-domain-that-is-writable-now-is-refused", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		stub := &p85R4DomainObservationStub{
			creationProof: assemblyManagedWriteProof,
			failFullSet:   true,
		}
		stub.install(t)

		_, err := p85R4LoadFrozen(t, runtime)
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("an artifact writable by this identity must refuse the load even when the sealed domain still matches, got %v", err)
		}
		if !strings.Contains(err.Error(), "is not closed to the managed agent identity") {
			t.Fatalf("the refusal must be the current-domain verification, got %v", err)
		}
		// Both questions were asked: the sealed domain was consulted after the
		// current one, so the check is not "sealed equals recorded, done".
		if stub.calls < 1 || stub.fullCalls < 1 {
			t.Fatalf("the current (full) domain must have been observed, got %d call(s) and %d full call(s)", stub.calls, stub.fullCalls)
		}
	})

	t.Run("an-anchor-without-a-shape-is-refused", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		p85R4RemoveAnchorShapeField(t, runtime)
		p85StubManagedWriteProofObservation(t, assemblyManagedWriteProof)

		_, err := p85R4LoadFrozen(t, runtime)
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("an anchor without a declared shape must be refused (fail-closed), got %v", err)
		}
		if !strings.Contains(err.Error(), "no recognized write-proof shape") {
			t.Fatalf("the refusal must name the missing shape, got %v", err)
		}
	})
}

// TestP85R4M17ShortDomainShortcutIsGone covers M17: the sealed-domain comparison
// may never stand in for the current-domain verification.
func TestP85R4M17ShortDomainShortcutIsGone(t *testing.T) {
	t.Run("a-closed-domain-with-a-creation-anchor-is-adopted", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		if _, err := p85R4LoadFrozen(t, runtime); err != nil {
			t.Fatalf("a creation anchor over a closed domain must be adopted, got %v", err)
		}
	})

	t.Run("a-creation-anchor-over-a-writable-object-set-is-refused", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		stub := &p85R4DomainObservationStub{
			creationProof: assemblyManagedWriteProof,
			failFullSet:   true,
		}
		stub.install(t)

		_, err := p85R4LoadFrozen(t, runtime)
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("a creation anchor must not short-circuit the current-domain verification, got %v", err)
		}
		if !strings.Contains(err.Error(), "is not closed to the managed agent identity") {
			t.Fatalf("the refusal must be the current-domain verification, got %v", err)
		}
	})

	t.Run("a-reuse-anchor-over-a-writable-artifact-is-refused", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		p85R4RewriteAnchorShape(t, runtime, agentRunnerPolicyWriteProofShapeReuse)
		// The FIRST full-domain observation is the current-object-set verification
		// (the sealed reuse domain is observed second), so this stub states exactly
		// "the artifact is writable now while the anchor's own seal still matches".
		stub := &p85R4DomainObservationStub{
			creationProof:     assemblyManagedWriteProof,
			fullProof:         assemblyManagedWriteProof,
			failFirstFullCall: true,
		}
		stub.install(t)

		_, err := p85R4LoadFrozen(t, runtime)
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("a reuse anchor must not short-circuit the current-domain verification, got %v", err)
		}
		if !strings.Contains(err.Error(), "is not closed to the managed agent identity") {
			t.Fatalf("the refusal must be the current-domain verification, got %v", err)
		}
	})

	t.Run("a-reuse-anchor-over-an-unwritable-artifact-is-adopted", func(t *testing.T) {
		runtime := p85R4FrozenCreationRuntime(t)
		p85R4RewriteAnchorShape(t, runtime, agentRunnerPolicyWriteProofShapeReuse)
		p85StubManagedWriteProofObservation(t, assemblyManagedWriteProof)
		if _, err := p85R4LoadFrozen(t, runtime); err != nil {
			t.Fatalf("a reuse anchor over a closed domain must be adopted, got %v", err)
		}
	})
}
