package gates

// P44-T027B-GATE-READINESS (C7) test-first (RED): the three readiness levels
// (probe / acceptance / production) must stay strictly separated.
//
// This file must fail to compile until internal/gates/readiness.go provides every
// symbol pinned below, with exactly the names and signatures the ① architecture
// design fixed. Compile failure - not assertion failure - is the expected RED
// shape for this slice.
//
// Falsification rule applied to every case below: remove the corresponding guard
// from the implementation and the case must turn red. A case that stays green
// after its guard is deleted is not coverage and must be rewritten.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

// Compile-time pins: the published surface of the C7 design.
var (
	_ func(ReadinessGoal, ReadinessInput) (Readiness, error) = EvaluateT027BReadiness

	// The sealed interface must be satisfied by exactly these three concrete
	// levels; the accessors live on the level that granted them.
	_ Readiness = (*ProbeReadiness)(nil)
	_ Readiness = (*AcceptanceReadiness)(nil)
	_ Readiness = (*ProductionReadiness)(nil)

	// Compile-time pin replacing the never-failing runtime type assertions that
	// used to live in TestReadinessLevelsCannotBeMixed: each level exposes only
	// its own accessor (method expressions below), and the three levels are three
	// distinct struct types. "A *ProbeReadiness is never a *AcceptanceReadiness"
	// is therefore a type-system fact, not a runtime condition, and must be
	// pinned here - a runtime `if` over it can never turn red.
	_ func(*ProbeReadiness) *ProbePermit        = (*ProbeReadiness).Permit
	_ func(*AcceptanceReadiness) SealBinding    = (*AcceptanceReadiness).Acceptance
	_ func(*ProductionReadiness) ReleaseBinding = (*ProductionReadiness).Release

	_ = ReadinessInput{
		Permit:           (*ProbePermit)(nil),
		Seal:             (*SealBinding)(nil),
		AT23Hash:         "",
		ProbeEvidence:    nil,
		EvidenceRootHash: "",
	}

	_ = SealBinding{Attempt: 1, TaskKey: "task-key", SealHash: "seal-hash"}

	_ = []ReadinessGoal{GoalProbe, GoalAcceptance, GoalProduction}
)

const (
	t027BTestSealHash         = "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	t027BTestProbeEvidence    = "sha256:2222222222222222222222222222222222222222222222222222222222222222"
	t027BTestEvidenceRootHash = "sha256:3333333333333333333333333333333333333333333333333333333333333333"
)

func t027BTestSeal() *SealBinding {
	return &SealBinding{Attempt: 1, TaskKey: "task-key", SealHash: t027BTestSealHash}
}

// acceptanceMaterial is the minimal input the C7 design treats as complete
// level-2 (acceptance) material: a seal binding, probe evidence and an evidence
// root hash.
func acceptanceMaterial(root string) ReadinessInput {
	return ReadinessInput{
		Permit:           NewProbePermit("permit-t027b", filepath.Join(root, "probe-permit.json"), 1),
		Seal:             t027BTestSeal(),
		AT23Hash:         "",
		ProbeEvidence:    []string{t027BTestProbeEvidence},
		EvidenceRootHash: t027BTestEvidenceRootHash,
	}
}

// TestReadinessProbeConclusionCannotAdvanceAcceptance is C7's primary guard: an
// "available" probe run is level-1 material only and must never be promoted to a
// level-2 (acceptance) conclusion.
func TestReadinessProbeConclusionCannotAdvanceAcceptance(t *testing.T) {
	root := t.TempDir()

	probeOnly := ReadinessInput{
		Permit:           NewProbePermit("permit-probe-only", filepath.Join(root, "probe-permit.json"), 1),
		Seal:             nil,
		ProbeEvidence:    []string{"probe:available"},
		EvidenceRootHash: "",
	}

	ready, err := EvaluateT027BReadiness(GoalAcceptance, probeOnly)
	if !errors.Is(err, ErrReadinessProbeNotAcceptance) {
		t.Fatalf("a probe conclusion must not advance to acceptance: got readiness %T, err %v", ready, err)
	}
	if ready != nil {
		t.Fatalf("a rejected evaluation must not return a readiness value, got %T", ready)
	}
}

