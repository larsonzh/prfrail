# ⑦ 独立终审 r1（安全评审本体）· 逐字原文（P79）

> 归档说明（§13.1.13）：本文件为 ⑦ 独立终审层 r1（首轮盲态）最终消息的逐字转录，主控未按要点重建、未改写措辞、未重组结构。
> 派单：`agentName = independent-reviewer`，`model = GPT-6.2 Sol (copilot)`（深度档，Max）；轮次 = r1（终审，额度 1 + 3 中的「1 终审」）。
> 盲态：首轮输入 = 起跑消息原文 + ADR-047 原文 + 三面定义 + 安全基线与实现文件 + 方案书 §2/§6；未附 ⑤ / ⑥ 清单，未附主控工作笔记。

---

定点核对完成。三面都有能定位到源码和契约的发现；结论会区分已接线的 BYOK 路径和仍被上游门禁挡住的装配路径，并如实注明没有运行测试。

## 面 1 威胁模型
### Section A 结论
FACE 1: FINDINGS
### Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| Medium | [internal/gates/runner.go:208](internal/gates/runner.go#L208) | [docs/SECURITY.md:23](docs/SECURITY.md#L23)：「限制必须由 runner/OS 能力落实；不能仅靠提示词。」 | 将产物读取绑定到授权工作区目录句柄，在打开时拒绝越界；无法安全打开即拒绝收集，不入证据库。 |

