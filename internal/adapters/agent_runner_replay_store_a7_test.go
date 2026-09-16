package adapters

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// a7ReleaseLeakedTempHandles removes the temp residue the in-process crash model
// leaves behind. A real crash releases the temp file handle by dying; here the
// panic unwinds while the file is still open, so the finalizer has to close it
// before the temp directory can be removed (the a7native suite kills a child
// process instead and needs none of this).
func a7ReleaseLeakedTempHandles(t *testing.T, root string) {
	t.Helper()
	directory := filepath.Join(root, "requests")
	for attempt := 0; attempt < 50; attempt++ {
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
		entries, err := os.ReadDir(directory)
		if err != nil {
			return
		}
		pending := false
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".tmp") {
				continue
			}
			if err := os.Remove(filepath.Join(directory, entry.Name())); err != nil {
				pending = true
			}
		}
		if !pending {
			return
		}
	}
	t.Fatalf("temp residue under %s stayed undeletable, which would make the temp directory cleanup fail for reasons unrelated to this test", directory)
}

// a7StageStop models a process that stops existing at one exact publish stage.
// The in-process equivalent uses the publish stage hook, because a timing-based
// kill cannot hit a chosen stage deterministically. The real kill experiments
// live in the a7native suite.
type a7StageStop struct{ stage string }

func a7PublishAndDieAtStage(t *testing.T, store *AgentRunnerReplayStore, request AgentRunnerRequestRecord, stage string) {
	t.Helper()
	original := replayStorePublishStageHook
	replayStorePublishStageHook = func(current string, _ string) {
		if current == stage {
			panic(a7StageStop{stage: current})
		}
	}
	defer func() {
		replayStorePublishStageHook = original
		recovered := recover()
		if _, ok := recovered.(a7StageStop); !ok {
			t.Fatalf("publish did not stop at stage %s: recovered %v", stage, recovered)
		}
	}()
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatalf("publish before stage %s failed: %v", stage, err)
	}
	t.Fatalf("publish returned even though it must stop at stage %s", stage)
}

// TestA7PublishStageHookIsNilByDefaultAndSeesEveryStage pins the seam the
// falsification prototype depends on: it must be inert in production and it must
// report every publish stage in order, otherwise a crash-point experiment can
// silently stop exercising a stage it claims to cover.
func TestA7PublishStageHookIsNilByDefaultAndSeesEveryStage(t *testing.T) {
	if replayStorePublishStageHook != nil {
		t.Fatal("the publish stage hook must be nil outside experiments")
	}
	store := replayStoreMustNew(t, t.TempDir())
	var stages []string
	original := replayStorePublishStageHook
	replayStorePublishStageHook = func(stage string, _ string) { stages = append(stages, stage) }
	t.Cleanup(func() { replayStorePublishStageHook = original })

	if _, err := store.RecordRequest(replayStoreRequestRecord(t, "request-one")); err != nil {
		t.Fatal(err)
	}
	replayStorePublishStageHook = original

	want := strings.Join([]string{
		replayStorePublishStageTempWritten,
		replayStorePublishStageTempSynced,
		replayStorePublishStageTempClosed,
		replayStorePublishStageRecordLinked,
		replayStorePublishStageParentSynced,
	}, ",")
	if got := strings.Join(stages, ","); got != want {
		t.Fatalf("publish stages = %s, want %s", got, want)
	}
	// With the hook restored to nil a later publish must behave exactly as before.
	if _, err := store.RecordRequest(replayStoreRequestRecord(t, "request-two")); err != nil {
		t.Fatalf("publish with a nil hook must succeed: %v", err)
	}
	if len(stages) != 5 {
		t.Fatalf("a nil hook must stay inert, stages = %v", stages)
	}
}

