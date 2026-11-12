package console

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/gates"
)

// This file is the P69 T027B-E2E-WIRING test carrier for the probe-permit gate
// of the production assembly entries (readiness_gate.go) and of the CLI `run`
// route. Every leg is deterministic: temporary directories only, no network and
// no external binary unless the full chain is actually reachable on the
// platform.
//
// The gate has two independent obligations, and both are pinned here:
//
//   - presence is decided BEFORE any precondition, so a nil or already spent
//     permit refuses with zero footprint (no run-root marker, no spawn);
//   - quota consumption happens ONLY after every precondition passed, and the
//     release deferred by the entry spends the permit exactly once, so a second
//     run on the same permit observes ErrProbePermitConsumed and never
//     ErrProbePermitInUse.

// dropCLIFlag removes one flag and its value from a CLI argument slice. It is
// the same-package seam the wiring tests use to build a deliberately incomplete
// invocation without duplicating the seven agent-runner flags.
func dropCLIFlag(args []string, name string) []string {
	filtered := make([]string, 0, len(args))
	for index := 0; index < len(args); index++ {
		if args[index] == name {
			index++ // skip the value of the dropped flag as well
			continue
		}
		filtered = append(filtered, args[index])
	}
	return filtered
}

// TestProbePermitConsumedAfterSuccessfulAssembly is the positive leg: a complete
// AgentRunner assembly spends its one-shot permit, persists the durable marker
// and refuses every later claim with the consumed sentinel.
//
// The chain itself completes only where the replay store proves publication
// durability (the production dispatcher refuses R anywhere else, see
// TestAgentRunnerAssemblyFullChain). On a platform without that proof the
// assembly refuses with the durability sentinel AFTER the permit claim, so the
// consumption assertions below are still the production behaviour under test;
// the success assertions apply where the chain can run.
func TestProbePermitConsumedAfterSuccessfulAssembly(t *testing.T) {
	durability := agentRunnerPublicationDurability(t)
	knobs := assemblyKnobs{}
	if durability == adapters.PublishDurabilityProven {
		knobs.declaredExecutable = buildAssemblyStub(t)
		knobs.runArgs = []string{"-mode", "mutate-workspace", "-sleep", "100ms"}
	}
	fixture := newAssemblyFixture(t, knobs)
	permitPath := filepath.Join(t.TempDir(), "probe-permit.marker")
	permit := gates.NewProbePermit("t027-consumed", permitPath, 1)
	fixture.options.ProbePermit = permit
	if durability == adapters.PublishDurabilityProven {
		agentRunnerSeedChildPIDReport(t, fixture.runRoot)
	}

	ctx := context.Background()
	summary, err := ExecuteAgentRunnerRun(ctx, fixture.options)
	if durability == adapters.PublishDurabilityProven {
		if err != nil {
			t.Fatalf("the production assembly must complete on a durability-proven platform: %v", err)
		}
		if summary.ChainState != "COMPLETED" {
			t.Fatalf("chain state = %s, want COMPLETED (%+v)", summary.ChainState, summary)
		}
	} else if !errors.Is(err, adapters.ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("the run must either complete or refuse on unproven durability, got %v", err)
	}

	// The claim is durable: the marker file itself is the cross-process bit.
	if _, statErr := os.Stat(permitPath); statErr != nil {
		t.Fatalf("the assembly must persist the one-shot permit claim marker: %v", statErr)
	}
	loaded, loadErr := gates.LoadProbePermit(permitPath)
	if loadErr != nil {
		t.Fatalf("the persisted claim must load: %v", loadErr)
	}
	if requireErr := loaded.Require(ctx); !errors.Is(requireErr, gates.ErrProbePermitConsumed) {
		t.Fatalf("a reloaded claim must refuse with the consumed sentinel, got %v", requireErr)
	}

	// A second run on the same permit object must observe the consumed sentinel,
	// never the in-use one: Complete ran through the deferred release.
	if _, secondErr := ExecuteAgentRunnerRun(ctx, fixture.options); !errors.Is(secondErr, gates.ErrProbePermitConsumed) {
		t.Fatalf("a second run on the same permit must refuse with the consumed sentinel, got %v", secondErr)
	} else if errors.Is(secondErr, gates.ErrProbePermitInUse) {
		t.Fatalf("a completed permit must never be reported as in use: %v", secondErr)
	}
}

// TestAgentRunnerAssemblyRefusesNilProbePermit is reverse leg 1: a nil permit is
// refused by the Step 0 readiness evaluation before any precondition runs, so
// the run root carries no event-log marker at all.
func TestAgentRunnerAssemblyRefusesNilProbePermit(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{})
	fixture.options.ProbePermit = nil

	_, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if !errors.Is(err, gates.ErrProbePermitMissing) {
		t.Fatalf("a nil probe permit must be refused with the missing sentinel, got %v", err)
	}
	if _, statErr := os.Stat(runEventLogPath(fixture.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a readiness refusal must write no run-root marker (stat err=%v)", statErr)
	}
}

