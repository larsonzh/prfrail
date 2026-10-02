package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/gates"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// This file carries the production AgentRunner workspace capturer and the
// production chain.PostflightPort. Both are plain implementations of existing
// seams: no chain port, no DTO and no wire record is added. The capturer is the
// single source of the production capture/diff byte shape, so the artifacts the
// run collector writes and the artifacts the postflight re-derives are
// byte-identical by construction instead of by convention.

// AgentRunner postflight evidence domains. They are production domains: the
// harness-local proofrail:b4-* domains must never appear here, or an evidence
// reader could not tell a production reconciliation from a harness one.
const (
	agentRunnerPostflightLegManifestDomain   = "proofrail:agent-runner-postflight-leg-manifest:1\n"
	agentRunnerPostflightLegDiffDomain       = "proofrail:agent-runner-postflight-leg-diff:1\n"
	agentRunnerPostflightLegLogDomain        = "proofrail:agent-runner-postflight-leg-log:1\n"
	agentRunnerPostflightLegUsageDomain      = "proofrail:agent-runner-postflight-leg-usage:1\n"
	agentRunnerPostflightLegProcessDomain    = "proofrail:agent-runner-postflight-leg-process:1\n"
	agentRunnerPostflightLegSecretDomain     = "proofrail:agent-runner-postflight-leg-secret:1\n"
	agentRunnerPostflightManifestFailDomain  = "proofrail:agent-runner-postflight-manifest-mismatch:1\n"
	agentRunnerPostflightDiffFailDomain      = "proofrail:agent-runner-postflight-diff-mismatch:1\n"
	agentRunnerPostflightLogFailDomain       = "proofrail:agent-runner-postflight-log-mismatch:1\n"
	agentRunnerPostflightUsageFailDomain     = "proofrail:agent-runner-postflight-usage-mismatch:1\n"
	agentRunnerPostflightProcessFailDomain   = "proofrail:agent-runner-postflight-live-process:1\n"
	agentRunnerPostflightSecretFailDomain    = "proofrail:agent-runner-postflight-secret:1\n"
	agentRunnerPostflightUncertainDomain     = "proofrail:agent-runner-postflight-uncertain:1\n"
	agentRunnerWorkspaceBaselineDomain       = "proofrail:agent-runner-baseline:1\n"
	agentRunnerWorkspacePlaceholderParent    = "sha256:0000000000000000000000000000000000000000000000000000000000000000"
	agentRunnerWorkspaceManifestKind         = "candidate"
	agentRunnerPostflightProcessPIDFileNames = 2
)

// agentRunnerMissingPIDReportDetail is the exact verdict text the console
// assembly also uses for the same condition (internal/console/runtime_agent.go,
// agentRunnerMissingPIDReportDetail). The two copies must stay byte-identical:
// it is a hand-synchronised coupling point, like the pid report file names.
const agentRunnerMissingPIDReportDetail = "the run is missing a required pid report so the managed processes cannot be proven gone"

// agentRunnerPostflightArtifactBound bounds every recorded artifact the
// postflight re-reads: the pre/post manifests, the diff, the usage artifact and
// both log streams. The run side already bounds the same artifacts before they
// are digested (usage and other metadata to agentRunnerMetadataBound, log
// streams to the console's fixed MaxLogBytes of 1 MiB), so 16 MiB is generous
// above every legitimate size while keeping any single re-read bounded. An
// artifact past the bound did not come from a byte-identical replay of what the
// run wrote: the leg that reads it degrades to uncertain instead of loading it
// into memory.
const agentRunnerPostflightArtifactBound int64 = 1 << 24

// agentRunnerParentPIDReportFileName is where the offline stub reports its own
// pid. It is the counterpart of agentRunnerChildPIDReportFileName
// ("stub-child.pid") and is a stub convention, not a wire contract: the name a
// real CLI uses for its parent pid report must be re-evaluated when a real
// candidate is wired in (out of scope for this slice).
const agentRunnerParentPIDReportFileName = "stub.pid"

