# P85 `SECURITY-PROOF-SUBSYSTEM` · ① 架构增补方案书 Rev 4（single-round · 身份统一五项）

- 切片：P85 · ① 架构师 = deep-reasoner（DeepSeek V4 Pro）；本文件为 **single-round 增补设计**，
  只做设计、不实现代码，**只覆盖「身份统一」五项**，不重复 rev3 其余条目。
- 承接输入：⑦ 复审 r2 原文（`04-reviews/07-re-review-r2.md`：NOT-CLOSED 四组 + 三条新发现
  + ⑧/⑨ 指令）+ `01-design-rev3.md` + 磁盘现状。
- 行号口径：**以 2026-11-23 磁盘现状为准**（本轮实测）。r2 的行号系其所审补丁快照，与本版
  不一致处在各条注明（如 r2 :451 ≈ 现盘 :482；r2 :93 ≈ 现盘 bounddir_windows.go:95）。
- 编码：UTF-8 **无 BOM + LF**（tmp 工作区文件）。

## 五项共用的底座（每项直接引用，不重复展开）

1. **取证句柄链** `agentRunnerPolicyDomainChain`（console，§1 定义）：目录/父目录两个
   `evidence.BoundDir`（各带 FileID）+ 制品/锚两个 `evidence.BoundFile`（各带句柄身份），
   一次绑定、全程同一批句柄，任何阶段之后可 `ReassertIdentities()`（重取 Fstat 身份逐一对证）。
2. **四槽固定域摘要**（adapters，§2 定义）：写守卫观察向量形状固定为
   `[策略目录 | 父目录 | 制品 | 锚]` 四槽，空槽写规范 `∅` 行；摘要域常量升级为
   `managedAgentPolicyWriteProofDomain = "proofrail:agent-runner-policy-write-guard:2\n"`。
3. **Windows 根名称链拒链**（evidence，§3 定义）：绑定目录路径自身的每个组件按重解析点
   分类，junction/symlink 一律拒绝。

---

## 1 取证句柄链贯穿（r2 新 High :558；F02/F06/F20 的「取证对象 ≠ 读写对象」残余）

**现状句。**
- `agent_runner_policy.go:558`（load 分支）制品经 `readAgentRunnerPolicyBytes(path)` →
  `openAgentRunnerPolicyFile`（:760）→ `OpenBoundDirVerified` **绑定 #1，函数返回即 Close**；
  `:600` 锚经 `loadAgentRunnerFrozenPolicyAnchor` → `readAgentRunnerPolicyBytes` **绑定 #2**；
  `:589` → `agentRunnerVerifyFrozenPolicyAnchor`（:623）→ `:646`
  `agentRunnerReverifyManagedWriteProof(policyDirectory, existing, recorded)` →
  `agentRunnerObserveManagedWriteProof` = `adapters.VerifyManagedAgentCannotWritePolicyDomain`
  ——**按名重新观察**（Linux `os.Lstat`；Windows 逐路径 CreateFile），观察所得身份（dev:ino /
  卷序列号+文件索引）与绑定 #1/#2 的句柄身份**从不比较**。
- `publishAgentRunnerPolicyFile`（:817）再次独立 `OpenBoundDirVerified`；Windows 半壁
  `agentRunnerAssertDirectoryIdentity`（publish_other.go:100）用 name-based `FileIDAt` 复核。
- 即：取证（写域观察）、读（制品/锚字节）、写（发布）三段的目录身份各自独立获取，互不
  对证——rev3 §1/§2 只修了「每段内部的同柄」，未修「跨段同一对象」。

**根因。** 各阶段按名重获身份；「观察判定可写性的对象」与「字节实际来自/写入的对象」之间
没有 FileID 等式，替换保持名称一致即可通过全部检查。

