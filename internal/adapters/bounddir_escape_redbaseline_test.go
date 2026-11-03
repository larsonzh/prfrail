//go:build !windows && !darwin

package adapters

// P53 FIX-GAP-77-SYMLINK-ESCAPE - RED BASELINE VARIANT (runtime gated).
//
// This is the verbatim red criterion that the ③ round recorded before the fix; its
// output is archived in
// docs/validation/evidence/FIX-GAP-77-SYMLINK-ESCAPE/03-tests/red-baseline.md.
//
// Purpose: it is the mutation check. The ⑧ reviewer mutates the production tree
// (for example, degrading openBoundArtifactLocked back to a single
// openat(O_NOFOLLOW)) and replays this test; the escape then has to succeed, the
// returned handle has to yield exactly the 23 outside bytes, and only then does the
// test fail with gap77RegressionMessage. Any other outcome fails with a different
// message - "red baseline not reproduced" when the open is refused, "canary
// mismatch" when the read fails, truncates or returns different bytes - so a broken
// payload or environment can never be mistaken for a reproduced vulnerability.
//
// It is still NOT part of the default suite, but the isolation moved from a private
// build tag to a runtime gate: the file belongs to the ordinary unix build (an
// editor tooling such as gopls resolves it normally - a bespoke tag made the file
// look like "no packages found"), it carries its own test name so it cannot shadow
// the permanent regression TestBoundDirMiddleSymlinkEscapesBound, and it skips
// itself unless PRFRAIL_P53_RED_BASELINE=1. Enable it explicitly:
//
//	PRFRAIL_P53_RED_BASELINE=1 go test -run TestBoundDirRedBaselineMiddleSymlinkEscape -v -count=1 ./internal/adapters/
//
// The retired `-tags p53_red_baseline` form is a no-op now: the old tag no longer
// selects anything, so automation left over from the ③b round has to switch to the
// environment variable above, otherwise it silently runs the green suite instead of
// the mutation check.
//
// On the FIXED tree that run is expected to FAIL with "red baseline not
// reproduced": the vulnerable open no longer succeeds, which is the whole point.
// It only reaches gap77RegressionMessage when the fix has been mutated away.
//
// The payload (canaries, tree builder) is shared with the regression through
// bounddir_escape_payload_unix_test.go, which is compiled unconditionally, so the
// mutation run compares against the same bytes as the green run.

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// redBaselineEnv is the one switch that arms this mutation check. It is read by
// nothing in the production tree: setting it only changes which tests run.
const redBaselineEnv = "PRFRAIL_P53_RED_BASELINE"

// TestBoundDirRedBaselineMiddleSymlinkEscape is the P53 red baseline (design r2 §5),
// armed by PRFRAIL_P53_RED_BASELINE=1 and skipped by every ordinary run.
//
// On the unfixed tree the escape must succeed, the returned handle must yield
// exactly the 23 outside bytes, and only then does the test fail with
// gap77RegressionMessage. Every deviation from that script fails with a different
// message, so a failing setup can never masquerade as "the vulnerability
// reproduced".
func TestBoundDirRedBaselineMiddleSymlinkEscape(t *testing.T) {
	if os.Getenv(redBaselineEnv) != "1" {
		t.Skipf("mutation check only (⑧): set %s=1 to replay the pre-fix red criterion; the default suite has to stay green", redBaselineEnv)
	}
	// The mutation target of design r2 section 2 is linux/amd64: every other unix
	// build answers ENOSYS by constant gate, so replaying the red criterion there
	// would report "red baseline not reproduced" for an environment reason instead
	// of a mutation. gap77SkipUnlessOpenat2 is deliberately not used here: its kernel
	// probe would skip on a fail-closed host, and a skipped mutation check reads the
	// same as a passing one in a report.
	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		t.Skipf("the GAP-77 mutation target is linux/amd64 (design r2 section 2); this build is %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	tree := buildGAP77Tree(t)

	dir, err := OpenBoundDir(tree.boundPath)
	if err != nil {
		t.Fatalf("OpenBoundDir: %v", err)
	}
	defer func() { _ = dir.Close() }()

	attackName := filepath.Join(gap77LinkName, gap77LeafName)
	file, err := OpenBoundDirFile(dir, attackName, 1024)
	if err != nil {
		// The red criterion needs the vulnerable open to succeed. Being refused here
		// means the payload or the environment changed, not that the bug is fixed.
		t.Fatalf("red baseline not reproduced: OpenBoundDirFile(%q) was refused on the vulnerable tree (%v); the red criterion needs the escape to succeed before the fix", attackName, err)
	}
	defer func() { _ = file.Close() }()

	data, truncated, readErr := file.ReadAll()
	if readErr != nil {
		t.Fatalf("canary mismatch: reading through the escaped handle failed (%v); the red baseline needs the outside bytes", readErr)
	}
	if truncated {
		t.Fatal("canary mismatch: the escaped read reported truncation; the 23 byte canary must arrive whole")
	}
	if string(data) != gap77OutsideCanary {
		t.Fatalf("canary mismatch: the escaped handle returned %d bytes %q, want the %d byte canary %q",
			len(data), data, len(gap77OutsideCanary), gap77OutsideCanary)
	}

	t.Fatal(gap77RegressionMessage)
}
