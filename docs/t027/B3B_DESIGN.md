# B3b 装置与协议设计冻结稿（candidate injection matrix + durability verdict）

> 状态：**设计权威**（本文件冻结 B3b 的协议、矩阵与判定口径；证据包与本报告冲突时以证据为准）。
> 日期：2026-09-20。来源：① V4 Pro（`DeepSeek V4 Pro (deepseek)`，Max）前置分析裁定 + 主控实现偏离注记（§11）。
> 上游：`docs/validation/B3A-CALIBRATION.md`（装置标定，§11 移交条件）与 `docs/t027/B3A_RIG_DESIGN.md`（B3a 装置设计）。

## 0. 目的与范围
对 A7 遗留的 U1/U2（候选 C2/C3 的断电耐久）做**设备级注入裁决**，并按三档口径（`proven`/`disproven`/`unresolved`，带同一前提）给出每个候选的结论。
**不做**：不改 CONTRACTS/schema/fixtures；不动平台耐久声明（Windows 保持 `unproven`）；不解锁 first-dispatch；不清理 `.tmp` 残留；不做并发 fencing 结论。

## 1. 冻结前提 P（结论只覆盖此前提）
- hypervisor：**VirtualBox 7.2.18 r175117**；guest：**Windows 11 专业版 25H2 build 26200.8037**（零售、未激活）；文件系统 **NTFS**。
- 存储：SATA/AHCI `useHostIOCache=false`，两块盘均 `--nonrotational on`；注入卷 **`R:` = INJECT**（固定盘 `data.vdi`，7.98 GB，4096 B 群集）。
- 注入方式：**guest 硬断电**（`VBoxManage controlvm poweroff`），重启至 `guestcontrol` 可用约 48 s。
- 基线快照 `pristine` = `065a25e1-d95c-4f33-8081-4f06216ecce8`（**不重取**：本片不改 VM 配置）。
- 网络：`nic1=none`（客机无网络栈）。
- **信令写域声明**：本片在客机 `C:\prfrail-prep\` 下新增**第二写域**用于阶段标记（§4）。`C:` 是信令域、**不参与测量**；测量域只有 `R:` 上的 replay store 根，且每轮只发生**一次**被测发布。

## 2. 候选与生产接缝
| 候选 | 语义 | 落点 |
|---|---|---|
| **C1** | 现行基线：`os.Link` no-replace + 平台父目录步（Windows 为空操作） | 默认路径本身 |
| **C2** | 用 **`MoveFileExW` + `MOVEFILE_WRITE_THROUGH`（`0x8`）** 的 no-replace move 取代 `os.Link` | 发布原语钩子 |
| **C3** | 用**目录句柄 `FlushFileBuffers`** 取代 Windows 上的空操作父目录步 | 父目录钩子 |
| **C2C3** | C2+C3 组合 | 两个钩子 |

两个接缝（**均为 nil 默认、零语义**，生产路径逐字节不变）：
- `replayStorePublishPrimitiveHook func(temp, path string) error` + 派发器 `replayStorePublishRecordNoReplace`（nil ⇒ `os.Link`）。
- `replayStoreParentSyncHook func(directory string) error` + 派发器 `replayStoreSyncParent`（nil ⇒ 既有 `replayStoreSyncParentDirectory`）。
  **只挂发布路径第 6 步**，不挂 `bootstrapReplayStoreRoot` / `agent_runner_replay_root.go`（bootstrap 不是发布）。
  既有函数变量 `replayStoreSyncParentDirectory` **保持不动**（10+ 处既有测试在替换它）。

**5 个阶段定义不变**（A7 已埋点）：`temp-written` → `temp-synced` → `temp-closed` → **`record-linked`** → **`parent-synced`**。
**候选决定格**：C2 = `{4,5}`（格 4 为核心：move 已返回、父目录步未跑）；C3 = `{5}`；**C2C3 = `{5}`**（**修订 R1**：目录 flush 在父目录步内执行，格 4 时 flush 尚未跑，若格 4 计入 C2C3 主张则该格每轮都缺机制令牌、臂结构性不可证；且格 4 的 C2C3 与 C2 逐字节等价，收窄不损失信息）；**正对照 = `{}`**（不主张耐久，只标定装置）；格 1–3 对所有臂等价（信息性）。

**生产零候选代码**：C1/C2/C3 的候选实现只存在于 `b3bnative` 实验二进制（§11 偏离注记）。
**绕开门禁声明**：写入器直接调用 `writeReplayRecordNoReplace`，并在此二进制内把 `store.publicationDurability` 置为 `proven`（与既有实验同款测试期覆盖）。**门禁是策略层、不是发布原语**，此绕过不影响被测对象，但必须在报告与证据包中逐句声明。

## 3. 客机写入器（`b3bnative` 测试二进制）
构建：host 即 `windows/amd64` ⇒ `go test -c -tags b3bnative -o writer.exe ./internal/adapters/`（静态、单文件），经既有 `guestcontrol copyto` 部署。

**env 契约**：`B3B_GUEST_WRITER`（必设，否则 skip）、`B3B_ROUND`、`B3B_CANDIDATE`（`c1|c2|c3|c2c3|positive-control`；**修订 R2**：控制臂必须由同一写入器驱动，否则区分力门结构性永假 ⇒ 全矩阵只能以错误理由得 `unresolved`）、`B3B_STAGE`（5 阶段之一）、`B3B_ROOT`、`B3B_MARKER`、`B3B_PLAN_ONLY`、`B3B_BLOCK_SECONDS`（默认 300）。
**记录确定性**：requestID = `b3b-<round>`；记录由 `replayStoreRequestRecord` 构造（`CreatedAt` 为钉死常量 `2026-09-10T02:00:00.000Z`）⇒ **宿主可预承诺**。

**标记文件**（`B3B_MARKER`，JSON Lines，逐行追加；**只要求客机内可见**，不要求耐久）：
- 首行 = **预承诺** `{round,candidate,stage,path,digest,size}`；`digest` = 发布字节（canonical JSON + `\n`，由 `evidence.EncodeCanonical` 复现）的 SHA-256。
- 之后每阶段一行 `{round,candidate,stage,path,target,trace[],timestamp}`；`trace` = 至该阶段为止的 syscall 轨迹（`os.Link` / `MoveFileExW(WRITE_THROUGH,no-replace)` / `platform-parent-sync` / `FlushFileBuffers(dir)`）。
- **顺序硬约束**：先写标记、后阻塞。

**阻塞与超时**：目标阶段写入标记后 `time.Sleep(B3B_BLOCK_SECONDS)`（默认 300 s，**不用 `select{}`**，避免 runtime 死锁检测误杀——A7 教训），超时打印 `TIMEOUT <stage>` 并以**专用退出码 90** 退出，**绝不自续推进**。
**输出行**：`PLAN …` / `TIMEOUT <stage>` / `SURVIVED <stage>`（后者表示宿主未在目标阶段内断电 ⇒ void）。

## 4. 宿主握手协议
1. **S6.5 轮询**：宿主以 ~2 s 间隔用 `guestcontrol run cmd /c type <marker>` 读取标记；**不以 stdout 流为主通道**（`guestcontrol` 是否流式回传中间输出未经核实）。
2. **cut 前置**：目标阶段标记已见 ∧ 无更晚阶段 ∧ `plan` 记录已 ACK ∧ 2.0 s 冻结延时已过 ∧ **存活探针**（`tasklist` 显示写入器仍在）。
3. **cut**：`VBoxManage controlvm poweroff`。
4. **阶段命中证明链**（四条同时成立）：① 宿主 cut 前归档的标记转录末条 = 目标阶段；② `plan` 记录 ACK 早于 cut；③ cut 时写入器仍存活；④ 重启后盘点落在该阶段的允许集内。
5. **void 三类 + 一条会话级阻断**：`marker-missed`（窗口内未见目标阶段）、`writer-exited-before-cut`、`done-before-cut`（见到 `SURVIVED`）；**`stage-escape`**（见到**晚于**目标阶段的标记）= 协议缺陷 ⇒ **会话级阻断**，不得继续矩阵。
6. **时序锚点变更**：2.0 s 数值保留，但**计时起点改为"检出目标阶段标记后"**（B3a 为"写入返回后"）。
7. **轮前预授权**：每轮开始前确认本轮 `{候选, 阶段, 目标路径}`；`-NoConfirm` 为会话级选项（与 B3a 一致）。
8. **同格 void 上限 3 轮**（沿用 B3a §3.3），超出即停并上报。

## 5. 轮次协议 S1–S9（B3b 变体）
| 步 | 内容 | 与 B3a 的差异 |
|---|---|---|
| S0 | 会话前置：版本断言、VM 事实、磁盘标志、**lib 哈希 pin**、函数碰撞检查、journal 链自证 | 新增 pin 与碰撞检查 |
| S1 | `restore pristine` | 同 |
| S2 | `startvm` | 同 |
| S3 | 就绪等待（`guestcontrol` 可用） | 同 |
| S3.5 | 工具部署：**writer.exe** + `11-shutdown-audit.ps1` + `12-round-inventory.ps1`，客机 `Get-FileHash` 比对 | 由 3 个 PS1 改为 exe + 2 个 PS1 |
| S4 | 注入前审计（uptime、许可、卷状态） | 同 |
| S4.5 | 轮目录布置（`C:\prfrail-prep\runs\rNN\`） | 信令域 |
| S5 | **`plan` 记录 + ACK**：以写入器 `-plan` 模式取得 `{path,digest,size}`（不发布），落 journal 并 ACK | **先记后写语义恢复**（可比 B3a） |
| S6 | 启动写入器（真实发布 + 目标阶段阻塞） | 由 PS1 写入改为驱动 exe |
| S6.5 | 标记轮询 + 转录归档 + 存活探针 | **新增** |
| S7 | 预授权 → 冻结延时 2.0 s → `controlvm poweroff` → 归档 `cut` 记录 | 计时起点变更 |
| S8 | 重启 → `11-shutdown-audit.ps1` → 盘点 ×2（`12-round-inventory.ps1`） | 同（原样复用） |
| S9 | 逐轮对账：发布字节 vs 预承诺摘要、允许集判定、void/有效标记 | 判定对象由"探针摘要"改为"记录路径 + 预承诺摘要" |

## 6. 矩阵与预算（N=5）
| 臂 | 格 | 轮数 |
|---|---|---|
| 标定前奏 | 正对照（基线原语 + 卷句柄 flush；无卷句柄写权限时回退为记录句柄 flush 并在 trace 单独标记，**修订 R2**） | 5 |
| C1 参照臂 | 5 阶段 × 5 | 25（**负对照并入阶段 1 格**） |
| C2 | 5 阶段 × 5 | 25 |
| C3 | 5 阶段 × 5 | 25 |
| C2C3（可选，已批） | 仅 `parent-synced`（**修订 R1**：格 4 合计入则结构不可证）× 5 | 5 |
| U4 子矩阵（可选，已批） | 3 候选 × 阶段{4,5} × 3（双发布序列：R 完整 → C 目标阶段断电） | 18 |
| 变异臂（装置自证伪） | 故意无 WRITE_THROUGH / flush 改 no-op（**人工臂**，不在矩阵调度器内） | 2 |
| 握手冒烟 | 通道定型 + δ 分布实测（独立 smoke 会话，不入主证据） | 3 |
| **合计** | 主矩阵 17 格×5 + U4 6 格×3 = **103 轮**；加冒烟与变异共 **108 轮** | **≈6.6–8.3 h** |

单轮 ≈ 3.8–4.6 min（含两次冷启动等待、两次盘点、cut 前后探针；B3a 实测 11 轮 38 m 42 s 为 256 KiB 载荷口径）。
**④ 终审复核（Codex）已指出**：6.5–7 h 属乐观边界，不是稳态值；含 void/运行重试时 8 h 以上更现实。长跑前须在离午夜 ≥9 h 的窗口内启动，且矩阵已把会话 ID 定一次后透传（**修订 R15**），因此跨午夜不会再分裂会话。

## 7. 三档判定口径
**允许集（每格 = 候选 × 阶段 × N 有效轮）**
- 阶段 1–3：`{记录不可见（absent，`.tmp` 残留惰性）, 记录完整可见（survived）}`；违反形态 = **撕裂**（`lost-torn-size`/`lost-torn-content`）。
- 阶段 4–5：`{survived}`；违反形态 = `lost-absent`、**`lost-zero-filled`**（目录项在而数据全零 ⇒ "原语返回即落盘"被直接推翻）、`lost-torn-*`。
- **逐轮全局不变量**（任一违反 = 该轮反例）：无孤儿 C、审计 `HARD-POWER-LOSS`、盘点 ×2 逐字节一致、卷不脏、journal 全链完整。
- 阶段 1–3 出现 `survived` **不是**反例（早于声明的耐久是加强）。

**机制执行证据（每轮强制）**：写入器轨迹必须命名该候选的真实原语（C2：`MoveFileExW(WRITE_THROUGH,no-replace)`；C3：`FlushFileBuffers(dir)` 位于阶段 4/5 之间；C1：`os.Link` 且无目录 flush）。**缺条目 ⇒ 该轮对该候选无效**。

**区分力门（防假 `proven`）**：同一 Go 写入器 + 同一记录级载荷下，负对照（C1 阶段 1）出现 **≥1 丢失**、正对照 **0 丢失**。门不成立 ⇒ 所有候选格**封顶 `unresolved`**。
（依据：B3a 的 256 KiB 载荷对照按 §10.1 **不得跨载荷迁移**。）

**`disproven`**：出现任何违反允许集的反例（≥1 次，须全门禁有效 + 违反类），或候选声称的保证被直接推翻。建议首个反例补 1–2 轮复现（记录，不改档位）。
**`proven`**：五条件全真——① 全部 5 格 × N=5 有效轮落在允许集；② 决定格机制证据齐备；③ 区分力门真；④ 全局不变量真；⑤ journal 全链与 plan/cut/reconcile 完整。任一缺 ⇒ 最高 `unresolved`。
**`unresolved`**：任一格有效轮 <5 且无反例；阶段归属无证据；区分力门失败；结论依赖不可证伪断言（如厂商文档语义）。
**U4**：子矩阵任一有效轮出现"**C 在而 R 不在**（真孤儿 C）" ⇒ 该候选 `disproven` + U4 `disproven`；全 0 孤儿 ⇒ U4 只能 `unresolved`（有限样本不能证有序，只能证伪）。

## 8. 装置落点与复用边界
**新建 `tools/b3b-rig/`**，B3a 的两份 lib 作为**只读依赖**点源复用（不改其字节），从而保证：B3a 证据包的 `rig/` 冻结副本仍等于仓库字节，且 journal 约定只有一份实现。
**硬边界**：
1. **不得新增 journal 字段**：冻结的 canonical 构造只遍历字段序列表，未列键**静默丢弃**；`New-B3aJournalRecord -Op` 被 `ValidateSet('plan','cut','reconcile')` 锁死。⇒ B3b 语义**映射到既有键**（下表），阶段命中/机制轨迹存**轮目录工件**。
2. **lib pin**：S0 断言两个 lib + `guest/12` + `guest/13` 的 SHA-256 = 启动冻结值，漂移即 fail-closed；证据包归档运行当次副本并逐字节核对。
3. **点源纪律**：被点源的模块**不得有顶层 `param()`**；B3b 新模块同规则；自检比对点源前后函数清单差集（防同名静默覆盖）。
4. `Invoke-B3aRound.ps1` **不可复用**（`-Candidate` 被 `ValidateSet` 锁死、S4.5/S6 固定）⇒ 自写 `Invoke-B3bRound.ps1`；`guest/12-round-inventory.ps1` 与 `guest/11-shutdown-audit.ps1` **原样复用**。

**journal 字段映射表**
| 既有键 | B3b 语义 |
|---|---|
| `candidate` | `c1`/`c2`/`c3`/`c2c3`/`positive-control` |
| `stage` | 5 阶段之一 |
| `targetPath` | 被测记录路径（`R:\…\requests\request.b3b-rNN.jsonl`） |
| `contentDigest` / `size` | 预承诺的发布字节摘要 / 字节数 |
| `outcome` / `verdict` / `voidReason` | 格判定与 void 原因 |
| `round` / `timestamp` | 同 B3a |

## 9. 证据包结构
`docs/validation/evidence/b3b-<date>/`（**修订 R5**：目录结构改为对真实会话根布局的忠实映射，因为驱动器的 `journal.jsonl`/`session.log`/`pin.json`/`matrix.json`/`verdict.json` 都直接位于会话根，只有 `rounds/` 是子目录）：
`README.md`（含布局映射、`claimScope` 逐字引用、链验结论与可复算边界）、`MANIFEST.md`、`SHA256SUMS.txt`、
`session/{session.log,pin.json}`、`journal/{journal.jsonl,chain-verify.txt}`（由 `journal.jsonl` 改名入包）、
`rounds/rNN/`（标记转录、cut 前归档、writer stdout/stderr、审计、盘点 ×2、对账、cut 记录）、
`matrix/{matrix.json,verdict.json}`、`calibration/`・`u4/`・`variant/`（存在则复制）、
`verification/{rig-vs-repo.txt,analyzer.txt}`、`rig/`（运行当次 B3b 装置 + 被复用 B3a 文件副本）。

## 10. 风险与不可证边界
1. **装置区分力风险（最大）**：记录级（~1 KiB）未刷目录项在 2 s 窗口可能高概率存活 ⇒ 主张格与参照格同为零丢失 ⇒ 只能 `unresolved`。B3a 的 4/5 只属 256 KiB 载荷。
2. **统计分辨率**：格内 0 丢失的单侧 95% 上界 p<0.45 ⇒ `proven` 只能读作"该前提下、该装置分辨率内未观测到反例 + 机制确执行 + 装置有区分力"，**不得**读作长期保证。
3. **阶段归属依赖协议证据**：阶段 4/5 在盘点观察上等价（A7 CP4≡CP5）⇒ 标记通道失效即只能 `unresolved`。
4. 结论只覆盖前提 P；不覆盖物理硬件、其它 hypervisor/存储栈、宿主自身掉电（U3）、SMB/共享卷（U5）。
5. 不覆盖：`.tmp` 残留清理、并发多写发布、resume/结算等其它链路。
6. 即使 `proven`：B3b 不构成 AT-23 证据，不解锁 first-dispatch，不解除其它平台 `unproven`；契约修订属**条件性后续**子步骤。

## 11. 实现偏离注记
- **D1（候选原语落点）**：① 建议落生产文件（`//go:build windows`、不接线）；主控改落**测试作用域**（`_test.go` + `windows` tag）。理由：避免把未上线的发布原语带进产品二进制（架构规则禁止投机接口），同时保留 Windows 腿 CI 覆盖；`proven` 后由实现切片提升为生产字节（同源即可）。**待 ③ 预审裁定**。
- **D2（`C:` 信令域）**：B3a 无第二写域；本片为阶段握手新增，已在 §1 声明。
- **D3（门禁绕过）**：写入器在实验进程内置 `publicationDurability=proven` 并直呼 `writeReplayRecordNoReplace`；已在 §2 声明。
- **D4（证据包可选分区，2026-09-21 实测）**：本装置的所有格——**含正对照标定格与 U4 配对轮**——全部写入 `<会话>/rounds/`（配对轮以 `paired=true` 标记），**不产出** `calibration/`・`u4/` 子目录。
  证据包 README 会如实列出 `optionalSectionsAbsent=calibration,u4,guest`；各格角色与 U4 汇总见 `matrix/verdict.json` 的 `cells` 与 `u4`。
  **读取约束**：不得把"可选分区缺席"误读为"标定或 U4 未执行"；也不得因此声称证据不完整（`rigMismatches=0 missingRequired=0`）。
  唯一存在可选分区的实例是 `variant/`（变异臂产物，已按 D1 并入并附自撰 `README.md`）。

