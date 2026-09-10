# ProofRail S1 Candidate Release Notes

[中文](S1_RELEASE_NOTES.md)

Date: 2026-09-09. Status: `DRAFT / NOT RELEASED`. Version, download links, formal artifact filenames, support start date, and approver remain pending S1 exit. This file does not authorize a tag, GitHub Release, publication, or installation.

## Candidate Capabilities

- Local single-user, single-writer CLI baseline: `init`, `validate`, `config explain`, `preview`, `approvals`, `run`, and `report`.
- Unauthorized capabilities are denied by default; static preview does not invoke commands, networking, models, or credential probes.
- Noop-only chains can complete; executable `code/build/verify` steps fail closed.
- Evidence, snapshots, process/lease guards, managed change sets, transactional apply, authorization revocation, backup/restore, and uninstall rules that retain data by default.
- Minimal generic/C/Go examples and a redacted frozen-fixture read-only whois shadow.
- Candidate release chain with fixed verifier/candidate commits, dual-platform candidate probe/self-host, checksums, SPDX SBOM, and license manifest.

## Support and Known Limitations

See the [S1 candidate support matrix](S1_SUPPORT_MATRIX_EN.md). The S1 installation model is a [manually extracted Windows portable ZIP](INSTALLATION_EN.md), invoked by explicit path without modifying PATH. No formal release package exists and the current interface is line-oriented CLI. Neither formal AgentRunner/CLI Agent nor supervised visible black-box candidates are delivered. Even when implemented, the latter warrants post-return artifact acceptance only, not tools, network, cost or external effects, and requires explicit user risk acknowledgment and supervision. S1 also excludes signing/attestation, production Linux support, automatic updates, a Web console, the full TUI, and arbitrary executable hooks.

## Remaining Pre-Release Gates

1. Rerun verification, first run, upgrade, rollback, and uninstall against the final ZIP.
2. Bind version, artifact filenames, download location, EOL notification channel, and support date to these notes and the support matrix.
3. Obtain the independent-review verdict and explicit product-owner S1 exit/release approval.
4. Complete T025–T028 and pass AT-22 operator interaction, AT-23 AgentRunner, and AT-24 black-box candidate acceptance.