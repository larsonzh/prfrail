# P77 方案书（design-r1）

切片：`ADR010-APPROVAL`（S 档文档类 `[SLICE]`）。本方案书所有锚点行号与原文均经磁盘字节级复核；复核中发现起跑消息事实 C/D 所述「软换行空格」在磁盘上不存在（详见 §Q2 与「风险与开放项」第 1 条），落盘一律以磁盘实际字节为准。

---

## Q1：ADR-010 批准后状态栏的精确文本

### 1.1 骨架裁定：如何同时满足「首状态词 = APPROVED」与「四要素按 ADR-014 格式（分号分隔）」

裁定采用**形态 A 骨架**（与 ADR-014 / ADR-045 / ADR-046 / ADR-047 同款结构）：

`**状态词（括注）**。四要素：裁决人 = …；日期 = …；理由 = …；依据 = …`

- 「首状态词 = APPROVED」由**状态词位**承载：单元格第一个加粗词取 `APPROVED`。
- 「四要素按 ADR-014 格式、分号分隔」由**四要素句**承载：四要素句结构与 ADR-014 完全一致（裁决人、日期、理由、依据四项，项内 `=`、项间 `；`）。
- 两条要求在骨架中分处不同位置，**不冲突**：起跑消息要求的是四要素句的结构，并非要求状态词本身必须写成「已裁/已登记」。
- 括注取 ADR-045 同款「用户直接裁定，属决策记录」，**不写**「不得作为先例」（那是临时授权例外的标注，决策记录不适用；ADR-045 先例即如此）。
- 历史完成条件原文（分号后一段）**逐字保留**，尾部追加闭括号。

### 1.2 CN 目标文本（`docs/ADR_REGISTER.md:20` 第 4 列整格，逐字）

改前：

`\`PROPOSED / NOT APPROVED\`；90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准`

改后：

`**APPROVED（用户直接裁定，属决策记录）**。四要素：裁决人 = 用户；日期 = 2026-11-16；理由 = v0.1.0 进入 Release 流程前冻结 EOL 通知方式；依据 = 用户 2026-11-16 P77 起跑消息同轮裁定原文。90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准（已于 2026-11-16 落实；P77）。**性质**：决策记录（同 ADR-045 / ADR-046 / ADR-047），不是临时授权例外（§1.7.4 五要素不适用）。**可引用**：作为 S-G4「ADR-010 获产品所有者批准」判据满足与 EOL 通知方式已冻结的决策依据。**生效日**：2026-11-16。**回退**：本 ADR 不追溯改写；未来修订 EOL 通知方式须走既有例外规则（带日期、理由和到期日的 ADR，经产品所有者明确批准后同步两处）。`

### 1.3 EN 目标文本（`docs/ADR_REGISTER_EN.md:20` 第 4 列整格，逐字）

改前：

`\`PROPOSED / NOT APPROVED\`; 90 days fits the early prototype, but the notice channels and evolution rule require explicit product-owner approval before the first release`

改后：

`**APPROVED (a direct user ruling; a decision record)**. Four elements: ruling party = the user; date = 2026-11-16; reason = freezing the EOL notification channel before v0.1.0 enters the release flow; basis = the user's same-turn P77 kickoff ruling of 2026-11-16. 90 days fits the early prototype, but the notice channels and evolution rule require explicit product-owner approval before the first release (fulfilled on 2026-11-16; P77). **Nature**: a decision record (same class as ADR-045 / ADR-046 / ADR-047), not a temporary-authorization exception (the §1.7.4 five elements do not apply). **Citable**: as the decision basis that the S-G4 criterion "ADR-010 is approved by the product owner" is met and that the EOL notification channel is now frozen. **Effective**: 2026-11-16. **Rollback**: this ADR is not rewritten retroactively; future revisions to the EOL notification channel follow the existing exception rule (a dated ADR with rationale and expiry, explicitly approved by the product owner, synchronized in both channels).`

### 1.4 核对点

