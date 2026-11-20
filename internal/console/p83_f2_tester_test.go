package console

// P83 F2 counterexample: the dispatch-time ledger snapshots must reread the
// operator-written files, not return the assembly-time in-memory copies.
//
// The end-to-end leg of F2 (a revocation appended after R is published) needs
// publish durability, which the Windows replay store declares unproven, so the
// production dispatch refuses before R there. These tests instead drive the two
// snapshots the assembly hands to the dispatcher - the only place the reread can
// happen - and prove they observe the on-disk state and fail closed when it
// cannot be read.

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/snapshot"
)

// p83TesterDispatcher binds the assembly ports for a single-task plan and returns
// the dispatcher the production dispatch would use.
func p83TesterDispatcher(t *testing.T, options AgentRunnerAssemblyOptions) *adapters.AgentRunnerReplayDispatcher {
	t.Helper()
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(options.RunRoot)); err != nil {
		t.Fatal(err)
	}
	store, err := adapters.NewAgentRunnerReplayStoreForRun(options.RunRoot, options.RunID)
	if err != nil {
		t.Fatal(err)
	}
	snapshotStore, err := snapshot.NewStore(filepath.Join(options.RunRoot, "tester-snapshot-store"), 0)
	if err != nil {
		t.Fatal(err)
	}
	snapshotRoot := filepath.Join(options.RunRoot, "tester-snapshots")
	if err := os.MkdirAll(snapshotRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(agentRunnerWorkspacePath(options.RunRoot, assemblyTaskID, 1), 0o755); err != nil {
		t.Fatal(err)
	}
	tasks, err := agentRunnerBindTasks(
		options,
		[]agentRunnerTaskBinding{{TaskID: assemblyTaskID, StepID: assemblyStepID}},
		[]adapters.AgentRunnerRequestRecord{options.RequestRecord},
		store, nil, snapshotStore, snapshotRoot,
		func() time.Time { return time.Now().UTC() },
	)
	if err != nil {
		t.Fatalf("bind the task ports: %v", err)
	}
	dispatcher, ok := tasks[0].dispatcher.(*adapters.AgentRunnerReplayDispatcher)
	if !ok {
		t.Fatalf("the bound dispatch port is not a replay dispatcher: %T", tasks[0].dispatcher)
	}
	return dispatcher
}

// TestAgentRunnerDispatchLedgerSnapshotRereadsOperatorRevocation is the F2
// counterexample that runs on Windows: the assembly holds an active grant in
// memory, the operator appends a revocation to the ledger file, and the dispatch
// snapshot must observe it.
// Mutation check: returning options.AuthorizationLedger from the closure keeps
// the grant active and turns this red.
func TestAgentRunnerDispatchLedgerSnapshotRereadsOperatorRevocation(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{declaredExecutable: os.Args[0]})
	grant := fixture.options.AuthorizationLedger.Grants[0]
	ledgerPath := filepath.Join(fixture.runRoot, "authorization-ledger.jsonl")
	if err := appendAuthorizationRecord(ledgerPath, grant); err != nil {
		t.Fatal(err)
	}
	options := fixture.options
	options.AuthorizationLedgerPath = ledgerPath
	dispatcher := p83TesterDispatcher(t, options)

	now := time.Now().UTC()
	snapshotLedger, err := dispatcher.LedgerSnapshot()
	if err != nil {
		t.Fatalf("the reread must accept a readable ledger: %v", err)
	}
	if status, evalErr := chain.EvaluateGrant(grant.Grant.AuthorizationID, snapshotLedger, now); evalErr != nil || status != chain.GrantActive {
		t.Fatalf("the operator grant must be active before the revocation, status=%v err=%v", status, evalErr)
	}

	revocation, err := chain.NewRevocationAuthorizationRecord(chain.AuthorizationRevocation{
		RecordID:          "p83-tester-revocation-1",
		Kind:              "revocation",
		IssuedAt:          now.Format(timestampLayout),
		IssuedBy:          evidence.Actor{Type: "operator", ID: "p83-tester-operator"},
		AuthorizationID:   grant.Grant.AuthorizationID,
		AuthorizationHash: grant.RecordHash,
		ReasonCode:        "operator-revoked",
		StopDisposition:   "not-required",
	}, time.Now, func() string { return "p83-tester-revocation-record-1" })
	if err != nil {
		t.Fatalf("build revocation: %v", err)
	}
	if err := appendAuthorizationRecord(ledgerPath, revocation); err != nil {
		t.Fatal(err)
	}

	// The assembly copy is unchanged: only the reread can see the revocation.
	if status, _ := chain.EvaluateGrant(grant.Grant.AuthorizationID, options.AuthorizationLedger, now); status != chain.GrantActive {
		t.Fatalf("the fixture must still hold the grant in memory, status=%v", status)
	}
	snapshotLedger, err = dispatcher.LedgerSnapshot()
	if err != nil {
		t.Fatalf("the reread must not fail on a readable ledger: %v", err)
	}
	if status, evalErr := chain.EvaluateGrant(grant.Grant.AuthorizationID, snapshotLedger, now); evalErr != nil || status == chain.GrantActive {
		t.Fatalf("the dispatch snapshot must observe the revocation appended after assembly, status=%v err=%v", status, evalErr)
	}
}

// TestAgentRunnerDispatchCostSnapshotRefusesUnreadableLedger is the budget half of
// F2: a configured cost ledger that cannot be read is a refusal, never a silent
// fallback to the assembly-time copy (which would still show an outstanding
// reservation and permit the launch).
// Mutation check: returning options.CostLedger when the read fails turns this red.
func TestAgentRunnerDispatchCostSnapshotRefusesUnreadableLedger(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{declaredExecutable: os.Args[0]})
	options := fixture.options
	options.CostLedgerPath = filepath.Join(fixture.runRoot, "cost-ledger-that-does-not-exist.json")
	dispatcher := p83TesterDispatcher(t, options)

	ledger, err := dispatcher.CostLedgerSnapshot()
	if err == nil {
		t.Fatalf("an unreadable cost ledger must refuse the launch, got ledger=%+v", ledger)
	}
	if ledger != nil {
		t.Fatalf("a refused cost snapshot must return no ledger, got %+v", ledger)
	}
}
