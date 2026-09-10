package chain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	operatorInteractionDomain        = "proofrail:operator-interaction:1\n"
	operatorInteractionSchemaVersion = "1.0.0"
	operatorInteractionTimeLayout    = "2006-01-02T15:04:05.000Z"
)

var (
	ErrInvalidOperatorInteraction   = errors.New("invalid operator interaction")
	ErrDuplicateInteractionResponse = errors.New("duplicate operator interaction response")
)

type OperatorInteractionRequest struct {
	RecordID         string         `json:"recordId"`
	Kind             string         `json:"kind"`
	InteractionID    string         `json:"interactionId"`
	RequestedAt      string         `json:"requestedAt"`
	RequestedBy      evidence.Actor `json:"requestedBy"`
	Operator         evidence.Actor `json:"operator"`
	RunID            string         `json:"runId"`
	TaskID           string         `json:"taskId"`
	StepID           string         `json:"stepId"`
	Attempt          int            `json:"attempt"`
	WorkspaceHash    string         `json:"workspaceHash"`
	ConversationID   string         `json:"conversationId"`
	RequestID        string         `json:"requestId"`
	ContextHash      string         `json:"contextHash"`
	Question         string         `json:"question"`
	Reason           string         `json:"reason"`
	AllowedResponses []string       `json:"allowedResponses"`
	Risk             string         `json:"risk"`
	RequiredAction   string         `json:"requiredAction"`
	Evidence         []string       `json:"evidence"`
}

type OperatorInteractionResponse struct {
	RecordID       string         `json:"recordId"`
	Kind           string         `json:"kind"`
	InteractionID  string         `json:"interactionId"`
	RespondedAt    string         `json:"respondedAt"`
	RespondedBy    evidence.Actor `json:"respondedBy"`
	RunID          string         `json:"runId"`
	TaskID         string         `json:"taskId"`
	StepID         string         `json:"stepId"`
	Attempt        int            `json:"attempt"`
	WorkspaceHash  string         `json:"workspaceHash"`
	ConversationID string         `json:"conversationId"`
	RequestID      string         `json:"requestId"`
	ContextHash    string         `json:"contextHash"`
	RequestHash    string         `json:"requestHash"`
	Selection      string         `json:"selection"`
	Evidence       []string       `json:"evidence"`
}

type OperatorInteractionRecord struct {
	SchemaVersion string                       `json:"schemaVersion"`
	Request       *OperatorInteractionRequest  `json:"-"`
	Response      *OperatorInteractionResponse `json:"-"`
	RecordHash    string                       `json:"recordHash"`
}

func (record OperatorInteractionRecord) MarshalJSON() ([]byte, error) {
	var body any
	switch {
	case record.Request != nil && record.Response == nil:
		body = record.Request
	case record.Response != nil && record.Request == nil:
		body = record.Response
	default:
		return nil, fmt.Errorf("%w: exactly one request or response required", ErrInvalidOperatorInteraction)
	}
	return json.Marshal(struct {
		SchemaVersion string `json:"schemaVersion"`
		Interaction   any    `json:"interaction"`
		RecordHash    string `json:"recordHash"`
	}{record.SchemaVersion, body, record.RecordHash})
}

func (record *OperatorInteractionRecord) UnmarshalJSON(input []byte) error {
	var wire struct {
		SchemaVersion string          `json:"schemaVersion"`
		Interaction   json.RawMessage `json:"interaction"`
		RecordHash    string          `json:"recordHash"`
	}
	if err := evidence.DecodeStrictJSON(input, &wire); err != nil {
		return err
	}
	var kind struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(wire.Interaction, &kind); err != nil {
		return err
	}
	record.SchemaVersion = wire.SchemaVersion
	record.RecordHash = wire.RecordHash
	switch kind.Kind {
	case "request":
		var request OperatorInteractionRequest
		if err := evidence.DecodeStrictJSON(wire.Interaction, &request); err != nil {
			return err
		}
		record.Request = &request
	case "response":
		var response OperatorInteractionResponse
		if err := evidence.DecodeStrictJSON(wire.Interaction, &response); err != nil {
			return err
		}
		record.Response = &response
	default:
		return fmt.Errorf("%w: unknown kind %q", ErrInvalidOperatorInteraction, kind.Kind)
	}
	return nil
}

