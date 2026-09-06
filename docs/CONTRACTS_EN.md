# ProofRail Architecture Contracts

[简体中文](CONTRACTS.md)

Date: 2026-09-06; S0 design-constraint draft. Authority: [RFC](RFC-proofrail-unattended-ai-engineering-product.md) sections 9-13 and 16-17. This document combines schema, snapshot/evidence, hooks, tickets/repair and adapters to reduce repeated reading.

## 1. Maturity and Compatibility

Requirements below come from the RFC; pending decisions block production implementation. T002 now provides the first structural Schemas, but this document is not a complete JSON Schema or a substitute for an independent validator. New wire fields/enums, canonical algorithms and error codes require an [ADR](ADR_REGISTER_EN.md), RFC revision, schemas and positive/negative goldens before Go implementation.

| Contract | Established semantics | Pending S0 artifact |
|---|---|---|
| Configuration/task/step | Draft 2020-12, `schemaVersion=1.0.0`, strict chain/task/four step kinds and workspace/component/language-scope topology | Harness/toolchain registries, references, merge rules and remaining schemas |
| change-set | Sequential in-memory prevalidation, first-fail-stop, markers/assertions, group transaction | Operations, line-ending semantics, before/after hashes, duplicate markers |
| snapshot/evidence | SHA-256 addressing, parent/candidate/evidence binding, immutability | Canonical paths/hash encoding, manifests/receipts |
| ticket/repair | Categories, fingerprint budgets, leases, Prepare/Inspect/Validate/Promote | Fingerprint algorithm, thresholds, ledger format, full transition table |
| adapter/context | Strict request/claim/result envelopes, recordHash, idempotency key, generation fencing, dispatch/takeover receipts and minimal context | Cross-record checker, signature receipts and platform atomicity probes |

Each persistent object has an independent `schemaVersion`, distinct from CLI and SessionBridge versions. The first version is `1.0.0`; unknown majors are rejected, and unknown minor/patch versions are read-only rejected by default. Repository, schema, runtime and wire JSON uniformly use UTF-8 without BOM plus LF. Canonical bytes additionally contain no formatting whitespace or trailing newline. Only a fixture explicitly testing BOM input compatibility may contain a BOM, and it cannot be reused as a business object.

## 2. Schema and Checker Responsibilities

| Object | Minimum semantics | Invalid/dynamic constraints |
|---|---|---|
| chain definition | `chain.id/profile`, at least one ordered task, minimal workspace, optional documentation | Cross-array task/step ID uniqueness, references and dependency cycles belong to the checker |
| run manifest | Frozen effective configuration, default origins and runtime bindings | No in-place plan edits during runs; exact fields await a later T002 slice |
| task | Stable id, nonempty steps, review and documentationPolicy | Duplicate IDs; skip-on-fail without failureIndependent |
| step | id, code/build/verify/noop; code execution defaults autonomous | Empty steps, unknown kind; noop requires reason and cannot start hooks/agents |
| hook | Independent strict registry; discriminated process/container runners; argument, failure, timeout, resource, network and artifact policies | Checker/preflight verifies IDs, references, tool digests and capabilities; execution belongs in receipts |
| target | Independent strict registry; stable id, class, access, unique path globs and one-way generated sources | Checker rejects escapes, overlapping writes, generation cycles and unresolved references |
| component/language scope | Component root; scope language/harness/toolchain/target IDs; dependencies; intersected step scopes | Unique IDs, resolved references, no cycles, available harness; shared roots allowed with unique writable targets |
| state event | Append-only envelope; chain/task/step entity; strict before/after states; sequence, predecessor hash, evidence and reason | Checker verifies cross-event sequence/hash/entity-state continuity, actor authorization and evidence existence |
| error | Strict category/code prefix, subject, evidence, redacted message, retryMode and suggestedAction | Checker/result contracts verify references, default retry restrictions, primary-error selection and concrete code catalog |
| documentation | required/if-affected/optional/forbidden and impactRules | Missing effective doc diff or stale generation; whitespace/mtime alone does not count |
| handoff | execution, scope, inputPolicy, deadline, returnActions, hooksAfterReturn | Single lease, stopped writers, fresh manifest; Schema alone cannot prove return checks |
| review/waiver | approve/reject/waive, actor, reason, policy, expiry, bound candidate | Self-approval, expiry, changed candidate or mismatched scope blocks |

