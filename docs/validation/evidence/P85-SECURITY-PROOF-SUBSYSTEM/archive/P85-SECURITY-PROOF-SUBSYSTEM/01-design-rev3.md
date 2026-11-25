# P85 `SECURITY-PROOF-SUBSYSTEM` · ① 架构增补方案书 Rev 3（single-round）

- 切片：P85 · ① 架构师 = deep-reasoner（DeepSeek V4 Pro）；本文件为 **single-round 增补设计**，
  只做设计、不实现代码。
- 承接输入：⑦ 复审 r1 原文（`04-reviews/07-re-review-r1.md`，8 项出处）+ ⑦ 终审原文
  （`04-reviews/07-final-review.md`）+ 磁盘现状（`internal/{evidence,console,adapters,guard}`）。
- 行号口径：**以 2026-11-23 磁盘现状为准**（本轮 PowerShell Select-String 实测；r1 原文的
  行号系旧快照，若与本版不一致以本版为准并在各条注明）。
- 编码：UTF-8 **无 BOM + LF**（tmp 工作区文件）。
- 8 项逐字承接 r1 的 8 条（其中第 8 项并入第 4 项，测试矩阵单列一行）；每项固定五段：
  现状句（file:line）→ 根因一句 → 原语设计（精确到函数与调用序）→ TOCTOU 论证 → 反例设计（③ 可执行）。

---

## 1 F01 残余：重入取证扩域 + `ManagedWriteProof` 去导出

**现状句。**
- `internal/console/runtime_agent.go:73`（`AgentRunnerAssemblyOptions`）、`:148-160`：
  `ManagedWriteProof string` 是**导出**字段，注释自称 test seam——任何外部包构造 options 时
  可传入任意合形摘要，使 `agentRunnerProvePolicyWriteProtection`（`agent_runner_policy.go:778`
  「injected」分支）跳过真实平台取证。
- `agent_runner_policy.go:411-416`（复用分支）与 `:458-465`（创建分支）只对 **policy 目录**
  取写权证明；既有制品/锚文件自身可写、父目录（runRoot 内）可替换的情形不在观察域内。
- `internal/adapters/agent_runner_policy_write_guard_unix.go:29-62`：观察向量仅
  `path:mode:uid:gid`（:53），**不含 dev:ino**；只观察目录本身，不含父目录与既有文件。

**根因。** 写域证明的观察对象只有「目录本体」，且注入缝是导出字段——重入时无法证明
「既有制品/锚不可被受管身份改写、父目录不可替换」，也无法排除外部伪造合形证明。

**原语设计。**
1. 去导出：删除 `AgentRunnerAssemblyOptions.ManagedWriteProof`（`runtime_agent.go:160`），
   改为非导出字段 `managedWriteProof string`（struct 仍在 :73）。`persistAgentRunnerFrozenPolicyProven`
   （`agent_runner_policy.go:356`，签名已收 `managedWriteProof string` 参数）改从
   `options.managedWriteProof` 取值；生产装配路径（`cmd/prfrail` → console 命令构造 options）
  无此字段可设 ⇒ **编译期封锁旁路**。包内测试（`package console`，已确认
   `runtime_agent_e2e_test.go:398` 等全部 in-package）直接设 `options.managedWriteProof`。
2. adapters 新原语（替代 `VerifyManagedAgentCannotWriteDirectory`）：
   `VerifyManagedAgentCannotWritePolicyDomain(directory string, existing []string) (string, error)`。
   观察向量（逐对象一行，按 `目录\n父目录\nfile1\nfile2…` 固定序，`existing` 排序后拼入）：
   - 对象 1 = policy 目录：`dev:ino:mode:uid:gid` + 既有 access-ACL 拒绝（F08 规则）；
   - 对象 2 = 父目录 `filepath.Dir(directory)`：同样五元组 + ACL；其变更权被
     `installSourceEffectiveWritability`（属主/other-write/group-write-in-group）拒绝判定，
     **覆盖「rename/replace policy 目录项」**；
   - 对象 3..n = `existing` 内每个既有文件（复用分支传 artifact + anchor 两路径）：叶级
     `Lstat` 五元组 + ACL + 有效写权（属主=运行用户即拒绝——既有文件可被 in-place 改写）。
   `dev:ino` 入向量 ⇒ 换目录、换文件、换父目录必改变摘要。Linux 半壁实现放
   `agent_runner_policy_write_guard_unix.go`；Windows 半壁对应扩展。
