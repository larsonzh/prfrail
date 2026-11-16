# P76 ⑤ 预审 r1

审查者：deep-reasoner + DeepSeek V4 Pro（独立于实现方）

## 输入（实际读取的文件）

1. `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md`（P76 起跑消息原文，磁盘实读）
2. `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/01-arch/design-r1.md`（① 方案书，磁盘实读）
3. `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-ruling/ruling-2026-11-16.md`（②a 裁定登记，磁盘实读）
4. `tmp/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/b-diff.patch`（本次全部暂存改动 diff = `git diff --cached` 输出，全文实读）
5. `docs/validation/s1-exit-closeout-assessment.md` / `_EN.md`（本片报告，diff 全文 + 磁盘 grep 定点复核）
6. `docs/validation/s1-exit.md:40-80` 与 `docs/validation/s1-exit_EN.md:52`（修文对象，磁盘实读）
7. `docs/S1_SUPPORT_MATRIX.md` / `_EN.md` 全文（磁盘实读）
8. `docs/ADR_REGISTER.md:14-25`、`:40-70` 与 ADR-045/046/047 行号定位（grep）
9. `docs/SECURITY.md:1-12`、`docs/PRODUCT_REQUIREMENTS.md:100-120`（磁盘实读 + grep 行号定位）
10. `docs/validation/p57-closeout-assessment.md:84-115` 与 `_EN.md:101`（磁盘实读）
11. `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/freedom-list.md`（磁盘实读）
12. `docs/DELIVERY_DIRECTIVE.md` §1.7.3（:250）、§4.2.1（:457-465）、§5.1 时序说明（:518）、§7.2.1（:713-718）、§12.4 台账（:1027 起）与变更日志（grep 定点）；§1.7.4 五要素经 ADR-033/036 等既有登记行的措辞对照核对
13. 机械读数复核源：`docs/t027/CHECKLIST_v57.md:45`、`v58`、`v59.md:47/:145`；`docs/ADR_REGISTER.md`；实跑 `node tools/gates/gate.js --selftest` 与滚动窗口行数计数（终端输出留存于会话资源文件）

## 一、最低核验动作（五条逐条结果）

### ① 报告声明的 `file:line` 与磁盘实际内容是否一致 —— 通过（抽查全部关键锚点）

| 报告引用 | 磁盘实读 | 结论 |
|---|---|---|
| `s1-exit.md:52` 状态列 = `BLOCKED`，修文后「T027 未完成、AT-23 整体未通过；」 | 磁盘 :52 逐字一致（grep + 实读） | ✓ |
| `s1-exit.md:43/:44`（AT-23/AT-24 行） | 磁盘 :43=阶段 1 条件闭合达成、:44=条件闭合 | ✓ |
| `DEV_PLAN.md:94`（T018 依赖「T017、T025、T026、T027、T028 及 whois 输入授权」）、`:96`（`IN_PROGRESS / BLOCKED`）、`:119`（T028 `COMPLETED`） | 磁盘实读逐字一致 | ✓ |
| `ADR_REGISTER.md:19`（ADR-009）/`:20`（ADR-010 `PROPOSED / NOT APPROVED`）/`:56`（ADR-045） | grep 定位 :56=ADR-045（ADR-046=:57、ADR-047=:58） | ✓ |
| `PRODUCT_REQUIREMENTS.md:69`（REQ-007）/`:71`（REQ-009）/`:82`（REQ-020）/`:107`（「精确依赖…仍待后续任务落实」）/`:109`（「独立安全评审结论仍是进入 S1 前的阻断门禁」） | grep 定位 :109、:111（注记）逐字一致 | ✓ |
| `SECURITY.md:5`（「未完成独立安全签署」）、注记行 :7 | 磁盘实读一致 | ✓ |
| p57 §4.2 门表 :86-:95、P71 指针 :97、S-G8 指针 :99、P76 指针 :101 | 磁盘实读：:101=P76 指针（EN :101 同位） | ✓ |
| `DEV_PLAN.md:27`（T001 行「独立安全评审仍按后续门禁落实」） | batch-03 引用块全文已复核，行号成立 | ✓ |

