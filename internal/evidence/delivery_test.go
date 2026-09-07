package evidence

import (
	"errors"
	"testing"
)

func TestExportRecordCompletedRoundTrip(t *testing.T) {
	export := validDeliveryExportCompleted()
	record, err := NewExportRecord(export)
	if err != nil {
		t.Fatal(err)
	}
	if record.SchemaVersion != SchemaVersion || !ValidHash(record.ExportHash) {
		t.Fatalf("unexpected export record header: %+v", record)
	}
	if err := ValidateExportRecord(record); err != nil {
		t.Fatalf("validate export record: %v", err)
	}
}

func TestExportRecordCompletedRequiresEntriesAndManifest(t *testing.T) {
	export := validDeliveryExportCompleted()
	export.Entries = nil
	if _, err := NewExportRecord(export); !errors.Is(err, ErrInvalidExportRecord) {
		t.Fatalf("expected completed-without-entries rejection, got %v", err)
	}

	export = validDeliveryExportCompleted()
	export.PackageManifestHash = nil
	if _, err := NewExportRecord(export); !errors.Is(err, ErrInvalidExportRecord) {
		t.Fatalf("expected completed-without-package hash rejection, got %v", err)
	}
}

func TestExportRecordFailedRequiresErrorEvidence(t *testing.T) {
	export := validDeliveryExportCompleted()
	export.Outcome = "failed"
	export.Entries = nil
	export.PackageManifestHash = nil
	export.Evidence = nil
	export.ErrorEvidence = nil
	if _, err := NewExportRecord(export); !errors.Is(err, ErrInvalidExportRecord) {
		t.Fatalf("expected failed-without-errorEvidence rejection, got %v", err)
	}
}

func TestExportRecordRejectsDestinationPreconditionAndHashMismatch(t *testing.T) {
	export := validDeliveryExportCompleted()
	export.DestinationPrecondition = "present"
	if _, err := NewExportRecord(export); !errors.Is(err, ErrInvalidExportRecord) {
		t.Fatalf("expected destination precondition rejection, got %v", err)
	}

	record, err := NewExportRecord(validDeliveryExportCompleted())
	if err != nil {
		t.Fatal(err)
	}
	record.ExportHash = digest("tampered")
	if err := ValidateExportRecord(record); !errors.Is(err, ErrInvalidExportRecord) {
		t.Fatalf("expected export hash mismatch rejection, got %v", err)
	}
}

func validDeliveryExportCompleted() DeliveryExport {
	packageHash := digest("package-manifest")
	return DeliveryExport{
		ExportID:                "export-one",
		CreatedAt:               "2026-09-08T12:00:00.000Z",
		RecordedBy:              Actor{Type: "system", ID: "proofrail"},
		RunID:                   "run-one",
		TaskID:                  "task-one",
		Attempt:                 1,
		SourceSnapshotHash:      digest("snapshot"),
		ReviewReceiptHash:       digest("review"),
		PromotionReceiptHash:    digest("promotion"),
		VerificationReportHash:  digest("verification"),
		DestinationRef:          "delivery",
		DestinationIdentityHash: digest("destination"),
		DestinationPrecondition: "absent",
		Entries: []ExportEntry{{
			Path:        "bin/prfrail.exe",
			Type:        "regular-file",
			ContentHash: digest("binary"),
			SizeBytes:   1024,
		}},
		Exclusions: []ExportExclusion{{
			Category:     "cache",
			Pattern:      "tmp/cache",
			MatchedCount: 0,
		}},
		Outcome:             "completed",
		PackageManifestHash: &packageHash,
		Evidence:            []string{digest("evidence")},
		ErrorEvidence:       []string{},
	}
}

func digest(label string) string {
	return Digest("", []byte(label))
}
