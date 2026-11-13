# 批量辅助 ②（GLM-5.3-Flash）AgentRunner 事件写入点与字段枚举

派单载体 / 模型：`deep-reasoner` + `GLM-5.3-Flash (glm)`（§2.2 批量辅助首选）
任务类型：枚举（§2.2 四类限定之①）；不含判断与择优。

## A. 事件类型闭集

`internal/adapters/agent_runner_event.go:19-35`：

```go
type AgentRunnerEvent struct {
	Kind        string   `json:"kind"`
	EventID     string   `json:"eventId"`
	RequestID   string   `json:"requestId"`
	RequestHash string   `json:"requestHash"`
	RunID       string   `json:"runId"`
	TaskID      string   `json:"taskId"`
	StepID      string   `json:"stepId"`
	Attempt     int      `json:"attempt"`
	AdapterID   string   `json:"adapterId"`
	SessionID   string   `json:"sessionId"`
	Sequence    int      `json:"sequence"`
	OccurredAt  string   `json:"occurredAt"`
	EventType   string   `json:"eventType"`
	PayloadHash string   `json:"payloadHash"`
	Evidence    []string `json:"evidence"`
}
```

`EventType` 闭集（`agent_runner_event.go:168-176` 的 switch）：`session-created` / `output` / `tool` / `usage` / `diagnostic` / `operator-action-required` / `process`；default 返回 `invalid event type`。

## B. 写入链路

**关键事实**：prfrail 生产 Go 代码**不写** `events.jsonl`；写入方为被 spawn 的外部 pinned CLI。库内唯一写入实现是确定性 stub `tools/agent-stub/main.go`。

- 常量：`internal/adapters/agent_runner_evidence.go:32` `agentRunnerEventsFileName = "events.jsonl"`；`tools/agent-stub/main.go:37` 同名常量
- 写入函数 `tools/agent-stub/main.go:218-232` `emitEvents`：一次性 `os.WriteFile(dir/events.jsonl, body, 0o600)`，body 预置两行（`session-created` / `output`）
- 追加函数 `tools/agent-stub/main.go:235-246` `appendOperatorAsk`：`OpenFile(O_APPEND|O_WRONLY|O_CREATE)` + `WriteString`（仅 `mode == "ask"`）
- **Sync/Flush：两个写入函数均无**
- 调用点 `tools/agent-stub/main.go:147-152`：`emitArtifacts(); emitEvents(); if ask { appendOperatorAsk() }; emitMarker()` —— 位于子 CLI 运行末段
- prfrail 侧只读：`agent_runner_evidence.go:263` `readBoundedAgentRunnerArtifact(filepath.Join(spec.EvidenceDir, agentRunnerEventsFileName), ...)`
- 跳位：写入发生在被 spawn 的子 CLI 进程运行期间；prfrail 读取在**结束跳**（`agent_runner_run.go:56` Run → `:90` StartAgentRunnerProcess → `:94` WaitAgentRunnerProcess → `:120` RunEvidenceDir → `:124` CollectAgentRunnerEvidence）

## C. 与 console 的关系

- `internal/console/` 中无直接读写 `events.jsonl` 的代码；`events.jsonl` 字样在 console 仅 `runtime.go:90` 的 `state-events.jsonl`（不同文件）
- console 内无 adapters 事件流解析代码（`operator-action-required` / `CollectAgentRunnerEvidence` 生产命中 0；后者唯一生产调用点 `agent_runner_run.go:124`）
- `ExecuteAgentRunnerRun` 主循环 `runtime_agent.go:415-459`：`for { runErr := engine.Run(ctx); ... }`；`engine.Run` 返回 `chain.ErrAwaitingAgentRunnerTerminal`（`:423`）后经 `:434` `WaitAgentRunnerDispatched` 阻塞等待子 CLI；其生产实现 `AgentRunnerPinnedCLIRun`（`:390`）

## D. requestId 与目录

- `internal/adapters/agent_runner_replay_root.go:17` `agentRunnerReplayRootDirectoryName = "agent-runner-replay"`；`:68` `replayRoot := filepath.Join(cleanRoot, agentRunnerReplayRootDirectoryName)`
- `:25` `agentRunnerReplayEvidenceDirectoryName = "runs"`；`:300-301` `base := filepath.Join(store.root, ...); directory := filepath.Join(base, requestID)`（函数 `RunEvidenceDir` 起于 `:293`）
- 时间点：`cli.go:573` `RunRoot: runDir`；`cli.go:514` 解码 request record 后 requestId 已知；replay store 于 `runtime_agent.go:204` 派生；完整 `runs/<requestId>` 首次构造于等待跳（`agent_runner_run.go:90` → `agent_runner_pinned_cli.go:163`）

（实得字符数：约 3600）
