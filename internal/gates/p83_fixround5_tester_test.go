package gates

// P83 FIX ROUND 5 counterexamples (items T4a, T4b, T6a, T6b).
//
// T4: the workspace is bound before the hook runs and a hook without artifacts
// opens nothing. T6: the supplied HookDefinitionHash must be the digest of the
// frozen hook, so a substituted definition never reaches Execute.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// p83Fix5SwapExecutor mutates the workspace inside Execute, which is the window
// between "the root was bound" and "artifacts are collected".
type p83Fix5SwapExecutor struct {
	swap  func()
	calls int
}

func (executor *p83Fix5SwapExecutor) Capabilities() Capabilities { return allCapabilities() }

func (executor *p83Fix5SwapExecutor) Execute(context.Context, Prepared, Hook) Execution {
	executor.calls++
	executor.swap()
	return successfulExecution()
}

// TestP83T4aNoArtifactGlobsNeverOpensTheWorkspace is the T4a counterexample: with
// no artifact globs the runner must not open the workspace root at all, so a root
// that cannot be bound (a regular file) is not a reason to fail.
//
// Mutation check: deleting the `len(request.Hook.ArtifactGlobs) == 0` early return
// in collectBoundArtifacts opens the un-bindable root, fails the run as
// artifact-missing and turns this red.
func TestP83T4aNoArtifactGlobsNeverOpensTheWorkspace(t *testing.T) {
	root := t.TempDir()
	notADirectory := filepath.Join(root, "p83-fix5-regular-file")
	if err := os.WriteFile(notADirectory, []byte("not a directory"), 0o600); err != nil {
		t.Fatal(err)
	}
	request := validRequest(notADirectory)
	request.Hook.ArtifactGlobs = nil
	refreshHookDefinition(t, &request)
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	record, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if err != nil {
		t.Fatalf("a hook without artifacts must not be refused by the workspace binding, got %v", err)
	}
	if executor.calls != 1 {
		t.Fatalf("the hook must execute exactly once, calls=%d", executor.calls)
	}
	if record.Result.Assessment != "passed" {
		t.Fatalf("expected a passed result, got %+v", record.Result)
	}
	if len(record.Result.Artifacts) != 0 {
		t.Fatalf("a hook without globs collects no artifact, got %+v", record.Result.Artifacts)
	}
}

// TestP83T4bRootSwappedDuringHookIsNotFollowed is the T4b counterexample: the
// workspace root is bound BEFORE the hook runs, so a rename-and-replace performed
// by the hook cannot redirect the collection to the replacement tree.
//
// The accepted outcomes are the two the contract allows: the artifact comes from
// the original bound directory, or collection fails closed. An artifact carrying
// the replacement bytes is never acceptable.
//
// Mutation check: moving the evidence.OpenBoundDir call below Execute re-binds the
// replacement directory and collects its bytes, turning this red.
func TestP83T4bRootSwappedDuringHookIsNotFollowed(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "p83-fix5-workspace")
	moved := workspace + "-original"
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	const originalBytes = "P83-FIX5-ORIGINAL-ARTIFACT"
	const replacementBytes = "P83-FIX5-REPLACEMENT-ARTIFACT"
	if err := os.WriteFile(filepath.Join(workspace, "artifact.txt"), []byte(originalBytes), 0o600); err != nil {
		t.Fatal(err)
	}
	request := validRequest(workspace)
	request.Hook.ArtifactGlobs = []string{"artifact.txt"}
	refreshHookDefinition(t, &request)

	var swapErr error
	executor := &p83Fix5SwapExecutor{swap: func() {
		if swapErr = os.Rename(workspace, moved); swapErr != nil {
			return
		}
		if swapErr = os.MkdirAll(workspace, 0o755); swapErr != nil {
			return
		}
		swapErr = os.WriteFile(filepath.Join(workspace, "artifact.txt"), []byte(replacementBytes), 0o600)
	}}
	store := &memoryResults{}
	objects := &memoryObjects{}
	runner := testRunner(executor, fakeScanner{}, store)
	runner.Objects = objects

	_, err := runner.Run(context.Background(), request)
	if swapErr != nil {
		t.Skipf("the host cannot replace the bound workspace during the hook: %v", swapErr)
	}
	if executor.calls != 1 {
		t.Fatalf("the hook must execute exactly once, calls=%d", executor.calls)
	}
	for _, data := range objects.values {
		if strings.Contains(string(data), replacementBytes) {
			t.Fatalf("an artifact from the replacement root must never be collected: %q", data)
		}
	}
	if err != nil {
		// Fail-closed collection is the other allowed answer.
		if !errors.Is(err, ErrGateFailed) {
			t.Fatalf("a refused collection must be a gate failure, got %v", err)
		}
		return
	}
	found := false
	for _, data := range objects.values {
		if strings.Contains(string(data), originalBytes) {
			found = true
		}
	}
	if !found {
		t.Fatalf("the artifact must come from the directory bound before the hook ran, objects=%d", len(objects.values))
	}
}

