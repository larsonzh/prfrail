# ① 架构设计原文（P54 CONTRACTS-S7-DURABILITY）

载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读，2026-11-04）。以下为 ① 返回文本逐字归档。

**批次 1/1（单批完成，无分批）**。以下为 P54 `CONTRACTS-S7-DURABILITY` ① 架构环节只读设计（未写任何文件）。

## A. §7 修订文本（CN 权威）

**A1. replay store 段最小措辞调整**（`docs/CONTRACTS.md:230` 段，「每个文件必须以原子、无替换发布；协议单调，不得删除、替换或追加已发布文件，且完成方必须先发布 R 后发布 C。」之后追加）：
> （新增）「如某（平台 × 前提）经本节「C6 判定路径」完成协议先行修订并经独立评审后激活预建槽位实现，则该平台记录文件的发布可见性以槽内头标志提交为准，其目录项在槽创建时一次性建立、此后永不重建或替换；未激活的平台与激活前的实现维持本段既有口径逐字不变。」

**A2. 新增段落「C6 判定路径」（bootstrap 预建固定槽，ADR-013）**（插入 `docs/CONTRACTS.md:232` 段之后）：
> **C6 判定路径**：（新增）C6 指每 requestId 记录文件作为槽文件在首次发布前由 bootstrap 类程序预建（目录项一次性建立），发布 = 写槽内数据 + 文件级 flush + 槽内头标志提交，可见性以槽内长度/标志界定而非目录项创建（出处 `docs/validation/t027-windows-publication-durability.md:28`、`docs/ADR_REGISTER.md:23`；冻结契约下不可即插的判定 `docs/validation/t027-windows-publication-durability.md:80`）。（新增）C6 与 C1/C2/C3 的关系沿用 ADR-013 既有枚举与 B3b 结论账本（`docs/validation/t027-b3b-durability.md:44-48`、`:122-128`），本条款不改写其中任何结论；引用 C2 的候选级 `proven` 必须带完整括注「仅限前提 P′；P′ 覆盖物理机目标部署环境未论证」。（新增）C6 从结构性候选升为可注入候选的前置条件：协议层完成仲裁原语替换——槽文件以 no-replace 一次性建立且永不删除、替换或追加，槽内头标志只允许 0→1 单向翻转并以槽级排他锁仲裁并发写者；该替换必须先修订本节并经独立评审后方可进入注入取证。（新增）C6 逐项取证：(a) 可复现注入或对照证据——在同一（平台 × 前提）装置上对槽创建点与各发布阶段（槽写入/数据 flush/头提交）注入，全部有效轮落在允许集内、证据包可由轮次记录与 journal 复算、无真孤儿；(b) 区分力门——同一装置、同一写入器与同一记录级载荷下，负对照（槽创建省略耐久原语或省略数据 flush）出现 ≥1 次丢失，正对照（耐久槽创建 + 文件级 flush 完整执行）零丢失，门不成立即封顶 `unresolved`；(c) 机制证据——每轮记录机制令牌（槽创建原语、flush 调用位置与次数、头提交写），变异臂（编译期降级原语）下装置必须拒绝把「机制未执行的存活」兑现为证据。任一判据缺失即 `unresolved`，存在可归属反例即 `disproven`；三判据同时满足仅构成该（平台 × 前提）下 C6 的候选级结论，不自动改变平台发布耐久等级（三判据本体与解锁边界沿用既有句，原文不动）。（新增）平台等级分层：C6 取证完成前，Windows 平台发布耐久等级保持 `unproven`，first-dispatch、launch-intent 与终局发布继续按 `unproven` 拒绝（沿用既有句：「没有任何前提满足全部判据时平台等级为 `unproven`，运行期 `unresolved` 与 `disproven` 一律按 `unproven` 行为执行」）；取证完成后按本节既有三判据判定并受解锁边界约束（沿用既有句：升级仅在该平台、该前提覆盖的部署环境内解除对应拒绝，且必须协议先行修订、证据可追溯并通过独立评审）——「落地 C6」本身不改变任何等级。（新增）非 Windows 零影响：Unix 路径与既有 proven 实例「本地 Unix 文件系统上，只有新 link 创建成功后其父目录 `fsync` 成功，才可把该次发布记为 `proven`」（沿用既有句）不受本条款任何影响；C6 条款仅描述 Windows 分支的可注入判定路径，C6 生产实现未激活前生产构造器与发布路径的逻辑路径与执行行为与现状等价（新增守卫行不改变未激活构建的行为）。

## B. C6 落地设计（函数级，供 ②）

**架构歧义先行报备**：requestId 在 bootstrap 时刻未知，「bootstrap 预建固定槽」字面不可行；本设计解析为「首次发布前预建 + store 打开时核验既有槽」（见结论行报 Sol 事项）。

