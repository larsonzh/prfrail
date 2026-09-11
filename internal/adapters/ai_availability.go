package adapters

import (
	"errors"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const aiAvailabilityDomain = "proofrail:ai-availability:1\n"

var (
	ErrInvalidAIAvailability  = errors.New("invalid AI availability report")
	ErrAIAvailabilityConflict = errors.New("AI availability report conflict")
	ErrAIUnavailable          = errors.New("AI unavailable")
)

type AIAvailability struct {
	Kind              string   `json:"kind"`
	ProbeID           string   `json:"probeId"`
	ProfileID         string   `json:"profileId"`
	ProfileConfigHash string   `json:"profileConfigHash"`
	Channel           string   `json:"channel"`
	ProbedAt          string   `json:"probedAt"`
	Status            string   `json:"status"`
	Reason            *string  `json:"reason,omitempty"`
	RequestsUsed      int      `json:"requestsUsed"`
	Evidence          []string `json:"evidence"`
}

type AIAvailabilityRecord struct {
	SchemaVersion string         `json:"schemaVersion"`
	Availability  AIAvailability `json:"availability"`
	RecordHash    string         `json:"recordHash"`
}

type AIAvailabilityPolicy struct {
	ProfileID         string
	ProfileConfigHash string
	Channel           string
	EvaluatedAt       time.Time
	MaximumAge        time.Duration
	MaximumRequests   int
}

type AIAvailabilityIndex struct {
	hashes map[string]string
}

func NewAIAvailabilityRecord(availability AIAvailability) (AIAvailabilityRecord, error) {
	if availability.Kind == "" {
		availability.Kind = "ai-availability"
	}
	availability.Evidence = sortedClone(availability.Evidence)
	if err := validateAIAvailability(availability); err != nil {
		return AIAvailabilityRecord{}, err
	}
	hash, err := digestMessage(aiAvailabilityDomain, availability)
	if err != nil {
		return AIAvailabilityRecord{}, err
	}
	return AIAvailabilityRecord{SchemaVersion: evidence.SchemaVersion, Availability: availability, RecordHash: hash}, nil
}

func DecodeAIAvailabilityRecord(input []byte) (AIAvailabilityRecord, error) {
	var record AIAvailabilityRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AIAvailabilityRecord{}, err
	}
	if err := ValidateAIAvailabilityRecord(record); err != nil {
		return AIAvailabilityRecord{}, err
	}
	return record, nil
}

func ValidateAIAvailabilityRecord(record AIAvailabilityRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAIAvailability)
	}
	if err := validateAIAvailability(record.Availability); err != nil {
		return err
	}
	expected, err := digestMessage(aiAvailabilityDomain, record.Availability)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAIAvailability)
	}
	return nil
}

func NewAIAvailabilityIndex(records []AIAvailabilityRecord) (*AIAvailabilityIndex, error) {
	index := &AIAvailabilityIndex{hashes: make(map[string]string, len(records))}
	for _, record := range records {
		if _, err := index.Record(record); err != nil {
			return nil, err
		}
	}
	return index, nil
}

func (index *AIAvailabilityIndex) Record(record AIAvailabilityRecord) (bool, error) {
	if index == nil {
		return false, fmt.Errorf("%w: nil replay index", ErrInvalidAIAvailability)
	}
	if err := ValidateAIAvailabilityRecord(record); err != nil {
		return false, err
	}
	if index.hashes == nil {
		index.hashes = make(map[string]string)
	}
	key := record.Availability.ProbeID
	if existing, found := index.hashes[key]; found {
		if existing != record.RecordHash {
			return false, fmt.Errorf("%w: probeId %s", ErrAIAvailabilityConflict, key)
		}
		return true, nil
	}
	index.hashes[key] = record.RecordHash
	return false, nil
}

func ValidateAIAvailabilityAdmission(record AIAvailabilityRecord, policy AIAvailabilityPolicy) error {
	if err := ValidateAIAvailabilityRecord(record); err != nil {
		return err
	}
	if !evidence.ValidID(policy.ProfileID) || !evidence.ValidHash(policy.ProfileConfigHash) || !validAIChannel(policy.Channel) || policy.EvaluatedAt.IsZero() || policy.MaximumAge <= 0 || policy.MaximumRequests < 1 {
		return fmt.Errorf("%w: invalid admission policy", ErrAIUnavailable)
	}
	availability := record.Availability
	if availability.ProfileID != policy.ProfileID || availability.ProfileConfigHash != policy.ProfileConfigHash || availability.Channel != policy.Channel {
		return fmt.Errorf("%w: profile or channel mismatch", ErrAIUnavailable)
	}
	probedAt, _ := time.Parse(TimestampLayout, availability.ProbedAt)
	if probedAt.After(policy.EvaluatedAt) || policy.EvaluatedAt.Sub(probedAt) > policy.MaximumAge {
		return fmt.Errorf("%w: stale or future evidence", ErrAIUnavailable)
	}
	if availability.RequestsUsed > policy.MaximumRequests {
		return fmt.Errorf("%w: request budget exceeded", ErrAIUnavailable)
	}
	if availability.Status != "available" {
		return fmt.Errorf("%w: %s", ErrAIUnavailable, *availability.Reason)
	}
	return nil
}

func validateAIAvailability(availability AIAvailability) error {
	if availability.Kind != "ai-availability" || !evidence.ValidID(availability.ProbeID) || !evidence.ValidID(availability.ProfileID) {
		return fmt.Errorf("%w: invalid identity", ErrInvalidAIAvailability)
	}
	if !evidence.ValidHash(availability.ProfileConfigHash) || !validAIChannel(availability.Channel) {
		return fmt.Errorf("%w: invalid profile or channel", ErrInvalidAIAvailability)
	}
	if _, err := time.Parse(TimestampLayout, availability.ProbedAt); err != nil {
		return fmt.Errorf("%w: invalid timestamp", ErrInvalidAIAvailability)
	}
	if availability.RequestsUsed < 0 || !validSortedUniqueHashes(availability.Evidence) || len(availability.Evidence) == 0 {
		return fmt.Errorf("%w: invalid request count or evidence", ErrInvalidAIAvailability)
	}
	switch availability.Status {
	case "available":
		if availability.Reason != nil || availability.RequestsUsed < 1 {
			return fmt.Errorf("%w: available requires a live response without a failure reason", ErrInvalidAIAvailability)
		}
	case "unavailable", "unknown":
		if availability.Reason == nil || !validAIAvailabilityReason(*availability.Reason) {
			return fmt.Errorf("%w: unavailable or unknown requires a reason", ErrInvalidAIAvailability)
		}
	default:
		return fmt.Errorf("%w: invalid status", ErrInvalidAIAvailability)
	}
	return nil
}

func validAIChannel(channel string) bool {
	return channel == "agent-runner-cli" || channel == "sessionbridge-silent" || channel == "sessionbridge-visible"
}

func validAIAvailabilityReason(reason string) bool {
	switch reason {
	case "configuration_invalid", "secret_missing", "credential_rejected", "account_unavailable", "quota_exhausted", "network_unreachable", "tls_or_proxy_failure", "provider_unavailable", "model_unavailable", "policy_blocked", "cli_unavailable", "sessionbridge_unavailable", "lm_api_unavailable", "visible_delivery_unavailable", "response_invalid", "timeout", "unknown":
		return true
	default:
		return false
	}
}
