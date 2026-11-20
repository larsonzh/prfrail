package console

// P83 FIX ROUND 4 (F6) - the probe receiver gate against the operator's trusted
// authorization ledger.
//
// F6 says the probe may only reach a receiver the operator approved, and that the
// refusal happens BEFORE any credential is resolved and before any process is
// started (internal/console/ai.go authorizeAIProbeReceiver). The earlier rounds
// only pinned the shape check (validateAIProbeBinding), which is self-computed and
// therefore proves nothing about approval. This file walks every way a ledger can
// fail to approve the receiver and asserts the same three facts each time: the
// command fails, the availability checker (the only thing that resolves the
// secret) is never entered, and no availability record is written.
//
// Mutation checks:
//   - deleting the authorizeAIProbeReceiver call, or lowering any one of its four
//     conditions (active status, network scope, toolIds, executable digest), turns
//     the corresponding subtest red because the checker is entered and the spy
//     secret store records a resolution;
//   - deleting the loadProbeAuthorizationLedger os.Stat guard turns
//     TestP83Fix4LoadProbeLedgerNeverDegradesToNoGrant red, not the table: an
//     absent ledger would still be refused one layer later, but it would be
//     reported as "no grant" instead of "unreadable ledger", which is exactly the
//     degradation the load guard exists to prevent.

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/adapters"
	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

// p83Fix4Grant builds the single operator grant every subtest starts from. The
// zero value is a grant that DOES approve the probe: the active window covers the
// CLI clock (newTestCLI pins it at 2026-09-08T10:30:00Z), the network scope names
// https://api.deepseek.com:443 (the configured profile receiver), the tool id is
// the probe identity and the evidence carries the executable digest.
func p83Fix4Grant(t *testing.T, executableHash string, mutate func(*chain.AuthorizationGrant)) chain.AuthorizationRecord {
	t.Helper()
	issued := time.Date(2026, 9, 8, 10, 0, 0, 0, time.UTC)
	grant := chain.AuthorizationGrant{
		RecordID:        "p83-fix4-grant-1",
		Kind:            "grant",
		IssuedAt:        issued.Format(timestampLayout),
		IssuedBy:        evidence.Actor{Type: "operator", ID: "p83-fix4-operator"},
		AuthorizationID: "p83-fix4-authorization-1",
		RunID:           "p83-fix4-run",
		RunManifestHash: evidence.Digest("", []byte("p83-fix4-manifest")),
		PolicyHash:      evidence.Digest("", []byte("p83-fix4-policy")),
		Scope: chain.AuthorizationScope{
			ToolIDs: []string{aiProbeToolID},
			Network: []chain.AuthorizationNetworkScope{
				{Protocol: "https", Host: "api.deepseek.com", Port: 443, Purpose: "probe-provider"},
			},
			Budget: chain.AuthorizationBudgetLimit{
				ModelCalls: 1, Tokens: 1000, WallClockMs: 600000, Attempts: 1, Currency: "USD",
			},
		},
		ExpiresAt: issued.Add(time.Hour).Format(timestampLayout),
		Evidence:  p83Fix4ApprovalMaterial(t, executableHash),
	}
	if mutate != nil {
		mutate(&grant)
	}
	record, err := chain.NewGrantAuthorizationRecord(grant, nil, nil)
	if err != nil {
		t.Fatalf("build the operator grant: %v", err)
	}
	return record
}

// p83Fix4ApprovalMaterial is the grant evidence of a fully approving receiver
// approval. T3 (finding F6) requires the complete receiver identity: the profile
// hash (provider/baseUrl/model/secretRef), the receiver binding digest (profile +
// channel + executable) and the executable digest itself.
func p83Fix4ApprovalMaterial(t *testing.T, executableHash string) []string {
	t.Helper()
	profile, found := configuredAIProfile(configuredAIChain(), "agent-runner-cli")
	if !found {
		t.Fatal("the shared AI chain must configure the agent-runner-cli channel")
	}
	profileHash, err := adapters.AIProviderProfileHash(profile)
	if err != nil {
		t.Fatalf("hash the approved profile: %v", err)
	}
	binding := evidence.Digest(aiProbeReceiverBindingDomain, []byte(profileHash+"\n"+"agent-runner-cli"+"\n"+executableHash))
	return []string{executableHash, profileHash, binding}
}

