package console

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/chain"
	"github.com/larsonzh/prfrail/internal/evidence"
)

const defaultApprovalsLedgerName = "authorization-ledger.jsonl"

// ApprovalItem is one row of the approvals inbox view.
type ApprovalItem struct {
	AuthorizationID string `json:"authorizationId"`
	Status          string `json:"status"`
	ExpiresAt       string `json:"expiresAt,omitempty"`
	IssuedBy        string `json:"issuedBy,omitempty"`
	ReasonCode      string `json:"reasonCode,omitempty"`
	Disposition     string `json:"disposition,omitempty"`
	NextAction      string `json:"nextAction"`
}

// ApprovalsReport is the offline read-only view over an authorization ledger.
type ApprovalsReport struct {
	LedgerPath   string         `json:"ledgerPath"`
	Grants       []ApprovalItem `json:"grants"`
	Revocations  []ApprovalItem `json:"revocations"`
	PendingCount int            `json:"pendingCount"`
	RevokedCount int            `json:"revokedCount"`
	Warnings     []string       `json:"warnings"`
}

// RevocationOutcome is returned by the revoke subcommand.
type RevocationOutcome struct {
	AuthorizationID string                        `json:"authorizationId"`
	ReasonCode      string                        `json:"reasonCode"`
	Disposition     string                        `json:"disposition"`
	RecordHash      string                        `json:"recordHash"`
	Revocation      chain.AuthorizationRevocation `json:"revocation"`
}

type stopperAdapter func(context.Context, string) ([]string, error)

func (adapter stopperAdapter) Stop(ctx context.Context, runID string) ([]string, error) {
	return adapter(ctx, runID)
}

func (cli CLI) executeApprovals(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: prfrail approvals <list|revoke> [options]")
		return exitUsage
	}
	switch args[0] {
	case "list":
		return cli.executeApprovalsList(args[1:], stdout, stderr)
	case "revoke":
		return cli.executeApprovalsRevoke(ctx, args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown approvals subcommand: %s\n", args[0])
		fmt.Fprintln(stderr, "usage: prfrail approvals <list|revoke> [options]")
		return exitUsage
	}
}

func (cli CLI) executeApprovalsList(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("approvals list")
	ledgerPath := set.String("ledger", "", "path to authorization ledger file")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "approvals list", *jsonOutput, parseErr.String(), "usage: prfrail approvals list [--ledger <path>] [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "approvals list", *jsonOutput, "approvals list does not accept positional arguments", "usage: prfrail approvals list [--ledger <path>] [--json]")
	}
	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "approvals list", *jsonOutput, err)
	}
	resolvedPath, err := resolveApprovalsLedgerPath(*ledgerPath, cwd)
	if err != nil {
		return writeCommandError(stdout, stderr, "approvals list", *jsonOutput, err)
	}
	report, err := BuildApprovalsReport(resolvedPath, cli.now().UTC())
	if err != nil {
		return writeCommandError(stdout, stderr, "approvals list", *jsonOutput, err)
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "approvals list", OK: true, ExitCode: exitSuccess, Data: report})
	}
	fmt.Fprintf(stdout, "ledger: %s\n", report.LedgerPath)
	fmt.Fprintf(stdout, "grants: %d pending=%d revoked=%d\n", len(report.Grants), report.PendingCount, report.RevokedCount)
	for _, item := range report.Grants {
		if item.Status == string(chain.GrantActive) {
			fmt.Fprintf(stdout, "authorization: %s status=%s expires=%s next=%s\n", item.AuthorizationID, item.Status, item.ExpiresAt, item.NextAction)
		} else {
			fmt.Fprintf(stdout, "authorization: %s status=%s next=%s\n", item.AuthorizationID, item.Status, item.NextAction)
		}
	}
	for _, item := range report.Revocations {
		fmt.Fprintf(stdout, "revocation: %s reason=%s disposition=%s next=%s\n", item.AuthorizationID, item.ReasonCode, item.Disposition, item.NextAction)
	}
	for _, warning := range report.Warnings {
		fmt.Fprintf(stdout, "warning: %s\n", warning)
	}
	return exitSuccess
}

