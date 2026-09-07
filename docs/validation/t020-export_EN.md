# T020 PC-02 Export Validation Report

[简体中文](t020-export.md)

Date: 2026-09-08. Verdict: T020 COMPLETE.

## Scope

1. internal/snapshot/export.go
2. internal/snapshot/export_test.go
3. internal/evidence/delivery.go
4. internal/evidence/delivery_test.go

## AT-17 Coverage

1. Accepted snapshot binding: export requires acceptedSnapshotHash to match snapshot manifest hash, blocking candidate-as-accepted misuse.
2. Destination precondition: destination must be absent; existing directory or file is rejected.
3. Overlap blocking: fail-close when destination overlaps source/run/store.
4. Offline integrity: verifies object closure and exported file hashes; missing objects and tampering are rejected.
5. Secret blocking: paths matching secret/token/password/credential/.env patterns are rejected.
6. Interrupted write has no completion receipt: injected interruption does not produce a completed summary and destination remains absent.
7. Source unchanged: export is read-only for source/store; tests assert source bytes remain unchanged.

## Key Tests

1. TestExportAcceptedPackageRoundTrip
2. TestExportAcceptedPackageRejectsHashMismatchAndExistingDestination
3. TestExportAcceptedPackageRejectsOverlapSecretMissingAndTamper
4. TestExportAcceptedPackageInterruptedWriteHasNoCompletion
5. TestExportRecordCompletedRoundTrip
6. TestExportRecordCompletedRequiresEntriesAndManifest
7. TestExportRecordFailedRequiresErrorEvidence
8. TestExportRecordRejectsDestinationPreconditionAndHashMismatch

## Gates

1. go test ./internal/evidence ./internal/snapshot: passed.
2. go build ./...: passed.
3. go vet ./...: passed.
4. go test -count=1 ./...: passed.

## Boundary

1. This slice provides library-level export APIs only; CLI export command is planned for later slices.
2. Secret-path screening uses a fail-close keyword policy and blocks on any match.
3. Export completion does not change chain/task acceptance state and is not equivalent to external release publication.
4. The secret keyword policy may produce false positives (e.g. `tokenizer.txt`) or false negatives (e.g. `.npmrc`/`.pem`); it stays fail-close and should later merge with the secret-policy exclusion mechanism.
5. `destinationIdentityHash` binds exported entry content only, not the destination absolute path; the CLI layer must make this semantic explicit.
6. Exported directories use 0755 and hardlinks are written as independent files; use Restore rather than export when bit-level fidelity is required.
7. No dedicated tests yet for Windows short-path/junction aliases or secret-policy exclusions.

## Review Log (2026-09-08)

1. Removed the unused `NormalizeExportHashes` helper (dead-code cleanup).
2. Gates re-run and all passed.
