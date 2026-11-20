package adapters

// P83 FIX ROUND 5 counterexamples (items T3c and T8b).
//
// T3c: the executable that was approved must still be the file the process is
// created from. T8b: the refusal evidence of a logging argument must not vary with
// the caller's value.

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

// p83Fix5EntryCount counts the directory entries, so a refusal can prove it left
// no footprint (no temporary directory).
func p83Fix5EntryCount(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	return len(entries)
}

// TestP83T3cExecutableReplacedAfterApprovalRefused is the T3c counterexample: the
// caller verifies the binding, the executable is then replaced, and Execute must
// refuse at the bound open/identity re-check without creating any footprint.
//
// Mutation check: deleting the bound-open hash comparison in Execute (or the
// ExpectedExecutableHash re-verification) lets the swapped image start and turns
// this red.
func TestP83T3cExecutableReplacedAfterApprovalRefused(t *testing.T) {
	dir := t.TempDir()
	executable := filepath.Join(dir, "copilot.exe")
	if err := os.WriteFile(executable, []byte("p83-fix5-approved-image"), 0o700); err != nil {
		t.Fatal(err)
	}
	approvedHash, err := CopilotCLIExecutableHash(executable)
	if err != nil {
		t.Fatal(err)
	}
	executor := ManagedCopilotCLIProbeExecutor{
		ExpectedExecutableHash: approvedHash,
		Timeout:                5 * time.Second,
		TemporaryDirectory:     dir,
	}
	// The approval happens here; the swap happens strictly after it.
	if err := executor.VerifyExecutableBinding(executable); err != nil {
		t.Fatalf("the approved executable must verify, got %v", err)
	}
	if err := os.WriteFile(executable, []byte("p83-fix5-swapped-image-with-a-different-size"), 0o700); err != nil {
		t.Fatal(err)
	}
	before := p83Fix5EntryCount(t, dir)

	execution := executor.Execute(context.Background(), CopilotCLIProbeInvocation{
		Executable:       executable,
		WorkingDirectory: dir,
		Args:             []string{"--version"},
	})
	if execution.Outcome != "start-failed" {
		t.Fatalf("a replaced executable must refuse before the spawn, got outcome=%q evidence=%v", execution.Outcome, execution.Evidence)
	}
	want := copilotCLIExecutorEvidence("executable-binding-mismatch")
	if len(execution.Evidence) != 1 || execution.Evidence[0] != want {
		t.Fatalf("the refusal must be the executable binding mismatch %q, got %v", want, execution.Evidence)
	}
	if after := p83Fix5EntryCount(t, dir); after != before {
		t.Fatalf("a refused executable must leave no temporary directory behind: before=%d after=%d", before, after)
	}
}

// TestP83T3cBoundIdentityRecheckRefusesReplacedFile pins the last line of defence
// on its own: neither the object identity nor the bytes read through that object
// still match the approved open, so the spawn is refused.
//
// Mutation check: deleting the identity or the digest comparison in
// requireSameExecutableFile returns nil and turns this red.
func TestP83T3cBoundIdentityRecheckRefusesReplacedFile(t *testing.T) {
	dir := t.TempDir()
	executable := filepath.Join(dir, "copilot.exe")
	if err := os.WriteFile(executable, []byte("p83-fix5-approved-image"), 0o700); err != nil {
		t.Fatal(err)
	}
	file, err := evidence.OpenBoundFile(executable, copilotCLIExecutableMaxBytes)
	if err != nil {
		t.Fatal(err)
	}
	fileID, err := file.FileID()
	if err != nil {
		t.Fatal(err)
	}
	approvedHash, err := file.Hash(copilotCLIExecutableDomain)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if err := requireSameExecutableFile(executable, fileID, approvedHash); err != nil {
		t.Fatalf("the unchanged file must pass the identity re-check, got %v", err)
	}
	if err := os.WriteFile(executable, []byte("p83-fix5-replacement-image-of-another-length"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := requireSameExecutableFile(executable, fileID, approvedHash); !errors.Is(err, ErrCopilotCLIExecutableBinding) {
		t.Fatalf("a replaced file must fail the identity re-check, got %v", err)
	}
	// Two distinct files with identical bytes must also be refused: the approved
	// identity names one object, and a copy - however byte-identical - is another.
	original := filepath.Join(dir, "p83-fix5-original-image")
	if err := os.WriteFile(original, []byte("p83-fix5-approved-image"), 0o700); err != nil {
		t.Fatal(err)
	}
	originalFile, err := evidence.OpenBoundFile(original, copilotCLIExecutableMaxBytes)
	if err != nil {
		t.Fatal(err)
	}
	originalID, err := originalFile.FileID()
	if err != nil {
		t.Fatal(err)
	}
	originalHash, err := originalFile.Hash(copilotCLIExecutableDomain)
	if err != nil {
		t.Fatal(err)
	}
	if err := originalFile.Close(); err != nil {
		t.Fatal(err)
	}
	copied := filepath.Join(dir, "p83-fix5-copied-image")
	if err := os.WriteFile(copied, []byte("p83-fix5-approved-image"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := requireSameExecutableFile(copied, originalID, originalHash); !errors.Is(err, ErrCopilotCLIExecutableBinding) {
		t.Fatalf("a byte-identical copy must fail the object identity re-check, got %v", err)
	}
}

// TestP83T8bRefusedLoggingEvidenceIsStable is the T8b counterexample: two
// different values of the same refused logging flag must produce IDENTICAL
// evidence, because the evidence classifies the flag and never echoes the value.
//
// Mutation check: replacing the flag class in the evidence with the caller's token
// makes the two evidence lists differ and turns this red.
func TestP83T8bRefusedLoggingEvidenceIsStable(t *testing.T) {
	dir := t.TempDir()
	executor := ManagedCopilotCLIProbeExecutor{Timeout: 5 * time.Second, TemporaryDirectory: dir}
	run := func(value string) CopilotCLIProbeExecution {
		t.Helper()
		return executor.Execute(context.Background(), CopilotCLIProbeInvocation{
			Executable:       filepath.Join(dir, "copilot.exe"),
			WorkingDirectory: dir,
			Args:             []string{"--log=" + value},
		})
	}
	first := run(filepath.Join(dir, "p83-fix5-probe-a.log"))
	second := run("/var/log/p83-fix5-probe-b-elsewhere.log")
	for _, execution := range []CopilotCLIProbeExecution{first, second} {
		if execution.Outcome != "start-failed" {
			t.Fatalf("a refused logging destination must not start, got outcome=%q", execution.Outcome)
		}
	}
	if !slices.Equal(first.Evidence, second.Evidence) {
		t.Fatalf("two refused values of the same flag must yield identical evidence: %v vs %v", first.Evidence, second.Evidence)
	}
	want := copilotCLIExecutorEvidence("logging-argument-refused:log")
	if len(first.Evidence) != 1 || first.Evidence[0] != want {
		t.Fatalf("the evidence must be the stable flag class %q, got %v", want, first.Evidence)
	}
	if strings.Contains(strings.Join(first.Evidence, " "), "p83-fix5-probe-a.log") {
		t.Fatalf("the evidence must not echo the refused value, got %v", first.Evidence)
	}
}
