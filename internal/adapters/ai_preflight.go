package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type AIProbeResult struct {
	Status       string
	Reason       string
	RequestsUsed int
	Evidence     []string
}

type AIChannelProbe interface {
	Probe(context.Context, AIProviderProfile, []byte, string, int) (AIProbeResult, error)
}

type AIAvailabilityChecker struct {
	Clock   func() time.Time
	IDs     func() string
	Secrets AISecretResolver
	Probe   AIChannelProbe
}

func (checker *AIAvailabilityChecker) Check(ctx context.Context, profile AIProviderProfile, channel string, maximumRequests int) (AIAvailabilityRecord, error) {
	profileHash, err := AIProviderProfileHash(profile)
	if err != nil || !validAIChannel(channel) || maximumRequests < 1 {
		return AIAvailabilityRecord{}, ErrInvalidAIAvailability
	}
	now := time.Now
	if checker != nil && checker.Clock != nil {
		now = checker.Clock
	}
	idSource := func() string { return "probe-missing" }
	if checker != nil && checker.IDs != nil {
		idSource = checker.IDs
	}

	credential := []byte(nil)
	if profile.AuthMode == "secret" {
		if checker == nil || checker.Secrets == nil {
			return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unavailable", "secret_missing", 0)
		}
		credential, err = checker.Secrets.Resolve(ctx, profile.SecretRef)
		if err != nil || len(credential) == 0 {
			reason := "unknown"
			if err == nil || errors.Is(err, ErrAISecretNotFound) {
				reason = "secret_missing"
			}
			clear(credential)
			return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unavailable", reason, 0)
		}
		defer clear(credential)
	}

	if checker == nil || checker.Probe == nil {
		return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unavailable", unavailableProbeReason(channel), 0)
	}
	result, probeErr := checker.Probe.Probe(ctx, profile, credential, channel, maximumRequests)
	if probeErr != nil {
		return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unknown", "unknown", result.RequestsUsed)
	}
	if !validAIProbeResult(result, maximumRequests) {
		return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unavailable", "response_invalid", max(result.RequestsUsed, 0))
	}
	availability := AIAvailability{
		ProbeID:           idSource(),
		ProfileID:         profile.ProfileID,
		ProfileConfigHash: profileHash,
		Channel:           channel,
		ProbedAt:          now().UTC().Format(TimestampLayout),
		Status:            result.Status,
		RequestsUsed:      result.RequestsUsed,
		Evidence:          result.Evidence,
	}
	if result.Status != "available" {
		availability.Reason = &result.Reason
	}
	return NewAIAvailabilityRecord(availability)
}

func validAIProbeResult(result AIProbeResult, maximumRequests int) bool {
	if result.RequestsUsed < 0 || result.RequestsUsed > maximumRequests || len(result.Evidence) == 0 {
		return false
	}
	if result.Status == "available" {
		return result.RequestsUsed > 0 && result.Reason == ""
	}
	return (result.Status == "unavailable" || result.Status == "unknown") && validAIAvailabilityReason(result.Reason)
}

func newLocalAIAvailability(profile AIProviderProfile, profileHash, channel, probeID string, probedAt time.Time, status, reason string, requestsUsed int) (AIAvailabilityRecord, error) {
	evidenceHash := evidence.Digest("proofrail:ai-availability-local-check:1\n", []byte(profileHash+"\n"+channel+"\n"+reason+"\n"+probedAt.UTC().Format(TimestampLayout)))
	return NewAIAvailabilityRecord(AIAvailability{
		ProbeID:           probeID,
		ProfileID:         profile.ProfileID,
		ProfileConfigHash: profileHash,
		Channel:           channel,
		ProbedAt:          probedAt.UTC().Format(TimestampLayout),
		Status:            status,
		Reason:            &reason,
		RequestsUsed:      requestsUsed,
		Evidence:          []string{evidenceHash},
	})
}

func unavailableProbeReason(channel string) string {
	if channel == "agent-runner-cli" {
		return "cli_unavailable"
	}
	return "sessionbridge_unavailable"
}
