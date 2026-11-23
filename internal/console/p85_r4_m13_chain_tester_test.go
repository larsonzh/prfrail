package console

// P85 design rev 4, section 1 / test row M13: the EVIDENCE HANDLE CHAIN of one
// frozen-policy task attempt.
//
// The r2 review (High, agent_runner_policy.go:558) says the defect is "the object
// whose write protection was judged" versus "the object the bytes came from": every
// stage used to resolve the same names again. Design rev 4 replaced that with ONE
// binding - directory, parent, artifact, anchor - and a terminal
// ReassertIdentities that re-takes each identity from its handle AND requires each
// name to still address that very object.
//
// The counterexamples below are the cross-segment replacement and the two named
// mutations:
//
//   - a directory renamed away and replaced by a same-named stand-in BETWEEN the
//     observation segment and the adoption segment => the production load path
//     refuses on ReassertIdentities (mutation: deleting the name re-assertions
//     turns it green);
//   - an artifact replaced by a same-named file after the binding => the chain's
//     file slot refuses even though the handle still addresses the recorded object
//     (mutation: keeping only the directory checks - "只对目录不对制品" - turns it
//     green).
//
// The positive control below is the other half of the same fact: the bytes a
// snapshot read returns are the BOUND object's, not whatever the name addresses at
// read time, which is what makes the identity comparison meaningful at all.

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
)

// p85R4FrozenCreationRuntime freezes one creation-shaped anchor for the task
// attempt every P85 fixture uses and returns the runtime that reads it back.
//
// The freeze installs the constant observation stand-in (p85StubManagedWriteProofObservation)
// because on this host the assembly and the managed agent share one OS identity;
// each case below replaces that seam with the stub it needs, so the identity
// property under test is never measured through the shared-identity host.
func p85R4FrozenCreationRuntime(t *testing.T) *agentRunnerRuntime {
	t.Helper()
	runtime := p83Runtime(t, chain.AuthorizationLedger{})
	p85FreezeAgentRunnerPolicy(t, runtime, p83TaskID, p83StepID, 1)
	return runtime
}

// p85R4LoadFrozen runs the production load path of the fixture's one attempt.
func p85R4LoadFrozen(t *testing.T, runtime *agentRunnerRuntime) (agentRunnerFrozenPolicy, error) {
	t.Helper()
	return loadAgentRunnerFrozenPolicy(runtime.store, runtime.runRoot, p83RunID, p83TaskID, p83StepID, 1)
}

// p85R4ReplaceArtifact retires the artifact and puts a byte-identical stand-in at
// the same name, which is a DIFFERENT object behind the same path.
//
// Measured on this host: Windows refuses to rename a DIRECTORY while a file inside
// it is open (ERROR_ACCESS_DENIED), so the directory half of the cross-segment
// counterexample cannot be staged while the chain holds its file handles. The
// directory half therefore lives in p85_r4_linux_native_test.go, where renaming a
// bound directory is allowed and the comparison is a real dev:ino comparison; the
// artifact half below stages the same fact - a name that no longer addresses the
// bound object - on both platforms.
func p85R4ReplaceArtifact(t *testing.T, runtime *agentRunnerRuntime) error {
	t.Helper()
	path := agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1)
	payload, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := os.Rename(path, path+".retired"); err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0o600)
}

// TestP85R4M13TerminalReassertRefusesAStandInSwappedMidEntry is the cross-segment
// counterexample of M13: the replacement happens AFTER the write-domain observation
// (the seam runs inside agentRunnerReverifyManagedWriteProof, which is reached after
// both files have been read through their handles) and BEFORE the adoption. The
// original implementation - which compared nothing at the end - adopted the
// artifact; the chain must refuse.
//
// Mutation check: deleting the agentRunnerReassertNameIdentity calls from
// ReassertIdentities - or dropping the artifact/anchor loop ("only the directory,
// not the artifact") - makes the load adopt the artifact and turns this red.
func TestP85R4M13TerminalReassertRefusesAStandInSwappedMidEntry(t *testing.T) {
	runtime := p85R4FrozenCreationRuntime(t)

	previous := agentRunnerObserveManagedWriteProof
	swapped := false
	var swapErr error
	agentRunnerObserveManagedWriteProof = func(*agentRunnerPolicyDomainChain, []adapters.ManagedPolicyDomainBoundFile) (string, error) {
		if !swapped {
			swapped = true
			swapErr = p85R4ReplaceArtifact(t, runtime)
		}
		return assemblyManagedWriteProof, nil
	}
	t.Cleanup(func() { agentRunnerObserveManagedWriteProof = previous })

	_, err := p85R4LoadFrozen(t, runtime)
	if swapErr != nil {
		t.Fatalf("the fixture must be able to replace the artifact: %v", swapErr)
	}
	if !swapped {
		t.Fatal("the observation seam must have run, so the swap happened between the segments and not before them")
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("an artifact replaced between the segments must refuse the load, got %v", err)
	}
	if !strings.Contains(err.Error(), "frozen policy artifact") {
		t.Fatalf("the refusal must name the artifact slot whose identity moved, got %v", err)
	}
	if !strings.Contains(err.Error(), "no longer names the object that was bound") {
		t.Fatalf("the refusal must come from the terminal identity re-assertion, got %v", err)
	}

	// The stand-in is byte-identical, so only the OBJECT identity can tell the two
	// files apart - which is exactly the property the chain provides.
	path := agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1)
	standIn, statErr := os.Stat(path)
	if statErr != nil {
		t.Fatalf("the stand-in artifact must exist: %v", statErr)
	}
	retired, statErr := os.Stat(path + ".retired")
	if statErr != nil {
		t.Fatalf("the retired artifact must exist: %v", statErr)
	}
	if os.SameFile(standIn, retired) {
		t.Fatal("the fixture must produce a DIFFERENT object at the same name")
	}
}

