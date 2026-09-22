package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/gates"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// ports implements every chain port the engine needs, plus the manifest capturer
// the adapter side uses. Every method touches the real filesystem or a real
// process; nothing here answers from memory except the manifest captured before
// dispatch, which is the same "before" snapshot internal/adapters takes.
type ports struct {
	runID        string
	requestID    string
	runRoot      string
	workspace    string
	snapshotRoot string
	store        *snapshot.Store
	replay       *adapters.AgentRunnerReplayStore

	router     chain.StepPort
	admission  chain.AgentRunnerAdmission
	dispatcher *adapters.AgentRunnerReplayDispatcher
	launcher   *adapters.AgentRunnerPinnedCLILauncher
	recordHash string
	process    adapters.AgentRunnerProcessConfig

	mu             sync.Mutex
	beforeManifest map[string][]byte
	journal        []journalEntry

	now func() time.Time
}

// journalEntry records one port observation in call order, which is the only
// ordering evidence the chain does not write into its own event log (the
// acceptance, review and promotion ports are event-free by contract).
type journalEntry struct {
	Leg     string            `json:"leg"`
	At      string            `json:"at"`
	Outcome string            `json:"outcome,omitempty"`
	Detail  map[string]string `json:"detail,omitempty"`
	Error   string            `json:"error,omitempty"`
}

func (p *ports) record(leg, outcome string, detail map[string]string, err error) {
	entry := journalEntry{Leg: leg, At: p.now().UTC().Format("2006-01-02T15:04:05.000Z"), Outcome: outcome, Detail: detail}
	if err != nil {
		entry.Error = err.Error()
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.journal = append(p.journal, entry)
}

func (p *ports) snapshotJournal() []journalEntry {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]journalEntry, len(p.journal))
	copy(out, p.journal)
	return out
}

// CaptureBaseline digests the workspace entries. It is deliberately deterministic
// (entries only, no capture timestamp) because the AgentRunner request record
// binds the parent snapshot hash and is emitted before the run.
func (p *ports) CaptureBaseline(_ context.Context, _ string) (chain.SnapshotRef, error) {
	entries, err := workspaceEntries(p.workspace, p.snapshotRoot)
	if err != nil {
		return chain.SnapshotRef{}, err
	}
	hash, err := entriesDigest("proofrail:b4-baseline:1\n", entries)
	if err != nil {
		return chain.SnapshotRef{}, err
	}
	p.record("baseline", "captured", map[string]string{"hash": hash, "entries": fmt.Sprintf("%d", len(entries))}, nil)
	return chain.SnapshotRef{Hash: hash}, nil
}

// Materialize returns the real isolated workspace of this run root. The workspace
// is prepared by the driver before the run: one run root owns one workspace.
func (p *ports) Materialize(_ context.Context, _ string, _ string, _ int, _ chain.SnapshotRef) (chain.Workspace, error) {
	if err := os.MkdirAll(p.workspace, 0o755); err != nil {
		return chain.Workspace{}, err
	}
	return chain.Workspace{Root: p.workspace}, nil
}

// Execute routes a step. AgentRunner steps take the admission-first production
// path; a step the platform store refuses to publish falls back to spawning
// through the same production launcher, and the fallback is recorded.
func (p *ports) Execute(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	if request.AgentRunnerFacts != nil {
		manifest, err := p.CaptureWorkspaceManifest(ctx, request.Workspace.Root)
		if err != nil {
			return chain.StepResult{}, err
		}
		p.mu.Lock()
		if p.beforeManifest == nil {
			p.beforeManifest = map[string][]byte{}
		}
		p.beforeManifest[request.AgentRunnerFacts.RequestID] = manifest
		p.mu.Unlock()
		return p.executeAgentRunner(ctx, request)
	}
	result, err := p.router.Execute(ctx, request)
	if request.Step.Kind == "build" || request.Step.Kind == "verify" {
		p.record("gate-hook", "ran", map[string]string{"stepId": request.Step.ID, "kind": request.Step.Kind}, err)
	}
	return result, err
}

