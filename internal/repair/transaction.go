package repair

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	repairTransactionDomain = "proofrail:repair-transaction:1\n"
	repairStageDomain       = "proofrail:repair-stage:1\n"
	timestampLayout         = "2006-01-02T15:04:05.000Z"
)

var (
	ErrInvalidTransaction  = errors.New("invalid repair transaction")
	ErrInvalidStage        = errors.New("invalid repair stage")
	ErrStageOrder          = errors.New("invalid repair stage order")
	ErrStaleCandidate      = errors.New("stale repair candidate")
	ErrTerminalTransaction = errors.New("repair transaction is terminal")
)

type Stage struct {
	StageID               string         `json:"stageId"`
	Stage                 string         `json:"stage"`
	StartedAt             string         `json:"startedAt"`
	CompletedAt           string         `json:"completedAt"`
	Actor                 evidence.Actor `json:"actor"`
	PreviousStageHash     any            `json:"previousStageHash"`
	Outcome               string         `json:"outcome"`
	Evidence              []string       `json:"evidence"`
	ErrorEvidence         []string       `json:"errorEvidence"`
	StageHash             string         `json:"stageHash"`
	ParentSnapshotHash    string         `json:"parentSnapshotHash,omitempty"`
	BeforeManifestHash    string         `json:"beforeManifestHash,omitempty"`
	TargetIDs             []string       `json:"targetIds,omitempty"`
	WriterStopEvidence    []string       `json:"writerStopEvidence,omitempty"`
	LeaseEvidence         []string       `json:"leaseEvidence,omitempty"`
	CandidateManifestHash any            `json:"candidateManifestHash,omitempty"`
	DiffHash              string         `json:"diffHash,omitempty"`
	ScopeEvidence         []string       `json:"scopeEvidence,omitempty"`
	OwnershipEvidence     []string       `json:"ownershipEvidence,omitempty"`
	ValidationPlanHash    string         `json:"validationPlanHash,omitempty"`
	HookResultEvidence    []string       `json:"hookResultEvidence,omitempty"`
	ValidateStageHash     string         `json:"validateStageHash,omitempty"`
	TicketLedgerHash      string         `json:"ticketLedgerHash,omitempty"`
	ResultingManifestHash any            `json:"resultingManifestHash,omitempty"`
}

func (stage Stage) MarshalJSON() ([]byte, error) {
	fields := map[string]any{
		"stageId":           stage.StageID,
		"stage":             stage.Stage,
		"startedAt":         stage.StartedAt,
		"completedAt":       stage.CompletedAt,
		"actor":             stage.Actor,
		"previousStageHash": stage.PreviousStageHash,
		"outcome":           stage.Outcome,
		"evidence":          nonNilStrings(stage.Evidence),
		"errorEvidence":     nonNilStrings(stage.ErrorEvidence),
		"stageHash":         stage.StageHash,
	}
	switch stage.Stage {
	case "prepare":
		fields["parentSnapshotHash"] = stage.ParentSnapshotHash
		fields["beforeManifestHash"] = stage.BeforeManifestHash
		fields["targetIds"] = nonNilStrings(stage.TargetIDs)
		fields["writerStopEvidence"] = nonNilStrings(stage.WriterStopEvidence)
		fields["leaseEvidence"] = nonNilStrings(stage.LeaseEvidence)
		fields["candidateManifestHash"] = stage.CandidateManifestHash
	case "inspect":
		fields["candidateManifestHash"] = stage.CandidateManifestHash
		fields["diffHash"] = stage.DiffHash
		fields["scopeEvidence"] = nonNilStrings(stage.ScopeEvidence)
		fields["ownershipEvidence"] = nonNilStrings(stage.OwnershipEvidence)
	case "validate":
		fields["candidateManifestHash"] = stage.CandidateManifestHash
		fields["validationPlanHash"] = stage.ValidationPlanHash
		fields["hookResultEvidence"] = nonNilStrings(stage.HookResultEvidence)
	case "promote":
		fields["candidateManifestHash"] = stage.CandidateManifestHash
		fields["validateStageHash"] = stage.ValidateStageHash
		fields["ticketLedgerHash"] = stage.TicketLedgerHash
		fields["writerStopEvidence"] = nonNilStrings(stage.WriterStopEvidence)
		fields["leaseEvidence"] = nonNilStrings(stage.LeaseEvidence)
		fields["resultingManifestHash"] = stage.ResultingManifestHash
	}
	return json.Marshal(fields)
}

