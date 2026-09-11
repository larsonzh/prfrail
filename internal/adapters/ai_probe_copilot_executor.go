package adapters

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"runtime"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

const (
	defaultCopilotCLIProbeTimeout = 45 * time.Second
	defaultCopilotCLIOutputBytes  = int64(1024 * 1024)
)

var copilotCLIEnvironmentAllowlist = []string{
	"APPDATA", "COMSPEC", "HOMEDRIVE", "HOMEPATH", "HOME", "HTTPS_PROXY", "HTTP_PROXY",
	"LOCALAPPDATA", "NODE_EXTRA_CA_CERTS", "NO_PROXY", "PROGRAMDATA", "SYSTEMROOT", "TEMP",
	"TMP", "TMPDIR", "USERPROFILE", "WINDIR",
}

type ManagedCopilotCLIProbeExecutor struct {
	Timeout            time.Duration
	OutputBytes        int64
	TemporaryDirectory string
}

func (executor ManagedCopilotCLIProbeExecutor) Execute(ctx context.Context, invocation CopilotCLIProbeInvocation) CopilotCLIProbeExecution {
	if invocation.Executable == "" || invocation.WorkingDirectory == "" {
		return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("invalid-invocation")}}
	}
	temporaryDirectory, err := os.MkdirTemp(executor.TemporaryDirectory, "proofrail-ai-probe-")
	if err != nil {
		return CopilotCLIProbeExecution{Outcome: "executor-failed", Evidence: []string{copilotCLIExecutorEvidence("temporary-directory")}}
	}
	defer os.RemoveAll(temporaryDirectory)
	logsDirectory := temporaryDirectory + string(os.PathSeparator) + "logs"
	if err := os.Mkdir(logsDirectory, 0o700); err != nil {
		return CopilotCLIProbeExecution{Outcome: "executor-failed", Evidence: []string{copilotCLIExecutorEvidence("log-directory")}}
	}
	usagePath := temporaryDirectory + string(os.PathSeparator) + "usage.json"
	args := append(slices.Clone(invocation.Args), "--usage-output-file", usagePath, "--log-dir", logsDirectory)
	environment, environmentBuffers := copilotCLIManagedEnvironment(invocation.Environment)
	defer func() {
		for _, buffer := range environmentBuffers {
			clear(buffer)
		}
		runtime.KeepAlive(environmentBuffers)
	}()

	timeout := executor.Timeout
	if timeout <= 0 {
		timeout = defaultCopilotCLIProbeTimeout
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	outputLimit := executor.OutputBytes
	if outputLimit <= 0 {
		outputLimit = defaultCopilotCLIOutputBytes
	}
	output := newCopilotCLIProbeOutput(outputLimit, cancel)
	processResult, runErr := guard.RunManaged(runCtx, guard.ProcessSpec{
		Command: invocation.Executable,
		Args:    args,
		Dir:     invocation.WorkingDirectory,
		Env:     environment,
		Stdout:  output.stdout,
		Stderr:  output.stderr,
	}, 500*time.Millisecond)

	stdout, stderr := output.bytes()
	defer clear(stdout)
	defer clear(stderr)
	usage, _ := os.ReadFile(usagePath)
	defer clear(usage)
	execution := parseCopilotCLIProbeOutput(stdout, usage)
	execution.Evidence = append(execution.Evidence,
		evidence.Digest("proofrail:copilot-cli-probe-stdout:1\n", stdout),
		evidence.Digest("proofrail:copilot-cli-probe-stderr:1\n", stderr),
	)
	if len(usage) > 0 {
		execution.Evidence = append(execution.Evidence, evidence.Digest("proofrail:copilot-cli-probe-usage:1\n", usage))
	}
	switch {
	case !processResult.Started:
		execution.Outcome = "start-failed"
	case output.exceeded():
		execution.Outcome = "resource-limited"
	case errors.Is(runErr, context.DeadlineExceeded):
		execution.Outcome = "timed-out"
	case errors.Is(runErr, context.Canceled):
		execution.Outcome = "cancelled"
	case errors.Is(runErr, guard.ErrProcessTreeRemained), errors.Is(runErr, guard.ErrTerminationUncertain):
		execution.Outcome = "termination-uncertain"
	default:
		execution.Outcome = "exited"
		execution.ExitCode = processResult.ExitCode
	}
	return execution
}

