package adapters

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	// copilotCLIExecutableMaxBytes bounds the executable the probe binds.
	copilotCLIExecutableMaxBytes = int64(64 << 20)
	// copilotCLIExecutableDomain separates the executable digest from every other
	// probe digest.
	copilotCLIExecutableDomain = "proofrail:ai-probe-executable:1\n"
)

// CopilotCLIExecutableHash digests the pinned probe executable through one bound
// handle. F6: the receiver binding covers the exact file that will be executed,
// so the console can refuse a probe whose executable is not the approved one
// before any credential is resolved. A missing or oversized file is an error, not
// a zero digest.
func CopilotCLIExecutableHash(path string) (string, error) {
	file, err := evidence.OpenBoundFile(path, copilotCLIExecutableMaxBytes)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()
	return file.Hash(copilotCLIExecutableDomain)
}

// requireSameExecutableFile proves, immediately before the spawn, that the path
// still names the very object the executor opened and compared with the approved
// digest (T3, finding F6).
//
// P83 FIX ROUND 6 (T3 "identity is not real"): the earlier form compared size,
// modification time and mode by path, and every one of those is a property a
// replacer controls. The proof is now two facts taken from ONE fresh open of the
// path: the object identity (volume serial plus file index on Windows, device
// plus inode on Unix) has to equal the identity derived from the approved handle,
// and the bytes read through that same handle have to digest to the approved
// digest. The identity rejects a substituted object; the digest rejects an
// in-place rewrite of the approved object. An identity that cannot be derived
// refuses - there is no fallback to a weaker comparison.
func requireSameExecutableFile(path, approvedFileID, approvedHash string) error {
	if approvedFileID == "" || approvedHash == "" {
		return fmt.Errorf("%w: the approved executable carries no object identity", ErrCopilotCLIExecutableBinding)
	}
	file, err := evidence.OpenBoundFile(path, copilotCLIExecutableMaxBytes)
	if err != nil {
		return fmt.Errorf("%w: reopen the approved executable: %v", ErrCopilotCLIExecutableBinding, err)
	}
	defer func() { _ = file.Close() }()
	actualFileID, err := file.FileID()
	if err != nil {
		return fmt.Errorf("%w: identify the executable at the spawn path: %v", ErrCopilotCLIExecutableBinding, err)
	}
	if actualFileID != approvedFileID {
		return fmt.Errorf("%w: the executable path no longer names the approved object", ErrCopilotCLIExecutableBinding)
	}
	actualHash, err := file.Hash(copilotCLIExecutableDomain)
	if err != nil {
		return fmt.Errorf("%w: digest the executable at the spawn path: %v", ErrCopilotCLIExecutableBinding, err)
	}
	if actualHash != approvedHash {
		return fmt.Errorf("%w: the approved object no longer holds the approved bytes", ErrCopilotCLIExecutableBinding)
	}
	return nil
}

var copilotCLIEnvironmentAllowlist = []string{
	"APPDATA", "COMSPEC", "HOMEDRIVE", "HOMEPATH", "HOME", "HTTPS_PROXY", "HTTP_PROXY",
	"LOCALAPPDATA", "NODE_EXTRA_CA_CERTS", "NO_PROXY", "PROGRAMDATA", "SYSTEMROOT", "TEMP",
	"TMP", "TMPDIR", "USERPROFILE", "WINDIR",
}

// ErrCopilotCLIExecutableBinding reports an executable that no longer matches the
// pinned receiver binding. It is raised before any credential is materialized and
// again immediately before the spawn (F6).
var ErrCopilotCLIExecutableBinding = errors.New("Copilot CLI executable binding mismatch")

// CopilotCLIProbeReceiverBinder is the optional capability an executor exposes so
// its caller can verify the receiver binding before it resolves a secret (F6).
// An executor without a pinned binding is not a binder and verifies nothing.
type CopilotCLIProbeReceiverBinder interface {
	VerifyExecutableBinding(executable string) error
}

// copilotCLIProbeRefusedLoggingFlags is the closed set of logging destinations a
// caller-supplied probe argument may not name. F8: a child that is handed a log
// destination writes its transcript to disk, and a deferred directory removal
// cannot prove nothing was ever written (CONTRACTS.md:214).
var copilotCLIProbeRefusedLoggingFlags = []string{"log", "log-dir", "logdir", "log-file", "logfile"}

