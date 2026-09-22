// Command b4-e2e is the B4 slice assembly harness: it drives the real ProofRail
// chain engine against the pinned deterministic CLI (tools/agent-stub) through
// the production adapter ports, on the real filesystem, with every managed
// process spawned by internal/guard.
//
// The harness is deliberately a *thin* assembly: every port it hands to
// chain.Options is either a production implementation (internal/adapters,
// internal/gates, internal/snapshot, internal/guard, internal/evidence) or the
// smallest possible real implementation, and the file runmap.go carries the one
// piece of production logic this harness has to reproduce (the run outcome
// mapping, which lives in unexported functions).
//
// Declared inputs note: admission (AgentRunnerCompositeAdmission) requires
// persisted preflight facts. This harness DECLARES them for the pinned offline
// CLI; it does not probe a candidate. See fixtures.go for the exact field values
// and declared-inputs.json in the evidence pack.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
	"github.com/larsonzh/prfrail/internal/tickets"
)

// The two effect-mapping constants are mirrored from unexported constants in
// internal/adapters/agent_runner.go (agentRunnerEffectMappingVersion /
// agentRunnerEffectMappingHash). The request record is refused without them, and
// they are contract values rather than tunables, so a mirror is the only
// external-API way to assemble a legal request. Any divergence is a defect: the
// harness records the mirrored pair in the verdict so a reviewer can diff it.
const (
	mirroredEffectMappingVersion = "1"
	mirroredEffectMappingHash    = "sha256:2c9609ee374bfc79bd0ed997aea6cfe4865c82ed0216d1a4a288139e0b0c92d4"
)

// Declared identities of the B4 workload. The adapter id and availability
// profile id are named "declared" on purpose: nothing here is a probe result of
// an AI candidate, and the B4 report states that explicitly.
const (
	declaredAdapterID      = "b4-pinned-cli-declared"
	declaredProfileID      = "b4-pinned-cli-declared"
	declaredGrantID        = "b4-authorization-1"
	declaredProviderID     = "b4-guard-job-object"
	declaredCapabilityNote = "declared for the pinned offline CLI; not an AI candidate probe"
)

// declaredInputs is everything admission and terminal settlement need, built
// from measured host/executable facts plus the declared capability statement.
type declaredInputs struct {
	record      adapters.AgentRunnerRequestRecord
	admission   adapters.AgentRunnerCompositeAdmission
	ledger      *tickets.CostLedger
	grantHash   string
	enfRecord   *adapters.AgentRunnerEnforcementRecord
	enfDeclared bool
}

// requestSpec pins the identity one harness invocation assembles.
type requestSpec struct {
	RunID      string
	TaskID     string
	StepID     string
	CreatedAt  time.Time
	ParentHash string
	// WorkspaceHash and ContextHash are digests of the input material.
	WorkspaceHash string
	ContextHash   string
}

// buildRequestRecord assembles and validates the AgentRunner request record. The
// authorization and budget digests are record hashes of records built from
// measured inputs, so the record can only be assembled once both exist.
func buildRequestRecord(spec requestSpec, authorizationHash, budgetHash string) (adapters.AgentRunnerRequestRecord, error) {
	return adapters.NewAgentRunnerRequestRecord(adapters.AgentRunnerRequest{
		RequestID:            "request-" + spec.RunID,
		RunID:                spec.RunID,
		TaskID:               spec.TaskID,
		StepID:               spec.StepID,
		Attempt:              1,
		CreatedAt:            spec.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z"),
		AdapterID:            declaredAdapterID,
		Mode:                 "create",
		WorkspaceHash:        spec.WorkspaceHash,
		ContextHash:          spec.ContextHash,
		ParentSnapshotHash:   spec.ParentHash,
		AuthorizationHash:    authorizationHash,
		BudgetHash:           budgetHash,
		AllowedTargets:       []string{"workspace-files"},
		AllowedEffects:       []string{"local-process", "workspace-write"},
		EffectMappingVersion: mirroredEffectMappingVersion,
		EffectMappingHash:    mirroredEffectMappingHash,
		Evidence:             []string{digestOf("proofrail:b4-declared-request-evidence:1\n", spec.RunID+"\n"+spec.TaskID+"\n"+spec.StepID)},
	})
}

