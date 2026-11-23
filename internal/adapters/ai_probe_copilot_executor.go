package adapters

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
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

// CopilotCLIProbeBYOKConfinementVerifier is the optional capability an executor
// exposes so its caller can prove the confinement of the child that will receive
// the BYOK credential BEFORE the credential is resolved (T8a, finding F8).
//
// The method carries no context because the caller runs it from the preflight,
// where no dispatch context exists yet; the platform check it performs is
// stateless.
type CopilotCLIProbeBYOKConfinementVerifier interface {
	VerifyBYOKConfinement() error
}

// copilotCLIProbeSpecBoundConfinementVerifier is the specification-bound form of
// the capability above (P85 F13). An executor that can answer it receives the
// digest of the specification that was FROZEN before the credential was resolved
// and must hand exactly that digest to the platform check, so its verdict is an
// answer about that specification instead of about "a probe in general". The plain
// method stays published so an executor without frozen inputs keeps working; the
// probe prefers this one whenever the executor implements it.
type copilotCLIProbeSpecBoundConfinementVerifier interface {
	VerifyBYOKConfinementFor(specDigest string) error
}

// copilotCLIProbeSandboxSource is implemented by an executor that carries its own
// check layer (P85 F18). The probe reads the layer through this capability, so the
// pre-Resolve verdict and the spawn-time verdict are answers of the SAME instance
// rather than of two independent platform lookups.
type copilotCLIProbeSandboxSource interface {
	CopilotCLIProbeSandbox() guard.Sandbox
}

// copilotCLIProbeSandboxOf reads the check layer an executor carries, or nil when it
// carries none (nil is the platform implementation).
func copilotCLIProbeSandboxOf(executor CopilotCLIProbeExecutor) guard.Sandbox {
	source, ok := executor.(copilotCLIProbeSandboxSource)
	if !ok || source == nil {
		return nil
	}
	return source.CopilotCLIProbeSandbox()
}

// copilotCLIProbeRunManaged is the LAST gate a credential-bearing probe passes
// before a child exists (P85 F18): the managed-process launcher. It is a package
// variable so a test can OBSERVE the start decision without creating a process;
// production never reassigns it, and the assignment sits after every other guard
// in Execute, so a count of one start means every preceding gate admitted.
var copilotCLIProbeRunManaged = guard.RunManaged

// copilotCLIProbeOutputBoundReader is the optional capability an executor exposes so
// the launch template frozen before the credential is resolved carries the SAME
// output bound the spawn will enforce (P85 F13). An executor that cannot answer is
// bounded by the package default, which is what Execute uses as well.
type copilotCLIProbeOutputBoundReader interface {
	CopilotCLIProbeOutputBound() int64
}

var (
	// ErrAIBYOKConfinement reports a BYOK spawn whose confinement could not be
	// proven. It is raised before any credential is resolved and again immediately
	// before the spawn, and the availability record names it by its closed-set
	// class byok_confinement_unproven.
	ErrAIBYOKConfinement = errors.New("AI BYOK spawn confinement is not proven")
)

const (
	// copilotCLIProbeCredentialEnvironmentKey names the one environment entry that
	// carries the BYOK credential. T8a (finding F8): its presence is what makes the
	// confinement proof mandatory for a spawn.
	copilotCLIProbeCredentialEnvironmentKey = "COPILOT_PROVIDER_API_KEY"
	// copilotCLIProbeBYOKFailure is the closed-set failure class of a refused
	// confinement proof. It is the only string a refusal may name, so two refusals
	// of the same class are indistinguishable in the record.
	copilotCLIProbeBYOKFailure = "byok_confinement_unproven"
)

// copilotCLIProbeBYOKEnvironmentKeys is the closed set of environment entries the
// probe may hand a credential-bearing child: exactly the BYOK profile keys built
// by copilotCLIProbeEnvironment.
var copilotCLIProbeBYOKEnvironmentKeys = []string{
	"COPILOT_MODEL",
	"COPILOT_PROVIDER_API_KEY",
	"COPILOT_PROVIDER_BASE_URL",
	"COPILOT_PROVIDER_TYPE",
	"COPILOT_PROVIDER_WIRE_API",
}

// copilotCLIProbeRefusedLoggingFlags is the closed set of logging destinations a
// caller-supplied probe argument may not name. F8: a child that is handed a log
// destination writes its transcript to disk, and a deferred directory removal
// cannot prove nothing was ever written (CONTRACTS.md:214).
var copilotCLIProbeRefusedLoggingFlags = []string{"log", "log-dir", "logdir", "log-file", "logfile"}

// copilotCLIProbeWhitelistedFlags is the closed set of argv flags the probe template
// may carry (P85 F13). A blacklist cannot prove that a flag added by a future CLI
// version is refused, so the envelope is built from an allowlist: an argv token that
// names a flag outside this set refuses the probe before any credential is resolved.
// The existing logging-destination blacklist stays as defence in depth; the two
// gates do not replace each other.
var copilotCLIProbeWhitelistedFlags = []string{
	"-p",
	"-C",
	"--name",
	"--model",
	"--output-format",
	"--stream",
	"--max-ai-credits",
	"--disable-builtin-mcps",
	"--disallow-temp-dir",
	"--no-custom-instructions",
	"--no-ask-user",
	"--no-auto-update",
	"--no-color",
	"--no-remote",
	"--no-remote-export",
	"--secret-env-vars",
	"--deny-url",
	"--available-tools",
	"--deny-tool",
}

// copilotCLIProbeValueFlags is the subset of the allowlist whose value travels in
// the NEXT argv token (or after `=`). Every other allowlisted flag is a switch.
var copilotCLIProbeValueFlags = []string{
	"-p",
	"-C",
	"--name",
	"--model",
	"--output-format",
	"--stream",
	"--max-ai-credits",
	"--secret-env-vars",
	"--deny-url",
	"--available-tools",
	"--deny-tool",
}

// copilotCLIProbeCreditCeiling is the upper bound of `--max-ai-credits`. The floor is
// the candidate CLI's documented minimum; the ceiling only bounds a value class, it
// is not a budget (the budget stays maximumRequests, enforced one request later).
const copilotCLIProbeCreditCeiling = 1 << 20

