# design-r1.md — P65 GATE-PIT-26-LANDING · ① 执行计划确认（轻量化，三案对比）

- 角色：① deep-reasoner + DeepSeek V4 Pro；只读，本片未修改任何文件。
- 工作目录：D:\LZProjects\prfrail；产出由主控落盘 `01-arch/design-r1.md`。

## 1. 输入读取确认

| # | 输入 | 实读位置 | 关键事实 |
|---|---|---|---|
| 1 | v53 §五 观察项原文 | `CHECKLIST_v53.md:120`（F2-FIX 观察项 1） | 症状=连续三片首跑 index 报 `G5(d) staged is a subset of declared`；根因=步骤 0 建新证据包未同轮补 `artifacts` 围栏行；候选处置=①起跑消息步骤 0 强制同轮补行 ②门禁坑第 26 条；累计 3 次 |
| 2 | v53 §九 门禁坑 | 表头 `:235`、版本注 `:237`、第 26 条 `:278` | 1–25 条格式惯例 = **加粗标题（实证来源标注）＋症状段＋处置段**；第 26 条现为候选形态 |
| 3 | 准则附录 B.7.1 | `DELIVERY_DIRECTIVE.md:1376–1391` | 模板 9 行；槽位先例 = 「⑥⑦ 输入包范围」（`:1387`，引用 §7.3） |
| 4 | 准则 §5.0 | `DELIVERY_DIRECTIVE.md:481–493` | 步 0=前置硬门（SW-1/2/3）；步 1=切片定义；步 2=主控校验；步 4=生成任务包。**§5.0 不枚举起跑消息字段** ⇒ 与 B.7.1 无结构重叠 |
| 5 | 五处命中 | F2-FIX `04-reviews/09-gates-r1.txt:63–65`；P62/P63 见 `at23-restricted-candidate.md:394`；P64×2 见 P65 briefing §零.5 + P64 briefing `:169` | F2-FIX 漏网形态 = 证据包 `MANIFEST.md`/`SHA256SUMS.txt` **自身两文件**未入围栏 |

**发现 1（briefing 文件映射错误）**：P65 briefing 3.1/3.2 把落点 A 的改动文件写作 `DELIVERY_DIRECTIVE{,_EN}.md 的 §九`——仓库事实：**门禁坑列表只存在于 `docs/t027/CHECKLIST_v53{,_EN}.md` §九**（主控机械核实：`DELIVERY_DIRECTIVE.md` 中「门禁坑」出现 0 次），DELIVERY_DIRECTIVE 无编号门禁坑列表（其 §9.3 为「已知异常模式」，非同物）。按本任务书必读输入 2 字面，落点 A 实体目标 = `CHECKLIST_v53{,_EN}.md` §九。

**发现 2（冻结矛盾）**：briefing 【四】14 声明 CHECKLIST_v53「快照冻结、v54 生成时继承」（主控机械核实：v53 由单一提交 `cef2ddb` 创建、此后未改动）。本设计推荐**对 v53 做最小化编辑**（v53 为现行版；门禁坑各版本均有在位新增先例：v49 新增第 22 条等）；备选路线 = 同一编辑集写成「v54 生成继承指令」落 `REMAINING_SLICES{,_EN}.md` P65 段。二选一由主控在 ② 前裁定，不属 3.4 停升级五类（主控裁定：因涉范围授权，已按 §1.7.2 停下升级，见 `01-arch/escalation-r1.md`）。

## 2. 三案逐项对比表

