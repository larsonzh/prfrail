package adapters

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
)

type recordingLaunchStub struct {
	mu        sync.Mutex
	calls     int
	requests  []chain.AgentRunnerLaunchRequest
	err       error
	launchID  string
	processID string
	startedAt time.Time
}

func (stub *recordingLaunchStub) StartAgentRunnerProcess(_ context.Context, launch chain.AgentRunnerLaunchRequest) (chain.AgentRunnerLaunchResult, error) {
	stub.mu.Lock()
	defer stub.mu.Unlock()
	stub.calls++
	stub.requests = append(stub.requests, launch)
	if stub.err != nil {
		return chain.AgentRunnerLaunchResult{}, stub.err
	}
	processID := stub.processID
	if processID == "" {
		processID = "process-stub-one"
	}
	startedAt := stub.startedAt
	if startedAt.IsZero() {
		startedAt = time.Date(2026, 9, 11, 8, 6, 0, 0, time.UTC)
	}
	launchID := "launch-" + launch.RequestID
	if stub.launchID != "" {
		launchID = stub.launchID
	}
	return chain.AgentRunnerLaunchResult{LaunchID: launchID, ProcessID: processID, StartedAt: startedAt}, nil
}

func (stub *recordingLaunchStub) callCount() int {
	stub.mu.Lock()
	defer stub.mu.Unlock()
	return stub.calls
}

func replayDispatcherFixture(t *testing.T) (*AgentRunnerReplayDispatcher, *AgentRunnerCompositeAdmission, *recordingLaunchStub, chain.StepRequest) {
	t.Helper()
	admission, request := compositeAdmissionFixture(t)
	store := replayStoreMustNew(t, t.TempDir())
	launcher := &recordingLaunchStub{}
	dispatcher := &AgentRunnerReplayDispatcher{
		Store:          store,
		Launcher:       launcher,
		LedgerSnapshot: func() (chain.AuthorizationLedger, error) { return admission.AuthorizationLedger, nil },
		CostLedger:     admission.CostLedger,
		RequestRecord:  admission.RequestRecord,
		Clock:          admission.Clock,
	}
	return dispatcher, &admission, launcher, request
}

func TestAgentRunnerReplayDispatcherFirstDispatchPublishesReceiptsAndLaunches(t *testing.T) {
	dispatcher, _, launcher, request := replayDispatcherFixture(t)
	result, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if err != nil {
		t.Fatalf("expected first dispatch to succeed, got %v", err)
	}
	if launcher.callCount() != 1 {
		t.Fatalf("expected exactly one launch, got %d", launcher.callCount())
	}
	receipt, found, err := dispatcher.Store.LaunchReceipt("request-one")
	if err != nil || !found {
		t.Fatalf("expected launch receipt: found=%v err=%v", found, err)
	}
	if receipt.Receipt.ReconfirmedAt != "2026-09-11T08:05:00Z" {
		t.Fatalf("reconfirmedAt = %s, want reconfirmation clock time", receipt.Receipt.ReconfirmedAt)
	}
	if receipt.Receipt.RequestHash != dispatcher.RequestRecord.RecordHash {
		t.Fatalf("receipt requestHash = %s, want %s", receipt.Receipt.RequestHash, dispatcher.RequestRecord.RecordHash)
	}
	identity, found, err := dispatcher.Store.LaunchIdentity("launch-request-one")
	if err != nil || !found {
		t.Fatalf("expected process identity: found=%v err=%v", found, err)
	}
	if identity.Identity.ProcessID != "process-stub-one" || identity.Identity.IntentRecordHash != receipt.RecordHash {
		t.Fatalf("identity not bound to intent: %+v", identity.Identity)
	}
	if len(result.Evidence) != 3 ||
		result.Evidence[0] != dispatcher.RequestRecord.RecordHash ||
		result.Evidence[1] != receipt.RecordHash ||
		result.Evidence[2] != identity.RecordHash {
		t.Fatalf("unexpected evidence chain: %v", result.Evidence)
	}
	state, err := dispatcher.Store.State("request-one")
	if err != nil || state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("expected dispatched-unknown state, got state=%v err=%v", state, err)
	}
}

