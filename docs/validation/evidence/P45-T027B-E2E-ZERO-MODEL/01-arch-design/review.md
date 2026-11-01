# ① 架构方案（原样归档）

> 载体：runSubagent(agentName="deep-reasoner", model="DeepSeek V4 Pro (deepseek)")。transcript 偏移 121811。输入：`tmp/P45-T027B-E2E-ZERO-MODEL/batch1-assembly-census.md`、`batch2-port-contracts.md`。

## A 分析与方案

**E2E 落点取舍结论：主落点 (a) `internal/chain` 引擎级 E2E；(b) 不改、也不应改 console 单任务限制；console 侧用既有多任务能力的 `ExecuteNoopRun` 补两条旁证腿。**

- 装配层「exactly one task」是**引擎存在前的拒绝前提**：`agentRunnerAssemblyStep`（`runtime_agent.go:774-790`）与 `supervisedAssemblyStep`（`runtime_supervised.go:335-338`）均以 `ErrInvalidConfig` 拒绝多任务，是 P1「装配层半壁」的既有契约。
- 放开 `:774`/`:335` 与 A6 启动器排他设计**自相冲突**：durable identity mirror 每个 runRoot 唯一（`agent_runner_pinned_cli.go:28-30` 注释），同 runRoot 第二 launch 被 `registry.live()` 拒绝（`pinned_cli.go:123-129`）；多任务装配将同 runRoot 并发双 launch，直接踩 DR13-FIX 刚钉住的排他契约。
- 回归风险：放开会红 `runtime_agent`/`runtime_supervised` 现有「多任务拒绝」用例并扰动 DR-9 敏感面。本片方案对 39 个链用例与 `runtime_agent_e2e_test.go` **零既有文件修改**——新文件同包编译、只读复用 `threeTaskDefinition`（`engine_test.go:202-207`）与 `testOptions`（`213-232`），回归风险≈0。
- 关键事实：`ExecuteNoopRun`（`runtime.go:33-64`）**本身无单任务限制**，`localRuntime`（`runtime_noop.go:17-46`）零 AI；console 已有多任务合法通道，无需改装配。

## B 三级/三任务判定表

| 级 | 状态 | 触发 | 断言 |
|---|---|---|---|
| T1 | `PASSED` | 默认全绿 fakes | `TaskStates["task-one|1"]=="PASSED"`；`workspaces.calls[0].parent==baselineHash`；T1 PASSED 事件 `InputEvidence[0]==acceptedOne` |
| T2 | `FAILED` | stub 在 `StepPort.Execute` 按 `request.TaskID=="task-two"` 返回 `errors.New("runner failed")`（`engine_test.go:64-66` 先例 `failTask` 分派） | T2 FAILED 事件 `FromState=="STEPS_RUNNING"`、`Reason.Code=="step-failed"`；`ChainState=="FAILED"`；`Run` 返回非 nil |
| T3 | 不启动 | 引擎回路短路 | `len(workspaces.calls)==2 && workspaces.calls[1].task=="task-two"`；事件流无 `{Kind:task, TaskID:task-three}` 事件；`TaskStates` 无 `taskKey("task-three",1)` 键；`steps.requests` 无 task-three |

## C 改动点清单

- **新增** `internal/chain/engine_zero_model_e2e_test.go`：主 E2E（T1/T2/T3 + 父快照 + 事件链 + 零模型）。复用包内 `testOptions`/`threeTaskDefinition`/`fakeSteps.failTask`，**不改** `engine_test.go`。
- **新增** `internal/console/runtime_zero_model_e2e_test.go`：① 三 noop 任务链 `ExecuteNoopRun` → `COMPLETED`、3×`PASSED`、`LoadRunSummary` 验链；② T2 含 `code` 步 → `localRuntime.Execute` 返回 `"local runtime cannot execute step kind"`（`runtime_noop.go:41-43`），走真实失败机 → 摘要 `FAILED`、task-three 缺席；③ `executePreview` 跑 `examples/no-ai-preview/`（README:1-16 判据）断言五计数器全 0。
- **新增** `internal/console/runtime_agent_concurrent_e2e_test.go`：DR13-FIX 强输入，复用 `buildAssemblyStub`（`sync.Once`，`runtime_agent_e2e_test.go:47`）。
- **修改/删除：无**。`internal/adapters/**` 零改动（仅只读调用导出 API）。

## D 红测清单（8 条）

