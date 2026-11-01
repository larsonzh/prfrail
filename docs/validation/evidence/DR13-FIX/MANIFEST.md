﻿﻿﻿﻿# MANIFEST — DR13-FIX 证据包

切片：`DR13-FIX`（修订类修复片 `[SLICE]`；P45 前置修复片，不占 v34 队列 1–4 编号）。建包时点 2026-11-01（`HEAD` = `8e834b2`，准则 v1.60）。

**包状态**：步骤 0 骨架落盘；后续载荷（01 / 03 / 05 / 06 / 07 / 08）按轮追加，并以一次性脚本重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息原文（含 DR-13 裁定与定位声明） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，含载体铁律与升级条件） |
| `01-diagnosis/review.md` | ① 定性结论原文（`DIAGNOSIS: B`，含证据链、W1/W2 交错时序、D1–D3 修法与 P45 输入） |
| `05-prereview/review.md` | ⑤ 预审 r1 + r2 输出原样（r1 2 Medium → r2 PASS） |
| `06-scan/review.md` | ⑥ 独立扫描输出原样（Section A–E，§B 零缺陷） |
| `07-final/review.md` | ⑦ 终审输出原样（`FINAL REVIEW: PASS` + 放行条件 6 条） |
| `03-tests/red-baseline.md` | ③ 红测基线原样：新增反例在未修复代码上确定性红（Windows 失败点 `spawn_recovery_test.go:281`） |

## 待追加（按轮）

| 路径 | 说明 |
|---|---|
| `08-ci/ci-evidence.txt` | ⑨ CI 取证 |

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）。
