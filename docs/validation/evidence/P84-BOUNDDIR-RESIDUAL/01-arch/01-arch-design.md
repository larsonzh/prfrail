# P84 BOUNDDIR-RESIDUAL ① 修复方案书（01-arch-design）

- 载体：① 架构师（deep-reasoner，DeepSeek V4 Pro Max；用户 2026-11-21 排片裁定）
- 落盘：主控（2026-11-21）；原文逐字转录自 ① 的最终回复。
- 结论：T4 = DOES-NOT-TOUCH / IN-BOUNDARY；T5#4 = DOES-NOT-TOUCH / IN-BOUNDARY。

---

P84-① ARCHITECTURE DESIGN

## ① 契约影响判定（§1.4 协议先行）

| 项 | 判定 | 理由 |
|---|---|---|
| T4（F1 残余） | **DOES-NOT-TOUCH / IN-BOUNDARY** | 被引用句 `SECURITY.md:23`（「限制必须由 runner/OS 能力落实」）只定要求不定机制；句柄相对打开、Ordinal 比较、中间 reparse 拒绝均是 `internal/evidence` 对该句的实现强化。与 P83 ① ruling-triage T4 判定一致。无新 schema、新 CLI；`ntdll!NtCreateFile` 走 `syscall.NewLazyDLL`，与既有 `kernel32!GetFinalPathNameByHandleW`（`bounddir_windows.go:50`）同型 OS API 调用，不是第三方包，go.mod 不动 ⇒ 不触守则 OUT-OF-BOUNDARY 三条件。 |
| T5#4（F4 残余） | **DOES-NOT-TOUCH / IN-BOUNDARY** | 被引用句 `CONTRACTS.md:91` 与 `:150` 已由 store + proof 架构落实；本项只是让 store-local 投影可恢复。`publishLeaseProof` 注释自述「not a wire record」（`runtime_agent.go:1886-1888`），不入 recordHash 集，租约文件属 `internal/guard` store-local 布局——与 P83 ① ruling-triage T5 行理由逐字同源。retained 代次文件（`lease.go:310`、`:328`）已存在，仅加内部链式校验 helper；无新 schema/CLI/依赖 ⇒ IN-BOUNDARY。 |

两项目前均无需 CONTRACTS{,_EN}/schema/夹具先行修改；若实现中发现必须改 schema/CLI/依赖，停下升级。

## ② T4 设计（逐改动点）

**(a) Windows 句柄相对后代打开（机制选型）**
选型 = **逐组件 `NtCreateFile` 走链**（RootDirectory=已验证父句柄），拒绝另两案：
- 单次 NtCreateFile/CreateFile 带 `FILE_OPEN_REPARSE_POINT` 开整条相对名：只作用于末级，中间 reparse 仍被跟随 ⇒ 不满足 (c)；
- `GetFinalPathNameByHandleW` 派生名 + 重开 + 比较：名与开之间可换根，正是 P42 注释（`bounddir_windows.go:10-16`）与 ⑦ 复审「不能排除中途换根」点名的形状。
走链论证：每一步以上一级已验证目录句柄为 RootDirectory 开下一组件 ⇒ 内核只在该目录对象内解析，case 变体孪生目录是绑定根的兄弟、不可达 ⇒ 无需路径比较即证归属。ABI：`ntdll!NtCreateFile`，`OBJECT_ATTRIBUTES{Length=unsafe.Sizeof, RootDirectory=父句柄, ObjectName=UNICODE_STRING(组件名), Attributes=OBJ_CASE_INSENSITIVE(0x40)}`，`ShareAccess=READ|WRITE|DELETE`，`CreateDisposition=FILE_OPEN(1)`，`CreateOptions=FILE_OPEN_REPARSE_POINT(0x200000)|FILE_OPEN_FOR_BACKUP_INTENT(0x4000)`，目录组件加 `FILE_DIRECTORY_FILE(1)`、末级加 `FILE_NON_DIRECTORY_FILE(0x40)`；NTSTATUS 映射后包 `*os.PathError`（保 `errors.Is(fs.ErrNotExist)`，同 `bounddir_unix.go:143-147` 口径）。仅 `//go:build windows`，Go 1.22 + 现 `syscall`/`unsafe` 可编译，无新依赖。
改动点：
- `bounddir_windows.go:86-121`（`openBoundDirFile`）：删 `:98` `filepath.Join(base, clean)` 与 `:101-102` 路径 `CreateFile`，改为 `openDescendantRelative(dir.handle, clean, wantDir=false)`；
- `bounddir_windows.go:135-190`（`BoundDirNames`）：`clean != "."` 分支的 `:153-166` 路径打开改走同一 `openDescendantRelative(..., wantDir=true)`；根分支（句柄直读）不变；
- 新增私有 helper（同文件）逐组件循环：开 → `GetFileInformationByHandle` → `FILE_ATTRIBUTE_REPARSE_POINT` 即关+拒 → 下一组件；
- `:117` 与 `:198` 的 `confirmOpenedArtifact`/`confirmEnumeratedDirectory` 保留（卷号 + 名称 belt）。

