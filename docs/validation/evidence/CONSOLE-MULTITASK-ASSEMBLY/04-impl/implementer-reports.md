# P50 ② 实施报告合集（②a–②d，逐字归档）

> 载体：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`；方案 A 拆分派单（用户 2026-11-03 裁定）。以下为各步回执原文。

## ②a 协议先行（仅 CONTRACTS 对）

**插入位置**：CN `docs/CONTRACTS.md:255`（A6 段末尾其后），前置句为 A6 段末「…本片全程离线：不得连接真实候选、不得访问网络，本段结论不得引为候选兼容性证明。」；EN `docs/CONTRACTS_EN.md:249`，前置句 “…nothing here may be cited as candidate-compatibility proof.”。设计草案 blockquote 已按 A6 普通正文风格去 `>` 呈现，CN 与草案逐字一致（程序比对 `true`）。

**CN 新增段（逐字）**
**多任务串行装配（P50）**：runRoot 级身份镜像与启动槽是「**同一时刻至多一个在跑 launch**」的排他点，不构成 run 生命周期内的一次性唯一。同一 `runId` 的一个 chain 可含多个 task（每 task 一个 isolated-workspace code step）；console 装配层为每个 `(taskId, stepId, attempt)` 绑定各自已持久化的 `agent-runner-request`（各 `requestId` 互异，R/C/终局/工件按 requestId 隔离），并由 chain engine **串行**调度：前一 task 的 launch 必须已证明停机（槽位释放、镜像所指进程死亡或镜像已由新 launch 原子轮换）后，后一 task 才允许发布 R 与 launch-intent；镜像轮换=以新 launch 身份原子覆盖旧镜像，与既有「镜像命名死亡进程时为陈旧、可被替换」语义一致。同一 runRoot 内禁止并发在跑的 launch；并发 Start 败者依旧不得创建日志工件、不得覆盖身份镜像。事件日志仍是该 runRoot 唯一 `events/state-events.jsonl` 全序源；各 task 六枚冻结事实按各自终局路由事件的位置化证据（A5 `evidence[2..7]`）归属，不得跨 task 对账。

**EN 新增段（逐字）**
**Multi-task serial assembly (P50)**: the runRoot-level identity mirror and the launch slot are the exclusion point for "**at most one running launch at any instant**", not a one-time uniqueness over the run lifetime. One chain under the same `runId` may contain multiple tasks (one isolated-workspace code step per task); the console assembly layer binds each `(taskId, stepId, attempt)` to its own already-persisted `agent-runner-request` (distinct `requestId` values, with R/C/terminals/artifacts isolated by requestId), and the chain engine schedules them **serially**: only after the previous task's launch has been provably stopped (slot released, the process named by the mirror is dead, or the mirror already atomically rotated by the new launch) may the next task publish R and launch-intent; mirror rotation means atomically overwriting the old mirror with the new launch identity, consistent with the existing semantics that "a process named by the mirror that is dead is stale and may be replaced". Concurrent running launches within one runRoot are forbidden; a loser of a concurrent Start still must not create log artifacts or overwrite the identity mirror. The event log remains the single total-order source `events/state-events.jsonl` for that runRoot; each task's six frozen facts attach via the positional evidence of that task's own terminal routing event (A5 `evidence[2..7]`) and must not be reconciled across tasks.

**git diff --numstat**
```
1	0	docs/CONTRACTS.md
1	0	docs/CONTRACTS_EN.md
```
（各 +1/−0；node 复核 added=1/deleted=0；两文件 BOM 保留、纯 LF；加粗标记各 3 对对称；段落间距与既有正文一致。工作区其余脏项为既有、非本角色改动。）

IMPLEMENT: DONE
IMPL-2A: DONE

## ②b 装配层 agent（仅 runtime_agent.go）

**新增/改名（均 `internal/console/runtime_agent.go`）**

- `:964` `func agentRunnerAssemblyPlan(definition chain.Definition) ([]agentRunnerTaskBinding, error)`
  0 task → `ErrInvalidConfig: …requires at least one task`（专用文案）；逐 task：`len(Steps)!=1` → 原 `requires exactly one step`；形状门 `Kind!="code" || Mode!=IsolatedWorkspace` → 原 `requires one isolated-workspace code step`。
- `:988` `agentRunnerAssemblyRecords(options, plan) ([]adapters.AgentRunnerRequestRecord, error)`
  `RequestRecords` 非空 → 按 `taskID/stepID` 查表（未命中即 `request record is missing or invalid`）；否则仅当 `len(plan)==1` 走 `RequestRecord` 兼容路径，多 task 无注册表同报该文案（红测翻转预期）；逐记录 `ValidateAgentRunnerRequestRecord` + RunID/TaskID/StepID 绑定 + `requestId` 互异。
- `:1023` `agentRunnerBindTasks(options, plan, records, store, launcher, snapshotStore, snapshotRoot, clock)`
  每 task 一套：`workspaceRoot` + `RejectWriterRoots` + `SnapshotCapturer` + `StepIntentPreparer` + `AgentRunnerPostflight` + `ReplayDispatcher` + `CompositeAdmission` + `TerminalPublisher`。
- `:664` `agentRunnerRequestRecordKey(taskID, stepID) = taskID+"/"+stepID`（稳定编码，公开给注册表构造方）。
- 路由多路复用（console-local，不改 adapters/chain）：`:1124` `agentRunnerDispatchMux`、`agentRunnerAdmissionMux`（(taskID,stepID) 键）、`agentRunnerIntentMux`（(taskID,stepID) 键，因 adapters preparer 对非匹配 step 返回零值故必须在外部选路）、`:1174` `agentRunnerPostflightMux`（(taskID,attempt) 键）。
- `:390` `agentRunnerRuntime`：`requestRecords/workspaceRoots/capturers` + `beforeManifests/acceptedEntries map[string]…`（键 `taskID#attempt`）+ 派发观测 `dispatchedTaskID/StepID/Attempt`；`managedProcessesGone(taskID,attempt)`、`:828` `proveTaskStopped(taskID,attempt)`。

