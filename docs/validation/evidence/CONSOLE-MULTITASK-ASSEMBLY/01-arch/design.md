# P50 `CONSOLE-MULTITASK-ASSEMBLY` — ① 架构师设计文档（01-arch/design.md）

## §1 结论摘要

- **破局裁决**：采用 **(c) 强化形态（(c)+(d) 组合）**——runRoot 边界不变、A6 身份镜像与启动槽机制**零改动**，多 task 由 chain engine **串行**调度；每 `(taskID, stepID, attempt)` 绑定各自已持久化的 `agent-runner-request`（`requestId` 互异）。拒绝 (a)（改镜像按 taskID 隔离，破坏 `console/stopper.go:47` 唯一读取点与 CONTRACTS.md:254「镜像按 run root 唯一」字面）与 (b)（按 runID 重划 runRoot，破坏 replay store 构造前提：CONTRACTS.md 要求 runRoot 含 `events/state-events.jsonl` durable 标记、store 由 `(runRoot, runID)` 唯一推导，且 chain 每 run 仅一个事件日志）。
- **协议先行判定：必须修改 `docs/CONTRACTS.md`**（§2 给出逐段草案，`_EN` 同步）。不改任何 wire 记录、Schema、`testdata/contracts/` 夹具（镜像/槽位/starting-marker 均为 store-local，不进 wire、不参与 recordHash，CONTRACTS.md 已钉；`schemas/` 无对应文件，夹具扫描确认零引用）。
- **红测**：两条（AgentRunner + Supervised 各一），卡 `runtime_agent.go:772` 与 `runtime_supervised.go:333` 的精确拒绝文案（§3）。
- **流分隔键**：wire/磁盘层 `requestId`；装配层 `(runID, taskID, stepID, attempt)`；chain 层 `runID`（单一事件日志全序）。Facts 不串号靠 A5 位置化路由证据 + 事件日志全序 + requestId 绑定（§4）。
- 风险定级 L 档成立：动 console 装配主路径与 adapters intent 层；不动 DR-9/DR13 敏感面（slot/mirror/恢复实现）。

## §2 命题一：破局方案与契约修订草案

**逐候选裁决**：
- **(a) 镜像按 taskID 隔离**：否。镜像文件名与位置是 A6 与 `stopper.go` 的既有契约（`managed-process.identity.json` 于 runRoot 根），operator 停机入口无 taskID 语义；且槽位仍是 runRoot 级单例，并行仍冲突。改造面最大、收益最小。
- **(b) 每 task 独立 runRoot**：否。与 replay store 派生、event-log durable 标记、`NewAgentRunnerReplayStoreForRun(options.RunRoot, options.RunID)`（`runtime_agent.go:152`）及 `resolveApprovalsRunDir(cwd, runID, hint)`（`stopper.go:92`）全部冲突；会制造 N 个事件日志，撕裂 chain 投影。
- **(c) runRoot 唯一 + 串行调度**：**是**。A6 字面「同一 run root **同时**只允许一个在跑的 launch」（CONTRACTS.md:254）本即**时刻级**排他；`refuseLiveMirrorProcess`（`agent_runner_pinned_cli.go:587`）已实现「镜像命名死亡进程=陈旧、可被替换」，任务间镜像轮换在现行代码内已合法。多 task 串行 = 逐 task 重复既有「槽位→R→intent→镜像→spawn→C→终局路由」五步，**零新锁**。
- **(d) 其他**：即 (c) + 装配层每 task 请求记录注册表 + 每 task 运行时态映射化。这是 (c) 的落地形态，非独立方案。

**为何仍必须改 CONTRACTS.md（不能只靠字面论证）**：P45 D9 与 DR13-FIX 已把「镜像按 run root 唯一」操作化为 console 单任务限制的依据；只论证「同时」字面即可容纳串行，会迫使 ⑦ 在旧裁决与新字面间裁量，违反「协议先行、不自行选择」纪律。故显式修订，明确 D9 判读被本段取代。

**契约修订草案（追加到 CONTRACTS.md:254 A6 段末，`_EN` 同步翻译）**：

