# ProofRail Architecture Decision Register

[简体中文](ADR_REGISTER.md)

Date: 2026-09-06. Distinguish RFC-decided, proposed and experimental decisions. This is not S1 approval. Record approver, date, RFC changes and evidence for each accepted decision; models cannot sign for humans.

| ID | Decision/proposal | Alternatives/rationale | Status/completion |
|---|---|---|---|
| ADR-001 | Go modular monolith, CLI/TUI, single writer, no default Git writes | No script translation/distribution complexity | Decided by RFC 4/11/16; T001 confirms dependency baseline |
| ADR-002 | Freeze Schema 2020-12 toolchain, canonical JSON/paths/errors/config precedence | Use `canonicalize` 4.0.0 to reproduce RFC 8785 vectors; Ajv 8.20.0 remains independent of the future Go core | T002 froze 29 structural Schemas; all 82 current hand-authored fixtures and both canonical vectors pass |
| ADR-003 | Event/journal acceptance protocol and Windows replacement/stop primitives | Rename covers one file only; prove durability, rollback and cross-file reader semantics | T004 observed that an open reader blocks ordinary replacement and directory `Sync` is unavailable on NTFS, while ext4 supports open-reader replacement and directory `Sync`; rename-before-receipt remains uncertain. S1 must use platform primitives and add power-loss injection; directory separation cannot replace required OS enforcement |
| ADR-004 | SessionBridge v0.1.1 silent plus file queue, core durable idempotency | No fork, visible/auto or GUI fallback; cache is not exactly-once | T004 observed multi-winner rename-only claims on Windows, making fixed-claim `O_EXCL`, requestId deduplication and generation fencing mandatory. Forty no-model SessionBridge tests pass; live LM API remains a host capability gate |
| ADR-005 | Bootstrap seed via fixed commit, normal build, independent Go tests and human audit; N builds N+1 with isolated replay | No candidate self-certification or host overwrite | RFC 16.5 decided; T017 records seed commit/hash, external oracle and rollback path |
| ADR-006 | Windows 11 amd64 production, Linux amd64 core CI, pure-Go release, pinned tooling/dependencies, and S1 commit+CI+SHA256SUMS+SBOM/license evidence | RFC 16.2 mentions Win10 unlike 9.6/14; no S1 offline key system, with signing/attestation deferred | Owner simplified policy 2026-09-08; T017 supplies baseline release evidence |
| ADR-007 | Minimal context, offline first, task budgets, bilingual IDs | Avoid costly repeated whole-repository reviews; English is not a second authority | Owner approved 2026-09-06; zero new paid S0 calls, separate authorization for increases |
| ADR-008 | RFC section 19 six minimum product flows in S1; local modules, no commercial backend | Reject engine-only delivery, implicit source/deployment writes or silent updates; explicitly review added cost | Owner approved scope 2026-09-06; T002 froze records, T003 added goldens and T019–T024 run checks; no S1 coding authorization |
| ADR-009 | S1 Windows amd64 uses a manually extracted portable ZIP in a user-selected independent directory, invoked by full `prfrail.exe` path without modifying PATH, the registry, or system directories | Reject installers, package managers, automatic updates, and implicit version switching; minimize installation side effects and support directory-level coexistence/rollback | Product owner explicitly approved 2026-09-09; bind and rerun the final ZIP before release |
| ADR-010 | Product versions use `vMAJOR.MINOR.PATCH` (first candidate `v0.1.0`); during the prototype phase, a new minor announces the previous minor's EOL in both GitHub Release notes and the repository support matrix, initially 90 calendar days after release | Patch releases do not start a new EOL countdown; before each minor release, review stability, adoption, and maintenance capacity to decide whether future windows should grow. Never retroactively shorten an announced window; exceptions require a dated ADR with rationale and expiry, explicit product-owner approval, and synchronized publication in both channels | `PROPOSED / NOT APPROVED`; 90 days fits the early prototype, but the notice channels and evolution rule require explicit product-owner approval before the first release |

## 1. Release and Dependencies

Do not raise the module's Go 1.22 minimum just because local Go is 1.27.0. Pin a security-supported patch toolchain at implementation time and verify minimum compatibility separately. A higher TUI/dependency requirement needs an ADR/RFC change before go.mod. Prefer standard library and mature pure-Go TOML/schema/canonical/TUI libraries; T001/T002 verify exact versions, licenses, vulnerabilities and minimum Go before locking. No project dependencies are installed here.

The owner approved this support policy on 2026-09-06: security fixes for the current stable minor and previous minor, retaining the latter for 90 days; historical unknown runs stay read-only without automatic migration. ADR-010 proposes three-part `vMAJOR.MINOR.PATCH` versions, beginning with candidate `v0.1.0`; patch releases do not reset EOL, while each new minor announces the previous minor's EOL in both GitHub Release notes and the repository support matrix. Ninety days is the initial early-prototype window, with pre-release reviews allowed to lengthen only future windows as the product matures. Exceptions require an owner-approved ADR with a date, rationale, and expiry, and cannot retroactively shorten an announced window. The proposal is not approved, and this is not yet a published SLA.

The owner simplified release trust on 2026-09-08: S1 does not maintain an offline Ed25519 key system and instead uses a fixed commit SHA, GitHub CI, SHA256SUMS, SBOM and license manifest as baseline release evidence. SHA256SUMS proves content agreement with the manifest, not publisher identity. Signing or keyless attestation is deferred and does not block T017/T018 or S1 exit. The optional `signature-receipt` capability remains for cross-host or later signed scenarios.

