# MANIFEST — GOVERNANCE-BASELINE-VERIFY 证据包

切片：`GOVERNANCE-BASELINE-VERIFY`（P47，M 档评估/治理类 `[SLICE]`）。建包时点 2026-11-02（`HEAD` = `0e46bd0`，准则 v1.60）。

**包状态**：步骤 0 骨架 + ① 基线核验载荷已落盘（R-11 经用户裁定为 B 档后执行）；② 补登记载荷按阶段追加，并以一次性脚本重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户起跑消息逐字原文（3,087 bytes，BOM + LF，`sha256=7a6992df…4d65`） |
| `01-baseline/correction.md` | G4 / G5 / G7 三段更正（4,031 bytes，BOM + LF，`sha256=29da66a4…0aa5`） |
| `01-baseline/verification.md` | 七项基线核验结论表与检索范围（3,632 bytes，BOM + LF，`sha256=903fbe80…07c8`） |
| `02-backfill/registration.md` | 补登记落地记录（OB-133 至 OB-136 与三行指针追加；2,839 bytes，BOM + LF，`sha256=52ec43c9…58ad`） |
| `freedom-list.md` | 执行自由度清单（5 项行使 / 3 项明确不做；含 R-11 未裁时的冻结声明） |
| `input-list.md` | 输入包清单（10 项，逐件字节数，OB-93 口径） |

## 待追加（按阶段）

| 路径 | 说明 |
|---|---|
| `08-ci/ci-evidence.txt` | ⑨ CI 取证 |

## 未入库文件清单

无。

## 重建说明

载荷为 UTF-8 + BOM + LF 纯文本（`SHA256SUMS.txt` 无 BOM）；`.yml`、Go 源码与 CI 工件不复制入包，仅登记路径。
