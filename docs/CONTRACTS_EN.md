# ProofRail Architecture Contracts

[简体中文](CONTRACTS.md)

Date: 2026-09-07; S1 architecture-contract baseline. This document is the normative authority for wire formats, schemas, state transitions, events, receipts, and cross-record invariants. The [project proposal](RFC-proofrail-unattended-ai-engineering-product.md) retains provenance, scope, and design rationale. This document combines schema, snapshot/evidence, hooks, tickets/repair, and adapters to reduce repeated reading.

## 1. Maturity and Compatibility

This document establishes the semantics below while retaining project-proposal provenance; pending decisions block production implementation. T002 now provides the first structural schemas, but this document is not a complete JSON Schema or a substitute for an independent validator. New wire fields/enums, canonical algorithms, and error codes require an [ADR](ADR_REGISTER_EN.md), a revision here, schemas, and positive/negative goldens before Go implementation.

| Contract | Established semantics | Pending S0 artifact |
|---|---|---|
| Configuration/task/step | Draft 2020-12, strict chain/task/four step kinds, workspace topology, harness/toolchain registries and deterministic ownership/merge | Independent positive/negative fixtures |
| change-set | Five strict operations, sequential in-memory prevalidation, first-fail-stop, markers/exact assertions, raw-byte hashes and group transaction | Cross-record checker and positive/negative fixtures |
| snapshot/evidence | SHA-256 addressing, parent/candidate/evidence binding, immutability | Canonical paths/hash encoding, manifests/receipts |
| ticket/repair | Stable fingerprint, append-only ledger and three-phase budget; strict Prepare/Inspect/Validate/Promote chain | Cross-record checker and positive/negative fixtures |
| adapter/context/agent-runner | Strict request/claim/result envelopes, recordHash, idempotency key, generation fencing, dispatch/takeover receipts and minimal context; CLI Agent process, session, capability and result mapping | Existing adapter trust-chain fixtures; AgentRunner Schemas, positive/negative fixtures and real CLI capability probes remain to freeze |
| product/lifecycle | Plan preview, export, authorization/revocation, effect/recovery, cost ledger and release/backup/retirement records | Independent positive/negative fixtures and S1 runtime acceptance |

Each persistent object has an independent `schemaVersion`, distinct from CLI and SessionBridge versions. The first version is `1.0.0`; unknown majors are rejected, and unknown minor/patch versions are read-only rejected by default. Repository, schema, runtime and wire JSON uniformly use UTF-8 without BOM plus LF. Canonical bytes additionally contain no formatting whitespace or trailing newline. Only a fixture explicitly testing BOM input compatibility may contain a BOM, and it cannot be reused as a business object.

## 2. Schema and Checker Responsibilities

| Object | Minimum semantics | Invalid/dynamic constraints |
|---|---|---|
| chain definition | `chain.id/profile`, at least one ordered task, minimal workspace, optional documentation | Cross-array task/step ID uniqueness, references and dependency cycles belong to the checker |
| run manifest | Input/effective chain digests, per-leaf origins, policy/parent snapshot and versioned runtime bindings | Checker verifies pointers/source closure, loaded objects, capability proof and one manifest per run |
| task | Stable id, nonempty steps, review and documentationPolicy | Duplicate IDs; skip-on-fail without failureIndependent |
| step | id, code/build/verify/noop; code execution defaults autonomous | Empty steps, unknown kind; noop requires reason and cannot start hooks/agents |
| hook | Independent strict registry; discriminated process/container runners; argument, failure, timeout, resource, network and artifact policies | Checker/preflight verifies IDs, references, tool digests and capabilities; execution belongs in receipts |
| harness/toolchain | Harness composes languages, parameters, hooks and artifacts; toolchain declares source, platforms, probe and evidence policy | Checker/preflight verifies unique references, versions, source, hashes, authorization and actual availability |
| target | Independent strict registry; stable id, class, access, unique path globs and one-way generated sources | Checker rejects escapes, overlapping writes, generation cycles and unresolved references |
| component/language scope | Component root; scope language/harness/toolchain/target IDs; dependencies; intersected step scopes | Unique IDs, resolved references, no cycles, available harness; shared roots allowed with unique writable targets |
| state event | Append-only envelope; chain/task/step entity; strict before/after states; sequence, predecessor hash, evidence and reason | Checker verifies cross-event sequence/hash/entity-state continuity, actor authorization and evidence existence |
| error/error-set | 37 closed codes with retry/action pairs; non-empty error index and primaryErrorId | Checker verifies references, uniqueness and time/category/code/errorId primary ordering |
| verification report | Sixteen fixed-order cross-record/trust-chain checks; passed/failed/incomplete | Checker rereads objects and verifies hashes, applicable checks, order and aggregate outcome; reports do not change state |
| product records | Strict preview/export/authorization/effect/cost/lifecycle discriminated records | Checker verifies authorization time, reference closure, arithmetic, platform facts and the S1 external-write prohibition |
| documentation | required/if-affected/optional/forbidden and impactRules | Missing effective doc diff or stale generation; whitespace/mtime alone does not count |
| handoff | execution, `handoffPolicy` (allowedTargets referencing read-write targets, inputPolicy, handoffTimeoutMs, returnActions, hooksAfterReturn) | Single lease, stopped writers, fresh manifest; Schema alone cannot prove return checks |
| review/waiver | approve/reject/waive, actor, reason, policy, expiry, bound candidate | Self-approval, expiry, changed candidate or mismatched scope blocks |

