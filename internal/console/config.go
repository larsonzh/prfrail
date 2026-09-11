package console

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	DefaultChainConfigName   = "proofrail.chain.json"
	FallbackChainConfigName  = "proofrail.json"
	defaultDocumentationRule = "if-affected"
)

var ErrInvalidConfig = errors.New("invalid proofrail configuration")

type ChainConfig struct {
	SchemaVersion string               `json:"schemaVersion"`
	Chain         ChainSection         `json:"chain"`
	AI            *AIConfig            `json:"ai,omitempty"`
	Documentation *DocumentationConfig `json:"documentation,omitempty"`
	Tasks         []TaskConfig         `json:"tasks"`
	Workspace     WorkspaceConfig      `json:"workspace"`
}

type AIConfig struct {
	Profiles []adapters.AIProviderProfile `json:"profiles"`
	Channels []AIChannelBinding           `json:"channels"`
}

type AIChannelBinding struct {
	Channel   string `json:"channel"`
	ProfileID string `json:"profileId"`
}

type ChainSection struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Profile string `json:"profile"`
}

type DocumentationConfig struct {
	Policy      string                    `json:"policy"`
	ImpactRules []DocumentationImpactRule `json:"impactRules,omitempty"`
}

type DocumentationImpactRule struct {
	ID             string   `json:"id"`
	WhenTargetTags []string `json:"whenTargetTags"`
	RequireTargets []string `json:"requireTargets"`
	RequiredHooks  []string `json:"requiredHooks"`
}

type TaskConfig struct {
	ID                  string       `json:"id"`
	Name                string       `json:"name,omitempty"`
	Model               string       `json:"model,omitempty"`
	Review              string       `json:"review"`
	DocumentationPolicy string       `json:"documentationPolicy,omitempty"`
	Steps               []StepConfig `json:"steps"`
}

type StepConfig struct {
	ID            string              `json:"id"`
	Kind          string              `json:"kind"`
	Reason        string              `json:"reason,omitempty"`
	Model         string              `json:"model,omitempty"`
	Execution     string              `json:"execution,omitempty"`
	Components    []string            `json:"components,omitempty"`
	LanguageScope []string            `json:"languageScopes,omitempty"`
	Targets       []string            `json:"targets,omitempty"`
	Hooks         []string            `json:"hooks,omitempty"`
	Order         string              `json:"order,omitempty"`
	HandoffPolicy *HandoffPolicyInput `json:"handoffPolicy,omitempty"`
}

type HandoffPolicyInput struct {
	AllowedTargets   []string `json:"allowedTargets"`
	InputPolicy      string   `json:"inputPolicy"`
	HandoffTimeoutMs int      `json:"handoffTimeoutMs"`
	ReturnActions    []string `json:"returnActions"`
	HooksAfterReturn []string `json:"hooksAfterReturn"`
}

type WorkspaceConfig struct {
	Agent      *AgentConfig      `json:"agent,omitempty"`
	Components []ComponentConfig `json:"components"`
}

type AgentConfig struct {
	Channel string `json:"channel"`
}

type ComponentConfig struct {
	ID             string                `json:"id"`
	Root           string                `json:"root"`
	LanguageScopes []LanguageScopeConfig `json:"languageScopes"`
	DependsOn      []string              `json:"dependsOn,omitempty"`
}

type LanguageScopeConfig struct {
	ID        string   `json:"id"`
	Language  string   `json:"language"`
	Harness   string   `json:"harness"`
	Toolchain string   `json:"toolchain,omitempty"`
	Targets   []string `json:"targets"`
	DependsOn []string `json:"dependsOn,omitempty"`
}

type ExecutableStep struct {
	TaskID string `json:"taskId"`
	StepID string `json:"stepId"`
	Kind   string `json:"kind"`
}

