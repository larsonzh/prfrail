# T012 Tickets/Repair 验证报告

[English](t012-tickets-repair_EN.md)

日期：2026-09-07。结论：`T012 COMPLETE`。本任务实现失败指纹与 ticket ledger、预算闸门、repair transaction 四阶段链，并完成 AT-07 关键场景验证。

## 实现范围

1. `internal/tickets/fingerprint.go`：失败指纹输入校验与稳定摘要（`proofrail:failure-fingerprint:1`）。
2. `internal/tickets/ledger.go`：append-only ledger（failure/override-granted/resolved）、派生 `currentBudgetState`、`ticketLedgerHash` 计算（`proofrail:ticket-ledger:1`）。
3. `internal/tickets/budget.go`：attempt 预约与预算门禁（hard-block、review gate、override window、墙钟/费用耗尽）。
4. `internal/repair/transaction.go`：Prepare/Inspect/Validate/Promote 顺序阶段链、相邻 `stageHash` 绑定（`proofrail:repair-stage:1`）、transaction 摘要（`proofrail:repair-transaction:1`）。

## AT-07 场景验证

1. 相同指纹耗尽：failure 计数达到阈值后进入 hard-block，resolved 不重置预算。
2. attempt 耗尽：override 授权窗口内超过 `maxAdditionalAttempts` 拒绝；同 attempt 重复消费拒绝。
3. 墙钟/费用耗尽：`ReserveAttempt` 在 wall-clock 或 cost 耗尽时立即阻断。
4. 假修复阻断：Promote 阶段 `validateStageHash`/`ticketLedgerHash` 绑定不一致拒绝。
5. 陈旧候选阻断：Inspect/Validate/Promote 的 candidate hash 与 Prepare 产物不一致时拒绝。
6. Promote 中断：`outcome=uncertain` 进入终态，后续阶段追加被拒绝。

## 测试用例

- `internal/tickets/ledger_test.go`
- `internal/repair/transaction_test.go`
- `internal/repair/transaction_fingerprint_budget_test.go`

## 门禁

- `go build ./...` 通过。
- `go vet ./...` 通过。
- `go test -count=1 ./...` 通过。

## 审查修复

2026-09-07 复核发现并修复下列 schema/语义合规问题：

1. 阶段序列化字段裁剪：新增 `Stage.MarshalJSON`，按阶段名只输出该阶段字段，并保证 prepare failed/uncertain 的 `candidateManifestHash` 与 promote failed/uncertain 的 `resultingManifestHash` 输出 `null`（此前被 omitempty 省略，违反 schema required）。
2. 证据切片非 nil：`uniqueHashes` 对空输入返回 `[]`，避免 nil 序列化为 `null`（schema 要求 array）。
3. 阶段 evidence 非空：Prepare/Inspect/Validate/Promote 的 `evidence` 至少一项（nonemptyHashSet）。
4. `targetIds` 唯一性：重复目标在写入前拒绝（idSet uniqueItems）。
5. 账本 ticketId 唯一：重复 ticketId 拒绝。
6. 预算状态派生：override 已耗尽时 `currentBudgetState` 回落 `pending-review`（契约规则 2 要求“未耗尽”）。

## 边界

本切片提供核心 ledger/repair 语义与本地门禁，不包含跨 run 成本计量聚合与 adapter 结算对账（后续 T023/T013 继续接入）。
