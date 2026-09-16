package adapters

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	// ErrInvalidAgentRunnerEvidence reports an unusable evidence specification.
	ErrInvalidAgentRunnerEvidence = errors.New("AgentRunner evidence invalid")
	// ErrAgentRunnerEvidenceGap reports a run whose evidence cannot be completed.
	// A gap must degrade the outcome instead of producing empty facts.
	ErrAgentRunnerEvidenceGap = errors.New("AgentRunner evidence gap")
)

// Offline run artifact names. They are store-local evidence: no wire contract
// and no record hash depends on them.
const (
	agentRunnerManifestPreFileName  = "manifest.pre.json"
	agentRunnerManifestPostFileName = "manifest.post.json"
	agentRunnerDiffFileName         = "diff.json"
	agentRunnerLogsDirectoryName    = "logs"
	agentRunnerStdoutFileName       = "stdout.log"
	agentRunnerStderrFileName       = "stderr.log"
	agentRunnerEventsFileName       = "events.jsonl"
	agentRunnerUsageFileName        = "usage.json"
	// agentRunnerChildPIDReportFileName is where the CLI reports the child
	// processes it spawned, one pid per line, so the run can verify the tree.
	agentRunnerChildPIDReportFileName = "stub-child.pid"
)

// agentRunnerLogCompletenessMarker must terminate each log stream. Its absence
// means the writer was interrupted, so the log cannot be called complete.
const agentRunnerLogCompletenessMarker = "proofrail-run-complete"

// AgentRunnerManifestCapturer captures workspace manifests and diffs them. It is
// the seam the launcher implements with the snapshot package; the collector owns
// the artifact layout and the digests.
type AgentRunnerManifestCapturer interface {
	CaptureWorkspaceManifest(ctx context.Context, workspaceRoot string) ([]byte, error)
	DiffWorkspaceManifests(ctx context.Context, before, after []byte) ([]byte, error)
}

// AgentRunnerUsage is the closed usage observation the pinned CLI writes.
type AgentRunnerUsage struct {
	Calls      int `json:"calls"`
	Tokens     int `json:"tokens"`
	DurationMs int `json:"durationMs"`
}

// AgentRunnerEvidence is the collected artifact set of one offline run. Facts
// exist only when every artifact they cover was collected completely: a gap is
// never expressed as an empty digest.
type AgentRunnerEvidence struct {
	ManifestHash            string
	DiffHash                string
	LogHash                 string
	UsageHash               string
	ProcessStopEvidenceHash string
	LogsComplete            bool
	UsageComplete           bool
	EventsComplete          bool
	// ManifestsComplete is true only when the pre manifest, the post manifest and
	// the diff were all collected. A gap here is never expressed as an empty
	// digest: the run must degrade instead of claiming a completed terminal.
	ManifestsComplete bool
	// Usage is the observation behind UsageHash. It is only meaningful when
	// UsageComplete is true; an incomplete observation is never settled as observed.
	Usage         AgentRunnerUsage
	OperatorAsk   bool
	ErrorEvidence []string
}

// AgentRunnerEvidenceSpec pins one collection.
type AgentRunnerEvidenceSpec struct {
	EvidenceDir      string
	WorkspaceRoot    string
	RequestID        string
	Capturer         AgentRunnerManifestCapturer
	StopEvidenceHash string
	// BeforeManifest is the workspace manifest the caller captured before the run
	// started. Its absence is a manifest gap: the collector must never pass off the
	// post-run workspace as the pre-state, because that would under-report the diff.
	BeforeManifest []byte
	// MaxLogBytes bounds each log artifact. Exceeding it truncates the artifact
	// and marks the logs incomplete rather than keeping an unbounded file.
	MaxLogBytes int64
}

