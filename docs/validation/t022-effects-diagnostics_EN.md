# T022 PC-04 Effects and Recovery Diagnosis Validation Report

[中文](t022-effects-diagnostics.md)

Date: 2026-09-08. Verdict: `T022 COMPLETE`.

## Implementation Scope

1. `internal/gates/effects.go` (side-effect classification/observation model, S1 admission and recovery action)
2. `internal/gates/effects_test.go`
3. `internal/evidence/diagnostics.go` (recovery-diagnosis record model and validation)
4. `internal/evidence/diagnostics_test.go`
5. `internal/chain/recover_test.go` (chain recovery-plan integration tests)
6. `docs/validation/t022-effects-diagnostics.md` / `t022-effects-diagnostics_EN.md` (this report)

## AT-19 Coverage

1. Read-only / local-discardable / external-write separation: `EffectClass` matches `effect-record.schema.json`; classifications bind `policyHash` and `subjectDefinitionHash` from trusted policy — model self-report is never accepted.
2. S1 external-write denial: `AdmitEffect` denies `external-write` and unknown classes (fail closed); only `read-only`/`local-discardable` are admitted. Command names or dry-run wording never influence the decision.
3. Unprovable limits are rejected: validation forces `external-write` classifications to claim only `reconcile-only`/`compensation-only` recovery, forbidding `none-required`/`discard-workspace`.
4. Unknown effects are not blindly retried: observations with `outcome=uncertain/failed` must carry `errorEvidence`; `ObservationRecoveryAction` returns `reconcile-only` for every non-completed outcome.
5. Diagnosis is read-only and redacted: `RecoveryDiagnosis` enforces `mode=read-only-redacted` and `mutationsPerformed=false`; `clear`/`blocked` oneOf evidence rules follow the schema.
6. Recovery plans never delete locks or logs: `TestRecoverDoesNotMutateEventLog` proves `Recover` neither deletes nor rewrites existing events; `TestRecoverProducesBlockedDiagnosisPlan` proves the plan anchors the last accepted snapshot and presents `blocked` with `inspect`/`controlled-stop` actions.

## Gates

1. `go test -count=1 ./internal/gates ./internal/evidence ./internal/chain`: passed.
2. `go build ./...`: passed.
3. `go vet ./...`: passed.
4. `go test -count=1 ./...`: passed (all packages ok).
5. Editor diagnostics: no errors in new/changed Go files.

## Key Command Evidence

```powershell
gofmt -w internal/gates/effects.go internal/gates/effects_test.go internal/evidence/diagnostics.go internal/evidence/diagnostics_test.go internal/chain/recover_test.go
go test -count=1 ./internal/gates ./internal/evidence ./internal/chain
go build ./...
go vet ./...
go test -count=1 ./...
```

## Boundaries

1. effects/diagnostics are library-level capabilities with no CLI entrypoint yet; external-write runners, idempotent external APIs and compensation protocols belong to the S2 ADR scope per RFC and are explicitly denied in S1.
2. Observation records are now wired in the gate runner assessment path: when a request includes an effect scope, `internal/gates/runner.go` builds and validates classification/observation, then binds `effectObservationHash` and `effectRecoveryAction` into the hook result. T023/T024 still own broader product-level orchestration and presentation entrypoints.
3. Managed-process evidence in a diagnosis (`processes`) is filled by callers from guard stop evidence; diagnosis never starts or stops processes itself.
4. Rollback covers controlled files only; network publishes, email, database writes, package uploads and device operations are never claimed as "rolled back".

## Review Notes (2026-09-08)

1. First-round compile fix: `StateEvent` contains non-comparable fields, so `TestRecoverDoesNotMutateEventLog` compares event snapshots with `reflect.DeepEqual`.
2. All gates re-run and passed.
3. A later compile regression exposed a test-signature drift after `runner.assess` gained the classification parameter (WrongArgCount). `internal/gates/runner_test.go` (`assessForTest`) now passes the extra classification argument (`nil`), and both `go test -count=1 ./internal/gates` and full repository gates pass again.
4. To eliminate implicit allow paths when effect metadata is omitted, hook policy now requires explicit `effectClass` (missing value fails closed with `ErrInvalidHook`), and `buildEffectScope` no longer applies implicit default classes. Added tests: `TestPrepareRejectsMissingEffectClass`, `TestRunnerDeniesExternalWriteEffectWithoutExecuting`, and `TestChainPortDeniesExternalWriteHook`.
