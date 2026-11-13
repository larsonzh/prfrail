# P71 起跑消息原文（00-briefing）

切片：S1-SCOPE-ANNOTATE（P71 · 层 2 治理侧 · S1 范围阶段归属列回写 + visible 调试指引落地 · 文档类 [SLICE] · 用户同轮显式授权已内含于本起跑消息）
=====================================================================

【〇·零 起跑前自检与上下文锚定】

1. 请先完整阅读 docs/t027/CHECKLIST_v56.md（当前快照）与 docs/DELIVERY_DIRECTIVE.md（v1.67，含 §13）。

2. **本片触产品侧发布结论表述**（REQ 阶段归属 / S1_SUPPORT_MATRIX / s1-exit 八门适用性）——按 §1.7.2 需用户同轮显式授权，**本起跑消息即为该授权**。回执首行须明写「用户 2026-11-13 同轮显式授权已内含于起跑消息」。

3. 核对 v56 §一 的 8 个锚点。**注：v56 §一 记录的 HEAD / origin 是 v56 的生成时点锚点（`759c934`）；v56 的三个提交（主 `d06f229` + `eb8568a` + 回填 `ee7f260`）与 P70 的三个提交（主 `f77ab49` + 修复 `6a405bd` + 回填 `ca88fa3`）落地后 HEAD 已按预期前进**。核对时按以下现时值判定：
   - HEAD = **`ca88fa3`**；origin/main = **`ca88fa3`**
   - **若 HEAD 不是 `ca88fa3` 且原因不是上述六个提交 ⇒ 报实况（不按 759c934 报不符，也不静默接受未知前进）**
   - 准则 **v1.67**；SELFTEST **737/737**；三作用域 `TOTAL_FAIL=0`
   - v53 = 292 / v54 = 180 / v55 = 182（LF，CN/EN 对称）
   - `tmp` 顶层 = `gates-cache` + `.gitkeep`
   - **若 ≥2 项不符（除上述 HEAD/origin 已声明前进外）⇒ 停止并报实况。**

4. **主控模型声明**：本片为文档类 M 档，按 v1.67 §2.2 默认 **`DeepSeek V4.1 Flash (deepseek)`**。**不登记 ADR**（§13.1.8 字面未触发；同类先例 P65 / P66 纯文档片主控用 Flash）。**若 ① 阶段判定触面超出 M 档（如必须改 REQ 表结构增删条目、必须改 MVA 作为验收复合标准、必须改准则正文、必须改 S1 Exit 八门判定本身）⇒ 按 §13.1.8 停下提示用户。** 回执首行标 [自执行]。

5. **记忆唤醒（红线复述）**：
   - 本片受 v1.67 §11.1（历史不改写）、§7.5（M 档 ⑦ 1+2 + 盲审另计 +1）、§7.2.1（⑥ 可跳过三条件）、§7.7（按比例抽查）、§13.1 第 7 条（批量辅助 GLM-5.3-Flash ≥2 次）、§13.4（门禁坑 1–30）、§13.5（快照结构模板）、§5.1（临时件位置闭集）约束。
   - **E-1 流程纪律（P70 确认口径，本片执行）**：⑧c 后到 ⑩ 回执之间，须**二次回读** `tmp` 五条读数，防编辑器缓冲二次落盘；若发生，如实登记不记 D（位置合法 + 工具行为 + 已零残留三要件成立时）。
   - 单次输出 ≤ 10,000 字符；不自行 commit/push；不替用户裁定；不启动任何新切片。
   - **临时件与门禁捕获文件一律落 `tmp/P71-S1-SCOPE-ANNOTATE/`**，不落 `tmp/` 根（P70 D=2 教训：D=3 触 §5.1 累计上界强制停片）。

6. **R/D 起始值**：本片从零起算 **R = 0、D = 0**（不继承 P70 的 R=1 / D=2）。

【〇 本片环境声明与流水线裁剪】

1. **环境**：Windows 本机；**不开启 Ubuntu VM**；**不开启网络代理**；**不启用 Sol 双打**（纯文档，无不可测语义）。Unix 侧无关（纯 `.md`）。

