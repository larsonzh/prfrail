# T028 方案书 — supervised-black-box 显式模式（[PILOT]，零探针）

> 角色：① 架构师（只读交付，驱动器落盘）。验证依据：`docs/CONTRACTS.md` §7、`docs/DELIVERY_DIRECTIVE.md` §7.4、`docs/DEV_PLAN.md` T028 行、`docs/TEST_STRATEGY.md` AT-24 行，以及对 `internal/chain|console|adapters|guard|evidence`、`schemas/`、`testdata/contracts/` 的实读。

---

## 1. 目标与范围 / 边界复核

**目标**：在既有链/控制台/适配器面之上，**加法式**实现 supervised-black-box 显式模式，使 AT-24 全部子句闭合：启动前版本化风险提示 + 绑定确认、隔离 workspace、SessionBridge visible 投递、人工监督与显式归还、归还后全量扫描、独立 gates/review、reduced/unknown 报告、四条禁止项机械可验。

**核心判定（先说结论，详见 §9）**：**不需要真实外部通道**。本片全部验收可由 `StubVisibleClient` + 确定性夹具覆盖；真实 SessionBridge visible 接线属 T027-B/后续切片。AT-24 的“visible 投递”在本片被建模为一次**受监督的确定性投递事实**，其可用性结论走既有 `ai_availability` 的 `visible_delivery_unavailable` 失败码族，不发起任何真实 IPC。

**边界复核（全部遵守，逐条落点）**：
- 不触 §2.2/§2.4/§5.2/§6.3/§7.5：本片不改 `docs/DELIVERY_DIRECTIVE.md`（除 §12.4 台账登记候选 OB——那属切片例行的台账元数据行，与判据/白名单无关，由驱动器按 §6.3 惯例执行）。
- 不做 enforcement 生产接线（T027-B）：本片**不 spawn 任何进程**、不写 replay store、不 touch `agent-runner-enforcement`。
- 不改 SessionBridge 本体：只在 prfrail 侧新增 `VisibleClient` 接口与 stub 实现。
- 不改 AgentRunner 契约语义：`agent-runner-*` schema/记录/路由零改动，仅新增“互斥性”负例测试。
- 不动角色集/载体/`.github/agents/**`。
- 不重构 `internal/**` 既有核心链路：唯一两处核心文件触碰 = `internal/chain/postflight.go`（加法分支）与 `internal/chain/models.go`（+1 常量、+1 校验条件），**已在 §9 显式申报**；若驱动器/用户判定这已属“改变既有核心路径语义”，按切片边界**停下**，本方案不再推进。

**不在范围内（显式不触碰清单）**：`docs/TEST_STRATEGY.md`（AT-24 行已存在）、`docs/DELIVERY_DIRECTIVE.md` 正文、`internal/evidence/event.go`（转移表无需新转移）、`internal/adapters/ai_availability.go`（已含 visible 通道与失败码）、`internal/taskdef/preview.go`（枚举闭集，见 §9 歧义 A）、AgentRunner 全部记录/Schema/路由、`.github/agents/**`。

---

## 2. 现状与可复用面（AT-24 可直接建立的既有资产）

| 资产 | 位置 | 本片用法 |
|---|---|---|
| supervised-black-box 契约语义 | `CONTRACTS.md` §7（“supervised-black-box 是正式支持但保证降级…”整段 + “SessionBridge visible 已投递 | 仅是 UI 投递/人工观察事实…”表行） | 设计的规范性底座；本片只需把“风险确认记录、显式归还=handoff 完整回执、supervised postflight 事实、reduced/unknown 报告”补进该段 |
| 显式归还 = manual-handoff 完整路径 | `internal/chain/handoff.go`（`HandoffSession/HandoffDecision/BuildHandoffReceiptRecord/OpenHandoffTransition/ReturnHandoffTransition`）、`schemas/handoff-receipt.schema.json`（outcome=complete 要求非空 hookResultEvidence；前后 manifest/diff/lease 全绑定） | **零新记录复用**：显式归还就是 handoff `complete` 回执 |
| 操作员交互机器（确认/恢复） | `internal/chain/operator_interaction.go`（`OpenOperatorInteraction/ResumeOperatorInteraction/PrepareOperatorInteractionResume`）、`internal/console/interactions.go`（list/respond/tui、pending inbox、重启重建）、`schemas/operator-interaction.schema.json` | 风险确认 = 一次 `allowedResponses=[confirm,reject]` 的结构化交互；确认响应直接复用 `prfrail interactions respond` |
| 配置声明面 | `internal/console/config.go`（step `execution` 白名单 + `manual-handoff`→`chain.ManualHandoff` 映射、`validateAIConfig`/`validConfiguredAIChannel` 已含 `sessionbridge-visible`）、`schemas/chain.schema.json`（step oneOf: autonomous/manualHandoff） | 加第三种执行声明 `supervised-black-box`（沿用 handoffPolicy） |
| AgentRunner 装配范式 | `internal/console/runtime_agent.go`（`ExecuteAgentRunnerRun` 的 fail-closed 装配、`agentRunnerAssemblyStep` 单步限制、运行时端口实现、`nextEventID`）、`internal/adapters/agent_runner_*.go` | 新装配函数逐段镜像：哨兵错误驱动、端口闭包、零 spawn、幂等重入 |
| postflight 门禁（A5） | `internal/chain/postflight.go`（五事实位置化布局、`frozenFactsFromEvidence`、`requirePostflightQualifiedReview`、`runTaskPostflight`） | 新增 supervised 分支：同样的“链只执法可证明部分、端口重推导”纪律 |
| 通道与失败码 | `internal/adapters/ai_availability.go`（`validAIChannel` 含 `sessionbridge-visible`；`validAIAvailabilityReason` 含 `visible_delivery_unavailable`）、`internal/adapters/sessbridge.go`（`SilentClient`/`SessbridgeAdapter`/`mapDispatchOutcome`） | 照 `SilentClient` 形状新建 `VisibleClient`；**不产生任何回执记录**（投递事实只进确认记录，见 §3/§5） |
| 扫描/停机/门禁 | `internal/gates/`（`Runner`、`SecretScanner`、`ChainPort`）、`internal/guard/`（`ProcessAlive/LeaseStore/StopProcessIdentity`）、`internal/evidence/diagnostics.go`（只读脱敏诊断，T022）、`internal/snapshot`（manifest/diff 捕获） | 归还扫描 = 重捕 manifest/diff + 秘密扫描 + 范围/文件类型检查 + 受管进程停止证明 + `hooksAfterReturn` 门禁 hook |
| 契约夹具与 oracle | `testdata/contracts/{valid,invalid}/{core-config,runtime-records,...}.json`、`tools/contracts/`（Ajv + canonicalize；“全 schema 覆盖”要求新 schema 必有正反夹具） | 契约先行落点（§5） |

