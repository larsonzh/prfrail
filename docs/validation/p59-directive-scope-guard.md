# P59 DIRECTIVE-SCOPE-GUARD 防越界防线（验证报告）

**状态：🔄 进行中**（阶段 B 完成；⑩ 停点等授权；未 commit/push）。

- 切片：`DIRECTIVE-SCOPE-GUARD`（治理类 `[SLICE]`，层 2 的 `G-3` 项）
- 起跑：2026-11-06；锚点 HEAD = `33bd849`；准则 v1.62 → v1.63（R-12 裁定授权，非切片内自行决定）
- 主控模型：`DeepSeek V4 Pro`（ADR-032，步骤 0 登记）
- 裁定：R-12 = 执行侧 A2 × 审查侧 B1 × 纪律 Y Y1；规模 S
- 两段式：阶段 A = 步骤 0 → ① → 停等 R-12；阶段 B = ②→⑩（已按裁定完成）

## 1. 异常登记区

- ⑦ 首轮结论行取值 `CHANGES_REQUESTED` 不在 §4.1 闭集（格式偏离）；按 v48 §7.5 的同 stage 重发计入上限口径：首轮占 1 次终审、r2 占 1 次复审，合计 1 终审 + 1 复审，均在 1+3 额度内；2 项实质发现已整改并经 r2 `RE-REVIEW: PASS` 闭合。
- 盲审口径更正（用户 ⑩ 裁定）：⑦ 首轮不附清单系 §7.3 锚定硬隔离；§7.7 盲审必抽条件（所有权 / 停机 / 身份）本片未触发 ⇒ 未单独抽盲审；起跑消息「必抽盲审」表述超出 §7.7 字面，以 §7.7 为准（如实留痕，不追责）。
- F2 处置（R-12 条目翻「已裁」+ P58 注记）超出 R-12 七交付物清单（⑤ r1 要求）；已在 ⑩ 回执向用户声明。
- 派单表未机械校验：本会话 debug 日志不含 transcript（`deep-reasoner` 命中 0，与 P56/P58 同型）；派单表以主控工具调用参数为事实。
- ⑤ r2 新发现 N1–N3（CHECKLIST §五 OB-137 观察清单 / G-3 行内含「最终形态由 R-12 裁定」/ v48 锚点陈旧）移交 ⑩ 事实包，本片受裁定范围约束不再加改。
- sixstage 转录格式适配（主控执行，⑥ 归档原文 `04-reviews/06-scan-r1.md` 保持逐字；本条含引文）：① 五段标题全角括注改空格（`## Section A（…）` → `## Section A …`，适配 G7-f 的标题正则）；② Section A 结论句的聚合计数措辞「共 17 个路径」改「覆盖 17 个路径」（适配 G6-4 的「共 NN」聚合计数正则）。白名单文件 `tools/gates/g6-whitelist.txt` 未动（本片边界不含 `tools/gates/**`）。

## 2. ① 三案对比表

摘要（全文见 `docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/01-arch/design-r1.md`）：

| 面 | 案 | 一句话 | 成本 | 新增机械判据 |
|---|---|---|---|---|
| 执行侧 | A1 记录态约束 | §1.7.2 段末加记录态条款，靠 ⑤/⑦ 人工核验 | S | 否 |
| 执行侧 | A2 正文条款 | §1.7.2 升级边界列表 +3 条 + §0.2 一句 | S | 否 |
| 执行侧 | A3 判据 | §6.3 新 G10/G6-9（升版形态 + 引用形态；授权真伪与超额派单不可机械） | M | 是 |
| 审查侧 | B1 ⑤ 第 7 条 | §4.2.1 增设「授权可回溯性」必核动作 | S | 否 |
| 审查侧 | B2 判据 | §6.3 G7-g（降级为引用形态判据） | M | 是 |
| 审查侧 | B3 记录态 | 仅 §12.4 OB-137 处置列加注 | S | 否 |
| 纪律 Y | Y1 升正文 / Y2 保留候选 / Y3 并入 A1 | 见 design-r1.md §六 | S | 否 |

① 推荐（仅供参考，最终由用户裁定）：执行侧 A2、审查侧 B1、纪律 Y Y1；规模 S。

**R-12 已于 2026-11-06 裁定：A2×B1×Y1；本片按裁定进入阶段 B 并落地（v1.63）。**

## 3. 阶段 B 执行摘要

- A2：§1.7.2 升级边界列表末新增三条 bullet（未经同轮授权的准则升版 / 授权依据表述须可回溯 / 超出授权数量或范围的派单；P58 的 D-6 为两次重复型实证）；§0.2 增「准则升版属 §1.7.2 升级边界」指针句。
- B1：§4.2.1 最低核验动作六条 → 七条（第 7 条 = 授权可回溯性）；§12.4 OB-137 处置列翻 `已处置` 并补落地注记。
- Y1：纪律 Y 升为正文条款（正文由 A2 三条承载）；v48 §八纪律 Y 行只留指针「已升格至 §1.7.2（v1.63）」；G-3 行标「由 P59 闭环」。
- 版本：v1.62 → v1.63（头部 / §0.3 的 v1.63 条目，均标注「用户 2026-11-06 的 R-12 裁定授权」）。

