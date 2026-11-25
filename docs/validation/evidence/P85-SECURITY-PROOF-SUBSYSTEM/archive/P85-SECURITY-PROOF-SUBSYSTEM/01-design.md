# P85 `SECURITY-PROOF-SUBSYSTEM` · ① 架构修复方案书

- 切片：P85 · L 档实现类 `[SLICE]` · 协议先行；① 架构师 = deep-reasoner / DeepSeek V4 Pro Max
- 承接：P83 报告 §8 顺延清单 P84-a（T2/F3）、P84-b（T3/F6）、P84-c（T8a/F8），不扩不缩
- 输入：起跑消息原文（`docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/00-briefing.md`）、
  `security-findings-remediation.md` §4/§5/§8、P83 终审与三轮复审及追加复审原文、
  现状代码（review.go / runtime_agent.go / runtime_supervised.go / ai.go /
  ai_probe_copilot_executor.go / ai_probe_copilot.go / ai_preflight.go / guard / evidence）、
  `CONTRACTS.md`、`SECURITY.md`、`DELIVERY_DIRECTIVE.md` §1.4
- 行号口径：以本片步骤 0 仓库现状（HEAD `e35e50f`，由主控锚点核对）为准；② 落地后需对最终树重钉。

---

## ① 契约影响判定

判据（起跑消息 §一 + `DELIVERY_DIRECTIVE.md` §1.4）：凡语义变化，顺序固定
「CONTRACTS 双语 → schema / fixtures → 代码 → 测试 → 文档回写」；任何实现若需
新 schema 记录、新 CLI 命令、新依赖或第三方包 ⇒ 判「触或需升级」，不得自行扩域。

### 1.1 T2（F3）：不触（DOES-NOT-TOUCH）

- 被满足句（现行契约已含所要求语义，无需增改）：
  - `CONTRACTS.md:144-146`：「`review-receipt.schema.json` 是唯一产生评审结论的对象……`reviewMode=policy`
    必须绑定**已批准的** `policyHash`」；
  - `CONTRACTS.md:147-148`：「评审者与候选变更产生者是否同一身份、`evidenceRootHash`/`policyHash`
    引用**是否真实存在**及 waiver 到期后不得复用，均由独立 checker 判定」——
    「引用是否真实存在」正是本项要补齐的验证；本项让 `grant.PolicyHash` 被一个真实存在的
    冻结策略体对证，是对 :147-148 的实现性补全，不是新语义。
- 不需要的工件（判据逐项排除）：
  - 新 schema：无。冻结策略制品是 **store-local 协议**，与 `handoff-policy.json` 同构
    （`runtime_supervised.go:73` 与 :63-66 注释明示「store-local protocol, not wire records:
    no schema and no external consumer depends on them」）。本项制品写入 `<runRoot>/agent-runner-policy/`，
    不进入任何 wire 记录，不改 `run-manifest.schema.json`（其 `policyHash` 字段 :142 与
    `bindings` kind `policy` :94 已存在，本片**不写 run manifest**）。
  - 新 CLI 命令：无。制品由既有 `ExecuteAgentRunnerRun` 装配路径在首次派发前写出；
    grant 仍由既有 approvals 流程产生，`authorization-record.schema.json:53-59` 的
    `policyHash`/`runManifestHash` 字段不变。
  - 新依赖 / 第三方包：无（Go 1.22 stdlib）。
- 判定：**不触 CONTRACTS 语义**。协议先行义务 = 无契约改动，② 直接进代码；但 §② 的
  「执行者不可写」与「发布前重读」两个性质必须在实现与测试中字面成立，否则⑦ 可据
  `CONTRACTS.md:147` 判回退。

### 1.2 T3（F6）：触（TOUCHES）

- 被违反/待补句：`CONTRACTS.md:214`（§7 长段）现行只要求「启动前必须探测并记录：规范化
  os/arch、CLI 平台包/构建指纹 platformHash、可执行文件身份/版本、非交互调用方式、工作目录
  绑定……」，**未要求**「安装来源不可变」这一 F6 顺延义务（Windows 名称→启动映像绑定窗口、
  Linux 同 inode ≠ 字节不可变）；`SECURITY.md:23`「限制必须由 runner/OS 能力落实」与
  `SECURITY.md:60`「网络默认拒绝，例外授权须列……接收端」亦未成文绑定「ACL / 不可写安装目录」
  这一 OS 能力到预检判据。
- 触发物：① 新的持久证明载体 `ai-install-source-proof.json`（store-local，但需要契约写明
  位置、读写双方、判定位与平台语义）；② `ai check` 的预检判据收紧（无证明即拒绝 `available`）；
  ③ 平台专属语义（Windows ACL 查询 / Linux 逐组件模式与属主检查 / 其余平台 fail-closed）。
  虽**不需要**新 schema、新 CLI 命令与新依赖（载体 store-local、`ai check` 语义扩展、
  syscall 实现），但上述③属安全语义变化 ⇒ 按 §1.4 协议先行。
- 判定：**触**。先行契约修订草案见 §1.5-A。

### 1.3 T8a（F8）：触（TOUCHES）

