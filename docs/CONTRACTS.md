# ProofRail 架构契约规范

[English](CONTRACTS_EN.md)

日期：2026-09-07；S1 架构契约基线。本文是 wire、Schema、状态转换、事件、receipt 与跨记录不变式的规范性权威；[项目建议书](RFC-proofrail-unattended-ai-engineering-product.md) 保留来源、范围和设计理由。本文合并 Schema、快照/证据、hook、票据/修复、adapter 五个协议包，减少重复阅读。

## 1. 契约成熟度与兼容

下表语义由本文确立，并保留项目建议书来源追踪；“待冻结”是生产实现阻断项。T002 已建立首批结构 Schema，
但本文不是完整 JSON Schema，也不能代替独立 validator。新增 wire 字段、枚举、canonical 算法和错误码
必须先经 [ADR](ADR_REGISTER.md) 更新本文，再建立 Schema 与正反黄金样例，最后写 Go 代码。

| 契约 | 已确定的语义 | S0 待冻结产物 |
|---|---|---|
| 配置/任务/步骤 | Draft 2020-12；严格 chain/task/四类 step；workspace 拓扑；harness/toolchain 注册表；确定性所有权/合并 | 独立正反 fixtures |
| change-set | 五类严格操作、顺序内存预验证、first-fail-stop、marker/精确断言、原始字节摘要与整组事务 | 跨记录 checker 与正反 fixtures |
| snapshot/evidence | SHA-256 内容寻址、父快照/候选/证据绑定、不可变 | 路径 canonical、hash 输入编码、manifest/receipt schema |
| ticket/repair | 稳定指纹、append-only ledger、三段预算；严格 Prepare→Inspect→Validate→Promote 阶段链 | 跨记录 checker 与正反 fixtures |
| adapter/context/agent-runner | 严格 request/claim/result 信封、recordHash、幂等键、generation fencing、dispatch/takeover receipt、最小上下文；CLI Agent 进程、会话、能力与结果映射 | 现有 adapter 独立信任链 fixtures；AgentRunner Schema/正反样例与真实 CLI 能力探针待冻结 |
| product/lifecycle | 计划预览、导出、授权/撤销、副作用/恢复、成本账本、发行/备份/退役记录 | 独立正反 fixtures 与 S1 运行验收 |

各持久对象独立 `schemaVersion`，不与 CLI 版本或 SessionBridge schemaVersion 混为一谈；首版为
`1.0.0`，未知主版本拒绝，未知 minor/patch 默认只读拒绝。仓库、Schema、运行时与 wire JSON 统一为
UTF-8 无 BOM + LF；canonical 字节还不得含格式空白或尾随换行。仅明确验证 BOM 输入兼容性的 fixture
可带 BOM，且不能复用为业务对象。

## 2. Schema 与静态检查责任