// copilotCLIProbeSpawnSpec is the FROZEN launch template of one probe (P85 F13): the
// allowlisted argv, the closed environment-key set, the declared authentication mode,
// the working-directory rule (ProcessSpec.Dir is the executor-owned temporary
// directory), the bounded in-memory stream class and the output bound. The checker
// builds and proves it BEFORE the credential is resolved; Execute instantiates it and
// refuses unless the spawn it is about to perform equals the frozen template.
type copilotCLIProbeSpawnSpec struct {
	Executable       string
	WorkingDirectory string
	AuthMode         string
	Args             []string
	// EnvironmentKeys is the closed OUTER environment-key set the template was
	// proven with. A spawn may hand the child a strict subset of it - the wire-API
	// entry is optional, so a profile that does not declare one emits four of the
	// five keys - but never a key outside it (P85 F13 / T8a C2).
	EnvironmentKeys []string
	OutputBound     int64
	// PlatformWriteProof is the digest of the platform write-restriction verdict
	// the pre-Resolve proof observed (P85 F13 fix round 3). It is empty in the pure
	// constructor's output and filled only after the verdict was observed, so a
	// spawn can require the SAME verdict instead of treating "the template was
	// proven" and "the sandbox was allowed" as unrelated facts.
	PlatformWriteProof string
}

// copilotCLIProbeSpawnSpecEqual compares the constructor's pure fields. The
// platform verdict is deliberately excluded: it is not an input of the pure
// constructor, it is the observation attached to one frozen template.
func copilotCLIProbeSpawnSpecEqual(left, right copilotCLIProbeSpawnSpec) bool {
	return left.Executable == right.Executable &&
		left.WorkingDirectory == right.WorkingDirectory &&
		left.AuthMode == right.AuthMode &&
		left.OutputBound == right.OutputBound &&
		slices.Equal(left.Args, right.Args) &&
		slices.Equal(left.EnvironmentKeys, right.EnvironmentKeys)
}

// copilotCLIProbeSpawnSpecDomain separates the frozen launch-template digest from
// every other probe digest (P85 F13).
const copilotCLIProbeSpawnSpecDomain = "proofrail:ai-probe-spawn-spec:1\n"

// copilotCLIProbeSpawnSpecFacts is the canonical field set a specification digest
// is taken over. It is the SAME field set copilotCLIProbeSpawnSpecEqual compares,
// so "the specs are equal" and "the digests are equal" cannot drift apart; the
// platform verdict is excluded because it is an observation ABOUT the template,
// not a part of it.
//
// The facts are canonically encoded (RFC 8785) instead of being joined with
// separators: an argv element or a path is a byte string this layer does not
// control, and a separator-joined form would let two different specifications
// produce one digest.
type copilotCLIProbeSpawnSpecFacts struct {
	Executable       string   `json:"executable"`
	WorkingDirectory string   `json:"workingDirectory"`
	AuthMode         string   `json:"authMode"`
	Args             []string `json:"args"`
	EnvironmentKeys  []string `json:"environmentKeys"`
	OutputBound      int64    `json:"outputBound"`
}

// copilotCLIProbeSpawnSpecDigest is the stable digest of one frozen launch
// template (P85 F13). It is what the platform write verdict is bound to, so two
// specifications on the same platform - two different output bounds, two different
// argv, two different key sets - can never share one verdict or one
// PlatformWriteProof. An unencodable fact set returns "", which every consumer
// treats as "no specification": the BYOK gate refuses an empty digest rather than
// admitting an unbound verdict.
func copilotCLIProbeSpawnSpecDigest(spec copilotCLIProbeSpawnSpec) string {
	encoded, err := evidence.EncodeCanonical(copilotCLIProbeSpawnSpecFacts{
		Executable:       spec.Executable,
		WorkingDirectory: spec.WorkingDirectory,
		AuthMode:         spec.AuthMode,
		Args:             spec.Args,
		EnvironmentKeys:  spec.EnvironmentKeys,
		OutputBound:      spec.OutputBound,
	})
	if err != nil {
		return ""
	}
	return evidence.Digest(copilotCLIProbeSpawnSpecDomain, encoded)
}

// copilotCLIProbePlatformWriteVerdict observes the platform write-restriction
// verdict for a credential-bearing probe and returns its stable digest (P85 F13).
//
// The verdict is bound to the specification it answers for: the check receives the
// frozen template's digest, guard.CheckSandbox refuses a verdict that echoes a
// different one, and the digest below carries that specification digest into the
// proof. "The sandbox was allowed" therefore becomes "the sandbox allowed THIS
// frozen spawn", and a template proven under one specification can never be
// spawned under another.
func copilotCLIProbePlatformWriteVerdict(ctx context.Context, sandbox guard.Sandbox, specDigest string) (string, error) {
	checkSpec := guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: specDigest}
	if err := guard.CheckSandbox(ctx, sandbox, checkSpec); err != nil {
		return "", fmt.Errorf("%w: %v", ErrAIBYOKConfinement, err)
	}
	canonical := runtime.GOOS + "\n" + runtime.GOARCH + "\n" + checkSpec.Reason + "\n" + specDigest + "\n" + guard.VerdictAllowed.String() + "\n"
	return evidence.Digest("proofrail:ai-byok-platform-verdict:1\n", []byte(canonical)), nil
}

