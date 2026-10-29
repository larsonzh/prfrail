# 第 2 轮输入包清单（治理与流程规划）

| 路径 | 行数 | 字节 | 用途 |
|---|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 2160 | 442418 | 交付执行准则全文本体 |
| `docs/t027/REMAINING_SLICES.md` | 1399 | 243577 | 切片台账与观察项 |
| `docs/DEV_PLAN.md` | 213 | 76012 | 计划进度与状态表 |
| `docs/ADR_REGISTER.md` | 99 | 38813 | 架构决策登记表 |
| `docs/CONTRACTS.md` | 327 | 73142 | 契约与协议条款 |
| `docs/SECURITY.md` | 65 | 8570 | 安全前提与边界 |
| `docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/02-governance/review.md` | 68 | 17822 | P35 第 2 轮（治理层）评审原文 |
| `tmp/STRATEGIC-PLANNING-2026-10-31/00-checklist-v25.md` | 261 | 14432 | v25 清单全文（对话附件的落盘副本） |

**合计**：8 件、4592 行、914786 字节（正文按本清单在派单时组装；证据包只留清单，不复写大文档，避免双源）。

**不包含**：`internal/**` 源码与 schema（本片无实现）。
**模型**：`GPT-6.1 Sol (copilot)`，档位 Max。
**口径**：行数按 `Get-Content` 记录数（末行无换行符时不虚减一行）；字节为原始字节数（含 BOM）。两项均由 `tmp/STRATEGIC-PLANNING-2026-10-31/step0-artifacts.js` 实测。
