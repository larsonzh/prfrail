package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/gates"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// postflight is the real chain-owned postflight port of this harness. The port
// contract (internal/chain/postflight.go) requires it to re-derive the facts
// itself and reconcile them field by field, without writing chain state:
//
//	leg 1 manifest   re-capture the workspace and require the entry set to equal the
//	                 frozen post manifest, exactly (added/removed/changed reported)
//	leg 2 diff       recompute the diff from the recorded pre manifest and the
//	                 fresh capture and require the recorded diff artifact
//	leg 3 logs       re-read both log streams and require the frozen log digest
//	leg 4 usage      re-read usage.json and require the frozen usage digest
//	leg 5 secrets    scan the workspace with the production gates secret scanner
//	leg 6 side effect require every pid the CLI reported to be gone
//
// A leg that cannot be proven makes the verdict uncertain; a leg that proves the
// facts wrong makes it failed. Both are non-passing and are never silent.
type postflight struct {
	requestID    string
	workspace    string
	snapshotRoot string
	owner        *ports
}

func newPostflight(requestID, workspace, snapshotRoot string, owner *ports) *postflight {
	return &postflight{requestID: requestID, workspace: workspace, snapshotRoot: snapshotRoot, owner: owner}
}

// RunPostflight implements chain.PostflightPort.
func (p *postflight) RunPostflight(ctx context.Context, request chain.PostflightRequest) (chain.PostflightDecision, error) {
	evidenceDir, err := p.owner.replay.RunEvidenceDir(p.requestID)
	if err != nil {
		return chain.PostflightDecision{}, err
	}
	storedPre, err := os.ReadFile(filepath.Join(evidenceDir, "manifest.pre.json"))
	if err != nil {
		return p.uncertain("postflight cannot read the recorded pre manifest: " + err.Error())
	}
	storedPost, err := os.ReadFile(filepath.Join(evidenceDir, "manifest.post.json"))
	if err != nil {
		return p.uncertain("postflight cannot read the recorded post manifest: " + err.Error())
	}
	storedDiff, err := os.ReadFile(filepath.Join(evidenceDir, "diff.json"))
	if err != nil {
		return p.uncertain("postflight cannot read the recorded diff: " + err.Error())
	}

	fresh, err := p.owner.CaptureWorkspaceManifest(ctx, p.workspace)
	if err != nil {
		return p.uncertain("postflight cannot re-capture the workspace: " + err.Error())
	}
	storedEntries, err := manifestEntries(storedPost)
	if err != nil {
		return p.uncertain("postflight cannot decode the recorded post manifest: " + err.Error())
	}
	freshEntries, err := manifestEntries(fresh)
	if err != nil {
		return p.uncertain("postflight cannot decode the fresh manifest: " + err.Error())
	}
	if changed := changedPaths(diffManifests(storedEntries, freshEntries)); len(changed) > 0 {
		return p.failed("proofrail:b4-postflight-manifest-mismatch:1\n", strings.Join(changed, "\n"))
	}

	preEntries, err := manifestEntries(storedPre)
	if err != nil {
		return p.uncertain("postflight cannot decode the recorded pre manifest: " + err.Error())
	}
	freshDiff, err := marshalDiff(diffManifests(preEntries, freshEntries))
	if err != nil {
		return p.uncertain("postflight cannot recompute the diff: " + err.Error())
	}
	if string(freshDiff) != string(storedDiff) {
		return p.failed("proofrail:b4-postflight-diff-mismatch:1\n", digestOf("", string(freshDiff)))
	}

	logsHash, err := hashLogs(evidenceDir)
	if err != nil {
		return p.uncertain("postflight cannot read the log artifacts: " + err.Error())
	}
	if logsHash != request.Facts.LogHash {
		return p.failed("proofrail:b4-postflight-log-mismatch:1\n", logsHash)
	}
	usageHash, err := fileHash(filepath.Join(evidenceDir, "usage.json"))
	if err != nil {
		return p.uncertain("postflight cannot read the usage artifact: " + err.Error())
	}
	if usageHash != request.Facts.UsageHash {
		return p.failed("proofrail:b4-postflight-usage-mismatch:1\n", usageHash)
	}

	if err := p.scanWorkspace(ctx); err != nil {
		return p.failed("proofrail:b4-postflight-secret:1\n", err.Error())
	}

	alive, err := alivePIDs(evidenceDir)
	if err != nil {
		return p.uncertain("postflight cannot inspect the reported process ids: " + err.Error())
	}
	if len(alive) > 0 {
		return p.failed("proofrail:b4-postflight-live-process:1\n", strings.Join(alive, "\n"))
	}

	legs := []string{
		digestOf("proofrail:b4-postflight-leg-manifest:1\n", digestOf("", string(fresh))),
		digestOf("proofrail:b4-postflight-leg-diff:1\n", digestOf("", string(freshDiff))),
		digestOf("proofrail:b4-postflight-leg-logs:1\n", logsHash),
		digestOf("proofrail:b4-postflight-leg-usage:1\n", usageHash),
		digestOf("proofrail:b4-postflight-leg-secret-scan:1\n", "no marker found"),
		digestOf("proofrail:b4-postflight-leg-processes:1\n", "all reported pids gone"),
	}
	decision := chain.PostflightDecision{
		Outcome:  chain.PostflightPassed,
		Evidence: append(factHashes(request.Facts), legs...),
	}
	p.owner.record("postflight", "passed", map[string]string{"legs": fmt.Sprintf("%d", len(legs))}, nil)
	return decision, nil
}