- 被违反/待补句：`CONTRACTS.md:214` 的失败原因闭集（「失败原因至少区分 `configuration_invalid`、
  ……`timeout` 和 `unknown`」）不含「BYOK 写限制未证明」这一分类；该段亦未成文
  「无法证明子进程禁写未脱敏落盘 ⇒ 不得解析凭据、不得启动」的零启动门禁。`SECURITY.md:23`
  「限制必须由 runner/OS 能力落实；不能仅靠提示词」是本门禁的设计依据，但产品契约侧需在
  CONTRACTS 成文。
- 触发物：① 新增闭集原因类（若复用 `configuration_invalid` 会误分类，违反 T8b 的闭集分类纪律）；
  ② 新增「BYOK 证明门禁」语义（预检阶段、解析凭据之前判定）。不需要新 schema、新 CLI 命令、
  新依赖。
- 判定：**触**。先行契约修订草案见 §1.5-B。

### 1.4 汇总

| 项 | 判定 | 依据 |
|---|---|---|
| T2 (F3) | 不触 | `CONTRACTS.md:144-148` 已含「已批准的 policyHash + 引用真实存在」语义；制品 store-local 同构 `handoff-policy.json`；无新 schema/命令/依赖 |
| T3 (F6) | 触 | `CONTRACTS.md:214` 缺「安装来源不可变」判据与平台语义；新增持久载体 + 预检收紧 |
| T8a (F8) | 触 | `CONTRACTS.md:214` 闭集缺新原因类；新增 BYOK 零启动门禁语义 |

### 1.5 先行契约修订草案（②a 契约，先于代码；双语严格镜像）

#### A. T3：CONTRACTS.md §7 在 :214 段（「……`timeout` 和 `unknown`。」之后）追加一句

> （草案·CN）对 `agent-runner-cli` 通道的可用性探测，探针可执行文件的安装来源必须可证明不可变：
> Windows 上以稳定绝对路径命名且其所在安装目录的 ACL 对本用户令牌关闭写入/删除；Linux 上其安装
> 目录全路径各组件对本用户不可写，且可执行文件为 ELF 或对其 shebang 解释器同样满足不可变条件；
> 其余平台 fail-closed。该证明以 chain 配置目录派生路径下的 store-local 载体
> `ai-install-source-proof.json` 持久化（写入方 = ProofRail 检查命令、读取方 = 检查命令与探针执行器，
> 每次启动前重读并实时重验路径→文件身份→摘要；载体不得由调用方自选路径）。安装来源不可证明时，
> 不得记录 `available`（以 `cli_unavailable` 分类拒绝）。本载体的 Windows ACL 与 Linux 模式/属主
> 判定分别由平台专属实现承担，判定失败一律 fail-closed。

> （Draft·EN, mirror, to be finalized by ②a with the _EN line-number anchor)

- 改动点：`docs/CONTRACTS.md` §7（:214 段后）；`docs/CONTRACTS_EN.md` 同位镜像。
- schema / 夹具：**不涉**（载体 store-local，无 wire 记录；不需新 schema；如后续评审要求
  goldens，另按协议先行补，不自行扩域）。

#### B. T8a：CONTRACTS.md §7 两处

> （草案·CN·1，追加在 :214 段失败原因闭集）……`response_invalid`、`timeout`、
> `byok_confinement_unproven` 和 `unknown`。

> （草案·CN·2，追加在 :216 段（`windows-credential:<target>` 段）之后）secret 认证模式的
> `agent-runner-cli` 探测在无法证明子进程禁写未脱敏落盘时，不得解析凭据、不得启动子进程
> （零启动拒发 BYOK）：证明依据 = 白名单 argv（无日志目的地）、环境白名单、仅继承标准流句柄、
> 受限工作目录与有界内存输出，外加平台写限制判定（Sandbox 检查层）；任一项不可证即拒发，
> 拒绝证据只用稳定闭集分类。host-auth 路径不适用本门禁。

> （Draft·EN, mirror)

- 改动点：`docs/CONTRACTS.md` :214（闭集）、:216 段后；`_EN` 同位镜像。
- schema / 夹具：**不涉**。

---

## ② T2 设计：AgentRunner 持久冻结策略制品 + `grant.PolicyHash` 独立对照

### 2.1 现状句

- `runtime_agent.go:818-825`（T2 STOP 注释）：「NOT DONE: requiring grant.PolicyHash to equal the
  run's FROZEN policy digest. This revision has no durable frozen-policy artifact……」；
- `runtime_agent.go:826-869` `operatorApprovedPolicyHash`：只从操作员账本取最新 active grant 的
  `grant.PolicyHash`（:864 `policyHash = grant.PolicyHash`），**无任何独立对照物**；
- `runtime_agent.go:762-785` `Review`：绑定 `policyHash` 于 `reviewPolicyHashes`（:774-777）并写入
  决策（:783）；`runtime_agent.go:989-1004` `verifyReviewPolicyBinding`：Publish 时重读账本、
  重取 grant 摘要并与已录值比对（:998-1002）——**两处比较都出自同一账本来源**，故「摘要合形、
  策略非独立获批对象」仍可通过（⑦ T2 STOP 原文）。