| 对象 | 最小语义 | 非法/动态约束 |
|---|---|---|
| chain definition | `chain.id/profile`、至少一个有序 task、最小 workspace、可选 documentation | task/step ID 跨数组唯一、引用及依赖无环由 checker 判定 |
| run manifest | chain 输入/有效摘要、逐叶值来源、policy/父快照及版本化运行绑定 | checker 验证 pointer/来源闭包、实际加载对象、能力证据和 run 唯一 manifest |
| task | 稳定 id、非空 steps、review、documentationPolicy | id 重复、skip-on-fail 未声明 failureIndependent 拒绝 |
| step | id、kind=code/build/verify/noop；code execution 默认 autonomous | 空 steps、未知 kind 拒绝；noop 必须 reason 且不能启动 hook/agent |
| hook | 独立严格注册表；process/container 判别联合、逐项参数、失败/超时/资源/网络/产物策略 | ID/引用、工具摘要与能力由 checker/preflight 判定；执行结果进入 receipt |
| harness/toolchain | harness 组合语言、参数、hook 和产物；toolchain 声明来源、平台、probe 与证据策略 | checker/preflight 验证唯一引用、版本、来源、摘要、授权与实际可用性 |
| target | 独立严格注册表；稳定 id、class、access、唯一路径 glob；generated 单向来源 | 越界、重叠写、循环生成、未知引用由 checker 拒绝 |
| component/language scope | component root；scope 的 language/harness/toolchain/target ID；依赖；step 取 scope 交集 | ID 唯一、引用存在、无环、harness 可用；同 root 允许不同唯一 target 所有者 |
| state event | append-only envelope；chain/task/step 实体；严格前后状态；序号、前序 hash、证据与原因 | 跨事件序号/hash/实体状态连续性、主体授权和证据存在性由 checker 判定 |
| error/error-set | 37 个封闭 code 及 retry/action 组合；非空错误索引与 primaryErrorId | checker 验证引用、唯一性及按时间/category/code/errorId 的主错误排序 |
| verification report | 16 类固定顺序的跨记录/信任链检查；passed/failed/incomplete | checker 复读对象并验证摘要、适用检查集、顺序与总结果；报告不改变状态 |
| product records | preview/export/authorization/effect/cost/lifecycle 严格判别记录 | checker 验证授权时点、引用闭包、算术、平台事实及禁止的 S1 external-write |
| documentation | required/if-affected/optional/forbidden；impactRules | 命中规则缺有效 doc diff、生成物陈旧拒绝；空白/mtime 不算更新 |
| handoff | execution、`handoffPolicy`（allowedTargets 引用 read-write target、inputPolicy、handoffTimeoutMs、returnActions、hooksAfterReturn） | 单写租约、停机、新鲜 manifest；归还后复检不能由 Schema 单独判定 |
| review/waiver | approve/reject/waive；主体、原因、策略依据、有效期、绑定候选 | 同主体自批、过期、候选变化、范围不符均阻断 |

三层检查不可合并：JSON Schema 判结构；语义 checker 判引用、所有权、循环、静态策略；运行期 gate 判文件 diff、进程、能力、秘密、预算与评审。非法黄金样例标注拒绝层，不能声称 JSON Schema 会检查实际文件或进程。

`schemas/` 已创建 chain、hook、target、workspace、state-event、error、adapter-envelope、adapter-receipt、
snapshot-manifest、evidence-manifest、review-receipt、promotion-receipt、handoff-receipt、hook-result、run-manifest、signature-receipt、error-set、ticket-ledger、repair-transaction、harness、toolchain、change-set、verification-report、plan-preview、export-record、authorization-record、effect-record、cost-ledger、lifecycle-record 二十九份结构 Schema；`testdata/contracts/valid/` 与
`testdata/contracts/invalid/` 已由 T003 建立并通过 `tools/contracts/` 独立 validator。每例携带 fixture ID、契约版本、
expected accept/reject、拒绝层和原因。必须含三任务链、四 kind、C+Go、C+JavaScript+Python、人工交接、
代码文档协同；每一条件至少一正一反。文档中的路径/片段是设计输入，不是已验证可运行样例。

### 2.1 Chain 状态、事件与恢复（T009 权威）

本节是 T009 对 chain/task/step 状态、顺序调度、事件投影和恢复行为的规范性来源。项目建议书
§10.3、§16.9.8 保留为设计来源；两者与本节冲突时阻断实现并修订文档，不由实现者自行选择。

Schema 只接受以下单步转换，未列出的转换、自转换和跨 attempt 原地重跑全部拒绝：

- chain：`NONE→CREATED`；`CREATED→BASELINED|FAILED|CANCELLED`；
	`BASELINED→RUNNING|FAILED|CANCELLED`；`RUNNING→PAUSED|COMPLETED|FAILED|CANCELLED`；
	`PAUSED→RUNNING|FAILED|CANCELLED`。
- task：`NONE→PENDING`；`PENDING→PRECHECK|FAILED|CANCELLED`；
	`PRECHECK→STEPS_RUNNING|FAILED|CANCELLED`；
	`STEPS_RUNNING→WAITING_FOR_OPERATOR|REVIEW_PENDING|FAILED|CANCELLED`；
	`WAITING_FOR_OPERATOR→STEPS_RUNNING|FAILED|CANCELLED`；
	`REVIEW_PENDING→PASSED|REPAIR_PENDING|FAILED|CANCELLED`；`FAILED→REPAIR_PENDING`；
	`REPAIR_PENDING→STEPS_RUNNING|FAILED|CANCELLED`。
