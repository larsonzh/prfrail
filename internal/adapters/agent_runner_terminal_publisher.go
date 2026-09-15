package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/tickets"
)

// AgentRunnerTerminalSettlement is the caller-observed usage/cost fact for one
// terminal publication. Status is one of the cost-ledger settlement statuses.
type AgentRunnerTerminalSettlement struct {
	Status              string
	ChargedAmountMicros *int64
	ObservedCalls       *int
	ObservedTokens      *int
	ProviderEvidence    []string
}

// AgentRunnerTerminalOutcome carries the completion claim, the usage
// observation and (for a completed claim) the frozen facts a caller wants to
// terminalize.
type AgentRunnerTerminalOutcome struct {
	Completion AgentRunnerCompletion
	Settlement AgentRunnerTerminalSettlement
	Facts      *AgentRunnerFrozenFacts
}

// AgentRunnerTerminalResult reports the durable chain hashes that prove the
// terminal decision. Replayed is true when a stored terminal decision was
// converged instead of newly decided.
type AgentRunnerTerminalResult struct {
	Evidence []string
	Replayed bool
}

// AgentRunnerTerminalPublisher owns the terminal publication protocol: publish
// the terminal-intent, settle the cost plan, publish C, and close the chain.
// It is invoked after dispatch; it never launches processes, never wires the
// Engine, and never publishes a completion without a settlement decision
// record. Store-local records are the only recovery source; type-safe value
// semantics keep repeated calls idempotent.
type AgentRunnerTerminalPublisher struct {
	Store         *AgentRunnerReplayStore
	CostLedger    *tickets.CostLedger
	RequestRecord AgentRunnerRequestRecord
	Clock         func() time.Time

	// Test-only single-shot probes. beforeSettle fires at the single settlement
	// choke point, so any branch that skips settlement can be proven write-free.
	beforeSettle     func() error
	afterIntentWrite func() error
	afterSettle      func() error
}