3. console 调用序：复用分支（`agent_runner_policy.go:411-416`）在**读取/采纳任何字节之前**
   以 `existing = [artifactPath, anchorPath]` 调新原语；创建分支（`:458-465`）以
   `existing = nil` 调。两者都发生在 `agentRunnerConfirmExistingPolicyDirectory` /
   `agentRunnerConfirmPolicyDirectory`（:827/:810，第 2 项改为句柄化）之后。
4. 取证缝收紧：`agentRunnerProvePolicyWriteProtection`（:778）保留包内参数但注入值
   仅测试可见（调用点只在本包）；生产分支改调新原语。

**TOCTOU 论证。** 目录观察（对象 1）在 bound 句柄上做（第 2 项原语）；父目录与文件观察
是 name-based，其窗口由**同一次观察的其它项封闭**：若目录已被证明「受管身份不可写」且
父目录「不可替换」，则目录内条目集合不会被受管身份改变 ⇒ 对 `existing` 文件的 Lstat 与
随后 handle-based 读取指向同一对象。若目录可写或父目录可替换，观察本身失败 ⇒ 拒绝。
替换者若非受管身份则属 SECURITY.md 本机进程边界（登记，不豁免受管身份要求）。

**反例设计（③）。**
- 既有制品可写：构造「目录不可写（root 属主、0700）、artifact 文件属主=运行用户」布局，
  复用分支须拒绝（现实现绿：只验目录）。恢复布局 → 绿。
- 父目录可替换：policy 目录不可写但父目录属主=运行用户 → 创建与复用均拒绝（现实现绿）。
- 旁路封锁：全仓 grep `ManagedWriteProof` 仅剩本包与测试引用；`go build ./...` 通过即证明
  外部包无入口（编译期证据，不需运行时反例）。
- 单变异红：① 删观察对象 2（父目录）⇒ 父目录反例转绿；② 删观察对象 3..n ⇒ 制品反例转绿。

---

## 2 F02/F06/F20 残余：`VerifyBoundDir` 绑定前置检查 fileID + T2 发布句柄化

**现状句。**
- `internal/evidence/bounddir_unix.go:460-489`：`boundDirVerify`/`VerifyBoundDir(dir, expected)`
  用 `os.Stat(expected)` **按当前名解析**（跟随链接）与 `Fstat(dir.fd)` 比较——比较对象是
  「名字现在的解析结果」，不是「前置检查批准的对象身份」。
- 调用点 `agent_runner_policy.go:693`（`:816/:836` 同族）、`console/ai_install_proof.go:211/231`、
  `ai_install_proof_publish_linux.go:49` 都在 `OpenBoundDirNoFollow`（`bounddir_unix.go:144`，
  内部 `openBoundDirNoFollow` :336）之后把 **path** 传回 `VerifyBoundDir`。
- T2 发布 `agent_runner_policy.go:732 publishAgentRunnerPolicyFile`：
  `os.CreateTemp(directory,…)` + `os.Link(temporaryPath, path)` + `agentRunnerSyncDirectory`
  （:849，`os.Open(path)` **按名**）——三次按名操作，与 T3 的句柄内发布
  （`bounddir_publish_linux.go:35/59/105`）不同形。

**根因。** 前置检查所得目录身份被丢弃，验证只比较「两次按名解析一致」，换目录后两次
都命中替换对象仍绿；发布与同步不经过已验证句柄。

**原语设计。**
1. evidence 新身份原语：
   - `type FileID string`（"dev:ino" 格式，与 `BoundFile.FileID()` 一致）；
   - `(dir *BoundDir) FileID() (FileID, error)`：mu 下 `Fstat(dir.fd)` 派生；
   - `VerifyBoundDir(dir *BoundDir, expected FileID) error`：**改签名**，Fstat 句柄并与
     `expected` 逐字比较，**删除 os.Stat/路径参数**（bounddir_unix.go:460-489 重写；
     全仓调用点同步，本包与 console/adapters 均已清点）。
   - `OpenBoundDirVerified(path) (*BoundDir, FileID, error)` = `openBoundDirNoFollow`
     （一次 `openat2`，`RESOLVE_NO_SYMLINKS|NO_MAGICLINKS|BENEATH`）+ `Fstat` 取 fileID。
     内核在**同一次解析**内完成组件级拒链并返回对象 ⇒ 「前置检查」与「句柄」同源，
     没有「检查后换对象」的中间窗口。
