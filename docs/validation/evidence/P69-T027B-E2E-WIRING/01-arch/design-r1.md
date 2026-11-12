# P69 `T027B-E2E-WIRING` ① 架构师方案书（design-r1）

载体：`runSubagent(agentName='deep-reasoner', model='DeepSeek V4 Pro (deepseek)')`（只读；最终消息逐字转录落盘）。
结论行：**DESIGN: DONE**。

## §1 现状核对（已实测 HEAD dcc02aa，仅列事实）

1. `internal/gates/readiness.go:159` `EvaluateT027BReadiness` 纯函数；`GoalProbe` 分支只查 `Permit != nil`（:166-173 注释明示：已 Spent 的 permit 仍产出 `ProbeReadiness`，取额度必须走 `Require`）。
2. `internal/gates/probe_permit.go`：`NewProbePermit`:126 / `LoadProbePermit`:140（缺 marker ⇒ Unused；有 marker ⇒ Spent）/ `Require`:210（O_EXCL 线性化）/ `Complete`:253（claimed→spent，幂等；unused 态返回 `ErrProbePermitInUse`）；哨兵 `Missing/InUse/Consumed`:66/:69/:71。
3. 装配入口 `internal/console/runtime_agent.go:126`：ctx 检查 → clock/ids/stub → Step 1 前置校验 → Step 2 `ensureAgentRunnerEventLogMarker` → … → engine.Run 驱动循环。`AgentRunnerAssemblyOptions` :79-124（无 permit 字段）。
4. `internal/console/runtime_supervised.go:267` `ExecuteSupervisedRun`：ctx → `supervisedAssemblyPlan` → RunRoot/VisibleSender 校验 → marker（:282-283）→ `newSupervisedEngine` → `PAUSED` 则 `resumeSupervisedRun`。`ExecuteSupervisedReturn`:310 独立入口。
5. **关键事实（起跑消息未覆盖）**：`ExecuteSupervisedRun` 是**可重入状态机**——同一 runRoot 会多次调用（park→confirm→open→return，见 `runtime_multitask_supervised_test.go:30-55` 4 连调）；门禁若无条件放在入口，第 2 次调用必遭 `ErrProbePermitConsumed`。续跑判别可行：`openConfirmationInteraction`（runtime_supervised.go:776）在首次 park 前把绑定请求写入 `SupervisedInteractionsLedgerPath(runRoot)`，故 **ledger 存在 = 非首启**。
6. CLI `internal/console/cli.go:274`：判定块 :328-346（互斥/无装配步骤/监督互斥/输入不完整拒绝 :342-346）；eventLog 已存在拒绝 :348-374；三路分流 :375-381。`cli.executeAgentRunnerRun`:454 构建 options（:544-564）；`cli.executeSupervisedRun` 在 `supervised_cmd.go:25`，options 由 `supervisedAssemblyOptions`:42 构建（`executeSupervisedReturn`:96 也复用它）。
7. 回归面实测（大于起跑消息口径）：`ExecuteAgentRunnerRun` 直调 **18 处**——`runtime_agent_e2e_test.go` 8 处 + 源码形态测试 1 处（:1027，计数 `return fail(`==8，本方案新增的 gate 早于 engine、不加 fail 调用、不加 summarizeProjection 返回，**不影响该测试**）、`runtime_multitask_driver_test.go` 5 处、`runtime_multitask_red_test.go` 4 处、`runtime_supervised_test.go`:884 1 处；`ExecuteSupervisedRun` 直调 **约 34 处**（`runtime_supervised_test.go` 约 26、multitask_red 5、multitask_supervised 3）。CLI 面仅 `cli_multitask_wiring_test.go`（args 统一经 `p51SevenAgentRunnerArgs`:140）与 `runtime_supervised_test.go`:873（互斥拒绝在 gate 前触发，不受影响）。
8. `TestAgentRunnerAssemblyRefusesSupervisedDefinition`（runtime_supervised_test.go:884）断言拒绝后 runRoot **零文件**；`newAssemblyFixture`（runtime_agent_e2e_test.go:157，options 字面量 :355）为 AgentRunner fixture 唯一注入点；`newSupervisedHarness`（runtime_supervised_test.go:102，options :120）为 supervised 唯一注入点。

