package adapters

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	copilotCLIProbePrompt = "Reply with exactly PONG and nothing else. Do not call any tool."
	copilotCLIProbeReply  = "PONG"
)

var ErrInvalidCopilotCLIProbe = errors.New("invalid Copilot CLI probe")

type CopilotCLIProbeInvocation struct {
	Executable       string
	WorkingDirectory string
	Args             []string
	Environment      map[string][]byte
}

type CopilotCLIProbeExecution struct {
	Outcome               string
	ExitCode              int
	FailureCode           string
	ResultSeen            bool
	ResultExit            int
	ResultUsageSeen       bool
	Reply                 string
	ToolCalls             int
	RequestsUsed          int
	UsageSeen             bool
	UsageUserRequests     int
	UsageProviderRequests int
	Evidence              []string
}

type CopilotCLIProbeExecutor interface {
	Execute(context.Context, CopilotCLIProbeInvocation) CopilotCLIProbeExecution
}

type CopilotCLIProbe struct {
	Executable       string
	WorkingDirectory string
	Executor         CopilotCLIProbeExecutor
}

func (probe CopilotCLIProbe) Probe(ctx context.Context, profile AIProviderProfile, credential []byte, channel string, maximumRequests int) (AIProbeResult, error) {
	if channel != "agent-runner-cli" {
		return AIProbeResult{}, fmt.Errorf("%w: unsupported channel", ErrInvalidCopilotCLIProbe)
	}
	if err := ValidateAIProviderProfile(profile); err != nil || maximumRequests < 1 {
		return AIProbeResult{}, fmt.Errorf("%w: profile or request budget", ErrInvalidCopilotCLIProbe)
	}
	if probe.Executable == "" || probe.WorkingDirectory == "" || probe.Executor == nil {
		return unavailableCopilotCLIProbe("cli_unavailable", 0, nil, CopilotCLIProbeInvocation{}), nil
	}
	if profile.AuthMode == "secret" && len(credential) == 0 {
		return unavailableCopilotCLIProbe("secret_missing", 0, nil, CopilotCLIProbeInvocation{}), nil
	}

	invocation := CopilotCLIProbeInvocation{
		Executable:       probe.Executable,
		WorkingDirectory: probe.WorkingDirectory,
		Args:             copilotCLIProbeArgs(profile.Model, probe.WorkingDirectory, maximumRequests),
		Environment:      copilotCLIProbeEnvironment(profile, credential),
	}
	defer clearCopilotCLIProbeEnvironment(invocation.Environment)
	execution := probe.Executor.Execute(ctx, invocation)
	evidenceHashes := copilotCLIProbeEvidence(invocation, execution)
	requestsUsed := max(execution.RequestsUsed, 0)

	switch execution.Outcome {
	case "start-failed":
		return unavailableCopilotCLIProbe("cli_unavailable", requestsUsed, evidenceHashes, invocation), nil
	case "timed-out":
		return unknownCopilotCLIProbe("timeout", requestsUsed, evidenceHashes), nil
	case "exited":
	default:
		return unknownCopilotCLIProbe("unknown", requestsUsed, evidenceHashes), nil
	}
	if execution.FailureCode != "" {
		if validCopilotCLIProbeFailure(execution.FailureCode) {
			return unavailableCopilotCLIProbe(execution.FailureCode, requestsUsed, evidenceHashes, invocation), nil
		}
		return unknownCopilotCLIProbe("unknown", requestsUsed, evidenceHashes), nil
	}
	if execution.ExitCode != 0 {
		return unknownCopilotCLIProbe("unknown", requestsUsed, evidenceHashes), nil
	}
	if requestsUsed != 1 || requestsUsed > maximumRequests || !execution.UsageSeen || execution.UsageUserRequests != 1 || execution.UsageProviderRequests != 1 || !execution.ResultSeen || !execution.ResultUsageSeen || execution.ResultExit != 0 || execution.ToolCalls != 0 || execution.Reply != copilotCLIProbeReply {
		return unavailableCopilotCLIProbe("response_invalid", requestsUsed, evidenceHashes, invocation), nil
	}
	return AIProbeResult{Status: "available", RequestsUsed: requestsUsed, Evidence: evidenceHashes}, nil
}

func copilotCLIProbeArgs(model, workingDirectory string, maximumRequests int) []string {
	return []string{
		"-p", copilotCLIProbePrompt,
		"-C", workingDirectory,
		"--name", "proofrail-availability-probe",
		"--model", model,
		"--output-format", "json",
		"--stream", "on",
		"--max-ai-credits", strconv.Itoa(maximumRequests),
		"--disable-builtin-mcps",
		"--disallow-temp-dir",
		"--no-custom-instructions",
		"--no-ask-user",
		"--no-auto-update",
		"--no-color",
		"--no-remote",
		"--no-remote-export",
		"--secret-env-vars", "COPILOT_PROVIDER_API_KEY",
		"--deny-url=https://*",
		"--deny-url=http://*",
		"--available-tools=powershell",
		"--deny-tool=shell",
	}
}

