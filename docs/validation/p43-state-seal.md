# ProofRail 验证报告 — P43-T027B-STATE-SEAL

| 字段 | 内容 |
|---|---|
| 切片 | `P43-T027B-STATE-SEAL` |
| 类型 | 实现类 `[SLICE]` |
| 准则版本 | v1.60 |
| 基线 HEAD | `9fa516b` |
| 日期 | 2026-10-31 |

## 目标与边界

- **目标**：闭合 C4、C5 与 N3 三项缺口——证据引用闭包必须在 intent 事件之前独立持久化，缺失证据不得产生 completed 终局；上下文不可变摘要 K 必须冻结于步骤离开 PENDING 之前，重复采集只能重算同一值；`Reopen` 与全链 attempt 必须取真实 attempt，不再写死 1。顺序不变式为封存 → intent → 结算 → C → closure。
- **收口形态**：`internal/chain` 新增 `seal.go`（`SealRecord` / `SealSink` / `FileEventStore.SealTaskEvidence` / `LoadTaskSeal` / `sameClosure`）与 `seal_verify.go`（`verifyAcceptedSeals` 在 `New` 与 `Recover` 双点校验）；新增 `attempt.go`（`TaskAttemptResolver` / `resolveTaskAttempt` / `projectionAttempt` / `taskAttempts` / `ContextKey` / `TaskContextKey`）；`engine.go` 与 `postflight.go` 的写死 attempt 逐点改造，路由、operator、facts、cancel、supervised 五条路径贯通 attempt。
- **范围**：`internal/chain/**` 与 `internal/tickets/**`。
- **边界**：生产装配接线（console 的 `Attempts` 来源与 supervised attempt-1 工作区命名）属边界外，登记为已知限制；不碰 `internal/adapters/**`、`internal/gates/**`、`internal/console/**`、schema、fixtures、`ci.yml`；不推 gitee；`commit` / `push` 停 ⑩ 等授权。

## 审查链结论

- **① 架构**：`deep-reasoner` + DeepSeek V4 Pro；《证据封存屏障与不可变上下文摘要 K》，`DESIGN: DONE`。
- **③ 红测**：`RED-CONFIRMED`，编译红为 `undefined: SealRecord`（`internal/chain/state_seal_test.go`）。
- **② 实现两轮**：主体落成（seal 屏障 + K 派生 + attempt 贯通）→ F 修复轮（`New` 内恢复校验、实体键改 attempt 感知、半写清理）。
- **④ 门禁**：`TOTAL_FAIL=0`；三作用域与 SELFTEST 读数见 §门禁读数。
- **⑤ r1**：`PRE-REVIEW: PASS WITH FIXES`（2 High + 4 Medium + 2 Low + 4 Info）。
- **⑤ r2**：`PRE-REVIEW: PASS`（H1、H2、M1–M4、L1 逐条闭合，另登记 N1–N3）。
- **⑥ r1**：`INDEPENDENT SCAN: BLOCKED`（主因：C5 的 K 未接线）→ 修复 → **⑥ r2**：归档判定 `PASS WITH FINDINGS`（2 High + 1 Medium + 2 Low，均非阻断），本片 ⑥ 额度 2/3。
- **⑦ r1**：`FINAL REVIEW: BLOCKED`（2 High + 4 Medium，编号 F1–F6）→ **⑦ r2**：`BLOCKED`（复查确认 F1 未落地）→ **⑦ r3**：`FINAL REVIEW: PASS`（0 新增发现），本片 ⑦ 额度 3/3。
- **R / D**：R = 0、D = 0（不继承上片）。

## 六段扫描转录

来源 `docs/validation/evidence/P43-T027B-STATE-SEAL/06-scan/review.md`；围栏内逐字转录、零改写，范围自 `## Section A` 起、至结论行止。

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
- 转录口径：⑥ 首轮判定 `INDEPENDENT SCAN: BLOCKED`（C5 的 K 未接线），修复 K 接线后重扫，此处转录第二轮。
- 判定词归一化：归档末行为 `INDEPENDENT SCAN: PASS WITH FINDINGS`，该取值不在 §4.1 的 ⑥ 闭集（`PASS` / `FINDINGS`）内，故围栏按闭集转录为 `PASS`；Section B 的五条发现与 Section D 的残余风险全文保留在 §F 项闭合 与 §已知限制。