**驱动循环（`:300`–`:352`）**

```
active = tasks[0]
loop:
  err = engine.Run(ctx)
  err == nil            → return summary (chain COMPLETED)
  !errors.Is(awaiting)  → fail(active.requestId, err)
  task = agentRunnerAwaitingTask(tasks, runtime)   // 用派发端口观测，非定义序假设
  wait  = WaitAgentRunnerDispatched(record, process, workspace, beforeManifests[key])
  task.publisher.PublishTerminal(outcome) → LoadChainAgentRunnerTerminal(store, record) → engine.SubmitAgentRunnerTerminal
  status != completed → fail
  task != tasks[last] → proveTaskStopped(prev) 失败即 fail-closed
```

**`git diff --numstat`**

```
568  173  internal/console/runtime_agent.go
```
（`git status` 仅此一文件新增改动；`docs/CONTRACTS*.md` 为 ②a 已暂存，`runtime_multitask_red_test.go` 为 ③ 新增未跟踪，本步均未触碰。`gofmt -l` 空输出。）

**与契约字面逐条对应**

| 契约 A6「多任务串行装配（P50）」字面 | 落地 |
|---|---|
| 同一 chain 可含多个 task（每 task 一个 isolated-workspace code step） | `agentRunnerAssemblyPlan` 逐 task 校验并产出绑定表 |
| 每 `(taskId, stepId, attempt)` 绑定各自已持久化的 request | `agentRunnerBindTasks` + `recordFor(taskID,attempt)`；intent/postflight/publisher 全部 record 级 |
| 各 `requestId` 互异，R/C/终局/工件按 requestId 隔离 | `agentRunnerAssemblyRecords` 强制互异；`replayEvidenceDir(taskID,attempt)` 经 record 的 requestId |
| 前一 launch 必须已证明停机后，后一 task 才可发布 R 与 launch-intent | 循环内 `proveTaskStopped`（复用 `managedProcessesGone`，证据缺口亦 fail-closed）在下一 `engine.Run` 之前执行 |
| 同一 runRoot 内禁止并发在跑 launch | 驱动为单线程串行，无并发 spawn、无新锁 |
| 事件日志仍是唯一全序源；六枚事实按各 task 终局路由位置化证据归属 | 未改 chain；postflight 走 `agentRunnerPostflightMux`，按 task 各自 requestId 证据目录 |
| 禁止并发 Start 败者创建日志/覆盖镜像 | 未触碰 launcher/mirror/slot（`internal/adapters` 零改动） |