// copilotCLIProbeLaunchSpec is the ONE constructor of a probe launch configuration.
// The pre-Resolve confinement proof and the spawn instantiate the same value through
// it, so "the specification that was proven" and "the specification that is spawned"
// cannot drift apart. It is a pure function: no environment lookup, no filesystem.
func copilotCLIProbeLaunchSpec(executable, workingDirectory, authMode, model string, maximumRequests int, outputLimit int64) (copilotCLIProbeSpawnSpec, error) {
	if executable == "" || workingDirectory == "" {
		return copilotCLIProbeSpawnSpec{}, fmt.Errorf("%w: the probe executable or working directory is missing", ErrInvalidCopilotCLIProbe)
	}
	if outputLimit <= 0 {
		return copilotCLIProbeSpawnSpec{}, fmt.Errorf("%w: the in-memory output capture is unbounded", ErrAIBYOKConfinement)
	}
	spec := copilotCLIProbeSpawnSpec{
		Executable:       executable,
		WorkingDirectory: workingDirectory,
		AuthMode:         authMode,
		OutputBound:      outputLimit,
	}
	switch authMode {
	case "secret":
		spec.EnvironmentKeys = slices.Clone(copilotCLIProbeBYOKEnvironmentKeys)
		sort.Strings(spec.EnvironmentKeys)
	case "host":
		// Host authentication hands the child no credential of ours, so the closed
		// key set is empty.
		spec.EnvironmentKeys = nil
	default:
		return copilotCLIProbeSpawnSpec{}, fmt.Errorf("%w: the probe authentication mode %q is not declared", ErrAIBYOKConfinement, authMode)
	}
	spec.Args = copilotCLIProbeArgs(model, workingDirectory, maximumRequests, authMode)
	if err := copilotCLIProbeValidateWhitelistedArgs(spec.Args, authMode); err != nil {
		return copilotCLIProbeSpawnSpec{}, err
	}
	return spec, nil
}

// copilotCLIProbeValidateWhitelistedArgs applies the argv allowlist (P85 F13): every
// flag token must be allowlisted, every value is checked against its class, a switch
// carries no value, and a credential-bearing argv may not name a working directory
// (`-C` / `--working-directory`), because the only directory a BYOK child may work in
// is ProcessSpec.Dir - the executor-owned temporary directory.
func copilotCLIProbeValidateWhitelistedArgs(args []string, authMode string) error {
	for index := 0; index < len(args); index++ {
		name, value, inline := copilotCLIProbeSplitFlag(args[index])
		if name == "" {
			return fmt.Errorf("%w: the probe argv carries a bare token that is not a flag value", ErrAIBYOKConfinement)
		}
		if !slices.Contains(copilotCLIProbeWhitelistedFlags, name) {
			return fmt.Errorf("%w: the probe argv names a flag outside the allowlist", ErrAIBYOKConfinement)
		}
		if authMode == "secret" && (name == "-C" || name == "--working-directory") {
			return fmt.Errorf("%w: the credential-bearing argv names a working directory", ErrAIBYOKConfinement)
		}
		if !slices.Contains(copilotCLIProbeValueFlags, name) {
			if inline {
				return fmt.Errorf("%w: the probe argv gives a value to a switch flag", ErrAIBYOKConfinement)
			}
			continue
		}
		if !inline {
			if index+1 >= len(args) {
				return fmt.Errorf("%w: the probe argv names a flag without a value", ErrAIBYOKConfinement)
			}
			value = args[index+1]
			index++
		}
		if err := copilotCLIProbeValidateFlagValue(name, value); err != nil {
			return err
		}
	}
	return nil
}

// copilotCLIProbeSplitFlag splits one argv token into its flag name and inline value.
// A token that does not begin with "-" or "/" is data, not a flag.
func copilotCLIProbeSplitFlag(arg string) (name, value string, inline bool) {
	token := strings.TrimSpace(arg)
	if token == "" || (token[0] != '-' && token[0] != '/') {
		return "", "", false
	}
	if index := strings.IndexByte(token, '='); index >= 0 {
		return token[:index], token[index+1:], true
	}
	return token, "", false
}

// copilotCLIProbeValidateFlagValue applies the value class of one allowlisted flag.
// Every class is CLOSED: a value that is not recognised refuses, it is never passed
// through as "probably fine".
func copilotCLIProbeValidateFlagValue(name, value string) error {
	switch name {
	case "--model":
		if !copilotCLIProbeModelToken(value) {
			return fmt.Errorf("%w: the probe model value is outside its class", ErrAIBYOKConfinement)
		}
	case "--max-ai-credits":
		credits, err := strconv.Atoi(value)
		if err != nil || credits < copilotCLICreditFloor || credits > copilotCLIProbeCreditCeiling {
			return fmt.Errorf("%w: the probe credit value is outside its class", ErrAIBYOKConfinement)
		}
	case "--deny-url":
		if value != "https://*" && value != "http://*" {
			return fmt.Errorf("%w: the probe url rule is outside its class", ErrAIBYOKConfinement)
		}
	case "--output-format":
		if value != "json" {
			return fmt.Errorf("%w: the probe output format is outside its class", ErrAIBYOKConfinement)
		}
	case "--stream":
		if value != "on" {
			return fmt.Errorf("%w: the probe stream mode is outside its class", ErrAIBYOKConfinement)
		}
	case "--name", "--available-tools", "--deny-tool":
		if !evidence.ValidID(value) {
			return fmt.Errorf("%w: the probe %s value is outside its class", ErrAIBYOKConfinement, name)
		}
	case "--secret-env-vars":
		if !slices.Contains(copilotCLIProbeBYOKEnvironmentKeys, value) {
			return fmt.Errorf("%w: the probe secret variable is outside its class", ErrAIBYOKConfinement)
		}
	case "-p":
		if value != copilotCLIProbePrompt {
			return fmt.Errorf("%w: the probe prompt is outside its class", ErrAIBYOKConfinement)
		}
	case "-C":
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%w: the probe working directory value is empty", ErrAIBYOKConfinement)
		}
	}
	return nil
}

