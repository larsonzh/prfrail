package guard

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

func TestLeaseFencingRenewAndRelease(t *testing.T) {
	store, err := NewLeaseStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	process := ProcessIdentity{PID: 101, StartToken: "start-one"}
	lease, err := store.Claim(context.Background(), "workspace-one", "worker-one", "claim-one", process, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Claim(context.Background(), "workspace-one", "worker-two", "claim-two", process, time.Minute); !errors.Is(err, ErrLeaseHeld) {
		t.Fatalf("second claim must fail: %v", err)
	}
	renewed, err := store.Renew(context.Background(), "workspace-one", lease.Token, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if renewed.Token.Generation != 2 || renewed.Token.ClaimID != lease.Token.ClaimID {
		t.Fatalf("unexpected renewed token: %+v", renewed.Token)
	}
	if err := store.Validate(context.Background(), "workspace-one", lease.Token); !errors.Is(err, ErrFencingTokenMismatch) {
		t.Fatalf("stale token accepted: %v", err)
	}
	if err := store.Release(context.Background(), "workspace-one", lease.Token); !errors.Is(err, ErrFencingTokenMismatch) {
		t.Fatalf("stale release accepted: %v", err)
	}
	if err := store.Release(context.Background(), "workspace-one", renewed.Token); err != nil {
		t.Fatal(err)
	}
}

func TestLeaseConcurrentClaimHasSingleWinner(t *testing.T) {
	store, _ := NewLeaseStore(t.TempDir())
	var wait sync.WaitGroup
	results := make(chan error, 16)
	for index := 0; index < 16; index++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			_, err := store.Claim(context.Background(), "workspace-race", fmt.Sprintf("worker-%d", index), fmt.Sprintf("claim-%d", index), ProcessIdentity{PID: index + 1, StartToken: "start"}, time.Minute)
			results <- err
		}(index)
	}
	wait.Wait()
	close(results)
	winners := 0
	for err := range results {
		if err == nil {
			winners++
		} else if !errors.Is(err, ErrLeaseHeld) {
			t.Fatalf("unexpected claim error: %v", err)
		}
	}
	if winners != 1 {
		t.Fatalf("got %d successful claims, want exactly one", winners)
	}
}

func TestLeaseTakeoverRequiresFreshStopAndAuthorization(t *testing.T) {
	authorization := evidence.Digest("", []byte("operator authorization"))
	store, _ := NewLeaseStore(t.TempDir(), WithTakeoverVerifier(func(_ Lease, proof TakeoverProof) error {
		if len(proof.AuthorizationEvidence) != 1 || proof.AuthorizationEvidence[0] != authorization {
			return errors.New("authorization evidence was not approved")
		}
		return nil
	}))
	fixedNow := time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)
	store.now = func() time.Time { return fixedNow }
	oldProcess := ProcessIdentity{PID: 201, StartToken: "old-start"}
	lease, err := store.Claim(context.Background(), "workspace-two", "worker-one", "claim-one", oldProcess, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	store.now = func() time.Time { return fixedNow.Add(time.Hour) }
	if _, err := store.Claim(context.Background(), "workspace-two", "worker-two", "claim-two", ProcessIdentity{PID: 202, StartToken: "new-start"}, time.Minute); !errors.Is(err, ErrLeaseHeld) {
		t.Fatalf("expiry alone allowed takeover: %v", err)
	}
	proof := makeTerminationEvidence(t, oldProcess, store.now())
	if _, err := store.Takeover(context.Background(), "workspace-two", "worker-two", "claim-two", ProcessIdentity{PID: 202, StartToken: "new-start"}, time.Minute, TakeoverProof{TerminationEvidence: proof}); !errors.Is(err, ErrTakeoverEvidence) {
		t.Fatalf("missing authorization accepted: %v", err)
	}
	wrongWriter := makeTerminationEvidence(t, ProcessIdentity{PID: 999, StartToken: "other"}, store.now())
	if _, err := store.Takeover(context.Background(), "workspace-two", "worker-two", "claim-two", ProcessIdentity{PID: 202, StartToken: "new-start"}, time.Minute, TakeoverProof{TerminationEvidence: wrongWriter, AuthorizationEvidence: []string{authorization}}); !errors.Is(err, ErrTakeoverEvidence) {
		t.Fatalf("wrong writer evidence accepted: %v", err)
	}
	stale := makeTerminationEvidence(t, oldProcess, store.now().Add(-10*time.Minute))
	if _, err := store.Takeover(context.Background(), "workspace-two", "worker-two", "claim-two", ProcessIdentity{PID: 202, StartToken: "new-start"}, time.Minute, TakeoverProof{TerminationEvidence: stale, AuthorizationEvidence: []string{authorization}}); !errors.Is(err, ErrTakeoverEvidence) {
		t.Fatalf("stale termination evidence accepted: %v", err)
	}
	taken, err := store.Takeover(context.Background(), "workspace-two", "worker-two", "claim-two", ProcessIdentity{PID: 202, StartToken: "new-start"}, time.Minute, TakeoverProof{TerminationEvidence: proof, AuthorizationEvidence: []string{authorization}})
	if err != nil {
		t.Fatal(err)
	}
	if taken.Token.Generation != lease.Token.Generation+1 || taken.Token.ClaimID == lease.Token.ClaimID {
		t.Fatalf("invalid takeover token: %+v", taken.Token)
	}
}

func TestLeaseTakeoverFailsClosedWithoutVerifier(t *testing.T) {
	store, _ := NewLeaseStore(t.TempDir())
	process := ProcessIdentity{PID: 301, StartToken: "old-start"}
	if _, err := store.Claim(context.Background(), "workspace-three", "worker-one", "claim-one", process, time.Minute); err != nil {
		t.Fatal(err)
	}
	proof := makeTerminationEvidence(t, process, store.now())
	_, err := store.Takeover(context.Background(), "workspace-three", "worker-two", "claim-two", ProcessIdentity{PID: 302, StartToken: "new-start"}, time.Minute, TakeoverProof{TerminationEvidence: proof, AuthorizationEvidence: []string{evidence.Digest("", []byte("authorization"))}})
	if !errors.Is(err, ErrTakeoverEvidence) {
		t.Fatalf("takeover without independent verifier must fail closed: %v", err)
	}
}

func makeTerminationEvidence(t *testing.T, identity ProcessIdentity, now time.Time) TerminationEvidence {
	t.Helper()
	result := TerminationEvidence{Identity: identity, Requested: formatTime(now.Add(-time.Second)), Verified: formatTime(now), Outcome: "stopped", Actions: []string{"verified"}}
	canonical, err := evidence.EncodeCanonical(struct {
		Identity  ProcessIdentity `json:"identity"`
		Requested string          `json:"requestedAt"`
		Verified  string          `json:"verifiedAt"`
		Outcome   string          `json:"outcome"`
		Actions   []string        `json:"actions"`
	}{result.Identity, result.Requested, result.Verified, result.Outcome, result.Actions})
	if err != nil {
		t.Fatal(err)
	}
	result.Hash = evidence.Digest("proofrail:termination-evidence:1\n", canonical)
	return result
}
