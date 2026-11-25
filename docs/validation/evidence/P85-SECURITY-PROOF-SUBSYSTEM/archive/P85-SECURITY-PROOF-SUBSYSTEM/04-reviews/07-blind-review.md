# ⑦ 盲审原文（P85 · GPT-6.2 Sol Max · 盲隔离）

已完成静态核对，三项各有未闭合点。我会在最终回执中把确定缺陷、条件性风险和未获实证的变异分开列出；本轮未改文件，也未运行测试或探针。

BLIND REVIEW: FINDINGS

编号 | 级别 | 位置 | 事实 | 建议
---|---|---|---|---
F01 | High | [agent_runner_policy.go:133]、[agent_runner_pinned_cli.go:218] | T2 只检查目录不相交；受管进程以同一 OS 身份启动，`cwd` 和 `0600` 不限制它写制品。重读只证明**当前**制品与 grant 相合，不证明它仍是装配时冻结的对象。 | 建立执行令牌不可写的存储边界，并在发布时核验独立的冻结锚点。
F02 | High | [agent_runner_policy.go:330]、[agent_runner_policy.go:362] | T2 目录路径检查不含制品末组件，实际读句柄会跟随链接；目录检查与打开也不绑定。安全目录下的名字仍可能指向可写对象。 | 从绑定目录句柄打开，拒绝末端链接/reparse，并验证文件所有权及可写别名。
F03 | High | [ai_install_proof_windows.go:167] | 对合并权限掩码的一次 `ACCESS_DENIED` 只说明**未同时获得全部权限**，不说明逐项均被拒绝。仅缺改属主权、仍可写或删除的 ACL 也可获证。 | 对有效令牌逐项检查变更权，区分文件权与目录权，并增加只授予单项变更权的原生反例。
F04 | High | [ai_install_proof_windows.go:97]、[process_windows.go:492] | Windows 只检查文件和直接父目录，未证明祖先命名空间、reparse 或父级删除权；拒写删句柄锁住的是叶对象，不锁住 `CreateProcessW` 重新解析的整条名称链。 | 验证并绑定最终启动名及必要的祖先目录，直到创建映像。
F05 | High | [ai.go:411]、[ai_probe_copilot_executor.go:367]、[executable_binding.go:62] | 载体对账完成后只传摘要；执行器丢弃实时证明事实，最终守卫另开对象也只验摘要。读取载体后丢失载体、同摘要异 `fileID`、解释器事实漂移，均不保证拒绝。 | 把持久载体中的身份与平台事实带到启动边界，和最终打开的对象及依赖逐项对证。
F06 | High | [ai_install_proof_unix.go:94]、[ai_install_proof.go:277] | Linux 对词法组件用 `Stat`，会跟随链接而漏查目标的真实祖先；权限、摘要和 shebang 又分次按名取证。主映像使用 fd 启动不能补足这些证明。 | 逐组件绑定真实路径或拒绝链接，从最终执行的同一对象取身份、权限、摘要和头部。
F07 | High | [ai_install_proof_unix.go:152]、[ai_install_proof_unix.go:169] | shebang 解析不等价于内核；`env` 分支跳过首个实际解释器，且不验证递归解释链。ELF 分支只验魔数，动态 ELF 的 `PT_INTERP` 也未纳入运行映像证明。 | 拒绝不能精确证明的解释形式；验证实际加载链，或限制为能证明的格式。
F08 | Medium | [ai_install_proof_unix.go:102] | `0222` 检查的是任何主体的写位，不是当前用户的有效权限。常规 Ubuntu 的 root-owned `/` 为 `0755`，非 root 也会在首组件失败，host-auth 同样被挡。模式位也不证明既有写句柄已失效。 | 按有效权限和文件系统能力判定，并用标准非 root 安装布局跑真实正向。
F09 | Medium | [agent_runner_policy.go:412]、[ai_install_proof.go:170] | 两个写入器均在硬链接失败后用 `Lstat` 加可覆盖的 `Rename` 降级，并发下不具备 no-replace 保证；事后读回不能撤销覆盖。 | 仅用原子 no-replace；无法证明该原语可用时拒绝。
F10 | Medium | [agent_runner_policy.go:284]、[agent_runner_policy.go:429] | T2 新建目录后只同步其内部，不同步创建该目录的父项；Windows 跳过目录同步，且已存在分支可绕过前次失败的同步。不能证明承诺的持久化线性化点。 | 证明完整目录链的持久性，重入时不把未确认持久的文件视为已冻结；T3 同类写入器一并核验。
F11 | Medium | [ai-availability.schema.json:18]、[CONTRACTS.md:214] | 双语契约与 Go 已添加 `byok_confinement_unproven`，wire Schema 仍拒绝。方案「不涉 Schema/夹具」的判断不成立。 | 同步闭集并补 Go 生成、Schema 校验的契约夹具。
F12 | Medium | [ai_preflight.go:110]、[ai_preflight.go:170]、[p85_t8a_tester_test.go:574] | S1 实际前置拒绝的 `Evidence` 混入 profile 摘要及时间，不是纯闭集；同类同证测试只覆盖证明替身通过后的执行器拒绝。 | 专用纯分类证据；profile/time 留在记录元数据。覆盖真实前置拒绝的不同值对照。
F13 | Medium | [ai_probe_copilot_executor.go:242]、[ai_probe_copilot_executor.go:466]、[sandbox.go:61] | 解析前只验静态条件和带标签的沙箱 verdict，不绑定实际 spawn。argv 是黑名单加自我相等，不是白名单；默认 `-C` 仍指调用目录。未来仅把 verdict 改成 `Allowed`，会先读凭据，再因具体配置拒绝。 | 在解析前构造并证明同一不可变启动规格，要求平台证明实际约束该规格；否则维持拒发。
F14 | Medium | [ai_probe_copilot_executor.go:401]、[ai_probe_copilot_executor.go:652] | 默认 `ai check` 用固定键，当前 S1 路径会拒发；但执行器按大小写精确识别凭据键、原样传递覆盖项。Windows 环境键不区分大小写，非默认调用可漏过最后门禁。 | 显式传递认证模式，按平台归一键并拒绝歧义。
F15 | Medium | [ai.go:411]、[cli.go:847] | T3 在控制台前置拒绝时走通用错误输出；畸形、读失败等路径没有结构化 `cli_unavailable`，部分分支只在自由文本附该字样。 | 统一输出契约原因码，区分诊断和分类。
F16 | Medium | [p85_t3_tester_test.go:120]、[ai_install_proof.go:192] | 共享测试无 Windows build tag，却把盘符路径作为解码正向；Linux 的 `filepath.IsAbs` 对该路径返回 false，正向在进入变异断言前即失败。 | 使用宿主绝对路径，或将该正向移至 Windows 专属测试。
F17 | Medium | [p85_t3_linux_tester_test.go:166]、[p85_t3_tester_test.go:51]、[p85_t3_tester_test.go:72] | 同 inode 改写无条件跳过；合成 `Matches` 只测比较函数。Linux 可写目录夹具还会因非 ELF 拒绝，不能单独杀目录守卫变异；Windows 启动标记用 `.cmd`，不能证明撤守卫后原生创建成功。 | 使用实际可启动的正对照，分别实测部分 ACL、同名换对象、同 inode 改写及解释器漂移。
F18 | Medium | [runtime_agent.go:1109]、[p85_t2_frozen_policy_tester_test.go:255]、[p85_t8a_tester_test.go:496] | T2 三类制品反例和撤末端重读具备静态判别力，但只撤 Publish 的直接重读仍会经 grant helper 重读。T8a 撤 checker 门禁有可判别反例；撤执行器门禁的用例却会被宿主 T3 条件跳过。 | 对 T2 保留账本校验、撤除 Publish 的全部制品重读；给 T8a 独立可控的启动观察点。提交逐项红、恢复绿及字节复原证据。
F19 | Low | [agent_runner_policy.go:79]、[agent_runner_policy.go:315] | T2 用连字符拼接合法 ID，不同 task/step 对可碰撞并误拒绝；所谓逐字节复用比较的是重编码后的字节，不是磁盘原始字节。 | 使用无歧义键；按承诺比较原字节，或明确只比较规范语义。
F20 | Low | [trusted_ledger.go:34]、[ai_install_proof.go:117] | T3 载体路径由配置派生属实，但派生只是路径字符串；读写没有组件级链接检查，不能证明载体实际留在派生目录内。 | 用绑定目录相对读写并拒绝链接重定向。
F21 | Low | [ai_install_proof_unix.go:1]、[process_linux.go:100]、[CONTRACTS.md:220] | Windows 实现隔离正确；Linux 专属 `Stat_t` 判定用 `!windows` 编译，运行时拒绝不能替代缺该类型的平台桩。另，Linux 固定继承 fd 3 与双语契约「仅标准流句柄」字面不符，当前 S1 拒绝使其未成为实际放行。 | Linux 专属 build tag 加不支持平台拒绝桩；明确 fd 3 的契约处置后再启用 BYOK。
F22 | Info | [gate-ci.txt:1]、[selftest.txt:1] | `index/tree` 为 `TOTAL_FAIL=0`，自测 `740/740`；`scope=ci` 对应 v69 快照，不是 P85 双平台记录。包内没有本片原生 Go 测试、lint 或变异执行日志。 | 不以门禁摘要代替行为证据。本轮只读，未修改或补跑。

- **T2 — NOT-CLOSED**：独立磁盘摘要与 run/task/step/attempt 对账已接入，但唯一写者、执行者不可写及持久冻结尚未证明。
- **T3 — NOT-CLOSED**：Windows 权限判定和名称绑定不足；Linux 未证明实际解释链、同 inode 字节及载体到启动对象的连续绑定。
- **T8a — NOT-CLOSED**：默认 S1 secret 路径确实零解析、零启动，host-auth 不进 BYOK 门禁；但实际拒绝证据不只含闭集，且新原因不符合 wire Schema。
