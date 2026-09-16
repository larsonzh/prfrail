package adapters

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// fakeManifestCapturer returns deterministic manifests and a diff so the
// collector can be exercised without a content-addressed snapshot store.
type fakeManifestCapturer struct {
	before     []byte
	after      []byte
	diff       []byte
	captureErr error
	diffErr    error
	captures   int
}

func (capturer *fakeManifestCapturer) CaptureWorkspaceManifest(context.Context, string) ([]byte, error) {
	if capturer.captureErr != nil {
		return nil, capturer.captureErr
	}
	if capturer.before == nil {
		// No pre-run value was configured, so each call serves the post-run capture.
		return capturer.after, nil
	}
	capturer.captures++
	if capturer.captures == 1 {
		return capturer.before, nil
	}
	return capturer.after, nil
}

func (capturer *fakeManifestCapturer) DiffWorkspaceManifests(context.Context, []byte, []byte) ([]byte, error) {
	if capturer.diffErr != nil {
		return nil, capturer.diffErr
	}
	return capturer.diff, nil
}

func newAgentRunnerEvidenceDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func writeAgentRunnerArtifactForTest(t *testing.T, dir, name string, data []byte) {
	t.Helper()
	target := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// completeAgentRunnerArtifacts writes the artifacts of a clean run.
func completeAgentRunnerArtifacts(t *testing.T, dir string) {
	t.Helper()
	writeAgentRunnerArtifactForTest(t, dir, filepath.Join(agentRunnerLogsDirectoryName, agentRunnerStdoutFileName), []byte("hello\n"+agentRunnerLogCompletenessMarker+"\n"))
	writeAgentRunnerArtifactForTest(t, dir, filepath.Join(agentRunnerLogsDirectoryName, agentRunnerStderrFileName), []byte(agentRunnerLogCompletenessMarker+"\n"))
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerUsageFileName, []byte(`{"calls":2,"tokens":123,"durationMs":5}`))
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerEventsFileName, []byte("{\"type\":\"session-created\"}\n{\"type\":\"output\"}\n"))
}

func completeAgentRunnerEvidenceSpec(dir string) AgentRunnerEvidenceSpec {
	return AgentRunnerEvidenceSpec{
		EvidenceDir:      dir,
		WorkspaceRoot:    dir,
		RequestID:        "request-one",
		Capturer:         &fakeManifestCapturer{after: []byte("after"), diff: []byte("diff")},
		StopEvidenceHash: evidence.Digest("", []byte("stop")),
		BeforeManifest:   []byte("before"),
		MaxLogBytes:      1 << 20,
	}
}

func TestAgentRunnerEvidenceCollectsCompleteRun(t *testing.T) {
	dir := newAgentRunnerEvidenceDir(t)
	completeAgentRunnerArtifacts(t, dir)
	collected, err := CollectAgentRunnerEvidence(context.Background(), completeAgentRunnerEvidenceSpec(dir))
	if err != nil {
		t.Fatal(err)
	}
	if !collected.LogsComplete || !collected.UsageComplete || !collected.EventsComplete || !collected.ManifestsComplete {
		t.Fatalf("a clean run must be complete: %+v", collected)
	}
	if len(collected.ErrorEvidence) != 0 {
		t.Fatalf("a clean run must carry no error evidence: %v", collected.ErrorEvidence)
	}
	facts, err := collected.FrozenFacts()
	if err != nil {
		t.Fatal(err)
	}
	if facts.ManifestHash != evidence.Digest("", []byte("after")) || facts.DiffHash != evidence.Digest("", []byte("diff")) {
		t.Fatalf("manifest/diff facts must come from the captured artifacts: %+v", facts)
	}
	if err := facts.Validate(); err != nil {
		t.Fatalf("collected facts must validate: %v", err)
	}
	for _, name := range []string{agentRunnerManifestPreFileName, agentRunnerManifestPostFileName, agentRunnerDiffFileName} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("expected artifact %s: %v", name, err)
		}
	}
}