---

## 3. 设计概要

### 3.1 状态机（零新增状态、零新增转移）

全部复用 §2.1 既有转移表；模式状态刻画：

```
                        ① 风险确认              ② 工作窗口（人工监督）        ③ 显式归还          ④ 尾链
step:  PENDING→RUNNING ──port→ [送 visible 提示] → RUNNING→WAITING_FOR_OPERATOR →(complete) RUNNING→PASSED
                        ErrAwaitingSupervisedConfirmation      ▲                    │
task:  STEPS_RUNNING ──────────────────────→ WAITING_FOR_OPERATOR ──────────────── STEPS_RUNNING ─→ postflight→REVIEW_PENDING→PASSED
chain: RUNNING ──────────────→ PAUSED ────────────────(保持 PAUSED)────────────────────────────────→ RUNNING→COMPLETED
        confirmation interaction: OpenOperatorInteraction / ResumeOperatorInteraction（既有入口）
        工作窗口开/关: 新增 engine 入口 OpenSupervisedWorkWindow / SubmitSupervisedReturn（见 §4）
```

**入口/出口条件（谁确认什么）**：
- **① 启动**：`prfrail run` 检测定义含 supervised 步骤 → 路由到 supervised 装配（§4-C5）。装配物化隔离 workspace、捕获 before-manifest、经 `VisibleClient` 投递版本化提示，然后经既有 `OpenOperatorInteraction` 把 task/step 置 WAITING_FOR_OPERATOR、chain PAUSED。**未确认 = 停在 PAUSED，不可能越过**。
- **② 确认**：操作员 `prfrail interactions respond --selection confirm`。装配在下一轮驱动中：校验响应绑定 → 落盘 `SupervisedRiskConfirmationRecord`（新记录）→ `OpenSupervisedWorkWindow`（**engine 级检查确认记录**：绑定 run/task/step/attempt/workspaceHash、selection=confirm、noticeVersion/noticeHash 与当前版本常量一致，任一不符拒绝）→ 开 handoff 窗口（task/step WAITING_FOR_OPERATOR），持久化 store-local 会话文件。
- **③ 显式归还**：`prfrail supervised return --outcome complete`：重捕 manifest/diff、范围/秘密/文件类型/副作用扫描、受管进程停止证明、`hooksAfterReturn` 门禁 hook，全部通过后写 `HandoffReceiptRecord(complete)` + `SupervisedReport`（reduced/unknown），再 `SubmitSupervisedReturn`（engine 校验回执 outcome=complete + 五事实摘要互异）→ 步骤 PASSED（证据 = 五事实位置化布局，reason=`supervised-return-completed`）。
- **④ 尾链**：`engine.Run` → `runTaskPostflight` **supervised 分支**（从步骤 PASSED 证据重建五事实，SupervisedPostflight 端口重推导并逐字段对账，任一缺失即暂停/拒绝）→ REVIEW_PENDING → Accept → Review（既有 policy review，**无任何快捷通道**）→ Publish → task PASSED → chain COMPLETED。

**“投递 ≠ 完成”如何被强制**（结构性，非承诺）：
1. 投递事实不产生任何状态转移、不回执记录（只进确认记录的 `delivery` 字段）；链上不存在“投递→PASSED”的转移路径；
2. 步骤 PASS 的唯一 supervised 入口是 `SubmitSupervisedReturn`，它要求完整 handoff 回执 + 五事实；
3. REVIEW_PENDING 的唯一资格是 supervised postflight 分支通过（镜像 A5 的 `requirePostflightQualifiedReview`）；
4. 任何 exit code/聊天文本/确认响应都不出现在④的任一输入里。

