package console

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

var ErrAIChannelNotConfigured = errors.New("AI channel not configured")

// ErrAIReceiverBinding reports a probe whose receiver (approved profile, channel
// and executable) cannot be bound. F6: no credential is resolved and no process
// is started when this error is returned.
var ErrAIReceiverBinding = errors.New("AI probe receiver binding invalid")

// aiProbeReceiverBindingDomain separates the receiver binding digest from every
// other probe digest.
const aiProbeReceiverBindingDomain = "proofrail:ai-probe-receiver-binding:1\n"

// aiProbeToolID is the tool identity the operator grant must name for the
// Copilot CLI availability probe to be allowed to reach its receiver (F6).
const aiProbeToolID = "tool-ai-probe"

// aiProbeDefaultPort is the port an HTTPS base URL without an explicit port
// resolves to, so the receiver a grant has to name is unambiguous.
const aiProbeDefaultPort = 443

// validateAIProbeBinding proves the probe receiver is well formed before anything
// is read from the secret store or spawned.
//
// F6 (SECURITY.md:60 "network is denied by default; an exception must name its
// receiver"): the allowed exception is the operator-configured profile plus the
// operator-supplied executable, so the binding digest covers exactly the
// profile, the channel and the digest of the executable that will run. Any
// failure - an invalid profile, an unreadable executable, an oversized
// executable - refuses the probe instead of degrading to "assume approved".
//
// P83 FIX ROUND 4 (F6 not closed): this is only the SHAPE check. A self-computed
// digest proves nothing about approval, so it is not the gate; the caller hands
// both returned digests to authorizeAIProbeReceiver, which compares them against
// the operator's trusted authorization ledger before the secret store is read.
func validateAIProbeBinding(profile adapters.AIProviderProfile, channel, executablePath string) (binding, executableHash string, err error) {
	if err := adapters.ValidateAIProviderProfile(profile); err != nil {
		return "", "", fmt.Errorf("%w: %w", ErrAIReceiverBinding, err)
	}
	executableHash, err = adapters.CopilotCLIExecutableHash(executablePath)
	if err != nil {
		return "", "", fmt.Errorf("%w: hash the probe executable: %v", ErrAIReceiverBinding, err)
	}
	profileHash, err := adapters.AIProviderProfileHash(profile)
	if err != nil {
		return "", "", fmt.Errorf("%w: %w", ErrAIReceiverBinding, err)
	}
	binding = evidence.Digest(aiProbeReceiverBindingDomain, []byte(profileHash+"\n"+channel+"\n"+executableHash))
	return binding, executableHash, nil
}

// aiProbeReceiverAddress reduces the profile base URL to the network receiver a
// probe would contact: the host and port the operator grant has to name.
// ValidateAIProviderProfile already forced the HTTPS scheme, so only an empty
// base URL is refused here - a profile that names no receiver cannot be approved.
func aiProbeReceiverAddress(profile adapters.AIProviderProfile) (string, int, error) {
	baseURL := profile.BaseURL
	if baseURL == "" {
		return "", 0, fmt.Errorf("%w: the profile declares no base URL, so the probe receiver cannot be named", ErrAIReceiverBinding)
	}
	remainder := strings.TrimPrefix(baseURL, "https://")
	authority, _, _ := strings.Cut(remainder, "/")
	host, portText, hasPort := strings.Cut(authority, ":")
	if host == "" {
		return "", 0, fmt.Errorf("%w: the profile base URL %q names no host", ErrAIReceiverBinding, baseURL)
	}
	if !hasPort {
		return host, aiProbeDefaultPort, nil
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return "", 0, fmt.Errorf("%w: the profile base URL %q names an invalid port", ErrAIReceiverBinding, baseURL)
	}
	return host, port, nil
}

