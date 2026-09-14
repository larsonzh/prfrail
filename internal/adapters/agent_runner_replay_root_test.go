package adapters

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func replayRootTestRunRoot(t *testing.T) string {
	t.Helper()
	runRoot := t.TempDir()
	events := filepath.Join(runRoot, "events")
	if err := os.MkdirAll(events, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(events, "state-events.jsonl"), []byte("{\"schemaVersion\":\"1.0.0\"}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return runRoot
}

func replayRootTestStore(t *testing.T, runRoot, runID string) *AgentRunnerReplayStore {
	t.Helper()
	store, err := NewAgentRunnerReplayStoreForRun(runRoot, runID)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS == "windows" {
		store.publicationDurability = PublishDurabilityProven
	}
	return store
}

func replayRootTestReadFile(t *testing.T, path string) []byte {
	t.Helper()
	wire, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return wire
}

func replayRootTestTamperMarker(t *testing.T, markerPath string, mutate func(map[string]any)) {
	t.Helper()
	wire := replayRootTestReadFile(t, markerPath)
	var marker map[string]any
	if err := json.Unmarshal(wire, &marker); err != nil {
		t.Fatal(err)
	}
	mutate(marker)
	payload, err := json.Marshal(marker)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(markerPath, payload, 0o644); err != nil {
		t.Fatal(err)
	}
}

func replayRootTestCopyRunRoot(t *testing.T, source, destination string) {
	t.Helper()
	err := filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		payload, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, payload, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAgentRunnerReplayStoreForRunDerivesStableRoot(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	first := replayRootTestStore(t, runRoot, "run-one")
	wantRoot := filepath.Join(runRoot, "agent-runner-replay")
	if first.root != wantRoot {
		t.Fatalf("derived replay root = %s, want %s", first.root, wantRoot)
	}
	if first.runID != "run-one" || first.runRoot != runRoot {
		t.Fatalf("store binding = (%s, %s), want (run-one, %s)", first.runID, first.runRoot, runRoot)
	}
	for _, name := range []string{"requests", "completions"} {
		info, err := os.Stat(filepath.Join(wantRoot, name))
		if err != nil || !info.IsDir() {
			t.Fatalf("replay directory %s unavailable: info=%v err=%v", name, info, err)
		}
	}
	markerPath := filepath.Join(wantRoot, "ownership.json")
	wire := replayRootTestReadFile(t, markerPath)
	var marker AgentRunnerReplayOwnershipRecord
	if err := evidence.DecodeStrictJSON(wire, &marker); err != nil {
		t.Fatal(err)
	}
	if marker.Kind != "agent-runner-replay-ownership" || marker.SchemaVersion != evidence.SchemaVersion {
		t.Fatalf("marker identity = (%s, %s)", marker.Kind, marker.SchemaVersion)
	}
	if marker.RunID != "run-one" || marker.RunRootPath != runRoot || !evidence.ValidHash(marker.RunRootHash) {
		t.Fatalf("marker binding = (%s, %s, %s)", marker.RunID, marker.RunRootPath, marker.RunRootHash)
	}
	second := replayRootTestStore(t, runRoot, "run-one")
	if second.root != first.root {
		t.Fatalf("reopened root = %s, want %s", second.root, first.root)
	}
	if after := replayRootTestReadFile(t, markerPath); string(after) != string(wire) {
		t.Fatalf("ownership marker mutated on reopen:\nbefore=%s\nafter=%s", wire, after)
	}
}

func TestAgentRunnerReplayStoreForRunRejectsInvalidInputs(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	fileRoot := filepath.Join(runRoot, "not-a-directory")
	if err := os.WriteFile(fileRoot, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name    string
		runRoot string
		runID   string
	}{
		{"empty run id", runRoot, ""},
		{"invalid run id", runRoot, "run one"},
		{"empty run root", "", "run-one"},
		{"relative run root", filepath.Join("relative", "run"), "run-one"},
		{"missing run root", filepath.Join(runRoot, "missing"), "run-one"},
		{"file run root", fileRoot, "run-one"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			if _, err := NewAgentRunnerReplayStoreForRun(testCase.runRoot, testCase.runID); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
				t.Fatalf("expected invalid store error, got %v", err)
			}
		})
	}
}

func TestAgentRunnerReplayStoreForRunRequiresDurableRunRootMarker(t *testing.T) {
	const runID = "run-one"
	t.Run("missing events directory", func(t *testing.T) {
		if _, err := NewAgentRunnerReplayStoreForRun(t.TempDir(), runID); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
			t.Fatalf("expected invalid store error, got %v", err)
		}
	})
	t.Run("events path is a file", func(t *testing.T) {
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "events"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := NewAgentRunnerReplayStoreForRun(root, runID); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
			t.Fatalf("expected invalid store error, got %v", err)
		}
	})
	t.Run("missing event log", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "events"), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := NewAgentRunnerReplayStoreForRun(root, runID); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
			t.Fatalf("expected invalid store error, got %v", err)
		}
	})
	t.Run("event log is a directory", func(t *testing.T) {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "events", "state-events.jsonl"), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := NewAgentRunnerReplayStoreForRun(root, runID); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
			t.Fatalf("expected invalid store error, got %v", err)
		}
	})
}