2. console 迁移调用序：
   - `openAgentRunnerPolicyFile`（:674-696）→ `OpenBoundDirVerified(directory)`；
     `VerifyBoundDir(bound, fileID)` 仅作持有期一致性断言；
   - `openAIInstallSourceProof`（`console/ai_install_proof.go:205`）、
     `aiInstallSourceConfirmDirectory`（:228）、`writeAIInstallSourceProofFile`
     （linux :43）同改。
   - `agentRunnerRejectUnsafePolicyDirectory`（策略级：workspace 重叠、组件不安全）保留，
     但身份绑定不再依赖它——身份来自 bind 本身。
3. T2 发布句柄化（对齐 T3）：`publishAgentRunnerPolicyFile`（:732）改为：
   `OpenBoundDirVerified(directory)` → `bound.CreateFileExclusive(temp, 0o600)`
   （bounddir_publish_linux.go:35）→ write + `Sync` → `bound.LinkNoReplace(temp, base)`（:59）
   → `bound.Remove(temp)`（:105）→ `bound.SyncDirectory()`（**新方法**：`fsync(dir.fd)`，
   替代 `agentRunnerSyncDirectory` 的按名 `os.Open`）→ Close。临时名生成复用
   `aiInstallSourceProofTemporaryName` 同族随机名。Windows 半壁无句柄 publish 原语，
   保持 path 形式并走已登记 F04 豁免 + fileID 断言（平台分派注明）。
4. 契约与文档：`CONTRACTS.md` 若描述「目录复核按名」措辞，同步为「绑定前置检查所得
   目录身份（dev:ino）」；属 docs 同步，随本片一并复核。

**TOCTOU 论证。** `openat2` 单次解析返回 fd，此后 `CreateFileExclusive`/`LinkNoReplace`/
`Remove`/`SyncDirectory` 全部相对**同一 dirfd**（O_EXCL/linkat/unlinkat/fsync，内核只在该
inode 上查单条目名）；名称被 rename/替换不迁移 fd；`VerifyBoundDir` 的 Fstat 断言与句柄
同源。因此「检查的目录」「读的目录」「写的目录」「同步的目录」是一个对象，且无任何
按名重解析步骤。

**反例设计（③）。**
- 目录替换：bind 后 rename 原目录、同名放替身目录 → 所有相对操作仍落原 inode（读回字节
  一致）。变异：把 T2 发布退回 `os.CreateTemp`+`os.Link` ⇒ 写落替身目录 ⇒ 红。
- `VerifyBoundDir` 变异：改回 `os.Stat` 按名比较 ⇒ 构造「前置检查批准 D1 → 检查后改名
  D1、D2 就位 → bind 得 D2」场景，变异版绿 / 正版红（Fstat≠captured fileID）。
- 同对象目录替换反例（r1 点名）：前置检查 + bind 之间换目录（第 2 项原语下单次解析即
  拒绝，需以「检查与 bind 分两步」的旧实现作变异对照）。

---

## 3 F07 残余 + 新 High：原始名称链不再被 `filepath.Dir/Clean` 消解 + shebang 256 字节窗口

**现状句。**
- `internal/adapters/ai_install_proof.go:366-372 openInstallSourceLeaf`：先
  `filepath.Dir(path)`（:367）与 `filepath.Base(path)`（:372）再 bind/开叶——Dir/Base
  内部 Clean，**可消去原串中的 `..`/`.`/重复分隔符与链接组件**；内核却按 shebang 或
  PT_INTERP 的**原串**解析。
- `ai_install_proof_unix.go:291 installSourcePathComponents` 同样先 `filepath.Clean`；
  `:448 installSourceVerifyInterpreter` 对解释器原串复用 `openInstallSourceLeaf`。
- `:393-400 installSourceShebangInterpreter`：以「去掉 `#!` 的余串长度 `>= 256`
  （`installSourceShebangLineLimit` :88）」判截断，**不计 `#!` 的 2 字节**，也**不验证
  换行是否落在内核 256 字节窗口内**。

