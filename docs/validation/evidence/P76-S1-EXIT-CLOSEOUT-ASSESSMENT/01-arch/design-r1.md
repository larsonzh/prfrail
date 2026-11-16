# P76 S1 Exit 收口路线图 — 评估方案书 r1

日期：2026-11-16。切片：P76 S1-EXIT-CLOSEOUT-ASSESSMENT。级别：评估类 M 档。
基线：HEAD b90de21 / 准则 v1.69 / SELFTEST 737/737 / 三作用域 TOTAL_FAIL=0

**行号口径**：本方案书所有 `file:line` 取自 2026-11-16 磁盘实读（批量辅助三次 + 本阶段终端 `Select-String` 复核，见各节注）。引用以「文件名 + 节号」为主、行号为辅（同 `docs/validation/s1-reach-reachability.md:68` 声明）。

**本片边界（不越界）**：真实 enforcement 留 S2；不重跑 AT-23 / AT-24；不启真实 CLI 端到端探针；不改 `s1-reach-reachability{,_EN}`；不改准则（维持 v1.69）；②b 阶段才改 `s1-exit{,_EN}.md:52`（本方案书只给修正建议）。**`s1-exit.md:52` 未修正前不得宣告 S1 Exit 通过；S-G8 保持 blocked 至 S2。**

---

## 0. 执行摘要

- 现状总判断：P75 闭环后，S1 Exit 八门中 **S-G3 / S-G4 / S-G5 / S-G6 为 S1 内可推进项**；S-G1 / S-G8 根因是「无真实 enforcement」且已由 ADR-045 固定留 S2；S-G7 为合成门，随 S-G1 未通过而保持不可达。
- 门级现状（p57 §4.2，`:86`–`:95`）：S-G1 = 不可达/未通过（AT-23 阶段 1 条件闭合为永久事实，但判据「CLI 真实启动」未满足）；S-G2 = 口径已滞后（`s1-exit.md:44` 已更正，门级口径待裁定）；S-G3 = 未通过（ADR-009 绑定复验未做）；S-G4 = 未通过（ADR-010 `PROPOSED / NOT APPROVED`）；S-G5 = 未知（人工门，`s1-exit.md:60` = `BLOCKED`）；S-G6 = 不可达（门级与任务级不同源）；S-G7 = 不可达；S-G8 = blocked 至 S2（P75 A2 / ADR-045）。
- 推荐收口顺序：**① 准入修文（`s1-exit.md:52`，本片 ②b）→ ② S-G2 门级口径裁定（推荐判据字面口径）→ ③ S-G4（ADR-010 产品所有者批准）→ ④ S-G3（正式 Release + 最终 ZIP 绑定与复验）→ ⑤ S-G6 门级侧并陈推进 → ⑥ S-G5 分步批准 → ⑦ S-G1/S-G8/S-G7 留 S2 后重评**（详见 §12）。
- 最关键 3 个裁定点：**Q2（S-G1/S-G2 门级口径三选一）**、**Q3（`s1-exit.md:52` 修文采纳）**、**Q10（产品所有者批准时机与形式）**。
- Q3 判定结果：`s1-exit.md:52` 状态列 **`BLOCKED` 保持正确**；修文为**加法式澄清**（「AT-23 未通过」后补注「阶段 1 条件闭合已达成、阶段 2 不适用」），**不改历史判定、不改状态列**；对 S-G7 判定结论**无影响**（仍不可达）。
- Q7 历史状态核实结论：**从未执行**（正式独立安全评审结论/签署不存在；`SECURITY.md:5` 明写「未完成独立安全签署」）；该门禁不直接阻断 S-G5 字面判据，但须产品所有者在批准时一并裁定处置。
- 审查层：⑦ 盲审 **不抽**（不触必抽三语义 所有权/停机/身份，且新窗口 P76 不触发按比例抽查）；⑥ **可跳过**（§7.2.1 三条件成立），**但须 ⑦ 确认并记录理由**；若用户裁定改判八门状态 ⇒ ⑥ 必做（条件性风险见 §13.2）。
- 术语核实：P75 EN §11.2 无 `mapping`/`comparison` 混用（两词位于 §3/§12 且与 CN 成对一致）；v59 EN 全文一致使用 `closure assessment`。均**无需修正**（详见 §14）。
- 边界重申：S1 Exit 宣告的准入条件 = `s1-exit.md:52` 已修正 + 产品所有者批准 + 其余可推进门收口；S-G1/S-G8 未解锁前 S1 Exit 整体判定不可能为通过（S-G7 依赖 S-G1）。
- 越界防护：REQ-007 黑箱子项 / REQ-009 / REQ-020 三语言实跑等推 S2 项**不得**计入 S1 任何门判定（§11）；AT-13 自身仍在 S1（P71 注记明写）且 T017 已完成，不得误标为缺口。
- 以下各节每项判定均附 `file:line` 或逐字原文；凡属推断而非已实测事实的，在 §16 分界清单中显式列出。

---

## 1. Q1 八门与十门映射

依据：p57 §4.2 门表（`docs/validation/p57-closeout-assessment.md:86`–`:95`，行号 2026-11-16 磁盘实读，与 batch-01 一致）与 s1-reach §2 十门表（`docs/validation/s1-reach-reachability.md:22`–`:31`）。s1-reach 自身**未提供**八门↔十门映射（batch-02 缺失性记录：两文件内无映射原文；「八门」字样仅 `:14` 一处）。