> **多任务串行装配（P50）**：runRoot 级身份镜像与启动槽是「**同一时刻至多一个在跑 launch**」的排他点，不构成 run 生命周期内的一次性唯一。同一 `runId` 的一个 chain 可含多个 task（每 task 一个 isolated-workspace code step）；console 装配层为每个 `(taskId, stepId, attempt)` 绑定各自已持久化的 `agent-runner-request`（各 `requestId` 互异，R/C/终局/工件按 requestId 隔离），并由 chain engine **串行**调度：前一 task 的 launch 必须已证明停机（槽位释放、镜像所指进程死亡或镜像已由新 launch 原子轮换）后，后一 task 才允许发布 R 与 launch-intent；镜像轮换=以新 launch 身份原子覆盖旧镜像，与既有「镜像命名死亡进程时为陈旧、可被替换」语义一致。同一 runRoot 内禁止并发在跑的 launch；并发 Start 败者依旧不得创建日志工件、不得覆盖身份镜像。事件日志仍是该 runRoot 唯一 `events/state-events.jsonl` 全序源；各 task 六枚冻结事实按各自终局路由事件的位置化证据（A5 `evidence[2..7]`）归属，不得跨 task 对账。

**Schema/夹具联动**：无。镜像/槽位为 store-local（CONTRACTS.md 已述「不进 wire、不参与 R/C 的 recordHash」）；`agent-runner-request.schema.json` 已有 taskId/stepId/requestId 字段，多记录复用。修订范围仅 CONTRACTS{,_EN}.md + 本设计 + P50 台账行；**不得改写** p45/DR13 历史报告（以 P50 报告记录取代事实）。

## §3 命题二：③ 红测规格

**文件**：`internal/console/runtime_multitask_red_test.go`（`package console`，与既有 `runtime_zero_model_e2e_test.go` 同包风格）。

| 用例 | 入口 | 断言（当前实现下） |
|---|---|---|
| `TestAgentRunnerAssemblyRefusesMultiTaskDefinition` | `ExecuteAgentRunnerRun`（`runtime_agent.go:103`） | 夹具：2 个 task，各 1 个 `Kind=code`+`Mode=IsolatedWorkspace` step 的合法 `chain.Definition`；`AgentRunnerAssemblyOptions{Definition:…}`（其余零值）。断言 `errors.Is(err, ErrInvalidConfig)` 且 `err.Error()` 含 `"the AgentRunner assembly requires exactly one task"` |
| `TestSupervisedAssemblyRefusesMultiTaskDefinition` | `ExecuteSupervisedRun`（`runtime_supervised.go:248`） | 夹具：2 个 task，各 1 个 supervised-black-box code step 且带 `HandoffPolicy`；断言 `errors.Is(err, ErrInvalidConfig)` 且含 `"the supervised assembly requires exactly one step"` 或 `"exactly one task"`（实现序：task 检查先于 step 检查，`:330-333`） |
| `TestMultiTaskDefinitionShapeIsOtherwiseValid`（控制腿） | 无入口（纯夹具校验） | 同一 2-task 夹具通过 chain 定义级校验（引用 P45 已证 3-noop-task 引擎腿 `engine_zero_model_e2e_test.go`）；并证 1-task 同形夹具越过步骤选择、落在**后续**校验（`"request record is missing or invalid"`，`runtime_agent.go:119`）——证明 2-task 夹具之红**只**来自单任务门，不误伤他项 |

**稳**：两断言点在 ctx 检查之后、任何 I/O 之前（`runtime_agent.go:113-116`、`runtime_supervised.go:251-254`），零文件、零 goroutine、零时序，确定性红。
**准**：精确文案 + 控制腿；不动 adapters；不依赖 runRoot/时钟。
**狠（绿形态）**：② 改完后两拒绝文案消失 → 红测**翻转**为：同夹具经 `ExecuteAgentRunnerRun` 不再在步骤选择被拒，而是在**离线 harness 内两 task 串行各自完成一次终局路由**；③ red-baseline 必须记录该翻转形态。同时 ② 必须新增 `TestAgentRunnerAssemblyAcceptsMultiTaskDefinitionShape`（任何 `len(definition.Tasks)!=1` 回归即红），使「狠」经翻转保留而非删除。

## §4 命题三：安全边界

