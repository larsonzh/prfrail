package console

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/guard"
	"github.com/larsonzh/prfrail/internal/snapshot"
	"github.com/larsonzh/prfrail/internal/tickets"
)

// This file is the ⑨ evidence carrier for the production AgentRunner assembly.
// Every record is DECLARED (declared, not probed): nothing here calls an AI
// provider, nothing touches the network, and the pinned workload is the
// deterministic offline agent-stub. The reverse legs prove that a refusal writes
// no replay identity and spawns no process.

const (
	assemblyRunID      = "run-t027"
	assemblyRequestID  = "request-run-t027"
	assemblyTaskID     = "task-assembly"
	assemblyStepID     = "code-assembly"
	assemblyAdapterID  = "t027-declared-adapter"
	assemblyProfileID  = "t027-declared-profile"
	assemblyStubVer    = "agent-stub-0.1.0"
	assemblyDeclaredDo = "proofrail:t027-assembly-declared:1\n"
	// Effect mapping pair mirrored from unexported constants in
	// internal/adapters/agent_runner.go. They are contract values, not tunables,
	// and a request record is refused without them.
	assemblyEffectMappingVersion = "1"
	assemblyEffectMappingHash    = "sha256:2c9609ee374bfc79bd0ed997aea6cfe4865c82ed0216d1a4a288139e0b0c92d4"
)

var (
	assemblyStubOnce sync.Once
	assemblyStubPath string
	assemblyStubErr  error
)

// buildAssemblyStub builds tools/agent-stub once. It is the pinned offline CLI
// the production launcher runs; no external dependency and no network is used.
func buildAssemblyStub(t *testing.T) string {
	t.Helper()
	assemblyStubOnce.Do(func() {
		root, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			assemblyStubErr = err
			return
		}
		target := filepath.Join(os.TempDir(), "proofrail-agent-stub-t027")
		if err := os.MkdirAll(target, 0o755); err != nil {
			assemblyStubErr = err
			return
		}
		name := "agent-stub"
		if strings.EqualFold(filepath.Ext(os.Args[0]), ".exe") {
			name += ".exe"
		}
		assemblyStubPath = filepath.Join(target, name)
		build := exec.Command("go", "build", "-o", assemblyStubPath, "./tools/agent-stub")
		build.Dir = root
		if output, err := build.CombinedOutput(); err != nil {
			assemblyStubErr = fmt.Errorf("build agent-stub: %v: %s", err, output)
		}
	})
	if assemblyStubErr != nil {
		t.Skipf("the assembly tests need a buildable agent-stub: %v", assemblyStubErr)
	}
	return assemblyStubPath
}

type assemblyKnobs struct {
	capabilityDisposition string
	capabilityPlatformOS  string
	enforcementMode       string // "" | "valid" | "stale"
	availabilityStatus    string // "" means available
	workspaceOverride     string
	declaredExecutable    string
	declaredVersion       string
	runArgs               []string
}

type assemblyFixture struct {
	runRoot string
	options AgentRunnerAssemblyOptions
}

func assemblyDefinition() chain.Definition {
	return chain.Definition{
		ID: "chain-assembly",
		Tasks: []chain.Task{{
			ID:    assemblyTaskID,
			Steps: []chain.Step{{ID: assemblyStepID, Kind: "code", Mode: chain.IsolatedWorkspace}},
		}},
	}
}

func assemblyDeclaredDigest(parts ...string) string {
	return digest(assemblyDeclaredDo, parts...)
}

func assemblyFileHash(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return evidence.Digest("", data)
}

func assemblyVerifiedFinding(note string) adapters.AgentRunnerCapabilityFinding {
	return adapters.AgentRunnerCapabilityFinding{Status: "verified", Evidence: []string{assemblyDeclaredDigest("finding", note)}}
}

func assemblyUnsupportedFinding(note string) adapters.AgentRunnerCapabilityFinding {
	return adapters.AgentRunnerCapabilityFinding{Status: "unsupported", Evidence: []string{assemblyDeclaredDigest("finding", note)}}
}

// assemblyCapabilityMatrix declares every capability dimension. A compatible
// disposition requires all twelve findings verified; the blocked disposition is
// the T026 reality (tool/network control unsupported), which forces admission to
// demand external enforcement.
func assemblyCapabilityMatrix(disposition string) adapters.AgentRunnerCapabilityMatrix {
	matrix := adapters.AgentRunnerCapabilityMatrix{
		Noninteractive:            assemblyVerifiedFinding("noninteractive"),
		CWD:                       assemblyVerifiedFinding("cwd"),
		EventStreamOrCompleteLogs: assemblyVerifiedFinding("events"),
		SessionCreate:             assemblyVerifiedFinding("session-create"),
		SessionResume:             assemblyVerifiedFinding("session-resume"),
		Cancellation:              assemblyVerifiedFinding("cancellation"),
		ProcessTreeStop:           assemblyVerifiedFinding("tree-stop"),
		PermissionControl:         assemblyVerifiedFinding("permission"),
		Usage:                     assemblyVerifiedFinding("usage"),
		UnattendedConfirmations:   assemblyVerifiedFinding("unattended"),
		ToolControl:               assemblyVerifiedFinding("tool-control"),
		NetworkControl:            assemblyVerifiedFinding("network-control"),
	}
	if disposition != "compatible" {
		matrix.ToolControl = assemblyUnsupportedFinding("tool-control")
		matrix.NetworkControl = assemblyUnsupportedFinding("network-control")
	}
	return matrix
}