- 历史完成条件原文逐字保留：CN「90 天适合当前早期原型，通知渠道与演进规则须在首版发布前获产品所有者明确批准」、EN「90 days fits the early prototype, but the notice channels and evolution rule require explicit product-owner approval before the first release」——均未增删一字，仅在其后追加闭括号。
- 新文本不含裸管道符 `|`，不破坏表格结构；长度低于该表既有最长行（ADR-012 状态列），无渲染风险。
- 第 2 列（决策/方案）与第 3 列（替代与理由）**不动**——EOL 通知方式本身零改动。
- 状态列同时满足：首状态词 `APPROVED` ✓、四要素分号分隔 ✓、历史原文保留 + 闭括号 ✓。

### 1.5 本 ADR 性质标注（参照 ADR-045 / ADR-046 / ADR-047 逐项）

| 问 | 裁定 |
|---|---|
| 性质 | **决策记录**（ADR-045 / ADR-046 / ADR-047 同类），批准一项已编号的待批提案，不含规则语义变更 |
| 是否临时授权例外 | **不是**（§1.7.4 五要素不适用，不写「不得作为先例」） |
| 是否可作先例 | **可引用**：作为「EOL 通知方式已冻结、S-G4 判据满足」的决策依据，可按决策性 ADR 常规属性引用 |
| 生效日 | 2026-11-16（用户同轮裁定日） |
| 回退 | 不追溯改写；未来修订走既有例外规则（带日期、理由和到期日的 ADR，经产品所有者明确批准后同步两处） |

---

## Q2：候选文档与事实 G 两处的逐份判定

> 判定三档：第一层必改 / 第二层视结果 / 不触。对每处「升级」或「降级」显式披露理由。行号均为磁盘实测（与起跑消息一致；软换行空格经字节核实不存在，见末列注记）。