type TransactionBody struct {
	TransactionID    string  `json:"transactionId"`
	CreatedAt        string  `json:"createdAt"`
	RunID            string  `json:"runId"`
	TaskID           string  `json:"taskId"`
	SourceAttempt    int     `json:"sourceAttempt"`
	NextAttempt      int     `json:"nextAttempt"`
	LedgerID         string  `json:"ledgerId"`
	TicketLedgerHash string  `json:"ticketLedgerHash"`
	Fingerprint      string  `json:"fingerprint"`
	PolicyHash       string  `json:"policyHash"`
	Stages           []Stage `json:"stages"`
	Status           string  `json:"status"`
}

type TransactionRecord struct {
	SchemaVersion   string          `json:"schemaVersion"`
	Transaction     TransactionBody `json:"transaction"`
	TransactionHash string          `json:"transactionHash"`
}

type PrepareInput struct {
	StageID               string
	StartedAt             time.Time
	CompletedAt           time.Time
	Actor                 evidence.Actor
	Outcome               string
	Evidence              []string
	ErrorEvidence         []string
	ParentSnapshotHash    string
	BeforeManifestHash    string
	TargetIDs             []string
	WriterStopEvidence    []string
	LeaseEvidence         []string
	CandidateManifestHash string
}

type InspectInput struct {
	StageID               string
	StartedAt             time.Time
	CompletedAt           time.Time
	Actor                 evidence.Actor
	Outcome               string
	Evidence              []string
	ErrorEvidence         []string
	CandidateManifestHash string
	DiffHash              string
	ScopeEvidence         []string
	OwnershipEvidence     []string
}

type ValidateInput struct {
	StageID               string
	StartedAt             time.Time
	CompletedAt           time.Time
	Actor                 evidence.Actor
	Outcome               string
	Evidence              []string
	ErrorEvidence         []string
	CandidateManifestHash string
	ValidationPlanHash    string
	HookResultEvidence    []string
}

type PromoteInput struct {
	StageID               string
	StartedAt             time.Time
	CompletedAt           time.Time
	Actor                 evidence.Actor
	Outcome               string
	Evidence              []string
	ErrorEvidence         []string
	CandidateManifestHash string
	ValidateStageHash     string
	TicketLedgerHash      string
	WriterStopEvidence    []string
	LeaseEvidence         []string
	ResultingManifestHash string
}

type Transaction struct {
	body TransactionBody
}

func NewTransaction(transactionID, runID, taskID string, sourceAttempt, nextAttempt int, ledgerID, ticketLedgerHash, fingerprint, policyHash string, createdAt time.Time) (*Transaction, error) {
	if !evidence.ValidID(transactionID) || !evidence.ValidID(runID) || !evidence.ValidID(taskID) || sourceAttempt < 1 || nextAttempt <= sourceAttempt || !evidence.ValidID(ledgerID) {
		return nil, ErrInvalidTransaction
	}
	if !evidence.ValidHash(ticketLedgerHash) || !evidence.ValidHash(fingerprint) || !evidence.ValidHash(policyHash) {
		return nil, ErrInvalidTransaction
	}
	return &Transaction{body: TransactionBody{
		TransactionID:    transactionID,
		CreatedAt:        createdAt.UTC().Format(timestampLayout),
		RunID:            runID,
		TaskID:           taskID,
		SourceAttempt:    sourceAttempt,
		NextAttempt:      nextAttempt,
		LedgerID:         ledgerID,
		TicketLedgerHash: ticketLedgerHash,
		Fingerprint:      fingerprint,
		PolicyHash:       policyHash,
		Stages:           []Stage{},
		Status:           "in-progress",
	}}, nil
}