- step：`NONE→PENDING`；`PENDING→RUNNING|NOOP_RECORDED|FAILED|CANCELLED`；
	`RUNNING→WAITING_FOR_OPERATOR|PASSED|FAILED|CANCELLED`；
	`WAITING_FOR_OPERATOR→RUNNING|FAILED|CANCELLED`。

`NONE` 只用于创建事件，不是投影可停留状态。`COMPLETED`、`PASSED`、`NOOP_RECORDED` 和
`CANCELLED` 是终态。只有 task 的 `FAILED→REPAIR_PENDING` 可离开失败态；step 修复必须增加 attempt。

- 每次状态转换先持久化 append-only state event，再更新可丢弃、可重建的 projection；事件按 run 使用
	从 1 开始的连续 sequence 和 previousEventHash 串联，并绑定 entity、前后状态、actor、证据与 reason。
- 创建 run 的首事件只能是 chain `NONE→CREATED`。baseline-0 成功且其证据已持久化后才能追加
	`CREATED→BASELINED`；不能从工作目录名称、文件存在或进程退出码推断状态。
- task 按定义顺序调度。首个 task 从 baseline-0 物化；后续 task 只能从前一 task 经有效 review 和 completed
	promotion receipt 产生的 accepted snapshot 物化，candidate、失败或不确定快照不得成为父快照。
- 每个 task 内 step 按 sequence 执行；`noop` 只允许 `PENDING→NOOP_RECORDED` 且必须记录原因，不启动
	agent、hook 或子进程。T009 使用 fake agent/runner，但其事件与状态规则不得形成替代协议。
- `pause` 在当前原子边界生效，停止启动新 step；`resume` 仅从表中允许的暂停态继续。`cancel` 必须先停止
	受管进程树并归档证据，再写 `CANCELLED`，且同一 run 不可恢复。
- 恢复先证明旧写者失效并取得有效 fencing token，再按 journal 处理未完成写入、重放完整事件链并重建
	projection。断尾、序号/hash/状态不连续、未知 before/after 或外部副作用不确定时保持 `PAUSED` 或
	`REPAIR_PENDING`；禁止猜测成功、重复投递或发布候选。
- 心跳、进度、诊断、review、handoff、takeover 和 promotion receipt 不是状态转换。Schema 判单条形状与
	状态对；checker 判事件链连续性、当前投影、唯一 genesis、actor 权限、证据存在和调度前置条件。

## 3. 配置解析与冻结

`proofrail.toml` 独占 chain/task/step、文档、策略和注册表引用；`workspace.toml` 独占位置、拓扑、平台及
adapter/harness/toolchain 选择，重复/未知键拒绝。优先级固定为 builtin-default < profile < chain < task <
step：缺失继承，scalar 替换，已知 object 递归合并，array 整体替换且绝不拼接；空数组仅在 Schema 允许时
表示明确无项，空字符串不等于 unset，不支持 null 删除。每个最终叶在 run manifest 中必须有唯一 pointer
来源；`config explain` 只读展示值/脱敏引用及覆盖链，不执行或回写。

解析完成再按 RFC 8785 JCS 计算 canonical run manifest 摘要，以 UTF-8 无 BOM 字节输入 SHA-256，
文本表示为 `sha256:` 加 64 位小写十六进制。运行期心跳、credentials、机器探测不得回写用户配置。
父 snapshot、hook 摘要、adapter/harness 版本和已授权策略绑定 run。变更配置只能新建 run，或走本文
明确的暂停、审批和新 manifest 流程；不能覆盖旧事实。

`run-manifest.schema.json` 以 `chainDefinitionHash`/`effectiveChainHash` 区分通过检查的输入定义与完全
展开默认值后的 chain；`resolution` 用 RFC 6901 pointer 记录每个最终叶值来自 builtin/profile/chain/
task/step，`bindings` 固定 adapter/harness/toolchain/hook/policy 的版本、对象摘要和能力证据。绝对路径、
凭据、环境变量值及可变机器探测不进入 manifest；pointer 存在性、来源覆盖完整性、实际加载对象和
runId 唯一 manifest 由 checker 判定。

`harness.schema.json` 只组合语言、toolchain、声明式参数、已注册 hook 与 artifact glob，不复制 runner 或
命令；`toolchain.schema.json` 声明来源、平台、分离的 executable/args probe 和证据策略。注册表存在不证明
工具已安装；引用、版本、来源、摘要及实际能力由 checker/preflight 验证，探测仍需授权。S1 随附
generic/C/Go 不形成核心语言枚举，环境值、凭据和绝对主机路径不得进入注册表或 run manifest。

