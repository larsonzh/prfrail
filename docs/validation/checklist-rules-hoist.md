# 切片验证报告 · P66 CHECKLIST-RULES-HOIST

- 准则版本：v1.65（步骤 0 时点；步骤 G 已升至 v1.66）
- 日期：2026-11-11
- 分类：治理类 `[SLICE]`；规模 L 档；探针额度 0（纯文档本机）
- 状态：**✅ 已完成**（2026-11-11 收口；停于 ⑩ 等待授权）
- 主控模型：`DeepSeek V4.1 Flash`（未切 Pro；按起跑消息【〇·零】5 不登记 ADR）
- 证据包：`docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/`
- 落点：准则新增 §13 五节（§13.1–§13.5）；v54 起快照只留时点数据

## 1. 执行摘要

一句话结论：v53 §七 / §八 / §九 中长期固化的执行纪律已按「全部迁移」方案结构性迁入准则 §13；v53 保持冻结，v54 起快照只保留时点数据。

关键读数：

- §13 五节条目数：13.1 = 20、13.2 = 26、13.3 = 25、13.4 = 26、13.5 三部分 8 / 4 / 5。
- 门禁坑 1–26 编号连续，症状与处置逐字保留（仅 3 处已登记改写，见 §5 与 §9）。
- CN/EN 对称：`git diff --numstat` = `184/3` 对 `184/3`。
- 门禁：三作用域 `TOTAL_FAIL=0`；`SELFTEST: PASS 722/722`。
- R 计数：R = 3（§9.2 口径；逐次构成见 §11）。

## 2. 三案对比表与最终落点

| 方案 | 主张 | ① 结论 | 理由 |
|---|---|---|---|
| A 全部迁移 | 源行数 98 全量纳入分析，合并后 82 条规则逐条处置（迁入 / 改指针 / 不迁各得其所） | **采纳** | 快照→准则对照完整、可审计；快照读者能在 §13 找到条款去向 |
| B 只迁 §七 §八 | 只迁起跑通则与协作纪律，清单元规则与门禁坑仍留快照 | 不采纳 | 门禁坑与清单元规则仍留在快照，快照膨胀未解 |
| C 只迁门禁坑 | 仅把坑表 1–26 移入准则 | 不采纳 | 其余纪律仍是软约束，与「升格为硬约束」的目标不符 |

① 结论：**同意 A**。最终落点 = 准则新增 §13 五节（§13.1–§13.5）+ v54 起快照精简（§七 / §八 / §九 改为一行指针 + 节名对照）。

## 3. 去重分析

规则级分母（引文，取自 `01-arch/dedup-analysis.md`）：本次共 82 条规则（源文件行数 98，合并后规则数 82）。14 处合并的对象是「同一条规则的多行描述」；非规则结构性行 2 个（坑表表头与版本注）不计入三分类。

| 分类 | 计数 | 处置口径 |
|---|---|---|
| 完全重复 | 19 | 改指针（§13 只留指针行，不复制实体） |
| 部分重叠 | 43 | 迁入（独立条目 + 指针） |
| 快照独有 | 20 | 迁入（全量正文） |

逐条表（82 条，含规则号 / 源位置 / 首句 / 对应节号 / 关系 / 备注）见 `01-arch/dedup-analysis.md` §1；本报告只放摘要表，不复制逐条行。

## 4. 迁移映射

| 落地状态 | 计数 | 构成 |
|---|---|---|
| 迁入 | 70 | 部分重叠 43 + 快照独有 20 + 门禁坑完全重复 7（逐字入坑表 + 指针括注） |
| 改指针 | 11 | 完全重复且非坑、纪律 Y 除外 |
| 不迁 | 1 | 纪律 Y 注记（已由 P59 升格为 §1.7.2，注记本身即指针） |
| 非规则行 | 2 | 坑表表头、版本注（不参与三分类） |

另补登 `G-0` 一条：§九 表前的引用总则句（不在 24 条表行内）→ `§13.3.1`，原样迁入。逐条映射与节号重指记录见 `01-arch/migration-mapping.md`。

## 5. 分步迁移逐项

