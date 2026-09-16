# T027 · A7 · Windows 发布耐久候选协议与反证 — 验证报告

日期：2026-09-16。状态：`COMPLETE`（③ 预审 / ③.5 扫描 / ④ 终审见 §9，按切片顺序回填）。
结论（**级别与前提必须一起读**）：在前提 **P = 进程级崩溃、无断电、单主机 Windows 11 / NTFS 本卷** 下——**可见性与仲裁层 `proven`**（P1–P5，构造性论证 + 观测）；**断电耐久 `unresolved`**（U1–U5，含候选 C2/C3 的耐久仍待 B3 注入，**不等同于“候选已被证伪”**）；**被证伪的断言集合 D1–D6 均为 `disproven`**；**Windows 发布耐久保持 `unproven` 不变，first-dispatch 仍被拒绝**。任何摘要句、ADR 与台账均须携带同一前提与同一三档口径，**不得把 `proven` 读成断电级保证**。A7 另有一条**推翻初始假设**的实测发现：Windows 上“目录句柄 flush”并非不可用（见 §4 E4），故 ADR-013 的候选排序须据此修订。

## 0. 诚实准则（本报告所有结论的前提）

1. **进程 kill 留下的是完整的 OS 页缓存，不是断电后的盘面**——被 TerminateProcess 结束的进程不丢任何已完成写入；
2. **`record-linked` 与 `parent-synced` 两个崩溃点在 kill 下观察等价**（§4 E2 实测同分布），因此 A7 的任何绿色实验**都不能**作为"目录项耐久"的证据；
3. 所以 A7 不改变 Windows `unproven`，也不解除 `first-dispatch`；断电类断言一律记 `unresolved` 并移交 B3（真实崩溃/断电注入）。

## 1. 范围与边界

- **做**：ADR 候选枚举、四行反证矩阵、最小 falsification prototype（默认套件模型实验 + `a7native` 真进程实验）、结论账本。
- **不做**（硬门）：不改变 Windows `unproven`；不开放真实 dispatch；**不改动生产发布原语**（`writeReplayRecordNoReplace` 的可见性/持久性步骤语义不变）；不改 CONTRACTS 双语；不改 schema/fixtures/CI 工作流；不接真实候选、不联网。
- 唯一生产文件改动是一处**nil 默认测试接缝** `replayStorePublishStageHook`（§6），用于把崩溃点定到确定阶段——纯计时 kill 无法稳定命中指定阶段。

## 2. 候选协议枚举与 A7 结果

| 候选 | 机制 | 可见性断言 | 耐久性断言 | A7 结果 |
|---|---|---|---|---|
| C1 现行基线 | 现有 `temp → file.Sync → close → os.Link → 父目录耐久步（Windows 为空操作）→ 删 temp` | 原子 no-replace（既有实现） | **无**（Windows 步骤为空操作，靠门禁拒绝写入） | 可见性 `proven`（§4 E2/E3）；耐久 `unresolved` |
| C2 `MoveFileEx(MOVEFILE_WRITE_THROUGH)` | 以带 WRITE_THROUGH 的 no-replace move 取代 `os.Link` | 原子出现全量内容 | 文档称"文件确实移动落盘后才返回"，但同文档另一句把保证限定在 copy-and-delete 移动 | **可用性：实测（证据注释）**（无提权成功、EEXIST 映射 `fs.ErrExist`、开 reader 时拒绝覆盖，§4 E6）；落盘保证 `unresolved` |
| C3 目录句柄 `FlushFileBuffers` | `CreateFile(dir, GENERIC_READ\|GENERIC_WRITE, FILE_FLAG_BACKUP_SEMANTICS)` + Flush | —（持久性步骤） | 声称等价 dir fsync | **A7 推翻"不可用"假设**：只读句柄 `Access is denied`，**读写/只写句柄返回成功且无需提权**（§4 E4）⇒ 与 C2 并列为 B3 注入候选；是否真正保证目录项落盘仍 `unresolved` |
| C4 卷句柄 flush | 开 `\\.\X:` 卷 + Flush | — | 语义上覆盖目录项 | **不予采用**：需管理员/备份特权，与便携产品定位（ADR-009）冲突；E5 未执行（未提权不具代表性，如实登记） |
| C5 两阶段 marker | R 之外再发一个 marker 目录项 | 多一个名字 | 声称 marker 可见性可判提交点 | **作为耐久机制 `disproven`**（§4 E7）：marker 是"新名字=新目录项"，与原语同耐久类；未引入新状态，也未消除跨文件一致面 |
| C6 预建槽位单文件 | bootstrap 预建固定槽，发布=写槽内数据 + 文件级 flush | 槽内长度/标志 | 把目录项耐久转成**文件数据**耐久（有标准保证） | 记入 ADR 作为后续候选；**现行契约下不可即插**（与"每 requestId 独立文件 + no-replace"冲突，仲裁需换锁原语） |
| 排除 | `ReplaceFile`（违背 no-replace 单调发布）、`FILE_FLAG_WRITE_THROUGH` 于 temp 句柄（只覆盖文件数据、不覆盖目录项）、TxF（已弃用）、USN journal 读侧验证（自身同受 flush 语义约束，仅 B3 可选辅助） | | | 见 §3 行 D |

