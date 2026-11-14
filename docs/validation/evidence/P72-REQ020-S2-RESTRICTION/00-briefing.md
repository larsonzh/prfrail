切片：REQ020-S2-RESTRICTION（P72 · 层 2 治理侧 · REQ-020 的 S2 承诺收窄登记 + 关联文档同步 · 文档类 [SLICE] · 用户同轮显式授权已内含于本起跑消息）
=====================================================================

【〇·零 起跑前自检与上下文锚定】

1. 请先完整阅读 docs/t027/CHECKLIST_v57.md（当前快照）与 docs/DELIVERY_DIRECTIVE.md（v1.68，含 §13）。

2. **本片触产品侧发布结论表述**（REQ-020 阶段归属 / TEST_STRATEGY AT-01 / S1_SUPPORT_MATRIX 语言行）——按 §1.7.2 需用户同轮显式授权，**本起跑消息即为该授权**。回执首行须明写「用户 2026-11-14 同轮显式授权已内含于起跑消息」。

3. 核对 v57 §一 的 8 个锚点。**注：v57 §一 记录的 HEAD / origin 是 v57 的生成时点锚点（`147008a`，A 提交后）；v57 的 B 提交（`e8678ae`）与层 3 修正提交（`70eccee`）落地后 HEAD 已按预期前移**。核对时按以下现时值判定：
   - HEAD = **`70eccee`**；origin/main = **`70eccee`**
   - **若 HEAD 不是 `70eccee` 且原因不是 v57 的 B 提交与修正提交 ⇒ 报实况（不按 147008a 报不符，也不静默接受未知前进）**
   - 准则 **v1.68**；SELFTEST **737/737**；三作用域 `TOTAL_FAIL=0`
   - v54 = 180 / v55 = 182 / v56 = 185（LF，CN/EN 对称）
   - `tmp` 顶层 = `gates-cache` + `.gitkeep`
   - **若 ≥2 项不符（除上述 HEAD/origin 已声明前移外）⇒ 停止并报实况。**

4. **主控模型声明**：本片为文档类 S 档，按 v1.68 §2.2 默认 **`DeepSeek V4.1 Flash (deepseek)`**。**不登记 ADR**（§13.1.8 字面未触发；同类先例 P65 / P66 / P71 纯文档片主控用 Flash）。**若 ① 阶段判定触面超出 S 档（如必须改 REQ 表结构增删条目、必须改 MVA 作为验收复合标准、必须改准则正文、必须改 §7.5 判据）⇒ 按 §13.1.8 停下提示用户。** 回执首行标 [自执行]。

5. **记忆唤醒（三条继承口径，本片必须遵守）**：
   - **D=2 防再犯句**（承 P71 报告 §8）：「执行角色（`prfrail-implementer` / `prfrail-tester` / `prfrail-documenter`）在 `tmp/` 下操作时，一律使用 `tmp/<切片ID>/`，不得在 `tmp/` 根创建任何文件、目录或 `git worktree`。P70/P71 共 4 次同型实例（主控 2 + 执行角色 2）已按 §5.1 自愈档处置；片内 D 累计达 3 次 ⇒ 强制停片。」
   - **OB-139「观察中」**（承 v57 §五）：「临时件位置违规（`tmp/` 根），P70 两次 + P71 两次共四次、跨 2 片、涉两类违规者；首状态词 `观察中`；P72/P73 再犯即第 5 次 ⇒ 停下升级用户裁决。」
   - **切片级建议第 2 项**（承 v57 §五）：「`--progress` 只覆盖状态转换，不覆盖日志类输出（不转发门禁命令实时 stdout/stderr、不消费 AgentRunner `events.jsonl`、无心跳）。P72 是文档片，不涉及长任务观测，仅作背景阅读；若 P72 中触及操作指引（如 OPERATIONS）需评估是否提示该边界。」
   - **E-1 流程纪律**（承 v56/P71）：「⑧c 后到 ⑩ 回执之间，二次回读 `tmp` 五条读数，防编辑器缓冲二次落盘；若发生，如实登记不记 D（位置合法 + 工具行为 + 已零残留三要件成立时）。」
   - 其余红线：v1.68 §11.1（历史不改写）、§7.5（S 档 ⑦ 1+1）、§7.2.1（⑥ 可跳过三条件）、§7.7（按比例抽查）、§13.1 第 7 条（批量辅助 GLM-5.3-Flash ≥2 次）、§13.4（门禁坑 1–30）、§13.5（快照结构模板）、§5.1（临时件位置闭集）约束。