func copilotCLIManagedEnvironment(overrides map[string][]byte) ([]string, [][]byte) {
	values := make(map[string][]byte, len(copilotCLIEnvironmentAllowlist)+len(overrides))
	for _, name := range copilotCLIEnvironmentAllowlist {
		if value, found := os.LookupEnv(name); found {
			values[name] = []byte(value)
		}
	}
	for name, value := range overrides {
		values[name] = slices.Clone(value)
	}
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	environment := make([]string, 0, len(names))
	buffers := make([][]byte, 0, len(names))
	for _, name := range names {
		entry := make([]byte, 0, len(name)+1+len(values[name]))
		entry = append(entry, name...)
		entry = append(entry, '=')
		entry = append(entry, values[name]...)
		buffers = append(buffers, entry)
		environment = append(environment, unsafe.String(unsafe.SliceData(entry), len(entry)))
		clear(values[name])
	}
	return environment, buffers
}

func parseCopilotCLIProbeOutput(stdout, usageJSON []byte) CopilotCLIProbeExecution {
	execution := CopilotCLIProbeExecution{}
	resultCount := 0
	modelCalls := 0
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	scanner.Buffer(make([]byte, 64*1024), int(defaultCopilotCLIOutputBytes))
	for scanner.Scan() {
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var event struct {
			Type     string          `json:"type"`
			Data     json.RawMessage `json:"data"`
			ExitCode int             `json:"exitCode"`
			Usage    json.RawMessage `json:"usage"`
		}
		if err := json.Unmarshal(line, &event); err != nil || event.Type == "" {
			execution.FailureCode = "response_invalid"
			continue
		}
		if execution.FailureCode == "" {
			execution.FailureCode = copilotCLIEventFailure(line)
		}
		switch event.Type {
		case "model.call_start":
			modelCalls++
		case "assistant.message":
			var data struct {
				Content      string            `json:"content"`
				ToolRequests []json.RawMessage `json:"toolRequests"`
			}
			if err := json.Unmarshal(event.Data, &data); err != nil {
				execution.FailureCode = "response_invalid"
				continue
			}
			execution.Reply = strings.TrimSpace(data.Content)
			execution.ToolCalls += len(data.ToolRequests)
		case "tool.execution_start":
			execution.ToolCalls++
		case "result":
			resultCount++
			execution.ResultSeen = true
			execution.ResultExit = event.ExitCode
			execution.ResultUsageSeen = len(event.Usage) > 0 && string(event.Usage) != "null"
		}
	}
	if scanner.Err() != nil || resultCount > 1 {
		execution.FailureCode = "response_invalid"
	}
	execution.RequestsUsed = modelCalls
	if len(usageJSON) > 0 {
		var usage struct {
			TotalUserRequests int `json:"totalUserRequests"`
			ModelMetrics      map[string]struct {
				Requests struct {
					Count int `json:"count"`
				} `json:"requests"`
			} `json:"modelMetrics"`
		}
		if err := json.Unmarshal(usageJSON, &usage); err != nil {
			execution.FailureCode = "response_invalid"
		} else {
			providerRequests := 0
			for _, metrics := range usage.ModelMetrics {
				providerRequests += metrics.Requests.Count
			}
			execution.UsageSeen = true
			execution.UsageUserRequests = usage.TotalUserRequests
			execution.UsageProviderRequests = providerRequests
			execution.RequestsUsed = max(execution.RequestsUsed, usage.TotalUserRequests, providerRequests)
		}
	}
	return execution
}

func copilotCLIEventFailure(eventJSON []byte) string {
	var value any
	if json.Unmarshal(eventJSON, &value) != nil {
		return "response_invalid"
	}
	return findCopilotCLIEventFailure(value)
}