**(b) Ordinal 替换 EqualFold**
`bounddir_windows.go:290-300`（`isUnderDirectory`）：`:297` 的 `strings.EqualFold` 前缀比较改为 **Ordinal**（字节相等）。语义：`base`/`candidate` 均来自 `handleFinalPath`（`:361-382`，同 API 同形态、真实磁盘大小写），真后代的最终名前缀字节必然相等 ⇒ Ordinal 不误拒；case 变体孪生字节不同 ⇒ 拒绝。卷号 belt 已在 `confirmOpenedArtifact`（`:235-288` 内 `VolumeSerialNumber` 比较）保留。连带：`boundDirVerify` 的 `:330` `EqualFold` 快路径同样可放行 case 孪生 ⇒ 改为仅字节相等才走快路径，其余一律 `sameObjectAs`。

**(c) 中间 reparse 拒绝**：由 (a) 走链的逐组件 `FILE_OPEN_REPARSE_POINT` + 属性检查落实；末级 reparse 拒绝（现 `:101` 标志语义）随之并入走链，不丢失 P83 已闭合的 T4b。

**(d) 根 rel 语义一致**：`bounddir_clean.go` 新增共享分类器 `boundRootName(rel) = strings.TrimSpace(rel)=="" || filepath.Clean(filepath.FromSlash(rel))=="."`。两侧统一口径：`BoundDirNames`（win `:135-147`、unix `:254-262`、darwin `:62-66`）用分类器把 `""` 与 `"."` 都映射为根；`openBoundDirFile`（win `:86-95`、unix `:129-140`）对根统一以同一错误拒绝。`cleanBoundRelativePath` 本体不动（`bounddir_clean.go:27-31` 对 `"."` 的拒绝保留为产物名守卫，`bounddir_identity_test.go:66` 不受影响）；三处注释（win `:135-137`、unix `:254-258`）改为与行为逐字一致。

**(e) 无产物不打开根**：**现状已短路**——`runner.go:76-86`（仅 globs>0 才 bind）、`:271-274`、`:288-290`（collect 两条 no-glob 早退）。无需改码，只需钉住（④ W4，现有 `gates/p83_fixround5_tester_test.go:37-63` 已以不可绑定根钉住，其变异声明已写）。

**(f) Unix 侧**：`Readdirnames` 已落地（`bounddir_unix.go:286`）；openat2 掩码 `boundResolveMask` 已拒中间 symlink/跨卷。Unix 仅需 (d) 的分类器统一 + 注释对齐；无功能性改动。darwin 保持 fail-closed（`bounddir_darwin.go:42-50`、`:62-66`），零改动。

## ③ T5#4 设计（逐改动点）

**(a) 批量回写逐条前驱校验**
`persistRenewedPublishLease`（`runtime_supervised.go:1167-1208`）现仅 `:1194` 比较 `ResourceID`。改：
- 签名加 `ctx` + `store`（三调用点 `:1144`、`:1155`、`:1896` 均有 runtime，可用 `newPublishLeaseStore(runtime.runRoot)`）；
- 每项：`!found` → continue（现状保留）；`ResourceID` 不符 → continue（永不碰异 run 记录）；否则 `store.Current(ctx, refreshed.ResourceID)`（错则 return），再经新 helper `verifyPublishLeaseHistory(existing, current)`：同 claim、`Generation < current.Generation`、hash 合法、且沿 `current.PreviousHash → RetainedLease`（`lease.go:354`、`:362`）有界链走（上限 64）命中某 retained 记录 `R` 满足 `R.LeaseHash==existing.LeaseHash && R.Token==existing.Token`；
- 语义 = 报错而非跳过（与 `:1186-1190` 不可读同口径）：任一不满足 ⇒ 整批返回错误，错误文本含 proof 路径 + 失败检查名；该 proof 文件字节不改写。异资源/无文件两类才是 skip。
- 说明：链走同时覆盖 adoption 路径 `:1155`（Reclaim 后旧 proof 恰为 G-1 同 claim）与多代 skew（个别 task 连续多次未随批，仍属真实历史，可追平；伪造 claim/hash/异 run 永远不可达链）。