【〇 本片环境声明与流水线裁剪】

1. **环境**：Windows 本机；**不开启 Ubuntu VM**；**不开启网络代理**；**不启用 Sol 双打**（纯文档，无不可测语义）。Unix 侧无关（纯 `.md`）。

2. **流水线（文档类 S 档）**：① 现状核对与差异清单 → ② 逐文档回写（按纪律 V 分小步）→ ③ 一致性检查 → ④ 门禁 → ⑤ 预审（1+1）→ ⑥ 独立扫描（按 §7.2.1 三条件判定，须 ⑦ 确认）→ ⑦ 终审（1+1，S 档）→ ⑧ 文档装配 → ⑨ 原生验证 → ⑧c → ⑩。

3. **⑥ 判定**：P72 是纯 `.md` 回写（无 `.go`/`.json`/`.yml`/schema/fixtures/脚本逻辑语义变更）⇒ §7.2.1 三条件 ①② 成立；③「不涉及状态机、门禁判据、证据模型、角色/授权契约的语义」——P72 改的是 REQ-020 阶段归属列（产品范围）+ AT-01 的语言相关断言（测试覆盖范围）+ S1_SUPPORT_MATRIX 语言行（支持矩阵），**不触状态机/门禁判据/证据模型/角色/授权契约的本体**；**注：§7.2.1 ③ 的「角色/授权契约」指执行准则层面的角色定义与授权契约（§2.2 / §1.7.4），不指产品层面的「产品范围」；P72 改的是产品范围的阶段归属表述，不触执行准则语义** ⇒ ③ 成立 ⇒ **⑥ 按 §7.2.1 判定可跳过，须 ⑦ 确认**。若 ① 判定触 3.3 停升级边界 ⇒ 停下升级，本片不继续。

4. **⑦ 盲审判定（本片不抽）**：按 §7.7 按比例抽查 + P70/P71 的窗口计数口径（P67 计入已抽片 ⇒ P68/P69/P70 为窗口内未抽片 ⇒ P71 为第 4 片已抽），**P71 抽后窗口重新起算，P72 为第 1 片**，未满 4 片 ⇒ **P72 不抽盲审**。另：P72 是 P71 抽后新窗口的**第 1 片**，窗口未满 4 片，不抽；**即便后续窗口内有评估片需选「触及面最广者」，P72 是纯文档片，触及面最窄，也不会被选中**。**回执须显式记录此判定依据，不得静默跳过。**

【一、授权边界（流水线角色与模型）】

- 主控：`DeepSeek V4.1 Flash (deepseek)`（默认，不登记 ADR；若 ① 判定须升级按 §13.1.8 停下提示用户）
- ① 现状核对与差异清单：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- ② 逐文档回写：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（文档类 ② 退化为文档工程师，§5.2）
- ③ 一致性检查：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`（文档类 ③ 退化为文档一致性检查，§5.2）
- ⑤ 预审：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（额度 1+1）
- ⑥ 独立扫描：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`（额度 3；按 §7.2.1 判定可跳过，须 ⑦ 确认）
- ⑦ 终审：`independent-reviewer` + `GPT-5.3-Codex (copilot)`（S 档 1 终审 + 1 复审；**不抽盲审**，见【〇】第 4 点）
- ⑧ 文档装配：`prfrail-documenter` + `DeepSeek V4.1 Flash (deepseek)`
- ⑨ / ⑧c / ⑩：主控（Flash）执行
- 每次 `runSubagent` 必须同时传 `agentName` 与 `model`（限定名）；**严格串行**（R5.1）

【二、批量辅助（强制 ≥2 次）】

