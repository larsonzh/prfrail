# P85 `SECURITY-PROOF-SUBSYSTEM` · ① 架构修复方案书 Rev 2

- 切片：P85 · `[SLICE]` · 实现类 L 档 · 协议先行；① 架构师 = deep-reasoner（DeepSeek V4 Pro）
- 承接：① Rev 1 方案书（`01-design.md`）+ ② 已落地 T2/T3/T8a + ② 整改单 1（F01/F02/F09/F10/F19/F20）
- 本版输入：⑦ 盲审原文（`04-reviews/07-blind-review.md`，22 条 = 7H + 11M + 3L + 1I）、
  ⑤ 预审 r1/r2 原文、磁盘现状（`internal/{adapters,console,guard,evidence}`、`schemas/`、
  `testdata/contracts/`、`tools/contracts/contracts.test.js`、`docs/CONTRACTS{,_EN}.md`）
- 行号口径：**以本版步骤 0 磁盘现状为准**；⑦ 盲审锚点对磁盘的重钉见各条「现状句」。
- 编码：UTF-8 **无 BOM + LF**（本 tmp 工作区文件，用户指定）。
- 硬约束核对见 §7；残余风险成文见 §6（交报告 §8）。

---

## 1 覆盖关系表

### 1.1 总表（盲审 22 条 → 处置）

| 编号 | 级别 | 处置 | 本版落点 |
|---|---|---|---|
| F01 | High | **单1已修**（只复核） | §1.2-1 |
| F02 | High | **单1已修**（只复核） | §1.2-2 |
| F03 | High | **本版整改** | §2.1 |
| F04 | High | **本版整改** | §2.2 |
| F05 | High | **本版整改** | §2.3 |
| F06 | High | **本版整改** | §2.4 |
| F07 | High | **本版整改** | §2.5 |
| F08 | Medium | **本版整改** | §3.1 |
| F09 | Medium | **单1已修**（只复核） | §1.2-3 |
| F10 | Medium | **单1已修**（只复核） | §1.2-4 |
| F11 | Medium | **本版整改**（协议先行；判定失误登记 §5.1） | §3.2 |
| F12 | Medium | **本版整改**（证据公式收紧，影响登记 §5.4） | §3.3 |
| F13 | Medium | **本版整改** | §3.4 |
| F14 | Medium | **本版整改** | §3.5 |
| F15 | Medium | **本版整改**（控制台局部输出形状，登记 §5.3） | §3.6 |
| F16 | Medium | **交 ③**（测试整改） | §4.1 |
| F17 | Medium | **交 ③**（测试整改） | §4.2 |
| F18 | Medium | **交 ③**（测试整改） | §4.3 |
| F19 | Low | **单1已修**（只复核） | §1.2-5 |
| F20 | Low | **单1已修**（只复核） | §1.2-6 |
| F21 | Low | **本版整改**（契约补例外措辞） | §3.7 |
| F22 | Info | **登记**（口径 + ⑨ 行为证据清单） | §3.8、§5.2 |

### 1.2 整改单 1 六条一致性复核（不重新设计）

1. **F01（装配冻结锚 + 四等式）已修**：`internal/console/agent_runner_policy.go`（F01 标记于
   :96/:154/:398/:426/:517/:549 区，装配锚 `assembly anchor` 于 :398/:517/:530/:582 区）写出
   装配时冻结锚；`internal/console/runtime_agent.go:1116` 起 `verifyReviewPolicyBinding` 把
   `Digest(磁盘制品 body) == 制品 policyHash == recorded == grant.PolicyHash` 接到锚上
   （F01 四等式 + 锚共五等式）。Rev 2 复核动作：重读锚写入位点与五等式逐项比较，确认
   「撤锚 ⇒ F01 反例红」的变异测试仍在位；不再新增设计。
2. **F02（bound-dir 末组件打开）已修**：`agent_runner_policy.go:487/:594/:628` 区经
   `evidence.OpenBoundDirFile` 打开制品末组件（拒 symlink/reparse）；BYOK 侧 `-C` 断言
   `ai_probe_copilot_executor.go:550` 起（`assertCopilotCLIProbeArgvWorkingDirectory`）在位。
   Rev 2 复核：F13 的 `-C` 冻结裁定（§3.4-3）**收紧但不推翻** F02——F02 断言保留为纵深防御。
3. **F09（删 Rename 降级）已修**：`internal/console/ai_install_proof.go:165-190`
   `writeAIInstallSourceProofFile` 硬链接失败即拒绝（无 `os.Rename` 回退）；`agent_runner_policy.go`
   写入器同族。复核 = grep 确认两写入器均无 `os.Rename` 降级分支。
4. **F10（父目录 sync + 重入冻结判据）已修**：原子写线性化点 = rename/link + 父目录 sync，
   重入不把未确认持久的文件视为已冻结。复核 = 线性化点注释与既有分支的跳过条件读回。
5. **F19（无碰撞键 + 原字节比较）已修**：`agent_runner_policy.go:81/:390/:443/:462` 区键拼接
   不含可碰撞分隔符（或显式转义），比较走磁盘原字节。复核 = 键格式与比较原语读回。
6. **F20（载体组件级检查）已修**：`internal/console/ai_install_proof.go:125-160`（bound-dir
   读）+ `:196-227`（`aiInstallSourceConfirmDirectory` 写前身份复核）。复核 = 与 §2.4 的
   叶级 bound-open 复用关系一致（F06 改造后 F20 的实现方式不变）。

---

## 2 设计级 High（逐条完整设计）

每条按固定四段：现状句 → 改动点 file:line → 原语选择 → 反例与变异测试点。

### 2.1 F03：Windows 合并权限掩码一次 ACCESS_DENIED ≠ 逐项拒绝

**现状句。** `internal/adapters/ai_install_proof_windows.go:156-180`
`installSourceDeniesMutation` 以合并掩码 `installSourceMutationMask`（:57-72，8 个权利位）
做**一次** `CreateFile`（:167）；一次 `ERROR_ACCESS_DENIED` 只证明「未同时获得全部权利」，
不证明逐项被拒。仅缺 `WRITE_OWNER` 仍可写或可删的 ACL 也会获得「不可变」结论。调用方
`:88-118` 只产出 `fileDenies/directoryDenies` 两布尔并入 `ACLEvidenceHash`。

