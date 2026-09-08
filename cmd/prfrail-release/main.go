package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/larsonzh/prfrail/internal/release"
)

func main() {
	if err := execute(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execute(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: prfrail-release <hash-tree|probe|selfhost|generate|verify>")
	}
	switch args[0] {
	case "hash-tree":
		set := flag.NewFlagSet("hash-tree", flag.ContinueOnError)
		root := set.String("root", "", "tree root")
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		hash, err := release.HashTree(*root)
		if err != nil {
			return err
		}
		fmt.Println(hash)
		return nil
	case "probe":
		set := flag.NewFlagSet("probe", flag.ContinueOnError)
		binary := set.String("binary", "", "candidate binary")
		noop := set.String("noop-chain", "", "noop chain fixture")
		executable := set.String("executable-chain", "", "executable chain fixture")
		workRoot := set.String("work-root", "", "isolated run root")
		output := set.String("output", "", "actual oracle output")
		expectedVersion := set.String("expected-version", "", "expected candidate version")
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		return release.ProbeCandidate(context.Background(), release.CandidateProbePlan{BinaryPath: *binary, NoopChainPath: *noop, ExecutableChainPath: *executable, WorkRoot: *workRoot, OutputPath: *output, ExpectedVersion: *expectedVersion})
	case "selfhost":
		set := flag.NewFlagSet("selfhost", flag.ContinueOnError)
		commit := set.String("commit", "", "candidate source commit")
		seed := set.String("seed", "", "seed root")
		seedHash := set.String("seed-hash", "", "expected seed tree hash")
		candidate := set.String("candidate", "", "candidate root")
		expected := set.String("expected", "", "external expected oracle")
		actual := set.String("actual", "", "candidate actual oracle")
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		result, err := release.VerifySelfHost(release.SelfHostPlan{SourceCommit: *commit, SeedRoot: *seed, ExpectedSeedHash: *seedHash, CandidateRoot: *candidate, ExpectedOraclePath: *expected, ActualOraclePath: *actual})
		return printResult(result, err)
	case "generate":
		set := flag.NewFlagSet("generate", flag.ContinueOnError)
		commit := set.String("commit", "", "candidate source commit")
		root := set.String("root", "", "artifact root")
		binary := set.String("binary", "", "binary path relative to artifact root")
		policy := set.String("license-policy", "", "license policy path")
		generatedAt := set.String("generated-at", "", "RFC3339 generation time")
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		timestamp, err := time.Parse(time.RFC3339, *generatedAt)
		if err != nil {
			return err
		}
		return release.GenerateBundleMetadata(release.GeneratePlan{SourceCommit: *commit, ArtifactRoot: *root, BinaryPath: *binary, LicensePolicyPath: *policy, GeneratedAt: timestamp})
	case "verify":
		set := flag.NewFlagSet("verify", flag.ContinueOnError)
		commit := set.String("commit", "", "candidate source commit")
		root := set.String("root", "", "artifact root")
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		result, err := release.VerifyReleaseBundle(release.BundlePlan{SourceCommit: *commit, ArtifactRoot: *root, ChecksumManifest: "SHA256SUMS", SBOMPath: "sbom.spdx.json", LicenseManifestPath: "licenses.json"})
		return printResult(result, err)
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func printResult(value any, err error) error {
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
