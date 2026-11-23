package adapters

// P85 T3 (finding F6) - the install source counterexamples at the adapter boundary.
//
// Two layers are pinned here:
//
//  1. the observation itself (VerifyAIInstallSource) refuses a source this user can
//     mutate, and the carrier it produces cannot be decoded into a weaker claim;
//  2. Execute re-observes the source immediately before the spawn, so a probe that
//     was admitted earlier still refuses to start while the credential is in hand.
//
// The refusal evidence must stay a stable classification, never a path or a set of
// permission bits (T8b discipline).

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// p85T3WriteSpawnMarkerScript writes an executable script that creates marker when
// it runs. A refused probe must leave marker absent, which is the observable form
// of "the child was never created".
func p85T3WriteSpawnMarkerScript(t *testing.T, dir, marker string) string {
	t.Helper()
	var name, body string
	if runtime.GOOS == "windows" {
		name = "p85-t3-probe.cmd"
		body = "@echo off\r\n: > \"" + marker + "\" echo ran\r\n"
	} else {
		name = "p85-t3-probe.sh"
		body = "#!/bin/sh\n: > '" + marker + "'\n"
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

// TestP85T3WritableInstallDirectoryIsNotProven is the "the source became writable"
// counterexample at the observation layer: a temporary directory grants this user
// write, delete and permission rights, so immutability cannot be proven there.
//
// Mutation check: the Windows ACL capability probe and the Unix mode/owner walk are
// the only guards that turn this red when removed.
func TestP85T3WritableInstallDirectoryIsNotProven(t *testing.T) {
	dir := t.TempDir()
	executable := filepath.Join(dir, "copilot")
	if err := os.WriteFile(executable, []byte("p85-t3-writable-probe"), 0o700); err != nil {
		t.Fatal(err)
	}
	facts, err := VerifyAIInstallSource(executable)
	if err == nil {
		t.Fatalf("a writable installation source must not be proven, got facts=%+v", facts)
	}
	if !errors.Is(err, ErrAIInstallSourceProof) {
		t.Fatalf("the refusal must carry ErrAIInstallSourceProof, got %v", err)
	}
}

// TestP85T3ExecuteRefusesUnprovableInstallSourceWithoutSpawning is the pre-spawn
// counterexample: the child never runs, the evidence is the stable classification,
// and the temporary root the executor created is cleaned up.
//
// Mutation check: deleting (or short-circuiting) the VerifyAIInstallSource call in
// Execute lets the script run and turns this red on the marker assertion.
func TestP85T3ExecuteRefusesUnprovableInstallSourceWithoutSpawning(t *testing.T) {
	dir := t.TempDir()
	marker := filepath.Join(dir, "spawned.marker")
	executable := p85T3WriteSpawnMarkerScript(t, dir, marker)

	executor := ManagedCopilotCLIProbeExecutor{TemporaryDirectory: dir, Timeout: 5 * time.Second}
	execution := executor.Execute(context.Background(), CopilotCLIProbeInvocation{
		Executable:       executable,
		WorkingDirectory: dir,
	})
	if execution.Outcome != "start-failed" {
		t.Fatalf("an unprovable install source must refuse the start, got outcome=%q evidence=%v", execution.Outcome, execution.Evidence)
	}
	want := copilotCLIExecutorEvidence("install-source-unproven")
	if len(execution.Evidence) != 1 || execution.Evidence[0] != want {
		t.Fatalf("the refusal must be the closed-set install-source classification %q, got %v", want, execution.Evidence)
	}
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("the probe must not be spawned, marker stat err=%v", err)
	}
}

// p85T3WindowsProof builds a well-formed carrier body for the platform named by the
// proof. It never touches the platform it names, so it can be built on any host.
func p85T3WindowsProof(executablePath string) AIInstallSourceProof {
	hash := p85T3Hash("p85-t3-executable")
	return NewAIInstallSourceProof(AIInstallSourceFacts{
		ExecutablePath: executablePath,
		ExecutableHash: hash,
		FileID:         "win32:1:2:3",
		Platform:       AIInstallSourcePlatformFor("windows", "amd64"),
		Immutability: AIInstallSourceImmutability{Windows: &AIInstallSourceWindowsFacts{
			DirACLDeniesWriteDelete: true,
			ACLEvidenceHash:         p85T3Hash("p85-t3-acl"),
		}},
	}, time.Date(2026, 9, 8, 10, 0, 0, 0, time.UTC))
}

func p85T3Hash(seed string) string {
	return AIInstallSourceProofDigest([]byte(seed))
}

// p85T3HostAbsoluteExecutable returns a host-absolute path under a fresh temporary
// directory.
//
// P85 F16: the carrier fixture used to be the literal `C:\Program Files\...`, and
// AIInstallSourceProof.Validate decides "absolute" with filepath.IsAbs - which is
// false on Linux for a Windows drive letter, so the SHARED positive case failed at
// encode on the Linux leg. The fixture must therefore come from the host's own temp
// root, which is absolute on every platform the shared test compiles for.
func p85T3HostAbsoluteExecutable(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join(t.TempDir(), name)
}

// TestP85T3CarrierDecodeRefusesAWeakenedOrMixedClaim pins the carrier as a claim
// that cannot be decoded into something weaker than it says: a dropped ACL verdict,
// a platform/fact mixture, an unknown platform and a malformed body are all refusals.
//
// Mutation check: deleting the corresponding branch of AIInstallSourceProof.Validate
// turns the matching case green and this red.
func TestP85T3CarrierDecodeRefusesAWeakenedOrMixedClaim(t *testing.T) {
	valid := p85T3WindowsProof(p85T3HostAbsoluteExecutable(t, "copilot.exe"))
	wire, err := EncodeAIInstallSourceProof(valid)
	if err != nil {
		t.Fatalf("a well-formed carrier must encode: %v", err)
	}
	decoded, err := DecodeAIInstallSourceProof(wire)
	if err != nil {
		t.Fatalf("a well-formed carrier must decode: %v", err)
	}
	if decoded.ExecutableHash != valid.ExecutableHash || decoded.FileID != valid.FileID ||
		decoded.Platform != valid.Platform || decoded.Immutability.Windows == nil ||
		!decoded.Immutability.Windows.DirACLDeniesWriteDelete {
		t.Fatalf("the decoded carrier lost a fact: %+v", decoded)
	}

	weakenedACL := valid
	weakenedACL.Immutability = AIInstallSourceImmutability{Windows: &AIInstallSourceWindowsFacts{
		DirACLDeniesWriteDelete: false,
		ACLEvidenceHash:         valid.Immutability.Windows.ACLEvidenceHash,
	}}
	mixedPlatform := valid
	mixedPlatform.Immutability = AIInstallSourceImmutability{Linux: &AIInstallSourceLinuxFacts{
		UnwritablePath:        true,
		OwnershipEvidenceHash: p85T3Hash("p85-t3-ownership"),
	}}
	unknownPlatform := valid
	unknownPlatform.Platform = AIInstallSourcePlatformFor("darwin", "arm64")
	relativePath := valid
	relativePath.ExecutablePath = "copilot.exe"

	for name, proof := range map[string]AIInstallSourceProof{
		"acl-does-not-deny-write": weakenedACL,
		"platform-fact-mismatch":  mixedPlatform,
		"unknown-platform":        unknownPlatform,
		"relative-path":           relativePath,
	} {
		body, encodeErr := EncodeAIInstallSourceProof(proof)
		if encodeErr != nil {
			t.Fatalf("%s: the body must still encode so the DECODER is what refuses: %v", name, encodeErr)
		}
		if _, decodeErr := DecodeAIInstallSourceProof(body); decodeErr == nil {
			t.Fatalf("%s: a weakened or mixed claim must be refused at decode", name)
		}
	}
	if _, err := DecodeAIInstallSourceProof([]byte(`{"schemaVersion":1`)); err == nil {
		t.Fatal("a malformed body must be refused at decode")
	}
}

// TestP85T3MatchesDetectsAChangedObservation pins the comparison the command uses:
// any change of digest, identity, platform or permission fact makes the memo stop
// describing the live source.
func TestP85T3MatchesDetectsAChangedObservation(t *testing.T) {
	proof := p85T3WindowsProof(p85T3HostAbsoluteExecutable(t, "copilot.exe"))
	facts := AIInstallSourceFacts{
		ExecutablePath: proof.ExecutablePath,
		ExecutableHash: proof.ExecutableHash,
		FileID:         proof.FileID,
		Platform:       proof.Platform,
		Immutability:   proof.Immutability,
	}
	if !facts.Matches(proof) {
		t.Fatal("an identical observation must match the carrier")
	}
	changedHash := facts
	changedHash.ExecutableHash = p85T3Hash("p85-t3-swapped-bytes")
	changedIdentity := facts
	changedIdentity.FileID = "win32:1:2:4"
	changedPlatform := facts
	changedPlatform.Platform = AIInstallSourcePlatformFor("windows", "arm64")
	changedACL := facts
	changedACL.Immutability = AIInstallSourceImmutability{Windows: &AIInstallSourceWindowsFacts{
		DirACLDeniesWriteDelete: true,
		ACLEvidenceHash:         p85T3Hash("p85-t3-acl-after-a-loosening"),
	}}
	changedPath := facts
	changedPath.ExecutablePath = p85T3HostAbsoluteExecutable(t, "other.exe")

	for name, observation := range map[string]AIInstallSourceFacts{
		"changed-digest":   changedHash,
		"changed-identity": changedIdentity,
		"changed-platform": changedPlatform,
		"changed-acl":      changedACL,
		"changed-path":     changedPath,
	} {
		if observation.Matches(proof) {
			t.Fatalf("%s: a changed observation must not match the carrier", name)
		}
	}
}
