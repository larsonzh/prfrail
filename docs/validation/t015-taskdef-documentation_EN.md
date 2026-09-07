# T015 Taskdef Documentation Validation Report

[简体中文](t015-taskdef-documentation.md)

Date: 2026-09-08. Verdict: `T015 COMPLETE`. This slice delivers `harnesses/{generic,c,go}/` and `internal/taskdef/documentation.go` with locked template A checks, deterministic documentation impact-rule enforcement, generated-artifact freshness/anti-manual-edit guards, and C+Go task-level recovery aggregation.

## Scope

1. `internal/taskdef/documentation.go`
   - Adds documentation policy and impact-rule checks for `required/if-affected/optional/forbidden`.
   - Adds a locked template A contract (`template-a@1.0.0`) with hash binding.
   - Adds generated-artifact gates blocking `generated.self-approval`, `generated.manual-edit`, `generated.stale`, `generated.actor-missing`, and template-lock mismatches.
   - Adds task-level aggregation so any failed language scope requires whole-task recovery.
2. `internal/taskdef/documentation_test.go`
   - Covers impacted-doc requirements and missing-doc blocking.
   - Covers template A lock, manual generated edits, self-approval, and stale generated outputs.
   - Covers C+Go task-wide recovery semantics.
   - Verifies presence and template defaults in `harnesses/generic|c|go/harness.json`.
3. `harnesses/{generic,c,go}/harness.json`
   - Adds three external S1 harness definitions (generic/C/Go).
   - Each harness declares `generated-hook-template=template-a@1.0.0` and `allow-generated-hooks=false`.

## AT-11 Coverage

1. generic/C/Go: three external harness profiles are present and validated by tests.
2. Locked template A: generated artifacts are accepted only with `template-a@1.0.0` and the locked template hash.
3. Blocking cases:
   - impact rule matched but required documentation target missing.
   - required hook missing.
   - manual edits on generated artifacts, stale generated outputs, or self-approval.
   - generated artifacts with a missing generator or approver actor.
4. C+Go same-task aggregation: any language-scope failure triggers whole-task recovery.

## Tests

- `internal/taskdef/documentation_test.go`
  - `TestCheckDocumentationIfAffectedPassesWithFreshGeneratedArtifacts`
  - `TestCheckDocumentationBlocksMissingDocumentationAndHook`
  - `TestCheckDocumentationBlocksGeneratedViolations`
  - `TestCheckDocumentationRequiresImpactRulesWhenPolicyNeedsDocumentation`
  - `TestEvaluateTaskRecoveryForCAndGoScopes`
  - `TestHarnessProfilesProvideGenericCGoAndTemplateALock`

## Gates

1. `go test -count=1 ./internal/taskdef`: passed.
2. `go build ./...`: passed.
3. `go vet ./...`: passed.
4. `go test -count=1 ./...`: passed.
5. `2026-09-08` remote high-count race regression (Ubuntu VM `10.0.0.199`, Go 1.27.1 + gcc 13.3.0): `go test -race -count=20 ./internal/applier ./internal/console ./internal/taskdef` and `go test -race -count=1 ./...` passed on first pass (no proxy fallback needed).

## Boundary

1. This slice only delivers static constraints for generated scenario A and does not introduce arbitrary generated script B.
2. This slice does not introduce whois-specific terms or runtime semantics.
3. No commit and no push in this run.
