package evidence

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	lifecycleRecordDomain    = "proofrail:lifecycle-record:1\n"
	lifecycleInventoryDomain = "proofrail:release-inventory:1\n"
	lifecycleTimestampLayout = "2006-01-02T15:04:05.000Z"
)

var (
	ErrInvalidLifecycleRecord   = errors.New("invalid lifecycle record")
	ErrInvalidReleaseInventory  = errors.New("invalid release inventory")
	ErrLifecycleUnknownKind     = errors.New("unknown lifecycle record kind")
	ErrLifecycleUnsupportedLink = errors.New("unsupported inventory link")
)

type LifecycleRecord struct {
	SchemaVersion string `json:"schemaVersion"`
	Record        any    `json:"record"`
	RecordHash    string `json:"recordHash"`
}

type LifecycleDependency struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	ContentHash     string `json:"contentHash"`
	LicenseEvidence string `json:"licenseEvidence"`
}

type LifecyclePlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type LifecycleSupport struct {
	ComponentID      string              `json:"componentId"`
	VersionRange     string              `json:"versionRange"`
	Status           string              `json:"status"`
	SupportExpiresAt *string             `json:"supportExpiresAt"`
	CapabilityIDs    []string            `json:"capabilityIds"`
	Platforms        []LifecyclePlatform `json:"platforms"`
}

type LifecycleRelease struct {
	RecordID             string                `json:"recordId"`
	Kind                 string                `json:"kind"`
	CreatedAt            string                `json:"createdAt"`
	Version              string                `json:"version"`
	BinaryHashes         []string              `json:"binaryHashes"`
	Dependencies         []LifecycleDependency `json:"dependencies"`
	SBOMHash             string                `json:"sbomHash"`
	LicenseManifestHash  string                `json:"licenseManifestHash"`
	ChecksumManifestHash string                `json:"checksumManifestHash"`
	SignatureReceiptHash *string               `json:"signatureReceiptHash"`
	SupportMatrix        []LifecycleSupport    `json:"supportMatrix"`
	KnownLimitationIDs   []string              `json:"knownLimitationIds"`
	RevocationFreshness  string                `json:"revocationFreshness"`
	RevocationEvidence   []string              `json:"revocationEvidence"`
}

type LifecycleBackup struct {
	RecordID                string   `json:"recordId"`
	Kind                    string   `json:"kind"`
	CreatedAt               string   `json:"createdAt"`
	RecordedBy              Actor    `json:"recordedBy"`
	SourceStoreHash         string   `json:"sourceStoreHash"`
	SourceSchemaVersion     string   `json:"sourceSchemaVersion"`
	DestinationRef          string   `json:"destinationRef"`
	DestinationIdentityHash string   `json:"destinationIdentityHash"`
	WriterStopEvidence      []string `json:"writerStopEvidence"`
	ObjectHashes            []string `json:"objectHashes"`
	EventSegmentHashes      []string `json:"eventSegmentHashes"`
	ReferenceRootHashes     []string `json:"referenceRootHashes"`
	SecretDisposition       string   `json:"secretDisposition"`
	SecretScanEvidence      []string `json:"secretScanEvidence"`
	ClosureStatus           string   `json:"closureStatus"`
	VerificationReportHash  string   `json:"verificationReportHash"`
	RestoreEvidence         []string `json:"restoreEvidence"`
	Outcome                 string   `json:"outcome"`
	BackupManifestHash      *string  `json:"backupManifestHash"`
	ErrorEvidence           []string `json:"errorEvidence"`
}

type LifecycleRetirement struct {
	RecordID                  string   `json:"recordId"`
	Kind                      string   `json:"kind"`
	CreatedAt                 string   `json:"createdAt"`
	RecordedBy                Actor    `json:"recordedBy"`
	Action                    string   `json:"action"`
	InstallationHash          string   `json:"installationHash"`
	ProcessStopEvidence       []string `json:"processStopEvidence"`
	AdapterRevocationEvidence []string `json:"adapterRevocationEvidence"`
	InventoryHash             string   `json:"inventoryHash"`
	RetentionDisposition      string   `json:"retentionDisposition"`
	DeletionAuthorizationHash *string  `json:"deletionAuthorizationHash"`
	AuditConflictDecisionHash *string  `json:"auditConflictDecisionHash"`
	SharedToolsDisposition    string   `json:"sharedToolsDisposition"`
	SecretIncidentEvidence    []string `json:"secretIncidentEvidence"`
	Outcome                   string   `json:"outcome"`
	ErrorEvidence             []string `json:"errorEvidence"`
}

