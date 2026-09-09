package release

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type workflowDocument struct {
	Name string `yaml:"name"`
	On   struct {
		Push struct {
			Branches []string `yaml:"branches"`
		} `yaml:"push"`
		PullRequest      any `yaml:"pull_request"`
		WorkflowDispatch struct {
			Inputs map[string]workflowInput `yaml:"inputs"`
		} `yaml:"workflow_dispatch"`
	} `yaml:"on"`
	Permissions map[string]string      `yaml:"permissions"`
	Env         map[string]string      `yaml:"env"`
	Jobs        map[string]workflowJob `yaml:"jobs"`
}

type workflowInput struct {
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Type        string `yaml:"type"`
	Default     any    `yaml:"default"`
}

type workflowJob struct {
	Name     string `yaml:"name"`
	If       string `yaml:"if"`
	Needs    string `yaml:"needs"`
	Strategy struct {
		FailFast *bool          `yaml:"fail-fast"`
		Matrix   map[string]any `yaml:"matrix"`
	} `yaml:"strategy"`
	RunsOn string            `yaml:"runs-on"`
	Env    map[string]string `yaml:"env"`
	Steps  []workflowStep    `yaml:"steps"`
}

type workflowStep struct {
	Name             string         `yaml:"name"`
	Uses             string         `yaml:"uses"`
	If               string         `yaml:"if"`
	Shell            string         `yaml:"shell"`
	WorkingDirectory string         `yaml:"working-directory"`
	With             map[string]any `yaml:"with"`
	Run              string         `yaml:"run"`
}

func TestSelfhostWorkflowUsesIndependentTrustedVerifier(t *testing.T) {
	workflowPath := filepath.Join("..", "..", ".github", "workflows", "ci.yml")
	data, err := os.ReadFile(workflowPath)
	if err != nil {
		t.Fatal(err)
	}
	workflow, err := decodeWorkflow(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateWorkflowContract(workflow); err != nil {
		t.Fatal(err)
	}
}

func decodeWorkflow(data []byte) (workflowDocument, error) {
	var workflow workflowDocument
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&workflow); err != nil {
		return workflowDocument{}, err
	}
	if err := decoder.Decode(&yaml.Node{}); !errors.Is(err, io.EOF) {
		return workflowDocument{}, fmt.Errorf("workflow must contain exactly one YAML document")
	}
	var document yaml.Node
	if err := yaml.Unmarshal(data, &document); err != nil {
		return workflowDocument{}, err
	}
	if err := validateWorkflowNode(&document, ""); err != nil {
		return workflowDocument{}, err
	}
	return workflow, nil
}

func validateWorkflowNode(node *yaml.Node, path string) error {
	if node.Kind == yaml.AliasNode || node.Anchor != "" || node.Tag == "!!merge" {
		return fmt.Errorf("workflow aliases, anchors and merges are not allowed at %q", path)
	}
	if node.Tag == "!!null" && path != "on/pull_request" {
		return fmt.Errorf("explicit null is not allowed at %q", path)
	}
	if node.Kind != yaml.MappingNode {
		for _, child := range node.Content {
			if err := validateWorkflowNode(child, path); err != nil {
				return err
			}
		}
		return nil
	}
	fields := make(map[string]*yaml.Node, len(node.Content)/2)
	for index := 0; index < len(node.Content); index += 2 {
		fields[node.Content[index].Value] = node.Content[index+1]
	}
	if path == "on" && (fields["push"] == nil || fields["pull_request"] == nil || fields["workflow_dispatch"] == nil) {
		return fmt.Errorf("workflow must declare push, pull_request and workflow_dispatch")
	}
	if fields["default"] != nil {
		return fmt.Errorf("workflow input defaults are not allowed at %q", path)
	}
	if fields["uses"] != nil {
		for _, field := range []string{"run", "shell", "working-directory"} {
			if fields[field] != nil {
				return fmt.Errorf("Action at %q must not declare %q, even when empty", path, field)
			}
		}
	}
	for index := 0; index < len(node.Content); index += 2 {
		key := node.Content[index]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
			return fmt.Errorf("workflow mapping key at %q must be a string", path)
		}
		childPath := key.Value
		if path != "" {
			childPath = path + "/" + key.Value
		}
		if err := validateWorkflowNode(node.Content[index+1], childPath); err != nil {
			return err
		}
	}
	return nil
}

