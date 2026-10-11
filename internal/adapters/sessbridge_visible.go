package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// This file is the visible-delivery edge of the supervised mode. It
// is deliberately receipt-free: a visible delivery proves only that a notice was
// handed to the operator UI, never that any work happened, so it produces no
// record of its own. The delivery fact is carried into the supervised risk
// confirmation record by the caller and nowhere else, which is what keeps a
// visible acknowledgement from ever masquerading as a completed step (contracts
// section 7.1a/7.1b). Nothing here references a chain transition.

var (
	// ErrVisibleClientUnavailable reports an unassembled visible client.
	ErrVisibleClientUnavailable = errors.New("sessbridge visible client unavailable")
	// ErrVisibleDeliveryRejected reports a visible dispatch the transport refused.
	ErrVisibleDeliveryRejected = errors.New("visible delivery rejected")
	// ErrVisibleDeliveryUncertain reports a visible dispatch whose outcome could
	// not be established (transport error, unknown status or request id drift).
	ErrVisibleDeliveryUncertain = errors.New("visible delivery uncertain")
)

// VisibleDeliveryUnavailableReason is the frozen ai-availability failure code a
// failed visible delivery reports (internal/adapters/ai_availability.go,
// validAIAvailabilityReason). Reusing the existing code keeps the availability
// vocabulary single-sourced.
const VisibleDeliveryUnavailableReason = "visible_delivery_unavailable"

// VisibleDelivery status values. They are the mapped dispatch outcomes of the
// shared caliber, not new transport statuses.
const (
	VisibleDeliveryAccepted  = "accepted"
	VisibleDeliveryRejected  = "rejected"
	VisibleDeliveryUncertain = "uncertain"
)

// VisibleClient sends one visible request and reports the transport answer. It
// mirrors SilentClient but owns the operator-visible channel.
type VisibleClient interface {
	SendVisible(context.Context, RequestEnvelope) (VisibleResponse, error)
}

// VisibleResponse is the visible transport answer.
type VisibleResponse struct {
	Status             string
	RequestID          string
	TransportRequestID string
	OutputEvidence     []string
	ErrorEvidence      []string
}

// VisibleDelivery is the factual result of one visible delivery attempt: the
// mapped status, the transport request id when the transport supplied one, the
// frozen failure reason when the delivery did not succeed, and the error
// evidence digests. It carries no completion claim.
type VisibleDelivery struct {
	DeliveryStatus     string
	TransportRequestID *string
	FailureReason      string
	ErrorEvidence      []string
}

// DeliverVisible maps one visible dispatch onto the delivery fact. The status
// mapping reuses the silent dispatch caliber (mapDispatchOutcome), a request id
// drift degrades an accepted answer to uncertain, and a non-accepted delivery
// always reports the frozen visible_delivery_unavailable failure code.
func DeliverVisible(ctx context.Context, client VisibleClient, envelope RequestEnvelope) (VisibleDelivery, error) {
	if client == nil {
		return VisibleDelivery{}, ErrVisibleClientUnavailable
	}
	response, callErr := client.SendVisible(ctx, envelope)
	status := mapDispatchOutcome(response.Status, callErr)
	transportID := strings.TrimSpace(response.TransportRequestID)
	errorHashes := sanitizeHashes(response.ErrorEvidence)
	transportError := ""
	if callErr != nil {
		transportError = callErr.Error()
	}
	if response.RequestID != "" && response.RequestID != envelope.Message.RequestID {
		transportError = fmt.Sprintf("requestId mismatch: expected %s got %s", envelope.Message.RequestID, response.RequestID)
		if status == VisibleDeliveryAccepted {
			status = VisibleDeliveryUncertain
		}
	}
	delivery := VisibleDelivery{DeliveryStatus: status, ErrorEvidence: errorHashes}
	if status == VisibleDeliveryAccepted {
		delivery.ErrorEvidence = []string{}
		if transportID != "" {
			delivery.TransportRequestID = &transportID
		}
		return delivery, nil
	}
	delivery.FailureReason = VisibleDeliveryUnavailableReason
	if len(delivery.ErrorEvidence) == 0 {
		delivery.ErrorEvidence = []string{visibleDeliveryErrorHash(response.Status, transportError, envelope.RecordHash)}
	}
	if status == VisibleDeliveryRejected {
		return delivery, fmt.Errorf("%w: %s: %s", ErrVisibleDeliveryRejected, VisibleDeliveryUnavailableReason, transportError)
	}
	return delivery, fmt.Errorf("%w: %s: %s", ErrVisibleDeliveryUncertain, VisibleDeliveryUnavailableReason, transportError)
}

// visibleDeliveryErrorHash binds the raw transport status and detail into a
// digest, mirroring transportErrorHash of the silent edge.
func visibleDeliveryErrorHash(status, detail, requestHash string) string {
	input := fmt.Sprintf("status=%s\ndetail=%s\nrequestHash=%s", strings.TrimSpace(status), strings.TrimSpace(detail), requestHash)
	return evidence.Digest("proofrail:adapter-visible-delivery-error:1\n", []byte(input))
}

// StubVisibleClient is the deterministic scripted visible client: it performs no
// network call and writes no file. It replays a scripted status sequence (the
// last entry repeats once the script is exhausted) so a test can drive accept,
// unavailable and unknown outcomes without any real transport.
type StubVisibleClient struct {
	// Script is the ordered sequence of scripted outcomes, e.g.
	// []string{VisibleStubAccepted, VisibleStubUnavailable}. An empty script
	// accepts every delivery.
	Script []string
	// ScriptedRequestID, when set, is returned as the transport request id for
	// every accepted delivery.
	ScriptedRequestID string

	callCount int
}

// The scripted stub outcomes.
const (
	VisibleStubAccepted    = "accepted"
	VisibleStubRejected    = "rejected"
	VisibleStubUnavailable = "unavailable"
	VisibleStubUnknown     = "unknown"
)

// SendVisible replays the next scripted outcome. The same script always yields
// the same answers in the same order.
func (stub *StubVisibleClient) SendVisible(_ context.Context, envelope RequestEnvelope) (VisibleResponse, error) {
	scripted := VisibleStubAccepted
	if len(stub.Script) > 0 {
		index := stub.callCount
		if index >= len(stub.Script) {
			index = len(stub.Script) - 1
		}
		if stub.Script[index] != "" {
			scripted = stub.Script[index]
		}
	}
	stub.callCount++
	requestID := envelope.Message.RequestID
	response := VisibleResponse{Status: scripted, RequestID: requestID}
	switch scripted {
	case VisibleStubAccepted:
		transportID := stub.ScriptedRequestID
		if transportID == "" {
			transportID = fmt.Sprintf("visible-%d", stub.callCount)
		}
		response.Status = "accepted"
		response.TransportRequestID = transportID
		return response, nil
	case VisibleStubRejected:
		response.Status = "rejected"
		return response, nil
	case VisibleStubUnavailable:
		response.Status = "unavailable"
		return response, ErrSessbridgeClientUnavailable
	default:
		response.Status = "unknown"
		return response, nil
	}
}

// Calls reports how many deliveries the stub has answered. It exists so a test
// can assert determinism and call counts without rewriting the stub.
func (stub *StubVisibleClient) Calls() int { return stub.callCount }