func TestAgentRunnerReplayDispatcherBlocksWithoutClaiming(t *testing.T) {
	t.Run("request only", func(t *testing.T) {
		dispatcher, _, launcher, request := replayDispatcherFixture(t)
		if _, err := dispatcher.Store.RecordRequest(dispatcher.RequestRecord); err != nil {
			t.Fatal(err)
		}
		if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
			t.Fatalf("expected unknown block, got %v", err)
		}
		if launcher.callCount() != 0 {
			t.Fatalf("no launch may happen for a replay, got %d", launcher.callCount())
		}
	})
	t.Run("intent present", func(t *testing.T) {
		dispatcher, _, launcher, request := replayDispatcherFixture(t)
		if _, err := dispatcher.Store.RecordRequest(dispatcher.RequestRecord); err != nil {
			t.Fatal(err)
		}
		receipt := launchReceiptRecordFor(t, dispatcher.RequestRecord)
		if _, err := dispatcher.Store.RecordLaunchReceipt(receipt); err != nil {
			t.Fatal(err)
		}
		if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
			t.Fatalf("expected unknown block for existing intent, got %v", err)
		}
		if launcher.callCount() != 0 {
			t.Fatalf("retained intent must never be deleted or relaunched, got %d launches", launcher.callCount())
		}
	})
	t.Run("identity present", func(t *testing.T) {
		dispatcher, _, launcher, request := replayDispatcherFixture(t)
		if _, err := dispatcher.Store.RecordRequest(dispatcher.RequestRecord); err != nil {
			t.Fatal(err)
		}
		receipt := launchReceiptRecordFor(t, dispatcher.RequestRecord)
		if _, err := dispatcher.Store.RecordLaunchReceipt(receipt); err != nil {
			t.Fatal(err)
		}
		if _, err := dispatcher.Store.RecordLaunchIdentity(launchIdentityRecordFor(t, receipt)); err != nil {
			t.Fatal(err)
		}
		if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
			t.Fatalf("expected unknown block for existing identity, got %v", err)
		}
		if launcher.callCount() != 0 {
			t.Fatalf("identity-present state must never relaunch, got %d launches", launcher.callCount())
		}
	})
	t.Run("terminal receipt", func(t *testing.T) {
		dispatcher, _, launcher, request := replayDispatcherFixture(t)
		if _, err := dispatcher.Store.RecordRequest(dispatcher.RequestRecord); err != nil {
			t.Fatal(err)
		}
		completion := replayStoreCompletionRecord(t, dispatcher.RequestRecord, "completion-one")
		if _, err := dispatcher.Store.RecordCompletion(dispatcher.RequestRecord, completion); err != nil {
			t.Fatal(err)
		}
		if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchTerminalReceiptPresent) {
			t.Fatalf("expected terminal receipt block, got %v", err)
		}
		if launcher.callCount() != 0 {
			t.Fatalf("terminal state must never launch, got %d", launcher.callCount())
		}
	})
}

func TestAgentRunnerReplayDispatcherReconfirmationFailureLeavesUnknownBlock(t *testing.T) {
	dispatcher, admission, launcher, request := replayDispatcherFixture(t)
	grant := admission.AuthorizationLedger.Grants[0]
	revocation, err := chain.NewRevocationAuthorizationRecord(chain.AuthorizationRevocation{
		RecordID:          "revocation-one",
		Kind:              "revocation",
		IssuedAt:          "2026-09-11T08:00:00.000Z",
		IssuedBy:          evidence.Actor{Type: "operator", ID: "operator-one"},
		AuthorizationID:   grant.Grant.AuthorizationID,
		AuthorizationHash: grant.RecordHash,
		ReasonCode:        "operator-request",
		StopDisposition:   "not-required",
	}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	admission.AuthorizationLedger.Revocations = []chain.AuthorizationRecord{revocation}

	_, err = dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) || !strings.Contains(err.Error(), "revoked") {
		t.Fatalf("expected revoked reconfirmation failure, got %v", err)
	}
	if launcher.callCount() != 0 {
		t.Fatalf("revoked authorization must never launch, got %d launches", launcher.callCount())
	}
	if _, found, err := dispatcher.Store.LaunchReceipt("request-one"); err != nil || found {
		t.Fatalf("no launch intent may be published on failed reconfirmation: found=%v err=%v", found, err)
	}
	state, err := dispatcher.Store.State("request-one")
	if err != nil || state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("request record must be retained, got state=%v err=%v", state, err)
	}
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
		t.Fatalf("failed reconfirmation must not auto-retry, got %v", err)
	}
	if launcher.callCount() != 0 {
		t.Fatalf("no relaunch after failed reconfirmation, got %d launches", launcher.callCount())
	}
}

func TestAgentRunnerReplayDispatcherDurabilityGatePrecedesLaunch(t *testing.T) {
	dispatcher, _, launcher, request := replayDispatcherFixture(t)
	dispatcher.Store.publicationDurability = PublishDurabilityUnproven
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("expected durability gate, got %v", err)
	}
	if launcher.callCount() != 0 {
		t.Fatalf("unproven store must never launch, got %d", launcher.callCount())
	}
	state, err := dispatcher.Store.State("request-one")
	if err != nil || state != AgentRunnerReplayStateAbsent {
		t.Fatalf("unproven store must not publish request records, got state=%v err=%v", state, err)
	}
}

func TestAgentRunnerReplayDispatcherStartFailureRetainsIntentWithoutRelaunch(t *testing.T) {
	dispatcher, _, launcher, request := replayDispatcherFixture(t)
	launcher.err = errors.New("spawn refused")
	_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchFailed) || !strings.Contains(err.Error(), "spawn refused") {
		t.Fatalf("expected start failure, got %v", err)
	}
	if launcher.callCount() != 1 {
		t.Fatalf("expected one start attempt, got %d", launcher.callCount())
	}
	if _, found, err := dispatcher.Store.LaunchReceipt("request-one"); err != nil || !found {
		t.Fatalf("launch intent must be retained on unspawned failure: found=%v err=%v", found, err)
	}
	if _, found, err := dispatcher.Store.LaunchIdentity("launch-request-one"); err != nil || found {
		t.Fatalf("no identity may exist without a spawn: found=%v err=%v", found, err)
	}
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
		t.Fatalf("retained intent must block redispatch, got %v", err)
	}
	if launcher.callCount() != 1 {
		t.Fatalf("retained intent must never be relaunched, got %d launches", launcher.callCount())
	}
}

