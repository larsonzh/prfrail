# ProofRail Security and Permission Design

[简体中文](SECURITY.md)

Date: 2026-09-07; security-design baseline without independent security sign-off. This document is the normative authority for permissions, threat modeling, and hard security gates. [Project proposal](RFC-proofrail-unattended-ai-engineering-product.md) sections 12 and 16.7-16.8 retain risk and design provenance. Requirements: REQ-004/005/007/013/019/021-026/028.

## 1. Assets and Trust

Protect source trees, baseline/accepted snapshots, events/journals/receipts, secrets, approval authority and cost budgets. S1 trusts the local OS and core binary, not models, repository text, hooks, adapters or external results. Same-user directory separation is not a malicious-process sandbox; hashes are not authentication. Reject missing capabilities rather than silently replacing required isolation with risk acknowledgment.

## 2. Write Permissions

| Actor | Source | run-workspace | outbox | Core state/store | Acceptance/publication |
|---|---|---|---|---|---|
| chain/core | Read-only capture | Materialize/restore after stopping writers | Read/validate | Locked writes | Only after independent authorization |
| Managed agent | No writes | Through applier | Candidate/response only | Denied | Denied |
| Isolated agent/hook | Denied | Authorized targets/cwd only | Own artifacts only | Denied | Denied |
| Handoff operator | Outside handoff authorization | Valid exclusive lease scope | Structured return | Denied | No automatic approval |
| Console/reviewer | Read-only references | Read-only diff | No direct writes | Control API only | Authorized commands, no receipt editing |

Runners/OS capabilities enforce restrictions, not prompts alone. Unknown processes or uncertain writer termination block restore, handoff and publication. An executor cannot combine policy mutation with independent approval of the same change.

## 3. STRIDE and Required Tests

| ID/category | Threat/fault | Control and owner | Tests |
|---|---|---|---|
| SEC-01 Spoofing | Forged review, old requestId/lease | Actor/hash/attempt/expiry binding; chain/tickets/adapters | AT-05/08/09 |
| SEC-02 Tampering | Replaced hooks/manifests/receipts, traversal/symlink/reparse escape | Pinned hashes, real-path/open-time checks, offline chain; taskdef/snapshot/evidence | AT-02/04/06 |
| SEC-03 Repudiation | Deleted/reordered events, anonymous waiver | Sequence/previous hash, actor/reason/expiry; evidence/chain | AT-04/05 |
| SEC-04 Disclosure | Tokens in logs, prompts or snapshots | Environment allowlist, secret references, pre-archive/export scan and size bounds; gates/evidence/adapters | AT-06/10 |
| SEC-05 Denial of service | Full disk, unlimited output, stuck process trees, retry storms | Quotas, runtime monitoring, process-tree stop, shared budgets; guard/gates/tickets | AT-03/06/07 |
| SEC-06 Elevation | Prompt-based self-approval, shell injection, generated self-validator | Data is not authority, argv, separated roles, templates; taskdef/gates/repair | AT-05/06/11 |
| SEC-07 Compound | Human/AI concurrent writes, expired handoff accepted | Single lease, stop proof, full return checks; chain/guard | AT-09 |
| SEC-08 Supply chain | Candidate overwrites seed/verifier, self-certifies release, tampers through a same-runner background process, drifts CI Actions, or spoofs checksum provenance | Separate runners for candidate probing and trusted verification, an independently reviewed earlier verifier commit, fixed candidate/bootstrap commits, full Action SHAs, GitHub CI, external oracle, SHA256SUMS/SBOM/license evidence, and human approval; release/harness | AT-13/14 |
| SEC-09 Elevation/disclosure | Preview executes malicious hooks; export aliases overwrite source or leak secrets | Zero-execution static preview, real destination checks/integrity/export scanning; taskdef/snapshot | AT-16/17 |
| SEC-10 Spoofing/effects | Dispatch after revocation, expired approval, uploads claimed reversible | Reauthorize each effect, in-flight stop evidence, default-deny external writes; chain/guard/gates | AT-18/19 |
| SEC-11 Denial/repudiation | Assumed timeout refunds, cross-run budget bypass, duplicate settlements | Durable reservations, unknown holds, shared locks/deduplication; tickets/adapters | AT-20 |
| SEC-12 Loss/disclosure | Incomplete backups, uninstall deletes shared tools, secrets in old evidence | Closure checks, new-store drills, retention by default/explicit disposition; snapshot/evidence/release | AT-21 |
| SEC-13 Spoofing/elevation/disclosure | Forged or stale operator replies, UI text treated as authority, or secrets persisted from dialogue | Actor/attempt/context-hash/authority binding, structured allowed responses, secret-direct, and replay deduplication; chain/console/evidence | AT-22 |
| SEC-14 Tampering/elevation/repudiation | Substituted CLI Agent binary/config, source/store escape, tool/network escalation, forged/missing events, ambiguous stop/resume, or self commit/push/publish | Executable/config digests, capability probe, isolated workspace, process-tree supervision, events/logs and before/after manifests, independent gates/review; adapters/guard/evidence | AT-23 |
| SEC-15 Disclosure/repudiation/effects | A visible black-box Agent reads out-of-scope resources, leaks secrets, hides network/cost/background work, performs external Git/publication effects, or treats risk notice as a waiver | Explicit versioned risk acknowledgment, isolated workspace, user supervision/host controls, full post-return scan, unknown labels, independent gates/review; console/chain/guard/evidence | AT-24 |

