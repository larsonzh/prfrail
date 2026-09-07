package guard

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrLeaseHeld            = errors.New("writer lease already held")
	ErrLeaseNotHeld         = errors.New("writer lease not held")
	ErrFencingTokenMismatch = errors.New("writer lease fencing token mismatch")
	ErrTakeoverEvidence     = errors.New("takeover evidence insufficient")
)

type FencingToken struct {
	ClaimID    string `json:"claimId"`
	Generation uint64 `json:"generation"`
}

type Lease struct {
	SchemaVersion  string          `json:"schemaVersion"`
	ResourceID     string          `json:"resourceId"`
	ConsumerID     string          `json:"consumerId"`
	Token          FencingToken    `json:"token"`
	Process        ProcessIdentity `json:"process"`
	ClaimedAt      string          `json:"claimedAt"`
	LeaseExpiresAt string          `json:"leaseExpiresAt"`
	PreviousHash   *string         `json:"previousHash"`
	LeaseHash      string          `json:"leaseHash"`
}

type leasePayload struct {
	SchemaVersion  string          `json:"schemaVersion"`
	ResourceID     string          `json:"resourceId"`
	ConsumerID     string          `json:"consumerId"`
	Token          FencingToken    `json:"token"`
	Process        ProcessIdentity `json:"process"`
	ClaimedAt      string          `json:"claimedAt"`
	LeaseExpiresAt string          `json:"leaseExpiresAt"`
	PreviousHash   *string         `json:"previousHash"`
}

type TakeoverProof struct {
	TerminationEvidence   TerminationEvidence
	AuthorizationEvidence []string
}

type TakeoverVerifier func(current Lease, proof TakeoverProof) error

type LeaseOption func(*LeaseStore)

func WithTakeoverVerifier(verifier TakeoverVerifier) LeaseOption {
	return func(store *LeaseStore) { store.verifyTakeover = verifier }
}

type LeaseStore struct {
	root           string
	now            func() time.Time
	verifyTakeover TakeoverVerifier
}

func NewLeaseStore(root string, options ...LeaseOption) (*LeaseStore, error) {
	if root == "" {
		return nil, errors.New("lease root is required")
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return nil, fmt.Errorf("create lease root: %w", err)
	}
	store := &LeaseStore{root: root, now: func() time.Time { return time.Now().UTC() }}
	for _, option := range options {
		option(store)
	}
	return store, nil
}

func (store *LeaseStore) Claim(ctx context.Context, resourceID, consumerID, claimID string, process ProcessIdentity, ttl time.Duration) (Lease, error) {
	return store.mutate(ctx, resourceID, func(current *Lease) (Lease, bool, error) {
		if current != nil {
			return Lease{}, false, ErrLeaseHeld
		}
		return store.newLease(resourceID, consumerID, FencingToken{ClaimID: claimID, Generation: 1}, process, ttl, nil)
	})
}

func (store *LeaseStore) Renew(ctx context.Context, resourceID string, token FencingToken, ttl time.Duration) (Lease, error) {
	return store.mutate(ctx, resourceID, func(current *Lease) (Lease, bool, error) {
		if current == nil {
			return Lease{}, false, ErrLeaseNotHeld
		}
		if current.Token != token {
			return Lease{}, false, ErrFencingTokenMismatch
		}
		previous := current.LeaseHash
		return store.newLease(resourceID, current.ConsumerID, FencingToken{ClaimID: token.ClaimID, Generation: token.Generation + 1}, current.Process, ttl, &previous)
	})
}

func (store *LeaseStore) Takeover(ctx context.Context, resourceID, consumerID, claimID string, process ProcessIdentity, ttl time.Duration, proof TakeoverProof) (Lease, error) {
	return store.mutate(ctx, resourceID, func(current *Lease) (Lease, bool, error) {
		if current == nil {
			return Lease{}, false, ErrLeaseNotHeld
		}
		if claimID == current.Token.ClaimID || !store.validTakeoverProof(*current, proof) {
			return Lease{}, false, ErrTakeoverEvidence
		}
		if store.verifyTakeover == nil || store.verifyTakeover(*current, proof) != nil {
			return Lease{}, false, ErrTakeoverEvidence
		}
		previous := current.LeaseHash
		return store.newLease(resourceID, consumerID, FencingToken{ClaimID: claimID, Generation: current.Token.Generation + 1}, process, ttl, &previous)
	})
}