func (transaction *Transaction) AppendPrepare(input PrepareInput) error {
	if len(transaction.body.Stages) != 0 || transaction.isTerminal() {
		return ErrStageOrder
	}
	if !validStageCommon(input.StageID, input.Actor, input.StartedAt, input.CompletedAt) || !evidence.ValidHash(input.ParentSnapshotHash) || !evidence.ValidHash(input.BeforeManifestHash) || len(input.TargetIDs) == 0 || len(input.WriterStopEvidence) == 0 || len(input.LeaseEvidence) == 0 || len(input.Evidence) == 0 {
		return ErrInvalidStage
	}
	seenTargets := make(map[string]struct{}, len(input.TargetIDs))
	for _, targetID := range input.TargetIDs {
		if !evidence.ValidID(targetID) {
			return ErrInvalidStage
		}
		if _, exists := seenTargets[targetID]; exists {
			return ErrInvalidStage
		}
		seenTargets[targetID] = struct{}{}
	}
	if !allHashes(input.Evidence) || !allHashes(input.ErrorEvidence) || !allHashes(input.WriterStopEvidence) || !allHashes(input.LeaseEvidence) {
		return ErrInvalidStage
	}
	if !validateOutcome(input.Outcome, "completed", "failed", "uncertain") {
		return ErrInvalidStage
	}
	candidate := any(nil)
	if input.Outcome == "completed" {
		if !evidence.ValidHash(input.CandidateManifestHash) || len(input.ErrorEvidence) != 0 {
			return ErrInvalidStage
		}
		candidate = input.CandidateManifestHash
	} else {
		if input.CandidateManifestHash != "" || len(input.ErrorEvidence) == 0 {
			return ErrInvalidStage
		}
	}
	statement := map[string]any{
		"stageId":               input.StageID,
		"stage":                 "prepare",
		"startedAt":             input.StartedAt.UTC().Format(timestampLayout),
		"completedAt":           input.CompletedAt.UTC().Format(timestampLayout),
		"actor":                 input.Actor,
		"previousStageHash":     nil,
		"outcome":               input.Outcome,
		"evidence":              uniqueHashes(input.Evidence),
		"errorEvidence":         uniqueHashes(input.ErrorEvidence),
		"parentSnapshotHash":    input.ParentSnapshotHash,
		"beforeManifestHash":    input.BeforeManifestHash,
		"targetIds":             input.TargetIDs,
		"writerStopEvidence":    uniqueHashes(input.WriterStopEvidence),
		"leaseEvidence":         uniqueHashes(input.LeaseEvidence),
		"candidateManifestHash": candidate,
	}
	hash, err := digest(repairStageDomain, statement)
	if err != nil {
		return err
	}
	stage := Stage{
		StageID:               input.StageID,
		Stage:                 "prepare",
		StartedAt:             input.StartedAt.UTC().Format(timestampLayout),
		CompletedAt:           input.CompletedAt.UTC().Format(timestampLayout),
		Actor:                 input.Actor,
		PreviousStageHash:     nil,
		Outcome:               input.Outcome,
		Evidence:              uniqueHashes(input.Evidence),
		ErrorEvidence:         uniqueHashes(input.ErrorEvidence),
		StageHash:             hash,
		ParentSnapshotHash:    input.ParentSnapshotHash,
		BeforeManifestHash:    input.BeforeManifestHash,
		TargetIDs:             append([]string(nil), input.TargetIDs...),
		WriterStopEvidence:    uniqueHashes(input.WriterStopEvidence),
		LeaseEvidence:         uniqueHashes(input.LeaseEvidence),
		CandidateManifestHash: candidate,
	}
	transaction.body.Stages = append(transaction.body.Stages, stage)
	transaction.body.Status = statusFromOutcome(input.Outcome)
	return nil
}

