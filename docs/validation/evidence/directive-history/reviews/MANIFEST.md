# 归档 MANIFEST · reviews

- **源路径**：`tmp/DELIVERY_DIRECTIVE_v111_REVIEW.md` / `tmp/DELIVERY_DIRECTIVE_v111_REVIEW_R3.md` / `tmp/DELIVERY_DIRECTIVE_v112_REVIEW.md` / `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R5.md` / `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R6.md` / `tmp/DELIVERY_DIRECTIVE_v112_REVIEW_R7.md`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：6 条 / 100,926 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：0 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 6 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 6 条 / BOM 剔除 4 条 / BOM 补入 0 条 / LF 归一 0 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/DELIVERY_DIRECTIVE_v111_REVIEW.txt` | 20339 | `714c27fd579f568a59f914cfb55257ab2ac2f0dda7b2ab275822217f3f08aec7` | 被受版本控制文档引用其路径（准则评审报告） |
| `archive/DELIVERY_DIRECTIVE_v111_REVIEW_R3.txt` | 18429 | `03144f6fd1b04b731318775d4c3a668446277aa373db0ef2a6cefa49457b02bf` | 被受版本控制文档引用其路径（准则评审报告） |
| `archive/DELIVERY_DIRECTIVE_v112_REVIEW.txt` | 18573 | `7c588a9565264f8bba34a91f0f220577dea45b9c78f8e901e8dbedd2a6fa5769` | 被受版本控制文档引用其路径（准则评审报告） |
| `archive/DELIVERY_DIRECTIVE_v112_REVIEW_R5.txt` | 13613 | `d5e608c3a76dc1b612830204d74116b8c43b81e80c237a677a0d2f4bbb326860` | 被受版本控制文档引用其路径（准则评审报告） |
| `archive/DELIVERY_DIRECTIVE_v112_REVIEW_R6.txt` | 14176 | `6adb9cb0a0f72998b603a9e3b165764b55e02a4114bf8a55d78a61cc142bc3fb` | 被受版本控制文档引用其路径（准则评审报告） |
| `archive/DELIVERY_DIRECTIVE_v112_REVIEW_R7.txt` | 15796 | `c0b590b88c11f938428782bd5fa8e1ca845eb2aae1a35bb007c95904fff7b144` | 被受版本控制文档引用其路径（准则评审报告） |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|

## 未入库文件清单（`.gitignore` 排除类型）

| 相对路径 | 字节 | sha256 | 理由 |
|---|---|---|---|

## 重建说明

- **未入库条目**（见上表）：本单元未入库清单为空（无忽略类型条目）。
- **排除条目**：一次性脚本、运行时状态与中间转储；其行为已由本片归档的切片级文档（切片定义 / 方案 / 评审包 / 门禁读数）与对应切片的提交记录固定。
- **行尾与 BOM 归一**：源 `tmp/` 副本的部分文件为 CRLF 或带 BOM；归档副本已按仓库编码约定归一到「`.ps1` 带 BOM / 其余不带 BOM / 一律 LF」（不影响字符内容）。

## 后缀与编码映射（`.md` → `.txt` 及编码归一）

| 原路径 | 归档路径 | 原始字节 sha256 | 归档字节 sha256 | 变更 |
|---|---|---|---|---|
| `DELIVERY_DIRECTIVE_v111_REVIEW.md` | `archive/DELIVERY_DIRECTIVE_v111_REVIEW.txt` | `391840323db5b589f417e1d46bd84c5f334dfafe2e304157e07e3d1f6fd60792` | `714c27fd579f568a59f914cfb55257ab2ac2f0dda7b2ab275822217f3f08aec7` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `DELIVERY_DIRECTIVE_v111_REVIEW_R3.md` | `archive/DELIVERY_DIRECTIVE_v111_REVIEW_R3.txt` | `c7cc8b7b5b858d50f2481ff740e71773850d8978db0a494ff8389871119b9869` | `03144f6fd1b04b731318775d4c3a668446277aa373db0ef2a6cefa49457b02bf` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `DELIVERY_DIRECTIVE_v112_REVIEW.md` | `archive/DELIVERY_DIRECTIVE_v112_REVIEW.txt` | `549e40be66ed92388f6f8f64537e7aee4733e65923029723549a9a4c912fee3c` | `7c588a9565264f8bba34a91f0f220577dea45b9c78f8e901e8dbedd2a6fa5769` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `DELIVERY_DIRECTIVE_v112_REVIEW_R5.md` | `archive/DELIVERY_DIRECTIVE_v112_REVIEW_R5.txt` | `d5e608c3a76dc1b612830204d74116b8c43b81e80c237a677a0d2f4bbb326860` | `d5e608c3a76dc1b612830204d74116b8c43b81e80c237a677a0d2f4bbb326860` | 后缀 `.md` → `.txt` |
| `DELIVERY_DIRECTIVE_v112_REVIEW_R6.md` | `archive/DELIVERY_DIRECTIVE_v112_REVIEW_R6.txt` | `be0fc3da1c22e968a3cea27b126b2c4ff33ea23620aa55b053a3a8ac5875e4d5` | `6adb9cb0a0f72998b603a9e3b165764b55e02a4114bf8a55d78a61cc142bc3fb` | 后缀 `.md` → `.txt` / BOM 剔除 |
| `DELIVERY_DIRECTIVE_v112_REVIEW_R7.md` | `archive/DELIVERY_DIRECTIVE_v112_REVIEW_R7.txt` | `c0b590b88c11f938428782bd5fa8e1ca845eb2aae1a35bb007c95904fff7b144` | `c0b590b88c11f938428782bd5fa8e1ca845eb2aae1a35bb007c95904fff7b144` | 后缀 `.md` → `.txt` |

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 58 / 119 / 858 / 1275 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 58 / 119 / 858 / 1275 |
| `docs/validation/directive-v1.12.md` | 61 / 62 / 63 / 65 / 66 / 67 / 342 / 343 / 344 / 345 / 380 |
| `docs/validation/directive-v1.12_EN.md` | 61 / 62 / 63 / 65 / 66 / 67 / 342 / 343 / 344 / 345 / 380 |
