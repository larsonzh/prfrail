# T027 · A4 · chain terminal routing and resume — 验证报告

日期：2026-09-15。结论：`A4 本地实现与门禁 PASS`；提交/推送与 CI 证据需用户同轮显式授权后回写。
切片目标：把终局语义收敛到 chain 核心路由与恢复，禁止 adapter receipt 直接驱动任务结论或跳过下游 acceptance。

## 1. 交付范围

- 新增 chain 自有终局事实 `AgentRunnerTerminal`（五态：completed/failed/cancelled/operator-action-required/uncertain），绑定 RequestID/RequestHash/CompletionHash/Run/Task/Step/Attempt/SessionID/Prior* 与去重有序证据哈希，不含结算摘要与 wire 细节。
- 新增 step 状态 `TERMINAL_PENDING`：`RUNNING→TERMINAL_PENDING`（证据=dispatch 结果哈希）表示“已派发、等终局”；`TERMINAL_PENDING→WAITING_FOR_OPERATOR|PASSED|FAILED|CANCELLED`。非可停留调度态。
- 新增唯一外部终局写入口 `Engine.SubmitAgentRunnerTerminal`；runStep 对 AgentRunner 步骤改为“派发成功→`TERMINAL_PENDING`→返回 `ErrAwaitingAgentRunnerTerminal`”，不再走 step-passed。
- 新增 adapters 单向翻译 `ToChainAgentRunnerTerminal` / `LoadChainAgentRunnerTerminal`（adapters→chain），chain 不读 replay store。
- 新增/改动文件：`internal/chain/agent_runner_terminal.go`、`internal/chain/agent_runner_terminal_route.go`、`internal/chain/agent_runner_terminal_test.go`、`internal/adapters/agent_runner_terminal_chain.go`、`internal/adapters/agent_runner_terminal_chain_test.go`；`internal/chain/engine.go`、`internal/chain/operator_interaction.go`、`internal/chain/recover.go`、`internal/evidence/event.go`、`internal/chain/engine_test.go`、`schemas/state-event.schema.json`、`testdata/contracts/{valid,invalid}/runtime-records.json`、`tools/contracts/contracts.test.js`、`docs/CONTRACTS{,_EN}.md`、`docs/t027/REMAINING_SLICES{,_EN}.md`。

## 2. 协议先行产物

- `docs/CONTRACTS.md` / `_EN.md` §2.1：step 行更新为 `RUNNING→TERMINAL_PENDING|WAITING_FOR_OPERATOR|PASSED|FAILED|CANCELLED`、新增 `TERMINAL_PENDING→WAITING_FOR_OPERATOR|PASSED|FAILED|CANCELLED`，并注明 `TERMINAL_PENDING` 非可停留态、重启无终局按不确定暂停。
- 同章节新增 §7“chain terminal 路由（A4）”段：路由表、operator 路由“路由即暂停”、Open 侧加法式接受与账本排他边界、uncertain 暂停边界、幂等/冲突、resume 连续性判据与 U-9 边界。
- `schemas/state-event.schema.json`：`stepTransition` 新增 `TERMINAL_PENDING` 分支并扩展 `RUNNING` 目标集。
- 契约夹具：新增 valid `state-event-step-terminal-pending-valid`（RUNNING→TERMINAL_PENDING）、valid `state-event-step-terminal-pending-passed-valid`（TERMINAL_PENDING→PASSED）、invalid `state-event-step-terminal-pending-resume-invalid`（TERMINAL_PENDING→RUNNING 必须拒绝）；`tools/contracts/contracts.test.js` 夹具计数 123→126，与 `internal/evidence/event.go` 转移表逐行一致。

## 3. 路由表与核心不变量（实现口径）

| 终局 | step | task | chain | 备注 |
|---|---|---|---|---|
| completed | `TERMINAL_PENDING→PASSED` | 不变（`STEPS_RUNNING`） | 不变（`RUNNING`） | 任务仍需 freeze/gates/review/completed promotion，下一次 `Run` 走既有下游流 |
| failed | `TERMINAL_PENDING→FAILED` | `FAILED` | `FAILED` | 不自动重试/relaunch |
| uncertain | `TERMINAL_PENDING→FAILED` | `FAILED` | `RUNNING→PAUSED`（reason `recovery-uncertain`） | 保持暂停且禁止重试（`Run` 直接 `ErrRecoveryUncertain`）；chain 已暂停时不改写暂停原因 |
| cancelled | `TERMINAL_PENDING→CANCELLED` | `CANCELLED` | `CANCELLED` | 先归档停机证据（`Stopper.Stop`），失败则 step 保持 parked 且 chain 暂停 |
| operator-action-required | `TERMINAL_PENDING→WAITING_FOR_OPERATOR` | `STEPS_RUNNING→WAITING_FOR_OPERATOR` | `RUNNING→PAUSED` | 路由即暂停，交接给 T025 机器；证据含 C 摘要 |

