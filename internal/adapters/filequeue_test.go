package adapters

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const (
	queueHashOne   = "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	queueHashTwo   = "sha256:2222222222222222222222222222222222222222222222222222222222222222"
	queueHashThree = "sha256:3333333333333333333333333333333333333333333333333333333333333333"
	queueHashFour  = "sha256:4444444444444444444444444444444444444444444444444444444444444444"
)

func queueTime() time.Time {
	return time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
}

func queueRequest(requestID, contextHash string) RequestMessage {
	return RequestMessage{
		Type:        "request",
		RequestID:   requestID,
		RunID:       "run-one",
		TaskID:      "task-one",
		Attempt:     1,
		CreatedAt:   queueTime().Format(TimestampLayout),
		Adapter:     "local-adapter",
		ContextHash: contextHash,
	}
}

func TestFileQueueDispatchClaimResultFlow(t *testing.T) {
	queue, err := NewFileQueue(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	request, err := queue.Dispatch(queueRequest("request-one", queueHashOne))
	if err != nil {
		t.Fatal(err)
	}
	claim, err := queue.Claim("request-one", "claim-one", "consumer-one", 1, queueTime().Add(time.Minute), queueTime().Add(2*time.Minute), nil)
	if err != nil {
		t.Fatal(err)
	}
	if claim.Message.RequestHash != request.RecordHash {
		t.Fatalf("claim request hash mismatch: %s != %s", claim.Message.RequestHash, request.RecordHash)
	}
	result, err := queue.SubmitResult(ResultMessage{
		Type:           "result",
		RequestID:      "request-one",
		ClaimID:        "claim-one",
		Generation:     1,
		CompletedAt:    queueTime().Add(3 * time.Minute).Format(TimestampLayout),
		RequestHash:    request.RecordHash,
		Status:         "completed",
		OutputEvidence: []string{queueHashTwo},
		ErrorEvidence:  []string{},
	})
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := queue.LoadResult("request-one")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.RecordHash != result.RecordHash {
		t.Fatalf("result hash mismatch: %s != %s", loaded.RecordHash, result.RecordHash)
	}
	if _, err := os.Stat(queue.requestPath("request-one")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("request file should be moved to inflight, stat err=%v", err)
	}
	if _, err := os.Stat(queue.inflightRequestPath("request-one")); err != nil {
		t.Fatalf("inflight request missing: %v", err)
	}
}

func TestFileQueueDispatchIsIdempotentForSamePayload(t *testing.T) {
	queue, err := NewFileQueue(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	first, err := queue.Dispatch(queueRequest("request-one", queueHashOne))
	if err != nil {
		t.Fatal(err)
	}
	second, err := queue.Dispatch(queueRequest("request-one", queueHashOne))
	if err != nil {
		t.Fatal(err)
	}
	if first.RecordHash != second.RecordHash {
		t.Fatalf("idempotent request hash mismatch: %s != %s", first.RecordHash, second.RecordHash)
	}
}

func TestFileQueueDispatchRejectsDifferentPayloadForSameRequestID(t *testing.T) {
	queue, err := NewFileQueue(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Dispatch(queueRequest("request-one", queueHashOne)); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Dispatch(queueRequest("request-one", queueHashTwo)); !errors.Is(err, ErrRecordConflict) {
		t.Fatalf("expected record conflict, got %v", err)
	}
}

func TestFileQueueClaimTakeoverRequiresReceipt(t *testing.T) {
	queue, err := NewFileQueue(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Dispatch(queueRequest("request-one", queueHashOne)); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Claim("request-one", "claim-one", "consumer-one", 1, queueTime().Add(time.Minute), queueTime().Add(2*time.Minute), nil); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Claim("request-one", "claim-two", "consumer-two", 2, queueTime().Add(3*time.Minute), queueTime().Add(4*time.Minute), nil); !errors.Is(err, ErrTakeoverRequired) {
		t.Fatalf("expected takeover receipt requirement, got %v", err)
	}
}

func TestFileQueueSubmitResultRejectsFencingMismatch(t *testing.T) {
	queue, err := NewFileQueue(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	request, err := queue.Dispatch(queueRequest("request-one", queueHashOne))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Claim("request-one", "claim-one", "consumer-one", 1, queueTime().Add(time.Minute), queueTime().Add(2*time.Minute), nil); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Claim("request-one", "claim-one", "consumer-one", 2, queueTime().Add(3*time.Minute), queueTime().Add(4*time.Minute), nil); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.SubmitResult(ResultMessage{
		Type:           "result",
		RequestID:      "request-one",
		ClaimID:        "claim-one",
		Generation:     1,
		CompletedAt:    queueTime().Add(5 * time.Minute).Format(TimestampLayout),
		RequestHash:    request.RecordHash,
		Status:         "completed",
		OutputEvidence: []string{queueHashThree},
		ErrorEvidence:  []string{},
	}); !errors.Is(err, ErrFencingConflict) {
		t.Fatalf("expected fencing conflict, got %v", err)
	}
}

func TestFileQueueSubmitResultRejectsWrongRequestHash(t *testing.T) {
	queue, err := NewFileQueue(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Dispatch(queueRequest("request-one", queueHashOne)); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Claim("request-one", "claim-one", "consumer-one", 1, queueTime().Add(time.Minute), queueTime().Add(2*time.Minute), nil); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.SubmitResult(ResultMessage{
		Type:           "result",
		RequestID:      "request-one",
		ClaimID:        "claim-one",
		Generation:     1,
		CompletedAt:    queueTime().Add(3 * time.Minute).Format(TimestampLayout),
		RequestHash:    queueHashFour,
		Status:         "completed",
		OutputEvidence: []string{queueHashTwo},
		ErrorEvidence:  []string{},
	}); !errors.Is(err, ErrInvalidEnvelope) {
		t.Fatalf("expected invalid envelope, got %v", err)
	}
}

func TestWriteCanonicalLineNoReplaceRejectsExistingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "record.jsonl")
	first := ResultEnvelope{SchemaVersion: "1.0.0", RecordHash: queueHashOne}
	if err := writeCanonicalLineNoReplace(path, first); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	second := ResultEnvelope{SchemaVersion: "1.0.0", RecordHash: queueHashTwo}
	if err := writeCanonicalLineNoReplace(path, second); !errors.Is(err, fs.ErrExist) {
		t.Fatalf("expected fs.ErrExist, got %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(before, after) {
		t.Fatal("existing record must never be overwritten")
	}
}

func TestFileQueueClaimGeneration1IsImmutable(t *testing.T) {
	queue, err := NewFileQueue(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Dispatch(queueRequest("request-one", queueHashOne)); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Claim("request-one", "claim-one", "consumer-one", 1, queueTime().Add(time.Minute), queueTime().Add(2*time.Minute), nil); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Claim("request-one", "claim-one", "consumer-one", 1, queueTime().Add(3*time.Minute), queueTime().Add(4*time.Minute), nil); !errors.Is(err, ErrRecordConflict) {
		t.Fatalf("expected record conflict for immutable generation-1 claim, got %v", err)
	}
}

func TestFileQueueClaimGeneration1RejectsTakeoverReceiptHash(t *testing.T) {
	queue, err := NewFileQueue(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Dispatch(queueRequest("request-one", queueHashOne)); err != nil {
		t.Fatal(err)
	}
	takeoverHash := queueHashTwo
	if _, err := queue.Claim("request-one", "claim-one", "consumer-one", 1, queueTime().Add(time.Minute), queueTime().Add(2*time.Minute), &takeoverHash); !errors.Is(err, ErrInvalidEnvelope) {
		t.Fatalf("expected invalid envelope for generation-1 takeover hash, got %v", err)
	}
}

func TestFileQueueClaimRenewRejectsTakeoverReceiptHash(t *testing.T) {
	queue, err := NewFileQueue(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Dispatch(queueRequest("request-one", queueHashOne)); err != nil {
		t.Fatal(err)
	}
	if _, err := queue.Claim("request-one", "claim-one", "consumer-one", 1, queueTime().Add(time.Minute), queueTime().Add(2*time.Minute), nil); err != nil {
		t.Fatal(err)
	}
	takeoverHash := queueHashTwo
	if _, err := queue.Claim("request-one", "claim-one", "consumer-one", 2, queueTime().Add(3*time.Minute), queueTime().Add(4*time.Minute), &takeoverHash); !errors.Is(err, ErrInvalidEnvelope) {
		t.Fatalf("expected invalid envelope for same-claim takeover hash, got %v", err)
	}
}
