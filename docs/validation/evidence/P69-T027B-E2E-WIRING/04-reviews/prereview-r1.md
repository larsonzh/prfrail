# P69 ⑤ 预审 r1 归档（逐字原文）

载体：`runSubagent(agentName='deep-reasoner', model='DeepSeek V4 Pro (deepseek)')`（只读）。结论行：**PRE-REVIEW: PASS**（1 Low + 7 Info，无 Medium+）。
以下为预审最终消息逐字转录。

已完整读取全部审查输入包（方案书、diff 全文、门禁读数、变异证据、契约 §4.2 原句）并逐条核对工作区实现。以下为 ⑤ 预审结论。

## 必核面逐项结论

**1. 接线点一致性 —— 已核，一致**
- 据主控落盘的 `tmp/P69-T027B-E2E-WIRING/05-prereview/diff.txt`，design-r1 §2.1 七项接线点全部命中且无清单外生产改动：
  - #1 `readiness_gate.go` 新建，`requireProbePermit` = Evaluate(GoalProbe)→Require→release（内调 Complete）三件套单一来源（diff 中 22 行新文件）；
  - #2 `cli.go`：`--probe-permit` 旗标（agentVersion 后、`set.Parse` 前）✓；gate 块位于 eventLog 已存在检查与 `cfg.ToDefinition()` **之后**、三路分流**之前**，条件 `len(supervised)>0 || len(executable)>0`，缺旗标/解析失败/加载失败/Evaluate 失败均 `writeCommandError` 拒绝，**CLI 不消费**（无任何 Require 调用）✓；`gatePermit` 传入两装配调用点 ✓；import gates ✓；
  - #3 `runtime_agent.go`：options 增 `ProbePermit` 字段；Step 0 判定在 stub 初始化后、Step 1 前；`requireProbePermit`（Require）在 Step 1 全部校验后、Step 2 marker 前；`defer releaseProbePermit()` 紧随其后（工作区 `runtime_agent.go:143-195` 与 diff 一致）；
  - #4 `runtime_supervised.go`：同名字段；fresh 门 `os.Stat(SupervisedInteractionsLedgerPath(...))`，`ErrNotExist` 才判定+Require+defer，非 NotExist 错误 fail-closed；位置在 VisibleSender 校验后、marker 前（工作区 `runtime_supervised.go:284-305` 与 diff 一致）；
  - #5 `supervised_cmd.go`：签名加 `permit` 并赋 `options.ProbePermit`；`supervisedAssemblyOptions`/`ExecuteSupervisedReturn` 零改动 ✓；
  - #6/#7 测试见必核面 2/6。
- 无清单外改动：据 `05-prereview/numstat-only.txt`，21 个文件与 diff 清单完全一致。

**2. 反例成立性 —— 已核，成立且自洽**
- 据 `docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt`（主控落盘）与 diff 中 `readiness_gate_test.go` 全文：
  - 实落 7 条（正例 1 + 反例 1a/1b/2a/2b/3a/3b），design-r1 §3 五条反例的"supervised 同型"拆分所致（见发现表 L2 计数注记）；
  - 每条断言均可证伪：1a/1b 哨兵+marker 断言（M1 删整块 ⇒ 红在 :113）、2a/2b Consumed 非 InUse+marker+zero-spawn（M2 删 Require ⇒ 红在 :152）、3a stderr 子串（M3 删 CLI gate ⇒ 红在 :200）、3b 既有措辞+marker 不烧（M4 删 :342-346 拒绝 ⇒ 红在 :220）、正例二跑 Consumed 非 InUse（M5 defer→`_=` ⇒ 红在 :98，InUse 泄漏）；
  - M1–M5 红行号与测试内 t.Fatalf 行逐一对上（diff 可数），守卫=变异对象一一对应，恢复 sha256 一致且恢复后绿；M2 首轮文件锁重试已如实登记。自洽 ✓。

**3. 越界审计 —— 已核，零越界**
- numstat-only.txt 全部路径 ∈ {`docs/ADR_REGISTER{,_EN}.md`, `docs/validation/evidence/P69-.../**`, `internal/console/**`}；对 `internal/gates`、`chain`、`adapters`、`tickets`、`guard`、`schemas`、`testdata`、`.github`、`go.mod` 均为零改动；
- `internal/console/**` 内只 import gates **导出符号**（`EvaluateT027BReadiness`/`ReadinessInput`/`GoalProbe`/`ProbePermit`/`LoadProbePermit`/`NewProbePermit`/哨兵；测试另复用 adapters 既有导出），未触及任何判定本体（`readiness.go`/`probe_permit.go` 零 diff）。