func NewOperatorInteractionRequestRecord(request OperatorInteractionRequest, clock Clock, ids IDSource) (OperatorInteractionRecord, error) {
	if clock == nil {
		clock = defaultClock
	}
	if ids == nil {
		ids = randomID
	}
	if request.RecordID == "" {
		request.RecordID = ids()
	}
	if request.Kind == "" {
		request.Kind = "request"
	}
	if request.RequestedAt == "" {
		request.RequestedAt = clock().UTC().Format(operatorInteractionTimeLayout)
	}
	if err := validateOperatorInteractionRequest(request); err != nil {
		return OperatorInteractionRecord{}, err
	}
	hash, err := digestRecord(operatorInteractionDomain, request)
	if err != nil {
		return OperatorInteractionRecord{}, err
	}
	return OperatorInteractionRecord{SchemaVersion: operatorInteractionSchemaVersion, Request: &request, RecordHash: hash}, nil
}

func NewOperatorInteractionResponseRecord(requestRecord OperatorInteractionRecord, response OperatorInteractionResponse, clock Clock, ids IDSource) (OperatorInteractionRecord, error) {
	if err := ValidateOperatorInteractionRecord(requestRecord); err != nil || requestRecord.Request == nil {
		return OperatorInteractionRecord{}, fmt.Errorf("%w: invalid request record", ErrInvalidOperatorInteraction)
	}
	if clock == nil {
		clock = defaultClock
	}
	if ids == nil {
		ids = randomID
	}
	if response.RecordID == "" {
		response.RecordID = ids()
	}
	if response.RequestID == "" {
		response.RequestID = ids()
	}
	if response.Kind == "" {
		response.Kind = "response"
	}
	if response.RespondedAt == "" {
		response.RespondedAt = clock().UTC().Format(operatorInteractionTimeLayout)
	}
	if err := validateOperatorInteractionResponse(*requestRecord.Request, requestRecord.RecordHash, response); err != nil {
		return OperatorInteractionRecord{}, err
	}
	hash, err := digestRecord(operatorInteractionDomain, response)
	if err != nil {
		return OperatorInteractionRecord{}, err
	}
	return OperatorInteractionRecord{SchemaVersion: operatorInteractionSchemaVersion, Response: &response, RecordHash: hash}, nil
}

func DecodeOperatorInteractionRecord(input []byte) (OperatorInteractionRecord, error) {
	var record OperatorInteractionRecord
	if err := json.Unmarshal(input, &record); err != nil {
		return OperatorInteractionRecord{}, err
	}
	if err := ValidateOperatorInteractionRecord(record); err != nil {
		return OperatorInteractionRecord{}, err
	}
	return record, nil
}

func ValidateOperatorInteractionRecord(record OperatorInteractionRecord) error {
	if record.SchemaVersion != operatorInteractionSchemaVersion || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidOperatorInteraction)
	}
	var body any
	switch {
	case record.Request != nil && record.Response == nil:
		if err := validateOperatorInteractionRequest(*record.Request); err != nil {
			return err
		}
		body = *record.Request
	case record.Response != nil && record.Request == nil:
		if err := validateStandaloneOperatorInteractionResponse(*record.Response); err != nil {
			return err
		}
		body = *record.Response
	default:
		return fmt.Errorf("%w: exactly one request or response required", ErrInvalidOperatorInteraction)
	}
	expected, err := digestRecord(operatorInteractionDomain, body)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidOperatorInteraction)
	}
	return nil
}

type OperatorInteractionLedger struct {
	Records []OperatorInteractionRecord
}

type OperatorInteractionBinding struct {
	RunID          string
	TaskID         string
	StepID         string
	Attempt        int
	WorkspaceHash  string
	ConversationID string
	ContextHash    string
}