func TestAgentRunnerReplayDispatcherLaunchIDMismatchIsUnproven(t *testing.T) {
	dispatcher, _, launcher, request := replayDispatcherFixture(t)
	launcher.launchID = "launch-foreign-one"
	_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) || !strings.Contains(err.Error(), "launch-foreign-one") {
		t.Fatalf("expected unproven identity for mismatched launchId, got %v", err)
	}
	if launcher.callCount() != 1 {
		t.Fatalf("expected exactly one start attempt, got %d", launcher.callCount())
	}
	if _, found, err := dispatcher.Store.LaunchIdentity("launch-request-one"); err != nil || found {
		t.Fatalf("no identity may be published for a mismatched launchId: found=%v err=%v", found, err)
	}
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
		t.Fatalf("mismatched launchId must block redispatch, got %v", err)
	}
	if launcher.callCount() != 1 {
		t.Fatalf("mismatched launchId must never relaunch, got %d launches", launcher.callCount())
	}
}

func TestAgentRunnerReplayDispatcherConcurrentDispatchLaunchesOnce(t *testing.T) {
	dispatcher, _, launcher, request := replayDispatcherFixture(t)
	const workers = 8
	results := make([]error, workers)
	var wg sync.WaitGroup
	for index := 0; index < workers; index++ {
		wg.Add(1)
		go func(slot int) {
			defer wg.Done()
			_, results[slot] = dispatcher.DispatchAgentRunner(context.Background(), request)
		}(index)
	}
	wg.Wait()
	successes := 0
	for _, err := range results {
		if err == nil {
			successes++
			continue
		}
		if !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
			t.Fatalf("unexpected concurrent dispatch error: %v", err)
		}
	}
	if successes != 1 {
		t.Fatalf("expected exactly one winning dispatch, got %d", successes)
	}
	if launcher.callCount() != 1 {
		t.Fatalf("concurrent dispatchers must launch once, got %d", launcher.callCount())
	}
	if _, found, err := dispatcher.Store.LaunchIdentity("launch-request-one"); err != nil || !found {
		t.Fatalf("expected single identity publication: found=%v err=%v", found, err)
	}
}

func TestAgentRunnerReplayDispatcherRejectsBindingMismatch(t *testing.T) {
	dispatcher, _, launcher, request := replayDispatcherFixture(t)
	tampered := request
	tampered.ParentHash = queueHashOne
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), tampered); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected binding mismatch rejection, got %v", err)
	}
	noFacts := request
	noFacts.AgentRunnerFacts = nil
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), noFacts); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected facts requirement rejection, got %v", err)
	}
	foreignFacts := request
	facts := *request.AgentRunnerFacts
	facts.ContextHash = queueHashOne
	foreignFacts.AgentRunnerFacts = &facts
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), foreignFacts); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("expected facts binding rejection, got %v", err)
	}
	if launcher.callCount() != 0 {
		t.Fatalf("binding mismatch must never launch, got %d", launcher.callCount())
	}
	state, err := dispatcher.Store.State("request-one")
	if err != nil || state != AgentRunnerReplayStateAbsent {
		t.Fatalf("rejected dispatch must not publish request records, got state=%v err=%v", state, err)
	}
}

func TestAgentRunnerReplayDispatcherIdentityConflictSurfacesAsUnproven(t *testing.T) {
	dispatcher, _, launcher, request := replayDispatcherFixture(t)
	receipt := launchReceiptRecordFor(t, dispatcher.RequestRecord)
	foreign, err := NewAgentRunnerLaunchIdentityRecord(AgentRunnerLaunchIdentity{
		LaunchID:         receipt.Receipt.LaunchID,
		RequestID:        receipt.Receipt.RequestID,
		IntentRecordHash: receipt.RecordHash,
		ProcessID:        "process-other-one",
		StartedAt:        "2026-09-11T08:05:02Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	path, err := dispatcher.Store.launchIdentityPath(receipt.Receipt.LaunchID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(path, foreign); err != nil {
		t.Fatal(err)
	}
	_, err = dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) || !errors.Is(err, ErrAgentRunnerLaunchIdentityConflict) {
		t.Fatalf("expected identity conflict to surface as unproven, got %v", err)
	}
	if launcher.callCount() != 1 {
		t.Fatalf("identity conflict must not trigger kill-and-retry, got %d launches", launcher.callCount())
	}
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
		t.Fatalf("unproven identity must block redispatch, got %v", err)
	}
	if launcher.callCount() != 1 {
		t.Fatalf("unproven identity must never relaunch, got %d launches", launcher.callCount())
	}
}

