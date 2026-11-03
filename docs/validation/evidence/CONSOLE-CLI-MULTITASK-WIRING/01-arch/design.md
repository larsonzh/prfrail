# ① 架构设计原文 —— P51 `CONSOLE-CLI-MULTITASK-WIRING`

载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读，额度 1/1）。实得输出字符数 2896（自报 `OUTPUT_CHARS=2896`）。

主控机械复核的事实点：`runtime_agent.go:88` `RequestRecords`、`:750` `agentRunnerRequestRecordKey`、`:1076-1080` `agentRunnerAssemblyRecords`、`config.go:36` `Tasks []TaskConfig`、`cli.go:307-308` `FindExecutableSteps/FindSupervisedSteps`、`cli.go:321/325` 互斥与 noop-only 判定、`cli.go:387/407` `complete()/anyProvided()` 均与设计文档所述一致。

---

# P51 `CONSOLE-CLI-MULTITASK-WIRING` 设计文档（① 架构，只读）

## 1. 调用链现状（文件:行号）

```
cmd/prfrail/main.go:15-19   main() → console.NewCLI(version) + cli.Execute(ctx, os.Args[1:], ...)
internal/console/cli.go:54  NewCLI
internal/console/cli.go:70  Execute；case "run"（约 :89）→ executeRun
internal/console/cli.go:274 executeRun
internal/console/cli.go:280 requestRecord := set.String("request-record", ...)   ← 唯一 String flag，无多记录注册
internal/console/cli.go:307-308 FindExecutableSteps / FindSupervisedSteps 决定三路由
internal/console/cli.go:309-316 agentRunnerRunInputs 八字段构造（无多记录字段）
internal/console/cli.go:314-318 P1 互斥（supervised × anyProvided）
internal/console/cli.go:320-323 complete() 全有判定 + noop-only refusal
internal/console/cli.go:355 executeAgentRunnerRun（executable>0 路由）
internal/console/cli.go:428 executeAgentRunnerRun 解码八件套
internal/console/cli.go:504-523 AgentRunnerAssemblyOptions 构造：只填 RequestRecord（约 :509），无 RequestRecords
internal/console/runtime_agent.go:72-98   AgentRunnerAssemblyOptions{RequestRecord :80, RequestRecords map :88}
internal/console/runtime_agent.go:126    ExecuteAgentRunnerRun（Step 1 预条件）
internal/console/runtime_agent.go:146-149 agentRunnerAssemblyPlan → agentRunnerAssemblyRecords
internal/console/runtime_agent.go:1076-1105 agentRunnerAssemblyRecords：注册表非空走 :1079-1085，否则单记录仅限单 task（:1090 起）
internal/console/runtime_agent.go:750    agentRunnerRequestRecordKey = taskID+"/"+stepID
```

**缺口确切位置**：`cli.go:504-523` 的 options 构造块——该块赋值了 `RequestRecord` 却从不赋值 `RequestRecords`；上游 `cli.go:280` 未注册任何可重复的 request-record flag，`agentRunnerRunInputs`（:376）无对应字段。装配层（`runtime_agent.go`）多记录能力完整，缺口 100% 在 CLI 输入面。

## 2. 缺口判定

(a) **能**。`ChainConfig.Tasks []TaskConfig`（`internal/console/config.go:36`），`ToDefinition` 循环映射全部 task；`agentRunnerAssemblyPlan`（runtime_agent.go:1052-1072）显式接受多 task（每 task 恰一个 isolated-workspace code 步），CONTRACTS.md:255「多任务串行装配（P50）」已定契约。
(b) 多 task 定义 + 单 `--request-record`：八件套 `complete()` 通过 → `cli.go:355` 路由 → options 仅填 `RequestRecord` → `ExecuteAgentRunnerRun` Step 1 `agentRunnerAssemblyRecords`（runtime_agent.go:1076）走 `len(plan)>1` 分支，在 **:1091-1093** 返回 `request record is missing or invalid: 2 tasks need a request record registry`。此时 run-root 标记（Step 2）未创建，零写入、零进程——安全拒绝，但错误发生在装配层、文案指向 CLI 并不存在的「registry」概念，且 CLI 无任何途径表达多记录 ⇒ 功能性缺口。
(c) **不必改 supervised**。`SupervisedAssemblyOptions`（runtime_supervised.go:215）确有 `TaskPolicies map[string]chain.HandoffPolicy`（:219-232），但它是 handoff 策略的**一致性钉**，不是记录注册表；supervised 模式无 request-record 概念。CLI 侧 `supervisedAssemblyOptions`（supervised_cmd.go:53-70）不填 `TaskPolicies` 是合法兼容路径。P1 互斥（cli.go:314-318）已隔离两族输入，新 flag 并入 `anyProvided()` 后自动继承互斥。结论：仅接 AgentRunner 线。

