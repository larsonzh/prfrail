package gates

import (
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type Hook struct {
	ID                string
	Kind              string
	Runner            ProcessRunner
	CWD               string
	EnvAllowlist      []string
	EffectClass       string
	ExternalSystems   []string
	RecoveryGuarantee string
	PolicyHash        string
	AuthorizationHash string
	EffectEvidence    []string
	OnFail            FailurePolicy
	ArtifactGlobs     []string
	Timeout           time.Duration
	Resources         ResourceLimits
	Network           NetworkPolicy
}

type EffectScope struct {
	EffectClass       string
	ExternalSystems   []string
	RecoveryGuarantee string
	PolicyHash        string
	AuthorizationHash string
	OperationKey      string
	Evidence          []string
}

type ProcessRunner struct {
	Type       string
	Executable string
	Args       []string
}

type FailurePolicy struct {
	Action      string
	MaxAttempts int
}

type ResourceLimits struct {
	MemoryBytes  int64
	OutputBytes  int64
	ProcessCount int
	CPUTime      time.Duration
}

type NetworkPolicy struct {
	Mode  string   `json:"mode"`
	Hosts []string `json:"hosts,omitempty"`
	Ports []int    `json:"ports,omitempty"`
}

type Capabilities struct {
	Process          bool
	Timeout          bool
	OutputLimit      bool
	MemoryLimit      bool
	ProcessLimit     bool
	CPUTimeLimit     bool
	NetworkDeny      bool
	NetworkLoopback  bool
	NetworkAllowlist bool
}

type Execution struct {
	StartedAt           *time.Time
	FinishedAt          time.Time
	Outcome             string
	ExitCode            *int
	Stdout              []byte
	Stderr              []byte
	OutputTruncated     bool
	RunnerEvidence      []string
	TerminationEvidence []string
	ErrorEvidence       []string
}

type Request struct {
	RunID              string
	TaskID             string
	StepID             string
	Attempt            int
	ExecutionAttempt   int
	WorkspaceRoot      string
	HookDefinitionHash string
	Effect             *EffectScope
	Hook               Hook
}

type Output struct {
	ObjectHash string `json:"objectHash"`
	ByteLength int64  `json:"byteLength"`
	Capture    string `json:"capture"`
	Redaction  string `json:"redaction"`
}

type Artifact struct {
	ID         string `json:"id"`
	ObjectHash string `json:"objectHash"`
	MediaType  string `json:"mediaType"`
	ByteLength int64  `json:"byteLength"`
	Redaction  string `json:"redaction"`
}

type Result struct {
	ResultID              string         `json:"resultId"`
	RecordedBy            evidence.Actor `json:"recordedBy"`
	RunID                 string         `json:"runId"`
	TaskID                string         `json:"taskId"`
	StepID                string         `json:"stepId"`
	Attempt               int            `json:"attempt"`
	HookID                string         `json:"hookId"`
	HookDefinitionHash    string         `json:"hookDefinitionHash"`
	ExecutionAttempt      int            `json:"executionAttempt"`
	AttemptedAt           string         `json:"attemptedAt"`
	StartedAt             *string        `json:"startedAt"`
	FinishedAt            string         `json:"finishedAt"`
	ExecutionOutcome      string         `json:"executionOutcome"`
	ExitCode              *int           `json:"exitCode"`
	Assessment            string         `json:"assessment"`
	FailureKind           *string        `json:"failureKind"`
	PolicyDisposition     string         `json:"policyDisposition"`
	EffectObservationHash *string        `json:"effectObservationHash,omitempty"`
	EffectRecoveryAction  *string        `json:"effectRecoveryAction,omitempty"`
	Stdout                *Output        `json:"stdout"`
	Stderr                *Output        `json:"stderr"`
	Artifacts             []Artifact     `json:"artifacts"`
	RunnerEvidence        []string       `json:"runnerEvidence"`
	TerminationEvidence   []string       `json:"terminationEvidence"`
	ErrorEvidence         []string       `json:"errorEvidence"`
}

type ResultRecord struct {
	SchemaVersion string `json:"schemaVersion"`
	Result        Result `json:"result"`
	ResultHash    string `json:"resultHash"`
}
