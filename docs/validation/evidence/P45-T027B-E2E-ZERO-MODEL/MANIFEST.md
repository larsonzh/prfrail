﻿﻿# MANIFEST — P45-T027B-E2E-ZERO-MODEL 证据包

切片：`P45-T027B-E2E-ZERO-MODEL`（M 档装配验证片 `[SLICE]`）。建包时点 2026-11-01（`HEAD` = `eb36992`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；后续载荷（01 / 03 / 05 / 06 / 07 / 08）按轮追加，并以一次性脚本重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（含 DR13-FIX 强输入与文档回写要求） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含并发断言与零模型证明要求） |
| `01-arch-design/review.md` | ① 架构《零模型三任务 E2E 装配方案》原文（含落点取舍、红测清单、并发断言设计、风险） |
| `03-tests/red-baseline.md` | ③ 首轮红测基线原样（7 绿 1 红；红点=夹具缺失） |
| `05-prereview/review.md` | ⑤ 预审输出原样（2 MINOR，均已处置） |
| `06-scan/review.md` | ⑥ 独立扫描输出原样（§B 零缺陷） |
| `07-final/review.md` | ⑦ 终审输出原样（`FINAL REVIEW: PASS` + 6 条放行条件） |

## 待追加（按轮）

| 路径 | 说明 |
|---|---|
| `08-ci/ci-evidence.txt` | ⑨ CI 取证 |

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）。
