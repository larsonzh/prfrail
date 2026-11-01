﻿﻿﻿# MANIFEST — STRATEGIC-REVIEW-2026-11-02 证据包

切片：`STRATEGIC-REVIEW-2026-11-02`（L 档评估类 `[SLICE]`，流水线裁剪：跳过 ①–⑦）。建包时点 2026-11-02（`HEAD` = `f497ff3`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；轮 1–3 产出与补偿复核、CI 取证按轮追加，并以一次性脚本重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户最终版起跑消息原文（含前言三项说明；15,900 bytes，BOM + LF，`sha256=86c4b37c…824b`） |
| `freedom-list.md` | 执行自由度清单（含补偿复核边界与输出长度防御） |
| `input-list.md` | 输入包清单（10 项，逐件字节数，OB-93 口径） |
| `01-round1/review.md` | 轮 1 双断点核验原文（12 项判定表；`断点 A = 部分闭合；断点 B = 未闭合`） |
| `02-round2/review.md` | 轮 2 缺口台账原文（76 项，2 片，片界标注保留） |
| `03-round3/review.md` | 轮 3 路径与优先级原文（三档序列；主链 = console 多任务装配扩展，其次 = `CONTRACTS-S7-DURABILITY`；先做三项） |
| `04-compensation/review.md` | 补偿复核原文（`COMPENSATION REVIEW: PASS`，仅格式/结构，不核战略判断质量） |

## 待追加（按轮）

| 路径 | 说明 |
|---|---|
| `08-ci/ci-evidence.txt` | ⑨ CI 取证 |

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM；`.yml` 与代码类输入不复制入包，仅登记路径与字节数）。