**对称性声明**：Unix 侧 `proven` 也只是"本地 Unix 文件系统"类上的声明（tmpfs/网络文件系统上 `directory.Sync` 同样不保证）。所有耐久结论必须声明在 `(OS, 文件系统, 存储栈)` 三元组上；A7 本机证据仅覆盖 **Windows 11 + NTFS 本卷 + 本机存储栈**。

## 3. 逐项反证矩阵（切片定义的四项）

| 行 | 断言 | A7 可证伪性 | 实验 | 判据 | 结果 |
|---|---|---|---|---|---|
| A | torn / stale writes：已可见记录要么不存在、要么完整且不再改变 | **部分可证伪**（可见性层可证伪；断电级撕裂不可） | E2 kill 矩阵（5 阶段 × N 轮，读侧轮询已发布路径） | 红=出现部分内容/解码失败/valid→different；绿=N 轮全在允许集 | **绿**：前 3 阶段恒 `absent`、后 2 阶段恒 `dispatched-unknown`（§4 E2） |
| B | R 与 marker 跨文件一致性：崩溃后状态 ∈ {absent, R-only, R+C}，无孤儿 C，marker 不引入需新规则的中间态 | **可证伪**（进程级）；断电下目录项**丢失有序**不可证 → B3 | **E9 完成记录崩溃矩阵**（同一原语 + 同一注入方式）+ E7 marker 模型 | 红=出现孤儿或非法态；绿=恒在允许集且 marker 中间态与无 marker 同构 | **绿**（E9）：link 前被杀 → R 存活且状态 `dispatched-unknown`（C 不可见，无孤儿）；link 后被杀 → `terminal-receipt-present`；E7：marker 可见但无收据时状态仍 `dispatched-unknown`；**断电孤儿留 B3** |
| C | 多进程 fencing / reuse：同槽恰一赢家；赢家被 kill 后不卡死、不解锁二次派发；同 hash 幂等、异 hash 冲突 | **可证伪** | E3（进程内 8 并发 + 杀赢家）+ **E10 双进程同轮并发发布同一 requestId**（barrier 仅同步起跑，**同时性未被强制**，见 §4 E10） | 红=双赢家/槽不可读/重放非幂等；绿=不变量恒成立 | **绿**：E3 `first=1 / unknown-block=7`、杀赢家后仍 `dispatched-unknown` 且再调用仍 `unknown-block`；**E10 每轮恰一个 first-dispatch、另一个 unknown-block**，store 保持可读（§4） |
| D | 各 crash point 可证性（元行） | **CP4 与 CP5 在 kill 下观察等价** ⇒ 耐久步骤是否生效**不可证伪**；这是 A7 无法把任何候选升为 `proven` 的根因 | E2 的分阶段分布对比 | 红=两阶段分布不同；绿=分布相同 | **绿（等价成立）**：两阶段均为 `dispatched-unknown` × N ⇒ 耐久断言记 `unresolved`，留 B3 |

## 4. 实验记录（本机 Windows 11 / NTFS；`a7native`，真进程）

> 表尾 `×N` 为**本报告采样轮数**；门禁运行使用 `PROOFRAIL_A7_ROUNDS=20`（请求矩阵/完成记录矩阵/双进程 fencing 均全绿）。两者口径统一说明见 §8。