- 对照（既有正确先例）：`runtime_supervised.go:2311-2323` `persistFrozenHandoffPolicy` 把
  handoff policy 冻结为 store-local 快照（策略体 + 摘要），Review（:1871）与
  Publish（:1024-1044 `verifyReviewPolicyBinding`）都从**磁盘快照**重导摘要。T2 即把同一
  模式移植到 AgentRunner。

### 2.2 证明载体

- 位置：`<runRoot>/agent-runner-policy/policy-<taskId>-<stepId>-attempt-<n>.json`
  （runRoot 下新目录，与 `agent-runner-replay/`、`workspaces/` 平级；`n` = 该 task 的 attempt）。
- 内容（store-local，无 schema）：
  `{runId, taskId, stepId, attempt, policyBody, policyHash}`；
  `policyBody` = 与本 run 的链配置同源的**策略体**（采用与 `preview.go:267-290`
  `previewPolicySnapshot` 相同的 canonical 结构：chainId/chainProfile/chainDocumentationPolicy/
  tasks[{taskId, review, documentationPolicy}]），`policyHash = evidence.Digest("", EncodeCanonical(policyBody))`
  ——与操作员审批侧 `preview.go:123` 的 `PlanPreview.PolicyHash = evidence.Digest("", policyCanonical)`
  **同域同派生**，因此 `grant.PolicyHash` 与冻结制品摘要是可对证的同一量。
- 读写双方：
  - 写方 = ProofRail 装配路径（`ExecuteAgentRunnerRun`，`runtime_agent.go:152` 起）：在
    `agentRunnerBindTasks`（:1471）完成后、首次派发（dispatcher 启动）**之前**，为每个
    task 绑定写出其制品（原子写：同目录临时文件 → fsync → rename → 父目录 fsync；
    「持久化成功是线性化点」——rename+父目录 sync 成功之前，一律视为未冻结，不得派发；
    失败 ⇒ 整个 run fail-closed、零启动，与 C6 预建失败语义一致）。
  - 读方 = ① Review（`operatorApprovedPolicyHash`，:826）；② Publish（`verifyReviewPolicyBinding`，
    :989）。两处都**从磁盘重读**并重导摘要，绝不用内存副本顶替（防漂移）。
- 执行者不可写保护：
  - 路径不相交：受管 agent 的写入根 = `<runRoot>/workspaces/<task>-attempt-n`
    （`runtime_agent.go:49`、:1717 `agentRunnerWorkspacePath`）；政策目录与其双向不相交，
    并在装配时显式断言（复用 `store.RejectWriterRoots` 同一断言家族，`runtime_agent.go:1513`）；
  - 组件级 symlink/reparse 拒绝：写与读均对 `agent-runner-policy/` 全路径做既有
    bounddir/lease 的组件级检查，发现链接即 fail-closed；
  - no-replace：已存在同键制品时只允许逐字节同摘要幂等（重开复用），不同内容拒绝；
  - 残余如实登记：同用户恶意进程的主动写入属 S1 沙箱外（`SECURITY.md:10`「相同 OS 用户下
    目录分离不是恶意进程沙箱」）；本载体的防线是「写者路径封堵 + no-replace + 发布时磁盘重读
    对账」，被篡改会在 Publish 的重读对账处转红并阻断。

### 2.3 改动点（file:line，现状句 → 改法）

1. `runtime_agent.go:826-869` `operatorApprovedPolicyHash`：
   在返回 `grant.PolicyHash` 之前，追加加载 `(runID,taskID,stepID,attempt)` 的冻结制品并
   要求 `grant.PolicyHash == 制品 policyHash == Digest("",canonical(制品 policyBody))`；
   制品缺失 / 不可读（解码或 canonical 失败）/ 合形但摘要不匹配 ⇒ 返回错误（沿用
   `ErrInvalidConfig` 包装），Review 拒绝 → 任务转 repair，不 approve。
2. `runtime_agent.go:762-785` `Review`：无需改结构（:774-777 记录的 `policyHash` 现在
   已经过 1 的独立对证）。
3. `runtime_agent.go:989-1004` `verifyReviewPolicyBinding`：改为**磁盘重读制品 + 重导摘要**，
   三等式：`Digest(body) == 制品 policyHash == recorded == grant.PolicyHash`；任一不等或
   制品不可读 ⇒ 拒绝（`blockedPromotion`，`runtime_agent.go:920-930`），零发布。
4. 新增 `internal/console/agent_runner_policy.go`（或并入 runtime_agent.go）：
   `agentRunnerFrozenPolicyPath`、`persistAgentRunnerFrozenPolicy`（原子写 + 幂等）、
   `loadAgentRunnerFrozenPolicy`（重导摘要校验）；域常量
   `proofrail:agent-runner-frozen-policy:1\n`。
5. `ExecuteAgentRunnerRun`（:152 起，Step 2 与派发之间）：插入「冻结全部 task 策略制品」步骤，
   失败即整体 fail-closed（在 run-root marker 建立后、任何 spawn 前）。