func validateWorkflowContract(workflow workflowDocument) error {
	if err := validateWorkflowStructure(workflow); err != nil {
		return err
	}
	if err := validateWorkflowActions(workflow); err != nil {
		return err
	}
	return validateWorkflowScripts(workflow)
}

func validateWorkflowStructure(workflow workflowDocument) error {
	if !reflect.DeepEqual(workflow.Permissions, map[string]string{"contents": "read"}) || !reflect.DeepEqual(workflow.Env, map[string]string{"GO_VERSION": "1.22.12", "CGO_ENABLED": "0", "GOTOOLCHAIN": "local"}) {
		return fmt.Errorf("workflow permissions or toolchain environment differ from the reviewed contract")
	}
	if !reflect.DeepEqual(workflow.On.Push.Branches, []string{"main"}) || workflow.On.PullRequest != nil {
		return fmt.Errorf("workflow push or pull_request triggers differ from the reviewed contract")
	}
	expectedInputs := []string{
		"candidate_commit", "verifier_commit", "bootstrap_commit", "bootstrap_manifest_sha256",
		"expected_oracle_sha256", "license_policy_sha256", "noop_chain_sha256", "executable_chain_sha256",
	}
	inputs := workflow.On.WorkflowDispatch.Inputs
	if len(inputs) != len(expectedInputs) {
		return fmt.Errorf("workflow_dispatch inputs=%d, want %d", len(inputs), len(expectedInputs))
	}
	for _, inputName := range expectedInputs {
		input, exists := inputs[inputName]
		if !exists || !input.Required || input.Type != "string" || input.Default != nil {
			return fmt.Errorf("input %q must be a required string without a default", inputName)
		}
	}
	const dispatchOnMain = "github.event_name == 'workflow_dispatch' && github.ref == 'refs/heads/main'"
	type jobContract struct {
		condition   string
		needs       string
		environment map[string]string
		steps       []string
	}
	expectedJobs := map[string]jobContract{
		"test": {
			steps: []string{"Checkout source", "Setup Go", "Build", "Vet", "Test", "Contract fixtures"},
		},
		"candidate-build": {
			condition:   dispatchOnMain,
			environment: map[string]string{"CANDIDATE_COMMIT": "${{ inputs.candidate_commit }}"},
			steps:       []string{"Validate candidate pin", "Checkout candidate", "Setup Go", "Build candidate without executing it", "Upload immutable candidate binary"},
		},
		"candidate-probe": {
			condition: dispatchOnMain, needs: "candidate-build",
			environment: map[string]string{
				"CANDIDATE_COMMIT": "${{ inputs.candidate_commit }}", "VERIFIER_COMMIT": "${{ inputs.verifier_commit }}",
				"NOOP_CHAIN_SHA256": "${{ inputs.noop_chain_sha256 }}", "EXECUTABLE_CHAIN_SHA256": "${{ inputs.executable_chain_sha256 }}",
			},
			steps: []string{"Validate repository trust pins", "Checkout trusted verifier", "Setup Go", "Build trusted probe harness", "Download immutable candidate binary", "Probe untrusted candidate", "Upload untrusted probe oracle"},
		},
		"selfhost": {
			condition: dispatchOnMain, needs: "candidate-probe",
			environment: map[string]string{
				"CANDIDATE_COMMIT": "${{ inputs.candidate_commit }}", "VERIFIER_COMMIT": "${{ inputs.verifier_commit }}", "BOOTSTRAP_COMMIT": "${{ inputs.bootstrap_commit }}",
				"BOOTSTRAP_MANIFEST_SHA256": "${{ inputs.bootstrap_manifest_sha256 }}", "EXPECTED_ORACLE_SHA256": "${{ inputs.expected_oracle_sha256 }}", "LICENSE_POLICY_SHA256": "${{ inputs.license_policy_sha256 }}",
			},
			steps: []string{"Validate repository trust pins", "Checkout trusted verifier", "Checkout fixed bootstrap seed", "Checkout candidate metadata", "Setup Go", "Build trusted verifier and seed", "Download immutable candidate binary", "Download untrusted probe oracle", "Verify evidence and generate release metadata", "Upload verified CI artifact"},
		},
	}
	if len(workflow.Jobs) != len(expectedJobs) {
		return fmt.Errorf("workflow jobs=%d, want %d", len(workflow.Jobs), len(expectedJobs))
	}
	for name, expected := range expectedJobs {
		job, exists := workflow.Jobs[name]
		if !exists || job.If != expected.condition || job.Needs != expected.needs || !reflect.DeepEqual(job.Env, expected.environment) {
			return fmt.Errorf("job %q conditions, dependencies or environment differ from the reviewed contract", name)
		}
		if job.RunsOn != "${{ matrix.os }}" || job.Strategy.FailFast == nil || *job.Strategy.FailFast || !reflect.DeepEqual(job.Strategy.Matrix, map[string]any{"os": []any{"ubuntu-latest", "windows-latest"}}) {
			return fmt.Errorf("job %q must use the reviewed native runner matrix", name)
		}
		if len(job.Steps) != len(expected.steps) {
			return fmt.Errorf("job %q steps=%d, want %d", name, len(job.Steps), len(expected.steps))
		}
		for index, expectedName := range expected.steps {
			if job.Steps[index].Name != expectedName {
				return fmt.Errorf("job %q step %d=%q, want %q", name, index, job.Steps[index].Name, expectedName)
			}
		}
	}
	return nil
}