func TestAgentRunnerReplayDispatcherLedgerSnapshotFailureBlocksLaunch(t *testing.T) {
	dispatcher, _, launcher, request := replayDispatcherFixture(t)
	dispatcher.LedgerSnapshot = func() (chain.AuthorizationLedger, error) {
		return chain.AuthorizationLedger{}, errors.New("ledger offline")
	}
	_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchReconfirmation) || !strings.Contains(err.Error(), "ledger offline") {
		t.Fatalf("expected reconfirmation failure for offline ledger, got %v", err)
	}
	if launcher.callCount() != 0 {
		t.Fatalf("offline ledger must never launch, got %d", launcher.callCount())
	}
	state, err := dispatcher.Store.State("request-one")
	if err != nil || state != AgentRunnerReplayStateDispatchedUnknown {
		t.Fatalf("request record must be retained without a launch, got state=%v err=%v", state, err)
	}
}

func TestAgentRunnerReplayDispatcherFailsClosedOnMissingFields(t *testing.T) {
	dispatcher, _, _, request := replayDispatcherFixture(t)
	var nilDispatcher *AgentRunnerReplayDispatcher
	if _, err := nilDispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("nil dispatcher must fail closed, got %v", err)
	}
	cases := map[string]func(*AgentRunnerReplayDispatcher){
		"store":    func(d *AgentRunnerReplayDispatcher) { d.Store = nil },
		"launcher": func(d *AgentRunnerReplayDispatcher) { d.Launcher = nil },
		"ledger":   func(d *AgentRunnerReplayDispatcher) { d.LedgerSnapshot = nil },
		"clock":    func(d *AgentRunnerReplayDispatcher) { d.Clock = nil },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			copyDispatcher := *dispatcher
			mutate(&copyDispatcher)
			if _, err := copyDispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
				t.Fatalf("missing %s must fail closed, got %v", name, err)
			}
		})
	}
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := dispatcher.DispatchAgentRunner(cancelled, request); !errors.Is(err, ErrInvalidAgentRunnerRequest) {
		t.Fatalf("cancelled context must fail closed, got %v", err)
	}
}

// reclaimingLaunchStub is a launcher that also exposes the private stop
// capability the dispatcher reclaims an unproven spawn through. It embeds
// recordingLaunchStub so the start behaviour of both stubs cannot drift apart.
// A stub can only report what it was told to: it records the stop calls it
// received and never claims anything about an operating-system process.
type reclaimingLaunchStub struct {
	*recordingLaunchStub
	stopMu      sync.Mutex
	stopCalls   int
	stopIDs     []string
	stopErr     error
	stopOutcome string
	// cancelOnStart, when set, cancels the dispatch context from inside the
	// start call, before any result is returned. It exists only so a test can
	// prove the reclaim shield survives a caller that has already given up.
	cancelOnStart context.CancelFunc
	// The stop observation records what the reclaim path actually handed to the
	// stop capability: a cancelled context, an unbounded context, or neither.
	stopCtxErr         error
	stopCtxHasDeadline bool
	stopCtxDeadline    time.Time
}

var _ agentRunnerLaunchReclaimer = (*reclaimingLaunchStub)(nil)

func newReclaimingLaunchStub() *reclaimingLaunchStub {
	return &reclaimingLaunchStub{recordingLaunchStub: &recordingLaunchStub{}}
}

// StartAgentRunnerProcess cancels the dispatch context first when a test asked
// for that, then delegates to the shared recording start so the start behaviour
// of both stubs cannot drift apart.
func (stub *reclaimingLaunchStub) StartAgentRunnerProcess(ctx context.Context, launch chain.AgentRunnerLaunchRequest) (chain.AgentRunnerLaunchResult, error) {
	stub.stopMu.Lock()
	cancel := stub.cancelOnStart
	stub.stopMu.Unlock()
	if cancel != nil {
		cancel()
	}
	return stub.recordingLaunchStub.StartAgentRunnerProcess(ctx, launch)
}

func (stub *reclaimingLaunchStub) StopAgentRunnerProcess(ctx context.Context, launchID string) (guard.TerminationEvidence, error) {
	deadline, hasDeadline := ctx.Deadline()
	ctxErr := ctx.Err()
	stub.stopMu.Lock()
	defer stub.stopMu.Unlock()
	stub.stopCalls++
	stub.stopIDs = append(stub.stopIDs, launchID)
	stub.stopCtxErr = ctxErr
	stub.stopCtxHasDeadline = hasDeadline
	stub.stopCtxDeadline = deadline
	if stub.stopErr != nil {
		return guard.TerminationEvidence{Outcome: "uncertain"}, stub.stopErr
	}
	outcome := stub.stopOutcome
	if outcome == "" {
		outcome = "stopped"
	}
	return guard.TerminationEvidence{Outcome: outcome}, nil
}

func (stub *reclaimingLaunchStub) stopCallCount() int {
	stub.stopMu.Lock()
	defer stub.stopMu.Unlock()
	return stub.stopCalls
}

func (stub *reclaimingLaunchStub) stopCallIDs() []string {
	stub.stopMu.Lock()
	defer stub.stopMu.Unlock()
	return append([]string(nil), stub.stopIDs...)
}

// stopContextObservation reports the state of the context the stop capability
// received: whether it carried a deadline, what its error was, and that
// deadline. It never infers anything the stub did not observe directly.
func (stub *reclaimingLaunchStub) stopContextObservation() (deadline time.Time, hasDeadline bool, ctxErr error) {
	stub.stopMu.Lock()
	defer stub.stopMu.Unlock()
	return stub.stopCtxDeadline, stub.stopCtxHasDeadline, stub.stopCtxErr
}

