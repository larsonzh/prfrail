# MANIFEST — P84-BOUNDDIR-RESIDUAL 证据包

切片：`BOUNDDIR-RESIDUAL`（P84，实现类 `[SLICE]`，M 档）。**定稿**（2026-11-21；结论 `PARTIALLY CLOSED`，见报告对 `docs/validation/bounddir-residual{,_EN}.md`；⑧c 前只增不删）。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：每次新增载荷后按磁盘派生重建 `SHA256SUMS.txt`，路径相对本目录。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 载荷（payload (start message verbatim, section 13.1.10)） |
| `input-list.md` | 载荷（payload (input list and anchor readings)） |
| `freedom-list.md` | 载荷（payload (execution freedom list)） |
| `01-arch/01-arch-design.md` | 载荷（payload ((1) architecture design, with Rev 1 + Rev 2 revision notes)） |
| `04-reviews/05-pre-review-r1.md` | 载荷（payload ((5) pre-review r1 verbatim)） |
| `04-reviews/05-pre-review-r2.md` | 载荷（payload ((5) pre-review r2 verbatim)） |
| `04-reviews/06-scan-r1.md` | 载荷（payload ((6) independent scan r1 verbatim)） |
| `04-reviews/06-scan-r2.md` | 载荷（payload ((6) independent scan r2 verbatim)） |
| `04-reviews/07-blind-review.md` | 载荷（payload ((7) blind review verbatim, mandatory draw)） |
| `04-reviews/07-final-review-r1.md` | 载荷（payload ((7) final review r1 verbatim)） |
| `04-reviews/07-re-review-r1.md` | 载荷（payload ((7) re-review r1 verbatim)） |
| `04-reviews/07-re-review-r2.md` | 载荷（payload ((7) re-review r2 verbatim)） |
| `04-reviews/07-re-review-r3.md` | 载荷（payload ((7) re-review r3 verbatim, final quota)） |
| `04-reviews/07-narrow-re-review.md` | 载荷（payload ((7) narrow re-review verbatim, user-ruled quota exception)） |

## 自洽性

- 载荷条目数 = 磁盘实测（由脚本派生，不写固定数）；
- 逐条 sha256 由 `SHA256SUMS.txt` 复算；
- 元数据两件不列入载荷集，避免自指哈希。