**根因。** 证明的字符串 ≠ 内核解析的字符串（Clean 消解原始名称链；shebang 窗口口径与
内核 BINPRM_BUF_SIZE 不一致）。

**原语设计。**
1. 规范拼写校验原语（adapters，linux 文件；evidence 侧 `openBoundDirNoFollow` :336 同步
   获得同一校验）：`installSourceRequireCanonicalSpelling(path string) error` —— 必须绝对
   路径、首字符 `/`、逐 `/` 切分后每个组件非空且非 `"."`/`".."`、无 NUL、尾斜杠仅根
   允许。**非规范即拒绝**（fail-closed：不能证明的拼写不接纳）。规范串上
   `filepath.Dir/Base/Clean` 均为恒等 ⇒ 可继续使用且不改变任何组件。
2. `openInstallSourceLeaf`（:366）：进入即校验原串；`installSourceVerifyInterpreter`（:448）
   对 shebang/PT_INTERP 给出的解释器原串同样先校验再 bind。`installSourcePathComponents`
   （:291）与 `verifyUnwritableInstallPathAgainstLeaf`（:154）的 Clean 前置改为
   「校验后按原串拆分」。逐组件 `Lstat` 拒链（:154 区，既有）不变。
3. shebang 窗口重写（:393）：输入改为 `data[0:min(256, len)]`（整文件前 256 字节窗口，
   **含 `#!`**，`installSourceShebangWindow = 256` = BINPRM_BUF_SIZE，常量注释改为内核
   引用）；先验 `data[0:2]=="#!"`；在窗口内找首个 `'\n'`（文件偏移 ≤ 255）——找不到 ⇒
   拒绝（内核 ENOEXEC，执行对象不可描述）；`spec = data[2:n]`（n ≤ 255 ⇒ |spec| ≤ 253，
   内核语义精确复刻）；按 space/tab 拆首 token（:403-408 既有，保留）；env 拒绝
   （:417）、绝对路径、递归深度（:333 `installSourceInterpreterFromBound` 链）不变。
   删除「余串长度 ≥ 256」判据（:398-400）。

**TOCTOU 论证。** 内核 exec 解析原串（逐组件、遇链接跟随）；我方 `openat2`
（RESOLVE_NO_SYMLINKS）在**同一个原串**上解析——无链接则两解析必然到达同一对象；有链接
则我方先拒 ⇒ 「证明的对象」与「内核将执行的对象」同对象或整体拒绝。Clean 移除后不存在
「我方证明消解后对象、内核解析原串」的错配；窗口判据与内核字节级一致 ⇒ 不存在「证明
的是内核会截断/不执行的行」。

**反例设计（③，Linux 原生）。**
- 原串含 `..`：shebang 写 `#!/a/../b`（`/a` 为链接）→ 新校验拒绝（现实现 Clean 成 `/b`
  证明后放行——绿）。
- 原串含链接组件：解释器路径含 symlink → RESOLVE_NO_SYMLINKS ELOOP 拒绝（保留行为）。
- 窗口边界夹具：换行在文件偏移 255（spec=253 字节）→ 通过；换行在偏移 256（spec=254，
  旧判据 `len<256` 曾放行）→ 新实现拒绝。两夹具钉死「含 `#!`」口径。
- 变异：① 删 canonical 校验 ⇒ 反例 1 绿；② 窗口改回「余串 vs 256」⇒ 反例 3 绿。

---

## 4 F10 残余 + 第 8 项（反例可判别性）：重入三操作收束到同一目录句柄

**现状句。**
- `internal/console/ai_install_proof.go:88-106`（found 分支）：`readAIInstallSourceProof(path)`
  （:155，内部经 `openAIInstallSourceProof` :205 **第一次 bind**）、
  `aiInstallSourceConfirmDirectory`（:228，**第二次 bind**）、`aiInstallSourceSyncDirectory`
  （:243，`os.Open(path)` **按名第三次打开**）——三次独立打开，不能证明同步的是载体所在目录。
- 首轮发布：`ai_install_proof_publish_linux.go:43 writeAIInstallSourceProofFile` 自开 handle
  写完 link 后，再 `aiInstallSourceSyncDirectory(filepath.Dir(path))`（按名，:86 区）。首轮
  link 可见但目录 sync 失败 → 拒绝；重入会采纳载体且**不复补同步**（:88-106）。