| 维度 | A：只落 §九（候选→正式） | B：只落 B.7.1 槽位 | C：A+B |
|---|---|---|---|
| 改动文件 | `CHECKLIST_v53{,_EN}.md` §九（+可选 §五 状态行） | `DELIVERY_DIRECTIVE{,_EN}.md` B.7.1（+1 行） | 两处（CN/EN 共 4 文件） |
| 防护强度 | 弱——提醒性质，开场读一次，无前置强制 | 中——起草起跑消息时必填槽位，纪律前置到失效时点 | 强——提醒层 + 前置层双覆盖 |
| 维护成本 | 低 | 低 | 中——两处静态互引（一次性对齐） |
| 可机械核验性 | 无（文档级） | 无（文档级） | 无；编号连续可由 G2/G5 类判据复核 |
| 与历史先例一致性 | 高——1–25 条全部只落 §九；第 11 条「围栏与磁盘枚举必须一致」为最邻近同族先例 | 高——「⑥⑦ 输入包范围」（P33）与「纪律 D 预置说明」（P32，指针+最小槽位名、不复述全文）两个槽位先例 | 双落点先例 = v1.54 纪律 D（正文 + 三处模板） |
| 与 §5.0 的冲突面 | 无 | 无——§5.0 管主控流程、B.7.1 管用户起跑消息内容，无字段重叠；③ 做措辞对照 | 无 |
| 对报告围栏/G8-b 的连带 | 无直接连带 | 正向——槽位含「MANIFEST 与 SHA256SUMS 自身两文件」，直击 F2-FIX 漏网形态；与 G8-b 同族 | 正向且双覆盖 |
| 与 G5(d) 判据本体距离 | 文档级，不触 `tools/gates/**` | 文档级，仅引用 §6.3 G5(d) 节号 | 文档级 |
| 对起跑消息长度影响 | 无 | +1 行（约 60 字），对 10,000 字符输出上限无实质影响 | +1 行 |
| 残余风险 | **P64 反例**：坑已登记为候选（v53 §五+§九，P64 briefing 亦点名）仍再中 2 次 ⇒ 单一提醒通道实证不足 | 模板被跳过不填 | 最低；残余 = 步骤 0 之后中途建包（判据本体范畴，不入本片） |
| 推荐度 | 不推荐单独 | 不推荐单独 | **推荐** |

## 3. 推荐案：C（A+B 双落点）

理由（引用证据）：
1. **A 单独不足有直接反例**：v53 §五/§九 已登记候选，P64 起跑消息（`:169`）甚至点名该候选，P64 前两轮仍再中 2 次（P65 briefing §零.5：累计 5 次）⇒ 仅「已知坑清单」不防复发。
2. **根因时点 = 步骤 0 建包时**（v53 §五原文）；B 把声明义务前置到起跑消息起草时刻（正是失效时点），A 提供可回溯条目与编号，互补。
3. **先例一致**：A 沿第 11 条同族先例；B 沿「⑥⑦ 输入包范围」槽位先例（引用节号、不复制实体，满足 B.7 硬约束「本卡不含规则实体」）；双落点在 v1.54 纪律 D 已有先例。
4. **成本可接受**：4 文件、每文件 1–3 行；两处互引是静态一次性对齐，措辞权威 = v53 §五/§九 原文，两处只引用不复述。

## 4. 推荐案的精确改动描述

### 4.1 A-part：`docs/t027/CHECKLIST_v53{,_EN}.md`（3 处编辑）

**A1（CN `:278`，第 26 条候选 → 正式；逐字文本）**：

> 26. **G5(d) 新增证据包围栏声明坑**（F2-FIX 实证，连续三片；P64 再证）：步骤 0 新建证据包时未同轮在变更报告 `artifacts` 围栏补行 ⇒ 首跑 index 报 `G5(d) staged is a subset of declared`（P62 `d0971e5` / P63 `de37c58` / F2-FIX `57631a3` 连续三次同坑；P64 `2dbc54a` 前两轮再中两次，累计 5 次）。**起跑消息步骤 0 必须强制「新增证据包 → 报告围栏同轮补行（含 MANIFEST 与 SHA256SUMS 自身两文件）」**（附录 B.7.1 已设对应槽位）。

与 v53 §五 原始记录关系：症状/根因/处置三点措辞与候选原文逐字一致，仅两处事实性增补——① 证据计数 3→5（P64×2，P65 briefing 既有口径）；② 删候选尾句「或将其升格为门禁坑正式条目（本条即候选）」、补「附录 B.7.1 已设对应槽位」互引。**未改写为更松或更严**。

**A2（CN `:235` 表头、`:237` 版本注）**：
- `:235`：`…+ F2-FIX 一条（第 26 条候选））` → `…+ F2-FIX 一条（第 26 条））`
- `:237`：`**v53 新增第 26 条候选（来自 F2-FIX）**` → `**v53 新增第 26 条候选（来自 F2-FIX）；P65 正式化（GATE-PIT-26-LANDING）**`

**A3（CN `:120` §五 观察项 1 尾句，可选，随主控裁定）**：
`累计 3 次，落点（起跑消息纪律 or 门禁坑）由后续治理片评估。` → `累计 5 次（含 P64 前两轮 ×2）；落点已由 P65 裁定：双落点（§九 第 26 条正式条目 + 附录 B.7.1 槽位）。`

