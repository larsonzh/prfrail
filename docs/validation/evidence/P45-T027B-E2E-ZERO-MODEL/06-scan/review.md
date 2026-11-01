# ⑥ 独立扫描（原样归档）

> 载体：runSubagent(agentName="independent-reviewer", model="MAI-Code-1.1-Flash (copilot)")。transcript 偏移 122618。

## A 结论
INDEPENDENT SCAN: PASS

## B 缺陷

| 编号 | 严重度 | file:line | 缺陷描述 | 修法 |
| --- | --- | --- | --- | --- |
|  |  |  |  |  |

## C 核验结论

1. T3 不启动的多路互证成立，且没有假绿空间。
   - 在 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 中，T2 失败后会检查：`workspaces.calls` 只有 T1/T2，`projection.TaskStates` 不含 T3，`store.events` 中不出现 T3 任务实体，`steps.requests` 中不出现 T3，并且 `len(steps.requests)==3` 这个断言强制卡住“额外发起”。
   - 在 [internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go) 中，控制台对 T2 用 `kind="code"`、`Mode=ManagedChangeSet` 的真实拒绝路径断言了 `localRuntime` 自然失败，且 `task-three` 不在 `taskIDs` 中。
   - 这不是单一字段判读，而是“任务启动序/事件流/投影键/步骤请求集”四重交叉，真实覆盖了“T3 never started”。

2. 零模型证明机械、且没有把计数器误当哨兵。
   - 引擎腿在 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 里直接断言：`options.StepIntents == nil`；并在执行时检查 `request.ExecutionTarget == DefaultExecution`，`request.AgentRunnerFacts == nil`，并且事件中不存在 `agentRunnerDispatchedReason`。
   - 预览腿在 [internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go) 中，计数器依然保留为 0，但测试明确说明它只是“README 运营承诺的静态锁定”，真正机械判断是 `RunnableInCLI`、`Outcome == ready`、`BlockingEvidence` 为空、且所有 `preview.Steps` 都是 `noop`。这符合“可证伪 + 计数器非哨兵”的要求。

3. 父快照绑定与事件链语义一致，且无越界。
   - [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 断言：T1 的父是 `baselineHash`，T2 的父是 `acceptedOne`，T3 的父是 `acceptedTwo`，终态 `AcceptedParent == acceptedThree`；这与真实链语义一致。
   - 事件链断言也直接调用了 [internal/evidence/event.go](internal/evidence/event.go) 的 `VerifyEventChain`，校验 `Sequence` 单调、`LastEventHash` 与最终事件 hash 一致，并确认 `seal` 数量 3。该证明与真实 engine state machine 一致，不是“测试造假”。

4. 并发腿是同 launcher/runRoot 的真正排他场景，且无墙钟判据。
   - [internal/console/runtime_agent_concurrent_e2e_test.go](internal/console/runtime_agent_concurrent_e2e_test.go) 复现的是同一 `launcher` 与同一 `runRoot` 下 8 个并发 start；赢家唯一，失败者只可看到 `ErrAgentRunnerLaunchConflict` / `ErrAgentRunnerProcessDuplicate`，并验证 `ProcessAlive`、slot owner 与 mirror 一致。
   - 这与 DR13-FIX 强输入一致，且没有使用 10s freshness 作为排他依据；它使用“状态/存在性/完成序”而非墙钟，是合理的。

5. 夹具和 README 与严格解码一致，且最小必要变更。
   - [examples/no-ai-preview/proofrail.chain.json](examples/no-ai-preview/proofrail.chain.json) 结构符合 [internal/console/config.go](internal/console/config.go) 的 `LoadChainConfig` 与 [schemas/chain.schema.json](schemas/chain.schema.json) 的 schema：`schemaVersion`、`chain`、`tasks`、`workspace` 全部合法，且文件大小 1044 bytes，起始 `0x7B 0x0A 0x20` 表明无 BOM，符合要求。
   - [examples/no-ai-preview/README.md](examples/no-ai-preview/README.md) 的修正是命令路径修正，最小必要，且与预览示例的行为一致。

6. 边界/卫生与证据包自洽性。
   - 实际只读边界被尊重：新测试文件是 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go)、[internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go)、[internal/console/runtime_agent_concurrent_e2e_test.go](internal/console/runtime_agent_concurrent_e2e_test.go)，生产代码未被改动；`git diff --name-only -- internal/` 在本轮检查中没有输出 tracked diff，只有这 3 个新测试文件分布在 internal tree，符合“零改动生产代码”的要求。
   - 证据包的 [docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/MANIFEST.md](docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/MANIFEST.md) 明确写着“步骤 0 骨架落盘；后续载荷按轮追加”，所以它当前是合法的 start-of-slice skeleton，而不是最终闭环包；该状态在文件中自说明，自洽。
   - Fresh verification evidence: `go test -count=1 ./...` returned all package lines as `ok` with no `FAIL` lines, including the internal/chain and internal/console packages.

## D 未覆盖面

- 这次只读复核未对真实 VS Code UI/浏览器路径做端到端点击验证；它依赖代码和事件日志的静态断言。
- 未对 Windows 之外的 OS/locale 运行环境做差异验证；并发 snapshot 与 ProcessAlive 的行为可能受平台差异影响。
- 证据包未进入最终 06/07/08 载荷阶段，因此未覆盖终审与 CI 归档闭环。
- 例子中的 CLI 命令路径已被核验为合法 JSON 例子，但未运行 `go run ./cmd/prfrail preview --chain ...` 的真实命令行交互输出。

## E 可证伪性

- 只要把 [examples/no-ai-preview/proofrail.chain.json](examples/no-ai-preview/proofrail.chain.json) 中任一 step 的 `kind` 改成 `code`，[internal/console/runtime_zero_model_e2e_test.go](internal/console/runtime_zero_model_e2e_test.go) 会立即让 `preview.Steps` 或 `BlockingEvidence` 断言失败。
- 只要让 T2 在 [internal/chain/engine_zero_model_e2e_test.go](internal/chain/engine_zero_model_e2e_test.go) 里继续推进到 T3 事件，`len(workspaces.calls)` 或 `task-three` 事件断言会立刻变红。
- 只要把并发测试里的“唯一赢家”逻辑去掉，或让 loser 观察到非 `Conflict`/`Duplicate` 错误，[internal/console/runtime_agent_concurrent_e2e_test.go](internal/console/runtime_agent_concurrent_e2e_test.go) 会直接失败。

§A: INDEPENDENT SCAN: PASS