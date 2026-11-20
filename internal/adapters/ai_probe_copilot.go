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

// copilotCLICreditFloor is the candidate CLI's documented minimum for
// --max-ai-credits (B1 probe #5: "Use at least 30 AI credits"). It is a
// candidate/platform fact, not the product request budget: the product budget
// stays the caller's maximumRequests and is still enforced one request later.
const copilotCLICreditFloor = 30

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
	// Measurements is what the executor OBSERVED, kept apart from Evidence so a
	// refusal never has to carry it (T8b, P83 fix round 8). It is a measurement of
	// the run (the validated outcome, exit class and usage counters), never a value
	// the probe refused, and it is deliberately not part of a refusal's evidence.
	Measurements []string
}

type CopilotCLIProbeExecutor interface {
	Execute(context.Context, CopilotCLIProbeInvocation) CopilotCLIProbeExecution
}

type CopilotCLIProbe struct {
	Executable       string
	WorkingDirectory string
	Executor         CopilotCLIProbeExecutor
}

// VerifyReceiverBinding implements AIReceiverBindingVerifier for the Copilot CLI
// probe. F6: the checker calls it before it resolves the credential, so a
// receiver that no longer matches the approved executable is refused while the
// secret is still untouched. An executor without a pinned binding (or without the
// binder capability) verifies nothing here and keeps its own pre-spawn check as
// the last line of defence.
func (probe CopilotCLIProbe) VerifyReceiverBinding() error {
	if probe.Executable == "" || probe.WorkingDirectory == "" {
		return fmt.Errorf("%w: probe executable or working directory missing", ErrInvalidCopilotCLIProbe)
	}
	binder, ok := probe.Executor.(CopilotCLIProbeReceiverBinder)
	if !ok || binder == nil {
		return nil
	}
	return binder.VerifyExecutableBinding(probe.Executable)
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
	requestsUsed := max(execution.RequestsUsed, 0)

	switch execution.Outcome {
	case "start-failed":
		// T8b (P83 fix round 6): a refusal is classified, never described. The
		// generic evidence below digests the executable, the working directory and
		// the full argv, so two refusals of different values (a refused logging
		// flag with two different destinations, two different models, two different
		// workspaces) produced two different records for one classification. A
		// start-failed probe therefore reports the stable classification only, which
		// makes "the same kind of refusal" produce identical evidence.
		//
		// T8b (P83 fix round 7): the refusal result is built directly instead of
		// through unavailableCopilotCLIProbe, whose empty-evidence fallback digests the
		// invocation (executable, working directory, argv). That fallback can no longer
		// reach the final evidence of a refusal, so the value independence holds end to
		// end and not only for the executor markers.
		return unavailableCopilotCLIRefusal("cli_unavailable", requestsUsed, copilotCLIProbeRefusalEvidence(execution)), nil
	case "timed-out":
		return unknownCopilotCLIProbe("timeout", requestsUsed, copilotCLIProbeEvidence(invocation, execution)), nil
	case "exited":
	default:
		return unknownCopilotCLIProbe("unknown", requestsUsed, copilotCLIProbeEvidence(invocation, execution)), nil
	}
	evidenceHashes := copilotCLIProbeEvidence(invocation, execution)
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
		"--max-ai-credits", strconv.Itoa(max(copilotCLICreditFloor, maximumRequests)),
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

// unavailableCopilotCLIRefusal builds a refusal result from exactly the evidence it
// was handed. T8b (P83 fix round 7): the shared unavailableCopilotCLIProbe falls
// back to an invocation digest when the caller passes no evidence, and that digest
// covers the executable, the working directory and the argv - the values a refusal
// must never carry. A refusal therefore has its own constructor, so no path can
// reach the final evidence with a value-dependent digest.
func unavailableCopilotCLIRefusal(reason string, requestsUsed int, evidenceHashes []string) AIProbeResult {
	if len(evidenceHashes) == 0 {
		// A refusal still has to name its class, and the class is value-independent.
		evidenceHashes = []string{evidence.Digest("proofrail:copilot-cli-availability-refusal:1\n", []byte(reason))}
	}
	return AIProbeResult{Status: "unavailable", Reason: reason, RequestsUsed: requestsUsed, Evidence: evidenceHashes}
}

func validCopilotCLIProbeFailure(reason string) bool {
	switch reason {
	case "credential_rejected", "account_unavailable", "quota_exhausted", "network_unreachable", "tls_or_proxy_failure", "provider_unavailable", "model_unavailable", "policy_blocked", "response_invalid":
		return true
	default:
		return false
	}
}

// copilotCLIProbeRefusalEvidence is the evidence of a REFUSED probe. T8b (P83 fix
// round 6): it digests only the stable classification - the outcome and the
// failure code class the executor reported - plus the executor's own classified
// markers. Nothing about the values that were refused (the executable path, the
// working directory, the model, the flag values) enters it, so two refusals of
// the same class are indistinguishable in the record and the record cannot become
// a channel that echoes a refused value.
//
// T8b (P83 fix round 8): "the executor reported" is not trusted. The executor is an
// interface, so its Evidence slice and its FailureCode string are caller-supplied
// data: an arbitrary digest merged here would put an unclassified value into the
// final refusal record, and a free-text failure code would put a value-dependent
// string into it. Both are dropped. The refusal digests exactly two CLOSED
// classifications - the outcome class and the failure-code class, each mapped
// through a fixed table - and nothing else. The executor's measurements stay on
// CopilotCLIProbeExecution.Measurements, which no refusal path reads.
func copilotCLIProbeRefusalEvidence(execution CopilotCLIProbeExecution) []string {
	canonical, err := evidence.EncodeCanonical(struct {
		OutcomeClass string `json:"outcomeClass"`
		FailureClass string `json:"failureClass"`
	}{
		OutcomeClass: copilotCLIProbeOutcomeClass(execution.Outcome),
		FailureClass: copilotCLIProbeFailureClass(execution.FailureCode),
	})
	if err != nil {
		return []string{evidence.Digest("proofrail:copilot-cli-availability-refusal:1\n", []byte("unencodable"))}
	}
	return []string{evidence.Digest("proofrail:copilot-cli-availability-refusal:1\n", canonical)}
}

// copilotCLIProbeOutcomeClass maps an executor outcome onto the closed set of
// outcomes a refusal record may name. T8b (P83 fix round 8): anything outside the
// set is "other", so a caller cannot widen the vocabulary of the record.
func copilotCLIProbeOutcomeClass(outcome string) string {
	switch outcome {
	case "start-failed", "executor-failed", "not-started":
		return outcome
	default:
		return "other"
	}
}

// copilotCLIProbeFailureClass maps a failure code onto the closed set a refusal
// record may name. T8b (P83 fix round 8): a code that is not one of the classified
// values is recorded as "unclassified" and its text never enters the evidence.
func copilotCLIProbeFailureClass(failureCode string) string {
	if validCopilotCLIProbeFailure(failureCode) {
		return failureCode
	}
	return "unclassified"
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
	// T8b (P83 fix round 8): the executor's observation moved to Measurements, which
	// is a separate, clearly-named field. It is folded back in here - and only here,
	// on the non-refusal paths - so a timed-out or unknown probe keeps describing
	// what was validated without a refusal ever carrying it.
	hashes = append(hashes, sanitizeHashes(execution.Measurements)...)
	hashes = append(hashes, evidence.Digest("proofrail:copilot-cli-availability-probe:1\n", canonical))
	sort.Strings(hashes)
	return slices.Compact(hashes)
}