**`go build ./...`**：`BUILD_EXIT=0`（无任何输出）

**偏差点名**
1. **`return fail(` 由 7 处变 8 处**：新增失败点是「前一 task 停机证明失败」（`:348`）。源码级不变量 `TestAgentRunnerAssemblyFunnelsEveryPostSpawnFailureThroughTheFallback` 的 `wantFailureReturnPoints=7` 需在 ②d 改测试时同步 bump 到 8；「单条 naked failure return 且只在 `fail` 闭包内」形态未变。
2. **`fail` 闭包签名加 `requestID` 形参**：多 task 下只能对「当前正在驱动的 task」做 try-stop；单 task 时取值与旧实现的 `options.RequestRecord.Request.RequestID` 完全相同。
3. **`Chain Baseline` 仍为 run 级**：`CaptureBaseline` 沿用首 task 的 record/workspace（多 task 时第二 task 从链的 accepted parent 物化）；契约未定义 run 级基线的多 task 口径，未自行扩展。
4. **等待 task 用「派发端口观测」而非定义序索引**：可自愈 PASSED 跳过；但若 resume 场景下某 step 已处于 `TERMINAL_PENDING` 而本进程未派发（无观测），会 fail-closed 报 `cannot bind the awaiting step to a dispatch`，不猜测。
5. **`WorkspaceRoot` 覆盖仅在单 task 计划生效**（测试强制 workspace/replay 重合用），多 task 忽略。
6. 既有 `internal/console` 测试因 `agentRunnerRuntime` 字段映射化与 `agentRunnerAssemblyStep` 移除而**编译失败属预期**（②d 统一改写）；本步按指令未跑 `go vet`/`go test`，未触碰任何测试文件。

IMPL-2B: DONE

## ②c 装配层 supervised（仅 runtime_supervised.go）

**装配计划**
- `type supervisedTaskBinding struct{ TaskID, StepID string; Policy chain.HandoffPolicy }`
- `supervisedAssemblyPlan(options SupervisedAssemblyOptions) ([]supervisedTaskBinding, error)`：`len(Tasks)==0` → `"the supervised assembly requires at least one task"`（专用文案，`ErrInvalidConfig`）；逐 task 校验 step 数=1 → 原 shape 文案 → `HandoffPolicy` 非空 → 原 policy 文案；`TaskPolicies` 非空时按 `supervisedTaskRecordKey(taskID,stepID)`（`"<taskID>/<stepID>"`，与 `agentRunnerRequestRecordKey` 同风格）查表，未登记或与 step 声明策略摘要不一致即拒绝（一致性钉死，非新增授权）。
- `SupervisedAssemblyOptions` 新增 `TaskPolicies map[string]chain.HandoffPolicy`；既有单记录路径（step 声明的 `HandoffPolicy`）在注册表为空时生效 ⇒ `cli.go` 零改动。