### Section C 逐条证据与推理
读取 [docs/SECURITY.md:11](docs/SECURITY.md#L11)、[docs/SECURITY_EN.md:11](docs/SECURITY_EN.md#L11) 和 [docs/TEST_STRATEGY.md:29](docs/TEST_STRATEGY.md#L29)：双语基线一致区分了目录分离与进程隔离、内容完整性与身份认证。STRIDE 已列出路径逃逸，缺口在于 SEC-02 承诺的打开时检查未覆盖所有读取入口。读取 [internal/gates/runner.go:194](internal/gates/runner.go#L194)：产物收集按路径展开，只对末端做 `Lstat`，再按路径重新打开。它既未约束中间目录解析，也未把检查与读取绑定到同一对象。现有 `GuardExecutor` 缺必需能力而拒绝执行，这限制了当前可达性，不补足该读取边界。
### Section D 可证伪性审计
- 删除 **BoundDir** 的中间路径限制，`TestBoundDirMiddleSymlinkRefusedWithClosedErrno` 会失去预期拒绝；它不经过 Gate 产物收集。删除证据对象摘要校验，`TestVerifyManifestObjectsFailsClosed` 的 `tampered` 子项会红；删除事件归属检查，`TestVerifyManifestStateEventBinding` 的跨 run 子项会红。
- 删除 C6 的魔数检查或未提交槽判定，`TestC6SlotValidationOrderBadMagicIsCorruption` 或 `TestC6SlotValidationOrderUncommittedPayloadIsAbsent` 会红。基线已明示装配记录丢失后的漂移检测降级，不把它记为完整性已证明。
- UNCOVERED：未见 Gate 产物读取跨目录边界时同时断言拒绝、证据库无越界对象、任务不推进。现有 AT 编号不能代替这些断言。
### Section E 我无法验证的部分 + 什么会推翻我的判断
未验证原生平台上的完整 STRIDE 运行证据。若能指出所有 `Runner` 产物读取调用都由不可绕开的根目录绑定保护，并有针对这条读取路径的反例，才会推翻本项。

## 面 2 硬安全门禁
### Section A 结论
FACE 2: FINDINGS
### Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| High | [internal/console/runtime_agent.go:1245](internal/console/runtime_agent.go#L1245) | [docs/CONTRACTS.md:212](docs/CONTRACTS.md#L212)：「生产 admission 是纯 preflight 与不可变绑定：不得消费 request replay 身份、不得发布 R/C 或写 replay store、不得启动进程；首次 dispatch 资格与 replay 身份消费由 dispatch 路径在发布 R 之后依据 replay store 判定；admission 的结论是评估时刻的 preflight verdict，不是锁，可能在 dispatch 前失效，dispatch 在发布 R 之后、允许 launch 之前必须重新确认授权未被撤销且预算 reservation 仍 outstanding，任一不满足即 fail-closed，不得把先前的 admission 通过当作 launch 许可。」 | 在 launch 前向可信账本读取当前授权与预留；读取失败、撤销或已结算均保留 R，拒绝启动。 |
| High | [internal/chain/review.go:116](internal/chain/review.go#L116) | [docs/CONTRACTS.md:147](docs/CONTRACTS.md#L147)：「评审者与候选变更产生者是否同一身份、`evidenceRootHash`/`policyHash` 引用是否真实存在及 waiver 到期后不得复用，均由独立 checker 判定。」 | 发布前从执行者不可写的可信来源验证策略批准与候选证据根；没有有效绑定则拒绝，或停在待人工批准。 |
| Medium | [internal/console/runtime_agent.go:731](internal/console/runtime_agent.go#L731)、[internal/console/runtime_supervised.go:1642](internal/console/runtime_supervised.go#L1642) | [docs/CONTRACTS.md:150](docs/CONTRACTS.md#L150)：「`promotion-receipt.schema.json` 将父/候选 snapshot、evidence root、review receipt、停机与租约证据绑定为单个发布事务。」 | 发布前核对当前工作区的排他租约或等价写入隔离证明，将其真实记录绑定到本次 attempt；缺失或不一致则拒绝。 |
| Medium | [internal/gates/runner.go:73](internal/gates/runner.go#L73) | [docs/CONTRACTS.md:391](docs/CONTRACTS.md#L391)：「只读/本地可丢弃/外部写分离；S1 外部写拒绝，未知副作用暂停对账；诊断不变更锁/日志」 | 执行前从冻结 hook 策略导出有效副作用分类；分类缺失、矛盾或不允许时不能进入 `Execute`。 |

### Section C 逐条证据与推理
- **撤销独立性。** 读取 [internal/console/cli.go:562](internal/console/cli.go#L562) 和 [internal/adapters/agent_runner_replay_dispatcher.go:150](internal/adapters/agent_runner_replay_dispatcher.go#L150)：dispatcher 在正确的位置调用复核，但生产回调返回 CLI 启动时读取的授权副本。复核可发现随时间过期，不能发现初读后另行持久化的撤销。预算也只接到已装载对象；未证明可观察跨进程结算。
- **策略与批准。** 读取 [internal/console/runtime_agent.go:678](internal/console/runtime_agent.go#L678)、[internal/console/runtime_supervised.go:1598](internal/console/runtime_supervised.go#L1598) 和 [internal/chain/engine.go:254](internal/chain/engine.go#L254)：两条装配路径在 postflight 合格后本地计算策略摘要并返回 `approve`。同名拒绝是真守卫，但摘要合形不证明策略事先获独立批准；postflight 与批准不能互作原件。
- **租约与副作用。** 读取 [internal/guard/lease.go:104](internal/guard/lease.go#L104) 和 [internal/console/runtime_supervised.go:1491](internal/console/runtime_supervised.go#L1491)：租约库缺独立 verifier 会拒绝接管，但发布端的 `LeaseEvidence` 只由 run ID 派生，不证明当前独占。归还 hook 接线未填 `Effect`，而 runner 仅在有分类记录时调用准入。现版 CLI 另有 hook 定义和 runner 能力拒绝，不能把这个缺口说成现版已经执行了外部写入。
### Section D 可证伪性审计
- 删除能力检查，`TestPrepareRejectsUnavailablePolicy`、`TestRunnerDoesNotExecuteWhenPolicyIsUnavailable` 会红；删除启用后的沙箱拒绝，`TestAgentRunnerRunSandboxRefusesBeforePermit`、`TestSupervisedRunSandboxRefusesBeforePermit` 会红。契约明订的默认 `off` 不算绕过启用后的检查。
- 删除同路径独占或损坏 marker 拒绝，`TestProbePermitCrossProcessSingleClaim`、`TestProbePermitDamagedMarkerStaysSpent` 分别会红。删除就绪级别阻断，`TestReadinessProbeConclusionCannotAdvanceAcceptance`、`TestReadinessProductionRequiresAT23` 会红；这些是级别和格式断言，不是证据签发者认证。
- 删除同名自批、过期 waiver、postflight 事实绑定或接管 verifier 检查，分别由 `TestRejectSelfReview`、`TestRejectExpiredOrWrongCandidateWaiver`、`TestPostflightGateRefusesAnUnprovenPass`、`TestLeaseTakeoverFailsClosedWithoutVerifier` 反证。
- UNCOVERED：生产接线缺少这四项反例：初读后持久化撤销且 **R 保留、零启动**；策略未获独立批准且 **零发布**；租约证明缺失且 **无完成回执**；副作用分类缺失且 **零执行**。账本读取失败测试和显式外部写分类测试覆盖的是相邻分支。
- UNCOVERED：探针测试证明同一路径只认领一次，不证明独立签发或不同路径不能代表同一授权；`--probe-permit` 的调用方信任边界不能用 marker 的存在代证。
### Section E 我无法验证的部分 + 什么会推翻我的判断
未验证独立控制面是否在这两条装配路径外实施了不可绕过的实时账本、策略批准、排他写入或 hook 分类。能反驳上述判断的是实际接线与对应的零启动、零发布或零执行断言。未交付的路径仍应保持阻断。

## 面 3 秘密与凭据
### Section A 结论
FACE 3: FINDINGS
### Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| High | [internal/console/ai.go:94](internal/console/ai.go#L94) | [docs/SECURITY.md:60](docs/SECURITY.md#L60)：「网络默认拒绝，例外授权须列 host/port/purpose、提供方数据策略及期限。」 | 解析凭据前，由可信本机控制路径校验已批准的 profile、凭据引用、接收端与子进程身份绑定；不符即不读库、不启动。 |
| Medium | [internal/guard/process_windows.go:457](internal/guard/process_windows.go#L457) | [docs/CONTRACTS.md:216](docs/CONTRACTS.md#L216)：「adapter 只接收解析后的短暂字节副本，使用后清零；系统错误只返回稳定分类，不得拼接 credential blob。」 | 对 Windows 启动所用的原生环境块及编码中间缓冲保留所有权；成功或失败均及时清零，并核对凭据读取与终端输入的临时副本。 |
| Medium | [internal/adapters/ai_probe_copilot_executor.go:53](internal/adapters/ai_probe_copilot_executor.go#L53) | [docs/CONTRACTS.md:214](docs/CONTRACTS.md#L214)：「profile 只保存 provider 类型、base URL、模型、传输选项及不透明 `secretRef`；API key/token 不得进入配置、摘要、日志、错误或证据，运行时只能由系统 SecretStore 解析并临时注入子进程环境。」 | 只保留经过校验的结构化事实，不对未脱敏流生成持久证据；关闭未证明脱敏的子进程落盘日志，无法保证则拒绝 BYOK 探针。 |

### Section C 逐条证据与推理
- 读取 [internal/adapters/ai_provider.go:115](internal/adapters/ai_provider.go#L115)、[internal/adapters/ai_secret_windows.go:139](internal/adapters/ai_secret_windows.go#L139) 和 [internal/console/ai.go:66](internal/console/ai.go#L66)：profile 验证只检查引用形状和 HTTPS 语法，Windows 后端按当前用户的目标名解析。`ai check` 的配置选择和可执行文件路径没有独立核对凭据与接收方的批准绑定。单请求授权、HTTPS 和配置摘要均不能证明接收方有权持有凭据。
- 读取 [internal/adapters/ai_preflight.go:57](internal/adapters/ai_preflight.go#L57)、[internal/adapters/ai_probe_copilot.go:80](internal/adapters/ai_probe_copilot.go#L80) 和 [internal/guard/process_windows.go:221](internal/guard/process_windows.go#L221)：SecretStore 返回的 Go 切片、探针环境副本及 executor 的字节缓冲有清零；真正传给 Windows 的 UTF-16 环境块及转换副本没有。`CredFree` 调用只能据实记为释放，不能据此认定擦除。
- 读取 [internal/adapters/ai_probe_copilot_executor.go:92](internal/adapters/ai_probe_copilot_executor.go#L92) 和 [internal/adapters/ai_secret_unsupported.go:22](internal/adapters/ai_secret_unsupported.go#L22)：ProofRail 未把解析值放进探针 argv，非 Windows SecretStore 拒绝解析，结构化原因也不转录原始错误。但子进程可写临时日志，原始 stdout/stderr/usage 未经秘密检查即入摘要，目录清理不能证明从未落盘。未观察到实际泄漏。
### Section D 可证伪性审计
- 去掉 SecretStore 返回副本或探针调用环境副本的清零，`TestAIAvailabilityCheckUsesOnlyRequestedChannelAndClearsSecret`、`TestCopilotCLIProbeBuildsIsolatedDeepSeekSmoke` 分别会红。去掉缺秘密时的拒绝，`TestAIAvailabilityCheckBlocksMissingSecretWithoutProbe` 会红。去掉 explain 遮蔽，`TestAIConfigRoundTripAndExplainRedactsSecretReference` 会红。
- UNCOVERED：未见生产入口对未批准的凭据/接收方绑定断言零解析、零启动；也未见覆盖原生环境块的成功及失败清零、真实平台后端拒绝后零启动，以及合成秘密在子进程输出、日志及最终证据中的端到端不泄漏检查。
### Section E 我无法验证的部分 + 什么会推翻我的判断
未验证 Windows 原生调用和子进程自身的日志行为。可推翻相应发现的是生产入口已有不可绕过的批准绑定、原生缓冲的完整清零路径，或能证明日志与摘要只消费脱敏事实的端到端反例。

## 合成
FINAL REVIEW: FINDINGS

读取 [docs/ADR_REGISTER.md:58](docs/ADR_REGISTER.md#L58) 与 [docs/PRODUCT_REQUIREMENTS.md:109](docs/PRODUCT_REQUIREMENTS.md#L109)：三面各有不符合判据之处，不能登记为安全门通过，也不能以此替代 S1 Exit 评审。以上是只读判断，未执行命令、未变更文件。方案书定位回显夹带了许可节之外的片段，未采信；未读取 ⑤/⑥ 清单或 P79 报告结论。