// p83Fix4Revocation revokes the grant record it is handed, so EvaluateGrant sees
// a revocation whose authorizationHash matches and reports GrantRevoked.
func p83Fix4Revocation(t *testing.T, grant chain.AuthorizationRecord) chain.AuthorizationRecord {
	t.Helper()
	revocation, err := chain.NewRevocationAuthorizationRecord(chain.AuthorizationRevocation{
		RecordID:          "p83-fix4-revocation-1",
		Kind:              "revocation",
		IssuedAt:          time.Date(2026, 9, 8, 10, 15, 0, 0, time.UTC).Format(timestampLayout),
		IssuedBy:          evidence.Actor{Type: "operator", ID: "p83-fix4-operator"},
		AuthorizationID:   "p83-fix4-authorization-1",
		AuthorizationHash: grant.RecordHash,
		ReasonCode:        "policy-change",
		StopDisposition:   "not-required",
	}, nil, nil)
	if err != nil {
		t.Fatalf("build the revocation: %v", err)
	}
	return revocation
}

// writeP83Fix4Ledger persists the records at path in append order.
func writeP83Fix4Ledger(t *testing.T, path string, records ...chain.AuthorizationRecord) {
	t.Helper()
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("create the ledger: %v", err)
	}
	for _, record := range records {
		if err := appendAuthorizationRecord(path, record); err != nil {
			t.Fatalf("append the ledger record: %v", err)
		}
	}
}