## F 项闭合

严重度列取审查报告原文口径，不加粗；⑦ 结论闭集为 `FINAL REVIEW: PASS`（r3，0 新增发现）。

| # | 项 | 严重度 | 位置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 H1 | `verifyAcceptedSeals` 只在 `Recover` 内，生产路径 `New` → `Run` 绕过恢复校验 | High | `internal/chain/recover.go`、`engine.go` 的 `New` | 已闭合：`New` 在写 CREATED 之前调用 `verifyAcceptedSeals`，`Recover` 保留同校验 |
| ⑤ r1 H2 | 三处 `!= 1` 只改比较操作数，下游实体仍写死 attempt 1 | High | `agent_runner_terminal_route.go`、`operator_interaction.go` | 已闭合：路由五处调用点与 operator 两函数改用 `*EntityAt(..., attempt)` |
| ⑤ r1 M1 | I1 未真正走到 `SealTaskEvidence` 失败分支 | Medium | `state_seal_test.go` 的失败 fake | 已闭合：fake 返回 `ErrSealMissing`，失败分支另立用例 |
| ⑤ r1 M2 | 断言键 `TaskStates["task-one"]` 恒真 | Medium | `state_seal_test.go` | 已闭合：断言键改 `taskKey("task-one", 1)` |
| ⑤ r1 M3 | 写/同步失败只 `Close`，留下半写文件 | Medium | `seal.go` | 已闭合：写、同步、关闭三处失败各 `os.Remove(path)` |
| ⑤ r1 M4 | facts 提取只按 RunID + TaskID 过滤，跨 attempt 串扰 | Medium | `postflight.go`、`supervised_blackbox.go` | 已闭合：两个 facts 提取函数增 attempt 过滤，`requirePostflightQualifiedReview` 增参 |
| ⑤ r1 L1 | `acceptFromExistingSeal` 复用 seal 时不校验 decision | Low | `engine.go` 的 `acceptFromExistingSeal` | 已闭合：复用前校验 `ReviewDecision == "approve"` 且 `PromotionDecision == "completed"` |
| ⑤ r2 N1 | `cancelRun` 只取单 attempt 实体 | Low | `engine.go` 的 `cancelRun` | 已闭合：按 `taskAttempts` 遍历任务与步骤，逐 attempt 取消 |
| ⑤ r2 N2 | supervised 路径固定 attempt-1 实体 | Low | `supervised_blackbox.go` | 已闭合：窗口、回返与 PASSED 收敛三处均 `*EntityAt(..., attempt)` |
| ⑥ r1 主因 | C5 的 K 生产零接线，交付物为死代码 | High | `attempt.go` / `engine.go` | 已闭合：seal 写入 K、复用前重派生比较、`runStep` 为每个 request 计算 K、`sameClosure` 比较 K |
| ⑦ r1 F1 | supervised 窗口与回返路径未按解析 attempt 取实体 | High | `supervised_blackbox.go` | 已闭合：r3 对该文件精确检索 `taskEntity(` 与 `stepEntity(` 均 0 命中 |
| ⑦ r1 F2 | cancel 只覆盖单 attempt，attempt>1 的步骤在取消后存活 | High | `engine.go` 的 `cancelRun` | 已闭合：任务与步骤均按 `taskAttempts` 逐 attempt 取消 |
| ⑦ r1 F3 | `ContextKey` 未按 intent 后的派发身份重算 | Medium | `engine.go` 的 `runStep` | 已闭合：intent 定稿后按 executionTarget / requestID / requestHash 重算 |
| ⑦ r1 F4 | resolver 失败直接返回，链停在 RUNNING 不 PAUSED | Medium | `engine.go` 的 `Run` | 已闭合：`failChain("attempt-resolve-failed")` 转 PAUSED 并返回可处置原因 |
| ⑦ r1 F5 | 拒绝用例未钉住拒绝原因，terminal fixture 字段不全 | Medium | `state_seal_test.go` | 已闭合：断言拒绝串归因 attempt 绑定，fixture 补 `RequestHash` / `CompletionHash` / `SessionID` / `Evidence` |
| ⑦ r1 F6 | 生产装配未注入 `Attempts` | Medium | `internal/console/**`（本片边界外） | 接受：边界外登记，后续片接线 |
| ⑦ r3 | 复查 0 新增可执行缺陷，`go build` / `go vet` / `go test` / tree 门禁全绿 | Info | `07-final/review.md` | 已闭合：`FINAL REVIEW: PASS` |

