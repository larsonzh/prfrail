package guard

// P83 FIX ROUND 6 counterexamples (items T5b and T5c).
//
// Every test below is written so that reverting the one guard it names turns it
// red; the mutation contract is stated in each test's doc comment.

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// TestP83Fix6T5bRenewFromAForeignProcessIsRefused is the T5b counterexample: the
// fencing token alone does not prove ownership, so a renewal presented by a
// process that is not the recorded holder must be refused with the holder
// mismatch instead of extending the lease.
//
// Mutation check: deleting the `current.Process != process` comparison in
// LeaseStore.Renew lets the foreign identity renew the record and turns this red.
func TestP83Fix6T5bRenewFromAForeignProcessIsRefused(t *testing.T) {
	store, err := NewLeaseStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	const resource = "workspace-p83-fix6-t5b"
	holder := ProcessIdentity{PID: 401, StartToken: "p83-fix6-holder"}
	claimed, err := store.Claim(ctx, resource, "worker-one", "claim-p83-fix6-t5b", holder, time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	foreign := ProcessIdentity{PID: 402, StartToken: "p83-fix6-foreign"}
	if _, err := store.Renew(ctx, resource, claimed.Token, foreign, time.Minute); !errors.Is(err, ErrLeaseHolderMismatch) {
		t.Fatalf("a renewal by a non-holder must refuse with the holder mismatch, got %v", err)
	}
	// The refusal must leave the record byte-identical: no generation advance, no
	// new holder and no new lease hash.
	current, err := store.Current(ctx, resource)
	if err != nil {
		t.Fatal(err)
	}
	if current.LeaseHash != claimed.LeaseHash || current.Token != claimed.Token || current.Process != holder {
		t.Fatalf("a refused renewal must not touch the record: got %+v (want %+v)", current, claimed)
	}

	// Positive control: the recorded holder still renews, so the refusal above is
	// the ownership proof and not a general renewal failure.
	renewed, err := store.Renew(ctx, resource, claimed.Token, holder, time.Minute)
	if err != nil {
		t.Fatalf("the holder must still be able to renew, got %v", err)
	}
	if renewed.Token.Generation != claimed.Token.Generation+1 || renewed.Process != holder {
		t.Fatalf("the holder renewal must advance the generation and keep the holder, got %+v", renewed.Token)
	}
}

// TestP83Fix6T5cCorruptRetainedGenerationBlocksTheRenewal is the T5c
// counterexample: an existing retained generation body is only acceptable when it
// re-reads as the record whose hash names it. A corrupt body must refuse the
// renewal and leave the current record in place instead of standing in for an
// earlier generation.
//
// Mutation check: replacing the read-back verification in retainGeneration with
// `return nil` (path-exists treated as success) lets the renewal overwrite the
// current record and turns this red.
func TestP83Fix6T5cCorruptRetainedGenerationBlocksTheRenewal(t *testing.T) {
	store, err := NewLeaseStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	const resource = "workspace-p83-fix6-t5c"
	process := ProcessIdentity{PID: 501, StartToken: "p83-fix6-t5c-holder"}
	claimed, err := store.Claim(ctx, resource, "worker-one", "claim-p83-fix6-t5c", process, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	// The first renewal retains the claimed generation, so retention itself is
	// proven to work before the corruption leg is judged.
	current, err := store.Renew(ctx, resource, claimed.Token, process, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.retainedLease(resource, claimed.LeaseHash); err != nil {
		t.Fatalf("the replaced generation must be retained and verifiable, got %v", err)
	}

	// Corrupt the body that the CURRENT record's hash names: the next renewal has
	// to retain exactly this record, and the occupied name no longer verifies.
	retainedPath := store.generationPath(resource, current.LeaseHash)
	if err := os.WriteFile(retainedPath, []byte("{\"schemaVersion\":\"1.0.0\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Renew(ctx, resource, current.Token, process, time.Minute); !errors.Is(err, evidence.ErrInvalidRecord) {
		t.Fatalf("a corrupt retained generation must refuse the renewal, got %v", err)
	}
	// The renewal is refused before the write, so the current record did not move.
	after, err := store.Current(ctx, resource)
	if err != nil {
		t.Fatal(err)
	}
	if after.LeaseHash != current.LeaseHash || after.Token != current.Token {
		t.Fatalf("a refused retention must not overwrite the current record: got %+v (want %+v)", after, current)
	}

	// Control: with the corrupt body removed the same renewal succeeds, so the
	// refusal was the corruption and nothing else.
	if err := os.Remove(retainedPath); err != nil {
		t.Fatal(err)
	}
	next, err := store.Renew(ctx, resource, current.Token, process, time.Minute)
	if err != nil {
		t.Fatalf("the same renewal must succeed once the retained body is gone, got %v", err)
	}
	if next.Token.Generation != current.Token.Generation+1 {
		t.Fatalf("the control renewal must advance the generation, got %+v", next.Token)
	}
	// The control renewal retained the previously corrupt generation as a valid
	// body, which is what makes it independently re-verifiable again.
	if _, err := store.retainedLease(resource, current.LeaseHash); err != nil {
		t.Fatalf("the re-written retained generation must verify, got %v", err)
	}
}
