package release

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func workflowFixture(t *testing.T) ([]byte, workflowDocument) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	workflow, err := decodeWorkflow(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateWorkflowContract(workflow); err != nil {
		t.Fatalf("baseline workflow: %v", err)
	}
	return data, workflow
}

func rejectWorkflowMutation(t *testing.T, data []byte, baseline workflowDocument, mutate func(*workflowDocument)) {
	t.Helper()
	mutated, err := decodeWorkflow(data)
	if err != nil {
		t.Fatal(err)
	}
	mutate(&mutated)
	if reflect.DeepEqual(baseline, mutated) {
		t.Fatal("mutation did not change the workflow")
	}
	if err := validateWorkflowContract(mutated); err == nil {
		t.Fatal("mutated workflow was accepted")
	}
}

func TestWorkflowScriptMutationsAreRejected(t *testing.T) {
	data, baseline := workflowFixture(t)
	mutations := map[string]func(*workflowStep){
		"quoted-commands":       func(step *workflowStep) { step.Run = "$ignored = @'\n" + step.Run + "\n'@\n" },
		"uninvoked-scriptblock": func(step *workflowStep) { step.Run = "$ignored = {\n" + step.Run + "\n}\n" },
		"unreachable-commands":  func(step *workflowStep) { step.Run = "if ($false) {\n" + step.Run + "\n}\n" },
		"early-success":         func(step *workflowStep) { step.Run = "exit 0\n" + step.Run },
		"fake-oracle-write": func(step *workflowStep) {
			step.Run += "\nSet-Content -LiteralPath (Join-Path $env:GITHUB_WORKSPACE 'generations/candidate/actual.json') -Value '{}'\n"
		},
		"alternate-variable": func(step *workflowStep) {
			step.Run = "$replacementTool = 'replacement'\n" + strings.ReplaceAll(step.Run, "& $releaseTool", "& $replacementTool")
		},
		"shadow-command":          func(step *workflowStep) { step.Run = "function global:go { }\n" + step.Run },
		"empty-script":            func(step *workflowStep) { step.Run = "" },
		"skip-step":               func(step *workflowStep) { step.If = "${{ false }}" },
		"different-shell":         func(step *workflowStep) { step.Shell = "bash" },
		"different-directory":     func(step *workflowStep) { step.WorkingDirectory = "candidate-source" },
		"extra-action-parameters": func(step *workflowStep) { step.With = map[string]any{"unreviewed": true} },
	}
	for jobName, job := range baseline.Jobs {
		for index, step := range job.Steps {
			if step.Run == "" {
				continue
			}
			for mutationName, mutate := range mutations {
				t.Run(jobName+"/"+step.Name+"/"+mutationName, func(t *testing.T) {
					rejectWorkflowMutation(t, data, baseline, func(workflow *workflowDocument) {
						mutate(&workflow.Jobs[jobName].Steps[index])
					})
				})
			}
		}
	}
}

func TestWorkflowActionParametersAreLocked(t *testing.T) {
	data, baseline := workflowFixture(t)
	for jobName, job := range baseline.Jobs {
		for index, step := range job.Steps {
			if step.Uses == "" {
				continue
			}
			for parameter := range step.With {
				for _, mutation := range []string{"delete", "change"} {
					t.Run(jobName+"/"+step.Name+"/"+parameter+"/"+mutation, func(t *testing.T) {
						rejectWorkflowMutation(t, data, baseline, func(workflow *workflowDocument) {
							parameters := workflow.Jobs[jobName].Steps[index].With
							if mutation == "delete" {
								delete(parameters, parameter)
							} else if value, ok := parameters[parameter].(bool); ok {
								parameters[parameter] = !value
							} else {
								parameters[parameter] = "unreviewed"
							}
						})
					})
				}
			}
		}
	}
}