type OperatorResumeCommand struct {
	InteractionID  string `json:"interactionId"`
	ConversationID string `json:"conversationId"`
	RequestID      string `json:"requestId"`
	Selection      string `json:"selection"`
	ResponseHash   string `json:"responseHash"`
}

func (ledger *OperatorInteractionLedger) Append(record OperatorInteractionRecord) error {
	if err := ValidateOperatorInteractionRecord(record); err != nil {
		return err
	}
	requests := make(map[string]OperatorInteractionRecord)
	responses := make(map[string]struct{})
	for _, existing := range ledger.Records {
		if existing.Request != nil {
			requests[existing.Request.InteractionID] = existing
		} else if existing.Response != nil {
			responses[existing.Response.InteractionID] = struct{}{}
		}
	}
	if record.Request != nil {
		if _, exists := requests[record.Request.InteractionID]; exists {
			return fmt.Errorf("%w: duplicate interaction request", ErrInvalidOperatorInteraction)
		}
	} else {
		interactionID := record.Response.InteractionID
		if _, exists := responses[interactionID]; exists {
			return ErrDuplicateInteractionResponse
		}
		request, exists := requests[interactionID]
		if !exists {
			return fmt.Errorf("%w: response without request", ErrInvalidOperatorInteraction)
		}
		if err := validateOperatorInteractionResponse(*request.Request, request.RecordHash, *record.Response); err != nil {
			return err
		}
	}
	ledger.Records = append(ledger.Records, record)
	return nil
}

func (ledger OperatorInteractionLedger) Pending() []OperatorInteractionRecord {
	answered := make(map[string]struct{})
	for _, record := range ledger.Records {
		if record.Response != nil {
			answered[record.Response.InteractionID] = struct{}{}
		}
	}
	pending := make([]OperatorInteractionRecord, 0)
	for _, record := range ledger.Records {
		if record.Request != nil {
			if _, exists := answered[record.Request.InteractionID]; !exists {
				pending = append(pending, record)
			}
		}
	}
	return pending
}

func PrepareOperatorInteractionResume(requestRecord, responseRecord OperatorInteractionRecord, current OperatorInteractionBinding, taskState, stepState string) (OperatorResumeCommand, error) {
	command, err := prepareOperatorInteractionResume(requestRecord, responseRecord, current)
	if err != nil {
		return OperatorResumeCommand{}, err
	}
	if taskState != "WAITING_FOR_OPERATOR" || stepState != "WAITING_FOR_OPERATOR" {
		return OperatorResumeCommand{}, fmt.Errorf("%w: interaction resume requires WAITING_FOR_OPERATOR", ErrInvalidState)
	}
	return command, nil
}

func prepareOperatorInteractionResume(requestRecord, responseRecord OperatorInteractionRecord, current OperatorInteractionBinding) (OperatorResumeCommand, error) {
	if err := ValidateOperatorInteractionRecord(requestRecord); err != nil || requestRecord.Request == nil {
		return OperatorResumeCommand{}, fmt.Errorf("%w: invalid request record", ErrInvalidOperatorInteraction)
	}
	if err := ValidateOperatorInteractionRecord(responseRecord); err != nil || responseRecord.Response == nil {
		return OperatorResumeCommand{}, fmt.Errorf("%w: invalid response record", ErrInvalidOperatorInteraction)
	}
	request := *requestRecord.Request
	response := *responseRecord.Response
	if err := validateOperatorInteractionResponse(request, requestRecord.RecordHash, response); err != nil {
		return OperatorResumeCommand{}, err
	}
	if current.RunID != request.RunID || current.TaskID != request.TaskID || current.StepID != request.StepID || current.Attempt != request.Attempt || current.WorkspaceHash != request.WorkspaceHash || current.ConversationID != request.ConversationID || current.ContextHash != request.ContextHash {
		return OperatorResumeCommand{}, fmt.Errorf("%w: current interaction binding changed", ErrInvalidOperatorInteraction)
	}
	return OperatorResumeCommand{
		InteractionID: request.InteractionID, ConversationID: request.ConversationID,
		RequestID: response.RequestID, Selection: response.Selection, ResponseHash: responseRecord.RecordHash,
	}, nil
}

