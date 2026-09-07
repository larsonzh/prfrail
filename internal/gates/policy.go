package gates

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrInvalidHook            = errors.New("invalid hook")
	ErrPolicyUnavailable      = errors.New("hook policy enforcement unavailable")
	ErrWorkspaceEscape        = errors.New("hook cwd escapes workspace")
	ErrEnvironmentUnavailable = errors.New("allowed environment variable unavailable")
)

var environmentName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type Prepared struct {
	Executable string
	Args       []string
	CWD        string
	Env        []string
}

func Prepare(request Request, capabilities Capabilities, environment []string) (Prepared, error) {
	if err := validateRequest(request); err != nil {
		return Prepared{}, err
	}
	if err := requireCapabilities(request.Hook, capabilities); err != nil {
		return Prepared{}, err
	}
	cwd, err := resolveCWD(request.WorkspaceRoot, request.Hook.CWD)
	if err != nil {
		return Prepared{}, err
	}
	env, err := selectEnvironment(request.Hook.EnvAllowlist, environment)
	if err != nil {
		return Prepared{}, err
	}
	return Prepared{Executable: request.Hook.Runner.Executable, Args: append([]string(nil), request.Hook.Runner.Args...), CWD: cwd, Env: env}, nil
}

func validateRequest(request Request) error {
	hook := request.Hook
	if !evidence.ValidID(request.RunID) || !evidence.ValidID(request.TaskID) || !evidence.ValidID(request.StepID) || !evidence.ValidID(hook.ID) || !evidence.ValidHash(request.HookDefinitionHash) || request.Attempt < 1 || request.ExecutionAttempt < 1 || request.WorkspaceRoot == "" || hook.Runner.Executable == "" {
		return ErrInvalidHook
	}
	if filepath.IsAbs(hook.Runner.Executable) || strings.ContainsAny(hook.Runner.Executable, "\\\x00") {
		return ErrInvalidHook
	}
	if hook.Runner.Type != "process" {
		return fmt.Errorf("%w: runner %q", ErrPolicyUnavailable, hook.Runner.Type)
	}
	if hook.Timeout <= 0 || hook.Resources.MemoryBytes <= 0 || hook.Resources.OutputBytes <= 0 || hook.Resources.ProcessCount <= 0 {
		return ErrInvalidHook
	}
	switch hook.Kind {
	case "precheck", "build", "test", "verify", "review", "cleanup":
	default:
		return ErrInvalidHook
	}
	switch hook.OnFail.Action {
	case "fail-stop", "warn", "manual":
		if hook.OnFail.MaxAttempts != 0 {
			return ErrInvalidHook
		}
	case "retry":
		if hook.OnFail.MaxAttempts < 2 || hook.OnFail.MaxAttempts > 10 {
			return ErrInvalidHook
		}
	default:
		return ErrInvalidHook
	}
	for _, pattern := range hook.ArtifactGlobs {
		clean := filepath.Clean(filepath.FromSlash(pattern))
		if pattern == "" || filepath.IsAbs(pattern) || strings.ContainsAny(pattern, "\\\x00") || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			return ErrInvalidHook
		}
	}
	return nil
}

func requireCapabilities(hook Hook, capabilities Capabilities) error {
	missing := ""
	switch {
	case !capabilities.Process:
		missing = "process runner"
	case !capabilities.Timeout:
		missing = "timeout"
	case !capabilities.OutputLimit:
		missing = "output limit"
	case !capabilities.MemoryLimit:
		missing = "memory limit"
	case !capabilities.ProcessLimit:
		missing = "process limit"
	case hook.Resources.CPUTime > 0 && !capabilities.CPUTimeLimit:
		missing = "CPU time limit"
	case hook.Network.Mode == "deny" && !capabilities.NetworkDeny:
		missing = "network deny"
	case hook.Network.Mode == "loopback" && !capabilities.NetworkLoopback:
		missing = "loopback network"
	case hook.Network.Mode == "allowlist" && !capabilities.NetworkAllowlist:
		missing = "network allowlist"
	case hook.Network.Mode != "deny" && hook.Network.Mode != "loopback" && hook.Network.Mode != "allowlist":
		return ErrInvalidHook
	}
	if missing != "" {
		return fmt.Errorf("%w: %s", ErrPolicyUnavailable, missing)
	}
	return nil
}

func resolveCWD(root, relative string) (string, error) {
	rootPath, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	rootPath, err = filepath.EvalSymlinks(rootPath)
	if err != nil {
		return "", err
	}
	candidate := rootPath
	if relative != "" {
		if filepath.IsAbs(relative) || filepath.Clean(relative) == ".." || strings.HasPrefix(filepath.Clean(relative), ".."+string(filepath.Separator)) {
			return "", ErrWorkspaceEscape
		}
		candidate = filepath.Join(rootPath, filepath.FromSlash(relative))
	}
	candidate, err = filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(rootPath, candidate)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", ErrWorkspaceEscape
	}
	return candidate, nil
}

func selectEnvironment(allowlist, environment []string) ([]string, error) {
	available := make(map[string]string, len(environment))
	for _, entry := range environment {
		name, _, found := strings.Cut(entry, "=")
		if found {
			available[name] = entry
		}
	}
	selected := make([]string, 0, len(allowlist))
	seen := make(map[string]struct{}, len(allowlist))
	for _, name := range allowlist {
		if !environmentName.MatchString(name) {
			return nil, ErrInvalidHook
		}
		if _, duplicate := seen[name]; duplicate {
			return nil, ErrInvalidHook
		}
		entry, exists := available[name]
		if !exists {
			return nil, fmt.Errorf("%w: %s", ErrEnvironmentUnavailable, name)
		}
		seen[name] = struct{}{}
		selected = append(selected, entry)
	}
	return selected, nil
}

func CurrentEnvironment() []string { return os.Environ() }