**改动点 file:line。**
1. `ai_install_proof_windows.go:57-72`：把合并常量拆为两张**逐权利表**：
   - 文件表 `installSourceFileMutationRights = []uint32{FILE_WRITE_DATA, FILE_APPEND_DATA,
     FILE_WRITE_EA, FILE_WRITE_ATTRIBUTES, STANDARD_DELETE, STANDARD_WRITE_DAC, STANDARD_WRITE_OWNER}`；
   - 目录表 `installSourceDirectoryMutationRights = []uint32{FILE_ADD_FILE, FILE_ADD_SUBDIRECTORY,
     FILE_DELETE_CHILD, FILE_WRITE_ATTRIBUTES, STANDARD_DELETE, STANDARD_WRITE_DAC, STANDARD_WRITE_OWNER}`
     （目录上「写内容」= 加子项；`FILE_DELETE_CHILD` 独立成项）。
2. `ai_install_proof_windows.go:156-180`：`installSourceDeniesMutation` 改为
   `installSourceMutationDenials(name string, rights []uint32) ([]uint32, error)`——对表内每个
   权利**单独一次** `CreateFile(desiredAccess=该单个权利, share=READ|WRITE|DELETE, OPEN_EXISTING,
   directory?FILE_FLAG_BACKUP_SEMANTICS:FILE_ATTRIBUTE_NORMAL)`：handle 成功 = 该权利可获
   （立即 `CloseHandle`，零写入），记入「可获集合」；`ERROR_ACCESS_DENIED` = 该权利被拒；
   其它错误 = 不能判定 ⇒ 整体 fail-closed。
3. `ai_install_proof_windows.go:88-118`：`verifyInstallSourceImmutability` 改为对文件跑文件表、
   对安装目录跑目录表；`denies = 两张表的可获集合均为空`；`ACLEvidenceHash` 输入改为
   **逐权利向量**（逐字，见下）。`AIInstallSourceWindowsFacts` 结构不变
   （`ai_install_proof.go:59-64`），仅证据输入变。
   ```
   // 逐字：拒绝向量编码
   ACLEvidenceHash = evidence.Digest(aiInstallSourceACLDomain, []byte(
       strings.ToLower(path) + "\n" +
       "objects=" + strings.Join(objectStates, ";") + "\n" +
       "trustee=" + sidText))
   // objectStates[i] = fmt.Sprintf("%s|denied:%s|dacl:%s",
   //     objectPath, strings.Join(deniedRightsHex, ","), descriptorDigest)
   // deniedRightsHex = 该对象上逐权利探测中被拒权利的十六进制字符串，按表序排列；
   // descriptorDigest = installSourceDescriptorDigest(objectPath)（既有 :182 起，保留）。
   ```

**原语选择。** 维持「直接问内核访问检查」的既有哲学（文件头注释已论证）：逐权利
`CreateFile` 与 `CreateProcessW` 走同一访问检查路径，规避 `GetEffectiveRightsFromAclW` 的
CREATOR OWNER 误报；每次探测零写入、即时关柄，不留足迹。不引入新依赖。

**反例与变异测试点。**
- 原生反例（⑦ 点名）：Windows 测试（`//go:build windows`）构造只授予**单一变更权**的 DACL
  （如仅 `FILE_ADD_FILE` 或仅 `WRITE_OWNER`，其余全拒）→ 该权利探测获得 handle ⇒ 判定
  可写 ⇒ 拒绝。夹具构造：临时目录 + `icacls`/SDDL 或纯 Go advapi32 设置受限 DACL，用后还原。
- 变异：① 把某权利删出表（并入合并掩码）⇒ 对应单一权利反例转绿；② 把「可获集合为空」
  改成「合并掩码一次探测」⇒ 同一反例转绿。
- 恢复：`git diff` 零增量；同夹具下新旧载体 `ACLEvidenceHash` 差异仅来自证据输入格式
  （载体为 store-local，无 wire 影响，见 §5.4）。

### 2.2 F04：Windows 祖先命名空间 / reparse / 父级删除权 + 最终启动名绑定

**现状句。**
- `ai_install_proof_windows.go:88-118`：只探测文件与**直接**父目录（`filepath.Dir(path)`）；
  祖先各级目录既不查 reparse，也不查本令牌对每级目录的变更权。
- `internal/guard/process_windows.go:500-507` `openApprovedExecutable`：按最终路径名一次
  `CreateFile`（`FILE_SHARE_READ` 拒写拒删锁**叶对象**）；`CreateProcessW`（`start` :212 起，
  `resolveManagedExecutable` :485）在 spawn 时重新解析整条名称链，链上任何 reparse/替换
  都在该锁之外。
- 既有原语：`evidence.OpenBoundDir`（`bounddir_windows.go:52`）与 `OpenBoundDirFile`（:76）
  打开时拒绝 reparse 点，P83 fix round 8 夹具实证（`p83_fixround8_windows_test.go:126-190`）。

**改动点 file:line。**
1. `ai_install_proof_windows.go:88` 之前新增 `verifyWindowsAncestorImmutability(path)`：
   把绝对路径拆为全组件链（盘根起、每级目录到叶文件），对**每级目录**做三件事：
   ① 拒 reparse——用 `evidence.OpenBoundDir` 打开该级目录（组件级 reparse 拒绝复用既有原语；
     若为叶文件级则用 `OpenBoundDirFile` 语义，见 §2.4-2 的统一改造）；
   ② 逐权利变更权探测——复用 F03 的 `installSourceMutationDenials`（目录表）；
   ③ DACL 指纹——`installSourceDescriptorDigest`（既有）。全链结果并入 `ACLEvidenceHash`
   （§2.1-3 的 `objects=` 列表扩展为**全链**，盘根含在内）。
2. `ai_install_proof.go:257-280` `VerifyAIInstallSource` 不变形；`DirACLDeniesWriteDelete`
   语义覆盖为「文件 + 安装目录 + 全部祖先」（Validate 不改）。
3. 「最终启动名绑定」落到既有 `openApprovedExecutable` 链：`process_windows.go:500-507`
   打开成功后追加 reparse 自检——`GetFileInformationByHandle` 断言对象非
   `FILE_ATTRIBUTE_REPARSE_POINT` 且非目录（目录/reparse ⇒ 关柄拒绝）；`FILE_SHARE_READ`
   拒写拒删句柄行为不变，`buildCreationAttributes`（:166 起）不变。`CreateProcessW` 自身
   的解析瞬间不在用户态可封范围内，残余登记 §6.1。

