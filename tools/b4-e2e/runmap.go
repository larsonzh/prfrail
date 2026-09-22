package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/tickets"
)

// chainTerminal is one built terminal chain: the three durable records plus the
// translated chain-owned fact the engine routes.
type chainTerminal struct {
	fact       chain.AgentRunnerTerminal
	intent     adapters.AgentRunnerTerminalIntentRecord
	closure    adapters.AgentRunnerTerminalClosureRecord
	completion adapters.AgentRunnerCompletionRecord
}

// This file reproduces four unexported production functions, because the harness
// must wait for a run the *dispatcher* started (AgentRunnerPinnedCLIRun.Run starts
// its own process, so it cannot be used after DispatchAgentRunner).
//
// Mirrored production code, byte-class-isomorphic in behaviour:
//
//	agentRunnerNeedsDurableStop   internal/adapters/agent_runner_run.go:177
//	agentRunnerTerminalStatus     internal/adapters/agent_runner_run.go:190
//	agentRunnerVerifyWait         internal/adapters/agent_runner_run.go:245
//	agentRunnerVerifyTreeGoneWith internal/adapters/agent_runner_run.go:261
//	readAgentRunnerChildPIDs      internal/adapters/agent_runner_run.go:302
//	agentRunnerOutcomeEvidenceFor internal/adapters/agent_runner_run.go:325
//	sortedUniqueAgentRunnerHashes internal/adapters/agent_runner_run.go:331
//	agentRunnerCompletionFor      internal/adapters/agent_runner_run.go:339
//	agentRunnerSettlementFor      internal/adapters/agent_runner_run.go:375
//
// Two deliberate, documented divergences (both recorded in the verdict):
//
//	(a) the natural-exit branch of AgentRunnerPinnedCLIRun.Run retires the launch
//	    through the unexported releaseLaunch, which this process cannot call; the
//	    harness stops through the reviewed public stop entry instead, which returns
//	    the same already-stopped proof for an exited process and also retires it.
//	(b) agentRunnerSettlementFor leaves ChargedAmountMicros to the caller; the
//	    harness supplies the zero-cost amount of the declared offline workload.

// runObservation is everything one managed run produced.
type runObservation struct {
	LaunchID    string
	RequestID   string
	Elapsed     time.Duration
	Watch       adapters.AgentRunnerWatchResult
	Process     guard.ProcessResult
	StopProof   *guard.TerminationEvidence
	StopError   string
	Collected   adapters.AgentRunnerEvidence
	CollectErr  string
	Status      string
	TreeStatus  string
	ErrorEvi    []string
	Completion  adapters.AgentRunnerCompletion
	Settlement  adapters.AgentRunnerTerminalSettlement
	Facts       *adapters.AgentRunnerFrozenFacts
	CollectedAt time.Time
}

// mirrorNeedsDurableStop mirrors agentRunnerNeedsDurableStop.
func mirrorNeedsDurableStop(watch adapters.AgentRunnerWatchResult, result guard.ProcessResult) bool {
	if result.Termination == nil {
		return true
	}
	return errors.Is(watch.Err, adapters.ErrAgentRunnerLaunchUnsettled)
}

// mirrorVerifyWait mirrors agentRunnerVerifyWait.
func mirrorVerifyWait(wait time.Duration) time.Duration {
	if wait <= 0 {
		return 2 * time.Second
	}
	return wait
}

