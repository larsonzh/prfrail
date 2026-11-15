//go:build windows

package guard

import "context"

// alphaSandbox is the Windows implementation of stage S1 (form alpha): S1 runs
// no check at all, so every check reports Unknown and says so. Unknown is
// fail-closed exactly like Denied, and it is the honest answer: S1 does not
// know whether the environment complies, so it must not claim that the
// environment provably does not comply.
type alphaSandbox struct{}

func (alphaSandbox) Check(_ context.Context, _ CheckSpec) (Verdict, error) {
	return Verdict{
		Status: VerdictUnknown,
		Reason: "S1 sandbox enforcement checks not implemented; environment compliance unknown",
	}, nil
}

func newPlatformSandbox() Sandbox { return alphaSandbox{} }
