# P50 ⑤ 预审原位重跑（r2）报告

## §1 自跑读数（磁盘快照锚点）

- **锚点**：HEAD `8bd8224`；工具链 go1.27.0/windows（仓库门禁为 1.22，CI 权威，见 INFO-1）。
- **快照**：CONTRACTS 对已暂存 +1/−0；`runtime_agent.go` +748/−177、`runtime_supervised.go` +485/−219（较 ②b/②c 回执增量 = 修复轮改动）；e2e +8/−7、supervised_test +14/−6；`REMAINING_SLICES{,_EN}` +22/−0；6 项未跟踪（4 个 console 测试文件、registry 测试、证据包目录）。
- **自跑**：`gofmt -l` 10 文件空输出；`go build ./...` EXIT=0；`go vet` 三包 EXIT=0；聚焦测试全 PASS（driver 三夹具 0.55s / H1 0.19s / M1 0.03s / M2 harness 0.33s / offline harness / funnels / registry；`TestAgentRunnerAssemblyFullChain` Windows 按既有 durability 闸 SKIP）；**全量 `go test -count=1 ./...` 15 包全 ok、EXIT=0**（console 22.6s / adapters 51.6s / chain 6.6s）。
- **机械核验**：`return fail(` = 8（与不变量测试一致）；`requires exactly one task` 两生产文件 = 0 处；5 个新 `.go` 无 BOM+LF、证据 `.md` BOM+LF，编码全部合规。

## §2 逐项 CLOSED 判定表

| 项 | 判定 | 证据（file:line） | 新发现 |
|---|---|---|---|
| H1 | ✅ 属实 | `runtime_agent.go` `agentRunnerDrivingRequestID`（观测优先，fallback=active requestId）+ `currentRequestID` 闭包；非 awaiting 分支 `fail(currentRequestID(), runErr)`；观测在 `preManifestAgentRunnerPort.DispatchAgentRunner` 内、inner 派发前。`runtime_multitask_driver_test.go:394` `TestAgentRunnerDriverStopsTheLaunchOfTheTaskThatJustFailed`：failByTask 注入 task 2 → 断言 receipts 仅见 `request-two`、stopper 仅 `launch-request-two`、原错误保留 ⇒ 真实卡住「task 2 失败 ⇒ 止损 task 2」（实测 PASS） | 无 |
| H2 | ✅ 属实 | 三夹具 `runtime_multitask_driver_test.go` `TestAgentRunnerDriverTwoTaskSerialAssembly` 均直调生产 `ExecuteAgentRunnerRun`；非导出 `testStub` 仅换 chain 端口，plan/record/marker/replay store/engine/event log/串行驱动/停机屏障全生产。`proveTaskStopped` 三态全行为覆盖：proven stop→COMPLETED 双 PASS；live→`ErrInvalidConfig`+"still alive" 且仅 1 次 dispatch；missing→`adapters.ErrInvalidAgentRunnerEvidence`+"cannot be proven gone" 且仅 1 次 dispatch（实测 PASS） | 无 |
| M1 | ✅ 属实 | `agentRunnerAssemblyRecords` 内 `record.Request.Attempt != 1` 拒绝（Step 1，先于 marker Step 2 / engine Step 8，注释明示 red baseline 语义）。`runtime_multitask_driver_test.go:438` 断言空摘要（`ChainState==""` 且计数 0）+ marker 不存在（`os.Stat` ErrNotExist）（实测 PASS） | 无 |
| M2 | ✅ 属实 | `runtime_multitask_supervised_test.go:141-171` 以生产同源 `newSupervisedRuntime`+`deriveSupervisedFacts`+`supervisedFactsList` 重推并按位比较 5 元组；位序与 chain `supervisedFactsEvidence` 逐字一致（`supervised_blackbox.go:84` ↔ `runtime_supervised.go:1937`）。`newSupervisedRuntime` 对照 HEAD 抽取纯函数：构造顺序（snapshotStore+runtime 先于 eventStore/engine）与字段集同源、无引擎依赖、测试侧重推为只读（实测 PASS） | 无 |
| I5 | ✅ 属实 | `runtime_agent.go` dispatched 字段组 `PREMISE (pinned by I5)` 注释 + `awaitingDispatch()` 引用同注 | 无 |

**L1 延后裁定确认**：可接受，不新增风险面——P50 交付物行已含「文档同步（`ARCHITECTURE`/`DEV_PLAN` 历史欠账补齐）」，⑧ 步骤覆盖该回写；注记纯文档项，不触及语义。
**L2 裁定记录确认**：`REMAINING_SLICES.md` P50「裁定记录」行与 `_EN`「Rulings recorded」行均存在且内容属实（方案 A、2026-11-03、范围限装配层+离线 harness+测试改写、CLI 后置、②a–d 串行派单），与 ②a–d 回执和实际改动一致。

## §3 新回归检查结论

- 静态/全量全绿；fail-site 不变量与 red-baseline 串机械核验通过；编码合规。
- 新 seam 最小且无泄漏：`testStub`/`agentRunnerAssemblyStub`/`agentRunnerDispatchWaiter`/`agentRunnerLaunchReceiptReader/Stopper` 全为非导出、console-local，生产载体为既有 adapters 类型（接口满足经编译证明）；`internal/chain` 零改动；adapters 唯一新增 `AgentRunnerStepIntentRegistry` 为纯增量（逐记录复用单记录构造器，不改变既有路径，I4 后置接用不变）。
- H1 修复语义封闭：无观测回退与旧行为等价（未达 dispatch 的失败仍回退 active），不引入新错配；止损为 best-effort 且永不掩盖原错误。

## §4 总判定

H1/H2/M1/M2/I5 五项 CLOSED 声明全部属实，L1 延后登记有效，L2 裁定记录两文件属实；修复未引入新回归。唯一环境备注为本地 1.27 vs CI 1.22 工具链差（FullChain 写路径与 race 由 CI 腿覆盖，⑧c/推送后以 CI 为准）。

`SEVERITY=H0/M0/L0/INFO2`
`PRE-REVIEW r2: PASS`

> INFO-1：本地自跑工具链为 Go 1.27.0，仓库门禁为 1.22；本报告读数需以 CI 1.22 腿为最终判据。INFO-2：② 修复轮回执未单独落盘证据包（会话内交付），本 r2 按硬约束以磁盘实文件为唯一锚点逐项核验，结论不受影响。