// agentRunnerDiffEntry and agentRunnerDiffReport are the production diff shape.
// Both the capturer and the postflight marshal through the same functions, so
// the frozen diff.json and the re-derived diff are comparable byte for byte.
type agentRunnerDiffEntry struct {
	Path        string `json:"path"`
	BeforeHash  string `json:"beforeHash,omitempty"`
	AfterHash   string `json:"afterHash,omitempty"`
	BeforeBytes int64  `json:"beforeBytes,omitempty"`
	AfterBytes  int64  `json:"afterBytes,omitempty"`
}

type agentRunnerDiffReport struct {
	Added   []agentRunnerDiffEntry `json:"added"`
	Removed []agentRunnerDiffEntry `json:"removed"`
	Changed []agentRunnerDiffEntry `json:"changed"`
}

// AgentRunnerSnapshotCapturer implements AgentRunnerManifestCapturer over
// internal/snapshot with the production capture protocol. One instance owns one
// workspace binding; every capture is a candidate capture of that workspace with
// a placeholder parent, because the manifest artifacts describe workspace state
// around a run and are not candidate freezes (the acceptance port owns those).
type AgentRunnerSnapshotCapturer struct {
	store   *snapshot.Store
	runID   string
	taskID  string
	attempt int
	counter atomic.Int64
}

var _ AgentRunnerManifestCapturer = (*AgentRunnerSnapshotCapturer)(nil)

func NewAgentRunnerSnapshotCapturer(store *snapshot.Store, runID, taskID string, attempt int) (*AgentRunnerSnapshotCapturer, error) {
	if store == nil {
		return nil, fmt.Errorf("%w: snapshot store is required", ErrInvalidAgentRunnerEvidence)
	}
	if !evidence.ValidID(runID) || !evidence.ValidID(taskID) || attempt < 1 {
		return nil, fmt.Errorf("%w: run, task and attempt are required", ErrInvalidAgentRunnerEvidence)
	}
	return &AgentRunnerSnapshotCapturer{store: store, runID: runID, taskID: taskID, attempt: attempt}, nil
}

// CaptureWorkspaceManifest captures the workspace and returns the marshalled
// snapshot manifest the collector stores as manifest.pre.json/manifest.post.json.
func (capturer *AgentRunnerSnapshotCapturer) CaptureWorkspaceManifest(ctx context.Context, workspaceRoot string) ([]byte, error) {
	if capturer == nil || capturer.store == nil {
		return nil, fmt.Errorf("%w: nil capturer", ErrInvalidAgentRunnerEvidence)
	}
	id := capturer.counter.Add(1)
	manifest, err := agentRunnerCaptureWorkspace(ctx, capturer.store, workspaceRoot, capturer.runID,
		fmt.Sprintf("agent-runner-%s-%d", capturer.runID, id), capturer.taskID, capturer.attempt)
	if err != nil {
		return nil, err
	}
	return json.Marshal(manifest)
}

// DiffWorkspaceManifests reports the entry-level difference between two captured
// manifests. It never reads the filesystem: the manifests are the input.
func (capturer *AgentRunnerSnapshotCapturer) DiffWorkspaceManifests(_ context.Context, before, after []byte) ([]byte, error) {
	return AgentRunnerDiffWorkspaceManifests(before, after)
}

// AgentRunnerDiffWorkspaceManifests is the production manifest diff. It decodes
// both manifests as snapshot manifests and marshals the same report shape the
// postflight re-derives.
func AgentRunnerDiffWorkspaceManifests(before, after []byte) ([]byte, error) {
	beforeManifest, err := snapshot.DecodeSnapshotManifest(before)
	if err != nil {
		return nil, fmt.Errorf("%w: decode before manifest: %v", ErrInvalidAgentRunnerEvidence, err)
	}
	afterManifest, err := snapshot.DecodeSnapshotManifest(after)
	if err != nil {
		return nil, fmt.Errorf("%w: decode after manifest: %v", ErrInvalidAgentRunnerEvidence, err)
	}
	return agentRunnerMarshalDiff(agentRunnerDiffManifests(beforeManifest.Manifest.Entries, afterManifest.Manifest.Entries))
}