Separate three validators: JSON Schema for structure; semantic checker for references, ownership, cycles and static policy; runtime gates for diffs, processes, capabilities, secrets, budgets and review. Each invalid golden identifies its rejection layer; JSON Schema cannot check live files or processes.

Twenty-nine structural Schemas for chain, hook, target, workspace, state-event, error, adapter-envelope, adapter-receipt,
snapshot-manifest, evidence-manifest, review-receipt, promotion-receipt, handoff-receipt, hook-result, run-manifest,
signature-receipt, error-set, ticket-ledger, repair-transaction, harness, toolchain, change-set, verification-report,
plan-preview, export-record, authorization-record, effect-record, cost-ledger and lifecycle-record now exist under `schemas/`;
T003 created `testdata/contracts/valid/` and `testdata/contracts/invalid/` and executed them with the independent validator under `tools/contracts/`. Each
fixture records ID, contract version, expected accept/reject, rejection layer and reason. Include three tasks, all four kinds,
C+Go, C+JavaScript+Python, handoff and code/docs collaboration; at least one positive and negative per condition.
Documentation paths/snippets are design inputs, not verified executable fixtures.

### 2.1 Chain State, Events, and Recovery (T009 Authority)

This section is the normative source for T009 chain/task/step states, ordered scheduling, event projections, and recovery behavior. Project proposal sections 10.3 and 16.9.8 remain design sources. Any conflict blocks implementation and requires a documentation correction; implementers must not choose between them.

Schemas accept only the following single-step transitions. Unlisted transitions, self-transitions, and in-place retries across attempts are rejected:

- chain: `NONE→CREATED`; `CREATED→BASELINED|FAILED|CANCELLED`;
	`BASELINED→RUNNING|FAILED|CANCELLED`; `RUNNING→PAUSED|COMPLETED|FAILED|CANCELLED`;
	`PAUSED→RUNNING|FAILED|CANCELLED`.
- task: `NONE→PENDING`; `PENDING→PRECHECK|FAILED|CANCELLED`;
	`PRECHECK→STEPS_RUNNING|FAILED|CANCELLED`;
	`STEPS_RUNNING→WAITING_FOR_OPERATOR|REVIEW_PENDING|FAILED|CANCELLED`;
	`WAITING_FOR_OPERATOR→STEPS_RUNNING|FAILED|CANCELLED`;
	`REVIEW_PENDING→PASSED|REPAIR_PENDING|FAILED|CANCELLED`; `FAILED→REPAIR_PENDING`;
	`REPAIR_PENDING→STEPS_RUNNING|FAILED|CANCELLED`.
- step: `NONE→PENDING`; `PENDING→RUNNING|NOOP_RECORDED|FAILED|CANCELLED`;
	`RUNNING→WAITING_FOR_OPERATOR|PASSED|FAILED|CANCELLED`;
	`WAITING_FOR_OPERATOR→RUNNING|FAILED|CANCELLED`.

`NONE` is a creation-event sentinel, never a stable projection state. `COMPLETED`, `PASSED`, `NOOP_RECORDED`, and `CANCELLED` are terminal. Only task `FAILED→REPAIR_PENDING` can leave a failed state; repairing a step increments its attempt.