func (cli CLI) executeApprovalsRevoke(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("approvals revoke")
	authorizationID := set.String("authorization-id", "", "authorization id to revoke")
	reasonCode := set.String("reason-code", "", "revocation reason code")
	ledgerPath := set.String("ledger", "", "path to authorization ledger file")
	runDir := set.String("run-dir", "", "run directory for managed process identity")
	stopDisposition := set.String("stop", "requested", "stop disposition: not-required or requested")
	actorID := set.String("actor-id", "local-operator", "revoking operator identity")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "approvals revoke", *jsonOutput, parseErr.String(), "usage: prfrail approvals revoke --authorization-id <id> --reason-code <code> [--ledger <path>] [--run-dir <path>] [--stop not-required|requested] [--actor-id <id>] [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "approvals revoke", *jsonOutput, "approvals revoke does not accept positional arguments", "usage: prfrail approvals revoke --authorization-id <id> --reason-code <code> [--ledger <path>] [--run-dir <path>] [--stop not-required|requested] [--actor-id <id>] [--json]")
	}
	if *authorizationID == "" || *reasonCode == "" {
		return writeUsageError(stdout, stderr, "approvals revoke", *jsonOutput, "--authorization-id and --reason-code are required", "usage: prfrail approvals revoke --authorization-id <id> --reason-code <code> [--ledger <path>] [--run-dir <path>] [--stop not-required|requested] [--actor-id <id>] [--json]")
	}
	if *stopDisposition != "not-required" && *stopDisposition != "requested" {
		return writeUsageError(stdout, stderr, "approvals revoke", *jsonOutput, "--stop must be not-required or requested", "usage: prfrail approvals revoke --authorization-id <id> --reason-code <code> [--ledger <path>] [--run-dir <path>] [--stop not-required|requested] [--actor-id <id>] [--json]")
	}
	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "approvals revoke", *jsonOutput, err)
	}
	resolvedPath, err := resolveApprovalsLedgerPath(*ledgerPath, cwd)
	if err != nil {
		return writeCommandError(stdout, stderr, "approvals revoke", *jsonOutput, err)
	}
	ledger, err := loadAuthorizationLedger(resolvedPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "approvals revoke", *jsonOutput, err)
	}
	grantRecord, grant, exists := ledger.LatestGrant(*authorizationID)
	if !exists {
		return writeCommandError(stdout, stderr, "approvals revoke", *jsonOutput, fmt.Errorf("no grant for authorization %q", *authorizationID))
	}
	revocation := chain.AuthorizationRevocation{
		Kind:                 "revocation",
		IssuedAt:             cli.now().UTC().Format(timestampLayout),
		IssuedBy:             evidence.Actor{Type: "operator", ID: *actorID},
		AuthorizationID:      *authorizationID,
		AuthorizationHash:    grantRecord.RecordHash,
		ReasonCode:           *reasonCode,
		StopDisposition:      *stopDisposition,
		StopEvidence:         []string{},
		ResidualRiskEvidence: []string{},
	}
	switch *stopDisposition {
	case "not-required":
		_ = grant
	case "requested":
		revocation.ResidualRiskEvidence = []string{digest("proofrail:stop-uncertain:1\n", *authorizationID, *reasonCode)}
	}
	stopResult, stopErr := chain.RequestControlledStop(ctx, revocation, stopperAdapter(func(stopCtx context.Context, id string) ([]string, error) {
		return cli.stopRun(stopCtx, id, *runDir)
	}), grant.RunID)
	if stopErr != nil && len(stopResult.StopEvidence) == 0 {
		stopResult.StopEvidence = []string{digest("proofrail:local-approvals-stop:1\n", grant.RunID)}
	}
	revocation.StopDisposition = stopResult.Disposition
	revocation.StopEvidence = append([]string{}, stopResult.StopEvidence...)
	revocation.ResidualRiskEvidence = append([]string{}, stopResult.ResidualRiskEvidence...)
	record, err := chain.NewRevocationAuthorizationRecord(revocation, cli.now, nil)
	if err != nil {
		return writeCommandError(stdout, stderr, "approvals revoke", *jsonOutput, err)
	}
	if err := appendAuthorizationRecord(resolvedPath, record); err != nil {
		return writeCommandError(stdout, stderr, "approvals revoke", *jsonOutput, err)
	}
	outcome := RevocationOutcome{
		AuthorizationID: *authorizationID,
		ReasonCode:      *reasonCode,
		Disposition:     revocation.StopDisposition,
		RecordHash:      record.RecordHash,
		Revocation:      revocation,
	}
	if stopErr != nil || revocation.StopDisposition == "uncertain" {
		return writeCommandError(stdout, stderr, "approvals revoke", *jsonOutput, fmt.Errorf("controlled stop uncertain (revocation recorded, run stays paused): %s", record.RecordHash))
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "approvals revoke", OK: true, ExitCode: exitSuccess, Message: "authorization revoked", Data: outcome})
	}
	fmt.Fprintf(stdout, "revoked: %s disposition=%s reason=%s hash=%s\n", *authorizationID, revocation.StopDisposition, *reasonCode, record.RecordHash)
	return exitSuccess
}

