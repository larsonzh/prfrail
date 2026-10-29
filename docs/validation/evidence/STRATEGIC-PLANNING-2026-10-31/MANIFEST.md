# MANIFEST — STRATEGIC-PLANNING-2026-10-31 证据包

切片：`STRATEGIC-PLANNING-2026-10-31`（P36，评估类 `[SLICE]`）。**草稿**：步骤 0 建包时点 2026-10-29（`HEAD` = `70cb193`），准则 v1.57。

**包状态**：分轮追加包 —— 每轮载荷落盘后按 `tmp/STRATEGIC-PLANNING-2026-10-31/step0-sums.js` 重算 `SHA256SUMS.txt` 并与本清单同步，⑧c 前保持自洽。

**载荷集**：本包根 `MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，余下全部文件为载荷；清单条目集 = `SHA256SUMS.txt` 条目集 = 载荷集（当前 8 件）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-doc-inventory.md` | 文档清单与继承事实（P35 报告对 + 两次 CI run 的可复核记录）＋ P36 全量文档清单 |
| `01-product-planning/input-list.md` | 第 1 轮（产品与架构规划）输入清单与字节数 |
| `01-product-planning/review.md` | 第 1 轮（产品与架构规划）评审正文（主控逐字转录，GPT-6.1 Sol） |
| `02-governance-planning/input-list.md` | 第 2 轮（治理与流程规划）输入清单与字节数 |
| `02-governance-planning/review.md` | 第 2 轮（治理与流程规划）评审正文（主控逐字转录，GPT-6.1 Sol） |
| `03-synthesis-astra/input-list.md` | 第 3 轮（综合裁决）输入清单——预期路径占位 |
| `03-synthesis-astra/review.md` | 第 3 轮（综合裁决）正文（主控逐字转录，GPT-6 Astra；OB-129） |
| `freedom-list.md` | 执行自由度清单（本片无申请） |

## 未入库文件清单

无（本包当前无二进制或忽略类型载荷）。

## 重建说明

载荷全部为可直接重写的纯文本（UTF-8 + BOM + LF）；清单与字节数由 `tmp/STRATEGIC-PLANNING-2026-10-31/` 的一次性脚本实测生成，脚本按 §5.1 在 ⑧c 清除，不属载荷。
