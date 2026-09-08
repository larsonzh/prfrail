package release

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrInvalidSelfHost = errors.New("invalid self-host evidence")
	ErrInvalidBundle   = errors.New("invalid release bundle")
	commitPattern      = regexp.MustCompile(`^[0-9a-f]{40}$`)
	checksumPattern    = regexp.MustCompile(`^([0-9a-f]{64})  ([^\\]+)$`)
)

type SelfHostPlan struct {
	SourceCommit       string
	SeedRoot           string
	ExpectedSeedHash   string
	CandidateRoot      string
	ExpectedOraclePath string
	ActualOraclePath   string
}

type SelfHostEvidence struct {
	SourceCommit  string
	SeedHash      string
	CandidateHash string
	OracleHash    string
}

type BundlePlan struct {
	SourceCommit        string
	ArtifactRoot        string
	ChecksumManifest    string
	SBOMPath            string
	LicenseManifestPath string
}

type BundleEvidence struct {
	SourceCommit                   string
	ChecksumManifestHash           string
	SBOMHash                       string
	LicenseManifestHash            string
	PublisherIdentityAuthenticated bool
}

type oracleDocument struct {
	SchemaVersion string       `json:"schemaVersion"`
	Cases         []oracleCase `json:"cases"`
}

type oracleCase struct {
	ID      string `json:"id"`
	Outcome string `json:"outcome"`
}

type licenseManifest struct {
	SchemaVersion string          `json:"schemaVersion"`
	Modules       []licenseModule `json:"modules"`
}

type licenseModule struct {
	Path    string `json:"path"`
	Version string `json:"version"`
	License string `json:"license"`
}

func VerifySelfHost(plan SelfHostPlan) (SelfHostEvidence, error) {
	if !commitPattern.MatchString(plan.SourceCommit) {
		return SelfHostEvidence{}, fmt.Errorf("%w: source commit must be a full lowercase SHA", ErrInvalidSelfHost)
	}
	seedRoot, err := existingRealPath(plan.SeedRoot)
	if err != nil {
		return SelfHostEvidence{}, fmt.Errorf("%w: seed root: %v", ErrInvalidSelfHost, err)
	}
	candidateRoot, err := existingRealPath(plan.CandidateRoot)
	if err != nil {
		return SelfHostEvidence{}, fmt.Errorf("%w: candidate root: %v", ErrInvalidSelfHost, err)
	}
	if pathsOverlap(seedRoot, candidateRoot) {
		return SelfHostEvidence{}, fmt.Errorf("%w: seed and candidate roots overlap", ErrInvalidSelfHost)
	}

	expectedOracle, err := existingRealPath(plan.ExpectedOraclePath)
	if err != nil {
		return SelfHostEvidence{}, fmt.Errorf("%w: expected oracle: %v", ErrInvalidSelfHost, err)
	}
	actualOracle, err := existingRealPath(plan.ActualOraclePath)
	if err != nil {
		return SelfHostEvidence{}, fmt.Errorf("%w: actual oracle: %v", ErrInvalidSelfHost, err)
	}
	if pathWithin(candidateRoot, expectedOracle) || !pathWithin(candidateRoot, actualOracle) {
		return SelfHostEvidence{}, fmt.Errorf("%w: oracle ownership is not independent", ErrInvalidSelfHost)
	}

	seedHash, err := HashTree(seedRoot)
	if err != nil {
		return SelfHostEvidence{}, fmt.Errorf("%w: hash seed: %v", ErrInvalidSelfHost, err)
	}
	if seedHash != plan.ExpectedSeedHash {
		return SelfHostEvidence{}, fmt.Errorf("%w: seed changed", ErrInvalidSelfHost)
	}
	expected, expectedHash, err := readOracle(expectedOracle)
	if err != nil {
		return SelfHostEvidence{}, fmt.Errorf("%w: expected oracle: %v", ErrInvalidSelfHost, err)
	}
	actual, _, err := readOracle(actualOracle)
	if err != nil {
		return SelfHostEvidence{}, fmt.Errorf("%w: actual oracle: %v", ErrInvalidSelfHost, err)
	}
	if !equalOracle(expected, actual) {
		return SelfHostEvidence{}, fmt.Errorf("%w: candidate result differs from external oracle", ErrInvalidSelfHost)
	}
	candidateHash, err := HashTree(candidateRoot)
	if err != nil {
		return SelfHostEvidence{}, fmt.Errorf("%w: hash candidate: %v", ErrInvalidSelfHost, err)
	}
	return SelfHostEvidence{
		SourceCommit:  plan.SourceCommit,
		SeedHash:      seedHash,
		CandidateHash: candidateHash,
		OracleHash:    expectedHash,
	}, nil
}