// loadProbeAuthorizationLedger loads the ledger the probe must be authorised by.
//
// F6: a ledger that is absent, unreadable or malformed refuses the probe; it must
// never degrade to "no grant needed" (which is what an empty ledger would look
// like if the load error were swallowed).
func loadProbeAuthorizationLedger(path string) (chain.AuthorizationLedger, error) {
	if _, err := os.Stat(path); err != nil {
		return chain.AuthorizationLedger{}, fmt.Errorf("%w: read the authorization ledger %s: %v", ErrAIReceiverBinding, path, err)
	}
	ledger, err := loadAuthorizationLedger(path)
	if err != nil {
		return chain.AuthorizationLedger{}, fmt.Errorf("%w: read the authorization ledger %s: %v", ErrAIReceiverBinding, path, err)
	}
	return ledger, nil
}

// aiProbeReceiver is the receiver the CLI admitted for one probe: the path the
// approved object was opened at, and the digest the approval actually carried.
//
// T3 (P83 fix round 6): the digest is a value that travels FROM the approval to
// the executor. It is never recomputed from Executable further down the call
// chain, because a digest derived from a caller-chosen path at spawn time proves
// nothing about what the operator approved.
type aiProbeReceiver struct {
	Executable     string
	ExecutableHash string
}

// authorizeAIProbeReceiver is the F6 gate: the probe receiver is accepted only
// when the operator's trusted authorization ledger carries an ACTIVE grant whose
// scope names the receiver, whose toolIds name the probe, and whose evidence
// carries the FULL receiver identity.
//
// T3 (finding F6 "complete receiver binding"): matching the host/port plus the
// executable digest is not enough - two profiles can share an address and differ
// in provider, model or secretRef. The grant must therefore also carry the
// profile hash (which covers provider/baseUrl/model/secretRef) and the receiver
// binding digest (which additionally covers the channel and the executable). A
// same-address profile with a different secretRef therefore no longer passes.
//
// T9 (finding: over-broad wiring): the gate applies only when a network receiver
// is actually used. A profile that names no base URL authenticates by host
// identity, reads no credential and contacts no receiver, so it is not refused
// here.
//
// The grant shape is the existing schemas/authorization-record.schema.json grant,
// whose evidence is an open hashSet (minItems 1), so no new schema and no new
// approval document is introduced. Any failure - no grant, no matching receiver,
// no matching tool, a missing executable/profile/binding digest, expired,
// revoked, an unevaluable ledger - refuses BEFORE the secret store is read and
// before any process is spawned.
func authorizeAIProbeReceiver(ledger chain.AuthorizationLedger, profile adapters.AIProviderProfile, binding, executableHash string, now time.Time) error {
	_, err := authorizeAIProbeReceiverApproval(ledger, profile, binding, executableHash, now)
	return err
}

