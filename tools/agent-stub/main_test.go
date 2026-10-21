package main

// These tests pin the three states of the child pid report contract:
//
//   - no descendant -> the report exists and is exactly 0 bytes;
//   - a descendant  -> the report lists the descendant's decimal pid;
//   - unwritable    -> the run fails closed with its dedicated exit code.
//
// A missing report is a defect on the consumer side, so "no descendant" has to be
// a positive statement on disk rather than an absence. The assertions read the
// real process exit status and the real bytes, because the stub is the artifact
// under test: re-implementing its behaviour here would prove nothing.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	// Contract names are spelled out literally instead of reusing the package
	// constants, so that renaming or changing a constant cannot silently move the
	// contract out from under these assertions.
	childPIDReportName = "stub-child.pid"
	ownPIDReportName   = "stub.pid"
	evidenceDirVarName = "PROOFRAIL_EVIDENCE_DIR"

	// childReportWriteExit is the stub's dedicated exit code for a pid report it
	// cannot write. writeOwnPID uses the same code, so the failure case also
	// asserts that stub.pid was written: that is what proves the exit came from
	// the child report and not from the own-pid step.
	childReportWriteExit = 11

	// spawnChildLife is how long the spawned descendant stays alive. It bounds the
	// only wait in this file.
	spawnChildLife = 200 * time.Millisecond
)

var (
	stubBuildOnce sync.Once
	stubBinary    string
	stubBuildErr  error
	stubBuildRoot string
)

// TestMain removes the build directory once, after every test in this binary.
func TestMain(m *testing.M) {
	code := m.Run()
	if stubBuildRoot != "" {
		_ = os.RemoveAll(stubBuildRoot)
	}
	os.Exit(code)
}

// buildStub builds the stub from source once per test binary. The stub is the
// artifact under test, so it is built rather than mocked; an unbuildable package
// skips with the build output instead of failing every assertion for the wrong
// reason.
func buildStub(t *testing.T) string {
	t.Helper()
	stubBuildOnce.Do(func() {
		root, err := os.MkdirTemp("", "proofrail-agent-stub-childpid-")
		if err != nil {
			stubBuildErr = err
			return
		}
		stubBuildRoot = root

		pkgDir, err := os.Getwd()
		if err != nil {
			stubBuildErr = err
			return
		}
		if _, err := os.Stat(filepath.Join(pkgDir, "main.go")); err != nil {
			stubBuildErr = fmt.Errorf("not running in the agent-stub package directory %s: %w", pkgDir, err)
			return
		}

		name := "agent-stub"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		target := filepath.Join(root, name)

		build := exec.Command("go", "build", "-o", target, ".")
		build.Dir = pkgDir
		if output, err := build.CombinedOutput(); err != nil {
			stubBuildErr = fmt.Errorf("go build -o %s . in %s: %v: %s", target, pkgDir, err, output)
			return
		}
		stubBinary = target
	})
	if stubBuildErr != nil {
		t.Skipf("these tests need a buildable agent-stub: %v", stubBuildErr)
	}
	return stubBinary
}

// stubRun is one observed invocation.
type stubRun struct {
	// parentPID is the pid of the stub process this test started. The descendant
	// overwrites stub.pid with its own pid, so the "the report names a different
	// process" assertion can only use the pid the harness observed.
	parentPID int
	// executed reports whether a process was started at all.
	executed bool
	// exitCode is the observed exit status and is meaningful only when executed.
	exitCode int
	// output is stdout+stderr, kept for failure messages.
	output string
}

// runStub runs the built stub against a dedicated evidence directory. The
// evidence directory is passed through the environment of this one process, so no
// state leaks between subtests and t.Parallel stays available.
func runStub(t *testing.T, evidenceDir string, args ...string) stubRun {
	t.Helper()
	cmd := exec.Command(buildStub(t), args...)
	cmd.Env = append(os.Environ(), evidenceDirVarName+"="+evidenceDir)
	output, err := cmd.CombinedOutput()

	run := stubRun{output: string(output)}
	if cmd.Process != nil {
		run.parentPID = cmd.Process.Pid
		run.executed = true
	}
	if err == nil {
		return run
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("could not execute the stub at all: %v (output: %q)", err, run.output)
	}
	run.exitCode = exitErr.ExitCode()
	return run
}

// readReport reads one artifact and fails with what was actually observed.
func readReport(t *testing.T, path string) []byte {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("artifact %s must exist: %v", path, err)
	}
	return body
}