Separate three validators: JSON Schema for structure; semantic checker for references, ownership, cycles and static policy; runtime gates for diffs, processes, capabilities, secrets, budgets and review. Each invalid golden identifies its rejection layer; JSON Schema cannot check live files or processes.

Eight structural Schemas for chain, hook, target, workspace, state-event, error, adapter-envelope and adapter-receipt now exist under `schemas/`; T003 still needs to create `testdata/contracts/valid/` and `testdata/contracts/invalid/` and execute them with an independent validator. Each fixture records ID, contract version, expected accept/reject, rejection layer and reason. Include three tasks, all four kinds, C+Go, C+JavaScript+Python, handoff and code/docs collaboration; at least one positive and negative per condition. Documentation paths/snippets are design inputs, not verified executable fixtures.

## 3. Configuration Resolution

Profiles provide defaults only; explicit task values override chain defaults; explicit step values override task defaults. List replacement/merge, empty values and workspace.toml versus proofrail.toml ownership await ADR-002. config explain reports values and origins. Silent array merging must not grant extra permissions.

After resolution, apply RFC 8785 JCS and hash the UTF-8 no-BOM bytes with SHA-256; represent the digest as `sha256:` plus 64 lowercase hexadecimal digits. Heartbeats, credentials and machine probes never rewrite user config. Bind parent snapshot, hook digests, adapter/harness versions and authorized policy. Changes require a new run or the RFC's pause/approval/new-manifest flow, never overwriting historical facts.

RFC section 16.9.3.1 freezes two byte-level vectors, `jcs-001` and `state-event-001`, covering JCS ordering/escaping/Unicode and the state-event domain-separated digest. T003 must still preserve them as independent fixtures and recompute them with a non-core implementation; documented digest values are not independent-validator evidence.

## 4. Snapshots, Evidence and Acceptance

`snapshot-manifest.schema.json` freezes baseline/candidate discrimination, parent and task/attempt binding, three path-entry
types, exclusions and secret-free environment facts. Regular-file hashes cover original bytes without BOM/EOL/encoding
conversion and explicitly mark package manifests, lockfiles, generated files and hardlink groups; directories preserve empty
directories, while symlink hashes cover the original relative target text. A manifest cannot declare itself accepted: only a
valid review plus promotion receipt authorizes downstream consumption. The independent semantic checker and cross-platform
goldens validate path ordering/uniqueness, Unicode/case collisions, parent/object closure, exclusion matches and restoration
feasibility. Git metadata is explicit evidence only, never a restoration source.

Events contain run/task, before/after state, monotonic sequence, time, actor, input evidence and reason. Hash chaining detects omission/reordering. Write temporary objects, verify and atomically publish; persist events before rebuildable projections. Torn-tail handling must not discard acknowledged facts and guess success.

Reader contract: the next task may consume a snapshot only when valid review binds parent, task/attempt, candidate and evidence root, and a completed publication receipt exists. Crashes between candidate/review/receipt/event/projection require journal replay. Directory names or a PASSED field alone prove nothing. Hashes prove integrity, not publisher identity.

`evidence-manifest.schema.json` freezes the evidence root as an immutable pre-review index for one task attempt. It binds
the task definition, resolved policy and parent/candidate snapshots, then references state-event, error, adapter, hook,
artifact, change-set, diff, ticket and handoff object hashes. It neither embeds logs nor infers success from object presence,
and it excludes review/waiver/promotion receipts; those objects point one-way to the evidence root to avoid a hash cycle.
The independent checker validates item ordering/uniqueness, same-attempt ownership, object existence, required-evidence
coverage and truthful redaction status.