## §2 接线设计

### 2.1 接线点清单

| # | 路径 | 增/改/删 | 内容 | 影响 |
|---|---|---|---|---|
| 1 | `internal/console/readiness_gate.go` | **增** | `requireProbePermit(ctx, p *gates.ProbePermit) (release func(), err)`：Evaluate(GoalProbe)→Require→返回 defer 用 release（内部调 Complete）。三件套单一来源 | 新文件 |
| 2 | `internal/console/cli.go` | 改 | ① `--probe-permit <path>` 旗标（:296 后）；② **gate 块**插在 `cfg.ToDefinition()` 之后、:375 三路分流之前：仅当 `len(supervised)>0 \|\| len(executable)>0` 时生效——缺旗标⇒fail-closed 拒绝；`LoadProbePermit(resolvePathFromCWD(...))` 失败⇒拒绝；`EvaluateT027BReadiness(GoalProbe)` 失败⇒拒绝；**CLI 不消费**（消费点在入口）；③ 把 permit 传入 `executeAgentRunnerRun`（:565 调用）与 `executeSupervisedRun`（:378 调用）；④ import gates | 新增拒绝路径，noop 分流不变 |
| 3 | `internal/console/runtime_agent.go` | 改 | ① `AgentRunnerAssemblyOptions` 增字段 `ProbePermit *gates.ProbePermit`；② 判定：stub 初始化后（~:141）、Step 1（:163）前；③ `Require`：Step 1 全部通过后、Step 2 marker（:173-178）前；④ `defer release()` 紧随 Require 成功 | 生产装配消费点 |
| 4 | `internal/console/runtime_supervised.go` | 改 | ① `SupervisedAssemblyOptions` 增同名字段；② **fresh 门**：`os.Stat(SupervisedInteractionsLedgerPath(options.RunRoot))` 为 `ErrNotExist` 才执行 判定+Require+defer Complete（非 NotExist 错误⇒fail-closed 拒绝）；位置在 VisibleSender 校验（:279）后、marker（:282）前 | 可重入安全 |
| 5 | `internal/console/supervised_cmd.go` | 改 | `executeSupervisedRun` 签名加 `permit` 参数并赋 `options.ProbePermit`；`supervisedAssemblyOptions` 与 `ExecuteSupervisedReturn` **不动**（归还路径 permit 留 nil，未消费路径不查） | CLI 面 |
| 6 | `internal/console/readiness_gate_test.go` | **增** | 正例 2 + 反例 1/2/3 共 5 条（见 §3） | 新测试 |
| 7 | 6 个既有测试文件 | 改 | 注入点收敛：`runtime_agent_e2e_test.go:355` fixture + `newTestProbePermit(t)`；`runtime_multitask_driver_test.go`（复用 fixture，0 处直改）；`runtime_multitask_red_test.go` 10 处裸 options 字面量补字段；`runtime_supervised_test.go` harness :120 + 3 处裸字面量（:884/:906/:965）；`cli_multitask_wiring_test.go` 的 `p51SevenAgentRunnerArgs` 追加 `--probe-permit <dir>\probe-permit.marker`（8 组测试一次覆盖） | 回归适配 |

### 2.2 调用顺序（Q4 结论）