type ConfigValidationSummary struct {
	ConfigPath        string           `json:"configPath"`
	ChainID           string           `json:"chainId"`
	TaskCount         int              `json:"taskCount"`
	StepCount         int              `json:"stepCount"`
	RunnableInCLI     bool             `json:"runnableInCli"`
	ExecutableSteps   []ExecutableStep `json:"executableSteps"`
	Warnings          []string         `json:"warnings"`
	ResolvedChainPath string           `json:"resolvedChainPath,omitempty"`
}

type ConfigExplainReport struct {
	ConfigPath          string              `json:"configPath"`
	ConfigHash          string              `json:"configHash"`
	ChainID             string              `json:"chainId"`
	Profile             string              `json:"profile"`
	AIProfiles          []AIProfileExplain  `json:"aiProfiles"`
	AIChannels          []AIChannelExplain  `json:"aiChannels"`
	DocumentationPolicy string              `json:"documentationPolicy"`
	RunnableInCLI       bool                `json:"runnableInCli"`
	ExecutableSteps     []ExecutableStep    `json:"executableSteps"`
	TaskPolicies        []TaskPolicyExplain `json:"taskPolicies"`
	Resolution          []ResolutionExplain `json:"resolution"`
	ConfigSearchOrder   []string            `json:"configSearchOrder"`
}

type AIProfileExplain struct {
	ProfileID         string `json:"profileId"`
	ProfileConfigHash string `json:"profileConfigHash"`
}

type AIChannelExplain struct {
	Channel           string `json:"channel"`
	ProfileID         string `json:"profileId"`
	ProfileConfigHash string `json:"profileConfigHash"`
}

type TaskPolicyExplain struct {
	TaskID        string  `json:"taskId"`
	Policy        string  `json:"policy"`
	SourceKind    string  `json:"sourceKind"`
	SourceHash    *string `json:"sourceHash"`
	SourcePointer *string `json:"sourcePointer"`
}

type ResolutionExplain struct {
	EffectivePointer string  `json:"effectivePointer"`
	SourceKind       string  `json:"sourceKind"`
	SourceHash       *string `json:"sourceHash"`
	SourcePointer    *string `json:"sourcePointer"`
}

func DefaultChainConfig(chainID string) ChainConfig {
	if !evidence.ValidID(chainID) {
		chainID = "proofrail"
	}
	return ChainConfig{
		SchemaVersion: "1.0.0",
		Chain: ChainSection{
			ID:      chainID,
			Profile: "minimal",
		},
		Tasks: []TaskConfig{{
			ID:     "bootstrap-task",
			Review: "manual",
			Steps: []StepConfig{{
				ID:     "noop-bootstrap",
				Kind:   "noop",
				Reason: "Bootstrap template generated by prfrail init",
			}},
		}},
		Workspace: WorkspaceConfig{
			Components: []ComponentConfig{{
				ID:   "workspace-main",
				Root: ".",
				LanguageScopes: []LanguageScopeConfig{{
					ID:       "generic-main",
					Language: "generic",
					Harness:  "generic-standard",
					Targets:  []string{"source"},
				}},
			}},
		},
	}
}

func ResolveChainConfigPath(explicitPath, cwd string) (string, error) {
	if explicitPath != "" {
		if cwd == "" {
			return filepath.Abs(explicitPath)
		}
		return resolvePathFromCWD(cwd, explicitPath)
	}
	if cwd == "" {
		return "", fmt.Errorf("%w: working directory is required", ErrInvalidConfig)
	}
	candidates := []string{
		filepath.Join(cwd, DefaultChainConfigName),
		filepath.Join(cwd, FallbackChainConfigName),
	}
	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", err
		}
		if info.IsDir() {
			continue
		}
		return filepath.Abs(candidate)
	}
	return "", fmt.Errorf("%w: config not found; looked for %s then %s", ErrInvalidConfig, DefaultChainConfigName, FallbackChainConfigName)
}

