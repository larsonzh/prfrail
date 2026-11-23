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
	// AuthMode DECLARES the authentication mode of this spawn (P85 F14): "host" or
	// "secret". It is set by the probe, and Execute refuses a declared mode that
	// disagrees with the credential the environment actually carries. An empty value
	// is the compatibility path for callers that predate the declaration.
	AuthMode string
	// FrozenSpawn is the launch template the confinement was PROVEN against before
	// the credential was resolved (P85 F13). Execute instantiates it and refuses
	// unless the spawn it is about to perform equals it element by element. Nil is
	// the compatibility path for an executor that is called directly.
	FrozenSpawn *copilotCLIProbeSpawnSpec
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
	// AuthMode, Model and MaximumRequests are the DECLARED spawn inputs (P85 F13).
	// The checker proves the confinement of a concrete launch template before it
	// resolves a credential, and the only way it can do that is if the probe knows
	// the inputs the spawn will use. They are taken from the same profile and the
	// same flag the checker is called with, and Probe refuses when they disagree with
	// what it is handed. All three empty is the compatibility path for tests and for
	// callers that only use the executor.
	AuthMode        string
	Model           string
	MaximumRequests int
	// frozen is the launch template the pre-Resolve proof built and validated
	// (P85 F13 fix round 3). It is stored on the SAME probe object the checker holds
	// (production passes a *CopilotCLIProbe), so the object that was proven and the
	// object that is spawned are one, not two derivations that happen to agree. Nil
	// is the compatibility path: no proof ran, so Probe builds the template and
	// documents the equivalence argument.
	frozen *copilotCLIProbeSpawnSpec
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

// VerifyImmutableInstallSource implements AIImmutableInstallSourceVerifier for the
// Copilot CLI probe. T3 (finding F6): the checker calls it after the receiver
// binding and before the credential is resolved, so an executable whose
// installation source this user can write, replace or re-permission is refused
// while the secret store is still untouched.
func (probe CopilotCLIProbe) VerifyImmutableInstallSource() error {
	if probe.Executable == "" {
		return fmt.Errorf("%w: probe executable missing", ErrInvalidCopilotCLIProbe)
	}
	_, err := VerifyAIInstallSource(probe.Executable)
	return err
}

// VerifyBYOKConfinement implements AIBYOKConfinementVerifier for the Copilot CLI
// probe. T8a (finding F8): the checker calls it after the receiver binding and the
// install source and strictly before the credential is resolved, so a probe whose
// executor cannot prove the spawn confinement is refused while the secret store is
// still untouched and no child exists.
//
// A missing capability is a refusal, never an admission: the gate must not degrade
// to "assume confined" when the executor cannot answer.
//
// P85 F13 (fix round 3): when the spawn inputs are DECLARED, the launch template is
// built ONCE here, validated, bound to the platform write verdict and stored on the
// probe, so Probe spawns that exact object instead of constructing a second one.
// The receiver is a pointer for exactly that reason: a value receiver would freeze
// a template onto a copy the spawn never reads.
func (probe *CopilotCLIProbe) VerifyBYOKConfinement() error {
	if probe == nil {
		return fmt.Errorf("%w: the probe is missing", ErrAIBYOKConfinement)
	}
	if probe.Executable == "" || probe.WorkingDirectory == "" || probe.Executor == nil {
		return fmt.Errorf("%w: probe executable, working directory or executor missing", ErrAIBYOKConfinement)
	}
	verifier, ok := probe.Executor.(CopilotCLIProbeBYOKConfinementVerifier)
	if !ok || verifier == nil {
		return fmt.Errorf("%w: the executor cannot prove the BYOK spawn confinement", ErrAIBYOKConfinement)
	}
	// The compatibility path declares no spawn inputs, so there is nothing to freeze;
	// the executor's own proof is still required.
	if probe.AuthMode == "" {
		return verifier.VerifyBYOKConfinement()
	}
	spec, err := probe.frozenLaunchTemplate()
	if err != nil {
		return err
	}
	// P85 F13: the frozen template's digest is what the platform verdict is bound
	// to. Every proof below names it, so an answer about one specification can never
	// be read as an answer about another.
	specDigest := copilotCLIProbeSpawnSpecDigest(spec)
	if specDigest == "" {
		return fmt.Errorf("%w: the frozen launch template carries no digest, so no platform verdict can be bound to it", ErrAIBYOKConfinement)
	}
	if bound, ok := probe.Executor.(copilotCLIProbeSpecBoundConfinementVerifier); ok && bound != nil {
		if err := bound.VerifyBYOKConfinementFor(specDigest); err != nil {
			return err
		}
	} else if err := verifier.VerifyBYOKConfinement(); err != nil {
		return err
	}
	// The platform write verdict is part of the frozen specification, so the spawn
	// can only proceed when it re-observes the same verdict for the same
	// specification. Before this, the template and the Sandbox verdict were two
	// independent facts. P85 F18: the layer is read from the executor, so the
	// pre-Resolve verdict and the spawn-time verdict are answers of ONE instance.
	verdict, err := copilotCLIProbePlatformWriteVerdict(context.Background(), copilotCLIProbeSandboxOf(probe.Executor), specDigest)
	if err != nil {
		return err
	}
	spec.PlatformWriteProof = verdict
	probe.frozen = &spec
	return nil
}

