package chain

import (
	"errors"
	"fmt"
	"slices"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	// ErrInvalidAgentRunnerTerminal reports a terminal fact that cannot be
	// trusted: a zero value, a mixed create/resume binding, a missing or
	// malformed digest, evidence without the completion digest, or a resume
	// whose session does not continue the prior one.
	ErrInvalidAgentRunnerTerminal = errors.New("invalid AgentRunner terminal")
	// ErrAgentRunnerTerminalConflict reports a terminal fact that cannot be
	// attached to the step's parked dispatch, or a second terminal for the same
	// request that claims a different completion.
	ErrAgentRunnerTerminalConflict = errors.New("AgentRunner terminal conflict")
	// ErrAwaitingAgentRunnerTerminal reports a step parked in TERMINAL_PENDING:
	// external execution was dispatched and the chain awaits the terminal fact.
	// Callers must never dispatch that step again or conclude it from the
	// dispatch result alone.
	ErrAwaitingAgentRunnerTerminal = errors.New("awaiting AgentRunner terminal")
	// ErrAgentRunnerResumeContinuityUnproven reports a resume terminal whose
	// prior completion cannot be proven from the step's event evidence history.
	// The caller must build a new attempt from the durable envelope; the chain
	// never guesses that a session continued.
	ErrAgentRunnerResumeContinuityUnproven = errors.New("AgentRunner resume continuity unproven")
)

// Step-transition reasons owned by AgentRunner dispatch and terminal routing.
// They stay distinct from the operator-interaction reasons so the evidence
// trail always shows which actor produced a waiting transition.
const (
	agentRunnerDispatchedReason        = "agent-dispatched"
	agentRunnerTerminalPassedReason    = "agent-terminal-completed"
	agentRunnerTerminalFailedReason    = "agent-terminal-failed"
	agentRunnerTerminalCancelledReason = "agent-terminal-cancelled"
	agentRunnerTerminalUncertainReason = "agent-terminal-uncertain"
	agentRunnerTerminalWaitingReason   = "agent-terminal-operator-action-required"
	agentRunnerStopUncertainReason     = "stop-uncertain"
	agentRunnerRecoveryUncertainReason = "recovery-uncertain"
	operatorInteractionAnsweredReason  = "operator-interaction-answered"
)

// AgentRunnerTerminalStatus is the chain-owned classification of an external
// terminal fact. It never carries settlement, wire, or provider detail.
type AgentRunnerTerminalStatus string

const (
	AgentRunnerTerminalCompleted              AgentRunnerTerminalStatus = "completed"
	AgentRunnerTerminalFailed                 AgentRunnerTerminalStatus = "failed"
	AgentRunnerTerminalCancelled              AgentRunnerTerminalStatus = "cancelled"
	AgentRunnerTerminalOperatorActionRequired AgentRunnerTerminalStatus = "operator-action-required"
	AgentRunnerTerminalUncertain              AgentRunnerTerminalStatus = "uncertain"
)

// AgentRunnerTerminal is the platform-neutral terminal fact the core routes.
// It is produced by the adapters layer and is the only shape that may move an
// AgentRunner step out of TERMINAL_PENDING. Evidence is deduplicated and order
// preserving, and always contains the request bound to the parked dispatch and
// the completion digest that proves the external execution ended.
type AgentRunnerTerminal struct {
	RequestID           string
	RequestHash         string
	CompletionHash      string
	RunID               string
	TaskID              string
	StepID              string
	Attempt             int
	SessionID           string
	PriorSessionID      string
	PriorCompletionHash string
	Status              AgentRunnerTerminalStatus
	Evidence            []string
}