| 实验 | 主张 | 观察 | 结果 |
|---|---|---|---|
| E2 kill 矩阵 | 每个发布阶段被杀后，replay root 只落在允许集内 | `temp-written`/`temp-synced`/`temp-closed` → `absent`×5；`record-linked`/`parent-synced` → `dispatched-unknown`×5 | 绿；CP4≡CP5 |
| E3 仲裁与杀赢家 | 并发单赢家；赢家被杀不解锁派发 | 8 并发 → `first=1, unknown-block=7`；杀赢家后状态 `dispatched-unknown`，再调用 `unknown-block` | 绿 |
| E4 目录句柄 flush 探针 | （初始假设）目录 flush 在 Windows 不可用 | 只读句柄 → `Access is denied`；**读写句柄 → `<nil>`；只写句柄 → `<nil>`** | **推翻假设**：非提权可调用成功；是否真落盘 `unresolved`（→ B3） |
| E5 卷句柄 flush | 提权下是否可行 | **未执行**（需管理员，与便携定位冲突，不具代表性） | 登记为未执行；采用性判据见 §2 C4 |
| E6 `MoveFileEx` 探针 | C2 可用性、EEXIST 语义、开放 reader 行为 | 无提权发布成功；目标已存在 → `fs.ErrExist`（收敛语义存活）；**目标已存在但未打开 → 同样拒绝**；目标被打开 → `ERROR_ALREADY_EXISTS`；映射到 `fs.ErrExist` 为**硬断言**（不匹配即失败） | 绿（可用性：实测，证据注释）；**open 变量未被隔离**：拒绝的原因是“目标已存在”，与是否被打开无关（报告不声称未测量的因果） |
| E7 marker 模型（默认套件） | marker 不能提升耐久类、不引入新状态 | marker 走同一耐久接缝（`syncCalls>0`）；`PublicationDurability()` 前后不变；状态仍 `dispatched-unknown`；Windows 声明仍 `unproven` 且写入门禁仍拒绝 | `disproven`（marker 作为耐久机制） |
| E8 临时残留 | kill 在 link 前留下的 temp 是否惰性 | 残留 1 个 `.tmp`；不被当作记录、不阻塞后续发布；**无启动期清理 ⇒ 会累积** | 绿 + 发现项（§7） |
| E9 完成记录崩溃矩阵（真 kill，**仅匹配 `completions/` 路径**） | C 发布各阶段被杀后不得产生孤儿、R 不得丢失 | `temp-written`/`temp-synced`/`temp-closed` → `dispatched-unknown`×5（C 不可见、R 存活）；`record-linked`/`parent-synced` → `terminal-receipt-present`×5 | 绿（无孤儿） |
| E10 双进程 fencing（真双进程，同轮并发） | 两个进程抢同一 requestId | 每轮恰一个 `first-dispatch` + 一个 `unknown-block`；store 保持可读，状态 `dispatched-unknown` | 绿；**bidirectional-readiness barrier（双向就绪栅栏）：两子进程各自 ready，父进程确认双 ready 后只放行一次**，5 轮均恰一 first-dispatch（并发语义不变，`unknown-block` 而非 `unknown-unlock`） |

## 5. 结论账本

**级别口径（与切片验收口径一致）**：结论等级**一律取 `proven` / `disproven` / `unresolved` 之一**；`observed`、`falsified`、`disproven-by-policy` **不作为等级**，仅作**证据注释**（如“可用性：实测”、“由实测推翻”、“政策排除”）。`proven` = 在一个明确声明的**前提**下，由构造性论证 + 可复现观测共同支撑。本报告 P1–P5 的 `proven` **一律带前提 P（进程级崩溃、无断电、单主机 Windows 11 / NTFS 本卷）**，不涉及断电耐久。

**构造性论证（P1–P5 的依据）**：

- **P1**：`os.Link`（Windows 下 `CreateHardLink`）在同一目录内原子建立新名字且**目标存在即失败**（`fs.ErrExist` 原样上传），因此“可见”与“独占”是同一个原子步骤的两个方面。
- **P2**：记录先写入同目录 temp、`file.Sync`、`close` 之后才 `link`，即**可见即完整**；写入失败/同步失败时 temp 被丢弃，不存在半写内容可见的路径（E2 采样见 §4/§8；报告采样 5 轮，门禁 20 轮）。
- **P3**：R 的发布先于 C（`RecordCompletion` 先 `ensureRequestLocked`），且 C 走同一 no-replace 原语 ⇒ 崩溃后不可能出现“有 C 无 R”的唯一顺序缺口（E9 采样见 §4/§8；报告采样 5 轮、零孤儿，门禁 20 轮）。
- **P4**：no-replace 碰撞后进入有界重读收敛（≤16 次），并比对 digest ⇒ 同内容幂等、异内容 `ErrAgentRunnerRequestConflict` 且零写（E3/E10 + 既有套件）。
- **P5**：`PublicationDurability() != Proven` 在**全部 8 处写路径入口**即拒绝（request、completion、launch receipt、launch identity、terminal intent、terminal closure、publisher 的发布与恢复），且平台声明不可由写入内容改变（E7）。

| Claim | 级别 | 证据 |
|---|---|---|
| P1 原子 no-replace 可见性（**进程内 8 路真并发** + **跨进程同轮并发尝试**；跨进程强同步未被证明） | `proven`（构造性论证 + 观测，带前提 P） | 论证见上 + E2/E3/E10 |
| P2 已发布路径无撕裂可见 | `proven`（构造性论证 + 观测） | 论证见上 + E2（N 轮零撕裂） |
| P3 进程级崩溃后状态 ∈ {absent, R-only, R+C}、无孤儿 | `proven`（构造性论证 + 观测） | 论证见上 + E2/E7/E9 |
| P4 单赢家仲裁、幂等重放、异 hash 冲突 | `proven`（构造性论证 + 观测） | 论证见上 + E3/E10 + 既有套件 |
| P5 unproven 门禁忠实（拒绝且删守卫必红） | `proven`（构造性论证 + 观测） | 论证见上 + 既有两个平台/调度器测试 + A7 复核 |
| D1 “目录句柄 flush 在 Windows 不可用” | `disproven`（证据注释：由实测推翻） | E4（读写/只写句柄返回成功，**每个可写模式成功均为硬断言**） |
| D2 “卷 flush 对便携非提权产品可用” | `disproven`（证据注释：**政策排除，非实验证伪**） | 需管理员特权（既有知识）+ ADR-009 便携定位；**E5 未执行** |
| D3 "marker 提升耐久类" | `disproven` | E7 |
| D4 "预建槽位在冻结契约下即插" | `disproven` | 契约冲突（§2 C6） |
| D5 "`FILE_FLAG_WRITE_THROUGH` 覆盖目录项" | `disproven` | 语义（只覆盖文件数据）+ E6 观察 |
| D6 "`ReplaceFile` 兼容 no-replace" | `disproven` | 契约（不得替换已发布文件） |
| U1 候选 C2/C3 的断电耐久 | `unresolved` | **A7 不可证**，B3 真实崩溃/断电注入 |
| U2 C2 文档措辞歧义（"落盘" vs "copy-and-delete 移动"） | `unresolved` | 厂商文档自相歧义，A7 无权裁决 |
| U3 本机存储栈（设备缓存/固件）flush 语义 | `unresolved` | A7 不可观测 |
| U4 断电下目录项**丢失有序**（真孤儿） | `unresolved` | A7 完全不可触，B3 最小实验见 §7 |
| U5 共享卷（SMB/集群）fencing | `unresolved` | 不在产品范围 |