## 3. CLI 参数形态设计（确定方案）

**flag**：新增 `--request-records <path>`，**可重复**，经 `set.Func` 注册（先例 `ai.go:122`）。值**仅为记录文件路径**，不带 `<taskID>/<stepID>` 前缀——键由记录内容推导（记录自带 `TaskID/StepID`，`agent_runner.go:26-48`）。`--request-record` 保持 `set.String` 单值不变。两者**互斥**。

**数据结构**（小写、包内）：

```go
type agentRunnerRunInputs struct {
    requestRecord       string
    requestRecords      []string // repeatable paths, decoded in executeAgentRunnerRun
    capabilityRecord    string
    enforcementRecord   string
    availabilityRecord  string
    authorizationLedger string
    costLedger          string
    agentCLI            string
    agentVersion        string
}

type agentRunnerRequestBindings struct {
    records []adapters.AgentRunnerRequestRecord
    byKey   map[string]adapters.AgentRunnerRequestRecord
}
```

**键推导**：`agentRunnerRequestRecordKey(record.Request.TaskID, record.Request.StepID)`（console 同包可直接调用 :750，与装配层编码天然一致）。

**fail-closed 校验**（在 `executeAgentRunnerRun` 解码后、构造 options 前；错误文案纯英文、无占位词）：
- 空路径：Func 回调即拒 `request record path is empty`（同 ai.go 先例）。
- 重复键：`request records must bind distinct task/step keys: duplicate key %q`
- 键未命中定义：`request record does not bind a task/step of the chain definition: %q`
- 定义缺覆盖：`request record is missing for task %q step %q`（CLI 侧早失败，避免落入装配层 registry 文案）
- runId 不一致：`request record runId %q does not match run id %q`
- `--request-records` 非空但链无 executable 步：`request records require at least one executable step in the chain definition`（防静默降级；不改变旧 `--request-record` 对 noop 链的既有行为）

**选择规则**：`requestRecords` 非空 ⇒ 填 `RequestRecords`（`RequestRecord` 留零值；装配层注册表非空时忽略单记录，:86-94 已定）。仅 `requestRecord` 非空 ⇒ 填 `RequestRecord`（旧行为逐字节不变）。两者同给 ⇒ `executeRun` 解析后立即拒绝：`--request-record and --request-records are mutually exclusive`。

**complete()/anyProvided() 演进**：`anyProvided()` 追加 `len(inputs.requestRecords) > 0`；`complete()` 将 requestRecord 维度改为 `TrimSpace(requestRecord) != "" || len(requestRecords) > 0`，**其余七件保持不变**——多记录下 capability/enforcement/availability/ledger/agent-cli 仍是共享单值必填（装配层 process config 是共享 pin，见 `agentRunnerBindTasks` runtime_agent.go:1127 起），不得减负。互斥判定放在 `executeRun`（不入 complete()）。

**usage 逐字**：
`usage: prfrail run [--chain <path>] [--run-id <id>] [--run-dir <path>] [--json] (--request-record <path> | --request-records <path>)... --capability-record <path> --enforcement-record <path> --availability-record <path> --authorization-ledger <path> --cost-ledger <path> --agent-cli <path> --agent-version <version>`

## 4. 需要修改的文件清单

| 文件 | 性质 | 函数/区块 | 预估行数 | 回归风险 |
|---|---|---|---|---|
| `internal/console/cli.go` | 修改 | `executeRun` :275-356（flag 注册、互斥、判定）；`runUsage` :370；`agentRunnerRunInputs` :376-426（字段+complete+anyProvided）；`executeAgentRunnerRun` :428-529（多记录解码、校验、`RequestRecords` 赋值） | +120~160 | 低：仅输入面；旧单值路径判定逻辑不动 |
| `internal/console/cli_multitask_wiring_test.go` | ③ 新增 | 红测 R1-R4 + 绿腿 G1 | — | 无 |
| `internal/console/runtime_agent.go` | **不改** | 装配层是 P50 已验证语义，禁止触碰（见 §6） | 0 | — |
| `internal/console/runtime_supervised.go` / `supervised_cmd.go` | **不改** | supervised 无记录输入端，缺口不在此 | 0 | — |
| `internal/adapters/agent_runner.go` | **不改** | `AgentRunnerRequest` 已带 `TaskID/StepID`，解码/校验已存在 | 0 | — |
| `internal/console/config.go` | **不改** | 多 task 解析与路由已支持 | 0 | — |
| `cmd/prfrail/main.go` | **不改** | 纯转发 | 0 | — |
| `docs/CONTRACTS{,_EN}.md` | **不改** | 见 §6 | 0 | — |