func resolveApprovalsLedgerPath(ledgerPath, cwd string) (string, error) {
	if ledgerPath == "" {
		return filepath.Abs(filepath.Join(cwd, defaultApprovalsLedgerName))
	}
	return resolvePathFromCWD(cwd, ledgerPath)
}

// BuildApprovalsReport evaluates every grant in the ledger against the
// current instant and summarizes pending/revoked rows with next actions.
func BuildApprovalsReport(ledgerPath string, now time.Time) (ApprovalsReport, error) {
	report := ApprovalsReport{LedgerPath: ledgerPath, Grants: []ApprovalItem{}, Revocations: []ApprovalItem{}}
	if _, err := os.Stat(ledgerPath); os.IsNotExist(err) {
		report.Warnings = append(report.Warnings, "ledger not found, showing empty inbox")
		return report, nil
	}
	ledger, err := loadAuthorizationLedger(ledgerPath)
	if err != nil {
		return ApprovalsReport{}, err
	}
	latestGrants := make(map[string]chain.AuthorizationRecord, len(ledger.Grants))
	grantOrder := make([]string, 0, len(ledger.Grants))
	for _, record := range ledger.Grants {
		authorizationID := record.Grant.AuthorizationID
		if _, exists := latestGrants[authorizationID]; !exists {
			grantOrder = append(grantOrder, authorizationID)
		}
		latestGrants[authorizationID] = record
	}
	for _, authorizationID := range grantOrder {
		record := latestGrants[authorizationID]
		status, err := chain.EvaluateGrant(authorizationID, ledger, now)
		if err != nil {
			return ApprovalsReport{}, err
		}
		item := ApprovalItem{
			AuthorizationID: authorizationID,
			Status:          string(status),
			ExpiresAt:       record.Grant.ExpiresAt,
			IssuedBy:        fmt.Sprintf("%s:%s", record.Grant.IssuedBy.Type, record.Grant.IssuedBy.ID),
		}
		switch status {
		case chain.GrantPending:
			item.NextAction = "await-issue"
			report.PendingCount++
		case chain.GrantActive:
			item.NextAction = "continue-under-scope"
		case chain.GrantExpired:
			item.NextAction = "reissue"
		case chain.GrantRevoked:
			item.NextAction = "stop-verify"
			report.RevokedCount++
		default:
			item.NextAction = "blocked"
		}
		report.Grants = append(report.Grants, item)
	}
	for _, record := range ledger.Revocations {
		item := ApprovalItem{
			AuthorizationID: record.Revocation.AuthorizationID,
			Status:          "recorded",
			ReasonCode:      record.Revocation.ReasonCode,
			Disposition:     record.Revocation.StopDisposition,
			IssuedBy:        fmt.Sprintf("%s:%s", record.Revocation.IssuedBy.Type, record.Revocation.IssuedBy.ID),
		}
		if record.Revocation.StopDisposition == "uncertain" {
			item.NextAction = "reconcile-stop"
		} else {
			item.NextAction = "none"
		}
		report.Revocations = append(report.Revocations, item)
	}
	sort.Slice(report.Grants, func(i, j int) bool { return report.Grants[i].AuthorizationID < report.Grants[j].AuthorizationID })
	sort.Slice(report.Revocations, func(i, j int) bool {
		return report.Revocations[i].AuthorizationID < report.Revocations[j].AuthorizationID
	})
	return report, nil
}

func loadAuthorizationLedger(path string) (chain.AuthorizationLedger, error) {
	ledger := chain.AuthorizationLedger{}
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return ledger, nil
	}
	if err != nil {
		return ledger, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		record, err := chain.DecodeAuthorizationRecord([]byte(text))
		if err != nil {
			return ledger, fmt.Errorf("authorization ledger line %d: %w", line, err)
		}
		if err := chain.ValidateAuthorizationRecord(record); err != nil {
			return ledger, fmt.Errorf("authorization ledger line %d: %w", line, err)
		}
		if record.Grant != nil {
			ledger.Grants = append(ledger.Grants, record)
		} else {
			ledger.Revocations = append(ledger.Revocations, record)
		}
	}
	if err := scanner.Err(); err != nil {
		return ledger, err
	}
	return ledger, nil
}

func appendAuthorizationRecord(path string, record chain.AuthorizationRecord) error {
	content, err := json.Marshal(record)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(append(content, '\n')); err != nil {
		return err
	}
	return file.Sync()
}