| S1 Exit 八门（p57 §4.2） | s1-reach 十门 | 映射依据（file:line） | 判定差异 / 注意 |
|---|---|---|---|
| S-G1 AT-23 真实候选腿与整体（p57:88，判据「CLI 可真实启动经准入的 AI 候选并计费」） | **一对三**：G-A 真实候选腿（:22，不可达）+ G-B 机制腿（:23，已执行不独立成立）+ G-C AT-23 整体（:24，未通过） | p57:88 现状栏引 `s1-exit.md:43`；s1-reach :22–:24 把 AT-23 拆三条腿 | p57 用双词「不可达 / 未通过」；s1-reach 三腿分列。注意：p57 S-G1 现状栏锚点 `s1-exit.md:43` 现行内容已为「阶段 1 条件闭合达成」（P61/F2/P75 后），p57 判定句按 §11.1 不改，门级判定与「阶段 1 闭合」是两个维度（v59 附段 B） |
| S-G2 AT-24 与 T028（p57:89，判据「`supervised-black-box` 实现完成且 AT-24 判定成立」） | G-D T028 与 AT-24（:25） | p57:89 引 `DEV_PLAN.md:119` 与 `s1-exit.md:44`；s1-reach :25 证据同源 | p57 S-G2 现状栏写「**口径滞后**」（依据已过时：`:44` 已由 P64 更正为「条件闭合」，P75 §11.2 第 1 项关闭）；s1-reach G-D 状态 = 已执行不独立成立（解锁条件含「真实 CLI 端到端依赖 T027-B」，即 S-G8 语义折入） |
| S-G3 最终 ZIP 绑定与复验（p57:90） | G-E 最终 ZIP 绑定与复验（:26） | p57:90 引 `s1-exit.md:56`（`PARTIAL`）+ `ADR_REGISTER.md:19`（ADR-009）；s1-reach :26 同源 | 判定一致：均为「未通过」；依据完全相同 |
| S-G4 EOL 通知方式（p57:91） | G-F EOL 通知方式（:27） | p57:91 引 `ADR_REGISTER.md:20`（`PROPOSED / NOT APPROVED`）；s1-reach :27 同源 | 判定一致：均为「未通过」 |
| S-G5 产品所有者 S1 exit 批准（p57:92） | G-G 产品所有者 S1 exit 批准（:28，人工门） | p57:92 引 `s1-exit.md:60`；s1-reach :28 同源 | 判定词差异：p57/v59 = 「未知（人工门）」，`s1-exit.md:60` 状态列 = `BLOCKED`——同一门在人工门闭集词（未知）与门表状态词（BLOCKED）两种表述并存，语义一致（未批准且不可推进） |
| S-G6 T018 整体（p57:93） | G-H T018 整体（:29） | p57:93 引 `DEV_PLAN.md:94`（依赖字段）、`:96`（`IN_PROGRESS / BLOCKED`）；s1-reach :29 同源 | 判定一致：均为「不可达」；两者都显式写「门级与任务级**不同源**」（p57:93 / s1-reach :29 与 :33） |
| S-G7 S1 exit 总体（RFC §14）（p57:94） | G-I S1 exit 总体（:30） | p57:94 引 `s1-exit.md:52`（`BLOCKED`）、`:56`–`:58`（三项 `PARTIAL`）、`:60`；s1-reach :30 引 `s1-exit.md` 头部结论行（`BLOCKED`）与 RFC §14 表 | 判定一致：均为「不可达」；s1-reach G-I 判据字面更完整（「AT-01 至 AT-24 全过且 T018 完成且 RFC §14 全部 exit」） |
| S-G8 真实 CLI 端到端（p57:95，判据「受监督路径可在真实 CLI 与真实候选下端到端运行」） | **无对应门** | s1-reach 十门中不存在 S-G8；其语义折入 G-A/G-C 的解锁条件（真实 E2E 重跑）与 G-D 的解锁条件（「真实 CLI 端到端仍依赖 T027-B（D-18），需 P5 的 R-3 解锁」，:25） | S-G8 已由 P75 裁定保持 blocked 至 S2（p57:99 加法式指针 + ADR-045，`ADR_REGISTER.md:56`）；s1-reach 无独立表达，属覆盖缺口（事实，非矛盾） |
| （s1-reach 独有） | **OB-21**（:31，①② 两个子项是否仍成立） | 八门无对应行 | OB-21 是跨片观察项（「未通过（处置词保持待议）」），不属 S1 Exit 八门；不得并入任何 S-Gx 判定 |
| （s1-reach 独有） | G-B 机制腿单独成门 | 八门中 S-G1 聚合了机制腿 | s1-reach 把「机制腿已执行」单列（已执行不独立成立），八门无此粒度 |

**判定闭集差异（需在收口报告显式声明）**：八门侧（p57 §4.2 + `s1-exit` 门表）混用 `PASS / PARTIAL / BLOCKED / 不可达 / 未通过 / 未知`；十门侧闭集 = 「不可达 / 未通过 / 已执行不独立成立 / 未知」（s1-reach :18 逐字：「判定闭集 = 不可达 / 未通过 / 已执行不独立成立 / 未知。`未知` 仅用于人工门。」）。映射时不得跨闭集自动换词，逐门以各自定义源的原文为准。

---

## 2. Q2 S-G1 / S-G2 门级口径

张力来源（v59 §四）：S-G1 判据「CLI 可真实启动…」vs 状态「阶段 1 条件闭合达成；门级判定待重评」；S-G2 判据已满足 vs 状态「条件闭合；门级判定待重评」。三口径如下。

**口径 (a) 判据字面口径**：门级判定严格按 p57 §4.2 判据字面——S-G1 判据（真实启动经准入候选并计费）**不满足** ⇒ 不可达/未通过；S-G2 判据（T028 实现完成 **且** AT-24 判定成立）**已满足** ⇒ 通过（条件闭合口径）。AT 层成就（阶段 1 闭合 / 条件闭合）只作证据，不自动折算为门级通过；「真实 CLI 端到端」不属 S-G2 判据，归 S-G8。

- 现状：p57 §4.2 表头明写「判据（可判真假）」（:87）；(a) 下两门均可判——S-G1 = 假（R-3 = (A)：`r3-real-cli-assessment.md` §9 事实链「S1 内可达性：**不可达**…」，`s1-exit.md:43` 亦明写「真实候选腿…不构成真实候选兼容性证明」）；S-G2 = 真（`DEV_PLAN.md:119` T028 `COMPLETED` + `s1-exit.md:44`「条件闭合（2026-11-11，P64 验收；⑦ 复审 RE-REVIEW: PASS）」）。
- 影响：`s1-exit.md:52` 内「S-G2 口径（用户裁决③）…S-G2 门级判定待 v54 快照，不得据此声称该门通过」句须同步改判（v54 快照早已生成，该句已过时）；v59 §四 S-G2 行的「门级判定待重评」在新收口报告中改写为「通过（条件闭合口径）」。S-G1 保持不通过 ⇒ S-G7 仍不可达。

**口径 (b) 折中口径**：S-G1 = `PARTIAL`（阶段 1 条件闭合作为部分证据，真实候选腿留 S2）；S-G2 = 通过。影响：引入 p57 八门现状词中未用过的 `PARTIAL` 到 S-G1 行（p57 S-G1 现状词 = 「不可达 / 未通过」），且与 `s1-exit` 门表状态闭集混用（门表 52 行只有 BLOCKED/PASS/PARTIAL，而 PARTIAL 在门表中已用于 56/57/58 三行——若 52 也标 PARTIAL 会模糊「核心任务链阻断」与「局部成果未绑定」两类语义）。

