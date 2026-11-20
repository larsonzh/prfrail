package gates

// P83 FIX ROUND 7 counterexamples (items T4a and T4b): the artifact enumeration
// reads its entry names from the bound directory object, and a pattern that
// escapes the bound object is refused instead of being skipped.
//
// Mutation contracts are stated on each test.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// p83Fix7SwapExecutor renames the bound root away and recreates a directory with
// DIFFERENT contents at the same name while the hook "runs". That is the exact
// window the bound enumeration has to survive: the bind happened before Execute,
// the collection happens after it, and the caller's path now addresses a tree the
// reads never touch.
type p83Fix7SwapExecutor struct {
	root         string
	replacements map[string]string
}

func (executor *p83Fix7SwapExecutor) Capabilities() Capabilities { return allCapabilities() }

func (executor *p83Fix7SwapExecutor) Execute(context.Context, Prepared, Hook) Execution {
	if err := os.Rename(executor.root, executor.root+"-original"); err != nil {
		panic("rename the bound root: " + err.Error())
	}
	if err := os.MkdirAll(executor.root, 0o755); err != nil {
		panic("recreate the root: " + err.Error())
	}
	for name, content := range executor.replacements {
		target := filepath.Join(executor.root, name)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			panic("recreate the replacement tree: " + err.Error())
		}
		if err := os.WriteFile(target, []byte(content), 0o600); err != nil {
			panic("write the replacement artifact: " + err.Error())
		}
	}
	return successfulExecution()
}

// TestP83Fix7T4aGlobNamesComeFromTheBoundObject is the T4a counterexample: after
// the root path is renamed away and re-created with different contents, the
// enumeration still lists the ORIGINAL object's entry names and the collection
// hashes the ORIGINAL bytes.
//
// Mutation check: replacing the boundArtifactGlob call in collectBoundArtifacts
// with the path glob `filepath.Glob(filepath.Join(request.WorkspaceRoot, pattern))`
// enumerates the replacement tree, collects "replacement.txt" and turns this red.
func TestP83Fix7T4aGlobNamesComeFromTheBoundObject(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "artifact.txt"), []byte("original bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "sub", "nested.txt"), []byte("original nested"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Direct enumeration half: names, not bytes.
	bound, err := evidence.OpenBoundDir(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = bound.Close() }()
	executor := &p83Fix7SwapExecutor{root: root, replacements: map[string]string{
		"replacement.txt":         "replacement bytes",
		"sub/replacement-sub.txt": "replacement nested",
	}}
	executor.Execute(context.Background(), Prepared{}, Hook{})

	names, err := boundArtifactGlob(bound, "*.txt")
	if err != nil {
		t.Fatalf("enumerating the bound object: %v", err)
	}
	if !slices.Equal(names, []string{"artifact.txt"}) {
		t.Fatalf("the enumeration must list the bound object's entries, got %v", names)
	}
	nested, err := boundArtifactGlob(bound, "sub/*.txt")
	if err != nil {
		t.Fatalf("enumerating the bound subdirectory: %v", err)
	}
	if !slices.Equal(nested, []string{"sub/nested.txt"}) {
		t.Fatalf("the nested enumeration must list the bound object's entries, got %v", nested)
	}

	// End-to-end half: the collected artifact carries the original bytes.
	swapRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(swapRoot, "artifact.txt"), []byte("original bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	request := validRequest(swapRoot)
	request.Hook.ArtifactGlobs = []string{"*.txt"}
	refreshHookDefinition(t, &request)
	store := &memoryResults{}
	objects := &memoryObjects{}
	runner := Runner{
		Executor: &p83Fix7SwapExecutor{root: swapRoot, replacements: map[string]string{"replacement.txt": "replacement bytes"}},
		Scanner:  fakeScanner{},
		Objects:  objects,
		Store:    store,
	}
	record, err := runner.Run(context.Background(), request)
	if err != nil {
		t.Fatalf("the run must collect the bound object's artifact, got %v", err)
	}
	if record.Result.Assessment != "passed" || len(record.Result.Artifacts) != 1 {
		t.Fatalf("exactly one bound artifact must be collected, got %+v", record.Result)
	}
	wanted := evidence.Digest("", []byte("original bytes"))
	if record.Result.Artifacts[0].ObjectHash != wanted {
		t.Fatalf("the collected artifact must hash the bound object's bytes, got %s want %s", record.Result.Artifacts[0].ObjectHash, wanted)
	}
	for _, data := range objects.values {
		if strings.Contains(string(data), "replacement bytes") {
			t.Fatal("a replacement artifact reached the evidence object store")
		}
	}
}

// TestP83Fix7T4bEscapingGlobIsRefused is the T4b counterexample: a pattern that
// leaves the bound object - through ".." or an absolute path - is refused with the
// out-of-workspace error, never silently skipped into an empty match.
//
// Mutation check: deleting the escape refusal in boundArtifactGlob (both the
// cleaned-pattern check and the per-component check) lets the pattern fall through
// to a no-match refusal (ErrArtifactMissing) and turns this red.
func TestP83Fix7T4bEscapingGlobIsRefused(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "artifact.txt"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	bound, err := evidence.OpenBoundDir(root)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = bound.Close() }()

	patterns := []string{
		"..",
		"../outside.txt",
		"sub/../../outside.txt",
		filepath.Join(t.TempDir(), "absolute.txt"),
		"/etc/passwd",
	}
	for _, pattern := range patterns {
		if _, err := boundArtifactGlob(bound, pattern); !errors.Is(err, ErrArtifactOutOfWorkspace) {
			t.Fatalf("pattern %q must be refused as out-of-workspace, got %v", pattern, err)
		}
	}
	// Control: a pattern that stays inside the bound object still enumerates, so the
	// refusals above are the escape check and not a general glob failure.
	names, err := boundArtifactGlob(bound, "artifact.txt")
	if err != nil || !slices.Equal(names, []string{"artifact.txt"}) {
		t.Fatalf("an inside pattern must still enumerate, got %v (err %v)", names, err)
	}
}
