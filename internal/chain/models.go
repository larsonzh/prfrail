package chain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrInvalidDefinition = errors.New("invalid chain definition")
	ErrInvalidState      = errors.New("invalid chain state")
	ErrPaused            = errors.New("chain paused")
	ErrCancelled         = errors.New("chain cancelled")
	ErrStepFailed        = errors.New("chain step failed")
	ErrRecoveryUncertain = errors.New("chain recovery uncertain")
)

type ChangeMode string

const (
	ManagedChangeSet  ChangeMode = "managed-change-set"
	IsolatedWorkspace ChangeMode = "isolated-workspace"
)

type Step struct {
	ID     string
	Kind   string
	Mode   ChangeMode
	Reason string
}

type Task struct {
	ID    string
	Steps []Step
}

type Definition struct {
	ID    string
	Tasks []Task
}

type SnapshotRef struct {
	Hash string
}

type Workspace struct {
	Root string
}

type StepRequest struct {
	RunID      string
	TaskID     string
	Step       Step
	Attempt    int
	ParentHash string
	Workspace  Workspace
}

type StepResult struct {
	Evidence []string
}

type CandidateResult struct {
	Snapshot          SnapshotRef
	EvidenceRootHash  string
	Evidence          []string
	CandidateProducer evidence.Actor
}

type ReviewRequest struct {
	RunID                 string
	TaskID                string
	Attempt               int
	ParentSnapshotHash    string
	CandidateSnapshotHash string
	EvidenceRootHash      string
	CandidateProducer     evidence.Actor
}

type ReviewReason struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

type WaiverAuthorization struct {
	AuthorizedBy    evidence.Actor `json:"authorizedBy"`
	PolicyBasisHash string         `json:"policyBasisHash"`
	Scope           []string       `json:"scope"`
	ExpiresAt       string         `json:"expiresAt"`
}

type ReviewDecision struct {
	ReceiptID             string
	OccurredAt            string
	RecordedBy            evidence.Actor
	ReviewMode            string
	PolicyHash            string
	Outcome               string
	CandidateSnapshotHash string
	Evidence              []string
	ErrorEvidence         []string
	Reason                *ReviewReason
	WaiverAuthorization   *WaiverAuthorization
}

type PromotionRequest struct {
	RunID                 string
	TaskID                string
	Attempt               int
	ParentSnapshotHash    string
	CandidateSnapshotHash string
	EvidenceRootHash      string
	ReviewReceiptHash     string
}

type PromotionDecision struct {
	ReceiptID            string
	OccurredAt           string
	Outcome              string
	AcceptedSnapshotHash string
	WriterStopEvidence   []string
	LeaseEvidence        []string
	Evidence             []string
	ErrorEvidence        []string
}

type BaselinePort interface {
	CaptureBaseline(context.Context, string) (SnapshotRef, error)
}

type WorkspacePort interface {
	Materialize(context.Context, string, string, int, SnapshotRef) (Workspace, error)
}

type StepPort interface {
	Execute(context.Context, StepRequest) (StepResult, error)
}

type AcceptancePort interface {
	Accept(context.Context, string, string, int, SnapshotRef, Workspace) (CandidateResult, error)
}

type ReviewerPort interface {
	Review(context.Context, ReviewRequest) (ReviewDecision, error)
}

type PublisherPort interface {
	Publish(context.Context, PromotionRequest) (PromotionDecision, error)
}

type Stopper interface {
	Stop(context.Context, string) ([]string, error)
}

type Reconciler interface {
	Reconcile(context.Context, string) error
}

type Clock func() time.Time
type IDSource func() string

func (definition Definition) validate() error {
	if !evidence.ValidID(definition.ID) || len(definition.Tasks) == 0 {
		return ErrInvalidDefinition
	}
	seen := map[string]struct{}{}
	for _, task := range definition.Tasks {
		if !evidence.ValidID(task.ID) || len(task.Steps) == 0 {
			return ErrInvalidDefinition
		}
		if _, exists := seen[task.ID]; exists {
			return fmt.Errorf("%w: duplicate ID %q", ErrInvalidDefinition, task.ID)
		}
		seen[task.ID] = struct{}{}
		for _, step := range task.Steps {
			if !evidence.ValidID(step.ID) {
				return ErrInvalidDefinition
			}
			if _, exists := seen[step.ID]; exists {
				return fmt.Errorf("%w: duplicate ID %q", ErrInvalidDefinition, step.ID)
			}
			seen[step.ID] = struct{}{}
			switch step.Kind {
			case "code":
				if step.Mode != ManagedChangeSet && step.Mode != IsolatedWorkspace {
					return fmt.Errorf("%w: code step %q has invalid mode", ErrInvalidDefinition, step.ID)
				}
			case "build", "verify":
				if step.Mode != "" {
					return fmt.Errorf("%w: non-code step %q has a change mode", ErrInvalidDefinition, step.ID)
				}
			case "noop":
				if step.Mode != "" || step.Reason == "" {
					return fmt.Errorf("%w: noop step %q requires only a reason", ErrInvalidDefinition, step.ID)
				}
			default:
				return fmt.Errorf("%w: step %q has unknown kind", ErrInvalidDefinition, step.ID)
			}
		}
	}
	return nil
}