func LoadChainConfig(path string) (ChainConfig, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ChainConfig{}, nil, err
	}
	var cfg ChainConfig
	if err := evidence.DecodeStrictJSON(data, &cfg); err != nil {
		return ChainConfig{}, nil, fmt.Errorf("%w: %v", ErrInvalidConfig, err)
	}
	if err := ValidateChainConfig(cfg); err != nil {
		return ChainConfig{}, nil, err
	}
	return cfg, data, nil
}

func ValidateChainConfig(cfg ChainConfig) error {
	if cfg.SchemaVersion != "1.0.0" {
		return fmt.Errorf("%w: unsupported schemaVersion %q", ErrInvalidConfig, cfg.SchemaVersion)
	}
	if !evidence.ValidID(cfg.Chain.ID) {
		return fmt.Errorf("%w: invalid chain.id", ErrInvalidConfig)
	}
	if cfg.Chain.Profile != "minimal" && cfg.Chain.Profile != "standard" && cfg.Chain.Profile != "strict" {
		return fmt.Errorf("%w: invalid chain.profile %q", ErrInvalidConfig, cfg.Chain.Profile)
	}
	if cfg.Documentation != nil {
		if err := validateDocumentationPolicy(cfg.Documentation.Policy); err != nil {
			return err
		}
	}
	if cfg.AI != nil {
		if err := validateAIConfig(cfg.AI); err != nil {
			return err
		}
	}
	if len(cfg.Tasks) == 0 {
		return fmt.Errorf("%w: tasks must not be empty", ErrInvalidConfig)
	}
	for _, task := range cfg.Tasks {
		if !evidence.ValidID(task.ID) {
			return fmt.Errorf("%w: invalid task id %q", ErrInvalidConfig, task.ID)
		}
		if task.Review != "manual" && task.Review != "policy" {
			return fmt.Errorf("%w: task %q has invalid review mode %q", ErrInvalidConfig, task.ID, task.Review)
		}
		if task.DocumentationPolicy != "" {
			if err := validateDocumentationPolicy(task.DocumentationPolicy); err != nil {
				return err
			}
		}
		if len(task.Steps) == 0 {
			return fmt.Errorf("%w: task %q has no steps", ErrInvalidConfig, task.ID)
		}
		for _, step := range task.Steps {
			if !evidence.ValidID(step.ID) {
				return fmt.Errorf("%w: invalid step id %q", ErrInvalidConfig, step.ID)
			}
			switch step.Kind {
			case "code":
				execution := step.Execution
				if execution == "" {
					execution = "autonomous"
				}
				if execution != "autonomous" && execution != "manual-handoff" {
					return fmt.Errorf("%w: step %q has invalid execution mode %q", ErrInvalidConfig, step.ID, step.Execution)
				}
				if execution == "manual-handoff" {
					if step.HandoffPolicy == nil {
						return fmt.Errorf("%w: manual-handoff step %q requires handoffPolicy", ErrInvalidConfig, step.ID)
					}
					if err := validateHandoffPolicy(step.HandoffPolicy); err != nil {
						return err
					}
				}
				if execution == "autonomous" && step.HandoffPolicy != nil {
					return fmt.Errorf("%w: autonomous step %q must not define handoffPolicy", ErrInvalidConfig, step.ID)
				}
			case "build", "verify":
				if len(step.Hooks) == 0 {
					return fmt.Errorf("%w: %s step %q requires hooks", ErrInvalidConfig, step.Kind, step.ID)
				}
				for _, hookID := range step.Hooks {
					if !evidence.ValidID(hookID) {
						return fmt.Errorf("%w: step %q has invalid hook id %q", ErrInvalidConfig, step.ID, hookID)
					}
				}
				if step.Execution != "" || step.HandoffPolicy != nil {
					return fmt.Errorf("%w: %s step %q must not define execution or handoffPolicy", ErrInvalidConfig, step.Kind, step.ID)
				}
			case "noop":
				if strings.TrimSpace(step.Reason) == "" {
					return fmt.Errorf("%w: noop step %q requires reason", ErrInvalidConfig, step.ID)
				}
				if step.Execution != "" || step.HandoffPolicy != nil || len(step.Hooks) > 0 {
					return fmt.Errorf("%w: noop step %q must not define execution, handoffPolicy, or hooks", ErrInvalidConfig, step.ID)
				}
			default:
				return fmt.Errorf("%w: step %q has unknown kind %q", ErrInvalidConfig, step.ID, step.Kind)
			}
		}
	}
	if cfg.Workspace.Agent != nil {
		if cfg.Workspace.Agent.Channel != "ipc" && cfg.Workspace.Agent.Channel != "file-queue" {
			return fmt.Errorf("%w: workspace.agent.channel must be ipc or file-queue", ErrInvalidConfig)
		}
	}
	if len(cfg.Workspace.Components) == 0 {
		return fmt.Errorf("%w: workspace.components must not be empty", ErrInvalidConfig)
	}
	for _, component := range cfg.Workspace.Components {
		if !evidence.ValidID(component.ID) {
			return fmt.Errorf("%w: invalid component id %q", ErrInvalidConfig, component.ID)
		}
		if !validComponentRoot(component.Root) {
			return fmt.Errorf("%w: invalid component root %q", ErrInvalidConfig, component.Root)
		}
		if len(component.LanguageScopes) == 0 {
			return fmt.Errorf("%w: component %q has no languageScopes", ErrInvalidConfig, component.ID)
		}
		for _, scope := range component.LanguageScopes {
			if !evidence.ValidID(scope.ID) || !evidence.ValidID(scope.Language) || !evidence.ValidID(scope.Harness) {
				return fmt.Errorf("%w: invalid language scope in component %q", ErrInvalidConfig, component.ID)
			}
			if scope.Toolchain != "" && !evidence.ValidID(scope.Toolchain) {
				return fmt.Errorf("%w: invalid toolchain %q in scope %q", ErrInvalidConfig, scope.Toolchain, scope.ID)
			}
			if len(scope.Targets) == 0 {
				return fmt.Errorf("%w: language scope %q has no targets", ErrInvalidConfig, scope.ID)
			}
			for _, targetID := range scope.Targets {
				if !evidence.ValidID(targetID) {
					return fmt.Errorf("%w: invalid target id %q in scope %q", ErrInvalidConfig, targetID, scope.ID)
				}
			}
		}
	}
	definition, err := cfg.ToDefinition()
	if err != nil {
		return err
	}
	if err := chain.ValidateDefinition(definition); err != nil {
		return err
	}
	return nil
}

