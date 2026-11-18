# P80-SG6-GATELEVEL-CONSOLIDATE - ⑦ 终审 r1（只读）

日期：2026-11-18
范围：仅复核已暂存快照与证据包；除本文件外未修改任何文件。

## 两项确认

1. SKIP CONFIRMED（§7.2.1 三条件）
- 依据 1：报告已显式声明「⑥ 跳过：§7.2.1 三条件全部成立」（docs/validation/sg6-gatelevel-consolidate.md:163）。
- 依据 2：整改快照文件集为 17/17 文档路径（tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.numstat:1-17），未见产品代码、脚本语义或契约语义改动载体。
- 结论：成立。

2. BLIND REVIEW CONFIRMED（盲审文件存在且发现已处置）
- 依据 1：盲审文件存在且结论行为 `BLIND REVIEW: FINDINGS`（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/blind-r1.md:44）。
- 依据 2：主报告 §13.1 已登记盲审轮次与处置口径（docs/validation/sg6-gatelevel-consolidate.md:166-171）。
- 依据 3：盲审与预审文件已进入 `artifacts` 围栏（docs/validation/sg6-gatelevel-consolidate.md:202-203）。
- 结论：成立。

## 十项逐条裁定

1. 改动范围最小化（仅本片载体，无越权）
- 裁定：PASS。
- 证据：暂存集 17 文件均在本片报告/台账/证据包范围内（tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.numstat:1-17）。

2. 门级判定与判据一致（“门级侧成立”且明示可证伪条件）
- 裁定：PASS。
- 证据：§5 给出「S-G6 门级侧成立」并列出 4 条可证伪条件（docs/validation/sg6-gatelevel-consolidate.md:61-67）。

3. 来源口径一致（p57 / s1-exit / v62 / ADR-046 / P76 Q8(c)）
- 裁定：PASS。
- 证据：§5、§6 给出来源落点与逐字引证，且声明 p57 与 s1-exit 门表零改动（docs/validation/sg6-gatelevel-consolidate.md:52-83）。

4. 并陈口径准确（门级侧成立 ≠ 判据通过 ≠ T018 完成）
- 裁定：PASS。
- 证据：§6 明确并陈句、依赖字段与不勾选状态保持（docs/validation/sg6-gatelevel-consolidate.md:72-84）。

5. 下一步清单与 ADR 约束一致（S-G5 两步、ADR-047 到期）
- 裁定：PASS。
- 证据：§7 明确 S-G5 第二步前置与 ADR-047 到期约束（docs/validation/sg6-gatelevel-consolidate.md:87-95）。

6. 对称性（CN/EN 报告同构，台账双侧条数一致）
- 裁定：PASS。
- 证据 1：报告对均为 +214 行（tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.numstat:16-17）。
- 证据 2：门禁记录 G3(a) 对报告对与台账对给出对称通过（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-tree.txt:10）。

7. 门禁读数一致性（报告 §12 vs 06-gates 原始输出）
- 裁定：FINDING（Medium）。
- 发现：`06-gates` 中的快照口径与当前整改快照不一致，存在证据时点漂移：
  - 当前整改快照为 17 文件（tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.numstat:1-17）。
  - 但 `gate-index.txt` 与 `gate-tree.txt` 均为 `changed=15`（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-index.txt:1，docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-tree.txt:1）。
  - `gate-index.txt` 记录 `G7-b ... 30 artifact row(s)`，与当前双报告 `artifacts` 各 17 行（共 34 行）不一致（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-index.txt:34；docs/validation/sg6-gatelevel-consolidate.md:196-214；docs/validation/sg6-gatelevel-consolidate_EN.md:196-214）。
  - `gate-ci.txt` 头部为 `changed=2` 且内容聚焦 `CHECKLIST_v62`，显示与本片快照混域（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-ci.txt:1,10）。
- 影响：§12 的“见⑨实测”虽可指向绿读数，但 `06-gates` 归档与当前最终快照不完全同批，证据链完整性不足。

8. 前轮发现闭合（⑤ Low/Info 与盲审 Medium/Low）
- 裁定：PASS。
- 证据：
  - ⑤ 预审结论与计数存在（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/prereview-r1.md:70）。
  - ⑦ 盲审结论与计数存在（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/blind-r1.md:44）。
  - 主报告已登记前轮处置闭环并将盲审/预审纳入围栏（docs/validation/sg6-gatelevel-consolidate.md:170-171,202-203）。

9. 额度口径（S 档 1+1 + 盲审另计）
- 裁定：PASS。
- 证据：报告抬头已按 §7.5/§7.7 声明（docs/validation/sg6-gatelevel-consolidate.md:10）。

10. 格式合规（artifacts 围栏与暂存集双向一致、无伪造路径）
- 裁定：PASS（就当前文档集合逐项比对）。
- 证据：
  - 暂存集 17 路径（tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.numstat:1-17）。
  - CN/EN 报告 `artifacts` 各 17 路径，路径均可在仓库中定位（docs/validation/sg6-gatelevel-consolidate.md:196-214；docs/validation/sg6-gatelevel-consolidate_EN.md:196-214）。
  - 双向比对结论：stagedNotDeclared = 0；declaredButMissing = 0。

## 严重度计数

- High: 0
- Medium: 1
- Low: 0
- Info: 0

## 是否需要 r2

需要。触发条件：r1 存在 Medium（第 7 项）。

## r2 最小修复目标

1. 在当前 17 文件快照上重跑并回填 `06-gates`（index/tree/ci/selftest），确保同批一致。
2. 使门禁归档与当前快照对齐：`changed` 口径、G7-b artifact 行数、scope=ci 作用域对象一致。
3. 回填后复核 §12 与 `06-gates` 的同批可追溯性，再做 r2 结论。

FINAL REVIEW: FINDINGS
