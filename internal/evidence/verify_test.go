package evidence

import (
	"context"
	"errors"
	"testing"
)

type memoryReader map[string][]byte

func (reader memoryReader) ReadObject(_ context.Context, hash string) ([]byte, error) {
	data, exists := reader[hash]
	if !exists {
		return nil, ErrObjectNotFound
	}
	return append([]byte(nil), data...), nil
}

func manifestWithArtifact(t *testing.T, data []byte, mediaType string) (EvidenceManifest, memoryReader) {
	t.Helper()
	body := validManifestBody()
	body.Items = []EvidenceItem{{ID: "artifact-one", Kind: "artifact", ObjectHash: Digest("", data), MediaType: mediaType, Redaction: "not-required"}}
	manifest, err := NewEvidenceManifest(body)
	if err != nil {
		t.Fatal(err)
	}
	return manifest, memoryReader{body.Items[0].ObjectHash: data}
}

func TestVerifyManifestObjects(t *testing.T) {
	manifest, reader := manifestWithArtifact(t, []byte(`{"a":1,"b":2}`), "application/json")
	result := VerifyManifestObjects(context.Background(), manifest, reader)
	if result.Outcome != VerificationPassed || len(result.Checks) != 1 || result.Checks[0].CheckID != "reference-closure" || len(result.Checks[0].Evidence) != 1 {
		t.Fatalf("result = %#v", result)
	}
}

func TestVerifyManifestObjectsOrdersEventsBySequence(t *testing.T) {
	first, err := NewStateEvent(Event{
		EventID: "event-001", RunID: "run-one", Sequence: 1, OccurredAt: "2026-09-07T00:00:00.000Z",
		Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "run-one"},
		FromState: "NONE", ToState: "CREATED", InputEvidence: []string{}, Reason: Reason{Code: "run-created"},
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewStateEvent(Event{
		EventID: "event-002", RunID: "run-one", Sequence: 2, OccurredAt: "2026-09-07T00:00:01.000Z",
		Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "run-one"},
		FromState: "CREATED", ToState: "BASELINED", PreviousEventHash: &first.EventHash, InputEvidence: []string{}, Reason: Reason{Code: "run-baselined"},
	})
	if err != nil {
		t.Fatal(err)
	}
	firstData, _ := EncodeCanonical(first)
	secondData, _ := EncodeCanonical(second)
	body := validManifestBody()
	body.Items = []EvidenceItem{
		{ID: "a-second", Kind: "state-event", ObjectHash: Digest("", secondData), MediaType: "application/json", Redaction: "not-required"},
		{ID: "z-first", Kind: "state-event", ObjectHash: Digest("", firstData), MediaType: "application/json", Redaction: "not-required"},
	}
	manifest, err := NewEvidenceManifest(body)
	if err != nil {
		t.Fatal(err)
	}
	reader := memoryReader{body.Items[0].ObjectHash: secondData, body.Items[1].ObjectHash: firstData}
	result := VerifyManifestObjects(context.Background(), manifest, reader)
	if result.Outcome != VerificationPassed || result.Checks[0].CheckID != "event-chain" {
		t.Fatalf("result = %#v", result)
	}
}

func TestVerifyManifestStateEventBinding(t *testing.T) {
	event, err := NewStateEvent(Event{
		EventID: "event-001", RunID: "run-one", Sequence: 1, OccurredAt: "2026-09-07T00:00:00.000Z",
		Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "run-one"},
		FromState: "NONE", ToState: "CREATED", InputEvidence: []string{}, Reason: Reason{Code: "run-created"},
	})
	if err != nil {
		t.Fatal(err)
	}
	data, err := EncodeCanonical(event)
	if err != nil {
		t.Fatal(err)
	}
	body := validManifestBody()
	body.Items = []EvidenceItem{{ID: "event-one", Kind: "state-event", ObjectHash: Digest("", data), MediaType: "application/json", Redaction: "not-required"}}
	manifest, err := NewEvidenceManifest(body)
	if err != nil {
		t.Fatal(err)
	}
	result := VerifyManifestObjects(context.Background(), manifest, memoryReader{body.Items[0].ObjectHash: data})
	if result.Outcome != VerificationPassed {
		t.Fatalf("result = %#v", result)
	}

	event, _ = NewStateEvent(Event{
		EventID: "event-001", RunID: "other-run", Sequence: 1, OccurredAt: "2026-09-07T00:00:00.000Z",
		Actor: Actor{Type: "system", ID: "proofrail"}, Entity: Entity{Kind: "chain", RunID: "other-run"},
		FromState: "NONE", ToState: "CREATED", InputEvidence: []string{}, Reason: Reason{Code: "run-created"},
	})
	data, _ = EncodeCanonical(event)
	body.Items[0].ObjectHash = Digest("", data)
	manifest, _ = NewEvidenceManifest(body)
	result = VerifyManifestObjects(context.Background(), manifest, memoryReader{body.Items[0].ObjectHash: data})
	if result.Outcome != VerificationFailed || result.Checks[2].CheckID != "attempt-binding" || result.Checks[2].ErrorCode != "integrity.reference-missing" {
		t.Fatalf("tampered event result = %#v", result)
	}
}

func TestVerifyManifestObjectsFailsClosed(t *testing.T) {
	manifest, _ := manifestWithArtifact(t, []byte(`{"a":1,"b":2}`), "application/json")
	tests := []struct {
		name   string
		reader ObjectReader
		status VerificationStatus
		code   string
	}{
		{"missing", memoryReader{}, VerificationIncomplete, "integrity.reference-missing"},
		{"tampered", memoryReader{manifest.Manifest.Items[0].ObjectHash: []byte(`{"a":2}`)}, VerificationFailed, "integrity.hash-mismatch"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := VerifyManifestObjects(context.Background(), manifest, test.reader)
			if result.Outcome != test.status || result.Checks[0].ErrorCode != test.code {
				t.Fatalf("result = %#v, want %s/%s", result, test.status, test.code)
			}
		})
	}

	noncanonical := []byte(`{ "a":1 }`)
	noncanonicalManifest, _ := manifestWithArtifact(t, noncanonical, "application/json")
	result := VerifyManifestObjects(context.Background(), noncanonicalManifest, memoryReader{Digest("", noncanonical): noncanonical})
	if result.Outcome != VerificationFailed || result.Checks[0].ErrorCode != "schema.invalid-document" {
		t.Fatalf("noncanonical result = %#v", result)
	}

	result = VerifyManifestObjects(context.Background(), manifest, nil)
	if result.Outcome != VerificationIncomplete || result.Checks[0].ErrorCode != "integrity.reference-missing" {
		t.Fatalf("nil reader result = %#v", result)
	}
}

func TestDecodeRecordsRejectsUnknownFieldsAndVersion(t *testing.T) {
	if _, err := DecodeStateEvent([]byte(`{"schemaVersion":"2.0.0","event":{},"eventHash":"sha256:0000000000000000000000000000000000000000000000000000000000000000"}`)); !errors.Is(err, ErrUnsupportedVersion) {
		t.Fatalf("error = %v, want unsupported version", err)
	}
	manifest, _ := NewEvidenceManifest(validManifestBody())
	encoded, err := EncodeCanonical(manifest)
	if err != nil {
		t.Fatal(err)
	}
	encoded = append(encoded[:len(encoded)-1], []byte(`,"unknown":true}`)...)
	if _, err := DecodeEvidenceManifest(encoded); !errors.Is(err, ErrInvalidJSON) {
		t.Fatalf("error = %v, want invalid JSON", err)
	}
}
