# P84 ⑦ 盲审（§7.7 必抽，终审前）逐字转录 · 2026-11-21

- 载体：`independent-reviewer`（只读）；模型 `GPT-6.2 Sol (copilot)`，档位 Max（用户 2026-11-21 排片裁定；§8.4 记录设定时刻 = 2026-11-21 本切片 ⑦ 环节）。
- 额度：盲审必抽 1 次（另计 +1、终审前执行）。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。
- 盲审纪律：输入包不含 ⑤/⑥ 审查发现。

---

W4 的现有实现确实短路，但测试只断言运行成功，感知不到一次失败后被忽略的根打开。证据核对已完成，我会把这个覆盖缺口和产品阻断项分开报告。

FINAL REVIEW: FINDINGS

**发现**
| ID | 严重度 | 位置 | 问题 | 阻断 |
|---|---|---|---|---|
| F1 | High | [internal/evidence/bounddir_windows.go:395](internal/evidence/bounddir_windows.go#L395) | 中间 reparse 的界内证明仍有 TOCTOU 窗口。 | 是 |
| F2 | High | [internal/console/runtime_agent.go:2104](internal/console/runtime_agent.go#L2104) | 同一任务带真实旧 proof 跨进程重入时，在验证前驱和尝试接管之前就被拒绝。 | 是 |
| F3 | Medium | [internal/guard/lease.go:279](internal/guard/lease.go#L279) | 代次提交后、锁目录清理前进程退出，会留下后继无法获取的租约锁。 | 是 |
| F4 | Medium | [internal/console/runtime_supervised.go:1949](internal/console/runtime_supervised.go#L1949) | 发布先写本任务 proof，后做 M1 相等校验；批内校验也未与写入串行化。 | 是 |
| F5 | Medium | [internal/evidence/bounddir_windows.go:167](internal/evidence/bounddir_windows.go#L167) | 根枚举仍按派生路径重开，不满足根分支由绑定对象直接取得枚举句柄的设计。 | 是 |
| F6 | Low | [internal/evidence/bounddir_windows.go:662](internal/evidence/bounddir_windows.go#L662) | `UNICODE_STRING` 字节长度未经上限检查就转为 `uint16`，可能将过长组件截成另一个名字。 | 否 |

**逐项证据**
1. **F1：** 普通组件确实通过 `RootDirectory` 在父句柄下解析，末级也保留 `FILE_OPEN_REPARSE_POINT` 拒绝。但 [followIntermediateReparse](internal/evidence/bounddir_windows.go#L628) 的放行依据仍是：先确认 `base` 当时指向绑定对象，再取目标句柄的最终名做 Ordinal 比较。两次查询之间没有固定目录名或祖先关系，句柄又允许 share-delete。两次观察可以各自为真，却不共同证明目标属于同一个根对象；末级复核重复了这个非原子判据。静态审查不能签署 [SECURITY.md:23](docs/SECURITY.md#L23) 所要求的 OS 级归属证明。
2. **F2：** 前任已提交 `N+1` 并退出，任务 proof 仍是可验证的 `N`。[ensurePublishLease](internal/console/runtime_supervised.go#L1134) 在有 proof 时只走 renew，[Publish](internal/console/runtime_supervised.go#L1932) 也一样。token/hash 不相等时，PID 不同直接返回，既不验证 retained 前驱，也不到 [Reclaim](internal/console/runtime_agent.go#L2120)。无 proof 的 claim 路径可以恢复，不能替代同一任务重入。部分批写也未必收敛：若本进程重入的任务已是当前代次且无须续租，[相等短路](internal/console/runtime_supervised.go#L1191) 不会补写尚未更新的兄弟任务。
3. **F3：** [withLock](internal/guard/lease.go#L262) 先用 `Mkdir` 占锁，随后 [提交当前代次](internal/guard/lease.go#L296)，最后才执行 `defer` 删除锁。进程在提交后、清理前被终止，后继 `Claim`/`Reclaim` 只能等锁或等 context 取消，进不到锁内的停机证明。这与 F2 不同：即使重入选对入口，也仍可能卡在锁外。
4. **F4：** [M1](internal/console/runtime_supervised.go#L1242) 对读到的 `Current` 正确比较了 token 和 hash，但发生在本任务 proof 已写之后。批内也只读一次 [Current](internal/guard/lease.go#L221)，之后写入并不受同一资源锁保护。因而 M1 能拒绝一部分过时代次，不能保证整条发布路径零写入或批写期间代次不变。[发布末尾复验](internal/console/runtime_supervised.go#L1960) 仍能拒绝过期 proof，没有证据表明这会生成 `completed` receipt。
5. **F5：** `""` 和 `"."` 现在确实同路，`readDirNames` 也确实读取获准的句柄。但根分支先从绑定句柄取名字，再用 `CreateFile` 重开。后续卷号和文件索引比较能拒绝已打开的替身，不能把按名重开变成按对象打开。重复枚举需要独立游标，这不要求重新解析名字。
6. **F6：** [相对路径清洗](internal/evidence/bounddir_clean.go#L21) 不检查组件的 UTF-16 字节长度。原生 API 按计数字符串解析，不能指望它拒绝 Go 侧已经截短的长度。常规 glob 使用枚举得到的文件名，本轮未发现该路径会产生这种超长组件。

**可证伪性审计**
- **W1/W2：** [W1 字符串用例](internal/evidence/p84_bounddir_residual_test.go#L39) 能区分 Ordinal 与 `EqualFold`；[孪生目录用例](internal/evidence/p84_bounddir_residual_test.go#L74) 可因环境跳过。[W2 出界再返回](internal/evidence/p84_bounddir_residual_test.go#L183) 能区分中间校验与末级校验，但不覆盖并发改名；其 [枚举用例](internal/evidence/p84_bounddir_residual_test.go#L172) 实际命中的是末级 reparse 拒绝。P52 的界内读取和 P83 的末级拒绝有各自独立的断言。
- **W3/W4：** [W3](internal/evidence/p84_bounddir_residual_test.go#L246) 锁定两种根拼写和文件拒绝，不区分直接用绑定对象还是按名重开。[W4](internal/gates/p83_fixround5_tester_test.go#L37) 锁定无产物时运行成功，但单删 [Runner 的预绑定条件](internal/gates/runner.go#L77) 仍可通过：多余打开失败会被忽略。当前源码确实不打开根，测试尚不能独立证明这一点。
- **C1/C2：** [C1](internal/console/p84_t5h4_residual_test.go#L127) 能抓仅按资源覆盖的回退，但单删 claim 比较仍可能被 retained token 比较兜住。[C2](internal/console/p84_t5h4_residual_test.go#L199) 的真实两代历史能抓退回只认 `N-1`。逐条验证会报错，不会跳过同资源的非法 proof；允许此前已验证的文件先落盘，不等于整批回滚。proof 是 run 级凭证，不证明内容完全相同的合法 proof 属于某个 task。
- **C3/C6/C7：** [C3](internal/console/p84_t5h4_residual_test.go#L268) 直接走 claim 和批写 helper，子进程在续租返回后退出，且允许 skip，因此不覆盖 F2、F3。[C6](internal/console/p84_t5h4_residual_test.go#L357) 钉住读取 `Current` 前已推进的情形，不覆盖发布预写或读后竞态。[C7](internal/console/p84_t5h4_residual_test.go#L423) 证明 65 跳拒绝、1 跳成功，不锁定第 64 跳必须成功。

**边界与未验证**
- Windows `class=10`、`FileName` 偏移断言和 NTSTATUS 非成功拒绝均已落实；常规顺序路径未见句柄泄漏。[Close](internal/evidence/bounddir_windows.go#L752) 未与走链互斥，不据此签署并发关闭安全。归档包级读数也不能证明多缓冲区枚举、异常记录布局或竞态行为。
- 未见 P83 的 holder、锁内 Reclaim、资源归属和发布复验被放宽。Unix 有预期的根拼写行为改动，不能称绝对零功能变化；未见 `openat2` 防线变化，Darwin 仍拒绝。归档 diff 的六个路径未涉及 schema、CLI、依赖或 `go.mod`。
- 已从现状源码确认 M1。当前模式仅允许读和搜索，未执行实时 `git diff`，不能确认最新工作树的完整变更集。[树门禁](tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-post3b.txt) 为 `TOTAL_FAIL=0`，[测试归档](tmp/P84-BOUNDDIR-RESIDUAL/test-all-post3.txt) 为包级 `ok`；未核实最新版本的逐例 PASS/SKIP、Windows 实机竞态或变异结果。
- 未打开⑤/⑥报告或清单。指定测试源码的注释含历史审查描述，不能记录为输入完全脱敏；以上判定以实现和断言为依据。

未运行测试 / 未修改文件
