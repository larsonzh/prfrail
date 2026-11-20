package guard

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

var (
	ErrLeaseHeld            = errors.New("writer lease already held")
	ErrLeaseNotHeld         = errors.New("writer lease not held")
	ErrFencingTokenMismatch = errors.New("writer lease fencing token mismatch")
	ErrTakeoverEvidence     = errors.New("takeover evidence insufficient")
	// ErrLeaseHolderMismatch reports a renewal presented by a process that is not
	// the holder the record names (T5).
	ErrLeaseHolderMismatch = errors.New("writer lease holder mismatch")
	// ErrLeaseReclaimUnproven reports an adoption attempted without the ownership
	// proof the transfer requires (T5).
	ErrLeaseReclaimUnproven = errors.New("writer lease reclaim unproven")
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

// Renew extends the lease the caller currently holds and advances its
// generation.
//
// T5 (P83 fix round 6): a token alone does not prove ownership. The old form
// renewed on a matching token and kept the recorded process identity, so any
// process that could read the token out of the run root - including a second
// successor of a stopped holder - could keep extending a lease it never held.
// The caller therefore has to present its own process identity, and renewal
// refuses unless that identity is the holder the record already names.
func (store *LeaseStore) Renew(ctx context.Context, resourceID string, token FencingToken, process ProcessIdentity, ttl time.Duration) (Lease, error) {
	return store.mutate(ctx, resourceID, func(current *Lease) (Lease, bool, error) {
		if current == nil {
			return Lease{}, false, ErrLeaseNotHeld
		}
		if current.Token != token {
			return Lease{}, false, ErrFencingTokenMismatch
		}
		if current.Process != process {
			return Lease{}, false, ErrLeaseHolderMismatch
		}
		previous := current.LeaseHash
		return store.newLease(resourceID, current.ConsumerID, FencingToken{ClaimID: token.ClaimID, Generation: token.Generation + 1}, process, ttl, &previous)
	})
}

// Reclaim transfers a held lease to a new holder and advances the generation.
//
// T5 (P83 fix round 6): adoption used to hand the successor the OLD record, so
// two successors of a stopped holder could both "adopt" it, both verify the same
// token, and both believe they hold the writer. Reclaiming makes adoption a
// transfer: the record names the new holder and a new generation, and the second
// successor - which still presents the previous token - is refused with
// ErrFencingTokenMismatch.
//
// proof runs under the store lock, so the "the previous holder is stopped"
// decision and the transfer cannot be separated by a competing successor. A
// missing proof is a refusal: the store never transfers on trust.
func (store *LeaseStore) Reclaim(ctx context.Context, resourceID string, token FencingToken, process ProcessIdentity, ttl time.Duration, proof func(Lease) error) (Lease, error) {
	return store.mutate(ctx, resourceID, func(current *Lease) (Lease, bool, error) {
		if current == nil {
			return Lease{}, false, ErrLeaseNotHeld
		}
		if current.Token != token {
			return Lease{}, false, ErrFencingTokenMismatch
		}
		if proof == nil {
			return Lease{}, false, ErrLeaseReclaimUnproven
		}
		if err := proof(*current); err != nil {
			return Lease{}, false, err
		}
		previous := current.LeaseHash
		return store.newLease(resourceID, current.ConsumerID, FencingToken{ClaimID: token.ClaimID, Generation: token.Generation + 1}, process, ttl, &previous)
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
	// The record this call is about to replace is read before the action runs, so
	// T5 retention below can keep the outgoing generation addressable by its hash.
	replaced, replacedErr := readLease(path)
	lease, write, err := action(path)
	if err != nil || !write {
		return lease, err
	}
	// T5 (P83 fix round): the record a renewal replaces stays addressable by the
	// hash it carries, so a receipt bound to the previous generation remains
	// independently re-verifiable after the renewal moved the current record on.
	if replacedErr == nil && replaced.LeaseHash != lease.LeaseHash {
		if err := store.retainGeneration(replaced); err != nil {
			return Lease{}, err
		}
	}
	if err := writeLease(path, lease); err != nil {
		return Lease{}, err
	}
	return lease, nil
}

// generationPath names the retained body of one lease generation by its hash. It
// sits next to the current record and never replaces it.
//
// NOTE (P83 fix round 6, registered as a non-blocking finding): retention is per
// generation and never pruned, so a long-lived run grows this directory by one
// file per renewal. A retention policy belongs HERE - generationPath is the only
// place a retained name is derived, so a policy would delete by age or sequence
// next to it and would never rewrite a name that is itself a hash.
func (store *LeaseStore) generationPath(resourceID, leaseHash string) string {
	return filepath.Join(store.root, resourceID+"."+strings.TrimPrefix(leaseHash, "sha256:")+".lease.json")
}

// retainGeneration writes the record a renewal is about to replace under its own
// hash. T5: renewal used to overwrite the only copy, so an already completed
// receipt that bound generation N-1 could no longer be re-verified against the
// body it named.
//
// P83 FIX ROUND 6 (T5 "retention treats path-exists as success"): an existing
// retained body is not accepted because the name exists. It is read back and has
// to be the record whose hash names it - fields, chronology, hash - so a
// corrupted or truncated retained file refuses the renewal instead of silently
// standing in for an earlier generation that can no longer be verified. Writing a
// different body under an occupied name is never an option: the name IS the hash.
//
// Retention is idempotent for a valid body: re-reading one that already verifies
// keeps it exactly as it is.
func (store *LeaseStore) retainGeneration(lease Lease) error {
	if !evidence.ValidHash(lease.LeaseHash) {
		return fmt.Errorf("%w: invalid retained lease hash", evidence.ErrInvalidRecord)
	}
	path := store.generationPath(lease.ResourceID, lease.LeaseHash)
	// T5#3 (P83 fix round 7): the existence probe must not follow a link either, or
	// a symlinked name would be judged "already retained" and the verification below
	// would be the first thing to notice. Lstat keeps the decision on the name the
	// store itself derives, and the read that follows refuses a link outright.
	if _, err := os.Lstat(path); err == nil {
		if _, verifyErr := store.retainedLease(lease.ResourceID, lease.LeaseHash); verifyErr != nil {
			return fmt.Errorf("%w: the retained generation %s cannot be verified: %v", evidence.ErrInvalidRecord, lease.LeaseHash, verifyErr)
		}
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return writeLease(path, lease)
}

// RetainedLease reads one earlier generation of a resource back by the hash it
// carries, verifying it exactly like the current record. It is the exported form of
// the re-verification seam: T5#4 (P83 fix round 8, item C) reconciles a proof that
// was left on a superseded token only when the predecessor it names is a record the
// store itself retained, so the reconciliation can never adopt a token that was
// never issued for this resource.
func (store *LeaseStore) RetainedLease(resourceID, leaseHash string) (Lease, error) {
	return store.retainedLease(resourceID, leaseHash)
}

// retainedLease reads one lease generation back by the hash it carries. It is the
// re-verification seam for earlier generations: the body is checked exactly like
// the current record (fields, chronology and hash), and a body that does not hash
// to the requested identity is refused instead of standing in for it.
func (store *LeaseStore) retainedLease(resourceID, leaseHash string) (Lease, error) {
	if err := validateLeaseID(resourceID); err != nil {
		return Lease{}, err
	}
	if !evidence.ValidHash(leaseHash) {
		return Lease{}, fmt.Errorf("%w: invalid retained lease hash", evidence.ErrInvalidRecord)
	}
	lease, err := readLease(store.generationPath(resourceID, leaseHash))
	if err != nil {
		return Lease{}, err
	}
	if lease.ResourceID != resourceID || lease.LeaseHash != leaseHash {
		return Lease{}, fmt.Errorf("%w: retained lease hash mismatch", evidence.ErrInvalidRecord)
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

// leaseRecordMaxBytes bounds one lease record read. A record is a few hundred
// bytes, so the bound exists only to stop a redirected or hostile file from making
// the store read without end while it holds the resource lock (T5#3, P83 fix round
// 7).
const leaseRecordMaxBytes int64 = 64 << 10

// readLease reads one lease record through a bounded, link-refusing handle.
//
// T5#3 (P83 fix round 7): the record used to be read with os.ReadFile, which
// follows a link and reads until end of file. Retention and re-verification both run
// while the store lock is held, so a record redirected to another file - or to a
// stream that never ends - could redirect the verification or block every writer of
// the resource. The read now opens without following a link, requires a regular file
// on the returned descriptor and stops at leaseRecordMaxBytes; every one of those
// violations is ErrInvalidRecord, never a shorter or a redirected read.
func readLease(path string) (Lease, error) {
	data, err := readLeaseRecord(path)
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

// readLeaseRecord returns the bounded bytes of one record. A missing record keeps
// its "not held" meaning so callers can still tell absence from corruption.
func readLeaseRecord(path string) ([]byte, error) {
	file, err := openLeaseRecord(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrLeaseNotHeld
		}
		return nil, fmt.Errorf("%w: open lease record: %v", evidence.ErrInvalidRecord, err)
	}
	defer func() { _ = file.Close() }()
	info, statErr := file.Stat()
	if statErr != nil {
		return nil, fmt.Errorf("%w: stat lease record: %v", evidence.ErrInvalidRecord, statErr)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("%w: lease record is not a regular file", evidence.ErrInvalidRecord)
	}
	data, readErr := io.ReadAll(io.LimitReader(file, leaseRecordMaxBytes+1))
	if readErr != nil {
		return nil, fmt.Errorf("%w: read lease record: %v", evidence.ErrInvalidRecord, readErr)
	}
	if int64(len(data)) > leaseRecordMaxBytes {
		return nil, fmt.Errorf("%w: lease record exceeds its %d byte bound", evidence.ErrInvalidRecord, leaseRecordMaxBytes)
	}
	return data, nil
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