func reclaimingDispatcherFixture(t *testing.T) (*AgentRunnerReplayDispatcher, *reclaimingLaunchStub, chain.StepRequest) {
	t.Helper()
	admission, request := compositeAdmissionFixture(t)
	store := replayStoreMustNew(t, t.TempDir())
	launcher := newReclaimingLaunchStub()
	dispatcher := &AgentRunnerReplayDispatcher{
		Store:          store,
		Launcher:       launcher,
		LedgerSnapshot: func() (chain.AuthorizationLedger, error) { return admission.AuthorizationLedger, nil },
		CostLedger:     admission.CostLedger,
		RequestRecord:  admission.RequestRecord,
		Clock:          admission.Clock,
	}
	return dispatcher, launcher, request
}

// assertReclaimOnce proves the reclaim ran exactly once against the
// dispatcher-derived launch id, and that the unproven dispatch was neither
// retried nor relaunched by a second attempt.
func assertReclaimOnce(t *testing.T, dispatcher *AgentRunnerReplayDispatcher, launcher *reclaimingLaunchStub, request chain.StepRequest) {
	t.Helper()
	if calls, ids := launcher.stopCallCount(), launcher.stopCallIDs(); calls != 1 {
		t.Fatalf("expected exactly one reclaim, got %d (ids=%v)", calls, ids)
	}
	if ids := launcher.stopCallIDs(); len(ids) != 1 || ids[0] != "launch-request-one" {
		t.Fatalf("reclaim must stop the dispatcher-derived launch id launch-request-one, got %v", ids)
	}
	if calls := launcher.callCount(); calls != 1 {
		t.Fatalf("expected exactly one start, got %d", calls)
	}
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
		t.Fatalf("an unproven identity must block redispatch, got %v", err)
	}
	if calls := launcher.callCount(); calls != 1 {
		t.Fatalf("an unproven identity must never relaunch, got %d starts", calls)
	}
	if calls := launcher.stopCallCount(); calls != 1 {
		t.Fatalf("a blocked redispatch must not reclaim again, got %d reclaims", calls)
	}
	if _, found, err := dispatcher.Store.LaunchReceipt("request-one"); err != nil || !found {
		t.Fatalf("the launch intent must stay retained: found=%v err=%v", found, err)
	}
}

// assertNoIdentityPublished proves the failed dispatch published nothing under
// its own launch id, which is the case whenever no identity record pre-exists.
func assertNoIdentityPublished(t *testing.T, dispatcher *AgentRunnerReplayDispatcher) {
	t.Helper()
	identity, found, err := dispatcher.Store.LaunchIdentity("launch-request-one")
	if err != nil {
		t.Fatalf("identity lookup failed: %v", err)
	}
	if found {
		t.Fatalf("no identity may be published for an unproven spawn, got %+v", identity.Identity)
	}
}

// TestAgentRunnerReplayDispatcherReclaimNeverFiresOnSuccess is the positive
// control: a launcher with the stop capability must not be stopped at all when
// the identity publishes normally.
func TestAgentRunnerReplayDispatcherReclaimNeverFiresOnSuccess(t *testing.T) {
	dispatcher, launcher, request := reclaimingDispatcherFixture(t)
	result, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if err != nil {
		t.Fatalf("expected the dispatch to succeed, got %v", err)
	}
	if calls := launcher.callCount(); calls != 1 {
		t.Fatalf("expected exactly one start, got %d", calls)
	}
	if calls := launcher.stopCallCount(); calls != 0 {
		t.Fatalf("a published identity must not reclaim, got %d reclaims (ids=%v)", calls, launcher.stopCallIDs())
	}
	if _, found, err := dispatcher.Store.LaunchIdentity("launch-request-one"); err != nil || !found {
		t.Fatalf("the identity must be published: found=%v err=%v", found, err)
	}
	if len(result.Evidence) != 3 {
		t.Fatalf("expected the full evidence chain, got %v", result.Evidence)
	}
}

func TestAgentRunnerReplayDispatcherLaunchIDMismatchReclaimsOnce(t *testing.T) {
	dispatcher, launcher, request := reclaimingDispatcherFixture(t)
	launcher.launchID = "launch-foreign-one"
	_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) {
		t.Fatalf("expected an unproven identity for a mismatched launchId, got %v", err)
	}
	if !strings.Contains(err.Error(), "launch-foreign-one") || !strings.Contains(err.Error(), "spawned process reclaimed (stopped)") {
		t.Fatalf("expected the mismatch cause and a proven reclaim in the message, got %v", err)
	}
	assertReclaimOnce(t, dispatcher, launcher, request)
	assertNoIdentityPublished(t, dispatcher)
}

func TestAgentRunnerReplayDispatcherIdentityConstructionFailureReclaimsOnce(t *testing.T) {
	dispatcher, launcher, request := reclaimingDispatcherFixture(t)
	// A whitespace-only process identity is refused by the record constructor, so
	// the failure happens after the spawn and before any publication.
	launcher.processID = "   "
	_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) {
		t.Fatalf("expected an unproven identity for an invalid process identity, got %v", err)
	}
	if !errors.Is(err, ErrInvalidAgentRunnerLaunchIdentity) || !strings.Contains(err.Error(), "empty process identity") {
		t.Fatalf("expected the construction cause to be preserved, got %v", err)
	}
	assertReclaimOnce(t, dispatcher, launcher, request)
	assertNoIdentityPublished(t, dispatcher)
}

