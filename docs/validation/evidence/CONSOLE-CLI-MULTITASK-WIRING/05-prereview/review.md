# ⑤ 预审结论原文 r1 —— P51 `CONSOLE-CLI-MULTITASK-WIRING`

载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读）。实得输出约 3200 字符（载体自报为估算值）。

结论：`PRE-REVIEW: PASS WITH FIXES`（1 High / 2 Low / 3 Info）。

---

## 1. 审查范围与方法

**只读命令**：`git rev-parse HEAD`（= `3aa53d5`，与切片锚点一致）、`git status --short`、`git diff HEAD -- internal/console/cli.go`（+116/−9，与声称一致）、`git diff --stat`、`go test -count=1 ./internal/console/`（`ok … 8.375s`，全绿）、`go build ./...`（ok）、`gofmt -l`（无输出，格式干净）。

**读过的文件**：`internal/console/cli.go` 全文（重点 `executeRun`:273–381、`agentRunnerRunInputs`:386–430、`executeAgentRunnerRun`:443–563、`readAgentRunnerRequestRecords`:564–616、`agentRunnerExecutableBindings`:618–636、`runUsage`:382、`newFlagSet`/`writeUsageError`/`resolvePathFromCWD`:697–774）；`runtime_agent.go`（`AgentRunnerAssemblyOptions`:72–111、`ExecuteAgentRunnerRun` 预条件、`agentRunnerAssemblyPlan`:1052–1072、`agentRunnerAssemblyRecords`:1076–1119、`agentRunnerRequestRecordKey`:750）；`config.go` `FindExecutableSteps`/`FindSupervisedSteps`:575–604；`adapters/agent_runner.go` `Decode/ValidateAgentRunnerRequestRecord`、`ErrInvalidAgentRunnerRequest` 文案；`cli_multitask_wiring_test.go` 全文（5 条）；`cli_test.go` 测试基座（`newTestCLI`/`runCLI` 走生产 `Execute`）；`docs/CONTRACTS.md` A6；证据包 `design.md`、`red-baseline.md`、`implementer-reports.md`。

## 2. 发现

| # | 严重度 | 位置 | 问题 | 证据 | 建议 |
|---|---|---|---|---|---|
| F1 | **High** | `cli.go:352–362`（路由） | 设计 §6 ⑦「`--request-records` 配无 executable 步的链必须拒绝」**未实现**：noop-only 链（`FindExecutableSteps` 为空）上 `executeRun` 三个路由守卫均不触发，`--request-records`（+七件套或不带七件套）被静默忽略，走 `ExecuteNoopRun` 并以 exit 0 成功收尾，写入事件日志。这正是切片要防的「静默降级」。实现的拒绝文案 `request records require at least one executable step...` 位于 `readAgentRunnerRequestRecords`（:601–603），该函数仅在 `len(executable)>0` 路由后才可达——noop-only 链根本走不到。设计 §3 括注仅豁免旧 `--request-record` 的 noop 链行为，未豁免新 flag。 | 代码路径逐行核验：`executeRun` 中 `len(executable)>0` 为假 → 跳过互斥后所有 AgentRunner 判定 → `ExecuteNoopRun` → `writeRunSummary` exit 0；全程无任何 `requestRecords` 引用。设计 `design.md` §3/§6 明文要求 exit≠0。 | 在 `executeRun` 互斥判定之后加：`if len(inputs.requestRecords) > 0 && len(executable)==0 && len(supervised)==0 { return writeCommandError(..., "request records require at least one executable step in the chain definition") }`，并补一条红→绿测试腿。 |
| F2 | Low（不构成 Medium） | `cli_multitask_wiring_test.go` | 设计 §6 的 7 条 fail-closed 中，③「键未命中定义」（:611）与④「定义缺覆盖」（:616）两条核心文案**无测试腿**；⑥空路径、⑦noop 链也无腿。判定：**不构成 Medium**——设计 §5 测试规格明文只定 R1–R4+G1 共 5 条，实施与③均按规格交付并在③报告偏差 4 登记「规格未要求」；③④ 是确定性纯文案拒绝分支，装配层 `agentRunnerAssemblyRecords`（runtime_agent.go:1079–1085）对 registry 缺失仍有独立拒绝。若切片质量标准要求「每条 fail-closed 文案必有腿」，可升级为 Medium。 | 5 条测试与 `red-baseline.md` 清单逐条吻合（4 红 1 绿，红因均为 unknown flag）；③报告偏差 4 自认未覆盖。 | 可选：后续片补 ③④ 两条腿；F1 修复必须同时补 noop 链腿。 |
| F3 | Low | `cli.go:609–613` | 注释声称「refusal order is fixed so the CLI failure is deterministic」，但「键未命中定义」检查遍历 `records` map，**多条 stray 记录时报告哪个键不绑定是非确定的**（Go map 迭代序）；覆盖检查遍历 `bindings` 切片是确定的。错误**种类**确定、**文本**可能漂移。 | 逐行核验循环对象；`for key := range records`。 | 改为先收集 stray 键并排序后报错，或遍历 `paths` 保持输入序。 |
| F4 | Info | `cli.go:382` | `runUsage` 中 `(--request-record <path> \| --request-records <path>)...` 的 `...` 位于分组外，读作整组可重复，暗示 `--request-record` 也可重复（实为单值 String flag）。文案为设计「usage 逐字」给定，非实施缺陷。 | 与 `design.md` §3 逐字一致。 | 回馈设计修订文案。 |
| F5 | Info | `cli.go:387–389` | 结构体上方注释仍写「every one of the eight values」，现已 9 字段（八件套 + 多记录维度）。②a 报告已登记，规格未要求。 | 注释与字段数不符。 | 顺手改注释。 |
| F6 | Info | `cli.go:593` | 多记录路径 `DecodeAgentRunnerRequestRecord` 错误原样返回（无 `request record` 标签）；但与单值路径（:487–489）**完全一致**——解码错误在两路均不带标签，标签只覆盖路径解析/文件读取。诊断一致性达标，仅可读性可提升。 | `ErrInvalidAgentRunnerRequest` 文案 `invalid AgentRunner request: ...` 两路相同。 | 可选：多路径解码错误包一层 `request record: %w`。 |

