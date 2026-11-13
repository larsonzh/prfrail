package console

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// P70-S1 --progress tests.
//
// The flag is a pure observer: it may only add lines to the progress sink. Every
// assertion below is chosen so that a decorator which leaks into stdout, alters
// the persisted bytes, or lets a broken sink fail the run turns red.
//
//   - R1: default-off output is unchanged, and --progress streams to stderr.
//   - R2: --json output stays byte-identical while progress streams.
//   - R3: a failing progress writer never fails Append.
//   - R4: the decorator does not change the persisted event bytes.
//   - R5: formatProgressLine keeps its exact shape (placeholders included).
//   - R6: the decorator keeps (and forwards) the chain.SealSink capability.
//   - R7: the flag is documented and accepted by `run`.

// progressStoreEvents builds a small valid event chain (chain -> task -> step) on
// a fixed clock. Append does not verify continuity, but Load does, so an intact
// chain keeps the fixture loadable without special cases.
func progressStoreEvents(t *testing.T) []evidence.StateEvent {
	t.Helper()
	const (
		occurredAt = "2026-11-13T00:00:00.000Z"
		runID      = "run-progress-store"
	)
	records := make([]evidence.StateEvent, 0, 3)
	var previous string
	build := func(entity evidence.Entity, from, to, code string) {
		event := evidence.Event{
			EventID:       fmt.Sprintf("event-%03d", len(records)+1),
			RunID:         runID,
			Sequence:      len(records) + 1,
			OccurredAt:    occurredAt,
			Actor:         evidence.Actor{Type: "system", ID: "proofrail"},
			Entity:        entity,
			FromState:     from,
			ToState:       to,
			InputEvidence: []string{},
			Reason:        evidence.Reason{Code: code},
		}
		if previous != "" {
			hash := previous
			event.PreviousEventHash = &hash
		}
		record, err := evidence.NewStateEvent(event)
		if err != nil {
			t.Fatalf("build progress test event %d: %v", len(records)+1, err)
		}
		previous = record.EventHash
		records = append(records, record)
	}
	build(evidence.Entity{Kind: "chain", RunID: runID}, "NONE", "CREATED", "run-created")
	build(evidence.Entity{Kind: "task", RunID: runID, TaskID: "task-one", Attempt: 1}, "NONE", "PENDING", "task-pending")
	build(evidence.Entity{Kind: "step", RunID: runID, TaskID: "task-one", StepID: "step-one", Attempt: 1}, "NONE", "PENDING", "step-pending")
	return records
}

// progressBrokenWriter is the sink shape R3 needs: every write fails.
type progressBrokenWriter struct{}

func (progressBrokenWriter) Write([]byte) (int, error) {
	return 0, errors.New("progress sink is broken")
}

// progressTestChainConfig writes the default single-noop-task chain config into
// root and returns its path. That definition emits one chain, one task and one
// step entity, which is exactly the three progress lines R1/R2 require.
func progressTestChainConfig(t *testing.T, root string) string {
	t.Helper()
	content, err := EncodeChainConfig(DefaultChainConfig("progress-chain"))
	if err != nil {
		t.Fatalf("encode the test chain config: %v", err)
	}
	path := filepath.Join(root, DefaultChainConfigName)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write the test chain config: %v", err)
	}
	return path
}

// progressEntityCounts counts progress lines per entity keyword. The line shape
// is fixed by formatProgressLine, so the second field after the prefix is always
// the entity kind.
func progressEntityCounts(stderr string) map[string]int {
	counts := map[string]int{}
	for _, line := range strings.Split(stderr, "\n") {
		if !strings.HasPrefix(line, "progress: ") {
			continue
		}
		fields := strings.Split(strings.TrimPrefix(line, "progress: "), " ")
		if len(fields) < 2 {
			counts[""]++
			continue
		}
		counts[fields[1]]++
	}
	return counts
}

func progressLineTotal(stderr string) int {
	total := 0
	for _, line := range strings.Split(stderr, "\n") {
		if strings.HasPrefix(line, "progress: ") {
			total++
		}
	}
	return total
}

// maskRunDir replaces the run directory token so two invocations that only differ
// by --run-dir can still be compared byte for byte.
func maskRunDir(output, runDir string) string {
	return strings.ReplaceAll(output, runDir, "<run-dir>")
}