- 载体 / 模型：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 白名单内批量辅助首选，无需临时授权）
- 实派次数：**≥ 2 次**（§13.1.7 与 §13.3.25；实派不足同轮补派）
- 任务限定：
  ① **REQ-020 现状枚举 + 关联引用**：从 `docs/PRODUCT_REQUIREMENTS{,_EN}.md` 找到 **REQ-020** 的当前行号、「阶段」列、「可验收要求」列原文；**另枚举 REQ-011 / REQ-022 / PC-01 / PC-05 / PC-06 等可能引用三语言或语言支持的行**（每项给 file:line + 原文摘要）；**不判断、只列事实**
  ② **AT 与支持矩阵的语言相关枚举**：从 `docs/TEST_STRATEGY{,_EN}.md` 找到 **AT-01 / AT-11 / AT-20** 三条的当前行号、「需求」列、「场景与断言」列原文；从 `docs/S1_SUPPORT_MATRIX{,_EN}.md` 找到 **generic/C/Go harness** 行（及其他语言相关行，如有）的当前行号与全文；**不判断、只列事实**
  ③ 不允许派「排序与 Top-N 择优」与结论性判断（§2.2 排除项）
- max_tokens ≥ 2048；回执报实得字符数；返回空则如实登记并改主控自执行

【三、五件评估任务（① 阶段）】

## 3.1 现状核对（只列事实，不判断）

- **REQ-020 现状**：`docs/PRODUCT_REQUIREMENTS{,_EN}.md` 的 REQ-020 当前行号、「阶段」列、「可验收要求」列原文
- **关联引用现状**：REQ-011（多语言 harness）/ REQ-022（文档协同）/ PC-01 / PC-05 / PC-06 等是否引用三语言或语言支持（每项给 file:line + 原文）
- **AT 现状**：`docs/TEST_STRATEGY{,_EN}.md` 的 AT-01 / AT-11 / AT-20 行原文（标注哪些含三语言/语言相关断言）
- **S1_SUPPORT_MATRIX 现状**：`docs/S1_SUPPORT_MATRIX{,_EN}.md` 的 generic/C/Go 行原文；是否提及 Python/Node/三语言
- **P68 R-07 b1 Q3 语言冻结原文**：`docs/validation/remaining-rulings{,_EN}.md` §5 的 Q3 冻结范围原文；`01-arch/r07-refined.md` 的 Q3 原文
- **P71 ⑤ Info 建议原文**：`docs/validation/s1-scope-annotate{,_EN}.md` §5 的「REQ-020 S2 承诺收窄登记」待办项原文
- **P57 §4.2 八门与语言的关系**：`docs/validation/p57-closeout-assessment{,_EN}.md` §4.2 是否有语言支持门（预期无，但须核实）
- **s1-exit AT-01 行现状**：`docs/validation/s1-exit{,_EN}.md` 的 AT-01 行（若 P72 改 AT-01，此处是否需同步）

## 3.2 差异清单与回写方案（五件齐全）

- **差异清单**（逐项）：每处文档的「现状 vs R-07 b1 Q3 冻结要求」的差异，标注改动性质（新增 / 修改 / 不动）
- **回写方案**（逐文档）：
  - `PRODUCT_REQUIREMENTS`：**REQ-020 的「阶段」列**改为「S1 不要求（仅 C+Go 契约）；S2 承诺收窄（Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结；不冻既有工具脚本 / 开放语言 ID / generic 能力）」——**只改「阶段」列**，不改「可验收要求」列、不删条目、不改编号
  - `TEST_STRATEGY`：AT-01 的三语言断言按阶段归属同步（**① 阶段须核实** AT-01 原文是否含「C+Go/三语言」字样；若含，改为「C+Go 契约（三语言实跑冻结至 S2，见 REQ-020）」；不含则标不适用）；AT-11 / AT-20 ① 阶段核实是否需同步（**只改含三语言/语言相关断言的行**，不改断言本体语义，只改阶段标注）
  - `S1_SUPPORT_MATRIX`：generic/C/Go 行的约束列**加法式**补一句「Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结至 S2（REQ-020 S2 承诺收窄）；不冻既有工具脚本 / 开放语言 ID / generic 能力」；**不改该行的支持状态判定**
  - `REMAINING_SLICES`：P72 段状态回写（步骤 0 ⬜→🔄；⑧ 🔄→✅）
- **REQ-011 / REQ-022 / PC-01 / PC-05 / PC-06 的现状核对结果**：若核实它们**未引用三语言/语言支持的阶段归属或可验收要求** ⇒ 标「不适用」；**若引用** ⇒ 停下升级（P72 边界外，可能需要另片处理）。
- **不变量**：
  - REQ 条目数不变（原 28 项）；REQ 编号连续性不变
  - REQ-020 的「可验收要求」列**零改动**
  - MVA 作为验收复合标准不动（P68 b1 明确不批 MVA 正式定义）
  - S1 Exit 八门判定本身不动（只加适用性说明）
  - `S1_SUPPORT_MATRIX` 的行数不变（不新增 Python/Node 独立行）