**新原语与调用序。**
1. console 新结构（替代 load/verify/publish 三处的独立绑定）：
   ```text
   bindAgentRunnerPolicyDomain(runRoot, taskID, stepID, attempt) (*agentRunnerPolicyDomainChain, error)
     dir, dirID      := evidence.OpenBoundDirVerified(policyDirectory)          // 绑定 #1，仅此一次
     parent, parentID:= evidence.OpenBoundDirVerified(filepath.Dir(policyDirectory))
     artifact, artifactID := dir 下 OpenBoundDirFile(base, bound) + BoundFile.FileID()   // 仅当存在
     anchor, anchorID     := dir 下 OpenBoundDirFile(anchorBase, bound) + BoundFile.FileID()
     ReassertIdentities(): VerifyBoundDir(dir,dirID) / VerifyBoundDir(parent,parentID)
                           / artifact/anchor FileID() 重取 == 记录值
   ```
2. adapters 新句柄形式原语（name-based 入口保留给无句柄调用点）：
   `VerifyManagedAgentCannotWritePolicyDomainHandles(dir, dirID, parent, parentID, files []*BoundFile) (string, error)`
   - Linux：每对象 `Fstat(fd)` 五元组（dev:ino:mode:uid:gid）+ 有效写权判定；ACL 检查保持
     path 形式（`installSourceDenyAccessACLPath`），按「身份三明治」执行：Fstat(id₀) → ACL
     → Fstat(id₁)，id₀≠id₁ 即拒；向量记 id₀。
   - Windows：目录身份取 `dir.FileID()`、文件取 `BoundFile.FileID()`；逐权利探测保持
     path 形式，同样 id₀→探测→id₁ 三明治。
3. 调用序迁移：
   - `loadAgentRunnerFrozenPolicy`：bind 一次 → 制品/锚经 chain 句柄读取 → `ObserveDomain`
     → `ReassertIdentities` → 才解码/采纳；
   - `publishAgentRunnerPolicyFile`（:817）签名改 `publishAgentRunnerPolicyFileFrom(bound,
     fileID, path, payload, label)`，发布后 `ReassertIdentities`（Linux 半壁
     publish_linux.go 已句柄内操作，仅入口去二次绑定；Windows 半壁的 name-based
     `FileIDAt` 复核按 §3 改为根链重分类复核）；
   - `agentRunnerVerifyFrozenPolicyAnchor` 签名收 chain（替代 `policyDirectory string,
     existing []string`），内部观察走 `VerifyManagedAgentCannotWritePolicyDomainHandles`。

**TOCTOU 论证。** 目录/父目录/制品/锚四个对象在链存续期内身份全部由各 fd 的 Fstat 派生；
观察、读取、发布、复验全经这些 fd；每步后 `ReassertIdentities` 重取比较 ⇒ 任一对象在任两
步之间被替换 ⇒ 终态身份变 ⇒ 拒绝。名称只在 bind 时刻解析一次（Linux openat2 单次解析；
Windows 见 §3 根链）。「观察判定对象」与「字节读写对象」是同一批句柄 ⇒ 「取证对象 ≠ 读写
对象」不存在。Linux 身份三明治的「换入换回」残余由非受管身份边界封闭（同 rev3 §1 论证：
能换目录的对象本身已使写域证明失败）。

**反例设计（③，M13）。**
- 正例：合法布局（非 root 运行）首发+重入绿。
- 反例 1（跨段换对象）：读制品后、读锚前 rename 掉 policy 目录、同名放替身（含替身锚）→
  正版因终态 Reassert 拒绝。变异：锚读取退回独立 `readAgentRunnerPolicyBytes(path)` →
  变异版绿。
- 反例 2（观察/读写错配）：`ObserveDomain` 目录身份改回 `FileIDAt(name)` → 构造「bind D1 →
  rename D1、D2 就位 → 观察按名得 D2 → 读取经 D1 句柄」→ 正版红 / 变异版绿。
- 单变异红：删 `ReassertIdentities` 调用 ⇒ 反例 1 转绿。

