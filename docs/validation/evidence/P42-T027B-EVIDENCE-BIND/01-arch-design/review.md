# P42-T027B-EVIDENCE-BIND ① 架构设计书（C3+N2）

切片：`P42-T027B-EVIDENCE-BIND`。环节：① 架构。载体：`deep-reasoner`；模型：`DeepSeek V4 Pro (deepseek)`。输出原样归档，内容未作改写。

---

# P42-T027B-EVIDENCE-BIND · ① 架构《C3+N2 句柄与目录绑定方案》

只读完成：未修改任何文件、无 git 操作。已核对 C3/N2 判定（P38 review.md）、`internal/evidence` 全目录、`internal/adapters` 全部读点、CONTRACTS §4/§7、P42 16 字段块。

## 一、现状与断裂点（file:line）

| 点位 | 断裂 |
|---|---|
| `internal/adapters/agent_runner_evidence.go:378-398` `readAgentRunnerPIDReport` | `os.Stat(path)`→限长判→`os.ReadFile(path)`：检查与使用两次解析路径（C3 核心点） |
| `internal/adapters/agent_runner_evidence.go:411-436` `readBoundedAgentRunnerArtifact` | 已 `os.Open`+`file.Stat()`（句柄上），但 `info.Size()<=bound` 分支**又按路径 `os.ReadFile(path)`** 重开——检查在句柄、读取在路径 |
| `internal/adapters/agent_runner_postflight.go:579/593` `agentRunnerPostflightReadArtifact`/`agentRunnerFileHash` | 经上述原语，同源 |
| `internal/adapters/agent_runner_postflight.go:552` `agentRunnerReadPIDs` | 无界 `os.ReadFile`，无身份/限长 |
| `internal/adapters/agent_runner_postflight.go:440` `scanWorkspace` | 每条目 `os.ReadFile(workspaceRoot/entry.Path)` |
| `internal/adapters/agent_runner_replay_store.go:571/589` `loadRequestRecord`/`loadCompletionRecord` | 无界 `os.ReadFile` |
| `internal/adapters/agent_runner_replay_root.go:381-399` `loadAgentRunnerReplayOwnership` | `os.Lstat(path)`→`os.ReadFile(path)` |
| `internal/adapters/filequeue.go:528-540/545-556` `loadExistingRequest`/`ensureInflightRequest` | `os.Stat` 探测→`loadRequest`（`filequeue.go:617/626/633` 内 `os.ReadFile`） |
| `internal/adapters/ai_probe_copilot_executor.go:85` | 无界 `os.ReadFile(usagePath)` |
| `internal/adapters/agent_runner_replay_root.go:289-311` `RunEvidenceDir` | `verifyPathSafetyLocked`（按路径 Lstat/GetFileAttributes 爬组件）→mkdir→`os.Stat`→**返回字符串**，绑定在返回时终止（N2 核心点） |
| `internal/evidence/disposition.go:247/256-281` `BuildReleaseInventory` | `os.Stat(rootAbs)`→`filepath.WalkDir`→每文件 `os.ReadFile(path)`：stat-后-walk + 按路径读 |

CONTRACTS 依据：A2 已显式登记「纯用户态组件检查与文件操作之间存在无法完全消除的 TOCTOU 残余窗口，属本切片显式记录的已知边界」；P42 即把该残余中「数据读」一类收敛到句柄语义。postflight A5 要求 pid 报告腿「重读重算对账」——本片必须保持该语义（重读仍重读，但读绑定同一句柄）。

## 二、Q1/Q2/Q3

**Q1 统一读取原语**（新增 `internal/evidence/boundread.go`，仅增不改）：

```go
type BoundFile struct{ f *os.File; info os.FileInfo }
func OpenBoundFile(path string, bound int64) (*BoundFile, error)  // 唯一 open
func (b *BoundFile) ReadAll() ([]byte, bool, error)              // truncated 标志
func (b *BoundFile) Hash(domain string) (string, error)          // 同一句柄
func (b *BoundFile) Identity() (size int64, mod time.Time, mode os.FileMode)
```

组合方式：`os.Open` 只开一次；身份取自 `f.Stat()`（句柄，非路径）；`ReadAll = Seek(0,0) + io.ReadAll(io.LimitReader(f, bound+1))`，`len==bound+1 ⇒ truncated=true`；`Hash = Seek(0,0) + io.Copy(io.MultiWriter(sha256, io.Discard), io.LimitReader(f, bound+1))` 单遍读+摘要；打开时 `f.Stat()` 判目录/非 regular/`Size()>bound` 即 fail-closed（日志类走 truncated 语义，见下）。**禁则：打开后任何 `os.Stat`/`os.Lstat`/`os.ReadFile`/`filepath.Join` 重新按路径解析。**

