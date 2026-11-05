# MANIFEST — G4B-CORPUS-DECOUPLE 证据包

切片：`G4B-CORPUS-DECOUPLE`（判据优化类 `[SLICE]`；M 档；层 3 备选·挂起项，用户 2026-11-05 同轮裁定激活）。建包时点 2026-11-05（`HEAD` = `f021af1`，准则 v1.60）。

**包状态**：步骤 0 已按正式起跑消息完成（2026-11-05 起跑）。② 实现、③ 测试、④ 门禁、⑤⑥⑦ 审查与 ⑧ 文档的产物按轮追加到下表，并同步重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（`MANIFEST` 归档条目集合 ≡ `SHA256SUMS` 条目集合 ≡ 包内载荷文件集合）。件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数。

**临时件**：本片临时脚本与中间产物一律落 `tmp/G4B-CORPUS-DECOUPLE/`（⑧c 统一清理）；本包内不含临时件。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户 2026-11-05 起跑消息逐字转录（含【〇·零】至【九】全部小节；围栏承载原文，围栏外为性质说明） |
| `03-test-report.md` | ③ 测试报告（新增 T114–T118 的测试点与变异敏感性、对既有断言的必要调整、四阶段 `SELFTEST` 读数、覆盖缺口） |
| `04-reviews/05-pre-review-r1.md` | ⑤ 预审第 1 轮原文（`PRE-REVIEW: PASS WITH FIXES`；3 项 Low，无 Medium+） |
| `04-reviews/05-pre-review-r2.md` | ⑤ 预审第 2 轮原文（整改后复审；`PRE-REVIEW: PASS WITH FIXES`；无 Medium+，1 项 Low = 读数叙述滞后） |
| `04-reviews/06-independent-scan.md` | ⑥ 独立扫描原文（Section A–E；Section B = `(no findings)`；`INDEPENDENT SCAN: PASS`） |
| `04-reviews/07-final-review-r1.md` | ⑦ 独立终审第 1 轮原文（1 High + 1 Medium + 1 Low；`FINAL REVIEW: FINDINGS`） |
| `04-reviews/07-final-review-r2.md` | ⑦ 整改后复审第 1 轮原文（1 High + 1 Low；`RE-REVIEW: FINDINGS`） |
| `04-reviews/07-final-review-r3.md` | ⑦ 整改后复审第 2 轮原文（`(no findings)`；三方计数 8/8/8、行数抽查 3/3；`RE-REVIEW: PASS`） |
| `input-list.md` | 输入包清单、锚点核对、G4b 判据文本与实现口径的差异对照、影响面（批量辅助第 1 次结论 + 主控复核）、改动前基线读数 |
| `freedom-list.md` | 执行自由度清单（步骤 0 起草；⑥ 后 ⑦ 前定稿）+ 不派发环节 + 纪律登记 + 升级边界 |

后续环节的产物按轮追加到上表并同步重算 `SHA256SUMS.txt`。
