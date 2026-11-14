## 面 1　C1–C7 与 N1–N4 的当前闭环状态

**`docs/validation/strategic-planning-2026-10-31.md`**：C1–C7、N1–N4 均**未出现**（全文 0 处编号命中；唯一正则命中为 :9 提交哈希 `397dc30` 内的 `c3` 字符序列，非编号语义）。

**`docs/validation/t027b-breakthrough.md`**：C1–C7、N1–N4 均**未出现**（全文 0 处命中）。

**`docs/validation/t027b-tactical-blueprint.md`**（C1–C7、N1–N4 全部出现于本文件）：

定义句原文：

```text
docs/validation/t027b-tactical-blueprint.md:17:- **批次 A（C1–C4）**：C1 以每进程唯一 reaper 与有界停机替换无界回收等待；C2 以 `ProcessSpec.OnSpawn` 与预登记 starting 所有者消除「已 spawn 却按未启动重试」的路径；C3 把限长读取与 PID 读取绑定同一已开对象；C4 以「原件封存先于 intent」重排终局提交顺序。
docs/validation/t027b-tactical-blueprint.md:18:- **批次 B（C5–C7）**：C5 以执行上下文摘要 `K=(run,task,step,attempt,requestId,requestHash)` 选择不可变上下文并共用 `outcomeOnce`；C6 以持久包装 `DurableCostLedger` 提供 `Reserve`、`RequireOutstandingReservation`、`Settle`；C7 以 `EvaluateT027BReadiness` 分离探针、真实验收与生产发布，消除 AT-23 自依赖。
docs/validation/t027b-tactical-blueprint.md:24:- **新漏洞**：N1 为 Job 归属建立前的宿主崩溃窗口；N2 为目录 helper 返回后的祖先替换；N3 为 `Reopen` 写死 attempt；N4 为 `ProbePermit` 额度可重复取得。
```

状态判定原文：

```text
docs/validation/t027b-tactical-blueprint.md:10:- **核心结论**：施工图覆盖 C1–C7，逐项给出文件与函数级改法、新增测试与断言点；Astra 判 `VERDICT: REVISE`——3 项 BLOCKED（C1 / C4 / C6）、4 项 PARTIAL（C2 / C3 / C5 / C7），并给出新漏洞 N1–N4。
docs/validation/t027b-tactical-blueprint.md:23:- **判定**：C1 / C4 / C6 为 BLOCKED；C2 / C3 / C5 / C7 为 PARTIAL，分别残留宿主崩溃窗口、目录替换窗口、attempt 写死与探针额度可重复消费。
docs/validation/t027b-tactical-blueprint.md:38:| 1 | Astra 判 REVISE 的 N1–N4 是否在本片内补一轮修订 | 待裁 | 11 月实施 |
docs/validation/t027b-tactical-blueprint.md:39:| 2 | `CONTRACT-BLOCK` 是否立项为契约修订片 | 待裁 | C4 与 C6 落地 |
docs/validation/t027b-tactical-blueprint.md:40:| 3 | C2 / C3 / C5 / C7 的 PARTIAL 是否可接受为 11 月起点 | 待裁 | T027-B 实施 |
```

逐项状态词（取自上述行）：C1=BLOCKED；C2=PARTIAL；C3=PARTIAL；C4=BLOCKED；C5=PARTIAL；C6=BLOCKED；C7=PARTIAL；N1–N4=待裁（:38）。

## 面 2　P39–P44 各自结论

编号事实：`docs/DEV_PLAN.md` 使用 **P39–P44** 编号；`T039`–`T044` 字样在该文件 0 命中。

小节标题逐字：

```text
docs/DEV_PLAN.md:177:### 11 月重构里程碑（T027-B）
```

P39–P44 与 P69 行逐字：

```text
docs/DEV_PLAN.md:180:- P39 持久预算：C6（DurableCostLedger 与跨进程共享锁，锁内重读后决策；报告 `docs/validation/p39-durable-budget{,_EN}.md`；准则 v1.60）
docs/DEV_PLAN.md:181:- P40 锁序：C1（已完成，报告 `docs/validation/p40-lock-reaper{,_EN}.md`）
docs/DEV_PLAN.md:182:- P41 归属：C2 与 N1（已完成，报告 `docs/validation/p41-spawn-ownership{,_EN}.md`）
docs/DEV_PLAN.md:183:- P42 绑定：C3 与 N2（已完成，报告 `docs/validation/p42-evidence-bind{,_EN}.md`）
docs/DEV_PLAN.md:184:- P43 封存：C4 与 C5 与 N3（已完成，报告 `docs/validation/p43-state-seal{,_EN}.md`）
docs/DEV_PLAN.md:185:- P44 门序：C7 与 N4（判定实现与测试已完成，报告 `docs/validation/p44-gate-readiness{,_EN}.md`；端到端接线已由 P69 完成）
docs/DEV_PLAN.md:186:- P69 接线：P44 readiness / permit 门接到生产装配与 CLI 入口（报告 `docs/validation/t027b-e2e-wiring{,_EN}.md`）
```

