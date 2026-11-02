# P50 `CONSOLE-MULTITASK-ASSEMBLY` · ③ 红测报告

## §A 测试文件与用例清单

**新增文件**：`internal/console/runtime_multitask_red_test.go`（`package console`，无 BOM + LF，`gofmt -l` 无输出，`go vet ./internal/console/` 无告警）
**生产代码**：零改动（`git status` 仅 `?? internal/console/runtime_multitask_red_test.go` 为本片新增）。

### 夹具（三个用例共用）

| 夹具 | 内容 |
|---|---|
| `multiTaskIsolatedDefinition()` | `Definition{ID:"chain-multitask-assembly"}`，2 task：`task-one/code-one`、`task-two/code-two`，均 `Kind="code"` + `Mode=chain.IsolatedWorkspace` |
| `multiTaskSupervisedDefinition()` | 同上 ID 与 task/step 名，均 `Kind="code"` + `Mode=chain.SupervisedBlackBox`，`HandoffPolicy{AllowedTargets:["target-one"], InputPolicy:"structured", HandoffTimeoutMs:60000, ReturnActions:["complete","abort"], HooksAfterReturn:["hook-one"]}`（满足 `validateHandoffPolicy` 全部条件） |
| `singleTaskIsolatedDefinition(task)` | 从 2-task 夹具裁出单 task，step 形状不变 |

### 断言字符串（逐字取自源码）

| 常量 | 字面值 | 来源 |
|---|---|---|
| `agentRunnerMultiTaskRefusal` | `the AgentRunner assembly requires exactly one task` | `runtime_agent.go:772` |
| `agentRunnerSingleStepRefusal` | `the AgentRunner assembly requires exactly one step` | `runtime_agent.go:776` |
| `agentRunnerStepShapeRefusal` | `the AgentRunner assembly requires one isolated-workspace code step` | `runtime_agent.go:780` |
| `agentRunnerPostSelectionRefusal` | `request record is missing or invalid` | `runtime_agent.go:119` |
| `supervisedMultiTaskRefusal` | `the supervised assembly requires exactly one task` | `runtime_supervised.go:330` |
| `supervisedSingleStepRefusal` | `the supervised assembly requires exactly one step` | `runtime_supervised.go:334` |
| `supervisedStepShapeRefusal` | `the supervised assembly requires one supervised-black-box code step` | `runtime_supervised.go:338` |
| `supervisedPostSelectionRefusal` | `run root is required` | `runtime_supervised.go:255` |

### ① `TestAgentRunnerAssemblyRefusesMultiTaskDefinition`
1. `chain.ValidateDefinition(def) == nil`（夹具本身合法，红不来自形状）
2. `len(def.Tasks) == 2`
3. `ExecuteAgentRunnerRun(context.Background(), AgentRunnerAssemblyOptions{Definition: def})`：`err != nil`
4. `errors.Is(err, ErrInvalidConfig)`
5. `strings.Contains(err.Error(), agentRunnerMultiTaskRefusal)`
6. 空摘要：`RunID=="" && ChainState=="" && Sequence==0 && EventLogPath==""`（拒绝发生在 run-root marker 与引擎之前）

### ② `TestSupervisedAssemblyRefusesMultiTaskDefinition`
1–2 同上；3 `ExecuteSupervisedRun` 返回 `err != nil`；4 `errors.Is(ErrInvalidConfig)`；5 含 `supervisedMultiTaskRefusal`；6 **负断言**：不含 `supervisedSingleStepRefusal`（钉死「task 检查先于 step 检查」的次序）；7 空摘要同上。

### ③ `TestMultiTaskDefinitionShapeIsOtherwiseValid`（控制腿）
1. 两个 2-task 夹具均通过 `chain.ValidateDefinition`（引擎启动前的同一套语义校验；等价 P45 多 task 引擎腿的前提）
2. **逐 task 单独喂入**：2-task 夹具的每个 task 单独成 definition 时，越过步骤选择、不含任一单任务/单 step/step 形状拒绝文案，且落在**后续**校验——AgentRunner 落在 `request record is missing or invalid`，Supervised 落在 `run root is required`（均 `errors.Is(ErrInvalidConfig)`）
3. 2-task 两入口的错误**只**含 task 门文案，不含 step / step 形状文案 → 红的唯一来源是单任务门

## §B 红基线运行日志（原文）

**聚焦三条：**
```
=== RUN   TestAgentRunnerAssemblyRefusesMultiTaskDefinition
--- PASS: TestAgentRunnerAssemblyRefusesMultiTaskDefinition (0.00s)
=== RUN   TestSupervisedAssemblyRefusesMultiTaskDefinition
--- PASS: TestSupervisedAssemblyRefusesMultiTaskDefinition (0.00s)
=== RUN   TestMultiTaskDefinitionShapeIsOtherwiseValid
--- PASS: TestMultiTaskDefinitionShapeIsOtherwiseValid (0.00s)
PASS
ok      github.com/larsonzh/prfrail/internal/console    1.405s
```

