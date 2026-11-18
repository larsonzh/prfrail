# MANIFEST — P81-DIRECTIVE-MODEL-UPGRADE 证据包

切片：`DIRECTIVE-MODEL-UPGRADE`（P81，⑥ 与 ⑦ 深度档载体模型直接替换 + P80 裁定 ① 留痕补落，治理片 `[SLICE]`，M 档）。建包时点 2026-11-18（`HEAD` = `ac15263`，准则 v1.69，本片升至 v1.70）。

**包状态**：步骤 0 建包（三件初始载荷）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` 载荷；⑧c 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等。两者不在下表内（避免与载荷集混计）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户 2026-11-18 起跑消息逐字原文（§13.1.10；授权可回溯性依据） |
| `01-arch/design-r1.md` | ① 阶段方案书原文（Q1 至 Q7） |
| `02-batch/batch-01-inventory.md` | 批量辅助落点清单（现役落点 / 历史行 / 冻结件三层归属） |
| `04-reviews/blind-note.md` | ⑦ 盲审不抽的判定记录（§7.7） |
| `04-reviews/final-r1.md` | ⑦ 终审 r1 输出原文 |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 输出原文 |
| `04-reviews/scan-r1.md` | ⑥ 独立扫描输出原文 |
| `04-reviews/scan-r2.md` | 本片载荷（见验证报告相应章节） |
| `06-gates/gate-ci.txt` | `--all --scope=ci` 原始输出 |
| `06-gates/gate-index.txt` | `--all --scope=index` 原始输出 |
| `06-gates/gate-selftest.txt` | `--selftest` 原始输出 |
| `06-gates/gate-tree.txt` | `--all --scope=tree` 原始输出 |
| `freedom-list.md` | 执行自由度清单（§1.7.3；步骤 0c 初稿、⑥ 之后定稿） |
| `input-list.md` | 审查输入清单与逐轮递交记录（OB-93 的逐件断言输入） |

## 后续追加计划

- 按 ⑤ / ⑥ / ⑦ / ⑨ 各轮追加；⑧c 只重算 `SHA256SUMS.txt` 并核验五条读数。
- 新增载荷必须同轮补本表条目并重算 `SHA256SUMS.txt`（§13.1.12）。

## 编码

全部 `.md` 为 UTF-8 with BOM + LF（仓库约定）；`SHA256SUMS.txt` 与 `06-gates/*.txt` 为无 BOM + LF。

## 对应引用点

- 准则 `docs/DELIVERY_DIRECTIVE.md` 的 §12.4 表尾引用块（P80 裁定 ① 的补落）；验证报告 `docs/validation/directive-model-upgrade{,_EN}.md` 的「证据包」节。
