# P54 `CONTRACTS-S7-DURABILITY` · ① 架构第二轮设计（design-r2，只读返回文本）

载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读，2026-11-04，未写任何文件）。批次 1/1（单批）。输入包五件已逐字读取：red-team-sol-review.md、requestid-timing-table.md、ruling-2026-11-04-q1.md、design.md（r1）、CONTRACTS §7 现全文、ADR-013、B3b §4。

## A. §7 修订文本 v2 全文（CN 权威）

**A1. replay store 段最小措辞调整**（在「……且完成方必须先发布 R 后发布 C。」之后追加）：
> （新增）「如某（平台 × 前提）经本节「C6 判定路径」完成协议先行修订并经独立评审后激活预建槽位实现，则该平台**仅 R/C 记录文件**的发布可见性以槽内头标志提交为准，其目录项在首次派发前由装配层预建一次性建立、此后永不重建或替换；未激活的平台与激活前的实现维持本段既有口径逐字不变。」

**A2. 新增段落「C6 判定路径」（bootstrap 预建固定槽，ADR-013）**（插入 CONTRACTS.md:232 段之后）：
> **C6 判定路径**：（C6 定义）C6 指每 requestId 的 R/C 记录文件作为槽文件在首次派发前预建：目录项一次性建立、此后永不重建或替换；发布 = 槽内写 payload + 数据 flush + 槽内头标志提交（payloadLen 0→len）+ 头 flush，可见性以槽内头标志界定而非目录项创建（出处 `docs/validation/t027-windows-publication-durability.md:28`、`docs/ADR_REGISTER.md:23`）。（预建生命周期，Q1=B，2026-11-04 用户裁定）**本裁定推翻 r1 的 High 6 判断**：r1「requestId 在 bootstrap 时刻未知」前提不成立（`runtime_agent.go:148→:181` 次序证明 store 构造时声明式 requestId 全集已确定），r1 的「首次发布前预建」解析作废。预建由生产装配层触发：`ExecuteAgentRunnerRun` 在 `NewAgentRunnerReplayStoreForRun` 返回之后、首次派发之前，以 `agentRunnerAssemblyRecords` 已校验的 requestId 全集调用新增可选方法 `PrebuildReplaySlots(requestIDs []string) error`；不得改 `NewAgentRunnerReplayStoreForRun` 签名，调用方必须传全集、不得分批。预建失败语义：任一槽预建失败 ⇒ 整个 run fail-closed（store 不可用、不进入首次派发）；部分成功不得视为可用；已建成的空槽允许保留（供重开幂等核验），但 run 必须停。重开语义：重开时对已存在槽只核验（magic + 长度 + payloadLen 合法性）不重建；分配记录（预建全集）已持久化而槽缺失 = corruption fail-closed；从未预建的槽缺失 = 正常。适用边界：预建只对生产装配路径要求；包内测试、A2 dispatch 测试等非生产入口不强制预建。（A7 偏离登记）本实现较 A7 原文的偏离点 = 预建由装配层触发而非 bootstrap 程序内部；该偏离不改 C6 的结构性优势（目录项仍在首次派发前一次性建立）。（预建全集持久化）replay root 内落一条 `prebuild-set.jsonl` store-local 记录；重开时声明集必须与该记录逐项一致（同为全集、无增无减），不一致 fail-closed（计划漂移按 corruption）；记录不存在 ⇒ 按首次预建执行（既有槽幂等核验、缺失槽补建，全部成功后发布记录）。（接缝范围）仅 R/C 记录采用槽容器：新增槽发布函数只由 R/C 发布点调用，`writeReplayRecordNoReplace` 主路径不动，ownership/launch-intent/launch-identity/terminal-intent/terminal-closure 等 store-local 记录保持原发布与原读取逐字节不变；显式禁止把其它记录迁入槽容器，若后续评估认为必须迁，停下升级用户裁决，不得自行扩域。（槽容器格式）槽文件 = 16B 固定头（8B magic `PRSLOT1\0` + uint64le payloadLen）+ payload（canonical JSONL，与既有 R/C 文件内容逐字节一致）；「一条 canonical JSONL」（`:230`）由 payload 维持，槽头属容器层不属于记录层。已发布 = payloadLen>0 且总长 = 16+payloadLen；空槽 = 总长 16 且 payloadLen=0。（未提交中间态冻结与校验顺序）「写 payload 后、头提交前」的读者可见结果恒冻结为 absent；校验顺序冻结：① 先查 magic（错 = corruption fail-closed）；② 再查总长是否 ∈ {16, 16+payloadLen}——总长 16 且 payloadLen=0 = absent；payloadLen=0 且总长>16 = 未提交中间态 ⇒ absent（不判 corruption）；总长<16，或 payloadLen>0 且总长≠16+payloadLen = corruption；③ 再按 payloadLen 有界读 payload 并严格解码 + requestId 校验，失败 = corruption。（三时点）头写可见点 = 头字段写入即可能对读者可见（不等待 flush #2，允许集据此冻结）；flush #2 耐久点 = 头 flush 成功（阶段 `slot-header-synced`）；解锁点 = flush #2 成功后释放槽锁与 store 锁。（仲裁协议）写者与读写同句柄（`CreateFileW(OPEN_EXISTING, GENERIC_READ|GENERIC_WRITE, FILE_SHARE_READ|FILE_SHARE_WRITE)`）；写者取 `LockFileEx` 排他锁（范围覆盖槽文件全域）、读者同句柄读；锁序固定 `store.mu → 槽文件锁`，禁止反向取锁、禁止多槽嵌套；取锁等待有界（超出按 uncertain 收敛）；解锁 = 同句柄精确范围解锁后关闭；映射视图不属保护范围（本设计不使用映射 I/O）。（失败重试）未提交中间态允许回收：锁内先截断至 16 再写新 payload（长载荷失败后重试短载荷不残留尾部）；已提交槽（payloadLen>0）绝不改写、绝不截断，按 `fs.ErrExist` 走既有碰撞收敛路径。（关系与引用纪律）C6 与 C1/C2/C3 的关系沿用 ADR-013 既有枚举与 B3b 结论账本（`docs/validation/t027-b3b-durability.md:44-48`、`:122-128`），本条款不改写其中任何结论；引用 C2 的候选级 `proven` 必须带完整括注「（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）」，不得截断为「C2 proven」。（前置条件）C6 从结构性候选升为可注入候选的前置条件：协议层完成仲裁原语替换——槽文件以 no-replace 一次性建立且永不删除、替换或追加，槽内头标志只允许 0→1 单向翻转并以槽级排他锁仲裁并发写者；该替换必须先修订本节并经独立评审后方可进入注入取证。（逐项取证）(a) 可复现注入或对照证据——在同一（平台 × 前提）装置上对槽预建点与各发布阶段注入，全部有效轮落在允许集内、证据包可由轮次记录与 journal 复算、无真孤儿；(b) 区分力门——同一装置、同一写入器与同一记录级载荷下，负对照（槽创建省略耐久原语或省略任一 flush）出现 ≥1 次丢失，正对照（完整 C6 序列）零丢失，门不成立即封顶 `unresolved`；(c) 机制证据——每轮记录机制令牌（槽创建原语、flush 调用位置与次数、头提交写），变异臂（编译期降级原语）下装置必须拒绝把「机制未执行的存活」兑现为证据。任一判据缺失即 `unresolved`，存在可归属反例即 `disproven`；三判据同时满足仅构成该（平台 × 前提）下 C6 的候选级结论，不自动改变平台发布耐久等级（三判据本体与解锁边界沿用既有句，原文不动）。（平台等级分层）C6 取证完成前，Windows 平台发布耐久等级保持 `unproven`，first-dispatch、launch-intent 与终局发布继续按 `unproven` 拒绝；取证完成后按本节既有三判据判定并受解锁边界约束——「落地 C6」本身不改变任何等级。（非 Windows 零影响）Unix 路径与既有 proven 实例不受本条款任何影响；C6 仅描述 Windows 分支的可注入判定路径；未激活构建的构造器与发布路径的逻辑路径与执行行为与现状等价（新增守卫行不改变未激活构建的行为）。

