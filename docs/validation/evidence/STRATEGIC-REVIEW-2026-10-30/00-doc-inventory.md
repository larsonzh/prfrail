# P35 步骤 0 · 文档清单（00-doc-inventory）

切片 `STRATEGIC-REVIEW-2026-10-30`（P35）步骤 0 阶段 B 产出。锚点：`HEAD` = `ebec7d1`、准则 v1.56、`SELFTEST: PASS 697/697`。

**九字段口径**：路径 / 类别 / 状态 / 一句话摘要 / 行数-字节 / 入链 / 出链 / 最后修改 / 已知过期点（对照 HEAD 主动检查路径、节号、状态、commit 是否存在，无则写「无」）。

**机械口径**：行数与字节为脚本实测（含 BOM）；入链 = `docs/**/*.md`（不含 `evidence/`）中提到该文件名基名的文件数；出链 = 该文档内引用的不同 `.md` 文件名数；最后修改 = `git log -1 --date=short`。

| 路径 | 类别 | 状态 | 一句话摘要 | 行数 / 字节 | 入链 | 出链 | 最后修改 | 已知过期点 |
|---|---|---|---|---|---|---|---|---|
| `README.md` | 项目概览 | 单语 | ProofRail（证轨） | 172 / 11316 | 21 | 15 | 2026-10-25 | 无 |
| `docs/t027/REMAINING_SLICES.md` | 台账（切片与完成记录） | CN/EN 成对 | T027 剩余切片清单（临时执行台账） | 1377 / 240898（EN 1364 / 271436） | 106 | 88 | 2026-10-28（工作区未提交） | 记录 `HEAD = aa49e0c`；记录 `HEAD = 1819e2d`；记录 `HEAD = aa07eb4`；记录 `HEAD = 53b2a90`；记录 `HEAD = caf3b0f`；记录 `HEAD = e91a8b0` |
| `docs/validation/directive-p28-followup.md` | 验证报告（切片证据） | CN/EN 成对 | 证轨 P29 治理片验证报告：空子目录缝隙收口与三条协作纪律落点（准则 v1.51） | 202 / 20674（EN 202 / 24423） | 3 | 9 | 2026-10-26 | 无 |
| `docs/validation/directive-whitelist-policy.md` | 验证报告（切片证据） | CN/EN 成对 | DIRECTIVE-WHITELIST-POLICY 切片验证报告（协作纪律 D 升格为准则正文条款） | 222 / 22464（EN 222 / 26442） | 3 | 10 | 2026-10-26 | 无 |
| `docs/validation/directive-ledger-field-split.md` | 验证报告（切片证据） | CN/EN 成对 | DIRECTIVE-LEDGER-FIELD-SPLIT 切片验证报告（§9.3 计数序列钉死与台账字段统一） | 231 / 22128（EN 231 / 25301） | 1 | 11 | 2026-10-27 | 无 |
| `docs/validation/directive-d-mechanical-guard.md` | 验证报告（切片证据） | CN/EN 成对 | DIRECTIVE-D-MECHANICAL-GUARD 切片验证报告（准则 v1.54：纪律 D 义务结构化 / 记录化） | 228 / 28757（EN 228 / 33349） | 1 | 11 | 2026-10-28 | 无 |
| `docs/validation/directive-v21-obs-sweep.md` | 验证报告（切片证据） | CN/EN 成对 | 切片验证报告 — `DIRECTIVE-V21-OBS-SWEEP`（准则 v1.55：五条 v21 观察项整合） | 275 / 29137（EN 275 / 34670） | 1 | 8 | 2026-10-28 | 无 |
| `docs/validation/directive-tmp-self-heal.md` | 验证报告（切片证据） | CN/EN 成对 | 证轨 P34 治理片验证报告：临时件自愈特例（准则 v1.56） | 291 / 21067（EN 291 / 24913） | 1 | 7 | 2026-10-28（工作区未提交） | 无 |
| `docs/ADR_REGISTER.md` | ADR 登记表 | CN/EN 成对 | ProofRail 架构决策记录 | 98 / 38813（EN 99 / 45525） | 43 | 13 | 2026-10-27 | 无 |
| `docs/ARCHITECTURE.md` | 系统架构 | CN/EN 成对 | ProofRail 概要与架构设计 | 95 / 11871（EN 95 / 13245） | 6 | 7 | 2026-09-11 | 无 |
| `docs/BUSINESS_WORKFLOWS.md` | 业务流程 | CN/EN 成对 | ProofRail 业务流程 | 106 / 10035（EN 106 / 11529） | 6 | 3 | 2026-10-11 | 无 |
| `docs/CODING_CONVENTIONS.md` | 编码规范 | CN/EN 成对 | ProofRail（证轨）编码与工程规范（Coding Conventions） | 88 / 6523（EN 56 / 5629） | 11 | 5 | 2026-09-10 | 无 |
| `docs/CONTRACTS.md` | 契约（协议） | CN/EN 成对 | ProofRail 架构契约规范 | 326 / 73142（EN 318 / 85095） | 44 | 6 | 2026-10-22 | 无 |
| `docs/DELIVERY_DIRECTIVE.md` | 执行准则（治理权威） | CN/EN 成对 | ProofRail 交付执行准则（DELIVERY_DIRECTIVE） | 2158 / 440872（EN 2158 / 537740） | 114 | 64 | 2026-10-28 | 无 |
| `docs/DEV_PLAN.md` | 实施计划 | CN/EN 成对 | ProofRail 迭代开发计划 | 213 / 76012（EN 214 / 87120） | 35 | 99 | 2026-10-26 | 无 |
| `docs/DOCUMENTATION_PLAN.md` | 文档治理 | CN/EN 成对 | ProofRail 文档导航与维护计划 | 81 / 9443（EN 81 / 10243） | 7 | 33 | 2026-10-03 | 无 |
| `docs/INSTALLATION.md` | 安装指南 | CN/EN 成对 | ProofRail S1 Windows 便携 ZIP 安装指南 | 205 / 15157（EN 205 / 16748） | 9 | 5 | 2026-10-11 | 无 |
| `docs/INSTALLATION_PLAN.md` | 安装计划 | CN/EN 成对 | ProofRail 安装与部署规划 | 62 / 7399（EN 62 / 8848） | 5 | 7 | 2026-10-11 | 无 |
| `docs/OPERATIONS.md` | 运维 | CN/EN 成对 | ProofRail 用户、操作员与可信发布手册 | 141 / 16577（EN 141 / 17940） | 10 | 5 | 2026-10-11 | 无 |
| `docs/PRODUCT_REQUIREMENTS.md` | 产品需求（REQ） | CN/EN 成对 | ProofRail 产品需求与调研分析 | 129 / 16087（EN 127 / 17678） | 8 | 5 | 2026-09-10 | 无 |
| `docs/RFC-proofrail-unattended-ai-engineering-product.md` | RFC（项目初心记录） | 单语 | ProofRail（证轨）项目建议书——AI 无人值守工程通用产品抽取方案 | 2446 / 189575 | 22 | 3 | 2026-09-08 | 无 |
| `docs/S1_RELEASE_NOTES.md` | 发布说明 | CN/EN 成对 | ProofRail S1 候选发行说明 | 24 / 2120（EN 24 / 2481） | 6 | 3 | 2026-10-11 | 无 |
| `docs/S1_SUPPORT_MATRIX.md` | 支持矩阵 | CN/EN 成对 | ProofRail S1 候选支持矩阵 | 26 / 2557（EN 26 / 3032） | 7 | 2 | 2026-10-11 | 无 |
| `docs/SECURITY.md` | 安全 | CN/EN 成对 | ProofRail 安全与权限设计 | 64 / 8570（EN 64 / 9387） | 6 | 3 | 2026-09-10 | 无 |
| `docs/SLICE-CLOSEOUT-TEMPLATE.md` | 模板 | CN/EN 成对 | 任务包 / 阶段结项报告模板（SLICE-CLOSEOUT-TEMPLATE） | 43 / 1757（EN 43 / 2175） | 4 | 1 | 2026-10-03 | 无 |
| `docs/TEST_STRATEGY.md` | 测试与验收 | CN/EN 成对 | ProofRail 测试验证方案 | 73 / 10633（EN 73 / 11361） | 14 | 5 | 2026-10-14 | 无 |

## 健康度小结

1. 文档对形态健康：本清单所列入库文档均在库，CN 侧全部存在；缺 `_EN` 者已在「状态」列标为单语，供 5 轮评审按单语引用，不作镜像假设。
2. 过期点集中在历史 `HEAD` 记录：多份验证报告按 §11.1 保留成文时的 `HEAD`，不做追溯改写 —— 评审引用时须以本清单的锚点为准。
3. 台账与准则为唯一权威：切片状态取 `REMAINING_SLICES{,_EN}.md`，规则语义取 `DELIVERY_DIRECTIVE{,_EN}.md`；其余文档与二者冲突时以二者为准。
4. 本清单不覆盖 `docs/validation/evidence/**`（冻结证据树）与 `docs/checkpoints/**`、`docs/t027/**` 之外的子目录，评审判定以主清单为准。
