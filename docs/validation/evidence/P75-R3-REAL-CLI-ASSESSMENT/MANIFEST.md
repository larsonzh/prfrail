# MANIFEST — P75-R3-REAL-CLI-ASSESSMENT 证据包

切片：`R3-REAL-CLI-ASSESSMENT`（P75，层 1 T027-B 主线里程碑余项评估，评估类 `[SLICE]`，M 档）。建包时点 2026-11-15（`HEAD` = `c5c3a78`，准则 v1.69）。

**包状态**：步骤 0 建包（三件载荷：`00-briefing.md` / `input-list.md` / `freedom-list.md`）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` 载荷；⑧ 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等。两者不在下表内（避免与载荷集混计）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10；含 2026-11-15 同轮显式授权） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘初稿；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `01-arch/design-r1.md` | ① 阶段评估原文（`deep-reasoner` + `DeepSeek V4 Pro`；`DESIGN: DONE`；含主控行号核验修正节） |
| `02-batch/batch-01-t027b-at23-b1b2-facts.md` | 批量辅助 ①（`GLM-5.3-Flash`）T027-B / AT-23 / B1-B2 事实枚举原文（实得 5952 字符） |
| `02-batch/batch-02-r08-r09-sandbox-facts.md` | 批量辅助 ②（`GLM-5.3-Flash`）R-08 / R-09 与 Sandbox 边界事实枚举原文（实得 5808 字符） |
| `02-batch/ruling-2026-11-15-r3.md` | 用户 2026-11-15 的 R-3 裁定原文归档（节录：裁定主体；⑤ r1 发现 #7 整改） |
| `04-reviews/blind-r1.md` | ⑦ 盲审原文（`BLIND REVIEW: PASS`，0 发现；§7.7 按比例抽查） |
| `04-reviews/final-r1.md` | ⑦ 终审原文（`FINAL REVIEW: PASS`，0 发现；`⑥ SKIP: SKIP CONFIRMED`） |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 原文（`PASS WITH FIXES`，7 Medium + 2 Low） |
| `04-reviews/prereview-r2.md` | ⑤ 预审 r2 原文（`PASS WITH FIXES`，3 Low，均已整改） |
| `06-gates/gate-ci.txt` | ④ 门禁 ci 作用域原始输出（点时刻捕获；包自洽行反映声明前状态） |
| `06-gates/gate-index.txt` | ④ 门禁 index 作用域原始输出（同上） |
| `06-gates/gate-selftest.txt` | ④ SELFTEST 原始输出（同上） |
| `06-gates/gate-tree.txt` | ④ 门禁 tree 作用域原始输出（同上） |

## 后续追加计划

- `01-arch/design-r1.md`：① 阶段评估原文（`deep-reasoner` + `DeepSeek V4 Pro`）。
- `02-batch/batch-01-*.md` / `batch-02-*.md`：批量辅助事实枚举原文（`deep-reasoner` + `GLM-5.3-Flash`，≥2 次）。
- `04-reviews/prereview-r1.md` / `scan-r1.md`（若 ⑥ 未跳过）/ `blind-r1.md` / `final-r1.md`：⑤ ⑥ ⑦ 原文。
- `06-gates/*.txt`：三作用域与 SELFTEST 原始输出。

## 编码

全部 `.md` 为 UTF-8 with BOM + LF（仓库约定）；`SHA256SUMS.txt` 为无 BOM + LF。