// copilotCLIProbeModelToken reports whether one model value stays inside its closed
// character class. The design named evidence.ValidID, which is the lowercase kebab
// form; real CLI model ids are dotted and sometimes mixed case (`gpt-5.6-luna`,
// `claude-sonnet-4.5`), so the class is the same closed alphabet extended with `.`
// and upper case. It still excludes whitespace, quoting, `=` and every shell
// metacharacter, which is what the class exists for.
func copilotCLIProbeModelToken(value string) bool {
	if value == "" || len(value) > 128 || value[0] == '-' {
		return false
	}
	for _, character := range value {
		switch {
		case character >= 'a' && character <= 'z':
		case character >= 'A' && character <= 'Z':
		case character >= '0' && character <= '9':
		case character == '-' || character == '_' || character == '.' || character == ':':
		default:
			return false
		}
	}
	return true
}

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
	// ApprovedInstallSource is the install-source FACT the admitted carrier records
	// (P85 F05): executable path, digest, object identity, platform and immutability
	// facts, plus the interpreter identity when the executable is a script or a
	// dynamic ELF. When it is set, the pre-spawn observation must MATCH IT ITEM BY
	// ITEM, and the object identity / digest that guard the spawn are taken from it
	// instead of being re-derived by this call. Nil is the compatibility path for
	// callers that never admitted a carrier.
	ApprovedInstallSource *AIInstallSourceFacts
	// Sandbox is the write-restriction check layer this executor answers the BYOK
	// gate with (P85 F18). Nil keeps the platform implementation - the production
	// shape, where no caller injects a check layer. When it is set, EVERY gate of
	// this executor uses this one instance, so the verdict the pre-Resolve proof
	// froze and the verdict the spawn re-observes are answers of the same layer.
	Sandbox guard.Sandbox
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

// CopilotCLIProbeOutputBound implements copilotCLIProbeOutputBoundReader (P85 F13).
func (executor ManagedCopilotCLIProbeExecutor) CopilotCLIProbeOutputBound() int64 {
	if executor.OutputBytes > 0 {
		return executor.OutputBytes
	}
	return defaultCopilotCLIOutputBytes
}

// VerifyBYOKConfinement implements CopilotCLIProbeBYOKConfinementVerifier. T8a
// (finding F8): the checker calls it before it resolves a credential, so a probe
// whose confinement cannot be proven never reads the secret store and never
// creates a child.
//
// The proof has two halves that share one implementation. This method owns the
// half that does not depend on the invocation - the static spawn envelope
// (inherited environment disjoint from the credential keys, a non-empty
// logging-destination class, a bounded in-memory capture) plus the platform
// write-restriction verdict. The half that depends on the concrete argv, working
// directory and stream set is asserted again in Execute immediately before the
// spawn, where those values exist. Both halves are fail-closed: an Unknown
// sandbox verdict refuses instead of admitting.
func (executor ManagedCopilotCLIProbeExecutor) VerifyBYOKConfinement() error {
	// P85 F13: the plain capability carries NO specification. It stays published for
	// executors called without frozen inputs, and it refuses here - the probe uses
	// VerifyBYOKConfinementFor whenever it has a frozen template, and a platform
	// verdict that answers for no specification is not a proof of anything.
	return executor.VerifyBYOKConfinementFor("")
}

// VerifyBYOKConfinementFor is the specification-bound half (P85 F13). It runs the
// same static envelope check and then asks the SAME check layer this executor will
// use at spawn time whether THIS specification may hand a credential to a child.
func (executor ManagedCopilotCLIProbeExecutor) VerifyBYOKConfinementFor(specDigest string) error {
	if err := copilotCLIProbeBYOKStaticConfinement(); err != nil {
		return err
	}
	if err := guard.CheckSandbox(context.Background(), executor.Sandbox, guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: specDigest}); err != nil {
		return fmt.Errorf("%w: %v", ErrAIBYOKConfinement, err)
	}
	return nil
}

// CopilotCLIProbeSandbox implements copilotCLIProbeSandboxSource (P85 F18).
func (executor ManagedCopilotCLIProbeExecutor) CopilotCLIProbeSandbox() guard.Sandbox {
	return executor.Sandbox
}