func TestWorkflowStructureMutationsAreRejected(t *testing.T) {
	data, baseline := workflowFixture(t)
	globalMutations := map[string]func(*workflowDocument){
		"write-permission":    func(workflow *workflowDocument) { workflow.Permissions["contents"] = "write" },
		"extra-permission":    func(workflow *workflowDocument) { workflow.Permissions["id-token"] = "write" },
		"floating-toolchain":  func(workflow *workflowDocument) { workflow.Env["GO_VERSION"] = "stable" },
		"automatic-toolchain": func(workflow *workflowDocument) { workflow.Env["GOTOOLCHAIN"] = "auto" },
		"cgo-enabled":         func(workflow *workflowDocument) { workflow.Env["CGO_ENABLED"] = "1" },
		"extra-environment":   func(workflow *workflowDocument) { workflow.Env["GOFLAGS"] = "-mod=mod" },
		"wrong-push-branch":   func(workflow *workflowDocument) { workflow.On.Push.Branches = []string{"unreviewed"} },
		"missing-job":         func(workflow *workflowDocument) { delete(workflow.Jobs, "selfhost") },
		"extra-job":           func(workflow *workflowDocument) { workflow.Jobs["unreviewed"] = workflow.Jobs["test"] },
		"missing-input":       func(workflow *workflowDocument) { delete(workflow.On.WorkflowDispatch.Inputs, "candidate_commit") },
	}
	for name, mutate := range globalMutations {
		t.Run(name, func(t *testing.T) { rejectWorkflowMutation(t, data, baseline, mutate) })
	}
	for inputName := range baseline.On.WorkflowDispatch.Inputs {
		for _, mutation := range []string{"optional", "default", "type"} {
			t.Run("inputs/"+inputName+"/"+mutation, func(t *testing.T) {
				rejectWorkflowMutation(t, data, baseline, func(workflow *workflowDocument) {
					input := workflow.On.WorkflowDispatch.Inputs[inputName]
					switch mutation {
					case "optional":
						input.Required = false
					case "default":
						input.Default = "unreviewed"
					case "type":
						input.Type = "boolean"
					}
					workflow.On.WorkflowDispatch.Inputs[inputName] = input
				})
			})
		}
	}
	jobMutations := map[string]func(*workflowJob){
		"skip-job":           func(job *workflowJob) { job.If = "${{ false }}" },
		"wrong-dependency":   func(job *workflowJob) { job.Needs = "unreviewed" },
		"self-hosted-runner": func(job *workflowJob) { job.RunsOn = "self-hosted" },
		"one-platform":       func(job *workflowJob) { job.Strategy.Matrix["os"] = []any{"windows-latest"} },
		"extra-matrix-field": func(job *workflowJob) { job.Strategy.Matrix["include"] = []any{map[string]any{"os": "self-hosted"}} },
		"fail-fast": func(job *workflowJob) {
			enabled := true
			job.Strategy.FailFast = &enabled
		},
		"wrong-environment": func(job *workflowJob) { job.Env = map[string]string{"CANDIDATE_COMMIT": "unreviewed"} },
		"extra-step":        func(job *workflowJob) { job.Steps = append(job.Steps, workflowStep{Run: "Write-Output 'replacement'"}) },
		"missing-step":      func(job *workflowJob) { job.Steps = job.Steps[1:] },
		"duplicate-name":    func(job *workflowJob) { job.Steps[1].Name = job.Steps[0].Name },
		"unnamed-step":      func(job *workflowJob) { job.Steps[0].Name = "" },
		"step-reordering":   func(job *workflowJob) { job.Steps[0], job.Steps[1] = job.Steps[1], job.Steps[0] },
	}
	for jobName := range baseline.Jobs {
		for mutationName, mutate := range jobMutations {
			t.Run(jobName+"/"+mutationName, func(t *testing.T) {
				rejectWorkflowMutation(t, data, baseline, func(workflow *workflowDocument) {
					job := workflow.Jobs[jobName]
					mutate(&job)
					workflow.Jobs[jobName] = job
				})
			})
		}
	}
}

func TestWorkflowYAMLMutationsAreRejected(t *testing.T) {
	data, _ := workflowFixture(t)
	source := string(data)
	const setupAction = "        uses: actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e # v7.0.0\n"
	mutations := map[string]string{
		"missing-fail-fast":          strings.ReplaceAll(source, "      fail-fast: false\n", ""),
		"empty-run-with-action":      strings.Replace(source, setupAction, setupAction+"        run: ''\n", 1),
		"null-run-with-action":       strings.Replace(source, setupAction, setupAction+"        run: null\n", 1),
		"empty-shell-with-action":    strings.Replace(source, setupAction, setupAction+"        shell: ''\n", 1),
		"null-directory-with-action": strings.Replace(source, setupAction, setupAction+"        working-directory: null\n", 1),
		"null-step-condition":        strings.Replace(source, setupAction, setupAction+"        if: null\n", 1),
		"null-input-default":         strings.Replace(source, "        required: true\n", "        default: null\n        required: true\n", 1),
		"duplicate-step-field":       strings.Replace(source, setupAction, setupAction+setupAction, 1),
		"duplicate-job":              strings.Replace(source, "  test:\n", "  test: {}\n  test:\n", 1),
		"unknown-field":              source + "unreviewed: true\n",
		"trailing-document":          source + "\n---\nname: unreviewed\n",
		"empty-trailing-document":    source + "\n---\n",
		"missing-pull-request":       strings.Replace(source, "  pull_request:\n", "", 1),
		"filtered-pull-request":      strings.Replace(source, "  pull_request:\n", "  pull_request:\n    branches: [main]\n", 1),
		"anchored-permissions":       strings.Replace(source, "permissions:\n", "permissions: &permissions\n", 1),
		"null-cgo":                   strings.Replace(source, "  CGO_ENABLED: \"0\"\n", "  CGO_ENABLED: null\n", 1),
	}
	for name, changed := range mutations {
		t.Run(name, func(t *testing.T) {
			if changed == source {
				t.Fatal("YAML mutation did not change the workflow")
			}
			workflow, err := decodeWorkflow([]byte(changed))
			if err == nil {
				err = validateWorkflowContract(workflow)
			}
			if err == nil {
				t.Fatal("mutated YAML was accepted")
			}
		})
	}
}
