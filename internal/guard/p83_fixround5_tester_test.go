package guard

// P83 FIX ROUND 5 counterexample (item T5b): a renewal must not destroy the
// generation it replaced, so a receipt bound to the previous lease hash stays
// re-verifiable through the retained record.
//
// Mutation check: deleting the store.retainGeneration call in withLock removes the
// retained body, retainedLease fails with ErrNotExist and this turns red.

import (
	"context"
	"testing"
	"time"
)

func TestP83T5bRetainedGenerationStaysReVerifiable(t *testing.T) {
	store, err := NewLeaseStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	process := ProcessIdentity{PID: 401, StartToken: "p83-fix5-start"}
	claimed, err := store.Claim(ctx, "workspace-t5b", "worker-one", "claim-t5b", process, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	renewed, err := store.Renew(ctx, "workspace-t5b", claimed.Token, process, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if renewed.LeaseHash == claimed.LeaseHash {
		t.Fatal("test precondition failed: a renewal must move the lease hash")
	}

	retained, err := store.retainedLease("workspace-t5b", claimed.LeaseHash)
	if err != nil {
		t.Fatalf("the replaced generation must stay addressable by its hash, got %v", err)
	}
	if retained.LeaseHash != claimed.LeaseHash || retained.Token != claimed.Token {
		t.Fatalf("retainedLease must return the previous generation, got %+v (want token %+v)", retained, claimed.Token)
	}
	// A second renewal must retain the second generation too, so the retention is a
	// property of every replaced record and not of the first one only.
	twice, err := store.Renew(ctx, "workspace-t5b", renewed.Token, process, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.retainedLease("workspace-t5b", renewed.LeaseHash)
	if err != nil {
		t.Fatalf("the second replaced generation must stay addressable, got %v", err)
	}
	if second.LeaseHash != renewed.LeaseHash || second.Token != renewed.Token {
		t.Fatalf("retainedLease must return the second generation, got %+v", second)
	}
	if twice.Token.Generation != renewed.Token.Generation+1 {
		t.Fatalf("the current generation must keep advancing, got %+v", twice.Token)
	}

	// A hash the store never wrote is refused instead of standing in for a body.
	missing := "sha256:0000000000000000000000000000000000000000000000000000000000000000"
	if _, err := store.retainedLease("workspace-t5b", missing); err == nil {
		t.Fatal("an unwritten generation hash must be refused")
	}
}
