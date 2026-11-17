# P77 输入清单（§13.1.1）

切片：P77 `ADR010-APPROVAL`（`[SLICE]` · 文档类 S 档）
基线：HEAD = origin/main = `fb90a20`；准则 v1.69；SELFTEST 737/737；三作用域 `TOTAL_FAIL=0`

## 一、核心判定对象

| # | 文件 | 用途 |
|---|---|---|
| 1 | `docs/ADR_REGISTER.md` / `ADR_REGISTER_EN.md` | ADR-010 当前状态栏（`PROPOSED / NOT APPROVED`）+ ADR-014–ADR-045 的四要素格式基准 + ADR-006–ADR-009 的批准记录范例（Q1） |
| 2 | `docs/S1_SUPPORT_MATRIX.md` / `S1_SUPPORT_MATRIX_EN.md` | 「版本与安全维护候选政策」段：待审批提议句 + 「该提议尚未获批」句（Q2 第 1 项） |

## 二、第二层候选（视实读结果判定）

| # | 文件 | 用途 |
|---|---|---|
| 3 | `docs/INSTALLATION_PLAN.md` / `INSTALLATION_PLAN_EN.md` | `:32` 附近 EOL 项实读（Q2 第 2 项） |
| 4 | `docs/PRODUCT_REQUIREMENTS.md` / `PRODUCT_REQUIREMENTS_EN.md` | §7「EOL 通知」待落实句（Q2 第 3 项） |
| 5 | `docs/S1_RELEASE_NOTES.md` / `S1_RELEASE_NOTES_EN.md` | 「发布前剩余门禁」段是否含 EOL 内容（Q2 第 4 项） |
| 6 | `docs/INSTALLATION.md` / `INSTALLATION_EN.md` | 是否含 EOL 表述（Q2 第 5 项） |
| 7 | `docs/OPERATIONS.md` / `OPERATIONS_EN.md` | 是否含 EOL 表述（Q2 第 6 项） |

## 三、门级与状态依据

| # | 文件 | 用途 |
|---|---|---|
| 8 | `docs/validation/p57-closeout-assessment.md` / `_EN.md` | §4.2 的 S-G4 判据原文（Q3 基准；**本片零改动**） |
| 9 | `docs/validation/s1-exit.md` / `s1-exit_EN.md` | S1 Exit 门表现状（Q3 对照；**本片零改动**） |
| 10 | `docs/validation/s1-exit-closeout-assessment.md` / `_EN.md` | P76 报告对：S-G4 解锁条件、ADR-046 的 Q10=(C) 分两步（Q3 / ③ 一致性） |
| 11 | `docs/t027/CHECKLIST_v60.md` / `_EN.md` | 当前快照（§四 层 1 队列项 P77 六要素；**本片零改动**） |

## 四、执行纪律与台账

| # | 文件 | 用途 |
|---|---|---|
| 12 | `docs/DELIVERY_DIRECTIVE.md` / `_EN.md`（v1.69） | §1.7 自由度 / §5.1 流水线 / §5.2 批次 / §7.5 额度 / §9.2 高风险 / §11.2 回填 / §13 输入清单 |
| 13 | `docs/t027/REMAINING_SLICES.md` / `_EN.md` | 切片台账：P76 段格式基准 + P77 段登记（步骤 0c）与完成记录（⑧b） |
| 14 | `docs/CODING_CONVENTIONS.md` / `_EN.md` | 编码 / 行尾 / Git 纪律（本片落盘执行） |

## 五、辅助材料

| # | 文件 | 用途 |
|---|---|---|
| 15 | `docs/validation/evidence/P77-ADR010-APPROVAL/00-briefing.md` | 起跑消息原文 |
| 16 | `docs/validation/evidence/P77-ADR010-APPROVAL/freedom-list.md` | 执行自由度清单（⑥ 后 / ⑦ 前定稿） |

## 六、取证方式说明

- 行号引用以**实读为准**；行号漂移时以「文件名 + 节号 / 段落标题」为主锚点。
- 组合 grep 模式集（Q4）：`ADR-010` / `待审批.*ADR-010` 或 `ADR-010.*待审批` / `尚未获批` 或 `未获批` / `EOL` 与 `ADR` 同行同现 / `EOL 通知` 或 `EOL 公告`；**禁止**单独 grep `EOL` 或 `90`。
- 全部输入均为**只读**；① 阶段不改任何文档。