func TestAgentRunnerEvidenceGapsNeverProduceFacts(t *testing.T) {
	cases := map[string]func(t *testing.T, dir string){
		"missing marker": func(t *testing.T, dir string) {
			writeAgentRunnerArtifactForTest(t, dir, filepath.Join(agentRunnerLogsDirectoryName, agentRunnerStdoutFileName), []byte("partial output"))
			writeAgentRunnerArtifactForTest(t, dir, filepath.Join(agentRunnerLogsDirectoryName, agentRunnerStderrFileName), []byte(agentRunnerLogCompletenessMarker+"\n"))
		},
		"missing stdout": func(t *testing.T, dir string) {
			writeAgentRunnerArtifactForTest(t, dir, filepath.Join(agentRunnerLogsDirectoryName, agentRunnerStderrFileName), []byte(agentRunnerLogCompletenessMarker+"\n"))
		},
		"invalid usage": func(t *testing.T, dir string) {
			completeAgentRunnerArtifacts(t, dir)
			writeAgentRunnerArtifactForTest(t, dir, agentRunnerUsageFileName, []byte(`{"calls":2,"tokens":123}`))
		},
		"usage with unknown field": func(t *testing.T, dir string) {
			completeAgentRunnerArtifacts(t, dir)
			writeAgentRunnerArtifactForTest(t, dir, agentRunnerUsageFileName, []byte(`{"calls":2,"tokens":123,"durationMs":5,"extra":1}`))
		},
		"negative usage": func(t *testing.T, dir string) {
			completeAgentRunnerArtifacts(t, dir)
			writeAgentRunnerArtifactForTest(t, dir, agentRunnerUsageFileName, []byte(`{"calls":-1,"tokens":123,"durationMs":5}`))
		},
	}
	for name, arrange := range cases {
		t.Run(name, func(t *testing.T) {
			dir := newAgentRunnerEvidenceDir(t)
			arrange(t, dir)
			collected, err := CollectAgentRunnerEvidence(context.Background(), completeAgentRunnerEvidenceSpec(dir))
			if err != nil {
				t.Fatal(err)
			}
			if len(collected.ErrorEvidence) == 0 {
				t.Fatalf("a gap must record error evidence: %+v", collected)
			}
			if _, err := collected.FrozenFacts(); !errors.Is(err, ErrAgentRunnerEvidenceGap) {
				t.Fatalf("a gap must never yield frozen facts, got %v", err)
			}
		})
	}
}

func TestAgentRunnerEvidenceLogsHashRawBytesAndBoundTruncation(t *testing.T) {
	dir := newAgentRunnerEvidenceDir(t)
	// Invalid UTF-8 must be opaque content, not a failure.
	raw := append([]byte{0xff, 0xfe, 0x00}, []byte("\n"+agentRunnerLogCompletenessMarker+"\n")...)
	writeAgentRunnerArtifactForTest(t, dir, filepath.Join(agentRunnerLogsDirectoryName, agentRunnerStdoutFileName), raw)
	writeAgentRunnerArtifactForTest(t, dir, filepath.Join(agentRunnerLogsDirectoryName, agentRunnerStderrFileName), []byte(agentRunnerLogCompletenessMarker+"\n"))
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerUsageFileName, []byte(`{"calls":1,"tokens":1,"durationMs":1}`))
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerEventsFileName, []byte("{\"type\":\"output\"}\n"))
	spec := completeAgentRunnerEvidenceSpec(dir)
	collected, err := CollectAgentRunnerEvidence(context.Background(), spec)
	if err != nil {
		t.Fatal(err)
	}
	if !collected.LogsComplete {
		t.Fatalf("non-UTF-8 output must not be treated as a gap: %+v", collected)
	}
	want := evidence.Digest("", append(append([]byte{}, raw...), []byte(agentRunnerLogCompletenessMarker+"\n")...))
	if collected.LogHash != want {
		t.Fatalf("log hash must cover the raw bytes of both streams: got %s want %s", collected.LogHash, want)
	}

	bounded := completeAgentRunnerEvidenceSpec(dir)
	bounded.MaxLogBytes = 8
	truncated, err := CollectAgentRunnerEvidence(context.Background(), bounded)
	if err != nil {
		t.Fatal(err)
	}
	if truncated.LogsComplete {
		t.Fatalf("a truncated log must never be complete: %+v", truncated)
	}
	if _, err := truncated.FrozenFacts(); !errors.Is(err, ErrAgentRunnerEvidenceGap) {
		t.Fatalf("a truncated log must not yield facts, got %v", err)
	}

	// A bounded read can stop exactly on the marker, leaving a prefix that looks
	// complete. Only the truncation flag can reject it, so this pins that flag.
	aligned := newAgentRunnerEvidenceDir(t)
	writeAgentRunnerArtifactForTest(t, aligned, filepath.Join(agentRunnerLogsDirectoryName, agentRunnerStdoutFileName), []byte(agentRunnerLogCompletenessMarker+"\n"+strings.Repeat("x", 64)))
	writeAgentRunnerArtifactForTest(t, aligned, filepath.Join(agentRunnerLogsDirectoryName, agentRunnerStderrFileName), []byte(agentRunnerLogCompletenessMarker+"\n"))
	writeAgentRunnerArtifactForTest(t, aligned, agentRunnerUsageFileName, []byte(`{"calls":1,"tokens":1,"durationMs":1}`))
	writeAgentRunnerArtifactForTest(t, aligned, agentRunnerEventsFileName, []byte("{\"type\":\"output\"}\n"))
	alignedSpec := completeAgentRunnerEvidenceSpec(aligned)
	alignedSpec.MaxLogBytes = int64(len(agentRunnerLogCompletenessMarker) + 1)
	collected, err = CollectAgentRunnerEvidence(context.Background(), alignedSpec)
	if err != nil {
		t.Fatal(err)
	}
	if collected.LogsComplete {
		t.Fatalf("a log cut on its marker boundary must not be complete: %+v", collected)
	}
}

