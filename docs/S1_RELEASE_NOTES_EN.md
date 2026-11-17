# ProofRail S1 Candidate Release Notes

[中文](S1_RELEASE_NOTES.md)

Date: 2026-09-09; updated 2026-11-17 (P78 `RELEASE-ZIP-BINDING`). Status: `CANDIDATE / RELEASED`. The v0.1.0 candidate release is created and bound: version `v0.1.0`, candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`, CI run `44043705691`, filename `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`; download entry = GitHub Release `v0.1.0`. The support start date and whole S1 exit approval remain pending before S2; this file documents the candidate release and does not constitute whole S1 exit approval.

## Candidate Capabilities

- Local single-user, single-writer CLI baseline: `init`, `validate`, `config explain`, `preview`, `approvals`, `run`, and `report`.
- Unauthorized capabilities are denied by default; static preview does not invoke commands, networking, models, or credential probes.
- Noop-only chains can complete; executable `code/build/verify` steps fail closed.
- Evidence, snapshots, process/lease guards, managed change sets, transactional apply, authorization revocation, backup/restore, and uninstall rules that retain data by default.
- Minimal generic/C/Go examples and a redacted frozen-fixture read-only whois shadow.
- Candidate release chain with fixed verifier/candidate commits, dual-platform candidate probe/self-host, checksums, SPDX SBOM, and license manifest.

## Support and Known Limitations

See the [S1 candidate support matrix](S1_SUPPORT_MATRIX_EN.md). The S1 installation model is a [manually extracted Windows portable ZIP](INSTALLATION_EN.md), invoked by explicit path without modifying PATH. The current artifact is the `v0.1.0` candidate release (not a published SLA) and the interface is a line-oriented CLI. Formal AgentRunner/CLI Agent are still undelivered; the supervised visible black-box candidate's product side is implemented by T028 (AT-24 conditionally closed) but its real CLI end-to-end depends on T027-B and it is not yet a runnable release. Even when implemented, the latter warrants post-return artifact acceptance only, not tools, network, cost or external effects, and requires explicit user risk acknowledgment and supervision. S1 also excludes signing/attestation, production Linux support, automatic updates, a Web console, the full TUI, and arbitrary executable hooks.

## Remaining Pre-Release Gates

1. Rerun verification, first run, upgrade, rollback, and uninstall against the final ZIP. — Done on 2026-11-17 by P78 `RELEASE-ZIP-BINDING` (evidence: `validation/release-zip-binding_EN.md`).
2. Bind version, artifact filenames, download location, EOL notification channel, and support date to these notes and the support matrix. — Version `v0.1.0`, filename `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`, and the download location (GitHub Release `v0.1.0`) are bound; the EOL notification channel was frozen by ADR-010 (2026-11-16); the support start date awaits S1 exit approval.
3. Obtain the independent-review verdict and explicit product-owner S1 exit/release approval. — The release part was approved by the product owner (2026-11-17, the P78 kickoff message in the same turn; ADR-046 Q10=(C) step one); whole S1 exit approval remains pending S2.
4. Complete T027 and pass AT-23; T028 is complete with AT-24 conditionally closed; S1 remains blocked by T027/AT-23. — T027-B is complete (partly not applicable, ADR-045); the AT-23 stage-1 condition closure is a permanent fact (ADR-045); S-G1 / S-G8 stay blocked; whole S1 Exit approval remains pending S2.

> Slice sync (P55 `EXIT-SEED-BUILD-EVIDENCE`): not applicable (no release artifact or feature change, evidence supplement only); reason: this slice produces no tag, artifact or feature item and only archives the AT-13 supplement under `docs/validation/evidence/EXIT-SEED-BUILD-EVIDENCE/`, so the candidate capabilities and remaining gates stay unchanged.

> Slice sync (P57 `P5-P6-CLOSEOUT-ASSESSMENT`): not applicable (no tag, artifact or feature item); reason: P57 only consolidates the S1 Exit remaining-gate closeout inventory (see `docs/validation/p57-closeout-assessment_EN.md` §4) and leaves this file's "gates remaining before release" wording untouched.

> Slice sync (P78 `RELEASE-ZIP-BINDING`): GitHub Release `v0.1.0` was created and its release record bound (version `v0.1.0` / candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86` / CI run `44043705691` / filename `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`); all four reverification items passed; the gate list above is updated accordingly.
