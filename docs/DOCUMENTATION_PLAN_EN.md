# ProofRail Documentation Index and Maintenance Plan

[简体中文](DOCUMENTATION_PLAN.md)

Date: 2026-09-07; continuously maintained documentation baseline. This is the first entry point for low-cost models, not a requirement to reread the full RFC every turn.

## 1. Authority and Reading Order

The former `RFC-proofrail-unattended-ai-engineering-product.md` is the [project proposal and historical design source](RFC-proofrail-unattended-ai-engineering-product.md), not the sole authority for every domain. Normative authority is divided below. Before implementation, revise and review the applicable Chinese authority, then synchronize English, schemas/goldens, tests, implementation, and user documentation. Bilingual or cross-domain conflicts block work until corrected; implementations must not choose an interpretation.

| Domain | Normative authority | Role of the proposal |
|---|---|---|
| Product scope, actors, requirements, and acceptance outcomes | [Product requirements](PRODUCT_REQUIREMENTS_EN.md) | Project rationale, problem source, and stage vision |
| End-to-end business rules and user workflows | [Business workflows](BUSINESS_WORKFLOWS_EN.md) | whois experience source and early flow draft |
| Module boundaries, dependency direction, and deployment topology | [Architecture](ARCHITECTURE_EN.md) | Rationale and candidate directions |
| Wire formats, schemas, state machines, events, receipts, and cross-record invariants | [Architecture contracts](CONTRACTS_EN.md); `schemas/` governs machine structure | Historical integrated design and migration source |
| Permissions, threats, and hard security gates | [Security design](SECURITY_EN.md) | Risk background and principle source |
| ATs, platforms, and fault-injection acceptance | [Test strategy](TEST_STRATEGY_EN.md) | Stage-goal source |
| Currently executable behavior | Verified implementation and [Operations](OPERATIONS_EN.md) | Not evidence of availability |
| Implementation order, dependencies, and completion | [Development plan](DEV_PLAN_EN.md) | Original roadmap source |
| Unresolved cross-domain decisions | [ADR register](ADR_REGISTER_EN.md) | Candidate solution source |

Read this page, the current T task in [DEV_PLAN_EN.md](DEV_PLAN_EN.md), the applicable authority sections/AT, and nearby code/tests. Read the proposal only to trace rationale or handle content not yet migrated. Do not load both complete language sets each round.

## 2. Delivery Index and Responsibility Coverage

| Package | Chinese/English | Coverage | Trigger/owner |
|---|---|---|---|
| Product | [需求](PRODUCT_REQUIREMENTS.md) / [Requirements](PRODUCT_REQUIREMENTS_EN.md) | Research, roles, use cases, scope, 28 functional/NFR items | Scope; product owner |
| Workflows | [流程](BUSINESS_WORKFLOWS.md) / [Workflows](BUSINESS_WORKFLOWS_EN.md) | Independent product main flow, failure/recovery, handoff, delivery, lifecycle | Business/state-flow changes; product owner/architect |
| Architecture | [架构](ARCHITECTURE.md) / [Architecture](ARCHITECTURE_EN.md) | Layers, modules, domain references, flow/transactions/storage | Boundaries/dependencies; architect |
| Protocols | [契约](CONTRACTS.md) / [Contracts](CONTRACTS_EN.md) | Schemas, snapshots/evidence, hooks, tickets/repair, adapters, handoff/docs | Wire/behavior; contract maintainer |
| Security | [安全](SECURITY.md) / [Security](SECURITY_EN.md) | STRIDE, permissions, secrets, audit, supply chain | Permissions/risks; independent reviewer |
| Usage/migration | [操作](OPERATIONS.md) / [Operations](OPERATIONS_EN.md) | Install/start/config, operations/debugging, whois mapping, self-host/release | CLI/deployment; implementer/release owner |
| Installation and release | [安装指南](INSTALLATION.md) / [Installation guide](INSTALLATION_EN.md); [安装规划](INSTALLATION_PLAN.md) / [Installation plan](INSTALLATION_PLAN_EN.md); [支持矩阵](S1_SUPPORT_MATRIX.md) / [Support](S1_SUPPORT_MATRIX_EN.md); [发行说明](S1_RELEASE_NOTES.md) / [Release notes](S1_RELEASE_NOTES_EN.md) | Executable portable-ZIP flow, pending release binding, upgrade/uninstall gates, candidate support and limitations | Artifact/path/signing/support decisions; release owner |
| Testing | [验证](TEST_STRATEGY.md) / [Testing](TEST_STRATEGY_EN.md) | 24 AT groups, tiers, platforms, faults, evidence | Requirements/contracts; test maintainer |
| Engineering | [规范](CODING_CONVENTIONS.md) / [Conventions](CODING_CONVENTIONS_EN.md) | Encoding, Go, Git, low-cost workflow | Tools/conventions; maintainer |
| Implementation | [计划](DEV_PLAN.md) / [Plan](DEV_PLAN_EN.md) | 28 tasks, dependencies, scope, traceability | Task/stage acceptance; implementation lead |
| Decisions | [ADR/readiness](ADR_REGISTER.md) / [Decisions](ADR_REGISTER_EN.md) | Trade-offs, bootstrap, release/support, 13 readiness gaps | Decisions/evidence; designated approver |
| Governance | [文档计划](DOCUMENTATION_PLAN.md) / this document | Bilingual rules, scoped reading, gates and future docs | Documentation changes; doc maintainer |

