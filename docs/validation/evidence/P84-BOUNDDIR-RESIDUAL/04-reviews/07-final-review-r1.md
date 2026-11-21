# P84 ⑦ 终审（FINAL REVIEW）r1 逐字转录 · 2026-11-21

- 载体：`independent-reviewer`（只读）；模型 `GPT-6.2 Sol (copilot)`，档位 Max（用户 2026-11-21 排片裁定）。
- 额度：⑦ = 1 终审 + 3 复审（本轮 = 终审，用 1/1）；盲审另计 +1 已在前（`07-blind-review.md`）。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。

---

FINAL REVIEW: FINDINGS

P84 ⑦ 终审第 1 轮。**T4、T5#4 均未闭合。** F5 的实现偏差可接受，F6 已闭合；其余阻断项见下表。

**盲审复核**
| 项 | 裁定 | 证据 |
|---|---|---|
| F1 | OPEN | [走链](internal/evidence/bounddir_windows.go#L642) 已改为读取 reparse 数据，从绑定句柄重走后缀，不再让内核跟随后事后核名。[W5b](internal/evidence/p84_bounddir_residual_test.go#L314) 有两跳界内正例和二跳越界反例，但不区分显示名与实际目标。当前实现仍可能校验错目标，见 N3。 |
| F2 | OPEN | [失配分支](internal/console/runtime_agent.go#L2128) 先验真实历史，另一持有者走 `Reclaim`。[C8](internal/console/p84_t5h4_residual_test.go#L613) 能区分恢复旧拒绝分支的变异。原跨进程入口缺口已修，但部分批写后的同进程重入不保证收敛，见 N4。 |
| F3 | OPEN | [Windows](internal/guard/lease_lock_windows.go#L33) 与 [Unix](internal/guard/lease_lock_unix.go#L20) 已接入活性检查；[Windows 标记读取](internal/guard/lease_lock_windows.go#L68) 也启用了 delete sharing。[C9](internal/console/p84_t5h4_residual_test.go#L673) 能发现删掉清锁调用，但只覆盖一个竞争者遇到完整死 PID 标记，不覆盖 N1、N2。 |
| F4 | OPEN | [本任务写前检查](internal/console/runtime_supervised.go#L1960) 顺序已修正；[C10](internal/console/p84_t5h4_residual_test.go#L718) 能抓等值判断被删，[C6](internal/console/p84_t5h4_residual_test.go#L358) 能抓批写开始前代次已推进。两者不覆盖检查后、写入前的推进，见 N5。 |
| F5 | CLOSED | **② 偏差 1 可接受。** [生产根分支](internal/evidence/bounddir_windows.go#L169) 使用绑定句柄，不进入仅为 `owned` 句柄执行的复核；[保留的根 helper 分支](internal/evidence/bounddir_windows.go#L365) 仍供 [P83 冻结反例](internal/evidence/p83_fixround8_windows_test.go#L109) 校验替身句柄与错误文本。保留不恢复按名重开。 |
| F6 | CLOSED | [长度守卫](internal/evidence/bounddir_windows.go#L907) 按 UTF-16 字节计数，为 `MaximumLength` 的终止 NUL 留出两字节。[W7](internal/evidence/p84_bounddir_residual_test.go#L444) 同时断言边界值和超限在打开前被拒绝。 |

**独立发现**
| ID | 严重度 | 位置与反例 | 最小修复 |
|---|---|---|---|
| N1 | High | [lease.go:397](internal/guard/lease.go#L397) 的清锁函数先按名字判死，再执行 [重命名](internal/guard/lease.go#L413)：`if err := os.Rename(lockPath, quarantine); err != nil {`。两个竞争者读到同一陈旧标记时，落后的检查可能移走已重新取得的**活锁**。移走后才复核，回迁失败又被忽略；锁名空缺期间可再次进入临界区，破坏 `Reclaim` 的单写保证。 | 让占锁、释放与进程退出由同一 OS 锁对象保证，或使用等价的原子锁身份判据。不得先移走活锁再尝试复原。 |
| N2 | Medium | [lease.go:263](internal/guard/lease.go#L263) 的 [占锁点](internal/guard/lease.go#L276) 是 `if err := os.Mkdir(lockPath, 0700); err == nil {`。建目录成功但标记尚未落位时进程终止，后继因无标记永远无法判死；超时只中止本次等待，不解锁。C9 预先写好了标记。 | 与 N1 一并消除“已占锁、无可恢复身份”状态；不要把缺标记直接判为死亡。 |
| N3 | Medium | [bounddir_windows.go:746](internal/evidence/bounddir_windows.go#L746) 解析 reparse 数据后，在 [目标选择](internal/evidence/bounddir_windows.go#L799) 使用 `target := printName`。`PrintName` 是显示名，`SubstituteName` 才是实际目标。两者不一致时，可以把实际应拒绝的重定向解释成界内路径，或读取另一份界内产物。这里**未证实直接读取界外字节**，但 P52 的目标归属判定不成立。 | 只用合法的 `SubstituteName` 作归属判定；`PrintName` 仅用于诊断。补两者不一致的解析和拒绝反例。 |
| N4 | Medium | [runtime_supervised.go:1190](internal/console/runtime_supervised.go#L1190) 在 `if refreshed == previous {` 直接返回，[`Publish` 也只在 proof 变化时补写](internal/console/runtime_supervised.go#L1941)。若部分批写留下主任务已是当前代次、兄弟任务落后两代，主任务重入不补批，兄弟晋升又被[严格 N−1](internal/console/runtime_agent.go#L2132) 拒绝。Rev 2 的“各自晋升时追平”对此不成立。 | 即使本任务 proof 未变，也允许已验证的当前持有者重跑幂等批量对账；保留 P83 的严格 N−1 判据。 |
| N5 | Medium | [runtime_supervised.go:1190](internal/console/runtime_supervised.go#L1190) 的批处理在 [读取当前记录](internal/console/runtime_supervised.go#L1232) 后无锁地复用该记录并[写入](internal/console/runtime_supervised.go#L1253)。代次若在两者间推进，旧批可回写过时 proof，甚至覆盖已刷新的 proof。[发布末验](internal/console/runtime_supervised.go#L1975) 能拒绝它看到的失配，但不能撤销先前写入。 | 将等值校验、目标 proof 重读和落盘放在与代次推进共用的资源锁内，覆盖本任务和兄弟项。 |
| N6 | Medium | [runtime_agent.go:2074](internal/console/runtime_agent.go#L2074) 的新增重入路径没有沿用 [`claimPublishLease` 的 run/consumer 守卫](internal/console/runtime_agent.go#L1961)。当前记录若同资源、拥有真实 retained 链，但 `ConsumerID` 不是发布 consumer，[`Reclaim` 回调](internal/console/runtime_agent.go#L2155) 只问旧持有者是否停止。真实历史不是该 consumer 的发布授权。 | 续租和发布复验统一检查 `publishLeaseBelongsToRun`，并在锁内对实际被接管的记录重复检查。 |

**整体判定**
- **T4：** 普通组件由 [父句柄约束](internal/evidence/bounddir_windows.go#L944)，末级 reparse 拒绝和 P83 替身目录拒绝仍在。根两种拼写同路，[无 glob 不绑定工作区](internal/gates/runner.go#L77)；Linux 的 [openat2 限制](internal/evidence/bounddir_unix.go#L65) 与 Darwin 的[拒绝行为](internal/evidence/bounddir_darwin.go#L62) 未见放宽。N3 仍阻断 [SECURITY.md:23](docs/SECURITY.md#L23) 所要求的归属证明，也不能签署完整的 [P52 界内允许、界外拒绝](internal/adapters/bounddir_escape_windows_test.go#L40)。
- **T5#4：** [逐项历史校验](internal/console/runtime_agent.go#L2247) 与 [retained 原件校验](internal/guard/lease.go#L485) 已落实；无记录或异资源略过，同资源假前驱拒绝且不改写该项。同进程真 N−1、跨进程真前驱两条路径已具备；第三条部分批写恢复尚不闭合。P83 异资源、假前驱和活前任拒绝未见放宽，但 N1、N6 使所有权结论不能签署。受审改动可在现有[恢复契约](docs/CONTRACTS.md#L91)和[发布契约](docs/CONTRACTS.md#L150)内修正，未见需要新增 schema、CLI 或依赖。

**可证伪性审计**
- **W1/W4/W6/W7：** Ordinal 改回 `EqualFold` 可被 [W1](internal/evidence/p84_bounddir_residual_test.go#L48) 区分；[W4](internal/gates/p83_fixround5_tester_test.go#L40) 只断言运行成功，可能看不到一次失败后被忽略的根打开。[W6](internal/evidence/p84_bounddir_residual_test.go#L383) 能抓按旧名字重开，但 `currentName()` 可随改名找到原对象，精确回退 pre-F5 仍可能绿，与主控所报 `GREEN` 相容。[W7 注释](internal/evidence/p84_bounddir_residual_test.go#L442) 的 `>` 改 `>=` 不会变红：UTF-16 字节数恒为偶数，阈值为奇数；删掉两字节预留才是有效变异。
- **C1/C2/C3/C7–C10：** C1/C2 能区分仅比资源及只认 N−1 的**批写**变异；[C3](internal/console/p84_t5h4_residual_test.go#L239) 不等于真实提交点中断；[C7](internal/console/p84_t5h4_residual_test.go#L424) 不钉第 64 跳必须成功。C8 不证实部分批写收敛，C9 不证实竞争清锁安全，C10 单独绕过 `Publish` 调用点仍可绿。复审应补确定性的断点覆盖，并同时断言代次、持有者及磁盘 proof 字节；仅凭压力测试全绿不能排除这些交错。

**无法验证的部分**
读取的[树门禁](tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-postR6c.txt#L45) 为 `TOTAL_FAIL=0`；[全量测试归档](tmp/P84-BOUNDDIR-RESIDUAL/test-all-postR6.txt) 显示包级 `ok`。未独立核实逐例 PASS/SKIP、`SELFTEST 740/740`、九项变异及 SHA256 复原、`16×100` 压测，或 Windows/Linux 的实际竞态。首轮按[审查隔离](docs/DELIVERY_DIRECTIVE.md#L713)未读取⑤/⑥原文。现有[差异快照](tmp/P84-BOUNDDIR-RESIDUAL/05-review/git-diff-postL.txt)早于盲审整改；只读载体不能执行实时 `git diff`，全工作树无越界须待最新静态 diff 和新增路径清单核实。

未运行测试 / 未修改文件
