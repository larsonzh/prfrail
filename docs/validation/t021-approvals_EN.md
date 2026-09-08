# T021 PC-03 Authorization and Approvals Validation Report

[中文](t021-approvals.md)

Date: 2026-09-08. Verdict: `T021 COMPLETE`.

## Implementation Scope

1. `internal/chain/authorization.go`
2. `internal/chain/authorization_test.go`
3. `internal/console/approvals.go`
4. `internal/console/approvals_test.go`
5. `internal/console/cli.go` (approvals command routing and usage)
6. `docs/validation/t021-approvals.md` / `t021-approvals_EN.md` (this report)

## AT-18 Coverage

1. Authorization binds principal, run/manifest hashes, tool/path/network scope, budget and validity: `AuthorizationGrant` fields align with `authorization-record.schema.json`; invalid principals, scopes, budgets or timestamps are rejected on write.
2. Revocation/expiry/hash mismatch block further actions and acceptance: `EvaluateGrant` returns `revoked` permanently once a revocation matches (a later grant with the same id cannot resurrect it), `expired` at `expiresAt`, and fails closed when the revocation `authorizationHash` matches no grant record for that id.
3. Failed in-flight stop stays paused: `RequestControlledStop` wires `stopDisposition=requested` to `Stopper.Stop`; stop failure returns `uncertain` plus `ErrRecoveryUncertain`, never `completed`; `approvals revoke` exits non-zero on uncertain stop while the revocation is already persisted (restart-visible).
4. Hard gates cannot be waived: the model has no waiver field or bypass path; `issuedBy` accepts only `operator`/`policy`, so agent self-authorization is rejected.
5. Pending inbox survives restart: `approvals list` rebuilds the pending queue (`pendingCount`, status, next action) from the append-only `authorization-ledger.jsonl` ledger, which persists across restarts.
6. Waiting never auto-approves: grants evaluated before `issuedAt` are `pending` (`await-issue`), expiry yields `expired` (`reissue`); no automatic approval path exists.

## Gates

1. `go test -count=1 ./internal/chain ./internal/console`: passed.
2. `go build ./...`: passed.
3. `go vet ./...`: passed.
4. `go test -count=1 ./...`: passed (all packages ok).
5. Editor diagnostics: `internal/chain/authorization.go`, `internal/console/approvals.go`, `internal/console/cli.go` have no errors.

## Key Command Evidence

```powershell
gofmt -w internal/chain/authorization.go internal/chain/authorization_test.go internal/console/approvals.go internal/console/approvals_test.go internal/console/cli.go
go test -count=1 ./internal/chain ./internal/console
go build ./...
go vet ./...
go test -count=1 ./...
```

## Boundaries

1. `approvals` exposes only `list`/`revoke`; grants are written by the trusted policy issuance flow — the CLI has no `grant` entrypoint (chat text or repository files can never promote into authorization).
2. Stop wiring uses the `chain.Stopper` abstraction: `approvals revoke` resolves `managed-process.identity.json` via `internal/console/stopper.go` and calls the guard stop path by default; local terminal noop runs return terminal evidence digests, and missing process identity fails closed as `uncertain`.
3. `not-required` revocations never call the stopper; `requested` success records `completed`, failure records `uncertain` requiring `reconcile-stop`.
4. The ledger is an append-only JSONL file: one canonical authorization record per line, parse failures reported with line numbers; no in-place edits or deletions.
5. Validity boundaries use strict timestamp comparison (`expiresAt` must be after `issuedAt`; evaluation before issuance yields `pending`).

## Review Notes (2026-09-08)

1. `EvaluateGrant` revocation matching was corrected to match against any grant record of the same authorizationId (multi-grant case); mismatches still fail closed. Test updated: `TestEvaluateGrantRevocationPermanentAndHashMismatchFailsClosed`.
2. `newDefaultStopRun` now wraps controlled stop with a bounded timeout context (`defaultStopTimeout=12s`, `defaultStopGrace=2s`) to prevent potential unbounded waiting during revoke stop verification. Added test: `TestDefaultStopRunAppliesBoundedStopTimeout`.
3. `BuildApprovalsReport` now aggregates by the latest grant per `authorizationId`, so `pendingCount`/`revokedCount` reflect the current approvals inbox instead of historical grant row count. Added test: `TestApprovalsReportUsesLatestGrantPerAuthorization`.
4. All gates re-run and passed.