**运行时态映射化**
- `type supervisedTaskState{ taskID, stepID, attempt, policy, workspaceRoot, acceptedEntries, acceptedHash, workspaceReady }`，键函数 `supervisedTaskStateKey(taskID, attempt) = "<taskID>#<attempt>"`（对齐 `agentRunnerTaskStateKey`）。
- `supervisedRuntime`：`plan []*supervisedTaskState`（定义序）+ `tasks map[string]*supervisedTaskState` + `current string`（仅 failpoint 需要，因该 seam 不带 task 身份，由 `Materialize`/`Execute` 标记）。
- 解析器：`stateFor(taskID,attempt)`、`stateForStep(taskID,stepID)`、`currentState()`、`parkedTaskState(ledger)`（取 ledger 中**最新**一条绑定 request 的 task ⇒ 停放 task；从未停放且单 task 时回退该唯一 task）。

**串行驱动**
- `ExecuteSupervisedRun`：`plan` → runRoot/visible 校验 → marker → 引擎；`PAUSED` → `resumeSupervisedRun`（内部按 `parkedTaskState` 逐 task 决定 B/C/D/replay 相位，确认按 `(runID,taskID,stepID,attempt)` 匹配，不再以「确认账本为空」判定相位）；否则 `engine.Run`，`errSupervisedParked` 视为合法停放。
- `ExecuteSupervisedReturn`：`loadLedger` → `parkedTaskState` → 该 task `loadSession` → abort（run 级不变）/ `enforceFrozenHandoffPolicy(state, session)` → `completeSupervisedRun(state,...)`；其中尾链 `engine.Run` 若因**下一 task** 再次停放（`errSupervisedParked`）不再报错 ⇒ 前一 task 归还后才进入下一 task。

**逐 task 冻结**：`persistFrozenHandoffPolicy(state)` / `enforceFrozenHandoffPolicy(state, session)` 使用该 task 的冻结文件与策略，成功后仅改写 `state.policy`。

**`git diff --numstat`**：`459	205	internal/console/runtime_supervised.go`（仅此一文件）。

**与契约字面对应**
1. 「同一 runId 一个 chain 可含多 task，每 task 一个 isolated-workspace code step」→ 计划逐 task 生成绑定，每 task 独立 workspace `<runRoot>/workspaces/<taskID>-attempt-1`。
2. 「由 chain engine 串行调度；前一 task 停机后才允许后一 task」→ 复用 engine 的定义序串行循环；本层不再有跨 task 复用，归还后才继续（`completeSupervisedRun` 容忍下一 task 停放）。
3. 「事件日志仍是唯一全序源；各 task 事实按各自终局路由事件归属，不得跨 task 对账」→ `RunSupervisedPostflight` 由 `request.TaskID/Attempt` 取该 task 工件重算；会话/冻结/回执/报告/扫描全部逐 task 落盘，`acceptedEntries/acceptedHash` 逐 task 键控。
4. 「runRoot 级身份镜像与槽位为时刻级排他」→ 本层未触碰镜像/槽位，无新增锁、无并发 spawn。
5. 单 task 兼容：`taskDir` 在 `len(plan)<=1` 时返回原 `supervised/` 平铺目录，`SupervisedReportPath`、`supervised/scans/*` 等既有 store-local 路径与消费者（`supervised_cmd.go`）不变。

**`go build ./...`**：`BUILD_EXIT=0`；`gofmt -l` 无输出；编码 BOM=False CRLF=False。

**偏差点名**
1. **注册表语义**：`TaskPolicies` 取「与 step 声明策略同摘要」的一致性钉死，而非策略供给/覆盖；原因是本角色不得新增策略权威，且红测要求「逐 task `HandoffPolicy` 非空」字面保留。若主控需要真正的策略供给语义，需另立契约条款。
2. **停放 task 的解析**：不使用 chain 内部 `taskKey` 编码（避免越界耦合），改用「交互账本中最新绑定 request」；语义等价于串行引擎的停放顺序，单 task 有唯一回退。
3. **磁盘布局**：多 task 时逐 task 子目录 `supervised/tasks/<taskID>-attempt-<n>/`（session/freeze/receipt/report/scans/delivery），单 task 保持平铺 ⇒ 既有单 task 测试的文件断言与 `SupervisedReportPath` 行为不变；两个 run 级账本（interactions/confirmations）保持平铺，因其记录自带 task/step。
4. **attempt 限定**：装配仅绑定 attempt=1（与既有 session/confirmation 的 `Attempt:1` 一致）；`Materialize` 收到 attempt≠1 会 fail-closed 报 `ErrInvalidConfig`，不猜测重试续接。
5. **abort 语义**：维持 run 级（`RequestCancel` → CANCELLED），任一 task 中止后后续 task 不再启动；未改为 task 级。
6. **已知中间态**：`ErrSupervisedStepRequired` 仍为导出但未使用的旧哨兵（未删除，避免扩大改动面）；`runtime_supervised_test.go` 因 `supervisedAssemblyStep`→`supervisedAssemblyPlan` 与 `newSupervisedEngine` 签名变化将编译失败，属 ②d 范围；`go build ./...` 全绿，未跑 `go vet`/`go test`。