// TestP83Fix4AICheckRefusesEveryUnauthorisedReceiver is the F6 ledger table: every
// failing ledger shape refuses the probe before the secret store is read.
func TestP83Fix4AICheckRefusesEveryUnauthorisedReceiver(t *testing.T) {
	otherHash := evidence.Digest("", []byte("p83-fix4-some-other-executable"))
	cases := []struct {
		name  string
		write func(t *testing.T, ledgerPath, executableHash string)
	}{
		{
			name:  "a-ledger-file-absent",
			write: func(*testing.T, string, string) {},
		},
		{
			name: "b-ledger-present-without-any-grant",
			write: func(t *testing.T, ledgerPath, _ string) {
				writeP83Fix4Ledger(t, ledgerPath)
			},
		},
		{
			name: "c-grant-expired",
			write: func(t *testing.T, ledgerPath, executableHash string) {
				writeP83Fix4Ledger(t, ledgerPath, p83Fix4Grant(t, executableHash, func(grant *chain.AuthorizationGrant) {
					grant.RecordID = "p83-fix4-grant-expired"
					grant.IssuedAt = time.Date(2026, 9, 8, 8, 0, 0, 0, time.UTC).Format(timestampLayout)
					grant.ExpiresAt = time.Date(2026, 9, 8, 9, 0, 0, 0, time.UTC).Format(timestampLayout)
				}))
			},
		},
		{
			name: "d-grant-revoked",
			write: func(t *testing.T, ledgerPath, executableHash string) {
				grant := p83Fix4Grant(t, executableHash, nil)
				writeP83Fix4Ledger(t, ledgerPath, grant, p83Fix4Revocation(t, grant))
			},
		},
		{
			name: "e-grant-network-scope-names-other-host",
			write: func(t *testing.T, ledgerPath, executableHash string) {
				writeP83Fix4Ledger(t, ledgerPath, p83Fix4Grant(t, executableHash, func(grant *chain.AuthorizationGrant) {
					grant.Scope.Network[0].Host = "api.attacker.example"
				}))
			},
		},
		{
			name: "e-grant-network-scope-names-other-port",
			write: func(t *testing.T, ledgerPath, executableHash string) {
				writeP83Fix4Ledger(t, ledgerPath, p83Fix4Grant(t, executableHash, func(grant *chain.AuthorizationGrant) {
					grant.Scope.Network[0].Port = 8443
				}))
			},
		},
		{
			name: "e-grant-network-scope-names-other-protocol",
			write: func(t *testing.T, ledgerPath, executableHash string) {
				writeP83Fix4Ledger(t, ledgerPath, p83Fix4Grant(t, executableHash, func(grant *chain.AuthorizationGrant) {
					grant.Scope.Network[0].Protocol = "http"
				}))
			},
		},
		{
			name: "f-grant-tool-ids-lack-the-probe",
			write: func(t *testing.T, ledgerPath, executableHash string) {
				writeP83Fix4Ledger(t, ledgerPath, p83Fix4Grant(t, executableHash, func(grant *chain.AuthorizationGrant) {
					grant.Scope.ToolIDs = []string{"tool-some-other-tool"}
				}))
			},
		},
		{
			name: "g-grant-evidence-lacks-the-executable-hash",
			write: func(t *testing.T, ledgerPath, executableHash string) {
				writeP83Fix4Ledger(t, ledgerPath, p83Fix4Grant(t, otherHash, nil))
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			root := t.TempDir()
			configPath := writeAIProbeTestChain(t, root)
			executablePath := filepath.Join(root, "copilot.exe")
			if err := os.WriteFile(executablePath, []byte("p83-fix4-approved-probe"), 0o700); err != nil {
				t.Fatal(err)
			}
			executableHash, err := adapters.CopilotCLIExecutableHash(executablePath)
			if err != nil {
				t.Fatalf("hash the probe executable: %v", err)
			}
			ledgerPath := filepath.Join(root, defaultApprovalsLedgerName)
			testCase.write(t, ledgerPath, executableHash)

			// The spy stand-in for the secret store: any entry proves the gate let the
			// probe through to the credential.
			secrets := &recordingProbeSecrets{}
			cli := newTestCLI(root)
			cli.secretManager = secrets
			checkerEntries := 0
			cli.checkAIAvailability = func(ctx context.Context, profile adapters.AIProviderProfile, _ string, _ int, _ aiProbeReceiver, _ string) (adapters.AIAvailabilityRecord, error) {
				checkerEntries++
				_, _ = cli.secretManager.Resolve(ctx, profile.SecretRef)
				return adapters.AIAvailabilityRecord{}, nil
			}

			outputPath := filepath.Join(root, "p83-fix4-availability.json")
			code, stdout, stderr := runCLI(t, cli,
				"ai", "check", "--chain", configPath,
				"--authorization-ledger", ledgerPath,
				"--channel", "agent-runner-cli", "--copilot", executablePath,
				"--max-requests", "1", "--out", outputPath, "--json",
			)
			if code != exitFailure {
				t.Fatalf("an unapproved receiver must fail closed, got code=%d stdout=%q stderr=%q", code, stdout, stderr)
			}
			if !strings.Contains(stdout, "receiver binding") {
				t.Fatalf("the receiver binding must be the refusing guard, got stdout=%q stderr=%q", stdout, stderr)
			}
			if checkerEntries != 0 {
				t.Fatalf("the availability checker must not run for an unapproved receiver, entries=%d", checkerEntries)
			}
			if secrets.resolveCalls != 0 {
				t.Fatalf("the secret store must not be read for an unapproved receiver, resolutions=%d", secrets.resolveCalls)
			}
			if _, statErr := os.Lstat(outputPath); !os.IsNotExist(statErr) {
				t.Fatalf("no availability record may be written, stat err=%v", statErr)
			}
		})
	}
}