func (transaction *Transaction) AppendInspect(input InspectInput) error {
	if err := transaction.ensureAdvance("inspect", 1, "prepare", "completed"); err != nil {
		return err
	}
	if !validStageCommon(input.StageID, input.Actor, input.StartedAt, input.CompletedAt) || !evidence.ValidHash(input.CandidateManifestHash) || !evidence.ValidHash(input.DiffHash) || len(input.ScopeEvidence) == 0 || len(input.OwnershipEvidence) == 0 || len(input.Evidence) == 0 {
		return ErrInvalidStage
	}
	if !allHashes(input.Evidence) || !allHashes(input.ErrorEvidence) || !allHashes(input.ScopeEvidence) || !allHashes(input.OwnershipEvidence) {
		return ErrInvalidStage
	}
	if !validateOutcome(input.Outcome, "passed", "failed", "uncertain") {
		return ErrInvalidStage
	}
	if !transaction.candidateMatches(input.CandidateManifestHash) {
		return ErrStaleCandidate
	}
	if input.Outcome == "passed" && len(input.ErrorEvidence) != 0 {
		return ErrInvalidStage
	}
	if input.Outcome != "passed" && len(input.ErrorEvidence) == 0 {
		return ErrInvalidStage
	}
	previous := transaction.body.Stages[len(transaction.body.Stages)-1].StageHash
	statement := map[string]any{
		"stageId":               input.StageID,
		"stage":                 "inspect",
		"startedAt":             input.StartedAt.UTC().Format(timestampLayout),
		"completedAt":           input.CompletedAt.UTC().Format(timestampLayout),
		"actor":                 input.Actor,
		"previousStageHash":     previous,
		"outcome":               input.Outcome,
		"evidence":              uniqueHashes(input.Evidence),
		"errorEvidence":         uniqueHashes(input.ErrorEvidence),
		"candidateManifestHash": input.CandidateManifestHash,
		"diffHash":              input.DiffHash,
		"scopeEvidence":         uniqueHashes(input.ScopeEvidence),
		"ownershipEvidence":     uniqueHashes(input.OwnershipEvidence),
	}
	hash, err := digest(repairStageDomain, statement)
	if err != nil {
		return err
	}
	stage := Stage{
		StageID:               input.StageID,
		Stage:                 "inspect",
		StartedAt:             input.StartedAt.UTC().Format(timestampLayout),
		CompletedAt:           input.CompletedAt.UTC().Format(timestampLayout),
		Actor:                 input.Actor,
		PreviousStageHash:     previous,
		Outcome:               input.Outcome,
		Evidence:              uniqueHashes(input.Evidence),
		ErrorEvidence:         uniqueHashes(input.ErrorEvidence),
		StageHash:             hash,
		CandidateManifestHash: input.CandidateManifestHash,
		DiffHash:              input.DiffHash,
		ScopeEvidence:         uniqueHashes(input.ScopeEvidence),
		OwnershipEvidence:     uniqueHashes(input.OwnershipEvidence),
	}
	transaction.body.Stages = append(transaction.body.Stages, stage)
	transaction.body.Status = statusFromOutcome(input.Outcome)
	return nil
}

