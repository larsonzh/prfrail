# MANIFEST — P40-T027B-LOCK-REAPER 证据包

切片：`P40-T027B-LOCK-REAPER`（P40，实现类 `[IMPL]`，C1 死锁与锁序）。建包时点 2026-10-31（`HEAD` = `38a5cba`，准则 v1.60）。

**包状态**：步骤 0 骨架 + ① 架构与批量辅助 #1 已落盘；② 实现两轮（含方案 A 与 F-01..F-06 修复）后，① 文件追加 F-07b 与措辞收口注记；⑤ 两轮、⑥、⑦ 输出原样归档；SHA256SUMS 由 `tmp/P40-T027B-LOCK-REAPER/rehash.js` 重算。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（当前 5 件）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-lock-census.md` | 批量辅助 #1《现有锁点清单》（GLM 枚举 internal/ 全部 Mutex/Lock/Unlock 与文件锁） |
| `01-arch-design/review.md` | ① 架构《C1 锁序与收割重构方案》（V4 Pro，主控逐字转录；含 F-07b 与措辞收口注记） |
| `05-prereview/review.md` | ⑤ 预审两轮输出原样（r1 `FINAL REVIEW: PASS WITH FIXES`；r2 `PRE-REVIEW: PASS`） |
| `06-scan/review.md` | ⑥ 独立扫描输出原样（`INDEPENDENT SCAN: PASS`） |
| `07-final/review.md` | ⑦ 终审输出原样（`FINAL REVIEW: PASS`） |
| `freedom-list.md` | 执行自由度清单（本片无申请，落盘以满足 §5.1 时序） |

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）；由 `tmp/P40-T027B-LOCK-REAPER/` 一次性脚本实测生成，脚本按 §5.1 在 ⑧c 清除。