// PublishTerminal publishes or converges the terminal chain for one request.
// Fresh decisions run intent -> settle -> C -> closure; replays re-verify the
// stored decision and only complete missing chain steps idempotently.
func (publisher *AgentRunnerTerminalPublisher) PublishTerminal(ctx context.Context, outcome AgentRunnerTerminalOutcome) (AgentRunnerTerminalResult, error) {
	if publisher == nil {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: nil terminal publisher", ErrInvalidAgentRunnerRequest)
	}
	if err := ctx.Err(); err != nil {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: terminal publication context cancelled: %w", ErrInvalidAgentRunnerRequest, err)
	}
	if publisher.Store == nil {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: replay store missing", ErrInvalidAgentRunnerRequest)
	}
	if publisher.CostLedger == nil {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: cost ledger missing", ErrInvalidAgentRunnerRequest)
	}
	if publisher.Clock == nil {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: evaluation clock missing", ErrInvalidAgentRunnerRequest)
	}
	if err := ValidateAgentRunnerRequestRecord(publisher.RequestRecord); err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	now := publisher.Clock().UTC()
	if now.IsZero() {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: evaluation clock returned zero time", ErrInvalidAgentRunnerRequest)
	}
	request := publisher.RequestRecord.Request
	completion := outcome.Completion
	if completion.CompletionID == "" {
		completion.CompletionID = "completion-" + request.RequestID
	}
	completionRecord, err := NewAgentRunnerCompletionRecord(completion)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	if err := validateAgentRunnerTerminalCompletionBinding(publisher.RequestRecord, completionRecord); err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	plan, err := terminalSettlementIntentFields(publisher.RequestRecord, completionRecord, outcome.Settlement)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	if err := attachTerminalIntentFacts(&plan, completionRecord, outcome.Facts); err != nil {
		return AgentRunnerTerminalResult{}, err
	}

	state, err := publisher.Store.State(request.RequestID)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	intent, foundIntent, err := publisher.Store.TerminalIntent(request.RequestID)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	closure, foundClosure, err := publisher.Store.TerminalClosure(request.RequestID)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	if foundIntent {
		return publisher.convergeTerminal(intent, closure, foundClosure, completionRecord, plan)
	}
	if foundClosure {
		return AgentRunnerTerminalResult{}, publisher.terminalCorruptionError(request.RequestID, "terminal closure present without a terminal intent")
	}
	switch state {
	case AgentRunnerReplayStateAbsent:
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: terminal publication requires a persisted request record for requestId %s", ErrAgentRunnerRequestConflict, request.RequestID)
	case AgentRunnerReplayStateTerminalReceiptPresent:
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: completion %s carries no settlement decision record", ErrAgentRunnerTerminalOrphan, request.RequestID)
	case AgentRunnerReplayStateDispatchedUnknown:
	default:
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: unknown replay state %s", ErrInvalidAgentRunnerRequest, state)
	}
	// Durability is a pre-settlement gate: an unproven store must never accept a
	// new terminal decision, because neither the intent nor the completion could
	// be published afterwards.
	if publisher.Store.PublicationDurability() != PublishDurabilityProven {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: terminal publication", ErrAgentRunnerReplayStoreDurabilityUnproven)
	}
	if err := validateAgentRunnerBudget(publisher.CostLedger, request); err != nil {
		// A concurrent publisher may have settled this reservation between our
		// intent read and this precheck. If the winner's intent is now visible,
		// converge on it instead of failing the whole publication.
		stored, found, loadErr := publisher.Store.TerminalIntent(request.RequestID)
		if loadErr != nil {
			return AgentRunnerTerminalResult{}, loadErr
		}
		if found {
			storedClosure, storedClosureFound, closureErr := publisher.Store.TerminalClosure(request.RequestID)
			if closureErr != nil {
				return AgentRunnerTerminalResult{}, closureErr
			}
			return publisher.convergeTerminal(stored, storedClosure, storedClosureFound, completionRecord, plan)
		}
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: %w", ErrAgentRunnerTerminalSettlementBlocked, err)
	}
	intentRecord, err := NewAgentRunnerTerminalIntentRecord(AgentRunnerTerminalIntent{
		RequestID:                request.RequestID,
		RequestHash:              publisher.RequestRecord.RecordHash,
		RunID:                    request.RunID,
		TaskID:                   request.TaskID,
		StepID:                   request.StepID,
		Attempt:                  request.Attempt,
		AdapterID:                request.AdapterID,
		Completion:               completionRecord,
		Facts:                    plan.Facts,
		SettlementEntryID:        plan.SettlementEntryID,
		SettlementIdempotencyKey: plan.SettlementIdempotencyKey,
		ReservationHash:          plan.ReservationHash,
		SettlementStatus:         plan.SettlementStatus,
		ChargedAmountMicros:      plan.ChargedAmountMicros,
		ObservedCalls:            plan.ObservedCalls,
		ObservedTokens:           plan.ObservedTokens,
		ProviderEvidence:         plan.ProviderEvidence,
		DecidedAt:                now.Format(time.RFC3339),
	})
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	decision, err := publisher.Store.RecordTerminalIntent(intentRecord)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	if decision == AgentRunnerTerminalIntentDecisionPresent {
		stored, found, err := publisher.Store.TerminalIntent(request.RequestID)
		if err != nil {
			return AgentRunnerTerminalResult{}, err
		}
		if !found {
			return AgentRunnerTerminalResult{}, fmt.Errorf("%w: terminal intent winner not visible for requestId %s", ErrAgentRunnerReplayStoreConvergence, request.RequestID)
		}
		storedClosure, storedClosureFound, err := publisher.Store.TerminalClosure(request.RequestID)
		if err != nil {
			return AgentRunnerTerminalResult{}, err
		}
		return publisher.convergeTerminal(stored, storedClosure, storedClosureFound, completionRecord, plan)
	}
	if publisher.afterIntentWrite != nil {
		if err := publisher.afterIntentWrite(); err != nil {
			return AgentRunnerTerminalResult{}, err
		}
	}
	if err := publisher.settleTerminalIntent(intentRecord.Intent); err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	if publisher.afterSettle != nil {
		if err := publisher.afterSettle(); err != nil {
			return AgentRunnerTerminalResult{}, err
		}
	}
	if _, err := publisher.Store.RecordCompletion(publisher.RequestRecord, intentRecord.Intent.Completion); err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	closureRecord, err := buildAgentRunnerTerminalClosure(intentRecord)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	replayed, err := publisher.Store.RecordTerminalClosure(closureRecord)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	return AgentRunnerTerminalResult{
		Evidence: []string{intentRecord.RecordHash, intentRecord.Intent.Completion.RecordHash, closureRecord.RecordHash},
		Replayed: replayed,
	}, nil
}

