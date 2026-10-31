package gates

// P44-T027B-GATE-READINESS (N4): a probe permit is a one-shot quota.
//
// N4 recorded that a reusable permit lets the same permit hand out a fresh
// quota on every call, so a probe could be re-admitted after a failure, a
// timeout or a process restart. The permit therefore linearises through a
// single durable on-disk claim:
//
//	unused --Require--> claimed --Complete--> spent
//
// Only Require can leave unused, only Complete can leave claimed, and no path
// returns to unused: a failed probe, a cancelled context or a restarted process
// keeps the quota consumed (fail-closed).
//
// Linearisation point: creating the marker file with O_CREATE|O_EXCL and
// flushing it with Sync. The in-memory atomic state is only a same-process
// fast-reject cache in front of that file; it is never the admission gate.
// Cross-process exclusion falls out of O_CREATE|O_EXCL itself, so two processes
// racing on the same path cannot both win.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

// Sentinel errors of the probe permit. Callers match them with errors.Is, but
// not every returned error carries a sentinel, and the two shapes must be told
// apart.
//
// Permit-state errors wrap exactly one sentinel (ErrProbePermitMissing,
// ErrProbePermitInUse or ErrProbePermitConsumed). These describe the state of a
// permit.
//
// Context cancellation and deadline expiry are passed through from ctx.Err(), so
// callers match them with errors.Is(err, context.Canceled) or
// errors.Is(err, context.DeadlineExceeded); they wrap no permit sentinel.
//
// I/O failures are not uniform, and the boundary matters. "Returned as-is, with
// no permit sentinel" applies only to read failures inside LoadProbePermit (for
// example a marker body that cannot be read): the permit state is still whatever
// the marker says, so the read error is reported unchanged and carries no permit
// sentinel.
//
// Persistence failures inside Require have a different shape and are NOT passed
// through. The durable marker is the linearisation point of the permit, so a
// claim that cannot be created (empty path, a marker already owned by another
// claim, or MkdirAll, OpenFile, Write or Sync failing) must never degrade into a
// silent admission. That is treated as a permit-state error: the permit is
// invalidated (spent) and the returned error wraps ErrProbePermitConsumed. This
// is fail-closed on purpose.
//
// Do not assume that every error returned by this file wraps a probe permit
// sentinel.
var (
	// ErrProbePermitMissing marks a probe conclusion (or a permit operation)
	// that has no bound permit to authorise it.
	ErrProbePermitMissing = errors.New("probe permit missing")
	// ErrProbePermitInUse marks a permit whose single claim is already held by
	// an attempt that has not completed.
	ErrProbePermitInUse = errors.New("probe permit already in use")
	// ErrProbePermitConsumed marks a permit whose quota is spent for good.
	ErrProbePermitConsumed = errors.New("probe permit already consumed")
)

const probePermitMarkerSchema = "1.0.0"

// probePermitState is the closed state set of a permit. The numeric values are
// load-bearing: they are the compare-and-swap operands used by Complete, and the
// values Require's fast-reject cache compares against.
type probePermitState uint32

const (
	// probePermitUnused is the only state from which Require may claim.
	probePermitUnused probePermitState = 0
	// probePermitClaimed means one attempt owns the quota right now.
	probePermitClaimed probePermitState = 1
	// probePermitSpent means the quota can never be claimed again.
	probePermitSpent probePermitState = 2
)

// probePermitMarker is the durable record written when a claim is granted. The
// existence of the marker file, not its body, is the cross-process claim bit:
// LoadProbePermit treats any marker it finds as an already spent permit.
type probePermitMarker struct {
	SchemaVersion string `json:"schemaVersion"`
	Kind          string `json:"kind"`
	ID            string `json:"id"`
	State         string `json:"state"`
	MaxRequests   int64  `json:"maxRequests"`
	ClaimedAt     string `json:"claimedAt"`
}

