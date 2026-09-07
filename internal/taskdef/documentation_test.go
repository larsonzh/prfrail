package taskdef

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestCheckDocumentationIfAffectedPassesWithFreshGeneratedArtifacts(t *testing.T) {
	input := validDocumentationInput()
	result, err := CheckDocumentation(input)
	if err != nil {
		t.Fatal(err)
	}
	if result.EffectivePolicy != DocumentationIfAffected {
		t.Fatalf("unexpected policy: %s", result.EffectivePolicy)
	}
	if result.Blocked {
		t.Fatalf("unexpected blocking issues: %+v", result.BlockingIssues)
	}
	if len(result.RequiredTargets) != 2 || result.RequiredTargets[0] != "cli-docs" || result.RequiredTargets[1] != "cli-help-generated" {
		t.Fatalf("unexpected required targets: %+v", result.RequiredTargets)
	}
	if len(result.RequiredHooks) != 2 || result.RequiredHooks[0] != "docs-check" || result.RequiredHooks[1] != "generated-freshness-check" {
		t.Fatalf("unexpected required hooks: %+v", result.RequiredHooks)
	}
}

func TestCheckDocumentationBlocksMissingDocumentationAndHook(t *testing.T) {
	input := validDocumentationInput()
	input.UpdatedTarget = []string{"api-source", "cli-help-generated"}
	input.PassedHooks = []string{"generated-freshness-check"}

	result, err := CheckDocumentation(input)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Blocked {
		t.Fatal("expected documentation checker to block")
	}
	assertIssueCodes(t, result.BlockingIssues, "documentation.required-hook-missing", "documentation.target-required-missing")
}

func TestCheckDocumentationBlocksGeneratedViolations(t *testing.T) {
	cases := []struct {
		name       string
		mutate     func(*DocumentationCheckInput)
		expectCode string
	}{
		{
			name: "template-required",
			mutate: func(input *DocumentationCheckInput) {
				input.Generated[0].TemplateID = "template-b"
			},
			expectCode: "generated.template-a-required",
		},
		{
			name: "template-lock-mismatch",
			mutate: func(input *DocumentationCheckInput) {
				input.Generated[0].TemplateHash = evidence.Digest("", []byte("mismatch"))
			},
			expectCode: "generated.template-lock-mismatch",
		},
		{
			name: "self-approval",
			mutate: func(input *DocumentationCheckInput) {
				actor := evidence.Actor{Type: "operator", ID: "same-user"}
				input.Generated[0].GeneratedBy = actor
				input.Generated[0].ApprovedBy = actor
			},
			expectCode: "generated.self-approval",
		},
		{
			name: "missing-approver",
			mutate: func(input *DocumentationCheckInput) {
				input.Generated[0].ApprovedBy = evidence.Actor{}
			},
			expectCode: "generated.actor-missing",
		},
		{
			name: "manual-edit",
			mutate: func(input *DocumentationCheckInput) {
				input.Generated[0].ManualEdit = true
			},
			expectCode: "generated.manual-edit",
		},
		{
			name: "stale",
			mutate: func(input *DocumentationCheckInput) {
				input.Generated[0].Fresh = false
			},
			expectCode: "generated.stale",
		},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			input := validDocumentationInput()
			item.mutate(&input)
			result, err := CheckDocumentation(input)
			if err != nil {
				t.Fatal(err)
			}
			if !result.Blocked {
				t.Fatal("expected generated artifact violation")
			}
			assertIssueCodes(t, result.BlockingIssues, item.expectCode)
		})
	}
}

func TestCheckDocumentationRequiresImpactRulesWhenPolicyNeedsDocumentation(t *testing.T) {
	input := validDocumentationInput()
	input.ChainPolicy = DocumentationRequired
	input.ImpactRules = nil
	input.TaskPolicy = ""

	result, err := CheckDocumentation(input)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Blocked {
		t.Fatal("expected blocking when required policy has no impact rules")
	}
	assertIssueCodes(t, result.BlockingIssues, "documentation.impact-rule-required")
}