6. 策略体来源：`run` CLI（`cli.go`）在装载链配置后计算 `previewPolicySnapshot(cfg)` 同源
   canonical 并传入新增的 `AgentRunnerAssemblyOptions` 字段（内部 API，非 wire）；
   若 `ChainDefinition` 字段不足以重建策略体，② 必须从 `ChainConfig` 侧透传，不得另造派生。

### 2.4 三类反例的检测路径与拒绝行为

| 反例 | 检测路径 | 拒绝行为 |
|---|---|---|
| 缺失 | `loadAgentRunnerFrozenPolicy` 返回 not-found（Review/Publish 均重读磁盘） | Review 拒 approve（task→repair）；Publish `blockedPromotion` 零发布 |
| 不可读 | 解码失败 / canonical 失败 / 组件级链接检查失败 | 同上，错误经 `ErrInvalidConfig` 分类 |
| 合形但不匹配 | `Digest(canonical(body)) != grant.PolicyHash`，或制品内 `policyHash` 与重导值不等，或 Review 后制品被改（Publish 重读时不等） | 同上；**零发布、无 accepted snapshot** |

### 2.5 可证伪闭合断言

> **F3-T2**：任一 AgentRunner run 的任务经 `Review` 产出 `policy` 模式 approve 之后、
> `Publish` 完成（outcome=completed）之前，以下恒成立：`Digest(磁盘制品的 policyBody) ==
> grant.PolicyHash == 该 approve 回执的 policyHash`；违反（制品缺失 / 不可读 / 合形但不匹配，
> 含 Review 后篡改）⇒ Review 不 approve 或 Publish 返回 failed 且零发布。

### 2.6 回归测试点

- 反例：制品缺失、制品不可读（畸形 JSON / canonical 失败）、合形但不匹配（grant.PolicyHash 与
  制品摘要互异）、Review 后制品字节被替换（Publish 转红）。
- 变异（逐项撤守卫 ⇒ 对应反例红；恢复 ⇒ 逐字节一致）：
  (1) 删除 `operatorApprovedPolicyHash` 的制品比对 ⇒ 「合形但不匹配」反例由红变绿（被抓）；
  (2) 删除 `verifyReviewPolicyBinding` 的磁盘重读（改读内存） ⇒ 「Review 后篡改」反例由红变绿；
  恢复后 `git diff` 零增量、测试逐字节一致。

---

## ③ T3 设计：不可变安装来源（Windows 稳定路径 + ACL；Linux 不可写安装目录）

### 3.1 现状核对（⑦ r3 High 要求核实现状）

- `internal/guard/process_linux.go:57-101`：**已修**。`cmd.Path` 仅在 `approved != nil` 分支改写为
  `/proc/self/fd/3`（:91-101），`ExtraFiles` 只对该子进程授 fd 3（round 8 修复），`cmd.Start`
  失败不回退按名启动（:104-109）。「无条件改写 cmd.Path」问题**不存在**。
- `internal/guard/process_windows.go:492-506` `openApprovedExecutable`：`FILE_SHARE_READ` only ⇒
  持有句柄期间拒写拒删，verify→spawn 窗口已由句柄闭合。
- `internal/adapters/ai_probe_copilot_executor.go:56-88` `requireSameExecutableFile`：一次 bound
  open 同时取文件身份（`fileid_windows.go:20-28` `win32:vol:indexHi:indexLo`；
  `fileid_unix.go:20-27` `unix:dev:ino`）与摘要，身份+摘要双比对。
- **仍缺**（P84-b）：① 句柄窗口之外的「名称→映像」稳定证明（Windows 换文件在两次 check 之间
  仍可替换同名文件——句柄只在最后一次 verify→spawn 内持有）；② Linux 同 inode 字节重写的
  证明（可写文件可原地改写字节保 inode）；③ shebang 解释器重开问题（脚本 fd 被执行时内核按
  shebang 重新按名解析解释器）。

### 3.2 证明载体：`ai-install-source-proof.json`

- 位置：chain 配置目录派生的受信根（与授权账本同派生：`trustedLedgerPath` 家族），
  文件名 `ai-install-source-proof.json`；**路径恒由派生得出，调用方不可自选**（防自批，同
  `ai.go:360-376` 的账本路径纪律）。
- 内容（store-local）：`{schemaVersion, executablePath, executableHash, fileID,
  platform: {os, arch, platformHash}, immutability: {windows: {dirAclDeniesWriteDelete: bool, aclEvidenceHash}, linux: {unwritablePath: bool, ownershipEvidenceHash, interpreter?: {path, hash, fileID}}}, recordedAt}`。
  载体字节入 `proofrail:ai-install-source-proof:1\n` 域摘要。
- 读写双方：
  - 写方 = `ai check`（ProofRail 检查命令，操作员调用）：在批准接收端成功且启动探测前，
    现查平台事实并原子写（临时文件 → fsync → rename → 目录 fsync；持久化成功是线性化点）；
    已存在时重验后仅当事实变化才拒绝（幂等收敛，不同内容拒绝）。
  - 读方 = ① `ai check` 下次运行时；② 探针执行器 pre-spawn（实时重验，见下）。
    载体只作「上次已验证事实的持久备忘」，**判定一律以实时重验为准**（防 TOCTOU）。
