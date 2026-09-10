# ProofRail S1 Candidate Support Matrix

[中文](S1_SUPPORT_MATRIX.md)

Date: 2026-09-09. Status: `CANDIDATE`, effective only after the formal release carrier and S1 exit are approved. This is not a published SLA.

| Scope | S1 candidate status | Constraints |
|---|---|---|
| Windows 11 amd64 | Production-support candidate | Pure-Go `CGO_ENABLED=0` portable ZIP; manual extraction, explicit-path invocation, no PATH modification; rerun the lifecycle against the final ZIP |
| Linux amd64 | Core CI | Build, test, candidate probe, and self-host only; no S1 production-support claim |
| macOS / Windows on Arm / Linux arm64 | Unsupported | Not included in the S1 acceptance matrix |
| generic/C/Go harnesses | Configuration and noop orchestration | Language scopes, read-only preview, and noop-only loops are verified; executable steps still fail closed |
| CLI Agent / AgentRunner | Planned, not delivered | Target formal AI execution channel; requires pinned version/config, capability probe, isolated workspace, process/session evidence and independent acceptance; T026/T027/AT-23 block S1 |
| VS Code / Copilot Chat / SessionBridge visible black-box candidate | Planned, reduced assurance, not delivered | Not needed for core offline CLI; bridges an isolated candidate only with explicit risk acknowledgment and human supervision, followed by independent acceptance; does not warrant tools/network/cost/external effects; T028/AT-24 block S1 |

## Candidate Version and Security-Maintenance Policy

- Product versions use `vMAJOR.MINOR.PATCH`, beginning with candidate `v0.1.0`; patch releases do not start a new EOL countdown.
- The current stable minor and previous minor receive security fixes; the candidate maintenance window for the previous minor is 90 days.
- Unknown historical runs remain read-only and are not automatically migrated.
- Proposed ADR-010: when a new minor is released, announce the previous minor's EOL date in both the GitHub Release notes and this support matrix, with EOL on the 90th calendar day after release.
- Ninety days is the initial prototype-phase window; before each minor release, review stability, adoption, and maintenance capacity, and only lengthen future newly announced windows.
- Deviations require a dated ADR with rationale and expiry, explicit product-owner approval, and synchronized publication in both channels, and must not retroactively shorten an announced window. This proposal is not approved, and this page does not replace the final notice.

## Release Trust Boundary

A candidate release requires a full commit SHA, GitHub CI, `SHA256SUMS`, an SPDX SBOM, and a license manifest. Checksums prove agreement with a manifest, not publisher identity; S1 does not claim signing or attestation. See the [portable ZIP installation guide](INSTALLATION_EN.md); filenames, download locations, and the support start date must be approved with the final release notes.