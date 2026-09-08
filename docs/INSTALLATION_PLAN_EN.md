# ProofRail Installation and Deployment Plan

[简体中文](INSTALLATION_PLAN.md)

Date: 2026-09-07. Status: planning, not an installation guide. No formal release package, stable installation command, or supported production deployment exists yet. This document records open decisions and pre-release gates so target designs are not presented as current capability.

## 1. Currently Available Path

Only developer source builds and tests are currently supported, requiring Go 1.22+, repository sources, and the target project's own toolchain. `cmd/prfrail` remains a development entry point and does not constitute an installed product. See [Operations](OPERATIONS_EN.md) for accurate current commands.

## 2. Target Deployment Model

The S1 target is a local, single-user, single-writer CLI/TUI with no persistent cloud service. The ProofRail binary, state/store, isolated run workspaces, and read-only source have explicit boundaries. SessionBridge is an optional adapter, not a core installation dependency. Compilers, interpreters, and test tools belong to harness declarations and are never silently installed by ProofRail.

## 3. Open Decisions

| Decision | Candidate | Freeze condition/owner | Status |
|---|---|---|---|
| S1 artifact | Windows amd64 archive; evaluate winget later | T017 native acceptance, SBOM, checksum, and signing | TBD |
| Linux support | S1 core CI/preview; production support candidate in S2 | Native compatibility, path/permission/upgrade tests | TBD |
| Install directory and PATH | Per-user standalone or package-manager directory | No admin requirement; side-by-side upgrade and rollback | TBD |
| Default state/store location | User data directory or explicit workspace configuration | Verifiable non-overlap, permissions, and backup policy | TBD |
| Configuration search order | Explicit argument, workspace, user profile | T016 freezes CLI and config-explain behavior | TBD |
| Signature/trust root | Offline Ed25519 direction approved; identity/key unresolved | ADR, key governance, revocation, and rotation rehearsal | TBD |
| Update mechanism | Manual verified download; package manager later | Support matrix, schema compatibility, backup/recovery gates | TBD |
| Uninstall/data retention | Separate binary removal from run/store disposition | T024 lifecycle record and explicit deletion authorization | Library rules frozen; installation entrypoint awaits T017 |
| VS Code integration | Optional extension/SessionBridge adapter | CLI-independent loop passes; no authority beyond core API | TBD |

Until ADRs and executable acceptance evidence exist, do not choose default directories, publish one-click installers, claim automatic updates, or require administrator privileges.

## 4. Planned User Path

1. Select an OS/architecture-matching artifact from a trusted release location.
2. Independently verify checksum, signature, signer identity, support status, and revocation freshness.
3. Extract/install into a separate directory without overwriting a running host; run version and read-only self-checks.
4. Initialize explicit state/store and workspace configuration, proving no overlap with source or run workspaces.
5. Run side-effect-free static preview before separately authorizing capability probes and execution.
6. Before upgrade, stop writers and make a complete recoverable backup to a new store. Preserve the old binary and store on failure.
7. Uninstall retains evidence and configuration by default; data deletion is separately authorized.

These remain acceptance targets, not current copy-and-run instructions. T024 now supplies library-level complete backup, new-store verification, retention-by-default, and explicit deletion-authorization rules; T017 still owns installation commands, release filenames, and download locations.

## 5. Pre-Release Installation Gates

- Native Windows 11 amd64 install, first run, upgrade, rollback, and uninstall pass. Linux claims match the current support matrix only.
- Artifacts support offline checksum, signature, SBOM, license, and exact-version verification. Unknown or stale revocation information never displays fully verified.
- A clean user machine runs ProofRail without Go. External tools required by project tasks are listed and never automatically installed.
- Installation neither modifies target source/Git nor removes SessionBridge/user toolchains, and secrets never enter configuration, logs, or backups.
- State/store backup covers the complete reference closure and has been restored into a new directory.
- Every installation command in README and operations docs has been exercised against the corresponding artifact on a supported platform.

## 6. Documentation Completion

T016 has frozen CLI, configuration, and current source-run behavior; T024 has frozen library-level backup, pre-upgrade blocking, retirement, and data-retention rules. T017 still needs to supply candidate artifacts, CI, signing, and the release support matrix. Until T017 completes, this remains a plan and [Operations](OPERATIONS_EN.md) continues to distinguish executable behavior from future design.