// mirrorVerifyTreeGone mirrors agentRunnerVerifyTreeGoneWith with the production
// process observations (guard.InspectProcess / guard.ProcessAlive).
func mirrorVerifyTreeGone(ctx context.Context, evidenceDir string, wait time.Duration) (bool, error) {
	report := filepath.Join(evidenceDir, "stub-child.pid")
	deadline := time.Now().Add(wait)
	for {
		pids, err := readPIDs(report)
		if err != nil {
			return false, err
		}
		stillAlive := false
		for _, pid := range pids {
			identity, err := guard.InspectProcess(pid)
			if err != nil {
				continue
			}
			if running, err := guard.ProcessAlive(identity); err != nil || running {
				stillAlive = true
				break
			}
		}
		if !stillAlive {
			return true, nil
		}
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

// mirrorOutcomeEvidenceFor mirrors agentRunnerOutcomeEvidenceFor.
func mirrorOutcomeEvidenceFor(requestID, reason string) string {
	return evidence.Digest("proofrail:agent-runner-outcome:1\n", []byte(requestID+"\n"+reason))
}

// mirrorSortedUnique mirrors sortedUniqueAgentRunnerHashes.
func mirrorSortedUnique(values []string) []string {
	sorted := append([]string(nil), values...)
	slices.Sort(sorted)
	return slices.Compact(sorted)
}

// mirrorTerminalStatus mirrors agentRunnerTerminalStatus.
func mirrorTerminalStatus(requestID string, watch adapters.AgentRunnerWatchResult, result guard.ProcessResult, stopProof *guard.TerminationEvidence, collected adapters.AgentRunnerEvidence, treeGone bool) (string, string, []string) {
	treeStatus := "unknown"
	if stopProof != nil && stopProof.Outcome == "stopped" && treeGone {
		treeStatus = "stopped"
	}
	naturalFailure := result.Termination == nil && result.ExitCode != 0
	status := ""
	switch {
	case collected.OperatorAsk:
		status = "operator-action-required"
	case watch.TimedOut:
		status = "uncertain"
	case errors.Is(watch.Err, context.Canceled) || errors.Is(watch.Err, context.DeadlineExceeded):
		status = "cancelled"
	case watch.Err != nil && !naturalFailure:
		status = "uncertain"
	case naturalFailure:
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
	if evidence.ValidHash(collected.ProcessStopEvidenceHash) {
		errorEvidence = append(errorEvidence, collected.ProcessStopEvidenceHash)
	}
	errorEvidence = append(errorEvidence, mirrorOutcomeEvidenceFor(requestID, status))
	return status, treeStatus, mirrorSortedUnique(errorEvidence)
}

// mirrorCompletionFor mirrors agentRunnerCompletionFor.
func mirrorCompletionFor(request adapters.AgentRunnerRunRequest, result guard.ProcessResult, collected adapters.AgentRunnerEvidence, status, treeStatus string, errorEvidence []string, completedAt time.Time) adapters.AgentRunnerCompletion {
	exitCode := result.ExitCode
	completion := adapters.AgentRunnerCompletion{
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
		CompletedAt:       completedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		Status:            status,
		ProcessTreeStatus: treeStatus,
		LogsComplete:      collected.LogsComplete,
		UsageComplete:     collected.UsageComplete,
		Evidence:          []string{request.RequestHash},
		ErrorEvidence:     errorEvidence,
	}
	if result.Termination == nil {
		completion.ExitCode = &exitCode
	}
	if status == "completed" {
		manifest := collected.ManifestHash
		completion.OutputManifestHash = &manifest
	}
	return completion
}

// mirrorSettlementFor mirrors agentRunnerSettlementFor and then applies the
// documented divergence (b): the zero-cost charged amount of this workload.
func mirrorSettlementFor(status string, collected adapters.AgentRunnerEvidence) adapters.AgentRunnerTerminalSettlement {
	if status != "completed" || !collected.UsageComplete {
		return adapters.AgentRunnerTerminalSettlement{Status: tickets.CostSettlementUnknown}
	}
	calls := collected.Usage.Calls
	tokens := collected.Usage.Tokens
	zero := int64(0)
	return adapters.AgentRunnerTerminalSettlement{
		Status:              tickets.CostSettlementObserved,
		ChargedAmountMicros: &zero,
		ObservedCalls:       &calls,
		ObservedTokens:      &tokens,
		ProviderEvidence:    []string{collected.UsageHash},
	}
}

// waitAndMap waits for the dispatched process, collects the artifacts, maps the
// outcome exactly as the production run does, and returns the observation.
func waitAndMap(ctx context.Context, launcher *adapters.AgentRunnerPinnedCLILauncher, store *adapters.AgentRunnerReplayStore, capturer adapters.AgentRunnerManifestCapturer, runID, requestID string, taskID, stepID string, requestHash string, workspaceRoot string, beforeManifest []byte, timeout, verifyWait time.Duration, maxLogBytes int64, clock func() time.Time) (*runObservation, error) {
	observation := &runObservation{RequestID: requestID, LaunchID: "launch-" + requestID}
	started := time.Now()
	watch, processResult, err := launcher.WaitAgentRunnerProcess(ctx, observation.LaunchID, requestID, timeout)
	observation.Elapsed = time.Since(started)
	observation.Watch = watch
	if err != nil {
		return observation, err
	}
	observation.Process = processResult
	classifyCtx := context.WithoutCancel(ctx)
	stopProof := processResult.Termination
	switch {
	case mirrorNeedsDurableStop(watch, processResult):
		// Divergence (a): the production branch retires the launch without stopping
		// when the guard already settled the kill; the harness always goes through
		// the reviewed public stop entry, which returns the already-stopped proof
		// for an exited process and retires the launch either way.
		proof, stopErr := launcher.StopAgentRunnerProcess(classifyCtx, observation.LaunchID)
		observation.StopProof = &proof
		stopProof = &proof
		if stopErr != nil {
			observation.StopError = stopErr.Error()
			return observation, stopErr
		}
	default:
		// The guard stopped the tree and proved it gone: the production path only
		// retires the launch here (unexported releaseLaunch). Nothing is stopped a
		// second time; the guard's own termination proof is the stop proof.
	}
	evidenceDir, err := store.RunEvidenceDir(requestID)
	if err != nil {
		return observation, err
	}
	collected, err := adapters.CollectAgentRunnerEvidence(classifyCtx, adapters.AgentRunnerEvidenceSpec{
		EvidenceDir:      evidenceDir,
		WorkspaceRoot:    workspaceRoot,
		RequestID:        requestID,
		Capturer:         capturer,
		StopEvidenceHash: stopProofHash(stopProof),
		BeforeManifest:   beforeManifest,
		MaxLogBytes:      maxLogBytes,
	})
	if err != nil {
		observation.CollectErr = err.Error()
		return observation, err
	}
	observation.Collected = collected
	treeGone := stopProof != nil && stopProof.Outcome == "stopped"
	if treeGone {
		treeGone, err = mirrorVerifyTreeGone(classifyCtx, evidenceDir, mirrorVerifyWait(verifyWait))
		if err != nil {
			return observation, err
		}
	}
	status, treeStatus, outcomeEvidence := mirrorTerminalStatus(requestID, watch, processResult, stopProof, collected, treeGone)
	runRequest := adapters.AgentRunnerRunRequest{
		Launch:      chain.AgentRunnerLaunchRequest{RequestID: requestID, RunID: runID, AdapterID: declaredAdapterID},
		RequestHash: requestHash,
		TaskID:      taskID,
		StepID:      stepID,
		Attempt:     1,
	}
	observation.Status = status
	observation.TreeStatus = treeStatus
	observation.ErrorEvi = outcomeEvidence
	observation.CollectedAt = clock()
	observation.Completion = mirrorCompletionFor(runRequest, processResult, collected, status, treeStatus, outcomeEvidence, observation.CollectedAt)
	observation.Settlement = mirrorSettlementFor(status, collected)
	if status != "completed" {
		return observation, nil
	}
	facts, err := collected.FrozenFacts()
	if err != nil {
		return observation, err
	}
	// The facts the collector produces are the adapter-side shape; the chain-side
	// shape is the same five digests.
	observation.Facts = &adapters.AgentRunnerFrozenFacts{
		ManifestHash:            facts.ManifestHash,
		DiffHash:                facts.DiffHash,
		LogHash:                 facts.LogHash,
		UsageHash:               facts.UsageHash,
		ProcessStopEvidenceHash: facts.ProcessStopEvidenceHash,
	}
	return observation, nil
}

func stopProofHash(proof *guard.TerminationEvidence) string {
	if proof == nil {
		return ""
	}
	return proof.Hash
}

// terminalChain builds the durable terminal chain (completion, intent, closure)
// through the public constructors and translates it into the chain-owned terminal
// fact. It is the alternative-evidence path for the publication leg: the
// production AgentRunnerTerminalPublisher refuses a fresh publication on a
// platform whose replay-store durability is unproven (measured on Windows), and
// this harness cannot change that platform declaration.
func terminalChain(inputs *declaredInputs, observation *runObservation, decidedAt time.Time) (chainTerminal, error) {
	request := inputs.record
	completionRecord, err := adapters.NewAgentRunnerCompletionRecord(observation.Completion)
	if err != nil {
		return chainTerminal{}, fmt.Errorf("completion record: %w", err)
	}
	providerEvidence := []string{completionRecord.RecordHash, request.RecordHash}
	providerEvidence = append(providerEvidence, observation.Settlement.ProviderEvidence...)
	intent := adapters.AgentRunnerTerminalIntent{
		RequestID:                request.Request.RequestID,
		RequestHash:              request.RecordHash,
		RunID:                    request.Request.RunID,
		TaskID:                   request.Request.TaskID,
		StepID:                   request.Request.StepID,
		Attempt:                  request.Request.Attempt,
		AdapterID:                request.Request.AdapterID,
		Completion:               completionRecord,
		SettlementEntryID:        "settle-" + completionRecord.RecordHash[len("sha256:"):len("sha256:")+20],
		SettlementIdempotencyKey: "settle-" + request.Request.RequestID,
		ReservationHash:          request.Request.BudgetHash,
		SettlementStatus:         observation.Settlement.Status,
		ChargedAmountMicros:      observation.Settlement.ChargedAmountMicros,
		ObservedCalls:            observation.Settlement.ObservedCalls,
		ObservedTokens:           observation.Settlement.ObservedTokens,
		ProviderEvidence:         dedupe(providerEvidence),
		DecidedAt:                decidedAt.UTC().Format(time.RFC3339),
		Facts:                    observation.Facts,
	}
	intentRecord, err := adapters.NewAgentRunnerTerminalIntentRecord(intent)
	if err != nil {
		return chainTerminal{}, fmt.Errorf("terminal intent record: %w", err)
	}
	closureRecord, err := adapters.NewAgentRunnerTerminalClosureRecord(adapters.AgentRunnerTerminalClosure{
		RequestID:                intent.RequestID,
		IntentRecordHash:         intentRecord.RecordHash,
		CompletionHash:           completionRecord.RecordHash,
		SettlementEntryID:        intent.SettlementEntryID,
		SettlementIdempotencyKey: intent.SettlementIdempotencyKey,
		ClosedAt:                 intent.DecidedAt,
	})
	if err != nil {
		return chainTerminal{}, fmt.Errorf("terminal closure record: %w", err)
	}
	fact, err := adapters.ToChainAgentRunnerTerminal(request, intentRecord, closureRecord)
	if err != nil {
		return chainTerminal{}, fmt.Errorf("translate terminal: %w", err)
	}
	return chainTerminal{fact: fact, intent: intentRecord, closure: closureRecord, completion: completionRecord}, nil
}

func dedupe(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