未发现任何 `file:line` 与磁盘不符的引用。

### ② diff 与报告描述是否一致 —— 一项不符（见 H-1）

- 报告 :127「2. S-G2 门级落定（本片已成文，Q2 = a）」与磁盘不符：`s1-exit.md:52` 的 S-G2 句未更新（方案书 §3 附条件修文未落地），详 H-1。
- 其余声称与 diff 一致：p57 只加指针（门表八行零改动）✓；S1_SUPPORT_MATRIX 只改 :13/:14 句尾 ✓；SECURITY/PRODUCT_REQUIREMENTS 加法式注记 ✓；s1-exit 只改 :52 一行 ✓；证据包新文件与 SHA256SUMS 条目集 ✓。

### ③ CN / EN 改动逐字同位、语义等价 —— 通过

- `s1-exit:52`：CN「AT-23 整体未通过」= EN「AT-23 as a whole is not passed」，同位同义 ✓。
- p57 P76 指针：CN/EN 均落 :101，语义等价 ✓。
- `S1_SUPPORT_MATRIX:13/:14`：CN/EN 同位、语义等价（「R-3 = (A) 暂不可起跑」=「R-3 = (A), cannot start for now」；「blocked 至 S2 enforcement 接入」=「stays blocked until S2 enforcement lands」）✓。
- `SECURITY{,_EN}` 注记、`PRODUCT_REQUIREMENTS{,_EN}` 注记：同位、语义等价 ✓。
- `ADR_REGISTER{,_EN}`：ADR-046/047 两行均追加于 ADR-045 之后，语义等价 ✓。
- 报告对：CN/EN 均为 206 行，逐节对称 ✓。

### ④ 边界声明（未改清单）与 diff 实际相符 —— 通过

diff 实际触及文件 = `ADR_REGISTER{,_EN}`、`PRODUCT_REQUIREMENTS{,_EN}`、`S1_SUPPORT_MATRIX{,_EN}`、`SECURITY{,_EN}`、`p57-closeout-assessment{,_EN}`、`s1-exit{,_EN}:52`、新报告对、evidence 包新文件。全部落在 briefing【八】第一/二层清单内；第三层「确定不触」（准则 / CONTRACTS / s1-reach / CHECKLIST_v59 / v50–v58 快照 / at23 / at24 / r08-r09 / 代码与 CI）零触及 ✓；`REMAINING_SLICES` 未触（⑧b 待做，符合流程）✓；s1-exit 其余行零改动 ✓。

### ⑤ 机械读数可复算 —— 通过（全部实跑/磁盘核对）

| 项 | 报告 §2 值 | 本环复算 | 结论 |
|---|---|---|---|
| SELFTEST | 737/737 | 实跑 `node tools/gates/gate.js --selftest` = `SELFTEST: PASS 737/737` | ✓ |
| 滚动窗口 | v57=160 / v58=148 / v59=180 | `(Get-Content).Count` 实测 160/148/180 | ✓ |
| ADR 最大号 | 47（新增 046/047） | 磁盘 `ADR_REGISTER.md:58`=ADR-047 | ✓ |
| OB 最大号 | 139 | `DELIVERY_DIRECTIVE.md` v1.68 行 + CHECKLIST_v59:47「OB max = 139」 | ✓ |
| DR 最大号 | 14 | CHECKLIST_v59:47「DR max = 14」（P70 新增 DR-14） | ✓ |
| 试点顺延 | 第 25 次 | CHECKLIST_v59:47/:145「试点顺延 = 第 25 次」 | ✓ |

（三作用域 `TOTAL_FAIL=0` 属步骤 0 干净工作区的基线读数，当前 staged 状态下不可直接复现，属预期，不计为不符。）

## 二、发现清单

### H-1（High）：`s1-exit.md:52` 的 S-G2 句未随 Q2=(a) 裁定更新，与同一 diff 内三处新口径冲突

