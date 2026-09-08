# T023 PC-05 成本账本验证报告

[English](t023-cost-ledger_EN.md)

日期：2026-09-08。结论：`T023 COMPLETE`。

## 实现范围

1. `internal/tickets/cost.go`：实现成本账本模型、append-only 记录、reservation/settlement 规则、summary 重放校验与共享上限判定。
2. `internal/tickets/cost_test.go`：覆盖跨 run 上限、unknown hold 重启保留、settlement 幂等、reservation 生命周期与 summary 漂移阻断。
3. `internal/tickets/budget.go`：新增 `CostBudgetExhausted` 与 `CostTokenBudgetExhausted` 以支持费用/Token fail-close 判定。
4. `internal/adapters/cost_usage.go` 与 `internal/adapters/cost_usage_test.go`：从 durable result/dispatch receipt 提取可结算 usage，并生成稳定 idempotency key。
5. `internal/console/cost.go` 与 `internal/console/cost_test.go`：实现 `prfrail cost report` 的 JSON/文本输出、无默认遥测、订阅模式账单声明。
6. `internal/console/cli.go`：接入 `cost` 子命令路由。
7. `docs/OPERATIONS.md` / `docs/OPERATIONS_EN.md`：补齐可执行命令与 PC-05 操作说明，消除中英状态漂移。
8. `docs/validation/t023-cost-ledger.md` / `docs/validation/t023-cost-ledger_EN.md`：本次验证报告。

## AT-20 覆盖

1. 调用前持久 reservation：`Reserve` 在账本中先追加 reservation 再进入后续结算路径；`TestCostLedgerCrashPointReservationLifecycle` 验证 reservation 在结算前保持占用。
2. settlement 去重与冲突阻断：`Settle` 按 `idempotencyKey` 幂等去重，同 key 变更 payload 返回冲突；`TestCostLedgerSettlementDeduplicates` 覆盖重复结算与冲突分支。
3. unknown hold 重启保留：`Snapshot` + `LoadCostLedger` 重放后保留 unknown hold；`TestCostLedgerUnknownHoldSurvivesRestart` 验证重启后仍占用并继续阻断超限申请。
4. 共享上限跨 run 生效：`validateSharedCaps` 使用 settled + outstanding 聚合判定；`TestCostLedgerSharedCapAcrossRuns` 覆盖跨 run 累计超限阻断。
5. 预算 fail-close：`CostBudgetExhausted` 在 committed 货币未知时直接视为耗尽，`CostTokenBudgetExhausted` 以 settled + hold token 判定；`TestCostBudgetExhaustionHelpers` 覆盖边界。
6. 本地报表无默认遥测：`BuildCostReport` 固定 `TelemetryExported=false`，并给出本地账本说明；`TestCostReportJSONAndText`、`TestCostReportSubscriptionDisclaimer`、`TestCostReportMissingLedgerIsNonFatal` 覆盖输出事实与免责声明。
7. usage 证据可追溯：`UsageFromResultEnvelope`/`UsageFromDispatchReceipt` 把 request/receipt/result 证据整合进 settlement-ready 结构；对应测试验证证据完整与 key 稳定。

## 门禁结果

1. `go test -count=1 ./internal/tickets ./internal/adapters ./internal/console`：通过。
2. `go build ./...`：通过。
3. `go vet ./...`：通过。
4. `go test -count=1 ./...`：通过（全部包 ok）。

## 关键命令证据

```powershell
go test -count=1 ./internal/tickets ./internal/adapters ./internal/console
go build ./...
go vet ./...
go test -count=1 ./...
```

## 边界

1. T023 仅实现本地账本与成本可观测性，不直接对接供应商账单 API；未知费用保持 hold；同一 reservation 一旦按 unknown 结算即不可再次结算，后续对账需进入新的可验证记账周期。
2. 订阅模式明确只约束授权调用与 token 上限，不承诺货币账单精确一致。
3. `cost report` 不会创建或外发遥测数据；缺失账本时返回 warning 与空报表，不自动生成历史事实。
4. 本次交付不包含 T024 生命周期处置、备份回放与发布清单能力。

## 复审记录（2026-09-08）

1. 复核确认 `docs/OPERATIONS_EN.md` 中“T023 pending”表述已改为“已实现 T023，T024 待实现”，与中文一致。
2. 全量门禁复跑通过，未发现新增编译、vet 或测试回归。