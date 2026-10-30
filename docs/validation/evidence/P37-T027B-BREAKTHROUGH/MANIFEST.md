# MANIFEST — P37-T027B-BREAKTHROUGH 证据包

切片：`P37-T027B-BREAKTHROUGH`（P37，评估类 `[SLICE]`）。建包时点 2026-10-30（`HEAD` = `00ddc53`，准则 v1.58）。

**包状态**：分轮追加包 —— 每轮载荷落盘后按 `tmp/P37-T027B-BREAKTHROUGH/step0-sums.js` 重算 `SHA256SUMS.txt` 并与本清单同步，⑧c 前保持自洽。

**载荷集**：本包根 `MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，余下全部文件为载荷；清单条目集 = `SHA256SUMS.txt` 条目集 = 载荷集（当前 6 件）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-doc-inventory.md` | 文档清单与输入包摘要 |
| `01-sol-design/input-list.md` | 第 1 轮（架构设计）输入清单与字节数 |
| `01-sol-design/review.md` | 第 1 轮《T027-B 生产级架构设计书》（主控逐字转录） |
| `02-astra-critique/input-list.md` | 第 2 轮（批判与重构）输入清单与字节数 |
| `02-astra-critique/review.md` | 第 2 轮《T027-B 设计批注与重构方案》（主控逐字转录） |
| `freedom-list.md` | 执行自由度清单（本片无申请） |

## 未入库文件清单

无（本包当前无二进制或忽略类型载荷）。

## 重建说明

载荷全部为可直接重写的纯文本（UTF-8 + BOM + LF）；清单与字节数由 `tmp/P37-T027B-BREAKTHROUGH/` 的一次性脚本实测生成，脚本按 §5.1 在 ⑧c 清除，不属载荷。
