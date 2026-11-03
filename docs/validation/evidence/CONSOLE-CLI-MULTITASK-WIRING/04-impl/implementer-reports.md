# ② 实施角色报告合集 —— P51 `CONSOLE-CLI-MULTITASK-WIRING`

载体：`prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)`。按纪律 V 拆分为 ②a（CLI 参数面）与 ②b（多记录解码与装配注入）两步串行派单。仅修改 `internal/console/cli.go`。

---

## ②a 报告原文（实得 2380 字符，自报 `OUTPUT_CHARS=2380`）

### 1. 改动点（仅 `internal/console/cli.go`）

| 位置 | 改动 |
|---|---|
| `cli.go:281-288` | 新增局部 `requestRecordPaths := make([]string, 0)` 与 `set.Func("request-records", ...)`（仿 `ai.go:121-127` 的 `prior-record` 先例）；回调内 `strings.TrimSpace(value) == ""` 时返回 `fmt.Errorf("request record path is empty")`，否则 append。 |
| `cli.go:319` | `agentRunnerRunInputs` 构造中新增 `requestRecords: requestRecordPaths`（紧随 `requestRecord` 之后）。 |
| `cli.go:328-330` | 在 inputs 构造之后、`len(supervised) > 0 && inputs.anyProvided()` 判定之前加入互斥拒绝，调用 `writeUsageError(stdout, stderr, "run", *jsonOutput, "--request-record and --request-records are mutually exclusive", runUsage)`。 |
| `cli.go:382` | `runUsage` 改为规格逐字给定的字符串。 |
| `cli.go:390` | 结构体新增 `requestRecords []string`，紧接 `requestRecord` 之后。 |
| `cli.go:401-403` | `complete()`：requestRecord 维度改为 `TrimSpace(requestRecord) != "" \|\| len(requestRecords) > 0`，提前返回；其余七个条件保持原样。 |
| `cli.go:423-425` | `anyProvided()`：先判 `len(inputs.requestRecords) > 0` 返回 true，其余八项判断逐字不变。 |

`executeAgentRunnerRun` 未改动：其形参本就是 `inputs agentRunnerRunInputs`，`requestRecords` 已随结构体透传，本步未使用。

### 2. 验证结果（真实输出摘要）

**① `gofmt -l internal/console/cli.go`** → 无输出。

**② `go build ./...`** → `build exit=0`。

**③ `go vet ./internal/console/`** → `vet exit=0`。

**④ 绿腿组**

```
=== RUN   TestRunRefusesRequestRecordWithRequestRecords
--- PASS: TestRunRefusesRequestRecordWithRequestRecords (0.01s)
=== RUN   TestRunSingleRequestRecordStaysCompatible
--- PASS: TestRunSingleRequestRecordStaysCompatible (0.01s)
ok      github.com/larsonzh/prfrail/internal/console    1.107s
```

**⑤ ②b 范围组**

```
=== RUN   TestRunRequestRecordsPathsAreRead
--- PASS: TestRunRequestRecordsPathsAreRead (0.01s)
=== RUN   TestRunRefusesDuplicateRequestRecordKey
    cli_multitask_wiring_test.go:213: the refusal must name the duplicate key, stderr="run: request record: read C:\\Users\\...\\Temp\\TestRunRefusesDuplicateRequestRecordKey3330035656\\001: Incorrect function.\n"
--- FAIL: TestRunRefusesDuplicateRequestRecordKey (0.01s)
=== RUN   TestRunRefusesRequestRecordRunIDMismatch
    cli_multitask_wiring_test.go:233: the refusal must state the run id mismatch, stderr="run: request record: read C:\\Users\\...\\Temp\\TestRunRefusesRequestRecordRunIDMismatch2271077711\\001: Incorrect function.\n"
--- FAIL: TestRunRefusesRequestRecordRunIDMismatch (0.01s)
FAIL    github.com/larsonzh/prfrail/internal/console    1.084s
```

关键确认：三条 stderr 均不再含 `flag provided but not defined`，失败原因已前移到读取/解码阶段（`request record: read ...`），说明可重复 flag 已生效。

### 3. 偏差

1. `TestRunRequestRecordsPathsAreRead` 已转 PASS，而派单规格第 5 项期望它仍 FAIL。原因：该测试只断言 `stderr` 含 `request record` 且不含 unknown flag；②a 让 flag 被接受后即进入 `executeAgentRunnerRun`，首个 `readRecord` 因单值为空立报 `request record: read ...`，断言即满足。这与架构设计文档 R1 的"实现后绿"预期一致，属派单文本与设计文档口径不一致，已按设计文档执行。
2. 结构体上方注释仍写 "every one of the eight values"，未改（规格未要求）。
3. 工作区另有台账/证据包既有改动，非本步产生，未触碰。

