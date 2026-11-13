# 01-arch · design-r1（P70 S1-STREAM-PROGRESS）

派单：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（① 架构/设计，只读）
输入包：`02-batch/batch-01-state-events-enumeration.md`、`02-batch/batch-02-agentrunner-events-enumeration.md`、`00-briefing.md`、`internal/console/**` 源码、`docs/CONTRACTS.md` 原句、P68 前提。

## 1. 现状核对（仅事实）

- 入口与 flag 解析：`internal/console/cli.go:275` `executeRun`；flag 定义 :277–:287（`jsonOutput` :280、`probePermit` :287），`set.Parse` :288；三条路由 :402（supervised）、:405（agent）、:407（noop）；`const runUsage` :414（不含 `--progress`）。
- 结束摘要 `writeRunSummary` `cli.go:678`，文本分支 6 行：`run: completed` / `runId: %s` / `runDir: %s` / `chainState: %s` / `events: %d` / `eventLog: %s`（:682–:687）。
- `runEventLogPath`：`runtime.go:89`，`filepath.Join(runDir, "events", "state-events.jsonl")`。
- 事件写入点：`chain/events.go:34-53` `FileEventStore.Append`（:42 O_APPEND 打开、:47 Write、:50 `file.Sync()`、:46 defer Close）；`chain/state.go:61-84` `stateWriter.append`（:83 `writer.store.Append`）；`chain/engine.go:584` `transition` → `state.append`，`Run` :122，生产写入仅此一路（引擎 `New` :77 首事件 CREATED，:85 `engine := &Engine{... state: stateWriter{store: options.Events...}}`）。
- AgentRunner 事件：写入方是外部 pinned CLI（库内仅 stub：`tools/agent-stub/main.go:218-232` `emitEvents` 一次性 `WriteFile`、:235-246 `appendOperatorAsk`，**均无 Sync**）；prfrail 只在结束跳读 `adapters/agent_runner_evidence.go:263` ← `agent_runner_run.go:124`；常量 `agentRunnerEventsFileName="events.jsonl"` :32；目录 `agent_runner_replay_root.go:17/25/68`（`RunEvidenceDir` :293 起）。
- `--json`：`cli.go:793` `writeJSONResponse`（`json.NewEncoder` + SetIndent，写 stdout）。
- P68 前提：`remaining-rulings.md:86`——`run` 跑完才输出、无 per-task/per-step 行、无 `--verbose` / `--progress`；`state-events.jsonl` 为唯一全序源但**无订阅 API**。

## 2. 数据源与可读性（静态代码依据）

- `state-events.jsonl`：**可被运行中读取**。依据：`events.go:42` 追加打开（非独占）、:47 每事件 Write 后 :50 `file.Sync()`、:46 每次 defer Close；`cli.go:366-370` run 前置仅要求文件「不存在」。`CONTRACTS.md:257` 原句：「事件日志仍是该 runRoot 唯一 `events/state-events.jsonl` 全序源」。
- AgentRunner `events.jsonl`：**不可依赖运行中读取**。写入方是 spawn 的外部子 CLI，stub 写时点在子进程末段且无 Sync；prfrail 生产代码仅在结束跳读（`agent_runner_evidence.go:263`）。

## 3. 实现形态判定：**形态 A**

装饰器实现于 `internal/console/**`：包装 console 三处构造的 store（`runtime.go:40`、`runtime_agent.go:345`、`runtime_supervised.go:620`），`Append` 转发内层成功后向 stderr 打印一行。引擎所有生产事件都经 `engine.go:584 → state.go:83 → Append`，同步打印天然按事件序、无重复、无 goroutine、无 poll 成本。

**SealSink 风险显式处理**：装饰器必须同时实现 `chain.SealSink`（`seal.go:66-71`：`SealTaskEvidence` / `LoadTaskSeal`）并转发内层。若不转发：`engine.go:340` 断言失败后静默返回 `(false, nil)`（封存重放关闭）；`engine.go:374` 断言失败返回 `ErrSealUnavailable` ⇒ 每个接受闭包 fail-closed；`seal_verify.go:25` 断言失败静默跳过恢复核验。`FileEventStore` 已实现 SealSink（`seal.go:117/173`），转发为纯委托。