// AgentRunnerWorkspaceBaselineDigest is the deterministic baseline digest of a
// workspace: entries only, never the capture timestamp. It is exported because
// the assembly and the caller that emitted the request record must agree on the
// parent snapshot hash the record binds.
func AgentRunnerWorkspaceBaselineDigest(ctx context.Context, store *snapshot.Store, workspaceRoot, runID, taskID string, attempt int) (string, error) {
	manifest, err := agentRunnerCaptureWorkspace(ctx, store, workspaceRoot, runID, "agent-runner-baseline-"+runID, taskID, attempt)
	if err != nil {
		return "", err
	}
	canonical, err := evidence.EncodeCanonical(manifest.Manifest.Entries)
	if err != nil {
		return "", err
	}
	return evidence.Digest(agentRunnerWorkspaceBaselineDomain, canonical), nil
}

func agentRunnerCaptureWorkspace(ctx context.Context, store *snapshot.Store, workspaceRoot, runID, snapshotID, taskID string, attempt int) (snapshot.SnapshotManifest, error) {
	if err := ctx.Err(); err != nil {
		return snapshot.SnapshotManifest{}, err
	}
	parent := agentRunnerWorkspacePlaceholderParent
	return snapshot.Capture(ctx, snapshot.CaptureOptions{
		SourceDir:          workspaceRoot,
		Kind:               agentRunnerWorkspaceManifestKind,
		SnapshotID:         snapshotID,
		RunID:              runID,
		ParentSnapshotHash: &parent,
		Task:               &snapshot.TaskBinding{TaskID: taskID, Attempt: attempt},
		Store:              store,
	})
}

func agentRunnerMarshalDiff(report agentRunnerDiffReport) ([]byte, error) {
	return json.Marshal(report)
}

// agentRunnerDiffManifests computes the entry-level difference of two sorted
// entry lists. Directories are included: a created or removed directory is part
// of the workspace state.
func agentRunnerDiffManifests(before, after []snapshot.Entry) agentRunnerDiffReport {
	report := agentRunnerDiffReport{Added: []agentRunnerDiffEntry{}, Removed: []agentRunnerDiffEntry{}, Changed: []agentRunnerDiffEntry{}}
	beforeByPath := make(map[string]snapshot.Entry, len(before))
	for _, entry := range before {
		beforeByPath[entry.Path] = entry
	}
	afterByPath := make(map[string]snapshot.Entry, len(after))
	for _, entry := range after {
		afterByPath[entry.Path] = entry
	}
	for path, entry := range afterByPath {
		prior, existed := beforeByPath[path]
		if !existed {
			report.Added = append(report.Added, agentRunnerDiffEntry{Path: path, AfterHash: entry.ContentHash, AfterBytes: entry.SizeBytes})
			continue
		}
		if agentRunnerEntrySignature(prior) != agentRunnerEntrySignature(entry) {
			report.Changed = append(report.Changed, agentRunnerDiffEntry{Path: path, BeforeHash: prior.ContentHash, AfterHash: entry.ContentHash, BeforeBytes: prior.SizeBytes, AfterBytes: entry.SizeBytes})
		}
	}
	for path, entry := range beforeByPath {
		if _, exists := afterByPath[path]; !exists {
			report.Removed = append(report.Removed, agentRunnerDiffEntry{Path: path, BeforeHash: entry.ContentHash, BeforeBytes: entry.SizeBytes})
		}
	}
	sort.Slice(report.Added, func(i, j int) bool { return report.Added[i].Path < report.Added[j].Path })
	sort.Slice(report.Removed, func(i, j int) bool { return report.Removed[i].Path < report.Removed[j].Path })
	sort.Slice(report.Changed, func(i, j int) bool { return report.Changed[i].Path < report.Changed[j].Path })
	return report
}

func agentRunnerEntrySignature(entry snapshot.Entry) string {
	return fmt.Sprintf("%s|%s|%s|%d|%t|%t|%s|%s|%s",
		entry.Path, entry.Type, entry.ContentHash, entry.SizeBytes, entry.Executable, entry.ReadOnly, entry.Role, entry.Target, entry.TargetHash)
}

// AgentRunnerPostflight is the production chain.PostflightPort. It is a binding
// port: the request record it is constructed with supplies the requestId and the
// binding the incoming PostflightRequest must agree with, so no field is added
// to the PostflightRequest DTO and no requestId crosses the wire.
type AgentRunnerPostflight struct {
	store        *AgentRunnerReplayStore
	request      AgentRunnerRequestRecord
	snapshotRoot string
	capturer     *AgentRunnerSnapshotCapturer
}