type ReleaseInventoryEntry struct {
	Path        string `json:"path"`
	Type        string `json:"type"`
	ContentHash string `json:"contentHash"`
	SizeBytes   int64  `json:"sizeBytes"`
}

type ReleaseInventory struct {
	RootPath      string                  `json:"rootPath"`
	GeneratedAt   string                  `json:"generatedAt"`
	Entries       []ReleaseInventoryEntry `json:"entries"`
	InventoryHash string                  `json:"inventoryHash"`
}

func NewLifecycleRecord(record any) (LifecycleRecord, error) {
	normalized, err := normalizeLifecycleRecord(record)
	if err != nil {
		return LifecycleRecord{}, err
	}
	if err := validateLifecycleRecordBody(normalized); err != nil {
		return LifecycleRecord{}, err
	}
	canonical, err := canonicalValue(normalized)
	if err != nil {
		return LifecycleRecord{}, fmt.Errorf("%w: canonicalize record: %v", ErrInvalidLifecycleRecord, err)
	}
	return LifecycleRecord{
		SchemaVersion: SchemaVersion,
		Record:        normalized,
		RecordHash:    Digest(lifecycleRecordDomain, canonical),
	}, nil
}

func NewRetirementLifecycleRecord(record LifecycleRetirement, retentionConflict bool) (LifecycleRecord, error) {
	if record.RetentionDisposition == "" {
		record.RetentionDisposition = "retain"
	}
	if record.SharedToolsDisposition == "" {
		record.SharedToolsDisposition = "preserved"
	}
	if retentionConflict && record.AuditConflictDecisionHash == nil {
		return LifecycleRecord{}, fmt.Errorf("%w: retention conflict requires auditConflictDecisionHash", ErrInvalidLifecycleRecord)
	}
	return NewLifecycleRecord(record)
}

func DecodeLifecycleRecord(input []byte) (LifecycleRecord, error) {
	type lifecycleRecordWire struct {
		SchemaVersion string          `json:"schemaVersion"`
		Record        json.RawMessage `json:"record"`
		RecordHash    string          `json:"recordHash"`
	}

	var wire lifecycleRecordWire
	if err := decodeStrictJSON(input, &wire); err != nil {
		return LifecycleRecord{}, err
	}
	if wire.SchemaVersion != SchemaVersion {
		return LifecycleRecord{}, fmt.Errorf("%w: unsupported schemaVersion %q", ErrInvalidLifecycleRecord, wire.SchemaVersion)
	}

	var kindProbe struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(wire.Record, &kindProbe); err != nil {
		return LifecycleRecord{}, fmt.Errorf("%w: decode lifecycle kind: %v", ErrInvalidLifecycleRecord, err)
	}

	var body any
	switch kindProbe.Kind {
	case "release":
		var release LifecycleRelease
		if err := decodeStrictJSON(wire.Record, &release); err != nil {
			return LifecycleRecord{}, err
		}
		body = release
	case "backup":
		var backup LifecycleBackup
		if err := decodeStrictJSON(wire.Record, &backup); err != nil {
			return LifecycleRecord{}, err
		}
		body = backup
	case "retirement":
		var retirement LifecycleRetirement
		if err := decodeStrictJSON(wire.Record, &retirement); err != nil {
			return LifecycleRecord{}, err
		}
		body = retirement
	default:
		return LifecycleRecord{}, fmt.Errorf("%w: %q", ErrLifecycleUnknownKind, kindProbe.Kind)
	}

	record := LifecycleRecord{SchemaVersion: wire.SchemaVersion, Record: body, RecordHash: wire.RecordHash}
	if err := ValidateLifecycleRecord(record); err != nil {
		return LifecycleRecord{}, err
	}
	return record, nil
}