// TestP83T6aWrongDefinitionHashRefusesBeforeExecute is the T6a counterexample: a
// shape-valid but wrong digest is a contradiction, refused before Prepare and
// before Execute.
//
// Mutation check: deleting the DefinitionHash comparison in Run lets the run
// execute under an unproven definition and turns this red.
func TestP83T6aWrongDefinitionHashRefusesBeforeExecute(t *testing.T) {
	request := validRequest(t.TempDir())
	request.HookDefinitionHash = digestOfString(t, "p83-fix5-shape-valid-but-wrong-definition")
	executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
	store := &memoryResults{}
	_, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
	if !errors.Is(err, ErrEffectClassificationConflict) {
		t.Fatalf("a wrong hook definition hash must be refused as a conflict, got %v", err)
	}
	if executor.calls != 0 {
		t.Fatalf("a contradictory definition must never execute, calls=%d", executor.calls)
	}
	if len(store.records) != 0 {
		t.Fatalf("a contradictory definition must persist nothing, got %+v", store.records)
	}
}

// TestP83T6bHookEditedWithoutRebindingDigestRefuses is the T6b counterexample: a
// hook whose classified fields (the recovery guarantee and the external systems)
// differ from the definition the supplied digest names is refused before Execute.
//
// Mutation check: deleting the DefinitionHash comparison in Run executes the
// substituted definition and turns this red.
func TestP83T6bHookEditedWithoutRebindingDigestRefuses(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*Hook)
	}{
		{
			name:   "a-different-recovery-guarantee",
			mutate: func(hook *Hook) { hook.RecoveryGuarantee = "reconcile-only" },
		},
		{
			name:   "b-different-external-systems",
			mutate: func(hook *Hook) { hook.ExternalSystems = []string{"ci-server-p83-fix5"} },
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			request := validRequest(t.TempDir())
			// The digest is deliberately NOT refreshed: this is the stale-digest shape
			// a caller would present after editing the definition.
			testCase.mutate(&request.Hook)
			executor := &fakeExecutor{executions: []Execution{successfulExecution()}}
			store := &memoryResults{}
			_, err := testRunner(executor, fakeScanner{}, store).Run(context.Background(), request)
			if !errors.Is(err, ErrEffectClassificationConflict) {
				t.Fatalf("an edited definition with a stale digest must be refused, got %v", err)
			}
			if executor.calls != 0 {
				t.Fatalf("a substituted definition must never execute, calls=%d", executor.calls)
			}
		})
	}
}

// digestOfString builds a shape-valid SHA-256 digest of one string.
func digestOfString(t *testing.T, value string) string {
	t.Helper()
	hash, err := DefinitionHash(Hook{ID: value})
	if err != nil {
		t.Fatal(err)
	}
	return hash
}