**口径 (c) 挂起口径**：维持 v59 现状「门级判定待重评」直到 S-G8 解锁再统一判定。影响：S-G2 的既有成就无法落定，且 S-G8 解锁在 S2 ⇒ S1 Exit 收口路线图将失去「S1 内可推进项全部落定」的终点，违背本片目的。

**推荐：(a)**。理由：① p57 判据按「可判真假」设计，两门事实已足以判定，挂起没有新增信息；② (a) 与 s1-reach G-A/G-D 同向（G-A 不可达、G-D 已执行不独立成立），不产生跨文档新状态词；③ (b) 会把「核心任务链」的阻断强度从 BLOCKED 弱化为 PARTIAL，与 `s1-exit.md:52` 状态列、G-I 依据句（s1-reach :30「核心任务链 `BLOCKED`」）冲突；④ (c) 使收口无限期。**无论采纳何口径，S-G1 均不通过、S-G7 保持不可达、S-G8 保持 blocked**；差异只在 S-G2 是否标记通过及其表述落点。

---

## 3. Q3 `s1-exit.md:52` 正确口径判定

逐字原文（2026-11-16 磁盘实读，行 52 单行；状态列 = `BLOCKED`）：

CN：`| 核心任务链与 fail-close 场景 | BLOCKED | 既有 AT-01–AT-22 通过；T026 完成但固定候选为 10 verified、2 unsupported、0 unknown，整体不兼容；T027 未完成、AT-23 未通过；T028 已完成、AT-24 已验收 = 条件闭合（P64，2026-11-11；[P64 报告](at24-supervised-black-box.md)）；**依赖链**：AT-23 阶段 1 条件闭合达成（前提① F2 修复片满足、前提② P61 已满足；阶段 2 不适用，用户 2026-11-09 裁定 A1）⇒ **S-G1 依据 = AT-23 阶段 1 条件闭合**（…）；…**S-G2 口径（用户裁决③）**：AT-24 条件闭合（P64）⇒ v54 快照生成时可标记「AT-24 条件闭合可推进」；S-G2 门级判定待 v54 快照，不得据此声称该门通过 |`

EN：`| Core task chain and fail-closed scenarios | BLOCKED | Existing AT-01–AT-22 pass; T026 is complete, but its pinned candidate has ten verified, two unsupported, and zero unknown capabilities and is incompatible; T027 is incomplete and AT-23 is not passed, while T028 is complete and AT-24 is accepted = conditional closure (P64, 2026-11-11; [P64 report](at24-supervised-black-box_EN.md)); **dependency chain**: AT-23 stage-1 condition closure achieved (…) ⇒ **S-G1 basis = AT-23 stage-1 condition closure** (…); …**S-G2 caliber (user ruling ③)**: AT-24 conditionally closed (P64) ⇒ the v54 snapshot may mark "AT-24 condition closure can proceed"; the S-G2 gate-level verdict is pending the v54 snapshot and must not be cited as that gate passing |`

**判定：`BLOCKED`（状态列保持），行内口径不一致须修；修文性质 = 加法式，不改历史判定、不改状态列。** 三条依据：

1. **与本文件其它行一致**：门表三行 `PARTIAL`（:56/:57/:58）都是「有条件成果、最终载体未绑定」型局部门；:52 的核心任务链含 T027 未完成（任务级事实，`DEV_PLAN.md:94` 依赖字段含 T027、T027 状态 `BLOCKED / NOT IMPLEMENTED`）与 AT-23 未通过（`s1-exit.md:43`「真实候选腿…不构成真实候选兼容性证明」）⇒ 核心任务链真实阻断 ⇒ `BLOCKED` 正确；若改 `PARTIAL` 会把「核心链阻断」降格为「部分成果」，与 :7 头部结论「不得把部分通过标记为 S1 完成」的既有判词方向相反。
2. **与 p57 §4.3 原项一致**：原第 2 项指控「:52 称 T027/T028 仍未实现，T028 部分滞后」（p57 §4.3，batch-01 逐字）；T028 部分已由 P64 更正；残余 = 前半句「T027 未完成、AT-23 未通过」与同一单元格后半句「AT-23 阶段 1 条件闭合达成」的行内口径不一致（P75 §11.2 第 3 项 = 唯一残余，`r3-real-cli-assessment.md:165` 逐字：「该行前半仍写「T027 未完成、AT-23 未通过」，同一单元格后半写「AT-23 阶段 1 条件闭合达成」⇒ 行内口径不一致；修正时机 = S1 exit 收口片「重述门结论」时一并修正」）。修正对象 = 表述对齐，不是状态改判。
3. **与 s1-reach G-I 判定一致**：G-I = 不可达，依据含「RFC §14 表（核心任务链 `BLOCKED`、三项 `PARTIAL`、业主批准 `BLOCKED`）」（s1-reach :30 逐字）⇒ 若把 :52 状态改为 `PARTIAL`，将直接与 G-I 依据句冲突 ⇒ 必须保持 `BLOCKED`。

**可直接替换的修文（CN，加法式）**——把「T027 未完成、AT-23 未通过；」替换为：

`T027 未完成、AT-23 未通过（阶段 1 条件闭合已达成、阶段 2 不适用，2026-11-09 裁定 A1）；`

**可直接替换的修文（EN，加法式）**——把「T027 is incomplete and AT-23 is not passed, while」替换为：

`T027 is incomplete and AT-23 is not passed (stage-1 condition closure achieved; stage 2 not applicable, ruling A1, 2026-11-09), while`

说明：仅在前半句补注「阶段 1 条件闭合已达成」并点明「阶段 2 不适用」的时点裁定，使前后半句口径对齐；不动状态列 `BLOCKED`、不动门表 9 行结构（P71 注记 `s1-exit.md:62`「门表 9 行的行名、行数与状态判定全部不变」）、不改任何历史判定（P61 PARTIAL / P63 豁免原文保留）。**是否影响 S-G7：否**——S-G7 依据栏只引 :52 的状态词 `BLOCKED`，修文不改状态词 ⇒ S-G7 判定结论不变（仍不可达，见 §4）。

