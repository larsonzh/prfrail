# ⑥ 独立扫描 r1 原文（P74）

- 载体：`independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`
- 调用：1（额度 ≤3）
- 结论：`INDEPENDENT SCAN: PASS`（Section B 无发现）；摘要句为「依据契约与主控落盘静态材料，本片满足 §7.2 的接口、时序、fail-closed 语义与可证伪性要求」。

---

## Section A

INDEPENDENT SCAN: PASS

## Section B

无发现

## Section C

- 读取主控落盘的 [docs/CONTRACTS.md](docs/CONTRACTS.md) 与 [docs/CONTRACTS.md](docs/CONTRACTS.md#L321-L359)：§7.2 的八条契约彼此一致，且与实现对齐：
  - `guard.Sandbox` 只做 `Check(ctx, CheckSpec) (Verdict, error)` 判定，不提供生命周期方法，符合 (a)。
  - `Allowed / Denied(reason) / Unknown(reason)` 三态都按 fail-closed 处理，符合 (b)。
  - `ErrSandboxUnavailable` / `ErrSandboxRejected` 均不得被解释为放行，符合 (c)。
  - 检查点放在 probe permit 认领之前，且在 run-root marker 之前，符合 (d)。
  - 触发源仅是 `PRFRAIL_SANDBOX_CHECK`，默认 `off`，符合 (e)。
  - 非 Windows / Windows S1 都返回 `Unknown`，不允许谎报 `Denied`，符合 (f)。
  - 边界不碰 chain/adapters 判定本体，也不落真实 enforcement，符合 (g)。
  - B 段不解锁 S-G1/S-G8，符合 (h)。

- 读取主控落盘的 [internal/guard/sandbox.go](internal/guard/sandbox.go) 与 [internal/guard/sandbox_unsupported.go](internal/guard/sandbox_unsupported.go)：`CheckSandbox` 统一把 `nil` / 非 `Allowed` / `err != nil` 全部包装成 fail-closed 语义；`Unsupported` 平台返回 `VerdictUnknown + ErrSandboxUnavailable`，与契约一致。

- 读取主控落盘的 [internal/console/runtime_agent.go](internal/console/runtime_agent.go) 与 [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go)：两条派发路径都在 permit 之前调用 `preflightSandbox`，且在 run-root marker / operator ledger 之前返回，符合“环境/配置类拒绝不得烧配额、不建 marker、不 spawn”的时序。

- 读取主控落盘的 [tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/red-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-agent-runner.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/red-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-supervised.txt) 与 [tmp/P74-SANDBOX-INTERFACE/03-tests/green-guard-contract.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/green-guard-contract.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/green-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/green-agent-runner.txt)、[tmp/P74-SANDBOX-INTERFACE/03-tests/green-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/green-supervised.txt)：红测先红、绿测后绿，说明 TDD 方向正确，且检查点确实位于 permit 之前。

- 读取主控落盘的 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt) 与 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt)：删掉调用点后，两个测试都从绿转红；这说明调用点不是“无意义装饰”，而是真实的安全断言。

- 读取主控落盘的 [tmp/P74-SANDBOX-INTERFACE/gate-index.txt](tmp/P74-SANDBOX-INTERFACE/gate-index.txt)、[tmp/P74-SANDBOX-INTERFACE/gate-tree.txt](tmp/P74-SANDBOX-INTERFACE/gate-tree.txt)、[tmp/P74-SANDBOX-INTERFACE/gate-ci.txt](tmp/P74-SANDBOX-INTERFACE/gate-ci.txt)、[tmp/P74-SANDBOX-INTERFACE/gate-selftest.txt](tmp/P74-SANDBOX-INTERFACE/gate-selftest.txt) 与 [tmp/P74-SANDBOX-INTERFACE/04-full-test.txt](tmp/P74-SANDBOX-INTERFACE/04-full-test.txt)：门禁四项均为 PASS，且全量测试输出里 `internal/guard` / `internal/console` 均为 `ok`，未见逆转证据。

## Section D