// ProbePermit is a single-use admission token for a probe run. Create one with
// NewProbePermit (fresh quota) or LoadProbePermit (quota read back from disk).
// The zero value is not usable: it carries no identity and no marker path.
//
// ProbePermit must not be copied after first use because it embeds an atomic.
type ProbePermit struct {
	id          string
	maxRequests int64
	state       atomic.Uint32
	path        string
}

// NewProbePermit returns an Unused permit bound to id whose claim is persisted
// at path. maxRequests is normalised to at least 1: a probe can never be
// admitted with a zero or negative quota.
//
// Neither id nor path is validated here. With an empty path there is nowhere to
// persist the claim marker, so the claim can never become durable: Require
// spends the permit and returns ErrProbePermitConsumed. A missing marker store
// fails closed and can never degrade into a silent admission.
//
// The marker file belongs to the caller for the whole life of the permit.
// Deleting the marker re-issues the quota to the next claim, so production
// paths must never delete it.
func NewProbePermit(id, path string, maxRequests int64) *ProbePermit {
	if maxRequests < 1 {
		maxRequests = 1
	}
	return &ProbePermit{id: id, maxRequests: maxRequests, path: path}
}

// LoadProbePermit rebuilds a permit from its marker file.
//
// A missing marker is not an error: it is a permit that was never claimed, so
// the returned permit is Unused and still bound to path. An existing marker
// means the quota was already claimed by this or another process, so the
// returned permit is Spent and Require refuses with ErrProbePermitConsumed -
// the restart must not reset the quota.
func LoadProbePermit(path string) (*ProbePermit, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("%w: probe permit marker path is required", ErrProbePermitMissing)
	}
	permit := &ProbePermit{id: filepath.Base(path), maxRequests: 1, path: path}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Zero value of state is probePermitUnused: a brand-new quota.
			return permit, nil
		}
		return nil, fmt.Errorf("load probe permit %q: %w", path, err)
	}

	// Fail closed: an unreadable body must never demote a claimed permit.
	var marker probePermitMarker
	if json.Unmarshal(data, &marker) == nil {
		if strings.TrimSpace(marker.ID) != "" {
			permit.id = marker.ID
		}
		if marker.MaxRequests > 0 {
			permit.maxRequests = marker.MaxRequests
		}
	}
	permit.state.Store(uint32(probePermitSpent))
	return permit, nil
}

// Require reserves the quota for one probe attempt. Persistence comes first:
// the durable marker is the linearisation point, so Require never admits an
// attempt whose claim is not already on disk.
//
// Order of operations:
//  1. pre-checks - nil permit, nil context, and an already cancelled or expired
//     context all refuse without touching the filesystem;
//  2. fast reject - a cache that already reads claimed or spent refuses with
//     ErrProbePermitInUse / ErrProbePermitConsumed and does no disk I/O;
//  3. persistClaim - O_CREATE|O_EXCL, Write and Sync of the marker file. Only a
//     successful persist admits: the cache is then set to claimed and nil is
//     returned. A marker that already exists (errors.Is(err, fs.ErrExist)) means
//     another claim holds the quota, and any other persistence failure means the
//     claim cannot be made durable: both invalidate the permit (spent) and
//     return an error wrapping ErrProbePermitConsumed.
//
// Invariant: Require returns nil if and only if the marker file is already on
// disk (Sync completed). Every crash therefore lands in one of exactly two
// shapes:
//
//	(a) Require did not return nil - no probe executed, so the quota was never
//	    used and loading the permit after a restart yields Unused, which is
//	    the correct reading of "nothing happened";
//	(b) the marker is on disk - a restart loads the permit as Spent and no
//	    further claim can be granted.
//
// There is no window in which a probe has already executed and the quota can
// still be brought back to life.
//
// Cancellation and timeout are reported straight through: the returned error
// satisfies errors.Is(err, context.Canceled) or
// errors.Is(err, context.DeadlineExceeded) and wraps no permit sentinel.
// Cancelling before the claim leaves the quota unconsumed; cancelling after the
// marker is durable keeps it consumed, because an abandoned attempt never
// refunds.
//
// Two constructor behaviours surface here. An empty path cannot persist the
// claim, so the permit lands in Spent and this call returns
// ErrProbePermitConsumed (fail-closed, never a silent admission). And the marker
// file is the quota itself: deleting it re-issues the quota to the next claim,
// so production paths must never delete it.
func (p *ProbePermit) Require(ctx context.Context) error {
	if p == nil {
		return fmt.Errorf("%w: nil probe permit", ErrProbePermitMissing)
	}
	if ctx == nil {
		return fmt.Errorf("%w: probe permit %q require needs a context", ErrProbePermitMissing, p.id)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("probe permit %q: require: %w", p.id, err)
	}

	// Fast reject only: the cache can be ahead of the disk, never behind it, so a
	// claimed or spent cache refuses without racing the filesystem for the claim.
	switch probePermitState(p.state.Load()) {
	case probePermitClaimed, probePermitSpent:
		return p.refusal()
	}

	// The marker file, not the cache, is the linearisation point: persisting the
	// claim is what makes an admission final.
	if err := p.persistClaim(); err != nil {
		// No durable claim means no admission. Record the loss in the cache so this
		// process stops asking, and fail closed towards the caller.
		p.state.Store(uint32(probePermitSpent))
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("%w: probe permit %q: marker %q is already held by another claim: %v",
				ErrProbePermitConsumed, p.id, p.path, err)
		}
		return fmt.Errorf("%w: probe permit %q: %v", ErrProbePermitConsumed, p.id, err)
	}
	p.state.Store(uint32(probePermitClaimed))

	if err := ctx.Err(); err != nil {
		// The attempt was abandoned after the claim became durable; the quota stays
		// consumed (an abandoned attempt keeps the claim, it never refunds).
		return fmt.Errorf("probe permit %q: require: %w", p.id, err)
	}
	return nil
}