func (engine *Engine) OpenOperatorInteraction(ctx context.Context, requestRecord OperatorInteractionRecord) error {
	if err := ValidateOperatorInteractionRecord(requestRecord); err != nil || requestRecord.Request == nil {
		return fmt.Errorf("%w: invalid request record", ErrInvalidOperatorInteraction)
	}
	request := *requestRecord.Request
	if request.RunID != engine.options.RunID || request.Attempt != 1 {
		return fmt.Errorf("%w: request does not belong to engine run", ErrInvalidOperatorInteraction)
	}
	task := taskEntity(request.RunID, request.TaskID)
	step := stepEntity(request.RunID, request.TaskID, request.StepID)
	taskState, stepState := engine.state.current(task), engine.state.current(step)
	initial := taskState == "STEPS_RUNNING" && stepState == "RUNNING"
	partial := taskState == "WAITING_FOR_OPERATOR" && (stepState == "RUNNING" || stepState == "WAITING_FOR_OPERATOR")
	if !initial && !partial {
		return fmt.Errorf("%w: interaction open requires active task and step", ErrInvalidState)
	}
	if engine.state.projection.ChainState == "RUNNING" {
		if err := engine.transition(ctx, chainEntity(request.RunID), "PAUSED", []string{requestRecord.RecordHash}, "operator-interaction-required"); err != nil {
			return err
		}
	} else if engine.state.projection.ChainState != "PAUSED" {
		return fmt.Errorf("%w: interaction open requires RUNNING or PAUSED chain", ErrInvalidState)
	} else if matches, err := engine.latestTransitionMatches(ctx, chainEntity(request.RunID), "PAUSED", "operator-interaction-required", requestRecord.RecordHash); err != nil {
		return err
	} else if !matches {
		return fmt.Errorf("%w: paused chain belongs to another control action", ErrInvalidOperatorInteraction)
	}
	if taskState == "WAITING_FOR_OPERATOR" {
		if matches, err := engine.latestTransitionMatches(ctx, task, taskState, "operator-interaction-required", requestRecord.RecordHash); err != nil {
			return err
		} else if !matches {
			return fmt.Errorf("%w: waiting task belongs to another interaction", ErrInvalidOperatorInteraction)
		}
	}
	if stepState == "WAITING_FOR_OPERATOR" {
		if matches, err := engine.latestTransitionMatches(ctx, step, stepState, "operator-interaction-required", requestRecord.RecordHash); err != nil {
			return err
		} else if !matches {
			return fmt.Errorf("%w: waiting step belongs to another interaction", ErrInvalidOperatorInteraction)
		}
	}
	if taskState == "STEPS_RUNNING" {
		if err := engine.transition(ctx, task, "WAITING_FOR_OPERATOR", []string{requestRecord.RecordHash}, "operator-interaction-required"); err != nil {
			return err
		}
	}
	if stepState == "RUNNING" {
		if err := engine.transition(ctx, step, "WAITING_FOR_OPERATOR", []string{requestRecord.RecordHash}, "operator-interaction-required"); err != nil {
			return err
		}
	}
	return nil
}

