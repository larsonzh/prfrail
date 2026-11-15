# ProofRail R3-REAL-CLI-ASSESSMENT 报告（P75 · 层 1 T027-B 主线里程碑余项评估 · 评估类 [SLICE] · M 档）

日期：2026-11-15。状态：② 阶段成文（① 已停点裁定；②a–②c 分步落地）。切片标识：`R3-REAL-CLI-ASSESSMENT`。

R-3 裁定：**(A) 暂不可起跑**（用户 2026-11-15 同轮裁定；附带裁定 A1–A5 见 §4）。

## 1. 执行摘要

- 本片回答 T027-B 主线里程碑的 R-3 余项三问：里程碑可否宣告完成、S-G8 可否推进、AT-23 阶段 1 闭合状态是否维持。
- 结论：R-3 = (A) 暂不可起跑；T027-B 宣告「完成（部分不适用）」；S-G8 保持 blocked 直到 S2 enforcement 接入；AT-23 阶段 1 条件闭合作为永久事实登记。
- 依据链：B1 已评估范围无可接受候选 + B2 外部 enforcement 边界尚未接入 + P74 仅落检查层（S1 默认 off）+ P61 / F2-FIX 已达成阶段 1 闭合。
- 下游：真实 CLI 端到端在 S1 内不可达；S1 Exit 收口时点不因本裁定前移；`p57` 与 `DEV_PLAN` 通过加法式指针与段落更新承接裁定。
- 本片不消耗探针、不改门表判定、不改契约与产品代码。

## 2. 现状核对

### 2.1 T027-B：C1–C7 / N1–N4 / R-1 / R-2

- C1–C7：P39 C6 / P40 C1 / P41 C2 / P42 C3 / P43 C4 与 C5 / P44 C7 均已落地（`docs/DEV_PLAN.md:180-185`）。
- N1–N4：定义见 `docs/validation/t027b-tactical-blueprint.md:24`；分别由 P41 / P42 / P43 / P44 落地。
- R-1：P69 `T027B-E2E-WIRING` 收口（`docs/validation/p57-closeout-assessment.md:72`）。
- R-2：P73 裁定 = 方案 b（`docs/validation/p57-closeout-assessment.md:74`），P74 落地 `Sandbox` 检查层（`docs/validation/sandbox-interface.md:8-9`）。

### 2.2 AT-23 阶段 1 条件闭合

- 现状：阶段 1 条件闭合达成（`docs/validation/s1-exit.md:43`；`docs/validation/at23-restricted-candidate.md:499`）。
- 前提①：机制腿全绿，由 F2 修复片 ⑦ 判定满足（`docs/validation/at23-restricted-candidate.md:423`）。
- 前提②：真实候选腿 ≥3 次记录 unsupported，P61 时点已满足（`docs/validation/at23-restricted-candidate.md:19-20`）。
- 边界：R-08.5 明写「条件闭合不构成真实候选兼容性证明」（`docs/validation/r08-r09-ruling.md:38-40`）。

### 2.3 S-G1 / S-G8 与八门判定

- S-G1：不可达 / 未通过（`docs/validation/p57-closeout-assessment.md:88`）。
- S-G8：未通过，依赖 T027-B（`docs/validation/p57-closeout-assessment.md:95`）。
- 其余门：S-G2 口径滞后；S-G3 与 S-G4 未通过；S-G5 未知（人工门）；S-G6 与 S-G7 不可达（同文件 `:89-94`）。
- P74 边界：`Sandbox` 检查层不解锁 S-G1 / S-G8（`docs/validation/sandbox-interface.md:133`）。

### 2.4 B1 / B2 与 R-08 / R-09

