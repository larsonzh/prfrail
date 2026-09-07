# T014 Handoff Validation Report

[简体中文](t014-handoff.md)

Date: 2026-09-07. Verdict: `T014 COMPLETE`. This slice delivers the manual-handoff core model and input-policy guardrails, including handoff receipt semantics, WAITING_FOR_OPERATOR state transitions, and secure terminal-only handling for secret-direct inputs.

## Scope

1. `internal/chain/handoff.go`
   - Adds `HandoffPolicy` validation for `allowedTargets/inputPolicy/handoffTimeoutMs/returnActions/hooksAfterReturn`.
   - Adds `BuildHandoffReceiptRecord` to validate/build handoff receipts with operator identity, policy hash, manifest/diff hashes, lease/hook/evidence hashes, and canonical receipt hash.
   - Adds `OpenHandoffTransition` and `ReturnHandoffTransition` to model the WAITING_FOR_OPERATOR open/return loop.
2. `internal/chain/models.go`, `internal/chain/router.go`
   - Adds `manual-handoff` as a code step mode and binds it to mandatory `handoffPolicy` validation.
   - Extends the step router with a dedicated manual-handoff port, keeping it separated from managed-change-set, isolated-workspace, and hook execution paths.
3. `internal/console/handoff.go`
   - Adds handoff input policy enforcement for `structured` and `secret-direct`.
   - Enforces secret-direct fail-closed behavior unless the channel is terminal and secure input capability is present.
   - Converts secret values into hash evidence only; plaintext is never returned in structured outputs.
4. Tests
   - Adds `internal/chain/handoff_test.go` and `internal/console/handoff_test.go`.
   - Extends `internal/chain/state_test.go` and `internal/chain/router_test.go` for manual-handoff definition and routing coverage.

## AT-09 Coverage

1. Open transition is only allowed as `STEPS_RUNNING/RUNNING -> WAITING_FOR_OPERATOR/WAITING_FOR_OPERATOR`.
2. Return transitions cover all return actions:
   - `complete -> STEPS_RUNNING/RUNNING`
   - `abort -> FAILED/FAILED`
   - `request-agent -> FAILED/FAILED` (the repair flow then opens a new attempt and injects confirmed conclusions; never maps to CANCELLED to avoid poisoning the whole run)
3. Handoff receipt constraints enforced:
   - `recordedBy=system`, `operator=operator`, valid timestamps, and `closedAt > openedAt`.
   - `complete` requires non-empty `hookResultEvidence`.
   - `abort/request-agent` require `hookResultEvidence=[]`.
4. `AllowedReturnActions` is bound to `outcome`; disallowed outcomes are rejected.

## AT-10 Coverage

1. Undeclared prompt fields are rejected.
2. `secret-direct` returns `ErrSecureInputRequired` (pause/fail-close) on non-terminal channels or without secure-input capability.
3. Secret inputs are represented only by hash evidence and do not appear in structured output (validated with a canary secret).

## Tests

- `internal/chain/handoff_test.go`
  - `TestBuildHandoffReceiptRecordComplete`
  - `TestBuildHandoffReceiptRecordRejectsOutcomeRules`
  - `TestBuildHandoffReceiptRecordRejectsDisallowedOutcome`
  - `TestBuildHandoffReceiptRecordRejectsEmptyAllowedReturnActions`
  - `TestBuiltHandoffReceiptSerializesToSchemaShape`
  - `TestHandoffTransitions`
- `internal/chain/state_test.go`
  - `TestDefinitionValidationFourKindsAndModes` (manual-handoff positive/negative checks added)
- `internal/chain/router_test.go`
  - `TestStepRouterSeparatesChangeModesAndHooks` (manual-handoff routing assertion added)
- `internal/console/handoff_test.go`
  - `TestPrepareHandoffInputStructured`
  - `TestPrepareHandoffInputRejectsUndeclaredPrompt`
  - `TestPrepareHandoffInputRejectsUnknownChannel`
  - `TestPrepareHandoffInputSecretDirectRequiresSecureTerminal`
  - `TestPrepareHandoffInputSecretDirectDoesNotExposeSecret`

## Gates

1. `go test -count=1 ./internal/chain ./internal/console`: passed.
2. `go build ./...`: passed.
3. `go vet ./...`: passed.
4. `go test -count=1 ./...`: passed.
5. Ubuntu VM remote full validation (2026-09-08, temporary proxy chain):
   - `go build ./...`: passed.
   - `go vet ./...`: passed.
   - `go test -count=1 ./...`: passed.
   - `CGO_ENABLED=1 CC=gcc go test -count=1 -race ./...`: passed.
   - `CGO_ENABLED=1 CC=clang go test -count=1 -race ./internal/chain ./internal/console`: passed.
   - Result summary: `RESULT build=0 vet=0 test=0 race_gcc=0 race_clang=0 race_all=0`.
   - After validation, the remote temporary proxy wrapper script and session proxy variables were cleaned up to avoid residual proxy configuration.

## Review Fixes

A 2026-09-08 review found and fixed the following issues:

1. Wrong request-agent return mapping: it previously mapped to `CANCELLED/CANCELLED`, but CANCELLED is terminal and would poison the whole run; the contract requires request-agent to open a new attempt and continue. It now maps to `FAILED/FAILED` (the only state-machine precondition for the repair flow to open a new attempt), with tests and this report updated accordingly.
2. An empty AllowedReturnActions silently skipped outcome binding checks (fail-open): the contract requires the outcome to fall within the policy `returnActions` set, and the schema enforces `minItems: 1`. The session now requires a valid non-empty set; empty or invalid sets are rejected, with a new `TestBuildHandoffReceiptRecordRejectsEmptyAllowedReturnActions`.
3. The input channel enum was not validated: unknown channel values passed silently for structured input. Non-terminal/model channels are now rejected fail-closed, with a new `TestPrepareHandoffInputRejectsUnknownChannel`.
4. Added regression coverage: `TestBuiltHandoffReceiptSerializesToSchemaShape` (abort receipts must serialize `hookResultEvidence` as `[]` rather than `null`, and keys must be camelCase); `state_test` now covers manual-handoff-without-policy and managed-with-policy negatives.

## Boundary

1. This slice does not introduce CLI-level handoff interaction commands yet; it focuses on chain semantics, receipt constraints, and secure input policy guards.
2. Runtime-dependent AT-09 items (lease conflicts, operator absence/disconnection, crash replay, and re-running gates on return) belong to later engine wiring with guard/lease and gates; real terminal secure input capability belongs to T016 CLI/TUI.
3. Validation in this run is offline unit/gate validation; interactive authorized-host smoke remains a separate planned activity.