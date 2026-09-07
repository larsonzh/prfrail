package adapters

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrSessbridgeClientUnavailable = errors.New("sessbridge client unavailable")
	ErrDispatchRejected            = errors.New("adapter dispatch rejected")
	ErrDispatchUncertain           = errors.New("adapter dispatch uncertain")
)

type SilentClient interface {
	SendSilent(context.Context, RequestEnvelope) (SilentResponse, error)
}

type SilentResponse struct {
	Status             string
	RequestID          string
	TransportRequestID string
	OutputEvidence     []string
	ErrorEvidence      []string
}

type SessbridgeAdapter struct {
	AdapterID  string
	RecordedBy evidence.Actor
	Clock      func() time.Time
	IDs        func() string
	Client     SilentClient
}

func (adapter *SessbridgeAdapter) Dispatch(ctx context.Context, request RequestMessage) (RequestEnvelope, DispatchReceiptRecord, error) {
	if adapter == nil || adapter.Client == nil {
		return RequestEnvelope{}, DispatchReceiptRecord{}, ErrSessbridgeClientUnavailable
	}
	envelope, err := NewRequestEnvelope(request)
	if err != nil {
		return RequestEnvelope{}, DispatchReceiptRecord{}, err
	}
	now := time.Now
	if adapter.Clock != nil {
		now = adapter.Clock
	}
	idSource := adapter.IDs
	if idSource == nil {
		idSource = func() string { return "id-missing" }
	}
	actor := adapter.RecordedBy
	if actor.Type == "" {
		actor.Type = "system"
	}
	if actor.ID == "" {
		actor.ID = "proofrail"
	}
	adapterID := adapter.AdapterID
	if adapterID == "" {
		adapterID = envelope.Message.Adapter
	}

	response, callErr := adapter.Client.SendSilent(ctx, envelope)
	transportID := strings.TrimSpace(response.TransportRequestID)
	if transportID == "" {
		transportID = "sess-" + idSource()
	}
	outcome := mapDispatchOutcome(response.Status, callErr)
	evidenceHashes := sanitizeHashes(response.OutputEvidence)
	errorHashes := sanitizeHashes(response.ErrorEvidence)
	transportError := ""
	if callErr != nil {
		transportError = callErr.Error()
	}
	if response.RequestID != "" && response.RequestID != envelope.Message.RequestID {
		transportError = fmt.Sprintf("requestId mismatch: expected %s got %s", envelope.Message.RequestID, response.RequestID)
		if outcome == "accepted" {
			outcome = "uncertain"
		}
	}
	if outcome == "accepted" {
		errorHashes = []string{}
	} else if len(errorHashes) == 0 {
		errorHashes = []string{transportErrorHash(response.Status, transportError, envelope.RecordHash)}
	}

	receiptRecord, err := NewDispatchReceiptRecord(DispatchReceipt{
		Type:               "dispatch",
		ReceiptID:          idSource(),
		OccurredAt:         now().UTC().Format(TimestampLayout),
		RecordedBy:         actor,
		RequestID:          envelope.Message.RequestID,
		RequestHash:        envelope.RecordHash,
		Evidence:           evidenceHashes,
		ErrorEvidence:      errorHashes,
		Adapter:            adapterID,
		Transport:          "sessbridge",
		TransportRequestID: &transportID,
		Outcome:            outcome,
	})
	if err != nil {
		return RequestEnvelope{}, DispatchReceiptRecord{}, err
	}

	switch outcome {
	case "accepted":
		return envelope, receiptRecord, nil
	case "rejected":
		if callErr != nil {
			return envelope, receiptRecord, errors.Join(ErrDispatchRejected, callErr)
		}
		return envelope, receiptRecord, ErrDispatchRejected
	default:
		if callErr != nil {
			return envelope, receiptRecord, errors.Join(ErrDispatchUncertain, callErr)
		}
		if transportError != "" {
			return envelope, receiptRecord, fmt.Errorf("%w: %s", ErrDispatchUncertain, transportError)
		}
		return envelope, receiptRecord, ErrDispatchUncertain
	}
}

func mapDispatchOutcome(status string, callErr error) string {
	if callErr != nil {
		return "uncertain"
	}
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "ok", "accepted", "completed":
		return "accepted"
	case "rejected", "failed":
		return "rejected"
	case "busy", "timeout", "poll_timeout", "lm_api_unavailable", "uncertain", "error":
		return "uncertain"
	default:
		return "uncertain"
	}
}

func sanitizeHashes(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !evidence.ValidHash(value) {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func transportErrorHash(status, detail, requestHash string) string {
	input := fmt.Sprintf("status=%s\ndetail=%s\nrequestHash=%s", strings.TrimSpace(status), strings.TrimSpace(detail), requestHash)
	return evidence.Digest("proofrail:adapter-transport-error:1\n", []byte(input))
}