P39 结论句（前置状态段）逐字：

```text
docs/DEV_PLAN.md:165:T027 前置状态更新（2026-10-31，P39 持久预算）：P39-T027B-DURABLE-BUDGET 完成，C6 已落盘——新增 DurableCostLedger 与跨进程共享锁，`Reserve` / `Settle` / `RequireOutstandingReservation` 改为锁内重读后决策，提交失败不签发确认，`unknown` 占用重启后仍计入；准则 v1.60。证据：[中文](validation/p39-durable-budget.md) / [English](validation/p39-durable-budget_EN.md)。
```

P44 目标与总体结论（`docs/validation/p44-gate-readiness.md`）逐字：

```text
docs/validation/p44-gate-readiness.md:13:- **目标**：闭合 C7 与 N4——readiness 严格分离探针就绪、真实验收就绪与生产发布就绪三级，探针结论不得推进验收；`ProbePermit` 一次性消费，调用前原子预留，消费后失效，失败、超时与重启均不重置。
docs/validation/p44-gate-readiness.md:14:- **口径**：不得宣称 C7/N4 端到端闭合，只能宣称门控就绪与许可路径已满足当前切片验收前提。
docs/validation/p44-gate-readiness.md:30:- **⑦ r3**：`FINAL REVIEW: PASS`（0 新增缺陷），本片额度 3/3。
docs/validation/p44-gate-readiness.md:165:- **R-3**：C7/N4 未端到端接线，只能宣称门控就绪与许可路径已满足当前切片验收前提，不得宣称端到端闭合。
docs/validation/p44-gate-readiness.md:166:- **端到端接线已由 P69 `T027B-E2E-WIRING` 完成（报告：`docs/validation/t027b-e2e-wiring{,_EN}.md`）**
```

## 面 3　P69 接线范围

改动文件清单（§2 接线点清单表 7 行）逐字：

```text
docs/validation/t027b-e2e-wiring.md:18:| 1 | `internal/console/readiness_gate.go` | 增 | `requireProbePermit` 三件套（Evaluate(GoalProbe) → Require → release/Complete）单一来源 | 新文件（+29） |
docs/validation/t027b-e2e-wiring.md:19:| 2 | `internal/console/cli.go` | 改 | `--probe-permit` 旗标 + 路由前 presence 判定（只判定不消费）+ `runUsage` + 两处调用点转发 permit + import gates | +32 / −4 |
docs/validation/t027b-e2e-wiring.md:20:| 3 | `internal/console/runtime_agent.go` | 改 | options 增 `ProbePermit` 字段；Step 0 判定（stub 后、Step 1 前）；Require（Step 1 全过后、Step 2 marker 前）+ `defer release` | +23 / −0 |
docs/validation/t027b-e2e-wiring.md:21:| 4 | `internal/console/runtime_supervised.go` | 改 | options 增同名字段；fresh 门（`SupervisedInteractionsLedgerPath` 不存在 = 首次起跑才 判定+Require+defer；可重入安全） | +23 / −0 |
docs/validation/t027b-e2e-wiring.md:22:| 5 | `internal/console/supervised_cmd.go` | 改 | `executeSupervisedRun` 签名加 permit 并赋 `options.ProbePermit`；`ExecuteSupervisedReturn` 不动 | +4 / −1 |
docs/validation/t027b-e2e-wiring.md:23:| 6 | `internal/console/readiness_gate_test.go` | 增 | 正例 1 + 反例 3 组共 7 条 | 新文件（+225） |
docs/validation/t027b-e2e-wiring.md:24:| 7 | 4 个既有测试文件 | 改 | 只补注入点、断言零改动（fixture `newAssemblyFixture`、harness `newSupervisedHarness`、`runtime_multitask_red_test.go` 10 处裸字面量、`p51SevenAgentRunnerArgs` 补旗标） | +32 / −16 |
```

接线范围原文：