**falsifiability（推翻条件）**：任一 kill 轮出现撕裂可见、valid→different 或孤儿 C ⇒ 推翻 P1/P2/P3；删守卫而测试不红 ⇒ 推翻 P5；出现一种对 marker 成立而对记录不成立的同类原语 ⇒ 推翻 D3（理论上不可能）；只有 B3 真实断电证据可推翻 U1（A7 任何实验都不得充当）。

## 6. 改动点与测试点（含变异绑定）

**改动点**：
1. `internal/adapters/agent_runner_replay_store.go`：新增 nil 默认包级接缝 `replayStorePublishStageHook` + 5 个阶段常量与 `replayStorePublishStage()` 调用点（`temp-written`/`temp-synced`/`temp-closed`/`record-linked`/`parent-synced`）。**nil 时零语义变化**。
2. 新增 `internal/adapters/agent_runner_replay_store_a7_test.go`（默认套件，进 CI）。
3. 新增 `internal/adapters/agent_runner_native_a7_test.go`（`//go:build a7native`）。
3.5 新增 `internal/adapters/agent_runner_native_a7_fencing_test.go`（`//go:build a7native`；E9 完成记录矩阵与 E10 双进程 fencing）。
4. 新增 `internal/adapters/agent_runner_native_a7_windows_test.go`（`//go:build a7native && windows`）。
5. `docs/ADR_REGISTER.md` / `_EN.md`：ADR-013 行 + §5 简述。
6. 本报告及其 `_EN` 对照版。
7. 收尾：`docs/DEV_PLAN{,_EN}.md` 与 `docs/t027/REMAINING_SLICES{,_EN}.md` 状态回写。

**测试点与变异绑定**：

| 测试 | 守护的守卫 | 变异动作 | 预期 |
|---|---|---|---|
| `TestA7PublishStageHookIsNilByDefaultAndSeesEveryStage` | 接缝生产期必须为 nil；5 个阶段调用点齐备 | 删除任一阶段调用点 | 红 |
| `TestA7CrashAtEveryStageLeavesOnlyLegalReplayStates` | 可见性边界（link 前后）；状态集封闭 | 让 link 前即可见 / 放宽状态集 | 红 |
| `TestA7MarkerPublicationAddsNoDurabilityClassAndNoNewState` | 耐久类只能由平台声明决定；marker 不产生新状态 | 让 marker 绕过 `replayStoreSyncParentDirectory` / 把 Windows 声明改 proven | 红（Windows 腿） |
| `TestA7LeftoverTempArtifactsAreInvisibleAndDoNotBlockPublication` | loader 只读精确路径、不扫目录 | 让 loader 扫描目录拾取 `.tmp` | 红 |
| `TestA7ReplayIdempotenceAndConflictGuardsStillBind` | `ensureRequestLocked` 的 digest 比对 | 删除 digest 比对 | 红 |
| `TestA7NativePublishCrashMatrix` | 阶段语义与允许集（真 kill） | 让 stage 判定与可见性脱钩 | 红 |
| `TestA7NativeSingleWinnerAndKilledWinnerDoNotUnlockDispatch` | 单赢家 + 杀赢家不解锁 | 让被杀赢家的槽被视为可重派 | 红 |
| `TestA7NativeCompletionPublicationCrashMatrix` | 完成记录（C）发布与请求（R）分离：`completions/` 路径被杀不得产孤儿、不得丢 R | 去掉 `completions/` 路径匹配（回到 R 重发布点）/ 让 C 在 link 前可见 | 红 |
| `TestA7NativeConcurrentProcessFencing` | 双进程同轮并发恰一 first-dispatch（双向就绪栅栏） | 把双向就绪栅栏改回单向起跑门（串行起跑也能绿 ⇒ 伪绿） | 红 |
| `TestA7NativeTempResidueIsInertButAccumulates` | `.tmp` 残留对可见性与仲裁完全惰性 | 让 loader 扫描目录并拾取 `.tmp` | 红 |
| `TestA7NativeDirectoryHandleFlushProbe` | 访问模式依赖性：**每个可写模式必须成功**且只读必被拒 | 把断言放宽为“至少一种模式成功”/ 删除只读必拒断言 | 红 |
| `TestA7NativeMoveFileExWriteThroughProbe` | no-replace 拒绝覆盖 | 允许覆盖 | 红 |