| # | 文档 `file:line` | 现状关键句（逐字） | 判定 | 目标文本（改前 → 改后） |
|---|---|---|---|---|
| 1 | `docs/S1_SUPPORT_MATRIX.md:21` / `_EN:21` | CN：`待审批的 ADR-010 提议：新 minor 发布时在 GitHub Release notes 与本支持矩阵同时公告前一 minor 的 EOL 日期，日期为发布后第 90 个自然日。`<br>EN：`Proposed ADR-010: when a new minor is released, announce the previous minor's EOL date in both the GitHub Release notes and this support matrix, with EOL on the 90th calendar day after release.` | **第一层必改**（与原判一致） | CN：`待审批的 ADR-010 提议：` → `已批准的 ADR-010（2026-11-16）：`，其余逐字不动<br>EN：`Proposed ADR-010:` → `Approved ADR-010 (2026-11-16):`，其余逐字不动 |
| 2 | `docs/S1_SUPPORT_MATRIX.md:23` / `_EN:23` | CN 尾句：`该提议尚未获批，本页不能替代最终公告。`<br>EN 尾句：`This proposal is not approved, and this page does not replace the final notice.` | **第一层必改**（与原判一致） | CN：`该提议尚未获批，本页不能替代最终公告。` → `ADR-010 已获产品所有者批准（2026-11-16）；本页不能替代最终公告。`（句首其余原文逐字不动）<br>EN：`This proposal is not approved, and this page does not replace the final notice.` → `ADR-010 was approved by the product owner on 2026-11-16; this page does not replace the final notice.` |
| 3 | `docs/INSTALLATION_PLAN.md:32` / `_EN:32` | CN：`S1 不指定默认安装目录，不发布一键安装脚本，不修改 PATH，不声明自动更新，也不要求管理员权限。state/store 默认位置、正式下载入口与 EOL 通知仍待冻结。`<br>EN：`S1 defines no default install directory, publishes no one-click installer, does not modify PATH, claims no automatic updates, and requires no administrator privileges. The default state/store location, formal download entry point, and EOL notification channel remain open.` | **第一层必改**（与原判一致；最小切分：三项并列仅 EOL 项变更） | CN 第二句：`state/store 默认位置、正式下载入口与 EOL 通知仍待冻结。` → `state/store 默认位置与正式下载入口仍待冻结；EOL 通知方式已由 ADR-010 获产品所有者批准（2026-11-16）。`（第一句逐字不动）<br>EN 第二句：`The default state/store location, formal download entry point, and EOL notification channel remain open.` → `The default state/store location and formal download entry point remain open; the EOL notification channel was approved by the product owner via ADR-010 (2026-11-16).` |
| 4 | `docs/PRODUCT_REQUIREMENTS.md:107` / `_EN:105` | CN 尾段：`精确依赖、EOL 通知、安全评审者和具体 whois 输入仍待后续任务落实，不因已有实现而自动取得付费调用、数据访问或发布授权。`<br>EN 尾段：`Exact dependencies, EOL notice, security reviewer, and concrete whois input remain later-task details; existing implementation grants no automatic authority for paid calls, data access, or release.` | **第一层必改，采用句内最小切分**（升级披露：台账登记块交付物④预判为「加法式指针」独立段，实读裁定改用句内切分） | CN：`精确依赖、EOL 通知、安全评审者和具体 whois 输入仍待后续任务落实，不因已有实现而自动取得付费调用、数据访问或发布授权。` → `精确依赖、安全评审者和具体 whois 输入仍待后续任务落实；EOL 通知方式已由 ADR-010 获产品所有者批准（2026-11-16）。不因已有实现而自动取得付费调用、数据访问或发布授权。`<br>EN：`Exact dependencies, EOL notice, security reviewer, and concrete whois input remain later-task details; existing implementation grants no automatic authority for paid calls, data access, or release.` → `Exact dependencies, security reviewer, and concrete whois input remain later-task details; the EOL notice was approved by the product owner via ADR-010 (2026-11-16). Existing implementation grants no automatic authority for paid calls, data access, or release.` |
| 5 | `docs/S1_RELEASE_NOTES.md:23` / `_EN:23` | CN：`2. 将版本、artifact 文件名、下载位置、EOL 通知方式和支持日期绑定到本说明及支持矩阵。`<br>EN：`2. Bind version, artifact filenames, download location, EOL notification channel, and support date to these notes and the support matrix.` | **不触**（降级披露：原判「第二层视结果」，实读判不触） | 不改。理由：该行是「发布前剩余门禁」的**绑定义务**项，不声称「未批准」；ADR-010 批准 ≠ 绑定完成（绑定对象 = 最终 ZIP 的版本/文件名/下载位置/EOL 通知方式/支持日期，由 P78 `RELEASE-ZIP-BINDING` 执行）。改动会把「批准」误记为「绑定」，属语义污染。 |
| 6 | `docs/INSTALLATION{,_EN}.md` | （0 命中） | **不触**（与原判一致） | 不改。磁盘实测无 EOL 表述，安装模型（ADR-009）与本次批准无交集。 |
| 7 | `docs/OPERATIONS{,_EN}.md` | （0 命中） | **不触**（与原判一致） | 不改。磁盘实测无 EOL 表述。 |
| 8 | `docs/ADR_REGISTER.md:62` / `_EN:63`（事实 G 第一处） | CN 尾句：`该提议尚未获批，当前仍不是已发布 SLA。`<br>EN 尾句：`The proposal is not approved, and this is not yet a published SLA.` | **纳入，第一层必改**（升级披露：P76 报告 §8 清单未含此句，实读判定该句是与批准事实直接矛盾的存量断言，必须消除；段首至「且不得追溯缩短已公告窗口」逐字不动） | CN：`该提议尚未获批，当前仍不是已发布 SLA。` → `ADR-010 已于 2026-11-16 获产品所有者批准；当前仍不是已发布 SLA。`<br>EN：`The proposal is not approved, and this is not yet a published SLA.` → `ADR-010 was approved by the product owner on 2026-11-16; this is not yet a published SLA.`（保留「不是已发布 SLA」：批准 ≠ 发布） |
| 9 | `docs/ADR_REGISTER.md:102` / `_EN:103`（事实 G 第二处） | CN 第三列：`首版前冻结 EOL 通知方式和例外`<br>EN 第三列：`Freeze EOL notice and exceptions before first release` | **纳入，追加闭括号**（升级披露：P76 §8 清单未含此格；该列是 T001 决策单的「当前约束待落实」状态列，先例 = ADR-002/ADR-003 等状态列落实后追加记录；原文不动只加闭括号，保留历史） | CN：`首版前冻结 EOL 通知方式和例外` → `首版前冻结 EOL 通知方式和例外（已于 2026-11-16 由 ADR-010 批准落实；P77）`<br>EN：`Freeze EOL notice and exceptions before first release` → `Freeze EOL notice and exceptions before first release (fulfilled by the ADR-010 approval on 2026-11-16; P77)` |