2. **流水线（文档类 M 档）**：① 现状核对与差异清单 → ② 逐文档回写（按纪律 V 分小步）→ ③ 一致性检查 → ④ 门禁 → ⑤ 预审（1+1）→ ⑥ 独立扫描（按 §7.2.1 三条件判定，须 ⑦ 确认）→ ⑦ **盲审 1 + 终审 1 + 复审 2**（M 档额度 1+2 + 盲审另计 +1）→ ⑧ 文档装配 → ⑨ 原生验证 → ⑧c → ⑩。

3. **⑥ 判定**：P71 是纯 `.md` 回写（无 `.go`/`.json`/`.yml`/schema/fixtures/脚本逻辑语义变更）⇒ §7.2.1 三条件 ①② 成立；③「不涉及状态机、门禁判据、证据模型、角色/授权契约的语义」——P71 改的是 REQ 阶段归属列 / S1_SUPPORT_MATRIX 口径 / TEST_STRATEGY 的 AT 条目 / s1-exit 适用性说明 / OPERATIONS 的 visible 段，**不触状态机/门禁判据/证据模型/角色/授权契约的本体**（只加适用性说明，不改判定本身）⇒ ③ 成立 ⇒ **⑥ 按 §7.2.1 判定可跳过，须 ⑦ 确认**。**若 ① 阶段判定触 3.3 停升级边界 ⇒ 停下升级，本片不继续。** 不得在 ① 未判定前预设「⑥ 必跳过」或「⑥ 必跑」。

4. **⑦ 盲审判定（本片必抽 1 次）**：按 §7.7 按比例抽查 + **P70 报告 §8 已明写的口径**「P68 计入 §7.7 窗口 ⇒ P68 / P69 / P70 / P71 为四片未抽 ⇒ P71 应抽盲审」——**本片 ⑦ 必抽 1 次盲审**（不计入 M 档 1+2 的复审额度，按 §7.5「盲审另计 +1」处理）。
   - **执行时机**：⑦ 盲审先于 ⑦ 终审（盲审是「⑦ 之前追加」，§7.7）
   - **输入**：只给 P71 的切片定义（本消息）+ diff（5 处文档回写 + 1 处可选）+ 每处改动的相关上游约束原文 + 只读约束；**不附任何 ⑤/⑥ 清单**
   - **每处给什么上游原文**（逐处）：
     - REQ-009 / 019 / 007 / 021 阶段归属：给 P68 报告 §5 的移出决定原文（`docs/validation/remaining-rulings{,_EN}.md` §5 + `01-arch/r07-refined.md` 的 b1 组合原文）
     - S1_SUPPORT_MATRIX TUI 行：给 P68 报告 §5 的 REQ-009 移出决定 + `docs/PRODUCT_REQUIREMENTS.md` 的 REQ-009 原文
     - S1_SUPPORT_MATRIX visible 行：给 ADR-012 原文（`docs/ADR_REGISTER.md:22`）+ `docs/S1_SUPPORT_MATRIX.md` 现状原文
     - TEST_STRATEGY AT-13/AT-24：给两行原文 + `docs/TEST_STRATEGY.md:69`（AT-13 行）
     - s1-exit：给 S-G1~S-G8 原文 + S-G2 现时状态
     - OPERATIONS visible 段：给 ADR-012 原文 + P68 技术前提 D5 结论（`r07-refined.md` §0.D5）
   - **输出**：按 §4.1 闭集 `BLIND REVIEW: PASS / PASS WITH FIXES / PARTIAL / FINDINGS`

【一、授权边界（流水线角色与模型）】

- 主控：`DeepSeek V4.1 Flash (deepseek)`（默认，不登记 ADR；若 ① 判定须升级按 §13.1.8 停下提示用户）
- ① 现状核对与差异清单：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ② 逐文档回写：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（文档类 ② 退化为文档工程师，§5.2「说明（文档型切片）」）
- ③ 一致性检查：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（文档类 ③ 退化为文档一致性检查，§5.2）
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1，M 档）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3；按 §7.2.1 判定，本片倾向可跳过，须 ⑦ 确认）
- ⑦ 盲审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（**必抽 1 次**，按 §7.7 + P70 报告 §8 口径）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（M 档 1 终审 + 2 复审）
  - **终审输入含盲审结论，但须在任务包中明写**：「盲审结论是独立第二意见，不替代你的独立判断；请先基于切片定义 + diff + 契约原句独立完成审查，再看盲审结论，二者对照后给出终审结论」
