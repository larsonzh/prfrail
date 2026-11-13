# 批量辅助 ①（GLM-5.3-Flash）state-events.jsonl 写入点与字段枚举

派单载体 / 模型：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 批量辅助首选；非流水线角色）
任务类型：枚举（§2.2 四类限定之①）；不含判断与择优。

## A. 路径定义

`internal/console/runtime.go:89-91`：

```go
func runEventLogPath(runDir string) string {
	return filepath.Join(runDir, "events", "state-events.jsonl")
}
```

非测试调用点：`runtime.go:34`（ExecuteNoopRun）、`runtime.go:63`、`runtime.go:67`（LoadRunSummary）、`cli.go:366`、`runtime_agent.go:198`、`runtime_supervised.go:305/343/616/810/904`。

## B. 写入链路

- `runtime.go:34`：`store, err := chain.NewFileEventStore(runEventLogPath(runDir))`
- `runtime_agent.go:341`：`eventStore, err := chain.NewFileEventStore(eventLog)`（`:198` 定义 `eventLog := runEventLogPath(options.RunRoot)`）
- `runtime_supervised.go:616`：`eventStore, err := chain.NewFileEventStore(runEventLogPath(options.RunRoot))`；`:904`：`store, err := chain.NewFileEventStore(runEventLogPath(runtime.runRoot))`
- `internal/chain/events.go:34-53` `FileEventStore.Append`：O_CREATE|O_APPEND|O_WRONLY → `file.Write(append(data, '\n'))` → **`file.Sync()`** → return
- `internal/chain/state.go:61-84` `stateWriter.append`：构造 `evidence.NewStateEvent` → `writer.store.Append(ctx, record)` → `writer.apply(record)`
- `internal/chain/engine.go:24` `Events EventStore`；`:85` `engine := &Engine{options: options, state: stateWriter{store: options.Events, ...}}`

## C. 事件类型闭集

`internal/evidence/event.go:58-62`：

```go
type StateEvent struct {
	SchemaVersion string `json:"schemaVersion"`
	Event         Event  `json:"event"`
	EventHash     string `json:"eventHash"`
}
```

- `Entity.Kind` 闭集（`event.go:205-220` 的 `validateEntity`）：`chain` / `task` / `step`，default 返回 `invalid entity kind`
- `Reason.Code`：无取值闭集，仅 `idPattern` 语法约束（`event.go:18`，`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`）
- JSON 字段名：`schemaVersion` / `event` / `eventHash`；`eventId` / `runId` / `sequence` / `occurredAt` / `actor` / `entity` / `fromState` / `toState` / `previousEventHash` / `inputEvidence` / `reason`；`type` / `id`；`kind` / `runId` / `taskId,omitempty` / `stepId,omitempty` / `attempt,omitempty`；`code` / `message,omitempty`

## D. 写入点的语义分类

生产写入点仅 1 处：`engine.go:585`（`Engine.transition` → `engine.state.append(...)`）。

`transition(` 调用点按所属函数与所在跳位（调用链锚定 `Run` 起于 `engine.go:122`）：

- 构造期（非 Run）：`engine.go:93`（`New`）
- Run 前段 / 循环外：`engine.go:115`（Initialize ← Run:136）、`:132`（Run:140 PAUSED 分支）、`:145`（BASELINED→RUNNING，进任务循环前）
- 任务循环内（Run:149 `for _, task := range engine.options.Definition.Tasks`）：`162`→runTask `173/179/189/204/221/329`；runTask:212→runStep `410/416/427/434/440/467/474/483/488`；runTask:217→runTaskPostflight `postflight.go:410/414`（含 runSupervisedTaskPostflight `:451/455`）；runTask:230→acceptFromExistingSeal `engine.go:363`；markTaskForRepair `389/393`；failTask `568`；boundary `519`；cancelRun `532/544/552/558`；failChain `577`
- 结束时（Run 尾）：`engine.go:166`（COMPLETED）
- 非 Run 直达（console 驱动）：`SubmitAgentRunnerTerminal`（`agent_runner_terminal_route.go:85/92`）与 convergeAgentTerminalRoute（`:148/155/161/169/173`），调用点 `runtime_agent.go:445`（engine.Run 于 `:422`）；`OpenOperatorInteraction`（`operator_interaction.go:351/376/381`）；`ResumeOperatorInteraction`（`:431/440`）；`OpenSupervisedWorkWindow`（`supervised_blackbox.go:250/253`）；`SubmitSupervisedReturn`（`:319`）与 convergeSupervisedReturn（`:331/336`）；`pauseUncertain`（`recover.go:37`）

## E. 已存在的调试/输出设施

- `internal/console/` 生产文件内 `progress|verbose|watch|ticker|heartbeat` 命中 **0**
- `chain.Options` 的装饰器/包装类型：**未找到**；直接字面量构造仅 `runtime.go:40`、`runtime_agent.go:345`、`runtime_supervised.go:620`

（实得字符数：约 4300）