// copilotCLIProbeBYOKStaticConfinement asserts the parts of the BYOK confinement
// proof that do not depend on one invocation. They are facts about the executor's
// own envelope, so they hold for every spawn it can construct.
func copilotCLIProbeBYOKStaticConfinement() error {
	for _, name := range copilotCLIEnvironmentAllowlist {
		if slices.Contains(copilotCLIProbeBYOKEnvironmentKeys, name) {
			return fmt.Errorf("%w: the inherited environment allowlist admits a credential key", ErrAIBYOKConfinement)
		}
	}
	if len(copilotCLIProbeRefusedLoggingFlags) == 0 {
		return fmt.Errorf("%w: the logging-destination class is empty, so argv cannot be screened", ErrAIBYOKConfinement)
	}
	if defaultCopilotCLIOutputBytes <= 0 {
		return fmt.Errorf("%w: the in-memory output capture is unbounded", ErrAIBYOKConfinement)
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
	switch {
	case executor.ApprovedInstallSource != nil:
		// P85 F05: the admitted carrier - not a re-derivation from the path - names
		// the object the operator approved. Its identity and digest are what
		// requireSameExecutableFile re-proves immediately before the spawn, and its
		// full fact set (path, platform, immutability, interpreter) is compared with
		// the live observation further down, as close to the spawn as possible.
		carrier := *executor.ApprovedInstallSource
		if !strings.ContainsAny(invocation.Executable, `\/`) {
			return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("executable-path-required")}}
		}
		if carrier.FileID == "" || carrier.ExecutableHash == "" {
			return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("executable-identity-unavailable")}}
		}
		if executor.ExpectedExecutableHash != "" && carrier.ExecutableHash != executor.ExpectedExecutableHash {
			return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("executable-binding-mismatch")}}
		}
		approvedFileID, approvedExecutableHash = carrier.FileID, carrier.ExecutableHash
	case executor.ExpectedExecutableHash != "":
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
	environment, environmentBuffers, environmentErr := copilotCLIManagedEnvironment(invocation.Environment)
	if environmentErr != nil {
		return CopilotCLIProbeExecution{Outcome: "start-failed", FailureCode: copilotCLIProbeBYOKFailure, Evidence: []string{copilotCLIExecutorEvidence("byok-confinement-unproven")}}
	}
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
	// T3 (finding F6, P84-b): the identity-and-digest binding above proves that the
	// object is the approved one at this instant. It does not prove that this user
	// cannot replace the object between two checks (Windows) or rewrite its bytes in
	// place while keeping the inode (Linux), so the installation source is
	// re-verified live immediately before the spawn. An unprovable source refuses
	// the start; the evidence stays a stable classification and never echoes the
	// path or the observed permission bits.
	if live, err := VerifyAIInstallSource(invocation.Executable); err != nil {
		return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("install-source-unproven")}}
	} else if carrier := executor.ApprovedInstallSource; carrier != nil && !live.MatchesFacts(*carrier) {
		// P85 F05: the carrier facts compete item by item - path, digest, object
		// identity, platform, immutability and (on Linux) the interpreter identity -
		// so a same-digest-different-inode copy, a platform drift and an interpreter
		// swap are all refused here, with the credential already materialized only in
		// memory and no child created.
		return CopilotCLIProbeExecution{Outcome: "start-failed", Evidence: []string{copilotCLIExecutorEvidence("install-source-mismatch")}}
	}
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
	spec := guard.ProcessSpec{
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
		// P85 F05: the carrier identity travels with the spawn too, so the guard's
		// final handle must be the SAME OBJECT (same file id), not just a
		// same-digest copy. The live re-verification above already compared the
		// fact set; this pins the guard's own open.
		ExecutableFileID: approvedFileID,
	}
	// T8a (finding F8): a child that receives the BYOK credential runs in the
	// executor-owned directory, so its default write location is the directory the
	// executor destroys, and the only artifact it is handed a path to is the bounded
	// usage file. The caller's workspace stays reachable only through the explicit
	// -C argument of the screened argv.
	// P85 F14 (fix round 2): the credential is located by PLATFORM-NORMALISED key
	// comparison, never by one exact spelling. On Windows the environment block is
	// case-insensitive, so a caller that passes "copilot_provider_api_key" still
	// hands the child a credential; under exact matching the whole BYOK gate would
	// have been skipped. Two case variants of the same key (or a variant that
	// differs from a BYOK key only in case) are ambiguous and refused.
	//
	// P85 F14 (fix round 3): the classification is taken from the SERIALIZED block
	// the child will parse - after the key-name gate in copilotCLIManagedEnvironment
	// refused empty names, names carrying '=' and platform-normalised duplicates -
	// so a key name cannot change its meaning at the first '=' and slip past the
	// gate. This runs on EVERY branch below (declared secret, declared host and the
	// compatibility path), so no branch skips the end gate any more.
	credentialPresent, credentialErr := copilotCLIProbeCarriesCredential(environment)
	if credentialErr != nil {
		return CopilotCLIProbeExecution{Outcome: "start-failed", FailureCode: copilotCLIProbeBYOKFailure, Evidence: []string{copilotCLIExecutorEvidence("byok-confinement-unproven")}}
	}
	// P85 F14: the authentication mode is DECLARED on the invocation. When it is
	// declared, the observed credential presence must agree with it; it is never
	// inferred from the environment alone.
	switch invocation.AuthMode {
	case "":
		// Compatibility path for callers that predate the declaration.
	case "secret":
		if !credentialPresent {
			return CopilotCLIProbeExecution{Outcome: "start-failed", FailureCode: copilotCLIProbeBYOKFailure, Evidence: []string{copilotCLIExecutorEvidence("byok-confinement-unproven")}}
		}
	case "host":
		if credentialPresent {
			return CopilotCLIProbeExecution{Outcome: "start-failed", FailureCode: copilotCLIProbeBYOKFailure, Evidence: []string{copilotCLIExecutorEvidence("byok-confinement-unproven")}}
		}
	default:
		return CopilotCLIProbeExecution{Outcome: "start-failed", FailureCode: copilotCLIProbeBYOKFailure, Evidence: []string{copilotCLIExecutorEvidence("byok-confinement-unproven")}}
	}
	if credentialPresent {
		spec.Dir = temporaryDirectory
	}
	// T8a (finding F8): the five-element confinement proof plus the platform
	// write-restriction verdict is the last gate before a credential-bearing child
	// exists. Unknown is a refusal: S1 cannot prove the environment, so a BYOK probe
	// never starts. The evidence is a stable classification - it never echoes the
	// argv, the environment or the observed permission bits.
	if credentialPresent {
		if err := assertCopilotCLIProbeBYOKConfinementIn(runCtx, invocation, spec, temporaryDirectory, usagePath, outputLimit, executor.Sandbox); err != nil {
			return CopilotCLIProbeExecution{Outcome: "start-failed", FailureCode: copilotCLIProbeBYOKFailure, Evidence: []string{copilotCLIExecutorEvidence("byok-confinement-unproven")}}
		}
	}
	processResult, runErr := copilotCLIProbeRunManaged(runCtx, spec, 500*time.Millisecond)

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

// assertCopilotCLIProbeBYOKConfinement is the platform-default form of the
// spawn-time proof: it asks the platform check layer, exactly like an executor
// that carries no injected layer (P85 F18 keeps this shape so a caller that only
// has the free function still gets the production chain).
func assertCopilotCLIProbeBYOKConfinement(ctx context.Context, invocation CopilotCLIProbeInvocation, spec guard.ProcessSpec, temporaryDirectory, usagePath string, outputLimit int64) error {
	return assertCopilotCLIProbeBYOKConfinementIn(ctx, invocation, spec, temporaryDirectory, usagePath, outputLimit, nil)
}

