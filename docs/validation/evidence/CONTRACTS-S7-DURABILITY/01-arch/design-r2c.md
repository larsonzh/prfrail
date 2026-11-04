# ① 架构微修 design-r2c 原文（P54 CONTRACTS-S7-DURABILITY）

载体：`deep-reasoner` + `DeepSeek V4 Pro`（只读，2026-11-04）。以下为返回文本逐字归档。

# P54 `CONTRACTS-S7-DURABILITY` · ① 架构微修 design-r2c（只读返回文本）

载体：`deep-reasoner` + `DeepSeek V4 Pro`；2026-11-04；只读，未写任何文件。已逐字读取：design-r2b（A2 v3 适用边界句、允许集表、逐条闭合表）、design-r2 E 节激活设计、`agent_runner_replay_store.go`（`replayStoreSlotEnabled` 声明 :152、读者守卫 :695/:721、`publishReplayRCNoReplace` :752、`PrebuildReplaySlots` 接缝 :329）、`_c6_enable.go`（init 置位 :73-74、`prebuildReplaySlots` 全流程）、`c6_windows.go`（init 接线 :96-103）、`c6_windows_test.go`（26 个 TestC6* + 辅助函数）、03-tests 三份清单（93 失败分布 / 26 绿 / 红线编译）。

## 1. 裁定：选 ②（store 级激活），排除 ① 与 ③

A2 v3 适用边界句逐字：「预建只对生产装配路径要求；包内测试、A2 dispatch 测试等非生产入口不强制预建。」

- 「只对…要求」：构建级 flag 把预建前置条件放大为**整个构建**的要求——凡经 `publishReplayRCNoReplace` 与 R/C 读者的路径都被强制进槽容器，要求面不再限于生产装配路径。
- 「不强制预建」：该句语义必须 = 非生产入口可不预建**且仍能正常工作**（否则成空话）。构建级激活下不预建 = 发布 file-not-found（93 失败主形态）、读取 corruption/missing（`TestAgentRunnerReplayStoreRequestCollisionVisibilityExhaustionIsConvergenceFailure` 的 magic mismatch 即此）——非生产入口被事实强制为「要么预建、要么失效」，与句子字面直接冲突。
- 排除 ①（验收收窄 `-run 'TestC6'`）：只把冲突移出验收视野，不消除契约违背；`go test -tags c6slot ./...` 恒红，任何带 tag 的 CI/回归腿失效，未来新增非生产入口（fuzz/example/benchmark）继承同一强制。适用边界句是契约条款而非测试范围注记，不能靠「不验收」兑现。
- 排除 ③（构建级 flag 下未预建回落普通路径）：回落会向已激活 store 静默写入普通 JSONL（未声明 requestId 或半预建态），而读者仍先走槽读者 ⇒ magic mismatch corruption（即 93 失败中那条的形态）；且瓦解 A2「声明缺失 = corruption fail-closed」与「先建全槽、后发记录」顺序论证的防线（静默回落 = 未预建也可发布）。回落若叠加 store 级状态，等价于 ② 去掉 fail-closed 严格性——严格性更弱，不取。

## 2. store 级激活精确设计

1. **激活位位置**：`AgentRunnerReplayStore` 新增字段 `slotsActive bool`（`store.mu` 保护；仅内存、不落盘、运行期不清零）。现有包级 `replayStoreSlotEnabled` 语义**收窄**为「构建能力位」（实现存在与否；仅 `_c6_enable.go` 置位；不再承担调度判定）。
2. **层叠关系（两者都成立才走槽路径）**：是。槽路径 = `replayStoreSlotEnabled`（编译期能力）∧ `store.slotsActive`（运行期实例态）∧ 对应 hook 非 nil。生产 Windows 构建不带 c6slot ⇒ 能力位 false ⇒ 恒普通路径，零语义不变（`PrebuildReplaySlots` 接缝 nil 返回 nil、激活位永不置位且无消费者）；Unix 同。
3. **置位时机**：仅在 `prebuildReplaySlots`（已持 `store.mu`）成功返回前置位——首建分支（`firstReplayPrebuild` 全槽就位 + 集合记录发布后）与重开核验分支（`verifyReplayPrebuildSetSlots` 通过后）都算成功。**不按 prebuild-set 存在从盘上恢复**：记录丢失是已接受降级路径，若激活依赖记录存在，该路径下普通读者会把槽文件读为 corruption；激活必须每次运行由同一调用重建（装配层每次都调用，首建/重开同一入口）。
4. **未激活 store 回落路径**：发布 = `store.publishReplayRCNoReplace` 未激活 ⇒ 原样调 `writeReplayRecordNoReplace`；读取 = 未激活 ⇒ 既有变量 `replayStoreLoadRequestRecord`/`replayStoreLoadCompletionRecord`（默认 = 纯普通读者）。实现形态：调度判定从自由函数上移到 store 方法——`publishReplayRCNoReplace` 改为 store 方法（仅两处调用 :493/:562，均持 `store.mu`）；自由函数 `loadRequestRecord`/`loadCompletionRecord` 删除槽守卫回归纯普通读者；六处读取点（`loadRequestLocked`、`loadCompletionLocked`、`ensureRequestLocked` 碰撞重读、`RecordCompletion` 碰撞重读、`launches.go:53`、`terminals.go:52`，已逐一核对全部为持 `store.mu` 的 store 方法）改走新增 store 级读入口（能力位 ∧ 激活位 ∧ hook）。`replayStoreLoadRequestRecord` 变量替换测试（store_test.go:350/:379）语义恢复前 C6 口径。
5. **测试三态**：默认构建（无 tag）= 能力位 false、hooks nil、激活位永不置位 ⇒ 与现状在限定观察面上等价（沿用 r2b B-1 口径，不含指令级保证）；linux = 三个 c6 文件均带 `windows` 约束、编译集不变，共享层仅 +1 字段与调度行，nil 缝下等价；c6slot = 93 项既有测试未调用预建 ⇒ 未激活 ⇒ 普通路径全绿，26 项 TestC6 经 `c6PrebuiltStore`/显式 `PrebuildReplaySlots` 激活 ⇒ 槽路径照常。③ 验收口径升格为 `go test -tags c6slot -count=1 ./internal/adapters/` 全绿，不再需要 `-run 'TestC6'` 收窄。

