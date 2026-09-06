# ProofRail Product Requirements and Research

[简体中文](PRODUCT_REQUIREMENTS.md)

Date: 2026-09-06. Status: S0 review draft, not a release commitment. Authority: [RFC](RFC-proofrail-unattended-ai-engineering-product.md) sections 1, 9, 14 and 17; engineering rules: [CODING_CONVENTIONS.md](CODING_CONVENTIONS.md). REQ IDs trace implementation; they do not replace RFC R IDs.

## 1. Research Scope and Evidence

| Source | Finding | Limitation |
|---|---|---|
| Current user confirmation | whois v3.4.0 released and frozen; SessionBridge v0.1.1 released on GitHub/Gitee/Marketplace and installed | Accepted project input; releases were not independently audited online here |
| RFC sections 2 and 9 | Existing experience identifies runaway retries, baseline pollution, missing approvals and broken evidence chains | Experience, not external interviews or market statistics |
| Current repository | Go 1.22 module; version command plus init/run placeholders; 11 placeholder internal directories | No execution, isolation, recovery or verification implementation |
| Adjacent SessionBridge RFC sections 4 and 7 | Reusable silent text request/receipt transport, without ProofRail task or approval semantics | Installation proves neither model availability nor file-editing tools |

This is desk research and contract analysis, without paid model experiments, competitor trials or broad web research. Claims of market scarcity in the RFC are positioning hypotheses, not established competitive evidence.

| Alternative category | Useful capability | Control responsibility ProofRail must establish |
|---|---|---|
| IDE assistants/autonomous agents | Propose or perform edits | Independent acceptance, budgets, trusted baselines and recovery evidence; no blanket claim that vendors lack them |
| CI/script orchestration | Repeatable verification | Recoverable local long-running tasks, operator handoff and candidate acceptance |
| Git worktrees/commits | Versioning and collaboration | Current-tree snapshots and uncommitted-file protection without Git write permission |
| SessionBridge | Host messaging and conversation context | Task ownership, write leases, structured changes, gates, review and publication |

H1: Small tasks with deterministic acceptance reduce supervision. Compare interventions, elapsed time and causes against a fixed manual-script task set; promise no saving percentage before measurement.
H2: Minimal context envelopes reduce cost. Record input/output tokens, model, price source and total cost per accepted task. Label unavailable provider usage as estimates.
H3: Model independence lowers integration cost. Exercise the same contracts through a no-AI file queue and SessionBridge silent transport.

## 2. Users and Scope

Primary users are individual maintainers and small teams. The owner approves scope, budget and releases; operators run, pause, hand off and recover; independent reviewers approve/reject candidates; AI is a restricted execution actor. One person may hold maintenance roles, but the change producer cannot independently authorize acceptance.

S0 completes specifications, schemas/goldens and spikes, then obtains S1 approval. S1 delivers Windows 11 amd64 CLI/TUI, Linux amd64 core CI, generic/C/Go harnesses, both change modes, manual-handoff, code/documentation coordination and two-generation self-hosting. S2 adds production Linux, composed languages, Web and supervised execution; S3 adds arbitrary generated scripts, more editors and inception. P is priority, S is delivery stage.

This work performs inception activities for ProofRail itself; it does not implement the deferred R2 inception product. Do not modify frozen whois, released SessionBridge or Git state. No chat system, GUI injection, cloud service or database is introduced.

## 3. User Scenarios and Testing

All scenarios are required for S1. Ordering represents implementation dependencies, not reduced priority.

| Story | Given / When / Then | Value |
|---|---|---|
| US1 Safe start | Given uncommitted files, capture baseline and run three tasks; source bytes stay unchanged and tasks consume accepted snapshots only | Protect daily work |
| US2 Controlled failure | Given a blocking T2 failure, stop, archive and pause; T3 never starts and recovery never consumes the failed candidate | Prevent cascading pollution |
| US3 Human participation | Given manual-handoff, stop writers before granting a lease; complete triggers checks and gates, not automatic PASS | Support real debugging |
| US4 Aggregate review | Given a CLI change matching documentation rules, missing docs or broken examples block the whole task; independent acceptance publishes both together | Prevent contract drift |
| US5 Offline verification | Given an archived run, verify events, objects and receipts offline; tampering, omission and reordering cannot PASS | Support audit and handover |
| US6 Self-hosting | Given an externally audited seed, build and replay an isolated candidate; external-oracle failure leaves seed and production unchanged | Avoid circular trust |

