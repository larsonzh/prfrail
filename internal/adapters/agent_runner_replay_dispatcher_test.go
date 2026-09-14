package adapters

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
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