// authorizeAIProbeReceiverApproval is the approval gate with the approved
// identity returned instead of discarded.
//
// T3 (P83 fix round 6): the digest that admitted the receiver is the grant's own
// evidence entry, and that exact value - not a digest re-derived from the path -
// is what the executor pins. The two are equal whenever the gate admits, but the
// provenance is what matters: the value the process is started under is the one
// the operator wrote into the ledger.
func authorizeAIProbeReceiverApproval(ledger chain.AuthorizationLedger, profile adapters.AIProviderProfile, binding, executableHash string, now time.Time) (aiProbeReceiver, error) {
	if profile.BaseURL == "" {
		return aiProbeReceiver{ExecutableHash: executableHash}, nil
	}
	host, port, err := aiProbeReceiverAddress(profile)
	if err != nil {
		return aiProbeReceiver{}, err
	}
	profileHash, err := adapters.AIProviderProfileHash(profile)
	if err != nil {
		return aiProbeReceiver{}, fmt.Errorf("%w: %w", ErrAIReceiverBinding, err)
	}
	for _, record := range ledger.Grants {
		grant := record.Grant
		if grant == nil {
			continue
		}
		// A superseded grant record is history, not an approval: a re-issued
		// authorization may narrow its scope, so only the latest record for an
		// authorizationId may be read as that operator's current intent.
		latest, _, latestExists := ledger.LatestGrant(grant.AuthorizationID)
		if !latestExists || latest.RecordHash != record.RecordHash {
			continue
		}
		// A grant is only an approval while it is active: an expired or revoked
		// authorization never authorises a receiver, and an unevaluable one is
		// refused rather than skipped.
		status, err := chain.EvaluateGrant(grant.AuthorizationID, ledger, now)
		if err != nil {
			return aiProbeReceiver{}, fmt.Errorf("%w: evaluate grant %q: %v", ErrAIReceiverBinding, grant.AuthorizationID, err)
		}
		if status != chain.GrantActive {
			continue
		}
		if !authorizationCoversReceiver(grant.Scope.Network, host, port) {
			continue
		}
		if !authorizationNamesTool(grant.Scope.ToolIDs, aiProbeToolID) {
			continue
		}
		if !authorizationCarriesHash(grant.Evidence, executableHash) {
			continue
		}
		// T3: the grant must approve this exact profile (provider, baseUrl, model
		// and secretRef all live in its hash) and this exact receiver binding
		// (profile + channel + executable). Either one missing means the grant was
		// written for a different receiver identity.
		if !authorizationCarriesHash(grant.Evidence, profileHash) {
			continue
		}
		if !authorizationCarriesHash(grant.Evidence, binding) {
			continue
		}
		// T3: hand back the grant's own entry, so the executor pins the approved
		// digest rather than a value this process happened to compute.
		if approved, found := authorizationApprovedHash(grant.Evidence, executableHash); found {
			return aiProbeReceiver{ExecutableHash: approved}, nil
		}
		return aiProbeReceiver{}, fmt.Errorf("%w: the grant records no executable digest", ErrAIReceiverBinding)
	}
	return aiProbeReceiver{}, fmt.Errorf("%w: no active grant authorises receiver %s:%d for tool %s with profile hash %s, executable digest %s and receiver binding %s (%d grant record(s))",
		ErrAIReceiverBinding, host, port, aiProbeToolID, profileHash, executableHash, binding, len(ledger.Grants))
}

// authorizationApprovedHash returns the grant's own evidence entry that matches
// wanted. T3: the entry - not the caller's argument - is the approved digest.
func authorizationApprovedHash(hashes []string, wanted string) (string, bool) {
	for _, candidate := range hashes {
		if candidate == wanted {
			return candidate, true
		}
	}
	return "", false
}

// authorizationCoversReceiver reports whether one network scope names exactly the
// receiver a probe contacts: HTTPS to the same host and port.
func authorizationCoversReceiver(scope []chain.AuthorizationNetworkScope, host string, port int) bool {
	for _, entry := range scope {
		if entry.Protocol == "https" && entry.Host == host && entry.Port == port {
			return true
		}
	}
	return false
}

// authorizationNamesTool reports whether the scope names toolID.
func authorizationNamesTool(toolIDs []string, toolID string) bool {
	for _, candidate := range toolIDs {
		if candidate == toolID {
			return true
		}
	}
	return false
}

// authorizationCarriesHash reports whether the grant records the exact digest.
// A grant that carries no matching digest names no executable and therefore
// approves none: the probe refuses instead of assuming the grant is about the
// file it was handed.
func authorizationCarriesHash(hashes []string, wanted string) bool {
	for _, candidate := range hashes {
		if candidate == wanted {
			return true
		}
	}
	return false
}

func (cli CLI) executeAI(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		writeAIUsage(stderr)
		return exitUsage
	}
	switch args[0] {
	case "check":
		return cli.executeAICheck(ctx, args[1:], stdout, stderr)
	case "verify":
		return cli.executeAIVerify(args[1:], stdout, stderr)
	default:
		writeAIUsage(stderr)
		return exitUsage
	}
}

