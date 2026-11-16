# MANIFEST — P76-S1-EXIT-CLOSEOUT-ASSESSMENT 证据包

切片：`S1-EXIT-CLOSEOUT-ASSESSMENT`（P76，S1 Exit 收口路线图评估，评估类 `[SLICE]`，M 档）。建包时点 2026-11-16（`HEAD` = `b90de21`，准则 v1.69）。

**包状态**：步骤 0 建包（三件载荷：`00-briefing.md` / `input-list.md` / `freedom-list.md`）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` 载荷；⑧c 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等。两者不在下表内（避免与载荷集混计）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10；2026-11-16 P76 起跑，仅授权步骤 0 + ① 阶段 + 停点） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘初稿；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `01-arch/design-r1.md` | ① 阶段评估方案书原文（`deep-reasoner` + `DeepSeek V4 Pro`；含 §17 主控抽验附注） |
| `02-batch/batch-01-s1-exit-p57.md` | 批量辅助 1 事实枚举（`GLM-5.3-Flash`）：s1-exit / p57 §4.2 八门现状 |
| `02-batch/batch-02-s1-reach.md` | 批量辅助 2 事实枚举（`GLM-5.3-Flash`）：s1-reach 十门 |
| `02-batch/batch-03-lagging-docs.md` | 批量辅助 3 事实枚举（`GLM-5.3-Flash`）：滞后口径 + 独立安全评审检索 |
| `02-ruling/ruling-2026-11-16.md` | ②a 用户裁定登记（八项裁定原文逐项照录，CN 单语） |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 原文（`PRE-REVIEW: FINDINGS`，0 危 1 高 0 中 5 低） |
| `04-reviews/prereview-r2.md` | ⑤ 预审 r2 原文（`PRE-REVIEW: PASS WITH FIXES`，0 中及以上，2 低） |
| `04-reviews/final-r1.md` | ⑦ 终审 r1 原文（`FINAL REVIEW: PASS WITH FIXES`，0 中及以上，1 低；含 `SKIP CONFIRMED`） |
| `06-gates/gate-index.txt` | ⑨ index 作用域门禁原始输出 |
| `06-gates/gate-tree.txt` | ⑨ tree 作用域门禁原始输出 |
| `06-gates/gate-ci.txt` | ⑨ ci 作用域门禁原始输出 |
| `06-gates/gate-selftest.txt` | ⑨ SELFTEST 原始输出 |

## 后续追加计划

- 无新增载荷；⑧c 只重算 `SHA256SUMS.txt` 并核验五条读数。

## 编码

全部 `.md` 为 UTF-8 with BOM + LF（仓库约定）；`SHA256SUMS.txt` 为无 BOM + LF。