Test IDs resolve in [TEST_STRATEGY_EN.md](TEST_STRATEGY_EN.md). Every mandatory negative case asserts no publication, no advancement and intact evidence, not merely a nonzero exit.

## 4. Secrets, Network and Retention

Passwords/tokens/MFA/passphrases never pass through models or config/logs. secret-direct goes only to a controlled terminal/provider; wait if unavailable. Record command metadata, hashes and redacted conclusions, not screens, keystrokes or complete REPL history by default. Scanning cannot recognize every secret; combine excluded paths, minimal allowlisted input and human review.

Network is default-deny; exceptions identify host/port/purpose, provider data policy and expiry. Reject runners unable to enforce required policy. The released ProofRail CLI and core runtime do not directly access GitHub, the GitHub API, raw content, package registries, or other public release sites; release verification consumes local artifacts, manifests, SBOMs, and license files only. GitHub Actions checkout, dependency downloads, and artifact uploads belong to the CI control plane and do not expand product runtime network authority. Pin dependency versions, sources, licenses and hashes. Never auto-install unknown plugins at runtime.

Retention follows the RFC; cleanup only unreferenced expired objects, never audit-locked data. Scan failure blocks export. On leakage, pause, isolate affected artifacts and notify the operator to revoke credentials; do not repeat secrets in diagnostics.

## 5. Review and Response

Independent review is required for broader write/network/tool authority, approval changes, generated hooks and canonical/signing/recovery protocols. Black-box risk acknowledgment is not a liability waiver: ProofRail still owns its promised isolation, scanning, gates, evidence and review, while the user owns only the stated host-supervision and unobservable external-behavior risks. Unfixed fail-close bypass, acceptance of failed candidates, secret leakage or unrecoverable overwrite blocks release.

Proposed response targets pending ADR-006: acknowledge private reports within two business days and classify impact within five. These are not published SLAs. Establish a private channel and maintainers before release; do not put credential-bearing reports or exploit details in public issues meanwhile. If signing/attestation is added later, its private keys or short-lived credentials must never enter repositories, model context or ordinary logs.

## 6. Product Lifecycle Boundaries

RFC section 19 PC-01–PC-06 do not authorize remote approval services or deployment engines. A trusted local control path creates manifest-bound authorization; models/repository text cannot authorize themselves. Same-user malicious-process defenses still require actual permission separation. Waivers cannot bypass integrity, stopped writers, exclusivity or secret protection.

Revocation blocks new actions and requests stopping; it cannot retract dispatched network calls. Unknown costs are not free. Exports/diagnostics are redacted but scanning cannot find all secrets. If historical evidence contains secrets, quarantine/restrict access, revoke credentials and block export. Human-approved deletion/retention conflict decisions produce secret-free disposition records, not rewritten history presented as intact.

Backup/restore must prove closure and consistency; uninstall preserves evidence and does not delete SessionBridge/Go. SSDs, snapshots and external backups may retain copies: logical deletion is not secure erasure. S1 separately displays content integrity and support status and states that SHA256SUMS does not prove publisher identity. If signing is added later, signature authenticity and revocation freshness remain separate signals.