- **反例清单**（至少 3 条）：
  - ① 改 REQ-020 阶段列**不得**删条目或改编号（否则触发停升级）
  - ② `S1_SUPPORT_MATRIX` 的 generic/C/Go 行**不得**改成「支持 Python/Node」或「三语言正式支持」（必须是「冻结」）
  - ③ REQ-020 的「可验收要求」列**不得**改动（原「S0 验证双/三语言配置；S1 C+Go 契约；S2 三语言实跑，失败整任务阻断」须逐字保留）
- 结论行：`DESIGN: DONE`（或 `DESIGN: NEEDS_CLARIFICATION` + 待澄清项），落 `tmp/P72-REQ020-S2-RESTRICTION/01-arch/design-r1.md`

## 3.3 停升级边界（必须遵守）

- ① 判定需改 REQ 表结构（增删条目、改编号、改「可验收要求」列）⇒ **停下升级**
- ① 判定需改 MVA 作为验收复合标准 ⇒ **停下升级**
- ① 判定需改 S1 Exit 八门判定本身 ⇒ **停下升级**
- ① 判定需改 `S1_SUPPORT_MATRIX` 的正式支持范围（不是加法式加注）⇒ **停下升级**
- ① 判定需改准则正文 / CONTRACTS 语义 / SECURITY 安全边界 ⇒ **停下升级**
- ① 判定 REQ-020 当前阶段列已不是「S0-S2」或与 v57 §四 描述不符 ⇒ 停下报实况
- ① 判定需新增 `S1_SUPPORT_MATRIX` 的 Python/Node 独立行 ⇒ **停下升级**（不在 P72 边界）
- ① 判定 REQ-011 / REQ-022 / PC-01 / PC-05 / PC-06 引用三语言/语言支持的阶段归属或可验收要求 ⇒ **停下升级**（P72 边界外）
- ⑤ / ⑦ 报 Medium+ ⇒ 修复后原位重跑该环
- 触任一情形 ⇒ 停下升级用户裁决

【四、执行步骤（纪律 V：② 分小步串行）】

1. **步骤 0（起跑登记）**：
   - 建 `tmp/P72-REQ020-S2-RESTRICTION/`；落 `input-list.md` 与 `freedom-list.md`
   - **台账 P72 段从 ⬜ 翻 🔄**（`docs/t027/REMAINING_SLICES{,_EN}.md`）
   - 起跑消息逐字落盘为 `docs/validation/evidence/P72-REQ020-S2-RESTRICTION/00-briefing.md`（UTF-8 with BOM + LF），并入 `MANIFEST.md` 与 `SHA256SUMS.txt`，在 `artifacts` 围栏声明
   - 读现状（3.1 各文档），进入 ①

2. **① 现状核对与差异清单**（V4 Pro）：按【三】出五件，落 `01-arch/design-r1.md`。若触 3.3 停升级边界 ⇒ 停下升级。

3. **② 逐文档回写**（documenter，按纪律 V 拆步，每步独立核验）：
   - ②a：`PRODUCT_REQUIREMENTS{,_EN}.md` 的 REQ-020 阶段列（CN/EN 各 1 行，逐句对称）
   - ②b：`TEST_STRATEGY{,_EN}.md` 的 AT-01（+ AT-11 / AT-20 视 ① 判定；CN/EN 各 1–3 行）
   - ②c：`S1_SUPPORT_MATRIX{,_EN}.md` 的 generic/C/Go 行约束列（CN/EN 各 1 行加法式）
   - ②d：`REMAINING_SLICES{,_EN}.md` 的 P72 段状态回写（步骤 0 已翻 🔄，本步不重复）
   - **每步完成即跑 `gate.js --all --scope=tree` + CN/EN 六形状核对**；任一步红即修，不进入下一步
   - 若某步发现需要改 REQ 表结构 / MVA 复合标准 / 八门判定本身 / 新增支持矩阵独立行 ⇒ 停下升级（边界外）