**根因。** 读载体、确认目录、同步是三个独立打开；重入缺少同柄强制 sync；现有链接反例
先被读载体前置拒绝，撤掉重入确认或同步仍绿（反例无可判别性）。

**原语设计。**
1. `aiInstallSourceBindCarrierDirectory(path) (*evidence.BoundDir, evidence.FileID, error)`
   = `OpenBoundDirVerified(filepath.Dir(path))`（第 2 项原语）。
2. `readAIInstallSourceProof` 拆为 `readAIInstallSourceProofFrom(bound, base, boundBytes)`
   —— 读相对**传入句柄**；`openAIInstallSourceProof` 仅保留给无句柄调用点并内部走组合原语。
3. found 分支（:88）调用序改为：
   bind 一次 → `readAIInstallSourceProofFrom(bound, base, …)` → `VerifyBoundDir(bound, fileID)`
   （同柄 Fstat 断言）→ `bound.SyncDirectory()`（同柄 fsync；**失败 ⇒ 拒绝，不采纳**）→
   `facts.Matches(existing)`。
4. 首发发布改为 `writeAIInstallSourceProofFileIn(bound, base, payload)`：
   CreateFileExclusive → write → Sync → LinkNoReplace → Remove → **SyncDirectory**，全在同一
   句柄；同步失败传播拒绝（现有 :43 的按名 sync 删除）。
5. 测试缝（③ 可判别性，第 8 项）：包内
   `var aiInstallSourceCarrierDirectorySync = (*evidence.BoundDir).SyncDirectory` 与
   `var aiInstallSourceReentryRequiresSync = true`（均注释为 test-only；生产不读外部配置）。

**TOCTOU 论证。** 单 fd 全链：读载体、身份断言、fsync 作用于同一目录对象；fsync 的对象 =
读的对象 = linkat 的目录；名称 rename 不迁移 fd；目录若被替换，下次 bind 的 fileID 不同
且载体缺失/不匹配 ⇒ 拒绝。不存在「确认的是 A、同步的是 B」窗口，也不存在「link 可见但
未持久即被采纳」的路径。

**反例设计（③，两项可执行）。**
- 反例 A（同步守卫红）：首轮 LinkNoReplace 后注入 `aiInstallSourceCarrierDirectorySync`
  失败 → 首轮拒绝；次轮健康 sync、同柄重入 → sync 成功 → 采纳（红→恢复绿）。
- 反例 B（守卫变异，判别性）：首轮 sync 失败后置 `aiInstallSourceReentryRequiresSync=false`
  → 重入直接采纳 ⇒ 证明拒绝来自同步守卫而非读载体守卫（读载体在此场景通过）。恢复
  `true` → 再次红。此项即第 8 项，测试矩阵单列。
- Windows 注释（:240 区）补充「原子可见 ≠ 断电持久」的 F04/F10 区分措辞。

---

## 5 F13 残余：`CheckSpec` 携带冻结规格摘要，Sandbox 判定绑定该规格

**现状句。**
- `internal/guard/sandbox.go:62 CheckSpec{Reason string}`——唯一输入是原因标签；`:49`
  `Verdict{Status, Reason}` 无规格回显。
- `internal/adapters/ai_probe_copilot_executor.go:251 copilotCLIProbePlatformWriteVerdict`：
  verdict 摘要 = `GOOS+GOARCH+reason+"allowed"`——只绑定平台、标签与固定 allowed，
  不含规格；:252/:523/:938 三处 `CheckSandbox` 只传 Reason；平台实现
  （`sandbox_windows.go` alphaSandbox / `sandbox_unsupported.go`）不读取规格事实。

**根因。** Sandbox 判定与冻结规格解耦：同平台不同写限制事实不可区分，且无法证明
Sandbox 对该规格作过判定。

**原语设计。**
1. guard：
   - `CheckSpec` 增 `SpecDigest string`（BYOK 检查必填，其它 reason 可选）；
   - `Verdict` 增 `SpecDigest string`（实现回显其判定的规格摘要）；
   - `CheckSandbox`（:95）：`spec.SpecDigest != ""` 时要求 `verdict.SpecDigest ==
     spec.SpecDigest`，空/不符 ⇒ 拒绝（复用既有 sentinel，fail-closed）；`Reason ==
     SandboxReasonBYOKProbe` 且 `SpecDigest == ""` ⇒ 直接拒绝（防旧调用点漏传）。