func TestAgentRunnerReplayDispatcherIdentityPublicationFailureReclaimsOnce(t *testing.T) {
	dispatcher, launcher, request := reclaimingDispatcherFixture(t)
	// Pre-write a foreign identity for the same launch id: the publication then
	// conflicts, which is the only identity failure the store reports after a
	// spawn that did publish its own process id.
	receipt := launchReceiptRecordFor(t, dispatcher.RequestRecord)
	foreign, err := NewAgentRunnerLaunchIdentityRecord(AgentRunnerLaunchIdentity{
		LaunchID:         receipt.Receipt.LaunchID,
		RequestID:        receipt.Receipt.RequestID,
		IntentRecordHash: receipt.RecordHash,
		ProcessID:        "process-other-one",
		StartedAt:        "2026-09-11T08:05:02Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	path, err := dispatcher.Store.launchIdentityPath(receipt.Receipt.LaunchID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(path, foreign); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	_, err = dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) || !errors.Is(err, ErrAgentRunnerLaunchIdentityConflict) {
		t.Fatalf("expected the identity conflict to surface as unproven, got %v", err)
	}
	assertReclaimOnce(t, dispatcher, launcher, request)

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(before) {
		t.Fatalf("a failed publication must not rewrite the foreign record:\nbefore=%s\nafter=%s", before, after)
	}
}

func TestAgentRunnerReplayDispatcherReclaimFailureKeepsBothChains(t *testing.T) {
	dispatcher, launcher, request := reclaimingDispatcherFixture(t)
	launcher.launchID = "launch-foreign-one"
	stopErr := errors.New("stop channel refused")
	launcher.stopErr = stopErr
	_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) {
		t.Fatalf("the original cause must survive a failed reclaim, got %v", err)
	}
	if !strings.Contains(err.Error(), "spawned process reclaim failed") || !strings.Contains(err.Error(), stopErr.Error()) {
		t.Fatalf("expected the stop error in the message, got %v", err)
	}
	assertReclaimOnce(t, dispatcher, launcher, request)
}

func TestAgentRunnerReplayDispatcherReclaimUnprovenIsReported(t *testing.T) {
	dispatcher, launcher, request := reclaimingDispatcherFixture(t)
	launcher.launchID = "launch-foreign-one"
	launcher.stopOutcome = "uncertain"
	_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) {
		t.Fatalf("an unproven reclaim must keep the identity sentinel, got %v", err)
	}
	if !strings.Contains(err.Error(), "reclaim unproven") || !strings.Contains(err.Error(), "uncertain") {
		t.Fatalf("expected the unproven outcome in the message, got %v", err)
	}
	assertReclaimOnce(t, dispatcher, launcher, request)
}

func TestAgentRunnerReplayDispatcherWithoutStopCapabilityReportsUnprovenReclaim(t *testing.T) {
	dispatcher, _, launcher, request := replayDispatcherFixture(t)
	reclaimer := reflect.TypeOf((*agentRunnerLaunchReclaimer)(nil)).Elem()
	if reflect.TypeOf(launcher).Implements(reclaimer) {
		t.Fatalf("the plain stub must not expose the stop capability, it implements %v", reclaimer)
	}
	launcher.launchID = "launch-foreign-one"
	_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) {
		t.Fatalf("expected an unproven identity without a stop capability, got %v", err)
	}
	if !strings.Contains(err.Error(), "no stop capability") || !strings.Contains(err.Error(), "was not reclaimed") {
		t.Fatalf("expected the missing capability to be stated, got %v", err)
	}
	if calls := launcher.callCount(); calls != 1 {
		t.Fatalf("expected exactly one start, got %d", calls)
	}
	if _, err := dispatcher.DispatchAgentRunner(context.Background(), request); !errors.Is(err, ErrAgentRunnerDispatchUnknownBlock) {
		t.Fatalf("an unproven identity must block redispatch, got %v", err)
	}
	if calls := launcher.callCount(); calls != 1 {
		t.Fatalf("a launcher without a stop capability must never be relaunched, got %d starts", calls)
	}
	if _, found, err := dispatcher.Store.LaunchReceipt("request-one"); err != nil || !found {
		t.Fatalf("the launch intent must stay retained: found=%v err=%v", found, err)
	}
}