- **流分隔键**：`requestId` 是 wire 与磁盘工件流键（A2：R/C 去重键；A6：`<runRoot>/agent-runner-replay/runs/<requestId>/` 工件；launch-intent/identity/terminal 均 requestId 键控）；装配侧用 `(runID, taskID, stepID, attempt)` 映射到唯一 requestId；`launchID` 仅单 launch 生命周期内有效。
- **Facts/事件流不串号**：① chain 事件日志单一 JSONL、单写者串行追加，全序即定义序；② A5 六事实是**该 task 终局路由事件**的 `evidence[2..7]` 位置化载荷，`SubmitAgentRunnerTerminal`+`LoadChainAgentRunnerTerminal(store, requestRecord)` 双重绑定 requestId；③ A5 已钉「同一 task 的 completed 终局来自多 step 则整体拒绝」；④ 每 task 运行时态（`beforeManifest`/`acceptedEntries`/`capturer`）从结构体字段改为按 `(taskID, attempt)` 键控映射，杜绝跨 task 复用（`agentRunnerRuntime`，`runtime_agent.go:311-330`）。
- **并发安全边界**：装配内**禁止并发 spawn**（S1 纪律：engine 串行调度，一 launch 存活期间无第二 launch）；跨进程并发 Start 由 A6 槽位 O_EXCL + `refuseLiveMirrorProcess` 兜底，DR13-FIX D1/D2 已封同进程误回收；镜像原子写 temp+rename（`agent_runner_process_registry.go:263`）；槽位释放仅当 owner 匹配。
- **线性化点**（「持久化成功即线性化点」）：① 槽位独占创建 = launch 线性化点；② R 发布（`first-dispatch`）= dispatch 线性化点；③ 镜像发布 = 可寻址线性化点；④ C 发布 = 终局线性化点；⑤ 事件追加 = chain 状态线性化点。多任务 = 逐 task 重放该五步；任务间轮换前必须完成停机的持久证明（`managedProcessesGone`，`runtime_agent.go:536`，或镜像 stale 判定），失败即 fail-closed 暂停 run。
- **平台差异**：Ubuntu `rename(2)` 原子替换；Windows `MoveFileEx(REPLACE_EXISTING)` 对打开中的目标可能 ACCESS_DENIED——镜像轮换失败必须 fail-closed 不自动重试；unix flock 的 close-fd 释放语义（DR13 W2）vs Windows LockFileEx；槽位 mtime 墙钟新鲜度受挂起/回拨影响（契约已定镜像兜底）；`ProcessAlive` 跨平台语义差异以 ubuntu CI 腿为最终判据（P45 已知限制）；console `-race` 扩包为既有待办片，不阻塞。

## §5 ② 实施步骤（文件级、按序）

1. `docs/CONTRACTS{,_EN}.md`：A6 段追加 §2 草案（协议先行，先于任何代码）。
2. `internal/console/runtime_multitask_red_test.go`：三条用例落盘，跑红存证（red-baseline）。
3. `internal/console/runtime_agent.go`：`agentRunnerAssemblyStep`(:769)→`agentRunnerAssemblyPlan`（产出 task 绑定表；**仍拒 0 task**）；`AgentRunnerAssemblyOptions.RequestRecord`→`RequestRecords`（按 (taskID,stepID) 注册表）；`agentRunnerRuntime` 运行时态映射化；驱动循环改为逐 await「定位记录→等待→发布终局→提交」；CLI 接线（`cmd/prfrail`）随后片完成。聚焦测试：`TestAgentRunnerAssemblyAcceptsMultiTaskDefinitionShape`。
4. `internal/adapters/agent_runner_intent.go`：新增多记录 preparer 注册表，复用既有 `matches==1` 逐记录校验（:20-38）。聚焦测试：多记录各绑定正确、跨记录 drift 拒绝。
5. `internal/console/runtime_multitask_agent_test.go`（绿测）：两 task 串行 stub 夹具——顺序、每 task 六事实独立、R/C 独立、事件全序。
6. `internal/console/runtime_supervised.go`：同构改造（supervised 会话/workspace 状态按 task 键控）；绿测 `runtime_multitask_supervised_test.go`。
7. 门禁：`node tools/gates/gate.js --all --skip-base`、`gofmt -l`、`go vet ./...`、`go test -count=1 ./...`。
8. 文档：`docs/t027/REMAINING_SLICES{,_EN}.md` P50 行、ARCHITECTURE 装配路径注记。

## §6 风险与边界

- **禁止事项**：不降级本片为 `[ROUTINE]`（L 档核心架构，白名单角色照常）；不推 gitee；`commit`/`push` 停 ⑩ 等同一轮授权；不新铸错误哨兵（复用 `ErrInvalidConfig`/`ErrAgentRunnerLaunchConflict`）；不改 `internal/adapters` 的 slot/mirror/恢复实现（DR-9/DR13 敏感面）；S1 内不做并行 spawn；不改 A2/A4/A5 语义；不改写历史验证报告。
- **主要风险**：① supervised 多 task 会话恢复复杂度最高（`loadSession` 按 task 键控），若进度吃紧可拆后置片；② Windows 镜像轮换在打开句柄下的失败模式需 CI 实测；③ 每 task 请求记录的 CLI 供给契约未定（谁写、何时写），必须在 ② 前经主控确认，否则退回「仅离线 harness 装配」范围。

DESIGN: DONE