**附：Q2 裁定后的条件修文（第二处，仅当 Q2 采纳口径 (a) 时执行）**——CN 把「**S-G2 口径（用户裁决③）**：AT-24 条件闭合（P64）⇒ v54 快照生成时可标记「AT-24 条件闭合可推进」；S-G2 门级判定待 v54 快照，不得据此声称该门通过」替换为「**S-G2 门级口径（P76 Q2 裁定 (a)）**：判据「`supervised-black-box` 实现完成且 AT-24 判定成立」已满足（T028 `COMPLETED`（`DEV_PLAN.md:119`）+ AT-24 条件闭合（P64））⇒ S-G2 = 通过；真实 CLI 端到端归 S-G8（blocked 至 S2），不并入 S-G2」；EN 对应替换为「**S-G2 gate caliber (P76 Q2 ruling (a))**: the criterion "the `supervised-black-box` mode is implemented and the AT-24 verdict holds" is met (T028 `COMPLETED` (`DEV_PLAN.md:119`) + AT-24 conditionally closed (P64)) ⇒ S-G2 = passed; the real CLI end-to-end path belongs to S-G8 (blocked until S2) and is not folded into S-G2」。该处为条件修文，须用户 Q2 裁定后由 ②b 执行。

---

## 4. Q4 S-G7 依据栏同步更新

- **结论：S-G7 判定结论不变**（仍「不可达」）。理由：S-G7 判据 = 「全 AT 通过且 T018 完成且 §14 全 exit」（p57:94），依赖 S-G1 至 S-G6；S-G1 未通过 + S-G8 blocked（S-G1 同源解锁，p57:95）+ S-G5 未批准 ⇒ 合成门不可达，与 Q3 修文无关。
- **p57 §4.2 S-G7 行依据栏无需改**：现状句「不可达：`s1-exit.md:52`（核心任务链 `BLOCKED`）、`:56` 至 `:58`（三项 `PARTIAL`）、`:60`」（p57:94）只引状态词，Q3 修文不改状态词 ⇒ 锚点依然成立。
- **若需在 v59 §四 或新收口报告同步描述，建议新表述**（加法式，不改 p57 本体除非用户另裁）：「不可达：`s1-exit.md:52`（核心任务链 `BLOCKED`；行内口径已按 P75 §11.2 第 3 项修正为「AT-23 未通过（阶段 1 条件闭合达成、阶段 2 不适用）」）、`:56`–`:58`（三项 `PARTIAL`）、`:60`（`BLOCKED`）」。

---

## 5. Q5 S-G3 解锁路径（最终 ZIP 绑定发行记录 + 复验）

依据链：`INSTALLATION.md:5`「当前尚无正式 GitHub Release；新目录结构必须通过新的可信候选 artifact 复验后，才能视为与最终发行物结构相同」；`S1_RELEASE_NOTES.md:5`「状态：`DRAFT / NOT RELEASED`。版本、下载链接、正式 artifact 文件名、支持起始日期和批准人均待 S1 exit 冻结；本文件不授权 tag、GitHub Release、发布或安装」；`INSTALLATION_PLAN.md:50`「正式发行文件名和下载地址仍须在发行批准时绑定」；p57:90 依赖 = 「S-G2、S-G4、发行记录」。分四步：

| 步 | 动作 | 前置条件 | 责任主体 | 证据形态 |
|---|---|---|---|---|
| ① | 产品所有者批准正式 Release（含版本号 v0.1.0 与发行授权） | S-G4 已批准（ADR-010 含版本号方案）或单独批准版本号 | 产品所有者 | ADR-010 状态更新（`ADR_REGISTER.md:20`）+ 发行授权记录（版本/日期/主体；ADR 表头要求「记录主体、日期…禁止模型代签」） |
| ② | 构建最终 ZIP artifact | ① 完成；CI candidate-build 作业产出候选 artifact（T017 已冻结发布证据链：固定 commit、SHA256SUMS、SBOM、许可证清单） | 主控 + CI | GitHub Actions artifact + `SHA256SUMS` + `sbom.spdx.json` + `licenses.json`（`INSTALLATION.md:152` 列明四文件） |
| ③ | 绑定发行记录：最终版本 / commit / run / 文件名写入三处 | ② 完成 | 文档角色 | `INSTALLATION.md:206`「正式发行前，以上命令须用最终 ZIP 再执行一次，并在 [S1 exit 报告]、[支持矩阵] 和 [发行说明] 中绑定最终版本、commit、run 和文件名」的落盘行；`S1_SUPPORT_MATRIX.md:5` 的 `CANDIDATE` 状态同步解除；`S1_RELEASE_NOTES.md:5` 的 `DRAFT / NOT RELEASED` 状态同步 |
| ④ | 安装生命周期复验：用最终 ZIP 重跑安装 / 升级 / 回滚 / 卸载 | ③ 完成 | 测试角色 | `INSTALLATION_PLAN.md:52`–`:59` 发布前安装门禁逐条证据 + `s1-install-drill` 复验报告（:38「正式下载入口和最终 ZIP 尚未发布，仍须用最终载体重跑」） |

要点：S-G3 依赖 S-G4（版本号在 ADR-010 内）⇒ S-G4 必须先于或并行于 S-G3（并行仅限「审批动作」与「ZIP 构建准备」并行，绑定/复验必须等版本号冻结）；S-G3 收口后门表 :56（正式安装包教程）与 :57（安装/升级/回滚/卸载实测）可同步转 `PASS`（该转判属于门级判定落地，建议列入 Q15 裁定或随 ②b 第二层文档同步，不属本方案书自裁）。

---

## 6. Q6 S-G4 解锁路径（ADR-010）

- **批准执行方式**：产品所有者明确批准。ADR 表头（`ADR_REGISTER.md` 头部）逐字：「每项批准时记录主体、日期、RFC 变更和测试证据，禁止模型代签」⇒ 执行者只能准备批准材料，批准动作与记录主体必须是产品所有者。ADR-010 行（:20）当前状态逐字：「`PROPOSED / NOT APPROVED`；90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准」。
- **文档改动点（批准后）**：
  1. `ADR_REGISTER.md:20` 状态列 → 批准记录（主体 / 日期 / 依据）；
  2. `INSTALLATION_PLAN.md:32`「state/store 默认位置、正式下载入口与 EOL 通知仍待冻结」中 EOL 部分 → 已冻结；
  3. `s1-exit.md:66`「EOL 通知与例外已有 ADR-010 待审批方案，但尚未获产品所有者批准」→ 已批准；
  4. `PRODUCT_REQUIREMENTS.md:107`「精确依赖、EOL 通知、安全评审者和具体 whois 输入仍待后续任务落实」中 EOL 通知项 → 已落实（该行其余项不动）；
  5. （视第二层清单）`S1_SUPPORT_MATRIX{,_EN}` 尾部支持策略段同步。