- ⑧ 文档装配：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：主控（Flash）执行
- 每次 `runSubagent` 必须同时传 `agentName` 与 `model`（限定名）；**严格串行**（R5.1）

【二、批量辅助（强制 ≥2 次）】

- 载体 / 模型：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 白名单内批量辅助首选，无需临时授权）
- 实派次数：**≥ 2 次**（§13.1.7 与 §13.3.25；实派不足同轮补派）
- 任务限定：
  ① **REQ 现状枚举**：从 `docs/PRODUCT_REQUIREMENTS.md` 找到 **REQ-009 / REQ-019 / REQ-007 / REQ-021** 四条的当前行号、「阶段」列当前值、「可验收要求」列原文摘要；**标注四条是否已有 S1/S2 双标注**（据 P68 报告 §5，REQ-021 原文已有「S1/S2 … S2 扩展 supervised」）；**另枚举 PC-01 / PC-06 行的「REQ」列引用**（P68 报告 §5 明写 `:119` / `:124` 处有 PC 表引用）；**不判断、只列事实**
  ② **AT 与门禁现状枚举**：从 `docs/TEST_STRATEGY.md` 找到 **AT-13 / AT-24** 的当前行号、「需求」列、「场景与断言」列、「所有者/计划」列原文；从 `docs/validation/s1-exit.md` 找到 **S-G1~S-G8** 的当前行号、状态词、门禁名与证据摘要；**另枚举 s1-exit 的 AT-13 / AT-24 行**（AT 表不是八门，是独立段）；**不判断、只列事实**
  ③ 不允许派「排序与 Top-N 择优」与结论性判断（§2.2 排除项）
- max_tokens ≥ 2048；回执报实得字符数；返回空则如实登记并改主控自执行

【三、五件评估任务（① 阶段）】

## 3.1 现状核对（只列事实，不判断）

- **REQ 四条 + PC 引用现状**：`docs/PRODUCT_REQUIREMENTS{,_EN}.md` 中 REQ-009 / REQ-019 / REQ-007 / REQ-021 的当前行号、「阶段」列、「可验收要求」列原文；**PC-01 / PC-06 行的「REQ」列引用**原文（P68 §5 明写 `:119` / `:124`）
- **S1_SUPPORT_MATRIX 现状**：`docs/S1_SUPPORT_MATRIX{,_EN}.md` 的表格行（Windows / Linux / macOS / generic-C-Go / CLI Agent-AgentRunner / VS Code-Copilot-SessionBridge visible 黑箱候选 六行）；**「版本与安全维护候选政策」段与「发行信任边界」段**是否需同步（预期不适用）
- **TEST_STRATEGY 现状**：AT-13 / AT-24 的当前行号与断言
- **s1-exit 现状**：S-G1~S-G8 的当前行号与状态词；**AT 表的 AT-13 / AT-24 行**（独立段，不是八门）；**S-G2 现时状态**（P64 后可能已条件闭合，不是 BLOCKED）
- **OPERATIONS 现状**：`docs/OPERATIONS{,_EN}.md` 中是否已有 visible 的用法说明
- **S1_RELEASE_NOTES 现状**：`docs/S1_RELEASE_NOTES{,_EN}.md` 的「支持与已知限制」段 + **「发布前剩余门禁」段**（4 条）现状
- **remaining-rulings 现状**：`docs/validation/remaining-rulings{,_EN}.md` §5 的 R-07 b1 附条件 2 行原文
- **P57 §4.2 现状**：`docs/validation/p57-closeout-assessment{,_EN}.md` §4.2 八门清单的原文
- **s1-reach-reachability 现状**：`docs/validation/s1-reach-reachability{,_EN}.md` §2 门可达性判定的原文
- **DEV_PLAN T018 段现状**：`docs/DEV_PLAN{,_EN}.md` 的 T018 段（`:93`-`:96` 附近）
- **s1-stream-progress（P70 报告）§8 现状**：P70 报告 §8 的「P71 应抽盲审」原文
- **P68 报告 §5 与 `r07-refined.md`**：R-07 b1 的 X 方案（阶段归属列标注）与 visible 三条边界原文

