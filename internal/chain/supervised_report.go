package chain

import (
	"fmt"
	"slices"
	"strings"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// The reduced/unknown report is the factual record of what a supervised
// black-box run can and cannot prove (contracts section 7.1d). The assurance
// level is permanently "reduced": it is the ceiling of the mode, never a
// degraded state. The unknown dimension set must carry every member of the
// frozen closed set (tool usage, resource reads, network access, cost,
// background processes, session continuity and external side effects, ruling
// R-09.8), so no unobservable dimension can be dropped from the record and no
// unknown dimension can ever be rendered as safe.
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
// always the reduced constant and can never be raised. An empty unknown
// dimension set falls back to the full closed set; a caller-supplied set must
// carry all seven required dimensions, so a partial set is refused instead of
// being silently completed.
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
// canonical record hash. A report that claims full assurance, omits any of the
// seven required unknown dimensions, or lists a dimension outside the closed
// set is refused.
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
	if len(body.UnknownDimensions) == 0 {
		return fmt.Errorf("%w: unknown dimensions must be non-empty", ErrInvalidSupervisedReport)
	}
	if !slices.Contains(body.UnknownDimensions, supervisedUnknownCost) {
		return fmt.Errorf("%w: unknown dimensions must carry cost", ErrInvalidSupervisedReport)
	}
	for _, dimension := range body.UnknownDimensions {
		if !slices.Contains(supervisedUnknownDimensionSet, dimension) {
			return fmt.Errorf("%w: unknown dimension %q is outside the closed set", ErrInvalidSupervisedReport, dimension)
		}
	}
	if len(uniqueHashes(body.UnknownDimensions)) != len(body.UnknownDimensions) {
		return fmt.Errorf("%w: unknown dimensions must be unique", ErrInvalidSupervisedReport)
	}
	if missing := missingSupervisedUnknownDimensions(body.UnknownDimensions); len(missing) > 0 {
		return fmt.Errorf("%w: unknown dimensions must carry all seven required dimensions (missing: %s)", ErrInvalidSupervisedReport, strings.Join(missing, ", "))
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

// missingSupervisedUnknownDimensions returns the required dimensions that are
// absent from the supplied set, in the canonical closed-set order, so the
// rejection message names exactly what the report failed to record.
func missingSupervisedUnknownDimensions(dimensions []string) []string {
	missing := make([]string, 0, len(supervisedUnknownDimensionSet))
	for _, dimension := range supervisedUnknownDimensionSet {
		if !slices.Contains(dimensions, dimension) {
			missing = append(missing, dimension)
		}
	}
	return missing
}
