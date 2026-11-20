package guard

// P83 FIX ROUND 7 counterexamples (items T5c1, T5c2 and T5c3): every retained
// generation body is read through a bounded, link-refusing, regular-file handle,
// so a redirected, unbounded or non-regular body refuses the renewal instead of
// standing in for the earlier generation.
//
// Mutation contracts are stated on each test.

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

// p83Fix7ExecutableDomain is the digest domain the approved executable is bound in.
// It is a test constant, not a product domain, and it is shared by the
// windows-tagged (T3a/T3c) and linux-tagged (T3d) executable tests.
const p83Fix7ExecutableDomain = "proofrail:test-approved-executable:1\n"

// TestP83Fix7ExecutableDomainIsStable keeps the test-only domain referenced on
// every platform: the executable tests that use it are build-tagged to windows
// and linux, and an untagged constant that nothing else reads is dead code here.
func TestP83Fix7ExecutableDomainIsStable(t *testing.T) {
	if p83Fix7ExecutableDomain != "proofrail:test-approved-executable:1\n" {
		t.Fatalf("unexpected test digest domain %q", p83Fix7ExecutableDomain)
	}
}

// p83Fix7Retainable drafts one store, claims generation 1 and renews to
// generation 2, so the current record is exactly the body the NEXT renewal has to
// retain under its own hash.
func p83Fix7Retainable(t *testing.T, resource string) (*LeaseStore, Lease) {
	t.Helper()
	store, err := NewLeaseStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	process := ProcessIdentity{PID: 9001, StartToken: "p83-fix7-retain-holder"}
	claimed, err := store.Claim(ctx, resource, "worker-one", "claim-p83-fix7-"+resource, process, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	current, err := store.Renew(ctx, resource, claimed.Token, process, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if current.Token.Generation != claimed.Token.Generation+1 {
		t.Fatalf("the fixture must reach generation 2, got %d", current.Token.Generation)
	}
	return store, current
}

// TestP83Fix7T5c1RetainedSymlinkIsRefused is the T5c1 counterexample: a link at the
// retained-generation path refuses with the invalid-record error even when its
// target is a perfectly valid record carrying the hash the name asks for - the
// retention must describe the store's own object, never a target somebody chose.
//
// Mutation check: replacing the Lstat probe in retainGeneration with a follow-link
// existence check (or dropping the link refusal in openLeaseRecord) lets the target
// verify and the renewal succeed, which turns this red.
func TestP83Fix7T5c1RetainedSymlinkIsRefused(t *testing.T) {
	const resource = "workspace-p83-fix7-t5c1"
	store, current := p83Fix7Retainable(t, resource)
	ctx := context.Background()
	process := current.Process

	retainedPath := store.generationPath(resource, current.LeaseHash)
	// A valid body with exactly the hash the retained name carries, written OUTSIDE
	// the store, so following the link would verify perfectly.
	targetDirectory := t.TempDir()
	target := filepath.Join(targetDirectory, "outside-body.json")
	if err := writeLease(target, current); err != nil {
		t.Fatalf("write the outside body: %v", err)
	}
	if err := os.Symlink(target, retainedPath); err != nil {
		t.Skipf("this host cannot create a symlink: %v", err)
	}

	if _, err := store.Renew(ctx, resource, current.Token, process, time.Minute); !errors.Is(err, evidence.ErrInvalidRecord) {
		t.Fatalf("a linked retained generation must refuse the renewal with the invalid-record error, got %v", err)
	}
	after, err := store.Current(ctx, resource)
	if err != nil {
		t.Fatal(err)
	}
	if after.LeaseHash != current.LeaseHash || after.Token != current.Token {
		t.Fatalf("a refused retention must not advance the current record: got %+v (want %+v)", after, current)
	}

	// Control: with the link removed the same renewal retains the same body and
	// succeeds, so the refusal above was the link and nothing else.
	if err := os.Remove(retainedPath); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Renew(ctx, resource, current.Token, process, time.Minute); err != nil {
		t.Fatalf("the same renewal must succeed once the link is gone, got %v", err)
	}
}

// TestP83Fix7T5c2OversizedRetainedBodyIsRefused is the T5c2 counterexample: a
// retained body past the record bound refuses the renewal, so a hostile file cannot
// make the store read without end while it holds the resource lock.
//
// Mutation check: dropping the leaseRecordMaxBytes limit (or reading the file with
// os.ReadFile) makes the read succeed and turns the byte-bound assertion red.
func TestP83Fix7T5c2OversizedRetainedBodyIsRefused(t *testing.T) {
	const resource = "workspace-p83-fix7-t5c2"
	store, current := p83Fix7Retainable(t, resource)
	ctx := context.Background()

	// Valid JSON followed by padding: only the bound, not the parser, can refuse it
	// for its size.
	body, err := json.Marshal(current)
	if err != nil {
		t.Fatal(err)
	}
	oversized := append(append([]byte{}, body...), []byte(strings.Repeat(" ", int(leaseRecordMaxBytes)))...)
	retainedPath := store.generationPath(resource, current.LeaseHash)
	if err := os.WriteFile(retainedPath, oversized, 0o600); err != nil {
		t.Fatal(err)
	}

	_, err = store.Renew(ctx, resource, current.Token, current.Process, time.Minute)
	if !errors.Is(err, evidence.ErrInvalidRecord) {
		t.Fatalf("an oversized retained body must refuse with the invalid-record error, got %v", err)
	}
	if !strings.Contains(err.Error(), "byte bound") {
		t.Fatalf("the refusal must report the record byte bound, got %v", err)
	}
}

// TestP83Fix7T5c3NonRegularRetainedPathIsRefused is the T5c3 counterexample: a
// directory at the retained-generation path refuses with the invalid-record error
// before anything is read.
//
// Mutation check: dropping the `info.Mode().IsRegular()` check in readLeaseRecord
// turns the refusal into a read error of a different class and fails the message
// assertion below.
func TestP83Fix7T5c3NonRegularRetainedPathIsRefused(t *testing.T) {
	const resource = "workspace-p83-fix7-t5c3"
	store, current := p83Fix7Retainable(t, resource)
	ctx := context.Background()

	retainedPath := store.generationPath(resource, current.LeaseHash)
	if err := os.Mkdir(retainedPath, 0o700); err != nil {
		t.Fatal(err)
	}

	_, err := store.Renew(ctx, resource, current.Token, current.Process, time.Minute)
	if !errors.Is(err, evidence.ErrInvalidRecord) {
		t.Fatalf("a non-regular retained path must refuse with the invalid-record error, got %v", err)
	}
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("the refusal must report the non-regular type, got %v", err)
	}
}