**原语选择。** 组件级 reparse 拒绝复用 SEC-02a 家族（`OpenBoundDir` / `OpenBoundDirFile`），
不另造路径遍历逻辑；每级目录变更权与 DACL 指纹沿用 F03 的逐权利 `CreateFile` +
`GetNamedSecurityInfoW`（同访问检查、零写入）；叶级拒写拒删仍由既有 held
`FILE_SHARE_READ` 句柄承担（`openApprovedExecutable` 不改打开语义，只加对象类别自检）。

**反例与变异测试点。**
- 祖先 reparse：在安装目录**上级**放 junction（`mklink /J`）指向可写目标 → 探测拒绝
  （现实现只查直接父目录会放行）。
- 单级祖先可授权：仅父级的父级授予 `FILE_ADD_FILE`（F03 夹具同源构造）→ 拒绝。
- 叶级 reparse：可执行文件本身是 reparse 点 → `openApprovedExecutable` 自检拒绝。
- 变异：① 删祖先循环（只查直接父目录）⇒ 前两反例绿；② 删 `openApprovedExecutable`
  的 reparse 自检 ⇒ 第三反例绿。恢复逐字节一致。

### 2.3 F05：载体事实（身份/平台/解释器）未带到启动边界

**现状句。**
- `internal/console/ai.go:411`：`recordOrConfirmAIInstallSourceProof(installProofPath,
  resolvedExecutable, receiver.ExecutableHash, cli.now())` 只返回 error；载体内
  `fileID / platform / interpreter` 事实停留在 `console/ai_install_proof.go:125`
  `readAIInstallSourceProof` 与 `Matches` 层，不外传。
- `ai.go:118` `aiProbeReceiver` 只含 `Executable / ExecutableHash`；`:536-554`
  `defaultAICheck` 只把摘要传给 `ManagedCopilotCLIProbeExecutor{ExpectedExecutableHash}`。
- `ai_probe_copilot_executor.go:270-445` `Execute`：pre-spawn 的 `requireSameExecutableFile`
  （:61-85）用**本调用内**新开 bound handle 派生 fileID 再自比（approvedFileID 来自 :367 区
  同一 open），**从不与载体记录的 fileID 对证**；平台事实只用于实时不可变判定，不与载体
  platform 对证；解释器事实只验「实时观察成立」，不与载体解释器 (path,hash,fileID) 逐项对证。
- `internal/guard/executable_binding.go:62-94` `verifyApprovedExecutable`：只比对摘要（held
  句柄上重导）——句柄即被执行对象，「执行对象=批准对象」已成立，但**脚本解释器**不在其视野。

**改动点 file:line。**
1. `console/ai_install_proof.go:60` `recordOrConfirmAIInstallSourceProof` 签名改为返回
   `(adapters.AIInstallSourceFacts, error)`：实时观察成功且（无载体 → 落盘后 / 有载体 →
   `facts.Matches(existing)` 收敛后）返回**载体记录的 facts**（它是「批准对象」的规范描述）。
2. `ai.go:118` `aiProbeReceiver` 增字段 `InstallSource *adapters.AIInstallSourceFacts`；
   `executeAICheck`（:411 区）回填；`defaultAICheck`（:536-554）构造
   `ManagedCopilotCLIProbeExecutor{ExpectedExecutableHash, ApprovedInstallSource: receiver.InstallSource}`。
   启动边界完整事实清单（载体五项）：`executablePath / executableHash / fileID /
   platform(os,arch,platformHash) / interpreter(path,hash,fileID)`。
3. `ai_probe_copilot_executor.go`：`ManagedCopilotCLIProbeExecutor` 增字段
   `ApprovedInstallSource *AIInstallSourceFacts`；`Execute`（:367 区）pre-spawn 改为
   **三等式**：实时 `VerifyAIInstallSource` 结果 `live.Matches(*ApprovedInstallSource)`
   （路径+摘要+fileID+platform+immutability+解释器逐项，`ai_install_proof.go:244` 既有
   `Matches` 承担）+ 既有 `requireSameExecutableFile` 用载体 fileID/hash 重验（:61-85，
   approvedFileID/approvedHash 改从载体取，不再本调用自派生）。空 `ApprovedInstallSource`
   走兼容路径（行为同现状）。
4. `adapters/ai_install_proof.go`：新增 `(proof AIInstallSourceProof) Facts() AIInstallSourceFacts`
   构造器（纯字段映射），供 1 的返回与 3 的注入共用。
5. `guard/executable_binding.go:62`：**不改**（摘要-on-held-handle 已是最强绑定；guard 层
   不验解释器是既有分工）。缺口登记：解释器链由 adapters 层 pre-spawn 全链对证承担，§6.1。

**原语选择。** 复用既有 `Matches`（路径平台化比较 + 其余逐字节）、`VerifyAIInstallSource`
（同柄身份+摘要+平台事实）、`requireSameExecutableFile`（pre-spawn 同柄重验）；不新增
原语、不扩 `ProcessSpec`/guard 契约。

**反例与变异测试点。**
- 同摘要异 fileID：cp 生成字节相同的第二 inode，载体记 A、现场换 B → pre-spawn 拒绝
  （现实现放行）。
- 平台事实漂移：载体 platform=linux/amd64、现场 windows（合成 `Matches` 输入）→ 拒绝。
- 解释器漂移：载体记录 /bin/sh(hash A)，现场同路径换文件（异摘要或同摘要异 inode）→ 拒绝。
- 变异：删 `Matches` 的 fileID/platform/解释器任一比较 ⇒ 对应反例绿；删 `Execute` 的载体
  注入 ⇒ 「同摘要异 fileID」反例绿。恢复逐字节一致。

### 2.4 F06：Linux 词法组件 Stat 跟随链接 + 分次按名取证

**现状句。**
- `ai_install_proof_unix.go:86-106` `verifyUnwritableInstallPath`：逐组件 `os.Stat`（跟随
  链接）——组件为 symlink 时检查的是**目标**的模式/属主，名称链真实祖先被漏查。
- `adapters/ai_install_proof.go:257-280` `VerifyAIInstallSource`：`evidence.OpenBoundFile` 的
  unix 打开器是裸 `os.Open`（`boundopen_unix.go:11-13`，跟随链接）；身份、摘要与 shebang
  头部是**分次按名取证**（fileID/hash 一次 open、`installSourceInterpreter` 再 `os.Open`
  读头），两次之间可被换名。
- 既有原语（SEC-02a 家族）：`evidence.OpenBoundDir`（`bounddir_unix.go:96`）+
  `OpenBoundDirFile`（:116）叶级 `O_NOFOLLOW`（P83 fix round 8 实证）。