func copilotCLIProbeEnvironment(profile AIProviderProfile, credential []byte) map[string][]byte {
	if profile.AuthMode != "secret" {
		return nil
	}
	environment := map[string][]byte{
		"COPILOT_PROVIDER_TYPE":     []byte(profile.ProviderType),
		"COPILOT_PROVIDER_BASE_URL": []byte(profile.BaseURL),
		"COPILOT_PROVIDER_API_KEY":  slices.Clone(credential),
		"COPILOT_MODEL":             []byte(profile.Model),
	}
	if profile.WireAPI != "" {
		environment["COPILOT_PROVIDER_WIRE_API"] = []byte(profile.WireAPI)
	}
	return environment
}

func clearCopilotCLIProbeEnvironment(environment map[string][]byte) {
	for _, value := range environment {
		clear(value)
	}
}

func unavailableCopilotCLIProbe(reason string, requestsUsed int, evidenceHashes []string, invocation CopilotCLIProbeInvocation) AIProbeResult {
	if len(evidenceHashes) == 0 {
		evidenceHashes = copilotCLIProbeEvidence(invocation, CopilotCLIProbeExecution{Outcome: "not-started", FailureCode: reason})
	}
	return AIProbeResult{Status: "unavailable", Reason: reason, RequestsUsed: requestsUsed, Evidence: evidenceHashes}
}

func unknownCopilotCLIProbe(reason string, requestsUsed int, evidenceHashes []string) AIProbeResult {
	return AIProbeResult{Status: "unknown", Reason: reason, RequestsUsed: requestsUsed, Evidence: evidenceHashes}
}

func validCopilotCLIProbeFailure(reason string) bool {
	switch reason {
	case "credential_rejected", "account_unavailable", "quota_exhausted", "network_unreachable", "tls_or_proxy_failure", "provider_unavailable", "model_unavailable", "policy_blocked", "response_invalid":
		return true
	default:
		return false
	}
}

func copilotCLIProbeEvidence(invocation CopilotCLIProbeInvocation, execution CopilotCLIProbeExecution) []string {
	environmentNames := make([]string, 0, len(invocation.Environment))
	for name := range invocation.Environment {
		environmentNames = append(environmentNames, name)
	}
	sort.Strings(environmentNames)
	canonical, _ := evidence.EncodeCanonical(struct {
		Executable       string   `json:"executable"`
		WorkingDirectory string   `json:"workingDirectory"`
		Args             []string `json:"args"`
		EnvironmentNames []string `json:"environmentNames"`
		Outcome          string   `json:"outcome"`
		ExitCode         int      `json:"exitCode"`
		FailureCode      string   `json:"failureCode,omitempty"`
		ResultSeen       bool     `json:"resultSeen"`
		ResultExit       int      `json:"resultExit"`
		ResultUsageSeen  bool     `json:"resultUsageSeen"`
		ReplyMatched     bool     `json:"replyMatched"`
		ToolCalls        int      `json:"toolCalls"`
		RequestsUsed     int      `json:"requestsUsed"`
		UsageSeen        bool     `json:"usageSeen"`
		UsageUser        int      `json:"usageUserRequests"`
		UsageProvider    int      `json:"usageProviderRequests"`
	}{
		Executable:       invocation.Executable,
		WorkingDirectory: invocation.WorkingDirectory,
		Args:             invocation.Args,
		EnvironmentNames: environmentNames,
		Outcome:          execution.Outcome,
		ExitCode:         execution.ExitCode,
		FailureCode:      execution.FailureCode,
		ResultSeen:       execution.ResultSeen,
		ResultExit:       execution.ResultExit,
		ResultUsageSeen:  execution.ResultUsageSeen,
		ReplyMatched:     execution.Reply == copilotCLIProbeReply,
		ToolCalls:        execution.ToolCalls,
		RequestsUsed:     execution.RequestsUsed,
		UsageSeen:        execution.UsageSeen,
		UsageUser:        execution.UsageUserRequests,
		UsageProvider:    execution.UsageProviderRequests,
	})
	hashes := sanitizeHashes(execution.Evidence)
	hashes = append(hashes, evidence.Digest("proofrail:copilot-cli-availability-probe:1\n", canonical))
	sort.Strings(hashes)
	return slices.Compact(hashes)
}