Edge cases include empty tasks/steps, unknown fields/versions, path escape, case collisions, shared targets, generation cycles, lease conflicts, partial writes, uncertain process liveness, disk exhaustion, unknown cost, model timeout, expired approval and unavailable secret-direct. Every case needs explicit rejection/pause evidence, never guessed success.

## 4. Functional Requirements

| ID | RFC | Stage | Testable requirement |
|---|---|---|---|
| REQ-001 | R1 | S0/S1 | No whois A/B or D/V semantics in core contracts; a non-C example passes |
| REQ-002 | R2 | S1 | Replayable baseline-to-report execution, excluding inception |
| REQ-003 | R3 | S1 | Ordered, resumable chains of three or more tasks; no fixed schema task-count ceiling |
| REQ-004 | R4 | S1 | Read-only source capture; accepted-snapshot propagation/recovery without Git; reject nonportable paths |
| REQ-005 | R5 | S1 | Technical success enters REVIEW_PENDING; independent approval/valid waiver is required for PASSED |
| REQ-006 | R6 | S1-S3 | S1 Windows production/Linux core CI; static model selection and preauthorized fallback; pause on missing capability |
| REQ-007 | R7 | S1 | IPC and file queue share ProofRail receipts; SessionBridge silent only, no GUI/auto fallback |
| REQ-008 | R8 | S1 | Rebuild context envelopes after host restart; conversation history is not authoritative |
| REQ-009 | R9 | S1 | Unified TUI for tasks, gates, tickets, logs and budgets; no separate monitor windows |
| REQ-010 | R10 | S1 | init/validate/config explain generate, validate and explain configuration; runtime state never rewrites it |
| REQ-011 | R11 | S1/S2 | External generic/C/Go harnesses in S1, no language branching in core; extend languages in S2 |
| REQ-012 | R12 | S1 | Separate executable/args, preflight capabilities, enforce timeout/output/network/environment policy |
| REQ-013 | R13 | S1/S3 | S1 template-generated hooks only, disabled by default, hash-bound and independently checked/approved; arbitrary scripts in S3 |
| REQ-014 | R14 | S0-S3 | Apply RFC section 14 exit gates; iteration demos are not stage completion |
| REQ-015 | R15 | S1 | Pure Go, Go 1.22 compatibility baseline, CGO_ENABLED=0 releases and TUI; review dependencies first |
| REQ-016 | R16 | S1/S2 | Complete CLI flow without IDE; Web/extension views gain no extra state-write authority |
| REQ-017 | R17 | S0 | Complete P0 documents, schemas/goldens, spikes and authorization before S1 |
| REQ-018 | R18 | S1 | Explicit ordered code/build/verify/noop; noop has a reason and receipt and starts no process |
| REQ-019 | R19 | S1 | Seed builds candidate, external oracle, isolated clean-room replay and independent release approval |
| REQ-020 | R20 | S0-S2 | S0 dual/triple-language configs; S1 C+Go contracts; S2 three-language execution with whole-task failure |
| REQ-021 | R21 | S1/S2 | S1 single-writer handoff, structured non-secret input, timeout/disconnection recovery and return checks; S2 supervised |
| REQ-022 | R22 | S1/S2 | First-class docs, deterministic impact rules, generation direction, waivers and aggregate acceptance; S2 intelligent hints |

## 5. Nonfunctional Requirements

| ID | Source | Measurable requirement |
|---|---|---|
| REQ-023 | RFC section 12 | All mandatory adversarial cases block; boundary violations, missing evidence and forged PASS never advance a chain |
| REQ-024 | RFC sections 10.3-10.4 | Faults at every durability boundary lead to replay or explicit pause; accepted references are never lost |
| REQ-025 | RFC section 12.3 | Offline tamper/omission/reordering detection; secret canaries absent from archives/model context |
| REQ-026 | RFC sections 9.3, 12.2 | Check wall time, attempts, fingerprints, storage and call/cost budgets before execution; pause safely on exhaustion |
| REQ-027 | RFC section 13.3 | Stable parseable --json; keyboard, no-color and narrow-terminal state/pause access |
| REQ-028 | RFC sections 6, 16.2 | Reject writes with unknown schema major versions; read-only historical verification; backup and compatibility checks before upgrade |