- Persist each append-only state event before updating its disposable, rebuildable projection. Events use a contiguous per-run sequence starting at 1 and previousEventHash linkage, and bind entity, before/after state, actor, evidence, and reason.
- The first run event is chain `NONE→CREATED`. Append `CREATED→BASELINED` only after baseline-0 and its evidence are durable. Directory names, file presence, and process exit codes never imply state.
- Schedule tasks in definition order. Materialize the first task from baseline-0. Materialize each later task only from the previous task's accepted snapshot backed by a valid review and completed promotion receipt. Candidate, failed, or uncertain snapshots cannot parent later tasks.
- Execute task steps by sequence. A `noop` only performs `PENDING→NOOP_RECORDED`, records a reason, and starts no agent, hook, or child process. T009 uses fake agents/runners without creating an alternate event or state protocol.
- `pause` takes effect at the current atomic boundary and starts no new step. `resume` continues only from a listed paused state. `cancel` stops the managed process tree and archives evidence before recording `CANCELLED`; the same run cannot resume.
- Recovery first proves the old writer inactive and acquires a valid fencing token, then resolves incomplete writes from journals, replays the complete event chain, and rebuilds projections. Torn tails, sequence/hash/state discontinuity, unknown before/after content, or uncertain external effects remain `PAUSED` or `REPAIR_PENDING`; never guess success, redispatch, or publish a candidate.
- Heartbeats, progress, diagnostics, review, handoff, takeover, and promotion receipts are not state transitions. Schema validates one record's shape and state pair; the checker validates chain continuity, current projection, unique genesis, actor authority, evidence existence, and scheduling preconditions.

## 3. Configuration Resolution

`proofrail.toml` exclusively owns chain/task/step, documentation, policy and registry references; `workspace.toml` owns
locations, topology, platform and adapter/harness/toolchain selection. Duplicate or unknown keys are rejected. Precedence is
builtin-default < profile < chain < task < step: absence inherits, scalars replace, known objects merge recursively and arrays
replace wholesale, never concatenate. Empty arrays mean explicitly none only where Schema permits; empty strings are not
unset and null deletion is unsupported. Every effective leaf has exactly one run-manifest pointer origin. `config explain`
is read-only, showing values/redacted references and override chains without execution or writeback.

After resolution, apply RFC 8785 JCS and hash the UTF-8 no-BOM bytes with SHA-256; represent the digest as `sha256:` plus 64 lowercase hexadecimal digits. Heartbeats, credentials and machine probes never rewrite user config. Bind parent snapshot, hook digests, adapter/harness versions and authorized policy. Changes require a new run or this document's pause/approval/new-manifest flow, never overwriting historical facts.

`run-manifest.schema.json` separates the validated input definition from the fully default-expanded chain through
`chainDefinitionHash` and `effectiveChainHash`. `resolution` uses RFC 6901 pointers to attribute every final leaf to a
builtin/profile/chain/task/step source; `bindings` freezes adapter/harness/toolchain/hook/policy versions, object digests
and capability evidence. Absolute paths, credentials, environment values and mutable machine probes stay outside. The
checker owns pointer existence, complete source coverage, loaded-object matching and one manifest per runId.

`harness.schema.json` composes languages, toolchains, declarative parameters, registered hooks and artifact globs without
copying runners or commands. `toolchain.schema.json` declares source, platforms, a separated executable/args probe and
evidence policy. Registry presence does not prove installation; checker/preflight verifies references, versions, source,
hashes and actual capability, and probing still requires authorization. Bundled generic/C/Go support does not close the
core language set. Environment values, credentials and absolute host paths stay out of registries and run manifests.

RFC section 16.9.3.1 freezes the `jcs-001` and `state-event-001` byte-level vectors, covering JCS ordering/escaping/Unicode and the state-event domain-separated digest; the later `agent-runner-enforcement-001` vector pins the enforcement domain digest. A non-core RFC 8785 implementation reproduces the independent fixtures' canonical bytes and digests.

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
artifact, change-set, diff, ticket, repair, authorization, effect, cost and handoff object hashes. Verification/export/
lifecycle records created after review cannot point backward into this root. It neither embeds logs nor infers success from object presence,
and it excludes review/waiver/promotion receipts; those objects point one-way to the evidence root to avoid a hash cycle.
The independent checker validates item ordering/uniqueness, same-attempt ownership, object existence, required-evidence
coverage and truthful redaction status.

