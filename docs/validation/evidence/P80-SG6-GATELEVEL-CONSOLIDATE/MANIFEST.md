# MANIFEST — P80-SG6-GATELEVEL-CONSOLIDATE 证据包

切片：`SG6-GATELEVEL-CONSOLIDATE`（P80，S-G6 门级侧合并判定 + 并陈口径 + S1 Exit 收口下一步清单，文档类 `[SLICE]`，S 档）。建包时点 2026-11-18（`HEAD` = `b5ff8ee`，准则 v1.69）。

**包状态**：步骤 0 建包（三件初始载荷）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` 载荷；⑧c 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等。两者不在下表内（避免与载荷集混计）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10；2026-11-18 P80 起跑，S 档文档类） |
| `01-arch/design-r1.md` | ① 阶段方案书原文（`deep-reasoner` + `DeepSeek V4 Pro`；Q1 至 Q6） |
| `02-batch/batch-01-sg6-sources.md` | 批量辅助事实枚举（`GLM-5.3-Flash`）：八门判定与证据落点现状 |
| `04-reviews/blind-r1.md` | ⑦ 盲审输出（§7.7 按比例抽查，额度另计） |
| `04-reviews/final-r1.md` | ⑦ 终审 r1 输出 |
| `04-reviews/final-r2.md` | ⑦ 复审 r2 输出 |
| `04-reviews/narrow-r1.md` | ⑦ 收窄读 r1 输出（§7.8.3 用户授权例外，仅本片有效） |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 输出 |
| `06-gates/gate-ci.txt` | `--scope=ci` 原始输出 |
| `06-gates/gate-index.txt` | `--scope=index` 原始输出 |
| `06-gates/gate-selftest.txt` | `--selftest` 原始输出 |
| `06-gates/gate-tree.txt` | `--scope=tree` 原始输出 |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘初稿；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |

## 后续追加计划

- 按 ⑤ / ⑦ / ⑨ 各轮追加；⑧c 只重算 `SHA256SUMS.txt` 并核验五条读数。

## 编码

全部 `.md` 为 UTF-8 with BOM + LF（仓库约定）；`SHA256SUMS.txt` 与 `06-gates/*.txt` 为无 BOM + LF。
