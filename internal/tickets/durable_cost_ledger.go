package tickets

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// ErrDurableNotCommitted reports that the canonical ledger snapshot for this
// operation was not durably committed. It is the umbrella classification, so
// errors.Is(err, ErrDurableNotCommitted) holds for all three outcomes below. A
// snapshot that is already visible is never removed or rolled back.
var ErrDurableNotCommitted = errors.New("durable cost ledger snapshot not committed")

// ErrDurableVisibleUnconfirmed reports that the new snapshot is (or may be)
// visible on disk while the durability sequence did not finish, so the caller
// holds no confirmation. See docs/CONTRACTS.md and the design §3 failure matrix.
var ErrDurableVisibleUnconfirmed = fmt.Errorf("%w: snapshot visible but unconfirmed", ErrDurableNotCommitted)

// ErrDurableOutcomeUnknown reports that the in-lock re-read could prove neither
// the target snapshot nor the pre-commit snapshot was observed.
var ErrDurableOutcomeUnknown = fmt.Errorf("%w: commit outcome unknown", ErrDurableNotCommitted)

// DurableCostLedgerConfig binds one durable ledger to exactly one shared scope.
// Root is the fixed budget root; Ledger is the immutable ledger identity and
// Initial is the allocation appended when the canonical ledger is created for
// the first time. Re-opening an initialized scope never re-issues Initial.
type DurableCostLedgerConfig struct {
	Root    string
	Ledger  CostLedgerConfig
	Initial AllocationInput

	// CommitFault is an optional test-only fault injection surface. It is
	// called with one of the Stage* constants at the corresponding commit
	// stage; a non-nil return value is treated as that stage failing. The
	// production default is nil, which means "never inject".
	CommitFault func(stage string) error
}

// DurableAck is the opaque durability confirmation for a single ledger entry.
// Its fields are private on purpose: a zero value is invalid and a caller can
// only check it, never construct it.
type DurableAck struct {
	valid        bool
	scopeKind    string
	scopeHash    string
	snapshotHash string
	entryHash    string
	confirmed    bool
}

// Valid reports whether this ack is a real durability confirmation.
func (ack DurableAck) Valid() bool {
	return ack.valid && ack.confirmed && ack.snapshotHash != "" && evidence.ValidHash(ack.entryHash)
}

// RequireEntry verifies that this confirmation covers entryHash.
func (ack DurableAck) RequireEntry(entryHash string) error {
	if !ack.Valid() {
		return fmt.Errorf("%w: durable ack is not confirmed", ErrInvalidCostLedger)
	}
	if !evidence.ValidHash(entryHash) {
		return fmt.Errorf("%w: invalid entry hash %q", ErrInvalidCostLedger, entryHash)
	}
	if entryHash != ack.entryHash {
		return fmt.Errorf("%w: durable ack covers entry %q, not %q", ErrInvalidCostLedger, ack.entryHash, entryHash)
	}
	return nil
}

// DurableCostLedger wraps the pure CostLedger rule model with a scope-wide
// cross-process lock, reload-before-decide and a durable commit sequence.
type DurableCostLedger struct {
	mu          sync.Mutex
	root        string
	path        string
	lockPath    string
	identity    CostLedgerConfig
	initial     AllocationInput
	commitFault func(stage string) error
	// confirmed holds the entries this instance has durably confirmed and
	// confirmedModel is the rule model over exactly that prefix. Read-only
	// queries answer from confirmedModel only, so a reloaded tail can never
	// manufacture a confirmation.
	confirmed      []any
	confirmedModel *CostLedger
	closed         bool
}

