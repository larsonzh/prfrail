# T027 · A5 · chain postflight acceptance — 验证报告

日期：2026-09-15。结论：`A5 完成`；已提交 `b8bad47` 并推送，GitHub Actions Ubuntu（含 `Race` 与 `Contract fixtures` 步骤）/Windows 全绿（run 34950623175）。
切片目标：把 task 的 review 资格从 adapter 记录收敛到 chain 自有 postflight 门——只有“可重建的冻结事实 + 端口独立判定通过”才允许写 `REVIEW_PENDING`。

## 1. 交付范围

- 新增 `internal/chain/postflight.go`：chain 自有冻结事实 DTO `AgentRunnerFrozenFacts`（manifest/diff/log/usage/process-stop 五枚摘要，要求两两互异且与请求/完成摘要互异）、`PostflightOutcome`（passed/failed/uncertain）、`PostflightRequest`、`PostflightDecision.Validate(facts)`、必填端口 `PostflightPort`、证据布局常量（`postflightFactsPrefix=2`、`postflightFactsCount=5`、`postflightFactsMinEvidence=7`）、`agentRunnerFactsEvidence`、`frozenFactsFromEvidence`、`agentRunnerTaskFacts`、`runTaskPostflight`、`requirePostflightQualifiedReview`，以及哨兵 `ErrInvalidPostflightFacts` / `ErrInvalidPostflightDecision` / `ErrPostflightUnavailable` / `ErrUnqualifiedReview`。
- 新增 `internal/chain/postflight_test.go`（事实校验矩阵、证据重建、决策校验、门通过/拒绝/uncertain/端口失败/不完整决策、部分写收敛、`REVIEW_PENDING` 已持久化的重入收敛、多步终局拒绝、损坏证据重放冲突、review 资格三态与重入拒绝、逐任务归属）。
- 改动 `internal/chain/engine.go`：`Options.Postflight` 必填（nil → `ErrPostflightUnavailable`）；`runTask` 在全部 step 通过后、`Acceptance.Accept` 之前运行门，且仅当 task 尚未处于 `REVIEW_PENDING`；重入 `REVIEW_PENDING` 时改为要求资格证明。
- 改动 `internal/chain/agent_runner_terminal.go` / `agent_runner_terminal_route.go`：终局 DTO 增 `Facts`（completed 必带、非 completed 必须为 nil），路由证据扩展为 `[RequestHash, CompletionHash, 五枚事实摘要, ...Evidence]`（去重保序），completed 重放按事实比对。
- 改动 adapters：`agent_runner_terminal_record.go`（intent 增 `frozenFacts` 块与校验）、`agent_runner_terminal_publisher.go`（completed 必须在写 intent 前携带事实）、`agent_runner_terminal_chain.go`（单向投影含 `toChain()` 事实重建）、`agent_runner_replay_store_terminals.go`（结算计划比较注释说明为何排除事实）及对应测试。
- 改动 `internal/console/runtime.go`：`chain.Options.Postflight = runtime`，`localRuntime.RunPostflight` 为 fail-closed 桩（返回 `ErrPostflightUnavailable`，noop-only 运行不会触发）。
- 改动 `internal/evidence/event.go`：task 转移表新增 `STEPS_RUNNING→REPAIR_PENDING`。

## 2. 协议先行产物

- `docs/CONTRACTS.md` / `_EN.md` §2.1：task 行新增 `STEPS_RUNNING→REPAIR_PENDING`，并注明该转移仅由 postflight 拒绝写入。
- 同章节新增 §7“postflight acceptance（A5）”段：冻结事实 DTO 与“completed 必带/非 completed 必须为空”、五枚互异且与请求/完成互异、证据前缀按位置重建、端口必须独立重推（重捕 manifest 并与父快照算 diff、范围/秘密/文件类型与副作用检查）并逐项对账、失败路由表、幂等与崩溃窗口、**无变更执行的合法性**（chain 不以“事实摘要等于父快照摘要”为拒绝理由，该判别只属于端口）、**A5 前路由证据的兼容边界**、**`REVIEW_PENDING` 重入的资格证明**。
- `schemas/state-event.schema.json`：`taskTransition` 新增 `REPAIR_PENDING` 目标。
- 契约夹具：valid 新增 `state-event-task-postflight-rejected-valid`，invalid 新增 `state-event-task-repair-before-steps-invalid`（`PRECHECK→REPAIR_PENDING` 必须拒绝，reason 文本同步为“只能来自 REVIEW_PENDING、FAILED 或 postflight 拒绝”）；`tools/contracts/contracts.test.js` 夹具计数 126→128。