不变量证据（聚焦测试）：

- 不直通：`TestAgentRunnerCompletedTerminalAloneDoesNotCompleteTask`（Submit 后 task 仍 `STEPS_RUNNING`、下游端口零调用；下一次 `Run` 才完成）；路由矩阵对五种终局逐一断言“路由期零下游调用”。
- 派发只 park：`TestAgentRunnerDispatchParksStepInsteadOfPassingIt`（`Run` 重入不重派发、不写事件）。
- 绑定与冲突：`TestAgentRunnerTerminalRefusesUnboundAndConflictingReceipts`（未 park、非 AgentRunner step、RequestHash 不匹配、同 request 异 completion 全部零写拒绝；重复终局 `Converged=true` 零写）。
- 崩溃收敛：`TestAgentRunnerTerminalConvergesAfterPartialWrite`（step 已写/task 未写）、`TestAgentRunnerTerminalConvergesAfterTaskWriteBeforeChainWrite`（step+task 已写/chain 未写，重放只补 1 条缺失转移）。
- resume 连续性：`TestAgentRunnerTerminalResumeContinuityProof`（`PriorCompletionHash` 不在 step 事件证据历史→零写拒绝；证明后放行；operator 已答复后的终局重放不回压等待态）。
- operator 接力：`TestAgentRunnerOperatorRouteSurvivesRecoveryAndInteractionOpen`（route→重启→`Recover`→`OpenOperatorInteraction`→`ResumeOperatorInteraction`→`Run` 全链闭环，且 `RecoveryUncertain` 保持 false）。
- 停机证据：`TestAgentRunnerTerminalCancelledArchivesStopEvidence`（CANCELLED 转移证据含 stop 摘要与 C 摘要）、`TestAgentRunnerTerminalCancelledArchiveFailuresStayParked`（停机失败→保持 parked + chain 暂停 + `ErrRecoveryUncertain`）。
- 重启不确定：`TestAgentRunnerRecoverTreatsParkedStepAsUncertain`（`TERMINAL_PENDING` 进入 recover 不确定清单；之后仍可路由已证明终局，但不静默续跑）。
- 翻译层：`TestToChainAgentRunnerTerminalMapsThePublishedChain`（身份/状态/会话/证据映射）、`TestToChainAgentRunnerTerminalFailsClosed`（零值、intent 绑定、结算槽位不一致、resume 会话不等）、`TestLoadChainAgentRunnerTerminalRequiresASettledClosure`（未结算/未关闭/缺 store 拒绝；翻译对 store 严格只读——整目录 path→size 快照前后相等）。

## 4. 审查链（协议 ③④ 闭环）

1. **V4 Pro 预审（③）**：`PRE-REVIEW: FINDINGS` — 1 Medium + 2 Low。
   - Medium：operator 路由只写 task/step 等待、chain 仍 `RUNNING`，若在交互记录出现前崩溃，`Recover` 会写 `recovery-uncertain` 暂停，此后 `OpenOperatorInteraction` 与 `Run` 都无法离开 → 路由永不能收敛。
   - Low 1：Open 侧对 agent 等待态不校验交互记录哈希（排他性由控制台账本承担），需文档化；Low 2：uncertain 终局在 chain 已因其他原因暂停时不会到达 `PAUSED`，属边界需文档化。
2. **整改**：路由同步 `RUNNING→PAUSED`（证据含 C 摘要）；新增加法式 `operatorWaitingPauseMatches`（经典绑定优先）；契约双语补齐两条边界。
3. **V4 Pro 复审**：`RE-REVIEW: FINDINGS` — 无 Medium+，但新增 1 Low（真实缺口）：operator 机器已归还控制权后重放同一终局会把 task/chain 重新压回等待。
4. **整改**：`convergeAgentTerminalRoute` 在 step 已不处于 `WAITING_FOR_OPERATOR` 时视为已收敛（零写），新增 `TestAgentRunnerTerminalResumeContinuityProof` 内的 superseded 重放断言。
5. **Codex 独立终审（④）**：`FINAL REVIEW: FINDINGS` — 1 Medium + 1 Low + 3 项测试覆盖建议。
   - Medium：翻译层只绑定 closure 的 requestId/intentRecordHash/completionHash，未绑定结算槽位（`settlementEntryId`/`settlementIdempotencyKey`），重算 recordHash 的篡改 closure 可穿过翻译。
   - Low：声称 4 个 `.md` 缺 BOM（见 §5 边界 5，字节级实测为误报）。
   - 覆盖建议：新增“task 已写、chain 未写”收敛反例；把 route→restart→Recover→Open→Resume 补上最终 `Run`。
