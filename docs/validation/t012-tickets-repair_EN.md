# T012 Tickets/Repair Validation Report

[简体中文](t012-tickets-repair.md)

Date: 2026-09-07. Verdict: `T012 COMPLETE`. This task implements failure fingerprinting, append-only ticket ledger, budget gating, and the four-stage repair transaction chain, with AT-07 critical scenarios validated.

## Scope

1. `internal/tickets/fingerprint.go`: failure fingerprint input validation and stable digest (`proofrail:failure-fingerprint:1`).
2. `internal/tickets/ledger.go`: append-only ledger entries (failure/override-granted/resolved), derived `currentBudgetState`, and `ticketLedgerHash` (`proofrail:ticket-ledger:1`).
3. `internal/tickets/budget.go`: repair attempt reservation and budget gates (hard-block, review gate, override window, wall-clock/cost exhaustion).
4. `internal/repair/transaction.go`: Prepare/Inspect/Validate/Promote ordered stage chain, adjacent `stageHash` binding (`proofrail:repair-stage:1`), and transaction hash (`proofrail:repair-transaction:1`).

## AT-07 Scenario Coverage

1. Same-fingerprint exhaustion: failure counts reach hard-block threshold; resolved does not reset budget.
2. Attempt exhaustion: override window rejects attempts beyond `maxAdditionalAttempts`; duplicate consumption of the same attempt is blocked.
3. Wall-clock/cost exhaustion: `ReserveAttempt` fails immediately when either budget is exhausted.
4. Fake-repair blocking: Promote rejects mismatched `validateStageHash` or `ticketLedgerHash` bindings.
5. Stale-candidate blocking: Inspect/Validate/Promote reject candidate hashes that diverge from Prepare output.
6. Interrupted Promote: `outcome=uncertain` becomes terminal and rejects further stage appends.

## Tests

- `internal/tickets/ledger_test.go`
- `internal/repair/transaction_test.go`
- `internal/repair/transaction_fingerprint_budget_test.go`

## Gates

- `go build ./...` passed.
- `go vet ./...` passed.
- `go test -count=1 ./...` passed.

## Review Fixes

A 2026-09-07 review found and fixed the following schema/semantic conformance issues:

1. Stage serialization field scoping: added `Stage.MarshalJSON` to emit only stage-specific fields and to serialize prepare failed/uncertain `candidateManifestHash` and promote failed/uncertain `resultingManifestHash` as `null` (previously omitted by omitempty, violating schema required).
2. Non-nil evidence slices: `uniqueHashes` returns `[]` for empty input, avoiding nil serializing as `null` (schema requires arrays).
3. Non-empty stage evidence: Prepare/Inspect/Validate/Promote `evidence` requires at least one item (nonemptyHashSet).
4. `targetIds` uniqueness: duplicate targets are rejected before writing (idSet uniqueItems).
5. Ledger `ticketId` uniqueness: duplicate ticket IDs are rejected.
6. Budget state derivation: an exhausted override falls back to `pending-review` (contract rule 2 requires "not exhausted").

## Boundary

This slice delivers the core ledger/repair semantics and local gates. Cross-run cost aggregation and adapter settlement reconciliation remain for later integration in T023/T013.
