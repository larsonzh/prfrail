# ProofRail Architecture Decision Register

[简体中文](ADR_REGISTER.md)

Date: 2026-09-06. Distinguish RFC-decided, proposed and experimental decisions. This is not S1 approval. Record approver, date, RFC changes and evidence for each accepted decision; models cannot sign for humans.

| ID | Decision/proposal | Alternatives/rationale | Status/completion |
|---|---|---|---|
| ADR-001 | Go modular monolith, CLI/TUI, single writer, no default Git writes | No script translation/distribution complexity | Decided by RFC 4/11/16; T001 confirms dependency baseline |
| ADR-002 | Freeze Schema 2020-12 toolchain, canonical JSON/paths/errors/config precedence | Evaluate a mature RFC 8785 implementation; sorted keys alone are not canonical; validator independent of core | Proposed; T002 freezes types, vectors, BOM inputs, exit codes and compatibility, revising RFC first |
| ADR-003 | Event/journal acceptance protocol and Windows replacement/stop primitives | Rename covers one file only; prove durability, rollback and cross-file reader semantics | Experimental; T004 native Windows/Linux fault results; directory separation cannot replace required OS enforcement |
| ADR-004 | SessionBridge v0.1.1 silent plus file queue, core durable idempotency | No fork, visible/auto or GUI fallback; cache is not exactly-once | RFC boundary decided; framing/claim/errors/capabilities freeze in T002/T004 |
| ADR-005 | Bootstrap seed via normal build, independent Go tests and human audit; N builds N+1 with isolated replay | No candidate self-certification or host overwrite | RFC 16.5 decided; identify seed/hash/oracle/signer/rollback before T017 |
| ADR-006 | Windows 11 amd64 production, Linux amd64 core CI, pure-Go release, pinned tooling/signing | RFC 16.2 mentions Win10 unlike 9.6/14; propose section 14 scope, no Win10 promise before review | Main platform route decided; T001/T002 freeze support, private reporting, licenses and signing |
| ADR-007 | Minimal context, offline first, task budgets, bilingual IDs | Avoid costly repeated whole-repository reviews; English is not a second authority | Proposed; owner approves cost caps; follow doc plan |
| ADR-008 | RFC section 19 six minimum product flows in S1; local modules, no commercial backend | Reject engine-only delivery, implicit source/deployment writes or silent updates; explicitly review added cost | Proposed, unapproved; T001 budget/scope, T002/T003 records/goldens, T019–T024 runtime checks |

## 1. Release and Dependencies

Do not raise the module's Go 1.22 minimum just because local Go is 1.27.0. Pin a security-supported patch toolchain at implementation time and verify minimum compatibility separately. A higher TUI/dependency requirement needs an ADR/RFC change before go.mod. Prefer standard library and mature pure-Go TOML/schema/canonical/TUI libraries; T001/T002 verify exact versions, licenses, vulnerabilities and minimum Go before locking. No project dependencies are installed here.

Proposed support: security fixes for current stable minor and previous minor, retaining the latter for 90 days; historical unknown runs stay read-only without automatic migration. Owner approval is required before external commitment.

Proposed signing: maintainer-held offline Ed25519 signature over the checksum manifest, published public-key fingerprint and rotation/revocation procedure. A hosted signing service is an alternative with external dependencies. Choice, custodian and bootstrap identity remain unapproved. SHA256SUMS alone is not identity authentication.

## 2. Ambiguities to Close

T002 makes one focused RFC revision for empty chains, array merging, duration/size/budget units/defaults, ID/attempt/requestId types, error codes, path/symlink/hardlink policy, signing coverage, review/waiver authorization, ticket thresholds, runner capabilities, lease renewal/takeover and Windows 10 wording. Until frozen, write documents/isolated experiments, not production defaults.

## 3. Current S0 Readiness

Verdict: **NOT_READY**. This does not mean this documentation round failed or authorize skipping gates.

| RFC 17 item | Evidence/gap | Closing tasks |
|---|---|---|
| 1 P0 document review | Combined bilingual package supplied; no human sign-off | T001 |
| 2 Schemas/goldens/validator | No complete machine contracts yet | T002/T003 |
| 3 STRIDE | Risks/owners/AT mapping written, independent review pending | T001/T003 |
| 4 whois/non-C | Mapping planned; validated fixtures absent | T003 |
| 5 Spikes | Go available; TUI/replacement/process trees/transports/paths unverified | T004 |
| 6 Traceable backlog | Development plan exists, review pending | T001 |
| 7 Platform/Go/license/support/response/signing | ADR-006 unapproved, dependencies not pinned | T001/T002 |
| 8 User approval | Current authorization covers documents, not S1 risk/cost/shadow execution | T001 |
| 9 step/TUI/bootstrap | Designed; prototype and seed/oracle identity absent | T003/T004 |
| 10 Polyglot schemas | Ownership rules written; dual/triple-language goldens not executed | T003 |
| 11 Handoff | Contract/threats written; machine positive/negative fixtures not executed | T003 |
| 12 Documentation | Policy/impact/generation rules written; machine checks not executed | T003 |
| 13 Product/lifecycle | PC-01–PC-06/AT-16–AT-21 mapped; authority/revocation races, measurement/shared locks, export/backup consistency and deletion conflicts not frozen/verified | T001/T002/T003 |

Next allowed work: review this package, complete S0 specifications and independent spikes. S1 needs evidence for every gate and explicit owner approval; document existence cannot automatically authorize it.

## 4. T001 Owner Decision Sheet

Status: `READY_FOR_OWNER_REVIEW`, not an approval record. The 2026-09-06 automated review can only recommend. No human selections were received, so T001 remains incomplete and T002 does not begin freezing implementation fields.

| Decision | Recommended default | Conservative behavior while unapproved | Approval record requirement |
|---|---|---|---|
| S1 scope | R1–R22 plus PC-01–PC-06; no S2/S3 pull-forward | Maintain specification drafts only; no S1 production code | Owner/date/accepted or removed items/budget impact |
| S0 budget | Current authorized session/local tools only; zero new paid calls | No paid model/service calls or automatic dependency installation | Measurement mode, call/token/money cap and expiry |
| Product owner | Current user is final scope/release approver | Never self-sign gate 17.8 | Actor identity and approval date |
| Security review | A reviewer who did not produce the candidate; disclose insufficient independence if one person holds both roles | Gate 17.3 remains NOT_READY | Reviewer, scope, conclusion and open risks |
| Platform/support | Windows 11 amd64 production; Linux amd64 core CI; current/previous minor with 90-day previous-minor window | No public support-duration promise | Start/end rule, EOL notice and exceptions |
| Release signing | Maintainer offline Ed25519 over checksum manifest; no keys in repository/model | No production release without a signer | Signer, public-key fingerprint and rotation/revocation process |
| whois shadow | In principle, later use fixed redacted read-only fixtures; each execution still needs path authorization | Do not read/copy/change whois assets | Input scope/hash, redactor, execution window and stop conditions |
| Dependency licenses | Prefer standard library; review each TOML/Schema/canonical/TUI candidate for license, minimum Go and maintenance | No go.mod edits or guessed version locks | Name/exact version/source/license/hash/review date |

Approval must explicitly distinguish “approve S0 specification freeze” from “approve S1 coding”; neither implies the other. A generic “continue” permits review-material improvements and side-effect-free checks, not role, budget, signing or whois-data authorization. No approval here implies Git commit/push/release authorization.