// progressSummaryJSON marshals a run summary and normalises the two fields that
// identify one run instead of describing its outcome:
//
//   - eventLogPath is derived from the run directory, which is masked down to the
//     "<run-dir>" placeholder while its run-relative tail is kept;
//   - lastEventHash covers the per-run evidence inventory (evidenceRootHash walks
//     the request evidence directory), so two runs on two run directories hash to
//     different values by construction even when every transition is identical.
//
// The bytes of every other field -- runId, chainState, sequence and both
// state-count projections -- are returned untouched, so a decorator that changes
// the projection still turns the comparison red. The guards keep the comparison
// honest: an empty image, or one that lost the chainState key or a real run
// directory, would make the equality assertion vacuous.
func progressSummaryJSON(t *testing.T, summary RunSummary, runDir string) string {
	t.Helper()
	image, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("marshal the run summary: %v", err)
	}
	if len(image) == 0 {
		t.Fatal("the run summary must marshal to a non-empty JSON image")
	}
	if !strings.Contains(string(image), `"chainState"`) {
		t.Fatalf("the JSON image must carry the chainState key: %s", image)
	}

	fields := map[string]json.RawMessage{}
	if err := json.Unmarshal(image, &fields); err != nil {
		t.Fatalf("decode the JSON image %s: %v", image, err)
	}

	var eventLogPath string
	if err := json.Unmarshal(fields["eventLogPath"], &eventLogPath); err != nil {
		t.Fatalf("decode eventLogPath from the JSON image %s: %v", image, err)
	}
	if !strings.HasPrefix(eventLogPath, runDir) {
		t.Fatalf("eventLogPath %q must live under the run directory %q", eventLogPath, runDir)
	}
	maskedPath, err := json.Marshal("<run-dir>" + strings.TrimPrefix(eventLogPath, runDir))
	if err != nil {
		t.Fatalf("mask eventLogPath %q: %v", eventLogPath, err)
	}
	fields["eventLogPath"] = maskedPath

	var lastEventHash string
	if err := json.Unmarshal(fields["lastEventHash"], &lastEventHash); err != nil || lastEventHash == "" {
		t.Fatalf("the JSON image must carry a non-empty lastEventHash: %s", image)
	}
	fields["lastEventHash"] = json.RawMessage(`"<run-identity>"`)

	normalised, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("re-marshal the normalised summary image: %v", err)
	}
	return string(normalised)
}

// progressRunReportJSON applies progressSummaryJSON to the summary a `run --json`
// report carries under "data" and keeps the surrounding envelope, so the whole
// serialised report -- command, ok, exitCode, message and normalised summary --
// is compared byte for byte.
func progressRunReportJSON(t *testing.T, report, runDir string) string {
	t.Helper()
	if !json.Valid([]byte(report)) {
		t.Fatalf("the report must be valid JSON: %s", report)
	}
	envelope := map[string]json.RawMessage{}
	if err := json.Unmarshal([]byte(report), &envelope); err != nil {
		t.Fatalf("decode the report envelope %s: %v", report, err)
	}
	data, ok := envelope["data"]
	if !ok {
		t.Fatalf("the report must carry a data object: %s", report)
	}
	summary := RunSummary{}
	if err := json.Unmarshal(data, &summary); err != nil {
		t.Fatalf("decode the report data as a run summary: %v", err)
	}
	envelope["data"] = json.RawMessage(progressSummaryJSON(t, summary, runDir))
	normalised, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("re-marshal the normalised report: %v", err)
	}
	return string(normalised)
}

// progressLastEventHash returns the hash of the last record an event log holds,
// which is what a run summary must report as its LastEventHash.
func progressLastEventHash(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the event log %s: %v", path, err)
	}
	hash := ""
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event evidence.StateEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("decode the event log line %q: %v", line, err)
		}
		hash = event.EventHash
	}
	if hash == "" {
		t.Fatalf("the event log %s must hold at least one event", path)
	}
	return hash
}

// R1: --progress defaults to off (stdout unchanged, stderr empty) and, when on,
// streams one line per chain/task/step entity to stderr only.
func TestProgressFlagStreamsToStderrAndKeepsStdout(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := progressTestChainConfig(t, root)
	const runID = "run-progress-r1"

	plainDir := filepath.Join(root, "tmp", "prfrail-runs", "r1-plain")
	code, plainStdout, plainStderr := runCLI(t, cli, "run", "--chain", chainPath, "--run-id", runID, "--run-dir", plainDir)
	if code != exitSuccess {
		t.Fatalf("the run without --progress must succeed: code=%d stderr=%s", code, plainStderr)
	}
	if plainStderr != "" {
		t.Fatalf("without --progress stderr must stay empty, got %q", plainStderr)
	}

	progressDir := filepath.Join(root, "tmp", "prfrail-runs", "r1-progress")
	code, progressStdout, progressStderr := runCLI(t, cli, "run", "--chain", chainPath, "--run-id", runID, "--run-dir", progressDir, "--progress")
	if code != exitSuccess {
		t.Fatalf("the run with --progress must succeed: code=%d stderr=%s", code, progressStderr)
	}

	// Not one progress byte may reach stdout, with or without the flag.
	if strings.Contains(plainStdout, "progress: ") || strings.Contains(progressStdout, "progress: ") {
		t.Fatalf("progress must never reach stdout:\nplain:\n%s\nwith progress:\n%s", plainStdout, progressStdout)
	}
	// stdout stays byte-identical once the run directory -- the only intended
	// difference between the two invocations -- is masked out.
	if got, want := maskRunDir(progressStdout, progressDir), maskRunDir(plainStdout, plainDir); got != want {
		t.Fatalf("--progress changed stdout:\n got:\n%s\nwant:\n%s", got, want)
	}

	counts := progressEntityCounts(progressStderr)
	for _, kind := range []string{"chain", "task", "step"} {
		if counts[kind] == 0 {
			t.Fatalf("expected at least one %q progress line, got %+v (stderr=%q)", kind, counts, progressStderr)
		}
	}
	if total := progressLineTotal(progressStderr); total < 3 {
		t.Fatalf("expected at least 3 progress lines, got %d (stderr=%q)", total, progressStderr)
	}
}