func newAssemblyFixture(t *testing.T, knobs assemblyKnobs) *assemblyFixture {
	t.Helper()
	if knobs.capabilityDisposition == "" {
		knobs.capabilityDisposition = "compatible"
	}
	if knobs.capabilityPlatformOS == "" {
		knobs.capabilityPlatformOS = runtime.GOOS
	}
	if knobs.declaredExecutable == "" {
		knobs.declaredExecutable = os.Args[0]
	}
	if knobs.declaredVersion == "" {
		knobs.declaredVersion = assemblyStubVer
	}

	runRoot := t.TempDir()
	workspace := knobs.workspaceOverride
	if workspace == "" {
		workspace = agentRunnerWorkspacePath(runRoot, assemblyTaskID, 1)
	}
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	baselineStore, err := snapshot.NewStore(filepath.Join(runRoot, "baseline-store"), 0)
	if err != nil {
		t.Fatal(err)
	}
	parent, err := adapters.AgentRunnerWorkspaceBaselineDigest(context.Background(), baselineStore, workspace, assemblyRunID, assemblyTaskID, 1)
	if err != nil {
		t.Fatal(err)
	}

	base := time.Now().UTC()
	executableHash := assemblyFileHash(t, knobs.declaredExecutable)
	configHash := assemblyDeclaredDigest("config", knobs.declaredExecutable, knobs.declaredVersion)
	platformHash := agentRunnerPlatformHash()
	policyHash := assemblyDeclaredDigest("policy")
	amount := int64(1000)

	grant, err := chain.NewGrantAuthorizationRecord(chain.AuthorizationGrant{
		RecordID:        "t027-grant-1",
		Kind:            "grant",
		IssuedAt:        base.Add(-time.Minute).UTC().Format(timestampLayout),
		IssuedBy:        evidence.Actor{Type: "operator", ID: "t027-declared-operator"},
		AuthorizationID: "t027-authorization-1",
		RunID:           assemblyRunID,
		RunManifestHash: parent,
		PolicyHash:      policyHash,
		Scope: chain.AuthorizationScope{
			TaskIDs:       []string{assemblyTaskID},
			StepIDs:       []string{assemblyStepID},
			TargetIDs:     []string{"workspace-files"},
			EffectClasses: []string{"local-discardable", "read-only"},
			Budget: chain.AuthorizationBudgetLimit{
				ModelCalls: 2, Tokens: 1000, WallClockMs: 600000, Attempts: 1,
				Currency: "USD", AmountMicros: &amount,
			},
		},
		ExpiresAt: base.Add(time.Hour).UTC().Format(timestampLayout),
		Evidence:  []string{assemblyDeclaredDigest("grant-evidence", "t027-grant-1")},
	}, func() time.Time { return base }, func() string { return "t027-grant-record-1" })
	if err != nil {
		t.Fatalf("declared grant: %v", err)
	}

	ledger, err := tickets.NewCostLedger(tickets.CostLedgerConfig{
		LedgerID:       "t027-cost-ledger",
		CreatedAt:      base,
		ScopeKind:      tickets.CostScopeRun,
		ScopeHash:      parent,
		Currency:       "USD",
		PricingMode:    tickets.CostPricingDocumented,
		PricingVersion: base.Format("2006-01-02"),
	})
	if err != nil {
		t.Fatalf("declared cost ledger: %v", err)
	}
	if _, err := ledger.AppendAllocation(tickets.AllocationInput{
		EntryID:           "t027-allocation-1",
		OccurredAt:        base,
		AuthorizationHash: grant.RecordHash,
		Limits: tickets.CostLimits{
			AmountMicros: &amount, ModelCalls: 2, Tokens: 1000, WallClockMs: 600000, Attempts: 1,
		},
	}); err != nil {
		t.Fatalf("declared allocation: %v", err)
	}
	reservation, err := ledger.Reserve(tickets.ReservationInput{
		EntryID:              "t027-reservation-1",
		OccurredAt:           base,
		RunID:                assemblyRunID,
		RequestID:            assemblyRequestID,
		IdempotencyKey:       "t027-reserve-1",
		AuthorizationHash:    grant.RecordHash,
		ReservedAmountMicros: &amount,
		ReservedCalls:        2,
		ReservedTokens:       1000,
		PricingEvidence:      []string{assemblyDeclaredDigest("pricing", "declared-zero-cost-offline-workload")},
	})
	if err != nil {
		t.Fatalf("declared reservation: %v", err)
	}

	availabilityStatus := knobs.availabilityStatus
	if availabilityStatus == "" {
		availabilityStatus = "available"
	}
	availabilityBody := adapters.AIAvailability{
		ProbeID:           "t027-availability",
		ProfileID:         assemblyProfileID,
		ProfileConfigHash: configHash,
		Channel:           "agent-runner-cli",
		ProbedAt:          base.Add(-time.Minute).UTC().Format(timestampLayout),
		Status:            availabilityStatus,
		RequestsUsed:      1,
		Evidence:          []string{assemblyDeclaredDigest("availability-evidence", "declared, not probed")},
	}
	if availabilityStatus != "available" {
		reason := "cli_unavailable"
		availabilityBody.Reason = &reason
	}
	availability, err := adapters.NewAIAvailabilityRecord(availabilityBody)
	if err != nil {
		t.Fatalf("declared availability: %v", err)
	}

	capability, err := adapters.NewAgentRunnerCapabilityRecord(adapters.AgentRunnerCapability{
		ProbeID:        "t027-capability",
		AdapterID:      assemblyAdapterID,
		ProbedAt:       base.Add(-time.Minute).UTC().Format(timestampLayout),
		OS:             knobs.capabilityPlatformOS,
		Arch:           runtime.GOARCH,
		PlatformHash:   platformHash,
		ExecutableHash: executableHash,
		Version:        knobs.declaredVersion,
		ConfigHash:     configHash,
		Disposition:    knobs.capabilityDisposition,
		Matrix:         assemblyCapabilityMatrix(knobs.capabilityDisposition),
		Evidence:       []string{assemblyDeclaredDigest("capability-evidence", "declared, not probed")},
	})
	if err != nil {
		t.Fatalf("declared capability: %v", err)
	}

	var enforcement *adapters.AgentRunnerEnforcementRecord
	if knobs.enforcementMode != "" {
		probedAt := base
		if knobs.enforcementMode == "stale" {
			probedAt = base.Add(-2 * time.Hour)
		}
		record, err := adapters.NewAgentRunnerEnforcementRecord(adapters.AgentRunnerEnforcement{
			EnforcementID:         "t027-enforcement",
			ProviderID:            "t027-provider",
			ProbedAt:              probedAt.UTC().Format(timestampLayout),
			CapabilityRecordHash:  capability.RecordHash,
			ExecutableHash:        executableHash,
			ConfigHash:            configHash,
			PlatformHash:          platformHash,
			EnforcementConfigHash: assemblyDeclaredDigest("enforcement-config", "declared"),
			Disposition:           "verified",
			Matrix: adapters.AgentRunnerEnforcementMatrix{
				ToolControl:    assemblyVerifiedFinding("enforcement-tool"),
				NetworkControl: assemblyVerifiedFinding("enforcement-network"),
			},
			Evidence: []string{assemblyDeclaredDigest("enforcement-evidence", "declared")},
		})
		if err != nil {
			t.Fatalf("declared enforcement: %v", err)
		}
		enforcement = &record
	}

	requestRecord, err := adapters.NewAgentRunnerRequestRecord(adapters.AgentRunnerRequest{
		RequestID:            assemblyRequestID,
		RunID:                assemblyRunID,
		TaskID:               assemblyTaskID,
		StepID:               assemblyStepID,
		Attempt:              1,
		CreatedAt:            base.UTC().Format(timestampLayout),
		AdapterID:            assemblyAdapterID,
		Mode:                 "create",
		WorkspaceHash:        parent,
		ContextHash:          assemblyDeclaredDigest("context", assemblyRunID),
		ParentSnapshotHash:   parent,
		AuthorizationHash:    grant.RecordHash,
		BudgetHash:           reservation.ReservationHash,
		AllowedTargets:       []string{"workspace-files"},
		AllowedEffects:       []string{"local-process", "workspace-write"},
		EffectMappingVersion: assemblyEffectMappingVersion,
		EffectMappingHash:    assemblyEffectMappingHash,
		Evidence:             []string{assemblyDeclaredDigest("request-evidence", assemblyRunID)},
	})
	if err != nil {
		t.Fatalf("declared request record: %v", err)
	}

	return &assemblyFixture{
		runRoot: runRoot,
		options: AgentRunnerAssemblyOptions{
			RunID:               assemblyRunID,
			RunRoot:             runRoot,
			Definition:          assemblyDefinition(),
			RequestRecord:       requestRecord,
			CapabilityRecord:    capability,
			EnforcementRecord:   enforcement,
			AvailabilityRecord:  availability,
			AvailabilityPolicy:  AgentRunnerDefaultAvailabilityPolicy(availability, base),
			AuthorizationLedger: chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{grant}},
			CostLedger:          ledger,
			AgentCLI: AgentRunnerCLIOptions{
				Executable:  knobs.declaredExecutable,
				Version:     knobs.declaredVersion,
				VersionArgs: []string{"-version"},
				RunArgs:     append([]string(nil), knobs.runArgs...),
				Timeout:     30 * time.Second,
				MaxLogBytes: 1 << 20,
				VerifyWait:  1500 * time.Millisecond,
			},
			WorkspaceRoot: knobs.workspaceOverride,
			Clock:         func() time.Time { return time.Now().UTC() },
		},
	}
}