---

## 4. 变更点清单（精确到文件）

### A. 契约文档（先行，顺序见 §5）

| # | 文件 | 变更 |
|---|---|---|
| A1 | `docs/CONTRACTS.md`、`docs/CONTRACTS_EN.md`（同位镜像） | §7 supervised-black-box 段：补 (a) 绑定风险确认记录定义与启动前阻断语义；(b) 显式归还 = handoff `complete` 回执（引用 §6）；(c) supervised 五事实与 postflight 分支；(d) reduced/unknown 报告记录与“禁止把未知项显示为安全”；§2 记录表加“supervised records”行；§7 行为表加两行（“supervised 风险确认 confirmed | 仅解锁工作窗口，不产生任何通过结论”、“supervised report | reduced/unknown 事实记录，assuranceLevel 恒 reduced”） |

### B. Schema（先行）

| # | 文件 | 变更 |
|---|---|---|
| B1 | `schemas/chain.schema.json`（新） | `codeStep` 的 oneOf 增第三变体 `supervisedBlackBoxCodeStep`：`execution` const `"supervised-black-box"`，`required: ["id","kind","execution","handoffPolicy"]`，其余字段同 `manualHandoffCodeStep` |
| B2 | `schemas/supervised-risk-confirmation.schema.json`（新） | 字段：`schemaVersion`、`record{kind:"supervised-risk-confirmation", recordId, issuedAt, issuedBy(operatorActor), runId, taskId, stepId, attempt, workspaceHash, riskNoticeVersion(integer≥1), riskNoticeHash(sha256), selection const "confirm", delivery{channel const "sessionbridge-visible", deliveredAt, deliveryStatus enum [accepted,unavailable,unknown], transportRequestId, errorEvidence(hashSet)}, interactionEvidence(长度2 的 hashSet：请求+响应记录哈希), evidence(非空 hashSet)}`、`recordHash`；域 `proofrail:supervised-risk-confirmation:1\n` |
| B3 | `schemas/supervised-report.schema.json`（新） | `record{kind:"supervised-report", reportId, recordedAt, recordedBy(systemActor), runId, taskId, stepId, attempt, mode const "supervised-black-box", assuranceLevel const "reduced", verifiedDimensions(array, 闭集枚举), unknownDimensions(array, minItems 1, 闭集枚举且必须含 cost), evidence(非空 hashSet)}`、`recordHash`；域 `proofrail:supervised-report:1\n`；闭集 = `tool-use, resource-reads, network-access, cost, background-processes, session-continuity, external-side-effects`；verified 闭集 = `workspace-isolation, manifest-before, manifest-after, diff, scope-scan, secret-scan, file-type-scan, side-effect-scan, process-stop-proof, independent-gates, independent-review` |

### C. Go 实现

