package taskdef

import (
	"errors"
	"fmt"
	"sort"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	DocumentationRequired   = "required"
	DocumentationIfAffected = "if-affected"
	DocumentationOptional   = "optional"
	DocumentationForbidden  = "forbidden"

	TargetClassSource        = "source"
	TargetClassTest          = "test"
	TargetClassConfig        = "config"
	TargetClassDocumentation = "documentation"
	TargetClassGenerated     = "generated"

	TemplateAID      = "template-a"
	TemplateAVersion = "1.0.0"

	templateAHashDomain = "proofrail:template-a:1\n"
)

var (
	ErrInvalidDocumentationInput = errors.New("invalid documentation input")
	templateALock                = mustTemplateALock()
)

type TemplateALock struct {
	ID      string
	Version string
	Hash    string
}

type DocumentationTarget struct {
	ID            string
	Class         string
	Tags          []string
	GeneratedFrom []string
	GeneratorHook string
}

type DocumentationImpactRule struct {
	ID             string
	WhenTargetTags []string
	RequireTargets []string
	RequiredHooks  []string
}

type GeneratedArtifactState struct {
	TargetID     string
	TemplateID   string
	Version      string
	TemplateHash string
	GeneratedBy  evidence.Actor
	ApprovedBy   evidence.Actor
	ManualEdit   bool
	Fresh        bool
}

type DocumentationCheckInput struct {
	ChainPolicy string
	TaskPolicy  string

	Targets       []DocumentationTarget
	ImpactRules   []DocumentationImpactRule
	ChangedTarget []string
	UpdatedTarget []string
	PassedHooks   []string
	Generated     []GeneratedArtifactState
}

type DocumentationIssue struct {
	Code     string
	TargetID string
	Detail   string
}

type DocumentationCheckResult struct {
	EffectivePolicy string
	RequiredTargets []string
	RequiredHooks   []string
	BlockingIssues  []DocumentationIssue
	Blocked         bool
}

type LanguageScopeResult struct {
	ScopeID  string
	Language string
	Harness  string
	Passed   bool
}

type TaskRecoveryDecision struct {
	Passed           bool
	RequiresRecovery bool
	FailedScopes     []string
}

func LockedTemplateA() TemplateALock {
	return templateALock
}

func CheckDocumentation(input DocumentationCheckInput) (DocumentationCheckResult, error) {
	policy, err := resolveDocumentationPolicy(input.ChainPolicy, input.TaskPolicy)
	if err != nil {
		return DocumentationCheckResult{}, err
	}

	targets, err := validateTargets(input.Targets)
	if err != nil {
		return DocumentationCheckResult{}, err
	}
	if err := validateTargetReferences(input.ChangedTarget, targets, "changed target"); err != nil {
		return DocumentationCheckResult{}, err
	}
	if err := validateTargetReferences(input.UpdatedTarget, targets, "updated target"); err != nil {
		return DocumentationCheckResult{}, err
	}

	rules, err := validateImpactRules(input.ImpactRules, targets)
	if err != nil {
		return DocumentationCheckResult{}, err
	}

	updatedSet := toSet(input.UpdatedTarget)
	passedHookSet := toSet(input.PassedHooks)
	changedTags := collectChangedTags(input.ChangedTarget, targets)

	requiredTargets := map[string]struct{}{}
	requiredHooks := map[string]struct{}{}
	for _, rule := range rules {
		if !intersects(changedTags, rule.WhenTargetTags) {
			continue
		}
		for _, targetID := range rule.RequireTargets {
			requiredTargets[targetID] = struct{}{}
		}
		for _, hookID := range rule.RequiredHooks {
			requiredHooks[hookID] = struct{}{}
		}
	}

	issues := make([]DocumentationIssue, 0)
	appendIssue := func(code, targetID, detail string) {
		issues = append(issues, DocumentationIssue{Code: code, TargetID: targetID, Detail: detail})
	}

	if (policy == DocumentationRequired || policy == DocumentationIfAffected) && len(rules) == 0 {
		appendIssue("documentation.impact-rule-required", "", "documentation policy requires explicit impact rules")
	}

	if policy == DocumentationRequired && !hasDocumentationUpdate(updatedSet, targets) {
		appendIssue("documentation.required-target-class-missing", "", "no documentation target update in a required policy task")
	}

	if policy == DocumentationForbidden {
		for targetID := range updatedSet {
			target := targets[targetID]
			if target.Class == TargetClassDocumentation || target.Class == TargetClassGenerated {
				appendIssue("documentation.forbidden-update", targetID, "documentation policy forbids documentation updates")
			}
		}
	}

	for targetID := range requiredTargets {
		if _, ok := updatedSet[targetID]; !ok {
			appendIssue("documentation.target-required-missing", targetID, "target required by impact rule is missing from updated targets")
		}
	}
	for hookID := range requiredHooks {
		if _, ok := passedHookSet[hookID]; !ok {
			appendIssue("documentation.required-hook-missing", "", fmt.Sprintf("required hook %q did not pass", hookID))
		}
	}

	for _, generated := range input.Generated {
		target, ok := targets[generated.TargetID]
		if !ok {
			return DocumentationCheckResult{}, fmt.Errorf("%w: generated target %q is unknown", ErrInvalidDocumentationInput, generated.TargetID)
		}
		if target.Class != TargetClassGenerated {
			return DocumentationCheckResult{}, fmt.Errorf("%w: target %q is not generated", ErrInvalidDocumentationInput, generated.TargetID)
		}
		if generated.TemplateID != TemplateAID || generated.Version != TemplateAVersion {
			appendIssue("generated.template-a-required", generated.TargetID, "generated artifact must use locked template A")
		}
		if generated.TemplateHash != templateALock.Hash {
			appendIssue("generated.template-lock-mismatch", generated.TargetID, "generated artifact template hash does not match the lock")
		}
		if generated.GeneratedBy.Type == "" || generated.GeneratedBy.ID == "" || generated.ApprovedBy.Type == "" || generated.ApprovedBy.ID == "" {
			appendIssue("generated.actor-missing", generated.TargetID, "generated artifacts require non-empty generator and approver actors")
		}
		if sameActor(generated.GeneratedBy, generated.ApprovedBy) {
			appendIssue("generated.self-approval", generated.TargetID, "generator and approver must be different actors")
		}
		if generated.ManualEdit {
			appendIssue("generated.manual-edit", generated.TargetID, "manual edits on generated artifacts are not allowed")
		}
		if !generated.Fresh {
			appendIssue("generated.stale", generated.TargetID, "generated artifact is stale")
		}
	}

	sortIssueList(issues)
	result := DocumentationCheckResult{
		EffectivePolicy: policy,
		RequiredTargets: sortedKeys(requiredTargets),
		RequiredHooks:   sortedKeys(requiredHooks),
		BlockingIssues:  dedupeIssues(issues),
	}
	result.Blocked = len(result.BlockingIssues) > 0
	return result, nil
}