1. 删掉 [internal/console/runtime_agent.go](internal/console/runtime_agent.go) 的检查块
   - 会红：`TestAgentRunnerRunSandboxRefusesBeforePermit` 的三个子case（denied / unknown / unavailable）要求 `ErrSandboxRejected` 或 `ErrSandboxUnavailable`，而 mutation 输出中 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt) 显示 “got <nil>” 和 “the dispatch must be refused with ...”。
   - 不会红：仅看“CLI 能起”或“build 能过”的测试；这类断言不覆盖“在 permit 前拒绝并不烧配额”的安全语义。
   - 结论：删掉检查块后，安全断言会红，且这不是“假红”。

2. 删掉 [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go) 的检查块
   - 会红：`TestSupervisedRunSandboxRefusesBeforePermit` 要求在 freshStart 分支中先被 sandbox 拒绝；mutation 输出 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt) 明确显示 “got supervised visible delivery is not accepted …”，说明代码已越过安全门而进入 visible delivery。
   - 不会红：“重入路径不检查、不消费配额” 的正向路径（非 freshStart）不会被这条测试覆盖到；它被设计为仅在 freshStart 上检查。
   - 结论：删掉该检查块后，fresh-start 的安全断言会红；“再入路径不检查”本身仍需单独证明。

3. 删掉 [internal/guard/sandbox_unsupported.go](internal/guard/sandbox_unsupported.go) 的 fail-closed 返回
   - 会红：`TestSandboxInterfaceContract` 明确要求 `NewSandbox()` 不为 nil，并要求 `CheckSandbox(..., nil, ...)` 以 fail-closed 拒绝；[tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/red-guard-contract.txt) 也显示 “NewSandbox must never return nil: a nil default would turn every check into a silent pass”。
   - 不会红：单纯 `go build ./...` 没有覆盖 fail-closed 安全语义；编译不证明“拒绝”的正确性。
   - 结论：删掉 fail-closed 返回，就是把默认实现从“安全拒绝”变成“静默放行/空路径”。

4. 把 `VerdictUnknown` 改成 `Allowed`
   - 会红：契约测试在 [internal/guard/sandbox_interface_test.go](internal/guard/sandbox_interface_test.go) 明确断言 “an Unknown verdict must not be allowed”；这条用例会立刻红。
   - 也会红：`TestAgentRunnerRunSandboxRefusesBeforePermit/unknown` 与 `TestSupervisedRunSandboxRefusesBeforePermit/unknown` 的语义会被破坏；它们原本依赖 “Unknown = fail-closed”，一改成 `Allowed`，运行路径就不再拒绝。
   - 不会红：只有“Happy path 允许 dispatch”的测试；它们并不会验证 fail-closed 性。
   - 结论：`Unknown` 是安全语义而非允许值，改成 `Allowed` 会直接破坏约束。

没有覆盖的不变量（我认为未被本次静态材料完全覆盖）：
- `PRFRAIL_SANDBOX_CHECK` 的非法值到真正 CLI 层的零副作用，只有 resolver 单元测试覆盖；未看到全链路 CLI 证据。
- `supervised` 的“非 freshStart 再入路径不检查、不消费配额”没有独立回归样本；现有 mutation 只覆盖 fresh-start 分支。
- 真实 Windows 主机上的实际行为仍未直接验证；代码和日志都证明了 S1 的平台语义，但没有在 Windows 真实环境上复现。

## Section E

- 我无法验证的部分：
  - 实机 Windows 行为：当前静态材料没有真实 Windows 运行证据；代码上只看到 `//go:build windows` 的 α 形态返回 `Unknown`，但没有真实宿主验证。
  - 真实探针配额/脚本层的跨进程竞争：材料证明了单元/白盒行为，但没有直接复现并发抢占后一个 claim 是否永远不被“悄悄放行”。

- 什么会推翻我的 PASS：
  - 任何一个检查点移到 permit 之后或 run-root marker 之后；
  - 任何实现把 `Unknown` / `ErrSandboxUnavailable` 当成 `Allowed`；
  - 删掉 [internal/guard/sandbox_unsupported.go](internal/guard/sandbox_unsupported.go) 的 fail-closed 返回，或让默认实现返回 nil；
  - `PRFRAIL_SANDBOX_CHECK` 的非法值被接受而不报 CLI usage error；
  - [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-agent-runner.txt) 与 [tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt](tmp/P74-SANDBOX-INTERFACE/03-tests/mutation-supervised.txt) 不再从绿转红；
  - 真实环境中出现“非 Windows 平台返回 Denied”或“Windows S1 返回 Denied”而不是 `Unknown`。

INDEPENDENT SCAN: PASS