// R2: the JSON report is byte-identical with and without --progress.
//
// Both invocations share one run directory (recreated in between) on purpose: a
// second directory would change every path-derived hash in the report, so the
// comparison would no longer isolate the flag.
func TestProgressFlagKeepsJSONOutputIdentical(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := progressTestChainConfig(t, root)
	const runID = "run-progress-r2"
	runDir := filepath.Join(root, "tmp", "prfrail-runs", "r2")

	code, plainJSON, plainStderr := runCLI(t, cli, "run", "--chain", chainPath, "--run-id", runID, "--run-dir", runDir, "--json")
	if code != exitSuccess {
		t.Fatalf("the JSON run without --progress must succeed: code=%d stderr=%s", code, plainStderr)
	}
	if plainStderr != "" {
		t.Fatalf("without --progress stderr must stay empty, got %q", plainStderr)
	}
	if err := os.RemoveAll(runDir); err != nil {
		t.Fatalf("clear the run directory between the two invocations: %v", err)
	}

	code, progressJSON, progressStderr := runCLI(t, cli, "run", "--chain", chainPath, "--run-id", runID, "--run-dir", runDir, "--json", "--progress")
	if code != exitSuccess {
		t.Fatalf("the JSON run with --progress must succeed: code=%d stderr=%s", code, progressStderr)
	}

	if !json.Valid([]byte(plainJSON)) || !json.Valid([]byte(progressJSON)) {
		t.Fatalf("both reports must stay valid JSON:\nplain:\n%s\nwith progress:\n%s", plainJSON, progressJSON)
	}
	response := decodeResponse(t, progressJSON)
	if !response.OK {
		t.Fatalf("the JSON report must stay successful: %s", progressJSON)
	}
	if progressJSON != plainJSON {
		t.Fatalf("--progress changed the JSON report:\n got:\n%s\nwant:\n%s", progressJSON, plainJSON)
	}
	if strings.Contains(progressJSON, "progress: ") {
		t.Fatalf("progress must never reach the JSON report: %s", progressJSON)
	}
	if total := progressLineTotal(progressStderr); total < 3 {
		t.Fatalf("expected at least 3 progress lines on stderr, got %d (stderr=%q)", total, progressStderr)
	}
}

// R3: a broken progress sink is swallowed -- Append returns nil and the event log
// is byte-identical to what the undecorated store produced.
func TestProgressStoreSwallowsWriterFailure(t *testing.T) {
	root := t.TempDir()
	events := progressStoreEvents(t)
	plainPath := filepath.Join(root, "plain", "state-events.jsonl")
	decoratedPath := filepath.Join(root, "decorated", "state-events.jsonl")

	plainStore, err := chain.NewFileEventStore(plainPath)
	if err != nil {
		t.Fatalf("open the plain event store: %v", err)
	}
	for _, event := range events {
		if err := plainStore.Append(context.Background(), event); err != nil {
			t.Fatalf("append to the plain store: %v", err)
		}
	}

	innerStore, err := chain.NewFileEventStore(decoratedPath)
	if err != nil {
		t.Fatalf("open the decorated event store: %v", err)
	}
	decorated := newProgressStore(innerStore, progressBrokenWriter{})
	for _, event := range events {
		if err := decorated.Append(context.Background(), event); err != nil {
			t.Fatalf("a broken progress sink must never fail the append: %v", err)
		}
	}

	plainBytes, err := os.ReadFile(plainPath)
	if err != nil {
		t.Fatalf("read the plain event log: %v", err)
	}
	decoratedBytes, err := os.ReadFile(decoratedPath)
	if err != nil {
		t.Fatalf("read the decorated event log: %v", err)
	}
	if len(plainBytes) == 0 {
		t.Fatal("the plain event log must not be empty")
	}
	if !bytes.Equal(plainBytes, decoratedBytes) {
		t.Fatalf("a broken sink changed the persisted events:\nplain:\n%s\ndecorated:\n%s", plainBytes, decoratedBytes)
	}
}

