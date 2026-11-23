# ⑦ 追加三读·Neon 交叉复核原文（P85 · Gemini 4 Neon Extra High · 试水第二腿）

FINAL REVIEW: FINDINGS

P85 ⑦ 追加复审（R=6 批，只读独立判定，未改文件、未读其它模型结论，仅核三组修复与对应终裁）。

## 一、三组修复逐组复核

1. **Linux ACL 取证对证 · NOT-CLOSED**：文件槽已通过 BoundFile.AccessACL（SYS_FGETXATTR）完成 fd 化，Linux 句柄路径在按名读目录 ACL 后亦执行了 Fstat 与闭名对证及调用端终态 ReassertIdentities；但三明治窗口登记不如实（把内部未给 BoundDir 加 AccessACL() 写成 Linux 不支持目录 fd 读 ACL、本层无法消除，且把窗口写反），Windows 侧存在闭名对证之后的按名读，且 Linux 闭名对证与目录槽无测试/变异覆盖（见新发现 1、2、4）。
2. **摘要槽位对齐契约 · CLOSED**：managedAgentPolicyDomainPlaceFileSlots 将向量固定为四槽 [策略目录 | 父目录 | 制品 | 锚]，域常量与空槽标记 ∅ 与 CONTRACTS 双语逐字一致；不按路径排序，只传锚时制品槽留 ∅，双角色争槽拒绝；双平台用例与 MUT-1 均已佐证。
3. **T3 Windows 取证与发布复核走根链 · CLOSED（仅限 r3 两处根链替换点；T3 整体仍受新发现 3 与 Windows 按路径发布残余阻断）**：可执行文件打开已改用 OpenBoundDirVerified 逐组件拒 reparse 并以同柄打开叶、错误携带 directoryID，载体发布复核已前后各执行一次根链重绑定比对且经 MUT-3 证伪；但 Windows 下 verifyInstallSourceImmutability 丢弃已绑定 leaf 句柄，且 CreateTemp/os.Link 仍按路径执行。

## 二、新发现

- **Medium · agent_runner_policy_write_guard_unix.go#L95-L105 与 #L227-L241**：目录槽 ACL 的「S1 边界登记」不实且窗口本可消除。注释宣称「Linux 尚不能对目录描述符查访问 ACL」「该按名重解析本层无法消除」；事实是 Linux fgetxattr(2) 原生支持 O_RDONLY|O_DIRECTORY 描述符（仅 O_PATH 报 EBADF），而 BoundDir 正是以 O_RDONLY|O_CLOEXEC|O_NOFOLLOW|O_NONBLOCK|O_DIRECTORY 打开且 :107 已拿到 fd。真正原因是 evidence 尚未为 BoundDir 暴露 AccessACL()，并非内核限制。另 :100 将窗口写反为「between the second comparison and getxattr(2)」（实为开/闭两次比对之间），且该边界目前仅在代码注释、未入文档。建议：在 evidence/acl_linux.go 为 BoundDir 增加基于 SYS_FGETXATTR 的 AccessACL() 直接消除目录槽按名窗口；若确作 S1 延期，须改写为如实的内部 API 边界并在报告显式登记。
- **Medium · agent_runner_policy_write_guard_windows.go#L115-L126**：Windows 写保护观测的闭名对证未覆盖全部按名读。after 对证在 installSourceMutationAccess 之后即结束，而按名读取 DACL 并进入摘要的 installSourceDescriptorDigest(slot.path)（GetNamedSecurityInfoW）落在闭名对证之后（且闭端未复查 ClassifyRootChainNoReparse）。建议将闭名身份比对（及根链拒 reparse 复查）移至 installSourceDescriptorDigest 之后。
- **Medium · ai_install_proof_windows.go#L157-L243**：Windows T3 取证中 ACL/DACL 观测与已取摘要/身份的叶句柄未对证。VerifyAIInstallSource 先以 openInstallSourceLeaf 拿到 file 句柄并取 FileID/ExecutableHash，随后传给 verifyInstallSourceImmutability(path, file)；但 Windows 实现于 :158 直接 _ = leaf 丢弃句柄，后续全部组件 reparse、逐权限 CreateFile 与 GetNamedSecurityInfoW 均按路径执行，且无任何 FileIDAt(path) == leaf.FileID() 校验。建议在 verifyInstallSourceImmutability 中取 leaf.FileID() 并在探测前后对叶路径做身份对证。
- **Medium · p85_r6_sandwich_linux_test.go#L25-L60 与 mutation-r6.txt**：测试与变异覆盖缺口。① Linux 三明治用例复用 p85R6SandwichFixture 构造的是文件槽（Linux 实际走 fgetxattr 而非目录按名读），且生产路径调用腿提前替换导致仅命中开名对证；删除 Linux 闭名对证不会令任何用例转红；② mutation-r6.txt 仅跑 3 项，缺 Linux 三明治与 OpenBoundDirVerified 的变异验证。
- **延续边界核对（T2/T3 Windows 发布按路径写入）· agent_runner_policy_publish_other.go#L50-L72、ai_install_proof_publish_other.go#L53-L77**：两处 Windows 发布虽已前后改走根链重绑定，但 os.CreateTemp/os.Link 仍按路径执行，前后对证之间被插链/替换时链接已发生、后置拒绝无法撤回。注释将其统一称作「registered F04 exemption」扩大了 F04 原有边界，不得默认以 F04 豁免，须由用户显式裁定为 S1 Windows 发布路径边界并独立登记，或改句柄相对创建/链接。
- **Low（不阻断）**：① 显式 Role 与路径后缀矛盾未拒绝，且后缀匹配区分大小写；② 非生产名义入口在按名 Getxattr 后无二次 Lstat 对证。

## 三、三项终裁复核（仅基于三组修复后状态）

- **T2 · NOT-CLOSED**：四槽摘要已闭合（组 2 CLOSED）；但 Linux 目录槽 ACL 边界失实与闭名对证缺覆盖、Windows 写保护 installSourceDescriptorDigest 落在闭名对证之外，以及 Windows 策略发布 CreateTemp/os.Link 按路径写入残余尚未闭合或获独立边界裁定。
- **T3 · NOT-CLOSED**：根链两处改动落地；但 Windows verifyInstallSourceImmutability 仍丢弃 leaf 句柄致 ACL/DACL 与已取摘要/身份未对证，缺 OpenBoundDirVerified 变异证据，且 Windows 载体发布按路径写入残余同 T2 待定性；原子可见仍不计作断电耐久。
- **T8a · CLOSED（代码复审）**：生产装配仍未注入 Sandbox，双平台默认返回 VerdictUnknown 在 Resolve 前零启动拒发 BYOK，本次 R=6 未引入回退；维持 BYOK 禁用、不宣称 OS 写限制已获证明。
- **输入读数核对**：numstat 与 gate-index/gate-tree 一致（96 路径，TOTAL_FAIL=0，index=worktree），SELFTEST 740/740；gate-ci 仍为 v69 历史快照；test-r6 仅跑 Windows 3 包；vm-r6-legs 见 15 包 ok 但无退出码与候选提交绑定。

## 四、⑨ 口径声明

维持 F17 / F22 既有口径：两项证据债不单独阻断代码复审，也不得记为闭合；获准推进至 ⑨ 后，须补齐并归档同一候选提交的 Windows + Ubuntu VM 双端全仓完整日志与退出码、非提升 Windows / 非 root Ubuntu 原生正对照与逐腿 PASS/SKIP 明细、全部必核项的变异红与逐字节复原记录，以及双腿 CI Gates 与 adapters/console/guard 包级 ok 行——缺证不得关片。
