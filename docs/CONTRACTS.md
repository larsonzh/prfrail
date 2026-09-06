# ProofRail 架构契约规范

[English](CONTRACTS_EN.md)

日期：2026-09-06；S0 设计约束草案。权威：[RFC](RFC-proofrail-unattended-ai-engineering-product.md) §9-13、§16-17。本文合并 Schema、快照/证据、hook、票据/修复、adapter 五个协议包，减少重复阅读。

## 1. 契约成熟度与兼容

下表“必须”均来自 RFC；“待冻结”是生产实现阻断项。T002 已建立首批结构 Schema，
但本文不是完整 JSON Schema，也不能代替独立 validator。新增 wire 字段、枚举、canonical 算法和错误码
必须先经 [ADR](ADR_REGISTER.md) 更新 RFC，再建立 Schema 与正反黄金样例，最后写 Go 代码。

| 契约 | 已确定的语义 | S0 待冻结产物 |
|---|---|---|
| 配置/任务/步骤 | Draft 2020-12；`schemaVersion=1.0.0`；严格 chain/task/四类 step；workspace/component/language scope 拓扑 | harness/toolchain 注册表、引用与合并规则、其余 schema |
| change-set | 顺序内存文本预验证、first-fail-stop、marker/精确断言、整组事务 | 操作枚举、换行语义、前后哈希、重复 marker 规则 |
| snapshot/evidence | SHA-256 内容寻址、父快照/候选/证据绑定、不可变 | 路径 canonical、hash 输入编码、manifest/receipt schema |
| ticket/repair | 分类、指纹预算、租约、Prepare→Inspect→Validate→Promote | 指纹算法、阈值、ledger 格式、完整转换表 |
| adapter/context | 严格 request/claim/result 信封、recordHash、幂等键、generation fencing、dispatch/takeover receipt、最小上下文 | 跨记录 checker、签名 receipt、平台原子能力实测 |

各持久对象独立 `schemaVersion`，不与 CLI 版本或 SessionBridge schemaVersion 混为一谈；首版为
`1.0.0`，未知主版本拒绝，未知 minor/patch 默认只读拒绝。仓库、Schema、运行时与 wire JSON 统一为
UTF-8 无 BOM + LF；canonical 字节还不得含格式空白或尾随换行。仅明确验证 BOM 输入兼容性的 fixture
可带 BOM，且不能复用为业务对象。

## 2. Schema 与静态检查责任

| 对象 | 最小语义 | 非法/动态约束 |
|---|---|---|
| chain definition | `chain.id/profile`、至少一个有序 task、最小 workspace、可选 documentation | task/step ID 跨数组唯一、引用及依赖无环由 checker 判定 |
| run manifest | 冻结后的有效配置、默认值来源和运行绑定 | 运行中不得原地改计划；确切字段待后续 T002 切片 |
| task | 稳定 id、非空 steps、review、documentationPolicy | id 重复、skip-on-fail 未声明 failureIndependent 拒绝 |
| step | id、kind=code/build/verify/noop；code execution 默认 autonomous | 空 steps、未知 kind 拒绝；noop 必须 reason 且不能启动 hook/agent |
| hook | 独立严格注册表；process/container 判别联合、逐项参数、失败/超时/资源/网络/产物策略 | ID/引用、工具摘要与能力由 checker/preflight 判定；执行结果进入 receipt |
| target | 独立严格注册表；稳定 id、class、access、唯一路径 glob；generated 单向来源 | 越界、重叠写、循环生成、未知引用由 checker 拒绝 |
| component/language scope | component root；scope 的 language/harness/toolchain/target ID；依赖；step 取 scope 交集 | ID 唯一、引用存在、无环、harness 可用；同 root 允许不同唯一 target 所有者 |
| state event | append-only envelope；chain/task/step 实体；严格前后状态；序号、前序 hash、证据与原因 | 跨事件序号/hash/实体状态连续性、主体授权和证据存在性由 checker 判定 |
| error | 严格 category/code 前缀、subject、evidence、脱敏 message、retryMode、suggestedAction | 引用、默认重试限制、主错误选择与具体 code 目录由 checker/result 契约判定 |
| documentation | required/if-affected/optional/forbidden；impactRules | 命中规则缺有效 doc diff、生成物陈旧拒绝；空白/mtime 不算更新 |
| handoff | execution、允许写范围、inputPolicy、deadline、returnActions、hooksAfterReturn | 单写租约、停机、新鲜 manifest；归还后复检不能由 Schema 单独判定 |
| review/waiver | approve/reject/waive；主体、原因、策略依据、有效期、绑定候选 | 同主体自批、过期、候选变化、范围不符均阻断 |