**EN 要点（`CHECKLIST_v53_EN.md` 同位三处；关键术语译法）**：
- `staged is a subset of declared` = 判据输出串**原样保留不译**；围栏 = `fence`（`artifacts fence`）；同轮 = `same round`；新增证据包 = `new evidence pack`；自身两文件 = `the two files MANIFEST and SHA256SUMS themselves`；门禁坑第 26 条 = `gate pitfall No. 26`；候选/正式 = `candidate / formal entry`。
- EN 第 26 条正式文本（与 CN 全句对齐）：
> 26. **the G5(d) new-evidence-fence declaration pitfall** (F2-FIX evidence, three consecutive slices; reconfirmed by P64): step 0 created the new evidence pack without adding the fence rows to the changed report's `artifacts` fence in the same round ⇒ the first index run reports `G5(d) staged is a subset of declared` (P62 `d0971e5` / P63 `de37c58` / F2-FIX `57631a3`, the same pitfall three times in a row; P64 `2dbc54a` hit it twice again in its first two rounds, 5 hits in total). **The briefing's step 0 must force "new evidence pack → report fence rows added in the same round (including the two files MANIFEST and SHA256SUMS themselves)"** (Appendix B.7.1 now carries the corresponding slot).
- EN 表头：`…+ F2-FIX one (candidate #26))` → `…+ F2-FIX one (#26))`；版本注：`**v53 adds candidate #26 (from F2-FIX)**` → `**v53 adds candidate #26 (from F2-FIX); formalized by P65 (GATE-PIT-26-LANDING)**`。

### 4.2 B-part：`docs/DELIVERY_DIRECTIVE{,_EN}.md` 附录 B.7.1

**插入位置**：CN **第 1387 行之后**（现有槽位「⑥⑦ 输入包范围」之后）、第 1388 行「授权边界」之前；EN 同位。
**槽位名**：`新增证据包 → 围栏同轮声明` / `New evidence pack → same-round fence declaration`
**最小措辞（1 行，只引用不复述）**：
- CN：`新增证据包 → 围栏同轮声明：[本片是否新建证据包；若新建 ⇒ 同轮在变更报告 artifacts 围栏补行（含 MANIFEST 与 SHA256SUMS 自身两文件）；先例 = 门禁坑第 26 条（docs/t027/CHECKLIST 现行版 §九），口径 = §6.3 G5(d)]`
- EN：`New evidence pack → same-round fence declaration: [whether this slice creates a new evidence pack; if yes ⇒ add rows to the changed report's artifacts fence in the same round (incl. MANIFEST and SHA256SUMS themselves); precedent = gate pitfall No. 26 (docs/t027/CHECKLIST current §9), calibre = §6.3 G5(d)]`

**与「⑥⑦ 输入包范围」的先后关系**：置于其后——该槽位声明**输入**包范围，新槽位声明**输出**包（证据包）围栏义务，同属「包范围」类别相邻；置于「授权边界」（commit/push/付费，另一类别）之前。符合 P32 纪律 D 槽位「指针 + 最小槽位名、不复述全文」先例与 B.7 硬约束。

## 5. 与 3.4 停升级边界的关系

- **不触 G5(d) 判据本体**：不改 `tools/gates/**`、不新增白名单/机器判据；B 槽位对 §6.3 G5(d) 仅节号引用。
- **不触 CONTRACTS / schema / fixtures**、**不改产品代码**：两处落点均为 `.md` 文档。
- ⇒ **不触发升级裁决**；结论行保持 `DESIGN: DONE`。
- 发现 1/发现 2 不属 3.4 五类停升级项，按 3.3 作为文本缺陷落 `03-test/consistency-r1.md` 逐条列明，由 ⑧ 修正或 ⑩ 如实登记。

## 6. 未采纳替代（不推荐，仅留痕）

1. 改判据本体加白名单/机器判据——超界（3.4 明令停下升级），且坑属执行纪律缺陷而非判据误报。
2. 只落 §9.3 故障模式库——与门禁坑列表是两套清单，再落第三处制造三处漂移。
3. 只落 §五 观察项——无强制面，P64 反例已证不足。

## 7. 风险与移交主控项

- ② 前需主控裁定「A-part 编辑 v53 在位 vs v54 继承指令」二选一（本设计推荐前者）。
- 第 26 条与第 11 条同族：正式条目以互引避免重复正文；③ 核验确认两句不互斥。
- ③ 需按 3.3 核验：编号 1–26 连续、CN/EN 同位、B.7.1 与 §5.0 无结构性冲突。

DESIGN: DONE