func validateWorkflowScripts(workflow workflowDocument) error {
	type scriptContract struct {
		digest           string
		shell            string
		condition        string
		workingDirectory string
	}
	expectedScripts := map[string]scriptContract{
		"test/Build": {digest: "854004a0c269ca4993dad170dc3f7d6e996493547b7c6b5c5c83cf20d113b42a"},
		"test/Vet":   {digest: "3d4b424f6fcd11f530525ddea9a46aeaffec7c883b6c04b6baf7df36bf47800b"},
		"test/Test":  {digest: "a8496b1836c1e6e44f6f4dfa908c3942ddebce2e0be73e5411aba70c6471789a"},
		"test/Contract fixtures": {
			digest: "94577fb7f37f9fbfcef7207756eebf7999694fea78cfe21c23b7e904d2ce7d5d", condition: "runner.os == 'Linux'", workingDirectory: "tools/contracts",
		},
		"candidate-build/Validate candidate pin": {
			digest: "5b4930d6ebbea82cb11be937304ae3cddb7bfb77236c276596451c9ba711b216", shell: "pwsh",
		},
		"candidate-build/Build candidate without executing it": {
			digest: "a5fda73329c34f4d3dafb939421f736281104b75319e8a44b4b0c5a4ff3c1a56", shell: "pwsh",
		},
		"candidate-probe/Validate repository trust pins": {
			digest: "17b4bbcd42ebc8ca74cc3124208f277b74245d1cfb3cc9ba4936d0af422fba39", shell: "pwsh",
		},
		"candidate-probe/Build trusted probe harness": {
			digest: "03c12f0eda7f57745960056efeb9fffe801efe0c3d685094bb6e82958237d3c6", shell: "pwsh", workingDirectory: "verifier-source",
		},
		"candidate-probe/Probe untrusted candidate": {
			digest: "61a15643bb3351c9aabdf2a4f461b56d2c084d91e0a86bf712889d25d12cda15", shell: "pwsh",
		},
		"selfhost/Validate repository trust pins": {
			digest: "6bf96d6f73fba9d8bd7acabfca0a4e7269d0ae8658a79cd3bf5a9db9e8cd4f48", shell: "pwsh",
		},
		"selfhost/Build trusted verifier and seed": {
			digest: "27a4f6a51f91cc35fdcab8f3b7dc1f509a28d0d9fff04c038a45131d1df70931", shell: "pwsh",
		},
		"selfhost/Verify evidence and generate release metadata": {
			digest: "7989dd9584abe48bcea870609d6f6d2df75399db9ac70ceded92636feab5a664", shell: "pwsh",
		},
	}
	seen := make(map[string]bool, len(expectedScripts))
	for jobName, job := range workflow.Jobs {
		for _, step := range job.Steps {
			key := jobName + "/" + step.Name
			expected, exists := expectedScripts[key]
			if !exists {
				if step.Uses == "" {
					return fmt.Errorf("unexpected script step %q", key)
				}
				continue
			}
			if seen[key] {
				return fmt.Errorf("duplicate script step %q", key)
			}
			seen[key] = true
			if step.Uses != "" || step.With != nil || step.Shell != expected.shell || step.If != expected.condition || step.WorkingDirectory != expected.workingDirectory {
				return fmt.Errorf("script step %q execution fields differ from the reviewed contract", key)
			}
			digest := fmt.Sprintf("%x", sha256.Sum256([]byte(step.Run)))
			if digest != expected.digest {
				return fmt.Errorf("script step %q differs from the reviewed body: SHA-256=%s, want %s", key, digest, expected.digest)
			}
		}
	}
	if len(seen) != len(expectedScripts) {
		return fmt.Errorf("workflow has %d reviewed scripts, want %d", len(seen), len(expectedScripts))
	}
	return nil
}