**AgentRunner**：ctx → 判定（`Evaluate(GoalProbe, ReadinessInput{Permit: options.ProbePermit})`，失败⇒fail-closed、零痕迹）→ Step 1 前置校验 → `Require(ctx)`（失败⇒fail-closed、零 run-root 痕迹）→ `defer release()` → Step 2 marker → Step 3-8 照旧。
**Supervised**：ctx → plan → RunRoot/VisibleSender 校验 → **fresh 门**成立时 判定→Require→defer release → marker → engine（PAUSED 则 resume）。
**理由**：判定先于一切（纯函数，缺失 permit 即拒，零 IO）；Require 置于前置校验之后——配置/记录类错误**不烧一次性额度**（操作员修复后可原路径重试，marker 未落盘仍是 Unused）；Require 在 marker 之前——拒绝的装配对 runRoot **零写入**（保住 `runtime_supervised_test.go:884` 的零文件断言与 `TestAgentRunnerAssemblyRefusesOnIncompleteRecords` 的无 marker 断言）。defer 使成功与失败路径都不漏 Complete；N4 语义「已 Require 不退额度」自动成立。
**失败路径三情形**：① 判定失败（permit 缺失/旗标缺失/加载失败）⇒ 入口或 CLI 拒绝，零消费零痕迹；② permit 已被消费 ⇒ CLI 判定仍过（presence-only，P44 有意设计）但入口 `Require` 以 `ErrProbePermitConsumed` 拒绝、不双启动（O_EXCL 兜底跨进程）；③ 判定成功但 permit 未消费 ⇒ 只可能发生在「判定后、Require 前」的前置校验拒绝——该次起跑 fail-closed、未 spawn、额度完好，**不存在消费缺口路径**（spawn ⇒ Require 已成功 ⇒ defer 必 Complete）。`ExecuteSupervisedReturn` 不受门控（归还/续跑非新起跑；额度已在首启消耗，N4 一致）。

### 2.3 Q1/Q2/Q3 结论

- **Q1**：三处均 `GoalProbe`。理由：装配起跑即「一次 probe 起跑」，L1 由其 permit 定义；装配入口无 SealBinding/AT23 材料，L2/L3 在该点结构上不可达（`readiness.go:166-190`）。
- **Q2**：两 options 增 `ProbePermit *gates.ProbePermit` 指针字段（非 `ReadinessInput`——它是纯判定入参，装配选项持原始材料，判定点现场构造）。**D3 触发条件满足**：P44 报告 :156 处方（注入 `ReadinessInput`）被起跑消息口径取代，注入形态自此**定稿**为「装配选项绑定 permit 指针 + 判定点构造 input」，D3 可解除。此偏差在报告中如实登记即可（报告属验证文档，非契约）。
- **Q3**：新增 `--probe-permit <path>`，用 `LoadProbePermit`（生产不得每跑重签额度，marker 即额度）。缺旗标且路由到装配 ⇒ fail-closed（反例 3a）。noop 路径与 `ai check` 路径**本片不动**：② 步骤清单未含；noop 无装配无 spawn 无可准入的 probe run，`ai check` 是 P44 设计的另一 L1 消费接缝（`ai.go`），其 permit 生命周期（按检查/按跑）需独立设计，属后续切片。`--request-record` 单表与 `--request-records` 多表两条路由同受门控（gate 块在三路分流前，单点覆盖两路 + supervised 路）。

### 2.4 回退方案

撤销 5 个生产文件 hunk（cli.go 旗标+gate+两处调用点+import；runtime_agent.go 字段+gate 块+import；runtime_supervised.go 字段+fresh gate；supervised_cmd.go 签名+调用点）+ 删除 `readiness_gate.go`/`readiness_gate_test.go` + 恢复 6 个测试文件注入点（§2.1 #7）。**注意**：回退必须是编译级全量——若只撤生产 hunk 而留测试改动，`runtime_agent_e2e_test.go` 等引用 `ProbePermit` 字段的文件**编译失败**（首当其冲红：`runtime_agent_e2e_test.go`、`runtime_multitask_red_test.go`、`runtime_supervised_test.go`、`cli_multitask_wiring_test.go` 未知旗标、`readiness_gate_test.go` 全红）；完整回退至 dcc02aa 则全绿。

## §3 反例清单（≥3，含变异检验目标）

