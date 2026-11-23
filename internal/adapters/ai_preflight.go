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

// AIImmutableInstallSourceVerifier is the optional capability a probe exposes so
// the checker can prove the executable's installation source before it resolves
// any credential (T3, finding F6). It runs right after the receiver binding and
// strictly before the secret store is touched, so an executable this user can
// replace or rewrite in place is refused while the credential is still unread.
// A probe without it keeps the executor's pre-spawn check as the last line of
// defence.
type AIImmutableInstallSourceVerifier interface {
	VerifyImmutableInstallSource() error
}

// AIBYOKConfinementVerifier is the optional capability a probe exposes so the
// checker can prove, before it resolves any credential, that the child which will
// receive the BYOK credential cannot write unredacted output to disk (T8a,
// finding F8). It runs after the receiver binding and the install source and
// strictly before the secret store is touched: an unprovable confinement refuses
// the probe while the credential is still unread, and no child is created.
//
// A probe without it is refused as unproven instead of admitted: the gate is
// fail-closed by design, so a missing capability is a refusal and never a
// downgrade to "assume confined".
type AIBYOKConfinementVerifier interface {
	VerifyBYOKConfinement() error
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
	// T3 ordering (CONTRACTS.md:216): receiver binding first, then the install
	// source, then the credential. A source that cannot be proven immutable is the
	// cli_unavailable refusal, recorded without resolving a secret or starting a
	// child.
	if verifier, ok := checker.Probe.(AIImmutableInstallSourceVerifier); ok && verifier != nil {
		if err := verifier.VerifyImmutableInstallSource(); err != nil {
			return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unavailable", unavailableProbeReason(channel), 0)
		}
	}
	// T8a ordering (CONTRACTS.md:217): a secret-auth agent-runner-cli probe may only
	// resolve a credential when the confinement of the child that will receive it is
	// PROVEN - a whitelisted argv with no logging destination, an environment
	// whitelist, only the standard stream handles, a controller-owned working
	// directory with one bounded artifact and a bounded in-memory capture, plus the
	// platform write-restriction verdict. The proof runs before the secret store is
	// touched, so an unprovable confinement never reads a credential and never
	// starts a child. Host authentication hands the child no credential of ours, so
	// the gate does not apply to it, and neither does it apply to the other
	// channels.
	if channel == "agent-runner-cli" && profile.AuthMode == "secret" {
		verifier, ok := checker.Probe.(AIBYOKConfinementVerifier)
		if !ok || verifier == nil {
			return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unavailable", "byok_confinement_unproven", 0)
		}
		if err := verifier.VerifyBYOKConfinement(); err != nil {
			return newLocalAIAvailability(profile, profileHash, channel, idSource(), now(), "unavailable", "byok_confinement_unproven", 0)
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
	// P85 F12: the refusal evidence is the PURE closed-set classification - the
	// status and the reason, both from closed sets - and nothing else. The earlier
	// form mixed the profile digest, the channel and the probe timestamp into it, so
	// two refusals of the same class produced two different evidence digests and the
	// "same class, same evidence" discipline (T8b) did not hold for local refusals.
	//
	// The profile hash, the channel and the timestamp are NOT lost: they remain the
	// record's metadata (ProfileConfigHash / Channel / ProbedAt), which is where a
	// value that legitimately differs per record belongs.
	evidenceHash := evidence.Digest("proofrail:ai-availability-local-check:1\n", []byte(status+"\n"+reason))
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
