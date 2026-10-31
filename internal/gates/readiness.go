package gates

// P44-T027B-GATE-READINESS (C7): probe readiness, acceptance readiness and
// production readiness are three separate levels, and a lower level can never
// be promoted into a higher one.
//
// C7 recorded that the gate sequence could treat an "available" probe run as if
// it were a real acceptance conclusion, and that production could then be
// released on top of that. The fix is structural rather than documentary:
//
//   - Readiness is a sealed interface (unexported tag method), so only this
//     package can introduce a readiness level;
//   - each level is its own concrete type with only its own accessor, so a
//     *ProbeReadiness value cannot be stored in a level-2 or level-3 slot and
//     then read back as one;
//   - EvaluateT027BReadiness is the only constructor and it is a pure decision:
//     it reads its input, never writes the filesystem and never consumes the
//     probe permit.
//
// Level requirements:
//
//	level 1 (probe)      : a bound *ProbePermit
//	level 2 (acceptance) : a valid *SealBinding (attempt >= 1, task key, seal hash)
//	level 3 (production) : level 2 plus a valid AT23 acceptance hash
//
// Level-1 probe evidence is never accepted as level-2 material: N4/C7 exist
// because that substitution silently upgraded a probe run into a release.

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Sentinel errors of readiness evaluation. Except for the single exception noted
// below, every error returned by EvaluateT027BReadiness wraps exactly one of
// them, so callers use errors.Is.
//
// The exception is the GoalProbe branch. It returns an error that wraps
// ErrProbePermitMissing, which is the fifth sentinel and belongs to the N4
// probe-permit domain rather than to the four readiness sentinels declared
// below. The reason is that "the probe has no bound permit" is fundamentally a
// permit-state error, not a readiness-evidence error: level 1 is defined by
// holding a permit, so the absence of one is reported in the vocabulary of
// permits. A caller evaluating GoalProbe must therefore match with
// errors.Is(err, ErrProbePermitMissing) and never with errors.Is against a
// readiness sentinel. Every other branch (GoalAcceptance, GoalProduction) wraps
// exactly one readiness sentinel.
var (
	// ErrReadinessProbeNotAcceptance marks an acceptance evaluation that only
	// holds level-1 probe evidence.
	ErrReadinessProbeNotAcceptance = errors.New("probe readiness is not acceptance readiness")
	// ErrReadinessEvidenceIncomplete marks missing or malformed readiness
	// material (seal binding fields, hash syntax, unknown goal).
	ErrReadinessEvidenceIncomplete = errors.New("readiness evidence incomplete")
	// ErrReadinessAcceptanceMissing marks a production evaluation that has no
	// level-2 material: level 3 may not skip level 2.
	ErrReadinessAcceptanceMissing = errors.New("acceptance readiness missing")
	// ErrReadinessProductionAT23Required marks a production evaluation whose
	// AT23 acceptance hash is absent.
	ErrReadinessProductionAT23Required = errors.New("production readiness requires the AT23 acceptance hash")
)

// readinessHashPattern accepts the canonical evidence hash form and its bare
// digest form; anything else is not evidence.
var readinessHashPattern = regexp.MustCompile(`^(?:sha256:)?[0-9a-f]{64}$`)

// ReadinessGoal selects the readiness level to evaluate.
type ReadinessGoal int

const (
	// GoalProbe asks for level 1: a probe may run.
	GoalProbe ReadinessGoal = iota
	// GoalAcceptance asks for level 2: work counts as accepted.
	GoalAcceptance
	// GoalProduction asks for level 3: a production release is admitted.
	GoalProduction
)

// Readiness is the sealed result of EvaluateT027BReadiness. The unexported tag
// method keeps the level set closed to this package: no other package can
// declare a fourth level or a fake implementation of an existing one.
type Readiness interface {
	readinessTag()
}

// SealBinding identifies the seal attempt that carries level-2 material.
type SealBinding struct {
	Attempt  int
	TaskKey  string
	SealHash string
}

// ReleaseBinding is the level-3 admission material handed to a release step.
type ReleaseBinding struct {
	Seal             SealBinding
	AT23Hash         string
	EvidenceRootHash string
}

// ReadinessInput is the raw material offered for evaluation. EvaluateT027BReadiness
// only reads it: it is never mutated and the permit is never consumed.
type ReadinessInput struct {
	Permit           *ProbePermit
	Seal             *SealBinding
	AT23Hash         string
	ProbeEvidence    []string
	EvidenceRootHash string
}

// ProbeReadiness is level-1 readiness: a probe run is admitted by a permit.
type ProbeReadiness struct {
	permit *ProbePermit
}

// AcceptanceReadiness is level-2 readiness: real acceptance material exists.
type AcceptanceReadiness struct {
	seal SealBinding
}

