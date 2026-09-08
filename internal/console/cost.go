package console

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/larsonzh/prfrail/internal/tickets"
)

const defaultCostLedgerName = "cost-ledger.json"

type CostReport struct {
	LedgerPath                  string              `json:"ledgerPath"`
	ScopeKind                   string              `json:"scopeKind,omitempty"`
	Currency                    string              `json:"currency,omitempty"`
	PricingMode                 string              `json:"pricingMode,omitempty"`
	PricingVersion              string              `json:"pricingVersion,omitempty"`
	Limits                      *tickets.CostLimits `json:"limits,omitempty"`
	ReservedAmountMicros        *int64              `json:"reservedAmountMicros"`
	SettledAmountMicros         *int64              `json:"settledAmountMicros"`
	UnknownReservedAmountMicros *int64              `json:"unknownReservedAmountMicros"`
	ReservedCalls               int                 `json:"reservedCalls"`
	SettledCalls                int                 `json:"settledCalls"`
	ReservedTokens              int                 `json:"reservedTokens"`
	SettledTokens               int                 `json:"settledTokens"`
	OutstandingReservations     int                 `json:"outstandingReservations"`
	UnknownHoldReservations     int                 `json:"unknownHoldReservations"`
	TelemetryExported           bool                `json:"telemetryExported"`
	BillingDisclaimer           string              `json:"billingDisclaimer"`
	Warnings                    []string            `json:"warnings"`
}

func (cli CLI) executeCost(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: prfrail cost report [--ledger <path>] [--json]")
		return exitUsage
	}
	switch args[0] {
	case "report":
		return cli.executeCostReport(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown cost subcommand: %s\n", args[0])
		fmt.Fprintln(stderr, "usage: prfrail cost report [--ledger <path>] [--json]")
		return exitUsage
	}
}

func (cli CLI) executeCostReport(args []string, stdout, stderr io.Writer) int {
	set, parseErr := newFlagSet("cost report")
	ledgerPath := set.String("ledger", "", "path to cost ledger file")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "cost report", *jsonOutput, parseErr.String(), "usage: prfrail cost report [--ledger <path>] [--json]")
	}
	if len(set.Args()) != 0 {
		return writeUsageError(stdout, stderr, "cost report", *jsonOutput, "cost report does not accept positional arguments", "usage: prfrail cost report [--ledger <path>] [--json]")
	}
	cwd, err := cli.getwd()
	if err != nil {
		return writeCommandError(stdout, stderr, "cost report", *jsonOutput, err)
	}
	resolvedPath, err := resolveCostLedgerPath(*ledgerPath, cwd)
	if err != nil {
		return writeCommandError(stdout, stderr, "cost report", *jsonOutput, err)
	}
	report, err := BuildCostReport(resolvedPath)
	if err != nil {
		return writeCommandError(stdout, stderr, "cost report", *jsonOutput, err)
	}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "cost report", OK: true, ExitCode: exitSuccess, Data: report})
	}
	fmt.Fprintf(stdout, "ledger: %s\n", report.LedgerPath)
	if report.ScopeKind != "" {
		fmt.Fprintf(stdout, "scope: %s currency=%s pricing=%s version=%s\n", report.ScopeKind, report.Currency, report.PricingMode, report.PricingVersion)
	}
	if report.Limits != nil {
		fmt.Fprintf(stdout, "limits: amountMicros=%s calls=%d tokens=%d wallClockMs=%d attempts=%d\n", formatMicros(report.Limits.AmountMicros), report.Limits.ModelCalls, report.Limits.Tokens, report.Limits.WallClockMs, report.Limits.Attempts)
	}
	fmt.Fprintf(stdout, "summary: reservedAmountMicros=%s settledAmountMicros=%s unknownReservedAmountMicros=%s\n", formatMicros(report.ReservedAmountMicros), formatMicros(report.SettledAmountMicros), formatMicros(report.UnknownReservedAmountMicros))
	fmt.Fprintf(stdout, "calls: reserved=%d settled=%d\n", report.ReservedCalls, report.SettledCalls)
	fmt.Fprintf(stdout, "tokens: reserved=%d settled=%d\n", report.ReservedTokens, report.SettledTokens)
	fmt.Fprintf(stdout, "holds: outstanding=%d unknown=%d\n", report.OutstandingReservations, report.UnknownHoldReservations)
	fmt.Fprintf(stdout, "telemetryExported: %t\n", report.TelemetryExported)
	fmt.Fprintf(stdout, "billing: %s\n", report.BillingDisclaimer)
	for _, warning := range report.Warnings {
		fmt.Fprintf(stdout, "warning: %s\n", warning)
	}
	return exitSuccess
}