1. `TestZeroModelE2ET1PassT2FailT3NeverStarted`：上表三级断言。未实现时：`workspaces.calls==3` 或出现 task-three 事件。
2. `TestZeroModelT1ParentBinding`：`calls[0].parent==baselineHash` + T1 PASSED 事件 `InputEvidence[0]==acceptedOne` + T2 物化 `calls[1].parent==acceptedOne`。红形态：父子链断链。
3. `TestZeroModelEventChainSealedMonotonic`：`evidence.VerifyEventChain==nil`（`event.go:113` 拒绝重复 ID）、`events[i].Sequence==i+1`、`projection.Sequence==len(events)`、`LastEventHash==末事件 EventHash`。红形态：`ErrBrokenEventChain` 或跳号。
4. `TestZeroModelNoAgentRunnerDispatch`：所有 `steps.requests` 的 `ExecutionTarget==DefaultExecution && AgentRunnerFacts==nil`，`fakeStepIntents` 不产 facts。红形态：任何 intent 泄漏即红。
5. `TestZeroModelConsoleNoopThreeTaskE2E`：console 三任务全 noop → 3×`PASSED`。红形态：console 多任务驱动回归。
6. `TestZeroModelConsoleNaturalFailureStopsLaterTasks`：T2 code 步 → 摘要 `FAILED`、task-three 缺席、FAILED 事件 `Reason.Code=="step-failed"`。红形态：T3 被启动或非 FAILED 终态。
7. `TestZeroModelPreviewCountersAllZero`：五计数器全 0。红形态：任一 >0 即红。
8. `TestZeroModelConcurrentStartOneWinner`：见 E。红形态：2 赢家 / 败者错误未命中 Conflict\|Duplicate / 赢家未存活 / 停后 slot 未释放。

## E 并发 Start 断言设计（DR13-FIX 强输入）

**调用面（全部 exported，只读）**：`adapters.NewAgentRunnerReplayStoreForRun(runRoot, runID)`（`runtime_agent.go:153-157` 先例）→ `adapters.NewAgentRunnerPinnedCLILauncher(config, store)`（`pinned_cli.go:89`）→ `StartAgentRunnerProcess(ctx, chain.AgentRunnerLaunchRequest{RequestID,...})`（`:112`）→ `StopAgentRunnerProcess(ctx, launchID)`（`:291`）；错误哨兵 `ErrAgentRunnerLaunchConflict`/`ErrAgentRunnerProcessDuplicate`（`:27-34`）；存活判定 `guard.InspectProcess`+`guard.ProcessAlive`（`guard/process.go:361`；console 先例 `runtime_agent.go:643-650`）。

**断言流程**（同一 launcher + 同一 runRoot，8 goroutine 经 barrier 同步起跑、互异 `RequestID`）：
1. **恰一赢家**：恰好 1 个 Start 返回 `(result, nil)`；
2. **败者全命中**：其余 7 个 `errors.Is(err, Conflict) || errors.Is(err, Duplicate)`；
3. **赢家 `ProcessAlive`**：`ProcessID` 是 `guard.ProcessIdentity` 的 canonical JSON（`agent_runner_process_registry.go:222-236` 编码），测试用 `evidence.DecodeStrictJSON` 解码后 `guard.ProcessAlive(identity)==true`（**无任何墙钟等待/重试**）；
4. **closeout 时 registry.len()==1 的包外等价**：赢家存活期间追加第 9 个 Start（新 requestID）→ `Conflict`（证 registry 非空且唯一 live 即赢家）；随后 `StopAgentRunnerProcess(winner.LaunchID)` → 解码身份 `ProcessAlive==false` → 同 launcher 新 Start 成功（证 slot 已释放，行为等价 len==0）；
5. **slot owner token 与 mirror 一致**：字面断言仅包内可达（`registry`/mirror 读写器均未导出）——已由 DR13-FIX 包内测试钉住（`agent_runner_pinned_cli_test.go:273-274`、`:242`/`:305`），P45 侧以「Stop 按 LaunchID 精确终止赢家且随后可重启」作行为等价，并以 `go test ./internal/adapters -run Concurrent -count=10` 重跑作回归腿；
6. **禁止 10s freshness 墙钟**：`agentRunnerLaunchSlotFreshness`（`pinned_cli.go:52-56`）不参与任何断言；无 `time.Sleep`、无按时间采样。

**不可行说明**：`registry.len()` 与 slot-token/mirror 为 `internal/adapters` 包私有态，P45 边界禁改 adapters，故只能在包内测试断言（已存在）或在包外以行为等价覆盖——建议按上述替代方案执行，无需开边界。

## F 风险与未确证项

1. **39 个链用例**：新测试若向共享 fakes（`fakeSteps` 等）加字段/改行为，会扩散影响；对策：新文件**只读复用**、必要时自带最小 stub，不动 `engine_test.go` 任何类型。
2. **`runtime_agent_e2e_test.go`**：并发新测试不得 `t.Parallel`（共享 agent-stub 构建与 runRoot）；`sync.Once` 已防重复构建，但并发 Start 测试自身需独占 runRoot。
3. **DR-9 敏感面**：identity mirror 每 runRoot 唯一，并发测试必须单 launcher 单 runRoot；跨平台差异（`guard.ProcessAlive` 在 Windows/Unix 语义不同，DR13-FIX 已遇 `//go:build unix` 本机不编译问题）——以 ubuntu CI 腿为最终判据（DR13-FIX R2 先例）。
4. **未确证项**：① `engine.go:162-164`（Run 短路行号）与 `pinned_cli.go:112/:291` 行号来自批量普查+本次内容核对，实现时须以 `grep -n` 复核；② `localRuntime.Execute` 拒绝非 noop 种类的报错文案若变，红测 6 需同步；③ `reason=="step-failed"` 依赖 `runStep`→`failTask` 传参（`engine.go:564-573`），钉住后属有意契约；④ console 预览腿依赖 `examples/no-ai-preview` 配置不被改（现 README 判据成立）。

DESIGN: DONE