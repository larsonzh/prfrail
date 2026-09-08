# T024 PC-06 Lifecycle Validation Report

[中文](t024-lifecycle.md)

Date: 2026-09-08. Verdict: `T024 COMPLETE`.

## Implementation Scope

1. `internal/snapshot/backup.go` / `backup_test.go`: complete object/event-segment/reference-root closure backup, new-store read-back verification, and staging-before-atomic-publication.
2. `internal/evidence/disposition.go` / `disposition_test.go`: release/backup/retirement lifecycle records, strict decoding and hash validation, release inventory, retirement, and retention rules.
3. `docs/DEV_PLAN*.md`, `docs/OPERATIONS*.md`, and `docs/INSTALLATION_PLAN*.md`: synchronized T024 completion while preserving the pending T017 release boundary.

## AT-21 Coverage

1. Active-writer and unknown-format blocking: backup requires `WritersStopped=true`, non-empty stop evidence, and `SourceSchemaVersion=1.0.0`; unknown versions are neither migrated nor written.
2. Complete closure: object, event-segment, and reference-root sets must be non-empty, uniquely hashed, and valid; missing objects/events/roots fail closed before publication.
3. New-store restore verification: objects are copied into an independent `Store`; events and references are written by hash and read back. Restore evidence and a completed manifest appear only after every check passes.
4. Original preservation and interruption safety: tests compare source objects before/after backup; injected failure leaves only removable staging and publishes neither destination nor completion manifest.
5. Lifecycle records: completed backup requires `closureStatus=complete`, excluded secrets, object/root/restore evidence, and a manifest hash; failed/uncertain records cannot carry a completion manifest.
6. Retention by default: retirement construction defaults to `retentionDisposition=retain` and `sharedToolsDisposition=preserved`, preventing SessionBridge or user-toolchain deletion.
7. Deletion and audit conflicts: `delete-authorized` requires separate human authorization; retention conflicts require an independent human decision hash.
8. Release inventory and offline revocation: inventory is path-sorted, content-stable, and rejects symlinks; release requires signature receipt, SBOM/license/checksum/support data, while offline revocation is represented only as `unknown`.

## Gate Results

1. `go test -count=1 ./internal/snapshot ./internal/evidence`: passed.
2. `go build ./...`: passed.
3. `go vet ./...`: passed.
4. `go test -count=1 ./...`: passed (all packages ok).

## Test List (8 Cases)

1. `internal/snapshot/backup_test.go`: `TestCreateBackupVerifiesNewStoreAndPreservesSource`
2. `internal/snapshot/backup_test.go`: `TestCreateBackupRejectsActiveWriterUnknownSchemaAndMissingClosure`
3. `internal/snapshot/backup_test.go`: `TestCreateBackupInterruptedHasNoCompletedManifest`
4. `internal/evidence/disposition_test.go`: `TestLifecycleBackupRoundTrip`
5. `internal/evidence/disposition_test.go`: `TestRetirementDefaultsRetainAndPreserveSharedTools`
6. `internal/evidence/disposition_test.go`: `TestRetirementDeletionAndConflictRequireHumanEvidence`
7. `internal/evidence/disposition_test.go`: `TestReleaseRequiresSignatureAndRepresentsOfflineRevocationAsUnknown`
8. `internal/evidence/disposition_test.go`: `TestBuildReleaseInventoryIsDeterministicAndRejectsLinks`

## Boundaries

1. T024 supplies library-level capabilities, not an installer, signing keys, release action, or copyable backup/uninstall CLI. Those remain governed by T017 and explicit release authorization.
2. Unknown schema versions are write-rejected while preserving the original store; explicit cross-version migration belongs to S2.
3. Logical deletion does not guarantee secure erasure from SSDs, snapshots, or external copies.
4. No `git commit`, `git push`, signing, or release action was performed.

## Review Notes

1. Resumed the previously incomplete `disposition.go`; fixed its package-local slice-cloning compile gap while retaining the contract implementation.
2. The VS Code test adapter did not discover the new Go test files, so repository-native `go test` was used and passed.
3. Related bilingual documents are synchronized; the user-edited T023 bilingual validation reports were not modified.