RFC §16.9.3.1 已冻结 `jcs-001` 与 `state-event-001` 两条字节级向量，分别覆盖 JCS 排序/转义/Unicode
及 state-event 域分隔摘要。T003 已将其固化为独立 fixture，并由非核心 RFC 8785 实现复算 canonical 字节与摘要。

## 4. 快照、证据与接受事务

`snapshot-manifest.schema.json` 冻结 baseline/candidate 判别、父快照与 task/attempt 绑定、三类路径条目、
排除规则及无秘密环境事实。普通文件摘要直接覆盖原始字节，不转换 BOM/行尾/编码，并显式标记包清单、
锁文件、生成物和 hardlink 组；目录保留空目录，symlink 摘要覆盖原始相对 target 文本。manifest 本身不含
accepted 状态，只有有效 review + promotion receipt 才能授权下游消费。路径排序/唯一、Unicode/大小写
碰撞、父目录与对象闭包、排除命中及可恢复性由独立语义 checker 和跨平台 goldens 验证；Git 元数据只作
显式证据，不能成为恢复来源。

每次状态事件含 run/task、前后状态、单调序号、时间、主体、输入证据和原因；事件哈希链用于发现缺失/重排。对象先临时写、校验、原子落位；事件持久后才更新可重建投影。断尾日志处理不得丢弃已确认事件而猜测成功。

接受事务的 reader contract：只有绑定父快照、task/attempt、候选及证据根的有效 review 和完成发布的 receipt 同时存在，才允许作为下一任务父快照。崩溃在候选、review、receipt、事件或投影之间时均重放 journal 判定；不得靠目录名或单个 PASSED 字段认定接受。哈希证明完整性，不证明发布者身份。

`evidence-manifest.schema.json` 将 evidence root 冻结为单个 task attempt 的评审前不可变索引：绑定任务定义、
解析策略、父/候选 snapshot，并引用状态事件、错误、adapter、hook、artifact、change-set、diff、ticket、
repair、authorization、effect、cost 与 handoff 对象摘要。评审后生成的 verification/export/lifecycle 记录
不得反向进入该 root。它不内嵌日志，不从对象存在推断成功，也不包含 review/waiver/promotion receipt；后者
单向引用 evidence root，避免摘要环。条目排序/唯一、同 attempt 归属、对象存在、必需证据覆盖及脱敏
真实性由独立 checker 判定。

`review-receipt.schema.json` 是唯一产生评审结论的对象，直接引用 `evidenceRootHash` 并输出 `approve/reject/waive` 之一。
`recordedBy.type` 限定 `operator|policy`，禁止 `agent|system` 自批；`reviewMode=policy` 必须绑定已批准的
`policyHash`。`waive` 必须同时给出人工 `authorizedBy`、`policyBasisHash`、豁免 `scope` 和非空 `expiresAt`，
不得等价于匿名自动批准。评审者与候选变更产生者是否同一身份、`evidenceRootHash`/`policyHash` 引用
是否真实存在及 waiver 到期后不得复用，均由独立 checker 判定。

`promotion-receipt.schema.json` 将父/候选 snapshot、evidence root、review receipt、停机与租约证据绑定为
单个发布事务。只有 outcome=completed、acceptedSnapshotHash 等于候选摘要，且 review 为 approve 或仍有效
waive 时，候选才可供下游消费；failed/uncertain 均保持 acceptedSnapshotHash=null 并阻断推进。发布顺序、
同 attempt 唯一 completed receipt、崩溃 reconcile 和 accepted 引用闭包由 journal/checker 判定。

`handoff-receipt.schema.json` 冻结一次 `manual-handoff` 归还的完成事实：绑定 `handoffPolicyHash`、操作员
身份、租约证据与归还前后 manifest/diff；`complete` 要求非空 `hookResultEvidence`，`abort/request-agent`
则必须为空。越界、秘密或未知进程扫描失败不产生该 receipt，只产生 `error` 记录并保持暂停。