// R4: the decorator forwards every event untouched, so a healthy sink changes
// nothing in the persisted bytes.
func TestProgressStoreForwardsEventsUnchanged(t *testing.T) {
	root := t.TempDir()
	events := progressStoreEvents(t)
	plainPath := filepath.Join(root, "plain", "state-events.jsonl")
	decoratedPath := filepath.Join(root, "decorated", "state-events.jsonl")

	plainStore, err := chain.NewFileEventStore(plainPath)
	if err != nil {
		t.Fatalf("open the plain event store: %v", err)
	}
	for _, event := range events {
		if err := plainStore.Append(context.Background(), event); err != nil {
			t.Fatalf("append to the plain store: %v", err)
		}
	}

	sink := &bytes.Buffer{}
	innerStore, err := chain.NewFileEventStore(decoratedPath)
	if err != nil {
		t.Fatalf("open the decorated event store: %v", err)
	}
	decorated := newProgressStore(innerStore, sink)
	for _, event := range events {
		if err := decorated.Append(context.Background(), event); err != nil {
			t.Fatalf("append through the decorator: %v", err)
		}
	}
	if sink.Len() == 0 {
		t.Fatal("the decorator must emit one progress line per appended event")
	}

	plainBytes, err := os.ReadFile(plainPath)
	if err != nil {
		t.Fatalf("read the plain event log: %v", err)
	}
	decoratedBytes, err := os.ReadFile(decoratedPath)
	if err != nil {
		t.Fatalf("read the decorated event log: %v", err)
	}
	// Without this guard two empty logs would compare equal and the assertion
	// below would hold vacuously.
	if len(plainBytes) == 0 {
		t.Fatal("the plain event log must not be empty")
	}
	if !bytes.Equal(plainBytes, decoratedBytes) {
		t.Fatalf("the decorator changed the persisted events:\nplain:\n%s\ndecorated:\n%s", plainBytes, decoratedBytes)
	}
}

// R5: the progress line shape is pinned exactly, placeholders included.
func TestFormatProgressLineSnapshot(t *testing.T) {
	cases := []struct {
		name string
		body evidence.Event
		want string
	}{
		{
			name: "chain",
			body: evidence.Event{
				OccurredAt: "2026-11-13T00:00:00.000Z",
				Entity:     evidence.Entity{Kind: "chain", RunID: "run-progress"},
				FromState:  "NONE",
				ToState:    "CREATED",
				Reason:     evidence.Reason{Code: "run-created"},
			},
			want: "progress: 2026-11-13T00:00:00.000Z chain NONE->CREATED run-created",
		},
		{
			name: "task",
			body: evidence.Event{
				OccurredAt: "2026-11-13T00:00:00.000Z",
				Entity:     evidence.Entity{Kind: "task", RunID: "run-progress", TaskID: "task-one", Attempt: 1},
				FromState:  "NONE",
				ToState:    "PENDING",
				Reason:     evidence.Reason{Code: "task-pending"},
			},
			want: "progress: 2026-11-13T00:00:00.000Z task task-one NONE->PENDING task-pending",
		},
		{
			name: "step",
			body: evidence.Event{
				OccurredAt: "2026-11-13T00:00:00.000Z",
				Entity:     evidence.Entity{Kind: "step", RunID: "run-progress", TaskID: "task-one", StepID: "step-one", Attempt: 1},
				FromState:  "NONE",
				ToState:    "PENDING",
				Reason:     evidence.Reason{Code: "step-pending"},
			},
			want: "progress: 2026-11-13T00:00:00.000Z step task-one/step-one NONE->PENDING step-pending",
		},
		{
			// Every optional field is empty: each one must fall back to the "-"
			// placeholder so the line keeps its column shape.
			name: "empty fields",
			body: evidence.Event{Entity: evidence.Entity{Kind: "chain", RunID: "run-progress"}},
			want: "progress: - chain -->- -",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := formatProgressLine(evidence.StateEvent{Event: tc.body})
			if got != tc.want {
				t.Fatalf("unexpected progress line:\n got: %q\nwant: %q", got, tc.want)
			}
		})
	}
}

