已完成三份输入（design-r2 全文含 A2 长行尾段、红队 r2 全文、Q1 裁定）的逐字读取与需改短语的全量定位核对。以下为 design-r2b 全文（只读，未写任何文件）。分为两个逻辑批次，两批合计约 9,500 字符（< 10,000 上限）；批次 1/2 = 交付 A，批次 2/2 = 交付 B/C/D/E。

---

# P54 `CONTRACTS-S7-DURABILITY` · ① 架构第三轮修订（design-r2b，只读返回文本）

载体：`deep-reasoner` + `DeepSeek V4 Pro`（只读，2026-11-04，未写任何文件）。输入三件已逐字读取：design-r2.md（含增补节与 A2 长行尾段）、red-team-sol-review-r2.md、ruling-2026-11-04-q1.md。改动点图例：①=红队 §3 High（历史全集绑定）；②=红队 §1 Medium（允许集不一致，design-r2.md:60）；③=红队 §4 Medium（BA3，design-r2.md:109）；④=红队 §3 Low（顺序论证，design-r2.md:95）；⑤=红队 §3 Low（哈希口径，design-r2.md:73 vs :93）；⑥=红队 §4 Low（API 形态，design-r2.md:105/:118）。方括号【】为 r2b 评审标注，非 §7 正式文本。本片边界不变：只改 CONTRACTS §7 的 C6 条款 + Windows 平台实现 + 相关测试；不改三判据本体、Unix 路径、ADR-013 结论、`NewAgentRunnerReplayStoreForRun` 签名、平台等级 unproven；C2 引用一律带完整括注「（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）」。

## 批次 1/2 —— A. §7 修订文本 v3 全文（CN 权威）

**A1. replay store 段最小措辞调整**（与 design-r2 v2 逐字一致，无改动点）：
> （新增）「如某（平台 × 前提）经本节「C6 判定路径」完成协议先行修订并经独立评审后激活预建槽位实现，则该平台**仅 R/C 记录文件**的发布可见性以槽内头标志提交为准，其目录项在首次派发前由装配层预建一次性建立、此后永不重建或替换；未激活的平台与激活前的实现维持本段既有口径逐字不变。」