| # | 文件 | 变更（符号级） |
|---|---|---|
| C1 | `internal/chain/models.go`（编辑） | +`SupervisedBlackBox ChangeMode = "supervised-black-box"`；`StepExecutionIntent.validate` 增加条件：`step.Mode == SupervisedBlackBox` 时 `intent.ExecutionTarget` 必须为 `DefaultExecution` 且 `AgentRunnerFacts == nil`（否则 ErrInvalidDefinition）——禁止任何路径把 supervised 步骤送往 AgentRunner |
| C2 | `internal/chain/router.go`（编辑） | +接口 `SupervisedBlackBoxPort{ OpenSupervisedBlackBox(context.Context, StepRequest) (StepResult, error) }`；`StepRouter` +字段 `Supervised SupervisedBlackBoxPort`；`Execute` 的 `case SupervisedBlackBox:` 分支（port 为 nil → ErrInvalidDefinition） |
| C3 | `internal/chain/supervised_blackbox.go`（新） | 哨兵 `ErrAwaitingSupervisedConfirmation`、`ErrAwaitingSupervisedReturn`；reason 常量 `supervisedReturnCompletedReason = "supervised-return-completed"`；`SupervisedFrozenFacts{HandoffReceiptHash, ScanEvidenceHash, ProcessStopEvidenceHash, GateEvidenceHash, ReportHash}` + `Validate()`（5 摘要全有效且互异）+ `DisjointFrom(...)`；位置化重建 `supervisedFactsFromEvidence(input []string)`；engine 方法 `(e *Engine) OpenSupervisedWorkWindow(ctx, session HandoffSession, confirmation SupervisedRiskConfirmationRecord) error`（前置：task=STEPS_RUNNING、step=RUNNING、chain=PAUSED；确认记录校验：`Validate` + 绑定四元组 + `Selection=="confirm"` + notice 版本/哈希==当前常量；转移=复用 `OpenHandoffTransition`；证据=[confirmationHash]）与 `(e *Engine) SubmitSupervisedReturn(ctx, receipt HandoffReceiptRecord, facts SupervisedFrozenFacts) error`（前置：task/step=WAITING_FOR_OPERATOR；校验 receipt.Outcome=="complete" + facts.Validate + facts 与 receipt 哈希 Disjoint；转移=`ReturnHandoffTransition(complete)` 后 step→PASSED，InputEvidence=facts.hashes()，reason=`supervised-return-completed`；步骤已 PASSED 时按证据相等幂等收敛、不等则冲突拒绝）；风险提示常量 `SupervisedRiskNoticeVersion=1` + 版本化提示文本常量（七项不可证明维度枚举，与 CONTRACTS §7 措辞逐字对齐）+ `SupervisedRiskNoticeHash()`（域 `proofrail:supervised-risk-notice:1\n`） |
| C4 | `internal/chain/supervised_risk_confirmation.go`（新） | `SupervisedRiskConfirmation` 结构 + `NewSupervisedRiskConfirmation(...)`/`Decode`/`Validate`（严格解码、recordHash 复算、绑定/枚举/时间戳校验）；`SupervisedRiskConfirmationLedger`（追加式 JSONL、no-replace、Append 去重：同 (run,task,step,attempt) 同摘要幂等、异摘要冲突；`Confirmed(binding)` 查询） |
| C5 | `internal/chain/supervised_report.go`（新） | `SupervisedReport` 结构 + `NewSupervisedReport(...)`（恒定 `AssuranceLevel="reduced"`，`unknownDimensions` 恒含 `cost` 且全部出自闭集）+ `Decode`/`Validate`；域 `proofrail:supervised-report:1\n` |
| C6 | `internal/chain/engine.go`（编辑，仅 Options） | `Options` +字段 `SupervisedPostflight SupervisedPostflightPort`（无该端口且命中 supervised 分支时 fail-closed `ErrPostflightUnavailable`） |
| C7 | `internal/chain/postflight.go`（编辑，加法分支，**已申报**） | +接口 `SupervisedPostflightPort{ RunSupervisedPostflight(context.Context, SupervisedPostflightRequest) (SupervisedPostflightDecision, error) }` + `SupervisedPostflightRequest{RunID,TaskID,Attempt,ParentSnapshotHash,Workspace,Facts SupervisedFrozenFacts}` + `SupervisedPostflightDecision.Validate(facts)`（passed 必须携带全部五事实摘要）；+`supervisedTaskFacts(events,runID,taskID)`（扫描 reason=`supervised-return-completed` 的 step 事件；多步骤/多回执拒绝）；`runTaskPostflight` 插入分支：agent 事实与 supervised 事实同现 → 拒绝（混合执行任务）；supervised 事实在 → 走 SupervisedPostflight 端口并以其证据写 REVIEW_PENDING；`requirePostflightQualifiedReview` 同步加 supervised 分支（REVIEW_PENDING 证据必须含五事实）。既有 agent 分支与 `!found→"steps-completed"` 路径零改动 |
| C8 | `internal/console/config.go`（编辑） | `ValidateChainConfig`：`execution` 白名单 +`"supervised-black-box"`；supervised 步骤必须带 `handoffPolicy`（复用 `validateHandoffPolicy`）；**新增** `validateSupervisedChannel(cfg)`：定义含 supervised 步骤时 `cfg.AI != nil` 且 `channels[]` 含 `sessionbridge-visible`（否则 ErrInvalidConfig，fail-closed）；`ToDefinition`：`execution=="supervised-black-box"` → `Mode=chain.SupervisedBlackBox` + HandoffPolicy 映射；+`FindSupervisedSteps(cfg) []ExecutableStep` |
| C9 | `internal/console/cli.go`（编辑） | `executeRun`：在 noop-only 拒绝**之前**，`if len(FindSupervisedSteps(cfg))>0 { return cli.executeSupervisedRun(...) }`；supervised 定义 + 任何 agent-runner 输入 flag 同时给出 → 拒绝（互斥，机械可验 P1） |
| C10 | `internal/console/runtime_supervised.go`（新，本片主体） | `SupervisedAssemblyOptions{RunID, RunRoot, Definition, VisibleSender adapters.VisibleClient, GateRunner *gates.Runner, Clock, IDs}`；`ExecuteSupervisedRun(ctx, options) (RunSummary, error)`：磁盘状态机驱动（`<runRoot>/supervised/` 下 `risk-confirmations.jsonl`、`session.json`、`handoff-receipt.json`、`report.json`、`scans/`）：未确认→开交互后返回等待态；已确认未开窗→写确认记录+开窗；已开窗→等待归还；`supervisedRuntime` 实现 chain 各端口（Baseline/Workspace/Steps/Acceptance/Reviewer/Publisher/Stopper/Reconciler，镜像 `agentRunnerRuntime`；Steps 内 `StepRouter{Supervised: 端口闭包}`）；端口闭包 `OpenSupervisedBlackBox`：捕获 before-manifest → `VisibleSender` 投递 → 返回哨兵 `ErrAwaitingSupervisedConfirmation`（**不写任何转移**）；归还扫描实现：`snapshot` 重捕 manifest/diff、`gates.SecretScanner`、范围/文件类型检查、`guard` 受管进程停止证明（无受管进程=平凡证明，镜像 `managedProcessesGone`）、`hooksAfterReturn` 门禁 hook（`gates.Runner`）；`SupervisedPostflightPort` 实现（重推导五事实并逐字段对账）；`ExecuteSupervisedReturn(...)`（写回执+报告 → `SubmitSupervisedReturn` → `engine.Run` 尾链）；超时判定（`closedAt-openedAt > HandoffTimeoutMs` → 拒绝 complete，仅允许 abort） |
| C11 | `internal/console/supervised_cmd.go`（新） | `executeSupervisedRun` 装配 + 命令组 `prfrail supervised <return|report>`：`return --run-id ... --outcome complete|abort`（显式归还）；`report`（打印持久化报告，golden 渲染：未知维度一律输出 `unknown`） |
| C12 | `internal/console/preview.go`（编辑） | `buildPreviewSteps`：`execution=="supervised-black-box"` → 预览 kind 映射为 `"manual-handoff"`（预览层近似，代码注释声明；见 §9 歧义 A） |
| C13 | `internal/adapters/sessbridge_visible.go`（新） | `VisibleClient{ SendVisible(context.Context, RequestEnvelope) (VisibleResponse, error) }`；`VisibleResponse{Status, TransportRequestID, OutputEvidence, ErrorEvidence}`；`DeliverVisible(...) → {DeliveryStatus: accepted|rejected|uncertain, ErrorEvidence}`（映射复用 `mapDispatchOutcome` 口径，失败码 `visible_delivery_unavailable`）；`StubVisibleClient`（确定性脚本桩：按 scripted status 返回，零网络、零文件写入）。**不产生任何回执记录**（投递事实只进入确认记录的 `delivery` 字段） |

