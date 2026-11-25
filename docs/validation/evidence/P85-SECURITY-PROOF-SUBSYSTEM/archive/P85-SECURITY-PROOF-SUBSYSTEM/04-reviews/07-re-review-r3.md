# ⑦ 复审 r3 原文（P85 · GPT-6.2 Sol Max · 复审 3/3 · 最后一轮）

BYOK 默认拒发已核实：生产装配未注入 Sandbox，两个平台默认都返回 Unknown。终裁会把代码仍未闭合的项和仅待⑨归档的证据分开列出。

RE-REVIEW: FINDINGS

P85 ⑦ 复审 r3（3/3）。只读核对冻结 diff 的相关段落、工作区源码及实跑日志；未改文件、未重跑测试，不重开 F07、F10 实现、F13 或 F18。

**新发现**
- **Medium · agent_runner_policy_write_guard_unix.go:118**：Linux 的 ACL 仍按路径读取，前后两次 Fstat 查的是同一个已打开的 fd。它们能证明 fd 没变，却不能证明 ACL 读取的是该 fd 所指的对象；最后一次名称身份复核也只能证明终态。
- **Medium · agent_runner_policy_write_guard.go:217**：`:2` 的文件槽仍按路径排序后追加。实际 .anchor.json 排在 .policy.json 前，全集变成 [目录 | 父目录 | 锚 | 制品]；只传锚时，锚还会占制品槽。这不是契约规定的固定槽位，M15/M17 的摘要替身不检验实际入哈希向量。
- **Medium · ai_install_proof.go:416**：T3 可执行文件取证仍调用 OpenBoundDirNoFollow，其 Windows 实现转入普通 OpenBoundDir，未走新的逐组件拒链。后续独立按名分类可以拒绝持续存在的链，但没有把分类结果与已经取摘要的对象绑定。这不是 F04 的 CreateProcessW 残余。

**r2 四组复核**
- **F01 · NOT-CLOSED。** 两份文件已相对同一 BoundDir 打开，终态也复核句柄与名称身份，但 Linux ACL 没有与取证句柄对证。
- **F02/F06/F20 · NOT-CLOSED。** Linux 发布继承取证句柄；Windows 策略发布虽在前后走拒链复核，CreateTemp/Link 仍按路径写入，后置拒绝不能撤回已经发生的链接。
- **短域捷径 · CLOSED。** 重入复验先验当前全集，再比较封存域；Publish 入口和完成前都走真实复验，M17 恢复捷径后按预期转红。
- **F10 反例 · CLOSED。** M16 的 R/G 腿保持同步失败，只改变守卫；删除同步只打中 R，删除守卫条件只打中 G，G 腿同时确认载体字节不变。

**r2 三条发现**
- **High :558 · NOT-CLOSED。** 原先制品、锚各自重绑目录的问题已消除，但真实权限观察仍有上述对象绑定缺口；M13 双平台替身腿证明终态能发现制品替换，Ubuntu 原生腿证明目录身份能区分，均未证明 ACL 取证同对象。
- **Medium :451 · CLOSED。** Shape 分支按 creation 重建空文件槽、按 reuse 重建全集；M15 的未改 creation 锚获采纳，缺 Shape 与错域摘要被拒。该实跑证明控制流，不等于真实权限正例。
- **Medium Windows 根链 · NOT-CLOSED。** 初次绑定及策略发布复核已拒链，M14 原生用例有区分力；但 T3 载体发布复核仍用会跟随 junction 的 FileIDAt，M14 测的是策略发布助手，不是这条路径。

**契约与终裁**
- 根链 CN/EN：双语一致，T3 实际取证绑定未满足。Shape CN/EN：与实现一致，缺失或未知即拒，不迁移旧锚。四槽 CN/EN：域版本与 ∅ 一致，槽序及缺槽位置不一致。
- **T2 · NOT-CLOSED。** 冻结和 Publish 复验已接通，但权限取证、Windows 实际写入及摘要槽位仍有缺口。
- **T3 · NOT-CLOSED。** Linux 同柄发布和 M16 判别性成立；Windows 安装来源绑定与载体发布复核仍未闭合。原子可见不计作断电耐久。
- **T8a · CLOSED（代码复审）。** 生产装配未注入 Sandbox，两平台默认 Unknown，secret 分支在 Resolve 前拒绝。保持 BYOK 禁用，不宣称 OS 写限制已获证明。

**双平台读数**
- Windows：M13–M17 对应本机腿有 PASS，M14 包括中间、末端及同对象 junction；日志可见 adapters、console、evidence 三包 ok，不据此推定全仓通过或令牌未提升。
- Ubuntu VM：M13/M16 原生腿、ELF 正对照与权限／解释器拒绝腿有 PASS；全量段 15 个包级 ok、可见范围无 FAIL，F17 特权腿为 SKIP。摘录未给出完整退出码或双端候选绑定。
- 变异：8 项均 1 → 0，四个变异文件恢复 SHA-256 与基线一致。index/tree 各 89 路径，三份门禁日志均 TOTAL_FAIL=0，SELFTEST 740/740；本地 scope=ci 仍是 v69 的 10 路径快照，不是 P85 双腿 CI。

**⑧／⑨ 指令**
- **⑧ 停在 ⑦，不进正式⑨、不翻闭环。** 保留用户 2026-11-23 裁定的 R=5 及授权留痕；r3 已到额，不自动追加返工、⑥或复审，先报用户裁定。
- **⑨ 维持 F17/F22 口径。** 两项证据债不单独阻断代码复审，也不得记为闭合。获准推进后，归档同一候选的两端完整日志与退出码、非提升 Windows／非 root Ubuntu 的原生正对照和逐腿 PASS/SKIP、变异红与逐字节复原，以及双腿 Gates 和 adapters、console、guard 的包级 ok。缺证不得关片。