**A2. 新增段落「C6 判定路径」（bootstrap 预建固定槽，ADR-013）**（插入 CONTRACTS.md:232 段之后）：
> **C6 判定路径**：（C6 定义）C6 指每 requestId 的 R/C 记录文件作为槽文件在首次派发前预建：目录项一次性建立、此后永不重建或替换；发布 = 槽内写 payload + 数据 flush + 槽内头标志提交（payloadLen 0→len）+ 头 flush，可见性以槽内头标志界定而非目录项创建（出处 `docs/validation/t027-windows-publication-durability.md:28`、`docs/ADR_REGISTER.md:23`）。（预建生命周期，Q1=B，2026-11-04 用户裁定）**本裁定推翻 r1 的 High 6 判断**：r1「requestId 在 bootstrap 时刻未知」前提不成立（`runtime_agent.go:148→:181` 次序证明 store 构造时声明式 requestId 全集已确定），r1 的「首次发布前预建」解析作废。预建由生产装配层触发：`ExecuteAgentRunnerRun` 在 `NewAgentRunnerReplayStoreForRun` 返回之后、首次派发之前，以 `agentRunnerAssemblyRecords` 已校验的 requestId 全集调用新增可选方法 `PrebuildReplaySlots(requestIDs []string) error`；不得改 `NewAgentRunnerReplayStoreForRun` 签名，调用方必须传全集、不得分批。预建失败语义：任一槽预建失败 ⇒ 整个 run fail-closed（store 不可用、不进入首次派发）；部分成功不得视为可用；已建成的空槽允许保留（供重开幂等核验），但 run 必须停。重开语义：重开时对已存在槽只核验（magic + 长度 + payloadLen 合法性）不重建；分配记录（预建全集）已持久化而槽缺失 = corruption fail-closed；从未预建的槽缺失 = 正常。适用边界：预建只对生产装配路径要求；包内测试、A2 dispatch 测试等非生产入口不强制预建。（A7 偏离登记）本实现较 A7 原文的偏离点 = 预建由装配层触发而非 bootstrap 程序内部；该偏离不改 C6 的结构性优势（目录项仍在首次派发前一次性建立）。（预建全集持久化）replay root 内落一条 `prebuild-set.jsonl` store-local 记录；重开时声明集必须与该记录逐项一致（同为全集、无增无减），不一致 fail-closed（计划漂移按 corruption）；记录不存在 ⇒ 按首次预建执行（既有槽幂等核验、缺失槽补建，全部成功后发布记录）。【改动点①：红队 §3 High（历史全集绑定）；删「补建出的集合与丢失前逐项相同」「漂移检测的丢失无实际安全后果」，以下列句取代】prebuild-set 丢失时，漂移检测一并丢失；此为有意降级，安全后果 = 计划漂移检测失效，但 R/C 绑定、orphan、runID 校验照常执行，且 Windows unproven 阶段 first-dispatch 被拒，该降级无实际运行风险。接受该降级。本降级经用户 2026-11-04 裁定接受。（接缝范围）仅 R/C 记录采用槽容器：新增槽发布函数只由 R/C 发布点调用，`writeReplayRecordNoReplace` 主路径不动，ownership/launch-intent/launch-identity/terminal-intent/terminal-closure 等 store-local 记录保持原发布与原读取逐字节不变；显式禁止把其它记录迁入槽容器，若后续评估认为必须迁，停下升级用户裁决，不得自行扩域。（槽容器格式）槽文件 = 16B 固定头（8B magic `PRSLOT1\0` + uint64le payloadLen）+ payload（canonical JSONL，与既有 R/C 文件内容逐字节一致）；「一条 canonical JSONL」（`:230`）由 payload 维持，槽头属容器层不属于记录层。已发布 = payloadLen>0 且总长 = 16+payloadLen；空槽 = 总长 16 且 payloadLen=0。（未提交中间态冻结与校验顺序）「写 payload 后、头提交前」的读者可见结果恒冻结为 absent；校验顺序冻结：① 先查 magic（错 = corruption fail-closed）；② 再查总长是否 ∈ {16, 16+payloadLen}——总长 16 且 payloadLen=0 = absent；payloadLen=0 且总长>16 = 未提交中间态 ⇒ absent（不判 corruption）；总长<16，或 payloadLen>0 且总长≠16+payloadLen = corruption；③ 再按 payloadLen 有界读 payload 并严格解码 + requestId 校验，失败 = corruption。（三时点）头写可见点 = 头字段写入即可能对读者可见（不等待 flush #2）；【改动点②：红队 §1 Medium（允许集不一致）】允许集按两场景分别冻结：进程 kill（页缓存仍在）——头写后读者可见恒 `dispatched-unknown`，`absent` 不落该场景允许集；断电恢复（页缓存丢失）——头可能未达盘，读者可见 `absent` 亦落允许集，即 `{absent, dispatched-unknown}`；两场景逐切点允许集见本修订附件 C，验收器必须按场景选表、不得混用。flush #2 耐久点 = 头 flush 成功（阶段 `slot-header-synced`）；解锁点 = flush #2 成功后释放槽锁与 store 锁。（仲裁协议）写者与读者同句柄（`CreateFileW(OPEN_EXISTING, GENERIC_READ|GENERIC_WRITE, FILE_SHARE_READ|FILE_SHARE_WRITE)`）；写者取 `LockFileEx` 排他锁（范围覆盖槽文件全域）、读者同句柄读；锁序固定 `store.mu → 槽文件锁`，禁止反向取锁、禁止多槽嵌套；取锁等待有界（超出按 uncertain 收敛）；解锁 = 同句柄精确范围解锁后关闭；映射视图不属保护范围（本设计不使用映射 I/O）。（失败重试）未提交中间态允许回收：锁内先截断至 16 再写新 payload（长载荷失败后重试短载荷不残留尾部）；已提交槽（payloadLen>0）绝不改写、绝不截断，按 `fs.ErrExist` 走既有碰撞收敛路径。（关系与引用纪律）C6 与 C1/C2/C3 的关系沿用 ADR-013 既有枚举与 B3b 结论账本（`docs/validation/t027-b3b-durability.md:44-48`、`:122-128`），本条款不改写其中任何结论；引用 C2 的候选级 `proven` 必须带完整括注「（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）」，不得截断为「C2 proven」。（前置条件）C6 从结构性候选升为可注入候选的前置条件：协议层完成仲裁原语替换——槽文件以 no-replace 一次性建立且永不删除、替换或追加，槽内头标志只允许 0→1 单向翻转并以槽级排他锁仲裁并发写者；该替换必须先修订本节并经独立评审后方可进入注入取证。（逐项取证）(a) 可复现注入或对照证据——在同一（平台 × 前提）装置上对槽预建点与各发布阶段注入，全部有效轮落在允许集内、证据包可由轮次记录与 journal 复算、无真孤儿；(b) 区分力门——同一装置、同一写入器与同一记录级载荷下，负对照（槽创建省略耐久原语或省略任一 flush）出现 ≥1 次丢失，正对照（完整 C6 序列）零丢失，门不成立即封顶 `unresolved`；(c) 机制证据——每轮记录机制令牌（槽创建原语、flush 调用位置与次数、头提交写），变异臂（编译期降级原语）下装置必须拒绝把「机制未执行的存活」兑现为证据。任一判据缺失即 `unresolved`，存在可归属反例即 `disproven`；三判据同时满足仅构成该（平台 × 前提）下 C6 的候选级结论，不自动改变平台发布耐久等级（三判据本体与解锁边界沿用既有句，原文不动）。（平台等级分层）C6 取证完成前，Windows 平台发布耐久等级保持 `unproven`，first-dispatch、launch-intent 与终局发布继续按 `unproven` 拒绝；取证完成后按本节既有三判据判定并受解锁边界约束——「落地 C6」本身不改变任何等级。（非 Windows 零影响）Unix 路径与既有 proven 实例不受本条款任何影响；C6 仅描述 Windows 分支的可注入判定路径；未激活构建的构造器与发布路径的逻辑路径与执行行为与现状等价（新增守卫行不改变未激活构建的行为）。