// assertCopilotCLIProbeBYOKConfinementIn is the spawn-time half of the T8a (finding
// F8) proof: it checks the exact configuration the credential-bearing child would
// be created with - the five elements below - and then the platform
// write-restriction verdict. Every element is a fact about the configuration, not
// about the child's internal behaviour, which is exactly the boundary the
// contract describes: what cannot be configured away is refused, not assumed.
//
// sandbox is the check layer the caller's executor carries (P85 F18). A nil layer
// is the platform implementation, which is the production shape.
//
// The errors name the element, never a value: a refusal record must stay a closed
// classification.
func assertCopilotCLIProbeBYOKConfinementIn(ctx context.Context, invocation CopilotCLIProbeInvocation, spec guard.ProcessSpec, temporaryDirectory, usagePath string, outputLimit int64, sandbox guard.Sandbox) error {
	// P85 intentional deviation (declared to reviewer): against the architecture
	// design §5.1-3 this implementation does NOT add a
	// "proofrail:ai-byok-confinement:1\n" domain constant. Refusal evidence reuses
	// the existing T8b closed-set mechanism so no assertion-set value can leak into
	// a record; the deviation is deliberate, not an omission.
	if err := copilotCLIProbeBYOKStaticConfinement(); err != nil {
		return err
	}
	// Element 0 (P85 F13): the spawn must instantiate the FROZEN template. The
	// template was built and proven before the credential was resolved; an argv,
	// authentication mode, environment-key set, output bound, executable or working
	// directory that no longer equals it is a refusal, so "what was proven" and
	// "what is started" cannot drift apart.
	authMode := copilotCLIProbeEffectiveAuthMode(invocation, spec.Env)
	// P85 F13: the digest of the specification this spawn instantiates. The FROZEN
	// template is used when one was proven (element 0 already compared every field
	// of it, so the digest describes the spawned specification too); the compatibility
	// path derives the same fact shape from the invocation itself.
	specDigest := copilotCLIProbeAssertedSpecDigest(invocation, authMode, outputLimit)
	if frozen := invocation.FrozenSpawn; frozen != nil {
		if frozen.AuthMode != authMode || frozen.Executable != invocation.Executable || frozen.WorkingDirectory != invocation.WorkingDirectory {
			return fmt.Errorf("%w: the spawn no longer matches the frozen launch template", ErrAIBYOKConfinement)
		}
		if !slices.Equal(frozen.Args, invocation.Args) {
			return fmt.Errorf("%w: the spawn argv is not the frozen argv", ErrAIBYOKConfinement)
		}
		// T8a C2 (P85 fix round): the frozen key set is the closed OUTER envelope,
		// not an exact set every spawn must reproduce. A profile that omits the
		// optional wire-API entry emits a strict subset of it, and refusing that
		// subset refused a legitimate BYOK spawn (deepSeekProfile). What the binding
		// still proves is "nothing spills outside the proven envelope": every key the
		// child would receive must be one the template froze. The comparison stays
		// byte-exact - a platform-normalised spelling of a key is refused here and
		// again by the closed BYOK set below - and the per-key ambiguity gate keeps
		// the Windows normalisation semantics intact.
		if err := copilotCLIProbeValidateEnvironmentKeys(invocation.Environment); err != nil {
			return err
		}
		if !copilotCLIProbeEnvironmentKeysWithin(invocation.Environment, frozen.EnvironmentKeys) {
			return fmt.Errorf("%w: the spawn environment keys spill outside the frozen key set", ErrAIBYOKConfinement)
		}
		if frozen.OutputBound != outputLimit {
			return fmt.Errorf("%w: the spawn output bound is not the frozen bound", ErrAIBYOKConfinement)
		}
		// The verdict half of the binding: the template was proven under one platform
		// write verdict, and the spawn only proceeds when the same verdict is observed
		// now. An empty recorded verdict is not a verdict.
		verdict, verdictErr := copilotCLIProbePlatformWriteVerdict(ctx, sandbox, specDigest)
		if verdictErr != nil {
			return verdictErr
		}
		if frozen.PlatformWriteProof == "" || frozen.PlatformWriteProof != verdict {
			return fmt.Errorf("%w: the spawn no longer runs under the platform write verdict the template was proven with", ErrAIBYOKConfinement)
		}
	}
	// Element 1: a whitelisted argv with no logging destination. The spawn argv
	// must be exactly the screened invocation argv plus the one bounded artifact.
	if _, refused := copilotCLIProbeRefusedLoggingArg(spec.Args); refused {
		return fmt.Errorf("%w: argv names a logging destination", ErrAIBYOKConfinement)
	}
	if !slices.Equal(spec.Args, copilotCLIProbeExecArgs(invocation.Args, usagePath)) {
		return fmt.Errorf("%w: the spawn argv is not the screened probe argv", ErrAIBYOKConfinement)
	}
	if !copilotCLIProbePathInside(temporaryDirectory, usagePath) {
		return fmt.Errorf("%w: the usage artifact is outside the executor directory", ErrAIBYOKConfinement)
	}
	// Element 2: the environment whitelist. The caller's overrides stay inside the
	// closed BYOK profile key set, and the child environment the executor built
	// carries nothing beyond the inherited allowlist and those keys.
	for name := range invocation.Environment {
		if !slices.Contains(copilotCLIProbeBYOKEnvironmentKeys, name) {
			return fmt.Errorf("%w: the BYOK environment carries a key outside the profile set", ErrAIBYOKConfinement)
		}
	}
	for _, entry := range spec.Env {
		name, _, _ := strings.Cut(entry, "=")
		if !slices.Contains(copilotCLIEnvironmentAllowlist, name) && !slices.Contains(copilotCLIProbeBYOKEnvironmentKeys, name) {
			return fmt.Errorf("%w: the child environment carries a key outside the allowlist", ErrAIBYOKConfinement)
		}
	}
	// Element 3: only the standard stream handles. ProcessSpec binds no stream
	// other than stdout and stderr, and both must be the bounded in-memory capture;
	// the executable handle belongs to the guard (Windows standard handles, Linux
	// descriptor 3) and is never caller-supplied.
	//
	// P85 F21: this is the exception CONTRACTS.md:220 now spells out - "standard
	// stream handles plus the verified executable-file descriptor" on Linux, the
	// verified fd behind /proc/self/fd/3 that the guard execs. That descriptor is the
	// T3 binding mechanism itself and carries no object that was not hashed and
	// approved, so naming it strengthens no claim; the standard-stream rule below is
	// what closes every OTHER handle.
	if spec.Stdout == nil || spec.Stderr == nil {
		return fmt.Errorf("%w: the child streams are not bound", ErrAIBYOKConfinement)
	}
	if _, isFile := spec.Stdout.(*os.File); isFile {
		return fmt.Errorf("%w: stdout is bound to a file handle", ErrAIBYOKConfinement)
	}
	if _, isFile := spec.Stderr.(*os.File); isFile {
		return fmt.Errorf("%w: stderr is bound to a file handle", ErrAIBYOKConfinement)
	}
	// Element 4: the working directory is the executor-owned temporary directory,
	// the only directory the child may reach without an explicit argument.
	if spec.Dir == "" || filepath.Clean(spec.Dir) != filepath.Clean(temporaryDirectory) {
		return fmt.Errorf("%w: the child working directory is not the executor directory", ErrAIBYOKConfinement)
	}
	// Element 5: the in-memory capture is bounded; a zero or negative bound is an
	// unbounded buffer and therefore an unprovable one.
	if outputLimit <= 0 {
		return fmt.Errorf("%w: the in-memory output capture is unbounded", ErrAIBYOKConfinement)
	}
	// Element 6 (P85 fix round 1, F02): the working directory of the child is not
	// only ProcessSpec.Dir - the screened argv also carries the CLI's own
	// working-directory argument, which the child resolves itself. A `-C` /
	// `--working-directory` naming any other directory would let the
	// credential-bearing child reach a tree the executor does not own (or destroy) and
	// would contradict the element-4 claim. Every occurrence must therefore name the
	// executor's temporary directory; an argv that carries no such argument keeps the
	// ProcessSpec.Dir binding alone.
	if err := assertCopilotCLIProbeArgvWorkingDirectory(spec.Args, temporaryDirectory); err != nil {
		return err
	}
	// The argv allowlist applies to the argv that will actually be spawned, on every
	// path - including the compatibility path that carries no frozen template, and
	// including an executor called directly with caller-supplied args. It runs after
	// the F02 working-directory assertion so that assertion keeps its own, more
	// precise diagnosis for the shape it owns: the F13 ruling narrows F02 (a
	// credential-bearing argv may not name a working directory at all) without
	// removing it.
	if err := copilotCLIProbeValidateWhitelistedArgs(invocation.Args, authMode); err != nil {
		return err
	}
	// The platform write-restriction verdict: on S1 every check is Unknown, and
	// Unknown is a refusal exactly like Denied. P85 F13: the check is asked about the
	// specification this spawn instantiates, so the verdict is an answer about the
	// frozen template rather than about "a probe".
	if err := guard.CheckSandbox(ctx, sandbox, guard.CheckSpec{Reason: guard.SandboxReasonBYOKProbe, SpecDigest: specDigest}); err != nil {
		return fmt.Errorf("%w: %v", ErrAIBYOKConfinement, err)
	}
	return nil
}

