## C5
- 读取 [agent_runner_intent.go:20-79](internal/adapters/agent_runner_intent.go#L20) `NewAgentRunnerStepIntentPreparer/PrepareStepIntent`、[runtime_agent.go:103-349](internal/console/runtime_agent.go#L103) `ExecuteAgentRunnerRun` → 新增 `ResolveAgentRunnerExecution`，按 `K=(run,task,step,attempt,requestId,requestHash)` 选择不可变上下文，逐项核验其余摘要与工作区身份；显式 Agent 目标缺绑定即拒绝，普通 hook 保持原路由，后续任务只在前一候选接受后绑定父快照，同一 Engine 循环处理派发与终局 → `execution, err := preparer.ResolveAgentRunnerExecution(ctx, request); if err != nil { return chain.StepExecutionIntent{}, err }; return execution.Intent, nil`
- 读取 [agent_runner_run.go:56-162](internal/adapters/agent_runner_run.go#L56) `Run`、[agent_runner_postflight.go:616-713](internal/adapters/agent_runner_postflight.go#L616) `WaitAgentRunnerDispatched`、[agent_runner_pinned_cli.go:442-505](internal/adapters/agent_runner_pinned_cli.go#L442) `WaitAgentRunnerProcess` → 两条运行入口共用执行上下文的 `outcomeOnce/outcomeDone`，仅原所有者启动等待、收尾和采集；截止点固定为原 `startedAt+Timeout`，重复调用只等待结果或取消自身等待，返回深拷贝，缓存不随进程句柄退休而丢失；无 intent 且不能证明原所有者连续性时阻断 → `execution.outcomeOnce.Do(func() { go execution.collectOnce() }); return execution.awaitOutcome(ctx)`
- 读取 [agent_runner_terminal_publisher.go:61-296](internal/adapters/agent_runner_terminal_publisher.go#L61) `PublishTerminal/convergeTerminal` → 新增显式 `RecoverTerminal(ctx)`，装配先查原 intent，再决定恢复或等待；恢复校验请求绑定及批次 A 的 `Artifacts.Verify` 后，仅用原完成体、事实和结算计划收敛，禁止重新采集、生成时间或从账本反推完成体；正常发布仍拒绝异终局，私有收敛方法透传 ctx → `return publisher.convergeTerminal(ctx, intent, closure, foundClosure, intent.Intent.Completion, intent.Intent)`
- 读取 [engine.go:109-215](internal/chain/engine.go#L109) `Run/runTask`、[runtime_agent.go:424-430](internal/console/runtime_agent.go#L424) `Materialize` → `Run` 在暂停恢复转移前先核验待评审资格，`runTask` 仅在 PRECHECK 调用幂等 Materialize，其余执行态使用新增 `WorkspacePort.Reopen`；首次物化从接受父快照建立工作区，重开只核验原 run/task/attempt/parent/工作区绑定，缺失不重建，未认证 REVIEW_PENDING 零写拒绝 → `workspace, err := engine.options.Workspaces.Reopen(ctx, engine.options.RunID, task.ID, 1, parent)`
- 读取 [runtime_agent.go:439-464](internal/console/runtime_agent.go#L439) `Accept` → 按原任务绑定加载已持久冻结的候选及证据根，核对工作区内容后返回原引用；不存在候选时才首次冻结，漂移不得重新冻结后沿用旧批准，已认证 REVIEW_PENDING 不重跑 postflight → `if found { return runtime.verifyFrozenCandidate(ctx, frozen, workspace) }`
- 新增/修改测试：[agent_runner_postflight_test.go](internal/adapters/agent_runner_postflight_test.go#L1) `TestWaitReentrySharesOriginalOutcome` → 输入同 K 的两个并发等待、后加入者取消及延迟重入；期望原执行继续并收敛。
  断言：底层运行等待与证据采集各一次，第二次调用不延长原截止点。
  断言：正常取得结果的调用完成时间、完成摘要与六事实完全相同，修改返回值不污染缓存。
- 新增/修改测试：[agent_runner_terminal_publisher_test.go](internal/adapters/agent_runner_terminal_publisher_test.go#L1) `TestRecoverUsesOnlyStoredIntent` → 输入分别在 intent 后、结算后崩溃的记录，重建 publisher 并改变时钟；期望恢复原决定。
  断言：原完成体、六事实、结算键和 intent 摘要不变，采集与 spawn 均为零。
  断言：持久结算条目仅一条；异完成体发布返回冲突，无 intent、无连续性证明的重入不补造终局。
- 新增/修改测试：[engine_test.go](internal/chain/engine_test.go#L1) `TestReviewReentryPreservesWorkspace` → 输入已暂停的 REVIEW_PENDING，分别提供无资格证据、合格且未变内容、合格但漂移内容。
  断言：无资格时返回 `ErrUnqualifiedReview`，事件追加、Materialize、工作区写入均为零。
  断言：合格重入只 Reopen，不重跑 postflight；漂移时不调用 Reviewer/Publisher，未变内容使用原候选引用。
- 新增/修改测试：[runtime_agent_e2e_test.go](internal/console/runtime_agent_e2e_test.go#L1) `TestTaskBindingsFollowAcceptedParent` → 输入两任务及声明 hook，T1 接受为 P1，T2 分别使用正确绑定、缺绑定、旧父快照绑定。
  断言：正确路径中 T2 的请求、workspace、reservation 与端口上下文均独立且绑定 P1。
  断言：缺绑定或父快照漂移时 T2 零 spawn；hook 不进入 Agent 路由。

## C6
- 读取 [cost.go:245-340](internal/tickets/cost.go#L245) `DecodeCostLedgerRecord/LoadCostLedger`、[cost.go:426-556](internal/tickets/cost.go#L426) `Snapshot/Reserve/Settle` → 保留 CostLedger 纯规则模型，新增持久包装 `DurableCostLedger`：宿主按共享 scope 固定唯一规范账本路径和稳定旁路 OS 锁，锁内重读校验、执行规则、校验旧条目前缀未变、写完整快照、Sync/Close、原子提交及平台所需目录屏障后才返回私有字段的 `DurableAck`；幂等重试也须重新确认耐久，确认失败保留可见结果且不回滚 → `record, err := DecodeCostLedgerRecord(data); if err != nil { return err }; ledger, err := LoadCostLedger(record)`
- 读取 [agent_runner_admission.go:42-116](internal/adapters/agent_runner_admission.go#L42) `AdmitAgentRunner`、[agent_runner_admission.go:199-210](internal/adapters/agent_runner_admission.go#L199) `validateAgentRunnerBudget`、[agent_runner_launch_reconfirmer.go:35-57](internal/adapters/agent_runner_launch_reconfirmer.go#L35) `ReconfirmLaunchAuthorization` → 用持久预算端口替换裸内存指针，提供 `Reserve(ctx,input)→outcome,ack,error`、只读 `RequireOutstandingReservation(ctx,hash)` 和 `Settle(ctx,input)→outcome,ack,error`；预留在 admission 外持久确认，admission 零写，R 后重新读取当前授权与预算；真实调用/token 上限不能调用前强制即拒绝，CLI credits 不替代产品预算 → `reservation, err := budget.RequireOutstandingReservation(ctx, request.BudgetHash)`
- 读取 [cli.go:428-535](internal/console/cli.go#L428) `executeAgentRunnerRun`、[runtime_agent.go:103-278](internal/console/runtime_agent.go#L103) `ExecuteAgentRunnerRun`、[agent_runner_replay_dispatcher.go:106-168](internal/adapters/agent_runner_replay_dispatcher.go#L106) `DispatchAgentRunner` → CLI 打开宿主固定 scope 的持久预算存储，所有生产端口共享该存储而非各持旧快照；首次打开完成必要持久确认，零值确认无效，R 前验证绑定 reservation 已持久存在，授权回调每次重新加载而非返回装配快照 → `LedgerSnapshot: func() (chain.AuthorizationLedger, error) { return loadAuthorizationLedger(ledgerPath) }`
- 读取 [agent_runner_terminal_publisher.go:227-329](internal/adapters/agent_runner_terminal_publisher.go#L227) `convergeTerminal/settleTerminalIntent` → 结算输入仅取原 intent，核对存储返回确认与条目摘要后才允许发布 C；保留原 `(requestId,C.recordHash)` 派生键及 R/C 证据，unknown 重启后继续占用，完整链在 unproven 下仍只读重放，任何需补写的恢复保持拒绝 → `settled, ack, err := publisher.Budget.Settle(ctx, input); if err != nil { return err }; return ack.RequireEntry(settled.EntryHash)`
- 读取 [agent_runner_postflight.go:695-755](internal/adapters/agent_runner_postflight.go#L695) `WaitAgentRunnerDispatched/agentRunnerSettlementForDeclaredOfflineWorkload`、[agent_runner_terminal_publisher.go:348-401](internal/adapters/agent_runner_terminal_publisher.go#L348) `terminalSettlementIntentFields` → 将 usage 映射器随执行上下文冻结，声明式离线才可补零；真实路径必须验证调用/token 来源及 observed 金额或合法 estimated 依据，费用或用量不完整映射 uncertain、unknown、Facts=nil，不让账本的缺省值冒充观测 → `if !usageVerified { status = "uncertain"; settlement = AgentRunnerTerminalSettlement{Status: tickets.CostSettlementUnknown} }`
- 新增/修改测试：[cost_test.go](internal/tickets/cost_test.go#L1) `TestDurableReserveAcrossProcesses` → 输入同共享 scope、金额上限 100，两个独立进程同时各预留 70；期望只准入一个。
  断言：恰一成功，另一返回 `ErrCostSharedCapExceeded`；重载账本占用为 70，无条目覆盖。
  断言：赢家按原键重试返回原 reservationHash、不追加；路径别名不能取得第二把独立 scope 锁。
- 新增/修改测试：[runtime_agent_e2e_test.go](internal/console/runtime_agent_e2e_test.go#L1) `TestReservationDurabilityGatesDispatch` → 输入预留提交时的 Sync、Close、原子提交、目录屏障失败或不受支持后端；期望没有有效持久确认。
  断言：不发布 R、launch-intent，不 spawn；可见但未确认的预留不删除、不释放。
  断言：重启按原键重新确认后至多保留一条预留，不能靠重新加载内存模型绕过拒绝。
- 新增/修改测试：[agent_runner_terminal_publisher_test.go](internal/adapters/agent_runner_terminal_publisher_test.go#L1) `TestSettlementDurabilityAndUnknownSurviveRestart` → 输入合法 intent，分别注入结算提交失败、提交成功但确认丢失，以及 unknown 计划后重启。
  断言：确认失败时无 C/closure；原键恢复后仅一条结算，同 reservation 异键结算被拒绝。
  断言：unknown 的金额仍为空、占用仍计入共享上限，reservation 不再 outstanding，不能复用或 relaunch。
- 新增/修改测试：[agent_runner_postflight_test.go](internal/adapters/agent_runner_postflight_test.go#L1) `TestPaidUsageCannotBecomeOfflineZero` → 输入付费候选，分别缺金额、缺调用/token 来源；另给合法声明式离线样例。
  断言：付费缺口均为 uncertain、unknown、Facts=nil，不能补零或用预留值充当 observed。
  断言：仅离线映射允许零金额；启动前无法证明调用上限时 admission 拒绝，零 spawn。

## C7
- 读取 [agent_runner_admission.go:42-116](internal/adapters/agent_runner_admission.go#L42) `AdmitAgentRunner`、[runtime_agent.go:103-166](internal/console/runtime_agent.go#L103) `ExecuteAgentRunnerRun` → 新增纯函数 `EvaluateT027BReadiness(goal,input)`，仅消费已核验且绑定候选、平台/P 的证据引用：G0 冻结授权与边界，F 是另获授权的非生产可行性探针，M1/M2/M3 分别验证生命周期、持久提交、串行装配；真实验收要求合法平台前置但不要求 AT-23 已通过，生产发布另要求 AT-23 与独立接受，判定函数不探针、不写记录、不提供绕门参数 → `decision, err := EvaluateT027BReadiness(goal, input); if err != nil { return RunSummary{}, err }; if !decision.Allowed { return RunSummary{}, ErrInvalidConfig }`
- 读取 [ai_probe_copilot.go:62-115](internal/adapters/ai_probe_copilot.go#L62) `CopilotCLIProbe.Probe` → 在 executor 调用前新增内部 ProbePermit 校验，绑定明确授权、profile/候选、平台/P、机器操作、请求上限及停止条件，使 F 可在机制工程之前独立获准执行；不经生产 dispatcher、不伪造 first-dispatch，PONG 只证明 availability，不能冒充会话、工具或完整 enforcement 取证，无合法同边界入口则保持阻塞 → `if err := probe.Permit.Require(ctx, profile, channel, maximumRequests); err != nil { return AIProbeResult{}, err }`
- 读取 [runtime_agent.go:772-786](internal/console/runtime_agent.go#L772) `agentRunnerAssemblyStep`、[agent_runner_enforcement.go:132-222](internal/adapters/agent_runner_enforcement.go#L132) `ValidateAgentRunnerAdmission/nonControlCapabilityFindings` → 将单任务限制收窄为每任务最多一个 Agent 步骤并允许声明 hook，使用 C5 的逐任务上下文驱动同一 Engine；readiness 必须复用原能力校验，保留 sessionResume 等十项要求和当前 resume 路由拒绝，机制替身不得注入生产或抵扣平台发布耐久 → `if err := ValidateAgentRunnerAdmission(capability, enforcement, policy); err != nil { return err }`
- 新增/修改测试：[agent_runner_admission_test.go](internal/adapters/agent_runner_admission_test.go#L1) `TestReadinessDoesNotCreateGateCycle` → 输入 F 无授权、F 有合法授权但生产平台未解锁、真实验收前置齐全但 AT-23 尚未签收三组证据。
  断言：无授权 F 拒绝且 executor 调用为零；有授权 F 的允许不改变生产拒绝或平台等级。
  断言：真实验收不因尚无 AT-23 形成循环；Production 仍因缺 AT-23 拒绝，readiness 本身零读写副作用。
- 新增/修改测试：[agent_runner_enforcement_test.go](internal/adapters/agent_runner_enforcement_test.go#L1) `TestCreateStillRequiresResumeCapability` → 输入新鲜 availability 与有效其余能力，仅 SessionResume 分别为 unknown、unsupported；请求 mode=create。
  断言：两种输入均 admission blocked，R/launch-intent/spawn 数均为零。
  断言：补偿 tool/network、首次不执行 resume 或 F 成功均不能改变该结论。
- 新增/修改测试：[runtime_agent_e2e_test.go](internal/console/runtime_agent_e2e_test.go#L1) `TestThreeTaskMechanismChainStopsAtSecondFailure` → 输入同一 Engine 的三个零模型任务及 hooks，T1 接受、T2 自然失败、T3 待执行。
  断言：仅一份 Engine/事件链，T2 绑定 T1 接受父快照，T3 的物化、预留、派发均为零。
  断言：T2/task 与 chain 为 FAILED；机制通过证据不产生 AT-23 签收或真实生产准入。

## 跨项约束
- 批次 A 的句柄、封存及读取 helper 不改名。
- 核心单写，业务锁不跨外部等待。
- 共享预算锁不嵌套 replay 或 registry 锁。
- 封存→intent→持久结算→C→closure。
- 重放不启动、不重复计费、不替换候选。
- 全部测试为拟新增，未执行测试或 lint。

CONTRACT-BLOCK: Windows 平台发布耐久升级须先修订 §7 并经独立评审，本批不得仅凭机制或探针成功解锁真实生产链。