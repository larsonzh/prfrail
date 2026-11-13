# MANIFEST — P71-S1-SCOPE-ANNOTATE 证据包

切片：`S1-SCOPE-ANNOTATE`（P71，层 2 治理侧，文档类 `[SLICE]`，M 档；S1 范围阶段归属列回写 + visible 调试指引落地）。建包时点 2026-11-13（`HEAD` = `ca88fa3`，准则 v1.67）。

**包状态**：步骤 0 建包（三件载荷：`00-briefing.md` / `input-list.md` / `freedom-list.md`）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/`（**必含盲审原文**，§13.1.13）/ `06-gates/` 载荷；⑧ 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10 起跑消息即工件；含 2026-11-13 同轮显式授权） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `01-arch/design-r2.md` | ① 架构方案书 design-r2 原文（含 F-1 裁定结果与逐项回写口径；`DESIGN: DONE`） |
| `02-batch/batch-01-req-enumeration.md` | 批量辅助 ① 转录（REQ 现状与 PC 引用枚举；主控脚本已复算一致） |
| `02-batch/batch-02-at-gates-enumeration.md` | 批量辅助 ② 转录（AT 与八门现状枚举） |
| `04-reviews/blind-r1.md` | ⑦ 盲审 r1 原文（不附 ⑤ / ⑥ 清单） |
| `04-reviews/final-r1.md` | ⑦ 终审 r1 原文（含 ⑥ 跳过判定确认） |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 原文（5 Low + 若干 Info） |