## 7. 已知边界与移交 B3

- **U1/U2/U3/U4 全部移交 B3**；B3 的最小实验：真实断电（或等价设备级注入）+ 每轮崩溃点 journal（落**另一物理设备**且自身可证耐久）+ 重启盘点 R/C 对与目录项；逐候选（C2、C3）分别注入；任一候选要写进契约为 `proven` 必须 B3 全绿。- **a7native 不进 CI**：`a7native` 实验（含 E4/E6/E9/E10）**不进 CI 默认路径也不进 `-race` 子集**（与 `a6native` 同例，防抖动）；本节结论是**单主机**证据，需在 Windows 主机手动重跑（`$env:PROOFRAIL_A7_ROUNDS=20; go test -tags a7native -count=1 -run TestA7Native ./internal/adapters/`）。- **临时残留累积**（E8 发现）：link 前崩溃留下的 `.tmp` 无人清理。本切片**不修**（属生产行为变更，需独立切片决策）；建议 B3/操作切片评估"启动期残留清理"是否引入新的读侧语义。
- **E5 未执行**：卷 flush 需提权，与便携产品定位（ADR-009）冲突；若后续策略变化须重新立项。
- **断言的适用范围**：所有耐久结论仅覆盖 `Windows 11 + NTFS 本卷 + 本机存储栈`；Unix 侧 `proven` 同样只在本地文件系统类上成立。
- **A7 不改变**：Windows `unproven`、`first-dispatch` 拒绝、CONTRACTS 冻结句、T027 `BLOCKED / NOT IMPLEMENTED`、AT-23 结论。

## 8. 成本与计量

- 探针额度：**0/10 使用**（A7 全离线，无付费探针）。
- 模型调用：① V4 Pro 前置分析 ×1；③ V4 Pro 预审 ×1（`PASS WITH FIXES`）+ 复审 ×3（第 1 轮 `RE-REVIEW: FINDINGS`：1 Medium + 4 Low；第 2 轮 3 Low；第 3 轮 `RE-REVIEW: PASS`）；③.5 MAI-Code-1.1-Flash ×2（`INDEPENDENT SCAN: FINDINGS` → `PASS`）；④ Codex ×6 轮（第 1–5 轮 `FINDINGS`，逐轮计数见下表，合计 **1 High + 10 Medium + 5 Low**，均已整改）→ **第 6 轮 `RE-REVIEW: PASS`**；**逐轮结论与整合以 §9 为单一权威**，`×N` 恒等于已回填轮数）；另 1 次 V4 Pro 调用因提供方连接超时失败（不消耗探针额度，按 §3.9 可用性预检口径登记）。
- 实验轮数：默认套件 A7 测试 5 项；`a7native` **7 项**；矩阵按 `5 阶段 × 轮数`（报告采样 5 轮、门禁 20 轮）；双进程 fencing 按 `轮数 × 2 子进程`（同口径）。

## 9. 审查记录（③ / ③.5 / ④）

### ③ V4 Pro 实现预审（2026-09-16）

结论：**`PASS WITH FIXES`**（无 High；无生产语义回归；冻结契约句未被违反）。

| 级别 | 发现 | 整改 |
|---|---|---|
| Medium | M1 行 B 证据引用不实（E2 从不发布 C），P3 “无孤儿”缺 C 崩溃证据 | **补 E9 完成记录崩溃矩阵**（真 kill，5 阶段 × N 轮，零孤儿），并把行 B 引用改为 E9 |
| Medium | M2 头部/寄存器 `proven` 与账本 `observed` 两层不一致 | 账本 P1–P5 升 `proven` 并**补构造性论证段**，明确前提“进程级崩溃、无断电”（§5） |
| Medium | M3 E6 的 open-reader 变量未隔离（目标存在即拒绝，与是否打开无关） | 探针**补对照**（已存在但未打开 → 同样拒绝），报告措辞改为“拒绝的原因是目标已存在，本探针不隔离 open 变量”（§4 E6） |
| Medium | M4 多进程 fencing 缺实验：行 C 仅进程内 8 并发 + 事后杀赢家 | **补 E10 双进程 fencing**（两子进程同刻发布同一 requestId，每轮恰一个 first-dispatch） |
| Low | L1 `a7ReleaseLeakedTempHandles` 静默失败（Windows 清理可能 flaky 且变异不红） | 轮询耗尽即 `t.Fatal` |
| Low | L2 目录 flush 探针 green-on-inversion（`flushAccepted==0` 只 log） | 改为直接 `t.Fatal`（该文件不进 CI，本机 NTFS 上必须红） |
| Low | L3 seam 契约未声明（持锁执行、不得回调 store、覆盖全部 no-replace 发布、失败路径不触发阶段） | seam 注释补齐该契约 |
| Low | L4 `_ = errors.Is` 占位 | 删除并清理 import |
| Low | L5 D2 的 `disproven` 标签过界（E5 未执行） | 改为 `disproven-by-policy`（非实验证伪），并写明 E5 未执行（**④ 第 1 轮再收敛为 `disproven` + 证据注释“政策排除”**） |