var _ chain.PostflightPort = (*AgentRunnerPostflight)(nil)

// NewAgentRunnerPostflight binds the port to one request record and one snapshot
// root. The record names the evidence directory (requestId) and the binding the
// postflight must agree with.
func NewAgentRunnerPostflight(store *AgentRunnerReplayStore, request AgentRunnerRequestRecord, snapshotRoot string) (*AgentRunnerPostflight, error) {
	if store == nil {
		return nil, fmt.Errorf("%w: replay store is required", ErrInvalidAgentRunnerReplayStore)
	}
	if err := ValidateAgentRunnerRequestRecord(request); err != nil {
		return nil, err
	}
	if strings.TrimSpace(snapshotRoot) == "" {
		return nil, fmt.Errorf("%w: snapshot root is required", ErrInvalidAgentRunnerEvidence)
	}
	snapshotStore, err := snapshot.NewStore(snapshotRoot, 0)
	if err != nil {
		return nil, err
	}
	capturer, err := NewAgentRunnerSnapshotCapturer(snapshotStore, request.Request.RunID, request.Request.TaskID, request.Request.Attempt)
	if err != nil {
		return nil, err
	}
	return &AgentRunnerPostflight{store: store, request: request, snapshotRoot: snapshotRoot, capturer: capturer}, nil
}