`review-receipt.schema.json` is the only object allowed to reference `evidenceRootHash`, producing exactly one of
`approve/reject/waive`. `recordedBy.type` is restricted to `operator|policy`, excluding `agent|system` self-approval;
`reviewMode=policy` must bind an approved `policyHash`. `waive` must carry a human `authorizedBy`, a `policyBasisHash`,
a non-empty `scope` and a non-null `expiresAt`, never an anonymous automatic approval. The independent checker verifies
reviewer identity differs from the change producer, that `evidenceRootHash`/`policyHash` references actually exist, and
that an expired waiver is never reused.

managed-change-set: read all targets, simulate operations in global order, validate markers/prehash/assertions/boundaries, persist journal, atomically replace each file, verify all postconditions, then record receipt. Roll back the entire group on failure; failed rollback means quarantine/pause. Do not overwrite external changes or restore while processes are alive.

isolated-workspace: direct tool writes stay in a disposable run directory, with complete before/after manifests, diff, process and gate evidence. Acceptance remains whole-task: never publish selected languages/files or successful code with failed documentation.

## 5. Gate Hooks

Kinds are precheck/build/test/verify/review/cleanup, executed in order. executable plus args uses no implicit shell; cwd stays in the authorized workspace. Explicit environment allowlist, default-deny network and enforceable resource/output/time limits are required. shell/container/remote/PTY are declared capabilities: reject unsupported requirements instead of weakening policy.

Results prove hook identity/digest, start/end, exit status, stdout/stderr summary, artifact references/hashes, policy and failure reason; exact wire fields freeze in schemas. Startup failure, timeout, unstoppable processes, missing artifacts and scan failure differ from exit 0. warn continues only with explicit nonblocking policy and is not PASS; retry(n) spends shared budgets; cleanup failure cannot erase the primary failure.

Generated hook A uses locked templates, syntax/template/dependency/dangerous-capability checks, independent approval and hash-bound mounting. Disabled by default. The generator cannot alter its validator or approve itself. Arbitrary generated scripts B stay outside S1.

## 6. Tickets, Repair and Handoff

Categories are task-static/code-fix/noncode. Ticket context includes run/task/attempt, failure evidence, permitted actions, budgets and lease. The identical-fingerprint budget follows pending_review/override_window/hard_block; thresholds and valid-fix evidence must freeze before implementation. Models cannot reset budgets.

Prepare creates an isolated candidate and parent digest; Inspect checks scope/diff/ownership; Validate runs frozen gates; Promote rechecks termination, lease, parent/candidate/validation hashes and authorization, then publishes atomically with a promotion receipt. Changes invalidate stale downstream validation. Rejected review requires a new repair attempt, never editing an old receipt.

Handoff stops managed writers, flushes journals, captures manifest, enters WAITING_FOR_OPERATOR and grants an exclusive operator lease. complete revokes the lease, checks boundaries/secrets/unknown processes and runs gates; abort keeps failure evidence and discards this candidate; request-agent creates an attempt with confirmed conclusions. Expiry/disconnection never implies success or reassignment. Pause secret-direct without a safe terminal.

## 7. Agents and SessionBridge

The ProofRail port takes a versioned context envelope containing task contract, parent snapshot, confirmed decisions, current evidence, unresolved tickets, budgets, permitted tools and reference hashes. It returns candidate changes/evidence, not authoritative task PASS. Persist dispatch before sending; correlation is fixed within an attempt, and duplicate requests cannot reapply a change-set. Transport retries reuse requestId; a new business attempt gets new identity and cost accounting.

SessionBridge v0.1.1 is an independent dependency with a frozen public v1 contract. Explicitly use mode=silent and legacy=false, the same channel directory/target instance and requestId-bound receipts. Clear same-ID stale results before atomically writing commands. Verify exact wire data against its RFC, goldens and client implementation; explanatory RFC snippets may omit required fields.

| External result/capability | ProofRail action |
|---|---|
| status=ok response text | Validate output schema, checker, gates and review; never directly PASSED |
| busy | Bounded backoff, same requestId, waiting budget; serialize each target response slot |
| timeout/poll_timeout/disconnection | Record uncertain result and pause/reconcile; never assume nonexecution and resend indefinitely |
| lm_api_unavailable/unknown status/version | Explicit failure/pause; only preauthorized file-queue alternative, no auto/visible/clipboard fallback |
| Truncated history/host restart | Reconstruct facts from ProofRail envelope, not compressed history |
| Text-only silent | Structured managed-change-set candidate; no implied IDE editing/command tools |

