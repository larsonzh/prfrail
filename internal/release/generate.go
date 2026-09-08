package release

import (
	"crypto/sha256"
	"debug/buildinfo"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type GeneratePlan struct {
	SourceCommit      string
	ArtifactRoot      string
	BinaryPath        string
	LicensePolicyPath string
	GeneratedAt       time.Time
}

type licensePolicy struct {
	SchemaVersion string                `json:"schemaVersion"`
	Modules       []licensePolicyModule `json:"modules"`
}

type licensePolicyModule struct {
	Path    string `json:"path"`
	License string `json:"license"`
}

type spdxDocument struct {
	SPDXVersion       string        `json:"spdxVersion"`
	DataLicense       string        `json:"dataLicense"`
	SPDXID            string        `json:"SPDXID"`
	Name              string        `json:"name"`
	DocumentNamespace string        `json:"documentNamespace"`
	CreationInfo      creationInfo  `json:"creationInfo"`
	Packages          []spdxPackage `json:"packages"`
}

type creationInfo struct {
	Created  string   `json:"created"`
	Creators []string `json:"creators"`
}

type spdxPackage struct {
	Name             string `json:"name"`
	SPDXID           string `json:"SPDXID"`
	VersionInfo      string `json:"versionInfo"`
	DownloadLocation string `json:"downloadLocation"`
	FilesAnalyzed    bool   `json:"filesAnalyzed"`
	LicenseConcluded string `json:"licenseConcluded"`
	LicenseDeclared  string `json:"licenseDeclared"`
}

func GenerateBundleMetadata(plan GeneratePlan) error {
	if !commitPattern.MatchString(plan.SourceCommit) {
		return fmt.Errorf("%w: source commit must be a full lowercase SHA", ErrInvalidBundle)
	}
	created := plan.GeneratedAt.UTC()
	if created.IsZero() {
		return fmt.Errorf("%w: generatedAt is required", ErrInvalidBundle)
	}
	root, err := existingRealPath(plan.ArtifactRoot)
	if err != nil {
		return fmt.Errorf("%w: artifact root: %v", ErrInvalidBundle, err)
	}
	binaryPath, err := resolveBundlePath(root, plan.BinaryPath)
	if err != nil {
		return err
	}
	policyData, err := os.ReadFile(plan.LicensePolicyPath)
	if err != nil {
		return err
	}
	var policy licensePolicy
	if err := decodeStrict(policyData, &policy); err != nil {
		return fmt.Errorf("%w: invalid license policy: %v", ErrInvalidBundle, err)
	}
	if policy.SchemaVersion != "1.0.0" || len(policy.Modules) == 0 {
		return fmt.Errorf("%w: incomplete license policy", ErrInvalidBundle)
	}
	licenses := make(map[string]string, len(policy.Modules))
	for _, item := range policy.Modules {
		if item.Path == "" || item.License == "" {
			return fmt.Errorf("%w: incomplete license policy entry", ErrInvalidBundle)
		}
		if _, exists := licenses[item.Path]; exists {
			return fmt.Errorf("%w: duplicate license policy path %q", ErrInvalidBundle, item.Path)
		}
		licenses[item.Path] = item.License
	}

	info, err := buildinfo.ReadFile(binaryPath)
	if err != nil {
		return fmt.Errorf("%w: read Go build info: %v", ErrInvalidBundle, err)
	}
	modules := []licenseModule{{Path: info.Main.Path, Version: plan.SourceCommit, License: licenses[info.Main.Path]}}
	for _, dependency := range info.Deps {
		module := dependency
		if dependency.Replace != nil {
			module = dependency.Replace
		}
		modules = append(modules, licenseModule{Path: module.Path, Version: module.Version, License: licenses[module.Path]})
	}
	for _, module := range modules {
		if module.Path == "" || module.Version == "" || module.License == "" {
			return fmt.Errorf("%w: no approved license for module %q", ErrInvalidBundle, module.Path)
		}
	}
	sort.Slice(modules, func(left, right int) bool { return modules[left].Path < modules[right].Path })

	licenseOutput := licenseManifest{SchemaVersion: "1.0.0", Modules: modules}
	if err := writeJSON(filepath.Join(root, "licenses.json"), licenseOutput); err != nil {
		return err
	}
	packages := make([]spdxPackage, 0, len(modules))
	for index, module := range modules {
		packages = append(packages, spdxPackage{
			Name:             module.Path,
			SPDXID:           fmt.Sprintf("SPDXRef-Package-%d", index+1),
			VersionInfo:      module.Version,
			DownloadLocation: "NOASSERTION",
			FilesAnalyzed:    false,
			LicenseConcluded: module.License,
			LicenseDeclared:  module.License,
		})
	}
	sbom := spdxDocument{
		SPDXVersion:       "SPDX-2.3",
		DataLicense:       "CC0-1.0",
		SPDXID:            "SPDXRef-DOCUMENT",
		Name:              "prfrail-" + plan.SourceCommit,
		DocumentNamespace: "https://proofrail.dev/spdx/" + plan.SourceCommit,
		CreationInfo: creationInfo{
			Created:  created.Format("2006-01-02T15:04:05Z"),
			Creators: []string{"Tool: prfrail-release"},
		},
		Packages: packages,
	}
	if err := writeJSON(filepath.Join(root, "sbom.spdx.json"), sbom); err != nil {
		return err
	}
	return writeChecksumManifest(root, "SHA256SUMS")
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(path, append(data, '\n'))
}

func writeChecksumManifest(root, manifestName string) error {
	files, err := regularBundleFiles(root, manifestName)
	if err != nil {
		return err
	}
	var output strings.Builder
	for _, relative := range files {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			return err
		}
		fmt.Fprintf(&output, "%x  %s\n", sha256.Sum256(data), relative)
	}
	return writeAtomic(filepath.Join(root, manifestName), []byte(output.String()))
}
