package console

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

const managedProcessIdentityFileName = "managed-process.identity.json"

const (
	defaultStopGrace   = 2 * time.Second
	defaultStopTimeout = 12 * time.Second
)

type processIdentityStopper func(context.Context, guard.ProcessIdentity, time.Duration) (guard.TerminationEvidence, error)

func newDefaultStopRun(getwd func() (string, error)) func(context.Context, string, string) ([]string, error) {
	return newDefaultStopRunWithStopper(getwd, processIdentityStopper(guard.StopProcessIdentity))
}

func newDefaultStopRunWithStopper(getwd func() (string, error), stopIdentity processIdentityStopper) func(context.Context, string, string) ([]string, error) {
	if getwd == nil {
		getwd = os.Getwd
	}
	if stopIdentity == nil {
		stopIdentity = processIdentityStopper(guard.StopProcessIdentity)
	}
	return func(ctx context.Context, runID, runDirHint string) ([]string, error) {
		if !evidence.ValidID(runID) {
			return nil, fmt.Errorf("invalid run id %q", runID)
		}
		cwd, err := getwd()
		if err != nil {
			return []string{digest("proofrail:approvals-stop-cwd-error:1\n", runID)}, err
		}
		runDir, err := resolveApprovalsRunDir(cwd, runID, runDirHint)
		if err != nil {
			return []string{digest("proofrail:approvals-stop-run-dir-error:1\n", runID)}, err
		}
		if summary, err := LoadRunSummary(ctx, runDir); err == nil {
			if summary.ChainState == "COMPLETED" || summary.ChainState == "FAILED" || summary.ChainState == "CANCELLED" {
				return []string{digest("proofrail:approvals-stop-terminal:1\n", runID, summary.ChainState)}, nil
			}
		}
		identityPath := filepath.Join(runDir, managedProcessIdentityFileName)
		identity, err := loadManagedProcessIdentity(identityPath)
		if err != nil {
			return []string{digest("proofrail:approvals-stop-missing-identity:1\n", runID, identityPath)}, fmt.Errorf("managed process identity unavailable: %w", err)
		}
		stopCtx, cancel := context.WithTimeout(ctx, defaultStopTimeout)
		defer cancel()
		proof, stopErr := stopIdentity(stopCtx, identity, defaultStopGrace)
		evidenceHashes := []string{}
		if evidence.ValidHash(proof.Hash) {
			evidenceHashes = append(evidenceHashes, proof.Hash)
		} else {
			evidenceHashes = append(evidenceHashes, digest("proofrail:approvals-stop-no-proof:1\n", runID))
		}
		if stopErr != nil {
			return evidenceHashes, stopErr
		}
		return evidenceHashes, nil
	}
}

func resolveApprovalsRunDir(cwd, runID, runDirHint string) (string, error) {
	hint := strings.TrimSpace(runDirHint)
	if hint == "" {
		return filepath.Abs(filepath.Join(cwd, "tmp", "prfrail-runs", runID))
	}
	return resolvePathFromCWD(cwd, hint)
}

func loadManagedProcessIdentity(path string) (guard.ProcessIdentity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return guard.ProcessIdentity{}, err
	}
	var identity guard.ProcessIdentity
	if err := evidence.DecodeStrictJSON(data, &identity); err != nil {
		return guard.ProcessIdentity{}, err
	}
	if identity.PID <= 0 || identity.StartToken == "" {
		return guard.ProcessIdentity{}, guard.ErrInvalidProcess
	}
	return identity, nil
}