`review-receipt.schema.json` is the only object that produces a review verdict; it directly references `evidenceRootHash` and produces exactly one of
`approve/reject/waive`. `recordedBy.type` is restricted to `operator|policy`, excluding `agent|system` self-approval;
`reviewMode=policy` must bind an approved `policyHash`. `waive` must carry a human `authorizedBy`, a `policyBasisHash`,
a non-empty `scope` and a non-null `expiresAt`, never an anonymous automatic approval. The independent checker verifies
reviewer identity differs from the change producer, that `evidenceRootHash`/`policyHash` references actually exist, and
that an expired waiver is never reused.

`promotion-receipt.schema.json` binds parent/candidate snapshots, evidence root, review receipt, writer-stop proof and lease
proof into one publication transaction. A candidate becomes consumable only for outcome=completed, when
acceptedSnapshotHash equals the candidate hash and review is approve or a still-valid waive. Failed/uncertain outcomes keep
acceptedSnapshotHash null and block advancement. Journal/checker logic verifies publication ordering, one completed receipt
per attempt, crash reconciliation and the accepted-reference closure.

`handoff-receipt.schema.json` freezes the completion fact of one `manual-handoff` return: it binds `handoffPolicyHash`,
operator identity, lease evidence and the before/after manifest and diff. `complete` requires non-empty
`hookResultEvidence`; `abort/request-agent` require it empty. A failed boundary/secret/unknown-process scan never produces
this receipt; it produces an `error` record and keeps the handoff paused instead.

managed-change-set: read all targets, simulate operations in global order, validate markers/prehash/assertions/boundaries, persist journal, atomically replace each file, verify all postconditions, then record receipt. Roll back the entire group on failure; failed rollback means quarantine/pause. Do not overwrite external changes or restore while processes are alive.

`change-set.schema.json` closes operations to create/delete/replace-exact/insert-before/insert-after and binds the task
attempt, parent snapshot, before/after manifests, per-operation raw-byte hashes, exact assertions and line-ending policy.
Insert markers must move from zero occurrences to exactly one; optional replace markers obey the same rule. Existing files
preserve one existing line-ending style while creates explicitly select LF/CRLF. Mixed endings, fuzzy/regex matching,
sequence gaps, duplicate markers or any hash mismatch are rejected before writes.

isolated-workspace: direct tool writes stay in a disposable run directory, with complete before/after manifests, diff, process and gate evidence. Acceptance remains whole-task: never publish selected languages/files or successful code with failed documentation.

## 5. Gate Hooks

Kinds are precheck/build/test/verify/review/cleanup, executed in order. executable plus args uses no implicit shell; cwd stays in the authorized workspace. Explicit environment allowlist, default-deny network and enforceable resource/output/time limits are required. shell/container/remote/PTY are declared capabilities: reject unsupported requirements instead of weakening policy.

`hook-result.schema.json` separates `executionOutcome` (runner fact), `assessment` (gate verdict) and
`policyDisposition` (subsequent control flow), while binding the hook definition, execution ordinal, timing, exit code,
stdout/stderr, artifacts, runner/termination proof and error evidence. Startup failure, timeout, resource enforcement,
unstoppable processes, missing artifacts and scan failure have distinct failureKind values. warn remains failed and is not
PASS; retry spends shared budgets; cleanup failure cannot erase the primary failure. The checker/runtime owns cross-record
references, actual artifacts, sufficient scan/termination proof and retry budgets.

Generated hook A uses locked templates, syntax/template/dependency/dangerous-capability checks, independent approval and hash-bound mounting. Disabled by default. The generator cannot alter its validator or approve itself. Arbitrary generated scripts B stay outside S1.

## 6. Tickets, Repair and Handoff

Categories are task-static/code-fix/noncode. Ticket context includes run/task/attempt, failure evidence, permitted actions, budgets and lease. The identical-fingerprint budget follows pending_review/override_window/hard_block; thresholds and valid-fix evidence must freeze before implementation. Models cannot reset budgets.

Prepare creates an isolated candidate and parent digest; Inspect checks scope/diff/ownership; Validate runs frozen gates; Promote rechecks termination, lease, parent/candidate/validation hashes and authorization, then publishes atomically with a promotion receipt. Changes invalidate stale downstream validation. Rejected review requires a new repair attempt, never editing an old receipt.