// CollectAgentRunnerEvidence reads the artifacts of one finished run. It never
// decides the outcome: it reports what it could and could not prove.
func CollectAgentRunnerEvidence(ctx context.Context, spec AgentRunnerEvidenceSpec) (AgentRunnerEvidence, error) {
	if err := ctx.Err(); err != nil {
		return AgentRunnerEvidence{}, err
	}
	if strings.TrimSpace(spec.EvidenceDir) == "" || !evidence.ValidID(spec.RequestID) {
		return AgentRunnerEvidence{}, fmt.Errorf("%w: evidence directory and requestId are required", ErrInvalidAgentRunnerEvidence)
	}
	if spec.MaxLogBytes <= 0 {
		return AgentRunnerEvidence{}, fmt.Errorf("%w: log bound must be positive", ErrInvalidAgentRunnerEvidence)
	}
	collected := AgentRunnerEvidence{ProcessStopEvidenceHash: spec.StopEvidenceHash}
	if !evidence.ValidHash(collected.ProcessStopEvidenceHash) {
		// An unproven stop is a gap, never an invalid digest inside the evidence.
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerStopGapDomain, []byte(spec.RequestID)))
	}
	if err := collectAgentRunnerLogs(spec, &collected); err != nil {
		return AgentRunnerEvidence{}, err
	}
	if err := collectAgentRunnerUsage(spec, &collected); err != nil {
		return AgentRunnerEvidence{}, err
	}
	if err := collectAgentRunnerEvents(spec, &collected); err != nil {
		return AgentRunnerEvidence{}, err
	}
	if err := collectAgentRunnerManifests(ctx, spec, &collected); err != nil {
		return AgentRunnerEvidence{}, err
	}
	collected.ErrorEvidence = dedupeTerminalEvidence(collected.ErrorEvidence, nil)
	return collected, nil
}

// FrozenFacts returns the five facts, and fails when the run cannot prove them.
// A caller must treat that failure as "not completed", never as a partial fact
// set: the chain refuses a completion without all five.
func (collected AgentRunnerEvidence) FrozenFacts() (AgentRunnerFrozenFacts, error) {
	facts := AgentRunnerFrozenFacts{
		ManifestHash:            collected.ManifestHash,
		DiffHash:                collected.DiffHash,
		LogHash:                 collected.LogHash,
		UsageHash:               collected.UsageHash,
		ProcessStopEvidenceHash: collected.ProcessStopEvidenceHash,
	}
	if !collected.LogsComplete {
		return AgentRunnerFrozenFacts{}, fmt.Errorf("%w: logs are incomplete", ErrAgentRunnerEvidenceGap)
	}
	if !collected.UsageComplete {
		return AgentRunnerFrozenFacts{}, fmt.Errorf("%w: usage is incomplete", ErrAgentRunnerEvidenceGap)
	}
	if !collected.EventsComplete {
		return AgentRunnerFrozenFacts{}, fmt.Errorf("%w: events are incomplete", ErrAgentRunnerEvidenceGap)
	}
	if !collected.ManifestsComplete {
		return AgentRunnerFrozenFacts{}, fmt.Errorf("%w: the manifests are incomplete", ErrAgentRunnerEvidenceGap)
	}
	for _, digest := range facts.hashes() {
		if !evidence.ValidHash(digest) {
			return AgentRunnerFrozenFacts{}, fmt.Errorf("%w: a frozen fact digest is missing", ErrAgentRunnerEvidenceGap)
		}
	}
	if err := facts.Validate(); err != nil {
		return AgentRunnerFrozenFacts{}, err
	}
	return facts, nil
}

// collectAgentRunnerLogs hashes the raw bytes of both streams and requires the
// completeness marker. Bytes are never decoded, so non-UTF-8 output is opaque
// content rather than a failure.
func collectAgentRunnerLogs(spec AgentRunnerEvidenceSpec, collected *AgentRunnerEvidence) error {
	var logBytes []byte
	complete := true
	for _, name := range []string{agentRunnerStdoutFileName, agentRunnerStderrFileName} {
		data, truncated, err := readBoundedAgentRunnerArtifact(filepath.Join(spec.EvidenceDir, agentRunnerLogsDirectoryName, name), spec.MaxLogBytes)
		if err != nil {
			return err
		}
		if truncated || !bytes.HasSuffix(bytes.TrimRight(data, "\r\n"), []byte(agentRunnerLogCompletenessMarker)) {
			complete = false
		}
		logBytes = append(logBytes, data...)
	}
	collected.LogsComplete = complete
	collected.LogHash = evidence.Digest("", logBytes)
	if !complete {
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerLogGapDomain, []byte(spec.RequestID)))
	}
	return nil
}

