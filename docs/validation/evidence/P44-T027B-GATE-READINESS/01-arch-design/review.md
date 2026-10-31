# P44 ① 架构设计方案书（归档）

## 归档信息

| 项 | 值 |
| --- | --- |
| 切片 | `P44-T027B-GATE-READINESS` |
| 步骤 | ① 架构（只读角色） |
| 载体 | `runSubagent(agentName='deep-reasoner', model='DeepSeek V4 Pro (deepseek)')` |
| transcript 偏移 | jsonl 第 116795 行（assistant.message，7,196 字符） |
| 原始末行 | `DESIGN: DONE` |
| 归档动作 | 主控按"起跑消息即工件"同族纪律，将 ① 产出原文落盘（不改写正文），并追加 §八 修订注记 |

> 本文件是 **归档 + 注记** 复合件：§一～§七 为 ① 载体原文（仅做 Markdown 包裹以适配单文件归档，未改动任何结论、数值、命名与 file:line）；§八 为主控在 ⑤ 预审 r1 裁定后补记的派生决策，**标注来源与裁定方**，不属 ① 原文。

## 一 结论与范围

- 三级 readiness 纯函数与 `ProbePermit` 落 `internal/gates/`（新文件 `readiness.go`、`probe_permit.go`）；调用点仅 console 三装配入口 `runtime_agent.go:103`、`runtime_supervised.go:248`、`runtime.go:33` 与 `ai.go:39-90` 探针路径；判定零读写副作用。
- N4 两选项（一次性消费 / 窗口共享计数）取 **一次性消费**：P44 目标明示，且共享计数需改动 tickets 配额语义（边界外）。
- P38 蓝图把 `Permit.Require` 放在 `adapters/ai_probe_copilot.go:62-115`、`EvaluateT027BReadiness` 放在 `agent_runner_admission.go:42-116`——均属本片边界外；本片在 gates+console 等效落地，adapters 接线登记为边界外限制（P43 D3 同口径，`p43-state-seal.md:104`）。
- 依据：`review.md:4-8`（C7/N4）、`review-batch-b.md:40-46`；gates 哨兵 9 枚逐枚核实（见四）。

## 二 A 三级 readiness

| 级 | 名称 | 判定输入 | 可推进 | 禁止反例 |
|---|---|---|---|---|
| L1 | `ProbeReadiness`（探针就绪） | `input.Permit != nil` 且 availability 证据绑定候选 | 执行 F 可行性探针（`ai check` 路径 `ai.go:39-90`） | 探针结论推进验收（AT-23 自依赖） |
| L2 | `AcceptanceReadiness`（真实验收就绪） | `input.Seal != nil`（attempt/taskKey/sealHash）且平台前置能力通过；**不要求 AT-23** | 真实验收链 postflight→freeze→review→promotion（唯一 seal 入口 `engine.go:306-318`） | 仅有探针证据、缺 seal、缺事实闭包 |
| L3 | `ProductionReadiness`（生产发布就绪） | L2 的 `AcceptanceReadiness` 值 + `AT23.Hash` + 独立接受 | 生产发布 | 缺 AT-23、无 L2 |

类型隔离（不可混用的四项机制）：

1. sealed interface：`type Readiness interface{ tag() readinessTag }`，`readinessTag` 与三实现均 unexported ⇒ 包外无法伪造实现。
2. 私有构造函数：三类型只能由 `EvaluateT027BReadiness` 返回。
3. 方法集差异：`ProbeReadiness` 只有 `Permit()`，`AcceptanceReadiness` 只有 `Acceptance()`，`ProductionReadiness` 只有 `Release()`——取错级别编译红。
4. 构造依赖：`newProductionReadiness(acc AcceptanceReadiness, at23 AT23Signoff)` 签名强制 L3 必先有 L2 值（类型层绑定）。

签名与分支：

```go
func EvaluateT027BReadiness(goal Goal, input ReadinessInput) (Readiness, error)
// GoalProbe      → (*ProbeReadiness, nil)；缺 Permit → (nil, ErrProbePermitMissing)
// GoalAcceptance → (*AcceptanceReadiness, nil)；证据不足 → (nil, ErrReadinessProbeNotAcceptance|ErrReadinessEvidenceIncomplete)
// GoalProduction → input.Acceptance==nil → (nil, ErrReadinessAcceptanceMissing)；缺 AT23 → (nil, ErrReadinessProductionAT23Required)；否则 (*ProductionReadiness, nil)
```

## 三 B ProbePermit

类型（`gates/probe_permit.go`）：

```go
type ProbePermit struct {
    id      string
    maxReq  int64
    state   atomic.Uint32 // 0=Unused 1=Claimed 2=Spent
    path    string        // 持久标记文件（O_EXCL no-replace）
}
```

