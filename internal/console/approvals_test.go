package console

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

func approvalsTestHash(seed byte) string {
	return "sha256:" + strings.Repeat(string(seed), 64)
}

func approvalsTestGrant() chain.AuthorizationGrant {
	amount := int64(500)
	return chain.AuthorizationGrant{
		RecordID:        "grant-r1",
		Kind:            "grant",
		IssuedAt:        "2026-09-08T10:00:00.000Z",
		IssuedBy:        evidence.Actor{Type: "operator", ID: "alice"},
		AuthorizationID: "auth-tool-build",
		RunID:           "run-one",
		RunManifestHash: approvalsTestHash('a'),
		PolicyHash:      approvalsTestHash('b'),
		Scope: chain.AuthorizationScope{
			TaskIDs:       []string{"task-build"},
			StepIDs:       []string{"step-build"},
			ToolIDs:       []string{"tool-go-build"},
			TargetIDs:     []string{"target-src"},
			Network:       []chain.AuthorizationNetworkScope{{Protocol: "https", Host: "example.com", Port: 443, Purpose: "fetch-module"}},
			EffectClasses: []string{"local-discardable"},
			Budget: chain.AuthorizationBudgetLimit{
				ModelCalls: 10, Tokens: 10000, WallClockMs: 60000, Attempts: 3,
				Currency: "USD", AmountMicros: &amount,
			},
		},
		ExpiresAt: "2026-09-09T10:00:00.000Z",
		Evidence:  []string{approvalsTestHash('c')},
	}
}

func appendTestGrant(t *testing.T, path string, cli CLI) chain.AuthorizationRecord {
	t.Helper()
	record, err := chain.NewGrantAuthorizationRecord(approvalsTestGrant(), cli.now, func() string { return "grant-record-1" })
	if err != nil {
		t.Fatal(err)
	}
	content, err := json.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if _, err := file.Write(append(content, '\n')); err != nil {
		t.Fatal(err)
	}
	return record
}

