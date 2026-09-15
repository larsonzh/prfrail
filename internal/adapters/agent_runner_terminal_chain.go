package adapters

import (
	"errors"
	"fmt"

	"github.com/larsonzh/prfrail/internal/chain"
)

// ErrInvalidAgentRunnerTerminalTranslation reports a terminal chain that cannot
// be translated into a chain terminal fact: a zero value, an unbound record, a
// status the chain does not know, or a completion that does not belong to the
// request.
var ErrInvalidAgentRunnerTerminalTranslation = errors.New("invalid AgentRunner terminal translation")

// ToChainAgentRunnerTerminal translates a published terminal chain into the
// chain-owned terminal fact. The translation is one way: the chain never reads
// the replay store, the session, or the process tree.
//
// Every input must already be valid and must bind the others: the closure binds
// the intent record and the completion digest, the completion binds the request
// (including the resume session rule), and only a status the chain models is
// accepted. A terminal fact that cannot prove the request and the completion is
// refused instead of being translated into something caller-friendly.
func ToChainAgentRunnerTerminal(request AgentRunnerRequestRecord, intent AgentRunnerTerminalIntentRecord, closure AgentRunnerTerminalClosureRecord) (chain.AgentRunnerTerminal, error) {
	if err := ValidateAgentRunnerRequestRecord(request); err != nil {
		return chain.AgentRunnerTerminal{}, translateAgentRunnerTerminalError(err)
	}
	if err := ValidateAgentRunnerTerminalIntentRecord(intent); err != nil {
		return chain.AgentRunnerTerminal{}, translateAgentRunnerTerminalError(err)
	}
	if err := ValidateAgentRunnerTerminalClosureRecord(closure); err != nil {
		return chain.AgentRunnerTerminal{}, translateAgentRunnerTerminalError(err)
	}
	body := request.Request
	intentBody := intent.Intent
	closureBody := closure.Closure
	if closureBody.RequestID != intentBody.RequestID ||
		closureBody.IntentRecordHash != intent.RecordHash ||
		closureBody.CompletionHash != intentBody.Completion.RecordHash ||
		closureBody.SettlementEntryID != intentBody.SettlementEntryID ||
		closureBody.SettlementIdempotencyKey != intentBody.SettlementIdempotencyKey {
		return chain.AgentRunnerTerminal{}, fmt.Errorf("%w: closure does not bind the intent, completion and settlement slot", ErrInvalidAgentRunnerTerminalTranslation)
	}
	if err := ValidateAgentRunnerCompletionBinding(request, intentBody.Completion); err != nil {
		return chain.AgentRunnerTerminal{}, translateAgentRunnerTerminalError(err)
	}
	status, err := translateAgentRunnerTerminalStatus(intentBody.Completion.Completion.Status)
	if err != nil {
		return chain.AgentRunnerTerminal{}, err
	}
	var facts *chain.AgentRunnerFrozenFacts
	if intentBody.Facts != nil {
		projected := intentBody.Facts.toChain()
		facts = &projected
	}
	completion := intentBody.Completion.Completion
	terminal := chain.AgentRunnerTerminal{
		RequestID:           body.RequestID,
		RequestHash:         request.RecordHash,
		CompletionHash:      intentBody.Completion.RecordHash,
		RunID:               body.RunID,
		TaskID:              body.TaskID,
		StepID:              body.StepID,
		Attempt:             body.Attempt,
		SessionID:           completion.SessionID,
		PriorSessionID:      body.PriorSessionID,
		PriorCompletionHash: body.PriorCompletionHash,
		Status:              status,
		Facts:               facts,
		Evidence: dedupeTerminalEvidence(
			[]string{request.RecordHash, intentBody.Completion.RecordHash, intent.RecordHash, closure.RecordHash},
			append(append(append([]string{}, completion.Evidence...), completion.ErrorEvidence...), intentBody.ProviderEvidence...),
		),
	}
	if err := terminal.Validate(); err != nil {
		return chain.AgentRunnerTerminal{}, translateAgentRunnerTerminalError(err)
	}
	return terminal, nil
}

// LoadChainAgentRunnerTerminal reads the durable terminal chain of one request
// and translates it. It never writes, settles, or concludes anything: an
// unsettled terminal (no intent), an unclosed terminal (no closure), and any
// corrupted record are refused before the chain is asked to route.
func LoadChainAgentRunnerTerminal(store *AgentRunnerReplayStore, request AgentRunnerRequestRecord) (chain.AgentRunnerTerminal, error) {
	if err := ValidateAgentRunnerRequestRecord(request); err != nil {
		return chain.AgentRunnerTerminal{}, translateAgentRunnerTerminalError(err)
	}
	if store == nil {
		return chain.AgentRunnerTerminal{}, fmt.Errorf("%w: replay store missing", ErrInvalidAgentRunnerTerminalTranslation)
	}
	requestID := request.Request.RequestID
	intent, found, err := store.TerminalIntent(requestID)
	if err != nil {
		return chain.AgentRunnerTerminal{}, err
	}
	if !found {
		return chain.AgentRunnerTerminal{}, fmt.Errorf("%w: terminal intent missing for requestId %s", ErrInvalidAgentRunnerTerminalTranslation, requestID)
	}
	closure, found, err := store.TerminalClosure(requestID)
	if err != nil {
		return chain.AgentRunnerTerminal{}, err
	}
	if !found {
		return chain.AgentRunnerTerminal{}, fmt.Errorf("%w: terminal closure missing for requestId %s", ErrInvalidAgentRunnerTerminalTranslation, requestID)
	}
	return ToChainAgentRunnerTerminal(request, intent, closure)
}

// toChain projects the store-local frozen facts onto the chain-owned DTO. The
// projection is one way and carries digests only: no decision, no settlement
// status, and no wire record ever crosses this boundary.
func (facts AgentRunnerFrozenFacts) toChain() chain.AgentRunnerFrozenFacts {
	return chain.AgentRunnerFrozenFacts{
		ManifestHash:            facts.ManifestHash,
		DiffHash:                facts.DiffHash,
		LogHash:                 facts.LogHash,
		UsageHash:               facts.UsageHash,
		ProcessStopEvidenceHash: facts.ProcessStopEvidenceHash,
	}
}

func translateAgentRunnerTerminalStatus(status string) (chain.AgentRunnerTerminalStatus, error) {
	switch status {
	case "completed":
		return chain.AgentRunnerTerminalCompleted, nil
	case "failed":
		return chain.AgentRunnerTerminalFailed, nil
	case "cancelled":
		return chain.AgentRunnerTerminalCancelled, nil
	case "operator-action-required":
		return chain.AgentRunnerTerminalOperatorActionRequired, nil
	case "uncertain":
		return chain.AgentRunnerTerminalUncertain, nil
	default:
		return "", fmt.Errorf("%w: unknown completion status %q", ErrInvalidAgentRunnerTerminalTranslation, status)
	}
}

// translateAgentRunnerTerminalError keeps both the translation sentinel and the
// underlying record error identifiable for callers.
func translateAgentRunnerTerminalError(err error) error {
	return fmt.Errorf("%w: %w", ErrInvalidAgentRunnerTerminalTranslation, err)
}