Handoff stops managed writers, flushes journals, captures manifest, enters WAITING_FOR_OPERATOR and grants an exclusive operator lease. complete revokes the lease, checks boundaries/secrets/unknown processes and runs gates; abort keeps failure evidence and discards this candidate; request-agent creates an attempt with confirmed conclusions. Expiry/disconnection never implies success or reassignment. Pause secret-direct without a safe terminal.

## 7. AgentRunner and Auxiliary Messaging Boundaries

The ProofRail port takes a versioned context envelope containing the task contract, parent snapshot, confirmed decisions, latest failure evidence, unresolved tickets, budgets, authorized targets, allowed effects and reference hashes. It returns external execution facts and candidate references, never authoritative task PASS. An `agent-runner-request` identifies a call by `requestId/runId/taskId/stepId/attempt`, binds immutable input through `workspaceHash/contextHash/parentSnapshotHash/authorizationHash/budgetHash`, and explicitly lists `allowedTargets/allowedEffects`. `mode=create` cannot carry a prior session; its first event establishes the request's sole session binding. `mode=resume` must bind `priorSessionId` and `priorCompletionHash`. Persist the request and authorization digest before launch. The same requestId and digest may replay, while a different digest conflicts; duplicate requests must not restart an agent in an unknown state or reapply the same candidate. A new business attempt receives new identity and cost accounting. Resume binds the original attempt, session and workspace; when continuity cannot be proven, reconstruct from the durable envelope as a new attempt rather than guessing that resume succeeded.

Before any AI-assisted entry point starts, ProofRail must perform a fail-closed availability preflight for the provider profile and channel that will actually be used. A profile stores only provider type, base URL, model, transport options, and an opaque `secretRef`; API keys and tokens cannot enter configuration, digests, logs, errors, or evidence, and may be resolved only from the system SecretStore for temporary injection into the child-process environment. `agent-runner-cli`, `sessionbridge-silent`, and `sessionbridge-visible` each produce an independent `available`, `unavailable`, or `unknown` verdict that cannot be inherited by another channel. Valid configuration, an existing secret, an installed CLI/extension, reachable IPC, or successful visible delivery is insufficient to prove AI availability. `available` requires a fresh, evidence-backed minimal live response on that channel, bound to the profile configuration digest, channel, and probe time, and executed within an explicit cost/request budget; stale, timed-out, skipped, or incomplete evidence blocks. Failure reasons include at least `configuration_invalid`, `secret_missing`, `credential_rejected`, `account_unavailable`, `quota_exhausted`, `network_unreachable`, `tls_or_proxy_failure`, `provider_unavailable`, `model_unavailable`, `policy_blocked`, `cli_unavailable`, `sessionbridge_unavailable`, `lm_api_unavailable`, `visible_delivery_unavailable`, `response_invalid`, `timeout`, and `unknown`. An unavailable channel cannot automatically switch provider, model, or channel. Non-AI features may continue, while the affected AI entry point remains disabled and displays its structured reason.

A `windows-credential:<target>` reference uses a generic credential in the current user's Windows Credential Manager with local-machine persistence; it is neither runtime isolation nor proof of provider availability. `prfrail secret set` may read a secret only from a verified local console with echo disabled; argv, redirection, pipes, JSON request bodies, and ordinary stdin injection are forbidden. `status` returns only the reference and existence, `delete` is idempotent, and none of the three commands displays the secret. An adapter receives only a temporary resolved byte copy and clears it after use. System failures return stable classifications without credential-blob content. Non-Windows platforms fail closed until their own SecretStore backend exists; they do not share or infer Windows credentials.

The formal AI execution channel is a CLI Agent adapter implementing the `AgentRunner` port. An adapter maps the unified request to an external CLI pinned by version and digest; it never promotes vendor transcripts, private session formats or exit codes directly into core state. Before launch it must probe and record normalized `os`/`arch`, the CLI platform-package/build fingerprint `platformHash`, executable identity/version, noninteractive invocation, working-directory binding, machine-readable events or complete logs, session create/resume, cancellation and process-tree termination, tool/network/permission controls, cost usage, and unattended confirmation behavior. Candidate identity includes at least version, `os`, `arch`, `platformHash`, `executableHash`, and `configHash`; another platform running the same version cannot inherit its conclusions. `platformHash` alone does not prove a kernel, filesystem, or OS enforcement boundary; enforcement configuration and runtime evidence retain that responsibility. A missing required capability blocks preflight. Prompt promises cannot substitute for OS or runner enforcement.

