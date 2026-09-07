# T011 Review/Publish 验证报告

[English](t011-chain_EN.md)

日期：2026-09-07。结论：`T011 COMPLETE`。本任务完成候选冻结、独立评审、发布收据绑定与 task 接受接线；并在复核中修复 2 个问题后再次通过门禁。

## 复核并修复的问题

1. receipt 序列化 schema 兼容性：`ReviewReason` 与 `WaiverAuthorization` 增加 JSON 键标签，避免出现 `Code` 等大写键；`message` 为空时省略，满足 schema 的 `minLength` 约束。
2. `REPAIR_PENDING` 恢复保护：阻断 `REPAIR_PENDING/WAITING_FOR_OPERATOR` 在同 attempt 内继续推进，链路回到 `PAUSED` 并要求后续新 attempt 或人工修复流程。

## 关键实现

1. `internal/chain/review.go`：评审收据构建、JCS 摘要、approve/reject/waive 校验。
2. `internal/chain/publish.go`：发布收据构建、completed/failed/uncertain 语义与证据校验。
3. `internal/chain/engine.go`：`REVIEW_PENDING` 后改为 candidate->review->promotion 流程；reject/无效 review/不完整 promotion 进入 `REPAIR_PENDING` 并暂停。
4. `internal/chain/models.go`：新增 candidate/review/publish 端口与数据模型，补齐 review 结构字段标签。

## 测试覆盖（AT-05）

1. 拒评：`TestReviewRejectMovesTaskToRepairPending`。
2. 自批拒绝：`TestRejectSelfReview`。
3. 过期 waiver 拒绝：`TestRejectExpiredOrWrongCandidateWaiver/expired_waiver`。
4. 错候选 waiver 拒绝：`TestRejectExpiredOrWrongCandidateWaiver/wrong_candidate`。
5. 有效独立批准后单次发布：`TestValidIndependentApprovalPublishesOnce`。
6. 修复回归：
   - schema 键名/空 message：`TestBuiltReceiptsSerializeToSchemaShape`。
   - repair pending 恢复阻断：`TestRepairPendingCannotResumeWithoutNewAttempt`。

## 门禁复核

- 定向复核：
  `go test -count=1 ./internal/chain -run "TestBuiltReceiptsSerializeToSchemaShape|TestRepairPendingCannotResumeWithoutNewAttempt|TestRejectSelfReview|TestRejectExpiredOrWrongCandidateWaiver|TestValidIndependentApprovalPublishesOnce"` 通过。
- 全量门禁：`go build ./...`、`go vet ./...`、`go test -count=1 ./...` 通过。

## 边界

当前实现保证同一次 `Run()` 中的单次发布行为；跨崩溃窗口的发布幂等与 receipt store fencing 由后续任务继续收敛，不在 T011 范围内。
