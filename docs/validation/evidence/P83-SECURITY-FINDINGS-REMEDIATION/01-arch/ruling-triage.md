# P83 ⑦ 终审 9 项 · 裁决分流（02-ruling-triage）

- 角色：① 架构（deep-reasoner，DeepSeek V4 Pro Max）。本文件是唯一落笔；对源代码与暂存区只读，不设计修复、不改代码、不重新设计切片。
- 输入权威：⑦ 终审 `FINAL REVIEW: FINDINGS` 逐字原文（`07-final-review.md`，3 High + 6 Medium，9 项全阻断）；
  P79 三面 Section B 逐字违反句（`docs/validation/s2-pre-independent-security-review.md` §7）；
  `docs/CONTRACTS.md`（:62、:91、:129、:147–148、:150、:212、:214、:216、:242 起、:391）；`docs/SECURITY.md`（:23、:60）。
- 编号映射（按终审顺序）：T1=F2、T2=F3、T3=F6、T4=F1、T5=F4、T6=F5、T7=F7、T8=F8、T9=范围与兼容性。

## 一 逐项裁决表

| # | 发现 | (a) 契约影响 | (b) 性质 | (c) 边界 |
|---|---|---|---|---|
| T1 | F2 | DOES-NOT-TOUCH | F1–F8 REMEDIATION | IN-BOUNDARY |
| T2 | F3 | DOES-NOT-TOUCH | F1–F8 REMEDIATION | IN-BOUNDARY |
| T3 | F6 | DOES-NOT-TOUCH | F1–F8 REMEDIATION | IN-BOUNDARY |
| T4 | F1 | DOES-NOT-TOUCH | F1–F8 REMEDIATION | IN-BOUNDARY |
| T5 | F4 | DOES-NOT-TOUCH | F1–F8 REMEDIATION | IN-BOUNDARY |
| T6 | F5 | DOES-NOT-TOUCH | F1–F8 REMEDIATION | IN-BOUNDARY |
| T7 | F7 | DOES-NOT-TOUCH | F1–F8 REMEDIATION | IN-BOUNDARY |
| T8 | F8 | DOES-NOT-TOUCH | PRODUCT DEFAULT AVAILABILITY DECISION | IN-BOUNDARY |
| T9 | 范围与兼容性 | DOES-NOT-TOUCH | F1–F8 REMEDIATION | IN-BOUNDARY |

## 二 逐项理由

### T1 = F2（dispatcher:153、assembly:1469）
- (a) DOES-NOT-TOUCH。既有句 `CONTRACTS.md:212`：「dispatch 在发布 R 之后、允许 launch 之前必须重新确认授权未被撤销且预算 reservation 仍 outstanding，任一不满足即 fail-closed」。内存副本观察不到新撤销，故「生产入口要求实时账本源」是该句「重新确认」对来源的必然要求；「读完再取时钟」是该句「重新确认」的正确评估次序（现在时钟取在账本读取前，读取期间到期的授权按旧时钟仍 active）。无需措辞、schema、夹具变更。
- (b) F1–F8 REMEDIATION。P79 F2 逐字：「在 launch 前向可信账本读取当前授权与预留；读取失败、撤销或已结算均保留 R，拒绝启动」——「当前」二字已含实时源与读取后再评估。
- (c) IN-BOUNDARY。CLI 启动路径已持有账本路径；assembly 把「路径必填（内存替身仅测试）」与 dispatcher 把取钟移到账本读取之后，均为现有子系统内接线。

