# T006 Snapshot and Storage Substrate Validation Report

[简体中文](t006-snapshot.md)

Date: 2026-09-07. Verdict: `T006 COMPLETE`. This report covers snapshot capture, content-addressable storage (CAS), Git-free restoration, and AT-02 acceptance; chain execution and transaction journaling belong to T008/T009.

## Implemented Scope

- `types.go`: `snapshot-manifest` data model and strict validation supporting baseline and candidate modes; Unicode code-point sorting, NFC-equivalent and case collision detection, directory closure enforcement, Windows reserved-name rejection, symlink escape defense, and domain-separated hashing with `proofrail:snapshot-manifest:1\n`.
- `store.go`: Content-Addressable Storage (CAS) with strict digest validation, temporary-file sync and same-volume hard-link no-replace publication, deduplication, read-time integrity checks, and corruption isolation. Existing corrupt objects cannot be overwritten. Byte quotas and deterministic GC dry-run retain reachable and audit-locked objects without physical deletion.
- `capture.go`: Read-only source workspace traversal capturing uncommitted and untracked files while preserving source tree metadata and contents unchanged; support for default (`.git/**`, `tmp/**`, `.prfrail/**`) and user exclusions, file role assignment, hardlink group detection, environment probing (OS/Arch/case sensitivity/whitelisted env names), and concurrent modification detection (mtime/size verification).
- `restore.go`: Isolated workspace materialization with pre-flight reference closure validation and same-parent staging followed by atomic publication. Hardlink or permission restoration failure neither degrades to a copy nor leaves a partial target. Reconstruction remains content-only and independent of Git history.
- `reparse_windows.go` / `reparse_other.go`: Platform-specific detection of reparse points and directory junctions; non-standard links and unportable reparse points fail closed.

## Acceptance Matrix (AT-02)

- Uncommitted file tree capture: Verified inclusion of modified, untracked, and lock files, generating schema-valid manifests and content digests.
- Source tree unchanged: Source workspace file sizes, mtimes, contents, and permissions matched identically before and after capture.
- Reserved names: CON, PRN, AUX, NUL, COM1-9, LPT1-9, and their extended names were strictly rejected (`ErrReservedPath`).
- Case collision: Conflicting paths with differing case were detected and rejected (`ErrCaseCollision`).
- Symbolic links: Relative workspace symlinks captured and materialized; escaping targets rejected. Skipped on unprivileged Windows and fully verified on Ubuntu.
- Reparse points/junctions: Non-standard reparse points and junctions detected and blocked (`ErrUnsupportedLink`).
- Hardlinks: Multiple files sharing inode/file index grouped under `hardlinkGroup`, materialized as hardlinks during restore.
- Long paths: Valid paths up to 1024 bytes supported; oversize paths rejected (`ErrInvalidPath`).
- Concurrent modification: A deterministic post-read mutation test proves that file size or mtime changes during capture halt with `ErrConcurrentChange`.
- CAS boundary: Invalid digests are rejected before path construction; a corrupt object at the expected digest path cannot be overwritten or used to escape the object directory.
- Atomic restoration: Complete materialization is published from a sibling staging directory; inability to restore hardlinks or permissions losslessly fails closed.
- Quota exhaustion: Exceeding `maxBytesQuota` during write or capture immediately fails closed (`ErrQuotaExceeded`).
- GC dry-run: Reachable and locked objects strictly preserved; only expired orphaned objects listed; zero files deleted.
- Restore closure verification: Blocked on non-empty target, missing objects, or tampered hashes; restored output bit-for-bit identical.

## Dual-Platform Validation Record

```powershell
go test -v ./internal/snapshot
go vet ./internal/snapshot
go test -cover ./internal/snapshot
```

Pre-audit baseline: Windows 11 (NTFS) passed 18 of 19 tests with one expected unprivileged-symlink skip and 76.3% statement coverage; Ubuntu 24.04 (ext4) passed all 19 tests. The audit added a CAS immutability counterexample and made concurrent mutation deterministic; final results are recorded with the T007 combined regression.

## Boundary

T006 establishes only the snapshot and storage substrate. It executes no Git commit/push, recovers no data from Git history, and performs no physical deletion during GC. Workspace write transactions and crash rollback points belong to T008.