## 3. 事实布局与失败路由

冻结事实只在 completed 终局出现，位置固定：

$$\text{routeEvidence} = [\text{RequestHash},\ \text{CompletionHash},\ \text{Manifest},\ \text{Diff},\ \text{Log},\ \text{Usage},\ \text{ProcessStop},\ \ldots\text{Evidence}]$$

门按位置 2..6 重建事实（跨重启可恢复、不读 replay store）；最短 7 枚，缺失或歧义一律 fail-closed。

| 触发 | task | chain | reason | 证据 |
|---|---|---|---|---|
| 端口返回 failed | `REPAIR_PENDING` | `RUNNING→PAUSED` | `postflight-rejected` | 事实 + 端口证据 |
| 端口返回 uncertain | `FAILED` | 仅 RUNNING 时 `PAUSED` | `postflight-uncertain` / `recovery-uncertain` | 端口错误证据；`Run` 拒绝重试 |
| 端口返回 error | `FAILED` | `FAILED` | `postflight-failed` | — |
| 决策非法（缺事实/未知结论/摘要非法） | `FAILED` | `FAILED` | `postflight-invalid` | — |
| 事实缺失、多步归属歧义、证据过短 | `FAILED` | `FAILED` | `postflight-facts-missing` | — |
| 端口通过 | `REVIEW_PENDING`（证据含五枚事实） | 不变 | `postflight-passed` | 事实 + 端口证据 |
| 无 AgentRunner 事实的任务 | 原路径零变化 | 不变 | `steps-completed` | — |
| `REVIEW_PENDING` 重入且资格无法证明 | 不变（零写入） | 不变 | — | 返回 `ErrUnqualifiedReview` |

## 4. 核心不变量与聚焦测试

- 门是 review 的唯一资格来源：`TestPostflightGatePassesAndBindsTheFrozenFacts`（通过后 `REVIEW_PENDING` 转移证据含五枚事实摘要，且逐任务归属序列为 `[task-one, task-two, task-three]`）；`TestPostflightGateRefusesAnUnprovenPass`（缺任一事实的 passed 决策不再进入 acceptance）。
- 端口必填：`TestEngineRequiresAPostflightPort`。
- 事实可证明与 fail-closed：`TestAgentRunnerFrozenFactsValidateFailsClosed`、`TestFrozenFactsRebuildFromRoutingEvidence`（短证据/歧义证据拒绝）、`TestCompletedTerminalWithoutFrozenFactsIsRefused`、`TestAgentRunnerTerminalFactsMustNotCollideWithItsBindings`。
- 归一 attribution：`TestPostflightGateRefusesTasksWithSeveralRoutedSteps`（同一 task 的 completed 终局来自两个 step → 拒绝而不是猜测，端口零调用，且被拒记录本身是合法事件链）。
- 幂等与崩溃窗口：`TestPostflightGateConvergesAfterReviewPendingWasPersisted`（`REVIEW_PENDING` 已持久化后重入必须收敛且门不重跑）、`TestPostflightGateConvergesAfterPartialWrite`（门通过后写 `REVIEW_PENDING` 前崩溃 → 重入重跑一次门并收敛）。
- 重放安全：`TestAgentRunnerReplayWithDifferentFactsConflicts`（同 request+completion、异事实 → 冲突且零写入；同事实仍收敛）、`TestAgentRunnerReplayWithDamagedRoutingEvidenceConflicts`（记录证据过短无法重建事实 → 冲突）。
- 兼容与资格边界：`TestRequirePostflightQualifiedReviewFailsClosed`（无 review 转移 / 有转移但无事实 / 有事实三态）、`TestPostflightGateRefusesAnUnqualifiedReviewState`（A5 前式 `REVIEW_PENDING` 重入 → `ErrUnqualifiedReview`、零写入、零下游调用、门不重跑）。
- 端口判别边界：`TestPostflightGateHandsAParentDigestFactSetToThePort`（manifest 摘要等于父快照摘要的事实集不被 chain 预拒，由端口判定决定结局）。
- 适配器侧事实绑定：`TestAgentRunnerTerminalIntentFactsBindingFailsClosed`（completed 缺事实 / 非 completed 带事实 / 事实不互异 → 全部拒绝）、`TestAgentRunnerTerminalRecordValidation`（done 缺事实反例走公开构造器）、`TestAgentRunnerTerminalPublisherRefusesACompletedOutcomeWithoutFacts`（写 intent 与结算之前拒绝）。