### D. 契约夹具与测试文件（详见 §5/§6）

### E. 交付文档

| # | 文件 | 变更 |
|---|---|---|
| E1 | `docs/validation/t028-supervised-black-box.md` + `_EN.md`（新，严格同位镜像，BOM+LF） | 验证报告对（审查链 ⑤⑥⑦ 结论、四条禁止项验证、门禁读数） |
| E2 | `docs/validation/evidence/T028-freedom-list.md`（新） | 自由度清单（⑥ 后/⑦ 前落盘并纳入 ⑦ 输入包，OB-48 口径）；建议 F 项：记录归属 chain 包、预览层 kind 近似、超时仅限归还时判、报告按 supervised 步骤单枚、stub 交付器随 adapter 入库、`<runRoot>/supervised/` store-local 协议为非 wire 文件 |
| E3 | `docs/t027/REMAINING_SLICES{,_EN}.md`（编辑） | T028 行勾选 + 状态更新 + 完成记录；**顺带修正** T027-ASSEMBLY-MIN 完成记录行的过期「未提交」字样（用户 2026-10-03 裁定并入下一片起跑回写） |
| E4 | `docs/DEV_PLAN.md` + `_EN.md`（编辑） | T028 行状态 `PLANNED / NOT IMPLEMENTED` → 已完成 + 证据引用（REQ 表不动） |

---

## 5. 契约先行判定

**判定：需要 2 个新记录类型；显式归还复用既有记录。**

1. **新记录 1 — `supervised-risk-confirmation`（绑定风险确认）**。理由：`operator-interaction`（T025 冻结）无“风险版本/通知哈希/投递事实”字段，`authorization-record` 是授权 grant 语义，均不能表达“绑定 task/attempt/workspace/风险版本的用户确认且为启动前置”；复用 convention 会把版本绑定塞进 `question` 文本，弱化机械可验性（P3/P4 需要记录级字段）。
2. **显式归还 — 复用 `handoff-receipt`（outcome=complete）**，**不新增** explicit-return receipt。理由：CONTRACTS §6 handoff 段已冻结“归还前后 manifest/diff、越界/秘密/未知进程扫描、规定 gates”语义，与 AT-24 归还子句逐字对应；T014 的 `BuildHandoffReceiptRecord`/`ReturnHandoffTransition` 现成。
3. **新记录 2 — `supervised-report`（reduced/unknown 报告）**。理由：`verification-report` 是 T002 冻结的 16 类 checker 输出，不可扩展；报告必须作为一等记录承载 `assuranceLevel const "reduced"` + `unknownDimensions`（含 cost），才能被 schema/Go 校验机械拒绝“显示为安全”。
4. **不新增 visible-delivery receipt**：投递事实内嵌于确认记录的 `delivery` 字段；结构性杜绝“回执冒充完成”（P2 的第一道防线就是没有这个对象）。

**顺序（硬性）**：A1（CONTRACTS 双语）→ B1/B2/B3（schema）→ D 契约夹具 → C（Go 代码）→ 测试 → E（文档）。任何一步的 schema/夹具与 CONTRACTS 措辞不一致，后序步骤不得开始。

