# batch-02 — s1-reach 十门现状枚举

生成者：deep-reasoner + GLM-5.3-Flash（P76 批量辅助 2，事实枚举，不含判定）

输入（实际读取的文件）：

- `docs/validation/s1-reach-reachability.md`（行数 152；全文读取，重点 §2 十门表，另核对头部状态行、§5、§13、§14）
- `docs/validation/s1-reach-reachability_EN.md`（行数 152；全文读取，双语对照）
- `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md`（仅用于确认批量 2 的任务范围与产出路径，非被测引文来源）

枚举口径说明（事实性）：

- 被测文档 §2 的表列名为「门 | 判据（可判真假） | 证据 | 判定 | 解锁条件」（CN `:20` / EN `:20`）；本文件把「判定」列称为「状态」，把「证据」列称为「依据 file:line」。
- 十门全部位于 §2 的表体内（CN `:22`–`:31` / EN `:22`–`:31`），无一门仅在其它章节出现；§2 之外的十门相关陈述另行枚举（见四）。
- 全部引文逐字照录；行号为 2026-11-16 当前磁盘快照。被测文档自身对行号时效的声明见三.4。
- 语境引文（总表前置）：报告头部状态行、判定闭集声明、一句话结论。

## 一、十门总表

| 门 | 名称 | 判据原文（逐字摘录） | 行号 | 当前状态（文档原文状态词） | 解锁条件（原文） | 依据 file:line |
|---|---|---|---|---|---|---|
| G-A | AT-23 真实候选腿 | CLI 可真实启动一个经兼容性准入的 AI 候选并计费 | CN `:22` / EN `:22` | 不可达 | 需要 T027 完成产品装配（当前证据 `runtime_noop.go:99` 显示 postflight 仍为 fail-closed 桩）且存在兼容候选（当前证据 T026 段显示 2 项 unsupported 未变） | `internal/console/runtime.go:33`（`ExecuteNoopRun`，noop-only 路径）、`internal/console/runtime_noop.go:99` 与 `:106`（`RunPostflight` fail-closed 桩）；`docs/DEV_PLAN.md` 的 T026 段（10 verified / 2 unsupported ⇒ 候选不兼容）与 T027 段（`BLOCKED / NOT IMPLEMENTED`） |
| G-B | AT-23 机制腿 | 隔离 workspace / timeout / 日志缺失 / 未知恢复 / exit 0 后重扫五面在真实宿主执行 | CN `:23` / EN `:23` | 已执行不独立成立 | 无需新解锁（已执行）；该腿只在 G-A 满足后才并入 AT-23 整体判定；不得写成"机制腿通过即 AT-23 通过" | `docs/t027/REMAINING_SLICES.md` 的 B4 小节；`docs/validation/t027-at23-e2e.md` |
| G-C | AT-23 整体 | G-A 与 G-B 同时满足并含真实计费证据 | CN `:24` / EN `:24` | 未通过 | 需要 G-A 解锁后重跑真实 E2E，当前证据 `s1-exit.md` 显示 AT-23 仍为 `BLOCKED` | `docs/validation/s1-exit.md` 的 AT-23 这一行（`BLOCKED`）；B4 小节的四事分立 |
| G-D | T028 与 AT-24 | `supervised-black-box` 显式模式实现完成且 AT-24 判定成立 | CN `:25` / EN `:25` | 已执行不独立成立 | 真实 CLI 端到端仍依赖 T027-B（D-18），需 P5 的 R-3 解锁；该门**与候选兼容准入无关**（黑箱模式不要求能力矩阵，见 `docs/ADR_REGISTER.md:22` 的 ADR-012） | `docs/DEV_PLAN.md:119`（T028 `COMPLETED`，AT-24 在 (b)/(c) 裁定下条件闭合）；`docs/validation/s1-exit.md:44`（已更正为「条件闭合」） |
| G-E | 最终 ZIP 绑定与复验 | 最终 ZIP 绑定发行记录并复验 | CN `:26` / EN `:26` | 未通过 | 需要 T018 其余门闭合后执行绑定与复验，当前证据显示该项仍缺失 | `docs/ADR_REGISTER.md:19`（ADR-009：最终 ZIP 仍须绑定发行记录并复验）；`docs/validation/s1-exit.md` 的「正式安装包教程」行（`PARTIAL`） |
| G-F | EOL 通知方式 | ADR-010 获产品所有者批准 | CN `:27` / EN `:27` | 未通过 | 需要产品所有者明确批准 ADR-010，当前证据显示其状态为 `PROPOSED / NOT APPROVED` | `docs/ADR_REGISTER.md:20`（ADR-010 状态 `PROPOSED / NOT APPROVED`） |
| G-G | 产品所有者 S1 exit 批准 | 产品所有者作出批准（人工门） | CN `:28` / EN `:28` | 未知 | 需要产品所有者作出 S1 exit 批准，当前证据显示尚未作出；本片不预测其时点 | `docs/validation/s1-exit.md` 的「产品所有者最终 S1 批准」行（`BLOCKED`） |
| G-H | T018 整体 | 依赖字段全部满足且四道剩余门闭合 | CN `:29` / EN `:29` | 不可达 | 需要 T027 与 T028 完成，当前证据显示两者未完成；**注意门级与任务级不同源**（见下方注） | `docs/DEV_PLAN.md:94`（依赖字段含 T027 与 T028）、`:96`（状态 `IN_PROGRESS / BLOCKED`） |
| G-I | S1 exit 总体 | AT-01 至 AT-24 全过且 T018 完成且 RFC §14 全部 exit | CN `:30` / EN `:30` | 不可达 | 需要 G-C 与 G-D 与 G-E 与 G-F 与 G-G 全部满足，当前证据显示其中任一均未满足 | `docs/validation/s1-exit.md` 头部结论行（`BLOCKED`）与 RFC §14 表（核心任务链 `BLOCKED`、三项 `PARTIAL`、业主批准 `BLOCKED`） |
| OB-21 | 自身 | ①② 两个子项是否仍成立 | CN `:31` / EN `:31` | 未通过（处置词保持待议） | 需要候选兼容或装配出现进展才触发其「待议时机」，当前证据显示两者均无变化 | ① 未装配（`runtime.go:33` 与 `runtime_noop.go:99` 未变）；② 中「`DR-1 未修`」子项已消除（DR-1 闭环，`docs/validation/evidence/dr-1-2026-09-24/`），「T026 不兼容」与「B2 §8.3 与 §8.4」未变（`docs/validation/B2-EXTERNAL-ENFORCEMENT.md` §8 第 3 与第 4 条） |