4. **③ 一致性检查**（documenter）：核对 3–5 处回写与 P68 R-07 b1 Q3 冻结范围逐句对应；核对 REQ 编号连续性 / 条目数不变 / REQ-020「可验收要求」列零改动 / MVA 复合标准不动 / S1 Exit 八门判定本身不动 / `S1_SUPPORT_MATRIX` 行数不变

5. **④ 门禁**：`gate.js --all --scope=index / tree / ci` + `gate.js --selftest` 全绿；G3(a.2) 六形状全等；G4b 无 newCjk 未登记

6. **⑤ 预审**（V4 Pro，1+1）：核对回写是否与 design-r1 一致、反例是否成立、有无越界改 REQ 表结构 / MVA 复合标准 / 八门判定本身 / 支持矩阵独立行；报 Medium+ ⇒ 修复后原位重跑 ⑤

7. **⑥ 独立扫描**（MAI，3 次额度）：按 §7.2.1 三条件判定（本片倾向可跳过；须 ⑦ 确认）；若判定可跳过，登记理由；若不可跳过，跑 §7.2.2 输出

8. **⑦ 终审**（Codex，1+1，不抽盲审）：按 §7.3 只给切片定义、diff、契约原句、只读约束；首轮不附任何 ⑤/⑥ 清单；复核 REQ-020 阶段列改动是否精确对应 Q3 冻结范围、有无越界改「可验收要求」列 / 编号 / 条目数、S1_SUPPORT_MATRIX 是否为加法式
   - 报 Medium+ ⇒ 修复后原位重跑 ⑦ 复审（额度 1+1 用尽 ⇒ 按 §7.5 升级用户裁决，不自费加轮）

9. **⑧ 文档装配**（documenter）：按【五】清单逐项，CN/EN 严格对称

10. **⑨ 原生验证**：三作用域 + SELFTEST 全绿（读数为 ⑩ 回执的前置）

11. **⑧c 收尾**：`tmp/P72-REQ020-S2-RESTRICTION/` 逐件精确清理，回读五条读数；**⑧c 后到 ⑩ 回执前二次回读一次**（E-1 防再犯）

12. **⑩ 停点**：停在 ⑩ 等授权；不自行 commit/push

【五、文档同步更新（本片范围；CN/EN 严格对称）】

**必改（6 项）**：

1. **新建报告对** `docs/validation/req020-s2-restriction{,_EN}.md`（核心交付物：§1 执行摘要 / §2 差异清单 / §3 逐文档回写方案 / §4 反例与验证 / §5 审查链结论 / §6 门禁读数 / §7 已知边界 / §8 明确未执行 / `artifacts` 围栏 / `sixstage` 转录）
2. **新建证据包** `docs/validation/evidence/P72-REQ020-S2-RESTRICTION/`（`00-briefing.md` / `01-arch/design-r1.md` / `02-batch/` / `04-reviews/` / `06-gates/` / `MANIFEST.md` / `SHA256SUMS.txt` / `freedom-list.md` / `input-list.md`）
3. **`docs/PRODUCT_REQUIREMENTS{,_EN}.md`** —— **只改 REQ-020 的「阶段」列**：从当前值改为「S1 不要求（仅 C+Go 契约）；S2 承诺收窄（Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结；不冻既有工具脚本 / 开放语言 ID / generic 能力）」；**不改「可验收要求」列**、不删条目、不改编号、不改其它 REQ 行
4. **`docs/TEST_STRATEGY{,_EN}.md`** —— **AT-01 的三语言断言按阶段归属同步**（① 阶段核实 AT-01 原文含「C+Go/三语言」字样则改为「C+Go 契约（三语言实跑冻结至 S2，见 REQ-020）」；不含则标不适用）；**AT-11 / AT-20 视 ① 判定**；**只改阶段标注，不改断言本体语义**
5. **`docs/S1_SUPPORT_MATRIX{,_EN}.md`** —— **generic/C/Go 行的约束列加法式补一句**：「Python/Node 官方 harness + 三语言产品任务 + 正式支持承诺冻结至 S2（REQ-020 S2 承诺收窄）；不冻既有工具脚本 / 开放语言 ID / generic 能力」；**不改该行的支持状态判定**（仍是「配置与 noop 编排」）；**不新增 Python/Node 独立行**
6. **`docs/t027/REMAINING_SLICES{,_EN}.md`** —— P72 段状态回写（步骤 0 时 ⬜ → 🔄；⑧ 阶段 🔄 → ✅）+ 完成记录行（报告路径，按 §11.2 不写 hash/run）