- **批准前 / 后差异**：
  - 前：S-G4 = 未通过（p57:91）；首版版本号未定 ⇒ S-G3 无法绑定；EOL 政策未生效（非已发布 SLA，`S1_SUPPORT_MATRIX.md:5`）。
  - 后：S-G4 = 通过；版本号 v0.1.0 冻结 ⇒ S-G3 进入绑定阶段；EOL 公告义务（新 minor 发布时双处公告，ADR-010 方案）自首版起可执行。
- 注意：ADR-010 的批准**不等于** S-G5（S1 exit 整体批准）；两者是不同授权对象（发行策略 vs 阶段退出）。

---

## 7. Q7 S-G5 与独立安全评审的关系

**历史状态核实结论：从未执行（作为正式的「独立安全评审」结论 / 签署）。** 检索事实与依据：

1. `docs/SECURITY.md:5`（2026-11-16 磁盘实读，行号实测）：「日期：2026-09-07；安全设计基线，未完成独立安全签署。」——**明写未完成签署**。
2. `docs/PRODUCT_REQUIREMENTS.md:107`：「精确依赖、EOL 通知、**安全评审者**和具体 whois 输入**仍待后续任务落实**」（评审者人选都未定）；`:109`：「独立安全评审结论仍是**进入 S1 前**的阻断门禁。」（EN 对应 `:105`/`:107`：「security reviewer… remain later-task details」「Independent security review remains a blocking gate before S1」）。
3. `docs/DEV_PLAN.md:27`（T001 行）：「精确依赖及独立安全评审仍按后续门禁落实」。
4. `docs/validation/` 全域检索（batch-03 §5，2026-11-16）：`security review` 唯一命中 = P76 本片 input-list 自引；`independent security` 唯一命中 = strategic-review EN:49 的条目句；`安全评审` 8 命中均非评审报告；`independent review` 49 命中 / 26+ 文件全部是各切片 ⑤/⑥/⑦ 审查记录与 T017 双主审，**无一为「独立安全评审报告」**。
5. 最接近的局部覆盖（**不构成** §7 意义的独立安全评审）：`s1-exit.md:59`「独立评审 | PASS | T017 基线及 Node 24 替换均经 Codex/DeepSeek 双主审」（发布信任面审查）；P52 `acceptance-security-counterexamples.md`（fail-closed 反例片，⑥/⑦ PASS，但发现并登记 GAP-77 HIGH——侧面证明未达到「零发现签署」态）。
6. `docs/TEST_STRATEGY{,_EN}.md` 全文检索「安全评审 / 独立评审 / security review」= 0 命中（本阶段 grep 实测）。
7. `docs/validation/strategic-review-2026-10-30.md:49`（CN，行号实测）：「| 2 | 在同一接受链上做独立安全反例验收 | `ACCEPTANCE-TRUST-NEGATIVES` |」——该条目已由 P52 执行（反例验收 ≠ 独立安全评审）。

**四种状态下对 S-G5 的影响**：已通过 ⇒ S-G5 无额外前置；已通过但须重新验证 ⇒ S-G5 附条件；未通过 ⇒ S-G5 被该门禁阻断；**从未执行（实况）** ⇒ S-G5 判据字面（p57:92「产品所有者作出批准」）不含独立安全评审前置，不直接阻断；但 §7 的「进入 S1 前阻断门禁」属 S0→S1 门，历史上 S1 已由产品所有者批准进入而该门禁未补闭环 ⇒ 构成**未闭环承诺项**。

**是否阻断 S1 Exit 宣告**：不构成八门判据的直接阻断（八门唯一人工门 = S-G5，判据只有「产品所有者作出批准」）；但建议作为 **S-G5 批准材料的附件裁定项**——产品所有者在批准时明确三选一：① 补一次独立安全评审（新切片）；② 正式豁免 / 顺延至 S2 并登记 ADR；③ 其它处置。未裁定前，收口报告须如实标注「独立安全评审从未执行，§7 门禁状态待产品所有者裁定」。（`CHECKLIST_v59.md:103` S-G5 行已预登记此待澄清关系。）

---

## 8. Q8 S-G6 门级与任务级不同源

- **DEV_PLAN 实测字段**（2026-11-16 磁盘实读，行号实测）：`:94` T018 任务行依赖字段 = 「依赖：T017、T025、T026、T027、T028 及 whois 输入授权」；`:96`「T018 状态：`IN_PROGRESS / BLOCKED`」；`:119` T028 = `COMPLETED`；T027 = `BLOCKED / NOT IMPLEMENTED`（防御性重构期）。
- **三方案**：
  - (a) 门级口径：四道剩余门（S-G3/S-G4/S-G5 及 S-G2 门级）清空即判 S-G6 门级通过，T018 任务级另行。
  - (b) 任务级口径：T018 勾选必须等 T027/T028 完成，S-G6 随任务级不可达。
  - (c) 门级与任务级并陈：门级进展与任务级勾选分开表述，不给单一结论。
- **影响**：(a) 会把 T018 提前标记，与 `DEV_PLAN.md:96` 冲突且需改任务状态字段（超本片边界，禁止）；(b) 丢失门级进展表达，与 p57 S-G6 行既有句冲突；(c) 与 s1-reach :33（「门级与任务级不同源（重要）…四门全部清空后 T018 **仍不可勾选**。本报告按门级与任务级**分开**表述，不给单一结论」）、p57:93、v59 §四 S-G6 行完全一致。
- **推荐：(c)**（现状维持，无新增改动）。S1 内后果：S-G3/S-G4/S-G5 收口只推进门级侧；T018 任务级因 T027（`BLOCKED / NOT IMPLEMENTED`）与 S-G1 未通过而**不可勾选**，S-G6 门级判定保持「不可达」直至 S2 解锁 S-G1/S-G8 并完成 T027。S-G6 的收口形态 = 门级侧记进展、任务级侧挂 S2，两者在收口报告中分列。

---

## 9. Q9 S1_SUPPORT_MATRIX 口径滞后修正

