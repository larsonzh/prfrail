# P84 ⑦ 复审 r2 逐字转录 · 2026-11-21

- 载体：`independent-reviewer`（只读）；模型 `GPT-6.2 Sol (copilot)`，档位 Max（用户 2026-11-21 排片裁定）。
- 额度：⑦ = 1 终审 + 3 复审（本轮 = 复审 2/3）。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。

---

关键证据已对齐：本任务锁外写在实现上已消除，但 C16 未覆盖设计要求的发布入口；Windows 侧也还不能把删除 `TrimSpace` 等同于全链路保真。我会分别给出裁定和最小修复。

RE-REVIEW: FINDINGS

P84 ⑦ 原位复审 r2（2/3）。**T4 未闭合**：原 N7 的空白裁剪已修，但归属判定仍会改写非规范目标。**T5#4 的锁内写实现已闭合，入口验收未闭合**：C16 能证伪 helper 绕锁，不能证伪 `Publish` 绕过 helper。共 2 项 Medium、1 项 Low。

**逐项复核**
| 项 | 原严重度 | 状态 | 独立核对 |
|---|---|---|---|
| F1 | High | CLOSED（原反例） | [bounddir_windows.go:638](internal/evidence/bounddir_windows.go#L638) 不再让内核跟随中间重定向，而是从绑定句柄重走后缀。[W5b:318](internal/evidence/p84_bounddir_residual_test.go#L318) 区分界内两跳和第二跳越界；原 N7 反例也已闭合。剩余归属问题见 N10。 |
| F2 | High | CLOSED | [runtime_agent.go:2135](internal/console/runtime_agent.go#L2135) 先验真实历史，同进程仍严查 N−1，异进程走[锁内接管:2162](internal/console/runtime_agent.go#L2162)。[C8:708](internal/console/p84_t5h4_residual_test.go#L708) 钉住接管后的代次、持有者和摘要；恢复旧的异进程早拒会使其失败。 |
| F3 | Medium | CLOSED | 限 Rev 3b 锁。[lease.go:285](internal/guard/lease.go#L285) 释放的是 OS 锁，不删除锁文件。[C9:780](internal/console/p84_t5h4_residual_test.go#L780) 区分子进程退出后的锁回收与已提交记录导致的 `ErrLeaseHeld`。旧目录另见 N8。 |
| F4 | Medium | OPEN（验收） | [Publish:1982](internal/console/runtime_supervised.go#L1982) 和 [claim/adopt:1173](internal/console/runtime_supervised.go#L1173) 均已消除原锁外写；[C16:1141](internal/console/p84_t5h4_residual_test.go#L1141) 不覆盖调用点回退，见 N11。 |
| F5 | Medium | CLOSED | [bounddir_windows.go:169](internal/evidence/bounddir_windows.go#L169) 根枚举直接用绑定句柄；[P83 替身反例:52](internal/evidence/p83_fixround8_windows_test.go#L52) 仍保留。[W6:387](internal/evidence/p84_bounddir_residual_test.go#L387) 能抓按旧名字重开，但不能单独排除从句柄重新取名再打开。 |
| F6 | Low | CLOSED | [bounddir_windows.go:918](internal/evidence/bounddir_windows.go#L918) 为终止 NUL 留两字节。[W7:453](internal/evidence/p84_bounddir_residual_test.go#L453) 检查上界与超限拒绝；有效变异是删两字节预留，不是把 `>` 改成 `>=`。 |
| N1 | High | CLOSED | [Windows:81](internal/guard/lease_lock_windows.go#L81) 与 [Unix:50](internal/guard/lease_lock_unix.go#L50) 对稳定文件对象取排他锁，不再按 PID 清锁。[C11:903](internal/console/p84_t5h4_residual_test.go#L903) 能抓 OS 锁被移除。结论不涵盖外部替换锁文件。 |
| N2 | Medium | CLOSED | [lease_lock_windows.go:60](internal/guard/lease_lock_windows.go#L60) 打开文件不等于持锁；取锁后进程退出由内核回收，没有「已占锁、标记未落盘」状态。[C9:780](internal/console/p84_t5h4_residual_test.go#L780) 断言的是无残留持锁状态，不是无残留文件。 |
| N3 | Medium | CLOSED | [bounddir_windows.go:802](internal/evidence/bounddir_windows.go#L802) 只接受非空 `SubstituteName`；`PrintName` 仅用于诊断。[C12:616](internal/evidence/p84_bounddir_residual_test.go#L616) 的双向不一致用例能抓名称字段被换。 |
| N4 | Medium | CLOSED | [Publish:2001](internal/console/runtime_supervised.go#L2001) 恒等也批，[批内:1259](internal/console/runtime_supervised.go#L1259) 跳过已当前项。[C13:967](internal/console/p84_t5h4_residual_test.go#L967) 能抓 helper 恢复恒等早退。维持既定边界：非法同资源 proof 阻断且不改写该项，不承诺回滚此前合法写入。 |
| N5 | Medium | OPEN（验收） | [RunLocked:408](internal/guard/lease.go#L408) 与[代次变更:246](internal/guard/lease.go#L246) 共用 `withLock`。本任务是[锁内重读 current → 比较 token/hash → 落盘:2364](internal/console/runtime_agent.go#L2364)；[批处理:1230](internal/console/runtime_supervised.go#L1230) 还在锁内逐项重读目的 proof、校验历史后落盘。单项 helper 不重读目的 proof，不是对该文件做比较后交换。C14/C16 钉住各自的锁，入口缺口见 N11。 |
| N6 | Medium | CLOSED | [入口:2096](internal/console/runtime_agent.go#L2096)、[接管回调:2166](internal/console/runtime_agent.go#L2166)、[终验:2401](internal/console/runtime_agent.go#L2401) 均检查资源和 consumer。[C15:1085](internal/console/p84_t5h4_residual_test.go#L1085) 使用真实前驱与异 consumer，断言拒绝且当前 token/hash 未变。 |
| N7 | Medium | CLOSED（原反例） | [normalizeReparseTarget:861](internal/evidence/bounddir_windows.go#L861) 及[扩展前缀 helper:1168](internal/evidence/bounddir_windows.go#L1168) 均无 `TrimSpace`。[W8:677](internal/evidence/p84_bounddir_residual_test.go#L677) 检查字节保留、近名兄弟拒绝，以及界内对照读到 `BOUND`；恢复裁剪会改变这些断言。不能据此签署整个归属链保真，见 N10。 |
| N8 | Medium | CLOSED | [Windows:75](internal/guard/lease_lock_windows.go#L75) 与 [Unix:44](internal/guard/lease_lock_unix.go#L44) 识别旧目录，[lease.go:367](internal/guard/lease.go#L367) 将哨兵直接返回，不进入争用等待，也不删除。[C17:1223](internal/console/p84_t5h4_residual_test.go#L1223) 钉住哨兵、目录仍在、租约记录未写及迁移后的正例；回退到等待会失败。 |
| N9 | Low | CLOSED | 归档已含 [Rev 3b:125](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/01-arch/01-arch-design.md#L125) 和 [Rev 4:134](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/01-arch/01-arch-design.md#L134)。清单和摘要索引各列 11 项，归档 G8-b 通过；未独立重算摘要。 |

**独立发现**
| 项 | 严重度 | 位置与反例 | 最小修复 |
|---|---|---|---|
| N10 · T4 目标保真 | Medium | [bounddir_windows.go:891](internal/evidence/bounddir_windows.go#L891) 在 Ordinal 准入之前对存储目标调用 `filepath.Clean`。被折叠的点组件不会进入后置检查，原目标可被改写成另一条获准路径。最终句柄复核只验证重走后的对象，不证明原目标可无损解释。**这是静态归属误判，不是已证实读取界外字节。** | 在比较前拒绝无法无损解释的非规范目标；按支持的前缀归一化后，以原字节判定等值或严格后代。补解析层负例，保留 P52、W8 正例。 |
| N11 · T5#4 入口验收 | Medium | [Rev 4:138](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/01-arch/01-arch-design.md#L138) 要求走实际 `Publish`。[C16:1181](internal/console/p84_t5h4_residual_test.go#L1181) 直调 helper，能钉住 helper 内的锁，但即使 `Publish` 回退到旧写法也可全绿。「前段必先取锁」只适用于实际续租或接管：同进程的真实 N−1 proof 会先[追平:2142](internal/console/runtime_agent.go#L2142)，租约未到续租阈值就[直接返回:2180](internal/console/runtime_agent.go#L2180)，仍会进入本任务写分支。 | 保留 C16 作为单元反例。复用[发布入口夹具:123](internal/console/p83_fixround3_tester_test.go#L123)，补上述免续租场景：持锁时不完成且 proof 字节不变，释放后完成并绑定当前摘要；变异只绕过 `Publish` 调用点，确认即使后续批处理仍超时，字节断言也会红。 |
| N12 · 清单描述 | Low | [MANIFEST.md:15](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/MANIFEST.md#L15) 仍只列 Rev 1 和 Rev 2，落后于实际载荷。不影响 N9 的内容闭合。 | 收尾时更新描述并按原口径复核证据包。 |

**可证伪性审计**
- **锁和回写：** C9 单独抓不到「从未取锁」，需与 C11 合看；C9 的已提交子例也不是在「已提交、尚未解锁」的精确位置中断。C10 测旧的只读校验入口，C16 测新写 helper 是否等锁；两者都不能独立证明生产调用点保留锁内等值校验。
- **路径：** [P52:40](internal/adapters/bounddir_escape_windows_test.go#L40) 与 [界内对照:84](internal/adapters/bounddir_escape_windows_test.go#L84) 钉住普通 junction 的拒绝与允许；W8 钉住空白字符，不覆盖 N10。[无 glob 短路:77](internal/gates/runner.go#L77) 静态仍在，但 [W4:40](internal/gates/p83_fixround5_tester_test.go#L40) 无法感知一次失败后被忽略的根打开。

**无法验证**
已读[R9 树门禁:45](tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-postR9.txt#L45) 的 `TOTAL_FAIL=0` 和[R9 全量归档:3](tmp/P84-BOUNDDIR-RESIDUAL/test-all-postR9.txt#L3) 的 15 个包级 `ok`。未独立核验逐例 PASS/SKIP、W8/C16/C17 变异红灯及恢复、`SELFTEST 740/740`、Windows/Linux 原生竞态或旧 run 实际升级。未找到包含本轮修复的 R9 原始 diff，不能签署完整工作树范围。两项 Medium 整改后按 §5.3 原位复审，不能以全绿门禁代替。

未运行测试 / 未修改文件