- 实时重验（每次带凭据探测的启动前，`ai_probe_copilot_executor.go` 的 pre-spawn 段）：
  Windows：`GetNamedSecurityInfo`/`GetEffectiveRightsFromAclW`（advapi32，syscall）现查可执行
  文件与其所在目录对本进程令牌是否无写入/删除/改权；再走既有 `requireSameExecutableFile`
  （身份+摘要同柄）。
  Linux：逐组件（安装目录全路径与文件自身）现查 mode/属主：本 euid 无写位且属主为 root 或
  非本用户；有 shebang 时再验证解释器路径存在、位于不可写目录并取 (dev,ino,digest)。
  任一步不可执行/不可判定 ⇒ fail-closed 拒绝（分类 `cli_unavailable`）。
- build tag：`ai_install_proof_windows.go`（`//go:build windows`，ACL 查询）、
  `ai_install_proof_unix.go`（`//go:build !windows`，Linux 判定；darwin 等一律返回不可证明 →
  fail-closed）。

### 3.3 改动点（file:line）

1. `internal/console/ai.go:464-487` `defaultAICheck` / `:392-411` `executeAICheck`：
   批准通过后、构造 checker 前，追加「读取或首次落盘安装来源证明」与实时重验；失败即
   `writeCommandError`（`cli_unavailable` 分类），**凭据不解析、进程不启动**。
2. `internal/adapters/ai_probe_copilot_executor.go:189-300` `Execute`：在 `requireSameExecutableFile`
   （:284 附近）前追加实时安装来源重验（Windows ACL / Linux 不可写 + 解释器）；未证明 ⇒
   `start-failed`，证据用 `copilotCLIExecutorEvidence("install-source-unproven")` 闭集分类。
3. 新增 `internal/adapters/ai_install_proof_windows.go` / `ai_install_proof_unix.go`：
   `verifyInstallSourceImmutability(executable string) (facts, error)` + 载体编解码。
4. 新增 console 侧载体路径派生与读写（`internal/console/ai_install_proof.go`，复用
   `trustedLedgerPath` 家族；`ai.go:336-341` 与 `:360-376` 同纪律）。
5. `internal/adapters/ai_preflight.go:57-63`（F6 先验证后解析的既有位点）：无需改顺序，
   T3 门禁并入 2 的 pre-spawn 段即可，但需测试钉住「拒发发生在 `Secrets.Resolve` 之前」。
   （判定顺序保留既有 F6 约定：先 VerifyReceiverBinding、再 T3 安装来源、最后 Resolve。）

### 3.4 三类反例与拒绝行为

| 反例 | 检测路径 | 拒绝行为 |
|---|---|---|
| 载体缺失 | `ai check` 批准后落盘失败 / 上次未落盘且现查无法成立 | 拒发 `available`；不解析凭据、不启动 |
| 载体不可读 | 解码失败 / 组件级链接检查失败 | 同上（fail-closed，不降级为「重算即可」） |
| 合形但不匹配 | 现查事实（ACL 被放开 / 目录变可写 / fileID 或摘要变）与载体、与 grant 批准摘要不符 | 同上；pre-spawn 重验失败即 `start-failed` 零启动 |

### 3.5 可证伪闭合断言

> **F6-T3-Windows**：`ai check` 的 agent-runner-cli 探测在 Windows 上启动带凭据的探针之后、
> 该探测产生 `available` 记录之前，被启动映像的文件身份与摘要恒等于操作员批准摘要，且其安装
> 目录 ACL 对本用户令牌恒无写入/删除；违反 ⇒ 探测零启动（或记录 `unavailable`），凭据不解析。
>
> **F6-T3-Linux**：Linux 上带凭据探测启动的映像对象恒由已验证 fd 创建（`/proc/self/fd/3`），
> 且安装目录与 shebang 解释器全路径对执行用户恒不可写；违反（目录可写 / 同 inode 字节被改写 /
> 解释器不可证明）⇒ 零启动。

### 3.6 回归测试点

- Windows（`//go:build windows` 测试）：目录 ACL 允许本用户写入 ⇒ 拒绝；同路径换文件
  （fileID 变化）⇒ 拒绝；摘掉 ACL 查询 ⇒ 变异红。
- Linux（`//go:build !windows` 测试，Ubuntu VM 原生跑）：安装目录 0755 属本用户 ⇒ 拒绝；
  目录 0555 root 后改写字节（需 root 夹具或临时 chown 还原）⇒ 正向通过；同 inode 改写
  （把目录改回可写再改字节）⇒ 拒绝；shebang 解释器位于可写目录 ⇒ 拒绝；删除解释器验证 ⇒ 变异红。
- 变异恢复：撤守卫 ⇒ 对应反例红；恢复 ⇒ 逐字节一致（含载体字节不变）。

### 3.7 残余风险登记（须写入报告 §8）

- shebang 解释器按名重解析由**内核**完成，ProofRail 无法对内核的这一次打开加锁；以
  「解释器路径不可写 + 载体记录解释器 (dev,ino,digest) 并 pre-spawn 重验」封堵，窗口残余 =
  内核解析瞬间的 TOCTOU（目录不可写 ⇒ 窗口内容为空）。如实登记。
