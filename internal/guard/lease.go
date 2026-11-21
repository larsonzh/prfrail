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
	path := store.leasePath(resourceID)
	// P84 fix round 8 (Rev 3b, the lock substrate behind F3): the lock is an OS-level
	// exclusive lock on ONE contended name, "<record>.lock" - LockFileEx with
	// LOCKFILE_EXCLUSIVE_LOCK|LOCKFILE_FAIL_IMMEDIATELY on Windows, flock(LOCK_EX|LOCK_NB)
	// on POSIX. All writers derive the same name from the same record path, so the kernel
	// is what serialises them, across processes and across goroutines of one process
	// alike (each acquisition uses its own handle).
	//
	// Round 7 made the lock the file NAME ("<record>.lock.<pid>", CREATE_NEW / O_EXCL),
	// which tied identity to the holder's pid: two processes derived two different names
	// and never blocked each other (tester round 4 measured a parent claiming while its
	// child held the lock). An OS lock has no such failure mode, and because the kernel
	// drops it when the holding process dies - crash included - there is no stale lock to
	// heal: no pid in the name, no liveness query, no removal. A competitor only waits,
	// ctx still bounds the wait, and the retry interval is unchanged.
	lock, err := acquireLeaseLock(ctx, path)
	if err != nil {
		return Lease{}, err
	}
	defer lock.release()
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

// errLeaseLockHeld is what createLeaseLock reports when another holder owns the
// resource lock right now. It is a sentinel rather than a bare "file exists" test
// because the lock is an OS lock on a file that always exists: only the kernel's
// refusal to grant the range means "someone else holds it" (P84 fix round 8).
var errLeaseLockHeld = errors.New("lease lock is held by another holder")

// ErrLeaseLockPathIsDirectory reports that the resource's lock path exists as a
// DIRECTORY instead of the lock file this substrate takes its byte range on.
//
// P84 fix round 9 (N8): it is the upgrade boundary left by the mkdir-style lock that
// the substrate replaced in round 8 - a legacy "<record>.lock" directory, not a lock.
// It is neither acquired nor deleted here: removing it while an old run still believed
// it held the lock would hand the resource to a second writer, and it will never turn
// into a free lock file, so waiting on it is pointless too. The caller stops and
// reports; the operator migrates the directory away only after the old run is proven
// stopped.
var ErrLeaseLockPathIsDirectory = errors.New("lease lock path is a legacy lock directory")

// legacyLeaseLockDirectoryError builds the N8 refusal. It names the path so the operator
// knows what to migrate, and states the precondition (the old run must be proven stopped)
// rather than suggesting a removal that would be unsafe on its own.
func legacyLeaseLockDirectoryError(recordPath string) error {
	return fmt.Errorf("%w: %s exists as a directory; an earlier run left a mkdir-style lock directory there, and it can only be migrated after that run is proven stopped (this substrate neither acquires nor deletes it)", ErrLeaseLockPathIsDirectory, leaseLockPath(recordPath))
}

// leaseLockPathIsDirectory reports whether the lock path currently names a directory,
// i.e. whether a legacy Mkdir lock directory occupies it. It is consulted only after an
// acquisition failed, so the ordinary path pays no extra stat.
func leaseLockPathIsDirectory(recordPath string) bool {
	info, err := os.Stat(leaseLockPath(recordPath))
	return err == nil && info.IsDir()
}

// leaseLockSuffix names the single contended lock file next to a lease record. Every
// process derives the same name from the same record path, which is what makes the OS
// lock taken on it a mutual exclusion between processes (P84 Rev 3b).
const leaseLockSuffix = ".lock"

// leaseLockWaitInterval is the retry interval of a competing acquisition. It is the
// interval the previous lock shape used, so ctx semantics do not change.
const leaseLockWaitInterval = 5 * time.Millisecond

// leaseLockPath names the lock file a resource's writers contend for.
func leaseLockPath(recordPath string) string {
	return recordPath + leaseLockSuffix
}

// acquireLeaseLock takes the resource lock for this process and returns it held.
//
// The wait is the whole mechanism: an OS-held lock cannot go stale, so a competitor
// simply retries every leaseLockWaitInterval until the holder releases it or ctx
// expires. "The holder" can also be this very process through another goroutine, and
// that is still a real wait - each acquisition uses its own handle, so the kernel
// serialises them like any other pair of holders.
//
// Two outcomes are NOT contention and end the wait immediately, because neither can
// ever resolve itself: a platform failure (reported as-is) and the legacy mkdir lock
// directory (ErrLeaseLockPathIsDirectory, P84 fix round 9 / N8). The directory is
// never removed here - see the sentinel's comment for why - and withLock passes the
// error straight back to its caller.
func acquireLeaseLock(ctx context.Context, recordPath string) (*leaseLock, error) {
	for {
		lock, err := createLeaseLock(recordPath)
		if err == nil {
			return lock, nil
		}
		if !leaseLockOccupied(err) {
			return nil, fmt.Errorf("acquire lease lock: %w", err)
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(leaseLockWaitInterval):
		}
	}
}

// leaseLockOccupied reports whether a failed acquisition lost the race for the lock
// rather than failing for another reason. Only the sentinel counts: a permission or
// platform failure is reported to the caller instead of turning into an endless wait
// on a lock that will never be granted.
func leaseLockOccupied(err error) bool {
	return errors.Is(err, errLeaseLockHeld)
}

// RunLocked runs fn with the resource lock held and the resource's current record
// read under that lock, and persists nothing: fn decides what the caller then does
// with the record it was given.
//
// P84 fix round 7 (N5): the store-local proof batch has to decide on the same record
// the store's own mutations are serialised against. Reading the record and then
// writing the proofs without the lock left a window in which a renewal - which does
// take the lock - advanced the generation between the read and the writes, so the
// batch projected a superseded generation. With fn running under the lock the record
// cannot move underneath it.
//
// fn MUST NOT call back into the store's locking operations (Claim, Renew, Reclaim,
// Takeover, Release, mutate, RunLocked itself): it already holds the resource lock and
// those calls would wait for it. The generation readers it needs - Current,
// RetainedLease - take no lock and are safe. A resource without a record passes the
// zero Lease, so fn can decide whether that is an error for its purpose.
func (store *LeaseStore) RunLocked(ctx context.Context, resourceID string, fn func(current Lease) error) error {
	if fn == nil {
		return errors.New("run-locked action is required")
	}
	_, err := store.withLock(ctx, resourceID, func(path string) (Lease, bool, error) {
		lease, err := readLease(path)
		if err != nil {
			if errors.Is(err, ErrLeaseNotHeld) {
				return Lease{}, false, fn(Lease{})
			}
			return Lease{}, false, err
		}
		if lease.ResourceID != resourceID {
			return Lease{}, false, fmt.Errorf("%w: lease resource mismatch", evidence.ErrInvalidRecord)
		}
		return Lease{}, false, fn(lease)
	})
	return err
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