否决：B（轮询）引入 goroutine 生命周期、ctx 取消清理、poll 成本与末行顺序竞态；C（仅 AgentRunner 事件）运行中可达性差且看不到 task/step 起止。

## 4. 与 `--json` 关系：**共存**

`--progress` 只写 **stderr**；`writeRunSummary` / `writeJSONResponse` 只写 **stdout** ⇒ stdout 字节流逐字节不变。无用法错误路径（反例「互斥返回用法错误」不适用）。

## 5. 反例清单

- **R1 默认关闭输出不变**：同一 definition 固定 runId、不同 runDir 两次运行（一次带 `--progress` 一次不带），stdout 逐字节相等且 stderr 均空。变红开关：删除「progress 为 nil 时不构造装饰器」守卫。
- **R2 `--progress` + `--json` 最终输出不变**：`--json` 与 `--json --progress` 两次运行 stdout 逐字节相等（后者 stderr 非空）。变红开关：把 progress 行写入 stdout。
- **R3 打印失败降级**：传入恒失败的 stderr writer，`ExecuteNoopRun` 返回 nil 且事件文件与无 progress 时逐字节相同。变红开关：删除「打印错误忽略」守卫。
- **R4 装饰器不改事件字节**：带/不带 progress 两次运行产出的事件文件逐字节相等。变红开关：装饰器在转发前重 marshal event。

## 6. 打印粒度可达性

- 逐 task 起止：**可达**（`task-created` engine.go:173、`task-precheck` :179、`steps-started` :204、`task-accepted` :329/:363、failTask / failChain）。
- 逐 step 起止：**可达**（`step-created` :410、`step-started` :427、`step-passed` :488、`step-failed` :467、`NOOP_RECORDED` :416）。
- 门禁命令起止：**部分可达＝已知边界**。门禁以 hook/gate 形态入任务步骤表时由 step 事件覆盖；作为 port 内命令（baseline/acceptance/review/scan）执行的门禁无起止事件。
- 门禁退出码：**不可达＝已知边界**。`StateEvent` 字段闭集无 exitCode；`reason.code` 是语义码非进程退出码。
- AI 工具调用摘要行：**不可达＝已知边界**（外部子 CLI 批写、prfrail 仅末跳读取）。

补法均需改 chain / adapters / schema ⇒ 属停升级边界外的能力，本片不实施，记为已知边界。

## 7. 改动面清单（仅 `internal/console/**`）

- `cli.go`：`executeRun` :275 新增 `progressOutput := set.Bool("progress", false, ...)`；:402/:405/:407 透传；`runUsage` :414 加 `[--progress]`；`writeRunSummary` :678 不动。
- 新增 `internal/console/progress.go`：`progressStore`（Append / Load + SealSink 两方法转发）、`formatProgressLine`、`wrapProgressStore(store, stderr, definition)`（nil 返回原 store）。
- `runtime.go`：`ExecuteNoopRun` :33 加 `progress io.Writer` 参数，:34 store 包裹后注入 :40。
- `runtime_agent.go`：`AgentRunnerAssemblyOptions` 加 `Progress io.Writer`（零值兼容既有测试）；:341 eventStore 包裹后注入 :345。
- `runtime_supervised.go`：`SupervisedAssemblyOptions` 加 `Progress`；:616 包裹后注入 :620；:904（只读重放）不包裹。
- `supervised_cmd.go`：`executeSupervisedRun` :34 透传 Progress。

## 8. 测试点清单

- 受影响：`cli_test.go`、`runtime_zero_model_e2e_test.go:76/115`、`stopper_test.go:24/59`（`ExecuteNoopRun` 签名的 4 处测试调用点机械更新）；AgentRunner / Supervised 走 options 字段加法，零改动。
- 新增：`internal/console/progress_test.go`（R1–R4 + SealSink 转发断言 + 打印格式快照 + `--progress` 解析）。

## 9. 停升级边界自检

① 需改 chain/adapters 事件订阅机制——**未触发**；② 两源均不可读且 console 无法解决——**未触发**（`state-events.jsonl` 每事件 Sync，形态 A 直接消费 Append 数据）；③ 需改 CONTRACTS 语义——**未触发**；④ 需改 schema/fixtures/workflows/go.mod——**未触发**；⑤ `--progress` 与 `--json` 无法调和——**未触发**（stderr/stdout 分离）。

DESIGN: DONE
