# P84 ⑦ 复审 r1 逐字转录 · 2026-11-21

- 载体：`independent-reviewer`（只读）；模型 `GPT-6.2 Sol (copilot)`，档位 Max（用户 2026-11-21 排片裁定）。
- 额度：⑦ = 1 终审 + 3 复审（本轮 = 复审 1/3）。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。

---

静态核查已完成。新锁在稳定的锁文件对象上消除了旧的清锁竞态；本任务 proof 的锁外写入仍存在，Windows 目标归一化也还有一个独立的误判反例。

RE-REVIEW: FINDINGS

P84 ⑦ 原位复审 r1：**T4、T5#4 均未闭合。** T4 已消除让内核跟随中间 reparse 后再核名的旧路径，但目标归一化仍可把界外名称判成界内。T5#4 的兄弟 proof 批写已锁内串行化，但 `Publish` 写本任务 proof 仍绕过同一把锁。归档门禁全绿不能排除这两个反例。

**逐项复核**
| 项 | 原严重度 | 状态 | 独立核对 |
|---|---|---|---|
| F1 | High | OPEN | [走链:638](internal/evidence/bounddir_windows.go#L638) 不再让内核跟随中间重定向，而是从绑定句柄重走后缀；[W5b:318](internal/evidence/p84_bounddir_residual_test.go#L318) 区分界内两跳和第二跳越界。但后缀准入仍受 N7 影响。 |
| F2 | High | CLOSED | [历史校验:2135](internal/console/runtime_agent.go#L2135) 后，异进程持有者走[锁内接管:2162](internal/console/runtime_agent.go#L2162)；同进程仍严查 N−1。[C8:708](internal/console/p84_t5h4_residual_test.go#L708) 核对接管后的代次、持有者和摘要；[C13:967](internal/console/p84_t5h4_residual_test.go#L967) 核对主项已当前时兄弟落后两代仍能追平。 |
| F3 | Medium | CLOSED | **限 Rev 3b 创建的锁。** [C9:780](internal/console/p84_t5h4_residual_test.go#L780) 的子进程持锁退出，不执行释放代码；锁文件仍在，后继能取得锁。已提交记录则使后继快速返回 `ErrLeaseHeld`。旧格式遗留见 N8。 |
| F4 | Medium | OPEN | [写前相等校验:1971](internal/console/runtime_supervised.go#L1971) 已补，但[本任务落盘:1979](internal/console/runtime_supervised.go#L1979) 仍在取锁前。C10 和 C14 都不经过此写点。 |
| F5 | Medium | CLOSED | [生产根分支:169](internal/evidence/bounddir_windows.go#L169) 直接枚举绑定句柄，不按名重开；[P83 冻结反例:109](internal/evidence/p83_fixround8_windows_test.go#L109) 所需的替身判据仍保留。 |
| F6 | Low | CLOSED | [长度守卫:915](internal/evidence/bounddir_windows.go#L915) 为终止 NUL 留两字节。[W7:453](internal/evidence/p84_bounddir_residual_test.go#L453) 检查最大可表示值和超限拒绝；去掉预留会使超限断言失败。 |
| N1 | High | CLOSED | 在核心 store 路径和锁文件对象不被替换的前提下，[Windows:74](internal/guard/lease_lock_windows.go#L74) 锁定同一对象的同一字节范围，[Unix:40](internal/guard/lease_lock_unix.go#L40) 锁定同一 inode。正常释放只解锁、关句柄，不重命名或删除；进程退出由内核回收持锁状态，无 PID 标记、无 heal。[C11:903](internal/console/p84_t5h4_residual_test.go#L903) 区分活持有者阻塞与移除 OS 锁后的直接成功。 |
| N2 | Medium | CLOSED | [打开锁文件:60](internal/guard/lease_lock_windows.go#L60) 不等于已占锁。进程若死在打开与锁定之间，没有持锁状态；死在取得锁之后由 OS 回收。C9 断言的是**无残留持锁状态**，不是无残留文件。 |
| N3 | Medium | CLOSED | [目标选择:802](internal/evidence/bounddir_windows.go#L802) 只接受非空 `SubstituteName`；`PrintName` 不参与选择。[C12:616](internal/evidence/p84_bounddir_residual_test.go#L616) 的两种反向样例均可抓到改读 `PrintName`。这不覆盖 N7 的归一化。 |
| N4 | Medium | CLOSED | [Publish:1989](internal/console/runtime_supervised.go#L1989) 恒等也批，[已当前项:1253](internal/console/runtime_supervised.go#L1253) 跳过。C13 能抓恢复 helper 早退。**接受 fail-closed 副作用**：可读异资源略过；同资源但非真实历史、或不可读的 proof 阻断晋升。坏项不改写，先前已验证项允许落盘，不承诺整批回滚。 |
| N5 | Medium | OPEN | [RunLocked:375](internal/guard/lease.go#L375) 与代次推进共用 `withLock`；[批内:1223](internal/console/runtime_supervised.go#L1223) 覆盖当前记录读取、等值校验、逐项重读、历史校验和落盘。锁外只收集路径本身不是缺陷；**本任务的独立写路径不在这个临界区内**。 |
| N6 | Medium | CLOSED | [入口:2096](internal/console/runtime_agent.go#L2096)、[接管回调:2166](internal/console/runtime_agent.go#L2166)、[终验:2366](internal/console/runtime_agent.go#L2366) 均查资源与 consumer。[C15:1085](internal/console/p84_t5h4_residual_test.go#L1085) 使用异 consumer、死持有者和真实历史，断言拒绝且当前 token/hash 未变。 |

**独立发现**
| 项 | 严重度 | 位置与反例 | 复审要求 |
|---|---|---|---|
| N5 / F4 残余 | Medium | [runtime_supervised.go:1979](internal/console/runtime_supervised.go#L1979)：A 校验 `N` 后暂停，B 取得资源锁并推进至 `N+1`，A 仍可把 `N` 写回本任务，甚至覆盖 B 已刷新的 proof。后续批写会拒绝 A，但不会撤销写入。无 proof 的 claim/adopt 分支也在[批锁释放后插入本任务 proof:1167](internal/console/runtime_supervised.go#L1167)。 | 本任务 proof 的重读、校验和落盘须与同资源代次推进串行化。用实际 `Publish` 入口验证持锁超时后 **proof 字节不变**；只断言没有 `completed` 不够。 |
| N7 | Medium | [bounddir_windows.go:852](internal/evidence/bounddir_windows.go#L852)：`TrimSpace` 会删除合法文件名末尾的 `U+00A0`。若 SubstituteName 指向末尾多该字符的界外兄弟目录，裁剪后可恰好等于绑定根，[后缀判定:877](internal/evidence/bounddir_windows.go#L877) 因而放行。最终句柄校验检查的是重走所得的界内对象，无法发现目标被换了。**这是误放行和可能读错产物，不是已证实读取界外字节。** | 保留有身份意义的字符，无法无损解释的目标拒绝；补普通界内正例与合法近名界外负例，断言实际打开行为。 |
| N8 | Medium | [基线差异:1524](tmp/P84-BOUNDDIR-RESIDUAL/05-review/git-diff-postR8.txt#L1524) 显示旧版在同一路径 `Mkdir`。旧写者崩溃留下目录后，新版按文件打开失败，[acquireLeaseLock:334](internal/guard/lease.go#L334) 把该错误直接返回，进不到租约恢复。[升级流程:194](docs/INSTALLATION.md#L194) 保留 state/store，未排除这一旧状态。 | 明确旧 run 的停机后迁移或拒绝边界，并用旧目录锁状态取证。不能靠在线盲删目录解决。 |
| N9 | Low | [归档方案:118](docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/01-arch/01-arch-design.md#L118) 仍止于已废弃的 Rev 3 锁方案；[Rev 3b:125](tmp/P84-BOUNDDIR-RESIDUAL/01-arch-design.md#L125) 只在临时方案中找到。 | 收尾时归档实际采用的方案，不把临时件当长期权威。 |

**可证伪性审计**
- **锁：** C11 能抓无 OS 锁；C9 单独抓不到，因为从未持锁也能通过死亡后获取的断言。C9 的已提交记录子例是在 `Claim` 返回、释放锁后退出，不是同一次「已提交、尚未解锁」中断。锁绑定文件对象而非路径；[Windows 允许 delete sharing:64](internal/guard/lease_lock_windows.go#L64)，故不能把内部不删锁推广成外部替换锁对象也安全。
- **回写：** C10 能抓删等值判断；C14 能抓批 helper 绕锁；C13 能抓 helper 恢复恒等早退。三者均不能抓 `Publish` 的锁外写入，C13 也不能抓调用方恢复条件调用。C15 能抓同时删掉入口和接管回调的 consumer 检查，不分别钉住回调与终验。
- **边界：** [P52:40](internal/adapters/bounddir_escape_windows_test.go#L40)、W2 和 W5b 区分普通界内外目标，C12 区分名称字段，均未覆盖 N7。根拼写和末级拒绝未见回退；[无 glob 分支:77](internal/gates/runner.go#L77) 静态仍短路，但 [W4:40](internal/gates/p83_fixround5_tester_test.go#L40) 的成功断言不能感知一次被忽略的根打开。

**无法验证的部分**
已读取[树门禁:45](tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-postR8.txt#L45) 的 `TOTAL_FAIL=0`，以及[全量归档:7](tmp/P84-BOUNDDIR-RESIDUAL/test-all-postR8.txt#L7) 的包级 `ok`。未独立核实逐例 PASS/SKIP、变异红灯及恢复、`SELFTEST 740/740`、`16×100` 压测、原生竞态和旧状态升级实测；未获取实时 `git diff` 或核验全工作树清单。上述阻断项是静态判定，不是本轮实测结果。

未运行测试 / 未修改文件