func TestEvaluateTaskRecoveryForCAndGoScopes(t *testing.T) {
	decision, err := EvaluateTaskRecovery([]LanguageScopeResult{
		{ScopeID: "c-main", Language: "c", Harness: "c-standard", Passed: true},
		{ScopeID: "go-main", Language: "go", Harness: "go-standard", Passed: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !decision.Passed || decision.RequiresRecovery || len(decision.FailedScopes) != 0 {
		t.Fatalf("unexpected pass decision: %+v", decision)
	}

	decision, err = EvaluateTaskRecovery([]LanguageScopeResult{
		{ScopeID: "c-main", Language: "c", Harness: "c-standard", Passed: false},
		{ScopeID: "go-main", Language: "go", Harness: "go-standard", Passed: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Passed || !decision.RequiresRecovery || len(decision.FailedScopes) != 1 || decision.FailedScopes[0] != "c-main" {
		t.Fatalf("expected whole-task recovery requirement, got %+v", decision)
	}
}

func TestHarnessProfilesProvideGenericCGoAndTemplateALock(t *testing.T) {
	templateParameter := "template-a@1.0.0"
	root := repositoryRoot(t)
	files := []struct {
		path      string
		harnessID string
	}{
		{path: filepath.Join(root, "harnesses", "generic", "harness.json"), harnessID: "generic-standard"},
		{path: filepath.Join(root, "harnesses", "c", "harness.json"), harnessID: "c-standard"},
		{path: filepath.Join(root, "harnesses", "go", "harness.json"), harnessID: "go-standard"},
	}

	for _, item := range files {
		data, err := os.ReadFile(item.path)
		if err != nil {
			t.Fatal(err)
		}
		var registry harnessRegistryFile
		if err := evidence.DecodeStrictJSON(data, &registry); err != nil {
			t.Fatalf("%s decode failed: %v", item.path, err)
		}
		if registry.SchemaVersion != "1.0.0" || len(registry.Harnesses) != 1 {
			t.Fatalf("unexpected harness registry shape for %s: %+v", item.path, registry)
		}
		harness := registry.Harnesses[0]
		if harness.ID != item.harnessID {
			t.Fatalf("unexpected harness id for %s: %s", item.path, harness.ID)
		}
		defaults := map[string]any{}
		for _, parameter := range harness.Parameters {
			defaults[parameter.ID] = parameter.Default
		}
		if value, ok := defaults["generated-hook-template"].(string); !ok || value != templateParameter {
			t.Fatalf("unexpected generated-hook-template default for %s: %#v", item.path, defaults["generated-hook-template"])
		}
		if value, ok := defaults["allow-generated-hooks"].(bool); !ok || value {
			t.Fatalf("unexpected allow-generated-hooks default for %s: %#v", item.path, defaults["allow-generated-hooks"])
		}
	}
}

type harnessRegistryFile struct {
	SchemaVersion string             `json:"schemaVersion"`
	Harnesses     []harnessFileEntry `json:"harnesses"`
}

type harnessFileEntry struct {
	ID            string                 `json:"id"`
	Version       string                 `json:"version"`
	Languages     []string               `json:"languages"`
	Toolchains    []string               `json:"toolchains"`
	Parameters    []harnessFileParameter `json:"parameters"`
	Hooks         []harnessFileHook      `json:"hooks"`
	ArtifactGlobs []string               `json:"artifactGlobs"`
}

type harnessFileParameter struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Default  any    `json:"default"`
}

type harnessFileHook struct {
	Phase    string `json:"phase"`
	HookID   string `json:"hookId"`
	Required bool   `json:"required"`
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("repository root with go.mod not found")
		}
		dir = parent
	}
}

func assertIssueCodes(t *testing.T, issues []DocumentationIssue, expected ...string) {
	t.Helper()
	set := map[string]struct{}{}
	for _, issue := range issues {
		set[issue.Code] = struct{}{}
	}
	for _, code := range expected {
		if _, ok := set[code]; !ok {
			t.Fatalf("missing issue code %q in %+v", code, issues)
		}
	}
}

func validDocumentationInput() DocumentationCheckInput {
	lock := LockedTemplateA()
	return DocumentationCheckInput{
		ChainPolicy: DocumentationIfAffected,
		Targets: []DocumentationTarget{
			{ID: "api-source", Class: TargetClassSource, Tags: []string{"public-api"}},
			{ID: "cli-docs", Class: TargetClassDocumentation},
			{ID: "cli-help-generated", Class: TargetClassGenerated, GeneratedFrom: []string{"api-source"}, GeneratorHook: "generate-cli-help"},
		},
		ImpactRules: []DocumentationImpactRule{
			{
				ID:             "cli-contract-docs",
				WhenTargetTags: []string{"public-api"},
				RequireTargets: []string{"cli-docs", "cli-help-generated"},
				RequiredHooks:  []string{"docs-check", "generated-freshness-check"},
			},
		},
		ChangedTarget: []string{"api-source"},
		UpdatedTarget: []string{"api-source", "cli-docs", "cli-help-generated"},
		PassedHooks:   []string{"docs-check", "generated-freshness-check"},
		Generated: []GeneratedArtifactState{
			{
				TargetID:     "cli-help-generated",
				TemplateID:   TemplateAID,
				Version:      TemplateAVersion,
				TemplateHash: lock.Hash,
				GeneratedBy:  evidence.Actor{Type: "agent", ID: "generator-bot"},
				ApprovedBy:   evidence.Actor{Type: "operator", ID: "reviewer-one"},
				Fresh:        true,
			},
		},
	}
}