## 5. 审查链与整改（闭环纪律：任一环节整改后重跑该环节直至通过）

- **③ V4 Pro 预审（第一轮）**：1 High（`runTask` 无条件写 `REVIEW_PENDING`，缺状态守卫，导致 `REVIEW_PENDING` 持久化后崩溃的重入尝试自转移被事件校验拒绝、永久失败）+ 1 Medium（completed 重放只比较 request+completion，事实分歧被静默收敛）+ 3 Low（`ErrorEvidence` 未校验摘要、事实与父摘要碰撞、夹具 reason 文本与计数过期）。整改：加 `current != "REVIEW_PENDING"` 守卫、completed 重放重建并比较事实、`ErrorEvidence` 摘要校验、夹具文本与计数修正，并补 `TestPostflightGateConvergesAfterReviewPendingWasPersisted`、`TestAgentRunnerReplayWithDifferentFactsConflicts` 等回归。
- **③ V4 Pro 复审（第二轮）**：High/Medium 清零；剩 4 条 Low（多步拒绝缺测试与契约句、`sameTerminalSettlementPlan` 排除事实缺注释、工作稿计数、损坏证据缺端到端测试）。整改：补 `TestPostflightGateRefusesTasksWithSeveralRoutedSteps` + 契约句、补注释、修正计数、补 `TestAgentRunnerReplayWithDamagedRoutingEvidenceConflicts`。
- **④ Codex 终审（第一轮）**：1 High + 2 Medium。(1) High：`facts.DisjointFrom(parent.Hash)` 是我自行加码、未经裁定的规则——`parent.Hash` 即已接受快照的 manifest 摘要，而按已拍板端口契约“重捕 manifest 并与父快照算 diff”，**无变更执行合法地与之相等**，该预拒可能误杀合法通过 → 撤销该检查，并把该判别明确划归端口（契约双语同步），原 ③ 的 Low 相应转为端口义务；(2) Medium：A5 前的 completed 路由证据不含事实前缀，重放按 conflict 零写入拒绝、不得就地续跑 → 写入契约作为兼容边界；(3) Medium：`TestPostflightDecisionRejectsMalformedEvidence` 首例被“缺第 5 枚事实”共因覆盖、对摘要循环不可证伪 → 改为“五枚事实全在 + 额外一枚非法摘要”。同时补齐适配器侧事实反例。
- **④ Codex 复审（第二轮）**：发现新 Medium——已处于 `REVIEW_PENDING` 的任务重入时跳过门直接进 acceptance，使 A5 前（或损坏存储中无门痕迹）的 review 状态可无 postflight 直达 PASSED，与新增契约句冲突 → 新增 `requirePostflightQualifiedReview` 与哨兵 `ErrUnqualifiedReview`：有 completed 路由终局的任务重入 `REVIEW_PENDING` 时，必须从该 `REVIEW_PENDING` 转移证据证明五枚事实齐全，否则零写入拒绝；无事实任务保持原路径。
- **④ Codex 复审（第三轮）**：`RE-REVIEW: PASS`，无 Medium+ finding；确认未新增状态/引擎入口、adapter 仍不能写策略或状态、门在 review 持久化后不重跑、无残留的父摘要预拒规则与悬挂测试名引用。

