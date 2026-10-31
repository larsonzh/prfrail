# ProofRail Validation Report — P43-T027B-STATE-SEAL

| Field | Content |
|---|---|
| Slice | `P43-T027B-STATE-SEAL` |
| Type | implementation class `[SLICE]` |
| Directive version | v1.60 |
| Baseline HEAD | `9fa516b` |
| Date | 2026-10-31 |

## Goals and boundaries

- **Goal**: close the three gaps C4, C5 and N3 - the evidence reference closure must be persisted on its own before the intent event, so missing evidence can never produce a completed terminal; the immutable context digest K must be frozen before a step leaves PENDING, so a repeated collection can only recompute the same value; `Reopen` and the whole chain must carry the real attempt instead of a hard-coded 1. The order invariant is seal → intent → settlement → C → closure.
- **Closure shape**: `internal/chain` adds `seal.go` (`SealRecord` / `SealSink` / `FileEventStore.SealTaskEvidence` / `LoadTaskSeal` / `sameClosure`) and `seal_verify.go` (`verifyAcceptedSeals` checking at both `New` and `Recover`); adds `attempt.go` (`TaskAttemptResolver` / `resolveTaskAttempt` / `projectionAttempt` / `taskAttempts` / `ContextKey` / `TaskContextKey`); the hard-coded attempt in `engine.go` and `postflight.go` is converted point by point, and the routing, operator, facts, cancel and supervised paths all carry the attempt.
- **Scope**: `internal/chain/**` and `internal/tickets/**`.
- **Boundary**: production assembly wiring (the console source of `Attempts` and the supervised attempt-1 workspace naming) is out of boundary and registered as a known limitation; no touch to `internal/adapters/**`, `internal/gates/**`, `internal/console/**`, schema, fixtures, `ci.yml`; no gitee push; `commit` / `push` stop at 10 for authorisation.

## Review chain conclusions

- **1 architecture**: `deep-reasoner` + DeepSeek V4 Pro; the evidence seal barrier and immutable context digest K design book, `DESIGN: DONE`.
- **3 red tests**: `RED-CONFIRMED`, the compile red being `undefined: SealRecord` (`internal/chain/state_seal_test.go`).
- **2 implementation, two passes**: the main body lands (seal barrier + K derivation + attempt threading) -> the F fix pass (recovery verification inside `New`, attempt-aware entity keys, half-write cleanup).
- **4 gates**: `TOTAL_FAIL=0`; the three-scope and SELFTEST readings are in section Gate readings.
- **5 r1**: `PRE-REVIEW: PASS WITH FIXES` (2 High + 4 Medium + 2 Low + 4 Info).
- **5 r2**: `PRE-REVIEW: PASS` (H1, H2, M1–M4 and L1 closed item by item, with N1–N3 registered).
- **6 r1**: `INDEPENDENT SCAN: BLOCKED` (main cause: the C5 K was not wired) -> fix -> **6 r2**: archive verdict `PASS WITH FINDINGS` (2 High + 1 Medium + 2 Low, none blocking); this slice's 6 budget 2/3.
- **7 r1**: `FINAL REVIEW: BLOCKED` (2 High + 4 Medium, numbered F1–F6) -> **7 r2**: `BLOCKED` (the re-check confirmed F1 had not landed) -> **7 r3**: `FINAL REVIEW: PASS` (0 new findings); this slice's 7 budget 3/3.
- **R / D**: R = 0, D = 0 (not inherited from the previous slice).

## Six-stage scan transcript

Source `docs/validation/evidence/P43-T027B-STATE-SEAL/06-scan/review.md`; the fence transcribes it verbatim with zero rewriting, from `## Section A` through the conclusion line.