func ValidateLifecycleRecord(record LifecycleRecord) error {
	if record.SchemaVersion != SchemaVersion {
		return fmt.Errorf("%w: unsupported schemaVersion %q", ErrInvalidLifecycleRecord, record.SchemaVersion)
	}
	if !validHash(record.RecordHash) {
		return fmt.Errorf("%w: invalid recordHash", ErrInvalidLifecycleRecord)
	}
	normalized, err := normalizeLifecycleRecord(record.Record)
	if err != nil {
		return err
	}
	if err := validateLifecycleRecordBody(normalized); err != nil {
		return err
	}
	canonical, err := canonicalValue(normalized)
	if err != nil {
		return fmt.Errorf("%w: canonicalize record: %v", ErrInvalidLifecycleRecord, err)
	}
	if expected := Digest(lifecycleRecordDomain, canonical); expected != record.RecordHash {
		return fmt.Errorf("%w: record hash mismatch", ErrInvalidLifecycleRecord)
	}
	return nil
}

func BuildReleaseInventory(rootPath string, generatedAt time.Time) (ReleaseInventory, error) {
	if strings.TrimSpace(rootPath) == "" {
		return ReleaseInventory{}, fmt.Errorf("%w: empty rootPath", ErrInvalidReleaseInventory)
	}
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}
	rootAbs, err := filepath.Abs(rootPath)
	if err != nil {
		return ReleaseInventory{}, err
	}
	rootInfo, err := os.Stat(rootAbs)
	if err != nil {
		return ReleaseInventory{}, err
	}
	if !rootInfo.IsDir() {
		return ReleaseInventory{}, fmt.Errorf("%w: rootPath must be a directory", ErrInvalidReleaseInventory)
	}

	entries := make([]ReleaseInventoryEntry, 0)
	err = filepath.WalkDir(rootAbs, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == rootAbs {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("%w: %q", ErrLifecycleUnsupportedLink, path)
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("%w: unsupported file type at %q", ErrInvalidReleaseInventory, path)
		}

		rel, err := filepath.Rel(rootAbs, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if !validRelativePath(rel) {
			return fmt.Errorf("%w: invalid inventory path %q", ErrInvalidReleaseInventory, rel)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		entries = append(entries, ReleaseInventoryEntry{
			Path:        rel,
			Type:        "regular-file",
			ContentHash: Digest("", data),
			SizeBytes:   int64(len(data)),
		})
		return nil
	})
	if err != nil {
		return ReleaseInventory{}, err
	}
	if len(entries) == 0 {
		return ReleaseInventory{}, fmt.Errorf("%w: inventory has no files", ErrInvalidReleaseInventory)
	}

	sort.Slice(entries, func(left, right int) bool {
		return entries[left].Path < entries[right].Path
	})
	seen := make(map[string]struct{}, len(entries))
	for _, item := range entries {
		if _, exists := seen[item.Path]; exists {
			return ReleaseInventory{}, fmt.Errorf("%w: duplicate inventory path %q", ErrInvalidReleaseInventory, item.Path)
		}
		seen[item.Path] = struct{}{}
	}

	payload := struct {
		Entries []ReleaseInventoryEntry `json:"entries"`
	}{Entries: entries}
	canonical, err := EncodeCanonical(payload)
	if err != nil {
		return ReleaseInventory{}, fmt.Errorf("%w: canonicalize inventory: %v", ErrInvalidReleaseInventory, err)
	}

	return ReleaseInventory{
		RootPath:      rootAbs,
		GeneratedAt:   generatedAt.UTC().Format(lifecycleTimestampLayout),
		Entries:       entries,
		InventoryHash: Digest(lifecycleInventoryDomain, canonical),
	}, nil
}