- **文件/标签**：新 `agent_runner_replay_store_c6_windows.go`（`//go:build windows`，全部 C6 实现，默认惰性）；新 `agent_runner_replay_store_c6_enable.go`（`//go:build windows && c6slot`，仅 init 置 `replayStoreSlotEnabled=true`）。生产 Windows 构建不带 `c6slot` ⇒ 零语义。
- **槽命名与布局**：路径与现有记录文件**同路径** `requests/request.<id>.jsonl`、`completions/completion.<id>.jsonl`（`pathForRequestID`，`agent_runner_replay_store.go:607-629`）⇒ 证据工件路径不变。文件 = 16B 固定头（8B magic `PRSLOT1\0` + `uint64le payloadLen`）+ payload（canonical JSONL）。空槽 = 长 16 且 payloadLen=0；已发布 = payloadLen>0 且长 = 16+payloadLen。
- **槽创建（C6 唯一涉及目录项的点）**：`createReplaySlotWindows(path)`——空头写同目录 temp → `FlushFileBuffers(temp)` → `MoveFileExW(MOVEFILE_WRITE_THROUGH)` no-replace 落到最终名（目标已存在 → `fs.ErrExist`，C2 契约同款；`t027-b3b-durability.md:45`）。每 requestId 一次。
- **写入路径**：`writeReplaySlotNoReplace(path, record)`：① `createReplaySlotWindows`（存在则复用），stage=`slot-created`；② `CreateFileW(OPEN_EXISTING, GENERIC_READ|GENERIC_WRITE, FILE_SHARE_READ|FILE_SHARE_WRITE)` + `LockFileEx` 排他锁，stage=`slot-opened`；③ 读头 payloadLen>0 ⇒ 解锁返回 `fs.ErrExist`（触发既有碰撞循环，`ensureRequestLocked` 零改动，`agent_runner_replay_store.go:496`）；④ 偏移 16 单次写 payload，stage=`slot-written`；⑤ `FlushFileBuffers` #1（数据），stage=`slot-data-synced`；⑥ 写头（payloadLen）+ `FlushFileBuffers` #2（头），stage=`slot-header-committed`，解锁返回。in-process 由 `store.mu` 排他，跨进程由 `LockFileEx` 排他。
- **五阶段对应**：C6 分支**替换**五阶段：无 temp、无 link、无父目录步；可见性自 `slot-header-committed` 起。旧阶段常量与两 B3b 接缝（`agent_runner_replay_store.go:87-133`）在 C6 分支**不触发**，二者仍 nil 默认。新阶段名放 c6 文件，复用既有 `replayStorePublishStage` 接缝（`:67-85`）。
- **异常逐项**：① 写槽中途崩溃（步骤 ④⑤ 间）→ payloadLen=0 ⇒ 读者见 absent，前一发布者未返回 first-dispatch ⇒ 无双派发（同 C1 temp 阶段允许集）；② flush #1 失败 → 不写头、返回错误、槽保持空，fail-closed；③ flush #2 失败 → 返回错误、槽保留已写数据，重新打开时头已提交则读者按已发布收敛（沿用既有句「可见性与耐久是两个事实」）；④ 外部篡改（长度≠16+payloadLen、magic 错、payload 不可解码）→ `replayStoreCorruption` fail-closed；⑤ 已发布槽目录项丢失 → 不静默重建（不可与「从未创建」区分），由 C 的 R+C 绑定/孤儿检测兜底。
- **接口与门禁不变**：`PublicationDurability()` 与 `RecordRequest`/`RecordCompletion` 的 unproven 拒绝（`agent_runner_replay_store.go:269-275`、`:314-316`、`:407`）逐字节不变；生产 Windows 仍 `unproven`（`agent_runner_replay_store_windows.go:5-7`）。实验写入器沿用 B3b D3 模式（直呼 `writeReplayRecordNoReplace`）。
- **新增接缝（共享层仅此一处）**，契约文本：
  ```
  // replayStoreSlotPublishHook replaces the whole no-replace publication with the
  // slot publication on Windows C6 experiment builds. It is nil in production and
  // on every non-Windows build; while nil, the logic path and observable behaviour
  // of writeReplayRecordNoReplace stay equivalent to today's, and the added guard
  // line does not affect the Unix path.
  // Seam contract: fires once per publication, before any stage fires. A
  // replacement must keep the header-commit arbitration (an already-committed slot
  // must surface as fs.ErrExist), must leave an uncommitted slot invisible, and
  // must leave a committed slot holding exactly the canonical payload. It may not
  // carry decisions and may not call back into the store.
  var replayStoreSlotPublishHook func(path string, record any) error
  ```
  在 `writeReplayRecordNoReplace`（`:671`）函数头加一行守卫。读者侧复用既有变量接缝 `replayStoreLoadRequestRecord`/`replayStoreLoadCompletionRecord`（`:45-46`），激活时改指槽读者（读头→空槽判 absent→否则有界读取+解码+requestId 校验）。机制令牌/变异臂用 c6 文件内私有缝 `slotFlushFn`/`slotCreateFn`。