// TestP83Fix4AICheckAdmitsMatchingActiveGrant is the positive control of the same
// table: with the grant untouched the probe proceeds, so the refusals above are
// attributable to the mutated condition and not to the gate refusing everything.
func TestP83Fix4AICheckAdmitsMatchingActiveGrant(t *testing.T) {
	root := t.TempDir()
	configPath := writeAIProbeTestChain(t, root)
	executablePath := filepath.Join(root, "copilot.exe")
	if err := os.WriteFile(executablePath, []byte("p83-fix4-approved-probe"), 0o700); err != nil {
		t.Fatal(err)
	}
	executableHash, err := adapters.CopilotCLIExecutableHash(executablePath)
	if err != nil {
		t.Fatalf("hash the probe executable: %v", err)
	}
	ledgerPath := filepath.Join(root, defaultApprovalsLedgerName)
	writeP83Fix4Ledger(t, ledgerPath, p83Fix4Grant(t, executableHash, nil))

	secrets := &recordingProbeSecrets{credential: []byte("p83-fix4-probe-credential")}
	cli := newTestCLI(root)
	cli.secretManager = secrets
	checkerEntries := 0
	cli.checkAIAvailability = func(ctx context.Context, profile adapters.AIProviderProfile, channel string, _ int, _ aiProbeReceiver, _ string) (adapters.AIAvailabilityRecord, error) {
		checkerEntries++
		if _, resolveErr := cli.secretManager.Resolve(ctx, profile.SecretRef); resolveErr != nil {
			return adapters.AIAvailabilityRecord{}, resolveErr
		}
		return adapters.NewAIAvailabilityRecord(adapters.AIAvailability{
			ProbeID:           "p83-fix4-probe",
			ProfileID:         profile.ProfileID,
			ProfileConfigHash: mustAIProfileHash(t, profile),
			Channel:           channel,
			ProbedAt:          "2026-09-11T08:00:00.000Z",
			Status:            "available",
			RequestsUsed:      1,
			Evidence:          []string{"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		})
	}

	outputPath := filepath.Join(root, "p83-fix4-availability.json")
	code, stdout, stderr := runCLI(t, cli,
		"ai", "check", "--chain", configPath,
		"--authorization-ledger", ledgerPath,
		"--channel", "agent-runner-cli", "--copilot", executablePath,
		"--max-requests", "1", "--out", outputPath, "--json",
	)
	if code != exitSuccess {
		t.Fatalf("a matching active grant must admit the probe, got code=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if checkerEntries != 1 {
		t.Fatalf("the approved receiver must reach the availability checker exactly once, entries=%d", checkerEntries)
	}
	if secrets.resolveCalls != 1 {
		t.Fatalf("the approved receiver must resolve its credential exactly once, resolutions=%d", secrets.resolveCalls)
	}
	if _, statErr := os.Lstat(outputPath); statErr != nil {
		t.Fatalf("an admitted probe must write its availability record: %v", statErr)
	}
}

// TestP83Fix4LoadProbeLedgerNeverDegradesToNoGrant pins the load half in isolation:
// an absent path is an error, not an empty ledger. A caller that swallowed this
// error would hand authorizeAIProbeReceiver an empty ledger, which looks exactly
// like "no grant needed" only because the gate refuses it - so the two layers are
// pinned separately.
func TestP83Fix4LoadProbeLedgerNeverDegradesToNoGrant(t *testing.T) {
	root := t.TempDir()
	missing := filepath.Join(root, "p83-fix4-missing-ledger.jsonl")
	ledger, err := loadProbeAuthorizationLedger(missing)
	if err == nil {
		t.Fatalf("an absent ledger must be an error, got ledger=%+v", ledger)
	}
	if !strings.Contains(err.Error(), "receiver binding") {
		t.Fatalf("an absent ledger must be reported as a receiver binding refusal, got %v", err)
	}

	present := filepath.Join(root, "p83-fix4-empty-ledger.jsonl")
	writeP83Fix4Ledger(t, present)
	ledger, err = loadProbeAuthorizationLedger(present)
	if err != nil {
		t.Fatalf("a present empty ledger must load: %v", err)
	}
	if len(ledger.Grants) != 0 || len(ledger.Revocations) != 0 {
		t.Fatalf("a present empty ledger must carry no records, got %+v", ledger)
	}
}