// buildDeclaredInputs builds the authorization grant, the cost reservation, the
// availability record, the declared capability statement and the enforcement
// statement, then binds them into one composite admission.
//
// Capability disposition is "compatible" because that is the only legal shape a
// capability record can take when every matrix finding is verified
// (internal/adapters/agent_runner_capability.go): a "blocked" disposition with a
// fully verified matrix is rejected as inconsistent, and a non-verified finding
// blocks admission outright. The report states the consequence: the matrix is a
// DECLARATION for the pinned offline CLI, and B4 claims the ProofRail mechanism
// leg only - never AI-candidate compatibility.
func buildDeclaredInputs(spec requestSpec, stubPath, stubVersion string, now time.Time) (*declaredInputs, error) {
	executableHash, err := fileHash(stubPath)
	if err != nil {
		return nil, err
	}
	platformHash := digestOf("proofrail:b4-declared-platform:1\n", runtime.GOOS+"\n"+runtime.GOARCH)
	policyHash := digestOf("proofrail:b4-declared-policy:1\n", "b4-admission-policy")
	configHash := digestOf("proofrail:b4-declared-config:1\n", declaredAdapterID+"\n"+stubVersion)

	amountMicros := int64(1000)
	grant, err := chain.NewGrantAuthorizationRecord(chain.AuthorizationGrant{
		RecordID:        declaredGrantID,
		Kind:            "grant",
		IssuedAt:        now.Add(-time.Minute).UTC().Format("2006-01-02T15:04:05.000Z"),
		IssuedBy:        evidence.Actor{Type: "operator", ID: "b4-declared-operator"},
		AuthorizationID: declaredGrantID,
		RunID:           spec.RunID,
		RunManifestHash: spec.ParentHash,
		PolicyHash:      policyHash,
		Scope: chain.AuthorizationScope{
			TaskIDs:       []string{spec.TaskID},
			StepIDs:       []string{spec.StepID},
			TargetIDs:     []string{"workspace-files"},
			EffectClasses: []string{"local-discardable", "read-only"},
			Budget: chain.AuthorizationBudgetLimit{
				ModelCalls:   2,
				Tokens:       1000,
				WallClockMs:  600000,
				Attempts:     1,
				Currency:     "USD",
				AmountMicros: &amountMicros,
			},
		},
		ExpiresAt: now.Add(time.Hour).UTC().Format("2006-01-02T15:04:05.000Z"),
		Evidence:  []string{digestOf("proofrail:b4-declared-grant-evidence:1\n", declaredGrantID)},
	}, func() time.Time { return now }, func() string { return declaredGrantID + "-record" })
	if err != nil {
		return nil, fmt.Errorf("build declared grant: %w", err)
	}

	ledger, err := tickets.NewCostLedger(tickets.CostLedgerConfig{
		LedgerID:       "b4-cost-ledger",
		CreatedAt:      now.UTC(),
		ScopeKind:      tickets.CostScopeRun,
		ScopeHash:      spec.ParentHash,
		Currency:       "USD",
		PricingMode:    tickets.CostPricingDocumented,
		PricingVersion: now.UTC().Format("2006-01-02"),
	})
	if err != nil {
		return nil, fmt.Errorf("build declared cost ledger: %w", err)
	}
	if _, err := ledger.AppendAllocation(tickets.AllocationInput{
		EntryID:           "b4-allocation-1",
		OccurredAt:        now.UTC(),
		AuthorizationHash: grant.RecordHash,
		Limits: tickets.CostLimits{
			AmountMicros: &amountMicros,
			ModelCalls:   2,
			Tokens:       1000,
			WallClockMs:  600000,
			Attempts:     1,
		},
	}); err != nil {
		return nil, fmt.Errorf("append declared allocation: %w", err)
	}
	reservation, err := ledger.Reserve(tickets.ReservationInput{
		EntryID:              "b4-reservation-1",
		OccurredAt:           now.UTC(),
		RunID:                spec.RunID,
		RequestID:            "request-" + spec.RunID,
		IdempotencyKey:       "b4-reserve-1",
		AuthorizationHash:    grant.RecordHash,
		ReservedAmountMicros: &amountMicros,
		ReservedCalls:        2,
		ReservedTokens:       1000,
		PricingEvidence:      []string{digestOf("proofrail:b4-declared-pricing:1\n", "zero-cost-offline-workload")},
	})
	if err != nil {
		return nil, fmt.Errorf("reserve declared budget: %w", err)
	}

	availability, err := adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
		ProbeID:           "b4-declared-availability",
		ProfileID:         declaredProfileID,
		ProfileConfigHash: configHash,
		Channel:           "agent-runner-cli",
		ProbedAt:          now.Add(-time.Minute).UTC().Format("2006-01-02T15:04:05.000Z"),
		Status:            "available",
		RequestsUsed:      1,
		Evidence:          []string{digestOf("proofrail:b4-declared-availability-evidence:1\n", declaredCapabilityNote)},
	})
	if err != nil {
		return nil, fmt.Errorf("build declared availability: %w", err)
	}

	capability, err := adapters.NewAgentRunnerCapabilityRecord(adapters.AgentRunnerCapability{
		ProbeID:        "b4-declared-capability",
		AdapterID:      declaredAdapterID,
		ProbedAt:       now.Add(-time.Minute).UTC().Format("2006-01-02T15:04:05.000Z"),
		OS:             runtime.GOOS,
		Arch:           runtime.GOARCH,
		PlatformHash:   platformHash,
		ExecutableHash: executableHash,
		Version:        stubVersion,
		ConfigHash:     configHash,
		Disposition:    "compatible",
		Matrix:         declaredCapabilityMatrix(),
		Evidence:       []string{digestOf("proofrail:b4-declared-capability-evidence:1\n", declaredCapabilityNote)},
	})
	if err != nil {
		return nil, fmt.Errorf("build declared capability: %w", err)
	}

	enforcement, err := adapters.NewAgentRunnerEnforcementRecord(adapters.AgentRunnerEnforcement{
		EnforcementID:         "b4-declared-enforcement",
		ProviderID:            declaredProviderID,
		ProbedAt:              now.UTC().Format("2006-01-02T15:04:05.000Z"),
		CapabilityRecordHash:  capability.RecordHash,
		ExecutableHash:        executableHash,
		ConfigHash:            configHash,
		PlatformHash:          platformHash,
		EnforcementConfigHash: digestOf("proofrail:b4-declared-enforcement-config:1\n", "guard-managed-process-job-containment"),
		Disposition:           "verified",
		Matrix: adapters.AgentRunnerEnforcementMatrix{
			ToolControl:    verifiedFinding("b4-declared-enforcement-tool"),
			NetworkControl: verifiedFinding("b4-declared-enforcement-network"),
		},
		Evidence: []string{digestOf("proofrail:b4-declared-enforcement-evidence:1\n", "internal/guard job containment; A6 native experiments")},
	})
	if err != nil {
		return nil, fmt.Errorf("build declared enforcement: %w", err)
	}

	record, err := buildRequestRecord(spec, grant.RecordHash, reservation.ReservationHash)
	if err != nil {
		return nil, fmt.Errorf("build request record: %w", err)
	}

	admission := adapters.AgentRunnerCompositeAdmission{
		RequestRecord:       record,
		AuthorizationLedger: chain.AuthorizationLedger{Grants: []chain.AuthorizationRecord{grant}},
		CostLedger:          ledger,
		AvailabilityRecord:  availability,
		AvailabilityPolicy: adapters.AIAvailabilityPolicy{
			ProfileID:         declaredProfileID,
			ProfileConfigHash: configHash,
			Channel:           "agent-runner-cli",
			EvaluatedAt:       now.UTC(),
			MaximumAge:        time.Hour,
			MaximumRequests:   1,
		},
		CapabilityRecord:  capability,
		EnforcementRecord: &enforcement,
		AdmissionPolicy: adapters.AgentRunnerAdmissionPolicy{
			OS:           runtime.GOOS,
			Arch:         runtime.GOARCH,
			PlatformHash: platformHash,
			EvaluatedAt:  now.UTC(),
			MaximumAge:   time.Hour,
		},
		Clock: func() time.Time { return now },
	}
	return &declaredInputs{
		record:      record,
		admission:   admission,
		ledger:      ledger,
		grantHash:   grant.RecordHash,
		enfRecord:   &enforcement,
		enfDeclared: true,
	}, nil
}