func validateWorkflowActions(workflow workflowDocument) error {
	expectedActions := map[string]string{
		"Checkout source":                     "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1",
		"Checkout candidate":                  "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1",
		"Checkout trusted verifier":           "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1",
		"Checkout fixed bootstrap seed":       "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1",
		"Checkout candidate metadata":         "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1",
		"Setup Go":                            "actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e",
		"Upload immutable candidate binary":   "actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a",
		"Download immutable candidate binary": "actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c",
		"Upload untrusted probe oracle":       "actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a",
		"Download untrusted probe oracle":     "actions/download-artifact@3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c",
		"Upload verified CI artifact":         "actions/upload-artifact@043fb46d1a93c77aae656e7c1c64a875d1fc6a0a",
	}
	expectedParameters := map[string]map[string]any{
		"test/Checkout source": {"persist-credentials": false},
		"test/Setup Go":        {"go-version": "${{ env.GO_VERSION }}", "cache": true},
		"candidate-build/Checkout candidate": {
			"ref": "${{ env.CANDIDATE_COMMIT }}", "path": "candidate-source", "persist-credentials": false,
		},
		"candidate-build/Setup Go": {"go-version": "${{ env.GO_VERSION }}", "cache": false},
		"candidate-build/Upload immutable candidate binary": {
			"name": "candidate-binary-${{ runner.os }}-${{ inputs.candidate_commit }}", "path": "generations/candidate/release/prfrail*", "if-no-files-found": "error",
		},
		"candidate-probe/Checkout trusted verifier": {
			"ref": "${{ env.VERIFIER_COMMIT }}", "path": "verifier-source", "persist-credentials": false,
		},
		"candidate-probe/Setup Go": {"go-version": "${{ env.GO_VERSION }}", "cache": false},
		"candidate-probe/Download immutable candidate binary": {
			"name": "candidate-binary-${{ runner.os }}-${{ inputs.candidate_commit }}", "path": "generations/candidate/release", "digest-mismatch": "error",
		},
		"candidate-probe/Upload untrusted probe oracle": {
			"name": "candidate-oracle-${{ runner.os }}-${{ inputs.candidate_commit }}", "path": "generations/candidate/actual.json", "if-no-files-found": "error",
		},
		"selfhost/Checkout trusted verifier": {
			"ref": "${{ env.VERIFIER_COMMIT }}", "path": "verifier-source", "persist-credentials": false,
		},
		"selfhost/Checkout fixed bootstrap seed": {
			"ref": "${{ env.BOOTSTRAP_COMMIT }}", "path": "seed-source", "persist-credentials": false,
		},
		"selfhost/Checkout candidate metadata": {
			"ref": "${{ env.CANDIDATE_COMMIT }}", "path": "candidate-source", "persist-credentials": false,
		},
		"selfhost/Setup Go": {"go-version": "${{ env.GO_VERSION }}", "cache": false},
		"selfhost/Download immutable candidate binary": {
			"name": "candidate-binary-${{ runner.os }}-${{ inputs.candidate_commit }}", "path": "generations/candidate/release", "digest-mismatch": "error",
		},
		"selfhost/Download untrusted probe oracle": {
			"name": "candidate-oracle-${{ runner.os }}-${{ inputs.candidate_commit }}", "path": "generations/candidate", "digest-mismatch": "error",
		},
		"selfhost/Upload verified CI artifact": {
			"name": "prfrail-${{ runner.os }}-${{ inputs.candidate_commit }}", "path": "generations/candidate/release/*", "if-no-files-found": "error",
		},
	}
	seen := make(map[string]bool, len(expectedParameters))
	for jobName, job := range workflow.Jobs {
		for _, step := range job.Steps {
			key := jobName + "/" + step.Name
			parameters, required := expectedParameters[key]
			if required {
				expected := expectedActions[step.Name]
				if seen[key] || step.Uses != expected || step.If != "" || step.Run != "" || step.Shell != "" || step.WorkingDirectory != "" {
					return fmt.Errorf("step %q must be one unconditional Action %q without run fields", key, expected)
				}
				seen[key] = true
				if !reflect.DeepEqual(step.With, parameters) {
					return fmt.Errorf("Action %q parameters differ from the reviewed contract: got %#v, want %#v", key, step.With, parameters)
				}
			} else if step.Uses != "" {
				return fmt.Errorf("job %s step %q has unapproved Action %q", jobName, step.Name, step.Uses)
			}
		}
	}
	if len(seen) != len(expectedParameters) {
		return fmt.Errorf("workflow has %d required Actions, want %d", len(seen), len(expectedParameters))
	}
	return nil
}