- `/proc` 未挂载的 Linux 主机：fd-exec 不可用 ⇒ `cmd.Start` 失败 ⇒ fail-closed 拒绝（不回退按名）。
- 管理员/SYSTEM 可无视 ACL：威胁模型内（`SECURITY.md:10`）同用户恶意进程不属 S1 承诺。

---

## ④ T8a 设计：BYOK 证明载体（不能证明 ⇒ 零启动拒发）

### 4.1 现状句

- `internal/adapters/ai_probe_copilot_executor.go:189-232`：已拒绝含日志目的地的 argv
  （:103-108、:202-211），环境白名单（:277-282），stdout/stderr 有界内存（:292-296、:517-565），
  但**未证明**「子进程不会自行落盘」，且仍无条件携带 BYOK 环境（`ai_probe_copilot.go:168-185`
  `copilotCLIProbeEnvironment`）启动；
- `internal/guard/sandbox_windows.go:17-25`：S1 沙箱 `alphaSandbox` 恒 `Unknown`（fail-closed）——
  当前任何平台都无法出具「写限制已验证」的判定；
- 复审口径（P83 r3）：「F8（T8a）未证明 CLI 默认落盘日志已关闭或脱敏，默认也未拒发 BYOK」。

### 4.2 门禁设计

- 位点：`internal/adapters/ai_preflight.go:57-63` 的既有 F6 位点之后、`Secrets.Resolve`（:66）之前，
  新增可选接口（镜像 `AIReceiverBindingVerifier` 模式）：
  `VerifyBYOKConfinement() error`。`CopilotCLIProbe` 实现之；执行器无该能力或判定不可证 ⇒
  返回错误，checker 以新闭集类 `byok_confinement_unproven` 记 `unavailable`，
  **凭据不解析、子进程不启动**（零启动拒发 BYOK）。
- 判定构成（五元 + 沙箱层，全部成立才可发 BYOK）：
  1. argv 白名单：无日志目的地参数（既有 `copilotCLIProbeRefusedLoggingArg` 拒发逻辑，
     门禁内重新断言）；
  2. 环境白名单：`copilotCLIEnvironmentAllowlist` + BYOK 覆盖键，无其他键（既有实现，门禁断言）；
  3. 继承句柄白名单：仅 stdin/stdout/stderr（guard `process_windows.go:282-284`
     `buildCreationAttributes` 句柄列表）+ Linux fd 3（执行器与 guard 契约）；门禁断言
     spawn 规格不含其他 ExtraFiles / 继承句柄；
  4. 工作目录限制：子进程 cwd = 执行器专属临时目录，唯一子可写产物 = `--usage-output-file`；
  5. 有界内存输出：stdout/stderr 进有界内存 writer，用完清零（既有 `:517-565`、`:293-296`）；
  6. 平台写限制判定：`CheckSandbox`（S1 = `Unknown` ⇒ 不成立）。
- S1 现实结论：第 6 项恒不成立 ⇒ **S1 上 BYOK 探测一律零启动拒发**；host-auth（无凭据）路径
  不适用本门禁，照常探测（F6/F8 边界「保留非 BYOK 路径」）。
- 拒绝证据：沿用 T8b 纪律——`unavailableCopilotCLIRefusal` 只带两个闭集分类（outcome 类 +
  failure 类，`ai_probe_copilot.go:241-263`），新类映射入 `validCopilotCLIProbeFailure`/
  `copilotCLIProbeFailureClass`；证据**不散列任何参数值**、不读 executor 的 Evidence/Measurements。
- 不落盘的证明依据：本门禁的「证明」= 对 spawn 配置的**封闭声明 + 机械断言**（1-5 项由
  executor/guard 既有机制机械执行并逐项断言）+ 沙箱层判定（6）。断言事实本身经既有
  `copilotCLIProbeEvidence` 摘要链进入可用性记录（非新落盘对象）。

### 4.3 改动点（file:line）

1. `internal/adapters/ai_preflight.go:57-63`：新增 `VerifyBYOKConfinement` 调用位点
   （仅 `profile.AuthMode == "secret"` 时调用，失败 ⇒ 拒发且不 Resolve）。
2. `internal/adapters/ai_probe_copilot.go:99-123`：`Probe` 内处理新拒绝类（分类表
   `validCopilotCLIProbeFailure` 增 `byok_confinement_unproven`；闭集映射
   `copilotCLIProbeFailureClass` 同步）。
3. `internal/adapters/ai_availability.go:220-242` `validAIAvailabilityReason`：闭集增
   `byok_confinement_unproven`（与 CONTRACTS:214 修订同步）。
4. `internal/adapters/ai_probe_copilot_executor.go:189-232`：`VerifyExecutableBinding` 旁新增
   `VerifyBYOKConfinement`（断言 1-5 + `CheckSandbox`）；`Execute` 在凭据存在时以同一断言
   作为 pre-spawn 最后防线（双保险，镜像 F6 的双重校验模式）。