// RunPostflight re-derives the frozen facts from the recorded evidence and the
// live workspace, reconciles them item by item, and never writes chain state.
//
// A leg that proves a fact wrong makes the verdict failed; a leg that cannot
// prove a fact makes it uncertain; both are non-passing and are never silent.
func (postflight *AgentRunnerPostflight) RunPostflight(ctx context.Context, request chain.PostflightRequest) (chain.PostflightDecision, error) {
	if postflight == nil || postflight.store == nil || postflight.capturer == nil {
		return chain.PostflightDecision{}, fmt.Errorf("%w: nil postflight port", ErrInvalidAgentRunnerReplayStore)
	}
	bound := postflight.request.Request
	if request.RunID != bound.RunID || request.TaskID != bound.TaskID {
		return chain.PostflightDecision{}, fmt.Errorf("%w: postflight request does not bind the recorded request %s", ErrInvalidAgentRunnerRequest, bound.RequestID)
	}
	if request.Attempt != bound.Attempt {
		return chain.PostflightDecision{}, fmt.Errorf("%w: postflight attempt %d does not bind attempt %d", ErrInvalidAgentRunnerRequest, request.Attempt, bound.Attempt)
	}
	if err := request.Facts.Validate(); err != nil {
		return chain.PostflightDecision{}, err
	}
	evidenceDir, err := postflight.store.RunEvidenceDir(bound.RequestID)
	if err != nil {
		return chain.PostflightDecision{}, err
	}
	storedPre, err := agentRunnerPostflightReadArtifact(filepath.Join(evidenceDir, agentRunnerManifestPreFileName), "pre manifest")
	if err != nil {
		return postflight.uncertain("postflight cannot read the recorded pre manifest: " + err.Error()), nil
	}
	storedPost, err := agentRunnerPostflightReadArtifact(filepath.Join(evidenceDir, agentRunnerManifestPostFileName), "post manifest")
	if err != nil {
		return postflight.uncertain("postflight cannot read the recorded post manifest: " + err.Error()), nil
	}
	storedDiff, err := agentRunnerPostflightReadArtifact(filepath.Join(evidenceDir, agentRunnerDiffFileName), "diff")
	if err != nil {
		return postflight.uncertain("postflight cannot read the recorded diff: " + err.Error()), nil
	}
	// Bind the stored artifacts to the frozen facts before anything is compared:
	// the collection side digests exactly the bytes it wrote, so a stored post
	// manifest or diff that does not digest to its fact is not the recorded
	// artifact. Without this binding a candidate able to write the evidence
	// directory could rewrite the artifact and the workspace together and the
	// manifest and diff legs would reconcile the forgery against itself.
	if evidence.Digest("", storedPost) != request.Facts.ManifestHash {
		return postflight.failed(agentRunnerPostflightManifestFailDomain, evidence.Digest("", storedPost)), nil
	}
	if evidence.Digest("", storedDiff) != request.Facts.DiffHash {
		return postflight.failed(agentRunnerPostflightDiffFailDomain, evidence.Digest("", storedDiff)), nil
	}
	preManifest, err := snapshot.DecodeSnapshotManifest(storedPre)
	if err != nil {
		return postflight.uncertain("postflight cannot decode the recorded pre manifest: " + err.Error()), nil
	}
	postManifest, err := snapshot.DecodeSnapshotManifest(storedPost)
	if err != nil {
		return postflight.uncertain("postflight cannot decode the recorded post manifest: " + err.Error()), nil
	}
	freshBytes, err := postflight.capturer.CaptureWorkspaceManifest(ctx, request.Workspace.Root)
	if err != nil {
		return postflight.uncertain("postflight cannot re-capture the workspace: " + err.Error()), nil
	}
	freshManifest, err := snapshot.DecodeSnapshotManifest(freshBytes)
	if err != nil {
		return postflight.uncertain("postflight cannot decode the fresh manifest: " + err.Error()), nil
	}
	// Leg 1: the workspace must still match the frozen post manifest exactly.
	if changed := agentRunnerChangedPaths(postManifest.Manifest.Entries, freshManifest.Manifest.Entries); len(changed) > 0 {
		return postflight.failed(agentRunnerPostflightManifestFailDomain, strings.Join(changed, "\n")), nil
	}
	// Leg 2: the frozen diff must be exactly the diff of the recorded pre state
	// and the fresh capture.
	freshDiff, err := agentRunnerMarshalDiff(agentRunnerDiffManifests(preManifest.Manifest.Entries, freshManifest.Manifest.Entries))
	if err != nil {
		return postflight.uncertain("postflight cannot recompute the diff: " + err.Error()), nil
	}
	if !bytes.Equal(freshDiff, storedDiff) {
		return postflight.failed(agentRunnerPostflightDiffFailDomain, evidence.Digest("", freshDiff)), nil
	}
	// Leg 3: both log streams, raw bytes, must digest to the frozen log fact.
	logsHash, err := agentRunnerHashLogs(evidenceDir)
	if err != nil {
		return postflight.uncertain("postflight cannot read the log artifacts: " + err.Error()), nil
	}
	if logsHash != request.Facts.LogHash {
		return postflight.failed(agentRunnerPostflightLogFailDomain, logsHash), nil
	}
	// Leg 4: the usage artifact must digest to the frozen usage fact.
	usageHash, err := agentRunnerFileHash(filepath.Join(evidenceDir, agentRunnerUsageFileName))
	if err != nil {
		return postflight.uncertain("postflight cannot read the usage artifact: " + err.Error()), nil
	}
	if usageHash != request.Facts.UsageHash {
		return postflight.failed(agentRunnerPostflightUsageFailDomain, usageHash), nil
	}
	// Leg 5: both pid reports must exist and every pid they report must be gone.
	// The stop-evidence bytes themselves cannot be rebuilt here (guard owns them);
	// only the claim they support is re-verified, and the frozen digest travels
	// with a pass. The "both reports must exist" criterion is word-for-word the
	// console assembly's (internal/console/runtime_agent.go, managedProcessesGone):
	// both reports are candidate written, so a run that left either missing cannot
	// be proven to have stopped its managed processes.
	alive, missing, err := agentRunnerAlivePIDs(evidenceDir)
	if err != nil {
		return postflight.uncertain("postflight cannot inspect the reported process ids: " + err.Error()), nil
	}
	if len(missing) > 0 {
		return postflight.uncertain(agentRunnerMissingPIDReportDetail), nil
	}
	if len(alive) > 0 {
		return postflight.failed(agentRunnerPostflightProcessFailDomain, strings.Join(alive, "\n")), nil
	}
	// Leg 6: the fresh capture must be free of secret markers under the
	// production scanner.
	if err := postflight.scanWorkspace(ctx, request.Workspace.Root, freshManifest.Manifest.Entries); err != nil {
		return postflight.failed(agentRunnerPostflightSecretFailDomain, err.Error()), nil
	}
	legs := []string{
		evidence.Digest(agentRunnerPostflightLegManifestDomain, []byte(evidence.Digest("", freshBytes))),
		evidence.Digest(agentRunnerPostflightLegDiffDomain, []byte(evidence.Digest("", freshDiff))),
		evidence.Digest(agentRunnerPostflightLegLogDomain, []byte(logsHash)),
		evidence.Digest(agentRunnerPostflightLegUsageDomain, []byte(usageHash)),
		evidence.Digest(agentRunnerPostflightLegProcessDomain, []byte("all reported pids gone")),
		evidence.Digest(agentRunnerPostflightLegSecretDomain, []byte("no marker found")),
	}
	return chain.PostflightDecision{
		Outcome:  chain.PostflightPassed,
		Evidence: append(agentRunnerFactHashes(request.Facts), legs...),
	}, nil
}

