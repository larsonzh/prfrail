# ProofRail Documentation Index and Maintenance Plan

[简体中文](DOCUMENTATION_PLAN.md)

Date: 2026-09-06; pre-prototype documentation baseline awaiting human review. This is the first entry point for low-cost models, not a requirement to reread the full RFC every turn.

## 1. Authority and Reading Order

The [RFC](RFC-proofrail-unattended-ai-engineering-product.md) is the sole protocol authority. Revise/review it before schemas/goldens, tests, implementation and user docs. The long-term English-first goal does not displace the current Chinese RFC. English documents are synchronized translations; conflicts block work until corrected, never model-selected interpretations.

Read this page, the current T task in [DEV_PLAN_EN.md](DEV_PLAN_EN.md), its contract sections/AT and nearby code/tests. Read the relevant RFC/ADR only for ambiguity or authority expansion. Do not load both complete language sets each round.

## 2. Delivery Index and RFC 16.3 Coverage

| Package | Chinese/English | Coverage | Trigger/owner |
|---|---|---|---|
| Product | [需求](PRODUCT_REQUIREMENTS.md) / [Requirements](PRODUCT_REQUIREMENTS_EN.md) | Research, roles, use cases, scope, 28 functional/NFR items | Scope; product owner |
| Architecture | [架构](ARCHITECTURE.md) / [Architecture](ARCHITECTURE_EN.md) | Layers, modules, domain references, flow/transactions/storage | Boundaries/dependencies; architect |
| Protocols | [契约](CONTRACTS.md) / [Contracts](CONTRACTS_EN.md) | Schemas, snapshots/evidence, hooks, tickets/repair, adapters, handoff/docs | Wire/behavior; contract maintainer |
| Security | [安全](SECURITY.md) / [Security](SECURITY_EN.md) | STRIDE, permissions, secrets, audit, supply chain | Permissions/risks; independent reviewer |
| Usage/migration | [操作](OPERATIONS.md) / [Operations](OPERATIONS_EN.md) | Install/start/config, operations/debugging, whois mapping, self-host/release | CLI/deployment; implementer/release owner |
| Testing | [验证](TEST_STRATEGY.md) / [Testing](TEST_STRATEGY_EN.md) | 21 AT groups, tiers, platforms, faults, evidence | Requirements/contracts; test maintainer |
| Engineering | [规范](CODING_CONVENTIONS.md) / [Conventions](CODING_CONVENTIONS_EN.md) | Encoding, Go, Git, low-cost workflow | Tools/conventions; maintainer |
| Implementation | [计划](DEV_PLAN.md) / [Plan](DEV_PLAN_EN.md) | 24 tasks, dependencies, scope, traceability | Task/stage acceptance; implementation lead |
| Decisions | [ADR/readiness](ADR_REGISTER.md) / [Decisions](ADR_REGISTER_EN.md) | Trade-offs, bootstrap, release/support, 13 readiness gaps | Decisions/evidence; designated approver |
| Governance | [文档计划](DOCUMENTATION_PLAN.md) / this document | Bilingual rules, scoped reading, gates and future docs | Documentation changes; doc maintainer |

Reuse RFC 10/18 terminology instead of maintaining a second glossary. Combined operations/testing/migration packages reduce reading cost but give every P0 topic a home. Schemas and executable examples are still absent: do not claim all P0 specifications are frozen. Audit policy lives in security, release policy in operations/ADRs; external policies need approval before release.

## 3. Phased Documentation

| Phase | Deliverable | Gate |
|---|---|---|
| This D0 round | 10 Chinese/English pairs, README links, requirement/task checkpoints | Links, encoding, IDs, mappings; honest status |
| S0 T001-T004 | Sign-off, RFC field revision, machine schemas/positive-negative fixtures, spikes, dependency/signing ADRs | All RFC 17 gates and independent validation |
| Every S1 task | Synchronized code/tests/schema/affected bilingual sections; replace planned commands only after execution | Valid links/commands, generated freshness, relevant AT |
| S1 exit | Real installation tutorial, runnable generic/C/Go examples, recovery evidence, self-host/shadow/support matrix, release notes | Independent review, real platforms, owner approval |
| S2/S3 | New platform/composed-language/supervised/Web/script B docs and migrations; complete English RFC as needed | Contracts/tests first, not marketing as support evidence |

## 4. Low-Cost Model Task Card

Fill real values in the session before each task; do not create placeholder-filled task files:

```text
Task: Txxx; REQ: list; AT: list; Authorized scope: exact paths
Input: RFC section + contract section + current test; exclusions: task row
Hypothesis: one falsifiable cause; smallest check: one real command
Budget: time/calls/cost cap; no paid calls without authorization
Evidence: command/cwd/exit/case counts/artifact hash/unverified items
```

Change one testable slice at a time. Inspect the first failure; avoid recursive whole-repository searches, repeated long logs and invented frameworks. Carry task contract/parent hash/confirmed decisions/latest failure, not entire chat history. Label estimated token usage and never drop safety constraints to save tokens.

Escalate contract ambiguity, broader permissions, unauthorized paths, two same-error attempts without new evidence, budget exhaustion or unrecoverable state. Handover contains changes, verification and next action, not a repeated project description.

## 5. Review Rules

Chinese keeps original names, English uses _EN, with reciprocal links. REQ/AT/T/ADR/SEC and wire identifiers match exactly. Counts describe planned requirements/tasks, not completion percentages. Review both translations together, not a summary-only English edition. Resolve conflicts in the authoritative Chinese RFC then synchronize.

Distinguish fact, proposal, measurement and plan; unexecuted examples are planned. Link only existing files; describe future paths as planned code text. Markdown uses UTF-8 BOM/LF; machine canonical/wire data does not inherit that encoding. Frozen snippets require independent validation/execution; prose is not machine validation.

This round runs terminal checks. T003 consolidates document/contract checks into maintained tooling; T017 wires CI. Keep stage reports/open issues and increment checkboxes only after acceptance. Traceability checkpoints prove mappings, not approval or implementation.

Enter the product completeness update through [requirements section 8](PRODUCT_REQUIREMENTS_EN.md#8-product-completeness-review-2026-09-06). It maps to RFC 19, contracts 9, architecture 6, security 6, operations 8, ADR-008, P1.7 and AT-16–AT-21. Keep 28 REQ IDs and 10 document pairs; take the union of PC-01–PC-06 supplements and primary mappings. S1 exit docs add onboarding, delivery packages, revocation, cost explanations, backup/restore/uninstall and SBOM. S2/S3 enhancements remain candidates, not supported S1 features.