- B1：已评估范围内无可接受候选（`docs/validation/t027-b1-candidate-live-discovery.md:4`、`:111`）。
- B2：采用外部 OS 强制，边界内外的划分见 `docs/validation/B2-EXTERNAL-ENFORCEMENT.md:12`、`:18-19`。
- B2 限制：零能力容器内的候选会话失败（`docs/validation/B2-EXTERNAL-ENFORCEMENT.md:162-164`）。
- 矩阵口径：固定候选 10 verified / 2 unsupported / 0 unknown，整体不兼容（`docs/validation/s1-exit.md:3`）。
- R-08 要件：R-08.3 只用 P；R-08.4 VM 与代理都不允许；R-08.5 条件闭合；R-08.6 BYOK 平行路径（`docs/validation/r08-r09-ruling.md:28-57`）。

### 2.5 R-3 出处与定义

- 出处：`docs/validation/p57-closeout-assessment.md:70`（层 1 第 8 项 `AT23-RESTRICTED-CANDIDATE`）。
- 定位：评估项，不是再次执行（用户 2026-11-15 裁定明确）。

## 3. R-3 三选一对比

- (A) 暂不可起跑：维持现状并登记；成本为零；S-G8 保持未通过；不动摇既有闭合；可回滚性 = 无动作。
- (A) 不做：不重跑真实候选腿、不引入 P′、不换候选、不改门判定。
- (B) 分阶段起跑：机制腿复用 F2-FIX 证据，真实候选腿重跑 ≥3 次；成本 = 探针逐次声明 + 一个 M 档实施片；结果至多「条件闭合」再确认。
- (B) 不做：不引入 P′、不换候选族、不改判据、不推进 S-G8。
- (C) 引入 P′ 或换候选：需新候选评估片或装置前提；成本与论证义务最高；是唯一可能真正推进 S-G8 的路径。
- (C) 不做：不维持固定候选单通道；不在本片内落地（需新片）。

| 轴 | (A) | (B) | (C) |
|---|---|---|---|
| 做什么 | 维持并登记 | 机制腿复用 + 真实腿重跑 | 新候选评估或 P′ 装置 |
| 不做 | 不重跑 / 不换候选 | 不换候选 / 不推进 S-G8 | 本片不落地 |
| 成本 | 0 | 探针 + 一个 M 片 | 高（多片 + 论证） |
| T027-B 完成判定 | 口径可落地 | 口径可落地 | 不受影响 |
| S-G8 推进 | 不推进 | 不推进 | 仅换候选成功才可能 |
| AT-23 闭合 | 维持 | 维持 | 成功则升级 |
| S1 Exit 时点 | 不前移 | 不前移 | 最早 |
| 可回滚性 | 无动作 | 切片可回滚 | 装置须成对纪律 |

- ① 阶段推荐为 (B)；用户 2026-11-15 裁定为 (A)，并指出「与 R-08/R-09 基线一致」不成立——R-08/R-09 是 AT-23 阶段 1 的目标，阶段 1 已达成。推荐不代替裁定。

## 4. R-3 裁定登记

| 要素 | 内容 |
|---|---|
| 裁决人 | 用户 |
| 日期 | 2026-11-15 |
| 理由 | 阶段 1 闭合已达成、不需重复执行；S-G1 / S-G8 根因是「无真实 enforcement」；P74 默认 off 使重跑结果仍为 unsupported；(B) 不推进 S-G8；(A) 成本为零且不动摇既有闭合 |
| 依据 | 用户 2026-11-15 同轮裁定原文（含五项附带裁定 A1–A5 与 O-7 = ② 的裁定）；原文归档：`docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/02-batch/ruling-2026-11-15-r3.md` |
| 裁定结论 | R-3 = (A) 暂不可起跑 |

### 4.1 附带裁定 A1–A5

| 编号 | 内容 |
|---|---|
| A1 | T027-B 宣告口径 = 「完成（部分不适用）」：主线已全部落地，R-3 部分不适用（真实 CLI 端到端需 S2 enforcement 前提） |
| A2 | S-G8 保持 blocked 直到 S2 enforcement 接入；`s1-exit` 门表零改动；`p57` §4.2 的 S-G8 行加法式加指针 |
| A3 | AT-23 阶段 1 条件闭合作为永久事实登记；阶段 2 不适用；`p57` §4.2 的 S-G1 行零改动 |
| A4 | R-3 处置落点三处：`DEV_PLAN` T027-B 段、`REMAINING_SLICES` P75 段、`p57` §3.3 的 R-3 行 |
| A5 | 新增 ADR-045（决策记录；不受「不得作为先例」约束） |