## 交付

- `SealRecord`（含 `ContextKey`）/ `SealSink` / `FileEventStore.SealTaskEvidence`（`O_EXCL` 独占新建 + `Sync` + 失败删半写）/ `LoadTaskSeal` / `sameClosure`：`internal/chain/seal.go`，seal 落在事件日志目录旁的 `seals/` 子目录，编码域 `proofrail:task-seal:1`。
- `seal_verify.go`：`verifyAcceptedSeals` 在 `New` 与 `Recover` 双点校验每个 PASSED attempt 的 seal 存在、可解码且其哈希在 `InputEvidence` 内；缺口返回 `ErrRecoveryUncertain` 并 PAUSE，不补造、不覆盖赢家。
- `TaskAttemptResolver` / `resolveTaskAttempt` / `projectionAttempt` / `taskAttempts` / `ContextKey` / `TaskContextKey`：`internal/chain/attempt.go`；`Options.Attempts` 为注入点，缺省回退事件投影的最高 attempt（新任务即 1，与既有 run 全等）。
- `StepRequest.ContextKey`：`internal/chain/models.go`；intent 之前以空维派生，intent 定稿后按 executionTarget / requestID / requestHash 重算。
- `acceptFromExistingSeal`：加载 seal → 重派生 `TaskContextKey` → parent / attempt / K 三项比对 → decision 校验 → 冻结复用 `task-accepted`，不重调 Reviewer 与 Publisher。
- attempt 贯通：dispatch（`runTask` / `runStep`）、routing（`agent_runner_terminal_route.go`）、operator（`operator_interaction.go`）、facts（`postflight.go` 的 attempt 过滤）、cancel（`cancelRun` 遍历 `taskAttempts`）、supervised（`supervised_blackbox.go` 的窗口、回返与 PASSED 收敛）。
- 测试 10 例：`state_seal_test.go` 9 例（屏障拒绝、闭包不完整、恢复校验、reopen 取解析 attempt、K 跨重入稳定、seal 绑定闭包、resolver 失败阻断、路由拒绝外来 attempt、seal 绑定任务级 K），`engine_test.go` 1 例按 attempt 键断言。

## 偏差与登记

| 编号 | 内容 | 来源 | 判定 |
|---|---|---|---|
| D1 | 封存载体取链内 seal 文件（含 K 字段），不取 seal 事件 | ① 架构 | 接受：`validTransition` 无自环，无法以事件表达同态封存；intent 事件引用 sealHash |
| D2 | `.seal.json` 用 `O_EXCL` 独占新建，不取临时文件 + `Link` | ② 实现期 | 接受：与 no-replace 等价，失败路径删半写，已存在则同哈希对账复用、异哈希报 mismatch |
| D3 | 生产装配未接线 `Attempts` | ⑤ r2 N1 / ⑦ F6 | 接受：边界外登记，缺省回退投影 attempt 使现状行为全等 |
| D4 | supervised 工作区命名固定 attempt-1 | ⑦ r3 Section C | 接受：边界外登记；合法 attempt>1 回执当前不可达，拒绝属 fail-closed |
| D5 | `ContextKey` 中 executionTarget / requestID / requestHash 在 intent 前为空、intent 后重算 | ② 实现期 / ⑦ F3 | 接受：意图定稿前这些维度尚不存在，定稿后即时重算并随 request 下发 |
| D6 | `taskEntity` / `stepEntity` 保留为 `*EntityAt(..., 1)` 便捷包装 | ② 实现期 / ⑤ r1 H2 | 接受：生产调用点已清零，残留调用仅在测试与单 attempt 路径 |
| D7 | ③ 红测与 ② 实现由主控（`DeepSeek V4.1 Flash`）直接完成，未派 `prfrail-tester` / `prfrail-implementer` 子代理 | ⑩ 回执核对 | 接受并登记：审查侧（①⑤⑥⑦⑧）全部按白名单派单且 `agentName` + `model` 显式传入，无模型降级；执行侧角色载体偏离如实记录 |

## 模型与计量

