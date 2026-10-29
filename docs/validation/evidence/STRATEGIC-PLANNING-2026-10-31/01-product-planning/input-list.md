# 第 1 轮输入包清单（产品与架构规划）

| 路径 | 行数 | 字节 | 用途 |
|---|---|---|---|
| `docs/PRODUCT_REQUIREMENTS.md` | 130 | 16087 | S1/S2 范围与 REQ 全集 |
| `docs/ARCHITECTURE.md` | 96 | 11871 | 架构分层与漂移面 |
| `docs/BUSINESS_WORKFLOWS.md` | 107 | 10035 | 业务流程与偏离面 |
| `docs/CONTRACTS.md` | 327 | 73142 | 契约与发布耐久三态口径 |
| `docs/RFC-proofrail-unattended-ai-engineering-product.md` | 2446 | 189575 | 历史设计来源（非权威） |
| `docs/DEV_PLAN.md` | 213 | 76012 | 计划进度与状态表 |
| `docs/TEST_STRATEGY.md` | 74 | 10633 | 验收口径 |
| `docs/S1_SUPPORT_MATRIX.md` | 27 | 2557 | S1 支持矩阵 |
| `docs/validation/evidence/STRATEGIC-REVIEW-2026-10-30/01-product/review.md` | 84 | 18764 | P35 第 1 轮（产品层）评审原文 |

**合计**：9 件、3504 行、408676 字节（正文按本清单在派单时组装；证据包只留清单，不复写大文档，避免双源）。

**不包含**：准则全文（避免治理视角污染产品视角）。
**模型**：`GPT-6.1 Sol (copilot)`，档位 Max。
**口径**：行数按 `Get-Content` 记录数（末行无换行符时不虚减一行）；字节为原始字节数（含 BOM）。两项均由 `tmp/STRATEGIC-PLANNING-2026-10-31/step0-artifacts.js` 实测。