## 4. 审查链

| 环 | 结论 |
|---|---|
| ① | `DESIGN: DONE`（三案对比表，停等 R-12） |
| ⑤ | r1 `PRE-REVIEW: FINDINGS`（2 Medium + 4 Low）→ 整改 → r2 `PRE-REVIEW: PASS`（F1–F6 CLOSED；N1–N3 移交） |
| ⑥ | `INDEPENDENT SCAN: PASS`（0 发现；sixstage 围栏已转录） |
| ⑦ | 首轮格式偏离（占 1 次终审）→ r2 `RE-REVIEW: PASS`（占 1 次复审；合计 1+1，额度 1+3 内）；未单独抽盲审（§7.7 条件未触发；首轮不附清单系 §7.3 锚定硬隔离） |

## 5. 文档同步（【四】15 项逐项）

1. `DELIVERY_DIRECTIVE{,_EN}.md`：已改（头部 v1.63 / §0.2 / §1.7.2 / §4.2.1 / §12.4 OB-137 / §0.3 新增行）。
2. `REMAINING_SLICES{,_EN}.md`：P59 段转 ✅ + 完成记录行 + 状态表行同步。
3. 报告对：本文件。
4. 证据包：MANIFEST / SHA256SUMS / 00-briefing / design-r1 / ruling / ⑤⑥⑦ 原文 / freedom-list / input-list / 06-gates 读数。
5. `ADR_REGISTER{,_EN}.md`：ADR-032（步骤 0）；ADR-033 待 ⑩ 用户裁定。
6. `DEV_PLAN`：不适用（治理片不回写 `DEV_PLAN`，按 §11.2）。
7. `ARCHITECTURE`：不适用（不涉及模块边界 / 数据流 / 存储）。
8. `SECURITY`：不适用（不涉及权限模型 / 威胁矩阵 / 秘密处理）。
9. `CONTRACTS`：不适用（不涉及 wire / schema / 状态机 / 判据语义；本片未落判据）。
10. `S1_SUPPORT_MATRIX` 与 `S1_RELEASE_NOTES`：不适用（无平台 / 支持层级 / EOL / 发布变更）。
11. `TEST_STRATEGY`：不适用（本片未落判据，不在 AT 体系）。
12. `CHECKLIST_v48`：纪律 Y 行标「已升格」（R-12 交付物）；R-12 条目翻「已裁」（超出清单，已声明）；不触发重新生成。
13. 其余产品文档（PRODUCT_REQUIREMENTS / BUSINESS_WORKFLOWS / RFC / DOCUMENTATION_PLAN / CODING_CONVENTIONS / OPERATIONS / INSTALLATION / INSTALLATION_PLAN）：不适用（只触执行纪律层）。
14. `SLICE-CLOSEOUT-TEMPLATE`：不适用（本片不是结项报告）。
15. 报告本体：artifacts 围栏（TAB 分隔、CN/EN 逐字节相同）+ sixstage 围栏。

## 6. 计量与派单

- R = 0；临时件 D = 0。
- 派单：①×1、⑤×2（1+1 额度内）、⑥×1（额度 1/3）、⑦×2（首轮 1 次终审 + r2 1 次复审；额度 1+3 内）；批量辅助 0（本片无枚举/摘录类任务）。
- 探针与付费外部调用：0 次真实外部调用。

## 7. 本片自身越界登记（D-6 类）

- 未自行升版（v1.63 系 R-12 裁定授权，§0.3 行已如实标注）。
- 未伪造授权（四处授权表述全部可回溯至 ruling 原文，⑤ r1 已核）。
- 未超额派单（⑤ 1+1、⑥ 3、⑦ 1+3 均在额度内）。

## 8. A-NN 编号完整性与偏离登记

| 编号 | 主体 | 偏离 | 处置 |
|---|---|---|---|
| D-1 | ⑦ 首轮 | 结论行取值 `CHANGES_REQUESTED` 表外（格式偏离） | 同 stage 重发计入上限：1 终审 + 1 复审（额度内）；实质发现整改后 r2 PASS |
| D-2 | 主控 | F2 处置（R-12 条目翻「已裁」）超出 R-12 七交付物清单 | ⑩ 回执声明 |
| D-3 | 主控 | ⑦ 输入 diff 生成早于围栏补行，致首轮 High 误报 | 已重新生成 diff 并复审 PASS |
| D-5 | 主控 | ⑩ 回执称「盲审已抽」实为 §7.3 锚定硬隔离；§7.7 盲审条件本片未触发 | 按用户 ⑩ 裁定更正口径并留痕，不追责 |

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/CHECKLIST_v48.md	repo
docs/t027/CHECKLIST_v48_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/00-briefing.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/01-arch/design-r1.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/01-arch/ruling-2026-11-06-r12.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/04-reviews/07-final-r2.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/06-gates/gate-readings.txt	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/MANIFEST.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/SHA256SUMS.txt	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/freedom-list.md	repo
docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/input-list.md	repo
docs/validation/p59-directive-scope-guard.md	repo
docs/validation/p59-directive-scope-guard_EN.md	repo
```

## ⑥ 独立扫描 sixstage 转录（逐字；含 2 处格式适配以过 G6-4 / G7-f 机械判据，适配明细见 §1 异常登记区；⑥ 归档原文 `04-reviews/06-scan-r1.md` 保持逐字原件）

```sixstage
## Section A 扫描范围与输入