复审（整改后重跑 ③）：**`RE-REVIEW: FINDINGS`**（无 High；1 Medium + 4 Low）。

| 级别 | 复审发现 | 整改 |
|---|---|---|
| Medium | **E9 前 3 阶段的 kill 实际落在 R 重发布路径**（`ensureRequestLocked` 对已存在 R 仍先写 temp、link 才碰撞），故“C 前 3 阶段崩溃”未被真正注入 | seam 改为 `func(stage, path string)`（仍 nil 默认、零语义变化），完成记录 helper **只匹配 `completions/` 路径**；重跑后前 3 阶段确为真 C 崩溃点 |
| Low | `_ = errors.Is` 占位残留 | 删除占位与 import |
| Low | §5 P5 论证“四个调用点”计数错误（实为 8 处） | 改为“全部 8 处写路径入口”并点名 |
| Low | §4 计数（×5/×3）与门禁 `ROUNDS=20` 不可对账 | §4 增采样说明、§8 统一口径 |
| Low | E10“同刻发布”未被实验强制（barrier 是起跑门） | 措辞收敛为“同轮并发”，并在 §3 行 C/§4 E10 注明重叠未被强制证明 |

复审（重跑 ③，第 2 次）：**`RE-REVIEW: FINDINGS`**（无 High、无 Medium；3 条 Low，均为文档口径/账本一致性）。

| 级别 | 复审发现 | 整改 |
|---|---|---|
| Low | §3 行 C 仍写“同刻发布”，与 §4 E10 收敛口径及 §9 自述三方冲突 | 行 C 改为“双进程**同轮并发**发布”，并加注“barrier 仅同步起跑，同时性未被强制”；EN 同步为 “in the same round” |
| Low | CN §3 存在**两条重复的 C 行**（旧 E3-only 行未删），与 EN 结构不一致 | 删除旧行，仅保留含 E10 的行（顺带修正行 B 的“跳文件”→“跨文件”） |
| Low | `docs/t027/REMAINING_SLICES{,_EN}.md` 仍写 A7“⬜ 未开始”，与实际进度不符 | 状态表行与小节标题均改为 🔄 进行中（完成回写仍在 ⑥ 执行） |

复审（第 3 次）：**`RE-REVIEW: PASS`**（无 High / Medium / Low；L-A/L-B/L-C 全部关闭，未发现新问题）。② 预审 `PASS WITH FIXES` → 复审×2 `FINDINGS`（1 Medium + 4 Low → 3 Low）→ 复审第 3 轮 `PASS`。

③ 复审的固有局限（写入 Section E 口径）：复审为**只读静态复审**（该环境下无终端执行权限，全目录静态诊断为零），门禁与 `a7native` 实测数字以主机日志为准，未在复审环境独立复跑；E4/E6/E9/E10 为 Windows 单主机证据（§7）。

### ③.5 MAI 低成本独立扫描（2026-09-16）

第 1 轮：**`INDEPENDENT SCAN: FINDINGS`**（2 High + 1 Medium，无 Critical；均为证据层级/措辞，无代码语义问题）。

| 级别 | 发现 | 整改 |
|---|---|---|
| High | 顶层摘要句未携带前提，易被读成“A7 整体 = proven” | 头部结论改为“**级别与前提必须一起读**”，显式声明前提 P（进程级崩溃、无断电、单主机 Windows 11 / NTFS 本卷）与三档口径；ADR-013 行与 §5 同步（不得把 `proven` 读成断电级保证） |
| High | P1 的“含并发”强于 E10 实际测量（barrier 仅为起跑门） | P1 改为“**进程内 8 路真并发 + 跨进程同轮并发尝试**；跨进程强同步未被证明”（§3 行 C 与 §4 E10 同口径） |
| Medium | ADR / 账本 / EN 三处级别口径可能漂移 | ADR-013 行新增“级别与前提”句（`proven` 带前提 P；候选级一律 `disproven`；`unresolved`），与报告头部、§5 逐档对齐；EN 同步 |

第 2 轮：**`INDEPENDENT SCAN: PASS`**（2 High + 1 Medium 全部关闭，无新增；属**文档口径/证据层级**层面的 PASS，不代表产品就绪）。

### ④ Codex 独立终审（2026-09-16，首轮不附任何清单）

第 1 轮：**`RE-REVIEW: FINDINGS`**（1 High + 2 Medium + 1 Low；无 Critical；无生产语义回归、冻结契约句未破）。

