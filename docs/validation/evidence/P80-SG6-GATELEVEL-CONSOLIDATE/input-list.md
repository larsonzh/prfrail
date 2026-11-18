# P80 输入清单（§13.1.1）

切片：P80 `SG6-GATELEVEL-CONSOLIDATE`（`[SLICE]` · S 档文档类）
基线：HEAD = origin/main = `b5ff8ee`；准则 v1.69；SELFTEST 737/737；三作用域 `TOTAL_FAIL=0`

## 一、核心判定对象

| # | 文件 | 用途 |
|---|---|---|
| 1 | `docs/validation/p57-closeout-assessment.md` / `_EN.md` | §4.2 八门表判据原文与当前判定（Q1 / Q2 判据字面；**本片零改动**） |
| 2 | `docs/validation/s1-exit.md` / `_EN.md` | 门表行（只读摘录）、`:52` 核心任务链、`:66` 发行绑定句（Q1 / Q3；**本片零改动**） |
| 3 | `docs/validation/s1-exit-closeout-assessment.md` / `_EN.md` | P76 报告：Q2 门级判定口径、Q5 四步、Q8=(c) 并陈口径、Q10=(C) 分两步、§12 收口路线图（Q1 / Q2 / Q3 / Q4） |
| 4 | `docs/validation/adr010-approval.md` / `_EN.md` | P77 报告：S-G4 收口章与 `ADR_REGISTER:20` 落点（Q1） |
| 5 | `docs/validation/release-zip-binding.md` / `_EN.md` | P78 报告：S-G3 收口章（§5 四元组、§7 收口）、发行绑定三处落点（Q1） |
| 6 | `docs/DEV_PLAN.md` / `_EN.md` | `:94` T018 依赖字段、`:96` T018 状态、`:98` P78 注记（Q2 依赖字段 / Q5 归属） |
| 7 | `docs/ADR_REGISTER.md` / `_EN.md` | `:20` ADR-009 与 ADR-010 状态列、ADR-046 / ADR-047（Q1 / Q4） |
| 8 | `docs/t027/CHECKLIST_v62.md` / `_EN.md` | §四 八门状态表与 P80 队列项（**冻结，本片零改动**；Q1 / Q2 口径来源） |

## 二、下一步清单依据（Q4）

| # | 文件 | 用途 |
|---|---|---|
| 9 | `docs/S1_SUPPORT_MATRIX.md` / `_EN.md` | 已绑定发行记录节；R-3 = (A) 与 S-G8 blocked 口径 |
| 10 | `docs/SECURITY.md` / `_EN.md` | 独立安全评审历史状态（从未执行）与豁免顺延口径（Q4：P79 前置） |
| 11 | `docs/PRODUCT_REQUIREMENTS.md` / `_EN.md` | §7 的 S0→S1 门禁（独立安全评审）与阶段归属（Q4） |
| 12 | `docs/validation/s1-reach-reachability.md` / `_EN.md` | 十门可达性刷新与「门级与任务级不同源」结论 |

## 三、门级与状态依据

| # | 文件 | 用途 |
|---|---|---|
| 13 | `docs/DELIVERY_DIRECTIVE.md` / `_EN.md`（v1.69） | §5.2 批次矩阵、§7.2.1 ⑥ 跳过三条件、§7.5 额度、§7.7 按比例抽查、§9.2 高风险、§11.2 回填、§13 输入清单（**本片零改动**） |
| 14 | `docs/ADR_REGISTER.md` 的 ADR-046 / ADR-047 | 分两步批准与安全评审豁免（Q2 / Q4） |

## 四、执行纪律与台账

| # | 文件 | 用途 |
|---|---|---|
| 15 | `docs/t027/REMAINING_SLICES.md` / `_EN.md` | P78 段格式基准 + P80 段登记（步骤 0c）与完成记录（⑧b） |
| 16 | `docs/CODING_CONVENTIONS.md` / `_EN.md` | 编码 / 行尾 / Git 纪律（本片落盘执行） |
| 17 | `docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/00-briefing.md` | 起跑消息原文 |
| 18 | `docs/validation/evidence/P80-SG6-GATELEVEL-CONSOLIDATE/freedom-list.md` | 执行自由度清单（步骤 0 初稿；⑥ 之后、⑦ 之前定稿） |

## 五、取证方式说明

- 行号引用以**实读为准**；行号漂移时以「文件名 + 节号 / 段落标题」为主锚点。
- Q6 组合 grep 模式集（独立执行、禁止单关键词联查）：`S-G6` / `SG6` / `门级侧`。
- 批量辅助 `batch-01` 只做事实枚举，不做判定；产出交 ①。
- 全部输入均为**只读**；① 阶段不改任何文档。
