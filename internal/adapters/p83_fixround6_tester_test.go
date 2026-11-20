package adapters

// P83 FIX ROUND 6 counterexamples (items T3a and T8b).
//
// T3a: the spawn-time proof requires BOTH the object identity and the digest of
// the approved executable, so neither a byte-identical copy (identity differs)
// nor an in-place rewrite (digest differs) may reach the spawn.
// T8b: the final start-failed evidence classifies the refusal and never echoes
// the refused value, so two different refused configurations of one class produce
// identical evidence.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// p83Fix6FileTruth returns the object identity and the executable digest of the
// file at path, taken from ONE bound open.
func p83Fix6FileTruth(t *testing.T, path string) (string, string) {
	t.Helper()
	file, err := evidence.OpenBoundFile(path, copilotCLIExecutableMaxBytes)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer func() { _ = file.Close() }()
	identity, err := file.FileID()
	if err != nil {
		t.Fatalf("identify %s: %v", path, err)
	}
	hash, err := file.Hash(copilotCLIExecutableDomain)
	if err != nil {
		t.Fatalf("digest %s: %v", path, err)
	}
	return identity, hash
}

// TestP83Fix6T3aIdentityAndDigestAreBothRequired is the T3a counterexample: the
// approved receiver names ONE object with ONE set of bytes, so the spawn-time
// check accepts the approved object and refuses both a same-content copy (the
// digest matches, the identity does not) and an in-place rewrite (the identity
// matches, the digest does not).
//
// Mutation checks. (1) Deleting the `actualFileID != approvedFileID` comparison
// admits the copy. (2) Deleting the `actualHash != approvedHash` comparison
// admits the rewritten object. Either turns this red.
func TestP83Fix6T3aIdentityAndDigestAreBothRequired(t *testing.T) {
	dir := t.TempDir()
	approved := filepath.Join(dir, "p83-fix6-approved-image.exe")
	if err := os.WriteFile(approved, []byte("p83-fix6-approved-image-bytes"), 0o700); err != nil {
		t.Fatal(err)
	}
	approvedID, approvedHash := p83Fix6FileTruth(t, approved)

	// Positive control: the approved object still passes, which is the launch.
	if err := requireSameExecutableFile(approved, approvedID, approvedHash); err != nil {
		t.Fatalf("the approved object must pass the spawn-time check, got %v", err)
	}

	// Identity differs, digest equal: a copy of the approved bytes elsewhere.
	copied := filepath.Join(dir, "p83-fix6-copied-image.exe")
	if err := os.WriteFile(copied, []byte("p83-fix6-approved-image-bytes"), 0o700); err != nil {
		t.Fatal(err)
	}
	copiedID, copiedHash := p83Fix6FileTruth(t, copied)
	if copiedHash != approvedHash {
		t.Fatalf("test precondition failed: the copy must digest identically (approved=%s copied=%s)", approvedHash, copiedHash)
	}
	if copiedID == approvedID {
		t.Fatalf("test precondition failed: two distinct files must not share an object identity %s", copiedID)
	}
	if err := requireSameExecutableFile(copied, approvedID, approvedHash); !errors.Is(err, ErrCopilotCLIExecutableBinding) {
		t.Fatalf("a byte-identical copy must be refused because its identity differs, got %v", err)
	}

	// Identity equal, digest differs: the approved object rewritten in place.
	file, err := os.OpenFile(approved, os.O_WRONLY|os.O_TRUNC, 0o700)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := file.WriteString("p83-fix6-rewritten-image-bytes"); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	rewrittenID, rewrittenHash := p83Fix6FileTruth(t, approved)
	if rewrittenID != approvedID {
		t.Fatalf("test precondition failed: an in-place rewrite must keep the object identity %s, got %s", approvedID, rewrittenID)
	}
	if rewrittenHash == approvedHash {
		t.Fatal("test precondition failed: the in-place rewrite must change the digest")
	}
	if err := requireSameExecutableFile(approved, approvedID, approvedHash); !errors.Is(err, ErrCopilotCLIExecutableBinding) {
		t.Fatalf("an in-place rewrite of the approved object must be refused because its digest differs, got %v", err)
	}
}

