# P43-T027B-STATE-SEAL · ① 架构方案（原样）

载体：`deep-reasoner`；模型：`DeepSeek V4 Pro (deepseek)`；只读。判定：`DESIGN: DONE`。

## 一 结论与范围
- 本片只动 `internal/chain/**`+`internal/tickets/**`：新增 `internal/chain/seal.go`；`events.go` 增 `SealSink`；`models.go` 增 `TaskAttemptResolver`/`WorkspaceReopener` 接口；`engine.go`/`postflight.go` 11 处写死 attempt 改造。
- 封存载体=链内 seal 文件（`FileEventStore` 路径旁 `seals/` 子目录），因 `validTransition` 无自环（event.go:212-228）不能做同态 seal 事件；intent 事件引用 sealHash。
- attempt 真源=`TaskAttemptResolver`（默认投影派生=1，与现状全等）；3 处 `!=1` 改为「等于解析值」，12 处 `<1` 全保留。
- 顺序不变式：封存→intent（task-accepted）→结算→C→closure；崩溃窗口靠 seal 文件幂等重放，不重调 Reviewer/Publisher。
- 协议先行条目 2 条（F①F⑤），供 ⑧ 落笔。

## 二 A 封存屏障
A1 对象与键：
- seal 记录（`chain/seal.go`）：`SealRecord{SchemaVersion,RunID,TaskID,Attempt,ParentHash,CandidateHash,EvidenceRootHash,ReviewReceiptHash,PromotionReceiptHash,ReviewDecision,PromotionDecision,Evidence[]}`，canonical 编码域 `proofrail:task-seal:1\n`（仿 review.go:12 域分隔），键=`<taskID>.<attempt>`。
- 落盘：`FileEventStore` 新增 `SealTaskEvidence(ctx,runID,taskID,attempt,record)→hash`，文件 `<事件日志目录>/seals/<runID>/<taskID>.<attempt>.seal.json`；临时文件+`Sync`+`Close`+`Link` no-replace（复刻 snapshot store.go:101 模式），已存在同哈希对账复用（store.go:107 语义）。
- 引擎侧经新接口 `SealSink`（events.go 内定义）类型断言 `options.Events`；生产路径全是 `FileEventStore`（console runtime.go:34/71、runtime_agent.go:269、runtime_supervised.go:448/685）。

A2 屏障判定：
- 「已封存」= seal 文件存在且 DecodeStrictJSON（review.go:87 同款重算）后重算哈希==文件名记录哈希，且其字段哈希与本次候选/收据哈希逐项相等。
- 防半封存：文件与事件双写双查——seal 文件先 fsync，intent 事件（events.go:34 `FileEventStore.Append` 内含 events.go:49 `file.Sync()`）后写且 InputEvidence 末尾附 sealHash；文件在而事件缺→忽略重封（幂等）；事件在而文件缺/哈希不等→缺口，禁 PASSED。
- 恢复验证：`Engine.New`（engine.go:69 Rebuild 后）与 `Recover`（recover.go:11 重放后）对每个 PASSED 任务执行 `verifySealedAccepted`：加载 seal、解码、核验哈希与 InputEvidence 闭包⊇；失败返回 `ErrRecoveryUncertain`。

A3 插入点：
- task：engine.go:280（`uniqueHashes(inputs)` 之后）与 281（`transition(…"PASSED",inputs,"task-accepted")`）之间插 `engine.sealTaskEvidence`：① 闭包哈希全 `ValidHash` 校验；② `SealSink.SealTaskEvidence`（fsync）；③ sealHash 追加进 inputs 再走 281。崩溃后重启：task 停在 REVIEW_PENDING，重入先 `LoadSeal`——存在则**跳过 Reviewer/Publisher 复调**（复用 sealed 两决策重建收据哈希）直接 281；不存在才走全量路径。窗口内任何崩溃都先于 281 事件写入，不越障。
- step：engine.go:377（Execute 成功）与 378（`step-passed`）之间只插闭包合法性校验（`result.Evidence` 全 ValidHash），**不加 seal 文件**——step 证据已由执行器先落盘（fact#4：gates/result.go:81 fsync、agent-runner 三记录先行），step-passed 事件的 InputEvidence 即闭包本身。

