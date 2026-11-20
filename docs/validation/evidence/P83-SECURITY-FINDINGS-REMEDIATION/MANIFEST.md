# MANIFEST — P83-SECURITY-FINDINGS-REMEDIATION 证据包

切片：`SECURITY-FINDINGS-REMEDIATION`（P83，实现类 `[SLICE]`，M-L 档）。结论：`PARTIALLY CLOSED`。

**载荷集**：本清单与 `SHA256SUMS.txt` 为元数据，其余为载荷；条目集三方相等（磁盘 / 本清单 / `SHA256SUMS.txt`）。
**重算口径**：每次新增载荷后按磁盘派生重建 `SHA256SUMS.txt`，路径相对本目录。

## 归档条目清单

| 路径 | 说明 |
|---|---|
| `00-briefing.md` | 载荷（起跑消息原文） |
| `01-arch/remediation-plan.md` | 载荷（① 修复方案书） |
| `01-arch/ruling-triage.md` | 载荷（① 终审九项契约分流判定书） |
| `04-reviews/07-final-review.md` | 载荷（⑦ 终审逐字转录） |
| `04-reviews/07-re-reviews-r1-r3.md` | 载荷（⑦ 复审 r1/r2/r3 逐字转录） |
| `04-reviews/07-additional-review.md` | 载荷（⑦ 追加复审逐字转录，用户直接裁定、不得作为先例） |
| `freedom-list.md` | 载荷（执行自由度清单） |
| `input-list.md` | 载荷（输入清单与锚点实测） |

## 自洽性

- 载荷条目数 = 磁盘实测（由脚本派生，不写固定数）；
- 逐条 sha256 由 `SHA256SUMS.txt` 复算；
- 元数据两件不列入载荷集，避免自指哈希；
- ⑤ 预审与 ⑥ 独立扫描的逐字原文未归档，其结论与轮次见报告「审查链」节；本片不宣称已归档该两环原文。

## 验收边界

⑤ 与 ⑥ 的 PASS 绑定整改前快照；最终快照仅由 ⑦ 各轮覆盖，且 ⑦ 仍为 FINDINGS。本包据此只证「部分闭合」，不证「安全闭合」。