**全包回归：**
```
ok      github.com/larsonzh/prfrail/internal/console    8.495s
```

**格式/静态检查：**
```
gofmt -l internal\console\runtime_multitask_red_test.go   → 无输出
go vet ./internal/console/                                → 无输出
```

## §C 与设计 §3 的偏差

1. **零偏差：零值 Options 即达单任务门。** 两入口的调用链实为 `ctx.Err()` → `agentRunnerAssemblyStep` / `supervisedAssemblyStep`，中间**无**其他前置校验（`ExecuteSupervisedRun` 的 `supervisedClockAndIDs` 只填默认值，不校验）。故 `AgentRunnerAssemblyOptions{Definition:…}` 与 `SupervisedAssemblyOptions{Definition:…}` 的其余零值可达门，无需补字段。
2. **Supervised 文案确认为 task 版**：设计 §3 写「`"exactly one step"` 或 `"exactly one task"`（按源码检查顺序）」，源码顺序确为 task→step（`supervisedAssemblyStep` 首分支即 `len(Tasks)!=1`），故断言 `requires exactly one task`，并加负断言锁定顺序。
3. **控制腿加强（设计未强制，属加项）**：设计只要求「AgentRunner 侧 1-task 夹具落在后续校验」；本实现额外对**每个 task 单独**跑两侧入口（4 次调用），并额外覆盖 supervised 侧「后续校验 = `run root is required`」，使「2-task 之红只来自单任务门」得到 per-task 粒度证明。
4. **执行侧编码修正**：文件初写为 CRLF，已转为 UTF-8 无 BOM + LF（仓库对 `.go` 的硬规则），内容语义未变。

## §D 绿形态（② 改完后这三条如何翻转）

| 用例 | ② 改完后的表现 | ② 必须做的改写 |
|---|---|---|
| `TestAgentRunnerAssemblyRefusesMultiTaskDefinition` | **变红**：2-task 定义不再被步骤选择拒绝，零 `RequestRecords` 下落到 `request record is missing or invalid`，第 5 断言（含 task 门文案）失败 | 按设计 §3 改写为 `TestAgentRunnerAssemblyAcceptsMultiTaskDefinitionShape`：同夹具**不得**含 `requires exactly one task`，且（补记录后）进入串行驱动；保留「拒绝文案已消失」的负断言，使回归可被捕获 |
| `TestSupervisedAssemblyRefusesMultiTaskDefinition` | **变红**：落到 `run root is required`，task 门文案消失 | 同构改写为 supervised 接受形；保留 `requires exactly one step` 与 `requires one supervised-black-box code step` 的负断言 |
| `TestMultiTaskDefinitionShapeIsOtherwiseValid` | **部分翻转**：第 2 步的 per-task 断言全绿（每个 task 单跑仍落在 `request record…` / `run root…`）；第 3 步（2-task 错误须含 task 门文案）变红 | 删除/反转第 3 步为「2-task 错误**不得**含任何单任务/单 step/step 形状文案」；第 1、2 步原样保留，作为永久形状回归腿 |

## §E 给 ② 实现的五条注意点

1. **文案精确性**：`agentRunnerMultiTaskRefusal` / `supervisedMultiTaskRefusal` 是逐字断言；② 删除门时**必须同步改写测试**，不得只留旧的「含门文案」断言（否则红测永久假绿）。
2. **保留 0-task 拒绝**：设计 §5 第 3 步要求 `agentRunnerAssemblyPlan` **仍拒 0 task**；该文案（`requires exactly one task`）与 `len(Tasks)!=1` 复用同一条，若改成 `len(Tasks)==0` 分支请新增专用文案并补 0-task 用例，避免与多 task 接受语义混淆。
3. **检查顺序不可回退**：supervised 侧 task→step 顺序已由负断言钉死；把 step 检查提到 task 之前会让 `TestSupervisedAssemblyRefusesMultiTaskDefinition` 的负断言（当前绿）在改写后暴露次序漂移。
4. **夹具复用**：三组夹具（`multiTaskIsolatedDefinition` / `multiTaskSupervisedDefinition` / `singleTaskIsolatedDefinition`）可直接被绿测 `runtime_multitask_agent_test.go`、`runtime_multitask_supervised_test.go` 复用；step ID `code-one`/`code-two`、task ID `task-one`/`task-two` 已满足 `evidence.ValidID`（`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`）且全局唯一。
5. **记录注册表与执行顺序**：`RequestRecords` 改为按 `(taskID, stepID)` 键控后，`request record is missing or invalid` 的**首个**失败点应保持「查表未命中即失败」，并保持 `record.RunID/TaskID/StepID` 绑定校验逐记录执行；同时 `AgentRunnerAssemblyOptions{Definition:…}` 其余零值仍应在单任务门之后**立即**被拒绝（不得先落盘 run-root marker），以维持本片第 6 断言（空摘要）在改写后的等价形态。

TDD-RED: ESTABLISHED