**(b) 代次→proof 可恢复事务**
线性化点（§13.1.11）：代次提交 = `writeLease` 原子 temp+rename（`lease.go:379`）持久成功即线性化点；proof 文件是投影（`runtime_agent.go:1886-1888`）。恢复不变量：任意已提交代次可由 {当前记录, retained 链} 重建投影。三条恢复路径并钉：① 同进程重入：`renewPublishLease` `:2087-2111` 前驱追平（前置 `verifyPublishLeasePredecessor`）；② 跨进程重入：`claimPublishLease` `:1945-1983` Reclaim（前任停机证明 `:1978-1981`）后 `:1155` 批量改写（经 (a) 链验证）；③ 批量中途崩溃：每项独立原子写（`writeSupervisedJSON` `:2389-2418` temp+rename+fsync），重入重跑链验证即收敛；retention 不清理（`lease.go:310-320` 注释）⇒ 历史链永存。改写幂等（重写同内容无副作用）。无新 schema/CLI/依赖。

**(c) 同进程分支现状钉住**：`verifyPublishLeasePredecessor`（`runtime_agent.go:2156-2182`）已含四道：claim（`:2157`）、N-1（`:2160`）、`PreviousHash`（`:2164`）、retained 原件 Token 复核（`:2168-2177`）；资源边界在 `:2075`/`:2184` 与 `retainedLease` `:362` 内复核。无缺口，保留现状，补测试映射（④ C4）。

**(d) 异主 proof 继续拒绝（显式声明）**：不改 `renewPublishLease` 的 `:2104`（当前记录非本进程 ⇒ 拒绝追平）与 `:2120`（需停机证明的 Reclaim）、`claimPublishLease` `:1961-1963`（资源/consumer 归属）、`verifyPublishLease` `:2202`（holder 复核）。批量改写仅限 `SupervisedDir(runRoot)` 下 + 资源匹配 + 真实历史三条件；不放宽任何接管条件。

## ④ 回归测试点（§13.1.2；反例 + 变异映射）

| # | 反例（平台） | 变异（移除哪步守卫 ⇒ 哪个测试变红） |
|---|---|---|
| W1 | `isUnderDirectory("C:\\Bound","C:\\BOUND\\x")` 必须 false；端到端 case 敏感目录孪生拒绝（Windows-only；fsutil 需权限则 skip+登记，VM/CI 补） | 把 Ordinal 还原为 `EqualFold`（`:297`）⇒ W1 红 |
| W2 | 中间 junction：`bound/mid→bound/real`，`OpenBoundDirFile(dir,"mid\\file.txt")` 与 `BoundDirNames(dir,"mid")` 均拒绝（Windows-only，本机可跑） | 删逐组件 reparse 检查或退回整名单次开 ⇒ W2 红 |
| W3 | 根一致性：`BoundDirNames(dir,"")` 与 `BoundDirNames(dir,".")` 均成功且等值（Windows 上绑 junction 根亦成功枚举目标）；`OpenBoundDirFile(dir,"")`/`(".")` 同错误拒绝（跨平台：本机 Windows + Linux CI） | 删 `boundRootName` 分类器 ⇒ `"."` 在 `BoundDirNames` 被拒 ⇒ W3 红 |
| W4 | 无产物不打开根：现 `gates/p83_fixround5_tester_test.go:37-63` 已钉（不可绑定根 + no-glob 仍 passed）。无新码，登记即可 | 删 `runner.go:76-86` 或 `:271-274`/`:288-290` 短路 ⇒ W4 红 |
| C1 | 同资源异 claim/哈希的其他任务 proof：批量改写必须报错且该文件字节不变（跨平台；按现实现 `:1194` 仅比 ResourceID ⇒ 预期现红） | 把 (a) 的链验证退回仅 ResourceID 比较 ⇒ C1 红 |
| C2 | 真实历史（同 claim N-1，含多代 skew）被追平成功、字节=refreshed（跨平台） | 删链走（仅认 N-1）⇒ C2 的多代 leg 红 |
| C3 | 代次提交后 proof 落盘前中断重入可恢复：同进程 leg 已被 `console/p83_fixround8_tester_test.go:56-83` 钉；跨进程 leg 用子进程自 exec（`-test.run` 标记）在 Windows 本机 + Linux CI 跑（跨平台；不稳定则降级 ⑨ VM/CI 验证并登记） | 删 `renewPublishLease:2087-2111` 追平或 `claimPublishLease` Reclaim ⇒ C3 红 |
| C4 | 同进程合形但不匹配 token/hash 仍拒绝：已钉（`console/p83_fixround8_tester_test.go:97-150`） | 删 `:2157`/`:2160`/`:2164`/`:2168-2177` 任一 ⇒ C4 对应子例红 |
| C5 | 异主拒绝不回退：现 `console/p83_fixround5_tester_test.go:432`（活跃前任）、`:419`（异主发布证明）保持绿 | 删 `runtime_agent.go:2104`/`:2120`/`:2202` holder 门 ⇒ C5 红 |