状态机：`Unused → Claimed → Spent`（terminal）。

- `Unused→Claimed`：`Require(ctx)` 内 CAS(0→1)——**线性化点**；随后写标记文件（`O_EXCL` 新建，先例 P43 D2 `seal.go` 的 O_EXCL 等价 no-replace，`p43-state-seal.md:103`）。标记写失败 → 直接转 `Spent`（额度作废，fail-closed）。
- `Claimed→Spent`：探针结束（成功/失败/超时任一）由 console 调 `Complete()` CAS(1→2)；同进程内未 Complete 前，Claimed 态的再次 `Require` 返回拒绝。
- **无任何转移回到 Unused** ⇒ 失败、超时、重启不得重置已消费额度。

并发「严格一次」：**用 `sync/atomic` CAS 而非互斥**。理由：permit 只有一个被争用的状态字、无多字段不变量；互斥是给多字段账本用的——`CostLedger.mu sync.RWMutex`（`cost.go:202-203`）保护 reservation 映射与汇总，`DurableCostLedger.mu`+锁序+跨进程文件锁（`durable_cost_ledger.go:101/:316`）。跨进程的线性化由**标记文件 O_EXCL** 承担（与 durable 锁同族：`durable_cost_lock_unix.go:17-38` flock、`windows.go:24` LockFileEx）：即使两进程各自 CAS 成功，O_EXCL 只允许一份标记存在，败者转 Spent。仓库无 atomic 先例（未确证），由 -race 竞态测试补证。

重启不重置的落点：`LoadProbePermit(path)`——标记文件存在 ⇒ 重建为 `Spent`；不存在 ⇒ `Unused`（新签发 permit 用新 id/path，旧 permit 不得重发）。失败/超时落点：探针返回错误或 ctx 超时**不**回滚 Claimed，console 必须 `Complete()` 至 Spent；若 console 崩溃于 Complete 前，标记已在盘 ⇒ 重启即 Spent。消耗模型：`Require` 一次性把 `maxReq` 额度全部预留给本次调用，调用后即失效（对应 N4「每次调用都可重新获得额度」的封堵）。

## 四 C 错误码与哨兵

新哨兵（7 枚，均 `internal/gates/readiness.go`/`probe_permit.go`）：

| 哨兵 | 返回分支 | 包装方式 |
|---|---|---|
| `ErrProbePermitMissing` | GoalProbe 且 input.Permit==nil | `fmt.Errorf("%w: …", …)` |
| `ErrProbePermitInUse` | `Require` 于 Claimed 态（并发第二人） | 同上 |
| `ErrProbePermitConsumed` | `Require` 于 Spent 态；失败/超时/重启后再 `Require` | 同上 |
| `ErrReadinessProbeNotAcceptance` | GoalAcceptance 输入仅探针证据 | 同上 |
| `ErrReadinessEvidenceIncomplete` | GoalAcceptance 缺 Seal 绑定/事实闭包 | 同上 |
| `ErrReadinessAcceptanceMissing` | GoalProduction 且 Acceptance==nil | 同上 |
| `ErrReadinessProductionAT23Required` | GoalProduction 且 AT23.Hash 空 | 同上 |

- 探针**执行期**失败/超时返回底层错误本身，不包 permit 哨兵——哨兵只守准入，与既有 `ErrSupervised*`（`runtime_supervised.go:150-199`）「门内/门外交错不混用」一致。
- 语义去重：既有 gates 9 枚 = `effects.go:18`、`policy.go:15-18`、`runner.go:18-20`、`scanner.go:9`（hook/effect/secret 域）；console 导出哨兵 `ai.go:15`、`config.go:29`、`handoff.go:14-15`、`runtime_supervised.go:150-199`（实测 19 枚，与批量值 18 差 1，口径未确证）。新增一律 `ErrReadiness*`/`ErrProbePermit*` 前缀，与既有全部前缀（`ErrInvalid*`/`ErrEnvironment*`/`ErrArtifact*`/`ErrSecret*`/`ErrSupervised*`/`ErrAIChannel*`/`ErrSecureInput*`/`ErrGate*`/`ErrScan*`/`ErrPolicy*`/`ErrWorkspace*`）零交集。
- 哨兵仅新增、不改既有；`errors.Is` 兼容性由 `%w` 链保持（先例 `policy.go:98`、`budget.go:47`）。

## 五 D 与 P43 的接缝

只读引用三点（不改 chain）：