6. **整改**：翻译层补齐结算槽位交叉绑定（与 `AgentRunnerReplayStore.validateTerminalClosureBinding` 同口径）并新增两条真反例（重算 recordHash 的内部有效记录）；新增 `TestAgentRunnerTerminalConvergesAfterTaskWriteBeforeChainWrite`；operator 恢复测试补齐最终 `Run`。
7. **Codex 复审（④ 闭环）**：`RE-REVIEW: PASS`，无发现；确认 Medium 已关闭、两条结算反例为真反例、崩溃窗口测试可证伪、无新增问题。R4（BOM）Codex 声明其工具无法取得字节级证据、拒绝在未证实情况下判定真伪（如实记录）。

结论：Medium+ 全部由独立终审方（Codex）修复后复审通过，无自审自证。

## 5. 已知边界（有意保留，已写入契约或说明）

1. **attempt 限第一轮**：chain 目前只建模 attempt 1（`taskEntity`/`stepEntity` 固定），提交非 1 的终局被拒；新 attempt 由 CLI 从持久 envelope 构建（U-9），A6/后续切片接管。
2. **Open 侧绑定强度**：agent 路由写入的等待转移绑定的是终局回执（当时交互记录尚不存在），故 Open 不校验交互记录哈希；交互排他性仍由控制台账本（`OperatorInteractionLedger`）承担。已写入 CONTRACTS §7。
3. **uncertain 与既有暂停**：chain 转移表不允许 `PAUSED→PAUSED`，故 uncertain 只暂停处于 `RUNNING` 的 chain；若 chain 已因其他原因暂停，恢复动作会在 task `FAILED` 边界停止（`runTask` 快速失败），不会重试该 step。已写入 CONTRACTS §7。
4. **Run 与等待态**：在交互等待期间调用 `Run` 会先把 chain 复为 `RUNNING`、随即以 `repair-required` 重新暂停，导致该暂停不再匹配交互绑定——与经典 T025 路径同构（既有行为），使用方不应在等待期调用 `Run`。
5. **BOM 误报澄清**：Codex Low 称 4 个 `.md` 缺 BOM。仓内字节级实测为 `EF BB BF` 起始（`[System.IO.File]::ReadAllBytes` 前三字节）且无 CRLF，符合 `docs/CODING_CONVENTIONS.md`；该 Low 判定为文本级读取（BOM 被读取器剥离）造成的误报，未做改动。
6. **publisher 与 chain 解耦**：A3 发布器不依赖 chain；A4 只通过翻译层把“已发布且已关闭”的终局事实交给 chain，receipt 本身零状态副作用。

## 6. 门禁证据（本地，Windows + Go 1.22 工具链）

| 命令 | 结果 |
|---|---|
| `gofmt -l internal` | 空输出 |
| `go build ./...` | 通过 |
| `go vet ./...` | 通过 |
| `go test -count=1 ./...` | 全部包通过（含 `internal/chain`、`internal/adapters`、`internal/release` 冻结契约测试） |
| `npm test`（`tools/contracts`） | 4/4 通过（126 个契约夹具，含新增 3 条） |
| `Select-String`/字节检查（`.md`） | `docs/CONTRACTS{,_EN}.md`、`docs/t027/REMAINING_SLICES{,_EN}.md`、本报告：BOM=True、CRLF=False |

Linux CI（含 `Race` 步骤与 `Contract fixtures` 步骤）证据待授权提交/推送后回写。

## 7. 复现方式

```text
go build ./... && go vet ./... && go test -count=1 ./...
# 聚焦：go test -count=1 -run 'AgentRunnerTerminal|ChainAgentRunner' ./internal/chain/... ./internal/adapters/...
cd tools/contracts && npm test
```

## 8. 未做 / 后续

- A5 postflight acceptance 集成；A6 pinned CLI adapter 与 dispatcher/admission 的 resume 放行（U-7 明确不在 A4 放开）；B 系列真实候选与耐久性证明。
- T027 仍为 `BLOCKED / NOT IMPLEMENTED`，AT-23 未通过。