func TestApprovalsListEmptyLedgerZeroPendingWithoutSideEffects(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "authorization-ledger.jsonl")

	code, stdout, stderr := runCLI(t, cli, "approvals", "list", "--ledger", ledgerPath, "--json")
	if code != 0 {
		t.Fatalf("list code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result := decodeResponse(t, stdout)
	if !result.OK {
		t.Fatalf("list failed: %s", stdout)
	}
	var report ApprovalsReport
	if err := json.Unmarshal(result.Data, &report); err != nil {
		t.Fatal(err)
	}
	if report.PendingCount != 0 || report.RevokedCount != 0 || len(report.Grants) != 0 {
		t.Fatalf("unexpected empty report: %+v", report)
	}
	if len(report.Warnings) == 0 {
		t.Fatal("expected missing-ledger warning")
	}
	if _, err := os.Stat(ledgerPath); !os.IsNotExist(err) {
		t.Fatal("list must not create the ledger file")
	}
}

func TestApprovalsRevokeRequestedPersistsAndListsRevoked(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "authorization-ledger.jsonl")
	appendTestGrant(t, ledgerPath, cli)
	cli.stopRun = func(_ context.Context, runID, _ string) ([]string, error) {
		return []string{approvalsTestHash('e'), digest("proofrail:local-approvals-stop:1\n", runID)}, nil
	}

	code, stdout, stderr := runCLI(t, cli, "approvals", "revoke", "--ledger", ledgerPath, "--authorization-id", "auth-tool-build", "--reason-code", "policy-change", "--json")
	if code != 0 {
		t.Fatalf("revoke code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result := decodeResponse(t, stdout)
	if !result.OK {
		t.Fatalf("revoke failed: %s", stdout)
	}
	var outcome RevocationOutcome
	if err := json.Unmarshal(result.Data, &outcome); err != nil {
		t.Fatal(err)
	}
	if outcome.Disposition != "completed" || outcome.AuthorizationID != "auth-tool-build" {
		t.Fatalf("unexpected outcome: %+v", outcome)
	}

	// Restart visibility: a fresh list over the same ledger must still see the revocation.
	code, stdout, stderr = runCLI(t, cli, "approvals", "list", "--ledger", ledgerPath, "--json")
	if code != 0 {
		t.Fatalf("list code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	var report ApprovalsReport
	if err := json.Unmarshal(decodeResponse(t, stdout).Data, &report); err != nil {
		t.Fatal(err)
	}
	if report.RevokedCount != 1 || len(report.Revocations) != 1 {
		t.Fatalf("unexpected report after revoke: %+v", report)
	}
	if report.Grants[0].Status != string(chain.GrantRevoked) || report.Grants[0].NextAction != "stop-verify" {
		t.Fatalf("grant must be revoked: %+v", report.Grants[0])
	}
	if report.Revocations[0].Disposition != "completed" || report.Revocations[0].ReasonCode != "policy-change" {
		t.Fatalf("unexpected revocation row: %+v", report.Revocations[0])
	}
}

func TestApprovalsRevokeStopUncertainRecordsRevocationAndFails(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "authorization-ledger.jsonl")
	appendTestGrant(t, ledgerPath, cli)
	cli.stopRun = func(_ context.Context, runID, _ string) ([]string, error) {
		return []string{approvalsTestHash('e')}, errors.New("process refused to stop")
	}

	code, stdout, stderr := runCLI(t, cli, "approvals", "revoke", "--ledger", ledgerPath, "--authorization-id", "auth-tool-build", "--reason-code", "security-alert", "--json")
	if code != 1 {
		t.Fatalf("uncertain revoke must fail, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result := decodeResponse(t, stdout)
	if result.OK || !strings.Contains(result.Error, "uncertain") {
		t.Fatalf("expected uncertain error response: %s", stdout)
	}
	report, err := BuildApprovalsReport(ledgerPath, cli.now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if report.Grants[0].Status != string(chain.GrantRevoked) {
		t.Fatalf("revocation must be recorded despite stop failure: %+v", report.Grants[0])
	}
	if report.Revocations[0].Disposition != "uncertain" || report.Revocations[0].NextAction != "reconcile-stop" {
		t.Fatalf("unexpected uncertain revocation row: %+v", report.Revocations[0])
	}
}

func TestApprovalsRevokeNotRequiredSkipsStopper(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "authorization-ledger.jsonl")
	appendTestGrant(t, ledgerPath, cli)
	calls := 0
	cli.stopRun = func(_ context.Context, _, _ string) ([]string, error) {
		calls++
		return nil, errors.New("must not be called")
	}

	code, stdout, stderr := runCLI(t, cli, "approvals", "revoke", "--ledger", ledgerPath, "--authorization-id", "auth-tool-build", "--reason-code", "scope-shrink", "--stop", "not-required", "--json")
	if code != 0 {
		t.Fatalf("revoke code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if calls != 0 {
		t.Fatalf("stopper called %d times for not-required", calls)
	}
	var outcome RevocationOutcome
	if err := json.Unmarshal(decodeResponse(t, stdout).Data, &outcome); err != nil {
		t.Fatal(err)
	}
	if outcome.Disposition != "not-required" {
		t.Fatalf("unexpected outcome: %+v", outcome)
	}
}

func TestApprovalsRevokeRequestedFailsClosedWhenManagedIdentityMissing(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "authorization-ledger.jsonl")
	appendTestGrant(t, ledgerPath, cli)

	code, stdout, stderr := runCLI(t, cli, "approvals", "revoke", "--ledger", ledgerPath, "--authorization-id", "auth-tool-build", "--reason-code", "runtime-stop", "--json")
	if code != 1 {
		t.Fatalf("missing managed identity must fail closed, code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result := decodeResponse(t, stdout)
	if result.OK || !strings.Contains(result.Error, "uncertain") {
		t.Fatalf("expected uncertain response: %s", stdout)
	}
	report, err := BuildApprovalsReport(ledgerPath, cli.now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Revocations) != 1 || report.Revocations[0].Disposition != "uncertain" {
		t.Fatalf("expected uncertain revocation row: %+v", report.Revocations)
	}
}

func TestApprovalsRevokeMissingGrantFails(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "authorization-ledger.jsonl")
	code, stdout, _ := runCLI(t, cli, "approvals", "revoke", "--ledger", ledgerPath, "--authorization-id", "auth-none", "--reason-code", "cleanup", "--json")
	if code != 1 {
		t.Fatalf("missing grant must fail, code=%d stdout=%s", code, stdout)
	}
	if !strings.Contains(decodeResponse(t, stdout).Error, "no grant") {
		t.Fatalf("unexpected error: %s", stdout)
	}
}

func TestApprovalsListTextHasNoANSI(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "authorization-ledger.jsonl")
	appendTestGrant(t, ledgerPath, cli)

	code, stdout, stderr := runCLI(t, cli, "approvals", "list", "--ledger", ledgerPath)
	if code != 0 {
		t.Fatalf("list code=%d stdout=%s stderr=%s", code, stdout, stderr)
	}
	if strings.Contains(stdout, "\x1b") || strings.Contains(stderr, "\x1b") {
		t.Fatalf("unexpected ANSI stdout=%q stderr=%q", stdout, stderr)
	}
	if !strings.Contains(stdout, "status=active") || !strings.Contains(stdout, "next=continue-under-scope") {
		t.Fatalf("text output missing facts: %s", stdout)
	}
}

func TestApprovalsReportUsesLatestGrantPerAuthorization(t *testing.T) {
	root := t.TempDir()
	cli := newTestCLI(root)
	ledgerPath := filepath.Join(root, "authorization-ledger.jsonl")

	first := appendTestGrant(t, ledgerPath, cli)
	secondGrant := approvalsTestGrant()
	secondGrant.RecordID = "grant-r2"
	secondGrant.IssuedAt = "2026-09-08T11:00:00.000Z"
	secondGrant.ExpiresAt = "2026-09-10T11:00:00.000Z"
	secondGrant.Evidence = []string{approvalsTestHash('f')}
	second, err := chain.NewGrantAuthorizationRecord(secondGrant, cli.now, func() string { return "grant-record-2" })
	if err != nil {
		t.Fatal(err)
	}
	if err := appendAuthorizationRecord(ledgerPath, second); err != nil {
		t.Fatal(err)
	}

	revocation := chain.AuthorizationRevocation{
		RecordID:             "rev-r1",
		Kind:                 "revocation",
		IssuedAt:             "2026-09-08T12:00:00.000Z",
		IssuedBy:             evidence.Actor{Type: "operator", ID: "bob"},
		AuthorizationID:      "auth-tool-build",
		AuthorizationHash:    first.RecordHash,
		ReasonCode:           "policy-change",
		StopDisposition:      "completed",
		StopEvidence:         []string{approvalsTestHash('d')},
		ResidualRiskEvidence: []string{},
	}
	revocationRecord, err := chain.NewRevocationAuthorizationRecord(revocation, cli.now, func() string { return "rev-record-1" })
	if err != nil {
		t.Fatal(err)
	}
	if err := appendAuthorizationRecord(ledgerPath, revocationRecord); err != nil {
		t.Fatal(err)
	}

	report, err := BuildApprovalsReport(ledgerPath, cli.now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Grants) != 1 {
		t.Fatalf("latest-state grants should deduplicate by authorizationId: %+v", report.Grants)
	}
	if report.RevokedCount != 1 || report.Grants[0].Status != string(chain.GrantRevoked) {
		t.Fatalf("deduplicated grant should reflect latest revoked state: %+v", report)
	}
}