// TestAgentRunnerReplayDispatcherReclaimsRealProcessWhenIdentityPublicationFails
// is the reality check for the reclaim path: a real pinned-CLI launcher spawns a
// real process, the identity publication conflicts on a foreign record, and the
// dispatcher must leave no live process, no owned registry handle and no claimed
// launch slot behind.
func TestAgentRunnerReplayDispatcherReclaimsRealProcessWhenIdentityPublicationFails(t *testing.T) {
	admission, request := compositeAdmissionFixture(t)
	request.Workspace.Root = t.TempDir()
	runRoot := t.TempDir()
	store := replayStoreForRunRoot(t, runRoot)
	// The dispatcher publishes request records, so the store needs a proven
	// durability gate; the launcher keeps using the very same store and run root.
	store.publicationDurability = PublishDurabilityProven
	launcher, err := NewAgentRunnerPinnedCLILauncher(agentRunnerPinnedCLIConfigForMode("run"), store)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		// Never leave a helper process behind, whatever the assertions below did.
		_, _ = launcher.StopAgentRunnerProcess(context.Background(), "launch-request-one")
	})
	dispatcher := &AgentRunnerReplayDispatcher{
		Store:          store,
		Launcher:       launcher,
		LedgerSnapshot: func() (chain.AuthorizationLedger, error) { return admission.AuthorizationLedger, nil },
		CostLedger:     admission.CostLedger,
		RequestRecord:  admission.RequestRecord,
		Clock:          admission.Clock,
	}

	receipt := launchReceiptRecordFor(t, admission.RequestRecord)
	foreign, err := NewAgentRunnerLaunchIdentityRecord(AgentRunnerLaunchIdentity{
		LaunchID:         receipt.Receipt.LaunchID,
		RequestID:        receipt.Receipt.RequestID,
		IntentRecordHash: receipt.RecordHash,
		ProcessID:        "process-other-one",
		StartedAt:        "2026-09-11T08:05:02Z",
	})
	if err != nil {
		t.Fatal(err)
	}
	identityPath, err := store.launchIdentityPath(receipt.Receipt.LaunchID)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeReplayRecordNoReplace(identityPath, foreign); err != nil {
		t.Fatal(err)
	}
	foreignBefore, err := os.ReadFile(identityPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = dispatcher.DispatchAgentRunner(context.Background(), request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) || !errors.Is(err, ErrAgentRunnerLaunchIdentityConflict) {
		t.Fatalf("expected the identity conflict to surface as unproven, got %v", err)
	}
	if !strings.Contains(err.Error(), "spawned process reclaimed") {
		t.Fatalf("expected the reclaim outcome in the message, got %v", err)
	}

	// The launcher published its own durable identity mirror for the spawn, and it
	// must now name a process that is gone.
	mirror, err := loadAgentRunnerProcessIdentityMirror(runRoot)
	if err != nil {
		t.Fatalf("the launcher must have mirrored the spawned process: %v", err)
	}
	alive, aliveErr := guard.ProcessAlive(mirror)
	if aliveErr != nil {
		t.Fatalf("the mirrored process identity must stay examinable: %+v (%v)", mirror, aliveErr)
	}
	if alive {
		t.Fatalf("the reclaimed process must not be alive: %+v", mirror)
	}
	if held := launcher.registry.len(); held != 0 {
		t.Fatalf("a reclaimed launch must release its handle, got %d held", held)
	}
	foreignAfter, err := os.ReadFile(identityPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(foreignAfter) != string(foreignBefore) {
		t.Fatalf("a failed publication must not rewrite the foreign record:\nbefore=%s\nafter=%s", foreignBefore, foreignAfter)
	}

	// A released slot is reusable: the same run root must accept a fresh launch.
	fresh := agentRunnerPinnedCLIRequest(request.Workspace.Root)
	fresh.RequestID = "request-two"
	relaunched, err := launcher.StartAgentRunnerProcess(context.Background(), fresh)
	if err != nil {
		t.Fatalf("the reclaimed launch must release its slot: %v", err)
	}
	t.Cleanup(func() {
		_, _ = launcher.StopAgentRunnerProcess(context.Background(), relaunched.LaunchID)
	})
	if held := launcher.registry.len(); held != 1 {
		t.Fatalf("the fresh launch must hold exactly one handle, got %d", held)
	}
}

// TestAgentRunnerReplayDispatcherReclaimSurvivesCallerCancellation pins the
// shield around the reclaim stop. The launcher cancels the dispatch context from
// inside the start call and then reports a mismatched launch id, so the reclaim
// branch is entered while the caller context is already dead. Deleting
// context.WithoutCancel from the production helper makes this test red.
func TestAgentRunnerReplayDispatcherReclaimSurvivesCallerCancellation(t *testing.T) {
	dispatcher, launcher, request := reclaimingDispatcherFixture(t)
	// A pinned grace makes the deadline assertion exact instead of "some deadline".
	dispatcher.Process.Grace = 250 * time.Millisecond
	dispatchCtx, cancelDispatch := context.WithCancel(context.Background())
	defer cancelDispatch()
	launcher.launchID = "launch-foreign-one"
	launcher.cancelOnStart = cancelDispatch

	_, err := dispatcher.DispatchAgentRunner(dispatchCtx, request)
	if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) {
		t.Fatalf("a cancelled caller must still surface the unproven identity, got %v", err)
	}
	if dispatchCtx.Err() == nil {
		t.Fatalf("test precondition failed: the stub must have cancelled the dispatch context, but it is still live")
	}
	if calls := launcher.stopCallCount(); calls != 1 {
		t.Fatalf("expected exactly one reclaim of the unproven spawn, got %d (ids=%v)", calls, launcher.stopCallIDs())
	}

	deadline, hasDeadline, ctxErr := launcher.stopContextObservation()
	if ctxErr != nil {
		t.Fatalf("shield failed: the reclaim stop must receive a context that is not cancelled, but Err()=%v (the caller context was %v)", ctxErr, dispatchCtx.Err())
	}
	if !hasDeadline {
		t.Fatalf("deadline failed: the reclaim stop must receive a time-bounded context, but Deadline() reported none")
	}
	remaining := time.Until(deadline)
	if remaining <= 0 {
		t.Fatalf("deadline failed: the reclaim stop context was already expired at stop time (deadline %s, remaining %s)", deadline, remaining)
	}
	// The bound is grace+1s measured from the reclaim; anything larger means the
	// stop is no longer bounded by the configured grace.
	if maxRemaining := dispatcher.Process.Grace + time.Second; remaining > maxRemaining {
		t.Fatalf("deadline failed: the reclaim stop context exceeds grace+1s, remaining %s > %s", remaining, maxRemaining)
	}
	assertNoIdentityPublished(t, dispatcher)
}

