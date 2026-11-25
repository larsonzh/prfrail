# ⑤ 预审 r1 原文（P85 · DeepSeek V4 Pro）

审查完成。所有输入已读取核验（方案书 462 行、主控注记、diff 全量、门禁读数、CONTRACTS 双语锚点、9 个实现文件实况、5 个测试文件及全部被改既有测试）。以下为最终结论。

---

`PRE-REVIEW: PASS WITH FIXES`

## 发现清单

- **F1 | Medium | `internal/adapters/p85_t3_linux_tester_test.go:99-110`** | 「同 inode 字节改写」反例（方案书 §3.6 承诺项）被写成**恒 `t.Skip`**——即使 Ubuntu VM 以非 root 运行也永不执行；Windows 侧「两次 check 之间换文件 ⇒ 预检拒绝」同样无端到端实测（仅由 P83 T3c 执行器层测试 + `Matches` 单测部分承接）。执行路径本身仍成立（不可写目录拒绝先于任何改写可能，root 改写属 S1 威胁模型外）。| 建议：⑨ VM 验证时补特权夹具实跑一次，或在报告 §8 显式豁免并载入 ⑦ 输入包，不得静默留白。
- **F2 | Low | `internal/adapters/ai_probe_copilot_executor.go:367-369`（对照 `docs/CONTRACTS.md:216`）** | 契约句「读取方 = 检查命令与**探针执行器**，每次启动前重读」——执行器 pre-spawn 只做**实时重观测**（`VerifyAIInstallSource` + `requireSameExecutableFile`），从不读载体文件。实时重验强度 ≥ 重读备忘，但字面未满足。| 建议：⑧ 前二选一——执行器补读载体并校验，或微调契约措辞为「检查命令重读、执行器实时重验」。
- **F3 | Low | `internal/adapters/ai_install_proof_unix.go:129-134`（`installSourceInterpreter` 无 shebang 即返回 nil）** | 契约要求「可执行文件为 ELF **或** 对 shebang 解释器满足不可变」；实现不校验 ELF 魔数，非 ELF 无 shebang 文件（binfmt_misc 注册格式）可通过且不记录解释器。威胁需内核级 binfmt 配置，实际探针为原生二进制，风险窄但未登记。| 建议：在报告残余章补一句（同方案书 §3.7 shebang 残余口径），或对无 shebang 文件加 ELF 魔数断言。
- **F4 | Low | `internal/console/ai_test.go:535-543`** | P83 T3c console 级「换可执行文件」反例降级为永久 `t.Skip`（74 行删除）。承接方 `TestP83T3cExecutableReplacedAfterApprovalRefused`（`p83_fixround5_tester_test.go:40`）确在且未跳，承接属实；但 console 层对该行为的变异点失去。| 建议：如实记入 ③ 报告与切片报告 §8；如可，后续补可注入安装源校验器后恢复该腿。
- **F5 | Low | ②b-3（域常量不落地）** | 方案书 §5.1-3 的 `proofrail:ai-byok-confinement:1\n` 未实现（全仓 grep 零命中），T8a 拒绝证据复用既有 T8b 闭集机制（`copilot-cli-probe-executor:1` / `availability-refusal:1`）。处置合理（新域摘要断言集有值泄漏风险，复用闭集机制反而更贴合「拒绝证据只用稳定闭集分类」），但代码注释未标「此为对方案书 §5.1-3 的偏离」。| 建议：② 报告如实申报即可（已申报则无需整改）。
- **F6 | Info | 审查输入包** | `docs/validation/evidence/P85-SECURITY-PROOF-SUBSYSTEM/archive/P85-SECURITY-PROOF-SUBSYSTEM/` 仅有 01-design / 02a 脚本 / 05-input，**③ 报告的变异红证据未随包提供**；变异真实性只能由测试注释的变异契约 + 代码推导确认（T2 两处、T3 四处、T8a 三处均为真变异：撤守卫确实使对应反例变绿）。| 建议：⑦ 输入包必须含 ③ 报告原文与红运行留痕。
- **F7 | Info | `internal/console/runtime_agent.go:785-791`** | `Review` 调用点改为 `operatorApprovedPolicyHashForAttempt`（方案书 §2.3.2 称「无需改结构」）——为满足 per-attempt 制品绑定所必需，非越界。

## 最低核验动作结论（逐条）

1. **实现与方案书一致性**：逐项 file:line 吻合（T2 制品路径/同域摘要/两端口重读/6b 冻结位点；T3 载体派生路径/pre-spawn 重验/平台 build tag；T8a 门禁位于 Resolve 前/五元断言/双保险/闭集分类）。②b-2 CreateFile 能力探测替代 `GetEffectiveRightsFromAclW` 如实且更强（直接问内核访问检查，规避 CREATOR OWNER 误报，见 windows 文件头注释）；②b-3、③b 均如实且合理（见 F4/F5）。
2. **三类反例覆盖**：T2 缺失/不可读/合形不匹配（含 Review 后篡改）四条全覆盖，断言「不 approve/零发布/证据字节不变」；T3 三反例覆盖（console 层），断言「checker 零进入/零解析/零落盘/载体不动」；T8a 证明缺失/不可读（Unknown≠放行）/五元逐一破坏全覆盖，断言「零解析零启动/拒绝证据闭集」。
3. **变异检验**：测试注释载明的 9 个变异契约与代码推导一致（撤守卫 ⇒ 反例绿）；恢复逐字节一致性由 T2/T3 载体的字节不变断言钉住。红运行留痕见 F6。
4. **契约先行**：CONTRACTS.md `:214`（闭集增 `byok_confinement_unproven`）/`:216`（T3 句）/`:220`（T8a 句）与方案书 §1.5 逐字一致；_EN 同位语义镜像（208/210/214，历史非同位行号偏移符合主控注记 a）；实现字面满足三处新句（F2/F3 为两处字面微隙）。
5. **边界**：`internal/chain/review.go`、`runtime_supervised.go` 零改动；go.mod/schemas 零改动；无新 CLI 命令；S1 `alphaSandbox` 恒 `Unknown`（`CheckSandbox` 非 Allowed 即拒）⇒ BYOK 恒拒发，与起跑消息一致；host-auth/非 cli 通道不受门禁（有负对照测试）。
6. **登记项**：见 F4/F5/F6——如实；③b 的 skip 均带明确理由，合理性成立但 F1 须闭环。

建议主控：F1 处置后（补验或豁免）即可交 ⑦；F2/F3/F4 随报告 §8 登记，不阻塞终审。