// agentRunnerPublicationDurability probes the platform constant through the
// public store API on a throwaway run root.
func agentRunnerPublicationDurability(t *testing.T) adapters.PublishDurability {
	t.Helper()
	probeRoot := t.TempDir()
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(probeRoot)); err != nil {
		t.Fatal(err)
	}
	store, err := adapters.NewAgentRunnerReplayStoreForRun(probeRoot, "run-durability-probe")
	if err != nil {
		t.Fatal(err)
	}
	return store.PublicationDurability()
}

func agentRunnerDirEntries(t *testing.T, path string) []string {
	t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	return names
}

// requireAgentRunnerZeroSpawn proves the structural half of a refusal: no replay
// identity, no launch slot, no frozen evidence and no dispatch transition.
func requireAgentRunnerZeroSpawn(t *testing.T, runRoot string) {
	t.Helper()
	replayRoot := filepath.Join(runRoot, "agent-runner-replay")
	for _, name := range []string{"requests", "completions", "launches", "terminals", "runs"} {
		if entries := agentRunnerDirEntries(t, filepath.Join(replayRoot, name)); len(entries) > 0 {
			t.Fatalf("a refused assembly must leave %s/ empty, found %v", name, entries)
		}
	}
	if _, err := os.Stat(filepath.Join(runRoot, managedProcessIdentityFileName)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("a refused assembly must leave no managed-process identity mirror (stat err=%v)", err)
	}
	data, err := os.ReadFile(runEventLogPath(runRoot))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "agent-dispatched") || strings.Contains(string(data), "TERMINAL_PENDING") {
		t.Fatalf("a refused assembly must write no dispatch transition: %s", data)
	}
}

func requireAgentRunnerEventReason(t *testing.T, runRoot, reason string) {
	t.Helper()
	data, err := os.ReadFile(runEventLogPath(runRoot))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), reason) {
		t.Fatalf("the event log must record %q: %s", reason, data)
	}
}