> A2 中除【改动点①②】两处外，其余全部子句与 design-r2 v2 逐字一致；「降级面封闭、无实际运行风险」为唯一安全措辞口径（全文档不再出现「无实际安全后果」）。

## 批次 2/2 —— B/C/D/E（续批次 1/2；B 为增补 A/B 修订点，旧→新）

### B. 补充 A/B 的修订文本

**补充 A · A-1 哈希口径（改动点⑤，红队 §3 Low，design-r2.md:73 vs :93）**
旧：「不进 wire、不参与任何 recordHash、不属于 R/C 已发布记录集」
新：「不进 wire、不参与 R/C 的 recordHash、自身 recordHash 按 A-4 校验、不属于 R/C 已发布记录集」——其余字句逐字不变。

**补充 A · A-3 ③（改动点①，红队 §3 High）**——删除两句（「补建出的集合与丢失前逐项相同」「漂移检测的丢失无实际安全后果」），以下列句逐字取代；① ② ④ 逐字不变：
> ③ **降级面封闭、无实际运行风险**：prebuild-set 丢失时，漂移检测一并丢失；此为有意降级，安全后果 = 计划漂移检测失效，但 R/C 绑定、orphan、runID 校验照常执行，且 Windows unproven 阶段 first-dispatch 被拒，该降级无实际运行风险。接受该降级。本降级经用户 2026-11-04 裁定接受。

**补充 A · A-4 顺序论证（改动点④，红队 §3 Low，design-r2.md:95）**——整段替换：
> 顺序论证（先建全槽、后发记录）**限定为：在首次预建 + 当前单阶段记录格式下，「先建全槽、后发记录」是避免未完成预建被误判 corruption 的顺序**，不作唯一安全序的绝对主张；论证前提「无首派发、无 R/C 已发布」仅首次预建成立——记录丢失后重开时可能已有已发布 R/C，该情形不套用本顺序论证，按 A2 改动点① 与 A-3 缺失分支语义处理（已发布槽只核验、缺失槽补建）。程序执行先后不等于断电后的持久化顺序；反序的代价是「记录写成、槽未建全」崩溃窗口在重开时命中二.4 触发前提（记录在 + 槽缺）制造假 corruption（误停），不得泛化为一般性不安全。其余论证不变：部分失败 ⇒ 无记录 ⇒ 重开按「从未预建」走幂等补建 = 正常，且该 run 已按预建失败语义 fail-closed；正序把该窗口整体折叠进安全的「记录不存在」分支。记录自身发布失败（碰撞不收敛等）同按预建失败语义 fail-closed；重开按记录存在与否分别收敛。

**补充 B · BA3（改动点③，红队 §4 Medium，design-r2.md:109）**
旧：「`State`/`RecordRequest`/`RecordCompletion` 在调用前后对任意 requestId 的返回完全一致（等值重放）」
新：「BA3：从两个等价初始状态（同一 store 逐字节快照或重建）分别运行**对照臂**（直接执行操作序列 S，含 `State`/`RecordRequest`/`RecordCompletion` 的任意合法序列）与**调用臂**（先 `store.PrebuildReplaySlots(A)` 再执行同一 S），比较两臂对同一操作序列的返回序列与落盘结果（replay root 目录树逐字节）一致。`RecordRequest`/`RecordCompletion` 是会写入的接口，**不得**当作无副作用探针在同一 store 上前后两次调用做「调用前后一致」的等值重放（首写 first-dispatch 与重放 unknown-block 本应不同），必须以两臂独立初始状态比较。」BA1/BA2/BA4/BA5 逐字不变。

