# T019 PC-01 预览验证报告

[English](t019-preview_EN.md)

日期：2026-09-08。结论：`T019 COMPLETE`。

## 实现范围

1. `internal/taskdef/preview.go`
2. `internal/console/preview.go`
3. `internal/taskdef/preview_test.go`
4. `internal/console/cli_test.go`（补充 preview 用例）
5. `examples/no-ai-preview/README.md`

## AT-16 覆盖

1. 无 AI 静态预览：`mode=offline-read-only`，不执行 steps/hooks。
2. 零调用计数：command/network/model/credential/versionProbe 全为 0。
3. unknown 可见：明确 reason 与 nextAction，不伪造探测结果。
4. 源与 store 不变：preview 仅读配置并生成摘要，不创建 run 工作目录。
5. JSON 与文本一致：两种输出共享同一 `previewHash` 与 outcome 事实。

## 门禁

1. `go test ./internal/taskdef ./internal/console`：通过。
2. `go build ./...`：通过。
3. `go vet ./...`：通过。
4. `go test -count=1 ./...`：通过。

## 关键命令证据

```powershell
go test ./internal/taskdef ./internal/console
go build ./...
go vet ./...
go test -count=1 ./...
```

## 边界

1. 不做命令探测、网络探测、模型调用或自动安装。
2. 对可执行 step 仅给出 `blocked` 预览，不触发 run。
3. 真实授权、导出与生命周期流程属于 T020+。
4. `reviewRequired` 恒为 `true`：预览采用保守口径，当前不表达免审步骤。
5. `blockingGateIds` 取自 step 的全部 `hooks`，不区分阻塞/非阻塞 gates。

## 复审记录（2026-09-08）

1. 对齐 `plan-preview.schema.json`：`networkRequirements[].host` 增加 `^[A-Za-z0-9.-]+$` 与长度 253 校验；`budget.pricingSource` 增加 1–256 长度校验；补充 `TestPlanPreviewRecordRejectsInvalidNetworkHost`、`TestPlanPreviewRecordRejectsLongPricingSource` 拒绝用例。
2. 门禁复跑全部通过。
