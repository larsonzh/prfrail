# MANIFEST — P77-ADR010-APPROVAL 证据包

切片：`ADR010-APPROVAL`（P77，ADR-010 批准后的 EOL 通知方式文档同步，文档类 `[SLICE]`，S 档）。建包时点 2026-11-17（`HEAD` = `fb90a20`，准则 v1.69）。

**包状态**：步骤 0 建包（三件载荷：`00-briefing.md` / `input-list.md` / `freedom-list.md`）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` 载荷；⑧c 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等。两者不在下表内（避免与载荷集混计）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10；2026-11-16 P77 起跑，S 档文档类，一口气跑完 ⑩） |
| `01-arch/design-r1.md` | ① 阶段方案书原文（`deep-reasoner` + `DeepSeek V4 Pro`；Q1–Q5） |
| `02-batch/batch-01-adr010-register-format.md` | 批量辅助 1 事实枚举（`GLM-5.3-Flash`）：ADR-010 原文 + ADR 主表 4 列 + 批准记录 / 四要素范例 |
| `02-batch/batch-02-eol-survey.md` | 批量辅助 2 事实枚举（`GLM-5.3-Flash`）：六份候选文档的 EOL 表述逐条枚举 |
| `04-reviews/final-r1.md` | ⑦ 终审 r1 原文（`GPT-5.3-Codex`；含 `SKIP CONFIRMED`） |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 原文（`DeepSeek V4 Pro`） |
| `04-reviews/prereview-r2.md` | ⑤ 预审 r2 原文（`DeepSeek V4 Pro`，如触发） |
| `06-gates/gate-ci.txt` | ⑨ ci 作用域门禁原始输出 |
| `06-gates/gate-index.txt` | ⑨ index 作用域门禁原始输出 |
| `06-gates/gate-selftest.txt` | ⑨ SELFTEST 原始输出 |
| `06-gates/gate-tree.txt` | ⑨ tree 作用域门禁原始输出 |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘初稿；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |

## 后续追加计划

- 按 ① / 批量辅助 / ⑤ / ⑦ / ⑨ 各轮追加；⑧c 只重算 `SHA256SUMS.txt` 并核验五条读数。

## 编码

全部 `.md` 为 UTF-8 with BOM + LF（仓库约定）；`SHA256SUMS.txt` 为无 BOM + LF。