// TestReadinessLevelsCannotBeMixed pins two things at once:
//   - level-1 and level-2/3 are distinct concrete types, so the sealed interface
//     offers no compile-time path from a probe readiness value into an acceptance
//     or production slot (each level exposes only its own accessor);
//   - level 3 cannot skip level 2.
func TestReadinessLevelsCannotBeMixed(t *testing.T) {
	root := t.TempDir()

	probeInput := ReadinessInput{
		Permit: NewProbePermit("permit-level-1", filepath.Join(root, "probe-permit.json"), 1),
	}
	ready, err := EvaluateT027BReadiness(GoalProbe, probeInput)
	if err != nil {
		t.Fatalf("goal probe must succeed with a permit: %v", err)
	}
	probeReady, ok := ready.(*ProbeReadiness)
	if !ok {
		t.Fatalf("goal probe must yield *ProbeReadiness, got %T", ready)
	}
	probeReady.Permit() // level-1 accessor must exist on the level-1 type only

	// The "a probe readiness value is never an acceptance or production
	// readiness value" half of this case is NOT asserted here on purpose: the
	// three levels are distinct struct types, so `x.(*AcceptanceReadiness)` on a
	// *ProbeReadiness value can only ever yield ok=false. Written as a runtime
	// `if` it was an unfalsifiable tautology (it could never fail, so it could
	// never prove anything). It is pinned at compile time in the var block at the
	// top of this file instead; only the falsifiable half remains below.

	productionWithoutAcceptance := ReadinessInput{
		Permit:        NewProbePermit("permit-level-3-no-l2", filepath.Join(root, "probe-permit-l3.json"), 1),
		Seal:          nil,
		ProbeEvidence: []string{t027BTestProbeEvidence},
		AT23Hash:      strings.Repeat("a", 64),
	}
	advanced, err := EvaluateT027BReadiness(GoalProduction, productionWithoutAcceptance)
	if !errors.Is(err, ErrReadinessAcceptanceMissing) {
		t.Fatalf("production without level-2 material must be refused: got readiness %T, err %v", advanced, err)
	}
	if advanced != nil {
		t.Fatalf("a rejected evaluation must not return a readiness value, got %T", advanced)
	}
}

// TestReadinessProbeGoalRequiresPermit keeps ErrProbePermitMissing reachable: a
// level-1 conclusion without a permit is the exact gap N4 exists to close.
func TestReadinessProbeGoalRequiresPermit(t *testing.T) {
	ready, err := EvaluateT027BReadiness(GoalProbe, ReadinessInput{})
	if !errors.Is(err, ErrProbePermitMissing) {
		t.Fatalf("a probe readiness conclusion without a permit must be refused: got readiness %T, err %v", ready, err)
	}
	if ready != nil {
		t.Fatalf("a rejected evaluation must not return a readiness value, got %T", ready)
	}
}

// TestReadinessProductionRequiresAT23 guards the level-3 admission rule while
// proving that the same input really does reach level 2 (otherwise the AT23
// assertion below could pass on an implementation that always fails).
func TestReadinessProductionRequiresAT23(t *testing.T) {
	root := t.TempDir()

	input := acceptanceMaterial(root)
	accepted, err := EvaluateT027BReadiness(GoalAcceptance, input)
	if err != nil {
		t.Fatalf("complete level-2 material must be accepted: %v", err)
	}
	acceptanceReady, ok := accepted.(*AcceptanceReadiness)
	if !ok {
		t.Fatalf("goal acceptance must yield *AcceptanceReadiness, got %T", accepted)
	}
	acceptanceReady.Acceptance() // level-2 accessor must exist

	ready, err := EvaluateT027BReadiness(GoalProduction, input)
	if !errors.Is(err, ErrReadinessProductionAT23Required) {
		t.Fatalf("production without the AT23 hash must be refused: got readiness %T, err %v", ready, err)
	}
	if ready != nil {
		t.Fatalf("a rejected evaluation must not return a readiness value, got %T", ready)
	}

	withAT23 := input
	withAT23.AT23Hash = strings.Repeat("b", 64)
	released, err := EvaluateT027BReadiness(GoalProduction, withAT23)
	if err != nil {
		t.Fatalf("level-2 material plus the AT23 hash must reach production readiness: %v", err)
	}
	productionReady, ok := released.(*ProductionReadiness)
	if !ok {
		t.Fatalf("goal production must yield *ProductionReadiness, got %T", released)
	}
	productionReady.Release() // level-3 accessor must exist
}

