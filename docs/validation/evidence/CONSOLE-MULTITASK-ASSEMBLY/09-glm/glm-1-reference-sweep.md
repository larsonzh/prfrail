# P50 批量辅助 ①：旧单任务限制引用面清查（只读）

## 1. `agentRunnerAssemblyStep` / `supervisedAssemblyStep` 定义与调用
- `internal/console/runtime_agent.go:772` — 定义（:769 注释「requires exactly one…」）
- `internal/console/runtime_agent.go:117` — 生产唯一调用（ExecuteAgentRunnerRun）
- `internal/console/runtime_supervised.go:333` — 定义（:330 注释）
- `internal/console/runtime_supervised.go:253` — 调用（ExecuteSupervisedRun 主路径）
- `internal/console/runtime_supervised.go:294` — 调用（resume 路径）
- `internal/console/runtime_supervised_test.go:1421` — 测试直调（harness.options.Definition）
- 文档：`docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/01-arch/design.md:53`（计划改 `agentRunnerAssemblyPlan` 仍拒 0 task）、`03-tests/red-baseline.md:72-73`、`input-list.md:12-13`；`docs/validation/evidence/P45-T027B-E2E-ZERO-MODEL/01-arch-design/review.md:9`；`P38-T027B-TACTICAL-BLUEPRINT/01-sol-blueprint/review-batch-b.md:42`；t028 历史存档（ORCH-NOTES.md:85、runtime_supervised.go.orig:209/281 等，只读存档）

## 2. 单任务拒绝文案
- `internal/console/runtime_agent.go:774/778/782` — task 门 / step 门 / 形状门三文案
- `internal/console/runtime_supervised.go:335/339/343` — 同构三文案（task 先于 step 检查）
- `internal/console/runtime_multitask_red_test.go:37-42` — 六文案常量（本片新增，排除）
- `docs/validation/evidence/CONSOLE-MULTITASK-ASSEMBLY/03-tests/red-baseline.md:20-26,81-83,88` — 文案登记与翻转计划（:88：0-task 须保留/新增专用文案）
- 同目录 `01-arch/design.md:33-34` — 红测断言文案
- `docs/validation/evidence/directive-history/assembly-min|t028` 各 diff 存档 — 历史来源，仅追溯

## 3. `AgentRunnerAssemblyOptions.RequestRecord` 构造/读取
- 定义：`internal/console/runtime_agent.go:77`
- 读取：`runtime_agent.go:121`（Validate）、`:124`（record:=）、`:257`（装配 runtime 注入）、`:316`（fail 闭包取 RequestID）、`:323`（agentRunnerRunRequest）、`:328`（LoadChainAgentRunnerTerminal）
- 构造：`internal/console/cli.go:504-522`（`RequestRecord:` 字面量；来源 `cli.go:280` `request-record` 旗标）
- 测试构造：`internal/console/runtime_agent_e2e_test.go:329/:355`、`runtime_supervised_test.go:884`（跨模式拒绝）
- `tools/b4-e2e/fixtures.go:58`、`tools/b4-e2e/main.go:112`（emitRequestRecord）
- 计划改名 `RequestRecords`（按 (taskID,stepID) 注册表）：design.md:53、red-baseline.md:91
- 注：adapters 内大量 `AgentRunnerRequestRecord` 类型使用为记录本体（构造在 `internal/adapters/agent_runner.go:56`），非 Options 字段；②只需改 console 侧字段与 cli.go/测试三构造点

## 4. 运行时态字段（映射化对象）
- `agentRunnerRuntime`：`runtime_agent.go:374-394`
  - `beforeManifest`：写 `:362/:364`（preManifestAgentRunnerPort :352-366）、读 `:324`（WaitAgentRunnerDispatched 入参）
  - `acceptedEntries`：写 `:453`、读 `:497`（agentRunnerEntryMismatches）
  - `capturer`：构造 `:178`、用 `:257/:264/:323`
  - 其余 `runID/store/workspaceRoot/snapshotStore/requestRecord/clock/ids` 为 run 级
- `supervisedRuntime`：`runtime_supervised.go:366-409`：`taskID/stepID/policy`（单任务绑定）、`acceptedEntries` 写 `:1294` 读 `:1330`、`beforeManifestPath` `:516/:1547`、`acceptedHash/workspaceReady/events`
- 测试直构：`runtime_agent_e2e_test.go:620`（newAgentRunnerPIDReportRuntime 返回 *agentRunnerRuntime）、`runtime_supervised_test.go:1107`