func normalizeLifecycleRecord(record any) (any, error) {
	switch typed := record.(type) {
	case LifecycleRelease:
		typed.BinaryHashes = cloneLifecycleStrings(typed.BinaryHashes)
		typed.Dependencies = cloneDependencies(typed.Dependencies)
		typed.SupportMatrix = cloneSupports(typed.SupportMatrix)
		typed.KnownLimitationIDs = cloneLifecycleStrings(typed.KnownLimitationIDs)
		typed.RevocationEvidence = cloneLifecycleStrings(typed.RevocationEvidence)
		return typed, nil
	case *LifecycleRelease:
		if typed == nil {
			return nil, fmt.Errorf("%w: nil release record", ErrInvalidLifecycleRecord)
		}
		return normalizeLifecycleRecord(*typed)
	case LifecycleBackup:
		typed.WriterStopEvidence = cloneLifecycleStrings(typed.WriterStopEvidence)
		typed.ObjectHashes = cloneLifecycleStrings(typed.ObjectHashes)
		typed.EventSegmentHashes = cloneLifecycleStrings(typed.EventSegmentHashes)
		typed.ReferenceRootHashes = cloneLifecycleStrings(typed.ReferenceRootHashes)
		typed.SecretScanEvidence = cloneLifecycleStrings(typed.SecretScanEvidence)
		typed.RestoreEvidence = cloneLifecycleStrings(typed.RestoreEvidence)
		typed.ErrorEvidence = cloneLifecycleStrings(typed.ErrorEvidence)
		return typed, nil
	case *LifecycleBackup:
		if typed == nil {
			return nil, fmt.Errorf("%w: nil backup record", ErrInvalidLifecycleRecord)
		}
		return normalizeLifecycleRecord(*typed)
	case LifecycleRetirement:
		typed.ProcessStopEvidence = cloneLifecycleStrings(typed.ProcessStopEvidence)
		typed.AdapterRevocationEvidence = cloneLifecycleStrings(typed.AdapterRevocationEvidence)
		typed.SecretIncidentEvidence = cloneLifecycleStrings(typed.SecretIncidentEvidence)
		typed.ErrorEvidence = cloneLifecycleStrings(typed.ErrorEvidence)
		return typed, nil
	case *LifecycleRetirement:
		if typed == nil {
			return nil, fmt.Errorf("%w: nil retirement record", ErrInvalidLifecycleRecord)
		}
		return normalizeLifecycleRecord(*typed)
	default:
		return nil, fmt.Errorf("%w: unsupported record type %T", ErrInvalidLifecycleRecord, record)
	}
}

func validateLifecycleRecordBody(record any) error {
	switch typed := record.(type) {
	case LifecycleRelease:
		return validateLifecycleRelease(typed)
	case LifecycleBackup:
		return validateLifecycleBackup(typed)
	case LifecycleRetirement:
		return validateLifecycleRetirement(typed)
	default:
		return fmt.Errorf("%w: unsupported record type %T", ErrInvalidLifecycleRecord, record)
	}
}