// convergeTerminal re-verifies the stored terminal decision against the caller's
// claim and idempotently completes any missing chain steps.
func (publisher *AgentRunnerTerminalPublisher) convergeTerminal(intent AgentRunnerTerminalIntentRecord, closure AgentRunnerTerminalClosureRecord, foundClosure bool, completionRecord AgentRunnerCompletionRecord, plan AgentRunnerTerminalIntent) (AgentRunnerTerminalResult, error) {
	body := intent.Intent
	state, err := publisher.Store.State(body.RequestID)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	if state == AgentRunnerReplayStateAbsent {
		// Without a persisted request record the recovery would silently recreate
		// R after booking the settlement; refuse instead.
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: terminal recovery requires the persisted request record for requestId %s", ErrAgentRunnerRequestConflict, body.RequestID)
	}
	if body.Completion.RecordHash != completionRecord.RecordHash {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: terminal slot for requestId %s already decided completion %s", ErrAgentRunnerTerminalIntentConflict, body.RequestID, body.Completion.RecordHash)
	}
	if !sameTerminalSettlementPlan(body, plan) {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: terminal slot for requestId %s already decided a different settlement plan", ErrAgentRunnerTerminalIntentConflict, body.RequestID)
	}
	if foundClosure {
		if err := publisher.Store.validateTerminalClosureBinding(intent, closure.Closure); err != nil {
			return AgentRunnerTerminalResult{}, publisher.terminalCorruptionError(body.RequestID, "terminal closure does not bind the stored terminal intent")
		}
		if state != AgentRunnerReplayStateTerminalReceiptPresent {
			// A closure without a visible completion contradicts the publication
			// order and cannot be repaired by replaying writes.
			return AgentRunnerTerminalResult{}, publisher.terminalCorruptionError(body.RequestID, "terminal closure present without a published completion")
		}
		if publisher.Store.PublicationDurability() != PublishDurabilityProven {
			// An unproven store allows only write-free replay: the closure proves the
			// settlement was booked in protocol order, so no ledger call is made.
			return AgentRunnerTerminalResult{
				Evidence: []string{intent.RecordHash, body.Completion.RecordHash, closure.RecordHash},
				Replayed: true,
			}, nil
		}
		if err := publisher.settleTerminalIntent(body); err != nil {
			return AgentRunnerTerminalResult{}, err
		}
		return AgentRunnerTerminalResult{
			Evidence: []string{intent.RecordHash, body.Completion.RecordHash, closure.RecordHash},
			Replayed: true,
		}, nil
	}
	// The chain is incomplete, so publishing the remaining steps requires a
	// proven store; recovery must never book a settlement it cannot close.
	if publisher.Store.PublicationDurability() != PublishDurabilityProven {
		return AgentRunnerTerminalResult{}, fmt.Errorf("%w: terminal recovery", ErrAgentRunnerReplayStoreDurabilityUnproven)
	}
	if err := publisher.settleTerminalIntent(body); err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	if publisher.afterSettle != nil {
		if err := publisher.afterSettle(); err != nil {
			return AgentRunnerTerminalResult{}, err
		}
	}
	if _, err := publisher.Store.RecordCompletion(publisher.RequestRecord, body.Completion); err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	closureRecord, err := buildAgentRunnerTerminalClosure(intent)
	if err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	if _, err := publisher.Store.RecordTerminalClosure(closureRecord); err != nil {
		return AgentRunnerTerminalResult{}, err
	}
	return AgentRunnerTerminalResult{
		Evidence: []string{intent.RecordHash, body.Completion.RecordHash, closureRecord.RecordHash},
		Replayed: true,
	}, nil
}