```text
docs/validation/t027b-e2e-wiring.md:7:- 目标：把 P44 的 readiness / permit 判定接到生产装配路径与 CLI 入口，收口 P57 §3.3 的 R-1。
docs/validation/t027b-e2e-wiring.md:8:- 结论：`DESIGN: DONE`（①）；接线落地 `internal/console/**` 5 个生产文件（4 改 1 增）+ 5 个测试文件（4 改 1 增）；验收五条反例全部有实测载体（③ 落 7 条测试）。
```

未接线部分原文：

```text
docs/validation/t027b-e2e-wiring.md:78:- noop 路径与 `ai check` 路径本片不接线（`ai check` 的 permit 生命周期需独立设计，属后续切片）。
docs/validation/t027b-e2e-wiring.md:87:未改准则正文（`DELIVERY_DIRECTIVE{,_EN}.md` 零改动）；未改 `internal/gates/**`、`internal/chain/**`、`internal/adapters/**`、`internal/tickets/**`、`internal/guard/**` 判定本体；未改 schema / fixtures / `.github/workflows/**` / `go.mod`；未改 CONTRACTS 语义（§4.2 末句时态更新属【五】第 10 项，见 ⑧a2）；未改 v50–v55 快照；未 commit / push；未启动新切片；未推 gitee。
```

结论行原文：

```text
docs/validation/t027b-e2e-wiring.md:9:- 审查链：⑤ `PRE-REVIEW: PASS`（1 Low 已以注释级修复、④ 原位重跑全绿）；⑥ `INDEPENDENT SCAN: PASS`（0 发现）；⑦ `FINAL REVIEW: PASS`（0 发现；计量 = 终审 1 + 复审 0 + 盲审 0）。
docs/validation/t027b-e2e-wiring.md:65:- 五要素留痕（§7.8.5）：对象 = ⑦ 第 N 轮；轮次 = r1；依据 = §4.3 + §5.3 正常收敛；结论行原文 = `FINAL REVIEW: PASS`。
docs/validation/t027b-e2e-wiring.md:195:INDEPENDENT SCAN: PASS
```

## 面 4　B 段（enforcement 生产落地）的边界描述

9. `docs/validation/t027b-tactical-blueprint.md`：「B 段」「enforcement」「Sandbox」「OB-55」均 **0 命中 ⇒ 未出现**。证据目录 `docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/` 内 3 个文件共 5 处 `enforcement` 命中，均为代码路径/代码引用（非 B 段边界段落）；逐字一行：

```text
docs/validation/evidence/P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/input-list.md:14:| `internal/adapters/agent_runner_enforcement.go` | 217 | 9800 | Go 源（非测试） |
```

其余命中位置：`01-sol-blueprint/review-batch-b.md:41,42,46`；`02-astra-redteam/input-list.md:14`。

10. `docs/DELIVERY_DIRECTIVE.md` §12.4 OB-55 整行逐字：

```text
docs/DELIVERY_DIRECTIVE.md:1092:| OB-55 | 2026-10-02 | **enforcement 生产落地与 S1/S2 的边界冲突（本片 ① 前的只读盘点发现；用户 2026-10-02 裁定）**：T027 装配要求「无 enforcement 记录即拒绝启动候选」，而外部 OS 强制边界**尚未接入产品**（`DEV_PLAN.md:112` 原文），其平台无关的 `Sandbox` 接口又按 `docs/validation/B2-EXTERNAL-ENFORCEMENT.md` §8 第 8 条**被明确推迟到 S2**（理由 = S1 提前落地会产生无人调用的死代码）⇒ 两条要求方向相反。**本片裁定**：`T027-ASSEMBLY-MIN` 只做**拒绝路径**（该路径在 console 的 `run` 命令有真实调用者，**不构成**无人调用的死代码 ⇒ 不违反上述推迟理由）；真实边界接入产品的归属（B 段）不在本片裁定。 | **待议**：B 段立项时由 ① 阶段评估（候选范围 = 独立切片 / 提前落 `Sandbox` 接口 / 维持 S2） |
```

11. `docs/DEV_PLAN.md` T027-B 段：「B 段」字样 0 命中；OB-55 相关两行逐字：

```text
docs/DEV_PLAN.md:151:T027-B 立项登记（2026-10-31，取代 2026-10-13 的待评估项登记）：`T027-B`（防御性重构期）已立项；边界冲突 OB-55 与真实 CLI 端到端依赖仍然有效；设计基线由 P37 与 P38 定稿（C1–C7 判定、N1–N4 与 CONTRACT-BLOCK），执行按「11 月重构里程碑（T027-B）」逐项推进。证据：[中文](validation/t027b-tactical-blueprint.md) / [English](validation/t027b-tactical-blueprint_EN.md)。
docs/DEV_PLAN.md:153:T027-B 收口评估（2026-11-05，P57 `P5-P6-CLOSEOUT-ASSESSMENT`）：P5 判定 `T027-B` **不新立片**——本里程碑已把 C1–C7 与 N1–N4 分配到 P39–P44（`docs/DEV_PLAN.md:176` 至 `:181`），余项三叉为 R-1 P44 端到端接线（无队列位）、R-2 OB-55 归属（待 R-06 裁定）、R-3 真实 CLI 端到端（层 1 第 8 项，待 R-08 与 R-09）。证据：[中文](validation/p57-closeout-assessment.md) / [English](validation/p57-closeout-assessment_EN.md)。
```

## 面 5　文件存在性（Test-Path）

```text
docs/validation/strategic-planning-2026-10-31.md => True
docs/validation/t027b-breakthrough.md => True
docs/validation/t027b-tactical-blueprint.md => True
docs/validation/p44-gate-readiness.md => True
docs/validation/t027b-e2e-wiring.md => True
docs/DEV_PLAN.md => True
docs/DELIVERY_DIRECTIVE.md => True
```

输出字符数： 8382
BATCH-01 DONE