A4 缺口处置：
- 哈希非法→reason `evidence-closure-incomplete`；写入/同步失败→`evidence-seal-failed`；均 `markTaskForRepair`→REPAIR_PENDING+chain PAUSED（engine.go:265-277），**不写 PASSED**——PASSED 唯一入口 281 被屏障独占，completed 自然被禁。
- 恢复期缺口（seal 文件缺失/篡改）：`ErrRecoveryUncertain`（复用 engine.go:340-343 的 fail-closed 语义），禁止补造、覆盖赢家、继续结算。

## 三 B 不可变摘要 K
- 构成：`K=Digest("proofrail:context-key:1\n", runID‖taskID‖stepID‖decimal(attempt)‖parentHash‖workspaceRoot‖executionTarget‖requestId‖requestHash)`；`‖`=`\x00` 域分隔（仿 state.go:127 taskKey），缺失字段空串占位；输出=64 hex（evidence.Digest）。
- 不可变：**无链内写点**——K 是纯派生函数，定稿于 `step-started`（engine.go:328）时其全部输入已在事件日志固化（run/task/step/attempt/parent/workspace 绑定）；唯一写者 engine（state.go:64 单 append 路径）。重复采集=从同一不可变输入重算：相等→复用；不等→`ErrRecoveryUncertain`（对齐 engine.go:340-343 不执行两次语义），永不改写。
- 与 ContextHash 对齐：operator_interaction.go:320 绑定比较追加 `current.K != requestK` 拒绝；operator_interaction.go:542 响应校验追加 `responseK == requestK`（K 由请求/响应记录字段+解析 attempt 重算，不信任适配器传值）；K 含 attempt ⇒ attempt=2 记录天然无法绑定 attempt=1 上下文。

## 四 C attempt 真源与改造点
真源：`TaskAttemptResolver`（models.go 新接口，`ResolveTaskAttempt(ctx,runID,taskID)(int,error)`）；默认实现=事件投影派生（取任务实体最大 attempt，state.go:127 键含 attempt），现状全 1⇒行为全等；tickets 授权（ledger.go:121 `AppendOverride` 载 attempt）供装配层接线（F④）。`runTask` 入口解析一次向下传递。

| file:line | 现状 | 改法 | 保留/删除理由 |
|---|---|---|---|
| engine.go:140 | `taskKey(task.ID,1)` | 先解析 attempt 再 `taskKey(task.ID,attempt)` | 跳过判定须按真实身份 |
| engine.go:183 | `Materialize(…,1,parent)` | 传 attempt | 工作区隔离（models.go:186 参数已备） |
| engine.go:211 | `Accept(…,1,…)` | 传 attempt | 候选隔离 |
| engine.go:221/233/250/262 | 四处 `Attempt:1` | 传 attempt | 收据身份真实 |
| engine.go:326 | `StepRequest…Attempt:1` | 传 attempt（runStep 增参） | 步骤身份真实 |
| engine.go:476/479 | taskEntity/stepEntity 写死 1 | 加 attempt 参数，全调用点改造 | 实体键隔离（state.go:127） |
| postflight.go:390/431 | 两处 `Attempt:1` | 传 attempt | 同族写死点 |
| agent_runner_terminal_route.go:33 | `terminal.Attempt != 1` | `!=`解析值 | 保留比较、换真源 |
| operator_interaction.go:334/389 | `request.Attempt != 1` | 同上 | 保留比较、换真源 |
| chain 8 处 `<1` | 范围断言 | 全保留 | 与真源无冲突 |
| tickets `<1`（budget.go:52、ledger.go:95/122/152） | 范围断言 | 全保留 | `ReserveAttempt`（budget.go:45）已把 attempt 当入参 |
| 无 | — | 新增 `WorkspaceReopener.Reopen(ctx,run,task,attempt,parent)`（models.go 独立接口） | P38 蓝图置 WorkspacePort，本片禁改 console 实现者故拆接口 |

