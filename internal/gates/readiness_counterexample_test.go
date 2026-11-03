package gates

// P52-S3 (ACCEPTANCE-SECURITY-COUNTEREXAMPLES): counterexamples for readiness
// evidence validation.
//
// These cases pin the order of the level-3 checks, the exact hash syntax the
// contract accepts, the no-trimming rule for the seal hash, and the trimming
// behaviour of the AT23 hash as it stands today (GAP-79 candidate).
//
// Falsification rule applied to every case below: remove the corresponding guard
// from internal/gates/readiness.go and the case must turn red.

import (
	"errors"
	"strings"
	"testing"
)

// TestReadinessProductionRejectsMalformedNonEmptyAT23 pins the AT23 syntax check
// of level 3 for a non-empty but malformed hash. A present AT23 value that is
// not a canonical digest is incomplete evidence, never a release admission.
//
// Contract basis: readinessHashPattern accepts `^(?:sha256:)?[0-9a-f]{64}$`,
// i.e. an optional lower-case "sha256:" prefix followed by exactly 64
// lower-case hex characters. Every form below violates one clause of that
// pattern while remaining non-empty, so the AT23Required branch cannot be the
// rejector and the syntax check must be.
func TestReadinessProductionRejectsMalformedNonEmptyAT23(t *testing.T) {
	// 64 lower-case hex characters: the digest shape the contract accepts.
	validHash := strings.Repeat("ab", 32)
	validSeal := SealBinding{Attempt: 1, TaskKey: "task-a", SealHash: validHash}

	cases := []struct {
		name   string
		at23   string
		reason string
	}{
		{
			name:   "prefixed hash one hex character short",
			at23:   "sha256:" + strings.Repeat("d", 63),
			reason: "a sha256-prefixed 63-character digest is one character short",
		},
		{
			name:   "prefixed hash one hex character long",
			at23:   "sha256:" + strings.Repeat("d", 65),
			reason: "a sha256-prefixed 65-character digest is one character past a digest",
		},
		{
			name:   "prefixed hash with non hex character",
			at23:   "sha256:" + strings.Repeat("d", 63) + "g",
			reason: "a non-hex character makes the AT23 hash unreadable as a digest",
		},
		{
			name:   "upper case prefix",
			at23:   "SHA256:" + strings.Repeat("d", 64),
			reason: "canonical evidence hashes use the lower-case sha256: prefix",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			seal := validSeal
			ready, err := EvaluateT027BReadiness(GoalProduction, ReadinessInput{Seal: &seal, AT23Hash: testCase.at23})
			if !errors.Is(err, ErrReadinessEvidenceIncomplete) {
				t.Fatalf("%s: malformed non-empty AT23 hash must be refused as %v, got readiness %T, err %v", testCase.reason, ErrReadinessEvidenceIncomplete, ready, err)
			}
			if ready != nil {
				t.Fatalf("%s: a rejected evaluation must not return a readiness value, got %T", testCase.reason, ready)
			}
		})
	}
}

// TestReadinessProductionValidatesSealBeforeAT23 pins the order of the level-3
// checks: the seal binding is validated before the AT23 hash is consumed. An
// invalid seal plus a perfectly good AT23 hash must be refused for the seal, so
// a malformed level-3 request can never be reported as "AT23 required/malformed"
// and hide the broken seal underneath.
//
// Falsification: move the Seal.validate call after the AT23 checks and this case
// still passes on the sentinel but turns red on the ordering guarantee only if
// the AT23 branch is also reachable - so the assertion below additionally pins
// that the AT23 path was not the reported reason.
func TestReadinessProductionValidatesSealBeforeAT23(t *testing.T) {
	validAT23 := strings.Repeat("ab", 32)
	// Attempt 0 identifies no attempt: the seal is invalid on its first field.
	invalidSeal := SealBinding{Attempt: 0, TaskKey: "task-a", SealHash: strings.Repeat("cd", 32)}

	ready, err := EvaluateT027BReadiness(GoalProduction, ReadinessInput{Seal: &invalidSeal, AT23Hash: validAT23})
	if !errors.Is(err, ErrReadinessEvidenceIncomplete) {
		t.Fatalf("an invalid seal must be refused as %v even when the AT23 hash is valid, got readiness %T, err %v", ErrReadinessEvidenceIncomplete, ready, err)
	}
	if ready != nil {
		t.Fatalf("a rejected evaluation must not return a readiness value, got %T", ready)
	}
	// The reported reason must be the seal, not the AT23 hash: a valid AT23 can
	// only be complained about if the seal check had been skipped.
	if !strings.Contains(err.Error(), "seal attempt") {
		t.Fatalf("the seal must be validated before the AT23 hash is consumed, got: %v", err)
	}
}

// TestReadinessAcceptanceRejectsPaddedSealHash pins the no-trimming rule of
// SealBinding.validate: the seal hash is stored verbatim, so a padded value must
// be rejected rather than accepted and then carried unnormalised.
//
// Falsification: wrap the SealHash match in strings.TrimSpace and this case
// turns red.
func TestReadinessAcceptanceRejectsPaddedSealHash(t *testing.T) {
	seal := SealBinding{Attempt: 1, TaskKey: "t", SealHash: strings.Repeat("ab", 32) + " "}

	ready, err := EvaluateT027BReadiness(GoalAcceptance, ReadinessInput{Seal: &seal})
	if !errors.Is(err, ErrReadinessEvidenceIncomplete) {
		t.Fatalf("a trailing-space seal hash must be refused as %v, got readiness %T, err %v", ErrReadinessEvidenceIncomplete, ready, err)
	}
	if ready != nil {
		t.Fatalf("a rejected evaluation must not return a readiness value, got %T", ready)
	}
}

// TestReadinessAcceptanceAdmitsTrimmedAT23 is a behaviour pin, not an
// endorsement. Level 3 currently trims the AT23 hash before matching it and
// stores the trimmed value, while the seal hash is never trimmed. This test
// records that asymmetry as it stands so that a future change cannot move it
// silently.
//
// GAP-79 candidate: the AT23 hash being trimmed while the seal hash is not is a
// symmetry gap worth a contract decision. This case pins the CURRENT behaviour
// only; it must be updated together with the contract if the decision goes the
// other way.
//
// Falsification: remove the strings.TrimSpace on input.AT23Hash and this case
// turns red.
func TestReadinessAcceptanceAdmitsTrimmedAT23(t *testing.T) {
	bare := strings.Repeat("ab", 32)
	validSeal := SealBinding{Attempt: 1, TaskKey: "task-a", SealHash: bare}

	ready, err := EvaluateT027BReadiness(GoalProduction, ReadinessInput{
		Seal:     &validSeal,
		AT23Hash: "  " + bare + "  ",
	})
	if err != nil {
		t.Fatalf("a whitespace-padded but otherwise valid AT23 hash is admitted today, got: %v", err)
	}
	productionReady, ok := ready.(*ProductionReadiness)
	if !ok {
		t.Fatalf("goal production must yield *ProductionReadiness, got %T", ready)
	}
	if got := productionReady.Release().AT23Hash; got != bare {
		t.Fatalf("the AT23 hash must be carried trimmed, got %q want %q", got, bare)
	}
}
