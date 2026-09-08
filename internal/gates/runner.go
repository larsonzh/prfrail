package gates

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrGateFailed      = errors.New("gate failed")
	ErrArtifactMissing = errors.New("required artifact missing")
	ErrScanFailed      = errors.New("secret scan failed")
)

type Scanner interface {
	Scan(context.Context, string, []byte) (string, error)
}

type ObjectStore interface {
	PutObject(context.Context, []byte) (string, error)
}

type Runner struct {
	Executor Executor
	Scanner  Scanner
	Objects  ObjectStore
	Store    ResultStore
	Clock    func() time.Time
	IDs      func() string
	Env      func() []string
}

func (runner Runner) Run(ctx context.Context, request Request) (ResultRecord, error) {
	if runner.Executor == nil || runner.Scanner == nil || runner.Objects == nil || runner.Store == nil {
		return ResultRecord{}, ErrPolicyUnavailable
	}
	clock := runner.Clock
	if clock == nil {
		clock = time.Now
	}
	ids := runner.IDs
	if ids == nil {
		ids = func() string { return fmt.Sprintf("result-%d", clock().UTC().UnixNano()) }
	}
	environment := runner.Env
	if environment == nil {
		environment = CurrentEnvironment
	}
	maxAttempts := request.ExecutionAttempt
	if request.Hook.OnFail.Action == "retry" {
		maxAttempts = request.Hook.OnFail.MaxAttempts
	}
	var record ResultRecord
	for executionAttempt := request.ExecutionAttempt; executionAttempt <= maxAttempts; executionAttempt++ {
		request.ExecutionAttempt = executionAttempt
		attempted := clock().UTC()
		prepared, err := Prepare(request, runner.Executor.Capabilities(), environment())
		if err != nil {
			return ResultRecord{}, err
		}
		classification, err := buildEffectClassification(request, attempted)
		if err != nil {
			return ResultRecord{}, err
		}
		if classification != nil && AdmitEffect(classification.Record) != EffectAdmit {
			denied := Execution{
				FinishedAt:     attempted,
				Outcome:        "start-failed",
				RunnerEvidence: []string{evidence.Digest("proofrail:gate-runner-evidence:1\n", []byte("effect-admission-denied"))},
				ErrorEvidence:  []string{evidence.Digest("proofrail:gate-error:1\n", []byte("effect admission denied by S1 policy"))},
			}
			result, err := runner.assess(ctx, request, ids(), attempted, denied, classification)
			if err != nil {
				return ResultRecord{}, err
			}
			record, err = makeRecord(result)
			if err != nil {
				return ResultRecord{}, err
			}
			if err := runner.Store.Append(ctx, record); err != nil {
				return ResultRecord{}, err
			}
			return record, ErrGateFailed
		}
		execution := runner.Executor.Execute(ctx, prepared, request.Hook)
		result, err := runner.assess(ctx, request, ids(), attempted, execution, classification)
		if err != nil {
			return ResultRecord{}, err
		}
		record, err = makeRecord(result)
		if err != nil {
			return ResultRecord{}, err
		}
		if err := runner.Store.Append(ctx, record); err != nil {
			return ResultRecord{}, err
		}
		if result.Assessment == "passed" {
			return record, nil
		}
		if result.PolicyDisposition != "retry" || executionAttempt == maxAttempts {
			return record, ErrGateFailed
		}
	}
	return record, ErrGateFailed
}

