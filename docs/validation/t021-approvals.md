# T021 PC-03 授权与审批验证报告

[English](t021-approvals_EN.md)

日期：2026-09-08。结论：`T021 COMPLETE`。

## 实现范围

1. `internal/chain/authorization.go`
2. `internal/chain/authorization_test.go`
3. `internal/console/approvals.go`
4. `internal/console/approvals_test.go`
5. `internal/console/cli.go`（`approvals` 命令接线与 usage）
6. `docs/validation/t021-approvals.md` / `t021-approvals_EN.md`（本报告）

## AT-18 覆盖

1. 授权绑定主体、run/manifest 摘要、工具/路径/网络范围、预算与有效期：`AuthorizationGrant` 全字段对齐 `authorization-record.schema.json`，非法主体、范围、预算或时间戳均拒绝写入。
2. 撤销/过期/错摘要阻断下一动作和接受：`EvaluateGrant` 在撤销命中时永久返回 `revoked`（同 ID 后续 grant 不得复活）；`expiresAt` 到期返回 `expired`；撤销 `authorizationHash` 与账本中任一该 ID grant 均不匹配时 fail-close 报错。
3. 在途进程停止失败保持暂停：`RequestControlledStop` 将 `stopDisposition=requested` 接线到 `Stopper.Stop`；停止失败返回 `uncertain` + `ErrRecoveryUncertain`，绝不谎报 `completed`；`approvals revoke` 在停止不确定时返回非零退出码，但撤销记录已持久化（重启可见）。
4. 硬门禁不可 waive：授权/撤销模型无任何 waiver 字段或旁路路径；`grant`/`revocation` 的 `issuedBy` 仅允许 `operator`/`policy`，`agent` 自授权被拒绝。
5. 待审批重启可见：`approvals list` 从 append-only 账本 `authorization-ledger.jsonl` 重建待审批队列（`pendingCount`、状态与 next action），账本跨进程重启保留。
6. 等待不会自动批准：`issuedAt` 之前的 grant 评估为 `pending`（`await-issue`），到期即 `expired`（`reissue`），不存在自动批准路径。

## 门禁

1. `go test -count=1 ./internal/chain ./internal/console`：通过。
2. `go build ./...`：通过。
3. `go vet ./...`：通过。
4. `go test -count=1 ./...`：通过（全部包 ok）。
5. 编辑器诊断：`internal/chain/authorization.go`、`internal/console/approvals.go`、`internal/console/cli.go` 均无错误。

## 关键命令证据

```powershell
gofmt -w internal/chain/authorization.go internal/chain/authorization_test.go internal/console/approvals.go internal/console/approvals_test.go internal/console/cli.go
go test -count=1 ./internal/chain ./internal/console
go build ./...
go vet ./...
go test -count=1 ./...
```

## 边界

1. `approvals` 仅提供 `list`/`revoke` 子命令；grant 记录由可信策略签发流程写入账本，CLI 不提供 `grant` 写入入口（防"聊天文本/仓库文件提升为授权"）。
2. 停机接线基于 `chain.Stopper` 抽象：`approvals revoke` 默认通过 `internal/console/stopper.go` 解析 `managed-process.identity.json` 并调用 guard 停机；本地 noop 终态 run 返回终态证据摘要，缺失进程身份时 fail-close 为 `uncertain`。
3. `not-required` 撤销不调用 stopper；`requested` 成功记录 `completed`，失败记录 `uncertain` 并要求 `reconcile-stop`。
4. 账本为 append-only JSONL：每行一条 canonical 授权记录，解析失败按行号报错；不做原地修改或删除。
5. 有效期边界按严格时间戳比较（`expiresAt` 必须晚于 `issuedAt`，签发前评估为 `pending`）。

## 复审记录（2026-09-08）

1. `EvaluateGrant` 撤销匹配修正为"与任一该 authorizationId 的 grant 记录哈希匹配"（同 ID 多 grant 场景），不匹配仍 fail-close；修正对应测试 `TestEvaluateGrantRevocationPermanentAndHashMismatchFailsClosed`。
2. `newDefaultStopRun` 对受控停机补充硬超时上下文（`defaultStopTimeout=12s`，`defaultStopGrace=2s`），避免 revoke 在停机确认阶段潜在无限等待；新增测试 `TestDefaultStopRunAppliesBoundedStopTimeout`。
3. `BuildApprovalsReport` 改为按 `authorizationId` 最新 grant 聚合，`pendingCount`/`revokedCount` 反映当前审批收件箱而非历史 grant 条数；新增测试 `TestApprovalsReportUsesLatestGrantPerAuthorization`。
4. 门禁复跑全部通过。