// TestAgentRunnerAssemblyFullChain is the positive production leg. It runs only
// where the replay store proves publication durability, because the production
// dispatcher refuses to publish R anywhere else.
func TestAgentRunnerAssemblyFullChain(t *testing.T) {
	if durability := agentRunnerPublicationDurability(t); durability != adapters.PublishDurabilityProven {
		t.Skipf("SKIP: publication durability is %s on %s, so the production dispatcher refuses before R and the full chain cannot run here. The Windows-native evidence is TestAgentRunnerAssemblyRefusesOnUnprovenDurability; the full chain runs on Ubuntu CI (local Unix).", durability, runtime.GOOS)
	}
	stub := buildAssemblyStub(t)
	fixture := newAssemblyFixture(t, assemblyKnobs{
		declaredExecutable: stub,
		runArgs:            []string{"-mode", "mutate-workspace", "-sleep", "100ms"},
	})
	// The pinned offline stub writes stub-child.pid only when it spawns a child, so
	// this workload leaves only stub.pid. The postflight process leg now requires
	// both reports, so seed the empty "reported no children" report the CLI would
	// write for a run with no descendants.
	agentRunnerSeedChildPIDReport(t, fixture.runRoot)
	summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err != nil {
		t.Fatalf("the production assembly must complete: %v", err)
	}
	if summary.ChainState != "COMPLETED" {
		t.Fatalf("chain state = %s, want COMPLETED (%+v)", summary.ChainState, summary)
	}
	requireAgentRunnerEventReason(t, fixture.runRoot, "TERMINAL_PENDING")
	requireAgentRunnerEventReason(t, fixture.runRoot, "postflight-passed")
	passed := false
	for _, count := range summary.TaskStateCounts {
		if count.State == "PASSED" {
			passed = true
		}
	}
	if !passed {
		t.Fatalf("the task must be PASSED: %+v", summary.TaskStateCounts)
	}
	if entries := agentRunnerDirEntries(t, filepath.Join(fixture.runRoot, "agent-runner-replay", "runs")); len(entries) == 0 {
		t.Fatal("the completed run must have frozen its evidence directory")
	}
}

// TestAgentRunnerAssemblyRefusesOnUnprovenDurability is the Windows-native leg:
// the production dispatch refuses at RecordRequest with zero spawn artifacts.
func TestAgentRunnerAssemblyRefusesOnUnprovenDurability(t *testing.T) {
	if durability := agentRunnerPublicationDurability(t); durability == adapters.PublishDurabilityProven {
		t.Skipf("SKIP: this leg measures the unproven-durability refusal; publication durability is proven on %s. The full chain is exercised by TestAgentRunnerAssemblyFullChain.", runtime.GOOS)
	}
	fixture := newAssemblyFixture(t, assemblyKnobs{})
	summary, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil {
		t.Fatalf("an unproven-durability platform must refuse the dispatch: %+v", summary)
	}
	if !errors.Is(err, adapters.ErrAgentRunnerReplayStoreDurabilityUnproven) {
		t.Fatalf("the refusal must be the durability sentinel, got %v", err)
	}
	// The best-effort mitigation must have attempted the reclaim: no launch receipt
	// exists on this refusal path, so it reports the stop outcome as unproven while
	// keeping the durability sentinel on the chain.
	if !strings.Contains(err.Error(), "unproven") {
		t.Fatalf("the refused dispatch must report the unproven stop mitigation, got %v", err)
	}
	requireAgentRunnerZeroSpawn(t, fixture.runRoot)
}

// TestAgentRunnerAssemblyRefusesOnUncompensatedCapability is the no-enforcement
// reverse leg: a non-compatible capability without an enforcement record is
// refused by the frozen admission, before any dispatch.
func TestAgentRunnerAssemblyRefusesOnUncompensatedCapability(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{capabilityDisposition: "blocked"})
	_, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil {
		t.Fatal("a non-compatible capability without enforcement must be refused")
	}
	if !errors.Is(err, adapters.ErrAgentRunnerCompositeAdmissionBlocked) || !errors.Is(err, adapters.ErrAgentRunnerAdmissionBlocked) {
		t.Fatalf("the refusal must be the composite admission sentinel, got %v", err)
	}
	if !strings.Contains(err.Error(), "external enforcement required") {
		t.Fatalf("the refusal must name the missing enforcement, got %v", err)
	}
	requireAgentRunnerZeroSpawn(t, fixture.runRoot)
}

// TestAgentRunnerAssemblyRefusesOnPlatformMismatch is the platform reverse leg:
// a declared capability for another platform is refused even when compatible.
func TestAgentRunnerAssemblyRefusesOnPlatformMismatch(t *testing.T) {
	other := "linux"
	if runtime.GOOS == "linux" {
		other = "windows"
	}
	fixture := newAssemblyFixture(t, assemblyKnobs{capabilityPlatformOS: other})
	_, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil {
		t.Fatal("a capability declared for another platform must be refused")
	}
	if !errors.Is(err, adapters.ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("the refusal must be the composite admission sentinel, got %v", err)
	}
	if !strings.Contains(err.Error(), "capability platform mismatch") {
		t.Fatalf("the refusal must name the platform mismatch, got %v", err)
	}
	requireAgentRunnerZeroSpawn(t, fixture.runRoot)
}

// TestAgentRunnerAssemblyRefusesOnStaleEnforcement is the stale-enforcement
// reverse leg: a bound enforcement record outside the maximum age is refused.
func TestAgentRunnerAssemblyRefusesOnStaleEnforcement(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{capabilityDisposition: "blocked", enforcementMode: "stale"})
	_, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil {
		t.Fatal("a stale enforcement record must be refused")
	}
	if !errors.Is(err, adapters.ErrAgentRunnerCompositeAdmissionBlocked) {
		t.Fatalf("the refusal must be the composite admission sentinel, got %v", err)
	}
	if !strings.Contains(err.Error(), "stale or future enforcement evidence") {
		t.Fatalf("the refusal must name the staleness, got %v", err)
	}
	requireAgentRunnerZeroSpawn(t, fixture.runRoot)
}

// TestAgentRunnerAssemblyRefusesOnUnavailableChannel is the no-availability
// reverse leg: an unavailable channel is refused before capability admission.
func TestAgentRunnerAssemblyRefusesOnUnavailableChannel(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{availabilityStatus: "unavailable"})
	_, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil {
		t.Fatal("an unavailable channel must be refused")
	}
	if !errors.Is(err, adapters.ErrAgentRunnerCompositeAdmissionBlocked) || !errors.Is(err, adapters.ErrAIUnavailable) {
		t.Fatalf("the refusal must carry the availability sentinel, got %v", err)
	}
	requireAgentRunnerZeroSpawn(t, fixture.runRoot)
}