## C. 非 Windows 零影响证明
- 全部 C6 代码在 `//go:build windows` 文件；`agent_runner_replay_store_unix.go` 与 Unix 平台函数零改动。
- 共享层仅 +1 变量 +1 行守卫。Unix 构建该变量恒 nil ⇒ nil 时 `writeReplayRecordNoReplace` 的逻辑路径与执行行为与现状等价；新增守卫行不影响 Unix 路径的字节级行为；读变量在 Unix 从不改指；旧阶段常量、两 B3b 接缝、门禁全未动。
- Unix 既有 proven 实例（`agent_runner_replay_store_unix.go:5-14`）不受新条款约束：A2 明示「Unix 路径不受本条款任何影响」，A1 以「未激活平台维持既有口径」为条件。验证：本片后 `go test ./...`（Linux 腿）结果不变（②/④ 执行）。

## D. 核心声明红队
「预建槽 + 文件级 flush 把目录项耐久转为文件数据耐久」——**条件成立，无不可闭合缺口**：
1. **成立点**：发布时点不产生任何目录变更 ⇒ 发布时耐久事实只剩文件数据 ⇒ `FlushFileBuffers` 的文件数据保证可覆盖；头提交独立 flush #2 使可见性标志与数据同批落盘。
2. **前提**：目录项耐久问题**未消失，而是被迁移**到槽创建点。声明只有在「槽创建过程的目录项耐久被独立取证」时才成立；若槽创建无耐久原语，断电丢项 = 已发布记录丢失（B3b `c1:record-linked` round 21 同型反例，`t027-b3b-durability.md:79`、`:144`）。A7 原文「有标准保证」是无前提的过度表述，A2 已把槽创建耐久取证写为 C6 前提。
3. **残余风险**：① 已发布槽目录项后续断电丢失（U4 类）——与 C2 发布后丢项同风险类，缓解 = 槽创建用 `MoveFileEx(WRITE_THROUGH)` + 独立注入取证 + 绝不静默重建 + 孤儿检测 fail-closed；② 头撕裂——16B 小写、窗口在 flush #1/#2 间：头未提交 ⇒ absent（安全方向），长度字段先落而 payload 未 flush ⇒ 解码失败按 corruption fail-closed（不产生假 absent）；③ U3 设备缓存——`FlushFileBuffers` 保证止于 OS 控制的存储栈，前提须携带存储栈声明（B3b §4.1 限定 (a) 同款）；④ 槽长度截断——长度校验 fail-closed。
4. 未发现「不写 proven 就无法落地 C6」的情形。

## E. 三判据对应（供 ③ 红测设计）
- **(a)**：B3 式注入，C6 决定格 = {slot-created, slot-data-synced, slot-header-committed}；R 允许集：头提交前恒 absent、后恒 dispatched-unknown；C 同构；U4 配对轮查孤儿。证据包由 rounds/+journal.jsonl 复算。
- **(b)**：负对照 = 槽创建省略 WRITE_THROUGH/省略 flush（预期 slot-created 格 ≥1 丢失）；正对照 = 完整 C6 序列零丢失。注意 B3b 门只证明 temp-written 一点区分力，C6 门须在 slot-created 点重新标定。
- **(c)**：令牌逐轮记录 `MoveFileExW(WRITE_THROUGH)`、`LockFileEx`、`FlushFileBuffers`×2、头提交写；变异臂（flush 降级 no-op）下「存活但机制未执行」必须被装置拒裁（B3b fail-closed 规则，`t027-b3b-durability.md:56-59`）。
- **升级路径**：三判据满足 ⇒ 仅 C6 候选级结论；平台级升级仍需协议先行（本片 §7 修订即先行文本）、证据可追溯、独立评审、AT-23 不豁免。C2 `proven`（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）不被本片改写或升级。

**报 Sol 事项**：① A7「bootstrap 预建」与 requestId 未知的矛盾（本设计取「首次发布前预建」解析）；② 共享层 +1 接缝是否可接受；③ `LockFileEx` 作为跨进程仲裁原语是否准入。
DESIGN: DONE（无升级事项；上述三项交 Sol 红队裁定，不影响 §7 文本主体）。
