package chain

// P43-T027B-STATE-SEAL: the C4 evidence seal barrier.
//
// docs/validation/evidence/P43-T027B-STATE-SEAL/01-arch-design/review.md section 2
// requires the evidence closure of a task to be durable before the accepted
// terminal is written. Before this slice the accepted event was the first and only
// place where the candidate closure, the evidence root and the two receipts became
// durable, so a reader could see a complete receipt whose evidence did not exist.
//
// The barrier makes that impossible: the accepted terminal is written only after a
// seal record that binds every hash of the closure has reached the disk with its
// own fsync. A missing seal is a gap, never a reason to reconstruct one.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	sealSchemaVersion = 1
	// sealDomain separates a seal digest from every other digest in the system.
	sealDomain  = "proofrail:task-seal:1\n"
	sealDirName = "seals"
	sealSuffix  = ".seal.json"
)

var (
	// ErrSealUnavailable reports that the seal could not be persisted.
	ErrSealUnavailable = errors.New("evidence seal unavailable")
	// ErrSealIncomplete reports that the closure carries a hash that is not a
	// valid evidence hash, so there is nothing trustworthy to seal.
	ErrSealIncomplete = errors.New("evidence closure incomplete")
	// ErrSealMismatch reports that a seal exists but does not describe the
	// closure the engine is about to accept.
	ErrSealMismatch = errors.New("evidence seal does not match the accepted closure")
	// ErrSealMissing reports that an accepted task has no seal on disk.
	ErrSealMissing = errors.New("evidence seal missing for an accepted task")
)

// SealRecord is the durable proof that one task attempt closed over a complete
// evidence set. Every field is a hash the accepted event also carries, so the
// record can be compared against the event instead of trusted on its own.
type SealRecord struct {
	SchemaVersion        int      `json:"schemaVersion"`
	RunID                string   `json:"runId"`
	TaskID               string   `json:"taskId"`
	Attempt              int      `json:"attempt"`
	ParentHash           string   `json:"parentHash"`
	CandidateHash        string   `json:"candidateHash"`
	EvidenceRootHash     string   `json:"evidenceRootHash"`
	ContextKey           string   `json:"contextKey"`
	ReviewReceiptHash    string   `json:"reviewReceiptHash"`
	PromotionReceiptHash string   `json:"promotionReceiptHash"`
	ReviewDecision       string   `json:"reviewDecision"`
	PromotionDecision    string   `json:"promotionDecision"`
	Evidence             []string `json:"evidence"`
}

// SealSink persists one seal record and reads it back. The event store owns the
// implementation because the seal lives next to the event log it protects.
type SealSink interface {
	SealTaskEvidence(ctx context.Context, record SealRecord) (string, error)
	LoadTaskSeal(ctx context.Context, runID, taskID string, attempt int) (SealRecord, string, error)
}

// ValidHash reports whether a hash is a usable evidence identifier.
func (record SealRecord) validate() error {
	if record.RunID == "" || record.TaskID == "" || record.Attempt < 1 {
		return fmt.Errorf("%w: seal identity is incomplete", ErrSealIncomplete)
	}
	if !evidence.ValidHash(record.CandidateHash) || !evidence.ValidHash(record.EvidenceRootHash) {
		return fmt.Errorf("%w: candidate or evidence root is not a hash", ErrSealIncomplete)
	}
	if record.ContextKey != "" && !evidence.ValidHash(record.ContextKey) {
		return fmt.Errorf("%w: context key is not a hash", ErrSealIncomplete)
	}
	if !evidence.ValidHash(record.ReviewReceiptHash) || !evidence.ValidHash(record.PromotionReceiptHash) {
		return fmt.Errorf("%w: a receipt hash is not a hash", ErrSealIncomplete)
	}
	for _, hash := range record.Evidence {
		if !evidence.ValidHash(hash) {
			return fmt.Errorf("%w: closure member %q is not a hash", ErrSealIncomplete, hash)
		}
	}
	return nil
}

// Hash returns the digest of the record under the seal domain. The schema version
// is part of the digest, so a field change cannot silently keep the old hash.
func (record SealRecord) Hash() (string, error) {
	record.SchemaVersion = sealSchemaVersion
	data, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("%w: encode seal record: %v", ErrSealUnavailable, err)
	}
	return evidence.Digest(sealDomain, data), nil
}

