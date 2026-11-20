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

// AIReceiverBindingVerifier is the optional capability a probe exposes so the
// checker can prove the receiver binding before it resolves any credential (F6).
// A probe without it keeps the pre-spawn verification inside its executor.
type AIReceiverBindingVerifier interface {
	VerifyReceiverBinding() error
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

	// F6: the probe's receiver binding and its executable hash are verified before
	// the credential is resolved, so a swapped executable is detected while the
	// secret store is still untouched. The executor keeps its own pre-spawn
	// re-verification as the last line of defence.
	if checker == nil || checker.Probe == nil {
		return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unavailable", unavailableProbeReason(channel), 0)
	}
	// F6 ordering: the verification runs before the secret is read, never after it.
	if verifier, ok := checker.Probe.(AIReceiverBindingVerifier); ok && verifier != nil {
		if err := verifier.VerifyReceiverBinding(); err != nil {
			return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unavailable", unavailableProbeReason(channel), 0)
		}
	}

	credential := []byte(nil)
	if profile.AuthMode == "secret" {
		// The nil-receiver case returned above (line 54), so the checker is known
		// to be non-nil here; only the resolver can still be absent (F6).
		if checker.Secrets == nil {
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