现状（P75 后口径滞后）：`S1_SUPPORT_MATRIX.md:13` 句尾「T026/T027/AT-23 阻断 S1」；`:14` 句尾「AT-24 条件闭合，S1 仍因 T027/AT-23 阻断」。P75 后事实：R-3 = (A)、T027-B = 完成（部分不适用）、S-G8 = blocked 至 S2、AT-23 阶段 1 条件闭合 = 永久事实（ADR-045，`ADR_REGISTER.md:56`；`DEV_PLAN.md:203`）。

**精确修正文字（CN）**：
- `:13` 句尾「T026/T027/AT-23 阻断 S1」→「T026 已完成（候选不兼容）；T027 未完成；AT-23 阶段 1 条件闭合达成、真实候选腿待 S2 enforcement 接入后重跑（ADR-045）；S1 内该通道不可运行发布」。
- `:14` 句尾「AT-24 条件闭合，S1 仍因 T027/AT-23 阻断」→「AT-24 条件闭合；真实 CLI 端到端留 S2（S-G8 blocked 至 S2，ADR-045）」。
- `:5` 保持（`CANDIDATE`，待正式发行载体和 S1 exit 批准后生效——该句与 S-G3/S-G5 路径一致，不滞后）。

**精确修正文字（EN）**：
- `:13`「T026/T027/AT-23 block S1」→「T026 is complete (candidate incompatible); T027 is incomplete; AT-23 stage-1 condition closure is achieved and the real-candidate leg re-runs only after S2 enforcement lands (ADR-045); this channel is not a runnable release inside S1」。
- `:14`「AT-24 conditionally closed, S1 still blocked by T027/AT-23」→「AT-24 conditionally closed; the real CLI end-to-end path is deferred to S2 (S-G8 blocked until S2, ADR-045)」。

**是否纳入 ②b**：**建议纳入**。理由：口径同步不改任何门级判定、不含语义变更；②b 第二层文档同步清单已含 `S1_SUPPORT_MATRIX`（briefing【八】）；若不纳入，将残留「T027/AT-23 阻断 S1」与 ADR-045 的矛盾口径，且 S-G3 ③ 步绑定支持矩阵时仍需改此处，届时两次改动不如一次到位。若用户裁定不纳入 ⇒ 收口报告登记为待办。

---

## 10. Q10 产品所有者批准时机与形式

| 选项 | 内容 | 风险 | 后果 |
|---|---|---|---|
| (A) 先行分段批准 | 现在即批准发行类事项（ADR-010 + 正式 Release），S-G5 标记「部分批准（发行部分）」 | 批准语义可能被扩大解释为「S1 Exit 通过」，需严格限定措辞 | 解除 S-G3 死锁，但引入「部分批准」新状态词 |
| (B) 一次性终批 | S-G3/S-G4 收口 + S2 完成后统一批准 S-G5 | 死锁：S-G3 需正式 Release，Release 需产品所有者批准，S-G5 需 S-G1（S2 才解锁）⇒ 首版发布无限期延迟，与 ADR-009「最终 ZIP 仍须绑定发行记录并复验」的 S1 定位冲突 | S1 内无法产生任何发行物 |
| (C) 分两步批准（推荐） | 第一步（S1 内）：批准 ADR-010（版本号 + EOL）+ 批准 v0.1.0 正式 Release 创建（含最终 ZIP 绑定）⇒ S-G3/S-G4 收口，S-G5 记「发行部分已批准，S1 Exit 整体批准待 S2」；第二步（S2 后）：批准 S1 Exit 整体 | 需两次授权动作；S-G5 的「未知」状态在 S1 内如实保持 | 首版发布可落地，批准语义不扩大，S-G5 边界清晰 |

**建议：(C)**。理由：既满足「S1 内可推进项收口」的本片目标，又保持「S1 Exit 整体通过」仅在 S2 解锁 S-G1/S-G8 后由产品所有者最终批准；与 §7 的独立安全评审附件裁定可合并在第一步批准材料中。**硬边界**：任何批准动作不得被解读为「S1 Exit 通过」；`s1-exit.md:52` 未修正前不得宣告；S-G8 保持 blocked。

---

## 11. Q11 REQ 阶段归属边界（不得算进 S1 判定的项）

| 项 | 阶段归属（原文） | 依据 file:line | 对 S-G1–S-G8 的约束 |
|---|---|---|---|
| REQ-007 黑箱候选子项 | 推迟 S2 | `PRODUCT_REQUIREMENTS.md:69`（R-1「S1/S2（黑箱候选子项推迟 S2）」）；`s1-exit.md:62`（P71 注记「黑箱候选子项（REQ-007/021）标推迟 S2」） | 不得以「黑箱子项未完成」判任一门失败或通过；S-G2 的门级判定不得引入黑箱子项口径（其判据只含 T028 + AT-24） |
| REQ-009 完整一体化 TUI | S1 不要求，推迟 S2 | `PRODUCT_REQUIREMENTS.md:71`（R-2）；`s1-exit.md:62`（「完整 TUI 移出不构成门级判定变更」） | 不得作为 S-G7 或任何门的缺口；门表 9 行判定不受 TUI 移出影响（P71 注记） |
| REQ-020 三语言实跑 | S1 仅 C+Go 契约；S2 三语言实跑 | `PRODUCT_REQUIREMENTS.md:82`（R-4「S1 不要求（仅 C+Go 契约）；S2 承诺收窄…S2 三语言实跑，失败整任务阻断」） | 不得算进 S-G6 / T018 的判定依据 |
| AT-13 相关 | **AT-13（自托管）仍在 S1**，T017 已完成 | `s1-exit.md:62`（「自托管（REQ-019 / AT-13）仍在 S1」）；`DEV_PLAN.md` T017 行 = 已完成 | 不得把「AT-13 未完成」当作 S1 缺口（它已完成）；与 AT-13 捆绑的推 S2 表述不得反向夹带回 S1 判定 |
| 附加：MVA 验收复合标准 | 不批正式定义 | `s1-exit.md:62` | 不得引用 MVA 作为 S1 门级判据 |

通用约束：任何 S-G1–S-G8 的判定依据都不得引用上述 S2 项作为「未通过原因」；S-G1/S-G8 的 blocked 依据只能来自 R-3 = (A)、T027 装配未完成与候选兼容性（ADR-045 固定），不得夹带 REQ-009 / REQ-020。

---

## 12. 收口顺序整体建议

推荐序列（依赖排序，含并行可能）：