1. `chain.LoadTaskSeal`（`seal.go`，`p43-state-seal.md:90`）：console 装配层（`runtime_agent.go:103` 附近）读取 seal 哈希与 `ContextKey`，填 `SealBinding{Attempt, TaskKey, SealHash}`；gates 判定函数只吃哈希、不做 IO。
2. `chain.ContextKey`/`TaskContextKey`（`attempt.go`，`p43-state-seal.md:92`）：L2 要求 SealBinding.TaskKey 与重派生值相等，等价复刻 `acceptFromExistingSeal` 的 stale 比对（`engine.go:349-353`）在装配前端的影子校验。
3. `Options.Attempts` 注入点（`attempt.go`）只读消费：readiness 不参与 attempt 解析，仅采信结果；当前缺省回退投影 attempt=1（D3 边界，`p43-state-seal.md:104`），attempt>1 的 L2 判定在装配层 fail-closed 不可达。
- gates 已可 import chain（`adapter.go:1-4`），但设计上**判定函数保持纯函数**（呼应 P38 `TestReadinessDoesNotCreateGateCycle` 零副作用断言，`review-batch-b.md:43`）。

## 六 E 红测清单

| 测试名 | 断言要点 | 未实现时应失败形态 |
|---|---|---|
| `TestReadinessProbeConclusionCannotAdvanceAcceptance` | GoalAcceptance+仅 availability → `errors.Is(ErrReadinessProbeNotAcceptance)` | 编译红 `undefined: EvaluateT027BReadiness`（新符号） |
| `TestReadinessLevelsCannotBeMixed` | 包外无法实现 `Readiness`；`ProbeReadiness` 上取 `Acceptance()` 编译红 | 同左（sealed 接口缺失即编译红） |
| `TestProbePermitConcurrentRequireSingleWinner` | 40 goroutine 同 permit：恰 1 次 nil、39 次 `ErrProbePermitInUse/Consumed`（先例 `cost_test.go:239`） | 赢家数 >1 |
| `TestProbePermitFailureDoesNotResetQuota` | Require ok→探针失败→再 Require = `ErrProbePermitConsumed` | 第二次 nil |
| `TestProbePermitTimeoutDoesNotResetQuota` | Require ok→ctx 超时→再 Require = Consumed | 可重入 |
| `TestProbePermitRestartDoesNotResetQuota` | 消费→`LoadProbePermit` 重建→Require = Consumed（双进程先例 `durable_cost_ledger_test.go:270/:688`） | 重启后 nil |
| `TestReadinessProductionRequiresAT23` | GoalProduction+有 L2+无 AT23 → `ErrReadinessProductionAT23Required` | 放行 |
| `TestReadinessHasNoSideEffects` | 判定前后工作目录/账本字节级不变 | 判定产生写 |

## 七 F 风险与未确证项

1. **装配回归面**：`runtime_agent_e2e_test.go:462-603` 共 12 处调用 `ExecuteAgentRunnerRun` 的装配用例需过 L2 判定；其 DECLARED 记录装配若缺 SealBinding 会 fail-closed 红。DR-9（`runtime_agent_e2e_test.go:464`）证明该测试对装配入口敏感 ⇒ ③ 红测前必须先定 `AgentRunnerAssemblyOptions` 注入 ReadinessInput 的形态。
2. **gates 4 判定点**（`adapter.go` ChainPort.RunHook、`policy.go` `requireCapabilities`、`runner.go` Run、`executor.go` Execute）只叠加判定、不改 hook 执行与 EffectClassification（`effects.go:45` 恢复动作闭集），避免效果分类漂移。
3. 未确证①：console 哨兵计数口径（实测 19 vs 批量 18），去重结论不受影响。
4. 未确证②：AT-23「签收」无权威工件（S1-REACH：真实候选腿不可达）⇒ L3 的 AT23 只能取声明式哈希，自证风险需绑定不可伪造证据摘要，否则 L3 名存实亡。
5. 未确证③：permit 若改走 `DurableCostLedger`（`durable_cost_ledger.go:190` Reserve）需 run root+锁文件，`ai check` 单发 CLI 无 run 上下文 ⇒ 本方案取 O_EXCL 标记，未验证其与 G 门禁/编码约定无冲突（新建 `.json` 需无 BOM+LF，见 copilot-instructions 硬规则）。
6. 未确证④：P43 报告残余（`p43-state-seal.md:143-146`）——装配未注入 `Attempts`、supervised 工作区固定 attempt-1（D4）——L2 的 attempt>1 路径不可达属 fail-closed，非本片可解。

---

## 八 修订注记（⑤ 预审 r1 裁定后补记，**非 ① 原文**）

> 记入依据：⑤ 预审 r1 判定 `PRE-REVIEW: PASS WITH FIXES`，其 B1 指出"① 方案未归档 ⇒ 无法对文核验，规格漂移必须留痕"。本节把 ① 原文与落地实现之间的三处差异逐条固化，供 ⑥/⑦/⑧ 复核。