func (cli CLI) executeAICheck(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	const usage = "usage: prfrail ai check --channel agent-runner-cli --copilot <path> --max-requests 1 --out <path> [--chain <path>] [--authorization-ledger <path>] [--workspace <path>] [--json]"
	set, parseErr := newFlagSet("ai check")
	chainPath := set.String("chain", "", "path to chain config file")
	// F6: the receiver approval lives in the operator's authorization ledger. The
	// path is optional and defaults to the ledger next to the chain configuration,
	// so no new REQUIRED flag is introduced; an absent ledger refuses the probe.
	authorizationLedger := set.String("authorization-ledger", "", "path to the authorization ledger that must carry the probe receiver grant")
	channel := set.String("channel", "", "configured AI channel")
	executable := set.String("copilot", "", "path to the pinned Copilot CLI executable")
	workspace := set.String("workspace", "", "probe working directory")
	outputPath := set.String("out", "", "write the immutable availability record")
	maximumRequests := set.Int("max-requests", 0, "maximum provider requests")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "ai check", *jsonOutput, parseErr.String(), usage)
	}
	if len(set.Args()) != 0 || *channel == "" || *executable == "" || *maximumRequests != 1 || *outputPath == "" {
		return writeUsageError(stdout, stderr, "ai check", *jsonOutput, "--channel, --copilot, --max-requests 1 and --out are required; positional values are forbidden", usage)
	}
	if *channel != "agent-runner-cli" {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, fmt.Errorf("%w: %s", ErrAIChannelNotConfigured, *channel))
	}
	if cli.checkAIAvailability == nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, errors.New("AI availability checker unavailable"))
	}

	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
	}
	resolvedConfig, err := ResolveChainConfigPath(*chainPath, cwd)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
	}
	config, _, err := LoadChainConfig(resolvedConfig)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
	}
	// F6/T3: the probe is authorised by the ledger of the chain configuration and
	// by nothing else. The path is DERIVED from the chain-config directory (the
	// durable root the run and its approvals belong to) and is re-read at the
	// moment of use; --authorization-ledger may only name that same file. A
	// disagreement is a refusal - honouring the caller's path is exactly the
	// self-approval hole finding F6 names.
	selectedLedger := ""
	if *authorizationLedger != "" {
		selectedLedger, err = resolvePathFromCWD(cwd, *authorizationLedger)
		if err != nil {
			return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
		}
	}
	ledgerPath, err := trustedLedgerPath(filepath.Dir(resolvedConfig), defaultApprovalsLedgerName, selectedLedger)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, fmt.Errorf("%w: %v", ErrAIReceiverBinding, err))
	}
	profile, found := configuredAIProfile(config, *channel)
	if !found {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, fmt.Errorf("%w: %s", ErrAIChannelNotConfigured, *channel))
	}
	resolvedExecutable, err := resolvePathFromCWD(cwd, *executable)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
	}
	workingDirectory := *workspace
	if workingDirectory == "" {
		workingDirectory = cwd
	} else {
		workingDirectory, err = resolvePathFromCWD(cwd, workingDirectory)
		if err != nil {
			return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
		}
	}
	resolvedOutput, err := resolvePathFromCWD(cwd, *outputPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
	}
	if _, err := os.Lstat(resolvedOutput); err == nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, fmt.Errorf("availability record output already exists: %s", resolvedOutput))
	} else if !errors.Is(err, os.ErrNotExist) {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
	}
	// F6: the receiver binding is verified before the availability check runs, so a
	// mismatched profile/executable pair never reaches the secret store and never
	// starts a child process. The binding digest alone is only a self-description,
	// so the ledger comparison below is what actually admits the receiver.
	//
	// T3 (P83 fix round 6): the digest computed here is the SHAPE check. The value
	// the probe is started under is the one the approval returns, and it is carried
	// in aiProbeReceiver all the way to the executor - never recomputed there.
	binding, executableHash, err := validateAIProbeBinding(profile, *channel, resolvedExecutable)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
	}
	receiver := aiProbeReceiver{Executable: resolvedExecutable, ExecutableHash: executableHash}
	// F6 (P83 fix round 4): the probe receiver is accepted only by an ACTIVE
	// operator grant in the trusted authorization ledger. The refusal happens here,
	// before the check runs, so no credential is resolved and no process is spawned.
	//
	// T9 (P83 fix round 6): the ledger is only consulted when a network receiver is
	// actually used. A host-auth profile names no base URL, contacts no receiver and
	// reads no grant, so requiring an operable ledger for it refused a legitimate
	// profile before the exemption could ever be reached.
	if profile.BaseURL != "" {
		ledger, err := loadProbeAuthorizationLedger(ledgerPath)
		if err != nil {
			return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
		}
		approved, err := authorizeAIProbeReceiverApproval(ledger, profile, binding, executableHash, cli.now())
		if err != nil {
			return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
		}
		receiver.ExecutableHash = approved.ExecutableHash
	}
	record, err := cli.checkAIAvailability(ctx, profile, *channel, *maximumRequests, receiver, workingDirectory)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, err)
	}
	if writeErr := adapters.WriteAIAvailabilityRecord(resolvedOutput, record); writeErr != nil {
		return writeCommandError(stdout, stderr, "ai check", *jsonOutput, writeErr)
	}
	available := record.Availability.Status == "available"
	exitCode := exitFailure
	if available {
		exitCode = exitSuccess
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "ai check", OK: available, ExitCode: exitCode, Data: record})
	}
	fmt.Fprintf(stdout, "ai check: channel=%s profile=%s status=%s requests=%d record=%s\n", record.Availability.Channel, record.Availability.ProfileID, record.Availability.Status, record.Availability.RequestsUsed, record.RecordHash)
	if record.Availability.Reason != nil {
		fmt.Fprintf(stdout, "reason: %s\n", *record.Availability.Reason)
	}
	return exitCode
}