**Q2 附带裁定（句内切分 vs 独立段注记）**：PRODUCT_REQUIREMENTS 采用**句内最小切分**，理由三条：① 该句是现状/待落实清单而非门禁判据——P76 独立段形态（「（P76 加法式注记…）：…原判据不变。」）的适用前提是判据原句不可动，此处不适用；② 本仓已有行内最小切分先例（P76 对 `s1-exit.md:52` 的行内修文、ADR-046 对 `S1_SUPPORT_MATRIX:13/:14` 句尾口径修正）；③ 句内切分使矛盾断言归零（不留「仍待落实」旧句），且 CN/EN 严格同构、句子数量与位置同构。

**字节注记（重要）**：起跑消息事实 C 称 `S1_SUPPORT_MATRIX:22/:23` 存在既有软换行空格（`新 公告`、`缩短 已公告窗口`）、事实 D 称 `INSTALLATION_PLAN:32` 存在 `权 限` 空格——经逐字节核实，**磁盘上这些空格不存在**（`:22` 为「仅可延长后续新公告的窗口」，`:23` 为「且不得追溯缩短已公告窗口」，`:32` 为「也不要求管理员权限」，均为连续字）。本方案所有目标文本以磁盘实际字节为准；落盘时**不得**按起跑消息的描述补加空格。

---

## Q3：S-G4 收口的事实承载位

### 3.1 P77 报告建议总节结构（供 ②c 直接照用）

- §1 切片标识与授权（切片 ID、起跑消息同轮裁定引用、R/D 基线）
- §2 边界与不改清单确认（本方案书「不改清单确认」节）
- §3 文档改动逐项证据（改前 → 改后、file:line、字节级核对）
- §4 门禁与验收证据（三作用域 `TOTAL_FAIL=0`、SELFTEST、编码/行尾、CN/EN 对称、`该提议尚未获批` 字符串零命中断言、ADR-010 状态栏 `APPROVED` 存在性断言）
- **§5 S-G4 收口（建议节号 = §5）**
- §6 下游解锁与残留（P78 解锁、S-G3 待 P78、S-G1/S-G8/S-G7 留 S2、S-G5 记账口径）
- §7 计量与审查链（①③④⑤⑦ 与预算）
- §8 风险与开放项
- §9 ⑩ 回执要素

### 3.2 S-G4 收口章（§5）四条内容骨架

1. **判据原文（不改）**：引用 `docs/validation/p57-closeout-assessment.md` §4.2 的 S-G4 行——`| S-G4 EOL 通知方式 | ADR-010 获产品所有者批准 | 未通过：\`docs/ADR_REGISTER.md:20\`（\`PROPOSED / NOT APPROVED\`） | 产品所有者 | 批准 ADR-010 |`，并注明本片对该行及其所在门表**零改动**。
2. **证据**：`docs/ADR_REGISTER.md:20` 状态列已按 Q1 目标文本落盘（新文本逐字引用）；本报告路径 `docs/validation/adr010-approval{,_EN}.md`；批准四要素（裁决人 = 用户；日期 = 2026-11-16；理由 = v0.1.0 进入 Release 流程前冻结 EOL 通知方式；依据 = 用户 2026-11-16 P77 起跑消息同轮裁定原文）。
3. **下游**：P78 `RELEASE-ZIP-BINDING` 解锁（其依赖项 = `ADR010-APPROVAL` 闭环，见 `CHECKLIST_v60:162`）；S-G3 待 P78 后通过；S-G1 / S-G8 / S-G7 留 S2；S-G5 记「发行部分已批准、S1 Exit 整体批准待 S2」，不得记为通过。
4. **不变项**：p57 §4.2 门表与 `s1-exit{,_EN}.md` 门表零改动；S-G4 判据未改；本批准**不得**被解读为 S1 Exit 通过；`s1-exit.md:52` 的 S-G2 句仍留后续切片 `S1-EXIT-SG2-SYNC`。