**4. 失败路径三情形 —— 已核，fail-closed 成立，无新 fail-open**
- ① 判定失败（缺旗标/加载失败/nil permit）⇒ CLI 或入口拒绝，零消费零痕迹（1a/1b/3a 断言 + M1/M3 证实）；
- ② permit 已被消费 ⇒ CLI 判定仍过（`readiness.go:166-173` presence-only 是 P44 有意设计），入口 `Require` 以 `ErrProbePermitConsumed` 拒绝、不双启动（2a/2b + M2）；跨进程 O_EXCL 兜底由 `probe_permit.go` 本体保证（未改动）；
- ③ 判定成功但 permit 未消费 ⇒ 仅 AgentRunner 的"判定后、Require 前"Step 1 拒绝，未 spawn、额度完好；supervised 判定+Require 原子无间隙；spawn ⇒ Require 已成功 ⇒ defer 必 Complete，无消费缺口；
- 特别审三项：fresh 门判别可靠——ledger 写入点 `parkSupervisedStep→openConfirmationInteraction→appendOperatorInteractionRecord`（工作区 `runtime_supervised.go:707/799`）在 fresh 门 Require 成功**之后**，故"ledger 存在 = 首启已过 Require"成立，重入跳过不构成 fail-open；崩溃窗口（Require 后、写 ledger 前）重启 ⇒ marker 已盘 ⇒ `Load`=Spent ⇒ Require Consumed ⇒ fail-closed，符合"失败/超时/重启不重置"；CLI 只判定不消费（diff 中无 Require 调用）与入口 Require 时序无竞态；defer Complete 为函数级注册，覆盖成功/失败/panic 全部路径。

**5. Unix 零影响 —— 已核，充分**
- diff 新增代码仅 `os.Stat`/`errors.Is`/`strings.TrimSpace`/`fmt.Errorf`/既有 `resolvePathFromCWD`，无 GOOS 分支、无 build tag、无平台文件；
- 据 `06-gates/p4-05-linux-build.txt` 与 `p4-06-linux-vet.txt`（均空输出）+ 无被漏审的平台分支 ⇒ 构成本机充分证据；Linux 运行时行为由 CI ubuntu 腿兜底（⑨ 事项）。

**6. 回归面 —— 已核，只补注入点、未改断言逻辑**
- `runtime_agent_e2e_test.go`（+9/-0）：import + fixture 字段 + `newTestProbePermit` helper，断言零改动；
- `runtime_multitask_red_test.go`（+10/-10）：恰 10 处裸 options 字面量补 `ProbePermit` 字段（与 design-r1 #7 的 10 处吻合），断言零改动；
- `runtime_supervised_test.go`（+12/-6）：harness options + 3 处裸字面量（:884/:906/:965 对应）补字段及 gofmt 对齐，断言零改动；
- `cli_multitask_wiring_test.go`（+1/-0）：`p51SevenAgentRunnerArgs` 追加旗标，断言零改动；
- `runtime_multitask_driver_test.go` 0 处直改（复用 fixture）✓；`runtime_supervised_test.go:884` 零文件断言不破——新 permit 的 marker 落在 `t.TempDir()` 独立目录，且该测试拒绝点在 Step 1（Require 前）。

