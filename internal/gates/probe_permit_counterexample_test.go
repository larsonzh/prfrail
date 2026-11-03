package gates

// P52-S3 (ACCEPTANCE-SECURITY-COUNTEREXAMPLES): counterexamples for the probe
// permit boundary.
//
// These cases pin the fail-closed edges that a happy-path suite leaves open:
// the nil receiver, the nil context, an out-of-set internal state, the
// non-ErrNotExist read failure that carries no permit sentinel, the empty marker
// path, and a claim that is already durable in another operating-system process.
//
// Falsification rule applied to every case below: remove the corresponding guard
// from internal/gates/probe_permit.go and the case must turn red.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// p52MarkerEnv carries the marker path into the cross-process helper.
const p52MarkerEnv = "P52_MARKER_PATH"

// TestProbePermitRequireNilPermitRejected pins the nil-receiver pre-check: a nil
// permit is a permit-state error, never a panic and never an admission.
func TestProbePermitRequireNilPermitRejected(t *testing.T) {
	var permit *ProbePermit

	err := permit.Require(context.Background())
	if !errors.Is(err, ErrProbePermitMissing) {
		t.Fatalf("nil permit must be refused as %v, got: %v", ErrProbePermitMissing, err)
	}
}

// TestProbePermitRequireNilContextRejected pins the nil-context pre-check
// (probe_permit.go, ctx == nil branch): the caller's intent cannot be known, so
// the request is refused instead of being treated as a live context.
func TestProbePermitRequireNilContextRejected(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "probe-permit.json")
	permit := NewProbePermit("permit-nil-context", marker, 1)

	err := permit.Require(nil)
	if !errors.Is(err, ErrProbePermitMissing) {
		t.Fatalf("nil context must be refused as %v, got: %v", ErrProbePermitMissing, err)
	}

	// The refusal must not have persisted a claim: a nil context is a pre-check,
	// so the quota stays available for a real caller.
	if _, statErr := os.Stat(marker); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a nil-context refusal must not write the claim marker, stat err = %v", statErr)
	}
}

// TestProbePermitCompleteUnknownStateIsConsumed is a white-box case. The closed
// state set is unused=0, claimed=1, spent=2 (probe_permit.go, probePermitState).
// A value outside that set can only come from memory corruption or a bug, and
// Complete must classify it as the fail-closed outcome rather than inventing a
// success: the permit is treated as consumed, never as completable.
//
// Falsification: make Complete's default branch return nil and this case turns
// red.
func TestProbePermitCompleteUnknownStateIsConsumed(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "probe-permit.json")
	permit := NewProbePermit("permit-unknown-state", marker, 1)

	// A legal claim first, so the permit is a genuine, already-claimed permit.
	if err := permit.Require(context.Background()); err != nil {
		t.Fatalf("setup: a fresh permit must be claimable, got: %v", err)
	}

	// 3 is one past probePermitSpent: outside the closed state set.
	permit.state.Store(uint32(3))

	if err := permit.Complete(); !errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("an out-of-set permit state must fail closed as %v, got: %v", ErrProbePermitConsumed, err)
	}
}

// TestProbePermitLoadNonNotExistReadFailure pins the non-uniform I/O boundary
// documented in probe_permit.go: only os.ErrNotExist is "never claimed". A read
// failure for any other reason is returned unchanged, wrapped only with the
// "load probe permit" context, and therefore carries NO permit sentinel.
//
// Contract basis: the sentinel doc block states explicitly that a marker body
// that cannot be read is "reported unchanged and carries no permit sentinel", so
// the absence of all three sentinels is the asserted behaviour - not a gap.
func TestProbePermitLoadNonNotExistReadFailure(t *testing.T) {
	// A directory exists, so this is not os.ErrNotExist, but ReadFile cannot
	// read it as a marker body.
	dirPath := t.TempDir()

	permit, err := LoadProbePermit(dirPath)
	if err == nil {
		t.Fatalf("reading a directory as a marker must fail, got permit %+v", permit)
	}
	if permit != nil {
		t.Fatalf("a failed load must not return a permit, got %+v", permit)
	}
	if errors.Is(err, ErrProbePermitMissing) || errors.Is(err, ErrProbePermitInUse) || errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("a non-ErrNotExist read failure must carry no permit sentinel, got: %v", err)
	}
	if want := "load probe permit"; !strings.Contains(err.Error(), want) {
		t.Fatalf("the read failure must name its operation %q, got: %v", want, err)
	}
}

// TestProbePermitLoadEmptyPathRejected pins the constructor guard: an empty
// marker path has nowhere to persist a claim, so it is refused as a
// permit-state error rather than silently producing an unbounded permit.
func TestProbePermitLoadEmptyPathRejected(t *testing.T) {
	permit, err := LoadProbePermit("")
	if !errors.Is(err, ErrProbePermitMissing) {
		t.Fatalf("an empty marker path must be refused as %v, got: %v", ErrProbePermitMissing, err)
	}
	if permit != nil {
		t.Fatalf("a rejected load must not return a permit, got %+v", permit)
	}
}

// TestProbePermitCrossProcessSingleClaim pins the cross-process half of the
// linearisation point: O_CREATE|O_EXCL makes the marker file itself the claim
// bit, so two operating-system processes racing on one path cannot both win.
//
// The child wins the claim and exits; the parent then loads the marker the child
// made durable and must be refused as consumed. This is the real-process version
// of the in-process race covered by TestProbePermitConcurrentRequireSingleWinner.
func TestProbePermitCrossProcessSingleClaim(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "probe-permit.json")

	child := exec.Command(os.Args[0], "-test.run=TestProbePermitCrossProcessHelper")
	child.Env = append(os.Environ(), p52MarkerEnv+"="+marker)
	// CombinedOutput keeps the child's own diagnostics (its t.Fatalf output goes
	// to stdout under `go test`) available for the failure message below.
	output, err := child.CombinedOutput()
	if err != nil {
		t.Fatalf("the child claim process failed: %v\nchild output:\n%s", err, output)
	}

	permit, loadErr := LoadProbePermit(marker)
	if loadErr != nil {
		t.Fatalf("loading the marker the child claimed must succeed: %v", loadErr)
	}
	if err := permit.Require(context.Background()); !errors.Is(err, ErrProbePermitConsumed) {
		t.Fatalf("a claim already durable in another process must refuse this process as %v, got: %v", ErrProbePermitConsumed, err)
	}
}

// TestProbePermitCrossProcessHelper is the child half of
// TestProbePermitCrossProcessSingleClaim. It is inert in a normal test run
// (skipped without the marker environment variable).
func TestProbePermitCrossProcessHelper(t *testing.T) {
	marker := os.Getenv(p52MarkerEnv)
	if marker == "" {
		t.Skip("helper process only: set " + p52MarkerEnv)
	}

	permit, err := LoadProbePermit(marker)
	if err != nil {
		t.Fatalf("helper load: %v", err)
	}
	if err := permit.Require(context.Background()); err != nil {
		t.Fatalf("helper require on an unclaimed marker must succeed: %v", err)
	}
	// Claim won and durable: report success to the parent through the exit code.
	os.Exit(0)
}
