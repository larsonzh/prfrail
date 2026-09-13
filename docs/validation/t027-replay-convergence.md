# T027 replay 收敛与发布耐久性前置验证报告

[English](t027-replay-convergence_EN.md)

日期：2026-09-14。结论：`PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`。

## 验证范围

本报告只覆盖已完成的离线 replay convergence/publication durability 前置切片，不宣称 AgentRunner 生产实现。effect mapping v1 已冻结：`local-process` 与 `workspace-write` 均映射到 `local-discardable`，request 绑定 mapping version/hash，缺失/未知/篡改绑定、未知 operation 和 grant class 覆盖不足均在 replay identity 消费或 dispatch 前 fail-closed。

`AgentRunnerReplayStore` 已覆盖：

- orphan completion 与磁盘 corruption 的 fail-closed 处理；
- no-replace 碰撞后的 bounded reread convergence；
- 同一 root 上多个 store 并发 request 的 single first-dispatch，其余调用阻断；
- Windows publication durability `unproven` 时的 pre-write 阻断，不创建 request record。

## 最终 verifier 结果

### Windows 本机执行

以下命令在 Windows 上执行：

| 命令 | 结果 |
|---|---|
| `node tools/contracts/contracts.test.js` | Node contract 4/4 通过 |
| `go test -count=1 ./internal/adapters -run '^TestAgentRunnerReplayStore'` | replay 32/32 项通过 |
| `go test -count=1 ./internal/adapters` | adapters 通过 |
| `go build ./...` | 通过 |
| `go vet ./...` | 通过 |
| `go test -count=1 ./...` | 全部通过 |

### Unix 目标：仅交叉编译

| 命令 | 结果 |
|---|---|
| `$env:GOOS='linux'; $env:GOARCH='amd64'; go build ./...` | Linux/amd64 cross-compile 通过 |
| `$env:GOOS='darwin'; $env:GOARCH='amd64'; go build ./...` | Darwin/amd64 cross-compile 通过 |

Linux/Darwin 命令只证明目标构建可完成，没有在对应系统执行 replay、fsync、崩溃注入或测试；不能把 cross-compile 记为 Unix runtime/platform PASS。

## 发布耐久性边界

Windows 可以提供 no-replace 的可见性，但当前不能证明新目录项持久化。因此任何可能返回 `first-dispatch` 的写操作在 pre-write 阶段拒绝；不会因为已看到文件就宣称 durable dispatch。

Unix 的 `proven` 路径要求调用方提供已存在且稳定的绝对 replay root，并在创建/确认 `requests/`、`completions/` 后成功 fsync parent directory。最终 verifier 在 Windows 执行，本轮没有 Unix 原生 crash-injection evidence；上述 Unix 条件是契约边界，不是本轮实机证明。

terminal receipt（包括 `uncertain`）只说明终局收据存在并阻止重复 dispatch，不是 task `PASS`，也不是 AT-23 通过证据。replay store 尚未接入 `AgentRunnerPort` 或 Engine。

## 明确未执行事项

- 未调用真实模型 CLI；
- 未访问网络；
- 未读取 SecretStore；
- 未执行 commit、push 或 publish；
- T027 仍为 `BLOCKED / NOT IMPLEMENTED`，AT-23 未通过。