2. adapters：
   - `copilotCLIProbeSpawnSpecDigest(spec copilotCLIProbeSpawnSpec) string` =
     `evidence.Digest(新域, canonical(Executable|WorkingDirectory|AuthMode|Args|
     EnvironmentKeys|OutputBound))`（字段集与 `copilotCLIProbeSpawnSpecEqual` :238 同界；
     不含 `PlatformWriteProof`）。
   - `copilotCLIProbePlatformWriteVerdict(ctx, sandbox, specDigest) (string, error)`：
     经 `CheckSandbox(CheckSpec{Reason, SpecDigest})` 后，verdict digest =
     `Digest(域, GOOS+"\n"+GOARCH+"\n"+reason+"\n"+specDigest+"\n"+status)`。
     `PlatformWriteProof` 因此绑定规格：同平台不同规格 ⇒ 不同 proof。
3. 调用链：`VerifyBYOKConfinement`（`ai_probe_copilot.go:136`）先
   `frozenLaunchTemplate()` → `specDigest` → `verifier.VerifyBYOKConfinement(specDigest)`；
   executor `VerifyBYOKConfinement`（:519）增 specDigest 参数；`assertCopilotCLIProbeBYOKConfinement`
   （:809）用 **frozen spec 的 digest**（Element 0 已比字段等 ⇒ digest 必然一致）重取 verdict。
   兼容路径（无 frozen spec）传自身构造 spec 的 digest。

**TOCTOU 论证。** 摘要由**同一 frozen 对象**派生；spawn 前重观测 verdict 仍对同一摘要；
mismatch ⇒ 拒绝。Sandbox 对「哪个规格」作答成为被核对的事实（verdict 回显 + 比较），
不再是调用方的声称。

**反例设计（③）。**
- guard 包：double 回显错误 `SpecDigest` ⇒ `CheckSandbox` 拒绝（Allowed 也不行）。
- adapters：两同平台规格仅 `OutputBound` 不同 ⇒ `PlatformWriteProof` 不同（现实现相同）。
- 变异：verdict digest 恢复常量 `"allowed"` 尾 ⇒ 上一反例绿。
- 既有 `p85T8ASandboxDouble`（`p85_t8a_tester_test.go:80`）补回显逻辑，旧用例补 SpecDigest。

---

## 6 新 Medium：Publish 对 `ManagedWriteProof` 核验真实冻结锚对证

**现状句。**
- `internal/console/agent_runner_policy.go:597-611 agentRunnerVerifyFrozenPolicyAnchor`：
  对 `anchor.ManagedWriteProof` 仅 `evidence.ValidHash`（:611）——**格式校验**，换成另一
  合形摘要仍通过；此函数是 load（`loadAgentRunnerFrozenPolicy`）与 Review/Publish
  （`runtime_agent.go` `verifyReviewPolicyBinding` 链）共用的锚判定口。
- 写守卫摘要（`agent_runner_policy_write_guard_unix.go:53-54`）不含 dev:ino ⇒ 另一目录
  的观察摘要也可合形通过。
- 现有测试只测非法摘要（`runtime_agent_e2e_test.go:52` 注入 `sha256:111…`）与重开装配
  合形变异，两者不能互证。

**根因。** Review/Publish 只验摘要格式、不对真实冻结锚（本目录 + 既有制品/锚）重取证对证。

**原语设计。**
1. 摘要对象化（与第 1 项同源）：观察向量含每对象 `dev:ino`（第 1 项-2）。
2. `agentRunnerVerifyFrozenPolicyAnchor`（:597）签名增 `policyDirectory string`（或直接收
   本读路径已持有的 bound 句柄 + fileID）；新增
   `agentRunnerReverifyManagedWriteProof(bound *evidence.BoundDir, fileID evidence.FileID,
   recorded string) error`：对**同一句柄**重观测（第 1 项扩展观察：目录+父目录+既有文件，
   句柄形式优先），摘要与 `anchor.ManagedWriteProof` 逐字比较；不等/不可观测 ⇒ 拒绝。
   load 与 Review/Publish 都经此单一判定口。
