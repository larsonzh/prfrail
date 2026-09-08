package adapters

import (
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type UsageMeasurement struct {
	RequestID           string   `json:"requestId"`
	IdempotencyKey      string   `json:"idempotencyKey"`
	Status              string   `json:"status"`
	ChargedAmountMicros *int64   `json:"chargedAmountMicros"`
	ObservedCalls       *int     `json:"observedCalls"`
	ObservedTokens      *int     `json:"observedTokens"`
	ProviderEvidence    []string `json:"providerEvidence"`
}

// UsageFromResultEnvelope normalizes one durable file-queue result into a
// settlement-ready usage observation with stable idempotency.
func UsageFromResultEnvelope(result ResultEnvelope) UsageMeasurement {
	calls := 1
	evidenceHashes := []string{result.RecordHash, result.Message.RequestHash}
	evidenceHashes = append(evidenceHashes, result.Message.OutputEvidence...)
	evidenceHashes = append(evidenceHashes, result.Message.ErrorEvidence...)
	return UsageMeasurement{
		RequestID:           result.Message.RequestID,
		IdempotencyKey:      stableUsageKey("result", result.Message.RequestID, result.RecordHash),
		Status:              "unknown",
		ChargedAmountMicros: nil,
		ObservedCalls:       &calls,
		ObservedTokens:      nil,
		ProviderEvidence:    sanitizeHashes(evidenceHashes),
	}
}

// UsageFromDispatchReceipt captures sessbridge dispatch evidence so a caller
// can persist a verifiable unknown-hold settlement when provider billing data
// is unavailable.
func UsageFromDispatchReceipt(receipt DispatchReceiptRecord) UsageMeasurement {
	calls := 1
	evidenceHashes := []string{receipt.ReceiptHash, receipt.Receipt.RequestHash}
	evidenceHashes = append(evidenceHashes, receipt.Receipt.Evidence...)
	evidenceHashes = append(evidenceHashes, receipt.Receipt.ErrorEvidence...)
	return UsageMeasurement{
		RequestID:           receipt.Receipt.RequestID,
		IdempotencyKey:      stableUsageKey("dispatch", receipt.Receipt.RequestID, receipt.ReceiptHash),
		Status:              "unknown",
		ChargedAmountMicros: nil,
		ObservedCalls:       &calls,
		ObservedTokens:      nil,
		ProviderEvidence:    sanitizeHashes(evidenceHashes),
	}
}

func stableUsageKey(kind, requestID, anchor string) string {
	hash := evidence.Digest("proofrail:adapter-usage-key:1\n", []byte(kind+"\n"+requestID+"\n"+anchor))
	hex := strings.TrimPrefix(hash, "sha256:")
	if len(hex) > 20 {
		hex = hex[:20]
	}
	return kind + "-" + hex
}