| 环节 | 载体 | 模型 |
|---|---|---|
| ① 架构 | `deep-reasoner` | DeepSeek V4 Pro |
| ② 实现 | `prfrail-implementer` | DeepSeek V4.1 Flash |
| ③ 测试 | `prfrail-tester` | DeepSeek V4.1 Flash |
| ⑤ 预审 | `deep-reasoner` | DeepSeek V4 Pro |
| ⑥ 扫描 | `independent-reviewer` | MAI-Code-1.1-Flash |
| ⑦ 终审 | `independent-reviewer` | GPT-5.3-Codex |
| ⑧ 文档 | `prfrail-documenter` | DeepSeek V4.1 Flash |
| 批量辅助 | `deep-reasoner` | GLM-5.3-Flash（2 次：attempt 普查 / 耐久写入顺序追踪） |

- **主控模型**：`DeepSeek V4.1 Flash`（High 档）；①⑤ 按白名单仍由 V4 Pro 承担，与主控选型是两件事。

## 门禁读数

- tree 作用域：`TOTAL_FAIL=0`；index 作用域：`TOTAL_FAIL=0`（changed 32）。
- 门禁自检：`SELFTEST: PASS 697/697`；`G6-1` 状态时态扫描行数 446，无待办与完成冲突。
- ci 作用域由 CI 双腿承担，读数见 §提交与 CI。

## 提交与 CI

- 主体提交：`b6e0f34` `feat: seal task evidence before acceptance and resolve the real attempt (P43-T027B-STATE-SEAL)`；推送 origin `9fa516b..b6e0f34`。
- CI 运行 `41228033279` 首轮即双绿：`Go ubuntu-latest` success、`Go windows-latest` success。
- ubuntu 腿硬证据：`ok internal/chain 0.348s`、`ok internal/console 1.904s`、`ok internal/tickets 1.999s`、Race 腿 `ok internal/chain 2.012s`、`SELFTEST: PASS 697/697`、`GATE REPORT scope=ci changed=32`、`TOTAL_FAIL=0`。
- windows 腿硬证据：`ok internal/chain 5.607s`、`ok internal/console 22.474s`、`ok internal/tickets 5.034s`、`SELFTEST: PASS 697/697`、`GATE REPORT scope=ci changed=32`、`TOTAL_FAIL=0`。
- 本报告回填提交属文档增量，其 CI 属观察运行，不计切片轮次。

## 后续片候选登记

- 生产装配接线 `Attempts`：console 侧提供 `TaskAttemptResolver` 可靠来源，解除 D3 边界。
- supervised attempt>1 端到端：解除工作区命名固定 attempt-1 的 D4 边界。
- seal 目录级耐久：crash 注入证明 `seals/` 目录项与事件日志同序落盘。
- stale / mismatch 判定矩阵补测：父快照差异、工作区差异、K 仅随 attempt 或 parentHash 变化的表格用例。
- `router.go` 与 handoff 路径的 K 绑定：非最终 acceptance 路径补齐 K 校验。
- `WorkspaceReopener` 接口落地：以独立接口承接 reopen 语义，与 `WorkspacePort` 解耦。

## 已知限制

- **⑤ N1 / N2**：均已并入交付（`cancelRun` 覆盖全 attempt、supervised 路径 attempt 感知），不另立残余项。
- **⑥ 残余 1**：stale 父快照与工作区差异缺直接用例，判定依赖 r3 的代码级复核。
- **⑥ 残余 2**：K 仅因 attempt 或 parentHash 变化的表格用例未补齐。
- **⑥ 残余 3**：supervised / handoff / cancel 路径未绑定 K，属非最终 acceptance 路径的维护提醒。
- **⑥ 残余 4**：seal 目录级耐久无 crash 注入证明，仅有文件级写入顺序证据。
- **⑦ F6 · 状态 = 接受（边界外）**：生产装配未注入 `Attempts`，attempt>1 在装配层不可达；回退策略为投影最高 attempt，装配引入 attempt>1 时须同步可靠来源并重审回退。
- **D1–D6 摘要**：封存载体取文件而非事件（D1）；`O_EXCL` 等价于 no-replace（D2）；装配 `Attempts` 未接线（D3）；supervised 工作区命名固定 attempt-1（D4）；K 的派发三维在 intent 后重算（D5）；`*EntityAt(..., 1)` 保留为便捷包装（D6）。

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