func (cli CLI) executeAIVerify(args []string, stdout, stderr io.Writer) int {
	const usage = "usage: prfrail ai verify --record <path> --channel <channel> --max-age <duration> --max-requests 1 [--prior-record <path>]... [--chain <path>] [--json]"
	set, parseErr := newFlagSet("ai verify")
	chainPath := set.String("chain", "", "path to chain config file")
	recordPath := set.String("record", "", "path to an AI availability record")
	priorRecordPaths := make([]string, 0)
	set.Func("prior-record", "path to a prior AI availability record; repeatable", func(value string) error {
		if value == "" {
			return errors.New("prior record path is empty")
		}
		priorRecordPaths = append(priorRecordPaths, value)
		return nil
	})
	channel := set.String("channel", "", "configured AI channel")
	maximumAge := set.Duration("max-age", 0, "maximum record age")
	maximumRequests := set.Int("max-requests", 0, "maximum provider requests")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "ai verify", *jsonOutput, parseErr.String(), usage)
	}
	if len(set.Args()) != 0 || *recordPath == "" || *channel == "" || *maximumAge <= 0 || *maximumRequests != 1 {
		return writeUsageError(stdout, stderr, "ai verify", *jsonOutput, "--record, --channel, a positive --max-age and --max-requests 1 are required; positional values are forbidden", usage)
	}
	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	resolvedConfig, err := ResolveChainConfigPath(*chainPath, cwd)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	config, _, err := LoadChainConfig(resolvedConfig)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	profile, found := configuredAIProfile(config, *channel)
	if !found {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, fmt.Errorf("%w: %s", ErrAIChannelNotConfigured, *channel))
	}
	profileHash, err := adapters.AIProviderProfileHash(profile)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	resolvedRecord, err := resolvePathFromCWD(cwd, *recordPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	wire, err := os.ReadFile(resolvedRecord)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	record, err := adapters.DecodeAIAvailabilityRecord(wire)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	replayIndex, err := adapters.NewAIAvailabilityIndex(nil)
	if err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	for _, priorPath := range priorRecordPaths {
		resolvedPrior, resolveErr := resolvePathFromCWD(cwd, priorPath)
		if resolveErr != nil {
			return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, resolveErr)
		}
		priorWire, readErr := os.ReadFile(resolvedPrior)
		if readErr != nil {
			return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, readErr)
		}
		priorRecord, decodeErr := adapters.DecodeAIAvailabilityRecord(priorWire)
		if decodeErr != nil {
			return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, decodeErr)
		}
		if _, replayErr := replayIndex.Record(priorRecord); replayErr != nil {
			return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, replayErr)
		}
	}
	if _, err := replayIndex.Record(record); err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	if err := adapters.ValidateAIAvailabilityAdmission(record, adapters.AIAvailabilityPolicy{
		ProfileID:         profile.ProfileID,
		ProfileConfigHash: profileHash,
		Channel:           *channel,
		EvaluatedAt:       cli.now().UTC(),
		MaximumAge:        *maximumAge,
		MaximumRequests:   *maximumRequests,
	}); err != nil {
		return writeCommandError(stdout, stderr, "ai verify", *jsonOutput, err)
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "ai verify", OK: true, ExitCode: exitSuccess, Data: record})
	}
	fmt.Fprintf(stdout, "ai verify: channel=%s profile=%s status=available record=%s\n", record.Availability.Channel, record.Availability.ProfileID, record.RecordHash)
	return exitSuccess
}

