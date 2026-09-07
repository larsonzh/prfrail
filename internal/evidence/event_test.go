package evidence

import (
	"errors"
	"testing"
)

func TestNewStateEventMatchesFrozenVector(t *testing.T) {
	event := Event{
		EventID: "event-001", RunID: "run-001", Sequence: 1, OccurredAt: "2026-09-06T12:34:56.789Z",
		Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "run-001"},
		FromState: "NONE", ToState: "CREATED", InputEvidence: []string{}, Reason: Reason{Code: "run-created"},
	}
	record, err := NewStateEvent(event)
	if err != nil {
		t.Fatal(err)
	}
	const want = "sha256:522dd307620e8aad5598045aba6a29d5bcbf60c8ff2d0945c70c4b480aa9502e"
	if record.EventHash != want {
		t.Fatalf("event hash = %s, want %s", record.EventHash, want)
	}
}

func TestVerifyEventChain(t *testing.T) {
	first, err := NewStateEvent(Event{
		EventID: "event-001", RunID: "run-001", Sequence: 1, OccurredAt: "2026-09-06T12:34:56.789Z",
		Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "run-001"},
		FromState: "NONE", ToState: "CREATED", InputEvidence: []string{}, Reason: Reason{Code: "run-created"},
	})
	if err != nil {
		t.Fatal(err)
	}
	previous := first.EventHash
	second, err := NewStateEvent(Event{
		EventID: "event-002", RunID: "run-001", Sequence: 2, OccurredAt: "2026-09-06T12:35:00.000Z",
		Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "run-001"},
		FromState: "CREATED", ToState: "BASELINED", PreviousEventHash: &previous, InputEvidence: []string{}, Reason: Reason{Code: "baseline-ready"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyEventChain([]StateEvent{first, second}); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		mutate func([]StateEvent)
	}{
		{"gap", func(records []StateEvent) { records[1].Event.Sequence = 3 }},
		{"reordered", func(records []StateEvent) { records[0], records[1] = records[1], records[0] }},
		{"corrupt", func(records []StateEvent) { records[1].EventHash = Digest("", []byte("wrong")) }},
		{"state discontinuity", func(records []StateEvent) { records[1].Event.FromState = "BASELINED" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			records := []StateEvent{first, second}
			test.mutate(records)
			if err := VerifyEventChain(records); !errors.Is(err, ErrBrokenEventChain) {
				t.Fatalf("error = %v, want broken chain", err)
			}
		})
	}
}

func TestNewStateEventRejectsInvalidTransitions(t *testing.T) {
	tests := []Event{
		{EventID: "event-001", RunID: "run-001", Sequence: 1, OccurredAt: "2026-09-06T12:34:56.789Z", Actor: Actor{Type: "agent", ID: "agent-1"}, Entity: Entity{Kind: "task", RunID: "run-001", TaskID: "task-1", Attempt: 1}, FromState: "NONE", ToState: "PENDING", InputEvidence: []string{}, Reason: Reason{Code: "created"}},
		{EventID: "event-001", RunID: "run-001", Sequence: 1, OccurredAt: "2026-09-06T12:34:56Z", Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "run-001"}, FromState: "NONE", ToState: "CREATED", InputEvidence: []string{}, Reason: Reason{Code: "created"}},
	}
	for _, event := range tests {
		if _, err := NewStateEvent(event); !errors.Is(err, ErrInvalidRecord) {
			t.Fatalf("error = %v, want invalid record", err)
		}
	}
}

func TestVerifyEventChainTracksInterleavedEntities(t *testing.T) {
	first, _ := NewStateEvent(Event{EventID: "event-001", RunID: "run-001", Sequence: 1, OccurredAt: "2026-09-06T12:00:00.000Z", Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "run-001"}, FromState: "NONE", ToState: "CREATED", InputEvidence: []string{}, Reason: Reason{Code: "created"}})
	previous := first.EventHash
	second, _ := NewStateEvent(Event{EventID: "event-002", RunID: "run-001", Sequence: 2, OccurredAt: "2026-09-06T12:00:01.000Z", Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "task", RunID: "run-001", TaskID: "task-one", Attempt: 1}, FromState: "NONE", ToState: "PENDING", PreviousEventHash: &previous, InputEvidence: []string{}, Reason: Reason{Code: "created"}})
	previous = second.EventHash
	third, _ := NewStateEvent(Event{EventID: "event-003", RunID: "run-001", Sequence: 3, OccurredAt: "2026-09-06T12:00:02.000Z", Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "run-001"}, FromState: "BASELINED", ToState: "RUNNING", PreviousEventHash: &previous, InputEvidence: []string{}, Reason: Reason{Code: "running"}})
	if err := VerifyEventChain([]StateEvent{first, second, third}); !errors.Is(err, ErrBrokenEventChain) {
		t.Fatalf("error = %v, want interleaved state discontinuity", err)
	}
}

func TestDecodeStateEventRejectsDuplicateField(t *testing.T) {
	_, err := DecodeStateEvent([]byte(`{"schemaVersion":"1.0.0","schemaVersion":"1.0.0","event":{},"eventHash":"sha256:0000000000000000000000000000000000000000000000000000000000000000"}`))
	if !errors.Is(err, ErrDuplicateKey) {
		t.Fatalf("error = %v, want duplicate key", err)
	}
}
