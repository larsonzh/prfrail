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
}

// IsAllowed reports whether the verdict lets a dispatch continue.
func (verdict Verdict) IsAllowed() bool { return verdict.Status == VerdictAllowed }

// CheckSpec is the input of one sandbox check. It carries the audit label of
// the dispatch being checked and nothing else: the check layer creates, runs
// and destroys no process.
type CheckSpec struct {
	// Reason labels the dispatch that triggered the check (for example
	// "agent-runner-dispatch").
	Reason string
}

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
	return nil
}
