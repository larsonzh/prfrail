package gates

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func validRequest(root string) Request {
	request := Request{
		RunID: "run-one", TaskID: "task-one", StepID: "step-one", Attempt: 1, ExecutionAttempt: 1, WorkspaceRoot: root,
		Hook: Hook{ID: "go-test", Kind: "test", Runner: ProcessRunner{Type: "process", Executable: "go", Args: []string{"test", "./...; touch injected"}}, EffectClass: "read-only", OnFail: FailurePolicy{Action: "fail-stop"}, Timeout: time.Second, Resources: ResourceLimits{MemoryBytes: 1, OutputBytes: 1024, ProcessCount: 1}, Network: NetworkPolicy{Mode: "deny"}},
	}
	// F5 (P83 fix round): Runner.Run now proves that the supplied digest is the
	// digest of the frozen hook definition, so the fixture binds the real digest.
	hash, err := DefinitionHash(request.Hook)
	if err != nil {
		panic(err)
	}
	request.HookDefinitionHash = hash
	return request
}

// refreshHookDefinition re-binds the request digest after a test mutates the
// frozen hook definition. F5: a stale digest is a contradiction, so every
// mutation of request.Hook must be followed by this call.
func refreshHookDefinition(t *testing.T, request *Request) {
	t.Helper()
	hash, err := DefinitionHash(request.Hook)
	if err != nil {
		t.Fatal(err)
	}
	request.HookDefinitionHash = hash
}

func allCapabilities() Capabilities {
	return Capabilities{Process: true, Timeout: true, OutputLimit: true, MemoryLimit: true, ProcessLimit: true, CPUTimeLimit: true, NetworkDeny: true, NetworkLoopback: true, NetworkAllowlist: true}
}

func TestPreparePreservesArgvAndAllowsOnlySelectedEnvironment(t *testing.T) {
	root := t.TempDir()
	request := validRequest(root)
	request.Hook.EnvAllowlist = []string{"PATH"}
	prepared, err := Prepare(request, allCapabilities(), []string{"SECRET=no", "PATH=tools"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(prepared.Args, request.Hook.Runner.Args) || !reflect.DeepEqual(prepared.Env, []string{"PATH=tools"}) {
		t.Fatalf("prepared command changed argv or leaked env: %+v", prepared)
	}
}

func TestPrepareRejectsUnavailablePolicy(t *testing.T) {
	_, err := Prepare(validRequest(t.TempDir()), Capabilities{Process: true, Timeout: true, OutputLimit: true, MemoryLimit: true, ProcessLimit: true}, nil)
	if !errors.Is(err, ErrPolicyUnavailable) {
		t.Fatalf("got %v", err)
	}
}

func TestPrepareRejectsEscapedAndSymlinkCWD(t *testing.T) {
	root := t.TempDir()
	request := validRequest(root)
	request.Hook.CWD = "../outside"
	if _, err := Prepare(request, allCapabilities(), nil); !errors.Is(err, ErrWorkspaceEscape) {
		t.Fatalf("relative escape: %v", err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "linked")); err == nil {
		request.Hook.CWD = "linked"
		if _, err := Prepare(request, allCapabilities(), nil); !errors.Is(err, ErrWorkspaceEscape) {
			t.Fatalf("symlink escape: %v", err)
		}
	}
}

func TestPrepareRejectsMissingAllowedEnvironment(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Hook.EnvAllowlist = []string{"TOKEN"}
	if _, err := Prepare(request, allCapabilities(), []string{"PATH=tools"}); !errors.Is(err, ErrEnvironmentUnavailable) {
		t.Fatalf("got %v", err)
	}
}

func TestPrepareRejectsMissingEffectClass(t *testing.T) {
	request := validRequest(t.TempDir())
	request.Hook.EffectClass = ""
	if _, err := Prepare(request, allCapabilities(), nil); !errors.Is(err, ErrInvalidHook) {
		t.Fatalf("missing effect class should fail closed, got %v", err)
	}
}