逐处替换清单：
- `agent_runner_evidence.go:378`：改 `OpenBoundFile(path, agentRunnerPostflightArtifactBound)`；`ErrNotExist→(nil,false,nil)` 保持「缺失非错误」契约；目录/超限→错误。
- `agent_runner_evidence.go:431`：**删除 `os.ReadFile(path)` 分支**，全走句柄 `ReadAll`。
- `agent_runner_postflight.go:552`：`agentRunnerReadPIDs` 改 `OpenBoundFile(path, agentRunnerPostflightArtifactBound).ReadAll`。
- `agent_runner_postflight.go:440`：`scanWorkspace` 逐条目 `OpenBoundFile`（目录绑定见 Q2）。
- `agent_runner_replay_store.go:571/589`、`agent_runner_replay_root.go:381`：三类记录载入改句柄读，类型判据从 `f.Stat()` 取（消除 `Lstat`+`ReadFile` 双解析）。
- `filequeue.go:528-556`：删除 `os.Stat` 预探测，`loadRequest` 单次 `OpenBoundFile` 内做 `ErrNotExist→(zero,false,nil)`；`ensureInflightRequest` 碰撞重试保留，但「存在性」只由 load 的 open 结果判定。
- `ai_probe_copilot_executor.go:85`：`OpenBoundFile` 小界（1 MiB），保持 `usage,_` 忽略错误语义。
- `evidence/disposition.go:247-281`：`BuildReleaseInventory` 改 Q2 目录绑定 + 句柄读（见下）。

失败与部分读语义：open 失败（除 ErrNotExist）→错误、零数据；打开即超限→**拒绝**（PID/记录类）或 **truncated 标志**（日志/事件类，保持既有降级口径）；`io.ReadAll` 中途 I/O 错误→丢弃全部字节、不产摘要；**摘要只在整读成功时计算**；truncated 数据永不进入六事实 `FrozenFacts`（既有规则不变）。

Windows/Unix 句柄差异：两者已打开的文件句柄都**绑定文件对象/inode**——路径随后被改名/替换/换 symlink 都不影响后续读，这正是「同句柄读到旧内容」的机制；Windows Go `os.Open` 携带 `FILE_SHARE_READ|WRITE|DELETE`，被替换/删除中的文件仍可读；Unix `open(2)` 后 unlink 不中断。注意点：① 打开**前**路径已是 symlink 属静态条件，由既有组件级拒绝覆盖，本片不重复；② Windows 若需事后复核身份可用 `GetFinalPathNameByHandle`，但正确性不依赖它——句柄即身份；③ 目录句柄差异见 Q2。

**Q2 祖先目录替换防护**。机制：
- Windows：`syscall.CreateFile(dir, GENERIC_READ, FILE_SHARE_READ|WRITE|DELETE, nil, OPEN_EXISTING, FILE_FLAG_BACKUP_SEMANTICS, 0)` 得目录句柄；相对打开用 `GetFinalPathNameByHandle(dirHandle)` 取绑定目录**对象当前真名** + `filepath.Join(name, rel)` 再 `CreateFile`（无 `openat`，用句柄派生名实现「绑定跟随对象而非名字」）。
- Unix（`//go:build !windows`）：`syscall.Openat(dirfd, rel, O_RDONLY|O_NOFOLLOW|O_CLOEXEC, 0)`；目录自身用 `fstat` 判 `S_IFDIR`。`openat` 使祖先被 rename/替换无法重定向解析；`O_NOFOLLOW` 防叶节点被换成 symlink。

Go 1.22 无 `os.Root`（Go 1.24+），走 syscall 方案；go.mod 即 `go 1.22`。平台边界：目录句柄打开、相对打开各放 `bounddir_windows.go`/`bounddir_unix.go`，各自可测（Windows 测试覆盖 junction/rename 替换；Unix 测试覆盖 rename+O_NOFOLLOW）；公共不变式测试跨平台同断言。