### T2 = F3（agent:801、supervised:995）
- (a) DOES-NOT-TOUCH。既有句 `CONTRACTS.md:147–148`：「`reviewMode=policy` 必须绑定已批准的 `policyHash`……`evidenceRootHash`/`policyHash` 引用是否真实存在……均由独立 checker 判定」。「未解析独立获批的策略对象」违反的是「已批准」三字（现在只验 hash 形状）；「执行者写不到该来源」是「独立 checker」的独立性落实。`CONTRACTS.md:62` 的 `supervised-risk-confirmation` 只字未提其可作为发布策略批准——把风险确认当批准源是对该行语义的误读，拒绝之即恢复该语义。
- (b) F1–F8 REMEDIATION。P79 F3 逐字：「发布前从执行者不可写的可信来源验证策略批准与候选证据根；没有有效绑定则拒绝，或停在待人工批准」——拒绝/停在人工是原文给出的合法处置，不是新的产品默认。
- (c) IN-BOUNDARY。授权账本与冻结 run manifest 的策略摘要均已存在，策略对象可从冻结策略解析并与 grant `PolicyHash` 对照；无需新记录类型或新 CLI。

### T3 = F6（ai.go:155、ai.go:435）
- (a) DOES-NOT-TOUCH。被违反句是 `SECURITY.md:60`（「网络默认拒绝，例外授权须列 host/port/purpose、提供方数据策略及期限」），不在 CONTRACTS。`CONTRACTS.md:214` 已要求 profile 只存不透明 `secretRef`、token 只经 SecretStore；`authorization-record.schema.json` 的 grant `evidence` 为开放 hashSet（minItems 1），完整 profile 摘要与 `secretRef` 身份摘要可作 evidence 承载，无需 schema/夹具变更。「固定可信审批来源」实现上沿用 `CONTRACTS.md:242 起` replay store 的既有原则「由 durable run root 确定性推导、不再接受任意 caller 自选 root」。
- (b) F1–F8 REMEDIATION。P79 F6 逐字：「解析凭据前，由可信本机控制路径校验已批准的 profile、凭据引用、接收端与子进程身份绑定；不符即不读库、不启动」。终审 F6 信任判定自述：「对原有 run CLI『调用者就是受信操作者』的模型，这不是实质变化」。
- (c) IN-BOUNDARY。ai.go / adapters / console 内接线 + 反例；「同一文件对象」绑定复用 `evidence.BoundFile`/Identity 既有原语。

### T4 = F1（runner.go:226、Darwin:26）
- (a) DOES-NOT-TOUCH。既有句 `SECURITY.md:23`：「限制必须由 runner/OS 能力落实；不能仅靠提示词」——执行前绑定根、之后沿用同一对象是该句落实；无产物不打开根不涉及该句。
- (b) F1–F8 REMEDIATION。P79 F1 逐字：「产物读取绑定授权工作区目录句柄，无法安全打开即拒绝收集」。macOS 无产物放行是恢复 P83 之前已合法、且本不读取任何文件的行为（回退本片引入的回归）。
- (c) IN-BOUNDARY。gates/runner.go + evidence/bounddir；无产物短路是现文件内改动。

### T5 = F4（runtime_agent.go:1885、lease.go:260）
- (a) DOES-NOT-TOUCH。既有句 `CONTRACTS.md:150`：「promotion-receipt.schema.json 将父/候选 snapshot、evidence root、review receipt、停机与租约证据绑定为单个发布事务」；`CONTRACTS.md:91`：「恢复先证明旧写者失效并取得有效 fencing token」。换 PID 须证前任已停、按 hash 保留每代原件（旧 hash 回执可独立复核）是这两句的落实。
- (b) F1–F8 REMEDIATION。P79 F4 逐字：「发布前核对当前工作区的排他租约或等价写入隔离证明，将其真实记录绑定到本次 attempt；缺失或不一致则拒绝」。
- (c) IN-BOUNDARY。租约文件是 `internal/guard` 的 store-local 记录，不在 CONTRACTS wire 对象清单，也不进 recordHash 集；代次归档是 store 布局调整，不改 schema。

### T6 = F5（runner.go:382、policy.go:50）
- (a) DOES-NOT-TOUCH。既有句 `CONTRACTS.md:391`（PC-04）：「只读/本地可丢弃/外部写分离；S1 外部写拒绝，未知副作用暂停对账」——分类字段只取冻结策略正是「分离」要求；「HookDefinitionHash 只校验形状」违反的是 run manifest 对 hook 摘要的绑定语义（§3）。
- (b) F1–F8 REMEDIATION。P79 F5 逐字：「执行前从冻结 hook 策略导出有效副作用分类；分类缺失、矛盾或不允许时不能进入 `Execute`」。
- (c) IN-BOUNDARY。gates/runner.go、policy.go 内核对真实 hook 定义摘要。

