package console

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/larsonzh/prfrail/internal/adapters"
)

var ErrAIChannelNotConfigured = errors.New("AI channel not configured")

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
	const usage = "usage: prfrail ai check --channel agent-runner-cli --copilot <path> --max-requests 1 --out <path> [--chain <path>] [--workspace <path>] [--json]"
	set, parseErr := newFlagSet("ai check")
	chainPath := set.String("chain", "", "path to chain config file")
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
	record, err := cli.checkAIAvailability(ctx, profile, *channel, *maximumRequests, resolvedExecutable, workingDirectory)
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

func (cli CLI) defaultAICheck(ctx context.Context, profile adapters.AIProviderProfile, channel string, maximumRequests int, executable, workingDirectory string) (adapters.AIAvailabilityRecord, error) {
	checker := adapters.AIAvailabilityChecker{
		Clock:   cli.now,
		IDs:     newAIProbeID,
		Secrets: cli.secretManager,
		Probe: adapters.CopilotCLIProbe{
			Executable:       executable,
			WorkingDirectory: workingDirectory,
			Executor:         adapters.ManagedCopilotCLIProbeExecutor{},
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