// Validate fails closed on any terminal fact that must not route. A zero value,
// a create statement carrying prior binding, a resume statement missing it, a
// non-continued resume session, a missing digest, and evidence that cannot
// prove the completion are all rejected before any state write.
func (terminal AgentRunnerTerminal) Validate() error {
	switch terminal.Status {
	case AgentRunnerTerminalCompleted, AgentRunnerTerminalFailed, AgentRunnerTerminalCancelled, AgentRunnerTerminalOperatorActionRequired, AgentRunnerTerminalUncertain:
	default:
		return fmt.Errorf("%w: invalid status %q", ErrInvalidAgentRunnerTerminal, terminal.Status)
	}
	if !evidence.ValidID(terminal.RequestID) || !evidence.ValidID(terminal.RunID) ||
		!evidence.ValidID(terminal.TaskID) || !evidence.ValidID(terminal.StepID) {
		return fmt.Errorf("%w: invalid identity binding", ErrInvalidAgentRunnerTerminal)
	}
	if !evidence.ValidHash(terminal.RequestHash) || !evidence.ValidHash(terminal.CompletionHash) {
		return fmt.Errorf("%w: invalid digest binding", ErrInvalidAgentRunnerTerminal)
	}
	if terminal.Attempt < 1 {
		return fmt.Errorf("%w: invalid attempt", ErrInvalidAgentRunnerTerminal)
	}
	if !evidence.ValidID(terminal.SessionID) {
		return fmt.Errorf("%w: invalid session identity", ErrInvalidAgentRunnerTerminal)
	}
	switch {
	case terminal.PriorSessionID == "" && terminal.PriorCompletionHash == "":
	case terminal.PriorSessionID != "" && terminal.PriorCompletionHash != "":
		if !evidence.ValidID(terminal.PriorSessionID) || !evidence.ValidHash(terminal.PriorCompletionHash) {
			return fmt.Errorf("%w: invalid prior binding", ErrInvalidAgentRunnerTerminal)
		}
		if terminal.SessionID != terminal.PriorSessionID {
			return fmt.Errorf("%w: resume must continue the prior session", ErrInvalidAgentRunnerTerminal)
		}
	default:
		return fmt.Errorf("%w: mixed create and resume binding", ErrInvalidAgentRunnerTerminal)
	}
	if len(terminal.Evidence) == 0 {
		return fmt.Errorf("%w: terminal evidence missing", ErrInvalidAgentRunnerTerminal)
	}
	if len(uniqueHashes(terminal.Evidence)) != len(terminal.Evidence) {
		return fmt.Errorf("%w: terminal evidence must be deduplicated", ErrInvalidAgentRunnerTerminal)
	}
	for _, hash := range terminal.Evidence {
		if !evidence.ValidHash(hash) {
			return fmt.Errorf("%w: invalid evidence digest", ErrInvalidAgentRunnerTerminal)
		}
	}
	if !slices.Contains(terminal.Evidence, terminal.RequestHash) {
		return fmt.Errorf("%w: terminal evidence must carry the request digest", ErrInvalidAgentRunnerTerminal)
	}
	if !slices.Contains(terminal.Evidence, terminal.CompletionHash) {
		return fmt.Errorf("%w: terminal evidence must carry the completion digest", ErrInvalidAgentRunnerTerminal)
	}
	return nil
}

// AgentRunnerTerminalRoute reports the states the core holds after routing one
// terminal fact. Converged marks an idempotent replay that wrote nothing.
type AgentRunnerTerminalRoute struct {
	Status     AgentRunnerTerminalStatus
	StepState  string
	TaskState  string
	ChainState string
	Converged  bool
}

// resume reports whether the terminal was produced by a resume attempt. The
// adapters layer already enforces create ⇒ no prior binding and resume ⇒ both
// prior fields, so a prior completion identifies a resume statement.
func (terminal AgentRunnerTerminal) resume() bool {
	return terminal.PriorCompletionHash != ""
}

// routeEvidence is the ordered, deduplicated evidence a routing transition
// carries: the request bound to the parked dispatch, the completion digest, and
// the adapter proof attached to the terminal.
func (terminal AgentRunnerTerminal) routeEvidence() []string {
	return uniqueHashes(append([]string{terminal.RequestHash, terminal.CompletionHash}, terminal.Evidence...))
}

// agentTerminalRoutingReason reports whether a step reason proves the step left
// TERMINAL_PENDING through core terminal routing.
func agentTerminalRoutingReason(code string) bool {
	switch code {
	case agentRunnerTerminalPassedReason, agentRunnerTerminalFailedReason,
		agentRunnerTerminalCancelledReason, agentRunnerTerminalUncertainReason,
		agentRunnerTerminalWaitingReason:
		return true
	default:
		return false
	}
}

// stepEvidenceHistoryContains reports whether any event of this step already
// carried the digest. It is the only continuity proof the chain accepts for a
// resume: the chain never reads the replay store or the session itself.
func stepEvidenceHistoryContains(events []evidence.StateEvent, entity evidence.Entity, hash string) bool {
	for _, record := range events {
		if record.Event.Entity != entity {
			continue
		}
		if slices.Contains(record.Event.InputEvidence, hash) {
			return true
		}
	}
	return false
}

// latestEntityEvent returns the newest recorded event for one entity.
func latestEntityEvent(events []evidence.StateEvent, entity evidence.Entity) (evidence.Event, bool) {
	for index := len(events) - 1; index >= 0; index-- {
		if events[index].Event.Entity == entity {
			return events[index].Event, true
		}
	}
	return evidence.Event{}, false
}

// latestEntityEventMatches mirrors latestTransitionMatches over an already
// loaded event log: the newest event of the entity must carry the expected
// state, reason, and every listed evidence digest.
func latestEntityEventMatches(events []evidence.StateEvent, entity evidence.Entity, state, reason string, hashes ...string) bool {
	event, found := latestEntityEvent(events, entity)
	if !found {
		return false
	}
	if event.ToState != state || event.Reason.Code != reason {
		return false
	}
	for _, hash := range hashes {
		if !slices.Contains(event.InputEvidence, hash) {
			return false
		}
	}
	return true
}
