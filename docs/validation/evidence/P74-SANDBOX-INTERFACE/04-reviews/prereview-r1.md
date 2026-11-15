# ⑤ 预审 r1 原文（P74）

- 载体：`deep-reasoner` + `DeepSeek V4 Pro (deepseek)`
- 调用：1（额度 1 + 1）
- 结论：`PRE-REVIEW: PASS`

---

PRE-REVIEW: PASS

**A. 接口 vs CONTRACTS §7.2 (a)–(h) 逐句对应**：对应成立。(a) `internal/guard/sandbox.go:80` 方法集仅 `Check(ctx, CheckSpec) (Verdict, error)`，无生命周期方法；`process.go` 零 diff，无类型重叠。(b) 三态齐备，`CheckSandbox` 将 Denied/Unknown 统一折为 `ErrSandboxRejected`，仅 Allowed 放行。(c) 两哨兵定义与映射正确。(d)(e)(f)(g)(h) 与实现一致（详见 B/C/F）。

**B. 检查点时序与 permit 兼容**：无发现。`runtime_agent.go:199-209` 与 `runtime_supervised.go:311-324` 均位于 `requireProbePermit`（`readiness_gate.go:14` → `ProbePermit.Require`，`probe_permit.go:210` `persistClaim` O_EXCL 为线性化点、配额不可退）与 Step 2 marker 之前；绿测以「拒绝后二次 `requireProbePermit` 成功 + marker `ErrNotExist`」实证不烧配额。Step 0 就绪门（仅存在性评估、不烧配额）先于检查运行，兼容。

**C. 双平台 fail-closed 一致性**：无发现。windows α 与 !windows 均恒返 Unknown；`CheckSandbox`（`sandbox.go:89`）统一折算拒绝，生产中唯一调用链为 `preflightSandbox` → `CheckSandbox`，不存在「未检查却 Allowed」路径；零值 `Verdict{}` 落在 Unknown（`VerdictUnknown = iota`），fail-closed 安全默认；nil Sandbox 回退平台实现且拒绝（`sandbox_interface_test.go:74` 已测）。

**D. 红测可证伪性**：成立。`mutation-agent-runner.txt` 删点后三例全红（err 得 `<nil>`）；`mutation-supervised.txt` 三例全红（err 得 visible-delivery 拒绝）。绿测断言覆盖：哨兵错误、零 summary、fake 恰 1 次且 reason 正确、Dispatcher 0 次、marker/ledger 不存在、配额可再认领——六项断言均存在且绿。

**E. 默认零行为变化**：无发现。未设/off ⇒ `sandboxEnforced=false`，两调用点仅多一次布尔比较；`executeRun` 每 run 多一次 `os.LookupEnv`（零副作用）；`cli.go:305-311` 读取位于 `LoadChainConfig`（首个文件读）之前，非法值仅走 `writeUsageError` 输出。

**F. 越界检查**：无发现。暂存文件集 = docs + `internal/console/**` + `internal/guard/sandbox*.go`；无 `internal/chain/**`、`internal/adapters/**`、schemas/fixtures/CI、s1-exit 改动；无真实 enforcement 与进程生命周期。

---

**Low**

1. `internal/guard/sandbox.go:94`：`CheckSandbox` 的 `err != nil` 分支丢弃了 Verdict.Reason——!windows 路径「sandbox check not supported on this platform」这一诚实理由不进入审计链，与 (b)「reason 随拒绝错误进入审计与诊断信息」有微差。建议错误文本补上 `verdict.Reason`（判空）。
2. `internal/guard/sandbox.go:94`：`Check` 直接返回 `ErrSandboxUnavailable` 时错误串自包裹为「guard: sandbox check unavailable: guard: sandbox check unavailable」。建议去重（`errors.Is` 相等时直接返回原错误）。
3. 无 windows 侧契约测试：α 契约（Unknown + nil error + 非空理由）仅经 `sandbox_interface_test.go` 的 nil 回退在 Windows CI 间接覆盖。建议补 `sandbox_windows_test.go` 直接钉死。

**Info**

4. 变异输出红在第一条错误断言处（`t.Fatalf` 短路），marker/配额断言未在变异轮被独立触发——属断言顺序固有行为，可证伪性结论不受影响。
5. `design-r1.md` §9.3 预测删点后 supervised「parks in the waiting state with nil error」，实测变异输出为 visible-delivery 拒绝——预测机理与实际不符，结论（变红）一致。
6. `red-agent-runner.txt`/`red-supervised.txt` 与同名 `mutation-*.txt` 内容完全一致（同一代码状态：调用点缺失）。③ 红与②b 后变异从内容不可区分，建议证据包记录捕获时序（文件时间戳）以便⑥/⑦复核 TDD 顺序。