func (postflight *AgentRunnerPostflight) failed(domain, detail string) chain.PostflightDecision {
	return chain.PostflightDecision{
		Outcome:       chain.PostflightFailed,
		ErrorEvidence: []string{evidence.Digest(domain, []byte(detail))},
	}
}

func (postflight *AgentRunnerPostflight) uncertain(detail string) chain.PostflightDecision {
	return chain.PostflightDecision{
		Outcome:       chain.PostflightUncertain,
		ErrorEvidence: []string{evidence.Digest(agentRunnerPostflightUncertainDomain, []byte(detail))},
	}
}

func (postflight *AgentRunnerPostflight) scanWorkspace(ctx context.Context, workspaceRoot string, entries []snapshot.Entry) error {
	scanner := gates.DefaultSecretScanner()
	for _, entry := range entries {
		if entry.Type != snapshot.EntryRegularFile {
			continue
		}
		data, err := os.ReadFile(filepath.Join(workspaceRoot, filepath.FromSlash(entry.Path)))
		if err != nil {
			return fmt.Errorf("%w: read %s: %v", ErrInvalidAgentRunnerEvidence, entry.Path, err)
		}
		if _, err := scanner.Scan(ctx, entry.Path, data); err != nil {
			return err
		}
	}
	return nil
}

// agentRunnerFactHashes returns the five frozen fact digests in routing order.
func agentRunnerFactHashes(facts chain.AgentRunnerFrozenFacts) []string {
	return []string{facts.ManifestHash, facts.DiffHash, facts.LogHash, facts.UsageHash, facts.ProcessStopEvidenceHash}
}

// agentRunnerChangedPaths lists the paths whose entry differs between the frozen
// post manifest and the fresh capture, including added and removed entries.
func agentRunnerChangedPaths(frozen, fresh []snapshot.Entry) []string {
	frozenByPath := make(map[string]string, len(frozen))
	for _, entry := range frozen {
		frozenByPath[entry.Path] = agentRunnerEntrySignature(entry)
	}
	freshByPath := make(map[string]string, len(fresh))
	for _, entry := range fresh {
		freshByPath[entry.Path] = agentRunnerEntrySignature(entry)
	}
	var changed []string
	for path, signature := range freshByPath {
		if prior, existed := frozenByPath[path]; !existed {
			changed = append(changed, "added "+path)
		} else if prior != signature {
			changed = append(changed, "changed "+path)
		}
	}
	for path := range frozenByPath {
		if _, exists := freshByPath[path]; !exists {
			changed = append(changed, "removed "+path)
		}
	}
	sort.Strings(changed)
	return changed
}

// agentRunnerHashLogs digests both log streams exactly as the collector does:
// raw bytes, no decode, in stdout-then-stderr order.
func agentRunnerHashLogs(evidenceDir string) (string, error) {
	var logBytes []byte
	for _, name := range []string{agentRunnerStdoutFileName, agentRunnerStderrFileName} {
		data, truncated, err := readBoundedAgentRunnerArtifact(filepath.Join(evidenceDir, agentRunnerLogsDirectoryName, name), agentRunnerPostflightArtifactBound)
		if err != nil {
			return "", err
		}
		if truncated {
			return "", fmt.Errorf("%w: recorded log %s exceeds the postflight artifact bound of %d bytes", ErrInvalidAgentRunnerEvidence, name, agentRunnerPostflightArtifactBound)
		}
		logBytes = append(logBytes, data...)
	}
	return evidence.Digest("", logBytes), nil
}