// TestSupervisedAssemblyRefusesNilProbePermit is the supervised twin of reverse
// leg 1: a fresh supervised start without a permit leaves the run root clean.
func TestSupervisedAssemblyRefusesNilProbePermit(t *testing.T) {
	harness := newSupervisedHarness(t, []string{"complete"}, false, 60000)
	harness.options.ProbePermit = nil

	_, err := ExecuteSupervisedRun(context.Background(), harness.options)
	if !errors.Is(err, gates.ErrProbePermitMissing) {
		t.Fatalf("a fresh supervised start without a permit must be refused with the missing sentinel, got %v", err)
	}
	if _, statErr := os.Stat(runEventLogPath(harness.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a readiness refusal must write no run-root marker (stat err=%v)", statErr)
	}
}

// TestAgentRunnerAssemblyRefusesSpentProbePermit is reverse leg 2: an already
// spent permit passes the presence evaluation (presence only) but is refused by
// Require, which runs before the run-root marker and before any dispatch.
func TestAgentRunnerAssemblyRefusesSpentProbePermit(t *testing.T) {
	ctx := context.Background()
	permit := newTestProbePermit(t)
	if err := permit.Require(ctx); err != nil {
		t.Fatalf("the pre-claimed permit fixture must claim: %v", err)
	}
	if err := permit.Complete(); err != nil {
		t.Fatalf("the pre-claimed permit fixture must complete: %v", err)
	}
	fixture := newAssemblyFixture(t, assemblyKnobs{})
	fixture.options.ProbePermit = permit

	_, err := ExecuteAgentRunnerRun(ctx, fixture.options)
	if !errors.Is(err, gates.ErrProbePermitConsumed) {
		t.Fatalf("a spent probe permit must be refused with the consumed sentinel, got %v", err)
	}
	if errors.Is(err, gates.ErrProbePermitInUse) {
		t.Fatalf("a spent permit must never be reported as in use: %v", err)
	}
	if _, statErr := os.Stat(runEventLogPath(fixture.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a quota refusal must write no run-root marker (stat err=%v)", statErr)
	}
	requireAgentRunnerZeroSpawn(t, fixture.runRoot)
}

// TestSupervisedAssemblyRefusesSpentProbePermit is the supervised twin of
// reverse leg 2.
func TestSupervisedAssemblyRefusesSpentProbePermit(t *testing.T) {
	ctx := context.Background()
	permit := newTestProbePermit(t)
	if err := permit.Require(ctx); err != nil {
		t.Fatalf("the pre-claimed permit fixture must claim: %v", err)
	}
	if err := permit.Complete(); err != nil {
		t.Fatalf("the pre-claimed permit fixture must complete: %v", err)
	}
	harness := newSupervisedHarness(t, []string{"complete"}, false, 60000)
	harness.options.ProbePermit = permit

	_, err := ExecuteSupervisedRun(ctx, harness.options)
	if !errors.Is(err, gates.ErrProbePermitConsumed) {
		t.Fatalf("a spent probe permit must be refused with the consumed sentinel, got %v", err)
	}
	if _, statErr := os.Stat(runEventLogPath(harness.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a quota refusal must write no run-root marker (stat err=%v)", statErr)
	}
}

// TestRunRefusesMissingProbePermitFlag is reverse leg 3a: a fully declared
// AgentRunner invocation that omits --probe-permit is refused fail-closed by the
// CLI presence gate, before any record is read.
func TestRunRefusesMissingProbePermitFlag(t *testing.T) {
	cli := newTestCLI(t.TempDir())
	chainPath := p51WriteChainConfig(t, p51SingleTaskChain()...)
	requestRecord := filepath.Join(t.TempDir(), "request-record.json")

	extra := append(dropCLIFlag(p51SevenAgentRunnerArgs(t), "--probe-permit"), "--request-record", requestRecord)
	code, stdout, stderr := runCLI(t, cli, p51RunArgs(chainPath, extra...)...)
	if code != exitFailure {
		t.Fatalf("run must fail without --probe-permit, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "--probe-permit") {
		t.Fatalf("the refusal must name the missing --probe-permit flag, stderr=%q", stderr)
	}
}

// TestRunRefusesIncompleteRecordsBeforeBurningPermit is reverse leg 3b: an
// incomplete agent-runner input set keeps the pre-existing noop-only refusal and
// never loads the permit, so the one-shot marker stays absent.
func TestRunRefusesIncompleteRecordsBeforeBurningPermit(t *testing.T) {
	cli := newTestCLI(t.TempDir())
	chainPath := p51WriteChainConfig(t, p51SingleTaskChain()...)
	requestRecord := filepath.Join(t.TempDir(), "request-record.json")
	permitPath := filepath.Join(t.TempDir(), "probe-permit.marker")

	extra := append(dropCLIFlag(p51SevenAgentRunnerArgs(t), "--agent-cli"), "--request-record", requestRecord)
	extra = append(dropCLIFlag(extra, "--probe-permit"), "--probe-permit", permitPath)
	code, stdout, stderr := runCLI(t, cli, p51RunArgs(chainPath, extra...)...)
	if code != exitFailure {
		t.Fatalf("run must refuse an incomplete agent-runner input set, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if !strings.Contains(stderr, "local runtime only supports noop steps") {
		t.Fatalf("the refusal must keep the noop-only wording, stderr=%q", stderr)
	}
	if _, statErr := os.Stat(permitPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("an incomplete declaration must not burn the probe permit (stat err=%v)", statErr)
	}
}
