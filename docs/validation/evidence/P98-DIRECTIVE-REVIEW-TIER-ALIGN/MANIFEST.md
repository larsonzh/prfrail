# P98 证据包 MANIFEST

切片：`DIRECTIVE-REVIEW-TIER-ALIGN`（P98）。日期：2026-11-28。

本包含 `MANIFEST.md` 与 `SHA256SUMS.txt` 两件元数据，其余为载荷。载荷在步骤 0 与后续阶段逐次并入；每次并入均同轮更新本清单与哈希（§13.1.12）。

## 归档条目清单

| 载荷 | 内容 |
|---|---|
| `00-briefing.md` | 起跑消息**逐字原文**（§13.1.10） |
| `input-list.md` | 输入清单与步骤 0 锚点实测 |
| `freedom-list.md` | 执行自由度清单（§1.7.3；F-1..F-12 待 ⑧b 终态回写） |
| `01-arch/01-arch-design.md` | ① 架构方案书（修订文本三件套 + 升版 / OB 处置 / ADR 判定；停点 1 四项候选） |
| `02-ruling-stop1.md` | 停点 1 用户裁定（逐字原文 + 结构化登记 + 与 01-arch 拟文差异） |
| `04-reviews/05-pre-review.txt` | ⑤ 预审报告逐字（PRE-REVIEW: PASS；.txt 归档） |
| `04-reviews/08-review-7-r1.txt` | ⑦ 终审第 1 轮原文逐字（FINAL REVIEW: FINDINGS；.txt 归档） |
| `04-reviews/09-scan-6.txt` | ⑥ 独立扫描报告逐字（补跑；INDEPENDENT SCAN: FINDINGS；.txt 归档） |
| `04-reviews/10-review-7-r2.txt` | ⑦ 复审第 2 轮原文逐字（RE-REVIEW: FINDINGS；.txt 归档） |

## 编码与行尾

载荷 `.md` 一律 UTF-8 with BOM + LF；哈希在编码归一之后计算。
