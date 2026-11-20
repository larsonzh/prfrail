package console

// P83 FIX ROUND 8 counterexamples (items T5#4a and T5#4b): a stale store-local
// proof is reconciled only when it is a verifiable earlier generation of the SAME
// claim. The record naming this process is the precondition; it is not consent to
// adopt any token a proof file carries.
//
// Mutation contracts are stated on each test.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

const p83Fix8RunID = "run-p83-fix8-t5c4"

// p83Fix8Superseded builds the exact state a generation commit leaves behind: the
// record moved one generation on while the proof on disk still names the one
// before it, and the record names THIS process as its holder.
func p83Fix8Superseded(t *testing.T) (string, *guard.LeaseStore, publishLeaseProof, guard.Lease) {
	t.Helper()
	runRoot := t.TempDir()
	ctx := context.Background()
	stale, err := claimPublishLease(ctx, runRoot, p83Fix8RunID)
	if err != nil {
		t.Fatalf("claim the publish lease: %v", err)
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	current, err := store.Renew(ctx, stale.ResourceID, stale.Token, publishLeaseProcess(p83Fix8RunID), publishLeaseTTL)
	if err != nil {
		t.Fatalf("commit the next generation: %v", err)
	}
	if current.Token.Generation != stale.Token.Generation+1 {
		t.Fatalf("test precondition failed: got generations %d -> %d", stale.Token.Generation, current.Token.Generation)
	}
	return runRoot, store, stale, current
}

// TestP83Fix8T5c4bGenuinePredecessorPublishProceeds is the T5#4b acceptance: a
// proof that IS the verified predecessor is reconciled, and the publish proceeds on
// the committed generation.
//
// Mutation check: deleting the reconcile block in renewPublishLease (returning the
// refusal for every stale proof) turns this red.
func TestP83Fix8T5c4bGenuinePredecessorPublishProceeds(t *testing.T) {
	runRoot, store, stale, current := p83Fix8Superseded(t)
	if err := verifyPublishLeasePredecessor(store, stale, current); err != nil {
		t.Fatalf("a genuine predecessor must verify, got %v", err)
	}
	ctx := context.Background()
	refreshed, err := renewPublishLease(ctx, runRoot, p83Fix8RunID, stale, time.Now().UTC())
	if err != nil {
		t.Fatalf("a genuine predecessor must be reconciled, got %v", err)
	}
	if refreshed.Token != current.Token || refreshed.LeaseHash != current.LeaseHash {
		t.Fatalf("the reconciled proof must carry the committed generation: got %+v want %+v", refreshed.Token, current.Token)
	}
	bound, err := verifyPublishLease(ctx, runRoot, p83Fix8RunID, refreshed, time.Now().UTC())
	if err != nil || bound != current.LeaseHash {
		t.Fatalf("the publish must proceed on the committed hash, got %q (err %v)", bound, err)
	}
}

// TestP83Fix8T5c4aUnverifiablePredecessorIsRefused is the T5#4a counterexample: a
// stale proof that is not a verifiable predecessor - wrong claim, generation not
// N-1, a hash the current record does not name, or a hash the store never retained
// - is refused even though the record names THIS process.
//
// Mutation check: replacing the verifyPublishLeasePredecessor call in
// renewPublishLease (or its body) with an unconditional acceptance lets every case
// below be reconciled and turns the end-to-end refusal assertion red.
func TestP83Fix8T5c4aUnverifiablePredecessorIsRefused(t *testing.T) {
	runRoot, store, stale, current := p83Fix8Superseded(t)
	ctx := context.Background()

	forged := evidence.Digest("proofrail:test-forged-lease:1\n", []byte("never issued by this store"))
	cases := []struct {
		name  string
		proof publishLeaseProof
		want  string
	}{
		{
			name:  "wrong-claim-id",
			proof: publishLeaseProof{ResourceID: stale.ResourceID, Token: guard.FencingToken{ClaimID: "run-p83-fix8-forged", Generation: stale.Token.Generation}, LeaseHash: stale.LeaseHash},
			want:  "not the current claim",
		},
		{
			name:  "previous-hash-mismatch",
			proof: publishLeaseProof{ResourceID: stale.ResourceID, Token: stale.Token, LeaseHash: current.LeaseHash},
			want:  "does not name",
		},
		{
			name:  "fabricated-hash",
			proof: publishLeaseProof{ResourceID: stale.ResourceID, Token: stale.Token, LeaseHash: forged},
			want:  "does not name",
		},
	}
	for _, testCase := range cases {
		if err := verifyPublishLeasePredecessor(store, testCase.proof, current); err == nil {
			t.Fatalf("%s: an unverifiable predecessor must be refused", testCase.name)
		} else if !strings.Contains(err.Error(), testCase.want) {
			t.Fatalf("%s: the refusal must say %q, got %v", testCase.name, testCase.want, err)
		}
	}

	// End to end: the same wrong-claim proof must not be reconciled, and the record
	// must not have moved.
	refusalProof := cases[0].proof
	if _, err := renewPublishLease(ctx, runRoot, p83Fix8RunID, refusalProof, time.Now().UTC()); !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a wrong-claim stale proof must keep the refusal, got %v", err)
	}
	after, err := store.Current(ctx, stale.ResourceID)
	if err != nil {
		t.Fatal(err)
	}
	if after.LeaseHash != current.LeaseHash || after.Token != current.Token {
		t.Fatalf("a refused reconciliation must not move the record: got %+v want %+v", after, current)
	}

	// Generation not N-1: a three-generation run whose proof is two generations old.
	_, store3, twoBehind, currentThree := p83Fix8ThreeGenerations(t)
	if twoBehind.Token.Generation+1 == currentThree.Token.Generation {
		t.Fatal("test precondition failed: the proof must not be the immediate predecessor")
	}
	if err := verifyPublishLeasePredecessor(store3, twoBehind, currentThree); err == nil {
		t.Fatal("a proof two generations behind must be refused")
	} else if !strings.Contains(err.Error(), "is not the generation before") {
		t.Fatalf("the refusal must name the generation gap, got %v", err)
	}

	// Retained-record proof: the predecessor names the right hash, but the store did
	// not retain it, so the claim cannot be verified from the store's own records.
	runRootRetained, storeRetained, staleRetained, currentRetained := p83Fix8Superseded(t)
	// The retained body sits next to the current record, named by its own hash.
	retainedPath := filepath.Join(publishLeaseRoot(runRootRetained), staleRetained.ResourceID+"."+strings.TrimPrefix(staleRetained.LeaseHash, "sha256:")+".lease.json")
	if err := os.Remove(retainedPath); err != nil {
		t.Fatal(err)
	}
	if err := verifyPublishLeasePredecessor(storeRetained, staleRetained, currentRetained); err == nil {
		t.Fatal("an unretained predecessor must be refused")
	} else if !strings.Contains(err.Error(), "retained records") {
		t.Fatalf("the refusal must name the retained records, got %v", err)
	}
}

// p83Fix8ThreeGenerations claims and renews twice, so the current record is three
// generations past the first proof.
func p83Fix8ThreeGenerations(t *testing.T) (string, *guard.LeaseStore, publishLeaseProof, guard.Lease) {
	t.Helper()
	runRoot := t.TempDir()
	ctx := context.Background()
	first, err := claimPublishLease(ctx, runRoot, p83Fix8RunID)
	if err != nil {
		t.Fatalf("claim the publish lease: %v", err)
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	holder := publishLeaseProcess(p83Fix8RunID)
	second, err := store.Renew(ctx, first.ResourceID, first.Token, holder, publishLeaseTTL)
	if err != nil {
		t.Fatalf("commit generation 2: %v", err)
	}
	third, err := store.Renew(ctx, second.ResourceID, second.Token, holder, publishLeaseTTL)
	if err != nil {
		t.Fatalf("commit generation 3: %v", err)
	}
	return runRoot, store, first, third
}