// TestAgentRunnerAssemblyRefusesOnWorkspaceReplayOverlap is the workspace/replay
// overlap reverse leg: an isolated workspace inside the replay root is refused by
// the replay store before the engine exists.
func TestAgentRunnerAssemblyRefusesOnWorkspaceReplayOverlap(t *testing.T) {
	// The override is only known once the run root exists, so the fixture is built
	// first and the overlap is injected through the run root layout.
	fixture := newAssemblyFixture(t, assemblyKnobs{workspaceOverride: ""})
	overlap := filepath.Join(fixture.runRoot, "agent-runner-replay", "workspaces")
	fixture.options.WorkspaceRoot = overlap
	_, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil {
		t.Fatal("a workspace inside the replay root must be refused")
	}
	if !errors.Is(err, adapters.ErrAgentRunnerReplayRootConflict) {
		t.Fatalf("the refusal must be the replay-root conflict sentinel, got %v", err)
	}
	requireAgentRunnerZeroSpawn(t, fixture.runRoot)
}

// TestAgentRunnerAssemblyRefusesOnIncompleteRecords is the assembly-precondition
// reverse leg: a request record that does not bind the definition step is
// refused before the engine exists.
func TestAgentRunnerAssemblyRefusesOnIncompleteRecords(t *testing.T) {
	fixture := newAssemblyFixture(t, assemblyKnobs{})
	fixture.options.RequestRecord = adapters.AgentRunnerRequestRecord{}
	_, err := ExecuteAgentRunnerRun(context.Background(), fixture.options)
	if err == nil {
		t.Fatal("a missing request record must be refused")
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("the precondition refusal must be the console config sentinel, got %v", err)
	}
	if _, statErr := os.Stat(runEventLogPath(fixture.runRoot)); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("a precondition refusal must not even create the run-root marker (stat err=%v)", statErr)
	}
}

// newAgentRunnerPIDReportRuntime builds a production-bound runtime over a real
// replay store and returns the evidence directory the pid reports live in. It
// exercises exactly the path Publish uses: the store is derived from the run
// root through adapters.NewAgentRunnerReplayStoreForRun and the evidence
// directory through RunEvidenceDir.
func newAgentRunnerPIDReportRuntime(t *testing.T) (*agentRunnerRuntime, string) {
	t.Helper()
	runRoot := t.TempDir()
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(runRoot)); err != nil {
		t.Fatal(err)
	}
	store, err := adapters.NewAgentRunnerReplayStoreForRun(runRoot, assemblyRunID)
	if err != nil {
		t.Fatal(err)
	}
	evidenceDir, err := store.RunEvidenceDir(assemblyRequestID)
	if err != nil {
		t.Fatal(err)
	}
	return &agentRunnerRuntime{
		runID: assemblyRunID,
		store: store,
		requestRecord: adapters.AgentRunnerRequestRecord{
			Request: adapters.AgentRunnerRequest{RequestID: assemblyRequestID},
		},
	}, evidenceDir
}

// agentRunnerSeedChildPIDReport writes the "reported no children" pid report the
// pinned offline stub does not write itself (it writes stub-child.pid only when it
// spawns a child). It exists so the full-chain fixture can satisfy the postflight
// process leg, which requires both reports.
func agentRunnerSeedChildPIDReport(t *testing.T, runRoot string) {
	t.Helper()
	if err := ensureAgentRunnerEventLogMarker(runEventLogPath(runRoot)); err != nil {
		t.Fatal(err)
	}
	store, err := adapters.NewAgentRunnerReplayStoreForRun(runRoot, assemblyRunID)
	if err != nil {
		t.Fatal(err)
	}
	evidenceDir, err := store.RunEvidenceDir(assemblyRequestID)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(evidenceDir, agentRunnerChildPIDReportFileName), nil, 0o600); err != nil {
		t.Fatal(err)
	}
}

// agentRunnerStubReceipts answers the launch-receipt lookup from a script.
type agentRunnerStubReceipts struct {
	receipt adapters.AgentRunnerLaunchReceiptRecord
	found   bool
	err     error
}

func (reader agentRunnerStubReceipts) LaunchReceipt(string) (adapters.AgentRunnerLaunchReceiptRecord, bool, error) {
	return reader.receipt, reader.found, reader.err
}

// agentRunnerRecordingStopper records the stop calls and answers with a scripted
// proof, so the mitigation can be driven without a live process.
type agentRunnerRecordingStopper struct {
	calls []string
	proof guard.TerminationEvidence
	err   error
}

func (stopper *agentRunnerRecordingStopper) StopAgentRunnerProcess(_ context.Context, launchID string) (guard.TerminationEvidence, error) {
	stopper.calls = append(stopper.calls, launchID)
	return stopper.proof, stopper.err
}