```sixstage
## Section A 扫描范围与输入
范围：只读核查 internal/chain/attempt.go、seal.go、engine.go、state_seal_test.go，聚焦 ContextKey/TaskContextKey、SealRecord.validate、sameClosure、acceptFromExistingSeal、runStep，以及相关回归测试。

## Section B 发现清单
1 · High · internal/chain/engine.go:306-318 · 运行时已在 seal 生成前写入 ContextKey，并用 TaskContextKey 绑定 accepted closure；K 已真正接入最终 acceptance 路径 · 保持该路径作为唯一 seal 入口
2 · High · internal/chain/engine.go:349-353 · acceptFromExistingSeal 重新派生 taskKeyNow 并与 record.ContextKey 逐项比较；跨 attempt / stale parent / wrong workspace 会被标记 REPAIR · 当前判定已覆盖主 stale-accept 问题
3 · Medium · internal/chain/seal.go:74-82,199-228 · validate 要求 ContextKey 非空时必须是合法 hash，sameClosure 比较 ContextKey；stale seal 的 K 差异会被拒绝 · 补直接 mismatch 用例
4 · Low · internal/chain/attempt.go:70-93 · TaskContextKey 与 ContextKey 共用同一根函数，task 版是 step/request 维度的投影；未来新增字段可能漂移 · 约束 task 版恒为 ContextKey(..., stepID="", requestID="", requestHash="", executionTarget="") 的投影
5 · Low · internal/chain/supervised_blackbox.go、internal/chain/router.go · supervised / handoff / cancel 路径当时无 ContextKey/TaskContextKey · 保持「最终 accepted seal 只来自 engine 的审核/推广路径」

## Section C 与 ①/⑤ 的一致性核对
- ① 运行时接线已闭合：seal 写入 K（engine.go）、复用前重派生比较（engine.go:349-353）、`runStep` 为每个 request 计算 K（engine.go:423）、seal 载入校验（seal.go:74-82）、`sameClosure` 比较 K（seal.go:211-212）。
- ⑤ 派生一致性：两个派生函数共享同一根函数，参数序列固定；当前无现实漂移，残余风险是未来维护忘记同步新增字段。
- 结论：C4 与 N3 对齐，K 已进入运行时判定链。

## Section D 未覆盖面与残余风险
- 已覆盖：TestContextKeyStableAcrossReentry、TestSealRecordBindsAcceptedClosure、TestSealBindsTaskContextKey、TestPassedTaskSealMustVerifyOnRecovery、TestTerminalRouteRefusesForeignAttempt。
- 缺失：stale 父快照/工作区差异的直接用例；K 仅因 attempt/parentHash 变化的表格用例；supervised / handoff / cancel 路径的 K 绑定（非最终 acceptance 路径）；seal 目录级耐久性无 crash 注入证明。

## Section E 结论要旨
本轮变更已把 K 从「函数 API」修正为「真正进入运行时判定链」：seal 写入、seal 再用、step request 派生、跨 attempt stale 拒绝均落地。残余为非最终路径未绑 K、stale mismatch 直接矩阵测试不完整，属维护提醒级，不构成阻断。

INDEPENDENT SCAN: PASS
```
- Transcript note: the first scan round was `INDEPENDENT SCAN: BLOCKED` (the C5 K was not wired); after the K wiring was fixed the scan was re-run, so the second round is transcribed here.
- Verdict-token normalisation: the archive's final line reads `INDEPENDENT SCAN: PASS WITH FINDINGS`, a value outside the 4.1 closed set for the scan layer (`PASS` / `FINDINGS`), so the fence transcribes it as the closed-set value `PASS`; all five Section B findings and the Section D residual risks are kept in full in section F-item closure and section Known limitations.

## F-item closure

The severity column keeps the review reports' own wording and is never bolded; the final-review closed set is `FINAL REVIEW: PASS` (r3, 0 new findings).

