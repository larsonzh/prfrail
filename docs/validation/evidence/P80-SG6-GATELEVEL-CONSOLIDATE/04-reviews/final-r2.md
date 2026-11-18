# P80-SG6-GATELEVEL-CONSOLIDATE - ⑦ 复审 r2（只读）

retryOf = final-r1
基线：`tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.diff` + `tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.numstat`（整改后快照）+ `docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/`

## 两项确认

1. SKIP CONFIRMED：成立。
- 依据：报告已声明「⑥ 跳过：§7.2.1 三条件全部成立」（docs/validation/sg6-gatelevel-consolidate.md:163）。
- 依据：暂存变更仅落在台账、报告与本片证据包（tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.numstat:1-18），未触及代码、CI、CONTRACTS、schema、fixtures 与门表源文档。
- 对照：§7.2.1 的三条件与「由 ⑦ 或用户确认」要求仍在（docs/DELIVERY_DIRECTIVE.md:713-717）。

2. BLIND REVIEW CONFIRMED：成立。
- 依据：盲审原文存在且结论行为 `BLIND REVIEW: FINDINGS`（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/04-reviews/blind-r1.md:1）。
- 依据：主报告已登记盲审轮次与结论（docs/validation/sg6-gatelevel-consolidate.md:171），并将盲审文件纳入 `artifacts` 围栏（docs/validation/sg6-gatelevel-consolidate.md:202）。

## r1 Medium 闭合判定

1. 同批快照一致性（`changed` 口径）
- r1 要求：`06-gates` 与整改后快照同批。
- 复核：CLOSED。
- 证据：`impl.numstat` 为 18 项（tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.numstat:1-18）；`gate-index/tree` 均为 `changed=18`（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-index.txt:1, docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-tree.txt:1）。

2. `G7-b` 行数与双报告围栏一致
- r1 要求：`G7-b` 与双报告 `artifacts` 行集合一致。
- 复核：CLOSED。
- 证据：`gate-index/tree` 均为 `36 artifact row(s)`（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-index.txt:34, docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-tree.txt:34）；CN/EN 两份报告围栏均与暂存集双向一致（stagedNotDeclared=EMPTY, declaredButMissing=EMPTY；基于 `impl.numstat` 与两份报告 `artifacts` 的独立比对）。

3. `--scope=ci` 本地退化口径是否如实说明
- r1 要求：在报告 §12 显式补充「读数作用域」说明，含 `--scope=ci` 本地退化口径。
- 复核：NOT CLOSED。
- 证据：`gate-ci` 明确为 `changed=2` 且对象是 `CHECKLIST_v62`（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-ci.txt:1,10）；报告 §12 仅写「`--scope=ci` 见⑨实测」与通用读数口径（docs/validation/sg6-gatelevel-consolidate.md:155,158），未显式声明该 `ci` 读数是本地退化口径（非本片快照同批判据）。

## 回归复核结果

1. 引证逐字抽检：通过（抽检项不少于 6 处，均与来源一致）。
- `p57` 判据原句「依赖字段满足且四道剩余门闭合」与报告 §5 引述一致（docs/validation/p57-closeout-assessment.md:93；docs/validation/sg6-gatelevel-consolidate.md:52）。
- `s1-exit` 的「T027 未完成、AT-23 整体未通过」与报告 §6 引述一致（docs/validation/s1-exit.md:52；docs/validation/sg6-gatelevel-consolidate.md:74）。
- v62 的「发行部分已批准（第一步）」与「门级侧可收口」与报告 §5/§7 口径一致（docs/t027/CHECKLIST_v62.md:119-120,196；docs/validation/sg6-gatelevel-consolidate.md:57,90）。
- P76 Q8(c) 并陈口径与报告 §6 一致（docs/validation/s1-exit-closeout-assessment.md:96-97；docs/validation/sg6-gatelevel-consolidate.md:76-77）。
- `DEV_PLAN` 依赖字段与状态字段引用一致（docs/DEV_PLAN.md:94,96,98；docs/validation/sg6-gatelevel-consolidate.md:61,74）。
- ADR 状态与到期约束引用一致（docs/ADR_REGISTER.md:19-20,57-58；docs/validation/sg6-gatelevel-consolidate.md:57,95）。

2. CN/EN 同构：通过。
- 依据：`G3(a)` 对报告对给出对称通过，六形状同构为真（docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-index.txt:10；docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-tree.txt:10）。

3. `artifacts` 围栏与暂存集双向一致：通过。
- 依据：围栏与 `impl.numstat` 双向差集均为空（独立脚本比对，口径同 r1 要求）。

4. 越界改动检查：通过。
- 依据：暂存集仅含本片允许载体（tmp/P80-SG6-GATELEVEL-CONSOLIDATE/impl.numstat:1-18），未出现 p57 门表、`s1-exit` 门表、v62 快照、准则、CONTRACTS、schema、fixtures、产品代码、CI、其它 ADR 的改动。

## 发现清单

| # | SEVERITY | file:line | 问题 | 最小修复 |
|---|---|---|---|---|
| 1 | Medium | docs/validation/sg6-gatelevel-consolidate.md:155; docs/validation/sg6-gatelevel-consolidate.md:158; docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/06-gates/gate-ci.txt:1,10 | §12 缺少 `--scope=ci` 的本地退化口径显式说明；当前文本不足以防止把 `changed=2` 的 ci 读数误解为本片同批快照读数。 | 在 §12 增补一条作用域声明：`--scope=ci` 为本地退化口径（对象是 `CHECKLIST_v62` 变更，`changed=2`），仅作环境与链路健康观测；本片同批一致性以 `--scope=index/tree`（`changed=18`）与 `impl.numstat` 对齐为准。 |

## 本轮额度与结论

- 本片档位：S 档；⑦ 额度 = 1 终审 + 1 复审（盲审另计）。
- 当前轮次：复审 r2（唯一一次复审）。
- 严重度计数：High 0 / Medium 1 / Low 0 / Info 0。
- 是否触发「停下升级」：是。依据：复审后仍存在 Medium，按本轮硬约束必须停下升级用户裁决，不得自行判定通过，不得开 r3。

RE-REVIEW: FINDINGS