## 12. 复现命令（宿主 PowerShell 5.1）
```powershell
# 构建客机写入器
go test -c -tags b3bnative -o tmp\b3b\writer.exe .\internal\adapters\
# 原生自测（无需 VM）
go vet -tags b3bnative .\internal\adapters\
go test -tags b3bnative -count=1 -run TestB3bNative .\internal\adapters\
```

## 13. ③ 预审整改记录（V4 Pro `changes-required` → 已闭环）

预审判定 `changes-required`，F1/F2/F3 为阻断级。逐条整改如下；每条都落到代码，并由 `tools/b3b-rig/Test-B3bRig.ps1`
的对应断言钉住（自检 **`checks=103 failures=0`**，含 6 个篡改拒裁夹具 + 1 个 cut 后 marker 截断容忍夹具 + 6 条入库证据包断言）。

> **计数口径注记（⚠、2026-09-20）**：本表早期版本写 `checks=61`/`4 个夹具`，已过期。断言数与夹具数随四轮整改增长：
> 当前 **103 checks / 6 篡改拒裁夹具 + 1 截断容忍夹具**。另有**证据包断言 6 条**（`bundle-*`，已入库，见 R30），
> 早期"证据包冒烟 22 项一次性手工运行、尚未入库"的说法在 R30 后作废。引用时一律以实跑末尾的 `checks=... failures=0` 为准。

