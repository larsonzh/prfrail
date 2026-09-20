# A6 · 前置分析（① V4 Pro，Max）与裁定核准

日期：2026-09-15。来源：`deep-reasoner`（DeepSeek V4 Pro）① 调用；主控核准并核实关键代码断言。
性质：本地工作稿（`docs/t027/` git-ignored）。⑥ 验收报告的输入之一。

## 0. 结论与协议先行范围

A6 是 A0/A2/A3/A5 冻结语义之上的**纯实现切片**：无新 task/step 状态、无新 wire 记录类型、不改 `.github/workflows`（`internal/release/workflow_test.go` 冻结契约不变）。

**协议先行仅一处**：`docs/CONTRACTS{,_EN}.md` §7 新增 A6 离线运行段，固定四件事——
1. 离线证据工件布局 `<runRoot>/agent-runner-replay/runs/<requestId>/`：`manifest.pre.json`、`manifest.post.json`、`diff.json`、`logs/stdout.log`、`logs/stderr.log`、`events.jsonl`、`usage.json`；五枚冻结事实口径 = `ManifestHash`=post manifest 摘要、`DiffHash`=diff.json 摘要、`LogHash`=stdout+stderr **原始字节**拼接摘要、`UsageHash`=usage.json 摘要、`ProcessStopEvidenceHash`=`guard.TerminationEvidence` 哈希；
2. 证据缺口 fail-closed 降级表（缺口 ⇒ 不得 completed）；
3. 超时/取消语义（看门狗 + 停机验证有界化）；
4. pinned CLI 版本绑定（版本常量 + `--print-version` 预检；运行时 `executableHash` **仅作证据**，不作准入判定）。

**JSON Schema 零改动、契约夹具零改动**（`agent-runner-completion.schema.json` 已含 `logsComplete/usageComplete/outputManifestHash/processTreeStatus`，状态枚举已冻结）。**零新持久记录类型**：运行期状态在进程内 registry + A2 二文件回执 + 确定性路径推导；跨重启只依赖 A2 回执、A3 terminal 记录与 `managed-process.identity.json`。

## 1. 进程生命周期状态机（adapter 运行层；chain 状态属 A4，零改动）

| 状态 | 含义 | 持久证据 |
|---|---|---|
| S0 ABSENT | 无进程、无 identity | `launches/intent.<requestId>.jsonl` 不存在 |
| S1 STARTING | intent 已发布（A2 单赢家），spawn 已调，identity 未发布 | `launches/intent`（RecordHash） |
| S2 RUNNING | 已 spawn 且 identity 已发布 | `launches/identity.<launchId>.jsonl` + 内存 `*guard.ManagedProcess` |
| S3 STOP_REQUESTED | 停机电信号已发、grace 计时中 | 无（过渡态，重启后由 `StopProcessIdentity` 重收敛） |
| S4 STOPPED | 树已死且身份校验通过（或 `already-stopped` 幂等命中） | `TerminationEvidence`（`proofrail:termination-evidence:1`） |
| S5 COLLECTING | 停后采集 manifest/diff/logs/usage | `runs/<requestId>/` 工件（可重导出，W6） |
| S6 FINALIZED | completion 已构建并交 A3 publisher / A4 翻译 | C（replay store）+ `terminals/terminal-intent|closure` |

| # | 转移 | 触发 | 结果 / 非法性 |
|---|---|---|---|
| T1 | S0→S1 | dispatcher 发布 intent 后调 `StartAgentRunnerProcess` | Windows `CREATE_SUSPENDED`→attach Job Object→resume；Linux `Setpgid`。失败返回 error 且**契约保证未 spawn**，可同进程重试 |
| T2 | S1→S2 | Start 成功 | identity 发布失败 ⇒ `ErrAgentRunnerLaunchIdentityUnproven`（不杀、不重试、不 relaunch；A2） |
| T3 | S2→S3 | 超时 / 取消（chain Stopper）/ `Terminate` | SIGTERM(pgid) grace → SIGKILL；Windows `TerminateJobObject`；杀失败 ⇒ uncertain |
| T4 | S3→S4 | 树死 + `waitIdentityGone` | 幂等：S4→S4 = `already-stopped` |
| T5 | S2→S4 | **自然退出** | Linux 残留子树 ⇒ `ErrProcessTreeRemained`；Windows job-close 为**异步**，A6 必须按 stub 自报子 PID 集补验，**不得**直接判 completed |
| T6 | S4→S5 | 采集证据 | 任何缺口 ⇒ 结局降级（U11） |
| T7 | S5→S6 | 构建 completion + Facts → A3 → A4 → `SubmitAgentRunnerTerminal` | 重复终局幂等收敛（A3/A4 已证） |