managed-change-set：读全部目标→按全局顺序在内存模拟→检查 marker/前置哈希/断言/边界→持久 journal→逐文件原子替换→写后全量核验→记录 receipt；失败整组回滚，回滚本身失败则隔离并暂停。禁止覆盖外部修改，禁止运行进程存活时恢复。

`change-set.schema.json` 封闭 create/delete/replace-exact/insert-before/insert-after 五类操作，并绑定 task
attempt、父 snapshot、前后 manifest、逐操作前后字节摘要、精确 assertion 和行尾策略。insert marker 必须
从 0 次变为恰好 1 次；replace 可选 marker 也遵守该规则。现有文件保留单一既有行尾，创建文件显式选择
LF/CRLF；混合行尾、模糊/regex 匹配、序号断裂、重复 marker 或任一摘要不符均在写入前拒绝。

isolated-workspace：直接工具写入仅限可丢弃运行目录；记录完整前后 manifest/diff/进程/gates；接受边界仍为整个 task。不能把部分语言、部分文件或“代码成功但文档失败”的子集发布。

## 5. Gate Hook 协议

kind 使用 precheck/build/test/verify/review/cleanup；按序执行。`executable + args[]` 不经过隐式 shell，cwd 必须落在授权运行工作区；环境继承白名单，网络默认拒绝，资源/输出/时长必须可执行限制。需要 shell/container/remote/PTY 是独立能力声明，不支持则拒绝而非弱化策略。

`hook-result.schema.json` 分离 `executionOutcome`（runner 事实）、`assessment`（门禁判断）与
`policyDisposition`（后续控制流），并绑定 hook 定义、执行序号、起止、退出码、stdout/stderr、产物、runner/
停机及错误证据。启动失败、超时、资源终止、进程无法停止、产物缺失和扫描失败均有独立 failureKind；
`warn` 仍为 failed 且不得冒充 PASS，retry 消耗统一预算，cleanup 失败不能抹去原失败。跨记录引用、实际
产物、扫描/停机充分性与重试预算由 checker/runtime 判定。

生成 hook A：只从锁定模板生成参数化命令→语法/模板/依赖/危险能力校验→独立批准→哈希绑定挂载。generated 默认禁用；生成器不能修改自己的验证规则或批准自己。任意生成脚本 B 不进入 S1。

## 6. 票据、修复与交接

ticket 分类沿用 task-static/code-fix/noncode；上下文包含所属 run/task/attempt、失败证据、允许动作、预算和租约。相同指纹预算按 pending_review→override_window→hard_block 限制，具体阈值/有效修复证明必须先冻结，禁止模型自行重置预算。

Prepare 在隔离区建立候选与父摘要；Inspect 比较范围/diff/所有权；Validate 运行冻结门禁；Promote 再校验停机、租约、父/候选/验证摘要与授权，原子提升并写 promotion receipt。任何步骤变化都令后续陈旧验证失效。评审拒绝后修复必须新 attempt，不能改写旧 receipt。

handoff：停止受管写者→flush journal→捕获 manifest→WAITING_FOR_OPERATOR→人工独占写租约。complete 收回租约并检查越界/秘密/未知进程、运行规定 gates；abort 保留失败证据、丢弃本次候选；request-agent 建新 attempt 并注入经确认结论。到期/断线不自动通过或转租；secret-direct 无安全终端则暂停。

## 7. AgentRunner 与辅助消息通道边界

ProofRail 端口接受版本化 context envelope：任务契约、父快照、已确认决策、最新失败证据、未决票据、预算、授权 target、允许副作用及引用摘要；返回外部执行事实和候选引用，不返回可信的 task PASS。每个 attempt 有固定关联身份；启动前持久化请求和授权摘要，重复请求不得重复启动未知状态的代理或应用同一候选。新业务 attempt 产生新身份并计费；恢复必须绑定原 attempt、会话和工作区，无法证明连续性时从持久 envelope 建立新 attempt，不猜测续接成功。

正式 AI 执行通道为 `AgentRunner` 端口的 CLI Agent adapter。adapter 只负责把统一请求映射到经固定版本和摘要绑定的外部 CLI，不把厂商 transcript、私有会话格式或退出码直接提升为核心状态。启动前必须探测并记录：可执行文件身份/版本、非交互调用方式、工作目录绑定、机器可读事件或完整日志、会话创建/恢复、取消与进程树停机、工具/网络/权限约束、费用用量以及无人值守确认行为。任一 required 能力不可验证即 preflight 阻断；不能用提示词承诺替代 OS/runner 限制。

