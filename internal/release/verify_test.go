package release

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestVerifySelfHostAcceptsIsolatedMatchingCandidate(t *testing.T) {
	fixture := newSelfHostFixture(t)

	evidence, err := VerifySelfHost(SelfHostPlan{
		SourceCommit:       fixture.commit,
		SeedRoot:           fixture.seedRoot,
		ExpectedSeedHash:   fixture.seedHash,
		CandidateRoot:      fixture.candidateRoot,
		ExpectedOraclePath: fixture.expectedOracle,
		ActualOraclePath:   fixture.actualOracle,
	})
	if err != nil {
		t.Fatal(err)
	}
	if evidence.SourceCommit != fixture.commit || evidence.SeedHash != fixture.seedHash {
		t.Fatalf("unexpected self-host evidence: %+v", evidence)
	}
}

func TestVerifySelfHostRejectsOverlapMismatchAndSeedMutation(t *testing.T) {
	fixture := newSelfHostFixture(t)

	plan := SelfHostPlan{
		SourceCommit:       fixture.commit,
		SeedRoot:           fixture.seedRoot,
		ExpectedSeedHash:   fixture.seedHash,
		CandidateRoot:      filepath.Join(fixture.seedRoot, "candidate"),
		ExpectedOraclePath: fixture.expectedOracle,
		ActualOraclePath:   fixture.actualOracle,
	}
	if _, err := VerifySelfHost(plan); err == nil {
		t.Fatal("seed/candidate overlap must fail closed")
	}

	plan.CandidateRoot = fixture.candidateRoot
	writeFile(t, fixture.actualOracle, []byte(`{"schemaVersion":"1.0.0","cases":[{"id":"accept","outcome":"failed"}]}`))
	if _, err := VerifySelfHost(plan); err == nil {
		t.Fatal("candidate/oracle mismatch must fail closed")
	}

	writeFile(t, fixture.actualOracle, []byte(`{"schemaVersion":"1.0.0","cases":[{"id":"accept","outcome":"passed"}]}`))
	writeFile(t, filepath.Join(fixture.seedRoot, "seed.txt"), []byte("mutated"))
	if _, err := VerifySelfHost(plan); err == nil {
		t.Fatal("seed mutation must fail closed")
	}
}

func TestVerifyReleaseBundleChecksCommitChecksumsSBOMAndLicenses(t *testing.T) {
	root := t.TempDir()
	commit := "0123456789abcdef0123456789abcdef01234567"
	writeFile(t, filepath.Join(root, "prfrail.exe"), []byte("binary"))
	writeFile(t, filepath.Join(root, "sbom.spdx.json"), []byte(`{"spdxVersion":"SPDX-2.3","dataLicense":"CC0-1.0","SPDXID":"SPDXRef-DOCUMENT","name":"prfrail-`+commit+`","documentNamespace":"https://proofrail.dev/spdx/`+commit+`","creationInfo":{"created":"2026-09-08T00:00:00Z","creators":["Tool: prfrail-release"]},"packages":[{"name":"golang.org/x/text","SPDXID":"SPDXRef-Package-1","versionInfo":"v0.22.0","downloadLocation":"NOASSERTION","filesAnalyzed":false,"licenseConcluded":"BSD-3-Clause","licenseDeclared":"BSD-3-Clause"}]}`))
	writeFile(t, filepath.Join(root, "licenses.json"), []byte(`{"schemaVersion":"1.0.0","modules":[{"path":"golang.org/x/text","version":"v0.22.0","license":"BSD-3-Clause"}]}`))
	writeChecksums(t, root, "prfrail.exe", "sbom.spdx.json", "licenses.json")

	if _, err := VerifyReleaseBundle(BundlePlan{
		SourceCommit:        commit,
		ArtifactRoot:        root,
		ChecksumManifest:    "SHA256SUMS",
		SBOMPath:            "sbom.spdx.json",
		LicenseManifestPath: "licenses.json",
	}); err != nil {
		t.Fatal(err)
	}

	if _, err := VerifyReleaseBundle(BundlePlan{
		SourceCommit:        "main",
		ArtifactRoot:        root,
		ChecksumManifest:    "SHA256SUMS",
		SBOMPath:            "sbom.spdx.json",
		LicenseManifestPath: "licenses.json",
	}); err == nil {
		t.Fatal("symbolic source revisions must be rejected")
	}

	writeFile(t, filepath.Join(root, "prfrail.exe"), []byte("tampered"))
	if _, err := VerifyReleaseBundle(BundlePlan{
		SourceCommit:        commit,
		ArtifactRoot:        root,
		ChecksumManifest:    "SHA256SUMS",
		SBOMPath:            "sbom.spdx.json",
		LicenseManifestPath: "licenses.json",
	}); err == nil {
		t.Fatal("checksum mismatch must fail closed")
	}
}