**夹具变更（精确）**：
- `testdata/contracts/valid/core-config.json` + 正例：含 `execution:"supervised-black-box"` + 完整 handoffPolicy 的链实例；
- `testdata/contracts/invalid/core-config.json` + 反例：supervised 步骤缺 handoffPolicy；supervised 步骤带未知字段（如 `waiveGates`）——additionalProperties 拒绝（P3 的 schema 层机械证）；
- `testdata/contracts/valid/runtime-records.json` + 正例：合法 confirmation（selection=confirm、delivery=accepted）与合法 report（reduced、unknownDimensions 含 cost）；
- `testdata/contracts/invalid/runtime-records.json` + 反例：selection=reject；recordHash 失配；notice 版本/hash 缺；report assuranceLevel=full；unknownDimensions 为空；unknownDimensions 缺 cost；未知维度 id；deliveryStatus 非法值。
- 校验：`tools/contracts` 下 `npm test`（Ajv 全 schema 覆盖检查会强制新 schema 必须有正反夹具）；canonical vectors（`testdata/contracts/vectors/canonical.json`）本片**不增**（PILOT 可选，注明即可）。

---

## 6. 测试计划（正例 + 反例，可执行）

统一前置：`go build ./...`、`go vet ./...`；夹具用固定小文件、`t.TempDir()`；新测试全部做变异检验（删守卫→变红）。

**① `go test -count=1 ./internal/chain/...`**
- 正例：
  - `TestSupervisedRiskConfirmationRoundTrip`：合法记录构造/解码/复算摘要；
  - `TestSupervisedRiskConfirmationLedgerIdempotentAppend`：同键同摘要幂等；
  - `TestSupervisedReportNewAlwaysReduced`：构造器恒 `reduced`、unknownDimensions 恒含 cost；
  - `TestOpenSupervisedWorkWindowWithValidConfirmation`：合法确认 → task/step=WAITING_FOR_OPERATOR、chain=PAUSED、转移证据含确认哈希；
  - `TestSubmitSupervisedReturnCompleteFlow`：complete 回执 + 五事实 → step PASSED、证据位置化布局可重建；
  - `TestSupervisedPostflightQualifiesReviewPending`：supervised 事实在 → 端口 passed（证据含五事实）→ REVIEW_PENDING 写入；
  - `TestSubmitSupervisedReturnIdempotentReplay`：步骤已 PASSED 且证据相同 → 幂等收敛。
- 反例：
  - `TestSupervisedRiskConfirmationRejectsBindingDrift`：run/task/step/attempt/workspaceHash 任一漂移 → 拒绝；
  - `TestSupervisedRiskConfirmationRejectsRejectSelection`：selection=reject 不可作为确认；
  - `TestOpenSupervisedWorkWindowWithoutConfirmation`：无确认记录 → 拒绝且零写入；
  - `TestOpenSupervisedWorkWindowWrongNoticeVersion`：notice 版本/哈希不符 → 拒绝（版本漂移即阻断）；
  - `TestSubmitSupervisedReturnRejectsAbortReceipt`：outcome=abort 回执 → 拒绝；
  - `TestSubmitSupervisedReturnRejectsFactCollisions`：五事实存在重复摘要 → `ErrInvalidPostflightFacts`；
  - `TestSubmitSupervisedReturnRejectsDeliveryOnlyEvidence`：以投递/确认类哈希充当五事实 → 拒绝（P2 机械证）；
  - `TestSupervisedPostflightDecisionMustConsumeAllFacts`：passed 缺任一事实 → `ErrInvalidPostflightDecision`；
  - `TestSupervisedTaskWithoutReturnFactsFailsClosed`：supervised 步骤无归还事实 → 任务 FAILED/REPAIR，REVIEW_PENDING 不写（P3 机械证）；
  - `TestSupervisedUnqualifiedReviewRefused`：REVIEW_PENDING 证据缺五事实 → `ErrUnqualifiedReview` 零写入；
  - `TestSupervisedStepRejectsAgentRunnerExecutionTarget`：`intent.validate` 拒绝 supervised+AgentRunnerExecution（P1 机械证）；
  - `TestRouterRefusesSupervisedWithoutPort`：port 为 nil → ErrInvalidDefinition。

**② `go test -count=1 ./internal/console/...`**
- 正例（stub 驱动全链，镜像 `TestAgentRunnerAssemblyFullChain`）：
  - `TestSupervisedAssemblyFullChain`：stub 投递 accepted → 开交互 → confirm → 开窗 → 写入工作区产物 → `supervised return`（扫描通过、hook 通过）→ postflight → policy review → promotion → run COMPLETED；报告 `reduced` 且 cost/process/side-effect 维度为 unknown；
  - `TestSupervisedRunRestartRecovery`：窗口期重启 → 状态重建 → return 正常闭环（T025 重启口径）；
  - `TestSupervisedReportRenderGolden`：`prfrail supervised report` 文本 golden：未知维度行含 `unknown`、不含 `safe/verified`（P4 机械证）。