func TestAgentRunnerReplayStoreForRunRejectsForeignRunWrites(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	body := agentRunnerRequest("request-foreign-run")
	body.RunID = "run-two"
	record, err := NewAgentRunnerRequestRecord(body)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordRequest(record); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("expected replay root conflict, got %v", err)
	}
	path, err := store.requestPath(record.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("foreign-run request must not be published, stat err=%v", err)
	}
}

func TestAgentRunnerReplayStoreForRunDetectsPersistedRunIDMismatch(t *testing.T) {
	root := t.TempDir()
	store := replayStoreMustNew(t, root)
	store.runID = "run-one"
	valid := replayStoreRequestRecord(t, "request-valid")
	if _, err := store.RecordRequest(valid); err != nil {
		t.Fatal(err)
	}
	if state, err := store.State(valid.Request.RequestID); err != nil || state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("state = %s err = %v", state, err)
	}

	requestBody := agentRunnerRequest("request-foreign")
	requestBody.RunID = "run-two"
	foreignRequest, err := NewAgentRunnerRequestRecord(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	requestWire, err := evidence.EncodeCanonical(foreignRequest)
	if err != nil {
		t.Fatal(err)
	}
	requestPath := filepath.Join(root, "requests", "request."+foreignRequest.Request.RequestID+".jsonl")
	if err := os.WriteFile(requestPath, append(requestWire, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := store.State(foreignRequest.Request.RequestID); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for foreign-run request, got %v", err)
	}

	completionBody := agentRunnerCompletion("completion-foreign")
	completionBody.RequestID = valid.Request.RequestID
	completionBody.RequestHash = valid.RecordHash
	completionBody.RunID = "run-two"
	completionBody.TaskID = valid.Request.TaskID
	completionBody.StepID = valid.Request.StepID
	completionBody.Attempt = valid.Request.Attempt
	completionBody.AdapterID = valid.Request.AdapterID
	foreignCompletion, err := NewAgentRunnerCompletionRecord(completionBody)
	if err != nil {
		t.Fatal(err)
	}
	completionWire, err := evidence.EncodeCanonical(foreignCompletion)
	if err != nil {
		t.Fatal(err)
	}
	completionPath := filepath.Join(root, "completions", "completion."+valid.Request.RequestID+".jsonl")
	if err := os.MkdirAll(filepath.Dir(completionPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(completionPath, append(completionWire, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := store.State(valid.Request.RequestID); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for foreign-run completion, got %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunOwnershipPreflight(t *testing.T) {
	markerPathFor := func(runRoot string) string {
		return filepath.Join(runRoot, "agent-runner-replay", "ownership.json")
	}
	t.Run("tampered run id conflicts", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		replayRootTestStore(t, runRoot, "run-one")
		replayRootTestTamperMarker(t, markerPathFor(runRoot), func(marker map[string]any) {
			marker["runId"] = "run-two"
		})
		if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
			t.Fatalf("expected replay root conflict, got %v", err)
		}
	})
	t.Run("tampered run root hash conflicts", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		replayRootTestStore(t, runRoot, "run-one")
		replayRootTestTamperMarker(t, markerPathFor(runRoot), func(marker map[string]any) {
			marker["runRootHash"] = queueHashOne
		})
		if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
			t.Fatalf("expected replay root conflict, got %v", err)
		}
	})
	t.Run("invalid createdAt is corruption", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		replayRootTestStore(t, runRoot, "run-one")
		replayRootTestTamperMarker(t, markerPathFor(runRoot), func(marker map[string]any) {
			marker["createdAt"] = "not-a-timestamp"
		})
		if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
			t.Fatalf("expected corruption, got %v", err)
		}
	})
	t.Run("undecodable marker is corruption", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		replayRootTestStore(t, runRoot, "run-one")
		if err := os.WriteFile(markerPathFor(runRoot), []byte("not json"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
			t.Fatalf("expected corruption, got %v", err)
		}
	})
	t.Run("missing marker with published records is corruption", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		store := replayRootTestStore(t, runRoot, "run-one")
		record := replayStoreRequestRecord(t, "request-published")
		if _, err := store.RecordRequest(record); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(markerPathFor(runRoot)); err != nil {
			t.Fatal(err)
		}
		if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
			t.Fatalf("expected corruption, got %v", err)
		}
	})
	t.Run("missing marker with launch intent is corruption", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		store := replayRootTestStore(t, runRoot, "run-one")
		record := replayStoreRequestRecord(t, "request-launched")
		if _, err := store.RecordRequest(record); err != nil {
			t.Fatal(err)
		}
		if _, err := store.RecordLaunchReceipt(launchReceiptRecordFor(t, record)); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(markerPathFor(runRoot)); err != nil {
			t.Fatal(err)
		}
		if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
			t.Fatalf("expected corruption, got %v", err)
		}
	})
	t.Run("missing marker without records is rebuilt", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		replayRootTestStore(t, runRoot, "run-one")
		if err := os.Remove(markerPathFor(runRoot)); err != nil {
			t.Fatal(err)
		}
		store, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(markerPathFor(runRoot)); err != nil {
			t.Fatalf("ownership marker not rebuilt: %v", err)
		}
		if store.runID != "run-one" {
			t.Fatalf("rebuilt store runID = %s", store.runID)
		}
	})
	t.Run("copied run root conflicts", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		replayRootTestStore(t, runRoot, "run-one")
		copyRoot := filepath.Join(t.TempDir(), "copied-run")
		replayRootTestCopyRunRoot(t, runRoot, copyRoot)
		if _, err := NewAgentRunnerReplayStoreForRun(copyRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
			t.Fatalf("expected replay root conflict, got %v", err)
		}
	})
	t.Run("different run id conflicts", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		replayRootTestStore(t, runRoot, "run-one")
		if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-two"); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
			t.Fatalf("expected replay root conflict, got %v", err)
		}
	})
	t.Run("marker directory is corruption", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		replayRootTestStore(t, runRoot, "run-one")
		if err := os.Remove(markerPathFor(runRoot)); err != nil {
			t.Fatal(err)
		}
		if err := os.Mkdir(markerPathFor(runRoot), 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
			t.Fatalf("expected corruption, got %v", err)
		}
	})
	t.Run("residual temp file without ownership is corruption", func(t *testing.T) {
		runRoot := replayRootTestRunRoot(t)
		requests := filepath.Join(runRoot, "agent-runner-replay", "requests")
		if err := os.MkdirAll(requests, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(requests, ".request.x.jsonl.123.tmp"), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
			t.Fatalf("expected corruption, got %v", err)
		}
	})
}