- 位置：`docs/validation/s1-exit.md:52`（与 `s1-exit_EN.md:52`）。
- 原文摘录（磁盘实读）：`…**S-G2 口径（用户裁决③）**：AT-24 条件闭合（P64）⇒ v54 快照生成时可标记「AT-24 条件闭合可推进」；S-G2 门级判定待 v54 快照，不得据此声称该门通过`
- 同 diff 内冲突证据：
  - `docs/validation/p57-closeout-assessment.md:101`（P76 指针）：「S-G2 判据满足 ⇒ 通过」；
  - `docs/validation/s1-exit-closeout-assessment.md:127`：「2. S-G2 门级落定（本片已成文，Q2 = a）」；
  - `docs/ADR_REGISTER.md:57`（ADR-046）：「Q2 = (a) 判据字面（…S-G2 判据满足 ⇒ 通过）」。
- 依据：方案书 §3「附：Q2 裁定后的条件修文（第二处，仅当 Q2 采纳口径 (a) 时执行）」——Q2 已裁定 (a)，条件成立；方案书 §12 第 2 步「更新 `s1-exit.md:52` 的 S-G2 句（§3 附条件修文）」被用户裁定六「认可 §12 序列」覆盖。且 v54 快照早已生成，「待 v54 快照」句已过时。
- 后果：S1 Exit 宣告的准入条件行内残留与最新裁定相反的表述（「不得据此声称该门通过」vs「S-G2 = 通过」）；报告 :127「本片已成文」的声称与磁盘不符。
- 建议：②b 按方案书 §3 附文本补齐 :52 的 S-G2 句更新（CN/EN），或经用户裁定明确保留该句；未补齐前不得宣告 S1 Exit 通过。

### L-1（Low）：报告 :11「六项门级口径修正」措辞不准

- 原文：`docs/validation/s1-exit-closeout-assessment.md:11`「…落地六项门级口径修正、两项 ADR、三处加法式注记与一处门表行内口径修文。」
- 说明：六项裁定中 Q7（安全评审豁免）、Q9（支持矩阵口径）、Q10（批准时机）不是「门级口径修正」。建议改为「六项裁定落地」或精确分列。

### L-2（Low）：报告 :87「全域无评审报告」字面不精确

- 原文：`docs/validation/s1-exit-closeout-assessment.md:87`「…以及 `docs/validation/` 全域无评审报告。」
- 说明：`docs/validation/` 下存在 49 处 `independent review` 记录（batch-03 §5.4），准确表述应为「全域无『独立安全评审』报告」（设计书 §7 第 4 点原文即此口径）。建议对齐措辞，避免被误读。

### L-3（Low）：SECURITY / PRODUCT_REQUIREMENTS 注记未含 ADR-047 的到期日上限

- 位置：`docs/SECURITY.md:7`、`docs/PRODUCT_REQUIREMENTS.md:111`（EN 对应行）。
- 原文摘录：`…已由 **ADR-047**（2026-11-16）豁免顺延至 S2，须在 S2 启动前补评。`
- 说明：ADR-047 到期日 =「S2 启动日与 2027-03-31 中较早者」，两处注记只写「S2 启动前」义务、未带 2027-03-31 上限。注记与 ADR 不矛盾（ADR 为权威），但建议注记补「且不晚于 2027-03-31」或注明以 ADR-047 为准。

### L-4（Low，流程观察）：freedom-list 仍为初稿，②b 阶段选择尚未回填

- 位置：`docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/freedom-list.md:8`（「## 一、已行使的执行自由度（初稿，待 ⑦ 前定稿）」）与 `:29`（「## 四、定稿（⑦ 前，待回填）——最终选择：待定。」）。
- 说明：当前处于 ⑤ 阶段，「待定」状态符合 §5.1 时序说明（OB-48：⑥ 之后、⑦ 之前定稿）。但 ⑦ 前必须回填：②b 阶段的选择（Q3 方案 A 落地、ADR-047 到期日 2027-03-31 的拟定、S1_SUPPORT_MATRIX 落地文本较设计书 §9 建议文本的简化），并按定稿规则重算 `SHA256SUMS.txt`；否则触发 OB-48 的时序偏离。

