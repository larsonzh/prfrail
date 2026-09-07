package evidence

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	exportRecordDomain    = "proofrail:export-record:1\n"
	exportTimestampLayout = "2006-01-02T15:04:05.000Z"
)

var ErrInvalidExportRecord = errors.New("invalid export record")

type ExportRecord struct {
	SchemaVersion string         `json:"schemaVersion"`
	Export        DeliveryExport `json:"export"`
	ExportHash    string         `json:"exportHash"`
}

type DeliveryExport struct {
	ExportID                string            `json:"exportId"`
	CreatedAt               string            `json:"createdAt"`
	RecordedBy              Actor             `json:"recordedBy"`
	RunID                   string            `json:"runId"`
	TaskID                  string            `json:"taskId"`
	Attempt                 int               `json:"attempt"`
	SourceSnapshotHash      string            `json:"sourceSnapshotHash"`
	ReviewReceiptHash       string            `json:"reviewReceiptHash"`
	PromotionReceiptHash    string            `json:"promotionReceiptHash"`
	VerificationReportHash  string            `json:"verificationReportHash"`
	DestinationRef          string            `json:"destinationRef"`
	DestinationIdentityHash string            `json:"destinationIdentityHash"`
	DestinationPrecondition string            `json:"destinationPrecondition"`
	Entries                 []ExportEntry     `json:"entries"`
	Exclusions              []ExportExclusion `json:"exclusions"`
	Outcome                 string            `json:"outcome"`
	PackageManifestHash     *string           `json:"packageManifestHash"`
	Evidence                []string          `json:"evidence"`
	ErrorEvidence           []string          `json:"errorEvidence"`
}

type ExportEntry struct {
	Path        string `json:"path"`
	Type        string `json:"type"`
	ContentHash string `json:"contentHash"`
	SizeBytes   int64  `json:"sizeBytes"`
}

type ExportExclusion struct {
	Category     string `json:"category"`
	Pattern      string `json:"pattern"`
	MatchedCount int    `json:"matchedCount"`
}

func NewExportRecord(export DeliveryExport) (ExportRecord, error) {
	if err := validateDeliveryExport(export); err != nil {
		return ExportRecord{}, err
	}
	canonical, err := EncodeCanonical(export)
	if err != nil {
		return ExportRecord{}, fmt.Errorf("%w: canonicalize export: %v", ErrInvalidExportRecord, err)
	}
	return ExportRecord{
		SchemaVersion: SchemaVersion,
		Export:        export,
		ExportHash:    Digest(exportRecordDomain, canonical),
	}, nil
}

func DecodeExportRecord(input []byte) (ExportRecord, error) {
	var record ExportRecord
	if err := decodeStrictJSON(input, &record); err != nil {
		return ExportRecord{}, err
	}
	if err := ValidateExportRecord(record); err != nil {
		return ExportRecord{}, err
	}
	return record, nil
}

func ValidateExportRecord(record ExportRecord) error {
	if record.SchemaVersion != SchemaVersion {
		return fmt.Errorf("%w: unsupported schemaVersion %q", ErrInvalidExportRecord, record.SchemaVersion)
	}
	if err := validateDeliveryExport(record.Export); err != nil {
		return err
	}
	canonical, err := EncodeCanonical(record.Export)
	if err != nil {
		return fmt.Errorf("%w: canonicalize export: %v", ErrInvalidExportRecord, err)
	}
	expected := Digest(exportRecordDomain, canonical)
	if record.ExportHash != expected {
		return fmt.Errorf("%w: export hash mismatch", ErrInvalidExportRecord)
	}
	return nil
}

func validateDeliveryExport(export DeliveryExport) error {
	if !validID(export.ExportID) || !validID(export.RunID) || !validID(export.TaskID) || !validID(export.DestinationRef) || export.Attempt < 1 {
		return fmt.Errorf("%w: invalid export identity", ErrInvalidExportRecord)
	}
	if err := validateExportTimestamp(export.CreatedAt); err != nil {
		return err
	}
	if export.RecordedBy.Type != "operator" && export.RecordedBy.Type != "system" {
		return fmt.Errorf("%w: invalid recordedBy.type", ErrInvalidExportRecord)
	}
	if !validID(export.RecordedBy.ID) {
		return fmt.Errorf("%w: invalid recordedBy.id", ErrInvalidExportRecord)
	}
	hashes := []string{
		export.SourceSnapshotHash,
		export.ReviewReceiptHash,
		export.PromotionReceiptHash,
		export.VerificationReportHash,
		export.DestinationIdentityHash,
	}
	for _, hash := range hashes {
		if !validHash(hash) {
			return fmt.Errorf("%w: invalid hash field", ErrInvalidExportRecord)
		}
	}
	if export.DestinationPrecondition != "absent" {
		return fmt.Errorf("%w: destinationPrecondition must be absent", ErrInvalidExportRecord)
	}
	if err := validateExportEntries(export.Entries); err != nil {
		return err
	}
	if err := validateExportExclusions(export.Exclusions); err != nil {
		return err
	}
	if err := validateHashSet(export.Evidence, "evidence"); err != nil {
		return err
	}
	if err := validateHashSet(export.ErrorEvidence, "errorEvidence"); err != nil {
		return err
	}
	switch export.Outcome {
	case "completed":
		if len(export.Entries) == 0 {
			return fmt.Errorf("%w: completed export requires entries", ErrInvalidExportRecord)
		}
		if export.PackageManifestHash == nil || !validHash(*export.PackageManifestHash) {
			return fmt.Errorf("%w: completed export requires packageManifestHash", ErrInvalidExportRecord)
		}
		if len(export.Evidence) == 0 || len(export.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: completed export requires evidence and empty errorEvidence", ErrInvalidExportRecord)
		}
	case "failed", "uncertain":
		if len(export.Entries) != 0 {
			return fmt.Errorf("%w: failed/uncertain export must not include entries", ErrInvalidExportRecord)
		}
		if export.PackageManifestHash != nil {
			return fmt.Errorf("%w: failed/uncertain export must not include packageManifestHash", ErrInvalidExportRecord)
		}
		if len(export.ErrorEvidence) == 0 {
			return fmt.Errorf("%w: failed/uncertain export requires errorEvidence", ErrInvalidExportRecord)
		}
	default:
		return fmt.Errorf("%w: invalid export outcome %q", ErrInvalidExportRecord, export.Outcome)
	}
	return nil
}