func (store *LeaseStore) validTakeoverProof(current Lease, proof TakeoverProof) bool {
	termination := proof.TerminationEvidence
	if termination.Identity != current.Process || termination.Outcome != "stopped" || !evidence.ValidHash(termination.Hash) || !validEvidenceHashes(proof.AuthorizationEvidence) {
		return false
	}
	verifiedAt, err := time.Parse("2006-01-02T15:04:05.000Z", termination.Verified)
	if err != nil || verifiedAt.Before(store.now().Add(-5*time.Minute)) || verifiedAt.After(store.now().Add(time.Minute)) {
		return false
	}
	canonical, err := evidence.EncodeCanonical(struct {
		Identity  ProcessIdentity `json:"identity"`
		Requested string          `json:"requestedAt"`
		Verified  string          `json:"verifiedAt"`
		Outcome   string          `json:"outcome"`
		Actions   []string        `json:"actions"`
	}{termination.Identity, termination.Requested, termination.Verified, termination.Outcome, termination.Actions})
	return err == nil && termination.Hash == evidence.Digest("proofrail:termination-evidence:1\n", canonical)
}

func (store *LeaseStore) Validate(ctx context.Context, resourceID string, token FencingToken) error {
	lease, err := store.Current(ctx, resourceID)
	if err != nil {
		return err
	}
	if lease.Token != token {
		return ErrFencingTokenMismatch
	}
	return nil
}

func (store *LeaseStore) Release(ctx context.Context, resourceID string, token FencingToken) error {
	_, err := store.withLock(ctx, resourceID, func(path string) (Lease, bool, error) {
		current, err := readLease(path)
		if err != nil {
			return Lease{}, false, err
		}
		if current.ResourceID != resourceID {
			return Lease{}, false, fmt.Errorf("%w: lease resource mismatch", evidence.ErrInvalidRecord)
		}
		if current.Token != token {
			return Lease{}, false, ErrFencingTokenMismatch
		}
		if err := os.Remove(path); err != nil {
			return Lease{}, false, err
		}
		return Lease{}, false, nil
	})
	return err
}

func (store *LeaseStore) Current(ctx context.Context, resourceID string) (Lease, error) {
	if err := validateLeaseID(resourceID); err != nil {
		return Lease{}, err
	}
	lease, err := readLease(store.leasePath(resourceID))
	if err == nil && lease.ResourceID != resourceID {
		return Lease{}, fmt.Errorf("%w: lease resource mismatch", evidence.ErrInvalidRecord)
	}
	return lease, err
}

func (store *LeaseStore) newLease(resourceID, consumerID string, token FencingToken, process ProcessIdentity, ttl time.Duration, previous *string) (Lease, bool, error) {
	if err := validateLeaseInputs(resourceID, consumerID, token.ClaimID, process, ttl); err != nil {
		return Lease{}, false, err
	}
	now := store.now()
	lease := Lease{SchemaVersion: "1.0.0", ResourceID: resourceID, ConsumerID: consumerID, Token: token, Process: process, ClaimedAt: formatTime(now), LeaseExpiresAt: formatTime(now.Add(ttl)), PreviousHash: previous}
	canonical, err := evidence.EncodeCanonical(payloadOf(lease))
	if err != nil {
		return Lease{}, false, err
	}
	lease.LeaseHash = evidence.Digest("proofrail:writer-lease:1\n", canonical)
	return lease, true, nil
}

func (store *LeaseStore) mutate(ctx context.Context, resourceID string, change func(*Lease) (Lease, bool, error)) (Lease, error) {
	return store.withLock(ctx, resourceID, func(path string) (Lease, bool, error) {
		var current *Lease
		lease, err := readLease(path)
		if err == nil {
			if lease.ResourceID != resourceID {
				return Lease{}, false, fmt.Errorf("%w: lease resource mismatch", evidence.ErrInvalidRecord)
			}
			current = &lease
		} else if !errors.Is(err, ErrLeaseNotHeld) {
			return Lease{}, false, err
		}
		return change(current)
	})
}