// TestP83Fix6T8bFinalRefusalEvidenceIsValueIndependent is the T8b counterexample:
// the start-failed evidence of a refusal is the same for two different refused
// configurations, because it names the stable refusal class and never the value
// that was refused.
//
// Mutation checks. (1) Using the refused token instead of the flag class in the
// logging refusal makes the first pair differ. (2) Adding the executable name to
// the path refusal makes the second pair differ. Either turns this red.
func TestP83Fix6T8bFinalRefusalEvidenceIsValueIndependent(t *testing.T) {
	dir := t.TempDir()
	approved := filepath.Join(dir, "copilot.exe")
	if err := os.WriteFile(approved, []byte("p83-fix6-t8b-approved-image"), 0o700); err != nil {
		t.Fatal(err)
	}
	approvedHash, err := CopilotCLIExecutableHash(approved)
	if err != nil {
		t.Fatal(err)
	}
	bound := ManagedCopilotCLIProbeExecutor{ExpectedExecutableHash: approvedHash, Timeout: 5 * time.Second, TemporaryDirectory: dir}
	run := func(executable string, args ...string) CopilotCLIProbeExecution {
		t.Helper()
		return bound.Execute(context.Background(), CopilotCLIProbeInvocation{
			Executable:       executable,
			WorkingDirectory: dir,
			Args:             args,
		})
	}

	// Pair 1: two refused logging destinations of DIFFERENT spelling and value.
	// The second one passes the flag and the value as separate tokens, which is a
	// different configuration of the same refusal.
	first := run(approved, "--LOG="+filepath.Join(dir, "p83-fix6-probe-a.log"))
	second := run(approved, "--log", "/var/log/p83-fix6-probe-b-elsewhere.log")
	for _, execution := range []CopilotCLIProbeExecution{first, second} {
		if execution.Outcome != "start-failed" {
			t.Fatalf("a refused logging destination must not start, got outcome=%q", execution.Outcome)
		}
	}
	if !slices.Equal(first.Evidence, second.Evidence) {
		t.Fatalf("two refused configurations of one class must yield identical evidence: %v vs %v", first.Evidence, second.Evidence)
	}
	wantLogging := copilotCLIExecutorEvidence("logging-argument-refused:log")
	if len(first.Evidence) != 1 || first.Evidence[0] != wantLogging {
		t.Fatalf("the logging refusal evidence must be the stable class %q, got %v", wantLogging, first.Evidence)
	}
	joined := strings.Join(first.Evidence, " ")
	for _, echoed := range []string{"p83-fix6-probe-a", "p83-fix6-probe-b", ".log", "/var/log"} {
		if strings.Contains(joined, echoed) {
			t.Fatalf("the evidence must not echo the refused value %q, got %v", echoed, first.Evidence)
		}
	}

	// Pair 2: two different bare command names, which are refused because a bare
	// name would be resolved through PATH at spawn time.
	third := run("copilot", "-p", "prompt")
	fourth := run("copilot-cli-probe", "-p", "prompt")
	for _, execution := range []CopilotCLIProbeExecution{third, fourth} {
		if execution.Outcome != "start-failed" {
			t.Fatalf("a bare command name must not start, got outcome=%q", execution.Outcome)
		}
	}
	if !slices.Equal(third.Evidence, fourth.Evidence) {
		t.Fatalf("two refused bare command names must yield identical evidence: %v vs %v", third.Evidence, fourth.Evidence)
	}
	wantPath := copilotCLIExecutorEvidence("executable-path-required")
	if len(third.Evidence) != 1 || third.Evidence[0] != wantPath {
		t.Fatalf("the path refusal evidence must be the stable class %q, got %v", wantPath, third.Evidence)
	}
	if strings.Contains(strings.Join(third.Evidence, " "), "copilot-cli-probe") {
		t.Fatalf("the evidence must not echo the refused command name, got %v", third.Evidence)
	}
}