**必须被测试钉死的非法转移**：S4/S5/S6→S2（永不 relaunch）、S0→S2（无 intent 不 spawn）、S1→S3（identity 未证不得杀）、任何 S6 出边、双 Start 同 launchID。

## 2. 崩溃窗口遍历（W1–W8）

| 窗口 | 组合 | 恢复 | 可证性 |
|---|---|---|---|
| W1 | intent 已发布、spawn 未发生 | `already-launched` 阻断，不 spawn | **可证** |
| W2 | spawn 已发生、identity 未发布 | 重启后被 A2 阻断；adapter 死亡时 job-close 兜底清树 | **不可证**（store 无“spawn 确已发生”记录；预写 state.json 同样证不了 → U10 否决新增记录） |
| W3 | RUNNING 中 adapter 崩溃 | identity 可停机（`taskkill /T /F` 或 `kill(-pgid)`）+ 租约 fencing | **可证** |
| W4 | RUNNING 中宿主崩溃/断电 | 仅知“曾 spawn”，当前存活不可知 → 不盲重启、unproven 阻断 | **不可证**（A7/B3 所有） |
| W5 | 停机 grace 中崩溃 | 重跑 `StopProcessIdentity` 收敛 | **可证** |
| W6 | 采集中途崩溃 | 采集方拥有文件（truncate+rewrite），stub 已死 ⇒ 字节确定、同 hash | **可证** |
| W7 | C 发布后、closure 前 | A3 幂等补写 | **可证** |
| W8 | dispatch 返回后、终局提交前 | step `TERMINAL_PENDING` → 终局重放/不确定暂停 | **可证** |

## 3. 证据缺口模式与 fail-closed 口径

| 工件 | 缺口模式 | 口径 |
|---|---|---|
| logs | 中断写尾缺、pipe 提前关、磁盘满、非 UTF-8、超大 | 双流各以**完整性 marker** 收尾才 `LogsComplete=true`（U8）；**原始字节直接哈希、不做解码**（非 UTF-8 不构成失败）；磁盘满/写错 ⇒ uncertain；有界截断 + `truncated=true` ⇒ `LogsComplete=false` |
| events | 无输出、损坏 JSONL、尾行撕裂 | 严格逐行 JSON 解码即完整性证明；事件不进五枚事实；completion 证据不得引用不存在的事件摘要 |
| manifest | post 采集失败、pre 不可得 | ProofRail 自己 `snapshot.Capture` 重捕；失败或父快照缺失 ⇒ diff 不可算 ⇒ uncertain |
| diff | “无变更” vs “缺失” | 无变更 = 合法结局（A5 已定）；工件生成失败 ⇒ uncertain |
| usage | 未写/写坏/磁盘满 | 严格 JSON + 闭合结构；缺失 ⇒ `UsageComplete=false` ⇒ 结算只能 `unknown`，结局非 completed |

**统一规则**：五枚事实只随 **completed** 存在（A5 已冻结）；任何缺口 ⇒ 终局降级为 uncertain / failed / cancelled 且 `Facts=nil`——**“事实不可用”不是 DTO 放空字段，而是整个终局不得 completed**。

## 4. Windows vs Unix 分异

| 维度 | Windows | Linux | 后果 |
|---|---|---|---|
| 树包含 | Job Object + `KILL_ON_JOB_CLOSE` | `Setpgid` + `kill(-pgid)` | Windows 无需树枚举 |
| 启动 | `CREATE_SUSPENDED` → attach → resume | `Setpgid:true` | 两者均“入树后才运行”，收窄 W2 |
| 树杀 | `TerminateJobObject` / `taskkill /T /F` | SIGTERM grace → SIGKILL | `Actions` 字段区分机制 |
| 退出码 | 存活判定用 259 哨兵 + StartToken 双校验 | 信号致死 `ExitCode()=-1` | **不得**以 -1 推断失败；被 stop/cancel 杀死的进程一律走 `TerminationEvidence.Outcome` |
| verifyGone | no-op（job-close 保证） | 轮询 pgid | Windows 异步 close ⇒ 需 PID 集补验（U14） |
| 时钟粒度 | ≈15.6ms | ≈1ms | 超时/取消断言一律容差（±50ms） |
| `process_unsupported.go` | — | — | **必须保持 fail-closed**；禁止加“无包含运行”回退 |

## 5. 生命周期所有权（一句话边界）