Continue referencing proposal sections 10/18 until their terminology is migrated; do not invent synonyms. Each implementation task moves its direct normative dependencies into the applicable authority. Combined operations/testing/migration packages reduce reading cost. Actual schema and machine-example completion comes from the development plan and validation reports, never document titles. Audit policy lives in security, release policy in operations/ADRs; external policies need approval before release.

## 3. Phased Documentation

| Phase | Deliverable | Gate |
|---|---|---|
| Documentation baseline | 12 Chinese/English pairs, README links, requirement/task checkpoints | Links, encoding, IDs, mappings; honest status |
| S0 T001-T004 | Sign-off, RFC field revision, machine schemas/positive-negative fixtures, spikes, dependency/signing ADRs | All RFC 17 gates and independent validation |
| Every S1 task | Synchronized code/tests/schema/affected bilingual sections; replace planned commands only after execution | Valid links/commands, generated freshness, relevant AT |
| S1 exit | Real installation tutorial, runnable generic/C/Go examples, recovery evidence, self-host/shadow/support matrix, release notes | Independent review, real platforms, owner approval |
| S2/S3 | New platform/composed-language/supervised/Web/script B docs and migrations; complete English RFC as needed | Contracts/tests first, not marketing as support evidence |

## 4. Low-Cost Model Task Card

Fill real values in the session before each task; do not create placeholder-filled task files:

```text
Task: Txxx; REQ: list; AT: list; Authorized scope: exact paths
Input: domain-authority section + required proposal source + current test; exclusions: task row
Hypothesis: one falsifiable cause; smallest check: one real command
Budget: time/calls/cost cap; no paid calls without authorization
Evidence: command/cwd/exit/case counts/artifact hash/unverified items
```

Change one testable slice at a time. Inspect the first failure; avoid recursive whole-repository searches, repeated long logs and invented frameworks. Carry task contract/parent hash/confirmed decisions/latest failure, not entire chat history. Label estimated token usage and never drop safety constraints to save tokens.

Escalate contract ambiguity, broader permissions, unauthorized paths, two same-error attempts without new evidence, budget exhaustion or unrecoverable state. Handover contains changes, verification and next action, not a repeated project description.

## 5. Review Rules

Chinese keeps original names, English uses _EN, with reciprocal links. REQ/AT/T/ADR/SEC and wire identifiers match exactly. Counts describe planned requirements/tasks, not completion percentages. Review both translations together, not a summary-only English edition. Resolve conflicts in the applicable Chinese authority identified in section 1, then synchronize English, schemas, and consumers.

Distinguish fact, proposal, measurement and plan; unexecuted examples are planned. Link only existing files; describe future paths as planned code text. Markdown uses UTF-8 BOM/LF; machine canonical/wire data does not inherit that encoding. Frozen snippets require independent validation/execution; prose is not machine validation.

This round runs terminal checks. T003 consolidates document/contract checks into maintained tooling; T017 wires CI. Keep stage reports/open issues and increment checkboxes only after acceptance. Traceability checkpoints prove mappings, not approval or implementation.

Enter the product completeness update through [requirements section 8](PRODUCT_REQUIREMENTS_EN.md#8-product-completeness-review-2026-09-06). Project proposal section 19 retains provenance; the update maps to contracts 9, architecture 6, security 6, operations 8, ADR-008, P1.7, and AT-16–AT-21. Keep 28 REQ IDs and 12 document pairs; take the union of PC-01–PC-06 supplements and primary mappings. S1 exit docs add onboarding, delivery packages, revocation, cost explanations, backup/restore/uninstall, and SBOM. S2/S3 enhancements remain candidates, not supported S1 features.