### 4. 未执行项

无。

---

## ②b 报告原文（实得 2140 字符，自报 `OUTPUT_CHARS=2140`）

### 1. 改动点（`internal/console/cli.go`，仅此一个文件）

| 行号 | 说明 |
|---|---|
| `cli.go:460-479` | `executeAgentRunnerRun` 内新增 `var (requestRecord ...; requestRecords map[string]adapters.AgentRunnerRequestRecord)`，并在 `readRecord` 定义之后、单值读取之前分支：`len(inputs.requestRecords) > 0` → 调用 `readAgentRunnerRequestRecords`，否则走原单值读取/解码（原代码逐字节保留，仅包进 `else`） |
| `cli.go:541` | `AgentRunnerAssemblyOptions` 构造内新增 `RequestRecords: requestRecords`（紧随 `RequestRecord`，原有字段顺序不变）；`requestRecord` 保持零值 |
| `cli.go:566-616` | 新增 `readAgentRunnerRequestRecords(cwd, runID, definition, paths)`：逐条 `resolvePathFromCWD` + `os.ReadFile`（失败包装 `request record: %w`）→ `DecodeAgentRunnerRequestRecord` → runId 一致性 → 重复键 → 全部解码后键命中定义 → 覆盖完整性 |
| `cli.go:618-636` | 新增 `agentRunnerExecutableBindings(definition)`：按装配层计划口径枚举 task/step |

错误文案逐字实现（无 `registry` 字样、无中文/TODO）。

### 2. 枚举口径说明

`agentRunnerExecutableBindings` 与 `runtime_agent.go` 的 `agentRunnerAssemblyPlan`（:1052-1072）同一口径：遍历 `definition.Tasks`；仅当该 task 的 `Steps` 恰好 1 个，且该 step `Kind == "code"` 且 `Mode == chain.IsolatedWorkspace` 时，才产出绑定。其余 task 不产出绑定，本就会被装配层拒绝。绑定集合为空且 `requestRecords` 非空 ⇒ `request records require at least one executable step in the chain definition`。

判定顺序：单条内先 runId 再重复键 → 全部解码后键命中 → 最后覆盖完整。

### 3. 验证结果

1. `gofmt -l internal/console/cli.go` → 无输出。
2. `go build ./...` → `build-exit=0`。
3. `go vet ./internal/console/` → `vet-exit=0`。
4. 聚焦测试：

```
=== RUN   TestRunRequestRecordsPathsAreRead
--- PASS: TestRunRequestRecordsPathsAreRead (0.01s)
=== RUN   TestRunRefusesRequestRecordWithRequestRecords
--- PASS: TestRunRefusesRequestRecordWithRequestRecords (0.01s)
=== RUN   TestRunRefusesDuplicateRequestRecordKey
--- PASS: TestRunRefusesDuplicateRequestRecordKey (0.01s)
=== RUN   TestRunRefusesRequestRecordRunIDMismatch
--- PASS: TestRunRefusesRequestRecordRunIDMismatch (0.02s)
=== RUN   TestRunSingleRequestRecordStaysCompatible
--- PASS: TestRunSingleRequestRecordStaysCompatible (0.01s)
PASS
ok      github.com/larsonzh/prfrail/internal/console    1.399s
```

5. `go test -count=1 ./internal/console/` → `ok github.com/larsonzh/prfrail/internal/console 9.996s`。

改动范围：`git status --porcelain` 中仅 `internal/console/cli.go` 由本角色改动；未执行任何 `git add/commit/push`。

### 4. 偏差

无。（实现取舍：`RequestRecords: requestRecords` 写在构造字面量内，flag 缺省时该值为 `nil` map，装配层 `len(options.RequestRecords) > 0` 为假 ⇒ 单值路径行为不变。）

---

## ②c 报告原文（⑤ 预审整改；实得 2732 字符，自报 `OUTPUT_CHARS=2732`）

### 1. 改动点

**`internal/console/cli.go`**