回归护栏：`evidence/p83_fixround8_windows_test.go:44/:130/:163`（T4a/T4b）、`evidence/p83_fixround8_linux_test.go:20`、`gates/p83_fixround5_tester_test.go:76`（根替换）、`guard/p83_fixround7_tester_test.go:70/:115/:148`（T5#3）、`guard/p83_fixround6_tester_test.go:25/:72`（T5b/c）全绿方闭合。平台标注：W1/W2 Windows-only（本机 + Windows CI）；Linux-only 的既有 T4c 归 Ubuntu CI；其余跨平台。

## ⑤ 可证伪闭合断言

- **T4**：「Windows 上任何经 `OpenBoundDirFile`/`BoundDirNames` 成功的后代打开，都源自以绑定句柄为根的逐组件句柄链；链上任一 reparse 组件要么被拒绝、要么其目标经卷号 + Ordinal 前缀 + 绑定根仍指向绑定对象的三重证明在界内（末级 reparse 一律拒绝），且最终句柄名与绑定根最终名按 Ordinal 前缀成立、卷号相同。」（Rev 1 修正后口径）
- **T5#4**：「store-local proof 的批量改写只接受同 ResourceID、同 claim、代次更小、且 hash 落在当前记录经 PreviousHash/retained 可达历史链上的 proof；任何其他 proof 不改写并报错。」

## ⑥ 依赖与顺序 + 风险

- **相互独立性**：T4（evidence/gates）与 T5#4（guard/console）无共享改动面 ⇒ 任意顺序、可并行落地；建议 T5#4 先行（纯 Go 逻辑、无 ABI 风险），T4 的 NtCreateFile 走链留足本机 Windows 验证窗口。
- **对 P83 已闭合守卫不回退**：T4 不得删 `:101` 末级 OPEN_REPARSE_POINT 语义、`confirmOpenedArtifact` 卷号/索引复核、Unix openat2 掩码；T5#4 不得动 Renew/Reclaim 停机证明（`lease.go:108`、`:136`）、retained 复核（`:328-352`）、T5#3 限读（`lease_open_*.go`）。合并回归跑 ④ 护栏全集。
- **Windows API 可行性风险**：① NtCreateFile 相对名 + OBJECT_ATTRIBUTES 长度/对齐 ABI（amd64/386），NTSTATUS 误读会全线误拒/误放 ⇒ W1/W2 本机先行 + `confirmOpenedArtifact` belt 兜底（误判方向 = fail-closed，可接受）；② `GetFinalPathNameByHandleW` 大小写形态因 FS 而异——真后代前缀字节相同由同 API 保证，若个别 FS 反例出现，Ordinal 检查误拒真后代（fail-closed 不泄界外），不得退回 EqualFold，改为记录并修 handleFinalPath 归一化；③ 目录句柄关闭竞态：走链期间组件句柄须在每步立即 close，防句柄复用（同 `BoundDir` Close 纪律）。
- **门禁**：只改 Go 文件（无 BOM + LF）；`gofmt -l`、`go build ./...`、`go vet ./...`、`go test ./...` + 聚焦门禁 `go test ./internal/evidence/... ./internal/guard/... ./internal/console/... ./internal/gates/... -run TestP84 -v`。

## 起跑输入与磁盘现状核对

- `persistRenewedPublishLease` 起跑「约 :1166-1208」= 磁盘 `:1167-1208` ✓；`verifyPublishLeasePredecessor` 起跑「约 :2104、:2156」= 磁盘调用 `:2107`、定义 `:2156` ✓；`bounddir_windows.go` 起跑「440 行」= 磁盘 ~438 行 ✓。
- 唯一偏差：起跑称 `renewPublishLease`「约 :2040-2100」，磁盘实为 `:2068-2142`（前移 ~28 行）；判定为约数误差，本方案书一律以磁盘行号为准，无实质影响。
- T4 项 4「无产物不打开根」起跑要求核对现状——已短路且已钉测试（`runner.go:76-86/:271-274/:288-290`；`gates/p83_fixround5_tester_test.go:37-63`），本片该项零代码改动。
- `BoundDirNames` 注释「rel == "" (or ".")」（win `:135-137`、unix `:254-258`）与 `cleanBoundRelativePath` 拒 `"."`（`bounddir_clean.go:29`）现况矛盾 = 起跑 T4 项 4 所述缺陷本体，按 ②(d) 处置，非新增矛盾。