IMPL-2C: DONE

## ②d 适配器与测试转绿

**改动统计（`git diff --numstat`）**

| 文件 | +/− | 性质 |
|---|---|---|
| `internal/adapters/agent_runner_intent.go` | +90 / −0 | 唯一生产文件（多记录 preparer 注册表） |
| `internal/adapters/agent_runner_intent_registry_test.go` | 新文件 | 注册表聚焦测试 |
| `internal/console/runtime_multitask_red_test.go` | 改写（未跟踪→已改写） | 接受形 |
| `internal/console/runtime_multitask_agent_test.go` | 新文件 | 两 task 串行离线 harness |
| `internal/console/runtime_multitask_supervised_test.go` | 新文件 | supervised 串行 |
| `internal/console/runtime_agent_e2e_test.go` | +8 / −7 | 夹具/计数修复 |
| `internal/console/runtime_supervised_test.go` | +14 / −6 | 直调/直构修复 |

`runtime_agent.go`(+568/−173)、`runtime_supervised.go`(+459/−205) 属 ②b/②c 既有改动，本步**未触碰**。

**GLM-1 第 7 组受影响测试逐条处置（全 CLOSED）**
1. `runtime_agent_e2e_test.go:448` FullChain — **CLOSED**：未改语义；Windows 因既有 durability 闸 SKIP（原文 skip 文案保留），Ubuntu CI 承担。
2. `:1012` FunnelsEveryPostSpawnFailure — **CLOSED**：`wantFailureReturnPoints` 7→8，注释新增「awaiting-task 绑定」「任务间停机证明」两处。
3. `:600` IncompleteRecords — **CLOSED**：单记录兼容路径仍有效，未改。
4. `:620` newAgentRunnerPIDReportRuntime — **CLOSED**：`requestRecord` 补 `TaskID=assemblyTaskID, Attempt=1`，使 `recordFor` 命中。
5. `:875` `managedProcessesGone()` — **CLOSED**：改为 `managedProcessesGone(assemblyTaskID, 1)`。
6. `runtime_supervised_test.go:1231` FullChain(ExtendedScans) — **CLOSED**：未改断言，PASS。
7. `:1421` 直调 `supervisedAssemblyStep` — **CLOSED**：改为 `supervisedAssemblyPlan(harness.options)` + `newSupervisedEngine(ctx, opts, plan, clock)`。
8. `:1107` 直构 `&supervisedRuntime{workspaceRoot,policy}` — **CLOSED**：改为 `&supervisedRuntime{}` + `&supervisedTaskState{...}`。
9. `:1134` `runSecretAndTypeScan(ctx, manifest)` — **CLOSED**：补 `state` 形参。
无 `t.Skip` 新增。