## 5. `ExecuteNoopRun` / `RunPostflight`（D9）
- 定义：`internal/console/runtime.go:33`；调用：`internal/console/cli.go:363`、`runtime_zero_model_e2e_test.go:76/:115`、`stopper_test.go:24/:59`
- `RunPostflight` 端口：`internal/chain/postflight.go:159`；fail-closed 桩：`internal/console/runtime_noop.go`（localRuntime）；supervised 实现 `runtime_supervised.go:1357`；生产实现 `internal/adapters/agent_runner_postflight.go:278`；`tools/b4-e2e/postflight.go:47`
- D9 登记：`docs/validation/p45-e2e-zero-model.md:132`（_EN:132）、`docs/t027/REMAINING_SLICES.md:1520`（_EN:1507）、briefing.md:22-23/53/59、freedom-list.md:6-7

## 6. 文档层引用
- `docs/ARCHITECTURE.md:74`（_EN:74）— 单任务限制保持 + ExecuteNoopRun 贯通
- `docs/t027/CHECKLIST_v35.md:60/66/143/165`；`CHECKLIST_v36.md:46/135`（_EN:46）；`v37:70`、`v38:58`、`v39:58`
- `docs/t027/REMAINING_SLICES.md:1516/1520`（_EN:1503/1507）
- `docs/validation/p45-e2e-zero-model.md:16/127/132`（_EN 同行号）
- `docs/validation/strategic-review-2026-11-02.md:9/38`（_EN:9/38）
- evidence：`STRATEGIC-REVIEW-2026-11-02/00-briefing.md:99`、`01-round1/review.md:14`、`03-round3/review.md:13`；`P45…/05-prereview/review.md:25-27`、`07-final/review.md:50`；CONSOLE-MULTITASK-ASSEMBLY（briefing/design/red-baseline/input-list/freedom-list，见组 1-3）；`directive-boundary-clarify*.md:87` 与 `directive-gates-patch*.md:81` 的 D9 为**同名异义**（编号/事故），勿误改
- `docs/RFC-proofrail-unattended-ai-engineering-product.md:587/706` — 「单任务」指多语言，与装配限制无关
- 历史存档（directive-history/ 各 diff、runtime_supervised.go.orig）只读不改

## 7. 既有测试受影响判定
1. `runtime_agent_e2e_test.go:448` TestAgentRunnerAssemblyFullChain — **受影响**：驱动循环重构的直接回归面，1-task 语义须逐断言保持。
2. 同 `:1012` FunnelsEveryPostSpawnFailure…Fallback — **受影响**：源码级计数钉住 fail 闭包形态，重构不得增删失败返回点。
3. 同 `:487/:511/:528/:549/:566/:581` RefusesOn* 六条 — 不受影响（守卫在记录/平台校验层）；但 `:600` IncompleteRecords 夹具若依赖单记录字段需随 `RequestRecords` 键控补形态（受影响-夹具）。
4. 同 `:620` newAgentRunnerPIDReportRuntime — **受影响**：直构 runtime 结构体，字段映射化须同步改。
5. `runtime_agent_concurrent_e2e_test.go:72` ConcurrentStartOneWinner — 不受影响（A6 契约在 launcher 层）；多任务落地后它是 A6 守护线。
6. `runtime_supervised_test.go:316/:1231` FullChain(ExtendedScans) — **受影响**：supervised 同构重构回归面。
7. 同 `:1421` 直调 supervisedAssemblyStep — **受影响**：函数改名/签名变化即编译红，须同步。
8. 同 `:1107` 直构 &supervisedRuntime{} — **受影响**：字段映射化须同步。
9. 同 `:882/:901` 跨模式拒绝两条 — 不受影响：错误模式定义仍被 step 形状门拒绝（文案保留）。
10. `runtime_zero_model_e2e_test.go:66/:103/:160` — 不受影响：ExecuteNoopRun/localRuntime 路径不在改动面。
11. `cli_test.go:90/:361`、`stopper_test.go:24/:59` — 不受影响：noop CLI 路径。
12. `internal/chain/engine_zero_model_e2e_test.go:21/:82/:114/:162` — 不受影响：引擎层本就多任务。
13. `internal/adapters/agent_runner_run_test.go:50/:651` 等 wait/map、Capturer 相关 — 间接受影响：②逐任务复用；函数签名不变则绿。

实得字符数 = 4386
GLM-1: DONE
