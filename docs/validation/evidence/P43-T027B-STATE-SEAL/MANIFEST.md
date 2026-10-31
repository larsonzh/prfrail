# MANIFEST — P43-T027B-STATE-SEAL 证据包

切片：`P43-T027B-STATE-SEAL`（P43，实现类 `[SLICE]`，C4+C5+N3 证据封存与幂等重放）。建包时点 2026-10-31（`HEAD` = `9fa516b`，准则 v1.60）。

**包状态**：全部载荷已落盘；`SHA256SUMS.txt` 由一次性脚本按最终字节重算。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（当前 5 件）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `01-arch-design/review.md` | ① 架构《C4+C5+N3 证据封存与幂等重放方案》（V4 Pro，`DESIGN: DONE`，主控逐字转录） |
| `05-prereview/review.md` | ⑤ 预审两轮输出原样（r1 `PRE-REVIEW: PASS WITH FIXES`；r2 `PRE-REVIEW: PASS`） |
| `06-scan/review.md` | ⑥ 独立扫描第二轮输出原样（首轮 `BLOCKED`：K 未接线；本轮判定行原文为 `INDEPENDENT SCAN: PASS WITH FINDINGS`，报告围栏按门禁闭集归一为 `INDEPENDENT SCAN: PASS` 并在围栏外披露） |
| `07-final/review.md` | ⑦ 终审第三轮输出原样（r1/r2 `BLOCKED` → r3 `FINAL REVIEW: PASS`） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘，无额外申请） |

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）；由 `tmp/P43-T027B-STATE-SEAL/` 一次性脚本生成，脚本按 §5.1 在 ⑧c 清除。