语境引文（逐字；在总表之前给出，供映射时定位状态口径）：

- 报告头部状态行：
  - CN（L4）：> **状态：已停止——未完成**（⑩ 停点；未完成清单见 §13；commit / push 须同轮授权）。
  - EN（L4）：> **Status: stopped — incomplete** (the ⑩ stop point; the incomplete list is in §13; commit and push need same-turn authorization).
- 判定闭集声明：
  - CN（L18）：> 判定闭集 = 不可达 / 未通过 / 已执行不独立成立 / 未知。`未知` 仅用于人工门。解锁句模板固定为「需要 X 具备 Y，当前证据 Z 显示…」，Z 必须是可定位的文件或证据包。
  - EN（L18）：> The verdict set is: unreachable / not passed / executed but not standing alone / unknown. `unknown` is reserved for human gates. The unlock sentence uses the fixed template "X needs Y, and evidence Z currently shows ...", where Z must be a locatable file or evidence pack.
- 一句话结论：
  - CN（L14）：- **一句话结论**：十门中**八门不可达或未通过**（G-A、G-C、G-D、G-E、G-F、G-H、G-I 与 OB-21 自身）、**一门未知**（G-G，人工门）、**一门已执行但不独立成立**（G-B）；DR-1 只消除了 OB-21 ② 的一个子项（`DR-1 未修`），**未改变 AT-23 的不可达判定**。
  - EN（L14）：- **One-line verdict**: **eight of the ten are unreachable or not passed** (G-A, G-C, G-D, G-E, G-F, G-H, G-I and OB-21 itself), **one is unknown** (G-G, a human gate) and **one was executed but does not stand alone** (G-B); DR-1 removed only one sub-item of OB-21 ② (the `DR-1 未修` clause) and did **not** change the AT-23 unreachability verdict.

## 二、各门详细摘录

### G-A