---

## 2 首发/重入摘要域统一（r2 Medium :451 ≈ 现盘 :482；消除「未变化首发锚被误拒」）

**现状句。**
- 创建分支 `agent_runner_policy.go:482` 封存**短域**（目录+父目录）摘要；复用分支 `:422`
  封存**全域**（+制品+锚）。封存域形状由「写锚时存在什么」隐式决定，锚内**不记录**。
- `agentRunnerReverifyManagedWriteProof`（:671）以双候选 `[nil, existing]` 猜域：若改「只按
  当前全集重建」，短域首发锚逐字节比较必然失败 ⇒ 未变化首发锚被误拒（:451 语义）；若保
  留短域候选，又回到 §5 的捷径。两者不可兼得，因为重入不知道锚封存的是哪个域。

**根因。** 摘要域形状不固定、不显式声明；封存域与重入重建域由调用点猜。

**新原语与调用序（采纳用户第二方案「空集固定域」为主、第一方案「同域对证」为比较规则，
并加显式 Shape 字段消除歧义）。**
1. 四槽固定域（底座 2）：`managedAgentPolicyDomainWriteProof` 总是产出四槽向量
   `[目录 | 父目录 | 制品 | 锚]`，不在 `existing` 内的槽写规范 `∅` 行；`existing` 排序规则
   不变；观察身份来自 §1 句柄链。摘要域常量升级 `:2`（向量形状变更 ⇒ 域版本随升，杜绝
   跨版本摘要碰撞）。
2. 锚结构 `agentRunnerFrozenPolicyAnchor` 增字段 `ManagedWriteProofShape string`
   ∈ {`"creation"`, `"reuse"`}，由写锚分支确定（复用分支恒 `"reuse"`）。**缺失（旧锚）即
   拒绝**，不猜测——S1 生产路径写域证明必失败、装配在写锚前即拒绝，磁盘不存在存量生产锚，
   无迁移负担；测试夹具随实现同步。`agentRunnerFrozenPolicyAnchorPayload`（:487）同步序列化。
3. 重入规则（`agentRunnerReverifyManagedWriteProof` 重写，与 §5 共用）：
   a. 按 `Shape` 重建**封存域**（四槽，槽位按 Shape 填 `∅`）→ 摘要与
      `anchor.ManagedWriteProof` **逐字节比较**，不等即拒。未变化首发锚：Shape=creation ⇒
      重建短域 ⇒ 相等 ⇒ 通过，**不再误拒**；
   b. 随后对**当前全部已存在对象**（目录+父目录+当前在盘制品+锚）强制验权（§5），必须
      成功，不与被封存摘要比较。

**TOCTOU 论证。** 域形状不再由「当时存在什么」推断：锚显式声明、重入按声明重建；封存摘要
与重建摘要对同一槽序、同一批句柄身份 ⇒ 对象未变则逐字节相等；对象已变 ⇒ 要么身份变化改
摘要（不等 ⇒ 拒），要么当前全集验权失败（拒）。两条路径覆盖「首发锚 + 制品/锚后来出现且
被改写」与「复用锚 + 任一对象被换」全部组合，且无捷径、无误拒。

**反例设计（③，M15）。**
- 正例 1：Shape=creation 首发 → 无注入重入（制品/锚在盘、未被改、写权封闭）→ 绿；断言实现
  内「同域对证 + 全集验权」双真（变异见下）。
- 正例 2：Shape=reuse 正常重入 → 绿。
- 反例：creation 锚改写成 reuse 摘要（或反之）→ 拒；旧锚无 Shape 字段 → 拒。
- 变异：删全集验权 ⇒ 构造「目录/父目录不可写、制品/锚可写且被改写」→ 变异版绿 / 正版红。

---

## 3 Windows OpenBoundDirVerified 根名称链（r2 Medium :93 ≈ 现盘 bounddir_windows.go:95）