// copilotCLIProbeRefusedLoggingArg reports the first caller-supplied argument
// that names a logging destination, covering the "--flag", "--flag=value" and
// "--flag value" forms: the flag token itself is refused in every form, so the
// separate value form cannot slip through either.
//
// P83 FIX ROUND 4 (F8): a token is only a flag when it begins with "-" or "/"
// after trimming. A bare value is data, so a legitimate argument whose value is
// literally "log" (or "log-dir") is neither refused nor mistaken for a flag.
func copilotCLIProbeRefusedLoggingArg(args []string) (string, bool) {
	for _, arg := range args {
		if _, refused := copilotCLIProbeLoggingFlagClass(arg); refused {
			return arg, true
		}
	}
	return "", false
}

// copilotCLIProbeLoggingFlagClass reports the stable class of an argument that
// names a logging destination: the normalised flag name, never the caller's
// value. T8b (P83 fix round): refusal evidence must not vary with the value, so
// the executor digests this class and two refused values of the same flag
// produce identical evidence.
func copilotCLIProbeLoggingFlagClass(arg string) (string, bool) {
	token := strings.TrimSpace(arg)
	if token == "" || (token[0] != '-' && token[0] != '/') {
		return "", false
	}
	name := token
	if index := strings.IndexByte(name, '='); index >= 0 {
		name = name[:index]
	}
	name = strings.ToLower(strings.TrimLeft(name, "-/"))
	for _, flag := range copilotCLIProbeRefusedLoggingFlags {
		if name == flag {
			return flag, true
		}
	}
	return "", false
}

type ManagedCopilotCLIProbeExecutor struct {
	Timeout            time.Duration
	OutputBytes        int64
	TemporaryDirectory string
	// ExpectedExecutableHash pins the executable the executor may start. F6: when
	// it is set the executor hashes the file again immediately before the spawn
	// and refuses on any mismatch, so a swapped executable never receives the
	// credentials the probe resolved for the approved profile. An empty value is
	// the compatibility path for callers that did not bind an executable.
	//
	// T3 (P83 fix round 6): the value is the digest the operator's approval
	// carried, not a digest recomputed from the path by the caller. The executor
	// derives the identity from the same handle it hashes, so the approved digest
	// and the object identity that guards the spawn describe one open file.
	ExpectedExecutableHash string
}

// VerifyExecutableBinding re-hashes the executable the probe is about to start
// and compares it with the pinned digest. F6: the caller invokes it before it
// resolves any credential, so a swapped executable never receives a secret, and
// Execute invokes it again immediately before the spawn. An executor that pins no
// digest returns nil: that is the compatibility path for callers that never bound
// an executable.
func (executor ManagedCopilotCLIProbeExecutor) VerifyExecutableBinding(executable string) error {
	if executor.ExpectedExecutableHash == "" {
		return nil
	}
	actual, err := CopilotCLIExecutableHash(executable)
	if err != nil {
		return fmt.Errorf("%w: hash the executable: %v", ErrCopilotCLIExecutableBinding, err)
	}
	if actual != executor.ExpectedExecutableHash {
		return fmt.Errorf("%w: the executable is no longer the approved one", ErrCopilotCLIExecutableBinding)
	}
	return nil
}