## 3. 半预建态语义与顺序相容性

预建中途失败 ⇒ 部分空槽已建、激活位保持 false、run fail-closed（装配层错误即停）。若调用方无视错误继续（生产不会）：该 store 走普通路径；普通读取碰到残留槽文件 ⇒ 解码失败 ⇒ corruption fail-closed，**不会静默错读**。重开：集合记录不存在 ⇒ `firstReplayPrebuild` 幂等核验已有空槽 + 补建缺失 + 补发记录 ⇒ 成功后激活——与 A2「已建成的空槽允许保留（供重开幂等核验）」一致。与「先建全槽、后发记录」序相容：激活位只在「全槽就位 +（首建）集合记录发布」后置位，且首条 R/C 发布只能发生在 `PrebuildReplaySlots` 成功返回之后 ⇒ 不存在「记录已发、激活位未置」或「反序崩溃窗口制造假 corruption」的新窗口；重开时已发布槽只核验、缺失槽补建的分支语义不变。

## 4. A2 是否需要澄清句：需要，补一句

理由：A2 v3 适用边界句只规定「谁必须预建」，未规定「未预建者在激活构建内的行为」——正是本次冲突的语义空白。在 A2 v3 适用边界句后追加（逐字）：

> （激活语义）槽路径仅在「构建携带槽实现」且「该 store 本次 `PrebuildReplaySlots` 成功完成（首次预建或重开核验）」两者同时成立时对该 store 启用；构建携带实现但未预建成功的 store（含未调用预建的非生产入口）维持普通记录路径逐字不变；预建失败留下的部分空槽不改变该判定，重开仍按「从未预建」幂等核验补建。

配套：design-r2 E.1 句同步修订为「`replayStoreSlotEnabled` 收窄为构建能力位；新增 store 级 `slotsActive`，仅 `prebuildReplaySlots` 成功前置位；槽路径 = 能力位 ∧ 激活位 ∧ hook」。A1「未激活的平台与激活前的实现维持本段既有口径逐字不变」不需改——store 级激活是「激活」粒度的细化，未激活 store 即「激活前」。

## 5. 对 ③ 八项覆盖的影响

1. **槽原语形状**（create 16B / 总长 / 二写拒绝 / 空 payload 拒绝）——无影响：直呼原语，不经调度守卫。
2. **校验顺序六项**——无影响：直呼 `loadRequestRecordSlot`/`loadCompletionRecordSlot`。
3. **写入失败与截断重试**（flush 失败头未提交、长写失败短重试无尾）——无影响：私有缝驱动。
4. **预建生命周期**（全建幂等、任一失败 fail-closed）——语义不变；建议各增一条激活位断言（成功 = 置位、失败 = 未置位），测试函数本身不改。
5. **prebuild-set 记录与缺失语义**（字段精确 / 声明缺失 corruption / 记录缺失补建 / 记录丢失降级）——无影响：全在预建实现内部，激活位在其上层。
6. **允许集表**（进程 kill / 断电 / corruption fail-closed）——无影响：`c6PrebuiltStore` 已激活，seam panic 后同 store 内存激活位仍在，读者仍走槽路径。
7. **BA3 两臂**——**需改**。现对照臂依赖「构建级激活 + 预建树快照」才能与调用臂同走槽路径；store 级激活下对照臂未激活 ⇒ 普通路径遇槽文件 corruption，两臂失配。修订（r2b 补充 B 的 BA3 句同步改）：初始状态 = 两个等价未预建 store（重建）；对照臂直接执行 S（普通路径）、调用臂先 `PrebuildReplaySlots(A)` 再执行 S（槽路径）；比较口径 = 返回序列逐项一致 + R/C 记录内容逐字节一致（对照臂 = 普通 JSONL 全字节，调用臂 = 槽 payload 去 16B 头后的字节）；`prebuild-set.jsonl` 与槽头属激活机制既有证据，不参与逐字节比较。测试重写 = `c6RunSequence` 两臂跑完后按 requestId 做 payload 投影比较，替代 `c6TreeManifest` 全树字节比较。
8. **门禁不变式**（激活不改 durability 等级、unproven 拒绝首派发、重放零写）——无影响：`PublicationDurability` 与激活位解耦（unproven 判定先于/独立于发布与读取）；`TestC6GateActivationKeepsDurabilityUnproven` 对 `replayStoreSlotEnabled` 的断言保留（能力位在 c6slot 构建恒真）。

另：93 项既有失败全部由 ② 消除；可选新增一条「c6slot 构建下未激活 store ≡ 普通路径」钉针（BA3 对照臂已实质覆盖）。

## 6. 结论

`DESIGN: DONE`

（风险登记：本轮只读，无 diff、无 build/test/lint 日志。本裁定不触及三判据本体、Unix 路径、ADR-013、`NewAgentRunnerReplayStoreForRun`/`PrebuildReplaySlots` 签名、平台等级 unproven；BA3 修订属 design-r2b 补充 B 设计文本 + ③ 测试重写范围，随本裁定交 ②/③ 落实；A2 新增激活语义句进 CONTRACTS 修订文本 v3。）
