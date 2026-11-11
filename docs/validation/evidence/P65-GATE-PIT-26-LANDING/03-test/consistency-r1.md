# P65 GATE-PIT-26-LANDING · ③ 文档一致性检查（r1）

- 角色：③ 文档一致性检查（prfrail-documenter 载体，只读核验 + 单一证据文件）
- 日期：2026-11-11
- 核验对象：② 新增的 4 个 +1 行 —— `docs/DELIVERY_DIRECTIVE{,_EN}.md` 附录 B.7.1 各 1 行；`docs/t027/REMAINING_SLICES{,_EN}.md` P65 段 A-part 各 1 行。
- 核验方式：纯文本 / 机械核验（node 只读探针，脚本落在 `tmp/P65-GATE-PIT-26-LANDING/probe*.js`）。**未执行门禁、未改代码、未改文档主体**。
- 工作区基线：`git status --porcelain` = 4 个 `M`（两对文档）+ 证据包未跟踪目录；`git diff --stat` = 4 文件 `+48 / -0`；`git diff --cached` 为空。

## 一、逐项结果表

| 核验项 | 核验内容 | 结果 |
|---|---|---|
| A | v53 §九 门禁坑编号 1..26 连续、无缺无重；EN 同位条目与编号一致 | PASS |
| B | v53 §九 第 26 条与 §五 行号 120 原始记录一致（症状 / 根因 / 候选处置三点） | PASS |
| C | 附录 B.7.1 新槽位位置正确、模板围栏未破坏、CN/EN 槽位语义一一对应 | PASS |
| D | 新槽位与 §5.0（步骤 0）无字段重叠与命令冲突 | PASS |
| E | 新槽位引用的「门禁坑第 26 条（docs/t027/CHECKLIST 现行版 §九）」指针可解析 | PASS |
| F | 台账 A-part 行位置正确、CN/EN 同构、引用指针存在性 | PASS |
| G | 准则头部版本行与附录 C 未被改动（升版未发生） | PASS |

## 二、逐项证据

### A. 门禁坑编号连续性

- CN 侧：`docs/t027/CHECKLIST_v53.md` 的「### 门禁坑…」子节起于行号 235，下节标题起于行号 279（区段终点，不含）。区内编号条目共 **26 条**，编号序列 `1,2,3,…,25,26`，`dups=none`（探针输出 `CN_CHK: count=26 contiguous=true dups=none`）。
- EN 侧：`docs/t027/CHECKLIST_v53_EN.md` 的「### Gate pitfalls…」子节同起于行号 235，下节同起于行号 279；`EN_CHK: count=26 contiguous=true dups=none`，编号序列与 CN 逐值相同。
- 第 26 条落点：CN 与 EN 均在 **行号 278**（`CN_CHK: item26Line=278` / `EN_CHK: item26Line=278`）。
- v53 未被本片改动：`git status --porcelain` 的 4 个 `M` 全在两对文档，`CHECKLIST_v53{,_EN}.md` 均不在其中。

### B. §九 第 26 条 ≡ §五 行号 120 原始记录

- §五 行号 120（`docs/t027/CHECKLIST_v53.md`）原文要点：「连续三片（P62 `d0971e5` / P63 `de37c58` / F2-FIX `57631a3`）首跑 index 都报 `G5(d) staged is a subset of declared`，根因 = 步骤 0 建新证据包时未同轮补 `artifacts` 围栏行。候选：① 起跑消息步骤 0 强制「新增证据包 → 报告围栏同轮补行」；② 或作为门禁坑第 26 条。」
- §九 行号 278 原文要点：「（候选，F2-FIX 实证，连续三片）：步骤 0 新建证据包时未同轮在变更报告 `artifacts` 围栏补行 ⇒ 首跑 index 报 `G5(d) staged is a subset of declared`（P62 `d0971e5` / P63 `de37c58` / F2-FIX `57631a3` 连续三次同坑）。**起跑消息步骤 0 必须强制「新增证据包 → 报告围栏同轮补行（含 MANIFEST 与 SHA256SUMS 自身两文件）」**；或将其升格为门禁坑正式条目（本条即候选）。」
- 三点逐项对照：
  - 症状：两侧同列三片提交号 + 同一报错串 `G5(d) staged is a subset of declared` ⇒ 一致。
  - 根因：两侧同为「步骤 0 新建证据包时未同轮补 `artifacts` 围栏行」⇒ 一致。
  - 候选处置：两侧同为「① 起跑消息步骤 0 强制」与「② 升格为门禁坑正式条目」两支，且都标注候选 ⇒ 一致（§九 额外写明括注「含 MANIFEST 与 SHA256SUMS 自身两文件」，属同一处置的细化，不构成矛盾）。
