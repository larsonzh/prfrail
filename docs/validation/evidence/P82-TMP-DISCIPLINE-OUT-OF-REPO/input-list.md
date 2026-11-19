# 审查输入清单 — P82 `TMP-DISCIPLINE-OUT-OF-REPO`

> 依据：§5.1 时序说明（OB-93）：审查输入包须逐件断言申报输入存在，包内附清单（文件名 + 行数或字节数）。

## 步骤 0 建包时的申报输入（① 阶段）

| 序 | 文件 | 行数 | 字节数 |
|---|---|---|---|
| 1 | `docs/DELIVERY_DIRECTIVE.md` | 2452 | 500519 |
| 2 | `docs/DELIVERY_DIRECTIVE_EN.md` | 2452 | 610510 |
| 3 | `docs/t027/REMAINING_SLICES.md` | 2420 | 377499 |
| 4 | `docs/t027/REMAINING_SLICES_EN.md` | 2407 | 423660 |
| 5 | `docs/t027/CHECKLIST_v64.md` | 224 | 28427 |
| 6 | `docs/t027/CHECKLIST_v64_EN.md` | 224 | 33158 |

## 逐轮递交记录

| 轮次 | 审查方 | 递交件 | 生成器断言 | 备注 |
|---|---|---|---|---|
| ① 架构 | `DeepSeek V4 Pro (deepseek)` | 上表六件 + 起跑消息 | 逐件存在，行数与字节数与上表一致 | 只读 |
| ⑤ r1 / r2 | `DeepSeek V4 Pro (deepseek)` | 起跑消息、方案书、变更集 diff 与 numstat、③ 报告、门禁原文 | 逐件存在 | 只读；附 ⑤ r1 清单于 r2 |
| ⑤ r3（用户授权收窄读） | `DeepSeek V4 Pro (deepseek)` | 仅：层 3 五处单元格 before/after、§一 第 ⑥ 项 before/after、三作用域门禁原始行 | 逐件存在 | 输入收窄声明见报告 A-3 |
| ⑥ r1 / r2 | `MAI Code 2 Flash (copilot)` | 变更集 diff 与 numstat、切片定义、准则与快照与台账、门禁原文 | 逐件存在 | 不附任何他人清单 |
| ⑦ r1 | `GPT-5.3-Codex (copilot)` | 切片定义、变更集 diff 与 numstat、准则与快照与台账、门禁原文 | 逐件存在 | **盲态首轮，不附 ⑤/⑥ 清单** |
| ⑦ r2 | `GPT-5.3-Codex (copilot)` | 上述 + ⑦ r1 原文 + 报告对 + ⑥ r2 原件 | 逐件存在 | 复审轮，允许附清单 |