func validateLifecycleRelease(record LifecycleRelease) error {
	if !validID(record.RecordID) || record.Kind != "release" {
		return fmt.Errorf("%w: invalid release identity", ErrInvalidLifecycleRecord)
	}
	if err := validateLifecycleTimestamp(record.CreatedAt, "createdAt"); err != nil {
		return err
	}
	if !validLifecycleVersion(record.Version) {
		return fmt.Errorf("%w: invalid release version", ErrInvalidLifecycleRecord)
	}
	if err := validateLifecycleHashSet(record.BinaryHashes, "binaryHashes", 1); err != nil {
		return err
	}
	if !validHash(record.SBOMHash) || !validHash(record.LicenseManifestHash) || !validHash(record.ChecksumManifestHash) {
		return fmt.Errorf("%w: invalid release hash field", ErrInvalidLifecycleRecord)
	}
	if record.SignatureReceiptHash != nil && !validHash(*record.SignatureReceiptHash) {
		return fmt.Errorf("%w: invalid signatureReceiptHash", ErrInvalidLifecycleRecord)
	}
	if len(record.SupportMatrix) == 0 {
		return fmt.Errorf("%w: supportMatrix must not be empty", ErrInvalidLifecycleRecord)
	}
	for _, dep := range record.Dependencies {
		if strings.TrimSpace(dep.Name) == "" || len(dep.Name) > 256 {
			return fmt.Errorf("%w: invalid dependency name", ErrInvalidLifecycleRecord)
		}
		if !validLifecycleVersion(dep.Version) || !validHash(dep.ContentHash) || !validHash(dep.LicenseEvidence) {
			return fmt.Errorf("%w: invalid dependency field", ErrInvalidLifecycleRecord)
		}
	}
	for _, support := range record.SupportMatrix {
		if !validID(support.ComponentID) || !validLifecycleVersion(support.VersionRange) {
			return fmt.Errorf("%w: invalid support identity", ErrInvalidLifecycleRecord)
		}
		switch support.Status {
		case "supported", "deprecated", "eol":
		default:
			return fmt.Errorf("%w: invalid support status %q", ErrInvalidLifecycleRecord, support.Status)
		}
		if support.SupportExpiresAt != nil {
			if err := validateLifecycleTimestamp(*support.SupportExpiresAt, "supportExpiresAt"); err != nil {
				return err
			}
		}
		if err := validateLifecycleIDSet(support.CapabilityIDs, "capabilityIds", 0); err != nil {
			return err
		}
		if len(support.Platforms) == 0 {
			return fmt.Errorf("%w: support platforms must not be empty", ErrInvalidLifecycleRecord)
		}
		seenPlatforms := make(map[string]struct{}, len(support.Platforms))
		for _, platform := range support.Platforms {
			if !validID(platform.OS) || !validID(platform.Arch) {
				return fmt.Errorf("%w: invalid support platform", ErrInvalidLifecycleRecord)
			}
			key := platform.OS + "\n" + platform.Arch
			if _, exists := seenPlatforms[key]; exists {
				return fmt.Errorf("%w: duplicate support platform", ErrInvalidLifecycleRecord)
			}
			seenPlatforms[key] = struct{}{}
		}
	}
	if err := validateLifecycleIDSet(record.KnownLimitationIDs, "knownLimitationIds", 0); err != nil {
		return err
	}
	if err := validateLifecycleHashSet(record.RevocationEvidence, "revocationEvidence", 0); err != nil {
		return err
	}
	switch record.RevocationFreshness {
	case "verified":
		if record.SignatureReceiptHash == nil || len(record.RevocationEvidence) == 0 {
			return fmt.Errorf("%w: verified revocation requires a signature receipt and evidence", ErrInvalidLifecycleRecord)
		}
	case "unknown":
		if len(record.RevocationEvidence) != 0 {
			return fmt.Errorf("%w: unknown revocation freshness cannot include evidence", ErrInvalidLifecycleRecord)
		}
	default:
		return fmt.Errorf("%w: invalid revocationFreshness %q", ErrInvalidLifecycleRecord, record.RevocationFreshness)
	}
	return nil
}

func validateLifecycleBackup(record LifecycleBackup) error {
	if !validID(record.RecordID) || record.Kind != "backup" {
		return fmt.Errorf("%w: invalid backup identity", ErrInvalidLifecycleRecord)
	}
	if err := validateLifecycleTimestamp(record.CreatedAt, "createdAt"); err != nil {
		return err
	}
	if !validLifecycleActor(record.RecordedBy) {
		return fmt.Errorf("%w: invalid recordedBy actor", ErrInvalidLifecycleRecord)
	}
	if !validHash(record.SourceStoreHash) || record.SourceSchemaVersion != SchemaVersion || !validID(record.DestinationRef) || !validHash(record.DestinationIdentityHash) || !validHash(record.VerificationReportHash) {
		return fmt.Errorf("%w: invalid backup binding fields", ErrInvalidLifecycleRecord)
	}
	if err := validateLifecycleHashSet(record.WriterStopEvidence, "writerStopEvidence", 0); err != nil {
		return err
	}
	if err := validateLifecycleHashSet(record.ObjectHashes, "objectHashes", 0); err != nil {
		return err
	}
	if err := validateLifecycleHashSet(record.EventSegmentHashes, "eventSegmentHashes", 0); err != nil {
		return err
	}
	if err := validateLifecycleHashSet(record.ReferenceRootHashes, "referenceRootHashes", 0); err != nil {
		return err
	}
	if record.SecretDisposition != "excluded" {
		return fmt.Errorf("%w: secretDisposition must be excluded", ErrInvalidLifecycleRecord)
	}
	if err := validateLifecycleHashSet(record.SecretScanEvidence, "secretScanEvidence", 1); err != nil {
		return err
	}
	switch record.ClosureStatus {
	case "complete", "incomplete", "unknown":
	default:
		return fmt.Errorf("%w: invalid closureStatus %q", ErrInvalidLifecycleRecord, record.ClosureStatus)
	}
	if err := validateLifecycleHashSet(record.RestoreEvidence, "restoreEvidence", 0); err != nil {
		return err
	}
	if err := validateLifecycleHashSet(record.ErrorEvidence, "errorEvidence", 0); err != nil {
		return err
	}
	switch record.Outcome {
	case "completed":
		if record.ClosureStatus != "complete" {
			return fmt.Errorf("%w: completed backup requires complete closureStatus", ErrInvalidLifecycleRecord)
		}
		if len(record.ObjectHashes) == 0 || len(record.ReferenceRootHashes) == 0 || len(record.RestoreEvidence) == 0 {
			return fmt.Errorf("%w: completed backup requires object/reference/restore evidence", ErrInvalidLifecycleRecord)
		}
		if record.BackupManifestHash == nil || !validHash(*record.BackupManifestHash) {
			return fmt.Errorf("%w: completed backup requires backupManifestHash", ErrInvalidLifecycleRecord)
		}
		if len(record.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: completed backup cannot include errorEvidence", ErrInvalidLifecycleRecord)
		}
	case "failed", "uncertain":
		if record.BackupManifestHash != nil {
			return fmt.Errorf("%w: failed/uncertain backup must not include backupManifestHash", ErrInvalidLifecycleRecord)
		}
		if len(record.ErrorEvidence) == 0 {
			return fmt.Errorf("%w: failed/uncertain backup requires errorEvidence", ErrInvalidLifecycleRecord)
		}
	default:
		return fmt.Errorf("%w: invalid backup outcome %q", ErrInvalidLifecycleRecord, record.Outcome)
	}
	return nil
}