## 修订记录（主控落盘；② 实现三轮实证后更正）

- **Rev 1（2026-11-21，② 修复轮 2 / 3 实证）**：
  1. ②(c) 语义更正：既有 P52 契约测试（`internal/adapters/bounddir_escape_windows_test.go`）钉住「中间 junction 目标是界内 ⇒ 必须放行；界外 ⇒ 拒绝且错误含 `outside the bound directory`」。故中间 reparse 的最终口径 = **末级 reparse 一律拒绝；中间 reparse 跟随，但目标必须经三重证明在界内，界外 / 不可证即拒绝（错误含 `outside the bound directory`）**。原方案书「链上任一 reparse 组件即拒绝」作废；⑤ T4 闭合断言同步改为：「链上任一 reparse 组件要么被拒绝，要么其目标经卷号 + Ordinal 前缀 + 绑定根仍指向绑定对象的三重证明在界内，且末级 reparse 一律拒绝」。
  2. 走链 DesiredAccess 由 `GENERIC_READ` 改为显式位（目录 = `FILE_LIST_DIRECTORY|FILE_READ_ATTRIBUTES|SYNCHRONIZE`；文件 = `FILE_READ_DATA|FILE_READ_ATTRIBUTES|SYNCHRONIZE`）：实测 `GENERIC_READ` 句柄上 `GetFileInformationByHandleEx(FileIdBothDirectoryInfo)` 返回 ERROR_INVALID_PARAMETER(87)，且 `os.File.ReadDir`（Go ≥1.23 句柄版）在 NtCreateFile 句柄上以 Restart 类永久阻塞（600s 挂死实证）；加 `FILE_SYNCHRONOUS_IO_*` 会令 NtCreateFile 返回 STATUS_INVALID_PARAMETER，不可加。
  3. 枚举改为自实现 `readDirNames`（`GetFileInformationByHandleEx`，class = 10 = `FileIdBothDirectoryInfo`，非 3——3 在该枚举中是 FileRenameInfo；派单初稿的 3 有误，以 ② 实证为准），替代 `os.File.ReadDir`：Go 1.22（CI 工具链）的 ReadDir 按名 FindFirstFileW，相对名会枚举错误目录；Go ≥1.23 句柄版在 NtCreateFile 句柄上挂死。过滤 `.` / `..`，终止码 = `ERROR_NO_MORE_FILES` / `ERROR_FILE_NOT_FOUND`，buffer 64 KiB 起 2 倍增长至 1 MiB 上限。
  4. ④ 测试点 W2 更正：界外中间 junction 两条路径（产物打开 / 枚举）均拒绝；**界内中间 junction 产物打开必须放行**（P52 钉住），枚举目标 reparse 一律拒绝（`p83_fixround8_windows_test.go:181` 钉住）。

## 修订记录 Rev 2（2026-11-21，⑦ 盲审 FINDINGS 整改设计；主控落盘）

盲审（GPT-6.2 Sol Max，§7.7 必抽，终审前）报 2 High + 3 Medium + 1 Low。整改设计如下（② 修复轮 6 逐条实现，③ 轮 3 补反例）：