func (store *LeaseStore) withLock(ctx context.Context, resourceID string, action func(string) (Lease, bool, error)) (Lease, error) {
	if err := validateLeaseID(resourceID); err != nil {
		return Lease{}, err
	}
	lockPath := store.leasePath(resourceID) + ".lock"
	for {
		if err := os.Mkdir(lockPath, 0700); err == nil {
			break
		} else if !os.IsExist(err) && !os.IsPermission(err) {
			return Lease{}, fmt.Errorf("acquire lease lock: %w", err)
		}
		select {
		case <-ctx.Done():
			return Lease{}, ctx.Err()
		case <-time.After(5 * time.Millisecond):
		}
	}
	defer os.Remove(lockPath)
	path := store.leasePath(resourceID)
	lease, write, err := action(path)
	if err != nil || !write {
		return lease, err
	}
	if err := writeLease(path, lease); err != nil {
		return Lease{}, err
	}
	return lease, nil
}

func writeLease(path string, lease Lease) error {
	data, err := evidence.EncodeCanonical(lease)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".lease-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err = tmp.Write(append(data, '\n')); err == nil {
		err = tmp.Sync()
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func readLease(path string) (Lease, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Lease{}, ErrLeaseNotHeld
	}
	if err != nil {
		return Lease{}, err
	}
	var lease Lease
	if err := evidence.DecodeStrictJSON(data, &lease); err != nil {
		return Lease{}, fmt.Errorf("decode lease: %w", err)
	}
	if lease.SchemaVersion != "1.0.0" || validateLeaseInputs(lease.ResourceID, lease.ConsumerID, lease.Token.ClaimID, lease.Process, time.Second) != nil || lease.Token.Generation == 0 || !evidence.ValidHash(lease.LeaseHash) {
		return Lease{}, fmt.Errorf("%w: invalid lease fields", evidence.ErrInvalidRecord)
	}
	claimedAt, claimErr := time.Parse("2006-01-02T15:04:05.000Z", lease.ClaimedAt)
	expiresAt, expiryErr := time.Parse("2006-01-02T15:04:05.000Z", lease.LeaseExpiresAt)
	if claimErr != nil || expiryErr != nil || !expiresAt.After(claimedAt) || (lease.PreviousHash != nil && !evidence.ValidHash(*lease.PreviousHash)) {
		return Lease{}, fmt.Errorf("%w: invalid lease chronology", evidence.ErrInvalidRecord)
	}
	canonical, err := evidence.EncodeCanonical(payloadOf(lease))
	if err != nil || lease.LeaseHash != evidence.Digest("proofrail:writer-lease:1\n", canonical) {
		return Lease{}, fmt.Errorf("%w: lease hash mismatch", evidence.ErrInvalidRecord)
	}
	return lease, nil
}

func payloadOf(lease Lease) leasePayload {
	return leasePayload{lease.SchemaVersion, lease.ResourceID, lease.ConsumerID, lease.Token, lease.Process, lease.ClaimedAt, lease.LeaseExpiresAt, lease.PreviousHash}
}

func (store *LeaseStore) leasePath(resourceID string) string {
	return filepath.Join(store.root, resourceID+".lease.json")
}

func validateLeaseID(value string) error {
	if !evidence.ValidID(value) {
		return fmt.Errorf("%w: invalid lease identifier %q", evidence.ErrInvalidRecord, value)
	}
	return nil
}

func validateLeaseInputs(resourceID, consumerID, claimID string, process ProcessIdentity, ttl time.Duration) error {
	if err := validateLeaseID(resourceID); err != nil {
		return err
	}
	if !evidence.ValidID(consumerID) || !evidence.ValidID(claimID) || process.PID <= 0 || process.StartToken == "" || ttl <= 0 {
		return fmt.Errorf("%w: invalid lease fields", evidence.ErrInvalidRecord)
	}
	return nil
}

func validEvidenceHashes(values []string) bool {
	if len(values) == 0 {
		return false
	}
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if !evidence.ValidHash(value) {
			return false
		}
		if _, ok := seen[value]; ok {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}
