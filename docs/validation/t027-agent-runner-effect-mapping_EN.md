# T027 AgentRunner Effect Mapping Prerequisite Validation Report

[中文](t027-agent-runner-effect-mapping.md)

Date: 2026-09-13. Verdict: `PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`.

## Implementation Scope

1. Froze the v1 canonical mapping body as `{"local-process":"local-discardable","workspace-write":"local-discardable"}`. The domain is `proofrail:agent-runner-effect-mapping:1\n`; its SHA-256 is `sha256:2c9609ee374bfc79bd0ed997aea6cfe4865c82ed0216d1a4a288139e0b0c92d4`.
2. `agent-runner-request` must bind `effectMappingVersion="1"` and that `effectMappingHash`. The mapping is a fixed implementation contract and introduces no runtime mapping record.
3. Schema, negative goldens, the canonical vector, and the semantic validator cover missing, unknown, or tampered bindings, unknown operations, grant-class under-coverage, and grant-class supersets.
4. Go request construction and strict decoding validate the same fixed binding. Composite admission maps every `allowedEffects` operation to an authorization class and validates grant coverage before request replay identity is recorded or dispatch can occur.

## Fail-Closed Results

- `local-process` and `workspace-write` map only to `local-discardable`.
- Missing, unknown, or tampered mapping versions or hashes and unknown operations are rejected during request-record validation.
- A valid authorization grant that lacks a mapped class is rejected during authorization admission. A healthy request using the same replay index is then admitted for its first dispatch, proving that rejection did not consume request identity.
- A grant may contain additional classes such as `read-only` and `external-write`; it passes when it includes every required `local-discardable` class.

## Validation Results

1. `node tools/contracts/contracts.test.js`: 3/3 passed; 36 Schemas, 123 fixtures, and 4 canonical vectors.
2. `go test -count=1 -json ./internal/adapters`: 124 passed, 0 failed, 0 skipped.
3. `go build ./...`: passed.
4. `go vet ./...`: passed.
5. `go test -count=1 -json ./...`: 393 passed, 0 failed, and 2 skipped (`TestApplyRejectsSymlinkParent` and `TestCaptureSymlinkValid`) across 13 tested packages; 3 additional packages had no tests.
6. `gofmt` on the changed Go files: passed.

## Boundaries and Remaining Blockers

1. This slice neither implements nor validates real tool/network enforcement; the pinned candidate still cannot dispatch.
2. The durable dispatch/completion ledger and terminal-receipt convergence are not wired to the AgentRunner port.
3. Resume continuity validation across prior completion, session, and workspace is not implemented.
4. No real model CLI was invoked, SecretStore was not read, no network access occurred, and no commit/push/publish action was performed.
5. T027 remains `BLOCKED / NOT IMPLEMENTED`; this report does not claim AT-23.
