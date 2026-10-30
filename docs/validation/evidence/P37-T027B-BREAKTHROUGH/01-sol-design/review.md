**T027-B 生产级架构设计书**
【设计决策】主旨：补齐可信执行边界与生产装配，复用既有核心状态机。
1. 耐久按“平台 × 前提 P”三态举证，Windows 默认仍为 `unproven`。
2. 固定官方 CLI 必须在实际强制边界内完成会话与工具迭代，不能移出边界换取成功。
3. dispatcher 独占启动；终局、技术门禁和独立接受分离，未知不得表示为 completed。

**D1 破局路径**
【事实】P35 的判断是“双断点”，见 [P35:9](docs/validation/strategic-review-2026-10-30.md#L9)；P36 明定立项裁定“不产生生产准入凭据”，见 [P36:42](docs/validation/strategic-planning-2026-10-31.md#L42)。

【推断】当前缺口是耐久、受限候选可运行性、真实证据及多任务装配的组合，单补 enforcement 记录不能贯通。

【设计决策】11 月按以下门顺序推进，不以日期替代条件：

| 门 | 准入与独立验收 | 失败判据 |
|---|---|---|
| G0 立项 | 用户裁定 OB-55 归属，冻结候选、Windows/P、资源费用及 AT-23 逐项映射 | 任一未裁，禁止实施 |
| G1 耐久 | 首选验证 C2；覆盖生产引导目录及 R、launch、完整 terminal 载荷；满足既有区分力、机制证据和独评要求 | 任一缺口维持 `unproven`，首次派发零写拒绝 |
| G2 边界 | 同一候选、平台、配置下，其他十项能力 verified，工具/网络分别有 verified 来源；官方会话及工具正例与越界反例同时成立 | 会话不可运行、旁路或证据陈旧即阻塞，不放宽权限 |
| G3 装配 | 零模型运行三任务；覆盖撤销、超预算、重复派发、终局发布中断、缺证、拒评等反例 | 反例仍启动、发布或推进即失败；不能只断言非零退出 |
| G4 真实闭环 | 前三门满足且真实调用另获授权；关闭 IDE，验证三任务成功、T2 失败不启 T3、重启不重发、交接复检、独立接受及源树零写 | 任一必需证据缺失，AT-23 仍未通过 |

现有 C2 证据只覆盖限定虚机/P 与约 1 KiB 载荷，不能直接兑现生产升级，见 [耐久报告:19](docs/validation/t027-b3b-durability.md#L19)。等级账本记录 `proven/disproven/unresolved`；后两者运行行为仍映射 `unproven`。

各门绑定输入摘要、命令、原始证据和裁定；机制、装配、真实候选分栏记账，不合并成“已通过”。ZIP、离线核验和备份恢复排在 G4 之后。

**D2 模块划分**
【事实】核心已有消费侧端口，见 [models.go:185](internal/chain/models.go#L185)；已有 replay-aware 派发与终局发布，见 [dispatcher:107](internal/adapters/agent_runner_replay_dispatcher.go#L107)、[publisher:61](internal/adapters/agent_runner_terminal_publisher.go#L61)。当前装配限定单任务单 Agent 步，见 [装配:113](internal/console/runtime_agent.go#L113)。

【设计决策】组件关系：
```text
CLI 装配[扩展] → chain / StepRouter[既有]
 → admission → dispatcher / replay[既有]
 → Windows 强制边界 launcher[新增] → 固定官方 CLI
 → 监督与厂商事件/用量转换[新增]
 → terminal publisher[既有] → chain terminal 路由[既有]
 → 独立 hooks → postflight → freeze → review → promotion
   [既有端口；补全生产实现和装配]
```
- `chain` 唯一拥有调度、状态和接受决定，不导入具体 adapter。
- `guard` 拥有受管进程与停机机制；Windows 令牌、Job、ACL 及网络实现留在平台层。
- adapter 只转换外部事实，不批准候选；宿主独占 replay、证据、策略及审批写权。
- 每任务从上一已接受快照物化，届时冻结独立 request、授权与 reservation；首验每任务一个 Agent step，不能复用单步装配三次冒充同一任务链。
- SessionBridge 不进入主执行链，不自动转黑箱。新增类型仅为内部 DTO，不新增 wire、状态或接受语义。

**D3 状态机**
【事实】终局路由及六事实 postflight 已冻结，见 [CONTRACTS:228](docs/CONTRACTS.md#L228)。现 admission 仍明确拒绝 resume 连续性，见 [admission:74](internal/adapters/agent_runner_admission.go#L74)。

【设计决策】正常路径为：

`chain CREATED→BASELINED→RUNNING`；task `PENDING→PRECHECK→STEPS_RUNNING`；Agent step `PENDING→RUNNING→TERMINAL_PENDING`。completed 只允许该 step PASSED；剩余声明 hook 通过后，依次 postflight、`REVIEW_PENDING`、freeze、独立 review、completed promotion，task 才 PASSED；全部任务接受才 chain COMPLETED。

不确定落点必须逐类保留：
- R-only、intent 缺 identity、失联、停机未证、收敛失败：保留原记录并阻断；不得删除意图或重新 spawn。
- 超时、日志/usage/manifest/pid 报告缺口：按既有优先级产生 uncertain；可路由时 task FAILED、chain PAUSED，禁止重试。无法合法发布收据时保留悬挂并报告错误，不伪造 C。
- 自然非零退出仍为 failed，不因缺证改成 completed；cancelled 先归档停机证据，停机失败保持暂停。
- operator 请求进入 task/step `WAITING_FOR_OPERATOR` 并暂停；只有绑定有效的结构化响应可继续控制流程。
- postflight 拒绝进入 `REPAIR_PENDING+PAUSED`；未知进入 `FAILED+PAUSED`；损坏或冲突零写拒绝。
- 重启只重放、核验和补齐自有终局事务，不恢复未知模型执行。首验维持核心 attempt 1 边界；新 attempt 与真实 session resume 单列后续实施门。

**D4 关键接口**
【事实】现等待入口使用声明式离线零费用映射，见 [postflight:735](internal/adapters/agent_runner_postflight.go#L735)，不能直接作为真实付费路径。

【设计决策】复用原签名：
```go
AdmitAgentRunner(context.Context, chain.StepRequest) error
DispatchAgentRunner(context.Context, chain.StepRequest) (chain.StepResult, error)
RunPostflight(context.Context, chain.PostflightRequest) (chain.PostflightDecision, error)
```
admission 可重复但只是当时预检；dispatch 重放返回阻断或终局存在，不重复执行；postflight 依既有资格证据收敛，不能用旧通过跨候选复用。终局 publisher、翻译与 `SubmitAgentRunnerTerminal` 保持原接口，同绑定幂等、异终局冲突。

新增内部接口建议：
```go
Prepare(context.Context, BoundarySpec) (PreparedBoundary, error)
WaitDispatched(context.Context, DispatchedRun) (AgentRunnerTerminalOutcome, error)
Decode(RawEvidence) (ObservedRun, error)
```
- `Prepare` 属平台执行层：绑定候选、平台、策略、target 与环境白名单；同绑定幂等，不扩大权限；准备失败不得启动。
- `WaitDispatched` 属 adapter：绑定既有 launch/R、前 manifest 和运行上界；只等待、采集及映射，绝不 spawn。重入只读取已定事实；无连续性证据即阻断。
- `Decode` 属固定版本转换器：同字节确定性输出真实 session、事件和用量；未知版本、撕裂或缺字段返回缺口，不补零、不合成成功。
- 所有失败保留底层分类；清理失败不得覆盖首次失败原因。

**D5 资源约束**
【事实】安全权威要求网络默认拒绝、例外明确 host/port/purpose，见 [SECURITY:49](docs/SECURITY.md#L49)。B2 明记旧候选容器会话失败及全回环豁免风险，见 [B2:156](docs/validation/B2-EXTERNAL-ENFORCEMENT.md#L156)。

【设计决策】
- 候选拟以官方 CLI `1.0.85` 为冻结验证起点，不表示 compatible；G0 补齐版本、平台包、可执行文件及配置摘要。禁 `latest`、自动更新、模型/provider/通道 fallback；十二项能力重新取证。
- 唯一 watchdog 管运行截止时间，取任务与授权墙钟上界较小值；grace 显式冻结，停机核验使用独立、有界且不受调用方取消打断的 context；保留既有取消优先落定规则。
- CPU、内存、子进程数、日志和磁盘上限纳入批准配置；无法强制则拒绝相应任务。输出在写入阶段限流，不能只在退出后截断。
- 凭据仅 SecretStore 临时解析、最小环境注入并清零；不复制整个用户目录或认证缓存。机器变更须授权 apply/revert，运行时不自动提权。
- 网络仅开放获批出口；全回环豁免不等于“只能访问代理”，监听变化或代理审计失败必须阻断。
- 探针与每次执行分别计量；调用前预留，发布 R 后重读授权/预算，运行中记录调用、token、墙钟及计费来源。结算按 observed/estimated/unknown 分列，unknown 保留占用，不默认为零。准则模型闭集及额度原样执行。

**D6 风险未决**
【事实】当前装配授权快照返回装配时账本，见 [装配:209](internal/console/runtime_agent.go#L209)；`Stop` 仅返回摘要、`Reconcile` 直接返回 nil，见 [装配:530](internal/console/runtime_agent.go#L530)。这些不能作为生产停机或恢复证明。

【设计决策】
- **死锁**：run 单写所有权常驻；预算/store 锁只覆盖短事务，固定锁序，禁止持锁等待模型、管道或回调重入。
- **孤儿进程**：受限边界与 Job 就绪后才恢复执行；spawn 后任一身份发布失败均由所有者有界回收。无空树证明不得接受或复用执行窗口。
- **TOCTOU**：保护目录与候选安装区用权限隔离，校验真实路径、文件同一性及打开句柄；每次发布复检。用户态检查残余窗口不得宣称消除。
- **竞态**：重确认必须读取持久账本当前值；共享预算不能只用 run 锁。启动槽 mtime 不是存活证明；超时/取消只能落一个终局。
- **证据伪造**：宿主采集并保护日志、用量和报告；pid 自报仅作辅助，不能替代 OS 进程树证明。空报告必须是真正零字节，缺失不能补空。
- **接受短路**：替换当前固定 approve 装配，接入独立 reviewer；补齐授权范围、秘密、文件类型/大小及可观测副作用检查，冻结后变化使旧批准失效。
- **恢复歧义**：逐个崩溃点验证 intent→结算→C→closure；只补自有事务，不收养孤儿、不退款未知费用、不重复启动。

**仍需用户裁定**
- OB-55 是否归入独立 T027-B 实施片，并允许提前落地有实际调用者的平台边界接口。
- 固定候选及目标 Windows/P；是否批准重新取得耐久与同边界真实探针证据。
- AT-23 三单元映射、签收责任人，以及首验不覆盖 session resume 的边界。
- 资源、费用、人工介入上限及失败后的停止/延期条件。

本轮已读完全部必读输入；未写文件、未运行测试或探针、未调用子代理。