`RunEvidenceDir` 升级：公开签名**不变**（Q3），内部在 `ensureAgentRunnerReplayDirectory` 确认后**立即打开目录句柄**并经句柄复核（`fstat` 是目录、Windows 上 `GetFileInformationByHandle` 属性含 DIRECTORY、`GetFinalPathNameByHandle` 与期望 clean 路径一致）——封住「mkdir 确认→返回」之间的换目录窗口；句柄注册进 `AgentRunnerReplayStore` 新增 `boundDirs map[string]*os.File`（`store.mu` 保护，requestID 为键）。adapters 内后续读取全部经新方法 `readBoundEvidence(requestID, rel, bound)` 走 `Openat`/`CreateFile` 相对打开。真源与生命周期：**持有者 = AgentRunnerReplayStore**（构造时绑定、run 全程有效）；新增 `func (store *AgentRunnerReplayStore) CloseBoundDirs()`，由 run 收尾路径调用，store 无显式 Close 则以进程生命周期为上限并登记；目录句柄以 share-delete 打开，不阻塞候选/工具继续写目录。

**Q3 契约兼容**。必须保持的签名：
1. `func (store *AgentRunnerReplayStore) RunEvidenceDir(requestID string) (string, error)`——`internal/console/runtime_agent.go:594` 与 `runtime_agent_e2e_test.go` 依赖字符串返回（console 超范围，不得破坏）。
2. `CollectAgentRunnerEvidence(ctx, AgentRunnerEvidenceSpec)` 及 spec 字段——run/postflight/console 装配三方调用。
3. 三个 seam 函数变量 `replayStoreLoadRequestRecord`/`replayStoreLoadCompletionRecord`/`replayStoreLoadOwnershipRecord`（`agent_runner_replay_store_test.go:350-384`、`agent_runner_replay_root_unix_test.go:202` 会打补丁）——**函数签名逐字不变**，只改其默认实现体。

演进方式：`internal/evidence` 只新增导出原语（零既有符号变更），adapters 内部换实现体；`internal/chain`/`gates`/`snapshot`/`console` 零改动。与 P41 分工：P41 的「重开+重锁+重解码」是**锁语义**（主机存活性判据：byte-range/flock 锁 + rename 保锁），与本片 TOCTOU 正交；`decodeAgentRunnerStartingMarker` 已是「打开文件上 Seek+LimitReader」的句柄式读取，本片原语可被其复用但不改其锁序；P41 拥有的锁/租约判定一条不动。

## 三、签名与不变式表

关键签名（≤3）：
```go
// internal/evidence
func OpenBoundFile(path string, bound int64) (*BoundFile, error)
func (b *BoundFile) ReadAll() ([]byte, bool, error)
// internal/adapters（+ 平台文件 OpenBoundDir/OpenBoundDirFile）
func (store *AgentRunnerReplayStore) readBoundEvidence(requestID, rel string, bound int64) ([]byte, bool, error)
```

| # | 不变式 | 位置/持有 |
|---|---|---|
| 1 | 每条证据读仅一次 open；限长、内容、摘要、身份同出自一个 `*os.File` | `BoundFile` |
| 2 | 打开即超限 ⇒ 拒绝（记录/PID 类）或 truncated 标志（日志类）；truncated 永不计入六事实 | `readBoundedAgentRunnerArtifact` |
| 3 | 摘要仅在整读成功时计算；部分读丢弃、错误不产摘要 | `BoundFile.Hash` |
| 4 | `RunEvidenceDir` 返回前目录已被句柄绑定；adapters 内后续 I/O 全部经 dirfd/目录句柄相对打开 | `AgentRunnerReplayStore.boundDirs` |
| 5 | 目录句柄生命周期 = store 生命周期；未 Close 前不得关闭 | `CloseBoundDirs` |
| 6 | 公开签名（`RunEvidenceDir` string 返回、三个 seam 变量）零变化 | 兼容性 |

## 四、红测不变式（③ 必须钉住，3–5 条）

- **I1（同句柄读旧内容）**：`OpenBoundFile` 后、`ReadAll` 前，将路径替换为指向敌意文件（Unix：rename+新文件；Windows：rename），`ReadAll` 必须返回打开时 inode 的字节且 `Hash` 与替换前一致。判定：读到替换后内容 ⇒ 红。
- **I2（替换攻击下拒绝服务）**：读前把文件替换为超 bound 文件或目录，读取必须报超限/类型错误、`FrozenFacts` 失败，绝不返回攻击者数据。判定：返回数据或 digest ⇒ 红。
- **I3（祖先目录替换）**：`RunEvidenceDir` 绑定后，把 `agent-runner-replay/` 或 evidence 目录的祖先 rename 走、原地放同名 symlink/junction 指向注入目录；随后 `readBoundEvidence` 必须仍读到绑定目录内原文件、缺失文件按原目录 ErrNotExist。判定：注入文件内容出现 ⇒ 红。
- **I4（身份自洽）**：`Identity().size` 与 `ReadAll` 字节长度一致（或与 truncated 标志一致），`Hash` 与未受攻击时 `os.ReadFile` 对照一致。
- **I5（等价回归）**：既有 postflight「重读重算对账」腿、seam 补丁测试、Windows/Unix marker 与回放测试仅换实现后全绿（等价性，不新增断言）。

