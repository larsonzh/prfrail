package console

// P85 F01 (design rev 3, section 8, rows M1-M3): the frozen-policy write domain is
// the policy directory, its PARENT and every existing artifact/anchor file, and the
// assembly option that used to let an out-of-package caller inject a well-formed
// proof is gone.
//
// The end-to-end counterexample the design names - "directory not writable while the
// artifact IS owned by the running user" - needs a directory this token cannot write
// and a file it can. Neither is constructible from an unprivileged test on either
// platform (a foreign owner needs chown/SeTakeOwnership), so the locally measurable
// form of the same fact is used instead: on this host every object of the domain is
// writable by the running identity, so the proof refuses and the refusal NAMES each
// observed object. Dropping an object from the observation therefore removes it from
// the refusal (M1: the artifact, M2: the parent) and changes the observation digest,
// which is exactly what the two named mutations have to break.

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// p85F01ArtifactAndAnchor returns the two names the reuse branch must observe for
// the one task attempt the T2 fixtures use.
func p85F01ArtifactAndAnchor(runtime *agentRunnerRuntime) (string, string) {
	return agentRunnerFrozenPolicyPath(runtime.runRoot, p83TaskID, p83StepID, 1),
		agentRunnerFrozenPolicyAnchorPath(runtime.runRoot, p83TaskID, p83StepID, 1)
}

// TestP85F01PolicyDomainObservesEachExistingFile is M1: the reuse branch must fold
// the existing artifact AND its anchor into the observation, before a single byte of
// either is read.
//
// Mutation checks:
//   - drop the `existing` argument from the reuse branch's
//     agentRunnerProvePolicyWriteProtection call => the refusal no longer names the
//     artifact or the anchor and "the-refusal-names-each-existing-file" goes red;
//   - drop the existing entries from managedAgentPolicyDomainObjects => the
//     observation digest stops depending on the artifact's identity and
//     "the-digest-tracks-the-artifact-identity" goes red.
func TestP85F01PolicyDomainObservesEachExistingFile(t *testing.T) {
	t.Run("the-refusal-names-each-existing-file", func(t *testing.T) {
		runtime, _ := p85FrozenReviewRuntime(t)
		artifactPath, anchorPath := p85F01ArtifactAndAnchor(runtime)
		body := p85PolicyBody(runtime.runID+"-chain", p83TaskID)

		// The production seam: no injected proof, so the real domain observation runs
		// on the reuse branch - which is the branch that must pass existing files.
		_, err := persistAgentRunnerFrozenPolicyProven(runtime.store, runtime.runRoot, runtime.runID, p83TaskID, p83StepID, 1, body, "")
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("a domain this identity can write must refuse the reuse branch, got %v", err)
		}
		message := err.Error()
		for _, name := range []string{artifactPath, anchorPath} {
			if !strings.Contains(message, fmt.Sprintf("%q", name)) {
				t.Fatalf("the refusal must name the observed existing file %s, got %v", name, err)
			}
		}
	})

	t.Run("the-digest-tracks-the-artifact-identity", func(t *testing.T) {
		root := t.TempDir()
		policy := filepath.Join(root, "agent-runner-policy")
		if err := os.MkdirAll(policy, 0o755); err != nil {
			t.Fatalf("create the policy directory: %v", err)
		}
		artifact := filepath.Join(policy, "artifact.json")
		if err := os.WriteFile(artifact, []byte("policy"), 0o600); err != nil {
			t.Fatalf("write the artifact: %v", err)
		}

		// Standing in for the replacement object, prepared OUTSIDE the observed
		// domain so the first observation is not disturbed by its presence.
		standIn := filepath.Join(t.TempDir(), "stand-in.json")
		if err := os.WriteFile(standIn, []byte("policy"), 0o600); err != nil {
			t.Fatalf("write the stand-in artifact: %v", err)
		}

		before, beforeErr := adapters.VerifyManagedAgentCannotWritePolicyDomain(policy, []string{artifact})
		if beforeErr == nil {
			t.Fatalf("the fixture must place the domain inside a directory this identity can write, so the proof refuses; got success")
		}
		if before == "" {
			t.Fatal("a refused observation must still carry the observation digest")
		}

		// Same path, same bytes, DIFFERENT object: only the per-file identity can
		// tell the two observations apart.
		//
		// Replace by RENAME, never by remove-then-recreate: measured on Linux/ext4, a
		// delete followed by a same-directory create inside a short window routinely
		// hands the freed inode back, so dev:ino - and therefore the digest - does not
		// move and this check reports a false red. A rename keeps the inode of the file
		// being renamed, so the new artifact carries the STAND-IN's inode while the
		// retired one still holds the old inode; the two identities cannot coincide.
		retired := artifact + ".old"
		if err := os.Rename(artifact, retired); err != nil {
			t.Fatalf("retire the artifact: %v", err)
		}
		if err := os.Rename(standIn, artifact); err != nil {
			t.Fatalf("promote the stand-in artifact: %v", err)
		}
		after, _ := adapters.VerifyManagedAgentCannotWritePolicyDomain(policy, []string{artifact})
		if after == before {
			t.Fatalf("replacing the observed artifact must change the observation digest, got %s twice", before)
		}
		if err := os.Remove(retired); err != nil {
			t.Fatalf("remove the retired artifact: %v", err)
		}
	})
}