5. `internal/guard/sandbox.go` 的 `CheckSpec`：新增 BYOK 探测 reason（不改沙箱判定实现）。

### 4.4 三类反例与拒绝行为

| 反例 | 检测路径 | 拒绝行为 |
|---|---|---|
| 证明缺失 | 执行器不实现 `VerifyBYOKConfinement` / 断言上下文不可得 | checker 在 Resolve 前拒发，`unavailable: byok_confinement_unproven` |
| 证明不可读 | 沙箱层判定不可得（S1 恒 `Unknown`）| 同上（fail-closed，不得把 Unknown 读作放行） |
| 合形但不匹配 | 断言与 spawn 规格不符（argv 含日志参数、环境有白名单外键、ExtraFiles 非空等）| pre-spawn 拒发，`start-failed` 零启动；拒绝证据不变随值 |

### 4.5 可证伪闭合断言

> **F8-T8a**：secret 认证模式下，`agent-runner-cli` 探测在「子进程禁写未脱敏日志」不可证明的
> 任一时刻（S1 恒不可证），`Secrets.Resolve` 不执行、子进程不创建、`available` 不产生；
> 违反 ⇒ 零启动拒绝且拒绝证据为闭集分类、不随被拒值变化。host-auth 探测不受本门禁影响。

### 4.6 回归测试点

- 反例：沙箱 `Unknown` + secret 模式 ⇒ 断言 Resolver 零调用、spawn 零次、记录
  `byok_confinement_unproven`；host-auth 无沙箱 ⇒ 探测放行（非 BYOK 路径钉住）；
  同类不同值拒绝（两个不同 model）证据相等（T8b 兼容）。
- 变异：删除 checker 的门禁调用 ⇒ spawn 发生（测试红）；删除 executor pre-spawn 断言 ⇒
  测试红；恢复 ⇒ 逐字节一致。

### 4.7 「证明」的边界与剩余风险登记（须写入报告）

- 证明的是**spawn 配置**（argv/env/句柄/cwd/输出五元 + 平台写限制判定），不是子进程内部行为；
  子进程绕过配置自行写 `%APPDATA%`、崩溃转储、OS 页面文件、provider 侧日志均不可证 ⇒
  正是因此 S1 上拒发 BYOK，而非「弱化后放行」。
- BYOK 拒发造成 S1 上 `ai check` secret 模式恒 unavailable ⇒ AgentRunner 的 BYOK 可用性
  预检被阻断——这是 F8 要求的 fail-closed 后果，如实写入报告并指明解锁路径 =
  后续沙箱片出具平台级写限制证明后经 `CheckSandbox` 自动放行（无需再改本片契约）。
- 探针外的 AgentRunner 真实运行路径（launcher 携带凭据）不在本片范围（T8a 只列两个探针文件），
  在报告边界章显式声明，不静默扩展。

---

## ⑤ 跨包共享契约（附录 B.1 第 6 项）

1. 新证据域常量（一律 `<域>:1` 后缀，与既有前缀编号兼容，不重排任何既有事实编号）：
   - `proofrail:agent-runner-frozen-policy:1\n`（T2 制品外哈希；制品内 `policyHash` 仍为
     `evidence.Digest("", canonicalBody)`，与 `preview.go:123` 同派生——**不得**再套一层新域，
     否则与操作员审批摘要不可对证）；
   - `proofrail:ai-install-source-proof:1\n`（T3 载体字节摘要）；
   - `proofrail:ai-byok-confinement:1\n`（T8a 门禁断言集摘要，进入可用性记录证据）。
2. 接口契约：
   - chain 包**不改**（`review.go` 的 `PolicyHash` 校验与 `ReviewEvidenceBinding` 不动）；
   - console → 新载体读写在 `internal/console`（T2 制品、T3 载体派生路径），
     adapters → 执行器与证明函数在 `internal/adapters`；
   - `evidence` 包复用：`OpenBoundFile` / `BoundFile.FileID` / `Canonicalize` / `Digest`
     / `ValidHash`，不新增 evidence 原语（Windows ACL 与 Linux 模式判定放 adapters 的平台文件，
     不得下沉进 evidence 语义）。
3. 写序与线性化点：
   - T2：策略制品**全部原子落位 + 父目录 sync 成功**（线性化点）之后才允许首次派发；
     Review/Publish 只消费已落位制品。
   - T3：安装来源载体原子落位（线性化点）之后才允许该次带凭据探测；pre-spawn 再实时重验。
   - T8a：无新增持久对象；门禁在 `Secrets.Resolve` 之前完成（现有 checker 顺序 :57-66 保持）。
4. 与既有事实/证据前缀兼容：所有新证据一律**追加**进 hashSet/证据切片，不改 A5 六枚事实的
   位置化布局（`evidence[2..7]`）、不改 `evidence-manifest.schema.json`、不改任何既有
   recordHash 输入集；既有回执与记录逐字节不变（重放兼容）。
5. 执行者不可写矩阵（对照 `SECURITY.md` §2）：T2 制品与 T3 载体均位于受管 agent 写入根之外，
   写方唯一 = ProofRail 装配/检查命令；agent（managed）对二者零路径可达；违反时由
   no-replace + 发布/启动前重读对账转红。