- EN 侧同形：`CHECKLIST_v53_EN.md` 行号 120 与行号 278 为同一对记录的英译，根因句与候选两支亦可逐点对应。

### C. 附录 B.7.1 新槽位

- 位置（CN，`docs/DELIVERY_DIRECTIVE.md`）：围栏内槽位顺序为 行号 1385 `档位` → 1386 `探针额度` → 1387 `⑥⑦ 输入包范围` → **1388 `新增证据包 → 围栏同轮声明`** → 1389 `授权边界` → 1390 `试点指标声明` → 1391 `备注`。
- 位置（EN，`docs/DELIVERY_DIRECTIVE_EN.md`）：`SLICE` → `Slice definition` → `Effort tier` → `Probe budget` → `⑥⑦ input-package scope` → **1388 `New evidence pack → same-round fence declaration`** → `Authorization boundary` → `[pilot period only] Pilot metrics declaration` → `Notes`。两侧同为行号 1388，逐行同位。
- 围栏结构：CN 与 EN 的 B.7.1 区块均为 行号 1376..1392（跨度 18），`fenceOpen=1381` / `fenceClose=1392`，围栏闭合正常、语言标记与块内行数未变。
- 语义对应：槽位名（`新增证据包 → 围栏同轮声明` ≡ `New evidence pack → same-round fence declaration`）、首句（`本片是否新建证据包` ≡ `whether this slice creates a new evidence pack`）、括注（`含 MANIFEST 与 SHA256SUMS 自身两文件` ≡ `incl. MANIFEST and SHA256SUMS themselves`）、尾句（`§6.3 G5(d)` 口径引用）四项均可一一对应。⇒ 成立。
- 差异（不阻断，见观察项 O-1）：EN 尾句在 `calibre = §6.3 G5(d)` 之后另有同指复述子句，CN 无对应子句。

### D. 与 §5.0（步骤 0）的结构性关系

- §5.0 区段 = 行号 481..576，步骤 0..5 依次为：前置硬门（SW-1 / SW-2 / SW-3）、用户提供切片定义、主控校验 + 角色就绪预检、档位设定、任务包生成、阻塞标记。行号 492 的「入口模板」注记把起跑消息字段的权威指向附录 B.7（启动卡）。
- 字段面：§5.0 步骤表内没有任何「证据包 / 围栏声明 / 报告 artifacts」字段或命令 ⇒ 新槽位与 §5.0 **无字段重叠**。
- 措辞面：新槽位文本不含「步骤 0」字样，故不与 §5.0 步骤 0（前置硬门）产生同形词冲突；门禁坑第 26 条中的「起跑消息步骤 0」沿用 `CHECKLIST_v53.md` 行号 148「## 七、起跑与步骤 0 通则」的切片步骤语义（既有文本，本片未改）。
- 命令面：`授权边界` 槽位（行号 1389）仍独立，新槽位不含授权或提交语义 ⇒ 无命令冲突。
- 卡面约束：B.7 首段硬约束为「只放入口、槽位与节号引用，不得复制规则实体」。新槽位以声明槽位形态承载「是否新建证据包」的入口事实，并把口径交给 §6.3 G5(d)，未复制判据判定逻辑 ⇒ 与硬约束不冲突。

### E. 「门禁坑第 26 条」指针可解析

- 指针原文（行号 1388）：「先例 = 门禁坑第 26 条（docs/t027/CHECKLIST 现行版 §九）」。
- 「现行版」核实：`docs/t027/` 下 CHECKLIST 版本序列最高为 v53（探针列出的版本族止于 `CHECKLIST_v53.md` 与 `CHECKLIST_v53_EN.md`，无 v54）⇒ v53 即现行版。
- 「§九」核实：`CHECKLIST_v53.md` 行号 204 为「## 九、清单元规则」，其「### 门禁坑…」子节起于行号 235，第 26 条位于行号 278，形态为「（候选，F2-FIX 实证，连续三片）」⇒ 指针可解析且为**候选形态**。
- §6.3 侧引用：`G5(d)` 在 `docs/DELIVERY_DIRECTIVE.md` 的 §6.3 G5 行（行号 612）以「(d) 变更集 ⊆ 声明集」定义，与实现层 `tools/gates/lib/criteria/g5.js` 的 `G5(d)` 记录 ID 同名同义（该文件行号 239 / 260 / 281 三处）⇒ 引用语义正确（记法风险见 O-2）。