// OpenDurableCostLedger resolves the canonical ledger path for the configured
// scope, creates the initial allocation when the ledger does not exist yet and
// otherwise reloads and verifies the existing snapshot.
func OpenDurableCostLedger(config DurableCostLedgerConfig) (*DurableCostLedger, error) {
	root := strings.TrimSpace(config.Root)
	if root == "" {
		return nil, fmt.Errorf("%w: durable root is required", ErrInvalidCostLedger)
	}
	if err := validateCostLedgerConfig(config.Ledger); err != nil {
		return nil, err
	}
	root = filepath.Clean(root)
	if err := os.MkdirAll(root, 0o700); err != nil {
		return nil, err
	}
	name := durableLedgerFileName(config.Ledger)
	ledger := &DurableCostLedger{
		root:        root,
		path:        filepath.Join(root, name),
		lockPath:    filepath.Join(root, name+".lock"),
		identity:    config.Ledger,
		initial:     config.Initial,
		commitFault: config.CommitFault,
	}
	if err := ledger.run(context.Background(), ledger.bootstrap); err != nil {
		return nil, err
	}
	return ledger, nil
}

// durableLedgerFileName derives the single canonical file name for a scope.
// Only the hexadecimal part of the scope hash may become a path component.
func durableLedgerFileName(config CostLedgerConfig) string {
	hex := strings.TrimPrefix(config.ScopeHash, "sha256:")
	var b strings.Builder
	for _, r := range hex {
		switch {
		case r >= '0' && r <= '9', r >= 'a' && r <= 'f':
			b.WriteRune(r)
		case r >= 'A' && r <= 'F':
			b.WriteRune(r + ('a' - 'A'))
		}
	}
	component := b.String()
	if component == "" {
		component = strings.TrimPrefix(
			evidence.Digest("proofrail:durable-cost-ledger-path:1\n", []byte(config.ScopeHash)),
			"sha256:",
		)
	}
	scope := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-':
			return r
		}
		return -1
	}, config.ScopeKind)
	if scope == "" {
		scope = "scope"
	}
	return "cost-" + scope + "-" + component + ".json"
}

// Path returns the canonical ledger path for this scope.
func (ledger *DurableCostLedger) Path() string {
	if ledger == nil {
		return ""
	}
	return ledger.path
}

// Close releases the durable ledger handle. The canonical snapshot and the
// side-car lock file are intentionally left on disk.
func (ledger *DurableCostLedger) Close() error {
	if ledger == nil {
		return nil
	}
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	ledger.closed = true
	return nil
}

// Reserve appends one reservation through the full durable protocol.
func (ledger *DurableCostLedger) Reserve(ctx context.Context, input ReservationInput) (ReservationOutcome, DurableAck, error) {
	var outcome ReservationOutcome
	var ack DurableAck
	err := ledger.run(ctx, func() error {
		current, record, err := ledger.loadLocked()
		if err != nil {
			return err
		}
		local, err := current.Reserve(input)
		if err != nil {
			return err
		}
		outcome = local
		snapshot, err := current.Snapshot()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrDurableNotCommitted, err)
		}
		if err := ledger.commitLocked(snapshot); err != nil {
			return err
		}
		ledger.confirmed = copyCostEntries(snapshot.Ledger.Entries)
		ledger.confirmedModel = current
		ack = newDurableAck(ledger.identity, snapshot, record, local.ReservationHash)
		return nil
	})
	if err != nil {
		if durableOutcomeVisible(err) {
			return outcome, DurableAck{}, err
		}
		return ReservationOutcome{}, DurableAck{}, err
	}
	return outcome, ack, nil
}

// Settle appends one settlement through the full durable protocol. An
// idempotent hit is re-confirmed against the latest locked snapshot.
func (ledger *DurableCostLedger) Settle(ctx context.Context, input SettlementInput) (SettlementOutcome, DurableAck, error) {
	var outcome SettlementOutcome
	var ack DurableAck
	err := ledger.run(ctx, func() error {
		current, record, err := ledger.loadLocked()
		if err != nil {
			return err
		}
		local, err := current.Settle(input)
		if err != nil {
			return err
		}
		outcome = local
		snapshot, err := current.Snapshot()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrDurableNotCommitted, err)
		}
		if err := ledger.commitLocked(snapshot); err != nil {
			return err
		}
		ledger.confirmed = copyCostEntries(snapshot.Ledger.Entries)
		ledger.confirmedModel = current
		ack = newDurableAck(ledger.identity, snapshot, record, local.EntryHash)
		return nil
	})
	if err != nil {
		if durableOutcomeVisible(err) {
			return outcome, DurableAck{}, err
		}
		return SettlementOutcome{}, DurableAck{}, err
	}
	return outcome, ack, nil
}