func TestAgentRunnerEvidenceTornEventLineIsAGap(t *testing.T) {
	dir := newAgentRunnerEvidenceDir(t)
	completeAgentRunnerArtifacts(t, dir)
	// The torn line comes first, so the scan must continue instead of stopping:
	// an operator request written after an unreadable line still outranks success.
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerEventsFileName, []byte("{\"type\":\"usa\n{\"type\":\"operator-action-required\"}\n"))
	collected, err := CollectAgentRunnerEvidence(context.Background(), completeAgentRunnerEvidenceSpec(dir))
	if err != nil {
		t.Fatal(err)
	}
	if collected.EventsComplete {
		t.Fatalf("a torn line must not be complete: %+v", collected)
	}
	if !collected.OperatorAsk {
		t.Fatalf("the scan must continue past an unreadable line: %+v", collected)
	}
	// An event gap is a gap like any other: it may not be expressed as a completed
	// terminal, so the facts are refused rather than published with a hole behind
	// them.
	if _, err := collected.FrozenFacts(); !errors.Is(err, ErrAgentRunnerEvidenceGap) {
		t.Fatalf("an event gap must refuse the facts, got %v", err)
	}
}

func TestAgentRunnerEvidenceMissingEventLogIsAnEmptyLog(t *testing.T) {
	dir := newAgentRunnerEvidenceDir(t)
	completeAgentRunnerArtifacts(t, dir)
	if err := os.Remove(filepath.Join(dir, agentRunnerEventsFileName)); err != nil {
		t.Fatal(err)
	}
	collected, err := CollectAgentRunnerEvidence(context.Background(), completeAgentRunnerEvidenceSpec(dir))
	if err != nil {
		t.Fatal(err)
	}
	// A CLI that emits no events has nothing to prove, so its absence is complete;
	// only a file that exists and cannot be decoded in full is a gap.
	if !collected.EventsComplete {
		t.Fatalf("an absent event log is an empty one: %+v", collected)
	}
	if _, err := collected.FrozenFacts(); err != nil {
		t.Fatalf("a run without events must still produce facts: %v", err)
	}
}

