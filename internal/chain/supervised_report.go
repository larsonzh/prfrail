package chain

import (
	"fmt"
	"slices"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// The reduced/unknown report is the factual record of what a supervised
// black-box run can and cannot prove (contracts section 7.1d). The assurance
// level is permanently "reduced": it is the ceiling of the mode, never a
// degraded state. The unknown dimension set is non-empty, always contains
// "cost", and every member comes from the frozen closed set, so an unknown
// dimension can never be rendered as safe.
const (
	supervisedReportDomain     = "proofrail:supervised-report:1\n"
	supervisedReportKind       = "supervised-report"
	supervisedReportSchema     = "1.0.0"
	supervisedReportMode       = "supervised-black-box"
	supervisedAssuranceReduced = "reduced"
	supervisedUnknownCost      = "cost"
)

// supervisedUnknownDimensionSet is the frozen closed set of dimensions that are
// unobservable for an external agent.
var supervisedUnknownDimensionSet = []string{
	"tool-use",
	"resource-reads",
	"network-access",
	"cost",
	"background-processes",
	"session-continuity",
	"external-side-effects",
}

// supervisedVerifiedDimensionSet is the frozen closed set of dimensions
// ProofRail proves about its own obligations.
var supervisedVerifiedDimensionSet = []string{
	"workspace-isolation",
	"manifest-before",
	"manifest-after",
	"diff",
	"scope-scan",
	"secret-scan",
	"file-type-scan",
	"side-effect-scan",
	"process-stop-proof",
	"independent-gates",
	"independent-review",
}

// SupervisedReport is the report body. AssuranceLevel is not an input field:
// the constructor always sets the reduced constant.
type SupervisedReport struct {
	Kind               string         `json:"kind"`
	ReportID           string         `json:"reportId"`
	RecordedAt         string         `json:"recordedAt"`
	RecordedBy         evidence.Actor `json:"recordedBy"`
	RunID              string         `json:"runId"`
	TaskID             string         `json:"taskId"`
	StepID             string         `json:"stepId"`
	Attempt            int            `json:"attempt"`
	Mode               string         `json:"mode"`
	AssuranceLevel     string         `json:"assuranceLevel"`
	VerifiedDimensions []string       `json:"verifiedDimensions"`
	UnknownDimensions  []string       `json:"unknownDimensions"`
	Evidence           []string       `json:"evidence"`
}

// SupervisedReportRecord is the persisted envelope.
type SupervisedReportRecord struct {
	SchemaVersion string           `json:"schemaVersion"`
	Record        SupervisedReport `json:"record"`
	RecordHash    string           `json:"recordHash"`
}

// SupervisedReportInput is the caller-supplied body of a new report.
type SupervisedReportInput struct {
	ReportID           string
	RecordedAt         string
	RunID              string
	TaskID             string
	StepID             string
	Attempt            int
	VerifiedDimensions []string
	UnknownDimensions  []string
	Evidence           []string
}

// NewSupervisedReport builds and validates a report. The assurance level is
// always the reduced constant, and the unknown dimension set always contains
// "cost": a caller can list fewer unobservable dimensions, never fewer than the
// mode demands, and can never raise the level.
func NewSupervisedReport(input SupervisedReportInput, clock Clock, ids IDSource) (SupervisedReportRecord, error) {
	if clock == nil {
		clock = defaultClock
	}
	if ids == nil {
		ids = randomID
	}
	unknown := cloneStrings(input.UnknownDimensions)
	if len(unknown) == 0 {
		unknown = cloneStrings(supervisedUnknownDimensionSet)
	}
	if !slices.Contains(unknown, supervisedUnknownCost) {
		unknown = append(unknown, supervisedUnknownCost)
	}
	body := SupervisedReport{
		Kind:               supervisedReportKind,
		ReportID:           input.ReportID,
		RecordedAt:         input.RecordedAt,
		RecordedBy:         evidence.Actor{Type: "system", ID: "proofrail"},
		RunID:              input.RunID,
		TaskID:             input.TaskID,
		StepID:             input.StepID,
		Attempt:            input.Attempt,
		Mode:               supervisedReportMode,
		AssuranceLevel:     supervisedAssuranceReduced,
		VerifiedDimensions: cloneStrings(input.VerifiedDimensions),
		UnknownDimensions:  unknown,
		Evidence:           cloneStrings(input.Evidence),
	}
	if body.ReportID == "" {
		body.ReportID = ids()
	}
	if body.RecordedAt == "" {
		body.RecordedAt = clock().UTC().Format(supervisedTimeLayout)
	}
	hash, err := digestRecord(supervisedReportDomain, body)
	if err != nil {
		return SupervisedReportRecord{}, err
	}
	record := SupervisedReportRecord{SchemaVersion: supervisedReportSchema, Record: body, RecordHash: hash}
	if err := ValidateSupervisedReport(record); err != nil {
		return SupervisedReportRecord{}, err
	}
	return record, nil
}

// DecodeSupervisedReport strictly decodes and validates one report.
func DecodeSupervisedReport(input []byte) (SupervisedReportRecord, error) {
	var record SupervisedReportRecord
	if err := evidence.DecodeStrictJSON(input, &record); err != nil {
		return SupervisedReportRecord{}, fmt.Errorf("%w: %v", ErrInvalidSupervisedReport, err)
	}
	if err := ValidateSupervisedReport(record); err != nil {
		return SupervisedReportRecord{}, err
	}
	return record, nil
}

// ValidateSupervisedReport checks the envelope and the body, and recomputes the
// canonical record hash. A report that claims full assurance, omits "cost", or
// lists a dimension outside the closed set is refused.
func ValidateSupervisedReport(record SupervisedReportRecord) error {
	if record.SchemaVersion != supervisedReportSchema || !evidence.ValidHash(record.RecordHash) {
		return fmt.Errorf("%w: schema version or record hash", ErrInvalidSupervisedReport)
	}
	if err := validateSupervisedReportBody(record.Record); err != nil {
		return err
	}
	expected, err := digestRecord(supervisedReportDomain, record.Record)
	if err != nil || expected != record.RecordHash {
		return fmt.Errorf("%w: recordHash mismatch", ErrInvalidSupervisedReport)
	}
	return nil
}

func validateSupervisedReportBody(body SupervisedReport) error {
	if body.Kind != supervisedReportKind {
		return fmt.Errorf("%w: invalid kind", ErrInvalidSupervisedReport)
	}
	if !evidence.ValidID(body.ReportID) || !evidence.ValidID(body.RunID) || !evidence.ValidID(body.TaskID) || !evidence.ValidID(body.StepID) || body.Attempt < 1 {
		return fmt.Errorf("%w: invalid binding identity", ErrInvalidSupervisedReport)
	}
	if body.RecordedBy.Type != "system" || !evidence.ValidID(body.RecordedBy.ID) {
		return fmt.Errorf("%w: invalid recording actor", ErrInvalidSupervisedReport)
	}
	if !validSupervisedTimestamp(body.RecordedAt) {
		return fmt.Errorf("%w: invalid record time", ErrInvalidSupervisedReport)
	}
	if body.Mode != supervisedReportMode {
		return fmt.Errorf("%w: invalid mode", ErrInvalidSupervisedReport)
	}
	if body.AssuranceLevel != supervisedAssuranceReduced {
		return fmt.Errorf("%w: assurance level is permanently reduced", ErrInvalidSupervisedReport)
	}
	if len(body.UnknownDimensions) == 0 || !slices.Contains(body.UnknownDimensions, supervisedUnknownCost) {
		return fmt.Errorf("%w: unknown dimensions must be non-empty and carry cost", ErrInvalidSupervisedReport)
	}
	for _, dimension := range body.UnknownDimensions {
		if !slices.Contains(supervisedUnknownDimensionSet, dimension) {
			return fmt.Errorf("%w: unknown dimension %q is outside the closed set", ErrInvalidSupervisedReport, dimension)
		}
	}
	if len(uniqueHashes(body.UnknownDimensions)) != len(body.UnknownDimensions) {
		return fmt.Errorf("%w: unknown dimensions must be unique", ErrInvalidSupervisedReport)
	}
	for _, dimension := range body.VerifiedDimensions {
		if !slices.Contains(supervisedVerifiedDimensionSet, dimension) {
			return fmt.Errorf("%w: verified dimension %q is outside the closed set", ErrInvalidSupervisedReport, dimension)
		}
	}
	if len(uniqueHashes(body.VerifiedDimensions)) != len(body.VerifiedDimensions) {
		return fmt.Errorf("%w: verified dimensions must be unique", ErrInvalidSupervisedReport)
	}
	if !validUniqueHashSet(body.Evidence, true) {
		return fmt.Errorf("%w: report evidence must be a non-empty unique hash set", ErrInvalidSupervisedReport)
	}
	return nil
}