A CLI capability record remains a statement about the candidate itself; external compensation never rewrites `unsupported` as `verified`. Preflight first normalizes the current environment to `os`/`arch` and computes `platformHash`, then requires an exact three-field match with the capability; even a compatible disposition cannot bypass platform validation. When an external execution boundary compensates tool or network control, a separate `agent-runner-enforcement` record must bind the same candidate by capability `recordHash`, `executableHash`, and `configHash`; require the enforcement, capability, and current evaluation environment to share one `platformHash`; carry runtime evidence for `toolControl` and `networkControl`; and use the `proofrail:agent-runner-enforcement:1\n` domain separator. Combined preflight admits execution only when the candidate's other ten entries are verified and each control is verified either by the candidate or by the bound external boundary. Enforcement `probedAt` must not precede capability `probedAt`, follow the evaluation time, or exceed the admission policy's maximum enforcement-evidence age. Missing records, wrong candidate/config/platform bindings, unknown or unsupported findings, duplicate or stale evidence, and any control without a verified source all block. Risk acknowledgment, directory separation, prompt instructions, and post-run scans cannot substitute for this record. External enforcement compensates only its explicit controls, does not change the capability disposition, and never proves task PASS.

The current AgentRunner product-development and delivery scope is `windows/amd64` only. After the Windows version is complete, `linux/amd64`, using an Ubuntu 24.04 VM or an equivalently proven environment, is the next independent development and probe target and remains unsupported until its full capability, enforcement, and AT-23 evidence pass. macOS, ARM, and other OS/architecture combinations retain extension points but have no committed schedule or support status; later probes may use native environments or GitHub Actions runners. A VM, emulator, or CI runner may provide complete enforcement evidence for a target platform only when target-kernel, architecture, filesystem, process, permission, and network-boundary equivalence is proven; otherwise it provides development feedback or capability-probe evidence only and cannot produce a compatible verdict. Narrow platform interfaces isolate platform-specific code, while shared wire formats, state machines, and fail-closed semantics remain free of Windows command and path assumptions.

The CLI Agent works only inside the disposable run-workspace materialized for the current task/attempt. It cannot write the source directory, store, policy, gate implementation or acceptance records. Within authorized scope it may search/read files, directly modify the isolated workspace, run allowed development tools and iterate on diagnostics. It cannot commit, push, publish, expand network/targets, approve its candidate or delete evidence. ProofRail records redacted launch argv, cwd, environment-allowlist digest, agent/config/version digests, session ID, process identity, events/logs, time/call/cost usage, stop evidence and before/after manifests/diff. Without machine events, raw stdout/stderr is opaque evidence and cannot prove a complete tool trace.

Agent-reported completion, exit zero or final text means only that external execution ended. ProofRail must prove process-tree termination, rescan the workspace, validate scope/secrets/effects, freeze a candidate, and independently run declared build/test/verify gates before review/promotion. Timeout, disconnection, failed termination, failed session resume, incomplete logs/costs or unknown external effects are uncertain and remain paused; switching sessions for a blind retry is forbidden. AgentRunner request, capability, enforcement, event and completion bodies use RFC 8785 canonical JSON with the respective `proofrail:agent-runner-request:1\n`, `proofrail:agent-runner-capability:1\n`, `proofrail:agent-runner-enforcement:1\n`, `proofrail:agent-runner-event:1\n` and `proofrail:agent-runner-completion:1\n` domain separators to compute `recordHash`. Sets are lexically sorted before hashing and reject duplicates. requestId, probeId, enforcementId, eventId and completionId are their primary deduplication keys: the same key and digest is an idempotent replay, while the same key with a different digest conflicts. Events additionally reject a session/sequence bound to another record or one request drifting across sessions; completions reject different terminal records for one request. A capability report is independently compatible only when every fixed matrix entry has runtime evidence and is verified; bound enforcement may provide combined admission under the preceding rule without rewriting capability facts. A completed receipt proves only external execution, process stop and evidence capture, never task PASS.