## 3.2 差异清单与回写方案（五件齐全）

- **差异清单**（逐项）：每处文档的「现状 vs R-07 b1 附条件 2 要求」的差异，标注改动性质（新增 / 修改 / 不动）
- **回写方案**（逐文档）：
  - `PRODUCT_REQUIREMENTS`：REQ-009 / REQ-019 / REQ-007 / REQ-021 的「阶段」列从「S1」改为「S1 不要求，推迟至 S2」（**X 方案：不删条目、不改编号、不改「可验收要求」列**）；**PC-01 / PC-06 行的「REQ」列引用**是否需同步阶段归属（**① 阶段判定：同步 / 不适用 / 需加注**）
  - `S1_SUPPORT_MATRIX`：同步「TUI 移出、基础终端保留」与「visible 降级保证」两处口径；**「版本与安全维护候选政策」段与「发行信任边界」段预期不适用**（① 阶段确认）
  - `TEST_STRATEGY`：AT-13 / AT-24 的「所有者/计划」列按阶段归属同步；**不改断言本体**
  - `s1-exit`：**① AT 表的 AT-13 / AT-24 行同步**（TEST_STRATEGY 改了，此处也要同步）；**② 八门适用性说明**（移出项如何从门禁映射中退出；S-G1 / S-G8 仍 BLOCKED 不变；S-G2 已条件闭合——须核实现时状态）；**不改门禁判定本身**
  - `OPERATIONS`：新增「长任务 / 调试期如何借 visible 观察 AI 行为」段落 + 三条边界
  - `S1_RELEASE_NOTES`（需评估）：**「支持与已知限制」段**（补阶段归属句）+ **「发布前剩余门禁」段**（4 条是否需同步）
  - `remaining-rulings`（需评估）：P68 §5 的 R-07 b1 附条件 2 行是否加「已落地」指针
  - `p57-closeout-assessment`（需评估）：§4.2 八门清单是否加 P71 指针
  - `s1-reach-reachability`（需评估）：§2 门可达性是否需同步
  - `DEV_PLAN` 的 T018 段（需评估）：S1 Exit 汇总是否需同步
  - `s1-stream-progress`（需评估）：P70 §8「P71 应抽盲审」是否加「已落地」指针
- **不变量**：
  - REQ 条目数不变（原 28 项）；REQ 编号连续性不变
  - MVA 作为验收复合标准不动（P68 的 b1 明确「不批 MVA 正式定义」）
  - S1 Exit 八门判定本身不动（只加适用性说明）
- **反例清单**（至少 3 条）：
  - ① 改 REQ 阶段归属列**不得**删条目或改编号（否则触发停升级）
  - ② `S1_SUPPORT_MATRIX` 的 visible 行**不得**改成「正式支持」（必须是「降级保证」）
  - ③ `s1-exit` 的 S-G1 / S-G8 状态**不得**因本片变成 PASS（移出项 ≠ 门禁通过）
- 结论行：`DESIGN: DONE`（或 `DESIGN: NEEDS_CLARIFICATION` + 待澄清项），落 `tmp/P71-S1-SCOPE-ANNOTATE/01-arch/design-r1.md`

## 3.3 停升级边界（必须遵守）

- ① 判定需改 REQ 表结构（增删条目、改编号、改「可验收要求」列）⇒ **停下升级**（P71 边界是 X 方案，只改阶段归属列）
- ① 判定需改 MVA 作为验收复合标准 ⇒ **停下升级**（P68 的 b1 明确不批 MVA）
- ① 判定需改 S1 Exit 八门判定本身（不只是适用性说明）⇒ **停下升级**
- ① 判定需改 `S1_SUPPORT_MATRIX` 的正式支持范围（不是口径同步）⇒ **停下升级**
- ① 判定需改准则正文 / CONTRACTS 语义 / SECURITY 安全边界 ⇒ **停下升级**
- ① 判定 `REQ-009 / 019 / 007 / 021` 中有条目已不是「S1」（如已被前片移出）⇒ 停下报实况
- ⑤ / ⑦ 报 Medium+ ⇒ 修复后原位重跑该环
- 触任一情形 ⇒ 停下升级用户裁决

