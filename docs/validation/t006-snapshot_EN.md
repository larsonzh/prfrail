# T006 Snapshot and Storage Substrate Validation Report

[简体中文](t006-snapshot.md)

Date: 2026-09-07. Verdict: `T006 COMPLETE`. This report covers snapshot capture, content-addressable storage (CAS), Git-free restoration, and AT-02 acceptance; chain execution and transaction journaling belong to T008/T009.

## Implemented Scope

- `types.go`: `snapshot-manifest` data model and strict validation supporting baseline and candidate modes; Unicode code-point sorting, directory closure enforcement, Windows reserved name rejection (CON/PRN/AUX/NUL/COM/LPT and suffixed variants), case collision detection, symlink escape defense, and domain-separated hashing with `proofrail:snapshot-manifest:1\n`.
- `store.go`: Content-Addressable Storage (CAS) with temporary-file sync, atomic rename, content deduplication, read-time hash integrity check and corruption isolation, storage byte quota enforcement, and deterministic GC dry-run (strictly retaining reachable and audit-locked objects, performing zero physical deletions).
- `capture.go`: Read-only source workspace traversal capturing uncommitted and untracked files while preserving source tree metadata and contents unchanged; support for default (`.git/**`, `tmp/**`, `.prfrail/**`) and user exclusions, file role assignment, hardlink group detection, environment probing (OS/Arch/case sensitivity/whitelisted env names), and concurrent modification detection (mtime/size verification).
- `restore.go`: Isolated workspace materialization with pre-flight reference closure validation (fail-close on missing objects); purely content-driven reconstruction of directories, regular files with executable/read-only permissions, hardlinks, and symlinks without calling Git or reading Git history.
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
- Concurrent modification: File size or mtime modification during capture detected and halted (`ErrConcurrentChange`).
- Quota exhaustion: Exceeding `maxBytesQuota` during write or capture immediately fails closed (`ErrQuotaExceeded`).
- GC dry-run: Reachable and locked objects strictly preserved; only expired orphaned objects listed; zero files deleted.
- Restore closure verification: Blocked on non-empty target, missing objects, or tampered hashes; restored output bit-for-bit identical.

## Dual-Platform Validation Record

```powershell
go test -v ./internal/snapshot
go vet ./internal/snapshot
go test -cover ./internal/snapshot
```

Windows 11 (NTFS): 18 of 19 tests passed; 1 symlink test skipped as expected due to missing unprivileged symlink capability; statement coverage 76.3%.
Ubuntu 24.04 (ext4): Executed via SSH in remote test environment; all 19 tests passed (including native symlink creation and resolution).

## Boundary

T006 establishes only the snapshot and storage substrate. It executes no Git commit/push, recovers no data from Git history, and performs no physical deletion during GC. Workspace write transactions and crash rollback points belong to T008.
