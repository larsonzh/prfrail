package adapters

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	agentRunnerTerminalIntentDomain  = "proofrail:agent-runner-terminal-intent:1\n"
	agentRunnerTerminalClosureDomain = "proofrail:agent-runner-terminal-closure:1\n"
)

var (
	// ErrInvalidAgentRunnerTerminalIntent reports an invalid terminal-intent record.
	ErrInvalidAgentRunnerTerminalIntent = errors.New("invalid AgentRunner terminal intent")
	// ErrInvalidAgentRunnerTerminalClosure reports an invalid terminal-closure record.
	ErrInvalidAgentRunnerTerminalClosure = errors.New("invalid AgentRunner terminal closure")
	// ErrAgentRunnerTerminalIntentConflict reports a terminal-intent slot already
	// owned by a different launch/outcome claim.
	ErrAgentRunnerTerminalIntentConflict = errors.New("AgentRunner terminal intent conflict")
	// ErrAgentRunnerTerminalClosureConflict reports a closure slot already owned by
	// a different closure record.
	ErrAgentRunnerTerminalClosureConflict = errors.New("AgentRunner terminal closure conflict")
	// ErrAgentRunnerTerminalOrphan reports a completion published without any
	// terminal-intent: no settlement decision record can be proven, so the state
	// must fail closed and must never be adopted by writing an intent back.
	ErrAgentRunnerTerminalOrphan = errors.New("AgentRunner terminal orphan completion")
	// ErrAgentRunnerTerminalSettlementBlocked reports a settlement that cannot be
	// booked (foreign settlement, conflict, missing reservation). The underlying
	// cost error is preserved for errors.Is.
	ErrAgentRunnerTerminalSettlementBlocked = errors.New("AgentRunner terminal settlement blocked")
)

// AgentRunnerFrozenFacts are the five digests an adapter froze after external
// execution ended. They are store-local facts, never decisions: the chain runs
// its own postflight and reconciles them before a task may be reviewed.
type AgentRunnerFrozenFacts struct {
	ManifestHash            string `json:"manifestHash"`
	DiffHash                string `json:"diffHash"`
	LogHash                 string `json:"logHash"`
	UsageHash               string `json:"usageHash"`
	ProcessStopEvidenceHash string `json:"processStopEvidenceHash"`
}

func (facts AgentRunnerFrozenFacts) hashes() []string {
	return []string{facts.ManifestHash, facts.DiffHash, facts.LogHash, facts.UsageHash, facts.ProcessStopEvidenceHash}
}

// Validate fails closed on missing digests and on digest collisions: the chain
// rebuilds the facts positionally from the routing evidence, so repeats would
// make the rebuild ambiguous.
func (facts AgentRunnerFrozenFacts) Validate() error {
	hashes := facts.hashes()
	for index, hash := range hashes {
		if !evidence.ValidHash(hash) {
			return fmt.Errorf("%w: frozen fact %d is not a valid hash", ErrInvalidAgentRunnerTerminalIntent, index)
		}
	}
	for index, hash := range hashes {
		if slices.Contains(hashes[:index], hash) {
			return fmt.Errorf("%w: frozen fact %d repeats an earlier digest", ErrInvalidAgentRunnerTerminalIntent, index)
		}
	}
	return nil
}

