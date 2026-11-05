# MANIFEST — P58-OB-TRIAGE 证据包

切片：`P58-OB-TRIAGE`（评估类 `[SLICE]`；S 档；层 2 的 `G-1` 项，用户 2026-11-05 点名起跑）。2026-11-05 步 0 建包（HEAD = `be1e89d`，准则 v1.62；该版本于 2026-11-05 经用户 ⑩ 追认）。

**包状态**：步骤 0、评估分析、⑤ 预审与 ⑨ 门禁取证已落盘；后续轮次的审查原文与终版读数按轮追加到下表，并同步重算 `SHA256SUMS.txt`。

**载荷集**：`MANIFEST.md` 与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（`MANIFEST` 归档条目集合 ≡ `SHA256SUMS` 条目集合 ≡ 包内载荷文件集合）。件数以 `SHA256SUMS.txt` 的哈希行为准，不设固定数。

**临时件**：本片临时脚本与中间产物一律落 `tmp/P58-OB-TRIAGE/`（⑧c 统一清理）；本包内不含临时件。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 用户 2026-11-05 起跑消息逐字转录（含【〇·零】至【六】全部小节） |
| `04-reviews/05-pre-review-r1.md` | ⑤ 预审 r1 逐字原文（`PRE-REVIEW: FINDINGS`；4 Critical + 2 Medium + 4 Low） |
| `04-reviews/05-pre-review-r2.md` | ⑤ 预审 r2 逐字原文（整改后复审；`PRE-REVIEW: PASS WITH FIXES`；r1 十条全部 CLOSED + 3 项 Low） |
| `04-reviews/05-pre-review-r3.md` | ⑤ 预审 r3 逐字原文（用户同轮指定的收窄读；六项修正点全部 CLOSED；`PRE-REVIEW: PASS`） |
| `04-reviews/05-pre-review-r4.md` | ⑤ 预审 r4 逐字原文（用户裁定后八项核验全部 CLOSED；`PRE-REVIEW: PASS`） |
| `06-gates/gate-readings.txt` | ⑨ 期 `--selftest` 与三作用域 `--all` 的原始输出（含退出码） |
| `freedom-list.md` | 执行自由度清单（步骤 0 落盘；含未行使选项留痕） |
| `input-list.md` | 审查输入清单、输入一致性断言、旧口径检索记录 |

后续环节的产物按轮追加到上表并同步重算 `SHA256SUMS.txt`。