| 步骤 | 内容 | `git diff --numstat`（CN = EN） | 落点 |
|---|---|---|---|
| A | §13 章框架（章导语 + 五个二级节空壳，零规则实体） | `24/0` | §13 章头与 §13.1–§13.5 标题 |
| B | v53 §七 20 条 | `21/0` | §13.1 二十条 |
| C | v53 §八 表行 + 纪律 Y 注记 + GLM 策略段 | `27/0` | §13.2 二十六条 |
| D | v53 §九 清单元 24 条表项 + 引用总则句 | `26/0` | §13.3 二十五条 |
| E | v53 §九 门禁坑 1–26 条 | `33/0` | §13.4 三列坑表 + 表前来源分组与版本注 |
| F | 快照结构模板 | `34/0` | §13.5 三部分（8 保留节 / 4 指针节 / 5 时点数据类） |
| G | 元数据 17 条：头部 1 + §0.3 1 + C.0at 13 + B.7.1 1 + B.7.4 1 | 元数据步骤 | 准则头部、§0.3 日志行、附录 C.0at、附录 B.7.1 / B.7.4 |

各步 numstat 取自该步核验文件（`03-test/step-*-check.md`）的实测值；步骤 G 只动元数据与槽位指针、不新增规则实体，其条目并入最终累计读数 `184/3` 对 `184/3`。

## 6. §13 新增章摘要

| 节 | 职责 | 条目数 |
|---|---|---|
| §13.1 起跑与步骤 0 通则 | 锚点核对、步骤 0 登记与起跑前置声明 | 20（13.1.1–13.1.20） |
| §13.2 跨片协作纪律 | 角色—载体—模型绑定、派单表机械读取、提交与临时件纪律 | 26（13.2.1–13.2.26） |
| §13.3 清单元规则 | 引用方式、重新生成时机与表述边界 | 25（13.3.1–13.3.25） |
| §13.4 门禁坑 | 事故登记册：编号 / 症状 / 处置，1–26 连续 | 26（含第 26 条 F2-FIX） |
| §13.5 快照结构模板 | v54 起快照的保留节、指针节与时点数据范围 | 三部分 8 / 4 / 5 |

## 7. §13.5 快照结构模板摘要

- 保留 8 节：§一 当前锚点、§二 事实包、§三 已完成切片、§四 执行队列、§五 观察项、§六 待用户裁定、§十 开场句、§十一 当前状态一句话。
- 指针 4 节：§七 → §13.1、§八 → §13.2、§九 清单元 → §13.3、§九 门禁坑 → §13.4（只留一行指针 + 节名对照，不复制实体）。
- 时点数据 5 类：锚点数 / 当期事实包 / 已完成切片清单 / 队列与观察项状态 / 规则文本（规则文本不入快照，只在 §13）。

## 8. 已知边界

- 规则升格后的语义变化——迁入 §13 的规则自下一片起为硬约束，不再仅作快照经验。
- §13.4 为事故登记册：症状与处置沿用事故原文；与正文或其他节同指者（第 11 条与 §13.1.12、第 21 条与 §9.3、第 8 条与 §13.1.19）以正文或既有条款为规范源，本表不另设第二个规范源。
- A-part：v53 §九 第 26 条候选已由 §13.4 正式化，P65 台账的「A-part 继承待办」由此落地。
- v53 冻结不追改（§11.1）：v50–v53 原样保留，本片不改快照正文。
- ⑥ 形式缺陷降级：第 3 轮标题未按 `Section A` 起头，按 §4.1 契约属格式偏离 ⇒ 结论降级为低权重旁证（见 §11）。
- EN 侧格式化补偿：坑表第 16 / 19 / 20 条为满足 G3(a.2) 六形状做了加粗补偿，与 EN 快照非逐字（CN 侧语义不受影响，已在 `03-test/step-E-check.md` 登记）。
- 指针精度留存：13.3.5 → §9.3、13.3.16 → §1.6 的指针措辞宽窄问题由 ⑤ 判为 Low 并 ACCEPT 留存，移交后续片收紧。
- 规则文本只活一处（§11.2）：v54 快照以指针替代 §七 / §八 / §九 正文。

## 9. 明确未执行事项

- 不改 v50–v53 历史快照（含 v53 §七 / §八 / §九 原文）；v54 不在本片生成。
- 不改判据本体 `tools/gates/**`、CJK 白名单、CONTRACTS / schema / fixtures、产品代码、CI。
- 不改准则 §0–§12 正文；仅 §0.3 新增日志行、头部版本行，以及 §13 内的两处收敛句（13.1.9 文档型退化指针、§13.4 登记册口径注记）。
- 不重跑 ⑥ 第 4 轮（额度用尽；恢复需后续新片按契约补做独立扫描）。
- 未 commit / push；未启动 VM / 代理 / Sol 双打；未新增探针。
- 未清理 `tmp/P66-CHECKLIST-RULES-HOIST/`（由 ⑧c 执行）；未新增临时件。