---

## Q4：grep 命中清单与归属

> 归属四类：本片改动 / 本片登记（台账）/ 不改（历史冻结）/ 不改（边界排除，需披露）。清单以起跑消息事实 H 的 23 文件为准（P4 按不限窗口同行同现口径复核 = 与主控实测一致）。

**（一）本片改动**

| 文件 | 理由 |
|---|---|
| `docs/ADR_REGISTER.md` / `_EN` | ADR-010 状态列（:20）、支持策略段尾句（:62/:63）、T001 决策单「平台/支持」状态列（:102/:103）——见 Q2 第 8/9 行 |
| `docs/S1_SUPPORT_MATRIX.md` / `_EN` | 尾部政策段 :21/:23——见 Q2 第 1/2 行（平台清单与支持层级不动） |
| `docs/INSTALLATION_PLAN.md` / `_EN` | :32 三项并列最小切分——见 Q2 第 3 行 |
| `docs/PRODUCT_REQUIREMENTS.md` / `_EN` | :107/:105 句内切分——见 Q2 第 4 行 |
| `docs/validation/s1-exit.md` / `_EN` | :7 简介段与 :66 剩余风险第 1 条的叙事句（`EOL 通知方式尚未冻结` / `尚未获产品所有者批准` 为矛盾断言）；**门表九行与 AT 汇总表零改动**。P76 报告 §8 明列 `s1-exit.md:66` 为批准后改动点；:7 与本片同为矛盾断言、与 :66 同款叙事，一并改（升级披露，理由同上）。目标文本：CN:7 `…但最终 ZIP 尚未绑定和复验，EOL 通知方式尚未冻结，产品所有者也未批准最终 S1 exit…` → `…但最终 ZIP 尚未绑定和复验，EOL 通知方式已由 ADR-010 获产品所有者批准（2026-11-16），产品所有者也未批准最终 S1 exit…`；EN:7 `…the EOL notification channel remains open, and…` → `…the EOL notification channel was approved by the product owner via ADR-010 (2026-11-16), and…`；CN:66 `…EOL 通知与例外已有 ADR-010 待审批方案，但尚未获产品所有者批准。` → `…EOL 通知与例外已由 ADR-010 获产品所有者批准（2026-11-16）。`；EN:66 `ADR-010 now proposes the EOL notice and exception policy, but product-owner approval remains outstanding.` → `The EOL notice and exception policy was approved by the product owner via ADR-010 (2026-11-16).` |

**（二）本片登记（台账）**

| 文件 | 理由 |
|---|---|
| `docs/t027/REMAINING_SLICES.md` / `_EN` | P77 自身条目的状态行（🔄 进行中 → ✅ 完成）与追加「完成记录」行（台账交付物⑥）。`:655` 的 `S1-EXIT-DOC` 历史行不动。时机建议：⑧c 前由主控执行。 |

**（三）不改（历史冻结）**

| 文件 | 理由 |
|---|---|
| `docs/t027/CHECKLIST_v59.md` / `_EN` | v59 已冻结快照（边界明令）；其 S-G4 行「未通过」为 v59 时点结论 |
| `docs/t027/CHECKLIST_v60.md` / `_EN` | v60 已冻结快照（边界明令）；八处命中均为 P76/P77/P78 登记与历史表述 |
| `docs/validation/p57-closeout-assessment.md` / `_EN` | S-Gx 门表定义源，P76 已明令门表零改动；:91 的 S-G4 行是判据本身，改动即改判据 |
| `docs/validation/s1-exit-closeout-assessment.md` / `_EN` | P76 时点报告（2026-11-16 产物）；其中 §8「批准后改动点」是 P76 对 P77 的建议、§12/§15 是裁定历史、§18「下一步」是时点路线图——均为历史记录，P77 不改写（§11.1 精神），批准证据由 ADR_REGISTER:20 与 P77 报告承载 |
| `docs/validation/s1-reach-reachability.md` / `_EN` | P5（2026-09-29）时点报告；:27 的 G-F 行「未通过」以 `ADR_REGISTER.md:20` 状态列为证据锚，读者可溯源，历史结论不改写 |

