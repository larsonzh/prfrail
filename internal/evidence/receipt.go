package evidence

import (
	"fmt"
	"regexp"
	"time"
)

const evidenceManifestDomain = "proofrail:evidence-manifest:1\n"

var mediaTypePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9!#$&^_.+\-]*/[a-z0-9][a-z0-9!#$&^_.+\-]*$`)

type EvidenceItem struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	ObjectHash string `json:"objectHash"`
	MediaType  string `json:"mediaType"`
	Redaction  string `json:"redaction"`
}

type EvidenceManifestBody struct {
	EvidenceID            string         `json:"evidenceId"`
	CreatedAt             string         `json:"createdAt"`
	RunID                 string         `json:"runId"`
	TaskID                string         `json:"taskId"`
	Attempt               int            `json:"attempt"`
	TaskDefinitionHash    string         `json:"taskDefinitionHash"`
	PolicyHash            string         `json:"policyHash"`
	ParentSnapshotHash    string         `json:"parentSnapshotHash"`
	CandidateSnapshotHash string         `json:"candidateSnapshotHash"`
	Items                 []EvidenceItem `json:"items"`
}

type EvidenceManifest struct {
	SchemaVersion    string               `json:"schemaVersion"`
	Manifest         EvidenceManifestBody `json:"manifest"`
	EvidenceRootHash string               `json:"evidenceRootHash"`
}

func NewEvidenceManifest(manifest EvidenceManifestBody) (EvidenceManifest, error) {
	if err := validateEvidenceManifestBody(manifest); err != nil {
		return EvidenceManifest{}, err
	}
	canonical, err := canonicalValue(manifest)
	if err != nil {
		return EvidenceManifest{}, err
	}
	return EvidenceManifest{
		SchemaVersion: SchemaVersion, Manifest: manifest,
		EvidenceRootHash: Digest(evidenceManifestDomain, canonical),
	}, nil
}

func DecodeEvidenceManifest(input []byte) (EvidenceManifest, error) {
	var manifest EvidenceManifest
	if err := decodeStrictJSON(input, &manifest); err != nil {
		return EvidenceManifest{}, err
	}
	if err := VerifyEvidenceManifest(manifest); err != nil {
		return EvidenceManifest{}, err
	}
	return manifest, nil
}

func VerifyEvidenceManifest(manifest EvidenceManifest) error {
	if manifest.SchemaVersion != SchemaVersion {
		return fmt.Errorf("%w: %q", ErrUnsupportedVersion, manifest.SchemaVersion)
	}
	if err := validateEvidenceManifestBody(manifest.Manifest); err != nil {
		return err
	}
	canonical, err := canonicalValue(manifest.Manifest)
	if err != nil {
		return err
	}
	if expected := Digest(evidenceManifestDomain, canonical); manifest.EvidenceRootHash != expected {
		return fmt.Errorf("%w: evidence root hash mismatch", ErrInvalidRecord)
	}
	return nil
}

func validateEvidenceManifestBody(manifest EvidenceManifestBody) error {
	if !validID(manifest.EvidenceID) || !validID(manifest.RunID) || !validID(manifest.TaskID) || manifest.Attempt < 1 {
		return fmt.Errorf("%w: invalid evidence identity", ErrInvalidRecord)
	}
	if _, err := time.Parse("2006-01-02T15:04:05.000Z", manifest.CreatedAt); err != nil {
		return fmt.Errorf("%w: invalid evidence timestamp", ErrInvalidRecord)
	}
	hashes := []string{manifest.TaskDefinitionHash, manifest.PolicyHash, manifest.ParentSnapshotHash, manifest.CandidateSnapshotHash}
	for _, hash := range hashes {
		if !validHash(hash) {
			return fmt.Errorf("%w: invalid manifest hash", ErrInvalidRecord)
		}
	}
	if manifest.ParentSnapshotHash == manifest.CandidateSnapshotHash {
		return fmt.Errorf("%w: parent and candidate snapshots are equal", ErrInvalidRecord)
	}
	if len(manifest.Items) == 0 {
		return fmt.Errorf("%w: evidence items are empty", ErrInvalidRecord)
	}
	seenIDs := make(map[string]struct{}, len(manifest.Items))
	for index, item := range manifest.Items {
		if !validID(item.ID) || !validEvidenceKind(item.Kind) || !validHash(item.ObjectHash) || !mediaTypePattern.MatchString(item.MediaType) {
			return fmt.Errorf("%w: invalid evidence item %d", ErrInvalidRecord, index)
		}
		if item.Redaction != "not-required" && item.Redaction != "passed" {
			return fmt.Errorf("%w: invalid redaction status", ErrInvalidRecord)
		}
		if _, exists := seenIDs[item.ID]; exists {
			return fmt.Errorf("%w: duplicate evidence item ID %q", ErrInvalidRecord, item.ID)
		}
		seenIDs[item.ID] = struct{}{}
		if index > 0 && compareEvidenceItems(manifest.Items[index-1], item) >= 0 {
			return fmt.Errorf("%w: evidence items are not strictly sorted", ErrInvalidRecord)
		}
	}
	return nil
}

func compareEvidenceItems(left, right EvidenceItem) int {
	if comparison := compareUTF16(left.Kind, right.Kind); comparison != 0 {
		return comparison
	}
	return compareUTF16(left.ID, right.ID)
}

func validEvidenceKind(kind string) bool {
	switch kind {
	case "state-event", "error", "adapter-request", "adapter-result", "adapter-receipt", "hook-result", "artifact", "change-set", "diff", "ticket", "repair-transaction", "authorization-record", "effect-record", "cost-ledger", "handoff-receipt":
		return true
	default:
		return false
	}
}