When human input is required, the CLI Agent adapter may report only structured `operator-action-required`. ProofRail stops or pauses the managed agent at an atomic boundary, persists the request and moves the task/step to `WAITING_FOR_OPERATOR`. ProofRail's CLI/TUI owns notification, response and return of control. Operator input becomes an applicable operator-interaction, review, authorization or handoff record after validation; terminal/chat prose cannot directly change state or authority. Before resuming the same Agent session, ProofRail revalidates attempt, workspace, session, context digest, lease and authorization. T025 freezes the general-clarification `operator-interaction` Schema, positive/negative fixtures, RFC 8785 digest, append-only JSONL replay, and Engine open/resume control API. The control API returns a resume command only when request/response, current binding, and waiting state agree; partial event writes keep the chain paused and converge idempotently by evidence hash. T027 wires this boundary to a real Agent session.

SessionBridge v0.1.1 `silent`, `visible` and `auto` modes are not formal AgentRunner implementations. `silent` may provide tool-free classification, summary, planning, structured analysis or managed-change-set suggestions. `visible` may support observation, diagnostics, notification, human takeover and the supervised black-box candidate flow below. Because `auto` makes the effective path and capability uncertain, it is excluded from the formal ProofRail flow. Any SessionBridge use still embeds the v1 file-IPC client, sets `legacy=false`, binds requestId receipts and owns durable idempotency. Continuing auxiliary conversations may use a stable `conversationId`, but history remains context only. A visible delivery cannot imply Agent completion, and CLI Agent failure must never automatically switch modes and continue writing.

`supervised-black-box` is a formally supported but reduced-assurance candidate-production/intake mode, not an AgentRunner execution substitute. ProofRail first materializes a disposable isolated-workspace and freezes authorized targets and the verification plan. The operator then supervises an external agent in the `visible` AI UI; that agent and its tool loop remain unobservable, untrusted external change producers. After the operator explicitly ends the work window and returns the workspace, ProofRail must prove known managed processes stopped, rebuild the complete manifest/diff, check scope, secrets, file types/sizes and observable effects, then independently run the frozen build/test/verify gates and review/promotion. Missing evidence pauses or rejects the flow. Chat text, visible receipts and verbal confirmation never prove execution completion, tool history or safety.

Before entry, CLI/TUI must display and persist explicit risk acknowledgment: ProofRail cannot prove which tools/commands ran, which resources outside OS isolation were read, network requests, actual tokens/cost/retries, every background process, session continuity, or external effects such as commit/push/publication. The product warrants only its own isolation, post-return artifact scan, independent gates, evidence and acceptance transaction. The user chooses a trusted host, supervises execution, constrains credentials/network and resolves unknown external effects. This responsibility notice cannot waive source/store isolation, secret scanning, independent gates or review, and unknowns cannot be shown as safe. Launch is forbidden without acknowledgment bound to task/attempt/workspace/risk-version; automatic fallback into this mode is forbidden.

| External result/capability | ProofRail action |
|---|---|
| CLI Agent completed/exit zero | Record external termination only; require stop proof, manifest/diff, scope, secret and independent gate/review checks before progress |
| CLI Agent timeout/disconnection/resume failure/failed stop | Record uncertain and pause for reconciliation; do not change session or restart to guess the result |
| CLI Agent requests more tools, network, targets or confirmation | Enter structured authorization/operator flow; do not execute without a valid record |
| SessionBridge silent status=ok | Validate the declared auxiliary-output Schema; it may produce analysis or a managed-change-set suggestion, never direct PASS |
| SessionBridge visible delivered | UI delivery/supervised-observation fact only; black-box mode also requires explicit workspace return and candidate reconstruction from disk, and still proves neither model, tool nor task completion |
| SessionBridge auto | Reject in formal flow; automatic fallback cannot change execution capability or evidence semantics |

Neither SessionBridge's success cache nor a CLI Agent's private session state is durable exactly-once. ProofRail persists dispatch, process, result and application-idempotency records itself. Context history, agent transcripts and UI content cannot replace the run manifest, events, receipts, workspace manifests or gate evidence. ProofRail owns secret scanning/redaction; do not assume an external Agent or SessionBridge implements it.