// AgentRunnerTerminalIntent is the pre-settlement terminal decision record. It
// embeds the full completion record and the settlement plan so crash recovery
// can re-drive the whole chain from store-local facts alone. It is store-local
// metadata: never on the wire and never part of any R/C recordHash.
type AgentRunnerTerminalIntent struct {
	Kind                     string                      `json:"kind"`
	RequestID                string                      `json:"requestId"`
	RequestHash              string                      `json:"requestHash"`
	RunID                    string                      `json:"runId"`
	TaskID                   string                      `json:"taskId"`
	StepID                   string                      `json:"stepId"`
	Attempt                  int                         `json:"attempt"`
	AdapterID                string                      `json:"adapterId"`
	Completion               AgentRunnerCompletionRecord `json:"completion"`
	Facts                    *AgentRunnerFrozenFacts     `json:"frozenFacts,omitempty"`
	SettlementEntryID        string                      `json:"settlementEntryId"`
	SettlementIdempotencyKey string                      `json:"settlementIdempotencyKey"`
	ReservationHash          string                      `json:"reservationHash"`
	SettlementStatus         string                      `json:"settlementStatus"`
	ChargedAmountMicros      *int64                      `json:"chargedAmountMicros"`
	ObservedCalls            *int                        `json:"observedCalls"`
	ObservedTokens           *int                        `json:"observedTokens"`
	ProviderEvidence         []string                    `json:"providerEvidence"`
	DecidedAt                string                      `json:"decidedAt"`
}

// AgentRunnerTerminalIntentRecord wraps the terminal-intent body with its local digest.
type AgentRunnerTerminalIntentRecord struct {
	SchemaVersion string                    `json:"schemaVersion"`
	Intent        AgentRunnerTerminalIntent `json:"intent"`
	RecordHash    string                    `json:"recordHash"`
}

// AgentRunnerTerminalClosure marks a completed terminal chain: the chain is
// complete if and only if a closure exists and binds the intent, the completion
// and the settlement keys. Store-local metadata, never on the wire.
type AgentRunnerTerminalClosure struct {
	Kind                     string `json:"kind"`
	RequestID                string `json:"requestId"`
	IntentRecordHash         string `json:"intentRecordHash"`
	CompletionHash           string `json:"completionHash"`
	SettlementEntryID        string `json:"settlementEntryId"`
	SettlementIdempotencyKey string `json:"settlementIdempotencyKey"`
	ClosedAt                 string `json:"closedAt"`
}

// AgentRunnerTerminalClosureRecord wraps the closure body with its local digest.
type AgentRunnerTerminalClosureRecord struct {
	SchemaVersion string                     `json:"schemaVersion"`
	Closure       AgentRunnerTerminalClosure `json:"closure"`
	RecordHash    string                     `json:"recordHash"`
}

// NewAgentRunnerTerminalIntentRecord builds a validated terminal-intent record.
func NewAgentRunnerTerminalIntentRecord(intent AgentRunnerTerminalIntent) (AgentRunnerTerminalIntentRecord, error) {
	if intent.Kind == "" {
		intent.Kind = "agent-runner-terminal-intent"
	}
	if err := validateAgentRunnerTerminalIntent(intent); err != nil {
		return AgentRunnerTerminalIntentRecord{}, err
	}
	hash, err := digestMessage(agentRunnerTerminalIntentDomain, intent)
	if err != nil {
		return AgentRunnerTerminalIntentRecord{}, err
	}
	return AgentRunnerTerminalIntentRecord{SchemaVersion: evidence.SchemaVersion, Intent: intent, RecordHash: hash}, nil
}

// DecodeAgentRunnerTerminalIntentRecord strictly decodes and validates a record.
func DecodeAgentRunnerTerminalIntentRecord(input []byte) (AgentRunnerTerminalIntentRecord, error) {
	var record AgentRunnerTerminalIntentRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AgentRunnerTerminalIntentRecord{}, err
	}
	if err := ValidateAgentRunnerTerminalIntentRecord(record); err != nil {
		return AgentRunnerTerminalIntentRecord{}, err
	}
	return record, nil
}

// ValidateAgentRunnerTerminalIntentRecord validates schema, hashes and body.
func ValidateAgentRunnerTerminalIntentRecord(record AgentRunnerTerminalIntentRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAgentRunnerTerminalIntent)
	}
	if err := validateAgentRunnerTerminalIntent(record.Intent); err != nil {
		return err
	}
	expected, err := digestMessage(agentRunnerTerminalIntentDomain, record.Intent)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAgentRunnerTerminalIntent)
	}
	return nil
}