func sealFileName(taskID string, attempt int) string {
	return fmt.Sprintf("%s.%d%s", taskID, attempt, sealSuffix)
}

func (store *FileEventStore) sealDir(runID string) string {
	return filepath.Join(filepath.Dir(store.path), sealDirName, runID)
}

// SealTaskEvidence writes one seal record and returns its digest. The write is
// exclusive: an existing record is never overwritten, it is either reused byte for
// byte when it describes the same closure or reported as a mismatch.
func (store *FileEventStore) SealTaskEvidence(ctx context.Context, record SealRecord) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	record.SchemaVersion = sealSchemaVersion
	if err := record.validate(); err != nil {
		return "", err
	}
	hash, err := record.Hash()
	if err != nil {
		return "", err
	}
	dir := store.sealDir(record.RunID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("%w: create seal directory: %v", ErrSealUnavailable, err)
	}
	path := filepath.Join(dir, sealFileName(record.TaskID, record.Attempt))
	data, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("%w: encode seal record: %v", ErrSealUnavailable, err)
	}
	data = append(data, '\n')
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			existing, existingHash, loadErr := store.LoadTaskSeal(ctx, record.RunID, record.TaskID, record.Attempt)
			if loadErr != nil {
				return "", loadErr
			}
			if existingHash != hash || !sameClosure(existing, record) {
				return "", fmt.Errorf("%w: %s", ErrSealMismatch, path)
			}
			return existingHash, nil
		}
		return "", fmt.Errorf("%w: open seal record: %v", ErrSealUnavailable, err)
	}
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		// A half-written seal would read back as a mismatch forever, so the failed
		// attempt removes it instead of leaving a poison file behind.
		_ = os.Remove(path)
		return "", fmt.Errorf("%w: write seal record: %v", ErrSealUnavailable, err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return "", fmt.Errorf("%w: sync seal record: %v", ErrSealUnavailable, err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("%w: close seal record: %v", ErrSealUnavailable, err)
	}
	return hash, nil
}

// LoadTaskSeal reads one seal record back and returns it with its digest.
func (store *FileEventStore) LoadTaskSeal(ctx context.Context, runID, taskID string, attempt int) (SealRecord, string, error) {
	if err := ctx.Err(); err != nil {
		return SealRecord{}, "", err
	}
	path := filepath.Join(store.sealDir(runID), sealFileName(taskID, attempt))
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return SealRecord{}, "", fmt.Errorf("%w: %s", ErrSealMissing, path)
		}
		return SealRecord{}, "", fmt.Errorf("%w: read seal record: %v", ErrSealUnavailable, err)
	}
	var record SealRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return SealRecord{}, "", fmt.Errorf("%w: decode seal record: %v", ErrSealMismatch, err)
	}
	if err := record.validate(); err != nil {
		return SealRecord{}, "", err
	}
	hash, err := record.Hash()
	if err != nil {
		return SealRecord{}, "", err
	}
	return record, hash, nil
}

// sameClosure reports whether two records describe the same accepted closure. It is
// used to decide whether an existing seal may be reused on re-entry.
func sameClosure(left, right SealRecord) bool {
	if left.RunID != right.RunID || left.TaskID != right.TaskID || left.Attempt != right.Attempt {
		return false
	}
	if left.ParentHash != right.ParentHash || left.CandidateHash != right.CandidateHash {
		return false
	}
	if left.EvidenceRootHash != right.EvidenceRootHash {
		return false
	}
	if left.ContextKey != right.ContextKey {
		return false
	}
	if left.ReviewReceiptHash != right.ReviewReceiptHash || left.PromotionReceiptHash != right.PromotionReceiptHash {
		return false
	}
	if left.ReviewDecision != right.ReviewDecision || left.PromotionDecision != right.PromotionDecision {
		return false
	}
	if len(left.Evidence) != len(right.Evidence) {
		return false
	}
	for index := range left.Evidence {
		if left.Evidence[index] != right.Evidence[index] {
			return false
		}
	}
	return true
}