// frozenLaunchTemplate builds the launch template from the probe's DECLARED inputs.
// It is the PURE constructor half of F13: VerifyBYOKConfinement calls it once and
// keeps the result, Probe calls it only to re-derive a comparison value (or when no
// proof ran at all).
func (probe CopilotCLIProbe) frozenLaunchTemplate() (copilotCLIProbeSpawnSpec, error) {
	spec, err := copilotCLIProbeLaunchSpec(probe.Executable, probe.WorkingDirectory, probe.AuthMode, probe.Model, probe.MaximumRequests, probe.cliOutputBound())
	if err != nil {
		return copilotCLIProbeSpawnSpec{}, fmt.Errorf("%w: %v", ErrAIBYOKConfinement, err)
	}
	return spec, nil
}

// launchSpecFor resolves the ONE specification the spawn instantiates (F13).
//
// When the pre-Resolve proof froze a template, that object is returned - the spawn
// consumes the template the credential gate approved - and the pure constructor is
// used only to assert that the declared inputs still derive the same template. A
// disagreement is a refusal, never a silent re-freeze. When no proof ran (the
// compatibility path) the constructor's output is the specification; because it is
// a pure function over immutable fields of this probe, and the probe is the same
// value the caller hands to Probe, that output is identical to what a proof would
// have produced.
func (probe CopilotCLIProbe) launchSpecFor(profile AIProviderProfile, maximumRequests int) (*copilotCLIProbeSpawnSpec, error) {
	built, err := copilotCLIProbeLaunchSpec(probe.Executable, probe.WorkingDirectory, profile.AuthMode, profile.Model, maximumRequests, probe.cliOutputBound())
	if err != nil {
		return nil, err
	}
	if probe.frozen == nil {
		return &built, nil
	}
	if !copilotCLIProbeSpawnSpecEqual(*probe.frozen, built) {
		return nil, fmt.Errorf("%w: the spawn no longer matches the launch template proven before the credential was resolved", ErrAIBYOKConfinement)
	}
	return probe.frozen, nil
}

// cliOutputBound reports the output bound the executor will enforce, so the frozen
// template carries the real bound instead of a guess. An executor that cannot answer
// is bounded by the package default, which is what Execute would use as well.
func (probe CopilotCLIProbe) cliOutputBound() int64 {
	if reader, ok := probe.Executor.(copilotCLIProbeOutputBoundReader); ok && reader != nil {
		return reader.CopilotCLIProbeOutputBound()
	}
	return defaultCopilotCLIOutputBytes
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

	// P85 F13: the spawn argv comes from the ONE launch-template constructor, so the
	// argv that is spawned is the argv the pre-Resolve proof validated.
	//
	// The DECLARED inputs must agree with what the checker hands the probe: a probe
	// that froze a template for one model/budget/mode may not spawn another. A
	// disagreement refuses instead of re-freezing silently.
	if probe.AuthMode != "" && (probe.AuthMode != profile.AuthMode || probe.Model != profile.Model || probe.MaximumRequests != maximumRequests) {
		return unavailableCopilotCLIRefusal(copilotCLIProbeBYOKFailure, 0, nil), nil
	}
	spec, err := probe.launchSpecFor(profile, maximumRequests)
	if err != nil {
		// A template that cannot be built cannot be spawned: the confinement class is
		// the honest classification (the argv/envelope is the thing that failed).
		if profile.AuthMode == "secret" {
			return unavailableCopilotCLIRefusal(copilotCLIProbeBYOKFailure, 0, nil), nil
		}
		return unavailableCopilotCLIRefusal("cli_unavailable", 0, nil), nil
	}
	invocation := CopilotCLIProbeInvocation{
		Executable:       probe.Executable,
		WorkingDirectory: probe.WorkingDirectory,
		Args:             spec.Args,
		Environment:      copilotCLIProbeEnvironment(profile, credential),
		AuthMode:         profile.AuthMode,
		FrozenSpawn:      spec,
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
		//
		// T8a (finding F8): an unprovable BYOK confinement is its own closed-set
		// class. It is reported under that class instead of collapsing into
		// cli_unavailable, so an operator can tell "the executable is unavailable"
		// apart from "the credential-bearing child could not be confined".
		if execution.FailureCode == copilotCLIProbeBYOKFailure {
			return unavailableCopilotCLIRefusal(copilotCLIProbeBYOKFailure, requestsUsed, copilotCLIProbeRefusalEvidence(execution)), nil
		}
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

// copilotCLIProbeArgs builds the screened child argv. P85 F13 (`-C` ruling): a
// credential-bearing (BYOK) spawn emits NO working-directory argument - its child
// working directory is decided by ProcessSpec.Dir alone, the executor-owned
// temporary directory - while a host-auth spawn keeps `-C <workingDirectory>`.
// Leaving `-C` in the secret argv pointed the child at the caller's workspace even
// after ProcessSpec.Dir had been moved to the temporary directory.
func copilotCLIProbeArgs(model, workingDirectory string, maximumRequests int, authMode string) []string {
	args := []string{
		"-p", copilotCLIProbePrompt,
	}
	if authMode != "secret" {
		args = append(args, "-C", workingDirectory)
	}
	return append(args,
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
	)
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
	case "credential_rejected", "account_unavailable", "quota_exhausted", "network_unreachable", "tls_or_proxy_failure", "provider_unavailable", "model_unavailable", "policy_blocked", "response_invalid", copilotCLIProbeBYOKFailure:
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
