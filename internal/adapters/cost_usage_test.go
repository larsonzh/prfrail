package adapters

import (
	"context"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestUsageFromResultEnvelopeIncludesVerifiableEvidence(t *testing.T) {
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

	usage := UsageFromResultEnvelope(result)
	if usage.Status != "unknown" || usage.ObservedCalls == nil || *usage.ObservedCalls != 1 {
		t.Fatalf("unexpected usage summary: %+v", usage)
	}
	if !evidence.ValidID(usage.IdempotencyKey) {
		t.Fatalf("idempotency key must be valid ID: %s", usage.IdempotencyKey)
	}
	if !containsHash(usage.ProviderEvidence, result.RecordHash) || !containsHash(usage.ProviderEvidence, request.RecordHash) || !containsHash(usage.ProviderEvidence, queueHashTwo) {
		t.Fatalf("usage evidence is incomplete: %+v", usage.ProviderEvidence)
	}
}

func TestUsageFromDispatchReceiptIncludesVerifiableEvidence(t *testing.T) {
	adapter := &SessbridgeAdapter{
		AdapterID:  "local-adapter",
		RecordedBy: evidence.Actor{Type: "system", ID: "proofrail"},
		Clock:      dispatchTime,
		IDs:        idSource("receipt-one"),
		Client: fakeSilentClient{response: SilentResponse{
			Status:             "busy",
			RequestID:          "request-one",
			TransportRequestID: "sess-003",
			ErrorEvidence:      []string{queueHashThree},
		}},
	}
	_, receipt, err := adapter.Dispatch(context.Background(), queueRequest("request-one", queueHashOne))
	if err == nil {
		t.Fatal("busy response must produce uncertain dispatch error")
	}
	usage := UsageFromDispatchReceipt(receipt)
	if usage.Status != "unknown" || usage.ObservedCalls == nil || *usage.ObservedCalls != 1 {
		t.Fatalf("unexpected usage summary: %+v", usage)
	}
	if !evidence.ValidID(usage.IdempotencyKey) {
		t.Fatalf("idempotency key must be valid ID: %s", usage.IdempotencyKey)
	}
	if !containsHash(usage.ProviderEvidence, receipt.ReceiptHash) || !containsHash(usage.ProviderEvidence, receipt.Receipt.RequestHash) {
		t.Fatalf("usage evidence is incomplete: %+v", usage.ProviderEvidence)
	}
}

func containsHash(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
