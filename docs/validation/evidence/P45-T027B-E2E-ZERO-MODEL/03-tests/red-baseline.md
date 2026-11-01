# ③ 首轮红测基线（原样归档）

> 载体：runSubagent(agentName="prfrail-tester", model="DeepSeek V4.1 Flash (deepseek)")。transcript 偏移 122083。首轮 7 绿 1 红，红点为 `TestZeroModelPreviewCountersAllZero`（`examples/no-ai-preview/proofrail.chain.json` 缺失，真实缺陷）。

## 1) 文件与用例清单

| 文件 | 行数 | 新增用例 |
|---|---|---|
| `internal/chain/engine_zero_model_e2e_test.go` | 192 | `TestZeroModelE2ET1PassT2FailT3NeverStarted` (L21)<br>`TestZeroModelT1ParentBinding` (L82)<br>`TestZeroModelEventChainSealedMonotonic` (L114)<br>`TestZeroModelNoAgentRunnerDispatch` (L162) |
| `internal/console/runtime_zero_model_e2e_test.go` | 183 | `TestZeroModelConsoleNoopThreeTaskE2E` (L66)<br>`TestZeroModelConsoleNaturalFailureStopsLaterTasks` (L103)<br>`TestZeroModelPreviewCountersAllZero` (L148) |
| `internal/console/runtime_agent_concurrent_e2e_test.go` | 164 | `TestZeroModelConcurrentStartOneWinner` (L72)（+ 自执行 helper `TestZeroModelAgentCLIHelper` L24） |

变更范围核对：`git status --porcelain` 中 `internal/` 下仅有这三个 `_test.go` 新增文件；未改任何生产代码、既有测试或 `internal/adapters/**`。

## 2) 四条命令输出

**① `gofmt -l internal/chain internal/console`** → 空（无输出，退出码 0）

**② `go vet ./internal/chain/ ./internal/console/`** → 无输出

**③ `go test -count=1 -run 'TestZeroModel' ./internal/chain/ ./internal/console/ -v`（过滤后）**
```
=== RUN   TestZeroModelE2ET1PassT2FailT3NeverStarted
--- PASS: TestZeroModelE2ET1PassT2FailT3NeverStarted (0.00s)
=== RUN   TestZeroModelT1ParentBinding
--- PASS: TestZeroModelT1ParentBinding (0.00s)
=== RUN   TestZeroModelEventChainSealedMonotonic
--- PASS: TestZeroModelEventChainSealedMonotonic (0.00s)
=== RUN   TestZeroModelNoAgentRunnerDispatch
--- PASS: TestZeroModelNoAgentRunnerDispatch (0.00s)
PASS
ok      github.com/larsonzh/prfrail/internal/chain      1.842s
=== RUN   TestZeroModelAgentCLIHelper
--- PASS: TestZeroModelAgentCLIHelper (0.00s)
=== RUN   TestZeroModelConcurrentStartOneWinner
--- PASS: TestZeroModelConcurrentStartOneWinner (1.41s)
=== RUN   TestZeroModelConsoleNoopThreeTaskE2E
--- PASS: TestZeroModelConsoleNoopThreeTaskE2E (0.07s)
=== RUN   TestZeroModelConsoleNaturalFailureStopsLaterTasks
--- PASS: TestZeroModelConsoleNaturalFailureStopsLaterTasks (0.05s)
=== RUN   TestZeroModelPreviewCountersAllZero
--- FAIL: TestZeroModelPreviewCountersAllZero (0.00s)
FAIL
FAIL    github.com/larsonzh/prfrail/internal/console      3.464s
FAIL
```

**④ RED 原文（`go test -run TestZeroModelPreviewCountersAllZero -v`）**
```
=== RUN   TestZeroModelPreviewCountersAllZero
    runtime_zero_model_e2e_test.go:151: the documented no-AI preview example config is missing
    (..\..\examples\no-ai-preview\proofrail.chain.json): GetFileAttributesEx
    ..\..\examples\no-ai-preview\proofrail.chain.json: The system cannot find the file specified.
    (README.md:9 tells the operator to run `--chain .\proofrail.chain.json`)
--- FAIL: TestZeroModelPreviewCountersAllZero (0.00s)
```
根因分析：`git ls-files examples/no-ai-preview` 只返回 `README.md`，示例目录**没有** `proofrail.chain.json`；README 第 9 行的用法命令指向一个不存在的文件 ⇒ 该示例按文档不可运行。计数器本身在 `preview.go` 中硬编码为 0（`BuildPreviewReport` 返回 `CallCounters{0,0,0,0,0}`），因此缺陷不在计数器，而在**样例夹具缺失**。

## 3) RED / GREEN 判定表