// copilotCLIProbeAssertedSpecDigest is the digest of the specification a spawn
// instantiates: the frozen template when the pre-Resolve proof produced one, the
// invocation-derived template otherwise (P85 F13). The two agree on every path that
// carries a frozen template, because element 0 of the proof compared each field
// before this digest is used.
func copilotCLIProbeAssertedSpecDigest(invocation CopilotCLIProbeInvocation, authMode string, outputLimit int64) string {
	if frozen := invocation.FrozenSpawn; frozen != nil {
		return copilotCLIProbeSpawnSpecDigest(*frozen)
	}
	return copilotCLIProbeSpawnSpecDigest(copilotCLIProbeSpawnSpec{
		Executable:       invocation.Executable,
		WorkingDirectory: invocation.WorkingDirectory,
		AuthMode:         authMode,
		Args:             invocation.Args,
		EnvironmentKeys:  copilotCLIProbeSortedEnvironmentKeys(invocation.Environment),
		OutputBound:      outputLimit,
	})
}

// assertCopilotCLIProbeArgvWorkingDirectory is the F02 half of the BYOK spawn
// proof: the child's implicit working directory (ProcessSpec.Dir) is only one of
// the two ways a copilot CLI child learns where to work, because the same argv
// can carry an explicit `-C <dir>` / `--working-directory <dir>` (or the
// `--flag=<dir>` form). This assertion requires every such occurrence to name the
// executor-owned temporary directory after cleaning; a dangling flag or a value
// that resolves anywhere else is refused. An argv with no such argument passes:
// the ProcessSpec.Dir element already covers that case.
func assertCopilotCLIProbeArgvWorkingDirectory(args []string, temporaryDirectory string) error {
	for index := 0; index < len(args); index++ {
		arg := strings.TrimSpace(args[index])
		name, value, inlineValue := copilotCLIProbeWorkingDirectoryArg(arg)
		if name == "" {
			continue
		}
		if !inlineValue {
			if index+1 >= len(args) {
				return fmt.Errorf("%w: the argv names a working directory without a value", ErrAIBYOKConfinement)
			}
			value = args[index+1]
			index++
		}
		if strings.TrimSpace(value) == "" || filepath.Clean(value) != filepath.Clean(temporaryDirectory) {
			return fmt.Errorf("%w: the argv names a working directory outside the executor directory", ErrAIBYOKConfinement)
		}
	}
	return nil
}

// copilotCLIProbeWorkingDirectoryArg classifies one argv token as a working
// directory flag. inlineValue reports whether the value travels inside the same
// token (`--working-directory=<dir>`); otherwise the value is the next token.
func copilotCLIProbeWorkingDirectoryArg(arg string) (name, value string, inlineValue bool) {
	for _, flag := range []string{"-C", "--working-directory"} {
		switch {
		case arg == flag:
			return flag, "", false
		case strings.HasPrefix(arg, flag+"="):
			return flag, arg[len(flag)+1:], true
		}
	}
	return "", "", false
}

// copilotCLIProbePathInside reports whether path stays inside directory. Both are
// cleaned first, so ".." cannot walk out of the executor-owned directory.
func copilotCLIProbePathInside(directory, path string) bool {
	relative, err := filepath.Rel(filepath.Clean(directory), filepath.Clean(path))
	if err != nil || filepath.IsAbs(relative) {
		return false
	}
	return relative != ".." && !strings.HasPrefix(relative, ".."+string(os.PathSeparator))
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

// copilotCLIEnvironmentKeyEqual compares two environment key names the way the
// platform compares them (P85 F14): case-insensitively on Windows, byte-exactly
// elsewhere. Comparing exact bytes on Windows is what let a lowercase credential key
// reach the child while the BYOK gate concluded "no credential".
func copilotCLIEnvironmentKeyEqual(left, right string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(left, right)
	}
	return left == right
}

