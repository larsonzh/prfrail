package console

// P83 FIX ROUND 7 counterexamples (items T5d and T5e): a stale local proof is
// re-synced only when the committed record names THIS process, and an unreadable
// proof fails the renewal batch instead of being skipped as success.
//
// Mutation contracts are stated on each test.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

// TestP83Fix7T5dStaleProofHeldByThisProcessIsResynced is the liveness half of
// T5d: a generation commit that was followed by an interrupted proof rewrite leaves
// a task on the superseded token, but the record still names this process - so the
// proof is re-synced to the committed generation and the publish proceeds.
//
// Mutation check: deleting the re-sync block in renewPublishLease (returning the
// ErrPublishLeaseUnproven refusal instead) turns this red.
func TestP83Fix7T5dStaleProofHeldByThisProcessIsResynced(t *testing.T) {
	runRoot := t.TempDir()
	ctx := context.Background()
	stale, err := claimPublishLease(ctx, runRoot, p83RunID)
	if err != nil {
		t.Fatal(err)
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	// The generation commit: exactly what a threshold renewal does to the record.
	committed, err := store.Renew(ctx, stale.ResourceID, stale.Token, publishLeaseProcess(p83RunID), publishLeaseTTL)
	if err != nil {
		t.Fatalf("commit a new generation: %v", err)
	}
	if committed.LeaseHash == stale.LeaseHash || committed.Token == stale.Token {
		t.Fatal("test precondition failed: the commit must supersede the stored proof")
	}

	refreshed, err := renewPublishLease(ctx, runRoot, p83RunID, stale, time.Now().UTC())
	if err != nil {
		t.Fatalf("a stale proof whose record names this process must be re-synced, got %v", err)
	}
	if refreshed.Token != committed.Token || refreshed.LeaseHash != committed.LeaseHash {
		t.Fatalf("the re-synced proof must carry the committed generation, got %+v want %+v", refreshed.Token, committed.Token)
	}
	// The publish proceeds on the re-synced proof, which is what "the record IS this
	// process's authority" has to mean end to end.
	bound, err := verifyPublishLease(ctx, runRoot, p83RunID, refreshed, time.Now().UTC())
	if err != nil || bound != committed.LeaseHash {
		t.Fatalf("the re-synced proof must verify and bind the committed hash, got %q (err %v)", bound, err)
	}
}

// TestP83Fix7T5dStaleProofHeldByAnotherProcessIsRefused is the refusal half of
// T5d: the same stale shape whose record names ANOTHER process keeps the refusal,
// even when that other process is provably stopped and the shared reclaim path
// would otherwise hand the record over.
//
// Mutation check: deleting the `lease.Process != holder` refusal inside the re-sync
// block lets the proof be re-synced onto the foreign record and then reclaimed, so
// the call returns a proof instead of the refusal and turns this red.
func TestP83Fix7T5dStaleProofHeldByAnotherProcessIsRefused(t *testing.T) {
	runRoot := t.TempDir()
	ctx := context.Background()
	const foreignRun = "run-p83-fix7-t5d-foreign"
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	// A real record of a DIFFERENT run, held by an identity that is provably gone:
	// the shared reclaim proof below would accept it.
	stopped := guard.ProcessIdentity{PID: 2147483646, StartToken: "p83-fix7-t5d-stopped"}
	resource := publishLeaseResourceID(foreignRun)
	if _, err := store.Claim(ctx, resource, publishLeaseConsumerID, "run-p83-fix7-t5d-claim", stopped, publishLeaseTTL); err != nil {
		t.Fatalf("seed the foreign record: %v", err)
	}
	stale := publishLeaseProof{
		ResourceID: resource,
		Token:      guard.FencingToken{ClaimID: "run-p83-fix7-t5d-stale-claim", Generation: 1},
		LeaseHash:  evidence.Digest("", []byte("p83-fix7-t5d-stale-proof")),
	}

	_, err = renewPublishLease(ctx, runRoot, foreignRun, stale, time.Now().UTC())
	if !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a stale proof whose record names another process must be refused, got %v", err)
	}
	// The refusal must not have transferred the record: it still names the stopped
	// holder, so a later successful reclaim is the only way it may move.
	current, err := store.Current(ctx, resource)
	if err != nil {
		t.Fatal(err)
	}
	if current.Process != stopped {
		t.Fatalf("a refused re-sync must not touch the record, got holder %+v", current.Process)
	}
}

// TestP83Fix7T5eUnreadableProofFailsTheRenewalBatch is the T5e counterexample: a
// store-local proof that cannot be parsed fails the whole renewal batch, so no task
// is left behind on a superseded token while the call reports success.
//
// Mutation check: replacing the `return err` after readSupervisedJSON with
// `continue` skips the unreadable proof and returns nil, which turns this red.
func TestP83Fix7T5eUnreadableProofFailsTheRenewalBatch(t *testing.T) {
	runRoot := t.TempDir()
	ctx := context.Background()
	proof, err := claimPublishLease(ctx, runRoot, p83RunID)
	if err != nil {
		t.Fatal(err)
	}
	supervised := SupervisedDir(runRoot)
	goodDirectory := filepath.Join(supervised, "task-a")
	badDirectory := filepath.Join(supervised, "task-b")
	for _, directory := range []string{goodDirectory, badDirectory} {
		if err := os.MkdirAll(directory, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := writeSupervisedJSON(filepath.Join(goodDirectory, publishLeaseProofName), proof); err != nil {
		t.Fatal(err)
	}
	badPath := filepath.Join(badDirectory, publishLeaseProofName)
	if err := os.WriteFile(badPath, []byte("{not-a-proof"), 0o600); err != nil {
		t.Fatal(err)
	}

	runtime := &supervisedRuntime{runID: p83RunID, runRoot: runRoot}
	refreshed := publishLeaseProof{
		ResourceID: proof.ResourceID,
		Token:      guard.FencingToken{ClaimID: proof.Token.ClaimID, Generation: proof.Token.Generation + 1},
		LeaseHash:  evidence.Digest("", []byte("p83-fix7-t5e-refreshed")),
	}
	if err := runtime.persistRenewedPublishLease(proof, refreshed); err == nil {
		t.Fatal("an unreadable store-local proof must fail the renewal batch")
	} else if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("the batch must report the unreadable proof, got %v", err)
	}
	// Control: the readable proof alone would have been rewritten, so the error above
	// is the unreadable sibling and not a general batch failure.
	if err := os.Remove(badPath); err != nil {
		t.Fatal(err)
	}
	if err := runtime.persistRenewedPublishLease(proof, refreshed); err != nil {
		t.Fatalf("the same batch must succeed once the unreadable proof is gone, got %v", err)
	}
}