| # | 预审发现 | 整改 | 钉住它的断言 |
| --- | --- | --- | --- |
| R1 | c2c3 决定格 `{4,5}` 与机制令牌规则自相矛盾 ⇒ 结构性不可证、10 轮空耗 | 决定格收窄为 `{5}`（§2 已改），组合格列表同步 | `decision-c2c3-parent-synced-only` |
| R2 | 正对照无驱动器 ⇒ 区分力门**确定性永假**，全臂只能以错误理由 `unresolved` | 写入器接受 `positive-control`（基线原语 + 目录 flush + 卷句柄 flush，无卷句柄权限则回退记录句柄并单独标记）；装置侧 `Get-B3bArmList`/`Get-B3bControlList` 纳入控制臂，控制臂不主张决定格 | `control-arm-cell-is-schedulable`・`arm-list-has-control`・`decision-control-claims-nothing` |
| R3 | 判定管线信任 `reconcile.json` 自身：`JournalChainOk` 硬编码 `$true`、journal 只记不门控 ⇒ 改两个 JSON 即可把 `disproven` 翻成 `proven` | 结果**重派生**（`plan.json`+`inventory-1.json`，与对账同源同一实现）+ 逐文件对位比对 `reconcile.json`/journal 记录的哈希（含 `marker-transcript.txt`/`marker.txt`/`audit.txt`）；不一致 ⇒ 该轮 `evidence-mismatch` 且 finalize **拒裁 exit 4** | `refuses-tampered-reconcile-outcome`・`refuses-tampered-inventory`・`tamper-changed-the-artifact` |
| R4 | 阶段归属读 cut **后**拷回的 `marker.txt`（尾行可能正被断电丢掉）；`stageEscape` 只记不阻断；违反轮在缺机制证据时仍判 `disproven` | 归属改以 cut 前归档 `marker-transcript.txt` 为主、差异显式记账；逃逸 ⇒ `sessionBlock` 且 finalize **exit 6**；反例要求 `stageHit && mechanismOk`，否则该格封顶 `unresolved` | `phase-filter-*`・`refuses-session-block`・`verdict-unattributable-violation-unresolved` |
| R5 | 证据包假设的 `journal/`、`matrix/` 子目录在真实会话根不存在 ⇒ 真会话必被 REFUSE | 包结构改为真实布局的忠实映射（§9 已改），并把链验结论写进 `journal/chain-verify.txt` 与 README | 证据包冒烟 22/22（`journal-copied`・`matrix-json-copied`・`chain-verify-ok` 等） |
| R6 | 复用 pin 只在 S0 查一次；证据包 `rig-vs-repo` 拿当前哈希比当前哈希（同义反复） | 证据包改为与 S0 记录的 `session/pin.json` **真比对**，缺 pin 记 `UNPINNED` 且计入 mismatches；`pin.json` 入包 | **证据包冒烟**（非自检）：`pin-drift-detected`・`rig-vs-repo-compares-against-session-pin` |
| R7 | U4 的 500 号分界纯数字约定、对账不记配对性 ⇒ 错号轮可把真反例移出臂裁决 | 对账记录 `paired`/`targetKind`；cell 模式强制区间（paired ⇒ ≥500，主矩阵 ⇒ <500）；finalize 按 `paired` 分类并断言号码区间与 `targetKind` 一致，不一致 **exit 4** | `refuses-paired-numbering-mismatch` |
| R8 | 区分力只在阶段 1 得证，却可让阶段 4/5 的 `proven` 通行；门失败无中途止损 | 判定文档强制携带 `claimScope`（逐字写明"仅阶段 1 证明区分力/非长期保证"）；新增 `-Checkpoint` 中途检查点（门不成立即 exit 4 并要求操作者决策） | `verdict-document-carries-claim-scope`・`checkpoint-exists-before-the-budget` |
| R9 | 臂级 `disproven` 把信息格（1–3）的违反也算到候选头上 | 仅决定格可推翻臂；信息格违反记 `infoViolations` 并封顶 `unresolved`（**修订 R23**：初版只记账未封顶，且自检把错误行为钉成规范） | `arm-info-cell-caps-the-arm-at-unresolved`・`arm-info-clean-can-still-be-proven`・`arm-c1-never-proven` |
| R10 | 存活探针与断电之间存在无复查窗口；写入器超时不改变轮有效性 | **先冻结延时、再复探、紧接着断电**（复探→cut 间隔仅为探针时长，`survival-probe-2.txt`，失败即 void）；cut 后从 writer stdout/stderr 回收 `TIMEOUT <stage>` 并转为 void | `second-cut-probe-runs-on-every-path`（AST 包含关系：不得位于 `-NoConfirm` 分支内） |
| R11 | 候选系统性失败被洗成 `marker-missed` | 检测超时时解析已归档 stdout/stderr，区分 `writer-failed-before-stage`（带错误文本）与真静默 | 同上 |
| R12 | `tasklist` 探针用宿主文件名而非客机固定名 | 探针改用客机固定名（由 `$GuestWriterExe` 派生） | 同上 |

