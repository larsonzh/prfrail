package guard

import (
	"context"
	"errors"
	"testing"
)

// sandboxCheckDouble answers one scripted verdict and records the specs it saw.
type sandboxCheckDouble struct {
	verdict Verdict
	err     error
	specs   []CheckSpec
}

func (double *sandboxCheckDouble) Check(_ context.Context, spec CheckSpec) (Verdict, error) {
	double.specs = append(double.specs, spec)
	return double.verdict, double.err
}

// TestSandboxInterfaceContract pins the check-layer contract: the default
// implementation always exists, the three verdicts are distinct, only Allowed
// is allowed, and the CheckSandbox gate turns every non-Allowed answer and
// every checker error into a fail-closed sentinel.
func TestSandboxInterfaceContract(t *testing.T) {
	if NewSandbox() == nil {
		t.Fatal("NewSandbox must never return nil: a nil default would turn every check into a silent pass")
	}
	if VerdictUnknown == VerdictAllowed || VerdictUnknown == VerdictDenied || VerdictAllowed == VerdictDenied {
		t.Fatalf("the three verdict statuses must be distinct: unknown=%d allowed=%d denied=%d",
			VerdictUnknown, VerdictAllowed, VerdictDenied)
	}
	if (Verdict{Status: VerdictUnknown}).IsAllowed() {
		t.Fatal("an Unknown verdict must not be allowed")
	}
	if (Verdict{Status: VerdictDenied}).IsAllowed() {
		t.Fatal("a Denied verdict must not be allowed")
	}
	if !(Verdict{Status: VerdictAllowed}).IsAllowed() {
		t.Fatal("an Allowed verdict must be allowed")
	}

	checkFailure := errors.New("checker exploded")
	cases := []struct {
		name    string
		verdict Verdict
		err     error
		want    error
	}{
		{name: "allowed", verdict: Verdict{Status: VerdictAllowed}},
		{name: "denied", verdict: Verdict{Status: VerdictDenied, Reason: "denied by double"}, want: ErrSandboxRejected},
		{name: "unknown", verdict: Verdict{Status: VerdictUnknown, Reason: "unknown by double"}, want: ErrSandboxRejected},
		{name: "unavailable", err: checkFailure, want: ErrSandboxUnavailable},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			double := &sandboxCheckDouble{verdict: testCase.verdict, err: testCase.err}
			err := CheckSandbox(context.Background(), double, CheckSpec{Reason: "unit"})
			if testCase.want == nil {
				if err != nil {
					t.Fatalf("an Allowed verdict must not refuse the dispatch: %v", err)
				}
			} else if !errors.Is(err, testCase.want) {
				t.Fatalf("the refusal must wrap %v, got %v", testCase.want, err)
			}
			if len(double.specs) != 1 || double.specs[0].Reason != "unit" {
				t.Fatalf("the gate must call the checker exactly once with the caller's spec, got %v", double.specs)
			}
		})
	}

	// A nil checker falls back to the platform implementation, which is
	// fail-closed on every platform in S1.
	if err := CheckSandbox(context.Background(), nil, CheckSpec{Reason: "nil"}); err == nil {
		t.Fatal("a nil sandbox must fail closed through the platform default")
	}
}