其余重点项结论（无发现）：

1. **路由完整性**：supervised 链 + `--request-records` → `anyProvided()` 命中 → `executeRun:334–338` 拒绝（exit 1），无回退；双 flag 同给先于路由以 exit 2 拒绝；executable 链 + 完整九件 → AgentRunner。严格互斥，除 F1 外无静默降级。
2. **错误处理**：`readAgentRunnerRequestRecords` 全失败面均经 `writeCommandError` 非零退出；无 panic/越界；成功时 map 必非 nil，失败时立即返回不触碰装配层；`RequestRecords: nil`（单值路径）时装配层 `len(...)>0` 为假，单值兼容路径不变。
3. **单任务不变**：`complete()` 与旧语义等价；`anyProvided()` 仅前置追加；互斥判定位于 `inputs` 构造后、runID 校验/`resolvedRunDir`/事件日志 stat/`ToDefinition` 之前，**早于 run-root 创建**。
4. **口径漂移**：`agentRunnerExecutableBindings` 三条件与 `agentRunnerAssemblyPlan` 逐字一致，不存在「CLI 放行、装配接受」的反向漂移；差异仅为「跳过 vs 整体拒绝」，装配层仍兜底拒绝且零写入。
5. **flag 语义**：`set.Func` 与位置参数校验无交互冲突；空/空白值在 Parse 期即拒（`request record path is empty`，exit 2）；同一路径重复两次 → 同键 → 重复键拒绝。
7. **回归面**：全仓无测试引用 `runUsage` 或旧 usage 字面量；`go test ./internal/console` 全绿、`go build ./...` ok、gofmt 干净。

## 3. 可证伪性声明

- **F1**：构造 noop-only 链 + 七件套 + 合法多记录 → 若实测 exit≠0 且 stderr 含 `request records require at least one executable step`，则发现被推翻（我基于代码路径静态核验，未落盘实测，因只读约束未创建临时链）。
- **F2**：若设计 §5 之外另有「每文案必测」的质量门禁条文（本片证据包未列），则「不构成 Medium」结论被推翻。
- **F3**：若两条以上 stray 记录时连续多次运行 stderr 报同一键，则非确定性结论被推翻。
- **F4/F5/F6**：文案/注释改动或两路标签要求变更即推翻。
- 其余「无发现」项：任何一条被实测出静默 exit 0 或零写入被破坏（run-root 提前创建）即推翻。

## 4. 结论

PRE-REVIEW: PASS WITH FIXES

OUTPUT_CHARS=约 3200（正文字符估算，含表格）