三层检查不可合并：JSON Schema 判结构；语义 checker 判引用、所有权、循环、静态策略；运行期 gate 判文件 diff、进程、能力、秘密、预算与评审。非法黄金样例标注拒绝层，不能声称 JSON Schema 会检查实际文件或进程。

`schemas/` 已创建 chain、hook、target、workspace、state-event、error、adapter-envelope、adapter-receipt、
snapshot-manifest、evidence-manifest、review-receipt、promotion-receipt 十二份结构 Schema；`testdata/contracts/valid/` 与
`testdata/contracts/invalid/` 仍待 T003 创建并由独立 validator 执行。每例携带 fixture ID、契约版本、
expected accept/reject、拒绝层和原因。必须含三任务链、四 kind、C+Go、C+JavaScript+Python、人工交接、
代码文档协同；每一条件至少一正一反。文档中的路径/片段是设计输入，不是已验证可运行样例。

## 3. 配置解析与冻结

profile 只提供默认值；task 显式值覆盖链默认；step 显式值覆盖 task 默认。列表覆盖/合并、空值含义、workspace.toml 与 proofrail.toml 的责任划分尚待 ADR-002 冻结。`config explain` 必须能输出值与来源，不允许悄悄合并数组造成更多权限。

解析完成再按 RFC 8785 JCS 计算 canonical run manifest 摘要，以 UTF-8 无 BOM 字节输入 SHA-256，
文本表示为 `sha256:` 加 64 位小写十六进制。运行期心跳、credentials、机器探测不得回写用户配置。
父 snapshot、hook 摘要、adapter/harness 版本和已授权策略绑定 run。变更配置只能新建 run，或走 RFC
明确的暂停、审批和新 manifest 流程；不能覆盖旧事实。

RFC §16.9.3.1 已冻结 `jcs-001` 与 `state-event-001` 两条字节级向量，分别覆盖 JCS 排序/转义/Unicode
及 state-event 域分隔摘要。T003 仍须把它们固化为独立 fixture，并以非核心实现复算；文档中的摘要值
不能替代独立 validator 证据。

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
解析策略、父/候选 snapshot，并引用状态事件、错误、adapter、hook、artifact、change-set、diff、ticket 与
handoff 对象摘要。它不内嵌日志，不从对象存在推断成功，也不包含 review/waiver/promotion receipt；后者
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

managed-change-set：读全部目标→按全局顺序在内存模拟→检查 marker/前置哈希/断言/边界→持久 journal→逐文件原子替换→写后全量核验→记录 receipt；失败整组回滚，回滚本身失败则隔离并暂停。禁止覆盖外部修改，禁止运行进程存活时恢复。

isolated-workspace：直接工具写入仅限可丢弃运行目录；记录完整前后 manifest/diff/进程/gates；接受边界仍为整个 task。不能把部分语言、部分文件或“代码成功但文档失败”的子集发布。

## 5. Gate Hook 协议

kind 使用 precheck/build/test/verify/review/cleanup；按序执行。`executable + args[]` 不经过隐式 shell，cwd 必须落在授权运行工作区；环境继承白名单，网络默认拒绝，资源/输出/时长必须可执行限制。需要 shell/container/remote/PTY 是独立能力声明，不支持则拒绝而非弱化策略。

结果至少证明 hook 身份/摘要、起止、退出状态、stdout/stderr 摘要、产物引用/哈希、策略及失败原因；确切 wire 字段在 schema 冻结。启动失败、超时、进程无法停止、产物缺失、扫描失败不同于 exit 0。`warn` 只有显式非阻断策略才能继续且不得冒充 PASS；retry(n) 消耗统一预算；cleanup 失败不能抹去原失败。

生成 hook A：只从锁定模板生成参数化命令→语法/模板/依赖/危险能力校验→独立批准→哈希绑定挂载。generated 默认禁用；生成器不能修改自己的验证规则或批准自己。任意生成脚本 B 不进入 S1。

## 6. 票据、修复与交接

ticket 分类沿用 task-static/code-fix/noncode；上下文包含所属 run/task/attempt、失败证据、允许动作、预算和租约。相同指纹预算按 pending_review→override_window→hard_block 限制，具体阈值/有效修复证明必须先冻结，禁止模型自行重置预算。