**D1 裁定（预审意见，已接受并写入报告约束）**：候选原语留在 `windows` tag 的测试作用域对 `proven` 主张**不构成无效**——被测对象是"生产发布路径 + 经 nil 默认接缝注入的原语"。使主张失效的条件（报告必须遵守）：(a) 后续提升切片改动原语语义（如给 `MoveFileExW` 加 `REPLACE_EXISTING`、丢弃 `WRITE_THROUGH`、放宽目录句柄访问模式）而不重跑矩阵——`TestB3bCandidateMoveFileExNoReplace` / `TestB3bCandidateDirectoryFlush` 钉死的契约必须逐字保留；(b) 报告把"被测量原语存活"写成"产品现已耐久"（Windows 上产品仍为 `unproven`，必须逐句声明）；(c) 提升前删除测试文件导致溯源断裂（当前由 `writer.exe` SHA-256 在 S0 与证据包两处钉住）。

**第二轮独立扫描（MAI ③.5）追加整改**：

| # | 发现 | 整改 | 钉住它的断言 |
| --- | --- | --- | --- |
| R13（blocker） | void 路径不做完整性校验：把实测轮改写为 `verdict='void'` 即可把它移出有效集而不触发 exit 4 | void 轮同样必须与 journal 记录一致（`contentDigest`/`candidate`/`stage`，以及 `extra.verdict`/`extra.voidReason`/`extra.outcome` 只要存在就必须一致），且理由必须属于驱动器可产生的 void 词表（`Get-B3bVoidReasonList`）；不满足 ⇒ 该轮 `evidence-mismatch` ⇒ finalize exit 4。注：该篡改方向的最坏后果是 `unresolved`（有效轮不足 N），不是假 `proven`，但仍按完整性缺陷封闭 | `refuses-tampered-void-relabel` |
| R14（major） | 完整性门仍把 cut **后**拷回的 `marker.txt` 当权威证据 ⇒ 断电截断尾行会误拒诚实轮 | `marker.txt` 退出完整性门（属告知性证据，`markerPostCutDiffers` 已逐轮记账）；阶段证明只认 cut 前归档 `marker-transcript.txt`（`markerPreCutDigest`）与 `audit.txt` | `post-cut-marker-truncated-fixture`・`matrix-c3-proven-despite-truncated-post-cut-marker` |