**改动点 file:line。**
1. `ai_install_proof_unix.go:86-106`：`os.Stat` → `os.Lstat`；组件 `Mode()&os.ModeSymlink
   != 0` ⇒ 拒绝（组件级拒链）。契约句（CONTRACTS:216「全路径各组件不可写」）语义不变，
   只把「组件」钉到**名称**而非链接目标；F08（§3.1）的有效权限改造与本条同区落地，一并改。
2. `adapters/ai_install_proof.go:257-280`：叶打开改为 bound-dir 链——
   `dir := evidence.OpenBoundDir(filepath.Dir(path))` → `file :=
   evidence.OpenBoundDirFile(dir, filepath.Base(path), copilotCLIExecutableMaxBytes)`；
   fileID/hash/头部**同柄**（fileID、Hash、ReadAll 全部走 `file` 这一个 BoundFile）。
3. `ai_install_proof_unix.go:135` `installSourceInterpreter` 拆为两层：
   `installSourceInterpreterFromBound(bound *evidence.BoundFile) (interp, headerClass, error)`
   ——从 `bound.ReadAll()`（受 `installSourceShebangHeaderLength` 界）读头部，shebang/ELF/
   PT_INTERP（§2.5）判定与解释器身份取证全部来自同一对象；解释器自身的打开同样走
   bound-dir 叶级 `O_NOFOLLOW`。保留 `installSourceInterpreter(path)` 薄包装（测试与兼容用），
   生产路径只调 FromBound。
4. Windows 侧同构：`ai_install_proof_windows.go` 叶打开复用 `OpenBoundDirFile`（其 Windows
   实现拒 reparse 叶），与 F04 的组件链检查合为一趟；`installSourceDeniesMutation` 的探测
   对象与 bound-open 的目录是同一 handle 链。

**原语选择。** 组件级拒链 = `os.Lstat`（stdlib，逐组件名）；叶级拒链 = 既有
`OpenBoundDir`/`OpenBoundDirFile`（SEC-02a 家族，Linux 上即 `openat`-相对 + `O_NOFOLLOW`
语义的仓库内实现，不新写 openat2 逻辑）；身份/摘要/头部同柄 = 既有 `BoundFile` 的
「单柄多面」契约（`boundread.go` 头注释）。不新增依赖、不动 `evidence` 包。

**反例与变异测试点。**
- 组件级 symlink：安装目录**中间组件**为软链指向不可写目标 → 拒绝（现实现跟随并通过）。
- 叶级 symlink：可执行文件本身是软链 → 拒绝。
- 换对象窗口：bound-open 与头部读取之间换文件——同柄语义由 `BoundFile` 测试族承接
  （`boundread_test.go`）；本片新增结构性断言：变异 = 把生产路径改回「独立 `os.Open`
  读头」⇒ 红。
- 变异：`Lstat` 改回 `Stat` ⇒ 组件 symlink 反例绿。恢复逐字节一致。

### 2.5 F07：shebang 解析不等价内核 + env 分支 + 动态 ELF PT_INTERP

**现状句。**
- `ai_install_proof_unix.go:135-240`：:169 `fields := strings.Fields(...)` 手写按空白切
  shebang 行；:175-207 `installSourceEnvInterpreter` 对 `env` 形式跳过选项、取首个 token 当
  解释器——**不等价内核**（内核先执行 `/usr/bin/env`，env 再按**用户可控 PATH** 解析；
  且 Linux 把 shebang 整行余部作为**一个**参数传给解释器，与多 token 解析不同）；
  递归解释链（解释器又是脚本）只验一层。
- ELF 分支只验 `\x7fELF` 魔数即返回 nil；动态 ELF 的 `PT_INTERP` 解释器未纳入不可变验证
  与载体记录。

**改动点 file:line。**
1. **env 形式裁定 = fail-closed 拒绝**（选 ①，理由见下）：`ai_install_proof_unix.go:175`
   `installSourceEnvInterpreter` 删除，改为 `filepath.Base(interpreterPath)=="env"` ⇒ 直接
   返回 `ErrAIInstallSourceProof`（「env 形式 shebang 在本规则下不可证明」）。
   理由：① env 的最终解释器选择依赖执行时 PATH（用户可写），不可变规则不覆盖 PATH，
   任何「精确解析」都是再实现内核+env 语义，漂移即 F07 复发；② 探针二进制由本仓库
   发行（ELF），env shebang 是外部夹具形态，拒绝只损失异形布局；③ CONTRACTS:216
   「可执行文件为 ELF 或对其 shebang 解释器同样满足不可变条件」允许对不可证明形式拒绝。
2. **shebang 行解析与内核对齐（能证部分）**：解释器路径必须绝对（既有）；`strings.Fields`
   切割保留（内核即按空白切首 token）；解释器自身再走 `installSourceInterpreterFromBound`
   （§2.4-3）做同构判定——若解释器也是脚本，其 shebang 解释器**同样逐项验证**，递归
   深度上限 = 4，超限拒绝（内核无深度限制，残余登记 §6.4）。
3. **动态 ELF PT_INTERP**：`ai_install_proof_unix.go` 新增
   `installSourceELFInterpreter(bound *evidence.BoundFile) (string, error)`：从同柄读 ELF 头
   （`encoding/binary` 解析 e_ident/e_phoff/e_phentsize/e_phnum，全部界内校验；Malformed ⇒
   拒绝）→ 遍历 program headers 找 `p_type==3`（PT_INTERP）→ 读 NUL 结尾路径（≤ p_filesz）；
   静态 ELF（无 PT_INTERP）返回 nil 直接通过。解释器路径必须绝对（相对 ⇒ 拒绝），随后与
   shebang 解释器同一验证（不可写全路径 + bound 叶级打开 + fileID/hash）并写入载体
   `AIInstallSourceInterpreter`（`ai_install_proof.go:66-72` 结构复用，不新增字段）。
4. pre-spawn：shebang 与 PT_INTERP 两类解释器都在 F05 的启动边界逐项对证范围内
   （§2.3-3 的 `Matches` 已含 `interpreter` 字段，无需再扩）。

**原语选择。** 头部/ELF 解析 = stdlib `bufio` + `encoding/binary`（bounded read，
`installSourceShebangHeaderLength` 界）；解释器不可变验证复用既有
`verifyUnwritableInstallPath` + bound 叶级打开；无新依赖、无新原语。