func TestAgentRunnerEvidenceWithoutAPreManifestIsAGap(t *testing.T) {
	dir := newAgentRunnerEvidenceDir(t)
	completeAgentRunnerArtifacts(t, dir)
	spec := completeAgentRunnerEvidenceSpec(dir)
	spec.BeforeManifest = nil
	collected, err := CollectAgentRunnerEvidence(context.Background(), spec)
	if err != nil {
		t.Fatal(err)
	}
	// Without a pre-run capture the diff would compare the post-run workspace with
	// itself, so the collector refuses to write one instead of claiming "no change".
	if collected.ManifestsComplete || collected.ManifestHash != "" || collected.DiffHash != "" {
		t.Fatalf("a missing pre manifest must be a gap: %+v", collected)
	}
	if _, err := os.Stat(filepath.Join(dir, agentRunnerDiffFileName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a gap must not write a diff artifact: %v", err)
	}
	if len(collected.ErrorEvidence) != 1 {
		t.Fatalf("the gap must be named: %v", collected.ErrorEvidence)
	}
}

func TestAgentRunnerEvidenceEventRequestBeyondTheBoundIsInvisible(t *testing.T) {
	dir := newAgentRunnerEvidenceDir(t)
	completeAgentRunnerArtifacts(t, dir)
	// The read stays bounded, so an operator request that lands past the bound is not
	// seen at all. This is the documented limit of the bounded prefix: the run still
	// degrades (the prefix is incomplete), which is why it can never be reported as a
	// completed run, but the request itself is invisible to the mapping.
	const lineSize = 512
	prefix := "{\"type\":\"output\",\"pad\":\""
	suffix := "\"}\n"
	pad := strings.Repeat("x", lineSize-len(prefix)-len(suffix))
	body := strings.Repeat(prefix+pad+suffix, int(agentRunnerMetadataBound)/lineSize+2)
	body += "{\"type\":\"operator-action-required\"}\n"
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerEventsFileName, []byte(body))
	collected, err := CollectAgentRunnerEvidence(context.Background(), completeAgentRunnerEvidenceSpec(dir))
	if err != nil {
		t.Fatal(err)
	}
	if collected.EventsComplete {
		t.Fatalf("a log beyond the read bound must be incomplete: %+v", collected)
	}
	if collected.OperatorAsk {
		t.Fatal("this test documents that a request beyond the bound is invisible")
	}
	if _, err := collected.FrozenFacts(); !errors.Is(err, ErrAgentRunnerEvidenceGap) {
		t.Fatalf("the incomplete prefix must withhold the facts, got %v", err)
	}
}

func TestAgentRunnerEvidenceCaptureFailuresAreGaps(t *testing.T) {
	dir := newAgentRunnerEvidenceDir(t)
	completeAgentRunnerArtifacts(t, dir)
	failing := completeAgentRunnerEvidenceSpec(dir)
	failing.Capturer = &fakeManifestCapturer{captureErr: errors.New("capture failed")}
	collected, err := CollectAgentRunnerEvidence(context.Background(), failing)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := collected.FrozenFacts(); !errors.Is(err, ErrAgentRunnerEvidenceGap) {
		t.Fatalf("a failed capture must be a gap, got %v", err)
	}

	diffFailing := completeAgentRunnerEvidenceSpec(dir)
	diffFailing.Capturer = &fakeManifestCapturer{before: []byte("b"), after: []byte("a"), diffErr: errors.New("diff failed")}
	collected, err = CollectAgentRunnerEvidence(context.Background(), diffFailing)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := collected.FrozenFacts(); !errors.Is(err, ErrAgentRunnerEvidenceGap) {
		t.Fatalf("a failed diff must be a gap, got %v", err)
	}

	missingCapturer := completeAgentRunnerEvidenceSpec(dir)
	missingCapturer.Capturer = nil
	collected, err = CollectAgentRunnerEvidence(context.Background(), missingCapturer)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := collected.FrozenFacts(); !errors.Is(err, ErrAgentRunnerEvidenceGap) {
		t.Fatalf("an absent capturer must be a gap, got %v", err)
	}
}