// TestRunWithoutChildWritesEmptyChildPIDReport proves the "no descendant" state
// is written, not merely implied by a missing file.
func TestRunWithoutChildWritesEmptyChildPIDReport(t *testing.T) {
	dir := t.TempDir()

	run := runStub(t, dir, "-mode", "run", "-sleep", "0s")
	if run.exitCode != 0 {
		t.Fatalf("run without a child exited with code %d, want 0 (output: %q)", run.exitCode, run.output)
	}

	reportPath := filepath.Join(dir, childPIDReportName)
	info, err := os.Stat(reportPath)
	if err != nil {
		t.Fatalf("child pid report %s must exist after a childless run, got stat error: %v (evidence dir %s)", reportPath, err, dir)
	}
	if info.IsDir() {
		t.Fatalf("child pid report %s is a directory, want a 0-byte regular file", reportPath)
	}
	if info.Size() != 0 {
		t.Fatalf("child pid report %s has %d bytes, want exactly 0 bytes as the positive statement \"no descendant\"; content: %q",
			reportPath, info.Size(), readReport(t, reportPath))
	}

	// The empty report only means "no child" if the run reached the step before it:
	// an empty report next to a missing stub.pid would be a truncated run.
	ownPIDPath := filepath.Join(dir, ownPIDReportName)
	ownPID := readReport(t, ownPIDPath)
	if len(ownPID) == 0 {
		t.Fatalf("own pid report %s is empty (%d bytes), want a decimal pid written before the child report", ownPIDPath, len(ownPID))
	}
	if _, err := strconv.Atoi(strings.TrimSpace(string(ownPID))); err != nil {
		t.Fatalf("own pid report %s holds %q, want a decimal pid: %v", ownPIDPath, ownPID, err)
	}
}

// TestSpawnWritesChildPIDList proves a real descendant is reported as a pid that
// differs from the stub's own pid.
func TestSpawnWritesChildPIDList(t *testing.T) {
	dir := t.TempDir()

	run := runStub(t, dir, "-mode", "spawn", "-child-life", spawnChildLife.String(), "-sleep", "0s")
	if run.exitCode != 0 {
		t.Fatalf("spawn exited with code %d, want 0 (output: %q)", run.exitCode, run.output)
	}
	if !run.executed || run.parentPID <= 0 {
		t.Fatalf("harness error: the spawned stub's own pid is unknown (executed=%t, pid=%d)", run.executed, run.parentPID)
	}

	reportPath := filepath.Join(dir, childPIDReportName)
	body := readReport(t, reportPath)
	if len(body) == 0 {
		t.Fatalf("child pid report %s is 0 bytes after spawning a child, want at least one decimal pid", reportPath)
	}

	text := strings.TrimSpace(string(body))
	if text == "" {
		t.Fatalf("child pid report %s holds only whitespace (%q), want one decimal pid", reportPath, body)
	}
	lines := strings.Split(text, "\n")
	if len(lines) != 1 {
		t.Fatalf("child pid report %s holds %d pid lines (%q), want exactly one for a single spawned child", reportPath, len(lines), body)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		t.Fatalf("child pid report %s holds %q, want one decimal pid: %v", reportPath, lines[0], err)
	}
	if pid <= 0 {
		t.Fatalf("child pid report %s names pid %d, want a positive pid", reportPath, pid)
	}
	if pid == run.parentPID {
		t.Fatalf("child pid report %s names pid %d, which is the stub's own pid, so it does not describe a descendant", reportPath, pid)
	}

	// The spawned descendant overwrites stub.pid with its own pid, so both
	// artifacts must name the same process; without this cross-check an arbitrary
	// positive pid that is not the parent would still pass the assertions above.
	ownPIDPath := filepath.Join(dir, ownPIDReportName)
	started := time.Now()
	deadline := started.Add(2 * time.Second)
	var ownPIDText string
	for {
		ownPIDBody, readErr := os.ReadFile(ownPIDPath)
		if readErr == nil {
			ownPIDText = strings.TrimSpace(string(ownPIDBody))
			if reported, convErr := strconv.Atoi(ownPIDText); convErr == nil && reported == pid {
				break
			}
		}
		if !time.Now().Before(deadline) {
			t.Fatalf("the child pid report names %d but stub.pid holds %q after %s: both artifacts must name the same spawned descendant",
				pid, ownPIDText, time.Since(started))
		}
		time.Sleep(10 * time.Millisecond)
	}

	// The descendant keeps the evidence directory as its working directory, so let
	// it finish before the harness removes the directory. The child sleeps exactly
	// spawnChildLife after publishing its pid, so waiting that long cannot be too
	// short.
	time.Sleep(spawnChildLife)
}

