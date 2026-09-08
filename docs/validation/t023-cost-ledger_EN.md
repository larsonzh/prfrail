# T023 PC-05 Cost Ledger Validation Report

[中文](t023-cost-ledger.md)

Date: 2026-09-08. Verdict: `T023 COMPLETE`.

## Implementation Scope

1. `internal/tickets/cost.go`: cost-ledger model, append-only records, reservation/settlement rules, replay summary validation, and shared-cap admission checks.
2. `internal/tickets/cost_test.go`: coverage for cross-run cap enforcement, restart-preserved unknown holds, settlement idempotency, reservation lifecycle, and summary-drift rejection.
3. `internal/tickets/budget.go`: added `CostBudgetExhausted` and `CostTokenBudgetExhausted` for fail-closed monetary/token admission checks.
4. `internal/adapters/cost_usage.go` and `internal/adapters/cost_usage_test.go`: settlement-ready usage extraction from durable result/dispatch receipts with stable idempotency keys.
5. `internal/console/cost.go` and `internal/console/cost_test.go`: `prfrail cost report` JSON/text output, no default telemetry export, and subscription billing disclaimer.
6. `internal/console/cli.go`: wired `cost` command routing.
7. `docs/OPERATIONS.md` and `docs/OPERATIONS_EN.md`: synchronized executable command examples and PC-05 operational guidance.
8. `docs/validation/t023-cost-ledger.md` and `docs/validation/t023-cost-ledger_EN.md`: this validation pair.

## AT-20 Coverage

1. Durable pre-call reservation: `Reserve` appends reservation records before settlement paths; `TestCostLedgerCrashPointReservationLifecycle` verifies holds persist until settlement.
2. Settlement deduplication and conflict blocking: `Settle` is idempotent on `idempotencyKey`; changed payload for the same key is rejected; `TestCostLedgerSettlementDeduplicates` covers both paths.
3. Restart-preserved unknown holds: `Snapshot` + `LoadCostLedger` replay keeps unknown holds; `TestCostLedgerUnknownHoldSurvivesRestart` verifies holds remain after restart and still block over-cap requests.
4. Shared cap across runs: `validateSharedCaps` checks settled + outstanding usage; `TestCostLedgerSharedCapAcrossRuns` validates cross-run cumulative blocking.
5. Fail-closed budget checks: `CostBudgetExhausted` treats unknown committed monetary usage as exhausted; `CostTokenBudgetExhausted` uses settled + hold tokens; `TestCostBudgetExhaustionHelpers` validates boundaries.
6. Local reporting with no default telemetry: `BuildCostReport` hard-sets `TelemetryExported=false`; `TestCostReportJSONAndText`, `TestCostReportSubscriptionDisclaimer`, and `TestCostReportMissingLedgerIsNonFatal` verify reporting facts and disclaimers.
7. Verifiable usage evidence: `UsageFromResultEnvelope` and `UsageFromDispatchReceipt` aggregate request/receipt/result evidence into settlement-ready usage structures; tests verify evidence completeness and stable keys.

## Gate Results

1. `go test -count=1 ./internal/tickets ./internal/adapters ./internal/console`: passed.
2. `go build ./...`: passed.
3. `go vet ./...`: passed.
4. `go test -count=1 ./...`: passed (all packages ok).

## Key Command Evidence

```powershell
go test -count=1 ./internal/tickets ./internal/adapters ./internal/console
go build ./...
go vet ./...
go test -count=1 ./...
```

## Boundaries

1. T023 delivers local ledgering and cost observability only; it does not integrate provider billing APIs directly. Unknown cost remains held; once a reservation is settled as unknown, that same reservation cannot be settled again and reconciliation must happen in a new verifiable accounting cycle.
2. Subscription mode explicitly guarantees authorized call/token caps, not exact monetary billing parity.
3. `cost report` does not create or export telemetry; if the ledger is missing, it returns a warning and an empty local report without fabricating history.
4. This delivery does not include T024 lifecycle disposition, backup replay, or release inventory capabilities.

## Review Notes (2026-09-08)

1. `docs/OPERATIONS_EN.md` was corrected from “T023 pending” to “T023 implemented; T024 pending” to match the Chinese operations document.
2. Full gates were re-run and passed with no new compile, vet, or test regressions.