| # | 反例 | 断言 | 守卫（file:line） | 变异检验目标 |
|---|---|---|---|---|
| 1 | 判定失败⇒不进入装配 | 完整有效 options + `ProbePermit=nil` ⇒ err `Is(gates.ErrProbePermitMissing)`；run-root event log marker **不存在**；零 spawn（`requireAgentRunnerZeroSpawn`）；supervised 同型（harness+`ProbePermit=nil`） | 判定块：runtime_agent.go（~:141 新插）/ runtime_supervised.go（:279 后新插）/ cli.go gate 块（:375 前新插） | 删入口 gate 块（判定+Require 一起删）⇒ 装配前进到 Step 2 marker ⇒ marker 存在断言红（**平台无关判别**；单独删判定因 Require 兜底不红，故变异对象为整块） |
| 2 | permit 已被消费⇒拒绝、不双启动 | permit 先 Require+Complete 预消费，再跑完整 options ⇒ err `Is(gates.ErrProbePermitConsumed)`；marker 不存在；零 spawn | Require 调用点：runtime_agent.go（Step 1 后新插）/ runtime_supervised.go（fresh 门内新插） | 删 `Require` 调用（保留判定）⇒ presence-only 判定放行 ⇒ marker 出现 ⇒ 红 |
| 3 | CLI 未传必要参数/缺旗标⇒fail-closed 不静默跳过 | 3a：装配链 + 无 `--probe-permit`（记录路径故意指向不存在文件）⇒ 拒绝消息含 probe permit 而非记录读取错误；3b：有 permit 但缺 `--agent-cli` ⇒ 拒绝消息仍为 :342-346 既有措辞、permit marker 不存在（未烧额度）、不落 noop | 3a=cli.go 新 gate 块；3b=既有 :342-346 拒绝 + `--probe-permit` 解析（:296 后） | 删 CLI gate ⇒ 3a 错误变记录读取错误红；删 :342-346 ⇒ 3b 消息/分流改变红 |
| 4 | 既有 P44 gates 判定测试不变 | `internal/gates/**` 零 diff ⇒ `go test ./internal/gates/` 全绿 | 边界 1（diff 审查门禁） | 对 gates 包任何 diff 由硬边界禁止，diff 审查即守卫 |
| 5 | Complete 不漏 | 正例（stub 驱动完整链，平台无关）后：permit marker 文件存在；**同一 permit 对象二跑** err `Is(ErrProbePermitConsumed)`（而非 InUse） | `defer release()`（两入口） | 删 defer Complete ⇒ 二跑 sentinel 由 Consumed 变 InUse ⇒ 红 |

**测试侧 permit 形态（Q5）**：统一 `gates.NewProbePermit(id, filepath.Join(t.TempDir(), "probe-permit.marker"), 1)`（helper `newTestProbePermit`）——免前置磁盘状态、Unused 起跑、Require 在 t.TempDir 落 marker，不影响 runRoot 零写入断言；`LoadProbePermit` 仅 CLI 测试经 `--probe-permit` 天然使用（缺 marker ⇒ Unused，符合生产语义）。既有测试最小改动即 §2.1 #7 六个注入点；被新门改变断言的仅「裸 options 拒绝类」用例（错误由 ErrInvalidConfig 序列变为 permit 缺失——补字段后原哨兵恢复），其余 fixture 用例断言不变。

## §4 边界判定

1. **只改 `internal/console/**`**：PASS。生产改动仅 console 5 文件（新文件也在 console 内）；import `gates` 为公开导出符号；不改 gates/chain/adapters/tickets/guard 判定本体、schema、fixtures、workflows、go.mod；测试改动全部在 `internal/console/*_test.go`。
2. **CONTRACTS/v50-v55/gitee**：PASS。本设计不触碰 CONTRACTS（§4.2 末句时态为 ⑧ 事项，不在本片）；不改快照；无任何 commit/push 步骤。
3. **P44 缺陷/接口冲突核查**：**未触发**。GoalProbe presence-only + Require 强制配额是 P44 有意分层（readiness.go:159 注释），非缺陷；唯一超预期事实——`ExecuteSupervisedRun` 可重入——用 ledger 存在性 fresh 门在装配层解决，不改判定本体；CLI 侧 spent permit 判过、入口拒绝是消费点后置的必然结果，非冲突。

## §5 结论

**DESIGN: DONE**
