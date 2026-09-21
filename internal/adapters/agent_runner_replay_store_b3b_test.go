package adapters

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// b3bTempResidue reports the temp files a publish left behind in a directory.
func b3bTempResidue(t *testing.T, directory string) []string {
	t.Helper()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read %s: %v", directory, err)
	}
	var residue []string
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tmp") {
			residue = append(residue, entry.Name())
		}
	}
	return residue
}

// TestB3bPublishPrimitiveHookIsNilByDefaultAndSeesEveryPublish pins the seam the
// B3b durability candidates are installed through: inert in production, and one
// call per no-replace publication handing over the temp file the store wrote and
// the record path it publishes. A seam that fires twice, or that hands over
// anything other than the record path, would let a candidate experiment measure a
// publication that never happened.
func TestB3bPublishPrimitiveHookIsNilByDefaultAndSeesEveryPublish(t *testing.T) {
	if replayStorePublishPrimitiveHook != nil {
		t.Fatal("the publish primitive hook must be nil outside experiments")
	}
	store := replayStoreMustNew(t, t.TempDir())
	path, err := store.requestPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	type call struct{ temp, published string }
	var calls []call
	original := replayStorePublishPrimitiveHook
	replayStorePublishPrimitiveHook = func(temp string, published string) error {
		calls = append(calls, call{temp: temp, published: published})
		return os.Link(temp, published)
	}
	t.Cleanup(func() { replayStorePublishPrimitiveHook = original })

	if _, err := store.RecordRequest(replayStoreRequestRecord(t, "request-one")); err != nil {
		t.Fatal(err)
	}
	replayStorePublishPrimitiveHook = original

	if len(calls) != 1 {
		t.Fatalf("the primitive hook must fire once per publication, got %d calls: %+v", len(calls), calls)
	}
	if calls[0].published != path {
		t.Fatalf("primitive target = %s, want the request record path %s", calls[0].published, path)
	}
	if filepath.Dir(calls[0].temp) != filepath.Dir(path) || !strings.HasSuffix(calls[0].temp, ".tmp") {
		t.Fatalf("primitive temp = %s, want a .tmp sibling of %s", calls[0].temp, path)
	}
	if residue := b3bTempResidue(t, filepath.Dir(path)); len(residue) != 0 {
		t.Fatalf("a successful publish must leave no temp residue, got %v", residue)
	}
	// With the hook restored to nil the store must publish exactly as before.
	if _, err := store.RecordRequest(replayStoreRequestRecord(t, "request-two")); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 1 {
		t.Fatalf("a nil hook must stay inert, calls = %+v", calls)
	}
}

// TestB3bPublishPrimitiveHookPreservesErrExistMapping pins the two error paths a
// candidate primitive has to keep distinct: an existing target is a no-replace
// collision the store is built to converge on, while every other failure must
// stay a failure. Reporting a generic failure as a collision would let the store
// converge on a slot that was never published.
func TestB3bPublishPrimitiveHookPreservesErrExistMapping(t *testing.T) {
	directory := t.TempDir()
	record := replayStoreRequestRecord(t, "request-one")
	published := filepath.Join(directory, "record-one.jsonl")
	if err := writeReplayRecordNoReplace(published, record); err != nil {
		t.Fatal(err)
	}
	if residue := b3bTempResidue(t, directory); len(residue) != 0 {
		t.Fatalf("a clean publish left temp residue %v", residue)
	}

	original := replayStorePublishPrimitiveHook
	t.Cleanup(func() { replayStorePublishPrimitiveHook = original })

	collided := filepath.Join(directory, "record-two.jsonl")
	replayStorePublishPrimitiveHook = func(string, string) error { return fs.ErrExist }
	err := writeReplayRecordNoReplace(collided, record)
	if !errors.Is(err, fs.ErrExist) {
		t.Fatalf("a no-replace collision must surface as fs.ErrExist, got %v", err)
	}
	if _, statErr := os.Stat(collided); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("a failed publish must not create the record, stat = %v", statErr)
	}
	if residue := b3bTempResidue(t, directory); len(residue) != 0 {
		t.Fatalf("a failed publish must clean up its temp file, residue %v", residue)
	}

	sentinel := errors.New("b3b primitive failure")
	failed := filepath.Join(directory, "record-three.jsonl")
	replayStorePublishPrimitiveHook = func(string, string) error { return sentinel }
	err = writeReplayRecordNoReplace(failed, record)
	if errors.Is(err, fs.ErrExist) {
		t.Fatalf("a failing primitive must not be reported as a no-replace collision: %v", err)
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("the primitive failure must reach the caller, got %v", err)
	}
	replayStorePublishPrimitiveHook = original
}