### 4.2 开放问题裁定

- O-1 至 O-6 与 O-8：由本轮裁定与 A1–A5 覆盖，属 P75 内裁定。
- O-7：裁定 = ②（不另立片）；义务传递至 S1 exit 收口片（详见 §11）。
- ⑥ 跳过：本片产出为纯文档（评估方案 + 裁定登记 + 报告装配），按 §7.2.1 三条件判定可跳过，该判定待 ⑦ 确认；若 ⑦ 不确认则补跑 ⑥ 并以 `sixstage` 围栏替换本行。

## 5. T027-B 完成判定

- 判据：里程碑构成（`docs/DEV_PLAN.md` 的 11 月重构里程碑清单）与 P57 §3.3 的余项三叉（`docs/validation/p57-closeout-assessment.md:64-76`）。
- 主线：C1–C7 与 N1–N4 已由 P39–P44 落地；R-1 由 P69 收口；R-2 由 P73 裁定并由 P74 落地检查层。
- R-3 定性：**部分不适用**——其承载位是「真实 CLI 与真实候选下的端到端运行」，前提为 S2 enforcement 接入，S1 内不具备该前提。
- 结论（A1）：T027-B 宣告「**完成（部分不适用）**」；本片不把它写成「未完成」，也不把它写成「S1 Exit 通过」。
- 边界：本判定不改变 `s1-exit` 门表、不改变 AT-23 判定、不改变任何 `unproven` 标记。

## 6. S-G8 推进条件

- 判据原文：受监督路径可在真实 CLI 与真实候选下端到端运行（`docs/validation/p57-closeout-assessment.md:95`）。
- 现状：未通过，依赖 P5 余项 R-3；S-G8 与 S-G1 同源解锁。
- A2 裁定：S-G8 **保持 blocked 直到 S2 enforcement 接入**；`s1-exit{,_EN}.md` 门表零改动；`p57` §4.2 的 S-G8 行加法式加指针（落点：`docs/validation/p57-closeout-assessment.md:99`）。
- 三选一横向核对：在 (A) / (B) / (C) 任一裁定下，本片均不推进 S-G8——(B) 至多再确认「条件闭合」，(C) 只有在换候选成功且真实端到端跑通后才可能推进。
- 边界：本片不改变 S-G8 的判据与判定；若后续片要求变更该判据，须按停升级边界由用户裁定。

## 7. AT-23 阶段 1 闭合状态维持

- **永久事实登记**（A3）：AT-23 阶段 1 条件闭合是历史事实，登记后不因 R-3 结论撤销、不因后续切片改写。
- P61 时点记录（2026-11-08）：前提②真实候选腿已跑过启动 / 执行 / 停机并记录 unsupported ⇒ 满足；前提①机制腿因 F2 继承漂移不满足（`docs/validation/at23-restricted-candidate.md:19-20`）。
- F2-FIX 时点记录（2026-11-09）：⑦ 终审 `FINAL: PASS WITH FIXES` ⇒ 前提①满足 ⇒ **AT-23 阶段 1 条件闭合达成**（`docs/validation/at23-restricted-candidate.md:423`、`:499`）。
- 阶段 2：不适用（用户 2026-11-09 裁定 A1），本片不改该裁定。
- 边界：R-08.5 明写「条件闭合不构成真实候选兼容性证明」；S-G1 的**门级判定**（不可达 / 未通过）与 AT-23 的**判定状态**是两个不同维度，本片不改变 S-G1 判定（A3）。

## 8. S1 Exit 收口时点影响