const agentRunnerLogGapDomain = "proofrail:agent-runner-log-gap:1\n"

// collectAgentRunnerUsage decodes the usage artifact strictly: an unknown field,
// a torn file or a missing artifact all mean the observation is unavailable, and
// an unavailable observation can never be settled as observed.
func collectAgentRunnerUsage(spec AgentRunnerEvidenceSpec, collected *AgentRunnerEvidence) error {
	target := filepath.Join(spec.EvidenceDir, agentRunnerUsageFileName)
	// The bounded read caps memory; whether the observation is usable is decided
	// by the closed-shape and strict-decode checks below, so a truncated file can
	// never pass them.
	data, _, err := readBoundedAgentRunnerArtifact(target, agentRunnerMetadataBound)
	if err != nil {
		return err
	}
	var usage AgentRunnerUsage
	if err := validateAgentRunnerUsageShape(data); err != nil {
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerUsageGapDomain, []byte(spec.RequestID)))
		return nil
	}
	if err := evidence.DecodeStrictJSON(data, &usage); err != nil {
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerUsageGapDomain, []byte(spec.RequestID)))
		return nil
	}
	if usage.Calls < 0 || usage.Tokens < 0 || usage.DurationMs < 0 {
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerUsageGapDomain, []byte(spec.RequestID)))
		return nil
	}
	collected.UsageComplete = true
	collected.Usage = usage
	collected.UsageHash = evidence.Digest("", data)
	return nil
}

// validateAgentRunnerUsageShape enforces a closed observation. A Go struct
// decode cannot distinguish an absent field from an explicit zero, so presence is
// checked on the raw object: a partial observation must never settle as observed.
func validateAgentRunnerUsageShape(data []byte) error {
	var fields map[string]any
	if err := evidence.DecodeStrictJSON(data, &fields); err != nil {
		return err
	}
	if fields == nil || len(fields) != 3 {
		return fmt.Errorf("%w: usage must carry exactly calls, tokens and durationMs", ErrAgentRunnerEvidenceGap)
	}
	for _, key := range []string{"calls", "tokens", "durationMs"} {
		if _, present := fields[key]; !present {
			return fmt.Errorf("%w: usage is missing %s", ErrAgentRunnerEvidenceGap, key)
		}
	}
	return nil
}

const agentRunnerUsageGapDomain = "proofrail:agent-runner-usage-gap:1\n"

// collectAgentRunnerEvents proves the event log is readable line by line. A torn
// or truncated line is a gap, which degrades the run: the completeness of the
// event log is part of the evidence contract, not a decoration. Every line that
// does decode is still scanned, so an operator request inside the retained prefix
// is never missed; the read stays bounded, so a request beyond the bound is not
// seen at all - that prefix is then incomplete, which degrades the run and can
// never produce a completed terminal. Events never become one of the five facts.
func collectAgentRunnerEvents(spec AgentRunnerEvidenceSpec, collected *AgentRunnerEvidence) error {
	data, truncated, err := readBoundedAgentRunnerArtifact(filepath.Join(spec.EvidenceDir, agentRunnerEventsFileName), agentRunnerMetadataBound)
	if err != nil {
		return err
	}
	// A missing artifact is an empty event log, which is complete: a CLI that emits
	// no events has nothing to prove. A file that exists but cannot be decoded in
	// full is the gap.
	complete := !truncated
	for _, line := range strings.Split(strings.TrimRight(string(data), "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var decoded map[string]any
		if err := evidence.DecodeStrictJSON([]byte(line), &decoded); err != nil {
			complete = false
			continue
		}
		// An operator request may never be ignored: it outranks success, so the
		// outcome mapping sees it instead of reporting a completed run.
		if decoded["type"] == "operator-action-required" {
			collected.OperatorAsk = true
		}
	}
	collected.EventsComplete = complete
	if !complete {
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerEventGapDomain, []byte(spec.RequestID)))
	}
	return nil
}

const agentRunnerEventGapDomain = "proofrail:agent-runner-event-gap:1\n"