（第三项 `inventory-2.json` 必填的顾虑经复核不成立：代码只在 reconcile 确实记录到第二份摘要时才要求它；实盘第二轮盘点在重启后运行，不存在被同一次断电截断的路径。）

**第三轮（④ Codex 终审 + 定向复审）与第四轮（收口复审）追加整改**：

| # | 发现 | 整改 | 钉住它的断言 |
| --- | --- | --- | --- |
| R15 | 会话目录按当日日期动态计算 ⇒ 跨午夜静默分裂会话 | 三脚本 `-SessionId` 贯通（矩阵解析一次并下传每轮）；预检允许已有 pin 胜出，**显式**冲突才拒绝；预检与矩阵同持会话锁 | `session-id-resolved-once-and-threaded`・`session-pin-pins-the-session-id`・`session-id-conflict-is-refused`・`preflight-takes-the-session-lock` |
| R16 | `plan-no-ack` 不在 void 词表，且静态提取正则抓不到括号拼接形式 | 词表补齐；正则改为 `Invoke-B3bVoid \(?'…'` 并裁剪尾冒号；加 AST 裸命令守卫 | `void-vocabulary-covers-driver`・`void-vocabulary-is-not-empty` |
| R17 | S0 未探客机存储根可写 ⇒ 权限/挂载问题拖到 S4/S6 才暴露 | 预检新增 `9-guest-store-root-writable`：先拉起 VM 并等 `guestcontrol` 就绪，GUID 后缀唯一文件名，guest 内 `$ErrorActionPreference='Stop'` 且**不无条件 exit 0** | S0 实机全绿（实机事实，非静态断言） |
| R18 | 预算口径与实现不一致 | §6 修正为 c2c3 1 格×5、变异臂为人工臂、主矩阵 17 格×5 + U4 6 格×3 = 103 轮；单轮 3.8–4.6 min、总 ≈6.6–8.3 h | 无（文档） |
| R19 | 我的整改自身引入 blocker：注释缺 `#`（能过语法解析、运行时当命令执行） | 修正注释，并新增 AST 守卫：任何无连字符的裸命令名若非已定义函数且不可解析即转红 | `no-unresolvable-bare-command` |
| R20 | 探针命令无条件 `exit 0` ⇒ 写入失败可伪装成功 | guest 内 `$ErrorActionPreference='Stop'` 且不再无条件退出 0；失败即预检失败 | 实机 S0 探针（R17） |
| R21 | 固定探针文件名可能覆盖既有文件 | 改为 GUID 后缀唯一名，先创建后删除 | 实机 S0 探针（R17） |
| R22 | pin 冲突静默采纳（掩盖真实不一致） | 保留原始入参标记；仅**显式** `-SessionId` 与 pin 冲突时拒绝 | `session-id-conflict-is-refused` |
| R23 | 信息格违反只记账、未封顶（自检还把错误行为钉成规范） | `proven` 追加 `-and ($infoViolations.Count -eq 0)`；自检改为双向（违反 ⇒ `unresolved`；信息格干净 ⇒ 仍可达 `proven`） | `arm-info-cell-caps-the-arm-at-unresolved`・`arm-info-clean-can-still-be-proven` |
| R24 | U4 的 `orphanCompletions` 未重派生 ⇒ 改一个字段即可把真证伪吞成 `unresolved` | finalize 改为从该轮 `inventory-1.json`（哈希已 pin）经 `Get-B3bReplayPairing` 重派生并对位，不一致 ⇒ **exit 4** | `refuses-tampered-u4-orphan`・`u4-tamper-applied` |
| R25 | R24 的初版对每个 paired 轮无条件要求盘点 ⇒ **诚实 void 的 U4 轮会被误拒** | void 轮跳过对位（void 路径恒早于盘点落盘，其篡改面由 void 词表 + journal 对账覆盖） | `void-u4-round-does-not-block-finalize` |
| R26 | `Start-Sleep` 仍位于复探之后（复探→cut 仍 2 s；注释与代码相反） | 顺序改为冻结 → 复探 → 断电 | `second-cut-probe-runs-on-every-path`（含实机 `survival-probe-2.txt`） |
| R27 | 零运行 VM 时全清 headless 为 fail-open（别的 VM 可能正在启动尚未入列） | 仅清「命令行匹配本 VM」或「无任何 VM 运行时命令行不可识别」的进程；`list runningvms` 失败即 fail-closed 不清 | `stale-session-match-is-boundary-anchored` |
| R28 | 会话锁的 PID 复用盲区未文档化 | 见下方"会话锁恢复动作" | — |
| R29 | ④ 收口复审：预检在**终止异常**路径上不释放会话锁（缺 finally 兜底）⇒ 留活 pid 锁，下一次诚实运行被 `REFUSED` | 脚本级 `trap`（覆盖所有终止错误，含助手内抛出）：`$SessionDir` 非空即释放 → 打印 `preflight aborted:` → `exit 1`；选 trap 而非 try/finally 是为避免把主流程整体缩进（语义等价、改动面最小） | `preflight-releases-the-lock-on-a-terminating-error` |
| R30 | ④ 收口复审：证据包构建器的 premise 回退与该构建器的"拒绝写出仓库根"守卫**均无入库断言**（此前靠一次性手工 smoke，不可重复） | 自检新增 K 段：合成最小会话 → 实调构建器，断言 `bundle-fixture-journal-ack`・`bundle-builds-from-the-repository-suite`・`bundle-layout-complete`・`bundle-readme-carries-the-premise`・`bundle-marks-a-missing-premise`（无 premise ⇒ README 必须写 `UNKNOWN (verdict.json carried no premise)`）・`bundle-refuses-to-write-outside-the-repo-root`；构建器拒绝是**终止错误而非退出码**，故统一经 `Invoke-B3bEvidenceBuilder` 包装（否则一次拒绝会打死整套自检，实测已发生） | 上述 6 条 + `bundle-*` |
| R31 | `-Cell` 的校验正则 `^([a-z0-9]+):([a-z-]+)$` 不接受含连字符的臂名，而控制臂名就叫 `positive-control` ⇒ **正对照标定格从 CLI 结构性不可调度**（判别力门的正半永远跑不了，所有臂会被非经验原因封顶为 `unresolved`）；而断言 `control-arm-cell-is-schedulable` 只 grep 了下一行的 `(Get-B3bArmList) -contains $arm`，因此一直是绿的 | 抽出单一权威解析器 `Resolve-B3bCellSpec`（臂名允许连字符 + 成员校验），矩阵与自检共用；断言改为**真调解析器**（行为式），并新增 `cell-resolver-still-refuses-unknown-arms`・`matrix-uses-the-shared-cell-resolver` | `control-arm-cell-is-schedulable`・`cell-resolver-still-refuses-unknown-arms`・`matrix-uses-the-shared-cell-resolver` |
| R32 | 会话号缺省回落为"当天日期"，而 pin 又存在会话目录**内部** ⇒ `-Finalize` 漏传 `-SessionId`（或跨零点后）会去另一空目录找轮次，**零轮却 exit 0 写出空判定**（读起来与诚实的 `unresolved` 完全同形，静默丢掉整场证据）；自检自身也因夹具写死 `2026-09-20` 而有时效性——跨零点后实测 7 条断言变红 | finalize 新增 fail-closed：解析出的会话目录下零轮即 `REFUSED` + exit 4；自检夹具会话号改为变量并**显式贯通**到每一次矩阵调用（含 finalize 与篡改夹具） | `finalize-refuses-a-session-with-no-rounds`・`finalize-honours-the-requested-session-id` |
| R33 | 盘点契约实现得比设计**本意更严**：B3a §2.3 要求的是"两次盘点逐字节一致"以证明**枚举态已定**、而非"头两次就读到定态"。实机 c1（temp-written）首读可能仍处未定态 ⇒ 三轮诚实负对照全部 void 且撞满预算；且旧代码在比较失败时**不落盘两份字节**，导致无法区分"存储真在变"与"读取时序伪影" | 新增 `Get-B3bSettledInventoryPair`（**向前扫描首个相邻一致对**，上界 `$B3bInventoryMaxReads = 4`）；轮驱动改为有界循环、每次读取均存档（`inventory-attempt-N.bin`）、记录稳定所需次数（`inventory-reads.txt`）、并以**稳定对**写 `inventory-1/2.json`；未稳定仍 void 且保留 `inventory-mismatch-1/2.bin` | 行为式：`inventory-settles-on-the-first-pair`・`inventory-settles-on-the-second-pair`・`inventory-never-settles`・`inventory-single-read-never-settles`・`inventory-retry-bound-leaves-room-to-settle`；静态：`round-preserves-every-inventory-read`・`round-records-the-settled-pair`・`round-records-how-many-reads-were-needed`・`round-preserves-the-unsettled-pair` |
| R34 | 单次盘点读取失败（重启后 guest 管理通道尚未就绪：实测返回很快、非超时）⇒ 旧代码直接 `exit 5`，整轮作废并触发整轮重试（每轮~5 min + 一次新 cut）；且日志只写 'failed'，**不记退出码与 stderr**（与 P4 同类可观测性缺陷） | 单次读取失败改为**轮内重试**（消耗同一个有界读取预算，不重跑整轮）；失败详情（exitCode + bounded stderr/stdout）逐次入日志；**全部**读取均失败才 `inventory-failures.txt` + `exit 5` | `round-retries-a-failed-inventory-read-in-place`・`inventory-failure-log-carries-the-exit-code`・`round-aborts-when-every-inventory-read-fails` |
| R35 | 变异臂（D1）文件命名 `agent_runner_b3b_variant_arm_test.go` 的**尾段 `arm` 被 Go 当作隐式 `GOARCH=arm`** ⇒ 该文件在 windows/amd64 上**永不参与构建**（显式 `//go:build windows && b3bnative && b3bvariant` 也无效），症状是 gopls 报 `No packages found for open file`（易误判为 buildFlags 未配） | 改名 `agent_runner_b3b_variant_degraded_test.go`（尾段避开 GOOS/GOARCH 词）；用 `go list -tags … -f "{{.TestGoFiles}}|{{.IgnoredGoFiles}}"` 核验归属（变体视图已含该文件、默认文件被排除） | 无需新断言：接线由 `go list` 两视图核验；坑位记入仓库记忆 `prfrail-go-gotchas.md` |