- 判据原文：CLI 可真实启动一个经兼容性准入的 AI 候选并计费
- 状态原文：不可达
- 解锁条件原文：需要 T027 完成产品装配（当前证据 `runtime_noop.go:99` 显示 postflight 仍为 fail-closed 桩）且存在兼容候选（当前证据 T026 段显示 2 项 unsupported 未变）
- 证据原文：`internal/console/runtime.go:33`（`ExecuteNoopRun`，noop-only 路径）、`internal/console/runtime_noop.go:99` 与 `:106`（`RunPostflight` fail-closed 桩）；`docs/DEV_PLAN.md` 的 T026 段（10 verified / 2 unsupported ⇒ 候选不兼容）与 T027 段（`BLOCKED / NOT IMPLEMENTED`）
- EN 对照（判据）：the CLI can really launch a compatibility-admitted AI candidate and bill it
- EN 对照（状态）：unreachable
- EN 对照（解锁条件）：T027 must complete the product assembly (evidence `runtime_noop.go:99` currently shows postflight still a fail-closed stub) and a compatible candidate must exist (evidence: the T026 block still shows two unsupported items)
- 出处：CN `:22` / EN `:22`

### G-B

- 判据原文：隔离 workspace / timeout / 日志缺失 / 未知恢复 / exit 0 后重扫五面在真实宿主执行
- 状态原文：已执行不独立成立
- 解锁条件原文：无需新解锁（已执行）；该腿只在 G-A 满足后才并入 AT-23 整体判定；不得写成"机制腿通过即 AT-23 通过"
- 证据原文：`docs/t027/REMAINING_SLICES.md` 的 B4 小节；`docs/validation/t027-at23-e2e.md`
- EN 对照（判据）：isolated workspace / timeout / missing logs / unknown recovery / rescan after exit 0 are executed on a real host
- EN 对照（状态）：executed but not standing alone
- EN 对照（解锁条件）：no new unlock needed (already executed); this leg enters the AT-23 verdict only once G-A holds; it must never be written as "the mechanism leg passed, so AT-23 passed"
- 出处：CN `:23` / EN `:23`

### G-C

- 判据原文：G-A 与 G-B 同时满足并含真实计费证据
- 状态原文：未通过
- 解锁条件原文：需要 G-A 解锁后重跑真实 E2E，当前证据 `s1-exit.md` 显示 AT-23 仍为 `BLOCKED`
- 证据原文：`docs/validation/s1-exit.md` 的 AT-23 这一行（`BLOCKED`）；B4 小节的四事分立
- EN 对照（判据）：G-A and G-B both hold and real billing evidence exists
- EN 对照（状态）：not passed
- EN 对照（解锁条件）：G-A must be unlocked first and the live E2E re-run; evidence `s1-exit.md` currently shows AT-23 still `BLOCKED`
- 出处：CN `:24` / EN `:24`

### G-D

- 判据原文：`supervised-black-box` 显式模式实现完成且 AT-24 判定成立
- 状态原文：已执行不独立成立
- 解锁条件原文：真实 CLI 端到端仍依赖 T027-B（D-18），需 P5 的 R-3 解锁；该门**与候选兼容准入无关**（黑箱模式不要求能力矩阵，见 `docs/ADR_REGISTER.md:22` 的 ADR-012）
- 证据原文：`docs/DEV_PLAN.md:119`（T028 `COMPLETED`，AT-24 在 (b)/(c) 裁定下条件闭合）；`docs/validation/s1-exit.md:44`（已更正为「条件闭合」）
- EN 对照（判据）：the explicit `supervised-black-box` mode is implemented and the AT-24 verdict holds
- EN 对照（状态）：executed but not standing alone
- EN 对照（解锁条件）：the real CLI end-to-end path still depends on T027-B (D-18) and needs R-3 of P5; this gate is **independent of candidate compatibility admission** (the black-box mode does not require a capability matrix, see ADR-012 at `docs/ADR_REGISTER.md:22`)
- 出处：CN `:25` / EN `:25`

### G-E

- 判据原文：最终 ZIP 绑定发行记录并复验
- 状态原文：未通过
- 解锁条件原文：需要 T018 其余门闭合后执行绑定与复验，当前证据显示该项仍缺失
- 证据原文：`docs/ADR_REGISTER.md:19`（ADR-009：最终 ZIP 仍须绑定发行记录并复验）；`docs/validation/s1-exit.md` 的「正式安装包教程」行（`PARTIAL`）
- EN 对照（判据）：the final ZIP is bound to the release record and re-verified
- EN 对照（状态）：not passed
- EN 对照（解锁条件）：the remaining T018 gates must close first and the binding and re-verification then be executed; evidence currently shows this item still missing
- 出处：CN `:26` / EN `:26`

