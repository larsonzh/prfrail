//go:build !windows

package guard

import "context"

// unsupportedSandbox is the check implementation of every non-Windows
// platform: there is no sandbox check capability, so every check reports
// Unknown together with ErrSandboxUnavailable (fail-closed).
type unsupportedSandbox struct{}

func (unsupportedSandbox) Check(_ context.Context, _ CheckSpec) (Verdict, error) {
	return Verdict{
		Status: VerdictUnknown,
		Reason: "sandbox check not supported on this platform",
	}, ErrSandboxUnavailable
}

func newPlatformSandbox() Sandbox { return unsupportedSandbox{} }
