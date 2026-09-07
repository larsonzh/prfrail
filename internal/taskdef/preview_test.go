package taskdef

import (
	"errors"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestPlanPreviewRecordReadyRoundTrip(t *testing.T) {
	preview := validPlanPreview()
	record, err := NewPlanPreviewRecord(preview)
	if err != nil {
		t.Fatal(err)
	}
	if record.SchemaVersion != "1.0.0" || !evidence.ValidHash(record.PreviewHash) {
		t.Fatalf("unexpected record header: %+v", record)
	}
	if err := ValidatePlanPreviewRecord(record); err != nil {
		t.Fatalf("record validation failed: %v", err)
	}
}

func TestPlanPreviewRecordBlockedRequiresEvidence(t *testing.T) {
	preview := validPlanPreview()
	preview.Outcome = "blocked"
	preview.BlockingEvidence = nil
	if _, err := NewPlanPreviewRecord(preview); !errors.Is(err, ErrInvalidPlanPreview) {
		t.Fatalf("expected blocked-without-evidence rejection, got %v", err)
	}
}

func TestPlanPreviewRecordCapabilityEvidenceRules(t *testing.T) {
	preview := validPlanPreview()
	preview.Capabilities = []PlanPreviewCapability{{
		CapabilityID: "network-access",
		Status:       "unknown",
		Evidence:     []string{digest("unexpected")},
	}}
	if _, err := NewPlanPreviewRecord(preview); !errors.Is(err, ErrInvalidPlanPreview) {
		t.Fatalf("expected unknown capability with evidence rejection, got %v", err)
	}
}

func TestPlanPreviewRecordDetectsHashMismatch(t *testing.T) {
	preview := validPlanPreview()
	record, err := NewPlanPreviewRecord(preview)
	if err != nil {
		t.Fatal(err)
	}
	record.PreviewHash = digest("tampered")
	if err := ValidatePlanPreviewRecord(record); !errors.Is(err, ErrInvalidPlanPreview) {
		t.Fatalf("expected hash mismatch rejection, got %v", err)
	}
}

func TestPlanPreviewRecordRejectsInvalidNetworkHost(t *testing.T) {
	preview := validPlanPreview()
	preview.NetworkRequirements = []PlanPreviewNetworkRequirement{{
		Protocol: "https",
		Host:     "bad_host",
		Port:     443,
		Purpose:  "fetch",
	}}
	if _, err := NewPlanPreviewRecord(preview); !errors.Is(err, ErrInvalidPlanPreview) {
		t.Fatalf("expected invalid network host rejection, got %v", err)
	}
}

func TestPlanPreviewRecordRejectsLongPricingSource(t *testing.T) {
	preview := validPlanPreview()
	longSource := strings.Repeat("s", 257)
	amount := int64(100)
	preview.Budget.EstimateStatus = "estimated"
	preview.Budget.EstimatedAmountMicros = &amount
	preview.Budget.PricingSource = &longSource
	if _, err := NewPlanPreviewRecord(preview); !errors.Is(err, ErrInvalidPlanPreview) {
		t.Fatalf("expected long pricingSource rejection, got %v", err)
	}
}

func validPlanPreview() PlanPreview {
	staticEvidence := digest("static")
	return PlanPreview{
		PreviewID:               "preview-one",
		CreatedAt:               "2026-09-08T12:00:00.000Z",
		Mode:                    "offline-read-only",
		ChainDefinitionHash:     digest("chain"),
		WorkspaceDefinitionHash: digest("workspace"),
		EffectiveConfigHash:     digest("config"),
		PolicyHash:              digest("policy"),
		Locations: PlanPreviewLocations{
			SourceRef: "source",
			RunRef:    "run",
			StoreRef:  "store",
		},
		ReadTargetIDs:  []string{"chain-config", "source"},
		WriteTargetIDs: []string{},
		Steps: []PlanPreviewStep{{
			TaskID:          "task-one",
			StepID:          "step-one",
			Kind:            "noop",
			BlockingGateIDs: []string{},
			ReviewRequired:  true,
		}},
		ReviewPointIDs: []string{"task-one"},
		Budget: PlanPreviewBudget{
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
		Capabilities: []PlanPreviewCapability{{
			CapabilityID: "static-preview",
			Status:       "verified",
			Evidence:     []string{staticEvidence},
		}},
		Outcome:          "ready",
		BlockingEvidence: []string{},
	}
}

func digest(label string) string {
	return evidence.Digest("", []byte(label))
}