## 发现表

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| Low | `internal/console/readiness_gate.go:20` | CONTRACTS §4.2："`Complete` 以 `CAS(1→2)` 落终态，对已 `Spent` 幂等，对从未 `Require` 的 `Unused` 返回 `ErrProbePermitInUse`" | release 闭包 `_ = permit.Complete()` 吞掉落终态失败，使 Complete 错误不可观测（单进程单次调用下失败不可达，无 fail-open）。最小修复：在闭包内捕获错误并走内部日志记录，或补注释论证该失败不可达（现注释只声明"completes exactly once"） |
| Info | design-r1 §2.1 #6 vs 实落 | design-r1 §2.1 #6："正例 2 + 反例 1/2/3 共 5 条" | 实落 7 条（正例 1 + 反例 1a/1b/2a/2b/3a/3b），系 §3 反例 1/2 的"supervised 同型"拆分；在 ⑧ 报告中如实登记计数即可，不改代码 |
| Info | `internal/console/runtime_agent.go:147,191` | design-r1 §2.2："判定先于一切……Require 置于前置校验之后" | AgentRunner 路径 Evaluate 执行两次（Step 0 + requireProbePermit 内），纯函数冗余、无副作用，与 design-r1 明文一致；无需修复 |
| Info | `runtime_supervised.go:284-305`（fresh 门） | design-r1 §3 反例表（变异检验目标仅列 5 行） | M1–M5 无 fresh 门专项变异；其守卫由既有 `TestSupervisedMultiTaskSerialHarness` 4 连调承担（fresh 门删除 ⇒ 第二次 ExecuteSupervisedRun 撞 `ErrProbePermitConsumed` ⇒ 该测试红，且据 `p4-04-test.txt` console 包 ok）。可接受；如需可在报告登记 |
| Info | design-r1 §3 反例 1 断言列 | design-r1 §3："run-root event log marker 不存在；零 spawn（requireAgentRunnerZeroSpawn）" | 反例 1a/1b 未含 zero-spawn 断言（2a 含）；marker 断言已足以证伪（M1 证实），且 marker 与 spawn 同因（gate 在 Step 2 前）。无需修复 |
| Info | `readiness_gate_test.go` 3b | design-r1 §3 反例 3 | 3b 无法区分 gate 块与 :342-346 拒绝的相对顺序（两种顺序都不烧额度、行为无差异）；M4 证实其守卫对象为 :342-346 的存在性。可接受 |
| Info | `docs/CONTRACTS.md` §4.2 末句 | CONTRACTS §4.2："本切片只落判定实现与测试，`internal/console/**` 未被改动，故 C7 / N4 的端到端接线尚未完成" | 时态更新为 ⑧ 回写事项（起跑消息【五】第 10 项），不在本片代码范围；预审不判缺陷，提示 ⑧ 勿遗漏 |
| Info | `runtime_supervised.go` fresh 门 | design-r1 §2.2："fresh 门成立时 判定→Require→defer release" | 重入（freshStart=false）不经 Evaluate、permit 可缺席；由"ledger 写时点在 Require 后"保证首启已消费额度，design-r1 §4.3 已裁定。无需修复 |

## 其余已核项

- 门禁读数（据 `06-gates/`）：`p4-01-gofmt.txt` 空；`p4-04-test.txt` 15 包全 ok（console 32.412s、gates 7.907s）；`p4-05-linux-build.txt`/`p4-06-linux-vet.txt` 空；`p4-07-selftest.txt` 末行 `SELFTEST: PASS 737/737`；`p4-08-index.txt`/`p4-09-tree.txt` 均 `TOTAL_FAIL=0`、`base gofmt/build/vet/test :: empty/0/0/ok=15 FAIL=0`、G4a 21 文件 OK、G4b unknown=[]、G8-a/b OK、G5(f) index=worktree（暂存与工作区一致）；`p4-10-ci.txt` 为 HEAD^1..HEAD 退化范围（changed=2），本片变更集覆盖由 index/tree 给出（仓库既有口径）。
- 测试构造核对：3a/3b 的 `dropCLIFlag` + `--request-record` 追加与 `p51SevenAgentRunnerArgs`（七项非 request 旗标）组合正确；正例测试平台感知分支复用既有 `agentRunnerPublicationDurability`/`buildAssemblyStub`/`seedChildPIDReport` fixture，非新逻辑；Windows 上耐久拒绝路径仍验证额度恰消耗一次与二跑 Consumed（M5 证实两端可证伪）。
- 契约一致性：`EvaluateT027BReadiness` 纯函数与 presence-only 语义由本体注释背书（`readiness.go:159/166-173`），装配未绕过；`Require` 线性化点与 `LoadProbePermit` 语义均由本体保证且零 diff。
- 变异证据自洽：M1/M2/M5 恢复 sha `1100E2DE…C3F336`、M3/M4 `5C8D6B15…620EE2B`（按文件对应一致），红行号与测试断言行逐一对上。

PRE-REVIEW: PASS