**⑤ 实机运行中才发现并修的七项**（逐级复审：P1–P4 过 ③ V4 Pro 收口；P5・P7 过 ④ Codex 定向复审 `VERDICT: PASS`；P6 见下方反证实验）：

| # | 现象（实机复现） | 整改 | 钉住它的断言 |
| --- | --- | --- | --- |
| P1 | 硬断电后残留 `VBoxHeadless` 会话使紧随的 `startvm` 报 `E_FAIL / SessionMachine` | `Wait-B3bVmPowerOff`（等 VM 真到 `poweroff`）+ `Clear-B3bStaleVmSession`（词边界匹配与 `$safeToClear` 边界，见 R27）+ `Start-B3bVm` 有界重试并记录 stderr；S1/S2/S8 三处生效 | `stale-session-match-is-boundary-anchored`・实机 `startvm ok=True attempts=1` |
| P2 | 装置无并发保护（两个实例同时操作同一 VM ⇒ `BLOCKER: cannot query VM`） | `Enter-B3bSessionLock`/`Exit-B3bSessionLock`（记录 pid + owner + **会话目录指纹**；按进程存活判陈旧；矩阵持锁不释放、预检短持后释放） | `lock-blocks-a-second-live-holder`・`inherited-lock-does-not-block-a-copied-session` |
| P3 | 冻结的 `Write-B3aTextFile -Text` 用参数绑定拒绝空串，而多处工件**设计上可为空** ⇒ 断电后未捕获异常 `exit 1` | 新增 `Write-B3bArtifactText`（`[AllowEmptyString]`，UTF-8 无 BOM + LF；空串不改变任何参与哈希/完整性门的字节）并替换轮驱动中所有可能为空的写入点 | 实机单轮 `exit 0` + `writer-stderr.txt`（0 字节且存在） |
| P4 | 矩阵在非预期退出码时丢弃轮 stderr ⇒ 只能看到裸退出码 | 非 0/3/5/6/7 时记录轮 stderr（该补丁直接定位了 P3） | 无（可观测性） |
| P5 | 正对照标定命令被 `REFUSED`（见 R31） | 共享 cell 解析器 + 行为式断言；五个 cell 实跑解析全部 OK | `control-arm-cell-is-schedulable` 等三条；反证实验：回退旧正则 ⇒ 该断言变红并给出 `arm=<threw>` 诊断 |
| P6 | 跨零点后自检 7 条 FAIL、finalize 写出空判定（见 R32） | 会话号显式贯通 + finalize 零轮 fail-closed | `finalize-refuses-a-session-with-no-rounds`・`finalize-honours-the-requested-session-id`；反证实验：禁用守卫 ⇒ 前者变红（exit=0 且指向空会话的 `matrix.json`） |
| P8 | c1 某轮 `inventory run 1 failed`，同一位置连续两次 ⇒ 整轮作废 3 次后单元格停止（见 R34） | 盘点读取失败改为轮内重试 + 退出码/stderr 入日志 | `round-retries-a-failed-inventory-read-in-place` 等三条 |