// TestB3bPublishPrimitiveHookEquivalenceWithLink drives a full store publication
// through a hook that performs exactly what the production primitive does. The
// record must be readable afterwards, and a second publication must still be
// treated as the same owned slot rather than replacing it - the equivalence a
// candidate experiment relies on when it swaps the primitive.
func TestB3bPublishPrimitiveHookEquivalenceWithLink(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	path, err := store.requestPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	original := replayStorePublishPrimitiveHook
	replayStorePublishPrimitiveHook = func(temp string, published string) error {
		calls++
		return os.Link(temp, published)
	}
	t.Cleanup(func() { replayStorePublishPrimitiveHook = original })

	decision, err := store.RecordRequest(replayStoreRequestRecord(t, "request-one"))
	if err != nil {
		t.Fatal(err)
	}
	if decision != AgentRunnerReplayDecisionFirstDispatch {
		t.Fatalf("decision = %s, want the first dispatch", decision)
	}
	published, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	decision, err = store.RecordRequest(replayStoreRequestRecord(t, "request-one"))
	if err != nil {
		t.Fatalf("republishing the same request through the hook must converge: %v", err)
	}
	if decision != AgentRunnerReplayDecisionUnknownBlock {
		t.Fatalf("decision = %s, want the owned slot to block a second dispatch", decision)
	}
	if calls != 2 {
		t.Fatalf("both publications must run through the hook, calls = %d", calls)
	}
	again, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(published, again) {
		t.Fatal("a republished record must not replace the owned slot")
	}
	if residue := b3bTempResidue(t, filepath.Dir(path)); len(residue) != 0 {
		t.Fatalf("convergence must leave no temp residue, got %v", residue)
	}
}

// TestB3bParentSyncHookIsNilByDefaultAndFiresOnceAfterPrimitive pins the second
// B3b seam: inert in production, fired once per publication, handed the directory
// holding the record, and fired between the record-linked and parent-synced
// stages. A hook that runs before the primitive or after the last stage would make
// a candidate experiment measure a durability step that ran at the wrong point.
func TestB3bParentSyncHookIsNilByDefaultAndFiresOnceAfterPrimitive(t *testing.T) {
	if replayStoreParentSyncHook != nil {
		t.Fatal("the parent sync hook must be nil outside experiments")
	}
	store := replayStoreMustNew(t, t.TempDir())
	path, err := store.requestPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	var events []string
	var directories []string
	originalStage := replayStorePublishStageHook
	originalPrimitive := replayStorePublishPrimitiveHook
	originalParent := replayStoreParentSyncHook
	replayStorePublishStageHook = func(stage string, _ string) { events = append(events, "stage:"+stage) }
	replayStorePublishPrimitiveHook = func(temp string, published string) error {
		events = append(events, "primitive")
		return os.Link(temp, published)
	}
	replayStoreParentSyncHook = func(directory string) error {
		events = append(events, "parent")
		directories = append(directories, directory)
		return nil
	}
	t.Cleanup(func() {
		replayStorePublishStageHook = originalStage
		replayStorePublishPrimitiveHook = originalPrimitive
		replayStoreParentSyncHook = originalParent
	})

	if _, err := store.RecordRequest(replayStoreRequestRecord(t, "request-one")); err != nil {
		t.Fatal(err)
	}
	replayStorePublishStageHook = originalStage
	replayStorePublishPrimitiveHook = originalPrimitive
	replayStoreParentSyncHook = originalParent

	want := strings.Join([]string{
		"stage:" + replayStorePublishStageTempWritten,
		"stage:" + replayStorePublishStageTempSynced,
		"stage:" + replayStorePublishStageTempClosed,
		"primitive",
		"stage:" + replayStorePublishStageRecordLinked,
		"parent",
		"stage:" + replayStorePublishStageParentSynced,
	}, ",")
	if got := strings.Join(events, ","); got != want {
		t.Fatalf("publication event order = %s, want %s", got, want)
	}
	if len(directories) != 1 || directories[0] != filepath.Dir(path) {
		t.Fatalf("the parent sync must receive the record directory once, got %v want %s", directories, filepath.Dir(path))
	}
	// With the hooks restored to nil a later publish must behave exactly as before.
	before := len(events)
	if _, err := store.RecordRequest(replayStoreRequestRecord(t, "request-two")); err != nil {
		t.Fatal(err)
	}
	if len(events) != before {
		t.Fatalf("nil hooks must stay inert, events = %v", events)
	}
}