// agentRunnerAlivePIDs reports the recorded pids that are still running and the
// pid reports that are missing. Both reports are written by the candidate, so a
// run that left either report missing cannot be verified and must never read as
// "nothing alive": the caller degrades to uncertain on a non-empty missing list.
// The criterion is word-for-word the console assembly's
// (internal/console/runtime_agent.go, managedProcessesGone).
func agentRunnerAlivePIDs(evidenceDir string) (alive, missing []string, err error) {
	return agentRunnerAlivePIDsWith(evidenceDir, guard.InspectProcess, guard.ProcessAlive)
}

// agentRunnerAlivePIDsWith takes the process observations as parameters so the
// fail-closed decisions can be exercised without racing a real process against
// its own exit. It mirrors agentRunnerVerifyTreeGoneWith in
// internal/adapters/agent_runner_run.go.
func agentRunnerAlivePIDsWith(evidenceDir string, inspect func(int) (guard.ProcessIdentity, error), aliveFn func(guard.ProcessIdentity) (bool, error)) (alive, missing []string, err error) {
	names := []string{agentRunnerChildPIDReportFileName, agentRunnerParentPIDReportFileName}
	for _, name := range names {
		pids, present, err := agentRunnerReadPIDs(filepath.Join(evidenceDir, name))
		if err != nil {
			return nil, nil, err
		}
		if !present {
			missing = append(missing, name)
			continue
		}
		for _, pid := range pids {
			identity, err := inspect(pid)
			if err != nil {
				// A pid that cannot be opened is the normal shape of a process that has
				// already exited: Windows OpenProcess reports ERROR_INVALID_PARAMETER and
				// Linux /proc/<pid> reports ENOENT. The run side reads the same signal the
				// same way (agentRunnerVerifyTreeGoneWith), so failing closed here would
				// make every completed run unverifiable. The liveness answer below is the
				// observation that must never be silently ignored.
				continue
			}
			running, err := aliveFn(identity)
			if err != nil {
				return nil, nil, fmt.Errorf("%w: cannot determine whether pid %d from %s is alive: %v", ErrInvalidAgentRunnerEvidence, pid, name, err)
			}
			if running {
				alive = append(alive, fmt.Sprintf("%s=%d", name, pid))
			}
		}
	}
	sort.Strings(alive)
	return alive, missing, nil
}

// agentRunnerReadPIDs reads the pids one report file lists, and whether the file
// exists at all. A missing file is not an error; the caller decides whether its
// absence is acceptable.
func agentRunnerReadPIDs(path string) ([]int, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}
	pids := make([]int, 0, agentRunnerPostflightProcessPIDFileNames)
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || pid <= 0 {
			return nil, false, fmt.Errorf("%w: unreadable pid %q in %s", ErrInvalidAgentRunnerEvidence, line, path)
		}
		pids = append(pids, pid)
	}
	return pids, true, nil
}

// agentRunnerPostflightReadArtifact reads one recorded artifact under the
// postflight bound. A missing or oversized artifact is an error, never a short
// success: the caller degrades to uncertain rather than reconciling against a
// silently truncated read.
func agentRunnerPostflightReadArtifact(path, label string) ([]byte, error) {
	data, truncated, err := readBoundedAgentRunnerArtifact(path, agentRunnerPostflightArtifactBound)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("%w: recorded %s is missing", ErrInvalidAgentRunnerEvidence, label)
	}
	if truncated {
		return nil, fmt.Errorf("%w: recorded %s exceeds the postflight artifact bound of %d bytes", ErrInvalidAgentRunnerEvidence, label, agentRunnerPostflightArtifactBound)
	}
	return data, nil
}

// agentRunnerFileHash returns the content digest of a recorded artifact, read
// under the postflight bound.
func agentRunnerFileHash(path string) (string, error) {
	data, err := agentRunnerPostflightReadArtifact(path, filepath.Base(path))
	if err != nil {
		return "", err
	}
	return evidence.Digest("", data), nil
}