func TestAgentRunnerReplayStoreForRunRejectsOverlappingWriterRoots(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	replayRoot := filepath.Join(runRoot, "agent-runner-replay")
	conflicts := []string{
		replayRoot,
		filepath.Join(replayRoot, "requests"),
		runRoot,
		filepath.Dir(runRoot),
		filepath.Join("relative", "workspace"),
	}
	for _, root := range conflicts {
		if err := store.RejectWriterRoots(root); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
			t.Fatalf("writer root %s: expected conflict, got %v", root, err)
		}
	}
	allowed := []string{
		filepath.Join(runRoot, "workspaces", "task-one-attempt-1"),
		t.TempDir(),
	}
	for _, root := range allowed {
		if err := store.RejectWriterRoots(root); err != nil {
			t.Fatalf("writer root %s: expected no error, got %v", root, err)
		}
	}
}

func TestAgentRunnerReplayStoreForRunRejectsProtectedRoots(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	replayRoot := filepath.Join(runRoot, "agent-runner-replay")
	conflicts := []string{
		replayRoot,
		filepath.Join(replayRoot, "requests"),
	}
	for _, root := range conflicts {
		if err := store.RejectProtectedRoots(root); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
			t.Fatalf("protected root %s: expected conflict, got %v", root, err)
		}
	}
	allowed := []string{
		runRoot,
		filepath.Dir(runRoot),
		t.TempDir(),
	}
	for _, root := range allowed {
		if err := store.RejectProtectedRoots(root); err != nil {
			t.Fatalf("protected root %s: expected no error, got %v", root, err)
		}
	}
	unbound := replayStoreMustNew(t, t.TempDir())
	if err := unbound.RejectProtectedRoots(filepath.Dir(unbound.root)); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("expected conflict without a bound run root, got %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunConvergesUnderConcurrentFirstOpen(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	const workers = 4
	stores := make([]*AgentRunnerReplayStore, workers)
	errs := make([]error, workers)
	start := make(chan struct{})
	var waitGroup sync.WaitGroup
	for index := 0; index < workers; index++ {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			<-start
			stores[index], errs[index] = NewAgentRunnerReplayStoreForRun(runRoot, "run-one")
		}(index)
	}
	close(start)
	waitGroup.Wait()
	for index := 0; index < workers; index++ {
		if errs[index] != nil {
			t.Fatalf("worker %d failed: %v", index, errs[index])
		}
		if stores[index].root != stores[0].root || stores[index].runID != "run-one" {
			t.Fatalf("worker %d binding = (%s, %s)", index, stores[index].root, stores[index].runID)
		}
	}
	markerPath := filepath.Join(runRoot, "agent-runner-replay", "ownership.json")
	if first, second := replayRootTestReadFile(t, markerPath), replayRootTestReadFile(t, markerPath); string(first) != string(second) {
		t.Fatal("ownership marker changed after concurrent first open")
	}
}