### F. 台账 A-part 行

- CN（`docs/t027/REMAINING_SLICES.md`）：行号 1871 `临时件落点声明` → **1872 `A-part 继承待办`** → 1873 `回写指针`；P65 字段表起于行号 1853，hunk 为 `@@ -1851,0 +1853,22 @@`（纯插入 22、删除 0）。
- EN（`docs/t027/REMAINING_SLICES_EN.md`）：行号 1858 `Temp artefact location` → **1859 `A-part carry-over`** → 1860 `Backfill pointer`；P65 字段表起于行号 1840，hunk 为 `@@ -1838,0 +1840,22 @@`（同为纯插入 22、删除 0）。
- 同构核对：两侧相对的「前一行 = 临时件落点声明」「后一行 = 回写指针」结构相同，语义逐句对应（v54/P66 正式化 + 本片不改 v53）。
- 指针存在性：
  - `docs/validation/evidence/P65-GATE-PIT-26-LANDING/01-arch/design-r1.md` **存在**，其 `### 4.1 A-part：…` 位于行号 46 ⇒ 解析成立。
  - `docs/validation/gate-pit-26-landing.md` 与 `docs/validation/gate-pit-26-landing_EN.md` 均**不存在**（`fs.existsSync=false`）⇒ 与「报告尚未生成、待 ⑧ 产出」的预期一致；台账已把该报告对同时列入「交付物」与「回写指针」，故登记为**预期未同步项**，不计 FAIL。

### G. 头部版本行与附录 C 未改动

- `docs/DELIVERY_DIRECTIVE.md` 行号 3：「版本：v1.64。日期：2026-11-06。状态：**CN 正文冻结**…」——**未改**。
- `docs/DELIVERY_DIRECTIVE_EN.md` 行号 3：「Version: v1.64. Date: 2026-11-06.…」——**未改**。
- 附录 C 标题：两文件均为行号 1443——**未改**。
- 机械证明：`git diff -U0` 两文件的 hunk 各仅一条 `@@ -1387,0 +1388 @@`（纯插入 1 行、删除 0 行；`--numstat` 两侧均 `1 0`）⇒ 行号 3 与行号 1443 均不在变更集内，逐字未动，**未发生升版**（符合「升版属 §1.7.2 升级边界、本片未授权」）。

## 三、发现清单

| 编号 | 位置 | 现象 | 处置 |
|---|---|---|---|
| O-1（观察，不阻断） | `docs/DELIVERY_DIRECTIVE_EN.md:1388` | EN 尾句在 `calibre = §6.3 G5(d)` 之后另带同指复述子句（`see DELIVERY_DIRECTIVE §6.3 G5(d)`），CN 侧无对应子句 | 已记录；建议 ⑧ 定稿时删去该重复子句（或 CN 同步补同句）后再提交，语义指向同一节，不构成语义不一致 |
| O-2（观察，不阻断） | `docs/DELIVERY_DIRECTIVE{,_EN}.md` §6.3 | 同一节并存两种记法：`G5(d)`（括号；G5 行子项 = 变更集 ⊆ 声明集，行号 612）与 `G5-d`（连字符；角色文件 frontmatter 五键，行号 640）。新槽位引用前者，语义正确但字面近似后者 | 已记录；系既有记法，非本片引入；建议后续治理片统一记法后再收紧引用 |

## 四、未同步项清单

1. `docs/validation/gate-pit-26-landing{,_EN}.md` 尚未生成 —— 属 ⑧ 交付物；台账「交付物」与「A-part 继承待办」已前置登记该路径（预期状态）。
2. 证据包自洽登记：本文件 `03-test/consistency-r1.md` 已在本轮同步登记进 `MANIFEST.md` 归档条目与 `SHA256SUMS.txt`（G8-b 包自洽）。
3. 报告 `artifacts` 围栏：本文件与证据包其余载荷须在 ⑧a 生成报告对时同轮补围栏行（G5(d) / G7-b）——即本片自身所治门禁坑的适用对象。
4. v53 §九 第 26 条仍为**候选形态** —— 由 v54/P66 生成时正式化（A-part 继承待办已登记）。

## 五、结论

CONSISTENCY: PASS
