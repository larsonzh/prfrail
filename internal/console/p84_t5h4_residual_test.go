package console

// P84 T5#4 residual counterexamples (items C1, C2, C3) for the store-local publish
// proof batch rewrite.
//
// C1: a proof that names this run's lease resource but is NOT a genuine superseded
// generation of it (another claim, a hash the store never retained) must fail the
// batch and keep its bytes. Resource identity alone is not consent to rewrite.
//
// C2: proofs that ARE genuine earlier generations of the same claim are brought onto
// the current record, including the multi-generation skew several tasks of one run
// can accumulate. The rewrite produces the refreshed bytes and the caught-up proof
// still verifies at promotion.
//
// C3: the generation commit and the proof rewrite are two writes. A process that
// commits a generation and dies before rewriting the proofs leaves them behind; the
// next process adopts the lease and the batch recovers those proofs through the
// retained history. The same-process leg is already pinned by
// internal/console/p83_fixround8_tester_test.go:56-83 and is registered only.
//
// Mutations that must turn these tests red are stated on each test.
//
// Round 2 of the tester adds the two findings of the independent scan:
//
// C6: the store record is the authority for the batch. A caller whose "refreshed"
// claim a concurrent advance has already left must fail the batch with zero
// rewrites (the M1 defect, fixed in product fix round 5).
//
// C7: the retained-history walk of verifyPublishLeaseHistory is bounded by
// publishLeaseHistoryLimit, and a genuine proof beyond that bound is refused with
// the depth error (the M2 boundary, registered rather than fixed, pinned here).

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

const p84T5h4RunID = "run-p84-t5h4-residual"

// p84T5h4State is one run's lease store with a superseded proof (stale) and the
// record the store currently carries (current), renewals generations apart.
type p84T5h4State struct {
	runRoot string
	store   *guard.LeaseStore
	stale   publishLeaseProof
	current publishLeaseProof
}