func (runner Runner) assess(ctx context.Context, request Request, resultID string, attempted time.Time, execution Execution, classification *EffectClassificationRecord) (Result, error) {
	if len(execution.RunnerEvidence) == 0 {
		execution.Outcome, execution.StartedAt, execution.ExitCode = "start-failed", nil, nil
	}
	if (execution.Outcome == "timed-out" || execution.Outcome == "resource-limited" || execution.Outcome == "termination-uncertain") && len(execution.TerminationEvidence) == 0 {
		execution.Outcome, execution.ExitCode = "termination-uncertain", nil
	}
	result := Result{ResultID: resultID, RecordedBy: evidence.Actor{Type: "system", ID: "proofrail"}, RunID: request.RunID, TaskID: request.TaskID, StepID: request.StepID, Attempt: request.Attempt, HookID: request.Hook.ID, HookDefinitionHash: request.HookDefinitionHash, ExecutionAttempt: request.ExecutionAttempt, AttemptedAt: formatTimestamp(attempted), StartedAt: formatOptionalTimestamp(execution.StartedAt), FinishedAt: formatTimestamp(execution.FinishedAt), ExecutionOutcome: execution.Outcome, ExitCode: execution.ExitCode, Assessment: "failed", PolicyDisposition: disposition(request.Hook.OnFail), Artifacts: []Artifact{}, RunnerEvidence: nonnil(execution.RunnerEvidence), TerminationEvidence: nonnil(execution.TerminationEvidence), ErrorEvidence: nonnil(execution.ErrorEvidence)}
	failure := executionFailure(execution)
	stdout, stdoutErr := runner.scanOutput(ctx, "stdout", execution.Stdout, request.Hook.Resources.OutputBytes, execution.OutputTruncated)
	stderr, stderrErr := runner.scanOutput(ctx, "stderr", execution.Stderr, request.Hook.Resources.OutputBytes, execution.OutputTruncated)
	if stdoutErr != nil && !errors.Is(stdoutErr, ErrScanFailed) {
		return Result{}, stdoutErr
	}
	if stderrErr != nil && !errors.Is(stderrErr, ErrScanFailed) {
		return Result{}, stderrErr
	}
	result.Stdout, result.Stderr = stdout, stderr
	if failure == "" && (stdoutErr != nil || stderrErr != nil) {
		failure = "scan-failed"
	}
	if failure == "" {
		artifacts, err := runner.collectArtifacts(ctx, request)
		result.Artifacts = artifacts
		if errors.Is(err, ErrArtifactMissing) {
			failure = "artifact-missing"
		} else if err != nil {
			failure = "scan-failed"
		}
	}
	if failure == "" {
		result.Assessment = "passed"
		result.PolicyDisposition = "pass"
		result.ErrorEvidence = []string{}
	}
	if failure != "" {
		result.FailureKind = &failure
		if len(result.ErrorEvidence) == 0 {
			result.ErrorEvidence = []string{evidence.Digest("proofrail:gate-error:1\n", []byte(failure))}
		}
	}
	if classification != nil {
		observation, recoveryAction, err := buildEffectObservation(request, *classification, execution, attempted)
		if err != nil {
			return Result{}, err
		}
		hash := observation.RecordHash
		result.EffectObservationHash = &hash
		result.EffectRecoveryAction = &recoveryAction
		if recoveryAction != "none-required" {
			recoveryEvidence := evidence.Digest("proofrail:effect-recovery-action:1\n", []byte(recoveryAction+"\n"+hash))
			result.ErrorEvidence = uniqueEvidenceHashes(append(result.ErrorEvidence, recoveryEvidence))
		}
	}
	return result, nil
}

func (runner Runner) scanOutput(ctx context.Context, name string, data []byte, limit int64, truncated bool) (*Output, error) {
	if int64(len(data)) > limit {
		return nil, ErrScanFailed
	}
	redaction, err := runner.Scanner.Scan(ctx, name, data)
	if err != nil {
		return nil, ErrScanFailed
	}
	hash, err := runner.Objects.PutObject(ctx, data)
	if err != nil {
		return nil, err
	}
	capture := "complete"
	if truncated {
		capture = "truncated"
	}
	return &Output{ObjectHash: hash, ByteLength: int64(len(data)), Capture: capture, Redaction: redaction}, nil
}