// TestReadinessHasNoSideEffects proves evaluation is a pure decision: two
// consecutive evaluations neither write the filesystem nor consume the permit.
func TestReadinessHasNoSideEffects(t *testing.T) {
	root := t.TempDir()
	input := acceptanceMaterial(root)
	withAT23 := input
	withAT23.AT23Hash = strings.Repeat("c", 64)

	before := snapshotT027BTree(t, root)
	if _, err := EvaluateT027BReadiness(GoalAcceptance, input); err != nil {
		t.Fatalf("goal acceptance must evaluate complete level-2 material: %v", err)
	}
	if _, err := EvaluateT027BReadiness(GoalProduction, withAT23); err != nil {
		t.Fatalf("goal production must evaluate level-2 material plus AT23: %v", err)
	}
	after := snapshotT027BTree(t, root)

	if diff := diffT027BSnapshot(before, after); diff != "" {
		t.Fatalf("readiness evaluation must not touch the filesystem, but:\n%s", diff)
	}

	if err := input.Permit.Require(context.Background()); err != nil {
		t.Fatalf("readiness evaluation must not consume the permit: %v", err)
	}
}

// TestReadinessUnknownGoalRejected pins the default branch of the goal switch: a
// goal outside the closed level set is incomplete evidence, never a silent
// success and never a level that happens to fall through.
func TestReadinessUnknownGoalRejected(t *testing.T) {
	ready, err := EvaluateT027BReadiness(ReadinessGoal(99), ReadinessInput{})
	if !errors.Is(err, ErrReadinessEvidenceIncomplete) {
		t.Fatalf("an unknown readiness goal must be refused as incomplete evidence: got readiness %T, err %v", ready, err)
	}
	if ready != nil {
		t.Fatalf("a rejected evaluation must not return a readiness value, got %T", ready)
	}
}

// TestReadinessProbeDoesNotConsumePermit separates the decision from the quota:
// evaluating goal probe admits the run but spends nothing - it must not write the
// durable claim marker and must not burn the single claim. Only Require spends.
func TestReadinessProbeDoesNotConsumePermit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "probe-permit.json")
	permit := NewProbePermit("permit-not-consumed", path, 1)

	ready, err := EvaluateT027BReadiness(GoalProbe, ReadinessInput{Permit: permit})
	if err != nil {
		t.Fatalf("goal probe must succeed with a permit: %v", err)
	}
	if _, ok := ready.(*ProbeReadiness); !ok {
		t.Fatalf("goal probe must yield *ProbeReadiness, got %T", ready)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("evaluation must not write the claim marker, stat err = %v", statErr)
	}

	if err := permit.Require(context.Background()); err != nil {
		t.Fatalf("evaluation must leave the quota claimable, require got: %v", err)
	}
	if err := permit.Complete(); err != nil {
		t.Fatalf("complete must succeed on the claim granted above: %v", err)
	}
	if err := permit.Require(context.Background()); !errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("only require may consume the quota, a second require must be refused, got: %v", err)
	}
}

