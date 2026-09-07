package console

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/taskdef"
)

type PreviewUnknown struct {
	UnknownID  string `json:"unknownId"`
	Reason     string `json:"reason"`
	NextAction string `json:"nextAction"`
}

type PreviewCallCounters struct {
	CommandCalls    int `json:"commandCalls"`
	NetworkCalls    int `json:"networkCalls"`
	ModelCalls      int `json:"modelCalls"`
	CredentialReads int `json:"credentialReads"`
	VersionProbes   int `json:"versionProbes"`
}

type PreviewReport struct {
	ConfigPath    string                    `json:"configPath"`
	RunnableInCLI bool                      `json:"runnableInCli"`
	PreviewRecord taskdef.PlanPreviewRecord `json:"previewRecord"`
	Unknowns      []PreviewUnknown          `json:"unknowns"`
	Warnings      []string                  `json:"warnings"`
	CallCounters  PreviewCallCounters       `json:"callCounters"`
}

func (cli CLI) executePreview(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("preview")
	chainPath := set.String("chain", "", "path to chain config file")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "preview", *jsonOutput, parseErr.String(), "usage: prfrail preview [--chain <path>] [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "preview", *jsonOutput, "preview does not accept positional arguments", "usage: prfrail preview [--chain <path>] [--json]")
	}

	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "preview", *jsonOutput, err)
	}
	resolvedPath, err := ResolveChainConfigPath(*chainPath, cwd)
	if err != nil {
		return writeCommandError(stdout, stderr, "preview", *jsonOutput, err)
	}
	cfg, rawConfig, err := LoadChainConfig(resolvedPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "preview", *jsonOutput, err)
	}

	report, err := BuildPreviewReport(resolvedPath, cfg, rawConfig, cli.now())
	if err != nil {
		return writeCommandError(stdout, stderr, "preview", *jsonOutput, err)
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "preview", OK: true, ExitCode: exitSuccess, Data: report})
	}
	fmt.Fprintf(stdout, "previewHash: %s\n", report.PreviewRecord.PreviewHash)
	fmt.Fprintf(stdout, "config: %s\n", report.ConfigPath)
	fmt.Fprintf(stdout, "mode: %s\n", report.PreviewRecord.Preview.Mode)
	fmt.Fprintf(stdout, "outcome: %s\n", report.PreviewRecord.Preview.Outcome)
	fmt.Fprintf(stdout, "runnableInCli: %t\n", report.RunnableInCLI)
	fmt.Fprintf(stdout, "steps: %d\n", len(report.PreviewRecord.Preview.Steps))
	fmt.Fprintf(stdout, "callCounters: command=%d network=%d model=%d credential=%d versionProbe=%d\n", report.CallCounters.CommandCalls, report.CallCounters.NetworkCalls, report.CallCounters.ModelCalls, report.CallCounters.CredentialReads, report.CallCounters.VersionProbes)
	for _, unknown := range report.Unknowns {
		fmt.Fprintf(stdout, "unknown: %s reason=%s next=%s\n", unknown.UnknownID, unknown.Reason, unknown.NextAction)
	}
	for _, warning := range report.Warnings {
		fmt.Fprintf(stdout, "warning: %s\n", warning)
	}
	return exitSuccess
}

func BuildPreviewReport(configPath string, cfg ChainConfig, rawConfig []byte, now time.Time) (PreviewReport, error) {
	definition, err := cfg.ToDefinition()
	if err != nil {
		return PreviewReport{}, err
	}
	chainCanonical, err := evidence.EncodeCanonical(definition)
	if err != nil {
		return PreviewReport{}, err
	}
	workspaceCanonical, err := evidence.EncodeCanonical(cfg.Workspace)
	if err != nil {
		return PreviewReport{}, err
	}
	policyCanonical, err := evidence.EncodeCanonical(previewPolicySnapshot(cfg))
	if err != nil {
		return PreviewReport{}, err
	}

	previewSteps, executableSteps := buildPreviewSteps(cfg)
	outcome := "ready"
	blockingEvidence := []string{}
	warnings := []string{}
	if len(executableSteps) > 0 {
		outcome = "blocked"
		blockingEvidence = make([]string, 0, len(executableSteps))
		for _, step := range executableSteps {
			blockingEvidence = append(blockingEvidence, evidence.Digest("", []byte(step.TaskID+"\n"+step.StepID+"\n"+step.Kind)))
		}
		warnings = append(warnings, "preview is read-only; executable steps require separate authorization and runtime capability probes")
	}

	sourceRef := cfg.Workspace.Components[0].ID
	preview := taskdef.PlanPreview{
		PreviewID:               NormalizeChainID(cfg.Chain.ID + "-preview"),
		CreatedAt:               now.UTC().Format(timestampLayout),
		Mode:                    "offline-read-only",
		ChainDefinitionHash:     evidence.Digest("", chainCanonical),
		WorkspaceDefinitionHash: evidence.Digest("", workspaceCanonical),
		EffectiveConfigHash:     evidence.Digest("", rawConfig),
		PolicyHash:              evidence.Digest("", policyCanonical),
		Locations: taskdef.PlanPreviewLocations{
			SourceRef: sourceRef,
			RunRef:    "run",
			StoreRef:  "store",
		},
		ReadTargetIDs:       buildReadTargetIDs(),
		WriteTargetIDs:      []string{},
		NetworkRequirements: []taskdef.PlanPreviewNetworkRequirement{},
		Steps:               previewSteps,
		ReviewPointIDs:      buildReviewPointIDs(cfg),
		Budget: taskdef.PlanPreviewBudget{
			ModelCallLimit:        0,
			TokenLimit:            0,
			WallClockMs:           60000,
			AttemptLimit:          1,
			Currency:              "USD",
			EstimatedAmountMicros: nil,
			EstimateStatus:        "not-applicable",
			PricingSource:         nil,
		},
		EstimatedStoreBytes: nil,
		ExclusionIDs:        []string{},
		Capabilities:        buildPreviewCapabilities(configPath, rawConfig),
		Outcome:             outcome,
		BlockingEvidence:    blockingEvidence,
	}
	previewRecord, err := taskdef.NewPlanPreviewRecord(preview)
	if err != nil {
		return PreviewReport{}, err
	}

	return PreviewReport{
		ConfigPath:    configPath,
		RunnableInCLI: len(executableSteps) == 0,
		PreviewRecord: previewRecord,
		Unknowns:      buildPreviewUnknowns(executableSteps),
		Warnings:      warnings,
		CallCounters: PreviewCallCounters{
			CommandCalls:    0,
			NetworkCalls:    0,
			ModelCalls:      0,
			CredentialReads: 0,
			VersionProbes:   0,
		},
	}, nil
}

