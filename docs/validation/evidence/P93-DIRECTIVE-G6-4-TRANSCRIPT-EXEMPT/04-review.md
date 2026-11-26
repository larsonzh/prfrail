# P93 审查回执摘要（⑤ / ⑥ / ⑦）

## ⑤ 预审（DeepSeek V4 Pro，1+2 额度内用 1）

- 首行：`PRE-REVIEW: PASS WITH FIXES`；`verdict: no Medium+`。
- 发现：1 条 Low —— `fence-exempt N` 仅在 findings=0 时进入 detail，与准则「豁免行数在 detail 中如实报出」不完全一致。
- 处置：把准则句精确化为「无 finding 时，豁免行数在 detail 中以 `fence-exempt N` 如实报出」（CN/EN 同改），不改共享 `emit`。

## ⑥ 独立扫描（MAI Code 2 Flash，≤3 用尽）

- r1：`INDEPENDENT SCAN: FINDINGS`，1 Medium —— 豁免只验闭合形态，外观围栏可躲 G6-4。
- 整改 r1：围栏体必须通过 G7-f 结构面（五段按序）；新增 T141 / T142。
- r2：`INDEPENDENT SCAN: FINDINGS`，1 Medium —— 标题行未 `.trim()`，与 G7-f 口径分叉（G7-f 绿而 G6-4 假红）。
- 整改 r2：`sixstageStructurallyValid` 逐行 `.trim()`；新增 T143。
- r3：`INDEPENDENT SCAN: PASS`，0 发现（原文逐字转录于报告对 §5 的 `sixstage` 围栏）。

## ⑦ 终审（GPT-5.3-Codex 常规档，1+2）

- r1（盲态）：`FINAL REVIEW: PASS`，0 发现，`verdict: no Medium+`。
- 独立取证读数：`--selftest` 757/757；tree `TOTAL_FAIL=0`；六形状全等；A/B 回归 `missing=0` / `FAIL=0`；两个最小反例分别为 PASS / FAIL；证据包三方集合一致、`hash_mismatch=0`；边界核查仅 7 个目标文件 + 证据目录。

## 批量辅助（GLM-5.3-Flash，2 次，载体 deep-reasoner）

- 双语条款镜像：`MIRROR: OK`（§0.3 行 / §6.3 句 / §11.3 句 / C.0bc 块四处语义等价）。
- 政策三处一致：`POLICY: CONSISTENT`（判据自身作用域、T23 零预载自洽、无任何处暗示可登记转录行）。