1. **准入修文**：`s1-exit{,_EN}.md:52` 按 §3 修文（本片 ②b；依赖用户 Q3 裁定 + ② 阶段授权）。
2. **S-G2 门级落定**：按 Q2 裁定（推荐 (a)）更新 `s1-exit.md:52` 的 S-G2 句（§3 附条件修文）+ 新收口报告 S-G2 表述。
3. **S-G4**：ADR-010 产品所有者批准 → `ADR_REGISTER.md:20` 及 §6 所列文档改动点。
4. **S-G3**：正式 Release 创建 → 最终 ZIP 构建 → 绑定（s1-exit / 支持矩阵 / 发行说明）→ 生命周期复验（§5 四步）。
5. **S-G6 门级侧**：随 S-G3/S-G4/S-G5 收口在收口报告中并陈推进；任务级保持不可勾选。
6. **S-G5**：按 Q10 裁定（推荐 (C)）分步批准。
7. **S-G1 / S-G8 / S-G7**：留 S2（ADR-045）；S1 Exit 整体宣告待 S2 解锁后由产品所有者最终批准。

并行可能：S-G4 审批动作与 S-G3 的 ZIP 构建准备（CI 候选 artifact）可并行；S-G3 的绑定/复验必须等版本号冻结（S-G4 完成）。每步依赖见上表。**预计阻塞点**：产品所有者批准动作（ADR-010 / Release）、S2 enforcement 接入（S-G1/S-G8）——后者不在 S1 控制内。

---

## 13. 审查层判定

### 13.1 ⑦ 盲审是否必抽

- **必抽条件**（`DELIVERY_DIRECTIVE` §7.7「必抽盲审」）：切片触及**所有权 / 停机 / 身份**语义。本片变更面 = `tmp/` 方案书 + `s1-exit{,_EN}:52` 修文 + 新收口报告 + 台账行，全部为文档措辞；不触碰进程所有权、停机原语、身份/授权语义 ⇒ **不触发必抽**。
- **按比例抽查**（§7.7「每 4 片抽 1 片」+ briefing【五】逐字：「新窗口自 P76 起算（P76/P77/P78 不触发按比例抽查，P79 触发）」）⇒ **P76 不触发**。
- **结论：⑦ 盲审不抽**（两条件均不触发）；⑦ 终审按 M 档 1 终审 + 2 复审（§7.5 表）执行，总额度上限 4（briefing【五】；R≥2 才升 §9.2 高风险 1+3）。若用户裁定引入门级状态改判，仍不触必抽三语义，但 ⑦ 终审输入包须附用户裁定原文。

### 13.2 ⑥ 是否可跳过（按 §7.2.1 三条件逐条）

- **① 改动只落文档与配置文案类文件**：本片全部改动为 `.md`（`tmp/` 方案书、`s1-exit{,_EN}:52` 修文、新收口报告、`ADR_REGISTER`/`S1_SUPPORT_MATRIX` 若按裁定同步）⇒ **成立**。
- **② 不含 `.go`/`.json`/`.yml`/schema/fixtures/脚本逻辑的语义变更**：本片不改任何代码/schema/fixtures/脚本 ⇒ **成立**。
- **③ 不涉及状态机、门禁判据、证据模型、角色/授权契约的语义**：`s1-exit.md:52` 修文只改门表行内表述，状态列 `BLOCKED` 不变、门表 9 行行名/行数/状态判定不变（P71 注记先例，`s1-exit.md:62`）、p57 §4.2 八门定义与判定不变（§11.1）、准则 v1.69 不动 ⇒ **成立**（须在报告中写明该边界论证）。
- **结论：⑥ 可跳过，但跳过必须由 ⑦ 确认并在报告中记录理由**（§7.2.1「跳过必须由 ⑦ 或用户确认（主控不得单独认定“可跳过”）」）。
- **条件性风险（显式标注）**：若用户在停点裁定中选择「改判八门状态」（如 S-G2 → 通过 落地到 p57 §4.2 状态栏、或 S-G3 相关门表行转 PASS）或其它触及门禁判定语义的选项 ⇒ ③ 被推翻 ⇒ **⑥ 必做**（SW-1/SW-2 反面先例：表面文档切片修订门禁判据语义，不得跳过，§7.2.1 反面样例）。本方案书推荐的所有路径均不引入状态改判落地 p57，故 ⑥ 跳过的前提是「用户裁定与推荐一致或至少不引入改判」。

---

## 14. 术语一致性核实结论

### 14.1 P75 EN §11.2 中 `mapping` vs `comparison`

- 实测（2026-11-16 磁盘读取 + 全文 grep）：`docs/validation/r3-real-cli-assessment_EN.md` §11.2（「O-7 obligation hand-off」表）内**既不出现 `mapping` 也不出现 `comparison`**。
- `comparison` 出现在 EN §3 标题「R-3 three-way comparison」（EN:51），CN 对应 §3 = 「R-3 三选一对比」（CN:51，grep 实测）；`mapping` 出现在 EN §12 区「the wording mapping between R-08.5 “conditional closure” and the S-G8 criterion “real end-to-end run” is not yet written」（EN:177），CN 对应 = 「R-08.5 的「条件闭合」与 S-G8 判据「真实端到端运行」之间的口径映射尚未成文」（CN:177，grep 实测）。
- **结论：无需修正**。两词分属不同概念（comparison = R-3 选项比较；mapping = R-08.5 与 S-G8 的口径映射），且 CN/EN 成对一致（comparison↔对比、mapping↔映射），不存在 §11.2 内混用；若此前观察认为 §11.2 有两词混用，系 §3/§12 的合法差异被误归位。

### 14.2 v59 EN 内部 `closure assessment` vs `closeout assessment`

- 实测（grep 全文，8 命中）：`docs/t027/CHECKLIST_v59_EN.md` 统一使用 `closure assessment`（行 1/103/110/150/176/180 等）；`closeout` 仅出现于文件名引用 `p57-closeout-assessment_EN.md`（行 36/38），不作为术语使用。CN 侧统一「收口评估」（无「收尾评估/闭口评估」变体，grep 实测）。
- **结论：无需修正**。v59 EN 内部术语一致。切片 ID `S1-EXIT-CLOSEOUT-ASSESSMENT` 与 ②c 将新建的报告文件名（`s1-exit-closeout-assessment{,_EN}.md`）属**切片/产物命名层**，与 v59 正文的 `closure assessment` 属两层命名，不构成同文混用；②c 新报告正文引用 v59 时按 v59 原文 `closure assessment` 即可，不强制统一。

---

## 15. 需用户裁定的选项汇总