// TestP85R4M13ChainSlotsRefuseAReplacedArtifact is the artifact half of M13: the
// directory is untouched, so only the file slot can notice. The bound handle keeps
// reading the recorded object (first assertion), while the name now addresses a new
// object (second assertion), which is precisely the "取证对象 != 读写对象" shape the
// finding names.
//
// Mutation check: dropping the artifact/anchor loop from ReassertIdentities ("only
// the directory, not the artifact") makes the final refusal disappear and turns
// this red.
func TestP85R4M13ChainSlotsRefuseAReplacedArtifact(t *testing.T) {
	runtime := p85R4FrozenCreationRuntime(t)
	bound, err := bindAgentRunnerPolicyDomain(runtime.runRoot, p83TaskID, p83StepID, 1)
	if err != nil {
		t.Fatalf("bind the policy domain chain: %v", err)
	}
	defer func() { _ = bound.Close() }()
	if bound.artifact == nil || bound.anchor == nil {
		t.Fatal("the fixture must bind both files")
	}
	recordedID := bound.artifactID
	original, err := bound.readBound(bound.artifact, bound.artifactPath, "frozen policy artifact")
	if err != nil {
		t.Fatalf("read the artifact through its handle: %v", err)
	}

	retired := bound.artifactPath + ".retired"
	if err := os.Rename(bound.artifactPath, retired); err != nil {
		t.Fatalf("retire the artifact: %v", err)
	}
	if err := os.WriteFile(bound.artifactPath, original, 0o600); err != nil {
		t.Fatalf("write the stand-in artifact: %v", err)
	}

	// The handle still addresses the object that was read: the chain's reads never
	// follow the name.
	again, err := bound.readBound(bound.artifact, bound.artifactPath, "frozen policy artifact")
	if err != nil {
		t.Fatalf("the bound handle must still read its object: %v", err)
	}
	if string(again) != string(original) {
		t.Fatalf("the bound handle must read the recorded object, got %q want %q", again, original)
	}
	identity, err := bound.artifact.FileID()
	if err != nil {
		t.Fatalf("re-identify the artifact on its handle: %v", err)
	}
	if identity != recordedID {
		t.Fatalf("the handle identity must not move, got %s want %s", identity, recordedID)
	}

	// The name, on the other hand, does not address it any more - and the terminal
	// re-assertion is what turns that into a refusal.
	err = bound.ReassertIdentities()
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("an artifact replaced after the binding must refuse, got %v", err)
	}
	if !strings.Contains(err.Error(), "frozen policy artifact") {
		t.Fatalf("the refusal must name the artifact slot, got %v", err)
	}
}

// TestP85R4M13BoundChainAcceptsAnUntouchedDomain is the positive control of the two
// legs above: without a replacement the same chain re-asserts cleanly, so the
// refusals are attributable to the replacement and not to the chain itself.
func TestP85R4M13BoundChainAcceptsAnUntouchedDomain(t *testing.T) {
	runtime := p85R4FrozenCreationRuntime(t)
	bound, err := bindAgentRunnerPolicyDomain(runtime.runRoot, p83TaskID, p83StepID, 1)
	if err != nil {
		t.Fatalf("bind the policy domain chain: %v", err)
	}
	defer func() { _ = bound.Close() }()
	if err := bound.ReassertIdentities(); err != nil {
		t.Fatalf("an untouched domain must re-assert cleanly, got %v", err)
	}
	if _, err := p85R4LoadFrozen(t, runtime); err != nil {
		t.Fatalf("an untouched creation anchor must load, got %v", err)
	}
}