// typedNilReclaimerStub is a launcher that stays usable through a nil receiver,
// so a typed-nil pointer can be driven through the reclaim path in this package.
// It starts nothing: it only reports a mismatched launch id to force the reclaim
// branch, and its stop reports that a nil launcher cannot prove a termination.
type typedNilReclaimerStub struct{}

var _ agentRunnerLaunchReclaimer = (*typedNilReclaimerStub)(nil)
var _ chain.AgentRunnerLauncher = (*typedNilReclaimerStub)(nil)

func (stub *typedNilReclaimerStub) StartAgentRunnerProcess(_ context.Context, _ chain.AgentRunnerLaunchRequest) (chain.AgentRunnerLaunchResult, error) {
	return chain.AgentRunnerLaunchResult{
		LaunchID:  "launch-foreign-one",
		ProcessID: "process-typed-nil-one",
		StartedAt: time.Date(2026, 9, 11, 8, 6, 0, 0, time.UTC),
	}, nil
}

func (stub *typedNilReclaimerStub) StopAgentRunnerProcess(_ context.Context, _ string) (guard.TerminationEvidence, error) {
	return guard.TerminationEvidence{Outcome: "uncertain"}, errors.New("typed nil launcher cannot prove a stop")
}

// TestAgentRunnerReplayDispatcherTypedNilReclaimerIsReportedFailClosed pins what
// a typed-nil launcher does. A typed nil satisfies the launcher interface and the
// optional reclaim interface, so the dispatcher's plain nil guard cannot see it;
// the requirement is that such a launcher still fails closed and is reported
// instead of panicking.
func TestAgentRunnerReplayDispatcherTypedNilReclaimerIsReportedFailClosed(t *testing.T) {
	// The concrete pinned-CLI launcher refuses a nil receiver itself, so its
	// start returns an error and the reclaim helper is never entered. That is the
	// nearest reachable behaviour for that type and is asserted plainly here.
	t.Run("pinned CLI nil pointer is refused before any reclaim", func(t *testing.T) {
		defer func() {
			if recovered := recover(); recovered != nil {
				t.Fatalf("a typed nil launcher must fail closed instead of panicking, got panic: %v", recovered)
			}
		}()
		dispatcher, _, _, request := replayDispatcherFixture(t)
		dispatcher.Launcher = (*AgentRunnerPinnedCLILauncher)(nil)
		if dispatcher.Launcher == nil {
			t.Fatalf("test precondition failed: a typed nil must not compare equal to the nil interface")
		}
		_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
		if !errors.Is(err, ErrAgentRunnerLaunchFailed) || !errors.Is(err, ErrInvalidAgentRunnerLauncher) {
			t.Fatalf("expected the pinned launcher to refuse a nil receiver as a failed launch, got %v", err)
		}
		if !strings.Contains(err.Error(), "nil launcher") {
			t.Fatalf("expected the nil receiver refusal to be stated, got %v", err)
		}
	})

	// A launcher that tolerates a nil receiver reaches the reclaim helper, where
	// the typed nil satisfies the reclaim interface. The requirement is no panic
	// plus an explicit reclaim outcome in the reported message.
	t.Run("typed nil reclaimer reports the failed reclaim", func(t *testing.T) {
		defer func() {
			if recovered := recover(); recovered != nil {
				t.Fatalf("a typed nil reclaimer must be reported instead of panicking, got panic: %v", recovered)
			}
		}()
		dispatcher, _, _, request := replayDispatcherFixture(t)
		var typedNil *typedNilReclaimerStub
		dispatcher.Launcher = typedNil
		if dispatcher.Launcher == nil {
			t.Fatalf("test precondition failed: a typed nil must not compare equal to the nil interface")
		}

		_, err := dispatcher.DispatchAgentRunner(context.Background(), request)
		if !errors.Is(err, ErrAgentRunnerLaunchIdentityUnproven) {
			t.Fatalf("expected the unproven identity sentinel for a typed nil reclaimer, got %v", err)
		}
		if !strings.Contains(err.Error(), "reclaim failed") || !strings.Contains(err.Error(), "typed nil launcher cannot prove a stop") {
			t.Fatalf("expected the reclaim outcome to be reported, got %v", err)
		}
	})
}