// TestAgentRunnerMitigateAbandonedLaunchStopsAndKeepsTheOriginalError pins the
// best-effort mitigation: a recorded launch is stopped, the original engine error
// always stays on the chain, and every way the stop can fail to be proven is
// reported as unproven instead of passing silently.
func TestAgentRunnerMitigateAbandonedLaunchStopsAndKeepsTheOriginalError(t *testing.T) {
	original := errors.New("the engine refused the dispatch")
	receipt := adapters.AgentRunnerLaunchReceiptRecord{Receipt: adapters.AgentRunnerLaunchReceipt{LaunchID: "launch-" + assemblyRequestID}}

	t.Run("a proven stop keeps the original error", func(t *testing.T) {
		stopper := &agentRunnerRecordingStopper{proof: guard.TerminationEvidence{Outcome: "stopped"}}
		err := agentRunnerMitigateAbandonedLaunch(context.Background(), agentRunnerStubReceipts{receipt: receipt, found: true}, stopper, assemblyRequestID, original)
		if !errors.Is(err, original) {
			t.Fatalf("the original error must stay on the chain: %v", err)
		}
		if len(stopper.calls) != 1 || stopper.calls[0] != receipt.Receipt.LaunchID {
			t.Fatalf("the mitigation must stop the recorded launch: %v", stopper.calls)
		}
		if strings.Contains(err.Error(), "unproven") {
			t.Fatalf("a proven stop must not report unproven: %v", err)
		}
	})

	t.Run("a missing receipt is unproven", func(t *testing.T) {
		stopper := &agentRunnerRecordingStopper{}
		err := agentRunnerMitigateAbandonedLaunch(context.Background(), agentRunnerStubReceipts{}, stopper, assemblyRequestID, original)
		if !errors.Is(err, original) {
			t.Fatalf("the original error must stay on the chain: %v", err)
		}
		if len(stopper.calls) != 0 {
			t.Fatalf("no receipt must mean no stop call: %v", stopper.calls)
		}
		if !strings.Contains(err.Error(), "unproven") {
			t.Fatalf("a missing receipt must report unproven: %v", err)
		}
	})

	t.Run("a stop error is unproven", func(t *testing.T) {
		stopper := &agentRunnerRecordingStopper{err: errors.New("the stop path refused")}
		err := agentRunnerMitigateAbandonedLaunch(context.Background(), agentRunnerStubReceipts{receipt: receipt, found: true}, stopper, assemblyRequestID, original)
		if !errors.Is(err, original) {
			t.Fatalf("the original error must stay on the chain: %v", err)
		}
		if len(stopper.calls) != 1 {
			t.Fatalf("a present receipt must be stopped: %v", stopper.calls)
		}
		if !strings.Contains(err.Error(), "unproven") {
			t.Fatalf("a stop error must report unproven: %v", err)
		}
	})

	t.Run("an unproven stop proof is unproven", func(t *testing.T) {
		stopper := &agentRunnerRecordingStopper{proof: guard.TerminationEvidence{Outcome: "uncertain"}}
		err := agentRunnerMitigateAbandonedLaunch(context.Background(), agentRunnerStubReceipts{receipt: receipt, found: true}, stopper, assemblyRequestID, original)
		if !errors.Is(err, original) {
			t.Fatalf("the original error must stay on the chain: %v", err)
		}
		if !strings.Contains(err.Error(), "unproven") {
			t.Fatalf("an unproven stop proof must report unproven: %v", err)
		}
	})
}

// TestAgentRunnerSeedChildPIDReportLandsInTheEvidenceDir proves the full-chain
// fixture helper writes the "reported no children" report into the evidence
// directory the store resolves, and that a second store construction over the
// same run root stays safe. Without it the full-chain leg (Linux only) would be
// the sole coverage of the seeding mechanics.
func TestAgentRunnerSeedChildPIDReportLandsInTheEvidenceDir(t *testing.T) {
	runRoot := t.TempDir()
	agentRunnerSeedChildPIDReport(t, runRoot)
	store, err := adapters.NewAgentRunnerReplayStoreForRun(runRoot, assemblyRunID)
	if err != nil {
		t.Fatal(err)
	}
	evidenceDir, err := store.RunEvidenceDir(assemblyRequestID)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(evidenceDir, agentRunnerChildPIDReportFileName))
	if err != nil {
		t.Fatalf("the seeded child report must exist: %v", err)
	}
	if len(data) != 0 {
		t.Fatalf("the seeded child report must be empty, got %q", string(data))
	}
}