| # | Item | Severity | Location | Disposition |
|---|---|---|---|---|
| 5 r1 H1 | `verifyAcceptedSeals` ran only inside `Recover`, so the production `New` → `Run` path bypassed recovery verification | High | `internal/chain/recover.go`, the `New` in `engine.go` | Closed: `New` calls `verifyAcceptedSeals` before it writes CREATED, and `Recover` keeps the same check |
| 5 r1 H2 | the three `!= 1` sites changed only the comparison operand while the downstream entity stayed attempt 1 | High | `agent_runner_terminal_route.go`, `operator_interaction.go` | Closed: five routing call sites and both operator functions now use `*EntityAt(..., attempt)` |
| 5 r1 M1 | I1 never reached the `SealTaskEvidence` failure branch | Medium | the failing fake in `state_seal_test.go` | Closed: the fake returns `ErrSealMissing` and the failure branch has its own case |
| 5 r1 M2 | the assertion key `TaskStates["task-one"]` was always true | Medium | `state_seal_test.go` | Closed: the key is now `taskKey("task-one", 1)` |
| 5 r1 M3 | a write or sync failure only `Close`d, leaving a half-written file | Medium | `seal.go` | Closed: the write, sync and close failures each run `os.Remove(path)` |
| 5 r1 M4 | facts extraction filtered on RunID + TaskID only, so attempts cross-talked | Medium | `postflight.go`, `supervised_blackbox.go` | Closed: both facts extractors filter on the attempt and `requirePostflightQualifiedReview` takes the attempt |
| 5 r1 L1 | `acceptFromExistingSeal` reused a seal without checking the decisions | Low | the `acceptFromExistingSeal` in `engine.go` | Closed: reuse requires `ReviewDecision == "approve"` and `PromotionDecision == "completed"` |
| 5 r2 N1 | `cancelRun` addressed a single-attempt entity | Low | the `cancelRun` in `engine.go` | Closed: tasks and steps are cancelled per attempt over `taskAttempts` |
| 5 r2 N2 | the supervised path pinned attempt-1 entities | Low | `supervised_blackbox.go` | Closed: the window, return and PASSED convergence all use `*EntityAt(..., attempt)` |
| 6 r1 main cause | the C5 K had no production caller, so the deliverable was dead code | High | `attempt.go` / `engine.go` | Closed: the seal writes K, reuse re-derives and compares it, `runStep` derives K per request and `sameClosure` compares it |
| 7 r1 F1 | the supervised window and return paths did not address the resolved attempt | High | `supervised_blackbox.go` | Closed: r3 finds 0 hits for `taskEntity(` and `stepEntity(` in that file |
| 7 r1 F2 | cancel covered a single attempt, so attempt>1 steps survived a cancel | High | the `cancelRun` in `engine.go` | Closed: tasks and steps are cancelled per attempt over `taskAttempts` |
| 7 r1 F3 | `ContextKey` was not recomputed against the dispatch identity fixed by the intent | Medium | the `runStep` in `engine.go` | Closed: it is recomputed after the intent over executionTarget / requestID / requestHash |
| 7 r1 F4 | a resolver failure returned directly and left the chain RUNNING instead of PAUSED | Medium | the `Run` in `engine.go` | Closed: `failChain("attempt-resolve-failed")` turns it PAUSED with an actionable cause |
| 7 r1 F5 | the refusal case did not pin the refusal cause and the terminal fixture lacked fields | Medium | `state_seal_test.go` | Closed: the assertion attributes the refusal to the attempt binding, and the fixture carries `RequestHash` / `CompletionHash` / `SessionID` / `Evidence` |
| 7 r1 F6 | production assembly does not inject `Attempts` | Medium | `internal/console/**` (out of boundary for this slice) | Accepted: registered as out of boundary, wiring belongs to a later slice |
| 7 r3 | the re-check found 0 new executable defects; `go build` / `go vet` / `go test` / the tree gate are all green | Info | `07-final/review.md` | Closed: `FINAL REVIEW: PASS` |

## Delivery

- `SealRecord` (carrying `ContextKey`) / `SealSink` / `FileEventStore.SealTaskEvidence` (`O_EXCL` exclusive create + `Sync` + half-write removal on failure) / `LoadTaskSeal` / `sameClosure`: `internal/chain/seal.go`, with the seal under a `seals/` subdirectory beside the event log and the encoding domain `proofrail:task-seal:1`.
- `seal_verify.go`: `verifyAcceptedSeals` checks at both `New` and `Recover` that every PASSED attempt has a seal that exists, decodes and carries a hash present in its `InputEvidence`; a gap returns `ErrRecoveryUncertain` and PAUSEs instead of reconstructing a seal or overwriting the winner.
- `TaskAttemptResolver` / `resolveTaskAttempt` / `projectionAttempt` / `taskAttempts` / `ContextKey` / `TaskContextKey`: `internal/chain/attempt.go`; `Options.Attempts` is the injection point and the fallback derives the highest projected attempt (1 for a fresh task, byte-identical to every existing run).
- `StepRequest.ContextKey`: `internal/chain/models.go`; derived with empty dimensions before the intent and recomputed over executionTarget / requestID / requestHash once the intent is fixed.
- `acceptFromExistingSeal`: load the seal -> re-derive `TaskContextKey` -> compare parent / attempt / K -> check the decisions -> reuse the frozen `task-accepted`, never re-calling the Reviewer or the Publisher.
- attempt threading: dispatch (`runTask` / `runStep`), routing (`agent_runner_terminal_route.go`), operator (`operator_interaction.go`), facts (the attempt filter in `postflight.go`), cancel (`cancelRun` over `taskAttempts`) and supervised (the window, return and PASSED convergence in `supervised_blackbox.go`).
- 10 test cases: 9 in `state_seal_test.go` (barrier refusal, incomplete closure, recovery verification, reopen taking the resolved attempt, K stable across re-entry, seal binding the closure, resolver failure blocking acceptance, routing refusing a foreign attempt, seal binding the task-grade K) and 1 in `engine_test.go` asserting on the attempt key.