## 10. 阶段转录（sixstage 围栏 A–E）

围栏外注记：以下 A–E 五段为 ⑥ r3 的逐段转录，标题已按 §4.1 契约归一（原文标题为 `## A. Scope` 等）；⑥ 的效力按 §11 降级为旁证。

```sixstage
## Section A

I re-read the amended text in [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md) and [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md), together with §5.2, §9.3 and §11.2, and checked whether the earlier two Medium findings were actually resolved in substance and whether the two Low findings remained acceptable.

## Section B

- None.
- The two earlier Medium findings are materially closed.
- The two Low findings remain within an acceptable range.

## Section C

- In [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md), 13.1.9 states that documentation-type slices degrade to prfrail-documenter per §5.2, and that this clause "sets up no further exception." This is the substantive fix: the governance rule remains in §5.2, and no new exception is created.
- In [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md), 13.4 explicitly says the pitfall table is an incident register and that the body text or existing section remains the normative source under §11.2, with "no second normative source." That directly removes the earlier authority conflict.
- In [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md), 13.3.5 explicitly says "A rule obligation ≠ a gate obligation," and points to §6.3 and §9.3. That is an acceptable boundary and matches the low-item intent.
- In [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md), 13.3.16 says a [ROUTINE] small fix takes no new slice number but loses no discipline, and points to §1.6 for the ring count. That is a pointer fix, not a semantic rewrite, and it is consistent with the original low-item concern being addressed by precision rather than by creating new rule semantics.
- The English mirror in [docs/DELIVERY_DIRECTIVE_EN.md](docs/DELIVERY_DIRECTIVE_EN.md) says the same thing, preserving the authority chain and the no-second-source rule.
- The direct review confirms the 13.1–13.5 sections are present and that the pitfalls list remains a 26-item register, not a second authority layer.

## Section D

- The critical authority question is not whether a note says "no second source," but whether the authority chain is now actually fault-free. It is: §11.2 remains the source-of-truth rule, §5.2 remains the controlling role matrix, and 13.4 explicitly says the register is not a second normative source.
- The earlier Medium 1 is therefore substantively closed: the table is a register, not an authority stream.
- The earlier Medium 2 is also substantively closed: 13.1.9 does not create a new exception; it explicitly says the document-type role degradation remains under §5.2.
- The Low items remain acceptable because they are clarifying boundaries and pointers, not new gate meaning or new rule authority.

## Section E

- No files were modified and no commit or push was made.
- This was a read-only audit of the governing text, and the read evidence supports closure of the earlier Mediums without any remaining High/Medium findings.

INDEPENDENT SCAN: PASS
```

## 11. 审查与扫描记录

- ⑤ 预审：r1 报 2 项须修（F-1 节号重指记录缺位、F-2 引用总则句未入分母）⇒ 修复后 r2 判定 `RE-REVIEW: PASS`（F-1 / F-2 CLOSED；F-3 / F-4 / F-5 三项 Low ACCEPT 留存）。
- ⑥ 独立扫描实派 3 轮：r1 用自定义段落名（`### Scope` / `### Validation` / `### Findings` / `### Final result`）⇒ 格式偏离，作废；r2 为有效 A–E 形状但结论为 `FINDINGS`（2 项 Medium：13.1.9 与 §5.2 字面冲突、§13.4 规范源二义），经 ③-r2 的收敛句核实为已闭合的陈旧读数；r3 结论 PASS，但标题写作 `## A. Scope` 而非 `## Section A`，按 §4.1 契约仍属格式偏离 ⇒ ⑥ 额度用尽 + 形式缺陷，其结论降级为低权重旁证。
- ⑦ 终审 r1：6 项判定（5 PASS + 1 ISSUE），发现清单含 [High] F-01（⑥ 的 PASS 不能作为有效通过结论），结论行 `FINAL REVIEW: FINDINGS`；整改项 = 本报告与台账如实登记「⑥ 额度用尽 + 形式缺陷」并把 ⑥ 结论降级为旁证。
- 闭环口径：以 ⑦ 结论与 ④ 门禁为主，⑥ 仅作低权重旁证（F-01 的整改项即本节与 §12 的登记）。
- R 计数（§9.2）：③-r1 → r2 修正 1 次、⑥-r2 的 FINDINGS → r3 修正 1 次、⑦-r1 的 F-01 整改 1 次 ⇒ R = 3；⑥-r1 作废为重派，不计入 R。
- ⑧ 装配后复核：改动仅涉报告与台账，规则与结构类判据（G1 / G3 / G4b / G6-4 / G7-e / G7-f）全绿；只余 G5(d) 围栏补行与 G8-b 的 MANIFEST / SHA256SUMS 重算两项，由主控在 ⑨ 前重同步。

