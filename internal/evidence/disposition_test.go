package evidence

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLifecycleBackupRoundTrip(t *testing.T) {
	manifestHash := Digest("", []byte("backup-manifest"))
	record, err := NewLifecycleRecord(LifecycleBackup{
		RecordID:                "backup-one",
		Kind:                    "backup",
		CreatedAt:               "2026-09-08T12:00:00.000Z",
		RecordedBy:              Actor{Type: "system", ID: "proofrail"},
		SourceStoreHash:         Digest("", []byte("source-store")),
		SourceSchemaVersion:     SchemaVersion,
		DestinationRef:          "backup-store",
		DestinationIdentityHash: Digest("", []byte("destination")),
		WriterStopEvidence:      []string{Digest("", []byte("stopped"))},
		ObjectHashes:            []string{Digest("", []byte("object"))},
		EventSegmentHashes:      []string{Digest("", []byte("events"))},
		ReferenceRootHashes:     []string{Digest("", []byte("references"))},
		SecretDisposition:       "excluded",
		SecretScanEvidence:      []string{Digest("", []byte("scan"))},
		ClosureStatus:           "complete",
		VerificationReportHash:  Digest("", []byte("verification")),
		RestoreEvidence:         []string{Digest("", []byte("restore"))},
		Outcome:                 "completed",
		BackupManifestHash:      &manifestHash,
		ErrorEvidence:           []string{},
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeCanonical(record)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeLifecycleRecord(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.RecordHash != record.RecordHash {
		t.Fatal("lifecycle record hash changed after round trip")
	}
}

func TestRetirementDefaultsRetainAndPreserveSharedTools(t *testing.T) {
	record, err := NewRetirementLifecycleRecord(validRetirement(), false)
	if err != nil {
		t.Fatal(err)
	}
	retirement := record.Record.(LifecycleRetirement)
	if retirement.RetentionDisposition != "retain" || retirement.SharedToolsDisposition != "preserved" {
		t.Fatalf("unsafe retirement defaults: %+v", retirement)
	}
}

func TestRetirementDeletionAndConflictRequireHumanEvidence(t *testing.T) {
	retirement := validRetirement()
	retirement.RetentionDisposition = "delete-authorized"
	if _, err := NewRetirementLifecycleRecord(retirement, false); !errors.Is(err, ErrInvalidLifecycleRecord) {
		t.Fatalf("expected missing deletion authorization rejection, got %v", err)
	}

	deleteAuthorization := Digest("", []byte("human-delete-authorization"))
	retirement.DeletionAuthorizationHash = &deleteAuthorization
	if _, err := NewRetirementLifecycleRecord(retirement, true); !errors.Is(err, ErrInvalidLifecycleRecord) {
		t.Fatalf("expected unresolved retention conflict rejection, got %v", err)
	}
	auditDecision := Digest("", []byte("human-audit-decision"))
	retirement.AuditConflictDecisionHash = &auditDecision
	if _, err := NewRetirementLifecycleRecord(retirement, true); err != nil {
		t.Fatalf("expected authorized deletion with conflict decision, got %v", err)
	}
}

func TestReleaseRequiresSignatureAndRepresentsOfflineRevocationAsUnknown(t *testing.T) {
	release := validRelease()
	release.SignatureReceiptHash = ""
	if _, err := NewLifecycleRecord(release); !errors.Is(err, ErrInvalidLifecycleRecord) {
		t.Fatalf("expected unsigned release rejection, got %v", err)
	}

	release = validRelease()
	release.RevocationFreshness = "unknown"
	release.RevocationEvidence = []string{}
	if _, err := NewLifecycleRecord(release); err != nil {
		t.Fatalf("offline unknown revocation should be representable: %v", err)
	}
	release.RevocationEvidence = []string{Digest("", []byte("stale-evidence"))}
	if _, err := NewLifecycleRecord(release); !errors.Is(err, ErrInvalidLifecycleRecord) {
		t.Fatalf("unknown revocation must not claim verification evidence: %v", err)
	}
}

func TestBuildReleaseInventoryIsDeterministicAndRejectsLinks(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "prfrail.exe"), []byte("binary"), 0644); err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
	first, err := BuildReleaseInventory(root, now)
	if err != nil {
		t.Fatal(err)
	}
	second, err := BuildReleaseInventory(root, now.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if first.InventoryHash != second.InventoryHash || len(first.Entries) != 1 {
		t.Fatalf("inventory content hash must be stable: %+v %+v", first, second)
	}

	linkPath := filepath.Join(root, "binary-link")
	if err := os.Symlink(filepath.Join("bin", "prfrail.exe"), linkPath); err == nil {
		if _, err := BuildReleaseInventory(root, now); !errors.Is(err, ErrLifecycleUnsupportedLink) {
			t.Fatalf("expected link rejection, got %v", err)
		}
	}
}

func validRetirement() LifecycleRetirement {
	return LifecycleRetirement{
		RecordID:                  "retirement-one",
		Kind:                      "retirement",
		CreatedAt:                 "2026-09-08T12:00:00.000Z",
		RecordedBy:                Actor{Type: "operator", ID: "owner"},
		Action:                    "uninstall",
		InstallationHash:          Digest("", []byte("installation")),
		ProcessStopEvidence:       []string{Digest("", []byte("process-stop"))},
		AdapterRevocationEvidence: []string{Digest("", []byte("adapter-revocation"))},
		InventoryHash:             Digest("", []byte("inventory")),
		SecretIncidentEvidence:    []string{},
		Outcome:                   "completed",
		ErrorEvidence:             []string{},
	}
}

func validRelease() LifecycleRelease {
	return LifecycleRelease{
		RecordID:             "release-one",
		Kind:                 "release",
		CreatedAt:            "2026-09-08T12:00:00.000Z",
		Version:              "0.1.0",
		BinaryHashes:         []string{Digest("", []byte("binary"))},
		Dependencies:         []LifecycleDependency{},
		SBOMHash:             Digest("", []byte("sbom")),
		LicenseManifestHash:  Digest("", []byte("licenses")),
		ChecksumManifestHash: Digest("", []byte("checksums")),
		SignatureReceiptHash: Digest("", []byte("signature")),
		SupportMatrix: []LifecycleSupport{{
			ComponentID:      "prfrail",
			VersionRange:     "0.1.x",
			Status:           "supported",
			SupportExpiresAt: nil,
			CapabilityIDs:    []string{"windows-cli"},
			Platforms:        []LifecyclePlatform{{OS: "windows", Arch: "amd64"}},
		}},
		KnownLimitationIDs:  []string{"no-production-package"},
		RevocationFreshness: "verified",
		RevocationEvidence:  []string{Digest("", []byte("revocation-check"))},
	}
}
