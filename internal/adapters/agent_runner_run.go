package adapters

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/tickets"
)

var (
	// ErrInvalidAgentRunnerRun reports an unusable offline run request.
	ErrInvalidAgentRunnerRun = errors.New("AgentRunner run invalid")
)

// agentRunnerOutcomeEvidenceDomain prefixes the digest that names why an outcome
// was degraded, so a non-completed receipt always carries provable error evidence.
const agentRunnerOutcomeEvidenceDomain = "proofrail:agent-runner-outcome:1\n"

// AgentRunnerRunRequest pins one offline run. TaskID, StepID and Attempt identify
// the step the completion receipt must bind to.
type AgentRunnerRunRequest struct {
	Launch      chain.AgentRunnerLaunchRequest
	RequestHash string
	TaskID      string
	StepID      string
	Attempt     int
	Timeout     time.Duration
	Capturer    AgentRunnerManifestCapturer
	MaxLogBytes int64
	// VerifyWait bounds the process-tree verification after a natural exit.
	VerifyWait time.Duration
}

// AgentRunnerPinnedCLIRun orchestrates one offline run: start, bounded wait,
// evidence collection and the outcome mapping. It never publishes records and
// never writes chain state: it returns the outcome the caller terminalizes.
type AgentRunnerPinnedCLIRun struct {
	Launcher *AgentRunnerPinnedCLILauncher
	Store    *AgentRunnerReplayStore
	Clock    func() time.Time
}