// WaitAgentRunnerDispatched waits for a run the dispatcher already started,
// collects its evidence and maps it into the terminal outcome. It is the
// post-dispatch half of the offline run: it never starts a process, never writes
// replay-store records (no R, no intent, no identity) and never touches chain
// state, so the caller owns terminal publication.
//
// beforeManifest must be the workspace manifest captured before the dispatch; a
// nil capture is a manifest gap that degrades the outcome exactly as
// AgentRunnerPinnedCLIRun.Run does, never a silent pass. Launcher and Store must
// be the same instances the dispatcher used, because the launch handle and the
// evidence directory must resolve to the run that is already live. The launch id
// is re-derived exactly as the dispatcher derived it, which is the only reason
// this entry can address the live launch without spawning one.
func (run *AgentRunnerPinnedCLIRun) WaitAgentRunnerDispatched(ctx context.Context, request AgentRunnerRunRequest, beforeManifest []byte) (AgentRunnerTerminalOutcome, error) {
	if run == nil || run.Launcher == nil || run.Store == nil {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: launcher and store are required", ErrInvalidAgentRunnerRun)
	}
	if err := ctx.Err(); err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	if !evidence.ValidHash(request.RequestHash) {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: the request record hash is required", ErrInvalidAgentRunnerRun)
	}
	if !evidence.ValidID(request.TaskID) || !evidence.ValidID(request.StepID) || request.Attempt < 1 {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: task, step and attempt are required", ErrInvalidAgentRunnerRun)
	}
	if request.Timeout <= 0 {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: a positive timeout is required", ErrInvalidAgentRunnerRun)
	}
	if request.MaxLogBytes <= 0 {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: a positive log bound is required", ErrInvalidAgentRunnerRun)
	}
	clock := run.Clock
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	// The run is assumed live: the launch id is re-derived from the request id the
	// same way the dispatcher derived it, and nothing is spawned here.
	launchID, err := agentRunnerLaunchID(request.Launch.RequestID)
	if err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	watch, processResult, err := run.Launcher.WaitAgentRunnerProcess(ctx, launchID, request.Launch.RequestID, request.Timeout)
	if err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	// The classification must complete even when the caller cancelled: the terminal
	// receipt is what records the cancellation.
	classifyCtx := context.WithoutCancel(ctx)
	stopProof := processResult.Termination
	switch {
	case agentRunnerNeedsDurableStop(watch, processResult):
		// A natural exit still owes a proven stop, and a run the watchdog could not
		// settle may still be alive: the reviewed stop path returns "already-stopped"
		// evidence for a process that already exited.
		proof, stopErr := run.Launcher.StopAgentRunnerProcess(classifyCtx, launchID)
		if stopErr != nil {
			return AgentRunnerTerminalOutcome{}, stopErr
		}
		stopProof = &proof
	default:
		// The guard stopped it and proved the tree gone, so the launch is only retired
		// (which closes the log artifacts before anything reads them).
		run.Launcher.releaseLaunch(launchID)
	}
	evidenceDir, err := run.Store.RunEvidenceDir(request.Launch.RequestID)
	if err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	collected, err := CollectAgentRunnerEvidence(classifyCtx, AgentRunnerEvidenceSpec{
		EvidenceDir:      evidenceDir,
		WorkspaceRoot:    agentRunnerWorkspaceRoot(request.Launch),
		RequestID:        request.Launch.RequestID,
		Capturer:         request.Capturer,
		StopEvidenceHash: stopProof.Hash,
		BeforeManifest:   beforeManifest,
		MaxLogBytes:      request.MaxLogBytes,
	})
	if err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	// A stop is only "stopped" once the tree is verified gone: a descendant may still
	// be dying when the receipt is written.
	treeGone := stopProof.Outcome == "stopped"
	if treeGone {
		treeGone, err = agentRunnerVerifyTreeGone(classifyCtx, evidenceDir, agentRunnerVerifyWait(request.VerifyWait))
		if err != nil {
			return AgentRunnerTerminalOutcome{}, err
		}
	}
	status, treeStatus, outcomeEvidence := agentRunnerTerminalStatus(request.Launch.RequestID, watch, processResult, stopProof, collected, treeGone)
	outcome := AgentRunnerTerminalOutcome{
		Completion: agentRunnerCompletionFor(request, processResult, collected, status, treeStatus, outcomeEvidence, clock()),
		Settlement: agentRunnerSettlementFor(status, collected),
	}
	if status != "completed" {
		return outcome, nil
	}
	facts, err := collected.FrozenFacts()
	if err != nil {
		// A completion claim without provable facts is refused instead of being
		// reported as a pass, because the chain would reject it anyway.
		return AgentRunnerTerminalOutcome{}, err
	}
	outcome.Facts = &facts
	return outcome, nil
}
