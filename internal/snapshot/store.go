package snapshot

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type Store struct {
	root          string
	objectsDir    string
	tmpDir        string
	maxBytesQuota int64
	mu            sync.RWMutex
}

type GCPlan struct {
	CandidateHashes []string `json:"candidateHashes"`
	ReclaimedBytes  int64    `json:"reclaimedBytes"`
	RetainedHashes  []string `json:"retainedHashes"`
}

func NewStore(root string, maxBytesQuota int64) (*Store, error) {
	objectsDir := filepath.Join(root, "objects")
	tmpDir := filepath.Join(root, "tmp")
	if err := os.MkdirAll(objectsDir, 0755); err != nil {
		return nil, fmt.Errorf("create store objects dir: %w", err)
	}
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("create store tmp dir: %w", err)
	}
	return &Store{
		root:          root,
		objectsDir:    objectsDir,
		tmpDir:        tmpDir,
		maxBytesQuota: maxBytesQuota,
	}, nil
}

func (s *Store) objectPath(hash string) string {
	// hash format: sha256:<64 hex>
	// file name is the raw sha256 hex or prefixed
	safeName := hash
	if len(safeName) > 7 && safeName[:7] == "sha256:" {
		safeName = safeName[7:]
	}
	return filepath.Join(s.objectsDir, safeName)
}

func (s *Store) PutObject(ctx context.Context, data []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash := evidence.Digest("", data)
	targetPath := s.objectPath(hash)

	// Check if already exists intact
	if info, err := os.Stat(targetPath); err == nil && info.Size() == int64(len(data)) {
		// Read and verify
		existing, err := os.ReadFile(targetPath)
		if err == nil && evidence.Digest("", existing) == hash {
			return hash, nil
		}
	}

	// Check quota
	if s.maxBytesQuota > 0 {
		currentBytes, err := s.totalBytesLocked()
		if err != nil {
			return "", err
		}
		if currentBytes+int64(len(data)) > s.maxBytesQuota {
			return "", fmt.Errorf("%w: current %d + new %d > quota %d", ErrQuotaExceeded, currentBytes, len(data), s.maxBytesQuota)
		}
	}

	// Write to tmp file and sync
	randBytes := make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		return "", fmt.Errorf("generate tmp name: %w", err)
	}
	tmpFile := filepath.Join(s.tmpDir, fmt.Sprintf("put-%s-%s", hex.EncodeToString(randBytes), hash[7:15]))

	f, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return "", fmt.Errorf("create tmp object: %w", err)
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(tmpFile)
	}()

	if _, err := f.Write(data); err != nil {
		return "", fmt.Errorf("write tmp object: %w", err)
	}
	if err := f.Sync(); err != nil {
		return "", fmt.Errorf("sync tmp object: %w", err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close tmp object: %w", err)
	}

	// Atomic rename to final path
	if err := os.Rename(tmpFile, targetPath); err != nil {
		return "", fmt.Errorf("rename to target object: %w", err)
	}
	return hash, nil
}

func (s *Store) GetObject(ctx context.Context, hash string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	targetPath := s.objectPath(hash)
	data, err := os.ReadFile(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrObjectNotFound, hash)
		}
		return nil, err
	}
	if computed := evidence.Digest("", data); computed != hash {
		return nil, fmt.Errorf("%w: object %s got %s", ErrObjectCorrupt, hash, computed)
	}
	return data, nil
}

func (s *Store) HasObject(hash string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	targetPath := s.objectPath(hash)
	info, err := os.Stat(targetPath)
	if err != nil || info.IsDir() {
		return false
	}
	// Verify content integrity
	data, err := os.ReadFile(targetPath)
	if err != nil {
		return false
	}
	return evidence.Digest("", data) == hash
}

func (s *Store) TotalBytes() (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.totalBytesLocked()
}

func (s *Store) totalBytesLocked() (int64, error) {
	var total int64
	entries, err := os.ReadDir(s.objectsDir)
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		total += info.Size()
	}
	return total, nil
}

func (s *Store) GCDryRun(ctx context.Context, reachableHashes map[string]bool, lockedHashes map[string]bool, ttl time.Duration) (*GCPlan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plan := &GCPlan{
		CandidateHashes: []string{},
		RetainedHashes:  []string{},
	}
	entries, err := os.ReadDir(s.objectsDir)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cutoff := now.Add(-ttl)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		rawHex := entry.Name()
		hash := "sha256:" + rawHex
		info, err := entry.Info()
		if err != nil {
			continue
		}

		// If reachable or locked, must retain
		if reachableHashes[hash] || lockedHashes[hash] {
			plan.RetainedHashes = append(plan.RetainedHashes, hash)
			continue
		}

		// Not reachable and not locked: check age vs TTL
		if info.ModTime().Before(cutoff) {
			plan.CandidateHashes = append(plan.CandidateHashes, hash)
			plan.ReclaimedBytes += info.Size()
		} else {
			plan.RetainedHashes = append(plan.RetainedHashes, hash)
		}
	}
	return plan, nil
}