func (engine *Engine) ResumeOperatorInteraction(ctx context.Context, requestRecord, responseRecord OperatorInteractionRecord, current OperatorInteractionBinding) (OperatorResumeCommand, error) {
	command, err := prepareOperatorInteractionResume(requestRecord, responseRecord, current)
	if err != nil {
		return OperatorResumeCommand{}, err
	}
	request := *requestRecord.Request
	if request.RunID != engine.options.RunID || request.Attempt != 1 {
		return OperatorResumeCommand{}, fmt.Errorf("%w: response does not belong to engine run", ErrInvalidOperatorInteraction)
	}
	task := taskEntity(request.RunID, request.TaskID)
	step := stepEntity(request.RunID, request.TaskID, request.StepID)
	taskState, stepState := engine.state.current(task), engine.state.current(step)
	waiting := taskState == "WAITING_FOR_OPERATOR" && (stepState == "WAITING_FOR_OPERATOR" || stepState == "RUNNING")
	resumed := taskState == "STEPS_RUNNING" && stepState == "RUNNING"
	if !waiting && !resumed {
		return OperatorResumeCommand{}, fmt.Errorf("%w: interaction resume requires waiting task and step", ErrInvalidState)
	}
	if engine.state.projection.ChainState != "PAUSED" {
		return OperatorResumeCommand{}, fmt.Errorf("%w: interaction resume requires PAUSED chain", ErrInvalidState)
	}
	inputs := []string{requestRecord.RecordHash, responseRecord.RecordHash}
	taskReason := "operator-interaction-required"
	taskHashes := []string{requestRecord.RecordHash}
	if resumed {
		taskReason = "operator-interaction-answered"
		taskHashes = inputs
	}
	if matches, matchErr := engine.latestTransitionMatches(ctx, task, taskState, taskReason, taskHashes...); matchErr != nil {
		return OperatorResumeCommand{}, matchErr
	} else if !matches {
		return OperatorResumeCommand{}, fmt.Errorf("%w: task belongs to another interaction", ErrInvalidOperatorInteraction)
	}
	if stepState == "WAITING_FOR_OPERATOR" {
		if matches, matchErr := engine.latestTransitionMatches(ctx, step, stepState, "operator-interaction-required", requestRecord.RecordHash); matchErr != nil {
			return OperatorResumeCommand{}, matchErr
		} else if !matches {
			return OperatorResumeCommand{}, fmt.Errorf("%w: waiting step belongs to another interaction", ErrInvalidOperatorInteraction)
		}
		if err := engine.transition(ctx, step, "RUNNING", inputs, "operator-interaction-answered"); err != nil {
			return OperatorResumeCommand{}, err
		}
	} else if matches, matchErr := engine.latestTransitionMatches(ctx, step, stepState, "operator-interaction-answered", inputs...); matchErr != nil {
		return OperatorResumeCommand{}, matchErr
	} else if !matches {
		return OperatorResumeCommand{}, fmt.Errorf("%w: running step belongs to another response", ErrInvalidOperatorInteraction)
	}
	if taskState == "WAITING_FOR_OPERATOR" {
		if err := engine.transition(ctx, task, "STEPS_RUNNING", inputs, "operator-interaction-answered"); err != nil {
			return OperatorResumeCommand{}, err
		}
	}
	return command, nil
}

func (engine *Engine) latestTransitionMatches(ctx context.Context, entity evidence.Entity, state, reason string, hashes ...string) (bool, error) {
	events, err := engine.options.Events.Load(ctx)
	if err != nil {
		return false, err
	}
	for index := len(events) - 1; index >= 0; index-- {
		event := events[index].Event
		if event.Entity != entity {
			continue
		}
		if event.ToState != state || event.Reason.Code != reason {
			return false, nil
		}
		for _, hash := range hashes {
			if !slices.Contains(event.InputEvidence, hash) {
				return false, nil
			}
		}
		return true, nil
	}
	return false, nil
}

func OpenOperatorInteractionTransition(taskState, stepState string) (nextTaskState, nextStepState string, err error) {
	if taskState != "STEPS_RUNNING" || stepState != "RUNNING" {
		return "", "", fmt.Errorf("%w: interaction open requires STEPS_RUNNING/RUNNING", ErrInvalidState)
	}
	return "WAITING_FOR_OPERATOR", "WAITING_FOR_OPERATOR", nil
}

func ResumeOperatorInteractionTransition(taskState, stepState string) (nextTaskState, nextStepState string, err error) {
	if taskState != "WAITING_FOR_OPERATOR" || stepState != "WAITING_FOR_OPERATOR" {
		return "", "", fmt.Errorf("%w: interaction resume requires WAITING_FOR_OPERATOR", ErrInvalidState)
	}
	return "STEPS_RUNNING", "RUNNING", nil
}