func (transaction *Transaction) AppendValidate(input ValidateInput) error {
	if err := transaction.ensureAdvance("validate", 2, "inspect", "passed"); err != nil {
		return err
	}
	if !validStageCommon(input.StageID, input.Actor, input.StartedAt, input.CompletedAt) || !evidence.ValidHash(input.CandidateManifestHash) || !evidence.ValidHash(input.ValidationPlanHash) || len(input.HookResultEvidence) == 0 || len(input.Evidence) == 0 {
		return ErrInvalidStage
	}
	if !allHashes(input.Evidence) || !allHashes(input.ErrorEvidence) || !allHashes(input.HookResultEvidence) {
		return ErrInvalidStage
	}
	if !validateOutcome(input.Outcome, "passed", "failed", "uncertain") {
		return ErrInvalidStage
	}
	if !transaction.candidateMatches(input.CandidateManifestHash) {
		return ErrStaleCandidate
	}
	if input.Outcome == "passed" && len(input.ErrorEvidence) != 0 {
		return ErrInvalidStage
	}
	if input.Outcome != "passed" && len(input.ErrorEvidence) == 0 {
		return ErrInvalidStage
	}
	previous := transaction.body.Stages[len(transaction.body.Stages)-1].StageHash
	statement := map[string]any{
		"stageId":               input.StageID,
		"stage":                 "validate",
		"startedAt":             input.StartedAt.UTC().Format(timestampLayout),
		"completedAt":           input.CompletedAt.UTC().Format(timestampLayout),
		"actor":                 input.Actor,
		"previousStageHash":     previous,
		"outcome":               input.Outcome,
		"evidence":              uniqueHashes(input.Evidence),
		"errorEvidence":         uniqueHashes(input.ErrorEvidence),
		"candidateManifestHash": input.CandidateManifestHash,
		"validationPlanHash":    input.ValidationPlanHash,
		"hookResultEvidence":    uniqueHashes(input.HookResultEvidence),
	}
	hash, err := digest(repairStageDomain, statement)
	if err != nil {
		return err
	}
	stage := Stage{
		StageID:               input.StageID,
		Stage:                 "validate",
		StartedAt:             input.StartedAt.UTC().Format(timestampLayout),
		CompletedAt:           input.CompletedAt.UTC().Format(timestampLayout),
		Actor:                 input.Actor,
		PreviousStageHash:     previous,
		Outcome:               input.Outcome,
		Evidence:              uniqueHashes(input.Evidence),
		ErrorEvidence:         uniqueHashes(input.ErrorEvidence),
		StageHash:             hash,
		CandidateManifestHash: input.CandidateManifestHash,
		ValidationPlanHash:    input.ValidationPlanHash,
		HookResultEvidence:    uniqueHashes(input.HookResultEvidence),
	}
	transaction.body.Stages = append(transaction.body.Stages, stage)
	transaction.body.Status = statusFromOutcome(input.Outcome)
	return nil
}

func (transaction *Transaction) AppendPromote(input PromoteInput) error {
	if err := transaction.ensureAdvance("promote", 3, "validate", "passed"); err != nil {
		return err
	}
	if !validStageCommon(input.StageID, input.Actor, input.StartedAt, input.CompletedAt) || !evidence.ValidHash(input.CandidateManifestHash) || !evidence.ValidHash(input.ValidateStageHash) || !evidence.ValidHash(input.TicketLedgerHash) || len(input.WriterStopEvidence) == 0 || len(input.LeaseEvidence) == 0 || len(input.Evidence) == 0 {
		return ErrInvalidStage
	}
	if !allHashes(input.Evidence) || !allHashes(input.ErrorEvidence) || !allHashes(input.WriterStopEvidence) || !allHashes(input.LeaseEvidence) {
		return ErrInvalidStage
	}
	if !validateOutcome(input.Outcome, "completed", "failed", "uncertain") {
		return ErrInvalidStage
	}
	if !transaction.candidateMatches(input.CandidateManifestHash) {
		return ErrStaleCandidate
	}
	previous := transaction.body.Stages[len(transaction.body.Stages)-1].StageHash
	if input.ValidateStageHash != previous || input.TicketLedgerHash != transaction.body.TicketLedgerHash {
		return ErrInvalidStage
	}
	result := any(nil)
	if input.Outcome == "completed" {
		if !evidence.ValidHash(input.ResultingManifestHash) || len(input.ErrorEvidence) != 0 {
			return ErrInvalidStage
		}
		result = input.ResultingManifestHash
	} else {
		if input.ResultingManifestHash != "" || len(input.ErrorEvidence) == 0 {
			return ErrInvalidStage
		}
	}
	statement := map[string]any{
		"stageId":               input.StageID,
		"stage":                 "promote",
		"startedAt":             input.StartedAt.UTC().Format(timestampLayout),
		"completedAt":           input.CompletedAt.UTC().Format(timestampLayout),
		"actor":                 input.Actor,
		"previousStageHash":     previous,
		"outcome":               input.Outcome,
		"evidence":              uniqueHashes(input.Evidence),
		"errorEvidence":         uniqueHashes(input.ErrorEvidence),
		"candidateManifestHash": input.CandidateManifestHash,
		"validateStageHash":     input.ValidateStageHash,
		"ticketLedgerHash":      input.TicketLedgerHash,
		"writerStopEvidence":    uniqueHashes(input.WriterStopEvidence),
		"leaseEvidence":         uniqueHashes(input.LeaseEvidence),
		"resultingManifestHash": result,
	}
	hash, err := digest(repairStageDomain, statement)
	if err != nil {
		return err
	}
	stage := Stage{
		StageID:               input.StageID,
		Stage:                 "promote",
		StartedAt:             input.StartedAt.UTC().Format(timestampLayout),
		CompletedAt:           input.CompletedAt.UTC().Format(timestampLayout),
		Actor:                 input.Actor,
		PreviousStageHash:     previous,
		Outcome:               input.Outcome,
		Evidence:              uniqueHashes(input.Evidence),
		ErrorEvidence:         uniqueHashes(input.ErrorEvidence),
		StageHash:             hash,
		CandidateManifestHash: input.CandidateManifestHash,
		ValidateStageHash:     input.ValidateStageHash,
		TicketLedgerHash:      input.TicketLedgerHash,
		WriterStopEvidence:    uniqueHashes(input.WriterStopEvidence),
		LeaseEvidence:         uniqueHashes(input.LeaseEvidence),
		ResultingManifestHash: result,
	}
	transaction.body.Stages = append(transaction.body.Stages, stage)
	transaction.body.Status = statusFromOutcome(input.Outcome)
	return nil
}

