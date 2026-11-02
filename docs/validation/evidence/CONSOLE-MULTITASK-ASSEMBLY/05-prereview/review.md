# P50 `CONSOLE-MULTITASK-ASSEMBLY` · ⑤ 预审报告

## §1 审查范围与方法（自跑读数）

**审查对象（磁盘快照）**：`git status --porcelain` + `git diff --numstat` 于 HEAD `8bd8224`：契约 2 文件（已暂存，各 +1/−0）、`runtime_agent.go` +568/−173、`runtime_supervised.go` +459/−205、`agent_runner_intent.go` +90/−0、`runtime_agent_e2e_test.go` +8/−7、`runtime_supervised_test.go` +14/−6、4 个未跟踪测试/证据文件、`REMAINING_SLICES{,_EN}` +21/+21。与 ②a–d 回执逐项一致，无掩藏。
**必读**：design.md、red-baseline.md、glm-1、implementer-reports、CONTRACTS A6 段（CN:255 / EN:249）、两生产文件全文、5 个测试文件、`git show HEAD:` 旧 `*AssemblyStep` 对照。
**自跑验证**：`go test -count=1` 三包全绿（console 11.279s / adapters 42.798s / chain 3.031s）；`go vet` 三包 EXIT=0；`gofmt -l` 9 文件空输出；`return fail(` 计数=8；`REMAINING_SLICES.md` 首字节 239,187,191,35（UTF-8+BOM，终锻乱码仅为渲染）。

## §2 发现清单

**High**
- **H1 多任务下 fail 闭包按「上一 task」错配 requestId，可致非首 task 的孤儿 launch 不被止损**。`runtime_agent.go:319` 非 awaiting 错误一律 `fail(active.record.Request.RequestID, …)`，而 `active` 仅在 awaiting 解析成功后更新（`:325`）；任务 2 在 `engine.Run` 内先记录 launch receipt 再失败（OB-57 形态：receipt 先于 spawn 落盘）时，`active` 仍是 task 1 → 止损读 task 1 的 receipt，task 2 存活进程不 try-stop。`:290-295` 注释「the only launch that can have been spawned」为假（dispatch 端口入口 `observeDispatch` 于 `:370` 已观测 task 2）。判据：向 `RequestRecords` 双 task 夹具注入 task 2 dispatch 后失败场景，task 2 requestId 无止损调用。修法：非 awaiting 分支优先用 `assemblyRuntime.awaitingDispatch()` 解析 requestId（观测先于 receipt 落盘），仅无观测时回退 `active`。
- **H2 生产多任务 AgentRunner 驱动路径零行为覆盖**。grep `internal/console/*_test.go` 中 `RequestRecords|proveTaskStopped|agentRunnerAwaitingTask` 命中 0 条：`ExecuteAgentRunnerRun` 的多记录注册表分支、绑定表、任务间停机屏障均从未被任何测试执行；唯一多任务绿测 `TestAgentRunnerMultiTaskSerialOfflineHarness` 手工重写驱动循环、直驱 `chain.New`，绕过生产驱动（`runtime_agent.go:307-352`）。契约核心句（屏障先于后一 task 的 R/launch-intent）仅经代码走查、无行为证明。修法：补生产入口级双 task 测试（利用既有 `agentRunnerLaunchReceiptReader/Stopper` seam 或抽 dispatcher seam），并为 `proveTaskStopped` 补 report 存在/缺失/alive 三态单测；否则在 ④ 明确登记为后置片。

**Medium**
- **M1 `agentRunnerAssemblyRecords` 不校验 Attempt，错配记录在 marker 落盘后才失败**。`:988-1022` 只查 RunID/TaskID/StepID；`recordsByTask` 按记录自带 attempt 键控（`:221-236`），引擎新 run 恒 dispatch attempt 1 → `recordFor(taskID,1)` 在 marker 与引擎创建后报「no request record is bound」（`:678-687`），违反 red-baseline 钉死的「拒绝先于 marker、空摘要」性质。supervised 侧已 fail-closed（②c 偏差 4）。修法：`:1009` 起追加 `record.Request.Attempt != 1` 拒绝并补用例。
- **M2 supervised 绿测未钉死五事实在事件证据中的冻结位置序**。`runtime_multitask_supervised_test.go:103-138` 取 `InputEvidence[:5]` 仅验非空+跨 task 两两互异；位置序若回归仍假绿（postflight 只对账 `request.Facts`，不对账事件证据序）。修法：断言该 5 元组等于按各 task 工件重推的 `supervisedFactsList(derived)` 有序结果。

