# T015 Taskdef Documentation 验证报告

[English](t015-taskdef-documentation_EN.md)

日期：2026-09-08。结论：`T015 COMPLETE`。本任务交付 `harnesses/{generic,c,go}/` 与 `internal/taskdef/documentation.go`，实现模板 A 锁定、impact rule 文档义务检查、生成物新鲜度与手改阻断，以及 C+Go 同 task 聚合恢复判定。

## 实现范围

1. `internal/taskdef/documentation.go`
   - 新增文档策略与 impact rule 检查：`required/if-affected/optional/forbidden`。
   - 新增模板 A 锁定：`template-a@1.0.0` 与锁定摘要绑定。
   - 新增生成物门禁：阻断 `generated.self-approval`、`generated.manual-edit`、`generated.stale`、`generated.actor-missing`、模板锁不匹配。
   - 新增多语言 task 聚合判定：任一 language scope 失败即 `RequiresRecovery=true`。
2. `internal/taskdef/documentation_test.go`
   - 覆盖 impact rule 命中与缺失文档阻断。
   - 覆盖模板 A 锁、手改生成物、自批、陈旧生成物阻断。
   - 覆盖 C+Go 作用域整任务恢复语义。
   - 覆盖 `harnesses/generic|c|go/harness.json` 的存在性与模板参数默认值。
3. `harnesses/{generic,c,go}/harness.json`
   - 新增三套 S1 harness 定义（generic/C/Go）。
   - 各 harness 统一声明 `generated-hook-template=template-a@1.0.0`、`allow-generated-hooks=false`。

## AT-11 覆盖

1. generic/C/Go：三套 harness 已外置并可被测试读取校验。
2. 锁模板生成 A：生成物仅允许 `template-a@1.0.0` 且必须匹配锁定摘要。
3. 阻断条件：
   - 代码命中 impact rule 但文档目标未更新。
   - required hook 未通过。
   - 生成物被手改、陈旧或自批。
   - 生成物缺少生成者或批准者 actor。
4. C+Go 同 task 聚合：任一 language scope 失败，整 task 进入恢复流程。

## 测试用例

- `internal/taskdef/documentation_test.go`
  - `TestCheckDocumentationIfAffectedPassesWithFreshGeneratedArtifacts`
  - `TestCheckDocumentationBlocksMissingDocumentationAndHook`
  - `TestCheckDocumentationBlocksGeneratedViolations`
  - `TestCheckDocumentationRequiresImpactRulesWhenPolicyNeedsDocumentation`
  - `TestEvaluateTaskRecoveryForCAndGoScopes`
  - `TestHarnessProfilesProvideGenericCGoAndTemplateALock`

## 门禁

1. `go test -count=1 ./internal/taskdef`：通过。
2. `go build ./...`：通过。
3. `go vet ./...`：通过。
4. `go test -count=1 ./...`：通过。
5. `2026-09-08` 远程 race 高次数回归（Ubuntu VM `10.0.0.199`，Go 1.27.1 + gcc 13.3.0）：`go test -race -count=20 ./internal/applier ./internal/console ./internal/taskdef` 与 `go test -race -count=1 ./...` 通过（首轮通过，未触发代理兜底）。

## 边界

1. 本切片只交付模板化生成场景 A 的静态约束，不引入任意生成脚本 B。
2. 本切片未引入 whois 专名或 whois 运行语义。
3. 本轮未提交、未推送。