func resolveCostLedgerPath(ledgerPath, cwd string) (string, error) {
	if strings.TrimSpace(ledgerPath) == "" {
		return filepath.Abs(filepath.Join(cwd, defaultCostLedgerName))
	}
	return resolvePathFromCWD(cwd, ledgerPath)
}

func BuildCostReport(path string) (CostReport, error) {
	report := CostReport{
		LedgerPath:                  path,
		ReservedAmountMicros:        int64Ptr(0),
		SettledAmountMicros:         int64Ptr(0),
		UnknownReservedAmountMicros: int64Ptr(0),
		TelemetryExported:           false,
		BillingDisclaimer:           "reserved and settled deltas are recorded locally; reports have no default telemetry export",
		Warnings:                    []string{},
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		report.Warnings = append(report.Warnings, "ledger not found, showing empty local cost report")
		return report, nil
	} else if err != nil {
		return CostReport{}, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return CostReport{}, err
	}
	record, err := tickets.DecodeCostLedgerRecord(content)
	if err != nil {
		return CostReport{}, err
	}
	ledger, err := tickets.LoadCostLedger(record)
	if err != nil {
		return CostReport{}, err
	}
	summary := ledger.Summary()
	report.ScopeKind = record.Ledger.ScopeKind
	report.Currency = record.Ledger.Currency
	report.PricingMode = record.Ledger.PricingMode
	report.PricingVersion = record.Ledger.PricingVersion
	report.ReservedAmountMicros = cloneMicros(summary.ReservedAmountMicros)
	report.SettledAmountMicros = cloneMicros(summary.SettledAmountMicros)
	report.UnknownReservedAmountMicros = cloneMicros(summary.UnknownReservedAmountMicros)
	report.ReservedCalls = summary.ReservedCalls
	report.SettledCalls = summary.SettledCalls
	report.ReservedTokens = summary.ReservedTokens
	report.SettledTokens = summary.SettledTokens
	report.OutstandingReservations = ledger.OutstandingReservations()
	report.UnknownHoldReservations = ledger.UnknownHoldReservations()
	if limits, ok := ledger.CurrentLimits(); ok {
		copyLimits := limits
		report.Limits = &copyLimits
	}
	report.BillingDisclaimer = deriveBillingDisclaimer(report.PricingMode, summary)
	return report, nil
}

func deriveBillingDisclaimer(pricingMode string, summary tickets.CostSummary) string {
	if pricingMode == tickets.CostPricingSubscription {
		return "subscription mode enforces authorized call/token caps and does not guarantee exact monetary bills"
	}
	if summary.ReservedAmountMicros == nil || summary.SettledAmountMicros == nil || summary.UnknownReservedAmountMicros == nil {
		return "monetary amount is partially unknown; unresolved reservations remain held until observed settlement"
	}
	return "reserved and settled deltas are recorded locally; reports have no default telemetry export"
}

func formatMicros(value *int64) string {
	if value == nil {
		return "unknown"
	}
	return fmt.Sprintf("%d", *value)
}

func cloneMicros(value *int64) *int64 {
	if value == nil {
		return nil
	}
	v := *value
	return &v
}

func int64Ptr(value int64) *int64 {
	v := value
	return &v
}