**（四）不改（边界排除，需披露）**

| 文件 | 理由 |
|---|---|
| `docs/DEV_PLAN.md` | 硬边界明令不改。三处命中（:96 T018 状态段、:98 T018 待议项、:149 B4 四事分立）均为**带日期的历史记录**（2026-09-22 / 09-29 / 09-23 标注）。处置结论：不改。**为何不构成「必须停下升级」**：S-G4 的判据对象是 `ADR_REGISTER.md:20` 状态列（门级），DEV_PLAN 三处属 T018 任务级历史进度；门级与任务级不同源是 p57 §4.2 S-G6 已确立的口径（P76 已按此收口）；:96 段内已有 2026-09-29 口径注记先例证明「历史段可带时点注记而不改写状态字段」。残留差异（:96「EOL 通知方式仍缺失」与批准事实的时点差）由 P77 报告如实登记，不追溯改写。 |
| `docs/S1_RELEASE_NOTES.md` | :23 为发布前门禁列表的绑定义务项（见 Q2 第 5 行降级判定），改动会污染「批准」与「绑定」语义边界 |

**镜像侧补充（并集之外的 EN 文件，同侧判定）**：`docs/S1_RELEASE_NOTES_EN.md`（判定同 CN = 不触）、`docs/INSTALLATION_PLAN_EN.md`（判定同 CN = 必改）、`docs/DEV_PLAN_EN.md`（判定同 CN = 不改）——三者未入中文模式并集，但镜像判定跟随 CN 侧，不产生额外文件面。

---

## Q5：批量辅助协调

- 主控 `tmp/P77-ADR010-APPROVAL/` 下无 batch-01 / batch-02 定义文件（现有辅助脚本为台账插入、BOM/LF 规范化、grep 并集与证据包编码类，均非批次定义）。**本方案书即本片判定的单一权威源**；若后续产生批次定义，以本方案逐条核对，不一致处停下回报。
- 本方案与主控既有预判（REMAINING_SLICES P77 登记块、v60 登记块）的差异清单：

| 差异 | 主控预判 | 本方案实读裁定 |
|---|---|---|
| PRODUCT_REQUIREMENTS §7 EOL 项落法 | 交付物④写「加法式指针」 | 句内最小切分（理由见 Q2 附带裁定），以实读为准 |
| S1_RELEASE_NOTES | 交付物⑤「视实读触项」含之 | 实读判不触（:23 是绑定义务，批准≠绑定），以实读为准 |
| s1-exit.md | P76 §8 只列 :66 | 实读扩展至 :7（同款矛盾断言），以实读为准 |
| ADR_REGISTER :62/:102 | P76 §8 清单未含 | 实读判必改（矛盾断言 + 状态列），以实读为准 |
| 软换行空格（事实 C/D） | 起跑消息称 `新 公告` / `缩短 已公告窗口` / `权 限` 为原文字节 | 逐字节核实不存在，落盘以磁盘字节为准，不补空格 |

---

## 改动清单总表

| 文件（CN/EN 成对） | 改动性质 | CN 增删预估 | EN 增删预估 | 是否对称 |
|---|---|---|---|---|
| `docs/ADR_REGISTER.md` / `_EN` | 替换（:20 状态列整格、:62/:63 尾句）+ 追加（:102/:103 第三列闭括号） | 不增行、不删行 | 不增行、不删行 | 是 |
| `docs/S1_SUPPORT_MATRIX.md` / `_EN` | 替换（:21、:23 行内） | 不增行、不删行 | 不增行、不删行 | 是 |
| `docs/INSTALLATION_PLAN.md` / `_EN` | 替换（:32 行内） | 不增行、不删行 | 不增行、不删行 | 是 |
| `docs/PRODUCT_REQUIREMENTS.md` / `_EN` | 替换（:107/:105 行内） | 不增行、不删行 | 不增行、不删行 | 是 |
| `docs/validation/s1-exit.md` / `_EN` | 替换（:7、:66 行内） | 不增行、不删行 | 不增行、不删行 | 是 |
| `docs/t027/REMAINING_SLICES.md` / `_EN` | 替换（P77 状态行）+ 追加（完成记录行） | 加一行（完成记录） | 加一行（完成记录） | 是 |
| `docs/validation/adr010-approval.md` / `_EN` | 新建（P77 报告对，②c 产物） | 新建 | 新建 | 是 |

