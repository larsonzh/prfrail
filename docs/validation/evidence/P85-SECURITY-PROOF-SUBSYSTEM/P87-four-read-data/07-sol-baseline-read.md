# ⑦ 追加三读·Sol 基准读原文（P85 · GPT-6.2 Sol Max · 试水第三腿）

FINAL REVIEW: FINDINGS

**三组复核**
1. **Linux ACL：NOT-CLOSED。** 文件 ACL 已由 `BoundFile.AccessACL` 在同一 fd 上读取；目录成功路径有开名、ACL、`fstat`、闭名，生产路径也有终态 `ReassertIdentities`。但 [ACL 拒绝分支](internal/adapters/agent_runner_policy_write_guard_unix.go#L155)直接返回，不是每次按名读后都执行闭名核验。
2. **摘要四槽：NOT-CLOSED。** 规范输入的 `:2` 摘要按 `[策略目录 | 父目录 | 制品 | 锚]` 排列，仅传规范锚时制品槽确为 `∅`；但[构槽仍有角色覆盖和路径排序兜底](internal/adapters/agent_runner_policy_write_guard.go#L243)，拒绝摘要也可能少槽，不能认定全路径与 [中文契约](docs/CONTRACTS.md#L233)及[英文契约](docs/CONTRACTS_EN.md#L227)一致。
3. **Windows 根链：CLOSED，限 r3 指定的两条根链修复。** [可执行文件取证](internal/adapters/ai_install_proof.go#L426)以 `OpenBoundDirVerified` 取得目录句柄和身份，相对该句柄打开叶，叶分类、身份与摘要取自同一次打开，叶打开失败带目录 ID；[载体发布](internal/console/ai_install_proof_publish_other.go#L50)在创建临时文件前、成功 `Link` 后各按根链重绑定并与原句柄身份对证。

**新发现**
- **Medium · [Windows 安装来源权限观察](internal/adapters/ai_install_proof_windows.go#L157)。** 不可变性检查忽略已绑定的 `leaf`，逐权利探测及 DACL 摘要都另按名称读取，未与已哈希对象对证。根链和叶分类已绑定，但不能由此推出权限事实也属于同一对象。
- **Low · [Linux ACL 闭名与窗口登记](internal/adapters/agent_runner_policy_write_guard_unix.go#L100)。** S1 残余确有源码登记，但不是注释所说的“第二次比对与 `getxattr` 之间”：未对证的名称解析发生在开名与闭名之间。ACL 存在、读取失败或后置 `fstat` 失败均会先拒绝，不会误放行，但也不会执行闭名核验。[Linux 用例](internal/adapters/p85_r6_sandwich_linux_test.go#L27)单测了闭名助手，生产调用只验证了开名拒绝。
- **Low · [角色优先级](internal/adapters/agent_runner_policy_write_guard.go#L342)。** 显式 `Role` 可覆盖文件后缀，无标记路径仍可排序占用空槽。当前生产调用传入规范路径和匹配角色，未见实际错槽，但构槽入口没有强制这一约束。
- **Low · [Linux](internal/adapters/agent_runner_policy_write_guard_unix.go#L60)与 [Windows](internal/adapters/agent_runner_policy_write_guard_windows.go#L51)拒绝摘要。** 已占槽观察失败且返回空行时，该槽被直接省略，仍产出 `:2` 摘要。调用方会拒绝，不会据此放行，但该摘要不再是固定四槽。

**终裁复核**
- **T2：NOT-CLOSED。** 文件 ACL 同柄、正常四槽及终态核验已落实；逐次闭名和固定槽契约仍有上述残差。
- **T3：NOT-CLOSED。** 两条 Windows 根链修复成立，但权限证据尚不能归于同一已取摘要对象。载体 `CreateTemp` / `Link` 仍按路径执行，后置拒绝不撤销已发生的链接；不把原子可见计作断电耐久。
- **T8a：CLOSED（代码复审）。** 生产未注入 Sandbox，两平台默认 `Unknown`，[预检在解析凭据前拒发](internal/adapters/ai_preflight.go#L112)。保持 BYOK 禁用，不宣称平台禁写已获证明。

**⑨ 口径**
只读核对冻结 diff、工作区和已有日志，未改文件、未重跑、不重开其他已闭项。Windows 摘录三包 `ok`、退出码 `0`，三项变异均 `1 → 0` 且逐字节复原；Ubuntu 摘录有 15 个包级 `ok`，未给出完整退出码或同一候选绑定。三份门禁 `TOTAL_FAIL=0`，自测 `740/740`；本地 `scope=ci` 是 v69 的 10 路径快照，不是 P85 双腿 CI。

**F17/F22 维持：不单独阻断代码复审，缺证不得关片。** F17 豁免只限同 inode 特权腿。⑨ 仍须归档同一候选的双端完整日志与退出码、非提升 Windows／非 root Ubuntu 原生正对照和逐腿 `PASS/SKIP`、变异红与逐字节复原，以及双腿 Gates 和 `adapters`、`console`、`guard` 的包级 `ok`。本次基准读临时授权不作先例。