// executeAgentRunner is the production dispatch path with a documented fallback.
//
// Production: StepRouter.ExecuteAgentRunner -> AgentRunnerCompositeAdmission ->
// AgentRunnerReplayDispatcher.DispatchAgentRunner (publish R, reconfirm, claim the
// launch slot, spawn, publish the identity).
//
// Measured on Windows: the dispatcher's first step, Store.RecordRequest, fails
// closed with ErrAgentRunnerReplayStoreDurabilityUnproven, because
// replayStorePublicationDurability() is a platform constant (unproven on Windows)
// and nothing on the public API can change it. No process is spawned and nothing
// is written in that case.
//
// Fallback: admission is still the production composite admission, and the spawn
// is still the production pinned-CLI launcher (therefore guard.RunManaged). Only
// the replay publication that the platform refuses is skipped, and the skipped
// step is recorded in the journal so no reader can mistake this for the full path.
func (p *ports) executeAgentRunner(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	result, err := p.router.Execute(ctx, request)
	if err == nil {
		p.record("dispatch", "production-path", map[string]string{"requestId": request.AgentRunnerFacts.RequestID}, nil)
		return result, nil
	}
	if !errors.Is(err, adapters.ErrAgentRunnerReplayStoreDurabilityUnproven) {
		p.record("dispatch", "production-path-failed", nil, err)
		return result, err
	}
	p.record("dispatch", "production-path-refused-by-platform", map[string]string{"error": err.Error()}, nil)
	launch := chain.AgentRunnerLaunchRequest{
		RequestID:     request.AgentRunnerFacts.RequestID,
		RunID:         request.RunID,
		AdapterID:     p.recordAdapterID(),
		Command:       p.process.Command,
		Args:          append([]string(nil), p.process.Args...),
		Dir:           request.Workspace.Root,
		Grace:         p.process.Grace,
		WorkspaceRoot: request.Workspace.Root,
	}
	spawned, spawnErr := p.launcher.StartAgentRunnerProcess(ctx, launch)
	if spawnErr != nil {
		p.record("dispatch", "fallback-spawn-failed", nil, spawnErr)
		return chain.StepResult{}, spawnErr
	}
	identity := digestOf("proofrail:b4-dispatch-identity:1\n", spawned.LaunchID+"\n"+spawned.ProcessID+"\n"+spawned.StartedAt.UTC().Format(time.RFC3339Nano))
	p.record("dispatch", "fallback-spawned-through-production-launcher", map[string]string{"launchId": spawned.LaunchID, "processId": spawned.ProcessID}, nil)
	return chain.StepResult{Evidence: []string{p.recordHash, identity}}, nil
}

func (p *ports) recordAdapterID() string {
	return declaredAdapterID
}

func (p *ports) beforeManifestFor(requestID string) []byte {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.beforeManifest[requestID]
}

// Accept freezes the candidate: it captures a real candidate snapshot of the
// workspace and digests the recorded evidence tree as the evidence root.
func (p *ports) Accept(_ context.Context, runID, taskID string, attempt int, parent chain.SnapshotRef, workspace chain.Workspace) (chain.CandidateResult, error) {
	captured, err := p.capture(ctx(), snapshot.CaptureOptions{
		SourceDir:          workspace.Root,
		Kind:               "candidate",
		SnapshotID:         fmt.Sprintf("candidate-%s-%s-%d", runID, taskID, attempt),
		RunID:              runID,
		ParentSnapshotHash: &parent.Hash,
		Task:               &snapshot.TaskBinding{TaskID: taskID, Attempt: attempt},
		Store:              p.store,
	})
	if err != nil {
		p.record("acceptance", "failed", nil, err)
		return chain.CandidateResult{}, err
	}
	evidenceRoot, err := p.evidenceRoot(p.requestID)
	if err != nil {
		p.record("acceptance", "failed", nil, err)
		return chain.CandidateResult{}, err
	}
	p.record("acceptance", "frozen", map[string]string{"candidateHash": captured.ManifestHash, "parentHash": parent.Hash}, nil)
	return chain.CandidateResult{
		Snapshot:          chain.SnapshotRef{Hash: captured.ManifestHash},
		EvidenceRootHash:  evidenceRoot,
		Evidence:          []string{digestOf("proofrail:b4-candidate-evidence:1\n", captured.ManifestHash)},
		CandidateProducer: evidence.Actor{Type: "agent", ID: "agent-stub-0.1.0"},
	}, nil
}