// Run executes one offline run. Every error it returns means the caller must not
// terminalize a success: the outcome it does return has already been mapped to
// whatever the evidence proves.
func (run *AgentRunnerPinnedCLIRun) Run(ctx context.Context, request AgentRunnerRunRequest) (AgentRunnerTerminalOutcome, error) {
	if run == nil || run.Launcher == nil || run.Store == nil {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: launcher and store are required", ErrInvalidAgentRunnerRun)
	}
	if err := ctx.Err(); err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	if !evidence.ValidHash(request.RequestHash) {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: the request record hash is required", ErrInvalidAgentRunnerRun)
	}
	if !evidence.ValidID(request.TaskID) || !evidence.ValidID(request.StepID) || request.Attempt < 1 {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: task, step and attempt are required", ErrInvalidAgentRunnerRun)
	}
	if request.Timeout <= 0 {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: a positive timeout is required", ErrInvalidAgentRunnerRun)
	}
	if request.MaxLogBytes <= 0 {
		return AgentRunnerTerminalOutcome{}, fmt.Errorf("%w: a positive log bound is required", ErrInvalidAgentRunnerRun)
	}
	clock := run.Clock
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}
	// The pre-run manifest must be captured before the CLI can touch the workspace:
	// a manifest captured afterwards would make the diff claim "no change" for a
	// run that did change it. A failed capture stays a gap and degrades the outcome
	// instead of aborting it, so the caller always receives a reason.
	var beforeManifest []byte
	if request.Capturer != nil {
		captured, captureErr := request.Capturer.CaptureWorkspaceManifest(ctx, agentRunnerWorkspaceRoot(request.Launch))
		if captureErr == nil {
			beforeManifest = captured
		}
	}
	launchResult, err := run.Launcher.StartAgentRunnerProcess(ctx, request.Launch)
	if err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	watch, processResult, err := run.Launcher.WaitAgentRunnerProcess(ctx, launchResult.LaunchID, request.Launch.RequestID, request.Timeout)
	if err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	// The classification must complete even when the caller cancelled: the terminal
	// receipt is what records the cancellation, so the caller's context must not be
	// able to abort the mapping half-way.
	classifyCtx := context.WithoutCancel(ctx)
	stopProof := processResult.Termination
	switch {
	case agentRunnerNeedsDurableStop(watch, processResult):
		// A run the watchdog could not settle may still be alive, and a natural exit
		// still owes a proven stop: the reviewed stop path returns "already-stopped"
		// evidence for a process that already exited, which is the process-stop fact
		// a completed terminal must carry. Retiring the launch without going through
		// it would let a second launch overwrite the identity of a live process.
		proof, stopErr := run.Launcher.StopAgentRunnerProcess(classifyCtx, launchResult.LaunchID)
		if stopErr != nil {
			return AgentRunnerTerminalOutcome{}, stopErr
		}
		stopProof = &proof
	default:
		// The guard stopped it and proved the tree gone, so the launch is only retired
		// (which closes the log artifacts before anything reads them).
		run.Launcher.releaseLaunch(launchResult.LaunchID)
	}
	evidenceDir, err := run.Store.RunEvidenceDir(request.Launch.RequestID)
	if err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	collected, err := CollectAgentRunnerEvidence(classifyCtx, AgentRunnerEvidenceSpec{
		EvidenceDir:      evidenceDir,
		WorkspaceRoot:    agentRunnerWorkspaceRoot(request.Launch),
		RequestID:        request.Launch.RequestID,
		Capturer:         request.Capturer,
		StopEvidenceHash: stopProof.Hash,
		BeforeManifest:   beforeManifest,
		MaxLogBytes:      request.MaxLogBytes,
	})
	if err != nil {
		return AgentRunnerTerminalOutcome{}, err
	}
	// A stop is only "stopped" once the tree is verified gone: on Windows the guard's
	// own verification is a no-op because job-close is asynchronous, so the run
	// re-checks the child processes the CLI reported - for a naturally exited run
	// and for a run the guard stopped, because in both cases a descendant may still
	// be dying when the receipt is written.
	treeGone := stopProof.Outcome == "stopped"
	if treeGone {
		treeGone, err = agentRunnerVerifyTreeGone(classifyCtx, evidenceDir, agentRunnerVerifyWait(request.VerifyWait))
		if err != nil {
			return AgentRunnerTerminalOutcome{}, err
		}
	}
	status, treeStatus, outcomeEvidence := agentRunnerTerminalStatus(request.Launch.RequestID, watch, processResult, stopProof, collected, treeGone)
	outcome := AgentRunnerTerminalOutcome{
		Completion: agentRunnerCompletionFor(request, processResult, collected, status, treeStatus, outcomeEvidence, clock()),
		Settlement: agentRunnerSettlementFor(status, collected),
	}
	if status != "completed" {
		return outcome, nil
	}
	facts, err := collected.FrozenFacts()
	if err != nil {
		// A completion claim without provable facts is refused instead of being
		// reported as a pass, because the chain would reject it anyway.
		return AgentRunnerTerminalOutcome{}, err
	}
	outcome.Facts = &facts
	return outcome, nil
}

func agentRunnerWorkspaceRoot(launch chain.AgentRunnerLaunchRequest) string {
	if launch.WorkspaceRoot != "" {
		return launch.WorkspaceRoot
	}
	return launch.Dir
}

// agentRunnerNeedsDurableStop reports whether the outcome still has to go through
// the durable identity stop. A natural exit owes the "already-stopped" proof, and a
// run the watchdog could not settle may still be running: in both cases the launch
// may only be retired after the reviewed stop path has spoken.
func agentRunnerNeedsDurableStop(watch AgentRunnerWatchResult, result guard.ProcessResult) bool {
	if result.Termination == nil {
		return true
	}
	// A settled guard stop is complete by itself, but an unsettled classification
	// means the run side never reconciled: its placeholder proof is not a stop.
	return errors.Is(watch.Err, ErrAgentRunnerLaunchUnsettled)
}