func EvaluateTaskRecovery(scopes []LanguageScopeResult) (TaskRecoveryDecision, error) {
	if len(scopes) == 0 {
		return TaskRecoveryDecision{}, fmt.Errorf("%w: at least one language scope result is required", ErrInvalidDocumentationInput)
	}

	failed := map[string]struct{}{}
	for _, scope := range scopes {
		if !evidence.ValidID(scope.ScopeID) || !evidence.ValidID(scope.Language) || !evidence.ValidID(scope.Harness) {
			return TaskRecoveryDecision{}, fmt.Errorf("%w: invalid scope, language, or harness identifier", ErrInvalidDocumentationInput)
		}
		if !scope.Passed {
			failed[scope.ScopeID] = struct{}{}
		}
	}

	decision := TaskRecoveryDecision{Passed: len(failed) == 0, FailedScopes: sortedKeys(failed)}
	decision.RequiresRecovery = !decision.Passed
	return decision, nil
}

func resolveDocumentationPolicy(chainPolicy, taskPolicy string) (string, error) {
	policy := taskPolicy
	if policy == "" {
		policy = chainPolicy
	}
	if policy == "" {
		policy = DocumentationIfAffected
	}
	switch policy {
	case DocumentationRequired, DocumentationIfAffected, DocumentationOptional, DocumentationForbidden:
		return policy, nil
	default:
		return "", fmt.Errorf("%w: unknown documentation policy %q", ErrInvalidDocumentationInput, policy)
	}
}

func validateTargets(items []DocumentationTarget) (map[string]DocumentationTarget, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("%w: at least one target is required", ErrInvalidDocumentationInput)
	}
	targets := make(map[string]DocumentationTarget, len(items))
	for _, item := range items {
		if !evidence.ValidID(item.ID) {
			return nil, fmt.Errorf("%w: invalid target ID %q", ErrInvalidDocumentationInput, item.ID)
		}
		if _, exists := targets[item.ID]; exists {
			return nil, fmt.Errorf("%w: duplicate target ID %q", ErrInvalidDocumentationInput, item.ID)
		}
		switch item.Class {
		case TargetClassSource, TargetClassTest, TargetClassConfig, TargetClassDocumentation, TargetClassGenerated:
		default:
			return nil, fmt.Errorf("%w: target %q has invalid class %q", ErrInvalidDocumentationInput, item.ID, item.Class)
		}
		for _, tag := range item.Tags {
			if !evidence.ValidID(tag) {
				return nil, fmt.Errorf("%w: target %q has invalid tag %q", ErrInvalidDocumentationInput, item.ID, tag)
			}
		}
		targets[item.ID] = item
	}

	for _, item := range items {
		if item.Class != TargetClassGenerated {
			continue
		}
		if len(item.GeneratedFrom) == 0 || !evidence.ValidID(item.GeneratorHook) {
			return nil, fmt.Errorf("%w: generated target %q requires generatedFrom and generatorHook", ErrInvalidDocumentationInput, item.ID)
		}
		for _, sourceID := range item.GeneratedFrom {
			if _, exists := targets[sourceID]; !exists {
				return nil, fmt.Errorf("%w: generated target %q references unknown source target %q", ErrInvalidDocumentationInput, item.ID, sourceID)
			}
		}
	}
	return targets, nil
}