// factHashes returns the five frozen fact digests in routing order.
func factHashes(facts chain.AgentRunnerFrozenFacts) []string {
	return []string{facts.ManifestHash, facts.DiffHash, facts.LogHash, facts.UsageHash, facts.ProcessStopEvidenceHash}
}

// failed and uncertain are methods so that every non-passing decision is also recorded
// in the run journal: a rejection has to be auditable from the artifacts (the reason,
// the leg and the offending paths), not only from the chain's own state transition.
func (p *postflight) failed(domain, detail string) (chain.PostflightDecision, error) {
	p.owner.record("postflight", "failed", map[string]string{"domain": domain, "detail": detail}, nil)
	return chain.PostflightDecision{
		Outcome:       chain.PostflightFailed,
		ErrorEvidence: []string{evidence.Digest(domain, []byte(detail))},
	}, nil
}

func (p *postflight) uncertain(detail string) (chain.PostflightDecision, error) {
	p.owner.record("postflight", "uncertain", map[string]string{"detail": detail}, nil)
	return chain.PostflightDecision{
		Outcome:       chain.PostflightUncertain,
		ErrorEvidence: []string{digestOf("proofrail:b4-postflight-uncertain:1\n", detail)},
	}, nil
}

// hashLogs digests both log streams exactly as the collector does.
func hashLogs(evidenceDir string) (string, error) {
	var logBytes []byte
	for _, name := range []string{"stdout.log", "stderr.log"} {
		data, err := os.ReadFile(filepath.Join(evidenceDir, "logs", name))
		if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		logBytes = append(logBytes, data...)
	}
	return evidence.Digest("", logBytes), nil
}

// scanWorkspace scans the workspace for secret markers through the production
// scanner, over the entries the snapshot capturer reports.
func (p *postflight) scanWorkspace(ctx context.Context) error {
	entries, err := workspaceEntries(p.workspace, p.snapshotRoot)
	if err != nil {
		return err
	}
	scanner := gates.DefaultSecretScanner()
	for _, entry := range entries {
		if entry.Type != snapshot.EntryRegularFile {
			continue
		}
		data, err := os.ReadFile(filepath.Join(p.workspace, filepath.FromSlash(entry.Path)))
		if err != nil {
			return err
		}
		if _, err := scanner.Scan(ctx, entry.Path, data); err != nil {
			return err
		}
	}
	return nil
}

// alivePIDs reports the recorded pids that are still running.
func alivePIDs(evidenceDir string) ([]string, error) {
	var alive []string
	for _, name := range []string{"stub.pid", "stub-child.pid"} {
		pids, err := readPIDs(filepath.Join(evidenceDir, name))
		if err != nil {
			return nil, err
		}
		for _, pid := range pids {
			identity, err := guard.InspectProcess(pid)
			if err != nil {
				continue
			}
			running, err := guard.ProcessAlive(identity)
			if err != nil {
				continue
			}
			if running {
				alive = append(alive, fmt.Sprintf("%s=%d", name, pid))
			}
		}
	}
	sort.Strings(alive)
	return alive, nil
}
