# ⑤ 预审 r1 原文（P54 CONTRACTS-S7-DURABILITY）

载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`（只读，2026-11-04）。以下为返回文本逐字归档。

# ⑤ 预审报告 — P54 `CONTRACTS-S7-DURABILITY`（只读，未写任何文件）

输入已实际读取：红队两轮原文、design r1/r2/r2b/r2c、`contracts-s7-durability.md` §5、磁盘实况 diff（`git diff --cached`，暂存区即磁盘最新，无未暂存改动）、CONTRACTS{,_EN} §7 现状、③ 测试全文（1264 行）、ADR-029 补充节、MANIFEST/SHA256SUMS。

## 一、逐条发现

**High：无。Medium：无。** 以下 4 条 Low。

**L1｜附件 C 悬空引用｜Low｜`docs/CONTRACTS.md:234` 与 `docs/CONTRACTS_EN.md:228`**
逐字引用：「两场景逐切点允许集见本修订附件 C，验收器必须按场景选表、不得混用」／EN「the per-cut-point allowed set for both scenarios is given in Annex C of this revision」。但 CONTRACTS{,_EN} 全文无附件 C（仅 `evidence/01-arch/design-r2b.md` C 节有表；`contracts-s7-durability.md:101` 与 design-r2b:15/:66 另有提及）。契约读者无法在契约内解析该表。
最小修复：把允许集表作为 §7 附件 C 嵌入 CONTRACTS{,_EN}；或改引用为「本修订附件 C（`docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2b.md` C 节）」。语义无缺口（三时点句已逐场景复述），仅可解析性。

**L2｜BA1/BA2 仍为未执行断言｜Low｜`internal/adapters/agent_runner_replay_store_c6_windows_test.go`（整体缺失，无行号）**
design-r2 补充 B 的 BA1（nil receiver + 任意输入恒 nil）与 BA2（调用前后 root 目录树逐字节不变）至今无对应测试——与红队 r2 §4 原判「BA1、BA2、BA4、BA5：……目前是待执行断言，不是已通过证据」一致（BA4=linux build/vet、BA5=赋值审计已由本环机械核对覆盖）。默认构建对 no-op 接缝无直接钉针（仅经 `internal/console` 装配测试间接触及 `runtime_agent.go:193`）。
最小修复：默认构建加一条 `TestC6PrebuildNoopWithoutActivation`（`(*AgentRunnerReplayStore)(nil).PrebuildReplaySlots(...)` 返回 nil + 普通 store 调用前后目录树逐文件 sha256 不变）。

**L3｜计划漂移/runId 分支无直接钉针｜Low｜`internal/adapters/agent_runner_replay_store_c6_enable.go:281-285`**
逐字引用：`:282`「prebuild set is bound to runId %s but the store requested %s」、`:285`「prebuild set does not match the declared requestId set (plan drift)」。High-1 闭合所依赖的**正向漂移检测**（记录在 + 声明集不同 → corruption）与 runId 不符 → root conflict 两分支均无测试；现有 `TestC6PrebuildDeclaredSlotMissingIsCorruption`（缺槽）与 `TestC6PrebuildLostSetRecordDegradesDriftDetectionOnly`（丢记录）均不覆盖。分支实现正确（`slices.Equal` + 既有 corruption 机制），属回归防护缺口。
最小修复：③ 补两条：预建后以不同声明集重开 → `ErrAgentRunnerReplayStoreCorruption` 含「plan drift」；改 runID 重开 → `ErrAgentRunnerReplayRootConflict`。

**L4｜A2 尾句措辞未带「限定观察面」限定｜Low｜`docs/CONTRACTS.md:234` 尾（EN `:228` 尾同）**
逐字引用：「未激活构建的构造器与发布路径的逻辑路径与执行行为与现状等价（新增守卫行不改变未激活构建的行为）」。r2b 补充 B-1/B-2 已裁定所有等价主张「限于限定观察面，不含任何指令级保证」，本句未携带该限定；且 r2c 后机制为 store 级调度（读写入口三条件判定），「新增守卫行」框架已陈旧。行为主张本身成立（未激活 = 普通路径），仅措辞纪律问题；该短语溯源为用户 2026-11-04 措辞修正后经两轮红队未点名的文本。
最小修复：括注改为「（新增守卫行在限定观察面上不改变未激活构建的行为）」；若按「判定表注记原文不动」先例，改写需用户授权，否则维持并交 ⑧ 注记。

## 二、按序复核结论

1. **红队闭合独立复核（不依赖自检表）**：对照两轮原文逐条核 r2b/r2c——二轮 High（历史全集绑定）**是降级不是修复**，已如实登记：CONTRACTS A2 含「接受该降级。本降级经用户 2026-11-04 裁定接受。」（CN:234 / EN:228 均有「accepted by the user ruling of 2026-11-04」），ADR-029 补充节 ② 登记红队结果（1H+2M+4L，计数口径与红队原文一致）与接受降级、③ 记 ⑦ 兜底路径；二轮 High 2/High 3 闭合（接缝收窄与 Q1=B 均已实现级验证：`writeReplayRecordNoReplace` 直呼点仅 ownership/launch×2/terminal×2/prebuild-set，R/C 仅经 `publishReplayRCNoReplace` 两调用点）；一轮 3 High + 4 Medium + 事项①②③ 全部在 r2/r2b 有对应句且无旧句残留（「均落允许集」「恒 dispatched-unknown」「不参与任何 recordHash」均无残余，CONTRACTS:240 处为既有无改动的 ownership 原句）。
2. **契约 ↔ 实现**：激活语义句 ↔ `store.slotsActive`（store.go:214-220；仅 `_c6_enable.go:111/:117` 成功前置位，能力位∧激活位∧hook 三条件）✓；允许集表 ↔ `c6AllowedSetCases`（五切点两场景逐格一致，含 header-written 断电双态与「corruption 不入允许集」）✓；校验顺序 ↔ `classifyReplaySlotFile`（magic→长度分类→有界 payloadLen，中间态 absent、总长<16/长度不符 corruption，与 A2 ①→②→③ 冻结一致）✓；接缝范围 ↔ 调用点归属 ✓；prebuild-set 格式 ↔ `_c6_enable.go` 字段表（kind/schemaVersion/runId/createdAt/count/requestIds 字典序+recordHash，域 `proofrail:agent-runner-replay-prebuild-set:1\n`）✓；先建全槽后发记录 ↔ `firstReplayPrebuild`（槽全就位后 `publishReplayPrebuildSet`，失败无记录→重开幂等补建）+ 测试钉「重开不重写记录」✓。
3. **边界合规**：三判据本体/解锁边界/平台等级默认态未改（diff 仅改 :230 句并插入 :234）；`agent_runner_replay_store_windows.go:5-6` 仍 `PublishDurabilityUnproven` + `TestC6GateActivationKeepsDurabilityUnproven` 钉死；ADR-013 结论「不改写」句在；`NewAgentRunnerReplayStoreForRun` 签名未动、`PrebuildReplaySlots(requestIDs []string) error` 与 A2 一致；Unix 无独立改动（共享层 +1 字段 + 调度行，nil 缝，linux build/vet 绿）；C2 完整括注 CN/EN 均在（「仅限前提 P′；P′ 覆盖物理机目标部署环境未论证」/ EN 对应句）+ 判定表子注记已落；CONTRACTS{,_EN} 【】=0。
4. **测试质量**：③ 可证伪性强、无假通过形态——BA3 两臂从等价未预建初始态出发、返回序列逐项 + 按 requestId 的 R/C 载荷逐字节比较，且载荷钉死 canonical JSONL 字节（防双空/双错共识）；允许集表测试带「每个允许态必须被观察到」完备性断言；阶段停点经 panic seam 强制（未停即 Fatal）；激活位四断言 + 未激活残留槽 corruption 断言齐备。缺口见 L2/L3。
5. **机械读数（本轮实测）**：默认 gofmt 空、build/vet 绿、`go test` 15 包 ok；`-tags c6slot` adapters ok；linux build/vet 绿；三作用域 gate `TOTAL_FAIL=0/0/0`；编码门禁逐文件核对（md=BOM+LF、go/txt=noBOM+LF）合规；SHA256SUMS 抽查 design-r2b/r2c/design-r2/red-team-r2 四条 hash 与磁盘一致。

## 三、结论

`PRE-REVIEW: PASS WITH FIXES` — 0 High / 0 Medium / 4 Low（L1 悬空引用、L2 未执行断言、L3 漂移分支钉针、L4 措辞限定）。按派单口径「报 Medium+ 时整改后原位重跑」**不触发**；建议主控将 L1–L4 作为 ③/⑧ 收尾项吸收（不改任何语义），不阻塞 ⑦ 终审。