- 反例：
  - `TestSupervisedAssemblyBlocksWithoutConfirmation`：不 confirm → 停在等待态，无任何 PASS（P2 机械证）；
  - `TestSupervisedAssemblyRejectsSelectionReject`：响应 reject → 阻断；
  - `TestSupervisedVisibleUnavailableFailsClosed`：stub 状态 unavailable/unknown（`visible_delivery_unavailable`）→ 不开窗、不启动；
  - `TestSupervisedReturnRefusesFailedGate`：确认存在 + 归还扫描通过但 `hooksAfterReturn` 门禁 hook 失败 → 拒绝（P3 机械证：确认不豁免门禁）；
  - `TestSupervisedReturnRefusesTimeout`：超窗 complete → 结构化拒绝，仅 abort 可；
  - `TestSupervisedReturnRefusesSecretFinding`：秘密扫描命中 → 拒绝；
  - `TestRunCommandRefusesMixedSupervisedAgentRunnerInputs`：supervised 定义 + agent-runner 输入 → 拒绝（P1 机械证）；
  - `TestAgentRunnerAssemblyRefusesSupervisedDefinition`：`ExecuteAgentRunnerRun` 收 supervised 定义 → ErrInvalidConfig、零文件（P1 机械证）；
  - `TestSupervisedAssemblyRefusesAgentRunnerDefinition`：`ExecuteSupervisedRun` 收 agent-runner 定义 → ErrInvalidConfig；
  - `TestNoSupervisedReferencesInAgentRunnerPaths`：源码级不变量（OB-61 先例）：`internal/adapters/*.go` 与 `runtime_agent.go` 源文本不含 `supervised-black-box`/`SupervisedBlackBox` token；`sessbridge_visible.go` 不引用任何 chain 转移（P1/P2 机械证，附“等价改写可规避”的限制声明）。

**③ `go test -count=1 ./internal/adapters/...`**
- 正例：`TestVisibleDeliveryMapsAccepted`（accepted → DeliveryStatus accepted）；`TestStubVisibleClientDeterministic`。
- 反例：`TestVisibleDeliveryUnavailableReason`（transport 错误 → uncertain + `visible_delivery_unavailable`）；`TestVisibleDeliveryRequestIDMismatch`（requestId 失配 → uncertain）。

**④ 全量与门禁（验收口径）**
`go test -count=1 ./...`；`tools/contracts` 下 `npm test`；`node tools/gates/gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0`；`--scope=index` ⇒ `TOTAL_FAIL=0`；`--selftest` 全绿。Linux 侧（GitHub Actions，含 `-race` 与 Contract fixtures 步骤）由 ⑨ 原生验证执行；本机无 gcc 不跑 -race（既有口径）。

---

## 7. 四条禁止项的机械可验映射

| 禁止项 | 机械可验检查（不依赖散文承诺） |
|---|---|
| **P1 不自动从 AgentRunner 降级** | ① `TestSupervisedStepRejectsAgentRunnerExecutionTarget`（intent.validate 拒绝，链级）；② `TestAgentRunnerAssemblyRefusesSupervisedDefinition`（装配入口拒绝）；③ `TestRunCommandRefusesMixedSupervisedAgentRunnerInputs`（CLI 路由拒绝）；④ `TestNoSupervisedReferencesInAgentRunnerPaths`（源级不变量：agent-runner 侧源文本零 supervised token；supervised 装配零 AgentRunner token）。四者任一删守卫即红 |
| **P2 visible 回执不冒充完成** | ① 结构层：投递不产生回执记录、无“投递→PASS”转移路径（§3 设计）；② `TestSubmitSupervisedReturnRejectsDeliveryOnlyEvidence`（投递哈希不能充当事实）；③ `TestSupervisedAssemblyBlocksWithoutConfirmation`（投递 accepted + 未归还 → 永无 PASSED）；④ `sessbridge_visible.go` 源级不变量（不引用链转移、不写回执）。验收态：e2e 断言事件日志中 delivery 阶段零状态事件 |
| **P3 风险确认不豁免安全门禁** | ① `TestSupervisedTaskWithoutReturnFactsFailsClosed` + `TestSupervisedUnqualifiedReviewRefused`（确认存在但无归还/postflight 事实 → REVIEW_PENDING 不可达）；② `TestSupervisedReturnRefusesFailedGate`（确认存在 + 门禁 hook 失败 → 任务 FAILED，确认不豁免）；③ schema 反例夹具（`waiveGates` 未知字段被 additionalProperties 拒绝；confirm 记录哈希不被 review/promotion 接受，域不相交）；④ postflight 分支的 five-facts 对账（passed 缺事实即拒） |
| **P4 未知费用/进程/副作用不显示为安全** | ① schema const：`assuranceLevel` 恒 `reduced`、`unknownDimensions` minItems≥1 且必含 `cost`（反例夹具：full/空/缺 cost/未知 id 均被拒）；② Go `NewSupervisedReport` 恒 reduced + `ValidateSupervisedReport` 拒绝变体；③ `TestSupervisedReportRenderGolden`（渲染输出对未知维度打印 `unknown`，grep 断言无 `safe/verified` 字样）；④ 版本化风险提示文本恒含七项不可证明维度，其哈希进入确认记录绑定（改提示=版本漂移=阻断，`TestOpenSupervisedWorkWindowWrongNoticeVersion`） |

---

## 8. reduced/unknown 报告口径