func validateLifecycleRetirement(record LifecycleRetirement) error {
	if !validID(record.RecordID) || record.Kind != "retirement" {
		return fmt.Errorf("%w: invalid retirement identity", ErrInvalidLifecycleRecord)
	}
	if err := validateLifecycleTimestamp(record.CreatedAt, "createdAt"); err != nil {
		return err
	}
	if !validLifecycleActor(record.RecordedBy) {
		return fmt.Errorf("%w: invalid recordedBy actor", ErrInvalidLifecycleRecord)
	}
	switch record.Action {
	case "deactivate", "uninstall", "delete-data", "quarantine-secret":
	default:
		return fmt.Errorf("%w: invalid retirement action %q", ErrInvalidLifecycleRecord, record.Action)
	}
	if !validHash(record.InstallationHash) || !validHash(record.InventoryHash) {
		return fmt.Errorf("%w: invalid retirement hash field", ErrInvalidLifecycleRecord)
	}
	if err := validateLifecycleHashSet(record.ProcessStopEvidence, "processStopEvidence", 1); err != nil {
		return err
	}
	if err := validateLifecycleHashSet(record.AdapterRevocationEvidence, "adapterRevocationEvidence", 1); err != nil {
		return err
	}
	if err := validateLifecycleHashSet(record.SecretIncidentEvidence, "secretIncidentEvidence", 0); err != nil {
		return err
	}
	switch record.RetentionDisposition {
	case "retain", "delete-authorized", "quarantine":
	default:
		return fmt.Errorf("%w: invalid retentionDisposition %q", ErrInvalidLifecycleRecord, record.RetentionDisposition)
	}
	if record.RetentionDisposition == "delete-authorized" {
		if record.DeletionAuthorizationHash == nil || !validHash(*record.DeletionAuthorizationHash) {
			return fmt.Errorf("%w: delete-authorized requires deletionAuthorizationHash", ErrInvalidLifecycleRecord)
		}
	} else if record.DeletionAuthorizationHash != nil {
		return fmt.Errorf("%w: deletionAuthorizationHash must be null unless delete-authorized", ErrInvalidLifecycleRecord)
	}
	if record.AuditConflictDecisionHash != nil && !validHash(*record.AuditConflictDecisionHash) {
		return fmt.Errorf("%w: invalid auditConflictDecisionHash", ErrInvalidLifecycleRecord)
	}
	if record.SharedToolsDisposition != "preserved" {
		return fmt.Errorf("%w: sharedToolsDisposition must be preserved", ErrInvalidLifecycleRecord)
	}
	if record.Action == "quarantine-secret" {
		if record.RetentionDisposition != "quarantine" {
			return fmt.Errorf("%w: quarantine-secret action requires quarantine retentionDisposition", ErrInvalidLifecycleRecord)
		}
		if len(record.SecretIncidentEvidence) == 0 {
			return fmt.Errorf("%w: quarantine-secret action requires secretIncidentEvidence", ErrInvalidLifecycleRecord)
		}
	}
	if err := validateLifecycleHashSet(record.ErrorEvidence, "errorEvidence", 0); err != nil {
		return err
	}
	switch record.Outcome {
	case "completed":
		if len(record.ErrorEvidence) != 0 {
			return fmt.Errorf("%w: completed retirement cannot include errorEvidence", ErrInvalidLifecycleRecord)
		}
	case "failed", "uncertain":
		if len(record.ErrorEvidence) == 0 {
			return fmt.Errorf("%w: failed/uncertain retirement requires errorEvidence", ErrInvalidLifecycleRecord)
		}
	default:
		return fmt.Errorf("%w: invalid retirement outcome %q", ErrInvalidLifecycleRecord, record.Outcome)
	}
	return nil
}

