# T019 PC-01 Preview Validation Report

[简体中文](t019-preview.md)

Date: 2026-09-08. Verdict: `T019 COMPLETE`.

## Scope

1. `internal/taskdef/preview.go`
2. `internal/console/preview.go`
3. `internal/taskdef/preview_test.go`
4. `internal/console/cli_test.go` (preview test cases)
5. `examples/no-ai-preview/README.md`

## AT-16 Coverage

1. No-AI static preview: `mode=offline-read-only`; no step/hook execution.
2. Zero call counters: command/network/model/credential/versionProbe are all 0.
3. Visible unknowns: explicit reason and nextAction, no fabricated probe result.
4. Source/store unchanged: preview reads config and computes hashes only; no run workspace creation.
5. JSON/text parity: both outputs report the same `previewHash` and outcome facts.

## Gates

1. `go test ./internal/taskdef ./internal/console`: passed.
2. `go build ./...`: passed.
3. `go vet ./...`: passed.
4. `go test -count=1 ./...`: passed.

## Key Command Evidence

```powershell
go test ./internal/taskdef ./internal/console
go build ./...
go vet ./...
go test -count=1 ./...
```

## Boundary

1. No command probes, network probes, model calls, or automatic installation.
2. Executable steps produce a `blocked` preview only and do not trigger run.
3. Real authorization/export/lifecycle flows belong to T020+.
4. `reviewRequired` is always `true`: the preview is intentionally conservative and cannot express review-free steps today.
5. `blockingGateIds` reflects all step `hooks` without distinguishing blocking from non-blocking gates.

## Review Log (2026-09-08)

1. Aligned with `plan-preview.schema.json`: `networkRequirements[].host` now enforces `^[A-Za-z0-9.-]+$` plus a 253-length limit; `budget.pricingSource` now enforces a 1–256 length limit; added `TestPlanPreviewRecordRejectsInvalidNetworkHost` and `TestPlanPreviewRecordRejectsLongPricingSource` rejection cases.
2. Gates re-run and all passed.