**新绿测断言要点**
- `TestAgentRunnerMultiTaskSerialOfflineHarness`：真链引擎 + 真 FileEventStore + 真 replay store + **真 `adapters.NewAgentRunnerStepIntentRegistry(records, def)`**，stub 仅承接 step/accept/review/publish/postflight。断言：①`executed==[task-one/code-one, task-two/code-two]`；②每 task 终局事件 `InputEvidence` 逐位等于 `[requestHash, completionHash, 六事实]`，且不含他 task 事实；③单日志全序 `start1<term1<start2<term2` 且 sequence 严格递增；④requestId/RecordHash/证据目录两两互异；durability 已证时再写 R/C 并断言两侧均为 `TerminalReceiptPresent`。
- `TestSupervisedMultiTaskSerialHarness`：生产入口全流程（park→confirm→open→return→park→…→COMPLETED）。断言：两 task 各 5 枚冻结事实、集合互不相交、`term1<open2<term2`、每 task 工件落在 `supervised/tasks/<taskID>-attempt-1/`（`handoff-receipt.json`+`session.json`）。
- 注册表 3 条：逐记录绑定正确（各 task 得自己 requestId/ContextHash）、跨记录漂移拒绝（同 task/step、同 requestId、跨 run、空集）、逐记录复用 `matches==1` 与 step 形状校验。

**测试结果原文关键行**

```
--- PASS: TestAgentRunnerAssemblyFunnelsEveryPostSpawnFailureThroughTheFallback (0.00s)
--- PASS: TestAgentRunnerMultiTaskSerialOfflineHarness (0.14s)
--- PASS: TestAgentRunnerAssemblyAcceptsMultiTaskDefinitionShape (0.00s)
--- PASS: TestSupervisedAssemblyAcceptsMultiTaskDefinitionShape (0.00s)
--- PASS: TestMultiTaskDefinitionShapeIsOtherwiseValid (0.00s)
--- PASS: TestSupervisedMultiTaskSerialHarness (0.36s)
--- PASS: TestSupervisedAssemblyFullChain (0.15s)
--- PASS: TestSupervisedAssemblyFullChainExtendedScans (0.14s)
--- PASS: TestAgentRunnerStepIntentRegistryBindsEachRecordToItsOwnStep (0.00s)
--- PASS: TestAgentRunnerStepIntentRegistryRejectsCrossRecordDrift (0.00s)
--- PASS: TestNewAgentRunnerStepIntentRegistryReusesPerRecordAdmission (0.00s)
ok  github.com/larsonzh/prfrail/internal/console    2.405s
ok  github.com/larsonzh/prfrail/internal/adapters   1.703s
```
聚焦三包：`ok console 11.119s / ok adapters 40.794s / ok chain 1.529s`。
全量：`ok ... internal/adapters 58.081s / chain 5.061s / console 22.111s`，其余包全部 `ok`（无 FAIL；`agent-stub 9.004s`）。
`gofmt -l`（改动 7 文件）无输出；`go build ./...`=0；`go vet` 三包 =0。
翻转机械证据：`requires exactly one task` 在两生产文件中已 0 处，仅剩 `requires at least one task`（`runtime_agent.go:966`、`runtime_supervised.go:380`）。

**偏差与限制（点名）**
1. **红→绿无法在本机原地重放**：②b/②c 未提交，无干净 HEAD 可回退到「拒绝存在」形态；翻转证据 = `03-tests/red-baseline.md` §B（改造前 3 条 PASS）＋本次接受形全绿＋旧文案 grep 清零。
2. **Windows 不跑 R/C 写路径**：`RecordRequest` 受 `PublishDurability` 门控（Windows 恒 unproven）。绿测在该平台只断言 per-requestId 路径/哈希互异并 `t.Logf`，**不 skip**；写路径由 Ubuntu CI 承担。
3. supervised 绿测用生产入口而非设计所述 stub 端口（等价更强），未改生产代码。
4. 注册表已实现并单测，但 `runtime_agent.go` 的 `agentRunnerIntentMux` 仍按记录逐个构造 preparer（该文件不在本步授权范围）；注册表供其后续选用，行为等价。
5. `runtime_multitask_red_test.go` 与两条 `runtime_multitask_{agent,supervised}_test.go` 为未跟踪新文件（未 `git add`）。

IMPL-2D: DONE
