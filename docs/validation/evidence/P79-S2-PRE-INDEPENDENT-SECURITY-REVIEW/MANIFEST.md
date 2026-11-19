# MANIFEST — P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW 证据包

切片：`S2-PRE-INDEPENDENT-SECURITY-REVIEW`（P79，ADR-047 补偿义务的独立安全评审补评，评估类 `[SLICE]`，M-L 档）。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：每次新增载荷后由脚本（`tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/build-pack.js`）按磁盘派生重建，不手工维护。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 载荷（P79 证据） |
| `01-arch/design-r1.md` | 载荷（P79 证据） |
| `01-arch/package-manifest.md` | 载荷（P79 证据） |
| `04-reviews/05-prereview-r1.md` | 载荷（P79 证据） |
| `04-reviews/05-prereview-r2.md` | 载荷（P79 证据） |
| `04-reviews/06-scan-r1.md` | 载荷（P79 证据） |
| `04-reviews/06-scan-r2.md` | 载荷（P79 证据） |
| `04-reviews/07-re-review-r1.md` | 载荷（P79 证据） |
| `04-reviews/07-re-review-r2.md` | 载荷（P79 证据） |
| `04-reviews/07-security-r1.md` | 载荷（P79 证据） |
| `06-gates/gate-enc-scan.txt` | 载荷（P79 证据） |
| `06-gates/gate-index-final.txt` | 载荷（P79 证据） |
| `06-gates/gate-selftest-final.txt` | 载荷（P79 证据） |
| `06-gates/gate-tree-final.txt` | 载荷（P79 证据） |
| `freedom-list.md` | 载荷（P79 证据） |
| `input-list.md` | 载荷（P79 证据） |

## 自洽性

- 载荷条目数 = 磁盘实测（由构建脚本派生，不写固定数）；
- 逐条 sha256 由 `SHA256SUMS.txt` 复算，脚本构建时即写入；
- 元数据两件不列入载荷集，避免自指哈希。
