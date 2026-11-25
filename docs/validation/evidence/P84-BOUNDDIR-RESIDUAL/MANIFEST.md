# MANIFEST — P84-BOUNDDIR-RESIDUAL 证据包

切片：`BOUNDDIR-RESIDUAL`（P84，实现类 `[SLICE]`，M 档）。**定稿**（2026-11-21；结论 `PARTIALLY CLOSED`，见报告对 `docs/validation/bounddir-residual{,_EN}.md`；⑧c 前只增不删）。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：每次新增载荷后按磁盘派生重建 `SHA256SUMS.txt`，路径相对本目录。

**归档来源（P89，2026-11-25）**：本包 `archive/P84-CI-FIX/**` 为旧临时目录（`tmp` 根下 `P84-CI-FIX`）的归档副本，原目录已清除；逐件的「原始字节 sha256 / 归档字节 sha256 / 处置」见 `docs/validation/evidence/P89-TMP-ARCHIVE-P84-P85/archive-byte-provenance.tsv`，引用改写对照见 `docs/validation/tmp-archive-p84-p85{,_EN}.md`。

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
| `archive/P84-CI-FIX/01-r2-contract-judgment.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/02-r2-verification-matrix.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/03-r1-design.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/04-ci-rerun-evidence.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/07-narrow-review-request.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/08-narrow-re-review.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/09-narrow-review-request-r2.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/10-final-review-request.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/11-narrow-re-review-r2.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/12-rectification-record.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/13-final-narrow-review.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/14-rectification-A.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/15-final-review-A.md` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/attempt1-full-log.zip` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/attempt1-shortname-anchors.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/attempt1-test-fail-list-normalized.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/attempt1-unzip/0_Go windows-latest.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/attempt1-unzip/1_Go ubuntu-latest.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/attempt1-unzip/Go ubuntu-latest/system.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/attempt1-unzip/Go windows-latest/system.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/attempt1-windows-failed-log-excerpt.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/attempt2-windows-failed-log-excerpt.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/full-go122-r4.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/full-go127-r4.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/full-suite-go122-r2.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/full-suite-go122-r3.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/full-suite-go122.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/full-suite-go127-r2.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/full-suite-go127-r3.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/full-suite-go127.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/gates-index-r4.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/gates-index-r6.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/gates-tree-r4.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/gates-tree-r6.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/r1-probe-output-go127.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/r1-probe-output.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/r1r2-final-diff-r2.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/r1r2-final-diff-r3.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/r1r2-final-diff-r4.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/r1r2-final-diff-r5.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/r1r2-final-diff-r6.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/r1r2-final-diff.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/w12-go122-r4.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |
| `archive/P84-CI-FIX/w12-go127-r4.txt` | 归档原件（P89 归档；处置见 P89 包 `archive-byte-provenance.tsv`） |

## 自洽性

- 载荷条目数 = 磁盘实测（由脚本派生，不写固定数）；
- 逐条 sha256 由 `SHA256SUMS.txt` 复算；
- 元数据两件不列入载荷集，避免自指哈希。