`guard` = 包含与停机**机制**；`adapter` = spawn / 证据 / 超时收尾；`chain` = 何时停与如何路由（`chain.Stopper` ← console `stopper.go` ← guard）。三者互不写对方状态。与 A2 的交互：spawn 唯一入口是 dispatcher 在 intent 发布后调用 launcher，**launcher 自己不写 intent/identity**。与 A4 的交互：`TERMINAL_PENDING` 期间进程存活由 adapter registry + 持久身份负责，chain 不读 replay store、不调度停机。

## 6. 改动点清单

**新建**：`internal/adapters/agent_runner_pinned_cli.go`（launcher，实现 `chain.AgentRunnerLauncher`）、`agent_runner_process_registry.go`（launchID→handle+phase；ProcessID 不透明编码/解码）、`agent_runner_timeout.go`（`AgentRunnerTimeoutManager.RunWithWatchdog`）、`agent_runner_evidence.go`（`CollectAgentRunnerEvidence` → `AgentRunnerEvidence`）、`agent_runner_run.go`（编排：dispatch→wait/watchdog→collect→状态映射→构建 `AgentRunnerTerminalOutcome`）及各自 `_test.go`；`tools/agent-stub/{main.go,version.go}`（确定性 stub CLI，仅 ⑤ 实验用）；⑥ 报告 `docs/validation/t027-offline-process{,_EN}.md`。

**修改**：`docs/CONTRACTS{,_EN}.md` §7（协议先行新段）；`internal/chain/agent_runner_launcher.go`（`AgentRunnerLaunchRequest` **加法式**扩展 `Command/Args/Dir/Env/Timeout/Grace/WorkspaceRoot`，replay/receipt 语义零变化）；`internal/adapters/agent_runner_replay_dispatcher.go`（填扩展字段）；`agent_runner_replay_root.go`（bootstrap 增 `runs/` 子目录）；`internal/guard/process.go`（U1 有界化，一行 + 回归测试）。

**明确出范围**：console 新命令与全链路 Engine 接线（B4）、真实候选与 capability/enforcement 准入（B1/B2）、resume 实现、新 Schema/夹具、`process_unsupported.go` 扩展、Windows 持久化证明（A7/B3）。

## 7. 测试点与 ⑤ 实验

**C.1 聚焦测试（全部进默认 `go test ./...`）**，逐条绑定变异检验目标：

| 测试 | 变异目标（删掉即必须变红） |
|---|---|
| Start 发布 identity 并运行 | `RecordLaunchIdentity` 或 launchID 比对 |
| Start 错误永不 spawn | `StartManaged` 错误路径的 close/Kill+Wait |
| stub 版本不匹配 fail-closed | `--print-version` 预检分支 |
| 经持久身份停机（Outcome=stopped、子树全灭） | `already-stopped` 或 `/T /F`/`-pgid` |
| 自然退出但子树残留 ⇒ 不得 completed | RunManaged stop 升级路径或 PID 集验证 |
| 超时 ⇒ 停机 + 结局 uncertain（Facts=nil） | `ctx.Done` 分支或 `TimedOut→uncertain` 映射 |
| 超时边界容差 | 把区间断言换成 `==` |
| cancel/timeout 竞态单赢家 | `already-stopped` 幂等或 registry 去重 |
| usage 缺失降级 | `UsageComplete=false` 分支 |
| 日志撕裂降级 | 完整性 marker 校验 |
| 非 UTF-8 日志仍成功 | 原始字节直写（改解码即红） |
| 五枚事实互异且与绑定互异 | `DisjointFrom` / 互异校验 |
| 完整闭环终局带 Facts 且 `Validate` 通过 | 任一事实摘要填充 |
| W1/W2 重放不 relaunch | dispatcher 的 already-launched/terminal 分支 |
| 停机验证有界（U1 回归） | 有界 ctx 还原为 `context.Background()` |

**C.2 ⑤ 原生实验（`//go:build a6native`，默认 `go test ./...` **不含**；轮数 env 可配、默认 20）**：E1 进程树泄漏（杀父后 stub 自报子 PID 集全部 `platformIdentityAlive=false`；最终态无容差）｜E2 超时边界（20 轮，耗时 ∈ [150ms, 600ms]，下界 ±50ms）｜E3 取消竞态（20 轮，恰一获胜方持有停机证据链）｜E4 崩溃注入（start 后 / stop 中 / 采集中各阶段 kill；不 relaunch、launch 计数=1、幂等收敛或 fail-closed）。

## 8. 裁定核准（U1–U14）