CLI Agent 只在为当前 task/attempt 物化的可丢弃 run-workspace 中工作，不得写源目录、store、策略、门禁实现或接受记录。它可在授权范围内检索/读取文件、直接修改 isolated-workspace、运行允许的开发工具并迭代诊断；不得自行 commit、push、发布、扩大网络/target、批准候选或清除证据。ProofRail 记录启动 argv 的脱敏形式、cwd、环境白名单摘要、agent/config/version 摘要、session ID、进程身份、事件/日志、时间/调用/费用、停止证据以及前后 manifest/diff。机器事件缺失时原始 stdout/stderr 只能作为不透明日志，不能推断工具调用完整性。

代理报告 completed、退出 0 或最终文本只表示外部执行结束。ProofRail 必须先证明进程树停止，重新扫描 workspace，校验范围/秘密/副作用，冻结 candidate，再独立运行声明的 build/test/verify gates；通过后仍进入独立 review/promotion。超时、失联、无法停止、会话恢复失败、日志或费用不完整、外部副作用不明均为 uncertain 并保持暂停，禁止换会话盲重试。AgentRunner 请求、能力报告、事件和 completion receipt 的 wire、canonical/hash、去重与重放规则尚未冻结；实现前须经 ADR、Schema、正反 golden 和 checker 更新。

需要人工输入时，CLI Agent adapter 只能上报结构化 `operator-action-required`；ProofRail 在原子边界停止或暂停受管代理，持久化请求并令 task/step 进入 `WAITING_FOR_OPERATOR`。通知、答复和控制权归还由 ProofRail 自有 CLI/TUI 承担。用户输入必须经 ProofRail 校验并持久化为适用的 operator interaction、review、authorization 或 handoff 记录；终端/聊天自由文本不能直接改变状态或授权。恢复同一 Agent 会话前须重新验证 attempt、workspace、session、上下文摘要、租约和授权；一般澄清问答的 `operator-interaction` wire 仍须先完成 T025 的 Schema、正反样例、摘要和重放契约。

SessionBridge v0.1.1 的 `silent`/`visible`/`auto` 均不是正式 AgentRunner：`silent` 可用于无工具的分类、摘要、计划、结构化分析或 managed-change-set 建议；`visible` 可用于观察、诊断、通知、人工接管及下述受监督黑箱候选流程；`auto` 不进入正式流程，因为其实际路径和能力不确定。任何 SessionBridge 使用仍须内置 v1 文件 IPC 客户端、`legacy=false`、requestId 回执绑定和自身持久幂等；连续辅助会话可使用稳定 `conversationId`，但 history 只作上下文。不得把 visible 投递成功解释为 Agent 完成，也不得在 CLI Agent 失败后自动切换模式继续写入。

`supervised-black-box` 是正式支持但保证降级的候选生产/导入模式，不是 AgentRunner 的执行替代。ProofRail 先物化可丢弃 isolated-workspace、冻结授权 target 和验证计划，再由操作员在 `visible` AI 界面中监督外部代理；该代理及其工具循环均视为不可观测、不可信的外部变更执行者。操作员显式结束并归还工作区后，ProofRail 必须证明已知受管进程停止，重新生成完整 manifest/diff，执行范围、秘密、文件类型/大小及可观测副作用检查，随后独立运行冻结的 build/test/verify gates 和 review/promotion。任何一步不完整即暂停或拒绝；不得从聊天文本、visible 回执或用户口头确认推断执行完成、工具轨迹或安全。

进入该模式前，CLI/TUI 必须展示并持久化显式风险确认：ProofRail 无法证明代理使用了哪些工具/命令、读取了哪些未被 OS 隔离的资源、是否发起网络请求、实际 token/费用/重试、全部后台进程、会话连续性，以及 commit/push/发布等外部副作用。产品只对其自身完成的隔离、归还后产物扫描、独立门禁、证据和接受事务作保证；用户负责选择可信宿主、监督过程、限制凭据/网络和确认未知外部副作用。该责任说明不能豁免源/store 隔离、秘密扫描、独立门禁或评审，也不能把未知项显示为安全。未取得绑定 task/attempt/workspace/风险版本的用户确认时禁止启动；不允许自动 fallback 到此模式。