3. 测试缝：`var agentRunnerObserveManagedWriteProof = adapters.VerifyManagedAgentCannotWritePolicyDomain`
   （包内 var，测试替换为 stub；生产永不读外部配置）。

**TOCTOU 论证。** 重观测在读取锚与制品时同句柄、同 fileID 断言链上进行（第 2 项原语）；
摘要含 fileID ⇒ 对象替换必改摘要；不存在「仅格式」的旁路。

**反例设计（③，Publish 层单独变异红）。** 夹具写冻结制品+锚（grant 不动），仅把锚的
`ManagedWriteProof` 换成**另一个合形摘要**（如另一目录的观察摘要）→ Review/Publish 拒绝；
恢复原字节 → 绿。与既有「非法摘要拒绝」不同源（格式对、对证错）。配合：重观测 stub
返回不同摘要的用例 + 观测失败的用例。

---

## 7 F18 残余：`Allowed` 替身进入真实执行器启动判定链

**现状句。**
- `internal/adapters/p85_t8a_tester_test.go:80-88`：`p85T8ASandboxDouble` 只进入
  `p85T8ASandboxProofProbe`（:108，替身组合），**不进真实 executor 的 Execute/assert 链**。
- `:586 TestP85T8AExecuteRefusesUnprovableVerdictWithoutSpawning`：用**不在白名单的
  shell argv**，即使撤 Sandbox 也被 argv/其它元素先拒（:587-605 区）——不能代替独立
  零启动变异。
- `ai_probe_copilot_executor.go:938`（`assertCopilotCLIProbeBYOKConfinement` 内）与
  `:519-523`（`VerifyBYOKConfinement`）的 `CheckSandbox` 传 `nil`（平台实现），无法注入
  Allowed；启动判定链尾为 `guard.RunManaged`（Execute 尾部，:940 区之后）。

**根因。** Allowed 替身与真实启动判定链脱钩，启动观察点不在链尾。

**原语设计。**
1. `ManagedCopilotCLIProbeExecutor` 增字段 `Sandbox guard.Sandbox`（nil = `guard.NewSandbox`，
   生产装配不设 ⇒ 行为不变）；`:519/:523/:938` 三处 `CheckSandbox` 改传 `executor.Sandbox`；
   `copilotCLIProbePlatformWriteVerdict`（:251）增 sandbox 参数（两处调用点同步）。
2. 独立启动观察点：包内 `var copilotCLIProbeRunManaged = guard.RunManaged`（test-only 注释）；
   Execute 尾部的 `guard.RunManaged(runCtx, spec, …)` 改调该 var。③ 以 stub 替换之：记录
   `ProcessSpec` 并返回 `Started: true`（不真 spawn、无副作用、双平台可跑）。

**TOCTOU 论证。** 观察点位于 RunManaged 之前最后一道门之后 ⇒ 计数=1 当且仅当平台守卫
放行且所有前置守卫通过；Sandbox 对象与真实链共用同一实例 ⇒ 替身判定即真实判定。

**反例设计（③，合规规格 + 只撤平台守卫）。**
- 绿：合规规格（白名单 argv、闭集 env keys、executor-owned workdir、输出界）+
  double `{Allowed, SpecDigest 回显}` → 观察点计数=1。
- 红（单变异）：同一 double 改 `VerdictUnknown`（或 `Sandbox=nil` 走平台 Unknown）→
  计数=0；恢复 Allowed → 计数=1。argv 等其它守卫保持通过 ⇒ 平台守卫是唯一被翻动的门。
- 对照反例：撤 Allowed 且 argv 不合规 ⇒ 仍在 argv 元素拒绝（记录元素级错误），证明新
  反例不被前置门禁掩盖。

---

## 8 ③ 测试矩阵增量