| Q# | 选项 | 后果 | 主控推荐 | 一句话理由 |
|---|---|---|---|---|
| Q2 | (a) 判据字面 / (b) S-G1=PARTIAL 折中 / (c) 挂起待重评 | 决定 S-G2 是否可标记通过、`:52` 后半句是否二次修文 | (a) | 判据「可判真假」已可判，且与 s1-reach G-A/G-D 同向 |
| Q3 | A 采纳加法式修文（BLOCKED 保持）/ B 改判 PARTIAL | `:52` 状态列是否变化、G-I 依据一致性 | A | 与 s1-reach :30「核心任务链 BLOCKED」一致，不改历史判定 |
| Q7 | 独立安全评审：补评 / 豁免顺延 S2 / 其它 | 未闭环承诺项如何处置、S-G5 批准材料附件 | 不预设（由产品所有者裁定）；补评或顺延二选一 | §7「进入 S1 前阻断门禁」从未执行，属产品所有者待裁项 |
| Q9 | 支持矩阵口径修正纳入 ②b / 不纳入 | 是否消除「T027/AT-23 阻断 S1」滞后口径 | 纳入 | 口径同步不改判定，避免与 ADR-045 残留矛盾 |
| Q10 | (C) 分两步批准 / (A) 先行分段 / (B) 一次终批 | S-G5 与发行批准的关系、首版发布时点 | (C) | 解除 S-G3 死锁且不扩大批准语义 |
| Q12 | 收口顺序认可 / 调整 | 后续执行片派单顺序 | 按 §12 序列 | 依赖关系已闭环，无平行替代更优 |

（Q1/Q4/Q8/Q11/Q13/Q14 为判定或事实性结论，不构成决策分叉，不在上表。）

---

## 16. 事实与推断分界

**已实测事实（2026-11-16 磁盘实读或批量辅助终端读取）**：

- HEAD `b90de21`、准则 v1.69、SELFTEST 737/737、三作用域 `TOTAL_FAIL=0`（P76 briefing 步骤 0，主控实测）。
- `s1-exit.md:52` 状态列 = `BLOCKED`，全文「P75」0 命中、「尚未实现」0 命中；`:43`/`:44`/`:56`–`:58`/`:60` 现状逐字（batch-01 附 A/B + §5.1）。
- p57 §4.2 八门行（`:86`–`:95`）、§3.3 三指针（`:72/:74/:76`）、S-G8 指针（`:99`）逐字；§4.3 原第 2 项 = `s1-exit.md:52`。
- P75 §11.2 第 3 项 = `s1-exit.md:52` 唯一残余（`r3-real-cli-assessment.md:165` / EN 对应段）。
- s1-reach 十门（`:22`–`:31`）、判定闭集（`:18`）、门级/任务级不同源（`:33`）、G-I 依据句（`:30`）。
- `ADR_REGISTER.md`：ADR-009（:19）/ ADR-010 `PROPOSED / NOT APPROVED`（:20，全文唯一命中）/ ADR-012（:22）/ ADR-045（:56）。
- `DEV_PLAN.md`：T018 依赖字段（:94）、状态（:96）、T028 `COMPLETED`（:119）、T027 `BLOCKED / NOT IMPLEMENTED`、T027-B 收口记录（:203）。
- `INSTALLATION.md:5/:152/:206`、`INSTALLATION_PLAN.md:32/:38/:50/:52`–`:59`、`S1_RELEASE_NOTES.md:5`（`DRAFT / NOT RELEASED`）、`S1_SUPPORT_MATRIX.md:5/:9/:13/:14` 逐字。
- `SECURITY.md:5`「未完成独立安全签署」、`:64`「以下事项需要独立评审…」；`PRODUCT_REQUIREMENTS.md:69/:71/:79/:82/:107/:109` 逐字。
- `docs/validation/` 独立安全评审检索 0 命中评审报告（batch-03 §5）；`TEST_STRATEGY{,_EN}` 0 命中（本阶段 grep）。
- `CHECKLIST_v59_EN.md` 用 `closure assessment`（8 命中）、CN 用「收口评估」；P75 EN §11.2 无 `mapping`/`comparison`。

**推断 / 待核（非实测事实，标注为判断）**：

- Q7 定性「从未执行」：事实基础 = 无签署记录 + `SECURITY.md:5` 明写；但「产品所有者历史上是否曾口头批准或豁免该门禁」为**待核**（未见记录 ≠ 确定没有，须产品所有者确认）。
- Q2 口径 (a) 的影响面（`:52` S-G2 句更新、新报告表述）依赖用户裁定，属规划推断。
- S-G3 四步中「正式 Release 创建由谁执行、CI artifact 是否可直接转正」为**待核**（T017 证据链覆盖候选 artifact，最终 ZIP 的构建动作需执行片落实）。
- Q10 (C) 的「第一步批准材料」构成（是否并入 §7 附件裁定）为建议性推断。
- ⑥ 可跳过的前提 = 用户裁定不引入门级状态改判（§13.2 条件性风险），该前提在停点前属待定。

---

## 17. 主控抽验附注（2026-11-16，主控 DeepSeek V4.1 Flash 实读核验）

独立抽验 6 处 `file:line`，与方案书正文一致：

- `docs/validation/s1-exit.md:52` 状态列 = `BLOCKED`，行内含「T027 未完成、AT-23 未通过」——一致。
- `docs/DEV_PLAN.md:94` 依赖字段全文 = 「T017、T025、T026、T027、T028 及 whois 输入授权」——一致。
- `docs/DEV_PLAN.md:96` = `T018 状态：IN_PROGRESS / BLOCKED`——一致。
- `docs/SECURITY.md:5` = 「安全设计基线，未完成独立安全签署」——一致。
- `docs/S1_SUPPORT_MATRIX.md:13` 句尾 = 「T026/T027/AT-23 阻断 S1」——一致。
- `docs/S1_SUPPORT_MATRIX.md:14` 句尾 = 「AT-24 条件闭合，S1 仍因 T027/AT-23 阻断」——一致。

编号瑕疵登记（不影响结论，② 阶段按 Q1–Q11 唯一体系执行）：

- §15 汇总表出现 `Q12`（收口顺序）行，并括注 `Q13/Q14`；本片问题集**只到 Q11**。读法 = §12 的「收口顺序整体建议」是否被用户认可。
- §5 第 ④ 步括注「（该转判属于门级判定落地，建议列入 Q15 裁定…）」为编号笔误，应读作「Q10 裁定 + ②b 第二层文档同步」。