## 2. Ambiguities to Close

T002 has frozen empty-chain, array merge, duration/size/budget units/defaults, ID/attempt/requestId, error, path/link, signing, review/waiver, ticket, runner, lease and Windows 10 rules. T003 independent fixtures/validator and T004 platform probes are complete; production defaults still require their S1 implementations, native failure tests and supply-chain gates.

## 3. Current S0 Readiness

Verdict: **NOT_READY**. This does not mean this documentation round failed or authorize skipping gates.

| RFC 17 item | Evidence/gap | Closing tasks |
|---|---|---|
| 1 P0 document review | Owner accepted the third-party read-only audit recommendations and approved S0 specification freeze/decision sheet on 2026-09-06; independent security sign-off remains item 3 | Closed (T001) |
| 2 Schemas/goldens/validator | All 29 structural Schemas have positive and negative cases; 82 current fixtures and two canonical vectors pass the independent oracle | Closed (T003) |
| 3 STRIDE | Risks/owners/AT mapping exists; T003 read-only security review found no fail-open path, while later gates still own runtime evidence | S0 static portion closed (T003) |
| 4 whois/non-C | whois remains a mapping example; Go/JavaScript/Python non-C fixtures pass | Closed (T003) |
| 5 Spikes | Windows 11/NTFS, Ubuntu 24.04/ext4, TUI, process trees, file queue and no-model SessionBridge contracts were exercised; limitations are in `validation/s0-spikes_EN.md` | Closed (T004) |
| 6 Traceable backlog | Owner approved current S0/S1 scope and traceability plan | Closed (T001) |
| 7 Platform/Go/license/support/response/release evidence | Platform/support/per-dependency review and simplified release-evidence policies approved; exact dependencies/channel remain open | T017 |
| 8 User approval | S0 freeze, S1 scope, zero-new-paid-call budget and whois-shadow principle approved; no S1 coding, concrete whois access or release authorization | Closed (T001); execution authorization remains separate |
| 9 step/TUI/bootstrap | Step/TUI contracts and a TUI prototype are verified; seed/oracle identity remains open | T017 |
| 10 Polyglot schemas | C+Go and C+JavaScript+Python ownership/dependency fixtures pass | Closed (T003) |
| 11 Handoff | Valid handoff and invalid read-only-target fixtures pass | Closed (T003) |
| 12 Documentation | Valid impact/waiver, missing-required-docs and generation-cycle fixtures pass | Closed (T003) |
| 13 Product/lifecycle | PC-01–PC-06 records have independent structural positive/negative fixtures; runtime acceptance remains an S1 exit gate | S0 contract portion closed (T003) |

Next allowed work: review this package, complete S0 specifications and independent spikes. S1 needs evidence for every gate and explicit owner approval; document existence cannot automatically authorize it.

## 4. T001 Owner Decision Sheet

Status: `APPROVED_FOR_S0_SPEC_FREEZE`. On 2026-09-06, the product owner (current repository user) explicitly approved the S0 specification freeze and this decision sheet based on a third-party read-only audit. T001 is complete and T002 may begin. This is not independent security sign-off or authorization for S1 production coding, paid calls, concrete whois access, Git operations or release.

| Decision | Approved content | Current constraint/follow-up | Approval record |
|---|---|---|---|
| S1 scope | R1–R22 plus PC-01–PC-06; no S2/S3 pull-forward | Scope freeze only; no S1 production coding authorization | Owner approved 2026-09-06 |
| S0 budget | Current authorized session/local tools only; zero new paid calls | No paid model/service calls; automatic dependency installation still needs task-specific grounds | Owner approved 2026-09-06 |
| Product owner | Current repository user is final scope/release approver | S1 coding, Git and release each still require explicit authorization | Owner confirmed 2026-09-06 |
| Security review | Designate a reviewer who did not produce the candidate; disclose insufficient independence if one person holds both roles | Reviewer not yet designated; gate 17.3 stays NOT_READY | Principle approved 2026-09-06; actor/conclusion pending |
| Platform/support | Windows 11 amd64 production; Linux amd64 core CI; current/previous minor with a 90-day previous-minor window | Freeze EOL notice and exceptions before first release | Owner approved 2026-09-06 |
| Release trust | S1 uses fixed commit SHA, GitHub CI, SHA256SUMS, SBOM and license manifest; no offline key system | Never describe SHA256SUMS as identity authentication; signing/attestation is a later enhancement | Owner approved simplified policy 2026-09-08 |
| Installation model | Manually extract the Windows amd64 portable ZIP into a user-selected independent directory; invoke by explicit path without changing PATH, registry, or system directories | The final ZIP still requires release-record binding and rerun; no publication authorization | Owner approved 2026-09-09 |
| whois shadow | In principle, later use fixed redacted read-only fixtures | Each run still requires input path/redaction/window authorization; no current asset access | Owner approved principle 2026-09-06 |
| Dependency licenses | Prefer standard library; review each TOML/Schema/canonical/TUI candidate for license, minimum Go and maintenance | T002 records name/version/source/license/hash/date before adoption | Owner approved policy 2026-09-06 |

This approval explicitly means “approve S0 specification freeze.” “Approve S1 coding” still requires a separate instruction; neither implies the other. This approval does not authorize Git commit/push/release, concrete whois asset access or new paid calls.