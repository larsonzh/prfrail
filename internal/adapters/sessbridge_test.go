package adapters

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func dispatchTime() time.Time {
	return time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
}

type fakeSilentClient struct {
	response SilentResponse
	err      error
}

func (client fakeSilentClient) SendSilent(_ context.Context, _ RequestEnvelope) (SilentResponse, error) {
	return client.response, client.err
}

func idSource(values ...string) func() string {
	index := 0
	return func() string {
		if index >= len(values) {
			return "fallback-id"
		}
		value := values[index]
		index++
		return value
	}
}

func TestSessbridgeDispatchAccepted(t *testing.T) {
	adapter := &SessbridgeAdapter{
		AdapterID:  "local-adapter",
		RecordedBy: evidence.Actor{Type: "system", ID: "proofrail"},
		Clock:      dispatchTime,
		IDs:        idSource("receipt-one"),
		Client: fakeSilentClient{response: SilentResponse{
			Status:             "ok",
			RequestID:          "request-one",
			TransportRequestID: "sess-001",
			OutputEvidence:     []string{queueHashTwo},
		}},
	}
	envelope, receipt, err := adapter.Dispatch(context.Background(), queueRequest("request-one", queueHashOne))
	if err != nil {
		t.Fatal(err)
	}
	if receipt.Receipt.Outcome != "accepted" {
		t.Fatalf("unexpected outcome: %s", receipt.Receipt.Outcome)
	}
	if receipt.Receipt.RequestHash != envelope.RecordHash {
		t.Fatalf("request hash mismatch: %s != %s", receipt.Receipt.RequestHash, envelope.RecordHash)
	}
	if receipt.Receipt.TransportRequestID == nil || *receipt.Receipt.TransportRequestID != "sess-001" {
		t.Fatalf("unexpected transport request id: %+v", receipt.Receipt.TransportRequestID)
	}
	if len(receipt.Receipt.ErrorEvidence) != 0 {
		t.Fatalf("accepted dispatch should not contain error evidence: %+v", receipt.Receipt.ErrorEvidence)
	}
}

func TestSessbridgeDispatchRejected(t *testing.T) {
	adapter := &SessbridgeAdapter{
		AdapterID:  "local-adapter",
		RecordedBy: evidence.Actor{Type: "system", ID: "proofrail"},
		Clock:      dispatchTime,
		IDs:        idSource("receipt-one"),
		Client: fakeSilentClient{response: SilentResponse{
			Status:             "rejected",
			RequestID:          "request-one",
			TransportRequestID: "sess-002",
			ErrorEvidence:      []string{queueHashThree},
		}},
	}
	_, receipt, err := adapter.Dispatch(context.Background(), queueRequest("request-one", queueHashOne))
	if !errors.Is(err, ErrDispatchRejected) {
		t.Fatalf("expected dispatch rejected, got %v", err)
	}
	if receipt.Receipt.Outcome != "rejected" {
		t.Fatalf("unexpected outcome: %s", receipt.Receipt.Outcome)
	}
	if len(receipt.Receipt.ErrorEvidence) == 0 {
		t.Fatal("rejected dispatch must keep error evidence")
	}
}

func TestSessbridgeDispatchBusyReturnsUncertain(t *testing.T) {
	adapter := &SessbridgeAdapter{
		AdapterID:  "local-adapter",
		RecordedBy: evidence.Actor{Type: "system", ID: "proofrail"},
		Clock:      dispatchTime,
		IDs:        idSource("receipt-one"),
		Client: fakeSilentClient{response: SilentResponse{
			Status:             "busy",
			RequestID:          "request-one",
			TransportRequestID: "sess-003",
		}},
	}
	_, receipt, err := adapter.Dispatch(context.Background(), queueRequest("request-one", queueHashOne))
	if !errors.Is(err, ErrDispatchUncertain) {
		t.Fatalf("expected dispatch uncertain, got %v", err)
	}
	if receipt.Receipt.Outcome != "uncertain" {
		t.Fatalf("unexpected outcome: %s", receipt.Receipt.Outcome)
	}
	if len(receipt.Receipt.ErrorEvidence) == 0 || !evidence.ValidHash(receipt.Receipt.ErrorEvidence[0]) {
		t.Fatalf("busy dispatch must have synthetic error evidence: %+v", receipt.Receipt.ErrorEvidence)
	}
}

func TestSessbridgeDispatchRequestIDMismatchReturnsUncertain(t *testing.T) {
	adapter := &SessbridgeAdapter{
		AdapterID:  "local-adapter",
		RecordedBy: evidence.Actor{Type: "system", ID: "proofrail"},
		Clock:      dispatchTime,
		IDs:        idSource("receipt-one"),
		Client: fakeSilentClient{response: SilentResponse{
			Status:             "ok",
			RequestID:          "request-two",
			TransportRequestID: "sess-004",
			OutputEvidence:     []string{queueHashTwo},
		}},
	}
	_, receipt, err := adapter.Dispatch(context.Background(), queueRequest("request-one", queueHashOne))
	if !errors.Is(err, ErrDispatchUncertain) {
		t.Fatalf("expected uncertain on mismatched request ID, got %v", err)
	}
	if receipt.Receipt.Outcome != "uncertain" || len(receipt.Receipt.ErrorEvidence) == 0 {
		t.Fatalf("mismatch must be uncertain with error evidence: %+v", receipt.Receipt)
	}
}

func TestSessbridgeDispatchClientErrorGeneratesTransportID(t *testing.T) {
	adapter := &SessbridgeAdapter{
		AdapterID:  "local-adapter",
		RecordedBy: evidence.Actor{Type: "system", ID: "proofrail"},
		Clock:      dispatchTime,
		IDs:        idSource("transport-one", "receipt-one"),
		Client: fakeSilentClient{
			err: errors.New("offline"),
		},
	}
	_, receipt, err := adapter.Dispatch(context.Background(), queueRequest("request-one", queueHashOne))
	if !errors.Is(err, ErrDispatchUncertain) {
		t.Fatalf("expected uncertain dispatch, got %v", err)
	}
	if receipt.Receipt.TransportRequestID == nil || *receipt.Receipt.TransportRequestID != "sess-transport-one" {
		t.Fatalf("expected generated transport request ID, got %+v", receipt.Receipt.TransportRequestID)
	}
	if receipt.Receipt.Outcome != "uncertain" {
		t.Fatalf("unexpected outcome: %s", receipt.Receipt.Outcome)
	}
}

func TestSessbridgeDispatchRequiresClient(t *testing.T) {
	adapter := &SessbridgeAdapter{}
	if _, _, err := adapter.Dispatch(context.Background(), queueRequest("request-one", queueHashOne)); !errors.Is(err, ErrSessbridgeClientUnavailable) {
		t.Fatalf("expected client unavailable, got %v", err)
	}
}