// agentRunnerTerminalStatus maps what happened into the terminal status, the
// process-tree status and the error evidence. Only a clean exit with complete
// evidence may be completed: an operator request, a timeout, a cancellation, an
// unproven stop, a failure and any gap all degrade the outcome instead.
func agentRunnerTerminalStatus(requestID string, watch AgentRunnerWatchResult, result guard.ProcessResult, stopProof *guard.TerminationEvidence, collected AgentRunnerEvidence, treeGone bool) (string, string, []string) {
	treeStatus := "unknown"
	if stopProof != nil && stopProof.Outcome == "stopped" && treeGone {
		treeStatus = "stopped"
	}
	// A natural failure is the run's own non-zero exit status; every other non-nil
	// watchdog error is a stop the run could not classify.
	naturalFailure := result.Termination == nil && result.ExitCode != 0
	status := ""
	switch {
	case collected.OperatorAsk:
		status = "operator-action-required"
	case watch.TimedOut:
		status = "uncertain"
	case errors.Is(watch.Err, context.Canceled) || errors.Is(watch.Err, context.DeadlineExceeded):
		// A caller cancellation and a caller deadline expiry are both "we were asked
		// to stop": neither may be reported as a failed or completed run.
		status = "cancelled"
	case watch.Err != nil && !naturalFailure:
		// The guard's residue and uncertainty sentinels, an unsettled run and any
		// error the run cannot attribute: fail closed rather than letting an unknown
		// stop reason fall through to the completion branch.
		status = "uncertain"
	case naturalFailure:
		// A natural non-zero exit is a first-class fact and outranks every gap: the tree
		// status has its own field on the receipt, so an unproven tree must not hide the
		// failure. Gaps can only block a completion.
		status = "failed"
	case treeStatus != "stopped":
		status = "uncertain"
	case result.Termination != nil:
		status = "uncertain"
	case !collected.LogsComplete || !collected.UsageComplete || !collected.EventsComplete || !collected.ManifestsComplete:
		status = "uncertain"
	default:
		status = "completed"
	}
	if status == "completed" {
		return status, treeStatus, nil
	}
	errorEvidence := append([]string(nil), collected.ErrorEvidence...)
	errorEvidence = append(errorEvidence, watch.Evidence...)
	// A degraded receipt must cite the stop evidence it archived: a stop that is
	// only described inside the outcome struct is not archived evidence.
	if evidence.ValidHash(collected.ProcessStopEvidenceHash) {
		errorEvidence = append(errorEvidence, collected.ProcessStopEvidenceHash)
	}
	// A non-completed receipt must carry provable error evidence even when the
	// collection found no artifact gap: the reason itself is the evidence.
	errorEvidence = append(errorEvidence, agentRunnerOutcomeEvidenceFor(requestID, status))
	return status, treeStatus, sortedUniqueAgentRunnerHashes(errorEvidence)
}

// agentRunnerVerifyWait bounds the tree verification. It is a tolerance rather
// than a correctness input, so a zero value falls back to a default.
func agentRunnerVerifyWait(wait time.Duration) time.Duration {
	if wait <= 0 {
		return 2 * time.Second
	}
	return wait
}

// agentRunnerVerifyTreeGone checks the child processes the CLI reported: a
// natural exit is only trustworthy when nothing it spawned is still running.
func agentRunnerVerifyTreeGone(ctx context.Context, evidenceDir string, wait time.Duration) (bool, error) {
	return agentRunnerVerifyTreeGoneWith(ctx, evidenceDir, wait, guard.InspectProcess, guard.ProcessAlive)
}

