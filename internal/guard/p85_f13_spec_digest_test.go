package guard

// P85 F13 (design rev 3, section 5 / test row M9, guard half): a verdict answers
// ONE frozen specification.
//
// The check spec carries the digest of the specification the dispatch is about, and
// the verdict echoes the digest it answered for. CheckSandbox compares the echo with
// the asked digest even on the Allowed path, so "the sandbox allowed a probe" can
// never be read as "the sandbox allowed THIS probe"; a BYOK check that names no
// specification is refused before the layer is consulted.
//
// Mutation check: delete the echo comparison (`spec.SpecDigest != "" &&
// verdict.SpecDigest != spec.SpecDigest`) or the empty-digest refusal for the BYOK
// reason => the named subtest goes red.

import (
	"context"
	"testing"
)

// p85F13EchoDouble answers one verdict and echoes the spec digest it was handed only
// when echo is true, so the "answers about another specification" case is scripted
// instead of delegated to a platform layer.
type p85F13EchoDouble struct {
	verdict Verdict
	echo    bool
	specs   []CheckSpec
}

func (double *p85F13EchoDouble) Check(_ context.Context, spec CheckSpec) (Verdict, error) {
	double.specs = append(double.specs, spec)
	verdict := double.verdict
	if double.echo {
		verdict.SpecDigest = spec.SpecDigest
	}
	return verdict, nil
}

func TestP85F13SandboxVerdictMustAnswerTheAskedSpecification(t *testing.T) {
	const asked = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

	t.Run("a-matching-echo-admits", func(t *testing.T) {
		double := &p85F13EchoDouble{verdict: Verdict{Status: VerdictAllowed}, echo: true}
		if err := CheckSandbox(context.Background(), double, CheckSpec{Reason: SandboxReasonBYOKProbe, SpecDigest: asked}); err != nil {
			t.Fatalf("an Allowed verdict that echoes the asked specification must admit, got %v", err)
		}
		if len(double.specs) != 1 || double.specs[0].SpecDigest != asked {
			t.Fatalf("the check must be asked about the caller's specification, got %+v", double.specs)
		}
	})

	t.Run("a-silent-allowed-verdict-is-refused", func(t *testing.T) {
		double := &p85F13EchoDouble{verdict: Verdict{Status: VerdictAllowed}, echo: false}
		err := CheckSandbox(context.Background(), double, CheckSpec{Reason: SandboxReasonBYOKProbe, SpecDigest: asked})
		if err == nil {
			t.Fatal("an Allowed verdict that answers for no specification must be refused")
		}
	})

	t.Run("an-allowed-verdict-for-another-specification-is-refused", func(t *testing.T) {
		// The double echoes a DIFFERENT digest: the shape a layer that answered about
		// some other spawn would produce.
		double := &p85F13EchoDouble{verdict: Verdict{Status: VerdictAllowed, SpecDigest: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}, echo: false}
		if err := CheckSandbox(context.Background(), double, CheckSpec{Reason: SandboxReasonBYOKProbe, SpecDigest: asked}); err == nil {
			t.Fatal("an Allowed verdict for another specification must be refused")
		}
	})

	t.Run("a-byok-check-without-a-specification-is-refused-before-the-layer", func(t *testing.T) {
		double := &p85F13EchoDouble{verdict: Verdict{Status: VerdictAllowed}, echo: true}
		if err := CheckSandbox(context.Background(), double, CheckSpec{Reason: SandboxReasonBYOKProbe}); err == nil {
			t.Fatal("a BYOK check that names no specification must be refused")
		}
		if len(double.specs) != 0 {
			t.Fatalf("the layer must not be consulted for an unbound check, got %+v", double.specs)
		}
	})

	t.Run("a-non-byok-check-may-carry-no-specification", func(t *testing.T) {
		double := &p85F13EchoDouble{verdict: Verdict{Status: VerdictAllowed}, echo: true}
		if err := CheckSandbox(context.Background(), double, CheckSpec{Reason: "agent-runner-dispatch"}); err != nil {
			t.Fatalf("an optional specification must stay optional for a non-BYOK label, got %v", err)
		}
	})
}