## 5. 红测规格（③ 照做；文件 `internal/console/cli_multitask_wiring_test.go`）

构造要求：`t.TempDir()` 内用 `DefaultChainConfig` + 两个 `TaskConfig`（各一个 `code` 步、默认落 `IsolatedWorkspace`）经 `EncodeChainConfig` 写链 JSON；记录文件用 `adapters.NewAgentRunnerRequestRecord(...)` 构造后 `json.Marshal` 落盘（`DecodeAgentRunnerRequestRecord` 走严格 JSON）。CLI 调用方式同 `TestRunCommandRefusesMixedSupervisedAgentRunnerInputs`（`NewCLI` + `Execute`）。

- **R1** `TestRunRequestRecordsPathsAreRead`：两 task 链 + `--request-records <不存在路径>` ×1 + 其余七件套。断言 exit≠0 且 stderr 含 `request record`、不含 `flag provided but not defined`。未实现红：`--request-records` 是 unknown flag，报 `flag provided but not defined: -request-records`；实现后绿：进入解码，报 `request record: open ...`。
- **R2** `TestRunRefusesRequestRecordWithRequestRecords`：单 task 链 + `--request-record a.json --request-records b.json` + 七件套。断言 stderr 含 `mutually exclusive`。未实现红（unknown flag），实现后绿（互斥拒绝）。
- **R3** `TestRunRefusesDuplicateRequestRecordKey`：两 task 链 + 两条有效记录（requestId 不同但同 `TaskID/StepID`）。断言 stderr 含 `duplicate key`。未实现红（unknown flag）。
- **R4** `TestRunRefusesRequestRecordRunIDMismatch`：两 task 链 + 两条有效记录、键正确，其一 `runId`≠`--run-id`。断言 stderr 含 `does not match run id`。未实现红（unknown flag）。
- **G1（防回归绿腿，现在即绿）** `TestRunSingleRequestRecordStaysCompatible`：单 task 链 + 八件套、仅单 `--request-record` 指向不存在文件。断言 exit≠0 且 stderr 含 `request record:`。现在绿：`executeAgentRunnerRun` 的 `readRecord` 报错；实现后必须仍绿。它是防回归而非哨兵：钉死单值 flag 的解析与错误标签不被 Func 化/重命名破坏，且不依赖任何新代码路径（新代码只在 `requestRecords` 非空时介入）。

## 6. 安全与边界

- **禁止触碰**（runtime_agent.go）：`ExecuteAgentRunnerRun`、`AgentRunnerAssemblyOptions` 结构、`agentRunnerAssemblyPlan`、`agentRunnerAssemblyRecords`、`agentRunnerBindTasks`、`agentRunnerRequestRecordKey` 编码。CLI 只喂输入，不做任何装配层语义复制（键编码复用同包函数即可）。
- **CONTRACTS A6：不改**。A6（CONTRACTS.md:254-255）已定义离线运行口径与多任务串行装配，本片不新增 task/step 状态、wire 记录、chain 语义，只是把装配层既有输入接到 CLI flag 面；`_EN` 同步不改。
- **fail-closed 清单**（必须拒绝、禁止静默降级）：① 双 flag 同给；② 重复键；③ 键未命中定义；④ 定义缺覆盖；⑤ 记录 `runId` 与 `--run-id` 不一致；⑥ 空路径；⑦ `--request-records` 配无 executable 步的链。任一违反均 exit≠0 且零写入（校验在 run-root 标记之前）。

## 7. 回滚与最小面

全部改动集中在 `cli.go` 单文件（测试文件另计），可独立单 commit 整体 revert；装配层零改动，回滚后 `--request-records` 恢复为 unknown flag，旧行为与回滚前逐字节一致。

OUTPUT_CHARS=2896