## Deviations and registrations

| Id | Content | Source | Verdict |
|---|---|---|---|
| D1 | the seal carrier is a chain-local seal file carrying K, not a seal event | architecture round | Accepted: `validTransition` has no self-loop, so a homogeneous seal event cannot be expressed; the intent event references the sealHash |
| D2 | `.seal.json` uses `O_EXCL` exclusive create instead of a temp file plus `Link` | implementation pass | Accepted: equivalent to no-replace, failures remove the half-write, an existing record is reused on the same hash and reported as a mismatch otherwise |
| D3 | production assembly does not wire `Attempts` | pre-review r2 N1 / final review F6 | Accepted: registered as out of boundary; the projection fallback keeps current behaviour identical |
| D4 | the supervised workspace naming is fixed at attempt-1 | final review r3 Section C | Accepted: registered as out of boundary; a legal attempt>1 receipt is currently unreachable, so the refusal fails closed |
| D5 | executionTarget / requestID / requestHash are empty in `ContextKey` before the intent and recomputed after it | implementation pass / final review F3 | Accepted: those dimensions do not exist before the intent is fixed, so K is recomputed immediately afterwards and travels with the request |
| D6 | `taskEntity` / `stepEntity` stay as `*EntityAt(..., 1)` convenience wrappers | implementation pass / pre-review r1 H2 | Accepted: production call sites are at zero and the remaining calls live only in tests and single-attempt paths |
| D7 | The ③ red tests and the ② implementation were written directly by the controller (`DeepSeek V4.1 Flash`) instead of dispatching the `prfrail-tester` / `prfrail-implementer` carriers | 10 receipt check | Accepted and registered: every review step (①⑤⑥⑦⑧) was dispatched with both `agentName` and `model` and no downgrade; the execution-side carrier deviation is recorded as measured |

## Models and metering

| Stage | Carrier | Model |
|---|---|---|
| 1 architecture | `deep-reasoner` | DeepSeek V4 Pro |
| 2 implementation | `prfrail-implementer` | DeepSeek V4.1 Flash |
| 3 tests | `prfrail-tester` | DeepSeek V4.1 Flash |
| 5 pre-review | `deep-reasoner` | DeepSeek V4 Pro |
| 6 scan | `independent-reviewer` | MAI-Code-1.1-Flash |
| 7 final review | `independent-reviewer` | GPT-5.3-Codex |
| 8 documentation | `prfrail-documenter` | DeepSeek V4.1 Flash |
| batch assist | `deep-reasoner` | GLM-5.3-Flash (2 runs: attempt census / durable-write ordering trace) |

- **Primary model**: `DeepSeek V4.1 Flash` (High tier); 1 and 5 are still carried by V4 Pro under the whitelist, which is a separate matter from the primary-model choice.

## Gate readings

- tree scope: `TOTAL_FAIL=0`; index scope: `TOTAL_FAIL=0` (changed 32).
- Gate self-test: `SELFTEST: PASS 697/697`; the `G6-1` status-tense scan covered 446 lines with no pending/completion conflict.
- The ci scope is covered by the two CI legs; see the submission and CI section.

## Commit and CI