【四、执行步骤（纪律 V：② 分小步串行）】

1. **步骤 0（起跑登记）**：
   - 建 `tmp/P71-S1-SCOPE-ANNOTATE/`；落 `input-list.md` 与 `freedom-list.md`
   - **台账 P71 段从 ⬜ 翻 🔄**（`docs/t027/REMAINING_SLICES{,_EN}.md`；本起跑登记的 R/D 起始 = 0/0）
   - 起跑消息逐字落盘为 `docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/00-briefing.md`（UTF-8 with BOM + LF），并入 `MANIFEST.md` 与 `SHA256SUMS.txt`，在 `artifacts` 围栏声明
   - 读现状（3.1 各文档），进入 ①

2. **① 现状核对与差异清单**（V4 Pro）：按【三】出五件，落 `01-arch/design-r1.md`。若触 3.3 停升级边界 ⇒ 停下升级。

3. **② 逐文档回写**（documenter，按纪律 V 拆步，每步独立核验）：
   - ②a：`PRODUCT_REQUIREMENTS{,_EN}.md` 的 REQ 四条阶段归属列 + PC-01/PC-06 引用（视 ① 判定）
   - ②b：`S1_SUPPORT_MATRIX{,_EN}.md` 的 TUI + visible 两处口径
   - ②c：`TEST_STRATEGY{,_EN}.md` 的 AT-13 / AT-24
   - ②d：`s1-exit{,_EN}.md` 的 AT 表（AT-13/AT-24）+ 八门适用性说明 + S-G2 现时状态核实
   - ②e：`OPERATIONS{,_EN}.md` 的 visible 段（+ 三条边界）
   - ②f（视 ① 判定）：`S1_RELEASE_NOTES{,_EN}.md` 的支持与已知限制段 + 发布前剩余门禁段（视评估结果）
   - ②g（视 ① 判定，属需评估项）：`remaining-rulings` / `p57-closeout-assessment` / `s1-reach-reachability` / `DEV_PLAN` T018 段 / `s1-stream-progress`（各按评估结果加法式加指针或不适用）
   - **每步完成即跑 `gate.js --all --scope=tree` + CN/EN 六形状核对**；任一步红即修，不进入下一步
   - 若某步发现需要改 REQ 表结构 / MVA 复合标准 / 八门判定本身 ⇒ 停下升级（边界外）

4. **③ 一致性检查**（documenter）：核对 5 处 + 需评估 6 处（视判定）与 P68 报告 R-07 b1 附条件 2 的逐句对应；核对 REQ 编号连续性 / 条目数不变 / MVA 复合标准不动 / S1 Exit 八门判定本身不动 / PC-01/PC-06 引用与 REQ 阶段归属一致

5. **④ 门禁**：`gate.js --all --scope=index / tree / ci` + `gate.js --selftest` 全绿；G3(a.2) 六形状全等；G4b 无 newCjk 未登记

6. **⑤ 预审**（V4 Pro，1+1）：核对 5 处 + 需评估项是否与 design-r1 一致、反例是否成立、有无越界改 REQ 表结构 / MVA 复合标准 / 八门判定本身；报 Medium+ ⇒ 修复后原位重跑 ⑤

7. **⑥ 独立扫描**（MAI，3 次额度）：按 §7.2.1 三条件判定（本片倾向可跳过；须 ⑦ 确认）；若判定可跳过，登记理由；若不可跳过，跑 §7.2.2 输出

8. **⑦ 盲审 1 + 终审 1 + 复审 2**（Codex）：
   - **盲审**（本片必抽 1 次，先于终审）：只给切片定义 + diff + 每处上游原文（见【〇】第 4 点）+ 只读约束；**不附任何 ⑤/⑥ 清单**；输出按 `BLIND REVIEW:` 闭集
   - **终审**：输入切片定义 + diff + 契约原句 + **盲审结论**；**任务包明写「盲审结论是独立第二意见，不替代你的独立判断；先基于切片定义 + diff + 契约原句独立完成审查，再看盲审结论」**；首轮不附任何 ⑤/⑥ 清单
   - **复审**：报 Medium+ 时重跑（额度 1+2 = 3 次复审；盲审不占复审额度）