func (publisher *AgentRunnerTerminalPublisher) settleTerminalIntent(intent AgentRunnerTerminalIntent) error {
	if publisher.beforeSettle != nil {
		if err := publisher.beforeSettle(); err != nil {
			return err
		}
	}
	occurredAt, err := time.Parse(time.RFC3339, intent.DecidedAt)
	if err != nil {
		return fmt.Errorf("%w: invalid decidedAt", ErrInvalidAgentRunnerTerminalIntent)
	}
	if _, err := publisher.CostLedger.Settle(tickets.SettlementInput{
		EntryID:             intent.SettlementEntryID,
		OccurredAt:          occurredAt,
		RunID:               intent.RunID,
		RequestID:           intent.RequestID,
		IdempotencyKey:      intent.SettlementIdempotencyKey,
		ReservationHash:     intent.ReservationHash,
		Status:              intent.SettlementStatus,
		ChargedAmountMicros: intent.ChargedAmountMicros,
		ObservedCalls:       intent.ObservedCalls,
		ObservedTokens:      intent.ObservedTokens,
		ProviderEvidence:    intent.ProviderEvidence,
	}); err != nil {
		return fmt.Errorf("%w: %w", ErrAgentRunnerTerminalSettlementBlocked, err)
	}
	return nil
}

func (publisher *AgentRunnerTerminalPublisher) terminalCorruptionError(requestID, reason string) error {
	path, err := publisher.Store.terminalClosurePath(requestID)
	if err != nil {
		return err
	}
	return replayStoreCorruption(path, reason, nil)
}

// validateAgentRunnerTerminalCompletionBinding reuses the shared completion
// binding rules (including the resume session check) so that nothing the
// replay store would reject after the settlement is accepted before it.
func validateAgentRunnerTerminalCompletionBinding(record AgentRunnerRequestRecord, completion AgentRunnerCompletionRecord) error {
	if err := ValidateAgentRunnerCompletionBinding(record, completion); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAgentRunnerTerminalIntent, err)
	}
	return nil
}

