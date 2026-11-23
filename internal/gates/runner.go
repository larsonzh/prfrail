package gates

import (
	"context"
	"errors"
	"fmt"
	"mime"
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
	// ErrEffectClassificationUnavailable reports a hook whose effect class or
	// policy hash cannot be derived. The gate refuses instead of executing an
	// unclassified side effect (F5).
	ErrEffectClassificationUnavailable = errors.New("effect classification unavailable")
	// ErrEffectClassificationConflict reports a declared scope that contradicts the
	// frozen hook policy, so the request cannot be classified at all (F5).
	ErrEffectClassificationConflict = errors.New("effect classification conflicts with the frozen hook policy")
	// ErrWorkspaceRootMoved reports a declared workspace root that no longer names
	// the directory the runner bound. T4: the attempt work directory would otherwise
	// be a different object than the artifact root the receipt binds.
	ErrWorkspaceRootMoved = errors.New("gate workspace root is no longer the bound directory")
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

// gateWorkspaceRootMoved proves the caller's declared root still addresses the
// directory the runner bound (P85 F02/F06/F20).
//
// The handle is asserted against the identity the binding produced, and the
// declared NAME is resolved once and compared with that identity. Two questions
// are kept apart on purpose: "is the handle still the approved object" is answered
// by the handle, and "does the caller's name still mean that object" is the only
// thing a name resolution may be used for.
func gateWorkspaceRootMoved(workspace *evidence.BoundDir, boundID evidence.FileID, declaredRoot string) error {
	if err := evidence.VerifyBoundDir(workspace, boundID); err != nil {
		return err
	}
	current, err := evidence.FileIDAt(declaredRoot)
	if err != nil {
		return err
	}
	if current != boundID {
		return fmt.Errorf("the declared root %q no longer names the bound directory (%s != %s)", declaredRoot, current, boundID)
	}
	return nil
}

func (runner Runner) Run(ctx context.Context, request Request) (ResultRecord, error) {
	if runner.Executor == nil || runner.Scanner == nil || runner.Objects == nil || runner.Store == nil {
		return ResultRecord{}, ErrPolicyUnavailable
	}
	// F5 (P83 fix round): a direct Runner.Run caller must prove that the hook it
	// hands in IS the frozen definition the supplied digest names. Prepare only
	// validates the shape of the values, so without this comparison a caller could
	// substitute a hook whose classified fields (recovery guarantee, external
	// systems) differ from the frozen policy while still presenting the approved
	// digest. The digest covers every field of the frozen definition, so any
	// contradiction is refused here, before Prepare and before Execute.
	frozenDefinitionHash, err := DefinitionHash(request.Hook)
	if err != nil {
		return ResultRecord{}, fmt.Errorf("%w: the hook definition cannot be digested: %v", ErrEffectClassificationUnavailable, err)
	}
	if frozenDefinitionHash != request.HookDefinitionHash {
		return ResultRecord{}, fmt.Errorf("%w: the supplied hook definition hash %q is not the digest %q of the frozen hook definition", ErrEffectClassificationConflict, request.HookDefinitionHash, frozenDefinitionHash)
	}
	// F1 (P83 fix round): the workspace is bound from the trusted root BEFORE the
	// hook runs, and that one bound object serves the whole run, so a root or an
	// ancestor replaced while the hook executes cannot redirect a later read. A
	// request that declares no artifact globs opens nothing at all: a hook without
	// artifacts must not be refused by a platform that cannot bind a directory.
	// A bind failure still refuses collection below (fail-closed), but only when
	// artifacts were declared.
	var workspace *evidence.BoundDir
	var workspaceID evidence.FileID
	if len(request.Hook.ArtifactGlobs) > 0 {
		if bound, bindErr := evidence.OpenBoundDir(request.WorkspaceRoot); bindErr == nil {
			// P85 F02/F06/F20: the identity is captured from the bound handle, so the
			// per-attempt check compares the caller's declared name with the object
			// that was actually bound instead of resolving the same name twice.
			identity, identityErr := bound.FileID()
			if identityErr != nil {
				_ = bound.Close()
			} else {
				workspace = bound
				workspaceID = identity
				defer func() { _ = workspace.Close() }()
			}
		}
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
	// declaredRoot is the root the caller named. It is re-proven against the bound
	// handle on every attempt, so the check keeps observing the caller's namespace
	// even after the loop repointed request.WorkspaceRoot at the handle-derived name.
	declaredRoot := request.WorkspaceRoot
	for executionAttempt := request.ExecutionAttempt; executionAttempt <= maxAttempts; executionAttempt++ {
		request.ExecutionAttempt = executionAttempt
		// T4 (P83 fix round 6): the attempt's work directory and the artifact
		// enumeration must resolve through the SAME bound root the artifacts are read
		// through. Without this proof a retry could execute in a replacement directory
		// while the receipt binds the artifacts of the directory that was bound before
		// the first attempt. The root is re-proven per attempt - a swap inside one
		// attempt's hook is still caught by the handle-bound read - and a root that no
		// longer addresses the bound object refuses BEFORE the hook runs, so no side
		// effect lands in an unproven directory.
		if workspace != nil {
			if verifyErr := gateWorkspaceRootMoved(workspace, workspaceID, declaredRoot); verifyErr != nil {
				return ResultRecord{}, fmt.Errorf("%w: %v", ErrWorkspaceRootMoved, verifyErr)
			}
			boundRoot, rootErr := workspace.RootPath()
			if rootErr != nil {
				return ResultRecord{}, fmt.Errorf("%w: %v", ErrWorkspaceRootMoved, rootErr)
			}
			// Prepare derives the work directory from this field, so pointing it at the
			// bound object's own name is what binds the execution to the artifact root.
			request.WorkspaceRoot = boundRoot
		}
		attempted := clock().UTC()
		prepared, err := Prepare(request, runner.Executor.Capabilities(), environment())
		if err != nil {
			return ResultRecord{}, err
		}
		classification, err := deriveEffectClassification(request, attempted)
		if err != nil {
			return ResultRecord{}, err
		}
		// F5 (CONTRACTS.md:391 "read-only / locally discardable / external write are
		// separated; S1 refuses external writes and pauses on an unknown effect"):
		// admission is unconditional. A hook whose effect class cannot be derived
		// never reaches Execute, and a derived class that is not admitted never does
		// either - there is no "no scope, therefore skip the check" branch.
		if AdmitEffect(classification.Record) != EffectAdmit {
			denied := Execution{
				FinishedAt:     attempted,
				Outcome:        "start-failed",
				RunnerEvidence: []string{evidence.Digest("proofrail:gate-runner-evidence:1\n", []byte("effect-admission-denied"))},
				ErrorEvidence:  []string{evidence.Digest("proofrail:gate-error:1\n", []byte("effect admission denied by S1 policy"))},
			}
			result, err := runner.assess(ctx, request, ids(), attempted, denied, classification, workspace)
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
		result, err := runner.assess(ctx, request, ids(), attempted, execution, classification, workspace)
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

func (runner Runner) assess(ctx context.Context, request Request, resultID string, attempted time.Time, execution Execution, classification *EffectClassificationRecord, workspace *evidence.BoundDir) (Result, error) {
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
		artifacts, err := runner.collectBoundArtifacts(ctx, request, workspace)
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

// gateArtifactMaxBytes bounds a single collected artifact. An artifact that
// exceeds it is refused (never truncated into evidence): F1 collects through a
// handle bounded like every other evidence read.
const gateArtifactMaxBytes int64 = 256 << 20

// ErrArtifactOutOfWorkspace reports a glob match that resolves outside the bound
// workspace root. It is a refusal, never a skip.
var ErrArtifactOutOfWorkspace = errors.New("artifact resolves outside the bound workspace")

func (runner Runner) collectArtifacts(ctx context.Context, request Request) ([]Artifact, error) {
	// F1 (P83 fix round): a request that declares no artifact globs opens nothing.
	// Opening the workspace unconditionally made an otherwise passable
	// no-artifact hook fail on a platform that cannot bind a directory.
	if len(request.Hook.ArtifactGlobs) == 0 {
		return []Artifact{}, nil
	}
	workspace, err := evidence.OpenBoundDir(request.WorkspaceRoot)
	if err != nil {
		return nil, ErrArtifactMissing
	}
	defer func() { _ = workspace.Close() }()
	return runner.collectBoundArtifacts(ctx, request, workspace)
}

// collectBoundArtifacts collects artifacts through a workspace handle the caller
// bound BEFORE the hook executed (F1). The handle - not the path - is the root
// every artifact is opened against, so a root or ancestor replaced while the hook
// ran cannot redirect a read. A missing handle is a refusal, never a fallback to
// a path-based read.
func (runner Runner) collectBoundArtifacts(ctx context.Context, request Request, workspace *evidence.BoundDir) ([]Artifact, error) {
	if len(request.Hook.ArtifactGlobs) == 0 {
		return []Artifact{}, nil
	}
	if workspace == nil {
		return nil, ErrArtifactMissing
	}
	// T4 (P83 fix round 7): the enumeration runs through the bound object, not through
	// the caller's path. Round 6 derived a base name from the handle and then globbed
	// that PATH, so a root renamed away and re-created between the bind and the
	// collection let the glob describe a directory the reads never touch; every
	// component is now matched against the entries read from the bound handle, so the
	// receipt and the enumeration come from one object.
	var paths []string
	for _, pattern := range request.Hook.ArtifactGlobs {
		matches, err := boundArtifactGlob(workspace, pattern)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			return nil, ErrArtifactMissing
		}
		paths = append(paths, matches...)
	}
	sort.Strings(paths)
	// F1 (SECURITY.md:23 "the limit must be enforced by the runner/OS capability,
	// never by a prompt"): every artifact is opened relative to the bound handle,
	// so the check and the read describe the same object. A path-based Lstat
	// followed by a path-based ReadFile let a middle symlink resolve outside the
	// workspace between the two calls.
	artifacts := make([]Artifact, 0, len(paths))
	seen := map[string]struct{}{}
	for _, rel := range paths {
		// T4 (P83 fix round 7): the name the collector reads is the name the bound
		// enumeration produced, so there is no path-to-relative mapping that could
		// describe a second root.
		if _, exists := seen[rel]; exists {
			continue
		}
		seen[rel] = struct{}{}
		file, err := evidence.OpenBoundDirFile(workspace, rel, gateArtifactMaxBytes)
		if err != nil {
			return nil, ErrArtifactMissing
		}
		data, truncated, readErr := file.ReadAll()
		closeErr := file.Close()
		if readErr != nil || truncated {
			return nil, ErrArtifactMissing
		}
		if closeErr != nil {
			return nil, closeErr
		}
		redaction, err := runner.Scanner.Scan(ctx, filepath.Base(rel), data)
		if err != nil {
			return nil, ErrScanFailed
		}
		hash, err := runner.Objects.PutObject(ctx, data)
		if err != nil {
			return nil, err
		}
		mediaType := mime.TypeByExtension(filepath.Ext(rel))
		if mediaType == "" {
			mediaType = "application/octet-stream"
		} else if semicolon := strings.IndexByte(mediaType, ';'); semicolon >= 0 {
			mediaType = mediaType[:semicolon]
		}
		artifacts = append(artifacts, Artifact{ID: fmt.Sprintf("artifact-%d", len(artifacts)+1), ObjectHash: hash, MediaType: mediaType, ByteLength: int64(len(data)), Redaction: redaction})
	}
	return artifacts, nil
}

// boundArtifactGlob enumerates one artifact pattern inside the bound workspace
// object and returns the matching names relative to it (P83 fix round 7, item T4).
//
// Each component of the pattern is matched against the entries read from the bound
// handle, so a directory that was replaced after the bind can never contribute a
// name: the enumeration and the reads describe one object. A pattern that escapes
// the bound directory is refused, never silently skipped.
func boundArtifactGlob(workspace *evidence.BoundDir, pattern string) ([]string, error) {
	clean := filepath.Clean(filepath.FromSlash(pattern))
	if clean == "." || clean == ".." || filepath.IsAbs(clean) ||
		strings.HasPrefix(clean, ".."+string(filepath.Separator)) || strings.HasPrefix(clean, "/") {
		return nil, ErrArtifactOutOfWorkspace
	}
	current := []string{""}
	for _, component := range strings.Split(filepath.ToSlash(clean), "/") {
		if component == "" || component == "." || component == ".." {
			return nil, ErrArtifactOutOfWorkspace
		}
		next := make([]string, 0, len(current))
		for _, dir := range current {
			names, err := evidence.BoundDirNames(workspace, dir)
			if err != nil {
				// A directory that cannot be enumerated contributes no match, exactly as
				// the previous Glob skipped an unreadable directory; a pattern that ends
				// with no match refuses the collection above.
				continue
			}
			for _, name := range names {
				matched, matchErr := filepath.Match(component, name)
				if matchErr != nil {
					return nil, ErrArtifactMissing
				}
				if !matched {
					continue
				}
				if dir == "" {
					next = append(next, name)
					continue
				}
				next = append(next, dir+"/"+name)
			}
		}
		current = next
	}
	return current, nil
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

// deriveEffectClassification exports the effect class of the frozen hook policy
// and optionally reconciles it with the request's declared scope.
//
// F5: the classification is always derived from request.Hook, the frozen policy
// the executor was prepared against, because that is the only source the gate can
// prove. A declared scope is a cross-check, never an override:
//
//   - a scope effect class that contradicts the hook is a conflict, not a
//     reclassification;
//   - a scope policy hash that is not a valid hash is unavailable, never a
//     placeholder digest invented by the gate;
//   - a valid scope policy hash that differs from the hook policy hash is a
//     conflict, because the gate cannot prove a policy the hook does not name;
//   - the operation key comes from the scope when it names one and is otherwise
//     derived from the hook definition hash, so an observation exists for every
//     admitted run.
func deriveEffectClassification(request Request, createdAt time.Time) (*EffectClassificationRecord, error) {
	hook := request.Hook
	effectClass := strings.TrimSpace(hook.EffectClass)
	if !validEffectClasses[EffectClass(effectClass)] {
		return nil, fmt.Errorf("%w: the frozen hook policy declares effect class %q", ErrEffectClassificationUnavailable, hook.EffectClass)
	}
	scope := EffectScope{}
	hasScope := request.Effect != nil
	if hasScope {
		scope = *request.Effect
		if scope.EffectClass != "" && scope.EffectClass != effectClass {
			return nil, fmt.Errorf("%w: the declared scope class %q contradicts the frozen hook class %q", ErrEffectClassificationConflict, scope.EffectClass, effectClass)
		}
	}
	policyHash := hook.PolicyHash
	if hasScope && scope.PolicyHash != "" {
		if !evidence.ValidHash(scope.PolicyHash) {
			// The old code fabricated a placeholder digest here, which made an
			// unprovable policy look classified. A scope that names an invalid policy
			// hash is refused instead.
			return nil, fmt.Errorf("%w: the declared scope carries an invalid policy hash", ErrEffectClassificationUnavailable)
		}
		// F5 (P83 fix round): a declared policy identity that differs from the frozen
		// hook policy is a contradiction, not an alternative truth. The hook is the
		// only policy the gate can prove, so acting on a scope that names a different
		// policy would classify the wrong thing; the request is refused instead of
		// executed.
		if scope.PolicyHash != policyHash {
			return nil, fmt.Errorf("%w: the declared scope policy hash %q contradicts the frozen hook policy hash %q", ErrEffectClassificationConflict, scope.PolicyHash, policyHash)
		}
	}
	if !evidence.ValidHash(policyHash) {
		// A hook without a declared policy hash is identified by the digest of its
		// own frozen definition, which is what the result record already binds and
		// what Prepare validated the hook against. This is a real identity of the
		// policy in force - unlike the deleted placeholder, which was constant for
		// every policy of a run and therefore forgeable.
		policyHash = request.HookDefinitionHash
	}
	if !evidence.ValidHash(policyHash) {
		return nil, fmt.Errorf("%w: no valid policy identity is bound to the frozen hook policy", ErrEffectClassificationUnavailable)
	}
	// F5 (P83 fix round): the classified fields come ONLY from the frozen hook
	// policy. A declared scope is a cross-check, never an override: a scope that
	// states a different recovery guarantee or a different external-system set is
	// a contradiction and is refused, because the gate can only prove the policy
	// the frozen hook names.
	recoveryGuarantee := hook.RecoveryGuarantee
	if recoveryGuarantee == "" {
		recoveryGuarantee = defaultEffectRecoveryGuarantee(effectClass)
	}
	externalSystems := append([]string{}, hook.ExternalSystems...)
	if hasScope {
		if scope.RecoveryGuarantee != "" && scope.RecoveryGuarantee != recoveryGuarantee {
			return nil, fmt.Errorf("%w: the declared scope recovery guarantee %q contradicts the frozen hook policy %q", ErrEffectClassificationConflict, scope.RecoveryGuarantee, recoveryGuarantee)
		}
		if !sameStringSet(scope.ExternalSystems, externalSystems) {
			return nil, fmt.Errorf("%w: the declared scope external systems contradict the frozen hook policy", ErrEffectClassificationConflict)
		}
	}
	evidenceSet := []string{}
	if hasScope {
		evidenceSet = uniqueEvidenceHashes(scope.Evidence)
	}
	if len(evidenceSet) == 0 {
		evidenceSet = uniqueEvidenceHashes(hook.EffectEvidence)
	}
	operationKey := operationKeyFor(request)
	if len(evidenceSet) == 0 {
		evidenceSet = []string{evidence.Digest("proofrail:effect-classification-evidence:1\n", []byte(request.RunID+"\n"+request.TaskID+"\n"+request.StepID+"\n"+operationKey))}
	}
	classification := EffectClassification{
		RecordID:              effectClassificationID(request, operationKey),
		Kind:                  "classification",
		CreatedAt:             formatTimestamp(createdAt),
		PolicyHash:            policyHash,
		SubjectKind:           "hook",
		SubjectID:             request.Hook.ID,
		SubjectDefinitionHash: request.HookDefinitionHash,
		EffectClass:           effectClass,
		ExternalSystems:       externalSystems,
		RecoveryGuarantee:     recoveryGuarantee,
		Evidence:              evidenceSet,
	}
	record, err := NewEffectClassificationRecord(classification)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// operationKeyFor names the operation the effect records describe. A declared
// scope key wins; otherwise the frozen hook definition hash is the stable
// identity, because a hook is the only thing the gate actually executes.
//
// The key must be a schema identifier (lowercase alphanumerics and dashes, at
// most 64 characters), so the digest is rendered as the usual prefixed hex form
// rather than as a "sha256:" hash string.
func operationKeyFor(request Request) string {
	if request.Effect != nil && request.Effect.OperationKey != "" {
		return request.Effect.OperationKey
	}
	digest := evidence.Digest("proofrail:effect-operation-key:1\n", []byte(request.HookDefinitionHash))
	return "hook-" + digest[len("sha256:"):len("sha256:")+24]
}

// effectiveAuthorizationHash reports the authorization the observation binds: a
// declared scope value wins, then the frozen hook. When neither names one, the
// only basis the gate can prove is the frozen policy itself, so the observation
// binds a digest of that policy's definition hash - an identity of the policy in
// force, never a run-scoped value that any policy would share (F5).
func effectiveAuthorizationHash(request Request) string {
	if request.Effect != nil && request.Effect.AuthorizationHash != "" {
		return request.Effect.AuthorizationHash
	}
	if request.Hook.AuthorizationHash != "" {
		return request.Hook.AuthorizationHash
	}
	return evidence.Digest("proofrail:effect-authorization-basis:1\n", []byte(request.HookDefinitionHash))
}

func buildEffectObservation(request Request, classification EffectClassificationRecord, execution Execution, attempted time.Time) (EffectObservationRecord, string, error) {
	// F5: the observation is derived from the same source as the classification, so
	// it exists for every admitted run and not only for a caller that supplied a
	// scope. Removing the "scope required" error is what closes the branch that let
	// an unscoped request execute unclassified.
	operationKey := operationKeyFor(request)
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
		RecordID:           effectObservationID(request, operationKey, classification.RecordHash),
		Kind:               "observation",
		OccurredAt:         formatTimestamp(occurredAt),
		RunID:              request.RunID,
		TaskID:             request.TaskID,
		Attempt:            request.Attempt,
		ClassificationHash: classification.RecordHash,
		AuthorizationHash:  effectiveAuthorizationHash(request),
		EffectClass:        classification.Record.EffectClass,
		OperationKey:       operationKey,
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

// sameStringSet reports whether two lists name the same distinct values, ignoring
// order and duplicates. F5 (P83 fix round): a declared scope restates the frozen
// policy, so an equal set is a cross-check and any difference is a contradiction.
func sameStringSet(left, right []string) bool {
	leftSet := make(map[string]struct{}, len(left))
	for _, value := range left {
		leftSet[value] = struct{}{}
	}
	rightSet := make(map[string]struct{}, len(right))
	for _, value := range right {
		rightSet[value] = struct{}{}
	}
	if len(leftSet) != len(rightSet) {
		return false
	}
	for value := range leftSet {
		if _, ok := rightSet[value]; !ok {
			return false
		}
	}
	return true
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