### R-1（派生决策）`GoalProduction` 缺 L2 材料的判据

- **① 原文**（§二）：`GoalProduction → input.Acceptance==nil → (nil, ErrReadinessAcceptanceMissing)`，与探针证据无关。
- **② 派单契约**（主控在 ② 输入中改写为 `Seal` 字段）：写为"`Seal == nil` 且 `ProbeEvidence` 为空 ⇒ `ErrReadinessAcceptanceMissing`"。
- **落地实现**（`readiness.go:177-186`）：仅 `Seal == nil` ⇒ `ErrReadinessAcceptanceMissing`，**不看** `ProbeEvidence`。
- **裁定**：⑤ 预审 r1 判定"② 的语义正确、不构成 C7 漏洞"——L3 的不变量是"不得缺 L2 材料"，探针证据存在与否不改变"L2 缺失"事实；若按合取式，带探针证据的输入会落到 `ErrReadinessProductionAT23Required`，反而给出误导性诊断。
- **影响**：行为等价于 ① 原文（同为 fail-closed），仅诊断哨兵不同；③ 的 `productionWithoutAcceptance` 用例（`readiness_test.go`）反向钉住该语义。
- **复核要求**：⑥ 须确认"探针证据 + 缺 seal"不会产出 L3 值；⑦ 须确认该派生决策不构成"探针结论推进验收"的通道。

### R-2（派生决策）`ReadinessInput` 字段由 chain 类型改为哈希串

- **① 原文**（§二·四）：`newProductionReadiness(acc AcceptanceReadiness, at23 AT23Signoff)` / `input.Acceptance==nil`，即 L3 以 **L2 值**为输入。
- **落地实现**：`ReadinessInput{Permit *ProbePermit, Seal *SealBinding, AT23Hash string, ProbeEvidence []string, EvidenceRootHash string}`，L3 以 `Seal` 绑定 + `AT23Hash` 为输入。
- **理由**：①§五 要求"判定函数保持纯函数、只吃哈希不做 IO"，若入参为 `AcceptanceReadiness` 值则必须先在包内构造该值，等于把"谁有权构造 L2"变成 gates 内部问题；改为哈希串后，纯函数性由签名直接保证（`TestReadinessHasNoSideEffects` 钉住）。
- **代价**：① 的"构造依赖（L3 签名强制先有 L2 值）"这一项类型层绑定**未落地**，改为运行时判据（`Seal == nil`）。此项在 ⑥/⑦ 须按"类型层隔离减弱"复核。
- **复核要求**：⑥ 须确认 `Readiness` sealed interface + 三类型方法集差异仍阻止跨层取用；若判定"类型层绑定缺失构成 C7 缺口"，须在本片内补 `ProductionReadiness` 的构造参数校验或在后续切片修复。

### R-3（边界登记）console 接线**未在本片落地**

- **① 原文**（§一）：调用点为 console 三装配入口（`runtime_agent.go:103`、`runtime_supervised.go:248`、`runtime.go:33`）与 `ai.go:39-90` 探针路径。
- **实际**：本片只落 `internal/gates/**` 判定实现与 15 例测试，**console 零改动**（⑤ B2 独立核实"全仓 0 生产调用点"）。
- **登记口径**：C7/N4 **不得宣称端到端闭合**；接线登记为后续切片输入（触发条件：`AgentRunnerAssemblyOptions` 注入 `ReadinessInput` 形态定稿，见 §七 风险 1）。
- **理由**：①§七 风险 1 已指出 `runtime_agent_e2e_test.go:462-603` 共 12 处装配用例对 L2 判定敏感（DR-9 前例）；在未定 `AgentRunnerAssemblyOptions` 注入形态前接线属高风险改动，且会突破"本片可验证范围"。

### R-4（环境限制）`-race` 未在本机执行

- 本机 `CGO_ENABLED=0` 且无可用 cgo 工具链（无 gcc/mingw/WSL/docker；clang 路径含空格 + MSVC link.exe 组合失败），`go test -race` 无法编译。
- 替代证据：`-count=50 -cpu 1,2,4` 于并发用例、`-count=20` 于全部 P44 用例，均绿；③ 增补轮另有 7 项变异验证（M1–M6 + M1M2 合并）证明断言可证伪。
- CI 的 Race 腿当前覆盖 `internal/adapters`/`internal/chain`（不含 `internal/gates`）⇒ 本片并发结论的 race 证据**缺失**，登记为已知限制；扩展 Race 腿属 `ci.yml` 变更，超本片边界。