| # | 对应项 | 测试（包/平台） | 反例/夹具 | 变异（红） | 恢复（绿） |
|---|---|---|---|---|---|
| M1 | 1 | console（linux/windows） | 既有制品属主=运行用户 + 目录不可写 | 删 per-file 观察 | 恢复 |
| M2 | 1 | console（linux/windows） | 父目录可替换（属主=运行用户） | 删 parent 观察 | 恢复 |
| M3 | 1 | 全仓 grep + `go build ./...` | 外部包无 `ManagedWriteProof` 入口 | 字段改回导出（编译/审查红） | 恢复 |
| M4 | 2 | evidence/console（linux） | bind 后 rename 目录放替身，写仍落原 inode | T2 发布退回 CreateTemp+Link | 恢复 |
| M5 | 2 | evidence/console（linux） | 前置检查 D1 → 改名 → bind D2 | VerifyBoundDir 退回 os.Stat | 恢复 |
| M6 | 3 | adapters（linux） | shebang `#!/a/../b`（a=链接） | 删 canonical 校验 | 恢复 |
| M7 | 3 | adapters（linux） | 换行偏移 255 过 / 256 拒 | 窗口改回「余串 vs 256」 | 恢复 |
| M8 | 4/8 | console（linux） | 首轮 sync 失败 → 重入同柄 sync 拒绝 | `aiInstallSourceReentryRequiresSync=false` → 采纳 | 恢复 |
| M9 | 5 | guard/adapters | 错误 SpecDigest 回显拒绝；同平台两规格 proof 不同 | verdict digest 恢复 "allowed" 尾 | 恢复 |
| M10 | 6 | console | 锚 ManagedWriteProof 换合形异目录摘要 → Review 拒 | 删对证重验（退回 ValidHash） | 恢复 |
| M11 | 7 | adapters（双平台） | Allowed double 达启动观察点；Unknown → 0 启动 | （即单变异本体） | Allowed 恢复 |
| M12 | 7 | adapters | 不合规 argv + 撤平台守卫 → 仍 argv 元素拒绝 | —（对照） | — |

- 每项执行口径：先正例绿 → 变异红 → 恢复绿 + 字节复原核对；M11/M12 用
  `copilotCLIProbeRunManaged` stub，不真 spawn。

---

## 9 实施次序（供 ⑧ 派单，本轮仅设计）

1. **evidence 原语层**：`FileID`、`BoundDir.FileID`、`VerifyBoundDir` 签名改 fileID、
   `OpenBoundDirVerified`、`BoundDir.SyncDirectory`、`openBoundDirNoFollow` 加 canonical
   校验（第 2、3、4 项公共底座）+ evidence 自测。
2. **第 2 项调用点迁移**：console/adapters 全部 `OpenBoundDirNoFollow+VerifyBoundDir`
   改组合原语；T2 发布句柄化（`publishAgentRunnerPolicyFile`）。
3. **第 3 项**：adapters canonical 校验接入 + shebang 256 窗口重写 + 窗口边界夹具。
4. **第 4/8 项**：console F10 同柄收束 + 两个测试缝 + 反例 A/B。
5. **第 1 项**：`ManagedWriteProof` 去导出 + `VerifyManagedAgentCannotWritePolicyDomain`
   扩域（目录/父目录/既有文件 + dev:ino）。
6. **第 6 项**：`agentRunnerReverifyManagedWriteProof` 对证重验 + Publish 层合形变异红
   （依赖 5）。
7. **第 5 项**：guard `CheckSpec.SpecDigest` / `Verdict.SpecDigest` / `CheckSandbox` 判定。
8. **第 7 项**：executor `Sandbox` 字段 + `copilotCLIProbeRunManaged` 观察点 + M11/M12
   （依赖 5、7 顺序内）。
9. 收口：`go build ./...`、`go vet ./...`、`go test ./...` + ③ 原生反例双平台执行 +
   docs/CONTRACTS 措辞同步复核 → 交 ⑦ 复审。

次序理由：evidence 底座先行避免重复改动；guard 改动放后减少既有测试面；1/6 共享观察
原语相邻落地；7 依赖 5 的 SpecDigest 机制。

---

## 10 登记事项（非设计变更）

- `CONTRACTS.md` 若含「目录复核按名」「sandbox 判定输入仅标签」类措辞，随第 2/5 项
  同步修订（协议先行，先文档后代码，由 ⑧/⑦ 复核）。
- S1 下 alphaSandbox/unsupportedSandbox 恒 Unknown ⇒ BYOK 保持禁用（本设计不改变该结论，
  仅使判定链可证）。
- Windows 的目录句柄同步与最终名称解析残余维持 F04/F10 已登记豁免，不在本增补内扩大。