### L-5（Low）：方案 A/B 文本定义的可回溯性弱；ADR-046 依据不完整

- 位置：`docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/02-ruling/ruling-2026-11-16.md:50`（X 授权「方案 A 优先、方案 B 仅当方案 A 实读后仍不足时用」）；`docs/ADR_REGISTER.md:57`（ADR-046 依据 =「P76 ① 停点六项裁定 + X 授权」）。
- 说明：① 方案 A（加「整体」）与方案 B 的具体文本对照未在 ruling / 报告 / ADR 中成文（方案 B 推测为设计书 §3 的修文文本）；报告 §5 只述结果不述定义，按 §4.2.1 第 7 条授权可回溯口径建议补一条方案 A/B 文本对照。② ADR-046 依据未提 ruling 第七节「附加（p57 不改只加指针）」，而其影响面含 p57 指针；建议依据补「+ 附加项」。

## 三、重点核对项逐项结论

1. **Q3 修文加法式**：✓ 通过。状态列仍为 `BLOCKED`（磁盘 :52 实读）；仅插入「整体」/「as a whole」二字，无删除、未新增「阶段 2 不适用」的重复表述（行内原已有两处，X 授权理由成立）；CN/EN 语义等价。
2. **ADR-046 / ADR-047 决策记录性质 + 四要素**：✓ 通过。两条均明写「决策记录（不是 §1.7.4 临时授权例外）」（ADR-047 另标「风险接受类」）；裁决人 = 用户 / 日期 = 2026-11-16 / 理由 / 依据四要素齐备（ADR-046 依据完整性见 L-5）。
3. **ADR-047 到期日自洽**：✓ 通过。「到期日 = S2 启动日与 2027-03-31 中较早者」在 ADR 内自洽（义务 + 上限 + 回退三句闭环）；与报告 §9、EN 对应行一致。注记缺上限见 L-3。
4. **S1_SUPPORT_MATRIX 只改状态口径**：✓ 通过。:13/:14 仅改第 3 列约束句尾；第 2 列支持层级（「规划、未交付」/「产品侧已实现（降级保证）」）零改动；:5 未动。
5. **p57 只加加法式指针、门表数据零改动**：✓ 通过。diff 仅表后 +2 行；门表 :86-:95 八行判定词零改动；P71/P75 两条既有指针保留。
6. **报告 §2 机械读数**：✓ 全部与磁盘一致并可复算（见一·⑤：ADR=47 / OB=139 / DR=14 / 顺延=25 / SELFTEST 737/737 / 160/148/180）。
7. **报告越界检查（REQ-007 黑箱子项 / REQ-009 / REQ-020 算进 S1 判定）**：✓ 未越界。报告 §4/§13 的 S-G1/S-G2/S-G8 判定依据仅来自 R-3=(A)、T027 装配、候选兼容性与 T028/AT-24 证据，未夹带推 S2 项；§13 明写禁止约束。
8. **自由度清单 F-1..F-5 与实际执行**：✓ 与 ① 阶段执行一致（单语方案书 / P75 体例目录 / 四列映射表 / 四状态分支 / 批量只作枚举）；未登记的自由度 = ②b 阶段选择（X 授权下方案 A 落地、ADR-047 到期日拟定、支持矩阵落地文本简化），须 ⑦ 前定稿回填（见 L-4/L-5）。

## 四、结论

PRE-REVIEW: FINDINGS

- Critical：0；High：1（H-1）；Medium：0；Low：5（L-1..L-5）。
- 整改要求：H-1 必须处置（补齐 `s1-exit{,_EN}.md:52` 的 S-G2 句或取得用户明确保留裁定，并同步报告 :127 的「已成文」表述）；L-1..L-5 建议同轮顺带处置。按 briefing【五】，r1 报 High（Medium+）⇒ 整改后须跑 ⑤ r2（retryOf = r1）。
- 本审查为只读；未修改任何受审文件。唯一新建文件 = 本审查原文。