- **F1（High，中间 reparse 界内证明 TOCTOU）→ 改机制**：废止「内核跟随 reparse + 事后名字核对」；改为 **FSCTL_GET_REPARSE_POINT 读取 junction 存储目标 → 目标绝对路径与绑定对象最终名做 Ordinal 前缀判定 → 界内则从绑定句柄按目标的相对后缀逐组件走链（全程句柄相对，界内由构造保证）；界外 / 前缀不匹配 / 相对目标 symlink ⇒ 拒绝**。支撑：junction（IO_REPARSE_TAG_MOUNT_POINT）目标恒为绝对路径；symlink（IO_REPARSE_TAG_SYMLINK）绝对目标同法、相对目标拒绝（fail-closed，需提权创建、测试 skip 登记）。末级 reparse 拒绝语义不变。
- **F2（High，同任务跨进程重入被拒）→ `renewPublishLease` 失配分支重排**：token/hash 失配时先用 `verifyPublishLeaseHistory`（多代）验证 proof 是当前记录的真实前驱；真实前驱且当前记录 holder == 本进程 ⇒ 追平（现状）；真实前驱且 holder ≠ 本进程 ⇒ 走 Reclaim（停机证明）后刷新全部 proof；非真实前驱 ⇒ 拒绝。同进程分支的严格 N-1 `verifyPublishLeasePredecessor` 与 C4 钉住不变。兄弟任务懒收敛（各自晋升时追平）登记为设计残余、不放宽。
- **F3（Medium，代次提交后进程死亡留下租约锁）→ 陈旧锁自愈**：`guard` 新增 PID 活性原语（Windows = OpenProcess+GetExitCodeProcess 经 kernel32 LazyDLL；Unix = kill(pid,0)），锁持有者把 PID 写入锁目录内 marker；竞争者在锁存在时读 marker，持锁进程已死 ⇒ 移除陈旧锁后重试；活性不可证 ⇒ 继续等待（fail-closed 不盲删）。不引入新依赖；`ctx` 超时语义不变。
- **F4（Medium，Publish 先写本任务 proof 后做相等校验）→ 写前校验**：新增 `assertPublishLeaseCurrent(store, resourceID, refreshed)`（store.Current 的 token/hash 与 refreshed 相等，否则错误），`Publish` 在写本任务 proof 之前调用；批内 M1 校验保留；残余窗口（校验后、写前并发推进）由发布末尾 `verifyPublishLease` 复验兜底（fail-closed），在代码注释与报告登记。
- **F5（Medium，根枚举按派生名重开）→ 直接用绑定句柄**：`BoundDirNames` 根分支删除 `currentName()+CreateFile(base)` 重开，直接 `readDirNames(dir.handle)`（绑定对象本身即枚举句柄，无名字解析）；`confirmEnumeratedDirectory` 根分支随之撤销（非根分支保留）。无名字解析即无换根窗口，严格更强。
- **F6（Low，UTF-16 长度截断）→ 长度守卫**：新增 `utf16ComponentLength(name) (uint16, bool)`，字节长 > `math.MaxUint16` ⇒ false；`openRelativeComponent` 据此拒绝超长组件。

## 修订记录 Rev 3（2026-11-21，⑦ 终审 FINDINGS 整改设计；主控落盘）

⑦ 终审第 1 轮：F5 / F6 判 CLOSED（F5 的 `confirmEnumeratedDirectory` 根分支保留获裁定可接受）；F1 / F2 / F3 / F4 OPEN + 新发现 N1（High）/ N2 / N3 / N4 / N5 / N6（Medium）。整改设计（② 修复轮 7 实现，③ 轮 4 补反例）：