func TestWorkflowActionMutationsAreRejected(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := decodeWorkflow(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateWorkflowContract(baseline); err != nil {
		t.Fatal(err)
	}
	mutations := map[string]func(*workflowStep){
		"missing-uses": func(step *workflowStep) { step.Uses = "" },
		"run-replacement": func(step *workflowStep) {
			step.Uses = ""
			step.With = nil
			step.Shell = "pwsh"
			step.Run = "Write-Output 'replacement'"
		},
		"action-with-run": func(step *workflowStep) { step.Run = "Write-Output 'replacement'" },
		"foreign-action":  func(step *workflowStep) { step.Uses = "untrusted/action@" + strings.Repeat("a", 40) },
		"changed-sha": func(step *workflowStep) {
			step.Uses = strings.Split(step.Uses, "@")[0] + "@" + strings.Repeat("a", 40)
		},
		"role-swap": func(step *workflowStep) {
			if strings.HasPrefix(step.Uses, "actions/checkout@") {
				step.Uses = "actions/setup-go@b7ad1dad31e06c5925ef5d2fc7ad053ef454303e"
			} else {
				step.Uses = "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1"
			}
		},
		"unnamed-action":        func(step *workflowStep) { step.Name = "" },
		"unknown-action-name":   func(step *workflowStep) { step.Name = "Unexpected Action" },
		"conditional-action":    func(step *workflowStep) { step.If = "${{ false }}" },
		"action-with-shell":     func(step *workflowStep) { step.Shell = "pwsh" },
		"action-with-directory": func(step *workflowStep) { step.WorkingDirectory = "candidate-source" },
		"missing-parameters":    func(step *workflowStep) { step.With = nil },
		"extra-parameter":       func(step *workflowStep) { step.With["unreviewed"] = true },
	}
	actionCount := 0
	for jobName, job := range baseline.Jobs {
		for index, step := range job.Steps {
			if step.Uses == "" {
				continue
			}
			actionCount++
			for mutationName, mutate := range mutations {
				t.Run(jobName+"/"+step.Name+"/"+mutationName, func(t *testing.T) {
					mutated, err := decodeWorkflow(data)
					if err != nil {
						t.Fatal(err)
					}
					mutate(&mutated.Jobs[jobName].Steps[index])
					if err := validateWorkflowContract(mutated); err == nil {
						t.Fatal("mutated Action must be rejected")
					}
				})
			}
		}
	}
	if actionCount == 0 {
		t.Fatal("mutation suite must exercise existing Actions")
	}
}

func TestWorkflowScriptsRejectCommentedCommands(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	baseline, err := decodeWorkflow(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := validateWorkflowContract(baseline); err != nil {
		t.Fatal(err)
	}
	for jobName, job := range baseline.Jobs {
		for index, step := range job.Steps {
			if step.Run == "" {
				continue
			}
			t.Run(jobName+"/"+step.Name, func(t *testing.T) {
				mutated, err := decodeWorkflow(data)
				if err != nil {
					t.Fatal(err)
				}
				mutated.Jobs[jobName].Steps[index].Run = "<#\n" + step.Run + "\n#>\nWrite-Output 'replacement'\n"
				if err := validateWorkflowContract(mutated); err == nil {
					t.Fatalf("commented commands were accepted; baseline script SHA-256=%x", sha256.Sum256([]byte(step.Run)))
				}
			})
		}
	}
}