// TestAgentRunnerAssemblyRefusesPublishWithoutPIDReport pins the fail-closed
// doctrine of the production Publish process proof. It is the same doctrine, and
// the same verdict, as leg 5 of the adapters postflight (agentRunnerAlivePIDs):
// BOTH pid reports must exist, and any missing report or unreadable entry blocks
// the promotion. A run that left either report missing cannot be proven to have
// stopped its managed processes, so the assembly refuses to publish.
func TestAgentRunnerAssemblyRefusesPublishWithoutPIDReport(t *testing.T) {
	// A dead but genuine pid: its report proves "reported", while the pid itself is
	// provably gone. This mirrors the dead-pid technique of
	// internal/adapters/agent_runner_run_test.go.
	dead := exec.Command(os.Args[0], "-test.run=^$")
	if err := dead.Run(); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name        string
		parent      string // "", "live", "dead", "garbage"
		child       string
		wantStopped bool
		wantInvalid bool
	}{
		{
			name:        "both pid reports deleted is a fail-closed refusal",
			wantInvalid: true,
		},
		{
			name:        "only a parent report is a fail-closed refusal",
			parent:      "dead",
			wantInvalid: true,
		},
		{
			// The ③ row of the previous table: a single child report can no longer
			// pass, because the missing parent report is itself a fail-closed gap.
			name:        "only a child report is a fail-closed refusal",
			child:       "dead",
			wantInvalid: true,
		},
		{
			name:        "both reports with a live parent keep the assembly blocked",
			parent:      "live",
			child:       "dead",
			wantStopped: false,
		},
		{
			name:        "both reports with gone pids prove stopped",
			parent:      "dead",
			child:       "dead",
			wantStopped: true,
		},
		{
			name:        "an unparsable pid report is a fail-closed refusal",
			parent:      "garbage",
			child:       "dead",
			wantInvalid: true,
		},
		{
			// The console parser used to read 12 from "12abc" via fmt.Sscanf, while the
			// adapters parser refused it: a drift that let the two sides answer different
			// verdicts for the same evidence. This row keeps the integration layer red if
			// the console side ever regresses to a numeric-prefix parser.
			name:        "a parent report with trailing garbage is a fail-closed refusal",
			parent:      "trailing",
			child:       "dead",
			wantInvalid: true,
		},
		{
			name:        "a parent report with a float pid is a fail-closed refusal",
			parent:      "float",
			child:       "dead",
			wantInvalid: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runtime, evidenceDir := newAgentRunnerPIDReportRuntime(t)
			writeReport := func(name, kind string) {
				t.Helper()
				if kind == "" {
					return
				}
				value := ""
				switch kind {
				case "live":
					value = fmt.Sprintf("%d\n", os.Getpid())
				case "dead":
					value = fmt.Sprintf("%d\n", dead.Process.Pid)
				case "garbage":
					value = "not-a-pid\n"
				case "trailing":
					value = "12abc\n"
				case "float":
					value = "7.0\n"
				}
				if err := os.WriteFile(filepath.Join(evidenceDir, name), []byte(value), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			writeReport(agentRunnerParentPIDReportFileName, tc.parent)
			writeReport(agentRunnerChildPIDReportFileName, tc.child)
			stopped, err := runtime.managedProcessesGone()
			if tc.wantInvalid {
				if stopped || !errors.Is(err, adapters.ErrInvalidAgentRunnerEvidence) {
					t.Fatalf("a run missing or unable to read a pid report must fail closed: stopped=%v err=%v", stopped, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("a run that left both pid reports must not error: %v", err)
			}
			if stopped != tc.wantStopped {
				t.Fatalf("stopped=%v, want %v", stopped, tc.wantStopped)
			}
		})
	}
}

// TestAgentRunnerReadPIDsMatrixConsoleTwin pins the exact pid-parsing contract of
// readAgentRunnerEvidencePIDs.
//
// COUPLING: this is a hand-synchronised twin of TestAgentRunnerReadPIDsMatrix in
// internal/adapters/agent_runner_postflight_test.go. Both sides read the same
// report format and must answer the same verdict, so the matrix and its
// expectations are kept equal by hand, exactly like agentRunnerMissingPIDReportDetail
// and the pid report file names. The parse is pure (t.TempDir only), so the case
// runs identically on Windows and Linux.
func TestAgentRunnerReadPIDsMatrixConsoleTwin(t *testing.T) {
	cases := []struct {
		name        string
		content     string
		missing     bool
		wantPIDs    []int
		wantPresent bool
		wantInvalid bool
	}{
		{name: "plain 7", content: "7", wantPIDs: []int{7}, wantPresent: true},
		{name: "signed plus 7", content: "+7", wantPIDs: []int{7}, wantPresent: true},
		{name: "padded 7", content: " 7 ", wantPIDs: []int{7}, wantPresent: true},
		{name: "zero padded 007", content: "007", wantPIDs: []int{7}, wantPresent: true},
		{name: "CRLF terminated 7", content: "7\r\n", wantPIDs: []int{7}, wantPresent: true},
		{name: "two pids with an interior blank line", content: "7\n\n8", wantPIDs: []int{7, 8}, wantPresent: true},
		{name: "trailing garbage 12abc", content: "12abc", wantInvalid: true},
		{name: "float 7.0", content: "7.0", wantInvalid: true},
		{name: "underscore 1_0", content: "1_0", wantInvalid: true},
		{name: "hex 0x10", content: "0x10", wantInvalid: true},
		{name: "negative zero -0", content: "-0", wantInvalid: true},
		{name: "negative -5", content: "-5", wantInvalid: true},
		{name: "int64 overflow", content: "9223372036854775808", wantInvalid: true},
		{name: "missing file", missing: true},
		{name: "empty file", content: "", wantPIDs: []int{}, wantPresent: true},
		{name: "whitespace only", content: "   \n\t\n  \n", wantPIDs: []int{}, wantPresent: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "stub.pid")
			if !tc.missing {
				if err := os.WriteFile(path, []byte(tc.content), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			pids, present, err := readAgentRunnerEvidencePIDs(path)
			if tc.missing {
				if err != nil || present || pids != nil {
					t.Fatalf("a missing report must answer present=false without error: pids=%v present=%v err=%v", pids, present, err)
				}
				return
			}
			if tc.wantInvalid {
				if !errors.Is(err, adapters.ErrInvalidAgentRunnerEvidence) || present {
					t.Fatalf("an unparsable report must fail closed: pids=%v present=%v err=%v", pids, present, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("a parsable report must not error: %v", err)
			}
			if present != tc.wantPresent || !agentRunnerSamePIDs(pids, tc.wantPIDs) {
				t.Fatalf("pids=%v present=%v, want pids=%v present=%v", pids, present, tc.wantPIDs, tc.wantPresent)
			}
		})
	}
}

// agentRunnerSamePIDs compares pid slices element by element, so an empty result
// compares equal whether it is nil or a zero-length slice.
func agentRunnerSamePIDs(got, want []int) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

// TestAgentRunnerMitigateAbandonedLaunchClosureCoversEveryUnprovenShape pins the
// two unproven shapes the fail closure in ExecuteAgentRunnerRun relies on but
// TestAgentRunnerMitigateAbandonedLaunchStopsAndKeepsTheOriginalError did not
// exercise: an unavailable reclaimer (nil reader or stopper) and an unreadable
// receipt. Together the two tests cover every branch of the single helper the
// closure delegates to, without duplicating the already pinned shapes.
func TestAgentRunnerMitigateAbandonedLaunchClosureCoversEveryUnprovenShape(t *testing.T) {
	original := errors.New("the engine refused after the spawn")

	t.Run("an unavailable reclaimer is unproven", func(t *testing.T) {
		err := agentRunnerMitigateAbandonedLaunch(context.Background(), nil, nil, assemblyRequestID, original)
		if !errors.Is(err, original) {
			t.Fatalf("the original error must stay on the chain: %v", err)
		}
		if !strings.Contains(err.Error(), "unproven") {
			t.Fatalf("an unavailable reclaimer must report unproven: %v", err)
		}
	})

	t.Run("an unreadable receipt is unproven", func(t *testing.T) {
		stopper := &agentRunnerRecordingStopper{}
		err := agentRunnerMitigateAbandonedLaunch(context.Background(), agentRunnerStubReceipts{err: errors.New("the receipt is unreadable")}, stopper, assemblyRequestID, original)
		if !errors.Is(err, original) {
			t.Fatalf("the original error must stay on the chain: %v", err)
		}
		if len(stopper.calls) != 0 {
			t.Fatalf("an unreadable receipt must mean no stop call: %v", stopper.calls)
		}
		if !strings.Contains(err.Error(), "unproven") {
			t.Fatalf("an unreadable receipt must report unproven: %v", err)
		}
	})
}

// TestAgentRunnerAssemblyFunnelsEveryPostSpawnFailureThroughTheFallback is a
// source-level invariant test. Windows cannot reach the awaiting-and-after paths
// of ExecuteAgentRunnerRun at runtime, so the compiler-visible shape stands in
// for the missing runtime evidence: the fail closure is the single exit for any
// failure after the engine exists, and no post-spawn failure may bypass the
// abandoned-launch mitigation the closure applies.
func TestAgentRunnerAssemblyFunnelsEveryPostSpawnFailureThroughTheFallback(t *testing.T) {
	// wantFailureReturnPoints is the guard rail. It is 7: the first engine.Run
	// non-awaiting error, the WaitAgentRunnerDispatched error, the PublishTerminal
	// error, the LoadChainAgentRunnerTerminal error, the SubmitAgentRunnerTerminal
	// error, the non-completed status error and the second engine.Run error. A new
	// failure path must return fail(...) and bump this constant in the same change.
	const wantFailureReturnPoints = 7
	const nakedFailureReturn = "summarizeProjection(engine.Projection(), eventLog), err"
	const successReturn = "return summarizeProjection(engine.Projection(), eventLog), nil"

	source, err := os.ReadFile("runtime_agent.go")
	if err != nil {
		t.Fatalf("the invariant reads runtime_agent.go from the package directory: %v", err)
	}
	body, ok := agentRunnerSourceFunctionBody(string(source), "func ExecuteAgentRunnerRun(")
	if !ok {
		t.Fatal("ExecuteAgentRunnerRun must be defined in runtime_agent.go")
	}

	// (a) The naked failure return appears exactly once, inside the fail closure.
	if got := strings.Count(body, nakedFailureReturn); got != 1 {
		t.Fatalf("the naked failure return must appear exactly once inside the fail closure, found %d", got)
	}
	closureStart := strings.Index(body, "fail := func(")
	if closureStart < 0 {
		t.Fatal("ExecuteAgentRunnerRun must define the fail closure")
	}
	closureOpen := strings.IndexByte(body[closureStart:], '{')
	if closureOpen < 0 {
		t.Fatal("the fail closure must have a body")
	}
	closureOpen += closureStart
	closureClose := agentRunnerSourceMatchingBrace(body, closureOpen)
	if closureClose < 0 {
		t.Fatal("the fail closure body must be balanced")
	}
	if !strings.Contains(body[closureStart:closureClose+1], nakedFailureReturn) {
		t.Fatal("the single failure return must live inside the fail closure")
	}

	// (b) Every failure return goes through the closure.
	if got := strings.Count(body, "return fail("); got != wantFailureReturnPoints {
		t.Fatalf("the fail closure must be the single entry for every failure return: found %d return fail(...) call sites, want %d (bump wantFailureReturnPoints when a failure path is added)", got, wantFailureReturnPoints)
	}

	// (c) No projection-based return bypasses the closure: every
	// return summarizeProjection(...) is either the one naked failure return or a
	// direct success return.
	offset := 0
	for {
		index := strings.Index(body[offset:], "return summarizeProjection(")
		if index < 0 {
			break
		}
		absolute := offset + index
		end := strings.IndexByte(body[absolute:], '\n')
		if end < 0 {
			end = len(body) - absolute
		}
		statement := strings.TrimSpace(body[absolute : absolute+end])
		if statement != "return "+nakedFailureReturn && statement != successReturn {
			t.Fatalf("a projection failure return bypasses the fail closure: %q", statement)
		}
		offset = absolute + len("return summarizeProjection(")
	}
}

// agentRunnerSourceFunctionBody returns the body of the named function in src,
// braces included, using a scanner that skips comments and string/rune literals
// so a brace inside any of them is not counted.
func agentRunnerSourceFunctionBody(src, signature string) (string, bool) {
	start := strings.Index(src, signature)
	if start < 0 {
		return "", false
	}
	open := strings.IndexByte(src[start:], '{')
	if open < 0 {
		return "", false
	}
	open += start
	close := agentRunnerSourceMatchingBrace(src, open)
	if close < 0 {
		return "", false
	}
	return src[open : close+1], true
}

// agentRunnerSourceMatchingBrace returns the index of the brace balancing the
// '{' at open, skipping line comments, block comments, string literals and rune
// literals.
func agentRunnerSourceMatchingBrace(src string, open int) int {
	depth := 0
	for i := open; i < len(src); i++ {
		switch src[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		case '"':
			i++
			for i < len(src) && src[i] != '"' {
				if src[i] == '\\' {
					i++
				}
				i++
			}
		case '`':
			i++
			for i < len(src) && src[i] != '`' {
				i++
			}
		case '\'':
			i++
			for i < len(src) && src[i] != '\'' {
				if src[i] == '\\' {
					i++
				}
				i++
			}
		case '/':
			if i+1 < len(src) && src[i+1] == '/' {
				for i < len(src) && src[i] != '\n' {
					i++
				}
			} else if i+1 < len(src) && src[i+1] == '*' {
				i += 2
				for i+1 < len(src) && !(src[i] == '*' && src[i+1] == '/') {
					i++
				}
				i++
			}
		}
	}
	return -1
}
