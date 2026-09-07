package taskdef

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	previewTimestampLayout = "2006-01-02T15:04:05.000Z"
	planPreviewHashDomain  = "proofrail:plan-preview:1\n"
)

var ErrInvalidPlanPreview = errors.New("invalid plan preview")

var networkHostPattern = regexp.MustCompile(`^[A-Za-z0-9.-]+$`)

type PlanPreviewRecord struct {
	SchemaVersion string      `json:"schemaVersion"`
	Preview       PlanPreview `json:"preview"`
	PreviewHash   string      `json:"previewHash"`
}

type PlanPreview struct {
	PreviewID               string                          `json:"previewId"`
	CreatedAt               string                          `json:"createdAt"`
	Mode                    string                          `json:"mode"`
	ChainDefinitionHash     string                          `json:"chainDefinitionHash"`
	WorkspaceDefinitionHash string                          `json:"workspaceDefinitionHash"`
	EffectiveConfigHash     string                          `json:"effectiveConfigHash"`
	PolicyHash              string                          `json:"policyHash"`
	Locations               PlanPreviewLocations            `json:"locations"`
	ReadTargetIDs           []string                        `json:"readTargetIds"`
	WriteTargetIDs          []string                        `json:"writeTargetIds"`
	NetworkRequirements     []PlanPreviewNetworkRequirement `json:"networkRequirements"`
	Steps                   []PlanPreviewStep               `json:"steps"`
	ReviewPointIDs          []string                        `json:"reviewPointIds"`
	Budget                  PlanPreviewBudget               `json:"budget"`
	EstimatedStoreBytes     *int64                          `json:"estimatedStoreBytes"`
	ExclusionIDs            []string                        `json:"exclusionIds"`
	Capabilities            []PlanPreviewCapability         `json:"capabilities"`
	Outcome                 string                          `json:"outcome"`
	BlockingEvidence        []string                        `json:"blockingEvidence"`
}

type PlanPreviewLocations struct {
	SourceRef string `json:"sourceRef"`
	RunRef    string `json:"runRef"`
	StoreRef  string `json:"storeRef"`
}

type PlanPreviewNetworkRequirement struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Purpose  string `json:"purpose"`
}

type PlanPreviewStep struct {
	TaskID          string   `json:"taskId"`
	StepID          string   `json:"stepId"`
	Kind            string   `json:"kind"`
	BlockingGateIDs []string `json:"blockingGateIds"`
	ReviewRequired  bool     `json:"reviewRequired"`
}

type PlanPreviewCapability struct {
	CapabilityID string   `json:"capabilityId"`
	Status       string   `json:"status"`
	Evidence     []string `json:"evidence"`
}

type PlanPreviewBudget struct {
	ModelCallLimit        int     `json:"modelCallLimit"`
	TokenLimit            int     `json:"tokenLimit"`
	WallClockMs           int     `json:"wallClockMs"`
	AttemptLimit          int     `json:"attemptLimit"`
	Currency              string  `json:"currency"`
	EstimatedAmountMicros *int64  `json:"estimatedAmountMicros"`
	EstimateStatus        string  `json:"estimateStatus"`
	PricingSource         *string `json:"pricingSource"`
}

func NewPlanPreviewRecord(preview PlanPreview) (PlanPreviewRecord, error) {
	if err := validatePlanPreview(preview); err != nil {
		return PlanPreviewRecord{}, err
	}
	canonical, err := evidence.EncodeCanonical(preview)
	if err != nil {
		return PlanPreviewRecord{}, fmt.Errorf("%w: canonicalize preview: %v", ErrInvalidPlanPreview, err)
	}
	return PlanPreviewRecord{
		SchemaVersion: "1.0.0",
		Preview:       preview,
		PreviewHash:   evidence.Digest(planPreviewHashDomain, canonical),
	}, nil
}

func ValidatePlanPreviewRecord(record PlanPreviewRecord) error {
	if record.SchemaVersion != "1.0.0" {
		return fmt.Errorf("%w: unsupported schemaVersion %q", ErrInvalidPlanPreview, record.SchemaVersion)
	}
	if err := validatePlanPreview(record.Preview); err != nil {
		return err
	}
	canonical, err := evidence.EncodeCanonical(record.Preview)
	if err != nil {
		return fmt.Errorf("%w: canonicalize preview: %v", ErrInvalidPlanPreview, err)
	}
	expected := evidence.Digest(planPreviewHashDomain, canonical)
	if record.PreviewHash != expected {
		return fmt.Errorf("%w: preview hash mismatch", ErrInvalidPlanPreview)
	}
	return nil
}