**现状句。**
- `bounddir_windows.go:95 OpenBoundDirVerified` → `OpenBoundDir`（:58）= 普通 `CreateFile`，
  无 `FILE_FLAG_OPEN_REPARSE_POINT`：目录路径**自身组件**的 junction/symlink 被内核静默跟随；
  相对子项的重解析分类（`openDescendantRelative` :687 逐组件 + `confirmOpenedArtifact`）只
  覆盖**绑定目录以下**。句柄身份与 name-based `FileIDAt` 都命中跟随后的对象 ⇒ 一致而绿。
- 已登记的 F04 CreateProcessW 最终名残留（spawn 名）不覆盖取证绑定名（r2 明确）。

**根因。** Windows 绑定原语不对「目录路径自身组件」分类重解析点；该缺口既不在 F04 残留
范围，也不被相对子项分类覆盖。

**选型：拒链**（不按 F04 残余扩展登记）。理由：① 冻结策略域是 P85 的安全锚点，绑定对象
身份是后续全部对证之根；② evidence 已具备逐组件相对打开+分类机制（`openRelativeComponent`，
`openDescendantRelative` 所用），根链行走只是把「卷根→目录」的组件走一遍，无新机制；
③ 再扩一条 F04 豁免会让「取证绑定」这个最内核的位置长期依赖豁免，r2 已言明既有 F04 不能
替这里背书；④ fail-closed、无兼容成本（生产路径不含重解析组件）。备用登记文案（不采用）：
「Windows 取证目录绑定的根名称链重解析点由内核跟随，绑定身份与 FileIDAt 同对象一致，属
F04 残余扩展：绑定名≠spawn 名，取证读取不产生执行面」——因 ②③ 弃用。