9. **⑧ 文档装配**（documenter）：按【五】清单逐项，CN/EN 严格对称

10. **⑨ 原生验证**：三作用域 + SELFTEST 全绿（读数为 ⑩ 回执的前置）

11. **⑧c 收尾**：`tmp/P71-S1-SCOPE-ANNOTATE/` 逐件精确清理，回读五条读数；**⑧c 后到 ⑩ 回执前二次回读一次**（E-1 防再犯）

12. **⑩ 停点**：停在 ⑩ 等授权；不自行 commit/push

【五、文档同步更新（本片范围；CN/EN 严格对称）】

**必改（8 项）**：

1. **新建报告对** `docs/validation/s1-scope-annotate{,_EN}.md`（核心交付物：§1 执行摘要 / §2 差异清单 / §3 逐文档回写方案 / §4 反例与验证 / §5 审查链结论（含盲审）/ §6 门禁读数 / §7 已知边界 / §8 明确未执行 / `artifacts` 围栏 / `sixstage` 转录）
2. **新建证据包** `docs/validation/evidence/P71-S1-SCOPE-ANNOTATE/`（`00-briefing.md` / `01-arch/design-r1.md` / `02-batch/` / `04-reviews/` / `06-gates/` / `MANIFEST.md` / `SHA256SUMS.txt` / `freedom-list.md` / `input-list.md`）；**`04-reviews/` 必须含盲审原文**（§13.1.13）
3. **`docs/OPERATIONS{,_EN}.md`** —— 新增「长任务 / 调试期如何借 visible 观察 AI 行为」段落 + **三条强制边界**：(1) visible **不是 S1 交付承诺**（S1_SUPPORT_MATRIX 明写「核心离线 CLI 不要求」）；(2) visible **不构成权威证据**（ADR-012 拒绝把 transcript / 完成声明 / 退出 0 当权威事实），仅调试便利；(3) 依赖 VS Code + Copilot Chat + 桥扩展 + 显式风险确认与人工监督，真实 CLI 端到端仍依赖 T027-B；不改其它内容
4. **`docs/PRODUCT_REQUIREMENTS{,_EN}.md`** —— REQ-009 / REQ-019 / REQ-007 / REQ-021 四条的「阶段」列：从「S1」改为「S1 不要求，推迟至 S2」（**X 方案：不删条目、不改编号、不改「可验收要求」列**；REQ-021 原文已有「S1/S2 … S2 扩展 supervised」双标注，据此对齐即可）；**PC-01 / PC-06 行的「REQ」列引用**按 ① 判定同步（视判定为同步 / 不适用 / 需加注）
5. **`docs/S1_SUPPORT_MATRIX{,_EN}.md`** —— 同步「TUI 移出、基础终端保留」（REQ-009 行 / 等价行的约束列）与「visible 降级保证」（VS Code / Copilot Chat / SessionBridge visible 行）两处口径；**不得**改成「正式支持」；「版本与安全维护候选政策」段与「发行信任边界」段预期不适用（① 阶段确认）
6. **`docs/TEST_STRATEGY{,_EN}.md`** —— AT-13 / AT-24 的「所有者/计划」列按阶段归属同步；**不改断言本体**
7. **`docs/validation/s1-exit{,_EN}.md`** —— **两个子项**：(i) **AT 表的 AT-13 / AT-24 行同步**（TEST_STRATEGY 改了，此处也要同步）；(ii) **八门适用性说明**（移出项如何从门禁映射中退出；S-G1 / S-G8 仍 BLOCKED 不变；S-G2 已条件闭合——须核实现时状态，不预设 BLOCKED）；**不改门禁判定本身**
8. **`docs/t027/REMAINING_SLICES{,_EN}.md`** —— P71 段状态回写（步骤 0 时 ⬜ → 🔄；⑧ 阶段 🔄 → ✅）+ 完成记录行（报告路径，按 §11.2 不写 hash/run）

