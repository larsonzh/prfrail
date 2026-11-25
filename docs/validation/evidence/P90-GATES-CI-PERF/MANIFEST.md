# MANIFEST — P90-GATES-CI-PERF 证据包

切片：`GATES-CI-PERF`（P90，工具片 `[SLICE]`，S 档）。**已闭环**（2026-11-25 起跑并收口；⑧c 已清理 `tmp/`）。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：新增载荷后按磁盘派生重建 `SHA256SUMS.txt`，路径相对本目录（POSIX 分隔）。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 起跑消息逐字原文与四要素（§13.1.10） |
| `01-design.md` | ① 方案书 §A–§J（含 ⑤ r1 后的更正注记） |
| `02-coverage-equivalence.md` | base 四件套覆盖等价论证、gofmt 处置、判据面与 TOTAL_FAIL 口径 |
| `03-verification/ab-A-all.txt` | ③ A 侧输出（--all，含四条 base 记录） |
| `03-verification/ab-B-skipbase.txt` | ③ B 侧输出（--all --skip-base） |
| `03-verification/ab-report.json` | ③ A/B 比对结论（identical=true、双侧同 sha256） |
| `03-verification/ab-warmup.txt` | ③ 预热轮输出（仅用于消除 G4b 缓存态差异） |
| `03-verification/ci-run-45659843522-legs.md` | CI run 45659843522 双腿步骤结论、Gates 硬读数与耗时对照 |
| `03-verification/comparator-sanity.json` | ③ 比对器可证伪性自检 |
| `03-verification/gate-ci-final.txt` | ⑨ 定稿态 ci 读数 |
| `03-verification/gate-index-final.txt` | ⑨ 定稿态 index 读数 |
| `03-verification/gate-tree-4.txt` | ④ 第一轮 tree 读数（证据包未建，G8-a 预期红） |
| `03-verification/gate-tree-final.txt` | ⑨ 定稿态 tree 读数（含 base 四条） |
| `03-verification/selftest-final.txt` | ⑨ 定稿态 SELFTEST 读数 |
| `03-verification/selftest.txt` | 步骤 0 SELFTEST 原始读数 |
| `04-reviews/05-pre-review-r1.md` | ⑤ 第 1 轮原文（PASS WITH FIXES） |
| `04-reviews/05-re-review-r2.md` | ⑤ 第 2 轮原文（RE-REVIEW: PASS WITH FIXES） |
| `04-reviews/06-scan-r1.md` | ⑥ 第 1 轮原文（FINDINGS） |
| `04-reviews/06-scan-r2.md` | ⑥ 第 2 轮原文（PASS） |
| `04-reviews/07-final-review-r1.md` | ⑦ 第 1 轮原文（FINDINGS） |
| `04-reviews/07-re-review-r2.md` | ⑦ 第 2 轮原文（RE-REVIEW: PASS） |
| `freedom-list.md` | 执行自由度清单（F-01~F-11，含代裁项） |
| `input-list.md` | 输入清单（步骤 0 锚点实读、改动前 CI 基线、审查输入包） |
