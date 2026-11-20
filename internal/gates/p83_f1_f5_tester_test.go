package gates

// P83 SECURITY-FINDINGS-REMEDIATION, tester counterexamples for F1 and F5.
//
// Every test in this file runs on a plain Windows host (no symlink, no
// developer mode), so each one is usable as mutation evidence locally.
//
//   - F1: an artifact match that resolves outside the bound workspace is refused
//     and never collected. The relative-escape leg needs no symlink at all, and
//     the unopenable-root leg proves the collector refuses instead of falling back
//     to path-based reads when the workspace handle cannot be opened.
//   - F5: an external-write hook with no declared scope used to skip the
//     admission check entirely (classification == nil). It must never execute.

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// linkDirectory makes link address target without needing symlink privileges:
// Windows uses an NTFS junction (no elevation required), elsewhere a symlink.
// It reports whether the link could be created.
func linkDirectory(link, target string) bool {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "mklink", "/J", link, target).Run() == nil
	}
	return os.Symlink(target, link) == nil
}

// TestRunnerArtifactRefusesJunctionEscape is the F1 counterexample that runs on
// Windows: the workspace holds "sub/link" as a junction to an outside directory,
// so the glob "sub/link/secret" names a real file outside the bound workspace. The
// old path-based Lstat + ReadFile collected it as evidence.
// Mutation check: replacing the bound open with os.ReadFile lets the secret into
// the object store and turns this red.
func TestRunnerArtifactRefusesJunctionEscape(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	outside := filepath.Join(root, "outside")
	if err := os.MkdirAll(filepath.Join(workspace, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	const secret = "SECRET-JUNCTION-ESCAPE"
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}
	if !linkDirectory(filepath.Join(workspace, "sub", "link"), outside) {
		t.Skip("the host cannot create a directory link")
	}

	request := validRequest(workspace)
	request.Hook.ArtifactGlobs = []string{"sub/link/secret"}
	refreshHookDefinition(t, &request)
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	objects := &memoryObjects{}
	runner := testRunner(executor, fakeScanner{}, store)
	runner.Objects = objects

	_, err := runner.Run(context.Background(), request)
	if !errors.Is(err, ErrGateFailed) {
		t.Fatalf("an artifact that resolves outside the workspace must fail the gate, got %v", err)
	}
	if len(store.records) != 1 {
		t.Fatalf("the refusal must persist exactly one result, got %d", len(store.records))
	}
	result := store.records[0].Result
	if result.Assessment == "passed" {
		t.Fatalf("an out-of-workspace artifact must never pass the gate, got %+v", result)
	}
	if len(result.Artifacts) != 0 {
		t.Fatalf("no artifact may be collected for an escaping match, got %+v", result.Artifacts)
	}
	for _, data := range objects.values {
		if strings.Contains(string(data), secret) {
			t.Fatal("the outside secret reached the evidence object store")
		}
	}
}

// TestCollectArtifactsRefusesMatchOutsideBoundRoot is the F1 collector guard,
// always runnable: the glob join produces a path outside the bound workspace and
// the collector must refuse it rather than read it by path. Prepare rejects such
// a pattern earlier, so this drives the collector directly - the boundary the
// finding names - and pins that no fallback to a path read exists.
// Mutation check: restoring the os.Lstat/os.ReadFile pair returns the outside
// bytes and no error, which turns this red.
func TestCollectArtifactsRefusesMatchOutsideBoundRoot(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	outside := filepath.Join(root, "outside")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	const secret = "SECRET-COLLECTOR-ESCAPE"
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}

	request := validRequest(workspace)
	request.Hook.ArtifactGlobs = []string{"../outside/secret"}
	objects := &memoryObjects{}
	runner := testRunner(&fakeExecutor{}, fakeScanner{}, &memoryResults{})
	runner.Objects = objects

	artifacts, err := runner.collectArtifacts(context.Background(), request)
	if !errors.Is(err, ErrArtifactOutOfWorkspace) {
		t.Fatalf("an outside match must be refused as an escape, got artifacts=%+v err=%v", artifacts, err)
	}
	if len(artifacts) != 0 {
		t.Fatalf("an escaped match must collect no artifact, got %+v", artifacts)
	}
	for _, data := range objects.values {
		if strings.Contains(string(data), secret) {
			t.Fatal("the outside secret reached the evidence object store")
		}
	}
}

// TestRunnerArtifactRefusesUnopenableWorkspaceRoot is the second F1 counterexample
// that runs on Windows: the workspace root exists (so Prepare accepts it) but is a
// regular file, so the workspace handle cannot be opened. The collector must
// refuse - the old code globbed to zero matches and reported a pass.
// Mutation check: restoring the path-based Lstat/ReadFile (no bound directory)
// turns this red because the run passes with no artifacts.
func TestRunnerArtifactRefusesUnopenableWorkspaceRoot(t *testing.T) {
	root := t.TempDir()
	rootFile := filepath.Join(root, "workspace-is-a-file")
	if err := os.WriteFile(rootFile, []byte("not a directory"), 0o600); err != nil {
		t.Fatal(err)
	}

	request := validRequest(rootFile)
	request.Hook.ArtifactGlobs = []string{"report.txt"}
	refreshHookDefinition(t, &request)
	store := &memoryResults{}
	_, err := testRunner(&fakeExecutor{executions: []Execution{successfulExecution()}}, fakeScanner{}, store).Run(context.Background(), request)
	if !errors.Is(err, ErrGateFailed) {
		t.Fatalf("an unopenable workspace root must refuse artifact collection, got %v", err)
	}
	result := store.records[0].Result
	if result.FailureKind == nil || *result.FailureKind != "artifact-missing" {
		t.Fatalf("an unopenable root must fail as artifact-missing, got %+v", result)
	}
	if len(result.Artifacts) != 0 {
		t.Fatalf("no artifact may be collected without a bound root, got %+v", result.Artifacts)
	}
}

// TestRunnerRefusesExternalWriteHookWithoutDeclaredScope is the F5 counterexample
// that runs on Windows: the frozen hook policy declares external-write and the
// request carries no scope, which is exactly the "classification == nil" branch
// that used to skip AdmitEffect and execute. The gate must deny without executing.
// Mutation check: restoring "classification != nil &&" in Run executes the
// external-write hook and turns this red.
func TestRunnerRefusesExternalWriteHookWithoutDeclaredScope(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Hook.EffectClass = "external-write"
	request.Hook.RecoveryGuarantee = "reconcile-only"
	request.Effect = nil
	refreshHookDefinition(t, &request)

	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	_, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if !errors.Is(err, ErrGateFailed) {
		t.Fatalf("an external-write hook without a scope must be denied, got %v", err)
	}
	if executor.calls != 0 {
		t.Fatalf("an external-write hook must never execute, calls=%d", executor.calls)
	}
	if len(store.records) != 1 {
		t.Fatalf("the denial must persist exactly one result, got %d", len(store.records))
	}
	result := store.records[0].Result
	if result.Assessment != "failed" || result.FailureKind == nil || *result.FailureKind != "start-failed" {
		t.Fatalf("the denial must persist a start-failed result, got %+v", result)
	}
	if len(result.ErrorEvidence) == 0 {
		t.Fatal("the denial must carry error evidence")
	}
}