func validateAIConfig(config *AIConfig) error {
	if config == nil || len(config.Profiles) == 0 || len(config.Channels) == 0 {
		return fmt.Errorf("%w: ai requires profiles and channels", ErrInvalidConfig)
	}
	profileHashes := make(map[string]string, len(config.Profiles))
	previousProfileID := ""
	for _, profile := range config.Profiles {
		if previousProfileID != "" && profile.ProfileID <= previousProfileID {
			return fmt.Errorf("%w: ai.profiles must be sorted by unique profileId", ErrInvalidConfig)
		}
		hash, err := adapters.AIProviderProfileHash(profile)
		if err != nil {
			return fmt.Errorf("%w: profile %q: %v", ErrInvalidConfig, profile.ProfileID, err)
		}
		profileHashes[profile.ProfileID] = hash
		previousProfileID = profile.ProfileID
	}
	previousChannel := ""
	for _, binding := range config.Channels {
		if !validConfiguredAIChannel(binding.Channel) {
			return fmt.Errorf("%w: invalid ai channel %q", ErrInvalidConfig, binding.Channel)
		}
		if previousChannel != "" && binding.Channel <= previousChannel {
			return fmt.Errorf("%w: ai.channels must be sorted by unique channel", ErrInvalidConfig)
		}
		if _, found := profileHashes[binding.ProfileID]; !found {
			return fmt.Errorf("%w: ai channel %q references unknown profile %q", ErrInvalidConfig, binding.Channel, binding.ProfileID)
		}
		previousChannel = binding.Channel
	}
	return nil
}