**需评估（5 项，① 阶段给出判定）**：

7. **`docs/validation/remaining-rulings{,_EN}.md` §5** —— ① 阶段核实 P68 报告 §5 的 Q3 语言冻结行是否需要加「已落地」指针（P68 报告是裁定记录，一般不改；若 ① 判定加，则**加法式**加 1 行指针「Q3 语言冻结已由 P72 落地（报告：`docs/validation/req020-s2-restriction{,_EN}.md`）」；**不改原行**）
8. **`docs/S1_RELEASE_NOTES{,_EN}.md`** —— ① 阶段核实「支持与已知限制」段是否已提到「Python/Node 官方支持」「三语言实跑」——**若已提到且与 REQ-020 阶段归属一致 ⇒ 不适用**；**若需补阶段归属句 ⇒ 加法式补 1 行**
9. **`docs/validation/s1-exit{,_EN}.md` 的 AT-01 行** —— ① 阶段核实 s1-exit 的 AT-01 行（若存在）是否需同步阶段归属；**若需 ⇒ 加法式加注**；**不需 ⇒ 标注不适用**
10. **`docs/validation/p57-closeout-assessment{,_EN}.md` §4.2** —— ① 阶段核实八门清单是否涉及语言支持门（预期无，但须核实）；**若涉及 ⇒ 加法式加 1 行指针**；**不涉及 ⇒ 标注不适用**
11. **`docs/validation/s1-scope-annotate{,_EN}.md`（P71 报告）** —— **需评估**：① 阶段核实 P71 报告 §5 的「REQ-020 S2 承诺收窄登记」待办项是否需要加指针（P71 报告是历史报告，一般不改；若 ① 判定加，则**加法式**加 1 行指针「已由 P72 落地（报告：`docs/validation/req020-s2-restriction{,_EN}.md`）」；**不改原行**）。

**不适用（9 项，理由见各行）**：

12. **`docs/validation/s1-stream-progress{,_EN}.md`（P70 报告）** —— 历史报告，不追改（§11.1）
13. **`docs/validation/t027b-e2e-wiring{,_EN}.md`（P69 报告）** —— 历史报告，不追改
14. **`docs/t027/CHECKLIST_v57{,_EN}.md`** —— 快照冻结，v58 生成时继承
15. **`docs/DELIVERY_DIRECTIVE{,_EN}.md`** —— 业务切片，非治理片；准则正文零改动；P72 主控默认 Flash，不登记 ADR
16. **`docs/CONTRACTS{,_EN}.md`** —— P72 触产品范围表述，不触 wire / Schema / 状态机语义
17. **`docs/SECURITY{,_EN}.md`** —— 不改安全边界
18. **`docs/ARCHITECTURE{,_EN}.md`** —— P72 改产品层（REQ/MATRIX/TEST_STRATEGY），不改架构层
19. **`docs/DEV_PLAN{,_EN}.md`** —— P72 是 S1 范围回写，不是 T027-B 里程碑分片（T027-B 已有 P69/P70 行，P72 不加）
20. **`docs/ADR_REGISTER{,_EN}.md`** —— 不登记 ADR（不触规则语义）

**其它文档（不适用，统一说明）**：`docs/OPERATIONS{,_EN}.md` / `docs/INSTALLATION{,_EN}.md` / `docs/INSTALLATION_PLAN{,_EN}.md` / `docs/BUSINESS_WORKFLOWS{,_EN}.md` / `docs/CODING_CONVENTIONS{,_EN}.md` / `docs/DOCUMENTATION_PLAN{,_EN}.md` / `docs/RFC-*.md` —— P72 不触操作指引 / 安装模型 / 业务流 / 编码规范 / 文档治理 / 历史来源文档。若 ① 阶段核实发现其中某文档与 REQ-020 阶段归属直接相关 ⇒ 在该文档的 ⑩ 回执行中如实登记「原判不适用，实际需评估，因 ___」，并按需评估路径处理。

> 第 12–20 项及「其它文档」段在 ⑩ 回执中逐项标注「不适用 + 理由」；第 7–11 项须给评估结论（适用或不适用 + 理由）。不得省略。

【六、⑩ 回执要求】

