# MANIFEST — P39-T027B-DURABLE-BUDGET 证据包

切片：`P39-T027B-DURABLE-BUDGET`（P39，实现类 `[IMPL]`，C6 预算落盘）。建包时点 2026-10-31（`HEAD` = `8b16e21`，准则 v1.59）。

**包状态**：① 架构设计分两批（批次 A 设计书 / 批次 B 反例矩阵）；载荷随 ②③④ 追加以一次性脚本重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（当前 4 件）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `01-arch-design/review-batch-b.md` | ① 批次 B《反例测试矩阵》（主控逐字转录） |
| `01-arch-design/review.md` | ① 批次 A《持久化架构设计书》（主控逐字转录） |
| `02-implement/voided-v4pro-waste.md` | 模型绑定偏离留证：误派 V4 Pro 的四件产物 sha256 与作废处置 |
| `freedom-list.md` | 执行自由度清单（本片无申请，落盘以满足 §5.1 时序） |

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）；由 `tmp/P39-T027B-DURABLE-BUDGET/` 一次性脚本实测生成，脚本按 §5.1 在 ⑧c 清除。