func validateAgentRunnerTerminalIntent(intent AgentRunnerTerminalIntent) error {
	if intent.Kind != "agent-runner-terminal-intent" {
		return fmt.Errorf("%w: invalid kind", ErrInvalidAgentRunnerTerminalIntent)
	}
	if !evidence.ValidID(intent.RequestID) || !evidence.ValidID(intent.RunID) ||
		!evidence.ValidID(intent.TaskID) || !evidence.ValidID(intent.StepID) ||
		!evidence.ValidID(intent.AdapterID) {
		return fmt.Errorf("%w: invalid identity binding", ErrInvalidAgentRunnerTerminalIntent)
	}
	if !evidence.ValidHash(intent.RequestHash) || !evidence.ValidHash(intent.ReservationHash) {
		return fmt.Errorf("%w: invalid digest binding", ErrInvalidAgentRunnerTerminalIntent)
	}
	if intent.Attempt < 1 {
		return fmt.Errorf("%w: invalid attempt", ErrInvalidAgentRunnerTerminalIntent)
	}
	if err := ValidateAgentRunnerCompletionRecord(intent.Completion); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidAgentRunnerTerminalIntent, err)
	}
	completion := intent.Completion.Completion
	if completion.RequestID != intent.RequestID || completion.RequestHash != intent.RequestHash ||
		completion.RunID != intent.RunID || completion.TaskID != intent.TaskID ||
		completion.StepID != intent.StepID || completion.Attempt != intent.Attempt ||
		completion.AdapterID != intent.AdapterID {
		return fmt.Errorf("%w: completion binding mismatch", ErrInvalidAgentRunnerTerminalIntent)
	}
	if !evidence.ValidID(intent.SettlementEntryID) || !evidence.ValidID(intent.SettlementIdempotencyKey) {
		return fmt.Errorf("%w: invalid settlement identity", ErrInvalidAgentRunnerTerminalIntent)
	}
	switch intent.SettlementStatus {
	case "observed", "estimated":
		if intent.ChargedAmountMicros == nil || *intent.ChargedAmountMicros < 0 {
			return fmt.Errorf("%w: %s settlement requires a non-negative charged amount", ErrInvalidAgentRunnerTerminalIntent, intent.SettlementStatus)
		}
	case "unknown":
		if intent.ChargedAmountMicros != nil {
			return fmt.Errorf("%w: unknown settlement cannot carry a charged amount", ErrInvalidAgentRunnerTerminalIntent)
		}
	default:
		return fmt.Errorf("%w: unknown settlement status %q", ErrInvalidAgentRunnerTerminalIntent, intent.SettlementStatus)
	}
	if len(intent.ProviderEvidence) == 0 {
		return fmt.Errorf("%w: settlement evidence missing", ErrInvalidAgentRunnerTerminalIntent)
	}
	if !slices.Contains(intent.ProviderEvidence, intent.Completion.RecordHash) {
		return fmt.Errorf("%w: settlement evidence must carry the completion digest", ErrInvalidAgentRunnerTerminalIntent)
	}
	if !slices.Contains(intent.ProviderEvidence, intent.RequestHash) {
		return fmt.Errorf("%w: settlement evidence must carry the request digest", ErrInvalidAgentRunnerTerminalIntent)
	}
	if _, err := time.Parse(time.RFC3339, intent.DecidedAt); err != nil {
		return fmt.Errorf("%w: invalid decidedAt", ErrInvalidAgentRunnerTerminalIntent)
	}
	return validateAgentRunnerTerminalIntentFacts(intent)
}