func validateLifecycleTimestamp(value string, field string) error {
	parsed, err := time.Parse(lifecycleTimestampLayout, value)
	if err != nil {
		return fmt.Errorf("%w: invalid %s", ErrInvalidLifecycleRecord, field)
	}
	if parsed.UTC().Format(lifecycleTimestampLayout) != value {
		return fmt.Errorf("%w: %s must be UTC millisecond timestamp", ErrInvalidLifecycleRecord, field)
	}
	return nil
}

func validateLifecycleIDSet(values []string, field string, minItems int) error {
	if len(values) < minItems {
		return fmt.Errorf("%w: %s requires at least %d item(s)", ErrInvalidLifecycleRecord, field, minItems)
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !validID(value) {
			return fmt.Errorf("%w: invalid %s value %q", ErrInvalidLifecycleRecord, field, value)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%w: duplicate %s value %q", ErrInvalidLifecycleRecord, field, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validateLifecycleHashSet(values []string, field string, minItems int) error {
	if len(values) < minItems {
		return fmt.Errorf("%w: %s requires at least %d item(s)", ErrInvalidLifecycleRecord, field, minItems)
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !validHash(value) {
			return fmt.Errorf("%w: invalid %s hash %q", ErrInvalidLifecycleRecord, field, value)
		}
		if _, exists := seen[value]; exists {
			return fmt.Errorf("%w: duplicate %s hash %q", ErrInvalidLifecycleRecord, field, value)
		}
		seen[value] = struct{}{}
	}
	return nil
}

func validLifecycleVersion(value string) bool {
	if len(value) == 0 || len(value) > 128 {
		return false
	}
	for i := 0; i < len(value); i++ {
		if value[i] < '!' || value[i] > '~' {
			return false
		}
	}
	return true
}

func validLifecycleActor(actor Actor) bool {
	if !validID(actor.ID) {
		return false
	}
	switch actor.Type {
	case "operator", "system":
		return true
	default:
		return false
	}
}

func cloneDependencies(values []LifecycleDependency) []LifecycleDependency {
	if len(values) == 0 {
		return []LifecycleDependency{}
	}
	cloned := make([]LifecycleDependency, len(values))
	copy(cloned, values)
	return cloned
}

func cloneSupports(values []LifecycleSupport) []LifecycleSupport {
	if len(values) == 0 {
		return []LifecycleSupport{}
	}
	cloned := make([]LifecycleSupport, len(values))
	for i, value := range values {
		cloned[i] = value
		cloned[i].CapabilityIDs = cloneLifecycleStrings(value.CapabilityIDs)
		if len(value.Platforms) == 0 {
			cloned[i].Platforms = []LifecyclePlatform{}
		} else {
			cloned[i].Platforms = make([]LifecyclePlatform, len(value.Platforms))
			copy(cloned[i].Platforms, value.Platforms)
		}
	}
	return cloned
}

func cloneLifecycleStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	cloned := make([]string, len(values))
	copy(cloned, values)
	return cloned
}