## B. High 1–3 逐条闭合说明（对照红队原文 file:line）

**High 1（design.md:20/22/24/55）**
- 空槽/损坏槽混淆（:24 vs :20）：闭合 = A2「未提交中间态冻结与校验顺序」——中间态（payloadLen=0 且总长>16）读作 absent 而非 corruption，但**必须先经长度分类**，不得仅凭 `payloadLen=0` 判 absent；异常长度（总长<16、payloadLen>0 且总长≠16+payloadLen）判 corruption。
- 丢槽后重建承诺无法执行（:22 + `agent_runner_replay_store.go:294` 的 R-only absent 反例）：闭合 = A2「预建全集持久化」——`prebuild-set.jsonl` 使「应存在而缺失」可机械检出；预建全集内缺失槽 = corruption fail-closed、**绝不重建**；删除 r1「不静默重建、由孤儿检测兜底」的表述，改为分配记录驱动的 fail-closed。
- 可见性与 flush 阶段混同（:55 + 页缓存反例）：闭合 = A2「三时点」——头写可见点（头写即可能可见）/flush #2 耐久点（`slot-header-synced`）/解锁点分开定义；「flush #2 前被 kill」注入格见 F；R 允许集自头写可见点起放宽为 `dispatched-unknown`（非恒 absent），三判据本体不动。
- 失败重试尾部处理：闭合 = A2「失败重试」——锁内 truncate→16 后重写；已提交槽永不改写。