**补充 B · API 形态与执行路径（改动点⑥，红队 §4 Low，design-r2.md:105/:118）**
B-1 定性句旧：「因此这是「新增方法」……不是「新增守卫行」（未在既有代码路径上插入条件判断）；区分点正是 ②……」
B-1 新（分三层分别说明，删除单句概括）：
> 因此 API 形态分三层陈述：① **新增 API** = 共享层新增导出符号 `PrebuildReplaySlots`，方法体在接缝为 nil 时仅「读接缝 → 判 nil → 返回 nil」，不触碰任何既有状态；② **装配调用** = `ExecuteAgentRunnerRun` 新增一行调用 + 一个错误检查分支，该分支在 Unix 上不可达；③ **不可达错误分支** = nil 接缝恒返回 nil ⇒「err 非 nil ⇒ fail-closed」在 Unix 上不可达、不改变错误返回面。本论证只主张**限定观察面**（store 字段、replay root 目录树逐文件字节与条目集合、外部文件、错误返回面）上的行为等价，不主张指令级或源码级保证。

B-2 旧句：「收益 = 指令层零调用，但可观察行为与共享调用点完全相同（同 B-1 判定的零差异）」
B-2 新句：「收益 = 调用点隔离于激活构建，在限定观察面上与共享调用点完全相同（同 B-1 判定的零差异，删除『指令层零调用』表述）」；并在 B-2 末尾追加：「本 B-2 全部等价主张限于限定观察面，不含任何指令级保证。」其余字句（含「默认维持『调用点共享、行为等价』」与替代形态取舍）逐字不变。

**B-3 与 A-2、A-3 ① ② ④、A-4 表、B-1 ①②③ 与 BA1/BA2/BA4/BA5 均逐字保留。**

### C. 允许集表（改动点②配套；按单槽自身写管道计，R 槽与 C 槽各自适用；盘点 = kill/crash 后恢复路径读取结果；未进入本槽发布流程的槽恒 absent，不属于任何 cut）

| 切点 | 进程 kill（页缓存仍在）允许集 | 断电恢复（页缓存丢失）允许集 |
|---|---|---|
| `slot-created`（16B 空槽已建） | {absent} | {absent}（目录项差分丢失仍读作 absent） |
| `slot-written`（payload 已写、头未提交） | {absent}（payloadLen=0 ⇒ 中间态冻结为 absent，不判 corruption） | {absent}（同上，落盘残留亦按长度分类判 absent） |
| `slot-data-synced`（flush #1 完成、头未提交） | {absent} | {absent} |
| `slot-header-written`（头已写、flush #2 未执行） | {dispatched-unknown}（头写经页缓存即时可见；absent 不落允许集） | {absent, dispatched-unknown}（头可能未达盘 ⇒ absent；已达盘 ⇒ dispatched-unknown） |
| `slot-header-synced`（flush #2 完成） | {dispatched-unknown} | {dispatched-unknown} |

表下注释（2026-11-04 用户建议追加，不改表本体）：corruption 不在允许集；验收器见到 corruption 一律判 fail-closed。

规则：同一（切点 × 场景）只允许一个允许集；验收器先定场景（进程 kill 注入 vs 断电模拟）再选表，不得在同轮判定中混用「头写后恒 dispatched-unknown」与「absent 亦允许」。本表与 A2 v3「三时点」句一致，不改变三判据本体。design-r2 F(a) 中「该轮盘点读者可见 ∈ {absent, dispatched-unknown} 均落允许集（R 允许集：头写前恒 absent、头写后恒 dispatched-unknown）」一句由本表与 A2 v3 三时点句取代，实施轮据此同步 F(a) 措辞。

### D. 逐条闭合表

| 红队发现 | 级 | r2b 闭合位置 |
|---|---|---|
| §3 High 历史全集绑定（design-r2.md:80） | High | A2 v3（预建全集持久化）改动点① 降级句 + 裁定留痕句；补充 A-3 ③ 重写 |
| §1 Medium 允许集不一致（design-r2.md:60） | Medium | A2 v3（三时点）改动点②；附件 C 允许集表逐切点自洽 |
| §4 Medium BA3 等值重放错误（design-r2.md:109） | Medium | 补充 B·BA3 重写为两臂对照（改动点③） |
| §3 Low 顺序论证适用域过宽（design-r2.md:95） | Low | 补充 A-4 顺序论证限定重写（改动点④） |
| §3 Low 自身哈希口径冲突（design-r2.md:73 vs :93） | Low | 补充 A-1 口径句修订（改动点⑤） |
| §4 Low API 形态与执行路径混写（design-r2.md:105/:118） | Low | 补充 B-1 三层分述 + B-2 删「指令层零调用」（改动点⑥） |

注：红队 §5 Low（判定表注记，contracts-s7-durability.md:66）属验证文档而非设计文本——本设计 A2（关系与引用纪律）与 F（升级路径）已带完整括注，判定表处括注由后续修订轮按同一括注补全，不属本片设计修订范围。

### E. 结论

`DESIGN: DONE`

（风险登记：本轮为只读设计修订，无 diff、无 build/test/lint 日志；r2b 全部改动限于用户裁定的 1 High + 2 Medium + 3 Low 六个点，未触及本片边界外任何语义。）