| 级别 | 发现 | 整改 |
|---|---|---|
| High | 切片验收口径要求结论只能是 `proven`/`disproven`/`unresolved`，而账本/ADR 引入了 `observed`、`falsified`、`disproven-by-policy` 作为**并列等级** | **等级三态收敛**：报告 §5 顶部重写为“等级一律取三态之一；`observed`/`falsified`/`disproven-by-policy` 仅作证据注释”；D1/D2 行改为 `disproven` + 注释；§4 E6 行的“可用性 `observed`”改为“可用性：实测（证据注释）”；ADR-013 行与 §5 同步（CN/EN） |
| Medium | E10 的子进程由**单向起跑门**释放，串行起跑也能误绿 | **双向就绪栅栏**：两子进程各自写 ready 并停在发布前，父进程等**两个 ready 均到**后只写一次 go 再放行；重跑 5 轮均恰一 first-dispatch |
| Medium | `MoveFileEx` 的“已存在目标 → `fs.ErrExist`”只记日志、无硬断言 | 该分支与 open-reader 分支均改为**硬断言**（`errors.Is(err, fs.ErrExist)`，不匹配即 `t.Fatal`） |
| Low | 请求矩阵用 `fmt.Sprint(map)` 比较 CP4/CP5，map 输出顺序不稳 | 改为**按键逐一比较**两状态计数（`absent` 与 `dispatched-unknown`） |

复审（④ 第 2 轮）：**`RE-REVIEW: FINDINGS`**（2 Medium + 1 Low；无 High/Critical；上轮 4 项整改逐条核对已落地，新问题集中在**文档口径一致性**与**审计计数闭环**）。

| 级别 | 复审发现 | 整改 |
|---|---|---|
| Medium | `DEV_PLAN{,_EN}.md` 中关于 Windows 发布耐久的句子未携带前提 P 与三态口径，存在跨文档口径漂移 | 该句补“**级别与前提必须一起读**：等级一律取 `proven`/`disproven`/`unresolved`，前提 P = …”，并在 DEV_PLAN 新增 **A7 段落**（含前提 P、三态、边界与 `unproven` 保持；中英同构） |
| Medium | §8 计数与 §9 记录不自洽（§8 写“复审×1”，§9 实为 ③ 复审×3；且 ④ 第 2 轮当时“待回填”） | 按实际调用数重写 §8；每轮 ④ 结果**结束后立即回填**本节（本轮即回填） |
| Low | E6 控制组（目标已存在且无人打开）只断言“非 nil”，未断言错误类型 | 控制组补 `errors.Is(closedErr, fs.ErrExist)` **硬断言**，使“有无 open reader 错误类相同”由断言而非推断支撑 |

复审（④ 第 3 轮）：**`RE-REVIEW: FINDINGS`**（2 Medium + 2 Low；无 High/Critical；上轮 3 项整改逐条确认已落地，新问题集中在**轮次计数闭环**与**跨文档候选处置口径**）。

| 级别 | 复审发现 | 整改 |
|---|---|---|
| Medium | §8 写 Codex 3 轮而 §9 仍“第 3 轮待回填”，且 DEV_PLAN A7 段只写“两轮”，三处口径不一 | **单一权威**：轮次明细以报告 §9 为准（§8 只记调用总数与指向），DEV_PLAN 不再自行枚举轮数，只引 §9；每轮结果返回后立即回填 §9 |
| Medium | DEV_PLAN A7 段“候选 C1–C6 一律 `disproven`”与报告/ADR 冲突（C2/C3 的断电耐久仍 `unresolved` 且为 B3 注入候选） | 改为与报告一致：**D1–D6 是被证伪的“断言”集合**；候选 C2/C3 的断电耐久仍 `unresolved` 并移交 B3（C1 为现行基线；C5 marker 作为耐久机制 `disproven`；C6 待契约修订） |
| Low | §6 追踪不完整：声明 `a7native` 7 项但只绑定 4 项，改动点未列 fencing 文件 | §6 补入 `agent_runner_native_a7_fencing_test.go` 与三个缺失实验（完成记录矩阵、双进程 fencing、temp 残留）的变异绑定；E4 绑定改为“放宽为至少一种可写模式成功” |
| Low | E4 文案“读写与只写都成功”强于测试断言（当时只断言“至少一种成功”+ 只读被拒） | **测试提升为硬断言**：每个可写访问模式都必须成功（`writableAccepted == writableModes`），只读仍必被拒 |

复审（④ 第 4 轮）：**`RE-REVIEW: FINDINGS`**（3 Medium + 1 Low；无 High/Critical；上轮 4 项整改逐条确认已落地；新问题均为**候选处置口径**、**证据分母**与**计数一致**）。

