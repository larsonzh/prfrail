# MANIFEST — P72-REQ020-S2-RESTRICTION 证据包

切片：`REQ020-S2-RESTRICTION`（P72，层 2 治理侧，文档类 `[SLICE]`，S 档；REQ-020 的 S2 承诺收窄登记 + 关联文档同步）。建包时点 2026-11-14（`HEAD` = `70eccee`，准则 v1.68）。

**包状态**：步骤 0 建包（三件载荷：`00-briefing.md` / `input-list.md` / `freedom-list.md`）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` 载荷；⑧ 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10 起跑消息即工件；含 2026-11-14 同轮显式授权） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `01-arch/design-r1.md` | ① 现状核对与差异清单 design-r1 原文（含 (h) 窄读裁决与八条停升级判定；`DESIGN: DONE`） |
| `02-batch/batch-01-req-enumeration.md` | 批量辅助 ① 转录（REQ-020 / REQ-011 / REQ-022 / PC-01 / PC-05 / PC-06 枚举，GLM-5.3-Flash） |
| `02-batch/batch-02-at-matrix-enumeration.md` | 批量辅助 ② 转录（AT-01 / AT-11 / AT-20 与 S1_SUPPORT_MATRIX 枚举，GLM-5.3-Flash） |
| `04-reviews/final-r1.md` | ⑦ 终审 r1 原文（不附 ⑤ 清单；含 ⑥ SKIP CONFIRMED） |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 原文（2 Low + 4 Info；`PRE-REVIEW: PASS WITH FIXES`） |
| `06-gates/gates-r1.txt` | 三作用域门禁捕获（index / tree / ci，`TOTAL_FAIL=0`） |
| `06-gates/selftest-r1.txt` | `gate.js --selftest` 捕获（`SELFTEST: PASS 737/737`） |
