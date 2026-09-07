package chain

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type memoryEvents struct {
	events []evidence.StateEvent
	err    error
}

func (store *memoryEvents) Append(_ context.Context, event evidence.StateEvent) error {
	if store.err != nil {
		return store.err
	}
	store.events = append(store.events, event)
	return nil
}
func (store *memoryEvents) Load(context.Context) ([]evidence.StateEvent, error) {
	return append([]evidence.StateEvent(nil), store.events...), store.err
}

func TestDefinitionValidationFourKindsAndModes(t *testing.T) {
	valid := Definition{ID: "chain-one", Tasks: []Task{{ID: "task-one", Steps: []Step{
		{ID: "code-managed", Kind: "code", Mode: ManagedChangeSet},
		{ID: "code-isolated", Kind: "code", Mode: IsolatedWorkspace},
		{ID: "code-handoff", Kind: "code", Mode: ManualHandoff, HandoffPolicy: &HandoffPolicy{
			AllowedTargets:   []string{"target-one"},
			InputPolicy:      "structured",
			HandoffTimeoutMs: 300000,
			ReturnActions:    []string{"complete", "abort"},
			HooksAfterReturn: []string{"hook-one"},
		}},
		{ID: "build-one", Kind: "build"}, {ID: "verify-one", Kind: "verify"},
		{ID: "noop-one", Kind: "noop", Reason: "not-applicable"},
	}}}}
	if err := valid.validate(); err != nil {
		t.Fatal(err)
	}
	invalid := valid
	invalid.Tasks[0].Steps[0].Mode = "direct"
	if !errors.Is(invalid.validate(), ErrInvalidDefinition) {
		t.Fatal("expected invalid change mode")
	}
	invalidHandoff := valid
	invalidHandoff.Tasks[0].Steps[2].HandoffPolicy.ReturnActions = []string{"complete", "complete"}
	if !errors.Is(invalidHandoff.validate(), ErrInvalidDefinition) {
		t.Fatal("expected invalid handoff policy")
	}
	missingPolicy := valid
	missingPolicy.Tasks[0].Steps[2].HandoffPolicy = nil
	if !errors.Is(missingPolicy.validate(), ErrInvalidDefinition) {
		t.Fatal("expected manual-handoff step without policy rejection")
	}
	managedWithPolicy := valid
	managedWithPolicy.Tasks[0].Steps[0].HandoffPolicy = &HandoffPolicy{
		AllowedTargets:   []string{"target-one"},
		InputPolicy:      "structured",
		HandoffTimeoutMs: 1000,
		ReturnActions:    []string{"abort"},
		HooksAfterReturn: []string{"hook-one"},
	}
	if !errors.Is(managedWithPolicy.validate(), ErrInvalidDefinition) {
		t.Fatal("expected managed step with handoff policy rejection")
	}
}

func TestStateEventPersistsBeforeProjection(t *testing.T) {
	store := &memoryEvents{err: errors.New("disk full")}
	writer := stateWriter{store: store, projection: Projection{TaskStates: map[string]string{}, StepStates: map[string]string{}}, clock: fixedClock, ids: func() string { return "event-one" }}
	err := writer.append(context.Background(), evidence.Actor{Type: "system", ID: "engine"}, evidence.Entity{Kind: "chain", RunID: "run-one"}, "CREATED", nil, "run-created")
	if err == nil || writer.projection.ChainState != "" {
		t.Fatal("projection changed before durable event")
	}
	store.err = nil
	if err := writer.append(context.Background(), evidence.Actor{Type: "system", ID: "engine"}, evidence.Entity{Kind: "chain", RunID: "run-one"}, "CREATED", nil, "run-created"); err != nil {
		t.Fatal(err)
	}
	if writer.projection.ChainState != "CREATED" {
		t.Fatal("projection was not updated")
	}
}

func TestInvalidTransitionDoesNotAppend(t *testing.T) {
	store := &memoryEvents{}
	writer := stateWriter{store: store, projection: Projection{TaskStates: map[string]string{}, StepStates: map[string]string{}}, clock: fixedClock, ids: func() string { return "event-one" }}
	err := writer.append(context.Background(), evidence.Actor{Type: "system", ID: "engine"}, evidence.Entity{Kind: "chain", RunID: "run-one"}, "RUNNING", nil, "invalid-start")
	if err == nil || len(store.events) != 0 || writer.projection.ChainState != "" {
		t.Fatal("invalid transition mutated durable or projected state")
	}
}

func TestRebuildRestoresAcceptedParent(t *testing.T) {
	store := &memoryEvents{}
	sequence := 0
	writer := stateWriter{store: store, projection: Projection{TaskStates: map[string]string{}, StepStates: map[string]string{}}, clock: fixedClock, ids: func() string { sequence++; return fmt.Sprintf("event-%d", sequence) }}
	actor := evidence.Actor{Type: "system", ID: "engine"}
	if err := writer.append(context.Background(), actor, chainEntity("run-one"), "CREATED", nil, "created"); err != nil {
		t.Fatal(err)
	}
	if err := writer.append(context.Background(), actor, chainEntity("run-one"), "BASELINED", []string{baselineHash}, "baselined"); err != nil {
		t.Fatal(err)
	}
	projection, err := Rebuild(store.events)
	if err != nil {
		t.Fatal(err)
	}
	if projection.AcceptedParent != baselineHash {
		t.Fatalf("accepted parent %q", projection.AcceptedParent)
	}
}

func TestFileEventStoreRejectsCorruptTail(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.jsonl")
	store, err := NewFileEventStore(path)
	if err != nil {
		t.Fatal(err)
	}
	record, err := evidence.NewStateEvent(evidence.Event{EventID: "event-one", RunID: "run-one", Sequence: 1, OccurredAt: fixedClock().Format("2006-01-02T15:04:05.000Z"), Actor: evidence.Actor{Type: "system", ID: "engine"}, Entity: evidence.Entity{Kind: "chain", RunID: "run-one"}, FromState: "NONE", ToState: "CREATED", InputEvidence: []string{}, Reason: evidence.Reason{Code: "run-created"}})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Append(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = file.WriteString("{\"torn\"")
	_ = file.Close()
	if _, err := store.Load(context.Background()); err == nil {
		t.Fatal("expected corrupt tail rejection")
	}
}

func fixedClock() time.Time { return time.Date(2026, 9, 7, 1, 2, 3, 0, time.UTC) }