func TestAgentRunnerEvidencePinsCompletenessGates(t *testing.T) {
	// A digest that looks present must never substitute for a proven observation:
	// the completeness flags are the proof, and each one alone can withhold facts.
	valid := evidence.Digest("", []byte("artifact"))
	base := AgentRunnerEvidence{
		ManifestHash:            valid,
		DiffHash:                evidence.Digest("", []byte("diff")),
		LogHash:                 evidence.Digest("", []byte("log")),
		UsageHash:               evidence.Digest("", []byte("usage")),
		ProcessStopEvidenceHash: evidence.Digest("", []byte("stop")),
		LogsComplete:            true,
		UsageComplete:           true,
		EventsComplete:          true,
		ManifestsComplete:       true,
	}
	if _, err := base.FrozenFacts(); err != nil {
		t.Fatalf("a complete evidence set must yield facts: %v", err)
	}
	for _, withheld := range []struct {
		name    string
		collect func(AgentRunnerEvidence) AgentRunnerEvidence
	}{
		{name: "usage", collect: func(collected AgentRunnerEvidence) AgentRunnerEvidence {
			collected.UsageComplete = false
			return collected
		}},
		{name: "logs", collect: func(collected AgentRunnerEvidence) AgentRunnerEvidence {
			collected.LogsComplete = false
			return collected
		}},
		{name: "events", collect: func(collected AgentRunnerEvidence) AgentRunnerEvidence {
			collected.EventsComplete = false
			return collected
		}},
		{name: "manifests", collect: func(collected AgentRunnerEvidence) AgentRunnerEvidence {
			collected.ManifestsComplete = false
			return collected
		}},
	} {
		t.Run(withheld.name, func(t *testing.T) {
			if _, err := withheld.collect(base).FrozenFacts(); !errors.Is(err, ErrAgentRunnerEvidenceGap) {
				t.Fatalf("an unproven %s observation must withhold facts, got %v", withheld.name, err)
			}
		})
	}
	missing := base
	missing.DiffHash = ""
	if _, err := missing.FrozenFacts(); !errors.Is(err, ErrAgentRunnerEvidenceGap) {
		t.Fatalf("a missing fact digest must withhold facts, got %v", err)
	}
}

func TestAgentRunnerEvidenceTruncatedEventsAreAGap(t *testing.T) {
	dir := newAgentRunnerEvidenceDir(t)
	completeAgentRunnerArtifacts(t, dir)
	// Every retained line must be well formed, otherwise a decoding failure would
	// mask the truncation flag: the line length is chosen so that the 1MiB bound
	// falls exactly on a line boundary and only the flag can reject the read.
	const lineSize = 256
	prefix := "{\"type\":\"output\",\"pad\":\""
	suffix := "\"}\n"
	pad := strings.Repeat("x", lineSize-len(prefix)-len(suffix))
	line := prefix + pad + suffix
	if len(line) != lineSize {
		t.Fatalf("line construction is off by %d bytes", lineSize-len(line))
	}
	lines := int(agentRunnerMetadataBound)/lineSize + 1
	writeAgentRunnerArtifactForTest(t, dir, agentRunnerEventsFileName, []byte(strings.Repeat(line, lines)))
	info, err := os.Stat(filepath.Join(dir, agentRunnerEventsFileName))
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() <= agentRunnerMetadataBound {
		t.Fatalf("the fixture must exceed the metadata bound: %d", info.Size())
	}
	collected, err := CollectAgentRunnerEvidence(context.Background(), completeAgentRunnerEvidenceSpec(dir))
	if err != nil {
		t.Fatal(err)
	}
	if collected.EventsComplete {
		t.Fatalf("a bounded read that lost the file tail must not report complete events: %+v", collected)
	}
}

func TestAgentRunnerEvidenceRefusesUnusableSpecs(t *testing.T) {
	dir := newAgentRunnerEvidenceDir(t)
	completeAgentRunnerArtifacts(t, dir)
	for name, mutate := range map[string]func(*AgentRunnerEvidenceSpec){
		"missing directory": func(spec *AgentRunnerEvidenceSpec) { spec.EvidenceDir = "" },
		"invalid request":   func(spec *AgentRunnerEvidenceSpec) { spec.RequestID = "Request One" },
		"unbounded logs":    func(spec *AgentRunnerEvidenceSpec) { spec.MaxLogBytes = 0 },
	} {
		t.Run(name, func(t *testing.T) {
			spec := completeAgentRunnerEvidenceSpec(dir)
			mutate(&spec)
			if _, err := CollectAgentRunnerEvidence(context.Background(), spec); !errors.Is(err, ErrInvalidAgentRunnerEvidence) {
				t.Fatalf("expected a fail-closed spec rejection, got %v", err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := CollectAgentRunnerEvidence(ctx, completeAgentRunnerEvidenceSpec(dir)); !errors.Is(err, context.Canceled) {
		t.Fatalf("a cancelled collection must refuse to start, got %v", err)
	}
}