**High 2（design.md:28/40，ownership 发布 `agent_runner_replay_root.go:425` / 解码 `:481`，CONTRACTS.md:230）**
- 闭合 = A2「接缝范围」+「槽容器格式」：r1 的「replaces the whole no-replace publication」作废；`writeReplayRecordNoReplace` 主路径不动，新槽发布函数只接 R/C 两个发布点；ownership/launch/terminal 原发布原读取逐字节不变，全局发布器与 R/C-only 读者的覆盖不匹配由此消除；「一条 canonical JSONL」由 payload 维持，槽头属容器层。

**High 3（design.md:17 / A7:28）**
- 闭合 = A2「预建生命周期」按 Q1=B 重写，覆盖七条边界，明写「本裁定推翻 r1 的 High 6 判断」，删除 r1「首次发布前预建」全部措辞。

**Medium 合并处置**：M4（design.md:49 核心声明）与 M5（区分力测量对象）闭合于 E/F（有条件声明 + 物理槽存活/未提交载荷存活/逻辑记录可见三测量分开）；M6（LockFileEx 协议不足）闭合于 A2「仲裁协议」（范围/同句柄/有界等待/锁序/解锁/映射视图排除）；M3（引用纪律）维持已修复状态，A2 新增引用处全部带完整括注。

## C. 二.1–7 → A2 句映射

| 条 | A2 对应句 |
|---|---|
| 二.1 | 「预建由生产装配层触发：`ExecuteAgentRunnerRun` 在 `NewAgentRunnerReplayStoreForRun` 返回之后、首次派发之前……不得改签名」 |
| 二.2 | 「以 `agentRunnerAssemblyRecords` 已校验的 requestId 全集……调用方必须传全集、不得分批」 |
| 二.3 | 「预建失败语义：任一槽预建失败 ⇒ 整个 run fail-closed……已建成的空槽允许保留……但 run 必须停」 |
| 二.4 | 「重开语义：重开时对已存在槽只核验……分配记录已持久化而槽缺失 = corruption fail-closed；从未预建的槽缺失 = 正常」 |
| 二.5 | 「适用边界：预建只对生产装配路径要求；包内测试、A2 dispatch 测试等非生产入口不强制预建」 |
| 二.6 | 「（A7 偏离登记）本实现较 A7 原文的偏离点 = 预建由装配层触发而非 bootstrap 程序内部；该偏离不改 C6 的结构性优势（目录项仍在首次派发前一次性建立）」 |
| 二.7 | 「（预建全集持久化）replay root 内落一条 `prebuild-set.jsonl`……不一致 fail-closed……全部成功后发布记录」+ 本设计 D 节 |