func writeAIUsage(output io.Writer) {
	fmt.Fprintln(output, "usage: prfrail ai <check|verify> [options]")
}

func (cli CLI) defaultAICheck(ctx context.Context, profile adapters.AIProviderProfile, channel string, maximumRequests int, receiver aiProbeReceiver, workingDirectory string) (adapters.AIAvailabilityRecord, error) {
	// F6: the executor pins the exact executable the console bound, so a file that
	// changes between the binding check and the spawn is refused rather than
	// started.
	//
	// T3 (P83 fix round 6): the pinned digest is the one the approval carried. It is
	// deliberately NOT recomputed here: hashing receiver.Executable at this point
	// would re-open a caller-chosen path and silently replace the approved identity
	// with "whatever is at the path now". A receiver without a digest would fall back
	// to the executor's unpinned compatibility path, which is exactly the silent
	// downgrade this item removes, so it is refused instead.
	if receiver.ExecutableHash == "" {
		return adapters.AIAvailabilityRecord{}, fmt.Errorf("%w: the approved executable digest is missing", ErrAIReceiverBinding)
	}
	checker := adapters.AIAvailabilityChecker{
		Clock:   cli.now,
		IDs:     newAIProbeID,
		Secrets: cli.secretManager,
		Probe: adapters.CopilotCLIProbe{
			Executable:       receiver.Executable,
			WorkingDirectory: workingDirectory,
			Executor:         adapters.ManagedCopilotCLIProbeExecutor{ExpectedExecutableHash: receiver.ExecutableHash},
		},
	}
	return checker.Check(ctx, profile, channel, maximumRequests)
}

func configuredAIProfile(config ChainConfig, channel string) (adapters.AIProviderProfile, bool) {
	if config.AI == nil {
		return adapters.AIProviderProfile{}, false
	}
	profileID := ""
	for _, binding := range config.AI.Channels {
		if binding.Channel == channel {
			profileID = binding.ProfileID
			break
		}
	}
	if profileID == "" {
		return adapters.AIProviderProfile{}, false
	}
	for _, profile := range config.AI.Profiles {
		if profile.ProfileID == profileID {
			return profile, true
		}
	}
	return adapters.AIProviderProfile{}, false
}

func newAIProbeID() string {
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		return "probe-unavailable"
	}
	return "probe-" + hex.EncodeToString(random)
}
