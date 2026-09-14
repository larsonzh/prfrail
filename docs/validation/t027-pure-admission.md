# T027 A1 pure admission 前置验证报告

[English](t027-pure-admission_EN.md)

日期：2026-09-14。结论：`PREREQUISITE COMPLETE / T027 BLOCKED / NOT IMPLEMENTED / AT-23 NOT PASS`。

## 验证范围

本报告只覆盖 T027 离线切片 A1 pure admission，不宣称 AgentRunner 生产实现。A1 已交付：

- `AgentRunnerCompositeAdmission` 删除 `RequestIndex` 字段与内存 replay 消费：admission 只做纯 preflight 与不可变绑定，无 replay 写入副作用，也不触发任何 launch。
- 删除仅被该消费使用的一次性 replay 阻断哨兵 `ErrAgentRunnerRequestReplayBlocked`，以及内存 `AgentRunnerRequestIndex` 类型、构造器与 `Record` 方法。同键同摘要幂等、同键异摘要冲突与 first-dispatch 语义仍由 A0 的耐久 replay store 在 dispatch 路径承担；`ErrAgentRunnerRequestConflict` 保留且仍被 replay store 测试覆盖。
- 责任边界固化：`docs/CONTRACTS.md` / `CONTRACTS_EN.md` 新增"生产 admission 是纯 preflight 与不可变绑定：不得消费 request replay 身份、不得发布 R/C 或写 replay store、不得启动进程；首次 dispatch 资格与 replay 身份消费由 dispatch 路径在发布 R 之后依据 replay store 判定；admission 结论是评估时刻 verdict 而非锁，dispatch 必须在发布 R 后、允许 launch 前重新确认授权未被撤销且预算 reservation 仍 outstanding"；`internal/chain/router.go` 的 `AgentRunnerAdmission` 接口文档与 `agent_runner_admission.go` 的结构/方法文档同步声明同一边界。
- 聚焦测试：`TestAgentRunnerCompositeAdmissionIsStatelessPreflight`（同一 admission 连续 3 次通过；同 requestId 异 digest 的冲突判定归 dispatch/store 边界而非 admission）；`...FailureLeavesNoConsumedState`（无效记录与缺失时钟阻断后，同一实例修复条件即通过）；`...ConcurrentRepeatsPass`（8 个并发重复 admission 全部通过）；`...LeavesInputsUnchanged`（成功与阻断两条路径前后对全部输入做深度相等断言，把"无写入副作用"升级为回归断言）；`...RejectsNilReceiverCancelledContextAndZeroClock`（nil receiver、已取消 ctx、零时间 clock 三个 fail-closed 分支）；删除内存索引专属用例并清理 `BeforeReplay` 旧命名。

## 执行结果

### Windows 本机执行（原生证据）

| 命令 | 结果 |
|---|---|
| `gofmt -l internal` | 干净（无输出） |
| `go build ./...` | 通过 |
| `go vet ./...` | 通过 |
| `go test -count=1 ./...` | 全部通过 |
| `go test -count=1 -run 'AgentRunnerCompositeAdmission' ./internal/adapters` | 通过 |

说明：Windows 本机无 gcc（cgo 不可用），未在本机运行 `-race`；Linux CI 已新增 `-race` 步骤（见下），推送后生效。并发正确性依据新增并发用例的常规执行，以及 `CostLedger.RequireOutstandingReservation` 的内部 `RWMutex` 读锁与 admission 全程只读的事实。

### CI 工作流变更（冻结契约同步）

- 新增 Linux CI 步骤 `Race`：`CGO_ENABLED=1 go test -race -count=1 ./internal/adapters/... ./internal/chain/...`，仅在 Runner 为 Linux 时执行。
- `internal/release` 的工作流冻结契约同步更新：`test` 作业步骤序列加入 `Race`，脚本摘要与条件写入冻结表；本地 `go test -count=1 ./...` 全绿（含 release 工作流契约与变异拒绝套件）。

### 审查结论

- V4 Pro 前置分析（补做）：A1 本体 PASS（实现与理想裁定逐项吻合）；提出 R1（契约补句，A2 前置）与 R2–R4 收尾项，已全部整改：R1 在 CONTRACTS 双语冻结"admission 结论非锁、dispatch 须在发布 R 后、launch 前重新确认授权未撤销与预算 reservation 仍 outstanding"；R2 补三例 fail-closed 分支测试；R3 补冲突哨兵归属注释；R4 补不变性 helper 判据注释；Codex 整改复审 PASS。
- V4 Pro 预审（补做）：PASS（无 Critical/High/Medium）。四项 Low 建议全部采纳：L1 Linux CI 新增 `-race` 门禁（工作流冻结契约同步更新）；L2 输入不变性回归测试；L3 A2 交接说明（admission 是 preflight 而非锁，dispatch 必须以发布 R 后的 replay store 判定为准）；L4 契约措辞精确化为"不得发布 R/C 或写 replay store"。

- Codex 独立终审：PASS（无 Critical/High/Medium）。全部 Low 建议已采纳：清理 `BeforeReplay` 命名、补"早失败无副作用"与"并发 admission"反例；"无 launch" 观测由既有 router "admission 阻断则不 dispatch" 测试承担。后续变动仅测试命名与新增用例、未改生产代码，因此未触发强制复审。

## 已知边界

- admission 的"无 launch"是结构性保证：`AgentRunnerCompositeAdmission` 不持有任何 port、进程或 replay 依赖；router 层已有 admission 阻断时不调用 `AgentRunnerPort` 的测试。
- replay 身份消费与 first-dispatch 资格尚未实现，属 A2 replay-aware dispatcher；在此之前不存在生产 dispatch 路径，因此本切片删除 admission 阻断不产生现实 redispatch 暴露。
- 本机 Windows 无 cgo/gcc，未在本机执行 `-race`；Linux CI 已固化 `-race` 步骤，推送后由 GitHub Actions 提供竞态回归证据。
- admission 是 preflight 而非锁：verdict 可能在 dispatch 前过期（授权撤销、预算结算）；A2 dispatch 必须以发布 R 后的 replay store 判定为准。

## 明确未执行事项

- 未引入真实 dispatch，未接入 `AgentRunnerPort` 或 Engine；
- 未调用真实模型 CLI，未访问网络，未读取 SecretStore；
- 本轮未执行 commit、push 或 publish；
- T027 仍为 `BLOCKED / NOT IMPLEMENTED`，AT-23 未通过。