**需评估（6 项，① 阶段给出判定）**：

9. **`docs/validation/remaining-rulings{,_EN}.md`** —— ① 阶段核实 P68 报告 §5 的 R-07 b1 附条件 2 行是否需要加「已落地」指针（P68 报告是裁定记录，一般不改；若 ① 判定加，则**加法式**加 1 行指针「附条件 2 已由 P71 落地（报告：`docs/validation/s1-scope-annotate{,_EN}.md`）」；**不改原行**）
10. **`docs/S1_RELEASE_NOTES{,_EN}.md`** —— ① 阶段核实**两段**：「支持与已知限制」段（是否已提到「完整 TUI 未交付」「visible 黑箱候选…尚不可运行发布」）+ **「发布前剩余门禁」段（4 条）**（是否需补阶段归属句）；**若已一致 ⇒ 不适用**；**若需补 ⇒ 加法式补**
11. **`docs/validation/p57-closeout-assessment{,_EN}.md` §4.2** —— ① 阶段核实八门清单是否需加 P71 指针（P57 是历史收口评估，一般不改；若加则加法式 1 行）
12. **`docs/validation/s1-reach-reachability{,_EN}.md` §2** —— ① 阶段核实门可达性判定是否需同步（P71 改 s1-exit 八门适用性，§2 是否需同步）
13. **`docs/DEV_PLAN{,_EN}.md` 的 T018 段** —— ① 阶段核实 T018 段（`:93`-`:96` 附近）是否需同步（S1 Exit 汇总）
14. **`docs/validation/s1-stream-progress{,_EN}.md`（P70 报告）§8** —— ① 阶段核实 P70 §8 的「P71 应抽盲审」是否需加「已落地」指针（P70 是历史报告，一般不改；若加则加法式 1 行）

**不适用（11 项，理由见各行）**：

15. **`docs/validation/t027b-e2e-wiring{,_EN}.md`** —— P71 不触 P69 接线
16. **`docs/t027/CHECKLIST_v56{,_EN}.md`** —— 快照冻结，v57 生成时继承
17. **`docs/DELIVERY_DIRECTIVE{,_EN}.md`** —— 业务切片，非治理片；准则正文零改动；P71 主控默认 Flash，不登记 ADR
18. **`docs/CONTRACTS{,_EN}.md`** —— P71 触产品侧发布结论表述，不触 wire / Schema / 状态机语义
19. **`docs/SECURITY{,_EN}.md`** —— 不改安全边界
20. **`docs/ARCHITECTURE{,_EN}.md`** —— P71 改操作层（OPERATIONS）与产品层（REQ/MATRIX/EXIT），不改架构层
21. **`docs/ADR_REGISTER{,_EN}.md`** —— 不登记 ADR；**ADR-012（visible 边界的上游）与 ADR-010（EOL 政策的上游）按 §11.1 不追改**
22. **`docs/validation/p44-gate-readiness{,_EN}.md`** —— 历史验收报告，不追改（§11.1）
23. **`docs/validation/at23-restricted-candidate{,_EN}.md`**（P61 报告）—— 历史验收报告，不追改
24. **`docs/validation/at24-supervised-black-box{,_EN}.md`**（P64 报告）—— 历史验收报告，不追改
25. **`docs/validation/t028-supervised-black-box{,_EN}.md`** —— 历史实现报告，不追改

**其它文档（不适用，统一说明）**：`docs/INSTALLATION{,_EN}.md` / `docs/INSTALLATION_PLAN{,_EN}.md` / `docs/BUSINESS_WORKFLOWS{,_EN}.md` / `docs/CODING_CONVENTIONS{,_EN}.md` / `docs/DOCUMENTATION_PLAN{,_EN}.md` / `docs/RFC-*.md` —— P71 不触安装模型 / 业务流 / 编码规范 / 文档治理 / 历史来源文档。若 ① 阶段核实发现其中某文档与 REQ 阶段归属直接相关 ⇒ 在该文档的 ⑩ 回执行中如实登记「原判不适用，实际需评估，因 ___」，并按需评估路径处理。

