# P76 输入清单（§13.1.1）

切片：P76 `S1-EXIT-CLOSEOUT-ASSESSMENT`（`[SLICE]` · 评估类 M 档）
基线：HEAD = origin/main = `b90de21`；准则 v1.69；SELFTEST 737/737；三作用域 `TOTAL_FAIL=0`

## 一、核心判定对象

| # | 文件 | 用途 |
|---|---|---|
| 1 | `docs/validation/s1-exit.md` / `s1-exit_EN.md` | S1 Exit 八门现状主体；**`:52` 核心任务链行内口径待修正**（Q3） |
| 2 | `docs/validation/p57-closeout-assessment.md` / `_EN.md` | §4.2 八门原始清单 + §4.3 原编号（Q1 / Q3 / Q4 基准） |
| 3 | `docs/validation/s1-reach-reachability.md` / `_EN.md` | 十门框架（G-A–G-I + OB-21）（Q1 映射对照；**本片零改动**） |

## 二、门级与状态依据

| # | 文件 | 用途 |
|---|---|---|
| 4 | `docs/validation/r3-real-cli-assessment.md` / `_EN.md` | P75 报告；§11.2 唯一残余项；R-3 = (A) 裁定 |
| 5 | `docs/validation/at23-restricted-candidate.md` / `_EN.md` | AT-23 阶段 1 条件闭合 = 永久事实（Q2） |
| 6 | `docs/validation/at24-supervised-black-box.md` / `_EN.md` | AT-24 条件闭合（Q2） |
| 7 | `docs/validation/r08-r09-ruling.md` / `_EN.md` | R-08 / R-09 裁定（Q2 / Q11 边界） |
| 8 | `docs/validation/acceptance-security-counterexamples.md` / `_EN.md` | 验收安全反例（Q7 独立安全评审历史状态核实对象之一） |
| 9 | `docs/SECURITY.md` / `SECURITY_EN.md` | 威胁模型 / 硬门禁 / 独立安全评审相关表述（Q7） |
| 10 | `docs/PRODUCT_REQUIREMENTS.md` / `_EN.md` | §7「独立安全评审结论」（**进入 S1 前**的阻断门禁）+ REQ 阶段归属（Q7 / Q11） |

## 三、门级解锁路径依据

| # | 文件 | 用途 |
|---|---|---|
| 11 | `docs/S1_RELEASE_NOTES.md` / `_EN.md` | 发布说明；S-G3 解锁前提（Q5） |
| 12 | `docs/INSTALLATION.md` / `INSTALLATION_EN.md` | 「当前尚无正式 GitHub Release」原文（Q5） |
| 13 | `docs/INSTALLATION_PLAN.md` / `INSTALLATION_PLAN_EN.md` | 安装生命周期复验路径（Q5） |
| 14 | `docs/ADR_REGISTER.md` / `ADR_REGISTER_EN.md` | ADR-009 / ADR-010 / ADR-044 / ADR-045（Q6 / Q10） |
| 15 | `docs/DEV_PLAN.md` / `DEV_PLAN_EN.md` | T018 任务级依赖字段与状态（Q8）；T027-B 段 |
| 16 | `docs/S1_SUPPORT_MATRIX.md` / `S1_SUPPORT_MATRIX_EN.md` | 「CLI Agent / AgentRunner」行口径滞后（Q6 / Q9） |
| 17 | `docs/TEST_STRATEGY.md` / `TEST_STRATEGY_EN.md` | AT 阶段归属（Q11） |

## 四、执行纪律与快照

| # | 文件 | 用途 |
|---|---|---|
| 18 | `docs/DELIVERY_DIRECTIVE.md` / `_EN.md`（v1.69） | §1.7 自由度 / §5.1 流水线 / §7.5 额度 / §7.7 盲审 / §9.2 高风险 / §13 输入清单 |
| 19 | `docs/t027/CHECKLIST_v59.md` / `_EN.md` | 当前快照（含 P75 补正；§二 ⑤ 编号体系已修正） |
| 20 | `docs/t027/CHECKLIST_v57.md` / `v58` | 滚动窗口基线（160 / 148） |
| 21 | `docs/t027/REMAINING_SLICES.md` / `_EN.md` | 切片台账（P76 完成记录在 ⑧b 补登） |

## 五、辅助材料

| # | 文件 | 用途 |
|---|---|---|
| 22 | `docs/CODING_CONVENTIONS.md` / `_EN.md` | 编码 / 行尾 / Git 纪律（本片落盘执行） |
| 23 | `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/00-briefing.md` | 起跑消息原文 |
| 24 | `docs/validation/evidence/P76-S1-EXIT-CLOSEOUT-ASSESSMENT/freedom-list.md` | 执行自由度清单（⑥ 后 / ⑦ 前定稿） |

## 六、取证方式说明

- 行号引用以**实读为准**；若行号漂移，以「文件名 + 节号」为主锚点。
- `docs/validation/` 下独立安全评审检索关键词：`security review` / `安全评审` / `independent review`（批量 3 承担）。
- 全部输入均为**只读**；本片 ① 阶段不改任何文档。