// TestP85F01PolicyDomainObservesTheParentIdentity is M2: the parent of the policy
// directory gates rename/replace of the directory entry itself, so its identity is
// part of the observation.
//
// The layout keeps the policy directory's OWN identity and spelling unchanged and
// replaces only its parent object (rename the parent away, recreate it, move the
// policy directory back to the same path). A digest that does not include the
// parent therefore cannot tell the two layouts apart.
//
// Mutation checks:
//   - drop the parent object from managedAgentPolicyDomainObjects => "the-digest-
//     tracks-the-parent-identity" goes red (identical digests);
//   - the creation branch's refusal stops naming the parent => "the-creation-branch-
//     names-the-parent" goes red.
func TestP85F01PolicyDomainObservesTheParentIdentity(t *testing.T) {
	t.Run("the-creation-branch-names-the-parent", func(t *testing.T) {
		runtime := p83Runtime(t, chain.AuthorizationLedger{})
		p85AttachFrozenPolicyStore(t, runtime)
		body := p85PolicyBody(runtime.runID+"-chain", p83TaskID)
		_, err := persistAgentRunnerFrozenPolicyProven(runtime.store, runtime.runRoot, runtime.runID, p83TaskID, p83StepID, 1, body, "")
		if !errors.Is(err, ErrInvalidConfig) {
			t.Fatalf("the creation branch must refuse an unprovable write domain, got %v", err)
		}
		parent := filepath.Dir(agentRunnerPolicyDirectory(runtime.runRoot))
		if !strings.Contains(err.Error(), fmt.Sprintf("%q", parent)) {
			t.Fatalf("the creation branch must observe the parent directory %s, got %v", parent, err)
		}
	})

	t.Run("the-digest-tracks-the-parent-identity", func(t *testing.T) {
		root := t.TempDir()
		parent := filepath.Join(root, "run-root")
		if err := os.MkdirAll(parent, 0o755); err != nil {
			t.Fatalf("create the parent: %v", err)
		}
		policy := filepath.Join(parent, "agent-runner-policy")
		if err := os.MkdirAll(policy, 0o755); err != nil {
			t.Fatalf("create the policy directory: %v", err)
		}
		before, _ := adapters.VerifyManagedAgentCannotWritePolicyDomain(policy, nil)
		if before == "" {
			t.Fatal("a refused observation must still carry the observation digest")
		}
		policyIDBefore, err := evidence.FileIDAt(policy)
		if err != nil {
			t.Fatalf("identify the policy directory: %v", err)
		}
		parentIDBefore, err := evidence.FileIDAt(parent)
		if err != nil {
			t.Fatalf("identify the parent directory: %v", err)
		}

		// Replace the PARENT object while keeping the policy directory's inode and
		// final spelling: rename the parent away, recreate it, move the policy
		// directory back to the exact same path.
		retired := parent + "-retired"
		if err := os.Rename(parent, retired); err != nil {
			t.Fatalf("rename the parent away: %v", err)
		}
		if err := os.Mkdir(parent, 0o755); err != nil {
			t.Fatalf("recreate the parent: %v", err)
		}
		if err := os.Rename(filepath.Join(retired, "agent-runner-policy"), policy); err != nil {
			t.Fatalf("move the policy directory back to its own spelling: %v", err)
		}

		policyIDAfter, err := evidence.FileIDAt(policy)
		if err != nil {
			t.Fatalf("identify the moved policy directory: %v", err)
		}
		parentIDAfter, err := evidence.FileIDAt(parent)
		if err != nil {
			t.Fatalf("identify the recreated parent: %v", err)
		}
		if policyIDAfter != policyIDBefore {
			t.Fatalf("the fixture must keep the policy directory object: %s -> %s", policyIDBefore, policyIDAfter)
		}
		if parentIDAfter == parentIDBefore {
			t.Fatalf("the fixture must replace the parent object: both identities are %s", parentIDAfter)
		}

		after, _ := adapters.VerifyManagedAgentCannotWritePolicyDomain(policy, nil)
		if after == before {
			t.Fatalf("replacing the observed parent must change the observation digest, got %s twice", before)
		}
	})
}

// TestP85F01AssemblyOptionsCarryNoExportedWriteProof is M3: the bypass review r1
// named ("an exported option accepts any well-formed digest") is closed at COMPILE
// time, not by a runtime check.
//
// Mutation check: rename the field back to the exported ManagedWriteProof => the
// reflection assertion goes red (and any out-of-package use would then compile, which
// the source scan below also refuses).
func TestP85F01AssemblyOptionsCarryNoExportedWriteProof(t *testing.T) {
	options := reflect.TypeOf(AgentRunnerAssemblyOptions{})
	if field, found := options.FieldByName("ManagedWriteProof"); found {
		t.Fatalf("AgentRunnerAssemblyOptions must not expose a write-proof field (found %s, exported=%v)", field.Name, field.PkgPath == "")
	}
	field, found := options.FieldByName("managedWriteProof")
	if !found {
		t.Fatal("the in-package test seam managedWriteProof must still exist")
	}
	if field.PkgPath == "" {
		t.Fatal("the test seam must stay unexported, otherwise any out-of-package caller can inject a proof")
	}

	// Compile-surface scan: no production file outside this package may name the
	// identifier. An out-of-package bypass needs the name, so this is the executable
	// form of the design's `grep ManagedWriteProof` evidence.
	moduleRoot := filepath.Join("..", "..")
	violations := []string{}
	err := filepath.WalkDir(moduleRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "tmp", "build", "node_modules":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		if strings.HasPrefix(filepath.ToSlash(path), filepath.ToSlash(filepath.Join(moduleRoot, "internal", "console"))+"/") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(data), "ManagedWriteProof") {
			violations = append(violations, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("scan the module for out-of-package write-proof entries: %v", err)
	}
	if len(violations) > 0 {
		t.Fatalf("only the console package may know the write-proof identifier, found %v", violations)
	}
}