func (runner Runner) collectArtifacts(ctx context.Context, request Request) ([]Artifact, error) {
	var paths []string
	for _, pattern := range request.Hook.ArtifactGlobs {
		matches, err := filepath.Glob(filepath.Join(request.WorkspaceRoot, filepath.FromSlash(pattern)))
		if err != nil || len(matches) == 0 {
			return nil, ErrArtifactMissing
		}
		paths = append(paths, matches...)
	}
	sort.Strings(paths)
	artifacts := make([]Artifact, 0, len(paths))
	seen := map[string]struct{}{}
	for _, path := range paths {
		if _, exists := seen[path]; exists {
			continue
		}
		seen[path] = struct{}{}
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() {
			return nil, ErrArtifactMissing
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		redaction, err := runner.Scanner.Scan(ctx, filepath.Base(path), data)
		if err != nil {
			return nil, ErrScanFailed
		}
		hash, err := runner.Objects.PutObject(ctx, data)
		if err != nil {
			return nil, err
		}
		mediaType := mime.TypeByExtension(filepath.Ext(path))
		if mediaType == "" {
			mediaType = "application/octet-stream"
		} else if semicolon := strings.IndexByte(mediaType, ';'); semicolon >= 0 {
			mediaType = mediaType[:semicolon]
		}
		artifacts = append(artifacts, Artifact{ID: fmt.Sprintf("artifact-%d", len(artifacts)+1), ObjectHash: hash, MediaType: mediaType, ByteLength: int64(len(data)), Redaction: redaction})
	}
	return artifacts, nil
}

func executionFailure(execution Execution) string {
	switch execution.Outcome {
	case "exited":
		if execution.ExitCode != nil && *execution.ExitCode == 0 {
			return ""
		}
		return "exit-nonzero"
	case "start-failed", "timed-out", "resource-limited", "termination-uncertain":
		return execution.Outcome
	default:
		return "start-failed"
	}
}

func disposition(policy FailurePolicy) string { return policy.Action }
func formatTimestamp(value time.Time) string  { return value.UTC().Format("2006-01-02T15:04:05.000Z") }
func formatOptionalTimestamp(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := formatTimestamp(*value)
	return &formatted
}
func nonnil(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func buildEffectClassification(request Request, createdAt time.Time) (*EffectClassificationRecord, error) {
	if request.Effect == nil {
		return nil, nil
	}
	scope := *request.Effect
	policyHash := scope.PolicyHash
	if !evidence.ValidHash(policyHash) {
		policyHash = evidence.Digest("proofrail:effect-policy-placeholder:1\n", []byte(request.HookDefinitionHash+"\n"+request.RunID+"\n"+request.TaskID+"\n"+request.StepID))
	}
	recoveryGuarantee := scope.RecoveryGuarantee
	if recoveryGuarantee == "" {
		recoveryGuarantee = defaultEffectRecoveryGuarantee(scope.EffectClass)
	}
	evidenceSet := uniqueEvidenceHashes(scope.Evidence)
	if len(evidenceSet) == 0 {
		evidenceSet = []string{evidence.Digest("proofrail:effect-classification-evidence:1\n", []byte(request.RunID+"\n"+request.TaskID+"\n"+request.StepID+"\n"+scope.OperationKey))}
	}
	classification := EffectClassification{
		RecordID:              effectClassificationID(request, scope.OperationKey),
		Kind:                  "classification",
		CreatedAt:             formatTimestamp(createdAt),
		PolicyHash:            policyHash,
		SubjectKind:           "hook",
		SubjectID:             request.Hook.ID,
		SubjectDefinitionHash: request.HookDefinitionHash,
		EffectClass:           scope.EffectClass,
		ExternalSystems:       append([]string{}, scope.ExternalSystems...),
		RecoveryGuarantee:     recoveryGuarantee,
		Evidence:              evidenceSet,
	}
	record, err := NewEffectClassificationRecord(classification)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func buildEffectObservation(request Request, classification EffectClassificationRecord, execution Execution, attempted time.Time) (EffectObservationRecord, string, error) {
	if request.Effect == nil {
		return EffectObservationRecord{}, "", fmt.Errorf("effect scope required for observation")
	}
	scope := *request.Effect
	occurredAt := execution.FinishedAt
	if occurredAt.IsZero() {
		occurredAt = attempted
	}
	outcome := effectObservationOutcome(execution)
	evidenceSet := []string{}
	errorEvidence := []string{}
	if outcome == "completed" {
		evidenceSet = uniqueEvidenceHashes(append(evidenceSet, execution.RunnerEvidence...))
	} else {
		errorEvidence = uniqueEvidenceHashes(append(errorEvidence, execution.ErrorEvidence...))
		if len(errorEvidence) == 0 {
			errorEvidence = []string{evidence.Digest("proofrail:effect-observation-error:1\n", []byte(request.RunID+"\n"+request.TaskID+"\n"+request.StepID+"\n"+execution.Outcome))}
		}
		evidenceSet = uniqueEvidenceHashes(append(evidenceSet, execution.RunnerEvidence...))
	}
	observation := EffectObservation{
		RecordID:           effectObservationID(request, scope.OperationKey, classification.RecordHash),
		Kind:               "observation",
		OccurredAt:         formatTimestamp(occurredAt),
		RunID:              request.RunID,
		TaskID:             request.TaskID,
		Attempt:            request.Attempt,
		ClassificationHash: classification.RecordHash,
		AuthorizationHash:  scope.AuthorizationHash,
		EffectClass:        classification.Record.EffectClass,
		OperationKey:       scope.OperationKey,
		Outcome:            outcome,
		Evidence:           evidenceSet,
		ErrorEvidence:      errorEvidence,
	}
	record, err := NewEffectObservationRecord(observation)
	if err != nil {
		return EffectObservationRecord{}, "", err
	}
	recoveryAction := ObservationRecoveryAction(record.Record)
	return record, recoveryAction, nil
}

func effectObservationOutcome(execution Execution) string {
	if execution.Outcome == "exited" && execution.ExitCode != nil && *execution.ExitCode == 0 {
		return "completed"
	}
	switch execution.Outcome {
	case "start-failed", "timed-out", "resource-limited", "termination-uncertain":
		return "uncertain"
	default:
		return "failed"
	}
}

func defaultEffectRecoveryGuarantee(effectClass string) string {
	switch EffectClass(effectClass) {
	case EffectReadOnly:
		return "none-required"
	case EffectLocalDiscardable:
		return "discard-workspace"
	case EffectExternalWrite:
		return "reconcile-only"
	default:
		return "reconcile-only"
	}
}

func effectClassificationID(request Request, operationKey string) string {
	suffix := evidence.Digest("proofrail:effect-classification-id:1\n", []byte(request.RunID+"\n"+request.TaskID+"\n"+request.StepID+"\n"+operationKey))
	return "class-" + suffix[len("sha256:"):len("sha256:")+24]
}

func effectObservationID(request Request, operationKey, classificationHash string) string {
	suffix := evidence.Digest("proofrail:effect-observation-id:1\n", []byte(request.RunID+"\n"+request.TaskID+"\n"+request.StepID+"\n"+operationKey+"\n"+classificationHash+"\n"+fmt.Sprintf("%d", request.ExecutionAttempt)))
	return "obs-" + suffix[len("sha256:"):len("sha256:")+24]
}

func uniqueEvidenceHashes(values []string) []string {
	if len(values) < 2 {
		if values == nil {
			return []string{}
		}
		return values
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