## D. 二.7 结论：**落**「预建全集」记录

- **理由**：二.4 的「分配记录存在而槽缺失 = corruption」必须依赖已持久化的分配记录才能机械检出；R-only 时刻 R 与 C 均缺失，R+C 绑定与孤儿检测无法兜底（红队 High 1 第二点反例）。顺序定为**先建全槽、后发记录**：部分失败时无记录 ⇒ 缺槽按「从未预建」= 正常，且 run 已 fail-closed；崩溃不会制造「记录存在而槽缺失」的假 corruption。
- **格式**：`<replayRoot>/prebuild-set.jsonl`，单条 canonical JSONL store-local 记录：`kind:"agent-runner-replay-prebuild-set"`、`schemaVersion`、`runId`、`createdAt`、`count`、`requestIds[]`（字典序、去重、逐项 `ValidID`）；recordHash 域 `proofrail:agent-runner-replay-prebuild-set:1\n`。
- **发布路径**：普通 no-replace 发布（同 ownership 原语，不进入槽容器、不要求槽头）；不进 wire、不参与 R/C 的 recordHash；碰撞后有界重读收敛，无法判定即 convergence failure；不要求自身 `proven`（与 ownership 同口径）。

## E. C6 函数级设计增量（相对 r1）

1. 保留：`agent_runner_replay_store_c6_windows.go`（`//go:build windows`，默认惰性）+ `_c6_enable.go`（`//go:build windows && c6slot`，init 置 `replayStoreSlotEnabled=true`）。
2. **新增方法**：`(*AgentRunnerReplayStore).PrebuildReplaySlots(requestIDs []string) error`（共享层定义，委托变量接缝 `replayStorePrebuildSlotsFn`，nil 时返回 nil）；c6 激活实现：校验全集 → 逐槽幂等预建（存在 = 核验 magic+长度+payloadLen；缺失且 prebuild-set 存在 = corruption；缺失且记录不存在 = `createReplaySlotWindows`）→ 全成后发布 `prebuild-set.jsonl`。未激活构建 = 无操作成功（现状等价）；Unix 恒无操作成功（Unix 路径零改动）。
3. **装配层调用点**：`runtime_agent.go` 步骤 3（`:181`）后、步骤 6 派发端口构造前加一行 `store.PrebuildReplaySlots(ids)`（ids 取自 `agentRunnerAssemblyRecords` 返回记录集的 requestId 全集）；错误 = 直接返回 run 失败（fail-closed）。
4. **发布接缝收窄**：删除 r1 的全局 `replayStoreSlotPublishHook` 替换方案；改为 R/C 发布点专用包裹 `publishReplayRCNoReplace`（仅 `RecordRequest`/`RecordCompletion` 调用；hook `replayStoreSlotPublishHook` nil 时原样调 `writeReplayRecordNoReplace`）。ownership/launch/terminal 继续直呼 `writeReplayRecordNoReplace`。
5. **写入路径**：`writeReplaySlotNoReplace` 六阶段 `slot-created`→`slot-written`→`slot-data-synced`→**`slot-header-written`（新）**→**`slot-header-synced`（原 slot-header-committed 拆分）**；加锁内 truncate→16 回收、flush #1 失败不写头、flush #2 失败头已写（重开按已发布收敛）。
6. **槽读者**：`loadRequestRecordSlot`/`loadCompletionRecordSlot` 按 A2 冻结校验顺序；重开核验走同一读者。
7. **机制令牌/变异臂**：保留私有缝 `slotFlushFn`/`slotCreateFn`；新增头写缝 `slotHeaderWriteFn`（变异臂：不写头 ⇒ 头可见点缺失可测）。
8. Unix 零影响口径维持 r1 已修正措辞：「nil 时逻辑路径与执行行为与现状等价；新增守卫行不影响 Unix 路径的字节级行为」。

## F. 三判据对应（B 下取证与注入）

