# ⑦ 复审 r2 原文（P85 · GPT-6.2 Sol Max · 复审 2/3）

RE-REVIEW: FINDINGS

P85 ⑦ 复审 r2（2/3）。只读核对暂存补丁、相关生产调用链、测试和日志；未改文件、未重跑测试、未重开 F14 等已 CLOSED 项。

**八项复核**
1. **F01 · NOT-CLOSED。** managedWriteProof 已非导出，未发现生产赋值旁路；重入也确实在读前观察目录、父目录及制品/锚并纳入文件身份，但取证与随后打开的对象没有身份对证，正向重入还存在下述摘要范围不一致。
2. **F02/F06/F20 · NOT-CLOSED。** VerifyBoundDir(*BoundDir, FileID) 已迁移全部调用点，Linux openat2 单次解析和单次发布的句柄内操作成立；但 T2 发布入口重新绑定目录，未继承保护域取证身份，Windows 根绑定也不拒链。
3. **F07 · CLOSED。** shebang 与 PT_INTERP 的原始解释器路径都在拆分前经过规范拼写校验，.. 等不会在这条生产链被 Clean 消掉；shebang 按含 #! 的前 256 字节保守拒绝窗口外换行。通用 evidence 入口仍会 Clean，不应宣称它也实施了同样的拒绝。
4. **F10 实现 · CLOSED。** T3 重入的读取、身份确认和同步确实使用同一 BoundDir，Linux 首发的链接与同步也同柄；Windows 只确认原子可见，不计断电持久。
5. **F13 · CLOSED。** 预检用 RFC 8785 规范编码派生规格摘要，CheckSandbox 拒绝 BYOK 缺摘要或 Allowed 回显不符，执行器以冻结规格重取并比较平台证明；这是规格绑定，不是 OS 写限制已获证明。
6. **新 Medium · NOT-CLOSED。** Publish 入口和完成前都进入真实复验，但短域匹配直接返回，可跳过制品与锚的当前写权检查；M10 只证明异目录合形摘要被拒。
7. **F18 · CLOSED。** M11 的替身进入真实 Execute，仅改变平台 verdict 即改变计数；观察点在最后一道 BYOK 门后，测的是到达启动调用，不是 OS 子进程已创建。
8. **F10 反例 · NOT-CLOSED。** M8 关闭守卫的同时把 failSync 设为 false，且没有「守卫开启、重入同步仍失败」的断言；删除重入同步门也可通过当前测试。

**三项终裁**
- **T2 · NOT-CLOSED。** 冻结和 Publish 双读已接通，但保护域与实际读写未绑定为同一对象，短域复验及重入收敛仍有缺口。
- **T3 · NOT-CLOSED。** Linux 解释器链和同柄同步实现已修；Windows carrier 根绑定及 M8 判别性仍阻断，不是单纯缺⑨日志。
- **T8a · CLOSED（代码复审）。** 生产装配未注入 Sandbox，平台仍返回 Unknown，secret 分支在 Resolve 前拒绝；保持 BYOK 禁用。

**新增发现**
- **High · agent_runner_policy.go:558**：制品、锚分别重新绑定目录，随后再按名观察保护域，没有比较三者的目录身份。M4 只证明取得句柄后的单次发布不会改投替身，不能证明取证对象就是读写对象。
- **Medium · agent_runner_policy.go:451**：首发锚封存「目录＋父目录」摘要；无注入重入即使完整域取证成功，也用「另加制品＋锚」的摘要重建期望锚并逐字节比较，未变化的首发锚也会被拒。
- **Medium · bounddir_windows.go:93**：OpenBoundDirVerified 用普通 CreateFile 打开 carrier 根，再读取该句柄的 ID；相对子项的重解析检查不覆盖根的名称链。发布前后两次 FileIDAt 也不是同柄发布。已登记的 F04 CreateProcessW 残余不能替这里背书。

**实跑读数**
Ubuntu 摘录有 15 个包级 ok，含 adapters、console、guard，可见范围无 FAIL；这是 tail -18，不含 go test 独立退出码或逐例 PASS/SKIP。M11 的 Linux 夹具固定用 /bin/sh，来源不可证会 Skip，不能由包级 ok 推定执行过。Windows 同候选的全量 15 包及 13 项变异红、复原后绿、字节一致，本轮仅有主控声明，未见对应原始日志。index/tree 各 83 路径，三作用域均 TOTAL_FAIL=0，SELFTEST 740/740；scope=ci 仍是 v69 的 10 路径快照，不是 P85 双腿 CI 证据。

**⑧ / ⑨ 指令**
- **⑧ 保持未闭环，不进正式⑨。** 保留用户已裁定的 R=4，按 OB-37 记录四要素；下一次应计返工先交用户裁决，不自动追加⑥。获准后只修上述 Medium+ 路径：统一取证与读写身份、对所有已存在对象强制验权并修正锚重入比较、补 Windows 根绑定、让 M8 在同步持续失败时只翻守卫。r3 为最后一轮额度，只核本轮 Medium+。
- **⑨ 补证不替代整改。** F17/F22 不单独阻断代码复审，也不得记闭环。同 inode 特权腿仍只登记未执行；其余原生反例须保留非提权 Windows、非 root Ubuntu 的有效正对照与明确 PASS/SKIP。归档同一候选的未裁剪日志、退出码、变异红/恢复绿及字节复原，以及双腿 Gates 和三包包级 ok。