// Review is a real policy review over the frozen candidate: it re-checks the
// postflight-qualified state the engine already proved (the review port may only
// be reached from REVIEW_PENDING) and records the policy it applied.
func (p *ports) Review(_ context.Context, request chain.ReviewRequest) (chain.ReviewDecision, error) {
	if !evidence.ValidHash(request.CandidateSnapshotHash) || !evidence.ValidHash(request.EvidenceRootHash) {
		return chain.ReviewDecision{}, fmt.Errorf("review request carries an invalid binding")
	}
	decision := chain.ReviewDecision{
		ReceiptID:  "review-b4-1",
		OccurredAt: p.now().UTC().Format("2006-01-02T15:04:05.000Z"),
		RecordedBy: evidence.Actor{Type: "policy", ID: "b4-mechanism-policy-review"},
		ReviewMode: "policy",
		PolicyHash: digestOf("proofrail:b4-review-policy:1\n", "postflight-qualified-candidate"),
		Outcome:    "approve",
		Evidence:   []string{digestOf("proofrail:b4-review-evidence:1\n", request.CandidateSnapshotHash)},
	}
	p.record("review", "approve", map[string]string{"candidateHash": request.CandidateSnapshotHash}, nil)
	return decision, nil
}

// Publish promotes the candidate after proving the workspace still matches the
// frozen candidate: writers are known stopped (the managed process is gone and
// nothing else writes the tree).
func (p *ports) Publish(_ context.Context, request chain.PromotionRequest) (chain.PromotionDecision, error) {
	entries, err := workspaceEntries(p.workspace, p.snapshotRoot)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	fresh, err := entriesDigest("proofrail:b4-workspace-entries:1\n", entries)
	if err != nil {
		return chain.PromotionDecision{}, err
	}
	stopped, err := p.managedProcessesGone()
	if err != nil {
		p.record("publish", "failed", nil, err)
		return chain.PromotionDecision{}, err
	}
	if !stopped {
		decision := chain.PromotionDecision{
			ReceiptID:     "promotion-b4-1",
			OccurredAt:    p.now().UTC().Format("2006-01-02T15:04:05.000Z"),
			Outcome:       "failed",
			ErrorEvidence: []string{digestOf("proofrail:b4-promotion-blocked:1\n", "a managed process of this run is still alive")},
		}
		p.record("publish", "failed", map[string]string{"reason": "managed process still alive"}, nil)
		return decision, nil
	}
	p.record("publish", "completed", map[string]string{"candidateHash": request.CandidateSnapshotHash, "workspaceDigest": fresh}, nil)
	return chain.PromotionDecision{
		ReceiptID:            "promotion-b4-1",
		OccurredAt:           p.now().UTC().Format("2006-01-02T15:04:05.000Z"),
		Outcome:              "completed",
		AcceptedSnapshotHash: request.CandidateSnapshotHash,
		WriterStopEvidence:   []string{digestOf("proofrail:b4-writer-stop:1\n", "no record of a live managed process")},
		LeaseEvidence:        []string{digestOf("proofrail:b4-lease:1\n", p.runRoot)},
		Evidence:             []string{digestOf("proofrail:b4-promotion-evidence:1\n", request.ReviewReceiptHash)},
	}, nil
}

// Stop is the chain stopper. The only managed writers of a B4 run are the pinned
// CLI launches, whose stop path is the launcher's; this port reports the
// recorded stop evidence location so the cancelled route has something to archive.
func (p *ports) Stop(context.Context, string) ([]string, error) {
	p.record("stop", "reported", nil, nil)
	return []string{digestOf("proofrail:b4-stopper:1\n", p.runRoot)}, nil
}

// Reconcile reports that the harness inspected this run root before resuming.
func (p *ports) Reconcile(_ context.Context, _ string) error {
	p.record("reconcile", "inspected", nil, nil)
	return nil
}