// RequireOutstandingReservation re-reads the canonical ledger under the scope
// lock and only answers for reservations that belong to the confirmed prefix of
// this instance. An unconfirmed tail is reported as ErrCostReservationNotFound:
// reading it must never promote it into a confirmation. The read-only path
// never creates the ledger or the side-car lock file.
func (ledger *DurableCostLedger) RequireOutstandingReservation(ctx context.Context, reservationHash string) (CostReservationEntry, error) {
	var entry CostReservationEntry
	err := ledger.runRead(ctx, func() error {
		if _, _, err := ledger.loadLocked(); err != nil {
			return err
		}
		if !evidence.ValidHash(reservationHash) {
			return fmt.Errorf("%w: invalid reservationHash", ErrInvalidCostLedger)
		}
		if ledger.confirmedModel == nil {
			return fmt.Errorf("%w: durable ledger has no confirmed prefix", ErrInvalidCostLedger)
		}
		found, err := ledger.confirmedModel.RequireOutstandingReservation(reservationHash)
		if err != nil {
			return err
		}
		entry = found
		return nil
	})
	if err != nil {
		return CostReservationEntry{}, err
	}
	return entry, nil
}

// durableOutcomeVisible reports whether a commit failure happened after the new
// snapshot may already be visible. Only such a failure returns the locally
// computed outcome; a caller may act on it only when the error is nil and the
// ack is valid.
func durableOutcomeVisible(err error) bool {
	return errors.Is(err, ErrDurableVisibleUnconfirmed) || errors.Is(err, ErrDurableOutcomeUnknown)
}

func newDurableAck(identity CostLedgerConfig, snapshot CostLedgerRecord, previous CostLedgerRecord, entryHash string) DurableAck {
	return DurableAck{
		valid:        true,
		scopeKind:    identity.ScopeKind,
		scopeHash:    identity.ScopeHash,
		snapshotHash: snapshot.LedgerHash,
		entryHash:    entryHash,
		confirmed:    previous.LedgerHash != "" || snapshot.LedgerHash != "",
	}
}

// run serialises in-process callers first and then takes the scope-wide
// cross-process lock for the whole read/decide/commit sequence.
func (ledger *DurableCostLedger) run(ctx context.Context, fn func() error) error {
	return ledger.runWithLock(ctx, true, fn)
}

// runRead is run for the read-only protocol: it joins the same scope lock but
// refuses to create a missing ledger or lock file.
func (ledger *DurableCostLedger) runRead(ctx context.Context, fn func() error) error {
	return ledger.runWithLock(ctx, false, fn)
}

func (ledger *DurableCostLedger) runWithLock(ctx context.Context, writable bool, fn func() error) error {
	if ctx != nil {
		if err := ctx.Err(); err != nil {
			return err
		}
	}
	ledger.mu.Lock()
	defer ledger.mu.Unlock()
	if ledger.closed {
		return fmt.Errorf("%w: durable cost ledger is closed", ErrInvalidCostLedger)
	}
	release, err := acquireDurableLock(ctx, ledger.lockPath, writable)
	if err != nil {
		return err
	}
	defer func() { _ = release() }()
	return fn()
}

// initMarkerPath is the durable "this scope was already initialized" marker.
// Its presence plus a missing canonical ledger is a fail-closed state: the
// initial allocation is never re-issued.
func (ledger *DurableCostLedger) initMarkerPath() string {
	return ledger.path + ".init"
}