---

## ⑥ 失败路径枚举（附录 B.1 第 7 项）

| 故障 | T2 | T3 | T8a |
|---|---|---|---|
| 磁盘满（写载体失败） | 临时文件写入/rename/sync 任一失败 ⇒ 冻结未线性化 ⇒ 整体 fail-closed 零派发；不留半写文件（temp 清理） | 载体落盘失败 ⇒ 该次探测拒发（不降级为「内存重算即可」） | 无新增落盘；可用性记录写失败沿用既有 `WriteAIAvailabilityRecord` 错误路径 |
| 句柄被替换（TOCTOU） | 读制品走组件级检查 + no-replace；被替换 ⇒ 重导摘要不等 ⇒ Review/Publish 拒绝 | 载体只作备忘，判定一律实时重验（ACL/文件身份同柄）；换文件 ⇒ fileID/摘要不等 ⇒ 拒发 | spawn 规格在 pre-spawn 由 executor 内存构造，不涉外部句柄；沙箱判定不可得即拒发 |
| 竞态（检查→启动窗口） | 制品写入后只读 + 发布前重读，无窗口收益 | Windows 由 deny-write/delete 句柄闭合；Linux 由不可写目录 + fd-exec 闭合；残余见 §3.7 | 门禁在凭据解析前完成，之后无可再「证明」的窗口；pre-spawn 双保险再断言一次 |
| 回滚/中断 | 部分制品落位（多任务）⇒ 重开时按幂等核验收敛，未冻结任务拒绝派发；不补造 | 载体原子写；中断 ⇒ 无载体或旧载体 ⇒ 下次重验后重写或拒发 | 无回滚对象 |
| 双平台差异 | 与平台无关（纯文件+摘要） | Windows=ACL 查询，Linux=模式/属主+解释器，其余 fail-closed；平台文件 build tag 隔离，语义表随 §1.5-A 契约成文 | S1 两平台沙箱均 Unknown ⇒ 均拒发 BYOK，行为一致 |

---

## ⑦ 测试矩阵（反例 + 变异逐项）

| 项 | 反例（缺失/不可读/合形不匹配） | 变异（撤守卫 ⇒ 红；恢复 ⇒ 逐字节一致） | 平台 |
|---|---|---|---|
| T2 | 制品缺失；制品畸形 JSON；`Digest(body)!=grant.PolicyHash`；Review 后制品被替换（Publish 红） | 撤 `operatorApprovedPolicyHash` 制品比对；撤 Publish 磁盘重读改内存 | 双平台（纯文件） |
| T3 | 载体缺失；载体不可读；ACL 放开（Win）/ 目录改可写（Linux）；同路径换文件；同 inode 改写字节；解释器可写目录 | 撤 ACL 查询；撤 Linux 模式检查；撤解释器验证；撤 pre-spawn 实时重验改读载体值 | Windows 本机 + Ubuntu VM + CI 双腿 |
| T8a | 沙箱 Unknown + secret ⇒ 零解析零启动；host-auth 放行；同类不同值拒绝证据相等 | 撤 checker 门禁调用；撤 executor pre-spawn 断言 | 双平台（沙箱层 fake 注入） |

验收口径（`SECURITY.md:45-46`）：每个反例断言「不发布 / 不推进 / 证据仍完整」，不只检查非零返回；
变异恢复后 `git diff` 对实现文件零增量（除预期改动点外）。

---

## ⑧ 实施次序建议（先契约后代码）

1. ②a 契约：先落地 §1.5 的 CONTRACTS 双语修订（T3 一句 + T8a 两处，含 `_EN` 严格镜像）——
   本片 schema / 夹具不涉，无 schema 改动。
2. ②b 代码（顺序）：T2（console：策略制品 + 两端口比对，不触契约、可先行）→ T3
   （adapters 平台文件 + console 载体路径 + `ai check` 接线）→ T8a（checker 门禁 +
   分类闭集 + executor 双保险）。
3. ③ 测试按 §⑦ 矩阵；Linux 原生用例必须在 Ubuntu VM（10.0.0.199，go1.22.12）跑，
   Windows 本机跑 `//go:build windows` 用例，CI 双平台腿兜底。
4. 文档回写：新增 `docs/validation/security-proof-subsystem{,_EN}.md`；台账
   `REMAINING_SLICES{,_EN}.md` P85 段闭环；`docs/SECURITY{,_EN}.md:5` 注记改为
   「F1–F8 全闭合（F3/F6/F8 由 P85 闭合）」；`DEV_PLAN{,_EN}.md` 回写 P85 一段；
   `security-findings-remediation{,_EN}.md` 的 P83 报告 §8 补时点注记（用户裁定项）。
5. 边界重申（不碰）：F1/F4（已全闭合）、F2/F5/F7（已闭合）；S1 Exit 八门判据；
   `SECURITY.md` 威胁模型基线；准则版本（不升版、不新增 OB/ADR）；CI 不动；无新依赖；
   T8a 不含 AgentRunner 真实运行路径的 BYOK。