**反例与变异测试点。**
- env shebang 脚本 → 拒绝（现实现会「解析」并可能放行）。
- 动态 ELF（`/bin/sh` 或 go test 二进制）→ 载体解释器非空且逐项可验证；PT_INTERP 解释器
  位于可写目录（复制 ld.so 夹具）→ 拒绝。
- 静态 ELF（无 PT_INTERP 夹具）→ 通过。
- 递归 shebang：解释器为脚本且链深 > 4 → 拒绝。
- 变异：恢复 env 解析分支 ⇒ env 反例绿；删 PT_INTERP 校验 ⇒ 动态 ELF 反例绿；删 ELF
  头界内校验（数组越界路径）⇒ 畸形头夹具红。恢复逐字节一致。

---

## 3 其余 11 条（逐条裁定 + 要点设计）

### 3.1 F08：`0222` 任何主体写位 → 有效权限 + 属主规则（裁定 = 整改）

- 现状句：`ai_install_proof_unix.go:102` `info.Mode().Perm()&0o222 != 0 || int(stat.Uid) ==
  euid`——任何主体的写位都拒；`0775 root:root` 的 `/usr` 对本用户不可写仍被拒；非常规挂载
  下 `host-auth` 也被挡。模式位也不证明既有写句柄失效。
- 改动点 `ai_install_proof_unix.go:86-106`（与 §2.4 的 Lstat 改造同区）：
  ```
  if stat.Uid == euid { 拒绝 }                      // 属主可 chmod 恢复
  if mode & 0o002 != 0 { 拒绝 }                    // other 写
  if mode & 0o020 != 0 && euid 组（syscall.Getgroups，含 egid）含 stat.Gid { 拒绝 }
  // owner 写位不再单独判定：owner==euid 已拒；owner≠euid 时 0o200 与本用户无关。
  ```
  根 euid 分支（:58-62）保留。
- 登记：既有写句柄失效（同 OS 用户已开写句柄）属 S1 威胁模型外（`SECURITY.md:10`），§6.2
  成文，不得声称可撤销。
- 反例/变异：`0775 root:root` 且用户不在组 → **通过**（新正例，VM 原生非 root 实测）；用户
  所在组有组写 → 拒绝；变异 = 删 `Getgroups` 判定 ⇒ 组写反例绿。

### 3.2 F11：schema 闭集同步 + 协议先行判定失误登记（裁定 = 整改）

- 现状句：`CONTRACTS.md:214` 与 Go（`ai_availability.go:221-226`）已含
  `byok_confinement_unproven`；`schemas/ai-availability.schema.json` `$defs/reason/enum`
  （:18 起列表）仍 17 项（无该值）→ wire 校验拒绝真实记录。Rev 1 §1.5-B 判「schema/夹具：
  不涉」**不成立**。
- **查清生成关系**：仓库**无** schema↔Go 生成器。`tools/contracts/contracts.test.js` 在
  测试期用 AJV 校验 `testdata/contracts/{valid,invalid}/ai-availability.json`；Go 闭集在
  `validAIAvailabilityReason` 手写；两者靠夹具间接对齐。逐字改法：
  1. `schemas/ai-availability.schema.json`：`"timeout",` 行后、`"unknown"` 前插入
     `"byok_confinement_unproven",`（与 CONTRACTS:214 顺序一致）。
  2. `testdata/contracts/valid/ai-availability.json`：追加 fixture
     `ai-availability-byok-confinement-unproven-valid`（status=unavailable、
     reason=`byok_confinement_unproven`、requestsUsed=0、evidence 单条合法 sha256、
     rejectionLayer=schema、expected=accept）。
  3. `testdata/contracts/invalid/ai-availability.json`：追加 fixture
     `ai-availability-reason-out-of-closed-set-invalid`（reason=闭集外值如
     `byok_confinement_unproven_typo`、expected=reject、rejectionLayer=schema）——钉住
     「闭集只增不滥」。
  4. `tools/contracts/contracts.test.js`：`assert.equal(fixtureCount, 144)` → `146`；
     `schemaFiles.length === 38` 不变（无新 schema 文件）；accepted/rejected 覆盖断言由
     扫描自动满足。
- 红绿顺序：先加夹具不动 schema → `npm test` 红（valid 夹具被旧 schema 拒）；再同步 schema
  → 绿；撤 Go 闭集项 → `go test ./internal/adapters/` 红（恢复绿）。
- 判定失误登记文案：见 §5.1（逐字）。

### 3.3 F12：本地拒绝证据纯闭集 vs CONTRACTS:214 绑定要求（裁定 = 整改）

- 现状句：`adapters/ai_preflight.go:169-181` `newLocalAIAvailability` 的 evidence =
  `Digest(localDomain, profileHash+"\n"+channel+"\n"+reason+"\n"+probedAt)`——混入 profile
  摘要与时间；同类拒绝在不同 profile/时间下证据互异，违反 T8b「同类同证」。
  `p85_t8a_tester_test.go:570-590` 只覆盖执行器层同类同证。
- 边界裁定：**profile/time 保留在记录元数据**（`ProfileConfigHash`/`ProbedAt` 是 schema
  必填字段，照常写入）；**证据字段纯闭集分类**（status + reason 两个分类值）。
- 证据构造新算法（逐字）：
  ```go
  // F12: refusal evidence is the pure closed-set classification.
  // Profile hash and timestamp remain record metadata (ProfileConfigHash / ProbedAt).
  evidenceHash := evidence.Digest("proofrail:ai-availability-local-check:1\n",
      []byte(status+"\n"+reason))
  ```
  （channel 亦降为元数据——它已在记录字段；reason 全部来自 `validAIAvailabilityReason`
  闭集，调用点现状已保证。）
- 新测试（③）：两个不同 profile、不同时间的同类 `byok_confinement_unproven` 拒绝 →
  `Evidence` 逐字节相等、`Reason` 相等，且 `ProfileConfigHash`/`ProbedAt` 仍互异（元数据
  保留）；不同类（`cli_unavailable` vs `byok_confinement_unproven`）证据互异。
- 变异：把证据改回含 profileHash/时间 ⇒ 同类同证反例红。
- 影响登记：证据公式收紧只影响**新产生的**本地拒绝记录；schema 不变、旧记录仍
  schema 合法（evidence 是不透明 hash 数组），§5.4 成文。

### 3.4 F13：解析前构造并证明同一不可变启动规格 + argv 白名单（裁定 = 整改）

