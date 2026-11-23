package adapters

// P85 F13 (design rev 3, section 5 / test row M9, adapters half): the platform write
// proof is bound to ONE frozen specification.
//
// Two specifications on the same platform that differ in a single field (here the
// output bound) must produce different specification digests AND different
// PlatformWriteProof digests. Before the fix the proof was GOOS+GOARCH+reason+
// "allowed", so both specifications shared one proof and "the sandbox allowed THIS
// spawn" was unprovable.
//
// Mutation check: drop specDigest from the canonical text in
// copilotCLIProbePlatformWriteVerdict (back to the constant "allowed" tail) =>
// "the-proof-follows-the-specification" goes red.

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/guard"
)

// p85F13SilentAllowedSandbox answers Allowed but echoes NO specification, which is
// the shape of a layer that judged "a probe in general".
type p85F13SilentAllowedSandbox struct{}

func (p85F13SilentAllowedSandbox) Check(_ context.Context, _ guard.CheckSpec) (guard.Verdict, error) {
	return guard.Verdict{Status: guard.VerdictAllowed}, nil
}

func p85F13BaseSpawnSpec(t *testing.T) copilotCLIProbeSpawnSpec {
	t.Helper()
	spec, err := copilotCLIProbeLaunchSpec(
		filepath.Join(t.TempDir(), "copilot"),
		t.TempDir(),
		"secret",
		p85T8ABYOKProfile().Model,
		1,
		1<<20,
	)
	if err != nil {
		t.Fatalf("a compliant launch template must be constructible: %v", err)
	}
	return spec
}

func TestP85F13PlatformWriteProofIsBoundToTheSpawnSpecification(t *testing.T) {
	base := p85F13BaseSpawnSpec(t)
	other := base
	other.OutputBound++

	baseDigest := copilotCLIProbeSpawnSpecDigest(base)
	otherDigest := copilotCLIProbeSpawnSpecDigest(other)
	if baseDigest == "" || otherDigest == "" {
		t.Fatalf("both specifications must encode a digest, got %q and %q", baseDigest, otherDigest)
	}
	if baseDigest == otherDigest {
		t.Fatalf("two specifications that differ in one field must not share a digest (%s)", baseDigest)
	}
	if again := copilotCLIProbeSpawnSpecDigest(base); again != baseDigest {
		t.Fatalf("the specification digest must be stable, got %s then %s", baseDigest, again)
	}

	sandbox := p85T8AAllowedSandbox{}
	baseProof, err := copilotCLIProbePlatformWriteVerdict(context.Background(), sandbox, baseDigest)
	if err != nil {
		t.Fatalf("an Allowed sandbox must answer the base specification, got %v", err)
	}
	otherProof, err := copilotCLIProbePlatformWriteVerdict(context.Background(), sandbox, otherDigest)
	if err != nil {
		t.Fatalf("an Allowed sandbox must answer the other specification, got %v", err)
	}
	if baseProof == otherProof {
		t.Fatalf("two specifications on one platform must not share a write proof (%s)", baseProof)
	}
	sameProof, err := copilotCLIProbePlatformWriteVerdict(context.Background(), sandbox, baseDigest)
	if err != nil || sameProof != baseProof {
		t.Fatalf("the proof must be stable for one specification, got %s then %s (err=%v)", baseProof, sameProof, err)
	}

	// A layer that answers without echoing the specification is refused, and so is a
	// request that names no specification at all: an unbound verdict is never a proof.
	if _, err := copilotCLIProbePlatformWriteVerdict(context.Background(), p85F13SilentAllowedSandbox{}, otherDigest); err == nil {
		t.Fatal("a verdict that does not echo the specification must be refused")
	}
	if _, err := copilotCLIProbePlatformWriteVerdict(context.Background(), sandbox, ""); err == nil {
		t.Fatal("an unbound specification must be refused")
	}
}