**Low**
- **L1 设计 §5 步骤 8 的 `ARCHITECTURE{,_EN}.md` 装配路径注记缺失**（git status 未见该二文件改动；仅 REMAINING_SLICES 已更新）。补注记或 ⑧ 显式登记延后。
- **L2 设计 §6 风险 ③「CLI 供给契约须 ② 前经主控确认，否则退回仅离线 harness」**：落地方案即回退形态，但主控裁定未见文字记录于 ②a–d。请在 ④/⑤c 补录，勿留默认证。

**Info**
- I1 `REMAINING_SLICES{,_EN}` P50 行仍标「步骤 0/进行中」（切片中段属预期，⑧ 更新）。
- I2 红→绿无法原地重放（②d 偏差 1）如实；建议 ⑧c 核对 red-baseline §B 真实性。
- I3 `ErrSupervisedStepRequired` 保留未用（②c 偏差 6）如实。
- I4 `AgentRunnerStepIntentRegistry` 已实现+三单测但生产未接（②d 偏差 4）：`agentRunnerIntentMux` 仍逐记录构造 preparer，行为等价，建议后片接用或删除其一。
- I5 `awaitingDispatch.dispatched` 标志不重置：串行+离线不可 resume 下无害，建议注释钉死该前提。

## §3 五个重点面逐面结论

1. **并发多任务上下文覆盖**：屏障次序正确（Submit 终局后、下一 `engine.Run` 前，`:344-350`）；跨 task 状态全部 `(taskID, attempt)` 键控（manifest/entries/capturer/workspace）；单 task 等价性经 HEAD 对照确认（文案、检查序 task→step→形状→policy 不变，零 task 专用新文案）；AgentRunner 离线不可 resume 与 A5 一致。但 H1、H2 使「正确性有代码依据、缺行为证明」。
2. **锁竞争与 A6 镜像**：adapters 零改动、无新锁、无并发 spawn；镜像/槽位/`refuseLiveMirrorProcess` 语义未触；`ConcurrentStartOneWinner` 守护线仍绿；A6「每 runRoot 唯一镜像」未被破坏（时刻级排他 + 串行轮换）。
3. **契约-代码对应**：CN/EN 新段与实现逐句可对（互异 requestId 强制 `:1012-1018`；停机证明取三者其一的「进程死亡」腿；`evidence[2..7]` 六事实由 harness 断言 `[requestHash, completionHash, 6 facts]` 严格位序；单日志全序两 harness 均断言 sequence 严格递增）。②b 映射表逐条属实。
4. **测试质量**：红→绿翻转符合 §D（负断言、零 task 文案、per-task step 规则均保留）；`wantFailureReturnPoints` 7→8 正确（独立计数=8、注释枚举一致、测试绿）；GLM-1 第 7 组 13 项处置经逐项核验（2/4/5/7/8/9 于 diff 确认，余经我独立全绿确认）如实；两 harness 断言均可证伪（迭代上限+COMPLETED 检查、双重 terminalize 防抖）。缺陷即 H2、M2、Windows R/C 跳过（已如实点名）。
5. **偏差点名复核**：②b 1-6 全真（唯缺 H1 之点名——fail 闭包注释的「唯一可 spawn launch」论证不成立，属 ②b 报告漏报）；②c 1-6 全真；②d 1-5 全真。未发现掩盖；R=0/D=0 与 numstat 互证。

## §4 总判定

核心架构裁决正确、契约修订与实现总体一致、红绿翻转与偏差登记诚实。阻塞项为 H1（止损错配，属跨 task 状态串用类，直接命中用户裁定审查面 1）与 H2（核心契约句无行为证明）。建议 ② 补修后交 ⑥ 复核：H1 一行级修法（观测优先解析 requestId），H2 可补最小生产路径测试或经主控显式裁决降为后置片并登记，其余按级处置。

SEVERITY=H2/M2/L2/INFO5
PRE-REVIEW: PASS WITH FIXES
