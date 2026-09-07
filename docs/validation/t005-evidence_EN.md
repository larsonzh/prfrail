# T005 Evidence-Layer Validation Report

[简体中文](t005-evidence.md)

Date: 2026-09-07. Verdict: `T005 COMPLETE`. This report covers only the evidence-layer portion of AT-04; journal, write transaction, rollback, and crash points belong to T008.

## Implemented Scope

- `canonical.go`: strict JSON reading; rejection of duplicate keys, invalid UTF-8, and unpaired surrogates; JCS UTF-16 key ordering; ECMAScript number rendering; and domain-separated SHA-256.
- `event.go`: `state-event` construction/decoding, version and field validation, hash/sequence continuity, and per-entity projected transition validation.
- `receipt.go`: `evidence-manifest` construction/decoding, immutable root digest, item uniqueness and fixed ordering, plus kind/media-type/redaction constraints.
- `verify.go`: rereads objects from a caller-supplied read-only content store and verifies reference closure, content digests, canonical JSON, event chains, and attempt/run binding. It returns only structured `passed|failed|incomplete` checks and frozen error codes.

## Fail-Closed Evidence

Tests cover the frozen canonical/state-event vectors plus duplicate keys, trailing values, invalid UTF-8, surrogates, negative zero, unsupported versions, unknown fields, manifest tampering/order/duplicates, missing objects, digest mismatches, noncanonical JSON, event gaps/reordering/corruption, cross-entity state discontinuity, manifest order differing from event sequence, run/attempt mismatches, and a nil reader. Missing capability produces `incomplete`; structural, digest, or chain conflicts produce `failed`.

## Validation Commands

```powershell
go test ./internal/evidence
go vet ./internal/evidence
go test -cover ./internal/evidence
```

Package tests and vet passed, with 79.5% statement coverage. The Windows race build was not run because this host has no `gcc`; `go test -race` stopped at the missing cgo compiler, after which ordinary tests passed. The final repository-wide gate results are recorded at task completion.

## Boundary

This implementation does not persist a complete `verification-report` or claim the full RFC section 16.9.26 identity/signature/review/promotion/snapshot/queue/lifecycle check set. Those checks are connected by their owning later tasks. T005 writes no state, object, or authorization and does not equate hash integrity with trusted identity.