> 第 15–25 项及「其它文档」段在 ⑩ 回执中逐项标注「不适用 + 理由」；第 9–14 项须给评估结论（适用或不适用 + 理由）。不得省略。

【六、⑩ 回执要求】

1. 派单表（`agentName` 逐字机械读取 + 机械读取脚本 stdout 原文；机械读取失败须如实标注「未机械校验」）
2. ① 五件摘要（现状核对 + 差异清单 + 回写方案 + 不变量 + 反例清单，每项一句话 + 落点）；**⑦ 盲审执行说明**（按 §7.7 + P70 报告 §8 口径，本片必抽 1 次）；**⑥ 跳过判定及理由**（须 ⑦ 确认）
3. ② 分步（②a–②g 视 ① 判定）每步的 `git diff --numstat` 与 CN/EN 六形状核对读数
4. ③ 一致性检查结果（5 处 + 需评估 6 处的逐句对应；REQ 编号连续性 / 条目数不变 / MVA 复合标准不动 / S1 Exit 八门判定本身不动 / PC-01/PC-06 引用与 REQ 阶段归属一致）
5. ⑤ / ⑥ / **⑦ 盲审** / ⑦ 终审 / ⑦ 复审结论摘要（⑦ 计量 = 盲审 N + 终审 N + 复审 M）
6. 三作用域 + SELFTEST 原始行
7. 文档同步 `git diff --stat`（逐项覆盖【五】25 项 + 其它文档段；必改 8 项须有 diff；需评估 6 项须给评估结论；不适用 11 项须列理由）
8. ⑧c 后五条终态读数（v1.67 §5.1「⑧c 后取证窗口」子款）+ **⑧c 后二次回读一次**（E-1）
9. 明确未执行事项：未改准则正文 / 判据本体 / `CONTRACTS` 语义 / schema / fixtures / CI / v50–v56 / REQ 表结构 / MVA 复合标准 / S1 Exit 八门判定本身；未替用户裁定；未 commit/push；未启动新切片
10. 本片自身越界登记：§1.7.2 三条 bullet 越界如实登记；未发生写「无」
11. 待授权项：提交授权；push 授权
12. 给 v57 的事实包（v57 生成时用：P71 闭环 hash + SELFTEST 读数 + 5 处回写的落点 + 需评估 6 项的判定结果 + ⑦ 盲审结论 + R/D 计数）

【七、执行节奏与红线】

1. ④ 全绿直接派 ⑤；⑤/⑥/⑦ 报 Medium+ 原位重跑该环；⑩ 是唯一常规停点
2. 本片自身受 §1.7.2 三条 bullet 约束
3. 临时件一律 `tmp/P71-S1-SCOPE-ANNOTATE/`；只增不删；⑧c 统一清理；**CI 日志与门禁捕获文件也落这里，不落 `tmp/` 根**（P70 D=2 教训）
4. **不改准则正文**：`DELIVERY_DIRECTIVE{,_EN}.md` 完全不改
5. **不改产品代码**：本片无 `.go` / `.json` / `.yml` / schema / fixtures / 脚本逻辑改动；不改 `internal/**`、`schemas/**`、`testdata/contracts/**`、`.github/workflows/**`、`go.mod`
6. **不改 v50–v56 历史快照**
7. **不改 REQ 表结构**（X 方案）：不删条目、不改编号、不改「可验收要求」列
8. **不改 MVA 作为验收复合标准**（P68 b1 明确不批）
9. **不改 S1 Exit 八门判定本身**（只加适用性说明）
10. **CN/EN 严格对称**：每步改动两侧逐字对照；G3(a) 六形状全等
11. **不替用户裁定**
12. **不启动新切片**
13. **纪律 V**：② 必须按 ②a–②g 分步串行，每步独立核验
14. **E-1 流程纪律**：⑧c 后到 ⑩ 回执之间，须二次回读 `tmp` 五条读数，防编辑器缓冲二次落盘
15. **临时件 D 计数**（§5.1）：临时件违规按 v1.67 §5.1 三段处置；**本片 D 累计达 3 次 ⇒ 强制停片**

现在开始执行 P71 步骤 0。