CLI Agent 退出 0 或报告完成不等于任务通过，只能触发 ProofRail 的独立停机、扫描、门禁、评审和发布流程。

| 外部结果/能力 | ProofRail 行为 |
|---|---|
| CLI Agent completed/exit 0 | 仅记录外部结束；停机、manifest/diff、范围、秘密和独立 gates/review 全部通过后才可推进 |
| CLI Agent timeout/失联/恢复失败/无法停机 | 记录 uncertain 并暂停核对；不得更换 session 或重复启动来猜测结果 |
| CLI Agent 请求额外工具、网络、target 或确认 | 结构化进入授权/操作员流程；未获有效记录不得执行 |
| SessionBridge silent status=ok | 按声明的辅助输出 Schema 检查；可形成分析或 managed-change-set 建议，绝不直接 PASSED |
| SessionBridge visible 已投递 | 仅是 UI 投递/人工观察事实；黑箱模式还须显式归还工作区并从磁盘重建候选，仍不证明模型、工具或任务完成 |
| SessionBridge auto | 正式流程拒绝；不得以自动回退改变执行能力或证据语义 |

SessionBridge 成功缓存和 CLI Agent 私有会话状态都不是持久 exactly-once。ProofRail 自己持久化投递、进程、结果及应用幂等记录；上下文历史、agent transcript 和 UI 内容均不可替代 run manifest、event、receipt、workspace manifest 与 gate evidence。Secret 扫描/脱敏由 ProofRail 保证，不假定外部 Agent 或 SessionBridge 已实现。

文件队列与 IPC 共享 ProofRail request/claim/result 信封，CLI 是消费者形态，不是第三套业务语义。文件队列
每文件恰好一条 canonical JSONL 记录，通过 request 原子移入 inflight 获取 claim；不可变 generation 文件与
`(claimId,generation)` fencing 约束续租、接管和 result。禁止多进程追加或覆盖正式文件。SessionBridge
传输 ID 只在 dispatch receipt 映射，不复用其单槽文件作为持久事实。takeover receipt 绑定旧/新 fencing、
停机证明与授权证据；只有 granted receipt 可被新 claim 引用。平台原子/崩溃边界仍须 T004 实测。
无 IDE、无 AI 测试用确定性 fixture consumer，不调用收费模型。

`signature-receipt.schema.json` 以 detached Ed25519 receipt 为跨主机 adapter 记录、bootstrap 产物和发布
校验和清单提供身份绑定。签名覆盖带 purpose/subjectKind/subjectHash/signer/keyId/公钥指纹/信任策略的
canonical statement，receiptHash 再覆盖 statement 与签名字节；cross-host/bootstrap/release 与其对象
类型和主体严格配对，64 字节签名只接受无 padding base64url。Schema 不证明密钥可信、当时未撤销或验签
成功；这些由独立 crypto validator、信任策略与验证证据判定。该能力在 S1 为可选增强，不是 T017/S1 exit 的强制发布门禁。

`lifecycle-record.schema.json` 的 release 记录始终强制 binary、SBOM、许可证和校验和清单摘要；`signatureReceiptHash` 可为 null。无签名时 `revocationFreshness` 必须为 `unknown` 且不得携带撤销核验证据；存在签名 receipt 时仍按上述严格契约验证。SHA256SUMS 只证明内容完整性，不证明发布者身份。

## 8. 错误与验收边界

持久错误区分 schema/static/policy/capability/transport/execution/integrity/storage/review/operator/internal 类原因。
`error.schema.json` 冻结 37 个首版 code，并逐组绑定 category、retryMode 与 suggestedAction；未知 code 按版本
fail-close。`error-set.schema.json` 绑定非空错误索引和 primaryErrorId；主错误按最早 occurredAt，再按固定
category 次序、code、errorId 决定，后续 cleanup 失败不能覆盖原始失败。CLI 按主错误 category 使用 10–20，
用法错误为 2；不得透传 SessionBridge 或 hook/OS 退出码。引用一致性与排序由 checker 判定，错误不得含秘密。