// TestA7CrashAtEveryStageLeavesOnlyLegalReplayStates is the core falsification
// experiment of the visibility layer: at every publish stage a process may stop
// existing, and the replay root must then hold either no record or a complete
// one - never a partial record, and never a state outside the legal set. The two
// stages after the link must be indistinguishable, which is exactly why a kill
// experiment cannot decide which of them ran.
func TestA7CrashAtEveryStageLeavesOnlyLegalReplayStates(t *testing.T) {
	stages := []string{
		replayStorePublishStageTempWritten,
		replayStorePublishStageTempSynced,
		replayStorePublishStageTempClosed,
		replayStorePublishStageRecordLinked,
		replayStorePublishStageParentSynced,
	}
	visible := map[string]bool{}
	states := map[string]AgentRunnerReplayState{}

	for _, stage := range stages {
		stage := stage
		root := t.TempDir()
		// Registered after t.TempDir so it runs before the temp directory removal.
		t.Cleanup(func() { a7ReleaseLeakedTempHandles(t, root) })
		store := replayStoreMustNew(t, root)
		request := replayStoreRequestRecord(t, "request-one")
		a7PublishAndDieAtStage(t, store, request, stage)

		restarted := replayStoreMustNew(t, root)
		state, err := restarted.State(request.Request.RequestID)
		if err != nil {
			t.Fatalf("stage %s: replay state after a crash must reconcile, got %v", stage, err)
		}
		if state != AgentRunnerReplayStateAbsent && state != AgentRunnerReplayStateDispatchedUnknown {
			t.Fatalf("stage %s: illegal replay state %s", stage, state)
		}
		path, err := restarted.requestPath(request.Request.RequestID)
		if err != nil {
			t.Fatal(err)
		}
		_, statErr := os.Stat(path)
		visible[stage] = statErr == nil
		states[stage] = state

		wantVisible := stage == replayStorePublishStageRecordLinked || stage == replayStorePublishStageParentSynced
		if visible[stage] != wantVisible {
			t.Fatalf("stage %s: record visible = %v, want %v", stage, visible[stage], wantVisible)
		}
		if wantVisible && state != AgentRunnerReplayStateDispatchedUnknown {
			t.Fatalf("stage %s: a visible request without a receipt must be %s, got %s",
				stage, AgentRunnerReplayStateDispatchedUnknown, state)
		}
		if !wantVisible && state != AgentRunnerReplayStateAbsent {
			t.Fatalf("stage %s: an unlinked record must stay invisible, got state %s", stage, state)
		}
	}

	linked := states[replayStorePublishStageRecordLinked]
	parentSynced := states[replayStorePublishStageParentSynced]
	if linked != parentSynced {
		t.Fatalf("the durability boundary must be observationally invisible under a crash: %s vs %s", linked, parentSynced)
	}
}

// TestA7MarkerPublicationAddsNoDurabilityClassAndNoNewState falsifies the
// two-phase marker as a durability mechanism rather than assuming it: a marker is
// one more directory entry published through the same no-replace primitive, so it
// inherits the same durability class and cannot upgrade it, and it must not
// introduce a replay state the protocol does not already have.
func TestA7MarkerPublicationAddsNoDurabilityClassAndNoNewState(t *testing.T) {
	root := t.TempDir()
	store := replayStoreMustNew(t, root)
	request := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(request); err != nil {
		t.Fatal(err)
	}
	before := store.PublicationDurability()

	originalSync := replayStoreSyncParentDirectory
	syncCalls := 0
	replayStoreSyncParentDirectory = func(path string) error {
		syncCalls++
		return originalSync(path)
	}
	t.Cleanup(func() { replayStoreSyncParentDirectory = originalSync })

	marker := struct {
		Marker    string `json:"marker"`
		RequestID string `json:"requestId"`
	}{Marker: "a7-two-phase", RequestID: request.Request.RequestID}
	markerPath := filepath.Join(root, "requests", "marker."+request.Request.RequestID+".jsonl")
	if err := writeReplayRecordNoReplace(markerPath, marker); err != nil {
		t.Fatalf("the marker publication must use the same primitive: %v", err)
	}
	replayStoreSyncParentDirectory = originalSync

	if syncCalls == 0 {
		t.Fatal("the marker went through no durability step at all, so it cannot be more durable than the record")
	}
	if after := store.PublicationDurability(); after != before {
		t.Fatalf("a marker publication must not change the durability class: %s -> %s", before, after)
	}
	state, err := store.State(request.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	if state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("a visible marker without a receipt must not create a new state, got %s", state)
	}

	// The platform declaration is the only thing that can lift a durability claim,
	// and a marker cannot reach it.
	if runtime.GOOS == "windows" {
		plain, err := newAgentRunnerReplayStoreAt(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if plain.PublicationDurability() != PublishDurabilityUnproven {
			t.Fatalf("the Windows declaration must stay unproven, got %s", plain.PublicationDurability())
		}
		if _, err := plain.RecordRequest(request); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
			t.Fatalf("an unproven store must refuse the write before publishing, got %v", err)
		}
	}
}