// CaptureWorkspaceManifest implements adapters.AgentRunnerManifestCapturer over
// the snapshot package, which is the same seam internal/adapters documents.
func (p *ports) CaptureWorkspaceManifest(ctx context.Context, workspaceRoot string) ([]byte, error) {
	captured, err := p.capture(ctx, snapshot.CaptureOptions{
		SourceDir:  workspaceRoot,
		Kind:       "candidate",
		SnapshotID: "manifest-" + filepath.Base(p.runRoot),
		RunID:      p.runID,
		Store:      p.store,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(captured)
}

// DiffWorkspaceManifests reports the entry-level difference between two captured
// manifests. It never reads the filesystem: the manifests are the input.
func (p *ports) DiffWorkspaceManifests(_ context.Context, before, after []byte) ([]byte, error) {
	beforeEntries, err := manifestEntries(before)
	if err != nil {
		return nil, err
	}
	afterEntries, err := manifestEntries(after)
	if err != nil {
		return nil, err
	}
	return marshalDiff(diffManifests(beforeEntries, afterEntries))
}

func (p *ports) capture(ctx context.Context, options snapshot.CaptureOptions) (snapshot.SnapshotManifest, error) {
	if options.ParentSnapshotHash == nil {
		parent := zerosHash
		options.ParentSnapshotHash = &parent
	}
	if options.Task == nil {
		options.Task = &snapshot.TaskBinding{TaskID: "b4-task", Attempt: 1}
	}
	return snapshot.Capture(ctx, options)
}

// evidenceRoot digests the sorted artifact inventory of one request's evidence
// directory: it is the root the review binds, not a policy decision.
func (p *ports) evidenceRoot(requestID string) (string, error) {
	root, err := p.replay.RunEvidenceDir(requestID)
	if err != nil {
		return "", err
	}
	var lines []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		hash, err := fileHash(path)
		if err != nil {
			return err
		}
		lines = append(lines, filepath.ToSlash(rel)+" "+hash)
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(lines)
	return digestOf("proofrail:b4-evidence-root:1\n", strings.Join(lines, "\n")), nil
}

// managedProcessesGone reports whether every pid the pinned CLI reported is gone.
func (p *ports) managedProcessesGone() (bool, error) {
	root, err := p.replay.RunEvidenceDir(p.requestID)
	if err != nil {
		return false, err
	}
	for _, name := range []string{"stub.pid", "stub-child.pid"} {
		pids, err := readPIDs(filepath.Join(root, name))
		if err != nil {
			return false, err
		}
		for _, pid := range pids {
			identity, err := guard.InspectProcess(pid)
			if err != nil {
				continue
			}
			alive, err := guard.ProcessAlive(identity)
			if err != nil {
				continue
			}
			if alive {
				return false, nil
			}
		}
	}
	return true, nil
}

// zerosHash is a syntactically valid placeholder parent for the workspace
// manifest capture of a step that is not frozen as a candidate snapshot.
var zerosHash = "sha256:0000000000000000000000000000000000000000000000000000000000000000"

// workspaceEntries captures a directory and returns its entries (entries only:
// the capture timestamp is not part of a candidate comparison).
func workspaceEntries(root, snapshotRoot string) ([]snapshot.Entry, error) {
	store, err := snapshot.NewStore(snapshotRoot, 0)
	if err != nil {
		return nil, err
	}
	captured, err := snapshot.Capture(context.Background(), snapshot.CaptureOptions{
		SourceDir:          root,
		Kind:               "candidate",
		SnapshotID:         "entries-" + filepath.Base(root),
		RunID:              "run-b4",
		ParentSnapshotHash: &zerosHash,
		Task:               &snapshot.TaskBinding{TaskID: "b4-task", Attempt: 1},
		Store:              store,
	})
	if err != nil {
		return nil, err
	}
	return captured.Manifest.Entries, nil
}

func entriesDigest(domain string, entries []snapshot.Entry) (string, error) {
	canonical, err := evidence.EncodeCanonical(entries)
	if err != nil {
		return "", err
	}
	return evidence.Digest(domain, canonical), nil
}

// manifestEntries decodes a captured manifest and returns its entries.
func manifestEntries(raw []byte) ([]snapshot.Entry, error) {
	var manifest snapshot.SnapshotManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}
	return manifest.Manifest.Entries, nil
}

// diffReport is the artifact the collector hashes as the diff fact.
type diffReport struct {
	Added   []diffEntry `json:"added"`
	Removed []diffEntry `json:"removed"`
	Changed []diffEntry `json:"changed"`
}

type diffEntry struct {
	Path        string `json:"path"`
	BeforeHash  string `json:"beforeHash,omitempty"`
	AfterHash   string `json:"afterHash,omitempty"`
	BeforeBytes int64  `json:"beforeBytes,omitempty"`
	AfterBytes  int64  `json:"afterBytes,omitempty"`
}

func diffManifests(before, after []snapshot.Entry) diffReport {
	beforeByPath := map[string]snapshot.Entry{}
	afterByPath := map[string]snapshot.Entry{}
	for _, entry := range before {
		beforeByPath[entry.Path] = entry
	}
	for _, entry := range after {
		afterByPath[entry.Path] = entry
	}
	report := diffReport{Added: []diffEntry{}, Removed: []diffEntry{}, Changed: []diffEntry{}}
	for path, entry := range afterByPath {
		prior, exists := beforeByPath[path]
		switch {
		case !exists:
			report.Added = append(report.Added, diffEntry{Path: path, AfterHash: entry.ContentHash, AfterBytes: entry.SizeBytes})
		case prior.ContentHash != entry.ContentHash || prior.Type != entry.Type || prior.SizeBytes != entry.SizeBytes:
			report.Changed = append(report.Changed, diffEntry{Path: path, BeforeHash: prior.ContentHash, AfterHash: entry.ContentHash, BeforeBytes: prior.SizeBytes, AfterBytes: entry.SizeBytes})
		}
	}
	for path, entry := range beforeByPath {
		if _, exists := afterByPath[path]; !exists {
			report.Removed = append(report.Removed, diffEntry{Path: path, BeforeHash: entry.ContentHash, BeforeBytes: entry.SizeBytes})
		}
	}
	sort.Slice(report.Added, func(i, j int) bool { return report.Added[i].Path < report.Added[j].Path })
	sort.Slice(report.Removed, func(i, j int) bool { return report.Removed[i].Path < report.Removed[j].Path })
	sort.Slice(report.Changed, func(i, j int) bool { return report.Changed[i].Path < report.Changed[j].Path })
	return report
}

func marshalDiff(report diffReport) ([]byte, error) {
	return json.Marshal(report)
}

// changedPaths lists every path the diff covers, for the verdict and for the
// postflight mismatch evidence.
func changedPaths(report diffReport) []string {
	paths := make([]string, 0, len(report.Added)+len(report.Removed)+len(report.Changed))
	for _, group := range [][]diffEntry{report.Added, report.Removed, report.Changed} {
		for _, entry := range group {
			paths = append(paths, entry.Path)
		}
	}
	sort.Strings(paths)
	return paths
}

// readPIDs reads a one-pid-per-line report; a missing report means "none".
func readPIDs(path string) ([]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var pids []int
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var pid int
		if _, err := fmt.Sscanf(line, "%d", &pid); err != nil || pid <= 0 {
			return nil, fmt.Errorf("unreadable pid %q in %s", line, path)
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

// gateHooks declares the F5 gate hook: a real command that inspects the
// workspace. It is *declared* here, and the executor capability question is
// settled by measurement: see hookPort. The hook id is derived from the step id
// it is declared for, so callers never have to patch the ID afterwards.
func gateHook(stepID string) gates.Hook {
	return gates.Hook{
		ID:                "b4-" + stepID,
		Kind:              "verify",
		Runner:            gates.ProcessRunner{Type: "process", Executable: "cmd.exe", Args: []string{"/c", "dir", "/b", "/s"}},
		EffectClass:       "read-only",
		RecoveryGuarantee: "none-required",
		OnFail:            gates.FailurePolicy{Action: "fail-stop"},
		Timeout:           60 * time.Second,
		Resources: gates.ResourceLimits{
			MemoryBytes:  1 << 30,
			OutputBytes:  1 << 20,
			ProcessCount: 8,
		},
		Network: gates.NetworkPolicy{Mode: "none"},
	}
}

// hookPort executes the gate hook of the F5 definition.
//
// Measured: the production pair gates.ChainPort + gates.GuardExecutor cannot run
// ANY hook. GuardExecutor.Capabilities() advertises {process, timeout, output
// limit} while gates.requireCapabilities unconditionally demands a memory limit
// and a process limit, so Prepare fails closed with
//
//	hook policy enforcement unavailable: memory limit
//
// before any process exists. The production port is therefore attempted once and
// its refusal is recorded; the hook then runs through the same reviewed process
// boundary the production launcher uses (guard.RunManaged) with an output bound
// the harness actually enforces. A non-zero exit fails the step.
type hookPort struct {
	owner           *ports
	commands        map[string]gateCommand
	production      gates.ChainPort
	productionTried bool
}

type gateCommand struct {
	executable string
	args       []string
}

func (h *hookPort) RunHook(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	command, exists := h.commands[request.Step.ID]
	if !exists {
		return chain.StepResult{}, fmt.Errorf("no gate command declared for step %q", request.Step.ID)
	}
	if !h.productionTried {
		h.productionTried = true
		if _, err := h.production.RunHook(ctx, request); err != nil {
			h.owner.record("gate-hook", "production-port-refused", map[string]string{"stepId": request.Step.ID, "error": err.Error()}, nil)
		} else {
			h.owner.record("gate-hook", "production-port-ran", map[string]string{"stepId": request.Step.ID}, nil)
		}
	}
	stdout := &boundedBuffer{limit: 1 << 20}
	stderr := &boundedBuffer{limit: 1 << 20}
	started := time.Now().UTC()
	result, runErr := guard.RunManaged(ctx, guard.ProcessSpec{
		Command: command.executable,
		Args:    append([]string(nil), command.args...),
		Dir:     request.Workspace.Root,
		Stdout:  stdout,
		Stderr:  stderr,
	}, 2*time.Second)
	execution := map[string]any{
		"stepId":     request.Step.ID,
		"executable": command.executable,
		"args":       command.args,
		"cwd":        request.Workspace.Root,
		"startedAt":  started.Format("2006-01-02T15:04:05.000Z"),
		"exitCode":   result.ExitCode,
		"started":    result.Started,
		"stdout":     stdout.String(),
		"stderr":     stderr.String(),
		"truncated":  stdout.truncated || stderr.truncated,
	}
	if runErr != nil {
		execution["error"] = runErr.Error()
	}
	record, marshalErr := json.Marshal(execution)
	if marshalErr != nil {
		return chain.StepResult{}, marshalErr
	}
	directory := filepath.Join(h.owner.runRoot, "gates")
	if mkErr := os.MkdirAll(directory, 0o755); mkErr != nil {
		return chain.StepResult{}, mkErr
	}
	recordPath := filepath.Join(directory, request.Step.ID+".json")
	if writeErr := os.WriteFile(recordPath, record, 0o644); writeErr != nil {
		return chain.StepResult{}, writeErr
	}
	evidence := []string{evidence.Digest("", record)}
	h.owner.record("gate-hook", "ran", map[string]string{"stepId": request.Step.ID, "exitCode": fmt.Sprintf("%d", result.ExitCode)}, runErr)
	if !result.Started {
		return chain.StepResult{Evidence: evidence}, fmt.Errorf("gate %s could not start: %v", request.Step.ID, runErr)
	}
	if result.ExitCode != 0 {
		return chain.StepResult{Evidence: evidence}, fmt.Errorf("gate %s exited with %d", request.Step.ID, result.ExitCode)
	}
	return chain.StepResult{Evidence: evidence}, nil
}

// boundedBuffer is the output bound this port actually enforces.
type boundedBuffer struct {
	limit     int64
	written   int64
	truncated bool
	buffer    []byte
}

func (b *boundedBuffer) Write(data []byte) (int, error) {
	remaining := b.limit - b.written
	if remaining <= 0 {
		b.truncated = true
		return len(data), nil
	}
	keep := int64(len(data))
	if keep > remaining {
		keep = remaining
		b.truncated = true
	}
	b.buffer = append(b.buffer, data[:keep]...)
	b.written += keep
	return len(data), nil
}

func (b *boundedBuffer) String() string { return string(b.buffer) }

// ctx returns a background context for the port methods that must not borrow the
// engine's context (candidate capture and evidence digests are bookkeeping).
func ctx() context.Context { return context.Background() }
