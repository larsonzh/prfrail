package release

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/gates"
)

func TestProductionCodeHasNoDirectNetworkImports(t *testing.T) {
	repositoryRoot := filepath.Join("..", "..")
	allowedExecFiles := map[string]bool{
		filepath.Clean(filepath.Join(repositoryRoot, "internal", "guard", "process.go")):             true,
		filepath.Clean(filepath.Join(repositoryRoot, "internal", "guard", "process_linux.go")):       true,
		filepath.Clean(filepath.Join(repositoryRoot, "internal", "guard", "process_unsupported.go")): true,
		filepath.Clean(filepath.Join(repositoryRoot, "internal", "guard", "process_windows.go")):     true,
		filepath.Clean(filepath.Join(repositoryRoot, "internal", "release", "probe.go")):             true,
	}
	for _, relativeRoot := range []string{"cmd", "internal"} {
		root := filepath.Join(repositoryRoot, relativeRoot)
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
			if err != nil {
				return err
			}
			for _, imported := range file.Imports {
				importPath, err := strconv.Unquote(imported.Path.Value)
				if err != nil {
					return err
				}
				if importPath == "net" || strings.HasPrefix(importPath, "net/") {
					t.Errorf("production runtime must not directly import network package %q: %s", importPath, path)
				}
				if importPath == "os/exec" && !allowedExecFiles[filepath.Clean(path)] {
					t.Errorf("production subprocess entry point is outside the reviewed process boundaries: %s", path)
				}
			}
			ast.Inspect(file, func(node ast.Node) bool {
				literal, ok := node.(*ast.BasicLit)
				if !ok || literal.Kind != token.STRING {
					return true
				}
				value, err := strconv.Unquote(literal.Value)
				if err != nil {
					return true
				}
				normalized := strings.ToLower(value)
				for _, forbidden := range []string{"curl", "wget", "invoke-webrequest", "start-bitstransfer", "git clone", "git fetch"} {
					if normalized == forbidden || strings.Contains(normalized, forbidden+" ") {
						t.Errorf("production runtime contains prohibited network command %q: %s", forbidden, path)
					}
				}
				return true
			})
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestCurrentProcessRunnerFailsClosedForEveryNetworkMode(t *testing.T) {
	capabilities := (gates.GuardExecutor{}).Capabilities()
	if capabilities.NetworkDeny || capabilities.NetworkLoopback || capabilities.NetworkAllowlist {
		t.Fatalf("current process runner must not claim unenforced network isolation: %+v", capabilities)
	}
	capabilities.MemoryLimit = true
	capabilities.ProcessLimit = true
	for _, mode := range []string{"deny", "loopback", "allowlist"} {
		request := gates.Request{
			RunID: "run-one", TaskID: "task-one", StepID: "step-one", Attempt: 1, ExecutionAttempt: 1,
			WorkspaceRoot: t.TempDir(), HookDefinitionHash: "sha256:1111111111111111111111111111111111111111111111111111111111111111",
			Hook: gates.Hook{
				ID: "offline-check", Kind: "test", Runner: gates.ProcessRunner{Type: "process", Executable: "go"},
				EffectClass: "read-only", OnFail: gates.FailurePolicy{Action: "fail-stop"}, Timeout: time.Second,
				Resources: gates.ResourceLimits{MemoryBytes: 1, OutputBytes: 1024, ProcessCount: 1}, Network: gates.NetworkPolicy{Mode: mode},
			},
		}
		if _, err := gates.Prepare(request, capabilities, nil); !errors.Is(err, gates.ErrPolicyUnavailable) {
			t.Errorf("network mode %q must fail closed, got %v", mode, err)
		}
	}
}
