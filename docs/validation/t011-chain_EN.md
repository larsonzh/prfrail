# T011 Review/Publish Validation Report

[简体中文](t011-chain.md)

Date: 2026-09-07. Verdict: `T011 COMPLETE`. This task delivers candidate freeze, independent review, promotion receipt binding, and task acceptance wiring. After review, two defects were fixed and all gates were re-validated.

## Reviewed and Fixed Defects

1. Receipt serialization schema conformance: JSON tags were added to `ReviewReason` and `WaiverAuthorization` to prevent upper-case keys such as `Code`; empty `message` is omitted to satisfy schema `minLength` constraints.
2. `REPAIR_PENDING` resume guard: `REPAIR_PENDING/WAITING_FOR_OPERATOR` cannot continue within the same attempt; the chain is paused and requires a new attempt or operator-driven repair flow.

## Key Implementation

1. `internal/chain/review.go`: review receipt construction, JCS hashing, and approve/reject/waive validation.
2. `internal/chain/publish.go`: promotion receipt construction, completed/failed/uncertain semantics, and evidence validation.
3. `internal/chain/engine.go`: replace direct acceptance after `REVIEW_PENDING` with candidate->review->promotion flow; reject/invalid review/incomplete promotion enters `REPAIR_PENDING` and pauses the chain.
4. `internal/chain/models.go`: add candidate/review/publish ports and models; complete JSON field tags for review structures.

## Test Coverage (AT-05)

1. Rejection path: `TestReviewRejectMovesTaskToRepairPending`.
2. Self-approval blocking: `TestRejectSelfReview`.
3. Expired waiver blocking: `TestRejectExpiredOrWrongCandidateWaiver/expired_waiver`.
4. Wrong-candidate waiver blocking: `TestRejectExpiredOrWrongCandidateWaiver/wrong_candidate`.
5. Exactly-once publish after valid independent approval: `TestValidIndependentApprovalPublishesOnce`.
6. Regression tests for review fixes:
   - Schema key casing and empty message: `TestBuiltReceiptsSerializeToSchemaShape`.
   - Repair-pending resume block: `TestRepairPendingCannotResumeWithoutNewAttempt`.

## Gate Re-Validation

- Targeted re-check:
  `go test -count=1 ./internal/chain -run "TestBuiltReceiptsSerializeToSchemaShape|TestRepairPendingCannotResumeWithoutNewAttempt|TestRejectSelfReview|TestRejectExpiredOrWrongCandidateWaiver|TestValidIndependentApprovalPublishesOnce"` passed.
- Full gates: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` passed.

## Boundary

Current behavior guarantees single publish within one `Run()` execution. Cross-crash publish idempotency and receipt-store fencing remain a later-scope hardening item outside T011.
