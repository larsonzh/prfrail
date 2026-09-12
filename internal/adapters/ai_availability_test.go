package adapters

import (
	"errors"
	"os"
	"path/filepath"
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

func TestAIAvailabilityIndexConcurrentReplay(t *testing.T) {
	record, err := NewAIAvailabilityRecord(aiAvailability())
	if err != nil {
		t.Fatal(err)
	}
	index := &AIAvailabilityIndex{}
	const workers = 24
	ready := make(chan struct{}, workers)
	start := make(chan struct{})
	results := make(chan struct {
		replayed bool
		err      error
	}, workers)
	for range workers {
		go func() {
			ready <- struct{}{}
			<-start
			replayed, err := index.Record(record)
			results <- struct {
				replayed bool
				err      error
			}{replayed: replayed, err: err}
		}()
	}
	for range workers {
		<-ready
	}
	close(start)

	inserted := 0
	replays := 0
	for range workers {
		result := <-results
		if result.err != nil {
			t.Fatalf("concurrent record failed: %v", result.err)
		}
		if result.replayed {
			replays++
		} else {
			inserted++
		}
	}
	if inserted != 1 || replays != workers-1 {
		t.Fatalf("expected one insert and %d replays, inserted=%d replays=%d", workers-1, inserted, replays)
	}
}

func TestAIAvailabilityIndexConcurrentConflict(t *testing.T) {
	firstAvailability := aiAvailability()
	secondAvailability := firstAvailability
	secondAvailability.Channel = "sessionbridge-silent"
	first, err := NewAIAvailabilityRecord(firstAvailability)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewAIAvailabilityRecord(secondAvailability)
	if err != nil {
		t.Fatal(err)
	}
	index := &AIAvailabilityIndex{}
	ready := make(chan struct{}, 2)
	start := make(chan struct{})
	results := make(chan struct {
		record   AIAvailabilityRecord
		replayed bool
		err      error
	}, 2)
	for _, record := range []AIAvailabilityRecord{first, second} {
		go func(record AIAvailabilityRecord) {
			ready <- struct{}{}
			<-start
			replayed, err := index.Record(record)
			results <- struct {
				record   AIAvailabilityRecord
				replayed bool
				err      error
			}{record: record, replayed: replayed, err: err}
		}(record)
	}
	for range 2 {
		<-ready
	}
	close(start)

	var winner AIAvailabilityRecord
	var loser AIAvailabilityRecord
	successes := 0
	conflicts := 0
	for range 2 {
		result := <-results
		switch {
		case result.err == nil:
			if result.replayed {
				t.Fatal("concurrent first record was reported as a replay")
			}
			winner = result.record
			successes++
		case errors.Is(result.err, ErrAIAvailabilityConflict):
			loser = result.record
			conflicts++
		default:
			t.Fatalf("unexpected concurrent record error: %v", result.err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("expected one success and one conflict, successes=%d conflicts=%d", successes, conflicts)
	}
	if replayed, err := index.Record(winner); err != nil || !replayed {
		t.Fatalf("winner was not replayed after insertion: replayed=%v err=%v", replayed, err)
	}
	if _, err := index.Record(loser); !errors.Is(err, ErrAIAvailabilityConflict) {
		t.Fatalf("loser did not remain in conflict: %v", err)
	}
}

func TestWriteAIAvailabilityRecordPublishesCanonicalRecordWithoutReplace(t *testing.T) {
	record, err := NewAIAvailabilityRecord(aiAvailability())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "records", "availability.json")
	if err := WriteAIAvailabilityRecord(path, record); err != nil {
		t.Fatal(err)
	}
	wire, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeAIAvailabilityRecord(wire)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.RecordHash != record.RecordHash || len(wire) == 0 || wire[len(wire)-1] != '\n' {
		t.Fatalf("record did not round trip canonically: %+v", decoded)
	}
	if err := WriteAIAvailabilityRecord(path, record); !errors.Is(err, os.ErrExist) {
		t.Fatalf("expected immutable destination, got %v", err)
	}
	wireAfter, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(wireAfter) != string(wire) {
		t.Fatal("existing availability record was modified")
	}
}