## 五 D 顺序不变式
| # | 步骤 | 落盘对象 | 失败后状态边 | 可重入语义 |
|---|---|---|---|---|
| 1 封存 | seal 文件（SealSink，fsync+Link no-replace） | `evidence-seal-failed`/`evidence-closure-incomplete`→REPAIR_PENDING（event.go:223 REVIEW_PENDING→REPAIR_PENDING 合法） | 存在即解码校验，跳过 Reviewer/Publisher 复调 |
| 2 intent | `task-accepted` 事件（engine.go:281，REVIEW_PENDING→PASSED 合法 event.go:225），InputEvidence=闭包∪sealHash | append 失败→REPAIR_PENDING；崩溃→重放走 A3 | 已 PASSED 跳过（engine.go:140） |
| 3 结算 | tickets：`ReserveAttempt`（budget.go:45）先于派发消费一次；授权 `AppendOverride`（ledger.go:121/122） | `ErrBudgetAttemptExhausted`（budget.go:56） | consumedByAttempt 去重（budget.go:55） |
| 4 C | `chain-completed`（engine.go:170） | 任一任务非 PASSED 不写 | 幂等 |
| 5 closure | `AcceptedParent←PASSED.InputEvidence[0]`（state.go:42/120） | 无效→`accepted-parent-missing` | 冻结父快照复用 |

## 六 E 红测清单（③ 先行）
| 测试名 | 断言要点 | 应失败形态 |
|---|---|---|
| `TestTaskSealBarrierPreventsAccepted` | SealSink 注入 Sync 失败：无 `task-accepted`、REPAIR_PENDING | 现红 |
| `TestClosureIncompleteCannotSeal` | EvidenceRootHash 非法：`evidence-closure-incomplete` | 现红 |
| `TestCrashAfterSealBeforeIntentReplaysDecisions` | seal 后新 Engine 重建：Reviewer/Publisher 调用 0 | 现红 |
| `TestPassedTaskSealMustVerifyOnRecovery` | 删/改 seal 文件后→`ErrRecoveryUncertain` | 现红 |
| `TestReopenUsesResolvedAttempt` | resolver=2：全链收 2、实体键 attempt=2 | 现红 |
| `TestAttemptIsolationAcrossWorkspaces` | attempt=1 与 attempt=2 键并存互不可见 | 现红 |
| `TestTerminalAttemptMismatchRefused` | terminal.Attempt=2 vs 解析 1：拒绝、零写 | 现红 |
| `TestOperatorInteractionAttemptAndKBinding` | 334 attempt 不齐拒绝；320/542 K 不齐拒绝 | 现红 |
| `TestContextKeyStableAcrossReentry` | 重建重算 K 相等复用；定义漂移→不改写 | 现红 |

## 七 F 风险与未确证项
- ① 旧事件日志无 sealHash（既有 run 目录）→ 恢复校验对新旧日志适用性未定；需协议先行：seal 强制生效边界，供 ⑧ 落笔。
- ② 新增 reason code `evidence-seal-failed`/`evidence-closure-incomplete` 及 `task-accepted` 证据布局追加 sealHash：CONTRACTS 若枚举 reason/证据语义需协议先行条目。
- ③ `<1` 计数：fact 称 13（tickets 5），实测 12（chain 8+tickets 4）；第 5 处未确证。
- ④ console 装配不接线 `TaskAttemptResolver` ⇒ 生产 attempt 恒 1，N3 仅链内闭合；装配片接线。
- ⑤ `WorkspaceReopener` 拆接口与 P38「WorkspacePort.Reopen」有差。
- ⑥ operator_interaction.go:320/542 为批量辅助锚点，未逐行编号（±3 行）。
- ⑦ 链 import tickets 无环（未确证导入后 vet 无告警）。
- ⑧ tickets 结算「真实 attempt 写入账本」的上游在 adapters/console，本片不可见，未确证。

DESIGN: DONE
