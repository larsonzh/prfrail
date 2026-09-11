package adapters

import (
	"errors"
	"testing"
	"time"
)

func aiAvailability() AIAvailability {
	return AIAvailability{
		ProbeID:           "probe-ai-one",
		ProfileID:         "deepseek-anthropic",
		ProfileConfigHash: queueHashOne,
		Channel:           "agent-runner-cli",
		ProbedAt:          "2026-09-11T08:00:00.000Z",
		Status:            "available",
		RequestsUsed:      1,
		Evidence:          []string{queueHashTwo},
	}
}

func aiAvailabilityPolicy() AIAvailabilityPolicy {
	return AIAvailabilityPolicy{
		ProfileID:         "deepseek-anthropic",
		ProfileConfigHash: queueHashOne,
		Channel:           "agent-runner-cli",
		EvaluatedAt:       time.Date(2026, 9, 11, 8, 5, 0, 0, time.UTC),
		MaximumAge:        10 * time.Minute,
		MaximumRequests:   1,
	}
}

func TestAIAvailabilityAdmissionAcceptsBoundFreshLiveProbe(t *testing.T) {
	record, err := NewAIAvailabilityRecord(aiAvailability())
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAIAvailabilityAdmission(record, aiAvailabilityPolicy()); err != nil {
		t.Fatalf("expected available admission, got %v", err)
	}
}

func TestAIAvailabilityAdmissionDoesNotCrossChannels(t *testing.T) {
	record, err := NewAIAvailabilityRecord(aiAvailability())
	if err != nil {
		t.Fatal(err)
	}
	policy := aiAvailabilityPolicy()
	policy.Channel = "sessionbridge-silent"
	if err := ValidateAIAvailabilityAdmission(record, policy); !errors.Is(err, ErrAIUnavailable) {
		t.Fatalf("expected channel mismatch to block, got %v", err)
	}
}

func TestAIAvailabilityAdmissionRejectsStaleAndOverBudgetEvidence(t *testing.T) {
	record, err := NewAIAvailabilityRecord(aiAvailability())
	if err != nil {
		t.Fatal(err)
	}
	policy := aiAvailabilityPolicy()
	policy.EvaluatedAt = policy.EvaluatedAt.Add(6 * time.Minute)
	if err := ValidateAIAvailabilityAdmission(record, policy); !errors.Is(err, ErrAIUnavailable) {
		t.Fatalf("expected stale evidence to block, got %v", err)
	}

	availability := aiAvailability()
	availability.RequestsUsed = 2
	record, err = NewAIAvailabilityRecord(availability)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAIAvailabilityAdmission(record, aiAvailabilityPolicy()); !errors.Is(err, ErrAIUnavailable) {
		t.Fatalf("expected over-budget evidence to block, got %v", err)
	}
}

func TestAIAvailabilityUnavailableReasonBlocksAdmission(t *testing.T) {
	availability := aiAvailability()
	reason := "secret_missing"
	availability.Status = "unavailable"
	availability.Reason = &reason
	availability.RequestsUsed = 0
	record, err := NewAIAvailabilityRecord(availability)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateAIAvailabilityAdmission(record, aiAvailabilityPolicy()); !errors.Is(err, ErrAIUnavailable) {
		t.Fatalf("expected unavailable report to block, got %v", err)
	}
}

func TestAIAvailabilityRequiresLiveEvidenceForAvailable(t *testing.T) {
	for name, mutate := range map[string]func(*AIAvailability){
		"no-request":  func(availability *AIAvailability) { availability.RequestsUsed = 0 },
		"no-evidence": func(availability *AIAvailability) { availability.Evidence = nil },
		"reason": func(availability *AIAvailability) {
			reason := "unknown"
			availability.Reason = &reason
		},
	} {
		t.Run(name, func(t *testing.T) {
			availability := aiAvailability()
			mutate(&availability)
			if _, err := NewAIAvailabilityRecord(availability); !errors.Is(err, ErrInvalidAIAvailability) {
				t.Fatalf("expected invalid available report, got %v", err)
			}
		})
	}
}

func TestAIAvailabilityReplayAndConflict(t *testing.T) {
	first, err := NewAIAvailabilityRecord(aiAvailability())
	if err != nil {
		t.Fatal(err)
	}
	index, err := NewAIAvailabilityIndex([]AIAvailabilityRecord{first})
	if err != nil {
		t.Fatal(err)
	}
	if replayed, err := index.Record(first); err != nil || !replayed {
		t.Fatalf("expected replay, replayed=%v err=%v", replayed, err)
	}
	changed := aiAvailability()
	changed.Channel = "sessionbridge-silent"
	conflicting, err := NewAIAvailabilityRecord(changed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := index.Record(conflicting); !errors.Is(err, ErrAIAvailabilityConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
}