`ticket-ledger.schema.json` 按 `(runId,fingerprint)` 记录唯一 append-only 账本。指纹只覆盖 error code、稳定
subject 和 failurePoint，不含 message、时间、attempt 或证据摘要；failure 跨 attempt 单调计数且 errorHash
不得重复。低于 reviewThreshold 为 pending-review；达到后须有效 override 才进入有限 override-window；达到
hardBlockThreshold 后在该 run 内不可回退。override 不得提高 hard block，resolved 也不重置预算。Schema
判三类 entry 形状；sequence、阈值关系、计数、授权窗口和唯一账本由 checker 判定。

`repair-transaction.schema.json` 将 Prepare/Inspect/Validate/Promote 固化为带前序摘要的 1–4 项阶段链。
只有前一阶段成功才能追加下一阶段；failed/uncertain 立即终止，completed 要求四阶段全部成功。Prepare
绑定父快照、隔离候选、目标、停机和租约；Inspect 绑定 diff/范围/所有权；Validate 绑定冻结计划和 hook
结果；Promote 再验停机、租约、账本与 Validate 摘要，只把候选装入下一 attempt 工作区。它不产生 accepted
snapshot、不进入 PASSED，也不替代独立 review 与 snapshot promotion。

契约验收不是“JSON 能解析”：须验证结构、语义、运行行为、崩溃点、兼容及安全。对应测试目录、用例与阶段见 [TEST_STRATEGY.md](TEST_STRATEGY.md) 和 [DEV_PLAN.md](DEV_PLAN.md)。T003 已运行独立结构/静态语义 validator；运行行为、崩溃点和平台能力仍由 T004/S1 gate 验证。

## 9. 产品与生命周期契约（RFC §19）

下列语义及 wire 字段已由 T002 冻结，并由 T003 补充正反样例与独立验证；未知类型拒绝写入。

| PC | 输入/权威 | 输出/不变式 | 验收 |
|---|---|---|---|
| PC-01 | 用户配置+静态能力声明 | 计划摘要、权限/预算/unknown；不执行 hook/模型/网络/凭据调用；版本探测须另授权 | AT-16 |
| PC-02 | 已接受 snapshot+review+完整引用 | 导出清单绑定父/目标/内容哈希/排除项；仅新目录，临时写后核验完成；失败不输出完成回执 | AT-17 |
| PC-03 | 可信主体签发的范围/manifest/预算/期限 | 新副作用和接受前校验；撤销阻断并请求受控停止，已有副作用不被“抹除”；waiver 不改硬安全门禁 | AT-18 |
| PC-04 | 可信策略定义 effect 类别+runner 可实施能力 | 只读/本地可丢弃/外部写分离；S1 外部写拒绝，未知副作用暂停对账；诊断不变更锁/日志 | AT-19 |
| PC-05 | 所有者预算、价格/计量来源、调用身份 | 调用前持久 reserve，结算去重，unknown 保留占用；共享上限跨 run 生效；不得用重启释放未结算预留 | AT-20 |
| PC-06 | 已核验引用闭包、支持矩阵、保留授权 | 备份包含对象+事件+引用但不含秘密；新 store 恢复；不改旧格式；卸载默认不删证据/共享工具 | AT-21 |

成本契约包含币种、价格版本、estimated/observed/unknown 和结算依据；订阅按调用/token cap 授权时必须标明“不保证货币账单”。预留与实际费用的差异留痕，重复回执不可重复结算；超供应商可控范围的硬币种保证不得宣称成立。

导出/备份 reader 必须验证对象闭包、路径/类型安全、版本、哈希和授权；导出完成不提升原 task/chain 状态，更不代表产品发布。对于活动 run，一致性未知时拒绝备份，不返回可恢复的假成功。处置记录不含被删除秘密值；删除审批和审计保留冲突进入人工决定。

对应结构分别为 `plan-preview.schema.json`、`export-record.schema.json`、`authorization-record.schema.json`、
`effect-record.schema.json`、`cost-ledger.schema.json` 与 `lifecycle-record.schema.json`。跨记录结论写入
`verification-report.schema.json`；固定检查顺序覆盖事件链、引用闭包、身份/attempt/snapshot、证据/review/
promotion、配置/注册表、change-set/repair/ticket、queue fencing、签名信任和生命周期授权。failed 或
incomplete 均 fail closed；报告自身不改变状态或授权。