- 八门现状（P57 时点 + P74 边界）：S-G1 与 S-G8 未通过；S-G3 与 S-G4 未通过；S-G5 为人工门（未知）；S-G6 与 S-G7 不可达；S-G2 口径滞后。
- R-3 = (A) 的影响：**收口时点不前移**；S1 Exit 的真实推进依赖 S2 enforcement 接入，而非 S1 内的重复评估。
- 准入约束（新增）：S1 exit 收口片起跑时，若 §11.2 所列残余时态滞后（第 3 项）尚未处置，则作为收口片的**准入条件之一**——收口片不得在未处置该残余项的情况下宣告 S1 Exit 通过。
- 边界：本片不给出 S1 Exit 的通过时间，也不改变任何门的判定。

## 9. 真实 CLI 端到端可达路径

- 事实链：S-G8 判据要求真实端到端 → 需要兼容候选 → B1 在已评估范围内未找到可接受候选（`docs/validation/t027-b1-candidate-live-discovery.md:4`、`:111`）→ B2 的外部队列边界尚未接入产品（`docs/validation/B2-EXTERNAL-ENFORCEMENT.md:12`、`:18-19`）→ 零能力容器内的候选会话失败（同文件 `:162-164`）。
- S1 内可达性：**不可达**——P74 的平台实现恒 `Unknown`，且触发源默认 `off`；即使置 `enforce` 也只是显式拒绝，不构成 enforcement。
- 可达路径：① S2 enforcement 接入 + 兼容候选（主路径）；② 经用户裁定 (C) 换候选族或引入 P′（需另裁 R-08.3 / R-08.4 的豁免前提）。
- 结论：本片不消耗探针、不启动端到端运行；S1 侧的「不做」由 ADR-045 固定。

## 10. 证据缺口

- 当前已有：机制腿全绿证据、P61 真实候选腿的 unsupported 记录、B1 / B2 结论、P69 接线证据、P74 检查层证据、AT-23 阶段 1 闭合判定。
- 仍然缺失：受监督路径的真实端到端运行证据（含计费）、兼容候选证据、S2 真实 enforcement、S-G1 / S-G8 的解锁证据。
- 结论：缺口的闭环条件全部落在 S2 或用户裁定的 (C) 路径上，S1 内无法补齐。

## 11. 开放问题与后续建议

### 11.1 O-1 至 O-8 处置结果

| 编号 | 处置 | 说明 |
|---|---|---|
| O-1 | P75 内裁定 | R-3 = (A)，见 §4 |
| O-2 | P75 内裁定 | 随 (A) 失效（探针通道与额度口径不适用） |
| O-3 | P75 内裁定 | 随 (A) 失效（候选族与 P′ 论证义务不适用），保留为 S2 后续入口 |
| O-4 | P75 内裁定 | 口径 = 「完成（部分不适用）」，见 §5 |
| O-5 | P75 内裁定 | S-G8 三裁定下均不推进，见 §6 |
| O-6 | P75 内裁定 | 新立 ADR-045（`docs/ADR_REGISTER.md:56`） |
| O-7 | 裁定 = ②（不另立片） | 义务传递见 §11.2 与 §11.3 |
| O-8 | P75 内裁定 | 未触停升级边界，见 §6 与 §7 |

### 11.2 O-7 义务传递（三项逐项，不得汇总）

p57 §4.3 为 2026-11-05 的时点登记；下列第 1、2 项已由 2026-11-15 磁盘复核更正，第 3 项为残余。

| P75 §11.2 项 | 对应 P57 §4.3 原项 | 处置 |
|---|---|---|
| 第 1 项（`s1-exit.md:44` 的 AT-24 所在行） | 原第 1 项 | 已由 P64 更正 ⇒ 关闭 |
| 第 3 项（`s1-exit.md:52` 的核心任务链行） | 原第 2 项 | 唯一残余 ⇒ 移交 S1 exit 收口片 |
| 第 2 项（`s1-reach` 的 G-D 行） | 原第 3 项 | 已由 s1-reach 更正 ⇒ 关闭 |