// p84T5h4Build claims the run's lease and commits renewals generations, so the
// first proof is renewals generations behind the current record.
func p84T5h4Build(t *testing.T, renewals int) *p84T5h4State {
	t.Helper()
	runRoot := t.TempDir()
	ctx := context.Background()
	first, err := claimPublishLease(ctx, runRoot, p84T5h4RunID)
	if err != nil {
		t.Fatalf("claim the publish lease: %v", err)
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	holder := publishLeaseProcess(p84T5h4RunID)
	proof := first
	for index := 0; index < renewals; index++ {
		renewed, err := store.Renew(ctx, proof.ResourceID, proof.Token, holder, publishLeaseTTL)
		if err != nil {
			t.Fatalf("commit renewal %d: %v", index+1, err)
		}
		proof = publishLeaseProof{ResourceID: renewed.ResourceID, Token: renewed.Token, LeaseHash: renewed.LeaseHash}
	}
	if renewals > 0 && proof.Token.Generation != first.Token.Generation+uint64(renewals) {
		t.Fatalf("test precondition failed: got generations %d -> %d for %d renewals", first.Token.Generation, proof.Token.Generation, renewals)
	}
	return &p84T5h4State{runRoot: runRoot, store: store, stale: first, current: proof}
}

func (state *p84T5h4State) runtime() *supervisedRuntime {
	return &supervisedRuntime{runID: p84T5h4RunID, runRoot: state.runRoot}
}

// proofPath is the store-local proof file of one task directory inside the
// supervised tree, exactly where persistRenewedPublishLease looks for it.
func (state *p84T5h4State) proofPath(taskDirectory string) string {
	return filepath.Join(SupervisedDir(state.runRoot), taskDirectory, publishLeaseProofName)
}

// p84T5h4ProofBytes is the exact byte image writeSupervisedJSON leaves behind.
func p84T5h4ProofBytes(t *testing.T, proof publishLeaseProof) []byte {
	t.Helper()
	canonical, err := evidence.EncodeCanonical(proof)
	if err != nil {
		t.Fatal(err)
	}
	return append(canonical, '\n')
}

func p84T5h4ReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

// TestP84T5h4C1ForeignProofIsNotRewritten is the C1 counterexample: a proof of this
// run's resource that is not a genuine earlier generation of the current claim
// (another claim id, or a hash the retained chain never carried) fails the batch,
// and its bytes are untouched.
//
// Mutation check: reducing verifyPublishLeaseHistory to a ResourceID comparison
// (the pre-P84 shape) accepts both cases, rewrites the file and turns the refusal
// and the byte assertions red.
func TestP84T5h4C1ForeignProofIsNotRewritten(t *testing.T) {
	forged := evidence.Digest("proofrail:p84-test-forged-lease:1\n", []byte("never issued by this store"))
	cases := []struct {
		name  string
		proof func(state *p84T5h4State) publishLeaseProof
	}{
		{
			name: "another-claim",
			proof: func(state *p84T5h4State) publishLeaseProof {
				return publishLeaseProof{
					ResourceID: state.stale.ResourceID,
					Token:      guard.FencingToken{ClaimID: "run-p84-t5h4-forged", Generation: state.stale.Token.Generation},
					LeaseHash:  state.stale.LeaseHash,
				}
			},
		},
		{
			name: "unretained-hash",
			proof: func(state *p84T5h4State) publishLeaseProof {
				return publishLeaseProof{
					ResourceID: state.stale.ResourceID,
					Token:      state.stale.Token,
					LeaseHash:  forged,
				}
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			state := p84T5h4Build(t, 1)
			ctx := context.Background()
			path := state.proofPath("task-foreign")
			foreign := testCase.proof(state)
			if err := writeSupervisedJSON(path, foreign); err != nil {
				t.Fatal(err)
			}
			before := p84T5h4ReadFile(t, path)

			err := state.runtime().persistRenewedPublishLease(ctx, state.store, state.stale, state.current)
			if !errors.Is(err, ErrPublishLeaseUnproven) {
				t.Fatalf("a non-historical proof of this resource must fail the batch, got %v", err)
			}
			if !strings.Contains(err.Error(), publishLeaseProofName) {
				t.Fatalf("the batch error must name the failing proof file, got %v", err)
			}
			if after := p84T5h4ReadFile(t, path); !bytes.Equal(after, before) {
				t.Fatalf("a refused proof must keep its bytes, got %q want %q", after, before)
			}

			// Control: the genuine superseded generation alone is accepted and rewritten,
			// so the refusal above is the history check and nothing else.
			if err := writeSupervisedJSON(path, state.stale); err != nil {
				t.Fatal(err)
			}
			if err := state.runtime().persistRenewedPublishLease(ctx, state.store, state.stale, state.current); err != nil {
				t.Fatalf("the same batch must succeed for a genuine superseded proof, got %v", err)
			}
			if got, want := p84T5h4ReadFile(t, path), p84T5h4ProofBytes(t, state.current); !bytes.Equal(got, want) {
				t.Fatalf("a recovered proof must carry the refreshed bytes, got %q want %q", got, want)
			}
		})
	}
}

// TestP84T5h4C2MultiGenerationHistoryIsCaughtUp is the C2 counterexample: proofs
// that are genuine earlier generations of the same claim - two generations behind
// here, which is the skew several tasks accumulate - are brought onto the current
// record, and the caught-up proof still verifies at promotion.
//
// Mutation check: replacing the verifyPublishLeaseHistory call in
// persistRenewedPublishLease with verifyPublishLeasePredecessor (only N-1 accepted)
// refuses the two-generation-old proof and turns this red.
func TestP84T5h4C2MultiGenerationHistoryIsCaughtUp(t *testing.T) {
	state := p84T5h4Build(t, 2)
	if gap := state.current.Token.Generation - state.stale.Token.Generation; gap < 2 {
		t.Fatalf("test precondition failed: the proof must be at least two generations behind, got %d", gap)
	}
	ctx := context.Background()
	directories := []string{"task-a", "task-b"}
	for _, directory := range directories {
		if err := writeSupervisedJSON(state.proofPath(directory), state.stale); err != nil {
			t.Fatal(err)
		}
	}

	if err := state.runtime().persistRenewedPublishLease(ctx, state.store, state.stale, state.current); err != nil {
		t.Fatalf("genuine earlier generations must be caught up: %v", err)
	}
	want := p84T5h4ProofBytes(t, state.current)
	for _, directory := range directories {
		if got := p84T5h4ReadFile(t, state.proofPath(directory)); !bytes.Equal(got, want) {
			t.Fatalf("%s must carry the refreshed bytes, got %q want %q", directory, got, want)
		}
	}
	// The recovered generation is the one a promotion binds.
	bound, err := verifyPublishLease(ctx, state.runRoot, p84T5h4RunID, state.current, time.Now().UTC())
	if err != nil || bound != state.current.LeaseHash {
		t.Fatalf("the caught-up proof must verify at promotion, got hash %q err %v", bound, err)
	}
}

// TestP84T5h4C3CrossProcessReentryRecoversProof is the C3 cross-process leg: one
// process commits a generation and dies before rewriting a task's proof; the next
// process adopts the lease (the predecessor is proven stopped) and its batch brings
// that proof onto the generation the adoption committed, after which the proof
// verifies at promotion.
//
// Mutation check: restricting the batch to N-1 (verifyPublishLeasePredecessor)
// leaves the two-generation-old proof unrecoverable and turns this red; returning
// the pre-transfer record from claimPublishLease leaves the record naming the dead
// process, so the promotion verification below goes red.
func TestP84T5h4C3CrossProcessReentryRecoversProof(t *testing.T) {
	runRoot := t.TempDir()
	child := exec.Command(os.Args[0], "-test.run=^TestP84T5h4C3ChildProcess$")
	child.Env = append(os.Environ(),
		"P84_T5H4_C3_RUN_ROOT="+runRoot,
		"P84_T5H4_C3_RUN_ID="+p84T5h4RunID,
	)
	var childOutput bytes.Buffer
	child.Stdout = &childOutput
	child.Stderr = &childOutput
	runErr := child.Run()
	childPID := 0
	if child.Process != nil {
		childPID = child.Process.Pid
		// The exited process must stop being addressable, or the adoption proof - which
		// refuses a holder whose stop cannot be disproved - cannot be exercised at all.
		_ = child.Process.Release()
	}
	if runErr != nil {
		t.Fatalf("the interrupted child must commit its generation: %v (%s)", runErr, strings.TrimSpace(childOutput.String()))
	}
	if !p84T5h4WaitUninspectable(childPID, 3*time.Second) {
		t.Skipf("this host keeps the exited child pid %d addressable, so its stop cannot be proven", childPID)
	}
	childGeneration, err := strconv.ParseUint(strings.TrimSpace(string(p84T5h4ReadFile(t, filepath.Join(runRoot, "child-generation.txt")))), 10, 64)
	if err != nil {
		t.Fatalf("read the child's generation: %v", err)
	}

	ctx := context.Background()
	adopted, err := claimPublishLease(ctx, runRoot, p84T5h4RunID)
	if err != nil {
		if strings.Contains(err.Error(), "previous publish lease holder") {
			t.Skipf("this host cannot prove the child holder stopped: %v", err)
		}
		t.Fatalf("the re-entry must adopt the lease of the stopped child: %v", err)
	}
	if adopted.Token.Generation <= childGeneration {
		t.Fatalf("the adoption must transfer the record to a new generation, got %d after the child's %d", adopted.Token.Generation, childGeneration)
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	runtime := &supervisedRuntime{runID: p84T5h4RunID, runRoot: runRoot}
	if err := runtime.persistRenewedPublishLease(ctx, store, publishLeaseProof{}, adopted); err != nil {
		t.Fatalf("the batch must recover the proof the crashed process left behind: %v", err)
	}
	path := filepath.Join(SupervisedDir(runRoot), "task-crashed", publishLeaseProofName)
	if got, want := p84T5h4ReadFile(t, path), p84T5h4ProofBytes(t, adopted); !bytes.Equal(got, want) {
		t.Fatalf("the recovered proof must carry the adopted generation, got %q want %q", got, want)
	}
	if bound, err := verifyPublishLease(ctx, runRoot, p84T5h4RunID, adopted, time.Now().UTC()); err != nil || bound != adopted.LeaseHash {
		t.Fatalf("the recovered proof must verify at promotion, got hash %q err %v", bound, err)
	}
}

// p84T5h4WaitUninspectable reports whether pid stopped resolving as a process
// within the bound. A process that exited can stay addressable while any handle on
// its object is still open, so the caller retries with a GC cycle in between.
func p84T5h4WaitUninspectable(pid int, bound time.Duration) bool {
	deadline := time.Now().Add(bound)
	for {
		if _, err := guard.InspectProcess(pid); err != nil {
			return true
		}
		if time.Now().After(deadline) {
			return false
		}
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
	}
}

// TestP84T5h4C3ChildProcess is the process that commits a generation and stops
// before rewriting the task proof. Every other run of this function skips.
func TestP84T5h4C3ChildProcess(t *testing.T) {
	runRoot := os.Getenv("P84_T5H4_C3_RUN_ROOT")
	if runRoot == "" {
		t.Skip("helper process for the P84 T5#4 cross-process re-entry test")
	}
	runID := os.Getenv("P84_T5H4_C3_RUN_ID")
	ctx := context.Background()
	stale, err := claimPublishLease(ctx, runRoot, runID)
	if err != nil {
		t.Fatalf("claim the publish lease: %v", err)
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	committed, err := store.Renew(ctx, stale.ResourceID, stale.Token, publishLeaseProcess(runID), publishLeaseTTL)
	if err != nil {
		t.Fatalf("commit the generation: %v", err)
	}
	// The commit is durable, the proof rewrite is not: this is the crash window.
	if err := writeSupervisedJSON(filepath.Join(SupervisedDir(runRoot), "task-crashed", publishLeaseProofName), stale); err != nil {
		t.Fatalf("leave the superseded proof behind: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runRoot, "child-generation.txt"), []byte(strconv.FormatUint(committed.Token.Generation, 10)), 0o600); err != nil {
		t.Fatalf("record the committed generation: %v", err)
	}
	os.Exit(0)
}

// TestP84T5h4C6ConcurrentAdvanceRefusesTheWholeBatch is the M1 counterexample of
// the independent scan: the caller renewed the lease and then a concurrent actor
// advanced the store one generation further, so the record the caller holds as
// "refreshed" is no longer the store's current one. The batch must read the
// store's record, see that the caller's claim is not it, and refuse the whole
// batch with zero rewrites - otherwise every proof would be stamped with a
// superseded generation and the projection would stop being a linearisation of
// the committed history.
//
// Mutation check: deleting the `current.Token != refreshed.Token ||
// current.LeaseHash != refreshed.LeaseHash` refusal in persistRenewedPublishLease
// (runtime_supervised.go, fix round 5) makes the call return nil and leave the
// proofs on the generation the caller passed in, which turns the refusal and the
// byte assertions below red.
func TestP84T5h4C6ConcurrentAdvanceRefusesTheWholeBatch(t *testing.T) {
	// One renewal: the store carries generation 2 and "stale" is generation 1, so a
	// later commit still leaves the stale proof reachable through the retained chain.
	state := p84T5h4Build(t, 1)
	ctx := context.Background()
	advanced, err := state.store.Renew(ctx, state.current.ResourceID, state.current.Token, publishLeaseProcess(p84T5h4RunID), publishLeaseTTL)
	if err != nil {
		t.Fatalf("commit the concurrent advance: %v", err)
	}
	if advanced.Token.Generation != state.current.Token.Generation+1 {
		t.Fatalf("test precondition failed: the concurrent advance must move the store one generation past the caller's, got %d after %d", advanced.Token.Generation, state.current.Token.Generation)
	}

	path := state.proofPath("task-raced")
	if err := writeSupervisedJSON(path, state.stale); err != nil {
		t.Fatal(err)
	}
	before := p84T5h4ReadFile(t, path)

	// The caller still believes generation 2 is current; the store says generation 3.
	err = state.runtime().persistRenewedPublishLease(ctx, state.store, state.stale, state.current)
	if !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a batch naming a generation the store has already left must be refused, got %v", err)
	}
	if !strings.Contains(err.Error(), "the store's current") {
		t.Fatalf("the refusal must be the record equality check, not a later per-proof check, got %v", err)
	}
	if after := p84T5h4ReadFile(t, path); !bytes.Equal(after, before) {
		t.Fatalf("a refused batch must rewrite nothing, got %q want %q", after, before)
	}

	// Control: the same proof is rewritten when the batch names the generation the
	// store actually carries, so the refusal above is the equality check and not the
	// history check refusing a generation the chain does not reach.
	refreshed := publishLeaseProof{ResourceID: advanced.ResourceID, Token: advanced.Token, LeaseHash: advanced.LeaseHash}
	if err := state.runtime().persistRenewedPublishLease(ctx, state.store, state.current, refreshed); err != nil {
		t.Fatalf("the batch must succeed against the store's current generation, got %v", err)
	}
	if got, want := p84T5h4ReadFile(t, path), p84T5h4ProofBytes(t, refreshed); !bytes.Equal(got, want) {
		t.Fatalf("the caught-up proof must carry the store's current bytes, got %q want %q", got, want)
	}
}

// TestP84T5h4C7HistoryDepthIsFailClosed pins the registered boundary of
// verifyPublishLeaseHistory (M2 of the independent scan, not fixed in this slice):
// the retained-history walk is bounded by publishLeaseHistoryLimit, and a genuine
// proof that sits further behind than the bound is refused with the depth error
// even though claim, generation and hash all match. Refusing is the intended
// fail-closed behaviour; widening the bound needs its own slice assessment.
//
// The chain is built for real, one store.Renew per generation, so every retained
// body is a record the store itself wrote and verified: the case exercises the
// depth branch of verifyPublishLeaseHistory and not a shape, claim, retained or
// token rejection. The immediate predecessor is rewritten successfully first, as a
// control that the chain this test builds is verifiable end to end.
//
// No skip conditions: the construction is a linear number of small file writes in
// the test's own temporary directory. If a host is ever slow enough for this to
// matter, the substitute (hand-built retained bodies chained by PreviousHash with
// self-proving hashes) would have to keep the bodies bit-identical to what
// writeLease produces, or the walk would fail on the retained check instead of the
// depth check and the case would stop proving anything.
//
// Mutation check: raising publishLeaseHistoryLimit to 256 (or any bound above the
// chain length) lets the walk reach the deepest proof, rewrites it, and turns the
// refusal and the byte assertions below red.
func TestP84T5h4C7HistoryDepthIsFailClosed(t *testing.T) {
	ctx := context.Background()
	runRoot := t.TempDir()
	first, err := claimPublishLease(ctx, runRoot, p84T5h4RunID)
	if err != nil {
		t.Fatalf("claim the publish lease: %v", err)
	}
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	holder := publishLeaseProcess(p84T5h4RunID)
	// One proof per generation: proofs[0] is the claim, proofs[len-1] the current record.
	proofs := []publishLeaseProof{first}
	current := first
	// One renewal past the registered bound of 64 (publishLeaseHistoryLimit). The
	// length is a literal on purpose: deriving it from the constant would let a
	// widened bound widen this chain with it, and the case would keep passing while
	// proving nothing.
	const renewals = 65
	started := time.Now()
	for index := 0; index < renewals; index++ {
		renewed, err := store.Renew(ctx, current.ResourceID, current.Token, holder, publishLeaseTTL)
		if err != nil {
			t.Fatalf("commit renewal %d of the deep chain: %v", index+1, err)
		}
		current = publishLeaseProof{ResourceID: renewed.ResourceID, Token: renewed.Token, LeaseHash: renewed.LeaseHash}
		proofs = append(proofs, current)
	}
	t.Logf("built %d retained generations in %s", renewals, time.Since(started))
	if gap := current.Token.Generation - first.Token.Generation; gap != uint64(renewals) {
		t.Fatalf("test precondition failed: the chain must span %d generations, got %d", renewals, gap)
	}

	supervisedState := &p84T5h4State{runRoot: runRoot, store: store, stale: first, current: current}
	runtime := supervisedState.runtime()

	// Control: the immediate predecessor is one hop from the current record and is
	// rewritten, so the chain is well formed and the bound is the only obstacle.
	predecessor := proofs[len(proofs)-2]
	controlPath := supervisedState.proofPath("task-predecessor")
	if err := writeSupervisedJSON(controlPath, predecessor); err != nil {
		t.Fatal(err)
	}
	if err := runtime.persistRenewedPublishLease(ctx, store, predecessor, current); err != nil {
		t.Fatalf("the immediate predecessor of a valid chain must be caught up, got %v", err)
	}
	if got, want := p84T5h4ReadFile(t, controlPath), p84T5h4ProofBytes(t, current); !bytes.Equal(got, want) {
		t.Fatalf("the predecessor control must carry the refreshed bytes, got %q want %q", got, want)
	}

	// The deepest proof is a genuine generation of this claim, renewals hops behind
	// the current record: one hop past the bound.
	deepestPath := supervisedState.proofPath("task-deepest")
	if err := writeSupervisedJSON(deepestPath, first); err != nil {
		t.Fatal(err)
	}
	before := p84T5h4ReadFile(t, deepestPath)
	err = runtime.persistRenewedPublishLease(ctx, store, first, current)
	if !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a proof %d generations behind must be refused by the bounded walk, got %v", renewals, err)
	}
	if !strings.Contains(err.Error(), "depth check failed") {
		t.Fatalf("the refusal must be the depth branch, not a shape/claim/retained/token rejection, got %v", err)
	}
	if after := p84T5h4ReadFile(t, deepestPath); !bytes.Equal(after, before) {
		t.Fatalf("a depth-refused proof must keep its bytes, got %q want %q", after, before)
	}
	// The refused generation is still a real record: the store reads it back, so the
	// depth error came from the walk bound and not from an unverifiable retained body.
	if retained, err := store.RetainedLease(first.ResourceID, first.LeaseHash); err != nil {
		t.Fatalf("the deepest generation must stay verifiable in the store: %v", err)
	} else if retained.Token != first.Token {
		t.Fatalf("the retained deepest generation must carry the proof's token, got %s/%d", retained.Token.ClaimID, retained.Token.Generation)
	}
}

// Round 3 of the tester adds the counterexamples for the blind review's F2/F3/F4
// findings:
//
// C8 (F2): a stale proof of a claim whose record names a DEAD predecessor process is
// recovered by transferring the record (Reclaim behind the stopped-holder proof)
// instead of being refused because the record names another process.
//
// C9 (F3): the lease lock is an OS-level lock on ONE contended name, so a holder that
// dies - crash included - has its lock released by the kernel and the next claimant
// never waits for a lock nobody owns. There is no heal, no pid marker and no liveness
// query to get wrong.
//
// C10 (F4): the write-before-publish equality check accepts only the store's current
// record.

// p84LockChild is one independent holder process started for the lease-lock cases.
type p84LockChild struct {
	cmd        *exec.Cmd
	output     *bytes.Buffer
	markerPath string
	waited     bool
	waitErr    error
}

// p84StartLockChild starts this test binary as a SEPARATE process that contends for one
// resource's lease lock exactly the way a second writer would.
//
// mode "die" takes the lock, writes the marker and dies without releasing, so from then
// on only the kernel can free the lock. mode "wait" takes the lock, writes the marker
// and blocks until the parent creates the release marker, which is a holder that stays
// alive. mode "claim" is the ordinary writer shape: Claim the lease and exit; Claim
// already returns the lock itself, so that leg is about what a committed record does to
// the next claimant.
func p84StartLockChild(t *testing.T, root, resourceID, markerPath, mode string) *p84LockChild {
	t.Helper()
	command := exec.Command(os.Args[0], "-test.run=^TestP84LeaseLockChildProcess$")
	command.Env = append(os.Environ(),
		"P84_LOCK_CHILD_ROOT="+root,
		"P84_LOCK_CHILD_RESOURCE="+resourceID,
		"P84_LOCK_CHILD_MARKER="+markerPath,
		"P84_LOCK_CHILD_MODE="+mode,
	)
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	if err := command.Start(); err != nil {
		t.Fatalf("start the lease lock helper: %v", err)
	}
	child := &p84LockChild{cmd: command, output: &output, markerPath: markerPath}
	t.Cleanup(child.terminate)
	return child
}

// terminate reaps the child when the case failed before it exited, so no helper
// outlives its test.
func (child *p84LockChild) terminate() {
	if child.waited {
		return
	}
	_ = child.cmd.Process.Kill()
	child.waitErr = child.cmd.Wait()
	child.waited = true
}

// waitExit reaps the child and reports its exit error. A watchdog kills the helper, so a
// stuck case surfaces as a failure instead of holding the suite.
func (child *p84LockChild) waitExit(t *testing.T) error {
	t.Helper()
	if child.waited {
		return child.waitErr
	}
	watchdog := time.AfterFunc(60*time.Second, func() { _ = child.cmd.Process.Kill() })
	defer watchdog.Stop()
	child.waitErr = child.cmd.Wait()
	child.waited = true
	return child.waitErr
}

// p84WaitForPath polls for a marker file, bounded by bound.
func p84WaitForPath(t *testing.T, path string, bound time.Duration) {
	t.Helper()
	deadline := time.Now().Add(bound)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if !time.Now().Before(deadline) {
			t.Fatalf("the helper never wrote %s within %s", path, bound)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

// TestP84LeaseLockChildProcess is the independent holder both lease-lock cases start.
// It writes the marker only after the resource lock is really held, so the parent's
// belief about the holder's state is an observation rather than a sleep. Every other
// run skips.
func TestP84LeaseLockChildProcess(t *testing.T) {
	root := os.Getenv("P84_LOCK_CHILD_ROOT")
	resourceID := os.Getenv("P84_LOCK_CHILD_RESOURCE")
	markerPath := os.Getenv("P84_LOCK_CHILD_MARKER")
	mode := os.Getenv("P84_LOCK_CHILD_MODE")
	if root == "" || resourceID == "" || markerPath == "" || mode == "" {
		t.Skip("helper process for the P84 lease-lock constructions")
	}
	store, err := guard.NewLeaseStore(root)
	if err != nil {
		t.Fatalf("open the lease store: %v", err)
	}
	if mode == "claim" {
		if _, err := store.Claim(context.Background(), resourceID, "p84-lock-child-consumer", "p84-lock-child-claim", publishLeaseProcess(p84T5h4RunID), publishLeaseTTL); err != nil {
			t.Fatalf("the helper must claim the lease: %v", err)
		}
		os.Exit(0)
	}
	if err := store.RunLocked(context.Background(), resourceID, func(guard.Lease) error {
		if err := os.WriteFile(markerPath, []byte(mode+"\n"), 0o600); err != nil {
			t.Fatalf("write the held marker: %v", err)
		}
		if mode == "die" {
			// Die while the lock is held: no release, no cleanup. From here on the
			// kernel owns the release.
			os.Exit(0)
		}
		releasePath := markerPath + ".release"
		deadline := time.Now().Add(60 * time.Second)
		for {
			if _, statErr := os.Stat(releasePath); statErr == nil {
				return nil
			}
			if !time.Now().Before(deadline) {
				t.Fatalf("the parent never released the live holder")
			}
			time.Sleep(5 * time.Millisecond)
		}
	}); err != nil {
		t.Fatalf("the helper must hold the resource lock: %v", err)
	}
}

// p84C8DeadIdentity starts a helper, lets it record its own OS identity and exit,
// and returns that identity once the host proves the process is gone. The identity
// is the real one the guard inspected, so the stopped-holder proof can be satisfied
// the same way it is satisfied in production.
func p84C8DeadIdentity(t *testing.T, directory string) guard.ProcessIdentity {
	t.Helper()
	path := filepath.Join(directory, "dead-identity.json")
	child := exec.Command(os.Args[0], "-test.run=^TestP84C8DeadIdentityChild$")
	child.Env = append(os.Environ(), "P84_C8_IDENTITY_PATH="+path)
	var output bytes.Buffer
	child.Stdout = &output
	child.Stderr = &output
	if err := child.Run(); err != nil {
		t.Fatalf("the identity helper must exit cleanly: %v (%s)", err, strings.TrimSpace(output.String()))
	}
	pid := child.Process.Pid
	if !p84T5h4WaitUninspectable(pid, 3*time.Second) {
		t.Skipf("this host keeps the exited helper pid %d addressable, so its stop cannot be proven", pid)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read the helper identity: %v", err)
	}
	var identity guard.ProcessIdentity
	if err := json.Unmarshal(data, &identity); err != nil {
		t.Fatalf("decode the helper identity: %v", err)
	}
	if identity.PID != pid || identity.StartToken == "" {
		t.Fatalf("the helper identity must name the helper process, got %+v want pid %d", identity, pid)
	}
	return identity
}

// TestP84C8DeadIdentityChild records this process's OS identity and exits. Every
// other run skips.
func TestP84C8DeadIdentityChild(t *testing.T) {
	path := os.Getenv("P84_C8_IDENTITY_PATH")
	if path == "" {
		t.Skip("helper process for the P84 C8 dead-holder construction")
	}
	identity, err := guard.InspectProcess(os.Getpid())
	if err != nil {
		t.Fatalf("inspect the helper process: %v", err)
	}
	data, err := json.Marshal(identity)
	if err != nil {
		t.Fatalf("encode the helper identity: %v", err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write the helper identity: %v", err)
	}
	os.Exit(0)
}

// TestP84C8DeadHolderReentryTransfersTheLease is the F2 cross-process counterexample
// this slice exists for. The record belongs to a process that has since died, one
// generation past the task proof that process left behind. The re-entering process
// presents that stale proof: the proof is real history of the claim, so the record
// is transferred to this process (Reclaim behind the stopped-holder proof) and the
// returned proof names the generation AFTER the dead holder's, which then verifies
// at promotion.
//
// Mutation check: restoring the pre-F2 mismatch branch - refusing whenever
// lease.Process != holder before any transfer - returns the "current lease ... is not
// the one this run claimed" error and turns the recovery assertion red. Deleting the
// history walk instead refuses the stale proof for the wrong reason, which the
// message assertion also catches.
func TestP84C8DeadHolderReentryTransfersTheLease(t *testing.T) {
	runRoot := t.TempDir()
	identity := p84C8DeadIdentity(t, runRoot)
	ctx := context.Background()
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	resourceID := publishLeaseResourceID(p84T5h4RunID)
	// The dead process claimed generation 1 and committed generation 2 (the crash
	// window: the commit is durable, the task proof is not rewritten).
	first, err := store.Claim(ctx, resourceID, publishLeaseConsumerID, publishLeaseClaimID(p84T5h4RunID), identity, publishLeaseTTL)
	if err != nil {
		t.Fatalf("claim the publish lease as the dead holder: %v", err)
	}
	second, err := store.Renew(ctx, resourceID, first.Token, identity, publishLeaseTTL)
	if err != nil {
		t.Fatalf("commit the dead holder's generation: %v", err)
	}
	if second.Token.Generation != first.Token.Generation+1 {
		t.Fatalf("test precondition failed: got generations %d -> %d", first.Token.Generation, second.Token.Generation)
	}
	stale := publishLeaseProof{ResourceID: resourceID, Token: first.Token, LeaseHash: first.LeaseHash}

	refreshed, err := renewPublishLease(ctx, runRoot, p84T5h4RunID, stale, time.Now().UTC())
	if err != nil {
		t.Fatalf("the stale proof of a dead holder must be recovered by a transfer, got %v", err)
	}
	if refreshed.Token.Generation != second.Token.Generation+1 {
		t.Fatalf("the transfer must advance past the dead holder's record, got generation %d after %d", refreshed.Token.Generation, second.Token.Generation)
	}
	current, err := store.Current(ctx, resourceID)
	if err != nil {
		t.Fatalf("read the transferred record: %v", err)
	}
	if current.Process != publishLeaseProcess(p84T5h4RunID) {
		t.Fatalf("the transfer must name this process as the holder, got %+v", current.Process)
	}
	if current.Token != refreshed.Token || current.LeaseHash != refreshed.LeaseHash {
		t.Fatalf("the returned proof must be the transferred record, got %+v", refreshed)
	}
	// The recovered proof is the one a promotion binds.
	if bound, err := verifyPublishLease(ctx, runRoot, p84T5h4RunID, refreshed, time.Now().UTC()); err != nil || bound != refreshed.LeaseHash {
		t.Fatalf("the recovered proof must verify at promotion, got hash %q err %v", bound, err)
	}
}

// TestP84C9DeadHolderLockIsReleasedByTheKernel pins N1/N2 against the OS-level lock:
// the lock is an exclusive kernel lock on ONE contended name, "<record>.lock", held
// through a handle and not through a marker in the file. A holder that dies - crashed
// or killed, without ever calling release - therefore leaves a lock FILE behind but no
// held lock: the kernel drops the lock when the holding process ends. The next claimant
// takes the resource in one shot instead of waiting for an owner that no longer exists,
// and nothing has to be healed, because there is no pid marker to read and no liveness
// query to answer.
//
// The lock file deliberately MAY survive the case: release never unlinks it (unlinking a
// lock file a waiter already has open would let a third process create a fresh file and
// hold a second, independent lock on the same resource), and no consumer enumerates it.
// The assertion is about the kernel's exclusion, never about the file's existence.
//
// The claim runs under a 2s context, so a lock that is NOT freed fails with the context
// error instead of hanging the suite.
//
// Mutation check (shared with C11, as recorded on purpose): the primary mutation is
// deleting the OS lock acquisition in createLeaseLock - the LockFileEx call on Windows,
// the Flock call on POSIX - which removes exclusion altogether. C11 turns red there
// (its live holder stops blocking the parent) while this case cannot: with no lock taken
// there is nothing left to recover, so the dead holder's release is trivially true. The
// shape that turns C9 red is the substrate this lock replaced - a holder identity kept in
// the file that the next claimant waits on with neither a liveness proof nor a heal,
// which returns context.DeadlineExceeded below.
func TestP84C9DeadHolderLockIsReleasedByTheKernel(t *testing.T) {
	const (
		resourceID = "ws-p84-c9-lock-recovery"
		consumerID = "p84-c9-consumer"
		claimID    = "p84-c9-claim"
	)
	root := t.TempDir()
	store, err := guard.NewLeaseStore(root)
	if err != nil {
		t.Fatal(err)
	}
	// The record path of a resource is "<root>/<resourceID>.lease.json", so the one
	// contended lock name is that path plus ".lock".
	lockPath := filepath.Join(root, resourceID+".lease.json.lock")
	markerPath := filepath.Join(root, "dead-holder.marker")
	holder := p84StartLockChild(t, root, resourceID, markerPath, "die")
	p84WaitForPath(t, markerPath, 30*time.Second)
	if err := holder.waitExit(t); err != nil {
		t.Fatalf("the dying holder must exit cleanly, got %v (%s)", err, strings.TrimSpace(holder.output.String()))
	}
	// The dead holder never released and never cleaned up, so its lock file is still
	// there. Only the kernel lock is gone, which is exactly the claim this case pins.
	if _, statErr := os.Stat(lockPath); statErr != nil {
		t.Fatalf("the dead holder's lock file must still be on disk, got %v", statErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	started := time.Now()
	lease, err := store.Claim(ctx, resourceID, consumerID, claimID, publishLeaseProcess(p84T5h4RunID), publishLeaseTTL)
	elapsed := time.Since(started)
	if err != nil {
		t.Fatalf("the dead holder's lock must be free for the next claimant, got %v after %s", err, elapsed)
	}
	if elapsed > time.Second {
		t.Fatalf("the dead holder's lock must be free immediately, took %s", elapsed)
	}
	if lease.Token.ClaimID != claimID || lease.Token.Generation != 1 {
		t.Fatalf("the freed lock must admit the ordinary claim, got %s/%d", lease.Token.ClaimID, lease.Token.Generation)
	}

	// Second leg, the ordinary writer shape: a holder that CLAIMS and then dies leaves a
	// committed record behind. The next claim must then be refused by that RECORD
	// (ErrLeaseHeld), not by a lock wait - and the refusal must arrive promptly rather
	// than at the context deadline, which is what tells the two failure modes apart.
	claimedResourceID := "ws-p84-c9-lock-claimed"
	claimedLockPath := filepath.Join(root, claimedResourceID+".lease.json.lock")
	claimer := p84StartLockChild(t, root, claimedResourceID, filepath.Join(root, "claim-holder.marker"), "claim")
	if err := claimer.waitExit(t); err != nil {
		t.Fatalf("the claiming holder must exit cleanly, got %v (%s)", err, strings.TrimSpace(claimer.output.String()))
	}
	claimCtx, claimCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer claimCancel()
	claimStarted := time.Now()
	_, err = store.Claim(claimCtx, claimedResourceID, consumerID, claimID, publishLeaseProcess(p84T5h4RunID), publishLeaseTTL)
	claimElapsed := time.Since(claimStarted)
	if !errors.Is(err, guard.ErrLeaseHeld) {
		t.Fatalf("the dead holder's committed record must refuse the next claim, got %v after %s", err, claimElapsed)
	}
	if claimElapsed > time.Second {
		t.Fatalf("the record refusal must not be delayed by the dead holder's lock, took %s", claimElapsed)
	}
	if _, statErr := os.Stat(claimedLockPath); statErr != nil {
		t.Fatalf("the lock file must outlive the process that created it, got %v", statErr)
	}
}

// TestP84C10AssertPublishLeaseCurrentRejectsASupersededProof pins F4: the check that
// stands between a proposal and the write of its task's store-local proof accepts
// exactly the store's current record.
//
// Mutation check: deleting the token/hash equality in assertPublishLeaseCurrent makes
// the superseded proof pass and turns the refusal assertion red.
func TestP84C10AssertPublishLeaseCurrentRejectsASupersededProof(t *testing.T) {
	state := p84T5h4Build(t, 1)
	ctx := context.Background()

	err := assertPublishLeaseCurrent(ctx, state.store, state.stale)
	if !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a proof the store has already left must be refused before publishing, got %v", err)
	}
	if !strings.Contains(err.Error(), "not the store's current") {
		t.Fatalf("the refusal must be the record equality, got %v", err)
	}
	if err := assertPublishLeaseCurrent(ctx, state.store, state.current); err != nil {
		t.Fatalf("the store's current record must pass the pre-publish check, got %v", err)
	}
	// A proof of a resource the store does not carry is refused through the same error.
	absent := publishLeaseProof{ResourceID: "ws-p84-c10-absent", Token: state.current.Token, LeaseHash: state.current.LeaseHash}
	if err := assertPublishLeaseCurrent(ctx, state.store, absent); !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a proof of an absent resource must be refused, got %v", err)
	}
}

// Round 4 of the tester adds the counterexamples of the final review's N1/N2, N4,
// N5 and N6 findings. Round 5 rewrites C9 and C11 for the OS-level lock:
//
// C11 (N1): a LIVE holder in another process keeps the resource lock; the claimant
// waits, fails on its context, and no lock file is removed.
//
// C13 (N4): a batch whose own refresh is already current still reconciles the run's
// other tasks instead of returning early.
//
// C14 (N5): the batch decides under LeaseStore.RunLocked, so a batch that runs while
// another holder owns the resource lock fails on its context instead of projecting a
// record it never locked.
//
// C15 (N6): a record of this run's resource issued for ANOTHER consumer is not this
// run's publish lease; a stale proof of it is refused and no transfer happens.

// TestP84C11LiveHolderKeepsTheLockUntilItExits pins the mutual exclusion the OS lock
// exists for, across processes: a holder in ANOTHER process keeps the resource until it
// exits. The competitor waits (5ms retries) and fails on its 500ms context instead of
// stealing the lock or proceeding unlocked, no lock file is removed, and once the live
// holder exits the very same claim succeeds.
//
// The holder is a real second process that blocks inside the locked section until the
// parent writes a release marker, so "held" is an observation of the child - the marker
// it writes after acquiring - and not a sleep.
//
// Mutation check (primary): deleting the OS lock acquisition in createLeaseLock - the
// LockFileEx call on Windows, the Flock call on POSIX - leaves no exclusion, so the
// parent's FIRST claim succeeds and the refusal assertion turns red.
func TestP84C11LiveHolderKeepsTheLockUntilItExits(t *testing.T) {
	const (
		resourceID = "ws-p84-c11-live-lock"
		consumerID = "p84-c11-consumer"
		claimID    = "p84-c11-claim"
	)
	root := t.TempDir()
	store, err := guard.NewLeaseStore(root)
	if err != nil {
		t.Fatal(err)
	}
	lockPath := filepath.Join(root, resourceID+".lease.json.lock")
	markerPath := filepath.Join(root, "live-holder.marker")
	holder := p84StartLockChild(t, root, resourceID, markerPath, "wait")
	p84WaitForPath(t, markerPath, 30*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	started := time.Now()
	_, err = store.Claim(ctx, resourceID, consumerID, claimID, publishLeaseProcess(p84T5h4RunID), publishLeaseTTL)
	elapsed := time.Since(started)
	if err == nil {
		t.Fatal("a live holder in another process must keep the resource lock")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("the claim must wait for the live holder and fail on its context, got %v", err)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("the claim must be bounded by its 500ms context, took %s", elapsed)
	}
	if _, statErr := os.Stat(lockPath); statErr != nil {
		t.Fatalf("the live holder's lock file must survive the failed claim, got %v", statErr)
	}

	// Release the live holder: its exit is what frees the kernel lock, exactly as a
	// crashed holder's exit does in C9.
	if err := os.WriteFile(markerPath+".release", []byte("release\n"), 0o600); err != nil {
		t.Fatalf("release the live holder: %v", err)
	}
	if err := holder.waitExit(t); err != nil {
		t.Fatalf("the released holder must exit cleanly, got %v (%s)", err, strings.TrimSpace(holder.output.String()))
	}
	retryCtx, retryCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer retryCancel()
	lease, err := store.Claim(retryCtx, resourceID, consumerID, claimID, publishLeaseProcess(p84T5h4RunID), publishLeaseTTL)
	if err != nil {
		t.Fatalf("the claim must succeed once the live holder exited, got %v", err)
	}
	if lease.Token.ClaimID != claimID || lease.Token.Generation != 1 {
		t.Fatalf("the freed lock must admit the waiting claim, got %s/%d", lease.Token.ClaimID, lease.Token.Generation)
	}
}

// TestP84T5h4C13CurrentRefreshStillReconcilesSiblings pins N4: a batch whose own
// refresh is already the current record must still reconcile the run's other tasks.
// Several tasks share one run-scoped record, so a task can fall behind while this
// task's own proof is already current - and this call is the only place it is
// caught up. The old shape returned as soon as refreshed == previous, which left
// that sibling behind forever.
//
// Mutation check: restoring `if refreshed == previous { return nil }` at the top of
// persistRenewedPublishLease leaves the sibling on its old generation and turns the
// byte assertion red (the main task's proof is already current either way, so only
// the sibling distinguishes the two shapes).
func TestP84T5h4C13CurrentRefreshStillReconcilesSiblings(t *testing.T) {
	// Two renewals: the stale proof is two generations behind the current record, so
	// the caught-up bytes cannot be confused with "nothing to do".
	state := p84T5h4Build(t, 2)
	ctx := context.Background()
	mainPath := state.proofPath("task-main")
	siblingPath := state.proofPath("task-sibling")
	if err := writeSupervisedJSON(mainPath, state.current); err != nil {
		t.Fatal(err)
	}
	if err := writeSupervisedJSON(siblingPath, state.stale); err != nil {
		t.Fatal(err)
	}

	// The caller's refresh IS the store's current record: previous == refreshed.
	if err := state.runtime().persistRenewedPublishLease(ctx, state.store, state.current, state.current); err != nil {
		t.Fatalf("an already-current refresh must still reconcile the run's tasks, got %v", err)
	}
	want := p84T5h4ProofBytes(t, state.current)
	if got := p84T5h4ReadFile(t, siblingPath); !bytes.Equal(got, want) {
		t.Fatalf("a sibling that fell behind must be caught up even though this task was current, got %q want %q", got, want)
	}
	if got := p84T5h4ReadFile(t, mainPath); !bytes.Equal(got, want) {
		t.Fatalf("the current proof must stay current, got %q want %q", got, want)
	}
}

// TestP84T5h4C14BatchWaitsForTheResourceLock pins N5: the batch decides on a record
// read UNDER the store's resource lock (LeaseStore.RunLocked), so it is ordered
// against the store's own mutations. While another holder owns that lock the batch
// cannot proceed, and its context bounds the wait: a batch that returned success
// there would be projecting a record it never locked.
//
// The holder is an exported RunLocked of this same process with a function that
// blocks until the test releases it, which is exactly the shape a concurrent
// renewal has for the report: same resource, same lock, different goroutine.
//
// Mutation check: reading the record with store.Current and running the batch body
// outside RunLocked (the pre-N5 shape) makes the call succeed while the lock is
// held and turns the refusal assertion red; the control afterwards keeps the case
// from passing merely because the batch is broken altogether.
func TestP84T5h4C14BatchWaitsForTheResourceLock(t *testing.T) {
	state := p84T5h4Build(t, 1)
	path := state.proofPath("task-locked")
	if err := writeSupervisedJSON(path, state.stale); err != nil {
		t.Fatal(err)
	}

	held := make(chan struct{})
	release := make(chan struct{})
	holderDone := make(chan error, 1)
	go func() {
		holderDone <- state.store.RunLocked(context.Background(), state.current.ResourceID, func(guard.Lease) error {
			close(held)
			<-release
			return nil
		})
	}()
	<-held
	// finish releases the holder exactly once, either from the control below or from
	// the deferred cleanup of an early failure, so the goroutine never outlives the
	// test and the channel is never closed twice.
	finished := false
	finish := func() {
		if finished {
			return
		}
		finished = true
		close(release)
		if err := <-holderDone; err != nil {
			t.Errorf("the lock holder must finish cleanly, got %v", err)
		}
	}
	defer finish()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	started := time.Now()
	err := state.runtime().persistRenewedPublishLease(ctx, state.store, state.stale, state.current)
	elapsed := time.Since(started)
	if err == nil {
		t.Fatal("a batch must not complete while another holder owns the resource lock")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("the batch must wait for the resource lock and fail on its context, got %v", err)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("the batch must be bounded by its 500ms context, took %s", elapsed)
	}
	if got := p84T5h4ReadFile(t, path); !bytes.Equal(got, p84T5h4ProofBytes(t, state.stale)) {
		t.Fatalf("a batch that never got the lock must rewrite nothing, got %q", got)
	}

	// Control: once the lock is free the same batch succeeds, so the refusal above is
	// the lock and not a broken call.
	finish()
	if err := state.runtime().persistRenewedPublishLease(context.Background(), state.store, state.stale, state.current); err != nil {
		t.Fatalf("the same batch must succeed once the lock is free, got %v", err)
	}
	if got, want := p84T5h4ReadFile(t, path), p84T5h4ProofBytes(t, state.current); !bytes.Equal(got, want) {
		t.Fatalf("the released batch must rewrite the proof, got %q want %q", got, want)
	}
}

// TestP84C15ForeignConsumerRecordIsNotThisRunsLease pins N6: a record of this run's
// resource that was issued for ANOTHER consumer is not this run's publish lease. An
// old proof of that record is genuine retained history, so the history walk alone
// would accept it and the transfer behind the stopped-holder proof would then hand
// the record to this process - a takeover of a lease this run never held.
//
// The record is built for real (claim + one renewal as a foreign consumer, with a
// holder the OS proves stopped), so the refusal cannot come from a shape, claim,
// generation or history check: only the consumer binding stands in the way.
//
// Mutation check: deleting the publishLeaseBelongsToRun guard from renewPublishLease
// AND from the Reclaim callback it delegates to (the guard is the pair of those two
// call sites) lets the transfer happen, the call returns a proof instead of
// ErrPublishLeaseUnproven, and the record assertions turn red.
func TestP84C15ForeignConsumerRecordIsNotThisRunsLease(t *testing.T) {
	const foreignConsumerID = "p84-c15-foreign-consumer"
	runRoot := t.TempDir()
	identity := p84C8DeadIdentity(t, runRoot)
	ctx := context.Background()
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	resourceID := publishLeaseResourceID(p84T5h4RunID)
	first, err := store.Claim(ctx, resourceID, foreignConsumerID, publishLeaseClaimID(p84T5h4RunID), identity, publishLeaseTTL)
	if err != nil {
		t.Fatalf("claim the resource for a foreign consumer: %v", err)
	}
	second, err := store.Renew(ctx, resourceID, first.Token, identity, publishLeaseTTL)
	if err != nil {
		t.Fatalf("advance the foreign consumer's record: %v", err)
	}
	stale := publishLeaseProof{ResourceID: resourceID, Token: first.Token, LeaseHash: first.LeaseHash}

	refreshed, err := renewPublishLease(ctx, runRoot, p84T5h4RunID, stale, time.Now().UTC())
	if !errors.Is(err, ErrPublishLeaseUnproven) {
		t.Fatalf("a record issued for another consumer must not be taken over, got %v (proof %+v)", err, refreshed)
	}
	if !strings.Contains(err.Error(), "not this run's publish lease") {
		t.Fatalf("the refusal must be the consumer binding, got %v", err)
	}
	current, err := store.Current(ctx, resourceID)
	if err != nil {
		t.Fatalf("read the record back: %v", err)
	}
	if current.ConsumerID != foreignConsumerID {
		t.Fatalf("the foreign consumer's record must keep its consumer, got %q", current.ConsumerID)
	}
	if current.Token != second.Token || current.LeaseHash != second.LeaseHash {
		t.Fatalf("the foreign consumer's record must not be transferred, got %s/%d hash %s want %s/%d hash %s", current.Token.ClaimID, current.Token.Generation, current.LeaseHash, second.Token.ClaimID, second.Token.Generation, second.LeaseHash)
	}
}

// TestP84C16PublishProofStaysUntouchedWhileTheResourceLockIsHeld pins the N5/F4
// residual of fix round 9: this task's OWN publish proof is written inside the
// store's resource critical section, exactly like the batch over its siblings.
//
// The write is exercised through its production entry, writePublishLeaseProofLocked
// - the call Publish makes at runtime_supervised.go:1982 and the claim/adopt branch
// makes at :1173. Driving the full Publish entry cannot isolate this write: its
// earlier renewal (renewPublishLease -> LeaseStore.Renew) takes the same resource
// lock, so a held lock fails the promotion before the proof write is reached and
// the entry would go red for the wrong reason. Holding RunLocked around the write
// itself separates "the proof write is serialised" from "something earlier in
// Publish happened to block"; the scope reduction is registered in the report.
//
// Mutation check: writing the proof outside LeaseStore.RunLocked (the pre-round-9
// shape: store.Current, compare, write) lets the bytes land while another holder
// owns the lock, turning the refusal assertion and the byte assertion red; the
// control leg keeps the case from passing merely because the write is broken.
func TestP84C16PublishProofStaysUntouchedWhileTheResourceLockIsHeld(t *testing.T) {
	state := p84T5h4Build(t, 1)
	path := state.proofPath("task-publish")
	if err := writeSupervisedJSON(path, state.stale); err != nil {
		t.Fatal(err)
	}
	staleBytes := p84T5h4ProofBytes(t, state.stale)
	if got := p84T5h4ReadFile(t, path); !bytes.Equal(got, staleBytes) {
		t.Fatalf("test precondition failed: the proof on disk must be the stale image, got %q", got)
	}

	held := make(chan struct{})
	release := make(chan struct{})
	holderDone := make(chan error, 1)
	go func() {
		holderDone <- state.store.RunLocked(context.Background(), state.current.ResourceID, func(guard.Lease) error {
			close(held)
			<-release
			return nil
		})
	}()
	<-held
	// finish releases the holder exactly once, from the control or from the deferred
	// cleanup of an early failure, so the goroutine never outlives the test.
	finished := false
	finish := func() {
		if finished {
			return
		}
		finished = true
		close(release)
		if err := <-holderDone; err != nil {
			t.Errorf("the lock holder must finish cleanly, got %v", err)
		}
	}
	defer finish()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	started := time.Now()
	err := writePublishLeaseProofLocked(ctx, state.store, path, state.current)
	elapsed := time.Since(started)
	if err == nil {
		t.Fatal("this task's proof must not be written while another holder owns the resource lock")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("the proof write must wait for the resource lock and fail on its context, got %v", err)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("the proof write must be bounded by its 500ms context, took %s", elapsed)
	}
	if got := p84T5h4ReadFile(t, path); !bytes.Equal(got, staleBytes) {
		t.Fatalf("a proof write that never got the lock must change no byte, got %q want %q", got, staleBytes)
	}

	// Control: with the lock free the same call lands the current record, so the
	// refusal above is the lock and not a broken write.
	finish()
	if err := writePublishLeaseProofLocked(context.Background(), state.store, path, state.current); err != nil {
		t.Fatalf("the same proof must be writable once the lock is free, got %v", err)
	}
	if got, want := p84T5h4ReadFile(t, path), p84T5h4ProofBytes(t, state.current); !bytes.Equal(got, want) {
		t.Fatalf("the released write must project the current record, got %q want %q", got, want)
	}
}

// TestP84C18PublishEntryHoldsTheResourceLockWithoutRenewing closes the N11
// acceptance gap of fix round 10: the proof write Publish performs for its own
// task has to be serialised on the resource lock on the path that does NOT renew
// the lease either.
//
// C16 pinned the same write through writePublishLeaseProofLocked, but its comment
// had to reduce the scope: driving the full Publish entry could not isolate the
// write, because the renewal that runs earlier takes the same resource lock and a
// held lock failed the promotion before the write was reached. This test closes the
// gap by removing the renewal from the picture instead of the entry.
//
// The stored proof is left exactly one generation behind the record while the
// record still names THIS process and still expires far beyond
// publishLeaseRenewBefore, so:
//
//   - renewPublishLease takes its same-process catch-up branch (the proof is real
//     history of the same claim: N-1, reachable through PreviousHash and the
//     retained body), which makes "refreshed" the current record without calling
//     Renew at all - the no-renewal case;
//   - refreshed != proof then sends Publish into this task's own proof write, which
//     is the branch under test.
//
// With another holder on the resource lock, Publish must not complete and the proof
// bytes must not move; once the lock is released the same call completes and the
// proof carries the current record.
//
// Mutation check: restoring the pre-round-9 write in Publish (assertPublishLeaseCurrent
// followed by writeSupervisedJSON, both outside LeaseStore.RunLocked) lets the proof
// land while the lock is held. The byte assertion catches the write itself and stays
// red even though the sibling batch that follows still times out, so the case cannot
// pass by the batch's failure alone.
func TestP84C18PublishEntryHoldsTheResourceLockWithoutRenewing(t *testing.T) {
	ctx := context.Background()
	harness := newSupervisedHarness(t, []string{adapters.VisibleStubAccepted}, false, 300000)
	harness.now = p83FixedNow()
	harness.driveToWindow()
	publishRuntime, state := p83F3SupervisedRuntime(t, harness)
	request := p83F3SupervisedReview(t, harness, publishRuntime, state)

	store, err := newPublishLeaseStore(publishRuntime.runRoot)
	if err != nil {
		t.Fatal(err)
	}
	resourceID := publishLeaseResourceID(supervisedTestRunID)
	claimed, err := store.Current(ctx, resourceID)
	if err != nil {
		t.Fatalf("read the publish lease the window claimed: %v", err)
	}
	stale := publishLeaseProof{ResourceID: claimed.ResourceID, Token: claimed.Token, LeaseHash: claimed.LeaseHash}
	path := publishRuntime.publishLeasePath(state)
	staleBytes := p84T5h4ProofBytes(t, stale)
	if got := p84T5h4ReadFile(t, path); !bytes.Equal(got, staleBytes) {
		t.Fatalf("test precondition failed: the window's proof must be the claimed generation, got %q want %q", got, staleBytes)
	}

	// Advance exactly one generation and deliberately leave the proof alone: it is now
	// the N-1 predecessor the same-process catch-up re-syncs.
	advanced, err := store.Renew(ctx, claimed.ResourceID, claimed.Token, publishLeaseProcess(supervisedTestRunID), publishLeaseTTL)
	if err != nil {
		t.Fatalf("commit one generation without rewriting the proof: %v", err)
	}
	current := publishLeaseProof{ResourceID: advanced.ResourceID, Token: advanced.Token, LeaseHash: advanced.LeaseHash}
	if current == stale {
		t.Fatal("test precondition failed: the record must have advanced by one generation")
	}
	if got := p84T5h4ReadFile(t, path); !bytes.Equal(got, staleBytes) {
		t.Fatalf("test precondition failed: the proof must still be the stale image, got %q", got)
	}
	// The renewal threshold must not be reached, or Publish would renew the record
	// instead of taking the catch-up path this test is built on.
	expiresAt, err := time.Parse(publishLeaseTimestamp, advanced.LeaseExpiresAt)
	if err != nil {
		t.Fatalf("read the record's expiry: %v", err)
	}
	if remaining := expiresAt.Sub(harness.clock()); remaining <= publishLeaseRenewBefore {
		t.Fatalf("test precondition failed: the record must outlive the renewal threshold, only %s remain", remaining)
	}

	held := make(chan struct{})
	release := make(chan struct{})
	holderDone := make(chan error, 1)
	go func() {
		holderDone <- store.RunLocked(context.Background(), resourceID, func(guard.Lease) error {
			close(held)
			<-release
			return nil
		})
	}()
	<-held
	// finish releases the holder exactly once, from the control or from the deferred
	// cleanup of an early failure, so the goroutine never outlives the test.
	finished := false
	finish := func() {
		if finished {
			return
		}
		finished = true
		close(release)
		if err := <-holderDone; err != nil {
			t.Errorf("the lock holder must finish cleanly, got %v", err)
		}
	}
	defer finish()

	blockedCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	started := time.Now()
	decision, publishErr := publishRuntime.Publish(blockedCtx, request)
	elapsed := time.Since(started)
	if publishErr == nil {
		t.Fatalf("Publish must not complete while another holder owns the resource lock, got %+v", decision)
	}
	// The context failure has to come from this task's proof write, i.e. from the
	// resource lock, and not from some earlier step that broke for its own reason.
	if !errors.Is(publishErr, context.DeadlineExceeded) {
		t.Fatalf("Publish must reach this task's proof write and fail on its context there, got %v", publishErr)
	}
	if elapsed > 5*time.Second {
		t.Fatalf("Publish must be bounded by its 500ms context, took %s", elapsed)
	}
	if decision.AcceptedSnapshotHash != "" || len(decision.LeaseEvidence) != 0 || len(decision.ErrorEvidence) != 0 {
		t.Fatalf("a publish that never got the lock must return no decision at all, got %+v", decision)
	}
	if got := p84T5h4ReadFile(t, path); !bytes.Equal(got, staleBytes) {
		t.Fatalf("a publish that never got the lock must change no proof byte, got %q want %q", got, staleBytes)
	}

	// Control: the same call completes once the lock is free, and the proof it lands is
	// the current record - the refusal above is the lock, not a broken entry.
	finish()
	decision, err = publishRuntime.Publish(context.Background(), request)
	if err != nil {
		t.Fatalf("the same publish must complete once the lock is free, got %v", err)
	}
	if decision.Outcome != "completed" || len(decision.LeaseEvidence) != 1 {
		t.Fatalf("the released publish must complete with its lease evidence, got %+v", decision)
	}
	if got, want := p84T5h4ReadFile(t, path), p84T5h4ProofBytes(t, current); !bytes.Equal(got, want) {
		t.Fatalf("the released publish must project the current record onto this task's proof, got %q want %q", got, want)
	}
}

// TestP84C17LegacyLockDirectoryIsReportedNotWaitedOn pins N8: a lock path that is a
// DIRECTORY is the upgrade boundary left by the mkdir-style lock of an older run.
// It can neither be acquired (this substrate locks a byte range of a file) nor
// deleted (an old run that still believes it holds the resource would then get a
// second writer), and it will never turn into a free lock file, so waiting on it is
// pointless. acquireLeaseLock therefore reports ErrLeaseLockPathIsDirectory
// immediately instead of entering the contention loop.
//
// The legacy shape is made by hand: mkdir at "<record>.lock", the same name guard
// derives from the record path (record = "<resource>.lease.json", lock = record +
// ".lock"). A short context bounds any wait, so an implementation that treats the
// directory as contention reports context.DeadlineExceeded instead of the sentinel.
//
// Mutation check: classifying the failed CreateFile as contention (returning
// errLeaseLockHeld for the directory) makes the claim wait for the 100ms context,
// and the sentinel assertion goes red.
func TestP84C17LegacyLockDirectoryIsReportedNotWaitedOn(t *testing.T) {
	const runID = "run-p84-c17-legacy-lock"
	runRoot := t.TempDir()
	store, err := newPublishLeaseStore(runRoot)
	if err != nil {
		t.Fatal(err)
	}
	resource := publishLeaseResourceID(runID)
	recordPath := filepath.Join(publishLeaseRoot(runRoot), resource+".lease.json")
	lockPath := recordPath + ".lock"
	if err := os.Mkdir(lockPath, 0o755); err != nil {
		t.Fatalf("create the legacy lock directory: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	started := time.Now()
	_, err = store.Claim(ctx, resource, publishLeaseConsumerID, publishLeaseClaimID(runID), publishLeaseProcess(runID), publishLeaseTTL)
	elapsed := time.Since(started)
	if err == nil {
		t.Fatal("a claim must not succeed while a legacy lock directory occupies the lock path")
	}
	if errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("the legacy lock directory must not be waited on until the context expires, took %s: %v", elapsed, err)
	}
	if !errors.Is(err, guard.ErrLeaseLockPathIsDirectory) {
		t.Fatalf("the legacy lock directory must be reported as its own sentinel, got %v after %s", err, elapsed)
	}
	// The substrate neither deletes nor renames the directory: migration is the
	// operator's step, after the old run is proven stopped.
	if info, statErr := os.Stat(lockPath); statErr != nil || !info.IsDir() {
		t.Fatalf("the legacy lock directory must be left in place for the operator, stat err=%v", statErr)
	}
	if _, statErr := os.Stat(recordPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the refused claim must not have written a record, stat err=%v", statErr)
	}

	// Control: once the operator has removed the directory the same claim succeeds,
	// so the refusal above is the directory and not a broken store.
	if err := os.Remove(lockPath); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Claim(context.Background(), resource, publishLeaseConsumerID, publishLeaseClaimID(runID), publishLeaseProcess(runID), publishLeaseTTL); err != nil {
		t.Fatalf("the same claim must succeed once the legacy directory is gone, got %v", err)
	}
}