func validConfiguredAIChannel(channel string) bool {
	switch channel {
	case "agent-runner-cli", "sessionbridge-silent", "sessionbridge-visible":
		return true
	default:
		return false
	}
}

func validateDocumentationPolicy(policy string) error {
	switch policy {
	case "required", "if-affected", "optional", "forbidden":
		return nil
	default:
		return fmt.Errorf("%w: invalid documentation policy %q", ErrInvalidConfig, policy)
	}
}

func validateHandoffPolicy(policy *HandoffPolicyInput) error {
	if policy == nil {
		return fmt.Errorf("%w: missing handoff policy", ErrInvalidConfig)
	}
	if len(policy.AllowedTargets) == 0 || len(policy.ReturnActions) == 0 || len(policy.HooksAfterReturn) == 0 {
		return fmt.Errorf("%w: handoff policy requires allowedTargets/returnActions/hooksAfterReturn", ErrInvalidConfig)
	}
	for _, target := range policy.AllowedTargets {
		if !evidence.ValidID(target) {
			return fmt.Errorf("%w: invalid handoff target %q", ErrInvalidConfig, target)
		}
	}
	switch policy.InputPolicy {
	case "structured", "secret-direct":
	default:
		return fmt.Errorf("%w: invalid handoff inputPolicy %q", ErrInvalidConfig, policy.InputPolicy)
	}
	if policy.HandoffTimeoutMs < 1 {
		return fmt.Errorf("%w: invalid handoff timeout", ErrInvalidConfig)
	}
	allowedActions := map[string]struct{}{"complete": {}, "abort": {}, "request-agent": {}}
	seenAction := map[string]struct{}{}
	for _, action := range policy.ReturnActions {
		if _, ok := allowedActions[action]; !ok {
			return fmt.Errorf("%w: invalid handoff return action %q", ErrInvalidConfig, action)
		}
		if _, exists := seenAction[action]; exists {
			return fmt.Errorf("%w: duplicate handoff return action %q", ErrInvalidConfig, action)
		}
		seenAction[action] = struct{}{}
	}
	for _, hook := range policy.HooksAfterReturn {
		if !evidence.ValidID(hook) {
			return fmt.Errorf("%w: invalid handoff hook %q", ErrInvalidConfig, hook)
		}
	}
	return nil
}

func validComponentRoot(root string) bool {
	if root == "." {
		return true
	}
	if root == "" || strings.Contains(root, "\\") || strings.ContainsRune(root, '\x00') || strings.HasPrefix(root, "/") {
		return false
	}
	if len(root) >= 2 && root[1] == ':' {
		return false
	}
	if index := strings.Index(root, ":"); index > 0 {
		return false
	}
	for _, segment := range strings.Split(root, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
	}
	return true
}

func (cfg ChainConfig) ToDefinition() (chain.Definition, error) {
	definition := chain.Definition{ID: cfg.Chain.ID, Tasks: make([]chain.Task, 0, len(cfg.Tasks))}
	for _, task := range cfg.Tasks {
		mappedTask := chain.Task{ID: task.ID, Steps: make([]chain.Step, 0, len(task.Steps))}
		for _, step := range task.Steps {
			mappedStep := chain.Step{ID: step.ID, Kind: step.Kind, Reason: step.Reason}
			switch step.Kind {
			case "code":
				execution := step.Execution
				if execution == "" {
					execution = "autonomous"
				}
				if execution == "manual-handoff" {
					if step.HandoffPolicy == nil {
						return chain.Definition{}, fmt.Errorf("%w: manual-handoff step %q requires handoffPolicy", ErrInvalidConfig, step.ID)
					}
					mappedStep.Mode = chain.ManualHandoff
					mappedStep.HandoffPolicy = &chain.HandoffPolicy{
						AllowedTargets:   append([]string{}, step.HandoffPolicy.AllowedTargets...),
						InputPolicy:      step.HandoffPolicy.InputPolicy,
						HandoffTimeoutMs: step.HandoffPolicy.HandoffTimeoutMs,
						ReturnActions:    append([]string{}, step.HandoffPolicy.ReturnActions...),
						HooksAfterReturn: append([]string{}, step.HandoffPolicy.HooksAfterReturn...),
					}
				} else {
					mappedStep.Mode = chain.IsolatedWorkspace
				}
			case "build", "verify", "noop":
			default:
				return chain.Definition{}, fmt.Errorf("%w: unsupported step kind %q", ErrInvalidConfig, step.Kind)
			}
			mappedTask.Steps = append(mappedTask.Steps, mappedStep)
		}
		definition.Tasks = append(definition.Tasks, mappedTask)
	}
	return definition, nil
}