func (ledger *DurableCostLedger) bootstrap() error {
	record, exists, err := ledger.readDiskLocked()
	if err != nil {
		return err
	}
	if exists {
		if err := ledger.verifyIdentityLocked(record); err != nil {
			return err
		}
		loaded, err := LoadCostLedger(record)
		if err != nil {
			return err
		}
		ledger.confirmed = copyCostEntries(record.Ledger.Entries)
		ledger.confirmedModel = loaded
		return nil
	}
	// No ledger on disk. Only a scope that was never initialized may create one:
	// a recovery that finds the ledger missing must fail closed.
	if _, statErr := os.Stat(ledger.initMarkerPath()); statErr == nil {
		return fmt.Errorf("%w: durable ledger %q is missing after initialization; refusing to re-issue the initial allocation", ErrInvalidCostLedger, ledger.path)
	} else if !errors.Is(statErr, os.ErrNotExist) {
		return statErr
	}
	if err := ledger.writeInitMarkerLocked(); err != nil {
		return err
	}
	fresh, err := NewCostLedger(ledger.identity)
	if err != nil {
		return err
	}
	if _, err := fresh.AppendAllocation(ledger.initial); err != nil {
		return err
	}
	snapshot, err := fresh.Snapshot()
	if err != nil {
		return err
	}
	if err := ledger.commitLocked(snapshot); err != nil {
		if !durableOutcomeVisible(err) {
			// Provably uncommitted: this attempt may be retried from scratch.
			_ = os.Remove(ledger.initMarkerPath())
		}
		return err
	}
	ledger.confirmed = copyCostEntries(snapshot.Ledger.Entries)
	ledger.confirmedModel = fresh
	return nil
}

// writeInitMarkerLocked creates the initialization marker exclusively and
// flushes it before the initial allocation is committed.
func (ledger *DurableCostLedger) writeInitMarkerLocked() error {
	payload := []byte("proofrail durable cost ledger initialization\n" +
		ledger.identity.LedgerID + "\n" + ledger.identity.ScopeKind + "\n" + ledger.identity.ScopeHash + "\n")
	path := ledger.initMarkerPath()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("%w: durable ledger %q is missing after initialization", ErrInvalidCostLedger, ledger.path)
		}
		return fmt.Errorf("%w: create initialization marker: %w", ErrDurableNotCommitted, err)
	}
	fail := func(cause error) error {
		_ = file.Close()
		_ = os.Remove(path)
		return fmt.Errorf("%w: initialization marker: %w", ErrDurableNotCommitted, cause)
	}
	if _, err := file.Write(payload); err != nil {
		return fail(err)
	}
	if err := file.Sync(); err != nil {
		return fail(err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("%w: initialization marker close: %w", ErrDurableNotCommitted, err)
	}
	if err := syncDurableDir(filepath.Dir(path)); err != nil {
		return fmt.Errorf("%w: initialization marker directory barrier: %w", ErrDurableNotCommitted, err)
	}
	return nil
}

// loadLocked re-reads the canonical ledger from disk (never a pre-replacement
// handle), verifies identity and that the confirmed prefix did not drift.
func (ledger *DurableCostLedger) loadLocked() (*CostLedger, CostLedgerRecord, error) {
	record, exists, err := ledger.readDiskLocked()
	if err != nil {
		return nil, CostLedgerRecord{}, err
	}
	if !exists {
		return nil, CostLedgerRecord{}, fmt.Errorf("%w: durable ledger %q is missing", ErrInvalidCostLedger, ledger.path)
	}
	if err := ledger.verifyIdentityLocked(record); err != nil {
		return nil, CostLedgerRecord{}, err
	}
	if err := ledger.verifyPrefixLocked(record); err != nil {
		return nil, CostLedgerRecord{}, err
	}
	loaded, err := LoadCostLedger(record)
	if err != nil {
		return nil, CostLedgerRecord{}, err
	}
	return loaded, record, nil
}