func (executor ManagedCopilotCLIProbeExecutor) Execute(ctx context.Context, invocation CopilotCLIProbeInvocation) CopilotCLIProbeExecution {
	if invocation.Executable == "" || invocation.WorkingDirectory == "" {
		return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("invalid-invocation")}}
	}
	// F8: a caller-supplied logging destination is refused, not stripped. Silently
	// dropping an argument would run the child with an argv the caller never
	// declared and could not observe, which is exactly the ambiguity this guard
	// exists to remove; a refusal is fail-closed and diagnosable.
	if refused, ok := copilotCLIProbeRefusedLoggingArg(invocation.Args); ok {
		// T8b (P83 fix round): the evidence is the stable flag class and never the
		// caller's token, so the refusal cannot vary with the value that was refused
		// and cannot become a channel that echoes a logging destination.
		class, _ := copilotCLIProbeLoggingFlagClass(refused)
		return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("logging-argument-refused:" + class)}}
	}
	// T3 (P83 fix round 6, "spawn must not silently fall back to a fresh path
	// lookup"): an approved digest names ONE object, and an object is only
	// addressable by a path. A bare command name would make the spawn resolve
	// through PATH - a lookup that happens later, in another namespace, and can
	// select a file the approval never saw - so it is refused instead of being
	// resolved behind the approval's back.
	if executor.ExpectedExecutableHash != "" && !strings.ContainsAny(invocation.Executable, `\/`) {
		return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("executable-path-required")}}
	}
	// F6: the receiver binding is verified before any temporary directory exists and
	// before the process is spawned: a mismatch is a refusal with no footprint. It
	// is also verified again here - after the caller already checked it before
	// resolving the credential - because this is the last line of defence.
	//
	// T3 (finding F6 "bind hash and file identity together"): the approved digest
	// is not re-derived from a path chosen later. ONE bound open yields both the
	// digest and the object identity of the file the executor is about to run, and
	// that identity is re-proven immediately before the spawn. A swap between "hash
	// the path" and "start that path" therefore refuses instead of executing an
	// image nobody approved.
	//
	// An executor that pins no digest is the documented compatibility path (a bare
	// command name resolved by the OS), so it keeps the pre-binding behaviour; every
	// production probe pins ExpectedExecutableHash and therefore takes this path.
	var boundExecutable *evidence.BoundFile
	approvedFileID, approvedExecutableHash := "", ""
	if executor.ExpectedExecutableHash != "" {
		file, openErr := evidence.OpenBoundFile(invocation.Executable, copilotCLIExecutableMaxBytes)
		if openErr != nil {
			return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("executable-unreadable")}}
		}
		boundExecutable = file
		defer func() { _ = boundExecutable.Close() }()
		approvedHash, hashErr := boundExecutable.Hash(copilotCLIExecutableDomain)
		if hashErr != nil {
			return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("executable-unreadable")}}
		}
		if approvedHash != executor.ExpectedExecutableHash {
			return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("executable-binding-mismatch")}}
		}
		// T3: the identity comes from the same handle the digest came from, so the
		// pre-spawn comparison is between the approved object and the object the
		// name resolves to now - never between two independent path resolutions.
		fileID, idErr := boundExecutable.FileID()
		if idErr != nil {
			return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("executable-identity-unavailable")}}
		}
		approvedFileID, approvedExecutableHash = fileID, approvedHash
	}
	temporaryDirectory, err := os.MkdirTemp(executor.TemporaryDirectory, "proofrail-ai-probe-")
	if err != nil {
		return CopilotCLIProbeExecution{Outcome: "executor-failed", Evidence: []string{copilotCLIExecutorEvidence("temporary-directory")}}
	}
	defer os.RemoveAll(temporaryDirectory)
	usagePath := temporaryDirectory + string(os.PathSeparator) + "usage.json"
	args := copilotCLIProbeExecArgs(invocation.Args, usagePath)
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
	// T3 (finding F6): the last check before the spawn proves the path still names
	// the object that was opened, hashed and approved above. os/exec can only start
	// a path, so this object-identity re-check is what ties the approved digest to
	// the file the child is created from. An approved digest without a provable
	// identity never reaches this point, and a path that cannot be re-opened is a
	// refusal rather than a fresh lookup of some other file.
	if boundExecutable != nil {
		if err := requireSameExecutableFile(invocation.Executable, approvedFileID, approvedExecutableHash); err != nil {
			return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("executable-binding-mismatch")}}
		}
	}
	processResult, runErr := guard.RunManaged(runCtx, guard.ProcessSpec{
		Command: invocation.Executable,
		Args:    args,
		Dir:     invocation.WorkingDirectory,
		Env:     environment,
		Stdout:  output.stdout,
		Stderr:  output.stderr,
		// T3 (P83 fix round 7): the approved digest travels with the spawn, so the
		// guard binds the child to the object it verifies itself - a held
		// deny-write/deny-delete handle on Windows, the verified descriptor on
		// Linux - instead of the launcher resolving the path a second time after
		// the checks above ran. An unbound executor keeps the compatibility path.
		ExecutableDigest:       approvedExecutableHash,
		ExecutableDigestDomain: copilotCLIExecutableDomain,
	}, 500*time.Millisecond)

	stdout, stderr := output.bytes()
	defer clear(stdout)
	defer clear(stderr)
	// P42 C3: the usage artifact is read through one bounded handle. The probe
	// keeps its "usage is best effort" semantics: a missing or refused artifact is
	// simply no observation, never a partial one.
	usage, _ := readAgentRunnerArtifactBytes(usagePath, agentRunnerMetadataBound)
	defer clear(usage)
	execution := parseCopilotCLIProbeOutput(stdout, usage)
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
	// P83 FIX ROUND 4 (F8): the child streams are only ever parsed and cleared. A
	// digest over raw stdout/stderr/usage is still a statement about unredacted
	// child output, which is exactly the evidence the contract forbids, so the
	// observation is derived from the validated structure instead.
	//
	// T8b (P83 fix round 8): the observation is a MEASUREMENT of the run, so it is
	// also published on its own field. A refusal reads neither field - its evidence
	// comes from the closed classification alone - so nothing an executor measured
	// can widen a refusal record, while the probe's own evidence set keeps its
	// existing shape for the non-refusal paths.
	observation := copilotCLIProbeObservationEvidence(execution)
	execution.Evidence = append(execution.Evidence, observation)
	execution.Measurements = append(execution.Measurements, observation)
	return execution
}