// collectAgentRunnerManifests writes the pre/post manifest artifacts and the
// diff artifact, and records their digests. Any capture or diff failure is a gap,
// and a missing pre-run capture is a gap too: the collector only ever records a
// diff between a manifest that was captured before the run and one captured after
// it, so it can never claim "no change" from a workspace it saw only once.
func collectAgentRunnerManifests(ctx context.Context, spec AgentRunnerEvidenceSpec, collected *AgentRunnerEvidence) error {
	if spec.Capturer == nil {
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerManifestGapDomain, []byte(spec.RequestID)))
		return nil
	}
	before := spec.BeforeManifest
	if before == nil {
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerManifestGapDomain, []byte(spec.RequestID)))
		return nil
	}
	after, err := spec.Capturer.CaptureWorkspaceManifest(ctx, spec.WorkspaceRoot)
	if err != nil {
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerManifestGapDomain, []byte(spec.RequestID)))
		return nil
	}
	diff, err := spec.Capturer.DiffWorkspaceManifests(ctx, before, after)
	if err != nil {
		collected.ErrorEvidence = append(collected.ErrorEvidence, evidence.Digest(agentRunnerDiffGapDomain, []byte(spec.RequestID)))
		return nil
	}
	for name, data := range map[string][]byte{
		agentRunnerManifestPreFileName:  before,
		agentRunnerManifestPostFileName: after,
		agentRunnerDiffFileName:         diff,
	} {
		if err := writeAgentRunnerArtifact(filepath.Join(spec.EvidenceDir, name), data); err != nil {
			return err
		}
	}
	collected.ManifestHash = evidence.Digest("", after)
	collected.DiffHash = evidence.Digest("", diff)
	collected.ManifestsComplete = true
	return nil
}

const (
	agentRunnerManifestGapDomain = "proofrail:agent-runner-manifest-gap:1\n"
	agentRunnerDiffGapDomain     = "proofrail:agent-runner-diff-gap:1\n"
	agentRunnerStopGapDomain     = "proofrail:agent-runner-stop-gap:1\n"
)

// agentRunnerMetadataBound bounds the small JSON artifacts (usage, events).
const agentRunnerMetadataBound int64 = 1 << 20

// readBoundedAgentRunnerArtifact reads at most bound bytes. A missing artifact is
// an empty, untruncated read: whether absence is a gap is decided by the collector
// that owns the artifact (a missing log fails its marker check, a missing usage
// fails its decode, and events tolerate absence).
func readBoundedAgentRunnerArtifact(path string, bound int64) ([]byte, bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, fmt.Errorf("%w: read %s: %v", ErrInvalidAgentRunnerEvidence, path, err)
	}
	defer func() { _ = file.Close() }()
	info, err := file.Stat()
	if err != nil {
		return nil, false, fmt.Errorf("%w: stat %s: %v", ErrInvalidAgentRunnerEvidence, path, err)
	}
	if info.Size() > bound {
		data := make([]byte, bound)
		if _, err := file.Read(data); err != nil {
			return nil, false, fmt.Errorf("%w: read %s: %v", ErrInvalidAgentRunnerEvidence, path, err)
		}
		return data, true, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false, fmt.Errorf("%w: read %s: %v", ErrInvalidAgentRunnerEvidence, path, err)
	}
	return data, false, nil
}

// writeAgentRunnerArtifact writes one artifact atomically.
func writeAgentRunnerArtifact(path string, data []byte) error {
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("%w: create %s: %v", ErrInvalidAgentRunnerEvidence, directory, err)
	}
	temporary, err := os.CreateTemp(directory, "artifact-*.tmp")
	if err != nil {
		return fmt.Errorf("%w: temp file in %s: %v", ErrInvalidAgentRunnerEvidence, directory, err)
	}
	name := temporary.Name()
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		_ = os.Remove(name)
		return fmt.Errorf("%w: write %s: %v", ErrInvalidAgentRunnerEvidence, name, err)
	}
	if err := temporary.Close(); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("%w: close %s: %v", ErrInvalidAgentRunnerEvidence, name, err)
	}
	if err := os.Rename(name, path); err != nil {
		_ = os.Remove(name)
		return fmt.Errorf("%w: publish %s: %v", ErrInvalidAgentRunnerEvidence, path, err)
	}
	return nil
}