func (ledger *DurableCostLedger) readDiskLocked() (CostLedgerRecord, bool, error) {
	raw, err := os.ReadFile(ledger.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return CostLedgerRecord{}, false, nil
		}
		return CostLedgerRecord{}, false, err
	}
	record, err := DecodeCostLedgerRecord(raw)
	if err != nil {
		return CostLedgerRecord{}, false, err
	}
	return record, true, nil
}

func (ledger *DurableCostLedger) verifyIdentityLocked(record CostLedgerRecord) error {
	body := record.Ledger
	identity := ledger.identity
	if body.LedgerID != identity.LedgerID ||
		body.ScopeKind != identity.ScopeKind ||
		body.ScopeHash != identity.ScopeHash ||
		body.Currency != identity.Currency ||
		body.PricingMode != identity.PricingMode ||
		body.PricingVersion != identity.PricingVersion {
		return fmt.Errorf("%w: durable ledger identity drift", ErrInvalidCostLedger)
	}
	if !identity.CreatedAt.IsZero() && body.CreatedAt != identity.CreatedAt.UTC().Format(timestampLayout) {
		return fmt.Errorf("%w: durable ledger createdAt drift", ErrInvalidCostLedger)
	}
	return nil
}

func (ledger *DurableCostLedger) verifyPrefixLocked(record CostLedgerRecord) error {
	if len(record.Ledger.Entries) < len(ledger.confirmed) {
		return fmt.Errorf("%w: durable ledger shrank from %d to %d entries", ErrInvalidCostLedger, len(ledger.confirmed), len(record.Ledger.Entries))
	}
	if len(ledger.confirmed) == 0 {
		return nil
	}
	if !sameCostEntries(ledger.confirmed, record.Ledger.Entries[:len(ledger.confirmed)]) {
		return fmt.Errorf("%w: durable ledger history drift", ErrInvalidCostLedger)
	}
	return nil
}

// commitLocked performs the ordered durability sequence: full snapshot write,
// temp Sync, pre-commit Close, atomic replace, directory barrier. Every failure
// before the atomic replace is classified ErrDurableNotCommitted; an inconclusive
// replace is re-read and classified in-lock; a post-replace failure is
// classified ErrDurableVisibleUnconfirmed. No path removes a visible snapshot.
func (ledger *DurableCostLedger) commitLocked(snapshot CostLedgerRecord) error {
	if err := ValidateCostLedgerRecord(snapshot); err != nil {
		return fmt.Errorf("%w: %w", ErrDurableNotCommitted, err)
	}
	payload, err := evidence.EncodeCanonical(snapshot)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDurableNotCommitted, err)
	}
	targetDigest := evidence.Digest(durableLedgerBytesDomain, payload)
	if err := ledger.injectFault(StageWrite); err != nil {
		return fmt.Errorf("%w: snapshot write: %w", ErrDurableNotCommitted, err)
	}
	before := ledger.observeDiskLocked()
	dir := filepath.Dir(ledger.path)
	tmp, err := os.CreateTemp(dir, ".cost-ledger-*.tmp")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDurableNotCommitted, err)
	}
	tmpName := tmp.Name()
	discard := func() {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
	}
	written, err := tmp.Write(payload)
	if err != nil {
		discard()
		return fmt.Errorf("%w: %w", ErrDurableNotCommitted, err)
	}
	if written != len(payload) {
		discard()
		return fmt.Errorf("%w: short snapshot write (%d of %d bytes)", ErrDurableNotCommitted, written, len(payload))
	}
	if err := ledger.injectFault(StageSync); err != nil {
		discard()
		return fmt.Errorf("%w: snapshot sync: %w", ErrDurableNotCommitted, err)
	}
	if err := tmp.Sync(); err != nil {
		discard()
		return fmt.Errorf("%w: %w", ErrDurableNotCommitted, err)
	}
	if err := ledger.injectFault(StageClose); err != nil {
		discard()
		return fmt.Errorf("%w: snapshot close: %w", ErrDurableNotCommitted, err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("%w: %w", ErrDurableNotCommitted, err)
	}
	if err := ledger.injectFault(StageRename); err != nil {
		_ = os.Remove(tmpName)
		return ledger.classifyReplaceFailureLocked(targetDigest, before, err)
	}
	if err := os.Rename(tmpName, ledger.path); err != nil {
		_ = os.Remove(tmpName)
		return ledger.classifyReplaceFailureLocked(targetDigest, before, err)
	}
	if err := ledger.injectFault(StageDirBarrier); err != nil {
		return fmt.Errorf("%w: directory barrier: %w", ErrDurableVisibleUnconfirmed, err)
	}
	if err := syncDurableDir(dir); err != nil {
		return fmt.Errorf("%w: %w", ErrDurableVisibleUnconfirmed, err)
	}
	if err := ledger.injectFault(StageDirClose); err != nil {
		return fmt.Errorf("%w: directory close: %w", ErrDurableVisibleUnconfirmed, err)
	}
	return nil
}

