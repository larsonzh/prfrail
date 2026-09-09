# ProofRail S1 Candidate Support Matrix

[中文](S1_SUPPORT_MATRIX.md)

Date: 2026-09-09. Status: `CANDIDATE`, effective only after the formal release carrier and S1 exit are approved. This is not a published SLA.

| Scope | S1 candidate status | Constraints |
|---|---|---|
| Windows 11 amd64 | Production-support candidate | Pure-Go `CGO_ENABLED=0` portable ZIP; manual extraction, explicit-path invocation, no PATH modification; rerun the lifecycle against the final ZIP |
| Linux amd64 | Core CI | Build, test, candidate probe, and self-host only; no S1 production-support claim |
| macOS / Windows on Arm / Linux arm64 | Unsupported | Not included in the S1 acceptance matrix |
| generic/C/Go harnesses | Configuration and noop orchestration | Language scopes, read-only preview, and noop-only loops are verified; executable steps still fail closed |
| VS Code / SessionBridge | Optional adapter | Not a core installation dependency and grants no authority beyond ProofRail |

## Candidate Version and Security-Maintenance Policy

- The current stable minor and previous minor receive security fixes; the candidate maintenance window for the previous minor is 90 days.
- Unknown historical runs remain read-only and are not automatically migrated.
- The EOL notification channel and exceptions must still be frozen before the first formal release; this page does not replace the final notice.

## Release Trust Boundary

A candidate release requires a full commit SHA, GitHub CI, `SHA256SUMS`, an SPDX SBOM, and a license manifest. Checksums prove agreement with a manifest, not publisher identity; S1 does not claim signing or attestation. See the [portable ZIP installation guide](INSTALLATION_EN.md); filenames, download locations, and the support start date must be approved with the final release notes.