1. 派单表（`agentName` 逐字机械读取 + 机械读取脚本 stdout 原文；机械读取失败须如实标注「未机械校验」）
2. ① 五件摘要（现状核对 + 差异清单 + 回写方案 + 不变量 + 反例清单，每项一句话 + 落点）；**⑦ 盲审判定依据**（按 §7.7 + P70/P71 窗口计数口径：P71 已抽 ⇒ P72 为第 1 片，未满 4 片不抽；须显式记录）；**⑥ 跳过判定及理由**（须 ⑦ 确认）
3. ② 分步（②a–②d 视 ① 判定）每步的 `git diff --numstat` 与 CN/EN 六形状核对读数
4. ③ 一致性检查结果（3–5 处回写的逐句对应；REQ 编号连续性 / 条目数不变 / REQ-020「可验收要求」列零改动 / MVA 复合标准不动 / S1 Exit 八门判定本身不动 / S1_SUPPORT_MATRIX 行数不变）
5. ⑤ / ⑥ / ⑦ 结论摘要（⑦ 计量 = 终审 N + 复审 M + 盲审 0）
6. 三作用域 + SELFTEST 原始行
7. 文档同步 `git diff --stat`（逐项覆盖【五】20 项 + 其它文档段；必改 6 项须有 diff；需评估 5 项须给评估结论；不适用 9 项须列理由）
8. ⑧c 后五条终态读数（v1.68 §5.1「⑧c 后取证窗口」子款）+ **⑧c 后二次回读一次**（E-1）
9. 明确未执行事项：未改准则正文 / 判据本体 / `CONTRACTS` 语义 / schema / fixtures / CI / v50–v57 / REQ 表结构 / MVA 复合标准 / S1 Exit 八门判定本身 / REQ-020「可验收要求」列；未替用户裁定；未 commit/push；未启动新切片
10. 本片自身越界登记：§1.7.2 三条 bullet 越界如实登记；未发生写「无」
11. **临时件纪律登记**：D 计数（若发生）；OB-139 相关（若再犯即第 5 次 ⇒ 停下升级）
12. 待授权项：提交授权；push 授权
13. 给 v58 的事实包（v58 生成时用：P72 闭环 hash + SELFTEST 读数 + 回写落点 + ⑦ 盲审未抽的依据 + R/D 计数 + OB-139 状态）

【七、执行节奏与红线】

1. ④ 全绿直接派 ⑤；⑤/⑦ 报 Medium+ 原位重跑该环；⑩ 是唯一常规停点
2. 本片自身受 §1.7.2 三条 bullet 约束
3. 临时件一律 `tmp/P72-REQ020-S2-RESTRICTION/`；只增不删；⑧c 统一清理；**CI 日志与门禁捕获文件也落这里，不落 `tmp/` 根**（P70/P71 教训；本片 D 累计达 3 次 ⇒ 强制停片）
4. **不改准则正文**：`DELIVERY_DIRECTIVE{,_EN}.md` 完全不改
5. **不改产品代码**：本片无 `.go` / `.json` / `.yml` / schema / fixtures / 脚本逻辑改动；不改 `internal/**`、`schemas/**`、`testdata/contracts/**`、`.github/workflows/**`、`go.mod`
6. **不改 v50–v57 历史快照**
7. **不改 REQ 表结构**：不删条目、不改编号、不改 REQ-020 以外的任何 REQ 行、不改 REQ-020「可验收要求」列
8. **不改 MVA 作为验收复合标准**（P68 b1 明确不批）
9. **不改 S1 Exit 八门判定本身**（只加适用性说明）
10. **不新增 `S1_SUPPORT_MATRIX` 的 Python/Node 独立行**（加法式加注到 generic/C/Go 行）
11. **CN/EN 严格对称**：每步改动两侧逐字对照；G3(a) 六形状全等
12. **不替用户裁定**
13. **不启动新切片**
14. **纪律 V**：② 必须按 ②a–②d 分步串行，每步独立核验
15. **E-1 流程纪律**：⑧c 后到 ⑩ 回执之间，须二次回读 `tmp` 五条读数
16. **D=2 防再犯句**：执行角色在 `tmp/` 下须用 `tmp/<切片ID>/`，不落 `tmp/` 根；片内 D=3 强制停片
17. **OB-139 观察中**：若 P72 再犯（第 5 次同型）⇒ 停下升级用户裁决

现在开始执行 P72 步骤 0。