SessionBridge's success cache is bounded in-memory state, not durable exactly-once. ProofRail persists dispatch/results and application idempotency. Complete-patch turns use noCompress according to the dependency contract plus independent completeness checks. History compression cannot bypass token/cost limits. ProofRail owns secret scanning/redaction; do not assume SessionBridge implements it.

File queue and IPC share ProofRail request/claim/result envelopes; CLI is a consumer form, not a third business protocol. Each queue file contains exactly one canonical JSONL record. An atomic request move into inflight acquires a claim; immutable generation files and `(claimId,generation)` fencing constrain renewal, takeover and result publication. Never append concurrently or overwrite formal files. Map SessionBridge transport IDs only in dispatch receipts; its single-slot files are not durable ProofRail facts. Takeover receipts bind old/new fencing, stopped-writer proof and authorization evidence; only a granted receipt may be referenced by a new claim. T004 still probes platform atomic/crash boundaries. Use a deterministic fixture consumer for no-IDE/no-AI testing, without paid calls.

## 8. Errors and Acceptance

Persistent errors distinguish schema/static/policy/capability/transport/execution/integrity/storage/review/operator/internal causes. `error.schema.json` freezes category/code prefixes, subjects, evidence, retry modes and suggested actions; CLI category exits are 10–20 and usage is 2. Never pass through SessionBridge 0/1/2/3 or hook/OS exits. Uncertain transport results require reconciliation, not blind retry; multi-error primary ordering awaits result/receipt Schemas. Errors never contain secrets.

Acceptance includes structural, semantic, runtime, crash, compatibility and security checks, not just parseable JSON. See [test strategy](TEST_STRATEGY_EN.md) and [development plan](DEV_PLAN_EN.md). No complete contract validator was generated or executed in this round; RFC gate 17.2 remains unmet.

## 9. Product and Lifecycle Contracts (RFC Section 19)

These are S0 semantics to freeze, not finalized wire fields. T002/T003 add versions, required fields, size limits, positive/negative fixtures and independent checks. Reject writes of unknown types.

| PC | Input/authority | Output/invariant | Acceptance |
|---|---|---|---|
| PC-01 | User configuration/static capabilities | Plan hash, permissions/budgets/unknowns; no hook/model/network/credential calls; version probes need separate consent | AT-16 |
| PC-02 | Accepted snapshot/review/complete references | Export binds parent/target/content hashes/exclusions; new directory only, verify temporary output before completion; no success receipt on failure | AT-17 |
| PC-03 | Trusted actor's scope/manifest/budget/expiry | Check before new effects/acceptance; revocation blocks and requests controlled stopping, never erases existing effects; waiver cannot bypass hard safety gates | AT-18 |
| PC-04 | Trusted effect policy/enforceable runner capabilities | Separate read-only/disposable local/external writes; reject S1 external writes, reconcile unknown outcomes; diagnostics do not mutate locks/logs | AT-19 |
| PC-05 | Owner budget/pricing or measurement source/call identity | Durable pre-call reserve, deduplicated settlement, unknown holds; shared caps across runs; restart cannot release unsettled reservations | AT-20 |
| PC-06 | Verified reference closure/support matrix/retention authority | Backup objects/events/references without secrets, restore to a new store, preserve old formats; uninstall keeps evidence/shared tools by default | AT-21 |

Cost records include currency, price version, estimated/observed/unknown and settlement evidence. Subscription call/token-cap authorization explicitly disclaims monetary-bill guarantees. Record reserved/actual differences; duplicate receipts cannot settle twice. Do not claim hard currency guarantees beyond provider control.

Export/backup readers verify reference closure, safe paths/types, versions, hashes and authority. Export completion does not change task/chain acceptance or imply product release. Reject live backups with uncertain consistency instead of reporting false recoverability. Disposition records exclude deleted secrets; deletion/retention conflicts require human decisions.