1. 第 1 项（`docs/validation/s1-exit.md:44` 的 AT-24 所在行）：**2026-11-15 磁盘复核 = 已更正**——该行现为「条件闭合（2026-11-11，P64 验收）」，与 `docs/DEV_PLAN.md:119` 的 `COMPLETED` 一致 ⇒ 本项关闭，不作义务传递。
2. 第 2 项（`docs/validation/s1-reach-reachability.md` §2 的 G-D 行）：**2026-11-15 磁盘复核 = 已更正**——该行现引用 `DEV_PLAN.md:119` 并注明 s1-exit 已更正 ⇒ 本项关闭，不作义务传递。
3. 第 3 项（`docs/validation/s1-exit.md:52` 的核心任务链行）：**残余项，唯一需要传递**——该行前半仍写「T027 未完成、AT-23 未通过」，同一单元格后半写「AT-23 阶段 1 条件闭合达成」⇒ 行内口径不一致；修正时机 = S1 exit 收口片「重述门结论」时一并修正。

### 11.3 准入条件与后续建议

- 准入条件：S1 exit 收口片起跑时，若 §11.2 第 3 项（`s1-exit` 核心任务链行的行内口径不一致）尚未处置，则该片不得宣告 S1 Exit 通过；收口片起跑消息起草时必须引用 §11.2 的三项复核结论。
- 后续建议：① S2 enforcement 接入片（真实端到端的前置）；② S1 exit 收口片（承载 §11.2 的清账与门结论重述）；③ 本片不另立清账片。
- 已识别的契约修订点（只登记）：R-08.5 的「条件闭合」与 S-G8 判据「真实端到端运行」之间的口径映射尚未成文。

## 12. 审查链结论

- ① `DESIGN: DONE`（`deep-reasoner` + `DeepSeek V4 Pro`；含主控行号核验修正节）。
- ②a / ②b / ②c 分步落地，每步 `gate.js --all --scope=tree` 与 CN/EN 六形状核对全绿。
- ③ 一致性检查：主控自执行（载体偏离，见 §14）。
- ④ 门禁：三作用域与 SELFTEST 全绿（读数见 §13）。
- ⑤ 预审：r1 `PASS WITH FIXES`（7 Medium + 2 Low）→ 整改 → r2 `PASS WITH FIXES`（3 Low）→ 整改；额度 1+1 用尽。
- ⑥ 独立扫描：按 §7.2.1 三条件判定跳过（纯文档 / 无机器结构语义 / 不涉四类语义），由 ⑦ 确认 `SKIP CONFIRMED`。
- ⑦ 盲审：`BLIND REVIEW: PASS`（0 发现；按 §7.7 按比例抽查命中本片，先于终审）。
- ⑦ 终审：`FINAL REVIEW: PASS`（0 发现）；无 Medium+ ⇒ 不触发复审。
- 盲审与终审结论一致（均为 PASS 且 0 发现）⇒ 无需根因仲裁。
- 计量：⑦ = 盲审 1 + 终审 1 = 2，未超 M 档 1+2 与盲审 +1 的合计 4；批量辅助 2 次（下限 2）。

## 13. 门禁读数与裁定落点

- 门禁：`--scope=index` / `--scope=tree` / `--scope=ci` 均 `TOTAL_FAIL=0`；`--selftest` = `SELFTEST: PASS 737/737`；基础 `go test ./...` = `ok=15 FAIL=0`。
- A1 落点：`docs/DEV_PLAN.md:203`（`_EN:204`）。
- A2 落点：`docs/validation/p57-closeout-assessment.md:99`（S-G8 加法式指针；`s1-exit` 门表零改动）。
- A3 落点：本报告 §7（永久事实登记）；`p57` §4.2 的 S-G1 行零改动。
- A4 落点：`docs/DEV_PLAN.md:203`、`docs/t027/REMAINING_SLICES.md` P75 段（⑧ 完成记录）、`docs/validation/p57-closeout-assessment.md:76`。
- A5 落点：`docs/ADR_REGISTER.md:56`（`_EN:57`）。
- 计量：R = 1（⑤ r1 报 Medium+ 后的整改与重跑）；D = 0（临时件零违规）。