- **(a) 注入**：决定格 = {`slot-created`（预建 MoveFileExW(WRITE_THROUGH)）, `slot-written`, `slot-data-synced`, `slot-header-written`, `slot-header-synced`}，B3 式 cut 逐格 kill；**「flush #2 前被 kill」窗口 = cut at `slot-header-written`**（头已写、flush #2 未执行），该轮盘点读者可见 ∈ {absent, dispatched-unknown} 均落允许集（R 允许集：头写前恒 absent、头写后恒 dispatched-unknown）；机制证据 = flush #2 未执行 ⇒ 该轮不得兑现耐久主张（B3b fail-closed 规则），仅证可见性/仲裁层。
- **(b) 区分力门**：负对照 = 槽创建省略 WRITE_THROUGH / 省略 flush #1 / 省略 flush #2（预期决定格 ≥1 丢失）；正对照 = 完整 C6 序列零丢失；门须分别在 `slot-created` 点与头提交点标定（B3b 只证 temp-written 一点）。测量对象分开（M5）：物理槽存活、未提交载荷存活、逻辑记录可见三测量各自报告，不以同一 absent 判槽丢失或存活。
- **(c) 机制证据**：令牌逐轮记录 `MoveFileExW(WRITE_THROUGH, no-replace)`、`LockFileEx`、`FlushFileBuffers`×2、头提交写（payloadLen 0→len）；变异臂降级（flush→no-op、WRITE_THROUGH→普通 move、头写→不写）下「存活但机制未执行」必须被装置拒裁。
- **升级路径**：三判据满足 ⇒ 仅 C6 候选级结论；平台级升级仍需协议先行（本片 §7 修订即先行文本）、证据可追溯、独立评审、AT-23 不豁免；C2 `proven`（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）不被改写或升级。

DESIGN: DONE

---

# design-r2 增补节（用户 2026-11-04 要求补 A/B 两处；① 第二轮增补返回文本逐字归档；不改 A2 主体）

（补充 A）prebuild-set.jsonl 自身的耐久等级与丢失降级语义（增补节，不改 A2 主体）

**A-1. 自身是否要求 `proven`：不要求。** `prebuild-set.jsonl` 是 store-local 元数据，与 `ownership.json` 既有口径逐项同构对齐：不进 wire、不参与任何 recordHash、不属于 R/C 已发布记录集，自身以普通 no-replace 单调发布（同 ownership 原语，不进槽容器、不要求槽头，D 节已定），且**不以发布耐久 `proven` 为前置（Windows 同样创建）**——与 CONTRACTS §7 ownership 句「且不以发布耐久 `proven` 为前置（Windows 同样创建）」同一口径。对齐理由两重：① **定位对齐**：§7 把 proven 前置只施加于 first-dispatch/launch/终局发布资格（R、C、launch-intent、terminal-intent 的发布权）；ownership 是所有权自愈器，prebuild-set 是计划漂移检测器，二者缺席都不改写 R/C 完整性判定，故无需、也不得给元数据单独立 proven 前置（立了反与 ownership 口径冲突）；② **事实对齐**：Unix 下普通 no-replace 发布附带父目录 fsync（平台事实，非记录级主张），prebuild-set 目录项在 Unix 上自然获得与 ownership 同等的平台级耐久；Windows 下与 ownership 同样不证明。**结论：prebuild-set 自身耐久不单独判定，随平台发布耐久口径走，不要求 proven。**

**A-2. 与二.4 是否冲突：不冲突。** 二.4 的 corruption 规则是条件句——触发前提是「分配记录（预建全集）**已持久化**」（记录在）且槽缺失；prebuild-set 自身丢失时该前提不成立，规则不触发。两分支按记录存在性互斥且穷尽：**记录在 + 槽缺 = corruption fail-closed**（二.4 原句）；**记录不在 + 槽缺 = 从未预建 = 正常**（A2「记录不存在 ⇒ 按首次预建执行」原句）。「记录在而槽缺」的唯一现实来源是记录写成之后槽目录项的外部丢失或电力差分丢失（含 unproven 平台的差分持久化窗口），该态**仍判 corruption**（保守方向，宁可误停不误建）——这是二.4 题中之义，不是与本补充的冲突；本补充只裁定反向情形（记录自身丢失）。因此二.4 无需改任何字句，只需把记录缺失分支的语义显式化。