// copilotCLIProbeObservationEvidence summarizes what the probe validated - the
// outcome, the exit class, the failure classification and the parsed usage
// counters - and never the bytes it saw. F8 (CONTRACTS.md:214): evidence must not
// be generated from an unredacted child stream, so no digest or copy of raw
// stdout, stderr or usage may be persisted.
func copilotCLIProbeObservationEvidence(execution CopilotCLIProbeExecution) string {
	canonical, err := evidence.EncodeCanonical(struct {
		Outcome         string `json:"outcome"`
		ExitCode        int    `json:"exitCode"`
		FailureCode     string `json:"failureCode,omitempty"`
		ResultSeen      bool   `json:"resultSeen"`
		ResultExit      int    `json:"resultExit"`
		ResultUsageSeen bool   `json:"resultUsageSeen"`
		ReplyMatched    bool   `json:"replyMatched"`
		ToolCalls       int    `json:"toolCalls"`
		RequestsUsed    int    `json:"requestsUsed"`
		UsageSeen       bool   `json:"usageSeen"`
		UsageUser       int    `json:"usageUserRequests"`
		UsageProvider   int    `json:"usageProviderRequests"`
	}{
		Outcome:         execution.Outcome,
		ExitCode:        execution.ExitCode,
		FailureCode:     execution.FailureCode,
		ResultSeen:      execution.ResultSeen,
		ResultExit:      execution.ResultExit,
		ResultUsageSeen: execution.ResultUsageSeen,
		// ReplyMatched is the validation verdict, not the reply: the model's text
		// itself never reaches evidence.
		ReplyMatched:  execution.Reply == copilotCLIProbeReply,
		ToolCalls:     execution.ToolCalls,
		RequestsUsed:  execution.RequestsUsed,
		UsageSeen:     execution.UsageSeen,
		UsageUser:     execution.UsageUserRequests,
		UsageProvider: execution.UsageProviderRequests,
	})
	if err != nil {
		return copilotCLIExecutorEvidence("observation-unencodable")
	}
	return evidence.Digest("proofrail:copilot-cli-probe-observation:1\n", canonical)
}

// copilotCLIProbeExecArgs builds the child argv of a probe.
//
// F8 (CONTRACTS.md:214 "API keys/tokens must never reach config, digests, logs,
// errors or evidence"): the probe passes no log destination of any kind. A
// "--log-dir" argument makes the child write its transcript to disk, and a
// deferred directory removal cannot prove nothing was ever written, so the only
// child-writable artifact left is the bounded, structured usage file. Execute
// refuses an invocation whose caller-supplied args name a logging destination
// (see copilotCLIProbeRefusedLoggingArg), so this function may assume the
// screened args; a future CLI that requires a log destination must refuse the
// probe instead of silently re-adding one here.
func copilotCLIProbeExecArgs(invocationArgs []string, usagePath string) []string {
	return append(slices.Clone(invocationArgs), "--usage-output-file", usagePath)
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
