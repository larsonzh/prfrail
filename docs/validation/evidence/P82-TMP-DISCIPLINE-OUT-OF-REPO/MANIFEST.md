# MANIFEST — P82-TMP-DISCIPLINE-OUT-OF-REPO 证据包

切片：`TMP-DISCIPLINE-OUT-OF-REPO`（P82，§5.1 临时件纪律补处置通道 + `tmp/` 根自愈口径，治理片 `[SLICE]`，M 档）。建包时点 2026-11-19（`HEAD` = `cd790b6`，准则 v1.70，本片升至 v1.71）。

**包状态**：步骤 0 建包三件初始载荷；⑧c 收口时归档 `01-arch/`、`02-batch/`、`04-reviews/`、`06-gates/` 并重算 `SHA256SUMS.txt`。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 载荷（P82 证据） |
| `01-arch/design-r1.md` | 载荷（P82 证据） |
| `02-batch/batch-01-instances.md` | 载荷（P82 证据） |
| `02-batch/batch-02-mirror.md` | 载荷（P82 证据） |
| `04-reviews/05-prereview-r1.md` | 载荷（P82 证据） |
| `04-reviews/05-prereview-r2.md` | 载荷（P82 证据） |
| `04-reviews/05-prereview-r3.md` | 载荷（P82 证据） |
| `04-reviews/06-scan-r1.md` | 载荷（P82 证据） |
| `04-reviews/06-scan-r2.md` | 载荷（P82 证据） |
| `04-reviews/07-final-r1.md` | 载荷（P82 证据） |
| `04-reviews/07-final-r2.md` | 载荷（P82 证据） |
| `06-gates/gate-ci.txt` | 载荷（P82 证据） |
| `06-gates/gate-index.txt` | 载荷（P82 证据） |
| `06-gates/gate-selftest.txt` | 载荷（P82 证据） |
| `06-gates/gate-tree.txt` | 载荷（P82 证据） |
| `freedom-list.md` | 载荷（P82 证据） |
| `input-list.md` | 载荷（P82 证据） |

## 编码

全部载荷为 `.md` / `.txt`：`.md` 为 UTF-8 with BOM + LF、`.txt` 为 UTF-8 without BOM + LF（§6.1 编码行 / §6.3 G4a）。哈希在编码归一之后计算。