**新原语与调用序。**
1. evidence（bounddir_windows.go）：
   `OpenBoundDirVerified(path)` 重写为根链行走 `openBoundDirChainVerified(path)`：
   a. 组件切分（复用 `installSourceWindowsComponents` 同构逻辑，卷根 `C:\` 起）；
   b. 卷根组件按现 `OpenBoundDir` 打开（卷根是挂载点但非 NTFS 重解析标记，可由
      `IsFixedVolumeRoot`/`ntfsRootRecordOf` 锚定，不拒）；
   c. 其余每组件 `openRelativeComponent(parent, component, true, 以 FILE_OPEN_REPARSE_POINT
      打开为自身)`：`FILE_ATTRIBUTE_REPARSE_POINT` ⇒ **拒绝**（junction/symlink/其它标记
      一律拒，不跟随、不比目标——根链无「站内」可证，与 `openDescendantRelative` 对站内
      中间组件的宽容策略不同）；非目录/不存在 ⇒ 拒绝；
   d. 末组件句柄即 BoundDir 句柄，`FileID()` 从句柄读取（现行为保留）；中间句柄逐一释放。
   - UNC/相对路径 ⇒ 拒绝（fail-closed）。
2. `publishAgentRunnerPolicyFileIn`（publish_other.go，Windows 半壁）的
   `agentRunnerAssertDirectoryIdentity` 中 name-based `FileIDAt(directory)` 复核改为
   `openBoundDirChainVerified(directory)` 重走根链 + FileID 与绑定值比较（发布前后各一次，
   复核与句柄同源）。
3. adapters Windows 半壁 `managedAgentPolicyDomainObservation`
   （agent_runner_policy_write_guard_windows.go:76）的 `installSourceComponentIsReparse`
   只分类末组件，升级为对每个对象的路径走同一根链分类器（evidence 暴露
   `BoundDirRootChainClassify(path)` 或 adapters 内同构实现），使「观察判定」与「绑定」对
   同一名称链给出同一结论。

**TOCTOU 论证。** 绑定 = 逐组件以「打开为自身」获取句柄链，最终句柄是「原串组件序列各自
为真实目录、无重解析」的产物；相对子项操作只发生在该对象内。根链组件被换成重解析点 ⇒
下一次绑定/复核分类即拒；被换成普通目录 ⇒ 终态 FileID（卷序列号+文件索引）变化 ⇒
`VerifyBoundDir`/`ReassertIdentities` 拒绝。F04 的 CreateProcessW 残留（spawn 名）与取证
绑定不再共享任何名称解析步骤。

**反例设计（③，M14，Windows 原生）。**
- 正例：无重解析的本地固定卷路径绑定绿（回归护栏）。
- 反例 1：runRoot 组件换成 junction → 绑定拒绝（现实现绿）。恢复 → 绿。
- 反例 2：发布前复核把 name-based `FileIDAt` 换回 → 「绑定后、发布前把目录换成同名
  junction」变异版绿 / 正版（根链复核）红。
- 变异：根链行走删「中间组件重解析拒绝」⇒ 反例 1 转绿。夹具需非提升令牌建 junction，
  不可建则 Skip（写明）。

---

## 4 M8 判别性修正（p85_t3_reentry_tester_test.go:247，M8 leg 2）

**现状句。**
- 生产缝（`ai_install_proof.go:139`）：
  `if syncErr := aiInstallSourceCarrierDirectorySync(bound); syncErr != nil && aiInstallSourceReentryRequiresSync { 拒绝 }`
  ——`requiresSync=false` 时同步失败不拒。
- M8 leg 1：首轮 sync 失败 → 拒；`failSync=false` 后重入 → 采纳（翻的是失败注入，守卫未动）。
- M8 leg 2（:242-260）：首轮 `failSync=true` 拒；随后 **failSync=false 且 guard=false 双翻** →
  采纳。**无「守卫开 + 同步持续失败 ⇒ 红」断言** ⇒ 删掉实现中的重入 sync（或删守卫条件），
  现有两腿均仍绿（r2 F10 反例 NOT-CLOSED）。

**根因。** 反例未把「失败注入」与「守卫」正交化：要证明守卫起作用，必须让失败注入保持
失败、只翻守卫一位。

**新原语与调用序（测试缝语义定稿，实现语义已满足、不要求改实现）。**
1. 测试缝语义（写为注释/文档，⑧ 落地核对）：
   - `aiInstallSourceCarrierDirectorySync` 只注入「同步调用结果」，不改变控制流；
   - `aiInstallSourceReentryRequiresSync` 只控制「重入分支是否把 sync 失败转为拒绝」，不影
     响首轮发布路径；
   - 生产语义恒为 `syncErr != nil && requiresSync ⇒ 拒`，守卫是唯一被翻动的门。
2. M8 重写为三腿（leg 2 替换）：
   - **leg-R（判别红）**：`failSync=true` 保持、`guard=true` → **重入须拒**，错误含
     `p85F10SyncFailureMarker`；
   - **leg-G（判别绿）**：**`failSync` 仍 = true**、只翻 `guard=false` → **重入须采纳**，且
     采纳字节与首发 carrier 逐字节一致；
   - leg-P（正向对照，保留现 leg 1）：首轮 sync 失败 → 拒；`failSync=false` + `guard=true`
     → 采纳。
3. 变异检查（单变异必红）：① 删实现中重入 sync 调用 ⇒ leg-R 绿；② 删 `&& requiresSync`
   条件（sync 失败恒拒）⇒ leg-G 红。双平台性：seam 腿平台无关（`p85F10SeamExecutable`
   已平台独立），原生耐久腿仍 Linux-only。

**TOCTOU 论证。** 本项无新增 TOCTOU 面（测试缝不改生产控制流）；论证点是判别性完备：
leg-R/leg-G 在 failSync 恒真前提下只改变守卫一位，任一实现删除「同步」或「守卫」都会令恰
一条腿翻转 ⇒ 反例对两处实现回归都具有判别力。

**反例设计（③，M16）。** 即 leg-R / leg-G / 变异①②，如上。

---

## 5 移除 Publish 短域匹配捷径（agent_runner_policy.go:688）

**现状句。**
- `agentRunnerReverifyManagedWriteProof`（:671-696）双候选 `[nil, existing]`，`:688`
  `if observed == recorded { return nil }` ——封存短域摘要的锚在 load/Review/Publish 时被
  `nil` 候选直接命中返回，**制品/锚当前写权不被验证**（rev3 §6 的对证被候选捷径架空）。
- 复用分支锚的载荷由制品+锚组成，锚自身也在同一域内——域不封闭时整对文件可被受管身份
  同步改写而不改变「短域摘要对证」结果。

**根因。** 候选集是「封存时刻可能存在的域」的枚举而非「当前存在对象的全集」；短域候选使
封存域=短域的锚跳过全集验权。

**新原语与调用序。**
1. 删除双候选循环；`agentRunnerReverifyManagedWriteProof` 重写为两步（与 §2 共用）：
   a. 按锚 `ManagedWriteProofShape` 重建封存域（§1 句柄链、§2 四槽）→ 与
      `anchor.ManagedWriteProof` 逐字节比较（同域对证）；
   b. **无条件**执行当前全集强制验权：`existing` = 当前在盘的制品+锚（load 已确认存在），
      观察必须成功；任一对象可写/不可判定 ⇒ 拒绝。不做摘要比较、不允许任何候选短路。
   等价表述：移除「封存域命中即返回」路径——全集合强制验权在每次 load/Review/Publish 都
   执行。
2. **性能可接受性说明**：对象数上界恒为 4（目录/父目录/制品/锚）。Linux：4×（Fstat + ACL
   读 + 权判定），约为一次锚载荷解码的零头。Windows：4 对象 × 每对象 ≤7 次权利探测
   （目录表 7 项、文件表 ≤7 项，`installSourceDirectoryMutationRights` /
   `installSourceFileMutationRights`）≈ 至多 28 次「开句柄即关」的 CreateFile 探测，单次
   <1ms；且 load/Review/Publish 均为低频人工节点，不在任务热路径。结论：开销可忽略，不
   需要缓存或摘要协商。
3. **Windows/Linux 差异登记**（写入观察语义注释 + SECURITY/CONTRACTS 同步项）：
   - Linux：mode/owner/group + POSIX access-ACL + dev:ino；root 直接拒绝（无法证明）；
     观察与读写经同一 openat2 句柄链（§1）。
   - Windows：逐权利探测 + 卷序列号/文件索引 + 每对象根链重解析分类（§3）；无 mode/owner
     规则（以 DACL 探测替代）；无目录句柄 fsync（F04/F10 已登记：原子可见 ≠ 断电持久）；
     卷根「排除 FILE_ADD_SUBDIRECTORY/ADD_FILE」的 F04 裁定**不适用于**策略域观察（策略
     域不是卷根）。

**TOCTOU 论证。** 全集验权不再被摘要比较短路 ⇒「锚摘要合形」不再足以通过；观察对象身份
来自 §1 句柄链、终态 `ReassertIdentities` 封闭交换窗口；封存域对证只决定「锚是否本域所
写」，全集验权决定「当前是否仍封闭」，两者独立且都必须为真。删除捷径不引入新窗口。

**反例设计（③，M17）。**
- 正例：creation 锚 + 当前全集封闭 → load/Publish 绿（短域锚不再误拒，接 §2 正例 1）。
- 反例 1（钉 :688 捷径）：creation 锚 + 制品/锚当前可写（目录/父目录不可写）→ 拒绝
  （现实现绿）。
- 反例 2：reuse 锚 + 制品当前可写 → 拒绝。
- 变异：恢复「短域命中即返回」⇒ 反例 1 转绿；「全集验权跳过」⇒ 反例 2 转绿。

---

## 6 ③ 测试增量（M13–M17，编号承接 rev3 M1–M12；M8 由 M16 重写承接）

| # | 对应项 | 测试（包/平台） | 反例/夹具 | 变异（红） | 恢复（绿） |
|---|---|---|---|---|---|
| M13 | §1 | console（linux 原生；seam 腿双平台） | 读制品后 rename 目录放替身（含替身锚）→ 正版拒；锚读取退独立 bind 变绿 | 删 `ReassertIdentities` 或退 path 读 | 恢复 |
| M14 | §3 | evidence（windows 原生） | runRoot 组件为 junction → `OpenBoundDirVerified` 拒（现绿）；发布前复核退 name-based 变体绿 | 删根链中间组件重解析拒绝 | 恢复 |
| M15 | §2 | console（双平台，注入 seam） | Shape=creation 首发 → 无注入重入绿；Shape 改写 → 拒；缺 Shape 旧锚 → 拒 | 删全集验权（「制品/锚可写且被改」转绿） | 恢复 |
| M16 | §4 | console（seam 腿双平台；原生腿 linux） | leg-R：failSync=true+guard=true 重入拒；leg-G：failSync=true+guard=false 重入采纳且字节一致 | ①删重入 sync ⇒ leg-R 绿；②删守卫条件 ⇒ leg-G 红 | 恢复 |
| M17 | §5 | console（双平台） | creation 锚+制品/锚可写 → load/Publish 拒（现绿）；reuse 锚+制品可写 → 拒 | 恢复短域捷径 ⇒ 转绿 | 恢复 |

- 执行口径同 rev3 §8：先正例绿 → 变异红 → 恢复绿 + 字节复原核对；M15/M17 用
  `agentRunnerObserveManagedWriteProof` 包内 var 注入 stub（现成 seam）；M14 夹具需非提升
  令牌，不可建 junction 则 Skip 并写明。

---

## 7 实施次序（供 ⑧ 派单，本轮仅设计）

1. **evidence Windows 根链**（§3 原语：`openBoundDirChainVerified` + 根链分类暴露）——
   底座，供 adapters 观察分类复用。
2. **console 取证句柄链**（§1：`agentRunnerPolicyDomainChain`、`readFrom`、
   `ReassertIdentities`）——load/publish/confirm 迁移到链。
3. **adapters 四槽固定域 + 句柄形式观察**（§2 底座）+ 锚 `ManagedWriteProofShape` 字段与
   序列化；摘要域常量升级 `:2`。
4. **reverify 重写**（§2+§5：同域对证 + 全集强制验权 + 删 :688 捷径）+ Windows 发布前
   根链复核（§3 调用序）。
5. **M8 三腿重写（§4）** + M13–M17 落地。
6. 收口：`go build ./...`、`go vet ./...`、`go test ./...` + 双平台原生反例执行 +
   docs 同步复核 → 交 ⑦ 复审。

次序理由：1 是 3 的观察分类依赖；2 是 3/4 的句柄来源；4 依赖 2/3 的 Shape 与四槽；测试
随实现步骤就地落地，避免事后补。

---

## 8 登记事项（非设计变更）

- `CONTRACTS.md`：锚结构增 `ManagedWriteProofShape`（协议先行，先文档后代码）；写守卫
  摘要域 `:2` 与四槽 `∅` 标记语义；Windows 取证绑定根名称链拒链（F04 边界表述更新：
  绑定名≠spawn 名，取证绑定无豁免）。
- `SECURITY.md`：Windows/Linux 观察差异（§5-3）与「目录句柄 fsync 不可表达」的既有
  F04/F10 措辞保持；明确「旧锚无 Shape ⇒ 拒绝」为 fail-closed 迁移规则。
- 不改变：BYOK 禁用结论、S1 写域证明必然失败的语义、F04 卷根排除裁定的适用范围。

