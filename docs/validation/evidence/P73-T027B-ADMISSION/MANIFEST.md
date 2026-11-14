# MANIFEST — P73-T027B-ADMISSION 证据包

切片：`T027B-ADMISSION`（P73，层 2 治理侧，M 档评估类 `[SLICE]`；OB-55 归属三选一评估 + B 段实施边界）。建包时点 2026-11-14（`HEAD` = `0cac174`，准则 v1.68）。

**包状态**：步骤 0 建包（三件载荷：`00-briefing.md` / `input-list.md` / `freedom-list.md`）；后续按轮追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` 载荷；⑧ 收口时归档定稿件并重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（§13.1.10 起跑消息即工件；含 2026-11-14 同轮显式授权） |
| `input-list.md` | 本片输入清单（步骤 0 落盘，§13.1.1） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘；⑥ 之后、⑦ 之前定稿，§1.7.3） |
| `01-arch/design-r1.md` | ① 现状核对与 OB-55 三选一评估 design-r1 原文（`DESIGN: DONE`；八件评估 + 推荐 + 需评估 6 项 + 停升级逐条裁决） |
| `02-batch/batch-01-p36-p69-facts.md` | 批量辅助 ① 转录（C1–C7 / N1–N4 / P39–P44 / P69 / OB-55 事实枚举，GLM-5.3-Flash，8383 字符） |
| `04-reviews/prereview-r1.md` | ⑤ 预审 r1 原文（1 High + 1 Medium + 1 Low + 2 Info；`PRE-REVIEW: PASS WITH FIXES`） |
| `04-reviews/final-r1.md` | ⑦ 终审 r1 原文（不附 ⑤ 清单；`FINAL REVIEW: PASS`，含 ⑥ SKIP CONFIRMED 与盲审不应抽） |
| `06-gates/gates-r1.txt` | 三作用域门禁捕获（index / tree / ci，`TOTAL_FAIL=0`） |
| `06-gates/selftest-r1.txt` | `gate.js --selftest` 捕获（`SELFTEST: PASS 737/737`） |
| `02-ruling/ob55-ruling.md` | OB-55 裁定登记（②a；四要素 + 七项裁定摘要） |
| `02-batch/batch-02-rulings-gates.md` | 批量辅助 ② 转录（R-08 / R-09 / p57 八门 / s1-exit / R-06 / 支持矩阵 / B2 §8，GLM-5.3-Flash，9073 字符） |