- 现状句：`ai_probe_copilot_executor.go:242-253` `VerifyBYOKConfinement`（静态半 + 沙箱
  判定）与 :466 起 `assertCopilotCLIProbeBYOKConfinement`（spawn 半）分属两个时刻；
  检查半不绑定具体 spawn 配置。argv 是黑名单（`copilotCLIProbeRefusedLoggingArg`）+
  「自我相等」（`slices.Equal(spec.Args, copilotCLIProbeExecArgs(invocation.Args,...))`）。
  默认 argv（`ai_probe_copilot.go:184` `copilotCLIProbeArgs`）带 `-C <probe.WorkingDirectory>`
  ——即使 `spec.Dir` 被改为 executor 临时目录，`-C` 仍指调用目录；未来 `CheckSandbox`
  变 Allowed 即「先读凭据、spawn 前才发现配置不合」。
- 改动点：
  1. 共享规格构造函数 `copilotCLIProbeLaunchSpec(invocation CopilotCLIProbeInvocation,
     usagePathPattern string, outputLimit int64) (frozenSpec, error)`（`ai_probe_copilot_executor.go`
     内）：纯函数，产出五元模板 = 白名单化 argv 模板、环境键闭集、流绑定类（bounded
     in-memory）、工作目录规则（executor 临时目录）、输出界。**checker 静态证明与 pre-spawn
     五元断言共用此构造函数**。
  2. `ai_probe_copilot.go:103-113` `VerifyBYOKConfinement`（probe 层）：Resolve 前构造并
     冻结该规格（含 argv 白名单校验 + `CheckSandbox`），封存于
     `CopilotCLIProbeInvocation.FrozenSpawn`（新字段）；`Probe`（:114 起）原样传给
     `Execute`。`Execute`（:270 起）不再自算配置，而是实例化后断言「实际五元 ==
     `instantiate(FrozenSpawn, temporaryDirectory, usagePath)`」（逐元素相等），随后 spawn。
     凭据值注入是唯一 Resolve 后变化（`COPILOT_PROVIDER_API_KEY` 值），键集已冻结。
  3. **`-C` 裁定**：secret 模式下 `copilotCLIProbeArgs` **不再输出 `-C`**（BYOK 子进程
     cwd 由 `ProcessSpec.Dir`（临时目录）单一决定）；冻结时白名单要求 credential 模式 argv
     **不含** `-C`/`--working-directory`（含即拒）。F02 的既有断言保留为纵深防御。host-auth
     路径维持 `-C <workingDirectory>` 不变。
  4. **argv 白名单裁定 = 白名单模板 + 参数校验（选白名单）**。理由：黑名单不能证明新
     CLI 版本新增的日志类旗标被拒（闭集之外静默放行）；白名单 = `copilotCLIProbeArgs`
     模板的旗标闭集（`-p/-C/--name/--model/--output-format/--stream/--max-ai-credits/
     --disable-builtin-mcps/--disallow-temp-dir/--no-custom-instructions/--no-ask-user/
     --no-auto-update/--no-color/--no-remote/--no-remote-export/--secret-env-vars/
     --deny-url/--available-tools/--deny-tool`）+ 逐项值类校验（`--model` 为 `ValidID`、
     `--max-ai-credits` 为 `[copilotCLICreditFloor, 上界]` 整数、`--deny-url` ∈ 两定值、
     其余固定常量），未知旗标或值类不符 ⇒ 拒绝；既有
     `copilotCLIProbeRefusedLoggingArg` 黑名单保留为纵深防御（双闸不互替）。
- 反例/变异：argv 增未知旗标（如 `--log-dir-new`）→ 冻结时拒绝；注入 fake 沙箱把 verdict
  改恒 Allowed → **Resolve 之前**即因配置不合拒绝（Resolver 计数 0、`byok_confinement_unproven`）；
  变异 = 删 `FrozenSpawn` 相等断言 ⇒ 该反例绿（凭据被读后拒）。恢复逐字节一致。

### 3.5 F14：Windows 环境键大小写不敏感（裁定 = 整改）

- 现状句：`ai_probe_copilot_executor.go:401-402` `credentialPresent` 按
  `invocation.Environment[copilotCLIProbeCredentialEnvironmentKey]` 精确大小写取值；
  `copilotCLIManagedEnvironment`（:652-690）同样精确匹配继承键。Windows 环境块大小写不敏感：
  调用方以 `copilot_provider_api_key` 传覆盖项时 `credentialPresent=false` ⇒ 整个 BYOK 门禁
  跳过，而子进程仍收到该凭据键（Windows 大小写不敏感解析）。
- 改动点：
  1. 新增平台归一比较 `copilotCLIEnvironmentKeyEqual(left, right string) bool`（Windows
     `strings.EqualFold`，其余逐字节）。
  2. `Execute`（:401-410）：`credentialPresent` 改为扫描 `invocation.Environment` 键、按归一
     比较找凭据键；发现**多个大小写变体**或与 allowlist/BYOK 键仅大小写不同的键 ⇒ 拒绝
     （歧义 fail-closed）。
  3. `copilotCLIManagedEnvironment`（:652-690）：覆盖键与继承 allowlist 键做归一冲突检查
     （Windows 上仅大小写不同的重叠加覆盖 ⇒ 拒绝）；输出只允许单一规范大小写。
  4. `CopilotCLIProbeInvocation`（`ai_probe_copilot.go:27`）新增 `AuthMode string`
     （"host"/"secret"，或等价的 `CarriesCredential bool`）：probe 构造时**显式声明**认证
     模式；`Execute` 断言「声明与凭据键存在性一致」，不一致拒绝——凭据存在性不再由环境键
     推断。
- 反例/变异：Windows 原生用例以 `copilot_provider_api_key` 小写键构造 invocation → 拒绝且
  marker 不落盘（现实现会 spawn）；变异 = 归一比较改回精确比较 ⇒ 该反例绿（spawn 发生）。

### 3.6 F15：console T3 拒绝路径统一结构化原因码 `cli_unavailable`（裁定 = 整改）

- 现状句：`console/cli.go:847-855` `writeCommandError` 只输出自由文本；`executeAICheck` 的
  T3 前置拒绝（`ai.go:411` 区）走同一通用路径；`ErrAIInstallSourceUnproven` 把
  `(cli_unavailable)` 嵌在自由文本（`console/ai_install_proof.go:70` 等）。畸形、读失败等
  路径无结构化分类。