- `cli.go:330-336`（F1）：在 `--request-record`/`--request-records` 互斥判定之后、supervised 与 executable 路由判定之前新增 fail-closed 分支：`len(requestRecords)>0 && len(executable)==0 && len(supervised)==0` ⇒ `writeCommandError(..., "request records require at least one executable step in the chain definition")`，exit=1，不再落入 `ExecuteNoopRun`。
- `cli.go:584`（F3）：`readAgentRunnerRequestRecords` 新增 `keys []string`，容量 `len(paths)`，按输入顺序记录解码出的 key。
- `cli.go:606`（F3）：每条记录入 map 后 `keys = append(keys, key)`，与 `paths` 同序。
- `cli.go:616`（F3）：「键未命中定义」检查由 `for key := range records` 改为 `for _, key := range keys`，报告顺序确定；错误文案不变。原有 `cli.go:610` 的「无 executable step」守卫保留，作为该函数被独立调用时的兜底。
- `cli.go:389-395`（F5）：`agentRunnerRunInputs` 注释改写为——request-record 维度须以「单值 `requestRecord` 或可重复 `requestRecords` 二选一」形式提供，其余七个值必须全部存在；不再声称「every one of the eight values」。

**`internal/console/cli_multitask_wiring_test.go`**

- `:241` 新增 `TestRunRefusesRequestRecordsWithoutExecutableSteps`（noop-only 链，`DefaultChainConfig`，先断言 `FindExecutableSteps==0`，再给不存在的 `--request-records` 路径；断言 exit≠0 且 stderr 含 `at least one executable step`）。
- `:273` 新增 `TestRunRefusesRequestRecordOutsideChainDefinition`（两 task 链 + 一条 `task-other/code-other` 记录；断言 exit≠0 且 stderr 含 `does not bind a task/step`）。
- `:293` 新增 `TestRunRefusesRequestRecordCoverageGap`（两 task 链 + 仅 `task-one/code-one` 记录；断言 exit≠0 且 stderr 含 `request record is missing for task`）。
- 复用既有 `p51*` 辅助函数，命名前缀 `TestRun...`；全部临时目录、无网络、无外部二进制。

### 2. 红→绿证据

**修复前（仅测试 3，cli.go 未改）**

```
=== RUN   TestRunRefusesRequestRecordsWithoutExecutableSteps
    cli_multitask_wiring_test.go:263: run must refuse --request-records on a noop-only chain, code=0 stdout=run: completed
        runId: run-p51-cli-multitask
        runDir: C:\Users\...\tmp\prfrail-runs\run-p51-cli-multitask
        chainState: COMPLETED
        events: 11
        eventLog: ...\events\state-events.jsonl
         stderr=
--- FAIL: TestRunRefusesRequestRecordsWithoutExecutableSteps (0.04s)
FAIL    github.com/larsonzh/prfrail/internal/console    1.232s
```

即：exit=0、`chainState: COMPLETED`、stderr 为空完全证实 F1 的「静默忽略并跑 noop 链」。

**修复后**

```
=== RUN   TestRunRefusesRequestRecordsWithoutExecutableSteps
--- PASS: TestRunRefusesRequestRecordsWithoutExecutableSteps (0.01s)
PASS
ok      github.com/larsonzh/prfrail/internal/console    1.781s
```

测试 4、5 属补充覆盖（实现后直接为绿，非红转绿），见 §3。

### 3. 其他验证结果

1. `gofmt -w` 已执行；`gofmt -l internal/console/cli.go internal/console/cli_multitask_wiring_test.go` → 无输出。
2. `go build ./...` → ok；`go vet ./...` → ok。
3. `go test -count=1 -run 'TestRun' -v ./internal/console/`：12 个 `TestRun*` 全 PASS。
4. `go test -count=1 ./internal/console/` → `ok github.com/larsonzh/prfrail/internal/console 10.778s`。
5. 改动范围核对：`git status --porcelain -- internal/console/` → 仅 `M internal/console/cli.go`、`?? internal/console/cli_multitask_wiring_test.go`；未执行任何 `git add/commit/push`。

### 4. 偏差

无。

---

## 主控机械复核（②a + ②b 后）

- `gofmt -l internal cmd`：无输出。
- `go build ./...`：exit 0；`go vet ./...`：exit 0。
- `go test -count=1 ./...`：15 个包全 `ok`（原始输出见 `04-impl/logs/04-full-suite-green.txt`）。
- 聚焦 5 条测试：4 红转 5 绿（红基线见 `03-tests/logs/03-tdd-red-run.txt`）。
- 改动面：`internal/console/cli.go` +116/−9（`git diff --stat`）。
