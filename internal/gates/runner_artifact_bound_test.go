package gates

// P83 F1 regression: the artifact collector must read through a handle bound to
// the workspace, not through a path that was checked and then resolved again.
//
// The workspace holds "sub/link" -> an outside directory that holds "secret". The
// glob "sub/link/secret" matches, a path-based Lstat/ReadFile pair follows the
// middle symlink and collects the outside bytes into the evidence object store.
// With the bound handle the middle component cannot be a symlink, so the open is
// refused and the run fails as artifact-missing with zero objects stored.
//
// Mutation check: restoring os.ReadFile in collectArtifacts turns this test red -
// the outside secret is collected and the assessment passes.
//
// The test is intentionally not unix-tagged: both halves of the primitive refuse
// a middle link (openat2 on unix, the handle identity re-check on Windows), so it
// runs wherever a directory symlink can be created and skips with the reason
// where the host denies it.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunnerArtifactRefusesEscapingLink(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	outside := filepath.Join(root, "outside")
	if err := os.MkdirAll(filepath.Join(workspace, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	const secret = "SECRET-OUTSIDE-WORKSPACE"
	if err := os.WriteFile(filepath.Join(outside, "secret"), []byte(secret), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(workspace, "sub", "link")); err != nil {
		t.Skipf("the host cannot create a symlink: %v", err)
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
	if err == nil {
		t.Fatal("an artifact reached through a middle symlink must be refused")
	}
	if !errors.Is(err, ErrGateFailed) {
		t.Fatalf("the refusal must be a gate failure, got %v", err)
	}
	if len(store.records) != 1 {
		t.Fatalf("the refusal must still persist exactly one result, got %d", len(store.records))
	}
	result := store.records[0].Result
	if result.Assessment != "failed" || result.FailureKind == nil || *result.FailureKind != "artifact-missing" {
		t.Fatalf("expected an artifact-missing failure, got %+v", result)
	}
	if len(result.Artifacts) != 0 {
		t.Fatalf("no artifact may be collected, got %+v", result.Artifacts)
	}
	for _, data := range objects.values {
		if strings.Contains(string(data), secret) {
			t.Fatal("the outside secret reached the evidence object store")
		}
	}
}