对称性结论：除报告对新建外，所有既有文档改动均为行内替换，两侧均不增删行；台账两侧各追加一行完成记录——`G3(a.1)` 判据满足（零非对称增删）。

---

## 不改清单确认

- S1 Exit 八门判据零改动；`docs/validation/p57-closeout-assessment{,_EN}.md` §4.2 门表与 `docs/validation/s1-exit{,_EN}.md` 门表零改动（本片对 s1-exit 的 :7/:66 为叙事段，非门表）。
- EOL 通知方式本身零改动（ADR-010 第 2/3 列逐字不动；各文档同步的仅是「已批准」状态，不含政策语义变更）。
- `S1_SUPPORT_MATRIX` 平台清单与支持层级零改动；`INSTALLATION` 安装模型零改动。
- 无 GitHub Release notes 外部动作（属 P78）。
- 其它 ADR 零改动；不新增 ADR / OB。
- 准则 `docs/DELIVERY_DIRECTIVE{,_EN}.md` 维持 v1.69 零改动；`CHECKLIST_v59/v60{,_EN}` 与 v50–v58 快照零改动；所有 `.go/.json/.yml`、schema、fixtures、CI、`go.mod` 零改动。
- `DEV_PLAN` / `SECURITY` / `ARCHITECTURE` / `CONTRACTS` / `TEST_STRATEGY` / `BUSINESS_WORKFLOWS` / `CODING_CONVENTIONS` / `DOCUMENTATION_PLAN` / `RFC-*` 零改动。
- 历史报告不改写：p57、s1-exit-closeout（P76）、s1-reach（P5）中的「未通过」类时点表述保留，靠 ADR_REGISTER:20 新状态溯源。

---

## 风险与开放项（供 ⑤/⑦ 重点核）

1. **软换行空格事实修正**：起跑消息事实 C/D 的 `新 公告` / `缩短 已公告窗口` / `权 限` 字节不存在；落盘若按起跑消息补空格会引入无源字节差。⑤ 须核所有目标文本与磁盘字节逐字一致。
2. **`**APPROVED（…）**` 形态为首例**：ADR-014~045 区间无 `APPROVED` 开头的状态列，本片首创形态 A + 状态词 `APPROVED` 组合。⑦ 须核：与四要素体例一致、对后续 ADR 状态词无扩大解释、与 ADR-012 的 `APPROVED / PARTIALLY IMPLEMENTED` 旧式形态不冲突。
3. **升级面超出 P76 §8 最小清单**（ADR_REGISTER:62/:102、s1-exit:7）：均为「矛盾断言 + 状态列」实读判定；⑤ 须核门表零改动边界未被破坏、最小切分充分性（无遗漏、无过度）。
4. **历史报告不改写与 §11.1 一致性**：p57 / s1-exit-closeout / s1-reach / CHECKLIST 快照 / DEV_PLAN 中的「未通过/待批准」时点表述保留；⑦ 须核该口径不构成「文档矛盾」判据项，必要时由 P77 报告 §6 集中登记残留差异清单。
5. **批准语义边界**：批准 ≠ S1 Exit 通过、批准 ≠ Release 完成、批准 ≠ 绑定完成；S-G5 记账口径「发行部分已批准、S1 Exit 整体批准待 S2」不得漂移；`s1-exit.md:52` 的 S-G2 句不得顺手改动（留 `S1-EXIT-SG2-SYNC`）。
6. **台账回写时机**：REMAINING_SLICES P77 状态行与完成记录行建议 ⑧c 前执行；若 ⑤/⑦ 输入包在台账回写前冻结，回写须单列 delta 复核，避免「台账已改但未经审查」缺口。

DESIGN: DONE