- **`reduced` = 模式天花板（恒定）**：supervised-black-box 是“保证降级”的候选生产/导入模式（CONTRACTS §7 原句）。只要 ProofRail 自身义务（隔离物化、归还后重捕 manifest/diff、范围/秘密/文件/副作用扫描、进程停止证明、独立 gates/review、接受事务）全部执行完毕，报告即 `assuranceLevel=reduced`——**永远不写 full**；schema const + 构造器 + 校验三层保证。它不是“降级状态”，而是该模式的**最高可达保证**。
- **`unknown` = 按维度、按门的判定**：外部代理的工具使用、资源读取、网络、token/费用/重试、后台进程、会话连续性、commit/push/发布等副作用**在原理上不可观测**，逐项列入 `unknownDimensions`（恒含 `cost`）。渲染与报告输出对每个未知维度只显示 `unknown`，禁止显示 `safe`/`verified`/数字承诺。
- **阻断与标注分界**：可检测的未知（进程停止证明失败、扫描不确定、断线超窗、投递未知）→ **阻断**（拒绝 complete / 暂停对账）；原理上不可观测的未知（费用、外部副作用）→ **如实标注 + 报告留痕**，不冒充安全、也不臆造“已清除”。判定规则：凡 ProofRail 有义务且有能力证明的维度未证明 ⇒ 阻断；无能力观测的维度 ⇒ unknown 标注。
- **“禁止自动降级”的落点**：指模式选择（P1：不得因 AgentRunner 不可用而自动落入本模式），与报告等级无关——报告恒 reduced 本身就是“不降级承诺”的表达。

---

## 9. 已知风险 / 歧义点 / 需升级项

**首项（切片定义 §1.7.2 要求先行回答）**：**真实外部通道 = 不需要，本方案不停下。** 判定依据：AT-24 全部行为可经 `StubVisibleClient`（确定性脚本桩）+ 文件夹具驱动闭环；可见通道在本片只是一次“受监督的投递事实”输入。真实 SessionBridge visible IPC 接线属 T027-B/后续生产接线，不在本片，亦不在验收证据链中。零探针成立。

其余风险/歧义：
1. **`internal/chain/postflight.go` + `models.go` 属核心文件触碰**（尽管是仅对新模式生效的加法分支）。按切片边界字面，若判定“加法分支也算改变既有核心路径语义”，须停下由用户裁决。我的判断：不触发（既有 agent 分支与无事实路径零改动、新分支仅对新增模式值生效），但**显式申报**。
2. **歧义 A（preview 层 kind 枚举）**：`internal/taskdef/preview.go:174` 的 kind 闭集不含 supervised；本方案选 C12 的映射近似（不改 taskdef）。若用户要预览层如实显示新 kind，需扩 taskdef 枚举（超出片列三包），须同轮授权。
3. **歧义 B（报告作用域）**：报告按“每个 supervised 步骤一枚”设计（PILOT 装配恰一枚）。多 supervised 任务/多枚报告是否需要 run 级聚合报告，留待后续切片。
4. **歧义 C（断线/超时语义）**：T014 冻结了 `HandoffTimeoutMs` 但无任何执行接线；本片仅在归还时点判超窗拒绝 complete（仅允许 abort），不做运行期定时器（无受管进程可停，也无“自动转租”）。这是 PILOT 简化，已列入自由度清单。
5. **歧义 D（“未知进程”扫描边界）**：guard 只能证明**受管**进程停止；外部代理进程本质上不可观测 ⇒ 归入报告的 `background-processes` 未知维度，而非硬扫描。若用户要求 OS 级进程发现扫描，需升级（涉及新 guard 能力）。
6. **歧义 E（确认记录的交付失败语义）**：deliveryStatus=unavailable/unknown 时确认记录**不产生**（阻断在开窗之前）；是否要为“已投递但未确认”单独留痕，本片用 operator-interaction 请求记录本身承载（已持久化），不另立记录。
7. 残余：确认交互与 handoff 会话的**幂等重放边界**（`session.json` 先写后转移的恢复矩阵，§3/§4）需在 ④ 阶段做崩溃注入测试；已知边界：非原子窗口（写文件与转移之间）由“先持久后转移”的顺序收敛，PILOT 不做断电级证明。

---

## 10. 开放问题（需用户/主控裁决）

1. **两个新记录类型**（risk-confirmation + report）是否接受？不接受则备选：确认复用 operator-interaction convention（版本绑定弱化）、报告退化为 console 输出（P4 机械可验性降级，需另行裁定可接受度）。
2. **歧义 A**：preview 层 kind 映射近似（C12，不改 taskdef）是否接受？
3. **CLI 面**：`prfrail supervised return|report` 命令组命名是否接受（备选：并入 `handoff` 命令族）？
4. **超时/断线 PILOT 语义**（歧义 C）是否接受？
5. 自由度清单落盘点（E2）与 ⑦ 输入包材料化（head-baseline/staged-diff/gate 读数）由驱动器执行——按 OB-48 口径，⑥ 后、⑦ 前，请主控在派 ⑦ 前核验其存在。

（完）