## 五、风险与反例

- **反例甲（推翻本方案）**：任何实现改回「先 stat 路径、再另开读取」或 `ReadAll` 内部重新 `filepath` 解析——I1/I2 立即红。机械门禁：`adapters`/`evidence` 生产代码 grep 禁 `os.ReadFile`/`os.Stat`（测试与 fixture 除外）。
- **反例乙（Windows 残余窗口）**：`GetFinalPathNameByHandle` 取名与 `CreateFile` 之间目录再被换 ⇒ 打开到的文件用句柄复核「最终路径仍在绑定目录之下」，不符即拒绝。该残余与 CONTRACTS 既有「纯用户态检查残余 TOCTOU」口径一致，登记。
- **反例丙（Unix O_NOFOLLOW 仅覆盖叶节点）**：中间组件 symlink 由 `openat` 天然不越出（相对 dirfd 解析，组件替换只在目录内生效）；目录本身被换不影响 dirfd。
- **风险一（句柄泄漏）**：绑定目录句柄未收敛关闭 → 随进程存活；share-delete 已含删除共享，但仍须 `CloseBoundDirs` 并加「关闭后可删目录」测试。
- **风险二（超范围调用方）**：console 仍以 `RunEvidenceDir` 字符串直接 open（`runtime_agent.go:633`），本片绑定不覆盖；报告如实登记，迁移属后续切片。
- **风险三（语义变化）**：打开时 `Size()>bound` 由「读满 bound 报 truncated」变为「立即拒绝」——对日志类必须保留 truncated 分支（读 bound+1），否则行为漂移；红测 I2 与既有日志超限测试并跑。

## 六、与 P43/P41 的边界

- **P41（已闭环 `caaf64d`）**：starting marker 锁序、重开重锁重解码、rename-away 释放全部不碰；本片原语可复用但零改动其语义。
- **P43 及后续**：console 包（`runtime_agent.go:594/633`）、`gates`/`snapshot` 的读取迁移、以及「读路径全面句柄化」的收口留给后续片；本片仅 `internal/adapters/**` 与 `internal/evidence/**`，不新增依赖、不触碰 chain。

DESIGN: DONE

## 修订注记（② 实现期与 ⑤ r1 后）

本节只追加实现期实测结论与口径修正，不改动上文原设计语义；上文表述与本节冲突时以本节为准。

- **D1（绑定对象收敛到 `runs` 基目录）**：Q2 原设计拟绑定 per-request 证据目录，实现期改为绑定 `runs` 基目录。实测依据：Windows 上子目录存在打开句柄时，重命名其祖先返回 `ERROR_ACCESS_DENIED`（即便以 `FILE_SHARE_DELETE` 打开），per-request 句柄会把「祖先替换攻击」变成拒绝服务；绑定基目录则让祖先可被改名而句柄继续指向原对象树，正是设计所要求的性质。读取解析随之定为 `<requestId>/<rel>`（相对基目录句柄解析）。残余面如实登记：`runs/<requestId>` 自身被同层替换（兄弟节点，非祖先）不在基目录句柄覆盖内；唯一生产调用方以冻结摘要比对（postflight `ManifestHash`/`DiffHash`）对读到的内容 fail-closed，替换退化为失败而非可信证据。per-request 子句柄列入后续片评估项。
- **D4（`BuildReleaseInventory` 不迁移）**：`internal/evidence/disposition.go` 的 `BuildReleaseInventory` 不在本片迁移。理由：adapters→evidence 的依赖方向不允许反向复用（本片口径是 evidence 新增原语、adapters 使用；不接受 evidence 反向依赖 adapters 的目录绑定能力）；且 release 树由操作员选定，不属候选可写面。因此该项属**设计清单项收窄**，归入后续读路径 sweep 片。
- **F1（Windows 逐读复核已实现）**：反例乙的修法已落地——`CreateFile` 打开产物后，对**文件句柄**调用 `GetFinalPathNameByHandle` 取最终路径，断言其（clean 后、剥 `\\?\` 前缀）仍位于绑定目录句柄**当前真名**之下；不符即关闭句柄并返回拒绝错误（fail-closed），封住「派生名 → `CreateFile`」之间的再替换窗口。Unix 侧无需后置复核：`openat(2)` 相对 dirfd 的解析与打开在同一系统调用内原子完成，不存在同名窗口；代码中以注释说明该平台差异。