// copilotCLIProbeValidateEnvironmentKeys refuses an override map whose key names
// cannot survive the environment block the child is handed (P85 F14 fix round 3).
//
// The child parses `NAME=VALUE` and takes everything before the FIRST `=` as the
// name, so a key that itself contains `=` is not the key the child sees: the
// caller could hand `COPILOT_PROVIDER_API_KEY=x` as a name and the child would
// receive `COPILOT_PROVIDER_API_KEY` with the credential while the raw map still
// reports "no credential". An empty name is not addressable either. On Windows
// the block is case-insensitive, so two names that differ only in case are two
// spellings of one variable and the child's value would depend on ordering.
func copilotCLIProbeValidateEnvironmentKeys(environment map[string][]byte) error {
	if len(environment) == 0 {
		return nil
	}
	seen := make(map[string]string, len(environment))
	for name := range environment {
		if name == "" {
			return fmt.Errorf("%w: the environment carries an empty key name", ErrAIBYOKConfinement)
		}
		if strings.ContainsRune(name, '=') {
			return fmt.Errorf("%w: an environment key name contains '=', so the key the child parses is not the key the caller named", ErrAIBYOKConfinement)
		}
		if strings.ContainsRune(name, 0) {
			return fmt.Errorf("%w: an environment key name contains a NUL byte", ErrAIBYOKConfinement)
		}
		normalized := name
		if runtime.GOOS == "windows" {
			normalized = strings.ToUpper(name)
		}
		if other, duplicate := seen[normalized]; duplicate {
			return fmt.Errorf("%w: the environment names %q and %q, which the platform cannot tell apart", ErrAIBYOKConfinement, other, name)
		}
		seen[normalized] = name
	}
	return nil
}

// copilotCLIProbeCarriesCredential reports whether the SERIALIZED environment block
// - the exact K=V entries the child receives - hands the child the BYOK credential
// (P85 F14 fix round 3).
//
// The classification is taken from the block, never from the caller's raw map: the
// child parses names at the first `=`, so the block is the only set of names that
// describes what the child really sees. A name that matches the credential key only
// under platform normalisation is refused: on a case-insensitive platform the child
// would receive a credential the gate cannot name, and on a case-sensitive one it is
// a different key that the closed environment set rejects.
//
// This runs on EVERY authentication branch - declared secret, declared host and the
// compatibility path - so no branch can skip the end gate.
func copilotCLIProbeCarriesCredential(environment []string) (bool, error) {
	present := false
	for _, entry := range environment {
		name, value, found := strings.Cut(entry, "=")
		if !found {
			return false, fmt.Errorf("%w: an environment entry carries no name/value separator", ErrAIBYOKConfinement)
		}
		if !copilotCLIEnvironmentKeyEqual(name, copilotCLIProbeCredentialEnvironmentKey) {
			continue
		}
		if name != copilotCLIProbeCredentialEnvironmentKey {
			return false, fmt.Errorf("%w: the environment names a case variant of the credential key", ErrAIBYOKConfinement)
		}
		if value != "" {
			present = true
		}
	}
	return present, nil
}

// copilotCLIProbeEffectiveAuthMode is the authentication mode the spawn actually has:
// the DECLARED mode when the invocation carries one, otherwise the mode implied by
// the credential key of the SERIALIZED environment block (P85 F14 fix round 3). The
// block, not the raw map, is what the child parses, so it is what decides.
func copilotCLIProbeEffectiveAuthMode(invocation CopilotCLIProbeInvocation, environment []string) string {
	if invocation.AuthMode != "" {
		return invocation.AuthMode
	}
	if present, err := copilotCLIProbeCarriesCredential(environment); err == nil && present {
		return "secret"
	}
	return "host"
}

// copilotCLIProbeEnvironmentKeysWithin reports whether every key of one override map
// is named by the frozen key set. P85 F13 follow-up (T8a C2): the frozen template
// carries the closed OUTER envelope and a profile may declare a strict subset of it
// (the wire-API entry is optional), so the binding is containment, not equality.
//
// The membership test is byte-exact on purpose: a platform-normalised spelling the
// frozen set does not name is refused here and again by the closed BYOK key set in
// assertCopilotCLIProbeBYOKConfinement, so loosening the equality does not loosen the
// Windows normalization semantics of copilotCLIProbeValidateEnvironmentKeys.
func copilotCLIProbeEnvironmentKeysWithin(environment map[string][]byte, frozen []string) bool {
	for name := range environment {
		if !slices.Contains(frozen, name) {
			return false
		}
	}
	return true
}

// copilotCLIProbeSortedEnvironmentKeys returns the sorted key set of one override map,
// so the frozen template's key set and the spawn's key set are compared as data.
func copilotCLIProbeSortedEnvironmentKeys(environment map[string][]byte) []string {
	if len(environment) == 0 {
		return nil
	}
	names := make([]string, 0, len(environment))
	for name := range environment {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// copilotCLIManagedEnvironment builds the child environment block.
//
// P85 F14 (fix round 3): the key names are validated FIRST, so no name that would
// change its meaning in the K=V serialization reaches the block. Then an override
// whose name collides with an inherited allowlist key under the PLATFORM comparison
// (case only, on Windows) is a conflict - the child would receive one variable under
// two spellings - and it refuses instead of letting one silently win. The output
// carries a single canonical spelling of every key.
func copilotCLIManagedEnvironment(overrides map[string][]byte) ([]string, [][]byte, error) {
	if err := copilotCLIProbeValidateEnvironmentKeys(overrides); err != nil {
		return nil, nil, err
	}
	for name := range overrides {
		for _, allowed := range copilotCLIEnvironmentAllowlist {
			if name != allowed && copilotCLIEnvironmentKeyEqual(name, allowed) {
				return nil, nil, fmt.Errorf("%w: the environment override collides with the inherited key %q", ErrAIBYOKConfinement, allowed)
			}
		}
	}
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
	return environment, buffers, nil
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