## 14. 已知边界

- `06-gates/` 的四个捕获为点时刻快照；包内自洽行反映捕获当时的声明状态，最终读数以 §13 与 ⑩ 回执为准。
- 裁定原文归档为**节录**：裁定主体（三选一裁定与五项附带裁定）逐字；过程指示段不转录，已在归档文件头注明范围。
- 批量辅助的子代理在计数时曾创建并即时删除临时文件（非本片临时目录），已确认零残留；不做为他片豁免依据。
- ③ 一致性检查由主控自执行（载体偏离），独立影响评估见报告异常处理章。

## 15. 明确未执行

- 未启动真实 CLI 端到端探针（(A) 方案下不消耗探针）。
- 未改动 `docs/validation/s1-exit{,_EN}.md` 与 `docs/validation/s1-reach-reachability{,_EN}.md`。
- 未改动 `CONTRACTS` / schema / fixtures / CI / 产品代码。
- 未替用户裁定；未 commit / push；未启动新切片。

### 15.1 异常处理与偏离登记

- P75-1（⑤ r1 报 Medium+ 后的整改）：O-7 三项未复核即搬运 p57 §4.3、两处悬空引用、缺「不做什么」、门禁捕获未刷新、裁定原文未归档 ⇒ 已整改并在原位重跑 ⑤（r2）。计入 R = 1。
- P75-2（作者期门禁口径修正，不计 R）：G8-b 元数据行混入载荷、G1(e2) 完成行、G7-f 新增报告声明、G6-4 过期字面量、G5(d) 围栏缺行、G8-b 自指哈希；均为同一撰写过程内的机械修正，已逐项复跑全绿。
- P75-3（载体偏离）：③ 一致性检查由主控自执行（§13.3.18）；独立影响评估 = 核对动作为事实比对，结论已由 ⑤ r2 与 ⑦ 盲审 / 终审独立复核；不得作为先例。
- P75-4（时序偏离）：freedom-list 定稿晚于 §1.7.3 要求的「⑥ 后 ⑦ 前」窗口（实际在 ⑧ 定稿）；影响评估 = 清单条目自步骤 0 起未变，⑦ 两读均按其内容核对越界；不得作为先例。
- P75-5（子代理临时文件）：批量辅助子代理在计数时创建并即时删除临时文件，已确认零残留；不作他片豁免依据。

## 16. 提交与 CI 事实链（回填）

- 提交：`ed2e26e`（`docs: rule R-3 as not startable and declare T027-B complete in part`），文件数 27、`+2287 / −0`。
- 推送：`origin/main` `c5c3a78..ed2e26e`（仅 origin；未推 gitee）。
- CI：run `43732384194`；`Go ubuntu-latest` success（3m19s）、`Go windows-latest` success（4m41s）；两腿 `GATE REPORT scope=ci changed=27`、`TOTAL_FAIL=0`、`SELFTEST: PASS 737/737`。
- 注记：本节由 CI 双绿后的 docs-only 提交补入（§11.2）；该回填提交自身的 CI run 只作观察记录。

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/00-briefing.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/01-arch/design-r1.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/02-batch/batch-01-t027b-at23-b1b2-facts.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/02-batch/batch-02-r08-r09-sandbox-facts.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/02-batch/ruling-2026-11-15-r3.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/04-reviews/blind-r1.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/04-reviews/final-r1.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/04-reviews/prereview-r1.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/04-reviews/prereview-r2.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/06-gates/gate-ci.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/06-gates/gate-index.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/06-gates/gate-selftest.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/06-gates/gate-tree.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/MANIFEST.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/SHA256SUMS.txt	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/freedom-list.md	repo
docs/validation/evidence/P75-R3-REAL-CLI-ASSESSMENT/input-list.md	repo
docs/validation/p57-closeout-assessment.md	repo
docs/validation/p57-closeout-assessment_EN.md	repo
docs/validation/r3-real-cli-assessment.md	repo
docs/validation/r3-real-cli-assessment_EN.md	repo
```
