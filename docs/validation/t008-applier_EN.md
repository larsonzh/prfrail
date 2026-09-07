# T008 Managed Change-Set and Transactional Applier Validation Report

[简体中文](t008-applier.md)

Date: 2026-09-07. Verdict: `T008 COMPLETE`. This task completes the RFC §16.9.25 managed change-set checker, fenced journal transaction, whole-group rollback, and AT-04 crash-recovery acceptance.

## Twelve Slices

1. Freeze the change-set, operation, assertion, marker, and pure in-memory write-plan models.
2. Strictly decode JSON, recompute the JCS/domain-separated changeSetHash, and bind the task attempt, parent snapshot, and before/after manifests.
3. Validate managed relative paths, regular-file targets, contiguous unique sequence/operationId values, and the 1–10000 operation bound.
4. Read all targets once and simulate only in memory in array order; repeated operations on one path bind each beforeHash to the prior simulated result.
5. Implement the closed create-file/delete-file/replace-exact/insert-before/insert-after operation set.
6. Enforce UTF-8, no NUL, internal LF, explicit create line endings, and preserve-existing for existing text without guessed conversion.
7. Enforce before/after assertions, exact match counts, and unique insert markers with 0→1 occurrence semantics; every failure performs zero writes.
8. Before workspace mutation, durably journal changeSetHash, before/after byte digests, modes, and planned created directories; reread blobs by digest.
9. Write each target through a same-directory temporary file, file fsync, and platform atomic replacement; Unix syncs the parent directory and Windows uses MoveFileExW REPLACE_EXISTING|WRITE_THROUGH.
10. Revalidate the T007 fencing token before every journal and workspace write; external changes, symlink parents, and stale writers fail closed.
11. Roll back the whole group in reverse order on ordinary failures; simulated crashes retain the journal, and restart recovery mutates only content equal to a known before/after state. Unknown content becomes uncertain and is not overwritten.
12. Cover journal-durable, before/after every item write, and after-full-verification crash points, then verify natively on Windows 11/NTFS and Ubuntu 24.04/ext4.

## Validation

- Windows: `go build ./...`, `go vet ./...`, and `go test -count=1 ./...` passed.
- Checker: all five operations, same-path sequential chains, strict JSON, digest/sequence/managed-path/UTF-8/NUL/assertion/marker counterexamples, and zero writes on failure passed.
- Applier: successful commit, ordinary-failure whole-group rollback, journal-durable, each before-write/after-write, after-verify, idempotent recovery, unknown external content, corrupt journal/blob, wrong token, mid-transaction fencing expiry, and symlink-parent counterexamples passed.
- Ubuntu: Go 1.27.1 on Ubuntu 24.04/ext4 passed offline vendored `go build ./...`, `go vet ./...`, `go test -count=1 ./...`, and `go test -race -count=1 ./internal/taskdef ./internal/applier`.
- `git diff --check` passed. Temporary archives, vendor content, and the remote validation directory were removed; repository `tmp/` contains only `.gitkeep`.

## Boundary

T008 does not schedule tasks or steps, produce chain projections or receipts, run gates, or promote candidate output to an accepted snapshot; T009 and later tasks retain those responsibilities. Windows `MoveFileExW(...WRITE_THROUGH)` supplies the available platform replacement durability but is not claimed equivalent to Unix directory fsync under power loss; every uncertain recovery remains fail closed.