func validateImpactRules(items []DocumentationImpactRule, targets map[string]DocumentationTarget) ([]DocumentationImpactRule, error) {
	seen := map[string]struct{}{}
	for _, item := range items {
		if !evidence.ValidID(item.ID) {
			return nil, fmt.Errorf("%w: invalid impact rule ID %q", ErrInvalidDocumentationInput, item.ID)
		}
		if _, exists := seen[item.ID]; exists {
			return nil, fmt.Errorf("%w: duplicate impact rule ID %q", ErrInvalidDocumentationInput, item.ID)
		}
		seen[item.ID] = struct{}{}
		if len(item.WhenTargetTags) == 0 || len(item.RequireTargets) == 0 || len(item.RequiredHooks) == 0 {
			return nil, fmt.Errorf("%w: impact rule %q requires non-empty tags, targets, and hooks", ErrInvalidDocumentationInput, item.ID)
		}
		for _, tag := range item.WhenTargetTags {
			if !evidence.ValidID(tag) {
				return nil, fmt.Errorf("%w: impact rule %q has invalid tag %q", ErrInvalidDocumentationInput, item.ID, tag)
			}
		}
		for _, targetID := range item.RequireTargets {
			if _, exists := targets[targetID]; !exists {
				return nil, fmt.Errorf("%w: impact rule %q references unknown target %q", ErrInvalidDocumentationInput, item.ID, targetID)
			}
		}
		for _, hookID := range item.RequiredHooks {
			if !evidence.ValidID(hookID) {
				return nil, fmt.Errorf("%w: impact rule %q has invalid hook ID %q", ErrInvalidDocumentationInput, item.ID, hookID)
			}
		}
	}
	return items, nil
}

func validateTargetReferences(values []string, targets map[string]DocumentationTarget, label string) error {
	for _, value := range values {
		if _, exists := targets[value]; !exists {
			return fmt.Errorf("%w: %s %q is unknown", ErrInvalidDocumentationInput, label, value)
		}
	}
	return nil
}

func collectChangedTags(changed []string, targets map[string]DocumentationTarget) map[string]struct{} {
	tags := map[string]struct{}{}
	for _, targetID := range changed {
		for _, tag := range targets[targetID].Tags {
			tags[tag] = struct{}{}
		}
	}
	return tags
}

func hasDocumentationUpdate(updated map[string]struct{}, targets map[string]DocumentationTarget) bool {
	for targetID := range updated {
		if targets[targetID].Class == TargetClassDocumentation {
			return true
		}
	}
	return false
}

func sameActor(left, right evidence.Actor) bool {
	if left.Type == "" || left.ID == "" || right.Type == "" || right.ID == "" {
		return false
	}
	return left.Type == right.Type && left.ID == right.ID
}

func toSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[value] = struct{}{}
	}
	return result
}

func intersects(tags map[string]struct{}, values []string) bool {
	for _, value := range values {
		if _, exists := tags[value]; exists {
			return true
		}
	}
	return false
}

func sortedKeys(values map[string]struct{}) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortIssueList(items []DocumentationIssue) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Code != items[j].Code {
			return items[i].Code < items[j].Code
		}
		if items[i].TargetID != items[j].TargetID {
			return items[i].TargetID < items[j].TargetID
		}
		return items[i].Detail < items[j].Detail
	})
}

func dedupeIssues(items []DocumentationIssue) []DocumentationIssue {
	if len(items) < 2 {
		return items
	}
	result := make([]DocumentationIssue, 0, len(items))
	for _, item := range items {
		if len(result) == 0 {
			result = append(result, item)
			continue
		}
		last := result[len(result)-1]
		if last.Code == item.Code && last.TargetID == item.TargetID && last.Detail == item.Detail {
			continue
		}
		result = append(result, item)
	}
	return result
}

func mustTemplateALock() TemplateALock {
	manifest := map[string]any{
		"templateId":                 TemplateAID,
		"templateVersion":            TemplateAVersion,
		"runner":                     "process",
		"fields":                     []string{"executable", "args", "cwd", "envAllowlist"},
		"shellMode":                  "forbidden",
		"defaultAllowGeneratedHooks": false,
	}
	canonical, err := evidence.EncodeCanonical(manifest)
	if err != nil {
		panic(err)
	}
	return TemplateALock{ID: TemplateAID, Version: TemplateAVersion, Hash: evidence.Digest(templateAHashDomain, canonical)}
}