Prepare 在隔离区建立候选与父摘要；Inspect 比较范围/diff/所有权；Validate 运行冻结门禁；Promote 再校验停机、租约、父/候选/验证摘要与授权，原子提升并写 promotion receipt。任何步骤变化都令后续陈旧验证失效。评审拒绝后修复必须新 attempt，不能改写旧 receipt。

handoff：停止受管写者→flush journal→捕获 manifest→WAITING_FOR_OPERATOR→人工独占写租约。complete 收回租约并检查越界/秘密/未知进程、运行规定 gates；abort 保留失败证据、丢弃本次候选；request-agent 建新 attempt 并注入经确认结论。到期/断线不自动通过或转租；secret-direct 无安全终端则暂停。

## 7. 代理与 SessionBridge 边界

ProofRail 端口接受版本化 context envelope：任务契约、父快照、已确认决策、最新证据、未决票据、预算、允许工具及引用摘要；返回候选变更/证据，不返回可信的任务 PASS。每个 attempt 有固定关联身份；投递前持久记录，重复请求不能重复应用同一 change-set。传输重试复用 requestId；新业务 attempt 产生新身份并计费。

SessionBridge v0.1.1 是独立外部依赖，冻结消费其公开 v1 契约：明确 `mode=silent`、`legacy=false`；使用同一 channel directory 和目标实例；`requestId` 绑定回执；写命令先清理同 ID 陈旧结果，再同目录原子写。准确 wire 以该产品 RFC、黄金样例及客户端实现共同核对，不能直接复制 RFC 中省略字段的说明片段。

| 外部结果/能力 | ProofRail 行为 |
|---|---|
| status=ok 的 response 文本 | 按输出 schema 检查，再做 checker/gates/review；绝不直接 PASSED |
| busy | 有限退避、同 requestId 重试，消耗等待预算；每目标响应槽串行化 |
| timeout/poll_timeout/失联 | 记录结果不确定并暂停/核对；不得假设未执行而无限重发 |
| lm_api_unavailable/未知状态/错误版本 | 显式失败/暂停；仅允许预授权文件队列替代，无 auto/visible/剪贴板兜底 |
| history 被截断/宿主重启 | 从 ProofRail envelope 重建事实，不把压缩聊天记录当完整补丁或决策 |
| silent 只有文本 | managed-change-set 候选输出；不得声称有 IDE 文件编辑/命令执行能力 |

SessionBridge 成功缓存是内存、限时机制，不是持久 exactly-once；ProofRail 自己持久化投递/结果及应用幂等记录。完整补丁轮次按依赖契约使用 noCompress 并独立验证输出完整性。不要借历史压缩绕过 token/费用上限。Secret 扫描/脱敏由 ProofRail 保证，不假定 SessionBridge 已实现。

文件队列与 IPC 共享 ProofRail request/claim/result 信封，CLI 是消费者形态，不是第三套业务语义。文件队列
每文件恰好一条 canonical JSONL 记录，通过 request 原子移入 inflight 获取 claim；不可变 generation 文件与
`(claimId,generation)` fencing 约束续租、接管和 result。禁止多进程追加或覆盖正式文件。SessionBridge
传输 ID 只在 dispatch receipt 映射，不复用其单槽文件作为持久事实。takeover receipt 绑定旧/新 fencing、
停机证明与授权证据；只有 granted receipt 可被新 claim 引用。平台原子/崩溃边界仍须 T004 实测。
无 IDE、无 AI 测试用确定性 fixture consumer，不调用收费模型。

## 8. 错误与验收边界

持久错误区分 schema/static/policy/capability/transport/execution/integrity/storage/review/operator/internal 类原因，
`error.schema.json` 冻结 category/code 前缀、对象、证据、重试模式和建议动作；CLI 按 category 使用 10–20，
用法错误为 2。不得透传 SessionBridge 0/1/2/3 或 hook/OS 退出码。transport 不确定结果必须先 reconcile，
不能盲重试；多错误主次顺序仍待 result/receipt Schema 冻结。错误不得包含秘密。

契约验收不是“JSON 能解析”：须验证结构、语义、运行行为、崩溃点、兼容及安全。对应测试目录、用例与阶段见 [TEST_STRATEGY.md](TEST_STRATEGY.md) 和 [DEV_PLAN.md](DEV_PLAN.md)。本轮未生成或运行完整契约 validator，S0 §17.2 仍不通过。

## 9. 产品与生命周期契约（RFC §19）

下列是 S0 必须冻结的语义，不是已确定 wire 字段。T002/T003 为每种记录补版本、必填字段、大小限制、正反样例与独立验证；未知类型拒绝写入。

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