### G-F

- 判据原文：ADR-010 获产品所有者批准
- 状态原文：未通过
- 解锁条件原文：需要产品所有者明确批准 ADR-010，当前证据显示其状态为 `PROPOSED / NOT APPROVED`
- 证据原文：`docs/ADR_REGISTER.md:20`（ADR-010 状态 `PROPOSED / NOT APPROVED`）
- EN 对照（判据）：ADR-010 is approved by the product owner
- EN 对照（状态）：not passed
- EN 对照（解锁条件）：the product owner must explicitly approve ADR-010; evidence currently shows `PROPOSED / NOT APPROVED`
- 出处：CN `:27` / EN `:27`

### G-G

- 判据原文：产品所有者作出批准（人工门）
- 状态原文：未知
- 解锁条件原文：需要产品所有者作出 S1 exit 批准，当前证据显示尚未作出；本片不预测其时点
- 证据原文：`docs/validation/s1-exit.md` 的「产品所有者最终 S1 批准」行（`BLOCKED`）
- EN 对照（判据）：the product owner grants approval (a human gate)
- EN 对照（状态）：unknown
- EN 对照（解锁条件）：the product owner must grant the S1 exit approval; evidence currently shows it has not been granted, and this slice predicts no date
- 出处：CN `:28` / EN `:28`

### G-H

- 判据原文：依赖字段全部满足且四道剩余门闭合
- 状态原文：不可达
- 解锁条件原文：需要 T027 与 T028 完成，当前证据显示两者未完成；**注意门级与任务级不同源**（见下方注）
- 证据原文：`docs/DEV_PLAN.md:94`（依赖字段含 T027 与 T028）、`:96`（状态 `IN_PROGRESS / BLOCKED`）
- EN 对照（判据）：the dependency fields are all satisfied and the four remaining gates are closed
- EN 对照（状态）：unreachable
- EN 对照（解锁条件）：T027 and T028 must complete; evidence currently shows neither has; note the **gate level and task level do not share a source** (see below)
- 出处：CN `:29` / EN `:29`

### G-I

- 判据原文：AT-01 至 AT-24 全过且 T018 完成且 RFC §14 全部 exit
- 状态原文：不可达
- 解锁条件原文：需要 G-C 与 G-D 与 G-E 与 G-F 与 G-G 全部满足，当前证据显示其中任一均未满足
- 证据原文：`docs/validation/s1-exit.md` 头部结论行（`BLOCKED`）与 RFC §14 表（核心任务链 `BLOCKED`、三项 `PARTIAL`、业主批准 `BLOCKED`）
- EN 对照（判据）：AT-01 to AT-24 all pass and T018 is complete and every RFC §14 exit holds
- EN 对照（状态）：unreachable
- EN 对照（解锁条件）：G-C, G-D, G-E, G-F and G-G must all hold; evidence currently shows none of them does
- 出处：CN `:30` / EN `:30`

### OB-21

- 判据原文：①② 两个子项是否仍成立
- 状态原文：未通过（处置词保持待议）
- 解锁条件原文：需要候选兼容或装配出现进展才触发其「待议时机」，当前证据显示两者均无变化
- 证据原文：① 未装配（`runtime.go:33` 与 `runtime_noop.go:99` 未变）；② 中「`DR-1 未修`」子项已消除（DR-1 闭环，`docs/validation/evidence/dr-1-2026-09-24/`），「T026 不兼容」与「B2 §8.3 与 §8.4」未变（`docs/validation/B2-EXTERNAL-ENFORCEMENT.md` §8 第 3 与第 4 条）
- EN 对照（判据）：whether its ① and ② sub-items still hold
- EN 对照（状态）：not passed (the disposition word stays 待议)
- EN 对照（解锁条件）：a candidate or assembly advance must occur to trigger its review timing; evidence currently shows neither has moved
- 出处：CN `:31` / EN `:31`

## 三、文档中的门级/任务级不同源陈述（逐字 + 行号）