func validatePlanPreview(preview PlanPreview) error {
	if !evidence.ValidID(preview.PreviewID) {
		return fmt.Errorf("%w: invalid previewId", ErrInvalidPlanPreview)
	}
	if err := validateTimestamp(preview.CreatedAt); err != nil {
		return err
	}
	if preview.Mode != "offline-read-only" {
		return fmt.Errorf("%w: invalid mode %q", ErrInvalidPlanPreview, preview.Mode)
	}
	if !evidence.ValidHash(preview.ChainDefinitionHash) ||
		!evidence.ValidHash(preview.WorkspaceDefinitionHash) ||
		!evidence.ValidHash(preview.EffectiveConfigHash) ||
		!evidence.ValidHash(preview.PolicyHash) {
		return fmt.Errorf("%w: invalid hash field", ErrInvalidPlanPreview)
	}
	if !evidence.ValidID(preview.Locations.SourceRef) || !evidence.ValidID(preview.Locations.RunRef) || !evidence.ValidID(preview.Locations.StoreRef) {
		return fmt.Errorf("%w: invalid location references", ErrInvalidPlanPreview)
	}
	if err := validateIDSet(preview.ReadTargetIDs, "readTargetIds"); err != nil {
		return err
	}
	if err := validateIDSet(preview.WriteTargetIDs, "writeTargetIds"); err != nil {
		return err
	}
	if err := validateIDSet(preview.ReviewPointIDs, "reviewPointIds"); err != nil {
		return err
	}
	if err := validateIDSet(preview.ExclusionIDs, "exclusionIds"); err != nil {
		return err
	}
	for _, requirement := range preview.NetworkRequirements {
		if requirement.Protocol != "http" && requirement.Protocol != "https" && requirement.Protocol != "ssh" && requirement.Protocol != "custom" {
			return fmt.Errorf("%w: invalid network protocol %q", ErrInvalidPlanPreview, requirement.Protocol)
		}
		if strings.TrimSpace(requirement.Host) == "" || len(requirement.Host) > 253 || !networkHostPattern.MatchString(requirement.Host) {
			return fmt.Errorf("%w: invalid network host %q", ErrInvalidPlanPreview, requirement.Host)
		}
		if requirement.Port < 1 || requirement.Port > 65535 {
			return fmt.Errorf("%w: invalid network port", ErrInvalidPlanPreview)
		}
		if !evidence.ValidID(requirement.Purpose) {
			return fmt.Errorf("%w: invalid network purpose", ErrInvalidPlanPreview)
		}
	}
	if len(preview.Steps) == 0 {
		return fmt.Errorf("%w: preview requires at least one step", ErrInvalidPlanPreview)
	}
	for _, step := range preview.Steps {
		if !evidence.ValidID(step.TaskID) || !evidence.ValidID(step.StepID) {
			return fmt.Errorf("%w: invalid step taskId/stepId", ErrInvalidPlanPreview)
		}
		switch step.Kind {
		case "code", "build", "verify", "noop", "manual-handoff":
		default:
			return fmt.Errorf("%w: invalid step kind %q", ErrInvalidPlanPreview, step.Kind)
		}
		if err := validateIDSet(step.BlockingGateIDs, "step.blockingGateIds"); err != nil {
			return err
		}
	}
	if err := validateBudget(preview.Budget); err != nil {
		return err
	}
	if preview.EstimatedStoreBytes != nil && *preview.EstimatedStoreBytes < 0 {
		return fmt.Errorf("%w: estimatedStoreBytes must be >= 0", ErrInvalidPlanPreview)
	}
	for _, capability := range preview.Capabilities {
		if !evidence.ValidID(capability.CapabilityID) {
			return fmt.Errorf("%w: invalid capabilityId %q", ErrInvalidPlanPreview, capability.CapabilityID)
		}
		if err := validateHashSet(capability.Evidence, "capability.evidence"); err != nil {
			return err
		}
		switch capability.Status {
		case "observed", "verified":
			if len(capability.Evidence) == 0 {
				return fmt.Errorf("%w: %s capability requires evidence", ErrInvalidPlanPreview, capability.Status)
			}
		case "planned", "unknown", "unavailable":
			if len(capability.Evidence) != 0 {
				return fmt.Errorf("%w: %s capability cannot include evidence", ErrInvalidPlanPreview, capability.Status)
			}
		default:
			return fmt.Errorf("%w: invalid capability status %q", ErrInvalidPlanPreview, capability.Status)
		}
	}
	if err := validateHashSet(preview.BlockingEvidence, "blockingEvidence"); err != nil {
		return err
	}
	switch preview.Outcome {
	case "ready":
		if len(preview.BlockingEvidence) != 0 {
			return fmt.Errorf("%w: ready preview cannot include blocking evidence", ErrInvalidPlanPreview)
		}
	case "blocked":
		if len(preview.BlockingEvidence) == 0 {
			return fmt.Errorf("%w: blocked preview requires blocking evidence", ErrInvalidPlanPreview)
		}
	default:
		return fmt.Errorf("%w: invalid outcome %q", ErrInvalidPlanPreview, preview.Outcome)
	}
	return nil
}