func StepCount(cfg ChainConfig) int {
	total := 0
	for _, task := range cfg.Tasks {
		total += len(task.Steps)
	}
	return total
}

func FindExecutableSteps(cfg ChainConfig) []ExecutableStep {
	steps := make([]ExecutableStep, 0)
	for _, task := range cfg.Tasks {
		for _, step := range task.Steps {
			if step.Kind == "noop" {
				continue
			}
			steps = append(steps, ExecutableStep{TaskID: task.ID, StepID: step.ID, Kind: step.Kind})
		}
	}
	return steps
}

func ValidateSummary(path string, cfg ChainConfig) ConfigValidationSummary {
	executable := FindExecutableSteps(cfg)
	warnings := []string{}
	if len(executable) > 0 {
		warnings = append(warnings, "local CLI runtime currently executes noop-only chains; use validate/config explain for executable steps")
	}
	return ConfigValidationSummary{
		ConfigPath:      path,
		ChainID:         cfg.Chain.ID,
		TaskCount:       len(cfg.Tasks),
		StepCount:       StepCount(cfg),
		RunnableInCLI:   len(executable) == 0,
		ExecutableSteps: executable,
		Warnings:        warnings,
	}
}

func BuildConfigExplain(path string, cfg ChainConfig, rawConfig []byte) ConfigExplainReport {
	configHash := evidence.Digest("", rawConfig)
	report := ConfigExplainReport{
		ConfigPath:          path,
		ConfigHash:          configHash,
		ChainID:             cfg.Chain.ID,
		Profile:             cfg.Chain.Profile,
		AIProfiles:          []AIProfileExplain{},
		AIChannels:          []AIChannelExplain{},
		DocumentationPolicy: effectiveChainDocumentationPolicy(cfg),
		RunnableInCLI:       len(FindExecutableSteps(cfg)) == 0,
		ExecutableSteps:     FindExecutableSteps(cfg),
		TaskPolicies:        []TaskPolicyExplain{},
		Resolution:          []ResolutionExplain{},
		ConfigSearchOrder:   []string{"explicit --chain path", filepath.Join("<cwd>", DefaultChainConfigName), filepath.Join("<cwd>", FallbackChainConfigName)},
	}
	addResolution := func(pointer, sourceKind string, sourcePointer *string) {
		var sourceHash *string
		if sourceKind != "builtin-default" {
			hash := configHash
			sourceHash = &hash
		}
		report.Resolution = append(report.Resolution, ResolutionExplain{
			EffectivePointer: pointer,
			SourceKind:       sourceKind,
			SourceHash:       sourceHash,
			SourcePointer:    sourcePointer,
		})
	}

	profileHashes := make(map[string]string)
	if cfg.AI != nil {
		for index, profile := range cfg.AI.Profiles {
			hash, err := adapters.AIProviderProfileHash(profile)
			if err != nil {
				continue
			}
			profileHashes[profile.ProfileID] = hash
			report.AIProfiles = append(report.AIProfiles, AIProfileExplain{ProfileID: profile.ProfileID, ProfileConfigHash: hash})
			pointer := fmt.Sprintf("/ai/profiles/%d", index)
			addResolution(pointer, "chain", &pointer)
		}
		for index, binding := range cfg.AI.Channels {
			report.AIChannels = append(report.AIChannels, AIChannelExplain{
				Channel:           binding.Channel,
				ProfileID:         binding.ProfileID,
				ProfileConfigHash: profileHashes[binding.ProfileID],
			})
			pointer := fmt.Sprintf("/ai/channels/%d", index)
			addResolution(pointer, "chain", &pointer)
		}
	}

	chainIDPointer := "/chain/id"
	chainProfilePointer := "/chain/profile"
	addResolution("/chain/id", "chain", &chainIDPointer)
	addResolution("/chain/profile", "chain", &chainProfilePointer)

	if cfg.Documentation != nil && cfg.Documentation.Policy != "" {
		policyPointer := "/documentation/policy"
		addResolution("/documentation/policy", "chain", &policyPointer)
	} else {
		addResolution("/documentation/policy", "builtin-default", nil)
	}

	for index, task := range cfg.Tasks {
		policy, sourceKind, sourcePointer := resolveTaskDocumentationPolicy(cfg, index)
		var sourceHash *string
		if sourceKind != "builtin-default" {
			hash := configHash
			sourceHash = &hash
		}
		report.TaskPolicies = append(report.TaskPolicies, TaskPolicyExplain{
			TaskID:        task.ID,
			Policy:        policy,
			SourceKind:    sourceKind,
			SourceHash:    sourceHash,
			SourcePointer: sourcePointer,
		})
		pointer := fmt.Sprintf("/tasks/%d/documentationPolicy", index)
		addResolution(pointer, sourceKind, sourcePointer)
	}

	if len(report.Resolution) > 1 {
		sort.SliceStable(report.Resolution, func(i, j int) bool {
			return report.Resolution[i].EffectivePointer < report.Resolution[j].EffectivePointer
		})
	}
	return report
}