// injectFault calls the optional test-only fault surface for one commit stage.
func (ledger *DurableCostLedger) injectFault(stage string) error {
	if ledger == nil || ledger.commitFault == nil {
		return nil
	}
	return ledger.commitFault(stage)
}

// durableDiskObservation is a byte-level view of the canonical ledger used to
// classify an inconclusive atomic replace. known is false when the read itself
// was inconclusive.
type durableDiskObservation struct {
	digest string
	exists bool
	known  bool
}

func (ledger *DurableCostLedger) observeDiskLocked() durableDiskObservation {
	raw, err := os.ReadFile(ledger.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return durableDiskObservation{known: true}
		}
		return durableDiskObservation{}
	}
	return durableDiskObservation{
		digest: evidence.Digest(durableLedgerBytesDomain, raw),
		exists: true,
		known:  true,
	}
}

// classifyReplaceFailureLocked re-reads the canonical path inside the same scope
// lock and decides which of the three outcomes the failed replace produced.
func (ledger *DurableCostLedger) classifyReplaceFailureLocked(targetDigest string, before durableDiskObservation, cause error) error {
	after := ledger.observeDiskLocked()
	switch {
	case !after.known:
		return fmt.Errorf("%w: atomic replace failed and the canonical ledger is unreadable: %w", ErrDurableOutcomeUnknown, cause)
	case after.exists && after.digest == targetDigest:
		return fmt.Errorf("%w: atomic replace failed but the target snapshot is visible: %w", ErrDurableVisibleUnconfirmed, cause)
	case !after.exists && before.known && !before.exists:
		return fmt.Errorf("%w: atomic replace failed, canonical ledger still absent: %w", ErrDurableNotCommitted, cause)
	case after.exists && before.known && before.exists && after.digest == before.digest:
		return fmt.Errorf("%w: atomic replace failed, previous snapshot still in place: %w", ErrDurableNotCommitted, cause)
	default:
		return fmt.Errorf("%w: atomic replace failed, neither snapshot observed: %w", ErrDurableOutcomeUnknown, cause)
	}
}

func copyCostEntries(entries []any) []any {
	if len(entries) == 0 {
		return nil
	}
	return append([]any{}, entries...)
}

func sameCostEntries(left, right []any) bool {
	leftBytes, err := evidence.EncodeCanonical(left)
	if err != nil {
		return false
	}
	rightBytes, err := evidence.EncodeCanonical(right)
	if err != nil {
		return false
	}
	return bytes.Equal(leftBytes, rightBytes)
}

// stableTerminalSettlementKey derives the settlement idempotency key for one
// terminal completion of one request.
func stableTerminalSettlementKey(requestID, completionHash string) string {
	hash := evidence.Digest("proofrail:agent-runner-settlement-key:1\n", []byte(requestID+"\n"+completionHash))
	hex := strings.TrimPrefix(hash, "sha256:")
	if len(hex) > 20 {
		hex = hex[:20]
	}
	return "settle-" + hex
}