func findCopilotCLIEventFailure(value any) string {
	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			if reason := findCopilotCLIEventFailure(item); reason != "" {
				return reason
			}
		}
	case map[string]any:
		for _, name := range []string{"code", "statusCode", "status"} {
			if candidate, found := typed[name]; found {
				if reason := classifyCopilotCLIEventFailure(candidate); reason != "" {
					return reason
				}
			}
		}
		for _, child := range typed {
			if reason := findCopilotCLIEventFailure(child); reason != "" {
				return reason
			}
		}
	}
	return ""
}

func classifyCopilotCLIEventFailure(value any) string {
	code := ""
	switch typed := value.(type) {
	case string:
		code = strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(typed, "-", "_"), " ", "_"))
	case float64:
		switch typed {
		case 401:
			code = "401"
		case 403:
			code = "403"
		case 429:
			code = "429"
		case 502, 503, 504:
			code = "503"
		}
	}
	switch code {
	case "401", "unauthorized", "authentication_failed", "invalid_api_key", "invalid_token":
		return "credential_rejected"
	case "403", "account_unavailable", "subscription_required":
		return "account_unavailable"
	case "429", "quota_exhausted", "insufficient_quota", "rate_limit_exceeded":
		return "quota_exhausted"
	case "network_unreachable", "econnrefused", "enotfound", "etimedout":
		return "network_unreachable"
	case "tls_error", "certificate_error", "proxy_error", "tls_or_proxy_failure":
		return "tls_or_proxy_failure"
	case "provider_unavailable", "service_unavailable", "bad_gateway", "gateway_timeout", "503":
		return "provider_unavailable"
	case "model_unavailable", "model_not_found", "unsupported_model":
		return "model_unavailable"
	case "policy_blocked", "policy_denied":
		return "policy_blocked"
	default:
		return ""
	}
}

func copilotCLIExecutorEvidence(category string) string {
	return evidence.Digest("proofrail:copilot-cli-probe-executor:1\n", []byte(category))
}

type copilotCLIProbeOutput struct {
	mu        sync.Mutex
	remaining int64
	overflow  bool
	cancel    context.CancelFunc
	stdout    *copilotCLIProbeWriter
	stderr    *copilotCLIProbeWriter
}

type copilotCLIProbeWriter struct {
	owner  *copilotCLIProbeOutput
	buffer bytes.Buffer
}

func newCopilotCLIProbeOutput(limit int64, cancel context.CancelFunc) *copilotCLIProbeOutput {
	output := &copilotCLIProbeOutput{remaining: limit, cancel: cancel}
	output.stdout = &copilotCLIProbeWriter{owner: output}
	output.stderr = &copilotCLIProbeWriter{owner: output}
	return output
}

func (writer *copilotCLIProbeWriter) Write(data []byte) (int, error) {
	writer.owner.mu.Lock()
	keep := min(int64(len(data)), writer.owner.remaining)
	if keep > 0 {
		_, _ = writer.buffer.Write(data[:keep])
		writer.owner.remaining -= keep
	}
	if keep < int64(len(data)) {
		writer.owner.overflow = true
	}
	overflow := writer.owner.overflow
	writer.owner.mu.Unlock()
	if overflow {
		writer.owner.cancel()
	}
	return len(data), nil
}

func (output *copilotCLIProbeOutput) exceeded() bool {
	output.mu.Lock()
	defer output.mu.Unlock()
	return output.overflow
}

func (output *copilotCLIProbeOutput) bytes() ([]byte, []byte) {
	output.mu.Lock()
	defer output.mu.Unlock()
	stdout := slices.Clone(output.stdout.buffer.Bytes())
	stderr := slices.Clone(output.stderr.buffer.Bytes())
	clear(output.stdout.buffer.Bytes())
	clear(output.stderr.buffer.Bytes())
	output.stdout.buffer.Reset()
	output.stderr.buffer.Reset()
	return stdout, stderr
}