1. 门级与任务级不同源（§2 表后独立陈述）
   - 出处：CN `:33` / EN `:33`
   - CN（L33）：- **门级与任务级不同源（重要）**：G-D 至 G-G 四道门**不依赖**候选或装配进展（T028 是独立模式），但 `docs/DEV_PLAN.md:94` 的 T018 **任务级依赖字段**含 T027 与 T028 ⇒ 四门全部清空后 T018 **仍不可勾选**。本报告按门级与任务级**分开**表述，不给单一结论。
   - EN（L33）：- **Gate level and task level do not share a source (important)**: gates G-D to G-G do **not** depend on candidate or assembly progress (T028 is a separate mode), yet the T018 **task-level dependency field** at `docs/DEV_PLAN.md:94` contains T027 and T028, so T018 **still cannot be ticked** once those four gates are empty. This report states the gate level and the task level **separately** and gives no single combined verdict.
2. 同一陈述在 G-H 行内的重复（行内片段）
   - 出处：CN `:29` / EN `:29`
   - CN（L29）：**注意门级与任务级不同源**（见下方注）
   - EN（L29）：note the **gate level and task level do not share a source** (see below)
3. OB-21 那一条的历史句与 §2 刷新结论不同源（§12.4 现象句仍含已被消除的子项）
   - 出处：CN `:35` / EN `:35`
   - CN（L35）：- **OB-21 这一行的历史句未改写**（指令 §11.1 的历史行口径；本片对授权歧义按 §1.7.2 从严处理）⇒ 其刷新结论落在本报告 §2 与 §12.4 的 OB-51 这一行，§12.4 该行的现象句仍含已被消除的子项。
   - EN（L35）：- **The historical sentence of the OB-21 row is not rewritten** (directive §11.1 historical-row wording; this slice handles the authorization ambiguity strictly per §1.7.2), so the refreshed verdicts live in §2 of this report and in the OB-51 row of §12.4, and that row's evidence cell still contains the sub-item that has since been removed.
4. 行号时效声明（§5 第 3 条；影响本文件全部行号的引用口径）
   - 出处：CN `:68` / EN `:68`
   - CN（L68）：3. **行号口径**：本报告引用的 `file:line` 取自 ① 阶段与 ② 阶段的实读；仓库存量文档的行号会随其它切片变动而漂移，引用时以**文件名与节号**为主、行号为辅。
   - EN（L68）：3. **Line-number basis**: the `file:line` references come from the real reads of the ① and ② stages; line numbers in existing documents drift with other slices, so **filename and section number are primary and line numbers secondary**.

## 四、文档中关于「十门」与「S1 Exit 八门」关系的原文陈述（若有的话）

关键词扫描（事实记录，2026-11-16 快照）：

- CN「十门」命中：`:14`、`:50`；CN「八门」命中：仅 `:14`。
- EN 侧无 "ten gates" 原词；对应陈述为 `:14` 的 "eight of the ten"。
- CN「S1 exit」命中：`:10`、`:28`、`:30`、`:130`（均为门名或结论句用语）。

逐字引文：

- CN（L50）：| ① 架构 | V4 Pro（`deep-reasoner`） | 已执行 | 方案书（十门判定框架、五点方法、验收判据、文件级改动点、停止条件） |
- EN（L50）：| ① architecture | V4 Pro (`deep-reasoner`) | executed | the plan (ten-gate framework, five-point method, acceptance criteria, file-level change points, stop conditions) |
- CN（L130）：- **本片结论不解除任何东西**：AT-23 真实候选腿仍不可达；T018 与 S1 exit 仍阻断。
- EN（L130）：- **This slice's verdicts unlock nothing**: the AT-23 live candidate leg stays unreachable and T018 and S1 exit stay blocked.
- CN（L14）与 EN（L14）的一句话结论见一、语境引文（该句在其后的括注中列出了 G-A、G-C、G-D、G-E、G-F、G-H、G-I 与 OB-21 自身，即「八门」在该句内指十门集合内的状态计数）。

缺失性记录（事实）：在上述两份输入文件内，未找到把「s1-reach 十门（G-A–G-I + OB-21）」与「S1 Exit 八门（S-G1–S-G8）」建立映射或对应关系的原文陈述；与「八门」相关的原文仅 `:14` 一处。
