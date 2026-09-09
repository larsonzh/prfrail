# ProofRail Installation and Deployment Plan

[简体中文](INSTALLATION_PLAN.md)

Date: 2026-09-09. Status: the S1 installation model is decided, while formal release remains pending. No formal release package or supported production deployment exists yet. This document records the frozen model, open decisions, and pre-release gates; see the [portable ZIP installation guide](INSTALLATION_EN.md) for operational steps.

## 1. Currently Available Path

Only developer source builds and tests are currently supported, requiring Go 1.22+, repository sources, and the target project's own toolchain. `cmd/prfrail` remains a development entry point and does not constitute an installed product. See [Operations](OPERATIONS_EN.md) for accurate current commands.

## 2. Target Deployment Model

The S1 target is a local, single-user, single-writer CLI/TUI with no persistent cloud service. The ProofRail binary, state/store, isolated run workspaces, and read-only source have explicit boundaries. SessionBridge is not an installation dependency for the core offline CLI, but it is the required adapter for connecting to Copilot Chat through VS Code. Compilers, interpreters, and test tools belong to harness declarations and are never silently installed by ProofRail. The current candidate delivers only a line-oriented CLI; the complete TUI remains future implementation.

## 3. Open Decisions

| Decision | Candidate | Freeze condition/owner | Status |
|---|---|---|---|
| S1 artifact | Manually extracted Windows amd64 portable ZIP; winget may be evaluated later | T017 native acceptance, fixed commit, SBOM, license manifest, and checksums | Owner approved 2026-09-09 |
| Linux support | S1 core CI/preview; production support candidate in S2 | Native compatibility, path/permission/upgrade tests | S1 scope decided |
| Install directory and PATH | User-selected independent directory; invoke by full path without changing PATH, registry, or system directories | No admin requirement; side-by-side upgrade and rollback | Owner approved 2026-09-09 |
| Default state/store location | User data directory or explicit workspace configuration | Verifiable non-overlap, permissions, and backup policy | TBD |
| Configuration search order | Explicit argument, workspace, user profile | T016 freezes CLI and config-explain behavior | TBD |
| Release trust | S1 uses fixed commit, GitHub CI, SHA256SUMS, SBOM and license manifest | Checksums prove only content integrity; signing/attestation is deferred | S1 baseline decided |
| Update mechanism | Manually verify and extract each version into a new directory, then explicitly switch paths; reconsider package managers later | Support matrix, schema compatibility, backup/recovery gates | S1 model decided; final ZIP awaits rerun |
| Uninstall/data retention | Manually remove the version directory; retain run/store by default and separately authorize data deletion | T024 lifecycle record and explicit deletion authorization | S1 model and library rules frozen |
| VS Code/AI integration | Optional for core offline mode; connecting to Copilot Chat requires VS Code 1.82+, usable Copilot Chat, and SessionBridge 0.1.1 | Matching channel/instance, fixed non-legacy `silent` mode, no authority beyond the core API, and executable-adapter product-loop acceptance | Host prerequisites defined; product loop pending |

S1 defines no default install directory, publishes no one-click installer, does not modify PATH, claims no automatic updates, and requires no administrator privileges. The default state/store location, formal download entry point, and EOL notification channel remain open.

## 4. Completed Candidate Portable-Package Drill

On 2026-09-09, two trusted Windows amd64 CI artifacts completed offline verifier checks plus first run, upgrade, rollback, and evidence-preserving uninstall under an explicit temporary directory. The old and new full candidate SHAs are `1e7af676e8e84028c7ecfef1ccde91213738c20c` and `603e8dc03e9fb74ce9c53243d38e704fd8469c6a`. See the [S1 Windows portable-package installation drill](validation/s1-install-drill_EN.md) for evidence.

This closes the engineering question of whether candidate artifacts can execute the lifecycle and supports the manually extracted portable-ZIP, no-PATH model approved on 2026-09-09. The formal download entry point and final ZIP do not yet exist, so the final carrier must rerun the drill.

## 5. Planned User Path

1. Select an OS/architecture-matching artifact from a trusted release location.
2. Independently verify the fixed commit, GitHub CI result, SHA256SUMS, SBOM/license manifest, and support status; state explicitly that checksums do not authenticate the publisher.
3. Extract/install into a separate directory without overwriting a running host; run version and read-only self-checks.
4. Initialize explicit state/store and workspace configuration, proving no overlap with source or run workspaces.
5. Run side-effect-free static preview before separately authorizing capability probes and execution.
6. Before upgrade, stop writers and make a complete recoverable backup to a new store. Preserve the old binary and store on failure.
7. Uninstall retains evidence and configuration by default; data deletion is separately authorized.

See the [portable ZIP installation guide](INSTALLATION_EN.md) for copyable commands. T024 supplies library-level complete backup, new-store verification, retention-by-default, and explicit deletion-authorization rules; formal release filenames and download locations remain bound to release approval.

## 6. Pre-Release Installation Gates

- Native Windows 11 amd64 install, first run, upgrade, rollback, and uninstall pass. Linux claims match the current support matrix only.
- Artifacts support offline SHA256SUMS, SBOM, license, and exact-version verification; this evidence does not authenticate the publisher. If signing/attestation is added later, unknown or stale revocation information never displays fully verified.
- A clean user machine runs ProofRail without Go. External tools required by project tasks are listed and never automatically installed.
- Installation neither modifies target source/Git nor removes SessionBridge/user toolchains, and secrets never enter configuration, logs, or backups.
- State/store backup covers the complete reference closure and has been restored into a new directory.
- Every installation command in README and operations docs has been exercised against the corresponding artifact on a supported platform.

## 7. Documentation Completion

T016 has frozen CLI, configuration, and current source-run behavior; T024 has frozen library-level backup, pre-upgrade blocking, retirement, and data-retention rules. T017 candidate CI artifacts, fixed-commit, SHA256SUMS, SBOM/license gates, and the dual-platform trust matrix are complete. Signing/attestation is a later enhancement. The S1 installation model was approved on 2026-09-09. This document still does not claim a production release until the final ZIP rerun and S1 exit are approved; [Operations](OPERATIONS_EN.md) continues to distinguish executable behavior from future design. See the candidate [support matrix](S1_SUPPORT_MATRIX_EN.md) and [release notes](S1_RELEASE_NOTES_EN.md) for release boundaries.