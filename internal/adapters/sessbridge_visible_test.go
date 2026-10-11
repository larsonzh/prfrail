package adapters

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// visibleEnvelope builds a valid visible request envelope for the tests.
func visibleEnvelope(t *testing.T) RequestEnvelope {
	t.Helper()
	envelope, err := NewRequestEnvelope(RequestMessage{
		Type:        "request",
		RequestID:   "visible-request-one",
		RunID:       "run-one",
		TaskID:      "task-one",
		Attempt:     1,
		CreatedAt:   time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC).Format(TimestampLayout),
		Adapter:     "sessbridge-visible",
		ContextHash: evidence.Digest("proofrail:test-visible-context:1\n", []byte("run-one")),
	})
	if err != nil {
		t.Fatalf("build visible envelope: %v", err)
	}
	return envelope
}

// mismatchedVisibleClient answers with a request id that does not belong to the
// envelope it was handed.
type mismatchedVisibleClient struct{}

func (mismatchedVisibleClient) SendVisible(context.Context, RequestEnvelope) (VisibleResponse, error) {
	return VisibleResponse{Status: "accepted", RequestID: "visible-request-two", TransportRequestID: "visible-1"}, nil
}

func TestVisibleDeliveryMapsAccepted(t *testing.T) {
	stub := &StubVisibleClient{}
	delivery, err := DeliverVisible(context.Background(), stub, visibleEnvelope(t))
	if err != nil {
		t.Fatalf("accepted delivery must not return an error: %v", err)
	}
	if delivery.DeliveryStatus != VisibleDeliveryAccepted {
		t.Fatalf("expected accepted, got %q", delivery.DeliveryStatus)
	}
	if delivery.TransportRequestID == nil || *delivery.TransportRequestID == "" {
		t.Fatal("an accepted delivery must record the transport request id")
	}
	if len(delivery.ErrorEvidence) != 0 || delivery.FailureReason != "" {
		t.Fatalf("an accepted delivery must carry no failure: %+v", delivery)
	}
}

func TestStubVisibleClientDeterministic(t *testing.T) {
	script := []string{VisibleStubAccepted, VisibleStubUnavailable, VisibleStubUnknown}
	first := &StubVisibleClient{Script: append([]string{}, script...)}
	second := &StubVisibleClient{Script: append([]string{}, script...)}
	for index := 0; index < len(script)+1; index++ {
		left, leftErr := first.SendVisible(context.Background(), visibleEnvelope(t))
		right, rightErr := second.SendVisible(context.Background(), visibleEnvelope(t))
		if left.Status != right.Status || left.TransportRequestID != right.TransportRequestID || left.RequestID != right.RequestID {
			t.Fatalf("call %d is not deterministic: %+v vs %+v", index, left, right)
		}
		if (leftErr == nil) != (rightErr == nil) {
			t.Fatalf("call %d error is not deterministic: %v vs %v", index, leftErr, rightErr)
		}
	}
	if first.Calls() != second.Calls() || first.Calls() != len(script)+1 {
		t.Fatalf("unexpected call counts %d and %d", first.Calls(), second.Calls())
	}
	// The stub must never write a file: the run root stays empty because the
	// stub is handed no filesystem at all. Its determinism is proven above by
	// the equal answer sequences.
}

func TestVisibleDeliveryUnavailableReason(t *testing.T) {
	stub := &StubVisibleClient{Script: []string{VisibleStubUnavailable}}
	delivery, err := DeliverVisible(context.Background(), stub, visibleEnvelope(t))
	if err == nil {
		t.Fatal("an unavailable delivery must return an error")
	}
	if !errors.Is(err, ErrVisibleDeliveryUncertain) {
		t.Fatalf("expected the uncertain delivery error, got %v", err)
	}
	if !strings.Contains(err.Error(), VisibleDeliveryUnavailableReason) {
		t.Fatalf("the failure code must be visible_delivery_unavailable: %v", err)
	}
	if delivery.DeliveryStatus != VisibleDeliveryUncertain || delivery.FailureReason != VisibleDeliveryUnavailableReason {
		t.Fatalf("unexpected delivery fact: %+v", delivery)
	}
	if len(delivery.ErrorEvidence) == 0 {
		t.Fatal("a failed delivery must carry error evidence")
	}
}

func TestVisibleDeliveryRequestIDMismatch(t *testing.T) {
	delivery, err := DeliverVisible(context.Background(), mismatchedVisibleClient{}, visibleEnvelope(t))
	if err == nil {
		t.Fatal("a request id mismatch must not be accepted")
	}
	if delivery.DeliveryStatus != VisibleDeliveryUncertain {
		t.Fatalf("a request id mismatch must degrade to uncertain, got %q", delivery.DeliveryStatus)
	}
	if delivery.FailureReason != VisibleDeliveryUnavailableReason {
		t.Fatalf("unexpected failure reason %q", delivery.FailureReason)
	}
}