func effectiveChainDocumentationPolicy(cfg ChainConfig) string {
	if cfg.Documentation != nil && cfg.Documentation.Policy != "" {
		return cfg.Documentation.Policy
	}
	return defaultDocumentationRule
}

func resolveTaskDocumentationPolicy(cfg ChainConfig, taskIndex int) (policy, sourceKind string, sourcePointer *string) {
	task := cfg.Tasks[taskIndex]
	if task.DocumentationPolicy != "" {
		pointer := fmt.Sprintf("/tasks/%d/documentationPolicy", taskIndex)
		return task.DocumentationPolicy, "task", &pointer
	}
	if cfg.Documentation != nil && cfg.Documentation.Policy != "" {
		pointer := "/documentation/policy"
		return cfg.Documentation.Policy, "chain", &pointer
	}
	return defaultDocumentationRule, "builtin-default", nil
}

func NormalizeChainID(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" {
		return "proofrail"
	}
	builder := strings.Builder{}
	builder.Grow(len(raw))
	previousDash := false
	for _, r := range raw {
		isLetter := r >= 'a' && r <= 'z'
		isDigit := r >= '0' && r <= '9'
		if isLetter || isDigit {
			builder.WriteRune(r)
			previousDash = false
			continue
		}
		if !previousDash {
			builder.WriteByte('-')
			previousDash = true
		}
	}
	id := strings.Trim(builder.String(), "-")
	if id == "" {
		id = "proofrail"
	}
	if id[0] < 'a' || id[0] > 'z' {
		id = "p-" + id
	}
	if len(id) > 64 {
		id = strings.Trim(id[:64], "-")
		if id == "" || id[0] < 'a' || id[0] > 'z' {
			id = "proofrail"
		}
	}
	if !evidence.ValidID(id) {
		return "proofrail"
	}
	return id
}

func EncodeChainConfig(cfg ChainConfig) ([]byte, error) {
	encoded, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(encoded, '\n'), nil
}