// TestReadinessAcceptanceRejectsInvalidSealBinding covers the level-2 field
// validation branch of goal acceptance: a seal that cannot be identified is not
// acceptance material, so every malformed binding must be refused as incomplete
// evidence and must not hand back a readiness value.
//
// Falsification: drop the SealBinding.validate call from the GoalAcceptance
// branch and all three sub-cases turn red.
//
// A zero attempt, a blank task key and a non-digest seal hash are the three
// ways a seal binding fails to identify the attempt it claims to carry.
func TestReadinessAcceptanceRejectsInvalidSealBinding(t *testing.T) {
	// 64 lower-case hex characters: the digest shape the contract accepts.
	validHash := strings.Repeat("ab", 32)

	cases := []struct {
		name   string
		seal   SealBinding
		reason string
	}{
		{
			name:   "zero attempt",
			seal:   SealBinding{Attempt: 0, TaskKey: "task-a", SealHash: validHash},
			reason: "a seal attempt of 0 identifies no attempt",
		},
		{
			name:   "blank task key",
			seal:   SealBinding{Attempt: 1, TaskKey: "   ", SealHash: validHash},
			reason: "a whitespace-only task key identifies no task",
		},
		{
			name:   "seal hash is not a digest",
			seal:   SealBinding{Attempt: 1, TaskKey: "task-a", SealHash: "seal-hash"},
			reason: "a seal hash that is not a 64-character digest is not evidence",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			seal := testCase.seal
			ready, err := EvaluateT027BReadiness(GoalAcceptance, ReadinessInput{Seal: &seal})
			if !errors.Is(err, ErrReadinessEvidenceIncomplete) {
				t.Fatalf("%s: %v, got readiness %T, err %v", testCase.reason, ErrReadinessEvidenceIncomplete, ready, err)
			}
			if ready != nil {
				t.Fatalf("%s: a rejected evaluation must not return a readiness value, got %T", testCase.reason, ready)
			}
		})
	}
}

// TestReadinessProductionRejectsMalformedAT23 covers the level-3 AT23 syntax
// check: a present but malformed AT23 hash is incomplete evidence, not a
// release admission.
//
// Falsification: delete the readinessHashPattern check on the AT23 hash in the
// GoalProduction branch and every sub-case turns red. The control at the end
// keeps the table honest: the same seal with a well-formed AT23 does reach level
// 3, so an implementation that simply rejected everything could not pass here.
func TestReadinessProductionRejectsMalformedAT23(t *testing.T) {
	validSeal := SealBinding{Attempt: 1, TaskKey: "task-a", SealHash: strings.Repeat("ab", 32)}

	cases := []struct {
		name   string
		at23   string
		reason string
	}{
		{
			name:   "too short",
			at23:   "abc",
			reason: "a three-character AT23 hash is not a digest",
		},
		{
			name:   "sixty three hex characters",
			at23:   strings.Repeat("d", 63),
			reason: "a 63-character AT23 hash is one character short of a digest",
		},
		{
			name:   "sixty five hex characters",
			at23:   strings.Repeat("e", 65),
			reason: "a 65-character AT23 hash is one character past a digest",
		},
		{
			name:   "upper case hex",
			at23:   strings.Repeat("A", 64),
			reason: "canonical evidence hashes are lower-case hex",
		},
		{
			name:   "non hex character",
			at23:   strings.Repeat("f", 63) + "g",
			reason: "a non-hex character makes the AT23 hash unreadable as a digest",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			seal := validSeal
			ready, err := EvaluateT027BReadiness(GoalProduction, ReadinessInput{Seal: &seal, AT23Hash: testCase.at23})
			if !errors.Is(err, ErrReadinessEvidenceIncomplete) {
				t.Fatalf("%s: malformed AT23 hash must be refused as incomplete evidence, got readiness %T, err %v", testCase.reason, ready, err)
			}
			if ready != nil {
				t.Fatalf("%s: a rejected evaluation must not return a readiness value, got %T", testCase.reason, ready)
			}
		})
	}

	t.Run("control: well formed AT23 is admitted", func(t *testing.T) {
		seal := validSeal
		ready, err := EvaluateT027BReadiness(GoalProduction, ReadinessInput{Seal: &seal, AT23Hash: strings.Repeat("a", 64)})
		if err != nil {
			t.Fatalf("the control input must reach production readiness, got: %v", err)
		}
		if _, ok := ready.(*ProductionReadiness); !ok {
			t.Fatalf("the control input must yield *ProductionReadiness, got %T", ready)
		}
	})
}