func TestVerifyReleaseBundleRejectsSBOMLicenseMismatch(t *testing.T) {
	root := t.TempDir()
	commit := "0123456789abcdef0123456789abcdef01234567"
	writeFile(t, filepath.Join(root, "prfrail"), []byte("binary"))
	writeFile(t, filepath.Join(root, "sbom.spdx.json"), []byte(`{"spdxVersion":"SPDX-2.3","dataLicense":"CC0-1.0","SPDXID":"SPDXRef-DOCUMENT","name":"prfrail-`+commit+`","documentNamespace":"https://proofrail.dev/spdx/`+commit+`","creationInfo":{"created":"2026-09-08T00:00:00Z","creators":["Tool: prfrail-release"]},"packages":[{"name":"example.com/module","SPDXID":"SPDXRef-Package-1","versionInfo":"v1.0.0","downloadLocation":"NOASSERTION","filesAnalyzed":false,"licenseConcluded":"MIT","licenseDeclared":"MIT"}]}`))
	writeFile(t, filepath.Join(root, "licenses.json"), []byte(`{"schemaVersion":"1.0.0","modules":[{"path":"example.com/module","version":"v1.0.0","license":"Apache-2.0"}]}`))
	writeChecksums(t, root, "prfrail", "sbom.spdx.json", "licenses.json")
	if _, err := VerifyReleaseBundle(BundlePlan{
		SourceCommit:        commit,
		ArtifactRoot:        root,
		ChecksumManifest:    "SHA256SUMS",
		SBOMPath:            "sbom.spdx.json",
		LicenseManifestPath: "licenses.json",
	}); err == nil {
		t.Fatal("SBOM/license mismatch must fail closed")
	}
}

func TestGenerateBundleMetadataRequiresApprovedLicenses(t *testing.T) {
	root := t.TempDir()
	binaryPath := filepath.Join(root, "release.test")
	copyFile(t, binaryPath, executablePath(t))
	policyPath := filepath.Join(root, "policy.json")
	writeFile(t, policyPath, []byte(`{"schemaVersion":"1.0.0","modules":[{"path":"github.com/larsonzh/prfrail","license":"MIT"},{"path":"golang.org/x/text","license":"BSD-3-Clause"},{"path":"gopkg.in/yaml.v3","license":"MIT"}]}`))
	commit := "0123456789abcdef0123456789abcdef01234567"

	err := GenerateBundleMetadata(GeneratePlan{
		SourceCommit:      commit,
		ArtifactRoot:      root,
		BinaryPath:        "release.test",
		LicensePolicyPath: policyPath,
		GeneratedAt:       time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := VerifyReleaseBundle(BundlePlan{
		SourceCommit:        commit,
		ArtifactRoot:        root,
		ChecksumManifest:    "SHA256SUMS",
		SBOMPath:            "sbom.spdx.json",
		LicenseManifestPath: "licenses.json",
	}); err != nil {
		t.Fatal(err)
	}

	writeFile(t, policyPath, []byte(`{"schemaVersion":"1.0.0","modules":[{"path":"golang.org/x/text","license":"BSD-3-Clause"}]}`))
	if err := GenerateBundleMetadata(GeneratePlan{
		SourceCommit:      commit,
		ArtifactRoot:      root,
		BinaryPath:        "release.test",
		LicensePolicyPath: policyPath,
		GeneratedAt:       time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC),
	}); err == nil {
		t.Fatal("missing dependency license must fail closed")
	}
}

func TestGenerateBundleMetadataRejectsMissingTimeBeforeWriting(t *testing.T) {
	root := t.TempDir()
	if err := GenerateBundleMetadata(GeneratePlan{
		SourceCommit: "0123456789abcdef0123456789abcdef01234567",
		ArtifactRoot: root,
	}); err == nil {
		t.Fatal("missing generatedAt must fail closed")
	}
	if _, err := os.Stat(filepath.Join(root, "licenses.json")); !os.IsNotExist(err) {
		t.Fatalf("metadata must not be written on invalid input: %v", err)
	}
}

type selfHostFixture struct {
	commit         string
	seedRoot       string
	seedHash       string
	candidateRoot  string
	expectedOracle string
	actualOracle   string
}

func newSelfHostFixture(t *testing.T) selfHostFixture {
	t.Helper()
	root := t.TempDir()
	seedRoot := filepath.Join(root, "seed")
	candidateRoot := filepath.Join(root, "candidate")
	expectedOracle := filepath.Join(root, "oracle", "expected.json")
	actualOracle := filepath.Join(candidateRoot, "actual.json")
	writeFile(t, filepath.Join(seedRoot, "seed.txt"), []byte("seed"))
	writeFile(t, expectedOracle, []byte(`{"schemaVersion":"1.0.0","cases":[{"id":"accept","outcome":"passed"}]}`))
	writeFile(t, actualOracle, []byte(`{"schemaVersion":"1.0.0","cases":[{"id":"accept","outcome":"passed"}]}`))
	seedHash, err := HashTree(seedRoot)
	if err != nil {
		t.Fatal(err)
	}
	return selfHostFixture{
		commit:         "0123456789abcdef0123456789abcdef01234567",
		seedRoot:       seedRoot,
		seedHash:       seedHash,
		candidateRoot:  candidateRoot,
		expectedOracle: expectedOracle,
		actualOracle:   actualOracle,
	}
}

func writeChecksums(t *testing.T, root string, paths ...string) {
	t.Helper()
	manifest := ""
	for _, path := range paths {
		data, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(data)
		manifest += fmt.Sprintf("%s  %s\n", hex.EncodeToString(sum[:]), path)
	}
	writeFile(t, filepath.Join(root, "SHA256SUMS"), []byte(manifest))
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}

func executablePath(t *testing.T) string {
	t.Helper()
	path, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func copyFile(t *testing.T, destination, source string) {
	t.Helper()
	input, err := os.Open(source)
	if err != nil {
		t.Fatal(err)
	}
	defer input.Close()
	output, err := os.Create(destination)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(output, input); err != nil {
		output.Close()
		t.Fatal(err)
	}
	if err := output.Close(); err != nil {
		t.Fatal(err)
	}
}