**A-3. 结论二选一：安全降级**（可接受静默补建）。A2 应补写：「prebuild-set 缺失时，缺失槽按从未预建处理，此为有意降级、不构成 corruption。」降级为何安全（四条）：
① **与既有事实对齐**：与 R-only 时 R+C 均缺失、无法区分「从未创建」与「丢失」的既有口径一致——store 从不承诺对「记录自身缺失」区分两种来历，缺失即按缺失态处理；
② **补建是幂等空槽、不覆盖已发布数据**：补建 = no-replace 创建 16B 空槽（payloadLen=0），只作用于**缺失**的槽；已存在槽（含已发布 R/C 的槽）一律只核验不重建，故不可能把已发布数据降级或伪造终局；
③ **降级面封闭**：丢失的只是漂移检测器本身，不是 R/C 完整性——R/C 绑定、orphan、runID 绑定校验照常执行；且预建全集来自已校验的声明式 request records（`agentRunnerAssemblyRecords`，确定性），补建出的集合与丢失前逐项相同，漂移检测的丢失无实际安全后果；
④ **成本不对称**：fail-closed 选项要求给元数据立单独耐久保证（自身 fsync 证明前置 + 分级证据），与 ownership 口径冲突，且在 Windows（unproven 平台）上不可满足——除非为一条元数据单独建设耐久机制，把一次元数据电力丢失升级为需人工干预的硬 corruption，可用性代价远大于零安全收益（安全不变量「已发布 R/C 不被覆盖、缺槽不产生假终局、缺失槽可幂等补齐」在降级语义下恒成立）。故不采用 fail-closed。

**A-4. 逐字段格式与顺序论证（仍成立）。** 逐字段（单条 canonical JSONL，`<replayRoot>/prebuild-set.jsonl`）：

| 字段 | 取值/约束 | 解码违规处置 |
|---|---|---|
| `kind` | 固定 `"agent-runner-replay-prebuild-set"` | 其它取值 = corruption |
| `schemaVersion` | `1` | 未知版本 = corruption |
| `runId` | store 构造绑定 runID | 与当前 runID 不符 = root conflict |
| `createdAt` | RFC3339 UTC | 不可解析 = corruption |
| `count` | = `len(requestIds)` | 与数组长度不符 = corruption |
| `requestIds[]` | 字典序、去重、逐项 `evidence.ValidID` | 乱序/重复/非法 ID = corruption |
| recordHash 域 | `proofrail:agent-runner-replay-prebuild-set:1\n` + canonical body | 与既有 recordHash 口径一致 |

顺序论证（先建全槽、后发记录）在降级语义下**不仅成立、而且是唯一安全序**：部分失败 ⇒ 无记录 ⇒ 重开按「从未预建」走幂等补建 = 正常，且该 run 已按预建失败语义 fail-closed（无首派发、无 R/C 已发布）；若改为先发记录后建槽，则「记录写成、槽未建全」的崩溃窗口会在重开时命中二.4 触发前提（记录在 + 槽缺）而制造**假 corruption**——正序把该窗口整体折叠进安全的「记录不存在」分支。记录自身发布失败（碰撞不收敛等）同按预建失败语义 fail-closed；重开按记录存在与否分别收敛。

---

（补充 B）共享层新增方法与装配层新增调用对 Unix 路径的影响论证（增补节，不改 A2 主体）