- 改动点：
  1. `cli.go:48-55` `commandResponse` 新增 `ErrorClass string \`json:"errorClass,omitempty"\``
     （**console 局部输出形状**：非 wire 记录、无 schema；登记 §5.3）。
  2. `cli.go:847` `writeCommandError` 增参 `class string`（或新增
     `writeAICheckRefusal(stdout, stderr, command, jsonOutput, class, err)`）；
     `ai check` 各拒绝位点显式传闭集分类：安装来源不可证 / binding / 账本 / 配置错误 →
     `cli_unavailable`（S1 口径：CLI 通道不可用）；BYOK 门禁已由 checker 产出记录
     （reason 字段承担），控制台不重复分类。文本模式输出 `ai check: cli_unavailable:
     <detail>`；JSON 模式输出 `errorClass` 字段。诊断（flag 用法错误）仍走
     `writeUsageError`，不附分类。
  3. 移除 `console/ai_install_proof.go` 等错误文案中的自由文本 `(cli_unavailable)` 嵌入，
     分类只走结构化字段。
- 反例/变异：console 包测试断言 JSON 输出含 `"errorClass":"cli_unavailable"` 且
  `resolveCalls==0`、checker 零进入；变异 = 删 class 传递 ⇒ 断言红。

### 3.7 F21：Linux fd 3 vs 契约「仅标准流句柄」+ 平台桩（裁定 = 整改）

- **裁定 = 契约补例外措辞**（非实现注释登记）。理由：fd 3 就是 T3 绑定机制本身
  （`guard/process_linux.go:91` `commandPath = "/proc/self/fd/3"`、:101 `cmd.ExtraFiles`），
  BYOK 契约句「仅继承标准流句柄」（`CONTRACTS.md:220`）字面对**每个**经批准可执行文件
  启动的 BYOK 子进程恒假——明知字面为假而不改契约，等于在契约里留一枚常驻违约；而 fd 3
  不含任何未验证对象（就是被哈希、被执行的同一文件描述符），写明不削弱任何保证。
- 改动点：
  1. `CONTRACTS.md:220`：「仅继承标准流句柄」改为「仅继承标准流句柄与经验证的可执行文件
     描述符（Linux 上为 `/proc/self/fd/3` 所引的已验证 fd；不得传递任何其它句柄）」；
     `CONTRACTS_EN.md:214`（据 ⑤r2 锚点，落地时以磁盘实况复核）同位镜像。协议先行：
     契约改 → 无 schema/夹具影响（句内措辞）→ 代码零改。
  2. 平台桩：`ai_install_proof_unix.go:1` 的 `//go:build !windows` 拆为
     `//go:build linux`；新增 `ai_install_proof_unsupported.go`
     （`//go:build !windows && !linux`）：`verifyInstallSourceImmutability` 恒返回
     `ErrAIInstallSourceProof`（含 `syscall.Stat_t` 引用的代码只在 linux 编译）；
     `runtime.GOOS != "linux"` 的运行时拒绝（:52-56）保留为第二道防线。
- 反例/变异：非 linux 平台 `go build ./...` 通过 + 桩函数单测恒拒；变异 = 桩返回 nil ⇒ 红。

### 3.8 F22：ci scope 口径（裁定 = 登记）

`05-input/gate-ci.txt`（v69 快照、`scope=ci changed=10`）与 `selftest.txt`（740/740）是
**步骤 0 锚点读数**，不是 P85 双平台行为记录；本片包内无原生 Go 测试/lint/变异执行日志。
登记口径与 ⑨ 行为证据清单见 §5.2。

---

## 4 ③ 测试要点（F16 / F17 / F18，逐条一句话设计）

- **F16（宿主绝对路径）**：共享解码正例（`adapters/p85_t3_tester_test.go:120` 区）把
  盘符路径改为宿主绝对路径构造 `filepath.Join(string(filepath.Separator), "opt",
  "proofrail", "copilot")`（两平台 `IsAbs` 恒真、`Validate` 通过），或把 Windows 盘符正例
  移入 `//go:build windows` 专属测试；断言不变（进入变异断言前不再失败）。
- **F17（实际可启动正对照 / 逐项分腿）**：① 正对照用真实可启动对象——Windows 用
  `%SystemRoot%\System32\cmd.exe`（`p85T8ASystemShell` 已证可证）、Linux 用复制 `/bin/sh`
  的 ELF（可写目录夹具改用它，让「目录守卫」单独被测量，不再被非 ELF 拒绝掩盖）；
  ② 逐项分腿：部分 ACL 单权利授予（F03 夹具）、同路径换对象（fileID 变化）、同 inode
  改写（维持 `EXEMPT-REGISTERED P85/F1`，四要素更新为本版新字段）、解释器漂移（同路径
  替换解释器字节）；③ Windows 启动标记改 `cmd /c echo ran> marker`（`p85T8AMarkerArgs`
  形态），证明撤守卫后**原生**创建成功（弃 `.cmd` 脚本）。
- **F18（T2 撤全部制品重读 / T8a 独立启动观察点）**：① T2 变异升级为「保留账本校验、
  撤除 Publish 路径的**全部**制品重读」——变异点 = `verifyFrozenPolicyArtifact` 的磁盘重读
  （含经 `operatorApprovedPolicyHashForAttempt` 的 grant-helper 路径，`runtime_agent.go:1116`
  起），「Review 后制品被替换」反例必须在此变异下转绿；② T8a 执行器门禁变异改用
  **独立可控启动观察点**：测试替身 executor 保留真实 BYOK 断言链（或反之）直测
  `assertCopilotCLIProbeBYOKConfinement`，脱离宿主 T3/sandbox 条件的跳过；正对照复用
  `p85T8ASystemShell`（可证主机上）；③ 提交逐项红、恢复绿及字节复原证据（变异前后
  `git diff` 对实现文件零增量）。

---

## 5 报告登记要点

### 5.1 F11「协议先行判定失误」登记文案（逐字，进报告）

> **登记：协议先行判定失误（P85 Rev 1 §1.5-B）。** 事实：Rev 1 判定 T8a「schema / 夹具：
> 不涉」，而失败原因闭集是 wire Schema `ai-availability.schema.json` 的 `reason` enum
> 成员；契约闭集（CONTRACTS:214）与 Go 闭集增 `byok_confinement_unproven` 后，schema 与
> 夹具未同步，真实记录会被 wire 校验拒绝。四要素：① 判据 = 协议先行要求「契约闭集增项
> 必须同步 schema + 契约夹具」（DELIVERY_DIRECTIVE §1.4）；② 偏离 = 只在 CONTRACTS 双语
> 与 Go 闭集增项，未改 schema/夹具；③ 纠正 = Rev 2 §3.2 同步 schema 闭集 + 增 valid/invalid
> 夹具 + 更新 `contracts.test.js` 计数；④ 声明 = 本项系 Rev 1 判定失误的直接纠正，
> **不得作为先例**放宽后续切片「schema 影响」判定——后续凡触失败原因/枚举/闭集即按
> 协议先行全链同步。

