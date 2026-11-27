# P97 对照评审记录（P87 基准侧 + Hydrogen 实测侧）

切片：P97 `HYDROGEN-CONTROL-PILOT`｜日期：2026-11-28｜对照对象：P87 证据包与准则 v1.81 的 P87 生效面｜性质：只记录、不重裁

## 1. 对照对象与输入包

对照对象为 P87 证据包 `docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/` 与准则 v1.81 的 P87 生效面。

输入缺口声明：`tmp/P87-DIRECTIVE-REVIEW-MODEL-POLICY/review-input/`（`diff.patch`、`fix-summary-r2.md`）已按 §1.5 清理不可得；本读以「落地形态 + 证据包」替代变更集 diff，与 P87 原 ⑦ 读**不等权**，不可比处一律标「不可比」。用户裁定：输入包若超 Hydrogen 单次输入上限，须先报用户收窄方案（P87 生效面优先），不得自行裁剪。

## 2. 四面定义与判读规则

| 对照面 | 记录口径 | 判读规则 |
|---|---|---|
| 结论 | 终审结论行原文 | 闭集：`FINAL REVIEW: PASS` / `FINAL REVIEW: FINDINGS` |
| 发现数 | Section B 计数 | Medium+ 与 Low 分计；UNCOVERED 不计入 |
| 严重度分布 | Critical/High/Medium/Low 四档 | 逐档给出条数 |
| 误判偏离面 | 与 P87 盲审同项比对 | 命中 P87 盲审那条 Medium 记「重合」；报出新 Medium+ 记「偏离」 |

## 3. P87 基准侧（逐字转录 + 证据锚点）

- ⑦ 终审 r1：结论 `FINAL REVIEW: PASS`；Section C 明确「未发现新增 Critical/High/Medium」。证据：`docs/validation/evidence/P87-DIRECTIVE-REVIEW-MODEL-POLICY/04-reviews/07-final-review-r1.md`。
- ⑦ 盲审 1 次：`FINDINGS`（1 Medium：备选额度口径与起跑消息双口径），已闭合，属终审前的发现。
- ⑤ 与 ⑥ 的发现不计入 ⑦ 基准，仅作背景：⑤ r1 6 Medium → 修复 → r2 1 Medium + 4 Low（额度用尽，交 ⑥⑦ 复核）；⑥ r1 1 Medium + 2 Low（`FINDINGS`）→ r2 `PASS`。
- 基准面取值：结论 = PASS；新增 Medium+ = 0；新增 Low = 0（据 P87 终审 Section C「未发现新增 Critical/High/Medium」与 Section D 全文无新发现行推出；P97 ⑤ 预审 r1 的 Low-4 要求补此锚点）。
- 记录节：`docs/validation/directive-review-model-policy.md`。

## 4. Hydrogen 侧（待试水填入）

| 待填字段 | 取值 |
|---|---|
| 结论行 | `FINAL REVIEW: FINDINGS`（Section A 首行逐字） |
| Section B 发现数 | Medium+ = 3（High 1 + Medium 2）；Low = 2；合计 5（UNCOVERED 不计入） |
| 严重度分布 | Critical 0 / High 1 / Medium 2 / Low 2 |
| 是否命中 P87 盲审同项 | 命中：第 4 条（§7.5 备选额度口径与起跑口径偏离）与 P87 盲审那条 Medium 同标的；等级不同（Low vs Medium） |
| 是否新 Medium+ | 3 条 —— #1 High（§5.2 深度档档位列 `Max` 与 §8.2 `gemini-4-neon` = Extra High 冲突）；#2 Medium（§5.2 矩阵字面可判性与档位互斥）；#3 Medium（§7.5 `gemini-4-neon` 两档上下文缺另一档单价） |
| 覆盖声明要点 | 目标五节（§2.2 / §5.2 / §5.4 / §6.3 G5-b / §7.5 / §8.2）CN 与 EN 均无截断完整读取；历史超长行在 2000 字符处被工具截断（非目标节）。注：原文 Section E 自述「目标五节」但实际列举 6 个条款号，本行按实际列举逐项转录、不改归档原文 |
| UNCOVERED 要点 | P87 改后 SELFTEST 输出件未入包；`--scope=ci` 的 `changed=2` 为上一提交而非 P87 自身；P87 报告对与 ADR-048 全文未入包 |
| 无效调用判定 | 有效（五段齐备，结论行逐字命中闭集） |

## 5. 对照总表

| # | 对照面 | P87 基准（证据锚点 + 逐字） | Hydrogen 实测（证据锚点 + 逐字） | 一致 / 偏离 / 不可比 | 备注（不等权标注） |
|---|---|---|---|---|---|
| 1 | 结论 | `07-final-review-r1.md`：`FINAL REVIEW: PASS` | `07-pilot-hydrogen.txt` Section A 首行：`FINAL REVIEW: FINDINGS` | 偏离 | 输入不等权（落地形态 vs 变更集 diff） |
| 2 | 发现数 | Section C：未发现新增 Critical/High/Medium | Section B：Medium+ = 3（High 1 + Medium 2）、Low = 2 | 偏离 | 同上 |
| 3 | 严重度分布 | 新增 Medium = 0；新增 Low = 0 | Critical 0 / High 1 / Medium 2 / Low 2 | 偏离 | 同上 |
| 4 | 误判偏离面 | 盲审 1 Medium（额度双口径）已闭合 | 第 4 条与盲审同标的、等级 Low；另 3 条为新发现 | 重合（等级不同）＋新发现 | 不重裁：P87 `PASS` 保持不变 |

## 6. 记录纪律（只记录不重裁）

1. 不改 P87 任何已闭环结论与报告 / 台账 / 版本行。
2. Hydrogen 新发现只登记为 P97 观察项，处置权归用户。
3. 措辞上限为「单次对照实测完成」；禁止写「轻量档已验证」「优于 P87」等外推表述（n=1）。

## 7. 观察项

O-1（用户 2026-11-28 裁定）：P96 额度外收窄复审的留痕位不足 §1.7.4 ② 字面——`freedom-list` / 台账有，§12.4 无。不追改历史，仅记录。

O-2（Hydrogen 第 1 条，High）：§5.2 深度档档位列写 `Max`，与 §8.2 的 `gemini-4-neon` = `Extra High` 冲突，并连带 §8.2「任何 `Max` 取档仍须同轮授权」的授权口径歧义。**P87 原 ⑦ 终审 `PASS` 未检出该矛盾，轻量档检出**；P87 结论不重裁，本项登记为观察项、处置权归用户。

O-3（Hydrogen 第 2 条，Medium）：§5.2 选型矩阵字面可判性不足（`定稿复核` / `定向小读` 无定义、`纯文档 / S 档` 连接词与档位优先级缺失、`其余三档` 计数差、Sol 备选与 Luna 备选缺客观门槛、`最高风险轮` 无界线）。不重裁 P87，登记为观察项、处置权归用户。

O-4（Hydrogen 第 3 条，Medium）：§7.5 `gemini-4-neon` 分 200K / 1M 两档但未标注主表单价所属档位、另一档单价缺失。不重裁 P87，登记为观察项、处置权归用户。

O-5（Hydrogen 第 4 / 5 条，Low）：§7.5「备选额度口径」措辞易误读（与 P87 盲审同标的）；`numstat` 证据件自身行数未最后回采。均登记为观察项、处置权归用户。