// TestUnwritableChildPIDReportFailsClosed proves a pid report that cannot be
// written stops the run with the dedicated exit code instead of leaving a
// missing file behind.
func TestUnwritableChildPIDReportFailsClosed(t *testing.T) {
	dir := t.TempDir()

	// A directory at the report path makes the write fail on every supported OS.
	reportPath := filepath.Join(dir, childPIDReportName)
	if err := os.Mkdir(reportPath, 0o755); err != nil {
		t.Fatalf("could not place the obstruction directory %s: %v", reportPath, err)
	}

	run := runStub(t, dir, "-mode", "run", "-sleep", "0s")
	if !run.executed {
		t.Fatalf("harness error: the stub was not executed (output: %q)", run.output)
	}
	if run.exitCode != childReportWriteExit {
		t.Fatalf("platform=%s: exit code = %d, want %d (the stub's dedicated exit code for an unwritable child pid report); output: %q",
			runtime.GOOS, run.exitCode, childReportWriteExit, run.output)
	}

	// Exit code 11 is shared with the own-pid step, so the artifact it must have
	// written first is what places the failure at the child report.
	ownPIDPath := filepath.Join(dir, ownPIDReportName)
	ownPID := readReport(t, ownPIDPath)
	if len(ownPID) == 0 {
		t.Fatalf("own pid report %s is empty (%d bytes) although exit code %d was observed, want the pid written before the failing child report",
			ownPIDPath, len(ownPID), childReportWriteExit)
	}
	if _, err := strconv.Atoi(strings.TrimSpace(string(ownPID))); err != nil {
		t.Fatalf("own pid report %s holds %q, want a decimal pid: %v", ownPIDPath, ownPID, err)
	}
}

// TestSpawnWithUnwritableChildPIDReportKillsChildAndFails pins the failure side of
// the spawn path: an empty report must never sit next to a live child; a failed
// spawn write must leave a gap, not a fabricated fact. So when the real pid cannot
// be written after a successful start, the run must kill the child and exit with
// its dedicated spawn-write code instead of degrading to the "no descendant" state.
func TestSpawnWithUnwritableChildPIDReportKillsChildAndFails(t *testing.T) {
	// spawnChildWriteExit is the stub's exit code for a child pid report that
	// cannot be written after the child started. It is deliberately distinct from
	// childReportWriteExit, so a regression that falls back to the empty report on
	// a live child fails here with 11 (or 0) instead of matching a shared code.
	const spawnChildWriteExit = 6

	dir := t.TempDir()

	// A directory at the report path makes the write fail on every supported OS,
	// and it fails only after the child is already alive because the spawn modes no
	// longer write the empty report up front.
	reportPath := filepath.Join(dir, childPIDReportName)
	if err := os.Mkdir(reportPath, 0o755); err != nil {
		t.Fatalf("could not place the obstruction directory %s: %v", reportPath, err)
	}

	run := runStub(t, dir, "-mode", "spawn", "-child-life", spawnChildLife.String(), "-sleep", "0s")
	if !run.executed {
		t.Fatalf("harness error: the stub was not executed (output: %q)", run.output)
	}
	if run.exitCode != spawnChildWriteExit {
		t.Fatalf("platform=%s: exit code = %d, want %d (dedicated code for a child pid report that cannot be written after a successful start); output: %q",
			runtime.GOOS, run.exitCode, spawnChildWriteExit, run.output)
	}

	// The obstruction must survive: replacing it with a regular file would mean the
	// run wrote a report over the failure instead of reporting it.
	info, err := os.Stat(reportPath)
	if err != nil {
		t.Fatalf("child pid report path %s must still exist: %v", reportPath, err)
	}
	if !info.IsDir() {
		t.Fatalf("child pid report path %s is no longer a directory (mode %v), want the obstruction left untouched", reportPath, info.Mode())
	}

	// The own-pid step runs before the spawn, so its artifact places the failure
	// after the start and proves the run was not truncated earlier.
	ownPIDPath := filepath.Join(dir, ownPIDReportName)
	ownPID := readReport(t, ownPIDPath)
	if len(ownPID) == 0 {
		t.Fatalf("own pid report %s is empty (%d bytes) although exit code %d was observed, want the pid written before the spawn",
			ownPIDPath, len(ownPID), spawnChildWriteExit)
	}
	if _, err := strconv.Atoi(strings.TrimSpace(string(ownPID))); err != nil {
		t.Fatalf("own pid report %s holds %q, want a decimal pid: %v", ownPIDPath, ownPID, err)
	}

	// os.ReadDir returns entries sorted by filename, so the expected list is stated
	// in that order. Exactly the artifacts of a run that reached the spawn may
	// exist: a regular file under any spelling of the report path would be the
	// fabricated "no descendant" fact this contract forbids.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("could not list the evidence directory %s: %v", dir, err)
	}
	want := []string{"events.jsonl", childPIDReportName, ownPIDReportName, "usage.json"}
	got := make([]string, 0, len(entries))
	for _, entry := range entries {
		got = append(got, entry.Name())
	}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("evidence directory %s holds %v, want exactly %v (a regular file at the report path would fabricate \"no descendant\")",
			dir, got, want)
	}
}
