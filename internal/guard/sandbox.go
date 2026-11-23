package guard

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// ErrSandboxUnavailable reports that the sandbox check capability is not
// available for this build or platform (no implementation, or the
// implementation cannot answer). Callers must treat it as fail-closed.
var ErrSandboxUnavailable = errors.New("guard: sandbox check unavailable")

// ErrSandboxRejected reports that the sandbox check refused this dispatch
// fail-closed: the verdict was Denied or Unknown. Callers must treat it as a
// refusal, never as a pass.
var ErrSandboxRejected = errors.New("guard: sandbox check rejected dispatch")

// VerdictStatus enumerates the outcome of one sandbox check.
type VerdictStatus uint8

const (
	// VerdictUnknown means compliance could not be determined. It is a
	// fail-closed outcome: the caller must refuse to dispatch.
	VerdictUnknown VerdictStatus = iota
	// VerdictAllowed means the environment satisfies the enforcement
	// requirements and dispatch may continue unchanged.
	VerdictAllowed
	// VerdictDenied means the environment provably does not satisfy the
	// enforcement requirements. It is a fail-closed outcome.
	VerdictDenied
)

// String renders the status for audit and diagnostic output.
func (status VerdictStatus) String() string {
	switch status {
	case VerdictAllowed:
		return "allowed"
	case VerdictDenied:
		return "denied"
	default:
		return "unknown"
	}
}

// Verdict is the answer of one Sandbox check. Denied and Unknown are both
// fail-closed: only Allowed lets a dispatch continue.
type Verdict struct {
	Status VerdictStatus
	// Reason is a human readable explanation that travels with the refusal
	// error into the audit and diagnostic output.
	Reason string
	// SpecDigest is the digest of the specification this verdict ANSWERED FOR
	// (P85 F13). A check layer echoes the CheckSpec.SpecDigest it was asked
	// about, so "which specification was judged" becomes a fact the caller can
	// compare instead of an assumption about the caller. CheckSandbox refuses a
	// verdict whose echo does not equal the digest it asked about, so an answer
	// about one specification can never be read as an answer about another.
	SpecDigest string
}

// IsAllowed reports whether the verdict lets a dispatch continue.
func (verdict Verdict) IsAllowed() bool { return verdict.Status == VerdictAllowed }

// CheckSpec is the input of one sandbox check. It carries the audit label of
// the dispatch being checked plus, when the caller has one, the digest of the
// FROZEN specification the answer has to be about: the check layer creates,
// runs and destroys no process.
type CheckSpec struct {
	// Reason labels the dispatch that triggered the check (for example
	// "agent-runner-dispatch").
	Reason string
	// SpecDigest is the digest of the specification the verdict must answer for
	// (P85 F13). It turns "the sandbox was allowed" into "the sandbox allowed
	// THIS specification": a verdict that echoes a different digest is refused,
	// so two specifications on the same platform can never share one verdict.
	// It is REQUIRED for SandboxReasonBYOKProbe and optional elsewhere.
	SpecDigest string
}

// SandboxReasonBYOKProbe labels the check a probe runs before it hands a
// credential to a child (T8a, finding F8). The dispatch may only continue when
// the environment's write restrictions are PROVEN, so an Unknown verdict is a
// refusal exactly like Denied: S1 admits no BYOK probe.
const SandboxReasonBYOKProbe = "ai-byok-probe"

// Sandbox is the sandbox check layer: it answers one question, namely whether
// the current environment satisfies the enforcement requirements. It enforces
// nothing itself.
//
// Contract:
//   - The method set carries no Create/Exec/Destroy/Terminate process lifecycle
//     method and overlaps no type with ProcessSpec, ProcessIdentity or
//     ManagedProcess.
//   - A non-nil error and any verdict other than Allowed are both fail-closed
//     refusals.
type Sandbox interface {
	Check(ctx context.Context, spec CheckSpec) (Verdict, error)
}

// NewSandbox returns the check implementation of the current platform (see
// newPlatformSandbox in sandbox_windows.go and sandbox_unsupported.go).
func NewSandbox() Sandbox { return newPlatformSandbox() }

// CheckSandbox is the dispatch-path gate: the check must be Allowed, and the
// checker must not error, or the dispatch is refused fail-closed. A nil
// implementation falls back to NewSandbox().
func CheckSandbox(ctx context.Context, sandbox Sandbox, spec CheckSpec) error {
	// P85 F13: a BYOK probe asks whether the environment may hand a credential to
	// a child, and that question is only meaningful for one FROZEN specification.
	// A caller that asks it without naming the specification is refused here, so a
	// call site that predates the binding can never read as "the sandbox allowed the
	// spawn we are about to make".
	if spec.Reason == SandboxReasonBYOKProbe && spec.SpecDigest == "" {
		return fmt.Errorf("%w: the BYOK probe check carries no specification digest, so the verdict cannot be bound to a frozen spawn", ErrSandboxRejected)
	}
	if sandbox == nil {
		sandbox = NewSandbox()
	}
	verdict, err := sandbox.Check(ctx, spec)
	if err != nil {
		// The verdict reason is the audit answer of the refusal, so it leads the
		// error text whenever the implementation supplies one. The sentinel is
		// never repeated and an already-wrapped sentinel is returned as is.
		reason := strings.TrimSpace(verdict.Reason)
		switch {
		case reason == "" && errors.Is(err, ErrSandboxUnavailable):
			return err
		case reason == "":
			return fmt.Errorf("%w: %v", ErrSandboxUnavailable, err)
		case errors.Is(err, ErrSandboxUnavailable):
			return fmt.Errorf("%w: %s", ErrSandboxUnavailable, reason)
		default:
			return fmt.Errorf("%w: %s: %v", ErrSandboxUnavailable, reason, err)
		}
	}
	if !verdict.IsAllowed() {
		return fmt.Errorf("%w: %s", ErrSandboxRejected, verdict.Reason)
	}
	// P85 F13: the answer has to be about the specification that was asked about.
	// The echo is compared even on the Allowed path - an Allowed verdict for a
	// different specification is not an Allowed verdict for this one.
	if spec.SpecDigest != "" && verdict.SpecDigest != spec.SpecDigest {
		return fmt.Errorf("%w: %s", ErrSandboxRejected, "the sandbox answered for a different specification than the one it was asked about")
	}
	return nil
}