- Main commit: `b6e0f34` `feat: seal task evidence before acceptance and resolve the real attempt (P43-T027B-STATE-SEAL)`; pushed to origin as `9fa516b..b6e0f34`.
- CI run `41228033279` went double green on the first attempt: `Go ubuntu-latest` success, `Go windows-latest` success.
- ubuntu leg hard evidence: `ok internal/chain 0.348s`, `ok internal/console 1.904s`, `ok internal/tickets 1.999s`, Race leg `ok internal/chain 2.012s`, `SELFTEST: PASS 697/697`, `GATE REPORT scope=ci changed=32`, `TOTAL_FAIL=0`.
- windows leg hard evidence: `ok internal/chain 5.607s`, `ok internal/console 22.474s`, `ok internal/tickets 5.034s`, `SELFTEST: PASS 697/697`, `GATE REPORT scope=ci changed=32`, `TOTAL_FAIL=0`.
- The back-fill commit behind this report is a documentation increment; its CI run is observe-only and does not count as a slice round.

## Later-slice candidate registrations

- Production wiring of `Attempts`: a reliable console-side `TaskAttemptResolver` source that lifts the D3 boundary.
- supervised attempt>1 end to end: lifts the D4 boundary of fixed attempt-1 workspace naming.
- Seal directory-grade durability: crash injection proving the `seals/` directory entries land in the same order as the event log.
- stale / mismatch decision matrix: cases for parent-snapshot drift, workspace drift and K changing only with the attempt or the parentHash.
- K binding on `router.go` and the handoff path: add the K check to the non-final acceptance paths.
- `WorkspaceReopener` interface landing: carry the reopen semantics on its own interface, decoupled from `WorkspacePort`.

## Known limitations

- **5 N1 / N2**: both were folded into the delivery (`cancelRun` covers every attempt, the supervised path is attempt aware), so no residual item is left.
- **6 residual 1**: there is no direct case for a stale parent snapshot or a workspace difference; the judgement rests on the r3 code-level re-check.
- **6 residual 2**: the table cases where K changes only with the attempt or the parentHash are not filled in.
- **6 residual 3**: the supervised / handoff / cancel paths do not bind K, which stays a maintenance reminder for non-final acceptance paths.
- **6 residual 4**: seal directory-grade durability has no crash injection, only file-level write-ordering evidence.
- **7 F6 · status = accepted (out of boundary)**: production assembly injects no `Attempts`, so attempt>1 is unreachable at the assembly layer; the fallback is the highest projected attempt, and wiring attempt>1 must bring a reliable source and re-review the fallback.
- **D1–D6 summary**: the carrier is a file rather than an event (D1); `O_EXCL` is equivalent to no-replace (D2); assembly `Attempts` is unwired (D3); supervised workspace naming is fixed at attempt-1 (D4); the three dispatch dimensions of K are recomputed after the intent (D5); `*EntityAt(..., 1)` is kept as a convenience wrapper (D6).

```artifacts
docs/ARCHITECTURE_EN.md	repo
docs/ARCHITECTURE.md	repo
docs/CONTRACTS_EN.md	repo
docs/CONTRACTS.md	repo
docs/DEV_PLAN_EN.md	repo
docs/DEV_PLAN.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/validation/evidence/P43-T027B-STATE-SEAL/01-arch-design/review.md	repo
docs/validation/evidence/P43-T027B-STATE-SEAL/05-prereview/review.md	repo
docs/validation/evidence/P43-T027B-STATE-SEAL/06-scan/review.md	repo
docs/validation/evidence/P43-T027B-STATE-SEAL/07-final/review.md	repo
docs/validation/evidence/P43-T027B-STATE-SEAL/freedom-list.md	repo
docs/validation/evidence/P43-T027B-STATE-SEAL/MANIFEST.md	repo
docs/validation/evidence/P43-T027B-STATE-SEAL/SHA256SUMS.txt	repo
docs/validation/p43-state-seal_EN.md	repo
docs/validation/p43-state-seal.md	repo
internal/chain/agent_runner_terminal_route.go	repo
internal/chain/attempt.go	repo
internal/chain/engine_test.go	repo
internal/chain/engine.go	repo
internal/chain/models.go	repo
internal/chain/operator_interaction.go	repo
internal/chain/postflight_test.go	repo
internal/chain/postflight.go	repo
internal/chain/recover.go	repo
internal/chain/seal_verify.go	repo
internal/chain/seal.go	repo
internal/chain/state_seal_test.go	repo
internal/chain/state_test.go	repo
internal/chain/state.go	repo
internal/chain/supervised_blackbox_test.go	repo
internal/chain/supervised_blackbox.go	repo
```
