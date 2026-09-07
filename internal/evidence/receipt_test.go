package evidence

import (
	"errors"
	"testing"
)

func validManifestBody() EvidenceManifestBody {
	return EvidenceManifestBody{
		EvidenceID: "evidence-one", CreatedAt: "2026-09-07T00:00:00.000Z", RunID: "run-one", TaskID: "task-one", Attempt: 1,
		TaskDefinitionHash: Digest("", []byte("task")), PolicyHash: Digest("", []byte("policy")),
		ParentSnapshotHash: Digest("", []byte("parent")), CandidateSnapshotHash: Digest("", []byte("candidate")),
		Items: []EvidenceItem{
			{ID: "artifact-one", Kind: "artifact", ObjectHash: Digest("", []byte("artifact")), MediaType: "application/octet-stream", Redaction: "not-required"},
			{ID: "event-one", Kind: "state-event", ObjectHash: Digest("", []byte("event")), MediaType: "application/json", Redaction: "passed"},
		},
	}
}

func TestEvidenceManifestRoundTrip(t *testing.T) {
	record, err := NewEvidenceManifest(validManifestBody())
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyEvidenceManifest(record); err != nil {
		t.Fatal(err)
	}
	record.Manifest.TaskID = "other-task"
	if err := VerifyEvidenceManifest(record); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("error = %v, want invalid record", err)
	}
}

func TestEvidenceManifestRejectsSemanticInvalidity(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*EvidenceManifestBody)
	}{
		{"unsorted", func(body *EvidenceManifestBody) { body.Items[0], body.Items[1] = body.Items[1], body.Items[0] }},
		{"duplicate ID", func(body *EvidenceManifestBody) { body.Items[1].ID = body.Items[0].ID }},
		{"same snapshot", func(body *EvidenceManifestBody) { body.CandidateSnapshotHash = body.ParentSnapshotHash }},
		{"unknown kind", func(body *EvidenceManifestBody) { body.Items[0].Kind = "review-receipt" }},
		{"bad media type", func(body *EvidenceManifestBody) { body.Items[0].MediaType = "Text/Plain" }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := validManifestBody()
			test.mutate(&body)
			if _, err := NewEvidenceManifest(body); !errors.Is(err, ErrInvalidRecord) {
				t.Fatalf("error = %v, want invalid record", err)
			}
		})
	}
}