// terminalSettlementIntentFields validates the settlement matrix and builds the
// deterministic settlement plan embedded in the terminal intent. The plan is a
// pure function of (request, completion digest, usage observation), so replays
// can prove whether a stored decision matches the caller's claim.
func terminalSettlementIntentFields(record AgentRunnerRequestRecord, completion AgentRunnerCompletionRecord, settlement AgentRunnerTerminalSettlement) (AgentRunnerTerminalIntent, error) {
	switch settlement.Status {
	case "observed", "estimated":
		if settlement.ChargedAmountMicros == nil || *settlement.ChargedAmountMicros < 0 {
			return AgentRunnerTerminalIntent{}, fmt.Errorf("%w: %s settlement requires a non-negative charged amount", ErrInvalidAgentRunnerTerminalIntent, settlement.Status)
		}
	case "unknown":
		if settlement.ChargedAmountMicros != nil {
			return AgentRunnerTerminalIntent{}, fmt.Errorf("%w: unknown settlement cannot carry a charged amount", ErrInvalidAgentRunnerTerminalIntent)
		}
	default:
		return AgentRunnerTerminalIntent{}, fmt.Errorf("%w: unknown settlement status %q", ErrInvalidAgentRunnerTerminalIntent, settlement.Status)
	}
	complete := completion.Completion
	if (complete.Status == "uncertain" || !complete.UsageComplete) && settlement.Status != "unknown" {
		return AgentRunnerTerminalIntent{}, fmt.Errorf("%w: outcome %s with usageComplete=%t permits only an unknown settlement", ErrInvalidAgentRunnerTerminalIntent, complete.Status, complete.UsageComplete)
	}
	body := record.Request
	var calls, tokens *int
	if settlement.Status != "unknown" {
		// The ledger clears observed values for unknown settlements; normalizing
		// here keeps the stored plan identical across replays.
		calls = settlement.ObservedCalls
		tokens = settlement.ObservedTokens
	}
	if calls != nil && *calls < 0 {
		return AgentRunnerTerminalIntent{}, fmt.Errorf("%w: observedCalls must be >= 0", ErrInvalidAgentRunnerTerminalIntent)
	}
	if tokens != nil && *tokens < 0 {
		return AgentRunnerTerminalIntent{}, fmt.Errorf("%w: observedTokens must be >= 0", ErrInvalidAgentRunnerTerminalIntent)
	}
	for _, digest := range settlement.ProviderEvidence {
		if !evidence.ValidHash(digest) {
			return AgentRunnerTerminalIntent{}, fmt.Errorf("%w: provider evidence must contain valid digests", ErrInvalidAgentRunnerTerminalIntent)
		}
	}
	providerEvidence := dedupeTerminalEvidence([]string{completion.RecordHash, record.RecordHash}, settlement.ProviderEvidence)
	key := stableTerminalSettlementKey(body.RequestID, completion.RecordHash)
	return AgentRunnerTerminalIntent{
		SettlementEntryID:        key,
		SettlementIdempotencyKey: key,
		ReservationHash:          body.BudgetHash,
		SettlementStatus:         settlement.Status,
		ChargedAmountMicros:      settlement.ChargedAmountMicros,
		ObservedCalls:            calls,
		ObservedTokens:           tokens,
		ProviderEvidence:         providerEvidence,
	}, nil
}

func buildAgentRunnerTerminalClosure(intent AgentRunnerTerminalIntentRecord) (AgentRunnerTerminalClosureRecord, error) {
	return NewAgentRunnerTerminalClosureRecord(AgentRunnerTerminalClosure{
		RequestID:                intent.Intent.RequestID,
		IntentRecordHash:         intent.RecordHash,
		CompletionHash:           intent.Intent.Completion.RecordHash,
		SettlementEntryID:        intent.Intent.SettlementEntryID,
		SettlementIdempotencyKey: intent.Intent.SettlementIdempotencyKey,
		ClosedAt:                 intent.Intent.DecidedAt,
	})
}

// attachTerminalIntentFacts binds the frozen facts to the settlement plan
// before the intent is published. A completed claim must carry them (the chain
// refuses to route a completion without provable facts) and only a completed
// claim may carry them.
func attachTerminalIntentFacts(plan *AgentRunnerTerminalIntent, completion AgentRunnerCompletionRecord, facts *AgentRunnerFrozenFacts) error {
	completed := completion.Completion.Status == "completed"
	switch {
	case completed && facts == nil:
		return fmt.Errorf("%w: a completed terminal must carry frozen facts", ErrInvalidAgentRunnerTerminalIntent)
	case !completed && facts != nil:
		return fmt.Errorf("%w: only a completed terminal carries frozen facts", ErrInvalidAgentRunnerTerminalIntent)
	case facts == nil:
		return nil
	}
	if err := facts.Validate(); err != nil {
		return err
	}
	plan.Facts = facts
	return nil
}

// dedupeTerminalEvidence composes the settlement evidence deterministically:
// the completion and request digests first, then the caller's observed
// evidence, dropping repeats so the ledger's duplicate-hash rule can never
// reject a plan the publisher already committed to disk.
func dedupeTerminalEvidence(prefix, values []string) []string {
	result := make([]string, 0, len(prefix)+len(values))
	seen := make(map[string]struct{}, len(prefix)+len(values))
	for _, value := range prefix {
		seen[value] = struct{}{}
		result = append(result, value)
	}
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
