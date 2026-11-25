# ⑦ 追加双读·主读原文（P85 · Gemini 4 Carbon High · 试水第一腿）

FINAL REVIEW: PASS WITH FIXES

P85 ⑦ 追加复审（主读，只读核对工作区源码、双语契约及 R6 实跑/门禁/变异日志；不重审已 CLOSED 项）。本轮无 Medium+ 阻断项，三组 r3 发现均闭合，新增 3 条 Low（记录在案、不阻断代码复审收敛）。

### 1. 三组修复逐组复核

- 组 1（Linux ACL 取证对证）· CLOSED。文件槽经 agent_runner_policy_write_guard_unix.go#L227-L235 调用 BoundFile.AccessACL（acl_linux.go#L36-L60 直走 SYS_FGETXATTR 读已绑定 fd，零路径解析）；目录槽在 #L124-L173 严格执行「Fstat(fd) → 开名对证 before → 按名读 AccessACLPath → Fstat(fd) 五元组复核 → 闭名对证 after」，闭名对证在 Linux 每次按名读后无条件执行，并由控制台层终态 ReassertIdentities（agent_runner_policy.go#L470、#L550、#L933）收口；#L99-L106 如实将目录槽 getxattr(2) 瞬间的换入换回重解析窗口登记为 S1 已知边界，未作越界声称。
- 组 2（摘要槽位对齐契约）· CLOSED。固定域标识 proofrail:agent-runner-policy-write-guard:2 与空标记 ∅，名称入口与句柄入口统一经 managedAgentPolicyDomainPlaceFileSlots 按角色（Role 字段或契约后缀）放入固定四槽 [策略目录 | 父目录 | 制品 | 锚]，只传锚时槽 3 保持 ∅、槽 4 放锚，同槽冲突直接拒绝，与 CONTRACTS 双语逐字一致；白盒与接缝测试直接钉住入哈希四槽向量与句柄身份，MUT-1 退回路径序即转红且恢复后逐字节一致。
- 组 3（T3 Windows 取证与发布复核走根链）· CLOSED。可执行文件取证 openInstallSourceLeaf 已改调 evidence.OpenBoundDirVerified（Windows 侧从卷根逐组件以 FILE_OPEN_REPARSE_POINT 打开并拒斥任一 reparse/UNC/扩展/相对路径），再相对同一目录句柄打开叶节点并在同一 BoundFile 上读取 FileID 与摘要，错误携带目录对象身份 directoryID；载体发布继承入口已绑定的 (bound, fileID)，在写前与链接后各调一次 aiInstallSourceAssertDirectoryIdentity 并最终在 bound 句柄内回读校验；原生反例与 MUT-3 均具区分力。

### 2. 新发现（0 Medium+，3 Low）

- Low · agent_runner_policy_write_guard_windows.go#L114-L125：Windows 写域观察中第二处按名读取 installSourceDescriptorDigest(slot.path) 位于闭名对证 after 之后，未包入该函数的闭名检查内。因写权判定 installSourceMutationAccess 已处于 [before, after] 内、跨对象污染 dacl: 只会导致摘要失配 fail-closed，且生产调用方紧接执行 chain.ReassertIdentities()，无放行绕过风险。建议修法（可延后）：将 after 身份重取与比较下移至 installSourceDescriptorDigest 之后。
- Low · agent_runner_policy_publish_other.go#L16-L18 及 ai_install_proof_publish_other.go#L11-L18：对应 r3 原文所指 Windows 发布阶段 os.CreateTemp / os.Link 仍按路径执行，两文件注释将其归入「已登记 F04 豁免」，与 CONTRACTS「与 spawn 名的既有表述互不替代」口径不符。在 S1 同用户非特权边界下，O_EXCL + no-replace 硬链接无覆盖、无秘密载荷，且发布后均经原 bound 句柄内回读校验，与全仓 Windows store-local 写入边界一致，不构成安全阻断。建议修法：P85 报告 §8 须将「Windows store-local 发布 CreateTemp/Link 用户态路径解析残余（采纳侧由绑定句柄回读 fail-closed 封闭，不覆盖既有文件）」独立于 F04 如实登记，后续顺延修正两处代码注释措辞。
- Low · p85_r6_sandwich_supported_test.go#L75-L82 及 p85_r6_sandwich_linux_test.go#L27-L55：三明治夹具构造的是文件槽（slot.file），而 R6 后 Linux 文件槽已走 fgetxattr、仅目录槽走按名 ACL；且生产入口子测在调用前替换名称，只命中 before 分支（无生产注入缝的边界已如实声明）。建议修法（可延后）：后续可补一条目录槽夹具覆盖 Linux slot.dir 分支。

### 3. 三项终裁复核（基于三组修复后状态）

- T2 · CLOSED（代码复审）。装配前策略冻结、Review/Publish 入口及终局重读、全集无条件验权 + Shape 同域对证、句柄链贯穿与终态 ReassertIdentities、Linux fgetxattr + 目录三明治及 :2 固定四槽均已对齐契约；Windows 发布路径残余按上述独立用户态 TOCTOU 边界入报告 §8 登记，不与 F04 混同。
- T3 · CLOSED（代码复审）。Linux 同柄发布、解释器原始链与 M16 判别性成立；Windows 可执行文件取证绑定与载体发布前后复核均已走 OpenBoundDirVerified 逐组件拒链且与后续句柄回读闭环；Windows 仅提供 no-replace 原子可见性，不宣称断电持久。
- T8a · CLOSED（代码复审）。维持已闭合结论：生产装配未注入 Sandbox，双平台默认返回 Unknown，secret 分支在凭据解析前零启动拒绝并输出闭集分类 byok_confinement_unproven；S1 保持 BYOK 禁用，不宣称 OS 写限制已获证明。

### 4. ⑨ 口径声明

维持 F17 / F22 既有口径：两项证据债不单独阻断本轮代码复审收敛，但缺证不得关片。进入 ⑨ 归档时须在同一最终候选上补齐并落盘：（1）Windows 本机与 Ubuntu VM 全量 go build/vet/test ./... 完整日志及明确退出码（含候选工作树一致性绑定）；（2）非提升令牌 Windows 与非 root Ubuntu 原生正对照及逐腿明确 PASS/SKIP 读数，F17 同 inode 特权夹具维持未执行四要素如实登记；（3）变异红、恢复绿及 SHA-256 逐字节复原证据；（4）提交后 GitHub Actions 双平台 CI 腿 Gates 步通过与 adapters、console、evidence、guard 包级 ok 硬证据（本地 scope=ci 仍为 v69 基线快照，不替代 P85 双腿 CI）。

FINAL REVIEW: PASS WITH FIXES
