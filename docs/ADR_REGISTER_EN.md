# ProofRail Architecture Decision Register

[简体中文](ADR_REGISTER.md)

Date: 2026-09-06. Distinguish RFC-decided, proposed and experimental decisions. This is not S1 approval. Record approver, date, RFC changes and evidence for each accepted decision; models cannot sign for humans.

| ID | Decision/proposal | Alternatives/rationale | Status/completion |
|---|---|---|---|
| ADR-001 | Go modular monolith, CLI/TUI, single writer, no default Git writes | No script translation/distribution complexity | Decided by RFC 4/11/16; T001 confirms dependency baseline |
| ADR-002 | Freeze Schema 2020-12 toolchain, canonical JSON/paths/errors/config precedence | Evaluate a mature RFC 8785 implementation; sorted keys alone are not canonical; validator independent of core | Partially frozen: chain/hook/target/workspace schemas, JCS, UTF-8 without BOM, version/path foundations are set; validator, vectors, exit codes and remaining objects await T002/T003 |
| ADR-003 | Event/journal acceptance protocol and Windows replacement/stop primitives | Rename covers one file only; prove durability, rollback and cross-file reader semantics | Experimental; T004 native Windows/Linux fault results; directory separation cannot replace required OS enforcement |
| ADR-004 | SessionBridge v0.1.1 silent plus file queue, core durable idempotency | No fork, visible/auto or GUI fallback; cache is not exactly-once | RFC boundary decided; framing/claim/errors/capabilities freeze in T002/T004 |
| ADR-005 | Bootstrap seed via normal build, independent Go tests and human audit; N builds N+1 with isolated replay | No candidate self-certification or host overwrite | RFC 16.5 decided; identify seed/hash/oracle/signer/rollback before T017 |
| ADR-006 | Windows 11 amd64 production, Linux amd64 core CI, pure-Go release, pinned tooling/signing | RFC 16.2 mentions Win10 unlike 9.6/14; section 14 is the first-release scope, with no Win10 commitment | Owner approved policy 2026-09-06; T002 freezes dependencies, T017 signer/key/revocation details |
| ADR-007 | Minimal context, offline first, task budgets, bilingual IDs | Avoid costly repeated whole-repository reviews; English is not a second authority | Owner approved 2026-09-06; zero new paid S0 calls, separate authorization for increases |
| ADR-008 | RFC section 19 six minimum product flows in S1; local modules, no commercial backend | Reject engine-only delivery, implicit source/deployment writes or silent updates; explicitly review added cost | Owner approved scope 2026-09-06; T002/T003 records/goldens, T019–T024 runtime checks; no S1 coding authorization |

## 1. Release and Dependencies

Do not raise the module's Go 1.22 minimum just because local Go is 1.27.0. Pin a security-supported patch toolchain at implementation time and verify minimum compatibility separately. A higher TUI/dependency requirement needs an ADR/RFC change before go.mod. Prefer standard library and mature pure-Go TOML/schema/canonical/TUI libraries; T001/T002 verify exact versions, licenses, vulnerabilities and minimum Go before locking. No project dependencies are installed here.

The owner approved this support policy on 2026-09-06: security fixes for the current stable minor and previous minor, retaining the latter for 90 days; historical unknown runs stay read-only without automatic migration. Freeze EOL notification and exceptions before first release; this is not yet a published SLA.

The owner approved this signing policy on 2026-09-06: maintainer-held offline Ed25519 signature over the checksum manifest, published public-key fingerprint and rotation/revocation procedure. Record the signer, custodian, public-key fingerprint and bootstrap identity before T017. SHA256SUMS alone is not identity authentication.

## 2. Ambiguities to Close

T002 makes one focused RFC revision for empty chains, array merging, duration/size/budget units/defaults, ID/attempt/requestId types, error codes, path/symlink/hardlink policy, signing coverage, review/waiver authorization, ticket thresholds, runner capabilities, lease renewal/takeover and Windows 10 wording. Until frozen, write documents/isolated experiments, not production defaults.

## 3. Current S0 Readiness

Verdict: **NOT_READY**. This does not mean this documentation round failed or authorize skipping gates.

| RFC 17 item | Evidence/gap | Closing tasks |
|---|---|---|
| 1 P0 document review | Owner accepted the third-party read-only audit recommendations and approved S0 specification freeze/decision sheet on 2026-09-06; independent security sign-off remains item 3 | Closed (T001) |
| 2 Schemas/goldens/validator | No complete machine contracts yet | T002/T003 |
| 3 STRIDE | Risks/owners/AT mapping written, independent review pending | T001/T003 |
| 4 whois/non-C | Mapping planned; validated fixtures absent | T003 |
| 5 Spikes | Go available; TUI/replacement/process trees/transports/paths unverified | T004 |
| 6 Traceable backlog | Owner approved current S0/S1 scope and traceability plan | Closed (T001) |
| 7 Platform/Go/license/support/response/signing | Platform/support/offline Ed25519/per-dependency review policies approved; exact dependencies/signer/key/channel remain open | T002/T017 |
| 8 User approval | S0 freeze, S1 scope, zero-new-paid-call budget and whois-shadow principle approved; no S1 coding, concrete whois access or release authorization | Closed (T001); execution authorization remains separate |
| 9 step/TUI/bootstrap | Designed; prototype and seed/oracle identity absent | T003/T004 |
| 10 Polyglot schemas | Ownership rules written; dual/triple-language goldens not executed | T003 |
| 11 Handoff | Contract/threats written; machine positive/negative fixtures not executed | T003 |
| 12 Documentation | Policy/impact/generation rules written; machine checks not executed | T003 |
| 13 Product/lifecycle | PC-01–PC-06/AT-16–AT-21 mapped; authority/revocation races, measurement/shared locks, export/backup consistency and deletion conflicts not frozen/verified | T001/T002/T003 |

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
| Release signing | Maintainer offline Ed25519 over checksum manifest; no keys in repository/model | Record signer, public-key fingerprint and rotation/revocation before T017 | Owner approved policy 2026-09-06 |
| whois shadow | In principle, later use fixed redacted read-only fixtures | Each run still requires input path/redaction/window authorization; no current asset access | Owner approved principle 2026-09-06 |
| Dependency licenses | Prefer standard library; review each TOML/Schema/canonical/TUI candidate for license, minimum Go and maintenance | T002 records name/version/source/license/hash/date before adoption | Owner approved policy 2026-09-06 |

This approval explicitly means “approve S0 specification freeze.” “Approve S1 coding” still requires a separate instruction; neither implies the other. This approval does not authorize Git commit/push/release, concrete whois asset access or new paid calls.