func validateExportTimestamp(value string) error {
	parsed, err := time.Parse(exportTimestampLayout, value)
	if err != nil {
		return fmt.Errorf("%w: invalid createdAt", ErrInvalidExportRecord)
	}
	if parsed.UTC().Format(exportTimestampLayout) != value {
		return fmt.Errorf("%w: createdAt must be UTC millisecond timestamp", ErrInvalidExportRecord)
	}
	return nil
}

func validateExportEntries(entries []ExportEntry) error {
	seenPath := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		if !validRelativePath(entry.Path) {
			return fmt.Errorf("%w: invalid entry path %q", ErrInvalidExportRecord, entry.Path)
		}
		switch entry.Type {
		case "regular-file", "directory", "symbolic-link":
		default:
			return fmt.Errorf("%w: invalid entry type %q", ErrInvalidExportRecord, entry.Type)
		}
		if !validHash(entry.ContentHash) {
			return fmt.Errorf("%w: invalid entry contentHash", ErrInvalidExportRecord)
		}
		if entry.SizeBytes < 0 {
			return fmt.Errorf("%w: entry sizeBytes must be >= 0", ErrInvalidExportRecord)
		}
		if _, exists := seenPath[entry.Path]; exists {
			return fmt.Errorf("%w: duplicate entry path %q", ErrInvalidExportRecord, entry.Path)
		}
		seenPath[entry.Path] = struct{}{}
	}
	return nil
}

func validateExportExclusions(exclusions []ExportExclusion) error {
	type key struct {
		category string
		pattern  string
	}
	seen := make(map[key]struct{}, len(exclusions))
	for _, exclusion := range exclusions {
		switch exclusion.Category {
		case "secret", "cache", "machine-environment", "unsupported-type":
		default:
			return fmt.Errorf("%w: invalid exclusion category %q", ErrInvalidExportRecord, exclusion.Category)
		}
		if !validRelativePath(exclusion.Pattern) {
			return fmt.Errorf("%w: invalid exclusion pattern %q", ErrInvalidExportRecord, exclusion.Pattern)
		}
		if exclusion.MatchedCount < 0 {
			return fmt.Errorf("%w: exclusion matchedCount must be >= 0", ErrInvalidExportRecord)
		}
		k := key{category: exclusion.Category, pattern: exclusion.Pattern}
		if _, exists := seen[k]; exists {
			return fmt.Errorf("%w: duplicate exclusion %s/%s", ErrInvalidExportRecord, exclusion.Category, exclusion.Pattern)
		}
		seen[k] = struct{}{}
	}
	return nil
}

func validateHashSet(values []string, field string) error {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !validHash(value) {
			return fmt.Errorf("%w: invalid %s hash %q", ErrInvalidExportRecord, field, value)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%w: duplicate %s hash %q", ErrInvalidExportRecord, field, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validRelativePath(path string) bool {
	if len(path) == 0 || len(path) > 1024 {
		return false
	}
	if strings.HasPrefix(path, "/") || strings.Contains(path, "\\") || strings.Contains(path, "\x00") {
		return false
	}
	if len(path) >= 2 && isASCIIAlpha(path[0]) && path[1] == ':' {
		return false
	}
	colon := strings.Index(path, ":")
	if colon > 0 && !strings.Contains(path[:colon], "/") && looksLikeScheme(path[:colon]) {
		return false
	}
	segments := strings.Split(path, "/")
	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return false
		}
	}
	return true
}

func isASCIIAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func looksLikeScheme(value string) bool {
	if value == "" || !isASCIIAlpha(value[0]) {
		return false
	}
	for i := 1; i < len(value); i++ {
		ch := value[i]
		if !isASCIIAlpha(ch) && !(ch >= '0' && ch <= '9') && ch != '+' && ch != '.' && ch != '-' {
			return false
		}
	}
	return true
}