**主控核准：全部 14 项按推荐选项通过。** 其中 4 项经代码实地核实（非仅采信）：

| U# | 裁定 | 核实证据 |
|---|---|---|
| U1 | 修 `RunManaged` 的 `verifyGone(context.Background())` 为 grace 派生有界 ctx + 回归测试 | ✅ `internal/guard/process.go:129` 确为无界 `context.Background()`（同函数 103/110 行已用有界 ctx，属不一致）；Linux 轮询式 `verifyGone` 遇 D 态子进程可永久挂住 |
| U2 | 证据工件放 replay root `runs/<requestId>/` | 复用 A0/A2 路径安全与所有权检查（零新增保护面） |
| U3 | `chain.AgentRunnerLaunchRequest` **加法式**扩展 | ✅ `internal/chain/agent_runner_launcher.go` 注释明写 “Later slices extend it with argv digests, cwd, environment allowlists, and guard parameters without changing replay/receipt semantics” |
| U4 | 看门狗在 adapter 包一层（guard 保持 T007 冻结） | 失败面最小；watchdog 到期即 uncertain |
| U5 | 单测走测试二进制 re-exec helper；真 stub 仅 ⑤ 实验 | ✅ `internal/guard/process_test.go:19-27` 已有 `PROOFRAIL_PROCESS_HELPER` + `os.Args[0]` 模式，可直接沿用 |
| U6 | stub 确定性契约（`run`/`sleep`/`spawn`/`crash`/`no-usage`/`ask`/`hang`/`--print-version`；只写 cwd、无网络） | 保证测试确定性 |
| U7 | pinning = 版本常量 + `--print-version` 预检；`executableHash` 仅记录 | 避免测试绑定构建产物路径 |
| U8 | 日志以 EOF marker 证明完整；events/usage 以严格解码为完整性证明 | 与 A3/A5 严格 JSON 纪律一致 |
| U9 | ProcessID = canonical JSON `{"pid":N,"startToken":"..."}`，并镜像写 `<runDir>/managed-process.identity.json` | ✅ `internal/console/stopper.go:15` 读同名文件、字段校验 `PID>0 && StartToken!=""`，字段名须与 `guard.ProcessIdentity` 的 JSON tag 一致 |
| U10 | **不**新增 store-local run-state 记录 | state.json 预写先于 spawn，同样证不了 W2；保持“零新持久记录类型” |
| U11 | 结局映射：completed ⇔ exit 0 + 树 stopped + logs/usage 完整 + post manifest/diff 齐；非零 ⇒ failed；外部 stop/cancel ⇒ cancelled；超时/缺口/残留/采集失败 ⇒ uncertain（`Facts=nil`）；`ask` ⇒ operator-action-required；非 completed 必带 ErrorEvidence | 与 A3/A4/A5 wire 校验一致 |
| U12 | 结算走 observed（stub 用量 × **离线占位单价**） | 占位单价不得被引用为“真实计费已验证”（见 E.6） |
| U13 | 原生实验用 `a6native` build tag 隔离 + 轮数可配 | 与试点 A4 硬性约束一致，防 CI 抖动 |
| U14 | 自然退出 + 子树残留一律非 completed；Windows job-close 后须按 stub 自报 PID 集补验才可写 `ProcessTreeStatus="stopped"` | 应对 close 异步性 |

## 9. 不可证与残余边界（E1–E7）

1. **W2（spawn 后 identity 前）与 W4（宿主崩溃时的存活状态）离线不可证**；信任边界 = 不盲重启 + unproven 阻断 + uncertain 暂停，真实持久化证明归 A7/B3。
2. **进程树全集不可证**：只能验证 stub 自报子 PID 集；逃逸进程不在本片可证范围，A5 postflight 端口的重扫是后续唯一兜底。
3. **Windows job-close 异步性**：以 `VerifyWait` 上界收敛，超过即 uncertain，不静默延长。
4. **时间断言非精确**：只证区间（±50ms 下界、宽上界），E2 因而不进默认测试路径。
5. `StopProcessIdentity` 的 `taskkill` 依赖是 T007 既有实现，A6 沿用、不新建等价机制。
6. **结算单价为离线占位**，B2 前不得被任何证据引用为真实计费。
7. **A6 不证明候选兼容**：stub 行为面远小于真实候选，本片结论不得引为 B1/B2/B4 的前置事实。

---
核准记录：主控于 2026-09-15 核准 U1–U14（含 4 项代码核实），下一步进入协议先行（CONTRACTS 双语 §7 新段），随后 ② 实现。
