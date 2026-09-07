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

func (s *Store) objectPath(hash string) (string, error) {
	if !evidence.ValidHash(hash) {
		return "", fmt.Errorf("%w: invalid object hash %q", ErrInvalidManifest, hash)
	}
	return filepath.Join(s.objectsDir, hash[7:]), nil
}

func (s *Store) PutObject(ctx context.Context, data []byte) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash := evidence.Digest("", data)
	targetPath, _ := s.objectPath(hash)

	if _, err := os.Stat(targetPath); err == nil {
		existing, readErr := os.ReadFile(targetPath)
		if readErr == nil && evidence.Digest("", existing) == hash {
			return hash, nil
		}
		return "", fmt.Errorf("%w: immutable object %s already exists with different bytes", ErrObjectCorrupt, hash)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat target object: %w", err)
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

	// A hard link publishes the fully synced temporary inode without replacing an
	// existing immutable object. Both directories are created on the same volume.
	if err := os.Link(tmpFile, targetPath); err != nil {
		if existing, readErr := os.ReadFile(targetPath); readErr == nil && evidence.Digest("", existing) == hash {
			return hash, nil
		}
		return "", fmt.Errorf("publish target object without replacement: %w", err)
	}
	return hash, nil
}

func (s *Store) GetObject(ctx context.Context, hash string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	targetPath, err := s.objectPath(hash)
	if err != nil {
		return nil, err
	}
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

	targetPath, err := s.objectPath(hash)
	if err != nil {
		return false
	}
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