// agentRunnerVerifyTreeGoneWith takes the process observations as parameters: the
// liveness decision is the whole point of the verification, so it must be
// testable without racing a real process against its own exit.
func agentRunnerVerifyTreeGoneWith(ctx context.Context, evidenceDir string, wait time.Duration, inspect func(int) (guard.ProcessIdentity, error), alive func(guard.ProcessIdentity) (bool, error)) (bool, error) {
	report := filepath.Join(evidenceDir, agentRunnerChildPIDReportFileName)
	deadline := time.Now().Add(wait)
	for {
		pids, err := readAgentRunnerChildPIDs(report)
		if err != nil {
			return false, err
		}
		stillAlive := false
		for _, pid := range pids {
			identity, err := inspect(pid)
			if err != nil {
				// A pid that cannot be observed is treated as gone: a CLI that spawns a
				// child and waits for it must not degrade merely because the child has
				// already exited. The platform containment (job close on Windows, the
				// process group elsewhere) stays the primary guarantee.
				continue
			}
			if running, err := alive(identity); err != nil || running {
				stillAlive = true
				break
			}
		}
		if !stillAlive {
			return true, nil
		}
		// A child may still be dying, so the answer is re-checked until the bound
		// expires; only then is the tree reported as unproven.
		if time.Now().After(deadline) {
			return false, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// readAgentRunnerChildPIDs reads the one-pid-per-line report the CLI writes. A
// missing report means it spawned nothing.
func readAgentRunnerChildPIDs(path string) ([]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: read %s: %v", ErrInvalidAgentRunnerRun, path, err)
	}
	pids := make([]int, 0, 1)
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil || pid <= 0 {
			return nil, fmt.Errorf("%w: unreadable child pid %q", ErrInvalidAgentRunnerRun, line)
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

// agentRunnerOutcomeEvidenceFor derives the digest that names one degraded outcome.
func agentRunnerOutcomeEvidenceFor(requestID, reason string) string {
	return evidence.Digest(agentRunnerOutcomeEvidenceDomain, []byte(requestID+"\n"+reason))
}

// sortedUniqueAgentRunnerHashes sorts and deduplicates digests, which is what a
// completion receipt requires for both evidence lists.
func sortedUniqueAgentRunnerHashes(values []string) []string {
	sorted := append([]string(nil), values...)
	slices.Sort(sorted)
	return slices.Compact(sorted)
}

// agentRunnerCompletionFor builds the completion claim. A completed claim carries
// the output manifest digest; every other status carries error evidence instead.
func agentRunnerCompletionFor(request AgentRunnerRunRequest, result guard.ProcessResult, collected AgentRunnerEvidence, status, treeStatus string, errorEvidence []string, completedAt time.Time) AgentRunnerCompletion {
	exitCode := result.ExitCode
	completion := AgentRunnerCompletion{
		Kind:              "agent-runner-completion",
		CompletionID:      "completion-" + request.Launch.RequestID,
		RequestID:         request.Launch.RequestID,
		RequestHash:       request.RequestHash,
		RunID:             request.Launch.RunID,
		TaskID:            request.TaskID,
		StepID:            request.StepID,
		Attempt:           request.Attempt,
		AdapterID:         request.Launch.AdapterID,
		SessionID:         "session-" + request.Launch.RequestID,
		CompletedAt:       completedAt.UTC().Format(TimestampLayout),
		Status:            status,
		ProcessTreeStatus: treeStatus,
		LogsComplete:      collected.LogsComplete,
		UsageComplete:     collected.UsageComplete,
		Evidence:          []string{request.RequestHash},
		ErrorEvidence:     errorEvidence,
	}
	if result.Termination == nil {
		// The process exited on its own, so the exit status is a fact about the run
		// rather than a placeholder for "we stopped it".
		completion.ExitCode = &exitCode
	}
	if status == "completed" {
		manifest := collected.ManifestHash
		completion.OutputManifestHash = &manifest
	}
	return completion
}

// agentRunnerSettlementFor maps the outcome to the usage observation. An
// observation that was never proven settles as unknown: an incomplete usage
// artifact may not be settled as observed.
func agentRunnerSettlementFor(status string, collected AgentRunnerEvidence) AgentRunnerTerminalSettlement {
	if status != "completed" || !collected.UsageComplete {
		return AgentRunnerTerminalSettlement{Status: tickets.CostSettlementUnknown}
	}
	calls := collected.Usage.Calls
	tokens := collected.Usage.Tokens
	return AgentRunnerTerminalSettlement{
		Status:           tickets.CostSettlementObserved,
		ObservedCalls:    &calls,
		ObservedTokens:   &tokens,
		ProviderEvidence: []string{collected.UsageHash},
	}
}