| 级别 | 复审发现 | 整改 |
|---|---|---|
| Medium | ADR-013 行与 ADR §5 仍写“候选级结论一律 `disproven`”，与同段“C2/C3 为 B3 候选”及 DEV_PLAN 新措辞不一，可能被读成“C2/C3 已被证伪” | ADR 两处（中英）改为与报告/DEV_PLAN 同词汇：**D1–D6 是被证伪的断言集合**；**候选 C2/C3 的断电耐久仍 `unresolved` 且为 B3 注入候选，不等同于“候选已被证伪”** |
| Medium | E4 探针把「已成功打开」的可写模式计数作分母，若某可写模式 `CreateFile` 失败便被排除在分母外，测试仍可能绿 —— 报告“每个可写模式成功均为硬断言”强于实际断言能力 | 分母**固定为 2**（`const a7WritableModes = 2`）：任一可写模式 `CreateFile` 失败即 `t.Fatal`；每个可写模式都要 `FlushFileBuffers` 成功（`writableFlushed != 2` 即失败）；只读模式仍必被拒 |
| Medium | §8 已声明 4 轮但 §9 仍是“第 4 轮”占位语，闭环在时点不成立 | 本节已回填第 4 轮完整条目（结论 + 发现 + 整改），且 §8 不再预宣告未回填轮次（每轮返回后立即回填，`×N` 恒等于已回填轮数） |
| Low | §5 构造性论证写“E2 三轮零撕裂 / E9 三轮无孤儿”，与 §4 的 ×5 采样、§8 的“报告采样 5 轮”不一致 | 改为“采样见 §4/§8；报告采样 5 轮、门禁 20 轮”，中英同口径 |

复审（④ 第 5 轮）：**`RE-REVIEW: FINDINGS`**（1 Medium；无 High/Critical；上轮 4 项整改逐条确认已落地，Q1/Q3–Q11 均满足，仅**头部摘要的候选处置措辞**未收敛）。

| 级别 | 复审发现 | 整改 |
|---|---|---|
| Medium | 头部结论仍写“若干候选 `disproven`（D1–D6）”，把 D1–D6 表述成候选结论；头部是最高优先阅读入口，会把 C2/C3 误读成“已被证伪候选” | 改为“**被证伪的断言集合 D1–D6 均为 `disproven`**”，并在同句保留“**候选 C2/C3 的耐久仍待 B3 注入，不等同于‘候选已被证伪’**”的防误读锚点；中英同步 |

复审（④ 第 6 轮，终）：**`RE-REVIEW: PASS`**（Medium 及以上为零；Q1–Q11 全部通过且逐条给出文件位置证据）。

复审明确核过的点（均给出位置）：三态口径（§5 与头部）、头部/账本/ADR 行/ADR §5/DEV_PLAN 的前提 P 与候选处置词汇一致（D1–D6 为被证伪断言、C2/C3 仍 `unresolved` 且交 B3）、冻结契约未破且接缝零语义（`replay_store.go` 接缝处与各写入口的 unproven 门禁、CONTRACTS 冻结段）、证据与结论可对账（§3 A–D / §4 E2–E10 / §5 P-D-U；E4 分母固定为 2）、E10 双向 rendezvous + 按轮唯一路径 + stdout 回传、E9 限定 `completions/` 路径、E6 三分支 `errors.Is(..., fs.ErrExist)` 硬断言且不声称未测因果、CP4≡CP5 按键比较、§6 改动点与变异绑定与 7 项实验一一对应、诚实边界齐全、计数闭环（§8 `×N` = §9 已回填轮数；DEV_PLAN 只引 §9）。

残留 Low 观察（复审自述，均为已知边界而非新问题）：复审为静态复核，未独立重跑 `a7native` 与全门禁（无终端）；`a7native` 证据仍为 Windows 单机且不进 CI，后续仍依赖 B3 断电注入与跨环境验证。

MAI 第 2 轮提出的**最小防误读边界句**（已采纳）：**“同一 requestId 在同机、同卷、同轮并发下已证实单赢家且不解锁二次派发；但跨进程强同步、跨主机强同步与断电耐久仍未证明。”**

**MAI 第 2 轮 Section D（仍无测试覆盖的不变量）**：断电耐久（真实断电/设备缓存/固件写回/目录项落盘顺序）；跨进程**强同步**（E10 只证“同机同卷同轮单赢家”）；非本卷存储栈（SMB/共享卷/云盘/集群）；`.tmp` 残留清理策略。**Section E（无法验证）**：真实断电注入与设备层语义（属 B3）。**推翻条件**：再出现“无前提的 `proven`”摘要、把 P1/E10 扩写为“跨进程强同步已证明”、或在未修订 CONTRACTS 前就要求解除 `unproven`。MAI 另确认本轮**未过度收敛**（`proven` 保留在前提 P 下、P1 保留了进程内真并发 + 跨进程同轮尝试）。

### 残差-抽查（§3.11 判据）

**未触发必抽**：A7 唯一生产改动是 nil 默认测试接缝（`replayStorePublishStageHook`，零语义变化），未改变所有权 / 停机 / 身份语义——接缝虽覆盖这些发布，但不改变其行为；因此本片不追加 §3.11 必抽的 Codex 盲审。④ 终审仍按常规执行，并把“接缝的覆盖面与零语义”作为**完整性**要求纳入。