func VerifyReleaseBundle(plan BundlePlan) (BundleEvidence, error) {
	if !commitPattern.MatchString(plan.SourceCommit) {
		return BundleEvidence{}, fmt.Errorf("%w: source commit must be a full lowercase SHA", ErrInvalidBundle)
	}
	root, err := existingRealPath(plan.ArtifactRoot)
	if err != nil {
		return BundleEvidence{}, fmt.Errorf("%w: artifact root: %v", ErrInvalidBundle, err)
	}
	manifestPath, err := resolveBundlePath(root, plan.ChecksumManifest)
	if err != nil {
		return BundleEvidence{}, err
	}
	checksums, manifestHash, err := readChecksums(manifestPath)
	if err != nil {
		return BundleEvidence{}, err
	}
	files, err := regularBundleFiles(root, plan.ChecksumManifest)
	if err != nil {
		return BundleEvidence{}, err
	}
	if len(checksums) != len(files) {
		return BundleEvidence{}, fmt.Errorf("%w: checksum manifest does not cover every artifact", ErrInvalidBundle)
	}
	for _, path := range files {
		expected, ok := checksums[path]
		if !ok {
			return BundleEvidence{}, fmt.Errorf("%w: missing checksum for %q", ErrInvalidBundle, path)
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(path)))
		if err != nil {
			return BundleEvidence{}, err
		}
		actual := sha256.Sum256(data)
		if hex.EncodeToString(actual[:]) != expected {
			return BundleEvidence{}, fmt.Errorf("%w: checksum mismatch for %q", ErrInvalidBundle, path)
		}
	}

	sbomPath, err := resolveBundlePath(root, plan.SBOMPath)
	if err != nil {
		return BundleEvidence{}, err
	}
	sbomHash, sbomModules, err := validateSBOM(sbomPath, plan.SourceCommit)
	if err != nil {
		return BundleEvidence{}, err
	}
	licensePath, err := resolveBundlePath(root, plan.LicenseManifestPath)
	if err != nil {
		return BundleEvidence{}, err
	}
	licenseHash, licenseModules, err := validateLicenseManifest(licensePath)
	if err != nil {
		return BundleEvidence{}, err
	}
	if !equalLicenseModules(sbomModules, licenseModules) {
		return BundleEvidence{}, fmt.Errorf("%w: SBOM and license manifest modules differ", ErrInvalidBundle)
	}
	return BundleEvidence{
		SourceCommit:                   plan.SourceCommit,
		ChecksumManifestHash:           manifestHash,
		SBOMHash:                       sbomHash,
		LicenseManifestHash:            licenseHash,
		PublisherIdentityAuthenticated: false,
	}, nil
}

func HashTree(root string) (string, error) {
	inventory, err := evidence.BuildReleaseInventory(root, time.Time{})
	if err != nil {
		return "", err
	}
	return inventory.InventoryHash, nil
}

func existingRealPath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("empty path")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(abs)
}

func pathsOverlap(left, right string) bool {
	return pathWithin(left, right) || pathWithin(right, left)
}

func pathWithin(root, path string) bool {
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	return relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)))
}

func readOracle(path string) (oracleDocument, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return oracleDocument{}, "", err
	}
	var document oracleDocument
	if err := decodeStrict(data, &document); err != nil {
		return oracleDocument{}, "", err
	}
	if document.SchemaVersion != "1.0.0" || len(document.Cases) == 0 {
		return oracleDocument{}, "", errors.New("invalid oracle document")
	}
	seen := make(map[string]struct{}, len(document.Cases))
	for _, item := range document.Cases {
		if item.ID == "" || item.Outcome == "" {
			return oracleDocument{}, "", errors.New("oracle case fields must not be empty")
		}
		if _, exists := seen[item.ID]; exists {
			return oracleDocument{}, "", fmt.Errorf("duplicate oracle case %q", item.ID)
		}
		seen[item.ID] = struct{}{}
	}
	return document, evidence.Digest("", data), nil
}

func equalOracle(expected, actual oracleDocument) bool {
	if expected.SchemaVersion != actual.SchemaVersion || len(expected.Cases) != len(actual.Cases) {
		return false
	}
	for index := range expected.Cases {
		if expected.Cases[index] != actual.Cases[index] {
			return false
		}
	}
	return true
}