func (transaction *Transaction) Snapshot() (TransactionRecord, error) {
	if len(transaction.body.Stages) == 0 {
		return TransactionRecord{}, ErrInvalidTransaction
	}
	hash, err := digest(repairTransactionDomain, transaction.body)
	if err != nil {
		return TransactionRecord{}, err
	}
	return TransactionRecord{SchemaVersion: evidence.SchemaVersion, Transaction: transaction.body, TransactionHash: hash}, nil
}

func (transaction *Transaction) ensureAdvance(nextStage string, expectedLen int, prevStage, prevOutcome string) error {
	if transaction.isTerminal() {
		return ErrTerminalTransaction
	}
	if len(transaction.body.Stages) != expectedLen {
		return ErrStageOrder
	}
	previous := transaction.body.Stages[len(transaction.body.Stages)-1]
	if previous.Stage != prevStage || previous.Outcome != prevOutcome {
		return ErrStageOrder
	}
	if nextStage == "promote" && transaction.body.Status != "in-progress" {
		return ErrStageOrder
	}
	return nil
}

func (transaction *Transaction) candidateMatches(candidateHash string) bool {
	if len(transaction.body.Stages) == 0 {
		return false
	}
	root, ok := transaction.body.Stages[0].CandidateManifestHash.(string)
	return ok && root == candidateHash
}

func (transaction *Transaction) isTerminal() bool {
	return transaction.body.Status == "failed" || transaction.body.Status == "uncertain" || transaction.body.Status == "completed"
}

func statusFromOutcome(outcome string) string {
	switch outcome {
	case "completed":
		return "in-progress"
	case "passed":
		return "in-progress"
	case "failed":
		return "failed"
	case "uncertain":
		return "uncertain"
	default:
		return "in-progress"
	}
}

func validStageCommon(stageID string, actor evidence.Actor, startedAt, completedAt time.Time) bool {
	if !evidence.ValidID(stageID) || !validActor(actor) {
		return false
	}
	if !completedAt.After(startedAt) {
		return false
	}
	return true
}

func validateOutcome(outcome string, allowed ...string) bool {
	for _, option := range allowed {
		if outcome == option {
			return true
		}
	}
	return false
}

func validActor(actor evidence.Actor) bool {
	if !evidence.ValidID(actor.ID) {
		return false
	}
	switch actor.Type {
	case "system", "operator", "agent", "policy":
		return true
	default:
		return false
	}
}

func allHashes(values []string) bool {
	for _, value := range values {
		if !evidence.ValidHash(value) {
			return false
		}
	}
	return true
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func uniqueHashes(values []string) []string {
	result := make([]string, 0, len(values))
	if len(values) < 2 {
		for _, value := range values {
			if evidence.ValidHash(value) {
				result = append(result, value)
			}
		}
		return result
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !evidence.ValidHash(value) {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func digest(domain string, body any) (string, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	canonical, err := evidence.Canonicalize(payload)
	if err != nil {
		return "", err
	}
	return evidence.Digest(domain, canonical), nil
}