// R6: the decorator must keep the chain.SealSink capability -- internal/chain
// type-asserts the event store for the seal barrier -- and forward both methods,
// the write path included. A wrapper that forwards only LoadTaskSeal answers every
// read with the bare store's data while persisting nothing, so the read path alone
// cannot prove the capability survived.
func TestProgressStoreKeepsSealSink(t *testing.T) {
	root := t.TempDir()

	// R6a: the capability is kept, and a missing seal is reported identically by
	// both stores. This pins the read path's error shape.
	t.Run("keeps capability and missing-seal shape", func(t *testing.T) {
		store, err := chain.NewFileEventStore(filepath.Join(root, "state-events.jsonl"))
		if err != nil {
			t.Fatalf("open the event store: %v", err)
		}
		decorated := newProgressStore(store, io.Discard)

		if _, ok := interface{}(store).(chain.SealSink); !ok {
			t.Fatal("the bare file event store must be a chain.SealSink")
		}
		sink, ok := interface{}(decorated).(chain.SealSink)
		if !ok {
			t.Fatal("progressStore must keep the chain.SealSink capability")
		}

		ctx := context.Background()
		wantRecord, wantHash, wantErr := store.LoadTaskSeal(ctx, "run-seal", "task-one", 1)
		gotRecord, gotHash, gotErr := sink.LoadTaskSeal(ctx, "run-seal", "task-one", 1)
		if wantErr == nil || gotErr == nil {
			t.Fatalf("a missing seal must be reported by both stores: bare=%v decorated=%v", wantErr, gotErr)
		}
		if !errors.Is(wantErr, chain.ErrSealMissing) || !errors.Is(gotErr, chain.ErrSealMissing) {
			t.Fatalf("both stores must report ErrSealMissing: bare=%v decorated=%v", wantErr, gotErr)
		}
		if gotErr.Error() != wantErr.Error() {
			t.Fatalf("the forwarded error shape must match:\n got: %v\nwant: %v", gotErr, wantErr)
		}
		if gotHash != wantHash || !reflect.DeepEqual(gotRecord, wantRecord) {
			t.Fatalf("the forwarded seal must match: got (%+v, %q) want (%+v, %q)", gotRecord, gotHash, wantRecord, wantHash)
		}
	})

	// R6b: the write path is forwarded for real. Both stores seal the same record
	// into their own directory. A decorator whose SealTaskEvidence returns
	// ("", nil) persists nothing, so the decorated seal directory stays empty and
	// the assertions below turn red.
	t.Run("forwards seal write path", func(t *testing.T) {
		const (
			runID   = "run-seal-write"
			taskID  = "task-one"
			attempt = 1
		)
		plainDir := filepath.Join(root, "plain")
		decoratedDir := filepath.Join(root, "decorated")

		plainStore, err := chain.NewFileEventStore(filepath.Join(plainDir, "state-events.jsonl"))
		if err != nil {
			t.Fatalf("open the bare event store: %v", err)
		}
		innerStore, err := chain.NewFileEventStore(filepath.Join(decoratedDir, "state-events.jsonl"))
		if err != nil {
			t.Fatalf("open the decorated event store: %v", err)
		}
		if _, ok := interface{}(innerStore).(chain.SealSink); !ok {
			t.Fatal("the bare file event store must be a chain.SealSink")
		}
		sink, ok := interface{}(newProgressStore(innerStore, io.Discard)).(chain.SealSink)
		if !ok {
			t.Fatal("progressStore must keep the chain.SealSink capability")
		}

		ctx := context.Background()
		record := progressSealRecord(runID, taskID, attempt)

		bareHash, bareErr := plainStore.SealTaskEvidence(ctx, record)
		if bareErr != nil {
			t.Fatalf("the bare store must seal the record: %v", bareErr)
		}
		forwardedHash, forwardedErr := sink.SealTaskEvidence(ctx, record)
		if forwardedErr != nil {
			t.Fatalf("the decorator must forward SealTaskEvidence: %v", forwardedErr)
		}
		if bareHash == "" || forwardedHash == "" {
			t.Fatalf("both stores must return a seal digest: bare=%q decorated=%q", bareHash, forwardedHash)
		}
		if forwardedHash != bareHash {
			t.Fatalf("the forwarded seal digest must match:\n got: %q\nwant: %q", forwardedHash, bareHash)
		}

		// The record must land at the same relative place below each store
		// directory, byte for byte: same run, same task, same attempt.
		bareSeals := progressSealFiles(t, plainDir)
		forwardedSeals := progressSealFiles(t, decoratedDir)
		if len(bareSeals) == 0 {
			t.Fatal("the bare store must have written one seal file")
		}
		if !reflect.DeepEqual(forwardedSeals, bareSeals) {
			t.Fatalf("the decorator must seal the same relative path with the same bytes:\n got: %v\nwant: %v", sealFileNames(forwardedSeals), sealFileNames(bareSeals))
		}

		// Reading back through each store must yield the same record and digest.
		wantRecord, wantHash, err := plainStore.LoadTaskSeal(ctx, runID, taskID, attempt)
		if err != nil {
			t.Fatalf("read the seal back from the bare store: %v", err)
		}
		gotRecord, gotHash, err := sink.LoadTaskSeal(ctx, runID, taskID, attempt)
		if err != nil {
			t.Fatalf("read the seal back through the decorator: %v", err)
		}
		if gotHash != wantHash {
			t.Fatalf("the seal digest must survive the decorator: got %q want %q", gotHash, wantHash)
		}
		if !reflect.DeepEqual(gotRecord, wantRecord) {
			t.Fatalf("the seal record must survive the decorator:\n got: %+v\nwant: %+v", gotRecord, wantRecord)
		}
		if wantHash != bareHash {
			t.Fatalf("the stored seal digest must equal the one SealTaskEvidence returned: stored %q returned %q", wantHash, bareHash)
		}
	})
}