Measurement plan, not measured performance: use fixed non-secret fixtures of 10,000 files/100 MiB and 100 noop tasks. Record CPU/RAM/disk/OS, capture/replay time, peak memory and evidence size. First successful measurements establish the baseline; subsequent same-machine median regressions above 20% trigger investigation. Never delete evidence for speed or promise unmeasured cross-platform latency.

## 6. Entities and Success Criteria

A chain owns ordered tasks; tasks own steps. Components/language scopes own targets. Snapshots define content ancestry; hooks produce evidence; tickets/attempts bound repair; review/promotion receipts prove acceptance; handoffs record write leases. Exact terminology follows RFC sections 10 and 18.

S1 succeeds only when US1-US10 (including section 8 additions) and all S1 mandatory negative cases pass, source Git state/content remain unwritten by the engine, failures/rejections block advancement, both change modes and handoffs recover, native Windows and Linux core CI evidence exists, and external self-hosting acceptance completes. Cross-compilation, an empty passing test suite or an AI completion claim is not product acceptance.

## 7. Approval and Further Research

The owner must approve budget measurement modes/caps, per-call limits, maintainers, independent security reviewers, signing identity and support duration. Until then, no paid calls, dependency upgrades or S0 completion claim. Monetary authorization requires bounded worst-case reservations and prohibits automatic calls without a bound. Subscription mode requires separately approved call/token hard caps and acknowledgment that monetary bills are not guaranteed; see RFC section 19.5.

One structured S0 interview confirms the three most common tasks, acceptable intervention points, evidence retention and machine limits. S1 shadow records validate H1-H3 and feed the stage report. Market size, willingness to pay and competitor performance do not block the safety loop but require separate evidence before promotional claims.

## 8. Product Completeness Review (2026-09-06)

An execution engine alone is not a usable product. RFC section 19 adds acceptance detail while retaining 28 REQ IDs; feature counts are not market evidence. PC IDs are product slices of existing requirements, not a separate protocol authority.

| Slice | Business gap and user outcome | REQ | S1 minimum / later boundary | Acceptance |
|---|---|---|---|---|
| PC-01 | New users cannot assess permissions, cost or readiness | REQ-009/010/016/027 | No-AI example, side-effect-free preview, unknown reasons/next action; S2 graphical onboarding/localization | AT-16 |
| PC-02 | Internal snapshots are not usable delivery | REQ-002/004/025 | Export accepted snapshots into a new directory with offline verification; no source writes/releases; S2 external PR integration evaluation | AT-17 |
| PC-03 | Revocation and pending approvals lack an operational flow | REQ-005/021/023 | Hash-bound run authorization, revocation/expiry blocking and approval inbox; S2 redacted notifications/RBAC | AT-18 |
| PC-04 | Rollback can be mistaken for undoing every side effect | REQ-012/023/024 | Classify effects, default-deny external writes, recovery plans/redacted diagnostics; S2 compensation separately reviewed | AT-19 |
| PC-05 | Timeouts/retries may be charged; costs/value are unclear | REQ-006/026 | Durable reservation/settlement, unknown holds, shared caps, local reports without default telemetry; S2 trends | AT-20 |
| PC-06 | Upgrade, backup, retirement and retention conflicts lack a complete flow | REQ-019/025/028 | SBOM/support matrix, complete backup closure/restore drill, evidence-preserving uninstall; S2 explicit migration | AT-21 |

New stories: US7 preview before authorization with no command/model calls; US8 export accepted results and verify in an empty directory without changing source; US9 revoke authorization to block new side effects while retaining in-flight stop evidence; US10 read historical evidence from verified backups after upgrade failure/retirement. These map to AT-16, AT-17, AT-18 and AT-21; AT-19/20 cover side-effect/cost negatives.

S1 minimum details join the revised exit gate, with renewed budget approval under RFC 17.8 before implementation. Keep S2 enhancements out of S1. S3 cloud services, commercial billing and team administration require demonstrated demand; no commercial backend is designed here. Local installation, trial, recovery, delivery and retirement define initial product usability.

Unknown prices are not free. Subscription providers without per-call pricing support only user-authorized call/token caps and an explicit monetary-bill limitation. This refines REQ-026 measurement modes without waiving authorization/budgets. Evaluate accepted-task total cost, first completion time, interventions and recovery success locally on fixed fixtures/machines; no measured results exist yet.