| # | 用例 | 判定 | 依据 / 证伪点 |
|---|---|---|---|
| 1 | `TestZeroModelE2ET1PassT2FailT3NeverStarted` | **GREEN** | 能力已存在（失败即停链已贯通）：T1 `PASSED`、T2 `FAILED`、`ChainState=FAILED`、`Run` 返回非 nil；`workspaces.calls` 恰 2 条（T1、T2），无 `task-three` 的 task/step 事件、无投影键、无 step 请求（`steps.requests` 恰 3 条）。删掉「runTask 出错即 return」的守卫 ⇒ T3 会被物化，必红。 |
| 2 | `TestZeroModelT1ParentBinding` | **GREEN** | `calls[0].parent == baselineHash`、`calls[1].parent == acceptedOne`，并追加 T2→T3 与终态 `AcceptedParent==acceptedThree`；若 parent 取错（例如仍用 baseline）必红。 |
| 3 | `TestZeroModelEventChainSealedMonotonic` | **GREEN** | `evidence.VerifyEventChain`（实际函数名即此）通过；`Sequence == index+1` 无跳号；`projection.Sequence == len(events)`；`LastEventHash == 末事件 EventHash`；另有 3 条 `ReviewDecision=="approve" && PromotionDecision=="completed"` 的封存记录（删掉 `sealAcceptedClosure` 即红）。 |
| 4 | `TestZeroModelNoAgentRunnerDispatch` | **GREEN** | 先断言装配未挂 `Options.StepIntents`（AgentRunner 目标唯一来源），再断言所有 `steps.requests` 的 `ExecutionTarget == DefaultExecution`（空串）且 `AgentRunnerFacts == nil`，且事件中无 `agentRunnerDispatchedReason`。 |
| 5 | `TestZeroModelConsoleNoopThreeTaskE2E` | **GREEN** | `ExecuteNoopRun` 三 noop 任务 `COMPLETED`、`PASSED` 计数 3；并用 `LoadRunSummary` 复核落盘摘要 `Sequence/LastEventHash` 一致（摘要与磁盘不一致会红）。 |
| 6 | `TestZeroModelConsoleNaturalFailureStopsLaterTasks` | **GREEN** | `localRuntime.Execute` 自然拒绝 `code` 步骤（错误含 `cannot execute step kind`）；重放事件得 `ChainState=FAILED`，任务态恰 2 个（1 PASSED + 1 FAILED），事件流无 `task-three`。若把非 noop 步骤静默跳过则必红。 |
| 7 | `TestZeroModelPreviewCountersAllZero` | **RED** | 夹具缺失（见上）；`examples/no-ai-preview/README.md:9` 承诺的 `proofrail.chain.json` 不存在。 |
| 8 | `TestZeroModelConcurrentStartOneWinner` | **GREEN** | 同 launcher / 同 runRoot、8 goroutine barrier 齐发：恰 1 赢家、7 输家且全为 `ErrAgentRunnerLaunchConflict`/`ErrAgentRunnerProcessDuplicate`；赢家 identity 解码后 `guard.ProcessAlive==true`；存活期追加 Start ⇒ `Conflict`；`Stop` 后 `ProcessAlive==false`；随后新 Start 成功且 LaunchID 不同。未用 `t.Parallel`，未用 `time.Sleep` 做同步（仅子进程自身生命周期 sleep），runRoot 为独占 `t.TempDir()`。 |

## 4) 因既有 helper 不符而改写的地方

- `evidence.VerifyEventChain`：已按实际函数名使用（`internal/evidence/event.go:96`），无改写。
- `testOptions` / `threeTaskDefinition` / `fakeSteps.failTask`（`engine_test.go:202-232`、`64-66`）：签名与预期一致，直接复用（`testOptions` 返回 9 个值）。
- **并发腿无法复用 adapters 包内测试 helper**（`agentRunnerPinnedCLIConfigForMode`、`TestAgentRunnerPinnedCLIHelper`、`replayStoreForRunRoot` 均为 `package adapters` 的测试内私有符号，console 测试二进制拿不到）。故在 console 侧自建：`zeroModelAgentCLIVersion` + 自执行 helper `TestZeroModelAgentCLIHelper`（`version`/`run` 两模式），配置 `Executable=os.Args[0]`；仅使用 adapters **导出** API（`NewAgentRunnerReplayStoreForRun` / `NewAgentRunnerPinnedCLILauncher` / `Start…` / `Stop…`）与导出哨兵。
- runRoot 前置条件：复用 console 既有 `ensureAgentRunnerEventLogMarker(runEventLogPath(runRoot))`（`runtime_agent.go:888`）写入 replay store 要求的 run-root 标记，未新建基础设施。
- `TaskStates` 键：console 侧不引用 chain 私有 `taskKey`，改用「投影键数量 + 状态分布计数」断言（避免依赖 `\x00` 键格式）。

## 5) 暴露的真实缺陷（未改生产代码）

- **缺陷 D1（夹具缺失）**：`examples/no-ai-preview/README.md:9` 声称在仓库根执行 `go run ./cmd/prfrail preview --chain .\proofrail.chain.json --json` 即可看到 `callCounters` 全 0，但 `examples/no-ai-preview/` 下只有 `README.md`，`proofrail.chain.json` 缺席（`git ls-files examples/no-ai-preview` 仅返回 README）。症状：按文档执行必报配置不存在；本片测试 `internal/console/runtime_zero_model_e2e_test.go:151` 即因此失败。修复方向（超出 ③ 权限，交 ①/⑧ 定）：补 `examples/no-ai-preview/proofrail.chain.json`（纯 noop 链，`outcome=ready`），或修正 README 指向真实样例。**注**：五个计数器本身在 `preview.go` 中是常量 0，能力层无缺陷。

TESTS: MIXED — RED: `TestZeroModelPreviewCountersAllZero`（样例夹具缺失）；GREEN: 其余 7 条。