// Complete marks a claimed permit as spent. It is idempotent for a permit that
// is already spent, and refuses (ErrProbePermitInUse) a permit that was never
// required, so completion can never invent a quota.
func (p *ProbePermit) Complete() error {
	if p == nil {
		return fmt.Errorf("%w: nil probe permit", ErrProbePermitMissing)
	}
	if p.state.CompareAndSwap(uint32(probePermitClaimed), uint32(probePermitSpent)) {
		return nil
	}
	switch probePermitState(p.state.Load()) {
	case probePermitSpent:
		return nil
	case probePermitUnused:
		return fmt.Errorf("%w: probe permit %q completed before it was required", ErrProbePermitInUse, p.id)
	default:
		return fmt.Errorf("%w: probe permit %q in unknown state", ErrProbePermitConsumed, p.id)
	}
}

// refusal classifies the lost claim race by the state the winner left behind.
func (p *ProbePermit) refusal() error {
	switch probePermitState(p.state.Load()) {
	case probePermitClaimed:
		return fmt.Errorf("%w: probe permit %q is held by an in-flight attempt", ErrProbePermitInUse, p.id)
	case probePermitSpent:
		return fmt.Errorf("%w: probe permit %q quota is spent", ErrProbePermitConsumed, p.id)
	default:
		return fmt.Errorf("%w: probe permit %q lost the claim race", ErrProbePermitInUse, p.id)
	}
}

// persistClaim writes the durable claim marker. O_EXCL makes the file itself
// the cross-process lock: exactly one writer can create it.
func (p *ProbePermit) persistClaim() error {
	if strings.TrimSpace(p.path) == "" {
		return errors.New("probe permit marker path is empty")
	}
	if dir := filepath.Dir(p.path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	marker := probePermitMarker{
		SchemaVersion: probePermitMarkerSchema,
		Kind:          "probe-permit",
		ID:            p.id,
		State:         "claimed",
		MaxRequests:   p.maxRequests,
		ClaimedAt:     time.Now().UTC().Format(time.RFC3339Nano),
	}
	data, err := json.Marshal(marker)
	if err != nil {
		return err
	}
	data = append(data, '\n')

	file, err := os.OpenFile(p.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if _, err := file.Write(data); err != nil {
		return err
	}
	return file.Sync()
}