func buildReadTargetIDs() []string {
	return []string{"chain-config", "source"}
}

func buildPreviewSteps(cfg ChainConfig) ([]taskdef.PlanPreviewStep, []ExecutableStep) {
	steps := make([]taskdef.PlanPreviewStep, 0, StepCount(cfg))
	executable := make([]ExecutableStep, 0)
	for _, task := range cfg.Tasks {
		for _, step := range task.Steps {
			kind := previewStepKind(step)
			gates := uniqueSortedIDs(step.Hooks)
			steps = append(steps, taskdef.PlanPreviewStep{
				TaskID:          task.ID,
				StepID:          step.ID,
				Kind:            kind,
				BlockingGateIDs: gates,
				ReviewRequired:  true,
			})
			if kind != "noop" {
				executable = append(executable, ExecutableStep{TaskID: task.ID, StepID: step.ID, Kind: kind})
			}
		}
	}
	return steps, executable
}

func previewStepKind(step StepConfig) string {
	if step.Kind == "code" {
		execution := step.Execution
		if execution == "" {
			execution = "autonomous"
		}
		if execution == "manual-handoff" {
			return "manual-handoff"
		}
		return "code"
	}
	return step.Kind
}

func buildReviewPointIDs(cfg ChainConfig) []string {
	reviewPoints := make([]string, 0, len(cfg.Tasks))
	for _, task := range cfg.Tasks {
		reviewPoints = append(reviewPoints, task.ID)
	}
	return uniqueSortedIDs(reviewPoints)
}

func buildPreviewCapabilities(configPath string, rawConfig []byte) []taskdef.PlanPreviewCapability {
	staticEvidence := evidence.Digest("", []byte("preview-read-only\n"+configPath+"\n"+evidence.Digest("", rawConfig)))
	return []taskdef.PlanPreviewCapability{
		{CapabilityID: "static-preview", Status: "verified", Evidence: []string{staticEvidence}},
		{CapabilityID: "command-execution", Status: "unknown", Evidence: []string{}},
		{CapabilityID: "network-access", Status: "unknown", Evidence: []string{}},
		{CapabilityID: "credential-input", Status: "unknown", Evidence: []string{}},
		{CapabilityID: "version-probe", Status: "unknown", Evidence: []string{}},
	}
}

func buildPreviewUnknowns(executable []ExecutableStep) []PreviewUnknown {
	unknowns := []PreviewUnknown{
		{UnknownID: "version-probe", Reason: "preview does not run tool/version probes", NextAction: "authorize validate/run for observed tool versions"},
		{UnknownID: "network-access", Reason: "preview does not perform network checks", NextAction: "authorize controlled preflight/network verification"},
		{UnknownID: "credential-input", Reason: "preview does not collect secret input", NextAction: "use manual-handoff secret-direct input in authorized execution"},
	}
	if len(executable) > 0 {
		unknowns = append(unknowns, PreviewUnknown{
			UnknownID:  "command-execution",
			Reason:     "preview never runs steps or hooks",
			NextAction: "authorize runtime execution after reviewing blocking steps",
		})
	}
	return unknowns
}

type previewPolicyTask struct {
	TaskID              string `json:"taskId"`
	Review              string `json:"review"`
	DocumentationPolicy string `json:"documentationPolicy"`
}

type previewPolicy struct {
	ChainID            string              `json:"chainId"`
	ChainProfile       string              `json:"chainProfile"`
	ChainDocumentation string              `json:"chainDocumentationPolicy"`
	Tasks              []previewPolicyTask `json:"tasks"`
}

func previewPolicySnapshot(cfg ChainConfig) previewPolicy {
	tasks := make([]previewPolicyTask, 0, len(cfg.Tasks))
	for index, task := range cfg.Tasks {
		policy, _, _ := resolveTaskDocumentationPolicy(cfg, index)
		tasks = append(tasks, previewPolicyTask{
			TaskID:              task.ID,
			Review:              task.Review,
			DocumentationPolicy: policy,
		})
	}
	return previewPolicy{
		ChainID:            cfg.Chain.ID,
		ChainProfile:       cfg.Chain.Profile,
		ChainDocumentation: effectiveChainDocumentationPolicy(cfg),
		Tasks:              tasks,
	}
}

func uniqueSortedIDs(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	set := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		set[value] = struct{}{}
	}
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