### 5.2 F22 口径登记 + ⑨ 行为证据清单（逐字，进报告 §8）

> **口径登记：** `05-input/gate-ci.txt`（`scope=ci`，v69 快照锚点）与 `selftest.txt`
> （740/740）为步骤 0 门禁读数，**不构成** P85 双平台行为证据；本片行为证据以 ⑨ 的
> 下列运行留痕为准：
> 1. Windows 本机：`go test ./internal/adapters/... -run P85 -count=1`（含 windows build
>    tag 腿：F03 单权利夹具、F04 junction、F05 同摘要异 fileID、F14 大小写腿、F17
>    cmd.exe 正对照）红绿留痕；
> 2. Ubuntu VM（10.0.0.199，go1.22.12）：`go test ./internal/adapters/... -run P85` linux
>    腿（F06 Lstat 拒链、F07 PT_INTERP/env/静态 ELF、F08 有效权限正反、F16 宿主绝对路径
>    正向）红绿留痕；
> 3. `go build ./...` 与 `go vet ./...` 双平台通过留痕；
> 4. `tools/contracts`：`npm test`（F11 先红后绿两段留痕）；
> 5. 变异证据：§2/§3 每条改动点的撤守卫红运行 + 恢复后 `git diff` 对实现文件零增量、
>    测试逐字节一致；
> 6. 门禁摘要（gate/index/tree 三作用域）仅作收口读数，不作行为证据替代。

### 5.3 F15 控制台输出形状扩展（登记）

`commandResponse` 增 `errorClass`（`omitempty`）字段：console 局部输出形状，非 wire 记录、
无 schema、不影响 `ai-availability` 记录结构；旧调用方不感知（省略即空）。登记以避免
后续审查误读为 wire 变更。

### 5.4 F12 本地拒绝证据公式收紧（登记）

本地拒绝记录（`newLocalAIAvailability` 产出）的证据摘要输入由
「profileHash+channel+reason+probedAt」收紧为「status+reason」纯闭集：仅影响**新产生的**
记录；schema 与 Go 校验（`validSortedUniqueHashes`）不变，旧记录仍合法；`ai verify` 准入
不依赖证据具体值。属语义收紧（T8b 类纪律落地），非 schema 变更。

### 5.5 契约改动清单（协议先行顺序）

| 项 | 文件 | 性质 |
|---|---|---|
| F11 | `schemas/ai-availability.schema.json` + `testdata/contracts/{valid,invalid}/ai-availability.json` + `tools/contracts/contracts.test.js` | 闭集枚举同步（非新字段） |
| F21 | `CONTRACTS{,_EN}.md`（CN :220 / EN :214 区） | 一句例外措辞 |
| F12/F15 | 无契约改动 | 实现语义收紧 / console 局部形状 |

顺序：契约与 schema/夹具先落地（②a），再代码（②b），再测试（③）。

---

## 6 残余风险登记（报告 §8 成文）

1. **内核名称解析窗口**（Windows `CreateProcessW` 重新解析整条名称链；Linux shebang 与
   `PT_INTERP` 解释器由内核按名解析的瞬间）：由「全链 ACL/模式判定（F03/F04/F06/F08）
   + 逐组件拒 reparse（F04/F06）+ held deny-write/delete 句柄 / fd-exec（既有 guard）
   + 解释器逐项对证（F05/F07）」封堵；解析瞬间的 TOCTOU 残余无法由用户态消除，**如实
   登记，不声称沙箱**。
2. **同 OS 用户恶意进程与既有写句柄**：同用户进程的主动写入、已打开写句柄的失效
   （F08）属 S1 威胁模型外（`SECURITY.md:10`）；本片防线 = 写者路径封堵 + no-replace +
   启动前重读/重验对账，被篡改会在对账处转红并阻断。
3. **env-shebang 一律拒绝的可用性成本**：发行物为 ELF，无实际影响；登记为 fail-closed
   决策（F07）。
4. **解释器递归链验证深度上限**：内核无深度限制；验证至载体记录层 + 每层不可写判定，
   超限拒绝；深层递归残余登记（F07）。
5. **载体格式演进**：F03/F04/F06/F07 改变平台事实与证据输入后，旧载体
   （`ai-install-source-proof.json`）与现场观察不再匹配 ⇒ `recordOrConfirm` 对旧载体
   fail-closed 拒绝（`cli_unavailable`），操作员重新跑 `ai check` 重建；不静默迁移、
   不降级放行（store-local，无 wire 影响）。
6. **F15 `errorClass` 仅覆盖控制台前置拒绝**：checker 内拒绝仍以记录 reason 为准，两处
   分类来源不重复、不冲突。

---

## 7 硬约束核对（不扩域 / 边界不碰）

- **无新依赖**：全量 stdlib（`syscall`/`encoding/binary`/`bufio`/`slices`）；无第三方包。
- **无新 schema 字段**：F11 仅为闭集枚举成员同步（`reason` enum 增一项）；无新文件、无
  新字段、无新 record。
- **无新 CLI 命令**：F15 只扩 `commandResponse` 可选输出字段。
- **不碰已闭合面**：S1 Exit 八门判据、`SECURITY.md` 威胁模型基线、F1/F2/F4/F5/F7 已闭合面
  （`internal/chain/review.go`、`runtime_supervised.go`、`guard` 契约、`ProcessSpec`）零改；
  F05 明确不改 guard 层语义。
- **登记而非声称**：§6 全部残余项一律「登记」，不得表述为沙箱或可证明无窗口。
- **载体与记录边界**：`ai-install-source-proof.json` 仍为 store-local（路径派生、无 schema、
  无 wire 记录）；`agent-runner-policy/` 制品同理（单1 已修，不动）。

---

*落盘：`docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/01-design-rev2.md`（UTF-8 无 BOM + LF）。*
*下一环节：② 按 §5.5 顺序落地（②a 契约/schema/夹具 → ②b 代码 → ③ 测试）。*
