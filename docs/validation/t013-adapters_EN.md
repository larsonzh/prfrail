# T013 Adapters Validation Report

[简体中文](t013-adapters.md)

Date: 2026-09-07. Verdict: `T013 COMPLETE`. This slice implements file-queue and sessbridge adapter entry points, unifies request/claim/result envelopes and dispatch/takeover receipt hashing rules, and validates AT-08 offline critical scenarios.

## Scope

1. `internal/adapters/filequeue.go`:
   - Implements validation, digesting, encoding, and decoding for request/claim/result envelopes.
   - Implements validation, digesting, encoding, and decoding for dispatch/takeover receipts.
   - Implements queue directories (requests/inflight/results/archive) and no-replace atomic writes.
   - Implements claim acquisition, generation monotonicity, fencing token checks, result submission, and idempotent conflict detection.
2. `internal/adapters/sessbridge.go`:
   - Defines a consumer-owned silent client port.
   - Maps sessbridge dispatch responses to adapter dispatch receipts.
   - Covers accepted/rejected/uncertain outcomes with requestId binding checks and transport error evidence fallback.
3. `internal/adapters/{filequeue_test,sessbridge_test}.go`:
   - Adds regression coverage for file-queue and sessbridge critical paths.

## AT-08 Coverage

1. Envelope parity: request/claim/result follow the same `adapter-envelope.schema.json` field set and per-message hash domains.
2. Repeated-apply prevention: same `requestId` + same payload is idempotent; same `requestId` + different payload is blocked as conflict.
3. Claim and fencing:
   - Initial claim requires the request to move into inflight.
   - generation must increase monotonically.
   - Result submission is rejected when `(claimId,generation)` does not match the current token.
4. Takeover constraint: changing claimId on later generations requires takeoverReceiptHash.
5. Sessbridge boundary:
   - `transport=sessbridge` requires non-empty transportRequestId.
   - busy/timeout/poll_timeout/lm_api_unavailable/unknown statuses map to uncertain.
   - Response `requestId` mismatch fail-closes to uncertain.
6. No GUI fallback: implementation only consumes a silent client port and does not add visible/GUI auto fallback.

## Tests

- `internal/adapters/filequeue_test.go`
  - `TestFileQueueDispatchClaimResultFlow`
  - `TestFileQueueDispatchIsIdempotentForSamePayload`
  - `TestFileQueueDispatchRejectsDifferentPayloadForSameRequestID`
  - `TestFileQueueClaimTakeoverRequiresReceipt`
  - `TestFileQueueSubmitResultRejectsFencingMismatch`
  - `TestFileQueueSubmitResultRejectsWrongRequestHash`
  - `TestWriteCanonicalLineNoReplaceRejectsExistingFile`
  - `TestFileQueueClaimGeneration1IsImmutable`
    - `TestFileQueueClaimGeneration1RejectsTakeoverReceiptHash`
    - `TestFileQueueClaimRenewRejectsTakeoverReceiptHash`
- `internal/adapters/sessbridge_test.go`
  - `TestSessbridgeDispatchAccepted`
  - `TestSessbridgeDispatchRejected`
  - `TestSessbridgeDispatchBusyReturnsUncertain`
  - `TestSessbridgeDispatchRequestIDMismatchReturnsUncertain`
  - `TestSessbridgeDispatchClientErrorGeneratesTransportID`
  - `TestSessbridgeDispatchRequiresClient`

## Gates

1. `go test -count=1 ./internal/adapters`: passed.
2. `go build ./...`: passed.
3. `go vet ./...`: passed.
4. `go test -count=1 ./...`: passed.

## Review Fixes

A 2026-09-07 review found and fixed the following issues:

1. no-replace publish did not hold: `os.Rename` replaces an existing destination on both Windows and Unix (violating "existing official files must never be overwritten"), which made the `fs.ErrExist` fallback unreachable. Publish now uses hard-link no-replace (`os.Link` fails when the destination exists) with best-effort temporary-file cleanup; the request→inflight move now uses link+remove with rollback on failure.
2. claim renewal semantics were incomplete: by contract, initial claims and same-claimId renewals require `takeoverReceiptHash=null`, but the previous implementation accepted non-null values. The `Claim` path now enforces this rule and only allows/needs `takeoverReceiptHash` when claimId changes.
3. Added regression tests: `TestWriteCanonicalLineNoReplaceRejectsExistingFile`, `TestFileQueueClaimGeneration1IsImmutable`, `TestFileQueueClaimGeneration1RejectsTakeoverReceiptHash`, and `TestFileQueueClaimRenewRejectsTakeoverReceiptHash`.

## Boundary

1. This slice does not modify the `sessbridge` repository.
2. This run validates offline contracts and unit tests; authorized-host smoke is executed in an authorized environment per plan, without changing this slice conclusion.