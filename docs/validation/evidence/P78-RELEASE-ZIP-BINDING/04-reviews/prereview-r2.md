# ⑤ 预审 r2（DeepSeek V4 Pro，只读；整改复审）

切片：P78 `RELEASE-ZIP-BINDING`。输入包：整改 diff `fix-r2.diff` + 全量变更集 + 报告对 + 方案书。

```
PRE-REVIEW: PASS WITH FIXES
```

（0 项 Medium 及以上；额度 1 + 1 用尽。）

## 整改闭合核验

| # | 原发现 | 结论 | 依据 |
|---|---|---|---|
| 1 | Medium §9 预声明读数 | 闭合 | `docs/validation/release-zip-binding.md:83` 表前新增说明段（⑨ 实跑口径 + 两处时序红因） |
| 2 | Medium `s1-exit` 陈旧句未登记 | 闭合 | 报告 `:102` §10 新增条（历史结论 vs 门表，指向 §7 与剩余风险第 1 条） |
| 3 | Low artifacts 未产出 | 闭合 | 报告 `:103` §10 新增条（时序声明） |
| 4 | Low 五处增补偏离未披露 | 闭合 | 报告 `:44` §4「口径偏离披露」段 |
| 5 | Low「无正式发行包」陈旧 | 闭合 | `INSTALLATION_PLAN:7` 与 `S1_RELEASE_NOTES:18` 改为候选发行口径，CN/EN 同步 |

特别核项（不得暗示 SLA 生效）：通过。`INSTALLATION_PLAN:7` 保留「（非已发布 SLA，S1 exit 整体批准仍待 S2）」；`S1_RELEASE_NOTES:5` 保留「不构成 S1 exit 整体通过」；`:18` 保留「（非已发布 SLA）」；`OPERATIONS:5` 保留「候选状态、非已发布 SLA」；`S1_SUPPORT_MATRIX:5` 与 `ADR_REGISTER:62` 原样未动。

## 新增发现

| # | 级别 | 问题 | 依据 |
|---|---|---|---|
| 1 | Low | 报告 §5 引 `INSTALLATION:206`，主实现使文件后移两行，实际绑定句位于 `:208` | `release-zip-binding.md:55` vs 磁盘 `INSTALLATION.md:206/208` |
| 2 | Low | `INSTALLATION_PLAN.md:27` 决策表「更新方式」状态列仍写「最终 ZIP 待复验」，与同文件 `:38` 矛盾；该格不含三个检索关键词 | `INSTALLATION_PLAN.md:27` vs `:38`（EN 同位置） |

两条均为 Low、非整改引入，不构成 Medium+；已在 ⑤ 之后由主控顺手闭合（报告 §5 改引 `:208` 并注明漂移；`INSTALLATION_PLAN:27` 改为「最终 ZIP 已复验（P78，2026-11-17）」，CN/EN 同步）。

## 我无法验证的部分

1. §9 门禁读数的实跑性（⑤ r2 时点 ⑨ 尚未发生）。
2. 外部事实（Release、CI run、四项复验的实际执行）。
3. 起跑消息同轮授权原文未直接见于输入包，仅由方案书与 ADR-046 承载。