// TestReadinessAcceptsSha256PrefixedHashes is the positive control for the two
// rejection tables above: both the seal hash and the AT23 hash may carry the
// canonical "sha256:" prefix, and a fully valid input must reach level 2 and
// level 3. Without this case an implementation that refused every hash would
// satisfy the rejection tables.
//
// Falsification: tighten readinessHashPattern to bare digests only and this
// test turns red.
func TestReadinessAcceptsSha256PrefixedHashes(t *testing.T) {
	digest := strings.Repeat("12", 32)
	prefixed := "sha256:" + digest

	seal := SealBinding{Attempt: 2, TaskKey: "task-a", SealHash: prefixed}
	input := ReadinessInput{Seal: &seal, AT23Hash: prefixed}

	accepted, err := EvaluateT027BReadiness(GoalAcceptance, input)
	if err != nil {
		t.Fatalf("a sha256-prefixed seal hash must be accepted at level 2: %v", err)
	}
	acceptanceReady, ok := accepted.(*AcceptanceReadiness)
	if !ok {
		t.Fatalf("goal acceptance must yield *AcceptanceReadiness, got %T", accepted)
	}
	if got := acceptanceReady.Acceptance().SealHash; got != prefixed {
		t.Fatalf("the seal hash must be carried verbatim, got %q want %q", got, prefixed)
	}

	released, err := EvaluateT027BReadiness(GoalProduction, input)
	if err != nil {
		t.Fatalf("a sha256-prefixed AT23 hash must be accepted at level 3: %v", err)
	}
	productionReady, ok := released.(*ProductionReadiness)
	if !ok {
		t.Fatalf("goal production must yield *ProductionReadiness, got %T", released)
	}
	release := productionReady.Release()
	if release.AT23Hash != prefixed {
		t.Fatalf("the AT23 hash must be carried verbatim, got %q want %q", release.AT23Hash, prefixed)
	}
	if release.Seal.SealHash != prefixed {
		t.Fatalf("the seal hash must be carried into the release binding, got %q want %q", release.Seal.SealHash, prefixed)
	}
}

type t027BFileState struct {
	IsDir   bool
	Size    int64
	ModTime time.Time
	Digest  string
}

func snapshotT027BTree(t *testing.T, root string) map[string]t027BFileState {
	t.Helper()

	states := make(map[string]t027BFileState)
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		state := t027BFileState{IsDir: entry.IsDir(), Size: info.Size(), ModTime: info.ModTime()}
		if !entry.IsDir() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			digest := sha256.Sum256(data)
			state.Digest = hex.EncodeToString(digest[:])
		}
		states[relative] = state
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot %s: %v", root, err)
	}
	return states
}

func diffT027BSnapshot(before, after map[string]t027BFileState) string {
	paths := make([]string, 0, len(before)+len(after))
	for path := range before {
		paths = append(paths, path)
	}
	for path := range after {
		if _, seen := before[path]; !seen {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)

	var lines []string
	for _, path := range paths {
		previous, hadPrevious := before[path]
		current, hadCurrent := after[path]
		switch {
		case !hadCurrent:
			lines = append(lines, fmt.Sprintf("removed %s", path))
		case !hadPrevious:
			lines = append(lines, fmt.Sprintf("created %s", path))
		case previous.Digest != current.Digest, previous.Size != current.Size, previous.IsDir != current.IsDir:
			lines = append(lines, fmt.Sprintf("modified %s", path))
		}
	}
	return strings.Join(lines, "\n")
}