## 6. 变异检验（证明新测试确实可证伪）

| 变异 | 期望 | 结果 |
|---|---|---|
| 禁用 `agentRunnerTaskFacts` 的多步拒绝分支 | `TestPostflightGateRefusesTasksWithSeveralRoutedSteps` 变红 | ✅ 变红（`facts routed by several steps must be refused, got <nil>`） |
| 去掉 `runTask` 的 `REVIEW_PENDING` 守卫 | `TestPostflightGateConvergesAfterReviewPendingWasPersisted` 变红 | ✅ 变红（自转移被事件校验拒绝） |
| 去掉 `PostflightDecision.Validate` 的 `Evidence` 摘要循环 | `TestPostflightDecisionRejectsMalformedEvidence` 变红 | ✅ 变红（`got <nil>`） |
| 恢复 `facts.DisjointFrom(parent.Hash)` 预拒 | `TestPostflightGateHandsAParentDigestFactSetToThePort` 变红 | ✅ 变红（端口未被调用） |
| 去掉 `requirePostflightQualifiedReview` 的事实包含循环 | `TestRequirePostflightQualifiedReviewFailsClosed/without_the_facts` 与 `TestPostflightGateRefusesAnUnqualifiedReviewState` 变红 | ✅ 双红 |

每次变异后源码均按 SHA-256 校验还原（备份只放在 `tmp/` 并在验证后删除）。

## 7. 接受的边界与残余风险

- **端口实现质量是唯一信任边界**：chain 只能证明结构（位置重建、五枚互异、与请求/完成摘要互异、review 资格留有证据），无法证明端口的重捕是真实的，也无法仅凭摘要关系区分“无变更执行”和“adapter 复述父快照”；该判别按契约属于端口。
- 引擎现有 gates（build/test/verify hook step）排在 `Acceptance.Accept` 之前，与 CONTRACTS §7“freeze 后跑 gates”的措辞仍有顺序差异；本切片不重排，交 B4 终审对账。
- `sameTerminalSettlementPlan` 有意不比较事实：持久化 intent 是结算计划的唯一恢复源，事实分歧由 chain 路由层的冲突拒绝，已在实现处注释。
- 兼容边界：A5 前的 completed 路由记录不含事实前缀，重放按 conflict 零写入拒绝；A5 前的 `REVIEW_PENDING` 若属于有事实的任务则按 `ErrUnqualifiedReview` 零写入拒绝。两类运行的共同处置是**按新 attempt 重新执行**，不得就地续跑。
- `REVIEW_PENDING→Accept` 之间的 Accept 幂等属既有假设；本切片未改变。
- review 资格守卫基于已加载事件做摘要包含证明，依赖事件日志在 `New`/`Recover` 校验后保持 append-only 且不被带外篡改。

## 8. 门禁证据与复现

- `gofmt -l internal` 无输出；`go build ./...`、`go vet ./...`、`go test -count=1 ./...` 全部通过（Windows 原生，`GOTOOLCHAIN=local`，工具链 1.22.12）。
- 契约夹具门禁：`cd tools/contracts; npm test` → 4 个用例通过、夹具计数 128。
- 编码/行尾硬规则逐文件字节校验：改动过的 `.md` 为 UTF-8 **with BOM** + LF，`.go`/`.json`/`.js` 为 UTF-8 **without BOM** + LF。
- 复现命令：`go test -count=1 ./internal/chain/... ./internal/adapters/... ./internal/console/... ./internal/evidence/...`；契约门禁见上。
- 本切片为离线切片，**未使用付费探针**（额度 0/0），未调用真实 CLI、未访问网络。
- 远端 CI：GitHub Actions run 34950623175（commit `b8bad47`）全绿——Ubuntu 作业实际执行了 Race 与 Contract fixtures 步骤，Windows 作业的这两个步骤按 `runner.os == 'Linux'` 条件跳过。