// progressSealRecord builds one valid chain.SealRecord. Every hash below is a real
// digest, so the record passes the seal validator instead of tripping its "not a
// hash" guard; the content is fixed, so two stores sealing it must agree.
func progressSealRecord(runID, taskID string, attempt int) chain.SealRecord {
	digest := func(label string) string {
		return evidence.Digest("proofrail:progress-test:1\n", []byte(label))
	}
	return chain.SealRecord{
		RunID:                runID,
		TaskID:               taskID,
		Attempt:              attempt,
		ParentHash:           digest("parent"),
		CandidateHash:        digest("candidate"),
		EvidenceRootHash:     digest("evidence-root"),
		ContextKey:           digest("context-key"),
		ReviewReceiptHash:    digest("review-receipt"),
		PromotionReceiptHash: digest("promotion-receipt"),
		ReviewDecision:       "approved",
		PromotionDecision:    "approved",
		Evidence:             []string{digest("closure-one"), digest("closure-two")},
	}
}

// progressSealFiles maps every seal file below storeDir to its content, keyed by
// the slash-separated path relative to storeDir. Comparing two stores this way
// pins location and bytes without hard-coding the seal directory layout.
func progressSealFiles(t *testing.T, storeDir string) map[string]string {
	t.Helper()
	files := map[string]string{}
	walkErr := filepath.WalkDir(storeDir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".seal.json") {
			return nil
		}
		relative, relErr := filepath.Rel(storeDir, path)
		if relErr != nil {
			return relErr
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		files[filepath.ToSlash(relative)] = string(data)
		return nil
	})
	if walkErr != nil {
		t.Fatalf("collect the seal files below %s: %v", storeDir, walkErr)
	}
	return files
}

// sealFileNames returns the sorted keys of a seal file map, so a mismatch reports
// paths instead of file contents.
func sealFileNames(files map[string]string) []string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// R7: the flag is documented in the usage line and `run --progress` is not a
// usage error.
func TestRunCommandAcceptsProgressFlag(t *testing.T) {
	if !strings.Contains(runUsage, "--progress") {
		t.Fatalf("the run usage line must document --progress: %s", runUsage)
	}

	root := t.TempDir()
	cli := newTestCLI(root)
	chainPath := progressTestChainConfig(t, root)
	runDir := filepath.Join(root, "tmp", "prfrail-runs", "r7")

	code, stdout, stderr := runCLI(t, cli, "run", "--chain", chainPath, "--run-id", "run-progress-r7", "--run-dir", runDir, "--progress")
	if code != exitSuccess {
		t.Fatalf("--progress must not be a usage error: code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if progressLineTotal(stderr) == 0 {
		t.Fatalf("--progress must stream progress lines, got %q", stderr)
	}
}

// progressEventShapes reduces an event log to the transition shape of every
// record: entity kind, task, step, from->to and the reason code. Ids, timestamps
// and run-root paths are dropped on purpose, so two runs that only differ by
// --progress are comparable without depending on either.
func progressEventShapes(t *testing.T, path string) []string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the event log %s: %v", path, err)
	}
	shapes := make([]string, 0)
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event evidence.StateEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("decode the event log line %q: %v", line, err)
		}
		shapes = append(shapes, fmt.Sprintf("%s %s/%s %s->%s %s",
			event.Event.Entity.Kind,
			event.Event.Entity.TaskID,
			event.Event.Entity.StepID,
			event.Event.FromState,
			event.Event.ToState,
			event.Event.Reason.Code,
		))
	}
	return shapes
}

// progressAgentRunnerRouteFixture builds the two-task production AgentRunner
// fixture on a fixed clock and, when progress is non-nil, wires it exactly the
// way the CLI route does (AgentRunnerAssemblyOptions.Progress).
func progressAgentRunnerRouteFixture(t *testing.T, progress io.Writer) (AgentRunnerAssemblyOptions, string) {
	t.Helper()
	options, _, evidenceDirs := newAssemblyDriverFixture(t)
	dead := fmt.Sprintf("%d\n", agentRunnerDeadPID(t))
	for _, taskID := range []string{multiTaskAgentTaskOne, multiTaskAgentTaskTwo} {
		seedStopReports(t, evidenceDirs[taskID], dead, "")
	}
	options.Clock = func() time.Time { return time.Date(2026, 11, 13, 0, 0, 0, 0, time.UTC) }
	options.Progress = progress
	return options, options.RunRoot
}