### T7 = F7（process_windows.go:518、secret_input_windows.go:51）
- (a) DOES-NOT-TOUCH。既有句 `CONTRACTS.md:216`：「adapter 只接收解析后的短暂字节副本，使用后清零；系统错误只返回稳定分类」。
- (b) F1–F8 REMEDIATION。P79 F7 逐字：「对原生环境块与编码中间缓冲保留所有权，成功或失败均及时清零，并核对凭据读取与终端输入的临时副本」。
- (c) IN-BOUNDARY。预留终止符空间、增长副本清零、测试直接观察 scratch/旧数组，均在现两文件内。

### T8 = F8（executor:175、executor:140）
- (a) DOES-NOT-TOUCH。既有句 `CONTRACTS.md:214`：token 不得进入配置、摘要、日志、错误或证据。「拒绝证据只用稳定分类、不散列参数值」是该句的直接实现。
- (b) PRODUCT DEFAULT AVAILABILITY DECISION。理由见 §四决策简报：P79 F8 虽含「无法保证则拒绝 BYOK 探针」，但当前不存在任何既有配置或操作员供给步骤能建立「子进程配置禁用未脱敏落盘」的证明（第三方 CLI 可写当前用户可写的任意路径；受控沙箱/日志重定向证明属新子系统），按字面实施 = BYOK 探测默认全拒 ⇒ AI 通道默认 unavailable，改变合法用户默认行为。
- (c) IN-BOUNDARY（拒绝门本身可在 executor 内落地）；但「证明载体」需要新子系统，本片不建，待用户裁定后另片。

### T9 = 范围与兼容性（review.go:145、ai.go:77）
- (a) DOES-NOT-TOUCH。`CONTRACTS.md:147–148` 定义的通用 review receipt 不要求 approve 额外携带 Evidence 绑定项（通用校验新增该要求即偏离）；`CONTRACTS.md:150` 已把证据根绑定放在发布事务层——证明要求属于发布路径，不属于通用回执层。`CONTRACTS.md:214`：「host 认证不得带 secretRef，secret 认证必须带不透明引用与 HTTPS base URL」——不读凭据的无 `baseUrl` host-auth profile 原本合法（无接收端可命名），新增拒绝即偏离。修复 = 回退到契约行为，无需措辞变更。
- (b) F1–F8 REMEDIATION。终审最小整改 3 与 6 已明示：「把新证明要求限于相应发布路径，保留通用人工回执兼容性」「保留非 BYOK 路径」——收窄即整改。
- (c) IN-BOUNDARY。即 T2/T3 的范围收窄，不新增任何能力。

## 三 分桶执行计划

### Bucket A（DOES-NOT-TOUCH + IN-BOUNDARY，本片内可修）：T1、T2、T3、T4、T5、T6、T7、T9
优先级（High 先、回归先）：
1. **T1（F2）**：生产装配账本路径必填（内存替身限测试）；dispatcher 先读两账本、后取时钟，用读后时钟评估。反例：无源、读期间过期、装配后磁盘结算，均断言保留 R、零启动。
2. **T2（F3）**：发布路径从执行者不可写来源解析独立获批策略对象并对照 digest；同 ID 只认最新 grant 记录为当前意图；supervised 发布拒绝「仅有风险确认」。反例：摘要合形但策略未独立批准、旧 grant scope 背书 → 零发布。
3. **T3（F6）**：grant evidence 核准完整 profile 摘要 + `secretRef` 身份 + 可执行文件身份；检查器沿用账本核准摘要，不按路径重取；启动前最后一次核验与 spawn 绑定同一文件对象。反例：自带账本、同地址换 `secretRef`、批准后换文件 → 零解析/零启动；保留非 BYOK 路径。
4. **T9（范围收窄，与 T2/T3 同轮）**：通用 approve 校验恢复契约形状（证据根绑定要求移到发布路径）；无 `baseUrl` host-auth profile 恢复放行（接收端门仅限 secret-auth/BYOK）。
5. **T4（F1）**：无产物直接返回空列表；有产物在执行前绑定可信根并沿用同一对象。反例：根替换、无产物 macOS 放行。
6. **T5（F4）**：发布事务持进程级排他证明；换 PID 须证前任已停；按 hash 保留每代租约原件并绑定 attempt。反例：双进程互斥、续租后旧回执独立复核。
7. **T6（F5）**：`Runner.Run` 核对真实 hook 定义摘要；`RecoveryGuarantee`/`ExternalSystems` 只取冻结策略，scope 不得覆盖；矛盾即拒。反例：元数据冲突零执行。
8. **T7（F7）**：预留终止符空间、清零增长副本与终端输入旧数组；测试直接观察 builder 成功/失败时的 scratch、部分块与旧数组（删掉清零即变红）。

