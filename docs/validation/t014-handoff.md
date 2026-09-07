# T014 Handoff 验证报告

[English](t014-handoff_EN.md)

日期：2026-09-07。结论：`T014 COMPLETE`。本任务完成 manual-handoff 的链路核心建模与输入策略守卫，覆盖 handoff receipt 语义、WAITING_FOR_OPERATOR 状态衔接，以及 secret-direct 的终端安全输入约束。

## 实现范围

1. `internal/chain/handoff.go`
   - 新增 `HandoffPolicy` 校验：`allowedTargets/inputPolicy/handoffTimeoutMs/returnActions/hooksAfterReturn` 全量约束。
   - 新增 `BuildHandoffReceiptRecord`：生成并校验 handoff receipt，绑定 operator、policy hash、manifest/diff、lease/hook/evidence，并输出 canonical receipt hash。
   - 新增 `OpenHandoffTransition` 与 `ReturnHandoffTransition`：显式建模 `WAITING_FOR_OPERATOR` 的开闭环。
2. `internal/chain/models.go`、`internal/chain/router.go`
   - `code` step 新增 `manual-handoff` 模式与 `handoffPolicy` 绑定校验。
   - Router 新增 manual-handoff 专用端口，避免与 managed-change-set / isolated-workspace / hooks 混流。
3. `internal/console/handoff.go`
   - 新增 handoff 输入策略守卫：`structured` 与 `secret-direct`。
   - secret-direct 强制 `terminal` 通道且要求安全输入能力，否则 fail-close 返回暂停信号。
   - 秘密值只输出摘要证据（hash），不进入结构化输出。
4. 测试
   - 新增 `internal/chain/handoff_test.go`、`internal/console/handoff_test.go`。
   - 扩展 `internal/chain/state_test.go`、`internal/chain/router_test.go` 覆盖 manual-handoff 路由与定义校验。

## AT-09 覆盖

1. open 过渡仅允许 `STEPS_RUNNING/RUNNING -> WAITING_FOR_OPERATOR/WAITING_FOR_OPERATOR`。
2. return 过渡覆盖三种归还动作：
   - `complete -> STEPS_RUNNING/RUNNING`
   - `abort -> FAILED/FAILED`
   - `request-agent -> FAILED/FAILED`（进入失败态后由修复流程建新 attempt 并注入经确认结论；不映射 CANCELLED 以免毒化整个 run）
3. handoff receipt 约束：
   - `recordedBy=system`、`operator=operator`、时间戳合法且 `closedAt > openedAt`。
   - `complete` 必须提供非空 `hookResultEvidence`。
   - `abort/request-agent` 必须 `hookResultEvidence=[]`。
4. `AllowedReturnActions` 与 `outcome` 绑定：不在允许集合内即拒绝。

## AT-10 覆盖

1. 未声明 prompt（额外字段）立即拒绝。
2. `secret-direct` 在非终端通道或无安全输入能力时返回 `ErrSecureInputRequired`（暂停而非降级）。
3. secret 输入仅生成摘要证据，不出现在结构化输出中（测试使用 canary 校验不泄漏）。

## 测试用例

- `internal/chain/handoff_test.go`
  - `TestBuildHandoffReceiptRecordComplete`
  - `TestBuildHandoffReceiptRecordRejectsOutcomeRules`
  - `TestBuildHandoffReceiptRecordRejectsDisallowedOutcome`
  - `TestBuildHandoffReceiptRecordRejectsEmptyAllowedReturnActions`
  - `TestBuiltHandoffReceiptSerializesToSchemaShape`
  - `TestHandoffTransitions`
- `internal/chain/state_test.go`
  - `TestDefinitionValidationFourKindsAndModes`（新增 manual-handoff 正反校验）
- `internal/chain/router_test.go`
  - `TestStepRouterSeparatesChangeModesAndHooks`（新增 manual-handoff 路由断言）
- `internal/console/handoff_test.go`
  - `TestPrepareHandoffInputStructured`
  - `TestPrepareHandoffInputRejectsUndeclaredPrompt`
  - `TestPrepareHandoffInputRejectsUnknownChannel`
  - `TestPrepareHandoffInputSecretDirectRequiresSecureTerminal`
  - `TestPrepareHandoffInputSecretDirectDoesNotExposeSecret`

## 门禁

1. `go test -count=1 ./internal/chain ./internal/console`：通过。
2. `go build ./...`：通过。
3. `go vet ./...`：通过。
4. `go test -count=1 ./...`：通过。
5. Ubuntu VM（2026-09-08，临时代理链路）远端全量验证：
   - `go build ./...`：通过。
   - `go vet ./...`：通过。
   - `go test -count=1 ./...`：通过。
   - `CGO_ENABLED=1 CC=gcc go test -count=1 -race ./...`：通过。
   - `CGO_ENABLED=1 CC=clang go test -count=1 -race ./internal/chain ./internal/console`：通过。
   - 结果汇总：`RESULT build=0 vet=0 test=0 race_gcc=0 race_clang=0 race_all=0`。
   - 验证完成后已清理远端临时代理包装脚本与会话变量，避免代理配置残留。

## 审查修复

2026-09-08 审查发现并修复以下问题：

1. request-agent 归还映射错误：此前映射为 `CANCELLED/CANCELLED`，但 CANCELLED 是终态且会毒化整个 run；契约规定 request-agent 建新 attempt 继续推进。已改为 `FAILED/FAILED`（失败态是状态机中进入修复流程建新 attempt 的唯一前置），并同步更新测试与本文档。
2. AllowedReturnActions 空集合静默跳过 outcome 绑定校验（fail-open）：契约要求 outcome 必须落在 `returnActions` 集合内，且 schema `minItems: 1`。现改为会话必填，空集合或非法动作集直接拒绝，并新增 `TestBuildHandoffReceiptRecordRejectsEmptyAllowedReturnActions`。
3. 输入通道枚举未校验：未知 channel 值在 structured 下会静默通过。现 fail-close 拒绝非 terminal/model 通道，新增 `TestPrepareHandoffInputRejectsUnknownChannel`。
4. 补充回归测试：新增 `TestBuiltHandoffReceiptSerializesToSchemaShape`（abort 时 `hookResultEvidence` 必须序列化为 `[]` 而非 `null`、键名必须 camelCase）；`state_test` 新增 manual-handoff 无 policy 与 managed 带 policy 两组负例。

## 边界

1. 本切片未引入 CLI 命令层 handoff 交互流程；当前交付聚焦链路模型、收据约束与输入安全守卫。
2. AT-09 中依赖运行时的租约冲突、操作员离席/断线、崩溃重放与归还后重跑 gates 由后续引擎接线（guard/lease + gates）承接；AT-10 的真实终端安全输入能力由 T016 CLI/TUI 承接。
3. 本轮为离线单元/门禁验证，授权宿主的交互式 smoke 保持后续计划执行。