已核对输入与边界口径，范围如下：

- 变更输入： [tmp/P59-DIRECTIVE-SCOPE-GUARD/06-scan-input.diff](tmp/P59-DIRECTIVE-SCOPE-GUARD/06-scan-input.diff)
- 裁定原文： [docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/01-arch/ruling-2026-11-06-r12.md](docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/01-arch/ruling-2026-11-06-r12.md#L1-L15)
- 起跑消息： [docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/00-briefing.md](docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/00-briefing.md#L7-L15)
- 契约原句： [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L39-L40), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L228-L241), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L452-L468), [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L1169-L1186)
- 约束对象： [docs/validation/p58-ob-triage-assessment.md](docs/validation/p58-ob-triage-assessment.md)

结论：本次 diff 覆盖 17 个路径，全部为文档/台账/报告/证据包文件，不含 [tools/gates](tools/gates)、[internal](internal)、[docs/CONTRACTS.md](docs/CONTRACTS.md)、[schemas](schemas) 及任何产品代码。核对后，变更落点符合 R-12 裁定的 A2 × B1 × Y1 组合，未顺手扩到 A3/B2，也未触碰判据实现或代码执行面。

## Section B 发现清单

无发现。

- 发现数：0
- 严重度：N/A
- 修法：不需要修正；当前 diff 与裁定边界一致，且未越过用户授权范围。

## Section C 对 R-12 裁定边界与契约原句的逐项对照

1. A2 执行侧：一致
   - R-12 明确要求在 §1.7.2 末新增 3 条 bullet，并在 §0.2 补一句；对应落地见 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L228-L241) 和 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L39-L40)。
   - 这与 [docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/01-arch/ruling-2026-11-06-r12.md](docs/validation/evidence/P59-DIRECTIVE-SCOPE-GUARD/01-arch/ruling-2026-11-06-r12.md#L5-L15) 的“执行侧 = A2”描述一致。

2. B1 审查侧：一致
   - §4.2.1 现已从六条增至七条，新增第 7 条“授权可回溯性”，位置与口径见 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L452-L468)。
   - R-12 要求“审查侧 = B1”，且将 §12.4 的 OB-137 处置列翻为已处置；这与 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L1169-L1186) 一致。

3. Y1 纪律 Y：一致
   - Y1 要求“升正文条款，但防重复描述”，且“v48 §八 只留指针”。对应落地见 [docs/t027/CHECKLIST_v48.md](docs/t027/CHECKLIST_v48.md#L167-L168)。
   - 这与 R-12 的“防重复描述”要求一致：正文在 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L228-L241)，而 v48 只保留指针，不重复正文。

4. OB-137 处置与授权回溯：一致
   - §12.4 的 OB-137 现已标为“已处置”，并追加说明，且与 B1 的授权可回溯动作相衔接，见 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L1169-L1186)。
   - R-12 明确写明：OB-137 由 B1 处置，并让 A2 承担执行侧；两者逻辑一致，无越界扩张。

5. v1.63 的正当升版：一致
   - 头部版本号和 §0.3 日志行均已写入 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L1-L3) 与 [docs/DELIVERY_DIRECTIVE.md](docs/DELIVERY_DIRECTIVE.md#L109-L110)。
   - 这与裁定原文“本片的 v1.63 升版是正当升版，由 R-12 授权”一致，且与 “非切片内自行决定” 的说明相符。

## Section D UNCOVERED

UNCOVERED: 本次只读扫描未审查任何运行时行为或 CI 结果，因为本 diff 未触碰 [tools/gates](tools/gates) 与代码路径，且本角色被限定为只读，不能做执行性验证。

UNCOVERED: 本次扫描没有覆盖未来若另开 A3/B2 判据片时的实现层面；这里仅验证当前 R-12 明确授权的 A2 × B1 × Y1 范围是否被遵守。

## Section E 结论与建议

本次扫描未发现越界、重复、重叠或误扩范围的问题。diff 封闭在 A2/B1/Y1 组合与对应证据包、报告、台账更新内，且与 R-12 原文、起跑消息、§1.7.2/§4.2.1/§12.4 的契约要求保持一致。建议接受并结束本轮范围核查，不再扩展到 A3/B2 或判据代码路径。

INDEPENDENT SCAN: PASS
```