**B-1. 共享层新增导出方法本身不改变 Unix 构建的可观察行为。** 「输入-输出等价 + 无副作用」判定句式：对任意 `requestIDs` 输入 A 与任意 store 状态 S，Unix 生产构建下 `store.PrebuildReplaySlots(A)` 满足——**输出 = `(nil)`，S 的迁移 = 恒等**（执行前后所有可观察状态：store 字段、replay root 目录树逐文件字节与条目集合、任何外部文件，均逐字节一致）。判定路径三段：
① 方法体在接缝为 nil 时仅执行「读 `replayStorePrebuildSlotsFn` → 判 nil → 返回 nil」，**不触碰任何既有数据/状态**——不解引用接收者字段、不读写文件、不调用既有函数，故无任何既有行为被共享代码路径改写；
② Unix 构建下 `replayStorePrebuildSlotsFn` **静态恒 nil 且无任何赋值路径**：共享层只声明不赋值，唯一赋值点位于 `//go:build windows && c6slot` 构建文件（`_c6_enable.go` 及其 c6 激活实现），Unix 构建不参与编译 ⇒ 运行时该变量不可能非 nil；
③ 返回 nil 无错误 ⇒ **不新增错误分支**：调用方「err 非 nil ⇒ fail-closed」的既有错误返回面在 Unix 上不可达，调用与不调用在可观察结果上等价。
因此这是「**新增方法**」（新增一个可被调用的导出符号），不是「新增守卫行」（未在既有代码路径上插入条件判断）；区分点正是 ②——守卫行会参与每个既有调用的判定，本方法只在自己被显式调用时执行一次 nil 判定。
可作 ② 测试的机械断言：
- BA1：Unix 构建下 `store.PrebuildReplaySlots(任意 requestIDs 集，含空集/非法 ID)` 恒返回 `nil`（接缝分支不校验参数、不触碰 receiver，`(*AgentRunnerReplayStore)(nil)` 调用亦返回 `nil`）；
- BA2：调用前后对 replay root 递归遍历（排序后逐文件 sha256 + 目录条目集合比对），**逐字节不变**；
- BA3：`State`/`RecordRequest`/`RecordCompletion` 在调用前后对任意 requestId 的返回完全一致（等值重放）；
- BA4：`GOOS=linux go build ./...`、`GOOS=linux go vet ./...` 通过（共享层方法签名在 Unix 上被编译，是防签名漂移的编译级护栏）；
- BA5：机械核对 `replayStorePrebuildSlotsFn` 的赋值仅出现在 `windows && c6slot` 构建文件中。

**B-2. 装配层新增一行对 Unix 路径是「多一次返回 nil 的调用」，且可观察行为等价于完全无操作。** 行为等价论证（四条）：
① **无前置条件失败路径**：nil 接缝分支在 store 有效与无效（含 nil receiver）时恒返回 nil（B-1 ②③），调用点在 Unix 上不存在可达的失败返回——「err 非 nil ⇒ fail-closed」分支不可达；
② **不改变后续步骤 4–7 的任何输入**：步骤 4–7 的输入（snapshotRoot、options、launcher 配置、plan/records、store、clock）均非该方法可写对象，nil 分支零副作用 ⇒ 步骤 4–7 收到的输入值与未调用时同一；
③ **不改变错误返回面**：恒 nil ⇒ 不新增错误分支、不改变任一既有分支取值 ⇒ `ExecuteAgentRunnerRun` 在 Unix 上的输入-输出关系逐点不变（同一 options ⇒ 同一 RunSummary/err）；
④ **调用点前置位置合法**：步骤 3 刚返回的 store 已通过全部构造校验，调用点的唯一依赖（store 存在）已满足，不引入新的时序或竞态条件。
若 ② 无法接受「多一次调用」，替代形态 = 构建标签隔离调用点（抽 `prebuildReplaySlotsIfActive(store, ids) error`：c6slot 激活构建提供真实现，非激活构建含 Unix 提供恒返回 nil 的空实现，调用点不变）。取舍：成本 = 调用点拆为两个平台文件、nil 语义在接缝与空实现两处重复（需同步维护）；收益 = 指令层零调用，但可观察行为与共享调用点完全相同（同 B-1 判定的零差异）；风险 = 共享调用点是签名漂移防护的最后一环，若调用点移入 c6slot 专属文件，Unix 构建将不再引用该方法，共享层签名可能只在 c6slot CI 中被编译校验，漂移防护减弱。**默认维持「调用点共享、行为等价」**；替代形态仅在 ② 硬性不接受时启用，并须在报告中如实登记该取舍。

**B-3. 一句话可审结论：「Unix 路径可观察行为与现状等价；新增方法仅在 c6slot 激活构建下产生任何可观察差异。」**