// declaredCapabilityMatrix states, per capability dimension, why the pinned
// offline CLI satisfies it. Every finding must be "verified" for a legal
// disposition, so this table is the complete declaration and nothing is implied.
func declaredCapabilityMatrix() adapters.AgentRunnerCapabilityMatrix {
	return adapters.AgentRunnerCapabilityMatrix{
		Noninteractive:            verifiedFinding("stub never reads stdin and never prompts"),
		CWD:                       verifiedFinding("stub resolves its working directory with os.Getwd"),
		EventStreamOrCompleteLogs: verifiedFinding("stub writes events.jsonl and terminates both log streams with the completeness marker"),
		SessionCreate:             verifiedFinding("one requestId is one session; the adapter derives session-<requestId>"),
		SessionResume:             verifiedFinding("declared: a re-run is reconstructible from the durable envelope; the stub itself has no resume surface"),
		Cancellation:              verifiedFinding("stub dies on the reviewed platform stop"),
		ProcessTreeStop:           verifiedFinding("guard job containment removes descendants; A6 native experiments"),
		ToolControl:               verifiedFinding("declared: the stub exposes no tool surface, so no tool can escape"),
		NetworkControl:            verifiedFinding("declared: the stub makes no network calls; it is built from source with no network code"),
		PermissionControl:         verifiedFinding("declared: the stub requests no elevated permission"),
		Usage:                     verifiedFinding("stub writes the closed-shape usage.json"),
		UnattendedConfirmations:   verifiedFinding("stub never requires a confirmation"),
	}
}

func verifiedFinding(note string) adapters.AgentRunnerCapabilityFinding {
	return adapters.AgentRunnerCapabilityFinding{
		Status:   "verified",
		Evidence: []string{digestOf("proofrail:b4-declared-finding:1\n", note)},
	}
}

// fileHash returns the content digest of a file.
func fileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

// digestOf mirrors evidence.Digest for harness-local domain strings.
func digestOf(domain, value string) string {
	return evidence.Digest(domain, []byte(value))
}