### Bucket B（TOUCHES ⇒ 协议先行）：∅
九项均不要求修改 `docs/CONTRACTS{,_EN}.md` 措辞、schema 或夹具文本；无契约新句。守则见 §五（防翻桶）。

### Bucket C（PRODUCT DEFAULT AVAILABILITY DECISION）：T8
**决策简报（待用户逐项拍板）**：
若按终审要求实施「不能证明该子进程配置禁用未脱敏落盘时拒绝 BYOK」，产品默认行为 = 所有携带用户自有 key 的 AI 可用性探测（Copilot CLI probe）在无法证明子进程禁写未脱敏日志时一律拒绝（fail-closed、零启动），AI 通道默认 `unavailable`，直到证明机制存在为止。操作员侧目前**不存在**任何供给步骤可建立该证明：第三方 CLI 可写当前用户可写的任意路径，本片既无受控日志重定向证明、也无沙箱子系统。两条路径：① 批准默认拒绝（BYOK 探测暂不可用，待未来引入可证明禁写机制后恢复；本片同步落地「拒绝证据只用稳定分类、不散列参数值」半句，该半句无争议、界内可做）；② 不批准默认拒绝，则「证明载体」（强制日志目录指向受控临时目录 + 事后审计，或等价机制）属新子系统，超出「只修 F1–F8，不扩大安全加固范围」边界，应顺延 P84。若既不裁决又默认实施，风险 = 发布后 AI 探测全链路 unavailable，且台账对「默认可用性改变」无用户授权记录。

## 四 依赖与顺序

- **信任源基座（T1 → T2/T3）**：先定 T1 的「实时账本源 + 读完取钟」语义；T2（策略批准源）与 T3（接收端审批源）复用同一条「账本路径固定、执行者不可写、实时重读」规则。三者共享同一判定，不得各自发明。
- **T9 与 T2/T3 同轮**：T9 两半分别是 F3 通用校验过度接线与 F6 无 baseUrl 拒绝过度接线，须与 T2/T3 的范围收窄同时落地，否则通用路径回归。
- **T4/T5/T6/T7 相互独立**，可与上组并行、任意顺序；T5（F4 租约代次归档）与 T2（F3 发布路径）都碰 promotion-receipt 证据绑定，合并回归时核对 `:150` 的收据证据引用不变。
- **T8 不阻塞 A 桶**：其「稳定分类」半句可随 A 桶并行实施；「拒发 BYOK」半句等用户裁定后再落。

## 五 守则（防翻桶）

- 任何 A 桶实现若需要新 schema 记录、新 CLI 命令、新依赖或第三方包，即 OUT-OF-BOUNDARY，停下升级，不自行扩域。
- 若 T3「固定可信审批来源」最终必须落在新位置或新记录类型（超出 `:242 起` 先例的推理范畴），该子项翻 TOUCHES ⇒ 协议先行：先改 `docs/CONTRACTS{,_EN}.md` 及 schema/夹具，再动代码。
- 本裁决不推翻切片边界「只修 F1–F8，不扩大安全加固范围」；T9 的方向是收窄而非新增。