### 11.1 观察项候选（P66 实证，非本片范围）——进 v54 事实包

- **⑥ 独立扫描输出合规率**：本片三轮 ⑥ 均存在形式缺陷（r1 未用 A–E 形状 ⇒ 作废；r2 用 A–E 但 2 项 Medium 属陈旧读数；r3 用 `## A. Scope` 而非 `## Section A`），额度用尽后被降为「低权重旁证」。
- **影响**：⑥ 作为独立扫描层的有效性在本片失效；⑦ 两次兜住（终审 FINDINGS + 复审 PASS），但 ⑥ 的「独立信号」未发挥。
- **候选处置**：v54 或后续治理片评估「⑥ 提示词模板是否需要加固」或「⑥ 格式偏离是否应触发自动重派而非计入额度」。
- **本片立场**：不处理，仅在 v54 事实包登记（用户 ⑩ 裁决 1）。

## 12. 越界与自由度登记

| 编号 | 项 | 是否触发 | 依据 / 处置 |
|---|---|---|---|
| D-1 | 步骤 A 由主控自执行（未派 `prfrail-documenter`） | 触发 | 纯骨架、零规则实体；独立影响 = 低；B–G 不延续（`03-test/step-A-check.md`） |
| D-2 | ③-r2 两处措辞级修正由主控自执行 | 触发 | 各 1 句 / 1 条注记，不新增规则实体；独立影响 = 低；⑤ 已独立复核（`03-test/consistency-r2.md`） |
| F-1 | 不自行 commit / push | 未触发 | 停于 ⑩ 等同一轮授权 |
| F-2 | 不启动新切片、不替用户裁定 | 未触发 | 三案由 ① 给建议、用户裁定走 A |
| F-3 | 不开 VM / 代理 / Sol 双打 | 未触发 | 纯文档本机执行 |
| F-4 | 不改判据本体与契约 | 未触发 | G4b 命中后改措辞，白名单零改动 |
| F-5 | 不改 §0–§12 正文 | 未触发 | 仅 §0.3 日志行与头部版本行 |
| F-6 | 不改 v50–v53 | 未触发 | v53 保持冻结 |
| F-7 | 门禁坑只改结构与位置 | 未触发 | 3 处已登记改写，症状与处置逐字保留 |
| F-8 | 模型绑定与串行派单 | 依令执行 | 每单同传 `agentName` 与 `model`，严格串行 |
| F-9 | 额度 ⑤ 1+1 / ⑥ 3 / ⑦ 1+3 | 用满 | ⑤ 两轮、⑥ 三轮、⑦ 一轮 |
| F-10 | 单次输出上限 | 依令分批 | 批量辅助回执报实得字符 |
| F-11 | 活台账只读、不主动改 | 本环按 ⑧ 授权例外 | 只改 P66 段，其余段只读 |
| F-12 | 步骤 0 围栏同轮声明 | 已执行 | 报告骨架与证据包同轮落 `artifacts` 围栏 |
| F-13 | 骨架期的围栏豁免声明以真实围栏替换 | 已执行 | 本环删除豁免行并落 `sixstage` 围栏 |
| 命中 1 | G4b 新增 CJK 单字 | 已处置 | §13.3 导语改「表述边界」，白名单零改动 |
| 命中 2 | G6-4 聚合计数措辞 | 已处置 | §13.5 脚注改「P48 起累计至第 26 条」 |

两类命中合计 2 次（G4b 一次、G6-4 一次），均以改措辞收口，未改判据本体。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/checklist-rules-hoist.md	repo
docs/validation/checklist-rules-hoist_EN.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/00-briefing.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/01-arch/decision-list.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/01-arch/dedup-analysis.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/01-arch/design-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/01-arch/migration-mapping.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/02-batch/glm-dedup-r1.txt	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/consistency-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/consistency-r2.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-A-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-B-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-C-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-D-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-E-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/03-test/step-F-check.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/05-prereview-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/05-prereview-r2.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/06-scan-r3.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/04-reviews/07-rereview-r1.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/06-gates/gates-r1.txt	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/06-gates/gates-r2.txt	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/MANIFEST.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/SHA256SUMS.txt	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/freedom-list.md	repo
docs/validation/evidence/P66-CHECKLIST-RULES-HOIST/input-list.md	repo
```