func validateTimestamp(value string) error {
	parsed, err := time.Parse(previewTimestampLayout, value)
	if err != nil {
		return fmt.Errorf("%w: invalid createdAt", ErrInvalidPlanPreview)
	}
	if parsed.UTC().Format(previewTimestampLayout) != value {
		return fmt.Errorf("%w: createdAt must be UTC millisecond timestamp", ErrInvalidPlanPreview)
	}
	return nil
}

func validateIDSet(values []string, field string) error {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !evidence.ValidID(value) {
			return fmt.Errorf("%w: invalid %s value %q", ErrInvalidPlanPreview, field, value)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%w: duplicate %s value %q", ErrInvalidPlanPreview, field, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validateHashSet(values []string, field string) error {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !evidence.ValidHash(value) {
			return fmt.Errorf("%w: invalid %s hash %q", ErrInvalidPlanPreview, field, value)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%w: duplicate %s hash %q", ErrInvalidPlanPreview, field, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validateBudget(budget PlanPreviewBudget) error {
	if budget.ModelCallLimit < 0 || budget.TokenLimit < 0 || budget.WallClockMs < 1 || budget.AttemptLimit < 1 {
		return fmt.Errorf("%w: invalid budget limits", ErrInvalidPlanPreview)
	}
	if len(budget.Currency) != 3 {
		return fmt.Errorf("%w: currency must be 3-letter ISO code", ErrInvalidPlanPreview)
	}
	for _, character := range budget.Currency {
		if character < 'A' || character > 'Z' {
			return fmt.Errorf("%w: currency must be uppercase letters", ErrInvalidPlanPreview)
		}
	}
	if budget.EstimatedAmountMicros != nil && *budget.EstimatedAmountMicros < 0 {
		return fmt.Errorf("%w: estimatedAmountMicros must be >= 0", ErrInvalidPlanPreview)
	}
	if budget.PricingSource != nil && (strings.TrimSpace(*budget.PricingSource) == "" || len(*budget.PricingSource) > 256) {
		return fmt.Errorf("%w: pricingSource must be 1-256 characters", ErrInvalidPlanPreview)
	}
	switch budget.EstimateStatus {
	case "estimated":
		if budget.EstimatedAmountMicros == nil || budget.PricingSource == nil {
			return fmt.Errorf("%w: estimated status requires amount and pricingSource", ErrInvalidPlanPreview)
		}
	case "unknown":
		if budget.EstimatedAmountMicros != nil {
			return fmt.Errorf("%w: unknown estimate status cannot carry amount", ErrInvalidPlanPreview)
		}
	case "not-applicable":
		if budget.EstimatedAmountMicros != nil {
			return fmt.Errorf("%w: not-applicable estimate cannot carry amount", ErrInvalidPlanPreview)
		}
	default:
		return fmt.Errorf("%w: invalid estimateStatus %q", ErrInvalidPlanPreview, budget.EstimateStatus)
	}
	return nil
}
