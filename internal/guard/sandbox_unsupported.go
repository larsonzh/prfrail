//go:build !windows

package guard

import "context"

// unsupportedSandbox is the check implementation of every non-Windows
// platform: there is no sandbox check capability, so every check reports
// Unknown together with ErrSandboxUnavailable (fail-closed).
type unsupportedSandbox struct{}

func (unsupportedSandbox) Check(_ context.Context, spec CheckSpec) (Verdict, error) {
	// P85 F13: the verdict echoes the specification it answered for, so a caller
	// can compare "which specification was judged" instead of assuming it.
	return Verdict{
		Status:     VerdictUnknown,
		Reason:     "sandbox check not supported on this platform",
		SpecDigest: spec.SpecDigest,
	}, ErrSandboxUnavailable
}

func newPlatformSandbox() Sandbox { return unsupportedSandbox{} }
