# 归档 MANIFEST · dce

- **源路径**：`tmp/dce/`
- **归档片**：`DIRECTIVE-TMP-LIFECYCLE`（v1.23，2026-10-04）
- **归档口径**：§1.5.1 的 ② 类「早期证据评估归档」——保留证据条目（被受版本控制文档引用其路径者 ∪ 切片级文档），排除运行时残留与一次性脚本
- **归档条目**：11 条 / 215,448 字节（归档载荷根 = `archive/`；逐条 sha256 见 `SHA256SUMS.txt`）
- **排除条目**：6 条（逐条见下「排除清单」）
- **未入库条目**：0 条（逐条见下「未入库文件清单」）
- **字节保真**：归档副本的**字符内容不改**；因 §6.3 的作用域（G4a 扫全部变更文件）与仓库编码约定，副本按**归档后类型**归一编码——`.ps1` 带 BOM，其余不带 BOM，一律 LF；逐条变更种类与**原始字节 sha256** 见下「后缀与编码映射」节
- **后缀口径**：归档副本的 `.md` 一律改存为 `.txt`（本单元 0 条），以避开只扫 `.md` 的文档判据；映射逐条见下节
- **编码归一计数**：后缀改写 0 条 / BOM 剔除 0 条 / BOM 补入 0 条 / LF 归一 1 条

## 归档条目清单

| 相对路径 | 字节 | sha256 | 归档理由 |
|---|---|---|---|
| `archive/change-set-r2.diff` | 82071 | `e27b960a51456dbf9835bfc9193dbca2f1056a740d87c80fa51c953e110f3c86` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/change-set.diff` | 82071 | `e27b960a51456dbf9835bfc9193dbca2f1056a740d87c80fa51c953e110f3c86` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/gate-index.txt` | 3599 | `e6bc81e4096a5c2e8c9f33637f2e587f467e05517f8d983ae062e98a4c3960aa` | 被受版本控制文档引用其路径 |
| `archive/gate-tree-first-fail.txt` | 3602 | `afbfa3716fe42d58d8b08da0c2939616f6d33a8b821ff67707ba77f965499a71` | 被受版本控制文档引用其路径 |
| `archive/gate-tree.txt` | 3598 | `f4af90127d65bc5ac1c7d3be06df22cf6339aaa05f6f21c1506e962776b406be` | 被受版本控制文档引用其路径 |
| `archive/gate-tree2.txt` | 3596 | `cc088ec72c1d978d29e50968edcd12c74d7eb1c2400d122fbd2ea73d59470055` | 切片级文档 / 评审包 / 门禁读数 |
| `archive/land.js` | 9474 | `f0b507b522860f4b4b467b8b8dfe3e308b1d145d1cf4833285e1544dddb3efc6` | 被受版本控制文档引用其路径 |
| `archive/peek.js` | 864 | `7d4e127a4c16f74a7ad1dbb6942dd9c6b0593f7a3e28a0428ae8fdfa439d3f58` | 被受版本控制文档引用其路径 |
| `archive/probe-id-rejections.txt` | 2402 | `69a26c3490434566047c679a869357311bb18cf7756cd7e396cb870f8eddb781` | 被受版本控制文档引用其路径 |
| `archive/selftest.txt` | 22989 | `f367da5d25d72739455dc259670de21ab6e9e80b725701fe36928a0e13007282` | 被受版本控制文档引用其路径 |
| `archive/test.txt` | 1182 | `3975ce69bff0d2cb99004f11aa625f51a169bbe0269baea4eb9f2dac4f85e965` | 被受版本控制文档引用其路径 |

## 排除清单（未归档）

| 相对路径 | 字节 | sha256 | 未归档理由 |
|---|---|---|---|
| `/ci-run-37100757524.txt` | 158766 | `279e707c817fa062ad2d7191f566751cd74efe20c5217a9541c292e187c73a86` | 中间转储（未被引用的过程输出） |
| `/fence.js` | 1485 | `4098e699d5b05d1f937bafb62a59658a52a8b283a2c081bca3cc23df058e0cc3` | 落地与修复脚本残留 |
| `/fix-r7.js` | 13305 | `720cf890bbf42aefe030f7c52c0b2543eee9b540d1332f4a7a3d0567ad38a63b` | 落地与修复脚本残留 |
| `/fix-r7b.js` | 9120 | `bf815bde7fada90c52e8ae759f310598b39b877cfc1d6ead92a3416fc3ccfe1d` | 落地与修复脚本残留 |
| `/fix-r8.js` | 4261 | `4e8d81953635dcb6b74c183715fd19e9e5c9ebfae3ef705163d5d3b8039c5888` | 落地与修复脚本残留 |
| `/mk-report.js` | 23590 | `1d21f6876004c5c876f1b2c2158ef6e7637c2b39f4ede674b973ba22efff162e` | 落地与修复脚本残留 |

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
| `test.txt` | `archive/test.txt` | `67d63fb48d6a23c98488e7970a2d0521c24047a06a1afd65661d8713b08f1421` | `3975ce69bff0d2cb99004f11aa625f51a169bbe0269baea4eb9f2dac4f85e965` | LF 归一 |

## 对应引用点

| 引用文件 | 行号 |
|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 121 / 897 |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 121 / 897 |
| `docs/validation/directive-cleanup-eval.md` | 53 / 160 / 161 / 162 / 163 / 164 / 165 / 166 / 167 |
| `docs/validation/directive-cleanup-eval_EN.md` | 42 / 160 / 161 / 162 / 163 / 164 / 165 / 166 / 167 |
| `docs/validation/evidence/directive-history/dce/archive/change-set-r2.diff` | 205 / 310 / 311 / 312 / 313 / 314 / 315 / 316 / 317 / 366 / 482 / 483 / 484 / 485 / 486 / 487 / 488 / 489 |
| `docs/validation/evidence/directive-history/dce/archive/change-set.diff` | 205 / 310 / 311 / 312 / 313 / 314 / 315 / 316 / 317 / 366 / 482 / 483 / 484 / 485 / 486 / 487 / 488 / 489 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/classification.txt` | 478 / 479 / 480 / 481 / 482 / 483 / 484 / 485 / 486 / 487 / 488 / 489 / 490 / 491 / 492 / 493 / 494 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/facts2.txt` | 12 / 43 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/inventory.txt` | 105 / 106 / 107 / 108 / 109 / 110 / 111 / 112 / 446 / 447 / 448 / 471 / 472 / 473 / 501 / 502 / 503 / 504 / 505 / 506 / 507 / 508 / 509 / 510 / 511 / 512 / 513 / 514 / 515 / 516 / 517 / 518 |
| `docs/validation/evidence/directive-history/tmp-lifecycle-workdir/archive/measure.txt` | 44 |