// ProductionReadiness is level-3 readiness: acceptance material plus the AT23
// acceptance hash.
type ProductionReadiness struct {
	release ReleaseBinding
}

func (*ProbeReadiness) readinessTag()      {}
func (*AcceptanceReadiness) readinessTag() {}
func (*ProductionReadiness) readinessTag() {}

// Permit returns the permit that admitted the probe run. It grants level 1
// only; it is not evidence for any higher level.
func (r *ProbeReadiness) Permit() *ProbePermit {
	if r == nil {
		return nil
	}
	return r.permit
}

// Acceptance returns the seal binding that carried the acceptance material.
func (r *AcceptanceReadiness) Acceptance() SealBinding {
	if r == nil {
		return SealBinding{}
	}
	return r.seal
}

// Release returns the material a production release is admitted on.
func (r *ProductionReadiness) Release() ReleaseBinding {
	if r == nil {
		return ReleaseBinding{}
	}
	return r.release
}

// EvaluateT027BReadiness decides whether input reaches goal. It is a pure
// function: no I/O, no mutation of input, no permit consumption.
func EvaluateT027BReadiness(goal ReadinessGoal, input ReadinessInput) (Readiness, error) {
	switch goal {
	case GoalProbe:
		// The decision checks presence only: it neither consumes the quota nor
		// inspects the permit state. Quota consumption and state validation belong
		// to ProbePermit.Require, which keeps the two responsibilities apart. An
		// already Spent permit therefore still yields a usable ProbeReadiness
		// here; that is by design rather than a quota bypass, because actually
		// taking the quota always goes through Require.
		if input.Permit == nil {
			return nil, fmt.Errorf("%w: goal probe needs a bound probe permit", ErrProbePermitMissing)
		}
		return &ProbeReadiness{permit: input.Permit}, nil

	case GoalAcceptance:
		if input.Seal == nil {
			if len(input.ProbeEvidence) > 0 {
				// The C7 hole itself: probe evidence is level 1 and stays there.
				return nil, fmt.Errorf("%w: probe evidence cannot stand in for a seal binding", ErrReadinessProbeNotAcceptance)
			}
			return nil, fmt.Errorf("%w: goal acceptance needs a seal binding", ErrReadinessEvidenceIncomplete)
		}
		if err := input.Seal.validate(); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrReadinessEvidenceIncomplete, err)
		}
		return &AcceptanceReadiness{seal: *input.Seal}, nil

	case GoalProduction:
		if input.Seal == nil {
			// Level 3 cannot skip level 2, whatever probe evidence is present.
			//
			// Derived decision: the "level-2 material is missing" test is
			// Seal == nil alone and deliberately ignores ProbeEvidence. Probe
			// evidence is never material for a higher level, so level 3 without
			// level 2 is always ErrReadinessAcceptanceMissing, even when probe
			// evidence is present, which matches C7 (a probe conclusion must never
			// advance acceptance).
			// See the revision note in
			// docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md
			// (derived decision).
			return nil, fmt.Errorf("%w: level 3 needs level-2 seal material", ErrReadinessAcceptanceMissing)
		}
		if err := input.Seal.validate(); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrReadinessEvidenceIncomplete, err)
		}
		at23Hash := strings.TrimSpace(input.AT23Hash)
		if at23Hash == "" {
			return nil, fmt.Errorf("%w: goal production needs the AT23 acceptance hash", ErrReadinessProductionAT23Required)
		}
		if !readinessHashPattern.MatchString(at23Hash) {
			return nil, fmt.Errorf("%w: invalid AT23 hash %q", ErrReadinessEvidenceIncomplete, input.AT23Hash)
		}
		return &ProductionReadiness{release: ReleaseBinding{
			Seal:             *input.Seal,
			AT23Hash:         at23Hash,
			EvidenceRootHash: input.EvidenceRootHash,
		}}, nil

	default:
		return nil, fmt.Errorf("%w: unknown readiness goal %d", ErrReadinessEvidenceIncomplete, goal)
	}
}

// validate checks the level-2 binding fields. A seal that cannot be identified
// is not acceptance material.
func (s SealBinding) validate() error {
	if s.Attempt < 1 {
		return fmt.Errorf("seal attempt %d must be at least 1", s.Attempt)
	}
	if strings.TrimSpace(s.TaskKey) == "" {
		return errors.New("seal task key is required")
	}
	// No trimming here: the seal hash is stored verbatim, so a padded value
	// must be rejected rather than accepted and then kept unnormalised.
	if !readinessHashPattern.MatchString(s.SealHash) {
		return fmt.Errorf("invalid seal hash %q", s.SealHash)
	}
	return nil
}