// TestB3bParentSyncHookFailureAbortsPublishAndCleansTemp pins what a failed
// parent-directory durability step has to do to the caller: report the failure
// instead of swallowing it, drop the temp file, and never leave a partial record.
// Swallowing the error would let a publication that never reached the durability
// boundary look like a completed one.
func TestB3bParentSyncHookFailureAbortsPublishAndCleansTemp(t *testing.T) {
	store := replayStoreMustNew(t, t.TempDir())
	path, err := store.requestPath("request-one")
	if err != nil {
		t.Fatal(err)
	}
	sentinel := errors.New("b3b parent sync failure")
	original := replayStoreParentSyncHook
	replayStoreParentSyncHook = func(string) error { return sentinel }
	t.Cleanup(func() { replayStoreParentSyncHook = original })

	if _, err := store.RecordRequest(replayStoreRequestRecord(t, "request-one")); !errors.Is(err, sentinel) {
		t.Fatalf("a failed parent-directory step must reach the decision layer, got %v", err)
	}
	if residue := b3bTempResidue(t, filepath.Dir(path)); len(residue) != 0 {
		t.Fatalf("a failed publication must remove its temp file, residue %v", residue)
	}
	// Visibility and durability are separate facts: the primitive already
	// published the record, so what must hold afterwards is the completeness
	// invariant, not invisibility.
	if _, statErr := os.Stat(path); statErr == nil {
		loaded, found, loadErr := replayStoreLoadRequestRecord(path, "request-one")
		if loadErr != nil || !found {
			t.Fatalf("a visible record must be complete and loadable, found=%v err=%v", found, loadErr)
		}
		if loaded.RecordHash == "" {
			t.Fatal("a visible record must carry its record hash")
		}
	} else if !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("an unexpected publication outcome: %v", statErr)
	}
	replayStoreParentSyncHook = original
}

// TestB3bParentSyncHookNilEquivalence pins the equivalence the candidate
// experiments rely on: while the hook is nil the dispatcher has to be the
// platform primitive itself, and a passthrough hook has to leave the same
// publication footprint as no hook at all.
func TestB3bParentSyncHookNilEquivalence(t *testing.T) {
	if replayStoreParentSyncHook != nil {
		t.Fatal("the parent sync hook must be nil outside experiments")
	}
	directory := t.TempDir()
	missing := filepath.Join(directory, "missing")
	for _, candidate := range []string{directory, missing} {
		wantErr := replayStoreSyncParentDirectory(candidate)
		gotErr := replayStoreSyncParent(candidate)
		if (wantErr == nil) != (gotErr == nil) {
			t.Fatalf("nil hook delegation mismatch for %s: dispatched=%v platform=%v", candidate, gotErr, wantErr)
		}
		if wantErr != nil && gotErr.Error() != wantErr.Error() {
			t.Fatalf("nil hook delegation changed the error for %s: dispatched=%v platform=%v", candidate, gotErr, wantErr)
		}
	}

	footprint := func(root string) string {
		var collected []string
		walkErr := filepath.WalkDir(root, func(current string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !entry.IsDir() {
				collected = append(collected, current)
			}
			return nil
		})
		if walkErr != nil {
			t.Fatal(walkErr)
		}
		sort.Strings(collected)
		return strings.Join(collected, ",")
	}

	plainRoot := t.TempDir()
	plain := replayStoreMustNew(t, plainRoot)
	if _, err := plain.RecordRequest(replayStoreRequestRecord(t, "request-one")); err != nil {
		t.Fatal(err)
	}
	passthroughRoot := t.TempDir()
	passthrough := replayStoreMustNew(t, passthroughRoot)
	original := replayStoreParentSyncHook
	replayStoreParentSyncHook = func(directory string) error { return replayStoreSyncParentDirectory(directory) }
	if _, err := passthrough.RecordRequest(replayStoreRequestRecord(t, "request-one")); err != nil {
		t.Fatal(err)
	}
	replayStoreParentSyncHook = original

	plainNames := strings.ReplaceAll(footprint(plainRoot), plainRoot, "<root>")
	passthroughNames := strings.ReplaceAll(footprint(passthroughRoot), passthroughRoot, "<root>")
	if plainNames != passthroughNames {
		t.Fatalf("passthrough footprint = %s, want %s", passthroughNames, plainNames)
	}
}