- **N1+N2（F3 的锁基座重设计）→ PID 名文件锁**：废弃 Mkdir 目录锁。锁身份 = 文件名自带（`<leasePath>.lock.<pid>`，`CREATE_NEW`/`O_CREAT|O_EXCL` 原子创建 ⇒ 身份与占锁同原子，消除「已占锁、无可恢复身份」窗口）；持有句柄以 `FILE_SHARE_READ|WRITE|DELETE` 打开（Windows）⇒ 自身可删除、释放 = 先删文件再关句柄（删除先于关闭 ⇒ 崩溃不残留活 PID 文件）；竞争者枚举 `lock.*` 解析 PID，**只删可证已死的 PID 锁**（活 PID 一律等待；PID 复用窗口登记为理论残余，威胁模型内不可达）；活性原语复用 Rev 2 的 OpenProcess/GetExitCodeProcess 与 kill(pid,0)。删除与再获取竞态：删除失败（共享冲突/不存在）⇒ 重试循环，绝不先移走活锁再复原。
- **N3（F1 的归属判定）→ SubstituteName 为唯一权威**：`reparseStoredTarget` 解析 REPARSE_DATA_BUFFER 时归属判定只用 **SubstituteName**（剥离 `\??\` / `\\?\` / `\GLOBAL??\` 前缀），`PrintName` 仅入诊断/错误文本；二者不一致时以 SubstituteName 判（界内则走链、界外则拒）。补解析反例（构造 PrintName ≠ SubstituteName 的 REPARSE_DATA_BUFFER 直测解析函数）。
- **N4（F2 的批写收敛）→ 恒等也批**：删除 `persistRenewedPublishLease` 的 `refreshed == previous` 早退；`ensurePublishLease` 与 `Publish` 在「本任务 proof 已当前」时同样重跑幂等批量对账（`existing == refreshed` 的项跳过、落后兄弟追平、异类 proof 照旧报错）；P83 严格 N-1 判据不动。
- **N5（F4 的锁内批写）→ RunLocked**：guard 新增导出 `LeaseStore.RunLocked(ctx, resourceID, fn)`（内部 `withLock` 读当前记录后调 fn）；`persistRenewedPublishLease` 整体改为在 RunLocked 内执行（锁内重读 current → M1 等值校验 → 走链校验 → 落盘），与代次推进共用同一把锁。
- **N6（重入的 run/consumer 守卫）→ 三处统一**：`renewPublishLease` 失配分支与 Reclaim 前、`verifyPublishLease` 终验前统一校验 `publishLeaseBelongsToRun(lease, runID)`（ResourceID + ConsumerID），Reclaim 回调内对锁定记录重复校验；真实历史 ≠ 该 consumer 的发布授权。
- ③ 轮 4 补 C11（N1 活锁不可删 / 死锁自愈）、C12（N3 解析反例）、C13（N4 恒等批追平）、C14（N5 批必须等锁：外部持 RunLocked 阻塞时批在短 ctx 内失败）、C15（N6 异 consumer 拒绝）；W7 的变异口径按 ⑦ 更正（阈值恒为偶数、有效变异 = 删两字节预留）。

## 修订记录 Rev 3b（2026-11-21，③ 轮 4 实证否决 PID 名文件锁；主控落盘）

③ 轮 4 实测：PID 名文件锁（`<record>.lock.<pid>` + CREATE_NEW）**破坏跨进程互斥**——两进程各建不同名文件互不阻塞（子进程持锁时父进程 Claim 7.8ms 成功），且 heal 永不触发。Rev 3 的 N1+N2 机制废止，改用 **OS 级锁（无陈旧锁、无 PID 文件、无 heal）**：

- **Windows**：单一争用名 `<record>.lock`，`CreateFile(OPEN_ALWAYS)` 后 `kernel32!LockFileEx`（EXCLUSIVE + LOCKFILE_FAIL_IMMEDIATELY）；争用 = ERROR_LOCK_VIOLATION ⇒ 等待（5ms / ctx 语义不变）；释放 = `UnlockFileEx` + `CloseHandle`（或仅关闭）。**进程死亡（含崩溃）由 OS 自动释放锁** ⇒ 陈旧锁问题按构造消失，N1/N2 同时消除。
- **Unix**：同一争用名 `os.OpenFile(O_CREATE|O_RDWR)` + `syscall.Flock(LOCK_EX|LOCK_NB)`；争用 = EWOULDBLOCK ⇒ 等待；释放 = Flock(LOCK_UN) + close（close 亦释放）；进程死亡由内核释放。
- 删除 `healStaleLeaseLocks` / PID 活性原语的锁用途（`lease_lock_*.go` 收缩为 acquire/release 包装；若 `leaseHolderAlive` 等不再被引用则删除，unsupported 平台保持拒绝语义）。
- ③ 轮 5 改写 C9 / C11：C9 = 子进程持锁后死亡 ⇒ 父进程无需任何 heal 即可取得（OS 回收实证）；C11 = 子进程存活持锁 ⇒ 父进程短 ctx 内等待失败（活锁互斥）；变异 = 把 LockFileEx/Flock 换成无操作 ⇒ C11 红。

## 修订记录 Rev 4（2026-11-21，⑦ 复审 r1 FINDINGS 整改设计；主控落盘）

复审 r1：F2 / F3 / F5 / F6 / N1 / N2 / N3 / N4 / N6 = CLOSED；F1 / F4 / N5 = OPEN + 新发现 N7 / N8 / N9。整改设计（② 修复轮 9 实现，③ 轮 6 补反例）：

- **N5 残余 / F4 残余（Medium，本任务 proof 锁外写）**：`Publish` 的「等值校验 + 本任务 proof 落盘」（runtime_supervised.go 约 :1971/:1979）移入 `store.RunLocked` 临界区（与代次推进共用同一把锁）；`ensurePublishLease` 的 claim/adopt 分支在批后插入本任务 proof 的写（约 :1167）同样移入 RunLocked。锁外仅收集路径。验收 = 实际 `Publish` 入口 + 外部持 RunLocked ⇒ 短 ctx 失败且 **proof 字节不变**（C16）。
- **N7（Medium，目标归一化 TrimSpace 误判）**：`reparseStoredTarget` 的归一化只做**精确前缀剥离**（`\??\` / `\\?\` / `\GLOBAL??\`），**删除 TrimSpace**；带合法空白类字符（如末尾 U+00A0）的名字保持字节原样，由 Ordinal 后缀判定决定（等于/严格后代才准入）。反例 = 界内正例 + 界外近名兄弟（绑定名 + 尾 U+00A0）拒绝（W8）。
- **N8（Medium，旧 Mkdir 锁目录迁移边界）**：`acquireLeaseLock` 检测「锁路径上是目录」（旧版 Mkdir 锁遗留）⇒ 返回**独立错误**（文本说明旧锁目录须在旧 run 停机证明后人工迁移），**不盲删、不无限等待**；边界写入代码注释与报告（升级边界：旧 run 的 state/store 保留时由操作员处置）。
- **N9（Low，归档方案滞后）**：主控把含 Rev 3b / Rev 4 的方案书同步进证据目录并重建包（本项无代码改动）。

## 修订记录 Rev 5（2026-11-21，⑦ 复审 r2 FINDINGS 整改设计；用户授权：裁决人 = 用户、日期 = 2026-11-21、理由 = r2 余 2 Medium + 1 Low 方向明确、整改成本低，末轮额度完整使用、依据 = 本消息）

复审 r2：2 Medium + 1 Low（F1 原反例已闭、F4/N5 实现已闭但验收面 OPEN）。整改设计（② 修复轮 10 实现，③ 轮 7 补反例）：

- **N10（Medium，T4 目标保真）**：`boundRelativeSuffix` / `normalizeReparseTarget` 的准入前处理废止 `filepath.Clean` 折叠。新口径 = **支持前缀（`\??\` / `\\?\` / `\GLOBAL??\`）精确剥离后，目标须字节等于其 `filepath.Clean` 形态（含 `..` 段、重复分隔符等一切非规范形态 ⇒ 无法无损解释 ⇒ 拒绝，fail-closed），再以原字节做 Ordinal 等于 / 严格后代判定**。保留 P52 / W8 正例。反例 = 解析层负例（`D:\Bound\x\..\outside` 类非规范目标拒绝；规范目标放行）。
- **N11（Medium，T5#4 入口验收）**：产品代码不动（r2 已确认锁内写实现闭合）。③ 轮 7 复用 `internal/console/p83_fixround3_tester_test.go:123` 发布入口夹具补**免续租场景**反例（同进程真实 N-1 proof 先追平、未到续租阈值直接进入本任务写分支）：外部持 `RunLocked` ⇒ `Publish` 不完成且 proof 字节不变；释放 ⇒ 完成并绑定当前摘要。变异 = 仅绕过 `Publish` 调用点（回退旧写路径），字节断言仍须红（即使批处理超时）。
- **INSTALLATION 文档联动（用户 2026-11-21 授权项 ③）**：`docs/INSTALLATION{,_EN}.md` 升级流程说明补「旧版本 Mkdir 目录锁遗留：升级后首次 Claim/Reclaim 若遇 `ErrLeaseLockPathIsDirectory`，须在旧 run 停机证明后人工移除遗留锁目录」段；列入本片文档同步清单（该对不在「不动」清单）。
- **N12（Low）**：⑧b 收尾时更新 MANIFEST 描述至 Rev 5 并复核证据包。

## 修订记录 Rev 6（2026-11-21，⑦ 复审 r3 FINDINGS 整改设计；用户裁定「额度例外 + 兜底」：裁决人 = 用户、日期 = 2026-11-21、理由 = N13 为确定性界内误拒、修复方向明确且小，收窄复审一轮值得、依据 = 本消息）

复审 r3：T5#4 全闭合；T4 余 N13（Medium 阻断）/ N14（Low）/ N12（Low 收尾）。整改设计（② 修复轮 11 实现，③ 轮 8 补反例；本授权超出 §7.5 与本次切片额度裁定，系用户直接裁定、不得作为先例）：

- **N13（Medium，规范盘符根误拒）** `boundRelativeSuffix`（bounddir_windows.go 约 :952-977）：等值比较改用 **保留规范根名尾分隔符** 的形态（`filepath.Clean(base)` 不 TrimRight）；仅在构造**严格后代前缀**时 `TrimRight` 分隔符后拼 `Separator`（`D:\` 根 ⇒ 严格后代 `D:\sub` 前缀 `D:\`，rest 非空放行；根自身等值分支先命中）。不放宽任何现有拒绝条件（`..`/正斜杠/非规范形仍拒）。
- **N14（Low，SubstituteName 内嵌 NUL）** `reparseStoredTarget`（约 :844）：解码前扫描声明 UTF-16 范围内的**内嵌 NUL 单元**，命中 ⇒ 拒绝（fail-closed，解析器错误文本含 NUL 偏移），再交 `UTF16ToString`。
- ③ 轮 8：W10 = 盘符根等值（`D:\` == `D:\` ⇒ 空后缀）+ 严格后代正例（`D:\sub` ⇒ `["sub"]`）；变异 = 恢复 TrimRight 等值形态 ⇒ W10 红。W11 = 解析器内嵌 NUL 反例（手工 REPARSE_DATA_BUFFER，SubstituteName 声明范围内含 NUL ⇒ 拒绝非截断）；变异 = 删除 NUL 检查 ⇒ W11 红。
- 收窄复审输入包 = 仅 N13/N14 两处 diff 与逐条闭合状态（用户裁定收窄）。兜底：若仍报 Medium+ ⇒ 不追加轮次，按 (b) 收口（N13/N14 登记残余顺延，结论 PARTIALLY CLOSED）。
