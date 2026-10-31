# MANIFEST — P41-T027B-SPAWN-OWNERSHIP 证据包

切片：`P41-T027B-SPAWN-OWNERSHIP`（P41，实现类 `[SLICE]`，C2+N1 孤儿进程与 Job 归属）。建包时点 2026-10-31（`HEAD` = `28f440f`，准则 v1.60）。

**包状态**：步骤 0 骨架 + ① 架构已落盘；⑤ 两轮、⑥、⑦ 输出原样归档；SHA256SUMS 由 `tmp/P41-T027B-SPAWN-OWNERSHIP/rehash.js` 重算。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（当前 5 件）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `01-arch-design/review.md` | ① 架构《C2+N1 归属与防孤儿设计书》（V4 Pro，主控逐字转录） |
| `05-prereview/review.md` | ⑤ 预审两轮输出原样（r1 `PRE-REVIEW: PASS WITH FIXES`；r2 `PRE-REVIEW: PASS`） |
| `06-scan/review.md` | ⑥ 独立扫描输出原样（`INDEPENDENT SCAN: PASS`） |
| `07-final/review.md` | ⑦ 终审输出原样（`FINAL REVIEW: PASS`） |
| `freedom-list.md` | 执行自由度清单（本片无申请，落盘以满足 §5.1 时序） |

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）；由 `tmp/P41-T027B-SPAWN-OWNERSHIP/` 一次性脚本实测生成，脚本按 §5.1 在 ⑧c 清除。