// R8: the AgentRunner route is covered at the production entry the `run` command
// calls, with the offline stub port set the repository's AgentRunner e2e tests
// use. Two runs differ only by AgentRunnerAssemblyOptions.Progress -- the field
// the CLI fills from the --progress writer -- and must produce the same event
// transitions, the same projection, and progress lines that reach the sink only.
func TestProgressFlagAgentRunnerRouteKeepsTheRunIdentical(t *testing.T) {
	plainOptions, plainRoot := progressAgentRunnerRouteFixture(t, nil)
	plainSummary, err := ExecuteAgentRunnerRun(context.Background(), plainOptions)
	if err != nil {
		t.Fatalf("the plain AgentRunner run must complete: %v", err)
	}

	sink := &bytes.Buffer{}
	progressOptions, progressRoot := progressAgentRunnerRouteFixture(t, sink)
	progressSummary, err := ExecuteAgentRunnerRun(context.Background(), progressOptions)
	if err != nil {
		t.Fatalf("the AgentRunner run with progress must complete: %v", err)
	}

	if plainSummary.ChainState != "COMPLETED" || progressSummary.ChainState != "COMPLETED" {
		t.Fatalf("both runs must complete: plain=%s progress=%s", plainSummary.ChainState, progressSummary.ChainState)
	}
	if got := zeroModelTaskCount(progressSummary, "PASSED"); got != 2 {
		t.Fatalf("both tasks must still pass, got %d: %+v", got, progressSummary.TaskStateCounts)
	}
	if progressSummary.Sequence != plainSummary.Sequence {
		t.Fatalf("progress changed the event count: got %d want %d", progressSummary.Sequence, plainSummary.Sequence)
	}
	if !reflect.DeepEqual(progressSummary.TaskStateCounts, plainSummary.TaskStateCounts) ||
		!reflect.DeepEqual(progressSummary.StepStateCounts, plainSummary.StepStateCounts) {
		t.Fatalf("progress changed the projection:\n got: %+v\nwant: %+v", progressSummary, plainSummary)
	}

	// The serialised summary -- the object `run --json` emits -- must be identical
	// too. Only the two run-identity fields are normalised (see
	// progressSummaryJSON); every other byte, including the projection above, is
	// compared as-is.
	plainJSON := progressSummaryJSON(t, plainSummary, plainRoot)
	progressJSON := progressSummaryJSON(t, progressSummary, progressRoot)
	if !bytes.Equal([]byte(progressJSON), []byte(plainJSON)) {
		t.Fatalf("--progress changed the JSON summary:\n got: %s\nwant: %s", progressJSON, plainJSON)
	}
	// The exempted hash still has to be faithful to its own log: each summary must
	// report the hash of the last event that run persisted.
	if got, want := plainSummary.LastEventHash, progressLastEventHash(t, runEventLogPath(plainRoot)); got != want {
		t.Fatalf("the plain summary must report the hash of its last persisted event: got %s want %s", got, want)
	}
	if got, want := progressSummary.LastEventHash, progressLastEventHash(t, runEventLogPath(progressRoot)); got != want {
		t.Fatalf("the progress summary must report the hash of its last persisted event: got %s want %s", got, want)
	}

	plainShapes := progressEventShapes(t, runEventLogPath(plainRoot))
	progressShapes := progressEventShapes(t, runEventLogPath(progressRoot))
	if !reflect.DeepEqual(progressShapes, plainShapes) {
		t.Fatalf("progress changed the persisted transitions:\n got: %v\nwant: %v", progressShapes, plainShapes)
	}
	if len(plainShapes) == 0 {
		t.Fatal("the plain run must persist at least one event")
	}

	// One line per persisted event, and never a byte of progress inside the log.
	if total := progressLineTotal(sink.String()); total != len(plainShapes) {
		t.Fatalf("expected one progress line per event: got %d want %d\n%s", total, len(plainShapes), sink.String())
	}
	counts := progressEntityCounts(sink.String())
	for _, kind := range []string{"chain", "task", "step"} {
		if counts[kind] == 0 {
			t.Fatalf("expected at least one %q progress line, got %+v\n%s", kind, counts, sink.String())
		}
	}
	logged, err := os.ReadFile(runEventLogPath(progressRoot))
	if err != nil {
		t.Fatalf("read the progress run event log: %v", err)
	}
	if strings.Contains(string(logged), "progress: ") {
		t.Fatalf("progress must never reach the event log: %s", logged)
	}
}

// progressSupervisedChainConfig writes the supervised-black-box fixture chain:
// one code step in supervised mode, the ai block with the visible channel the
// mode requires, and a handoff policy that validates.
func progressSupervisedChainConfig(t *testing.T, root string) string {
	t.Helper()
	cfg := configuredAIChain()
	cfg.Chain.ID = "progress-supervised"
	cfg.AI.Channels = append(cfg.AI.Channels, AIChannelBinding{Channel: supervisedRequiredChannel, ProfileID: "github"})
	cfg.Tasks = []TaskConfig{{
		ID:     "task-supervised-progress",
		Review: "manual",
		Steps: []StepConfig{{
			ID:        "code-supervised-progress",
			Kind:      "code",
			Execution: supervisedExecutionMode,
			HandoffPolicy: &HandoffPolicyInput{
				AllowedTargets:   []string{"target-one"},
				InputPolicy:      "structured",
				HandoffTimeoutMs: 60000,
				ReturnActions:    []string{"complete", "abort"},
				HooksAfterReturn: []string{"hook-one"},
			},
		}},
	}}
	if err := ValidateChainConfig(cfg); err != nil {
		t.Fatalf("the supervised progress fixture chain must validate: %v", err)
	}
	content, err := EncodeChainConfig(cfg)
	if err != nil {
		t.Fatalf("encode the supervised progress fixture chain: %v", err)
	}
	path := filepath.Join(root, DefaultChainConfigName)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write the supervised progress fixture chain: %v", err)
	}
	return path
}