// validateAgentRunnerTerminalIntentFacts enforces the frozen-facts binding: a
// completed terminal must carry the five facts, and no other status may carry
// them. Without facts a completion can never be audited by postflight, so the
// store refuses to hold the decision at all.
func validateAgentRunnerTerminalIntentFacts(intent AgentRunnerTerminalIntent) error {
	completed := intent.Completion.Completion.Status == "completed"
	switch {
	case completed && intent.Facts == nil:
		return fmt.Errorf("%w: a completed terminal intent must carry frozen facts", ErrInvalidAgentRunnerTerminalIntent)
	case !completed && intent.Facts != nil:
		return fmt.Errorf("%w: only a completed terminal intent carries frozen facts", ErrInvalidAgentRunnerTerminalIntent)
	case intent.Facts == nil:
		return nil
	default:
		return intent.Facts.Validate()
	}
}

// NewAgentRunnerTerminalClosureRecord builds a validated closure record.
func NewAgentRunnerTerminalClosureRecord(closure AgentRunnerTerminalClosure) (AgentRunnerTerminalClosureRecord, error) {
	if closure.Kind == "" {
		closure.Kind = "agent-runner-terminal-closure"
	}
	if err := validateAgentRunnerTerminalClosure(closure); err != nil {
		return AgentRunnerTerminalClosureRecord{}, err
	}
	hash, err := digestMessage(agentRunnerTerminalClosureDomain, closure)
	if err != nil {
		return AgentRunnerTerminalClosureRecord{}, err
	}
	return AgentRunnerTerminalClosureRecord{SchemaVersion: evidence.SchemaVersion, Closure: closure, RecordHash: hash}, nil
}

// DecodeAgentRunnerTerminalClosureRecord strictly decodes and validates a record.
func DecodeAgentRunnerTerminalClosureRecord(input []byte) (AgentRunnerTerminalClosureRecord, error) {
	var record AgentRunnerTerminalClosureRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return AgentRunnerTerminalClosureRecord{}, err
	}
	if err := ValidateAgentRunnerTerminalClosureRecord(record); err != nil {
		return AgentRunnerTerminalClosureRecord{}, err
	}
	return record, nil
}

// ValidateAgentRunnerTerminalClosureRecord validates schema, hashes and body.
func ValidateAgentRunnerTerminalClosureRecord(record AgentRunnerTerminalClosureRecord) error {
	if record.SchemaVersion != evidence.SchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidAgentRunnerTerminalClosure)
	}
	if err := validateAgentRunnerTerminalClosure(record.Closure); err != nil {
		return err
	}
	expected, err := digestMessage(agentRunnerTerminalClosureDomain, record.Closure)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidAgentRunnerTerminalClosure)
	}
	return nil
}

func validateAgentRunnerTerminalClosure(closure AgentRunnerTerminalClosure) error {
	if closure.Kind != "agent-runner-terminal-closure" {
		return fmt.Errorf("%w: invalid kind", ErrInvalidAgentRunnerTerminalClosure)
	}
	if !evidence.ValidID(closure.RequestID) {
		return fmt.Errorf("%w: invalid identity binding", ErrInvalidAgentRunnerTerminalClosure)
	}
	if !evidence.ValidHash(closure.IntentRecordHash) || !evidence.ValidHash(closure.CompletionHash) {
		return fmt.Errorf("%w: invalid digest binding", ErrInvalidAgentRunnerTerminalClosure)
	}
	if !evidence.ValidID(closure.SettlementEntryID) || !evidence.ValidID(closure.SettlementIdempotencyKey) {
		return fmt.Errorf("%w: invalid settlement identity", ErrInvalidAgentRunnerTerminalClosure)
	}
	if _, err := time.Parse(time.RFC3339, closure.ClosedAt); err != nil {
		return fmt.Errorf("%w: invalid closedAt", ErrInvalidAgentRunnerTerminalClosure)
	}
	return nil
}

func stableTerminalSettlementKey(requestID, completionHash string) string {
	hash := evidence.Digest("proofrail:agent-runner-settlement-key:1\n", []byte(requestID+"\n"+completionHash))
	hex := strings.TrimPrefix(hash, "sha256:")
	if len(hex) > 20 {
		hex = hex[:20]
	}
	return "settle-" + hex
}