// TestA7LeftoverTempArtifactsAreInvisibleAndDoNotBlockPublication covers the
// crash residue of the first three stages: a temp file is not a record. A loader
// that scanned the directory instead of reading exact paths would adopt it and
// invent a replay decision from residue.
func TestA7LeftoverTempArtifactsAreInvisibleAndDoNotBlockPublication(t *testing.T) {
	root := t.TempDir()
	store := replayStoreMustNew(t, root)
	published := replayStoreRequestRecord(t, "request-one")
	if _, err := store.RecordRequest(published); err != nil {
		t.Fatal(err)
	}
	publishedPath, err := store.requestPath(published.Request.RequestID)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := os.ReadFile(publishedPath)
	if err != nil {
		t.Fatal(err)
	}
	// A crash before the link leaves exactly this: a complete temp file that was
	// never linked into place.
	tempPath := filepath.Join(root, "requests", ".request.request-two.jsonl.123456.tmp")
	if err := os.WriteFile(tempPath, payload, 0o600); err != nil {
		t.Fatal(err)
	}

	state, err := store.State("request-two")
	if err != nil {
		t.Fatalf("residue must not be reported as corruption: %v", err)
	}
	if state != AgentRunnerReplayStateAbsent {
		t.Fatalf("a temp file is not a record, got state %s", state)
	}
	decision, err := store.RecordRequest(replayStoreRequestRecord(t, "request-two"))
	if err != nil {
		t.Fatalf("residue must not block a later publication: %v", err)
	}
	if decision != AgentRunnerReplayDecisionFirstDispatch {
		t.Fatalf("residue must not consume the request, decision = %s", decision)
	}
	if _, err := os.Stat(tempPath); err != nil {
		t.Fatalf("the experimental residue must be left alone by the store: %v", err)
	}
}

// TestA7ReplayIdempotenceAndConflictGuardsStillBind re-binds the arbitration
// guarantees the durability question must not weaken: an identical replay is
// idempotent, and a different digest for the same requestId is a conflict that
// writes nothing.
func TestA7ReplayIdempotenceAndConflictGuardsStillBind(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	request := replayStoreRequestRecord(t, "request-one")
	if decision, err := store.RecordRequest(request); err != nil || decision != AgentRunnerReplayDecisionFirstDispatch {
		t.Fatalf("first publication = %s, %v", decision, err)
	}
	decision, err := store.RecordRequest(request)
	if err != nil {
		t.Fatal(err)
	}
	if decision != AgentRunnerReplayDecisionUnknownBlock {
		t.Fatalf("an identical replay must not dispatch again, decision = %s", decision)
	}
	body := agentRunnerRequest("request-one")
	body.StepID = body.StepID + "-a7"
	conflicting, err := NewAgentRunnerRequestRecord(body)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.RecordRequest(conflicting); !errors.Is(err, ErrAgentRunnerRequestConflict) {
		t.Fatalf("a different digest for the same requestId must conflict, got %v", err)
	}
	state, err := store.State("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("a refused conflict must write nothing, state = %s", state)
	}
}