// R9: `run` on a supervised-black-box definition routes to the supervised
// assembly, which parks the run for the operator and must exit 0 -- without
// --progress stderr stays empty, with --progress stdout stays byte-identical and
// the lines stream to stderr only.
func TestProgressFlagSupervisedRouteKeepsStdoutAndStreamsToStderr(t *testing.T) {
	root := t.TempDir()
	chainPath := progressSupervisedChainConfig(t, root)
	const runID = "run-progress-supervised"

	run := func(name string, extra ...string) (int, string, string) {
		cli := newTestCLI(root)
		// The CLI defaults to the fail-closed placeholder; the deterministic
		// scripted stub answers the visible delivery so the route can park.
		cli.visibleSender = &adapters.StubVisibleClient{Script: []string{adapters.VisibleStubAccepted}}
		runDir := filepath.Join(root, "tmp", "prfrail-runs", name)
		args := []string{
			"run",
			"--chain", chainPath,
			"--run-id", runID,
			"--run-dir", runDir,
			"--probe-permit", filepath.Join(t.TempDir(), "probe-permit.marker"),
		}
		args = append(args, extra...)
		code, stdout, stderr := runCLI(t, cli, args...)
		return code, maskRunDir(stdout, runDir), stderr
	}

	code, plainStdout, plainStderr := run("r9-plain")
	if code != exitSuccess {
		t.Fatalf("the supervised run without --progress must exit 0: code=%d stderr=%s", code, plainStderr)
	}
	if plainStderr != "" {
		t.Fatalf("without --progress stderr must stay empty, got %q", plainStderr)
	}

	code, progressStdout, progressStderr := run("r9-progress", "--progress")
	if code != exitSuccess {
		t.Fatalf("the supervised run with --progress must exit 0: code=%d stderr=%s", code, progressStderr)
	}

	if progressStdout != plainStdout {
		t.Fatalf("--progress changed stdout:\n got:\n%s\nwant:\n%s", progressStdout, plainStdout)
	}
	if strings.Contains(plainStdout, "progress: ") || strings.Contains(progressStdout, "progress: ") {
		t.Fatalf("progress must never reach stdout:\nplain:\n%s\nwith progress:\n%s", plainStdout, progressStdout)
	}

	counts := progressEntityCounts(progressStderr)
	for _, kind := range []string{"chain", "task", "step"} {
		if counts[kind] == 0 {
			t.Fatalf("expected at least one %q progress line, got %+v (stderr=%q)", kind, counts, progressStderr)
		}
	}
	if total := progressLineTotal(progressStderr); total < 3 {
		t.Fatalf("expected at least 3 progress lines on stderr, got %d (stderr=%q)", total, progressStderr)
	}

	// The same route with --json: the serialised report must not depend on the
	// flag either. Each leg needs its own run directory, so the report is compared
	// through progressRunReportJSON, which masks the run root and normalises the
	// evidence-root-derived chain hash and leaves every other byte alone.
	code, plainJSONStdout, plainJSONStderr := run("r9-json-plain", "--json")
	if code != exitSuccess {
		t.Fatalf("the supervised JSON run without --progress must exit 0: code=%d stderr=%s", code, plainJSONStderr)
	}
	if plainJSONStderr != "" {
		t.Fatalf("without --progress the JSON route must keep stderr empty, got %q", plainJSONStderr)
	}
	code, progressJSONStdout, progressJSONStderr := run("r9-json-progress", "--json", "--progress")
	if code != exitSuccess {
		t.Fatalf("the supervised JSON run with --progress must exit 0: code=%d stderr=%s", code, progressJSONStderr)
	}
	if !json.Valid([]byte(plainJSONStdout)) || !json.Valid([]byte(progressJSONStdout)) {
		t.Fatalf("both JSON reports must stay valid JSON:\nplain:\n%s\nwith progress:\n%s", plainJSONStdout, progressJSONStdout)
	}
	if strings.Contains(plainJSONStdout, "progress: ") || strings.Contains(progressJSONStdout, "progress: ") {
		t.Fatalf("progress must never reach the JSON report:\nplain:\n%s\nwith progress:\n%s", plainJSONStdout, progressJSONStdout)
	}
	if got, want := progressRunReportJSON(t, progressJSONStdout, filepath.Join(root, "tmp", "prfrail-runs", "r9-json-progress")),
		progressRunReportJSON(t, plainJSONStdout, filepath.Join(root, "tmp", "prfrail-runs", "r9-json-plain")); got != want {
		t.Fatalf("--progress changed the JSON report:\n got: %s\nwant: %s", got, want)
	}
	if total := progressLineTotal(progressJSONStderr); total < 3 {
		t.Fatalf("expected at least 3 progress lines on stderr with --json, got %d (stderr=%q)", total, progressJSONStderr)
	}
}
