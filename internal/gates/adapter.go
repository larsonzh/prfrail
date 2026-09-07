package gates

import (
	"context"
	"fmt"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

type ChainPort struct {
	Runner *Runner
	Hooks  map[string]Hook
}

func (port ChainPort) RunHook(ctx context.Context, request chain.StepRequest) (chain.StepResult, error) {
	hook, exists := port.Hooks[request.Step.ID]
	if !exists || port.Runner == nil || !compatibleKind(request.Step.Kind, hook.Kind) {
		return chain.StepResult{}, fmt.Errorf("%w: hook for step %q", ErrInvalidHook, request.Step.ID)
	}
	hash, err := DefinitionHash(hook)
	if err != nil {
		return chain.StepResult{}, err
	}
	record, runErr := port.Runner.Run(ctx, Request{RunID: request.RunID, TaskID: request.TaskID, StepID: request.Step.ID, Attempt: request.Attempt, ExecutionAttempt: 1, WorkspaceRoot: request.Workspace.Root, HookDefinitionHash: hash, Hook: hook})
	result := chain.StepResult{}
	if evidence.ValidHash(record.ResultHash) {
		result.Evidence = []string{record.ResultHash}
	}
	return result, runErr
}

func DefinitionHash(hook Hook) (string, error) {
	definition := struct {
		ID     string `json:"id"`
		Kind   string `json:"kind"`
		Runner struct {
			Type       string   `json:"type"`
			Executable string   `json:"executable"`
			Args       []string `json:"args,omitempty"`
		} `json:"runner"`
		CWD           string         `json:"cwd,omitempty"`
		EnvAllowlist  []string       `json:"envAllowlist,omitempty"`
		OnFail        map[string]any `json:"onFail"`
		ArtifactGlobs []string       `json:"artifactGlobs,omitempty"`
		TimeoutMS     int64          `json:"timeoutMs"`
		Resources     struct {
			MemoryBytes  int64 `json:"memoryBytes"`
			OutputBytes  int64 `json:"outputBytes"`
			ProcessCount int   `json:"processCount"`
			CPUTimeMS    int64 `json:"cpuTimeMs,omitempty"`
		} `json:"resourceLimits"`
		Network NetworkPolicy `json:"networkPolicy"`
	}{ID: hook.ID, Kind: hook.Kind, CWD: hook.CWD, EnvAllowlist: hook.EnvAllowlist, ArtifactGlobs: hook.ArtifactGlobs, TimeoutMS: hook.Timeout.Milliseconds(), Network: hook.Network}
	definition.Runner.Type, definition.Runner.Executable, definition.Runner.Args = hook.Runner.Type, hook.Runner.Executable, hook.Runner.Args
	definition.OnFail = map[string]any{"action": hook.OnFail.Action}
	if hook.OnFail.Action == "retry" {
		definition.OnFail["maxAttempts"] = hook.OnFail.MaxAttempts
	}
	definition.Resources.MemoryBytes, definition.Resources.OutputBytes, definition.Resources.ProcessCount, definition.Resources.CPUTimeMS = hook.Resources.MemoryBytes, hook.Resources.OutputBytes, hook.Resources.ProcessCount, hook.Resources.CPUTime.Milliseconds()
	canonical, err := evidence.EncodeCanonical(definition)
	if err != nil {
		return "", err
	}
	return evidence.Digest("proofrail:hook-definition:1\n", canonical), nil
}

func compatibleKind(stepKind, hookKind string) bool {
	return stepKind == "build" && hookKind == "build" || stepKind == "verify" && (hookKind == "test" || hookKind == "verify")
}

var _ chain.HookPort = ChainPort{}