func resolveBundlePath(root, relative string) (string, error) {
	if relative == "" || filepath.IsAbs(relative) || filepath.Clean(relative) != relative || strings.Contains(relative, `\`) {
		return "", fmt.Errorf("%w: invalid relative path %q", ErrInvalidBundle, relative)
	}
	path := filepath.Join(root, relative)
	realPath, err := existingRealPath(path)
	if err != nil {
		return "", fmt.Errorf("%w: resolve %q: %v", ErrInvalidBundle, relative, err)
	}
	if !pathWithin(root, realPath) {
		return "", fmt.Errorf("%w: path escapes artifact root", ErrInvalidBundle)
	}
	return realPath, nil
}

func readChecksums(path string) (map[string]string, string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", err
	}
	checksums := make(map[string]string)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		matches := checksumPattern.FindStringSubmatch(scanner.Text())
		if matches == nil || !validManifestPath(matches[2]) {
			return nil, "", fmt.Errorf("%w: invalid SHA256SUMS line", ErrInvalidBundle)
		}
		if _, exists := checksums[matches[2]]; exists {
			return nil, "", fmt.Errorf("%w: duplicate checksum path %q", ErrInvalidBundle, matches[2])
		}
		checksums[matches[2]] = matches[1]
	}
	if err := scanner.Err(); err != nil {
		return nil, "", err
	}
	if len(checksums) == 0 {
		return nil, "", fmt.Errorf("%w: empty SHA256SUMS", ErrInvalidBundle)
	}
	return checksums, evidence.Digest("", data), nil
}

func validManifestPath(path string) bool {
	return path != "" && path == filepath.ToSlash(filepath.Clean(filepath.FromSlash(path))) && path != ".." && !strings.HasPrefix(path, "../") && !strings.HasPrefix(path, "/")
}

func regularBundleFiles(root, checksumManifest string) ([]string, error) {
	files := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("%w: symbolic links are not release artifacts", ErrInvalidBundle)
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("%w: unsupported artifact type", ErrInvalidBundle)
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if relative != checksumManifest {
			files = append(files, relative)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func validateSBOM(path, sourceCommit string) (string, []licenseModule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	var document spdxDocument
	if err := decodeStrict(data, &document); err != nil {
		return "", nil, fmt.Errorf("%w: invalid SBOM JSON: %v", ErrInvalidBundle, err)
	}
	if document.SPDXVersion != "SPDX-2.3" || document.DataLicense != "CC0-1.0" || document.SPDXID != "SPDXRef-DOCUMENT" {
		return "", nil, fmt.Errorf("%w: invalid SPDX document identity", ErrInvalidBundle)
	}
	if document.Name != "prfrail-"+sourceCommit || document.DocumentNamespace != "https://proofrail.dev/spdx/"+sourceCommit {
		return "", nil, fmt.Errorf("%w: SBOM is not bound to source commit", ErrInvalidBundle)
	}
	if document.CreationInfo.Created == "" || len(document.CreationInfo.Creators) != 1 || document.CreationInfo.Creators[0] != "Tool: prfrail-release" {
		return "", nil, fmt.Errorf("%w: incomplete SBOM creation info", ErrInvalidBundle)
	}
	modules := make([]licenseModule, 0, len(document.Packages))
	seen := make(map[string]struct{}, len(document.Packages))
	for _, item := range document.Packages {
		if item.Name == "" || item.SPDXID == "" || item.VersionInfo == "" || item.DownloadLocation == "" || item.FilesAnalyzed || item.LicenseDeclared == "" || item.LicenseDeclared != item.LicenseConcluded {
			return "", nil, fmt.Errorf("%w: incomplete SBOM package", ErrInvalidBundle)
		}
		key := item.Name + "@" + item.VersionInfo
		if _, exists := seen[key]; exists {
			return "", nil, fmt.Errorf("%w: duplicate SBOM package %q", ErrInvalidBundle, key)
		}
		seen[key] = struct{}{}
		modules = append(modules, licenseModule{Path: item.Name, Version: item.VersionInfo, License: item.LicenseDeclared})
	}
	if len(modules) == 0 {
		return "", nil, fmt.Errorf("%w: SBOM packages must not be empty", ErrInvalidBundle)
	}
	sort.Slice(modules, func(left, right int) bool { return modules[left].Path < modules[right].Path })
	return evidence.Digest("", data), modules, nil
}

func validateLicenseManifest(path string) (string, []licenseModule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, err
	}
	var manifest licenseManifest
	if err := decodeStrict(data, &manifest); err != nil {
		return "", nil, fmt.Errorf("%w: invalid license manifest: %v", ErrInvalidBundle, err)
	}
	if manifest.SchemaVersion != "1.0.0" || len(manifest.Modules) == 0 {
		return "", nil, fmt.Errorf("%w: incomplete license manifest", ErrInvalidBundle)
	}
	seen := make(map[string]struct{}, len(manifest.Modules))
	for _, module := range manifest.Modules {
		if module.Path == "" || module.Version == "" || module.License == "" {
			return "", nil, fmt.Errorf("%w: incomplete license entry", ErrInvalidBundle)
		}
		key := module.Path + "@" + module.Version
		if _, exists := seen[key]; exists {
			return "", nil, fmt.Errorf("%w: duplicate license entry %q", ErrInvalidBundle, key)
		}
		seen[key] = struct{}{}
	}
	sort.Slice(manifest.Modules, func(left, right int) bool { return manifest.Modules[left].Path < manifest.Modules[right].Path })
	return evidence.Digest("", data), manifest.Modules, nil
}

func equalLicenseModules(left, right []licenseModule) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func decodeStrict(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("trailing JSON value")
	}
	return nil
}