**⑤ 检查点结果（2026-09-21 02:55，会话 `2026-09-21`）**：`-Checkpoint` = **`CHECKPOINT PASSED: discrimination power demonstrated; the matrix may proceed.`**（exit 0）。
依据：正对照（`positive-control:parent-synced`）r00–r04 共 5 轮全部 `survived`（丢失 0）；负对照（`c1:temp-written`）有效轮均 `lost-absent` ⇒ 判别力门双半成立。
同会话盘点读取次数实测分布：`reads=2`（5 轮标定 + 2 轮 c1）与 `reads=3`（4 轮 c1），**未出现 4 读**；即 R33 的额外代价平均为 1 次盘点（≈ 27 s/轮）。
矩阵运行失败计入 `exit 5`（r07，2 次整轮重试同点失败）已在 R34 修复后清零（恢复补跑 3 轮全部 `exit 0`）。

**⑤ 用户裁决记录（2026-09-21）**

| # | 裁决 | 执行口径 |
| --- | --- | --- |
| D1 | **补跑变异臂 2 轮**（§6 表中"变异臂（装置自证伪）"） | 定位为**加强证据**：不阻塞主矩阵、不改变主矩阵结果；窗口 = 主矩阵完成后、`-Finalize` 之前；**若失败不回溯主矩阵**，只如实记录，并把 §10.1 结论降级为"判别力门成立但变异臂未通过"。实现口径：tag 门控变体构建（不污染默认构建）+ 独立 rig root（避免污染主会话判定）+ 事后并入 `<主会话>\variant\`（bundle 会复制该可选分区）。**判据**：降级原语下装置若仍给 `proven` ⇒ 装置失明（属严重发现，按降级登记）；若封顶 `unresolved`/不 proven ⇒ 自证伪成立。 |
| D2 | **④ 按实记账** | 报告「成本与计量」逐次列出 5 次 ④ 调用的目的与发现数，并写明超预算原因（闭环铁律 + P5/P7 新增审查点），明确"非执行者违规"；同时在 §10.1 或报告末尾登记一条**「准则待议」**：硬门切片（A7/B2/B3/B4）的 ④ 调用上限是否需单列——留待本轮结束后评估，**不在 A7 内修改准则**。 |
| D3 | **⑥ 档位切深度（Max），可提前切** | ⑥ 的 ADR-013 重排与反证原型类**证明式推理**需要 Max；切换只影响主控推理档位，**不影响**后台运行的矩阵进程（独立于会话）。用户可能不在场 ⇒ 提前切换；若切换导致会话上下文重建，状态可由 `/memories/session/b3b-t027-closeout.md` + 本设计说明 + 会话日志完全恢复。 |

**R33 原始证据（实机，主会话 `2026-09-20`，c1:temp-written r08，两次盘点的完整差异）**：
`len1=720 len2=723`；叶子 18/18；**差异叶子恰 2 个**，且只指向同一字段——
`probeEntries.files[0].size`：首读 **0** → 次读 **1082**（文件为未链接临时件 `requests\.request.b3b-r08.jsonl.819624024.tmp`）。
两读的判定相关字段**完全一致**：`probeState=entries`、`dirty='Volume - R: is NOT Dirty'`、`lastBoot=2026-09-21 00:54:24`、
`fileCount=4`、`tmpFiles` 列表、`requestCount=0`、`completionCount=0`、`requestsWithoutCompletion=[]`、`orphanCompletions=[]`、`runDirs=[]`。
即：差异集中在**不参与 c1 结论判定**的未链接临时件尺寸上，而 c1 的预期结论（记录未被发布 ⇒ `lost`）在两次读取下都成立。
取证前置事实：旧代码在比较失败时先 void、后落盘，三轮诚实轮次**没留下任何可比对字节**——这是可观测性缺陷，已一并修（`inventory-mismatch-1/2.bin`）。
反证实验：把"向前扫描首个相邻一致对"退回"只看首对" ⇒ `inventory-settles-on-the-second-pair` 变红（`stable=False tries=3 pairOk=False`），恢复后回到 `checks=100 failures=0`。
④ Codex 定向复审（带上述原始差异）判定：**PASS + GO**，并明确"现有证据更支持测量伪影而非掩盖设备保真度问题"。

**时间口径（R33 附带，须按实测回算）**：单轮时长不再能按"两次盘点"估算；以每轮 `inventory-reads.txt` 的分布乘以单次盘点实测耗时回算，
不得沿用 §6 旧的"两次盘点"口径。S8 单次盘点调用超时上限 120 s，理论上界因此抬高。

**设备保真度：本轮证据未证伪、亦未证明（R33 附带，⑥ 必须逐句限定）**：本次差异方向是"晚读可见到更多"（0 → 1082）。
两种解释都能产生该现象：① 客机首次枚举读到了尚未定型的元数据视图（测量伪影）；② 宿主/虚拟盘层在 cut 之后仍把已缓存写落盘（保真度问题）。
证据**偏向前者**：判定关键字段两读一致、且该尺寸变化不改变 c1 的 `lost` 结论；但装置本身**不能仅凭这一条差异**证明排除后者。
因此 ⑥ 不得据此声称"装置已证明等价于物理断电"；如需该主张，必须另有独立检验（例如对 cut 前后虚拟盘上的字节做外部比对），本片未做。

**会话锁恢复动作（R28）**：锁的陈旧性以「记录的 pid 是否仍存活」+「锁内 `sessionDir` 是否指向当前会话目录」判定。
已知盲区：上一轮崩溃后若 OS 将该 pid 复用给无关进程，下一轮诚实运行会被 `REFUSED: the session lock is held by a live process`。
恢复动作：确认确无装置进程在跑（`Get-CimInstance Win32_Process -Filter "Name = 'powershell.exe'"` 过滤 `Invoke-B3b`）后，删除该会话目录下的 `session.lock` 重跑；
不得在不确定时删除（否则会重新引入两实例并发操作同一 VM 的危险局面）。

**收口闭环记录（用户质询后补做，2026-09-20）**：按《Flash 主控执行准则》v3.1 §5 铁律，③/③.5/④ 三级在整改后**各自重跑收口复审**：

| 级 | 收口结果 | 备注 |
| --- | --- | --- |
| ③ V4 Pro | **PASS** | 2 轮内无新实质发现 |
| ③.5 MAI | **PASS** | — |
| ④ Codex 收口 #1 | `FINDINGS` | 5/5 历史项关闭；1 条新 MAJOR = R29；1 处覆盖缺口 = R30 |
| ④ Codex 收口 #2 | **PASS（明确 GO：103 轮可起跑）** | R29/R30 均按可证伪性逐条给出\"删守卫则哪条断言变红\" |

治理更正：此前\"整改后未重跑本级复审\"属**流程缺陷**，不是特殊考虑；本表为更正后的执行事实，成本与计量一并在切片报告记账。

**复审终止规则（防无限循环，⑥ 须采纳）**：仅当新一轮发现**改变被断言的行为**或**暴露新的可证伪性缺口**时才再跑同级复审；仅影响诊断信息、不改变任何断言与判定语义的细节记为**已接受残余风险**并入档，不再消耗复审预算。

**已接受的残余风险（R29 附带）**：`trap` 体内层未再包 try/catch。若释放时遇极端 I/O/ACL 异常，原始异常文本可能被替代——影响面仅限**诊断文本**：退出码仍为 1；`Exit-B3bSessionLock` 本身已含\"文件不存在即返回\"\"JSON 解析失败即返回\"\"仅当 pid 匹配才删除\"三层防御；且该情形下锁文件无论如何都删不掉，遮蔽不改变结果。判定不构成新缺陷，记录在案、不追修。

**报告不得使用的措辞**（预审第 5 项，⑥ 必须逐条规避）：不得出现长期保证/普遍耐久（`durable`/`guarantee`/`survives power loss`）——只能用 §10.2 句式；不得写 U4"已证明有序"（18 轮只能证伪，上限仍是 `unresolved`）；不得写"候选在决定格优于基线"（除非 C1 同阶段格出现丢失）；不得声称"journal 链可检出轮工件篡改"（链只验 journal 行，工件完整性靠 R3 的重派生 + 逐文件比对；对工件与 journal 的**整体重写**只能靠证据包校验和与 S0 pin 这两个外部锚点发现）；不得声称"证据包自动可复算"（只能说可由 `rounds/` + `journal.jsonl` 复算）。
**必须逐句声明的限定**（MAI 追加）：(a) 正对照实际执行的是卷句柄 flush 还是回退的记录句柄 flush——若为回退，区分力门只能读作"在该回退刷盘下装置可观测到干净结果"，**不得**读作"卷级设备缓存已证明被排空"；(b) `c2c3` 在 `parent-synced` 与 `c3` 同一决定点、且该阶段的组合原语与 C3 等价，故 `c2c3` 只能作为**诊断臂/重复臂**报告，不得当作独立的耐久主张；(c) 阶段证明取自 cut 前转录，而非断电后拷回件。