File queue and IPC share ProofRail request/claim/result envelopes; CLI is a consumer form, not a third business protocol. Each queue file contains exactly one canonical JSONL record. An atomic request move into inflight acquires a claim; immutable generation files and `(claimId,generation)` fencing constrain renewal, takeover and result publication. Never append concurrently or overwrite formal files. Map SessionBridge transport IDs only in dispatch receipts; its single-slot files are not durable ProofRail facts. Takeover receipts bind old/new fencing, stopped-writer proof and authorization evidence; only a granted receipt may be referenced by a new claim. T004 still probes platform atomic/crash boundaries. Use a deterministic fixture consumer for no-IDE/no-AI testing, without paid calls.

`signature-receipt.schema.json` provides detached Ed25519 identity binding for cross-host adapter records, bootstrap
artifacts and release checksum manifests. The signature covers a canonical statement containing purpose, subject kind/hash,
signer, key ID, public-key fingerprint and trust policy; receiptHash then covers the statement plus signature bytes.
cross-host/bootstrap/release are strictly paired with their subject and actor types, and a 64-byte signature accepts only
unpadded base64url. Schema cannot prove key trust, non-revocation at issuance or cryptographic validity; independent crypto
validation, trust policy and verification evidence own those checks. This capability is optional in S1 and is not a mandatory T017/S1-exit release gate.

The release branch of `lifecycle-record.schema.json` always requires binary, SBOM, license-manifest and checksum-manifest hashes; `signatureReceiptHash` may be null. Without a signature, `revocationFreshness` must be `unknown` with no revocation-verification evidence. When a signature receipt is present, the strict contract above still applies. SHA256SUMS proves content integrity, not publisher identity.

## 8. Errors and Acceptance

Persistent errors distinguish schema/static/policy/capability/transport/execution/integrity/storage/review/operator/internal
causes. `error.schema.json` freezes 37 initial codes and binds each group to category, retryMode and suggestedAction;
unknown codes fail closed by version. `error-set.schema.json` binds a non-empty error index and primaryErrorId. Primary
selection uses earliest occurredAt, then fixed category order, code and errorId, so later cleanup failures cannot replace the
original failure. CLI exits 10–20 follow the primary category; usage is 2, never SessionBridge or hook/OS exits. The checker
owns reference consistency and ordering. Errors never contain secrets.

`ticket-ledger.schema.json` records one append-only ledger per `(runId,fingerprint)`. The fingerprint covers only error
code, stable subject and failurePoint, excluding messages, time, attempts and evidence hashes. Failures accumulate across
attempts and an errorHash cannot be counted twice. Below reviewThreshold the state is pending-review; at or above it an
effective override is required for a bounded override-window; hardBlockThreshold is irreversible within the run. Overrides
cannot raise the hard block and resolution does not reset budget. Schema checks the three entry shapes; the checker owns
sequence, threshold relationships, counts, authorization windows and ledger uniqueness.

`repair-transaction.schema.json` fixes Prepare/Inspect/Validate/Promote as a one-to-four-entry stage chain linked by
predecessor hashes. A stage may follow only a successful predecessor; failed or uncertain ends the chain, and completed
requires all four stages. Prepare binds the parent snapshot, isolated candidate, targets, stopped writers and lease; Inspect
binds diff/scope/ownership; Validate binds the frozen plan and hook results; Promote rechecks stop, lease, ledger and Validate
hashes and installs only the next attempt workspace. It creates no accepted snapshot or PASSED result and cannot replace
independent review and snapshot promotion.

Acceptance includes structural, semantic, runtime, crash, compatibility and security checks, not just parseable JSON. See [test strategy](TEST_STRATEGY_EN.md) and [development plan](DEV_PLAN_EN.md). T003 executed the independent structural/static-semantic validator; runtime, crash and platform capabilities remain T004/S1 gates.

## 9. Product and Lifecycle Contracts (RFC Section 19)

T002 froze these semantics and wire fields, and T003 added positive/negative fixtures and independent checks. Reject writes of unknown types.

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

The corresponding structures are `plan-preview.schema.json`, `export-record.schema.json`,
`authorization-record.schema.json`, `effect-record.schema.json`, `cost-ledger.schema.json` and
`lifecycle-record.schema.json`. `verification-report.schema.json` records cross-record decisions in a fixed order spanning
event chains, reference closure, identity/attempt/snapshot, evidence/review/promotion, config/registries,
change-set/repair/ticket, queue fencing, signature trust and lifecycle authorization. Failed or incomplete results fail
closed; a report never changes state or authority.