func TestAgentRunnerReplayStoreForRunKeepsPublicationDurabilitySemantics(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one")
	if err != nil {
		t.Fatal(err)
	}
	want := PublishDurabilityProven
	if runtime.GOOS == "windows" {
		want = PublishDurabilityUnproven
	}
	if got := store.PublicationDurability(); got != want {
		t.Fatalf("publication durability = %s, want %s", got, want)
	}
	if runtime.GOOS != "windows" {
		return
	}
	record := replayStoreRequestRecord(t, "request-unproven")
	decision, err := store.RecordRequest(record)
	if !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) || decision == AgentRunnerReplayDecisionFirstDispatch {
		t.Fatalf("expected durability-unproven rejection, got decision=%s err=%v", decision, err)
	}
	path, pathErr := store.requestPath(record.Request.RequestID)
	if pathErr != nil {
		t.Fatal(pathErr)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("unproven durability must not create R, stat err=%v", statErr)
	}
}

func TestAgentRunnerReplayStoreForRunRejectsReplayRootAsFile(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	if err := os.WriteFile(filepath.Join(runRoot, "agent-runner-replay"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewAgentRunnerReplayStoreForRun(runRoot, "run-one"); !errors.Is(err, ErrInvalidAgentRunnerReplayStore) {
		t.Fatalf("expected invalid store error, got %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunRejectsForeignRunCompletionWrites(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	store := replayRootTestStore(t, runRoot, "run-one")
	request := replayStoreRequestRecord(t, "request-completion-run")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	completionBody := agentRunnerCompletion("completion-foreign-run")
	completionBody.RequestID = request.Request.RequestID
	completionBody.RequestHash = request.RecordHash
	completionBody.RunID = "run-two"
	completionBody.TaskID = request.Request.TaskID
	completionBody.StepID = request.Request.StepID
	completionBody.Attempt = request.Request.Attempt
	completionBody.AdapterID = request.Request.AdapterID
	foreignCompletion, err := NewAgentRunnerCompletionRecord(completionBody)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordCompletion(request, foreignCompletion); !errors.Is(err, ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("expected replay root conflict for foreign-run completion write, got %v", err)
	}
	wire, err := evidence.EncodeCanonical(foreignCompletion)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(runRoot, "agent-runner-replay", "completions", "completion."+request.Request.RequestID+".jsonl")
	if err := os.WriteFile(path, append(wire, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordCompletion(request, replayStoreCompletionRecord(t, request, "completion-new")); !errors.Is(err, ErrAgentRunnerReplayStoreCorruption) {
		t.Fatalf("expected corruption for persisted foreign-run completion, got %v", err)
	}
}

func TestAgentRunnerReplayStoreForRunConcurrentDistinctRunsConflict(t *testing.T) {
	runRoot := replayRootTestRunRoot(t)
	const runIDs = 2
	stores := make([]*AgentRunnerReplayStore, runIDs)
	errs := make([]error, runIDs)
	runIdentifiers := []string{"run-one", "run-two"}
	start := make(chan struct{})
	var waitGroup sync.WaitGroup
	for index := 0; index < runIDs; index++ {
		waitGroup.Add(1)
		go func(index int) {
			defer waitGroup.Done()
			<-start
			stores[index], errs[index] = NewAgentRunnerReplayStoreForRun(runRoot, runIdentifiers[index])
		}(index)
	}
	close(start)
	waitGroup.Wait()
	successes := 0
	conflicts := 0
	winner := ""
	for index := 0; index < runIDs; index++ {
		switch {
		case errs[index] == nil:
			successes++
			winner = runIdentifiers[index]
		case errors.Is(errs[index], ErrAgentRunnerReplayRootConflict):
			conflicts++
		default:
			t.Fatalf("unexpected worker %d error: %v", index, errs[index])
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("expected exactly one success and one conflict, got successes=%d conflicts=%d (errs=%v)", successes, conflicts, errs)
	}
	if errs[0] == nil && stores[0].runID != winner {
		t.Fatalf("winning store is bound to %s, want %s", stores[0].runID, winner)
	}
	if errs[1] == nil && stores[1].runID != winner {
		t.Fatalf("winning store is bound to %s, want %s", stores[1].runID, winner)
	}
	markerPath := filepath.Join(runRoot, "agent-runner-replay", "ownership.json")
	var marker AgentRunnerReplayOwnershipRecord
	if err := evidence.DecodeStrictJSON(replayRootTestReadFile(t, markerPath), &marker); err != nil {
		t.Fatal(err)
	}
	if marker.RunID != winner {
		t.Fatalf("ownership marker runId = %s, want %s", marker.RunID, winner)
	}
}