func validateOperatorInteractionRequest(request OperatorInteractionRequest) error {
	if request.Kind != "request" || !validOperatorInteractionIdentity(request.RecordID, request.InteractionID, request.RunID, request.TaskID, request.StepID, request.ConversationID, request.RequestID) || request.Attempt < 1 {
		return fmt.Errorf("%w: invalid request identity", ErrInvalidOperatorInteraction)
	}
	if request.RequestedBy.Type != "adapter" || !evidence.ValidID(request.RequestedBy.ID) || request.Operator.Type != "operator" || !evidence.ValidID(request.Operator.ID) {
		return fmt.Errorf("%w: invalid request actors", ErrInvalidOperatorInteraction)
	}
	if !evidence.ValidHash(request.WorkspaceHash) || !evidence.ValidHash(request.ContextHash) || !validUniqueHashSet(request.Evidence, true) {
		return fmt.Errorf("%w: invalid request hashes", ErrInvalidOperatorInteraction)
	}
	if _, err := time.Parse(operatorInteractionTimeLayout, request.RequestedAt); err != nil {
		return fmt.Errorf("%w: invalid request time", ErrInvalidOperatorInteraction)
	}
	if !validInteractionText(request.Question) || !validInteractionText(request.Reason) || !validInteractionText(request.Risk) || !validInteractionText(request.RequiredAction) || !validUniqueIDs(request.AllowedResponses) {
		return fmt.Errorf("%w: request must be structured", ErrInvalidOperatorInteraction)
	}
	return nil
}

func validateStandaloneOperatorInteractionResponse(response OperatorInteractionResponse) error {
	if response.Kind != "response" || !validOperatorInteractionIdentity(response.RecordID, response.InteractionID, response.RunID, response.TaskID, response.StepID, response.ConversationID, response.RequestID) || response.Attempt < 1 {
		return fmt.Errorf("%w: invalid response identity", ErrInvalidOperatorInteraction)
	}
	if response.RespondedBy.Type != "operator" || !evidence.ValidID(response.RespondedBy.ID) || !evidence.ValidHash(response.WorkspaceHash) || !evidence.ValidHash(response.ContextHash) || !evidence.ValidHash(response.RequestHash) || !validUniqueHashSet(response.Evidence, true) {
		return fmt.Errorf("%w: invalid response binding", ErrInvalidOperatorInteraction)
	}
	if _, err := time.Parse(operatorInteractionTimeLayout, response.RespondedAt); err != nil || !evidence.ValidID(response.Selection) {
		return fmt.Errorf("%w: invalid response content", ErrInvalidOperatorInteraction)
	}
	return nil
}

func validateOperatorInteractionResponse(request OperatorInteractionRequest, requestHash string, response OperatorInteractionResponse) error {
	if err := validateStandaloneOperatorInteractionResponse(response); err != nil {
		return err
	}
	requestedAt, _ := time.Parse(operatorInteractionTimeLayout, request.RequestedAt)
	respondedAt, _ := time.Parse(operatorInteractionTimeLayout, response.RespondedAt)
	if !respondedAt.After(requestedAt) || response.InteractionID != request.InteractionID || response.RespondedBy != request.Operator || response.RunID != request.RunID || response.TaskID != request.TaskID || response.StepID != request.StepID || response.Attempt != request.Attempt || response.WorkspaceHash != request.WorkspaceHash || response.ConversationID != request.ConversationID || response.ContextHash != request.ContextHash || response.RequestHash != requestHash || response.RequestID == request.RequestID || !slices.Contains(request.AllowedResponses, response.Selection) {
		return fmt.Errorf("%w: response does not match request", ErrInvalidOperatorInteraction)
	}
	return nil
}

func validOperatorInteractionIdentity(values ...string) bool {
	for _, value := range values {
		if !evidence.ValidID(value) {
			return false
		}
	}
	return true
}

func validInteractionText(value string) bool {
	return strings.TrimSpace(value) != "" && len(value) <= 4096 && !strings.ContainsAny(value, "\x00\r")
}
